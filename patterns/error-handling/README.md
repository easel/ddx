# Error Handling Patterns

This directory contains common error handling patterns for different languages and frameworks.

## General Principles

1. **Fail Fast**: Detect errors early and handle them appropriately
2. **Be Specific**: Provide meaningful error messages with context
3. **Log Appropriately**: Log errors at the right level with sufficient detail
4. **Graceful Degradation**: Handle errors in a way that maintains user experience
5. **Don't Hide Errors**: Ensure errors are visible to developers and operators

## Language-Specific Patterns

### Go
- Use explicit error returns: `func doSomething() (result, error)`
- Check errors immediately: `if err != nil { return err }`
- Wrap errors with context: `fmt.Errorf("failed to process %s: %w", item, err)`
- Use custom error types for specific cases

### JavaScript/TypeScript
- Use try/catch for synchronous operations
- Handle Promise rejections with `.catch()` or try/catch with async/await
- Create custom Error classes for different error types
- Use error boundaries in React for UI error handling

### Python
- Use specific exception types rather than generic Exception
- Follow the EAFP principle: "Easier to Ask for Forgiveness than Permission"
- Use context managers for resource cleanup
- Log exceptions with full stack traces

### Rust
- Use `Result<T, E>` for fallible operations
- Use `Option<T>` for values that might not exist
- Implement custom error types with `thiserror` or similar
- Use `?` operator for error propagation

## Common Anti-Patterns to Avoid

- **Silent Failures**: Catching exceptions without proper handling
- **Generic Error Messages**: "Something went wrong" tells users nothing
- **Exception Swallowing**: Catching errors without logging or handling
- **Over-Catching**: Using broad exception handlers that hide specific issues
- **Error Codes**: Using magic numbers instead of meaningful error types

## Best Practices

1. **Create Error Hierarchies**: Organize errors by type and severity
2. **Include Context**: Add relevant information to error messages
3. **Use Structured Logging**: Include error details in structured format
4. **Test Error Paths**: Write tests for error conditions
5. **Document Error Behavior**: Specify what errors functions can return
6. **Monitor Errors**: Set up alerting for critical error conditions

## Error Handling Checklist

- [ ] Errors are caught at appropriate levels
- [ ] Error messages provide sufficient context
- [ ] Critical errors are logged with full details
- [ ] User-facing errors are friendly and actionable
- [ ] Error recovery strategies are implemented where possible
- [ ] Error paths are tested
- [ ] Monitoring and alerting is in place for critical errors