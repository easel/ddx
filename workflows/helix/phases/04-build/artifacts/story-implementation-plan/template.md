# Implementation Plan: IP-XXX-[story-name]

*Step-by-step build plan for making tests pass*

## Story Reference

**User Story**: [[US-XXX-[story-name]]]
**Technical Design**: [[TD-XXX-[story-name]]]
**Test Plan**: [[TP-XXX-[story-name]]]

## Implementation Overview

**Objective**: Make all tests defined in TP-XXX pass using TDD approach

**Approach**: Red-Green-Refactor cycle for each test suite

**Time Estimate**: [X days/hours]

## Test-Driven Development Sequence

### Phase 1: Red (All Tests Failing)

Verify test suite is complete and failing:

```bash
# Run all tests for this story
make test STORY=XXX

# Expected: All tests fail
# - Unit tests: 0/N passing
# - Integration tests: 0/N passing
# - API tests: 0/N passing
```

### Phase 2: Green (Make Tests Pass)

#### Step 1: Unit Tests First

**Target Tests**: `test_[component]_*.py`

**Implementation Order**:
1. **File**: `src/[component]/[module].py`
   - Function: `[function_name]`
   - Purpose: [What it does]
   - Test: `test_[function_name]`

2. **File**: `src/[component]/validator.py`
   - Function: `validate_[input]`
   - Purpose: Input validation
   - Test: `test_validate_[input]`

**Code to Write**:
```python
# Minimal implementation to pass test
def function_name(param):
    # Just enough code to make test pass
    pass
```

#### Step 2: Integration Tests

**Target Tests**: `test_[story]_integration.py`

**Components to Connect**:
1. Component A → Component B
   - File: `src/integration/[connector].py`
   - Method: Wire up components
   - Test: `test_components_integrate`

2. Database Integration
   - File: `src/db/[repository].py`
   - Method: CRUD operations
   - Test: `test_database_operations`

#### Step 3: API Tests

**Target Tests**: `test_[resource]_api.py`

**Endpoints to Implement**:
1. **Endpoint**: `POST /api/v1/[resource]`
   - File: `src/api/[resource]_controller.py`
   - Method: `create_[resource]`
   - Test: `test_create_[resource]`

2. **Endpoint**: `GET /api/v1/[resource]/{id}`
   - File: `src/api/[resource]_controller.py`
   - Method: `get_[resource]`
   - Test: `test_get_[resource]`

### Phase 3: Refactor (Improve Code Quality)

After all tests pass, refactor for:

1. **Code Organization**
   - Extract common functions
   - Remove duplication
   - Improve naming

2. **Performance**
   - Optimize queries
   - Add caching where needed
   - Reduce complexity

3. **Maintainability**
   - Add documentation
   - Improve error messages
   - Enhance logging

## File Creation/Modification Order

### New Files to Create

```yaml
order: 1
files:
  - path: src/models/[model].py
    purpose: Data model for story
    tests: test_models_[model].py

order: 2
files:
  - path: src/services/[service].py
    purpose: Business logic
    tests: test_services_[service].py

order: 3
files:
  - path: src/api/[controller].py
    purpose: API endpoints
    tests: test_api_[controller].py
```

### Existing Files to Modify

```yaml
file: src/config/routes.py
changes:
  - Add new route for [endpoint]
  - Register controller

file: src/db/schema.py
changes:
  - Add new table/column
  - Update indexes
```

## Implementation Checklist

### Pre-Implementation
- [ ] All tests written and failing
- [ ] Technical design reviewed
- [ ] Dependencies identified
- [ ] Development environment ready

### During Implementation

#### Unit Test Implementation
- [ ] `test_[function1]` - passing
- [ ] `test_[function2]` - passing
- [ ] `test_validation` - passing
- [ ] `test_edge_cases` - passing

#### Integration Test Implementation
- [ ] `test_component_integration` - passing
- [ ] `test_database_operations` - passing
- [ ] `test_external_service` - passing

#### API Test Implementation
- [ ] `test_create_resource` - passing
- [ ] `test_get_resource` - passing
- [ ] `test_update_resource` - passing
- [ ] `test_delete_resource` - passing
- [ ] `test_error_handling` - passing

### Post-Implementation
- [ ] All tests passing
- [ ] Code refactored
- [ ] Documentation added
- [ ] Code review completed
- [ ] CI/CD pipeline green

## Incremental Commits

Structure commits for clear history:

```bash
# Commit 1: Model implementation
git add src/models/[model].py tests/test_models_[model].py
git commit -m "feat(US-XXX): implement [model] with passing tests"

# Commit 2: Service layer
git add src/services/[service].py tests/test_services_[service].py
git commit -m "feat(US-XXX): add [service] business logic"

# Commit 3: API endpoints
git add src/api/[controller].py tests/test_api_[controller].py
git commit -m "feat(US-XXX): implement API endpoints for [resource]"

# Commit 4: Integration
git add src/integration/*
git commit -m "feat(US-XXX): wire up components for story"

# Commit 5: Refactoring
git add [refactored files]
git commit -m "refactor(US-XXX): improve code quality and performance"
```

## Dependencies and Blockers

### Dependencies
- [ ] TD-XXX technical design approved
- [ ] TP-XXX test plan complete
- [ ] Required libraries installed
- [ ] Test data available

### Potential Blockers
- Issue: [Potential problem]
  - Mitigation: [How to handle]
- Issue: [Another problem]
  - Mitigation: [Solution]

## Development Environment

### Setup Requirements
```bash
# Install dependencies
pip install -r requirements.txt

# Set up test database
make setup-test-db

# Configure environment
cp .env.example .env.test
```

### Running Tests During Development
```bash
# Run specific test
pytest tests/test_[specific].py::test_[name]

# Run with coverage
pytest --cov=src tests/

# Watch mode for TDD
ptw tests/ -- --failed-first
```

## Definition of Done

Story implementation is complete when:

- [ ] All tests from TP-XXX are passing
- [ ] Code coverage meets minimum (80%)
- [ ] No lint errors or warnings
- [ ] Code follows project style guide
- [ ] Documentation is complete
- [ ] Performance requirements met
- [ ] Security scan passes
- [ ] Code review approved
- [ ] CI/CD pipeline is green

## Time Tracking

Estimated vs Actual:

| Phase | Estimated | Actual | Notes |
|-------|-----------|--------|-------|
| Unit Tests | 4h | [actual] | [notes] |
| Integration | 2h | [actual] | [notes] |
| API | 2h | [actual] | [notes] |
| Refactoring | 2h | [actual] | [notes] |
| **Total** | **10h** | **[actual]** | |

## Notes and Decisions

Document implementation decisions and learnings:

- [Decision 1]: [Rationale]
- [Learning 1]: [What we discovered]
- [TODO]: [Future improvements]

---

*This implementation plan guides the TDD process from red to green to refactored code.*