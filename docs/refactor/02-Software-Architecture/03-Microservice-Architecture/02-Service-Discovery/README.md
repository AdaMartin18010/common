# 02-服务发现 (Service Discovery)

## 目录

- [02-服务发现 (Service Discovery)](#02-服务发现-service-discovery)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 设计原则](#11-设计原则)
    - [1.2 架构模式](#12-架构模式)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 服务发现系统](#21-服务发现系统)
    - [2.2 服务状态](#22-服务状态)
    - [2.3 一致性模型](#23-一致性模型)
  - [3. 服务注册](#3-服务注册)
    - [3.1 注册接口](#31-注册接口)
    - [3.2 注册实现](#32-注册实现)
  - [4. 服务发现](#4-服务发现)
    - [4.1 发现接口](#41-发现接口)
    - [4.2 服务缓存](#42-服务缓存)
  - [5. 健康检查](#5-健康检查)
    - [5.1 健康检查接口](#51-健康检查接口)
  - [6. 负载均衡](#6-负载均衡)
    - [6.1 负载均衡接口](#61-负载均衡接口)
  - [7. 故障转移](#7-故障转移)
    - [7.1 故障转移策略](#71-故障转移策略)
  - [8. Go语言实现](#8-go语言实现)
    - [8.1 服务发现管理器](#81-服务发现管理器)
    - [8.2 使用示例](#82-使用示例)
  - [9. 应用场景](#9-应用场景)
    - [9.1 微服务架构](#91-微服务架构)
    - [9.2 容器编排](#92-容器编排)
  - [10. 总结](#10-总结)
    - [10.1 核心特性](#101-核心特性)
    - [10.2 技术优势](#102-技术优势)
    - [10.3 应用价值](#103-应用价值)

---

## 1. 概述

服务发现是微服务架构中的核心组件，负责服务的注册、发现、健康检查和负载均衡。它解决了分布式系统中服务定位和通信的问题。

### 1.1 设计原则

1. **高可用性**: 服务发现系统本身必须高可用
2. **一致性**: 确保服务信息的一致性
3. **实时性**: 及时反映服务状态变化
4. **可扩展性**: 支持大规模服务注册
5. **容错性**: 具备故障检测和恢复能力

### 1.2 架构模式

```text
服务发现架构
├── 服务注册中心 (Service Registry)
│   ├── 服务注册
│   ├── 服务注销
│   ├── 健康检查
│   └── 元数据管理
├── 服务发现客户端 (Service Discovery Client)
│   ├── 服务查询
│   ├── 负载均衡
│   ├── 故障转移
│   └── 缓存管理
└── 服务提供者 (Service Provider)
    ├── 服务注册
    ├── 心跳检测
    ├── 优雅下线
    └── 元数据更新
```

## 2. 形式化定义

### 2.1 服务发现系统

**定义 2.1 (服务发现系统)**: 服务发现系统是一个五元组 Σ = (S, R, D, H, L)，其中：

- S是服务集合
- R是注册表
- D是发现器
- H是健康检查器
- L是负载均衡器

**定义 2.2 (服务)**: 服务是一个六元组 s = (id, name, address, port, metadata, status)，其中：

- id是服务唯一标识符
- name是服务名称
- address是服务地址
- port是服务端口
- metadata是服务元数据
- status是服务状态

### 2.2 服务状态

**定义 2.3 (服务状态)**: 服务状态集合 Status = {UP, DOWN, STARTING, OUT_OF_SERVICE}

**定义 2.4 (状态转换)**: 状态转换函数 δ: Status × Event → Status

```math
δ(UP, health_check_fail) = DOWN
δ(DOWN, health_check_pass) = UP
δ(STARTING, ready) = UP
δ(UP, graceful_shutdown) = OUT_OF_SERVICE
```

### 2.3 一致性模型

**定义 2.5 (最终一致性)**: 对于任意两个节点 n₁, n₂，存在时间 t，使得：

```latex
∀s ∈ S: n₁.getService(s) = n₂.getService(s) after time t
```

**定理 2.1 (CAP定理)**: 在分布式服务发现系统中，最多只能同时满足一致性(Consistency)、可用性(Availability)和分区容错性(Partition tolerance)中的两个。

## 3. 服务注册

### 3.1 注册接口

```go
// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
    Register(service *Service) error
    Deregister(serviceID string) error
    Update(service *Service) error
    GetService(serviceID string) (*Service, error)
    ListServices() ([]*Service, error)
    Watch(serviceID string) (<-chan *ServiceEvent, error)
}

// Service 服务定义
type Service struct {
    ID       string            `json:"id"`
    Name     string            `json:"name"`
    Address  string            `json:"address"`
    Port     int               `json:"port"`
    Metadata map[string]string `json:"metadata"`
    Status   ServiceStatus     `json:"status"`
    Created  time.Time         `json:"created"`
    Updated  time.Time         `json:"updated"`
}

type ServiceStatus int

const (
    StatusStarting ServiceStatus = iota
    StatusUp
    StatusDown
    StatusOutOfService
)

// ServiceEvent 服务事件
type ServiceEvent struct {
    Type    ServiceEventType
    Service *Service
    Time    time.Time
}

type ServiceEventType int

const (
    EventRegistered ServiceEventType = iota
    EventDeregistered
    EventUpdated
    EventStatusChanged
)
```

### 3.2 注册实现

```go
// EtcdRegistry 基于Etcd的服务注册实现
type EtcdRegistry struct {
    client    *clientv3.Client
    prefix    string
    ttl       int64
    mutex     sync.RWMutex
    watchers  map[string][]chan *ServiceEvent
}

func NewEtcdRegistry(endpoints []string, prefix string) (*EtcdRegistry, error) {
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, err
    }

    return &EtcdRegistry{
        client:   client,
        prefix:   prefix,
        ttl:      30, // 30秒TTL
        watchers: make(map[string][]chan *ServiceEvent),
    }, nil
}

func (r *EtcdRegistry) Register(service *Service) error {
    service.Updated = time.Now()
    
    // 序列化服务信息
    data, err := json.Marshal(service)
    if err != nil {
        return err
    }

    // 创建租约
    lease, err := r.client.Grant(context.Background(), r.ttl)
    if err != nil {
        return err
    }

    // 注册服务
    key := r.getServiceKey(service.ID)
    _, err = r.client.Put(context.Background(), key, string(data), clientv3.WithLease(lease.ID))
    if err != nil {
        return err
    }

    // 保持租约活跃
    go r.keepAlive(lease.ID, key, data)

    // 通知观察者
    r.notifyWatchers(EventRegistered, service)

    return nil
}

func (r *EtcdRegistry) Deregister(serviceID string) error {
    key := r.getServiceKey(serviceID)
    
    // 获取服务信息用于通知
    service, err := r.GetService(serviceID)
    if err != nil {
        return err
    }

    // 删除服务
    _, err = r.client.Delete(context.Background(), key)
    if err != nil {
        return err
    }

    // 通知观察者
    r.notifyWatchers(EventDeregistered, service)

    return nil
}

func (r *EtcdRegistry) keepAlive(leaseID clientv3.LeaseID, key string, data []byte) {
    ticker := time.NewTicker(time.Duration(r.ttl/3) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            _, err := r.client.KeepAliveOnce(context.Background(), leaseID)
            if err != nil {
                // 租约过期，重新注册
                r.reRegister(key, data)
                return
            }
        }
    }
}

func (r *EtcdRegistry) reRegister(key string, data []byte) {
    // 重新创建租约并注册
    lease, err := r.client.Grant(context.Background(), r.ttl)
    if err != nil {
        return
    }

    _, err = r.client.Put(context.Background(), key, string(data), clientv3.WithLease(lease.ID))
    if err != nil {
        return
    }

    go r.keepAlive(lease.ID, key, data)
}

func (r *EtcdRegistry) getServiceKey(serviceID string) string {
    return fmt.Sprintf("%s/services/%s", r.prefix, serviceID)
}
```

## 4. 服务发现

### 4.1 发现接口

```go
// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
    Discover(serviceName string) ([]*Service, error)
    DiscoverOne(serviceName string) (*Service, error)
    Watch(serviceName string) (<-chan *ServiceEvent, error)
    Close() error
}

// DiscoveryClient 服务发现客户端
type DiscoveryClient struct {
    registry ServiceRegistry
    cache    *ServiceCache
    balancer LoadBalancer
    mutex    sync.RWMutex
}

func NewDiscoveryClient(registry ServiceRegistry) *DiscoveryClient {
    return &DiscoveryClient{
        registry: registry,
        cache:    NewServiceCache(),
        balancer: NewRoundRobinBalancer(),
    }
}

func (dc *DiscoveryClient) Discover(serviceName string) ([]*Service, error) {
    // 先从缓存获取
    services := dc.cache.Get(serviceName)
    if len(services) > 0 {
        return services, nil
    }

    // 从注册中心获取
    allServices, err := dc.registry.ListServices()
    if err != nil {
        return nil, err
    }

    // 过滤指定服务
    var result []*Service
    for _, service := range allServices {
        if service.Name == serviceName && service.Status == StatusUp {
            result = append(result, service)
        }
    }

    // 更新缓存
    dc.cache.Set(serviceName, result)

    return result, nil
}

func (dc *DiscoveryClient) DiscoverOne(serviceName string) (*Service, error) {
    services, err := dc.Discover(serviceName)
    if err != nil {
        return nil, err
    }

    if len(services) == 0 {
        return nil, fmt.Errorf("no service found: %s", serviceName)
    }

    // 使用负载均衡器选择服务
    return dc.balancer.Select(services), nil
}
```

### 4.2 服务缓存

```go
// ServiceCache 服务缓存
type ServiceCache struct {
    cache map[string][]*Service
    ttl   time.Duration
    mutex sync.RWMutex
}

func NewServiceCache() *ServiceCache {
    return &ServiceCache{
        cache: make(map[string][]*Service),
        ttl:   30 * time.Second,
    }
}

func (sc *ServiceCache) Get(serviceName string) []*Service {
    sc.mutex.RLock()
    defer sc.mutex.RUnlock()
    
    return sc.cache[serviceName]
}

func (sc *ServiceCache) Set(serviceName string, services []*Service) {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    
    sc.cache[serviceName] = services
    
    // 设置过期时间
    go func() {
        time.Sleep(sc.ttl)
        sc.mutex.Lock()
        delete(sc.cache, serviceName)
        sc.mutex.Unlock()
    }()
}

func (sc *ServiceCache) Invalidate(serviceName string) {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    
    delete(sc.cache, serviceName)
}
```

## 5. 健康检查

### 5.1 健康检查接口

```go
// HealthChecker 健康检查接口
type HealthChecker interface {
    Check(service *Service) error
    Start(service *Service) error
    Stop(serviceID string) error
}

// HealthCheck 健康检查定义
type HealthCheck struct {
    Type     string            `json:"type"`
    URL      string            `json:"url"`
    Interval time.Duration     `json:"interval"`
    Timeout  time.Duration     `json:"timeout"`
    Headers  map[string]string `json:"headers"`
}

// HTTPHealthChecker HTTP健康检查实现
type HTTPHealthChecker struct {
    client  *http.Client
    checks  map[string]*HealthCheck
    running map[string]chan bool
    mutex   sync.RWMutex
}

func NewHTTPHealthChecker() *HTTPHealthChecker {
    return &HTTPHealthChecker{
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
        checks:  make(map[string]*HealthCheck),
        running: make(map[string]chan bool),
    }
}

func (hc *HTTPHealthChecker) Check(service *Service) error {
    check, exists := hc.checks[service.ID]
    if !exists {
        return nil // 没有配置健康检查
    }

    // 构建健康检查URL
    url := fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, check.URL)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    // 添加自定义头部
    for key, value := range check.Headers {
        req.Header.Set(key, value)
    }

    resp, err := hc.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return nil
    }

    return fmt.Errorf("health check failed: status %d", resp.StatusCode)
}

func (hc *HTTPHealthChecker) Start(service *Service) error {
    hc.mutex.Lock()
    defer hc.mutex.Unlock()

    if _, exists := hc.running[service.ID]; exists {
        return fmt.Errorf("health check already running for service: %s", service.ID)
    }

    check, exists := hc.checks[service.ID]
    if !exists {
        return nil // 没有配置健康检查
    }

    stopCh := make(chan bool)
    hc.running[service.ID] = stopCh

    go func() {
        ticker := time.NewTicker(check.Interval)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                if err := hc.Check(service); err != nil {
                    // 健康检查失败，更新服务状态
                    hc.updateServiceStatus(service.ID, StatusDown)
                } else {
                    // 健康检查成功，更新服务状态
                    hc.updateServiceStatus(service.ID, StatusUp)
                }
            case <-stopCh:
                return
            }
        }
    }()

    return nil
}

func (hc *HTTPHealthChecker) Stop(serviceID string) error {
    hc.mutex.Lock()
    defer hc.mutex.Unlock()

    if stopCh, exists := hc.running[serviceID]; exists {
        close(stopCh)
        delete(hc.running, serviceID)
    }

    return nil
}

func (hc *HTTPHealthChecker) updateServiceStatus(serviceID string, status ServiceStatus) {
    // 这里应该调用注册中心更新服务状态
    // 实际实现中需要注入注册中心依赖
}
```

## 6. 负载均衡

### 6.1 负载均衡接口

```go
// LoadBalancer 负载均衡接口
type LoadBalancer interface {
    Select(services []*Service) *Service
    Update(services []*Service)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
    services []*Service
    index    int64
    mutex    sync.RWMutex
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
    return &RoundRobinBalancer{
        services: make([]*Service, 0),
        index:    0,
    }
}

func (rr *RoundRobinBalancer) Select(services []*Service) *Service {
    rr.mutex.Lock()
    defer rr.mutex.Unlock()

    if len(services) == 0 {
        return nil
    }

    // 原子递增索引
    current := atomic.AddInt64(&rr.index, 1)
    index := int(current) % len(services)

    return services[index]
}

func (rr *RoundRobinBalancer) Update(services []*Service) {
    rr.mutex.Lock()
    defer rr.mutex.Unlock()

    rr.services = services
}

// WeightedRoundRobinBalancer 加权轮询负载均衡器
type WeightedRoundRobinBalancer struct {
    services []*WeightedService
    mutex    sync.RWMutex
}

type WeightedService struct {
    Service  *Service
    Weight   int
    Current  int
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
    return &WeightedRoundRobinBalancer{
        services: make([]*WeightedService, 0),
    }
}

func (wrr *WeightedRoundRobinBalancer) Select(services []*Service) *Service {
    wrr.mutex.Lock()
    defer wrr.mutex.Unlock()

    if len(wrr.services) == 0 {
        return nil
    }

    // 找到当前权重最大的服务
    maxWeight := -1
    selectedIndex := -1

    for i, ws := range wrr.services {
        if ws.Current > maxWeight {
            maxWeight = ws.Current
            selectedIndex = i
        }
    }

    if selectedIndex == -1 {
        return nil
    }

    // 更新权重
    wrr.services[selectedIndex].Current -= wrr.getTotalWeight()

    // 增加所有服务的当前权重
    for _, ws := range wrr.services {
        ws.Current += ws.Weight
    }

    return wrr.services[selectedIndex].Service
}

func (wrr *WeightedRoundRobinBalancer) getTotalWeight() int {
    total := 0
    for _, ws := range wrr.services {
        total += ws.Weight
    }
    return total
}

func (wrr *WeightedRoundRobinBalancer) Update(services []*Service) {
    wrr.mutex.Lock()
    defer wrr.mutex.Unlock()

    // 转换为加权服务
    weightedServices := make([]*WeightedService, len(services))
    for i, service := range services {
        weight := 1 // 默认权重为1
        if w, ok := service.Metadata["weight"]; ok {
            if parsed, err := strconv.Atoi(w); err == nil {
                weight = parsed
            }
        }

        weightedServices[i] = &WeightedService{
            Service: service,
            Weight:  weight,
            Current: weight,
        }
    }

    wrr.services = weightedServices
}
```

## 7. 故障转移

### 7.1 故障转移策略

```go
// FailoverStrategy 故障转移策略
type FailoverStrategy interface {
    Select(services []*Service, failed *Service) *Service
}

// SimpleFailoverStrategy 简单故障转移策略
type SimpleFailoverStrategy struct{}

func (sfs *SimpleFailoverStrategy) Select(services []*Service, failed *Service) *Service {
    for _, service := range services {
        if service.ID != failed.ID && service.Status == StatusUp {
            return service
        }
    }
    return nil
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
    serviceID    string
    state        CircuitBreakerState
    failureCount int
    threshold    int
    timeout      time.Duration
    lastFailure  time.Time
    mutex        sync.RWMutex
}

type CircuitBreakerState int

const (
    StateClosed CircuitBreakerState = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(serviceID string, threshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        serviceID: serviceID,
        state:     StateClosed,
        threshold: threshold,
        timeout:   timeout,
    }
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
    if !cb.canExecute() {
        return fmt.Errorf("circuit breaker is open")
    }

    err := operation()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mutex.RLock()
    defer cb.mutex.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.mutex.RUnlock()
            cb.mutex.Lock()
            cb.state = StateHalfOpen
            cb.mutex.Unlock()
            cb.mutex.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err != nil {
        cb.failureCount++
        cb.lastFailure = time.Now()

        if cb.failureCount >= cb.threshold {
            cb.state = StateOpen
        }
    } else {
        cb.failureCount = 0
        cb.state = StateClosed
    }
}
```

## 8. Go语言实现

### 8.1 服务发现管理器

```go
// ServiceDiscoveryManager 服务发现管理器
type ServiceDiscoveryManager struct {
    registry     ServiceRegistry
    discovery    ServiceDiscovery
    healthChecker HealthChecker
    balancer     LoadBalancer
    failover     FailoverStrategy
    mutex        sync.RWMutex
}

func NewServiceDiscoveryManager(
    registry ServiceRegistry,
    healthChecker HealthChecker,
    balancer LoadBalancer,
    failover FailoverStrategy,
) *ServiceDiscoveryManager {
    discovery := NewDiscoveryClient(registry)
    
    return &ServiceDiscoveryManager{
        registry:      registry,
        discovery:     discovery,
        healthChecker: healthChecker,
        balancer:      balancer,
        failover:      failover,
    }
}

func (sdm *ServiceDiscoveryManager) RegisterService(service *Service) error {
    // 注册服务
    if err := sdm.registry.Register(service); err != nil {
        return err
    }

    // 启动健康检查
    if err := sdm.healthChecker.Start(service); err != nil {
        return err
    }

    return nil
}

func (sdm *ServiceDiscoveryManager) DeregisterService(serviceID string) error {
    // 停止健康检查
    sdm.healthChecker.Stop(serviceID)

    // 注销服务
    return sdm.registry.Deregister(serviceID)
}

func (sdm *ServiceDiscoveryManager) DiscoverService(serviceName string) (*Service, error) {
    return sdm.discovery.DiscoverOne(serviceName)
}

func (sdm *ServiceDiscoveryManager) DiscoverServiceWithFailover(serviceName string) (*Service, error) {
    services, err := sdm.discovery.Discover(serviceName)
    if err != nil {
        return nil, err
    }

    if len(services) == 0 {
        return nil, fmt.Errorf("no service found: %s", serviceName)
    }

    // 使用负载均衡器选择服务
    selected := sdm.balancer.Select(services)
    
    // 如果服务不可用，使用故障转移
    if selected.Status != StatusUp {
        selected = sdm.failover.Select(services, selected)
        if selected == nil {
            return nil, fmt.Errorf("no available service: %s", serviceName)
        }
    }

    return selected, nil
}
```

### 8.2 使用示例

```go
// 使用示例
func main() {
    // 创建Etcd注册中心
    registry, err := NewEtcdRegistry([]string{"localhost:2379"}, "/myapp")
    if err != nil {
        log.Fatal(err)
    }

    // 创建健康检查器
    healthChecker := NewHTTPHealthChecker()

    // 创建负载均衡器
    balancer := NewRoundRobinBalancer()

    // 创建故障转移策略
    failover := &SimpleFailoverStrategy{}

    // 创建服务发现管理器
    sdm := NewServiceDiscoveryManager(registry, healthChecker, balancer, failover)

    // 注册服务
    service := &Service{
        ID:      "user-service-1",
        Name:    "user-service",
        Address: "localhost",
        Port:    8080,
        Status:  StatusUp,
        Metadata: map[string]string{
            "version": "1.0.0",
            "weight":  "2",
        },
    }

    if err := sdm.RegisterService(service); err != nil {
        log.Fatal(err)
    }

    // 发现服务
    discovered, err := sdm.DiscoverService("user-service")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Discovered service: %s:%d\n", discovered.Address, discovered.Port)

    // 使用熔断器调用服务
    cb := NewCircuitBreaker(discovered.ID, 5, 30*time.Second)
    
    err = cb.Execute(func() error {
        // 调用服务
        return callService(discovered)
    })
    
    if err != nil {
        log.Printf("Service call failed: %v", err)
    }
}

func callService(service *Service) error {
    // 实际的服务调用逻辑
    url := fmt.Sprintf("http://%s:%d/health", service.Address, service.Port)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("service returned status: %d", resp.StatusCode)
    }
    
    return nil
}
```

## 9. 应用场景

### 9.1 微服务架构

```go
// 微服务应用示例
type MicroserviceApp struct {
    discovery *ServiceDiscoveryManager
    services  map[string]*Service
}

func NewMicroserviceApp() *MicroserviceApp {
    // 初始化服务发现
    registry, _ := NewEtcdRegistry([]string{"localhost:2379"}, "/myapp")
    healthChecker := NewHTTPHealthChecker()
    balancer := NewRoundRobinBalancer()
    failover := &SimpleFailoverStrategy{}
    
    discovery := NewServiceDiscoveryManager(registry, healthChecker, balancer, failover)
    
    return &MicroserviceApp{
        discovery: discovery,
        services:  make(map[string]*Service),
    }
}

func (app *MicroserviceApp) StartService(name, address string, port int) error {
    service := &Service{
        ID:      fmt.Sprintf("%s-%d", name, time.Now().Unix()),
        Name:    name,
        Address: address,
        Port:    port,
        Status:  StatusUp,
    }
    
    if err := app.discovery.RegisterService(service); err != nil {
        return err
    }
    
    app.services[service.ID] = service
    return nil
}

func (app *MicroserviceApp) CallService(serviceName string) error {
    service, err := app.discovery.DiscoverServiceWithFailover(serviceName)
    if err != nil {
        return err
    }
    
    // 调用服务
    return app.callService(service)
}

func (app *MicroserviceApp) callService(service *Service) error {
    // 实际的RPC调用逻辑
    return nil
}
```

### 9.2 容器编排

```go
// Kubernetes服务发现适配器
type KubernetesServiceDiscovery struct {
    clientset *kubernetes.Clientset
    namespace string
}

func NewKubernetesServiceDiscovery(namespace string) (*KubernetesServiceDiscovery, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }
    
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }
    
    return &KubernetesServiceDiscovery{
        clientset: clientset,
        namespace: namespace,
    }, nil
}

func (ksd *KubernetesServiceDiscovery) Discover(serviceName string) ([]*Service, error) {
    services, err := ksd.clientset.CoreV1().Services(ksd.namespace).List(context.Background(), metav1.ListOptions{
        LabelSelector: fmt.Sprintf("app=%s", serviceName),
    })
    if err != nil {
        return nil, err
    }
    
    var result []*Service
    for _, svc := range services.Items {
        service := &Service{
            ID:      string(svc.UID),
            Name:    svc.Name,
            Address: svc.Spec.ClusterIP,
            Port:    int(svc.Spec.Ports[0].Port),
            Status:  StatusUp,
        }
        result = append(result, service)
    }
    
    return result, nil
}
```

## 10. 总结

服务发现是微服务架构中的关键组件，通过Go语言的实现，提供了高性能、高可用的服务注册和发现能力。

### 10.1 核心特性

1. **服务注册**: 支持服务的注册、注销和更新
2. **服务发现**: 提供高效的服务查询和负载均衡
3. **健康检查**: 实时监控服务健康状态
4. **故障转移**: 自动处理服务故障和恢复
5. **高可用性**: 支持集群部署和容错

### 10.2 技术优势

1. **高性能**: 基于Etcd的高性能存储
2. **一致性**: 强一致性保证
3. **可扩展性**: 支持大规模服务注册
4. **容错性**: 内置故障检测和恢复机制
5. **易用性**: 简洁的API设计

### 10.3 应用价值

1. **微服务架构**: 支持大规模微服务部署
2. **容器编排**: 与Kubernetes等平台集成
3. **云原生应用**: 适配云环境部署
4. **分布式系统**: 解决服务定位问题
5. **高可用架构**: 提升系统可用性

服务发现通过标准化的接口和高效的实现，为现代分布式系统提供了可靠的服务治理能力。
