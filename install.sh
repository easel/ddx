#!/bin/bash

# DDx (Document-Driven Development eXperience) Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/yourusername/ddx/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DDX_HOME="${HOME}/.ddx"
DDX_REPO="https://github.com/easel/ddx"
DDX_BRANCH="main"

# Logging functions
log() {
    echo -e "${BLUE}[DDx]${NC} $1"
}

success() {
    echo -e "${GREEN}[DDx]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[DDx]${NC} $1"
}

error() {
    echo -e "${RED}[DDx]${NC} $1"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check for git
    if ! command -v git &> /dev/null; then
        error "Git is required but not installed. Please install git first."
    fi
    
    # Check for basic utilities (curl/wget for downloading binaries)
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        error "curl or wget is required but neither is installed."
    fi
    
    success "Prerequisites check passed"
}

# Clone DDx repository
clone_repository() {
    log "Installing DDx to ${DDX_HOME}..."
    
    if [ -d "${DDX_HOME}" ]; then
        warn "DDx already exists at ${DDX_HOME}. Updating..."
        cd "${DDX_HOME}"
        git pull origin "${DDX_BRANCH}"
    else
        git clone -b "${DDX_BRANCH}" "${DDX_REPO}" "${DDX_HOME}"
    fi
    
    success "Repository cloned successfully"
}

# Install CLI tool
install_cli() {
    log "Installing DDx CLI tool..."
    
    # Detect platform
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
    esac
    
    # Determine archive extension based on OS
    if [ "$OS" = "windows" ]; then
        ARCHIVE_EXT="zip"
        BINARY_NAME="ddx.exe"
    else
        ARCHIVE_EXT="tar.gz"
        BINARY_NAME="ddx"
    fi
    
    # Download appropriate archive
    ARCHIVE_NAME="ddx-${OS}-${ARCH}.${ARCHIVE_EXT}"
    DOWNLOAD_URL="${DDX_REPO}/releases/latest/download/${ARCHIVE_NAME}"
    
    log "Downloading ${ARCHIVE_NAME}..."
    
    # Create temp directory for download
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf ${TEMP_DIR}" EXIT
    
    if command -v curl &> /dev/null; then
        curl -fsSL "${DOWNLOAD_URL}" -o "${TEMP_DIR}/${ARCHIVE_NAME}"
    else
        wget -q "${DOWNLOAD_URL}" -O "${TEMP_DIR}/${ARCHIVE_NAME}"
    fi
    
    # Extract binary from archive
    log "Extracting binary..."
    cd "${TEMP_DIR}"
    
    if [ "$ARCHIVE_EXT" = "zip" ]; then
        unzip -q "${ARCHIVE_NAME}"
    else
        tar -xzf "${ARCHIVE_NAME}"
    fi
    
    # Move binary to DDx home
    mv "${BINARY_NAME}" "${DDX_HOME}/ddx"
    chmod +x "${DDX_HOME}/ddx"
    
    # Create symlink in user's local bin or add to PATH
    LOCAL_BIN="${HOME}/.local/bin"
    mkdir -p "${LOCAL_BIN}"
    ln -sf "${DDX_HOME}/ddx" "${LOCAL_BIN}/ddx"
    
    success "CLI tool installed"
}

# Set up shell completions
setup_completions() {
    log "Setting up shell completions..."
    
    # Detect shell
    SHELL_NAME=$(basename "$SHELL")
    
    case "$SHELL_NAME" in
        bash)
            COMPLETION_FILE="$HOME/.bash_completion"
            if [ -f "$COMPLETION_FILE" ]; then
                echo "# DDx completions" >> "$COMPLETION_FILE"
                echo "eval \"\$(ddx completion bash)\"" >> "$COMPLETION_FILE"
            fi
            ;;
        zsh)
            COMPLETION_DIR="$HOME/.zsh/completions"
            mkdir -p "$COMPLETION_DIR"
            ddx completion zsh > "$COMPLETION_DIR/_ddx" 2>/dev/null || true
            ;;
        fish)
            COMPLETION_DIR="$HOME/.config/fish/completions"
            mkdir -p "$COMPLETION_DIR"
            ddx completion fish > "$COMPLETION_DIR/ddx.fish" 2>/dev/null || true
            ;;
    esac
    
    success "Shell completions configured"
}

# Add to PATH if needed
update_path() {
    log "Checking PATH configuration..."
    
    # Local bin path
    LOCAL_BIN="${HOME}/.local/bin"
    
    # Check if already in PATH
    if [[ ":$PATH:" == *":$LOCAL_BIN:"* ]]; then
        success "PATH is already configured"
        return
    fi
    
    # Add to shell rc file
    SHELL_NAME=$(basename "$SHELL")
    case "$SHELL_NAME" in
        bash)
            RC_FILE="$HOME/.bashrc"
            ;;
        zsh)
            RC_FILE="$HOME/.zshrc"
            ;;
        fish)
            RC_FILE="$HOME/.config/fish/config.fish"
            ;;
        *)
            RC_FILE="$HOME/.profile"
            ;;
    esac
    
    if [ -f "$RC_FILE" ]; then
        echo "" >> "$RC_FILE"
        echo "# DDx CLI PATH" >> "$RC_FILE"
        echo "export PATH=\"\$PATH:$LOCAL_BIN\"" >> "$RC_FILE"
        success "Added DDx to PATH in $RC_FILE"
    else
        warn "Could not find shell config file. Please add $LOCAL_BIN to your PATH manually."
    fi
}

# Initial configuration
initial_config() {
    log "Running initial configuration..."
    
    # Create default config if it doesn't exist
    if [ ! -f "${HOME}/.ddx.yml" ]; then
        cp "${DDX_HOME}/.ddx.yml" "${HOME}/.ddx.yml"
        success "Created default configuration at ~/.ddx.yml"
    fi
}

# Show getting started information
show_getting_started() {
    echo ""
    echo "🎉 DDx (Document-Driven Development eXperience) installed successfully!"
    echo ""
    echo "📚 Getting Started:"
    echo "   ddx --help           Show available commands"
    echo "   ddx init             Initialize DDx in a project"
    echo "   ddx list             Show available templates and patterns"
    echo "   ddx diagnose         Analyze current project setup"
    echo ""
    echo "📖 Documentation:"
    echo "   ${DDX_HOME}/docs/    Local documentation"
    echo "   ${DDX_REPO}          Online repository"
    echo ""
    echo "🔧 Configuration:"
    echo "   ~/.ddx.yml           Global configuration file"
    echo "   ${DDX_HOME}          DDx installation directory"
    echo ""
    echo "⚡ Quick Start:"
    echo "   cd your-project"
    echo "   ddx init"
    echo ""
    
    if command -v ddx &> /dev/null; then
        success "DDx is ready to use! Run 'ddx --version' to verify."
    else
        warn "Please restart your shell or run 'source ~/.${SHELL_NAME}rc' to use ddx command."
    fi
}

# Main installation flow
main() {
    echo "🚀 Installing DDx - Document-Driven Development eXperience"
    echo ""
    
    check_prerequisites
    clone_repository
    install_cli
    setup_completions
    update_path
    initial_config
    show_getting_started
}

# Run installation
main "$@"