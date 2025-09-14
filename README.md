# BMAD OpenCode Extensions

Extensions for integrating the BMAD Method with the official `sst/opencode` application.

## Overview

This repository contains the foundational components for the BMAD OpenCode Engine:

- **Agent Loader**: TypeScript script that converts BMAD `.md` agent files into `opencode.json` configuration
- **Workflow Engine**: Go server plugin that executes BMAD YAML workflows within OpenCode
- **MMOA Implementation**: Multi-Model Orchestration Architecture planning document with 7 core task categories and cost-optimized model selection

## Project Structure

```
/bmad-opencode-extensions
├── /packages
│   └── /workflow-engine         # Go-based server plugin
├── /scripts
│   └── /agent-loader           # TypeScript config generator
├── /bmad-core                  # BMAD core definitions (submodule)
├── /example-project            # Sample project for testing
├── grok--mmao-implementation-planning.md  # MMOA architecture document
├── package.json                # TypeScript dependencies
└── go.mod                     # Go dependencies
```

## MMOA (Multi-Model Orchestration Architecture)

The repository includes comprehensive planning for intelligent AI model selection and cost optimization:

### Key Features
- **7 Core Task Categories**: Streamlined from 21 to fundamental categories (Simple/Quick, General Development, Coding/Implementation, Analysis/Research, Creative/Content, Complex Reasoning, Specialized)
- **Legacy Model Primaries**: `claude-sonnet-4` for development/coding, `gemini-2.5-pro` for analysis/research
- **Cost-Effective Backups**: Full fallback chains with verified pricing from official providers
- **Image Generation Support**: Multi-provider options for visual content creation

### Model Coverage
- **Anthropic**: Claude models (legacy + modern) for reliable performance
- **Google**: Gemini models for large context and analysis tasks
- **OpenAI**: GPT models for cost-effective and specialized tasks

### Implementation Benefits
- **60-70% Cost Reduction**: Intelligent model routing based on task characteristics
- **Proven Reliability**: Legacy models for core tasks, modern models for optimization
- **Comprehensive Coverage**: All major AI providers with automatic fallbacks

See [`grok--mmao-implementation-planning.md`](grok--mmao-implementation-planning.md) for complete technical specifications and rollout strategy.

## Prerequisites

- Node.js 18+ with Bun runtime
- Go 1.21+
- `sst/opencode` binary installed globally
- Git configured for development

## Quick Start

1. **Install Dependencies**
   ```bash
   bun install
   go mod tidy
   ```

2. **Generate Config** (Story 1.2)
   ```bash
   bun run generate-config --agents-dir ./bmad-core --output ./example-project/opencode.json
   ```

3. **Build Workflow Engine** (Story 1.3)
   ```bash
   bun run build
   ```

## Development Status

- [x] Repository setup and structure
- [x] **Story 1.1: Plugin development environment + Hello World test** ✅
  - OpenCode successfully installed and configured
  - Created working `opencode.json` agent configuration
  - Created `bmad-ping` test agent (both JSON and Markdown formats)
  - **Verified**: `opencode run "@bmad-ping Test integration"` → **"Pong! BMAD OpenCode Engine is working!"**
  - Proof of concept complete: BMAD agents integrate successfully with OpenCode
- [x] **Story 1.2: BMAD Agent Loader implementation** ✅
  - TypeScript script scans BMAD `.md` agent files using glob patterns
  - Parses YAML frontmatter with gray-matter library
  - Generates valid `opencode.json` configuration with `{file:}` references
  - **Command**: `npm run generate-config --agents-dir <path> --output <path>`
  - **Verified**: Generated config loads 3 sample agents successfully
  - **Test**: `@growth-marketer` and `@architect` agents respond correctly
- [x] **Story 1.3: Workflow Engine proof of concept** ✅
  - Go-based standalone workflow engine CLI tool
  - Parses YAML workflow files with structure validation
  - Executes single-step workflows (PoC scope per architecture)
  - **Command**: `./dist/workflow-engine <workflow-file.yaml>`
  - **Sample**: Created 3-step documentation workflow with multi-agent collaboration
  - **Integration**: `workflow-runner` agent can invoke the engine from OpenCode

## 🏆 **EPIC 1 COMPLETE**

**All foundational proof of concept components are working:**
- ✅ **Story 1.1**: Plugin development environment with Hello World agent
- ✅ **Story 1.2**: BMAD Agent Loader converts `.md` files to `opencode.json`
- ✅ **Story 1.3**: Workflow Engine executes single-step YAML workflows

**Epic 1 has successfully validated the core BMAD-OpenCode integration architecture. Ready for Epic 2 full engine implementation!**

## Contributing

This is the foundational proof of concept for the larger BMAD OpenCode Engine project. See our [architecture documentation](docs/architecture/epic-1-architecture.md) for technical details.

## License

MIT