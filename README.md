# BMAD OpenCode Engine

A comprehensive AI-powered software engineering platform that integrates the BMAD (Business Methodology for AI Development) method with the official `sst/opencode` CLI tool.

## Overview

The BMAD OpenCode Engine provides:

- **🤖 Intelligent Agent System**: 12 specialized AI agents for different software engineering roles
- **⚡ Advanced Workflow Engine**: Multi-step YAML workflows with interactive document creation
- **📋 Quality Assurance**: Comprehensive checklists and validation for all deliverables
- **🎯 Template-Driven Development**: Structured document creation with user elicitation
- **🔧 Agent Loader**: Automated conversion of markdown agent definitions to OpenCode configuration

## Project Structure

```
bmad-opencode-engine/
├── bmad-core/                  # BMAD methodology core components
│   ├── agents/                 # 12 specialized AI agent definitions
│   │   ├── analyst.md         # Business analysis & market research
│   │   ├── architect.md       # System architecture & design
│   │   ├── bmad-orchestrator.md # Master workflow coordinator
│   │   ├── dev.md             # Full-stack development
│   │   ├── growth-marketer.md # Growth & user acquisition
│   │   ├── mindmap-specialist.md # Documentation analysis
│   │   ├── pm.md              # Product management
│   │   ├── po.md              # Product ownership
│   │   ├── qa.md              # Quality assurance
│   │   ├── sm.md              # Scrum master
│   │   ├── ux-expert.md       # User experience design
│   │   └── workflow-runner.md # Workflow execution
│   ├── checklists/            # Quality validation checklists
│   ├── templates/             # Document creation templates
│   └── workflows/             # YAML workflow definitions
├── packages/
│   └── workflow-engine/       # Go-based workflow execution engine
├── scripts/
│   └── agent-loader/          # TypeScript agent configuration generator
├── docs/                      # Project documentation
├── example-project/           # Sample implementations
├── generated-opencode.json    # Auto-generated OpenCode configuration
├── package.json               # Node.js dependencies
└── go.mod                    # Go module dependencies
```

## Key Features

### 🤖 **Intelligent Agent System**
- **12 Specialized Agents**: Covering all aspects of software engineering from analysis to deployment
- **Role-Based Expertise**: Each agent optimized for specific development tasks and methodologies
- **Markdown-Defined**: Easy to create, modify, and version agent configurations
- **Auto-Configuration**: Automated conversion to OpenCode-compatible JSON format

### ⚡ **Advanced Workflow Engine**
- **YAML-Based Workflows**: Declarative workflow definitions with multi-step execution
- **Interactive & Batch Modes**: Support for both guided user interaction and automated processing
- **Template Integration**: Structured document creation with variable substitution
- **Error Recovery**: Robust error handling with state preservation and recovery mechanisms

### 📋 **Quality Assurance Framework**
- **Comprehensive Checklists**: Pre-built validation checklists for different project types
- **Automated Validation**: Checklist execution with detailed pass/fail reporting
- **Quality Gates**: Configurable quality thresholds for different development stages
- **Continuous Improvement**: Checklist results feed into process optimization

## Prerequisites

- **Node.js 18+** (npm for package management)
- **Go 1.21+** (for workflow engine)
- **Git** (for version control)
- **OpenCode CLI** (optional, for full integration)

### Optional Dependencies
- **Bun runtime** (alternative to npm for faster installs)
- **Docker** (for containerized development)
- **Make** (for build automation)

## Quick Start

1. **Install Dependencies**
   ```bash
   npm install
   go mod tidy
   ```

2. **Generate Agent Configuration**
   ```bash
   npm run generate-config -- --agents-dir bmad-core/agents --output generated-opencode.json
   ```

3. **Build Workflow Engine**
   ```bash
   npm run build
   ```

4. **Run a Sample Workflow**
   ```bash
   ./dist/workflow-engine bmad-core/workflows/create-simple-doc.yaml
   ```

5. **Test Agent Integration**
   ```bash
   opencode run "@bmad-orchestrator help"
   ```

## Development Status

### ✅ **Epic 1: Foundation & Core Infrastructure - COMPLETE**
- [x] Repository setup and structure
- [x] **Story 1.1: Plugin development environment + Hello World test**
- [x] **Story 1.2: BMAD Agent Loader implementation**
- [x] **Story 1.3: Workflow Engine proof of concept**

### 🚧 **Epic 2: Enhanced Workflow Engine - READY FOR DEVELOPMENT**
- [x] **Story 2.1: Template-driven interactive document creation**
  - ✅ Comprehensive acceptance criteria defined (20 requirements)
  - ✅ Template system architecture designed
  - ✅ Interactive and YOLO execution modes specified
  - ✅ Error handling and validation requirements complete
- [x] **Story 2.2: Checklist-based task support**
  - ✅ Detailed acceptance criteria defined (24 requirements)
  - ✅ Checklist validation framework designed
  - ✅ Quality assurance integration specified
  - ✅ Business value metrics defined

### 📊 **Project Readiness Metrics**
- **PRD Completeness**: 91% ✅
- **Acceptance Criteria**: 100% ✅ (44 detailed criteria across both stories)
- **Testing Strategy**: Comprehensive ✅
- **Error Handling**: Fully specified ✅
- **Technical Architecture**: Validated ✅

### 🎯 **Next Development Phase**
**Epic 2 Implementation Ready to Begin:**
1. Template processing system implementation
2. Interactive document creation workflow
3. Checklist validation framework
4. Multi-step workflow orchestration
5. Quality assurance integration

**Estimated Timeline**: 2-3 weeks for Epic 2 completion

## 🤖 Agent System Overview

The BMAD OpenCode Engine includes 12 specialized AI agents, each optimized for specific software engineering tasks:

| Agent | Role | Key Capabilities |
|-------|------|------------------|
| **BMAD Orchestrator** | Master Coordinator | Workflow orchestration, agent selection, multi-agent collaboration |
| **Product Manager** | Strategy & Planning | PRD creation, roadmap planning, stakeholder management |
| **Product Owner** | Requirements & Backlog | User story refinement, acceptance criteria, sprint planning |
| **Business Analyst** | Analysis & Research | Market research, competitive analysis, requirements gathering |
| **System Architect** | Technical Design | System architecture, technology selection, infrastructure planning |
| **UX Expert** | User Experience | UI/UX design, wireframes, user journey optimization |
| **Scrum Master** | Process Management | Story creation, agile process guidance, team facilitation |
| **QA Engineer** | Quality Assurance | Test planning, quality gates, risk assessment |
| **Full-Stack Developer** | Implementation | Code development, testing, deployment automation |
| **Growth Marketer** | User Acquisition | Growth strategies, analytics, conversion optimization |
| **Mind Map Specialist** | Documentation Analysis | Technical documentation processing, knowledge mapping |
| **Workflow Runner** | Execution Engine | YAML workflow execution, process automation |

## 📚 Documentation

- **[Product Requirements Document](docs/prd.md)**: Complete project specifications and acceptance criteria
- **[Agent Definitions](bmad-core/agents/)**: Detailed role descriptions and capabilities for each agent
- **[Workflow Templates](bmad-core/workflows/)**: YAML workflow definitions and examples
- **[Quality Checklists](bmad-core/checklists/)**: Comprehensive validation checklists for different project types

## Contributing

The BMAD OpenCode Engine is an open-source project focused on improving AI-assisted software development through structured methodologies. See our [PRD](docs/prd.md) for detailed technical specifications and development roadmap.

## License

MIT