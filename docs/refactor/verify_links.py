#!/usr/bin/env python3
import os
import re
import glob
from pathlib import Path

def get_all_files(directory):
    """Get all files in a directory and its subdirectories."""
    all_files = set()
    for root, dirs, files in os.walk(directory):
        for file in files:
            rel_path = os.path.relpath(os.path.join(root, file), directory)
            all_files.add(rel_path.replace('\\', '/'))  # Normalize path separators
    return all_files

def extract_links(file_path):
    """Extract all internal markdown links from a file."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except UnicodeDecodeError:
        try:
            # Try another common encoding
            with open(file_path, 'r', encoding='latin1') as f:
                content = f.read()
        except Exception as e:
            print(f"Error reading {file_path}: {e}")
            return []
    
    # Find all markdown links
    links = re.findall(r'\[([^\]]+)\]\(([^)]+)\)', content)
    
    # Filter to keep only local file links (not URLs)
    internal_links = []
    for link_text, link_path in links:
        if not (link_path.startswith('http://') or link_path.startswith('https://') or link_path.startswith('#')):
            # Only include file links, not anchors
            internal_links.append((link_text, link_path))
    
    return internal_links

def verify_link(link_path, file_path, all_files):
    """Verify if a link is valid."""
    # Handle absolute paths within the workspace
    if link_path.startswith('/'):
        target = link_path[1:]  # Remove leading /
        return target in all_files
    
    # Handle relative paths
    file_dir = os.path.dirname(file_path)
    target_path = os.path.normpath(os.path.join(file_dir, link_path))
    target_path = target_path.replace('\\', '/')  # Normalize path separators
    
    # Check if the target file exists
    return target_path in all_files

def process_directory(directory):
    """Process all markdown files in a directory and check their links."""
    all_files = get_all_files(directory)
    broken_links = []
    error_files = 0
    
    markdown_files = glob.glob(os.path.join(directory, '**/*.md'), recursive=True)
    
    for file_path in markdown_files:
        try:
            rel_file_path = os.path.relpath(file_path, directory).replace('\\', '/')
            links = extract_links(file_path)
            
            for link_text, link_path in links:
                if not verify_link(link_path, rel_file_path, all_files):
                    broken_links.append((rel_file_path, link_text, link_path))
        except Exception as e:
            error_files += 1
            print(f"Error processing {file_path}: {e}")
    
    return broken_links, error_files

def print_report(broken_links, error_files=0):
    """Print a report of broken links."""
    if error_files > 0:
        print(f"⚠️ Failed to process {error_files} files due to errors.")
    
    if not broken_links:
        print("✅ No broken links found.")
        return
    
    print(f"❌ Found {len(broken_links)} broken links:")
    current_file = None
    
    for file_path, link_text, link_path in sorted(broken_links):
        if file_path != current_file:
            print(f"\nIn file: {file_path}")
            current_file = file_path
        
        print(f"  - [{link_text}]({link_path})")

if __name__ == "__main__":
    refactor_dir = "docs/refactor"
    
    if not os.path.isdir(refactor_dir):
        print(f"Error: {refactor_dir} directory not found.")
        exit(1)
    
    broken_links, error_files = process_directory(refactor_dir)
    print_report(broken_links, error_files)
    
    if broken_links or error_files > 0:
        exit(1)
    else:
        print("All links verified successfully!") 