#!/usr/bin/env python3
"""
修复markdown文件中的数学表达式
添加缺失的LaTeX标签
"""

import os
import re
import glob
from pathlib import Path

def fix_math_expressions(content):
    """修复数学表达式，添加缺失的LaTeX标签"""
    
    # 修复行内数学表达式
    # 匹配 $...$ 格式但缺少 ```latex 标签的
    content = re.sub(
        r'(?<!```latex\n)\$([^$]+)\$(?!\n```)',
        r'```latex\n$\1$\n```',
        content
    )
    
    # 修复块级数学表达式
    # 匹配 $$...$$ 格式但缺少 ```latex 标签的
    content = re.sub(
        r'(?<!```latex\n)\$\$([^$]+)\$\$(?!\n```)',
        r'```latex\n$$\1$$\n```',
        content
    )
    
    # 修复行内数学表达式（反斜杠转义）
    # 匹配 \$...\$ 格式
    content = re.sub(
        r'(?<!```latex\n)\\\$([^$]+)\\\$(?!\n```)',
        r'```latex\n$\1$\n```',
        content
    )
    
    # 修复块级数学表达式（反斜杠转义）
    # 匹配 \$\$...\$\$ 格式
    content = re.sub(
        r'(?<!```latex\n)\\\$\$([^$]+)\\\$\$(?!\n```)',
        r'```latex\n$$\1$$\n```',
        content
    )
    
    # 修复表格中的数学表达式
    # 匹配表格中的数学表达式
    def fix_table_math(match):
        table_content = match.group(0)
        # 在表格中的数学表达式周围添加LaTeX标签
        table_content = re.sub(
            r'\$([^$]+)\$',
            r'```latex\n$\1$\n```',
            table_content
        )
        return table_content
    
    content = re.sub(
        r'\|.*\|.*\|',
        fix_table_math,
        content,
        flags=re.MULTILINE
    )
    
    return content

def process_file(file_path):
    """处理单个文件"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        content = fix_math_expressions(content)
        
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"✅ 修复: {file_path}")
            return True
        else:
            print(f"⏭️  跳过: {file_path}")
            return False
            
    except Exception as e:
        print(f"❌ 错误: {file_path} - {e}")
        return False

def main():
    """主函数"""
    # 获取所有markdown文件
    refactor_dir = Path("docs/refactor")
    markdown_files = list(refactor_dir.rglob("*.md"))
    
    print(f"🔍 找到 {len(markdown_files)} 个markdown文件")
    
    fixed_count = 0
    for file_path in markdown_files:
        if process_file(file_path):
            fixed_count += 1
    
    print(f"\n📊 修复完成: {fixed_count}/{len(markdown_files)} 个文件")

if __name__ == "__main__":
    main() 