# Smoke Test Prompt

Execute rapid validation tests immediately after deployment to verify that critical system functionality is working correctly in the target environment.

## Operational Action

This is an operational action that executes validation tests post-deployment. It does not generate documentation files but rather runs test suites and validation scripts to verify system health.

## Purpose

Smoke tests provide:
- Quick confidence that deployment succeeded
- Early detection of critical failures
- Go/no-go decision for continuing deployment
- Baseline functionality verification
- Rapid feedback loop

## Smoke Test Principles

### Fast and Focused
- Run in under 5 minutes total
- Test only critical paths
- Skip edge cases and comprehensive validation
- Fail fast on first critical error
- Provide clear pass/fail status

### Environment-Aware
- Adapt to target environment (staging/production)
- Use appropriate test data
- Respect rate limits and quotas
- Avoid destructive operations
- Clean up test artifacts

## Core Smoke Tests

### 1. Infrastructure Health

```bash
#!/bin/bash
# infrastructure-smoke.sh

echo "🔍 Running Infrastructure Smoke Tests..."

# Check all pods are running
check_pods() {
  echo -n "Checking pod status... "
  FAILED_PODS=$(kubectl get pods -n $NAMESPACE -o json | \
    jq -r '.items[] | select(.status.phase != "Running") | .metadata.name')

  if [ -z "$FAILED_PODS" ]; then
    echo "✅ All pods running"
    return 0
  else
    echo "❌ Failed pods: $FAILED_PODS"
    return 1
  fi
}

# Check service endpoints
check_services() {
  echo -n "Checking service endpoints... "
  SERVICES=("api" "web" "worker")

  for service in "${SERVICES[@]}"; do
    ENDPOINTS=$(kubectl get endpoints $service -n $NAMESPACE -o json | \
      jq -r '.subsets[].addresses | length')

    if [ "$ENDPOINTS" -eq 0 ]; then
      echo "❌ No endpoints for service: $service"
      return 1
    fi
  done

  echo "✅ All services have endpoints"
  return 0
}

# Check database connectivity
check_database() {
  echo -n "Checking database connectivity... "
  kubectl run db-check --image=postgres:14 --rm -i --restart=Never -- \
    psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1" &>/dev/null

  if [ $? -eq 0 ]; then
    echo "✅ Database accessible"
    return 0
  else
    echo "❌ Database connection failed"
    return 1
  fi
}

# Run all checks
check_pods && check_services && check_database
```

### 2. API Health Checks

```javascript
// api-smoke.js
const axios = require('axios');
const assert = require('assert');

const API_BASE_URL = process.env.API_URL || 'https://api.example.com';
const TIMEOUT = 5000; // 5 seconds per request

async function smokeTestAPI() {
  console.log('🔍 Running API Smoke Tests...');

  const tests = [
    {
      name: 'Health Check',
      request: {
        method: 'GET',
        url: '/health',
        timeout: TIMEOUT
      },
      validate: (response) => {
        assert.strictEqual(response.status, 200);
        assert.strictEqual(response.data.status, 'healthy');
      }
    },
    {
      name: 'Version Check',
      request: {
        method: 'GET',
        url: '/version',
        timeout: TIMEOUT
      },
      validate: (response) => {
        assert.strictEqual(response.status, 200);
        assert(response.data.version, 'Version should be present');
        console.log(`  Deployed version: ${response.data.version}`);
      }
    },
    {
      name: 'Database Connectivity',
      request: {
        method: 'GET',
        url: '/health/db',
        timeout: TIMEOUT
      },
      validate: (response) => {
        assert.strictEqual(response.status, 200);
        assert.strictEqual(response.data.database, 'connected');
      }
    },
    {
      name: 'Cache Connectivity',
      request: {
        method: 'GET',
        url: '/health/cache',
        timeout: TIMEOUT
      },
      validate: (response) => {
        assert.strictEqual(response.status, 200);
        assert.strictEqual(response.data.cache, 'connected');
      }
    }
  ];

  for (const test of tests) {
    try {
      const response = await axios({
        ...test.request,
        baseURL: API_BASE_URL
      });

      test.validate(response);
      console.log(`  ✅ ${test.name}`);
    } catch (error) {
      console.error(`  ❌ ${test.name}: ${error.message}`);
      throw new Error(`Smoke test failed: ${test.name}`);
    }
  }

  console.log('✅ All API smoke tests passed');
}
```

### 3. Authentication Flow

```javascript
// auth-smoke.js
async function smokeTestAuth() {
  console.log('🔍 Testing Authentication Flow...');

  try {
    // Test login
    const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
      email: 'smoke-test@example.com',
      password: process.env.SMOKE_TEST_PASSWORD
    });

    assert(loginResponse.data.token, 'Should receive auth token');
    const token = loginResponse.data.token;
    console.log('  ✅ Login successful');

    // Test authenticated request
    const profileResponse = await axios.get(`${API_BASE_URL}/profile`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    assert.strictEqual(profileResponse.status, 200);
    assert(profileResponse.data.email, 'Should return user profile');
    console.log('  ✅ Authenticated request successful');

    // Test logout
    const logoutResponse = await axios.post(`${API_BASE_URL}/auth/logout`, {}, {
      headers: { Authorization: `Bearer ${token}` }
    });

    assert.strictEqual(logoutResponse.status, 200);
    console.log('  ✅ Logout successful');

    // Verify token is invalidated
    try {
      await axios.get(`${API_BASE_URL}/profile`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      throw new Error('Token should be invalidated after logout');
    } catch (error) {
      if (error.response && error.response.status === 401) {
        console.log('  ✅ Token properly invalidated');
      } else {
        throw error;
      }
    }

  } catch (error) {
    console.error(`  ❌ Auth smoke test failed: ${error.message}`);
    throw error;
  }
}
```

### 4. Critical Business Flows

```javascript
// business-smoke.js
async function smokeTestBusinessFlows() {
  console.log('🔍 Testing Critical Business Flows...');

  const token = await getAuthToken();
  const headers = { Authorization: `Bearer ${token}` };

  // Test product listing
  const productsResponse = await axios.get(`${API_BASE_URL}/products`, { headers });
  assert(Array.isArray(productsResponse.data), 'Products should be an array');
  assert(productsResponse.data.length > 0, 'Should have products');
  console.log('  ✅ Product listing works');

  // Test search
  const searchResponse = await axios.get(`${API_BASE_URL}/products/search?q=test`, { headers });
  assert(Array.isArray(searchResponse.data), 'Search should return array');
  console.log('  ✅ Search functionality works');

  // Test cart operations (non-destructive)
  const cartResponse = await axios.get(`${API_BASE_URL}/cart`, { headers });
  assert.strictEqual(cartResponse.status, 200);
  console.log('  ✅ Cart accessible');

  // Test checkout flow (validation only)
  const checkoutValidation = await axios.post(`${API_BASE_URL}/checkout/validate`, {
    items: [{ productId: 'test-001', quantity: 1 }],
    dryRun: true
  }, { headers });

  assert.strictEqual(checkoutValidation.status, 200);
  console.log('  ✅ Checkout validation works');
}
```

### 5. UI Smoke Tests

```javascript
// ui-smoke.js
const puppeteer = require('puppeteer');

async function smokeTestUI() {
  console.log('🔍 Testing UI Smoke Tests...');

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  try {
    // Test homepage loads
    await page.goto(process.env.WEB_URL || 'https://example.com', {
      waitUntil: 'networkidle2',
      timeout: 10000
    });

    const title = await page.title();
    assert(title, 'Page should have a title');
    console.log('  ✅ Homepage loads');

    // Test critical elements exist
    const criticalElements = [
      { selector: 'header', name: 'Header' },
      { selector: 'nav', name: 'Navigation' },
      { selector: 'main', name: 'Main content' },
      { selector: 'footer', name: 'Footer' }
    ];

    for (const element of criticalElements) {
      const exists = await page.$(element.selector) !== null;
      assert(exists, `${element.name} should exist`);
      console.log(`  ✅ ${element.name} present`);
    }

    // Test page performance
    const metrics = await page.metrics();
    assert(metrics.TaskDuration < 5000, 'Page should load in under 5 seconds');
    console.log(`  ✅ Page load time: ${Math.round(metrics.TaskDuration)}ms`);

    // Test JavaScript errors
    const jsErrors = [];
    page.on('pageerror', error => jsErrors.push(error.message));
    await page.reload();

    assert(jsErrors.length === 0, `No JavaScript errors: ${jsErrors.join(', ')}`);
    console.log('  ✅ No JavaScript errors');

  } finally {
    await browser.close();
  }
}
```

### 6. External Service Integration

```javascript
// external-smoke.js
async function smokeTestExternalServices() {
  console.log('🔍 Testing External Service Integrations...');

  const services = [
    {
      name: 'Payment Gateway',
      test: async () => {
        const response = await axios.post(`${API_BASE_URL}/payments/validate`, {
          testMode: true,
          amount: 100,
          currency: 'USD'
        });
        assert.strictEqual(response.data.status, 'valid');
      }
    },
    {
      name: 'Email Service',
      test: async () => {
        const response = await axios.post(`${API_BASE_URL}/email/test`, {
          to: 'smoke-test@example.com',
          subject: 'Smoke Test',
          dryRun: true
        });
        assert.strictEqual(response.data.status, 'ready');
      }
    },
    {
      name: 'CDN',
      test: async () => {
        const response = await axios.head('https://cdn.example.com/static/app.js');
        assert.strictEqual(response.status, 200);
      }
    },
    {
      name: 'Analytics',
      test: async () => {
        const response = await axios.post(`${API_BASE_URL}/analytics/test`, {
          event: 'smoke_test',
          properties: { timestamp: Date.now() }
        });
        assert.strictEqual(response.status, 204);
      }
    }
  ];

  for (const service of services) {
    try {
      await service.test();
      console.log(`  ✅ ${service.name}`);
    } catch (error) {
      console.error(`  ❌ ${service.name}: ${error.message}`);
      // Don't fail deployment for non-critical external services
      if (service.critical) {
        throw error;
      }
    }
  }
}
```

## Smoke Test Orchestration

```javascript
// run-smoke-tests.js
const chalk = require('chalk');

async function runAllSmokeTests() {
  const startTime = Date.now();
  console.log(chalk.blue.bold('\n🚀 Starting Smoke Tests...\n'));

  const testSuites = [
    { name: 'Infrastructure', fn: smokeTestInfrastructure, critical: true },
    { name: 'API Health', fn: smokeTestAPI, critical: true },
    { name: 'Authentication', fn: smokeTestAuth, critical: true },
    { name: 'Business Flows', fn: smokeTestBusinessFlows, critical: true },
    { name: 'UI', fn: smokeTestUI, critical: false },
    { name: 'External Services', fn: smokeTestExternalServices, critical: false }
  ];

  const results = [];

  for (const suite of testSuites) {
    console.log(chalk.yellow(`\nRunning ${suite.name}...`));

    try {
      await suite.fn();
      results.push({ suite: suite.name, status: 'passed' });
      console.log(chalk.green(`✅ ${suite.name} passed`));
    } catch (error) {
      results.push({ suite: suite.name, status: 'failed', error: error.message });
      console.error(chalk.red(`❌ ${suite.name} failed: ${error.message}`));

      if (suite.critical) {
        console.error(chalk.red.bold('\n🛑 Critical smoke test failed. Stopping tests.'));
        process.exit(1);
      }
    }
  }

  const duration = Math.round((Date.now() - startTime) / 1000);
  const passed = results.filter(r => r.status === 'passed').length;
  const failed = results.filter(r => r.status === 'failed').length;

  console.log(chalk.blue.bold('\n📊 Smoke Test Summary:'));
  console.log(chalk.green(`  Passed: ${passed}`));
  console.log(chalk.red(`  Failed: ${failed}`));
  console.log(chalk.gray(`  Duration: ${duration}s`));

  if (failed === 0) {
    console.log(chalk.green.bold('\n✅ All smoke tests passed! Deployment validated.\n'));
    process.exit(0);
  } else {
    console.log(chalk.yellow.bold('\n⚠️ Some non-critical tests failed. Review before proceeding.\n'));
    process.exit(0);
  }
}

// Run if executed directly
if (require.main === module) {
  runAllSmokeTests().catch(error => {
    console.error(chalk.red.bold('Unexpected error:'), error);
    process.exit(1);
  });
}
```

## Environment-Specific Configuration

```javascript
// smoke-config.js
const config = {
  staging: {
    apiUrl: 'https://staging-api.example.com',
    webUrl: 'https://staging.example.com',
    testUser: 'smoke-staging@example.com',
    timeout: 10000,
    retries: 2
  },
  production: {
    apiUrl: 'https://api.example.com',
    webUrl: 'https://example.com',
    testUser: 'smoke-prod@example.com',
    timeout: 5000,
    retries: 1
  }
};

module.exports = config[process.env.ENVIRONMENT] || config.staging;
```

## CI/CD Integration

```yaml
# .github/workflows/smoke-tests.yml
name: Smoke Tests

on:
  deployment_status:

jobs:
  smoke-tests:
    if: github.event.deployment_status.state == 'success'
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v2

      - name: Setup Node.js
        uses: actions/setup-node@v2
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci

      - name: Run smoke tests
        env:
          ENVIRONMENT: ${{ github.event.deployment.environment }}
          API_URL: ${{ github.event.deployment.payload.api_url }}
        run: npm run smoke-test

      - name: Report results
        if: always()
        uses: actions/github-script@v6
        with:
          script: |
            const status = ${{ job.status == 'success' }} ? 'success' : 'failure';
            await github.repos.createDeploymentStatus({
              ...context.repo,
              deployment_id: context.payload.deployment.id,
              state: status,
              description: 'Smoke tests ' + status
            });
```

## Success Criteria

Smoke tests pass when:
- All critical tests pass (100% required)
- Non-critical test pass rate > 80%
- Total execution time < 5 minutes
- No timeouts or network errors
- Key metrics within expected ranges

Remember: Smoke tests are your first line of defense after deployment. Keep them fast, focused, and reliable. They should give you confidence to proceed or the wisdom to rollback.