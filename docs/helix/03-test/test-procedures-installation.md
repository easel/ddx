---
title: "Test Procedures - Cross-Platform Installation"
type: test-procedures
feature_id: FEAT-004
workflow_phase: test
artifact_type: test-procedures
tags:
  - helix/test
  - helix/artifact/test
  - helix/phase/test
  - installation
  - cross-platform
  - procedures
related:
  - "[[TS-004-installation-test-specification]]"
  - "[[FEAT-004-cross-platform-installation]]"
  - "[[SD-004-cross-platform-installation]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Test Procedures: FEAT-004 Cross-Platform Installation

## Test Execution Strategy

This document defines the specific procedures for executing the comprehensive test suite for the cross-platform installation system. These procedures must be followed to ensure proper Test-Driven Development and HELIX compliance.

## Pre-Implementation Test Setup (Red Phase)

### Step 1: Create Test Infrastructure

```bash
# Create test directory structure
mkdir -p cli/cmd
mkdir -p cli/internal/installation
mkdir -p cli/test/fixtures/installation
mkdir -p cli/test/mock

# Create test files for each component
touch cli/cmd/installation_acceptance_test.go
touch cli/cmd/doctor.go
touch cli/cmd/self_update.go
touch cli/cmd/setup.go
touch cli/cmd/uninstall.go

# Create internal package tests
touch cli/internal/installation/platform_detector_test.go
touch cli/internal/installation/binary_distributor_test.go
touch cli/internal/installation/environment_configurer_test.go
touch cli/internal/installation/installation_manager_test.go

# Create mock services
touch cli/test/mock/github_api.go
touch cli/test/mock/network_service.go
touch cli/test/mock/shell_environment.go
```

### Step 2: Set Up Test Fixtures

Create comprehensive test data for installation scenarios:

```bash
# Create platform-specific test fixtures
mkdir -p cli/test/fixtures/installation/platforms/{linux,darwin,windows}
mkdir -p cli/test/fixtures/installation/binaries
mkdir -p cli/test/fixtures/installation/configs

# Create mock binaries for testing
echo '#!/bin/bash\necho "ddx version 1.0.0"' > cli/test/fixtures/installation/binaries/ddx-linux-amd64
echo '#!/bin/bash\necho "ddx version 1.0.0"' > cli/test/fixtures/installation/binaries/ddx-darwin-arm64
echo 'echo "ddx version 1.0.0"' > cli/test/fixtures/installation/binaries/ddx-windows-amd64.exe

chmod +x cli/test/fixtures/installation/binaries/*

# Create shell profile fixtures
mkdir -p cli/test/fixtures/installation/shell-profiles
echo 'export PATH="$PATH:/usr/local/bin"' > cli/test/fixtures/installation/shell-profiles/.bashrc
echo 'export PATH="$PATH:/usr/local/bin"' > cli/test/fixtures/installation/shell-profiles/.zshrc
echo 'set -x PATH $PATH /usr/local/bin' > cli/test/fixtures/installation/shell-profiles/config.fish
```

### Step 3: Create Mock Services

#### GitHub API Mock

```bash
# Create mock GitHub releases
cat > cli/test/fixtures/installation/mock-github-releases.json << 'EOF'
{
  "tag_name": "v1.0.0",
  "assets": [
    {
      "name": "ddx-linux-amd64.tar.gz",
      "browser_download_url": "https://github.com/example/ddx/releases/download/v1.0.0/ddx-linux-amd64.tar.gz"
    },
    {
      "name": "ddx-darwin-arm64.tar.gz",
      "browser_download_url": "https://github.com/example/ddx/releases/download/v1.0.0/ddx-darwin-arm64.tar.gz"
    },
    {
      "name": "ddx-windows-amd64.zip",
      "browser_download_url": "https://github.com/example/ddx/releases/download/v1.0.0/ddx-windows-amd64.zip"
    }
  ]
}
EOF
```

## Test Execution Procedures

### Phase 1: Red Phase Execution (Failing Tests)

#### Procedure 1.1: Run Initial Test Suite

```bash
cd cli

# Run acceptance tests (should ALL fail initially)
go test ./cmd -run "TestAcceptance_US02[8-9]|TestAcceptance_US03[0-5]" -v

# Expected result: ALL TESTS FAIL
# This confirms we're in proper Red phase
```

**Success Criteria for Red Phase:**
- All 8 acceptance tests compile successfully
- All 8 acceptance tests fail with meaningful error messages
- Error messages clearly indicate what needs implementation
- No panics or unexpected errors during test execution

#### Procedure 1.2: Validate Test Coverage

```bash
# Run with coverage to ensure test structure is correct
go test ./cmd -run "TestAcceptance_US02[8-9]|TestAcceptance_US03[0-5]" -cover -v

# Check test files exist and are properly structured
ls -la cli/cmd/*test.go | grep -E "(installation|doctor|setup|uninstall)"
```

### Phase 2: Green Phase Execution (Incremental Implementation)

#### Procedure 2.1: Implement Core Commands First

**Step 1: Implement doctor command**

```bash
# 1. Create minimal doctor.go that makes TestAcceptance_US030 pass
# 2. Run specific test to verify
go test ./cmd -run "TestAcceptance_US030_InstallationVerification" -v

# Expected: Test should now PASS
```

**Step 2: Implement setup command**

```bash
# 1. Create minimal setup.go that makes TestAcceptance_US029 pass
# 2. Run specific test to verify
go test ./cmd -run "TestAcceptance_US029_AutomaticPathConfiguration" -v

# Expected: Test should now PASS
```

**Step 3: Implement self-update command**

```bash
# 1. Create minimal self_update.go that makes TestAcceptance_US032 pass
# 2. Run specific test to verify
go test ./cmd -run "TestAcceptance_US032_UpgradeExistingInstallation" -v

# Expected: Test should now PASS
```

**Step 4: Implement uninstall command**

```bash
# 1. Create minimal uninstall.go that makes TestAcceptance_US033 pass
# 2. Run specific test to verify
go test ./cmd -run "TestAcceptance_US033_UninstallDDX" -v

# Expected: Test should now PASS
```

#### Procedure 2.2: Implement Installation Scripts

**Step 1: Update install.sh for US-028**

```bash
# 1. Modify install.sh to use GitHub releases
# 2. Add platform detection
# 3. Run installation test
go test ./cmd -run "TestAcceptance_US028_OneCommandInstallation" -v

# Expected: Unix/Linux installation tests should PASS
```

**Step 2: Create install.ps1 for Windows**

```bash
# 1. Create PowerShell installation script
# 2. Test Windows installation scenario
go test ./cmd -run "TestAcceptance_US028.*windows" -v

# Expected: Windows installation tests should PASS
```

#### Procedure 2.3: Implement Advanced Features

**Step 1: Package Manager Support (US-031)**

```bash
# 1. Create package manager detection
# 2. Create installation packages
# 3. Test package manager installation
go test ./cmd -run "TestAcceptance_US031_PackageManagerInstallation" -v

# Expected: Package manager tests should PASS
```

**Step 2: Offline Installation (US-034)**

```bash
# 1. Create offline package bundling
# 2. Implement offline installation logic
# 3. Test offline scenarios
go test ./cmd -run "TestAcceptance_US034_OfflineInstallation" -v

# Expected: Offline installation tests should PASS
```

**Step 3: Enhanced Diagnostics (US-035)**

```bash
# 1. Enhance doctor command with diagnostics
# 2. Add problem detection and remediation
# 3. Test diagnostic scenarios
go test ./cmd -run "TestAcceptance_US035_InstallationDiagnostics" -v

# Expected: Diagnostic tests should PASS
```

### Phase 3: Integration Testing

#### Procedure 3.1: Cross-Platform Testing

**Linux Testing:**

```bash
# Set up Linux test environment (Docker or VM)
docker run -it --rm -v $(pwd):/workspace ubuntu:20.04 bash

# Inside container:
cd /workspace
./install.sh
ddx doctor
ddx version
```

**macOS Testing:**

```bash
# On macOS system:
curl -sSL https://raw.githubusercontent.com/yourusername/ddx/main/install.sh | sh
ddx doctor
ddx version
```

**Windows Testing:**

```powershell
# In PowerShell:
iwr -useb https://raw.githubusercontent.com/yourusername/ddx/main/install.ps1 | iex
ddx doctor
ddx version
```

#### Procedure 3.2: Performance Testing

```bash
# Test installation speed
time ./install.sh

# Should complete in < 60 seconds on 10Mbps connection
# Measure with different network speeds using network throttling
```

#### Procedure 3.3: Reliability Testing

```bash
# Test interrupted installation
./install.sh &
INSTALL_PID=$!
sleep 5
kill $INSTALL_PID

# Verify no corrupted state
ddx doctor || echo "Clean failure state confirmed"

# Test retry mechanism
./install.sh  # Should work cleanly
```

### Phase 4: Regression Testing

#### Procedure 4.1: Full Test Suite Execution

```bash
# Run complete test suite
cd cli
go test ./cmd -v

# All installation-related tests should PASS
# No existing functionality should be broken
```

#### Procedure 4.2: End-to-End Workflow Testing

```bash
# Test complete user journey:

# 1. Fresh installation
curl -sSL https://ddx.dev/install | sh

# 2. First use
ddx init my-project
cd my-project
ddx list

# 3. Update
ddx self-update

# 4. Verify after update
ddx doctor

# 5. Uninstall
ddx uninstall --preserve-data

# 6. Verify clean removal
which ddx  # Should return not found
```

## Test Environment Management

### Environment Setup Scripts

```bash
# Create automated test environment setup
cat > scripts/setup-test-env.sh << 'EOF'
#!/bin/bash
set -e

echo "Setting up installation test environment..."

# Clean previous test artifacts
rm -rf /tmp/ddx-test-*

# Create fresh test directory
TEST_DIR="/tmp/ddx-test-$(date +%s)"
mkdir -p "$TEST_DIR"

# Copy test fixtures
cp -r cli/test/fixtures/installation/* "$TEST_DIR/"

echo "Test environment ready at: $TEST_DIR"
EOF

chmod +x scripts/setup-test-env.sh
```

### Mock Service Management

```bash
# Start mock GitHub API server for testing
cat > scripts/start-mock-services.sh << 'EOF'
#!/bin/bash

# Start mock GitHub API
go run cli/test/mock/github_api.go &
GITHUB_MOCK_PID=$!

# Start mock download server
go run cli/test/mock/download_server.go &
DOWNLOAD_MOCK_PID=$!

echo "Mock services started:"
echo "GitHub API: http://localhost:8080"
echo "Download Server: http://localhost:8081"

# Store PIDs for cleanup
echo "$GITHUB_MOCK_PID" > /tmp/github-mock.pid
echo "$DOWNLOAD_MOCK_PID" > /tmp/download-mock.pid
EOF

chmod +x scripts/start-mock-services.sh
```

## Quality Gates

### Test Execution Checklist

#### Red Phase Checklist
- [ ] All 8 acceptance tests written and compile successfully
- [ ] All tests fail with meaningful error messages
- [ ] Test fixtures and mocks are properly set up
- [ ] Test coverage reports show test structure coverage
- [ ] No implementation code exists yet

#### Green Phase Checklist
- [ ] Tests pass incrementally as features are implemented
- [ ] No test is made to pass without proper implementation
- [ ] Each feature implementation is minimal and focused
- [ ] Integration tests pass on all target platforms
- [ ] Performance criteria are met (<60s installation)

#### Final Validation Checklist
- [ ] All 8 user stories have passing acceptance tests
- [ ] Cross-platform testing completed successfully
- [ ] Performance benchmarks met
- [ ] Regression tests pass
- [ ] Documentation is updated
- [ ] Manual testing scenarios validated

### Continuous Integration Integration

```bash
# Add to CI pipeline
cat > .github/workflows/installation-tests.yml << 'EOF'
name: Installation Tests

on: [push, pull_request]

jobs:
  test-installation:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]

    steps:
    - uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21

    - name: Run Installation Tests
      run: |
        cd cli
        go test ./cmd -run "TestAcceptance_US02[8-9]|TestAcceptance_US03[0-5]" -v

    - name: Test Actual Installation
      run: |
        ./install.sh
        ddx version
        ddx doctor
EOF
```

## Troubleshooting Guide

### Common Test Failures

1. **Network-related test failures:**
   ```bash
   # Check mock services are running
   curl http://localhost:8080/api/releases/latest

   # Restart mock services if needed
   ./scripts/start-mock-services.sh
   ```

2. **Platform detection failures:**
   ```bash
   # Check platform simulation
   uname -s
   uname -m
   ```

3. **Permission-related failures:**
   ```bash
   # Ensure test directories are writable
   ls -la /tmp/ddx-test-*

   # Fix permissions if needed
   chmod -R 755 /tmp/ddx-test-*
   ```

### Debug Mode Testing

```bash
# Run tests with verbose debugging
DDX_DEBUG=1 go test ./cmd -run "TestAcceptance_US028" -v

# Enable installation script debugging
DEBUG=1 ./install.sh

# Check detailed logs
tail -f /tmp/ddx-install.log
```

---

These procedures ensure proper TDD execution and HELIX compliance for the cross-platform installation system. All procedures must be followed in order to maintain quality and reliability standards.