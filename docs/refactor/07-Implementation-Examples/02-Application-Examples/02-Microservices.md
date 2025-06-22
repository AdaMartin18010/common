# 02-微服务 (Microservices)

## 1. 概述

### 1.1 微服务基础

**微服务** 是一种将应用程序构建为一组小型自治服务的架构风格，每个服务运行在自己的进程中，通过轻量级机制通信。

**核心原则**：

- **单一职责**：每个服务专注于一个业务功能
- **自治性**：服务独立部署和扩展
- **松耦合**：服务间通过API通信
- **技术多样性**：不同服务可使用不同技术栈

### 1.2 微服务架构模式

**形式化定义**：
设 ```latex
$S = \{s_1, s_2, ..., s_n\}$
``` 为服务集合，微服务系统：
$```latex
$\text{MicroserviceSystem}(S) = \bigcup_{i=1}^{n} \text{Service}(s_i) \cup \text{Communication}(S)$
```$

## 2. 服务发现

### 2.1 理论基础

**服务发现** 是微服务架构中的关键组件，负责动态注册和发现服务实例。

**功能**：

- 服务注册：新服务实例注册到注册中心
- 服务发现：客户端查询可用的服务实例
- 健康检查：监控服务实例的健康状态
- 负载均衡：在多个实例间分发请求

### 2.2 Go语言实现

```go
package microservices

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "sync"
 "time"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
 ID       string            `json:"id"`
 Name     string            `json:"name"`
 Host     string            `json:"host"`
 Port     int               `json:"port"`
 Metadata map[string]string `json:"metadata"`
 Status   string            `json:"status"`
 LastSeen time.Time         `json:"last_seen"`
}

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
 instances map[string]*ServiceInstance
 mu        sync.RWMutex
}

// NewServiceRegistry 创建服务注册中心
func NewServiceRegistry() *ServiceRegistry {
 return &ServiceRegistry{
  instances: make(map[string]*ServiceInstance),
 }
}

// Register 注册服务实例
func (sr *ServiceRegistry) Register(instance *ServiceInstance) error {
 sr.mu.Lock()
 defer sr.mu.Unlock()
 
 instance.LastSeen = time.Now()
 instance.Status = "healthy"
 sr.instances[instance.ID] = instance
 
 log.Printf("服务实例注册: %s (%s:%d)", instance.Name, instance.Host, instance.Port)
 return nil
}

// Deregister 注销服务实例
func (sr *ServiceRegistry) Deregister(instanceID string) error {
 sr.mu.Lock()
 defer sr.mu.Unlock()
 
 if instance, exists := sr.instances[instanceID]; exists {
  delete(sr.instances, instanceID)
  log.Printf("服务实例注销: %s (%s:%d)", instance.Name, instance.Host, instance.Port)
 }
 
 return nil
}

// GetInstances 获取服务实例
func (sr *ServiceRegistry) GetInstances(serviceName string) []*ServiceInstance {
 sr.mu.RLock()
 defer sr.mu.RUnlock()
 
 var instances []*ServiceInstance
 for _, instance := range sr.instances {
  if instance.Name == serviceName && instance.Status == "healthy" {
   instances = append(instances, instance)
  }
 }
 
 return instances
}

// HealthCheck 健康检查
func (sr *ServiceRegistry) HealthCheck() {
 ticker := time.NewTicker(30 * time.Second)
 defer ticker.Stop()
 
 for range ticker.C {
  sr.mu.Lock()
  for id, instance := range sr.instances {
   if time.Since(instance.LastSeen) > 60*time.Second {
    instance.Status = "unhealthy"
    log.Printf("服务实例不健康: %s", id)
   }
  }
  sr.mu.Unlock()
 }
}

// RegistryServer 注册中心服务器
type RegistryServer struct {
 registry *ServiceRegistry
}

// NewRegistryServer 创建注册中心服务器
func NewRegistryServer() *RegistryServer {
 registry := NewServiceRegistry()
 go registry.HealthCheck()
 
 return &RegistryServer{
  registry: registry,
 }
}

// Start 启动注册中心服务器
func (rs *RegistryServer) Start(port int) error {
 http.HandleFunc("/register", rs.handleRegister)
 http.HandleFunc("/deregister", rs.handleDeregister)
 http.HandleFunc("/discover", rs.handleDiscover)
 http.HandleFunc("/health", rs.handleHealth)
 
 log.Printf("注册中心服务器启动在端口 %d", port)
 return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// handleRegister 处理注册请求
func (rs *RegistryServer) handleRegister(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost {
  http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
  return
 }
 
 var instance ServiceInstance
 if err := json.NewDecoder(r.Body).Decode(&instance); err != nil {
  http.Error(w, "无效的请求体", http.StatusBadRequest)
  return
 }
 
 if err := rs.registry.Register(&instance); err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 
 w.WriteHeader(http.StatusOK)
}

// handleDeregister 处理注销请求
func (rs *RegistryServer) handleDeregister(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost {
  http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
  return
 }
 
 var req struct {
  InstanceID string `json:"instance_id"`
 }
 
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
  http.Error(w, "无效的请求体", http.StatusBadRequest)
  return
 }
 
 if err := rs.registry.Deregister(req.InstanceID); err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 
 w.WriteHeader(http.StatusOK)
}

// handleDiscover 处理服务发现请求
func (rs *RegistryServer) handleDiscover(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodGet {
  http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
  return
 }
 
 serviceName := r.URL.Query().Get("service")
 if serviceName == "" {
  http.Error(w, "缺少服务名称参数", http.StatusBadRequest)
  return
 }
 
 instances := rs.registry.GetInstances(serviceName)
 
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(instances)
}

// handleHealth 处理健康检查请求
func (rs *RegistryServer) handleHealth(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]string{
  "status": "healthy",
 })
}
```

### 2.3 服务发现客户端

```go
// ServiceDiscoveryClient 服务发现客户端
type ServiceDiscoveryClient struct {
 registryURL string
 client      *http.Client
}

// NewServiceDiscoveryClient 创建服务发现客户端
func NewServiceDiscoveryClient(registryURL string) *ServiceDiscoveryClient {
 return &ServiceDiscoveryClient{
  registryURL: registryURL,
  client: &http.Client{
   Timeout: 10 * time.Second,
  },
 }
}

// Discover 发现服务实例
func (sdc *ServiceDiscoveryClient) Discover(serviceName string) ([]*ServiceInstance, error) {
 url := fmt.Sprintf("%s/discover?service=%s", sdc.registryURL, serviceName)
 
 resp, err := sdc.client.Get(url)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()
 
 if resp.StatusCode != http.StatusOK {
  return nil, fmt.Errorf("服务发现失败: %d", resp.StatusCode)
 }
 
 var instances []*ServiceInstance
 if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
  return nil, err
 }
 
 return instances, nil
}

// Register 注册服务实例
func (sdc *ServiceDiscoveryClient) Register(instance *ServiceInstance) error {
 data, err := json.Marshal(instance)
 if err != nil {
  return err
 }
 
 resp, err := sdc.client.Post(sdc.registryURL+"/register", "application/json", bytes.NewReader(data))
 if err != nil {
  return err
 }
 defer resp.Body.Close()
 
 if resp.StatusCode != http.StatusOK {
  return fmt.Errorf("服务注册失败: %d", resp.StatusCode)
 }
 
 return nil
}

// Deregister 注销服务实例
func (sdc *ServiceDiscoveryClient) Deregister(instanceID string) error {
 req := struct {
  InstanceID string `json:"instance_id"`
 }{
  InstanceID: instanceID,
 }
 
 data, err := json.Marshal(req)
 if err != nil {
  return err
 }
 
 resp, err := sdc.client.Post(sdc.registryURL+"/deregister", "application/json", bytes.NewReader(data))
 if err != nil {
  return err
 }
 defer resp.Body.Close()
 
 if resp.StatusCode != http.StatusOK {
  return fmt.Errorf("服务注销失败: %d", resp.StatusCode)
 }
 
 return nil
}
```

## 3. 负载均衡

### 3.1 理论基础

**负载均衡** 是将请求分发到多个服务实例的技术，确保系统的高可用性和性能。

**算法类型**：

- **轮询 (Round Robin)**：依次分发请求
- **加权轮询 (Weighted Round Robin)**：根据权重分发
- **最少连接 (Least Connections)**：选择连接数最少的实例
- **随机 (Random)**：随机选择实例

### 3.2 Go语言实现

```go
// LoadBalancer 负载均衡器
type LoadBalancer struct {
 instances []*ServiceInstance
 strategy  LoadBalancingStrategy
 mu        sync.RWMutex
}

// LoadBalancingStrategy 负载均衡策略接口
type LoadBalancingStrategy interface {
 Select(instances []*ServiceInstance) *ServiceInstance
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer(strategy LoadBalancingStrategy) *LoadBalancer {
 return &LoadBalancer{
  instances: make([]*ServiceInstance, 0),
  strategy:  strategy,
 }
}

// UpdateInstances 更新实例列表
func (lb *LoadBalancer) UpdateInstances(instances []*ServiceInstance) {
 lb.mu.Lock()
 defer lb.mu.Unlock()
 
 lb.instances = instances
}

// GetInstance 获取实例
func (lb *LoadBalancer) GetInstance() *ServiceInstance {
 lb.mu.RLock()
 defer lb.mu.RUnlock()
 
 if len(lb.instances) == 0 {
  return nil
 }
 
 return lb.strategy.Select(lb.instances)
}

// RoundRobinStrategy 轮询策略
type RoundRobinStrategy struct {
 current int
 mu      sync.Mutex
}

// Select 选择实例
func (rr *RoundRobinStrategy) Select(instances []*ServiceInstance) *ServiceInstance {
 rr.mu.Lock()
 defer rr.mu.Unlock()
 
 if len(instances) == 0 {
  return nil
 }
 
 instance := instances[rr.current]
 rr.current = (rr.current + 1) % len(instances)
 
 return instance
}

// WeightedRoundRobinStrategy 加权轮询策略
type WeightedRoundRobinStrategy struct {
 current int
 mu      sync.Mutex
}

// Select 选择实例
func (wrr *WeightedRoundRobinStrategy) Select(instances []*ServiceInstance) *ServiceInstance {
 wrr.mu.Lock()
 defer wrr.mu.Unlock()
 
 if len(instances) == 0 {
  return nil
 }
 
 // 简化实现，实际应该考虑权重
 instance := instances[wrr.current]
 wrr.current = (wrr.current + 1) % len(instances)
 
 return instance
}

// RandomStrategy 随机策略
type RandomStrategy struct{}

// Select 选择实例
func (rs *RandomStrategy) Select(instances []*ServiceInstance) *ServiceInstance {
 if len(instances) == 0 {
  return nil
 }
 
 // 使用时间戳作为随机种子
 rand.Seed(time.Now().UnixNano())
 index := rand.Intn(len(instances))
 
 return instances[index]
}

// LeastConnectionsStrategy 最少连接策略
type LeastConnectionsStrategy struct {
 connectionCounts map[string]int
 mu               sync.RWMutex
}

// NewLeastConnectionsStrategy 创建最少连接策略
func NewLeastConnectionsStrategy() *LeastConnectionsStrategy {
 return &LeastConnectionsStrategy{
  connectionCounts: make(map[string]int),
 }
}

// Select 选择实例
func (lcs *LeastConnectionsStrategy) Select(instances []*ServiceInstance) *ServiceInstance {
 lcs.mu.RLock()
 defer lcs.mu.RUnlock()
 
 if len(instances) == 0 {
  return nil
 }
 
 var selected *ServiceInstance
 minConnections := int(^uint(0) >> 1) // 最大整数值
 
 for _, instance := range instances {
  count := lcs.connectionCounts[instance.ID]
  if count < minConnections {
   minConnections = count
   selected = instance
  }
 }
 
 return selected
}

// IncrementConnections 增加连接数
func (lcs *LeastConnectionsStrategy) IncrementConnections(instanceID string) {
 lcs.mu.Lock()
 defer lcs.mu.Unlock()
 
 lcs.connectionCounts[instanceID]++
}

// DecrementConnections 减少连接数
func (lcs *LeastConnectionsStrategy) DecrementConnections(instanceID string) {
 lcs.mu.Lock()
 defer lcs.mu.Unlock()
 
 if count := lcs.connectionCounts[instanceID]; count > 0 {
  lcs.connectionCounts[instanceID] = count - 1
 }
}
```

## 4. 熔断器

### 4.1 理论基础

**熔断器模式** 是一种保护系统免受级联故障影响的机制。

**状态**：

- **关闭 (Closed)**：正常状态，请求正常处理
- **开启 (Open)**：故障状态，快速失败
- **半开 (Half-Open)**：恢复状态，允许部分请求

### 4.2 Go语言实现

```go
// CircuitBreaker 熔断器
type CircuitBreaker struct {
 state           CircuitBreakerState
 failureCount    int
 lastFailureTime time.Time
 threshold       int
 timeout         time.Duration
 mu              sync.RWMutex
}

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int

const (
 StateClosed CircuitBreakerState = iota
 StateOpen
 StateHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
 return &CircuitBreaker{
  state:     StateClosed,
  threshold: threshold,
  timeout:   timeout,
 }
}

// Execute 执行请求
func (cb *CircuitBreaker) Execute(command func() error) error {
 cb.mu.Lock()
 defer cb.mu.Unlock()
 
 switch cb.state {
 case StateClosed:
  return cb.executeClosed(command)
 case StateOpen:
  return cb.executeOpen()
 case StateHalfOpen:
  return cb.executeHalfOpen(command)
 default:
  return fmt.Errorf("未知的熔断器状态")
 }
}

// executeClosed 关闭状态执行
func (cb *CircuitBreaker) executeClosed(command func() error) error {
 err := command()
 if err != nil {
  cb.failureCount++
  cb.lastFailureTime = time.Now()
  
  if cb.failureCount >= cb.threshold {
   cb.state = StateOpen
   log.Printf("熔断器开启，失败次数: %d", cb.failureCount)
  }
  return err
 }
 
 // 成功时重置失败计数
 cb.failureCount = 0
 return nil
}

// executeOpen 开启状态执行
func (cb *CircuitBreaker) executeOpen() error {
 if time.Since(cb.lastFailureTime) >= cb.timeout {
  cb.state = StateHalfOpen
  log.Println("熔断器进入半开状态")
  return cb.executeHalfOpen(nil)
 }
 
 return fmt.Errorf("熔断器开启，请求被拒绝")
}

// executeHalfOpen 半开状态执行
func (cb *CircuitBreaker) executeHalfOpen(command func() error) error {
 if command == nil {
  return fmt.Errorf("熔断器半开状态，需要提供命令")
 }
 
 err := command()
 if err != nil {
  cb.state = StateOpen
  cb.lastFailureTime = time.Now()
  log.Println("熔断器重新开启")
  return err
 }
 
 // 成功时关闭熔断器
 cb.state = StateClosed
 cb.failureCount = 0
 log.Println("熔断器关闭")
 return nil
}

// GetState 获取状态
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
 cb.mu.RLock()
 defer cb.mu.RUnlock()
 return cb.state
}

// GetFailureCount 获取失败次数
func (cb *CircuitBreaker) GetFailureCount() int {
 cb.mu.RLock()
 defer cb.mu.RUnlock()
 return cb.failureCount
}
```

## 5. 微服务客户端

### 5.1 基础客户端

```go
// MicroserviceClient 微服务客户端
type MicroserviceClient struct {
 discoveryClient *ServiceDiscoveryClient
 loadBalancer    *LoadBalancer
 circuitBreaker  *CircuitBreaker
 httpClient      *http.Client
}

// NewMicroserviceClient 创建微服务客户端
func NewMicroserviceClient(registryURL string) *MicroserviceClient {
 return &MicroserviceClient{
  discoveryClient: NewServiceDiscoveryClient(registryURL),
  loadBalancer:    NewLoadBalancer(&RoundRobinStrategy{}),
  circuitBreaker:  NewCircuitBreaker(5, 30*time.Second),
  httpClient: &http.Client{
   Timeout: 10 * time.Second,
  },
 }
}

// Call 调用服务
func (mc *MicroserviceClient) Call(serviceName, path string, method string, body interface{}) (*http.Response, error) {
 // 发现服务实例
 instances, err := mc.discoveryClient.Discover(serviceName)
 if err != nil {
  return nil, fmt.Errorf("服务发现失败: %v", err)
 }
 
 if len(instances) == 0 {
  return nil, fmt.Errorf("没有可用的服务实例: %s", serviceName)
 }
 
 // 更新负载均衡器
 mc.loadBalancer.UpdateInstances(instances)
 
 // 使用熔断器执行请求
 var response *http.Response
 err = mc.circuitBreaker.Execute(func() error {
  instance := mc.loadBalancer.GetInstance()
  if instance == nil {
   return fmt.Errorf("无法获取服务实例")
  }
  
  url := fmt.Sprintf("http://%s:%d%s", instance.Host, instance.Port, path)
  response, err = mc.makeRequest(url, method, body)
  return err
 })
 
 return response, err
}

// makeRequest 发送HTTP请求
func (mc *MicroserviceClient) makeRequest(url, method string, body interface{}) (*http.Response, error) {
 var req *http.Request
 var err error
 
 if body != nil {
  data, err := json.Marshal(body)
  if err != nil {
   return nil, err
  }
  
  req, err = http.NewRequest(method, url, bytes.NewReader(data))
  if err != nil {
   return nil, err
  }
  
  req.Header.Set("Content-Type", "application/json")
 } else {
  req, err = http.NewRequest(method, url, nil)
  if err != nil {
   return nil, err
  }
 }
 
 return mc.httpClient.Do(req)
}
```

### 5.2 泛型客户端

```go
// GenericMicroserviceClient 泛型微服务客户端
type GenericMicroserviceClient[T any] struct {
 baseClient *MicroserviceClient
}

// NewGenericMicroserviceClient 创建泛型微服务客户端
func NewGenericMicroserviceClient[T any](registryURL string) *GenericMicroserviceClient[T] {
 return &GenericMicroserviceClient[T]{
  baseClient: NewMicroserviceClient(registryURL),
 }
}

// Get 发送GET请求
func (gmc *GenericMicroserviceClient[T]) Get(serviceName, path string) (T, error) {
 var result T
 
 resp, err := gmc.baseClient.Call(serviceName, path, http.MethodGet, nil)
 if err != nil {
  return result, err
 }
 defer resp.Body.Close()
 
 if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
  return result, err
 }
 
 return result, nil
}

// Post 发送POST请求
func (gmc *GenericMicroserviceClient[T]) Post(serviceName, path string, data interface{}) (T, error) {
 var result T
 
 resp, err := gmc.baseClient.Call(serviceName, path, http.MethodPost, data)
 if err != nil {
  return result, err
 }
 defer resp.Body.Close()
 
 if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
  return result, err
 }
 
 return result, nil
}

// Put 发送PUT请求
func (gmc *GenericMicroserviceClient[T]) Put(serviceName, path string, data interface{}) (T, error) {
 var result T
 
 resp, err := gmc.baseClient.Call(serviceName, path, http.MethodPut, data)
 if err != nil {
  return result, err
 }
 defer resp.Body.Close()
 
 if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
  return result, err
 }
 
 return result, nil
}

// Delete 发送DELETE请求
func (gmc *GenericMicroserviceClient[T]) Delete(serviceName, path string) error {
 _, err := gmc.baseClient.Call(serviceName, path, http.MethodDelete, nil)
 return err
}
```

## 6. 微服务示例

### 6.1 用户服务

```go
// UserService 用户服务
type UserService struct {
 registryClient *ServiceDiscoveryClient
 instance       *ServiceInstance
 users          map[string]*User
 mu             sync.RWMutex
}

// User 用户模型
type User struct {
 ID       string    `json:"id"`
 Username string    `json:"username"`
 Email    string    `json:"email"`
 Created  time.Time `json:"created"`
}

// NewUserService 创建用户服务
func NewUserService(host string, port int) *UserService {
 instance := &ServiceInstance{
  ID:       fmt.Sprintf("user-service-%s-%d", host, port),
  Name:     "user-service",
  Host:     host,
  Port:     port,
  Metadata: make(map[string]string),
 }
 
 return &UserService{
  registryClient: NewServiceDiscoveryClient("http://localhost:8080"),
  instance:       instance,
  users:          make(map[string]*User),
 }
}

// Start 启动用户服务
func (us *UserService) Start() error {
 // 注册服务
 if err := us.registryClient.Register(us.instance); err != nil {
  return err
 }
 
 // 设置路由
 http.HandleFunc("/users", us.handleUsers)
 http.HandleFunc("/users/", us.handleUser)
 http.HandleFunc("/health", us.handleHealth)
 
 // 启动HTTP服务器
 addr := fmt.Sprintf("%s:%d", us.instance.Host, us.instance.Port)
 log.Printf("用户服务启动在 %s", addr)
 return http.ListenAndServe(addr, nil)
}

// handleUsers 处理用户列表请求
func (us *UserService) handleUsers(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodGet:
  us.getUsers(w, r)
 case http.MethodPost:
  us.createUser(w, r)
 default:
  http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
 }
}

// handleUser 处理单个用户请求
func (us *UserService) handleUser(w http.ResponseWriter, r *http.Request) {
 switch r.Method {
 case http.MethodGet:
  us.getUser(w, r)
 case http.MethodPut:
  us.updateUser(w, r)
 case http.MethodDelete:
  us.deleteUser(w, r)
 default:
  http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
 }
}

// getUsers 获取用户列表
func (us *UserService) getUsers(w http.ResponseWriter, r *http.Request) {
 us.mu.RLock()
 defer us.mu.RUnlock()
 
 users := make([]*User, 0, len(us.users))
 for _, user := range us.users {
  users = append(users, user)
 }
 
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(users)
}

// createUser 创建用户
func (us *UserService) createUser(w http.ResponseWriter, r *http.Request) {
 var user User
 if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
  http.Error(w, "无效的请求体", http.StatusBadRequest)
  return
 }
 
 us.mu.Lock()
 defer us.mu.Unlock()
 
 user.ID = fmt.Sprintf("user-%d", len(us.users)+1)
 user.Created = time.Now()
 us.users[user.ID] = &user
 
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(http.StatusCreated)
 json.NewEncoder(w).Encode(user)
}

// getUser 获取用户
func (us *UserService) getUser(w http.ResponseWriter, r *http.Request) {
 // 从URL路径提取用户ID
 userID := r.URL.Path[len("/users/"):]
 
 us.mu.RLock()
 defer us.mu.RUnlock()
 
 user, exists := us.users[userID]
 if !exists {
  http.Error(w, "用户不存在", http.StatusNotFound)
  return
 }
 
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(user)
}

// updateUser 更新用户
func (us *UserService) updateUser(w http.ResponseWriter, r *http.Request) {
 userID := r.URL.Path[len("/users/"):]
 
 var user User
 if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
  http.Error(w, "无效的请求体", http.StatusBadRequest)
  return
 }
 
 us.mu.Lock()
 defer us.mu.Unlock()
 
 if _, exists := us.users[userID]; !exists {
  http.Error(w, "用户不存在", http.StatusNotFound)
  return
 }
 
 user.ID = userID
 us.users[userID] = &user
 
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(user)
}

// deleteUser 删除用户
func (us *UserService) deleteUser(w http.ResponseWriter, r *http.Request) {
 userID := r.URL.Path[len("/users/"):]
 
 us.mu.Lock()
 defer us.mu.Unlock()
 
 if _, exists := us.users[userID]; !exists {
  http.Error(w, "用户不存在", http.StatusNotFound)
  return
 }
 
 delete(us.users, userID)
 w.WriteHeader(http.StatusNoContent)
}

// handleHealth 健康检查
func (us *UserService) handleHealth(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]string{
  "status": "healthy",
 })
}
```

## 7. 总结

### 7.1 微服务最佳实践

1. **服务设计**：单一职责、高内聚低耦合
2. **服务发现**：动态注册和发现
3. **负载均衡**：选择合适的负载均衡策略
4. **熔断器**：防止级联故障
5. **监控**：全面的服务监控

### 7.2 性能优化

1. **连接池**：复用HTTP连接
2. **缓存**：缓存服务发现结果
3. **异步**：使用异步通信
4. **压缩**：启用响应压缩

### 7.3 安全考虑

1. **认证**：服务间认证
2. **授权**：访问控制
3. **加密**：传输加密
4. **审计**：操作日志

---

**参考文献**：

1. Newman, S. (2021). Building Microservices
2. Richardson, C. (2018). Microservices Patterns
3. Go官方文档：net/http包
