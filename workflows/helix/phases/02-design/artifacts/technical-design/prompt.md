# Prompt: Create Technical Design for User Story

You are creating a technical design document for a specific user story. This design should detail HOW to implement the story within the context of the broader solution architecture.

## Context

You have access to:
1. **User Story** (`US-XXX`): Defines WHAT needs to be built and acceptance criteria
2. **Solution Design** (`SD-XXX`): Overall technical architecture for the feature
3. **Feature Specification** (`FEAT-XXX`): Complete feature requirements

## Your Task

Create a technical design document that:

### 1. Maps Story to Implementation
- Review the user story's acceptance criteria
- Identify which components need changes
- Define new components if needed
- Specify the implementation approach

### 2. Details Technical Changes
- **API Changes**: New or modified endpoints
- **Data Model**: Schema changes or new entities
- **Business Logic**: How the story's requirements will be implemented
- **Integration**: How components will interact

### 3. Addresses Quality Attributes
- **Security**: Authentication, authorization, data protection
- **Performance**: Expected load and optimization strategies
- **Testing**: What tests are needed to verify the story
- **Compatibility**: Backward compatibility and migration

### 4. Provides Implementation Guidance
- **Sequence**: Order of development steps
- **Dependencies**: What must be done first
- **Risks**: Technical challenges and mitigations

## Key Principles

1. **Story-Focused**: This design is for ONE user story, not the entire feature
2. **Concrete**: Provide specific technical details, not abstractions
3. **Traceable**: Every acceptance criterion must map to technical implementation
4. **Testable**: Design must enable comprehensive testing
5. **Incremental**: Should be implementable independently as a vertical slice

## Questions to Answer

Before creating the technical design, ensure you can answer:

1. **Acceptance**: How will each acceptance criterion be satisfied technically?
2. **Components**: Which existing components need modification?
3. **Interfaces**: What APIs or contracts need to be created/modified?
4. **Data**: What data structures or schema changes are required?
5. **Security**: What security controls are needed for this story?
6. **Performance**: What are the performance implications?
7. **Testing**: How will we verify this story works correctly?
8. **Dependencies**: What other stories or components does this depend on?
9. **Rollback**: How can this story be safely rolled back if needed?

## Structure Guidelines

### DO:
- Reference the parent feature's solution design
- Map each acceptance criterion to technical implementation
- Provide concrete API specifications
- Include actual schema definitions
- Specify exact component changes
- Define clear test scenarios
- Consider security from the start

### DON'T:
- Duplicate the overall solution design
- Include implementation for other stories
- Leave acceptance criteria unmapped
- Provide vague or abstract descriptions
- Ignore non-functional requirements
- Forget about rollback scenarios

## Example Acceptance Criteria Mapping

```
Acceptance Criteria:
"Given a user is logged in, when they click 'List MCP Servers',
then they see all available servers with installation status"

Technical Implementation:
- Endpoint: GET /api/v1/mcp/servers
- Component: MCPController.listServers()
- Data: Query mcp_servers table, join with installations
- Response: JSON array with server details and installed flag
- Auth: Requires valid session token
- Cache: 5-minute TTL on server list
```

## Output

Your technical design should be:
- **Complete**: All acceptance criteria addressed
- **Specific**: Concrete technical details provided
- **Actionable**: Developer can start implementing immediately
- **Verifiable**: Clear success criteria and test approach
- **Safe**: Security and rollback considered

Remember: This design enables a developer to implement exactly what's needed for this user story - no more, no less.