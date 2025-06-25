package main

import (
	"fmt"
	"net/http"
)

// 角色与权限定义
var rolePermissions = map[string][]string{
	"admin":  {"read", "write", "delete"},
	"editor": {"read", "write"},
	"viewer": {"read"},
}

// 用户结构体
type User struct {
	Username string
	Role     string
}

// 检查用户是否有权限
func hasPermission(user *User, perm string) bool {
	perms, ok := rolePermissions[user.Role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// RBAC中间件
func rbacMiddleware(perm string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟从请求获取用户
		user := &User{Username: "alice", Role: "editor"}
		if !hasPermission(user, perm) {
			http.Error(w, "无权限", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 示例处理器
func writeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "写入操作成功")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/write", rbacMiddleware("write", http.HandlerFunc(writeHandler)))
	fmt.Println("服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
