# Deploy to Production Prompt

Execute a controlled, monitored deployment to the production environment with appropriate safety measures, progressive rollout, and rollback capabilities.

## Operational Action

This is an operational action that performs critical production deployment activities. It does not generate documentation files but rather executes deployment procedures with safety measures and monitoring.

## Purpose

Production deployment is:
- The culmination of all previous work
- A high-stakes operation requiring precision
- An opportunity to deliver value to users
- A learning experience for future improvements
- A test of your deployment procedures

## Pre-Deployment Requirements

### Go/No-Go Checklist
All items must be checked before proceeding:
- [ ] Staging deployment successful for 24+ hours
- [ ] All stakeholders notified and approved
- [ ] Deployment window confirmed
- [ ] On-call team ready
- [ ] Rollback plan tested in staging
- [ ] Customer support briefed
- [ ] Status page prepared for updates
- [ ] Monitoring dashboards open
- [ ] War room established (if needed)
- [ ] Recent production backup verified

## Deployment Strategies

### 1. Blue-Green Deployment

#### Prepare Green Environment
```bash
# Verify blue is current production
kubectl get service production-lb -o jsonpath='{.spec.selector.deployment}'
# Output: blue

# Deploy to green environment
kubectl apply -f production-green-deployment.yaml

# Wait for green to be ready
kubectl wait --for=condition=available \
  --timeout=600s \
  deployment/api-green -n production

# Run health checks on green
./scripts/health-check.sh https://production-green.internal
```

#### Traffic Cutover
```bash
# Switch load balancer to green
kubectl patch service production-lb -n production \
  -p '{"spec":{"selector":{"deployment":"green"}}}'

# Monitor metrics immediately
watch -n 1 'kubectl top pods -n production -l deployment=green'

# Quick validation
curl -s https://api.example.com/health | jq '.'
```

#### Keep Blue for Rollback
```bash
# Scale blue down but keep ready
kubectl scale deployment api-blue -n production --replicas=2

# Blue remains available for instant rollback
# Can be scaled back up immediately if needed
```

### 2. Canary Deployment

#### Progressive Rollout
```bash
# Stage 1: Deploy canary (5% traffic)
kubectl apply -f production-canary.yaml
kubectl set image deployment/api-canary \
  api=registry.example.com/myapp:v${VERSION} -n production

# Configure traffic split
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: api
  namespace: production
spec:
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: api-canary
      weight: 100
  - route:
    - destination:
        host: api-stable
      weight: 95
    - destination:
        host: api-canary
      weight: 5
EOF
```

#### Monitor Canary Metrics
```javascript
// Monitor canary vs stable
const monitorCanary = async () => {
  const canaryMetrics = await getMetrics('deployment=canary');
  const stableMetrics = await getMetrics('deployment=stable');

  const comparison = {
    errorRate: {
      canary: canaryMetrics.errorRate,
      stable: stableMetrics.errorRate,
      threshold: stableMetrics.errorRate * 1.1 // 10% tolerance
    },
    latency: {
      canary: canaryMetrics.p95Latency,
      stable: stableMetrics.p95Latency,
      threshold: stableMetrics.p95Latency * 1.2 // 20% tolerance
    }
  };

  if (comparison.errorRate.canary > comparison.errorRate.threshold) {
    console.error('Canary error rate exceeds threshold!');
    await rollbackCanary();
  }

  if (comparison.latency.canary > comparison.latency.threshold) {
    console.error('Canary latency exceeds threshold!');
    await rollbackCanary();
  }

  return comparison;
};
```

#### Progressive Increase
```bash
# Stage 2: Increase to 25% (after 30 minutes)
kubectl patch virtualservice api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/http/1/route/1/weight", "value": 25}]'

# Stage 3: Increase to 50% (after 1 hour)
kubectl patch virtualservice api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/http/1/route/1/weight", "value": 50}]'

# Stage 4: Full deployment (after 2 hours)
kubectl patch virtualservice api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/http/1/route/1/weight", "value": 100}]'
```

### 3. Feature Flag Deployment

#### Deploy with Features Disabled
```javascript
// Deploy code with feature behind flag
if (featureFlags.isEnabled('new-checkout-flow')) {
  return renderNewCheckout();
} else {
  return renderOldCheckout();
}
```

#### Progressive Feature Enablement
```bash
# Enable for internal users
curl -X POST https://api.launchdarkly.com/flags/new-checkout-flow \
  -H "Authorization: ${LD_API_KEY}" \
  -d '{"variations": [{"targets": ["internal-users"], "variation": 0}]}'

# Enable for 10% of users
curl -X POST https://api.launchdarkly.com/flags/new-checkout-flow \
  -H "Authorization: ${LD_API_KEY}" \
  -d '{"rollout": {"variations": [{"variation": 0, "weight": 10000}]}}'

# Monitor and increase gradually
```

## Deployment Execution

### 1. Final Pre-Flight Checks

```bash
# Verify current production state
./scripts/production-health-check.sh

# Check resource availability
kubectl describe resourcequota -n production
kubectl top nodes -l node-role.kubernetes.io/production=true

# Verify external dependencies
for service in payment-gateway email-service cdn analytics; do
  curl -f https://$service.example.com/health || exit 1
done

# Confirm deployment image
docker pull registry.example.com/myapp:v${VERSION}
docker inspect registry.example.com/myapp:v${VERSION} | jq '.[0].Config.Labels'
```

### 2. Database Migration (if needed)

```bash
# Take backup before migration
pg_dump -h production-db.example.com -U admin -d production \
  --no-owner --no-privileges \
  -f backup-production-pre-deploy-$(date +%Y%m%d-%H%M%S).sql

# Run migration with transaction
psql -h production-db.example.com -U admin -d production <<EOF
BEGIN;
-- Migration SQL here
SELECT * FROM verify_migration();
COMMIT;
EOF

# Verify migration success
psql -h production-db.example.com -U admin -d production \
  -c "SELECT version, checksum FROM schema_migrations ORDER BY version DESC LIMIT 1;"
```

### 3. Execute Deployment

```bash
# Record deployment start
echo "Deployment started at $(date)" >> deployment.log

# Apply production configuration
kubectl apply -f production-config.yaml

# Deploy application
kubectl set image deployment/api \
  api=registry.example.com/myapp:v${VERSION} \
  -n production --record

# Monitor rollout
kubectl rollout status deployment/api -n production --timeout=10m

# Verify all pods running
kubectl get pods -n production -l app=api \
  -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.phase}{"\n"}{end}'
```

### 4. Post-Deployment Validation

#### Immediate Checks (0-5 minutes)
```bash
# Health endpoints
curl -f https://api.example.com/health
curl -f https://api.example.com/ready

# Version verification
DEPLOYED_VERSION=$(curl -s https://api.example.com/version | jq -r '.version')
if [ "$DEPLOYED_VERSION" != "v${VERSION}" ]; then
  echo "ERROR: Version mismatch!"
  ./scripts/rollback.sh
fi

# Critical path testing
./scripts/production-smoke-tests.sh
```

#### Extended Monitoring (5-30 minutes)
```javascript
// Monitor key metrics
const monitoring = setInterval(async () => {
  const metrics = await getProductionMetrics();

  console.log(`
    Error Rate: ${metrics.errorRate}% (threshold: 1%)
    P95 Latency: ${metrics.p95}ms (threshold: 500ms)
    Throughput: ${metrics.rps} req/s
    CPU Usage: ${metrics.cpu}%
    Memory: ${metrics.memory}MB
  `);

  if (metrics.errorRate > 1 || metrics.p95 > 500) {
    console.error('Metrics exceeded thresholds!');
    clearInterval(monitoring);
    await initiateRollback();
  }
}, 60000); // Check every minute

// Stop monitoring after 30 minutes
setTimeout(() => clearInterval(monitoring), 30 * 60 * 1000);
```

## Monitoring During Deployment

### Real-time Dashboards
Monitor these dashboards during deployment:
- Application metrics (errors, latency, throughput)
- Infrastructure metrics (CPU, memory, network)
- Business metrics (orders, sign-ups, revenue)
- User experience metrics (page load, interactions)
- External service status

### Alert Response
```yaml
# Critical alerts during deployment
alerts:
  - name: deployment_high_error_rate
    threshold: "> 5%"
    action: Immediate rollback

  - name: deployment_latency_spike
    threshold: "p95 > 2x baseline"
    action: Investigate, consider rollback

  - name: deployment_pod_crashes
    threshold: "> 2 restarts"
    action: Check logs, rollback if pattern

  - name: deployment_traffic_drop
    threshold: "< 50% of normal"
    action: Check routing, possible rollback
```

## Rollback Procedures

### Automatic Rollback Triggers
```bash
# Automated rollback script
#!/bin/bash
rollback_if_unhealthy() {
  local max_attempts=5
  local attempt=0

  while [ $attempt -lt $max_attempts ]; do
    if curl -f https://api.example.com/health; then
      echo "Health check passed"
      return 0
    fi
    attempt=$((attempt + 1))
    sleep 10
  done

  echo "Health checks failed. Initiating rollback..."
  kubectl rollout undo deployment/api -n production
  ./scripts/notify-rollback.sh
}
```

### Manual Rollback
```bash
# Immediate rollback to previous version
kubectl rollout undo deployment/api -n production

# Rollback to specific revision
kubectl rollout undo deployment/api -n production --to-revision=142

# For blue-green: switch back to blue
kubectl patch service production-lb -n production \
  -p '{"spec":{"selector":{"deployment":"blue"}}}'
```

## Communication During Deployment

### Status Updates
```markdown
**Production Deployment Status**

⏰ Start Time: 14:00 UTC
📦 Version: v2.3.4
👥 Deploy Lead: @johndoe

**Timeline:**
- 14:00 - Deployment initiated
- 14:05 - ✅ 25% of servers updated
- 14:10 - ✅ 50% of servers updated
- 14:15 - ✅ 75% of servers updated
- 14:20 - ✅ 100% deployed
- 14:25 - ✅ Health checks passing
- 14:30 - ✅ Monitoring stable

**Metrics:**
- Error Rate: 0.02% (normal)
- P95 Latency: 245ms (normal)
- Throughput: 1,250 req/s

**Status: SUCCESS** ✅
```

### Incident Communication
If issues arise:
1. Update status page immediately
2. Notify on-call team
3. Open war room
4. Send stakeholder updates every 15 minutes
5. Document timeline and decisions

## Post-Deployment Tasks

### Immediate (First Hour)
- [ ] Monitor dashboards continuously
- [ ] Review error logs
- [ ] Check customer support tickets
- [ ] Verify feature flags working
- [ ] Test critical user paths
- [ ] Update status page

### Short-term (First Day)
- [ ] Analyze performance metrics
- [ ] Review user feedback
- [ ] Check for any degradation
- [ ] Verify backups running
- [ ] Update documentation
- [ ] Schedule retrospective

### Long-term (First Week)
- [ ] Conduct deployment retrospective
- [ ] Update runbooks
- [ ] Analyze user adoption
- [ ] Performance optimization
- [ ] Security scan
- [ ] Plan next iteration

## Success Metrics

Deployment is successful when:
- Zero downtime achieved
- Error rate remains below 0.1%
- P95 latency within 10% of baseline
- No critical alerts fired
- No emergency rollback needed
- Positive user feedback
- Business metrics stable or improved

Remember: Production deployment is not just a technical exercise—it's about delivering value to users safely and reliably. Take your time, follow the process, and never skip steps under pressure.