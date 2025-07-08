import os

OUTPUT_FILE = 'navigation-index/TOC.md'
ROOT_DIR = 'knowledge-system'


def scan_dir(path, prefix=''):
    entries = []
    for name in sorted(os.listdir(path)):
        full_path = os.path.join(path, name)
        if os.path.isdir(full_path):
            entries.append(f'{prefix}- **{name}/**')
            entries.extend(scan_dir(full_path, prefix + '  '))
        else:
            entries.append(f'{prefix}- {name}')
    return entries


def main():
    toc_lines = ['# 知识体系目录索引 (TOC)', '']
    toc_lines += scan_dir(ROOT_DIR)
    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        f.write('\n'.join(toc_lines))
    print(f'已生成目录索引: {OUTPUT_FILE}')

if __name__ == '__main__':
    main() 