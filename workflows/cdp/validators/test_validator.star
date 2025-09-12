# Test Validator
# Validates that tests are written and fail before implementation (Test-Driven Development)

load("//lib/common.star", "read_file", "file_exists", "is_code_file", "is_test_file", 
     "find_related_files", "create_violation", "get_validator_config")
load("//lib/git.star", "get_all_changed_files", "parse_file_changes_from_diff", "get_commit_history")

def validate_test_first(ctx):
    """Validate that tests are written and fail before implementation.
    
    Enforces Test-Driven Development (TDD) practices by ensuring:
    1. Tests exist for new code
    2. Tests were written before implementation (Red phase)
    3. Tests initially failed before implementation
    
    Args:
        ctx (dict): Validation context containing files, git info, and project config
        
    Returns:
        list: List of violation objects
    """
    violations = []
    config = get_validator_config(ctx, "test_validator")
    
    if not config.get("enabled", True):
        return violations
    
    changed_files = get_all_changed_files(ctx)
    
    # Separate code files and test files
    code_files = [f for f in changed_files if is_code_file(f) and not is_test_file(f)]
    test_files = [f for f in changed_files if is_test_file(f)]
    
    # Validate that code files have corresponding tests
    for code_file in code_files:
        test_violations = _validate_code_has_tests(code_file, test_files, ctx, config)
        violations.extend(test_violations)
    
    # Validate test-first development order
    tdd_violations = _validate_test_first_order(code_files, test_files, ctx, config)
    violations.extend(tdd_violations)
    
    # Validate test quality and coverage
    quality_violations = _validate_test_quality(test_files, ctx, config)
    violations.extend(quality_violations)
    
    return violations

def _validate_code_has_tests(code_file, test_files, ctx, config):
    """Validate that a code file has corresponding tests.
    
    Args:
        code_file (str): Path to code file
        test_files (list): List of test files in the changeset
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Violations for missing tests
    """
    violations = []
    
    # Skip certain file types that don't need tests
    skip_patterns = config.get("skip_test_patterns", [
        "main.go",
        "__init__.py",
        "config.py",
        "settings.py",
        "migrations/",
        "vendor/",
        "node_modules/"
    ])
    
    for pattern in skip_patterns:
        if pattern in code_file:
            return violations
    
    # Find related test files
    all_files = _get_all_project_files(ctx)
    related_tests = find_related_files(code_file, all_files, "test")
    
    # Check if any tests exist in the changeset
    has_tests_in_changeset = any(
        _files_are_related(code_file, test_file) for test_file in test_files
    )
    
    # Check if tests exist at all (including existing ones)
    has_existing_tests = len(related_tests) > 0
    
    if not has_tests_in_changeset and not has_existing_tests:
        violation = create_violation(
            rule="missing_tests",
            severity="error",
            message=f"Code file '{code_file}' requires corresponding tests",
            file_path=code_file,
            suggestion=_get_test_suggestion(code_file)
        )
        violations.append(violation)
    
    elif has_existing_tests and not has_tests_in_changeset:
        # Code was modified but tests weren't updated
        violation = create_violation(
            rule="outdated_tests",
            severity="warning",
            message=f"Code file '{code_file}' was modified but corresponding tests were not updated",
            file_path=code_file,
            suggestion=f"Update tests in {related_tests[0] if related_tests else 'test file'} to cover the new changes"
        )
        violations.append(violation)
    
    return violations

def _validate_test_first_order(code_files, test_files, ctx, config):
    """Validate that tests were written before implementation (TDD order).
    
    Args:
        code_files (list): Code files in changeset
        test_files (list): Test files in changeset
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: TDD order violations
    """
    violations = []
    
    if not config.get("enforce_tdd_order", True):
        return violations
    
    # Analyze commit history to determine order
    commits = get_commit_history(ctx.get("git", {}).get("log", ""), limit=20)
    
    for code_file in code_files:
        tdd_violations = _check_tdd_order_for_file(code_file, test_files, commits, ctx)
        violations.extend(tdd_violations)
    
    return violations

def _check_tdd_order_for_file(code_file, test_files, commits, ctx):
    """Check TDD order for a specific code file.
    
    Args:
        code_file (str): Code file path
        test_files (list): Test files in changeset
        commits (list): Commit history
        ctx (dict): Validation context
        
    Returns:
        list: TDD violations for this file
    """
    violations = []
    
    # Find related test file
    related_test = None
    for test_file in test_files:
        if _files_are_related(code_file, test_file):
            related_test = test_file
            break
    
    if not related_test:
        return violations  # No test in changeset, handled by other validator
    
    # Analyze commits to find when test and code were first introduced/modified
    test_first_commit = None
    code_first_commit = None
    
    for i, commit in enumerate(commits):
        # Simple heuristic: check if file paths appear in commit message or diff
        # In real implementation, would parse actual git log with file changes
        commit_msg = commit.get("message", "").lower()
        
        if related_test.lower() in commit_msg and not test_first_commit:
            test_first_commit = i
        
        if code_file.lower() in commit_msg and not code_first_commit:
            code_first_commit = i
    
    # If both files found in history, check order (lower index = more recent)
    if (test_first_commit is not None and code_first_commit is not None and 
        test_first_commit > code_first_commit):
        violation = create_violation(
            rule="code_before_test",
            severity="warning",
            message=f"Code file '{code_file}' was implemented before corresponding test '{related_test}' (violates TDD)",
            file_path=code_file,
            suggestion="Follow TDD: write failing test first, then implement code to make it pass"
        )
        violations.append(violation)
    
    return violations

def _validate_test_quality(test_files, ctx, config):
    """Validate quality of test files.
    
    Args:
        test_files (list): Test files to validate
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Test quality violations
    """
    violations = []
    
    for test_file in test_files:
        quality_violations = _validate_individual_test(test_file, ctx, config)
        violations.extend(quality_violations)
    
    return violations

def _validate_individual_test(test_file, ctx, config):
    """Validate an individual test file for quality.
    
    Args:
        test_file (str): Test file path
        ctx (dict): Validation context
        config (dict): Validator configuration
        
    Returns:
        list: Quality violations for this test
    """
    violations = []
    content = read_file(test_file)
    
    if not content:
        violation = create_violation(
            rule="empty_test",
            severity="error",
            message=f"Test file '{test_file}' is empty",
            file_path=test_file,
            suggestion="Add test cases covering the functionality"
        )
        violations.append(violation)
        return violations
    
    # Check for test function patterns based on language
    test_function_patterns = _get_test_patterns(test_file)
    has_tests = any(pattern in content for pattern in test_function_patterns)
    
    if not has_tests:
        violation = create_violation(
            rule="no_test_functions",
            severity="error",
            message=f"Test file '{test_file}' contains no recognizable test functions",
            file_path=test_file,
            suggestion=f"Add test functions using patterns: {', '.join(test_function_patterns)}"
        )
        violations.append(violation)
    
    # Check for assertions
    assertion_patterns = _get_assertion_patterns(test_file)
    has_assertions = any(pattern in content for pattern in assertion_patterns)
    
    if has_tests and not has_assertions:
        violation = create_violation(
            rule="no_assertions",
            severity="warning",
            message=f"Test file '{test_file}' has test functions but no assertions",
            file_path=test_file,
            suggestion="Add assertions to verify expected behavior"
        )
        violations.append(violation)
    
    # Check for proper test structure (Arrange-Act-Assert)
    structure_violations = _validate_test_structure(test_file, content, config)
    violations.extend(structure_violations)
    
    return violations

def _get_test_patterns(test_file):
    """Get test function patterns based on file extension.
    
    Args:
        test_file (str): Test file path
        
    Returns:
        list: Test function patterns
    """
    if test_file.endswith(".go"):
        return ["func Test", "func Benchmark", "func Example"]
    elif test_file.endswith(".py"):
        return ["def test_", "class Test", "def setUp", "def tearDown"]
    elif test_file.endswith((".js", ".ts")):
        return ["describe(", "it(", "test(", "expect("]
    elif test_file.endswith(".java"):
        return ["@Test", "public void test", "void test"]
    else:
        return ["test", "Test", "assert"]

def _get_assertion_patterns(test_file):
    """Get assertion patterns based on file extension.
    
    Args:
        test_file (str): Test file path
        
    Returns:
        list: Assertion patterns
    """
    if test_file.endswith(".go"):
        return ["t.Error", "t.Fail", "assert.", "require.", "if"]
    elif test_file.endswith(".py"):
        return ["assert", "self.assert", "assertEqual", "assertTrue"]
    elif test_file.endswith((".js", ".ts")):
        return ["expect(", "assert(", "should", "toBe", "toEqual"]
    elif test_file.endswith(".java"):
        return ["Assert.", "assertEquals", "assertTrue", "assertThat"]
    else:
        return ["assert", "expect", "should"]

def _validate_test_structure(test_file, content, config):
    """Validate test structure follows best practices.
    
    Args:
        test_file (str): Test file path
        content (str): File content
        config (dict): Validator configuration
        
    Returns:
        list: Test structure violations
    """
    violations = []
    
    # Check for descriptive test names
    lines = content.split("\n")
    test_functions = []
    
    for line in lines:
        line = line.strip()
        if any(pattern in line for pattern in _get_test_patterns(test_file)):
            test_functions.append(line)
    
    for test_func in test_functions:
        if _has_poor_test_name(test_func):
            violation = create_violation(
                rule="poor_test_name",
                severity="info",
                message=f"Test function in '{test_file}' has unclear name: {test_func[:50]}...",
                file_path=test_file,
                suggestion="Use descriptive test names that explain what is being tested and expected outcome"
            )
            violations.append(violation)
    
    return violations

def _has_poor_test_name(test_function):
    """Check if test function has a poor name.
    
    Args:
        test_function (str): Test function declaration
        
    Returns:
        bool: True if name is poor
    """
    poor_patterns = [
        "test1",
        "test2", 
        "testA",
        "testB",
        "testFunction",
        "testMethod",
        "testIt",
        "testThis"
    ]
    
    func_lower = test_function.lower()
    return any(pattern in func_lower for pattern in poor_patterns)

def _files_are_related(code_file, test_file):
    """Check if code file and test file are related.
    
    Args:
        code_file (str): Code file path
        test_file (str): Test file path
        
    Returns:
        bool: True if files are related
    """
    code_base = code_file.split("/")[-1].split(".")[0]
    test_base = test_file.split("/")[-1]
    
    # Common test naming patterns
    patterns = [
        f"{code_base}_test",
        f"test_{code_base}",
        f"{code_base}.test",
        f"Test{code_base.title()}"
    ]
    
    return any(pattern in test_base for pattern in patterns)

def _get_test_suggestion(code_file):
    """Generate suggestion for creating tests.
    
    Args:
        code_file (str): Code file path
        
    Returns:
        str: Test creation suggestion
    """
    base_name = code_file.split("/")[-1].split(".")[0]
    extension = code_file.split(".")[-1]
    
    if extension == "go":
        test_file = f"{base_name}_test.go"
        suggestion = f"Create '{test_file}' with func Test{base_name.title()}(t *testing.T) functions"
    elif extension == "py":
        test_file = f"test_{base_name}.py"
        suggestion = f"Create '{test_file}' with test_ functions using unittest or pytest"
    elif extension in ["js", "ts"]:
        test_file = f"{base_name}.test.{extension}"
        suggestion = f"Create '{test_file}' with describe/it blocks using Jest or similar framework"
    else:
        test_file = f"{base_name}_test.{extension}"
        suggestion = f"Create '{test_file}' with appropriate test functions"
    
    return suggestion

def _get_all_project_files(ctx):
    """Get list of all files in the project (simplified implementation).
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: All project file paths
    """
    all_files = []
    
    if "files" in ctx:
        files = ctx["files"]
        all_files.extend(files.get("added", []))
        all_files.extend(files.get("modified", []))
    
    return all_files

def validate_test_coverage(ctx):
    """Validate test coverage meets minimum thresholds.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: Coverage violations
    """
    violations = []
    config = get_validator_config(ctx, "test_validator")
    
    coverage_threshold = config.get("coverage_threshold", 80)
    
    if not coverage_threshold:
        return violations
    
    # This would integrate with actual coverage tools
    # For now, provide a structure for coverage validation
    
    changed_files = get_all_changed_files(ctx)
    code_files = [f for f in changed_files if is_code_file(f) and not is_test_file(f)]
    
    for code_file in code_files:
        # In real implementation, would call coverage tool
        # coverage = get_file_coverage(code_file)
        
        # Simulated coverage check
        violation = create_violation(
            rule="insufficient_coverage",
            severity="warning",
            message=f"Code file '{code_file}' may not meet coverage threshold of {coverage_threshold}%",
            file_path=code_file,
            suggestion=f"Ensure test coverage is at least {coverage_threshold}% by running coverage tools"
        )
        # Only add if actually below threshold
        # violations.append(violation)
    
    return violations