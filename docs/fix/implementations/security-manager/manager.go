package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SecurityManager 安全管理器
type SecurityManager struct {
	authenticator    *Authenticator
	accessController *AccessController
	dataEncryptor    *DataEncryptor
	securityLogger   *SecurityLogger
	inputFilter      *InputFilter
	pathValidator    *PathValidator
	logger           *zap.Logger
	config           *SecurityConfig
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret           string        `json:"jwt_secret"`
	EncryptionKey       string        `json:"encryption_key"`
	SecurityLogPath     string        `json:"security_log_path"`
	MaxInputLength      int           `json:"max_input_length"`
	ForbiddenPatterns   []string      `json:"forbidden_patterns"`
	AllowedPaths        []string      `json:"allowed_paths"`
	SessionTimeout      time.Duration `json:"session_timeout"`
	PasswordMinLength   int           `json:"password_min_length"`
	RequireSpecialChars bool          `json:"require_special_chars"`
	EnableAudit         bool          `json:"enable_audit"`
	AuditLogPath        string        `json:"audit_log_path"`
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(config *SecurityConfig) (*SecurityManager, error) {
	authenticator := NewAuthenticator(config.JWTSecret)
	accessController := NewAccessController()

	dataEncryptor := NewDataEncryptor(config.EncryptionKey)

	securityLogger, err := NewSecurityLogger(config.SecurityLogPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create security logger: %w", err)
	}

	inputFilter, err := NewInputFilter(config.MaxInputLength, config.ForbiddenPatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to create input filter: %w", err)
	}

	pathValidator := NewPathValidator(config.AllowedPaths)

	return &SecurityManager{
		authenticator:    authenticator,
		accessController: accessController,
		dataEncryptor:    dataEncryptor,
		securityLogger:   securityLogger,
		inputFilter:      inputFilter,
		pathValidator:    pathValidator,
		logger:           zap.L().Named("security-manager"),
		config:           config,
	}, nil
}

// AuthenticateUser 用户认证
func (sm *SecurityManager) AuthenticateUser(username, password string) (*Session, error) {
	session, err := sm.authenticator.Authenticate(username, password)
	if err != nil {
		sm.securityLogger.LogAuthenticationEvent(username, "FAILED", err.Error())
		return nil, err
	}

	sm.securityLogger.LogAuthenticationEvent(username, "SUCCESS", "")
	return session, nil
}

// ValidateToken 验证令牌
func (sm *SecurityManager) ValidateToken(tokenString string) (*User, error) {
	return sm.authenticator.ValidateToken(tokenString)
}

// AuthorizeAccess 授权访问
func (sm *SecurityManager) AuthorizeAccess(session *Session, resource, action string) bool {
	user, err := sm.authenticator.ValidateToken(session.Token)
	if err != nil {
		sm.securityLogger.LogAuthorizationEvent(session.UserID, resource, action, "DENIED")
		return false
	}

	authorized := sm.accessController.CheckAccess(user, resource, action)
	if authorized {
		sm.securityLogger.LogAuthorizationEvent(user.Username, resource, action, "ALLOWED")
	} else {
		sm.securityLogger.LogAuthorizationEvent(user.Username, resource, action, "DENIED")
	}

	return authorized
}

// ValidateInput 验证输入
func (sm *SecurityManager) ValidateInput(input string) (string, error) {
	return sm.inputFilter.Filter(input)
}

// ValidatePath 验证路径
func (sm *SecurityManager) ValidatePath(path string) (string, error) {
	return sm.pathValidator.ValidatePath(path)
}

// EncryptData 加密数据
func (sm *SecurityManager) EncryptData(data []byte) ([]byte, error) {
	return sm.dataEncryptor.Encrypt(data)
}

// DecryptData 解密数据
func (sm *SecurityManager) DecryptData(data []byte) ([]byte, error) {
	return sm.dataEncryptor.Decrypt(data)
}

// LogSecurityEvent 记录安全事件
func (sm *SecurityManager) LogSecurityEvent(eventType, user, resource, action, result string) {
	sm.securityLogger.LogSecurityEvent(eventType, user, resource, action, result)
}

// Close 关闭安全管理器
func (sm *SecurityManager) Close() error {
	return sm.securityLogger.Close()
}

// Authenticator 身份认证器
type Authenticator struct {
	users     map[string]*User
	sessions  map[string]*Session
	jwtSecret []byte
	logger    *zap.Logger
	mu        sync.RWMutex
}

// User 用户
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Session 会话
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Created   time.Time `json:"created"`
}

// NewAuthenticator 创建身份认证器
func NewAuthenticator(jwtSecret string) *Authenticator {
	return &Authenticator{
		users:     make(map[string]*User),
		sessions:  make(map[string]*Session),
		jwtSecret: []byte(jwtSecret),
		logger:    zap.L().Named("authenticator"),
	}
}

// RegisterUser 注册用户
func (auth *Authenticator) RegisterUser(username, password, role string) error {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	if _, exists := auth.users[username]; exists {
		return errors.New("user already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		ID:       generateID(),
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
		Active:   true,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	auth.users[username] = user
	auth.logger.Info("user registered", zap.String("username", username))
	return nil
}

// Authenticate 认证
func (auth *Authenticator) Authenticate(username, password string) (*Session, error) {
	auth.mu.RLock()
	user, exists := auth.users[username]
	auth.mu.RUnlock()

	if !exists {
		auth.logger.Warn("authentication failed: user not found", zap.String("username", username))
		return nil, errors.New("invalid credentials")
	}

	if !user.Active {
		auth.logger.Warn("authentication failed: user inactive", zap.String("username", username))
		return nil, errors.New("user inactive")
	}

	// 验证密码
	if !auth.verifyPassword(password, user.Password) {
		auth.logger.Warn("authentication failed: invalid password", zap.String("username", username))
		return nil, errors.New("invalid credentials")
	}

	// 创建会话
	session := &Session{
		ID:        generateID(),
		UserID:    user.ID,
		Token:     auth.generateToken(user),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Created:   time.Now(),
	}

	auth.mu.Lock()
	auth.sessions[session.ID] = session
	auth.mu.Unlock()

	auth.logger.Info("user authenticated", zap.String("username", username))
	return session, nil
}

// verifyPassword 验证密码
func (auth *Authenticator) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken 生成令牌
func (auth *Authenticator) generateToken(user *User) string {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.jwtSecret)

	return tokenString
}

// ValidateToken 验证令牌
func (auth *Authenticator) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)

		auth.mu.RLock()
		user, exists := auth.users[userID]
		auth.mu.RUnlock()

		if !exists {
			return nil, errors.New("user not found")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}

// AccessController 访问控制器
type AccessController struct {
	policies map[string]*Policy
	logger   *zap.Logger
	mu       sync.RWMutex
}

// Policy 策略
type Policy struct {
	ID       string   `json:"id"`
	Resource string   `json:"resource"`
	Action   string   `json:"action"`
	Roles    []string `json:"roles"`
	Effect   string   `json:"effect"` // Allow or Deny
}

// NewAccessController 创建访问控制器
func NewAccessController() *AccessController {
	return &AccessController{
		policies: make(map[string]*Policy),
		logger:   zap.L().Named("access-controller"),
	}
}

// AddPolicy 添加策略
func (ac *AccessController) AddPolicy(policy *Policy) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.policies[policy.ID] = policy
	ac.logger.Info("policy added", zap.String("policy_id", policy.ID))
}

// RemovePolicy 移除策略
func (ac *AccessController) RemovePolicy(policyID string) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.policies[policyID]; !exists {
		return fmt.Errorf("policy %s not found", policyID)
	}

	delete(ac.policies, policyID)
	ac.logger.Info("policy removed", zap.String("policy_id", policyID))
	return nil
}

// CheckAccess 检查访问权限
func (ac *AccessController) CheckAccess(user *User, resource, action string) bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	for _, policy := range ac.policies {
		if policy.Resource == resource && policy.Action == action {
			for _, role := range policy.Roles {
				if role == user.Role {
					if policy.Effect == "Allow" {
						ac.logger.Info("access allowed",
							zap.String("user", user.Username),
							zap.String("resource", resource),
							zap.String("action", action))
						return true
					} else {
						ac.logger.Warn("access denied by policy",
							zap.String("user", user.Username),
							zap.String("resource", resource),
							zap.String("action", action))
						return false
					}
				}
			}
		}
	}

	ac.logger.Warn("access denied: no matching policy",
		zap.String("user", user.Username),
		zap.String("resource", resource),
		zap.String("action", action))
	return false
}

// DataEncryptor 数据加密器
type DataEncryptor struct {
	key    []byte
	logger *zap.Logger
}

// NewDataEncryptor 创建数据加密器
func NewDataEncryptor(key string) *DataEncryptor {
	return &DataEncryptor{
		key:    []byte(key),
		logger: zap.L().Named("data-encryptor"),
	}
}

// Encrypt 加密
func (de *DataEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(de.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	de.logger.Info("data encrypted", zap.Int("size", len(data)))

	return ciphertext, nil
}

// Decrypt 解密
func (de *DataEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(de.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	de.logger.Info("data decrypted", zap.Int("size", len(plaintext)))
	return plaintext, nil
}

// SecurityLogger 安全日志记录器
type SecurityLogger struct {
	logger *zap.Logger
	file   *os.File
	config *SecurityConfig
}

// NewSecurityLogger 创建安全日志记录器
func NewSecurityLogger(logPath string) (*SecurityLogger, error) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open security log file: %w", err)
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logPath}

	logger, err := config.Build()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &SecurityLogger{
		logger: logger.Named("security"),
		file:   file,
	}, nil
}

// LogSecurityEvent 记录安全事件
func (sl *SecurityLogger) LogSecurityEvent(eventType, user, resource, action, result string) {
	sl.logger.Info("security event",
		zap.String("event_type", eventType),
		zap.String("user", user),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.String("result", result),
		zap.Time("timestamp", time.Now()),
		zap.String("ip", sl.getClientIP()),
	)
}

// LogAuthenticationEvent 记录认证事件
func (sl *SecurityLogger) LogAuthenticationEvent(username, result, reason string) {
	sl.logger.Info("authentication event",
		zap.String("username", username),
		zap.String("result", result),
		zap.String("reason", reason),
		zap.Time("timestamp", time.Now()),
		zap.String("ip", sl.getClientIP()),
	)
}

// LogAuthorizationEvent 记录授权事件
func (sl *SecurityLogger) LogAuthorizationEvent(user, resource, action, result string) {
	sl.logger.Info("authorization event",
		zap.String("user", user),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.String("result", result),
		zap.Time("timestamp", time.Now()),
		zap.String("ip", sl.getClientIP()),
	)
}

// getClientIP 获取客户端IP
func (sl *SecurityLogger) getClientIP() string {
	// 这里应该从请求上下文中获取客户端IP
	return "unknown"
}

// Close 关闭
func (sl *SecurityLogger) Close() error {
	return sl.file.Close()
}

// InputFilter 输入过滤器
type InputFilter struct {
	maxLength int
	patterns  []*regexp.Regexp
	logger    *zap.Logger
}

// NewInputFilter 创建输入过滤器
func NewInputFilter(maxLength int, patterns []string) (*InputFilter, error) {
	compiledPatterns := make([]*regexp.Regexp, 0)
	for _, pattern := range patterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile pattern %s: %w", pattern, err)
		}
		compiledPatterns = append(compiledPatterns, compiled)
	}

	return &InputFilter{
		maxLength: maxLength,
		patterns:  compiledPatterns,
		logger:    zap.L().Named("input-filter"),
	}, nil
}

// Filter 过滤输入
func (fi *InputFilter) Filter(input string) (string, error) {
	// 检查长度
	if len(input) > fi.maxLength {
		fi.logger.Warn("input too long", zap.Int("length", len(input)), zap.Int("max", fi.maxLength))
		return "", errors.New("input too long")
	}

	// 检查模式
	for _, pattern := range fi.patterns {
		if pattern.MatchString(input) {
			fi.logger.Warn("input contains forbidden pattern",
				zap.String("input", input),
				zap.String("pattern", pattern.String()))
			return "", errors.New("input contains forbidden characters")
		}
	}

	// 转义特殊字符
	filtered := html.EscapeString(input)

	return filtered, nil
}

// ValidateEmail 验证邮箱
func (fi *InputFilter) ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

// ValidateURL 验证URL
func (fi *InputFilter) ValidateURL(url string) bool {
	_, err := http.ParseURL(url)
	return err == nil
}

// PathValidator 路径验证器
type PathValidator struct {
	allowedPaths []string
	logger       *zap.Logger
}

// NewPathValidator 创建路径验证器
func NewPathValidator(allowedPaths []string) *PathValidator {
	return &PathValidator{
		allowedPaths: allowedPaths,
		logger:       zap.L().Named("path-validator"),
	}
}

// ValidatePath 验证路径
func (pv *PathValidator) ValidatePath(path string) (string, error) {
	// 检查路径遍历攻击
	if strings.Contains(path, "..") {
		pv.logger.Warn("path traversal attempt detected", zap.String("path", path))
		return "", errors.New("invalid path")
	}

	// 检查绝对路径
	if filepath.IsAbs(path) {
		// 检查是否在允许的路径范围内
		for _, allowedPath := range pv.allowedPaths {
			if strings.HasPrefix(path, allowedPath) {
				return filepath.Clean(path), nil
			}
		}
		pv.logger.Warn("access to unauthorized path attempted", zap.String("path", path))
		return "", errors.New("access denied")
	}

	// 处理相对路径
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		pv.logger.Warn("relative path traversal attempt", zap.String("path", path))
		return "", errors.New("invalid path")
	}

	return cleanPath, nil
}

// IsPathSafe 检查路径是否安全
func (pv *PathValidator) IsPathSafe(path string) bool {
	_, err := pv.ValidatePath(path)
	return err == nil
}

// generateID 生成ID
func generateID() string {
	// 简单的ID生成，实际应用中应该使用UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
