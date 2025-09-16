# CLI-004: Library Management Contract

## Overview
This contract defines how DDx CLI commands interact with the library management system.

## Configuration Methods

### 1. Command-Line Flag
```bash
ddx --library-base-path /custom/path [command]
```
- **Scope**: Applies to single command execution
- **Priority**: Highest (overrides all other methods)
- **Use Case**: Testing, development, CI/CD

### 2. Environment Variable
```bash
export DDX_LIBRARY_BASE_PATH=/custom/library
ddx [command]
```
- **Scope**: Applies to shell session
- **Priority**: Second highest
- **Use Case**: Testing environments, containers

### 3. Automatic Detection
The system automatically detects the appropriate library:
- **Development**: `<git-repo>/library/` when in DDx repository
- **Project**: `.ddx/library/` in project root or parent directories
- **Global**: `~/.ddx/library/` as fallback

## Library Structure

```
library/
├── personas/           # AI personality definitions
│   └── *.md           # Persona files with YAML frontmatter
├── mcp-servers/       # MCP server configurations
│   ├── registry.yml   # Server registry
│   └── servers/       # Server definitions
│       └── *.yml
├── templates/         # Project templates
│   └── */            # Template directories
├── patterns/         # Reusable code patterns
│   └── */           # Pattern directories
├── prompts/         # AI prompts
│   └── */          # Prompt categories
└── configs/        # Configuration templates
    └── *.yml      # Config files
```

## API Contract

### Core Functions

#### GetLibraryPath
```go
func GetLibraryPath(overridePath string) (string, error)
```
**Returns**: Absolute path to library directory
**Priority Order**:
1. overridePath parameter
2. DDX_LIBRARY_BASE_PATH environment variable
3. Git repository library/ (if in DDx repo)
4. Nearest .ddx/library/
5. ~/.ddx/library/

#### Resource-Specific Helpers
```go
func GetPersonasPath(libraryOverride string) (string, error)
func GetMCPServersPath(libraryOverride string) (string, error)
func GetTemplatesPath(libraryOverride string) (string, error)
func GetPatternsPath(libraryOverride string) (string, error)
func GetPromptsPath(libraryOverride string) (string, error)
func GetConfigsPath(libraryOverride string) (string, error)
```

#### Resource Resolution
```go
func ResolveLibraryResource(resourcePath string, libraryOverride string) (string, error)
```
**Purpose**: Resolve a resource path relative to the library
**Security**: Prevents directory traversal attacks
**Validation**: Ensures resource exists

## Command Integration

### Commands Using Library

| Command | Library Usage |
|---------|--------------|
| `ddx persona list` | Lists from library/personas/ |
| `ddx persona load` | Loads from library/personas/ |
| `ddx mcp list` | Uses library/mcp-servers/registry.yml |
| `ddx mcp install` | Reads from library/mcp-servers/ |
| `ddx init --template` | Copies from library/templates/ |
| `ddx apply` | Applies from library/patterns/ |
| `ddx list` | Shows all library resources |

### Error Handling

| Error | Message | Resolution |
|-------|---------|-----------|
| Library not found | "Library path does not exist: [path]" | Create directory or fix path |
| Resource not found | "Resource not found: [resource]" | Check resource exists in library |
| Invalid path | "Invalid resource path: [path]" | Fix path traversal attempt |
| Permission denied | "Cannot access library: [path]" | Fix file permissions |

## Migration Contract

### From Old Structure
```bash
# Old structure (repository root)
personas/
templates/
patterns/
prompts/
mcp-servers/

# New structure (under library/)
library/
├── personas/
├── templates/
├── patterns/
├── prompts/
└── mcp-servers/
```

### Backward Compatibility
- Check old locations if library/ not found (deprecated)
- Log warnings when using old structure
- Provide migration command: `ddx migrate-library`

## Testing Contract

### Test Scenarios
1. **Development Mode**: Verify uses repo library/
2. **Override Flag**: Test --library-base-path works
3. **Environment Variable**: Verify DDX_LIBRARY_BASE_PATH
4. **Project Library**: Test .ddx/library/ discovery
5. **Global Fallback**: Verify ~/.ddx/library/ usage

### Test Utilities
```go
// Create test library structure
func CreateTestLibrary(t *testing.T, path string)

// Set library override for tests
func WithLibraryPath(path string) func()
```

## Security Contract

### Path Validation
- No directory traversal (../)
- Absolute paths resolved
- Symbolic links followed safely
- Permissions checked

### Sensitive Data
- No credentials in library files
- Environment variables for secrets
- Secure file permissions (0644 max)

## Performance Contract

### Benchmarks
- Path resolution: < 10ms
- Resource loading: < 100ms
- Library enumeration: < 500ms

### Caching
- Resolved paths cached per session
- Library structure cached with TTL
- Invalidation on filesystem changes

## Versioning Contract

### Library Version
- Version file: library/VERSION
- Semantic versioning (MAJOR.MINOR.PATCH)
- Compatibility checks on load

### Resource Versions
- Individual resources can have versions
- Specified in metadata/frontmatter
- Version constraints in dependencies