package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// 审计事件结构体
type AuditEvent struct {
	User      string
	Action    string
	Target    string
	Success   bool
	Timestamp time.Time
	Detail    string
}

// 审计日志记录器
type AuditLogger struct {
	logger *log.Logger
	file   *os.File
}

func NewAuditLogger(filename string) (*AuditLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "[AUDIT] ", log.LstdFlags)
	return &AuditLogger{logger: logger, file: file}, nil
}

func (a *AuditLogger) LogEvent(event AuditEvent) {
	a.logger.Printf("user=%s action=%s target=%s success=%v detail=%s time=%s",
		event.User, event.Action, event.Target, event.Success, event.Detail, event.Timestamp.Format(time.RFC3339))
}

func (a *AuditLogger) Close() {
	a.file.Close()
}

// 简单异常检测：检测失败操作
func detectAbnormalEvents(filename string) ([]AuditEvent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var abnormals []AuditEvent
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if n == 0 || err != nil {
			break
		}
		lines := strings.Split(string(buf[:n]), "\n")
		for _, line := range lines {
			if strings.Contains(line, "success=false") {
				abnormals = append(abnormals, AuditEvent{Detail: line})
			}
		}
	}
	return abnormals, nil
}

func main() {
	logger, err := NewAuditLogger("security_audit.log")
	if err != nil {
		fmt.Println("日志初始化失败:", err)
		return
	}
	defer logger.Close()

	// 记录正常和异常操作
	logger.LogEvent(AuditEvent{
		User:      "alice",
		Action:    "login",
		Target:    "system",
		Success:   true,
		Timestamp: time.Now(),
		Detail:    "登录成功",
	})
	logger.LogEvent(AuditEvent{
		User:      "bob",
		Action:    "delete_file",
		Target:    "confidential.txt",
		Success:   false,
		Timestamp: time.Now(),
		Detail:    "权限不足，删除失败",
	})

	// 检测异常事件
	abnormals, err := detectAbnormalEvents("security_audit.log")
	if err != nil {
		fmt.Println("异常检测失败:", err)
		return
	}
	fmt.Printf("检测到%d个异常事件\n", len(abnormals))
	for _, e := range abnormals {
		fmt.Println("异常详情:", e.Detail)
	}
}
