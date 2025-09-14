# Feature Registry Generation Prompt

Create and maintain a comprehensive registry of all features in the system.

## Storage Location

Store the feature registry at: `docs/01-frame/feature-registry.md`

## Purpose

The Feature Registry is the single source of truth for:
- Feature identification (FEAT-XXX IDs)
- Feature status tracking
- Dependency management
- Team ownership
- Cross-phase traceability

## When to Update

Update the Feature Registry when:
1. **New Feature Identified**: Assign next available FEAT-XXX ID
2. **Status Changes**: Feature moves to next phase
3. **Dependencies Discovered**: New relationships identified
4. **Ownership Changes**: Team or person responsible changes
5. **Feature Cancelled**: Mark as deprecated with reason

## ID Assignment Process

### Getting Next ID
1. Check the registry for highest assigned ID
2. Increment by 1
3. Format as FEAT-XXX (e.g., FEAT-007)
4. Never reuse retired IDs

### ID Reservation
- Can reserve ranges for teams (e.g., FEAT-100-199 for Team B)
- Document reserved ranges in registry

## Status Progression

Features typically progress through these statuses:
```
Draft → Specified → Designed → In Test → In Build → Built → Deployed
```

But can also be:
- **On Hold**: Temporarily paused
- **Cancelled**: Will not be built
- **Deprecated**: Being phased out

## Dependency Analysis

Before assigning a feature ID, analyze:
1. **Technical Dependencies**: What other features must exist?
2. **Data Dependencies**: What data from other features is needed?
3. **User Journey Dependencies**: What features must users complete first?

Document all dependencies to prevent circular references and ensure buildable order.

## Feature Categorization

Group features logically:
- By domain (Auth, Payment, Reporting)
- By user type (Admin, Customer, Partner)
- By platform (Web, Mobile, API)
- By priority (P0, P1, P2)

This helps with:
- Resource allocation
- Sprint planning
- Architecture decisions
- Testing strategies

## Quality Checklist

Before adding a feature to the registry:
- [ ] Clear, descriptive name
- [ ] Brief but complete description
- [ ] Correct status
- [ ] Priority assigned (P0/P1/P2)
- [ ] Owner identified
- [ ] Dependencies documented
- [ ] Category assigned

## Common Patterns

### Feature Bundles
Sometimes multiple related features are identified together:
```
FEAT-010: E-commerce Platform
├── FEAT-011: Product Catalog
├── FEAT-012: Shopping Cart
├── FEAT-013: Checkout Process
└── FEAT-014: Order Management
```

Document parent-child relationships when applicable.

### Cross-Cutting Features
Some features affect multiple areas:
```
FEAT-020: Audit Logging (affects all features)
FEAT-021: Rate Limiting (affects all APIs)
```

Mark these as "Infrastructure" or "Cross-cutting" category.

## Registry Maintenance

### Weekly Updates
- Review feature status
- Update progress
- Check for new dependencies
- Assign new features

### Monthly Review
- Archive completed features
- Analyze delivery metrics
- Identify bottlenecks
- Plan upcoming features

### Quarterly Planning
- Reserve ID ranges
- Deprecate old features
- Major status updates
- Dependency audit

## Integration with Workflow

The Feature Registry connects all phases:

1. **Frame Phase**: Feature gets ID and specification
2. **Design Phase**: Status → "Designed" when complete
3. **Test Phase**: Status → "In Test" when tests written
4. **Build Phase**: Status → "In Build" during implementation
5. **Deploy Phase**: Status → "Deployed" when released

## Traceability

Every feature ID (FEAT-XXX) should be traceable to:
- User Stories (US-XXX)
- Solution Designs (SD-XXX)
- API Contracts (API-XXX)
- ADRs (ADR-XXX) when applicable
- Test Suites (tests/FEAT-XXX/)
- Source Code (src/features/FEAT-XXX/)

## Anti-Patterns to Avoid

### ❌ Vague Features
**Bad**: "Improve Performance"
**Good**: "FEAT-025: API Response Caching"

### ❌ Feature Sprawl
**Bad**: 200 tiny features
**Good**: Logical groupings with sub-features

### ❌ Missing Dependencies
**Bad**: Starting FEAT-030 without required FEAT-029
**Good**: Clear dependency chain documented

### ❌ Zombie Features
**Bad**: Features in "Draft" for 6 months
**Good**: Regular review and cleanup

## Success Metrics

A well-maintained registry enables:
- Quick feature status checks
- Accurate progress reporting
- Efficient resource allocation
- Clear communication with stakeholders
- Smooth handoffs between teams

---
*The Feature Registry is critical for project coordination. Keep it accurate and current.*