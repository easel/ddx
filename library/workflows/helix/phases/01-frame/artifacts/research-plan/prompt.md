# Research Plan Generation Prompt

## Context
You are helping create a research plan for the Frame phase of the HELIX workflow. This is used when there are significant unknowns, unclear requirements, or need for investigation before proceeding with standard PRD development.

## Objective
Create a comprehensive research plan that transforms uncertainty into actionable insights for product development. The research should be time-boxed, evidence-based, and directly inform subsequent Frame phase artifacts.

## Key Principles

### 1. Research-Driven Clarity
- Transform vague problems into specific research questions
- Focus on actionable insights, not academic knowledge
- Ensure every question ties to a product decision

### 2. Appropriate Methods
- Match research methods to question types:
  - **User research**: Interviews, surveys, observations for user needs
  - **Market research**: Competitive analysis, sizing for business viability
  - **Technical research**: Spikes, prototypes for feasibility
  - **Domain research**: Expert interviews for business rules

### 3. Time-Boxed Investigation
- Set strict time limits for all research activities
- Prioritize high-impact, low-effort investigations first
- Plan for "good enough" answers rather than perfect knowledge

### 4. Evidence-Based Conclusions
- Require evidence for all findings
- Document methodology and limitations
- Separate facts from assumptions and opinions

## Input Analysis

Analyze the provided context to identify:

### Research Triggers
Look for these indicators that research is needed:
- "We think users want..." (assumption, needs validation)
- "The market might..." (uncertainty, needs investigation)
- "It should be technically possible..." (feasibility question)
- "Similar to [competitor]..." (needs competitive analysis)
- "Users complain about..." (needs deeper investigation)

### Knowledge Gaps
Identify what we don't know that prevents moving forward:
- **User needs**: Who are they? What problems do they have?
- **Market opportunity**: Is there demand? How much?
- **Technical feasibility**: Can we build this? How complex?
- **Business viability**: Will it generate value? At what cost?
- **Competitive landscape**: What exists? What's our differentiator?

### Risk Areas
Highlight areas where wrong assumptions could be costly:
- Large development investment without user validation
- Technical approach without feasibility verification
- Market entry without competitive understanding
- Resource allocation without scope clarity

## Research Plan Structure

### Research Objectives
Create 3-5 specific, measurable research questions:
- Each question should inform a specific product decision
- Success criteria must be measurable and time-bound
- Questions should be answerable with available resources

**Good Example**: "What are the top 3 pain points in current CLI workflow management tools among senior developers at startups (50-200 employees)?"

**Bad Example**: "What do developers want in CLI tools?" (too broad, unmeasurable)

### Method Selection
Choose appropriate research methods:

**User Research**:
- Interviews: Deep qualitative insights (5-8 participants)
- Surveys: Broad quantitative validation (50+ participants)
- Observations: Actual vs. reported behavior
- Jobs-to-be-Done: Understand user motivations

**Market Research**:
- Competitive analysis: Feature comparison, positioning
- Market sizing: TAM/SAM analysis, demand validation
- Trend analysis: Technology and market direction

**Technical Research**:
- Literature review: Existing solutions and approaches
- Prototyping: Feasibility and complexity validation
- Benchmarking: Performance and scalability analysis
- Architecture spikes: Technical approach validation

**Domain Research**:
- Expert interviews: Subject matter expertise
- Process analysis: Current state workflows
- Documentation review: Existing standards and practices

### Timeline and Resource Planning
- **Total duration**: 1-4 weeks maximum
- **Daily time allocation**: Specific hours per day
- **Milestone checkpoints**: Weekly progress reviews
- **Resource requirements**: People, tools, budget
- **Risk mitigation**: Backup plans for delays

### Success Criteria
Define what "done" looks like:
- All research questions answered with evidence
- Confidence level achieved for key decisions
- Stakeholder alignment on findings
- Clear recommendations for next steps

## Output Guidelines

### Research Questions Format
Use this structure for each research question:

```markdown
**Research Question**: [Specific, measurable question]
- **Decision Impact**: [What product decision this informs]
- **Success Criteria**: [How to know we have enough information]
- **Method**: [How we'll investigate this]
- **Timeline**: [When this will be completed]
```

### Method Description Format
For each research method:

```markdown
**Method**: [Name of method]
- **Objective**: [What specific question this addresses]
- **Approach**: [Step-by-step execution plan]
- **Participants/Sources**: [Who or what we'll study]
- **Duration**: [Time required]
- **Deliverable**: [Specific output produced]
```

### Risk Documentation
Document potential research risks:
- **Timeline risks**: What could delay research?
- **Quality risks**: What could compromise findings?
- **Resource risks**: What dependencies exist?
- **Scope risks**: What could expand beyond plan?

## Quality Checks

Before finalizing the research plan, verify:

### Completeness
- [ ] All major unknowns addressed
- [ ] Methods appropriate for questions
- [ ] Timeline realistic and bounded
- [ ] Resources clearly defined
- [ ] Success criteria measurable

### Actionability
- [ ] Questions inform specific decisions
- [ ] Methods will produce usable insights
- [ ] Timeline allows for action on findings
- [ ] Budget justified by decision impact

### Feasibility
- [ ] Participant recruitment realistic
- [ ] Timeline accounts for delays
- [ ] Skills available in team
- [ ] Tools and resources accessible

## Anti-Patterns to Avoid

### Research Theater
- Don't research to delay decisions
- Avoid questions that don't inform action
- Don't use research to confirm pre-existing beliefs

### Analysis Paralysis
- Set strict time limits
- Plan for "good enough" vs. perfect information
- Focus on highest-impact uncertainties first

### Method Misalignment
- Don't use surveys for complex qualitative insights
- Don't rely on interviews for quantitative validation
- Don't use technical spikes for business questions

### Scope Creep
- Stick to original research questions
- Resist expanding scope during execution
- Time-box all research activities

## Success Metrics

A good research plan should:
- Transform uncertainty into actionable insights
- Be completable within 1-4 weeks
- Directly inform Frame phase artifacts
- Provide evidence for key product decisions
- Reduce risk of wrong assumptions

## Example Research Questions by Type

### User Research
- "What are the top 3 workflow bottlenecks for CLI users in our target segment?"
- "How do senior developers currently share configuration across projects?"
- "What evidence exists that our target users would pay for this solution?"

### Technical Research
- "Can we achieve <2 second startup time with our proposed architecture?"
- "What are the integration complexities with the top 5 CI/CD platforms?"
- "How do similar tools handle cross-platform distribution and updates?"

### Market Research
- "How do existing solutions monetize and at what scale?"
- "What gaps exist in current competitive offerings?"
- "What adoption patterns exist for developer productivity tools?"

## Integration with HELIX

### Pre-Research
- Research plan should be triggered by uncertainty in problem definition
- Must have stakeholder alignment on research necessity
- Should have clear decision points that depend on findings

### Post-Research
- Findings directly inform PRD development
- Validated assumptions become foundation for requirements
- Technical insights guide principles and constraints
- User insights shape personas and user stories

### Transition Criteria
Research is complete when:
- All research questions answered with sufficient confidence
- Key assumptions validated or invalidated with evidence
- Stakeholders aligned on findings and implications
- Clear path forward defined for Frame phase

Remember: Research should enable decision-making, not delay it. Good research transforms uncertainty into confidence for product development decisions.