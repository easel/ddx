# API Contract: Core CLI Commands [FEAT-001]

**Contract ID**: CLI-001
**Feature**: FEAT-001
**Type**: CLI
**Status**: Approved
**Version**: 2.0.0
**Updated**: 2025-01-18

*Command-line interface contracts for DDx core commands*

## Command: init

**Purpose**: Initialize DDx in a project with optional template application
**Usage**: `ddx init [template] [options]`

### Arguments
- `template` (optional): Name of template to apply during initialization (e.g., "nextjs", "python-flask")

### Options
- `--repo, -r <url>`: Master repository URL (default: https://github.com/ddx-tools/ddx)
- `--branch, -b <name>`: Repository branch (default: main)
- `--path, -p <path>`: Project path (default: current directory)
- `--force, -f`: Overwrite existing .ddx.yml configuration
- `--no-git`: Skip git subtree setup
- `--minimal`: Create minimal configuration without subtree

### Input
- Interactive prompts if configuration values not provided
- Reads existing .ddx.yml if present (unless --force)

### Output Format
```
✓ Initialized DDx in [path]
✓ Created .ddx.yml configuration
✓ Set up git subtree at .ddx/
✓ Applied template: [template-name] (if specified)

Project ready for document-driven development!
Run 'ddx prompts list' or 'ddx templates list' to see available resources.
```

### Exit Codes
- `0`: Success
- `1`: General error (permission denied, invalid path)
- `2`: Configuration already exists (use --force to overwrite)
- `3`: Git repository required but not found
- `4`: Template not found
- `5`: Network error fetching repository

### Examples
```bash
# Basic initialization
$ ddx init
✓ Initialized DDx in /current/directory
✓ Created .ddx.yml configuration
✓ Set up git subtree at .ddx/

# Initialize with template
$ ddx init nextjs
✓ Initialized DDx in /current/directory
✓ Created .ddx.yml configuration
✓ Set up git subtree at .ddx/
✓ Applied template: nextjs

# Custom repository
$ ddx init --repo https://github.com/myorg/ddx-custom --branch develop
✓ Initialized DDx with custom repository
```

### Error Handling
```bash
# Already initialized
$ ddx init
Error: .ddx.yml already exists. Use --force to overwrite.

# Not a git repository
$ ddx init
Error: Current directory is not a git repository. Initialize git first or use --no-git option.

# Template not found
$ ddx init unknown-template
Error: Template 'unknown-template' not found. Run 'ddx list templates' to see available templates.
```

---

## Resource Commands (Noun-Verb Structure)

DDx follows a noun-verb command structure for resource management. Each resource type has its own command namespace with consistent actions.

### Command Pattern
**Usage**: `ddx <resource> <action> [name] [options]`

Resources: `prompts`, `templates`, `patterns`, `persona`, `mcp`, `workflows`
Actions: `list`, `show`, `apply` (where applicable)

---

## Command: prompts

**Purpose**: Manage AI prompts and instructions
**Usage**: `ddx prompts <action> [name] [options]`

### Subcommand: prompts list

**Purpose**: List available AI prompts
**Usage**: `ddx prompts list [options]`

### Options
- `--search, -s <term>`: Search for prompts containing term
- `--verbose, -v`: Show all prompt files recursively (not just directories)
- `--format <format>`: Output format: "table" (default), "json", "yaml"

### Output Format
```
Available Prompts:

claude/             Claude-specific prompts (2 files)
common/             General AI prompts (5 files)
ddx/                DDx workflow prompts (1 file)
```

### Output Format (Verbose)
```
Available Prompts:

claude/
  code-review.md    Code review prompt for Claude
  system-prompts/
    security.md     Security-focused review prompt
common/
  refactor.md       Refactoring assistance prompt
  docs.md           Documentation generation prompt
```

### Examples
```bash
# List all prompts
$ ddx prompts list
Available Prompts:
claude/             Claude-specific prompts (2 files)
common/             General AI prompts (5 files)

# List with file details
$ ddx prompts list --verbose
Available Prompts:
claude/
  code-review.md    Code review prompt for Claude
  system-prompts/
    security.md     Security-focused review prompt

# Search for specific prompts
$ ddx prompts list --search review
Matching Prompts:
claude/code-review.md
claude/system-prompts/security.md
```

### Subcommand: prompts show

**Purpose**: Display a specific prompt
**Usage**: `ddx prompts show <name>`

### Arguments
- `name` (required): Path to prompt file (e.g., "claude/code-review")

### Output Format
```
Prompt: claude/code-review.md
────────────────────────────────
# Code Review Prompt

You are a senior code reviewer...
[prompt content]
```

---

## Command: templates

**Purpose**: Manage project templates
**Usage**: `ddx templates <action> [name] [options]`

### Subcommand: templates list

**Purpose**: List available project templates
**Usage**: `ddx templates list [options]`

### Options
- `--search, -s <term>`: Search for templates
- `--tags, -t <tags>`: Filter by tags
- `--verbose, -v`: Show detailed information

### Output Format
```
Available Templates:

Name            Description                          Tags
────────────────────────────────────────────────────────────
nextjs          Next.js 14 with TypeScript          react, web
python-flask    Flask API with best practices       python, api
go-cli          Go CLI with Cobra framework         go, cli
```

### Subcommand: templates apply

**Purpose**: Apply a template to the current project
**Usage**: `ddx templates apply <name> [options]`

### Arguments
- `name` (required): Template name (e.g., "nextjs")

### Options
- `--force, -f`: Overwrite existing files
- `--dry-run`: Preview changes without applying
- `--variables, -v <key=value>`: Set template variables

### Examples
```bash
# List templates
$ ddx templates list

# Apply a template
$ ddx templates apply nextjs
✓ Applied template: nextjs
✓ Created 15 files
✓ Updated package.json

# Preview changes
$ ddx templates apply python-flask --dry-run
```

---

## Command: patterns

**Purpose**: Manage code patterns
**Usage**: `ddx patterns <action> [name] [options]`

### Subcommand: patterns list

**Purpose**: List available code patterns
**Usage**: `ddx patterns list [options]`

### Options
- `--search, -s <term>`: Search for patterns
- `--category, -c <name>`: Filter by category
- `--verbose, -v`: Show detailed information

### Subcommand: patterns apply

**Purpose**: Apply a code pattern to the project
**Usage**: `ddx patterns apply <name> [options]`

### Arguments
- `name` (required): Pattern name (e.g., "error-handling")

### Options
- `--target, -t <path>`: Target directory for pattern
- `--force, -f`: Overwrite existing files

---

## Command: persona

**Purpose**: Manage AI personas for consistent interactions
**Usage**: `ddx persona <action> [arguments] [options]`

### Subcommand: persona list

**Purpose**: List available AI personas
**Usage**: `ddx persona list [options]`

### Options
- `--role <role>`: Filter by role capability
- `--tag <tag>`: Filter by tags

### Output Format
```
Available Personas:

NAME                   ROLES                             DESCRIPTION
────────────────────────────────────────────────────────────────────
strict-code-reviewer   code-reviewer, security-analyst   Uncompromising code quality enforcer
test-engineer-tdd      test-engineer, quality-analyst    Test-driven development specialist
```

### Subcommand: persona bind

**Purpose**: Bind a persona to a role in the project
**Usage**: `ddx persona bind <role> <persona-name>`

### Arguments
- `role` (required): Role to bind (e.g., "code-reviewer")
- `persona-name` (required): Persona to bind (e.g., "strict-code-reviewer")

### Subcommand: persona load

**Purpose**: Load bound personas into CLAUDE.md
**Usage**: `ddx persona load [options]`

### Options
- `--role <role>`: Load only specific role

---

## Legacy Command: apply (Deprecated)

**Purpose**: Apply DDx resources (templates, patterns, prompts, configs) to the project
**Usage**: `ddx apply <resource-path> [options]`

### Arguments
- `resource-path` (required): Path to resource (e.g., "templates/nextjs", "patterns/auth-jwt", "prompts/claude/refactor")

### Options
- `--target, -t <path>`: Target directory for application (default: current directory)
- `--vars <key=value>`: Variable substitutions (can be repeated)
- `--vars-file <file>`: YAML/JSON file with variables
- `--force, -f`: Overwrite existing files without prompting
- `--dry-run`: Preview changes without applying
- `--interactive, -i`: Interactive mode for variable input
- `--merge-strategy <strategy>`: How to handle conflicts: "skip", "overwrite", "merge", "prompt" (default: prompt)

### Input
- Template variables from --vars, --vars-file, or interactive prompts
- User confirmation for file overwrites (unless --force)

### Output Format
```
Applying template: nextjs
Variables:
  project_name: my-app
  port: 3000

Files to create/modify:
  ✓ package.json (new)
  ✓ tsconfig.json (new)
  ⚠ README.md (exists - will merge)
  ✓ src/app/page.tsx (new)

Continue? (y/n): y

Applying changes...
  ✓ Created package.json
  ✓ Created tsconfig.json
  ✓ Merged README.md (3 conflicts resolved)
  ✓ Created src/app/page.tsx

Successfully applied template: nextjs
4 files created/modified
Run 'npm install' to install dependencies
```

### Exit Codes
- `0`: Success
- `1`: Resource not found
- `2`: Invalid resource path
- `3`: Variable validation failed
- `4`: File conflict (user cancelled)
- `5`: Permission denied
- `6`: Template syntax error

### Examples
```bash
# Apply template with variables
$ ddx apply templates/nextjs --vars project_name=my-app --vars port=3001
Applying template: nextjs
[output as above]

# Dry run to preview changes
$ ddx apply patterns/auth-jwt --dry-run
[DRY RUN] Would apply pattern: auth-jwt
Files that would be created:
  - src/auth/jwt.js
  - src/auth/middleware.js
  - src/auth/config.js

# Interactive mode
$ ddx apply templates/python-flask -i
Applying template: python-flask
Enter value for 'app_name' (default: app): my_api
Enter value for 'port' (default: 5000): 8080
Enter value for 'database' (postgres/mysql/sqlite) [sqlite]: postgres

# Apply with merge strategy
$ ddx apply configs/eslint --merge-strategy merge
Merging ESLint configuration...
✓ Merged .eslintrc.json (combined rules)

# From variables file
$ ddx apply templates/nextjs --vars-file project.yml
Loading variables from project.yml...
Applied template with 5 variables
```

### Variable Substitution
Variables in templates use `{{variable_name}}` syntax:
```javascript
// Before substitution
const PORT = {{port || 3000}};
const APP_NAME = "{{project_name}}";

// After substitution (with project_name=my-app, port=3001)
const PORT = 3001;
const APP_NAME = "my-app";
```

### Merge Strategies
- **skip**: Skip files that already exist
- **overwrite**: Replace existing files completely
- **merge**: Attempt intelligent merge (JSON, YAML, code)
- **prompt**: Ask user for each conflict (default)

### Error Handling
```bash
# Resource not found
$ ddx apply templates/unknown
Error: Resource 'templates/unknown' not found. Run 'ddx list templates' to see available templates.

# Missing required variables
$ ddx apply templates/nextjs
Error: Required variable 'project_name' not provided.
Use --vars project_name=value or --interactive mode.

# Permission denied
$ ddx apply configs/eslint --target /system
Error: Permission denied writing to /system/.eslintrc.json

# Conflict without force
$ ddx apply patterns/auth-jwt
Warning: Following files already exist:
  - src/auth/config.js
  - src/middleware/auth.js
Use --force to overwrite or --merge-strategy to specify handling.
```

## Contract Validation

All commands must:
1. Return consistent exit codes
2. Support --help flag for documentation
3. Support --version flag for version info
4. Output errors to stderr, normal output to stdout
5. Support JSON output format for automation
6. Validate all inputs before execution
7. Provide clear error messages with recovery suggestions

## Global Options

Available for all commands:
- `--config <file>`: Use alternate config file (default: .ddx.yml)
- `--verbose, -v`: Verbose output
- `--quiet, -q`: Suppress non-error output
- `--no-color`: Disable colored output
- `--help, -h`: Show help message
- `--version`: Show version information

## Environment Variables

- `DDX_CONFIG`: Path to configuration file
- `DDX_REPO`: Default repository URL
- `DDX_BRANCH`: Default branch name
- `DDX_NO_COLOR`: Disable colored output
- `DDX_DEBUG`: Enable debug logging