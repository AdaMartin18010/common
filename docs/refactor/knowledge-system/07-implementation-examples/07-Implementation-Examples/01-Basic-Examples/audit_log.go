package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 安全日志记录器
type AuditLogger struct {
	logger *log.Logger
	file   *os.File
}

// 创建新的日志记录器
func NewAuditLogger(filename string) (*AuditLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "[AUDIT] ", log.LstdFlags|log.Lshortfile)
	return &AuditLogger{logger: logger, file: file}, nil
}

// 记录敏感操作
type AuditEvent struct {
	User      string
	Action    string
	Target    string
	Timestamp time.Time
	Success   bool
	Detail    string
}

func (a *AuditLogger) LogEvent(event AuditEvent) {
	a.logger.Printf("user=%s action=%s target=%s success=%v detail=%s time=%s",
		event.User, event.Action, event.Target, event.Success, event.Detail, event.Timestamp.Format(time.RFC3339))
}

// 日志轮转（简单示例）
func (a *AuditLogger) Rotate(newFilename string) error {
	a.file.Close()
	file, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	a.logger.SetOutput(file)
	a.file = file
	return nil
}

func main() {
	logger, err := NewAuditLogger("audit.log")
	if err != nil {
		fmt.Println("日志初始化失败:", err)
		return
	}
	defer logger.file.Close()

	// 记录敏感操作
	event := AuditEvent{
		User:      "alice",
		Action:    "delete_user",
		Target:    "bob",
		Timestamp: time.Now(),
		Success:   true,
		Detail:    "管理员删除用户bob",
	}
	logger.LogEvent(event)

	fmt.Println("敏感操作已记录到audit.log")
}
