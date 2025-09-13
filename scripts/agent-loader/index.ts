#!/usr/bin/env node

/**
 * BMAD Agent Loader - Story 1.2
 *
 * Converts BMAD .md agent files into OpenCode JSON configuration.
 * Usage: bun run generate-config --agents-dir <path> --output <path>
 */

import { glob } from 'glob';
import matter from 'gray-matter';
import { readFileSync, writeFileSync } from 'fs';
import { resolve, relative, basename } from 'path';

interface AgentConfig {
  description?: string;
  mode?: 'primary' | 'subagent' | 'all';
  model?: string;
  temperature?: number;
  tools?: Record<string, boolean>;
  prompt?: string;
}

interface OpenCodeConfig {
  agent: Record<string, AgentConfig>;
  $schema?: string;
}

interface CliArgs {
  agentsDir: string;
  output: string;
}

function parseArgs(): CliArgs {
  const args = process.argv.slice(2);
  const agentsDirIndex = args.indexOf('--agents-dir');
  const outputIndex = args.indexOf('--output');

  if (agentsDirIndex === -1 || outputIndex === -1) {
    console.error('Usage: bun run generate-config --agents-dir <path> --output <path>');
    process.exit(1);
  }

  return {
    agentsDir: args[agentsDirIndex + 1],
    output: args[outputIndex + 1]
  };
}

async function scanAgentFiles(agentsDir: string): Promise<string[]> {
  const pattern = resolve(agentsDir, '**/*.md');
  console.log(`🔍 Scanning for agent files: ${pattern}`);

  const files = await glob(pattern);
  console.log(`📄 Found ${files.length} agent files`);

  return files;
}

function parseAgentFile(filePath: string): { name: string; config: AgentConfig } {
  const content = readFileSync(filePath, 'utf-8');
  const { data, content: markdownBody } = matter(content);

  // Agent name from filename (without .md extension)
  const name = basename(filePath, '.md');

  // Build agent configuration
  const config: AgentConfig = {
    description: data.description || `Agent: ${name}`,
    mode: data.mode || 'subagent',
    model: data.model || 'anthropic/claude-sonnet-4-20250514',
    temperature: data.temperature || 0.1,
    tools: data.tools || {
      write: true,
      edit: true,
      bash: true
    }
  };

  // Use {file:} reference for prompt instead of embedding content
  // This keeps the .md file as the source of truth
  // Make path relative to the output directory
  const outputDir = resolve(parseArgs().output).replace(/\/[^\/]+$/, ''); // Remove filename
  const relativePath = relative(outputDir, filePath);
  config.prompt = `{file:${relativePath}}`;

  console.log(`✅ Parsed agent: ${name}`);
  return { name, config };
}

function generateOpenCodeConfig(agents: Array<{ name: string; config: AgentConfig }>): OpenCodeConfig {
  const config: OpenCodeConfig = {
    agent: {},
    $schema: "https://opencode.ai/config.json"
  };

  for (const { name, config: agentConfig } of agents) {
    config.agent[name] = agentConfig;
  }

  return config;
}

async function main() {
  console.log('🏗️  BMAD Agent Loader - Story 1.2');
  console.log('   Converting BMAD .md agents to OpenCode JSON configuration\n');

  const { agentsDir, output } = parseArgs();

  try {
    // 1. Scan for agent files
    const agentFiles = await scanAgentFiles(agentsDir);

    if (agentFiles.length === 0) {
      console.warn('⚠️  No agent files found');
      return;
    }

    // 2. Parse each agent file
    const agents = agentFiles.map(parseAgentFile);

    // 3. Generate OpenCode configuration
    const config = generateOpenCodeConfig(agents);

    // 4. Write output file
    const outputPath = resolve(output);
    writeFileSync(outputPath, JSON.stringify(config, null, 2));

    console.log(`\n🎉 Generated OpenCode configuration:`);
    console.log(`   📁 Output: ${outputPath}`);
    console.log(`   🤖 Agents: ${agents.length}`);
    console.log(`   📋 Agent names: ${agents.map(a => a.name).join(', ')}`);

  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

// Check if this is the main module
if (process.argv[1] === import.meta.url.replace('file://', '')) {
  main().catch(console.error);
}