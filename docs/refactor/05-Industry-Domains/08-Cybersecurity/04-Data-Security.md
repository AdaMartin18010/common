# 4. 数据安全 (Data Security)

## 4.1 数据安全基础

### 4.1.1 数据安全定义与分类

数据安全是保护数据免受未授权访问、使用、披露、中断、修改或破坏的过程。

**形式化定义**：
设 $D$ 为数据集，$U$ 为用户集合，$A$ 为操作集合。
数据安全函数 $S: D \times U \times A \rightarrow \{0,1\}$ 定义为：
$$S(d,u,a) = \begin{cases}
1 & \text{if user } u \text{ is authorized to perform action } a \text{ on data } d \\
0 & \text{otherwise}
\end{cases}$$

### 4.1.2 数据分类

**数据敏感度分类**：
1. **公开数据**：可以公开访问的数据
2. **内部数据**：仅限内部人员访问的数据
3. **机密数据**：需要特殊授权访问的数据
4. **绝密数据**：最高级别保护的数据

**数据生命周期**：
```
创建 → 存储 → 使用 → 传输 → 归档 → 销毁
```

## 4.2 数据保护技术

### 4.2.1 数据加密

**加密类型**：
1. **对称加密**：使用相同密钥加密和解密
2. **非对称加密**：使用公钥加密，私钥解密
3. **同态加密**：在加密数据上进行计算

**Go语言数据加密实现**：

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "io"
)

// DataEncryption 数据加密
type DataEncryption struct {
    aesKey []byte
    rsaKey *rsa.PrivateKey
}

// NewDataEncryption 创建数据加密
func NewDataEncryption() (*DataEncryption, error) {
    // 生成AES密钥
    aesKey := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
        return nil, fmt.Errorf("failed to generate AES key: %v", err)
    }

    // 生成RSA密钥对
    rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, fmt.Errorf("failed to generate RSA key: %v", err)
    }

    return &DataEncryption{
        aesKey: aesKey,
        rsaKey: rsaKey,
    }, nil
}

// EncryptAES AES加密
func (d *DataEncryption) EncryptAES(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(d.aesKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 生成随机IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %v", err)
    }

    // 创建GCM模式
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    // 加密
    ciphertext := gcm.Seal(iv, iv, plaintext, nil)
    return ciphertext, nil
}

// DecryptAES AES解密
func (d *DataEncryption) DecryptAES(ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(d.aesKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    // 创建GCM模式
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    // 提取IV
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // 解密
    plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %v", err)
    }

    return plaintext, nil
}

// EncryptRSA RSA加密
func (d *DataEncryption) EncryptRSA(plaintext []byte) ([]byte, error) {
    // 使用OAEP填充
    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &d.rsaKey.PublicKey, plaintext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt: %v", err)
    }
    return ciphertext, nil
}

// DecryptRSA RSA解密
func (d *DataEncryption) DecryptRSA(ciphertext []byte) ([]byte, error) {
    // 使用OAEP填充
    plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, d.rsaKey, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %v", err)
    }
    return plaintext, nil
}

// HybridEncrypt 混合加密
func (d *DataEncryption) HybridEncrypt(plaintext []byte) ([]byte, error) {
    // 使用AES加密数据
    aesCiphertext, err := d.EncryptAES(plaintext)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt with AES: %v", err)
    }

    // 使用RSA加密AES密钥
    rsaCiphertext, err := d.EncryptRSA(d.aesKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt AES key: %v", err)
    }

    // 组合加密结果
    result := append(rsaCiphertext, aesCiphertext...)
    return result, nil
}

// HybridDecrypt 混合解密
func (d *DataEncryption) HybridDecrypt(ciphertext []byte) ([]byte, error) {
    // 分离RSA和AES密文
    rsaKeySize := 256 // RSA-2048加密后的密钥大小
    rsaCiphertext := ciphertext[:rsaKeySize]
    aesCiphertext := ciphertext[rsaKeySize:]

    // 使用RSA解密AES密钥
    aesKey, err := d.DecryptRSA(rsaCiphertext)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt AES key: %v", err)
    }

    // 使用AES解密数据
    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    iv := aesCiphertext[:aes.BlockSize]
    aesCiphertext = aesCiphertext[aes.BlockSize:]

    plaintext, err := gcm.Open(nil, iv, aesCiphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %v", err)
    }

    return plaintext, nil
}
```

### 4.2.2 数据脱敏

**脱敏技术**：
1. **静态脱敏**：永久性数据脱敏
2. **动态脱敏**：实时数据脱敏
3. **格式保持脱敏**：保持数据格式的脱敏

**Go语言数据脱敏实现**：

```go
package main

import (
    "crypto/md5"
    "fmt"
    "regexp"
    "strings"
)

// DataMasking 数据脱敏
type DataMasking struct{}

// NewDataMasking 创建数据脱敏
func NewDataMasking() *DataMasking {
    return &DataMasking{}
}

// MaskEmail 邮箱脱敏
func (d *DataMasking) MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }

    username := parts[0]
    domain := parts[1]

    if len(username) <= 2 {
        return email
    }

    maskedUsername := username[:2] + "***"
    return maskedUsername + "@" + domain
}

// MaskPhone 手机号脱敏
func (d *DataMasking) MaskPhone(phone string) string {
    if len(phone) < 7 {
        return phone
    }

    return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskIDCard 身份证脱敏
func (d *DataMasking) MaskIDCard(idCard string) string {
    if len(idCard) < 10 {
        return idCard
    }

    return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// MaskCreditCard 信用卡脱敏
func (d *DataMasking) MaskCreditCard(cardNumber string) string {
    if len(cardNumber) < 8 {
        return cardNumber
    }

    return cardNumber[:4] + " **** **** " + cardNumber[len(cardNumber)-4:]
}

// HashSensitiveData 敏感数据哈希
func (d *DataMasking) HashSensitiveData(data string) string {
    hash := md5.Sum([]byte(data))
    return fmt.Sprintf("%x", hash)
}

// AnonymizeData 数据匿名化
func (d *DataMasking) AnonymizeData(data string, dataType string) string {
    switch dataType {
    case "email":
        return d.MaskEmail(data)
    case "phone":
        return d.MaskPhone(data)
    case "idcard":
        return d.MaskIDCard(data)
    case "creditcard":
        return d.MaskCreditCard(data)
    default:
        return d.HashSensitiveData(data)
    }
}
```

## 4.3 数据访问控制

### 4.3.1 访问控制模型

**访问控制类型**：
1. **DAC (Discretionary Access Control)**：自主访问控制
2. **MAC (Mandatory Access Control)**：强制访问控制
3. **RBAC (Role-Based Access Control)**：基于角色的访问控制
4. **ABAC (Attribute-Based Access Control)**：基于属性的访问控制

**Go语言访问控制实现**：

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Permission 权限
type Permission struct {
    Resource string
    Action   string
}

// Role 角色
type Role struct {
    ID          string
    Name        string
    Permissions []Permission
}

// User 用户
type User struct {
    ID       string
    Username string
    Roles    []string
    Level    int // 安全级别
}

// DataAccessControl 数据访问控制
type DataAccessControl struct {
    users       map[string]*User
    roles       map[string]*Role
    dataLevels  map[string]int
    mu          sync.RWMutex
}

// NewDataAccessControl 创建数据访问控制
func NewDataAccessControl() *DataAccessControl {
    return &DataAccessControl{
        users:      make(map[string]*User),
        roles:      make(map[string]*Role),
        dataLevels: make(map[string]int),
    }
}

// AddUser 添加用户
func (d *DataAccessControl) AddUser(user *User) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.users[user.ID] = user
}

// AddRole 添加角色
func (d *DataAccessControl) AddRole(role *Role) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.roles[role.ID] = role
}

// SetDataLevel 设置数据安全级别
func (d *DataAccessControl) SetDataLevel(dataID string, level int) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.dataLevels[dataID] = level
}

// CheckAccess 检查访问权限
func (d *DataAccessControl) CheckAccess(userID, dataID, action string) bool {
    d.mu.RLock()
    defer d.mu.RUnlock()

    user, exists := d.users[userID]
    if !exists {
        return false
    }

    dataLevel, exists := d.dataLevels[dataID]
    if !exists {
        return false
    }

    // 检查安全级别
    if user.Level < dataLevel {
        return false
    }

    // 检查角色权限
    for _, roleID := range user.Roles {
        role, exists := d.roles[roleID]
        if !exists {
            continue
        }

        for _, permission := range role.Permissions {
            if permission.Resource == dataID && permission.Action == action {
                return true
            }
        }
    }

    return false
}

// RBACAccessControl 基于角色的访问控制
type RBACAccessControl struct {
    dac *DataAccessControl
}

// NewRBACAccessControl 创建RBAC访问控制
func NewRBACAccessControl() *RBACAccessControl {
    return &RBACAccessControl{
        dac: NewDataAccessControl(),
    }
}

// GrantPermission 授予权限
func (r *RBACAccessControl) GrantPermission(roleID, resource, action string) {
    role, exists := r.dac.roles[roleID]
    if !exists {
        return
    }

    permission := Permission{
        Resource: resource,
        Action:   action,
    }

    role.Permissions = append(role.Permissions, permission)
}

// RevokePermission 撤销权限
func (r *RBACAccessControl) RevokePermission(roleID, resource, action string) {
    role, exists := r.dac.roles[roleID]
    if !exists {
        return
    }

    for i, permission := range role.Permissions {
        if permission.Resource == resource && permission.Action == action {
            role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
            break
        }
    }
}

// AssignRole 分配角色
func (r *RBACAccessControl) AssignRole(userID, roleID string) {
    user, exists := r.dac.users[userID]
    if !exists {
        return
    }

    user.Roles = append(user.Roles, roleID)
}
```

## 4.4 数据隐私保护

### 4.4.1 隐私保护技术

**隐私保护方法**：
1. **差分隐私**：在数据中添加噪声
2. **k-匿名**：确保每个记录至少与k-1个其他记录相同
3. **l-多样性**：确保敏感属性有足够的多样性
4. **t-接近性**：限制敏感属性的分布

### 4.4.2 数据治理

**数据治理原则**：
1. **数据质量**：确保数据的准确性和完整性
2. **数据生命周期**：管理数据的整个生命周期
3. **数据合规**：确保符合法律法规要求
4. **数据安全**：保护数据免受安全威胁

**Go语言数据治理实现**：

```go
package main

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "time"
)

// DataGovernance 数据治理
type DataGovernance struct {
    policies map[string]DataPolicy
    audits   []DataAudit
}

// DataPolicy 数据策略
type DataPolicy struct {
    ID          string
    Name        string
    Description string
    Rules       []DataRule
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// DataRule 数据规则
type DataRule struct {
    Field       string
    Validation  string
    Encryption  bool
    Masking     bool
    Retention   time.Duration
}

// DataAudit 数据审计
type DataAudit struct {
    ID        string
    UserID    string
    Action    string
    Resource  string
    Timestamp time.Time
    Hash      string
}

// NewDataGovernance 创建数据治理
func NewDataGovernance() *DataGovernance {
    return &DataGovernance{
        policies: make(map[string]DataPolicy),
        audits:   make([]DataAudit, 0),
    }
}

// AddPolicy 添加数据策略
func (d *DataGovernance) AddPolicy(policy DataPolicy) {
    policy.CreatedAt = time.Now()
    policy.UpdatedAt = time.Now()
    d.policies[policy.ID] = policy
}

// ValidateData 验证数据
func (d *DataGovernance) ValidateData(data map[string]interface{}, policyID string) error {
    policy, exists := d.policies[policyID]
    if !exists {
        return fmt.Errorf("policy not found")
    }

    for _, rule := range policy.Rules {
        value, exists := data[rule.Field]
        if !exists {
            continue
        }

        // 验证规则
        if err := d.validateRule(value, rule); err != nil {
            return fmt.Errorf("validation failed for field %s: %v", rule.Field, err)
        }
    }

    return nil
}

// validateRule 验证规则
func (d *DataGovernance) validateRule(value interface{}, rule DataRule) error {
    switch rule.Validation {
    case "email":
        return d.validateEmail(value.(string))
    case "phone":
        return d.validatePhone(value.(string))
    case "required":
        if value == nil || value == "" {
            return fmt.Errorf("field is required")
        }
    }
    return nil
}

// validateEmail 验证邮箱
func (d *DataGovernance) validateEmail(email string) error {
    emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(emailPattern, email)
    if !matched {
        return fmt.Errorf("invalid email format")
    }
    return nil
}

// validatePhone 验证手机号
func (d *DataGovernance) validatePhone(phone string) error {
    phonePattern := `^1[3-9]\d{9}$`
    matched, _ := regexp.MatchString(phonePattern, phone)
    if !matched {
        return fmt.Errorf("invalid phone format")
    }
    return nil
}

// AuditAction 审计操作
func (d *DataGovernance) AuditAction(userID, action, resource string) {
    audit := DataAudit{
        ID:        fmt.Sprintf("audit_%d", time.Now().Unix()),
        UserID:    userID,
        Action:    action,
        Resource:  resource,
        Timestamp: time.Now(),
    }

    // 计算哈希
    data, _ := json.Marshal(audit)
    hash := sha256.Sum256(data)
    audit.Hash = fmt.Sprintf("%x", hash)

    d.audits = append(d.audits, audit)
}

// GetAuditLog 获取审计日志
func (d *DataGovernance) GetAuditLog(userID string, startTime, endTime time.Time) []DataAudit {
    var result []DataAudit
    for _, audit := range d.audits {
        if audit.UserID == userID && audit.Timestamp.After(startTime) && audit.Timestamp.Before(endTime) {
            result = append(result, audit)
        }
    }
    return result
}
```

## 4.5 数据备份与恢复

### 4.5.1 备份策略

**备份类型**：
1. **完全备份**：备份所有数据
2. **增量备份**：只备份变化的数据
3. **差异备份**：备份上次完全备份后的所有变化

### 4.5.2 数据恢复

**恢复策略**：
1. **即时恢复**：快速恢复到最近状态
2. **点时间恢复**：恢复到指定时间点
3. **灾难恢复**：从灾难中恢复数据

## 4.6 数据安全最佳实践

### 4.6.1 数据分类管理

**分类原则**：
1. **按敏感度分类**：根据数据敏感程度分类
2. **按用途分类**：根据数据用途分类
3. **按法规分类**：根据法规要求分类

### 4.6.2 数据生命周期管理

**生命周期阶段**：
1. **创建**：数据创建时的安全措施
2. **存储**：数据存储时的安全措施
3. **使用**：数据使用时的安全措施
4. **传输**：数据传输时的安全措施
5. **归档**：数据归档时的安全措施
6. **销毁**：数据销毁时的安全措施

### 4.6.3 数据安全监控

**监控内容**：
1. **访问监控**：监控数据访问行为
2. **异常检测**：检测异常访问模式
3. **合规监控**：监控合规性要求

---

**相关链接**：
- [1.1 安全基础理论](../01-Security-Foundations.md#1.1-安全基础理论)
- [2.1 网络安全基础理论](../02-Network-Security.md#2.1-网络安全基础理论)
- [3.1 应用安全基础](../03-Application-Security.md#3.1-应用安全基础)

**下一步**：
- [5. 密码学基础](../05-Cryptography.md)
- [6. 安全运营](../06-Security-Operations.md)
- [7. 威胁情报](../07-Threat-Intelligence.md)
