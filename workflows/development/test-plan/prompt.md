---
tags: [prompt, test-plan, development-workflow, quality-assurance, testing, ai-assisted]
references: template.md
created: 2025-01-12
modified: 2025-01-12
---

# Test Plan Creation Assistant

This prompt helps you create comprehensive test plans using the DDX test plan template. Work through each section systematically to develop effective testing strategies that ensure software quality and reliability.

## Using This Prompt

1. Start with the relevant feature specification or requirements document
2. Identify the testing scope and objectives
3. Answer the guiding questions for each section
4. Use the provided frameworks for systematic test planning
5. Review with stakeholders and adjust based on feedback

## Template Reference

This prompt uses the test plan template at [[template|template.md]]. You can also reference the comprehensive test planning guide at [[README|README.md]] and explore examples in the [[examples|examples/]] directory.

## Prerequisites Checklist

**Before creating a test plan, ensure you have:**
- [ ] Feature specifications or requirements documents
- [ ] Understanding of system architecture and components
- [ ] Knowledge of user workflows and acceptance criteria
- [ ] Access to test environments and tools
- [ ] Understanding of performance and security requirements
- [ ] List of stakeholders and their testing expectations
- [ ] Timeline and resource constraints

**Initial Assessment Questions:**
- What is the primary purpose of this testing effort?
- What are the main risks if this feature fails in production?
- Who are the target users and what are their critical workflows?
- What are the key quality attributes (performance, security, usability)?
- What testing constraints do we need to work within?

## Section-by-Section Guidance

### Executive Summary and Strategy

**Testing Objectives Framework:**
- **Functional Validation**: Verify that features work as specified
- **Quality Assurance**: Ensure quality attributes are met
- **Risk Mitigation**: Identify and prevent critical failures
- **User Experience**: Validate user workflows and satisfaction
- **System Integration**: Ensure components work together properly

**Scope Definition Questions:**
- What specific features, components, or systems will be tested?
- What aspects are explicitly excluded from testing?
- What are the testing boundaries and interfaces?
- What level of testing is required (unit, integration, system, acceptance)?
- What quality attributes need validation?

**Quality Goals Framework:**
Use SMART criteria for quality objectives:
- **Functional Quality**: Defect rates, test coverage, requirement compliance
- **Performance Quality**: Response times, throughput, resource utilization
- **Security Quality**: Vulnerability assessments, access control validation
- **Usability Quality**: User satisfaction, task completion rates, accessibility
- **Reliability Quality**: Availability, error rates, recovery times

### Test Strategy Development

**Test Level Planning:**

**Unit Testing Strategy:**
- What components need unit testing?
- What testing frameworks and tools will be used?
- What code coverage targets are appropriate?
- Who is responsible for unit test development and execution?
- How will unit test results be integrated into CI/CD?

**Integration Testing Strategy:**
- What integration points need testing?
- How will data flow between components be validated?
- What API contracts need verification?
- How will external dependencies be handled (mocking, test doubles)?
- What integration environments are required?

**System Testing Strategy:**
- What end-to-end workflows need validation?
- How will non-functional requirements be tested?
- What system configurations need testing?
- How will system testing be automated?
- What deployment scenarios need validation?

**Acceptance Testing Strategy:**
- Who will participate in acceptance testing?
- What business scenarios need validation?
- How will user feedback be collected and incorporated?
- What are the acceptance criteria for release?
- How will acceptance testing be managed and tracked?

### Test Scope Analysis

**Feature Prioritization Framework:**
- **Critical Features**: Must work perfectly (core business functions)
- **Important Features**: Should work well (supporting functions)
- **Nice-to-have Features**: Could have issues without major impact

**Risk-Based Testing Approach:**
For each feature, assess:
- **Business Impact**: What happens if this fails?
- **Technical Complexity**: How complex is the implementation?
- **Change Frequency**: How often does this area change?
- **User Visibility**: How visible are failures to users?
- **Integration Points**: How many other systems does this affect?

**Testing Boundary Questions:**
- What systems or components are in scope for testing?
- What external dependencies are out of scope?
- What legacy systems or features are excluded?
- What future features are not yet ready for testing?
- What third-party integrations need special consideration?

### Test Environment Planning

**Environment Strategy Framework:**

**Environment Types:**
- **Development Environment**: For early testing and debugging
- **Integration Environment**: For component interaction testing
- **System Test Environment**: For end-to-end testing
- **Performance Environment**: For load and stress testing
- **User Acceptance Environment**: For business user testing

**For Each Environment, Define:**
- **Purpose**: What type of testing will be performed?
- **Configuration**: What hardware, software, and network setup?
- **Data**: What test data is required and how is it managed?
- **Access**: Who has access and what are the security requirements?
- **Maintenance**: How is the environment maintained and refreshed?

**Environment Management Questions:**
- How closely should test environments mirror production?
- What are the data privacy and security requirements?
- How will environment configurations be versioned?
- What monitoring and alerting is needed?
- How will environment issues be reported and resolved?

### Test Data Strategy

**Test Data Categories:**

**Functional Test Data:**
- **Positive Test Data**: Valid inputs that should succeed
- **Negative Test Data**: Invalid inputs that should fail gracefully
- **Boundary Test Data**: Edge cases and limits
- **Volume Test Data**: Large datasets for performance testing

**Test Data Sources:**
- **Production Subsets**: Carefully selected and anonymized production data
- **Synthetic Data**: Generated data that mimics production characteristics
- **Manual Test Data**: Hand-crafted data for specific scenarios
- **API-Generated Data**: Data created through application APIs

**Data Management Framework:**
- **Data Creation**: How will test data be created and maintained?
- **Data Privacy**: How will sensitive data be protected?
- **Data Refresh**: How often will test data be updated?
- **Data Isolation**: How will test data be isolated between tests?
- **Data Cleanup**: How will test data be cleaned up after testing?

### Test Case Development

**Test Case Design Framework:**

**Functional Test Case Structure:**
For each test case, define:
- **Objective**: What specific functionality is being tested?
- **Preconditions**: What setup is required before the test?
- **Test Steps**: Detailed steps to execute the test
- **Expected Results**: What should happen at each step?
- **Test Data**: What specific data is needed?
- **Dependencies**: What other tests or systems are required?

**Test Case Categories:**
- **Happy Path Tests**: Normal, expected usage scenarios
- **Alternative Path Tests**: Different ways to accomplish goals
- **Error Path Tests**: How the system handles various failures
- **Boundary Tests**: Edge cases and limit conditions
- **Integration Tests**: Cross-system or cross-component interactions

**Test Coverage Analysis:**
- **Requirement Coverage**: Are all requirements tested?
- **Code Coverage**: Are sufficient code paths exercised?
- **User Journey Coverage**: Are critical user workflows tested?
- **Error Scenario Coverage**: Are error conditions properly tested?
- **Integration Coverage**: Are all integration points tested?

### Performance Testing Strategy

**Performance Test Types:**

**Load Testing:**
- **Objective**: Validate system performance under normal expected load
- **User Load**: How many concurrent users?
- **Duration**: How long should the test run?
- **Success Criteria**: What response times and throughput are acceptable?

**Stress Testing:**
- **Objective**: Determine system breaking point
- **Load Profile**: How will load be increased beyond normal capacity?
- **Failure Criteria**: What constitutes system failure?
- **Recovery Testing**: How quickly does the system recover?

**Volume Testing:**
- **Objective**: Test with large amounts of data
- **Data Volumes**: How much data should be processed?
- **Storage Testing**: How does storage affect performance?
- **Query Performance**: How do database queries perform with large datasets?

**Performance Metrics Framework:**
- **Response Time**: Time to complete individual operations
- **Throughput**: Number of operations per unit time
- **Concurrent Users**: Number of simultaneous active users
- **Resource Utilization**: CPU, memory, network, and disk usage
- **Error Rates**: Frequency of errors under load

### Security Testing Planning

**Security Testing Categories:**

**Authentication Testing:**
- **Valid Credentials**: Do valid credentials work properly?
- **Invalid Credentials**: Are invalid credentials properly rejected?
- **Password Policies**: Are password complexity rules enforced?
- **Session Management**: Are sessions properly created and destroyed?
- **Multi-Factor Authentication**: Do MFA mechanisms work correctly?

**Authorization Testing:**
- **Role-Based Access**: Do user roles control access properly?
- **Resource Protection**: Are protected resources properly secured?
- **Privilege Escalation**: Can users gain unauthorized access?
- **Cross-User Access**: Can users access other users' data?

**Input Validation Testing:**
- **SQL Injection**: Are database queries protected?
- **Cross-Site Scripting (XSS)**: Are user inputs properly sanitized?
- **Command Injection**: Are system commands protected?
- **File Upload Security**: Are uploaded files properly validated?

**Security Testing Tools:**
- **Static Analysis**: Automated code scanning for vulnerabilities
- **Dynamic Analysis**: Runtime security testing
- **Penetration Testing**: Simulated attack scenarios
- **Vulnerability Scanning**: Automated vulnerability detection

### Automation Strategy

**Automation Scope Decision Framework:**

**Good Automation Candidates:**
- **Repetitive Tests**: Tests executed frequently across releases
- **Regression Tests**: Tests that validate existing functionality
- **Data-Driven Tests**: Tests with multiple data scenarios
- **Performance Tests**: Load and stress testing scenarios
- **API Tests**: Service interface validation

**Manual Testing Candidates:**
- **Exploratory Testing**: Creative, investigative testing
- **Usability Testing**: User experience evaluation
- **Visual Testing**: UI appearance and layout validation
- **Ad-hoc Testing**: One-time or rarely executed tests
- **Complex Setup Tests**: Tests requiring complex manual setup

**Automation Framework Design:**
- **Tool Selection**: What tools and frameworks will be used?
- **Test Architecture**: How will tests be structured and organized?
- **Data Management**: How will test data be managed in automation?
- **Reporting**: How will automated test results be reported?
- **Maintenance**: How will automated tests be maintained and updated?

### Risk Assessment Framework

**Risk Identification Categories:**

**Technical Risks:**
- **Environment Risks**: Test environment availability and stability
- **Tool Risks**: Testing tool limitations or failures
- **Data Risks**: Test data quality and availability issues
- **Integration Risks**: Dependencies on external systems
- **Performance Risks**: System performance under load

**Process Risks:**
- **Schedule Risks**: Insufficient time for thorough testing
- **Resource Risks**: Inadequate testing skills or capacity
- **Communication Risks**: Poor stakeholder communication
- **Quality Risks**: Inadequate test coverage or depth
- **Change Risks**: Requirements or scope changes during testing

**For Each Risk, Document:**
- **Description**: What exactly is the risk?
- **Probability**: How likely is it to occur? (Low/Medium/High)
- **Impact**: What would be the consequences? (Low/Medium/High)
- **Risk Level**: Overall risk assessment (Low/Medium/High)
- **Mitigation**: How can we prevent or minimize the risk?
- **Contingency**: What will we do if the risk materializes?

### Metrics and Reporting Strategy

**Test Metrics Framework:**

**Execution Metrics:**
- **Test Case Execution Rate**: Percentage of planned tests executed
- **Test Pass Rate**: Percentage of executed tests that pass
- **Test Coverage**: Requirements and code coverage metrics
- **Defect Discovery Rate**: Number of defects found per time period

**Quality Metrics:**
- **Defect Density**: Number of defects per feature or component
- **Defect Escape Rate**: Defects found in production after testing
- **Test Effectiveness**: Percentage of defects caught by testing
- **Customer Impact**: Severity and impact of escaped defects

**Efficiency Metrics:**
- **Test Automation Coverage**: Percentage of tests automated
- **Test Execution Time**: Time required to execute test suites
- **Environment Utilization**: Test environment availability and usage
- **Resource Productivity**: Team output and efficiency measures

**Reporting Strategy:**
- **Daily Reports**: Quick status updates for the team
- **Weekly Reports**: Detailed progress and issue reports for stakeholders
- **Milestone Reports**: Comprehensive reports at key project gates
- **Final Report**: Complete test summary and quality assessment

## Template Usage

```markdown
{{include: template.md}}
```

## Quality Assurance Framework

### Test Plan Review Checklist

**Completeness:**
- [ ] All test objectives are clearly defined
- [ ] Test scope is comprehensive and realistic
- [ ] All test types and levels are addressed
- [ ] Test cases cover all critical requirements
- [ ] Resource and environment needs are identified
- [ ] Schedule and milestones are realistic

**Clarity:**
- [ ] Test objectives are specific and measurable
- [ ] Test procedures are easy to follow
- [ ] Success criteria are clearly defined
- [ ] Responsibilities are clearly assigned
- [ ] Communication plan is comprehensive

**Feasibility:**
- [ ] Test approach is technically feasible
- [ ] Resource requirements are available
- [ ] Timeline is realistic given scope
- [ ] Test environments are available or obtainable
- [ ] Tools and frameworks are appropriate

### Stakeholder Review Process

**Technical Review Focus:**
- Test approach alignment with system architecture
- Technical feasibility of test procedures
- Tool and framework appropriateness
- Environment and infrastructure requirements

**Business Review Focus:**
- Alignment with business requirements and priorities
- User acceptance testing approach
- Risk coverage and mitigation strategies
- Timeline and resource implications

**Quality Assurance Review Focus:**
- Test coverage adequacy
- Test procedure effectiveness
- Defect management process
- Quality metrics and reporting approach

## Common Test Planning Patterns

### Web Application Testing
Standard patterns for:
- Cross-browser compatibility testing
- Responsive design validation
- User interface testing
- API testing
- Security testing

### Mobile Application Testing
Patterns for:
- Device compatibility testing
- Performance on mobile networks
- Battery usage testing
- Touch interface testing
- App store validation

### API Testing
Patterns for:
- Functional API testing
- API contract validation
- Performance testing
- Security testing
- Error handling validation

### Database Testing
Patterns for:
- Data integrity testing
- Query performance testing
- Backup and recovery testing
- Migration testing
- Concurrency testing

## Integration with Development Process

### Agile Testing Integration
- **Sprint Planning**: Test planning integrated into development planning
- **Continuous Testing**: Tests executed throughout development cycle
- **Feedback Loops**: Rapid feedback on quality issues
- **Test-Driven Development**: Tests driving development implementation

### DevOps Integration
- **Continuous Integration**: Automated testing in CI/CD pipelines
- **Environment as Code**: Test environment provisioning automation
- **Shift-Left Testing**: Earlier testing in the development lifecycle
- **Production Monitoring**: Continuous quality monitoring in production

## Advanced Testing Considerations

### Exploratory Testing
- **Charter-Based Testing**: Structured exploratory testing approach
- **Session-Based Testing**: Time-boxed exploratory testing sessions
- **Risk-Based Exploration**: Focus on high-risk areas
- **Tool-Assisted Exploration**: Using tools to enhance exploration

### Accessibility Testing
- **Screen Reader Testing**: Validation with assistive technologies
- **Keyboard Navigation**: Testing without mouse interaction
- **Color Contrast**: Visual accessibility validation
- **WCAG Compliance**: Standards-based accessibility testing

### Internationalization Testing
- **Multi-Language Support**: Testing with different languages
- **Character Encoding**: Unicode and character set validation
- **Date/Time Formats**: Localized format testing
- **Cultural Considerations**: Testing for cultural appropriateness

## Examples and References

- [[examples/web-application-testing|Web Application Test Plan]]: Complete example for web applications
- [[examples/api-testing|API Testing Plan]]: Comprehensive API testing strategy
- [[examples/mobile-testing|Mobile Application Test Plan]]: Mobile-specific testing approach
- [[../feature-spec/README|Feature Specifications]]: How test plans relate to feature requirements
- [[../architecture/README|Architecture Decisions]]: How testing aligns with system architecture

## Need Help?

- For test planning best practices: [[README|Test Planning Guide]]
- For feature specification alignment: [[../feature-spec/README|Feature Specification Guide]]
- For architecture considerations: [[../architecture/README|Architecture Decision Records]]
- For release planning: [[../release/README|Release Planning Guide]]
- For examples: Check the `examples/` directory

Remember: Test plans ensure software quality through systematic validation of functionality, performance, and security. They should be comprehensive enough to catch critical issues while remaining practical and maintainable for long-term success. Focus on risk-based testing to maximize value while working within resource and schedule constraints.