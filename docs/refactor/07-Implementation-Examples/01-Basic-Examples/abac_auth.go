package main

import (
	"fmt"
	"net/http"
)

// 用户属性结构体
type User struct {
	Username string
	Role     string
	Dept     string
	Level    int
}

// 资源属性结构体
type Resource struct {
	Type   string
	Owner  string
	Secret bool
}

// ABAC策略决策函数
func abacDecision(user *User, res *Resource, action string) bool {
	// 示例策略：
	// 1. admin可访问所有资源
	if user.Role == "admin" {
		return true
	}
	// 2. 只有本部门且级别大于2的用户可读secret资源
	if action == "read" && res.Secret && user.Dept == res.Owner && user.Level > 2 {
		return true
	}
	// 3. 普通用户可读非secret资源
	if action == "read" && !res.Secret {
		return true
	}
	return false
}

// ABAC中间件
func abacMiddleware(user *User, res *Resource, action string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !abacDecision(user, res, action) {
			http.Error(w, "ABAC拒绝访问", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 示例处理器
func secretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "访问机密资源成功")
}

func main() {
	user := &User{Username: "bob", Role: "user", Dept: "dev", Level: 3}
	res := &Resource{Type: "doc", Owner: "dev", Secret: true}
	mux := http.NewServeMux()
	mux.Handle("/secret", abacMiddleware(user, res, "read", http.HandlerFunc(secretHandler)))
	fmt.Println("服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
