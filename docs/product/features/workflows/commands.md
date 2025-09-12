# Workflow Commands Specification

**Version**: 1.0  
**Status**: Specification  
**Date**: 2025-09-12  
**Related**: [Technical Overview](./overview.md), [Usage Guide](../../../usage/workflows/overview.md)

## Overview

This document specifies the complete command-line interface for DDX workflow operations. All workflow commands are implemented as subcommands under the `ddx workflow` namespace, providing a consistent and discoverable interface.

## Command Structure

### Primary Commands

#### `ddx workflow list`
Lists available workflows with filtering and search capabilities.

**Syntax**:
```bash
ddx workflow list [options]
```

**Options**:
- `--category <category>` - Filter by workflow category
- `--tag <tag>` - Filter by workflow tags (can be repeated)
- `--author <author>` - Filter by workflow author
- `--format <format>` - Output format: `table`, `json`, `yaml` (default: `table`)
- `--verbose, -v` - Show detailed workflow information
- `--search <term>` - Search workflow names and descriptions

**Output Format** (table):
```
NAME              CATEGORY        VERSION   AUTHOR      DESCRIPTION
development       software-dev    1.2.0     ddx-team    Complete SDLC workflow
incident-response ops             1.0.1     ops-team    Production incident handling
product-launch    product         1.1.0     pm-team     Product launch process
```

**Exit Codes**:
- `0` - Success
- `1` - No workflows found matching criteria
- `2` - Configuration error

---

#### `ddx workflow create`
Creates a new workflow interactively or from specification.

**Syntax**:
```bash
ddx workflow create [name] [options]
```

**Arguments**:
- `name` - Workflow name (optional, prompted if not provided)

**Options**:
- `--template <template>` - Base workflow template to extend
- `--category <category>` - Workflow category
- `--description <desc>` - Short description
- `--author <author>` - Author name (defaults to git config)
- `--interactive, -i` - Interactive creation mode (default)
- `--batch` - Non-interactive mode (requires all parameters)
- `--output-dir <path>` - Custom output directory
- `--dry-run` - Preview creation without writing files

**Interactive Flow**:
1. Workflow name (if not provided)
2. Category selection from available categories
3. Description entry
4. Phase definition (iterative)
5. Artifact identification
6. Template and prompt creation options
7. Review and confirmation

**Exit Codes**:
- `0` - Success
- `1` - Creation failed (validation errors, file conflicts)
- `2` - User cancelled interactive session
- `3` - Invalid parameters

---

#### `ddx workflow apply`
Applies an existing workflow to the current project.

**Syntax**:
```bash
ddx workflow apply <name>[:<artifact>] [options]
```

**Arguments**:
- `name` - Workflow name (required)
- `artifact` - Specific artifact to apply (optional)

**Options**:
- `--phase <phase>` - Start from specific phase
- `--variables <file>` - Variable values file (YAML/JSON)
- `--set <key=value>` - Set individual variables (can be repeated)
- `--output-dir <path>` - Custom output directory (default: `.ddx/workflows/<name>`)
- `--force` - Overwrite existing files
- `--interactive, -i` - Interactive variable collection
- `--validate-only` - Validate without executing
- `--resume` - Resume interrupted workflow
- `--skip-phase <phase>` - Skip specific phases (can be repeated)

**Execution Flow**:
1. Workflow validation and loading
2. Variable collection (interactive or from file)
3. Phase dependency resolution
4. Sequential/parallel phase execution
5. Artifact generation and validation
6. State persistence and reporting

**Exit Codes**:
- `0` - Success (all phases completed)
- `1` - Workflow execution failed
- `2` - Validation errors
- `3` - User cancelled during interactive mode
- `4` - Resume failed (corrupted state)

---

#### `ddx workflow validate`
Validates workflow structure and metadata.

**Syntax**:
```bash
ddx workflow validate <name> [options]
```

**Arguments**:
- `name` - Workflow name to validate

**Options**:
- `--strict` - Strict validation mode (fails on warnings)
- `--format <format>` - Output format: `text`, `json`, `junit` (default: `text`)
- `--output <file>` - Write results to file
- `--check-examples` - Validate example artifacts
- `--check-links` - Validate internal links
- `--schema-version <version>` - Validate against specific schema version

**Validation Checks**:
1. **Structure Validation**:
   - Required files present
   - Directory structure compliance
   - File naming conventions

2. **Metadata Validation**:
   - workflow.yml schema compliance
   - Artifact meta.yml validation
   - Version format validation

3. **Content Validation**:
   - Template variable consistency
   - Prompt template references
   - Phase dependency resolution

4. **Link Validation**:
   - Internal cross-references
   - Template includes
   - Example references

**Output Format** (text):
```
Validating workflow: development
✓ Structure validation passed
✓ Metadata validation passed  
✓ Content validation passed
⚠ Link validation: 2 warnings
  - phases/02-design.md: Reference to missing template
  - artifacts/prd/examples/: Empty examples directory

Summary: PASSED with warnings
```

**Exit Codes**:
- `0` - Validation passed
- `1` - Validation failed (errors found)
- `2` - Validation passed with warnings (--strict mode fails)
- `3` - Workflow not found

---

### Phase Management Commands

#### `ddx workflow phase list`
Lists phases for a workflow with status information.

**Syntax**:
```bash
ddx workflow phase list <workflow> [options]
```

**Options**:
- `--status <status>` - Filter by phase status: `pending`, `active`, `completed`, `failed`
- `--format <format>` - Output format: `table`, `json`, `yaml`
- `--show-criteria` - Show entry/exit criteria

---

#### `ddx workflow phase start`
Starts execution of a specific workflow phase.

**Syntax**:
```bash
ddx workflow phase start <workflow> <phase> [options]
```

**Options**:
- `--force` - Force start even if entry criteria not met
- `--variables <file>` - Phase-specific variable values
- `--timeout <duration>` - Phase timeout (overrides workflow default)

---

#### `ddx workflow phase complete`
Marks a phase as completed after validation.

**Syntax**:
```bash
ddx workflow phase complete <workflow> <phase> [options]
```

**Options**:
- `--artifacts <paths>` - Artifact paths to validate (comma-separated)
- `--force` - Force completion without validation
- `--notes <text>` - Completion notes

---

### Artifact Management Commands

#### `ddx workflow artifact list`
Lists artifacts for a workflow with generation status.

**Syntax**:
```bash
ddx workflow artifact list <workflow> [options]
```

**Options**:
- `--phase <phase>` - Filter by phase
- `--status <status>` - Filter by status: `pending`, `generated`, `validated`, `failed`
- `--type <type>` - Filter by artifact type: `document`, `code`, `config`, `data`

---

#### `ddx workflow artifact generate`
Generates a specific workflow artifact.

**Syntax**:
```bash
ddx workflow artifact generate <workflow> <artifact> [options]
```

**Options**:
- `--template-only` - Generate from template without AI assistance
- `--prompt-context <file>` - Additional context for AI prompts
- `--output <path>` - Custom output path
- `--model <model>` - AI model to use for generation

---

### Status and Information Commands

#### `ddx workflow status`
Shows current status of workflow execution.

**Syntax**:
```bash
ddx workflow status <workflow> [options]
```

**Options**:
- `--detailed` - Show phase and artifact details
- `--format <format>` - Output format: `table`, `json`, `yaml`
- `--watch` - Continuously monitor status

**Output Format**:
```
Workflow: development
Status: In Progress
Started: 2025-09-12 14:30:00
Progress: 3/6 phases completed

PHASE           STATUS      PROGRESS    ARTIFACTS
define          completed   100%        3/3 generated
design          completed   100%        2/2 generated  
implement       active      60%         4/7 generated
test            pending     0%          0/3 generated
release         pending     0%          0/2 generated
iterate         pending     0%          0/1 generated
```

---

#### `ddx workflow info`
Shows detailed information about a workflow.

**Syntax**:
```bash
ddx workflow info <workflow> [options]
```

**Options**:
- `--show-examples` - Include example listings
- `--show-variables` - Show variable definitions
- `--show-integrations` - Show external integrations

---

### Utility Commands

#### `ddx workflow reset`
Resets workflow state to initial conditions.

**Syntax**:
```bash
ddx workflow reset <workflow> [options]
```

**Options**:
- `--phase <phase>` - Reset from specific phase
- `--preserve-artifacts` - Keep generated artifacts
- `--confirm` - Skip confirmation prompt

**Warning**: This operation is destructive and removes workflow progress.

---

#### `ddx workflow export`
Exports workflow as a shareable package.

**Syntax**:
```bash
ddx workflow export <workflow> [options]
```

**Options**:
- `--format <format>` - Export format: `tar`, `zip`, `git-bundle`
- `--output <file>` - Output file path
- `--include-state` - Include execution state
- `--include-artifacts` - Include generated artifacts

---

#### `ddx workflow import`
Imports workflow from external package.

**Syntax**:
```bash
ddx workflow import <package> [options]
```

**Arguments**:
- `package` - Package file path or URL

**Options**:
- `--name <name>` - Override workflow name
- `--force` - Overwrite existing workflow
- `--validate` - Validate after import

---

## Global Options

All workflow commands support these global options:

- `--config <file>` - Custom configuration file
- `--verbose, -v` - Verbose output
- `--quiet, -q` - Suppress non-essential output
- `--no-color` - Disable colored output
- `--help, -h` - Show command help

## Configuration

Workflow commands respect DDX configuration in `.ddx/config.yml`:

```yaml
workflows:
  default_author: "Team Name"
  default_category: "development"
  auto_validate: true
  parallel_phases: true
  artifact_timeout: 300s
  
  # AI integration settings
  ai:
    default_model: "claude-3-sonnet"
    max_context_size: 100000
    temperature: 0.7
    
  # Storage settings  
  storage:
    workflow_dir: "workflows/"
    state_dir: ".ddx/workflows/"
    artifact_dir: "artifacts/"
```

## Error Handling

### Common Error Scenarios

1. **Workflow Not Found**:
   ```
   Error: Workflow 'invalid-name' not found
   Available workflows: development, incident-response, product-launch
   Use 'ddx workflow list' to see all available workflows
   ```

2. **Phase Dependency Failure**:
   ```
   Error: Cannot start phase 'test' - dependency 'implement' not completed
   Required entry criteria not met:
   - Implementation must be complete
   - Code review must be approved
   ```

3. **Validation Errors**:
   ```
   Error: Workflow validation failed
   - Missing required file: artifacts/prd/template.md
   - Invalid phase reference in workflow.yml: 'nonexistent-phase'  
   - Template variable mismatch: {{project_name}} not defined
   ```

4. **State Corruption**:
   ```
   Error: Workflow state corrupted or incompatible
   State version: 1.0.0, expected: 1.2.0
   Run 'ddx workflow reset <workflow>' to reinitialize
   ```

### Recovery Procedures

- **Corrupted State**: Use `ddx workflow reset` with appropriate options
- **Missing Dependencies**: Use `ddx workflow validate` to identify issues  
- **Interrupted Execution**: Use `ddx workflow apply --resume`
- **Configuration Issues**: Check `.ddx/config.yml` and global settings

## Performance Considerations

### Optimization Features
- **Lazy Loading**: Workflows loaded on-demand
- **Parallel Execution**: Independent phases run concurrently
- **Incremental Updates**: Only changed artifacts regenerated
- **Template Caching**: Compiled templates cached for reuse

### Resource Limits
- **Memory**: 512MB default limit for workflow execution
- **Timeout**: 30 minutes default workflow timeout  
- **Concurrency**: 4 parallel phases maximum
- **File Size**: 10MB limit for individual artifacts

## Security Considerations

### Access Controls
- Workflow execution requires project permissions
- Template modification restricted to workflow authors
- External integrations require explicit authorization

### Data Protection
- Sensitive variables encrypted at rest
- Artifact access logged for audit trails
- External system credentials managed securely

## Integration Points

### Git Integration
- Workflow state tracked in git history
- Commit hooks for validation
- Branch-specific workflow contexts

### CI/CD Integration  
- Workflow status exported for pipeline consumption
- Automated workflow execution triggers
- Quality gate integration

### External Tools
- Issue tracker integration (Jira, GitHub Issues)
- Documentation platforms (Confluence, Notion)
- Communication tools (Slack, Teams)

## Related Documentation

- [Technical Overview](./overview.md) - System architecture and design
- [Create Command Specification](./create-command.md) - Detailed create command spec
- [Apply Command Specification](./apply-command.md) - Detailed apply command spec
- [Usage Guide](../../../usage/workflows/overview.md) - User-facing documentation