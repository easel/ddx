# Risk Register Generation Prompt

Create a comprehensive risk register to identify, assess, and plan mitigation for all project risks.

## Storage Location

Store the risk register at: `docs/01-frame/risk-register.md`

This central location ensures risks are visible and actively managed throughout the project.

## Purpose

The Risk Register:
- Identifies potential threats and opportunities
- Assesses probability and impact
- Defines mitigation strategies
- Assigns ownership and accountability
- Tracks risk status over time

## Key Principles

### 1. Comprehensive Identification
- Consider all risk categories
- Include both threats and opportunities
- Think beyond obvious risks
- Consider cumulative and cascading risks

### 2. Honest Assessment
- Be realistic about probabilities
- Don't downplay impacts
- Consider worst-case scenarios
- Acknowledge uncertainty

### 3. Active Management
- Assign clear ownership
- Define specific actions
- Set review frequencies
- Track effectiveness

## Risk Categories to Consider

### Technical Risks
- Technology choices
- Integration complexity
- Performance requirements
- Security vulnerabilities
- Technical debt
- Scalability challenges

### Business Risks
- Market changes
- Competitive threats
- ROI uncertainty
- Stakeholder alignment
- Regulatory compliance
- Reputation impact

### Resource Risks
- Team availability
- Skill gaps
- Budget constraints
- Vendor dependencies
- Equipment/infrastructure

### Project Management Risks
- Schedule compression
- Scope creep
- Requirement changes
- Communication breakdown
- Decision delays

### External Risks
- Economic conditions
- Legal/regulatory changes
- Third-party failures
- Natural disasters
- Geopolitical events

## Assessment Guidelines

### Probability Assessment
Consider:
- Historical data from similar projects
- Expert judgment
- Industry benchmarks
- Current project conditions
- External factors

### Impact Assessment
Evaluate impact on:
- **Schedule**: Delays to timeline
- **Budget**: Cost overruns
- **Quality**: Defects, technical debt
- **Scope**: Feature reduction
- **Reputation**: Brand damage
- **Team**: Morale, retention

### Risk Scoring
- Use consistent scales (1-5 recommended)
- Calculate Risk Score = Probability × Impact
- Define clear thresholds for action
- Consider qualitative factors too

## Mitigation Strategy Development

### For Each High/Critical Risk

1. **Prevention** (Reduce Probability)
   - What can we do to prevent this?
   - How can we detect early warning signs?
   - What controls can we implement?

2. **Mitigation** (Reduce Impact)
   - How can we minimize damage?
   - What contingencies are needed?
   - What fallback options exist?

3. **Transfer** (Share Risk)
   - Can we insure against this?
   - Can we contractually transfer?
   - Can we partner to share risk?

4. **Acceptance** (Prepared Response)
   - What triggers acceptance?
   - What's our response plan?
   - What reserves are needed?

## Quality Checklist

Before completing the risk register:
- [ ] All major risk categories explored
- [ ] Each risk has an owner
- [ ] Mitigation strategies are specific
- [ ] Triggers/indicators identified
- [ ] Review schedule established
- [ ] Escalation criteria defined
- [ ] Budget impact estimated

## Common Pitfalls to Avoid

### ❌ Generic Risks
**Bad**: "Project might be delayed"
**Good**: "Integration with legacy system may take 2-3 weeks longer due to undocumented APIs"

### ❌ Vague Mitigation
**Bad**: "We'll manage this carefully"
**Good**: "Assign dedicated integration specialist, create API documentation, daily integration testing"

### ❌ Orphan Risks
**Bad**: No owner assigned
**Good**: Specific person accountable with clear actions

### ❌ Static Register
**Bad**: Created once, never updated
**Good**: Weekly reviews, status updates, trend tracking

## Questions to Answer

1. **Identification**
   - What could go wrong?
   - What dependencies exist?
   - What assumptions are we making?
   - What's outside our control?
   - What opportunities exist?

2. **Assessment**
   - How likely is this to occur?
   - What would the impact be?
   - Can we detect it early?
   - Is it getting better or worse?

3. **Response**
   - Who owns this risk?
   - What's our strategy?
   - What specific actions will we take?
   - How will we know it's working?
   - When do we escalate?

## Risk Register Maintenance

### Weekly Tasks
- Review high/critical risks
- Update risk status
- Identify new risks
- Close resolved risks

### Monthly Tasks
- Full register review
- Trend analysis
- Strategy effectiveness
- Budget impact update

### Quarterly Tasks
- Executive briefing
- Lessons learned
- Strategy adjustment
- Reserve analysis

## Using with AI

When working with AI to identify risks:

```bash
# Provide project context
"We're building a [system type] for [users]"
"Key dependencies include [list]"
"Critical success factors are [list]"

# Ask for specific categories
"What technical risks should we consider?"
"What regulatory risks apply?"
"What are common risks for similar projects?"

# Request mitigation strategies
"For risk X, what are proven mitigation approaches?"
"What early warning signs should we monitor?"
```

## Remember

Risk management is about:
1. **Preparation**, not paranoia
2. **Action**, not just documentation  
3. **Learning**, not blame
4. **Opportunity**, not just threats

A good risk register helps the team sleep better at night, knowing they're prepared for challenges ahead.