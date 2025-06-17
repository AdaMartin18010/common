#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速修复markdown文件中的数学表达式
"""

import os
import re
import glob

def fix_math_in_file(file_path):
    """修复单个文件中的数学表达式"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original = content
        
        # 修复行内数学表达式
        content = re.sub(
            r'(?<!```latex\n)\$([^$]+)\$(?!\n```)',
            r'```latex\n$\1$\n```',
            content
        )
        
        # 修复块级数学表达式
        content = re.sub(
            r'(?<!```latex\n)\$\$([^$]+)\$\$(?!\n```)',
            r'```latex\n$$\1$$\n```',
            content
        )
        
        if content != original:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"Fixed: {file_path}")
            return True
        return False
        
    except Exception as e:
        print(f"Error: {file_path} - {e}")
        return False

def main():
    """主函数"""
    # 获取所有markdown文件
    md_files = []
    for root, dirs, files in os.walk('.'):
        for file in files:
            if file.endswith('.md'):
                md_files.append(os.path.join(root, file))
    
    print(f"Found {len(md_files)} markdown files")
    
    fixed = 0
    for file_path in md_files:
        if fix_math_in_file(file_path):
            fixed += 1
    
    print(f"Fixed {fixed}/{len(md_files)} files")

if __name__ == "__main__":
    main() 