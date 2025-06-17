#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速修正数学表达式格式的Python脚本
"""

import os
import re
import glob

def fix_math_format(content):
    """快速修正数学表达式格式"""
    modified = False
    
    # 修正不在数学环境中的 \text{} 命令
    def fix_text_command(match):
        nonlocal modified
        modified = True
        text_content = match.group(1)
        return f'$\\text{{{text_content}}}$'
    
    # 查找并修正不在数学环境中的 \text{}
    pattern = r'(?<!\$)\\text\{([^}]*)\}(?!\$)'
    if re.search(pattern, content):
        content = re.sub(pattern, fix_text_command, content)
    
    # 修正表格中的数学表达式
    def fix_table_math(match):
        nonlocal modified
        cell_content = match.group(1).strip()
        if '\\' in cell_content and not re.search(r'\$.*\$', cell_content):
            modified = True
            return f'& ${cell_content}$ &'
        return match.group(0)
    
    # 修正表格中的数学表达式
    table_pattern = r'& ([^&]*\\[a-zA-Z]+\{[^}]*\}[^&]*) &'
    if re.search(table_pattern, content):
        content = re.sub(table_pattern, fix_table_math, content)
    
    return content, modified

def process_files():
    """处理所有markdown文件"""
    md_files = glob.glob('**/*.md', recursive=True)
    
    print(f"Found {len(md_files)} markdown files")
    
    fixed_count = 0
    for filepath in md_files:
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            content, modified = fix_math_format(content)
            
            if modified:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"Fixed: {filepath}")
                fixed_count += 1
                
        except Exception as e:
            print(f"Error processing {filepath}: {e}")
    
    print(f"\nFixed {fixed_count} files out of {len(md_files)} total files")
    return fixed_count

if __name__ == '__main__':
    process_files() 