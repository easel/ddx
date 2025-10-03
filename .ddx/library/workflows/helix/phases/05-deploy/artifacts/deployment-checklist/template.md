# Deployment Checklist

**Project**: {{ project_name }}  
**Version**: [Version Number]  
**Date**: [Deployment Date]  
**Deploy Lead**: [Name]

## Pre-Deployment

### Code Readiness
- [ ] All tests passing in CI/CD
- [ ] Code coverage meets requirements (>80%)
- [ ] Security scan completed (no critical issues)
- [ ] Code review approved by 2+ reviewers
- [ ] Performance benchmarks met
- [ ] Documentation updated

### Configuration
- [ ] Environment variables documented
- [ ] Secrets rotated and stored securely
- [ ] Feature flags configured
- [ ] Database migrations tested
- [ ] External service credentials verified
- [ ] Rate limits configured

### Dependencies
- [ ] All dependencies version-locked
- [ ] License compliance verified
- [ ] Vulnerability scan passed
- [ ] Breaking changes documented
- [ ] Compatibility matrix updated

### Communication
- [ ] Release notes prepared
- [ ] Customer communication drafted
- [ ] Internal team notified
- [ ] Support team briefed
- [ ] Maintenance window scheduled (if needed)

## Deployment Process

### Staging Deployment
- [ ] Deploy to staging environment
- [ ] Run smoke tests
- [ ] Verify monitoring/alerting
- [ ] Performance testing completed
- [ ] User acceptance testing passed
- [ ] Rollback tested

### Production Deployment

#### Phase 1: Preparation
- [ ] Backup production database
- [ ] Confirm rollback procedure
- [ ] Verify monitoring dashboards
- [ ] Check error budgets
- [ ] Ensure on-call availability

#### Phase 2: Deployment
- [ ] Deploy to canary instance (if applicable)
- [ ] Monitor canary metrics (15 min)
- [ ] Deploy to 10% of production
- [ ] Monitor metrics (30 min)
- [ ] Deploy to 50% of production
- [ ] Monitor metrics (30 min)
- [ ] Deploy to 100% of production
- [ ] Final verification

#### Phase 3: Validation
- [ ] Run production smoke tests
- [ ] Verify critical user journeys
- [ ] Check monitoring dashboards
- [ ] Confirm no increase in errors
- [ ] Validate performance metrics
- [ ] Test rollback capability

## Post-Deployment

### Immediate (0-2 hours)
- [ ] Monitor error rates
- [ ] Check performance metrics
- [ ] Review user feedback channels
- [ ] Verify all services healthy
- [ ] Confirm data integrity

### Short-term (2-24 hours)
- [ ] Analyze usage patterns
- [ ] Review any incidents
- [ ] Check for memory leaks
- [ ] Monitor resource usage
- [ ] Validate business metrics

### Follow-up (24-72 hours)
- [ ] Conduct retrospective
- [ ] Document lessons learned
- [ ] Update runbooks
- [ ] Close deployment ticket
- [ ] Archive deployment artifacts

## Rollback Criteria

Initiate rollback if:
- [ ] Error rate > 5% for 5 minutes
- [ ] P95 latency > 2x baseline
- [ ] Critical functionality broken
- [ ] Data corruption detected
- [ ] Security vulnerability discovered

## Rollback Procedure

1. **Decision**: Incident commander decides rollback
2. **Communication**: Notify all stakeholders
3. **Execute**: Run rollback automation
4. **Verify**: Confirm system restored
5. **Investigate**: Root cause analysis

## Sign-offs

| Role | Name | Approved | Date/Time |
|------|------|----------|-----------|
| Engineering Lead | | ☐ | |
| Product Owner | | ☐ | |
| QA Lead | | ☐ | |
| Operations | | ☐ | |
| Security | | ☐ | |

## Deployment Artifacts

- Build ID: [Build Number]
- Git SHA: [Commit Hash]
- Docker Image: [Image Tag]
- Release Notes: [Link]
- Rollback Build: [Previous Build ID]

## Notes and Issues

### Known Issues
- [Issue 1 and workaround]
- [Issue 2 and mitigation]

### Dependencies Changed
- [Dependency 1: old version → new version]
- [Dependency 2: old version → new version]

### Database Changes
- [Migration 1: description]
- [Migration 2: description]

### Configuration Changes
- [Config 1: old value → new value]
- [Config 2: added/removed]

---
*This checklist must be completed and archived for compliance*