package main

import (
	"fmt"
	"net/http"
	"sync"
)

// 封禁IP列表（线程安全）
var (
	bannedIPs = make(map[string]bool)
	mu        sync.Mutex
)

// 自动封禁IP
func banIP(ip string) {
	mu.Lock()
	defer mu.Unlock()
	bannedIPs[ip] = true
	fmt.Println("[AUTO-RESPONSE] 已自动封禁IP:", ip)
}

// 检查IP是否被封禁
func isBanned(ip string) bool {
	mu.Lock()
	defer mu.Unlock()
	return bannedIPs[ip]
}

// 模拟发送告警（实际可对接邮件、钉钉等）
func sendAlert(msg string) {
	fmt.Println("[ALERT] 发送告警:", msg)
}

// 自动化响应中间件
func autoResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if isBanned(ip) {
			http.Error(w, "您的IP已被封禁", http.StatusForbidden)
			return
		}
		// 模拟检测到恶意行为
		if r.URL.Query().Get("attack") == "1" {
			banIP(ip)
			sendAlert("检测到恶意行为，已自动封禁IP: " + ip)
			http.Error(w, "检测到恶意行为，已自动响应", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "安全自动化响应系统")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	fmt.Println("自动化响应服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", autoResponseMiddleware(mux))
}
