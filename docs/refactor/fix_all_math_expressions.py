#!/usr/bin/env python3
import os
import re
import glob

def fix_math_expressions(file_path):
    """Fix LaTeX math expressions in markdown files by adding proper code block tags."""
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
            return False

    # Replace inline math expressions $...$ with proper LaTeX tags
    inline_pattern = r'(?<!\`\`\`latex\n)\$([^$\n]+?)\$(?!\n\`\`\`)'
    if re.search(inline_pattern, content):
        content = re.sub(inline_pattern, r'```latex\n$\1$\n```', content)
    
    # Replace block math expressions $$...$$ with proper LaTeX tags
    block_pattern = r'(?<!\`\`\`latex\n)\$\$([^$]+?)\$\$(?!\n\`\`\`)'
    if re.search(block_pattern, content):
        content = re.sub(block_pattern, r'```latex\n$$\1$$\n```', content)
    
    try:
        # Write the fixed content back to the file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    except Exception as e:
        print(f"Error writing {file_path}: {e}")
        return False

def process_directory(directory):
    """Process all markdown files in a directory and its subdirectories."""
    fixed_files = 0
    error_files = 0
    markdown_files = glob.glob(os.path.join(directory, '**/*.md'), recursive=True)
    
    for file_path in markdown_files:
        if fix_math_expressions(file_path):
            fixed_files += 1
            print(f"Fixed math expressions in {file_path}")
        else:
            error_files += 1
            print(f"Failed to fix math expressions in {file_path}")
    
    return fixed_files, error_files

if __name__ == "__main__":
    refactor_dir = "docs/refactor"
    
    if not os.path.isdir(refactor_dir):
        print(f"Error: {refactor_dir} directory not found.")
        exit(1)
    
    fixed_count, error_count = process_directory(refactor_dir)
    print(f"\nFixed math expressions in {fixed_count} files.")
    if error_count > 0:
        print(f"Failed to process {error_count} files.")
    print("Done!") 