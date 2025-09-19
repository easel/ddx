# Build Plan Generation Prompt

Create a comprehensive build plan that defines how to systematically implement code to make failing tests pass. This plan takes the failing test suite from the Test phase and creates a roadmap for implementation using Test-Driven Development principles.

## Storage Location

Store the build plan at: `docs/helix/04-build/implementation-plan.md`

## Purpose

The build plan is an implementation roadmap that:
- Prioritizes which tests to make pass first
- Defines the component build sequence
- Plans refactoring milestones
- Establishes code organization patterns
- Creates checkpoints for progress tracking

## Key Principles

### 1. Test-Driven Implementation
No code without a failing test - every line of production code must be written to make a specific test pass.

### 2. Incremental Development
Build the system incrementally, making one test pass at a time, maintaining a working system throughout.

### 3. Red-Green-Refactor Cycle
Follow the TDD cycle religiously:
- Red: Start with a failing test
- Green: Write minimal code to pass
- Refactor: Improve code while keeping tests green

### 4. Continuous Integration
Every test that passes should be committed, ensuring the build is always in a deployable state.

## Plan Components

### Implementation Strategy

#### Build Order Priority
Define which tests to tackle first based on:
- **Dependencies**: Tests that unblock others
- **Risk**: High-risk areas first
- **Value**: Business-critical features
- **Complexity**: Balance easy wins with hard problems
- **Learning**: Areas that inform architecture

#### Component Sequence
Plan the order of component development:
1. Core data models
2. Business logic layer
3. Service interfaces
4. API endpoints
5. UI components
6. Integration points

#### Architecture Decisions
Document key implementation choices:
- Design patterns to use
- Code organization structure
- Dependency management
- Error handling approach
- Performance considerations

### Test-to-Code Mapping

#### Test Analysis
For each failing test, identify:
- What functionality it requires
- Which components need implementation
- Dependencies on other tests
- Estimated implementation effort

#### Implementation Tasks
Break down each test into tasks:
- Minimal code to make test pass
- Refactoring opportunities
- Documentation needs
- Performance optimization points

#### Progress Tracking
Define how to measure progress:
- Tests passing per day
- Coverage increasing
- Code quality metrics
- Technical debt tracking

### Development Workflow

#### Daily Cycle
Structure the daily development flow:
1. Review failing tests
2. Select next test to tackle
3. Implement minimal solution
4. Verify test passes
5. Refactor if needed
6. Commit and push
7. Update progress tracking

#### Code Review Points
Define when to pause for review:
- After each component completion
- Before major refactoring
- At integration boundaries
- When patterns emerge

#### Integration Checkpoints
Plan integration milestones:
- Component integration tests passing
- API contract tests passing
- End-to-end tests passing
- Performance benchmarks met

### Refactoring Strategy

#### Refactoring Triggers
When to refactor:
- Three strikes rule (duplication)
- Performance bottlenecks
- Complexity thresholds
- Test feedback
- Code review feedback

#### Refactoring Patterns
Common refactoring to apply:
- Extract method/class
- Introduce design patterns
- Simplify conditionals
- Remove duplication
- Improve naming

#### Safety Checks
Ensure refactoring safety:
- All tests remain green
- No functionality change
- Performance not degraded
- Code coverage maintained

### Code Organization

#### Project Structure
Define the codebase organization:
- Layer separation (MVC, Clean Architecture, etc.)
- Module boundaries
- Shared code location
- Configuration management
- Documentation placement

#### Coding Standards
Establish consistency rules:
- Naming conventions
- File organization
- Comment standards
- Error handling patterns
- Logging approach

#### Documentation Requirements
Plan documentation creation:
- API documentation
- Code comments
- README files
- Architecture diagrams
- Developer guides

### Quality Assurance

#### Definition of Done
When is a test truly "passing":
- Test passes consistently
- Code is clean and readable
- Documentation is updated
- Code review completed
- Performance acceptable

#### Quality Gates
Checkpoints before proceeding:
- Code coverage targets met
- No critical issues from linting
- Security checks passed
- Performance benchmarks met
- Documentation complete

#### Technical Debt Management
Handle debt strategically:
- Track debt as you create it
- Plan refactoring sprints
- Balance speed with quality
- Document compromises
- Schedule debt payment

### Risk Management

#### Implementation Risks
Identify potential blockers:
- Complex integrations
- Performance challenges
- Unclear requirements
- Technical limitations
- Resource constraints

#### Mitigation Strategies
Plan risk responses:
- Spike solutions for unknowns
- Prototype risky areas first
- Have backup approaches
- Time-box investigations
- Escalation procedures

### Resource Planning

#### Team Allocation
Plan resource distribution:
- Feature team assignments
- Pairing rotations
- Knowledge sharing sessions
- Code review assignments
- Support responsibilities

#### Timeline Estimation
Realistic time planning:
- Tests per day velocity
- Component completion dates
- Integration milestones
- Buffer for unknowns
- Refactoring time

#### Tool Requirements
Development environment needs:
- IDE setup and plugins
- Debugging tools
- Performance profilers
- Database tools
- API testing tools

## Integration with Test Plan

### Using the Test Plan
Reference the test plan for:
- Understanding test structure
- Test priority order
- Coverage requirements
- Test infrastructure
- Success metrics

### Test Plan Alignment
Ensure alignment on:
- Test organization mirrors code organization
- Priority order matches business value
- Coverage targets drive implementation
- Infrastructure supports development

### Feedback Loop
Provide feedback to test plan:
- Tests that need clarification
- Missing test scenarios discovered
- Test infrastructure improvements
- Test organization refinements

## Success Metrics

### Velocity Metrics
Track development speed:
- Tests passing per day
- Story points completed
- Features delivered
- Bugs discovered vs fixed

### Quality Metrics
Monitor code quality:
- Code coverage percentage
- Cyclomatic complexity
- Technical debt ratio
- Code review findings
- Performance metrics

### Progress Indicators
Measure overall progress:
- Percentage of tests passing
- Features complete
- Integration points working
- Documentation coverage

## Quality Checklist

Before starting implementation:
- [ ] All failing tests are understood
- [ ] Implementation priority is clear
- [ ] Architecture decisions are documented
- [ ] Team responsibilities are assigned
- [ ] Development environment is ready
- [ ] CI/CD pipeline is configured
- [ ] Code standards are defined
- [ ] Review process is established

## Common Pitfalls to Avoid

### ❌ Implementing Without Tests
- Bad: Writing code not demanded by tests
- Good: Only code that makes tests pass

### ❌ Over-Engineering
- Bad: Building for future requirements
- Good: Minimal code to pass current tests

### ❌ Skipping Refactoring
- Bad: Accumulating technical debt
- Good: Regular refactoring when tests are green

### ❌ Ignoring Test Feedback
- Bad: Fighting against test design
- Good: Let tests guide architecture

### ❌ Big Bang Integration
- Bad: Integrate everything at once
- Good: Continuous integration throughout

## Next Phase: Deploy

The build plan enables deployment by ensuring:
1. All tests are passing
2. Code quality standards are met
3. Documentation is complete
4. Performance targets achieved
5. System is production-ready

Remember: The build phase is about making tests pass systematically, not about writing tests (they already exist and are failing).