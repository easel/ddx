# Project Principles Generation Prompt

Generate project-specific principles based on the HELIX workflow principles and the SDD methodology.

## Storage Location

Store the principles document at: `docs/helix/01-frame/principles.md`

This central location ensures the principles are easily discoverable and remain the single source of truth for product principles.

## Context
You are establishing the foundational principles for a new project. These principles will govern all development decisions and must be followed throughout the project lifecycle.

## Instructions

1. **Review the core articles** from the template and ensure you understand their importance
2. **Add project-specific principles** based on:
   - The project's domain and requirements
   - Team capabilities and constraints
   - Performance and security needs
   - Regulatory or compliance requirements

3. **Document technology constraints** that will act as principles:
   - Choose primary language based on project needs
   - Select minimal framework set
   - Define data storage approach

4. **Be specific and enforceable**:
   - Each principle must be testable
   - Avoid vague statements
   - Include clear success/failure criteria

## Enforcement Reminder

These principles will be:
- Checked at every phase gate
- Referenced in all templates
- Validated by automated tools
- Used to guide AI assistance

## Questions to Consider

1. What makes this project unique?
2. What are the non-negotiable quality requirements?
3. What mistakes must be avoided?
4. What patterns should be encouraged?
5. What complexity is actually necessary?

Generate principles that will lead to maintainable, testable, and simple solutions.
