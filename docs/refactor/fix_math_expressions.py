#!/usr/bin/env python3
"""
修复markdown文件中的数学表达式
将所有没有```latex标签的数学表达式正确格式化
"""

import os
import re
import glob

def fix_math_expressions(file_path):
    """修复单个文件中的数学表达式"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        
        # 修复行内数学表达式
        # 匹配 $...$ 格式但不包含 ```latex 标签的
        content = re.sub(
            r'(?<!```latex\n)\$([^$\n]+)\$(?!\n```)',
            r'$`\1`$',
            content
        )
        
        # 修复块级数学表达式
        # 匹配 $$...$$ 格式但不包含 ```latex 标签的
        content = re.sub(
            r'(?<!```latex\n)\$\$([^$]+)\$\$(?!\n```)',
            r'```latex\n$$\1$$\n```',
            content
        )
        
        # 修复表格中的数学表达式
        content = re.sub(
            r'(\|[^|]*)\$([^$\n]+)\$([^|]*\|)',
            r'\1$`\2`$\3',
            content
        )
        
        # 修复条件表达式中的数学符号
        content = re.sub(
            r'\\text\{([^}]+)\}',
            r'\\text{\1}',
            content
        )
        
        # 修复反斜杠转义
        content = re.sub(
            r'\\\\([a-zA-Z]+)',
            r'\\\\\1',
            content
        )
        
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"修复: {file_path}")
            return True
        return False
        
    except Exception as e:
        print(f"处理文件 {file_path} 时出错: {e}")
        return False

def main():
    """主函数"""
    # 获取所有markdown文件
    md_files = []
    
    # 搜索docs/refactor目录下的所有.md文件
    for root, dirs, files in os.walk('docs/refactor'):
        for file in files:
            if file.endswith('.md'):
                md_files.append(os.path.join(root, file))
    
    print(f"找到 {len(md_files)} 个markdown文件")
    
    fixed_count = 0
    for file_path in md_files:
        if fix_math_expressions(file_path):
            fixed_count += 1
    
    print(f"修复了 {fixed_count} 个文件")

if __name__ == "__main__":
    main() 