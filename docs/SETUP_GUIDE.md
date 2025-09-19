# Complete Setup Guide for DDX and AI-Assisted Development

> From zero to productive: Setting up your development environment for DDX and modern AI-assisted coding

## Prerequisites Overview

Before we begin, here's what we'll be setting up:
- ✅ **Terminal emulator** - Your command-line interface
- ✅ **Homebrew** - Package manager for installing tools
- ✅ **Git** - Version control system
- ✅ **Claude Code** - AI-powered development assistant
- ✅ **DDX** - Document-Driven Development toolkit

**Time needed**: 30-45 minutes for complete setup

---

## Part 1: Terminal Setup

The terminal (command line) is where you'll interact with DDX and other development tools. Let's get you set up with a modern terminal.

### macOS

macOS comes with Terminal.app built-in, but you have better options:

#### Option 1: Built-in Terminal (Quick Start)
1. Press `Cmd + Space` and type "Terminal"
2. Press Enter to open Terminal.app
3. You're ready! (Though consider upgrading to Option 2 or 3)

#### Option 2: iTerm2 (Recommended)
```bash
# Download from: https://iterm2.com
# Or install via Homebrew (after Part 2):
brew install --cask iterm2
```

#### Option 3: Warp (Modern Alternative)
```bash
# Download from: https://www.warp.dev
# Or install via Homebrew (after Part 2):
brew install --cask warp
```

### Windows

Windows users should use Windows Terminal with WSL2 for the best experience.

#### Step 1: Install Windows Terminal
1. Open Microsoft Store
2. Search for "Windows Terminal"
3. Click Install
4. Launch Windows Terminal from Start Menu

#### Step 2: Install WSL2 (Windows Subsystem for Linux)
Open PowerShell as Administrator and run:
```powershell
# Enable WSL
wsl --install

# Restart your computer when prompted
# After restart, WSL will complete Ubuntu installation
# Create a username and password when prompted
```

#### Step 3: Configure Windows Terminal
1. Open Windows Terminal
2. Click the dropdown arrow → Settings
3. Set "Ubuntu" as default profile
4. Save settings

**Important**: Use the Ubuntu/WSL2 terminal for all following commands.

### Linux

Linux distributions come with capable terminals:

#### Ubuntu/Debian
- Default: GNOME Terminal (already installed)
- Alternative: `sudo apt install terminator`

#### Fedora/RHEL
- Default: GNOME Terminal (already installed)
- Alternative: `sudo dnf install tilix`

#### Arch
- Install: `sudo pacman -S alacritty`

---

## Part 2: Installing Homebrew

Homebrew is a package manager that makes installing development tools simple. Think of it as an app store for command-line tools.

### macOS and Linux

Open your terminal and run this single command:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

**What this does:**
- Downloads the Homebrew installer
- Sets up Homebrew in the correct location
- Configures your system to use Homebrew

#### Post-Installation Setup

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

#### Verify Installation
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

### Setting Up SSH for GitHub

SSH keys allow secure communication with GitHub without passwords.

#### Generate SSH Key
```bash
# Generate a new SSH key (replace with your email)
ssh-keygen -t ed25519 -C "your.email@example.com"

# Press Enter to accept default location
# Optionally set a passphrase (or press Enter for none)
```

#### Add SSH Key to SSH Agent
```bash
# Start the ssh-agent
eval "$(ssh-agent -s)"

# Add your SSH key
ssh-add ~/.ssh/id_ed25519
```

#### Add SSH Key to GitHub
```bash
# Copy your public key
cat ~/.ssh/id_ed25519.pub
# This displays your key - copy the entire output
```

Then:
1. Go to [GitHub SSH Settings](https://github.com/settings/keys)
2. Click "New SSH key"
3. Paste your key and save

#### Test Connection
```bash
ssh -T git@github.com
# Should see: "Hi username! You've successfully authenticated..."
```

---

## Part 4: Installing Claude Code

Claude Code is an AI-powered development assistant that works seamlessly with DDX.

### What is Claude Code?

Claude Code is a desktop application that provides:
- AI-powered code generation and review
- Intelligent refactoring suggestions
- Natural language to code translation
- Integration with your development workflow

### Installation Process

#### Step 1: Download Claude Code

Visit [claude.ai/download](https://claude.ai/download) and download for your platform:
- **macOS**: Claude-Code-x.x.x.dmg
- **Windows**: Claude-Code-Setup-x.x.x.exe
- **Linux**: Claude-Code-x.x.x.AppImage

#### Step 2: Install

**macOS:**
1. Open the .dmg file
2. Drag Claude Code to Applications folder
3. Launch from Applications

**Windows:**
1. Run the installer
2. Follow the setup wizard
3. Launch from Start Menu

**Linux:**
```bash
# Make AppImage executable
chmod +x Claude-Code-*.AppImage

# Run it
./Claude-Code-*.AppImage

# Optional: Move to Applications
sudo mv Claude-Code-*.AppImage /opt/claude-code
sudo ln -s /opt/claude-code /usr/local/bin/claude-code
```

#### Step 3: Initial Setup

1. Launch Claude Code
2. Sign in with your Anthropic account
3. Complete the welcome tour
4. Configure your preferences:
   - Theme (Light/Dark)
   - Font size
   - Key bindings

### Creating CLAUDE.md Files

CLAUDE.md files provide context to Claude about your project:

```bash
# Create a CLAUDE.md in your project root
cat > CLAUDE.md << 'EOF'
# Project Context for Claude

## Project Overview
This project uses DDX for development workflow management.

## Technology Stack
- Language: [Your language]
- Framework: [Your framework]
- Database: [Your database]

## Development Guidelines
- Follow test-driven development
- Use descriptive variable names
- Write comprehensive documentation

## DDX Integration
This project uses DDX personas and workflows.
Check .ddx.yml for configuration.
EOF
```

### Integrating Claude Code with Your Editor

#### VS Code Integration
```bash
# Install Claude Code extension
code --install-extension anthropic.claude-code
```

#### Cursor Integration
Cursor has Claude built-in:
```bash
# Install Cursor
brew install --cask cursor

# Or download from: https://cursor.sh
```

---

## Part 5: Installing DDX

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

## Part 6: First Steps with DDX

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

## Part 7: Troubleshooting

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

**SSH key not working:**
```bash
# Check if ssh-agent is running
eval "$(ssh-agent -s)"

# Re-add your key
ssh-add ~/.ssh/id_ed25519

# Test GitHub connection
ssh -T git@github.com
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

## Part 8: Next Steps

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