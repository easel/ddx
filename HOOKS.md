# Git Hooks with Lefthook

This project uses [Lefthook](https://github.com/evilmartians/lefthook) for managing git hooks. Lefthook is a fast, cross-platform Git hooks manager that works on Windows, macOS, and Linux without requiring shell scripts.

## Installation

### 1. Install Lefthook

Choose one of the following installation methods:

#### Using Go (recommended for this project)
```bash
go install github.com/evilmartians/lefthook@latest
```

#### Using Homebrew (macOS/Linux)
```bash
brew install lefthook
```

#### Using npm
```bash
npm install -g lefthook
```

#### Using Windows (Scoop)
```powershell
scoop install lefthook
```

### 2. Install Git Hooks

After installing Lefthook, run this command in the project root:

```bash
lefthook install
```

This will set up the git hooks defined in `lefthook.yml`.

## Available Hooks

### Pre-commit Hooks

The following checks run automatically before each commit:

#### Fast Checks (always run)
- **conflicts**: Detects merge conflict markers
- **debug-statements**: Finds debug/print statements in code

#### Go-specific Checks (only when .go files are staged)
- **go-fmt**: Ensures Go code is properly formatted
- **go-lint**: Runs golangci-lint or go vet
- **go-build**: Verifies the code compiles
- **go-test**: Runs tests for changed packages

#### Security Checks
- **secrets**: Scans for hardcoded secrets and API keys
- **binaries**: Detects large files and unexpected binaries

#### Project-specific Checks
- **ddx-validate**: Validates DDx configuration when relevant files change

### Pre-push Hooks
- **tests**: Runs the full test suite before pushing

## Usage

### Normal Workflow

Just commit as usual - hooks will run automatically:

```bash
git add .
git commit -m "Your commit message"
```

### Skipping Hooks

#### Skip specific hooks
```bash
LEFTHOOK_EXCLUDE=go-test,go-build git commit -m "Quick fix"
```

#### Skip all hooks
```bash
git commit --no-verify -m "Emergency fix"
# or
LEFTHOOK=0 git commit -m "Emergency fix"
```

### Local Configuration

You can customize hook behavior locally without affecting other developers:

1. Copy the example configuration:
```bash
cp .lefthook-local.yml.example .lefthook-local.yml
```

2. Edit `.lefthook-local.yml` to customize hooks (this file is gitignored)

Common customizations:
```yaml
# Skip slow checks during development
pre-commit:
  commands:
    go-test:
      skip: true
    go-build:
      skip: true
```

## Troubleshooting

### Hooks not running?

1. Ensure Lefthook is installed: `lefthook version`
2. Reinstall hooks: `lefthook install -f`
3. Check configuration: `lefthook run pre-commit`

### False positives in security checks?

1. For secrets detection, you can:
   - Use environment variables instead of hardcoding
   - Add exceptions to the patterns in `lefthook.yml`
   - Skip the check temporarily: `LEFTHOOK_EXCLUDE=secrets git commit ...`

2. For binary detection:
   - Add binary files to `.gitignore`
   - Use Git LFS for large files
   - Skip the check: `LEFTHOOK_EXCLUDE=binaries git commit ...`

### Performance issues?

1. Use local configuration to skip slow checks during development
2. Run hooks manually when needed: `lefthook run pre-commit`
3. Use `--no-verify` for emergency commits

## Manual Hook Execution

You can run hooks manually without committing:

```bash
# Run all pre-commit hooks
lefthook run pre-commit

# Run specific hook
lefthook run pre-commit --commands go-test

# Run with verbose output
lefthook run pre-commit --verbose
```

## Uninstalling

To remove git hooks:

```bash
lefthook uninstall
```

## Benefits over Shell Scripts

1. **Cross-platform**: Works natively on Windows, macOS, and Linux
2. **Fast**: Parallel execution of independent checks
3. **Smart**: Only runs relevant checks based on staged files
4. **Configurable**: Easy to customize without modifying scripts
5. **Maintainable**: Single YAML configuration file
6. **No dependencies**: Doesn't require bash, WSL, or Git Bash on Windows

## Contributing

When adding new checks:

1. Edit `lefthook.yml` to add your check
2. Test locally: `lefthook run pre-commit`
3. Document any new requirements in this file
4. Consider adding skip conditions for development workflow