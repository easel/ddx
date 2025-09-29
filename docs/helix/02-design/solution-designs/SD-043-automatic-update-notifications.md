# Solution Design: SD-043 - Automatic Update Notifications

**Design ID**: SD-043
**User Story**: US-043 - Automatic Update Notifications
**Status**: Draft
**Created**: 2025-09-29
**Updated**: 2025-09-29

## Overview

### Problem Statement
Users need to be notified when DDx updates are available without manually checking, while avoiding excessive network requests and maintaining a smooth, non-intrusive user experience.

### Solution Summary
Implement an automatic update notification system that:
1. Checks for updates once per 24 hours when any command runs
2. Caches check results to avoid GitHub API rate limits
3. Displays notifications after command completion
4. Can be disabled via environment variable or configuration
5. Fails silently on network errors to avoid disrupting workflows

---

## Architecture

### Component Overview

```
┌─────────────────────────────────────────────────────────┐
│                    DDx CLI                               │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │         command_factory.go                      │   │
│  │                                                  │   │
│  │  PreRunE: Check if update check needed ────────┼───┼─► internal/update
│  │  PostRunE: Display notification if available   │   │       │
│  └────────────────────────────────────────────────┘   │       │
│                                                          │       │
└──────────────────────────────────────────────────────────┘       │
                                                                   │
┌──────────────────────────────────────────────────────────────────┘
│
│  internal/update/
│  ├── checker.go         Check logic, GitHub API interaction
│  ├── cache.go           Cache file management, TTL logic
│  ├── types.go           Shared types (UpdateCheckResult, etc.)
│  └── version.go         Version comparison (from upgrade.go)
│
└──► GitHub Releases API (https://api.github.com/repos/easel/ddx/releases/latest)
     └──► Cache: ~/.cache/ddx/last-update-check.json
```

### Data Flow

1. **Command Execution Begins**
   - Root command PreRunE hook executes (`checkForUpdates`)
   - Checks environment variable `DDX_DISABLE_UPDATE_CHECK`
   - Loads configuration (check if `update_check.enabled`)

2. **Update Check Decision**
   - Attempts to load cache from `~/.cache/ddx/last-update-check.json`
   - If cache exists and < 24 hours old → Use cached result (instant)
   - If cache missing/expired/invalid → Perform new check (network call ~200-500ms)

3. **GitHub API Query** (if check needed)
   - **Synchronous call**: `GET https://api.github.com/repos/easel/ddx/releases/latest`
   - Compare returned version with current binary version
   - Store result in cache with current timestamp
   - Store checker in CommandFactory for PostRunE

4. **Command Executes**
   - User's command runs normally
   - Once per 24 hours: slight delay (200-500ms) from network check
   - All other times: instant (just read cache file)

5. **Notification Display**
   - Root command PostRunE hook executes (`displayUpdateNotification`)
   - Retrieves checker from CommandFactory
   - If update available → Display: `⬆️  Update available: v0.1.3 (run 'ddx upgrade' to install)`
   - If no update → Silent (no output)

---

## Detailed Design

### Cache Structure

**Location**: `~/.cache/ddx/last-update-check.json` (XDG Base Directory Specification)

**Schema**:
```json
{
  "last_check": "2025-09-29T21:00:00Z",
  "current_version": "v0.1.2",
  "latest_version": "v0.1.3",
  "update_available": true,
  "check_error": null
}
```

**Fields**:
- `last_check`: RFC3339 timestamp of last successful check
- `current_version`: Version of binary when check was performed
- `latest_version`: Latest version from GitHub API
- `update_available`: Boolean indicating if update exists
- `check_error`: Error message from last check (null if successful)

**TTL Logic**:
```go
func (c *UpdateCache) IsExpired() bool {
    return time.Since(c.LastCheck) > 24*time.Hour
}
```

### Configuration Schema

**File**: `.ddx/config.yaml`

```yaml
update_check:
  enabled: true      # Enable/disable automatic checks
  frequency: 24h     # Check frequency (Go duration format)
```

**Type Definition**:
```go
type UpdateCheckConfig struct {
    Enabled   bool   `yaml:"enabled"`
    Frequency string `yaml:"frequency"` // Duration: "24h", "12h", etc.
}
```

### API Contracts

#### Internal API: update.Checker

```go
package update

// Checker handles update checking logic
type Checker struct {
    config      *config.NewConfig
    cache       *Cache
    currentVer  string
    result      *UpdateCheckResult  // Stores check result
}

// CheckForUpdate determines if check is needed and performs it
// Returns immediately with cached result if cache is fresh
// Performs network call if cache expired (once per 24h)
func (c *Checker) CheckForUpdate(ctx context.Context) (*UpdateCheckResult, error)

// ShouldCheck returns true if check is needed (cache expired)
func (c *Checker) ShouldCheck() bool

// IsUpdateAvailable returns result from last CheckForUpdate call
func (c *Checker) IsUpdateAvailable() (bool, string, error)
```

#### Internal API: update.Cache

```go
package update

// Cache manages the update check cache file
type Cache struct {
    filePath string
    data     *CacheData
}

// CacheData represents cached update information
type CacheData struct {
    LastCheck       time.Time `json:"last_check"`
    CurrentVersion  string    `json:"current_version"`
    LatestVersion   string    `json:"latest_version"`
    UpdateAvailable bool      `json:"update_available"`
    CheckError      string    `json:"check_error,omitempty"`
}

// Load reads cache from disk
func (c *Cache) Load() error

// Save writes cache to disk
func (c *Cache) Save() error

// IsExpired checks if cache is older than TTL
func (c *Cache) IsExpired() bool
```

#### Shared Functions (from upgrade.go)

```go
// FetchLatestRelease queries GitHub API
func FetchLatestRelease() (*GitHubRelease, error)

// NeedsUpgrade compares semantic versions
func NeedsUpgrade(current, latest string) (bool, error)

// ParseVersion parses semantic version string
func ParseVersion(version string) ([3]int, error)
```

---

## Implementation Details

### Command Factory Integration

**File**: `cli/cmd/command_factory.go`

**Design Decision**: Synchronous check (no goroutines) because:
- Check happens at most once per 24 hours
- Most invocations just read cache (instant)
- Network check adds ~200-500ms max (acceptable for once-daily check)
- Simpler implementation, no race conditions or global state

```go
func (f *CommandFactory) NewRootCommand() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "ddx",
        Short: "Document-Driven Development eXperience",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            return f.checkForUpdates(cmd)
        },
        PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
            return f.displayUpdateNotification(cmd)
        },
    }
    // ... rest of command setup
    return rootCmd
}

func (f *CommandFactory) checkForUpdates(cmd *cobra.Command) error {
    // Check if disabled via env var
    if os.Getenv("DDX_DISABLE_UPDATE_CHECK") == "1" {
        return nil
    }

    // Load config
    cfg, err := config.Load()
    if err != nil {
        // Silent failure - use defaults
        cfg = config.DefaultNewConfig()
    }

    // Check if disabled via config
    if cfg.UpdateCheck != nil && !cfg.UpdateCheck.Enabled {
        return nil
    }

    // Create checker and perform check (synchronous)
    checker := update.NewChecker(f.Version, cfg)
    ctx := context.Background()

    // This is fast when cache is valid (just reads file)
    // Only slow once per 24 hours when cache expired (network call)
    _, err = checker.CheckForUpdate(ctx)

    // Log errors to stderr (don't let users get stranded on old versions)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Warning: Could not check for updates: %v\n", err)
    }

    // Store in factory for PostRunE to access
    f.updateChecker = checker

    // Never return error - don't disrupt workflow
    return nil
}

func (f *CommandFactory) displayUpdateNotification(cmd *cobra.Command) error {
    if f.updateChecker == nil {
        return nil
    }

    available, version, err := f.updateChecker.IsUpdateAvailable()
    if err != nil || !available {
        return nil
    }

    fmt.Fprintf(cmd.OutOrStdout(),
        "\n⬆️  Update available: %s (run 'ddx upgrade' to install)\n",
        version)

    return nil
}
```

### Error Handling Strategy

| Error Scenario | Behavior | Logging |
|----------------|----------|---------|
| Network unavailable | Silent to user, use cached result if available | Log to stderr: "Warning: Could not check for updates (network unavailable)" |
| GitHub API rate limited | Silent to user, retry on next check | Log to stderr: "Warning: Update check rate limited, will retry later" |
| Cache file corrupted | Recreate cache file, perform fresh check | Log to stderr: "Warning: Update cache corrupted, recreating" |
| Cache file permission denied | Silent to user, skip update checks entirely | Log to stderr: "Warning: Cannot access update cache (permission denied)" |
| Invalid version format | Silent to user, skip check | Log to stderr: "Warning: Invalid version format, skipping update check" |
| Configuration load failure | Use default settings (enabled=true) | No log (expected in some cases) |

**Key Principles**:
- Never disrupt user workflow due to update check failures
- Log all check failures to stderr so users aren't stranded on old versions
- Logs are visible but non-intrusive (stderr, prefixed with "Warning:")
- Users can diagnose issues with `ddx doctor` or by reviewing stderr

---

## Performance Considerations

### Optimization Strategies

1. **Cache-First Approach**
   - Check cache before any network call
   - 24-hour TTL prevents excessive API usage (1 request per day max)
   - Respects GitHub API rate limits (60/hour unauthenticated)
   - Cache read is nearly instant (<1ms)

2. **Synchronous with Acceptable Latency**
   - Cache hit (99.96% of invocations): <1ms overhead
   - Cache miss (0.04% of invocations, once per day): 200-500ms acceptable
   - Simpler than async: no goroutines, channels, or race conditions

3. **Lazy Loading**
   - Only load cache when enabled
   - Don't parse JSON if disabled via env var or config

4. **Minimal Memory Footprint**
   - Cache struct: ~200 bytes
   - Checker instance: ~1KB
   - Stored in CommandFactory, GC'd after command

### Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Check overhead (cached) | <1ms | Time to read and parse cache file |
| Check overhead (network) | 200-500ms | Acceptable once per 24 hours |
| Memory usage | <1MB | Total heap allocation for update logic |
| Network bandwidth | <10KB | Size of GitHub API response |
| Disk I/O | <1KB | Cache file size |

### Performance Justification

With 24-hour caching:
- User runs 100 commands per day
- 1 check requires network (500ms)
- 99 checks use cache (<1ms each)
- Average overhead: ~5ms per command
- Acceptable tradeoff for update notifications

---

## Security Considerations

### Data Security

1. **HTTPS Only**
   - All GitHub API calls use HTTPS
   - No plaintext version information transmitted

2. **No Authentication Required**
   - Uses public GitHub API endpoint
   - No tokens or credentials stored

3. **Cache File Permissions**
   - Cache directory: 0755 (drwxr-xr-x)
   - Cache file: 0644 (rw-r--r--)
   - No sensitive information in cache

### Privacy

1. **No Telemetry**
   - No user tracking
   - No analytics sent
   - Only version comparison performed

2. **User Control**
   - Easy to disable via env var
   - Easy to disable via config
   - Transparent behavior

### Rate Limiting

1. **GitHub API Limits**
   - Unauthenticated: 60 requests/hour
   - With cache: ~1 request/day per user
   - Shared IP addresses (e.g., corporate networks) protected by cache

---

## Testing Strategy

### Unit Tests

**File**: `cli/internal/update/cache_test.go`
- Test cache file creation and loading
- Test TTL expiration logic
- Test cache corruption handling
- Test permission error handling

**File**: `cli/internal/update/checker_test.go`
- Test update check logic with mocked GitHub API
- Test version comparison
- Test configuration override
- Test environment variable override

### Integration Tests

**File**: `cli/cmd/command_factory_test.go`
- Test PreRunE hook with cache
- Test PostRunE notification display
- Test interaction between multiple commands

### Acceptance Tests

**File**: `cli/cmd/version_acceptance_test.go`
- Test automatic check on version command (US-043)
- Test notification display format
- Test disable via env var
- Test disable via config
- Test cache behavior

---

## Migration Strategy

### Backwards Compatibility

**No Breaking Changes**:
- New feature, no existing functionality modified
- Configuration is additive (new `update_check` section)
- Defaults to enabled (matches user expectations)

### Existing Users

1. **First Run After Update**
   - No cache exists → Check runs
   - If update available → Notification shown
   - Cache created for future runs

2. **Configuration**
   - No configuration changes required
   - Works with default settings
   - Can opt-out if desired

---

## Deployment Considerations

### Rollout Strategy

1. **Phase 1**: Deploy with feature enabled by default
2. **Phase 2**: Monitor GitHub API usage and performance
3. **Phase 3**: Gather user feedback on notification UX
4. **Phase 4**: Adjust frequency or behavior based on feedback

### Monitoring

**Metrics to Track**:
- Cache hit rate (should be >95%)
- GitHub API request volume
- User opt-out rate
- Performance impact (command execution time)

**Success Criteria**:
- <5% increase in GitHub API requests
- <10ms overhead on command execution
- <5% user opt-out rate
- Zero crashes or errors related to update checks

---

## Future Enhancements

### Phase 2 Features

1. **Configurable Frequency**
   - Support durations: 12h, 6h, weekly
   - Validate duration format in config

2. **Changelog Preview**
   - Include first 3 lines of release notes
   - Format: "What's new: - Feature X, - Fix Y, ..."

3. **Pre-release Opt-In**
   - Configuration: `update_check.include_prerelease: true`
   - Show beta/RC versions

4. **Update Channels**
   - Stable (default)
   - Beta
   - Nightly

### Phase 3 Features

1. **Smart Notification Timing**
   - Only show once per session
   - Remember dismissed notifications

2. **Automatic Background Updates**
   - Configuration: `update_check.auto_update: true`
   - Download in background, prompt to restart

---

## Open Questions

1. **Should we randomize check timing within 24-hour window?**
   - **Decision**: No (Phase 1), revisit if API load becomes issue

2. **Should we cache the notification text?**
   - **Decision**: Yes, store in memory during command execution

3. **Should we support GitHub authentication tokens?**
   - **Decision**: No (Phase 1), users unlikely to hit rate limits

4. **Should we show notification on every command or just once?**
   - **Decision**: Every command (unless cached check shows no update)

---

## Dependencies

- GitHub Releases API (public, no authentication)
- XDG Base Directory specification for cache location
- Go standard library: `encoding/json`, `net/http`, `time`, `os`
- Internal: `internal/config` for configuration loading
- Command: `cmd/upgrade.go` for shared version logic

---

## References

- US-043: Automatic Update Notifications
- US-008: Check DDX Version
- Feature Spec: `/docs/product/features/self-update.md`
- XDG Base Directory: https://specifications.freedesktop.org/basedir-spec/
- GitHub API: https://docs.github.com/en/rest/releases

---

*This solution design documents the architecture for US-043: Automatic Update Notifications*