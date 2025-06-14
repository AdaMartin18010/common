# 08-网络安全 (Cybersecurity)

## 01-安全扫描工具 (Security Scanning Tools)

### 1. 概述

### 1.1 定义与目标

安全扫描工具是网络安全基础设施的核心组件，用于主动发现和识别系统中的安全漏洞、配置错误和潜在威胁。

**形式化定义**：
设 $S$ 为扫描目标系统，$V$ 为漏洞集合，$C$ 为配置集合，$T$ 为威胁集合，则安全扫描函数 $f$ 定义为：

$$f: S \rightarrow (V, C, T)$$

其中：

- $V = \{v_1, v_2, ..., v_n\}$ 为发现的漏洞集合
- $C = \{c_1, c_2, ..., c_m\}$ 为配置问题集合  
- $T = \{t_1, t_2, ..., t_k\}$ 为威胁集合

### 1.2 核心特性

1. **完整性**：$\forall s \in S, f(s) \neq \emptyset$
2. **准确性**：$P(f(s) = \text{true positive}) > 0.95$
3. **实时性**：$T(f(s)) < \text{threshold}$
4. **可扩展性**：$|S| \rightarrow \infty$ 时仍保持性能

### 2. 架构设计

### 2.1 分层架构

```text
┌─────────────────────────────────────┐
│           扫描调度层                  │
├─────────────────────────────────────┤
│           扫描引擎层                  │
├─────────────────────────────────────┤
│           插件管理层                  │
├─────────────────────────────────────┤
│           结果处理层                  │
└─────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 扫描调度器

```go
// ScannerScheduler 扫描调度器
type ScannerScheduler struct {
    queue       *PriorityQueue
    workers     []*ScannerWorker
    config      *ScannerConfig
    metrics     *ScannerMetrics
    mu          sync.RWMutex
}

// ScannerConfig 扫描配置
type ScannerConfig struct {
    MaxConcurrentScans int           `json:"max_concurrent_scans"`
    ScanTimeout        time.Duration `json:"scan_timeout"`
    RetryAttempts      int           `json:"retry_attempts"`
    RateLimit          int           `json:"rate_limit"`
}

// ScanRequest 扫描请求
type ScanRequest struct {
    ID          string                 `json:"id"`
    Target      string                 `json:"target"`
    ScanType    ScanType               `json:"scan_type"`
    Priority    Priority               `json:"priority"`
    Parameters  map[string]interface{} `json:"parameters"`
    CreatedAt   time.Time              `json:"created_at"`
    ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

// ScanType 扫描类型
type ScanType string

const (
    ScanTypeVulnerability ScanType = "vulnerability"
    ScanTypePort          ScanType = "port"
    ScanTypeWeb           ScanType = "web"
    ScanTypeNetwork       ScanType = "network"
    ScanTypeConfiguration ScanType = "configuration"
)

// Priority 优先级
type Priority int

const (
    PriorityLow    Priority = 1
    PriorityMedium Priority = 2
    PriorityHigh   Priority = 3
    PriorityCritical Priority = 4
)

// ScheduleScan 调度扫描
func (s *ScannerScheduler) ScheduleScan(req *ScanRequest) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 验证请求
    if err := s.validateRequest(req); err != nil {
        return fmt.Errorf("invalid scan request: %w", err)
    }
    
    // 添加到队列
    s.queue.Push(req)
    
    // 更新指标
    s.metrics.IncQueuedScans()
    
    return nil
}

// validateRequest 验证扫描请求
func (s *ScannerScheduler) validateRequest(req *ScanRequest) error {
    if req.Target == "" {
        return errors.New("target is required")
    }
    
    if req.ScanType == "" {
        return errors.New("scan type is required")
    }
    
    // 验证目标格式
    if !s.isValidTarget(req.Target) {
        return errors.New("invalid target format")
    }
    
    return nil
}

// isValidTarget 验证目标格式
func (s *ScannerScheduler) isValidTarget(target string) bool {
    // IP地址验证
    if net.ParseIP(target) != nil {
        return true
    }
    
    // 域名验证
    if strings.Contains(target, ".") && !strings.Contains(target, ":") {
        return true
    }
    
    // URL验证
    if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
        return true
    }
    
    return false
}
```

#### 2.2.2 扫描引擎

```go
// ScannerEngine 扫描引擎
type ScannerEngine struct {
    plugins     map[ScanType][]ScannerPlugin
    config      *EngineConfig
    results     chan *ScanResult
    errors      chan *ScanError
    ctx         context.Context
    cancel      context.CancelFunc
}

// EngineConfig 引擎配置
type EngineConfig struct {
    MaxConcurrency int           `json:"max_concurrency"`
    Timeout        time.Duration `json:"timeout"`
    RetryCount     int           `json:"retry_count"`
    RateLimit      int           `json:"rate_limit"`
}

// ScannerPlugin 扫描插件接口
type ScannerPlugin interface {
    Name() string
    Version() string
    Scan(ctx context.Context, target string, params map[string]interface{}) (*ScanResult, error)
    Validate(target string) error
}

// ScanResult 扫描结果
type ScanResult struct {
    ID          string                 `json:"id"`
    Target      string                 `json:"target"`
    ScanType    ScanType               `json:"scan_type"`
    Plugin      string                 `json:"plugin"`
    Findings    []Finding              `json:"findings"`
    Metadata    map[string]interface{} `json:"metadata"`
    StartedAt   time.Time              `json:"started_at"`
    CompletedAt time.Time              `json:"completed_at"`
    Duration    time.Duration          `json:"duration"`
}

// Finding 发现结果
type Finding struct {
    ID          string                 `json:"id"`
    Type        FindingType            `json:"type"`
    Severity    Severity               `json:"severity"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Evidence    string                 `json:"evidence"`
    Location    string                 `json:"location"`
    CVE         []string               `json:"cve,omitempty"`
    CVSS        *CVSSScore             `json:"cvss,omitempty"`
    Remediation string                 `json:"remediation"`
    Tags        []string               `json:"tags"`
    CreatedAt   time.Time              `json:"created_at"`
}

// FindingType 发现类型
type FindingType string

const (
    FindingTypeVulnerability FindingType = "vulnerability"
    FindingTypeMisconfiguration FindingType = "misconfiguration"
    FindingTypeInformation   FindingType = "information"
    FindingTypeWarning       FindingType = "warning"
)

// Severity 严重程度
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
    SeverityInfo     Severity = "info"
)

// CVSSScore CVSS评分
type CVSSScore struct {
    BaseScore    float64 `json:"base_score"`
    TemporalScore float64 `json:"temporal_score"`
    EnvironmentalScore float64 `json:"environmental_score"`
    Vector       string  `json:"vector"`
}

// ExecuteScan 执行扫描
func (e *ScannerEngine) ExecuteScan(req *ScanRequest) (*ScanResult, error) {
    ctx, cancel := context.WithTimeout(e.ctx, e.config.Timeout)
    defer cancel()
    
    // 获取插件
    plugins, exists := e.plugins[req.ScanType]
    if !exists {
        return nil, fmt.Errorf("no plugins found for scan type: %s", req.ScanType)
    }
    
    // 创建结果
    result := &ScanResult{
        ID:        req.ID,
        Target:    req.Target,
        ScanType:  req.ScanType,
        StartedAt: time.Now(),
        Findings:  make([]Finding, 0),
        Metadata:  make(map[string]interface{}),
    }
    
    // 并发执行插件
    var wg sync.WaitGroup
    resultChan := make(chan *ScanResult, len(plugins))
    errorChan := make(chan error, len(plugins))
    
    for _, plugin := range plugins {
        wg.Add(1)
        go func(p ScannerPlugin) {
            defer wg.Done()
            
            // 验证目标
            if err := p.Validate(req.Target); err != nil {
                errorChan <- fmt.Errorf("plugin %s validation failed: %w", p.Name(), err)
                return
            }
            
            // 执行扫描
            pluginResult, err := p.Scan(ctx, req.Target, req.Parameters)
            if err != nil {
                errorChan <- fmt.Errorf("plugin %s scan failed: %w", p.Name(), err)
                return
            }
            
            resultChan <- pluginResult
        }(plugin)
    }
    
    // 等待所有插件完成
    go func() {
        wg.Wait()
        close(resultChan)
        close(errorChan)
    }()
    
    // 收集结果
    for pluginResult := range resultChan {
        result.Findings = append(result.Findings, pluginResult.Findings...)
        for k, v := range pluginResult.Metadata {
            result.Metadata[k] = v
        }
    }
    
    // 检查错误
    var errors []error
    for err := range errorChan {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        return result, fmt.Errorf("scan completed with errors: %v", errors)
    }
    
    result.CompletedAt = time.Now()
    result.Duration = result.CompletedAt.Sub(result.StartedAt)
    
    return result, nil
}
```

### 2.3 插件系统

#### 2.3.1 漏洞扫描插件

```go
// VulnerabilityScanner 漏洞扫描插件
type VulnerabilityScanner struct {
    cveDB       *CVEDatabase
    signatures  []VulnerabilitySignature
    config      *VulnScannerConfig
}

// VulnScannerConfig 漏洞扫描配置
type VulnScannerConfig struct {
    EnableCVE     bool     `json:"enable_cve"`
    EnableExploit bool     `json:"enable_exploit"`
    SeverityFilter []Severity `json:"severity_filter"`
    Timeout       time.Duration `json:"timeout"`
}

// VulnerabilitySignature 漏洞特征
type VulnerabilitySignature struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Pattern     string                 `json:"pattern"`
    Severity    Severity               `json:"severity"`
    CVEs        []string               `json:"cves"`
    Exploit     *ExploitInfo           `json:"exploit,omitempty"`
}

// ExploitInfo 利用信息
type ExploitInfo struct {
    Available bool   `json:"available"`
    Type      string `json:"type"`
    Complexity string `json:"complexity"`
    Impact    string `json:"impact"`
}

// Scan 执行漏洞扫描
func (v *VulnerabilityScanner) Scan(ctx context.Context, target string, params map[string]interface{}) (*ScanResult, error) {
    result := &ScanResult{
        Target:   target,
        ScanType: ScanTypeVulnerability,
        Plugin:   v.Name(),
        Findings: make([]Finding, 0),
        Metadata: make(map[string]interface{}),
    }
    
    // 端口扫描
    openPorts, err := v.scanPorts(ctx, target)
    if err != nil {
        return nil, fmt.Errorf("port scan failed: %w", err)
    }
    
    result.Metadata["open_ports"] = openPorts
    
    // 服务识别
    services, err := v.identifyServices(ctx, target, openPorts)
    if err != nil {
        return nil, fmt.Errorf("service identification failed: %w", err)
    }
    
    result.Metadata["services"] = services
    
    // 漏洞检测
    for _, service := range services {
        findings, err := v.detectVulnerabilities(ctx, target, service)
        if err != nil {
            continue // 继续检测其他服务
        }
        
        result.Findings = append(result.Findings, findings...)
    }
    
    return result, nil
}

// scanPorts 端口扫描
func (v *VulnerabilityScanner) scanPorts(ctx context.Context, target string) ([]int, error) {
    var openPorts []int
    
    // 常见端口列表
    commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 3306, 5432, 6379, 27017}
    
    for _, port := range commonPorts {
        select {
        case <-ctx.Done():
            return openPorts, ctx.Err()
        default:
            if v.isPortOpen(target, port) {
                openPorts = append(openPorts, port)
            }
        }
    }
    
    return openPorts, nil
}

// isPortOpen 检查端口是否开放
func (v *VulnerabilityScanner) isPortOpen(target string, port int) bool {
    timeout := time.Duration(v.config.Timeout)
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, port), timeout)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

// identifyServices 识别服务
func (v *VulnerabilityScanner) identifyServices(ctx context.Context, target string, ports []int) ([]Service, error) {
    var services []Service
    
    for _, port := range ports {
        select {
        case <-ctx.Done():
            return services, ctx.Err()
        default:
            service, err := v.identifyService(target, port)
            if err != nil {
                continue
            }
            services = append(services, service)
        }
    }
    
    return services, nil
}

// Service 服务信息
type Service struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    Name        string `json:"name"`
    Version     string `json:"version"`
    Banner      string `json:"banner"`
}

// identifyService 识别单个服务
func (v *VulnerabilityScanner) identifyService(target string, port int) (Service, error) {
    service := Service{
        Port:     port,
        Protocol: "tcp",
    }
    
    // 连接并获取banner
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, port), v.config.Timeout)
    if err != nil {
        return service, err
    }
    defer conn.Close()
    
    // 发送探测数据
    probes := []string{
        "\r\n\r\n",
        "GET / HTTP/1.0\r\n\r\n",
        "SSH-2.0-OpenSSH_8.0\r\n",
    }
    
    for _, probe := range probes {
        conn.SetWriteDeadline(time.Now().Add(v.config.Timeout))
        _, err := conn.Write([]byte(probe))
        if err != nil {
            continue
        }
        
        // 读取响应
        conn.SetReadDeadline(time.Now().Add(v.config.Timeout))
        buffer := make([]byte, 1024)
        n, err := conn.Read(buffer)
        if err != nil {
            continue
        }
        
        response := string(buffer[:n])
        service.Banner = response
        
        // 基于响应识别服务
        if strings.Contains(response, "SSH") {
            service.Name = "SSH"
            service.Version = v.extractVersion(response, "SSH")
        } else if strings.Contains(response, "HTTP") {
            service.Name = "HTTP"
            service.Version = v.extractVersion(response, "HTTP")
        } else if strings.Contains(response, "FTP") {
            service.Name = "FTP"
            service.Version = v.extractVersion(response, "FTP")
        }
        
        break
    }
    
    return service, nil
}

// extractVersion 提取版本信息
func (v *VulnerabilityScanner) extractVersion(banner, service string) string {
    re := regexp.MustCompile(fmt.Sprintf(`%s/(\d+\.\d+(?:\.\d+)?)`, service))
    matches := re.FindStringSubmatch(banner)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}

// detectVulnerabilities 检测漏洞
func (v *VulnerabilityScanner) detectVulnerabilities(ctx context.Context, target string, service Service) ([]Finding, error) {
    var findings []Finding
    
    // 基于服务类型和版本检测漏洞
    if service.Name == "SSH" {
        sshFindings, err := v.detectSSHVulnerabilities(target, service)
        if err == nil {
            findings = append(findings, sshFindings...)
        }
    } else if service.Name == "HTTP" {
        httpFindings, err := v.detectHTTPVulnerabilities(target, service)
        if err == nil {
            findings = append(findings, httpFindings...)
        }
    }
    
    return findings, nil
}

// detectSSHVulnerabilities 检测SSH漏洞
func (v *VulnerabilityScanner) detectSSHVulnerabilities(target string, service Service) ([]Finding, error) {
    var findings []Finding
    
    // 检查弱加密算法
    if v.hasWeakCrypto(service.Banner) {
        findings = append(findings, Finding{
            ID:          uuid.New().String(),
            Type:        FindingTypeVulnerability,
            Severity:    SeverityMedium,
            Title:       "Weak SSH Encryption Algorithm",
            Description: "SSH service is using weak encryption algorithms",
            Evidence:    service.Banner,
            Location:    fmt.Sprintf("%s:%d", target, service.Port),
            Remediation: "Configure SSH to use strong encryption algorithms only",
            Tags:        []string{"ssh", "encryption", "weak-crypto"},
            CreatedAt:   time.Now(),
        })
    }
    
    // 检查默认端口
    if service.Port == 22 {
        findings = append(findings, Finding{
            ID:          uuid.New().String(),
            Type:        FindingTypeInformation,
            Severity:    SeverityInfo,
            Title:       "SSH Running on Default Port",
            Description: "SSH service is running on the default port 22",
            Evidence:    fmt.Sprintf("Port %d", service.Port),
            Location:    fmt.Sprintf("%s:%d", target, service.Port),
            Remediation: "Consider changing SSH to a non-standard port",
            Tags:        []string{"ssh", "default-port"},
            CreatedAt:   time.Now(),
        })
    }
    
    return findings, nil
}

// hasWeakCrypto 检查是否使用弱加密
func (v *VulnerabilityScanner) hasWeakCrypto(banner string) bool {
    weakAlgorithms := []string{"des", "3des", "md5", "sha1"}
    bannerLower := strings.ToLower(banner)
    
    for _, algo := range weakAlgorithms {
        if strings.Contains(bannerLower, algo) {
            return true
        }
    }
    
    return false
}

// detectHTTPVulnerabilities 检测HTTP漏洞
func (v *VulnerabilityScanner) detectHTTPVulnerabilities(target string, service Service) ([]Finding, error) {
    var findings []Finding
    
    // 检查HTTP版本
    if strings.Contains(service.Banner, "HTTP/1.0") {
        findings = append(findings, Finding{
            ID:          uuid.New().String(),
            Type:        FindingTypeVulnerability,
            Severity:    SeverityMedium,
            Title:       "Outdated HTTP Version",
            Description: "Server is using HTTP/1.0 which lacks security features",
            Evidence:    service.Banner,
            Location:    fmt.Sprintf("%s:%d", target, service.Port),
            Remediation: "Upgrade to HTTP/1.1 or HTTP/2",
            Tags:        []string{"http", "version", "outdated"},
            CreatedAt:   time.Now(),
        })
    }
    
    return findings, nil
}

// Name 插件名称
func (v *VulnerabilityScanner) Name() string {
    return "vulnerability-scanner"
}

// Version 插件版本
func (v *VulnerabilityScanner) Version() string {
    return "1.0.0"
}

// Validate 验证目标
func (v *VulnerabilityScanner) Validate(target string) error {
    if target == "" {
        return errors.New("target is required")
    }
    
    // 验证IP地址或域名格式
    if net.ParseIP(target) == nil && !strings.Contains(target, ".") {
        return errors.New("invalid target format")
    }
    
    return nil
}
```

## 3. 数学建模

### 3.1 漏洞评分模型

**CVSS评分计算**：

$$CVSS_{Base} = \min(10, \text{Impact} + \text{Exploitability})$$

其中：

- $\text{Impact} = 10.41 \times (1 - (1 - \text{ConfImpact}) \times (1 - \text{IntegImpact}) \times (1 - \text{AvailImpact}))$
- $\text{Exploitability} = 20 \times \text{AccessVector} \times \text{AccessComplexity} \times \text{Authentication}$

### 3.2 扫描效率模型

**扫描时间估计**：

$$T_{scan} = \sum_{i=1}^{n} (T_{connect_i} + T_{probe_i} + T_{analyze_i})$$

其中：

- $T_{connect_i}$ 为连接时间
- $T_{probe_i}$ 为探测时间
- $T_{analyze_i}$ 为分析时间

### 3.3 误报率模型

**精确度计算**：

$$Precision = \frac{TP}{TP + FP}$$

**召回率计算**：

$$Recall = \frac{TP}{TP + FN}$$

其中：

- $TP$ 为真正例
- $FP$ 为假正例
- $FN$ 为假负例

## 4. 性能优化

### 4.1 并发控制

```go
// RateLimiter 速率限制器
type RateLimiter struct {
    tokens     chan struct{}
    rate       int
    burst      int
    interval   time.Duration
    lastRefill time.Time
    mu         sync.Mutex
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate int, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, burst),
        rate:     rate,
        burst:    burst,
        interval: time.Second / time.Duration(rate),
    }
    
    // 初始化令牌
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }
    
    return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

// Wait 等待令牌可用
func (rl *RateLimiter) Wait(ctx context.Context) error {
    rl.refill()
    
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// refill 补充令牌
func (rl *RateLimiter) refill() {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(rl.lastRefill)
    tokensToAdd := int(elapsed / rl.interval)
    
    if tokensToAdd > 0 {
        for i := 0; i < tokensToAdd && len(rl.tokens) < rl.burst; i++ {
            select {
            case rl.tokens <- struct{}{}:
            default:
                break
            }
        }
        rl.lastRefill = now
    }
}
```

### 4.2 缓存机制

```go
// ScanCache 扫描缓存
type ScanCache struct {
    cache    *lru.Cache
    ttl      time.Duration
    maxSize  int
}

// CacheEntry 缓存条目
type CacheEntry struct {
    Result    *ScanResult
    ExpiresAt time.Time
}

// NewScanCache 创建扫描缓存
func NewScanCache(maxSize int, ttl time.Duration) (*ScanCache, error) {
    cache, err := lru.New(maxSize)
    if err != nil {
        return nil, err
    }
    
    return &ScanCache{
        cache:   cache,
        ttl:     ttl,
        maxSize: maxSize,
    }, nil
}

// Get 获取缓存结果
func (c *ScanCache) Get(key string) (*ScanResult, bool) {
    if entry, found := c.cache.Get(key); found {
        cacheEntry := entry.(*CacheEntry)
        if time.Now().Before(cacheEntry.ExpiresAt) {
            return cacheEntry.Result, true
        } else {
            c.cache.Remove(key)
        }
    }
    return nil, false
}

// Set 设置缓存结果
func (c *ScanCache) Set(key string, result *ScanResult) {
    entry := &CacheEntry{
        Result:    result,
        ExpiresAt: time.Now().Add(c.ttl),
    }
    c.cache.Add(key, entry)
}
```

## 5. 监控与指标

### 5.1 性能指标

```go
// ScannerMetrics 扫描指标
type ScannerMetrics struct {
    totalScans       int64
    successfulScans  int64
    failedScans      int64
    scanDuration     time.Duration
    findingsCount    int64
    falsePositives   int64
    mu               sync.RWMutex
}

// IncTotalScans 增加总扫描数
func (m *ScannerMetrics) IncTotalScans() {
    atomic.AddInt64(&m.totalScans, 1)
}

// IncSuccessfulScans 增加成功扫描数
func (m *ScannerMetrics) IncSuccessfulScans() {
    atomic.AddInt64(&m.successfulScans, 1)
}

// IncFailedScans 增加失败扫描数
func (m *ScannerMetrics) IncFailedScans() {
    atomic.AddInt64(&m.failedScans, 1)
}

// RecordScanDuration 记录扫描时长
func (m *ScannerMetrics) RecordScanDuration(duration time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.scanDuration = duration
}

// GetMetrics 获取指标
func (m *ScannerMetrics) GetMetrics() map[string]interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    successRate := float64(0)
    if m.totalScans > 0 {
        successRate = float64(m.successfulScans) / float64(m.totalScans) * 100
    }
    
    return map[string]interface{}{
        "total_scans":      m.totalScans,
        "successful_scans": m.successfulScans,
        "failed_scans":     m.failedScans,
        "success_rate":     successRate,
        "avg_duration":     m.scanDuration,
        "findings_count":   m.findingsCount,
        "false_positives":  m.falsePositives,
    }
}
```

## 6. 安全考虑

### 6.1 权限控制

```go
// PermissionManager 权限管理器
type PermissionManager struct {
    policies map[string]*ScanPolicy
    mu       sync.RWMutex
}

// ScanPolicy 扫描策略
type ScanPolicy struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Targets     []string `json:"targets"`
    ScanTypes   []ScanType `json:"scan_types"`
    Permissions []Permission `json:"permissions"`
    CreatedBy   string   `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}

// Permission 权限
type Permission struct {
    Resource string   `json:"resource"`
    Actions  []string `json:"actions"`
}

// CheckPermission 检查权限
func (pm *PermissionManager) CheckPermission(userID, target string, scanType ScanType) error {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    // 检查用户策略
    for _, policy := range pm.policies {
        if pm.userHasPolicy(userID, policy.ID) {
            if pm.policyAllows(policy, target, scanType) {
                return nil
            }
        }
    }
    
    return errors.New("permission denied")
}

// userHasPolicy 检查用户是否有策略
func (pm *PermissionManager) userHasPolicy(userID, policyID string) bool {
    // 实现用户策略检查逻辑
    return true // 简化实现
}

// policyAllows 检查策略是否允许
func (pm *PermissionManager) policyAllows(policy *ScanPolicy, target string, scanType ScanType) bool {
    // 检查目标
    targetAllowed := false
    for _, allowedTarget := range policy.Targets {
        if pm.matchesTarget(target, allowedTarget) {
            targetAllowed = true
            break
        }
    }
    
    if !targetAllowed {
        return false
    }
    
    // 检查扫描类型
    typeAllowed := false
    for _, allowedType := range policy.ScanTypes {
        if scanType == allowedType {
            typeAllowed = true
            break
        }
    }
    
    return typeAllowed
}

// matchesTarget 检查目标是否匹配
func (pm *PermissionManager) matchesTarget(target, pattern string) bool {
    // 支持通配符匹配
    if strings.Contains(pattern, "*") {
        regexPattern := strings.ReplaceAll(pattern, "*", ".*")
        matched, _ := regexp.MatchString(regexPattern, target)
        return matched
    }
    
    return target == pattern
}
```

## 7. 总结

安全扫描工具是网络安全基础设施的重要组成部分，通过系统化的漏洞发现和风险评估，为组织提供主动的安全防护能力。本模块提供了：

1. **完整的架构设计**：分层架构、插件系统、并发控制
2. **形式化数学建模**：CVSS评分、扫描效率、误报率计算
3. **高性能实现**：速率限制、缓存机制、指标监控
4. **安全考虑**：权限控制、输入验证、错误处理

通过Go语言的内存安全和并发特性，实现了高效、可靠的安全扫描工具，为网络安全防护提供了强有力的技术支撑。

---

**相关链接**：

- [02-入侵检测系统](../02-Intrusion-Detection-System/README.md)
- [03-加密服务](../03-Encryption-Services/README.md)
- [04-身份认证](../04-Identity-Authentication/README.md)
