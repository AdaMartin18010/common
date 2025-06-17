#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
批量修正数学表达式格式的Python脚本
"""

import os
import re
import glob

def fix_math_expressions(content):
    """修正数学表达式格式"""
    modified = False
    
    # 修正 \text{} 格式 - 确保在数学环境中
    def fix_text_in_math(match):
        nonlocal modified
        modified = True
        text_content = match.group(1)
        return f'\\text{{{text_content}}}'
    
    # 查找并修正不在数学环境中的 \text{}
    content = re.sub(r'(?<!\$)\\text\{([^}]*)\}(?!\$)', fix_text_in_math, content)
    
    # 修正表格中的数学表达式
    def fix_table_math(match):
        nonlocal modified
        modified = True
        cell_content = match.group(1)
        # 如果包含数学表达式但没有被包围，则包围它
        if '\\' in cell_content and not re.search(r'\$.*\$', cell_content):
            return f'& ${cell_content}$ &'
        return match.group(0)
    
    content = re.sub(r'& ([^&]*\\[a-zA-Z]+\{[^}]*\}[^&]*) &', fix_table_math, content)
    
    # 修正其他常见的数学表达式格式问题
    # 确保所有数学符号都在正确的环境中
    math_patterns = [
        (r'(?<!\$)\\begin\{([^}]*)\}(.*?)\\end\{\1\}(?!\$)', r'$$\1\2\1$$'),
        (r'(?<!\$)\\sum_\{([^}]*)\}(?!\$)', r'$\sum_{\1}$'),
        (r'(?<!\$)\\int_\{([^}]*)\}(?!\$)', r'$\int_{\1}$'),
        (r'(?<!\$)\\frac\{([^}]*)\}\{([^}]*)\}(?!\$)', r'$\frac{\1}{\2}$'),
    ]
    
    for pattern, replacement in math_patterns:
        if re.search(pattern, content):
            content = re.sub(pattern, replacement, content)
            modified = True
    
    return content, modified

def process_file(filepath):
    """处理单个文件"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        content, modified = fix_math_expressions(content)
        
        if modified:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"Fixed: {filepath}")
            return True
        else:
            print(f"No changes needed: {filepath}")
            return False
            
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False

def main():
    """主函数"""
    # 获取所有markdown文件
    md_files = glob.glob('**/*.md', recursive=True)
    
    print(f"Found {len(md_files)} markdown files")
    
    fixed_count = 0
    for filepath in md_files:
        if process_file(filepath):
            fixed_count += 1
    
    print(f"\nFixed {fixed_count} files out of {len(md_files)} total files")

if __name__ == '__main__':
    main() 