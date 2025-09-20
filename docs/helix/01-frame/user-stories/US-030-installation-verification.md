# User Story: US-030 - Installation Verification

**Feature**: FEAT-004 - Cross-Platform Installation
**Priority**: P0
**Status**: Draft

## Story

**As a** developer who just installed DDX
**I want to** verify DDX installed correctly
**So that** I know it's ready to use and can troubleshoot any issues

## Business Value

- Reduces support tickets from incomplete installations
- Increases user confidence in the installation process
- Enables self-service troubleshooting

## Acceptance Criteria

### AC-001: Automatic Post-Install Verification
**Given** DDX installation completes
**When** the installer finishes
**Then** it automatically runs verification checks and reports status

### AC-002: Version Command Verification
**Given** DDX is installed
**When** I run `ddx version`
**Then** it displays the correct installed version number

### AC-003: Health Check Command
**Given** DDX is installed
**When** I run `ddx doctor`
**Then** it checks and reports on installation health including PATH configuration

### AC-004: Dependency Verification
**Given** DDX is installed
**When** verification runs
**Then** it confirms git is available and reports any missing dependencies

### AC-005: Clear Issue Reporting
**Given** verification detects issues
**When** checks complete
**Then** specific problems and suggested fixes are displayed

## Definition of Done

- [ ] Automatic verification runs after installation
- [ ] `ddx version` command works reliably
- [ ] `ddx doctor` provides comprehensive health check
- [ ] Git availability is verified
- [ ] Clear error messages with actionable suggestions
- [ ] Verification completes within 5 seconds

## Notes

- Should check both installation and runtime environment
- Must provide specific remediation steps for common issues
- Verification should not require network access