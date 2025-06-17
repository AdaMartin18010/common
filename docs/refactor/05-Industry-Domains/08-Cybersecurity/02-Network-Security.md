# 2. 网络安全 (Network Security)

## 2.1 网络安全基础理论

### 2.1.1 网络安全定义与分类

网络安全是保护网络基础设施、网络服务和网络数据免受未授权访问、使用、披露、中断、修改或破坏的过程。

**形式化定义**：
设 $N = (V, E)$ 为网络图，其中 $V$ 为节点集合，$E$ 为边集合。
网络安全函数 $S: N \times T \rightarrow \{0,1\}$ 定义为：
$$S(N,t) = \begin{cases}
1 & \text{if } N \text{ is secure at time } t \\
0 & \text{otherwise}
\end{cases}$$

### 2.1.2 网络安全威胁模型

**威胁分类**：
1. **被动威胁**：信息窃听、流量分析
2. **主动威胁**：数据篡改、拒绝服务、重放攻击
3. **内部威胁**：恶意内部人员
4. **外部威胁**：外部攻击者

**威胁概率模型**：
$$P(T_i) = \sum_{j=1}^{n} P(T_i|V_j) \cdot P(V_j)$$
其中 $T_i$ 为威胁类型，$V_j$ 为漏洞类型。

### 2.1.3 网络安全原则

**CIA三元组**：
- **机密性 (Confidentiality)**：$C = \frac{|S|}{|T|}$
- **完整性 (Integrity)**：$I = \frac{|V|}{|D|}$
- **可用性 (Availability)**：$A = \frac{|U|}{|R|}$

其中 $S$ 为安全传输的数据，$T$ 为总数据，$V$ 为验证通过的数据，$D$ 为总数据，$U$ 为可用时间，$R$ 为总时间。

## 2.2 网络协议安全

### 2.2.1 TCP/IP协议栈安全

**协议层次安全**：
```
应用层: HTTPS, SSH, SFTP
传输层: TLS/SSL, DTLS
网络层: IPsec, VPN
链路层: WPA2, MAC过滤
物理层: 物理隔离
```

### 2.2.2 TLS/SSL协议

**TLS握手过程**：
1. **Client Hello**: 客户端发送支持的加密套件
2. **Server Hello**: 服务器选择加密套件
3. **Certificate**: 服务器发送证书
4. **Key Exchange**: 密钥交换
5. **Finished**: 握手完成

**Go语言TLS实现**：

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

// TLSConfig 配置TLS连接
type TLSConfig struct {
    CertFile   string
    KeyFile    string
    CAFile     string
    ServerName string
}

// NewTLSConfig 创建TLS配置
func NewTLSConfig(config TLSConfig) (*tls.Config, error) {
    // 加载证书
    cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
    if err != nil {
        return nil, fmt.Errorf("failed to load certificate: %v", err)
    }

    // 加载CA证书
    caCert, err := ioutil.ReadFile(config.CAFile)
    if err != nil {
        return nil, fmt.Errorf("failed to load CA certificate: %v", err)
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        ServerName:   config.ServerName,
        MinVersion:   tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }, nil
}

// TLSServer TLS服务器
type TLSServer struct {
    config *tls.Config
    server *http.Server
}

// NewTLSServer 创建TLS服务器
func NewTLSServer(addr string, config *tls.Config) *TLSServer {
    return &TLSServer{
        config: config,
        server: &http.Server{
            Addr:      addr,
            TLSConfig: config,
        },
    }
}

// Start 启动TLS服务器
func (s *TLSServer) Start() error {
    log.Printf("Starting TLS server on %s", s.server.Addr)
    return s.server.ListenAndServeTLS("", "")
}

// TLSServerHandler 处理TLS请求
func TLSServerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Secure connection established!\n")
    fmt.Fprintf(w, "Protocol: %s\n", r.TLS.Version)
    fmt.Fprintf(w, "Cipher Suite: %s\n", r.TLS.CipherSuite)
}
```

### 2.2.3 IPsec协议

**IPsec组件**：
1. **AH (Authentication Header)**：提供数据完整性
2. **ESP (Encapsulating Security Payload)**：提供机密性和完整性
3. **IKE (Internet Key Exchange)**：密钥管理

**IPsec模式**：
- **传输模式**：保护上层协议数据
- **隧道模式**：保护整个IP包

## 2.3 网络安全架构

### 2.3.1 纵深防御架构

**多层防护**：
```
┌─────────────────────────────────────┐
│           边界防护层                  │
│         (Firewall, IDS/IPS)         │
├─────────────────────────────────────┤
│           网络分段层                  │
│         (VLAN, DMZ)                 │
├─────────────────────────────────────┤
│           主机防护层                  │
│      (HIDS, Antivirus)             │
├─────────────────────────────────────┤
│           应用防护层                  │
│    (WAF, Input Validation)         │
├─────────────────────────────────────┤
│           数据防护层                  │
│      (Encryption, DLP)             │
└─────────────────────────────────────┘
```

### 2.3.2 零信任架构

**零信任原则**：
1. **永不信任，始终验证**
2. **最小权限原则**
3. **假设被攻破**

**零信任组件**：
- **身份验证**：多因子认证
- **设备验证**：设备健康检查
- **网络验证**：网络分段
- **应用验证**：应用访问控制

**Go语言零信任实现**：

```go
package main

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "log"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

// ZeroTrustPolicy 零信任策略
type ZeroTrustPolicy struct {
    UserID       string
    DeviceID     string
    ResourceID   string
    Permissions  []string
    ExpiresAt    time.Time
}

// ZeroTrustEngine 零信任引擎
type ZeroTrustEngine struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

// NewZeroTrustEngine 创建零信任引擎
func NewZeroTrustEngine() (*ZeroTrustEngine, error) {
    // 生成RSA密钥对
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key: %v", err)
    }

    return &ZeroTrustEngine{
        privateKey: privateKey,
        publicKey:  &privateKey.PublicKey,
    }, nil
}

// GenerateToken 生成访问令牌
func (e *ZeroTrustEngine) GenerateToken(policy ZeroTrustPolicy) (string, error) {
    claims := jwt.MapClaims{
        "user_id":     policy.UserID,
        "device_id":   policy.DeviceID,
        "resource_id": policy.ResourceID,
        "permissions": policy.Permissions,
        "exp":         policy.ExpiresAt.Unix(),
        "iat":         time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(e.privateKey)
}

// ValidateToken 验证访问令牌
func (e *ZeroTrustEngine) ValidateToken(tokenString string) (*ZeroTrustPolicy, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return e.publicKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %v", err)
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        policy := &ZeroTrustPolicy{
            UserID:     claims["user_id"].(string),
            DeviceID:   claims["device_id"].(string),
            ResourceID: claims["resource_id"].(string),
            ExpiresAt:  time.Unix(int64(claims["exp"].(float64)), 0),
        }

        // 检查权限
        if permissions, ok := claims["permissions"].([]interface{}); ok {
            for _, p := range permissions {
                policy.Permissions = append(policy.Permissions, p.(string))
            }
        }

        return policy, nil
    }

    return nil, fmt.Errorf("invalid token")
}

// AccessControl 访问控制
type AccessControl struct {
    engine *ZeroTrustEngine
}

// NewAccessControl 创建访问控制
func NewAccessControl(engine *ZeroTrustEngine) *AccessControl {
    return &AccessControl{engine: engine}
}

// CheckAccess 检查访问权限
func (ac *AccessControl) CheckAccess(tokenString, resourceID, action string) (bool, error) {
    policy, err := ac.engine.ValidateToken(tokenString)
    if err != nil {
        return false, err
    }

    // 检查资源ID
    if policy.ResourceID != resourceID {
        return false, fmt.Errorf("resource ID mismatch")
    }

    // 检查权限
    for _, permission := range policy.Permissions {
        if permission == action {
            return true, nil
        }
    }

    return false, fmt.Errorf("insufficient permissions")
}
```

## 2.4 网络安全监控

### 2.4.1 入侵检测系统 (IDS)

**IDS类型**：
1. **基于签名的检测**：匹配已知攻击模式
2. **基于异常的检测**：检测异常行为
3. **基于行为的检测**：分析用户行为模式

**Go语言IDS实现**：

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "regexp"
    "strings"
    "sync"
    "time"
)

// IDSRule IDS规则
type IDSRule struct {
    ID          string
    Pattern     string
    Description string
    Severity    string
    Action      string
}

// IDSEngine IDS引擎
type IDSEngine struct {
    rules    []IDSRule
    alerts   chan IDSAlert
    mu       sync.RWMutex
    patterns map[string]*regexp.Regexp
}

// IDSAlert IDS告警
type IDSAlert struct {
    RuleID      string
    SourceIP    string
    Destination string
    Payload     string
    Timestamp   time.Time
    Severity    string
}

// NewIDSEngine 创建IDS引擎
func NewIDSEngine() *IDSEngine {
    return &IDSEngine{
        alerts:   make(chan IDSAlert, 100),
        patterns: make(map[string]*regexp.Regexp),
    }
}

// AddRule 添加检测规则
func (e *IDSEngine) AddRule(rule IDSRule) error {
    e.mu.Lock()
    defer e.mu.Unlock()

    pattern, err := regexp.Compile(rule.Pattern)
    if err != nil {
        return fmt.Errorf("invalid pattern: %v", err)
    }

    e.rules = append(e.rules, rule)
    e.patterns[rule.ID] = pattern
    return nil
}

// AnalyzePacket 分析网络包
func (e *IDSEngine) AnalyzePacket(sourceIP, destination string, payload []byte) {
    e.mu.RLock()
    defer e.mu.RUnlock()

    payloadStr := string(payload)
    for _, rule := range e.rules {
        if pattern, exists := e.patterns[rule.ID]; exists {
            if pattern.MatchString(payloadStr) {
                alert := IDSAlert{
                    RuleID:      rule.ID,
                    SourceIP:    sourceIP,
                    Destination: destination,
                    Payload:     payloadStr,
                    Timestamp:   time.Now(),
                    Severity:    rule.Severity,
                }
                e.alerts <- alert
            }
        }
    }
}

// StartAlertHandler 启动告警处理器
func (e *IDSEngine) StartAlertHandler() {
    go func() {
        for alert := range e.alerts {
            log.Printf("[%s] IDS Alert: %s from %s to %s",
                alert.Severity, alert.RuleID, alert.SourceIP, alert.Destination)
        }
    }()
}

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
    engine *IDSEngine
    conn   net.Conn
}

// NewNetworkMonitor 创建网络监控器
func NewNetworkMonitor(engine *IDSEngine, conn net.Conn) *NetworkMonitor {
    return &NetworkMonitor{
        engine: engine,
        conn:   conn,
    }
}

// Monitor 监控网络流量
func (nm *NetworkMonitor) Monitor() {
    defer nm.conn.Close()

    scanner := bufio.NewScanner(nm.conn)
    for scanner.Scan() {
        line := scanner.Text()

        // 解析网络包信息
        parts := strings.Split(line, "|")
        if len(parts) >= 3 {
            sourceIP := parts[0]
            destination := parts[1]
            payload := parts[2]

            nm.engine.AnalyzePacket(sourceIP, destination, []byte(payload))
        }
    }
}
```

### 2.4.2 安全信息与事件管理 (SIEM)

**SIEM功能**：
1. **日志收集**：收集各种安全日志
2. **事件关联**：关联分析安全事件
3. **实时监控**：实时监控安全状态
4. **报告生成**：生成安全报告

## 2.5 网络安全测试

### 2.5.1 渗透测试

**渗透测试阶段**：
1. **信息收集**：收集目标信息
2. **漏洞扫描**：扫描系统漏洞
3. **漏洞利用**：利用发现的漏洞
4. **后渗透**：维持访问权限
5. **报告编写**：编写测试报告

### 2.5.2 安全评估

**评估方法**：
- **定性评估**：基于专家判断
- **定量评估**：基于数学模型
- **混合评估**：结合定性和定量

**风险评估公式**：
$$Risk = Threat \times Vulnerability \times Impact$$

## 2.6 网络安全最佳实践

### 2.6.1 网络分段

**分段原则**：
1. **按功能分段**：不同功能区域隔离
2. **按安全级别分段**：不同安全级别隔离
3. **按用户类型分段**：不同用户类型隔离

### 2.6.2 访问控制

**访问控制模型**：
- **DAC (Discretionary Access Control)**：自主访问控制
- **MAC (Mandatory Access Control)**：强制访问控制
- **RBAC (Role-Based Access Control)**：基于角色的访问控制

### 2.6.3 加密通信

**加密算法选择**：
- **对称加密**：AES-256-GCM
- **非对称加密**：RSA-2048, ECC-256
- **哈希算法**：SHA-256, SHA-3

## 2.7 网络安全趋势

### 2.7.1 新兴威胁

1. **AI驱动的攻击**：使用AI技术进行攻击
2. **供应链攻击**：攻击软件供应链
3. **量子计算威胁**：量子计算对加密的威胁

### 2.7.2 新兴技术

1. **区块链安全**：区块链在网络安全中的应用
2. **AI安全**：AI在网络安全中的应用
3. **零信任网络**：零信任架构的普及

---

**相关链接**：
- [1.1 安全基础理论](../01-Security-Foundations.md#1.1-安全基础理论)
- [3.1 应用安全](../03-Application-Security.md#3.1-应用安全基础)
- [4.1 数据安全](../04-Data-Security.md#4.1-数据安全基础)

**下一步**：
- [3. 应用安全](../03-Application-Security.md)
- [4. 数据安全](../04-Data-Security.md)
- [5. 密码学基础](../05-Cryptography.md)
