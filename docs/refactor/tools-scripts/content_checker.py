import os
import re

ROOT_DIR = 'knowledge-system'
REPORT_FILE = 'navigation-index/content-quality-report.md'

# 检查规则
def check_md_file(file_path):
    issues = []
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # 检查是否有标题
        if not re.search(r'^#\s+', content, re.MULTILINE):
            issues.append('缺少主标题')
            
        # 检查是否有摘要
        if '## 摘要' not in content and '摘要' not in content:
            issues.append('缺少摘要部分')
            
        # 检查是否有详细内容
        if '## 详细内容' not in content and '## 内容' not in content:
            issues.append('缺少详细内容部分')
            
        # 检查是否有参考文献
        if '## 参考文献' not in content and '## 参考' not in content:
            issues.append('缺少参考文献部分')
            
        # 检查是否有标签
        if '## 标签' not in content and '#' not in content:
            issues.append('缺少标签')
            
        # 检查文件大小
        if len(content) < 100:
            issues.append('内容过短（少于100字符）')
            
        # 检查是否有链接
        if not re.search(r'\[([^\]]+)\]\(([^)]+)\)', content):
            issues.append('缺少内部或外部链接')
            
    except Exception as e:
        issues.append(f'文件读取错误: {e}')
        
    return issues

# 扫描所有md文件
md_files = []
for root, dirs, files in os.walk(ROOT_DIR):
    for file in files:
        if file.endswith('.md'):
            md_files.append(os.path.join(root, file))

# 生成检查报告
report_lines = ['# 内容规范检查报告', '']
report_lines.append(f'## 检查统计')
report_lines.append(f'- 总文件数: {len(md_files)}')
report_lines.append('')

issues_count = 0
files_with_issues = 0

for md_file in md_files:
    rel_path = os.path.relpath(md_file, ROOT_DIR)
    issues = check_md_file(md_file)
    
    if issues:
        files_with_issues += 1
        issues_count += len(issues)
        report_lines.append(f'### {rel_path}')
        for issue in issues:
            report_lines.append(f'- ⚠️ {issue}')
        report_lines.append('')

report_lines.append(f'## 总结')
report_lines.append(f'- 有问题的文件: {files_with_issues}/{len(md_files)}')
report_lines.append(f'- 总问题数: {issues_count}')
report_lines.append(f'- 平均每文件问题数: {issues_count/len(md_files) if md_files else 0:.1f}')

with open(REPORT_FILE, 'w', encoding='utf-8') as f:
    f.write('\n'.join(report_lines))

print(f'已生成内容规范检查报告: {REPORT_FILE}')
print(f'检查了 {len(md_files)} 个文件，发现 {files_with_issues} 个文件有问题，共 {issues_count} 个问题') 