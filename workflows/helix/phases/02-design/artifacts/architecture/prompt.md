# Architecture Documentation Generation Prompt

Create comprehensive architecture documentation that visualizes system structure, components, and relationships using the C4 model and supporting diagrams.

## Storage Location

Store architecture documentation at: `docs/helix/02-design/architecture.md`

## Purpose

Architecture documentation provides:
- Visual system understanding for all stakeholders
- Clear component boundaries and responsibilities
- Data flow and interaction patterns
- Deployment and infrastructure context
- Reference for implementation decisions

## Key Requirements

### 1. C4 Model Implementation

Follow the C4 (Context, Containers, Components, Code) model for architectural views:

#### Level 1: System Context Diagram
Show how the system fits in the larger ecosystem:
- External users (personas, not generic "users")
- External systems and services
- High-level system purpose
- Key relationships and interactions

```mermaid
C4Context
    title System Context for DDX Platform

    Person(dev, "Developer", "Software developer using DDX for project templates")
    Person(contributor, "Contributor", "Community member sharing patterns")

    System(ddx, "DDX Platform", "CLI toolkit for AI-assisted development")

    System_Ext(github, "GitHub", "Code hosting and version control")
    System_Ext(ai_service, "AI Service", "Claude, GPT for AI assistance")
    System_Ext(registry, "Template Registry", "Shared templates and patterns")

    Rel(dev, ddx, "Uses CLI commands")
    Rel(contributor, ddx, "Shares templates")
    Rel(ddx, github, "Syncs templates and patterns")
    Rel(ddx, ai_service, "Gets AI assistance")
    Rel(ddx, registry, "Downloads/uploads resources")
```

#### Level 2: Container Diagram
Show major architectural components:
- Applications, databases, file systems
- Technology choices for each container
- Communication protocols and data formats
- Security boundaries

```mermaid
C4Container
    title Container Diagram for DDX Platform

    Person(dev, "Developer")

    Container(cli, "DDX CLI", "Go", "Command-line interface for DDX operations")
    Container(config, "Config Manager", "Go", "Manages .ddx.yml and user settings")
    Container(template_engine, "Template Engine", "Go", "Processes templates with variables")
    Container(git_sync, "Git Sync", "Go", "Handles git subtree operations")
    Container(ai_client, "AI Client", "Go", "Interfaces with AI services")

    ContainerDb(local_store, "Local Store", "File System", "Templates, patterns, cache")

    System_Ext(github, "GitHub")
    System_Ext(ai_service, "AI Service")

    Rel(dev, cli, "Commands", "CLI")
    Rel(cli, config, "Read/write config")
    Rel(cli, template_engine, "Process templates")
    Rel(cli, git_sync, "Sync resources")
    Rel(cli, ai_client, "AI requests")
    Rel(template_engine, local_store, "Read templates")
    Rel(git_sync, github, "Pull/push", "Git")
    Rel(ai_client, ai_service, "API calls", "HTTPS")
```

#### Level 3: Component Diagram
Detail internal structure of key containers:
- Classes, interfaces, services
- Component responsibilities
- Dependencies and relationships

```mermaid
C4Component
    title Component Diagram for DDX CLI

    Container(cli, "DDX CLI", "Go", "Main CLI application")

    Component(cmd_init, "Init Command", "Go", "Initialize DDX in project")
    Component(cmd_apply, "Apply Command", "Go", "Apply templates/patterns")
    Component(cmd_list, "List Command", "Go", "List available resources")
    Component(cmd_update, "Update Command", "Go", "Update from master repo")

    Component(config_mgr, "Config Manager", "Go", "Handle .ddx.yml files")
    Component(template_proc, "Template Processor", "Go", "Variable substitution")
    Component(git_ops, "Git Operations", "Go", "Git subtree commands")
    Component(validator, "Validator", "Go", "Validate configs and templates")

    Rel(cmd_init, config_mgr, "Create config")
    Rel(cmd_apply, template_proc, "Process templates")
    Rel(cmd_apply, git_ops, "Sync if needed")
    Rel(cmd_list, config_mgr, "Read available resources")
    Rel(cmd_update, git_ops, "Pull updates")
    Rel(template_proc, validator, "Validate before processing")
```

### 2. Data Flow Diagrams

Show how information moves through the system:

```mermaid
flowchart TD
    A[User runs 'ddx apply template'] --> B[Parse command arguments]
    B --> C{Template exists locally?}
    C -->|No| D[Fetch from git repository]
    C -->|Yes| E[Load template files]
    D --> E
    E --> F[Parse template variables]
    F --> G[Prompt user for missing variables]
    G --> H[Process template with variables]
    H --> I[Validate output files]
    I --> J[Write files to project]
    J --> K[Update .ddx.yml with applied template]
```

### 3. Deployment Architecture

Show infrastructure and deployment context:

```mermaid
C4Deployment
    title Deployment Diagram for DDX

    Deployment_Node(dev_machine, "Developer Machine", "macOS/Linux/Windows") {
        Container(ddx_cli, "DDX CLI", "Go binary")
        ContainerDb(local_cache, "Local Cache", ".ddx/ directory")
    }

    Deployment_Node(github_cloud, "GitHub", "Cloud hosting") {
        Container(master_repo, "Master Repository", "Git repository with all templates")
        Container(user_projects, "User Projects", "Individual project repositories")
    }

    Deployment_Node(ai_cloud, "AI Service", "Cloud API") {
        Container(ai_api, "AI API", "Claude/GPT API endpoints")
    }

    Rel(ddx_cli, master_repo, "git subtree pull/push", "HTTPS")
    Rel(ddx_cli, user_projects, "git operations", "SSH/HTTPS")
    Rel(ddx_cli, ai_api, "AI requests", "HTTPS")
```

## Diagram Types and Tools

### Mermaid Diagrams
Use Mermaid for most architectural diagrams:
- C4 diagrams for system architecture
- Flowcharts for data flow
- Sequence diagrams for interactions
- Entity relationship diagrams for data models

### Specialized Diagrams

#### Network Architecture
```mermaid
graph TB
    subgraph "Developer Environment"
        CLI[DDX CLI]
        Cache[Local Cache]
    end

    subgraph "External Services"
        GitHub[GitHub API]
        AI[AI Service API]
        Registry[Template Registry]
    end

    CLI -.->|HTTPS| GitHub
    CLI -.->|HTTPS| AI
    CLI -.->|HTTPS| Registry
    CLI -->|File I/O| Cache
```

#### Security Architecture
```mermaid
graph TD
    subgraph "Security Boundaries"
        subgraph "Local System"
            CLI[DDX CLI Process]
            Files[File System Access]
        end

        subgraph "Network Security"
            TLS[TLS 1.3 Encryption]
            Auth[API Key Authentication]
        end

        subgraph "External Services"
            GitHub[GitHub]
            AI[AI Service]
        end
    end

    CLI --> TLS
    TLS --> Auth
    Auth --> GitHub
    Auth --> AI
```

## Documentation Structure

### Main Architecture Document
```markdown
# System Architecture

## Overview
[Brief system description and purpose]

## Architecture Principles
[Key principles guiding design decisions]

## System Context
[Context diagram and external interfaces]

## Container Architecture
[Container diagram and technology choices]

## Component Design
[Key component diagrams and responsibilities]

## Data Architecture
[Data flow and storage design]

## Security Architecture
[Security measures and boundaries]

## Deployment Architecture
[Infrastructure and deployment model]

## Quality Attributes
[How architecture supports NFRs]

## Architecture Decisions
[Links to relevant ADRs]
```

## Quality Attributes Mapping

### Performance
- Latency requirements and design implications
- Throughput considerations
- Caching strategies
- Resource utilization

### Scalability
- Horizontal vs vertical scaling approaches
- Load distribution mechanisms
- Capacity planning considerations
- Performance bottlenecks

### Security
- Security boundaries and trust zones
- Authentication and authorization flows
- Data protection mechanisms
- Attack surface analysis

### Maintainability
- Module boundaries and coupling
- Dependency management
- Testing strategies
- Deployment automation

### Reliability
- Failure modes and recovery
- Monitoring and alerting
- Backup and disaster recovery
- Circuit breakers and resilience

## Architecture Validation

### Principle Compliance
Verify architecture aligns with established principles:
- Simplicity over complexity
- Composition over inheritance
- Explicit over implicit
- Fail fast and safe
- Single responsibility

### Scenario Analysis
Test architecture against key scenarios:
- Normal operation flows
- Error and exception handling
- Performance under load
- Security attack vectors
- Evolution and change scenarios

### Architecture Review Checklist
- [ ] All major components identified
- [ ] Component responsibilities clear
- [ ] Dependencies explicitly shown
- [ ] Data flow documented
- [ ] Security boundaries defined
- [ ] Performance characteristics addressed
- [ ] Scalability approach defined
- [ ] Deployment model documented
- [ ] Quality attributes mapped
- [ ] Principles compliance verified

## Integration with Design Phase

Architecture documentation enables the Design phase by:
1. Providing visual system understanding
2. Establishing component boundaries
3. Defining integration points
4. Supporting technical decision making
5. Enabling implementation planning

Remember: Architecture diagrams are communication tools. They should clarify, not complicate. Focus on the essential structures and relationships that matter for understanding and implementing the system.