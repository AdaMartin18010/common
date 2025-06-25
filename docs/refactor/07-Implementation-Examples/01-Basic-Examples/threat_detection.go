package main

import (
	"fmt"
	"net/http"
	"strings"
)

// 黑名单IP列表
var blacklist = map[string]bool{
	"192.168.1.100": true,
	"10.0.0.5":      true,
}

// 简单恶意请求特征
var maliciousPatterns = []string{
	"/etc/passwd",
	"select * from",
	"<script>",
}

// 入侵检测中间件
func intrusionDetectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		for k := range blacklist {
			if strings.HasPrefix(ip, k) {
				fmt.Println("[ALERT] 黑名单IP访问:", ip)
				http.Error(w, "禁止访问", http.StatusForbidden)
				return
			}
		}
		for _, pattern := range maliciousPatterns {
			if strings.Contains(r.RequestURI, pattern) || strings.Contains(r.URL.RawQuery, pattern) {
				fmt.Println("[ALERT] 检测到恶意请求:", r.RequestURI)
				http.Error(w, "检测到恶意请求", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// 示例处理器
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "欢迎访问安全系统")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	fmt.Println("威胁检测服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", intrusionDetectionMiddleware(mux))
}
