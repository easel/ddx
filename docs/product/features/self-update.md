---
tags: [feature-spec, self-update, cli, maintenance, development-workflow]
version: 1.0.0
---

# Feature Specification: DDx Self-Update

**PRD Reference**: [DDx PRD v1 - Cross-Platform Installation](../prd-ddx-v1.md#feature-cross-platform-installation)  
**Epic/Initiative**: CLI Infrastructure and Maintenance  
**Status**: Draft  
**Created**: 2025-01-12  
**Last Updated**: 2025-01-12  
**Owner**: DDx Team  
**Tech Lead**: TBD

## Overview

### Feature Description
The DDx self-update feature enables users to update their DDx binary to the latest version directly from the command line, without needing to re-run the installation script. This feature checks for new releases on GitHub, downloads the appropriate platform-specific binary, and safely replaces the current executable.

### Business Context
As DDx evolves rapidly with new features, bug fixes, and pattern improvements, users need a frictionless way to stay current. Manual updates through re-running install scripts create friction and reduce adoption of improvements. A built-in self-update mechanism ensures users benefit from the latest enhancements and security fixes with minimal effort.

### Success Criteria
- **Update Adoption Rate**: >70% of users on latest version within 30 days of release
- **Update Success Rate**: >95% successful updates without manual intervention
- **Time to Update**: <30 seconds for binary update on standard connection
- **User Satisfaction**: Zero friction update experience

## User Stories

### Primary User Story
**As a** DDx user  
**I want** to update my DDx binary with a simple command  
**So that** I can access the latest features and fixes without manual installation

**Acceptance Criteria:**
- [ ] Running `ddx self-update` checks for and installs the latest version
- [ ] Update process preserves my configuration and settings
- [ ] Clear feedback shows update progress and completion
- [ ] Update completes within 30 seconds on standard connection
- [ ] Rollback is possible if update fails

### Secondary User Stories

#### Check for Updates
**As a** cautious user  
**I want** to check if updates are available without installing  
**So that** I can review changes before updating

**Acceptance Criteria:**
- [ ] `ddx self-update --check` shows available version without installing
- [ ] Changelog or release notes are displayed
- [ ] Current vs available version comparison is shown
- [ ] No changes are made to the system

#### Automated Update Checks
**As a** regular DDx user  
**I want** to be notified when updates are available  
**So that** I don't miss important improvements

**Acceptance Criteria:**
- [ ] `ddx version` shows if newer version is available
- [ ] Update check is performed periodically (configurable)
- [ ] Notifications are non-intrusive
- [ ] Can disable automatic checks via configuration

#### CI/CD Integration
**As a** DevOps engineer  
**I want** to ensure DDx is up-to-date in CI pipelines  
**So that** builds use the latest stable version

**Acceptance Criteria:**
- [ ] `ddx self-update --ci` updates without interactive prompts
- [ ] Exit codes indicate update status
- [ ] Can specify version constraints
- [ ] Silent mode available for automated environments

## Functional Requirements

### Core Functionality

#### Version Management
**Purpose**: Track and compare DDx versions  
**Behavior**: Compare semantic versions and determine update availability  
**Inputs**: Current version (embedded), Latest version (from GitHub API)  
**Outputs**: Update availability status, version delta  
**Business Rules**: 
- Follow semantic versioning (major.minor.patch)
- Pre-release versions handled appropriately
- Version downgrade prevented by default
- Force flag allows specific version installation

#### Update Discovery
**Purpose**: Check GitHub releases for new versions  
**Behavior**: Query GitHub Releases API for latest stable release  
**Inputs**: Current version, GitHub repository URL  
**Outputs**: Latest version info, download URLs, changelog  
**Business Rules**:
- Only stable releases considered by default
- Pre-release opt-in via flag
- Rate limiting handled gracefully
- Offline mode degrades gracefully

#### Binary Download
**Purpose**: Retrieve platform-specific binary  
**Behavior**: Download correct binary for user's platform  
**Inputs**: Platform (OS/arch), Version to download  
**Outputs**: Downloaded binary file  
**Business Rules**:
- Platform detection automatic
- Architecture detection (amd64, arm64)
- Resume partial downloads
- Verify checksums before installation

#### Binary Replacement
**Purpose**: Safely replace running executable  
**Behavior**: Atomic replacement with rollback capability  
**Inputs**: New binary, Current binary location  
**Outputs**: Updated binary in place  
**Business Rules**:
- Backup current binary before replacement
- Preserve file permissions and ownership
- Handle running executable replacement
- Rollback on failure

### User Interactions

#### Workflow 1: Interactive Update
**Trigger**: User runs `ddx self-update`  
**Steps**:
1. Check current version
2. Query GitHub for latest release
3. Display version comparison and changelog
4. Prompt for confirmation
5. Download new binary
6. Verify checksum
7. Backup current binary
8. Replace binary
9. Verify installation
10. Display success message

**Alternative Flows**:
- **Already up-to-date**: Display message and exit
- **Download failure**: Retry with exponential backoff
- **Verification failure**: Abort and restore backup
- **Network unavailable**: Display offline message

#### Workflow 2: Check-Only Mode
**Trigger**: User runs `ddx self-update --check`  
**Steps**:
1. Check current version
2. Query GitHub for latest release
3. Compare versions
4. Display results without making changes
5. Show upgrade command if update available

**Alternative Flows**:
- **Rate limited**: Use cached check result if recent
- **Network error**: Display last known check result

## Technical Requirements

### Architecture Components

#### GitHub API Client
- RESTful API integration with github.com
- Rate limiting awareness (60/hour unauthenticated)
- Optional authentication via GITHUB_TOKEN
- Caching of API responses

#### Version Comparator
- Semantic version parsing and comparison
- Pre-release version handling
- Version constraint evaluation

#### Binary Manager
- Platform detection (runtime.GOOS, runtime.GOARCH)
- Download with progress indication
- Checksum verification (SHA256)
- Atomic file replacement
- Permission preservation

#### Configuration Integration
- Update check frequency setting
- Pre-release opt-in setting
- Proxy configuration support
- Custom mirror support (future)

### Dependencies
- No external Go dependencies beyond standard library
- GitHub API v3 (REST)
- HTTPS for secure downloads

### Platform Support
- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)
- Future: FreeBSD, ARM variants

## API Specifications

### CLI Commands

#### `ddx self-update`
Update DDx to the latest version

**Flags:**
- `--check`: Check for updates without installing
- `--force`: Force update even if up-to-date
- `--version <version>`: Install specific version
- `--pre-release`: Include pre-release versions
- `--yes`: Skip confirmation prompt
- `--no-backup`: Skip backup creation

**Exit Codes:**
- 0: Success (updated or already current)
- 1: Update check failed
- 2: Download failed
- 3: Installation failed
- 4: Version constraint not met

#### `ddx version`
Show version with update check

**Flags:**
- `--check-update`: Check for available updates
- `--json`: Output in JSON format

**Output Format:**
```
DDx v1.2.3
Commit: abc123
Built: 2025-01-12T10:00:00Z
Update available: v1.2.4 (run 'ddx self-update' to install)
```

### Internal APIs

#### Version Check Response
```go
type ReleaseInfo struct {
    Version     string    `json:"tag_name"`
    Name        string    `json:"name"`
    Prerelease  bool      `json:"prerelease"`
    PublishedAt time.Time `json:"published_at"`
    Body        string    `json:"body"` // Changelog
    Assets      []Asset   `json:"assets"`
}

type Asset struct {
    Name        string `json:"name"`
    DownloadURL string `json:"browser_download_url"`
    Size        int    `json:"size"`
}
```

## UI/UX Requirements

### Progress Indication
- Download progress bar with percentage
- Estimated time remaining
- Data transfer rate display

### Confirmation Prompts
- Clear version comparison (current → new)
- Changelog preview (first 10 lines)
- Explicit confirmation required
- Option to view full changelog

### Error Messages
- Clear, actionable error messages
- Troubleshooting suggestions
- Fallback to manual update instructions

### Success Feedback
- Version update confirmation
- What's new highlights
- Next steps guidance

## Performance Requirements

### Response Times
- Version check: <2 seconds
- Binary download: <30 seconds (10MB on 3Mbps)
- Binary replacement: <1 second
- Rollback: <1 second

### Resource Usage
- Memory: <50MB during update
- Disk: 2x binary size temporary space
- Network: Resumable downloads
- CPU: Minimal (primarily I/O bound)

### Scalability
- GitHub API rate limits (60/hour unauthenticated)
- CDN distribution for binaries
- Mirror support for enterprise environments

## Security Requirements

### Download Security
- HTTPS-only downloads
- Checksum verification (SHA256)
- Future: GPG signature verification
- No arbitrary code execution

### File System Security
- Preserve original file permissions
- No privilege escalation
- Secure temporary file handling
- Atomic operations prevent corruption

### Network Security
- Proxy support (HTTP_PROXY, HTTPS_PROXY)
- Certificate pinning (future)
- No credential transmission
- Optional GitHub token for rate limits

## Error Handling

### Network Errors
- Automatic retry with exponential backoff
- Offline mode detection
- Partial download resume
- Timeout handling

### File System Errors
- Permission denied handling
- Disk space verification
- Backup restoration on failure
- Lock file for concurrent updates

### Version Errors
- Invalid version format handling
- Downgrade prevention
- Constraint violation messages
- Pre-release warning

## Testing Strategy

### Unit Tests
- Version comparison logic
- Platform detection
- Checksum verification
- API response parsing

### Integration Tests
- GitHub API interaction
- Download functionality
- Binary replacement (test binary)
- Configuration integration

### End-to-End Tests
- Full update workflow
- Rollback scenarios
- Network failure simulation
- Platform-specific testing

### Manual Testing
- Real binary replacement
- Cross-platform verification
- Upgrade/downgrade paths
- Error recovery scenarios

## Monitoring and Analytics

### Metrics to Track
- Update check frequency
- Update success rate
- Version distribution
- Error rates by type
- Platform distribution

### Logging
- Update attempts and outcomes
- Error details for troubleshooting
- Version transition tracking
- Performance metrics

## Migration and Rollback

### Migration from Install Script
- Detect installation method
- Preserve existing configuration
- Update PATH if needed
- Handle symlinks correctly

### Rollback Strategy
- Automatic backup before update
- `ddx self-update --rollback` command
- Manual restoration instructions
- Version history tracking

## Future Enhancements

### Phase 2 Features
- Automatic background updates
- Delta updates for bandwidth saving
- Custom update channels (stable, beta, nightly)
- GPG signature verification
- Update scheduling

### Phase 3 Features
- Plugin system updates
- Dependency resolution
- A/B testing for releases
- Gradual rollout support
- Enterprise update server

## Implementation Notes

### File Replacement Strategy
On Unix systems, can replace running binary directly. On Windows, requires special handling:
1. Rename current binary to .old
2. Write new binary
3. Delete .old on next run

### Version Embedding
Use Go build flags to embed version information:
```bash
-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT}"
```

### GitHub API Rate Limiting
- Cache successful checks for 1 hour
- Support GitHub token via environment variable
- Implement exponential backoff for rate limits

### Cross-Compilation Considerations
- Build tags for platform-specific code
- Handle platform-specific path separators
- Test on all target platforms

## Success Metrics

### Launch Success Criteria
- 90% of update attempts succeed
- <1% of updates require rollback
- User satisfaction score >4.5/5
- Adoption rate >50% in first month

### Long-term Success Metrics
- Time to update reduces by 80%
- Support tickets related to updates decrease by 70%
- Version fragmentation reduces to <3 active versions
- Security patches adopted within 48 hours by >80% of users

## Risks and Mitigations

### Risk: Binary Corruption
**Mitigation**: Checksum verification, atomic replacement, automatic backup

### Risk: GitHub API Unavailability
**Mitigation**: Caching, mirrors, fallback to manual update

### Risk: Platform-Specific Issues
**Mitigation**: Extensive testing, gradual rollout, platform-specific error handling

### Risk: Security Vulnerabilities
**Mitigation**: HTTPS only, checksum verification, future signature verification

## Documentation Requirements

### User Documentation
- Update instructions in README
- Troubleshooting guide
- FAQ section
- Video tutorial

### Developer Documentation
- Implementation guide
- Testing procedures
- Release process updates
- API documentation

## Appendices

### A. Platform Binary Naming Convention
- `ddx-linux-amd64.tar.gz`
- `ddx-linux-arm64.tar.gz`
- `ddx-darwin-amd64.tar.gz`
- `ddx-darwin-arm64.tar.gz`
- `ddx-windows-amd64.zip`

### B. Example Update Flow
```bash
$ ddx self-update
Current version: v1.0.0
Latest version:  v1.1.0

Changelog:
- Added self-update feature
- Fixed bug in template application
- Improved performance

Update to v1.1.0? [y/N]: y
Downloading ddx-linux-amd64.tar.gz... 100%
Verifying checksum... OK
Installing update... Done
Successfully updated to v1.1.0
```

### C. Configuration Example
```yaml
# .ddx.yml
update:
  check_frequency: daily
  include_prereleases: false
  auto_update: false
  channel: stable
```