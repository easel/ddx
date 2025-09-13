---
tags: [prompt, release-notes, development-workflow, communication, user-documentation, ai-assisted]
references: template.md
created: 2025-01-12
modified: 2025-01-12
---

# Release Notes Creation Assistant

This prompt helps you create comprehensive release notes using the DDX release notes template. Work through each section systematically to communicate software changes effectively to users, stakeholders, and development teams.

## Using This Prompt

1. Gather all changes, features, and fixes since the last release
2. Identify your target audiences and their information needs
3. Answer the guiding questions for each section
4. Use the provided frameworks for consistent communication
5. Review with stakeholders and test with representative users

## Template Reference

This prompt uses the release notes template at [[template|template.md]]. You can also reference the comprehensive release notes guide at [[README|README.md]] and explore examples in the [[examples|examples/]] directory.

## Prerequisites Checklist

**Before creating release notes, ensure you have:**
- [ ] Complete list of all changes since the last release
- [ ] Feature specifications and implementation details
- [ ] Test results and quality metrics
- [ ] User feedback and beta testing results
- [ ] Breaking changes and migration requirements
- [ ] Security updates and vulnerability fixes
- [ ] Performance improvements and benchmarks
- [ ] Known issues and limitations

**Information Gathering Questions:**
- What are the most significant changes in this release?
- Who are the primary audiences for these release notes?
- What changes require user action or attention?
- Are there any security or breaking changes?
- What support resources are available for users?

## Section-by-Section Guidance

### Release Overview and Highlights

**Release Classification:**
- **Major Release**: Significant new features, possible breaking changes
- **Minor Release**: New features and enhancements, backward compatible
- **Patch Release**: Bug fixes and small improvements
- **Hotfix**: Critical security or bug fixes
- **Beta/Preview**: Early access features for testing

**Highlight Selection Framework:**
Choose highlights based on:
- **User Impact**: Changes that most significantly affect user experience
- **Strategic Importance**: Features that align with product strategy
- **Community Requests**: Features users have been asking for
- **Competitive Advantage**: Unique capabilities or improvements
- **Performance Gains**: Measurable improvements users will notice

**Quick Stats Guidelines:**
- Use concrete numbers when possible
- Focus on metrics that matter to users
- Include performance improvements as percentages
- Mention scale of testing or validation
- Keep numbers accurate and verifiable

### New Features Documentation

**Feature Description Framework:**
For each new feature, answer:
- **What**: What does this feature do?
- **Why**: Why is it valuable to users?
- **Who**: Who will benefit from this feature?
- **When**: When and how can users access it?
- **Where**: Where in the product is it located?

**Feature Communication Pattern:**
1. **Lead with Benefit**: Start with user value, not technical details
2. **Explain Purpose**: Why this feature was created
3. **Provide Context**: How it fits into user workflows
4. **Include Examples**: Concrete use cases or scenarios
5. **Guide Usage**: Clear steps for getting started

**User-Centric Language Guidelines:**
- Use "you" instead of "users" when possible
- Focus on actions users can take
- Avoid technical jargon unless necessary
- Use active voice and clear verbs
- Include specific benefits, not just capabilities

### Enhancements Communication

**Enhancement Categories:**
- **Performance Improvements**: Speed, efficiency, resource usage
- **User Interface Updates**: Design, usability, accessibility
- **Workflow Enhancements**: Process improvements, automation
- **Integration Improvements**: Better connectivity with other tools
- **Reliability Enhancements**: Stability, error handling, recovery

**Enhancement Documentation Pattern:**
- **Before/After**: Describe the previous state and new state
- **Quantified Benefits**: Include metrics where possible
- **User Impact**: How users will notice the improvement
- **Availability**: When and where improvements are active
- **Additional Actions**: Any steps users can take to benefit more

**Performance Communication Guidelines:**
- Use specific percentages or times when possible
- Compare to previous versions, not competitors
- Include context about typical usage scenarios
- Mention any conditions that affect performance
- Provide benchmarking methodology if relevant

### Bug Fixes and Issue Resolution

**Bug Fix Prioritization:**
- **High Priority**: Issues that blocked users or caused data problems
- **Security Related**: Any security vulnerabilities or weaknesses
- **User Experience**: Problems that degraded user experience
- **Compatibility**: Issues with browsers, devices, or integrations
- **Performance**: Problems that affected system performance

**Issue Documentation Framework:**
For significant fixes:
- **Problem Description**: What was wrong?
- **User Impact**: Who was affected and how?
- **Resolution**: How the problem was solved
- **Prevention**: Steps taken to prevent recurrence
- **User Actions**: Any steps users should take

**Communication Sensitivity:**
- Acknowledge problems without over-emphasizing failures
- Focus on resolution and improvement
- Thank users for reporting issues when appropriate
- Avoid technical blame or complex technical explanations
- Emphasize commitment to quality and user experience

### Security Updates

**Security Communication Balance:**
- **Be Transparent**: Acknowledge security updates clearly
- **Avoid Details**: Don't provide information that could enable attacks
- **Emphasize Action**: Clear guidance on any required user actions
- **Build Confidence**: Demonstrate commitment to security
- **Provide Support**: Clear channels for security questions

**Security Update Categories:**
- **Proactive Improvements**: General security enhancements
- **Vulnerability Fixes**: Specific known vulnerabilities addressed
- **Compliance Updates**: Changes for regulatory compliance
- **Authentication/Authorization**: Access control improvements
- **Data Protection**: Encryption, privacy, data handling improvements

**Required Information for Security Updates:**
- Severity level (Critical, High, Medium, Low)
- Affected components or data types
- Required user actions (if any)
- Timeline for automatic updates
- How to get help with security concerns

### Breaking Changes Management

**Breaking Change Communication Framework:**
1. **Clear Identification**: Make breaking changes highly visible
2. **Impact Assessment**: Exactly what functionality changes
3. **Timeline Communication**: When changes take effect
4. **Migration Guidance**: Step-by-step adaptation instructions
5. **Support Resources**: Where to get help with changes

**Migration Documentation Pattern:**
- **Assessment**: How to determine if you're affected
- **Preparation**: Steps to take before the change
- **Migration**: Specific actions to adapt to the change
- **Validation**: How to verify the migration was successful
- **Troubleshooting**: Common issues and solutions

**Change Communication Timeline:**
- **Advance Notice**: Early warning in previous releases
- **Detailed Guidance**: Comprehensive migration information
- **Implementation**: Clear communication when changes take effect
- **Follow-up Support**: Continued assistance after changes
- **Retrospective**: Learn from user experience with changes

### Technical Details Documentation

**Developer-Focused Information:**
- **API Changes**: New, modified, and deprecated endpoints
- **Database Changes**: Schema updates and migration requirements
- **Configuration Updates**: New settings and changed defaults
- **Integration Changes**: Updates to third-party integrations
- **Development Tools**: Changes to SDKs, libraries, or frameworks

**API Documentation Standards:**
- Include exact endpoint paths and methods
- Show request and response examples
- Document new parameters and fields
- Explain deprecation timelines
- Provide migration paths for changed endpoints

**Technical Audience Considerations:**
- Use precise technical terminology appropriately
- Include code examples where helpful
- Link to detailed technical documentation
- Provide troubleshooting guidance
- Offer direct support channels for technical issues

### System Requirements and Compatibility

**Requirement Updates:**
- **Operating System**: Minimum and recommended OS versions
- **Browser Support**: Supported browsers and version ranges
- **Hardware Requirements**: CPU, memory, storage, network needs
- **Software Dependencies**: Required frameworks, libraries, or tools
- **Network Requirements**: Bandwidth, connectivity, firewall considerations

**Compatibility Communication:**
- **Backward Compatibility**: What continues to work from previous versions
- **Forward Compatibility**: Considerations for future updates
- **Platform Support**: Multi-platform availability and differences
- **Integration Compatibility**: Third-party tool and service compatibility
- **Data Compatibility**: File format and data migration considerations

### Installation and Upgrade Guidance

**Installation Documentation:**
- **Prerequisites**: What needs to be in place before installation
- **Step-by-Step Process**: Clear, numbered installation steps
- **Verification**: How to confirm successful installation
- **Common Issues**: Typical problems and solutions
- **Support Resources**: Where to get help with installation

**Upgrade Strategy Communication:**
- **Backup Recommendations**: Essential data protection steps
- **Upgrade Options**: Automatic vs. manual upgrade paths
- **Downtime Expectations**: Expected service interruptions
- **Rollback Plans**: How to revert if problems occur
- **Post-Upgrade Actions**: Required steps after upgrade completion

### Known Issues and Support

**Issue Documentation Guidelines:**
- **Clear Description**: Exact nature of the issue
- **Impact Assessment**: Who is affected and under what conditions
- **Workaround Information**: Temporary solutions if available
- **Resolution Timeline**: Expected fix timeframe
- **Progress Updates**: How users will be informed of progress

**Support Resource Organization:**
- **Self-Service Options**: Documentation, FAQs, knowledge base articles
- **Community Resources**: Forums, user groups, community support
- **Direct Support**: How to contact technical support
- **Priority Support**: Escalation paths for critical issues
- **Educational Resources**: Training, webinars, tutorials

## Communication Strategy Framework

### Audience-Specific Messaging

**End User Focus:**
- **Headline Benefits**: What's most exciting and valuable
- **Visual Elements**: Screenshots, videos, interactive demos
- **Getting Started**: Simple first steps to try new features
- **Learning Resources**: Tutorials, guides, and help articles
- **Community Connection**: Ways to connect with other users

**Developer Focus:**
- **Technical Details**: API changes, code examples, integration guides
- **Migration Tools**: Automated or assisted migration resources
- **Development Resources**: Updated SDKs, documentation, samples
- **Testing Guidance**: How to validate integrations and implementations
- **Developer Community**: Technical forums, GitHub repositories, office hours

**Administrator Focus:**
- **Deployment Guidance**: Installation, configuration, and management
- **Security Considerations**: Security implications and best practices
- **Monitoring and Maintenance**: Operational considerations
- **User Management**: Changes affecting user administration
- **Compliance**: Regulatory and policy implications

### Multi-Channel Distribution Strategy

**Channel Selection Framework:**
- **In-Product Notifications**: For users actively using the product
- **Email Communications**: For broader user base and important updates
- **Website/Blog Posts**: For detailed information and SEO benefits
- **Social Media**: For community engagement and broader reach
- **Documentation Sites**: For detailed technical information
- **Partner Communications**: For integration partners and resellers

**Content Adaptation Guidelines:**
- **Social Media**: Highlight key benefits with visual elements
- **Email**: Executive summary with links to full details
- **In-Product**: Contextual notifications about relevant changes
- **Technical Docs**: Comprehensive implementation details
- **Blog Posts**: Storytelling approach with user benefits focus

### Feedback and Engagement Strategy

**Feedback Collection Methods:**
- **Direct Feedback Forms**: Specific feedback on release content
- **Community Discussions**: Open forums for user conversations
- **Support Channel Monitoring**: Track questions and issues
- **Usage Analytics**: Monitor feature adoption and usage patterns
- **User Interviews**: Deep dive conversations with key users

**Response Planning:**
- **FAQ Development**: Common questions and comprehensive answers
- **Content Updates**: Refinements based on user confusion
- **Additional Resources**: Supplementary guides and tutorials
- **Direct Outreach**: Proactive communication with affected users
- **Community Engagement**: Active participation in user discussions

## Template Usage

```markdown
{{include: template.md}}
```

## Quality Assurance Framework

### Content Review Checklist

**Accuracy and Completeness:**
- [ ] All significant changes are documented
- [ ] Technical details are accurate and verified
- [ ] Links and references are working and current
- [ ] Version numbers and dates are correct
- [ ] Security information is appropriately detailed

**Clarity and Usability:**
- [ ] Language is clear and appropriate for target audience
- [ ] Instructions are easy to follow
- [ ] Benefits are clearly communicated
- [ ] Examples and use cases are relevant
- [ ] Visual hierarchy supports scanning and reading

**Brand and Style:**
- [ ] Tone matches brand voice and previous communications
- [ ] Terminology is consistent throughout
- [ ] Formatting follows established standards
- [ ] Visual elements support the content
- [ ] Legal and compliance requirements are met

### User Testing and Validation

**Pre-Release Testing:**
- **Internal Review**: Technical accuracy and completeness
- **Stakeholder Review**: Business alignment and messaging
- **User Representative Review**: Clarity and usefulness from user perspective
- **Editorial Review**: Grammar, style, and consistency
- **Legal Review**: Compliance and risk assessment

**Post-Release Monitoring:**
- **User Feedback Analysis**: Themes and issues in user responses
- **Support Ticket Analysis**: Common questions and confusion points
- **Usage Analytics**: How users interact with new features
- **Community Sentiment**: Overall reception and discussion themes
- **Adoption Metrics**: Feature uptake and usage patterns

## Common Release Note Patterns

### Major Release Communications
- **Vision and Strategy**: How this release advances product vision
- **Comprehensive Feature Overview**: Detailed coverage of major new capabilities
- **Migration and Upgrade Focus**: Extensive guidance for significant changes
- **Community and Ecosystem Impact**: Effects on integrations and partners

### Minor Release Communications
- **Feature Enhancement Focus**: Improvements to existing capabilities
- **User Experience Improvements**: Usability and workflow enhancements
- **Performance and Reliability**: Technical improvements users will notice
- **Iterative Progress**: Building on previous releases and user feedback

### Patch Release Communications
- **Issue Resolution Focus**: Problems that have been fixed
- **Security and Stability**: Critical updates for user protection
- **Quick Deployment**: Streamlined update process
- **Minimal Disruption**: Changes that don't affect user workflows

### Emergency Release Communications
- **Urgency and Importance**: Clear communication of critical nature
- **Immediate Actions**: What users need to do right away
- **Problem Resolution**: Clear explanation of issues being addressed
- **Follow-up Plans**: Additional updates or improvements planned

## Advanced Communication Techniques

### Storytelling in Release Notes
- **User Journey Focus**: Frame changes in terms of user success stories
- **Problem-Solution Narrative**: Describe challenges and how new features address them
- **Progress Demonstration**: Show how the product is evolving toward user goals
- **Community Integration**: Include user feedback and community contributions

### Visual and Interactive Elements
- **Screenshots and Demos**: Show new features in action
- **Video Walkthroughs**: Guided tours of complex new functionality
- **Interactive Guides**: Step-by-step tutorials within the product
- **Comparison Charts**: Before and after feature comparisons

### Personalization and Segmentation
- **Role-Based Content**: Different information for different user types
- **Usage-Based Recommendations**: Relevant features based on user behavior
- **Geographic Customization**: Regional differences and availability
- **Plan-Based Information**: Features available to different subscription levels

## Examples and References

- [[examples/major-release|Major Release Example]]: Comprehensive major version release notes
- [[examples/security-update|Security Update Example]]: Security-focused release communication
- [[examples/api-changes|API Changes Example]]: Developer-focused technical updates
- [[../feature-spec/README|Feature Specifications]]: How release notes relate to implemented features
- [[../test-plan/README|Test Plans]]: How testing validates release note claims

## Need Help?

- For release note best practices: [[README|Release Notes Guide]]
- For feature specification alignment: [[../feature-spec/README|Feature Specification Guide]]
- For test validation: [[../test-plan/README|Test Planning Guide]]
- For architecture context: [[../architecture/README|Architecture Decision Records]]
- For examples: Check the `examples/` directory

Remember: Release notes are the primary communication vehicle between development teams and users. They should be clear, accurate, and user-focused while providing the technical detail necessary for successful adoption of new software versions. Effective release notes build trust, reduce support burden, and accelerate feature adoption. Focus on user benefits, be transparent about issues, and provide clear guidance for any required actions.