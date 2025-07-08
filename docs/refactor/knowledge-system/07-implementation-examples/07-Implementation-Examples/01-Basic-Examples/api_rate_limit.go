package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 令牌桶限流器
type TokenBucket struct {
	capacity   int
	tokens     int
	rate       int // 每秒生成令牌数
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	newTokens := int(elapsed * float64(tb.rate))
	if newTokens > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+newTokens)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 限流中间件
func rateLimitMiddleware(tb *TokenBucket, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.Allow() {
			http.Error(w, "请求过多，请稍后再试", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 示例处理器
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API请求成功")
}

func main() {
	tb := NewTokenBucket(5, 2) // 容量5，每秒2个令牌
	mux := http.NewServeMux()
	mux.Handle("/api", rateLimitMiddleware(tb, http.HandlerFunc(apiHandler)))
	fmt.Println("API限流服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
