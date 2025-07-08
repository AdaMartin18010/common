# 3. 应用安全 (Application Security)

## 3.1 应用安全基础

### 3.1.1 应用安全定义与分类

应用安全是保护应用程序免受安全威胁的过程，包括Web应用、移动应用、桌面应用等。

**形式化定义**：
设 ```latex
A
``` 为应用程序集合，```latex
T
``` 为威胁集合，```latex
V
``` 为漏洞集合。
应用安全函数 ```latex
S: A \times T \times V \rightarrow \{0,1\}
``` 定义为：
$S(a,t,v) = \begin{cases}
1 & \text{if application } a \text{ is secure against threat } t \text{ and vulnerability } v \\
0 & \text{otherwise}
\end{cases}$

### 3.1.2 应用安全威胁模型

**威胁分类**：
1. **注入攻击**：SQL注入、命令注入、XSS
2. **认证攻击**：暴力破解、会话劫持
3. **授权攻击**：权限提升、越权访问
4. **数据泄露**：敏感信息泄露、配置泄露

**威胁概率计算**：
$```latex
P(Attack) = P(Vulnerability) \times P(Exploit) \times P(Impact)
```$

## 3.2 Web应用安全

### 3.2.1 OWASP Top 10

**2021年OWASP Top 10**：
1. **A01:2021 - 失效的访问控制**
2. **A02:2021 - 加密失败**
3. **A03:2021 - 注入**
4. **A04:2021 - 不安全设计**
5. **A05:2021 - 安全配置错误**
6. **A06:2021 - 易受攻击和过时的组件**
7. **A07:2021 - 身份验证和验证失败**
8. **A08:2021 - 软件和数据完整性故障**
9. **A09:2021 - 安全日志和监控失败**
10. **A10:2021 - 服务器端请求伪造**

### 3.2.2 SQL注入防护

**SQL注入类型**：
- **布尔盲注**：基于布尔值的盲注
- **时间盲注**：基于时间延迟的盲注
- **联合查询注入**：使用UNION的注入
- **堆叠查询注入**：使用分号的注入

**Go语言SQL注入防护**：

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strings"

    _ "github.com/lib/pq"
)

// User 用户模型
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// SecureUserService 安全用户服务
type SecureUserService struct {
    db *sql.DB
}

// NewSecureUserService 创建安全用户服务
func NewSecureUserService(db *sql.DB) *SecureUserService {
    return &SecureUserService{db: db}
}

// ValidateInput 输入验证
func (s *SecureUserService) ValidateInput(input string) bool {
    // 检查SQL注入模式
    sqlPattern := regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript)`)
    if sqlPattern.MatchString(input) {
        return false
    }

    // 检查XSS模式
    xssPattern := regexp.MustCompile(`(?i)(<script|javascript:|vbscript:|onload=|onerror=)`)
    if xssPattern.MatchString(input) {
        return false
    }

    return true
}

// GetUserByID 通过ID获取用户（安全版本）
func (s *SecureUserService) GetUserByID(id int) (*User, error) {
    // 使用参数化查询防止SQL注入
    query := "SELECT id, username, email FROM users WHERE id = $1"
    row := s.db.QueryRow(query, id)

    user := &User{}
    err := row.Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %v", err)
    }

    return user, nil
}

// GetUserByUsername 通过用户名获取用户（安全版本）
func (s *SecureUserService) GetUserByUsername(username string) (*User, error) {
    // 输入验证
    if !s.ValidateInput(username) {
        return nil, fmt.Errorf("invalid input")
    }

    // 使用参数化查询
    query := "SELECT id, username, email FROM users WHERE username = $1"
    row := s.db.QueryRow(query, username)

    user := &User{}
    err := row.Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %v", err)
    }

    return user, nil
}

// CreateUser 创建用户（安全版本）
func (s *SecureUserService) CreateUser(username, email, password string) error {
    // 输入验证
    if !s.ValidateInput(username) || !s.ValidateInput(email) {
        return fmt.Errorf("invalid input")
    }

    // 密码哈希
    hashedPassword, err := s.hashPassword(password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %v", err)
    }

    // 使用参数化查询
    query := "INSERT INTO users (username, email, password_hash) VALUES (```latex
1,
```2, $3)"
    _, err = s.db.Exec(query, username, email, hashedPassword)
    if err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }

    return nil
}

// hashPassword 密码哈希
func (s *SecureUserService) hashPassword(password string) (string, error) {
    // 使用bcrypt进行密码哈希
    import "golang.org/x/crypto/bcrypt"

    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

// SecureHandler 安全处理器
func (s *SecureUserService) SecureHandler(w http.ResponseWriter, r *http.Request) {
    // 设置安全头
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    w.Header().Set("Content-Security-Policy", "default-src 'self'")

    // 获取用户ID参数
    userID := r.URL.Query().Get("id")
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    // 参数验证
    if !s.ValidateInput(userID) {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // 获取用户信息
    user, err := s.GetUserByUsername(userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // 返回用户信息
    fmt.Fprintf(w, "User: %s, Email: %s", user.Username, user.Email)
}
```

### 3.2.3 XSS防护

**XSS类型**：
1. **反射型XSS**：恶意脚本从请求中反射
2. **存储型XSS**：恶意脚本存储在服务器
3. **DOM型XSS**：恶意脚本修改DOM

**Go语言XSS防护**：

```go
package main

import (
    "fmt"
    "html"
    "net/http"
    "regexp"
    "strings"
)

// XSSProtection XSS防护
type XSSProtection struct{}

// NewXSSProtection 创建XSS防护
func NewXSSProtection() *XSSProtection {
    return &XSSProtection{}
}

// SanitizeHTML HTML净化
func (x *XSSProtection) SanitizeHTML(input string) string {
    // HTML实体编码
    sanitized := html.EscapeString(input)

    // 移除危险标签
    dangerousTags := []string{"script", "iframe", "object", "embed", "form"}
    for _, tag := range dangerousTags {
        pattern := regexp.MustCompile(fmt.Sprintf(`(?i)<%s[^>]*>.*?</%s>`, tag, tag))
        sanitized = pattern.ReplaceAllString(sanitized, "")
    }

    return sanitized
}

// ValidateInput 输入验证
func (x *XSSProtection) ValidateInput(input string) bool {
    // 检查XSS模式
    xssPatterns := []string{
        `<script[^>]*>`,
        `javascript:`,
        `vbscript:`,
        `onload=`,
        `onerror=`,
        `onclick=`,
        `onmouseover=`,
    }

    for _, pattern := range xssPatterns {
        matched, _ := regexp.MatchString(pattern, strings.ToLower(input))
        if matched {
            return false
        }
    }

    return true
}

// SecureResponse 安全响应
func (x *XSSProtection) SecureResponse(w http.ResponseWriter, r *http.Request) {
    // 设置安全头
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    w.Header().Set("Content-Security-Policy", "default-src 'self'")

    // 获取用户输入
    userInput := r.URL.Query().Get("input")
    if userInput == "" {
        http.Error(w, "Input is required", http.StatusBadRequest)
        return
    }

    // 输入验证
    if !x.ValidateInput(userInput) {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // HTML净化
    sanitizedInput := x.SanitizeHTML(userInput)

    // 安全输出
    fmt.Fprintf(w, "<h1>Safe Output</h1><p>%s</p>", sanitizedInput)
}
```

## 3.3 移动应用安全

### 3.3.1 移动应用安全威胁

**主要威胁**：
1. **恶意软件**：病毒、木马、间谍软件
2. **数据泄露**：敏感信息泄露
3. **逆向工程**：代码反编译
4. **中间人攻击**：网络通信拦截

### 3.3.2 移动应用安全防护

**防护措施**：
1. **代码混淆**：防止逆向工程
2. **证书固定**：防止中间人攻击
3. **数据加密**：保护敏感数据
4. **权限控制**：最小权限原则

**Go语言移动应用安全**：

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "net/http"
)

// MobileAppSecurity 移动应用安全
type MobileAppSecurity struct {
    certPins map[string]string
}

// NewMobileAppSecurity 创建移动应用安全
func NewMobileAppSecurity() *MobileAppSecurity {
    return &MobileAppSecurity{
        certPins: make(map[string]string),
    }
}

// AddCertPin 添加证书固定
func (m *MobileAppSecurity) AddCertPin(hostname, expectedHash string) {
    m.certPins[hostname] = expectedHash
}

// VerifyCertPin 验证证书固定
func (m *MobileAppSecurity) VerifyCertPin(hostname, actualHash string) bool {
    expectedHash, exists := m.certPins[hostname]
    if !exists {
        return false
    }
    return expectedHash == actualHash
}

// SecureHTTPClient 安全HTTP客户端
func (m *MobileAppSecurity) SecureHTTPClient() *http.Client {
    // 创建TLS配置
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }

    // 创建HTTP客户端
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
    }

    return client
}

// SecureAPIRequest 安全API请求
func (m *MobileAppSecurity) SecureAPIRequest(url string) (*http.Response, error) {
    client := m.SecureHTTPClient()

    // 创建请求
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // 设置安全头
    req.Header.Set("User-Agent", "SecureMobileApp/1.0")
    req.Header.Set("Accept", "application/json")

    // 发送请求
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }

    return resp, nil
}
```

## 3.4 API安全

### 3.4.1 API安全威胁

**主要威胁**：
1. **未授权访问**：API密钥泄露
2. **数据泄露**：敏感信息暴露
3. **速率限制绕过**：API滥用
4. **注入攻击**：参数注入

### 3.4.2 API安全防护

**防护措施**：
1. **身份验证**：JWT、OAuth2
2. **授权控制**：RBAC、ABAC
3. **速率限制**：防止API滥用
4. **输入验证**：参数验证

**Go语言API安全**：

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/time/rate"
)

// APISecurity API安全
type APISecurity struct {
    jwtSecret []byte
    rateLimit *rate.Limiter
    mu        sync.RWMutex
}

// NewAPISecurity 创建API安全
func NewAPISecurity(jwtSecret string) *APISecurity {
    return &APISecurity{
        jwtSecret: []byte(jwtSecret),
        rateLimit: rate.NewLimiter(rate.Every(time.Second), 10), // 每秒10个请求
    }
}

// GenerateJWT 生成JWT令牌
func (a *APISecurity) GenerateJWT(userID string, permissions []string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    userID,
        "permissions": permissions,
        "exp":        time.Now().Add(time.Hour * 24).Unix(),
        "iat":        time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.jwtSecret)
}

// ValidateJWT 验证JWT令牌
func (a *APISecurity) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return a.jwtSecret, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %v", err)
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

// RateLimitMiddleware 速率限制中间件
func (a *APISecurity) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !a.rateLimit.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    }
}

// AuthMiddleware 认证中间件
func (a *APISecurity) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 获取Authorization头
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        // 验证JWT令牌
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := a.ValidateJWT(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // 将用户信息添加到请求上下文
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

// SecureAPIHandler 安全API处理器
func (a *APISecurity) SecureAPIHandler(w http.ResponseWriter, r *http.Request) {
    // 获取用户信息
    user := r.Context().Value("user").(jwt.MapClaims)
    userID := user["user_id"].(string)

    // 检查权限
    permissions := user["permissions"].([]interface{})
    hasPermission := false
    for _, p := range permissions {
        if p.(string) == "read" {
            hasPermission = true
            break
        }
    }

    if !hasPermission {
        http.Error(w, "Insufficient permissions", http.StatusForbidden)
        return
    }

    // 返回安全数据
    fmt.Fprintf(w, `{"message": "Secure API response", "user_id": "%s"}`, userID)
}
```

## 3.5 应用安全测试

### 3.5.1 静态应用安全测试 (SAST)

**SAST工具**：
- **Go语言**：gosec, staticcheck
- **通用工具**：SonarQube, CodeQL

### 3.5.2 动态应用安全测试 (DAST)

**DAST工具**：
- **Web应用**：OWASP ZAP, Burp Suite
- **API测试**：Postman, Insomnia

### 3.5.3 交互式应用安全测试 (IAST)

**IAST特点**：
- 运行时检测
- 低误报率
- 实时反馈

## 3.6 应用安全最佳实践

### 3.6.1 安全开发生命周期 (SDLC)

**SDLC阶段**：
1. **需求分析**：安全需求识别
2. **设计**：安全架构设计
3. **开发**：安全编码实践
4. **测试**：安全测试
5. **部署**：安全部署
6. **维护**：安全维护

### 3.6.2 安全编码实践

**编码原则**：
1. **输入验证**：所有输入都要验证
2. **输出编码**：所有输出都要编码
3. **最小权限**：使用最小权限原则
4. **深度防御**：多层安全防护

### 3.6.3 安全配置管理

**配置原则**：
1. **默认安全**：默认配置应该是安全的
2. **最小配置**：只启用必要的功能
3. **定期更新**：定期更新安全配置
4. **配置验证**：验证配置的正确性

---

**相关链接**：
- [1.1 安全基础理论](../01-Security-Foundations.md#1.1-安全基础理论)
- [2.1 网络安全基础理论](../02-Network-Security.md#2.1-网络安全基础理论)
- [4.1 数据安全基础](../04-Data-Security.md#4.1-数据安全基础)

**下一步**：
- [4. 数据安全](../04-Data-Security.md)
- [5. 密码学基础](../05-Cryptography.md)
- [6. 安全运营](../06-Security-Operations.md)
