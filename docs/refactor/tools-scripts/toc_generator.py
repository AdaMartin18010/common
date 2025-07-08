import os

OUTPUT_FILE = 'navigation-index/TOC.md'
PROGRESS_FILE = 'navigation-index/progress.md'
ROOT_DIR = 'knowledge-system'

empty_dirs = []
no_md_dirs = []

def scan_dir(path, prefix='', level=1, max_level=4):
    entries = []
    if level > max_level:
        return entries
    has_md = False
    has_file = False
    for name in sorted(os.listdir(path)):
        full_path = os.path.join(path, name)
        if os.path.isdir(full_path):
            sub_entries = scan_dir(full_path, prefix + '  ', level+1, max_level)
            entries.append(f'{prefix}- **{name}/**')
            entries.extend(sub_entries)
            if not sub_entries:
                empty_dirs.append(full_path)
        else:
            has_file = True
            entries.append(f'{prefix}- {name}')
            if name.endswith('.md'):
                has_md = True
    if not has_file:
        empty_dirs.append(path)
    elif not has_md:
        no_md_dirs.append(path)
    return entries

def main():
    toc_lines = ['# 知识体系目录索引 (TOC)', '']
    toc_lines += scan_dir(ROOT_DIR)
    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        f.write('\n'.join(toc_lines))
    # 进度与内容检查
    progress_lines = ['# 目录内容检查与进度记录', '']
    progress_lines.append(f'## 空目录（无任何文件/子目录）: {len(empty_dirs)}')
    for d in empty_dirs:
        progress_lines.append(f'- {d}')
    progress_lines.append('')
    progress_lines.append(f'## 无md文件的目录: {len(no_md_dirs)}')
    for d in no_md_dirs:
        progress_lines.append(f'- {d}')
    with open(PROGRESS_FILE, 'w', encoding='utf-8') as f:
        f.write('\n'.join(progress_lines))
    print(f'已生成目录索引: {OUTPUT_FILE}')
    print(f'已生成进度记录: {PROGRESS_FILE}')

if __name__ == '__main__':
    main() 