# API Contract: Synchronization Commands [FEAT-002]

**Contract ID**: CLI-002
**Feature**: FEAT-002
**Type**: CLI
**Status**: Approved
**Version**: 1.0.0

*Command-line interface contracts for DDx synchronization commands*

## Command: update

**Purpose**: Pull latest changes from the master DDx repository while preserving local modifications
**Usage**: `ddx update [options]`

### Options
- `--strategy <strategy>`: Update strategy: "merge" (default), "rebase", "theirs", "ours"
- `--preview`: Show what would be updated without making changes
- `--force, -f`: Force update even with uncommitted changes
- `--branch <name>`: Update from specific branch (overrides config)
- `--repo <url>`: Update from specific repository (overrides config)
- `--components <list>`: Update only specific components (comma-separated: templates,patterns,prompts,configs)
- `--backup`: Create backup before updating (default: true)
- `--no-backup`: Skip backup creation

### Workflow
1. Check for uncommitted changes
2. Create backup (unless --no-backup)
3. Fetch latest from upstream
4. Apply update strategy
5. Resolve conflicts if needed
6. Update .ddx.yml metadata

### Output Format
```
Checking for updates from https://github.com/ddx-tools/ddx (main)...
Found 23 updates:
  Templates: 3 updated, 1 new
  Patterns: 5 updated, 2 new
  Prompts: 8 updated
  Configs: 4 updated

Local modifications detected in:
  - templates/nextjs/package.json (will be preserved)
  - patterns/auth/custom-auth.js (will be preserved)

Creating backup at .ddx-backup-20250115-143022...
Applying updates with merge strategy...
  ✓ Updated templates/python-flask
  ✓ Updated templates/go-cli
  ✓ Added templates/rust-wasm
  ⚠ Conflict in patterns/auth/jwt.js (local changes preserved)
  ✓ Updated prompts/claude/refactor

Update complete:
  ✓ 22 successful updates
  ⚠ 1 conflict (manual review required)

Review conflicts in:
  - patterns/auth/jwt.js.conflict

Run 'ddx doctor' to verify project health.
```

### Exit Codes
- `0`: Success (all updates applied)
- `1`: Partial success (some conflicts)
- `2`: No updates available
- `3`: Uncommitted changes (use --force)
- `4`: Network error
- `5`: Invalid strategy
- `6`: Backup failed

### Update Strategies
- **merge**: Three-way merge, preserving local changes (default)
- **rebase**: Reapply local changes on top of upstream
- **theirs**: Accept all upstream changes, discard local
- **ours**: Keep all local changes, ignore upstream

### Examples
```bash
# Standard update
$ ddx update
Checking for updates...
Found 5 updates, applying...

# Preview updates
$ ddx update --preview
Would update:
  - templates/nextjs (v1.2.0 → v1.3.0)
  - patterns/auth/jwt (security patch)
  - prompts/claude/test (improved prompts)
No changes made (preview mode)

# Force update with uncommitted changes
$ ddx update --force
Warning: Uncommitted changes in working directory
Creating backup...
Proceeding with update...

# Update specific components only
$ ddx update --components templates,patterns
Updating only: templates, patterns
Skipping: prompts, configs

# Update from different branch
$ ddx update --branch experimental
Updating from experimental branch...

# Handle conflicts
$ ddx update
Conflict detected in patterns/auth/jwt.js
Options:
  1. Keep local version
  2. Accept upstream version
  3. Open merge tool
  4. Skip this file
Choice [1-4]: 3
Opening merge tool...
```

### Conflict Resolution
When conflicts occur, DDx creates:
- `.conflict` files with upstream version
- `.backup` files with original version
- Marks conflicts in files with git-style markers

```javascript
// patterns/auth/jwt.js (with conflict)
<<<<<<<< LOCAL (without extra <)
const TOKEN_EXPIRY = '1h'; // Custom local change
======== (without extra =)
const TOKEN_EXPIRY = '30m'; // Upstream security update
>>>>>>>> UPSTREAM (without extra >)
```

---

## Command: contribute

**Purpose**: Share local improvements back to the DDx community repository
**Usage**: `ddx contribute [resource-path] [options]`

### Arguments
- `resource-path` (optional): Specific resource to contribute (default: all local changes)

### Options
- `--message, -m <msg>`: Contribution message/description
- `--type <type>`: Contribution type: "feature", "fix", "improvement", "pattern", "template"
- `--preview`: Preview what would be contributed
- `--branch <name>`: Target branch for contribution (default: contributions)
- `--author <name>`: Author attribution
- `--email <email>`: Author email
- `--validate`: Run validation before contributing
- `--pr`: Create pull request automatically (requires GitHub CLI)
- `--issue <number>`: Link to related issue

### Workflow
1. Identify local changes vs upstream
2. Validate changes meet contribution guidelines
3. Prepare contribution package
4. Push to contribution branch
5. Optionally create pull request

### Output Format
```
Analyzing local changes...
Found contributions:
  ✓ New template: templates/vue-composition
  ✓ Modified pattern: patterns/auth/oauth2
  ✓ New prompt: prompts/claude/security-review

Validating contributions...
  ✓ Template structure valid
  ✓ Pattern tests passing
  ✓ Documentation complete
  ✓ No sensitive data detected

Contribution Summary:
  Type: feature
  Message: Add Vue 3 Composition API template with TypeScript
  Author: Jane Developer <jane@example.com>

Files included (12):
  - templates/vue-composition/package.json
  - templates/vue-composition/tsconfig.json
  - templates/vue-composition/src/...
  - docs/templates/vue-composition.md

Proceed with contribution? (y/n): y

Preparing contribution...
  ✓ Created contribution branch: contributions/vue-composition-20250115
  ✓ Committed changes
  ✓ Pushed to upstream repository

Creating pull request...
  ✓ PR #234 created: https://github.com/ddx-tools/ddx/pull/234

Thank you for contributing to DDx!
Your contribution will be reviewed by maintainers.
Track status at: https://github.com/ddx-tools/ddx/pull/234
```

### Validation Checks
Before accepting contributions, DDx validates:
1. **Structure**: Follows DDx resource structure
2. **Documentation**: Has README and examples
3. **Testing**: Includes tests or test data
4. **Security**: No hardcoded secrets or credentials
5. **Licensing**: Compatible license
6. **Quality**: Passes linting and formatting

### Exit Codes
- `0`: Success
- `1`: No changes to contribute
- `2`: Validation failed
- `3`: Authentication required
- `4`: Network error
- `5`: Permission denied
- `6`: Contribution rejected by hooks

### Examples
```bash
# Contribute all local changes
$ ddx contribute -m "Add OAuth2 authentication pattern"
Found 3 local changes to contribute...

# Contribute specific resource
$ ddx contribute templates/vue-composition -m "Vue 3 Composition API template" --type template
Contributing: templates/vue-composition

# Preview contribution
$ ddx contribute --preview
Would contribute:
  - templates/custom-react (new)
  - patterns/cache/redis (modified)
  - Total: 15 files, 2.3 KB
No changes made (preview mode)

# Create PR automatically
$ ddx contribute --pr --issue 123 -m "Fix: JWT expiration handling"
Creating contribution with PR...
✓ PR #235 created and linked to issue #123

# Validation failure
$ ddx contribute patterns/insecure
Validating contribution...
✗ Security check failed:
  - Hardcoded API key found in line 45
  - Missing input validation in auth.js
Please fix these issues before contributing.

# With custom author
$ ddx contribute --author "Team Alpha" --email team@company.com
Contributing as: Team Alpha <team@company.com>
```

### Contribution Guidelines
Contributions must follow these guidelines:
1. **Naming**: Use kebab-case for directories, appropriate case for files
2. **Documentation**: Include README.md with usage examples
3. **Variables**: Use `{{variable}}` syntax for substitution
4. **Dependencies**: List all required dependencies
5. **Testing**: Include test files or test instructions
6. **Examples**: Provide working examples

### Contribution Package Structure
```
contribution/
├── metadata.json          # Contribution metadata
├── resources/
│   ├── templates/
│   │   └── vue-composition/
│   ├── patterns/
│   │   └── auth/oauth2/
│   └── prompts/
│       └── claude/security-review
├── docs/
│   └── contributions/
│       └── vue-composition.md
└── tests/
    └── vue-composition.test.js
```

### Pull Request Template
When creating PRs, DDx uses this template:
```markdown
## Contribution: [Title]

### Type
- [ ] New Template
- [ ] New Pattern
- [ ] New Prompt
- [ ] Improvement
- [ ] Bug Fix

### Description
[Detailed description of the contribution]

### Testing
- [ ] Tests included
- [ ] Documentation updated
- [ ] Examples provided

### Checklist
- [ ] Follows DDx guidelines
- [ ] No sensitive data
- [ ] Validated locally

Fixes #[issue-number]
```

## Synchronization State Management

Both commands maintain synchronization state in `.ddx.yml`:
```yaml
sync:
  last_update: 2025-01-15T14:30:00Z
  upstream_commit: abc123def456
  local_modifications:
    - templates/nextjs/package.json
    - patterns/auth/custom-auth.js
  contribution_history:
    - date: 2025-01-10
      pr: 234
      status: merged
```

## Conflict Management

DDx provides several conflict resolution strategies:
1. **Automatic**: Merge non-conflicting changes
2. **Interactive**: Prompt for each conflict
3. **Manual**: Create .conflict files for review
4. **Policy-based**: Apply predetermined rules

## Network and Authentication

- Uses git credentials for authentication
- Supports SSH and HTTPS protocols
- Handles proxy configurations
- Implements retry logic for network failures
- Caches authentication tokens appropriately

## Error Handling
```bash
# Network failure
$ ddx update
Error: Failed to connect to repository
Possible causes:
  - Network connection issue
  - Repository URL incorrect
  - Authentication required
Try: ddx update --repo [url] --debug

# Merge conflict
$ ddx update
Error: Automatic merge failed for patterns/auth/jwt.js
Manual resolution required:
  1. Review changes in patterns/auth/jwt.js.conflict
  2. Edit patterns/auth/jwt.js to resolve
  3. Remove .conflict file when resolved
  4. Run 'ddx update --continue'

# Contribution rejection
$ ddx contribute
Error: Contribution validation failed
  - Missing required documentation
  - Test coverage below threshold (60% required)
Run 'ddx contribute --validate' for detailed report
```