# Prometheus 集成实现方案

## 目录

1. [设计目标](#设计目标)
2. [架构设计](#架构设计)
3. [核心组件](#核心组件)
4. [实现细节](#实现细节)
5. [使用示例](#使用示例)
6. [最佳实践](#最佳实践)
7. [性能优化](#性能优化)

## 设计目标

### 1.1 核心目标

1. **统一指标收集**: 为所有组件提供统一的指标收集接口
2. **自动指标暴露**: 自动暴露指标到HTTP端点
3. **自定义指标支持**: 支持自定义业务指标
4. **性能优化**: 最小化指标收集的性能影响
5. **配置灵活**: 支持灵活的指标配置

### 1.2 设计原则

- **标准化**: 遵循Prometheus指标命名规范
- **可观测性**: 提供完整的系统可观测性
- **低侵入性**: 最小化对现有代码的影响
- **高性能**: 使用高效的指标收集机制

## 架构设计

### 2.1 整体架构

```text
┌─────────────────────────────────────┐
│           Application               │
├─────────────────────────────────────┤
│         Metrics Registry            │
├─────────────────────────────────────┤
│         Metrics Collector           │
├─────────────────────────────────────┤
│         HTTP Exporter               │
├─────────────────────────────────────┤
│         Prometheus Server           │
└─────────────────────────────────────┘
```

### 2.2 组件关系

```text
Component
├── MetricsProvider
│   ├── CounterProvider
│   ├── GaugeProvider
│   ├── HistogramProvider
│   └── SummaryProvider
├── MetricsRegistry
├── MetricsCollector
└── HTTPExporter
```

## 核心组件

### 3.1 指标提供者接口

```go
// 指标提供者接口
type MetricsProvider interface {
    Name() string
    Help() string
    Labels() []string
    Collect(ch chan<- prometheus.Metric)
}

// 计数器提供者
type CounterProvider interface {
    MetricsProvider
    Increment(labels prometheus.Labels)
    IncrementBy(value float64, labels prometheus.Labels)
    GetValue(labels prometheus.Labels) float64
}

// 仪表提供者
type GaugeProvider interface {
    MetricsProvider
    Set(value float64, labels prometheus.Labels)
    Increment(labels prometheus.Labels)
    Decrement(labels prometheus.Labels)
    GetValue(labels prometheus.Labels) float64
}

// 直方图提供者
type HistogramProvider interface {
    MetricsProvider
    Observe(value float64, labels prometheus.Labels)
    GetBuckets(labels prometheus.Labels) map[float64]uint64
}

// 摘要提供者
type SummaryProvider interface {
    MetricsProvider
    Observe(value float64, labels prometheus.Labels)
    GetQuantiles(labels prometheus.Labels) map[float64]float64
}
```

### 3.2 指标注册表

```go
// 指标注册表
type MetricsRegistry struct {
    registry  *prometheus.Registry
    providers map[string]MetricsProvider
    logger    *zap.Logger
    mu        sync.RWMutex
}

func NewMetricsRegistry() *MetricsRegistry {
    return &MetricsRegistry{
        registry:  prometheus.NewRegistry(),
        providers: make(map[string]MetricsProvider),
        logger:    zap.L().Named("metrics-registry"),
    }
}

func (mr *MetricsRegistry) Register(provider MetricsProvider) error {
    mr.mu.Lock()
    defer mr.mu.Unlock()
    
    if _, exists := mr.providers[provider.Name()]; exists {
        return fmt.Errorf("metrics provider %s already registered", provider.Name())
    }
    
    // 创建Prometheus指标
    metric := mr.createPrometheusMetric(provider)
    if err := mr.registry.Register(metric); err != nil {
        return fmt.Errorf("failed to register metric: %w", err)
    }
    
    mr.providers[provider.Name()] = provider
    mr.logger.Info("metrics provider registered", zap.String("name", provider.Name()))
    
    return nil
}

func (mr *MetricsRegistry) Unregister(name string) error {
    mr.mu.Lock()
    defer mr.mu.Unlock()
    
    if _, exists := mr.providers[name]; !exists {
        return fmt.Errorf("metrics provider %s not found", name)
    }
    
    delete(mr.providers, name)
    mr.logger.Info("metrics provider unregistered", zap.String("name", name))
    
    return nil
}

func (mr *MetricsRegistry) GetRegistry() *prometheus.Registry {
    return mr.registry
}

func (mr *MetricsRegistry) createPrometheusMetric(provider MetricsProvider) prometheus.Collector {
    switch p := provider.(type) {
    case CounterProvider:
        return prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: p.Name(),
                Help: p.Help(),
            },
            p.Labels(),
        )
    case GaugeProvider:
        return prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: p.Name(),
                Help: p.Help(),
            },
            p.Labels(),
        )
    case HistogramProvider:
        return prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: p.Name(),
                Help: p.Help(),
            },
            p.Labels(),
        )
    case SummaryProvider:
        return prometheus.NewSummaryVec(
            prometheus.SummaryOpts{
                Name: p.Name(),
                Help: p.Help(),
            },
            p.Labels(),
        )
    default:
        return nil
    }
}
```

### 3.3 指标收集器

```go
// 指标收集器
type MetricsCollector struct {
    registry *MetricsRegistry
    logger   *zap.Logger
    metrics  CollectorMetrics
}

func NewMetricsCollector(registry *MetricsRegistry) *MetricsCollector {
    return &MetricsCollector{
        registry: registry,
        logger:   zap.L().Named("metrics-collector"),
        metrics:  NewCollectorMetrics(),
    }
}

func (mc *MetricsCollector) CollectMetrics() {
    mc.metrics.CollectionCount.Inc()
    
    registry := mc.registry.GetRegistry()
    metrics, err := registry.Gather()
    if err != nil {
        mc.metrics.CollectionErrors.Inc()
        mc.logger.Error("failed to gather metrics", zap.Error(err))
        return
    }
    
    mc.metrics.MetricsCollected.Add(float64(len(metrics)))
    mc.logger.Debug("metrics collected", zap.Int("count", len(metrics)))
}

func (mc *MetricsCollector) StartPeriodicCollection(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            mc.CollectMetrics()
        }
    }
}
```

## 实现细节

### 4.1 计数器实现

```go
// 计数器实现
type CounterImpl struct {
    name   string
    help   string
    labels []string
    values map[string]float64
    mu     sync.RWMutex
}

func NewCounter(name, help string, labels []string) *CounterImpl {
    return &CounterImpl{
        name:   name,
        help:   help,
        labels: labels,
        values: make(map[string]float64),
    }
}

func (c *CounterImpl) Name() string {
    return c.name
}

func (c *CounterImpl) Help() string {
    return c.help
}

func (c *CounterImpl) Labels() []string {
    return c.labels
}

func (c *CounterImpl) Increment(labels prometheus.Labels) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    key := c.getLabelKey(labels)
    c.values[key]++
}

func (c *CounterImpl) IncrementBy(value float64, labels prometheus.Labels) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    key := c.getLabelKey(labels)
    c.values[key] += value
}

func (c *CounterImpl) GetValue(labels prometheus.Labels) float64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    key := c.getLabelKey(labels)
    return c.values[key]
}

func (c *CounterImpl) Collect(ch chan<- prometheus.Metric) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    for key, value := range c.values {
        labels := c.parseLabelKey(key)
        metric := prometheus.MustNewConstMetric(
            prometheus.NewDesc(c.name, c.help, c.labels, nil),
            prometheus.CounterValue,
            value,
            labels...,
        )
        ch <- metric
    }
}

func (c *CounterImpl) getLabelKey(labels prometheus.Labels) string {
    var keys []string
    for _, label := range c.labels {
        keys = append(keys, labels[label])
    }
    return strings.Join(keys, "|")
}

func (c *CounterImpl) parseLabelKey(key string) []string {
    return strings.Split(key, "|")
}
```

### 4.2 仪表实现

```go
// 仪表实现
type GaugeImpl struct {
    name   string
    help   string
    labels []string
    values map[string]float64
    mu     sync.RWMutex
}

func NewGauge(name, help string, labels []string) *GaugeImpl {
    return &GaugeImpl{
        name:   name,
        help:   help,
        labels: labels,
        values: make(map[string]float64),
    }
}

func (g *GaugeImpl) Name() string {
    return g.name
}

func (g *GaugeImpl) Help() string {
    return g.help
}

func (g *GaugeImpl) Labels() []string {
    return g.labels
}

func (g *GaugeImpl) Set(value float64, labels prometheus.Labels) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    key := g.getLabelKey(labels)
    g.values[key] = value
}

func (g *GaugeImpl) Increment(labels prometheus.Labels) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    key := g.getLabelKey(labels)
    g.values[key]++
}

func (g *GaugeImpl) Decrement(labels prometheus.Labels) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    key := g.getLabelKey(labels)
    g.values[key]--
}

func (g *GaugeImpl) GetValue(labels prometheus.Labels) float64 {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    key := g.getLabelKey(labels)
    return g.values[key]
}

func (g *GaugeImpl) Collect(ch chan<- prometheus.Metric) {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    for key, value := range g.values {
        labels := g.parseLabelKey(key)
        metric := prometheus.MustNewConstMetric(
            prometheus.NewDesc(g.name, g.help, g.labels, nil),
            prometheus.GaugeValue,
            value,
            labels...,
        )
        ch <- metric
    }
}

func (g *GaugeImpl) getLabelKey(labels prometheus.Labels) string {
    var keys []string
    for _, label := range g.labels {
        keys = append(keys, labels[label])
    }
    return strings.Join(keys, "|")
}

func (g *GaugeImpl) parseLabelKey(key string) []string {
    return strings.Split(key, "|")
}
```

### 4.3 HTTP导出器

```go
// HTTP导出器
type HTTPExporter struct {
    addr     string
    registry *MetricsRegistry
    logger   *zap.Logger
    metrics  ExporterMetrics
}

func NewHTTPExporter(addr string, registry *MetricsRegistry) *HTTPExporter {
    return &HTTPExporter{
        addr:     addr,
        registry: registry,
        logger:   zap.L().Named("http-exporter"),
        metrics:  NewExporterMetrics(),
    }
}

func (he *HTTPExporter) Start() error {
    // 创建HTTP处理器
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        he.handleMetrics(w, r)
    })
    
    // 启动HTTP服务器
    go func() {
        if err := http.ListenAndServe(he.addr, handler); err != nil {
            he.logger.Error("http exporter failed", zap.Error(err))
        }
    }()
    
    he.logger.Info("http exporter started", zap.String("addr", he.addr))
    return nil
}

func (he *HTTPExporter) handleMetrics(w http.ResponseWriter, r *http.Request) {
    he.metrics.RequestCount.Inc()
    
    // 设置响应头
    w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
    
    // 收集指标
    registry := he.registry.GetRegistry()
    metrics, err := registry.Gather()
    if err != nil {
        he.metrics.RequestErrors.Inc()
        http.Error(w, fmt.Sprintf("failed to gather metrics: %v", err), http.StatusInternalServerError)
        return
    }
    
    // 编码指标
    encoder := prometheus.NewEncoder(prometheus.NewTextEncoder(w))
    for _, metric := range metrics {
        if err := encoder.Encode(metric); err != nil {
            he.metrics.RequestErrors.Inc()
            he.logger.Error("failed to encode metric", zap.Error(err))
        }
    }
    
    he.metrics.RequestSuccess.Inc()
}
```

### 4.4 组件指标集成

```go
// 组件指标集成
type ComponentMetrics struct {
    statusGauge    *GaugeImpl
    uptimeGauge    *GaugeImpl
    startCounter   *CounterImpl
    stopCounter    *CounterImpl
    errorCounter   *CounterImpl
    lastErrorGauge *GaugeImpl
    startTime      time.Time
    stopTime       time.Time
}

func NewComponentMetrics(componentID string) *ComponentMetrics {
    return &ComponentMetrics{
        statusGauge: NewGauge(
            "component_status",
            "Component status (0=stopped, 1=starting, 2=running, 3=stopping, 4=error)",
            []string{"component_id", "component_type"},
        ),
        uptimeGauge: NewGauge(
            "component_uptime_seconds",
            "Component uptime in seconds",
            []string{"component_id"},
        ),
        startCounter: NewCounter(
            "component_start_total",
            "Total number of component starts",
            []string{"component_id"},
        ),
        stopCounter: NewCounter(
            "component_stop_total",
            "Total number of component stops",
            []string{"component_id"},
        ),
        errorCounter: NewCounter(
            "component_error_total",
            "Total number of component errors",
            []string{"component_id", "error_type"},
        ),
        lastErrorGauge: NewGauge(
            "component_last_error_timestamp",
            "Timestamp of last error",
            []string{"component_id"},
        ),
    }
}

func (cm *ComponentMetrics) SetStatus(status ComponentStatus, componentID, componentType string) {
    cm.statusGauge.Set(float64(status), prometheus.Labels{
        "component_id":   componentID,
        "component_type": componentType,
    })
}

func (cm *ComponentMetrics) SetUptime(componentID string) {
    if cm.startTime.IsZero() {
        return
    }
    
    var uptime float64
    if cm.stopTime.IsZero() {
        uptime = time.Since(cm.startTime).Seconds()
    } else {
        uptime = cm.stopTime.Sub(cm.startTime).Seconds()
    }
    
    cm.uptimeGauge.Set(uptime, prometheus.Labels{
        "component_id": componentID,
    })
}

func (cm *ComponentMetrics) IncrementStartCount(componentID string) {
    cm.startCounter.Increment(prometheus.Labels{
        "component_id": componentID,
    })
    cm.startTime = time.Now()
}

func (cm *ComponentMetrics) IncrementStopCount(componentID string) {
    cm.stopCounter.Increment(prometheus.Labels{
        "component_id": componentID,
    })
    cm.stopTime = time.Now()
}

func (cm *ComponentMetrics) IncrementErrorCount(componentID, errorType string) {
    cm.errorCounter.Increment(prometheus.Labels{
        "component_id": componentID,
        "error_type":   errorType,
    })
    
    cm.lastErrorGauge.Set(float64(time.Now().Unix()), prometheus.Labels{
        "component_id": componentID,
    })
}

func (cm *ComponentMetrics) GetProviders() []MetricsProvider {
    return []MetricsProvider{
        cm.statusGauge,
        cm.uptimeGauge,
        cm.startCounter,
        cm.stopCounter,
        cm.errorCounter,
        cm.lastErrorGauge,
    }
}
```

## 使用示例

### 5.1 基础使用

```go
func main() {
    // 创建指标注册表
    registry := NewMetricsRegistry()
    
    // 创建HTTP导出器
    exporter := NewHTTPExporter(":8080", registry)
    
    // 创建组件指标
    componentMetrics := NewComponentMetrics("user-service")
    
    // 注册指标提供者
    for _, provider := range componentMetrics.GetProviders() {
        if err := registry.Register(provider); err != nil {
            log.Fatal(err)
        }
    }
    
    // 启动HTTP导出器
    if err := exporter.Start(); err != nil {
        log.Fatal(err)
    }
    
    // 模拟组件操作
    componentMetrics.SetStatus(StatusRunning, "user-service", "service")
    componentMetrics.IncrementStartCount("user-service")
    
    // 更新运行时间
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                componentMetrics.SetUptime("user-service")
            }
        }
    }()
    
    // 等待信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}
```

### 5.2 自定义指标

```go
// 自定义业务指标
type BusinessMetrics struct {
    requestCounter   *CounterImpl
    responseTimeHistogram *HistogramImpl
    activeUsersGauge *GaugeImpl
}

func NewBusinessMetrics() *BusinessMetrics {
    return &BusinessMetrics{
        requestCounter: NewCounter(
            "business_requests_total",
            "Total number of business requests",
            []string{"service", "method", "status"},
        ),
        responseTimeHistogram: NewHistogram(
            "business_response_time_seconds",
            "Response time in seconds",
            []string{"service", "method"},
            []float64{0.1, 0.5, 1.0, 2.0, 5.0},
        ),
        activeUsersGauge: NewGauge(
            "business_active_users",
            "Number of active users",
            []string{"service"},
        ),
    }
}

func (bm *BusinessMetrics) RecordRequest(service, method, status string) {
    bm.requestCounter.Increment(prometheus.Labels{
        "service": service,
        "method":  method,
        "status":  status,
    })
}

func (bm *BusinessMetrics) RecordResponseTime(service, method string, duration time.Duration) {
    bm.responseTimeHistogram.Observe(duration.Seconds(), prometheus.Labels{
        "service": service,
        "method":  method,
    })
}

func (bm *BusinessMetrics) SetActiveUsers(service string, count int) {
    bm.activeUsersGauge.Set(float64(count), prometheus.Labels{
        "service": service,
    })
}

func (bm *BusinessMetrics) GetProviders() []MetricsProvider {
    return []MetricsProvider{
        bm.requestCounter,
        bm.responseTimeHistogram,
        bm.activeUsersGauge,
    }
}
```

### 5.3 中间件集成

```go
// HTTP中间件
func MetricsMiddleware(metrics *BusinessMetrics) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // 包装响应写入器以捕获状态码
            wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            // 调用下一个处理器
            next.ServeHTTP(wrappedWriter, r)
            
            // 记录指标
            duration := time.Since(start)
            service := "api"
            method := r.Method
            status := strconv.Itoa(wrappedWriter.statusCode)
            
            metrics.RecordRequest(service, method, status)
            metrics.RecordResponseTime(service, method, duration)
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## 最佳实践

### 6.1 指标命名规范

```go
// 指标命名规范
const (
    // 前缀规范
    PrefixComponent = "component_"
    PrefixBusiness  = "business_"
    PrefixSystem    = "system_"
    
    // 后缀规范
    SuffixTotal     = "_total"
    SuffixSeconds   = "_seconds"
    SuffixBytes     = "_bytes"
    SuffixRatio     = "_ratio"
    
    // 常用指标名称
    MetricRequestsTotal      = "requests_total"
    MetricResponseTime       = "response_time_seconds"
    MetricErrorRate          = "error_rate"
    MetricActiveConnections  = "active_connections"
    MetricMemoryUsage        = "memory_usage_bytes"
    MetricCPUUsage           = "cpu_usage_ratio"
)
```

### 6.2 标签使用规范

```go
// 标签使用规范
const (
    // 常用标签
    LabelService     = "service"
    LabelMethod      = "method"
    LabelStatus      = "status"
    LabelComponent   = "component"
    LabelInstance    = "instance"
    LabelVersion     = "version"
    LabelEnvironment = "environment"
    
    // 状态值
    StatusSuccess = "success"
    StatusError   = "error"
    StatusTimeout = "timeout"
    
    // 环境值
    EnvDevelopment = "development"
    EnvStaging     = "staging"
    EnvProduction  = "production"
)
```

### 6.3 性能优化

```go
// 性能优化配置
type MetricsConfig struct {
    // 收集间隔
    CollectionInterval time.Duration `yaml:"collection_interval"`
    
    // 缓冲区大小
    BufferSize int `yaml:"buffer_size"`
    
    // 并发收集器数量
    CollectorCount int `yaml:"collector_count"`
    
    // 指标缓存时间
    CacheDuration time.Duration `yaml:"cache_duration"`
    
    // 是否启用压缩
    EnableCompression bool `yaml:"enable_compression"`
}

func DefaultMetricsConfig() *MetricsConfig {
    return &MetricsConfig{
        CollectionInterval: 15 * time.Second,
        BufferSize:         1000,
        CollectorCount:     4,
        CacheDuration:      5 * time.Minute,
        EnableCompression:  true,
    }
}
```

## 性能优化

### 7.1 内存优化

```go
// 内存优化的指标存储
type OptimizedMetricsStorage struct {
    values sync.Map
    config *MetricsConfig
}

func NewOptimizedMetricsStorage(config *MetricsConfig) *OptimizedMetricsStorage {
    return &OptimizedMetricsStorage{
        config: config,
    }
}

func (oms *OptimizedMetricsStorage) Set(key string, value float64) {
    oms.values.Store(key, value)
}

func (oms *OptimizedMetricsStorage) Get(key string) (float64, bool) {
    if value, ok := oms.values.Load(key); ok {
        return value.(float64), true
    }
    return 0, false
}

func (oms *OptimizedMetricsStorage) Cleanup() {
    // 定期清理过期的指标
    oms.values.Range(func(key, value interface{}) bool {
        // 检查是否过期
        if oms.isExpired(key.(string)) {
            oms.values.Delete(key)
        }
        return true
    })
}

func (oms *OptimizedMetricsStorage) isExpired(key string) bool {
    // 实现过期检查逻辑
    return false
}
```

### 7.2 并发优化

```go
// 并发优化的指标收集器
type ConcurrentMetricsCollector struct {
    collectors []*MetricsCollector
    config     *MetricsConfig
    logger     *zap.Logger
}

func NewConcurrentMetricsCollector(registry *MetricsRegistry, config *MetricsConfig) *ConcurrentMetricsCollector {
    collectors := make([]*MetricsCollector, config.CollectorCount)
    for i := 0; i < config.CollectorCount; i++ {
        collectors[i] = NewMetricsCollector(registry)
    }
    
    return &ConcurrentMetricsCollector{
        collectors: collectors,
        config:     config,
        logger:     zap.L().Named("concurrent-metrics-collector"),
    }
}

func (cmc *ConcurrentMetricsCollector) StartCollection() {
    for i, collector := range cmc.collectors {
        go func(id int, c *MetricsCollector) {
            c.StartPeriodicCollection(cmc.config.CollectionInterval)
        }(i, collector)
    }
    
    cmc.logger.Info("concurrent metrics collection started", 
        zap.Int("collector_count", len(cmc.collectors)))
}
```

## 总结

Prometheus集成实现方案提供了以下核心功能：

1. **统一指标收集**: 标准化的指标收集接口
2. **自动指标暴露**: HTTP端点自动暴露指标
3. **自定义指标支持**: 灵活的自定义业务指标
4. **性能优化**: 内存和并发优化
5. **最佳实践**: 命名规范和标签使用规范
6. **中间件集成**: 简单的HTTP中间件集成

这个实现方案为Golang Common库提供了完整的Prometheus监控能力，可以显著提升系统的可观测性。
