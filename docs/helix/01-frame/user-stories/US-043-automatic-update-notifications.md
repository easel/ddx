# User Story: US-043 - Automatic Update Notifications

**Story ID**: US-043
**Feature**: FEAT-004 - Cross-Platform Installation
**Priority**: P2
**Status**: Draft
**Created**: 2025-09-29
**Updated**: 2025-09-29

## Story

**As a** DDx user
**I want** to be notified automatically when updates are available
**So that** I can stay current with features and security fixes without manual checking

## Acceptance Criteria

- [ ] **Given** I run any ddx command, **when** more than 24 hours have passed since last check, **then** the system checks for updates in the background
- [ ] **Given** an update is available, **when** my command completes, **then** a notification displays: "⬆️  Update available: vX.Y.Z (run 'ddx upgrade' to install)"
- [ ] **Given** the update check fails, **when** network is unavailable, **then** the failure is silent and doesn't disrupt my workflow
- [ ] **Given** I set DDX_DISABLE_UPDATE_CHECK=1, **when** I run any command, **then** no update check is performed
- [ ] **Given** I configure update_check.enabled: false, **when** I run any command, **then** no update check is performed
- [ ] **Given** I checked for updates recently, **when** I run another command within 24 hours, **then** the cached result is used (no network request)
- [ ] **Given** an update check runs, **when** I measure command performance, **then** overhead is less than 10ms
- [ ] **Given** I use any platform, **when** update checks run, **then** they work correctly on Linux, macOS, and Windows

## Definition of Done

- [ ] User story fully implemented
- [ ] Cache system stores last check timestamp
- [ ] 24-hour TTL prevents excessive API calls
- [ ] Environment variable override works
- [ ] Configuration file override works
- [ ] Notification displays after command completion
- [ ] Silent failure on network errors
- [ ] No performance degradation
- [ ] Unit tests written and passing
- [ ] Integration tests written and passing
- [ ] Acceptance tests written and passing
- [ ] Documentation updated (README, cli-commands.md)
- [ ] Cross-platform testing completed

## Technical Notes

### Implementation Considerations

- **Cache Location**: Use XDG Base Directory specification (`~/.cache/ddx/`)
- **Cache Format**: JSON file with timestamp, version info, check result
- **Cache TTL**: Default 24 hours (configurable)
- **Check Timing**: PreRunE hook in root command (non-blocking)
- **Display Timing**: PostRunE hook in root command (after output)
- **Shared Code**: Reuse version comparison and GitHub API code from upgrade.go
- **Rate Limiting**: Cache prevents hitting GitHub API rate limits (60/hour unauthenticated)

### Cache Structure

```json
{
  "last_check": "2025-09-29T21:00:00Z",
  "latest_version": "v0.1.2",
  "current_version": "v0.1.1",
  "update_available": true
}
```

### Configuration Schema

```yaml
# .ddx/config.yaml
update_check:
  enabled: true
  frequency: 24h  # Go duration format
```

### Error Scenarios

- Network unavailable during check → Silent failure, use cached result if available
- GitHub API rate limited → Silent failure, retry on next check
- Cache file corrupted → Recreate cache file
- Invalid cache timestamp → Treat as expired, perform new check
- File system permission denied → Silent failure, skip update checks

## Validation Scenarios

### Scenario 1: First Run Update Check
1. Install DDx fresh
2. Run `ddx version`
3. **Expected**: Update check runs (no cache exists), notification shows if update available

### Scenario 2: Cached Result Within 24 Hours
1. Run `ddx version` (check runs, cache created)
2. Immediately run `ddx list`
3. **Expected**: No network request, uses cached result, notification shows if applicable

### Scenario 3: Expired Cache
1. Run `ddx version` (cache created with timestamp)
2. Modify cache file to have timestamp >24 hours ago
3. Run `ddx version` again
4. **Expected**: New check runs, cache updated with current timestamp

### Scenario 4: Disable via Environment Variable
1. Set `export DDX_DISABLE_UPDATE_CHECK=1`
2. Run any ddx command
3. **Expected**: No update check, no notification, no cache file created

### Scenario 5: Disable via Configuration
1. Set `update_check.enabled: false` in `.ddx/config.yaml`
2. Run any ddx command
3. **Expected**: No update check, no notification

### Scenario 6: Silent Network Failure
1. Disconnect network
2. Run `ddx version` (cache expired or missing)
3. **Expected**: Command executes normally, no error messages, no notification

### Scenario 7: Update Notification Format
1. Ensure newer version exists on GitHub
2. Run `ddx version` with expired cache
3. **Expected**: After version output, see: "⬆️  Update available: v0.1.3 (run 'ddx upgrade' to install)"

### Scenario 8: No Update Available
1. Ensure running latest version
2. Run any command with expired cache
3. **Expected**: Check runs, no notification displayed, cache updated

## User Personas

### Primary: Busy Developer
- **Role**: Full-stack developer using DDx daily
- **Goals**: Stay updated without manual effort
- **Pain Points**: Forgetting to check for updates, missing important fixes
- **Technical Level**: Intermediate to advanced
- **Behavior**: Runs DDx commands frequently throughout the day

### Secondary: Security-Conscious Team Lead
- **Role**: Team lead responsible for tool security
- **Goals**: Ensure team uses patched versions
- **Pain Points**: Version fragmentation across team, security vulnerabilities
- **Technical Level**: Advanced
- **Behavior**: Wants automatic notifications but may customize frequency

### Tertiary: CI/CD Pipeline
- **Role**: Automated build system
- **Goals**: Predictable behavior, no interactive prompts
- **Pain Points**: Unwanted network requests, rate limiting
- **Technical Level**: N/A (automated)
- **Behavior**: May disable update checks via environment variable

## Dependencies

- US-032: Upgrade Existing Installation (provides upgrade.go code to share)
- GitHub Releases API (for version checking)
- File system access (for cache storage)
- Configuration system (for enable/disable settings)

## Related Stories

- US-008: Check DDX Version (automatic checking satisfies AC)
- US-032: Upgrade Existing Installation (provides upgrade command to run)
- FEAT-004: Cross-Platform Installation (parent feature)

## Success Metrics

- **Adoption Rate**: >80% of users see update notifications within 48 hours of release
- **API Usage**: <5% increase in GitHub API requests (cache effectiveness)
- **Performance**: <10ms overhead on command execution
- **User Satisfaction**: No complaints about intrusive notifications
- **Opt-Out Rate**: <5% of users disable update checks

## Open Questions

1. Should we randomize check timing within the 24-hour window to spread API load?
2. Should the notification include a brief changelog or just the version number?
3. Should we cache the notification text to show on subsequent commands until upgraded?
4. Should we support different channels (stable, beta) in the future?

---
*This user story is part of FEAT-004: Cross-Platform Installation*