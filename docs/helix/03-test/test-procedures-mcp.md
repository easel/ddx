# MCP Test Procedures

## Document Overview

**Document ID**: TP-PROC-MCP-001  
**Feature**: FEAT-001 (MCP Server Management)  
**Version**: 1.0.0  
**Created**: 2025-01-15  

## Test Writing Procedures

### 1. Unit Test Procedures

#### Setup
```go
// test/mcp/unit_test.go
package mcp_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/yourusername/ddx/cli/internal/mcp"
)

// Test fixture setup
func setupTest(t *testing.T) (*mcp.Registry, func()) {
    t.Helper()
    
    // Create temp directory
    tmpDir := t.TempDir()
    
    // Copy test fixtures
    copyFixtures(t, tmpDir)
    
    // Create registry
    registry, err := mcp.NewRegistry(tmpDir)
    require.NoError(t, err)
    
    // Return cleanup function
    return registry, func() {
        // Cleanup is automatic with t.TempDir()
    }
}
```

#### Test Structure Template
```go
func TestFeatureName(t *testing.T) {
    // Arrange
    registry, cleanup := setupTest(t)
    defer cleanup()
    
    testCases := []struct {
        name     string
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "github",
            expected: &mcp.Server{Name: "github"},
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "",
            expected: nil,
            wantErr:  true,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Act
            result, err := registry.GetServer(tc.input.(string))
            
            // Assert
            if tc.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tc.expected, result)
            }
        })
    }
}
```

### 2. Integration Test Procedures

#### Environment Setup
```bash
#!/bin/bash
# scripts/setup-integration-tests.sh

# Create mock Claude configurations
mkdir -p ~/.claude
cp test/fixtures/configs/claude-code.json ~/.claude/settings.local.json

# Set environment variables
export CLAUDE_CODE_CONFIG=~/.claude/settings.local.json
export DDX_MCP_TEST_MODE=true

# Start test registry server
go run test/servers/mock-registry/main.go &
REGISTRY_PID=$!

# Run tests
go test -tags=integration ./cli/internal/mcp/...

# Cleanup
kill $REGISTRY_PID
rm -rf ~/.claude/settings.local.json.test
```

#### Integration Test Template
```go
// +build integration

package mcp_test

func TestInstallServerIntegration(t *testing.T) {
    // Skip if not in integration mode
    if os.Getenv("DDX_MCP_TEST_MODE") != "true" {
        t.Skip("Skipping integration test")
    }
    
    // Setup real environment
    installer := mcp.NewInstaller()
    configPath := setupRealConfig(t)
    
    // Test real installation
    err := installer.Install("github", mcp.InstallOptions{
        ConfigPath: configPath,
        Environment: map[string]string{
            "GITHUB_TOKEN": "test_token",
        },
        NoBackup: true, // For testing
    })
    
    require.NoError(t, err)
    
    // Verify real config file
    config := readRealConfig(t, configPath)
    assert.Contains(t, config, "github")
}
```

### 3. Contract Test Procedures

#### CLI Contract Testing
```go
package cmd_test

import (
    "bytes"
    "encoding/json"
    "github.com/spf13/cobra"
)

func TestListCommandContract(t *testing.T) {
    // Create command
    cmd := cmd.NewMCPCommand()
    
    // Capture output
    var stdout, stderr bytes.Buffer
    cmd.SetOut(&stdout)
    cmd.SetErr(&stderr)
    
    // Set args
    cmd.SetArgs([]string{"list", "--format", "json"})
    
    // Execute
    err := cmd.Execute()
    require.NoError(t, err)
    
    // Validate contract
    var output contracts.ListOutput
    err = json.Unmarshal(stdout.Bytes(), &output)
    require.NoError(t, err)
    
    // Contract assertions
    assert.NotNil(t, output.Total)
    assert.NotNil(t, output.Installed)
    assert.NotNil(t, output.Servers)
    
    // Validate schema
    validateJSONSchema(t, stdout.Bytes(), "schemas/list-output.json")
}
```

### 4. E2E Test Procedures

#### E2E Test Setup
```go
func TestE2EInstallFlow(t *testing.T) {
    // Start fresh environment
    e2e := setupE2EEnvironment(t)
    defer e2e.Cleanup()
    
    // Execute commands as user would
    steps := []struct {
        name    string
        command string
        args    []string
        expect  string
    }{
        {
            name:    "list servers",
            command: "mcp",
            args:    []string{"list"},
            expect:  "Available MCP Servers",
        },
        {
            name:    "install github",
            command: "mcp",
            args:    []string{"install", "github", "--env", "GITHUB_TOKEN=test", "--yes"},
            expect:  "successfully installed",
        },
        {
            name:    "check status",
            command: "mcp",
            args:    []string{"status", "github"},
            expect:  "Configured",
        },
    }
    
    for _, step := range steps {
        t.Run(step.name, func(t *testing.T) {
            output := e2e.RunCommand(step.command, step.args...)
            assert.Contains(t, output, step.expect)
        })
    }
}
```

### 5. Security Test Procedures

#### Security Test Template
```go
func TestSecurityValidation(t *testing.T) {
    validator := mcp.NewValidator()
    
    // Test injection attacks
    injectionTests := []struct {
        name    string
        input   string
        field   string
    }{
        {"SQL injection", "'; DROP TABLE servers; --", "server_name"},
        {"Command injection", "$(rm -rf /)", "env_value"},
        {"Path traversal", "../../../etc/passwd", "config_path"},
        {"XSS attempt", "<script>alert('xss')</script>", "description"},
    }
    
    for _, test := range injectionTests {
        t.Run(test.name, func(t *testing.T) {
            err := validator.ValidateField(test.field, test.input)
            assert.Error(t, err, "Should reject: %s", test.input)
            assert.Contains(t, err.Error(), "invalid")
        })
    }
    
    // Test credential masking
    t.Run("credential masking", func(t *testing.T) {
        sensitive := "ghp_secrettoken123"
        masked := mcp.MaskSensitive(sensitive)
        assert.NotContains(t, masked, "secret")
        assert.Equal(t, "ghp_***", masked)
    })
}
```

### 6. Performance Test Procedures

#### Benchmark Template
```go
func BenchmarkRegistryOperations(b *testing.B) {
    registry := setupLargeRegistry(b) // 1000 servers
    
    b.Run("Load", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            r, _ := mcp.LoadRegistry("large-registry.yml")
            _ = r
        }
    })
    
    b.Run("Search", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            results := registry.Search("git")
            if len(results) == 0 {
                b.Fatal("No results found")
            }
        }
    })
    
    b.Run("Filter", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            results := registry.FilterByCategory("development")
            if len(results) == 0 {
                b.Fatal("No results found")
            }
        }
    })
}

// Performance assertions
func TestPerformanceRequirements(t *testing.T) {
    registry := setupLargeRegistry(t)
    
    // Test load time < 50ms
    start := time.Now()
    _, err := mcp.LoadRegistry("large-registry.yml")
    elapsed := time.Since(start)
    
    require.NoError(t, err)
    assert.Less(t, elapsed, 50*time.Millisecond, 
        "Registry load took %v, expected < 50ms", elapsed)
}
```

## Test Execution Procedures

### Local Development Testing

```bash
# Run all unit tests
make test-mcp-unit

# Run with coverage
make test-mcp-coverage

# Run specific test
go test -v -run TestRegistryLoad ./cli/internal/mcp/...

# Run with race detection
go test -race ./cli/internal/mcp/...

# Run benchmarks
go test -bench=. -benchmem ./cli/internal/mcp/...
```

### CI/CD Testing

```yaml
# .github/workflows/mcp-tests.yml
name: MCP Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make test-mcp-unit
      - uses: codecov/codecov-action@v3

  integration-tests:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v3
      - run: make test-mcp-integration

  security-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make test-mcp-security
      - run: gosec ./cli/internal/mcp/...
```

## Test Data Management

### Fixture Organization
```
test/fixtures/mcp/
├── registry/
│   ├── valid/
│   │   ├── github.yml
│   │   ├── postgres.yml
│   │   └── registry.yml
│   ├── invalid/
│   │   ├── missing-name.yml
│   │   ├── bad-command.yml
│   │   └── malformed.yml
│   └── large/
│       └── 1000-servers.yml
├── configs/
│   ├── claude-code-empty.json
│   ├── claude-code-with-servers.json
│   ├── claude-desktop-empty.json
│   └── corrupted.json
└── credentials/
    ├── valid-tokens.json
    └── invalid-tokens.json
```

### Test Data Generation
```go
// test/generators/registry.go
func GenerateLargeRegistry(numServers int) *mcp.Registry {
    registry := &mcp.Registry{
        Version: "1.0.0",
        Servers: make([]mcp.Server, numServers),
    }
    
    categories := []string{"development", "database", "filesystem", "productivity"}
    
    for i := 0; i < numServers; i++ {
        registry.Servers[i] = mcp.Server{
            Name:        fmt.Sprintf("server-%d", i),
            Description: fmt.Sprintf("Test server %d", i),
            Category:    categories[i%len(categories)],
            Command: mcp.CommandSpec{
                Executable: "echo",
                Args:       []string{"test"},
            },
        }
    }
    
    return registry
}
```

## Environment Configuration

### Test Environment Variables
```bash
# test.env
DDX_MCP_TEST_MODE=true
DDX_MCP_REGISTRY=test/fixtures/mcp/registry/valid/registry.yml
DDX_MCP_CACHE_DIR=/tmp/ddx-mcp-test-cache
DDX_MCP_NO_COLOR=true
CLAUDE_CODE_CONFIG=test/fixtures/mcp/configs/claude-code-empty.json
```

### Mock Services
```go
// test/mocks/claude_config.go
type MockClaudeConfig struct {
    mock.Mock
}

func (m *MockClaudeConfig) Load(path string) (*ClaudeConfig, error) {
    args := m.Called(path)
    return args.Get(0).(*ClaudeConfig), args.Error(1)
}

func (m *MockClaudeConfig) Save(path string, config *ClaudeConfig) error {
    args := m.Called(path, config)
    return args.Error(0)
}
```

## Test Quality Checks

### Pre-Test Validation
```bash
#!/bin/bash
# scripts/validate-tests.sh

# Check test coverage
go test -cover ./cli/internal/mcp/... | grep -E "coverage: [0-9]+\.[0-9]+%"

# Check for focused tests
if grep -r "FIt\|FDescribe\|FContext" test/; then
    echo "Error: Focused tests found"
    exit 1
fi

# Check for skipped tests
if grep -r "t.Skip" test/ | grep -v "integration"; then
    echo "Warning: Skipped tests found"
fi

# Lint test code
golangci-lint run test/...
```

### Test Review Checklist

- [ ] Test names are descriptive
- [ ] AAA pattern followed (Arrange, Act, Assert)
- [ ] Error cases tested
- [ ] Edge cases covered
- [ ] Cleanup code present
- [ ] No hardcoded values
- [ ] Assertions are specific
- [ ] No test interdependencies

## Troubleshooting Guide

### Common Issues

#### Issue: Tests fail on CI but pass locally
**Solution**: Check for environment differences
```bash
# Debug CI environment
go test -v -tags=debug ./cli/internal/mcp/...
```

#### Issue: Flaky tests
**Solution**: Add retry logic for network operations
```go
func TestWithRetry(t *testing.T) {
    require.Eventually(t, func() bool {
        err := networkOperation()
        return err == nil
    }, 5*time.Second, 100*time.Millisecond)
}
```

#### Issue: Test data pollution
**Solution**: Use t.TempDir() and proper cleanup
```go
func TestWithCleanData(t *testing.T) {
    dir := t.TempDir() // Automatically cleaned up
    // Use dir for test data
}
```

## Test Reporting

### Coverage Report Generation
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./cli/internal/mcp/...
go tool cover -html=coverage.out -o coverage.html

# Generate test report
go test -json ./cli/internal/mcp/... > test-results.json
go-test-report < test-results.json > test-report.html
```

### Test Metrics
```go
// test/metrics/collector.go
type TestMetrics struct {
    TotalTests     int
    PassedTests    int
    FailedTests    int
    SkippedTests   int
    Coverage       float64
    Duration       time.Duration
    FailureReasons map[string]int
}

func CollectMetrics(results []TestResult) *TestMetrics {
    // Aggregate test results
    // Generate metrics report
}
```

---

*These procedures ensure consistent, reliable testing of the MCP server management feature.*