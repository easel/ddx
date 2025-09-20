# ADR-013: Library Path Resolution Strategy

## Status
Accepted

## Context
DDx needs a flexible way to locate library resources (templates, patterns, prompts, personas, MCP servers) that works across different scenarios:
- Development (working on DDx itself)
- Testing (with custom library paths)
- Production (installed DDx)
- Project-specific (customized libraries)

Previously, resources were scattered across the repository root, making it unclear what was library content vs. implementation code.

## Decision
We will implement a hierarchical library path resolution system with the following priority order:

1. **Command-line flag** (`--library-base-path`): Highest priority for explicit overrides
2. **Environment variable** (`DDX_LIBRARY_BASE_PATH`): For testing and CI/CD scenarios
3. **Git repository detection**: When developing DDx, use `<git-root>/library/`
4. **Project-local library**: Look for nearest `.ddx/library/` (traverse upward)
5. **Global library**: Fallback to `~/.ddx/library/`

All library content will be centralized under a single `library/` directory containing:
```
library/
├── personas/       # AI personality definitions
├── mcp-servers/    # MCP server registry and definitions
├── templates/      # Project templates
├── patterns/       # Code patterns
├── prompts/        # AI prompts
└── configs/        # Configuration templates
```

## Consequences

### Positive
- **Clear separation**: Library content is clearly separated from implementation code
- **Flexible testing**: Easy to test with different library configurations
- **Project isolation**: Projects can have custom libraries without affecting global setup
- **Development clarity**: When working on DDx itself, library is always at `<repo>/library/`
- **Consistent structure**: Same library layout everywhere (dev, installed, project)

### Negative
- **Migration required**: Existing installations need to move content to new structure
- **Path complexity**: Multiple resolution paths could be confusing initially
- **Documentation updates**: All references to old paths need updating

## Implementation
The resolution logic is implemented in `cli/internal/config/library.go` with:
- `GetLibraryPath()`: Main resolution function
- Helper functions for specific resource types (personas, templates, etc.)
- Integration with existing loaders and registries

## References
- Feature Specification: FEAT-012 (Library Management)
- Solution Design: SD-012-library-management.md
- Test Plan: test-plan-library.md