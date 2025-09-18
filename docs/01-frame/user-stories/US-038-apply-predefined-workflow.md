# User Story: [US-038] - Developer Applying Standard Workflow

**Story ID**: US-038
**Priority**: P0
**Feature**: FEAT-005 - Workflow Execution Engine
**Created**: 2025-01-14
**Status**: Defined

## Story Description

**As a** developer working on a new project
**I want** to apply a predefined workflow to my project
**So that** I can follow proven patterns without manual setup and reduce time from project initiation to productive development

## Business Value

- Reduces project setup time from hours to seconds
- Ensures adherence to proven development methodologies
- Eliminates need to recreate workflows from scratch
- Provides immediate access to team's best practices

## Acceptance Criteria

### AC-001: Workflow Discovery
- **Given** I am in a project directory
- **When** I run workflow list command
- **Then** I can see available workflows with descriptions and categories

### AC-002: Workflow Application
- **Given** I have selected a workflow
- **When** I apply it to my current project
- **Then** the workflow is initialized with all required phases and templates

### AC-003: Variable Substitution
- **Given** a workflow requires project-specific variables
- **When** I apply the workflow
- **Then** I am prompted for required variables and they are substituted correctly

### AC-004: Progress Tracking
- **Given** I have applied a workflow
- **When** I check the workflow status
- **Then** I can see current phase and overall progress percentage

### AC-005: Artifact Generation
- **Given** a workflow phase generates artifacts
- **When** the phase executes
- **Then** expected files and documents are created in the correct locations

### AC-006: Session Recovery
- **Given** my workflow execution was interrupted
- **When** I restart the workflow
- **Then** I can resume from where I left off without losing progress

## Definition of Done

- [ ] User can discover available workflows
- [ ] User can apply workflow to current project
- [ ] Required variables are collected during application
- [ ] Progress is visible and trackable
- [ ] Generated artifacts match expectations
- [ ] Workflow can be resumed after interruption
- [ ] Application completes in under 10 seconds
- [ ] All acceptance criteria have automated tests

## Notes

This is a foundational story that enables all workflow-based functionality. Success here directly impacts user adoption and satisfaction.