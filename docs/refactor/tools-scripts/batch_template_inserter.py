import os
import shutil

PROGRESS_FILE = 'navigation-index/progress.md'
TEMPLATE_FILE = 'document-templates/knowledge-entry-template.md'

# 读取进度文件，获取空目录和无md目录
empty_dirs = []
no_md_dirs = []
try:
    with open(PROGRESS_FILE, 'r', encoding='utf-8') as f:
        content = f.read()
        lines = content.split('\n')
        in_empty_section = False
        in_no_md_section = False
        for line in lines:
            if '空目录' in line:
                in_empty_section = True
                in_no_md_section = False
            elif '无md文件的目录' in line:
                in_empty_section = False
                in_no_md_section = True
            elif line.startswith('- ') and in_empty_section:
                empty_dirs.append(line.strip('- '))
            elif line.startswith('- ') and in_no_md_section:
                no_md_dirs.append(line.strip('- '))
except Exception as e:
    print(f'⚠️ 读取进度文件失败: {e}')

# 批量插入模板
inserted_count = 0
for d in set(empty_dirs + no_md_dirs):
    target = os.path.join(d, 'README.md')
    if not os.path.exists(target):
        try:
            shutil.copyfile(TEMPLATE_FILE, target)
            print(f'✅ 已插入模板: {target}')
            inserted_count += 1
        except Exception as e:
            print(f'⚠️ 插入失败: {target} ({e})')

print(f'\n📊 批量补全完成: 成功插入 {inserted_count} 个模板文件') 