# Enforcement Mechanisms

## Overview

CDP enforcement mechanisms ensure consistent application of clinical development protocols across all projects and teams. Like medical regulatory bodies that enforce clinical standards, these mechanisms provide both automated and human oversight to maintain code quality and process compliance.

## Automated Enforcement

### Pre-commit Hooks

Pre-commit hooks provide the first line of defense, catching issues before they enter the repository.

**Installation**:
```bash
# Using Lefthook (recommended)
echo 'gem "lefthook"' >> Gemfile
bundle install
lefthook install

# Using pre-commit (Python-based alternative)
pip install pre-commit
pre-commit install
```

**Configuration** (`.lefthook.yml`):
```yaml
pre-commit:
  parallel: true
  
  commands:
    # Secrets detection
    secrets-check:
      glob: "*.{go,js,ts,py,yml,yaml,json}"
      run: gitleaks detect --source . --verbose
      fail_text: "❌ Secrets detected in commit"
    
    # Binary file prevention
    no-binaries:
      glob: "*"
      run: |
        if git diff --cached --name-only | xargs file | grep -q binary; then
          echo "❌ Binary files not allowed in commits"
          exit 1
        fi
    
    # Feature limit check
    feature-limit:
      run: |
        if ! ddx diagnose --features --max=3; then
          echo "❌ Too many concurrent features (max 3)"
          exit 1
        fi
    
    # Go specific checks
    go-format:
      glob: "*.go"
      run: |
        if ! gofmt -l . | wc -l | grep -q "^0$"; then
          echo "❌ Go code not formatted. Run: make fmt"
          exit 1
        fi
    
    go-lint:
      glob: "*.go"
      run: |
        if ! golangci-lint run; then
          echo "❌ Go linting failed"
          exit 1
        fi
    
    go-test:
      glob: "*.go"
      run: |
        if ! go test ./...; then
          echo "❌ Go tests failed"
          exit 1
        fi
    
    # Coverage check
    coverage-check:
      glob: "*.go"
      run: |
        coverage=$(go test -coverprofile=coverage.out ./... | tail -1 | awk '{print $3}' | sed 's/%//')
        if (( $(echo "$coverage < 80" | bc -l) )); then
          echo "❌ Test coverage ${coverage}% below required 80%"
          exit 1
        fi
    
    # Complexity check
    complexity-check:
      glob: "*.go"
      run: |
        if ! gocyclo -over 10 .; then
          echo "❌ Code complexity exceeds limit"
          exit 1
        fi
```

**Custom Hook Scripts**:
```bash
#!/bin/bash
# scripts/hooks/check-commit-message.sh

commit_regex='^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "❌ Invalid commit message format"
    echo "Format: type(scope): description"
    echo "Types: feat, fix, docs, style, refactor, test, chore"
    echo "Example: feat(auth): add user login validation"
    exit 1
fi
```

### CI/CD Pipeline Enforcement

**GitHub Actions Configuration** (`.github/workflows/enforcement.yml`):
```yaml
name: CDP Enforcement

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  enforce-constraints:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Full history for better analysis
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Cache Dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      # Constraint Enforcement
      - name: Check Feature Limits
        run: |
          echo "Checking concurrent feature limits..."
          if ! ddx diagnose --features --max=3; then
            echo "::error::Too many concurrent features in development"
            exit 1
          fi
      
      - name: Complexity Analysis
        run: |
          echo "Analyzing code complexity..."
          go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
          
          # Check cyclomatic complexity
          if ! gocyclo -over 10 .; then
            echo "::error::Code complexity exceeds maximum allowed (10)"
            exit 1
          fi
          
          # Generate complexity report
          gocyclo -avg . > complexity-report.txt
          echo "::notice::Average complexity: $(cat complexity-report.txt)"
      
      - name: Security Scanning
        run: |
          echo "Running security scans..."
          
          # Install security tools
          go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
          
          # Run security scan
          gosec -fmt json -out gosec-report.json ./...
          
          # Check for high/medium severity issues
          high_issues=$(jq '.Stats.found_issues' gosec-report.json)
          if [ "$high_issues" -gt 0 ]; then
            echo "::error::Security vulnerabilities found: $high_issues"
            jq -r '.Issues[] | "- \(.file):\(.line) - \(.details)"' gosec-report.json
            exit 1
          fi
      
      - name: Test Coverage Enforcement
        run: |
          echo "Checking test coverage..."
          
          # Run tests with coverage
          go test -coverprofile=coverage.out ./...
          
          # Extract coverage percentage
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
          
          echo "Current coverage: ${coverage}%"
          
          # Check minimum coverage requirement
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "::error::Test coverage ${coverage}% below required minimum (80%)"
            
            # Show uncovered functions
            echo "Uncovered functions:"
            go tool cover -func=coverage.out | grep -v "100.0%" | head -20
            exit 1
          fi
      
      - name: Documentation Check
        run: |
          echo "Checking documentation requirements..."
          
          # Check for README files in new packages
          for dir in $(find . -name "*.go" -exec dirname {} \; | sort -u); do
            if [ ! -f "$dir/README.md" ] && [ ! -f "$dir/doc.go" ]; then
              echo "::warning::Missing documentation in $dir"
            fi
          done
          
          # Check for godoc comments on public functions
          missing_docs=$(go doc -all . | grep -c "func.*{" || true)
          if [ "$missing_docs" -gt 0 ]; then
            echo "::warning::$missing_docs public functions may be missing documentation"
          fi

  performance-gates:
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
      
      - name: Setup Test Environment
        run: |
          docker-compose -f docker-compose.test.yml up -d
          sleep 30  # Wait for services to be ready
      
      - name: Run Performance Tests
        run: |
          echo "Running performance tests..."
          
          # Install k6
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
          
          # Run load tests
          k6 run tests/performance/load-test.js
      
      - name: Performance Budget Check
        run: |
          echo "Checking performance budgets..."
          
          # Extract performance metrics from test results
          p95_response_time=$(jq -r '.metrics.http_req_duration.values.p95' results.json)
          error_rate=$(jq -r '.metrics.http_req_failed.values.rate' results.json)
          
          # Check against budgets
          if (( $(echo "$p95_response_time > 500" | bc -l) )); then
            echo "::error::P95 response time ${p95_response_time}ms exceeds budget (500ms)"
            exit 1
          fi
          
          if (( $(echo "$error_rate > 0.01" | bc -l) )); then
            echo "::error::Error rate ${error_rate} exceeds budget (1%)"
            exit 1
          fi

  integration-gates:
    runs-on: ubuntu-latest
    needs: [enforce-constraints]
    
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
      
      - name: Integration Test Suite
        run: |
          echo "Running integration tests..."
          
          # Start test dependencies
          docker-compose -f docker-compose.test.yml up -d
          
          # Wait for readiness
          ./scripts/wait-for-services.sh
          
          # Run integration tests
          go test -tags=integration ./tests/integration/...
      
      - name: Contract Validation
        run: |
          echo "Validating API contracts..."
          
          # Generate OpenAPI spec from code
          swag init -g main.go -o ./api/docs
          
          # Validate against existing contracts
          if ! diff -q api/docs/swagger.yaml api/contracts/current.yaml; then
            echo "::error::API contract has changed without proper versioning"
            echo "Run: make update-contracts"
            exit 1
          fi

  deployment-gates:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    needs: [enforce-constraints, performance-gates, integration-gates]
    
    steps:
      - name: Security Final Check
        run: |
          echo "Final security validation..."
          
          # Container security scan
          docker build -t app:latest .
          
          # Trivy scan
          curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin
          trivy image --exit-code 1 --severity HIGH,CRITICAL app:latest
      
      - name: Deployment Readiness Check
        run: |
          echo "Checking deployment readiness..."
          
          # Verify all quality gates passed
          echo "✅ Constraints enforced"
          echo "✅ Performance tests passed"
          echo "✅ Integration tests passed"
          echo "✅ Security scans clean"
          echo "🚀 Ready for deployment"
```

### Branch Protection Rules

**GitHub Branch Protection**:
```yaml
# .github/branch-protection.yml
protection_rules:
  main:
    required_status_checks:
      strict: true
      contexts:
        - "enforce-constraints"
        - "performance-gates"  
        - "integration-gates"
    
    enforce_admins: true
    required_pull_request_reviews:
      required_approving_reviews: 2
      dismiss_stale_reviews: true
      require_code_owner_reviews: true
    
    restrictions:
      users: []
      teams: ["senior-developers", "architects"]
    
    allow_force_pushes: false
    allow_deletions: false
```

## Human Oversight

### Code Review Requirements

**Review Checklist Template**:
```markdown
## CDP Compliance Review

### Constraints Check
- [ ] Feature count within limits (≤3 concurrent)
- [ ] Complexity scores acceptable (≤10)
- [ ] Test coverage meets minimum (≥80%)
- [ ] No security vulnerabilities introduced

### Code Quality
- [ ] Code follows established patterns and conventions
- [ ] Appropriate error handling implemented
- [ ] Documentation updated for public APIs
- [ ] No obvious performance issues

### Testing
- [ ] Unit tests cover new functionality
- [ ] Integration tests updated where needed
- [ ] Edge cases and error conditions tested
- [ ] Performance impact assessed

### Architecture
- [ ] Changes align with system architecture
- [ ] Dependencies justified and minimal
- [ ] Backward compatibility maintained
- [ ] Migration strategy documented (if needed)

### Business Requirements
- [ ] Requirements clearly understood and met
- [ ] Acceptance criteria satisfied
- [ ] User experience considered
- [ ] Monitoring and alerting addressed

## Reviewer Notes
[Add specific feedback and suggestions here]

## Final Assessment
- [ ] ✅ Approved - Ready to merge
- [ ] ⚠️  Approved with minor changes
- [ ] ❌ Changes requested - Major issues found
```

**Reviewer Assignment Rules**:
```yaml
# .github/CODEOWNERS
# Global reviewers
* @senior-dev-team

# Critical system components require architect review  
/internal/core/ @architects @security-team
/api/ @api-team @architects

# Database changes require DBA review
**/migrations/ @database-team
**/schema/ @database-team

# Security-sensitive code requires security review
/internal/auth/ @security-team
/internal/crypto/ @security-team

# Performance-critical code requires performance review
/internal/cache/ @performance-team
/internal/queue/ @performance-team
```

### Architecture Review Board (ARB)

**ARB Charter**:
- Review all major architectural decisions
- Assess constraint violations and approve exceptions
- Evaluate system complexity and technical debt
- Guide technology choices and standards

**Review Process**:
```yaml
arb_review_process:
  triggers:
    - new_service_creation
    - major_architecture_change
    - constraint_violation_exception
    - technology_stack_addition
    - performance_degradation
  
  participants:
    required:
      - lead_architect
      - senior_developer
      - product_owner
    optional:
      - security_architect
      - devops_engineer
      - domain_expert
  
  deliverables:
    - architecture_decision_record
    - risk_assessment
    - implementation_timeline
    - monitoring_strategy
```

**Architecture Decision Record Template**:
```markdown
# ADR-XXX: [Title]

## Status
[Proposed | Accepted | Rejected | Superseded]

## Context
[Describe the problem and environmental forces]

## Decision
[Describe the chosen solution and rationale]

## Consequences
[Describe positive and negative outcomes]

## Implementation
- Timeline: [Expected completion date]
- Resources: [Required team members and time]
- Dependencies: [External requirements]
- Risks: [Potential issues and mitigations]

## Monitoring
- Success Metrics: [How to measure success]
- Alert Conditions: [When to raise concerns]
- Review Schedule: [When to reassess]

## References
- [Related ADRs, documents, or resources]
```

### Quality Assurance Reviews

**QA Review Process**:
1. **Feature Testing**: Comprehensive functional testing
2. **Regression Testing**: Ensure existing functionality intact
3. **Performance Testing**: Validate performance requirements
4. **Security Testing**: Identify potential vulnerabilities
5. **Usability Testing**: Assess user experience

**QA Sign-off Criteria**:
```yaml
qa_signoff:
  functional:
    - all_test_cases_executed: true
    - critical_bugs_resolved: true
    - acceptance_criteria_met: true
  
  non_functional:
    - performance_requirements_met: true
    - security_requirements_met: true
    - accessibility_requirements_met: true
  
  process:
    - test_documentation_complete: true
    - defect_tracking_current: true
    - deployment_guide_updated: true
```

## Monitoring and Alerting

### Compliance Monitoring

**Metrics Dashboard**:
```yaml
compliance_metrics:
  constraint_violations:
    feature_limit_breaches:
      current: 0
      target: 0
      trend: "stable"
    
    complexity_violations:
      current: 2
      target: 0
      trend: "improving"
    
    coverage_violations:
      current: 1
      target: 0
      trend: "stable"
  
  process_adherence:
    code_review_compliance:
      current: 98%
      target: 100%
      trend: "stable"
    
    testing_compliance:
      current: 95%
      target: 100%
      trend: "improving"
    
    documentation_compliance:
      current: 87%
      target: 95%
      trend: "improving"
```

**Alert Configuration**:
```yaml
alerts:
  - name: "Constraint Violation"
    condition: "constraint_violations > 0"
    severity: "high"
    channels: ["slack-dev", "email-leads"]
    message: "CDP constraint violation detected: {{ .violation_type }}"
  
  - name: "Coverage Drop"
    condition: "test_coverage < 80"
    severity: "medium"
    channels: ["slack-dev"]
    message: "Test coverage dropped to {{ .coverage }}% in {{ .component }}"
  
  - name: "Review Bottleneck"
    condition: "avg_review_time > 24h"
    severity: "medium"
    channels: ["slack-leads"]
    message: "Code review time exceeding SLA: {{ .avg_time }}"
```

### Performance Monitoring

**Performance Dashboards**:
```yaml
performance_monitoring:
  build_times:
    current_avg: "3m 45s"
    target: "5m"
    trend: "improving"
  
  test_execution:
    unit_tests: "45s"
    integration_tests: "8m 30s"
    system_tests: "23m 15s"
  
  deployment_frequency:
    current: "3.2 per day"
    target: "2+ per day"
    trend: "stable"
  
  lead_time:
    current: "2.1 days"
    target: "3 days"
    trend: "improving"
```

## Exception Handling

### Constraint Exception Process

**Exception Request Template**:
```yaml
constraint_exception:
  id: "EXC-2024-001"
  requestor: "john.doe@company.com"
  date_requested: "2024-01-15"
  
  violation:
    type: "complexity_score"
    current_value: 12
    limit: 10
    component: "src/user/authentication.go"
  
  justification:
    business_reason: "Critical security fix required before deadline"
    technical_reason: "Refactoring would require additional 2 weeks"
    risk_assessment: "Low - well-tested legacy code"
  
  mitigation_plan:
    immediate: "Add comprehensive unit tests"
    short_term: "Schedule refactoring in next sprint"
    long_term: "Redesign authentication module"
  
  approval:
    arb_review_date: "2024-01-16"
    approved_by: "jane.smith@company.com"
    expiry_date: "2024-02-15"
    conditions: ["Must add tests", "Refactor by Q2"]
```

**Exception Tracking**:
```bash
# CLI tool for exception management
ddx exceptions list --status=active
ddx exceptions create --type=complexity --component=auth
ddx exceptions approve EXC-2024-001 --reviewer="lead-architect"
ddx exceptions close EXC-2024-001 --reason="refactoring-complete"
```

### Emergency Procedures

**Hotfix Process**:
```yaml
hotfix_procedure:
  triggers:
    - production_outage
    - critical_security_vulnerability
    - data_corruption_risk
  
  expedited_process:
    - skip_feature_limit_check: true
    - reduce_review_requirement: 1  # Instead of 2
    - allow_direct_main_push: true  # With approval
    - bypass_system_tests: false    # Never skip
  
  mandatory_requirements:
    - security_scan: required
    - unit_tests: required  
    - rollback_plan: required
    - incident_documentation: required
  
  post_deployment:
    - immediate_monitoring: 30_minutes
    - follow_up_review: 24_hours
    - process_improvement: 7_days
```

## Continuous Improvement

### Enforcement Evolution

**Regular Reviews**:
- Weekly: Constraint violation analysis
- Monthly: Enforcement effectiveness assessment  
- Quarterly: Process improvement planning
- Annually: Complete enforcement framework review

**Improvement Process**:
```yaml
improvement_cycle:
  data_collection:
    - violation_metrics
    - developer_feedback
    - process_efficiency_data
    - business_impact_analysis
  
  analysis:
    - root_cause_analysis
    - trend_identification
    - cost_benefit_assessment
    - risk_evaluation
  
  implementation:
    - pilot_programs
    - gradual_rollout
    - training_programs
    - tool_updates
  
  validation:
    - effectiveness_measurement
    - developer_satisfaction
    - business_value_assessment
    - constraint_optimization
```

**Feedback Integration**:
```yaml
feedback_mechanisms:
  developer_surveys:
    frequency: quarterly
    focus: process_burden, tool_effectiveness
  
  retrospectives:
    frequency: sprint_end
    focus: constraint_impact, improvement_ideas
  
  metrics_analysis:
    frequency: weekly
    focus: violation_trends, resolution_times
  
  stakeholder_reviews:
    frequency: monthly
    focus: business_value, delivery_impact
```

This comprehensive enforcement system ensures that CDP principles are consistently applied while remaining flexible enough to handle exceptional circumstances and continuous improvement.