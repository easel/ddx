---
tags: [adr, architecture, development-workflow, decisions, ddx, starlark, validators, extensions]
template: false
version: 1.0.0
---

# ADR-004: Embed Starlark for Workflow Extensions and Validators

**Date**: 2025-01-12  
**Status**: Accepted  
**Deciders**: DDX Development Team  
**Technical Story**: Enable user-customizable workflow extensions, validators, and automation through embedded scripting

## Context

### Problem Statement
DDX needs to provide extensible workflow automation that allows users to:
- Define custom validation rules for their development processes
- Enforce project-specific policies and conventions
- Create workflow extensions without modifying DDX core
- Share reusable validation logic across projects
- Implement complex conditional logic in workflows
- Maintain security by sandboxing user-provided code

### Forces at Play
- **Extensibility**: Users need to customize workflows for their specific needs
- **Security**: User-provided code must not compromise system security
- **Performance**: Validators must execute quickly in pre-commit hooks
- **Simplicity**: Extension language should be easy to learn and use
- **Portability**: Extensions should work across all platforms
- **Determinism**: Same inputs should always produce same outputs
- **Integration**: Must integrate cleanly with Go CLI implementation

### Constraints
- Must be embeddable in Go binary
- Cannot allow filesystem or network access without explicit permissions
- Must have predictable resource consumption
- Should use familiar Python-like syntax for developer adoption
- Need to support configuration-as-code patterns
- Must be statically analyzable for security auditing

## Decision

### Chosen Approach
Embed Starlark (a Python-like configuration language) for workflow extensions, validators, and custom automation within DDX.

### Rationale
- **Python-like Syntax**: Familiar to most developers, reducing learning curve
- **Sandboxed by Design**: No filesystem, network, or system access by default
- **Deterministic**: No global state, guaranteed reproducible execution
- **Go Integration**: Native Go implementation with excellent embedding support
- **Battle-tested**: Used by Bazel, Buck, and other build systems
- **Performance**: Faster than embedded interpreters like Lua or JavaScript
- **Configuration Focus**: Designed specifically for configuration and rules
- **Type Safety**: Static type checking capabilities available

## Alternatives Considered

### Option 1: Lua
**Description**: Embed Lua scripting language for extensions

**Pros**:
- Mature and widely used embedding language
- Small footprint
- Fast execution
- Good Go bindings available (gopher-lua)
- Used successfully in many projects (nginx, Redis)

**Cons**:
- Less familiar syntax for many developers
- 1-based indexing confuses users
- Requires careful sandboxing setup
- Global state management complexity
- Less suited for configuration use cases

**Why rejected**: Requires more sandboxing work and less familiar to Python developers who form a large part of the user base

### Option 2: JavaScript (V8/QuickJS)
**Description**: Embed JavaScript engine for extensions

**Pros**:
- Most widely known programming language
- Rich ecosystem and tooling
- Powerful language features
- JSON native support
- Good debugging tools

**Cons**:
- Large binary size with V8 embedding
- Complex sandboxing requirements
- Async complexity not needed for validators
- Resource consumption harder to control
- Non-deterministic features (Date, Math.random)

**Why rejected**: Too heavyweight for simple validation logic and harder to sandbox effectively

### Option 3: WebAssembly (WASM)
**Description**: Use WASM as the extension runtime

**Pros**:
- Language agnostic (compile from any language)
- Strong sandboxing guarantees
- Good performance
- Growing ecosystem
- Deterministic execution

**Cons**:
- Complex toolchain for users
- Debugging is difficult
- Large runtime overhead for simple scripts
- Poor developer experience for configuration
- Overkill for validation logic

**Why rejected**: Too complex for users to write simple validators and poor developer experience

### Option 4: YAML/JSON with JSONPath
**Description**: Pure declarative configuration with JSONPath expressions

**Pros**:
- No code execution, purely declarative
- Simple to understand
- Easy to validate and audit
- No security concerns
- Version control friendly

**Cons**:
- Very limited logic capabilities
- Cannot express complex validation rules
- No loops or conditionals
- Quickly becomes unmaintainable for complex rules
- Cannot share logic between validators

**Why rejected**: Insufficient expressiveness for complex validation scenarios

### Option 5: Native Go Plugins
**Description**: Compile validators as Go plugins loaded at runtime

**Pros**:
- Native performance
- Full Go capabilities
- Type safe
- No additional runtime needed

**Cons**:
- Platform-specific (poor Windows support)
- Requires Go toolchain for users
- Version compatibility issues
- No sandboxing - full system access
- Complex distribution model
- Security concerns with arbitrary code

**Why rejected**: Platform limitations and security concerns make this unsuitable

### Option 6: CEL (Common Expression Language)
**Description**: Use Google's CEL for expression evaluation

**Pros**:
- Designed for security policies
- Non-Turing complete (guaranteed termination)
- Good Go support
- Used in Kubernetes policies

**Cons**:
- Limited to expressions, not full programs
- Less familiar syntax
- Limited control flow
- Cannot define functions
- Harder to compose complex validators

**Why rejected**: Too limited for complex validation logic requiring multiple steps

## Consequences

### Positive Consequences
- **Familiar Syntax**: Python developers can write extensions immediately
- **Security by Default**: Sandboxed execution prevents security issues
- **Reproducible**: Deterministic execution ensures consistent behavior
- **Fast Execution**: Compiled bytecode runs efficiently
- **Easy Integration**: Clean Go API for embedding
- **Reusable Logic**: Can share validator libraries across projects
- **Static Analysis**: Can analyze scripts for security before execution
- **Good Tooling**: Starlark LSP and formatters available

### Negative Consequences
- **Limited Features**: No classes, limited standard library
- **Learning Curve**: Subtle differences from Python can confuse
- **Debugging**: Limited debugging capabilities compared to full languages
- **No Package Manager**: Must implement our own module system
- **Documentation**: Less documentation than mainstream languages
- **Community Size**: Smaller community than Lua or JavaScript

### Neutral Consequences
- **Python Subset**: Familiar but not identical to Python
- **Configuration Language**: Optimized for configuration, not general programming
- **Bazel Association**: Some developers associate with Bazel complexity
- **Frozen Types**: Immutable data structures by default
- **Import System**: Different from Python's import system

## Implementation

### Required Changes
- Integrate Starlark interpreter into DDX CLI
- Create standard library for validator common operations
- Implement module loading system for sharing code
- Build context object marshaling from Go to Starlark
- Create debugging and testing utilities
- Document Starlark dialect and available functions
- Implement resource limits (execution time, memory)

### Migration Strategy
For existing projects:
1. Start with built-in validators (no Starlark required)
2. Gradually introduce custom validators as needed
3. Share successful patterns as reusable modules
4. Build library of common validation patterns

### Success Metrics
- **Execution Speed**: < 10ms for typical validator
- **Memory Usage**: < 10MB per validator execution
- **Adoption Rate**: > 50% of projects using custom validators
- **Security Incidents**: Zero security issues from validators
- **User Satisfaction**: > 80% find Starlark easy to use
- **Reuse Rate**: Average 3+ shared validators per project

## Compliance

### Security Requirements
- Complete sandboxing with no filesystem/network access
- Resource limits enforced (CPU time, memory)
- Static analysis before execution
- No code injection vulnerabilities
- Audit logging of validator execution

### Performance Requirements
- Validator execution < 10ms average
- Pre-commit hook total time < 500ms
- Memory usage < 10MB per validator
- Support parallel validator execution

### Regulatory Requirements
- No execution of untrusted code
- Clear security boundaries
- Audit trail for policy enforcement
- GDPR compliance for any data processing

## Monitoring and Review

### Key Indicators to Watch
- Validator execution time percentiles
- Memory usage patterns
- Security vulnerability reports
- User-reported confusion points
- Adoption and usage statistics
- Common error patterns

### Review Date
Q4 2025 - After 12 months of production use

### Review Triggers
- Security vulnerability discovered
- Performance regression > 50%
- Major Starlark version release
- User satisfaction drops below 70%
- Alternative technology breakthrough

## Related Decisions

### Dependencies
- ADR-002: Go CLI Implementation - Starlark must embed in Go
- ADR-001: Git Subtree - Validators distributed via subtree
- ADR-003: YAML Configuration - Validator config in YAML

### Influenced By
- Need for CDP (Continuous Documentation Process) validation
- User requests for customizable workflows
- Security requirements for multi-tenant usage
- Success of Starlark in Bazel/Buck

### Influences
- Validator library design
- Module distribution system
- Debugging tool requirements
- Documentation approach
- Community contribution model

## References

### Documentation
- [Starlark Language Specification](https://github.com/bazelbuild/starlark/blob/master/spec.md)
- [Starlark in Go](https://github.com/google/starlark-go)
- [DDX Validator Guide](/workflows/cdp/validators/README.md)
- [Bazel's Starlark Documentation](https://bazel.build/rules/language)

### External Resources
- [Why Starlark?](https://github.com/bazelbuild/starlark/blob/master/users.md)
- [Starlark vs Other Languages](https://github.com/google/starlark-go/blob/master/doc/README.md)
- [Configuration Languages Comparison](https://github.com/dhall-lang/dhall-lang/wiki/Comparisons)

### Discussion History
- Initial discussion on extensibility requirements
- Security review of embedding options
- Performance benchmarks of different solutions
- User survey on preferred syntax

## Notes

The choice of Starlark aligns well with DDX's medical metaphor - like medical diagnostic criteria that must be precise, reproducible, and auditable, Starlark validators provide deterministic, sandboxed execution of validation rules. The Python-like syntax makes it accessible to medical researchers and data scientists who often use Python.

Key insight: By choosing a configuration language rather than a general-purpose scripting language, we guide users toward writing simple, focused validators rather than complex programs. This constraint actually improves maintainability and sharing.

The sandboxed nature of Starlark means validators can be safely shared and executed without security review, enabling a marketplace of reusable validation patterns - similar to sharing medical diagnostic protocols across institutions.

Implementation tip: Start with a minimal standard library and grow it based on user needs. This prevents bloat and ensures every function has a real use case.

---

**Last Updated**: 2025-01-12  
**Next Review**: 2025-12-01