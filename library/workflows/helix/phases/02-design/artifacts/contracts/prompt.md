# Contract Definition Prompt

Define precise API contracts that will be tested before implementation.

## Storage Location

Store API contracts at: `docs/helix/02-design/contracts/API-{number}-{title-with-hyphens}.md`

## Naming Convention

Follow this consistent format:
- **File Format**: `API-{number}-{title-with-hyphens}.md`
- **Number**: Zero-padded 3-digit sequence (001, 002, 003...)
- **Title**: Descriptive, lowercase with hyphens
- **Examples**:
  - `API-001-cli-command-interface.md`
  - `API-002-template-management-api.md`
  - `API-003-ai-integration-api.md`

## Contract-First Principles

### 1. Define the Interface, Not the Implementation
- Specify WHAT the API does, not HOW
- Focus on inputs, outputs, and errors
- Leave implementation details undefined

### 2. Make Contracts Testable
Every contract must be verifiable:
- Clear input specifications
- Exact output formats
- Specific error conditions
- Measurable performance requirements

### 3. Design for Users, Not Implementers
- Make the API intuitive to use
- Provide clear, helpful error messages
- Include comprehensive examples
- Document all edge cases

## Required Contract Elements

### For CLI Interfaces
1. **Command syntax** with all options
2. **Input formats** (stdin, files, arguments)
3. **Output formats** (stdout, stderr, files)
4. **Exit codes** with meanings
5. **Real examples** with expected output

### For REST APIs
1. **Endpoint paths** with methods
2. **Request schemas** with validation rules
3. **Response schemas** for success and error
4. **Status codes** with meanings
5. **Authentication** requirements

### For Library APIs
1. **Function signatures** with types
2. **Parameter constraints** and validation
3. **Return values** with possible states
4. **Exceptions** with conditions
5. **Thread safety** guarantees

## Contract Quality Checklist

- [ ] Can we write tests from this contract alone?
- [ ] Are all edge cases documented?
- [ ] Are error messages helpful?
- [ ] Are examples complete and correct?
- [ ] Is versioning strategy clear?

## Anti-Patterns to Avoid

❌ **Vague Contracts**
```
"Returns data about the user"
```

✅ **Precise Contracts**
```json
{
  "userId": "string (UUID)",
  "name": "string (max 100 chars)",
  "email": "string (valid email format)",
  "createdAt": "string (ISO-8601)"
}
```

❌ **Implementation Leaking**
```
"Queries the PostgreSQL database for users"
```

✅ **Implementation Agnostic**
```
"Returns users matching the search criteria"
```

❌ **Untestable Requirements**
```
"Should be fast"
```

✅ **Testable Requirements**
```
"Response time < 100ms for requests with < 100 items"
```

## Example: Good CLI Contract

```bash
Command: parse
Purpose: Parse and validate JSON input
Usage: $ tool parse [--strict] [--schema FILE]

Input:
  Format: JSON via stdin or --file
  Max size: 10MB

Output:
  Success: Validated JSON to stdout
  Error: Error message to stderr

Exit Codes:
  0: Valid JSON
  1: Invalid JSON
  2: Schema validation failed
  3: File not found

Example:
  $ echo '{"key": "value"}' | tool parse
  {"key": "value"}
  
  $ echo 'invalid' | tool parse
  Error: Invalid JSON at position 0
  $ echo $?
  1
```

This contract is complete, testable, and user-friendly.