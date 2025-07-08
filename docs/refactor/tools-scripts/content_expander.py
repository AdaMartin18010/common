import os
import re
import shutil

ROOT_DIR = 'knowledge-system'
TEMPLATE_FILE = 'document-templates/knowledge-entry-template.md'

# 知识点关联映射
KNOWLEDGE_MAPPING = {
    '架构': ['微服务架构', '组件架构', '云原生架构', 'IoT架构'],
    '设计': ['创建型模式', '结构型模式', '行为型模式', '并发模式'],
    '理论': ['数学基础', '逻辑基础', '范畴论基础', '计算理论基础'],
    '语言': ['Go语言', '类型系统', '语义理论', '编译器理论'],
    '方法': ['形式化验证', '软件工程形式化', '工作流建模'],
    '系统': ['工作流系统', '分布式系统', '实时系统'],
    '质量': ['测试验证', '性能优化', '安全保证', '国际标准']
}

def find_related_topics(content):
    """基于内容查找相关主题"""
    related = []
    for keyword, topics in KNOWLEDGE_MAPPING.items():
        if keyword in content:
            related.extend(topics)
    return list(set(related))

def expand_content():
    """自动扩充内容"""
    expanded_count = 0
    
    # 扫描所有md文件
    for root, dirs, files in os.walk(ROOT_DIR):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                    
                    # 查找相关主题
                    related_topics = find_related_topics(content)
                    
                    # 为每个相关主题创建链接文件
                    for topic in related_topics:
                        topic_dir = os.path.join(root, topic)
                        topic_file = os.path.join(topic_dir, 'README.md')
                        
                        if not os.path.exists(topic_file):
                            # 创建目录
                            os.makedirs(topic_dir, exist_ok=True)
                            
                            # 复制模板并修改
                            shutil.copyfile(TEMPLATE_FILE, topic_file)
                            
                            # 读取并修改模板内容
                            with open(topic_file, 'r', encoding='utf-8') as f:
                                template_content = f.read()
                            
                            # 替换模板内容
                            new_content = template_content.replace(
                                '（请填写知识点名称）', 
                                topic
                            ).replace(
                                '（简要描述该知识点的核心内容和意义）',
                                f'与 {os.path.basename(root)} 相关的 {topic} 知识点'
                            )
                            
                            with open(topic_file, 'w', encoding='utf-8') as f:
                                f.write(new_content)
                            
                            expanded_count += 1
                            print(f'✅ 已扩充: {topic_file}')
                            
                except Exception as e:
                    print(f'⚠️ 扩充失败: {file_path} ({e})')
    
    print(f'\n📊 自动化内容扩充完成: 成功扩充 {expanded_count} 个知识点')

if __name__ == '__main__':
    expand_content() 