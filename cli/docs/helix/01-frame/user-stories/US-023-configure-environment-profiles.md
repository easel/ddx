# US-023: Configure Environment Profiles

## User Story

**As a** developer working with multiple environments
**I want** to easily create, manage, and switch between environment profiles
**So that** I can efficiently manage different configurations for development, staging, and production environments

## Background

Organizations typically deploy applications across multiple environments (development, staging, production, etc.), each requiring different configurations. While DDx supports environment-specific configuration files through DDX_ENV (implemented in US-019), users need a convenient way to create, manage, and switch between these profiles without manually managing configuration files.

## Acceptance Criteria

### AC1: Create Environment Profiles
- **Given** I want to set up a new environment profile
- **When** I run `ddx config profile create <profile-name>`
- **Then** a new `.ddx.<profile-name>.yml` configuration file should be created
- **And** the profile should inherit base configuration settings
- **And** I should be able to customize environment-specific settings

### AC2: List Available Profiles
- **Given** I have multiple environment profiles configured
- **When** I run `ddx config profile list`
- **Then** I should see all available profiles with their status
- **And** the currently active profile should be clearly indicated
- **And** profile details should include last modified date and validation status

### AC3: Activate Environment Profiles
- **Given** I have multiple environment profiles
- **When** I run `ddx config profile activate <profile-name>`
- **Then** the DDX_ENV environment variable should be set for the current session
- **And** subsequent DDx commands should use the activated profile's configuration
- **And** I should receive confirmation of the profile activation

### AC4: Copy Environment Profiles
- **Given** I want to create a new profile based on an existing one
- **When** I run `ddx config profile copy <source-profile> <destination-profile>`
- **Then** a new profile should be created with the source profile's configuration
- **And** I should be able to modify the copied profile independently
- **And** the copy operation should preserve all configuration sections

### AC5: Validate Environment Profiles
- **Given** I want to verify a specific profile's configuration
- **When** I run `ddx config profile validate <profile-name>`
- **Then** the profile's configuration should be validated for syntax and completeness
- **And** I should receive detailed validation results
- **And** validation errors should include specific line numbers and suggestions

### AC6: Show Profile Configuration
- **Given** I want to review a profile's configuration
- **When** I run `ddx config profile show <profile-name>`
- **Then** I should see the complete resolved configuration for that profile
- **And** inherited values should be clearly distinguished from profile-specific values
- **And** the output should be formatted for easy reading

### AC7: Compare Environment Profiles
- **Given** I want to understand differences between profiles
- **When** I run `ddx config profile diff <profile-a> <profile-b>`
- **Then** I should see a detailed comparison of the two profiles
- **And** differences should be highlighted with clear indicators
- **And** the comparison should include added, removed, and modified values

### AC8: Delete Environment Profiles
- **Given** I want to remove an obsolete environment profile
- **When** I run `ddx config profile delete <profile-name>`
- **Then** the profile configuration file should be safely removed
- **And** I should receive confirmation before deletion
- **And** the operation should prevent deletion of the currently active profile

## Business Value

### Primary Benefits
- **Streamlined Environment Management**: Developers can quickly switch between environments
- **Reduced Configuration Errors**: Profile validation prevents deployment issues
- **Improved Development Workflow**: Easy profile copying accelerates environment setup
- **Enhanced Visibility**: Profile listing and comparison improve configuration transparency

### Use Cases
- **Multi-Environment Development**: Different database connections, API endpoints, and feature flags per environment
- **Team Collaboration**: Shared profile templates for consistent environment setups
- **CI/CD Integration**: Profile activation in deployment pipelines
- **Configuration Debugging**: Profile comparison and validation for troubleshooting

## Technical Considerations

### Integration Points
- Builds on existing environment configuration system (US-019)
- Leverages configuration validation framework (US-022)
- Integrates with variable substitution system (US-018)
- Uses existing configuration loading and merging logic

### Implementation Notes
- Profile commands should be subcommands of `ddx config profile`
- Profile activation should set DDX_ENV for current shell session
- Profile validation should use enhanced validation from US-022
- Profile comparison should support multiple output formats (text, JSON)

## Definition of Done

- [ ] All profile management commands are implemented and functional
- [ ] Profile creation includes proper inheritance from base configuration
- [ ] Profile activation successfully sets environment variables
- [ ] Profile validation provides comprehensive error reporting
- [ ] Profile comparison shows clear, actionable differences
- [ ] Unit tests cover all profile management scenarios
- [ ] Integration tests verify profile switching and inheritance
- [ ] Documentation explains profile concepts and command usage
- [ ] All acceptance criteria are met and verified
- [ ] Performance is acceptable for typical profile operations

## Dependencies

- **US-017**: Initialize Configuration (foundation configuration system)
- **US-019**: Override Configuration (environment-specific configuration files)
- **US-022**: Validate Configuration (enhanced validation framework)

## Estimated Effort

**Medium** - Builds on solid existing foundation, primarily adding user interface commands

## Priority

**High** - Essential for multi-environment workflows and developer productivity