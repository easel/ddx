# DDx - Document-Driven Development eXperience

> A toolkit for AI-assisted development that helps you share templates, prompts, and patterns across projects.

DDx (Document-Driven Development eXperience) is like having a differential diagnosis system for your development workflow - it helps you quickly identify, apply, and share the right patterns, templates, and AI prompts across all your projects.

## 🚀 Quick Start

**One-line installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/easel/ddx/main/install.sh | bash
```

**Initialize in your project:**
```bash
cd your-project
ddx init
```

**See what's available:**
```bash
ddx list
ddx diagnose
```

## ✨ Features

- **🔄 Cross-Project Sync**: Share improvements across all your projects automatically
- **🤖 AI Integration**: Curated prompts and patterns for AI-assisted development
- **📋 Templates**: Project templates, code patterns, and configuration files
- **🔍 Project Diagnosis**: Analyze your setup and get improvement suggestions
- **⚡ Zero Dependencies**: Single binary, works everywhere (Mac, Linux, Windows)
- **📦 Git Integration**: Built on git subtree for reliable version control

## 🎯 Core Concepts

### Document-Driven Development
Like medical differential diagnosis (DDx), we use structured documentation to:
- **Diagnose** project issues and missing components
- **Prescribe** appropriate templates and patterns
- **Document** decisions and architectural choices
- **Share** improvements back to the community

### Medical Metaphor
- **Symptoms**: Project pain points and inefficiencies
- **Diagnosis**: `ddx diagnose` analyzes your project health
- **Treatment**: `ddx apply` prescribes and applies solutions
- **Rounds**: `ddx update` keeps your toolkit current with latest practices

## 📚 Commands

DDx follows a noun-verb command structure for better organization and discoverability:

### Core Commands
| Command | Description |
|---------|-------------|
| `ddx init` | Initialize DDx in current project |
| `ddx diagnose` | Analyze project and suggest improvements |
| `ddx update` | Update toolkit from master repository |
| `ddx contribute` | Share improvements back to community |

### Resource Commands (noun-verb structure)
| Command | Description |
|---------|-------------|
| **Prompts** | |
| `ddx prompts list` | List available AI prompts |
| `ddx prompts show <name>` | Display a specific prompt |
| **Templates** | |
| `ddx templates list` | List available project templates |
| `ddx templates apply <name>` | Apply a project template |
| **Patterns** | |
| `ddx patterns list` | List available code patterns |
| `ddx patterns apply <name>` | Apply a code pattern |
| **Personas** | |
| `ddx persona list` | List available AI personas |
| `ddx persona show <name>` | Show persona details |
| `ddx persona bind <role> <name>` | Bind persona to role |
| `ddx persona load` | Load personas into CLAUDE.md |
| **MCP Servers** | |
| `ddx mcp list` | List available MCP servers |
| `ddx mcp show <name>` | Show MCP server details |
| **Workflows** | |
| `ddx workflows list` | List available workflows |
| `ddx workflows show <name>` | Show workflow details |

## 📖 Usage Examples

```bash
# Start a new Next.js project with DDx patterns
ddx init --template nextjs

# List and apply templates
ddx templates list
ddx templates apply nextjs

# Browse and use AI prompts
ddx prompts list
ddx prompts show claude/code-review
ddx prompts list --verbose  # See all prompt files

# Work with code patterns
ddx patterns list
ddx patterns apply error-handling

# Manage AI personas
ddx persona list
ddx persona bind code-reviewer strict-code-reviewer
ddx persona load

# Check your project health
ddx diagnose --fix

# Share your improvements
ddx contribute patterns/my-new-pattern
```

## 🏗️ Project Structure

```
ddx/
├── library/           # DDx library resources
│   ├── prompts/       # AI prompts and instructions
│   │   ├── claude/    # Claude-specific prompts
│   │   └── general/   # Model-agnostic prompts
│   ├── templates/     # Project templates
│   │   ├── nextjs/    # Next.js starter
│   │   ├── python/    # Python projects
│   │   └── common/    # Common files (.gitignore, etc.)
│   ├── patterns/      # Code patterns and examples
│   │   ├── error-handling/
│   │   ├── testing/
│   │   └── ai-integration/
│   ├── personas/      # AI personality definitions
│   │   └── *.md       # Persona files
│   ├── mcp-servers/   # MCP server configurations
│   │   ├── registry.yml
│   │   └── servers/
│   └── configs/       # Tool configurations
│       ├── eslint/
│       ├── prettier/
│       └── typescript/
├── cli/               # DDx CLI implementation
├── docs/              # Documentation
├── scripts/           # Build and setup scripts
├── workflows/         # HELIX workflow definitions
└── install.sh         # One-line installer
```

## ⚙️ Configuration

DDx uses `.ddx.yml` files for configuration:

```yaml
version: 1.0
includes:
  - prompts/claude
  - scripts/hooks
  - templates/common
variables:
  project_name: "my-project"
  ai_model: "claude-3-opus"
overrides:
  "prompts/custom.md": "local/my-prompt.md"
```

### Library Path Resolution

DDx looks for library resources in the following order:

1. **Command flag**: `ddx --library-base-path /custom/path [command]`
2. **Environment variable**: `DDX_LIBRARY_BASE_PATH=/path ddx [command]`
3. **Development mode**: `<git-repo>/library/` when in DDx repository
4. **Project library**: `.ddx/library/` in current or parent directory
5. **Global library**: `~/.ddx/library/` (default installation)

This allows flexible testing and project-specific customization.

## 🔄 Git Subtree Integration

DDx uses git subtree to:
- Pull updates from the master toolkit
- Allow local modifications
- Contribute improvements back
- Maintain full version control history

```bash
# Manual subtree operations (handled by CLI)
git subtree add --prefix=.ddx https://github.com/easel/ddx main --squash
git subtree pull --prefix=.ddx https://github.com/easel/ddx main --squash
```

## 🤝 Contributing

### Adding New Resources

1. **Fork** this repository
2. **Add** your template, pattern, or prompt
3. **Test** with `ddx apply your-resource`
4. **Submit** a pull request

### Sharing from Projects

```bash
# From within a project
ddx contribute patterns/my-pattern
# Creates a branch and PR automatically
```

### Resource Guidelines

- **Templates**: Include README.md and clear variable substitution
- **Patterns**: Provide examples in multiple languages where applicable  
- **Prompts**: Include context about when and how to use
- **Scripts**: Make them cross-platform compatible

## 🏥 Health Check

Run `ddx diagnose` to check:
- ✅ DDx installation and configuration
- ✅ Git repository setup
- ✅ Project structure and conventions
- ✅ AI integration (CLAUDE.md, prompts, etc.)
- ✅ Development tooling (linting, testing, etc.)

Score: `85/100` ⭐ Excellent project health!

## 🚑 Emergency Procedures

### Broken Installation
```bash
# Reinstall DDx
curl -fsSL https://raw.githubusercontent.com/easel/ddx/main/install.sh | bash

# Reset project DDx
rm -rf .ddx
ddx init
```

### Merge Conflicts
```bash
# Reset to master state
ddx update --reset

# Or merge manually
git subtree pull --prefix=.ddx https://github.com/easel/ddx main
```

## 📋 Roadmap

- [ ] **VS Code Extension**: Integrated DDx commands
- [ ] **GitHub Actions**: Automated DDx updates and health checks  
- [ ] **Plugin System**: Custom resource types and processors
- [ ] **Community Hub**: Browse and share resources online
- [ ] **AI Model Adapters**: Support for different AI providers
- [ ] **Template Marketplace**: Curated collection of community templates

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

## 🏆 Recognition

Inspired by:
- Medical differential diagnosis methodology
- Infrastructure as Code principles
- Developer experience best practices
- Open source collaboration patterns

---

**Made with ❤️ for the AI-assisted development community**

*"Like a medical DDx system, but for your code - helping you diagnose issues and prescribe the right solutions."*