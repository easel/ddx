# Development Environment Setup

> **Last Updated**: 2025-09-11
> **Status**: Active
> **Owner**: DDx Team

## Overview

Complete guide for setting up a development environment for contributing to DDx.

## Context

This guide is for developers who want to contribute to the DDx project. For end-user installation, see [[usage/getting-started/quick-start]].

## Prerequisites

### Required Software

| Software | Version | Purpose |
|----------|---------|---------|
| Go | 1.21+ | Primary development language |
| Git | 2.30+ | Version control |
| Make | 3.81+ | Build automation |

### Optional Software

| Software | Version | Purpose |
|----------|---------|---------|
| golangci-lint | Latest | Code linting |
| air | Latest | Hot reload for development |
| gh | Latest | GitHub CLI for PRs |
| lefthook | Latest | Git hooks management |

## Setup Steps

### 1. Clone the Repository

```bash
# Clone via HTTPS
git clone https://github.com/yourusername/ddx.git

# Or via SSH
git clone git@github.com:yourusername/ddx.git

cd ddx
```

### 2. Install Go

#### macOS
```bash
brew install go
```

#### Linux
```bash
# Download from https://go.dev/dl/
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

#### Windows
Download installer from [https://go.dev/dl/](https://go.dev/dl/)

### 3. Install Development Dependencies

```bash
cd cli

# Install Go dependencies
make deps

# Install optional tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest
go install github.com/evilmartians/lefthook@latest
```

### 4. Configure Git Hooks

```bash
# From repository root
lefthook install
```

This sets up pre-commit hooks for:
- Code formatting
- Linting
- Testing
- Security checks

### 5. Verify Installation

```bash
# Run all checks
make all

# Expected output:
# ✓ Dependencies installed
# ✓ Tests passing
# ✓ Linting clean
# ✓ Build successful
```

## Development Workflow

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./cmd/... -run TestInit
```

### Linting

```bash
# Run linting
make lint

# Auto-fix issues
golangci-lint run --fix
```

### Hot Reload Development

```bash
# Start development server with hot reload
make dev

# Or manually with air
air
```

### Code Formatting

```bash
# Format all Go code
make fmt

# Or manually
go fmt ./...
```

## IDE Configuration

### VS Code

Install recommended extensions:
```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.makefile-tools",
    "eamodio.gitlens",
    "yzhang.markdown-all-in-one"
  ]
}
```

Settings (`.vscode/settings.json`):
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.formatTool": "goimports",
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand/IntelliJ

1. Open project root
2. Configure Go SDK: File → Project Structure → SDKs
3. Enable Go modules: Preferences → Go → Go Modules
4. Configure make: Preferences → Build → Make

## Project Structure

```
ddx/
├── cli/                    # Go CLI application
│   ├── cmd/               # Command implementations
│   │   ├── root.go       # Root command
│   │   ├── init.go       # Init command
│   │   ├── list.go       # List command
│   │   └── ...
│   ├── internal/          # Internal packages
│   │   ├── config/       # Configuration management
│   │   ├── template/     # Template processing
│   │   └── git/          # Git operations
│   ├── main.go           # Entry point
│   ├── go.mod            # Go module definition
│   └── Makefile          # Build automation
├── templates/             # Project templates
├── patterns/              # Code patterns
├── prompts/              # AI prompts
├── configs/              # Tool configurations
└── docs/                 # Documentation
```

## Common Tasks

### Adding a New Command

1. Create new file in `cli/cmd/`
2. Implement command using Cobra
3. Register in `cli/cmd/root.go`
4. Add tests in `cli/cmd/<command>_test.go`
5. Update documentation

### Updating Dependencies

```bash
# Update all dependencies
make update-deps

# Update specific dependency
go get -u github.com/spf13/cobra

# Tidy dependencies
go mod tidy
```

### Running CI Locally

```bash
# Run full CI pipeline locally
lefthook run pre-push

# Or manually
make lint
make test
make build
```

## Troubleshooting

### Common Issues

#### Go Module Errors
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download
```

#### Build Failures
```bash
# Clean build artifacts
make clean

# Rebuild from scratch
make all
```

#### Permission Errors
```bash
# Fix permissions for install
sudo make install

# Or install to user directory
make install PREFIX=~/.local
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DDX_HOME` | DDx installation directory | `~/.ddx` |
| `DDX_CONFIG` | Configuration file path | `~/.ddx/config.yml` |
| `DDX_DEBUG` | Enable debug logging | `false` |
| `GOOS` | Target OS for build | Current OS |
| `GOARCH` | Target architecture | Current arch |

## Related Documentation

- [[development/contributing/guidelines]] - Contributing guidelines
- [[development/standards/coding-standards]] - Coding standards
- [[development/testing/strategy]] - Testing strategy
- [[architecture/cli-architecture]] - CLI architecture details