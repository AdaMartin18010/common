# 01-安全基础 (Security Foundations)

## 概述

网络安全基础是构建安全系统的核心，包括密码学、认证授权、安全协议等技术。本章将介绍网络安全的理论基础、算法实现和Go语言应用。

## 目录

1. [理论基础](#1-理论基础)
2. [密码学基础](#2-密码学基础)
3. [认证与授权](#3-认证与授权)
4. [安全协议](#4-安全协议)
5. [密钥管理](#5-密钥管理)
6. [安全通信](#6-安全通信)
7. [安全审计](#7-安全审计)
8. [威胁防护](#8-威胁防护)

## 1. 理论基础

### 1.1 安全模型

**定义 1.1** (安全函数)
安全函数 ```latex
S: M \times K \rightarrow C
``` 将消息 ```latex
M
``` 使用密钥 ```latex
K
``` 加密为密文 ```latex
C
```：

```latex
$S(m, k) = c \text{ where } c = \text{encrypt}(m, k)
```$

**定理 1.1** (安全可逆性)
对于安全函数 ```latex
S
```，存在解密函数 ```latex
S^{-1}
``` 满足：

```latex
$S^{-1}(S(m, k), k) = m
```$

### 1.2 安全原则

```go
// 安全接口
type SecurityProvider interface {
    Encrypt(data []byte, key []byte) ([]byte, error)
    Decrypt(data []byte, key []byte) ([]byte, error)
    Hash(data []byte) ([]byte, error)
    Sign(data []byte, privateKey []byte) ([]byte, error)
    Verify(data []byte, signature []byte, publicKey []byte) (bool, error)
}

// 安全配置
type SecurityConfig struct {
    Algorithm     string            `json:"algorithm"`
    KeySize       int               `json:"key_size"`
    SaltLength    int               `json:"salt_length"`
    Iterations    int               `json:"iterations"`
    Options       map[string]string `json:"options"`
}

// 安全管理器
type SecurityManager struct {
    providers map[string]SecurityProvider
    config    SecurityConfig
    keyStore  KeyStore
}

func (sm *SecurityManager) GetProvider(name string) (SecurityProvider, error) {
    if provider, exists := sm.providers[name]; exists {
        return provider, nil
    }
    return nil, errors.New("security provider not found")
}

func (sm *SecurityManager) RegisterProvider(name string, provider SecurityProvider) {
    sm.providers[name] = provider
}
```

## 2. 密码学基础

### 2.1 对称加密

**定义 2.1** (对称加密)
对称加密使用相同密钥进行加密和解密：

```latex
$\text{Encrypt}(m, k) = c
```$
$```latex
\text{Decrypt}(c, k) = m
```$

```go
// AES加密器
type AESEncryptor struct {
    keySize int
    mode    cipher.BlockMode
}

func (ae *AESEncryptor) Encrypt(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    // 生成随机IV
    iv := make([]byte, aes.BlockSize)
    if _, err := rand.Read(iv); err != nil {
        return nil, err
    }
    
    // 填充数据
    paddedData := ae.pad(data)
    
    // 加密
    ciphertext := make([]byte, len(paddedData))
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext, paddedData)
    
    // 组合IV和密文
    result := make([]byte, 0)
    result = append(result, iv...)
    result = append(result, ciphertext...)
    
    return result, nil
}

func (ae *AESEncryptor) Decrypt(data []byte, key []byte) ([]byte, error) {
    if len(data) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    // 分离IV和密文
    iv := data[:aes.BlockSize]
    ciphertext := data[aes.BlockSize:]
    
    // 解密
    plaintext := make([]byte, len(ciphertext))
    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(plaintext, ciphertext)
    
    // 去除填充
    return ae.unpad(plaintext)
}

func (ae *AESEncryptor) pad(data []byte) []byte {
    padding := aes.BlockSize - len(data)%aes.BlockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

func (ae *AESEncryptor) unpad(data []byte) ([]byte, error) {
    length := len(data)
    if length == 0 {
        return nil, errors.New("invalid padding")
    }
    
    padding := int(data[length-1])
    if padding > length {
        return nil, errors.New("invalid padding")
    }
    
    return data[:length-padding], nil
}

func (ae *AESEncryptor) Hash(data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    return hash[:], nil
}
```

### 2.2 非对称加密

```go
// RSA加密器
type RSAEncryptor struct {
    keySize int
}

func (re *RSAEncryptor) GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
    return rsa.GenerateKey(rand.Reader, re.keySize)
}

func (re *RSAEncryptor) Encrypt(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
    return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

func (re *RSAEncryptor) Decrypt(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    return rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}

func (re *RSAEncryptor) Sign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    hash := sha256.Sum256(data)
    return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

func (re *RSAEncryptor) Verify(data []byte, signature []byte, publicKey *rsa.PublicKey) (bool, error) {
    hash := sha256.Sum256(data)
    err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
    return err == nil, err
}
```

### 2.3 哈希函数

```go
// 哈希管理器
type HashManager struct {
    algorithms map[string]hash.Hash
}

func (hm *HashManager) Hash(data []byte, algorithm string) ([]byte, error) {
    var h hash.Hash
    
    switch algorithm {
    case "sha256":
        h = sha256.New()
    case "sha512":
        h = sha512.New()
    case "md5":
        h = md5.New()
    default:
        return nil, errors.New("unsupported hash algorithm")
    }
    
    h.Write(data)
    return h.Sum(nil), nil
}

// 密码哈希
type PasswordHasher struct {
    saltLength int
    iterations int
}

func (ph *PasswordHasher) HashPassword(password []byte) ([]byte, error) {
    // 生成随机盐
    salt := make([]byte, ph.saltLength)
    if _, err := rand.Read(salt); err != nil {
        return nil, err
    }
    
    // 使用PBKDF2
    hash := pbkdf2.Key(password, salt, ph.iterations, 32, sha256.New)
    
    // 组合盐和哈希
    result := make([]byte, 0)
    result = append(result, salt...)
    result = append(result, hash...)
    
    return result, nil
}

func (ph *PasswordHasher) VerifyPassword(password, hashedPassword []byte) (bool, error) {
    if len(hashedPassword) < ph.saltLength {
        return false, errors.New("invalid hash format")
    }
    
    // 提取盐
    salt := hashedPassword[:ph.saltLength]
    storedHash := hashedPassword[ph.saltLength:]
    
    // 计算哈希
    hash := pbkdf2.Key(password, salt, ph.iterations, 32, sha256.New)
    
    return bytes.Equal(hash, storedHash), nil
}
```

## 3. 认证与授权

### 3.1 身份认证

```go
// 认证器接口
type Authenticator interface {
    Authenticate(credentials Credentials) (bool, error)
    GenerateToken(user User) (string, error)
    ValidateToken(token string) (bool, error)
}

// 凭证
type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Token    string `json:"token,omitempty"`
}

// 用户
type User struct {
    ID       string            `json:"id"`
    Username string            `json:"username"`
    Email    string            `json:"email"`
    Roles    []string          `json:"roles"`
    Metadata map[string]string `json:"metadata"`
}

// 密码认证器
type PasswordAuthenticator struct {
    userStore UserStore
    hasher    *PasswordHasher
}

func (pa *PasswordAuthenticator) Authenticate(credentials Credentials) (bool, error) {
    user, err := pa.userStore.GetByUsername(credentials.Username)
    if err != nil {
        return false, err
    }
    
    return pa.hasher.VerifyPassword([]byte(credentials.Password), user.PasswordHash)
}

func (pa *PasswordAuthenticator) GenerateToken(user User) (string, error) {
    // 创建JWT令牌
    claims := jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "roles":    user.Roles,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
        "iat":      time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("secret_key"))
}

func (pa *PasswordAuthenticator) ValidateToken(tokenString string) (bool, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret_key"), nil
    })
    
    if err != nil {
        return false, err
    }
    
    return token.Valid, nil
}

// 用户存储接口
type UserStore interface {
    GetByUsername(username string) (*User, error)
    GetByID(id string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id string) error
}

// 内存用户存储
type MemoryUserStore struct {
    users map[string]*User
    mu    sync.RWMutex
}

func (mus *MemoryUserStore) GetByUsername(username string) (*User, error) {
    mus.mu.RLock()
    defer mus.mu.RUnlock()
    
    for _, user := range mus.users {
        if user.Username == username {
            return user, nil
        }
    }
    
    return nil, errors.New("user not found")
}

func (mus *MemoryUserStore) GetByID(id string) (*User, error) {
    mus.mu.RLock()
    defer mus.mu.RUnlock()
    
    if user, exists := mus.users[id]; exists {
        return user, nil
    }
    
    return nil, errors.New("user not found")
}

func (mus *MemoryUserStore) Create(user *User) error {
    mus.mu.Lock()
    defer mus.mu.Unlock()
    
    if _, exists := mus.users[user.ID]; exists {
        return errors.New("user already exists")
    }
    
    mus.users[user.ID] = user
    return nil
}
```

### 3.2 访问控制

```go
// 授权器接口
type Authorizer interface {
    Authorize(user User, resource string, action string) (bool, error)
    GrantPermission(userID string, resource string, action string) error
    RevokePermission(userID string, resource string, action string) error
}

// RBAC授权器
type RBACAuthorizer struct {
    permissions map[string]map[string][]string // user -> resource -> actions
    roles       map[string][]string            // role -> permissions
    userRoles   map[string][]string            // user -> roles
    mu          sync.RWMutex
}

func (ra *RBACAuthorizer) Authorize(user User, resource string, action string) (bool, error) {
    ra.mu.RLock()
    defer ra.mu.RUnlock()
    
    // 检查直接权限
    if actions, exists := ra.permissions[user.ID]; exists {
        if resourceActions, exists := actions[resource]; exists {
            for _, a := range resourceActions {
                if a == action {
                    return true, nil
                }
            }
        }
    }
    
    // 检查角色权限
    if roles, exists := ra.userRoles[user.ID]; exists {
        for _, role := range roles {
            if rolePermissions, exists := ra.roles[role]; exists {
                for _, permission := range rolePermissions {
                    if permission == fmt.Sprintf("%s:%s", resource, action) {
                        return true, nil
                    }
                }
            }
        }
    }
    
    return false, nil
}

func (ra *RBACAuthorizer) GrantPermission(userID string, resource string, action string) error {
    ra.mu.Lock()
    defer ra.mu.Unlock()
    
    if ra.permissions[userID] == nil {
        ra.permissions[userID] = make(map[string][]string)
    }
    
    if ra.permissions[userID][resource] == nil {
        ra.permissions[userID][resource] = make([]string, 0)
    }
    
    // 检查是否已存在
    for _, a := range ra.permissions[userID][resource] {
        if a == action {
            return nil // 权限已存在
        }
    }
    
    ra.permissions[userID][resource] = append(ra.permissions[userID][resource], action)
    return nil
}

func (ra *RBACAuthorizer) RevokePermission(userID string, resource string, action string) error {
    ra.mu.Lock()
    defer ra.mu.Unlock()
    
    if actions, exists := ra.permissions[userID]; exists {
        if resourceActions, exists := actions[resource]; exists {
            for i, a := range resourceActions {
                if a == action {
                    ra.permissions[userID][resource] = append(resourceActions[:i], resourceActions[i+1:]...)
                    return nil
                }
            }
        }
    }
    
    return errors.New("permission not found")
}
```

## 4. 安全协议

### 4.1 TLS协议

```go
// TLS配置
type TLSConfig struct {
    CertFile    string
    KeyFile     string
    CAFile      string
    MinVersion  uint16
    MaxVersion  uint16
    CipherSuites []uint16
}

// TLS服务器
type TLSServer struct {
    config *TLSConfig
    server *http.Server
}

func (ts *TLSServer) Start(addr string) error {
    // 加载证书
    cert, err := tls.LoadX509KeyPair(ts.config.CertFile, ts.config.KeyFile)
    if err != nil {
        return err
    }
    
    // 创建TLS配置
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   ts.config.MinVersion,
        MaxVersion:   ts.config.MaxVersion,
        CipherSuites: ts.config.CipherSuites,
    }
    
    // 创建服务器
    ts.server = &http.Server{
        Addr:      addr,
        TLSConfig: tlsConfig,
    }
    
    return ts.server.ListenAndServeTLS("", "")
}

// TLS客户端
type TLSClient struct {
    config *TLSConfig
    client *http.Client
}

func (tc *TLSClient) CreateClient() error {
    // 加载CA证书
    caCert, err := ioutil.ReadFile(tc.config.CAFile)
    if err != nil {
        return err
    }
    
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    // 创建TLS配置
    tlsConfig := &tls.Config{
        RootCAs:    caCertPool,
        MinVersion: tc.config.MinVersion,
        MaxVersion: tc.config.MaxVersion,
    }
    
    // 创建客户端
    tc.client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
    }
    
    return nil
}

func (tc *TLSClient) Get(url string) (*http.Response, error) {
    return tc.client.Get(url)
}
```

### 4.2 OAuth2协议

```go
// OAuth2服务器
type OAuth2Server struct {
    clients    map[string]*OAuth2Client
    tokens     map[string]*OAuth2Token
    authCodes  map[string]*AuthCode
    mu         sync.RWMutex
}

type OAuth2Client struct {
    ID           string   `json:"id"`
    Secret       string   `json:"secret"`
    RedirectURIs []string `json:"redirect_uris"`
    Scopes       []string `json:"scopes"`
}

type OAuth2Token struct {
    AccessToken  string    `json:"access_token"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int       `json:"expires_in"`
    RefreshToken string    `json:"refresh_token"`
    Scope        string    `json:"scope"`
    UserID       string    `json:"user_id"`
    ClientID     string    `json:"client_id"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type AuthCode struct {
    Code        string    `json:"code"`
    ClientID    string    `json:"client_id"`
    UserID      string    `json:"user_id"`
    RedirectURI string    `json:"redirect_uri"`
    Scope       string    `json:"scope"`
    ExpiresAt   time.Time `json:"expires_at"`
}

func (oas *OAuth2Server) Authorize(clientID, redirectURI, scope, state string) (string, error) {
    oas.mu.Lock()
    defer oas.mu.Unlock()
    
    // 验证客户端
    client, exists := oas.clients[clientID]
    if !exists {
        return "", errors.New("invalid client")
    }
    
    // 验证重定向URI
    validRedirect := false
    for _, uri := range client.RedirectURIs {
        if uri == redirectURI {
            validRedirect = true
            break
        }
    }
    if !validRedirect {
        return "", errors.New("invalid redirect URI")
    }
    
    // 生成授权码
    code := generateRandomString(32)
    oas.authCodes[code] = &AuthCode{
        Code:        code,
        ClientID:    clientID,
        RedirectURI: redirectURI,
        Scope:       scope,
        ExpiresAt:   time.Now().Add(time.Minute * 10),
    }
    
    return code, nil
}

func (oas *OAuth2Server) ExchangeToken(clientID, clientSecret, code, redirectURI string) (*OAuth2Token, error) {
    oas.mu.Lock()
    defer oas.mu.Unlock()
    
    // 验证客户端
    client, exists := oas.clients[clientID]
    if !exists || client.Secret != clientSecret {
        return nil, errors.New("invalid client")
    }
    
    // 验证授权码
    authCode, exists := oas.authCodes[code]
    if !exists {
        return nil, errors.New("invalid authorization code")
    }
    
    if authCode.ClientID != clientID || authCode.RedirectURI != redirectURI {
        return nil, errors.New("invalid authorization code")
    }
    
    if time.Now().After(authCode.ExpiresAt) {
        delete(oas.authCodes, code)
        return nil, errors.New("authorization code expired")
    }
    
    // 生成访问令牌
    accessToken := generateRandomString(32)
    refreshToken := generateRandomString(32)
    
    token := &OAuth2Token{
        AccessToken:  accessToken,
        TokenType:    "Bearer",
        ExpiresIn:    3600,
        RefreshToken: refreshToken,
        Scope:        authCode.Scope,
        UserID:       authCode.UserID,
        ClientID:     clientID,
        ExpiresAt:    time.Now().Add(time.Hour),
    }
    
    oas.tokens[accessToken] = token
    delete(oas.authCodes, code)
    
    return token, nil
}

func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
```

## 5. 密钥管理

### 5.1 密钥存储

```go
// 密钥存储接口
type KeyStore interface {
    Store(keyID string, key []byte) error
    Retrieve(keyID string) ([]byte, error)
    Delete(keyID string) error
    List() ([]string, error)
}

// 加密密钥存储
type EncryptedKeyStore struct {
    storage    KeyStore
    masterKey  []byte
    encryptor  SecurityProvider
}

func (eks *EncryptedKeyStore) Store(keyID string, key []byte) error {
    // 加密密钥
    encryptedKey, err := eks.encryptor.Encrypt(key, eks.masterKey)
    if err != nil {
        return err
    }
    
    return eks.storage.Store(keyID, encryptedKey)
}

func (eks *EncryptedKeyStore) Retrieve(keyID string) ([]byte, error) {
    // 获取加密的密钥
    encryptedKey, err := eks.storage.Retrieve(keyID)
    if err != nil {
        return nil, err
    }
    
    // 解密密钥
    return eks.encryptor.Decrypt(encryptedKey, eks.masterKey)
}

// 密钥轮换
type KeyRotation struct {
    keyStore KeyStore
    current  string
    previous string
    next     string
}

func (kr *KeyRotation) RotateKeys() error {
    // 生成新密钥
    newKey := make([]byte, 32)
    if _, err := rand.Read(newKey); err != nil {
        return err
    }
    
    // 存储新密钥
    newKeyID := fmt.Sprintf("key_%d", time.Now().Unix())
    if err := kr.keyStore.Store(newKeyID, newKey); err != nil {
        return err
    }
    
    // 更新密钥引用
    kr.previous = kr.current
    kr.current = newKeyID
    
    return nil
}

func (kr *KeyRotation) GetCurrentKey() ([]byte, error) {
    return kr.keyStore.Retrieve(kr.current)
}

func (kr *KeyRotation) GetPreviousKey() ([]byte, error) {
    if kr.previous == "" {
        return nil, errors.New("no previous key")
    }
    return kr.keyStore.Retrieve(kr.previous)
}
```

## 6. 安全通信

### 6.1 安全通道

```go
// 安全通道
type SecureChannel struct {
    conn    net.Conn
    encrypt SecurityProvider
    decrypt SecurityProvider
    key     []byte
}

func (sc *SecureChannel) Send(data []byte) error {
    // 加密数据
    encrypted, err := sc.encrypt.Encrypt(data, sc.key)
    if err != nil {
        return err
    }
    
    // 发送长度
    length := uint32(len(encrypted))
    if err := binary.Write(sc.conn, binary.BigEndian, length); err != nil {
        return err
    }
    
    // 发送数据
    _, err = sc.conn.Write(encrypted)
    return err
}

func (sc *SecureChannel) Receive() ([]byte, error) {
    // 接收长度
    var length uint32
    if err := binary.Read(sc.conn, binary.BigEndian, &length); err != nil {
        return nil, err
    }
    
    // 接收数据
    encrypted := make([]byte, length)
    _, err := io.ReadFull(sc.conn, encrypted)
    if err != nil {
        return nil, err
    }
    
    // 解密数据
    return sc.decrypt.Decrypt(encrypted, sc.key)
}

// 安全RPC
type SecureRPC struct {
    channel *SecureChannel
}

func (sr *SecureRPC) Call(method string, args interface{}, reply interface{}) error {
    // 序列化参数
    data, err := json.Marshal(args)
    if err != nil {
        return err
    }
    
    // 创建请求
    request := RPCRequest{
        Method: method,
        Args:   data,
    }
    
    requestData, err := json.Marshal(request)
    if err != nil {
        return err
    }
    
    // 发送请求
    if err := sr.channel.Send(requestData); err != nil {
        return err
    }
    
    // 接收响应
    responseData, err := sr.channel.Receive()
    if err != nil {
        return err
    }
    
    // 解析响应
    var response RPCResponse
    if err := json.Unmarshal(responseData, &response); err != nil {
        return err
    }
    
    if response.Error != "" {
        return errors.New(response.Error)
    }
    
    // 解析结果
    return json.Unmarshal(response.Result, reply)
}

type RPCRequest struct {
    Method string          `json:"method"`
    Args   json.RawMessage `json:"args"`
}

type RPCResponse struct {
    Result json.RawMessage `json:"result"`
    Error  string          `json:"error"`
}
```

## 7. 安全审计

### 7.1 审计日志

```go
// 审计日志
type AuditLog struct {
    ID        string                 `json:"id"`
    Timestamp time.Time              `json:"timestamp"`
    UserID    string                 `json:"user_id"`
    Action    string                 `json:"action"`
    Resource  string                 `json:"resource"`
    Result    string                 `json:"result"`
    Details   map[string]interface{} `json:"details"`
    IP        string                 `json:"ip"`
    UserAgent string                 `json:"user_agent"`
}

// 审计器
type Auditor struct {
    logger AuditLogger
    config AuditConfig
}

type AuditLogger interface {
    Log(entry *AuditLog) error
    Query(filter AuditFilter) ([]*AuditLog, error)
}

type AuditConfig struct {
    Enabled     bool     `json:"enabled"`
    LogLevel    string   `json:"log_level"`
    SensitiveFields []string `json:"sensitive_fields"`
    Retention   time.Duration `json:"retention"`
}

type AuditFilter struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Result    string    `json:"result"`
}

func (a *Auditor) LogEvent(userID, action, resource, result string, details map[string]interface{}) error {
    if !a.config.Enabled {
        return nil
    }
    
    // 过滤敏感信息
    filteredDetails := a.filterSensitiveData(details)
    
    entry := &AuditLog{
        ID:        generateUUID(),
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Resource:  resource,
        Result:    result,
        Details:   filteredDetails,
    }
    
    return a.logger.Log(entry)
}

func (a *Auditor) filterSensitiveData(details map[string]interface{}) map[string]interface{} {
    filtered := make(map[string]interface{})
    
    for key, value := range details {
        sensitive := false
        for _, field := range a.config.SensitiveFields {
            if key == field {
                sensitive = true
                break
            }
        }
        
        if sensitive {
            filtered[key] = "[REDACTED]"
        } else {
            filtered[key] = value
        }
    }
    
    return filtered
}

func generateUUID() string {
    return fmt.Sprintf("%x-%x-%x-%x-%x", 
        rand.Uint32(), rand.Uint16(), rand.Uint16(), rand.Uint16(), rand.Uint64())
}
```

## 8. 威胁防护

### 8.1 入侵检测

```go
// 入侵检测系统
type IntrusionDetectionSystem struct {
    rules    []DetectionRule
    alerts   chan *SecurityAlert
    analyzer *TrafficAnalyzer
}

type DetectionRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Pattern     string                 `json:"pattern"`
    Threshold   int                    `json:"threshold"`
    TimeWindow  time.Duration          `json:"time_window"`
    Actions     []string               `json:"actions"`
    Enabled     bool                   `json:"enabled"`
}

type SecurityAlert struct {
    ID        string                 `json:"id"`
    RuleID    string                 `json:"rule_id"`
    Severity  string                 `json:"severity"`
    Message   string                 `json:"message"`
    Timestamp time.Time              `json:"timestamp"`
    Source    string                 `json:"source"`
    Details   map[string]interface{} `json:"details"`
}

type TrafficAnalyzer struct {
    patterns map[string]*Pattern
    counters map[string]*Counter
    mu       sync.RWMutex
}

type Pattern struct {
    Regex   *regexp.Regexp
    Matches int
    LastSeen time.Time
}

type Counter struct {
    Count     int
    FirstSeen time.Time
    LastSeen  time.Time
}

func (ids *IntrusionDetectionSystem) AnalyzeTraffic(data []byte, source string) error {
    ids.analyzer.mu.Lock()
    defer ids.analyzer.mu.Unlock()
    
    for _, rule := range ids.rules {
        if !rule.Enabled {
            continue
        }
        
        if ids.matchesRule(data, rule) {
            if ids.shouldAlert(rule, source) {
                alert := &SecurityAlert{
                    ID:        generateUUID(),
                    RuleID:    rule.ID,
                    Severity:  "high",
                    Message:   fmt.Sprintf("Rule '%s' triggered", rule.Name),
                    Timestamp: time.Now(),
                    Source:    source,
                    Details:   map[string]interface{}{"data": string(data)},
                }
                
                select {
                case ids.alerts <- alert:
                default:
                    // 通道满，丢弃告警
                }
            }
        }
    }
    
    return nil
}

func (ids *IntrusionDetectionSystem) matchesRule(data []byte, rule DetectionRule) bool {
    pattern, exists := ids.analyzer.patterns[rule.Pattern]
    if !exists {
        regex, err := regexp.Compile(rule.Pattern)
        if err != nil {
            return false
        }
        pattern = &Pattern{Regex: regex}
        ids.analyzer.patterns[rule.Pattern] = pattern
    }
    
    if pattern.Regex.Match(data) {
        pattern.Matches++
        pattern.LastSeen = time.Now()
        return true
    }
    
    return false
}

func (ids *IntrusionDetectionSystem) shouldAlert(rule DetectionRule, source string) bool {
    counterKey := fmt.Sprintf("%s:%s", rule.ID, source)
    counter, exists := ids.analyzer.counters[counterKey]
    
    if !exists {
        counter = &Counter{FirstSeen: time.Now()}
        ids.analyzer.counters[counterKey] = counter
    }
    
    counter.Count++
    counter.LastSeen = time.Now()
    
    // 检查时间窗口
    if time.Since(counter.FirstSeen) > rule.TimeWindow {
        counter.Count = 1
        counter.FirstSeen = time.Now()
    }
    
    return counter.Count >= rule.Threshold
}
```

## 总结

本章详细介绍了网络安全基础的核心技术，包括：

1. **理论基础**：安全模型、安全原则
2. **密码学基础**：对称加密、非对称加密、哈希函数
3. **认证与授权**：身份认证、访问控制
4. **安全协议**：TLS协议、OAuth2协议
5. **密钥管理**：密钥存储、密钥轮换
6. **安全通信**：安全通道、安全RPC
7. **安全审计**：审计日志、审计器
8. **威胁防护**：入侵检测、安全告警

这些技术为构建安全、可靠的网络系统提供了完整的理论基础和实现方案。

---

**相关链接**：

- [02-网络安全协议](../02-Security-Protocols.md)
- [03-安全架构](../03-Security-Architecture.md)
- [04-安全测试](../04-Security-Testing.md)
- [05-安全监控](../05-Security-Monitoring.md)
- [06-安全响应](../06-Security-Response.md)
- [07-安全合规](../07-Security-Compliance.md)
- [08-安全工具](../08-Security-Tools.md)
