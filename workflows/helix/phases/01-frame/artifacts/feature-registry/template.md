# Feature Registry

**Document Type**: Feature Registry
**Status**: [Active | Archived]
**Last Updated**: [Date]
**Maintained By**: [Team/Person]

## Purpose

This registry tracks all features in the system, their status, dependencies, and ownership. It serves as the single source of truth for feature identification and tracking.

## Active Features

| ID | Name | Description | Status | Priority | Owner | Created | Updated |
|----|------|-------------|--------|----------|-------|---------|---------|
| FEAT-001 | [Feature Name] | [Brief description] | [Draft/Specified/Designed/Built/Deployed] | P0 | [Team/Person] | [Date] | [Date] |
| FEAT-002 | | | | | | | |
| FEAT-003 | | | | | | | |

## Feature Status Definitions

- **Draft**: Initial concept, requirements being gathered
- **Specified**: Feature specification complete (Frame phase done)
- **Designed**: Technical design complete (Design phase done)
- **In Test**: Tests being written (Test phase active)
- **In Build**: Implementation in progress (Build phase active)
- **Built**: Implementation complete, ready for deployment
- **Deployed**: Released to production
- **Deprecated**: No longer supported, scheduled for removal

## Feature Dependencies

Document which features depend on others:

| Feature | Depends On | Dependency Type | Notes |
|---------|------------|-----------------|-------|
| FEAT-002 | FEAT-001 | Required | Payment needs authenticated users |
| FEAT-003 | FEAT-001, FEAT-002 | Required | Reports need auth and payment data |

## Feature Categories

Group features by type or domain:

### Authentication & Security
- FEAT-001: User Authentication
- FEAT-004: Two-Factor Authentication

### Payment & Billing
- FEAT-002: Payment Processing
- FEAT-005: Subscription Management

### Reporting & Analytics
- FEAT-003: Report Generation
- FEAT-006: Real-time Analytics

## ID Assignment Rules

1. **Sequential Numbering**: Features are numbered sequentially (001, 002, 003...)
2. **Never Reuse IDs**: Once assigned, an ID is permanent even if feature is cancelled
3. **Three Digits**: Use format FEAT-XXX (e.g., FEAT-001, not FEAT-1)
4. **Reserve Ranges**: Optionally reserve ranges for different teams or categories

## Deprecated/Cancelled Features

Track features that were cancelled or deprecated:

| ID | Name | Status | Reason | Date |
|----|------|--------|--------|------|
| [FEAT-XXX] | [Name] | Cancelled | [Why cancelled] | [Date] |

## Feature Metrics

Track high-level metrics:

| Quarter | Features Planned | Features Delivered | Success Rate |
|---------|-----------------|-------------------|--------------|
| Q1 2024 | 5 | 4 | 80% |
| Q2 2024 | 6 | 6 | 100% |

## Cross-References

### Related Documents
- **PRD**: `docs/01-frame/prd.md` - Overall product vision
- **Principles**: `docs/01-frame/principles.md` - Guiding principles
- **Feature Specs**: `docs/01-frame/features/FEAT-XXX-[name].md`

### Artifact Locations by Feature
For each feature, artifacts are located at:
- **Specification**: `docs/01-frame/features/FEAT-XXX-[name].md`
- **User Stories**: `docs/01-frame/user-stories/US-XXX-[name].md`
- **Solution Design**: `docs/design/solution-designs/SD-XXX-[name].md`
- **Contracts**: `docs/design/contracts/API-XXX-[name].md`
- **Tests**: `tests/FEAT-XXX-[name]/`
- **Implementation**: `src/features/FEAT-XXX-[name]/`

## Maintenance Notes

- Review and update weekly during planning sessions
- Archive completed features quarterly
- Audit dependencies before starting new features
- Update status as features progress through phases

---
*This is a living document. Update it whenever features are added, modified, or their status changes.*