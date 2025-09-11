# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DDx (Document-Driven Development eXperience) is a CLI toolkit for AI-assisted development that helps developers share templates, prompts, and patterns across projects. The project follows a medical differential diagnosis metaphor - using structured documentation to diagnose project issues, prescribe solutions, and share improvements.

## Architecture

The project has a dual structure:
- **CLI Application** (`/cli/`): Go-based command-line tool built with Cobra framework
- **Content Repository** (root): Templates, patterns, prompts, and configurations for the DDx toolkit

### Key Components

- `cli/` - Go CLI application source code
  - `cmd/` - Cobra command implementations (init, list, apply, diagnose, update, contribute)
  - `internal/` - Internal packages (config, templates, git utilities)
  - `main.go` - Application entry point
- `templates/` - Project templates (NextJS, Python, etc.)
- `patterns/` - Reusable code patterns and examples
- `prompts/` - AI prompts and instructions (Claude-specific and general)
- `scripts/` - Automation scripts and git hooks
- `configs/` - Tool configurations (ESLint, Prettier, TypeScript)

## Development Commands

### CLI Development (run from `/cli/` directory)

```bash
# Build and test
make build          # Build for current platform
make test           # Run Go tests
make lint           # Run golangci-lint (or go vet if not available)
make fmt            # Format Go code

# Development workflow
make all            # Clean, deps, test, build
make dev            # Development mode with file watching (requires air)
make run ARGS="..."  # Run CLI with arguments
make install        # Install locally to ~/.local/bin/ddx

# Dependencies
make deps           # Install and tidy Go modules
make update-deps    # Update all dependencies

# Multi-platform builds
make build-all      # Build for all platforms
make release        # Create release archives
```

### Project Structure Navigation

The CLI uses git subtree for managing the relationship between individual projects and the master DDx repository. The `.ddx.yml` configuration file defines:
- Repository URL and branch
- Included resources (prompts, scripts, templates, patterns)
- Template variables and overrides
- Git subtree settings

### Key Patterns

1. **Command Structure**: Each CLI command is implemented as a separate file in `cli/cmd/`
2. **Configuration Management**: Uses Viper for config file handling with YAML format
3. **Template Processing**: Variable substitution system for customizing templates
4. **Git Integration**: Built on git subtree for reliable version control and contribution workflows
5. **Cross-Platform Support**: Makefile supports building for multiple platforms (macOS, Linux, Windows)

### Testing and Quality

- Go tests are in `*_test.go` files alongside source code
- Linting uses golangci-lint (fallback to go vet)
- Code formatting with `go fmt`
- Cross-platform compatibility is maintained

### CLI Command Overview

- `ddx init` - Initialize DDx in a project (with optional template)
- `ddx list` - Show available templates, patterns, and prompts
- `ddx apply <resource>` - Apply templates, patterns, or configurations
- `ddx diagnose` - Analyze project health and suggest improvements
- `ddx update` - Update toolkit from master repository
- `ddx contribute` - Share improvements back to community

The CLI follows the medical metaphor throughout, treating projects as patients that need diagnosis and treatment through appropriate templates and patterns.
