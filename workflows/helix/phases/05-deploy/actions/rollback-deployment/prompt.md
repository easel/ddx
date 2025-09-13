# Rollback Deployment Prompt

Execute an emergency rollback to restore the system to a known good state when deployment issues are detected. This procedure must be fast, reliable, and well-documented.

## Purpose

Rollback procedures ensure:
- Rapid recovery from failed deployments
- Minimal user impact during incidents
- Data integrity is maintained
- Clear communication during crisis
- Lessons learned for prevention

## Rollback Decision Criteria

### Automatic Rollback Triggers
These conditions should trigger immediate automatic rollback:
- Error rate > 5% for 5 minutes
- P95 latency > 2x baseline for 10 minutes
- More than 50% of health checks failing
- Critical business metrics drop > 20%
- Data corruption detected
- Security breach identified

### Manual Rollback Triggers
Human decision needed for:
- User reports of critical bugs
- Partial feature failures
- Performance degradation < automatic thresholds
- External service incompatibilities
- Business decision (market timing, PR issues)

## Rollback Strategies

### 1. Kubernetes Rollback

#### Immediate Rollback to Previous Version
```bash
#!/bin/bash
# rollback.sh - Emergency rollback script

echo "🚨 INITIATING EMERGENCY ROLLBACK 🚨"
echo "Timestamp: $(date)"
echo "Operator: ${USER}"

# Capture current state for post-mortem
kubectl get deployment api -n production -o yaml > deployment-failed-$(date +%Y%m%d-%H%M%S).yaml
kubectl logs -n production -l app=api --tail=1000 > logs-failed-$(date +%Y%m%d-%H%M%S).log

# Perform rollback
echo "Rolling back deployment..."
kubectl rollout undo deployment/api -n production

# Monitor rollback progress
kubectl rollout status deployment/api -n production --timeout=5m

if [ $? -eq 0 ]; then
  echo "✅ Rollback completed successfully"
else
  echo "❌ Rollback failed - escalating to disaster recovery"
  ./disaster-recovery.sh
fi

# Verify health after rollback
sleep 30
if curl -f https://api.example.com/health; then
  echo "✅ Health check passed after rollback"
else
  echo "⚠️ Health check still failing - further investigation needed"
fi
```

#### Rollback to Specific Revision
```bash
# List available revisions
kubectl rollout history deployment/api -n production

# Rollback to specific revision
kubectl rollout undo deployment/api -n production --to-revision=142

# Verify the rollback
kubectl get deployment api -n production -o jsonpath='{.spec.template.spec.containers[0].image}'
```

### 2. Blue-Green Rollback

#### Instant Traffic Switch
```bash
# Current state: Green is failing, Blue is previous stable
echo "Switching traffic back to Blue environment..."

# Immediate switch
kubectl patch service production-lb -n production \
  -p '{"spec":{"selector":{"deployment":"blue"}}}' \
  --type merge

echo "Traffic switched to Blue in < 1 second"

# Verify traffic routing
for i in {1..10}; do
  VERSION=$(curl -s https://api.example.com/version | jq -r '.version')
  echo "Response from: $VERSION"
  sleep 1
done

# Scale down failed Green deployment
kubectl scale deployment api-green -n production --replicas=0
```

### 3. Canary Rollback

#### Progressive Traffic Reduction
```bash
# Gradually reduce canary traffic to 0%
echo "Reducing canary traffic..."

# Step 1: Reduce to 5%
kubectl patch virtualservice api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/http/1/route/1/weight", "value": 5}]'

sleep 60

# Step 2: Reduce to 0%
kubectl patch virtualservice api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/http/1/route/1/weight", "value": 0}]'

# Remove canary deployment
kubectl delete deployment api-canary -n production
```

### 4. Feature Flag Rollback

#### Disable Features Without Deployment
```javascript
// Emergency feature flag disable
async function disableFailingFeature() {
  const flagName = 'new-checkout-flow';

  // Disable immediately for all users
  await featureFlags.update(flagName, {
    enabled: false,
    rollout: { percentage: 0 }
  });

  // Log the action
  logger.critical('Emergency feature flag disabled', {
    flag: flagName,
    timestamp: new Date(),
    operator: process.env.USER
  });

  // Notify team
  await slack.send({
    channel: '#incidents',
    text: `🚨 Feature flag '${flagName}' disabled due to production issues`
  });
}
```

## Database Rollback

### Schema Rollback
```sql
-- Rollback migration script
BEGIN;

-- Verify we're rolling back the right version
SELECT version, applied_at
FROM schema_migrations
ORDER BY version DESC
LIMIT 1;

-- Rollback DDL changes
ALTER TABLE orders DROP COLUMN IF EXISTS new_field;
ALTER TABLE users ALTER COLUMN email TYPE varchar(255);
DROP INDEX IF EXISTS idx_orders_new_field;

-- Restore dropped objects
CREATE TABLE restored_table AS
SELECT * FROM backup.restored_table_backup;

-- Update migration tracking
DELETE FROM schema_migrations
WHERE version = '20240115123456';

-- Verify rollback
SELECT COUNT(*) FROM information_schema.columns
WHERE table_name = 'orders'
AND column_name = 'new_field';
-- Should return 0

COMMIT;
```

### Data Rollback
```bash
# Restore data from backup
pg_restore -h production-db.example.com \
  -U admin \
  -d production \
  -t affected_table \
  backup-production-pre-deploy.dump

# Verify data integrity
psql -h production-db.example.com -U admin -d production <<EOF
SELECT
  COUNT(*) as record_count,
  MAX(updated_at) as last_update
FROM affected_table;
EOF
```

## Communication During Rollback

### Internal Communication Template
```markdown
**🚨 PRODUCTION ROLLBACK IN PROGRESS 🚨**

**Time Detected:** 14:35 UTC
**Time Rollback Started:** 14:37 UTC
**Severity:** P1 - Critical
**Impact:** API errors affecting 15% of users

**Issue:**
- Elevated 500 errors after v2.3.4 deployment
- Payment processing failing for subset of users
- Database connection pool exhaustion

**Actions Taken:**
1. ✅ Initiated automatic rollback at 14:37
2. ✅ Notified on-call team
3. ⏳ Rollback in progress (ETA: 5 minutes)
4. ⏳ Monitoring recovery metrics

**Next Steps:**
- Complete rollback verification
- Root cause analysis
- Customer communication
- Post-mortem scheduling

**War Room:** https://meet.example.com/incident-2024-01-15
**Incident Doc:** https://docs.example.com/incident-2024-01-15
```

### Customer Communication
```markdown
**Service Disruption Notice**

We are currently experiencing technical difficulties with our API service.
Our team has identified the issue and is actively working on a resolution.

**Impact:** Some users may experience errors when placing orders
**Started:** 2:35 PM UTC
**Expected Resolution:** 2:45 PM UTC

We apologize for any inconvenience. Updates will be posted every 15 minutes.

Last Updated: 2:37 PM UTC
```

## Rollback Verification

### Health Checks
```bash
#!/bin/bash
# verify-rollback.sh

echo "Verifying rollback success..."

# Check application version
CURRENT_VERSION=$(curl -s https://api.example.com/version | jq -r '.version')
echo "Current version: $CURRENT_VERSION"

# Check health endpoints
ENDPOINTS=("/health" "/ready" "/metrics")
for endpoint in "${ENDPOINTS[@]}"; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://api.example.com$endpoint)
  if [ $STATUS -eq 200 ]; then
    echo "✅ $endpoint: OK"
  else
    echo "❌ $endpoint: Failed (HTTP $STATUS)"
  fi
done

# Check error rates
ERROR_RATE=$(curl -s https://metrics.example.com/api/v1/query \
  -d 'query=rate(http_requests_total{status=~"5.."}[5m])' \
  | jq -r '.data.result[0].value[1]')

echo "Current error rate: ${ERROR_RATE}%"

# Check key business metrics
ORDER_RATE=$(curl -s https://metrics.example.com/api/v1/query \
  -d 'query=rate(orders_created_total[5m])' \
  | jq -r '.data.result[0].value[1]')

echo "Order rate: ${ORDER_RATE} orders/min"
```

### Smoke Tests After Rollback
```javascript
// Post-rollback validation
async function validateRollback() {
  const tests = [
    {
      name: 'API Health',
      test: async () => {
        const res = await axios.get('/health');
        assert(res.data.status === 'healthy');
      }
    },
    {
      name: 'Authentication',
      test: async () => {
        const res = await axios.post('/auth/login', {
          email: 'test@example.com',
          password: 'test-password'
        });
        assert(res.data.token);
      }
    },
    {
      name: 'Core Functionality',
      test: async () => {
        const res = await axios.get('/api/products');
        assert(Array.isArray(res.data));
      }
    }
  ];

  for (const test of tests) {
    try {
      await test.test();
      console.log(`✅ ${test.name}: PASSED`);
    } catch (error) {
      console.error(`❌ ${test.name}: FAILED - ${error.message}`);
      return false;
    }
  }

  return true;
}
```

## Post-Rollback Actions

### Immediate (First 30 minutes)
1. Confirm system stability
2. Document timeline of events
3. Notify stakeholders of resolution
4. Capture logs and metrics
5. Update status page
6. Begin preliminary root cause analysis

### Short-term (First 24 hours)
1. Conduct incident review meeting
2. Create detailed incident report
3. Identify and document root cause
4. Plan fixes for the issues
5. Update runbooks with learnings
6. Test fixes in staging

### Long-term (First week)
1. Implement permanent fixes
2. Add missing tests/monitoring
3. Conduct blameless post-mortem
4. Share learnings with team
5. Update deployment procedures
6. Plan re-deployment strategy

## Disaster Recovery (If Rollback Fails)

### Escalation Path
```bash
#!/bin/bash
# disaster-recovery.sh

echo "⚠️ DISASTER RECOVERY MODE ACTIVATED ⚠️"

# 1. Switch to DR site
echo "Failing over to disaster recovery site..."
kubectl config use-context dr-cluster

# 2. Update DNS to point to DR
./update-dns-to-dr.sh

# 3. Restore from last known good backup
./restore-from-backup.sh

# 4. Page senior leadership
./page-executives.sh "Production down, DR activated"

# 5. Open bridge line
echo "Open emergency bridge: +1-555-911-HELP"
```

## Rollback Metrics

Track these metrics for every rollback:
- Time to detect issue (MTTD)
- Time to decision (human delay)
- Time to rollback (MTTR)
- User impact duration
- Data loss (if any)
- Revenue impact
- Root cause category

## Prevention Checklist

After rollback, ensure:
- [ ] Root cause identified and fixed
- [ ] Tests added to catch this issue
- [ ] Monitoring added for this scenario
- [ ] Deployment process updated
- [ ] Team trained on lessons learned
- [ ] Runbook updated with this scenario
- [ ] Similar risks identified and mitigated

Remember: A good rollback is fast, safe, and educational. Every rollback is an opportunity to improve your deployment process and system resilience.