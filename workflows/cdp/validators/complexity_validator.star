# Complexity Validator
# Validates complexity constraints to maintain manageable development pace

load("//lib/common.star", "read_file", "is_code_file", "calculate_complexity_score", 
     "create_violation", "get_validator_config")
load("//lib/git.star", "get_all_changed_files", "extract_branch_name", "is_feature_branch", 
     "count_commits_in_branch", "get_commit_history")

def validate_complexity(ctx):
    """Validate complexity constraints for manageable development.
    
    Enforces complexity limits to prevent:
    1. Too many concurrent features in development
    2. Overly complex individual changes
    3. Large, monolithic commits
    4. Excessive cognitive load
    
    Args:
        ctx (dict): Validation context containing files, git info, and project config
        
    Returns:
        list: List of violation objects
    """
    violations = []
    config = get_validator_config(ctx, "complexity_validator")
    
    if not config.get("enabled", True):
        return violations
    
    # Validate concurrent feature limits
    concurrent_violations = _validate_concurrent_features(ctx, config)
    violations.extend(concurrent_violations)
    
    # Validate individual file complexity
    file_complexity_violations = _validate_file_complexity(ctx, config)
    violations.extend(file_complexity_violations)
    
    # Validate commit complexity
    commit_complexity_violations = _validate_commit_complexity(ctx, config)
    violations.extend(commit_complexity_violations)
    
    # Validate change set size
    changeset_violations = _validate_changeset_size(ctx, config)
    violations.extend(changeset_violations)
    
    return violations

def _validate_concurrent_features(ctx, config):
    """Validate number of concurrent features in development.
    
    Args:
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Concurrent feature violations
    """
    violations = []
    
    max_concurrent = config.get("max_concurrent_features", 3)
    branch_name = extract_branch_name(ctx)
    
    if not is_feature_branch(branch_name):
        return violations  # Not a feature branch, skip check
    
    # In a real implementation, this would:
    # 1. List all open feature branches
    # 2. Count branches with recent activity
    # 3. Check against configured limit
    
    # Simulated check - would integrate with git branch listing
    active_features = _count_active_feature_branches(ctx)
    
    if active_features > max_concurrent:
        violation = create_violation(
            rule="too_many_concurrent_features",
            severity="warning",
            message=f"Too many concurrent features in development ({active_features} > {max_concurrent})",
            file_path=branch_name,
            suggestion=f"Complete existing features before starting new ones. Limit: {max_concurrent} concurrent features"
        )
        violations.append(violation)
    
    return violations

def _validate_file_complexity(ctx, config):
    """Validate complexity of individual files.
    
    Args:
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: File complexity violations
    """
    violations = []
    
    complexity_threshold = config.get("complexity_threshold", 50)
    changed_files = get_all_changed_files(ctx)
    
    for file_path in changed_files:
        if not is_code_file(file_path):
            continue
        
        complexity_score = calculate_complexity_score(file_path)
        
        if complexity_score > complexity_threshold:
            severity = "error" if complexity_score > complexity_threshold * 1.5 else "warning"
            
            violation = create_violation(
                rule="excessive_file_complexity",
                severity=severity,
                message=f"File '{file_path}' has high complexity score: {complexity_score} (threshold: {complexity_threshold})",
                file_path=file_path,
                suggestion=_get_complexity_reduction_suggestion(file_path, complexity_score)
            )
            violations.append(violation)
    
    return violations

def _validate_commit_complexity(ctx, config):
    """Validate complexity of individual commits.
    
    Args:
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Commit complexity violations
    """
    violations = []
    
    max_files_per_commit = config.get("max_files_per_commit", 10)
    max_lines_per_commit = config.get("max_lines_per_commit", 500)
    
    changed_files = get_all_changed_files(ctx)
    
    # Check number of files in commit
    if len(changed_files) > max_files_per_commit:
        violation = create_violation(
            rule="too_many_files_in_commit",
            severity="warning",
            message=f"Commit touches too many files ({len(changed_files)} > {max_files_per_commit})",
            file_path="commit",
            suggestion="Split large commits into smaller, focused changes"
        )
        violations.append(violation)
    
    # Check total lines changed (would need git diff stats)
    total_lines_changed = _calculate_total_lines_changed(ctx, changed_files)
    
    if total_lines_changed > max_lines_per_commit:
        violation = create_violation(
            rule="too_many_lines_in_commit",
            severity="warning", 
            message=f"Commit changes too many lines ({total_lines_changed} > {max_lines_per_commit})",
            file_path="commit",
            suggestion="Break large changes into smaller, reviewable commits"
        )
        violations.append(violation)
    
    return violations

def _validate_changeset_size(ctx, config):
    """Validate overall changeset size for the branch.
    
    Args:
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Changeset size violations
    """
    violations = []
    
    max_branch_complexity = config.get("max_branch_complexity", 200)
    branch_name = extract_branch_name(ctx)
    
    if not is_feature_branch(branch_name):
        return violations
    
    # Calculate total complexity for the entire branch
    branch_complexity = _calculate_branch_complexity(ctx)
    
    if branch_complexity > max_branch_complexity:
        violation = create_violation(
            rule="excessive_branch_complexity",
            severity="error",
            message=f"Branch '{branch_name}' has excessive complexity: {branch_complexity} (max: {max_branch_complexity})",
            file_path=branch_name,
            suggestion="Split this feature into smaller, more manageable branches"
        )
        violations.append(violation)
    
    return violations

def _count_active_feature_branches(ctx):
    """Count active feature branches (simulated).
    
    In real implementation, this would:
    1. Execute git branch -r to list remote branches
    2. Filter for feature/ prefix
    3. Check recent activity with git log
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        int: Number of active feature branches
    """
    # Simulated count - in real implementation would query git
    return 2  # Placeholder

def _get_complexity_reduction_suggestion(file_path, complexity_score):
    """Generate suggestions for reducing file complexity.
    
    Args:
        file_path (str): Path to complex file
        complexity_score (int): Current complexity score
        
    Returns:
        str: Suggestion for reducing complexity
    """
    suggestions = []
    
    if complexity_score > 80:
        suggestions.append("Consider splitting into multiple files or modules")
    
    if complexity_score > 60:
        suggestions.append("Extract complex functions into separate functions")
        suggestions.append("Reduce nesting levels and conditional complexity")
    
    suggestions.extend([
        "Simplify complex conditional logic",
        "Extract reusable components",
        "Consider using design patterns to reduce complexity"
    ])
    
    return ". ".join(suggestions[:2])  # Return first 2 suggestions

def _calculate_total_lines_changed(ctx, changed_files):
    """Calculate total lines changed in commit.
    
    Args:
        ctx (dict): Validation context
        changed_files (list): List of changed files
        
    Returns:
        int: Total lines changed
    """
    total = 0
    
    # In real implementation, would parse git diff --stat
    # For now, estimate based on file count and content
    for file_path in changed_files:
        if is_code_file(file_path):
            content = read_file(file_path)
            if content:
                # Rough estimate: assume 10% of file changed
                lines = len(content.split("\n"))
                total += max(1, lines // 10)
    
    return total

def _calculate_branch_complexity(ctx):
    """Calculate total complexity for the entire branch.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        int: Branch complexity score
    """
    total_complexity = 0
    changed_files = get_all_changed_files(ctx)
    
    # Sum complexity of all changed files
    for file_path in changed_files:
        if is_code_file(file_path):
            file_complexity = calculate_complexity_score(file_path)
            total_complexity += file_complexity
    
    # Add commit count penalty
    commits = get_commit_history(ctx.get("git", {}).get("log", ""), limit=50)
    branch_commits = count_commits_in_branch(commits)
    
    # Each commit adds to complexity
    total_complexity += branch_commits * 5
    
    return total_complexity

def validate_cognitive_load(ctx):
    """Validate cognitive load limits for reviewability.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: Cognitive load violations
    """
    violations = []
    config = get_validator_config(ctx, "complexity_validator")
    
    max_cognitive_load = config.get("max_cognitive_load", 100)
    changed_files = get_all_changed_files(ctx)
    
    # Calculate cognitive load based on various factors
    cognitive_load = 0
    
    # Factor 1: Number of different file types
    file_types = set()
    for file_path in changed_files:
        if "." in file_path:
            file_types.add(file_path.split(".")[-1])
    
    cognitive_load += len(file_types) * 10
    
    # Factor 2: Number of different directories
    directories = set()
    for file_path in changed_files:
        directory = "/".join(file_path.split("/")[:-1])
        directories.add(directory)
    
    cognitive_load += len(directories) * 5
    
    # Factor 3: Mix of new and modified files
    if "files" in ctx:
        files = ctx["files"]
        added = len(files.get("added", []))
        modified = len(files.get("modified", []))
        deleted = len(files.get("deleted", []))
        
        # Mixed operations increase cognitive load
        operations = sum(1 for x in [added, modified, deleted] if x > 0)
        cognitive_load += operations * 15
    
    if cognitive_load > max_cognitive_load:
        violation = create_violation(
            rule="excessive_cognitive_load",
            severity="warning",
            message=f"Change has high cognitive load: {cognitive_load} (max: {max_cognitive_load})",
            file_path="changeset",
            suggestion="Simplify changeset by focusing on single concern or splitting into multiple changes"
        )
        violations.append(violation)
    
    return violations

def validate_dependency_complexity(ctx):
    """Validate complexity introduced by dependencies.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: Dependency complexity violations
    """
    violations = []
    changed_files = get_all_changed_files(ctx)
    
    # Look for dependency-related files
    dependency_files = [
        "go.mod", "go.sum",
        "package.json", "package-lock.json", "yarn.lock",
        "requirements.txt", "Pipfile", "Pipfile.lock",
        "Cargo.toml", "Cargo.lock",
        "pom.xml", "build.gradle"
    ]
    
    changed_deps = [f for f in changed_files if any(df in f for df in dependency_files)]
    
    if changed_deps:
        # Check for major version changes or new dependencies
        for dep_file in changed_deps:
            content = read_file(dep_file)
            if content:
                dep_violations = _analyze_dependency_changes(dep_file, content, ctx)
                violations.extend(dep_violations)
    
    return violations

def _analyze_dependency_changes(dep_file, content, ctx):
    """Analyze dependency file changes for complexity issues.
    
    Args:
        dep_file (str): Dependency file path
        content (str): File content
        ctx (dict): Validation context
        
    Returns:
        list: Dependency-related violations
    """
    violations = []
    
    # Look for patterns that indicate major changes
    major_change_patterns = [
        "v0.", "v1.", "v2.", "v3.", "v4.", "v5.",  # Major version changes
        "beta", "alpha", "rc",  # Pre-release versions
        "SNAPSHOT",  # Maven snapshots
    ]
    
    for pattern in major_change_patterns:
        if pattern in content:
            violation = create_violation(
                rule="major_dependency_change",
                severity="warning",
                message=f"Dependency file '{dep_file}' contains potentially breaking changes ({pattern})",
                file_path=dep_file,
                suggestion="Ensure compatibility testing for major dependency changes"
            )
            violations.append(violation)
            break  # Only report once per file
    
    return violations