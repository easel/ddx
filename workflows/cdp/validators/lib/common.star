# Common utility functions for CDP validators
# Provides shared functionality for file operations, string manipulation, and validation helpers

load("@stdlib//os", "os")
load("@stdlib//re", "re")
load("@stdlib//strings", "strings")

def read_file(path):
    """Read file contents safely with error handling.
    
    Args:
        path (str): File path to read
        
    Returns:
        str: File contents or empty string if file doesn't exist
    """
    try:
        return os.read_file(path)
    except:
        return ""

def file_exists(path):
    """Check if file exists.
    
    Args:
        path (str): File path to check
        
    Returns:
        bool: True if file exists, False otherwise
    """
    try:
        os.read_file(path)
        return True
    except:
        return False

def is_code_file(path):
    """Check if file is a code file based on extension.
    
    Args:
        path (str): File path to check
        
    Returns:
        bool: True if file is a code file
    """
    code_extensions = [
        ".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".rs", 
        ".rb", ".php", ".cs", ".swift", ".kt", ".scala", ".clj"
    ]
    
    for ext in code_extensions:
        if path.endswith(ext):
            return True
    return False

def is_test_file(path):
    """Check if file is a test file based on naming patterns.
    
    Args:
        path (str): File path to check
        
    Returns:
        bool: True if file is a test file
    """
    test_patterns = [
        r".*_test\.(go|py|js|ts)$",
        r".*\.test\.(js|ts)$",
        r".*/test_.*\.py$",
        r".*/.*Test\.(java|kt|scala)$",
        r".*/.*_spec\.(rb|js|ts)$"
    ]
    
    for pattern in test_patterns:
        if re.match(pattern, path):
            return True
    return False

def is_spec_file(path):
    """Check if file is a specification document.
    
    Args:
        path (str): File path to check
        
    Returns:
        bool: True if file is a spec file
    """
    spec_patterns = [
        r".*/specs?/.*\.md$",
        r".*/requirements?/.*\.md$",
        r".*/design/.*\.md$",
        r".*/(PRD|BRD|FRD).*\.md$",
        r".*/.*_spec\.md$",
        r".*/.*_requirements\.md$"
    ]
    
    for pattern in spec_patterns:
        if re.match(pattern, path):
            return True
    return False

def get_file_extension(path):
    """Get file extension from path.
    
    Args:
        path (str): File path
        
    Returns:
        str: File extension including the dot
    """
    parts = path.split(".")
    if len(parts) > 1:
        return "." + parts[-1]
    return ""

def get_base_name(path):
    """Get base name of file without extension.
    
    Args:
        path (str): File path
        
    Returns:
        str: Base name without extension
    """
    name = path.split("/")[-1]
    if "." in name:
        return ".".join(name.split(".")[:-1])
    return name

def normalize_path(path):
    """Normalize file path by removing leading/trailing slashes.
    
    Args:
        path (str): File path to normalize
        
    Returns:
        str: Normalized path
    """
    return path.strip("/")

def create_violation(rule, severity, message, file_path, suggestion="", line=None):
    """Create a standardized violation object.
    
    Args:
        rule (str): Rule identifier
        severity (str): error|warning|info
        message (str): Human-readable violation description
        file_path (str): Path to violating file
        suggestion (str): Recommended fix
        line (int): Optional line number
        
    Returns:
        dict: Violation object
    """
    violation = {
        "rule": rule,
        "severity": severity,
        "message": message,
        "file": file_path,
        "suggestion": suggestion
    }
    
    if line != None:
        violation["line"] = line
    
    return violation

def find_related_files(target_file, file_list, relation_type="spec"):
    """Find files related to the target file based on naming patterns.
    
    Args:
        target_file (str): File to find relations for
        file_list (list): List of all files to search
        relation_type (str): Type of relation - "spec", "test", "impl"
        
    Returns:
        list: List of related file paths
    """
    base_name = get_base_name(target_file)
    target_dir = "/".join(target_file.split("/")[:-1])
    
    related = []
    
    if relation_type == "spec":
        # Look for specification files
        spec_patterns = [
            "docs/specs/" + base_name + ".md",
            "docs/design/" + base_name + ".md", 
            target_dir + "/" + base_name + "_spec.md",
            "specs/" + base_name + ".md"
        ]
        
        for pattern in spec_patterns:
            for file_path in file_list:
                if file_path.endswith(pattern) or pattern in file_path:
                    related.append(file_path)
    
    elif relation_type == "test":
        # Look for test files
        test_patterns = [
            base_name + "_test" + get_file_extension(target_file),
            base_name + ".test" + get_file_extension(target_file),
            "test_" + base_name + get_file_extension(target_file)
        ]
        
        for pattern in test_patterns:
            for file_path in file_list:
                if file_path.endswith(pattern):
                    related.append(file_path)
    
    elif relation_type == "impl":
        # Look for implementation files
        impl_ext = get_file_extension(target_file)
        impl_base = base_name.replace("_test", "").replace(".test", "").replace("test_", "")
        
        for file_path in file_list:
            if (get_base_name(file_path) == impl_base and 
                get_file_extension(file_path) == impl_ext and
                not is_test_file(file_path)):
                related.append(file_path)
    
    return related

def calculate_complexity_score(file_path, content=""):
    """Calculate complexity score for a file based on various metrics.
    
    Args:
        file_path (str): Path to the file
        content (str): File content (optional, will read if not provided)
        
    Returns:
        int: Complexity score (0-100)
    """
    if not content:
        content = read_file(file_path)
    
    if not content:
        return 0
    
    lines = content.split("\n")
    score = 0
    
    # Base complexity from line count
    score += min(len(lines) // 10, 20)
    
    # Complexity indicators
    complexity_patterns = [
        (r"\bif\b", 2),           # Conditional statements
        (r"\bfor\b", 2),          # Loops
        (r"\bwhile\b", 2),        # Loops
        (r"\btry\b", 3),          # Exception handling
        (r"\bcatch\b", 3),        # Exception handling
        (r"\bswitch\b", 4),       # Switch statements
        (r"\bclass\b", 5),        # Class definitions
        (r"\binterface\b", 3),    # Interface definitions
        (r"\bfunc\b", 1),         # Function definitions
        (r"\bdef\b", 1),          # Function definitions (Python)
    ]
    
    for pattern, weight in complexity_patterns:
        matches = re.findall(pattern, content)
        score += len(matches) * weight
    
    # Nested structure penalty
    max_indent = 0
    for line in lines:
        indent = len(line) - len(line.lstrip())
        max_indent = max(max_indent, indent)
    
    score += max_indent // 4
    
    return min(score, 100)

def get_project_config(ctx):
    """Extract project configuration from context.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        dict: Project configuration
    """
    if "project" in ctx and "config" in ctx["project"]:
        return ctx["project"]["config"]
    return {}

def get_validator_config(ctx, validator_name):
    """Get configuration for specific validator.
    
    Args:
        ctx (dict): Validation context
        validator_name (str): Name of validator
        
    Returns:
        dict: Validator-specific configuration
    """
    config = get_project_config(ctx)
    
    if "cdp" in config and "validators" in config["cdp"]:
        validators_config = config["cdp"]["validators"]
        if validator_name in validators_config:
            return validators_config[validator_name]
    
    return {"enabled": True}  # Default configuration