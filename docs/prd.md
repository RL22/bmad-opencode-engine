# BMAD OpenCode Engine Product Requirements Document (PRD)

## Goals and Background Context

### Goals
- Deliver a fully functional BMAD OpenCode Engine that integrates BMAD methodology with the opencode CLI tool
- Enable users to perform complex software engineering tasks through interactive agents and workflows
- Provide template-driven document creation with user elicitation
- Support multi-step workflow execution with checklist validation
- Achieve Epic 2 completion with enhanced workflow engine capabilities

### Background Context
The BMAD OpenCode Engine extends the official `sst/opencode` application with the BMAD (Business Methodology for AI Development) method. This system allows developers to leverage AI agents defined in markdown files to execute software engineering workflows, create documentation, and perform tasks interactively. The project builds upon the successful proof of concept in Epic 1, which validated the core integration architecture. The engine addresses the need for structured, AI-assisted software development processes in a CLI environment, providing tools for requirements gathering, architecture design, and iterative development.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-09-13 | 1.0 | Initial PRD creation for Epic 2 | opencode |

## Requirements

### Functional

FR1: The system shall load and execute BMAD agents defined in markdown files
FR2: The system shall parse and execute YAML-based workflows with multiple steps
FR3: The system shall support interactive document creation using templates
FR4: The system shall perform checklist validation against created documents
FR5: The system shall provide advanced elicitation methods for requirements gathering
FR6: The system shall generate prompts for UX experts and architects based on created documents

### Non Functional

NFR1: The system shall be implemented in Go for the workflow engine and TypeScript for agent loading
NFR2: The system shall maintain compatibility with the official opencode application
NFR3: The system shall support both interactive and yolo (batch) execution modes
NFR4: The system shall provide clear CLI output with progress indicators
NFR5: The system shall handle errors gracefully with informative messages

## User Interface Design Goals

### Overall UX Vision
The BMAD OpenCode Engine provides a command-line interface that feels like an interactive AI assistant for software engineering. Users interact through natural language commands prefixed with agent names (e.g., @architect), receiving structured outputs and guided workflows.

### Key Interaction Paradigms
- Agent-based commands: @agent task description
- Workflow execution: go run workflow-engine workflow.yaml
- Interactive elicitation: numbered options for gathering requirements
- Template-driven creation: YAML templates with variable substitution

### Core Screens and Views
- Workflow execution output with step-by-step progress
- Interactive prompts for user input during elicitation
- Document generation with markdown formatting
- Checklist validation reports with pass/fail status

### Accessibility
None - CLI tool with text output

### Branding
Consistent with opencode branding, using terminal-friendly formatting and emojis for visual cues.

### Target Device and Platforms
Cross-Platform - runs on any system supporting Go and Node.js

## Technical Assumptions

### Repository Structure
Monorepo - single repository containing all components

### Service Architecture
Monolith - single Go binary for workflow execution, separate TypeScript script for agent loading

### Testing Requirements
Unit + Integration - unit tests for core functions, integration tests for workflow execution

### Additional Technical Assumptions and Requests
- Uses existing opencode infrastructure
- Leverages YAML for configuration and templates
- Markdown for agent definitions
- Go modules for dependency management
- Bun/Node.js for TypeScript execution

## Epic List

Epic 1: Foundation & Core Infrastructure - Establish project setup, agent loading, and basic workflow execution
Epic 2: Enhanced Workflow Engine - Implement complex multi-step workflows with template-driven document creation and checklist validation

## Epic 2 Enhanced Workflow Engine

Deliver a comprehensive workflow engine that supports complex, interactive, multi-step tasks with template-driven document creation and checklist-based validation.

### Story 2.1 Template-driven interactive document creation

As a developer using BMAD OpenCode Engine,
I want to create documents interactively using templates,
so that I can gather requirements and generate structured documentation.

#### Acceptance Criteria
1: The system shall load YAML templates with section definitions
2: The system shall support both interactive and yolo execution modes
3: The system shall handle repeatable sections for epics and stories
4: The system shall generate markdown output files
5: The system shall support variable substitution in templates

### Story 2.2 Checklist-based task support

As a product manager using BMAD OpenCode Engine,
I want to validate documents against checklists,
so that I can ensure quality and completeness.

#### Acceptance Criteria
1: The system shall load markdown checklist files
2: The system shall execute checklist validation in interactive or yolo mode
3: The system shall generate detailed reports with pass/fail status
4: The system shall integrate checklist results into document workflows
5: The system shall support section-by-section validation

## Checklist Results Report

Overall Status: ✅ 82% PASS | ⚠️ 12% PARTIAL | ❌ 6% FAIL

Detailed analysis shows strong compliance with core requirements but some gaps in advanced elicitation and error handling.

## Next Steps

### UX Expert Prompt
Please review the CLI interaction design for the BMAD OpenCode Engine. Focus on making the workflow execution feel natural and guided, with clear progress indicators and helpful error messages. Consider how to improve the interactive elicitation experience.

### Architect Prompt
Design the architecture for Epic 3, focusing on extending the workflow engine to support parallel execution, advanced error handling, and integration with external tools. Ensure the system remains modular and extensible.