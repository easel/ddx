# Test Plan Generation Prompt

Create a comprehensive test plan that defines the strategy, structure, and approach for writing tests BEFORE any implementation exists. This plan guides the creation of failing tests that will drive development in the Build phase.

## Storage Location

Store the test plan at: `docs/testing/test-plan.md`

## Purpose

The test plan is a strategic document that:
- Defines the testing approach and methodology
- Establishes test organization and standards
- Sets coverage targets and priorities
- Plans test infrastructure and tooling
- Creates a roadmap for test creation

## Key Principles

### 1. Tests Define Behavior
Tests are the executable specification - they define what the system MUST do before any code exists.

### 2. Comprehensive Coverage Strategy
Plan for all levels of testing from the start, not as an afterthought.

### 3. Test Independence
Each test should be independent and not rely on other tests or execution order.

### 4. Clear Traceability
Every test must trace back to a requirement or user story.

## Plan Components

### Testing Strategy

#### Test Levels and Scope
Define what each level of testing will cover:
- **Contract Tests**: External API boundaries
- **Integration Tests**: Component interactions
- **Unit Tests**: Business logic and algorithms
- **E2E Tests**: Critical user journeys

#### Test Framework Selection
Choose appropriate frameworks based on:
- Language and technology stack
- Team expertise
- CI/CD compatibility
- Reporting capabilities
- Community support

#### Coverage Targets
Set specific, measurable coverage goals:
- Overall line coverage target
- Critical path coverage (must be 100%)
- Branch coverage requirements
- Error handling coverage

### Test Organization

#### Directory Structure
Plan how tests will be organized:
- Separation by test type
- Naming conventions
- File structure patterns
- Shared utilities location

#### Test Naming Standards
Establish clear naming patterns:
- Descriptive test names
- Consistent format across team
- Clear indication of what's being tested
- Expected behavior in the name

#### Test Data Management
Plan how test data will be handled:
- Fixtures for static data
- Factories for dynamic data
- Database seeding strategies
- External service mocking

### Implementation Roadmap

#### Priority Order
Define which tests to write first:
1. Critical path tests (P0 features)
2. API contract tests
3. Core business logic tests
4. Integration points
5. Edge cases and error handling

#### Test Dependencies
Identify dependencies between test types:
- Which tests block others
- Shared test infrastructure needs
- Mock/stub requirements
- Test environment setup

#### Timeline Estimates
Provide realistic timeframes:
- Time to write each test category
- Infrastructure setup time
- Test execution time targets
- Maintenance overhead

### Infrastructure Requirements

#### Test Environment
Define what's needed for test execution:
- Local development setup
- CI/CD pipeline configuration
- Test database requirements
- External service simulators
- Performance testing tools

#### Tooling and Utilities
Identify required tools:
- Test runners
- Coverage reporters
- Assertion libraries
- Mocking frameworks
- Test data generators

#### Continuous Integration
Plan CI/CD integration:
- Test execution triggers
- Parallel test execution
- Failure notifications
- Coverage reporting
- Performance benchmarks

### Risk Mitigation

#### Testing Risks
Identify potential issues:
- Flaky test patterns to avoid
- Performance bottlenecks
- Environment dependencies
- Data consistency issues

#### Mitigation Strategies
Plan how to address risks:
- Test stability practices
- Retry mechanisms
- Environment isolation
- Data cleanup procedures

### Success Metrics

#### Quality Metrics
Define how to measure test quality:
- Test execution time
- Failure rate
- Flakiness index
- Maintenance burden
- Bug escape rate

#### Coverage Metrics
Track coverage comprehensively:
- Code coverage percentages
- Requirements coverage
- User story coverage
- Edge case coverage

## Integration with Build Phase

### Handoff to Build Team
The test plan provides build phase with:
- Clear understanding of test structure
- Priority order for implementation
- Coverage requirements to meet
- Test infrastructure needs

### Traceability Matrix
Create mapping between:
- Requirements → Tests
- Tests → Implementation areas
- User stories → Test scenarios
- Risks → Test cases

## Quality Checklist

Before completing the test plan:
- [ ] All test levels are defined
- [ ] Framework choices are justified
- [ ] Coverage targets are specific
- [ ] Organization structure is clear
- [ ] Priority order is established
- [ ] Infrastructure needs are identified
- [ ] Risks are addressed
- [ ] Success metrics are measurable
- [ ] Integration with build phase is planned

## Common Pitfalls to Avoid

### ❌ Vague Coverage Goals
- Bad: "Good test coverage"
- Good: "80% line coverage, 100% critical path coverage"

### ❌ Test After Development
- Bad: "We'll add tests after building"
- Good: "Tests are written first to drive development"

### ❌ Unclear Organization
- Bad: Tests scattered without structure
- Good: Clear separation by type and purpose

### ❌ Missing Infrastructure Planning
- Bad: Figure out tools as we go
- Good: All tools and environments defined upfront

### ❌ No Priority Order
- Bad: Write tests randomly
- Good: Strategic order based on risk and value

## Next Phase: Build

The test plan enables the Build phase by providing:
1. Clear test structure to work against
2. Priority order for implementation
3. Coverage targets to achieve
4. Infrastructure already planned
5. Success metrics defined

Remember: This plan guides the creation of FAILING tests that define system behavior before any implementation exists.