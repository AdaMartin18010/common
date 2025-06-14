# 01-Web应用 (Web Application)

## 1. 概述

### 1.1 Web应用基础

**Web应用** 是基于HTTP协议的网络应用程序，Go语言通过标准库`net/http`提供强大的Web开发支持。

**架构模式**：
- **MVC模式**：Model-View-Controller
- **RESTful API**：Representational State Transfer
- **微服务架构**：Microservices Architecture

### 1.2 Go Web开发优势

1. **高性能**：原生HTTP服务器
2. **简洁性**：标准库支持
3. **并发性**：goroutine处理请求
4. **类型安全**：强类型系统

## 2. HTTP服务器

### 2.1 基础HTTP服务器

```go
package webapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// BasicServer 基础HTTP服务器
func BasicServer() {
	// 定义处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})
	
	// 启动服务器
	log.Println("服务器启动在 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CustomServer 自定义服务器
func CustomServer() {
	// 创建多路复用器
	mux := http.NewServeMux()
	
	// 注册路由
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/posts", handlePosts)
	
	// 创建服务器配置
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Println("自定义服务器启动在 :8080")
	log.Fatal(server.ListenAndServe())
}

// handleUsers 用户处理器
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// handlePosts 文章处理器
func handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPosts(w, r)
	case http.MethodPost:
		createPost(w, r)
	default:
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}
```

### 2.2 泛型HTTP处理器

```go
// GenericHandler 泛型处理器
type GenericHandler[T any] struct {
	service Service[T]
}

// Service 服务接口
type Service[T any] interface {
	GetAll() ([]T, error)
	GetByID(id string) (T, error)
	Create(item T) error
	Update(id string, item T) error
	Delete(id string) error
}

// NewGenericHandler 创建泛型处理器
func NewGenericHandler[T any](service Service[T]) *GenericHandler[T] {
	return &GenericHandler[T]{
		service: service,
	}
}

// HandleGetAll 处理获取所有请求
func (h *GenericHandler[T]) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// HandleGetByID 处理根据ID获取请求
func (h *GenericHandler[T]) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "缺少ID参数", http.StatusBadRequest)
		return
	}
	
	item, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// HandleCreate 处理创建请求
func (h *GenericHandler[T]) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	if err := h.service.Create(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}
```

## 3. 路由系统

### 3.1 自定义路由

```go
// Router 路由器
type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

// NewRouter 创建路由器
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// Register 注册路由
func (r *Router) Register(method, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

// ServeHTTP 实现http.Handler接口
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handlers, exists := r.routes[req.Method]; exists {
		if handler, exists := handlers[req.URL.Path]; exists {
			handler(w, req)
			return
		}
	}
	
	http.NotFound(w, req)
}

// GET 注册GET路由
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Register(http.MethodGet, path, handler)
}

// POST 注册POST路由
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Register(http.MethodPost, path, handler)
}

// PUT 注册PUT路由
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Register(http.MethodPut, path, handler)
}

// DELETE 注册DELETE路由
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Register(http.MethodDelete, path, handler)
}
```

### 3.2 参数化路由

```go
// ParamRouter 参数化路由器
type ParamRouter struct {
	routes map[string]map[string]*Route
}

// Route 路由定义
type Route struct {
	pattern string
	handler http.HandlerFunc
	params  []string
}

// NewParamRouter 创建参数化路由器
func NewParamRouter() *ParamRouter {
	return &ParamRouter{
		routes: make(map[string]map[string]*Route),
	}
}

// Register 注册参数化路由
func (pr *ParamRouter) Register(method, pattern string, handler http.HandlerFunc) {
	if pr.routes[method] == nil {
		pr.routes[method] = make(map[string]*Route)
	}
	
	params := extractParams(pattern)
	pr.routes[method][pattern] = &Route{
		pattern: pattern,
		handler: handler,
		params:  params,
	}
}

// extractParams 提取参数
func extractParams(pattern string) []string {
	var params []string
	// 简化实现，实际应该使用正则表达式
	return params
}

// ServeHTTP 实现http.Handler接口
func (pr *ParamRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if routes, exists := pr.routes[req.Method]; exists {
		for pattern, route := range routes {
			if matchPattern(req.URL.Path, pattern) {
				route.handler(w, req)
				return
			}
		}
	}
	
	http.NotFound(w, req)
}

// matchPattern 匹配模式
func matchPattern(path, pattern string) bool {
	// 简化实现，实际应该使用正则表达式
	return path == pattern
}
```

## 4. 中间件

### 4.1 中间件基础

```go
// Middleware 中间件类型
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain 中间件链
func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 调用下一个处理器
		next(w, r)
		
		// 记录日志
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}
		
		// 验证token（简化实现）
		if !validateToken(token) {
			http.Error(w, "无效token", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}

// validateToken 验证token
func validateToken(token string) bool {
	// 简化实现
	return token != ""
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limit int, window time.Duration) Middleware {
	requests := make(map[string][]time.Time)
	
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			now := time.Now()
			
			// 清理过期请求
			if times, exists := requests[clientIP]; exists {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) < window {
						valid = append(valid, t)
					}
				}
				requests[clientIP] = valid
			}
			
			// 检查限流
			if len(requests[clientIP]) >= limit {
				http.Error(w, "请求过于频繁", http.StatusTooManyRequests)
				return
			}
			
			// 记录请求
			requests[clientIP] = append(requests[clientIP], now)
			
			next(w, r)
		}
	}
}
```

### 4.2 泛型中间件

```go
// GenericMiddleware 泛型中间件
type GenericMiddleware[T any] func(http.HandlerFunc, T) http.HandlerFunc

// ContextMiddleware 上下文中间件
func ContextMiddleware[T any](key string, value T) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, value)
			next(w, r.WithContext(ctx))
		}
	}
}

// ValidationMiddleware 验证中间件
func ValidationMiddleware[T any](validator func(T) error) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var data T
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "无效的请求体", http.StatusBadRequest)
				return
			}
			
			if err := validator(data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			// 将验证后的数据存储到上下文中
			ctx := context.WithValue(r.Context(), "validated_data", data)
			next(w, r.WithContext(ctx))
		}
	}
}
```

## 5. 数据模型

### 5.1 基础模型

```go
// User 用户模型
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Post 文章模型
type Post struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// Response 响应模型
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse 分页响应
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
```

### 5.2 泛型模型

```go
// Model 基础模型接口
type Model interface {
	GetID() string
	SetID(id string)
	GetCreated() time.Time
	SetCreated(t time.Time)
	GetUpdated() time.Time
	SetUpdated(t time.Time)
}

// BaseModel 基础模型实现
type BaseModel struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// GetID 获取ID
func (bm *BaseModel) GetID() string {
	return bm.ID
}

// SetID 设置ID
func (bm *BaseModel) SetID(id string) {
	bm.ID = id
}

// GetCreated 获取创建时间
func (bm *BaseModel) GetCreated() time.Time {
	return bm.Created
}

// SetCreated 设置创建时间
func (bm *BaseModel) SetCreated(t time.Time) {
	bm.Created = t
}

// GetUpdated 获取更新时间
func (bm *BaseModel) GetUpdated() time.Time {
	return bm.Updated
}

// SetUpdated 设置更新时间
func (bm *BaseModel) SetUpdated(t time.Time) {
	bm.Updated = t
}

// GenericUser 泛型用户模型
type GenericUser struct {
	BaseModel
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenericPost 泛型文章模型
type GenericPost struct {
	BaseModel
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
```

## 6. 服务层

### 6.1 基础服务

```go
// UserService 用户服务
type UserService struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

// GetAll 获取所有用户
func (us *UserService) GetAll() ([]*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	users := make([]*User, 0, len(us.users))
	for _, user := range us.users {
		users = append(users, user)
	}
	
	return users, nil
}

// GetByID 根据ID获取用户
func (us *UserService) GetByID(id string) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	user, exists := us.users[id]
	if !exists {
		return nil, fmt.Errorf("用户不存在: %s", id)
	}
	
	return user, nil
}

// Create 创建用户
func (us *UserService) Create(user *User) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	if _, exists := us.users[user.ID]; exists {
		return fmt.Errorf("用户已存在: %s", user.ID)
	}
	
	user.Created = time.Now()
	user.Updated = time.Now()
	us.users[user.ID] = user
	
	return nil
}

// Update 更新用户
func (us *UserService) Update(id string, user *User) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	if _, exists := us.users[id]; !exists {
		return fmt.Errorf("用户不存在: %s", id)
	}
	
	user.ID = id
	user.Updated = time.Now()
	us.users[id] = user
	
	return nil
}

// Delete 删除用户
func (us *UserService) Delete(id string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	if _, exists := us.users[id]; !exists {
		return fmt.Errorf("用户不存在: %s", id)
	}
	
	delete(us.users, id)
	return nil
}
```

### 6.2 泛型服务

```go
// GenericService 泛型服务
type GenericService[T Model] struct {
	items map[string]T
	mu    sync.RWMutex
}

// NewGenericService 创建泛型服务
func NewGenericService[T Model]() *GenericService[T] {
	return &GenericService[T]{
		items: make(map[string]T),
	}
}

// GetAll 获取所有项目
func (gs *GenericService[T]) GetAll() ([]T, error) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	
	items := make([]T, 0, len(gs.items))
	for _, item := range gs.items {
		items = append(items, item)
	}
	
	return items, nil
}

// GetByID 根据ID获取项目
func (gs *GenericService[T]) GetByID(id string) (T, error) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	
	item, exists := gs.items[id]
	if !exists {
		var zero T
		return zero, fmt.Errorf("项目不存在: %s", id)
	}
	
	return item, nil
}

// Create 创建项目
func (gs *GenericService[T]) Create(item T) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	id := item.GetID()
	if _, exists := gs.items[id]; exists {
		return fmt.Errorf("项目已存在: %s", id)
	}
	
	now := time.Now()
	item.SetCreated(now)
	item.SetUpdated(now)
	gs.items[id] = item
	
	return nil
}

// Update 更新项目
func (gs *GenericService[T]) Update(id string, item T) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	if _, exists := gs.items[id]; !exists {
		return fmt.Errorf("项目不存在: %s", id)
	}
	
	item.SetID(id)
	item.SetUpdated(time.Now())
	gs.items[id] = item
	
	return nil
}

// Delete 删除项目
func (gs *GenericService[T]) Delete(id string) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	
	if _, exists := gs.items[id]; !exists {
		return fmt.Errorf("项目不存在: %s", id)
	}
	
	delete(gs.items, id)
	return nil
}
```

## 7. 完整应用示例

```go
// WebApp Web应用
type WebApp struct {
	router *Router
	userService *UserService
	postService *GenericService[*GenericPost]
}

// NewWebApp 创建Web应用
func NewWebApp() *WebApp {
	app := &WebApp{
		router:      NewRouter(),
		userService: NewUserService(),
		postService: NewGenericService[*GenericPost](),
	}
	
	app.setupRoutes()
	return app
}

// setupRoutes 设置路由
func (app *WebApp) setupRoutes() {
	// 用户路由
	app.router.GET("/api/users", Chain(
		app.handleGetUsers,
		LoggingMiddleware,
		CORSMiddleware,
	))
	
	app.router.POST("/api/users", Chain(
		app.handleCreateUser,
		LoggingMiddleware,
		CORSMiddleware,
		ValidationMiddleware[User](app.validateUser),
	))
	
	// 文章路由
	app.router.GET("/api/posts", Chain(
		app.handleGetPosts,
		LoggingMiddleware,
		CORSMiddleware,
	))
	
	app.router.POST("/api/posts", Chain(
		app.handleCreatePost,
		LoggingMiddleware,
		CORSMiddleware,
		AuthMiddleware,
		ValidationMiddleware[*GenericPost](app.validatePost),
	))
}

// handleGetUsers 处理获取用户请求
func (app *WebApp) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.userService.GetAll()
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, Response[[]*User]{
		Success: true,
		Data:    users,
	})
}

// handleCreateUser 处理创建用户请求
func (app *WebApp) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	if err := app.userService.Create(&user); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, Response[*User]{
		Success: true,
		Data:    &user,
	})
}

// handleGetPosts 处理获取文章请求
func (app *WebApp) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.postService.GetAll()
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, Response[[]*GenericPost]{
		Success: true,
		Data:    posts,
	})
}

// handleCreatePost 处理创建文章请求
func (app *WebApp) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post GenericPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		respondWithError(w, "无效的请求体", http.StatusBadRequest)
		return
	}
	
	if err := app.postService.Create(&post); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondWithJSON(w, Response[*GenericPost]{
		Success: true,
		Data:    &post,
	})
}

// validateUser 验证用户
func (app *WebApp) validateUser(user User) error {
	if user.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if user.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	return nil
}

// validatePost 验证文章
func (app *WebApp) validatePost(post *GenericPost) error {
	if post.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if post.Content == "" {
		return fmt.Errorf("内容不能为空")
	}
	return nil
}

// respondWithJSON 响应JSON
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// respondWithError 响应错误
func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, Response[interface{}]{
		Success: false,
		Error:   message,
	})
	w.WriteHeader(code)
}

// Run 运行应用
func (app *WebApp) Run(addr string) error {
	log.Printf("Web应用启动在 %s", addr)
	return http.ListenAndServe(addr, app.router)
}
```

## 8. 总结

### 8.1 最佳实践

1. **分层架构**：清晰分离路由、服务、模型层
2. **中间件链**：使用中间件处理横切关注点
3. **错误处理**：统一的错误响应格式
4. **类型安全**：使用泛型提高代码复用性
5. **并发安全**：使用适当的同步机制

### 8.2 性能优化

1. **连接池**：复用HTTP连接
2. **缓存**：缓存频繁访问的数据
3. **压缩**：启用gzip压缩
4. **限流**：防止过载

### 8.3 安全考虑

1. **输入验证**：验证所有用户输入
2. **认证授权**：实现适当的认证机制
3. **HTTPS**：使用TLS加密
4. **CORS**：正确配置跨域策略

---

**参考文献**：
1. Go官方文档：net/http包
2. Go Web开发最佳实践
3. RESTful API设计指南 