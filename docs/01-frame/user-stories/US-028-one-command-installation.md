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

## Definition of Done

- [ ] Installation works on macOS, Linux, and Windows
- [ ] No sudo/admin privileges required
- [ ] Installation completes in under 60 seconds
- [ ] `ddx version` command works after installation
- [ ] Installation process shows progress indicators
- [ ] Error messages are clear and actionable

## Notes

- Installation success rate target: >99%
- Must work on corporate networks with proxy settings
- Should handle interrupted downloads gracefully