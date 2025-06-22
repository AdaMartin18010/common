#!/usr/bin/env python3
import os
import re
import sys

def fix_latex_in_file(file_path):
    """
    Fix LaTeX formatting issues in a Markdown file.
    
    Issues to fix:
    1. ```latex without closing tag
    2. Nested LaTeX tags
    3. Improper LaTeX dollar sign formatting
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
        
        # Fix 1: Ensure LaTeX blocks are properly formatted
        # Replace ```latex\n$...$\n``` with ```latex\n...\n```
        content = re.sub(r'```latex\s*\n\s*\$(.+?)\$\s*\n\s*```', r'```latex\n\1\n```', content, flags=re.DOTALL)
        
        # Fix 2: Fix inline LaTeX expressions
        # Replace ```latex\n$...$\n``` with $...$
        content = re.sub(r'```latex\s*\n\s*(.+?)\s*\n\s*```', r'```latex\n\1\n```', content, flags=re.DOTALL)
        
        # Fix 3: Replace incorrect nested LaTeX tags
        # ```latex\n$```latex\n...\n```$\n``` with ```latex\n...\n```
        content = re.sub(r'```latex\s*\n\s*\$```latex\s*\n(.+?)\n```\$\s*\n\s*```', r'```latex\n\1\n```', content, flags=re.DOTALL)
        
        # Fix 4: Fix inline LaTeX expressions with incorrect formatting
        # Replace ```latex\n...\n``` with $...$
        content = re.sub(r'```latex\s*\n\s*\$(.+?)\$\s*\n\s*```', r'$\1$', content, flags=re.DOTALL)
        
        # Fix 5: Fix missing LaTeX tags
        # Look for math expressions that should be in LaTeX but aren't
        # This is more complex and might need manual review
        
        with open(file_path, 'w', encoding='utf-8') as file:
            file.write(content)
        
        return True
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def process_directory(directory):
    """
    Process all Markdown files in a directory and its subdirectories.
    """
    fixed_count = 0
    error_count = 0
    
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                print(f"Processing {file_path}...")
                if fix_latex_in_file(file_path):
                    fixed_count += 1
                else:
                    error_count += 1
    
    print(f"\nSummary:")
    print(f"Files processed successfully: {fixed_count}")
    print(f"Files with errors: {error_count}")

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