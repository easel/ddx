# Technical Spike Generation Prompt

## Context
You are helping create a technical spike for the Design phase of the HELIX workflow. Technical spikes are time-boxed investigations used to explore technical unknowns, validate architectural approaches, or reduce implementation risk before committing to detailed design decisions.

**Primary Focus**: Technology evaluation and selection to implement architectural decisions made in ADRs.

## Objective
Create a focused technical investigation that transforms uncertainty into actionable insights. The spike should be strictly time-boxed, evidence-based, and produce concrete recommendations that inform architectural decisions.

## Key Principles

### 1. Time-Boxed Investigation
- Set strict time boundaries (typically 1-5 days)
- Focus on answering specific technical questions
- Stop when time is up, regardless of completeness
- Prefer "good enough" answers to perfect knowledge

### 2. Evidence-Based Analysis
- Base all conclusions on concrete evidence
- Provide measurable results where possible
- Document methodology and limitations
- Separate facts from opinions and assumptions

### 3. Decision-Focused Outcomes
- Answer specific technical questions that block progress
- Provide clear recommendations with rationale
- Identify risks and trade-offs explicitly
- Enable confident architectural decisions

### 4. Minimal Viable Investigation
- Do the minimum work necessary to answer the question
- Create throwaway artifacts (prototypes, tests, etc.)
- Focus on learning, not building production code
- Preserve artifacts for future reference

## Spike Planning Framework

### Technical Question Analysis
Identify the specific uncertainty that needs investigation:

**Question Types**:
- **Technology Selection**: "Which library/framework should we use for X?" (most common)
- **Performance Validation**: "Can technology X handle our performance requirements?"
- **Integration Assessment**: "How complex is integrating with system Y?"
- **Feasibility Study**: "Is it possible to implement feature A with technology B?"
- **Compatibility Testing**: "Does framework X work with our existing stack?"

**Tech Spike Scope (What Belongs Here)**:
- Comparing specific libraries or frameworks (e.g., "Caliban vs Sangria for GraphQL")
- Performance benchmarking of technologies
- Integration complexity assessment
- Version compatibility testing
- License and ecosystem evaluation

**Not for Tech Spikes (Use Other Artifacts)**:
- Fundamental architectural decisions → **ADR** (e.g., "Should we use GraphQL?" → ADR)
- Implementation patterns and architecture → **Solution Design**
- Configuration and deployment → **Implementation Guide**

**Question Quality Criteria**:
- Specific and measurable
- Directly blocks architectural decisions
- Answerable within time budget
- Has clear success criteria

### Artifact Relationships

Tech Spikes work within this design artifact flow:

```
ADR (Why) → Tech Spike (What) → Solution Design (How)
```

**Typical Flow:**
1. **ADR decides approach**: "Use distributed caching for session management"
2. **Tech Spike selects technology**: "Redis vs Hazelcast vs Memcached evaluation"
3. **Solution Design defines implementation**: "Redis cluster architecture"

**Reference the supporting ADR** in your spike to maintain traceability. The architectural decision should already be made - your spike is choosing the specific technology to implement it.

### Spike Scope Definition

#### In Scope
Define what the spike will investigate:
- Specific technologies or approaches to evaluate
- Particular aspects to measure or test
- Concrete questions to answer
- Success criteria to meet

#### Out of Scope
Explicitly exclude:
- Production-ready implementations
- Comprehensive feature development
- Complete testing and hardening
- Performance optimization beyond measurement

#### Success Criteria
Define what "done" looks like:
- Specific questions answered
- Evidence level required
- Recommendations provided
- Risks identified

## Investigation Methodologies

### Prototype Development
**When to Use**: Validating architectural approaches, integration complexity
**Approach**:
- Build minimal working example
- Focus on proving/disproving hypothesis
- Use simplest possible implementation
- Measure specific characteristics

**Deliverables**:
- Working prototype code
- Performance measurements
- Integration complexity assessment
- Implementation time estimates

### Comparative Analysis
**When to Use**: Choosing between multiple technical options
**Approach**:
- Define evaluation criteria
- Implement minimal versions of each approach
- Measure and compare results
- Document trade-offs

**Deliverables**:
- Comparison matrix
- Trade-off analysis
- Recommendation with rationale
- Risk assessment for each option

### Benchmarking and Performance Testing
**When to Use**: Understanding performance characteristics, scalability limits
**Approach**:
- Create realistic test scenarios
- Measure key performance metrics
- Test under various load conditions
- Document performance characteristics

**Deliverables**:
- Performance benchmarks
- Scalability analysis
- Resource utilization data
- Performance recommendations

### Integration Testing
**When to Use**: Validating compatibility, understanding integration complexity
**Approach**:
- Test key integration points
- Validate data flow and compatibility
- Measure integration overhead
- Document integration requirements

**Deliverables**:
- Integration strategy
- Compatibility matrix
- Integration complexity assessment
- Required adaptations

### Expert Consultation
**When to Use**: Leveraging specialized knowledge, validating approaches
**Approach**:
- Prepare specific questions
- Consult with domain experts
- Validate findings with multiple sources
- Document expert recommendations

**Deliverables**:
- Expert insights summary
- Validated approaches
- Risk assessments
- Implementation recommendations

## Spike Execution Guidelines

### Time Management
- **Day 1**: Setup, initial investigation, early findings
- **Day 2**: Deep investigation, data collection, analysis
- **Day 3**: Synthesis, conclusions, recommendations
- **Final Hours**: Documentation, artifact preservation

### Investigation Structure

#### Setup Phase (10-20% of time)
- Environment setup
- Tool installation
- Initial research
- Approach refinement

#### Investigation Phase (60-70% of time)
- Execute planned activities
- Collect data and evidence
- Document findings as you go
- Adapt approach based on discoveries

#### Analysis Phase (15-20% of time)
- Synthesize findings
- Draw conclusions
- Formulate recommendations
- Identify follow-up needs

### Evidence Collection
Document everything that informs conclusions:

**Quantitative Evidence**:
- Performance measurements
- Resource utilization data
- Error rates and failure modes
- Time and complexity metrics

**Qualitative Evidence**:
- Integration complexity observations
- Developer experience insights
- Maintainability assessments
- Implementation challenges

**Artifact Evidence**:
- Working prototype code
- Test scripts and results
- Configuration examples
- Documentation samples

## Findings Documentation

### Structured Findings Format
For each significant discovery:

```markdown
**FINDING**: [Specific discovery statement]
- **Evidence**: [Concrete proof/data supporting finding]
- **Confidence Level**: High/Medium/Low
- **Implications**: [What this means for the project]
- **Risks**: [Any risks this finding introduces]
```

### Trade-off Analysis
Present choices systematically:

| Factor | Option A | Option B | Option C | Weight | Winner |
|--------|----------|----------|----------|--------|--------|
| Performance | Score/Data | Score/Data | Score/Data | High | Option X |
| Complexity | Assessment | Assessment | Assessment | Medium | Option Y |
| Maintainability | Assessment | Assessment | Assessment | High | Option Z |

### Risk Assessment
Document risks discovered during investigation:

| Risk | Probability | Impact | Evidence | Mitigation |
|------|-------------|--------|----------|------------|
| Specific risk | H/M/L | H/M/L | Concrete evidence | Specific strategy |

## Recommendation Formulation

### Recommendation Structure
```markdown
**RECOMMENDATION**: [Specific, actionable recommendation]
- **Rationale**: [Why this is the best choice based on evidence]
- **Confidence**: High/Medium/Low
- **Risks**: [Known risks with this approach]
- **Next Steps**: [What needs to happen to implement this]
- **Timeline**: [When this should be done]
```

### Decision Support
Provide decision makers with:
- Clear recommendation with confidence level
- Supporting evidence summary
- Risk assessment and mitigation strategies
- Implementation implications
- Timeline and resource impacts

## Quality Standards

### Evidence Quality
- All conclusions backed by concrete evidence
- Methodology documented and repeatable
- Limitations and confidence levels stated
- Assumptions made explicit

### Time Management
- Strict adherence to time budget
- Clear stopping criteria respected
- Investigation focused, not exhaustive
- Documentation concurrent with investigation

### Actionability
- Recommendations specific and implementable
- Next steps clearly defined
- Decision points identified
- Risks and trade-offs explicit

## Common Spike Patterns

### Architecture Validation Spike
**Objective**: Validate architectural approach feasibility
**Activities**:
- Implement key architectural components
- Test critical system interactions
- Measure performance characteristics
- Assess implementation complexity

### Technology Evaluation Spike
**Objective**: Compare technology alternatives
**Activities**:
- Implement same functionality with different technologies
- Benchmark performance and resource usage
- Assess learning curve and documentation quality
- Evaluate ecosystem maturity and support

### Integration Complexity Spike
**Objective**: Understand integration challenges
**Activities**:
- Implement integration with target systems
- Test data flow and format compatibility
- Measure integration overhead
- Document required adaptations

### Performance Characteristics Spike
**Objective**: Understand system performance profile
**Activities**:
- Create realistic load scenarios
- Measure response times and throughput
- Test scalability limits
- Identify performance bottlenecks

## Anti-Patterns to Avoid

### Scope Creep
- Don't expand investigation beyond original questions
- Resist urge to build production-ready code
- Stop when time budget is exhausted
- Focus on answering specific questions

### Analysis Paralysis
- Don't seek perfect information
- Make decisions based on available evidence
- Accept "good enough" answers within time constraints
- Move forward with reasonable confidence

### Solution Bias
- Don't investigate to confirm predetermined conclusions
- Remain open to unexpected findings
- Consider evidence that contradicts assumptions
- Report negative results honestly

### Insufficient Evidence
- Don't make recommendations without concrete support
- Avoid conclusions based solely on opinions
- Provide measurable evidence where possible
- Document methodology and limitations

## Success Indicators

A successful technical spike:
- Answers specific technical questions within time budget
- Provides evidence-based recommendations
- Reduces technical risk and uncertainty
- Enables confident architectural decisions
- Preserves knowledge for future reference

## Integration with HELIX Workflow

### Pre-Spike
- Spike should be triggered by specific technical uncertainty
- Must have clear decision points that depend on findings
- Should be approved by technical lead and stakeholder

### Post-Spike
- Findings inform architecture and solution design decisions
- Recommendations guide implementation planning
- Risks inform project planning and mitigation strategies
- Artifacts preserved for future reference and learning

### Decision Points
Spike results should directly inform:
- Architecture decisions (ADRs)
- Technology choices
- Implementation strategies
- Risk mitigation approaches
- Timeline and resource planning

Remember: The goal is learning, not building. Good technical spikes prevent bad architectural decisions by providing evidence for critical choices.