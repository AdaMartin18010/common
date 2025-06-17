# 1. 服务网格 (Service Mesh)

## 1.1 服务网格理论基础

### 1.1.1 服务网格形式化定义

**定义 1.1** (服务网格): 服务网格是一个四元组 $\mathcal{M} = (P, S, C, R)$，其中：

- $P$ 是代理集合 (Proxies)
- $S$ 是服务集合 (Services)
- $C$ 是控制平面 (Control Plane)
- $R$ 是路由规则集合 (Routing Rules)

**服务网格架构**:

```latex
\text{ServiceMesh} = \text{DataPlane} \times \text{ControlPlane}
```

其中：

- $\text{DataPlane} = \text{Proxy} \times \text{Service}$
- $\text{ControlPlane} = \text{Discovery} \times \text{Configuration} \times \text{Policy}$

### 1.1.2 服务网格核心概念

**定义 1.2** (服务代理): 服务代理是一个三元组 $\text{Proxy} = (I, O, F)$，其中：

- $I$ 是输入接口集合
- $O$ 是输出接口集合  
- $F$ 是过滤器集合

**定义 1.3** (服务发现): 服务发现是一个函数 $\text{Discovery}: \text{ServiceName} \rightarrow \mathcal{P}(\text{Endpoint})$，其中：

- $\text{ServiceName}$ 是服务名称集合
- $\text{Endpoint}$ 是服务端点集合

## 1.2 Go语言服务网格实现

### 1.2.1 代理实现

```go
// Proxy 服务代理接口
type Proxy interface {
    // 启动代理
    Start(ctx context.Context) error
    
    // 停止代理
    Stop(ctx context.Context) error
    
    // 添加过滤器
    AddFilter(filter Filter) error
    
    // 移除过滤器
    RemoveFilter(filterID string) error
    
    // 获取代理状态
    GetStatus() ProxyStatus
}

// ProxyStatus 代理状态
type ProxyStatus struct {
    ID        string    `json:"id"`
    Service   string    `json:"service"`
    Status    string    `json:"status"`
    StartTime time.Time `json:"start_time"`
    Metrics   Metrics   `json:"metrics"`
}

// Metrics 代理指标
type Metrics struct {
    RequestsTotal   int64         `json:"requests_total"`
    RequestsSuccess int64         `json:"requests_success"`
    RequestsFailed  int64         `json:"requests_failed"`
    AverageLatency  time.Duration `json:"average_latency"`
    LastUpdateTime  time.Time     `json:"last_update_time"`
}

// ServiceProxy 服务代理实现
type ServiceProxy struct {
    id           string
    serviceName  string
    config       *ProxyConfig
    filters      map[string]Filter
    discovery    ServiceDiscovery
    router       Router
    metrics      *Metrics
    mutex        sync.RWMutex
    stopChan     chan struct{}
}

// ProxyConfig 代理配置
type ProxyConfig struct {
    ListenPort    int           `json:"listen_port"`
    UpstreamPort  int           `json:"upstream_port"`
    BufferSize    int           `json:"buffer_size"`
    Timeout       time.Duration `json:"timeout"`
    MaxRetries    int           `json:"max_retries"`
    CircuitBreaker *CircuitBreakerConfig `json:"circuit_breaker"`
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
    Threshold     int           `json:"threshold"`
    Timeout       time.Duration `json:"timeout"`
    HalfOpenLimit int           `json:"half_open_limit"`
}

// NewServiceProxy 创建服务代理
func NewServiceProxy(id, serviceName string, config *ProxyConfig) *ServiceProxy {
    return &ServiceProxy{
        id:          id,
        serviceName: serviceName,
        config:      config,
        filters:     make(map[string]Filter),
        discovery:   NewServiceDiscovery(),
        router:      NewRouter(),
        metrics:     &Metrics{},
        stopChan:    make(chan struct{}),
    }
}

// Start 启动代理
func (sp *ServiceProxy) Start(ctx context.Context) error {
    // 启动监听器
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", sp.config.ListenPort))
    if err != nil {
        return fmt.Errorf("failed to start listener: %w", err)
    }
    
    go sp.acceptConnections(ctx, listener)
    
    return nil
}

// acceptConnections 接受连接
func (sp *ServiceProxy) acceptConnections(ctx context.Context, listener net.Listener) {
    defer listener.Close()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-sp.stopChan:
            return
        default:
            conn, err := listener.Accept()
            if err != nil {
                log.Printf("Accept connection error: %v", err)
                continue
            }
            
            go sp.handleConnection(ctx, conn)
        }
    }
}

// handleConnection 处理连接
func (sp *ServiceProxy) handleConnection(ctx context.Context, conn net.Conn) {
    defer conn.Close()
    
    // 创建请求上下文
    reqCtx := &RequestContext{
        ID:        uuid.New().String(),
        StartTime: time.Now(),
        Conn:      conn,
        Proxy:     sp,
    }
    
    // 应用入站过滤器
    if err := sp.applyInboundFilters(reqCtx); err != nil {
        log.Printf("Inbound filter error: %v", err)
        return
    }
    
    // 路由请求
    target, err := sp.router.Route(reqCtx)
    if err != nil {
        log.Printf("Routing error: %v", err)
        return
    }
    
    // 转发请求
    if err := sp.forwardRequest(reqCtx, target); err != nil {
        log.Printf("Forward error: %v", err)
        return
    }
    
    // 应用出站过滤器
    if err := sp.applyOutboundFilters(reqCtx); err != nil {
        log.Printf("Outbound filter error: %v", err)
        return
    }
    
    // 更新指标
    sp.updateMetrics(reqCtx)
}

// RequestContext 请求上下文
type RequestContext struct {
    ID        string
    StartTime time.Time
    EndTime   time.Time
    Conn      net.Conn
    Proxy     *ServiceProxy
    Headers   map[string]string
    Body      []byte
    Response  []byte
    Error     error
}

// applyInboundFilters 应用入站过滤器
func (sp *ServiceProxy) applyInboundFilters(ctx *RequestContext) error {
    sp.mutex.RLock()
    defer sp.mutex.RUnlock()
    
    for _, filter := range sp.filters {
        if inboundFilter, ok := filter.(InboundFilter); ok {
            if err := inboundFilter.OnInbound(ctx); err != nil {
                return err
            }
        }
    }
    return nil
}

// applyOutboundFilters 应用出站过滤器
func (sp *ServiceProxy) applyOutboundFilters(ctx *RequestContext) error {
    sp.mutex.RLock()
    defer sp.mutex.RUnlock()
    
    for _, filter := range sp.filters {
        if outboundFilter, ok := filter.(OutboundFilter); ok {
            if err := outboundFilter.OnOutbound(ctx); err != nil {
                return err
            }
        }
    }
    return nil
}

// forwardRequest 转发请求
func (sp *ServiceProxy) forwardRequest(ctx *RequestContext, target *Endpoint) error {
    // 连接到目标服务
    upstreamConn, err := net.Dial("tcp", target.Address)
    if err != nil {
        return fmt.Errorf("failed to connect to upstream: %w", err)
    }
    defer upstreamConn.Close()
    
    // 创建双向数据流
    errChan := make(chan error, 2)
    
    // 客户端到上游
    go func() {
        _, err := io.Copy(upstreamConn, ctx.Conn)
        errChan <- err
    }()
    
    // 上游到客户端
    go func() {
        _, err := io.Copy(ctx.Conn, upstreamConn)
        errChan <- err
    }()
    
    // 等待任一方向完成
    select {
    case err := <-errChan:
        return err
    case <-time.After(sp.config.Timeout):
        return fmt.Errorf("request timeout")
    }
}

// updateMetrics 更新指标
func (sp *ServiceProxy) updateMetrics(ctx *RequestContext) {
    sp.mutex.Lock()
    defer sp.mutex.Unlock()
    
    sp.metrics.RequestsTotal++
    if ctx.Error != nil {
        sp.metrics.RequestsFailed++
    } else {
        sp.metrics.RequestsSuccess++
    }
    
    duration := time.Since(ctx.StartTime)
    if sp.metrics.RequestsTotal > 0 {
        totalDuration := sp.metrics.AverageLatency * time.Duration(sp.metrics.RequestsTotal-1)
        sp.metrics.AverageLatency = (totalDuration + duration) / time.Duration(sp.metrics.RequestsTotal)
    }
    
    sp.metrics.LastUpdateTime = time.Now()
}
```

### 1.2.2 过滤器实现

```go
// Filter 过滤器接口
type Filter interface {
    GetID() string
    GetName() string
    GetType() FilterType
}

// FilterType 过滤器类型
type FilterType string

const (
    FilterTypeInbound  FilterType = "INBOUND"
    FilterTypeOutbound FilterType = "OUTBOUND"
    FilterTypeBoth     FilterType = "BOTH"
)

// InboundFilter 入站过滤器
type InboundFilter interface {
    Filter
    OnInbound(ctx *RequestContext) error
}

// OutboundFilter 出站过滤器
type OutboundFilter interface {
    Filter
    OnOutbound(ctx *RequestContext) error
}

// BaseFilter 基础过滤器
type BaseFilter struct {
    ID   string
    Name string
    Type FilterType
}

// GetID 获取过滤器ID
func (bf *BaseFilter) GetID() string {
    return bf.ID
}

// GetName 获取过滤器名称
func (bf *BaseFilter) GetName() string {
    return bf.Name
}

// GetType 获取过滤器类型
func (bf *BaseFilter) GetType() FilterType {
    return bf.Type
}

// LoggingFilter 日志过滤器
type LoggingFilter struct {
    BaseFilter
    logger *log.Logger
}

// NewLoggingFilter 创建日志过滤器
func NewLoggingFilter() *LoggingFilter {
    return &LoggingFilter{
        BaseFilter: BaseFilter{
            ID:   uuid.New().String(),
            Name: "logging",
            Type: FilterTypeBoth,
        },
        logger: log.New(os.Stdout, "[PROXY] ", log.LstdFlags),
    }
}

// OnInbound 入站处理
func (lf *LoggingFilter) OnInbound(ctx *RequestContext) error {
    lf.logger.Printf("Inbound request: %s from %s", ctx.ID, ctx.Conn.RemoteAddr())
    return nil
}

// OnOutbound 出站处理
func (lf *LoggingFilter) OnOutbound(ctx *RequestContext) error {
    lf.logger.Printf("Outbound response: %s, duration: %v", ctx.ID, time.Since(ctx.StartTime))
    return nil
}

// RateLimitFilter 限流过滤器
type RateLimitFilter struct {
    BaseFilter
    limiter *rate.Limiter
}

// NewRateLimitFilter 创建限流过滤器
func NewRateLimitFilter(rps int) *RateLimitFilter {
    return &RateLimitFilter{
        BaseFilter: BaseFilter{
            ID:   uuid.New().String(),
            Name: "rate_limit",
            Type: FilterTypeInbound,
        },
        limiter: rate.NewLimiter(rate.Limit(rps), rps),
    }
}

// OnInbound 入站处理
func (rlf *RateLimitFilter) OnInbound(ctx *RequestContext) error {
    if !rlf.limiter.Allow() {
        return fmt.Errorf("rate limit exceeded")
    }
    return nil
}

// CircuitBreakerFilter 熔断器过滤器
type CircuitBreakerFilter struct {
    BaseFilter
    breaker *CircuitBreaker
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
    config     *CircuitBreakerConfig
    state      CircuitBreakerState
    failures   int
    lastFailure time.Time
    mutex      sync.RWMutex
}

// CircuitBreakerState 熔断器状态
type CircuitBreakerState string

const (
    CircuitBreakerStateClosed   CircuitBreakerState = "CLOSED"
    CircuitBreakerStateOpen     CircuitBreakerState = "OPEN"
    CircuitBreakerStateHalfOpen CircuitBreakerState = "HALF_OPEN"
)

// NewCircuitBreakerFilter 创建熔断器过滤器
func NewCircuitBreakerFilter(config *CircuitBreakerConfig) *CircuitBreakerFilter {
    return &CircuitBreakerFilter{
        BaseFilter: BaseFilter{
            ID:   uuid.New().String(),
            Name: "circuit_breaker",
            Type: FilterTypeOutbound,
        },
        breaker: &CircuitBreaker{
            config: config,
            state:  CircuitBreakerStateClosed,
        },
    }
}

// OnOutbound 出站处理
func (cbf *CircuitBreakerFilter) OnOutbound(ctx *RequestContext) error {
    if !cbf.breaker.AllowRequest() {
        return fmt.Errorf("circuit breaker is open")
    }
    
    // 标记请求开始
    cbf.breaker.OnRequestStart()
    
    // 如果请求失败，记录失败
    if ctx.Error != nil {
        cbf.breaker.OnRequestFailure()
    } else {
        cbf.breaker.OnRequestSuccess()
    }
    
    return nil
}

// AllowRequest 是否允许请求
func (cb *CircuitBreaker) AllowRequest() bool {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()
    
    switch cb.state {
    case CircuitBreakerStateClosed:
        return true
    case CircuitBreakerStateOpen:
        if time.Since(cb.lastFailure) > cb.config.Timeout {
            cb.state = CircuitBreakerStateHalfOpen
            return true
        }
        return false
    case CircuitBreakerStateHalfOpen:
        return cb.failures < cb.config.HalfOpenLimit
    default:
        return false
    }
}

// OnRequestStart 请求开始
func (cb *CircuitBreaker) OnRequestStart() {
    // 可以在这里添加请求开始的处理逻辑
}

// OnRequestSuccess 请求成功
func (cb *CircuitBreaker) OnRequestSuccess() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    if cb.state == CircuitBreakerStateHalfOpen {
        cb.state = CircuitBreakerStateClosed
        cb.failures = 0
    }
}

// OnRequestFailure 请求失败
func (cb *CircuitBreaker) OnRequestFailure() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    cb.failures++
    cb.lastFailure = time.Now()
    
    if cb.state == CircuitBreakerStateClosed && cb.failures >= cb.config.Threshold {
        cb.state = CircuitBreakerStateOpen
    } else if cb.state == CircuitBreakerStateHalfOpen {
        cb.state = CircuitBreakerStateOpen
    }
}
```

### 1.2.3 服务发现实现

```go
// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
    // 注册服务
    Register(service *Service) error
    
    // 注销服务
    Deregister(serviceID string) error
    
    // 发现服务
    Discover(serviceName string) ([]*Endpoint, error)
    
    // 监听服务变化
    Watch(serviceName string) (<-chan ServiceEvent, error)
}

// Service 服务定义
type Service struct {
    ID      string            `json:"id"`
    Name    string            `json:"name"`
    Version string            `json:"version"`
    Address string            `json:"address"`
    Port    int               `json:"port"`
    Tags    []string          `json:"tags"`
    Metadata map[string]string `json:"metadata"`
    Health   HealthStatus     `json:"health"`
}

// Endpoint 服务端点
type Endpoint struct {
    ID      string            `json:"id"`
    Service string            `json:"service"`
    Address string            `json:"address"`
    Port    int               `json:"port"`
    Weight  int               `json:"weight"`
    Tags    []string          `json:"tags"`
    Metadata map[string]string `json:"metadata"`
}

// HealthStatus 健康状态
type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "HEALTHY"
    HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
    HealthStatusUnknown   HealthStatus = "UNKNOWN"
)

// ServiceEvent 服务事件
type ServiceEvent struct {
    Type     ServiceEventType `json:"type"`
    Service  *Service         `json:"service"`
    Timestamp time.Time       `json:"timestamp"`
}

// ServiceEventType 服务事件类型
type ServiceEventType string

const (
    ServiceEventAdded   ServiceEventType = "ADDED"
    ServiceEventUpdated ServiceEventType = "UPDATED"
    ServiceEventDeleted ServiceEventType = "DELETED"
)

// InMemoryServiceDiscovery 内存服务发现实现
type InMemoryServiceDiscovery struct {
    services map[string]*Service
    watchers map[string][]chan ServiceEvent
    mutex    sync.RWMutex
}

// NewServiceDiscovery 创建服务发现
func NewServiceDiscovery() ServiceDiscovery {
    return &InMemoryServiceDiscovery{
        services: make(map[string]*Service),
        watchers: make(map[string][]chan ServiceEvent),
    }
}

// Register 注册服务
func (sd *InMemoryServiceDiscovery) Register(service *Service) error {
    sd.mutex.Lock()
    defer sd.mutex.Unlock()
    
    // 检查服务是否已存在
    if existing, exists := sd.services[service.ID]; exists {
        // 更新服务
        sd.services[service.ID] = service
        sd.notifyWatchers(service.Name, ServiceEvent{
            Type:      ServiceEventUpdated,
            Service:   service,
            Timestamp: time.Now(),
        })
    } else {
        // 新增服务
        sd.services[service.ID] = service
        sd.notifyWatchers(service.Name, ServiceEvent{
            Type:      ServiceEventAdded,
            Service:   service,
            Timestamp: time.Now(),
        })
    }
    
    return nil
}

// Deregister 注销服务
func (sd *InMemoryServiceDiscovery) Deregister(serviceID string) error {
    sd.mutex.Lock()
    defer sd.mutex.Unlock()
    
    service, exists := sd.services[serviceID]
    if !exists {
        return fmt.Errorf("service not found: %s", serviceID)
    }
    
    delete(sd.services, serviceID)
    
    sd.notifyWatchers(service.Name, ServiceEvent{
        Type:      ServiceEventDeleted,
        Service:   service,
        Timestamp: time.Now(),
    })
    
    return nil
}

// Discover 发现服务
func (sd *InMemoryServiceDiscovery) Discover(serviceName string) ([]*Endpoint, error) {
    sd.mutex.RLock()
    defer sd.mutex.RUnlock()
    
    var endpoints []*Endpoint
    
    for _, service := range sd.services {
        if service.Name == serviceName && service.Health == HealthStatusHealthy {
            endpoint := &Endpoint{
                ID:      service.ID,
                Service: service.Name,
                Address: service.Address,
                Port:    service.Port,
                Weight:  1, // 默认权重
                Tags:    service.Tags,
                Metadata: service.Metadata,
            }
            endpoints = append(endpoints, endpoint)
        }
    }
    
    return endpoints, nil
}

// Watch 监听服务变化
func (sd *InMemoryServiceDiscovery) Watch(serviceName string) (<-chan ServiceEvent, error) {
    sd.mutex.Lock()
    defer sd.mutex.Unlock()
    
    eventChan := make(chan ServiceEvent, 10)
    sd.watchers[serviceName] = append(sd.watchers[serviceName], eventChan)
    
    return eventChan, nil
}

// notifyWatchers 通知监听器
func (sd *InMemoryServiceDiscovery) notifyWatchers(serviceName string, event ServiceEvent) {
    if watchers, exists := sd.watchers[serviceName]; exists {
        for _, watcher := range watchers {
            select {
            case watcher <- event:
            default:
                // 通道已满，跳过
            }
        }
    }
}
```

## 1.3 控制平面实现

### 1.3.1 配置管理

```go
// ControlPlane 控制平面接口
type ControlPlane interface {
    // 启动控制平面
    Start(ctx context.Context) error
    
    // 停止控制平面
    Stop(ctx context.Context) error
    
    // 配置路由规则
    ConfigureRouting(rules []RoutingRule) error
    
    // 配置策略
    ConfigurePolicy(policies []Policy) error
    
    // 获取代理状态
    GetProxyStatus() map[string]ProxyStatus
}

// RoutingRule 路由规则
type RoutingRule struct {
    ID          string                 `json:"id"`
    Service     string                 `json:"service"`
    Match       *MatchCondition        `json:"match"`
    Route       []RouteDestination     `json:"route"`
    Weight      map[string]int         `json:"weight"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// MatchCondition 匹配条件
type MatchCondition struct {
    Headers map[string]string `json:"headers"`
    Path    string            `json:"path"`
    Method  string            `json:"method"`
}

// RouteDestination 路由目标
type RouteDestination struct {
    Service string            `json:"service"`
    Version string            `json:"version"`
    Weight  int               `json:"weight"`
    Tags    map[string]string `json:"tags"`
}

// Policy 策略
type Policy struct {
    ID       string                 `json:"id"`
    Type     PolicyType             `json:"type"`
    Service  string                 `json:"service"`
    Config   map[string]interface{} `json:"config"`
    Enabled  bool                   `json:"enabled"`
}

// PolicyType 策略类型
type PolicyType string

const (
    PolicyTypeRateLimit    PolicyType = "RATE_LIMIT"
    PolicyTypeCircuitBreaker PolicyType = "CIRCUIT_BREAKER"
    PolicyTypeRetry        PolicyType = "RETRY"
    PolicyTypeTimeout      PolicyType = "TIMEOUT"
    PolicyTypeLoadBalancer PolicyType = "LOAD_BALANCER"
)

// ControlPlaneImpl 控制平面实现
type ControlPlaneImpl struct {
    discovery    ServiceDiscovery
    configStore  ConfigStore
    proxyManager ProxyManager
    apiServer    *APIServer
    stopChan     chan struct{}
}

// ConfigStore 配置存储接口
type ConfigStore interface {
    SaveRoutingRules(rules []RoutingRule) error
    LoadRoutingRules() ([]RoutingRule, error)
    SavePolicies(policies []Policy) error
    LoadPolicies() ([]Policy, error)
}

// ProxyManager 代理管理器
type ProxyManager interface {
    RegisterProxy(proxy Proxy) error
    UnregisterProxy(proxyID string) error
    GetProxies() map[string]Proxy
    UpdateProxyConfig(proxyID string, config interface{}) error
}

// NewControlPlane 创建控制平面
func NewControlPlane(discovery ServiceDiscovery, configStore ConfigStore) ControlPlane {
    return &ControlPlaneImpl{
        discovery:   discovery,
        configStore: configStore,
        proxyManager: NewProxyManager(),
        apiServer:    NewAPIServer(),
        stopChan:     make(chan struct{}),
    }
}

// Start 启动控制平面
func (cp *ControlPlaneImpl) Start(ctx context.Context) error {
    // 启动API服务器
    if err := cp.apiServer.Start(ctx); err != nil {
        return fmt.Errorf("failed to start API server: %w", err)
    }
    
    // 启动配置同步
    go cp.syncConfiguration(ctx)
    
    // 启动健康检查
    go cp.healthCheck(ctx)
    
    return nil
}

// syncConfiguration 同步配置
func (cp *ControlPlaneImpl) syncConfiguration(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-cp.stopChan:
            return
        case <-ticker.C:
            if err := cp.syncRoutingRules(); err != nil {
                log.Printf("Failed to sync routing rules: %v", err)
            }
            
            if err := cp.syncPolicies(); err != nil {
                log.Printf("Failed to sync policies: %v", err)
            }
        }
    }
}

// syncRoutingRules 同步路由规则
func (cp *ControlPlaneImpl) syncRoutingRules() error {
    rules, err := cp.configStore.LoadRoutingRules()
    if err != nil {
        return err
    }
    
    // 更新所有代理的路由规则
    proxies := cp.proxyManager.GetProxies()
    for _, proxy := range proxies {
        if serviceProxy, ok := proxy.(*ServiceProxy); ok {
            serviceProxy.router.UpdateRules(rules)
        }
    }
    
    return nil
}

// syncPolicies 同步策略
func (cp *ControlPlaneImpl) syncPolicies() error {
    policies, err := cp.configStore.LoadPolicies()
    if err != nil {
        return err
    }
    
    // 更新所有代理的策略
    proxies := cp.proxyManager.GetProxies()
    for _, proxy := range proxies {
        if serviceProxy, ok := proxy.(*ServiceProxy); ok {
            serviceProxy.UpdatePolicies(policies)
        }
    }
    
    return nil
}

// healthCheck 健康检查
func (cp *ControlPlaneImpl) healthCheck(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-cp.stopChan:
            return
        case <-ticker.C:
            cp.checkProxyHealth()
        }
    }
}

// checkProxyHealth 检查代理健康状态
func (cp *ControlPlaneImpl) checkProxyHealth() {
    proxies := cp.proxyManager.GetProxies()
    for proxyID, proxy := range proxies {
        status := proxy.GetStatus()
        
        // 检查代理是否响应
        if time.Since(status.Metrics.LastUpdateTime) > 5*time.Minute {
            log.Printf("Proxy %s appears to be unhealthy", proxyID)
            // 可以在这里实现代理重启逻辑
        }
    }
}
```

## 1.4 服务网格监控

### 1.4.1 指标收集

```go
// MetricsCollector 指标收集器
type MetricsCollector struct {
    metrics map[string]*ProxyMetrics
    mutex   sync.RWMutex
}

// ProxyMetrics 代理指标
type ProxyMetrics struct {
    ProxyID           string                 `json:"proxy_id"`
    ServiceName       string                 `json:"service_name"`
    RequestsTotal     int64                  `json:"requests_total"`
    RequestsSuccess   int64                  `json:"requests_success"`
    RequestsFailed    int64                  `json:"requests_failed"`
    AverageLatency    time.Duration          `json:"average_latency"`
    ErrorRate         float64                `json:"error_rate"`
    Throughput        float64                `json:"throughput"`
    LastUpdateTime    time.Time              `json:"last_update_time"`
    CustomMetrics     map[string]interface{} `json:"custom_metrics"`
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        metrics: make(map[string]*ProxyMetrics),
    }
}

// UpdateMetrics 更新指标
func (mc *MetricsCollector) UpdateMetrics(proxyID string, metrics *Metrics) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    if _, exists := mc.metrics[proxyID]; !exists {
        mc.metrics[proxyID] = &ProxyMetrics{
            ProxyID:       proxyID,
            CustomMetrics: make(map[string]interface{}),
        }
    }
    
    proxyMetrics := mc.metrics[proxyID]
    proxyMetrics.RequestsTotal = metrics.RequestsTotal
    proxyMetrics.RequestsSuccess = metrics.RequestsSuccess
    proxyMetrics.RequestsFailed = metrics.RequestsFailed
    proxyMetrics.AverageLatency = metrics.AverageLatency
    proxyMetrics.LastUpdateTime = metrics.LastUpdateTime
    
    // 计算错误率
    if metrics.RequestsTotal > 0 {
        proxyMetrics.ErrorRate = float64(metrics.RequestsFailed) / float64(metrics.RequestsTotal)
    }
    
    // 计算吞吐量（每秒请求数）
    if time.Since(proxyMetrics.LastUpdateTime) > 0 {
        proxyMetrics.Throughput = float64(metrics.RequestsTotal) / time.Since(proxyMetrics.LastUpdateTime).Seconds()
    }
}

// GetMetrics 获取指标
func (mc *MetricsCollector) GetMetrics() map[string]*ProxyMetrics {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    result := make(map[string]*ProxyMetrics)
    for k, v := range mc.metrics {
        result[k] = v
    }
    return result
}
```

### 1.4.2 分布式追踪

```go
// TraceCollector 追踪收集器
type TraceCollector struct {
    traces map[string]*Trace
    mutex  sync.RWMutex
}

// Trace 追踪信息
type Trace struct {
    ID        string        `json:"id"`
    Service   string        `json:"service"`
    StartTime time.Time     `json:"start_time"`
    EndTime   time.Time     `json:"end_time"`
    Duration  time.Duration `json:"duration"`
    Spans     []Span        `json:"spans"`
    Error     error         `json:"error,omitempty"`
}

// Span 追踪片段
type Span struct {
    ID       string            `json:"id"`
    ParentID string            `json:"parent_id"`
    Service  string            `json:"service"`
    Operation string           `json:"operation"`
    StartTime time.Time        `json:"start_time"`
    EndTime   time.Time        `json:"end_time"`
    Duration  time.Duration    `json:"duration"`
    Tags     map[string]string `json:"tags"`
    Logs     []Log             `json:"logs"`
}

// Log 日志条目
type Log struct {
    Timestamp time.Time                 `json:"timestamp"`
    Fields    map[string]interface{}    `json:"fields"`
}

// NewTraceCollector 创建追踪收集器
func NewTraceCollector() *TraceCollector {
    return &TraceCollector{
        traces: make(map[string]*Trace),
    }
}

// StartTrace 开始追踪
func (tc *TraceCollector) StartTrace(service string) *Trace {
    traceID := uuid.New().String()
    trace := &Trace{
        ID:        traceID,
        Service:   service,
        StartTime: time.Now(),
        Spans:     make([]Span, 0),
    }
    
    tc.mutex.Lock()
    tc.traces[traceID] = trace
    tc.mutex.Unlock()
    
    return trace
}

// EndTrace 结束追踪
func (tc *TraceCollector) EndTrace(traceID string, err error) error {
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    trace, exists := tc.traces[traceID]
    if !exists {
        return fmt.Errorf("trace not found: %s", traceID)
    }
    
    trace.EndTime = time.Now()
    trace.Duration = trace.EndTime.Sub(trace.StartTime)
    trace.Error = err
    
    return nil
}

// AddSpan 添加追踪片段
func (tc *TraceCollector) AddSpan(traceID string, span Span) error {
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    trace, exists := tc.traces[traceID]
    if !exists {
        return fmt.Errorf("trace not found: %s", traceID)
    }
    
    trace.Spans = append(trace.Spans, span)
    return nil
}
```

## 1.5 总结

服务网格模块涵盖了以下核心内容：

1. **理论基础**: 形式化定义服务网格的数学模型
2. **代理实现**: Go语言实现的服务代理和过滤器
3. **服务发现**: 动态服务注册和发现机制
4. **控制平面**: 配置管理和策略控制
5. **监控追踪**: 指标收集和分布式追踪

这个设计提供了一个完整的服务网格框架，支持微服务架构中的服务间通信、负载均衡、故障恢复等核心功能。
