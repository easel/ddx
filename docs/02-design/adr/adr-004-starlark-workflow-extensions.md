# ADR-004: Embed Starlark for Workflow Extensions and Validators

**Date**: 2025-01-12
**Status**: Accepted
**Deciders**: DDX Development Team
**Related Feature(s)**: Workflow Extensions and Validation
**Confidence Level**: High

## Context

DDX needs to provide extensible workflow automation that allows users to define custom validation rules, enforce project-specific policies, create workflow extensions without modifying DDX core, and implement complex conditional logic while maintaining security through sandboxing.

### Problem Statement

How do we enable user-customizable workflow extensions and validators that are secure, performant, easy to write, and portable across platforms without allowing arbitrary code execution that could compromise system security?

### Current State

Most tools either don't support customization, require plugins compiled in the same language, use shell scripts with security risks, or embed full scripting languages like JavaScript or Python with complex sandboxing requirements.

### Requirements Driving This Decision
- Enable custom validation rules and workflow logic
- Sandbox execution to prevent security issues
- Fast execution for pre-commit hooks (<10ms)
- Familiar syntax for developers
- Deterministic execution for reproducibility
- Work across all platforms
- Integrate cleanly with Go CLI

## Decision

We will embed Starlark (a Python-like configuration language) for workflow extensions, validators, and custom automation within DDX.

### Key Points
- Starlark provides Python-like syntax familiar to developers
- Sandboxed by design with no filesystem/network access
- Deterministic execution guarantees reproducibility
- Native Go implementation for easy embedding
- Battle-tested in Bazel and Buck build systems
- Designed specifically for configuration use cases
- Faster than embedded interpreters like Lua

## Alternatives Considered

### Option 1: Lua
**Description**: Embed Lua scripting language for extensions

**Pros**:
- Mature embedding language
- Small footprint
- Fast execution
- Good Go bindings
- Widely used

**Cons**:
- Less familiar syntax
- 1-based indexing confusion
- Requires sandboxing setup
- Global state complexity
- Not configuration-focused

**Evaluation**: Rejected - requires more sandboxing work and less familiar syntax

### Option 2: JavaScript (V8/QuickJS)
**Description**: Embed JavaScript engine for extensions

**Pros**:
- Most widely known language
- Rich ecosystem
- Powerful features
- JSON native support
- Good debugging tools

**Cons**:
- Large binary size with V8
- Complex sandboxing
- Async complexity unnecessary
- Resource consumption issues
- Non-deterministic features

**Evaluation**: Rejected - too heavyweight and hard to sandbox effectively

### Option 3: Starlark (Selected)
**Description**: Embed Starlark configuration language

**Pros**:
- Python-like familiar syntax
- Sandboxed by design
- Deterministic execution
- Native Go implementation
- Configuration-focused
- Battle-tested in Bazel

**Cons**:
- Limited standard library - mitigated by custom functions
- No classes - not needed for validators
- Smaller community - offset by simplicity

**Evaluation**: Selected for security, simplicity, and configuration focus

## Consequences

### Positive Consequences
- **Security by Default**: Sandboxed execution prevents issues
- **Familiar Syntax**: Python developers adapt quickly
- **Deterministic**: Same inputs always produce same outputs
- **Fast Execution**: Compiled bytecode runs efficiently
- **Easy Integration**: Clean Go API for embedding
- **Reusable Logic**: Can share validator libraries

### Negative Consequences
- **Limited Features**: No classes or full Python stdlib
- **Learning Curve**: Subtle differences from Python
- **Debugging Tools**: Limited compared to full languages
- **Documentation**: Less than mainstream languages

### Neutral Consequences
- **Python Subset**: Familiar but not identical
- **Configuration Language**: Not general purpose
- **Immutable Data**: Different from Python mutability

## Implementation Impact

### Development Impact
- **Effort**: Medium - Integrate interpreter and build stdlib
- **Time**: 2-3 weeks for full integration
- **Skills Required**: Go embedding, Starlark knowledge

### Operational Impact
- **Performance**: < 10ms typical validator execution
- **Scalability**: Excellent - stateless execution
- **Maintenance**: Low - stable language specification

### Security Impact
- Complete sandboxing by default
- No filesystem or network access
- Resource limits enforceable
- Static analysis possible
- No code injection risks

## Risks and Mitigation

| Risk | Probability | Impact | Mitigation Strategy |
|------|------------|--------|-------------------|
| User confusion with Python | Medium | Low | Clear documentation of differences |
| Performance issues | Low | Medium | Caching and optimization |
| Limited debugging | Medium | Low | Better error messages and tools |
| Feature requests | High | Low | Clear scope boundaries |

## Dependencies

### Technical Dependencies
- Starlark-go library
- Go 1.21+ for embedding
- Resource limit enforcement

### Decision Dependencies
- ADR-003: Go implementation enables embedding
- ADR-001: Validators for workflow structure
- ADR-008: Community validators sharing

## Validation

### How We'll Know This Was Right
- Validator execution < 10ms average
- Zero security incidents from validators
- 50%+ projects using custom validators
- High user satisfaction ratings
- Successful validator sharing

### Review Triggers
This decision should be reviewed if:
- Security vulnerability discovered
- Performance degrades significantly
- Alternative technology breakthrough
- User adoption is low
- Maintenance burden increases

## References

### Internal References
- [DDX Validator Guide](/workflows/validators/README.md)
- [DDX Security Model](/docs/security/README.md)
- Related ADRs: ADR-001, ADR-003, ADR-008

### External References
- [Starlark Language Spec](https://github.com/bazelbuild/starlark/blob/master/spec.md)
- [Starlark in Go](https://github.com/google/starlark-go)
- [Bazel Starlark Docs](https://bazel.build/rules/language)
- [Configuration Languages Comparison](https://github.com/dhall-lang/dhall-lang/wiki/Comparisons)

## Notes

### Meeting Notes
- Security was primary concern in selection
- Team valued Python-like syntax for adoption
- Determinism important for CI/CD integration

### Future Considerations
- Build standard library of validators
- Create validator marketplace
- Add debugging tools
- Consider visual validator builder

### Lessons Learned
*To be filled after 6 months of production use*

---

## Decision History

### 2025-01-12 - Initial Decision
- Status: Proposed
- Author: DDX Development Team
- Notes: Evaluation of embedding options

### 2025-01-12 - Review and Acceptance
- Status: Accepted
- Reviewers: DDX Core Team, Security Team
- Changes: Added resource limits requirement

### Post-Implementation Review
- *To be scheduled after Q4 2025*

---
*This ADR documents a significant architectural decision and its rationale for future reference.*