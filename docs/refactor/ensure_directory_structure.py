#!/usr/bin/env python3
import os
import re
import sys
import shutil

def validate_directory_structure(directory):
    """
    Validate that directories follow the correct numbering pattern:
    XX-Name-Of-Directory where XX is a number.
    """
    pattern = re.compile(r'^(\d{2})-[A-Za-z0-9-]+$')
    
    issues = []
    for item in os.listdir(directory):
        item_path = os.path.join(directory, item)
        if os.path.isdir(item_path) and not item.startswith('.'):
            if not pattern.match(item):
                issues.append((item_path, "Directory name doesn't follow pattern XX-Name-Of-Directory"))
            else:
                # Recursively check subdirectories
                sub_issues = validate_directory_structure(item_path)
                issues.extend(sub_issues)
    
    return issues

def fix_directory_name(old_path, pattern=re.compile(r'^(\d+)-(.+)$')):
    """
    Fix directory name to follow the correct pattern.
    """
    dir_name = os.path.basename(old_path)
    parent_dir = os.path.dirname(old_path)
    
    # If it's already a numbered directory but wrong format
    match = pattern.match(dir_name)
    if match:
        number, name = match.groups()
        # Ensure two digits
        new_number = number.zfill(2)
        
        # Convert to Title-Case-With-Hyphens
        words = re.findall(r'[A-Z][a-z0-9]*|[a-z0-9]+', name)
        new_name = '-'.join(word.capitalize() for word in words)
        
        new_dir_name = f"{new_number}-{new_name}"
        new_path = os.path.join(parent_dir, new_dir_name)
        
        return new_path
    
    # If it's not a numbered directory, we need more information to fix it
    return None

def fix_directory_structure(directory):
    """
    Fix directory structure to follow correct numbering pattern.
    """
    issues = validate_directory_structure(directory)
    
    if not issues:
        print("All directories follow the correct pattern.")
        return
    
    print(f"Found {len(issues)} directories with naming issues:")
    for path, issue in issues:
        print(f"  {path}: {issue}")
        
        new_path = fix_directory_name(path)
        if new_path and new_path != path:
            try:
                # Rename the directory
                shutil.move(path, new_path)
                print(f"  Renamed to: {new_path}")
            except Exception as e:
                print(f"  Error renaming: {e}")

def validate_file_structure(directory):
    """
    Validate that files in each directory follow the correct numbering pattern:
    XX-Name-Of-File.md where XX is a number.
    """
    pattern = re.compile(r'^(\d{2})-[A-Za-z0-9-]+\.md$')
    
    issues = []
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.md') and file != 'README.md':
                if not pattern.match(file):
                    issues.append((os.path.join(root, file), "File name doesn't follow pattern XX-Name-Of-File.md"))
    
    return issues

def fix_file_name(old_path, pattern=re.compile(r'^(\d+)-(.+)\.md$')):
    """
    Fix file name to follow the correct pattern.
    """
    file_name = os.path.basename(old_path)
    dir_path = os.path.dirname(old_path)
    
    # If it's already a numbered file but wrong format
    match = pattern.match(file_name)
    if match:
        number, name = match.groups()
        # Ensure two digits
        new_number = number.zfill(2)
        
        # Convert to Title-Case-With-Hyphens
        name_without_ext = os.path.splitext(name)[0]
        words = re.findall(r'[A-Z][a-z0-9]*|[a-z0-9]+', name_without_ext)
        new_name = '-'.join(word.capitalize() for word in words)
        
        new_file_name = f"{new_number}-{new_name}.md"
        new_path = os.path.join(dir_path, new_file_name)
        
        return new_path
    
    # If it's not a numbered file, we need more information to fix it
    return None

def fix_file_structure(directory):
    """
    Fix file structure to follow correct numbering pattern.
    """
    issues = validate_file_structure(directory)
    
    if not issues:
        print("All files follow the correct pattern.")
        return
    
    print(f"Found {len(issues)} files with naming issues:")
    for path, issue in issues:
        print(f"  {path}: {issue}")
        
        new_path = fix_file_name(path)
        if new_path and new_path != path:
            try:
                # Rename the file
                shutil.move(path, new_path)
                print(f"  Renamed to: {new_path}")
            except Exception as e:
                print(f"  Error renaming: {e}")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        directory = sys.argv[1]
    else:
        directory = "docs/refactor"
    
    if not os.path.isdir(directory):
        print(f"Error: {directory} is not a valid directory")
        sys.exit(1)
    
    print(f"Validating directory structure in {directory}...")
    fix_directory_structure(directory)
    
    print(f"\nValidating file structure in {directory}...")
    fix_file_structure(directory) 