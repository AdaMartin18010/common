#!/usr/bin/env python3
import os
import re
import sys

def get_all_markdown_files(directory):
    """
    Get all Markdown files in a directory and its subdirectories.
    Returns a dictionary mapping relative paths to absolute paths.
    """
    markdown_files = {}
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.md'):
                abs_path = os.path.join(root, file)
                rel_path = os.path.relpath(abs_path, directory)
                markdown_files[rel_path] = abs_path
    return markdown_files

def fix_links_in_file(file_path, all_files, base_dir):
    """
    Fix local links in a Markdown file.
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
        
        # Find all Markdown links
        links = re.findall(r'\[([^\]]+)\]\(([^)]+)\)', content)
        
        modified = False
        for link_text, link_target in links:
            # Only process local links (not URLs)
            if not link_target.startswith(('http://', 'https://', '#')):
                # Remove any fragment identifier for path checking
                path_part = link_target.split('#')[0]
                
                # Handle relative paths
                if not path_part.startswith('/'):
                    # Get the directory of the current file
                    current_dir = os.path.dirname(file_path)
                    target_path = os.path.normpath(os.path.join(current_dir, path_part))
                    
                    # Check if the target file exists
                    if not os.path.exists(target_path):
                        # Try to find a matching file
                        rel_target = os.path.relpath(target_path, base_dir)
                        possible_matches = []
                        
                        for file_path in all_files.keys():
                            if os.path.basename(file_path) == os.path.basename(rel_target):
                                possible_matches.append(file_path)
                        
                        if len(possible_matches) == 1:
                            # Found a unique match
                            new_rel_path = os.path.relpath(
                                os.path.join(base_dir, possible_matches[0]), 
                                os.path.dirname(file_path)
                            )
                            
                            # Preserve any fragment identifier
                            if '#' in link_target:
                                new_rel_path += '#' + link_target.split('#')[1]
                            
                            # Replace the link in the content
                            old_link = f'[{link_text}]({link_target})'
                            new_link = f'[{link_text}]({new_rel_path})'
                            content = content.replace(old_link, new_link)
                            modified = True
                            print(f"  Fixed link: {old_link} -> {new_link}")
        
        if modified:
            with open(file_path, 'w', encoding='utf-8') as file:
                file.write(content)
            return True
        return False
    
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def process_directory(directory):
    """
    Process all Markdown files in a directory and its subdirectories.
    """
    all_files = get_all_markdown_files(directory)
    fixed_count = 0
    processed_count = 0
    
    for rel_path, abs_path in all_files.items():
        print(f"Processing {rel_path}...")
        processed_count += 1
        if fix_links_in_file(abs_path, all_files, directory):
            fixed_count += 1
    
    print(f"\nSummary:")
    print(f"Files processed: {processed_count}")
    print(f"Files with fixed links: {fixed_count}")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        directory = sys.argv[1]
    else:
        directory = "docs/refactor"
    
    if not os.path.isdir(directory):
        print(f"Error: {directory} is not a valid directory")
        sys.exit(1)
    
    print(f"Processing Markdown files in {directory}...")
    process_directory(directory) 