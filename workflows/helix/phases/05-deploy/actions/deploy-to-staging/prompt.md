# Deploy to Staging Prompt

Execute a complete deployment to the staging environment, including building, packaging, deploying, and validating the application in a production-like setting.

## Purpose

Staging deployment serves as:
- Final validation before production
- Performance testing environment
- User acceptance testing platform
- Integration verification
- Rollback procedure testing

## Deployment Steps

### 1. Pre-Deployment Preparation

#### Environment Verification
```bash
# Check staging environment status
kubectl get nodes -l env=staging
kubectl get pods -n staging

# Verify database connectivity
psql -h staging-db.example.com -U app_user -d staging -c "SELECT version();"

# Check external service connectivity
curl -I https://staging-api.payment-provider.com/health
curl -I https://staging.email-service.com/status

# Verify required secrets exist
kubectl get secrets -n staging | grep -E "db-password|api-keys|certificates"
```

#### Build Artifacts
```bash
# Build application
npm run build:staging
# or
go build -tags staging -o app-staging

# Build Docker image
docker build -t myapp:staging-${VERSION} \
  --build-arg ENV=staging \
  --build-arg VERSION=${VERSION} .

# Push to registry
docker tag myapp:staging-${VERSION} registry.example.com/myapp:staging-${VERSION}
docker push registry.example.com/myapp:staging-${VERSION}
```

### 2. Database Migration

#### Pre-Migration Backup
```bash
# Backup current database
pg_dump -h staging-db.example.com -U admin -d staging \
  -f backup-staging-$(date +%Y%m%d-%H%M%S).sql

# Verify backup
pg_restore --list backup-staging-*.sql | head -20
```

#### Run Migrations
```bash
# Dry run first
npm run migrate:staging -- --dry-run

# Execute migrations
npm run migrate:staging

# Verify migration success
psql -h staging-db.example.com -U admin -d staging \
  -c "SELECT version, applied_at FROM schema_migrations ORDER BY version DESC LIMIT 5;"
```

### 3. Deployment Execution

#### Kubernetes Deployment
```yaml
# staging-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: staging
spec:
  replicas: 2
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: api
      env: staging
  template:
    metadata:
      labels:
        app: api
        env: staging
        version: ${VERSION}
    spec:
      containers:
      - name: api
        image: registry.example.com/myapp:staging-${VERSION}
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: staging
        - name: VERSION
          value: ${VERSION}
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

Execute deployment:
```bash
# Apply deployment
kubectl apply -f staging-deployment.yaml

# Watch rollout status
kubectl rollout status deployment/api -n staging

# Verify pods are running
kubectl get pods -n staging -l app=api
```

#### Blue-Green Deployment Alternative
```bash
# Deploy to green environment
kubectl apply -f staging-green-deployment.yaml

# Run smoke tests against green
./scripts/smoke-test.sh https://staging-green.example.com

# Switch traffic to green
kubectl patch service api -n staging \
  -p '{"spec":{"selector":{"color":"green"}}}'

# Keep blue for quick rollback
kubectl scale deployment api-blue -n staging --replicas=0
```

### 4. Post-Deployment Validation

#### Health Checks
```bash
# Check application health
curl -f https://staging.example.com/health || exit 1

# Check all endpoints
for endpoint in /api/users /api/products /api/orders; do
  response=$(curl -s -o /dev/null -w "%{http_code}" https://staging.example.com$endpoint)
  if [ $response -ne 200 ]; then
    echo "ERROR: $endpoint returned $response"
    exit 1
  fi
done

# Verify version deployed
deployed_version=$(curl -s https://staging.example.com/version | jq -r '.version')
if [ "$deployed_version" != "$VERSION" ]; then
  echo "ERROR: Wrong version deployed. Expected $VERSION, got $deployed_version"
  exit 1
fi
```

#### Smoke Tests
```javascript
// staging-smoke-tests.js
const axios = require('axios');
const assert = require('assert');

const BASE_URL = 'https://staging.example.com';

async function runSmokeTests() {
  console.log('Running staging smoke tests...');

  // Test authentication
  const loginResponse = await axios.post(`${BASE_URL}/api/auth/login`, {
    email: 'test@example.com',
    password: 'test-password'
  });
  assert(loginResponse.data.token, 'Login should return token');

  const token = loginResponse.data.token;
  const headers = { Authorization: `Bearer ${token}` };

  // Test core functionality
  const productsResponse = await axios.get(`${BASE_URL}/api/products`, { headers });
  assert(Array.isArray(productsResponse.data), 'Products should be an array');

  // Test database connectivity
  const userResponse = await axios.get(`${BASE_URL}/api/users/me`, { headers });
  assert(userResponse.data.email, 'Should retrieve user data');

  // Test external service integration
  const paymentTest = await axios.post(`${BASE_URL}/api/payments/test`, {
    amount: 100,
    currency: 'USD'
  }, { headers });
  assert(paymentTest.data.success, 'Payment test should succeed');

  console.log('✅ All smoke tests passed!');
}

runSmokeTests().catch(err => {
  console.error('❌ Smoke tests failed:', err.message);
  process.exit(1);
});
```

### 5. Performance Validation

#### Load Testing
```javascript
// k6 load test for staging
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 50 },  // Ramp up
    { duration: '5m', target: 50 },  // Sustain
    { duration: '2m', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.1'],
  },
};

export default function () {
  const response = http.get('https://staging.example.com/api/products');

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
```

### 6. Monitoring Setup

#### Configure Alerts
```yaml
# staging-alerts.yaml
alerts:
  - name: staging_high_error_rate
    expr: rate(http_requests_total{env="staging",status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: warning
      env: staging
    annotations:
      summary: "High error rate in staging"

  - name: staging_slow_response
    expr: histogram_quantile(0.95, http_request_duration_seconds{env="staging"}) > 1
    for: 10m
    labels:
      severity: warning
      env: staging
```

#### Verify Monitoring
```bash
# Check metrics are being collected
curl -s https://staging-metrics.example.com/api/v1/query \
  -d 'query=up{job="staging-api"}' | jq '.data.result'

# Test alert firing
curl -X POST https://staging-alerts.example.com/test \
  -d '{"alert":"staging_high_error_rate"}'
```

### 7. User Acceptance Testing

#### Enable Test Users
```sql
-- Create test accounts
INSERT INTO users (email, password_hash, role, is_test)
VALUES
  ('uat1@example.com', '$2b$10$...', 'user', true),
  ('uat2@example.com', '$2b$10$...', 'admin', true);

-- Enable feature flags for testing
UPDATE feature_flags
SET enabled = true
WHERE name IN ('new_checkout_flow', 'beta_dashboard')
AND environment = 'staging';
```

#### UAT Checklist
- [ ] Login/logout functionality
- [ ] Core user workflows
- [ ] Admin functions
- [ ] Payment processing
- [ ] Email notifications
- [ ] Report generation
- [ ] Mobile responsiveness
- [ ] Performance acceptable

### 8. Rollback Testing

#### Simulate Rollback
```bash
# Save current version
CURRENT_VERSION=$(kubectl get deployment api -n staging -o jsonpath='{.spec.template.spec.containers[0].image}')

# Deploy previous version
kubectl set image deployment/api api=registry.example.com/myapp:staging-previous -n staging

# Verify rollback
kubectl rollout status deployment/api -n staging

# Restore current version
kubectl set image deployment/api api=$CURRENT_VERSION -n staging
```

## Success Criteria

Staging deployment is successful when:
- [ ] All containers running and healthy
- [ ] Health checks passing
- [ ] Smoke tests passing
- [ ] Performance within thresholds
- [ ] No critical alerts firing
- [ ] Monitoring data flowing
- [ ] Logs accessible
- [ ] UAT sign-off received

## Troubleshooting

### Common Issues

#### Pods Not Starting
```bash
# Check pod events
kubectl describe pod <pod-name> -n staging

# Check logs
kubectl logs <pod-name> -n staging --previous

# Check resource constraints
kubectl top nodes
kubectl top pods -n staging
```

#### Database Connection Failures
```bash
# Test connectivity
kubectl run -it --rm debug --image=postgres:14 --restart=Never -n staging -- \
  psql -h staging-db.example.com -U app_user -d staging

# Check secrets
kubectl get secret db-credentials -n staging -o yaml
```

#### Performance Issues
```bash
# Check resource usage
kubectl top pods -n staging

# Review slow queries
psql -h staging-db.example.com -U admin -d staging \
  -c "SELECT query, mean_exec_time, calls
      FROM pg_stat_statements
      ORDER BY mean_exec_time DESC
      LIMIT 10;"

# Check cache hit rates
redis-cli -h staging-cache.example.com INFO stats
```

## Post-Deployment

### Documentation Updates
- Update deployment log
- Record version deployed
- Note any issues encountered
- Update runbook if needed
- Share learnings with team

### Preparation for Production
- Review staging metrics
- Confirm all tests passed
- Get stakeholder approval
- Schedule production deployment
- Prepare rollback plan

Remember: Staging is your last line of defense before production. Take it seriously, test thoroughly, and never skip steps even when under pressure.