# Git operation utilities for CDP validators
# Provides functionality for parsing diffs, analyzing commit history, and file change detection

load("@stdlib//re", "re")
load("@stdlib//strings", "strings")

def get_changed_files(diff):
    """Parse git diff output to extract changed files.
    
    Args:
        diff (str): Git diff output
        
    Returns:
        dict: Dictionary with added, modified, deleted file lists
    """
    changes = {
        "added": [],
        "modified": [],
        "deleted": []
    }
    
    if not diff:
        return changes
    
    lines = diff.split("\n")
    current_file = None
    
    for line in lines:
        # Parse diff headers to identify files
        if line.startswith("diff --git"):
            # Extract file path from "diff --git a/path b/path"
            parts = line.split(" ")
            if len(parts) >= 4:
                file_path = parts[3][2:]  # Remove "b/" prefix
                current_file = file_path
        
        elif line.startswith("new file mode"):
            if current_file and current_file not in changes["added"]:
                changes["added"].append(current_file)
        
        elif line.startswith("deleted file mode"):
            if current_file and current_file not in changes["deleted"]:
                changes["deleted"].append(current_file)
        
        elif line.startswith("index") and current_file:
            # File modification (not new or deleted)
            if (current_file not in changes["added"] and 
                current_file not in changes["deleted"] and
                current_file not in changes["modified"]):
                changes["modified"].append(current_file)
    
    return changes

def get_changed_files_from_status(status):
    """Parse git status output to extract changed files.
    
    Args:
        status (str): Git status output
        
    Returns:
        dict: Dictionary with added, modified, deleted file lists
    """
    changes = {
        "added": [],
        "modified": [],
        "deleted": []
    }
    
    if not status:
        return changes
    
    lines = status.split("\n")
    
    for line in lines:
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        
        # Parse git status format: "XY filename"
        if len(line) >= 3:
            status_code = line[:2]
            filename = line[3:]
            
            # First character is staged, second is working tree
            staged = status_code[0]
            working = status_code[1]
            
            if staged == "A" or working == "A":
                if filename not in changes["added"]:
                    changes["added"].append(filename)
            elif staged == "D" or working == "D":
                if filename not in changes["deleted"]:
                    changes["deleted"].append(filename)
            elif staged == "M" or working == "M":
                if filename not in changes["modified"]:
                    changes["modified"].append(filename)
            elif staged == "?" and working == "?":
                # Untracked file
                if filename not in changes["added"]:
                    changes["added"].append(filename)
    
    return changes

def get_commit_history(log_output, limit=10):
    """Parse git log output to extract commit information.
    
    Args:
        log_output (str): Git log output
        limit (int): Maximum number of commits to parse
        
    Returns:
        list: List of commit dictionaries
    """
    if not log_output:
        return []
    
    commits = []
    current_commit = {}
    
    lines = log_output.split("\n")
    
    for line in lines:
        line = line.strip()
        
        if line.startswith("commit "):
            # Start of new commit
            if current_commit:
                commits.append(current_commit)
                if len(commits) >= limit:
                    break
            
            current_commit = {
                "hash": line.split(" ")[1],
                "author": "",
                "date": "",
                "message": ""
            }
        
        elif line.startswith("Author: "):
            current_commit["author"] = line[8:]
        
        elif line.startswith("Date: "):
            current_commit["date"] = line[6:]
        
        elif line and not line.startswith("commit ") and not line.startswith("Author: ") and not line.startswith("Date: "):
            # Commit message line
            if current_commit.get("message"):
                current_commit["message"] += "\n" + line
            else:
                current_commit["message"] = line
    
    # Add last commit
    if current_commit and len(commits) < limit:
        commits.append(current_commit)
    
    return commits

def parse_file_changes_from_diff(diff, file_path):
    """Extract specific changes for a file from diff output.
    
    Args:
        diff (str): Git diff output
        file_path (str): Specific file to analyze
        
    Returns:
        dict: File change information with added/removed lines
    """
    changes = {
        "added_lines": [],
        "removed_lines": [],
        "total_added": 0,
        "total_removed": 0
    }
    
    if not diff:
        return changes
    
    lines = diff.split("\n")
    in_target_file = False
    line_number = 0
    
    for line in lines:
        if line.startswith("diff --git") and file_path in line:
            in_target_file = True
            continue
        
        elif line.startswith("diff --git") and file_path not in line:
            in_target_file = False
            continue
        
        if not in_target_file:
            continue
        
        if line.startswith("@@"):
            # Parse hunk header to get line numbers
            match = re.match(r"@@ -(\d+),?\d* \+(\d+),?\d* @@", line)
            if match:
                line_number = int(match.group(2))
            continue
        
        if line.startswith("+") and not line.startswith("+++"):
            changes["added_lines"].append({"line": line_number, "content": line[1:]})
            changes["total_added"] += 1
            line_number += 1
        
        elif line.startswith("-") and not line.startswith("---"):
            changes["removed_lines"].append({"line": line_number, "content": line[1:]})
            changes["total_removed"] += 1
        
        elif not line.startswith("\\"):
            line_number += 1
    
    return changes

def get_file_blame_info(blame_output, line_number):
    """Parse git blame output to get author information for specific line.
    
    Args:
        blame_output (str): Git blame output
        line_number (int): Line number to get info for
        
    Returns:
        dict: Author and commit information for the line
    """
    if not blame_output:
        return {}
    
    lines = blame_output.split("\n")
    
    if line_number <= 0 or line_number > len(lines):
        return {}
    
    blame_line = lines[line_number - 1]
    
    # Parse blame format: "hash (author date line_num) content"
    match = re.match(r"([a-f0-9]+)\s+\(([^)]+)\s+(\d+)\)\s*(.*)", blame_line)
    
    if match:
        return {
            "hash": match.group(1),
            "author_info": match.group(2),
            "line_number": int(match.group(3)),
            "content": match.group(4)
        }
    
    return {}

def is_merge_commit(commit_message):
    """Check if commit message indicates a merge commit.
    
    Args:
        commit_message (str): Commit message
        
    Returns:
        bool: True if merge commit
    """
    merge_patterns = [
        r"^Merge branch",
        r"^Merge pull request",
        r"^Merge remote-tracking branch",
        r"^Merge tag"
    ]
    
    for pattern in merge_patterns:
        if re.match(pattern, commit_message):
            return True
    
    return False

def extract_branch_name(ctx):
    """Extract current branch name from context.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        str: Branch name or empty string
    """
    if "git" in ctx and "branch" in ctx["git"]:
        return ctx["git"]["branch"]
    return ""

def get_all_changed_files(ctx):
    """Get all changed files from context, combining multiple sources.
    
    Args:
        ctx (dict): Validation context
        
    Returns:
        list: List of all changed file paths
    """
    all_files = []
    
    if "files" in ctx:
        files = ctx["files"]
        
        if "added" in files:
            all_files.extend(files["added"])
        
        if "modified" in files:
            all_files.extend(files["modified"])
        
        if "deleted" in files:
            all_files.extend(files["deleted"])
    
    # Remove duplicates while preserving order
    seen = set()
    unique_files = []
    for file_path in all_files:
        if file_path not in seen:
            seen.add(file_path)
            unique_files.append(file_path)
    
    return unique_files

def is_feature_branch(branch_name):
    """Check if branch name indicates a feature branch.
    
    Args:
        branch_name (str): Branch name
        
    Returns:
        bool: True if feature branch
    """
    feature_patterns = [
        r"^feature/",
        r"^feat/",
        r"^enhancement/",
        r"^add/",
        r"^new/"
    ]
    
    for pattern in feature_patterns:
        if re.match(pattern, branch_name):
            return True
    
    return False

def count_commits_in_branch(commits, base_branch="main"):
    """Count commits that are unique to current branch.
    
    Args:
        commits (list): List of commit objects
        base_branch (str): Base branch name
        
    Returns:
        int: Number of commits in feature branch
    """
    # Simple heuristic: count commits until we find a merge commit
    # or reach a reasonable limit
    feature_commits = 0
    
    for commit in commits:
        if is_merge_commit(commit.get("message", "")):
            break
        feature_commits += 1
        if feature_commits > 50:  # Reasonable limit
            break
    
    return feature_commits