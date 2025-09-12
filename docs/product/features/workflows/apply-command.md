# Workflow Apply Command Specification

**Version**: 1.0  
**Status**: Specification  
**Date**: 2025-09-12  
**Related**: [Commands Overview](./commands.md), [Create Command](./create-command.md), [Technical Overview](./overview.md)

## Overview

The `ddx workflow apply` command is the primary interface for executing existing workflows. It orchestrates the complete workflow lifecycle from initialization through completion, handling phase transitions, artifact generation, variable substitution, and state management.

## Command Signature

```bash
ddx workflow apply <name>[:<artifact>] [options]
```

### Arguments

- `<name>` (required): Workflow identifier
  - Must match an existing workflow in the catalog
  - Case-sensitive exact match
  - Examples: `development`, `incident-response`, `product-launch`
  
- `<artifact>` (optional): Specific artifact to apply
  - When specified, applies only the named artifact instead of the full workflow
  - Must be a valid artifact ID within the workflow
  - Useful for regenerating individual deliverables
  - Examples: `development:prd`, `incident-response:runbook`

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--phase <phase>` | string | first | Start from specific phase |
| `--variables <file>` | string | none | Variable values file (YAML/JSON) |
| `--set <key=value>` | string[] | none | Set individual variables (repeatable) |
| `--output-dir <path>` | string | `.ddx/workflows/<name>` | Custom output directory |
| `--force` | boolean | false | Overwrite existing files |
| `--interactive, -i` | boolean | false | Interactive variable collection |
| `--validate-only` | boolean | false | Validate without executing |
| `--resume` | boolean | false | Resume interrupted workflow |
| `--skip-phase <phase>` | string[] | none | Skip specific phases (repeatable) |
| `--parallel` | boolean | false | Enable parallel phase execution |
| `--timeout <duration>` | string | 30m | Workflow execution timeout |
| `--watch` | boolean | false | Monitor execution progress |
| `--save-state` | boolean | true | Persist workflow state |
| `--ai-model <model>` | string | config | AI model for prompt processing |

## Implementation Architecture

### Core Components

#### 1. WorkflowApplier
```go
type WorkflowApplier struct {
    catalogManager    CatalogManager
    templateEngine    TemplateEngine
    phaseOrchestrator PhaseOrchestrator
    artifactGenerator ArtifactGenerator
    stateManager      StateManager
    variableResolver  VariableResolver
    validator         WorkflowValidator
    progressReporter  ProgressReporter
}

type ApplyOptions struct {
    WorkflowName     string
    ArtifactName     string // Optional - for single artifact application
    StartPhase       string
    VariablesFile    string
    Variables        map[string]string
    OutputDir        string
    Force            bool
    Interactive      bool
    ValidateOnly     bool
    Resume           bool
    SkipPhases       []string
    Parallel         bool
    Timeout          time.Duration
    Watch            bool
    SaveState        bool
    AIModel          string
}
```

#### 2. Execution Pipeline
```go
func (wa *WorkflowApplier) Apply(options ApplyOptions) (*ApplyResult, error) {
    // 1. Load and validate workflow
    workflow, err := wa.loadWorkflow(options.WorkflowName)
    if err != nil {
        return nil, fmt.Errorf("failed to load workflow: %w", err)
    }
    
    // 2. Handle resume logic
    if options.Resume {
        return wa.resumeWorkflow(workflow, options)
    }
    
    // 3. Handle single artifact application
    if options.ArtifactName != "" {
        return wa.applySingleArtifact(workflow, options)
    }
    
    // 4. Execute full workflow
    return wa.executeFullWorkflow(workflow, options)
}
```

### Variable Resolution System

#### Variable Sources (Priority Order)
1. Command-line `--set` flags (highest priority)
2. Variables file (`--variables`)
3. Environment variables (DDX_VAR_*)
4. Project configuration (.ddx/config.yml)
5. Workflow defaults (workflow.yml)
6. Interactive prompts (lowest priority)

#### Variable Context Structure
```go
type VariableContext struct {
    Workflow      WorkflowVariables `json:"workflow"`
    Project       ProjectVariables  `json:"project"`
    Environment   EnvVariables      `json:"environment"`
    User          UserVariables     `json:"user"`
    Phase         PhaseVariables    `json:"phase"`
    Artifact      ArtifactVariables `json:"artifact"`
    Timestamp     time.Time         `json:"timestamp"`
    ExecutionID   string           `json:"execution_id"`
}

type WorkflowVariables struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Category    string `json:"category"`
    Author      string `json:"author"`
    Description string `json:"description"`
}

type ProjectVariables struct {
    Name        string            `json:"name"`
    Path        string            `json:"path"`
    GitRepo     string            `json:"git_repo,omitempty"`
    GitBranch   string            `json:"git_branch,omitempty"`
    Language    string            `json:"language,omitempty"`
    Framework   string            `json:"framework,omitempty"`
    Custom      map[string]string `json:"custom,omitempty"`
}
```

#### Variable Resolution Implementation
```go
func (vr *VariableResolver) Resolve(workflow *Workflow, options ApplyOptions) (*VariableContext, error) {
    context := &VariableContext{
        Timestamp:   time.Now(),
        ExecutionID: generateExecutionID(),
    }
    
    // Populate workflow variables
    context.Workflow = WorkflowVariables{
        Name:        workflow.Name,
        Version:     workflow.Version,
        Category:    workflow.Category,
        Author:      workflow.Author,
        Description: workflow.Description,
    }
    
    // Populate project variables
    project, err := vr.detectProject()
    if err != nil {
        return nil, fmt.Errorf("failed to detect project: %w", err)
    }
    context.Project = *project
    
    // Resolve custom variables by priority
    customVars := make(map[string]string)
    
    // 1. Workflow defaults
    for _, variable := range workflow.Variables {
        if variable.Default != nil {
            customVars[variable.Name] = fmt.Sprintf("%v", variable.Default)
        }
    }
    
    // 2. Project configuration
    projectConfig, err := vr.loadProjectConfig()
    if err == nil {
        for key, value := range projectConfig.Variables {
            customVars[key] = value
        }
    }
    
    // 3. Environment variables
    for key, value := range os.Environ() {
        if strings.HasPrefix(key, "DDX_VAR_") {
            varName := strings.TrimPrefix(key, "DDX_VAR_")
            customVars[strings.ToLower(varName)] = value
        }
    }
    
    // 4. Variables file
    if options.VariablesFile != "" {
        fileVars, err := vr.loadVariablesFile(options.VariablesFile)
        if err != nil {
            return nil, fmt.Errorf("failed to load variables file: %w", err)
        }
        for key, value := range fileVars {
            customVars[key] = value
        }
    }
    
    // 5. Command-line flags
    for key, value := range options.Variables {
        customVars[key] = value
    }
    
    // 6. Interactive prompts for missing required variables
    if options.Interactive {
        for _, variable := range workflow.Variables {
            if variable.Required && customVars[variable.Name] == "" {
                value, err := vr.promptForVariable(variable)
                if err != nil {
                    return nil, fmt.Errorf("failed to collect variable %s: %w", variable.Name, err)
                }
                customVars[variable.Name] = value
            }
        }
    }
    
    context.User.Custom = customVars
    
    return context, nil
}
```

### Phase Orchestration

#### Phase Execution Engine
```go
type PhaseOrchestrator struct {
    executor      PhaseExecutor
    validator     PhaseValidator
    stateManager  StateManager
    reporter      ProgressReporter
}

func (po *PhaseOrchestrator) Execute(workflow *Workflow, context *VariableContext, options ApplyOptions) error {
    phases := po.resolvePhaseDependencies(workflow.Phases, options)
    
    if options.Parallel {
        return po.executeParallel(phases, context, options)
    } else {
        return po.executeSequential(phases, context, options)
    }
}

func (po *PhaseOrchestrator) executeSequential(phases []Phase, context *VariableContext, options ApplyOptions) error {
    for i, phase := range phases {
        // Skip if requested
        if contains(options.SkipPhases, phase.ID) {
            po.reporter.ReportPhaseSkipped(phase.ID)
            continue
        }
        
        // Start phase
        po.reporter.ReportPhaseStarted(phase.ID, i+1, len(phases))
        
        // Validate entry criteria
        if err := po.validator.ValidateEntryConditions(phase, context); err != nil {
            return fmt.Errorf("phase %s entry conditions not met: %w", phase.ID, err)
        }
        
        // Execute phase
        phaseResult, err := po.executor.Execute(phase, context, options)
        if err != nil {
            po.stateManager.SavePhaseState(phase.ID, PhaseStateFailed, err)
            return fmt.Errorf("phase %s execution failed: %w", phase.ID, err)
        }
        
        // Validate exit criteria
        if err := po.validator.ValidateExitConditions(phase, phaseResult); err != nil {
            return fmt.Errorf("phase %s exit conditions not met: %w", phase.ID, err)
        }
        
        // Mark phase complete
        po.stateManager.SavePhaseState(phase.ID, PhaseStateCompleted, nil)
        po.reporter.ReportPhaseCompleted(phase.ID, phaseResult.Duration)
    }
    
    return nil
}
```

#### Parallel Phase Execution
```go
func (po *PhaseOrchestrator) executeParallel(phases []Phase, context *VariableContext, options ApplyOptions) error {
    dependencyGraph := po.buildDependencyGraph(phases)
    executionGroups := po.topologicalSort(dependencyGraph)
    
    for groupIndex, group := range executionGroups {
        po.reporter.ReportExecutionGroup(groupIndex+1, len(executionGroups), len(group))
        
        // Execute all phases in this group concurrently
        results := make(chan PhaseResult, len(group))
        errors := make(chan error, len(group))
        
        for _, phase := range group {
            go func(p Phase) {
                result, err := po.executor.Execute(p, context, options)
                if err != nil {
                    errors <- fmt.Errorf("phase %s failed: %w", p.ID, err)
                    return
                }
                results <- result
            }(phase)
        }
        
        // Wait for all phases in group to complete
        completed := 0
        for completed < len(group) {
            select {
            case result := <-results:
                po.reporter.ReportPhaseCompleted(result.PhaseID, result.Duration)
                completed++
            case err := <-errors:
                return err
            case <-time.After(options.Timeout):
                return fmt.Errorf("phase group %d timed out after %v", groupIndex+1, options.Timeout)
            }
        }
    }
    
    return nil
}
```

### Artifact Generation

#### Artifact Generator Interface
```go
type ArtifactGenerator interface {
    Generate(artifact Artifact, context *VariableContext, options GenerateOptions) (*GeneratedArtifact, error)
    Validate(artifact *GeneratedArtifact) (*ValidationResult, error)
    Save(artifact *GeneratedArtifact, outputPath string) error
}

type GenerateOptions struct {
    AIModel          string
    TemplateOnly     bool
    PromptContext    string
    OutputFormat     string
    ValidationLevel  ValidationLevel
}
```

#### Template Processing
```go
type TemplateProcessor struct {
    engine        *template.Template
    functions     template.FuncMap
    variableCache map[string]interface{}
}

func (tp *TemplateProcessor) Process(templateContent string, context *VariableContext) (string, error) {
    // Create template with custom functions
    tmpl, err := template.New("artifact").Funcs(tp.functions).Parse(templateContent)
    if err != nil {
        return "", fmt.Errorf("template parsing failed: %w", err)
    }
    
    // Prepare template context
    templateContext := tp.buildTemplateContext(context)
    
    // Execute template
    var output bytes.Buffer
    if err := tmpl.Execute(&output, templateContext); err != nil {
        return "", fmt.Errorf("template execution failed: %w", err)
    }
    
    return output.String(), nil
}

func (tp *TemplateProcessor) buildTemplateContext(context *VariableContext) map[string]interface{} {
    return map[string]interface{}{
        "workflow":     context.Workflow,
        "project":      context.Project,
        "environment":  context.Environment,
        "user":         context.User,
        "timestamp":    context.Timestamp,
        "execution_id": context.ExecutionID,
        
        // Utility functions
        "now":          func() string { return time.Now().Format(time.RFC3339) },
        "date":         func() string { return time.Now().Format("2006-01-02") },
        "uuid":         func() string { return generateUUID() },
        "slugify":      slugify,
        "title":        strings.Title,
        "upper":        strings.ToUpper,
        "lower":        strings.ToLower,
    }
}
```

#### AI-Powered Artifact Generation
```go
type AIArtifactGenerator struct {
    client       AIClient
    promptEngine PromptEngine
}

func (aag *AIArtifactGenerator) Generate(artifact Artifact, context *VariableContext, options GenerateOptions) (*GeneratedArtifact, error) {
    // Load and process prompt
    promptContent, err := aag.loadPrompt(artifact.PromptPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load prompt: %w", err)
    }
    
    processedPrompt, err := aag.promptEngine.Process(promptContent, context)
    if err != nil {
        return nil, fmt.Errorf("prompt processing failed: %w", err)
    }
    
    // Generate content using AI
    response, err := aag.client.Generate(AIRequest{
        Model:       options.AIModel,
        Prompt:      processedPrompt,
        MaxTokens:   artifact.MaxTokens,
        Temperature: artifact.Temperature,
        Context:     options.PromptContext,
    })
    if err != nil {
        return nil, fmt.Errorf("AI generation failed: %w", err)
    }
    
    // Parse and structure response
    content, err := aag.parseAIResponse(response, artifact.Type)
    if err != nil {
        return nil, fmt.Errorf("response parsing failed: %w", err)
    }
    
    return &GeneratedArtifact{
        ID:          artifact.ID,
        Name:        artifact.Name,
        Type:        artifact.Type,
        Content:     content,
        Metadata:    response.Metadata,
        GeneratedAt: time.Now(),
        Generator:   "ai:" + options.AIModel,
    }, nil
}
```

### State Management

#### Workflow State Structure
```go
type WorkflowState struct {
    WorkflowName   string                     `json:"workflow_name"`
    Version        string                     `json:"version"`
    ExecutionID    string                     `json:"execution_id"`
    Status         WorkflowStatus             `json:"status"`
    StartedAt      time.Time                  `json:"started_at"`
    CompletedAt    *time.Time                 `json:"completed_at,omitempty"`
    CurrentPhase   string                     `json:"current_phase"`
    Context        *VariableContext           `json:"context"`
    Phases         map[string]*PhaseState     `json:"phases"`
    Artifacts      map[string]*ArtifactState  `json:"artifacts"`
    Options        ApplyOptions               `json:"options"`
    Checkpoints    []Checkpoint               `json:"checkpoints"`
}

type PhaseState struct {
    ID          string        `json:"id"`
    Status      PhaseStatus   `json:"status"`
    StartedAt   *time.Time    `json:"started_at,omitempty"`
    CompletedAt *time.Time    `json:"completed_at,omitempty"`
    Duration    time.Duration `json:"duration"`
    Error       string        `json:"error,omitempty"`
    Attempts    int           `json:"attempts"`
}

type ArtifactState struct {
    ID           string           `json:"id"`
    Status       ArtifactStatus   `json:"status"`
    GeneratedAt  *time.Time       `json:"generated_at,omitempty"`
    Path         string           `json:"path,omitempty"`
    Hash         string           `json:"hash,omitempty"`
    ValidationResult *ValidationResult `json:"validation_result,omitempty"`
    Generator    string           `json:"generator,omitempty"`
}
```

#### State Persistence
```go
type StateManager struct {
    stateDir    string
    compression bool
}

func (sm *StateManager) SaveState(state *WorkflowState) error {
    stateFile := filepath.Join(sm.stateDir, fmt.Sprintf("%s-%s.json", state.WorkflowName, state.ExecutionID))
    
    // Create checkpoint before saving
    checkpoint := sm.createCheckpoint(state)
    state.Checkpoints = append(state.Checkpoints, checkpoint)
    
    // Marshal state
    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal state: %w", err)
    }
    
    // Compress if enabled
    if sm.compression {
        data, err = sm.compress(data)
        if err != nil {
            return fmt.Errorf("compression failed: %w", err)
        }
    }
    
    // Write atomically
    tmpFile := stateFile + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return fmt.Errorf("failed to write state file: %w", err)
    }
    
    if err := os.Rename(tmpFile, stateFile); err != nil {
        return fmt.Errorf("failed to commit state file: %w", err)
    }
    
    return nil
}
```

### Resume Functionality

#### Resume Logic
```go
func (wa *WorkflowApplier) resumeWorkflow(workflow *Workflow, options ApplyOptions) (*ApplyResult, error) {
    // Load existing state
    state, err := wa.stateManager.LoadLatestState(workflow.Name)
    if err != nil {
        return nil, fmt.Errorf("failed to load workflow state: %w", err)
    }
    
    // Validate state compatibility
    if err := wa.validateStateCompatibility(workflow, state); err != nil {
        return nil, fmt.Errorf("state compatibility check failed: %w", err)
    }
    
    // Find resume point
    resumePhase, err := wa.findResumePoint(state)
    if err != nil {
        return nil, fmt.Errorf("failed to determine resume point: %w", err)
    }
    
    wa.progressReporter.ReportResumeStarted(workflow.Name, resumePhase, state.ExecutionID)
    
    // Resume execution from the determined point
    resumeOptions := options
    resumeOptions.StartPhase = resumePhase
    
    return wa.executeFullWorkflow(workflow, resumeOptions)
}

func (wa *WorkflowApplier) findResumePoint(state *WorkflowState) (string, error) {
    // Find the first incomplete phase
    for phaseID, phaseState := range state.Phases {
        if phaseState.Status == PhaseStatusFailed || phaseState.Status == PhaseStatusActive {
            return phaseID, nil
        }
    }
    
    // If all phases are complete, workflow is done
    if state.Status == WorkflowStatusCompleted {
        return "", errors.New("workflow already completed")
    }
    
    // Find next phase to execute
    return wa.findNextPhase(state.CurrentPhase, state), nil
}
```

### Validation System

#### Multi-Level Validation
```go
type ValidationLevel int

const (
    ValidationBasic ValidationLevel = iota
    ValidationStandard
    ValidationStrict
    ValidationComprehensive
)

type WorkflowValidator struct {
    schemaValidator   SchemaValidator
    contentValidator  ContentValidator
    linkValidator     LinkValidator
    customValidators  []CustomValidator
}

func (wv *WorkflowValidator) Validate(workflow *Workflow, level ValidationLevel) (*ValidationResult, error) {
    result := &ValidationResult{
        WorkflowName: workflow.Name,
        Level:        level,
        StartedAt:    time.Now(),
        Errors:       []ValidationError{},
        Warnings:     []ValidationWarning{},
    }
    
    // Schema validation (always performed)
    if err := wv.schemaValidator.Validate(workflow); err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Type:    "schema",
            Message: err.Error(),
            Field:   "workflow",
        })
    }
    
    // Content validation (standard and above)
    if level >= ValidationStandard {
        contentErrors := wv.contentValidator.Validate(workflow)
        result.Errors = append(result.Errors, contentErrors...)
    }
    
    // Link validation (strict and above)
    if level >= ValidationStrict {
        linkWarnings := wv.linkValidator.Validate(workflow)
        result.Warnings = append(result.Warnings, linkWarnings...)
    }
    
    // Custom validation (comprehensive)
    if level >= ValidationComprehensive {
        for _, validator := range wv.customValidators {
            customResult, err := validator.Validate(workflow)
            if err != nil {
                return nil, fmt.Errorf("custom validation failed: %w", err)
            }
            result.Errors = append(result.Errors, customResult.Errors...)
            result.Warnings = append(result.Warnings, customResult.Warnings...)
        }
    }
    
    result.CompletedAt = time.Now()
    result.Duration = result.CompletedAt.Sub(result.StartedAt)
    result.IsValid = len(result.Errors) == 0
    
    return result, nil
}
```

### Progress Reporting

#### Real-time Progress Updates
```go
type ProgressReporter interface {
    ReportWorkflowStarted(workflowName string, totalPhases int)
    ReportPhaseStarted(phaseID string, phaseNum, totalPhases int)
    ReportPhaseProgress(phaseID string, progress float64, message string)
    ReportPhaseCompleted(phaseID string, duration time.Duration)
    ReportPhaseSkipped(phaseID string)
    ReportArtifactGenerated(artifactID string, path string)
    ReportWorkflowCompleted(workflowName string, totalDuration time.Duration)
    ReportError(err error)
}

type ConsoleProgressReporter struct {
    output     io.Writer
    spinner    *spinner.Spinner
    startTime  time.Time
    lastUpdate time.Time
}

func (cpr *ConsoleProgressReporter) ReportPhaseStarted(phaseID string, phaseNum, totalPhases int) {
    cpr.spinner.Start()
    fmt.Fprintf(cpr.output, "Phase %d/%d: %s\n", phaseNum, totalPhases, phaseID)
    cpr.lastUpdate = time.Now()
}

func (cpr *ConsoleProgressReporter) ReportPhaseProgress(phaseID string, progress float64, message string) {
    if time.Since(cpr.lastUpdate) < time.Second {
        return // Rate limit updates
    }
    
    progressBar := cpr.buildProgressBar(progress, 40)
    fmt.Fprintf(cpr.output, "\r%s [%s] %s", phaseID, progressBar, message)
    cpr.lastUpdate = time.Now()
}
```

### Error Handling and Recovery

#### Error Classification
```go
type ErrorCategory int

const (
    ErrorCategoryValidation ErrorCategory = iota
    ErrorCategoryTemplate
    ErrorCategoryAI
    ErrorCategoryFileSystem
    ErrorCategoryNetwork
    ErrorCategoryTimeout
    ErrorCategoryUserCancelled
    ErrorCategoryStateCorruption
)

type WorkflowError struct {
    Category    ErrorCategory `json:"category"`
    Phase       string        `json:"phase,omitempty"`
    Artifact    string        `json:"artifact,omitempty"`
    Message     string        `json:"message"`
    Cause       error         `json:"cause,omitempty"`
    Recoverable bool          `json:"recoverable"`
    Timestamp   time.Time     `json:"timestamp"`
}
```

#### Recovery Strategies
```go
type RecoveryStrategy interface {
    CanRecover(err *WorkflowError) bool
    Recover(err *WorkflowError, context *WorkflowContext) error
}

type RetryRecoveryStrategy struct {
    maxRetries int
    backoff    time.Duration
}

func (rrs *RetryRecoveryStrategy) Recover(err *WorkflowError, context *WorkflowContext) error {
    if !err.Recoverable {
        return err
    }
    
    for attempt := 1; attempt <= rrs.maxRetries; attempt++ {
        time.Sleep(rrs.backoff * time.Duration(attempt))
        
        // Retry the operation
        if retryErr := context.RetryLastOperation(); retryErr == nil {
            return nil // Recovery successful
        }
        
        if attempt == rrs.maxRetries {
            return fmt.Errorf("recovery failed after %d attempts: %w", rrs.maxRetries, err)
        }
    }
    
    return err
}
```

### Testing Framework

#### Integration Test Structure
```go
func TestWorkflowApply_FullExecution(t *testing.T) {
    // Setup test environment
    tmpDir := t.TempDir()
    workflowCatalog := setupTestCatalog(tmpDir)
    
    applier := NewWorkflowApplier(WorkflowApplierConfig{
        CatalogPath: workflowCatalog,
        StateDir:    filepath.Join(tmpDir, "state"),
        OutputDir:   filepath.Join(tmpDir, "output"),
    })
    
    // Test full workflow execution
    options := ApplyOptions{
        WorkflowName: "test-workflow",
        Interactive:  false,
        Variables: map[string]string{
            "project_name": "test-project",
            "author":       "test-user",
        },
    }
    
    result, err := applier.Apply(options)
    require.NoError(t, err)
    assert.Equal(t, WorkflowStatusCompleted, result.Status)
    assert.Len(t, result.GeneratedArtifacts, 3)
    
    // Verify artifacts were created
    for _, artifact := range result.GeneratedArtifacts {
        assert.FileExists(t, artifact.Path)
        
        // Validate artifact content
        content, err := os.ReadFile(artifact.Path)
        require.NoError(t, err)
        assert.Contains(t, string(content), "test-project")
        assert.Contains(t, string(content), "test-user")
    }
    
    // Verify state was saved
    stateFile := filepath.Join(tmpDir, "state", fmt.Sprintf("test-workflow-%s.json", result.ExecutionID))
    assert.FileExists(t, stateFile)
}
```

### Performance Optimization

#### Caching Strategy
```go
type CacheManager struct {
    templateCache    *lru.Cache
    promptCache      *lru.Cache
    validationCache  *lru.Cache
    workflowCache    *lru.Cache
}

func (cm *CacheManager) GetTemplate(templatePath string) (*Template, error) {
    if cached, ok := cm.templateCache.Get(templatePath); ok {
        return cached.(*Template), nil
    }
    
    template, err := cm.loadTemplate(templatePath)
    if err != nil {
        return nil, err
    }
    
    cm.templateCache.Add(templatePath, template)
    return template, nil
}
```

#### Async Processing
```go
type AsyncArtifactGenerator struct {
    generator    ArtifactGenerator
    workers      int
    queue        chan GenerationTask
    results      chan GenerationResult
    workerPool   *sync.WaitGroup
}

func (aag *AsyncArtifactGenerator) ProcessArtifacts(artifacts []Artifact, context *VariableContext) <-chan GenerationResult {
    results := make(chan GenerationResult, len(artifacts))
    
    // Start worker pool
    for i := 0; i < aag.workers; i++ {
        aag.workerPool.Add(1)
        go aag.worker(results)
    }
    
    // Queue all artifacts
    for _, artifact := range artifacts {
        aag.queue <- GenerationTask{
            Artifact: artifact,
            Context:  context,
        }
    }
    
    // Close queue and wait for completion
    close(aag.queue)
    go func() {
        aag.workerPool.Wait()
        close(results)
    }()
    
    return results
}
```

### CLI Output Format

#### Success Output
```
Applying workflow: development
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Workflow validation passed
✓ Variable resolution completed
✓ Output directory prepared: .ddx/workflows/development

Phase 1/6: Define
  ✓ Requirements Document generated → artifacts/requirements.md
  ✓ Project Charter generated → artifacts/charter.md
  ✓ Phase completed in 1m 32s

Phase 2/6: Design  
  ✓ Architecture Document generated → artifacts/architecture.md
  ✓ Database Schema generated → artifacts/schema.sql
  ✓ Phase completed in 2m 15s

Phase 3/6: Implement
  ⚠ Source Code generation skipped (no template)
  ✓ Configuration Files generated → config/
  ✓ Phase completed in 0m 45s

...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Workflow completed successfully! 🎉

Artifacts generated: 8
Total duration: 12m 34s
Output directory: .ddx/workflows/development

Next steps:
- Review generated artifacts
- Run 'ddx workflow status development' to see details
- Use 'ddx workflow validate development' to verify completeness
```

#### Error Output
```
Applying workflow: development
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✗ Workflow application failed

Error in phase: Design
Artifact: Architecture Document
Cause: Template processing failed

Details:
  Template variable '{{database_type}}' is undefined
  Required variables missing: database_type, deployment_target
  
Recovery options:
  1. Provide missing variables:
     ddx workflow apply development --set database_type=postgresql --set deployment_target=cloud
  
  2. Resume from Design phase:
     ddx workflow apply development --resume --phase design
     
  3. Skip Design phase:
     ddx workflow apply development --skip-phase design

For more details: ddx workflow status development
```

## Related Documentation

- [Commands Overview](./commands.md) - Complete command reference
- [Create Command Specification](./create-command.md) - Workflow creation
- [Technical Overview](./overview.md) - System architecture
- [Usage Guide](../../../usage/workflows/overview.md) - User documentation