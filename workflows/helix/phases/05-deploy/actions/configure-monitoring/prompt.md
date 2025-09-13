# Configure Monitoring Prompt

Set up comprehensive monitoring and observability infrastructure to track application health, performance, and business metrics in real-time across all environments.

## Purpose

Monitoring configuration ensures:
- Issues are detected before users complain
- Performance degradation is caught early
- Business impact is immediately visible
- Root cause analysis is data-driven
- Capacity planning is informed

## Monitoring Stack Setup

### 1. Metrics Collection

#### Prometheus Setup
```yaml
# prometheus-config.yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__

  - job_name: 'node-exporter'
    kubernetes_sd_configs:
      - role: node
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - '/etc/prometheus/rules/*.yml'
```

#### Application Metrics
```javascript
// Instrument application code
const prometheus = require('prom-client');

// Create metrics
const httpRequestDuration = new prometheus.Histogram({
  name: 'http_request_duration_seconds',
  help: 'Duration of HTTP requests in seconds',
  labelNames: ['method', 'route', 'status_code'],
  buckets: [0.1, 0.5, 1, 2, 5]
});

const httpRequestTotal = new prometheus.Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'route', 'status_code']
});

const businessMetrics = {
  ordersCreated: new prometheus.Counter({
    name: 'orders_created_total',
    help: 'Total number of orders created',
    labelNames: ['product_category', 'payment_method']
  }),

  orderValue: new prometheus.Histogram({
    name: 'order_value_dollars',
    help: 'Order value in dollars',
    buckets: [10, 50, 100, 500, 1000, 5000]
  })
};

// Middleware to track metrics
app.use((req, res, next) => {
  const start = Date.now();

  res.on('finish', () => {
    const duration = (Date.now() - start) / 1000;
    httpRequestDuration
      .labels(req.method, req.route?.path || 'unknown', res.statusCode)
      .observe(duration);

    httpRequestTotal
      .labels(req.method, req.route?.path || 'unknown', res.statusCode)
      .inc();
  });

  next();
});
```

### 2. Logging Infrastructure

#### Fluentd Configuration
```yaml
# fluentd-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      <parse>
        @type json
        time_format %Y-%m-%dT%H:%M:%S.%NZ
      </parse>
    </source>

    <filter kubernetes.**>
      @type kubernetes_metadata
    </filter>

    <filter kubernetes.**>
      @type record_transformer
      <record>
        hostname ${hostname}
        environment ${record["kubernetes"]["namespace_name"]}
        service ${record["kubernetes"]["labels"]["app"]}
      </record>
    </filter>

    <match kubernetes.**>
      @type elasticsearch
      host elasticsearch
      port 9200
      index_name fluentd
      type_name _doc
      logstash_format true
      logstash_prefix kubernetes
      <buffer>
        @type file
        path /var/log/fluentd-buffers/kubernetes.system.buffer
        flush_mode interval
        flush_interval 5s
      </buffer>
    </match>
```

#### Structured Logging
```javascript
// Use structured logging
const winston = require('winston');

const logger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: {
    service: 'api',
    version: process.env.VERSION,
    environment: process.env.NODE_ENV
  },
  transports: [
    new winston.transports.Console({
      format: winston.format.simple()
    })
  ]
});

// Log with context
logger.info('Order processed', {
  orderId: order.id,
  userId: user.id,
  amount: order.total,
  processingTime: endTime - startTime,
  paymentMethod: order.paymentMethod
});
```

### 3. Distributed Tracing

#### Jaeger Setup
```yaml
# jaeger-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jaeger
spec:
  selector:
    matchLabels:
      app: jaeger
  template:
    metadata:
      labels:
        app: jaeger
    spec:
      containers:
      - name: jaeger
        image: jaegertracing/all-in-one:latest
        ports:
        - containerPort: 5775
          protocol: UDP
        - containerPort: 6831
          protocol: UDP
        - containerPort: 6832
          protocol: UDP
        - containerPort: 5778
        - containerPort: 16686
        - containerPort: 14268
        env:
        - name: COLLECTOR_ZIPKIN_HTTP_PORT
          value: "9411"
        - name: SPAN_STORAGE_TYPE
          value: elasticsearch
        - name: ES_SERVER_URLS
          value: http://elasticsearch:9200
```

#### Application Tracing
```javascript
// Initialize tracing
const opentelemetry = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');

const provider = new NodeTracerProvider();

const jaegerExporter = new JaegerExporter({
  serviceName: 'api-service',
  agentHost: 'jaeger-agent',
  agentPort: 6831,
});

provider.addSpanProcessor(
  new BatchSpanProcessor(jaegerExporter)
);

provider.register();

// Trace operations
const tracer = opentelemetry.trace.getTracer('api-service');

async function processOrder(order) {
  const span = tracer.startSpan('process-order', {
    attributes: {
      'order.id': order.id,
      'order.value': order.total
    }
  });

  try {
    // Process order
    await validateOrder(order);
    await chargePayment(order);
    await updateInventory(order);
    await sendConfirmation(order);

    span.setStatus({ code: SpanStatusCode.OK });
  } catch (error) {
    span.setStatus({
      code: SpanStatusCode.ERROR,
      message: error.message
    });
    throw error;
  } finally {
    span.end();
  }
}
```

### 4. Dashboard Creation

#### Grafana Dashboards
```json
{
  "dashboard": {
    "title": "Application Overview",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{route}}"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status_code=~\"5..\"}[5m]) / rate(http_requests_total[5m])",
            "legendFormat": "Error Rate %"
          }
        ],
        "type": "graph",
        "alert": {
          "conditions": [
            {
              "evaluator": {
                "params": [0.01],
                "type": "gt"
              }
            }
          ]
        }
      },
      {
        "title": "Response Time (p95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Active Users",
        "targets": [
          {
            "expr": "active_users",
            "legendFormat": "Active Users"
          }
        ],
        "type": "stat"
      }
    ]
  }
}
```

### 5. Alert Rules

#### Prometheus Alert Rules
```yaml
# alerts.yml
groups:
  - name: application
    rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status_code=~"5.."}[5m]) > 0.05
      for: 5m
      labels:
        severity: critical
        team: backend
      annotations:
        summary: "High error rate detected"
        description: "Error rate is {{ $value | humanizePercentage }} for {{ $labels.instance }}"

    - alert: HighLatency
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
      for: 10m
      labels:
        severity: warning
        team: backend
      annotations:
        summary: "High latency detected"
        description: "95th percentile latency is {{ $value }}s"

    - alert: PodCrashLooping
      expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
      for: 5m
      labels:
        severity: critical
        team: platform
      annotations:
        summary: "Pod {{ $labels.namespace }}/{{ $labels.pod }} is crash looping"

  - name: business
    rules:
    - alert: LowOrderRate
      expr: rate(orders_created_total[1h]) < 10
      for: 30m
      labels:
        severity: warning
        team: business
      annotations:
        summary: "Order rate below threshold"
        description: "Only {{ $value }} orders per hour"

    - alert: PaymentFailureRate
      expr: rate(payment_failures_total[5m]) / rate(payment_attempts_total[5m]) > 0.1
      for: 5m
      labels:
        severity: critical
        team: payments
      annotations:
        summary: "High payment failure rate"
        description: "{{ $value | humanizePercentage }} of payments failing"
```

### 6. Log Aggregation Queries

#### Kibana Saved Searches
```json
{
  "saved_searches": [
    {
      "name": "Application Errors",
      "query": "level:ERROR AND service:api",
      "columns": ["timestamp", "message", "error.type", "trace_id"]
    },
    {
      "name": "Slow Queries",
      "query": "db.duration:>1000",
      "columns": ["timestamp", "db.query", "db.duration", "user_id"]
    },
    {
      "name": "Failed Payments",
      "query": "event:payment_failed",
      "columns": ["timestamp", "order_id", "amount", "error.message"]
    },
    {
      "name": "Security Events",
      "query": "event_type:security_*",
      "columns": ["timestamp", "event_type", "user_id", "ip_address"]
    }
  ]
}
```

## Monitoring Verification

### Health Check Endpoints
```javascript
// Comprehensive health check
app.get('/health', async (req, res) => {
  const health = {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    checks: {}
  };

  // Check database
  try {
    await db.query('SELECT 1');
    health.checks.database = 'healthy';
  } catch (error) {
    health.checks.database = 'unhealthy';
    health.status = 'degraded';
  }

  // Check cache
  try {
    await redis.ping();
    health.checks.cache = 'healthy';
  } catch (error) {
    health.checks.cache = 'unhealthy';
    health.status = 'degraded';
  }

  // Check external services
  try {
    await axios.get('https://payment-api.example.com/health');
    health.checks.payment_gateway = 'healthy';
  } catch (error) {
    health.checks.payment_gateway = 'unhealthy';
  }

  const statusCode = health.status === 'healthy' ? 200 : 503;
  res.status(statusCode).json(health);
});
```

### Synthetic Monitoring
```javascript
// Synthetic transaction monitoring
const syntheticTests = [
  {
    name: 'user_login',
    interval: 60000, // Every minute
    test: async () => {
      const response = await axios.post('/api/auth/login', {
        email: 'synthetic@example.com',
        password: 'synthetic-password'
      });
      assert(response.data.token, 'Login should return token');
    }
  },
  {
    name: 'product_search',
    interval: 300000, // Every 5 minutes
    test: async () => {
      const response = await axios.get('/api/products/search?q=test');
      assert(Array.isArray(response.data), 'Search should return array');
    }
  }
];

syntheticTests.forEach(test => {
  setInterval(async () => {
    const start = Date.now();
    try {
      await test.test();
      syntheticMetrics.success.labels(test.name).inc();
      syntheticMetrics.duration.labels(test.name).observe((Date.now() - start) / 1000);
    } catch (error) {
      syntheticMetrics.failure.labels(test.name).inc();
      logger.error(`Synthetic test failed: ${test.name}`, { error: error.message });
    }
  }, test.interval);
});
```

## Testing Monitoring

### Alert Testing
```bash
# Test alert firing
curl -X POST http://prometheus:9090/api/v1/admin/tsdb/create_block \
  -d '{
    "series": [{
      "labels": {"__name__": "test_metric", "severity": "critical"},
      "samples": [[1234567890, 100]]
    }]
  }'

# Verify alert received
curl http://alertmanager:9093/api/v1/alerts | jq '.data[0]'
```

### Load Testing Monitoring
```javascript
// Generate load to test monitoring
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    normal_load: {
      executor: 'constant-rate',
      rate: 100,
      duration: '5m',
      preAllocatedVUs: 10
    },
    spike_load: {
      executor: 'ramping-rate',
      startRate: 100,
      stages: [
        { duration: '2m', target: 500 },
        { duration: '1m', target: 1000 },
        { duration: '2m', target: 100 }
      ]
    }
  }
};

export default function() {
  // Generate various response codes
  const endpoint = Math.random() > 0.95 ? '/api/error' : '/api/products';
  const response = http.get(`https://staging.example.com${endpoint}`);

  check(response, {
    'status is 200': (r) => r.status === 200
  });
}
```

## Monitoring Maintenance

### Regular Tasks
- Review and tune alert thresholds weekly
- Archive old logs monthly
- Update dashboards based on incidents
- Prune unused metrics
- Optimize query performance
- Update documentation

### Cost Optimization
- Set appropriate retention periods
- Use sampling for high-volume metrics
- Compress old data
- Remove unused dashboards
- Optimize cardinality

Remember: Monitoring is not a one-time setup—it's an ongoing investment in system reliability. Good monitoring pays for itself by preventing outages and enabling rapid problem resolution.