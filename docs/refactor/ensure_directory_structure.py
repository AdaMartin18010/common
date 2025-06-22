#!/usr/bin/env python3
import os
import re

def check_directory_name(dir_name):
    """Check if directory name follows the numeric pattern: 01-Name-Format."""
    pattern = r'^\d{2}-[A-Z][a-zA-Z0-9-]+'
    return bool(re.match(pattern, os.path.basename(dir_name)))

def check_file_name(file_name):
    """Check if file name follows the pattern: 01-Name-Format.md."""
    if file_name == "README.md":
        return True
    pattern = r'^\d{2}-[A-Z][a-zA-Z0-9-]+\.md$'
    return bool(re.match(pattern, os.path.basename(file_name)))

def check_structure(directory):
    """Check directory structure for naming issues."""
    issues = []
    
    for root, dirs, files in os.walk(directory):
        # Skip the root directory itself
        if root == directory:
            continue
        
        # Check immediate subdirectories of root
        if os.path.dirname(root) == directory:
            if not check_directory_name(root):
                issues.append(f"Directory name issue: {os.path.relpath(root, directory)}")
        
        # Check markdown files
        for file in files:
            if file.endswith('.md') and file != "README.md":
                if not check_file_name(file):
                    issues.append(f"File name issue: {os.path.relpath(os.path.join(root, file), directory)}")
    
    return issues

def print_report(issues):
    """Print a report of structure issues."""
    if not issues:
        print("✅ Directory structure follows naming conventions.")
        return
    
    print(f"❌ Found {len(issues)} structure issues:")
    for issue in issues:
        print(f"  - {issue}")

if __name__ == "__main__":
    refactor_dir = "docs/refactor"
    
    if not os.path.isdir(refactor_dir):
        print(f"Error: {refactor_dir} directory not found.")
        exit(1)
    
    # Don't check top-level files in the refactor directory
    issues = check_structure(refactor_dir)
    print_report(issues)
    
    if issues:
        exit(1)
    else:
        print("Directory structure verified successfully!") 