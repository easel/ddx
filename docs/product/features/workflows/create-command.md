# Workflow Create Command Specification

**Version**: 1.0  
**Status**: Specification  
**Date**: 2025-09-12  
**Related**: [Commands Overview](./commands.md), [Technical Overview](./overview.md)

## Overview

The `ddx workflow create` command provides comprehensive functionality for creating new workflows interactively or programmatically. It handles the complete workflow scaffolding process, from initial metadata collection to file generation and validation.

## Command Signature

```bash
ddx workflow create [name] [options]
```

### Arguments

- `name` (optional): Workflow identifier following naming conventions
  - Must be lowercase with hyphens (e.g., `my-custom-workflow`)
  - 3-50 characters in length
  - Cannot conflict with existing workflow names
  - If omitted, user will be prompted during interactive mode

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--template <template>` | string | none | Base workflow template to extend |
| `--category <category>` | string | prompted | Workflow category classification |
| `--description <desc>` | string | prompted | Short workflow description |
| `--author <author>` | string | git config | Author name for metadata |
| `--interactive, -i` | boolean | true | Interactive creation mode |
| `--batch` | boolean | false | Non-interactive mode requiring all parameters |
| `--output-dir <path>` | string | `workflows/` | Custom output directory |
| `--dry-run` | boolean | false | Preview creation without writing files |
| `--force` | boolean | false | Overwrite existing workflow |
| `--validate` | boolean | true | Validate created workflow |
| `--open-editor` | boolean | false | Open created files in editor |
| `--git-init` | boolean | false | Initialize git repo for workflow |

## Implementation Requirements

### Input Validation

#### Name Validation
```go
func ValidateWorkflowName(name string) error {
    if len(name) < 3 || len(name) > 50 {
        return errors.New("workflow name must be 3-50 characters")
    }
    
    if !regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`).MatchString(name) {
        return errors.New("workflow name must be lowercase with hyphens")
    }
    
    reservedNames := []string{"help", "version", "config", "list"}
    for _, reserved := range reservedNames {
        if name == reserved {
            return fmt.Errorf("'%s' is a reserved workflow name", name)
        }
    }
    
    return nil
}
```

#### Category Validation
```go
type WorkflowCategory struct {
    ID          string `yaml:"id"`
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Icon        string `yaml:"icon"`
}

var PredefinedCategories = []WorkflowCategory{
    {ID: "development", Name: "Software Development", Description: "SDLC workflows"},
    {ID: "operations", Name: "DevOps & Operations", Description: "Infrastructure and deployment"},
    {ID: "product", Name: "Product Management", Description: "Product lifecycle management"},
    {ID: "business", Name: "Business Process", Description: "General business workflows"},
    {ID: "research", Name: "Research & Analysis", Description: "Investigation and study workflows"},
    {ID: "custom", Name: "Custom", Description: "User-defined workflows"},
}
```

### Interactive Creation Flow

#### Step 1: Workflow Metadata Collection
```
Creating new DDX workflow...

Workflow Name: my-awesome-workflow
  ✓ Name available and valid

Category Selection:
  1. Software Development    - SDLC workflows
  2. DevOps & Operations    - Infrastructure and deployment  
  3. Product Management     - Product lifecycle management
  4. Business Process       - General business workflows
  5. Research & Analysis    - Investigation and study workflows
  6. Custom                 - User-defined workflows

Select category [1-6]: 1

Description: A comprehensive workflow for awesome development processes
  ✓ Description accepted (75 characters)

Author: Jane Doe <jane@example.com>
  ✓ Author information confirmed
```

#### Step 2: Phase Definition
```
Phase Definition:
Define the sequential phases of your workflow. Each phase should represent a major milestone.

Phase 1:
  Name: Planning
  Description: Define requirements and project scope
  ✓ Phase added

Phase 2: 
  Name: Implementation  
  Description: Build and develop the solution
  ✓ Phase added

Phase 3:
  Name: Validation
  Description: Test and validate the implementation
  ✓ Phase added

Add another phase? [y/N]: n
  ✓ 3 phases defined
```

#### Step 3: Artifact Identification
```
Artifact Definition:
Identify key deliverables for each phase.

Phase: Planning
  Artifact 1:
    Name: Requirements Document
    Type: [document/code/config/data]: document
    Required: [y/N]: y
    ✓ Artifact added
    
  Artifact 2:
    Name: Project Plan
    Type: document
    Required: y
    ✓ Artifact added
    
  Add another artifact for Planning? [y/N]: n

Phase: Implementation
  Artifact 1:
    Name: Source Code
    Type: code
    Required: y
    ✓ Artifact added
    
  Add another artifact for Implementation? [y/N]: n

Phase: Validation  
  Artifact 1:
    Name: Test Report
    Type: document
    Required: y
    ✓ Artifact added
    
  Add another artifact for Validation? [y/N]: n

✓ 4 artifacts defined across 3 phases
```

#### Step 4: Template and Prompt Options
```
Template and Prompt Generation:
Choose how to create templates and prompts for your artifacts.

For each artifact, select:
1. Generate basic template structure
2. Use existing template as base
3. Skip template generation (create manually)

Requirements Document:
  Template option [1-3]: 1
  Generate AI prompt? [Y/n]: y
  ✓ Will generate template and prompt

Project Plan:
  Template option [1-3]: 1  
  Generate AI prompt? [Y/n]: y
  ✓ Will generate template and prompt

Source Code:
  Template option [1-3]: 3
  (Code artifacts typically don't need templates)
  ✓ Will skip template generation

Test Report:
  Template option [1-3]: 1
  Generate AI prompt? [Y/n]: y  
  ✓ Will generate template and prompt
```

#### Step 5: Review and Confirmation
```
Workflow Summary:
──────────────────
Name: my-awesome-workflow
Category: Software Development
Author: Jane Doe <jane@example.com>
Description: A comprehensive workflow for awesome development processes

Phases:
  1. Planning (2 artifacts)
  2. Implementation (1 artifact)  
  3. Validation (1 artifact)

Artifacts:
  📄 Requirements Document (template + prompt)
  📄 Project Plan (template + prompt)
  💻 Source Code (no template)
  📄 Test Report (template + prompt)

Output Directory: workflows/my-awesome-workflow/

Create this workflow? [Y/n]: y
```

### File Generation Process

#### Directory Structure Creation
```go
type WorkflowStructure struct {
    RootDir   string
    Phases    []string
    Artifacts []ArtifactSpec
}

func (ws *WorkflowStructure) CreateDirectories() error {
    dirs := []string{
        ws.RootDir,
        filepath.Join(ws.RootDir, "phases"),
    }
    
    for _, artifact := range ws.Artifacts {
        artifactDir := filepath.Join(ws.RootDir, artifact.ID)
        dirs = append(dirs, artifactDir)
        dirs = append(dirs, filepath.Join(artifactDir, "examples"))
    }
    
    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
    }
    
    return nil
}
```

#### Metadata File Generation

##### workflow.yml Generation
```go
func GenerateWorkflowMetadata(spec WorkflowSpec) (*WorkflowMetadata, error) {
    metadata := &WorkflowMetadata{
        Name:        spec.Name,
        Version:     "1.0.0",
        Description: spec.Description,
        Author:      spec.Author,
        Category:    spec.Category,
        Tags:        spec.Tags,
        Created:     time.Now(),
        
        Phases:    make([]PhaseSpec, len(spec.Phases)),
        Artifacts: make([]ArtifactSpec, len(spec.Artifacts)),
        
        Automation: AutomationSpec{
            InitCommand:     fmt.Sprintf("ddx workflow init %s", spec.Name),
            ValidateCommand: fmt.Sprintf("ddx workflow validate %s", spec.Name),
            CompleteCommand: fmt.Sprintf("ddx workflow complete %s", spec.Name),
        },
    }
    
    // Configure phases
    for i, phase := range spec.Phases {
        phaseSpec := PhaseSpec{
            ID:          phase.ID,
            Name:        phase.Name,
            Description: phase.Description,
            Artifacts:   phase.ArtifactIDs,
            EntryConditions: []string{
                "Previous phase completed successfully",
                "Required artifacts validated",
            },
            ExitConditions: []string{
                "All phase artifacts generated",
                "Phase validation passed",
            },
        }
        
        // Link to next phase
        if i < len(spec.Phases)-1 {
            phaseSpec.Next = spec.Phases[i+1].ID
        }
        
        metadata.Phases[i] = phaseSpec
    }
    
    // Configure artifacts
    for i, artifact := range spec.Artifacts {
        artifactSpec := ArtifactSpec{
            ID:          artifact.ID,
            Name:        artifact.Name,
            Type:        artifact.Type,
            Required:    artifact.Required,
            Template:    fmt.Sprintf("%s/template.md", artifact.ID),
            Prompt:      fmt.Sprintf("%s/prompt.md", artifact.ID),
            Validation: ValidationSpec{
                RequiredSections: []string{"Title", "Content"},
                MinLength:       100,
                Format:          "markdown",
            },
        }
        
        metadata.Artifacts[i] = artifactSpec
    }
    
    return metadata, nil
}
```

#### Template Generation

##### Basic Document Template
```go
func GenerateDocumentTemplate(artifact ArtifactSpec) string {
    return fmt.Sprintf(`# %s

**Workflow**: {{workflow_name}}  
**Phase**: {{phase_name}}  
**Created**: {{date}}  
**Author**: {{author}}  
**Status**: {{status:draft}}

## Overview

Brief description of this %s and its purpose within the {{workflow_name}} workflow.

## Content

### Section 1: [Title]

[Content placeholder]

### Section 2: [Title]

[Content placeholder]

### Section 3: [Title]

[Content placeholder]

## Metadata

- **Dependencies**: {{dependencies}}
- **Related Artifacts**: {{related_artifacts}}
- **Version**: {{version:1.0.0}}
- **Last Modified**: {{last_modified}}

## Notes

{{notes}}

---
*Generated by DDX Workflow System*
`, artifact.Name, strings.ToLower(artifact.Name))
}
```

##### AI Prompt Generation
```go
func GenerateArtifactPrompt(artifact ArtifactSpec, workflow WorkflowSpec) string {
    return fmt.Sprintf(`# %s Creation Assistant

**Template**: [[template.md|%s Template]]  
**Workflow**: %s  
**Phase**: %s

## Purpose

This prompt assists in creating a comprehensive %s as part of the %s workflow. The %s serves as [purpose description].

## Information Gathering

To complete this artifact effectively, gather the following information:

### Context Questions
- What is the specific scope of this %s?
- Who is the target audience?
- What constraints or requirements must be considered?

### Content Questions  
- What are the key points that must be covered?
- What format and structure work best?
- What examples or references should be included?

### Quality Criteria
- How will success be measured?
- What validation steps are required?
- Who needs to review and approve?

## Template Structure

The template includes these main sections:
- **Overview**: Purpose and context
- **Content**: Main sections with placeholders
- **Metadata**: Tracking and reference information

## Best Practices

### Content Guidelines
- Be specific and actionable
- Use clear, concise language
- Include relevant examples
- Consider the audience's expertise level

### Common Pitfalls to Avoid
- Overly generic content
- Missing context or assumptions
- Incomplete metadata
- Poor organization or structure

## Template

{{template:template.md}}

## Instructions

1. Review the template structure above
2. Gather information using the questions provided
3. Fill in each template section with appropriate content
4. Validate against the quality criteria
5. Save the completed artifact

Remember to replace all placeholder text with actual content and ensure all metadata fields are properly completed.
`, artifact.Name, artifact.Name, workflow.Name, "[Phase Name]", 
   strings.ToLower(artifact.Name), workflow.Name, strings.ToLower(artifact.Name),
   strings.ToLower(artifact.Name))
}
```

#### Documentation Generation

##### README.md Generation
```go
func GenerateWorkflowReadme(spec WorkflowSpec) string {
    return fmt.Sprintf(`# %s Workflow

**Category**: %s  
**Version**: 1.0.0  
**Author**: %s  
**Created**: %s

## Overview

%s

This workflow follows DDX's medical metaphor, treating your project as a patient that needs diagnosis and treatment through structured phases and proven practices.

## When to Use This Workflow

Use this workflow when:
- [Condition 1 - when this workflow applies]
- [Condition 2 - specific scenarios]
- [Condition 3 - problem indicators]

## Workflow Phases

%s

## Quick Start

1. **Apply the workflow**:
   ` + "```bash\n   ddx workflow apply %s\n   ```" + `

2. **Follow the phases**:
   Work through each phase systematically, completing required artifacts

3. **Validate progress**:
   ` + "```bash\n   ddx workflow status %s\n   ```" + `

## Customization

This workflow can be customized by:
- Modifying templates in artifact directories
- Adjusting phase requirements in workflow.yml
- Adding custom validation rules
- Creating project-specific examples

## Contributing

Found improvements? Contribute back to the community:
` + "```bash\nddx contribute workflows/%s\n```" + `

## Related Workflows

- [Related Workflow 1]: Brief description
- [Related Workflow 2]: Brief description

## Medical Metaphor

In DDX's medical theme:
- **Workflow** = Treatment Protocol
- **Phases** = Treatment Steps  
- **Artifacts** = Medical Records
- **Templates** = Forms
- **Validation** = Quality Check

## Support

- [Usage Guide](../../../usage/workflows/overview.md)
- [Creating Workflows](../../../usage/workflows/creating-workflows.md)
- [Community Forum](https://github.com/ddx-dev/community)
`, spec.Name, spec.Category, spec.Author, time.Now().Format("2006-01-02"),
   spec.Description, generatePhaseList(spec.Phases), spec.Name, spec.Name, spec.Name)
}
```

### Batch Mode Implementation

#### Parameter Validation
```go
type BatchModeParams struct {
    Name        string   `json:"name" validate:"required,workflowname"`
    Category    string   `json:"category" validate:"required,oneof=development operations product business research custom"`
    Description string   `json:"description" validate:"required,min=10,max=200"`
    Author      string   `json:"author" validate:"required"`
    Phases      []string `json:"phases" validate:"required,min=1,max=10"`
    Artifacts   []ArtifactDef `json:"artifacts" validate:"required,min=1"`
}

func ValidateBatchParams(params BatchModeParams) error {
    validate := validator.New()
    validate.RegisterValidation("workflowname", validateWorkflowName)
    
    if err := validate.Struct(params); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Additional business logic validation
    if len(params.Phases) != len(unique(params.Phases)) {
        return errors.New("duplicate phase names not allowed")
    }
    
    return nil
}
```

#### Configuration File Support
```yaml
# workflow-config.yml
name: my-batch-workflow
category: development
description: Workflow created in batch mode
author: CI/CD System

phases:
  - id: phase1
    name: Planning
    description: Planning and requirements gathering
  - id: phase2  
    name: Implementation
    description: Development and coding
  - id: phase3
    name: Testing
    description: Quality assurance and validation

artifacts:
  - id: requirements
    name: Requirements Document
    type: document
    phase: phase1
    required: true
    template: true
    prompt: true
    
  - id: code
    name: Source Code
    type: code
    phase: phase2
    required: true
    template: false
    prompt: false
```

### Error Handling

#### Error Categories
```go
type CreateError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Code    int       `json:"code"`
    Details map[string]interface{} `json:"details,omitempty"`
}

type ErrorType string

const (
    ValidationError  ErrorType = "validation"
    FileSystemError  ErrorType = "filesystem" 
    TemplateError    ErrorType = "template"
    InteractionError ErrorType = "interaction"
    ConfigError      ErrorType = "config"
)
```

#### Recovery Strategies
```go
func (cmd *CreateCommand) handleError(err error, context CreateContext) error {
    switch createErr := err.(type) {
    case *CreateError:
        switch createErr.Type {
        case ValidationError:
            return cmd.promptForCorrection(createErr, context)
        case FileSystemError:
            return cmd.handleFileSystemError(createErr, context)
        case TemplateError:
            return cmd.fallbackToBasicTemplates(context)
        default:
            return createErr
        }
    default:
        return fmt.Errorf("unexpected error during workflow creation: %w", err)
    }
}
```

### Testing Requirements

#### Unit Tests
```go
func TestWorkflowCreate_ValidInput(t *testing.T) {
    tests := []struct {
        name     string
        input    CreateInput
        expected WorkflowSpec
        wantErr  bool
    }{
        {
            name: "basic workflow creation",
            input: CreateInput{
                Name: "test-workflow",
                Category: "development", 
                Description: "Test workflow for unit testing",
                Interactive: false,
            },
            expected: WorkflowSpec{
                Name: "test-workflow",
                Category: "development",
                Description: "Test workflow for unit testing",
            },
            wantErr: false,
        },
        // Additional test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := NewCreateCommand()
            result, err := cmd.Execute(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected.Name, result.Name)
            assert.Equal(t, tt.expected.Category, result.Category)
            assert.Equal(t, tt.expected.Description, result.Description)
        })
    }
}
```

#### Integration Tests
```go
func TestWorkflowCreate_EndToEnd(t *testing.T) {
    tmpDir := t.TempDir()
    
    cmd := NewCreateCommand()
    cmd.OutputDir = tmpDir
    
    input := CreateInput{
        Name:        "integration-test-workflow",
        Category:    "development",
        Description: "End-to-end integration test workflow",
        Batch:       true,
        Phases:      []string{"plan", "execute", "validate"},
        Artifacts: []ArtifactDef{
            {Name: "Plan Document", Type: "document", Phase: "plan"},
            {Name: "Implementation", Type: "code", Phase: "execute"},  
            {Name: "Test Results", Type: "document", Phase: "validate"},
        },
    }
    
    result, err := cmd.Execute(input)
    require.NoError(t, err)
    
    // Verify directory structure
    workflowDir := filepath.Join(tmpDir, "integration-test-workflow")
    assert.DirExists(t, workflowDir)
    assert.FileExists(t, filepath.Join(workflowDir, "workflow.yml"))
    assert.FileExists(t, filepath.Join(workflowDir, "README.md"))
    assert.FileExists(t, filepath.Join(workflowDir, "GUIDE.md"))
    
    // Verify artifact directories
    assert.DirExists(t, filepath.Join(workflowDir, "plan-document"))
    assert.FileExists(t, filepath.Join(workflowDir, "plan-document", "template.md"))
    assert.FileExists(t, filepath.Join(workflowDir, "plan-document", "prompt.md"))
    
    // Verify workflow validation passes
    validator := NewWorkflowValidator()
    validationResult, err := validator.Validate(result.Name)
    require.NoError(t, err)
    assert.True(t, validationResult.IsValid())
}
```

### Performance Considerations

#### Optimization Targets
- **Creation Time**: <5 seconds for basic workflow
- **Template Generation**: <1 second per artifact
- **File I/O**: Batch writes to minimize system calls
- **Memory Usage**: <50MB for workflow creation process

#### Resource Management
```go
type ResourceManager struct {
    maxConcurrentTemplates int
    templateCache         map[string]Template
    fileWriteBuffer       int
}

func (rm *ResourceManager) OptimizeCreation(spec WorkflowSpec) error {
    // Pre-allocate template cache
    rm.templateCache = make(map[string]Template, len(spec.Artifacts))
    
    // Process artifacts in batches
    artifactBatches := rm.batchArtifacts(spec.Artifacts, rm.maxConcurrentTemplates)
    
    for _, batch := range artifactBatches {
        if err := rm.processBatch(batch); err != nil {
            return fmt.Errorf("batch processing failed: %w", err)
        }
    }
    
    return nil
}
```

### CLI Integration

#### Command Registration
```go
func init() {
    workflowCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new workflow",
    Long:  `Create a new workflow with interactive or batch mode options`,
    Args:  cobra.MaximumNArgs(1),
    RunE:  runCreateCommand,
}

func runCreateCommand(cmd *cobra.Command, args []string) error {
    // Extract options from command flags
    options := CreateOptions{
        Name:         getStringArg(args, 0),
        Template:     cmd.Flags().GetString("template"),
        Category:     cmd.Flags().GetString("category"),  
        Description:  cmd.Flags().GetString("description"),
        Author:       cmd.Flags().GetString("author"),
        Interactive:  cmd.Flags().GetBool("interactive"),
        Batch:        cmd.Flags().GetBool("batch"),
        OutputDir:    cmd.Flags().GetString("output-dir"),
        DryRun:       cmd.Flags().GetBool("dry-run"),
        Force:        cmd.Flags().GetBool("force"),
    }
    
    creator := workflow.NewCreator()
    return creator.Create(options)
}
```

## Related Documentation

- [Commands Overview](./commands.md) - All workflow commands
- [Apply Command Specification](./apply-command.md) - Workflow application
- [Technical Overview](./overview.md) - System architecture
- [Usage Guide](../../../usage/workflows/creating-workflows.md) - User documentation