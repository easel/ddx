---
tags: [adr, architecture, templates, variables, substitution, ddx]
template: false
version: 1.0.0
---

# ADR-007: Variable Substitution and Template Processing

**Date**: 2025-01-14
**Status**: Proposed
**Deciders**: DDX Development Team
**Technical Story**: Define the variable substitution system for processing templates and patterns with project-specific values

## Context

### Problem Statement
DDX needs a powerful yet safe variable substitution system that:
- Processes templates with project-specific values
- Supports nested variables and complex data structures
- Provides conditional logic for dynamic content
- Handles missing variables gracefully
- Prevents injection attacks and unsafe operations
- Maintains readability in templates
- Supports multiple substitution passes for complex scenarios

### Forces at Play
- **Power**: Need expressive substitution capabilities
- **Safety**: Must prevent code injection and path traversal
- **Simplicity**: Syntax should be intuitive and readable
- **Performance**: Processing must be fast for large templates
- **Debugging**: Clear error messages for template issues
- **Compatibility**: Should work with existing template formats
- **Extensibility**: Support custom functions and filters

### Constraints
- Cannot execute arbitrary code
- Must be deterministic and reproducible
- Should not conflict with target language syntax
- Must handle binary files appropriately
- Need to preserve file permissions and attributes
- Must work with version control systems

## Decision

### Chosen Approach
Implement a multi-stage template processing system using:
1. **Go text/template** as the primary engine
2. **Custom delimiters** to avoid conflicts
3. **Staged processing** for complex substitutions
4. **Safe function library** with sandboxed operations
5. **Validation layer** to prevent dangerous operations

### Template Syntax

#### Basic Variables
```
{{.ProjectName}}              # Simple variable
{{.Author.Name}}             # Nested variable
{{.Config.Database.Host}}    # Deep nesting
{{.Items[0]}}               # Array access
{{.Map["key"]}}             # Map access
```

#### Custom Delimiters (to avoid conflicts)
```
$[[.ProjectName]]            # Alternative delimiter
«.ProjectName»              # Unicode delimiter option
<%=.ProjectName%>           # ERB-style delimiter option
```

#### Conditional Logic
```
{{if .EnableAuth}}
  // Authentication code here
{{end}}

{{if eq .Database "postgres"}}
  // PostgreSQL specific code
{{else if eq .Database "mysql"}}
  // MySQL specific code
{{else}}
  // Default database code
{{end}}

{{range .Features}}
  - {{.Name}}: {{.Description}}
{{end}}
```

#### Built-in Functions
```
# String manipulation
{{upper .ProjectName}}           # UPPERCASE
{{lower .ProjectName}}           # lowercase
{{title .ProjectName}}           # Title Case
{{snake .ProjectName}}           # snake_case
{{camel .ProjectName}}           # camelCase
{{pascal .ProjectName}}          # PascalCase
{{kebab .ProjectName}}           # kebab-case
{{plural .EntityName}}           # Pluralization
{{singular .EntityName}}         # Singularization

# Path manipulation
{{base .FilePath}}               # Base filename
{{dir .FilePath}}                # Directory path
{{ext .FilePath}}                # File extension
{{join .PathParts "/"}}          # Path joining

# Data manipulation
{{default "value" .MaybeEmpty}}  # Default values
{{quote .StringValue}}           # Quote string
{{indent 4 .TextBlock}}          # Indent text
{{trim .StringValue}}            # Trim whitespace
{{split .CSVString ","}}         # Split string

# Date/Time
{{now}}                          # Current timestamp
{{date "2006-01-02" now}}       # Format date
{{dateModify "-1h" now}}        # Modify date

# Environment
{{env "HOME"}}                   # Environment variable
{{os}}                           # Operating system
{{arch}}                         # Architecture

# Encoding
{{b64enc .Secret}}               # Base64 encode
{{b64dec .Encoded}}              # Base64 decode
{{sha256sum .Content}}           # SHA256 hash
{{uuid}}                         # Generate UUID
```

### Variable Sources

#### Priority Order (highest to lowest)
1. Command-line flags (`--var key=value`)
2. Environment variables (`DDX_VAR_KEY`)
3. Local overrides (`.ddx.local.yml`)
4. Project configuration (`.ddx.yml`)
5. User defaults (`~/.ddx/defaults.yml`)
6. Template defaults (`template.defaults.yml`)

#### Variable Definition Schema
```yaml
variables:
  # Simple variables
  projectName: "my-project"
  author: "John Doe"

  # Nested structures
  database:
    type: "postgres"
    host: "localhost"
    port: 5432

  # Arrays
  features:
    - name: "auth"
      enabled: true
    - name: "api"
      enabled: false

  # Computed variables
  computed:
    projectSlug: "{{kebab .projectName}}"
    dbUrl: "{{.database.type}}://{{.database.host}}:{{.database.port}}"

  # Conditional defaults
  defaults:
    port: "{{if eq .environment \"production\"}}80{{else}}3000{{end}}"
```

### Processing Pipeline

```
1. Load Variables
   ├── Merge from all sources
   ├── Validate types
   └── Resolve computed values

2. Parse Template
   ├── Detect delimiter style
   ├── Parse syntax tree
   └── Validate functions

3. First Pass (Variables)
   ├── Substitute simple variables
   ├── Evaluate conditionals
   └── Process loops

4. Second Pass (Functions)
   ├── Execute functions
   ├── Apply filters
   └── Format output

5. Validation
   ├── Check for unresolved variables
   ├── Validate output
   └── Security scan

6. Write Output
   ├── Create directories
   ├── Write files
   └── Set permissions
```

### Rationale
- **Go text/template**: Mature, well-tested, good security model
- **Multi-stage**: Handles complex dependencies between variables
- **Custom delimiters**: Avoids conflicts with target languages
- **Safe functions**: Curated set prevents dangerous operations
- **Validation**: Catches errors before writing files

## Alternatives Considered

### Option 1: Mustache Templates
**Description**: Use Mustache logic-less templates

**Pros**:
- Very simple syntax
- Language agnostic
- Wide implementation support
- No logic in templates

**Cons**:
- Too limited for complex scenarios
- No conditional logic
- No functions or filters
- Can't handle nested structures well

**Why rejected**: Insufficient expressiveness for complex project templates

### Option 2: Jinja2-style Templates
**Description**: Implement Jinja2-compatible templating

**Pros**:
- Popular in Python ecosystem
- Powerful and expressive
- Good filter system
- Familiar to many developers

**Cons**:
- Complex implementation in Go
- Python-centric design
- Heavyweight for our needs
- Potential security concerns

**Why rejected**: Implementation complexity and Python-specific design

### Option 3: Handlebars Templates
**Description**: Use Handlebars templating system

**Pros**:
- Good balance of power and simplicity
- Helper system for extensions
- Block helpers for complex logic
- Popular in JavaScript ecosystem

**Cons**:
- JavaScript-centric design
- Limited Go implementations
- Custom helper complexity
- Less familiar to Go developers

**Why rejected**: Limited Go support and JavaScript-centric design

### Option 4: Liquid Templates
**Description**: Use Shopify's Liquid templating

**Pros**:
- Safe by design
- Good for user-generated content
- Clear syntax
- Good documentation

**Cons**:
- Ruby-centric origin
- Limited Go implementations
- Less powerful than text/template
- Smaller ecosystem

**Why rejected**: Limited Go ecosystem and less powerful than native solution

### Option 5: Custom Template Language
**Description**: Build our own template language

**Pros**:
- Perfect fit for our needs
- Complete control
- No external dependencies
- Optimized for our use case

**Cons**:
- Massive development effort
- Security risks
- No ecosystem
- Documentation burden
- Maintenance overhead

**Why rejected**: Unjustified development effort and security risks

## Consequences

### Positive Consequences
- **Powerful**: Handles complex templating scenarios
- **Safe**: Sandboxed execution prevents security issues
- **Familiar**: Go developers know text/template
- **Extensible**: Easy to add new functions
- **Fast**: Native Go performance
- **Maintainable**: Standard library quality

### Negative Consequences
- **Learning Curve**: Template syntax requires learning
- **Debugging**: Template errors can be cryptic
- **Complexity**: Multi-stage processing adds complexity
- **Go-specific**: Syntax familiar mainly to Go developers
- **Limited IDE Support**: Less tooling than mainstream options

### Neutral Consequences
- **Delimiter Choice**: Must document delimiter options
- **Function Library**: Requires careful curation
- **Error Handling**: Template errors need translation
- **Performance**: Template caching needed for speed

## Implementation

### Required Changes
1. Build template processing engine
2. Implement variable merging logic
3. Create safe function library
4. Add delimiter detection/configuration
5. Build validation layer
6. Implement caching system
7. Create debugging tools
8. Write comprehensive documentation

### Security Measures
```go
// Prohibited operations
- No file system access beyond output directory
- No network operations
- No system command execution
- No arbitrary code evaluation
- Path traversal prevention
- Size limits on expansions

// Validation checks
- Template size limits (10MB)
- Variable depth limits (10 levels)
- Loop iteration limits (1000)
- Output size limits (100MB)
- Execution time limits (30s)
```

### Error Handling
```
# Clear error messages
Error: Undefined variable "ProjectName" at line 15, column 8
  Did you mean "projectName"? (case sensitive)

Error: Function "random" is not available
  Available functions: uuid, now, env

Error: Template loop exceeded maximum iterations (1000)
  Check for infinite loops in range statements

Error: Invalid delimiter style in template
  Detected: {{...}}, Expected: $[[...]]
```

### Success Metrics
- **Processing Speed**: < 100ms for typical template
- **Error Clarity**: 90% errors understood without docs
- **Function Coverage**: All common operations supported
- **Security**: Zero template injection vulnerabilities
- **Compatibility**: Works with 95% of existing templates

## Compliance

### Security Requirements
- Complete sandboxing of template execution
- Input validation and sanitization
- Output validation for path traversal
- Resource limits enforcement
- Audit logging of template processing

### Performance Requirements
- Template processing < 100ms average
- Memory usage < 50MB per template
- Support for templates up to 10MB
- Concurrent template processing
- Template compilation caching

### Regulatory Requirements
- No execution of user code
- Data sanitization for PII
- Secure defaults for all operations
- Clear security boundaries

## Monitoring and Review

### Key Indicators to Watch
- Template processing performance
- Error rates and types
- Function usage statistics
- Security incident reports
- User feedback on syntax
- Cache hit rates

### Review Date
Q2 2025 - After initial production usage

### Review Triggers
- Security vulnerability discovered
- Performance degradation > 50%
- User satisfaction < 70%
- Major Go template changes
- Alternative technology breakthrough

## Related Decisions

### Dependencies
- ADR-001: Template structure definition
- ADR-003: Go implementation enables text/template
- ADR-005: Configuration defines variables
- ADR-011: Security requirements for sandboxing

### Influenced By
- Go text/template capabilities
- Security requirements
- User feedback on template complexity
- Common templating patterns

### Influences
- Template creation guidelines
- Variable naming conventions
- Documentation requirements
- IDE plugin development

## References

### Documentation
- [Go text/template](https://golang.org/pkg/text/template/)
- [Template Security](https://www.owasp.org/index.php/Template_Injection)
- [Sprig Template Functions](https://masterminds.github.io/sprig/)
- [Template Best Practices](https://docs.docker.com/develop/dev-best-practices/)

### External Resources
- [Template Engine Comparison](https://en.wikipedia.org/wiki/Comparison_of_web_template_engines)
- [Safe Template Processing](https://pragmaticwebsecurity.com/articles/spasecurity/template-injection.html)
- [Go Template Patterns](https://www.practical-go-lessons.com/chap-36-templates)

### Discussion History
- Template syntax evaluation
- Security review findings
- Performance benchmarking results
- User feedback on complexity

## Notes

The template system follows DDX's medical metaphor - like medical forms that are filled out with patient-specific information, templates are filled with project-specific values. The multi-stage processing is like a medical checklist ensuring all required information is properly collected and validated.

Key insight: By using Go's text/template with custom functions, we get a battle-tested engine with good security properties while maintaining the flexibility needed for complex templates.

Implementation tip: Start with a minimal function set and expand based on actual user needs. Each function added increases the API surface and potential security risk.

The variable precedence system mirrors medical prescription protocols - general guidelines (defaults) are overridden by specific patient needs (project config) and further adjusted by immediate circumstances (command-line flags).

---

**Last Updated**: 2025-01-14
**Next Review**: 2025-04-14