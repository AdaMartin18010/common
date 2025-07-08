import os
import re
from collections import defaultdict

ROOT_DIR = 'knowledge-system'
VISUALIZATION_FILE = 'navigation-index/knowledge-visualization.md'

# 主题分类映射
TOPIC_MAPPING = {
    '01-foundation-theory': '基础理论',
    '02-software-architecture': '软件架构', 
    '03-design-patterns': '设计模式',
    '04-programming-languages': '编程语言',
    '05-industry-domains': '行业领域',
    '06-formal-methods': '形式化方法',
    '07-implementation-examples': '实现示例',
    '08-software-engineering-formalization': '软件工程形式化',
    '09-programming-language-theory': '编程语言理论',
    '10-workflow-systems': '工作流系统',
    '11-advanced-topics': '高级主题',
    '12-international-standards': '国际标准',
    '13-quality-assurance': '质量保证'
}

def extract_topics(content):
    """提取内容中的主题标签"""
    topics = []
    # 提取 #标签 格式
    topic_matches = re.findall(r'#(\w+)', content)
    topics.extend(topic_matches)
    # 提取关键词
    keywords = ['架构', '设计', '模式', '理论', '实现', '方法', '系统', '语言', '工程']
    for keyword in keywords:
        if keyword in content:
            topics.append(keyword)
    return list(set(topics))

def build_knowledge_graph():
    """构建知识图谱"""
    graph = defaultdict(list)
    topic_files = defaultdict(list)
    
    # 扫描所有md文件
    for root, dirs, files in os.walk(ROOT_DIR):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                rel_path = os.path.relpath(file_path, ROOT_DIR)
                
                # 确定主题分类
                topic = None
                for key, value in TOPIC_MAPPING.items():
                    if key in rel_path:
                        topic = value
                        break
                
                if topic:
                    topic_files[topic].append(rel_path)
                
                # 提取主题标签
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                    topics = extract_topics(content)
                    for t in topics:
                        graph[t].append(rel_path)
                except:
                    pass
    
    return graph, topic_files

def generate_visualization():
    """生成可视化内容"""
    graph, topic_files = build_knowledge_graph()
    
    lines = ['# 知识体系可视化图谱', '']
    
    # 主题分类统计
    lines.append('## 主题分类统计')
    for topic, files in topic_files.items():
        lines.append(f'### {topic}')
        lines.append(f'- 文件数量: {len(files)}')
        for file in files[:5]:  # 只显示前5个文件
            lines.append(f'  - {file}')
        if len(files) > 5:
            lines.append(f'  - ... 还有 {len(files)-5} 个文件')
        lines.append('')
    
    # 主题关联网络
    lines.append('## 主题关联网络')
    for topic, files in graph.items():
        if len(files) > 1:  # 只显示有关联的主题
            lines.append(f'### #{topic}')
            lines.append(f'- 关联文件数: {len(files)}')
            for file in files[:3]:  # 只显示前3个文件
                lines.append(f'  - {file}')
            if len(files) > 3:
                lines.append(f'  - ... 还有 {len(files)-3} 个文件')
            lines.append('')
    
    # 知识层级关系
    lines.append('## 知识层级关系')
    lines.append('```mermaid')
    lines.append('graph TD')
    for topic in TOPIC_MAPPING.values():
        lines.append(f'    A[{topic}]')
    lines.append('```')
    
    with open(VISUALIZATION_FILE, 'w', encoding='utf-8') as f:
        f.write('\n'.join(lines))
    
    print(f'已生成知识图谱可视化: {VISUALIZATION_FILE}')

if __name__ == '__main__':
    generate_visualization() 