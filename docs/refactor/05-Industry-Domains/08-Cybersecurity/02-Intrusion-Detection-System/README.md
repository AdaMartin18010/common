# 02-入侵检测系统 (Intrusion Detection System)

## 1. 概述

### 1.1 定义与目标

入侵检测系统(IDS)是网络安全监控的核心组件，用于实时检测和响应网络中的恶意活动和安全威胁。

**形式化定义**：
设 $E$ 为事件流，$P$ 为模式集合，$A$ 为告警集合，则入侵检测函数 $f$ 定义为：

$$f: E \times P \rightarrow A$$

其中：

- $E = \{e_1, e_2, ..., e_n\}$ 为网络事件流
- $P = \{p_1, p_2, ..., p_m\}$ 为检测模式集合
- $A = \{a_1, a_2, ..., a_k\}$ 为告警集合

### 1.2 检测原理

**模式匹配**：
$$\text{Match}(e, p) = \begin{cases}
1 & \text{if } e \text{ matches pattern } p \\
0 & \text{otherwise}
\end{cases}$$

**异常检测**：
$$\text{Anomaly}(e) = \begin{cases}
1 & \text{if } \|e - \mu\| > k\sigma \\
0 & \text{otherwise}
\end{cases}$$

其中 $\mu$ 为正常行为均值，$\sigma$ 为标准差，$k$ 为阈值系数。

## 2. 架构设计

### 2.1 系统架构

```
┌─────────────────────────────────────┐
│           事件采集层                  │
├─────────────────────────────────────┤
│           预处理层                    │
├─────────────────────────────────────┤
│           检测引擎层                  │
├─────────────────────────────────────┤
│           告警处理层                  │
├─────────────────────────────────────┤
│           响应执行层                  │
└─────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 事件采集器

```go
// EventCollector 事件采集器
type EventCollector struct {
    sources    map[string]EventSource
    buffer     chan *SecurityEvent
    config     *CollectorConfig
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// CollectorConfig 采集器配置
type CollectorConfig struct {
    BufferSize    int           `json:"buffer_size"`
    BatchSize     int           `json:"batch_size"`
    FlushInterval time.Duration `json:"flush_interval"`
    MaxWorkers    int           `json:"max_workers"`
}

// EventSource 事件源接口
type EventSource interface {
    ID() string
    Type() SourceType
    Start(ctx context.Context) error
    Stop() error
    Events() <-chan *SecurityEvent
}

// SourceType 事件源类型
type SourceType string

const (
    SourceTypeNetwork    SourceType = "network"
    SourceTypeSystem     SourceType = "system"
    SourceTypeApplication SourceType = "application"
    SourceTypeLog        SourceType = "log"
)

// SecurityEvent 安全事件
type SecurityEvent struct {
    ID          string                 `json:"id"`
    Source      string                 `json:"source"`
    Type        EventType              `json:"type"`
    Timestamp   time.Time              `json:"timestamp"`
    SourceIP    string                 `json:"source_ip"`
    DestIP      string                 `json:"dest_ip"`
    SourcePort  int                    `json:"source_port"`
    DestPort    int                    `json:"dest_port"`
    Protocol    string                 `json:"protocol"`
    Payload     []byte                 `json:"payload"`
    Metadata    map[string]interface{} `json:"metadata"`
    Severity    EventSeverity          `json:"severity"`
}

// EventType 事件类型
type EventType string

const (
    EventTypeConnection EventType = "connection"
    EventTypePacket     EventType = "packet"
    EventTypeLogin      EventType = "login"
    EventTypeFileAccess EventType = "file_access"
    EventTypeProcess    EventType = "process"
)

// EventSeverity 事件严重程度
type EventSeverity string

const (
    EventSeverityLow    EventSeverity = "low"
    EventSeverityMedium EventSeverity = "medium"
    EventSeverityHigh   EventSeverity = "high"
    EventSeverityCritical EventSeverity = "critical"
)

// Start 启动事件采集
func (ec *EventCollector) Start() error {
    ec.ctx, ec.cancel = context.WithCancel(context.Background())

    // 启动所有事件源
    for _, source := range ec.sources {
        ec.wg.Add(1)
        go func(s EventSource) {
            defer ec.wg.Done()
            if err := s.Start(ec.ctx); err != nil {
                log.Printf("Failed to start source %s: %v", s.ID(), err)
            }
        }(source)
    }

    // 启动事件处理协程
    for i := 0; i < ec.config.MaxWorkers; i++ {
        ec.wg.Add(1)
        go ec.eventProcessor()
    }

    return nil
}

// Stop 停止事件采集
func (ec *EventCollector) Stop() error {
    ec.cancel()
    ec.wg.Wait()
    return nil
}

// eventProcessor 事件处理器
func (ec *EventCollector) eventProcessor() {
    batch := make([]*SecurityEvent, 0, ec.config.BatchSize)
    ticker := time.NewTicker(ec.config.FlushInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ec.ctx.Done():
            ec.flushBatch(batch)
            return
        case event := <-ec.buffer:
            batch = append(batch, event)
            if len(batch) >= ec.config.BatchSize {
                ec.flushBatch(batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                ec.flushBatch(batch)
                batch = batch[:0]
            }
        }
    }
}

// flushBatch 刷新批次
func (ec *EventCollector) flushBatch(batch []*SecurityEvent) {
    if len(batch) == 0 {
        return
    }

    // 发送到检测引擎
    // 这里简化实现，实际应该发送到检测引擎
    log.Printf("Flushed %d events", len(batch))
}
```

#### 2.2.2 网络事件源

```go
// NetworkEventSource 网络事件源
type NetworkEventSource struct {
    id       string
    interface string
    filter   string
    config   *NetworkSourceConfig
    events   chan *SecurityEvent
    ctx      context.Context
    cancel   context.CancelFunc
}

// NetworkSourceConfig 网络源配置
type NetworkSourceConfig struct {
    Interface     string        `json:"interface"`
    Filter        string        `json:"filter"`
    SnapLen       int           `json:"snap_len"`
    Promiscuous   bool          `json:"promiscuous"`
    Timeout       time.Duration `json:"timeout"`
    MaxPacketSize int           `json:"max_packet_size"`
}

// NewNetworkEventSource 创建网络事件源
func NewNetworkEventSource(id string, config *NetworkSourceConfig) *NetworkEventSource {
    return &NetworkEventSource{
        id:     id,
        config: config,
        events: make(chan *SecurityEvent, 1000),
    }
}

// Start 启动网络事件源
func (ns *NetworkEventSource) Start(ctx context.Context) error {
    ns.ctx, ns.cancel = context.WithCancel(ctx)

    // 打开网络接口
    handle, err := pcap.OpenLive(ns.config.Interface, int32(ns.config.SnapLen), ns.config.Promiscuous, ns.config.Timeout)
    if err != nil {
        return fmt.Errorf("failed to open interface: %w", err)
    }
    defer handle.Close()

    // 设置过滤器
    if ns.config.Filter != "" {
        if err := handle.SetBPFFilter(ns.config.Filter); err != nil {
            return fmt.Errorf("failed to set BPF filter: %w", err)
        }
    }

    // 开始捕获
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

    go func() {
        for {
            select {
            case <-ns.ctx.Done():
                return
            case packet := <-packetSource.Packets():
                ns.processPacket(packet)
            }
        }
    }()

    return nil
}

// Stop 停止网络事件源
func (ns *NetworkEventSource) Stop() error {
    if ns.cancel != nil {
        ns.cancel()
    }
    return nil
}

// Events 获取事件通道
func (ns *NetworkEventSource) Events() <-chan *SecurityEvent {
    return ns.events
}

// ID 获取源ID
func (ns *NetworkEventSource) ID() string {
    return ns.id
}

// Type 获取源类型
func (ns *NetworkEventSource) Type() SourceType {
    return SourceTypeNetwork
}

// processPacket 处理网络包
func (ns *NetworkEventSource) processPacket(packet gopacket.Packet) {
    // 解析网络层
    networkLayer := packet.NetworkLayer()
    if networkLayer == nil {
        return
    }

    // 解析传输层
    transportLayer := packet.TransportLayer()
    if transportLayer == nil {
        return
    }

    var sourceIP, destIP string
    var sourcePort, destPort int
    var protocol string

    // 提取IP信息
    switch v := networkLayer.(type) {
    case *layers.IPv4:
        sourceIP = v.SrcIP.String()
        destIP = v.DstIP.String()
        protocol = v.Protocol.String()
    case *layers.IPv6:
        sourceIP = v.SrcIP.String()
        destIP = v.DstIP.String()
        protocol = v.NextHeader.String()
    default:
        return
    }

    // 提取端口信息
    switch v := transportLayer.(type) {
    case *layers.TCP:
        sourcePort = int(v.SrcPort)
        destPort = int(v.DstPort)
    case *layers.UDP:
        sourcePort = int(v.SrcPort)
        destPort = int(v.DstPort)
    default:
        return
    }

    // 创建安全事件
    event := &SecurityEvent{
        ID:         uuid.New().String(),
        Source:     ns.id,
        Type:       EventTypePacket,
        Timestamp:  packet.Metadata().Timestamp,
        SourceIP:   sourceIP,
        DestIP:     destIP,
        SourcePort: sourcePort,
        DestPort:   destPort,
        Protocol:   protocol,
        Payload:    packet.Data(),
        Metadata:   make(map[string]interface{}),
        Severity:   ns.calculateSeverity(packet),
    }

    // 添加元数据
    event.Metadata["packet_size"] = len(packet.Data())
    event.Metadata["layer_count"] = len(packet.Layers())

    // 发送事件
    select {
    case ns.events <- event:
    default:
        // 通道满，丢弃事件
        log.Printf("Event channel full, dropping event")
    }
}

// calculateSeverity 计算事件严重程度
func (ns *NetworkEventSource) calculateSeverity(packet gopacket.Packet) EventSeverity {
    // 检查常见攻击模式
    payload := packet.Data()
    payloadStr := strings.ToLower(string(payload))

    // SQL注入检测
    sqlPatterns := []string{"union select", "drop table", "insert into", "delete from"}
    for _, pattern := range sqlPatterns {
        if strings.Contains(payloadStr, pattern) {
            return EventSeverityHigh
        }
    }

    // XSS检测
    xssPatterns := []string{"<script>", "javascript:", "onload=", "onerror="}
    for _, pattern := range xssPatterns {
        if strings.Contains(payloadStr, pattern) {
            return EventSeverityMedium
        }
    }

    // 端口扫描检测
    if ns.isPortScan(packet) {
        return EventSeverityMedium
    }

    return EventSeverityLow
}

// isPortScan 检测端口扫描
func (ns *NetworkEventSource) isPortScan(packet gopacket.Packet) bool {
    // 简化实现，实际需要更复杂的检测逻辑
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer == nil {
        return false
    }

    tcp, _ := tcpLayer.(*layers.TCP)
    if tcp == nil {
        return false
    }

    // 检查SYN标志
    return tcp.SYN && !tcp.ACK
}
```

#### 2.2.3 检测引擎

```go
// DetectionEngine 检测引擎
type DetectionEngine struct {
    rules       []DetectionRule
    patterns    map[string]*Pattern
    anomalies   *AnomalyDetector
    alerts      chan *SecurityAlert
    config      *EngineConfig
    mu          sync.RWMutex
}

// EngineConfig 引擎配置
type EngineConfig struct {
    MaxRules        int           `json:"max_rules"`
    AlertThreshold  int           `json:"alert_threshold"`
    TimeWindow      time.Duration `json:"time_window"`
    EnableAnomaly   bool          `json:"enable_anomaly"`
    EnableSignature bool          `json:"enable_signature"`
}

// DetectionRule 检测规则
type DetectionRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Type        RuleType               `json:"type"`
    Pattern     string                 `json:"pattern"`
    Conditions  []Condition            `json:"conditions"`
    Actions     []Action               `json:"actions"`
    Enabled     bool                   `json:"enabled"`
    Priority    int                    `json:"priority"`
    CreatedAt   time.Time              `json:"created_at"`
}

// RuleType 规则类型
type RuleType string

const (
    RuleTypeSignature RuleType = "signature"
    RuleTypeAnomaly   RuleType = "anomaly"
    RuleTypeBehavior  RuleType = "behavior"
)

// Condition 条件
type Condition struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}

// Action 动作
type Action struct {
    Type    string                 `json:"type"`
    Params  map[string]interface{} `json:"params"`
}

// SecurityAlert 安全告警
type SecurityAlert struct {
    ID          string                 `json:"id"`
    RuleID      string                 `json:"rule_id"`
    RuleName    string                 `json:"rule_name"`
    Severity    AlertSeverity          `json:"severity"`
    Source      string                 `json:"source"`
    Target      string                 `json:"target"`
    Description string                 `json:"description"`
    Evidence    []*SecurityEvent       `json:"evidence"`
    Timestamp   time.Time              `json:"timestamp"`
    Status      AlertStatus            `json:"status"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// AlertSeverity 告警严重程度
type AlertSeverity string

const (
    AlertSeverityLow    AlertSeverity = "low"
    AlertSeverityMedium AlertSeverity = "medium"
    AlertSeverityHigh   AlertSeverity = "high"
    AlertSeverityCritical AlertSeverity = "critical"
)

// AlertStatus 告警状态
type AlertStatus string

const (
    AlertStatusNew      AlertStatus = "new"
    AlertStatusAcknowledged AlertStatus = "acknowledged"
    AlertStatusResolved AlertStatus = "resolved"
    AlertStatusFalsePositive AlertStatus = "false_positive"
)

// ProcessEvent 处理安全事件
func (de *DetectionEngine) ProcessEvent(event *SecurityEvent) error {
    de.mu.RLock()
    defer de.mu.RUnlock()

    // 签名检测
    if de.config.EnableSignature {
        if alert := de.signatureDetection(event); alert != nil {
            select {
            case de.alerts <- alert:
            default:
                log.Printf("Alert channel full, dropping alert")
            }
        }
    }

    // 异常检测
    if de.config.EnableAnomaly {
        if alert := de.anomalyDetection(event); alert != nil {
            select {
            case de.alerts <- alert:
            default:
                log.Printf("Alert channel full, dropping alert")
            }
        }
    }

    return nil
}

// signatureDetection 签名检测
func (de *DetectionEngine) signatureDetection(event *SecurityEvent) *SecurityAlert {
    for _, rule := range de.rules {
        if !rule.Enabled || rule.Type != RuleTypeSignature {
            continue
        }

        if de.matchRule(event, rule) {
            return de.createAlert(event, rule)
        }
    }

    return nil
}

// matchRule 匹配规则
func (de *DetectionEngine) matchRule(event *SecurityEvent, rule DetectionRule) bool {
    for _, condition := range rule.Conditions {
        if !de.evaluateCondition(event, condition) {
            return false
        }
    }

    return true
}

// evaluateCondition 评估条件
func (de *DetectionEngine) evaluateCondition(event *SecurityEvent, condition Condition) bool {
    var value interface{}

    // 获取字段值
    switch condition.Field {
    case "source_ip":
        value = event.SourceIP
    case "dest_ip":
        value = event.DestIP
    case "source_port":
        value = event.SourcePort
    case "dest_port":
        value = event.DestPort
    case "protocol":
        value = event.Protocol
    case "severity":
        value = event.Severity
    default:
        return false
    }

    // 执行操作符
    switch condition.Operator {
    case "equals":
        return value == condition.Value
    case "not_equals":
        return value != condition.Value
    case "contains":
        if str, ok := value.(string); ok {
            if target, ok := condition.Value.(string); ok {
                return strings.Contains(str, target)
            }
        }
        return false
    case "greater_than":
        if num, ok := value.(int); ok {
            if target, ok := condition.Value.(int); ok {
                return num > target
            }
        }
        return false
    case "less_than":
        if num, ok := value.(int); ok {
            if target, ok := condition.Value.(int); ok {
                return num < target
            }
        }
        return false
    default:
        return false
    }
}

// createAlert 创建告警
func (de *DetectionEngine) createAlert(event *SecurityEvent, rule DetectionRule) *SecurityAlert {
    return &SecurityAlert{
        ID:          uuid.New().String(),
        RuleID:      rule.ID,
        RuleName:    rule.Name,
        Severity:    de.mapSeverity(event.Severity),
        Source:      event.SourceIP,
        Target:      event.DestIP,
        Description: rule.Description,
        Evidence:    []*SecurityEvent{event},
        Timestamp:   time.Now(),
        Status:      AlertStatusNew,
        Metadata:    make(map[string]interface{}),
    }
}

// mapSeverity 映射严重程度
func (de *DetectionEngine) mapSeverity(eventSeverity EventSeverity) AlertSeverity {
    switch eventSeverity {
    case EventSeverityCritical:
        return AlertSeverityCritical
    case EventSeverityHigh:
        return AlertSeverityHigh
    case EventSeverityMedium:
        return AlertSeverityMedium
    case EventSeverityLow:
        return AlertSeverityLow
    default:
        return AlertSeverityLow
    }
}

// anomalyDetection 异常检测
func (de *DetectionEngine) anomalyDetection(event *SecurityEvent) *SecurityAlert {
    if de.anomalies == nil {
        return nil
    }

    if isAnomaly := de.anomalies.Detect(event); isAnomaly {
        return &SecurityAlert{
            ID:          uuid.New().String(),
            RuleID:      "anomaly",
            RuleName:    "Anomaly Detection",
            Severity:    AlertSeverityMedium,
            Source:      event.SourceIP,
            Target:      event.DestIP,
            Description: "Anomalous behavior detected",
            Evidence:    []*SecurityEvent{event},
            Timestamp:   time.Now(),
            Status:      AlertStatusNew,
            Metadata:    make(map[string]interface{}),
        }
    }

    return nil
}
```

#### 2.2.4 异常检测器

```go
// AnomalyDetector 异常检测器
type AnomalyDetector struct {
    models      map[string]*AnomalyModel
    baseline    *BehaviorBaseline
    config      *AnomalyConfig
    mu          sync.RWMutex
}

// AnomalyConfig 异常检测配置
type AnomalyConfig struct {
    LearningRate    float64       `json:"learning_rate"`
    Threshold       float64       `json:"threshold"`
    WindowSize      int           `json:"window_size"`
    UpdateInterval  time.Duration `json:"update_interval"`
}

// AnomalyModel 异常模型
type AnomalyModel struct {
    ID       string
    Type     string
    Features []string
    Weights  []float64
    Bias     float64
    mu       sync.RWMutex
}

// BehaviorBaseline 行为基线
type BehaviorBaseline struct {
    patterns map[string]*PatternStats
    mu       sync.RWMutex
}

// PatternStats 模式统计
type PatternStats struct {
    Count     int64
    Mean      float64
    Variance  float64
    LastSeen  time.Time
}

// Detect 检测异常
func (ad *AnomalyDetector) Detect(event *SecurityEvent) bool {
    ad.mu.RLock()
    defer ad.mu.RUnlock()

    // 提取特征
    features := ad.extractFeatures(event)

    // 计算异常分数
    score := ad.calculateAnomalyScore(features)

    // 检查是否超过阈值
    return score > ad.config.Threshold
}

// extractFeatures 提取特征
func (ad *AnomalyDetector) extractFeatures(event *SecurityEvent) []float64 {
    features := make([]float64, 0)

    // 时间特征
    features = append(features, float64(event.Timestamp.Hour()))
    features = append(features, float64(event.Timestamp.Weekday()))

    // 网络特征
    features = append(features, float64(event.SourcePort))
    features = append(features, float64(event.DestPort))

    // 协议特征
    protocolScore := ad.getProtocolScore(event.Protocol)
    features = append(features, protocolScore)

    // 负载特征
    payloadSize := len(event.Payload)
    features = append(features, float64(payloadSize))

    return features
}

// getProtocolScore 获取协议分数
func (ad *AnomalyDetector) getProtocolScore(protocol string) float64 {
    protocolScores := map[string]float64{
        "TCP": 1.0,
        "UDP": 0.8,
        "ICMP": 0.6,
        "HTTP": 0.9,
        "HTTPS": 0.7,
    }

    if score, exists := protocolScores[protocol]; exists {
        return score
    }

    return 0.5
}

// calculateAnomalyScore 计算异常分数
func (ad *AnomalyDetector) calculateAnomalyScore(features []float64) float64 {
    // 简化实现：使用欧几里得距离
    baseline := ad.getBaselineFeatures()

    if len(baseline) != len(features) {
        return 0.0
    }

    var sum float64
    for i := 0; i < len(features); i++ {
        diff := features[i] - baseline[i]
        sum += diff * diff
    }

    return math.Sqrt(sum)
}

// getBaselineFeatures 获取基线特征
func (ad *AnomalyDetector) getBaselineFeatures() []float64 {
    // 简化实现：返回默认基线
    return []float64{12.0, 3.0, 1024.0, 80.0, 0.8, 512.0}
}

// UpdateBaseline 更新基线
func (ad *AnomalyDetector) UpdateBaseline(event *SecurityEvent) {
    ad.mu.Lock()
    defer ad.mu.Unlock()

    // 更新模式统计
    pattern := ad.getEventPattern(event)

    if stats, exists := ad.baseline.patterns[pattern]; exists {
        stats.Count++
        stats.LastSeen = event.Timestamp

        // 更新均值和方差
        oldMean := stats.Mean
        stats.Mean = (stats.Mean*float64(stats.Count-1) + float64(len(event.Payload))) / float64(stats.Count)
        stats.Variance = (stats.Variance*float64(stats.Count-2) + float64(len(event.Payload)-oldMean)*float64(len(event.Payload)-stats.Mean)) / float64(stats.Count-1)
    } else {
        ad.baseline.patterns[pattern] = &PatternStats{
            Count:    1,
            Mean:     float64(len(event.Payload)),
            Variance: 0.0,
            LastSeen: event.Timestamp,
        }
    }
}

// getEventPattern 获取事件模式
func (ad *AnomalyDetector) getEventPattern(event *SecurityEvent) string {
    return fmt.Sprintf("%s:%d->%s:%d", event.SourceIP, event.SourcePort, event.DestIP, event.DestPort)
}
```

## 3. 数学建模

### 3.1 异常检测模型

**马氏距离**：
$$D^2 = (x - \mu)^T \Sigma^{-1} (x - \mu)$$

其中：
- $x$ 为特征向量
- $\mu$ 为均值向量
- $\Sigma$ 为协方差矩阵

### 3.2 检测率模型

**检测率计算**：
$$DR = \frac{TP}{TP + FN}$$

**误报率计算**：
$$FPR = \frac{FP}{FP + TN}$$

### 3.3 性能指标

**精确度**：
$$Precision = \frac{TP}{TP + FP}$$

**F1分数**：
$$F1 = 2 \times \frac{Precision \times Recall}{Precision + Recall}$$

## 4. 响应机制

### 4.1 自动响应

```go
// ResponseEngine 响应引擎
type ResponseEngine struct {
    actions    map[string]ResponseAction
    policies   []ResponsePolicy
    config     *ResponseConfig
    alerts     <-chan *SecurityAlert
    ctx        context.Context
    cancel     context.CancelFunc
}

// ResponseConfig 响应配置
type ResponseConfig struct {
    EnableAutoResponse bool          `json:"enable_auto_response"`
    ResponseDelay      time.Duration `json:"response_delay"`
    MaxConcurrent      int           `json:"max_concurrent"`
    RetryAttempts      int           `json:"retry_attempts"`
}

// ResponseAction 响应动作
type ResponseAction interface {
    ID() string
    Name() string
    Execute(ctx context.Context, alert *SecurityAlert) error
    Validate(alert *SecurityAlert) error
}

// ResponsePolicy 响应策略
type ResponsePolicy struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Conditions  []Condition            `json:"conditions"`
    Actions     []string               `json:"actions"`
    Priority    int                    `json:"priority"`
    Enabled     bool                   `json:"enabled"`
}

// Start 启动响应引擎
func (re *ResponseEngine) Start() error {
    re.ctx, re.cancel = context.WithCancel(context.Background())

    // 启动响应处理协程
    for i := 0; i < re.config.MaxConcurrent; i++ {
        go re.responseProcessor()
    }

    return nil
}

// Stop 停止响应引擎
func (re *ResponseEngine) Stop() error {
    if re.cancel != nil {
        re.cancel()
    }
    return nil
}

// responseProcessor 响应处理器
func (re *ResponseEngine) responseProcessor() {
    for {
        select {
        case <-re.ctx.Done():
            return
        case alert := <-re.alerts:
            re.processAlert(alert)
        }
    }
}

// processAlert 处理告警
func (re *ResponseEngine) processAlert(alert *SecurityAlert) {
    // 查找匹配的策略
    policies := re.findMatchingPolicies(alert)

    // 按优先级排序
    sort.Slice(policies, func(i, j int) bool {
        return policies[i].Priority > policies[j].Priority
    })

    // 执行响应动作
    for _, policy := range policies {
        for _, actionID := range policy.Actions {
            if action, exists := re.actions[actionID]; exists {
                if err := re.executeAction(action, alert); err != nil {
                    log.Printf("Failed to execute action %s: %v", actionID, err)
                }
            }
        }
    }
}

// findMatchingPolicies 查找匹配的策略
func (re *ResponseEngine) findMatchingPolicies(alert *SecurityAlert) []ResponsePolicy {
    var matching []ResponsePolicy

    for _, policy := range re.policies {
        if !policy.Enabled {
            continue
        }

        if re.policyMatches(policy, alert) {
            matching = append(matching, policy)
        }
    }

    return matching
}

// policyMatches 策略匹配
func (re *ResponseEngine) policyMatches(policy ResponsePolicy, alert *SecurityAlert) bool {
    for _, condition := range policy.Conditions {
        if !re.evaluateCondition(condition, alert) {
            return false
        }
    }
    return true
}

// evaluateCondition 评估条件
func (re *ResponseEngine) evaluateCondition(condition Condition, alert *SecurityAlert) bool {
    var value interface{}

    switch condition.Field {
    case "severity":
        value = alert.Severity
    case "source":
        value = alert.Source
    case "target":
        value = alert.Target
    case "rule_id":
        value = alert.RuleID
    default:
        return false
    }

    switch condition.Operator {
    case "equals":
        return value == condition.Value
    case "not_equals":
        return value != condition.Value
    case "contains":
        if str, ok := value.(string); ok {
            if target, ok := condition.Value.(string); ok {
                return strings.Contains(str, target)
            }
        }
        return false
    default:
        return false
    }
}

// executeAction 执行动作
func (re *ResponseEngine) executeAction(action ResponseAction, alert *SecurityAlert) error {
    // 验证动作
    if err := action.Validate(alert); err != nil {
        return fmt.Errorf("action validation failed: %w", err)
    }

    // 执行动作
    ctx, cancel := context.WithTimeout(re.ctx, re.config.ResponseDelay)
    defer cancel()

    return action.Execute(ctx, alert)
}
```

### 4.2 阻断动作

```go
// BlockAction 阻断动作
type BlockAction struct {
    firewall FirewallManager
    duration time.Duration
}

// NewBlockAction 创建阻断动作
func NewBlockAction(firewall FirewallManager, duration time.Duration) *BlockAction {
    return &BlockAction{
        firewall: firewall,
        duration: duration,
    }
}

// Execute 执行阻断
func (ba *BlockAction) Execute(ctx context.Context, alert *SecurityAlert) error {
    // 创建阻断规则
    rule := &FirewallRule{
        ID:       uuid.New().String(),
        Action:   "DROP",
        SourceIP: alert.Source,
        DestIP:   alert.Target,
        Protocol: "ANY",
        Created:  time.Now(),
        Expires:  time.Now().Add(ba.duration),
    }

    // 添加到防火墙
    if err := ba.firewall.AddRule(rule); err != nil {
        return fmt.Errorf("failed to add firewall rule: %w", err)
    }

    // 设置过期清理
    go func() {
        time.Sleep(ba.duration)
        ba.firewall.RemoveRule(rule.ID)
    }()

    return nil
}

// Validate 验证动作
func (ba *BlockAction) Validate(alert *SecurityAlert) error {
    if alert.Source == "" {
        return errors.New("source IP is required")
    }

    if alert.Target == "" {
        return errors.New("target IP is required")
    }

    return nil
}

// ID 获取动作ID
func (ba *BlockAction) ID() string {
    return "block"
}

// Name 获取动作名称
func (ba *BlockAction) Name() string {
    return "Block IP"
}

// FirewallManager 防火墙管理器接口
type FirewallManager interface {
    AddRule(rule *FirewallRule) error
    RemoveRule(ruleID string) error
    ListRules() ([]*FirewallRule, error)
}

// FirewallRule 防火墙规则
type FirewallRule struct {
    ID       string    `json:"id"`
    Action   string    `json:"action"`
    SourceIP string    `json:"source_ip"`
    DestIP   string    `json:"dest_ip"`
    Protocol string    `json:"protocol"`
    Created  time.Time `json:"created"`
    Expires  time.Time `json:"expires"`
}
```

## 5. 监控与报告

### 5.1 性能监控

```go
// IDSMetrics IDS指标
type IDSMetrics struct {
    eventsProcessed   int64
    alertsGenerated   int64
    falsePositives    int64
    truePositives     int64
    avgProcessingTime time.Duration
    mu                sync.RWMutex
}

// RecordEvent 记录事件
func (m *IDSMetrics) RecordEvent(processingTime time.Duration) {
    atomic.AddInt64(&m.eventsProcessed, 1)

    m.mu.Lock()
    defer m.mu.Unlock()

    // 更新平均处理时间
    total := m.avgProcessingTime * time.Duration(m.eventsProcessed-1)
    m.avgProcessingTime = (total + processingTime) / time.Duration(m.eventsProcessed)
}

// RecordAlert 记录告警
func (m *IDSMetrics) RecordAlert(isTruePositive bool) {
    atomic.AddInt64(&m.alertsGenerated, 1)

    if isTruePositive {
        atomic.AddInt64(&m.truePositives, 1)
    } else {
        atomic.AddInt64(&m.falsePositives, 1)
    }
}

// GetMetrics 获取指标
func (m *IDSMetrics) GetMetrics() map[string]interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()

    totalAlerts := m.alertsGenerated
    precision := float64(0)
    if totalAlerts > 0 {
        precision = float64(m.truePositives) / float64(totalAlerts) * 100
    }

    return map[string]interface{}{
        "events_processed":    m.eventsProcessed,
        "alerts_generated":    m.alertsGenerated,
        "true_positives":      m.truePositives,
        "false_positives":     m.falsePositives,
        "precision":           precision,
        "avg_processing_time": m.avgProcessingTime,
    }
}
```

## 6. 总结

入侵检测系统是网络安全防护的重要组成部分，通过实时监控和分析网络流量，能够及时发现和响应安全威胁。本模块提供了：

1. **完整的检测架构**：事件采集、预处理、检测引擎、告警处理
2. **多种检测方法**：签名检测、异常检测、行为分析
3. **智能响应机制**：自动响应、策略匹配、动作执行
4. **性能监控**：指标收集、性能分析、报告生成

通过Go语言的高性能和并发特性，实现了高效、可靠的入侵检测系统，为网络安全防护提供了强有力的技术支撑。

---

**相关链接**：
- [01-安全扫描工具](../01-Security-Scanning-Tools/README.md)
- [03-加密服务](../03-Encryption-Services/README.md)
- [04-身份认证](../04-Identity-Authentication/README.md)
