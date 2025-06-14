# 01-安全扫描工具 (Security Scanning Tools)

## 目录

- [01-安全扫描工具 (Security Scanning Tools)](#01-安全扫描工具-security-scanning-tools)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 定义](#11-定义)
    - [1.2 核心功能](#12-核心功能)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 扫描空间定义](#21-扫描空间定义)
    - [2.2 扫描函数](#22-扫描函数)
    - [2.3 风险评估模型](#23-风险评估模型)
  - [3. 核心架构](#3-核心架构)
    - [3.1 系统架构](#31-系统架构)
    - [3.2 扫描配置](#32-扫描配置)
  - [4. 扫描引擎](#4-扫描引擎)
    - [4.1 基础扫描引擎](#41-基础扫描引擎)
    - [4.2 端口扫描器](#42-端口扫描器)
  - [5. 漏洞检测](#5-漏洞检测)
    - [5.1 漏洞检测引擎](#51-漏洞检测引擎)
  - [6. 端口扫描](#6-端口扫描)
    - [6.1 高级端口扫描](#61-高级端口扫描)
  - [7. 服务识别](#7-服务识别)
    - [7.1 服务指纹识别](#71-服务指纹识别)
  - [8. 安全评估](#8-安全评估)
    - [8.1 风险评估引擎](#81-风险评估引擎)
  - [9. 报告生成](#9-报告生成)
    - [9.1 报告生成器](#91-报告生成器)
  - [10. 性能优化](#10-性能优化)
    - [10.1 并发优化](#101-并发优化)
  - [11. 部署运维](#11-部署运维)
    - [11.1 容器化部署](#111-容器化部署)
    - [11.2 Kubernetes部署](#112-kubernetes部署)
    - [11.3 监控配置](#113-监控配置)

## 1. 概述

### 1.1 定义

安全扫描工具是一种自动化安全评估系统，用于识别网络、系统和应用程序中的安全漏洞、配置错误和潜在威胁。

### 1.2 核心功能

- **漏洞扫描**: 识别已知安全漏洞
- **端口扫描**: 发现开放端口和服务
- **服务识别**: 识别运行的服务和版本
- **配置审计**: 检查安全配置合规性
- **威胁评估**: 评估安全风险等级

## 2. 形式化定义

### 2.1 扫描空间定义

设 $S$ 为扫描空间，$T$ 为目标集合，$V$ 为漏洞集合：

$$S = \{T, V, C, R\}$$

其中：

- $T = \{t_1, t_2, ..., t_n\}$ 为目标集合
- $V = \{v_1, v_2, ..., v_m\}$ 为漏洞集合
- $C$ 为扫描配置
- $R$ 为扫描结果

### 2.2 扫描函数

扫描函数 $f: T \times C \rightarrow R$ 定义为：

$$f(t, c) = \{(v, s, d) | v \in V, s \in [0,1], d \in D\}$$

其中：

- $v$ 为漏洞标识
- $s$ 为严重程度 (0-1)
- $d$ 为详细描述

### 2.3 风险评估模型

风险评分函数 $R: V \times T \rightarrow \mathbb{R}$：

$$R(v, t) = \alpha \cdot S(v) + \beta \cdot E(v, t) + \gamma \cdot I(v, t)$$

其中：

- $S(v)$ 为漏洞严重程度
- $E(v, t)$ 为利用难度
- $I(v, t)$ 为影响程度
- $\alpha, \beta, \gamma$ 为权重系数

## 3. 核心架构

### 3.1 系统架构

```go
// 扫描引擎核心接口
type Scanner interface {
    Scan(ctx context.Context, target Target) ([]Vulnerability, error)
    GetCapabilities() []Capability
    GetVersion() string
}

// 目标定义
type Target struct {
    ID          string            `json:"id"`
    Host        string            `json:"host"`
    Ports       []int             `json:"ports"`
    Protocols   []string          `json:"protocols"`
    Metadata    map[string]string `json:"metadata"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

// 漏洞定义
type Vulnerability struct {
    ID          string                 `json:"id"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Severity    SeverityLevel          `json:"severity"`
    CVE         string                 `json:"cve,omitempty"`
    CVSS        float64                `json:"cvss,omitempty"`
    Category    string                 `json:"category"`
    References  []string               `json:"references"`
    Evidence    map[string]interface{} `json:"evidence"`
    DiscoveredAt time.Time             `json:"discovered_at"`
}

// 严重程度枚举
type SeverityLevel int

const (
    SeverityInfo SeverityLevel = iota
    SeverityLow
    SeverityMedium
    SeverityHigh
    SeverityCritical
)

func (s SeverityLevel) String() string {
    switch s {
    case SeverityInfo:
        return "Info"
    case SeverityLow:
        return "Low"
    case SeverityMedium:
        return "Medium"
    case SeverityHigh:
        return "High"
    case SeverityCritical:
        return "Critical"
    default:
        return "Unknown"
    }
}
```

### 3.2 扫描配置

```go
// 扫描配置
type ScanConfig struct {
    Timeout           time.Duration     `json:"timeout"`
    MaxConcurrency    int               `json:"max_concurrency"`
    RateLimit         int               `json:"rate_limit"`
    UserAgent         string            `json:"user_agent"`
    Headers           map[string]string `json:"headers"`
    Proxy             string            `json:"proxy,omitempty"`
    Credentials       *Credentials      `json:"credentials,omitempty"`
    ScanTypes         []ScanType        `json:"scan_types"`
    ExcludePatterns   []string          `json:"exclude_patterns"`
    IncludePatterns   []string          `json:"include_patterns"`
}

// 扫描类型
type ScanType string

const (
    ScanTypePort      ScanType = "port"
    ScanTypeVulnerability ScanType = "vulnerability"
    ScanTypeService   ScanType = "service"
    ScanTypeConfig    ScanType = "config"
    ScanTypeWeb       ScanType = "web"
)

// 认证信息
type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Token    string `json:"token,omitempty"`
}
```

## 4. 扫描引擎

### 4.1 基础扫描引擎

```go
// 基础扫描引擎
type BaseScanner struct {
    config     *ScanConfig
    plugins    map[string]Plugin
    results    chan ScanResult
    errors     chan error
    ctx        context.Context
    cancel     context.CancelFunc
}

// 扫描结果
type ScanResult struct {
    Target        Target         `json:"target"`
    Vulnerabilities []Vulnerability `json:"vulnerabilities"`
    Services      []Service      `json:"services"`
    ScanTime      time.Duration  `json:"scan_time"`
    CompletedAt   time.Time      `json:"completed_at"`
}

// 服务信息
type Service struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    Name        string `json:"name"`
    Version     string `json:"version,omitempty"`
    Banner      string `json:"banner,omitempty"`
    SSL         bool   `json:"ssl"`
}

// 插件接口
type Plugin interface {
    Name() string
    Version() string
    Scan(ctx context.Context, target Target, config *ScanConfig) ([]Vulnerability, error)
    GetCapabilities() []Capability
}

// 能力定义
type Capability struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    ScanTypes   []ScanType `json:"scan_types"`
    Protocols   []string `json:"protocols"`
}

// 创建基础扫描引擎
func NewBaseScanner(config *ScanConfig) *BaseScanner {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &BaseScanner{
        config:  config,
        plugins: make(map[string]Plugin),
        results: make(chan ScanResult, 100),
        errors:  make(chan error, 100),
        ctx:     ctx,
        cancel:  cancel,
    }
}

// 注册插件
func (s *BaseScanner) RegisterPlugin(plugin Plugin) {
    s.plugins[plugin.Name()] = plugin
}

// 执行扫描
func (s *BaseScanner) Scan(targets []Target) error {
    // 创建工作池
    workerPool := make(chan struct{}, s.config.MaxConcurrency)
    
    // 启动扫描任务
    for _, target := range targets {
        select {
        case <-s.ctx.Done():
            return s.ctx.Err()
        case workerPool <- struct{}{}:
            go func(t Target) {
                defer func() { <-workerPool }()
                s.scanTarget(t)
            }(target)
        }
    }
    
    return nil
}

// 扫描单个目标
func (s *BaseScanner) scanTarget(target Target) {
    start := time.Now()
    var allVulns []Vulnerability
    var allServices []Service
    
    // 执行各种类型的扫描
    for _, scanType := range s.config.ScanTypes {
        switch scanType {
        case ScanTypePort:
            services := s.scanPorts(target)
            allServices = append(allServices, services...)
        case ScanTypeVulnerability:
            vulns := s.scanVulnerabilities(target)
            allVulns = append(allVulns, vulns...)
        case ScanTypeService:
            services := s.scanServices(target)
            allServices = append(allServices, services...)
        }
    }
    
    result := ScanResult{
        Target:         target,
        Vulnerabilities: allVulns,
        Services:       allServices,
        ScanTime:       time.Since(start),
        CompletedAt:    time.Now(),
    }
    
    select {
    case s.results <- result:
    case <-s.ctx.Done():
    }
}
```

### 4.2 端口扫描器

```go
// 端口扫描器
type PortScanner struct {
    timeout time.Duration
    workers int
}

// 创建端口扫描器
func NewPortScanner(timeout time.Duration, workers int) *PortScanner {
    return &PortScanner{
        timeout: timeout,
        workers: workers,
    }
}

// 扫描端口
func (ps *PortScanner) ScanPorts(target Target) []Service {
    var services []Service
    var wg sync.WaitGroup
    results := make(chan Service, 100)
    
    // 常见端口列表
    commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 3306, 5432, 6379, 27017}
    
    // 启动工作协程
    for i := 0; i < ps.workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for _, port := range commonPorts {
                if service := ps.scanPort(target.Host, port); service != nil {
                    results <- *service
                }
            }
        }()
    }
    
    // 等待所有工作完成
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集结果
    for service := range results {
        services = append(services, service)
    }
    
    return services
}

// 扫描单个端口
func (ps *PortScanner) scanPort(host string, port int) *Service {
    address := fmt.Sprintf("%s:%d", host, port)
    
    conn, err := net.DialTimeout("tcp", address, ps.timeout)
    if err != nil {
        return nil
    }
    defer conn.Close()
    
    // 尝试获取服务banner
    banner := ps.getBanner(conn, port)
    
    return &Service{
        Port:     port,
        Protocol: "tcp",
        Name:     ps.identifyService(port, banner),
        Banner:   banner,
        SSL:      port == 443 || port == 993 || port == 995,
    }
}

// 获取服务banner
func (ps *PortScanner) getBanner(conn net.Conn, port int) string {
    // 发送探测数据
    var probe string
    switch port {
    case 22:
        probe = "\n"
    case 80, 443:
        probe = "GET / HTTP/1.0\r\n\r\n"
    case 21:
        probe = "QUIT\r\n"
    default:
        probe = "\n"
    }
    
    conn.SetWriteDeadline(time.Now().Add(ps.timeout))
    conn.Write([]byte(probe))
    
    // 读取响应
    conn.SetReadDeadline(time.Now().Add(ps.timeout))
    buffer := make([]byte, 1024)
    n, _ := conn.Read(buffer)
    
    if n > 0 {
        return strings.TrimSpace(string(buffer[:n]))
    }
    
    return ""
}

// 识别服务
func (ps *PortScanner) identifyService(port int, banner string) string {
    // 基于端口和banner识别服务
    switch port {
    case 22:
        return "SSH"
    case 80, 443:
        return "HTTP"
    case 21:
        return "FTP"
    case 25:
        return "SMTP"
    case 53:
        return "DNS"
    case 3306:
        return "MySQL"
    case 5432:
        return "PostgreSQL"
    case 6379:
        return "Redis"
    case 27017:
        return "MongoDB"
    default:
        // 基于banner内容识别
        if strings.Contains(banner, "SSH") {
            return "SSH"
        }
        if strings.Contains(banner, "HTTP") {
            return "HTTP"
        }
        return "Unknown"
    }
}
```

## 5. 漏洞检测

### 5.1 漏洞检测引擎

```go
// 漏洞检测引擎
type VulnerabilityScanner struct {
    rules       []VulnerabilityRule
    signatures  map[string]Signature
    database    *VulnDatabase
}

// 漏洞规则
type VulnerabilityRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Category    string                 `json:"category"`
    Severity    SeverityLevel          `json:"severity"`
    Conditions  []Condition            `json:"conditions"`
    Actions     []Action               `json:"actions"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// 检测条件
type Condition struct {
    Type        string `json:"type"`
    Field       string `json:"field"`
    Operator    string `json:"operator"`
    Value       string `json:"value"`
    Negate      bool   `json:"negate"`
}

// 检测动作
type Action struct {
    Type    string                 `json:"type"`
    Params  map[string]interface{} `json:"params"`
}

// 创建漏洞检测引擎
func NewVulnerabilityScanner() *VulnerabilityScanner {
    return &VulnerabilityScanner{
        rules:      make([]VulnerabilityRule, 0),
        signatures: make(map[string]Signature),
        database:   NewVulnDatabase(),
    }
}

// 加载规则
func (vs *VulnerabilityScanner) LoadRules(rules []VulnerabilityRule) {
    vs.rules = append(vs.rules, rules...)
}

// 检测漏洞
func (vs *VulnerabilityScanner) ScanVulnerabilities(target Target, services []Service) []Vulnerability {
    var vulnerabilities []Vulnerability
    
    // 对每个服务进行漏洞检测
    for _, service := range services {
        vulns := vs.scanService(target, service)
        vulnerabilities = append(vulnerabilities, vulns...)
    }
    
    // 执行通用规则检测
    vulns := vs.executeRules(target, services)
    vulnerabilities = append(vulnerabilities, vulns...)
    
    return vulnerabilities
}

// 扫描单个服务
func (vs *VulnerabilityScanner) scanService(target Target, service Service) []Vulnerability {
    var vulnerabilities []Vulnerability
    
    switch service.Name {
    case "HTTP":
        vulns := vs.scanWebService(target, service)
        vulnerabilities = append(vulnerabilities, vulns...)
    case "SSH":
        vulns := vs.scanSSHService(target, service)
        vulnerabilities = append(vulnerabilities, vulns...)
    case "FTP":
        vulns := vs.scanFTPService(target, service)
        vulnerabilities = append(vulnerabilities, vulns...)
    }
    
    return vulnerabilities
}

// 扫描Web服务
func (vs *VulnerabilityScanner) scanWebService(target Target, service Service) []Vulnerability {
    var vulnerabilities []Vulnerability
    
    // 检测常见Web漏洞
    if vuln := vs.detectSQLInjection(target, service); vuln != nil {
        vulnerabilities = append(vulnerabilities, *vuln)
    }
    
    if vuln := vs.detectXSS(target, service); vuln != nil {
        vulnerabilities = append(vulnerabilities, *vuln)
    }
    
    if vuln := vs.detectCSRF(target, service); vuln != nil {
        vulnerabilities = append(vulnerabilities, *vuln)
    }
    
    if vuln := vs.detectDirectoryTraversal(target, service); vuln != nil {
        vulnerabilities = append(vulnerabilities, *vuln)
    }
    
    return vulnerabilities
}

// 检测SQL注入
func (vs *VulnerabilityScanner) detectSQLInjection(target Target, service Service) *Vulnerability {
    // SQL注入检测逻辑
    payloads := []string{
        "' OR '1'='1",
        "'; DROP TABLE users; --",
        "' UNION SELECT * FROM users --",
    }
    
    for _, payload := range payloads {
        if vs.testSQLInjection(target, service, payload) {
            return &Vulnerability{
                ID:          generateID(),
                Title:       "SQL Injection Vulnerability",
                Description: "Potential SQL injection vulnerability detected",
                Severity:    SeverityHigh,
                Category:    "Injection",
                References:  []string{"https://owasp.org/www-community/attacks/SQL_Injection"},
                Evidence: map[string]interface{}{
                    "payload": payload,
                    "service": service.Name,
                    "port":    service.Port,
                },
                DiscoveredAt: time.Now(),
            }
        }
    }
    
    return nil
}

// 测试SQL注入
func (vs *VulnerabilityScanner) testSQLInjection(target Target, service Service, payload string) bool {
    // 构建测试URL
    testURL := fmt.Sprintf("http://%s:%d/test?param=%s", target.Host, service.Port, url.QueryEscape(payload))
    
    // 发送请求
    resp, err := http.Get(testURL)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    // 分析响应
    body, _ := io.ReadAll(resp.Body)
    response := string(body)
    
    // 检查SQL错误信息
    sqlErrors := []string{
        "SQL syntax",
        "mysql_fetch_array",
        "ORA-",
        "PostgreSQL",
        "SQLite",
    }
    
    for _, sqlError := range sqlErrors {
        if strings.Contains(response, sqlError) {
            return true
        }
    }
    
    return false
}
```

## 6. 端口扫描

### 6.1 高级端口扫描

```go
// 高级端口扫描器
type AdvancedPortScanner struct {
    baseScanner *PortScanner
    nmapScanner *NmapScanner
    synScanner  *SYNScanner
}

// SYN扫描器
type SYNScanner struct {
    timeout time.Duration
    workers int
}

// 创建SYN扫描器
func NewSYNScanner(timeout time.Duration, workers int) *SYNScanner {
    return &SYNScanner{
        timeout: timeout,
        workers: workers,
    }
}

// SYN扫描
func (ss *SYNScanner) Scan(target Target) []Service {
    var services []Service
    var wg sync.WaitGroup
    
    // 全端口扫描
    for port := 1; port <= 65535; port++ {
        wg.Add(1)
        go func(p int) {
            defer wg.Done()
            if ss.isPortOpen(target.Host, p) {
                service := Service{
                    Port:     p,
                    Protocol: "tcp",
                    Name:     "Unknown",
                }
                services = append(services, service)
            }
        }(port)
    }
    
    wg.Wait()
    return services
}

// 检查端口是否开放
func (ss *SYNScanner) isPortOpen(host string, port int) bool {
    // 使用原始套接字进行SYN扫描
    // 注意：这需要root权限
    return false // 简化实现
}
```

## 7. 服务识别

### 7.1 服务指纹识别

```go
// 服务指纹识别器
type ServiceFingerprinter struct {
    signatures map[string]ServiceSignature
    database   *ServiceDatabase
}

// 服务签名
type ServiceSignature struct {
    Service     string   `json:"service"`
    Version     string   `json:"version"`
    Patterns    []string `json:"patterns"`
    Ports       []int    `json:"ports"`
    Protocols   []string `json:"protocols"`
}

// 创建服务指纹识别器
func NewServiceFingerprinter() *ServiceFingerprinter {
    return &ServiceFingerprinter{
        signatures: make(map[string]ServiceSignature),
        database:   NewServiceDatabase(),
    }
}

// 识别服务
func (sf *ServiceFingerprinter) Fingerprint(service Service) Service {
    // 基于banner和端口进行服务识别
    if service.Banner != "" {
        service.Name = sf.identifyByBanner(service.Banner)
        service.Version = sf.extractVersion(service.Banner)
    } else {
        service.Name = sf.identifyByPort(service.Port)
    }
    
    return service
}

// 基于banner识别服务
func (sf *ServiceFingerprinter) identifyByBanner(banner string) string {
    // 常见服务banner模式
    patterns := map[string][]string{
        "Apache":     {"Apache", "httpd"},
        "Nginx":      {"nginx"},
        "IIS":        {"IIS", "Microsoft"},
        "SSH":        {"SSH", "OpenSSH"},
        "MySQL":      {"MySQL"},
        "PostgreSQL": {"PostgreSQL"},
        "Redis":      {"Redis"},
        "MongoDB":    {"MongoDB"},
    }
    
    for service, keywords := range patterns {
        for _, keyword := range keywords {
            if strings.Contains(strings.ToLower(banner), strings.ToLower(keyword)) {
                return service
            }
        }
    }
    
    return "Unknown"
}

// 提取版本信息
func (sf *ServiceFingerprinter) extractVersion(banner string) string {
    // 版本提取正则表达式
    versionPatterns := []string{
        `(\d+\.\d+\.\d+)`,
        `(\d+\.\d+)`,
        `version (\d+\.\d+\.\d+)`,
    }
    
    for _, pattern := range versionPatterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(banner)
        if len(matches) > 1 {
            return matches[1]
        }
    }
    
    return ""
}
```

## 8. 安全评估

### 8.1 风险评估引擎

```go
// 风险评估引擎
type RiskAssessmentEngine struct {
    riskModels  map[string]RiskModel
    calculators map[string]RiskCalculator
}

// 风险模型
type RiskModel struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Factors     []RiskFactor           `json:"factors"`
    Weights     map[string]float64     `json:"weights"`
    Thresholds  map[string]float64     `json:"thresholds"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// 风险因子
type RiskFactor struct {
    Name        string  `json:"name"`
    Weight      float64 `json:"weight"`
    MinValue    float64 `json:"min_value"`
    MaxValue    float64 `json:"max_value"`
    Description string  `json:"description"`
}

// 风险计算器接口
type RiskCalculator interface {
    Calculate(vulnerabilities []Vulnerability, target Target) RiskScore
}

// 风险评分
type RiskScore struct {
    Overall     float64            `json:"overall"`
    Components  map[string]float64 `json:"components"`
    Level       RiskLevel          `json:"level"`
    Details     []RiskDetail       `json:"details"`
    CalculatedAt time.Time         `json:"calculated_at"`
}

// 风险等级
type RiskLevel int

const (
    RiskLevelLow RiskLevel = iota
    RiskLevelMedium
    RiskLevelHigh
    RiskLevelCritical
)

// 风险详情
type RiskDetail struct {
    Factor      string  `json:"factor"`
    Value       float64 `json:"value"`
    Weight      float64 `json:"weight"`
    Contribution float64 `json:"contribution"`
    Description string  `json:"description"`
}

// 创建风险评估引擎
func NewRiskAssessmentEngine() *RiskAssessmentEngine {
    return &RiskAssessmentEngine{
        riskModels:  make(map[string]RiskModel),
        calculators: make(map[string]RiskCalculator),
    }
}

// 注册风险模型
func (rae *RiskAssessmentEngine) RegisterRiskModel(model RiskModel) {
    rae.riskModels[model.ID] = model
}

// 注册风险计算器
func (rae *RiskAssessmentEngine) RegisterCalculator(name string, calculator RiskCalculator) {
    rae.calculators[name] = calculator
}

// 评估风险
func (rae *RiskAssessmentEngine) AssessRisk(vulnerabilities []Vulnerability, target Target, modelID string) RiskScore {
    model, exists := rae.riskModels[modelID]
    if !exists {
        // 使用默认模型
        model = rae.getDefaultModel()
    }
    
    // 计算各风险因子
    components := make(map[string]float64)
    var details []RiskDetail
    
    for _, factor := range model.Factors {
        value := rae.calculateFactorValue(factor, vulnerabilities, target)
        weight := model.Weights[factor.Name]
        
        components[factor.Name] = value
        details = append(details, RiskDetail{
            Factor:       factor.Name,
            Value:        value,
            Weight:       weight,
            Contribution: value * weight,
            Description:  factor.Description,
        })
    }
    
    // 计算总体风险评分
    overall := rae.calculateOverallRisk(components, model.Weights)
    level := rae.determineRiskLevel(overall, model.Thresholds)
    
    return RiskScore{
        Overall:      overall,
        Components:   components,
        Level:        level,
        Details:      details,
        CalculatedAt: time.Now(),
    }
}

// 计算风险因子值
func (rae *RiskAssessmentEngine) calculateFactorValue(factor RiskFactor, vulnerabilities []Vulnerability, target Target) float64 {
    switch factor.Name {
    case "vulnerability_count":
        return float64(len(vulnerabilities))
    case "severity_score":
        return rae.calculateSeverityScore(vulnerabilities)
    case "exposure_time":
        return rae.calculateExposureTime(target)
    case "attack_surface":
        return rae.calculateAttackSurface(target)
    default:
        return 0.0
    }
}

// 计算严重程度评分
func (rae *RiskAssessmentEngine) calculateSeverityScore(vulnerabilities []Vulnerability) float64 {
    if len(vulnerabilities) == 0 {
        return 0.0
    }
    
    var totalScore float64
    for _, vuln := range vulnerabilities {
        switch vuln.Severity {
        case SeverityCritical:
            totalScore += 10.0
        case SeverityHigh:
            totalScore += 7.5
        case SeverityMedium:
            totalScore += 5.0
        case SeverityLow:
            totalScore += 2.5
        case SeverityInfo:
            totalScore += 1.0
        }
    }
    
    return totalScore / float64(len(vulnerabilities))
}

// 计算总体风险
func (rae *RiskAssessmentEngine) calculateOverallRisk(components map[string]float64, weights map[string]float64) float64 {
    var totalScore float64
    var totalWeight float64
    
    for factor, value := range components {
        weight := weights[factor]
        totalScore += value * weight
        totalWeight += weight
    }
    
    if totalWeight > 0 {
        return totalScore / totalWeight
    }
    
    return 0.0
}

// 确定风险等级
func (rae *RiskAssessmentEngine) determineRiskLevel(score float64, thresholds map[string]float64) RiskLevel {
    if score >= thresholds["critical"] {
        return RiskLevelCritical
    } else if score >= thresholds["high"] {
        return RiskLevelHigh
    } else if score >= thresholds["medium"] {
        return RiskLevelMedium
    } else {
        return RiskLevelLow
    }
}
```

## 9. 报告生成

### 9.1 报告生成器

```go
// 报告生成器
type ReportGenerator struct {
    templates map[string]ReportTemplate
    formatters map[string]ReportFormatter
}

// 报告模板
type ReportTemplate struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        ReportType             `json:"type"`
    Template    string                 `json:"template"`
    Variables   []string               `json:"variables"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// 报告类型
type ReportType string

const (
    ReportTypeHTML ReportType = "html"
    ReportTypePDF  ReportType = "pdf"
    ReportTypeJSON ReportType = "json"
    ReportTypeXML  ReportType = "xml"
    ReportTypeCSV  ReportType = "csv"
)

// 报告格式化器接口
type ReportFormatter interface {
    Format(report ScanReport) ([]byte, error)
    GetType() ReportType
}

// 扫描报告
type ScanReport struct {
    ID              string         `json:"id"`
    Title           string         `json:"title"`
    Description     string         `json:"description"`
    Target          Target         `json:"target"`
    ScanConfig      *ScanConfig    `json:"scan_config"`
    Results         []ScanResult   `json:"results"`
    RiskAssessment  RiskScore      `json:"risk_assessment"`
    Summary         ReportSummary  `json:"summary"`
    GeneratedAt     time.Time      `json:"generated_at"`
    GeneratedBy     string         `json:"generated_by"`
}

// 报告摘要
type ReportSummary struct {
    TotalTargets       int `json:"total_targets"`
    ScannedTargets     int `json:"scanned_targets"`
    TotalVulnerabilities int `json:"total_vulnerabilities"`
    CriticalVulns      int `json:"critical_vulnerabilities"`
    HighVulns          int `json:"high_vulnerabilities"`
    MediumVulns        int `json:"medium_vulnerabilities"`
    LowVulns           int `json:"low_vulnerabilities"`
    InfoVulns          int `json:"info_vulnerabilities"`
    TotalServices      int `json:"total_services"`
    ScanDuration       time.Duration `json:"scan_duration"`
}

// 创建报告生成器
func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{
        templates:  make(map[string]ReportTemplate),
        formatters: make(map[string]ReportFormatter),
    }
}

// 注册报告模板
func (rg *ReportGenerator) RegisterTemplate(template ReportTemplate) {
    rg.templates[template.ID] = template
}

// 注册报告格式化器
func (rg *ReportGenerator) RegisterFormatter(formatter ReportFormatter) {
    rg.formatters[formatter.GetType()] = formatter
}

// 生成报告
func (rg *ReportGenerator) GenerateReport(results []ScanResult, target Target, config *ScanConfig, reportType ReportType) ([]byte, error) {
    // 创建扫描报告
    report := rg.createScanReport(results, target, config)
    
    // 获取对应的格式化器
    formatter, exists := rg.formatters[string(reportType)]
    if !exists {
        return nil, fmt.Errorf("unsupported report type: %s", reportType)
    }
    
    // 格式化报告
    return formatter.Format(report)
}

// 创建扫描报告
func (rg *ReportGenerator) createScanReport(results []ScanResult, target Target, config *ScanConfig) ScanReport {
    // 统计摘要信息
    summary := rg.calculateSummary(results)
    
    // 风险评估
    var allVulns []Vulnerability
    for _, result := range results {
        allVulns = append(allVulns, result.Vulnerabilities...)
    }
    
    riskEngine := NewRiskAssessmentEngine()
    riskAssessment := riskEngine.AssessRisk(allVulns, target, "default")
    
    return ScanReport{
        ID:             generateID(),
        Title:          fmt.Sprintf("Security Scan Report - %s", target.Host),
        Description:    "Comprehensive security assessment report",
        Target:         target,
        ScanConfig:     config,
        Results:        results,
        RiskAssessment: riskAssessment,
        Summary:        summary,
        GeneratedAt:    time.Now(),
        GeneratedBy:    "Security Scanner v1.0",
    }
}

// 计算报告摘要
func (rg *ReportGenerator) calculateSummary(results []ScanResult) ReportSummary {
    summary := ReportSummary{}
    
    for _, result := range results {
        summary.TotalTargets++
        summary.ScannedTargets++
        summary.TotalVulnerabilities += len(result.Vulnerabilities)
        summary.TotalServices += len(result.Services)
        summary.ScanDuration += result.ScanTime
        
        // 统计各严重程度的漏洞数量
        for _, vuln := range result.Vulnerabilities {
            switch vuln.Severity {
            case SeverityCritical:
                summary.CriticalVulns++
            case SeverityHigh:
                summary.HighVulns++
            case SeverityMedium:
                summary.MediumVulns++
            case SeverityLow:
                summary.LowVulns++
            case SeverityInfo:
                summary.InfoVulns++
            }
        }
    }
    
    return summary
}
```

## 10. 性能优化

### 10.1 并发优化

```go
// 并发扫描管理器
type ConcurrentScanManager struct {
    workerPool *WorkerPool
    rateLimiter *RateLimiter
    resultCollector *ResultCollector
}

// 工作池
type WorkerPool struct {
    workers    int
    tasks      chan ScanTask
    results    chan ScanResult
    errors     chan error
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

// 扫描任务
type ScanTask struct {
    ID       string `json:"id"`
    Target   Target `json:"target"`
    Config   *ScanConfig `json:"config"`
    Priority int    `json:"priority"`
}

// 创建并发扫描管理器
func NewConcurrentScanManager(workers int) *ConcurrentScanManager {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ConcurrentScanManager{
        workerPool: &WorkerPool{
            workers: workers,
            tasks:   make(chan ScanTask, 1000),
            results: make(chan ScanResult, 1000),
            errors:  make(chan error, 1000),
            ctx:     ctx,
            cancel:  cancel,
        },
        rateLimiter:    NewRateLimiter(100), // 100 requests per second
        resultCollector: NewResultCollector(),
    }
}

// 启动工作池
func (csm *ConcurrentScanManager) Start() {
    for i := 0; i < csm.workerPool.workers; i++ {
        csm.workerPool.wg.Add(1)
        go csm.worker(i)
    }
}

// 工作协程
func (csm *ConcurrentScanManager) worker(id int) {
    defer csm.workerPool.wg.Done()
    
    for {
        select {
        case task := <-csm.workerPool.tasks:
            // 执行扫描任务
            result := csm.executeTask(task)
            csm.workerPool.results <- result
        case <-csm.workerPool.ctx.Done():
            return
        }
    }
}

// 执行扫描任务
func (csm *ConcurrentScanManager) executeTask(task ScanTask) ScanResult {
    // 应用速率限制
    csm.rateLimiter.Wait()
    
    // 创建扫描器
    scanner := NewBaseScanner(task.Config)
    
    // 执行扫描
    start := time.Now()
    vulnerabilities := scanner.ScanVulnerabilities(task.Target, []Service{})
    
    return ScanResult{
        Target:         task.Target,
        Vulnerabilities: vulnerabilities,
        ScanTime:       time.Since(start),
        CompletedAt:    time.Now(),
    }
}

// 速率限制器
type RateLimiter struct {
    rate       int
    tokens     chan struct{}
    ticker     *time.Ticker
    ctx        context.Context
    cancel     context.CancelFunc
}

// 创建速率限制器
func NewRateLimiter(rate int) *RateLimiter {
    ctx, cancel := context.WithCancel(context.Background())
    
    rl := &RateLimiter{
        rate:   rate,
        tokens: make(chan struct{}, rate),
        ctx:    ctx,
        cancel: cancel,
    }
    
    // 初始化令牌
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    // 启动令牌生成器
    rl.ticker = time.NewTicker(time.Second / time.Duration(rate))
    go rl.tokenGenerator()
    
    return rl
}

// 令牌生成器
func (rl *RateLimiter) tokenGenerator() {
    for {
        select {
        case <-rl.ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default:
                // 令牌桶已满
            }
        case <-rl.ctx.Done():
            return
        }
    }
}

// 等待令牌
func (rl *RateLimiter) Wait() {
    select {
    case <-rl.tokens:
        // 获取到令牌，可以继续
    case <-rl.ctx.Done():
        // 上下文取消
    }
}
```

## 11. 部署运维

### 11.1 容器化部署

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scanner .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/scanner .

# 复制配置文件
COPY configs/ ./configs/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./scanner"]
```

### 11.2 Kubernetes部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: security-scanner
  labels:
    app: security-scanner
spec:
  replicas: 3
  selector:
    matchLabels:
      app: security-scanner
  template:
    metadata:
      labels:
        app: security-scanner
    spec:
      containers:
      - name: scanner
        image: security-scanner:latest
        ports:
        - containerPort: 8080
        env:
        - name: SCAN_CONFIG_PATH
          value: "/root/configs/scanner.yaml"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        volumeMounts:
        - name: config-volume
          mountPath: /root/configs
      volumes:
      - name: config-volume
        configMap:
          name: scanner-config

---
apiVersion: v1
kind: Service
metadata:
  name: security-scanner-service
spec:
  selector:
    app: security-scanner
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

### 11.3 监控配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'security-scanner'
    static_configs:
      - targets: ['security-scanner-service:80']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

这个安全扫描工具的实现提供了：

1. **完整的扫描引擎架构**
2. **多种扫描类型支持**
3. **并发和性能优化**
4. **风险评估模型**
5. **报告生成功能**
6. **容器化部署方案**

所有代码都使用Go语言实现，符合现代软件工程最佳实践。

---

**相关链接**：

- [02-入侵检测系统](../02-Intrusion-Detection-System/README.md)
- [03-加密服务](../03-Encryption-Services/README.md)
- [04-身份认证](../04-Identity-Authentication/README.md)
