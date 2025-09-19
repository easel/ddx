# Project Principles

*These principles extend the HELIX workflow principles for this specific project.*

## Core Development Articles

### Article I: Library-First Principle
Every feature MUST begin as a standalone library with clear boundaries and minimal dependencies.

### Article II: CLI Interface Mandate  
All libraries MUST expose functionality through CLI:
- Accept text input (stdin/args/files)
- Produce text output (stdout)
- Support JSON for structured data exchange

### Article III: Test-First Imperative
NO implementation code shall be written before:
1. Tests are written and documented
2. Tests are reviewed and approved
3. Tests are confirmed to FAIL (Red phase of TDD)

### Article VII: Simplicity Gate
- Maximum 3 major components for initial implementation
- Additional complexity requires documented justification
- No premature optimization or future-proofing

### Article VIII: Anti-Abstraction Principle
- Use framework features directly
- No unnecessary wrapper layers
- Single source of truth for each concept

### Article IX: Integration-First Testing
- Prefer real environments over mocks
- Contract tests mandatory before implementation
- Use actual databases and services in tests

## Project-Specific Principles
<!-- Define additional principles specific to this project -->
[PROJECT_SPECIFIC_PRINCIPLE_1]
[PROJECT_SPECIFIC_PRINCIPLE_2]

## Technology Constraints
<!-- Document any technology decisions that act as principles -->
- Primary Language: [LANGUAGE]
- Framework: [FRAMEWORK]
- Database: [DATABASE]

## Exceptions Log
<!-- Document any necessary violations with justification and timeline for resolution -->
| Date | Principle | Exception | Justification | Resolution Timeline |
|------|-----------|-----------|---------------|-------------------|
| | | | | |