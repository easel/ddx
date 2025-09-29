# Implementation Plan: IP-043 - Automatic Update Notifications

**Plan ID**: IP-043
**User Story**: US-043 - Automatic Update Notifications
**Solution Design**: SD-043
**Status**: Draft
**Created**: 2025-09-29
**Updated**: 2025-09-29

## Overview

This implementation plan details the step-by-step approach to implementing automatic update notifications for DDx, following TDD principles (Test → Build → Refactor).

---

## Prerequisites

- [ ] US-043 user story reviewed and approved
- [ ] SD-043 solution design reviewed and approved
- [ ] Test specifications written
- [ ] Development environment set up

---

## Implementation Phases

### Phase 1: TEST - Write Failing Tests (TDD Red)

#### Step 1.1: Acceptance Tests
**File**: `cli/cmd/version_acceptance_test.go`

Add test functions for US-043:
```go
func TestAcceptance_US043_AutomaticUpdateNotifications(t *testing.T) {
    t.Run("automatic_check_runs_once_per_24_hours", func(t *testing.T) { ... })
    t.Run("notification_displays_after_command", func(t *testing.T) { ... })
    t.Run("disable_via_env_var", func(t *testing.T) { ... })
    t.Run("disable_via_config", func(t *testing.T) { ... })
    t.Run("cache_prevents_excessive_checks", func(t *testing.T) { ... })
    t.Run("silent_failure_on_network_error", func(t *testing.T) { ... })
    t.Run("notification_format", func(t *testing.T) { ... })
}
```

**Expected**: All tests FAIL (Red phase)

#### Step 1.2: Unit Tests - Cache Management
**File**: `cli/internal/update/cache_test.go`

```go
func TestCache_Load(t *testing.T) { ... }
func TestCache_Save(t *testing.T) { ... }
func TestCache_IsExpired(t *testing.T) { ... }
func TestCache_CorruptedFile(t *testing.T) { ... }
func TestCache_PermissionDenied(t *testing.T) { ... }
func TestCache_CreateDirectory(t *testing.T) { ... }
```

**Expected**: Compilation errors (package doesn't exist yet)

#### Step 1.3: Unit Tests - Update Checker
**File**: `cli/internal/update/checker_test.go`

```go
func TestChecker_ShouldCheck(t *testing.T) { ... }
func TestChecker_CheckForUpdate(t *testing.T) { ... }
func TestChecker_IsUpdateAvailable(t *testing.T) { ... }
func TestChecker_RespectConfig(t *testing.T) { ... }
func TestChecker_RespectEnvVar(t *testing.T) { ... }
```

**Expected**: Compilation errors (package doesn't exist yet)

---

### Phase 2: BUILD - Create Package Structure

#### Step 2.1: Create internal/update Package
```bash
mkdir -p cli/internal/update
```

#### Step 2.2: Create Type Definitions
**File**: `cli/internal/update/types.go`

```go
package update

import "time"

// CacheData represents cached update information
type CacheData struct {
    LastCheck       time.Time `json:"last_check"`
    CurrentVersion  string    `json:"current_version"`
    LatestVersion   string    `json:"latest_version"`
    UpdateAvailable bool      `json:"update_available"`
    CheckError      string    `json:"check_error,omitempty"`
}

// UpdateCheckResult represents the result of an update check
type UpdateCheckResult struct {
    UpdateAvailable bool
    LatestVersion   string
    Error           error
}
```

**Verification**: `go build ./internal/update` should succeed

#### Step 2.3: Implement Cache Management
**File**: `cli/internal/update/cache.go`

Implement:
- `type Cache struct`
- `func NewCache() *Cache`
- `func (c *Cache) Load() error`
- `func (c *Cache) Save() error`
- `func (c *Cache) IsExpired() bool`
- `func (c *Cache) getCacheFilePath() (string, error)` - XDG Base Directory logic

**Verification**: Run `go test ./internal/update -run TestCache` - tests should start passing

#### Step 2.4: Extract Shared Version Functions
**File**: `cli/internal/update/version.go`

Move from `cmd/upgrade.go`:
```go
func FetchLatestRelease() (*GitHubRelease, error)
func NeedsUpgrade(current, latest string) (bool, error)
func ParseVersion(version string) ([3]int, error)
```

Also define:
```go
type GitHubRelease struct {
    TagName string `json:"tag_name"`
    Name    string `json:"name"`
    Body    string `json:"body"`
    HTMLURL string `json:"html_url"`
}
```

**Verification**: `go test ./internal/update -run TestVersion` should pass

---

### Phase 3: BUILD - Implement Update Checker

#### Step 3.1: Implement Checker Logic
**File**: `cli/internal/update/checker.go`

```go
package update

import (
    "context"
    "sync"

    "github.com/easel/ddx/cli/internal/config"
)

var (
    lastChecker     *Checker
    lastCheckerLock sync.RWMutex
)

type Checker struct {
    currentVersion string
    config         *config.NewConfig
    cache          *Cache
    result         *UpdateCheckResult
}

func NewChecker(version string, cfg *config.NewConfig) *Checker
func (c *Checker) CheckForUpdate(ctx context.Context) (*UpdateCheckResult, error)
func (c *Checker) ShouldCheck() bool
func (c *Checker) IsUpdateAvailable() (bool, string, error)
func GetLastChecker() *Checker
```

**Implementation Steps**:
1. Implement `NewChecker` - initialize cache
2. Implement `ShouldCheck` - check cache expiry and config
3. Implement `CheckForUpdate` - GitHub API call + cache update
4. Implement `IsUpdateAvailable` - return cached result
5. Implement `GetLastChecker` - global instance access

**Verification**: `go test ./internal/update -run TestChecker` should pass

---

### Phase 4: BUILD - Refactor upgrade.go

#### Step 4.1: Update upgrade.go to Use Shared Code
**File**: `cli/cmd/upgrade.go`

Replace:
```go
// Remove these functions (moved to internal/update):
// - fetchLatestRelease()
// - needsUpgrade()
// - parseVersion()
// - GitHubRelease type
```

Add imports:
```go
import "github.com/easel/ddx/cli/internal/update"
```

Update function calls:
```go
latestRelease, err := update.FetchLatestRelease()
needsUpgrade, err := update.NeedsUpgrade(currentVersion, latestVersion)
```

**Verification**: `go test ./cmd -run TestUpgrade` should still pass

---

### Phase 5: BUILD - Extend Configuration

#### Step 5.1: Add UpdateCheck to Config Types
**File**: `cli/internal/config/types.go`

```go
type UpdateCheckConfig struct {
    Enabled   bool   `yaml:"enabled"`
    Frequency string `yaml:"frequency"` // Duration: "24h", "12h", etc.
}

type NewConfig struct {
    Version         string              `yaml:"version" json:"version"`
    Library         *LibraryConfig      `yaml:"library" json:"library"`
    PersonaBindings map[string]string   `yaml:"persona_bindings,omitempty" json:"persona_bindings,omitempty"`
    UpdateCheck     *UpdateCheckConfig  `yaml:"update_check,omitempty" json:"update_check,omitempty"`  // NEW
}
```

#### Step 5.2: Update Default Config
**File**: `cli/internal/config/types.go`

```go
func DefaultNewConfig() *NewConfig {
    return &NewConfig{
        // ... existing fields ...
        UpdateCheck: &UpdateCheckConfig{
            Enabled:   true,
            Frequency: "24h",
        },
    }
}
```

#### Step 5.3: Update ApplyDefaults
**File**: `cli/internal/config/types.go`

```go
func (c *NewConfig) ApplyDefaults() {
    // ... existing code ...

    if c.UpdateCheck == nil {
        c.UpdateCheck = &UpdateCheckConfig{
            Enabled:   true,
            Frequency: "24h",
        }
    } else {
        if c.UpdateCheck.Frequency == "" {
            c.UpdateCheck.Frequency = "24h"
        }
    }
}
```

**Verification**: `go test ./internal/config` should pass

---

### Phase 6: BUILD - Integrate with Command Factory

#### Step 6.1: Add Update Check Hooks
**File**: `cli/cmd/command_factory.go`

**Note**: Synchronous implementation (no goroutines) - simpler and acceptable latency.

```go
import (
    "context"
    "os"

    "github.com/easel/ddx/cli/internal/update"
)

// Add field to CommandFactory struct
type CommandFactory struct {
    WorkingDir    string
    Version       string
    updateChecker *update.Checker  // NEW: Store checker instance
}

func (f *CommandFactory) NewRootCommand() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "ddx",
        Short: "Document-Driven Development eXperience",
        Long: banner + "\nA toolkit for managing development assets...",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            return f.checkForUpdates(cmd)
        },
        PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
            return f.displayUpdateNotification(cmd)
        },
    }
    // ... rest of command setup ...
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

    // Fast when cache valid (just reads file)
    // Slow once per 24h when cache expired (network call ~200-500ms)
    _, err = checker.CheckForUpdate(ctx)

    // Log errors to stderr (don't let users get stranded on old versions)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Warning: Could not check for updates: %v\n", err)
    }

    // Store in factory for PostRunE
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

**Verification**: Run any ddx command and verify no errors

---

### Phase 7: TEST - Verify All Tests Pass (TDD Green)

#### Step 7.1: Run Unit Tests
```bash
cd cli
go test ./internal/update -v
go test ./internal/config -v
```

**Expected**: All tests PASS

#### Step 7.2: Run Integration Tests
```bash
go test ./cmd -run TestUpgrade -v
go test ./cmd -run TestCommandFactory -v
```

**Expected**: All tests PASS

#### Step 7.3: Run Acceptance Tests
```bash
go test ./cmd -run TestAcceptance_US043 -v
```

**Expected**: All tests PASS

#### Step 7.4: Run Full Test Suite
```bash
go test ./... -v
```

**Expected**: All tests PASS (or existing failures unrelated to this feature)

---

## Verification Steps

### Manual Testing

#### Test 1: First Run (No Cache)
```bash
rm -rf ~/.cache/ddx
./build/ddx version
# Expected: Check runs, notification shows if update available
```

#### Test 2: Cached Result
```bash
./build/ddx version
# Immediately run again:
./build/ddx list
# Expected: No network request, uses cache
```

#### Test 3: Expired Cache
```bash
# Modify cache file timestamp to >24 hours ago
# Run command
./build/ddx version
# Expected: New check runs
```

#### Test 4: Disable via Env Var
```bash
DDX_DISABLE_UPDATE_CHECK=1 ./build/ddx version
# Expected: No check, no notification
```

#### Test 5: Disable via Config
```bash
# Edit .ddx/config.yaml:
# update_check:
#   enabled: false
./build/ddx version
# Expected: No check, no notification
```

#### Test 6: Network Failure
```bash
# Disconnect network
./build/ddx version
# Expected: Command works, no error messages
```

---

## Rollback Plan

If issues arise:

1. **Immediate Rollback**:
   - Revert commits related to update checking
   - Remove PersistentPreRunE and PersistentPostRunE hooks
   - Keep refactored version code (no breaking changes)

2. **Partial Rollback**:
   - Disable by default in config
   - Add to documentation: "Set DDX_DISABLE_UPDATE_CHECK=1"

3. **Data Cleanup**:
   - Cache files are isolated in ~/.cache/ddx/
   - Safe to delete without affecting other functionality

---

## Performance Benchmarks

Run before and after implementation:

```bash
# Benchmark command execution time
time ./build/ddx version >/dev/null
time ./build/ddx list >/dev/null

# With cache hit (should be <1ms overhead)
# With cache miss (should be <10ms overhead due to goroutine)
```

---

## Documentation Updates

After implementation complete:

1. **README.md**: Mention automatic update checks
2. **docs/cli-commands.md**: Document disable mechanisms
3. **CLAUDE.md**: Update behavior description
4. **CHANGELOG.md**: Add feature to next release notes

---

## Definition of Done

- [ ] All unit tests written and passing
- [ ] All integration tests written and passing
- [ ] All acceptance tests written and passing
- [ ] Manual testing completed successfully
- [ ] Performance benchmarks meet targets (<10ms overhead)
- [ ] Code review completed
- [ ] Documentation updated
- [ ] PR approved and merged
- [ ] Feature tested on Linux, macOS, Windows

---

## Timeline Estimate

- Phase 1 (TEST): 2-3 hours
- Phase 2-3 (BUILD - Package): 2-3 hours
- Phase 4 (BUILD - Refactor): 1 hour
- Phase 5 (BUILD - Config): 1 hour
- Phase 6 (BUILD - Integration): 2 hours
- Phase 7 (TEST - Verification): 1-2 hours

**Total**: 9-12 hours of development time

---

## Dependencies

- US-032: Upgrade command (provides version comparison code)
- GitHub Releases API availability
- XDG Base Directory specification compliance

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| GitHub API rate limiting | Medium | Cache prevents excessive calls |
| Network latency | Low | Background goroutine, non-blocking |
| Cache file corruption | Low | Recreate on error |
| Test flakiness | Medium | Mock GitHub API in tests |
| Performance degradation | Medium | Benchmark before/after |

---

*This implementation plan details the step-by-step approach for US-043: Automatic Update Notifications*