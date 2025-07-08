import os
import re

OUTPUT_FILE = 'navigation-index/knowledge-graph.md'
ROOT_DIR = 'knowledge-system'

link_pattern = re.compile(r'\[([^\]]+)\]\(([^)]+)\)')

# 收集所有md文件
md_files = []
for root, dirs, files in os.walk(ROOT_DIR):
    for file in files:
        if file.endswith('.md'):
            md_files.append(os.path.join(root, file))

# 构建知识点引用关系
graph = []
warnings = []
for md_path in md_files:
    rel_path = os.path.relpath(md_path, ROOT_DIR)
    try:
        with open(md_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        warnings.append(f'⚠️ 跳过无法解码文件: {rel_path} ({e})')
        continue
    links = link_pattern.findall(content)
    for text, link in links:
        if link.endswith('.md'):
            graph.append(f'- [{rel_path}] → [{text}]({link})')

with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
    f.write('# 知识体系交叉引用图谱\n\n')
    if graph:
        f.write('\n'.join(graph))
    else:
        f.write('（暂无交叉引用）\n')
    if warnings:
        f.write('\n\n## 警告\n')
        f.write('\n'.join(warnings))
print(f'已生成知识图谱: {OUTPUT_FILE}') 