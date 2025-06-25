package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 日志聚合器
func aggregateLogs(files []string) ([]string, error) {
	var logs []string
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			logs = append(logs, scanner.Text())
		}
		f.Close()
	}
	return logs, nil
}

// 简单日志检索
func searchLogs(logs []string, keyword string) []string {
	var result []string
	for _, log := range logs {
		if strings.Contains(log, keyword) {
			result = append(result, log)
		}
	}
	return result
}

func main() {
	// 假设有两个日志文件
	files := []string{"app1.log", "app2.log"}
	logs, err := aggregateLogs(files)
	if err != nil {
		fmt.Println("日志聚合失败:", err)
		return
	}
	fmt.Printf("共聚合日志%d条\n", len(logs))

	// 检索包含"ERROR"的日志
	result := searchLogs(logs, "ERROR")
	fmt.Printf("包含ERROR的日志共%d条:\n", len(result))
	for _, line := range result {
		fmt.Println(line)
	}
}
