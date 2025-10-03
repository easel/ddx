# Story Refinement Log: {{STORY_ID}}-refinement-{{REFINEMENT_NUMBER}}

**Story ID**: {{STORY_ID}}
**Refinement Number**: {{REFINEMENT_NUMBER}}
**Date**: {{CURRENT_DATE}}
**Refinement Type**: {{REFINEMENT_TYPE}} <!-- bugs|requirements|enhancement|mixed -->
**Triggered By**: {{TRIGGER_REASON}} <!-- implementation|testing|stakeholder-feedback|production-issue -->
**Refinement Lead**: {{REFINEMENT_LEAD}}
**Status**: {{STATUS}} <!-- in-progress|completed|deferred -->

## Executive Summary

**Refinement Overview**: {{BRIEF_DESCRIPTION}}

**Primary Impact**: {{IMPACT_SUMMARY}}

**Phase Scope**: {{AFFECTED_PHASES}} <!-- frame|design|test|build|deploy|iterate -->

## Original Story State

### Story Reference
- **Original Story**: [{{STORY_ID}} - {{STORY_TITLE}}](../01-frame/user-stories/{{STORY_ID}}.md)
- **Story Version**: {{ORIGINAL_VERSION}}
- **Last Updated**: {{STORY_LAST_UPDATED}}
- **Implementation Status**: {{IMPLEMENTATION_STATUS}}

### Key Acceptance Criteria (Original)
{{#each ORIGINAL_ACCEPTANCE_CRITERIA}}
- **{{ac_id}}**: {{ac_description}}
{{/each}}

### Related Documents (Pre-Refinement)
- Frame Phase: {{FRAME_DOCUMENTS}}
- Design Phase: {{DESIGN_DOCUMENTS}}
- Test Phase: {{TEST_DOCUMENTS}}
- Build Phase: {{BUILD_DOCUMENTS}}

## Issues Identified

{{#each IDENTIFIED_ISSUES}}
### Issue {{issue_number}}: {{issue_title}}

**Category**: {{issue_category}} <!-- bug|gap|enhancement|constraint -->
**Priority**: {{issue_priority}} <!-- critical|high|medium|low -->
**Discovery Context**: {{discovery_context}}

**Description**:
{{issue_description}}

**Root Cause Analysis**:
{{root_cause}}

**Impact Assessment**:
- **User Impact**: {{user_impact}}
- **Technical Impact**: {{technical_impact}}
- **Business Impact**: {{business_impact}}
- **Risk Level**: {{risk_level}}

**Affected Acceptance Criteria**:
{{#each affected_ac}}
- {{ac_id}}: {{impact_description}}
{{/each}}

**Supporting Evidence**:
{{#each evidence}}
- {{evidence_item}}
{{/each}}

{{/each}}

## Refinement Analysis

### Requirements Impact Analysis

**Scope Evaluation**:
- **Within Original Scope**: {{in_scope_changes}}
- **Scope Extensions**: {{scope_extensions}}
- **Out of Scope**: {{out_of_scope_items}}

**Dependency Analysis**:
- **Affected Stories**: {{affected_stories}}
- **New Dependencies**: {{new_dependencies}}
- **Broken Dependencies**: {{broken_dependencies}}

**Backwards Compatibility**:
- **Breaking Changes**: {{breaking_changes}}
- **API Impact**: {{api_impact}}
- **Migration Required**: {{migration_needed}}

### Stakeholder Input

**Consultation Process**:
{{#each stakeholder_consultations}}
- **Stakeholder**: {{stakeholder_name}}
- **Role**: {{stakeholder_role}}
- **Input Date**: {{consultation_date}}
- **Key Feedback**: {{feedback_summary}}
- **Decisions Made**: {{decisions}}
{{/each}}

**Authority and Approvals**:
- **Refinement Authority**: {{refinement_authority}}
- **Approvals Required**: {{approvals_needed}}
- **Approval Status**: {{approval_status}}

## Refinement Resolutions

{{#each REFINEMENT_RESOLUTIONS}}
### Resolution {{resolution_number}}: {{resolution_title}}

**Addresses Issue(s)**: {{addressed_issues}}

**Resolution Strategy**: {{strategy}} <!-- harmonize|split-story|defer|clarify|enhance -->

**Detailed Resolution**:
{{resolution_description}}

**Requirements Changes**:
{{#each requirement_changes}}
- **Type**: {{change_type}} <!-- add|modify|remove|clarify -->
- **Target**: {{target_requirement}}
- **Original**: {{original_text}}
- **Refined**: {{refined_text}}
- **Rationale**: {{change_rationale}}
{{/each}}

**Implementation Impact**:
- **Code Changes Required**: {{code_changes}}
- **Test Updates Needed**: {{test_updates}}
- **Documentation Updates**: {{doc_updates}}
- **Effort Estimate**: {{effort_estimate}}

**Risk Mitigation**:
{{#each risk_mitigations}}
- **Risk**: {{risk_description}}
- **Mitigation**: {{mitigation_strategy}}
{{/each}}

{{/each}}

## Phase-Specific Updates

### Frame Phase Updates

**User Story Modifications**:
{{#each story_modifications}}
- **Section**: {{section_name}}
- **Change Type**: {{change_type}}
- **Original**: {{original_content}}
- **Updated**: {{updated_content}}
- **Rationale**: {{update_rationale}}
{{/each}}

**New Requirements Added**:
{{#each new_requirements}}
- **Requirement**: {{requirement_text}}
- **Source**: {{requirement_source}}
- **Priority**: {{requirement_priority}}
{{/each}}

**Requirements Documentation Updates**:
{{#each frame_doc_updates}}
- **Document**: {{document_path}}
- **Update Type**: {{update_type}}
- **Summary**: {{update_summary}}
{{/each}}

### Design Phase Updates

**Architecture Changes**:
{{#each architecture_changes}}
- **Component**: {{component_name}}
- **Change**: {{change_description}}
- **Impact**: {{architecture_impact}}
- **ADR Required**: {{adr_needed}}
{{/each}}

**API Contract Updates**:
{{#each api_updates}}
- **API**: {{api_name}}
- **Change**: {{api_change}}
- **Breaking**: {{is_breaking}}
- **Version Impact**: {{version_impact}}
{{/each}}

**Design Documentation Updates**:
{{#each design_doc_updates}}
- **Document**: {{document_path}}
- **Update Type**: {{update_type}}
- **Summary**: {{update_summary}}
{{/each}}

### Test Phase Updates

**New Test Cases**:
{{#each new_test_cases}}
- **Test**: {{test_name}}
- **Category**: {{test_category}}
- **Purpose**: {{test_purpose}}
- **Priority**: {{test_priority}}
{{/each}}

**Modified Test Cases**:
{{#each modified_test_cases}}
- **Test**: {{test_name}}
- **Change**: {{test_change}}
- **Reason**: {{change_reason}}
{{/each}}

**Test Documentation Updates**:
{{#each test_doc_updates}}
- **Document**: {{document_path}}
- **Update Type**: {{update_type}}
- **Summary**: {{update_summary}}
{{/each}}

### Build Phase Updates

**Implementation Plan Changes**:
{{#each implementation_changes}}
- **Component**: {{component_name}}
- **Change**: {{implementation_change}}
- **Effort**: {{change_effort}}
- **Risk**: {{change_risk}}
{{/each}}

**Technical Debt**:
- **Created**: {{technical_debt_created}}
- **Resolved**: {{technical_debt_resolved}}
- **Management Plan**: {{debt_management}}

**Build Documentation Updates**:
{{#each build_doc_updates}}
- **Document**: {{document_path}}
- **Update Type**: {{update_type}}
- **Summary**: {{update_summary}}
{{/each}}

## Validation and Quality Assurance

### Consistency Validation

**Cross-Phase Consistency**:
- [ ] Frame requirements align with refined story
- [ ] Design documents reflect requirement changes
- [ ] Test specifications cover all refined requirements
- [ ] Build plans address all requirement changes
- [ ] No orphaned references or broken links

**Requirement Validation**:
- [ ] All refined requirements are clear and testable
- [ ] No conflicts between original and refined requirements
- [ ] Acceptance criteria remain SMART (Specific, Measurable, Achievable, Relevant, Time-bound)
- [ ] Edge cases and error scenarios properly addressed

**Documentation Validation**:
- [ ] All referenced documents updated
- [ ] Cross-references between phases verified
- [ ] Traceability maintained end-to-end
- [ ] Version control properly managed

### Impact Assessment

**Project Impact**:
- **Timeline Impact**: {{timeline_impact}}
- **Resource Impact**: {{resource_impact}}
- **Budget Impact**: {{budget_impact}}
- **Quality Impact**: {{quality_impact}}

**Team Impact**:
- **Communication Completed**: {{team_communication}}
- **Training Required**: {{training_needed}}
- **Process Changes**: {{process_changes}}

**External Impact**:
- **Customer Impact**: {{customer_impact}}
- **Partner Impact**: {{partner_impact}}
- **Regulatory Impact**: {{regulatory_impact}}

## Lessons Learned

### Process Insights

**What Worked Well**:
{{#each process_positives}}
- {{positive_item}}
{{/each}}

**What Could Be Improved**:
{{#each process_improvements}}
- **Issue**: {{improvement_issue}}
- **Suggestion**: {{improvement_suggestion}}
{{/each}}

**Root Cause Prevention**:
{{#each prevention_measures}}
- **Root Cause**: {{root_cause}}
- **Prevention Strategy**: {{prevention_strategy}}
- **Implementation**: {{prevention_implementation}}
{{/each}}

### Refinement Quality

**Refinement Effectiveness**:
- **Issues Fully Resolved**: {{issues_resolved}}
- **Partial Resolutions**: {{partial_resolutions}}
- **Deferred Items**: {{deferred_items}}
- **New Issues Discovered**: {{new_issues}}

**Documentation Quality**:
- **Clarity Improved**: {{clarity_improvement}}
- **Completeness Enhanced**: {{completeness_enhancement}}
- **Consistency Achieved**: {{consistency_achievement}}

## Next Actions

### Immediate Actions
{{#each immediate_actions}}
- [ ] **Action**: {{action_description}}
- **Owner**: {{action_owner}}
- **Due Date**: {{action_due_date}}
- **Dependencies**: {{action_dependencies}}
{{/each}}

### Follow-up Items
{{#each followup_items}}
- [ ] **Item**: {{item_description}}
- **Timeline**: {{item_timeline}}
- **Owner**: {{item_owner}}
- **Related Stories**: {{related_stories}}
{{/each}}

### Monitoring and Review

**Success Metrics**:
{{#each success_metrics}}
- **Metric**: {{metric_name}}
- **Target**: {{metric_target}}
- **Measurement**: {{metric_measurement}}
{{/each}}

**Review Schedule**:
- **Next Review Date**: {{next_review_date}}
- **Review Participants**: {{review_participants}}
- **Review Criteria**: {{review_criteria}}

## Appendices

### A. Supporting Documents
{{#each supporting_documents}}
- [{{document_title}}]({{document_path}}) - {{document_description}}
{{/each}}

### B. Communication Records
{{#each communication_records}}
- **Date**: {{communication_date}}
- **Type**: {{communication_type}}
- **Participants**: {{participants}}
- **Summary**: {{communication_summary}}
- **Decisions**: {{decisions_made}}
{{/each}}

### C. Technical Analysis
{{technical_analysis_content}}

### D. Risk Register
{{#each risks}}
- **Risk**: {{risk_description}}
- **Probability**: {{risk_probability}}
- **Impact**: {{risk_impact}}
- **Mitigation**: {{risk_mitigation}}
- **Owner**: {{risk_owner}}
{{/each}}

---

**Refinement Completion**:
- **Completed By**: {{completion_lead}}
- **Completion Date**: {{completion_date}}
- **Quality Review**: {{quality_reviewer}}
- **Approval**: {{final_approval}}

**Related Refinements**:
- **Previous**: [{{STORY_ID}}-refinement-{{PREVIOUS_NUMBER}}]({{STORY_ID}}-refinement-{{PREVIOUS_NUMBER}}.md)
- **Next**: [{{STORY_ID}}-refinement-{{NEXT_NUMBER}}]({{STORY_ID}}-refinement-{{NEXT_NUMBER}}.md)

**Story Status After Refinement**: {{FINAL_STORY_STATUS}}