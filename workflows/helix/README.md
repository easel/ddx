# HELIX Workflow

An AI-assisted development workflow designed for iterative software development with human-AI collaboration.

## Overview

HELIX is a structured approach to building software where human creativity and AI capabilities work together throughout the development cycle. Each phase has clear responsibilities for both human developers and AI agents.

## Phases

1. **Frame** - Define the problem and establish context
2. **Design** - Architect the solution approach
3. **Implement** - Build the system with AI assistance
4. **Test** - Validate functionality and quality
5. **Deploy** - Release to production with monitoring
6. **Iterate** - Learn and improve for the next cycle

## Input Gates

Each phase (except Frame) has input gates that validate the previous phase's outputs before allowing progression:

- **Design** cannot start until Frame outputs are validated
- **Implement** cannot start until Design is reviewed and approved
- **Test** cannot start until Implementation is complete and reviewed
- **Deploy** cannot start until all Tests pass
- **Iterate** begins once the system is deployed and operational

This ensures quality at each step and prevents skipping crucial validation.

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

## Getting Started

```bash
ddx workflow apply helix
```

This will initialize a new project using the HELIX workflow, guiding you through each phase with appropriate templates and AI assistance.

## The Helix Concept

*Note: The workflow name comes from the double helix structure of DNA, which serves as an interesting metaphor for human-AI collaboration - two complementary strands that support each other, with connection points (quality gates) ensuring structural integrity. Like DNA, each iteration builds on the previous one, creating an ascending spiral of progress.*