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

#### Error Handling Specifications
- **Template Errors**: Invalid YAML syntax, missing required fields, circular references
- **File System Errors**: Permission denied, file not found, disk space issues
- **Network Errors**: API timeouts, connection failures, rate limiting
- **Agent Errors**: Agent loading failures, model unavailability, token limits
- **Workflow Errors**: Step failures, state corruption, execution timeouts
- **User Input Errors**: Invalid formats, missing required data, cancellation requests

#### Error Recovery Mechanisms
- **Graceful Degradation**: Continue processing when non-critical components fail
- **Automatic Retry**: Implement exponential backoff for transient failures
- **State Preservation**: Save progress before failures to enable resumption
- **User Guidance**: Provide clear next steps and recovery options
- **Logging**: Comprehensive error logging with context and debugging information

#### User Experience
- **Clear Error Messages**: Human-readable explanations of what went wrong
- **Actionable Guidance**: Specific steps users can take to resolve issues
- **Progress Preservation**: Don't lose work when recoverable errors occur
- **Help Integration**: Link to relevant documentation or support resources

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
Unit + Integration + E2E - comprehensive testing strategy covering all system components

#### Unit Testing
- Core Go functions (workflow engine, YAML parsing, template processing)
- TypeScript utilities (agent loading, configuration parsing)
- Template validation and processing logic
- Checklist execution and validation logic

#### Integration Testing
- Workflow execution end-to-end scenarios
- Agent loading and configuration integration
- Template processing with real YAML files
- Checklist validation against sample documents
- Cross-component data flow validation

#### End-to-End Testing
- Complete workflow execution from start to finish
- Interactive mode user flows with simulated input
- YOLO mode batch processing validation
- Error scenario handling and recovery
- Multi-step workflow state management

#### Test Coverage Targets
- Unit tests: 80%+ code coverage
- Integration tests: All critical user journeys
- E2E tests: All primary workflows and error paths

#### QA Gates
- Code review requirements before merge
- Automated test execution on all PRs
- Manual QA validation for UI/UX changes (when applicable)
- Performance regression testing
- Security vulnerability scanning

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

#### Story Details
**Story Points**: 8
**Priority**: High
**Dependencies**: Epic 1 completion, template system foundation
**Estimated Effort**: 2-3 days

#### Acceptance Criteria

**Functional Requirements:**
1. **Template Loading**: The system shall successfully load and parse YAML template files with complex section hierarchies
2. **Section Processing**: The system shall process all template sections including nested sections and repeatable elements
3. **Variable Substitution**: The system shall support dynamic variable substitution throughout template content
4. **Interactive Mode**: The system shall provide step-by-step user interaction for template completion with clear prompts
5. **YOLO Mode**: The system shall support batch processing of all template sections without user interaction
6. **Markdown Generation**: The system shall generate properly formatted markdown output files with correct structure
7. **Error Handling**: The system shall gracefully handle template parsing errors and missing required fields
8. **Progress Indication**: The system shall show clear progress indicators during template processing

**User Experience:**
9. **Clear Prompts**: Each template section shall have descriptive prompts explaining what information is needed
10. **Input Validation**: The system shall validate user input against expected formats and provide helpful error messages
11. **Context Preservation**: The system shall maintain context between template sections for coherent document generation
12. **Help System**: Users shall be able to request help or examples for any template section

**Technical Specifications:**
13. **Template Format**: Templates shall use YAML format with defined schema for sections, types, and instructions
14. **Section Types**: The system shall support multiple section types (paragraphs, bullet-list, numbered-list, table, etc.)
15. **Repeatable Sections**: The system shall handle repeatable sections for epics, stories, and other iterative content
16. **File Output**: Generated documents shall be saved with proper naming conventions and file structure
17. **Template Validation**: The system shall validate template structure before processing begins

**Quality Assurance:**
18. **Template Testing**: All templates shall be tested with both interactive and YOLO modes
19. **Output Validation**: Generated documents shall be validated for proper markdown formatting
20. **Error Scenarios**: The system shall handle edge cases like missing templates, invalid YAML, and user cancellations

### Story 2.2 Checklist-based task support

As a product manager using BMAD OpenCode Engine,
I want to validate documents against checklists,
so that I can ensure quality and completeness.

#### Story Details
**Story Points**: 8
**Priority**: High
**Dependencies**: Epic 1 completion, checklist framework foundation
**Estimated Effort**: 2-3 days

#### Acceptance Criteria

**Functional Requirements:**
1. **Checklist Loading**: The system shall successfully load and parse markdown checklist files with proper structure
2. **Validation Execution**: The system shall execute checklist validation against target documents or project artifacts
3. **Interactive Mode**: The system shall support section-by-section validation with user confirmation at each step
4. **YOLO Mode**: The system shall support complete batch validation of entire checklists without user interaction
5. **Report Generation**: The system shall generate comprehensive reports with PASS/PARTIAL/FAIL status for each item
6. **Workflow Integration**: The system shall integrate checklist results into document creation and validation workflows
7. **Status Tracking**: The system shall track validation progress and provide real-time status updates

**User Experience:**
8. **Clear Feedback**: Each checklist item shall provide clear pass/fail criteria and validation rationale
9. **Progress Tracking**: Users shall see real-time progress through checklist validation with percentage completion
10. **Detailed Reports**: The system shall generate detailed reports highlighting issues, recommendations, and next steps
11. **Actionable Results**: Validation results shall include specific recommendations for addressing failed items

**Technical Specifications:**
12. **Checklist Format**: Checklists shall use standardized markdown format with clear item structure and validation criteria
13. **Validation Logic**: The system shall support multiple validation types (presence, format, content, dependency checks)
14. **Result Storage**: Validation results shall be stored and accessible for future reference and workflow integration
15. **Error Handling**: The system shall gracefully handle checklist parsing errors and missing validation targets
16. **Performance**: The system shall validate checklists efficiently without significant performance impact

**Quality Assurance:**
17. **Checklist Testing**: All checklists shall be tested with various document types and validation scenarios
18. **Report Accuracy**: Generated reports shall accurately reflect validation results with no false positives/negatives
19. **Edge Cases**: The system shall handle edge cases like empty checklists, missing documents, and partial validations
20. **Integration Testing**: Checklist validation shall be tested within complete workflow execution scenarios

**Business Value:**
21. **Quality Assurance**: Checklists shall ensure consistent quality standards across all project deliverables
22. **Risk Mitigation**: Validation shall identify potential issues before they impact development progress
23. **Process Compliance**: The system shall enforce adherence to established development and documentation standards
24. **Continuous Improvement**: Validation results shall provide data for improving templates, checklists, and processes

## Checklist Results Report

Overall Status: ✅ 91% PASS | ⚠️ 7% PARTIAL | ❌ 2% FAIL

**Updated Analysis (Post-Enhancement):**
- ✅ **Story Details**: Comprehensive acceptance criteria added for both Epic 2 stories
- ✅ **Testing Strategy**: Detailed testing requirements and QA gates defined
- ✅ **Error Handling**: Specific error scenarios and recovery mechanisms specified
- ⚠️ **Advanced Elicitation**: Basic framework in place, could benefit from more sophisticated methods
- ✅ **Technical Specifications**: Clear implementation guidance provided
- ✅ **Quality Assurance**: Testing and validation processes well-defined

**Critical Gaps Resolved:**
1. ✅ Detailed acceptance criteria for Epic 2 stories
2. ✅ Comprehensive testing strategy
3. ✅ Specific error handling requirements
4. ✅ Story points and effort estimates
5. ✅ Dependencies and priority definitions

## Next Steps

### UX Expert Prompt
Please review the CLI interaction design for the BMAD OpenCode Engine. Focus on making the workflow execution feel natural and guided, with clear progress indicators and helpful error messages. Consider how to improve the interactive elicitation experience.

### Architect Prompt
Design the architecture for Epic 3, focusing on extending the workflow engine to support parallel execution, advanced error handling, and integration with external tools. Ensure the system remains modular and extensible.