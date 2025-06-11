# 监控修复方案

## 目录

1. [监控框架实现](#监控框架实现)
2. [指标收集实现](#指标收集实现)
3. [日志增强实现](#日志增强实现)
4. [追踪系统实现](#追踪系统实现)
5. [告警系统实现](#告警系统实现)

## 监控框架实现

### 1.1 统一监控管理器

```go
// 文件: monitoring/manager.go
package monitoring

import (
    "context"
    "fmt"
    "time"
    
    "go.uber.org/zap"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

// MonitoringManager 统一监控管理器
type MonitoringManager struct {
    metricsCollector    *SystemMetricsCollector
    businessMetrics     *BusinessMetricsCollector
    logManager          *EnhancedLogManager
    tracer              *Tracer
    prometheusExporter  *PrometheusExporter
    alertManager        *AlertManager
    logger              *zap.Logger
    config              MonitoringConfig
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
    PrometheusPort int           `json:"prometheus_port"`
    JaegerHost     string        `json:"jaeger_host"`
    JaegerPort     int           `json:"jaeger_port"`
    SamplingRate   float64       `json:"sampling_rate"`
    LogLevel       string        `json:"log_level"`
    MetricsInterval time.Duration `json:"metrics_interval"`
    AlertRules     []AlertRule   `json:"alert_rules"`
}

// NewMonitoringManager 创建监控管理器
func NewMonitoringManager(config MonitoringConfig) *MonitoringManager {
    // 创建Prometheus导出器
    prometheusExporter := NewPrometheusExporter(config.PrometheusPort)
    
    // 创建指标收集器
    systemMetrics := NewSystemMetricsCollector()
    businessMetrics := NewBusinessMetricsCollector()
    
    // 注册指标
    prometheusExporter.Register(systemMetrics.metrics)
    prometheusExporter.Register(businessMetrics)
    
    // 创建日志管理器
    logConfig := LogConfig{
        Level:      config.LogLevel,
        OutputPath: "logs/app.log",
        MaxSize:    100,
        MaxBackups: 3,
        MaxAge:     7,
        Compress:   true,
    }
    logManager := NewEnhancedLogManager(logConfig)
    
    // 创建追踪器
    tracer := NewTracer()
    
    // 创建告警管理器
    alertManager := NewAlertManager(businessMetrics)
    
    return &MonitoringManager{
        metricsCollector:    systemMetrics,
        businessMetrics:     businessMetrics,
        logManager:          logManager,
        tracer:              tracer,
        prometheusExporter:  prometheusExporter,
        alertManager:        alertManager,
        logger:              zap.L().Named("monitoring-manager"),
        config:              config,
    }
}

// Start 启动监控管理器
func (mm *MonitoringManager) Start() error {
    // 启动Prometheus导出器
    if err := mm.prometheusExporter.Start(); err != nil {
        return fmt.Errorf("failed to start prometheus exporter: %w", err)
    }
    
    // 启动系统指标收集
    mm.metricsCollector.Start()
    
    // 启动告警管理器
    mm.alertManager.Start()
    
    mm.logger.Info("monitoring manager started")
    return nil
}

// Stop 停止监控管理器
func (mm *MonitoringManager) Stop() error {
    // 停止系统指标收集
    mm.metricsCollector.Stop()
    
    // 停止告警管理器
    mm.alertManager.Stop()
    
    // 停止Prometheus导出器
    if err := mm.prometheusExporter.Stop(); err != nil {
        return fmt.Errorf("failed to stop prometheus exporter: %w", err)
    }
    
    mm.logger.Info("monitoring manager stopped")
    return nil
}

// GetMetricsCollector 获取指标收集器
func (mm *MonitoringManager) GetMetricsCollector() *BusinessMetricsCollector {
    return mm.businessMetrics
}

// GetLogManager 获取日志管理器
func (mm *MonitoringManager) GetLogManager() *EnhancedLogManager {
    return mm.logManager
}

// GetTracer 获取追踪器
func (mm *MonitoringManager) GetTracer() *Tracer {
    return mm.tracer
}
```

### 1.2 监控中间件

```go
// 文件: monitoring/middleware.go
package monitoring

import (
    "time"
    
    "go.uber.org/zap"
    "common/model/component"
)

// MonitoringMiddleware 监控中间件
type MonitoringMiddleware struct {
    manager *MonitoringManager
    logger  *zap.Logger
}

// NewMonitoringMiddleware 创建监控中间件
func NewMonitoringMiddleware(manager *MonitoringManager) *MonitoringMiddleware {
    return &MonitoringMiddleware{
        manager: manager,
        logger:  zap.L().Named("monitoring-middleware"),
    }
}

// WrapComponent 包装组件
func (mm *MonitoringMiddleware) WrapComponent(component component.Component) component.Component {
    return &MonitoredComponent{
        component: component,
        manager:   mm.manager,
        logger:    mm.logger,
    }
}

// MonitoredComponent 监控组件
type MonitoredComponent struct {
    component component.Component
    manager   *MonitoringManager
    logger    *zap.Logger
}

// ID 返回组件ID
func (mc *MonitoredComponent) ID() string {
    return mc.component.ID()
}

// Kind 返回组件类型
func (mc *MonitoredComponent) Kind() string {
    return mc.component.Kind()
}

// Start 启动组件
func (mc *MonitoredComponent) Start() error {
    startTime := time.Now()
    
    // 记录日志
    mc.manager.GetLogManager().LogComponentLifecycle(mc.component, "start", nil)
    
    // 开始追踪
    span := mc.manager.GetTracer().StartComponentSpan(mc.component, "component.start")
    defer span.Finish()
    
    // 执行启动
    err := mc.component.Start()
    
    // 记录指标
    duration := time.Since(startTime)
    mc.manager.GetMetricsCollector().RecordComponentStart(mc.component.Kind())
    
    if err != nil {
        span.SetTag("error", true)
        span.LogKV("error.message", err.Error())
        mc.manager.GetLogManager().LogComponentLifecycle(mc.component, "start", err)
    } else {
        span.SetTag("error", false)
        span.SetTag("duration", duration.Seconds())
    }
    
    return err
}

// Stop 停止组件
func (mc *MonitoredComponent) Stop() error {
    startTime := time.Now()
    
    // 记录日志
    mc.manager.GetLogManager().LogComponentLifecycle(mc.component, "stop", nil)
    
    // 开始追踪
    span := mc.manager.GetTracer().StartComponentSpan(mc.component, "component.stop")
    defer span.Finish()
    
    // 执行停止
    err := mc.component.Stop()
    
    // 记录指标
    duration := time.Since(startTime)
    mc.manager.GetMetricsCollector().RecordComponentStop(mc.component.Kind())
    
    if err != nil {
        span.SetTag("error", true)
        span.LogKV("error.message", err.Error())
        mc.manager.GetLogManager().LogComponentLifecycle(mc.component, "stop", err)
    } else {
        span.SetTag("error", false)
        span.SetTag("duration", duration.Seconds())
    }
    
    return err
}

// IsRunning 检查组件是否运行
func (mc *MonitoredComponent) IsRunning() bool {
    return mc.component.IsRunning()
}
```

## 指标收集实现

### 2.1 系统指标收集器

```go
// 文件: monitoring/system_metrics.go
package monitoring

import (
    "runtime"
    "time"
    
    "go.uber.org/zap"
    "github.com/prometheus/client_golang/prometheus"
)

// SystemMetricsCollector 系统指标收集器
type SystemMetricsCollector struct {
    metrics  *SystemMetrics
    logger   *zap.Logger
    interval time.Duration
    stopChan chan struct{}
}

// SystemMetrics 系统指标
type SystemMetrics struct {
    GoroutineCount  prometheus.Gauge
    MemoryUsage     prometheus.Gauge
    CPUUsage        prometheus.Gauge
    GCStats         prometheus.Counter
    FileDescriptors prometheus.Gauge
}

// NewSystemMetricsCollector 创建系统指标收集器
func NewSystemMetricsCollector() *SystemMetricsCollector {
    metrics := &SystemMetrics{
        GoroutineCount: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "go_goroutines",
            Help: "Number of goroutines",
        }),
        MemoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "go_memory_usage_bytes",
            Help: "Memory usage in bytes",
        }),
        CPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "go_cpu_usage_percent",
            Help: "CPU usage percentage",
        }),
        GCStats: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "go_gc_cycles_total",
            Help: "Total number of GC cycles",
        }),
        FileDescriptors: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "go_file_descriptors",
            Help: "Number of file descriptors",
        }),
    }
    
    return &SystemMetricsCollector{
        metrics:  metrics,
        logger:   zap.L().Named("system-metrics"),
        interval: time.Second * 30,
        stopChan: make(chan struct{}),
    }
}

// Start 启动指标收集
func (smc *SystemMetricsCollector) Start() {
    go smc.collectLoop()
    smc.logger.Info("system metrics collector started")
}

// Stop 停止指标收集
func (smc *SystemMetricsCollector) Stop() {
    close(smc.stopChan)
    smc.logger.Info("system metrics collector stopped")
}

// collectLoop 收集循环
func (smc *SystemMetricsCollector) collectLoop() {
    ticker := time.NewTicker(smc.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            smc.collect()
        case <-smc.stopChan:
            return
        }
    }
}

// collect 收集指标
func (smc *SystemMetricsCollector) collect() {
    // 收集goroutine数量
    smc.metrics.GoroutineCount.Set(float64(runtime.NumGoroutine()))
    
    // 收集内存使用
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    smc.metrics.MemoryUsage.Set(float64(m.Alloc))
    
    // 收集GC统计
    smc.metrics.GCStats.Add(float64(m.NumGC))
    
    smc.logger.Debug("system metrics collected")
}
```

### 2.2 业务指标收集器

```go
// 文件: monitoring/business_metrics.go
package monitoring

import (
    "time"
    
    "go.uber.org/zap"
    "github.com/prometheus/client_golang/prometheus"
)

// BusinessMetricsCollector 业务指标收集器
type BusinessMetricsCollector struct {
    ComponentStartTotal    prometheus.Counter
    ComponentStopTotal     prometheus.Counter
    ComponentErrorTotal    prometheus.Counter
    EventPublishTotal      prometheus.Counter
    EventSubscribeTotal    prometheus.Counter
    EventProcessDuration   prometheus.Histogram
    logger                 *zap.Logger
}

// NewBusinessMetricsCollector 创建业务指标收集器
func NewBusinessMetricsCollector() *BusinessMetricsCollector {
    return &BusinessMetricsCollector{
        ComponentStartTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_start_total",
            Help: "Total number of component starts",
        }),
        ComponentStopTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_stop_total",
            Help: "Total number of component stops",
        }),
        ComponentErrorTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_error_total",
            Help: "Total number of component errors",
        }),
        EventPublishTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "event_publish_total",
            Help: "Total number of events published",
        }),
        EventSubscribeTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "event_subscribe_total",
            Help: "Total number of event subscriptions",
        }),
        EventProcessDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "event_process_duration_seconds",
            Help:    "Event processing duration in seconds",
            Buckets: prometheus.DefBuckets,
        }),
        logger: zap.L().Named("business-metrics"),
    }
}

// RecordComponentStart 记录组件启动
func (bmc *BusinessMetricsCollector) RecordComponentStart(componentType string) {
    bmc.ComponentStartTotal.Inc()
    bmc.logger.Debug("component start recorded", zap.String("type", componentType))
}

// RecordComponentStop 记录组件停止
func (bmc *BusinessMetricsCollector) RecordComponentStop(componentType string) {
    bmc.ComponentStopTotal.Inc()
    bmc.logger.Debug("component stop recorded", zap.String("type", componentType))
}

// RecordComponentError 记录组件错误
func (bmc *BusinessMetricsCollector) RecordComponentError(componentType string, errorType string) {
    bmc.ComponentErrorTotal.Inc()
    bmc.logger.Warn("component error recorded", 
        zap.String("type", componentType),
        zap.String("error_type", errorType))
}

// RecordEventPublish 记录事件发布
func (bmc *BusinessMetricsCollector) RecordEventPublish(topic string) {
    bmc.EventPublishTotal.Inc()
    bmc.logger.Debug("event publish recorded", zap.String("topic", topic))
}

// RecordEventSubscribe 记录事件订阅
func (bmc *BusinessMetricsCollector) RecordEventSubscribe(topic string) {
    bmc.EventSubscribeTotal.Inc()
    bmc.logger.Debug("event subscribe recorded", zap.String("topic", topic))
}

// RecordEventProcessDuration 记录事件处理时长
func (bmc *BusinessMetricsCollector) RecordEventProcessDuration(duration time.Duration) {
    bmc.EventProcessDuration.Observe(duration.Seconds())
    bmc.logger.Debug("event process duration recorded", 
        zap.Duration("duration", duration))
}
```

## 日志增强实现

### 3.1 增强的日志管理器

```go
// 文件: monitoring/log_manager.go
package monitoring

import (
    "fmt"
    "time"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "common/model/component"
)

// EnhancedLogManager 增强的日志管理器
type EnhancedLogManager struct {
    logger  *zap.Logger
    config  LogConfig
    metrics *BusinessMetricsCollector
    tracer  *Tracer
}

// LogConfig 日志配置
type LogConfig struct {
    Level      string            `json:"level"`
    OutputPath string            `json:"output_path"`
    MaxSize    int               `json:"max_size"`
    MaxBackups int               `json:"max_backups"`
    MaxAge     int               `json:"max_age"`
    Compress   bool              `json:"compress"`
    Fields     map[string]string `json:"fields"`
    Sampling   LogSampling       `json:"sampling"`
}

// LogSampling 日志采样
type LogSampling struct {
    Initial    int `json:"initial"`
    Thereafter int `json:"thereafter"`
}

// NewEnhancedLogManager 创建增强的日志管理器
func NewEnhancedLogManager(config LogConfig) *EnhancedLogManager {
    // 配置Zap日志
    zapConfig := zap.NewProductionConfig()
    zapConfig.OutputPaths = []string{config.OutputPath}
    zapConfig.Level = getLogLevel(config.Level)
    
    // 添加采样配置
    if config.Sampling.Initial > 0 {
        zapConfig.Sampling = &zap.SamplingConfig{
            Initial:    config.Sampling.Initial,
            Thereafter: config.Sampling.Thereafter,
        }
    }
    
    logger, err := zapConfig.Build()
    if err != nil {
        panic(fmt.Sprintf("failed to build logger: %v", err))
    }
    
    return &EnhancedLogManager{
        logger:  logger,
        config:  config,
        metrics: NewBusinessMetricsCollector(),
        tracer:  NewTracer(),
    }
}

// LogComponentLifecycle 记录组件生命周期
func (elm *EnhancedLogManager) LogComponentLifecycle(component component.Component, event string, err error) {
    fields := []zap.Field{
        zap.String("component_id", component.ID()),
        zap.String("component_kind", component.Kind()),
        zap.String("event", event),
        zap.Time("timestamp", time.Now()),
    }
    
    if err != nil {
        fields = append(fields, zap.Error(err))
        elm.logger.Error("component lifecycle event", fields...)
        elm.metrics.RecordComponentError(component.Kind(), "lifecycle")
    } else {
        elm.logger.Info("component lifecycle event", fields...)
    }
    
    // 记录指标
    switch event {
    case "start":
        elm.metrics.RecordComponentStart(component.Kind())
    case "stop":
        elm.metrics.RecordComponentStop(component.Kind())
    }
}

// LogEvent 记录事件
func (elm *EnhancedLogManager) LogEvent(event Event, action string) {
    fields := []zap.Field{
        zap.String("event_id", event.ID),
        zap.String("event_type", event.Type),
        zap.String("action", action),
        zap.Time("timestamp", event.Timestamp),
        zap.String("source", event.Source),
    }
    
    elm.logger.Info("event action", fields...)
    
    // 记录指标
    switch action {
    case "publish":
        elm.metrics.RecordEventPublish(event.Type)
    case "subscribe":
        elm.metrics.RecordEventSubscribe(event.Type)
    }
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zap.AtomicLevel {
    switch level {
    case "debug":
        return zap.NewAtomicLevelAt(zapcore.DebugLevel)
    case "info":
        return zap.NewAtomicLevelAt(zapcore.InfoLevel)
    case "warn":
        return zap.NewAtomicLevelAt(zapcore.WarnLevel)
    case "error":
        return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
    default:
        return zap.NewAtomicLevelAt(zapcore.InfoLevel)
    }
}
```

## 追踪系统实现

### 4.1 分布式追踪器

```go
// 文件: monitoring/tracer.go
package monitoring

import (
    "fmt"
    "time"
    
    "go.uber.org/zap"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
    "common/model/component"
)

// Tracer 分布式追踪器
type Tracer struct {
    tracer  opentracing.Tracer
    logger  *zap.Logger
    closer  io.Closer
}

// Event 事件
type Event struct {
    ID        string
    Type      string
    Data      interface{}
    Timestamp time.Time
    Source    string
}

// NewTracer 创建追踪器
func NewTracer() *Tracer {
    // 配置Jaeger
    cfg := &config.Configuration{
        ServiceName: "golang-common",
        Sampler: &config.SamplerConfig{
            Type:  "probabilistic",
            Param: 0.1, // 采样率10%
        },
        Reporter: &config.ReporterConfig{
            LocalAgentHostPort: "localhost:6831",
        },
    }
    
    tracer, closer, err := cfg.NewTracer()
    if err != nil {
        panic(fmt.Sprintf("failed to create tracer: %v", err))
    }
    
    return &Tracer{
        tracer: tracer,
        logger: zap.L().Named("tracer"),
        closer: closer,
    }
}

// StartComponentSpan 开始组件追踪
func (t *Tracer) StartComponentSpan(component component.Component, operation string) opentracing.Span {
    span := t.tracer.StartSpan(operation)
    
    // 添加组件标签
    span.SetTag("component.id", component.ID())
    span.SetTag("component.kind", component.Kind())
    span.SetTag("operation", operation)
    
    return span
}

// StartEventSpan 开始事件追踪
func (t *Tracer) StartEventSpan(event Event, operation string) opentracing.Span {
    span := t.tracer.StartSpan(operation)
    
    // 添加事件标签
    span.SetTag("event.id", event.ID)
    span.SetTag("event.type", event.Type)
    span.SetTag("event.source", event.Source)
    span.SetTag("operation", operation)
    
    return span
}

// Close 关闭追踪器
func (t *Tracer) Close() error {
    if t.closer != nil {
        return t.closer.Close()
    }
    return nil
}
```

## 告警系统实现

### 5.1 告警管理器

```go
// 文件: monitoring/alert_manager.go
package monitoring

import (
    "fmt"
    "time"
    
    "go.uber.org/zap"
    "github.com/google/uuid"
)

// AlertManager 告警管理器
type AlertManager struct {
    rules     map[string]AlertRule
    metrics   *BusinessMetricsCollector
    logger    *zap.Logger
    stopChan  chan struct{}
}

// AlertRule 告警规则
type AlertRule struct {
    ID          string        `json:"id"`
    Name        string        `json:"name"`
    Description string        `json:"description"`
    Condition   string        `json:"condition"`
    Severity    string        `json:"severity"`
    Duration    time.Duration `json:"duration"`
    Enabled     bool          `json:"enabled"`
}

// Alert 告警
type Alert struct {
    ID          string    `json:"id"`
    RuleID      string    `json:"rule_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Severity    string    `json:"severity"`
    Timestamp   time.Time `json:"timestamp"`
}

// NewAlertManager 创建告警管理器
func NewAlertManager(metrics *BusinessMetricsCollector) *AlertManager {
    return &AlertManager{
        rules:   make(map[string]AlertRule),
        metrics: metrics,
        logger:  zap.L().Named("alert-manager"),
        stopChan: make(chan struct{}),
    }
}

// AddRule 添加告警规则
func (am *AlertManager) AddRule(rule AlertRule) {
    am.rules[rule.ID] = rule
    am.logger.Info("alert rule added", 
        zap.String("rule_id", rule.ID),
        zap.String("name", rule.Name))
}

// RemoveRule 移除告警规则
func (am *AlertManager) RemoveRule(ruleID string) {
    delete(am.rules, ruleID)
    am.logger.Info("alert rule removed", zap.String("rule_id", ruleID))
}

// Start 启动告警管理器
func (am *AlertManager) Start() {
    go am.evaluateLoop()
    am.logger.Info("alert manager started")
}

// Stop 停止告警管理器
func (am *AlertManager) Stop() {
    close(am.stopChan)
    am.logger.Info("alert manager stopped")
}

// evaluateLoop 评估循环
func (am *AlertManager) evaluateLoop() {
    ticker := time.NewTicker(time.Second * 30)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            alerts := am.EvaluateRules()
            for _, alert := range alerts {
                am.handleAlert(alert)
            }
        case <-am.stopChan:
            return
        }
    }
}

// EvaluateRules 评估规则
func (am *AlertManager) EvaluateRules() []Alert {
    var alerts []Alert
    
    for _, rule := range am.rules {
        if !rule.Enabled {
            continue
        }
        
        if am.evaluateCondition(rule.Condition) {
            alert := Alert{
                ID:          uuid.New().String(),
                RuleID:      rule.ID,
                Name:        rule.Name,
                Description: rule.Description,
                Severity:    rule.Severity,
                Timestamp:   time.Now(),
            }
            alerts = append(alerts, alert)
            
            am.logger.Warn("alert triggered", 
                zap.String("rule_id", rule.ID),
                zap.String("name", rule.Name),
                zap.String("severity", rule.Severity))
        }
    }
    
    return alerts
}

// evaluateCondition 评估条件
func (am *AlertManager) evaluateCondition(condition string) bool {
    // 简单的条件评估逻辑
    switch condition {
    case "high_error_rate":
        return am.metrics.ComponentErrorTotal.Counter().Get() > 100
    case "low_component_health":
        return am.metrics.ComponentStartTotal.Counter().Get() < 10
    default:
        return false
    }
}

// handleAlert 处理告警
func (am *AlertManager) handleAlert(alert Alert) {
    // 这里可以实现告警通知逻辑
    // 例如发送邮件、短信、Slack消息等
    
    am.logger.Error("alert triggered",
        zap.String("alert_id", alert.ID),
        zap.String("rule_id", alert.RuleID),
        zap.String("name", alert.Name),
        zap.String("severity", alert.Severity))
}
```

这个监控修复方案提供了完整的监控系统实现，包括：

1. **监控框架**: 统一的监控管理器和中间件
2. **指标收集**: 系统和业务指标收集器
3. **日志增强**: 结构化日志和生命周期记录
4. **追踪系统**: 分布式追踪和链路追踪
5. **告警系统**: 告警规则和通知机制

通过这些实现，可以显著提升Golang Common库的可观测性和运维能力。 