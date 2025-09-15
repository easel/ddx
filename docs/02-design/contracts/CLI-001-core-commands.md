# API Contract: Core CLI Commands [FEAT-001]

**Contract ID**: CLI-001
**Feature**: FEAT-001
**Type**: CLI
**Status**: Approved
**Version**: 1.0.0

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
Run 'ddx list' to see available resources.
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

## Command: list

**Purpose**: Display available DDx resources (templates, patterns, prompts, configs)
**Usage**: `ddx list [resource-type] [options]`

### Arguments
- `resource-type` (optional): Type to list - "templates", "patterns", "prompts", "configs", or "all" (default: all)

### Options
- `--filter, -f <query>`: Filter results by name or tag
- `--tags, -t <tags>`: Filter by comma-separated tags
- `--format <format>`: Output format: "table" (default), "json", "yaml", "simple"
- `--verbose, -v`: Show detailed information
- `--local`: Show only locally available resources
- `--remote`: Show only remote resources

### Output Format (Table - Default)
```
Available DDx Resources:

TEMPLATES (5)
Name            Description                          Tags           Status
─────────────────────────────────────────────────────────────────────────
nextjs          Next.js 14 with TypeScript          react, web     Available
python-flask    Flask API with best practices       python, api    Available
go-cli          Go CLI with Cobra framework         go, cli        Local only

PATTERNS (12)
Name            Description                          Tags           Status
─────────────────────────────────────────────────────────────────────────
auth-jwt        JWT authentication pattern          auth, security Available
error-handler   Centralized error handling          errors         Available
```

### Output Format (JSON)
```json
{
  "templates": [
    {
      "name": "nextjs",
      "description": "Next.js 14 with TypeScript",
      "tags": ["react", "web"],
      "status": "available",
      "path": ".ddx/templates/nextjs"
    }
  ],
  "patterns": [...],
  "prompts": [...],
  "configs": [...]
}
```

### Exit Codes
- `0`: Success
- `1`: Invalid resource type
- `2`: Filter syntax error
- `3`: Configuration error

### Examples
```bash
# List all resources
$ ddx list
Available DDx Resources:
[output as above]

# List only templates
$ ddx list templates
TEMPLATES (5)
nextjs          Next.js 14 with TypeScript          react, web     Available
python-flask    Flask API with best practices       python, api    Available

# Filter by tag
$ ddx list --tags react,typescript
Filtered Results (tag: react, typescript):
nextjs          Next.js 14 with TypeScript          react, web     Available
react-component React component pattern              react          Available

# JSON output
$ ddx list templates --format json
{"templates": [{"name": "nextjs", ...}]}

# Verbose output
$ ddx list patterns auth-jwt -v
PATTERN: auth-jwt
Description: JWT authentication pattern with refresh tokens
Tags: auth, security, jwt
Status: Available
Path: .ddx/patterns/auth/jwt
Files: 5 files, 2.3 KB
Dependencies: jsonwebtoken, bcrypt
Last Updated: 2025-01-10
```

---

## Command: apply

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