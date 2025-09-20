# DDx Architecture Overview

> **Last Updated**: 2025-01-15
> **Status**: Active
> **Owner**: DDx Team

## Overview

DDx follows a modular, extensible architecture designed to support document-driven development workflows across diverse project types. This document presents the architecture using the C4 model (Context, Container, Component, Code) for clear visualization at different abstraction levels.

## C4 Model Architecture

### Level 1: System Context Diagram

The highest level view showing DDx in the context of external systems and users.

```mermaid
graph TB
    subgraph "Users"
        DEV[Developer]
        TEAM[Development Team]
        CONTRIB[Contributor]
    end

    subgraph "DDx System"
        DDX[DDx CLI Toolkit]
    end

    subgraph "External Systems"
        GH[GitHub/GitLab]
        NPM[NPM Registry]
        PIP[PyPI]
        AI[AI Services]
    end

    DEV --> |Uses| DDX
    TEAM --> |Collaborates via| DDX
    CONTRIB --> |Contributes to| DDX

    DDX --> |Fetches/Pushes| GH
    DDX --> |Downloads packages| NPM
    DDX --> |Downloads packages| PIP
    DDX --> |Integrates with| AI

    style DDX fill:#f9f,stroke:#333,stroke-width:4px
```

### Level 2: Container Diagram

Shows the high-level shape of the software architecture and how responsibilities are distributed.

```mermaid
graph TB
    subgraph "DDx System Boundary"
        subgraph "CLI Application [Go Binary]"
            CLI[Command Interface<br/>Cobra Framework]
            CORE[Core Engine<br/>Business Logic]
            SYNC[Sync Manager<br/>Git Operations]
            TPL[Template Engine<br/>Variable Substitution]
            VAL[Validator<br/>Security & Quality]
        end

        subgraph "Configuration [File System]"
            CONF[.ddx.yml<br/>Project Config]
            STATE[Sync State<br/>JSON Database]
            CACHE[Resource Cache<br/>Local Storage]
        end

        subgraph "Content Repository [Git]"
            TMPL[Templates<br/>Project Boilerplates]
            PTRN[Patterns<br/>Code Examples]
            PRMT[Prompts<br/>AI Instructions]
            CFG[Configs<br/>Tool Settings]
        end
    end

    subgraph "User Project"
        PROJ[Project Files]
        GIT[Git Repository]
    end

    CLI --> CORE
    CORE --> SYNC
    CORE --> TPL
    CORE --> VAL

    CORE --> CONF
    SYNC --> STATE
    CORE --> CACHE

    SYNC --> TMPL
    SYNC --> PTRN
    TPL --> PRMT
    TPL --> CFG

    CORE --> PROJ
    SYNC --> GIT
```

### Level 3: Component Diagram - CLI Application

Zooms into the CLI application container to show its internal components.

```mermaid
graph TB
    subgraph "CLI Application Components"
        subgraph "Command Layer"
            INIT[init.go<br/>Initialize DDx]
            LIST[list.go<br/>List Resources]
            APPLY[apply.go<br/>Apply Resources]
            UPDATE[update.go<br/>Sync Upstream]
            CONTRIB[contribute.go<br/>Share Changes]
            DIAG[diagnose.go<br/>Health Check]
        end

        subgraph "Core Services"
            CFG_SVC[Config Service<br/>Viper Integration]
            GIT_SVC[Git Service<br/>Subtree Operations]
            TPL_SVC[Template Service<br/>Processing Engine]
            VAL_SVC[Validation Service<br/>Input & Security]
            ASSET_SVC[Asset Service<br/>Resource Discovery]
        end

        subgraph "Infrastructure"
            FS[File System<br/>I/O Operations]
            NET[Network<br/>HTTP/Git Client]
            LOG[Logger<br/>Structured Logging]
            ERR[Error Handler<br/>Recovery & Reporting]
        end
    end

    INIT --> CFG_SVC
    INIT --> GIT_SVC

    LIST --> ASSET_SVC
    LIST --> CFG_SVC

    APPLY --> TPL_SVC
    APPLY --> VAL_SVC
    APPLY --> FS

    UPDATE --> GIT_SVC
    UPDATE --> NET

    CONTRIB --> GIT_SVC
    CONTRIB --> VAL_SVC

    DIAG --> VAL_SVC
    DIAG --> ASSET_SVC

    CFG_SVC --> FS
    GIT_SVC --> NET
    TPL_SVC --> FS

    style INIT fill:#e1f5fe
    style LIST fill:#e1f5fe
    style APPLY fill:#e1f5fe
    style UPDATE fill:#e1f5fe
    style CONTRIB fill:#e1f5fe
    style DIAG fill:#e1f5fe
```

### Level 3: Component Diagram - Sync System

Details of the synchronization subsystem.

```mermaid
graph TB
    subgraph "Synchronization Components"
        subgraph "Sync Manager"
            PULL[Pull Service<br/>Fetch Updates]
            PUSH[Push Service<br/>Contribute]
            MERGE[Merge Engine<br/>3-Way Merge]
            CONFLICT[Conflict Resolver<br/>Resolution Strategies]
        end

        subgraph "State Management"
            STATE[State Tracker<br/>Sync Metadata]
            HIST[History Logger<br/>Audit Trail]
            BACKUP[Backup Manager<br/>Snapshots]
        end

        subgraph "Git Integration"
            SUBTREE[Subtree Wrapper<br/>Git Commands]
            DIFF[Diff Engine<br/>Change Detection]
            COMMIT[Commit Builder<br/>Message Generation]
        end
    end

    PULL --> SUBTREE
    PULL --> MERGE
    PULL --> STATE

    PUSH --> SUBTREE
    PUSH --> COMMIT
    PUSH --> STATE

    MERGE --> CONFLICT
    MERGE --> DIFF

    CONFLICT --> BACKUP

    STATE --> HIST

    style PULL fill:#fff3e0
    style PUSH fill:#fff3e0
    style MERGE fill:#fff3e0
```

### Level 4: Deployment Diagram

Shows how DDx is deployed across different platforms.

```mermaid
graph TB
    subgraph "Developer Machine"
        subgraph "macOS/Linux/Windows"
            BIN[DDx Binary<br/>~15MB]
            CONF_LOCAL[Local Config<br/>~/.ddx/]
            PROJ_LOCAL[Project Files]
        end
    end

    subgraph "Version Control"
        subgraph "GitHub/GitLab"
            MASTER[Master Repo<br/>ddx-tools/ddx]
            FORK[User Fork<br/>Contributions]
            PROJ_REPO[Project Repo<br/>With Subtree]
        end
    end

    subgraph "Package Registries"
        BREW[Homebrew<br/>macOS]
        APT[APT/YUM<br/>Linux]
        CHOCO[Chocolatey<br/>Windows]
    end

    BIN --> CONF_LOCAL
    BIN --> PROJ_LOCAL

    BIN --> |Pull/Push| MASTER
    BIN --> |Contribute| FORK
    PROJ_LOCAL --> |Sync| PROJ_REPO

    BREW --> BIN
    APT --> BIN
    CHOCO --> BIN

    style BIN fill:#f0f4c3
```

## Data Flow Diagrams

### Resource Application Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Config
    participant Template
    participant Validator
    participant FileSystem

    User->>CLI: ddx apply template/nextjs
    CLI->>Config: Load .ddx.yml
    Config-->>CLI: Configuration data
    CLI->>Template: Load template
    Template-->>CLI: Template files
    CLI->>CLI: Resolve variables
    CLI->>Validator: Validate inputs
    Validator-->>CLI: Validation result
    CLI->>FileSystem: Write files
    FileSystem-->>CLI: Success
    CLI-->>User: Template applied
```

### Synchronization Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Git
    participant Upstream
    participant Merger
    participant State

    User->>CLI: ddx update
    CLI->>State: Check last sync
    State-->>CLI: Last commit SHA
    CLI->>Git: Fetch upstream
    Git->>Upstream: Pull changes
    Upstream-->>Git: New commits
    Git-->>CLI: Changes fetched
    CLI->>Merger: Merge changes
    Merger-->>CLI: Merge result
    CLI->>State: Update sync state
    State-->>CLI: State saved
    CLI-->>User: Update complete
```

## Component Interaction Matrix

| Component | Interacts With | Protocol | Purpose |
|-----------|---------------|----------|---------|
| CLI Commands | Core Services | Function calls | Execute business logic |
| Core Services | File System | OS APIs | Read/write files |
| Git Service | Remote Repos | Git protocol | Sync operations |
| Template Engine | Variables | In-memory | Substitution |
| Validator | Security Rules | Pattern matching | Input validation |
| Config Service | YAML Parser | Viper library | Configuration management |
| Network Client | GitHub API | HTTPS | PR creation |
| Logger | Output Stream | stdout/stderr | User feedback |

## Core Components

### 1. CLI Application (`/cli/`)

The command-line interface built with Go and Cobra framework.

**Key Components:**
- **Command Layer** (`/cli/cmd/`) - Command implementations
- **Internal Packages** (`/cli/internal/`) - Core business logic
- **Configuration** (`/cli/internal/config/`) - Viper-based configuration

**Commands:**
- `init` - Initialize DDx in a project
- `list` - Display available resources
- `apply` - Apply templates/patterns
- `diagnose` - Analyze project health
- `update` - Update toolkit resources
- `contribute` - Share improvements

### 2. Content Repository

Centralized repository of reusable resources.

**Structure:**
- **Templates** (`/templates/`) - Project boilerplates
- **Patterns** (`/patterns/`) - Code patterns and examples
- **Prompts** (`/prompts/`) - AI assistance prompts
- **Configs** (`/configs/`) - Tool configurations

### 3. Project Integration

How DDx integrates with user projects.

**Integration Methods:**
- **Git Subtree** - Primary method for syncing resources
- **Direct Copy** - Alternative for simple resource application
- **Symlinks** - Development mode for local testing

## Data Flow

1. **Initialization**
   - User runs `ddx init`
   - Creates `.ddx.yml` configuration
   - Sets up git subtree (optional)

2. **Resource Application**
   - User runs `ddx apply <resource>`
   - CLI reads configuration
   - Fetches resource from repository
   - Applies with variable substitution
   - Updates project files

3. **Contribution Flow**
   - User modifies resources
   - Runs `ddx contribute`
   - Changes pushed to subtree
   - PR created to master repository

## Design Principles

### 1. Modularity
- Loosely coupled components
- Clear separation of concerns
- Plugin-based extensibility

### 2. Simplicity
- Minimal dependencies
- Clear, intuitive interfaces
- Convention over configuration

### 3. Portability
- Cross-platform support (macOS, Linux, Windows)
- No external runtime requirements
- Self-contained binaries

### 4. Version Control Integration
- Git-native workflows
- Subtree for reliable syncing
- Preserves project history

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | Go 1.21+ | Performance, portability, single binary |
| CLI Framework | Cobra | Industry standard, feature-rich |
| Configuration | Viper | Flexible configuration management |
| Version Control | Git | Universal adoption, subtree support |
| Build System | Make | Simple, cross-platform |
| Testing | Go testing | Built-in, comprehensive |

## Extension Points

### 1. Custom Templates
Users can create project-specific templates in `.ddx/templates/`.

### 2. Pattern Libraries
Organizations can maintain private pattern repositories.

### 3. Prompt Collections
Teams can develop domain-specific AI prompts.

### 4. Tool Configurations
Shareable configurations for linters, formatters, etc.

### 5. MCP Server Configuration
Project-local MCP server definitions for Claude Code integration.

## MCP Server Management Architecture

### MCP Installation Strategy

DDx implements a **project-local, dependency-managed** approach for MCP (Model Context Protocol) servers:

#### 1. Local Dependency Management

**Package Manager Abstraction:**
DDx supports multiple package managers through intelligent detection and abstraction:

- **Auto-Detection**: Automatically detects package manager from lock files
  - `pnpm-lock.yaml` → pnpm (recommended for efficiency)
  - `yarn.lock` → yarn
  - `bun.lockb` → bun
  - `package-lock.json` or none → npm (default)
- **Configuration Override**: Can be specified in `.ddx.yml`
- **Unified Interface**: Same DDx commands work with any package manager

**Package Installation Strategy:**
- **Local Dependencies**: MCP servers installed as project-local packages in `node_modules/`
- **Package.json Management**: Automatic `package.json` creation/update for MCP dependencies
- **Version Locking**: Specific versions pinned for reproducible environments
- **Team Synchronization**: Dependencies and lock files committed to version control

**Installation Flow:**
```
ddx mcp install <server> → Detect PM → [npm/pnpm/yarn/bun] install → Update Claude config → Verify
```

#### 2. Configuration Hierarchy Design

DDx implements a **local-first** configuration strategy for MCP servers:

1. **Project-Local** (Primary): `.claude/settings.local.json`
   - Project-specific MCP server definitions
   - References local `node_modules/` installations
   - Version-controlled with project code
   - Shared across team members
   - Isolated from other projects

2. **Global** (Secondary): `~/.claude/settings.local.json`
   - User-wide MCP server definitions
   - Uses global npm packages or system installations
   - Personal development preferences
   - Available across all projects

#### 3. MCP Server Lifecycle Management

**Installation Phase:**
```mermaid
graph LR
    A[ddx mcp install] --> B[Check package.json]
    B --> C[npm install package]
    C --> D[Update Claude config]
    D --> E[Verify connectivity]
    E --> F[Success confirmation]
```

**Configuration Generation:**
- **Command Path Resolution**: Resolves to local `node_modules/.bin/` executables
- **Environment Setup**: Configures project-specific environment variables
- **Path Variables**: Substitutes `$PWD` with actual project directory
- **Validation**: Ensures all required dependencies are available

#### 4. Dependency Architecture

**Package Management:**
```json
{
  "devDependencies": {
    "@modelcontextprotocol/server-filesystem": "^1.0.0",
    "@modelcontextprotocol/server-sequential-thinking": "^1.0.0",
    "@playwright/mcp": "^1.0.0"
  }
}
```

**Generated Configuration (adapts to package manager):**
```json
// npm/pnpm:
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",  // or "pnpx" for pnpm
      "args": ["@modelcontextprotocol/server-filesystem", "$PWD"],
      "env": {}
    }
  }
}

// yarn:
{
  "mcpServers": {
    "filesystem": {
      "command": "yarn",
      "args": ["dlx", "@modelcontextprotocol/server-filesystem", "$PWD"],
      "env": {}
    }
  }
}

// bun:
{
  "mcpServers": {
    "filesystem": {
      "command": "bunx",
      "args": ["@modelcontextprotocol/server-filesystem", "$PWD"],
      "env": {}
    }
  }
}
```

#### 5. Design Principles

- **Project Isolation**: Each project manages its own MCP dependencies and versions
- **Team Consistency**: `package.json` and `.claude/` configs ensure identical AI tooling across team
- **Reproducible Environments**: Locked dependency versions prevent drift
- **Graceful Degradation**: Falls back to global configuration when local not available
- **Zero-Config Defaults**: Projects work without explicit MCP configuration
- **Dependency Transparency**: All MCP requirements visible in `package.json`

#### 6. Security & Isolation Model

**Dependency Isolation:**
- Each project maintains separate MCP server versions
- No cross-project MCP server interference
- Local installations prevent version conflicts

**Path Security:**
- All file paths resolved relative to project root
- `$PWD` substitution ensures proper sandboxing
- No access to parent directories without explicit configuration

## Security Considerations

- **No Network Dependencies** - Works offline after initial setup
- **Local Execution** - All processing happens locally
- **Git Security** - Leverages git's security model
- **No Telemetry** - No data collection or phone-home

## Performance Characteristics

- **Startup Time** - < 100ms typical
- **Resource Application** - < 1s for most operations
- **Memory Usage** - < 50MB typical
- **Binary Size** - ~15MB compressed

## Related Documentation

- [[architecture/cli-architecture]] - Detailed CLI architecture
- [[architecture/decisions/]] - Architecture decision records
- [[implementation/setup/installation]] - Installation guide
- [[development/contributing/architecture]] - Contributing to architecture