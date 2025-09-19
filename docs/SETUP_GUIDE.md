# Complete Setup Guide for DDX and AI-Assisted Development

> From zero to productive: Setting up your development environment for DDX and modern AI-assisted coding

## Prerequisites Overview

Before we begin, here's what we'll be setting up:
- ✅ **Homebrew** - Package manager for installing tools
- ✅ **Terminal emulator** - Your command-line interface
- ✅ **Git** - Version control system
- ✅ **Container runtime** - Docker or Podman for isolated environments
- ✅ **Claude Code** - AI-powered development assistant
- ✅ **Language runtimes** - Python, Node.js, Go, or Rust (as needed)
- ✅ **DDX** - Document-Driven Development toolkit

**Time needed**: 45-60 minutes for complete setup

---

## Part 1: Installing Homebrew

[Homebrew](https://brew.sh) is the missing package manager for macOS and Linux. Think of it as an app store for command-line tools. We'll install this first, then use it to install everything else.

### Quick Start Terminal

To run the Homebrew installation, you'll need a terminal. Use your system's default for now:

**macOS**: Press `Cmd + Space`, type "Terminal", press Enter
**Windows**: Open PowerShell as Administrator (we'll set up WSL2 later)
**Linux**: Open your default terminal application

### macOS and Linux

Visit [brew.sh](https://brew.sh) for the official installation instructions, or run this command:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

**What this does:**
- Downloads the Homebrew installer
- Sets up Homebrew in the correct location
- Configures your system to use Homebrew

### Post-Installation Setup

After installation completes, you'll see instructions. Run these commands:

**On macOS (Apple Silicon):**
```bash
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

**On macOS (Intel):**
```bash
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/usr/local/bin/brew shellenv)"
```

**On Linux:**
```bash
echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.profile
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
```

### Verify Installation
```bash
brew --version
# Should output: Homebrew 4.x.x
```

✅ **Success indicator**: You see a version number

❌ **If it fails**: Close terminal, reopen, and try `brew --version` again

### Windows (Using WSL2)

In your Ubuntu/WSL2 terminal:

```bash
# Update package list
sudo apt update

# Install dependencies
sudo apt-get install build-essential curl file git

# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Add to path
echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> ~/.profile
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
```

### Windows (Alternative: Chocolatey)

If you prefer native Windows tools, use Chocolatey:

In PowerShell (as Administrator):
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

---

## Part 2: Modern Terminal Emulator

Now that you have Homebrew installed, let's upgrade to a modern terminal emulator with better features, performance, and customization options.

### macOS

#### Option 1: Ghostty (Primary Recommendation)
[Ghostty](https://ghostty.org) is a fast, feature-rich, native terminal emulator with excellent performance:

```bash
# Install with Homebrew
brew install --cask ghostty

# Launch Ghostty
open -a Ghostty
```

**Why Ghostty:**
- ⚡ Native performance (written in Zig)
- 🎨 Beautiful rendering with proper font support
- ⌨️ Excellent keyboard shortcuts
- 🔧 Simple configuration

#### Option 2: Kitty (Alternative Recommendation)
[Kitty](https://sw.kovidgoyal.net/kitty/) is a GPU-accelerated terminal with advanced features:

```bash
# Install with Homebrew
brew install --cask kitty

# Launch Kitty
open -a kitty
```

**Why Kitty:**
- 🚀 GPU-accelerated rendering
- 🖼️ Image and graphics support
- 📑 Tabs and splits
- 🎯 Extensive customization

#### Option 3: iTerm2 (Traditional Choice)
```bash
brew install --cask iterm2
```

#### Option 4: Warp (AI-Enhanced)
```bash
brew install --cask warp
```

### Linux

#### Ghostty (Recommended)
```bash
# For Linux, build from source or check ghostty.org for packages
# Alternatively, use Kitty:
```

#### Kitty (Recommended)
```bash
# Install with package manager
curl -L https://sw.kovidgoyal.net/kitty/installer.sh | sh /dev/stdin

# Or via package manager (Ubuntu/Debian)
sudo apt install kitty

# Or via Homebrew on Linux
brew install --cask kitty
```

#### Alternative Options
```bash
# Alacritty (Rust-based, fast)
brew install --cask alacritty

# WezTerm (feature-rich)
brew install --cask wezterm
```

### Windows (WSL2)

For Windows users, after setting up WSL2, you can use Windows Terminal with WSL2, or install a better terminal:

#### Windows Terminal + WSL2 (Built-in)
Already configured in Part 1 if you followed the WSL2 setup.

#### Ghostty or Kitty for Windows
Check the official websites for Windows builds:
- [Ghostty for Windows](https://ghostty.org)
- [Kitty for Windows](https://sw.kovidgoyal.net/kitty/)

### Terminal Configuration Tips

After installing your terminal, consider these quick improvements:

```bash
# Install a better shell prompt (Starship)
brew install starship
echo 'eval "$(starship init zsh)"' >> ~/.zshrc  # For zsh
echo 'eval "$(starship init bash)"' >> ~/.bashrc # For bash

# Install useful terminal tools
brew install eza  # Better ls
brew install bat  # Better cat
brew install fd   # Better find
brew install ripgrep # Better grep
```

---

## Part 3: Essential Development Tools

### Git Installation and Configuration

Git is essential for version control and using DDX.

#### Install Git
```bash
# macOS/Linux with Homebrew
brew install git

# Ubuntu/Debian
sudo apt install git

# Fedora
sudo dnf install git

# Windows Chocolatey
choco install git
```

#### Configure Git
Set your identity (required for commits):
```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

#### Verify Installation
```bash
git --version
# Should output: git version 2.x.x
```

### GitHub CLI (gh) Installation and Setup

The GitHub CLI provides secure authentication and powerful GitHub features directly from your terminal.

#### Install GitHub CLI
```bash
# macOS/Linux with Homebrew
brew install gh

# Ubuntu/Debian
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh

# Fedora/RHEL
sudo dnf install gh

# Windows with Chocolatey
choco install gh

# Windows with Scoop
scoop install gh
```

#### Authenticate with GitHub
```bash
# Login to GitHub via OAuth in your browser
gh auth login

# Follow the prompts:
# - Choose GitHub.com
# - Choose HTTPS for git protocol
# - Authenticate with your web browser
# - The browser will open for OAuth authentication
```

#### Verify Authentication
```bash
# Check authentication status
gh auth status
# Should show: "✓ Logged in to github.com as username"

# Test with a simple command
gh repo list --limit 5
# Lists your first 5 repositories
```

#### Using gh with Git
Once authenticated, gh automatically configures git to use HTTPS with your credentials:
```bash
# Clone a repository
gh repo clone owner/repository

# Create a new repository
gh repo create my-new-project --public

# Create a pull request
gh pr create --title "My feature" --body "Description"

# Check PR status
gh pr status
```

#### Benefits of GitHub CLI
- 🔒 **Secure OAuth**: No need to manage SSH keys
- 🌐 **Works everywhere**: HTTPS works behind firewalls and proxies
- ⚡ **Powerful features**: PR management, issue tracking, releases
- 🔄 **Automatic auth**: Credentials handled seamlessly

---

## Part 4: Container Runtimes (Docker/Podman)

Container runtimes are essential for modern development, allowing you to run isolated environments, test deployments, and ensure consistency across different systems. DDX uses containers for various workflows and testing scenarios.

### Why You Need Container Runtime

- 🔧 **Isolated Environments**: Run different versions of tools without conflicts
- 📦 **Consistent Development**: Same environment across all team members
- 🚀 **Easy Deployment**: Test production configurations locally
- 🧪 **Safe Testing**: Experiment without affecting your system

### Docker (macOS and Windows)

Docker Desktop provides an easy-to-use container runtime with a GUI.

#### macOS Installation
```bash
# Install Docker Desktop with Homebrew
brew install --cask docker

# Launch Docker Desktop
open -a Docker

# Wait for Docker to start, then verify
docker --version
docker run hello-world
```

#### Windows Installation

**Option 1: Docker Desktop with WSL2 (Recommended)**
1. Ensure WSL2 is installed (from Part 1)
2. Download [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
3. Run the installer
4. Enable "Use WSL 2 based engine" during setup
5. Launch Docker Desktop from Start Menu

Verify installation in WSL2 terminal:
```bash
docker --version
docker run hello-world
```

**Option 2: Install via Chocolatey**
```powershell
# In PowerShell as Administrator
choco install docker-desktop
```

### Podman (Fedora/RHEL/Linux)

Podman is a daemonless container runtime that's Docker-compatible and doesn't require root privileges.

#### Fedora Installation
```bash
# Install Podman
sudo dnf install -y podman podman-compose

# Optional: Install podman-docker for Docker compatibility
sudo dnf install -y podman-docker

# Verify installation
podman --version
podman run hello-world

# Enable Docker compatibility (optional)
# This creates 'docker' as an alias to 'podman'
sudo touch /etc/containers/nodocker
```

#### Ubuntu/Debian Installation
```bash
# Add Kubic repository
. /etc/os-release
echo "deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/ /" | sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
curl -L "https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/Release.key" | sudo apt-key add -

# Install Podman
sudo apt update
sudo apt install -y podman

# Verify
podman --version
```

#### macOS (Alternative to Docker)
```bash
# Install Podman with Homebrew
brew install podman

# Initialize and start Podman machine
podman machine init
podman machine start

# Verify
podman --version
podman run hello-world
```

### Container Runtime Configuration

After installing your container runtime, configure it for optimal development:

#### Docker Configuration
```bash
# Add your user to docker group (Linux only)
sudo usermod -aG docker $USER
newgrp docker

# Configure Docker resources (via Docker Desktop GUI on Mac/Windows)
# Recommended: 4GB RAM, 2 CPUs minimum
```

#### Podman Configuration
```bash
# Enable user namespaces for rootless containers
echo "$USER:100000:65536" | sudo tee -a /etc/subuid
echo "$USER:100000:65536" | sudo tee -a /etc/subgid

# Set Podman as Docker replacement (optional)
echo 'alias docker=podman' >> ~/.bashrc
echo 'alias docker-compose=podman-compose' >> ~/.bashrc
source ~/.bashrc
```

### Testing Your Container Runtime

Run these commands to ensure everything works:

```bash
# Pull and run a test container
docker run --rm alpine echo "Containers are working!"
# or with Podman
podman run --rm alpine echo "Containers are working!"

# Run an interactive container
docker run -it --rm ubuntu:latest bash
# Type 'exit' to leave

# Check running containers
docker ps
# or
podman ps
```

### Container Tips for DDX

DDX can leverage containers for:
- Running different Node.js versions for MCP servers
- Testing in clean environments
- Deploying HELIX workflow applications
- Isolating development dependencies

```bash
# Example: Run Node.js in a container for MCP servers
docker run -it --rm -v "$PWD":/workspace -w /workspace node:18 npm install

# Example: Test DDX in a clean environment
docker run -it --rm -v "$PWD":/ddx -w /ddx golang:1.21 go test ./...
```

---

## Part 5: Installing Claude Code

[Claude Code](https://claude.ai/code) is an AI-powered development assistant that transforms how you write and understand code. It works seamlessly with DDX for enhanced AI-assisted development.

### What is Claude Code?

Claude Code provides:
- 🤖 **AI Pair Programming**: Real-time coding assistance and suggestions
- 📁 **Project Context**: Understands your entire codebase
- 🔧 **Code Generation**: Write functions, tests, and documentation
- 🔍 **Code Analysis**: Debug issues and understand complex code
- ♻️ **Refactoring**: Improve code quality and performance
- 🚀 **Multi-file Editing**: Make changes across your entire project

### Installation

#### macOS
```bash
# Install with Homebrew (when available)
# brew install --cask claude-code

# Or download directly:
# Visit https://claude.ai/download
# Download Claude-x.x.x-mac.zip
# Extract and move to Applications folder
```

#### Windows
```bash
# Download from https://claude.ai/download
# Run Claude-Setup-x.x.x.exe
# Follow the installation wizard
```

#### Linux
```bash
# Download AppImage from https://claude.ai/download
wget https://claude.ai/download/claude-x.x.x-linux.AppImage
chmod +x claude-x.x.x-linux.AppImage

# Run directly
./claude-x.x.x-linux.AppImage

# Or install system-wide
sudo mv claude-x.x.x-linux.AppImage /opt/Claude.AppImage
sudo ln -s /opt/Claude.AppImage /usr/local/bin/claude
```

### Initial Setup

1. **Launch Claude Code**
   - macOS: Open from Applications or `open -a "Claude"`
   - Windows: Launch from Start Menu
   - Linux: Run from terminal or application menu

2. **Sign In**
   - Create or sign in to your Anthropic account
   - Choose your subscription plan (Pro recommended for development)

3. **Configure Settings**
   ```
   Settings → Preferences:
   - Theme: Dark/Light/System
   - Editor Font: Adjust size and family
   - Model: Claude 3.5 Sonnet (recommended)
   - Context Window: Maximum for large projects
   ```

### Creating CLAUDE.md Files

CLAUDE.md files provide persistent context about your project. DDX can automatically manage these:

```bash
# DDX automatically creates CLAUDE.md with persona bindings
ddx persona load

# Or create manually with project context
cat > CLAUDE.md << 'EOF'
# Project Context for Claude Code

## Overview
This project uses DDX for development workflow management with HELIX methodology.

## Architecture
- **Language**: Go 1.21+
- **Framework**: Cobra CLI
- **Testing**: Go test with high coverage requirements
- **Deployment**: Docker containers

## Coding Standards
- Follow Go idioms and best practices
- Maintain 80%+ test coverage
- Use meaningful variable names
- Document all public APIs

## DDX Integration
- Personas defined in .ddx.yml
- HELIX workflow active in Design phase
- MCP servers configured for enhanced capabilities
EOF
```

### Keyboard Shortcuts

Essential shortcuts for productivity:

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| New Chat | `Cmd+N` | `Ctrl+N` |
| Open File | `Cmd+O` | `Ctrl+O` |
| Save Chat | `Cmd+S` | `Ctrl+S` |
| Search | `Cmd+F` | `Ctrl+F` |
| Settings | `Cmd+,` | `Ctrl+,` |
| Toggle Sidebar | `Cmd+B` | `Ctrl+B` |

### Project Setup Best Practices

#### 1. Add Project to Claude Code
```bash
# Open your project in Claude Code
cd your-project
open -a "Claude" .  # macOS
# or
claude .  # If symlinked
```

#### 2. Configure .clignore
Create a `.clignore` file to exclude files from context:
```bash
cat > .clignore << 'EOF'
# Dependencies
node_modules/
vendor/
.venv/

# Build outputs
dist/
build/
*.exe
*.dll

# Large files
*.log
*.sqlite
*.mp4
*.zip

# Secrets
.env
*.key
*.pem
EOF
```

#### 3. Use DDX Personas
```bash
# Bind personas for consistent AI behavior
ddx persona bind code-reviewer strict-code-reviewer
ddx persona bind architect systems-thinker
ddx persona load  # Updates CLAUDE.md
```

### MCP Server Integration

Claude Code supports MCP (Model Context Protocol) servers for enhanced capabilities:

```bash
# Install MCP servers with DDX
ddx mcp install filesystem      # File system access
ddx mcp install github          # GitHub integration
ddx mcp install sequential-thinking  # Advanced reasoning

# Servers are automatically configured in:
# .claude/settings.local.json
```

### Tips for Effective Use

#### 1. Project Context
- Keep CLAUDE.md updated with project changes
- Use clear folder structures
- Document non-obvious design decisions

#### 2. Prompt Engineering
```markdown
Good: "Create a REST API endpoint for user authentication using JWT"
Better: "Create a REST API endpoint for user authentication using JWT,
        following our existing pattern in handlers/, with tests"
```

#### 3. Iterative Development
- Start with high-level design
- Break down into smaller tasks
- Use Claude Code to implement each piece
- Review and refine

#### 4. Code Review
```bash
# Use Claude Code for reviews
"Review this code for security vulnerabilities and performance issues"
```

### Integrating with Other Tools

#### Zed (Recommended Code Editor)
[Zed](https://zed.dev) is a high-performance, collaborative code editor with built-in AI features:
```bash
# Install Zed with Homebrew (macOS)
brew install --cask zed

# Or download directly from https://zed.dev/download
# Available for macOS and Linux

# Launch Zed
open -a Zed  # macOS
zed .        # If added to PATH
```

**Why Zed:**
- ⚡ Native performance (written in Rust)
- 🤝 Real-time collaboration
- 🤖 Built-in AI assistant supporting multiple models
- 🎨 Beautiful, minimal UI
- 🔌 Growing extension ecosystem

**Configure Zed for AI:**
```json
// Settings → assistant.provider
{
  "assistant": {
    "default_model": "claude-3.5-sonnet",
    "provider": "anthropic"
  }
}
```

#### VS Code
```bash
# Install VS Code with Homebrew
brew install --cask visual-studio-code

# Install Claude Code extension (when available)
code --install-extension anthropic.claude-code
```

#### Cursor
Cursor has Claude integration built-in:
```bash
# Install Cursor with Homebrew
brew install --cask cursor

# Configure Claude in Cursor
# Settings → Models → Claude 3.5 Sonnet
```

#### Terminal Integration
```bash
# Use Claude Code CLI (when available)
claude "explain this error: $(cat error.log)"
claude "write a script to process these CSV files"
```

### Troubleshooting

#### Common Issues

**"Cannot connect to Claude"**
- Check internet connection
- Verify subscription is active
- Try logging out and back in

**"Context window exceeded"**
- Use `.clignore` to exclude unnecessary files
- Break large projects into smaller contexts
- Focus on specific directories

**"Slow response times"**
- Reduce context size
- Check for large files in project
- Ensure good internet connection

### Advanced Features

#### Custom Instructions
Add persistent instructions in Settings:
```
Always:
- Follow TDD practices
- Write comprehensive tests
- Use meaningful variable names
- Add error handling
- Document complex logic
```

#### Project Templates
Save and reuse project configurations:
```bash
# Export current setup
cp CLAUDE.md ~/.claude/templates/my-template.md

# Use in new projects
cp ~/.claude/templates/my-template.md ./CLAUDE.md
```

---

## Part 6: Knowledge Management with Obsidian

[Obsidian](https://obsidian.md) is a powerful knowledge base that works on local markdown files. It's perfect for managing documentation, notes, and knowledge graphs for your projects.

### Installing Obsidian

```bash
# macOS with Homebrew
brew install --cask obsidian

# Windows with Chocolatey
choco install obsidian

# Linux (AppImage)
wget https://github.com/obsidianmd/obsidian-releases/releases/latest/download/Obsidian-*.AppImage
chmod +x Obsidian-*.AppImage
./Obsidian-*.AppImage

# Or download from https://obsidian.md/download
```

### Why Obsidian for Development

- 📝 **Markdown-based**: All notes are plain text markdown files
- 🔗 **Linked Knowledge**: Create connections between concepts with [[wikilinks]]
- 📊 **Graph View**: Visualize relationships between documentation
- 🔍 **Powerful Search**: Find anything across all your notes instantly
- 🔌 **Extensible**: Hundreds of community plugins
- 💾 **Local First**: Your data stays on your machine

### Setting Up Obsidian for DDX

1. **Create a Vault for Your Project**
   ```bash
   # Create a docs vault in your project
   mkdir -p ~/Documents/ObsidianVaults/ProjectDocs

   # Open Obsidian and select "Open folder as vault"
   # Navigate to the created folder
   ```

2. **Install Essential Plugins**
   - **Templater**: For consistent document templates
   - **Dataview**: Query and visualize your notes
   - **Git**: Version control for your vault
   - **Kanban**: Project management boards
   - **Excalidraw**: Technical diagrams

3. **DDX Integration**
   ```bash
   # Use DDX to convert docs to Obsidian format
   ddx obsidian migrate docs/

   # This adds frontmatter and converts links to wikilinks
   ```

4. **Templates for Development**
   Create templates for common documentation:
   - Feature specifications
   - API documentation
   - Meeting notes
   - Bug reports
   - Architecture decisions

### Obsidian Workflow Tips

- **Daily Notes**: Track progress and discoveries
- **Tags**: Use #ddx #helix #architecture for organization
- **Canvas**: Visual project planning and architecture
- **Sync**: Use git for version control of your vault

---

## Part 7: Additional AI Development Tools

Beyond Claude Code, several other AI tools can enhance your development workflow:

### OpenAI Tools

#### GPT CLI (OpenAI Codex Access)
```bash
# Install openai CLI
pip install openai-cli

# Configure with your API key
export OPENAI_API_KEY="your-key-here"

# Use for code generation
openai api completions.create -m "gpt-4" -p "Write a Python function to..."
```

#### GitHub Copilot CLI
```bash
# Install GitHub Copilot CLI
gh extension install github/gh-copilot

# Use for command suggestions
gh copilot suggest "how to rebase my branch"
gh copilot explain "git reflog"
```

### Google Gemini CLI

```bash
# Install Gemini CLI
npm install -g @google/generative-ai-cli

# Configure with your API key
gemini config set api_key YOUR_API_KEY

# Use for code assistance
gemini "explain this error: $(cat error.log)"
gemini "optimize this SQL query: $(cat query.sql)"
```

### Open Source Alternatives

#### Ollama (Local AI Models)
```bash
# Install Ollama
brew install ollama

# Pull and run models locally
ollama pull codellama
ollama pull mistral
ollama pull phi

# Use for code generation
ollama run codellama "write a REST API in Go"
```

#### LM Studio
```bash
# Download from https://lmstudio.ai
# Run models locally with a GUI
# Supports Code Llama, Mistral, and many others
```

#### Continue.dev (VS Code/JetBrains)
```bash
# Install the Continue extension in VS Code
code --install-extension continue.continue

# Configure with local or cloud models
# Supports Ollama, OpenAI, Anthropic, and more
```

### AI Tool Comparison

| Tool | Best For | Cost | Privacy |
|------|----------|------|---------|
| Claude Code | Full project understanding | Subscription | Cloud |
| GitHub Copilot | Inline completions | Subscription | Cloud |
| Gemini | Google ecosystem integration | Pay-per-use | Cloud |
| Ollama | Privacy-focused development | Free | Local |
| LM Studio | GUI for local models | Free | Local |

### Setting Up Multiple AI Tools

Create an AI configuration file:
```bash
# ~/.ai-tools
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# Aliases for quick access
alias ai-claude="claude"
alias ai-gpt="openai api"
alias ai-gemini="gemini"
alias ai-local="ollama run codellama"
```

---

## Part 8: Development Language Runtimes

If you're not using containers, you'll need to install the runtime for your programming language. Here's how to set up the most common development environments using Homebrew and version managers.

### Python

Python is essential for many development tasks, data science, and scripting.

#### Install Python with Homebrew
```bash
# Install latest Python 3
brew install python@3.12

# Verify installation
python3 --version
pip3 --version
```

#### Using pyenv (Python Version Manager)
For projects requiring different Python versions:
```bash
# Install pyenv
brew install pyenv

# Add to shell configuration
echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.zshrc
echo '[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.zshrc
echo 'eval "$(pyenv init -)"' >> ~/.zshrc
source ~/.zshrc

# Install Python versions
pyenv install 3.12.0
pyenv install 3.11.7
pyenv install 3.10.13

# Set global default
pyenv global 3.12.0

# Set project-specific version
cd your-project
pyenv local 3.11.7  # Creates .python-version file
```

#### Virtual Environments
```bash
# Create virtual environment
python3 -m venv .venv

# Activate it
source .venv/bin/activate  # macOS/Linux
# or
.venv\Scripts\activate  # Windows

# Install packages
pip install -r requirements.txt

# Deactivate when done
deactivate
```

### Node.js

Node.js is required for JavaScript/TypeScript development and many modern tools.

#### Install Node.js with Homebrew
```bash
# Install latest LTS Node.js
brew install node

# Verify installation
node --version
npm --version
```

#### Using nvm (Node Version Manager)
For managing multiple Node versions:
```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Add to shell configuration (if not auto-added)
echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.zshrc
echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.zshrc
echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"' >> ~/.zshrc
source ~/.zshrc

# Install Node versions
nvm install --lts  # Latest LTS
nvm install 20     # Specific major version
nvm install 18.19.0  # Specific version

# Set default
nvm alias default 20

# Use specific version in project
cd your-project
nvm use 18
echo "18" > .nvmrc  # Save version for project
```

#### Package Managers
```bash
# npm comes with Node.js

# Install pnpm (faster, more efficient)
brew install pnpm

# Install yarn
brew install yarn

# Install bun (all-in-one toolkit)
curl -fsSL https://bun.sh/install | bash
```

### Go

Go is perfect for building fast, reliable, and efficient software.

#### Install Go with Homebrew
```bash
# Install latest Go
brew install go

# Verify installation
go version

# Set up Go workspace (optional, Go 1.13+ uses modules)
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```

#### Using g (Go Version Manager)
```bash
# Install g
curl -sSL https://git.io/g-install | sh -s

# Install Go versions
g install 1.21.5
g install 1.20.12

# List installed versions
g list

# Switch versions
g use 1.21.5

# Set for project
cd your-project
echo "1.21.5" > .go-version
```

#### Go Module Setup
```bash
# Initialize new module
cd your-project
go mod init github.com/yourusername/yourproject

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Build project
go build ./...

# Run tests
go test ./...
```

### Rust

Rust provides memory safety and blazing performance.

#### Install Rust with rustup
```bash
# Install rustup (Rust's official installer)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Add to PATH (if not auto-added)
echo 'source "$HOME/.cargo/env"' >> ~/.zshrc
source ~/.zshrc

# Verify installation
rustc --version
cargo --version
```

#### Managing Rust Toolchains
```bash
# Update Rust
rustup update

# Install specific version
rustup toolchain install 1.75.0
rustup toolchain install nightly

# Set default toolchain
rustup default stable

# Use specific toolchain for project
cd your-project
rustup override set 1.75.0

# Or use rust-toolchain.toml
cat > rust-toolchain.toml << 'EOF'
[toolchain]
channel = "1.75.0"
components = ["rustfmt", "clippy"]
EOF
```

#### Rust Development Tools
```bash
# Install essential tools
rustup component add rustfmt   # Code formatter
rustup component add clippy    # Linter
rustup component add rust-src   # Source code for std library

# Install cargo extensions
cargo install cargo-watch   # Auto-rebuild on changes
cargo install cargo-edit    # Add/remove dependencies
cargo install cargo-audit   # Security audit
cargo install sccache      # Compilation cache
```

### Language-Specific DDX Integration

#### Python Projects
```bash
# Initialize DDX with Python template
ddx init --template python-api

# Common Python DDX patterns
ddx patterns apply python-error-handling
ddx patterns apply python-testing
```

#### Node.js Projects
```bash
# Initialize DDX with Node.js template
ddx init --template nodejs-api
ddx init --template nextjs

# Install MCP servers (Node.js based)
ddx mcp install filesystem
```

#### Go Projects
```bash
# Initialize DDX with Go template
ddx init --template go-cli
ddx init --template go-api

# Apply Go patterns
ddx patterns apply go-error-handling
ddx patterns apply go-testing
```

#### Rust Projects
```bash
# Initialize DDX with Rust template
ddx init --template rust-cli
ddx init --template rust-wasm

# Apply Rust patterns
ddx patterns apply rust-error-handling
ddx patterns apply rust-async
```

### Version Management Best Practices

1. **Use Version Files**: Create `.python-version`, `.nvmrc`, `.go-version`, or `rust-toolchain.toml` in your projects
2. **Document Requirements**: Always specify language versions in README
3. **CI/CD Alignment**: Ensure CI uses same versions as development
4. **Team Consistency**: Share version files in git

### Quick Setup Script

Save this as `setup-languages.sh`:
```bash
#!/bin/bash

echo "Setting up development languages..."

# Python
brew install python@3.12 pyenv
pyenv install 3.12.0

# Node.js
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.zshrc
nvm install --lts

# Go
brew install go

# Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

echo "Done! Restart your terminal to use the new tools."
```

---

## Part 9: Installing DDX

Now let's install DDX itself!

### Prerequisites Check

Verify you have everything needed:
```bash
# Check Git
git --version
# Need: 2.0 or higher

# Check Homebrew (macOS/Linux)
brew --version
# Should show version

# Check terminal
echo $SHELL
# Should show /bin/zsh, /bin/bash, or similar
```

### Installation Methods

#### Method 1: Quick Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/ddx-tools/ddx/main/install.sh | bash
```

This script:
- ✅ Downloads the latest DDX binary
- ✅ Places it in ~/.local/bin
- ✅ Adds to your PATH
- ✅ Verifies installation

#### Method 2: Homebrew
```bash
# Add DDX tap
brew tap ddx-tools/ddx

# Install DDX
brew install ddx
```

#### Method 3: Go Install
```bash
# Requires Go 1.21+
go install github.com/ddx-tools/ddx/cli@latest
```

#### Method 4: From Source
```bash
# Clone repository
git clone https://github.com/ddx-tools/ddx
cd ddx/cli

# Build and install
make install
```

### Verification

Confirm DDX is installed correctly:
```bash
# Check version
ddx --version
# Output: ddx version 1.x.x

# Check help
ddx --help
# Should display available commands
```

✅ **Success**: You see version and help information
❌ **If "command not found"**: Restart terminal or run `source ~/.bashrc`

---

## Part 10: First Steps with DDX

### Initialize Your First Project

#### Create a New Project
```bash
# Create and enter project directory
mkdir my-awesome-project
cd my-awesome-project

# Initialize DDX
ddx init
```

You'll be prompted for:
1. Project template (or skip)
2. Configuration options
3. Resource selection

#### With a Template
```bash
# Initialize with Next.js template
ddx init --template nextjs

# Initialize with Python API template
ddx init --template python-api
```

### Essential DDX Commands

#### Explore Available Resources
```bash
# List everything available
ddx list

# List AI prompts
ddx prompts list

# List project templates
ddx templates list

# List code patterns
ddx patterns list
```

#### Apply Resources
```bash
# View a specific prompt
ddx prompts show code-review

# Apply a template
ddx templates apply react-component

# Add a pattern to your project
ddx patterns apply error-handling
```

#### Project Diagnosis
```bash
# Analyze your project setup
ddx diagnose

# Get improvement suggestions
ddx diagnose --verbose
```

### Integration with Claude Code

#### Setting Up Personas
```bash
# List available personas
ddx persona list

# Bind a persona to a role
ddx persona bind code-reviewer strict-code-reviewer

# Load personas into CLAUDE.md
ddx persona load
```

Your `.ddx.yml` will now contain:
```yaml
personas:
  bindings:
    code-reviewer: strict-code-reviewer
```

#### Installing MCP Servers

MCP (Model Context Protocol) servers extend Claude's capabilities:

```bash
# List available MCP servers
ddx mcp list

# Install filesystem MCP server
ddx mcp install filesystem

# Install GitHub MCP server
ddx mcp install github
```

This automatically configures `.claude/settings.json`:
```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "$PWD"]
    }
  }
}
```

---

## Part 11: Troubleshooting

### Common Issues and Solutions

#### Terminal Issues

**"Command not found" after installation:**
```bash
# Reload shell configuration
source ~/.bashrc  # or ~/.zshrc for zsh

# Check PATH
echo $PATH
# Should include ~/.local/bin or /usr/local/bin
```

**"Permission denied" errors:**
```bash
# Fix permissions for DDX
chmod +x ~/.local/bin/ddx

# Fix Homebrew permissions (macOS)
sudo chown -R $(whoami) /usr/local/bin /usr/local/lib
```

#### Homebrew Issues

**"brew: command not found":**
```bash
# Re-run post-install setup
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

**"Error: Your Xcode is too outdated" (macOS):**
```bash
# Install Xcode Command Line Tools
xcode-select --install
```

#### Git Issues

**"Please tell me who you are" error:**
```bash
git config --global user.name "Your Name"
git config --global user.email "your@email.com"
```

**GitHub authentication issues:**
```bash
# Re-authenticate with GitHub CLI
gh auth logout
gh auth login

# Check authentication status
gh auth status

# If behind a proxy, set environment variables
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
```

#### DDX Issues

**"Failed to initialize DDX":**
```bash
# Check git is initialized
git init

# Manually create config
cat > .ddx.yml << 'EOF'
version: "1.0"
repository: https://github.com/ddx-tools/ddx-master
branch: main
EOF
```

**"Resource not found":**
```bash
# Update DDX resources
ddx update

# Check available resources
ddx list --verbose
```

### Getting Help

#### DDX Help Resources
- **Documentation**: Run `ddx help <command>`
- **GitHub Issues**: [github.com/ddx-tools/ddx/issues](https://github.com/ddx-tools/ddx/issues)
- **Discord Community**: [discord.gg/ddx](https://discord.gg/ddx)

#### Claude Code Help
- **In-app Help**: Press `Cmd/Ctrl + ?` in Claude Code
- **Documentation**: [docs.anthropic.com/claude-code](https://docs.anthropic.com/claude-code)
- **Support**: [support.anthropic.com](https://support.anthropic.com)

---

## Part 12: Next Steps

### Recommended Learning Path

1. **Week 1: Master the Basics**
   - Practice terminal navigation
   - Learn basic git commands
   - Explore DDX templates

2. **Week 2: AI Integration**
   - Set up Claude Code workflows
   - Experiment with personas
   - Try MCP servers

3. **Week 3: Advanced Workflows**
   - Apply the HELIX workflow
   - Create custom templates
   - Share improvements

### Essential Resources

#### Terminal/Command Line
- [The Art of Command Line](https://github.com/jlevy/the-art-of-command-line)
- [Command Line Crash Course](https://developer.mozilla.org/en-US/docs/Learn/Tools_and_testing/Understanding_client-side_tools/Command_line)

#### Git
- [Pro Git Book](https://git-scm.com/book)
- [GitHub Skills](https://skills.github.com/)

#### DDX & Claude
- [DDX Documentation](https://docs.ddx.dev)
- [Claude Code Best Practices](https://docs.anthropic.com/claude-code/best-practices)

### Quick Reference Card

Save this for quick reference:

```bash
# Daily DDX Commands
ddx list                    # See what's available
ddx diagnose               # Check project health
ddx update                 # Get latest improvements
ddx prompts show <name>    # View a prompt
ddx templates apply <name> # Apply a template

# Git Essentials
git status                 # Check changes
git add .                  # Stage changes
git commit -m "message"    # Commit changes
git push                   # Push to remote
git pull                   # Pull from remote

# Terminal Navigation
pwd                        # Current directory
ls -la                     # List files
cd <directory>             # Change directory
mkdir <name>              # Create directory
rm <file>                 # Delete file
```

---

## Congratulations! 🎉

You now have a complete development environment with:
- ✅ Modern terminal emulator
- ✅ Homebrew package manager
- ✅ Git version control
- ✅ Claude Code AI assistant
- ✅ DDX toolkit

You're ready to build amazing things with AI-assisted development!

**Your first challenge**: Initialize a new project with `ddx init` and explore what DDX can do for you.

---

*Remember: Every expert was once a beginner. Take it one step at a time, and don't hesitate to ask for help in the DDX community!*