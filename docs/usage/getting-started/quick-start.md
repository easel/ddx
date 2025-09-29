# Getting Started with DDx

> **Last Updated**: 2025-09-11
> **Status**: Active
> **Owner**: DDx Team

## Overview

Welcome to DDx (Document-Driven Development eXperience)! This guide will help you get up and running quickly.

## Context

This quick start guide is for developers who want to immediately start using DDx in their projects. For more detailed setup instructions, see [[development/tools/setup]].

## Installation

### Quick Install (Recommended)

Run this one-line command to install DDx:

```bash
curl -fsSL https://raw.githubusercontent.com/easel/ddx/main/install.sh | bash
```

This will:
- Download and install the DDx CLI
- Set up shell completions
- Create default configuration
- Add DDx to your PATH

### Manual Install

1. Download the appropriate binary for your platform from the [releases page](https://github.com/easel/ddx/releases)
2. Extract and move to a directory in your PATH:
   ```bash
   tar -xzf ddx-linux-amd64.tar.gz
   sudo mv ddx /usr/local/bin/
   ```
3. Verify installation:
   ```bash
   ddx version
   ```

### Build from Source

If you have Go 1.21+ installed:

```bash
git clone https://github.com/easel/ddx
cd ddx/cli
make build
make install
```

## First Steps

### 1. Initialize DDx in a Project

Navigate to your project directory and run:

```bash
cd your-project
ddx init
```

This creates:
- `.ddx/` directory with toolkit resources
- `.ddx.yml` configuration file
- Copies of selected prompts, templates, and scripts

### 2. Explore Available Resources

See what's available in your toolkit:

```bash
ddx list
```

Filter by type:
```bash
ddx list --type templates
ddx list --type patterns
ddx list --search react
```

### 3. Apply Resources

Apply a template or pattern to your project:

```bash
# Apply a project template
ddx apply nextjs

# Apply a code pattern
ddx apply error-handling

# Apply AI prompts
ddx apply prompts/claude
```

### 4. Check Project Health

Run a diagnostic to see how your project is doing:

```bash
ddx doctor
```

Get a detailed report:
```bash
ddx doctor --report
```

Auto-fix common issues:
```bash
ddx doctor --fix
```

## Understanding DDx Structure

### Local Project Structure

After running `ddx init`, your project will have:

```
your-project/
├── .ddx/                  # DDx resources (synced from master)
│   ├── prompts/
│   ├── templates/
│   ├── patterns/
│   └── scripts/
├── .ddx.yml              # Local configuration
└── your existing files...
```

### Configuration

The `.ddx.yml` file controls what resources are included and how they're configured:

```yaml
version: 1.0
includes:
  - prompts/claude         # Include Claude AI prompts
  - scripts/hooks          # Include git hooks
  - templates/common       # Include common templates
variables:
  project_name: "my-app"   # Used in template substitution
  ai_model: "claude-3-opus"
overrides:
  "prompts/custom.md": "local/my-custom-prompt.md"
```

### Template Variables

Templates use variable substitution with `{{variable_name}}` syntax:

```markdown
# {{project_name}}

This is a {{project_type}} project using {{tech_stack}}.
```

Variables can be:
- Defined in `.ddx.yml`
- Detected automatically (like project name from directory)
- Prompted for during `ddx apply`

## Common Workflows

### Starting a New Project

```bash
# Create project directory
mkdir my-new-project
cd my-new-project

# Initialize git
git init

# Initialize DDx with a template
ddx init --template nextjs

# Check everything looks good
ddx doctor
```

### Adding DDx to Existing Project

```bash
# Navigate to existing project
cd existing-project

# Initialize DDx (will detect project type)
ddx init

# Apply relevant patterns
ddx apply error-handling
ddx apply testing

# Set up git hooks
ddx apply scripts/hooks
```

### Keeping Up to Date

```bash
# Update your local toolkit
ddx update

# Check for issues after update
ddx doctor

# See what changed
git log --oneline .ddx/
```

### Contributing Back

```bash
# Create a new pattern
mkdir .ddx/patterns/my-pattern
# ... add files ...

# Contribute it back
ddx contribute patterns/my-pattern
```

## Best Practices

### 1. Version Control

- **Commit** the `.ddx.yml` configuration file
- **Consider** committing the `.ddx/` directory for reproducible builds
- **Don't commit** temporary DDx files (they're in the default .gitignore)

### 2. Team Adoption

- Use `ddx init` to ensure consistent setup across team members
- Share custom configurations via `.ddx.yml`
- Use `ddx doctor` in CI to enforce project health

### 3. Customization

- Use `overrides` in `.ddx.yml` for project-specific modifications
- Create local templates and patterns in your project
- Contribute improvements back to benefit the community

### 4. AI Integration

- Keep `CLAUDE.md` updated with project context
- Use DDx prompts as starting points for AI conversations
- Document AI usage patterns for your team

## Troubleshooting

### Common Issues

**DDx command not found:**
```bash
# Check PATH
echo $PATH
# Restart shell or source rc file
source ~/.bashrc  # or ~/.zshrc
```

**Permission denied:**
```bash
# Make sure binary is executable
chmod +x ~/.local/bin/ddx
```

**Template variables not substituted:**
```bash
# Check variable definitions in .ddx.yml
ddx config --show
```

**Git subtree conflicts:**
```bash
# Reset to clean state
ddx update --reset
```

### Getting Help

- Run `ddx help` for command information
- Use `ddx doctor` to check for configuration issues
- Check the [documentation](../README.md) for detailed guides
- Open an issue on [GitHub](https://github.com/easel/ddx/issues)

## Next Steps

Now that you have DDx set up:

1. [Explore available templates](templates.md)
2. [Learn about patterns](patterns.md) 
3. [Configure AI integration](ai-integration.md)
4. [Contribute your own resources](contributing.md)

Happy coding with DDx! 🚀