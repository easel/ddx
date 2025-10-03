# Production Runbook: {{ project_name }}

**Last Updated**: [Date]  
**Maintained By**: [Team Name]  
**On-Call Rotation**: [Link to schedule]

## Service Overview

### Description
[Brief description of what this service does]

### Business Impact
- **Critical Functions**: [What breaks if this service is down]
- **Affected Users**: [Who is impacted]
- **Revenue Impact**: [Financial impact per hour of downtime]

### Architecture
```
[Simple ASCII or mermaid diagram showing key components]
```

### Dependencies
- **Upstream**: [Services that call this service]
- **Downstream**: [Services this service calls]
- **Databases**: [Database dependencies]
- **External**: [Third-party services]

## Key Metrics and Alerts

### Golden Signals
1. **Latency**: p50 < 100ms, p99 < 1s
2. **Traffic**: Normal range 100-500 RPS
3. **Errors**: Error rate < 0.1%
4. **Saturation**: CPU < 70%, Memory < 80%

### Dashboard Links
- [Operations Dashboard](link)
- [Business Metrics](link)
- [Error Tracking](link)
- [APM/Tracing](link)

## Common Issues and Solutions

### Issue: High Latency
**Symptoms**: 
- P95 latency > 1s
- Increasing queue depth
- User complaints about slowness

**Diagnosis**:
1. Check database query performance
2. Review cache hit rates
3. Look for blocking I/O operations
4. Check for memory pressure/GC

**Resolution**:
```bash
# 1. Increase cache TTL temporarily
kubectl set env deployment/api CACHE_TTL=3600

# 2. Scale up instances
kubectl scale deployment/api --replicas=10

# 3. If database is slow, failover to replica
# See database runbook for procedure
```

### Issue: Service Not Responding
**Symptoms**:
- Health checks failing
- 503 errors from load balancer
- No metrics being reported

**Diagnosis**:
1. Check pod status: `kubectl get pods -l app=api`
2. Review recent deployments
3. Check for OOM kills: `kubectl describe pod [pod-name]`
4. Review error logs

**Resolution**:
```bash
# 1. Quick restart (if safe)
kubectl rollout restart deployment/api

# 2. Rollback if recent deployment
kubectl rollout undo deployment/api

# 3. Scale horizontally if resource constrained
kubectl scale deployment/api --replicas=15
```

### Issue: Database Connection Errors
**Symptoms**:
- "Connection pool exhausted" errors
- Timeouts on database operations
- Sporadic 500 errors

**Diagnosis**:
1. Check connection pool metrics
2. Review database CPU/connections
3. Look for slow queries
4. Check for deadlocks

**Resolution**:
```bash
# 1. Increase connection pool size
kubectl set env deployment/api DB_POOL_SIZE=50

# 2. Kill long-running queries (use with caution)
# Connect to database and run:
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE state = 'active' AND query_time > interval '5 minutes';

# 3. Temporary read traffic diversion to replica
kubectl set env deployment/api DB_READ_HOST=replica.db.example.com
```

### Issue: Memory Leak
**Symptoms**:
- Gradually increasing memory usage
- Eventually OOM kills
- Performance degradation over time

**Diagnosis**:
1. Review memory graphs over 24h period
2. Check for correlation with traffic
3. Capture heap dump if possible
4. Review recent code changes

**Resolution**:
```bash
# 1. Immediate mitigation - rolling restart
kubectl rollout restart deployment/api

# 2. Increase memory limits temporarily
kubectl set resources deployment/api --limits=memory=4Gi

# 3. Enable aggressive GC (Java example)
kubectl set env deployment/api JAVA_OPTS="-XX:+UseG1GC -XX:MaxGCPauseMillis=100"

# 4. Capture diagnostics for debugging
kubectl exec [pod-name] -- jmap -dump:format=b,file=/tmp/heapdump.hprof 1
kubectl cp [pod-name]:/tmp/heapdump.hprof ./heapdump.hprof
```

## Maintenance Procedures

### Deployment
```bash
# Standard deployment
./deploy.sh production v1.2.3

# Canary deployment
./deploy.sh production v1.2.3 --canary 10%

# Emergency hotfix
./deploy.sh production v1.2.3-hotfix --skip-tests --emergency
```

### Rollback
```bash
# Rollback to previous version
kubectl rollout undo deployment/api

# Rollback to specific version
kubectl rollout undo deployment/api --to-revision=42

# View rollout history
kubectl rollout history deployment/api
```

### Scaling
```bash
# Manual scaling
kubectl scale deployment/api --replicas=20

# Enable autoscaling
kubectl autoscale deployment/api --min=5 --max=50 --cpu-percent=70

# Vertical scaling (requires restart)
kubectl set resources deployment/api --requests=memory=2Gi,cpu=1000m
```

### Database Maintenance
```bash
# Run migrations
kubectl exec -it deploy/migration-job -- npm run migrate

# Backup database
kubectl exec -it deploy/backup-job -- ./backup.sh

# Analyze/vacuum (PostgreSQL)
kubectl exec -it [db-pod] -- psql -c "VACUUM ANALYZE;"
```

## Emergency Procedures

### Full Service Outage
1. **Declare incident** in #incidents channel
2. **Page secondary on-call** if needed
3. **Check recent changes** (deployments, configs, flags)
4. **Verify dependencies** are healthy
5. **Attempt quick restart** if safe
6. **Rollback** if recent deployment suspected
7. **Engage vendor support** if external dependency issue

### Data Corruption
1. **STOP all writes immediately**
   ```bash
   kubectl scale deployment/api --replicas=0
   ```
2. **Backup current state** before any fixes
3. **Identify scope** of corruption
4. **Restore from backup** if necessary
5. **Validate data integrity** before resuming

### Security Incident
1. **Isolate affected systems**
2. **Preserve evidence** (logs, memory dumps)
3. **Notify security team** immediately
4. **Follow security incident playbook**
5. **Do not attempt cleanup** without security approval

## Contact Information

### Escalation Path
1. Primary On-Call: [PagerDuty rotation]
2. Secondary On-Call: [PagerDuty rotation]
3. Team Lead: [Name] - [Phone]
4. Director: [Name] - [Phone]
5. VP Engineering: [Name] - [Phone]

### Vendor Support
- **Cloud Provider**: [Support Portal](link) | [Phone]
- **Database Vendor**: [Support Portal](link) | [Phone]
- **CDN Provider**: [Support Portal](link) | [Phone]
- **Monitoring Vendor**: [Support Portal](link) | [Phone]

### Internal Teams
- **Security**: #security-team | [On-call phone]
- **Database**: #database-team | [On-call phone]
- **Network**: #network-team | [On-call phone]
- **Platform**: #platform-team | [On-call phone]

## Disaster Recovery

### Backup Strategy
- **Frequency**: Every 6 hours
- **Retention**: 30 days
- **Location**: [Backup storage location]
- **Recovery Time Objective (RTO)**: 4 hours
- **Recovery Point Objective (RPO)**: 6 hours

### DR Procedure
1. **Assess damage** and determine recovery path
2. **Notify stakeholders** of expected downtime
3. **Provision new infrastructure** if needed
4. **Restore from backup**
5. **Validate data integrity**
6. **Update DNS/routing**
7. **Verify full functionality**
8. **Document incident** and lessons learned

## Appendix

### Configuration Files
- Production config: `/configs/production.yaml`
- Environment variables: `/configs/.env.production`
- Feature flags: [Feature flag service](link)

### Log Locations
- Application logs: [Log aggregator query](link)
- Access logs: [Log aggregator query](link)
- Error logs: [Log aggregator query](link)
- Audit logs: [Log aggregator query](link)

### Useful Commands
```bash
# Get pod logs
kubectl logs -f deployment/api --tail=100

# Execute command in pod
kubectl exec -it [pod-name] -- /bin/bash

# Port forward for debugging
kubectl port-forward [pod-name] 8080:8080

# Check resource usage
kubectl top pods -l app=api

# Get events
kubectl get events --sort-by='.lastTimestamp'
```

---
*Keep this runbook updated. Last review: [Date]*