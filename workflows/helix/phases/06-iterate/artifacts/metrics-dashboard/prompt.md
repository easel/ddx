# Metrics Dashboard Generation Prompt

Create a comprehensive metrics dashboard with AI-powered insights for production system monitoring and analysis.

## Storage Location

Store the dashboard at: `docs/helix/06-iterate/metrics-dashboard.md`

Update this dashboard regularly (daily/weekly) to track trends and identify patterns across iterations.

## Dashboard Purpose

The Metrics Dashboard is an **AI-enhanced monitoring tool** that:
- Provides real-time visibility into system health
- Detects anomalies before they become incidents
- Correlates metrics to identify root causes
- Predicts future trends and capacity needs
- Recommends optimizations based on data patterns

## Key Principles

### 1. Data-Driven Insights Over Raw Numbers
- ✅ "User engagement dropped 15% correlated with slower page loads"
- ❌ "DAU is 5,432"

### 2. Predictive Over Reactive
- Use historical patterns to forecast issues
- Alert on trending problems before thresholds
- Recommend preventive actions

### 3. Actionable Over Informational
- Every insight should suggest an action
- Prioritize by impact and effort
- Link metrics to business outcomes

### 4. Holistic Over Siloed
- Correlate across metric categories
- Understand cascading effects
- Consider external factors

## Section-by-Section Guidance

### Executive Summary
Generate an AI summary that:
- Highlights the most significant changes
- Identifies critical issues requiring attention
- Celebrates improvements and successes
- Provides a health score with justification

### Business Metrics
Focus on user value:
- Track engagement trends
- Identify user behavior patterns
- Correlate with feature releases
- Predict churn risks

### Technical Metrics
Monitor system health:
- Detect performance degradation
- Identify capacity constraints
- Track error patterns
- Optimize resource usage

### Anomaly Detection
Use AI to identify:
- Unusual patterns in normal ranges
- Correlations between seemingly unrelated metrics
- Early warning signs of issues
- Statistical outliers requiring investigation

### Predictive Analytics
Forecast based on:
- Historical trend analysis
- Seasonal patterns
- Growth trajectories
- External factor correlations

### Recommendations
Prioritize actions by:
- Impact on user experience
- Business value
- Implementation effort
- Risk mitigation

## AI Analysis Techniques

### Pattern Recognition
```python
# Example patterns to detect:
- Periodic spikes (daily, weekly patterns)
- Gradual degradation over time
- Sudden step changes after deployments
- Correlation between user actions and errors
```

### Anomaly Detection Methods
- **Statistical**: Z-score, moving averages
- **Machine Learning**: Isolation forests, clustering
- **Time Series**: ARIMA, Prophet
- **Correlation**: Pearson, Spearman coefficients

### Predictive Models
- **Regression**: Linear, polynomial for trends
- **Classification**: Random forest for alert prediction
- **Time Series**: LSTM for complex patterns
- **Ensemble**: Combine multiple models for accuracy

## Data Sources

### Required Integrations
- Application Performance Monitoring (APM)
- Infrastructure monitoring
- User analytics
- Business intelligence
- Log aggregation
- Error tracking

### Collection Frequency
- Real-time: Critical metrics (errors, uptime)
- Minutely: Performance metrics
- Hourly: User engagement
- Daily: Business metrics
- Weekly: Trends and patterns

## Quality Checklist

Before publishing the dashboard:
- [ ] All metrics have current values and trends
- [ ] Anomalies are investigated and explained
- [ ] Predictions include confidence levels
- [ ] Recommendations are specific and actionable
- [ ] Correlations are statistically significant
- [ ] Historical comparisons provide context

## Common Pitfalls to Avoid

### ❌ Vanity Metrics
**Bad**: Tracking metrics that look good but don't drive decisions
**Good**: Focus on actionable metrics tied to outcomes

### ❌ Alert Fatigue
**Bad**: Too many alerts for minor fluctuations
**Good**: Smart thresholds with anomaly detection

### ❌ Missing Context
**Bad**: Numbers without explanation
**Good**: Insights explaining why metrics changed

### ❌ Correlation vs Causation
**Bad**: Assuming correlation implies causation
**Good**: Investigate root causes before concluding

## Advanced Features

### Multi-Dimensional Analysis
- Segment metrics by user cohort
- Compare across different time periods
- Analyze by feature flag exposure
- Break down by geographic region

### Intelligent Alerting
- Dynamic thresholds based on patterns
- Severity classification using ML
- Alert suppression during known events
- Predictive alerting before issues occur

### Root Cause Analysis
- Automated correlation analysis
- Dependency mapping
- Change impact assessment
- Historical pattern matching

## Visualization Guidelines

### Effective Presentations
- Use sparklines for trends in tables
- Heat maps for correlation matrices
- Time series for historical data
- Scatter plots for correlations

### Color Coding
- 🟢 Green: Within normal range
- 🟡 Yellow: Approaching threshold
- 🔴 Red: Requires immediate attention
- 🔵 Blue: Informational/neutral

## Integration with Iteration

The dashboard should:
1. Track metrics defined in Frame phase success criteria
2. Monitor SLOs defined in Design phase
3. Validate improvements from Build phase
4. Confirm Deploy phase stability
5. Feed insights into next Frame phase

## Automation Opportunities

### Scheduled Reports
- Daily summary for team standup
- Weekly executive briefing
- Monthly trend analysis
- Quarterly business review

### Triggered Analysis
- Post-deployment validation
- Incident detection and diagnosis
- Capacity planning alerts
- Cost optimization recommendations

## Remember

The Metrics Dashboard is about **turning data into decisions**. It should answer:
1. Is the system healthy right now?
2. Are we meeting our business goals?
3. What problems are emerging?
4. What should we do next?
5. Are our improvements working?

A good dashboard prevents surprises, enables proactive decisions, and demonstrates continuous improvement.