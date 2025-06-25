#!/usr/bin/env python3
"""
代码示例性能优化工具

该工具用于分析和优化Go代码示例的性能，包括：
1. 性能基准测试生成
2. 内存使用分析
3. 并发安全性检查
4. 错误处理完善
"""

import os
import re
import json
import subprocess
from pathlib import Path
from typing import List, Dict, Any, Optional
import logging

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class GoCodeOptimizer:
    """Go代码优化器"""
    
    def __init__(self, docs_path: str):
        self.docs_path = Path(docs_path)
        self.optimization_results = {}
        
    def find_go_files(self) -> List[Path]:
        """查找所有Go代码文件"""
        go_files = []
        for root, dirs, files in os.walk(self.docs_path):
            for file in files:
                if file.endswith('.go'):
                    go_files.append(Path(root) / file)
        return go_files
    
    def analyze_code_performance(self, file_path: Path) -> Dict[str, Any]:
        """分析代码性能"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            analysis = {
                'file': str(file_path),
                'lines': len(content.split('\n')),
                'functions': self._count_functions(content),
                'goroutines': self._count_goroutines(content),
                'channels': self._count_channels(content),
                'mutexes': self._count_mutexes(content),
                'error_handling': self._analyze_error_handling(content),
                'memory_usage': self._analyze_memory_usage(content),
                'concurrency_safety': self._analyze_concurrency_safety(content),
                'performance_issues': self._detect_performance_issues(content)
            }
            
            return analysis
            
        except Exception as e:
            logger.error(f"分析文件 {file_path} 时出错: {e}")
            return {'file': str(file_path), 'error': str(e)}
    
    def _count_functions(self, content: str) -> int:
        """统计函数数量"""
        pattern = r'func\s+\w+\s*\([^)]*\)'
        return len(re.findall(pattern, content))
    
    def _count_goroutines(self, content: str) -> int:
        """统计goroutine数量"""
        pattern = r'go\s+\w+\s*\('
        return len(re.findall(pattern, content))
    
    def _count_channels(self, content: str) -> int:
        """统计channel数量"""
        pattern = r'make\s*\(\s*chan\s+'
        return len(re.findall(pattern, content))
    
    def _count_mutexes(self, content: str) -> int:
        """统计互斥锁数量"""
        pattern = r'sync\.(Mutex|RWMutex)'
        return len(re.findall(pattern, content))
    
    def _analyze_error_handling(self, content: str) -> Dict[str, Any]:
        """分析错误处理"""
        analysis = {
            'has_error_returns': bool(re.search(r'error\s*\)', content)),
            'has_error_checks': bool(re.search(r'if\s+err\s*!=', content)),
            'has_panic_handling': bool(re.search(r'defer\s+recover', content)),
            'error_handling_score': 0
        }
        
        score = 0
        if analysis['has_error_returns']:
            score += 1
        if analysis['has_error_checks']:
            score += 1
        if analysis['has_panic_handling']:
            score += 1
            
        analysis['error_handling_score'] = score
        return analysis
    
    def _analyze_memory_usage(self, content: str) -> Dict[str, Any]:
        """分析内存使用"""
        analysis = {
            'has_memory_allocation': bool(re.search(r'make\s*\(', content)),
            'has_garbage_collection': bool(re.search(r'runtime\.GC', content)),
            'has_memory_profiling': bool(re.search(r'pprof\.', content)),
            'potential_memory_leaks': self._detect_memory_leaks(content)
        }
        return analysis
    
    def _detect_memory_leaks(self, content: str) -> List[str]:
        """检测潜在的内存泄漏"""
        leaks = []
        
        # 检查未关闭的资源
        if re.search(r'os\.Open', content) and not re.search(r'\.Close\(\)', content):
            leaks.append("文件可能未正确关闭")
        
        if re.search(r'http\.Get', content) and not re.search(r'\.Close\(\)', content):
            leaks.append("HTTP响应可能未正确关闭")
        
        # 检查goroutine泄漏
        if re.search(r'go\s+\w+\s*\(', content) and not re.search(r'context\.Cancel', content):
            leaks.append("goroutine可能没有正确的取消机制")
            
        return leaks
    
    def _analyze_concurrency_safety(self, content: str) -> Dict[str, Any]:
        """分析并发安全性"""
        analysis = {
            'has_race_conditions': self._detect_race_conditions(content),
            'has_deadlock_risk': self._detect_deadlock_risk(content),
            'has_proper_synchronization': self._has_proper_synchronization(content),
            'concurrency_safety_score': 0
        }
        
        score = 0
        if not analysis['has_race_conditions']:
            score += 1
        if not analysis['has_deadlock_risk']:
            score += 1
        if analysis['has_proper_synchronization']:
            score += 1
            
        analysis['concurrency_safety_score'] = score
        return analysis
    
    def _detect_race_conditions(self, content: str) -> List[str]:
        """检测竞态条件"""
        races = []
        
        # 检查共享变量的并发访问
        if re.search(r'go\s+\w+\s*\(', content):
            # 查找可能的共享变量
            shared_vars = re.findall(r'var\s+(\w+)', content)
            for var in shared_vars:
                if re.search(f'{var}\\s*[+\\-*/=]', content):
                    races.append(f"变量 {var} 可能存在竞态条件")
                    
        return races
    
    def _detect_deadlock_risk(self, content: str) -> List[str]:
        """检测死锁风险"""
        deadlocks = []
        
        # 检查锁的嵌套使用
        if re.search(r'sync\.Mutex', content) and re.search(r'Lock\(\)', content):
            # 检查是否有多个锁
            lock_count = len(re.findall(r'\.Lock\(\)', content))
            unlock_count = len(re.findall(r'\.Unlock\(\)', content))
            
            if lock_count != unlock_count:
                deadlocks.append("锁的获取和释放不匹配")
                
        return deadlocks
    
    def _has_proper_synchronization(self, content: str) -> bool:
        """检查是否有适当的同步机制"""
        has_mutex = bool(re.search(r'sync\.Mutex', content))
        has_channel = bool(re.search(r'make\s*\(\s*chan', content))
        has_context = bool(re.search(r'context\.', content))
        
        return has_mutex or has_channel or has_context
    
    def _detect_performance_issues(self, content: str) -> List[str]:
        """检测性能问题"""
        issues = []
        
        # 检查字符串拼接
        if re.search(r'\\+\\s*["\']', content):
            issues.append("使用+进行字符串拼接，建议使用strings.Builder")
        
        # 检查频繁的内存分配
        if content.count('make(') > 5:
            issues.append("可能存在频繁的内存分配")
        
        # 检查不必要的goroutine
        if re.search(r'go\s+func\s*\(\s*\)\s*\{[^}]*fmt\.Print', content):
            issues.append("在goroutine中使用fmt.Print可能影响性能")
            
        return issues
    
    def generate_benchmark_template(self, file_path: Path) -> str:
        """生成基准测试模板"""
        template = f"""package main

import (
	"testing"
	"time"
)

// Benchmark for {file_path.name}
func Benchmark{file_path.stem}(b *testing.B) {{
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {{
		// TODO: Add your benchmark code here
		// Example:
		// result := YourFunction()
		// _ = result
	}}
}}

// Memory benchmark
func Benchmark{file_path.stem}Memory(b *testing.B) {{
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {{
		// TODO: Add your memory benchmark code here
	}}
}}

// Concurrent benchmark
func Benchmark{file_path.stem}Concurrent(b *testing.B) {{
	b.RunParallel(func(pb *testing.PB) {{
		for pb.Next() {{
			// TODO: Add your concurrent benchmark code here
		}}
	}})
}}
"""
        return template
    
    def generate_optimization_suggestions(self, analysis: Dict[str, Any]) -> List[str]:
        """生成优化建议"""
        suggestions = []
        
        # 错误处理建议
        if analysis.get('error_handling', {}).get('error_handling_score', 0) < 2:
            suggestions.append("建议完善错误处理机制，包括错误返回、检查和恢复")
        
        # 并发安全建议
        if analysis.get('concurrency_safety', {}).get('concurrency_safety_score', 0) < 2:
            suggestions.append("建议加强并发安全性，使用适当的同步机制")
        
        # 内存使用建议
        memory_analysis = analysis.get('memory_usage', {})
        if memory_analysis.get('potential_memory_leaks'):
            suggestions.append("建议检查并修复潜在的内存泄漏问题")
        
        # 性能问题建议
        performance_issues = analysis.get('performance_issues', [])
        for issue in performance_issues:
            suggestions.append(f"性能优化建议: {issue}")
        
        return suggestions
    
    def optimize_all_files(self) -> Dict[str, Any]:
        """优化所有文件"""
        go_files = self.find_go_files()
        logger.info(f"找到 {len(go_files)} 个Go文件")
        
        results = {
            'total_files': len(go_files),
            'analyses': [],
            'optimization_suggestions': [],
            'benchmark_templates': []
        }
        
        for file_path in go_files:
            logger.info(f"分析文件: {file_path}")
            
            # 分析代码
            analysis = self.analyze_code_performance(file_path)
            results['analyses'].append(analysis)
            
            # 生成优化建议
            suggestions = self.generate_optimization_suggestions(analysis)
            results['optimization_suggestions'].extend(suggestions)
            
            # 生成基准测试模板
            benchmark_template = self.generate_benchmark_template(file_path)
            results['benchmark_templates'].append({
                'file': str(file_path),
                'template': benchmark_template
            })
        
        return results
    
    def save_results(self, results: Dict[str, Any], output_file: str):
        """保存结果到文件"""
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        
        logger.info(f"结果已保存到: {output_file}")
    
    def generate_report(self, results: Dict[str, Any]) -> str:
        """生成分析报告"""
        report = []
        report.append("# Go代码性能优化分析报告")
        report.append("")
        
        # 总体统计
        report.append("## 总体统计")
        report.append(f"- 总文件数: {results['total_files']}")
        report.append("")
        
        # 文件分析
        report.append("## 文件分析")
        for analysis in results['analyses']:
            if 'error' in analysis:
                report.append(f"### {analysis['file']}")
                report.append(f"- 错误: {analysis['error']}")
                report.append("")
                continue
                
            report.append(f"### {analysis['file']}")
            report.append(f"- 代码行数: {analysis['lines']}")
            report.append(f"- 函数数量: {analysis['functions']}")
            report.append(f"- Goroutine数量: {analysis['goroutines']}")
            report.append(f"- Channel数量: {analysis['channels']}")
            report.append(f"- 互斥锁数量: {analysis['mutexes']}")
            
            # 错误处理评分
            error_score = analysis['error_handling']['error_handling_score']
            report.append(f"- 错误处理评分: {error_score}/3")
            
            # 并发安全评分
            concurrency_score = analysis['concurrency_safety']['concurrency_safety_score']
            report.append(f"- 并发安全评分: {concurrency_score}/3")
            
            # 性能问题
            if analysis['performance_issues']:
                report.append("- 性能问题:")
                for issue in analysis['performance_issues']:
                    report.append(f"  - {issue}")
            
            report.append("")
        
        # 优化建议
        if results['optimization_suggestions']:
            report.append("## 优化建议")
            for suggestion in results['optimization_suggestions']:
                report.append(f"- {suggestion}")
            report.append("")
        
        return "\n".join(report)

def main():
    """主函数"""
    # 使用当前目录作为docs_path
    docs_path = "."
    
    if not os.path.exists(docs_path):
        logger.error(f"路径不存在: {docs_path}")
        return
    
    optimizer = GoCodeOptimizer(docs_path)
    
    # 执行优化分析
    logger.info("开始代码性能优化分析...")
    results = optimizer.optimize_all_files()
    
    # 保存结果
    output_file = "code_optimization_results.json"
    optimizer.save_results(results, output_file)
    
    # 生成报告
    report = optimizer.generate_report(results)
    report_file = "code_optimization_report.md"
    with open(report_file, 'w', encoding='utf-8') as f:
        f.write(report)
    
    logger.info(f"分析报告已保存到: {report_file}")
    
    # 生成基准测试模板
    benchmark_dir = "benchmarks"
    os.makedirs(benchmark_dir, exist_ok=True)
    
    for template_info in results['benchmark_templates']:
        file_name = Path(template_info['file']).stem + "_benchmark_test.go"
        benchmark_file = os.path.join(benchmark_dir, file_name)
        
        with open(benchmark_file, 'w', encoding='utf-8') as f:
            f.write(template_info['template'])
    
    logger.info(f"基准测试模板已生成到: {benchmark_dir}/")

if __name__ == "__main__":
    main() 