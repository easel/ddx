# User Story: US-008 - Check DDX Version

**Story ID**: US-008
**Feature**: FEAT-001 - Core CLI Framework
**Priority**: P1
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-09-29

## Story

**As a** developer
**I want** to check the DDX version
**So that** I know if updates are available

## Acceptance Criteria

- [ ] **Given** I want version info, **when** I run `ddx version`, **then** the current DDX version number is displayed
- [ ] **Given** version is displayed, **when** I view the output, **then** build information (commit hash, build date) is included
- [ ] **Given** I'm online, **when** I check version, **then** the system checks for available updates automatically
- [ ] **Given** updates are available, **when** version is displayed, **then** changelog highlights for newer versions are shown
- [ ] **Given** I don't want update checks, **when** I run `ddx version --no-check`, **then** update checking is suppressed
- [ ] **Given** my version is outdated, **when** I check version, **then** a clear indication that the version is outdated is shown
- [ ] **Given** version changes may affect compatibility, **when** updates are available, **then** compatibility warnings are displayed

## Definition of Done

- [ ] Version command implemented
- [ ] Build information embedded in binary
- [ ] Update checking mechanism implemented
- [ ] Changelog fetching and display
- [ ] Flag to suppress update checks
- [ ] Outdated version warnings
- [ ] Compatibility warning system
- [ ] Unit tests written and passing
- [ ] Integration tests for version checking
- [ ] Documentation updated

## Technical Notes

### Implementation Considerations
- Version should be embedded at build time
- Use semantic versioning (major.minor.patch)
- Check for updates against GitHub releases API
- Cache update check results (check once per day) - **Note**: Automatic checking is implemented in US-043
- Include git commit hash in build info
- Support pre-release versions (alpha, beta, rc)

### Error Scenarios
- Network unavailable for update check
- GitHub API rate limited
- Invalid version format in binary
- Update check fails but shouldn't block version display

## Validation Scenarios

### Scenario 1: Basic Version Check
1. Run `ddx version`
2. **Expected**: Display version like "DDX v1.2.3 (commit: abc123, built: 2025-01-14)"

### Scenario 2: Update Available
1. Run `ddx version` when newer version exists
2. **Expected**: Current version shown plus "Update available: v1.2.4"

### Scenario 3: Suppress Update Check
1. Run `ddx version --no-check`
2. **Expected**: Version displayed without checking for updates

### Scenario 4: Offline Version Check
1. Disconnect from network
2. Run `ddx version`
3. **Expected**: Version displayed, update check skipped gracefully

### Scenario 5: Pre-release Version
1. Run `ddx version` on beta build
2. **Expected**: Display "DDX v1.2.3-beta.1" with appropriate warnings

## User Persona

### Primary: Maintenance-Conscious Developer
- **Role**: Developer who keeps tools updated
- **Goals**: Stay current with latest features and fixes
- **Pain Points**: Not knowing when updates are available
- **Technical Level**: All levels

### Secondary: Support Seeker
- **Role**: Developer troubleshooting issues
- **Goals**: Verify version for bug reports
- **Pain Points**: Unclear version information
- **Technical Level**: Beginner to intermediate

## Dependencies

- Build system must embed version info
- Network access for update checks (optional)
- GitHub API for release information

## Related Stories

- US-004: Update Assets from Master (natural next step after finding update)
- US-006: Get Command Help (often checked together)
- US-043: Automatic Update Notifications (implements automatic update checking mentioned in AC)

---
*This user story is part of FEAT-001: Core CLI Framework*