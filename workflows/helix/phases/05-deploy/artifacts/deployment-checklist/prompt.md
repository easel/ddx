# Deployment Checklist Generation Prompt

Create a comprehensive deployment checklist that ensures all prerequisites are met and procedures are followed for a safe, successful deployment to production.

## Storage Location

Store the deployment checklist at: `docs/helix/05-deploy/deployment-checklist.md`

## Purpose

The deployment checklist serves as:
- A go/no-go decision framework
- A systematic validation of readiness
- A communication tool for stakeholders
- A record of deployment preparation
- A risk mitigation document

## Key Requirements

### 1. Pre-Deployment Validation

#### Code Readiness
- Test results summary (unit, integration, E2E)
- Code coverage metrics
- Security scan results
- Performance test results
- Code review status
- Static analysis results

#### Configuration Readiness
- Environment variables documented and set
- Secrets rotated and securely stored
- Feature flags configured
- Database migrations tested
- API keys and certificates valid
- Rate limits and quotas configured

#### Infrastructure Readiness
- Server capacity verified
- Load balancers configured
- CDN cache rules set
- DNS records prepared
- SSL certificates valid
- Backup systems operational

### 2. Dependency Management

#### Application Dependencies
- All dependencies version-locked
- Security vulnerabilities addressed
- License compliance verified
- Breaking changes documented
- Compatibility matrix updated
- Dependency update notes

#### External Service Dependencies
- Third-party service status checked
- API rate limits verified
- Failover mechanisms tested
- Service accounts active
- Integration points validated
- SLA requirements met

### 3. Communication Plan

#### Internal Communication
- Development team notified
- Operations team briefed
- Support team prepared
- Management informed
- Deployment schedule shared
- Escalation paths defined

#### External Communication
- Customer notification drafted
- Status page updated
- Release notes prepared
- Documentation updated
- Support articles ready
- Maintenance window scheduled (if needed)

### 4. Deployment Process

#### Staging Validation
- Staging deployment successful
- Smoke tests passing
- Integration tests passing
- Performance acceptable
- User acceptance complete
- Rollback tested

#### Production Preparation
- Production backup completed
- Database migrations ready
- Cache warming planned
- Traffic routing prepared
- Monitoring enhanced
- Incident response ready

### 5. Rollback Planning

#### Rollback Triggers
Define conditions that trigger rollback:
- Error rate threshold exceeded
- Response time degradation
- Critical functionality broken
- Data corruption detected
- Security breach identified
- Business metrics impacted

#### Rollback Procedures
- Rollback steps documented
- Time to rollback estimated
- Data migration reversal planned
- Communication plan ready
- Post-rollback validation defined
- Lessons learned process ready

### 6. Risk Assessment

#### Technical Risks
- Identify potential failure points
- Assess impact of each risk
- Define mitigation strategies
- Prepare contingency plans
- Document decision criteria
- Establish risk tolerance

#### Business Risks
- Customer impact assessment
- Revenue impact analysis
- Regulatory compliance check
- Brand reputation considerations
- Competitive implications
- Market timing factors

## Checklist Structure

### Section Organization
1. **Header**: Version, date, deploy lead, approvers
2. **Executive Summary**: Key changes and impacts
3. **Pre-Deployment**: All validation checks
4. **Deployment Steps**: Detailed procedures
5. **Validation**: Post-deployment checks
6. **Rollback**: Emergency procedures
7. **Sign-offs**: Required approvals

### Item Format
Each checklist item should include:
- Clear, actionable description
- Responsible party
- Completion status checkbox
- Verification method
- Time estimate
- Dependencies

## Validation Requirements

### Automated Checks
Items that can be automatically verified:
- Test suite execution
- Security scans
- Performance benchmarks
- Configuration validation
- Certificate expiration
- Service health checks

### Manual Checks
Items requiring human verification:
- Business logic validation
- User experience testing
- Documentation review
- Communication approval
- Risk assessment
- Go/no-go decision

## Timeline Considerations

### Deployment Window
- Preferred deployment time
- Duration estimate
- Buffer time included
- Rollback window
- Post-deployment monitoring period
- Business hours consideration

### Critical Periods
Identify times to avoid deployment:
- Peak traffic periods
- End of month/quarter
- Holiday seasons
- Marketing campaigns
- Major events
- Maintenance windows

## Success Criteria

Define what constitutes successful deployment:
- All tests passing
- Error rates below threshold
- Response times within SLA
- Key features functional
- No critical alerts
- Business metrics stable

## Post-Deployment Actions

### Immediate Actions (0-1 hour)
- Smoke tests executed
- Monitoring verified
- Key metrics checked
- Team notifications sent
- Initial user feedback gathered
- Quick wins documented

### Short-term Actions (1-24 hours)
- Extended monitoring
- Performance analysis
- Error log review
- User feedback collection
- Support ticket monitoring
- Metrics trend analysis

### Long-term Actions (1-7 days)
- Full metrics analysis
- User satisfaction survey
- Performance optimization
- Documentation updates
- Lessons learned session
- Next iteration planning

## Quality Checklist

Before finalizing the deployment checklist:
- [ ] All critical systems covered
- [ ] Dependencies clearly mapped
- [ ] Rollback procedures tested
- [ ] Communication plan complete
- [ ] Risk assessment thorough
- [ ] Success criteria measurable
- [ ] Timeline realistic
- [ ] Responsibilities assigned
- [ ] Approval chain defined
- [ ] Documentation current

## Integration with Deploy Phase

This checklist enables the Deploy phase by:
1. Providing systematic validation
2. Ensuring nothing is forgotten
3. Facilitating clear communication
4. Enabling informed decisions
5. Reducing deployment risks

Remember: A good deployment checklist turns a stressful deployment into a routine procedure. Every item should add value and reduce risk.