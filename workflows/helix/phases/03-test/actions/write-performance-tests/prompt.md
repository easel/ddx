# Write Performance Tests Prompt

Create performance tests that validate system responsiveness, scalability, and resource efficiency under various load conditions. These tests ensure the application meets non-functional requirements for speed and capacity.

## Test Output Location

Generate performance tests in: `tests/performance/`

Organize tests by performance type:
- `tests/performance/load/` - Load testing scripts
- `tests/performance/stress/` - Stress testing scenarios
- `tests/performance/spike/` - Spike testing configurations
- `tests/performance/endurance/` - Endurance testing suites

## Purpose

Performance tests verify:
- Response time under normal load
- System behavior under peak load
- Breaking points and failure modes
- Resource consumption patterns
- Scalability characteristics

## Types of Performance Tests

### Load Testing
Verify system performance under expected normal conditions:
- Typical number of concurrent users
- Average transaction rates
- Normal data volumes
- Expected usage patterns

### Stress Testing
Determine system limits and breaking points:
- Maximum concurrent users
- Peak transaction rates
- Resource exhaustion points
- Recovery behavior after stress

### Spike Testing
Test sudden load increases:
- Flash sale scenarios
- Viral content situations
- Marketing campaign impacts
- Unexpected traffic surges

### Endurance Testing
Verify stability over extended periods:
- Memory leak detection
- Resource degradation
- Database connection pooling
- Cache effectiveness

## Performance Metrics

### Response Time Metrics
```javascript
// Define performance budgets
const performanceBudgets = {
  api: {
    p50: 100,   // 50th percentile: 100ms
    p95: 500,   // 95th percentile: 500ms
    p99: 1000,  // 99th percentile: 1s
    max: 3000   // Maximum: 3s
  },
  web: {
    firstContentfulPaint: 1500,
    timeToInteractive: 3000,
    largestContentfulPaint: 2500,
    totalBlockingTime: 300
  }
};
```

### Throughput Metrics
```javascript
const throughputTargets = {
  apiRequests: {
    readsPerSecond: 10000,
    writesPerSecond: 1000,
    searchQueriesPerSecond: 500
  },
  database: {
    queriesPerSecond: 5000,
    transactionsPerSecond: 500
  },
  messageQueue: {
    messagesPerSecond: 10000
  }
};
```

## Load Testing Implementation

### Using k6 for Load Testing
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const apiDuration = new Trend('api_duration');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp up to 200 users
    { duration: '5m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    errors: ['rate<0.01'], // Error rate < 1%
  },
};

export default function () {
  // Test user login
  const loginRes = http.post('https://api.example.com/auth/login', {
    email: 'user@example.com',
    password: 'password123'
  });

  check(loginRes, {
    'login successful': (r) => r.status === 200,
    'token received': (r) => r.json('token') !== '',
  });

  errorRate.add(loginRes.status !== 200);
  apiDuration.add(loginRes.timings.duration);

  const token = loginRes.json('token');

  // Test API endpoints
  const headers = { Authorization: `Bearer ${token}` };

  // Get user profile
  const profileRes = http.get('https://api.example.com/profile', { headers });
  check(profileRes, {
    'profile retrieved': (r) => r.status === 200,
  });

  // Search products
  const searchRes = http.get('https://api.example.com/products?q=laptop', { headers });
  check(searchRes, {
    'search successful': (r) => r.status === 200,
    'results returned': (r) => r.json('products').length > 0,
  });

  sleep(1); // Think time between requests
}
```

### Database Performance Testing
```javascript
// Test database query performance
const { performance } = require('perf_hooks');

async function testDatabasePerformance() {
  const iterations = 1000;
  const concurrency = 10;

  // Test read performance
  const readTimes = [];
  for (let i = 0; i < iterations; i++) {
    const start = performance.now();
    await db.query('SELECT * FROM users WHERE id = $1', [Math.floor(Math.random() * 10000)]);
    const end = performance.now();
    readTimes.push(end - start);
  }

  // Test write performance
  const writeTimes = [];
  for (let i = 0; i < iterations; i++) {
    const start = performance.now();
    await db.query('INSERT INTO logs (message, timestamp) VALUES ($1, $2)', ['test', new Date()]);
    const end = performance.now();
    writeTimes.push(end - start);
  }

  // Calculate statistics
  const stats = {
    read: calculateStats(readTimes),
    write: calculateStats(writeTimes)
  };

  // Assert performance requirements
  expect(stats.read.p95).toBeLessThan(50); // 95% of reads < 50ms
  expect(stats.write.p95).toBeLessThan(100); // 95% of writes < 100ms
}

function calculateStats(times) {
  times.sort((a, b) => a - b);
  return {
    min: times[0],
    max: times[times.length - 1],
    mean: times.reduce((a, b) => a + b) / times.length,
    p50: times[Math.floor(times.length * 0.5)],
    p95: times[Math.floor(times.length * 0.95)],
    p99: times[Math.floor(times.length * 0.99)]
  };
}
```

## Stress Testing

### Finding Breaking Points
```javascript
export const options = {
  stages: [
    { duration: '5m', target: 500 },   // Ramp to 500 users
    { duration: '5m', target: 1000 },  // Ramp to 1000 users
    { duration: '5m', target: 2000 },  // Ramp to 2000 users
    { duration: '5m', target: 3000 },  // Ramp to 3000 users
    { duration: '5m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_failed: [{
      threshold: 'rate<0.05',
      abortOnFail: true,
      delayAbortEval: '30s'
    }],
  },
};

export default function () {
  const res = http.get('https://api.example.com/products');

  // Track when system starts failing
  if (res.status !== 200) {
    console.error(`System failing at ${__VU} virtual users`);
  }

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time OK': (r) => r.timings.duration < 2000,
  });
}
```

## Memory and Resource Testing

### Memory Leak Detection
```javascript
// Node.js memory monitoring
const v8 = require('v8');
const { performance } = require('perf_hooks');

function monitorMemory() {
  const measurements = [];

  setInterval(() => {
    const heapStats = v8.getHeapStatistics();
    measurements.push({
      timestamp: Date.now(),
      heapUsed: heapStats.used_heap_size,
      heapTotal: heapStats.total_heap_size,
      external: heapStats.external_memory,
      rss: process.memoryUsage().rss
    });

    // Check for memory leak indicators
    if (measurements.length > 100) {
      const recent = measurements.slice(-100);
      const trend = calculateTrend(recent.map(m => m.heapUsed));

      if (trend > 0.5) { // Growing by 50% over 100 samples
        console.warn('Potential memory leak detected');
      }
    }
  }, 1000);
}

function calculateTrend(values) {
  const n = values.length;
  const sumX = n * (n - 1) / 2;
  const sumY = values.reduce((a, b) => a + b);
  const sumXY = values.reduce((sum, y, x) => sum + x * y, 0);
  const sumX2 = n * (n - 1) * (2 * n - 1) / 6;

  const slope = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX);
  return slope / (sumY / n); // Normalized slope
}
```

## Frontend Performance Testing

### Web Vitals Testing
```javascript
// Using Puppeteer for frontend performance
const puppeteer = require('puppeteer');

async function testWebPerformance() {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  // Enable performance monitoring
  await page.evaluateOnNewDocument(() => {
    window.performanceMetrics = [];

    // Observe all performance entries
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        window.performanceMetrics.push(entry);
      }
    });

    observer.observe({ entryTypes: ['paint', 'largest-contentful-paint', 'layout-shift', 'first-input'] });
  });

  // Navigate and measure
  await page.goto('https://example.com', { waitUntil: 'networkidle0' });

  // Get Core Web Vitals
  const metrics = await page.evaluate(() => {
    return {
      FCP: performance.getEntriesByName('first-contentful-paint')[0]?.startTime,
      LCP: window.performanceMetrics.find(m => m.entryType === 'largest-contentful-paint')?.startTime,
      CLS: window.performanceMetrics
        .filter(m => m.entryType === 'layout-shift')
        .reduce((sum, entry) => sum + entry.value, 0),
      TTFB: performance.timing.responseStart - performance.timing.requestStart
    };
  });

  // Assert performance budgets
  expect(metrics.FCP).toBeLessThan(1800);
  expect(metrics.LCP).toBeLessThan(2500);
  expect(metrics.CLS).toBeLessThan(0.1);
  expect(metrics.TTFB).toBeLessThan(600);

  await browser.close();
}
```

## API Performance Testing

### Endpoint-Specific Tests
```javascript
// Test specific API endpoints
export default function () {
  const scenarios = [
    {
      name: 'GET /users',
      request: () => http.get('https://api.example.com/users'),
      threshold: { p95: 200, p99: 500 }
    },
    {
      name: 'POST /users',
      request: () => http.post('https://api.example.com/users', { name: 'Test User' }),
      threshold: { p95: 300, p99: 700 }
    },
    {
      name: 'GET /products/search',
      request: () => http.get('https://api.example.com/products/search?q=laptop&limit=50'),
      threshold: { p95: 500, p99: 1000 }
    }
  ];

  scenarios.forEach(scenario => {
    const res = scenario.request();

    check(res, {
      [`${scenario.name} - status 200`]: (r) => r.status === 200,
      [`${scenario.name} - p95 < ${scenario.threshold.p95}ms`]: (r) => r.timings.duration < scenario.threshold.p95,
    });
  });
}
```

## Reporting and Analysis

### Performance Report Generation
```javascript
function generatePerformanceReport(results) {
  return {
    summary: {
      totalRequests: results.metrics.http_reqs.count,
      failedRequests: results.metrics.http_req_failed.count,
      errorRate: results.metrics.http_req_failed.rate,
      avgResponseTime: results.metrics.http_req_duration.avg,
      p95ResponseTime: results.metrics.http_req_duration.p95,
      p99ResponseTime: results.metrics.http_req_duration.p99,
      throughput: results.metrics.http_reqs.rate
    },
    thresholds: {
      passed: results.thresholds.passed,
      failed: results.thresholds.failed
    },
    recommendations: analyzeResults(results)
  };
}

function analyzeResults(results) {
  const recommendations = [];

  if (results.metrics.http_req_duration.p95 > 1000) {
    recommendations.push('Consider caching frequently accessed data');
    recommendations.push('Optimize database queries');
  }

  if (results.metrics.http_req_failed.rate > 0.01) {
    recommendations.push('Investigate error patterns');
    recommendations.push('Implement circuit breakers');
  }

  return recommendations;
}
```

## Configuration

### Performance Test Environment
```yaml
# docker-compose.perf.yml
version: '3.8'
services:
  app:
    image: myapp:latest
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
    environment:
      - NODE_ENV=performance
      - DB_POOL_SIZE=50
      - CACHE_SIZE=1000

  database:
    image: postgres:14
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
    environment:
      - POSTGRES_MAX_CONNECTIONS=200
      - POSTGRES_SHARED_BUFFERS=1GB

  cache:
    image: redis:7
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
```

## Quality Checklist

Before performance tests are complete:
- [ ] Load testing covers expected traffic patterns
- [ ] Stress testing identifies breaking points
- [ ] Response time targets are validated
- [ ] Throughput requirements are met
- [ ] Resource usage is within limits
- [ ] Memory leaks are tested
- [ ] Database performance is validated
- [ ] Frontend performance meets Core Web Vitals
- [ ] API endpoints meet SLA requirements
- [ ] Performance regression tests are automated

## Best Practices

### DO
- ✅ Test in production-like environments
- ✅ Use realistic data volumes
- ✅ Monitor all system resources
- ✅ Test gradually increasing loads
- ✅ Include think time in scripts
- ✅ Test both read and write operations
- ✅ Measure percentiles, not just averages
- ✅ Automate performance regression testing

### DON'T
- ❌ Test against production directly
- ❌ Use unrealistic test data
- ❌ Ignore warmup periods
- ❌ Focus only on happy paths
- ❌ Test without monitoring
- ❌ Ignore database performance
- ❌ Skip endurance testing
- ❌ Test only at maximum load

Remember: Performance is a feature. Test it early, test it often, and establish clear performance budgets that align with user expectations.