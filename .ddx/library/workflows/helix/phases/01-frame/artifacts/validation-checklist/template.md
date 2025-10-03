# Frame Phase Validation Checklist

**Phase**: Frame (01)
**Status**: [ ] Not Started | [ ] In Progress | [ ] Complete
**Validated By**: [Name]
**Validation Date**: [Date]
**Overall Result**: [ ] Pass | [ ] Pass with Conditions | [ ] Fail

## Purpose

This checklist ensures all Frame phase deliverables are complete, consistent, and ready for the Design phase.

## Validation Summary

| Category | Items | Completed | Pass Rate |
|----------|-------|-----------|----------|
| Documentation | XX | XX | XX% |
| Quality Gates | XX | XX | XX% |
| Stakeholder Approval | XX | XX | XX% |
| Technical Readiness | XX | XX | XX% |

## 1. Documentation Completeness

### Product Requirements Document (PRD)
- [ ] **Created**: PRD exists at `docs/helix/01-frame/prd.md`
- [ ] **Executive Summary**: Written and reflects current state
- [ ] **Problem Statement**: Clear and quantified
- [ ] **Success Metrics**: Specific, measurable, with targets
- [ ] **User Personas**: Based on research/data
- [ ] **Requirements**: Prioritized as P0/P1/P2
- [ ] **Risks**: Identified with mitigation strategies
- [ ] **Timeline**: Realistic milestones defined
- [ ] **Approval**: Signed off by stakeholders

**Notes**: [Any issues or conditions]

### Principles Document
- [ ] **Created**: Exists at `docs/helix/01-frame/principles.md`
- [ ] **Core Principles**: Defined and enforceable
- [ ] **Technology Constraints**: Specified
- [ ] **Quality Standards**: Clear and measurable
- [ ] **Exception Process**: Defined
- [ ] **Alignment**: Consistent with PRD

**Notes**: [Any issues or conditions]

### Feature Specifications
- [ ] **Coverage**: All P0 features have specifications
- [ ] **Location**: Stored in `docs/helix/01-frame/features/`
- [ ] **Naming**: Follow FEAT-XXX-title format
- [ ] **Completeness**: All sections filled
- [ ] **Clarifications**: No [NEEDS CLARIFICATION] in P0 specs
- [ ] **Testability**: All requirements are testable
- [ ] **NFRs**: Non-functional requirements defined

**Notes**: [Any issues or conditions]

### User Stories
- [ ] **Coverage**: P0 features have user stories
- [ ] **Location**: Stored in `docs/helix/01-frame/user-stories/`
- [ ] **Format**: Follow "As a... I want... So that..."
- [ ] **Acceptance Criteria**: Given/When/Then format
- [ ] **Definition of Done**: Specified for each story
- [ ] **Traceability**: Linked to features

**Notes**: [Any issues or conditions]

### Feature Registry
- [ ] **Created**: Exists at `docs/helix/01-frame/feature-registry.md`
- [ ] **Completeness**: All features registered
- [ ] **IDs**: Unique FEAT-XXX identifiers
- [ ] **Dependencies**: Documented
- [ ] **Ownership**: Assigned for each feature
- [ ] **Priority**: P0/P1/P2 classifications

**Notes**: [Any issues or conditions]

### Stakeholder Map
- [ ] **Created**: Exists at `docs/helix/01-frame/stakeholder-map.md`
- [ ] **Identification**: All stakeholders identified
- [ ] **RACI Matrix**: Complete and consistent
- [ ] **Communication Plan**: Defined
- [ ] **Engagement Strategy**: Specified per stakeholder
- [ ] **Contact Info**: Current and complete

**Notes**: [Any issues or conditions]

### Risk Register
- [ ] **Created**: Exists at `docs/helix/01-frame/risk-register.md`
- [ ] **Risk Identification**: Major risks identified
- [ ] **Assessment**: Probability and impact scored
- [ ] **Mitigation**: Strategies defined for high risks
- [ ] **Ownership**: Each risk has an owner
- [ ] **Monitoring**: Review schedule established

**Notes**: [Any issues or conditions]

## 2. Quality Gates

### Content Quality
- [ ] **Clarity**: Documents are clear and unambiguous
- [ ] **Consistency**: No contradictions between documents
- [ ] **Completeness**: No missing critical information
- [ ] **Measurability**: Metrics and requirements are quantified
- [ ] **Traceability**: Clear links between artifacts

### Technical Validation
- [ ] **Feasibility**: Requirements are technically achievable
- [ ] **Constraints**: Technical limitations acknowledged
- [ ] **Dependencies**: External dependencies identified
- [ ] **Assumptions**: Documented and validated

### Business Validation
- [ ] **Value**: Clear business value articulated
- [ ] **ROI**: Return on investment considered
- [ ] **Market Fit**: Solution addresses market need
- [ ] **Competitive**: Differentiation understood

## 3. Stakeholder Validation

### Approvals Obtained
- [ ] **Product Owner**: PRD approved
- [ ] **Technical Lead**: Feasibility confirmed
- [ ] **Key Stakeholders**: Buy-in secured
- [ ] **Sponsor**: Funding/resources committed

### Feedback Incorporated
- [ ] **Review Cycles**: Completed required reviews
- [ ] **Feedback Logged**: All feedback documented
- [ ] **Changes Made**: Feedback incorporated
- [ ] **Conflicts Resolved**: Disagreements addressed

## 4. Cross-Reference Validation

### Artifact Consistency
- [ ] **PRD → Features**: All PRD requirements have specs
- [ ] **Features → Stories**: Features have user stories
- [ ] **Stories → Acceptance**: Clear acceptance criteria
- [ ] **Registry → Features**: Registry matches specifications
- [ ] **Risks → Mitigation**: High risks have plans

### Naming Conventions
- [ ] **Feature IDs**: FEAT-XXX format consistent
- [ ] **Story IDs**: US-XXX format consistent
- [ ] **Risk IDs**: RISK-XXX format consistent
- [ ] **File Names**: Follow prescribed patterns

## 5. Exit Criteria Verification

### Must-Have Criteria
- [ ] **PRD Approved**: By all required parties
- [ ] **P0 Specified**: All P0 requirements detailed
- [ ] **Personas Validated**: Based on real data
- [ ] **Metrics Measurable**: With targets and methods
- [ ] **Principles Established**: And enforceable
- [ ] **Registry Initialized**: With all features
- [ ] **Stories Created**: For P0 features
- [ ] **Risks Assessed**: With mitigation plans
- [ ] **Stakeholders Aligned**: With RACI defined
- [ ] **Scope Clear**: In/out of scope defined

### Should-Have Criteria
- [ ] **Assumptions Documented**: In dedicated document
- [ ] **Glossary Created**: With key terms
- [ ] **Traceability Matrix**: Requirements to tests
- [ ] **Examples Provided**: For key artifacts

## 6. Common Issues Checklist

### Typical Problems to Check
- [ ] **No Zombie Requirements**: Remove dead requirements
- [ ] **No Conflicting Priorities**: P0s are truly critical
- [ ] **No Vague Terms**: "Fast", "Easy" are quantified
- [ ] **No Missing Owners**: Everything has accountability
- [ ] **No Unrealistic Timelines**: Schedule is achievable
- [ ] **No Hidden Dependencies**: All dependencies explicit
- [ ] **No Assumed Knowledge**: Context is provided

## 7. Readiness Assessment

### Design Phase Readiness
- [ ] **Requirements Stable**: Not expecting major changes
- [ ] **Team Available**: Design team identified
- [ ] **Tools Ready**: Design tools accessible
- [ ] **Knowledge Transfer**: Frame team can brief Design team

### Risk Assessment
- [ ] **Low Risk**: Proceeding to Design is low risk
- [ ] **Mitigation Ready**: For any remaining risks
- [ ] **Escalation Path**: Clear if issues arise

## 8. Validation Sign-Off

### Validation Team

| Role | Name | Signature | Date | Status |
|------|------|-----------|------|--------|
| Product Owner | [Name] | [Signature] | [Date] | [ ] Approved |
| Technical Lead | [Name] | [Signature] | [Date] | [ ] Approved |
| QA Lead | [Name] | [Signature] | [Date] | [ ] Approved |
| Project Manager | [Name] | [Signature] | [Date] | [ ] Approved |

### Conditions and Exceptions

**Conditions for Approval**:
1. [Any conditions that must be met]
2. [Any follow-up actions required]

**Approved Exceptions**:
1. [Any approved deviations from standard]
2. [Justification for exceptions]

### Overall Validation Result

- [ ] **PASS**: All criteria met, ready for Design phase
- [ ] **CONDITIONAL PASS**: Proceed with specific conditions
- [ ] **FAIL**: Must address issues before proceeding

**Validation Notes**:
[Overall assessment and any additional context]

## 9. Action Items

### Before Design Phase

| Action | Owner | Due Date | Status |
|--------|-------|----------|--------|
| [Action needed] | [Name] | [Date] | [ ] Open |

### For Design Phase

| Item | Description | Owner |
|------|-------------|-------|
| [Carryover item] | [What needs attention in Design] | [Who will handle] |

## 10. Lessons Learned

### What Went Well
1. [Positive outcome or practice]
2. [Effective approach to replicate]

### Areas for Improvement
1. [Issue encountered]
2. [Suggested improvement]

### Recommendations for Future Projects
1. [Key learning to apply]
2. [Process improvement suggestion]

---

**Validation Completed By**: [Name]
**Date**: [Date]
**Next Review**: [When to revalidate if needed]

*This checklist should be completed before formally transitioning to the Design phase.*