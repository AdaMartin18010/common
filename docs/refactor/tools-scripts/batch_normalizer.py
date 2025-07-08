import os
import re

ROOT_DIR = 'knowledge-system'

# 规范化内容模板
SUMMARY = '（本条目自动生成摘要：请补充该知识点的核心内容和意义）'
DETAIL = '- 背景与定义：\n- 关键概念：\n- 相关原理：\n- 实践应用：\n- 典型案例：\n- 拓展阅读：'
REFERENCES = '- [示例参考文献1](#)\n- [示例参考文献2](#)'
TAGS = '- #待补充 #知识点 #标签'

# 规范化处理
for root, dirs, files in os.walk(ROOT_DIR):
    for file in files:
        if file == 'README.md':
            file_path = os.path.join(root, file)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    content = f.read()
                changed = False
                # 补充摘要
                if '## 摘要' not in content:
                    content = content.replace('## 标题', '## 标题\n')
                    content = content.replace('## 摘要', f'## 摘要\n{SUMMARY}') if '## 摘要' in content else content.replace('（简要描述该知识点的核心内容和意义）', SUMMARY)
                    if '## 摘要' not in content:
                        content = content.replace('## 标题', f'## 标题\n\n## 摘要\n{SUMMARY}\n')
                    changed = True
                # 补充详细内容
                if '## 详细内容' not in content:
                    if '## 摘要' in content:
                        content = content.replace('## 摘要', f'## 摘要\n\n## 详细内容\n{DETAIL}')
                    else:
                        content += f'\n## 详细内容\n{DETAIL}'
                    changed = True
                # 补充参考文献
                if '## 参考文献' not in content:
                    content += f'\n\n## 参考文献\n{REFERENCES}'
                    changed = True
                # 补充标签
                if '## 标签' not in content:
                    content += f'\n\n## 标签\n{TAGS}'
                    changed = True
                if changed:
                    with open(file_path, 'w', encoding='utf-8') as f:
                        f.write(content)
                    print(f'✅ 已规范化: {file_path}')
            except Exception as e:
                print(f'⚠️ 规范化失败: {file_path} ({e})') 