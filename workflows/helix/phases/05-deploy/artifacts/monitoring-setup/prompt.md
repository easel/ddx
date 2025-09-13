# Monitoring Setup Generation Prompt

Create a comprehensive monitoring and observability setup that provides real-time insights into system health, performance, and user experience throughout the deployment and operation lifecycle.

## Storage Location

Store the monitoring setup at: `docs/monitoring/setup.md`

## Purpose

The monitoring setup ensures:
- System health is continuously tracked
- Issues are detected before users notice
- Performance degradation is caught early
- Business metrics are visible
- Incident response is data-driven

## Key Requirements

### 1. Metrics Collection

#### Application Metrics
Define what to measure:
- **Request metrics**: Rate, errors, duration (RED)
- **Resource metrics**: CPU, memory, disk, network
- **Business metrics**: Conversions, revenue, user actions
- **Custom metrics**: Application-specific measurements

#### Infrastructure Metrics
Monitor underlying systems:
- Server health and utilization
- Database performance
- Cache hit rates
- Queue depths
- Network latency
- Storage usage

### 2. Logging Strategy

#### Log Levels and Structure
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "ERROR",
  "service": "api",
  "trace_id": "abc123",
  "user_id": "user456",
  "message": "Payment processing failed",
  "error": {
    "type": "PaymentGatewayError",
    "code": "INSUFFICIENT_FUNDS",
    "details": {}
  }
}
```

#### Log Aggregation
- Centralized logging platform
- Log retention policies
- Search and filter capabilities
- Correlation across services
- Sensitive data masking
- Compliance requirements

### 3. Dashboard Design

#### Operations Dashboard
Real-time system health:
- Service availability
- Current error rates
- Active users
- Response times
- Resource utilization
- Deployment status

#### Business Dashboard
Key business metrics:
- User acquisition
- Revenue metrics
- Feature adoption
- User engagement
- Conversion funnels
- Customer satisfaction

#### Performance Dashboard
System performance details:
- API endpoint latencies
- Database query times
- Cache performance
- Third-party service calls
- Background job processing
- Resource consumption trends

### 4. Alert Configuration

#### Alert Rules
Define conditions and thresholds:

```yaml
alerts:
  - name: high_error_rate
    condition: error_rate > 1%
    duration: 5m
    severity: critical

  - name: slow_response_time
    condition: p95_latency > 1000ms
    duration: 10m
    severity: warning

  - name: low_disk_space
    condition: disk_usage > 85%
    duration: 5m
    severity: warning
```

#### Alert Routing
- Severity-based escalation
- On-call rotation integration
- Notification channels (email, SMS, Slack)
- Alert grouping and deduplication
- Acknowledgment tracking
- Escalation procedures

### 5. Distributed Tracing

#### Trace Implementation
Track requests across services:
- Request flow visualization
- Latency breakdown
- Error propagation
- Dependency mapping
- Performance bottlenecks
- Service dependencies

#### Trace Sampling
Balance visibility and overhead:
- 100% sampling for errors
- 10% sampling for normal traffic
- Dynamic sampling based on endpoints
- Trace retention policies
- Cost optimization

### 6. SLI/SLO Definition

#### Service Level Indicators (SLIs)
Measurable aspects of service:
- Availability: Successful requests / Total requests
- Latency: Requests under 200ms / Total requests
- Quality: Successful transactions / Total transactions
- Freshness: Data updated within SLA / Total updates

#### Service Level Objectives (SLOs)
Targets for SLIs:
- Availability: 99.9% (43.8 minutes downtime/month)
- Latency: 95% of requests < 200ms
- Quality: 99.5% transaction success rate
- Freshness: 99% of data < 5 minutes old

#### Error Budgets
Managing risk and innovation:
- Monthly error budget calculation
- Budget consumption tracking
- Alert when budget depleted
- Innovation freeze triggers
- Budget reset policies

## Monitoring Stack Components

### 1. Metrics Platform
- Time-series database (Prometheus, InfluxDB)
- Metrics collection agents
- Service discovery
- Data retention policies
- Query language and APIs

### 2. Logging Platform
- Log aggregation (ELK, Splunk)
- Log shippers and agents
- Index management
- Search interfaces
- Retention and archival

### 3. Tracing Platform
- Distributed tracing (Jaeger, Zipkin)
- Trace collectors
- Storage backend
- UI for trace analysis
- Integration with metrics

### 4. Visualization Platform
- Dashboard tools (Grafana, Datadog)
- Custom visualizations
- Mobile access
- Sharing and embedding
- Automated reports

## Implementation Guide

### Quick Start Monitoring
Essential monitoring to implement first:
1. Health check endpoints
2. Basic resource metrics
3. Error tracking
4. Key business metrics
5. Critical alerts

### Progressive Enhancement
Add sophisticated monitoring over time:
1. Distributed tracing
2. Custom dashboards
3. Predictive alerts
4. Anomaly detection
5. Capacity planning

## Alert Fatigue Prevention

### Smart Alerting
Reduce noise, increase signal:
- Alert only on user impact
- Group related alerts
- Use alert dependencies
- Implement quiet hours
- Auto-resolve when fixed

### Alert Quality
Each alert should be:
- **Actionable**: Clear steps to resolve
- **Urgent**: Requires immediate attention
- **Unique**: Not duplicate of other alerts
- **Documented**: Runbook linked
- **Tested**: Verified to work correctly

## Monitoring Checklist

Before deployment, ensure:
- [ ] All services have health checks
- [ ] Metrics are being collected
- [ ] Logs are aggregated centrally
- [ ] Dashboards are configured
- [ ] Alerts are set up and tested
- [ ] On-call rotation is configured
- [ ] Runbooks are linked to alerts
- [ ] SLOs are defined and tracked
- [ ] Tracing is implemented
- [ ] Documentation is complete

## Cost Optimization

### Data Management
Control monitoring costs:
- Appropriate retention periods
- Sampling strategies
- Compression and archival
- Metric cardinality limits
- Log level management

### Tool Selection
Balance features and cost:
- Open source vs. commercial
- Self-hosted vs. managed
- Pay-per-use vs. fixed cost
- Feature requirements
- Team expertise

## Incident Response Integration

### During Incidents
Monitoring provides:
- Initial detection and alerting
- Impact assessment
- Root cause investigation
- Progress tracking
- Resolution verification

### Post-Incident
Monitoring enables:
- Timeline reconstruction
- Impact analysis
- Prevention measures
- SLO impact calculation
- Report generation

## Quality Checklist

Before monitoring setup is complete:
- [ ] All critical paths monitored
- [ ] Dashboards load quickly
- [ ] Alerts are actionable
- [ ] No alert storms possible
- [ ] Data retention appropriate
- [ ] Costs are sustainable
- [ ] Team is trained
- [ ] Documentation complete
- [ ] Tested during deployment
- [ ] Integrated with incident response

Remember: Good monitoring is like insurance—you hope you don't need it, but you're glad it's there when you do. Invest in observability before you need it.