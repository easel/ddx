# DDx Test Procedures

**Document Type**: Test Procedures
**Version**: 1.0.0
**Status**: Draft
**Created**: 2025-01-15
**Updated**: 2025-01-15

## 1. Introduction

This document provides step-by-step procedures for writing, executing, and maintaining tests for the DDx CLI toolkit. These procedures ensure consistent, high-quality test development following TDD principles.

## 2. Test Development Procedures

### 2.1 Unit Test Development

#### Procedure: Writing Unit Tests

**Prerequisites:**
- Design documents reviewed
- Test plan approved
- Development environment setup

**Steps:**

1. **Identify test targets**
   ```bash
   # List all Go files needing tests
   find . -name "*.go" ! -name "*_test.go" | grep -E "internal|cmd"
   ```

2. **Create test file**
   ```bash
   # For each source file, create corresponding test
   # Example: config.go → config_test.go
   touch internal/config/config_test.go
   ```

3. **Write test structure**
   ```go
   package config

   import (
       "testing"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
   )

   func TestConfigLoad_ValidFile(t *testing.T) {
       // Test implementation
   }
   ```

4. **Implement table-driven tests**
   ```go
   func TestValidation(t *testing.T) {
       tests := []struct {
           name    string
           input   Config
           wantErr bool
       }{
           {
               name:    "valid config",
               input:   Config{Repository: "https://github.com/ddx/ddx"},
               wantErr: false,
           },
           {
               name:    "missing repository",
               input:   Config{},
               wantErr: true,
           },
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               err := tt.input.Validate()
               if tt.wantErr {
                   assert.Error(t, err)
               } else {
                   assert.NoError(t, err)
               }
           })
       }
   }
   ```

5. **Run tests (expect failure)**
   ```bash
   go test ./internal/config -v
   # Should fail - no implementation yet
   ```

6. **Document test purpose**
   ```go
   // TestConfigLoad_ValidFile validates that valid YAML configurations
   // are loaded correctly with all fields populated and defaults applied.
   func TestConfigLoad_ValidFile(t *testing.T) {
       // ...
   }
   ```

### 2.2 Integration Test Development

#### Procedure: Writing Integration Tests

1. **Setup test environment**
   ```go
   func setupTestEnvironment(t *testing.T) string {
       tempDir := t.TempDir()

       // Initialize git repo
       cmd := exec.Command("git", "init")
       cmd.Dir = tempDir
       require.NoError(t, cmd.Run())

       return tempDir
   }
   ```

2. **Write command tests**
   ```go
   func TestInitCommand_NewProject(t *testing.T) {
       // Arrange
       testDir := setupTestEnvironment(t)
       os.Chdir(testDir)
       defer os.Chdir("..")

       // Act
       cmd := exec.Command("ddx", "init", "--template", "nextjs")
       output, err := cmd.CombinedOutput()

       // Assert
       require.NoError(t, err)
       assert.Contains(t, string(output), "Initialized DDx")
       assert.FileExists(t, ".ddx.yml")
   }
   ```

3. **Test error scenarios**
   ```go
   func TestInitCommand_ExistingProject(t *testing.T) {
       testDir := setupTestEnvironment(t)

       // Create existing .ddx.yml
       require.NoError(t, os.WriteFile(
           filepath.Join(testDir, ".ddx.yml"),
           []byte("version: 1.0.0"),
           0644,
       ))

       // Should fail without --force
       cmd := exec.Command("ddx", "init")
       _, err := cmd.CombinedOutput()
       assert.Error(t, err)
   }
   ```

### 2.3 Contract Test Development

#### Procedure: Validating API Contracts

1. **Load contract specification**
   ```go
   type CLIContract struct {
       Command  string
       Args     []string
       Expected ExpectedOutput
   }

   func loadContract(t *testing.T) []CLIContract {
       // Load from design/contracts/CLI-001-core-commands.md
       // Parse and return contract specifications
   }
   ```

2. **Test contract compliance**
   ```go
   func TestCLI001_InitCommand_Contract(t *testing.T) {
       contracts := loadContract(t)

       for _, contract := range contracts {
           t.Run(contract.Command, func(t *testing.T) {
               // Execute command
               cmd := exec.Command("ddx", contract.Args...)
               output, err := cmd.CombinedOutput()

               // Validate output format
               assert.Equal(t, contract.Expected.ExitCode, cmd.ProcessState.ExitCode())
               assert.Regexp(t, contract.Expected.OutputPattern, string(output))
           })
       }
   }
   ```

3. **Validate error codes**
   ```go
   func TestExitCodes(t *testing.T) {
       tests := []struct {
           args     []string
           expected int
       }{
           {[]string{"init"}, 0},                    // Success
           {[]string{"init", "--invalid"}, 1},       // Invalid flag
           {[]string{"apply", "nonexistent"}, 2},    // Resource not found
       }

       for _, tt := range tests {
           cmd := exec.Command("ddx", tt.args...)
           cmd.Run()
           assert.Equal(t, tt.expected, cmd.ProcessState.ExitCode())
       }
   }
   ```

### 2.4 End-to-End Test Development

#### Procedure: Writing E2E Tests

1. **Define user journey**
   ```go
   func TestE2E_CompleteWorkflow(t *testing.T) {
       // User story: Initialize project, apply template, update, contribute

       // Step 1: Initialize
       runCommand(t, "ddx", "init")
       assertFileExists(t, ".ddx.yml")

       // Step 2: Apply template
       runCommand(t, "ddx", "apply", "templates/nextjs")
       assertFileExists(t, "package.json")

       // Step 3: Update from upstream
       runCommand(t, "ddx", "update")

       // Step 4: Make changes and contribute
       modifyFile(t, "templates/custom/my-template.md")
       runCommand(t, "ddx", "contribute", "--message", "Add custom template")
   }
   ```

2. **Test cross-platform compatibility**
   ```go
   func TestE2E_CrossPlatform(t *testing.T) {
       if runtime.GOOS == "windows" {
           t.Run("Windows paths", testWindowsPaths)
       }
       if runtime.GOOS == "darwin" {
           t.Run("macOS specifics", testMacOSSpecifics)
       }
       if runtime.GOOS == "linux" {
           t.Run("Linux permissions", testLinuxPermissions)
       }
   }
   ```

### 2.5 Security Test Development

#### Procedure: Writing Security Tests

1. **Test input validation**
   ```go
   func TestSecurity_PathTraversal(t *testing.T) {
       maliciousInputs := []string{
           "../../../etc/passwd",
           "..\\..\\..\\windows\\system32",
           "templates/../../../sensitive",
       }

       for _, input := range maliciousInputs {
           cmd := exec.Command("ddx", "apply", input)
           output, _ := cmd.CombinedOutput()

           // Should reject malicious input
           assert.Contains(t, string(output), "invalid path")
           assert.NotEqual(t, 0, cmd.ProcessState.ExitCode())
       }
   }
   ```

2. **Test secret detection**
   ```go
   func TestSecurity_SecretDetection(t *testing.T) {
       // Create file with secret
       content := `
           API_KEY=sk_live_abcd1234efgh5678
           password: super_secret_123
       `

       tempFile := writeTestFile(t, content)

       // Should detect and prevent
       cmd := exec.Command("ddx", "contribute", tempFile)
       output, _ := cmd.CombinedOutput()

       assert.Contains(t, string(output), "secret detected")
   }
   ```

## 3. Test Execution Procedures

### 3.1 Local Test Execution

#### Procedure: Running Tests Locally

1. **Run all tests**
   ```bash
   # Run all tests with coverage
   go test ./... -cover -v

   # Generate coverage report
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```

2. **Run specific test types**
   ```bash
   # Unit tests only
   go test ./internal/... -short

   # Integration tests
   go test ./cmd/... -run Integration

   # E2E tests
   go test ./tests/e2e/... -tags=e2e
   ```

3. **Run with race detection**
   ```bash
   go test ./... -race
   ```

4. **Benchmark tests**
   ```bash
   go test -bench=. -benchmem ./...
   ```

### 3.2 CI/CD Test Execution

#### Procedure: Automated Test Execution

1. **Setup GitHub Actions workflow**
   ```yaml
   name: Test
   on: [push, pull_request]

   jobs:
     test:
       runs-on: ${{ matrix.os }}
       strategy:
         matrix:
           os: [ubuntu-latest, macos-latest, windows-latest]
           go: [1.21, 1.22]

       steps:
         - uses: actions/checkout@v4
         - uses: actions/setup-go@v5
           with:
             go-version: ${{ matrix.go }}

         - name: Run tests
           run: |
             go test ./... -v -cover

         - name: Upload coverage
           uses: codecov/codecov-action@v3
   ```

2. **Monitor test results**
   - Check GitHub Actions tab
   - Review test output logs
   - Analyze coverage reports
   - Track test duration trends

### 3.3 Test Data Management

#### Procedure: Managing Test Fixtures

1. **Create test fixtures**
   ```bash
   mkdir -p testdata/configs/valid
   mkdir -p testdata/configs/invalid
   mkdir -p testdata/templates/simple
   ```

2. **Generate test data**
   ```go
   func generateTestConfig(t *testing.T) string {
       config := Config{
           Version: "1.0.0",
           Repository: Repository{
               URL:    "https://github.com/test/repo",
               Branch: "main",
           },
       }

       data, err := yaml.Marshal(config)
       require.NoError(t, err)

       path := filepath.Join(t.TempDir(), ".ddx.yml")
       require.NoError(t, os.WriteFile(path, data, 0644))

       return path
   }
   ```

3. **Clean up test data**
   ```go
   func TestWithCleanup(t *testing.T) {
       tempDir := t.TempDir() // Automatically cleaned up

       // Or manual cleanup
       defer func() {
           os.RemoveAll("test-output")
       }()
   }
   ```

## 4. Test Maintenance Procedures

### 4.1 Updating Tests

#### Procedure: Maintaining Test Suite

1. **Review test failures**
   ```bash
   # Identify failing tests
   go test ./... | grep FAIL

   # Run specific failing test
   go test -run TestName ./package -v
   ```

2. **Update test expectations**
   - Review requirement changes
   - Update test assertions
   - Add new test cases
   - Remove obsolete tests

3. **Refactor test code**
   ```go
   // Extract common setup
   func setupTest(t *testing.T) (*Config, func()) {
       // Setup code
       cleanup := func() {
           // Cleanup code
       }
       return config, cleanup
   }
   ```

### 4.2 Test Review Process

#### Procedure: Code Review for Tests

1. **Pre-review checklist**
   - [ ] Tests follow naming conventions
   - [ ] Tests have clear descriptions
   - [ ] Both happy and error paths tested
   - [ ] No hardcoded values
   - [ ] Tests are deterministic
   - [ ] Proper cleanup implemented

2. **Review criteria**
   - Test coverage adequate
   - Assertions meaningful
   - Test data appropriate
   - Performance acceptable
   - No test interdependencies

3. **Post-review actions**
   - Address review comments
   - Update test documentation
   - Re-run test suite
   - Update coverage reports

## 5. Troubleshooting Guide

### 5.1 Common Issues and Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| Flaky tests | Race conditions | Add synchronization, use -race flag |
| Slow tests | Large test data | Use smaller datasets, parallelize |
| Coverage gaps | Missing test files | Add tests for uncovered code |
| Environment issues | Missing dependencies | Document requirements, use containers |
| Path issues | OS differences | Use filepath package, not hardcoded paths |

### 5.2 Debugging Failed Tests

1. **Increase verbosity**
   ```bash
   go test -v -run TestName
   ```

2. **Add debug output**
   ```go
   t.Logf("Debug: value = %+v", value)
   ```

3. **Use debugger**
   ```bash
   dlv test ./package -- -test.run TestName
   ```

4. **Isolate test**
   ```bash
   go test -run "^TestName$" ./package
   ```

## 6. Quality Checklists

### 6.1 Unit Test Checklist

- [ ] Test covers all public methods
- [ ] Error conditions tested
- [ ] Boundary values tested
- [ ] Nil/empty inputs handled
- [ ] Concurrent access tested (if applicable)
- [ ] Mock objects properly configured
- [ ] Test data is minimal
- [ ] No external dependencies

### 6.2 Integration Test Checklist

- [ ] Real components used (minimal mocking)
- [ ] Environment properly isolated
- [ ] Database transactions rolled back
- [ ] File system changes cleaned up
- [ ] Network calls mocked or recorded
- [ ] Error propagation tested
- [ ] Performance within limits

### 6.3 E2E Test Checklist

- [ ] Complete user journey tested
- [ ] All platforms tested
- [ ] Performance benchmarked
- [ ] Error recovery tested
- [ ] Data persistence verified
- [ ] UI feedback validated
- [ ] Documentation matches behavior

## 7. Performance Testing

### 7.1 Benchmark Development

```go
func BenchmarkConfigLoad(b *testing.B) {
    configFile := "testdata/config.yml"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        LoadConfig(configFile)
    }
}
```

### 7.2 Load Testing

```bash
# Using vegeta for load testing
echo "GET http://localhost:8080/api/health" | \
    vegeta attack -duration=30s -rate=100 | \
    vegeta report
```

## 8. Continuous Improvement

### 8.1 Metrics Collection

- Test execution time
- Coverage percentage
- Failure rate
- Defect detection rate
- Test maintenance effort

### 8.2 Regular Reviews

- Weekly: Test failures and fixes
- Monthly: Coverage analysis
- Quarterly: Test strategy review
- Annually: Tool and process evaluation

## Appendix A: Test Templates

### Unit Test Template
```go
package package_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFunction_Scenario(t *testing.T) {
    // Arrange
    input := "test"
    expected := "result"

    // Act
    actual := Function(input)

    // Assert
    assert.Equal(t, expected, actual)
}
```

### Table-Driven Test Template
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "result", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            actual, err := Function(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, actual)
            }
        })
    }
}
```

## Appendix B: Useful Commands

```bash
# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=c.out ./...
go tool cover -html=c.out

# Run specific test
go test -run TestName ./package

# Run tests in parallel
go test -parallel 4 ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run with verbose output
go test -v ./...

# Run short tests only
go test -short ./...
```