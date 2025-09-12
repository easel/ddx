---
tags: [prompt, feature-spec, development-workflow, implementation, requirements, ai-assisted]
references: template.md
created: 2025-01-12
modified: 2025-01-12
---

# Feature Specification Creation Assistant

This prompt helps you create comprehensive feature specifications using the DDX feature specification template. Work through each section systematically to transform product requirements into actionable technical specifications.

## Using This Prompt

1. Start with the relevant Product Requirements Document (PRD)
2. Identify the specific feature or capability being specified
3. Answer the guiding questions for each section
4. Use the provided frameworks for consistent analysis
5. Review with stakeholders before finalizing the specification

## Template Reference

This prompt uses the feature specification template at [[template|template.md]]. You can also reference the comprehensive feature specification guide at [[README|README.md]] and explore examples in the [[examples|examples/]] directory.

## Prerequisites Checklist

**Before creating a feature specification, ensure you have:**
- [ ] Relevant PRD section(s) that describe the product requirements
- [ ] Understanding of existing system architecture
- [ ] Access to design mockups or wireframes (if applicable)
- [ ] List of key stakeholders who need to review the specification
- [ ] Understanding of technical constraints and dependencies
- [ ] Clear timeline and implementation priorities

**Questions to Answer First:**
- Which PRD section(s) does this feature implement?
- What is the primary business goal this feature supports?
- Who are the target users and what are their needs?
- What are the main technical challenges anticipated?
- How does this feature fit into the overall product roadmap?

## Section-by-Section Guidance

### Overview and Context

**Feature Naming:**
- Choose a clear, descriptive name that reflects user value
- Avoid internal jargon or technical terminology
- Keep it concise but specific
- Examples: "Advanced Search Filters", "Real-time Collaboration", "User Profile Management"

**Business Context Questions:**
- How does this feature support our business objectives?
- What user problems does it solve?
- What are the expected business outcomes?
- How will success be measured?
- What is the priority relative to other features?

**Success Criteria Framework:**
Use SMART criteria (Specific, Measurable, Achievable, Relevant, Time-bound):
- **Adoption Metrics**: How many users will use this feature?
- **Engagement Metrics**: How frequently will it be used?
- **Business Metrics**: What business KPIs will improve?
- **Quality Metrics**: What quality standards must be met?
- **Performance Metrics**: What performance targets are required?

### User Stories and Acceptance Criteria

**User Story Quality Checklist:**
- [ ] Follows the "As a [user], I want [functionality], so that [benefit]" format
- [ ] Describes user value, not technical implementation
- [ ] Is testable and verifiable
- [ ] Is appropriately sized (not too big or too small)
- [ ] Has clear acceptance criteria

**Acceptance Criteria Best Practices:**
- Start each criterion with a testable action verb
- Focus on behavior and outcomes, not implementation
- Include both positive and negative test cases
- Consider edge cases and error conditions
- Make criteria specific and unambiguous

**User Story Prioritization:**
- **Primary Stories**: Core functionality that must be implemented
- **Secondary Stories**: Important but not critical functionality
- **Future Stories**: Nice-to-have features for later iterations

**Questions for Each User Story:**
- What is the user trying to accomplish?
- What context or situation triggers this need?
- What would constitute success from the user's perspective?
- What could go wrong, and how should the system respond?
- Are there any assumptions about user knowledge or capabilities?

### Functional Requirements Deep-dive

**Component Breakdown Strategy:**
For each component, define:
- **Purpose**: What business need does this component serve?
- **Scope**: What is and isn't included in this component?
- **Inputs**: What data or events trigger this component?
- **Processing**: What business logic or transformation occurs?
- **Outputs**: What results or side effects are produced?
- **State**: What information does this component maintain?

**Business Rules Documentation:**
- **Validation Rules**: What data validation is required?
- **Business Logic**: What calculations or decisions are made?
- **Workflow Rules**: What sequence of actions must occur?
- **Permission Rules**: Who can perform what actions?
- **Data Rules**: What data relationships must be maintained?

**User Interaction Mapping:**
For each workflow:
1. **Happy Path**: Normal sequence of actions and responses
2. **Alternative Paths**: Different ways users might accomplish the goal
3. **Error Paths**: What happens when things go wrong
4. **Recovery Paths**: How users can get back on track

### Technical Requirements Analysis

**Architecture Alignment Questions:**
- How does this feature fit into our existing architecture?
- Are there any architectural decisions that need to be made?
- What existing patterns and components can be reused?
- Are there new architectural patterns being introduced?
- What are the implications for system complexity?

**Component Design Framework:**

**Frontend Components:**
- **Presentation Logic**: What UI elements and interactions?
- **State Management**: What data needs to be maintained?
- **Side Effects**: What external calls or actions are triggered?
- **Reusability**: Can components be used elsewhere?
- **Testing Strategy**: How will components be tested?

**Backend Services:**
- **Business Logic**: What business rules are enforced?
- **Data Access**: How is data retrieved and persisted?
- **External Integration**: What external systems are involved?
- **Error Handling**: How are errors detected and managed?
- **Performance Considerations**: What are the scalability needs?

**Database Design Questions:**
- What new data needs to be stored?
- How does new data relate to existing data?
- What are the query patterns and performance implications?
- Are there any data migration requirements?
- What are the data retention and deletion policies?

### API Specifications

**API Design Principles:**
- **RESTful Design**: Follow REST conventions for HTTP APIs
- **Consistency**: Use consistent naming and structure patterns
- **Versioning**: Plan for API evolution and backward compatibility
- **Documentation**: Ensure APIs are self-documenting
- **Error Handling**: Provide clear, actionable error responses

**For Each Endpoint, Define:**
- **Purpose**: What business function does this endpoint serve?
- **Authentication**: What authentication is required?
- **Authorization**: What permissions are checked?
- **Input Validation**: What validation rules apply?
- **Business Logic**: What processing occurs?
- **Response Format**: What data is returned and in what structure?
- **Error Scenarios**: What errors can occur and how are they communicated?

**API Documentation Standards:**
- Use consistent parameter and response formats
- Provide clear examples for all endpoints
- Document all possible error responses
- Include authentication and authorization details
- Specify rate limiting and usage guidelines

### UI/UX Requirements

**User Experience Design Questions:**
- What is the user's mental model for this feature?
- How does this feature fit into existing user workflows?
- What are the key user tasks that need to be optimized?
- What are the most common user paths through the feature?
- How can we minimize cognitive load and friction?

**Interface Design Considerations:**
- **Information Architecture**: How is information organized and presented?
- **Navigation**: How do users move through the feature?
- **Feedback**: How does the system provide feedback to users?
- **Error Prevention**: How can we prevent user errors?
- **Recovery**: How can users correct mistakes?

**Responsive Design Framework:**
- **Mobile First**: Start with mobile constraints and scale up
- **Breakpoints**: Define specific screen size breakpoints
- **Content Priority**: What content is most important at each size?
- **Interaction Adaptations**: How do interactions change across devices?
- **Performance**: How do design choices affect performance?

**Accessibility Checklist:**
- [ ] Keyboard navigation for all interactive elements
- [ ] Screen reader support with proper ARIA labels
- [ ] Color contrast meets WCAG 2.1 AA standards
- [ ] Focus indicators are clearly visible
- [ ] Error messages are clearly associated with form fields
- [ ] Images have meaningful alternative text
- [ ] Form labels are properly associated with inputs

### Performance Requirements

**Performance Metrics Framework:**
- **Response Time**: How fast should the system respond?
- **Throughput**: How many operations per unit of time?
- **Concurrent Users**: How many simultaneous users?
- **Resource Usage**: What are the CPU, memory, and storage limits?
- **Scalability**: How should performance scale with growth?

**Performance Testing Strategy:**
- **Load Testing**: Normal expected traffic patterns
- **Stress Testing**: Peak traffic scenarios
- **Spike Testing**: Sudden traffic increases
- **Volume Testing**: Large amounts of data
- **Endurance Testing**: Extended periods of use

**Performance Optimization Considerations:**
- **Caching**: What data can be cached and for how long?
- **Database Optimization**: What indexes and query optimizations are needed?
- **Network Optimization**: How can network usage be minimized?
- **Rendering Optimization**: How can UI rendering be optimized?
- **Asset Optimization**: How can static assets be optimized?

### Security Requirements Analysis

**Security Threat Modeling:**
- **Authentication**: How do we verify user identity?
- **Authorization**: How do we control access to resources?
- **Input Validation**: How do we prevent malicious input?
- **Data Protection**: How do we protect sensitive data?
- **Communication Security**: How do we secure data in transit?
- **Audit Logging**: What security events need to be logged?

**Common Security Patterns:**
- **Input Sanitization**: Prevent injection attacks
- **Output Encoding**: Prevent XSS attacks
- **CSRF Protection**: Prevent cross-site request forgery
- **Rate Limiting**: Prevent abuse and DoS attacks
- **Session Management**: Secure session handling
- **Encryption**: Protect sensitive data at rest and in transit

### Error Handling and Resilience

**Error Classification Framework:**
- **User Errors**: Mistakes made by users (validation failures)
- **System Errors**: Internal system failures (database connectivity)
- **External Errors**: Failures in external dependencies
- **Network Errors**: Communication failures
- **Resource Errors**: Insufficient resources (memory, disk space)

**Error Recovery Strategies:**
- **Graceful Degradation**: Maintain core functionality when components fail
- **Retry Logic**: Automatically retry transient failures
- **Circuit Breakers**: Prevent cascading failures
- **Fallback Options**: Alternative ways to accomplish user goals
- **User Communication**: Clear error messages and guidance

### Testing Strategy Planning

**Testing Pyramid Application:**
- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete user workflows
- **Performance Tests**: Validate performance requirements
- **Security Tests**: Verify security controls

**Test Case Design:**
- **Happy Path Testing**: Normal, expected usage scenarios
- **Edge Case Testing**: Boundary conditions and unusual inputs
- **Error Path Testing**: How the system handles various failures
- **Regression Testing**: Ensure existing functionality still works
- **Usability Testing**: Validate user experience design

### Risk Assessment and Mitigation

**Risk Identification Framework:**

**Technical Risks:**
- **Complexity Risk**: Is the solution too complex to implement reliably?
- **Performance Risk**: Can the system meet performance requirements?
- **Integration Risk**: Are there challenges integrating with other systems?
- **Security Risk**: Are there potential security vulnerabilities?
- **Scalability Risk**: Will the solution scale as needed?

**Business Risks:**
- **Schedule Risk**: Can this be delivered on time?
- **Resource Risk**: Do we have the necessary skills and capacity?
- **User Adoption Risk**: Will users actually use this feature?
- **Competitive Risk**: How might competitors respond?
- **Regulatory Risk**: Are there compliance considerations?

**Risk Assessment Matrix:**
For each risk, evaluate:
- **Probability**: How likely is this risk to occur? (Low/Medium/High)
- **Impact**: What would be the consequences? (Low/Medium/High)
- **Detectability**: How early can we detect this risk materializing?
- **Mitigation Strategy**: How can we prevent or minimize this risk?
- **Contingency Plan**: What will we do if the risk occurs?

## Template Usage

```markdown
{{include: template.md}}
```

## Quality Assurance Framework

### Specification Review Checklist

**Completeness:**
- [ ] All template sections are addressed
- [ ] User stories have clear acceptance criteria
- [ ] Technical requirements are specific and testable
- [ ] APIs are fully specified with examples
- [ ] Performance and security requirements are defined
- [ ] Testing strategy is comprehensive

**Clarity:**
- [ ] Requirements are unambiguous
- [ ] Technical terms are defined or clearly understood
- [ ] User workflows are easy to follow
- [ ] Success criteria are measurable
- [ ] Dependencies are clearly identified

**Consistency:**
- [ ] Naming conventions are consistent throughout
- [ ] Requirements align with existing system patterns
- [ ] User stories map to functional requirements
- [ ] API specifications match data requirements
- [ ] Testing strategy covers all requirements

**Feasibility:**
- [ ] Technical approach is realistic given constraints
- [ ] Timeline expectations are reasonable
- [ ] Resource requirements are available
- [ ] Performance targets are achievable
- [ ] Security requirements can be implemented

### Stakeholder Review Process

**Product Review Focus:**
- Business alignment and user value
- Completeness of user stories and acceptance criteria
- Success metrics and measurement approach
- Priority and scope appropriateness

**Engineering Review Focus:**
- Technical feasibility and architecture alignment
- Implementation approach and complexity
- Performance and scalability considerations
- Testing strategy adequacy

**Design Review Focus:**
- User experience and interface design
- Accessibility and inclusive design
- Responsive design considerations
- Design system consistency

**Security Review Focus:**
- Security requirements completeness
- Threat modeling adequacy
- Compliance with security standards
- Privacy and data protection considerations

## Common Feature Specification Patterns

### CRUD Features
Standard Create, Read, Update, Delete functionality with:
- Data validation and business rules
- User permissions and access control
- Audit logging and change tracking
- Search and filtering capabilities

### Integration Features
Features that connect with external systems:
- API design and data mapping
- Error handling and retry logic
- Authentication and security
- Monitoring and alerting

### User Management Features
Features related to user accounts and access:
- Registration and onboarding flows
- Authentication and session management
- Profile management and preferences
- Permission and role management

### Reporting and Analytics Features
Features that provide data insights:
- Data collection and aggregation
- Query and filtering capabilities
- Visualization and presentation
- Export and sharing functionality

### Real-time Features
Features that provide live updates:
- Event-driven architecture considerations
- WebSocket or Server-Sent Event implementation
- State synchronization across clients
- Connection management and fallback strategies

## Integration with Development Workflow

### PRD to Feature Spec Mapping
- Each PRD requirement should map to specific feature specifications
- Business goals should translate to technical success criteria
- User personas should inform user story creation
- Market requirements should influence prioritization

### Feature Spec to Implementation Planning
- Break feature specifications into development tasks
- Identify dependencies between components
- Estimate effort and complexity
- Plan implementation phases and milestones

### Feature Spec to Testing
- Transform acceptance criteria into test cases
- Create test data based on data requirements
- Plan test environments and dependencies
- Design performance and security test scenarios

## Documentation Standards

### Writing Guidelines
- Use clear, concise language
- Avoid technical jargon when possible
- Define terms and acronyms on first use
- Use consistent formatting and structure
- Include examples and illustrations

### Review and Update Process
- Regular reviews to ensure continued relevance
- Version control for all changes
- Stakeholder notification of significant updates
- Archival process for obsolete specifications

## Getting Started Checklist

**Preparation Phase:**
- [ ] Identify the PRD section(s) being implemented
- [ ] Gather relevant stakeholders for input
- [ ] Review existing system architecture and constraints
- [ ] Collect design mockups and user research
- [ ] Understand timeline and resource constraints

**Creation Phase:**
- [ ] Use the template structure
- [ ] Work through each section systematically
- [ ] Focus on user value and business outcomes
- [ ] Include specific, testable requirements
- [ ] Document assumptions and dependencies

**Review Phase:**
- [ ] Conduct stakeholder reviews
- [ ] Address feedback and concerns
- [ ] Validate technical feasibility
- [ ] Confirm alignment with business goals
- [ ] Get formal approval from key stakeholders

**Implementation Phase:**
- [ ] Break specification into development tasks
- [ ] Plan implementation phases
- [ ] Track progress against specification
- [ ] Update specification as needed
- [ ] Validate implementation against requirements

## Examples and References

- [[examples/user-authentication|User Authentication Feature Spec]]: Complete example of authentication features
- [[examples/search-filters|Advanced Search Feature Spec]]: Complex UI feature with multiple components
- [[examples/api-integration|External API Integration Spec]]: Integration with third-party services
- [[../prd/README|PRD Documentation]]: How feature specs relate to product requirements
- [[../architecture/README|Architecture Decisions]]: How specs align with architectural choices

## Need Help?

- For feature specification best practices: [[README|Feature Specification Guide]]
- For PRD alignment: [[../prd/README|PRD Documentation]]
- For architecture guidance: [[../architecture/README|Architecture Decision Records]]
- For testing strategy: [[../test-plan/README|Test Planning Guide]]
- For examples: Check the `examples/` directory

Remember: Feature specifications bridge product vision and technical implementation. They should be detailed enough to guide development while remaining flexible enough to accommodate necessary changes during implementation. Focus on user value, technical clarity, and measurable outcomes.