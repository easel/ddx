# Test Procedures Generation Prompt

Create comprehensive test procedures that define how tests will be written, executed, and maintained throughout the project lifecycle. These procedures establish standards and practices for the test implementation phase.

## Storage Location

Store the test procedures at: `docs/testing/test-procedures.md`

## Purpose

Test procedures provide the implementation team with:
- Step-by-step guidance for writing each type of test
- Standards for test structure and organization
- Execution and validation procedures
- Quality checklists and review criteria
- Troubleshooting guides for common issues

## Key Requirements

### 1. Test Writing Procedures

Define specific procedures for each test type:

#### Unit Test Procedures
- Setup and teardown patterns
- Mocking and stubbing guidelines
- Assertion best practices
- Test data management
- Coverage requirements

#### Integration Test Procedures
- Service integration patterns
- Database test data management
- External API simulation
- Transaction handling
- Cleanup procedures

#### End-to-End Test Procedures
- User flow simulation
- Browser automation setup
- Mobile testing approach
- Data persistence verification
- Performance monitoring

#### Performance Test Procedures
- Load generation setup
- Metrics collection
- Baseline establishment
- Results analysis
- Bottleneck identification

#### Security Test Procedures
- Vulnerability scanning setup
- Authentication testing
- Authorization verification
- Input validation testing
- Security regression testing

### 2. Execution Procedures

#### Local Development
- Environment setup steps
- Pre-test configuration
- Test runner commands
- Debug mode execution
- Partial test execution

#### CI/CD Pipeline
- Automated execution triggers
- Parallel execution setup
- Failure handling
- Retry mechanisms
- Results reporting

### 3. Validation Procedures

#### Test Quality Checks
- Code review checklist
- Test independence verification
- Determinism validation
- Performance impact assessment
- Maintenance burden evaluation

#### Coverage Validation
- Coverage tool configuration
- Target verification
- Gap analysis procedures
- Exemption handling
- Reporting standards

### 4. Maintenance Procedures

#### Test Updates
- When to update tests
- Refactoring guidelines
- Deprecation process
- Version compatibility
- Documentation updates

#### Failure Investigation
- Root cause analysis steps
- Flaky test identification
- Environment issue diagnosis
- Fix verification process
- Regression prevention

### 5. Quality Standards

#### Naming Conventions
- Test file naming patterns
- Test method naming
- Test data naming
- Fixture naming
- Report naming

#### Code Organization
- Test structure patterns
- Helper method organization
- Shared utility location
- Configuration management
- Documentation placement

#### Documentation Standards
- Test purpose description
- Precondition documentation
- Expected behavior notes
- Failure reason tracking
- Maintenance history

## Implementation Guidelines

### Test Independence
Each test must:
- Run in isolation
- Not depend on execution order
- Clean up after itself
- Use fresh test data
- Reset shared state

### Performance Considerations
- Keep tests fast (unit < 100ms)
- Minimize I/O operations
- Use in-memory databases where possible
- Parallelize when appropriate
- Profile slow tests regularly

### Reliability Requirements
- No random failures
- Consistent across environments
- Deterministic outcomes
- Clear failure messages
- Reproducible issues

## Troubleshooting Guide

### Common Issues and Solutions

#### Flaky Tests
- **Symptom**: Tests pass/fail randomly
- **Diagnosis**: Check for timing issues, shared state, external dependencies
- **Solution**: Add proper waits, isolate state, mock external services

#### Slow Tests
- **Symptom**: Test suite takes too long
- **Diagnosis**: Profile test execution, identify bottlenecks
- **Solution**: Optimize queries, use test doubles, parallelize execution

#### Environment Issues
- **Symptom**: Tests fail in CI but pass locally
- **Diagnosis**: Compare environments, check dependencies, verify data
- **Solution**: Align environments, document requirements, use containers

#### Coverage Gaps
- **Symptom**: Coverage below targets
- **Diagnosis**: Identify untested code paths
- **Solution**: Add missing tests, refactor untestable code, update targets

## Quality Checklist

Before marking procedures complete, ensure:
- [ ] All test types have detailed procedures
- [ ] Execution steps are clear and complete
- [ ] Validation criteria are measurable
- [ ] Maintenance processes are defined
- [ ] Quality standards are specific
- [ ] Troubleshooting covers common issues
- [ ] Examples demonstrate best practices
- [ ] Tools and frameworks are specified
- [ ] Integration with CI/CD is documented
- [ ] Review process is established

## Integration with Build Phase

These procedures enable the Build phase by providing:
1. Clear guidelines for test implementation
2. Standards that ensure consistency
3. Troubleshooting guides for common issues
4. Quality checks for test validation
5. Maintenance procedures for long-term success

Remember: Good test procedures make the difference between tests that help and tests that hinder development.