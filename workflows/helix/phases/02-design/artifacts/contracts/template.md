# API Contract: [Contract Name] [FEAT-XXX]

**Contract ID**: API-XXX
**Feature**: FEAT-XXX
**Type**: [REST | GraphQL | CLI | Library]
**Status**: [Draft | Reviewed | Approved]
**Version**: 1.0.0

*Define all external interfaces before implementation*

## CLI Interface Contract

### Command Structure
```bash
$ [program] [command] [options] [arguments]
```

### Commands

#### Command: [command-name]
**Purpose**: [What this command does]
**Usage**: `$ [program] [command] [options]`

**Options**:
- `--option1, -o` : [Description]
- `--option2` : [Description]

**Input**:
- Format: [JSON|Text|File]
- Schema: [Define structure]

**Output**:
- Format: [JSON|Text]
- Schema: [Define structure]

**Exit Codes**:
- `0`: Success
- `1`: General error
- `2`: [Specific error]

**Examples**:
```bash
# Example 1
$ [program] [command] --option value
[Expected output]

# Example 2
$ echo "input" | [program] [command]
[Expected output]
```

---

## REST API Contract (if applicable)

### Base URL
```
[protocol]://[host]:[port]/api/v1
```

### Endpoints

#### GET /[resource]
**Purpose**: [What this endpoint does]
**Authentication**: [Required|Optional|None]

**Request**:
```http
GET /[resource]?param1=value HTTP/1.1
Host: [host]
Accept: application/json
Authorization: Bearer [token]
```

**Response Success (200)**:
```json
{
  "field1": "value",
  "field2": 123
}
```

**Response Error (4xx/5xx)**:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

---

## Library API Contract

### Public Functions

#### Function: `functionName`
```[language]
function functionName(param1: Type, param2: Type): ReturnType
```

**Purpose**: [What this function does]
**Parameters**:
- `param1`: [Description, constraints]
- `param2`: [Description, constraints]

**Returns**: [Description of return value]

**Throws**: 
- `ErrorType1`: When [condition]
- `ErrorType2`: When [condition]

**Example**:
```[language]
const result = functionName(value1, value2);
// result contains...
```

---

## Data Contracts

### Input Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "field1": {
      "type": "string",
      "description": "..."
    }
  },
  "required": ["field1"]
}
```

### Output Schema
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "result": {
      "type": "string"
    }
  }
}
```

---

## Error Contracts

### Error Codes
| Code | Description | HTTP Status | Recovery Action |
|------|-------------|-------------|-----------------|
| ERR_001 | [Description] | 400 | [What user should do] |
| ERR_002 | [Description] | 500 | [What user should do] |

### Error Response Format
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": {},
    "timestamp": "ISO-8601"
  }
}
```

---

## Contract Validation

### Test Scenarios
1. **Happy Path**: [Expected behavior with valid input]
2. **Invalid Input**: [How system handles bad input]
3. **Edge Cases**: [Boundary conditions]
4. **Error Cases**: [How errors are reported]

### Backwards Compatibility
- [ ] All changes are additive only
- [ ] No breaking changes to existing contracts
- [ ] Version negotiation supported

---

## Feature Traceability

### Parent Feature
- **Feature Specification**: `docs/features/FEAT-XXX/specification.md`
- **User Stories Implemented**: US-XXX, US-XXX

### Related Artifacts
- **ADRs**: [Related architectural decisions]
- **Test Suites**: `tests/FEAT-XXX/contract/`
- **Implementation**: `src/features/FEAT-XXX/api/`

### Contract Naming Convention
- Format: `[feature]-[interface-type]-contract.md`
- Example: `auth-rest-api-contract.md`
- Example: `payment-cli-contract.md`

---
*Note: Create one contract document per major interface.*
*Some contracts may serve multiple features (mark as "Cross-cutting").*
*Contract ID (API-XXX) should be unique across the project.*