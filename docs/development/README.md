# Development Documentation

> **Last Updated**: 2025-09-11
> **Status**: Active
> **Owner**: DDx Team

## Overview

Developer documentation for contributing to and extending the DDx toolkit.

## Contents

### [Contributing](/docs/development/contributing/)
Guidelines for contributing code, documentation, and resources.

### [Standards](/docs/development/standards/)
Coding standards, conventions, and best practices.

### [Testing](/docs/development/testing/)
Testing strategies, frameworks, and guidelines.

### [CI/CD](/docs/development/ci-cd/)
Continuous integration and deployment documentation.

### [Tools](/docs/development/tools/)
Development tools, setup instructions, and configurations.

### [Release](/docs/development/release/)
Release process, versioning, and changelog management.

## Quick Start for Developers

1. [[tools/setup]] - Set up your development environment
2. [[contributing/guidelines]] - Read contribution guidelines
3. [[standards/coding-standards]] - Understand coding standards
4. [[testing/strategy]] - Learn about testing approach

## Development Workflow

```bash
# Clone the repository
git clone https://github.com/yourusername/ddx.git
cd ddx

# Set up development environment
cd cli
make deps

# Run tests
make test

# Build the CLI
make build

# Run linting
make lint

# Install locally
make install
```

## Key Commands

```bash
# Development mode with hot reload
make dev

# Run all checks (lint, test, build)
make all

# Build for all platforms
make build-all

# Create release archives
make release
```

## Project Structure

```
ddx/
├── cli/                 # Go CLI application
│   ├── cmd/            # Command implementations
│   ├── internal/       # Internal packages
│   └── main.go        # Entry point
├── templates/          # Project templates
├── patterns/           # Code patterns
├── prompts/           # AI prompts
├── configs/           # Tool configurations
└── docs/              # Documentation
```

## Technologies

- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **Testing**: Go testing package
- **Linting**: golangci-lint
- **Git Hooks**: Lefthook

## Related Documentation

- [[architecture/cli-architecture]] - CLI architecture details
- [[implementation/setup/installation]] - Installation guide
- [[usage/getting-started/quick-start]] - User quick start