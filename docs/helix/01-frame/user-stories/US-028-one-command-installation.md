# User Story: US-028 - One-Command Installation

**Feature**: FEAT-004 - Cross-Platform Installation
**Priority**: P0
**Status**: Draft

## Story

**As a** developer new to DDX
**I want to** install DDX with a single command
**So that** I can start using it immediately without complex setup procedures

## Business Value

- Reduces barrier to adoption by eliminating installation friction
- Minimizes support overhead from installation issues
- Increases conversion rate from trial to active usage

## Acceptance Criteria

### AC-001: Unix Installation
**Given** I am on a Unix-like system (macOS/Linux)
**When** I execute `curl -sSL https://ddx.dev/install | sh`
**Then** DDX is installed and ready to use without additional setup

### AC-002: Windows Installation
**Given** I am on Windows with PowerShell
**When** I execute `iwr -useb https://ddx.dev/install.ps1 | iex`
**Then** DDX is installed and ready to use without additional setup

### AC-003: Platform Auto-Detection
**Given** I run the installation command
**When** the installer executes
**Then** it automatically detects my OS and architecture without prompting

### AC-004: No Admin Privileges Required
**Given** I am a standard user without admin rights
**When** I run the installation command
**Then** installation completes successfully in my user directory

### AC-005: Installation Verification
**Given** installation completes
**When** I run `ddx version`
**Then** it displays the installed version number

### AC-006: Binary Download from GitHub Releases
**Given** the installer runs
**When** downloading DDX binary from GitHub releases
**Then** it downloads the correct platform-specific binary (e.g., ddx-linux-amd64.tar.gz, ddx-darwin-arm64.tar.gz)
**And** it verifies SHA256 checksum matches release assets

### AC-007: Correct Binary Selection
**Given** platform detection runs
**When** selecting binary from GitHub releases
**Then** the correct OS/architecture combination is automatically selected
**And** the binary is downloaded from https://github.com/easel/ddx/releases/latest/download/{binary-name}

## Definition of Done

- [ ] Installation works on macOS, Linux, and Windows
- [ ] No sudo/admin privileges required
- [ ] Installation completes in under 60 seconds
- [ ] `ddx version` command works after installation
- [ ] Installation process shows progress indicators
- [ ] Error messages are clear and actionable
- [ ] Binary is downloaded from GitHub releases (not repository cloning)
- [ ] SHA256 checksum verification prevents corrupted downloads
- [ ] Platform/architecture detection works for all supported combinations

## Notes

- Installation success rate target: >99%
- Must work on corporate networks with proxy settings
- Should handle interrupted downloads gracefully