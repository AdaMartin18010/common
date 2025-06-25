package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

// 生成CSRF Token
func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// 校验CSRF Token
func validateCSRFToken(r *http.Request, token string) bool {
	formToken := r.FormValue("csrf_token")
	return formToken == token
}

// CSRF中间件
func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken, err := generateCSRFToken()
		if err != nil {
			http.Error(w, "CSRF Token生成失败", http.StatusInternalServerError)
			return
		}
		// 设置token到cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "csrf_token",
			Value: csrfToken,
			Path:  "/",
		})
		r.ParseForm()
		if r.Method == http.MethodPost {
			if !validateCSRFToken(r, csrfToken) {
				http.Error(w, "CSRF校验失败", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// 示例处理器
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, CSRF Protection!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	fmt.Println("服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", csrfMiddleware(mux))
}
