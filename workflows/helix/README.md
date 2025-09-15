# HELIX Workflow

A test-driven development workflow with AI-assisted collaboration for building high-quality software iteratively.

## Overview

HELIX enforces Test-Driven Development (TDD) through a structured phase approach where tests are written BEFORE implementation. This ensures quality is built-in from the start and specifications are executable. Human creativity and AI capabilities collaborate throughout, with tests serving as the contract between design and implementation.

## Phases

1. **Frame** - Define the problem and establish context
2. **Design** - Architect the solution approach
3. **Test** - Write failing tests that define system behavior (Red phase)
4. **Build** - Implement code to make tests pass (Green phase)
5. **Deploy** - Release to production with monitoring
6. **Iterate** - Learn and improve for the next cycle

## Input Gates

Each phase (except Frame) has input gates that validate the previous phase's outputs before allowing progression:

- **Design** cannot start until Frame outputs are validated
- **Test** cannot start until Design is reviewed and approved
- **Build** cannot start until Tests are written and failing (Red phase)
- **Deploy** cannot start until all Tests pass (Green phase)
- **Iterate** begins once the system is deployed and operational

This test-first approach ensures specifications drive implementation and quality is built in from the start.

## Human-AI Collaboration

Throughout the workflow, responsibilities are shared:

### Human Responsibilities
- Problem definition and creative vision
- Strategic decisions and architecture choices
- Code review and quality assessment
- User experience and business logic

### AI Agent Responsibilities
- Pattern recognition and suggestions
- Code generation and refactoring
- Test case generation
- Documentation and analysis

## Security Integration

HELIX integrates security practices throughout every phase, following DevSecOps principles to ensure security is built-in rather than bolted-on:

### Security-First Approach
- **Frame**: Security requirements, threat modeling, and compliance analysis established upfront
- **Design**: Security architecture and controls designed into system structure
- **Test**: Security test suites created alongside functional tests
- **Build**: Secure coding practices and automated security scanning integrated
- **Deploy**: Security monitoring and incident response procedures activated
- **Iterate**: Security metrics tracked and security posture continuously improved

### Key Security Artifacts
- **Security Requirements**: Comprehensive security and compliance requirements
- **Threat Model**: STRIDE-based threat analysis with risk assessment
- **Security Architecture**: Defense-in-depth design with security controls
- **Security Tests**: Automated and manual security testing procedures
- **Security Monitoring**: Production security monitoring and alerting

### Security Quality Gates
Each phase includes security checkpoints that must be satisfied before progression, ensuring security requirements are met throughout the development lifecycle.

## Getting Started

```bash
ddx workflow apply helix
```

This will initialize a new project using the HELIX workflow, guiding you through each phase with appropriate templates and AI assistance.

## Why Test-First?

The HELIX workflow enforces writing tests before implementation because:

1. **Tests ARE the Specification** - Tests define exactly what the system should do
2. **Clear Definition of Done** - Implementation is complete when tests pass
3. **Prevents Over-Engineering** - Only write code needed to pass tests
4. **Built-in Quality** - Bugs are caught immediately, not later
5. **Refactoring Safety** - Green tests provide confidence to improve code

## The TDD Cycle

Within the Test and Build phases, we follow the Red-Green-Refactor cycle:

1. **Red** (Test Phase) - Write a failing test that defines desired behavior
2. **Green** (Build Phase) - Write minimal code to make the test pass
3. **Refactor** (Build Phase) - Improve code quality while keeping tests green

## The Helix Concept

*The workflow name comes from the double helix structure of DNA, which serves as an interesting metaphor for human-AI collaboration - two complementary strands that support each other, with connection points (quality gates) ensuring structural integrity. Like DNA, each iteration builds on the previous one, creating an ascending spiral of progress.*