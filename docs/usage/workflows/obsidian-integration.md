---
tags: [workflows, obsidian, integration, knowledge-management, tools]
aliases: ["Obsidian Integration", "Obsidian Setup", "Knowledge Management"]
created: 2025-01-12
modified: 2025-01-12
---

# Obsidian Integration for DDX Workflows

This guide shows how to set up Obsidian as your knowledge management system for DDX workflows, creating a powerful environment for documentation, tracking, and collaboration.

## Why Obsidian + DDX?

Obsidian's powerful linking, graph view, and plugin ecosystem make it an ideal companion for DDX workflows:

- **Bidirectional Linking**: Connect workflows, artifacts, and decisions
- **Graph Visualization**: See relationships between workflow components
- **Template System**: Integrate with DDX templates seamlessly
- **Plugin Ecosystem**: Extend functionality for project management
- **Version Control**: Git-friendly markdown files
- **Collaborative**: Share vaults with team members

Following DDX's medical metaphor, Obsidian becomes your "Electronic Health Record" system - a comprehensive, interconnected view of your project's development history and current state.

## Initial Setup

### 1. Install Obsidian

Download from [obsidian.md](https://obsidian.md) and install for your platform.

### 2. Create DDX Vault

Create a new vault specifically for your DDX workflows:

```bash
# Create vault directory structure
mkdir my-project-vault
cd my-project-vault

# Initialize git repository (optional but recommended)
git init
echo ".obsidian/workspace*" >> .gitignore

# Create initial folder structure
mkdir -p {workflows,artifacts,decisions,meetings,references}
```

### 3. Configure Vault Structure

Set up a folder structure that aligns with DDX workflows:

```
my-project-vault/
├── .obsidian/               # Obsidian configuration
├── workflows/              # DDX workflow instances
│   ├── current/           # Active workflows
│   ├── completed/         # Finished workflows
│   └── templates/         # Workflow templates
├── artifacts/             # Generated documents
│   ├── requirements/      # PRDs, specs, etc.
│   ├── design/           # Architecture, designs
│   ├── implementation/   # Code documentation
│   └── testing/          # Test plans, results
├── decisions/             # Architecture Decision Records
├── meetings/              # Meeting notes
├── references/            # External references
├── people/               # Team member profiles
└── projects/             # Project overviews
```

## Essential Plugins

### Core Workflow Plugins

#### 1. Templater
Advanced template system for DDX workflow integration.

**Installation:**
1. Settings → Community Plugins → Browse
2. Search "Templater" → Install → Enable

**Configuration:**
```javascript
// Templates/ddx-workflow-start.md
---
tags: [workflow, <% tp.file.title.toLowerCase() %>]
created: <% tp.date.now() %>
status: active
workflow_type: {{VALUE:development|debugging|optimization}}
---

# <% tp.file.title %>

## Workflow Overview
- **Type**: {{VALUE:development|debugging|optimization}}
- **Started**: <% tp.date.now() %>
- **Team**: [[People/{{VALUE:team-lead}}]]
- **Project**: [[Projects/{{VALUE:project-name}}]]

## Current Phase
- [ ] [[Workflows/Templates/phase-define|Define Phase]]
- [ ] [[Workflows/Templates/phase-design|Design Phase]]  
- [ ] [[Workflows/Templates/phase-implement|Implement Phase]]
- [ ] [[Workflows/Templates/phase-test|Test Phase]]
- [ ] [[Workflows/Templates/phase-release|Release Phase]]

## Artifacts
*Links to generated artifacts will appear here*

## Notes
*Workflow-specific notes and observations*
```

#### 2. Dataview
Query and display workflow data dynamically.

**Installation:** Same process as Templater

**Example Queries:**
```dataview
// Active workflows dashboard
TABLE workflow_type, started, status
FROM "workflows/current"
WHERE status = "active"
SORT started DESC
```

```dataview
// Artifacts by phase
LIST
FROM "artifacts"
WHERE contains(tags, "current-workflow")
SORT file.ctime DESC
```

#### 3. Kanban
Visual workflow management boards.

**Installation:** Enable Kanban plugin

**Workflow Board Template:**
```markdown
---
kanban-plugin: board
---

## Backlog
- [ ] [[Define Phase Tasks]]
- [ ] [[Design Phase Tasks]]

## In Progress  
- [ ] [[Current Phase Work]]

## Review
- [ ] [[Pending Reviews]]

## Done
- [x] [[Completed Tasks]]
```

#### 4. Calendar
Track workflow milestones and deadlines.

**Usage:**
- Create calendar events for phase deadlines
- Link to workflow documents
- Track team availability

### Advanced Plugins

#### 1. Advanced Tables
Better table editing for requirements matrices.

#### 2. Excalidraw
Embedded diagrams for architecture and workflows.

#### 3. Tasks
Advanced task management with due dates and dependencies.

#### 4. Git
Version control integration for team collaboration.

## Template System

### DDX Template Integration

Create Obsidian templates that mirror DDX workflow templates:

#### Workflow Starter Template
```markdown
# Templates/workflow-starter.md
---
tags: [workflow, template]
---

# {{title}} Workflow

## Metadata
- **Workflow ID**: {{title | lower | replace(" ", "-")}}
- **Type**: {{type}}
- **Started**: {{date}}
- **Owner**: [[People/{{owner}}]]
- **Project**: [[Projects/{{project}}]]

## Phase Checklist
- [ ] [[#Define Phase]] 
- [ ] [[#Design Phase]]
- [ ] [[#Implement Phase]]
- [ ] [[#Test Phase]]
- [ ] [[#Release Phase]]
- [ ] [[#Iterate Phase]]

## Define Phase
### Entry Criteria
- [ ] Problem identified
- [ ] Stakeholders aligned

### Activities
- [ ] Create [[Artifacts/PRD-{{title}}|Product Requirements Document]]
- [ ] Define [[Artifacts/Success-Metrics-{{title}}|Success Metrics]]

### Exit Criteria
- [ ] Requirements approved
- [ ] Success criteria defined

---

## Design Phase
*Continue for each phase...*
```

#### Artifact Templates

##### PRD Template
```markdown
# Templates/prd-template.md
---
tags: [template, prd, requirements]
---

# {{product-name}} - Product Requirements Document

## Problem Statement
{{problem-description}}

## Target Users
{{user-personas}}

## Solution Overview
{{solution-approach}}

## Success Metrics
{{success-criteria}}

## Dependencies
- [[Artifacts/{{dependency-1}}]]
- [[Artifacts/{{dependency-2}}]]

## Related
- **Workflow**: [[Workflows/{{workflow-name}}]]
- **Architecture**: [[Decisions/{{arch-decision}}]]
```

##### ADR Template
```markdown
# Templates/adr-template.md
---
tags: [template, adr, decision]
---

# ADR-{{number}}: {{title}}

**Status**: {{status}}
**Date**: {{date}}
**Deciders**: {{decision-makers}}

## Context
{{context-description}}

## Decision
{{decision-made}}

## Consequences
### Positive
- {{positive-consequence-1}}
- {{positive-consequence-2}}

### Negative  
- {{negative-consequence-1}}
- {{negative-consequence-2}}

## Related
- **Workflow**: [[Workflows/{{workflow-name}}]]
- **Artifacts**: [[Artifacts/{{related-artifact}}]]
```

### Template Hotkeys

Set up keyboard shortcuts for quick template insertion:

1. Settings → Hotkeys
2. Search "Templater: Insert Template"
3. Assign shortcuts:
   - `Ctrl+T, W` - Workflow starter
   - `Ctrl+T, P` - PRD template  
   - `Ctrl+T, A` - ADR template
   - `Ctrl+T, M` - Meeting notes

## Graph View Configuration

### Customizing the Graph

Configure the graph view to visualize workflow relationships:

#### Graph Settings
```json
{
  "collapse-color-groups": false,
  "colorGroups": [
    {
      "query": "tag:#workflow",
      "color": {
        "a": 1,
        "rgb": 5431378
      }
    },
    {
      "query": "tag:#artifact",
      "color": {
        "a": 1,
        "rgb": 5419488
      }
    },
    {
      "query": "tag:#decision",
      "color": {
        "a": 1,
        "rgb": 14725458
      }
    }
  ],
  "collapse-display": true,
  "showTags": true,
  "showAttachments": false,
  "hideUnresolved": false,
  "showOrphans": false
}
```

#### Workflow-Specific Views

Create filtered graph views for different contexts:

1. **Current Workflows**: Show only active workflow components
2. **Project History**: Show completed workflows and their artifacts
3. **Decision Tree**: Focus on ADRs and their relationships
4. **Team Knowledge**: Show people and their contributions

## Dataview Queries for Workflow Tracking

### Active Workflow Dashboard

```dataview
TABLE 
  workflow_type as "Type",
  started as "Started", 
  current_phase as "Phase",
  choice(status = "active", "🟢", choice(status = "blocked", "🔴", "🟡")) as "Status"
FROM "workflows/current"
WHERE status != "completed"
SORT started DESC
```

### Artifact Status Board

```dataview
TABLE
  artifact_type as "Type",
  status as "Status",
  owner as "Owner",
  due_date as "Due"
FROM "artifacts"  
WHERE contains(tags, "current")
SORT due_date ASC
```

### Decision History

```dataview
LIST
FROM "decisions"
WHERE status = "accepted"
SORT file.ctime DESC
LIMIT 10
```

### Team Contribution Analysis

```dataview
TABLE
  length(file.outlinks) as "Contributions",
  length(file.inlinks) as "References"
FROM "people"
SORT length(file.outlinks) DESC
```

### Phase Completion Tracking

```dataview
TASK
FROM "workflows/current"
WHERE !completed
GROUP BY file.name
SORT status
```

## Automation and Scripts

### Templater Scripts

#### Auto-link Creation
```javascript
// Templates/scripts/auto-link.js
function createWorkflowLinks(workflowName) {
  const phases = ['define', 'design', 'implement', 'test', 'release'];
  const artifacts = ['prd', 'architecture', 'test-plan'];
  
  let links = "## Quick Links\n";
  
  // Phase links
  links += "### Phases\n";
  phases.forEach(phase => {
    links += `- [[Workflows/${workflowName}/phases/${phase}|${phase} Phase]]\n`;
  });
  
  // Artifact links  
  links += "### Artifacts\n";
  artifacts.forEach(artifact => {
    links += `- [[Artifacts/${workflowName}-${artifact}|${artifact}]]\n`;
  });
  
  return links;
}
```

#### Workflow Status Updates
```javascript
// Templates/scripts/status-update.js
function updateWorkflowStatus() {
  const today = tp.date.now();
  const currentFile = tp.file.title;
  
  return `
## Status Update - ${today}
**Workflow**: ${currentFile}
**Updated By**: {{team-member}}
**Current Phase**: {{current-phase}}

### Progress This Week
- {{accomplishment-1}}
- {{accomplishment-2}}

### Next Week Plan
- {{plan-1}}
- {{plan-2}}

### Blockers
- {{blocker-1}}

### Metrics
- **Phase Completion**: {{completion-percentage}}%
- **Artifacts Created**: {{artifact-count}}
- **Decisions Made**: {{decision-count}}
`;
}
```

### QuickAdd Macros

Set up QuickAdd plugin for rapid workflow operations:

#### New Workflow Macro
```javascript
// Create new workflow with full structure
module.exports = async (params) => {
  const workflowName = await tp.system.prompt("Workflow name:");
  const workflowType = await tp.system.suggester(
    ["Development", "Debugging", "Optimization"],
    ["development", "debugging", "optimization"]
  );
  
  // Create workflow directory structure
  await app.vault.createFolder(`workflows/current/${workflowName}`);
  await app.vault.createFolder(`workflows/current/${workflowName}/phases`);
  
  // Create main workflow file
  const workflowContent = await tp.file.include("[[Templates/workflow-starter]]");
  await app.vault.create(
    `workflows/current/${workflowName}/${workflowName}.md`,
    workflowContent
  );
};
```

## Collaboration Features

### Team Vault Setup

For team collaboration:

#### 1. Shared Vault Structure
```
team-vault/
├── .obsidian/
│   ├── config.json          # Shared settings
│   ├── hotkeys.json         # Team hotkeys  
│   └── plugins/            # Team plugins
├── workflows/
│   ├── active/             # Current team workflows
│   └── archive/            # Completed workflows
├── people/
│   ├── alice.md           # Team member profiles
│   ├── bob.md
│   └── carol.md
└── team/
    ├── standards.md        # Team conventions
    ├── tools.md           # Tool configurations
    └── onboarding.md      # New member guide
```

#### 2. Git Workflow
```bash
# Daily workflow sync
git pull origin main
# Work on workflows
git add .
git commit -m "Update workflow progress"
git push origin main
```

#### 3. Conflict Resolution
Use Obsidian's merge conflict resolution:
- Configure external merge tool
- Resolve conflicts in familiar editor
- Maintain workflow continuity

### Real-time Collaboration

For live collaboration:
- **Obsidian Sync**: Official sync service
- **Self-hosted sync**: Using Syncthing or similar  
- **Git-based**: Regular commits and pulls
- **Shared drives**: OneDrive, Google Drive (with caution)

## Mobile Integration

### Obsidian Mobile Setup

Sync workflows to mobile devices:

#### 1. Quick Capture Templates
```markdown
# Templates/mobile-capture.md
---
tags: [mobile, capture, {{date:YYYY-MM-DD}}]
---

# Quick Capture - {{date}} {{time}}

## Context
*Where/what/why*

## Idea/Issue/Task
*Main content*

## Related
- **Workflow**: [[Workflows/{{workflow}}]]
- **Person**: [[People/{{person}}]]

## Actions
- [ ] Process this capture
- [ ] Add to appropriate workflow
```

#### 2. Status Check Dashboard
```markdown
# Mobile/dashboard.md
# Team Dashboard

```dataview
TABLE status, current_phase, owner
FROM "workflows/active"
LIMIT 5
```

## Recent Updates
```dataview
LIST
FROM ""
WHERE file.mtime > date(today) - dur(3 days)
SORT file.mtime DESC
LIMIT 10
```
```

## Workflow Validation

### Automated Checks

Use Dataview to validate workflow completeness:

#### Missing Artifacts Check
```dataview
TABLE
  workflow_type,
  status,
  "❌ Missing artifacts" as Issues
FROM "workflows/active"
WHERE !contains(file.outlinks, "artifacts")
```

#### Incomplete Phases Check
```dataview
TASK
FROM "workflows/active"
WHERE contains(tags, "incomplete-phase")
```

#### Broken Links Check
```dataview
LIST
FROM "workflows"
WHERE length(file.outlinks) > length(file.inlinks) * 2
```

### Quality Metrics

Track workflow quality with custom queries:

#### Documentation Coverage
```dataview
TABLE
  (length(file.outlinks) / 10) * 100 as "Documentation %",
  length(file.inlinks) as "References",
  file.size as "Size"
FROM "workflows/active"
SORT file.size DESC
```

## Best Practices

### Organization

1. **Consistent Naming**: Use kebab-case for files and folders
2. **Clear Tagging**: Develop team tagging conventions
3. **Regular Cleanup**: Archive completed workflows
4. **Cross-linking**: Connect related concepts liberally

### Performance

1. **Plugin Management**: Only enable needed plugins
2. **Large Vaults**: Use folder exclusion for performance
3. **Image Optimization**: Compress embedded images
4. **Regular Maintenance**: Clean up orphaned files

### Security

1. **Sensitive Data**: Use separate vault for sensitive info
2. **Access Control**: Manage team vault permissions
3. **Backup Strategy**: Regular backups of vault data
4. **Version History**: Leverage git for change tracking

## Troubleshooting

### Common Issues

#### Slow Performance
- Disable unnecessary plugins
- Exclude large folders from indexing
- Optimize Dataview queries

#### Sync Conflicts
- Establish team merge conventions
- Use descriptive commit messages
- Communicate about major changes

#### Plugin Compatibility
- Test plugins with team vault
- Document required plugin versions
- Have fallback plans for critical plugins

### Getting Help

#### Community Resources
- Obsidian Community Discord
- r/ObsidianMD subreddit
- DDX community forums

#### Documentation
- Official Obsidian docs
- Plugin-specific documentation  
- Team-specific conventions

## Advanced Workflows

### Multi-Project Management

For managing multiple projects:

#### Project Index
```markdown
# Projects/index.md

```dataview
TABLE
  status,
  current_phase,
  team_lead,
  last_update
FROM "projects"
WHERE status != "archived"
SORT last_update DESC
```
```

#### Cross-Project Dependencies
```dataview
TABLE
  project,
  dependency,
  status,
  impact
FROM "dependencies"
WHERE status = "active"
```

### Metrics and Reporting

#### Weekly Team Report
```markdown
# Reports/weekly-{{date:YYYY-MM-DD}}.md

## Team Velocity
```dataview
TABLE
  completed_this_week,
  planned_next_week,
  velocity_trend
FROM "workflows/active"
```

## Workflow Health
```dataview
TABLE
  choose(days_in_phase < 7, "🟢", choose(days_in_phase < 14, "🟡", "🔴")) as Health,
  current_phase,
  days_in_phase
FROM "workflows/active"  
```
```

### Integration with External Tools

#### Jira Integration
```markdown
# External ticket: PROJ-123
# Obsidian link: [[Workflows/feature-x]]
```

#### GitHub Integration
```markdown
# PR: https://github.com/team/repo/pull/456
# Related workflow: [[Workflows/bug-fix-y]]
```

## Medical Metaphor in Obsidian

Applying DDX's medical metaphor to Obsidian:

- **Vault** = Hospital Information System
- **Workflows** = Treatment Protocols  
- **Graph View** = Patient Relationship Map
- **Templates** = Medical Forms
- **Dataview Queries** = Diagnostic Reports
- **Tags** = Medical Classifications
- **Links** = Care Continuity
- **Backlinks** = Medical History

## Getting Started Checklist

### Initial Setup
- [ ] Install Obsidian
- [ ] Create DDX workflow vault
- [ ] Install essential plugins (Templater, Dataview, Kanban)
- [ ] Set up folder structure
- [ ] Configure graph view

### Template Creation  
- [ ] Create workflow starter template
- [ ] Create artifact templates (PRD, ADR, etc.)
- [ ] Set up template hotkeys
- [ ] Test template functionality

### Team Integration
- [ ] Share vault with team
- [ ] Document team conventions
- [ ] Set up collaboration workflow
- [ ] Train team on Obsidian basics

### Workflow Integration
- [ ] Import existing DDX workflows
- [ ] Create dashboard queries
- [ ] Set up automation scripts
- [ ] Test end-to-end workflow

### Ongoing Maintenance
- [ ] Regular vault cleanup
- [ ] Plugin updates
- [ ] Template improvements  
- [ ] Team feedback incorporation

## Next Steps

1. **Start Simple**: Begin with basic templates and folder structure
2. **Iterate Gradually**: Add complexity as team comfort grows
3. **Customize Heavily**: Adapt to your specific workflow needs
4. **Share Knowledge**: Document and share successful patterns
5. **Stay Updated**: Keep plugins and practices current

Remember: Obsidian is a powerful tool, but its value comes from consistent use and gradual improvement. Start with basic DDX workflow integration and evolve your setup over time based on real usage patterns and team feedback.