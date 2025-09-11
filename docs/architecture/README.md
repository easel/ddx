# Architecture Documentation

> **Last Updated**: 2025-09-11
> **Status**: Active
> **Owner**: DDx Team

## Overview

Technical architecture documentation for the DDx toolkit, including system design, component specifications, and architectural decisions.

## Contents

### [Decisions](/docs/architecture/decisions/)
Architecture Decision Records (ADRs) documenting key technical choices.

### [Diagrams](/docs/architecture/diagrams/)
System architecture diagrams and visual representations.

### [Components](/docs/architecture/components/)
Individual component specifications and interfaces.

### [Integrations](/docs/architecture/integrations/)
Integration patterns, APIs, and external system connections.

### [Data](/docs/architecture/data/)
Data models, schemas, and data flow documentation.

### [Security](/docs/architecture/security/)
Security architecture, threat models, and security controls.

## Key Documents

- [[architecture-overview]] - High-level system architecture
- [[cli-architecture]] - CLI application structure
- [[git-subtree-integration]] - Git subtree architecture

## Architecture Principles

1. **Modularity** - Loosely coupled, highly cohesive components
2. **Extensibility** - Plugin-based architecture for custom extensions
3. **Portability** - Cross-platform compatibility (macOS, Linux, Windows)
4. **Simplicity** - Minimal dependencies, clear interfaces
5. **Performance** - Efficient resource usage, fast execution

## Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **Configuration**: Viper (YAML)
- **Version Control**: Git (subtree-based)
- **Build System**: Make
- **Testing**: Go testing package

## Related Documentation

- [[implementation/setup/installation]] - Installation guide
- [[development/standards/coding-standards]] - Coding standards
- [[usage/getting-started/quick-start]] - Getting started guide