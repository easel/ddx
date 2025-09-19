# Runbook Generation Prompt

Create a comprehensive operational runbook that provides step-by-step procedures for deployment, maintenance, troubleshooting, and emergency response throughout the application lifecycle.

## Storage Location

Store the runbook at: `docs/helix/05-deploy/runbook.md`

## Purpose

The runbook serves as:
- The single source of truth for operations
- A guide for on-call engineers
- Training material for new team members
- Documentation for compliance
- A living document that evolves with the system

## Key Sections

### 1. System Overview

#### Architecture Summary
Provide context for operators:
- System components and their roles
- Data flow between services
- External dependencies
- Critical vs. non-critical services
- Service level agreements

#### Access Information
Document how to access systems:
- Production URLs and endpoints
- Admin panel locations
- Database connection strings
- SSH/RDP access procedures
- VPN requirements
- Required credentials/keys location

### 2. Deployment Procedures

#### Standard Deployment
Step-by-step for normal releases:

```bash
# 1. Pre-deployment checks
./scripts/pre-deploy-check.sh

# 2. Create deployment branch
git checkout -b deploy/v1.2.3

# 3. Run deployment
./scripts/deploy.sh production v1.2.3

# 4. Verify deployment
./scripts/verify-deployment.sh

# 5. Update status page
./scripts/update-status.sh deployed
```

#### Emergency Hotfix
Rapid deployment for critical fixes:
1. Assess severity and impact
2. Create hotfix branch from production
3. Apply minimal fix
4. Fast-track testing
5. Deploy with abbreviated process
6. Full testing post-deployment

#### Rollback Procedures
When things go wrong:
1. Identify rollback trigger
2. Notify stakeholders
3. Execute rollback script
4. Verify system stability
5. Investigate root cause
6. Document incident

### 3. Monitoring and Alerts

#### Dashboard Guide
Where to look for what:
- **System Health**: Main operations dashboard
- **Performance**: Response time and throughput graphs
- **Business Metrics**: Revenue and user activity
- **Errors**: Error rate and log aggregation
- **Infrastructure**: Resource utilization

#### Alert Response
For each alert type:

```yaml
alert: high_error_rate
severity: critical
dashboard: https://monitoring.example.com/errors
runbook_section: "#high-error-rate"
first_response:
  - Check error logs for patterns
  - Identify affected services
  - Check recent deployments
  - Assess user impact
resolution:
  - Rollback if deployment-related
  - Scale services if load-related
  - Fix and deploy if code issue
  - Update status page
```

### 4. Troubleshooting Guide

#### Common Issues

##### Service Won't Start
**Symptoms**: Service fails to initialize
**Diagnosis**:
1. Check logs for startup errors
2. Verify configuration files
3. Check database connectivity
4. Validate required services running
5. Confirm resource availability

**Resolution**:
- Fix configuration errors
- Restart dependent services
- Clear corrupted cache
- Increase resource limits
- Restore from backup if needed

##### High Latency
**Symptoms**: Slow response times
**Diagnosis**:
1. Check current load vs. capacity
2. Review slow query logs
3. Analyze trace data
4. Check external service latency
5. Review recent code changes

**Resolution**:
- Scale horizontally
- Optimize database queries
- Increase cache usage
- Implement circuit breakers
- Throttle non-critical features

##### Data Inconsistency
**Symptoms**: Data mismatch across systems
**Diagnosis**:
1. Identify affected data scope
2. Check replication status
3. Review transaction logs
4. Verify data pipeline status
5. Check for failed jobs

**Resolution**:
- Run reconciliation scripts
- Replay failed transactions
- Restore from consistent backup
- Fix and re-run ETL jobs
- Implement data validation

### 5. Maintenance Procedures

#### Routine Maintenance

##### Database Maintenance
```sql
-- Weekly maintenance
VACUUM ANALYZE;
REINDEX DATABASE production;

-- Monthly maintenance
-- Archive old data
INSERT INTO archive.orders
SELECT * FROM orders
WHERE created_at < NOW() - INTERVAL '90 days';

DELETE FROM orders
WHERE created_at < NOW() - INTERVAL '90 days';
```

##### Certificate Renewal
1. Check expiration dates (30 days ahead)
2. Generate new certificates
3. Test in staging environment
4. Deploy during maintenance window
5. Verify all services using new certs
6. Update monitoring for new expiry

##### Dependency Updates
1. Review security advisories
2. Test updates in development
3. Run full test suite
4. Deploy to staging
5. Monitor for 24 hours
6. Deploy to production

### 6. Disaster Recovery

#### Backup Procedures
- **Frequency**: Database (hourly), Files (daily)
- **Retention**: 30 days standard, 1 year for monthly
- **Testing**: Monthly restore test
- **Location**: Multi-region storage
- **Encryption**: AES-256 at rest

#### Recovery Scenarios

##### Complete System Failure
**RTO**: 4 hours | **RPO**: 1 hour

1. Declare disaster recovery mode
2. Activate DR team
3. Provision infrastructure in DR region
4. Restore database from backup
5. Deploy application code
6. Update DNS to DR endpoints
7. Verify system functionality
8. Communicate status to stakeholders

##### Data Corruption
**RTO**: 2 hours | **RPO**: Point before corruption

1. Identify corruption timestamp
2. Stop writes to affected systems
3. Restore from clean backup
4. Replay transactions post-backup
5. Validate data integrity
6. Resume normal operations

### 7. Security Response

#### Security Incident Response
1. **Detect**: Identify potential breach
2. **Contain**: Isolate affected systems
3. **Investigate**: Determine scope and impact
4. **Eradicate**: Remove threat
5. **Recover**: Restore normal operations
6. **Review**: Post-incident analysis

#### Common Security Tasks
- Rotate compromised credentials
- Block malicious IPs
- Apply security patches
- Review access logs
- Update firewall rules
- Implement additional monitoring

### 8. Communication Procedures

#### Incident Communication

##### Internal Communication
```
INCIDENT DETECTED
Time: [timestamp]
Severity: [P1/P2/P3]
Impact: [user-facing description]
Team: @oncall @engineering
Status Page: [updating/updated]
War Room: [link to call/chat]
```

##### Customer Communication
- Initial acknowledgment (within 15 minutes)
- Status updates (every 30 minutes)
- Resolution notification
- Post-incident report (within 48 hours)

#### Escalation Matrix
| Time | Level | Contact |
|------|-------|---------|
| 0 min | L1 - On-call Engineer | PagerDuty |
| 15 min | L2 - Team Lead | Phone |
| 30 min | L3 - Engineering Manager | Phone |
| 1 hour | L4 - VP Engineering | Phone |
| 2 hours | L5 - CTO | Phone |

### 9. Performance Tuning

#### Application Optimization
- Enable query caching
- Implement connection pooling
- Optimize database indexes
- Use CDN for static assets
- Enable compression
- Implement lazy loading

#### Infrastructure Scaling
- Horizontal scaling triggers
- Vertical scaling limits
- Auto-scaling configuration
- Load balancer tuning
- Cache sizing
- Queue optimization

### 10. Appendices

#### A. Important Commands
```bash
# View logs
kubectl logs -f deployment/api

# Scale deployment
kubectl scale deployment/api --replicas=10

# Database console
psql -h db.example.com -U admin -d production

# Clear cache
redis-cli FLUSHALL

# Force deployment
./deploy.sh --force --skip-tests production
```

#### B. Contact Information
- On-call rotation: [schedule link]
- Vendor support: [contact list]
- Management chain: [phone numbers]
- Customer success: [contact info]
- Legal/Compliance: [contact info]

#### C. External Dependencies
- Payment Gateway: [status page, support]
- Email Service: [status page, support]
- CDN Provider: [status page, support]
- Cloud Provider: [status page, support]
- Monitoring Service: [status page, support]

## Runbook Maintenance

### Update Triggers
- After each incident
- Following system changes
- During quarterly review
- When team members rotate
- After automation improvements

### Version Control
- Store in git repository
- Tag versions with deployments
- Review changes in PRs
- Maintain changelog
- Archive old versions

## Quality Checklist

Before runbook is complete:
- [ ] All procedures tested
- [ ] Commands verified to work
- [ ] Links are valid
- [ ] Contact info current
- [ ] Diagrams included
- [ ] Examples provided
- [ ] Index/TOC created
- [ ] Searchable format
- [ ] Team reviewed
- [ ] Accessible during outages

Remember: A runbook is only as good as its last update. Keep it current, keep it tested, and keep it accessible when everything else is on fire.