---
title: "Implementation Plan - US-028 One-Command Installation"
type: implementation-plan
user_story_id: US-028
feature_id: FEAT-004
workflow_phase: build
artifact_type: implementation-plan
tags:
  - helix/build
  - helix/artifact/implementation
  - installation
  - one-command
related:
  - "[[US-028-one-command-installation]]"
  - "[[SD-004-cross-platform-installation]]"
  - "[[TS-004-installation-test-specification]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Implementation Plan: US-028 One-Command Installation

## User Story Reference

**US-028**: As a developer new to DDX, I want to install DDX with a single command, so that I can start using it immediately without complex setup procedures.

**Acceptance Criteria**:
- AC-001: Unix Installation: `curl -sSL https://ddx.dev/install | sh`
- AC-002: Windows Installation: `iwr -useb https://ddx.dev/install.ps1 | iex`
- AC-003: Platform Auto-Detection
- AC-004: No Admin Privileges Required
- AC-005: Installation Verification

## Component Mapping

**Primary Component**: InstallationManager (from SD-004)
**Supporting Components**:
- PlatformDetector (platform detection)
- BinaryDistributor (GitHub releases)
- EnvironmentConfigurer (PATH setup)
- InstallationValidator (verification)

## Implementation Strategy

### Overview
Enhance the existing `install.sh` script and create a new `install.ps1` PowerShell script to provide one-command installation across all platforms with GitHub releases integration.

### Key Changes
1. Update `install.sh` to use GitHub releases instead of git clone
2. Create `install.ps1` for Windows PowerShell
3. Add platform detection and binary selection
4. Implement checksum verification
5. Add progress indicators and error handling

## Detailed Implementation Steps

### Step 1: Update Unix/Linux/macOS Installation Script

**File**: `/install.sh`

**Changes Required**:
1. **Replace git clone with binary download**:
   ```bash
   # OLD: Clone repository approach
   git clone -b "${DDX_BRANCH}" "${DDX_REPO}" "${DDX_HOME}"

   # NEW: Download binary from GitHub releases
   download_binary_from_github_releases()
   ```

2. **Add platform/architecture detection**:
   ```bash
   detect_platform() {
       OS=$(uname -s | tr '[:upper:]' '[:lower:]')
       ARCH=$(uname -m)

       case "$ARCH" in
           x86_64) ARCH="amd64" ;;
           aarch64|arm64) ARCH="arm64" ;;
           armv7l) ARCH="arm" ;;
           *) error "Unsupported architecture: $ARCH" ;;
       esac

       case "$OS" in
           linux|darwin) ;;
           *) error "Unsupported operating system: $OS" ;;
       esac
   }
   ```

3. **Implement GitHub API integration**:
   ```bash
   get_latest_release() {
       local api_url="https://api.github.com/repos/easel/ddx/releases/latest"
       if command -v curl &> /dev/null; then
           curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
       elif command -v wget &> /dev/null; then
           wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
       else
           error "Neither curl nor wget found"
       fi
   }
   ```

4. **Add binary download and verification**:
   ```bash
   download_and_verify_binary() {
       local version="$1"
       local binary_name="ddx-${OS}-${ARCH}"
       local download_url="https://github.com/easel/ddx/releases/download/${version}/${binary_name}.tar.gz"
       local checksum_url="https://github.com/easel/ddx/releases/download/${version}/${binary_name}.tar.gz.sha256"

       log "Downloading DDX ${version} for ${OS}/${ARCH}..."

       # Download binary and checksum
       download_file "$download_url" "${TEMP_DIR}/${binary_name}.tar.gz"
       download_file "$checksum_url" "${TEMP_DIR}/${binary_name}.tar.gz.sha256"

       # Verify checksum
       verify_checksum "${TEMP_DIR}/${binary_name}.tar.gz" "${TEMP_DIR}/${binary_name}.tar.gz.sha256"

       # Extract binary
       tar -xzf "${TEMP_DIR}/${binary_name}.tar.gz" -C "$TEMP_DIR"

       # Install binary
       install_binary "${TEMP_DIR}/ddx"
   }
   ```

5. **Add progress indicators**:
   ```bash
   show_progress() {
       local current="$1"
       local total="$2"
       local message="$3"
       local percent=$((current * 100 / total))
       printf "\r[%3d%%] %s" "$percent" "$message"
   }
   ```

### Step 2: Create Windows PowerShell Installation Script

**File**: `/install.ps1`

**Implementation**:
```powershell
# DDX (Document-Driven Development eXperience) Windows Installation Script
# Usage: iwr -useb https://ddx.dev/install.ps1 | iex

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:USERPROFILE\bin",
    [switch]$Force = $false
)

# Configuration
$ErrorActionPreference = "Stop"
$DDX_REPO = "easel/ddx"
$DDX_API_BASE = "https://api.github.com/repos/$DDX_REPO"

# Logging functions
function Write-Log {
    param([string]$Message)
    Write-Host "[DDX] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[DDX] $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "[DDX] $Message" -ForegroundColor Red
    exit 1
}

# Platform detection
function Get-Platform {
    $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    return @{
        OS = "windows"
        Arch = $arch
        Extension = "zip"
        Binary = "ddx.exe"
    }
}

# GitHub API functions
function Get-LatestRelease {
    try {
        $response = Invoke-RestMethod -Uri "$DDX_API_BASE/releases/latest" -UseBasicParsing
        return $response.tag_name
    }
    catch {
        Write-Error "Failed to get latest release: $_"
    }
}

# Download and installation
function Install-DDX {
    Write-Log "Installing DDX for Windows..."

    $platform = Get-Platform
    $version = if ($Version -eq "latest") { Get-LatestRelease } else { $Version }

    Write-Log "Installing DDX version $version..."

    # Create installation directory
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    # Download binary
    $binaryName = "ddx-$($platform.OS)-$($platform.Arch)"
    $downloadUrl = "https://github.com/$DDX_REPO/releases/download/$version/$binaryName.$($platform.Extension)"
    $tempFile = Join-Path $env:TEMP "$binaryName.$($platform.Extension)"

    Write-Log "Downloading from $downloadUrl..."
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
    }
    catch {
        Write-Error "Failed to download DDX: $_"
    }

    # Extract and install
    Write-Log "Extracting binary..."
    Expand-Archive -Path $tempFile -DestinationPath $env:TEMP -Force

    $ddxBinary = Join-Path $InstallDir $platform.Binary
    Move-Item -Path (Join-Path $env:TEMP $platform.Binary) -Destination $ddxBinary -Force

    # Add to PATH
    Add-ToPath $InstallDir

    # Verify installation
    try {
        & $ddxBinary version | Out-Null
        Write-Success "DDX installed successfully!"
        Write-Success "Run 'ddx --help' to get started."
    }
    catch {
        Write-Error "Installation verification failed: $_"
    }

    # Cleanup
    Remove-Item $tempFile -ErrorAction SilentlyContinue
}

# PATH management
function Add-ToPath {
    param([string]$Directory)

    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$Directory*") {
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Log "Added $Directory to user PATH"
        Write-Log "Please restart your terminal or run 'refreshenv' to use DDX"
    }
}

# Main execution
try {
    Install-DDX
}
catch {
    Write-Error "Installation failed: $_"
}
```

### Step 3: Create GitHub Actions for Release Automation

**File**: `/.github/workflows/release.yml`

**Implementation**:
```yaml
name: Release DDX

on:
  release:
    types: [published]

jobs:
  build-and-release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
        exclude:
          - goos: windows
            goarch: arm64

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          cd cli
          binary_name="ddx"
          if [ "$GOOS" = "windows" ]; then
            binary_name="ddx.exe"
          fi

          go build -o "$binary_name" .

          # Create archive
          archive_name="ddx-$GOOS-$GOARCH"
          if [ "$GOOS" = "windows" ]; then
            zip "$archive_name.zip" "$binary_name"
            echo "ASSET_PATH=$archive_name.zip" >> $GITHUB_ENV
          else
            tar -czf "$archive_name.tar.gz" "$binary_name"
            echo "ASSET_PATH=$archive_name.tar.gz" >> $GITHUB_ENV
          fi

      - name: Generate checksum
        run: |
          if [ "${{ matrix.goos }}" = "windows" ]; then
            sha256sum "$ASSET_PATH" > "$ASSET_PATH.sha256"
          else
            shasum -a 256 "$ASSET_PATH" > "$ASSET_PATH.sha256"
          fi

      - name: Upload Release Asset
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ github.event.release.upload_url }}
          asset_path: ./cli/${{ env.ASSET_PATH }}
          asset_name: ${{ env.ASSET_PATH }}
          asset_content_type: application/octet-stream

      - name: Upload Checksum
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ github.event.release.upload_url }}
          asset_path: ./cli/${{ env.ASSET_PATH }}.sha256
          asset_name: ${{ env.ASSET_PATH }}.sha256
          asset_content_type: text/plain
```

## Integration Points

### Test Coverage
- **Primary Test**: `TestAcceptance_US028_OneCommandInstallation`
- **Test Scenarios**:
  - Unix one-command install
  - macOS one-command install
  - Windows one-command install
- **Test Location**: `cli/cmd/installation_acceptance_test.go:140-180`

### Related Components
- **US-029**: Installation scripts will set up PATH (EnvironmentConfigurer)
- **US-030**: Installation will be verified by doctor command (InstallationValidator)
- **US-035**: Installation failures will be diagnosed (InstallationManager)

### File Dependencies
- Requires GitHub releases to be set up and populated
- Needs web server to host install scripts at ddx.dev domain
- Must coordinate with existing `install.sh` (update, don't replace)

## Implementation Sequence

### Phase 1: Update Existing Script
1. Modify `install.sh` to use GitHub releases
2. Add platform detection logic
3. Implement binary download and verification
4. Add progress indicators

### Phase 2: Create Windows Script
1. Implement `install.ps1` with equivalent functionality
2. Add PowerShell-specific PATH management
3. Test on Windows environments

### Phase 3: Release Automation
1. Create GitHub Actions workflow
2. Test release process
3. Generate initial releases

### Phase 4: Integration Testing
1. Test installation scripts manually
2. Run acceptance tests
3. Verify cross-platform compatibility

## Success Criteria

✅ **Implementation Complete When**:
1. `TestAcceptance_US028_OneCommandInstallation` passes on all platforms
2. Unix command `curl -sSL https://ddx.dev/install | sh` works
3. Windows command `iwr -useb https://ddx.dev/install.ps1 | iex` works
4. Installation completes in under 60 seconds
5. No admin privileges required for installation
6. Binary is executable and `ddx version` works after installation

## Risk Mitigation

### High-Risk Areas
1. **GitHub API Rate Limits**: Implement caching and fallback URLs
2. **Network Failures**: Add retry logic with exponential backoff
3. **Platform Detection Failures**: Provide manual override options
4. **Checksum Verification**: Ensure secure download verification

### Testing Strategy
1. **Local Testing**: Test scripts in isolated environments
2. **CI/CD Testing**: Automated testing in GitHub Actions
3. **Manual Verification**: Test on real Windows/macOS/Linux systems
4. **Edge Case Testing**: Test with slow networks, interruptions, permissions issues

---

This implementation plan provides step-by-step guidance for implementing US-028 while maintaining proper test-driven development practices and ensuring integration with the broader DDX installation system.