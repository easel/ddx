# Specification Validator
# Validates that specifications exist before implementation according to CDP principles

load("//lib/common.star", "read_file", "file_exists", "is_code_file", "is_spec_file", 
     "find_related_files", "create_violation", "get_validator_config")
load("//lib/git.star", "get_all_changed_files")

def validate_spec_before_code(ctx):
    """Validate that specifications exist before code implementation.
    
    Enforces the CDP principle that all code changes must have corresponding
    specifications or design documents that describe the intended functionality.
    
    Args:
        ctx (dict): Validation context containing files, git info, and project config
        
    Returns:
        list: List of violation objects
    """
    violations = []
    config = get_validator_config(ctx, "spec_validator")
    
    if not config.get("enabled", True):
        return violations
    
    changed_files = get_all_changed_files(ctx)
    all_project_files = _get_all_project_files(ctx)
    
    # Check each changed code file for corresponding specifications
    for file_path in changed_files:
        if not is_code_file(file_path):
            continue
        
        # Skip test files - they don't need separate specs
        if file_path.endswith("_test.go") or file_path.endswith(".test.js") or "test_" in file_path:
            continue
        
        spec_violations = _validate_file_has_spec(file_path, all_project_files, config, ctx)
        violations.extend(spec_violations)
    
    # Check for specification completeness
    spec_completeness_violations = _validate_spec_completeness(changed_files, ctx, config)
    violations.extend(spec_completeness_violations)
    
    return violations

def _get_all_project_files(ctx):
    """Get list of all files in the project.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: All project file paths
    """
    # In a real implementation, this would scan the project directory
    # For now, we'll use the files from context plus some common locations
    all_files = []
    
    if "files" in ctx:
        files = ctx["files"]
        all_files.extend(files.get("added", []))
        all_files.extend(files.get("modified", []))
        # Don't include deleted files in project scan
    
    # Add common spec locations that might exist
    common_spec_paths = [
        "docs/specs/",
        "docs/design/", 
        "docs/requirements/",
        "specifications/",
        "design/",
        "requirements/"
    ]
    
    # In real implementation, would scan these directories
    # For validation purposes, assume they exist if mentioned in config
    
    return all_files

def _validate_file_has_spec(file_path, all_files, config, ctx):
    """Validate that a specific file has corresponding specification.
    
    Args:
        file_path (str): Path to code file
        all_files (list): All project files
        config (dict): Validator configuration
        ctx (dict): Validation context
        
    Returns:
        list: Violations for this file
    """
    violations = []
    
    # Look for related specification files
    related_specs = find_related_files(file_path, all_files, "spec")
    
    # Check configured spec patterns
    spec_patterns = config.get("spec_patterns", [
        "docs/specs/*.md",
        "docs/design/*.md",
        "docs/requirements/*.md"
    ])
    
    # Check if any specs exist for this file
    has_spec = len(related_specs) > 0
    
    if not has_spec:
        # Try to find specs based on patterns and file structure
        has_spec = _check_spec_patterns(file_path, spec_patterns)
    
    if not has_spec:
        violation = create_violation(
            rule="spec_before_code",
            severity="error",
            message=f"Code file '{file_path}' requires a specification document before implementation",
            file_path=file_path,
            suggestion=_get_spec_suggestion(file_path, spec_patterns)
        )
        violations.append(violation)
    
    return violations

def _check_spec_patterns(file_path, spec_patterns):
    """Check if specification exists based on configured patterns.
    
    Args:
        file_path (str): Code file path
        spec_patterns (list): List of specification path patterns
        
    Returns:
        bool: True if specification found
    """
    # Extract component/module name from file path
    file_parts = file_path.split("/")
    file_name = file_parts[-1].split(".")[0]
    
    # Try different spec naming conventions
    spec_names = [
        file_name,
        file_name.replace("_", "-"),
        file_parts[-2] if len(file_parts) > 1 else file_name,  # Directory name
    ]
    
    for pattern in spec_patterns:
        for name in spec_names:
            # Simple pattern matching - in real implementation would use glob
            if "*" in pattern:
                spec_path = pattern.replace("*", name)
                if file_exists(spec_path):
                    return True
    
    return False

def _get_spec_suggestion(file_path, spec_patterns):
    """Generate suggestion for creating specification.
    
    Args:
        file_path (str): Code file path
        spec_patterns (list): Configured spec patterns
        
    Returns:
        str: Suggestion text
    """
    file_name = file_path.split("/")[-1].split(".")[0]
    
    if spec_patterns:
        suggested_path = spec_patterns[0].replace("*", file_name)
        return f"Create specification document at '{suggested_path}' describing the purpose, requirements, and design of this code"
    
    return f"Create a specification document for '{file_name}' in docs/specs/ or docs/design/ describing its requirements and design"

def _validate_spec_completeness(changed_files, ctx, config):
    """Validate that specifications are complete and up-to-date.
    
    Args:
        changed_files (list): List of changed files
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Specification completeness violations
    """
    violations = []
    
    # Check if specification files themselves are being modified
    spec_files = [f for f in changed_files if is_spec_file(f)]
    
    for spec_file in spec_files:
        spec_violations = _validate_individual_spec(spec_file, ctx, config)
        violations.extend(spec_violations)
    
    return violations

def _validate_individual_spec(spec_file, ctx, config):
    """Validate an individual specification file for completeness.
    
    Args:
        spec_file (str): Path to specification file
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Violations for this specification
    """
    violations = []
    content = read_file(spec_file)
    
    if not content:
        violation = create_violation(
            rule="empty_specification",
            severity="error", 
            message=f"Specification file '{spec_file}' is empty",
            file_path=spec_file,
            suggestion="Add content describing requirements, design, and acceptance criteria"
        )
        violations.append(violation)
        return violations
    
    # Check for required sections based on configuration
    required_sections = config.get("required_sections", [
        "## Purpose",
        "## Requirements", 
        "## Design",
        "## Acceptance Criteria"
    ])
    
    missing_sections = []
    for section in required_sections:
        if section not in content and section.lower() not in content.lower():
            missing_sections.append(section)
    
    if missing_sections:
        violation = create_violation(
            rule="incomplete_specification",
            severity="warning",
            message=f"Specification '{spec_file}' is missing required sections: {', '.join(missing_sections)}",
            file_path=spec_file,
            suggestion="Add the missing sections to provide complete documentation"
        )
        violations.append(violation)
    
    # Check for placeholder content
    placeholder_patterns = [
        "TODO",
        "TBD", 
        "To be determined",
        "Fill this in",
        "Lorem ipsum"
    ]
    
    for pattern in placeholder_patterns:
        if pattern.lower() in content.lower():
            violation = create_violation(
                rule="placeholder_content",
                severity="warning",
                message=f"Specification '{spec_file}' contains placeholder content: '{pattern}'",
                file_path=spec_file,
                suggestion="Replace placeholder content with actual requirements and design details"
            )
            violations.append(violation)
            break  # Only report one placeholder violation per file
    
    return violations

def validate_api_specifications(ctx):
    """Validate API-specific specification requirements.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: API specification violations
    """
    violations = []
    changed_files = get_all_changed_files(ctx)
    
    # Look for API-related files
    api_files = [f for f in changed_files if _is_api_file(f)]
    
    for api_file in api_files:
        api_violations = _validate_api_file_spec(api_file, ctx)
        violations.extend(api_violations)
    
    return violations

def _is_api_file(file_path):
    """Check if file is API-related based on path and naming.
    
    Args:
        file_path (str): File path to check
        
    Returns:
        bool: True if file is API-related
    """
    api_indicators = [
        "/api/",
        "/handlers/",
        "/routes/",
        "/controllers/",
        "api.go",
        "handler.go",
        "routes.go"
    ]
    
    for indicator in api_indicators:
        if indicator in file_path:
            return True
    
    return False

def _validate_api_file_spec(api_file, ctx):
    """Validate that API file has proper OpenAPI/Swagger specification.
    
    Args:
        api_file (str): API file path
        ctx (dict): Validation context
        
    Returns:
        list: API specification violations
    """
    violations = []
    
    # Look for OpenAPI/Swagger specs
    spec_files = [
        "docs/api/openapi.yaml",
        "docs/api/swagger.yaml", 
        "api/openapi.yml",
        "openapi.json",
        "swagger.json"
    ]
    
    has_api_spec = any(file_exists(spec) for spec in spec_files)
    
    if not has_api_spec:
        violation = create_violation(
            rule="missing_api_specification",
            severity="error",
            message=f"API file '{api_file}' requires OpenAPI/Swagger specification",
            file_path=api_file,
            suggestion="Create OpenAPI specification in docs/api/openapi.yaml defining endpoints, schemas, and responses"
        )
        violations.append(violation)
    
    return violations