package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// 模拟Prometheus指标采集
var (
	totalRequests int
	errorRequests int
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "# HELP total_requests 总请求数\n")
	fmt.Fprintf(w, "# TYPE total_requests counter\n")
	fmt.Fprintf(w, "total_requests %d\n", totalRequests)
	fmt.Fprintf(w, "# HELP error_requests 错误请求数\n")
	fmt.Fprintf(w, "# TYPE error_requests counter\n")
	fmt.Fprintf(w, "error_requests %d\n", errorRequests)
}

// 实时日志监控与异常告警
type Monitor struct {
	logger *log.Logger
	file   *os.File
}

func NewMonitor(filename string) (*Monitor, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "[MONITOR] ", log.LstdFlags)
	return &Monitor{logger: logger, file: file}, nil
}

func (m *Monitor) LogEvent(event string, isError bool) {
	m.logger.Printf("event=%s error=%v time=%s", event, isError, time.Now().Format(time.RFC3339))
	if isError {
		fmt.Println("[ALERT] 检测到异常事件:", event)
	}
}

func (m *Monitor) Close() {
	m.file.Close()
}

func main() {
	monitor, err := NewMonitor("security_monitor.log")
	if err != nil {
		fmt.Println("监控日志初始化失败:", err)
		return
	}
	defer monitor.Close()

	http.HandleFunc("/metrics", metricsHandler)

	// 模拟业务请求与监控
	for i := 0; i < 10; i++ {
		totalRequests++
		isError := rand.Intn(10) < 2 // 20%概率为异常
		if isError {
			errorRequests++
			monitor.LogEvent(fmt.Sprintf("请求%d发生错误", i+1), true)
		} else {
			monitor.LogEvent(fmt.Sprintf("请求%d正常", i+1), false)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("安全监控日志已生成，Prometheus指标接口: http://localhost:8080/metrics")
	http.ListenAndServe(":8080", nil)
}
