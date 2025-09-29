# DDx Test Plan

**Document Type**: Test Plan
**Standard**: IEEE 829-2008 compliant
**Version**: 1.0.0
**Status**: Draft
**Created**: 2025-01-15
**Updated**: 2025-01-18

## 1. Test Plan Identifier

**Plan ID**: TP-DDX-001
**Project**: DDx CLI Toolkit
**Phase**: Test (HELIX Phase 03)

## 2. Introduction

This test plan defines the comprehensive testing strategy for the DDx CLI toolkit, following Specification-Driven Development (SDD) principles and HELIX workflow requirements. All tests are written before implementation to serve as executable specifications.

### 2.1 Objectives

- Define test scope, approach, resources, and schedule
- Establish test coverage targets and quality gates
- Ensure all requirements from Frame and Design phases are validated
- Implement test-first development practices
- Achieve minimum 80% code coverage for critical paths

### 2.2 Scope

**In Scope:**
- Core CLI commands (init, doctor, update, contribute, version)
- Resource commands following noun-verb structure:
  - `prompts` commands (list, show)
  - `templates` commands (list, show, apply)
  - `patterns` commands (list, show, apply)
  - `persona` commands (list, show, bind, load)
  - `mcp` commands (list, show)
  - `workflows` commands (list, show, run)
- Internal packages (config, git, templates, persona)
- Library path resolution (development, config, environment, fallback)
- API contract validation
- Security controls and input validation
- Performance benchmarks
- Cross-platform compatibility

**Out of Scope:**
- Third-party library internals
- Operating system APIs
- Network infrastructure testing

## 3. Test Items

### 3.1 Components Under Test

| Component | Type | Priority | Coverage Target |
|-----------|------|----------|-----------------|
| Core Commands (init, doctor, update) | Integration | P0 | 100% |
| Prompts Commands (list, show) | Integration | P0 | 100% |
| Templates Commands (list, show, apply) | Integration | P0 | 100% |
| Patterns Commands (list, show, apply) | Integration | P0 | 100% |
| Persona Commands (list, show, bind, load) | Integration | P0 | 100% |
| Library Path Resolution | Unit/Integration | P0 | 100% |
| Config Package | Unit | P0 | 90% |
| Git Package | Unit/Integration | P0 | 85% |
| Templates Package | Unit | P0 | 90% |
| Persona Package | Unit | P0 | 90% |
| API Contracts | Contract | P0 | 100% |
| Security Controls | Security | P0 | 100% |
| Workflows | E2E | P1 | Critical paths |
| Error Handling | Unit/Integration | P1 | 80% |

### 3.2 Test Deliverables

1. Test plan document (this document)
2. Test procedures document
3. Test specifications
4. Test suites (unit, integration, contract, e2e, security)
5. Test reports
6. Coverage reports
7. Performance benchmarks

## 4. Testing Approach

### 4.1 Testing Pyramid

Following the testing pyramid principle for optimal test distribution:

```
         /\
        /E2E\       5% - End-to-end tests (critical user journeys)
       /------\
      /Contract\    10% - Contract tests (API compliance)
     /----------\
    /Integration \  25% - Integration tests (component interactions)
   /--------------\
  /     Unit      \ 60% - Unit tests (isolated component logic)
 /________________\
```

### 4.2 Test Types

#### 4.2.1 Unit Tests
**Purpose**: Test individual functions and methods in isolation
**Tools**: Go testing package, testify assertions
**Location**: Adjacent to source files (*_test.go)
**Approach**:
- Table-driven tests for multiple scenarios
- Mock external dependencies
- Test both happy and error paths
- Focus on pure business logic

#### 4.2.2 Integration Tests
**Purpose**: Test component interactions and CLI commands
**Tools**: Go testing, test fixtures
**Location**: cli/cmd/*_test.go
**Approach**:
- Test actual command execution
- Use temporary directories
- Verify file system changes
- Test with real git operations

#### 4.2.3 Contract Tests
**Purpose**: Validate API contracts from design phase
**Tools**: Go testing, schema validation
**Location**: tests/contracts/
**Approach**:
- Test CLI command interfaces
- Validate input/output formats
- Verify exit codes
- Check error messages

#### 4.2.4 End-to-End Tests
**Purpose**: Validate complete user workflows
**Tools**: Go testing, shell scripts
**Location**: tests/e2e/
**Approach**:
- Test full workflows
- Cross-platform validation
- Performance benchmarks
- Real repository interactions

#### 4.2.5 Security Tests
**Purpose**: Validate security controls
**Tools**: Go testing, security scanners
**Location**: tests/security/
**Approach**:
- Input validation testing
- Path traversal prevention
- Secret detection validation
- Injection attack prevention

### 4.3 Noun-Verb Command Testing

#### Testing Strategy for Resource Commands
Each resource command must be tested for:

1. **Command Structure Validation**
   - Correct noun-verb pattern recognition
   - Subcommand routing (e.g., `prompts list` vs `prompts show`)
   - Help text generation for each resource
   - Tab completion hints

2. **Library Path Resolution Testing**
   ```go
   func TestLibraryPathResolution(t *testing.T) {
       tests := []struct {
           name     string
           setup    func() // Set env, config, etc.
           expected string
       }{
           {"flag override", setupFlagOverride, "/custom/path"},
           {"env variable", setupEnvVar, "/env/path"},
           {"config file", setupConfigFile, "./library"},
           {"development mode", setupDevMode, "./library"},
           {"project library", setupProjectLib, ".ddx/library"},
           {"global fallback", setupNoConfig, "~/.ddx/library"},
       }
   }
   ```

3. **Resource-Specific Tests**
   - **Prompts**: Test recursive listing with `--verbose`
   - **Templates**: Test template application with variables
   - **Patterns**: Test pattern matching and filtering
   - **Personas**: Test binding persistence and loading

4. **Migration Testing**
   - Ensure no `getDDxHome()` function calls
   - Verify no hardcoded `~/.ddx` paths
   - Test that old commands return helpful migration messages

### 4.4 Test Data Management

#### Test Fixtures
```
testdata/
├── configs/          # Sample configuration files
│   ├── valid/       # Valid .ddx.yml examples
│   └── invalid/     # Invalid configurations for error testing
├── templates/       # Test templates
│   ├── simple/      # Basic templates
│   └── complex/     # Templates with variables
├── repos/           # Git repository fixtures
│   ├── empty/       # Empty repo
│   └── with-ddx/    # Repo with DDx initialized
└── secrets/         # Files with secrets for detection testing
```

#### Test Data Principles
- Use deterministic test data
- Clean up after tests
- Avoid external dependencies
- Use smallest viable datasets

## 5. Pass/Fail Criteria

### 5.1 Test Case Pass Criteria
- All assertions pass
- No unexpected errors
- Expected output matches
- Side effects verified
- Performance within limits

### 5.2 Test Suite Pass Criteria
- All tests pass
- Coverage targets met
- No flaky tests
- Performance benchmarks pass
- Security scans clean

### 5.3 Release Criteria
- 100% of P0 tests passing
- 95% of P1 tests passing
- Code coverage ≥ 80%
- No critical security issues
- Performance requirements met
- Cross-platform tests passing

## 6. Test Environment

### 6.1 Development Environment
- **OS**: macOS, Linux, Windows
- **Go Version**: 1.21+
- **Git Version**: 2.30+
- **Test Runner**: go test
- **Coverage Tool**: go cover
- **CI/CD**: GitHub Actions

### 6.2 Test Environments

| Environment | Purpose | Configuration |
|------------|---------|---------------|
| Local | Developer testing | Developer machine |
| CI | Automated testing | GitHub Actions runners |
| Cross-platform | Compatibility | Matrix: OS × Go version |
| Performance | Benchmarking | Dedicated resources |

### 6.3 Environment Setup
```bash
# Install dependencies
go mod download

# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test type
make test-unit
make test-integration
make test-e2e
```

## 7. Testing Tasks

### 7.1 Test Development Schedule

| Week | Task | Owner | Status |
|------|------|-------|--------|
| 1 | Test plan & procedures | Test Lead | In Progress |
| 1 | Unit test development | Dev Team | Pending |
| 2 | Integration tests | Dev Team | Pending |
| 2 | Contract tests | QA Team | Pending |
| 3 | E2E tests | QA Team | Pending |
| 3 | Security tests | Security Team | Pending |
| 4 | Performance tests | Performance Team | Pending |
| 4 | Test review & refinement | All | Pending |

### 7.2 Test Execution Schedule

| Phase | Duration | Activities |
|-------|----------|------------|
| Initial | 1 day | Smoke tests, environment validation |
| Unit | 2 days | All unit tests |
| Integration | 2 days | Component integration tests |
| System | 3 days | E2E, contract, security tests |
| Regression | 2 days | Full test suite |
| Performance | 2 days | Load and stress tests |

## 8. Responsibilities

### 8.1 Roles and Responsibilities

| Role | Responsibilities |
|------|------------------|
| Test Lead | Test strategy, plan, coordination |
| Developers | Unit tests, integration tests |
| QA Engineers | E2E tests, test execution |
| Security Team | Security test design and execution |
| DevOps | CI/CD setup, environment management |

### 8.2 Test Review Process
1. Developer writes tests following TDD
2. Peer review via pull request
3. Test lead reviews coverage
4. Security team reviews security tests
5. Automated CI validation

## 9. Risk Analysis

### 9.1 Testing Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Incomplete test coverage | High | Medium | Enforce coverage gates |
| Flaky tests | Medium | High | Fix immediately, quarantine if needed |
| Environment differences | High | Medium | Use containers, matrix testing |
| Test data corruption | Low | Low | Isolated test environments |
| Long test execution | Medium | Medium | Parallel execution, test optimization |

### 9.2 Product Risks

| Risk | Testing Strategy |
|------|------------------|
| Data loss during sync | Extensive backup/restore testing |
| Security vulnerabilities | Security test suite, SAST/DAST |
| Cross-platform issues | Matrix testing on all platforms |
| Performance degradation | Benchmark tests, profiling |
| Breaking changes | Contract tests, versioning tests |

## 10. Test Metrics

### 10.1 Coverage Metrics
- Line coverage: Target 80%, Critical paths 95%
- Branch coverage: Target 75%
- Function coverage: Target 90%
- Package coverage: 100% (all packages have tests)

### 10.2 Quality Metrics
- Test execution time: < 5 minutes for unit, < 15 minutes total
- Defect detection rate: > 90% before production
- Test effectiveness: Defect escape rate < 5%
- Test maintainability: < 20% test code churn per release

### 10.3 Reporting
Weekly test reports including:
- Test execution status
- Coverage trends
- Defect metrics
- Risk assessment
- Blockers and issues

## 11. Test Tools

### 11.1 Testing Framework
```go
// Standard Go testing
import "testing"

// Assertion library
import "github.com/stretchr/testify/assert"

// Mocking
import "github.com/stretchr/testify/mock"

// BDD-style tests (optional)
import "github.com/onsi/ginkgo"
import "github.com/onsi/gomega"
```

### 11.2 CI/CD Integration
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: make test-all
      - run: make coverage-report
```

### 11.3 Additional Tools
- **golangci-lint**: Static analysis
- **go-acc**: Accurate coverage across packages
- **gotestsum**: Better test output
- **go-mutesting**: Mutation testing
- **vegeta**: Load testing

## 12. Approvals

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Test Lead | [Name] | [Signature] | [Date] |
| Development Lead | [Name] | [Signature] | [Date] |
| Product Owner | [Name] | [Signature] | [Date] |
| Security Lead | [Name] | [Signature] | [Date] |

## Appendix A: Test Case Template

```go
func TestComponentName_MethodName_Scenario(t *testing.T) {
    // Arrange
    // Set up test data and dependencies

    // Act
    // Execute the function under test

    // Assert
    // Verify the results
}
```

## Appendix B: Test Naming Conventions

- Unit tests: `Test<Type>_<Method>_<Scenario>`
- Integration tests: `TestIntegration_<Feature>_<Scenario>`
- E2E tests: `TestE2E_<Workflow>_<Scenario>`
- Benchmark tests: `Benchmark<Operation>`

## Appendix C: Coverage Exclusions

Excluded from coverage requirements:
- Generated code
- Test helpers
- Main function
- Panic handlers
- OS-specific code (tested separately)