# Feedback Analysis Generation Prompt

Synthesize user feedback using AI-powered sentiment analysis, pattern recognition, and priority scoring to drive product improvements.

## Storage Location

Store the analysis at: `docs/06-iterate/feedback-analysis.md`

Update this analysis weekly or after significant feedback volumes to maintain pulse on user satisfaction.

## Analysis Purpose

The Feedback Analysis is an **AI-enhanced synthesis tool** that:
- Aggregates feedback from multiple channels
- Identifies patterns and trends in user sentiment
- Prioritizes improvements based on impact
- Predicts churn risks and satisfaction drivers
- Connects user voice to product decisions

## Key Principles

### 1. Patterns Over Individual Comments
- Look for recurring themes across feedback
- Weight by frequency and user segment value
- Identify root causes behind surface complaints

### 2. Sentiment as Leading Indicator
- Track sentiment trends before they impact metrics
- Correlate sentiment with user behavior
- Predict future NPS/CSAT movements

### 3. Actionable Over Interesting
- Every insight must drive a decision
- Prioritize by user impact and business value
- Connect feedback to specific features/changes

### 4. Segmented Over Averaged
- Different user segments have different needs
- Personalize improvements by segment
- Avoid one-size-fits-all solutions

## Data Collection Sources

### Primary Channels
- **In-app feedback**: Surveys, feedback widgets
- **Support tickets**: Email, chat, phone transcripts
- **App store reviews**: iOS, Android, desktop
- **Social media**: Twitter/X, Reddit, LinkedIn
- **User interviews**: Transcripts and notes
- **Community forums**: Discord, Slack, discussion boards
- **NPS/CSAT surveys**: Scores and comments

### Collection Best Practices
- Normalize data formats across channels
- Preserve context and metadata
- Track feedback source and user segment
- Maintain temporal ordering
- Handle multiple languages appropriately

## AI Analysis Techniques

### Sentiment Analysis
```python
# Multi-level sentiment scoring:
- Document level: Overall feedback sentiment
- Aspect level: Sentiment per feature/topic
- Temporal: Sentiment evolution over time
- Comparative: Sentiment vs competitors
```

### Topic Modeling
- **LDA**: Discover hidden topic structures
- **BERT**: Contextual topic understanding
- **Clustering**: Group similar feedback
- **Keyword Extraction**: Identify key terms

### Priority Scoring Algorithm
```python
priority_score = (
    user_segment_value * 0.3 +
    frequency_of_mention * 0.25 +
    sentiment_impact * 0.2 +
    business_alignment * 0.15 +
    implementation_feasibility * 0.1
)
```

### Churn Prediction
- Identify language patterns of churning users
- Correlate feedback with usage decline
- Score users by churn risk
- Recommend retention interventions

## Section-by-Section Guidance

### Executive Summary
AI should synthesize:
- Overall health of user satisfaction
- Most significant changes since last period
- Critical issues requiring immediate attention
- Success stories to celebrate and amplify

### Key Findings
Focus on insights that:
- Appear across multiple feedback channels
- Affect significant user segments
- Correlate with business metrics
- Suggest clear actions

### Feature Requests
Prioritize by:
- Number of unique users requesting
- Strategic alignment with product vision
- Technical feasibility assessment
- Competitive differentiation potential

### Pain Points
Analyze through:
- Root cause analysis (5 Whys)
- Impact on user journey
- Correlation with support costs
- Relationship to churn

### Sentiment Analysis
Segment by:
- User persona/segment
- Feature area
- Customer lifecycle stage
- Geographic region
- Pricing tier

## Advanced Analytics

### Emotion Detection
Beyond positive/negative:
- Frustration vs confusion
- Delight vs satisfaction
- Urgency vs patience
- Trust vs skepticism

### Intent Classification
Understand what users want:
- Bug report
- Feature request
- Praise
- Question
- Complaint
- Suggestion

### Comparative Analysis
- Benchmark against competitors
- Track against industry standards
- Compare with previous periods
- Analyze by cohort

## Integration Points

### Connect to Metrics
- Link sentiment to NPS scores
- Correlate with usage metrics
- Map to conversion rates
- Track against retention

### Feed to Planning
- Inform product roadmap
- Prioritize bug fixes
- Guide UX improvements
- Shape marketing messages

## Quality Checks

### Validation Requirements
- [ ] Statistical significance of patterns
- [ ] Balanced representation across segments
- [ ] Bias detection in AI analysis
- [ ] Human validation of AI insights
- [ ] Cross-reference with quantitative data

### Common Pitfalls

#### ❌ Vocal Minority Bias
**Bad**: Letting loudest voices drive decisions
**Good**: Weight feedback by user segment value

#### ❌ Recency Bias
**Bad**: Overweighting latest feedback
**Good**: Track trends over time

#### ❌ Sentiment Misclassification
**Bad**: Misreading sarcasm or context
**Good**: Human validation of edge cases

#### ❌ Analysis Paralysis
**Bad**: Endless analysis without action
**Good**: Time-boxed analysis with clear outcomes

## Response Strategies

### Closing the Loop
1. **Acknowledge**: Thank users for feedback
2. **Communicate**: Share what you're doing
3. **Deliver**: Ship improvements
4. **Follow-up**: Check if issues resolved
5. **Celebrate**: Share success stories

### Segment-Specific Responses
- **Promoters**: Amplify and reward
- **Passives**: Engage and convert
- **Detractors**: Rescue and retain

## Automation Opportunities

### Real-time Processing
- Stream feedback as it arrives
- Trigger alerts for critical issues
- Update sentiment scores continuously
- Auto-tag and categorize

### Intelligent Routing
- Route urgent issues to support
- Flag feature requests for product
- Identify PR opportunities
- Detect potential crises

### Predictive Actions
- Prevent churn before it happens
- Anticipate support volume spikes
- Forecast feature adoption
- Predict satisfaction trends

## Visualization Best Practices

### Effective Displays
- Word clouds for topic frequency
- Sentiment heat maps by feature
- Trend lines for satisfaction
- Sankey diagrams for user journeys
- Scatter plots for correlation

### Dashboard Integration
- Real-time sentiment widget
- Top issues tracker
- Feature request leaderboard
- Churn risk monitor

## Remember

Feedback Analysis is about **turning user voice into product value**. It should answer:
1. What do users love and hate?
2. What should we build next?
3. Who is at risk of churning?
4. How are we trending?
5. What blindspots do we have?

Great feedback analysis creates empathy, drives priorities, and keeps the team connected to users.