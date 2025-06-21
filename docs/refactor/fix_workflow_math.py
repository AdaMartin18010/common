#!/usr/bin/env python3
"""
修复工作流系统文件中的数学表达式
添加缺失的LaTeX标签
"""

import os
import re
from pathlib import Path

def fix_math_expressions(content):
    """修复数学表达式，添加缺失的LaTeX标签"""
    
    # 修复行内数学表达式 $...$
    content = re.sub(
        r'(\$[^$\n]+\$)',
        r'```latex\n\1\n```',
        content
    )
    
    # 修复块级数学表达式 $$...$$
    content = re.sub(
        r'(\$\$[^$]+\$\$)',
        r'```latex\n\1\n```',
        content
    )
    
    # 防止重复添加标签
    content = re.sub(
        r'```latex\n```latex\n',
        r'```latex\n',
        content
    )
    
    content = re.sub(
        r'\n```\n```',
        r'\n```',
        content
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
    # 获取工作流系统目录下的所有markdown文件
    workflow_dir = Path("docs/refactor/10-Workflow-Systems")
    markdown_files = list(workflow_dir.rglob("*.md"))
    
    print(f"🔍 找到 {len(markdown_files)} 个工作流系统markdown文件")
    
    fixed_count = 0
    for file_path in markdown_files:
        if process_file(file_path):
            fixed_count += 1
    
    print(f"\n📊 修复完成: {fixed_count}/{len(markdown_files)} 个文件")

if __name__ == "__main__":
    main() 