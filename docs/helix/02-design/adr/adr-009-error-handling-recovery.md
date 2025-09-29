---
tags: [adr, architecture, error-handling, recovery, resilience, ddx]
template: false
version: 1.0.0
---

# ADR-009: Error Handling and Recovery Strategy

**Date**: 2025-01-14
**Status**: Proposed
**Deciders**: DDX Development Team
**Technical Story**: Define comprehensive error handling, recovery mechanisms, and resilience patterns for DDX operations

## Context

### Problem Statement
DDX needs robust error handling that:
- Provides clear, actionable error messages to users
- Enables graceful recovery from failures
- Supports rollback of partial operations
- Maintains system consistency during failures
- Facilitates debugging and troubleshooting
- Prevents data corruption and loss
- Handles network and filesystem failures gracefully

### Forces at Play
- **Clarity**: Error messages must be understandable and actionable
- **Recovery**: System should recover gracefully from failures
- **Consistency**: Partial operations must not leave inconsistent state
- **Debugging**: Errors must provide sufficient context for diagnosis
- **Performance**: Error handling shouldn't degrade performance
- **User Experience**: Errors shouldn't frustrate or confuse users
- **Automation**: Support automated error recovery where possible

### Constraints
- Must work across different platforms
- Cannot lose user data during failures
- Should not expose sensitive information in errors
- Must handle interrupted operations gracefully
- Need to support offline error scenarios
- Must integrate with CI/CD error reporting

## Decision

### Chosen Approach
Implement a multi-layer error handling strategy:
1. **Structured error types** with error codes and categories
2. **Transaction-based operations** with rollback capability
3. **Retry mechanisms** with exponential backoff
4. **Recovery suggestions** with automated fixes where possible
5. **Comprehensive logging** with adjustable verbosity

### Error Classification System

```go
type ErrorCategory string

const (
    // User errors (1xxx)
    ErrUserInput      ErrorCategory = "USER_INPUT"       // 1001-1099
    ErrConfiguration  ErrorCategory = "CONFIGURATION"    // 1100-1199
    ErrPermission     ErrorCategory = "PERMISSION"       // 1200-1299

    // System errors (2xxx)
    ErrFilesystem     ErrorCategory = "FILESYSTEM"       // 2001-2099
    ErrNetwork        ErrorCategory = "NETWORK"          // 2100-2199
    ErrDependency     ErrorCategory = "DEPENDENCY"       // 2200-2299

    // Operation errors (3xxx)
    ErrValidation     ErrorCategory = "VALIDATION"       // 3001-3099
    ErrConflict       ErrorCategory = "CONFLICT"         // 3100-3199
    ErrTimeout        ErrorCategory = "TIMEOUT"          // 3200-3299

    // Internal errors (4xxx)
    ErrInternal       ErrorCategory = "INTERNAL"         // 4001-4099
    ErrPanic          ErrorCategory = "PANIC"            // 4100-4199
    ErrCorruption     ErrorCategory = "CORRUPTION"       // 4200-4299
)

type DDXError struct {
    Code        int                    // Error code
    Category    ErrorCategory          // Error category
    Message     string                 // User-facing message
    Details     string                 // Technical details
    Suggestion  string                 // Recovery suggestion
    Context     map[string]interface{} // Additional context
    Wrapped     error                  // Original error
    Timestamp   time.Time              // When error occurred
    Recoverable bool                   // Can be recovered
}
```

### Error Message Format

```
Error [1001]: Invalid project configuration

The .ddx.yml file contains invalid syntax at line 15, column 8.

Details:
  File: /path/to/.ddx.yml
  Issue: Unexpected indentation
  Expected: 2 spaces
  Found: 3 spaces

Suggestion:
  Fix the indentation at line 15 to use 2 spaces, or run:
  $ ddx validate --fix .ddx.yml

For more information, see: https://ddx.dev/errors/1001
```

### Transaction and Rollback System

```go
type Transaction struct {
    ID          string
    Operations  []Operation
    Checkpoints []Checkpoint
    StartTime   time.Time
    State       TransactionState
}

type Operation struct {
    Name     string
    Execute  func() error
    Rollback func() error
    Verify   func() error
}

type Checkpoint struct {
    Name      string
    Timestamp time.Time
    State     map[string]interface{}
    CanResume bool
}

// Usage example
tx := NewTransaction("apply-template")
tx.AddOperation("backup", backupFiles, restoreFiles)
tx.AddOperation("apply", applyTemplate, revertTemplate)
tx.AddOperation("validate", validateResult, nil)

if err := tx.Execute(); err != nil {
    if rollbackErr := tx.Rollback(); rollbackErr != nil {
        // Catastrophic failure, manual intervention needed
        return CatastrophicError(err, rollbackErr)
    }
    return TransactionError(err)
}
```

### Retry Strategy

```go
type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    Multiplier      float64
    RetryableErrors []ErrorCategory
}

var DefaultRetryConfig = RetryConfig{
    MaxAttempts:  3,
    InitialDelay: 1 * time.Second,
    MaxDelay:     30 * time.Second,
    Multiplier:   2.0,
    RetryableErrors: []ErrorCategory{
        ErrNetwork,
        ErrTimeout,
        ErrFilesystem, // Only specific filesystem errors
    },
}

func WithRetry(operation func() error, config RetryConfig) error {
    var lastErr error
    delay := config.InitialDelay

    for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
        if err := operation(); err == nil {
            return nil
        } else if !isRetryable(err, config) {
            return err
        } else {
            lastErr = err
            if attempt < config.MaxAttempts {
                time.Sleep(delay)
                delay = time.Duration(float64(delay) * config.Multiplier)
                if delay > config.MaxDelay {
                    delay = config.MaxDelay
                }
            }
        }
    }
    return RetryExhaustedError(lastErr, config.MaxAttempts)
}
```

### Recovery Mechanisms

#### Automatic Recovery
```yaml
recovery:
  auto_recovery:
    enabled: true
    strategies:
      - corruption_detection:
          action: "rebuild_index"
          conditions: ["index_checksum_mismatch"]
      - network_failure:
          action: "use_cache"
          conditions: ["network_timeout", "dns_failure"]
      - permission_error:
          action: "request_elevation"
          conditions: ["permission_denied"]
      - lock_conflict:
          action: "wait_and_retry"
          conditions: ["resource_locked"]
```

#### Manual Recovery Commands
```bash
# Recover from corrupted state
ddx doctor --doctor
ddx doctor --repair

# Recover from interrupted operation
ddx recover --transaction-id abc123
ddx recover --from-checkpoint

# Clean up after failure
ddx clean --orphaned-files
ddx clean --invalid-state

# Reset to known good state
ddx reset --soft  # Keep user data
ddx reset --hard  # Full reset
```

### Error Context and Debugging

```go
type ErrorContext struct {
    // Execution context
    Command     string
    Arguments   []string
    WorkingDir  string
    Environment map[string]string

    // System context
    Platform    string
    GoVersion   string
    DDXVersion  string
    GitVersion  string

    // Operation context
    Transaction string
    Operation   string
    Phase       string

    // Debug information
    StackTrace  string
    LogFile     string
    ConfigDump  string
}

// Rich error reporting
func (e *DDXError) Report() string {
    if verbose {
        return e.DetailedReport()
    }
    return e.UserFriendlyReport()
}
```

### Rationale
- **Structured Errors**: Enable programmatic handling and consistent formatting
- **Transactions**: Ensure consistency and enable rollback
- **Retry Logic**: Handle transient failures automatically
- **Recovery Tools**: Provide both automatic and manual recovery options
- **Rich Context**: Include sufficient information for debugging

## Alternatives Considered

### Option 1: Simple Error Strings
**Description**: Use basic error strings like standard Go errors

**Pros**:
- Simple implementation
- Lightweight
- Familiar to Go developers
- No additional complexity

**Cons**:
- No structure for parsing
- Inconsistent messages
- No error codes
- Poor internationalization
- Limited debugging info

**Why rejected**: Insufficient for complex error scenarios and poor UX

### Option 2: Panic/Recover Everywhere
**Description**: Use Go's panic/recover for all error handling

**Pros**:
- Automatic stack unwinding
- Centralized error handling
- Simple error propagation
- Clean code flow

**Cons**:
- Not idiomatic Go
- Performance overhead
- Difficult to handle gracefully
- Poor for expected errors
- Hard to test

**Why rejected**: Anti-pattern in Go and poor for expected errors

### Option 3: Error Channels
**Description**: Use channels for async error handling

**Pros**:
- Async error handling
- Decoupled error processing
- Good for concurrent operations
- Flexible routing

**Cons**:
- Complex implementation
- Timing issues
- Difficult debugging
- Memory overhead
- Not suitable for sync operations

**Why rejected**: Overcomplicated for primarily synchronous operations

### Option 4: Global Error Registry
**Description**: Maintain global registry of all errors

**Pros**:
- Centralized error management
- Consistent error codes
- Easy documentation
- Version control

**Cons**:
- Global state issues
- Tight coupling
- Difficult testing
- Concurrency concerns
- Maintenance overhead

**Why rejected**: Global state creates more problems than it solves

### Option 5: Exception-style with Result Types
**Description**: Use Result<T, E> types like Rust

**Pros**:
- Explicit error handling
- Type safety
- No hidden control flow
- Composable

**Cons**:
- Not idiomatic Go
- Verbose code
- Generic limitations in Go
- Learning curve
- Library incompatibility

**Why rejected**: Not idiomatic and poor ecosystem fit

## Consequences

### Positive Consequences
- **Clear Communication**: Users understand what went wrong
- **Graceful Degradation**: System remains usable during failures
- **Data Safety**: Transactions prevent corruption
- **Debuggability**: Rich context aids troubleshooting
- **Automation**: Many errors self-recover
- **Consistency**: Standardized error handling

### Negative Consequences
- **Complexity**: More code for error handling
- **Performance**: Slight overhead for transactions
- **Storage**: Error logs consume disk space
- **Maintenance**: Error catalog needs updates
- **Testing**: More test scenarios needed

### Neutral Consequences
- **Verbosity**: More detailed error messages
- **Learning Curve**: Users learn error codes
- **Documentation**: Requires error documentation
- **Monitoring**: Need error tracking systems

## Implementation

### Required Changes
1. Define error type hierarchy
2. Implement transaction system
3. Build retry mechanisms
4. Create recovery commands
5. Add error context collection
6. Implement rollback handlers
7. Build error documentation
8. Create debugging tools

### Error Handling Patterns

#### Pattern 1: Wrapped Errors
```go
if err := operation(); err != nil {
    return WrapError(err, ErrFilesystem,
        "Failed to apply template",
        map[string]interface{}{
            "template": templateName,
            "target": targetPath,
        })
}
```

#### Pattern 2: Transaction Pattern
```go
tx := BeginTransaction("complex-operation")
defer tx.Cleanup()

tx.Checkpoint("pre-modification")
if err := modifyFiles(); err != nil {
    return tx.RollbackTo("pre-modification", err)
}

tx.Checkpoint("post-modification")
if err := validate(); err != nil {
    return tx.RollbackTo("pre-modification", err)
}

return tx.Commit()
```

#### Pattern 3: Recovery Pattern
```go
err := operation()
if IsRecoverable(err) {
    if recovered := AttemptRecovery(err); recovered == nil {
        return nil
    }
}
return err
```

### Success Metrics
- **Error Clarity**: > 90% of errors understood without documentation
- **Recovery Rate**: > 70% of recoverable errors auto-recover
- **MTTR**: < 5 minutes mean time to recovery
- **Data Loss**: Zero data loss from DDX errors
- **User Satisfaction**: > 85% find error messages helpful

## Compliance

### Security Requirements
- No sensitive data in error messages
- Sanitize user input in errors
- No system information leakage
- Secure error log storage
- Rate limit error endpoints

### Performance Requirements
- Error handling < 10ms overhead
- Transaction rollback < 1s
- Retry delays configurable
- Logging asynchronous
- Memory bounded for error storage

### Regulatory Requirements
- GDPR compliance for error logs
- No PII in error messages
- Audit trail for critical errors
- Error retention policies

## Monitoring and Review

### Key Indicators to Watch
- Error frequency by category
- Recovery success rates
- Transaction rollback frequency
- Mean time to recovery
- User error report rates
- Critical error trends

### Review Date
Q2 2025 - After production usage patterns emerge

### Review Triggers
- Critical data loss incident
- Recovery rate drops below 60%
- User complaints about errors
- Performance degradation
- New error patterns emerge

## Related Decisions

### Dependencies
- ADR-003: Go implementation affects error handling
- ADR-005: Configuration errors need handling
- ADR-007: Template processing errors
- ADR-008: Contribution validation errors

### Influenced By
- Go error handling idioms
- Transaction patterns from databases
- Retry patterns from distributed systems
- UX research on error messages

### Influences
- CLI output formatting
- Logging strategy
- Monitoring implementation
- Testing strategies
- Documentation requirements

## References

### Documentation
- [Go Error Handling](https://blog.golang.org/error-handling-and-go)
- [Error Handling Best Practices](https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html)
- [Transaction Patterns](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [Retry Patterns](https://docs.microsoft.com/en-us/azure/architecture/patterns/retry)

### External Resources
- [Google SRE Error Budgets](https://sre.google/sre-book/error-budget/)
- [Error Message Guidelines](https://www.nngroup.com/articles/error-message-guidelines/)
- [Resilience Patterns](https://github.com/App-vNext/Polly/wiki/Resilience-policies)

### Discussion History
- Error handling philosophy discussion
- Transaction system design review
- Recovery mechanism evaluation
- User feedback on error messages

## Notes

The error handling system follows DDX's medical metaphor - like medical diagnosis and treatment, errors are classified (diagnosis), treated (recovery), and prevented (validation). The transaction system mirrors medical procedures with checkpoints and rollback capabilities ensuring patient (system) safety.

Key insight: By treating errors as first-class citizens with structure, context, and recovery paths, we transform frustrating failures into learning opportunities and maintain user trust.

Implementation tip: Start with comprehensive error collection and gradually add automated recovery based on actual failure patterns. Over-engineering recovery for hypothetical failures wastes effort.

The staged approach mirrors medical triage - immediate assessment (error classification), stabilization (rollback/recovery), and treatment (resolution) with clear escalation paths for critical issues.

---

**Last Updated**: 2025-01-14
**Next Review**: 2025-04-14