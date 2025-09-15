# ADR-011: Adopting "Persona" Terminology for AI Personality Templates

**ADR ID**: ADR-011
**Status**: Accepted
**Date**: 2025-01-15
**Decision Makers**: DDX Core Team

## Context

The DDX project needs a system for defining and managing reusable AI personality templates that can be applied to different roles within development workflows. We need to establish consistent terminology that:

1. Clearly communicates the purpose of these templates
2. Avoids confusion with existing concepts
3. Aligns with industry practices
4. Supports our abstraction model (concrete implementations vs abstract roles)

### Options Considered

1. **"Agents"** - Used by Claude Code, Microsoft AutoGen, CrewAI
2. **"Personas"** - Used in prompt engineering, GitHub Copilot community
3. **"Assistants"** - Used by OpenAI, general AI tools
4. **"Personalities"** - Descriptive but not standard
5. **"Roles"** - Could be confused with abstract function

### Industry Research

Our research revealed:
- Claude Code uses "agents" in `.claude/agents/` directory
- OpenAI Codex community uses "AGENTS.md" files
- GitHub Copilot community increasingly uses "persona-based development"
- Prompt engineering literature favors "persona" for role-based prompting
- "Agent" increasingly implies autonomous operation, not just personality

## Decision

We will adopt **"persona"** as the primary term for AI personality templates, with **"role"** reserved for abstract functions.

### Terminology Definitions

- **Persona**: A concrete AI personality implementation (e.g., `strict-code-reviewer`)
- **Role**: An abstract function that can be fulfilled by personas (e.g., `code-reviewer`)
- **Binding**: The connection between a role and a specific persona in a project

### Rationale

1. **Clarity of Purpose**: "Persona" immediately conveys that these are personality templates, not autonomous agents

2. **Avoid Agent Confusion**: "Agent" increasingly implies autonomous AI systems (AutoGPT, AutoGen) that can take independent action. Our templates are passive personality definitions.

3. **Industry Alignment**: The growing "persona-based development" movement, especially in the GitHub Copilot community, uses this terminology.

4. **Clean Abstraction**: Creates clear separation:
   - Workflows require **roles** (abstract)
   - Projects bind **personas** to roles (concrete)
   - Users load **personas** for sessions (application)

5. **Intuitive Understanding**: Developers immediately understand that a "persona" is a defined personality/approach, similar to user personas in UX design.

## Consequences

### Positive

- Clear, intuitive terminology that accurately describes the feature
- Avoids confusion with autonomous agent systems
- Aligns with emerging industry practices
- Supports clean abstraction between requirements and implementation
- Natural language: "This persona fulfills the code-reviewer role"

### Negative

- Differs from Claude Code's "agents" terminology (minor inconsistency)
- May need to educate users familiar with "agent" terminology
- Some tools use "agent" for similar concepts

### Neutral

- File structure uses `/personas/` instead of `/agents/`
- CLI commands use `ddx persona` instead of `ddx agent`
- Configuration uses `persona_bindings` instead of `agent_mappings`
- Documentation must be consistent in terminology usage

## Implementation

### File Structure
```
/personas/                    # Not /agents/
├── strict-code-reviewer.md
├── test-engineer-tdd.md
└── ...
```

### Configuration
```yaml
persona_bindings:             # Not agent_bindings
  code-reviewer: strict-code-reviewer
  test-engineer: test-engineer-tdd
```

### CLI Commands
```bash
ddx persona list              # Not ddx agent list
ddx persona bind
ddx persona load
```

### Documentation
- Always use "persona" for the personality template
- Always use "role" for the abstract function
- Never mix "agent" terminology except when referring to external systems

## Alternatives Considered

### Alternative 1: Use "Agent" Throughout
- **Pros**: Aligns with Claude Code, familiar term
- **Cons**: Implies autonomy, confusing with agent frameworks
- **Rejected because**: Creates wrong mental model

### Alternative 2: Use "Personality"
- **Pros**: Very descriptive, unambiguous
- **Cons**: Not standard, longer to type
- **Rejected because**: No industry precedent

### Alternative 3: Dual Terminology
- **Pros**: Could use both terms contextually
- **Cons**: Confusing, inconsistent
- **Rejected because**: Clarity requires single term

## References

- Claude Code subagents documentation
- GitHub Copilot persona-based development articles
- Prompt engineering best practices (2024-2025)
- OpenAI AGENTS.md specification
- Microsoft AutoGen agent documentation

## Review and Approval

- **Proposed by**: DDX Core Team
- **Reviewed by**: Architecture Team
- **Approved by**: Project Lead
- **Approval Date**: 2025-01-15

---

*This decision establishes "persona" as the standard terminology for AI personality templates in DDX, creating a clear and intuitive vocabulary for the feature.*