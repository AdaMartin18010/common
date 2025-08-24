# 组件架构 (Component Architecture)

## 概述

组件架构是一种软件设计模式，将系统分解为独立的、可重用的组件，每个组件具有明确定义的接口和职责。组件可以独立开发、测试、部署和维护，通过标准化的接口进行通信和协作，实现系统的模块化和可扩展性。

## 基本概念

### 核心特征

- **模块化设计**：系统分解为独立的组件模块
- **接口标准化**：组件间通过标准化接口通信
- **松耦合**：组件间最小化依赖关系
- **高内聚**：组件内部功能紧密相关
- **可重用性**：组件可在不同场景中重复使用
- **可替换性**：组件可以独立升级和替换

### 应用场景

- **企业级应用**：大型企业系统的模块化开发
- **微服务架构**：服务组件的独立部署和管理
- **插件系统**：可扩展的插件化架构
- **框架开发**：可配置的框架组件
- **第三方集成**：与外部系统的组件化集成

## 核心组件

### 组件接口 (Component Interface)

```go
// 组件接口定义
type Component interface {
    ID() string
    Name() string
    Version() string
    Initialize(config *ComponentConfig) error
    Start() error
    Stop() error
    Status() ComponentStatus
    GetInterface(name string) (interface{}, error)
}

// 组件配置
type ComponentConfig struct {
    ID      string                 `json:"id"`
    Name    string                 `json:"name"`
    Version string                 `json:"version"`
    Config  map[string]interface{} `json:"config"`
}

// 组件状态
type ComponentStatus struct {
    State     string `json:"state"`
    Health    string `json:"health"`
    StartTime int64  `json:"start_time"`
    Uptime    int64  `json:"uptime"`
}

// 基础组件实现
type BaseComponent struct {
    id      string
    name    string
    version string
    config  *ComponentConfig
    state   string
    startTime int64
    interfaces map[string]interface{}
}

func (bc *BaseComponent) ID() string {
    return bc.id
}

func (bc *BaseComponent) Name() string {
    return bc.name
}

func (bc *BaseComponent) Version() string {
    return bc.version
}

func (bc *BaseComponent) Initialize(config *ComponentConfig) error {
    bc.id = config.ID
    bc.name = config.Name
    bc.version = config.Version
    bc.config = config
    bc.state = "Initialized"
    bc.interfaces = make(map[string]interface{})
    return nil
}

func (bc *BaseComponent) Start() error {
    bc.state = "Running"
    bc.startTime = time.Now().Unix()
    return nil
}

func (bc *BaseComponent) Stop() error {
    bc.state = "Stopped"
    return nil
}

func (bc *BaseComponent) Status() ComponentStatus {
    uptime := int64(0)
    if bc.startTime > 0 {
        uptime = time.Now().Unix() - bc.startTime
    }
    
    return ComponentStatus{
        State:     bc.state,
        Health:    "Healthy",
        StartTime: bc.startTime,
        Uptime:    uptime,
    }
}

func (bc *BaseComponent) GetInterface(name string) (interface{}, error) {
    if iface, exists := bc.interfaces[name]; exists {
        return iface, nil
    }
    return nil, fmt.Errorf("interface %s not found", name)
}

func (bc *BaseComponent) RegisterInterface(name string, iface interface{}) {
    bc.interfaces[name] = iface
}
```

### 组件管理器 (Component Manager)

```go
// 组件管理器
type ComponentManager struct {
    components map[string]Component
    registry   map[string]ComponentFactory
    config     *ManagerConfig
}

// 组件工厂接口
type ComponentFactory func(config *ComponentConfig) (Component, error)

// 管理器配置
type ManagerConfig struct {
    AutoStart     bool   `json:"auto_start"`
    HealthCheck   bool   `json:"health_check"`
    CheckInterval int    `json:"check_interval"`
    LogLevel      string `json:"log_level"`
}

func NewComponentManager(config *ManagerConfig) *ComponentManager {
    return &ComponentManager{
        components: make(map[string]Component),
        registry:   make(map[string]ComponentFactory),
        config:     config,
    }
}

func (cm *ComponentManager) RegisterComponent(name string, factory ComponentFactory) {
    cm.registry[name] = factory
}

func (cm *ComponentManager) CreateComponent(name string, config *ComponentConfig) (Component, error) {
    factory, exists := cm.registry[name]
    if !exists {
        return nil, fmt.Errorf("component factory %s not found", name)
    }
    
    component, err := factory(config)
    if err != nil {
        return nil, err
    }
    
    cm.components[config.ID] = component
    return component, nil
}

func (cm *ComponentManager) GetComponent(id string) (Component, error) {
    component, exists := cm.components[id]
    if !exists {
        return nil, fmt.Errorf("component %s not found", id)
    }
    return component, nil
}

func (cm *ComponentManager) StartComponent(id string) error {
    component, err := cm.GetComponent(id)
    if err != nil {
        return err
    }
    
    return component.Start()
}

func (cm *ComponentManager) StopComponent(id string) error {
    component, err := cm.GetComponent(id)
    if err != nil {
        return err
    }
    
    return component.Stop()
}

func (cm *ComponentManager) ListComponents() []Component {
    components := make([]Component, 0, len(cm.components))
    for _, component := range cm.components {
        components = append(components, component)
    }
    return components
}
```

### 组件通信 (Component Communication)

```go
// 消息接口
type Message interface {
    ID() string
    Type() string
    Source() string
    Target() string
    Data() interface{}
    Timestamp() int64
}

// 基础消息实现
type BaseMessage struct {
    id        string
    msgType   string
    source    string
    target    string
    data      interface{}
    timestamp int64
}

func (bm *BaseMessage) ID() string {
    return bm.id
}

func (bm *BaseMessage) Type() string {
    return bm.msgType
}

func (bm *BaseMessage) Source() string {
    return bm.source
}

func (bm *BaseMessage) Target() string {
    return bm.target
}

func (bm *BaseMessage) Data() interface{} {
    return bm.data
}

func (bm *BaseMessage) Timestamp() int64 {
    return bm.timestamp
}

// 消息总线
type MessageBus struct {
    subscribers map[string][]MessageHandler
    messages    chan Message
    running     bool
}

// 消息处理器
type MessageHandler func(message Message) error

func NewMessageBus() *MessageBus {
    return &MessageBus{
        subscribers: make(map[string][]MessageHandler),
        messages:    make(chan Message, 1000),
        running:     false,
    }
}

func (mb *MessageBus) Subscribe(msgType string, handler MessageHandler) {
    mb.subscribers[msgType] = append(mb.subscribers[msgType], handler)
}

func (mb *MessageBus) Publish(message Message) error {
    if !mb.running {
        return fmt.Errorf("message bus not running")
    }
    
    mb.messages <- message
    return nil
}

func (mb *MessageBus) Start() {
    mb.running = true
    go mb.processMessages()
}

func (mb *MessageBus) Stop() {
    mb.running = false
    close(mb.messages)
}

func (mb *MessageBus) processMessages() {
    for message := range mb.messages {
        handlers := mb.subscribers[message.Type()]
        for _, handler := range handlers {
            go func(h MessageHandler, m Message) {
                if err := h(m); err != nil {
                    log.Printf("Error handling message: %v", err)
                }
            }(handler, message)
        }
    }
}
```

### 组件生命周期管理 (Component Lifecycle)

```go
// 生命周期管理器
type LifecycleManager struct {
    components map[string]Component
    order      []string
    started    map[string]bool
}

func NewLifecycleManager() *LifecycleManager {
    return &LifecycleManager{
        components: make(map[string]Component),
        order:      make([]string, 0),
        started:    make(map[string]bool),
    }
}

func (lm *LifecycleManager) AddComponent(component Component, dependencies []string) {
    lm.components[component.ID()] = component
    lm.order = append(lm.order, component.ID())
}

func (lm *LifecycleManager) StartAll() error {
    // 按依赖顺序启动组件
    for _, id := range lm.order {
        if err := lm.startComponent(id); err != nil {
            return fmt.Errorf("failed to start component %s: %v", id, err)
        }
    }
    return nil
}

func (lm *LifecycleManager) StopAll() error {
    // 按相反顺序停止组件
    for i := len(lm.order) - 1; i >= 0; i-- {
        id := lm.order[i]
        if err := lm.stopComponent(id); err != nil {
            return fmt.Errorf("failed to stop component %s: %v", id, err)
        }
    }
    return nil
}

func (lm *LifecycleManager) startComponent(id string) error {
    component := lm.components[id]
    if err := component.Start(); err != nil {
        return err
    }
    lm.started[id] = true
    return nil
}

func (lm *LifecycleManager) stopComponent(id string) error {
    component := lm.components[id]
    if err := component.Stop(); err != nil {
        return err
    }
    lm.started[id] = false
    return nil
}
```

### 具体组件实现示例

```go
// 数据库组件
type DatabaseComponent struct {
    BaseComponent
    db     *sql.DB
    config *DatabaseConfig
}

type DatabaseConfig struct {
    Driver   string `json:"driver"`
    DSN      string `json:"dsn"`
    MaxConns int    `json:"max_conns"`
}

func NewDatabaseComponent() Component {
    return &DatabaseComponent{}
}

func (dc *DatabaseComponent) Initialize(config *ComponentConfig) error {
    if err := dc.BaseComponent.Initialize(config); err != nil {
        return err
    }
    
    // 解析数据库配置
    dbConfig := &DatabaseConfig{}
    if configData, ok := config.Config["database"].(map[string]interface{}); ok {
        dbConfig.Driver = configData["driver"].(string)
        dbConfig.DSN = configData["dsn"].(string)
        dbConfig.MaxConns = int(configData["max_conns"].(float64))
    }
    
    dc.config = dbConfig
    return nil
}

func (dc *DatabaseComponent) Start() error {
    if err := dc.BaseComponent.Start(); err != nil {
        return err
    }
    
    // 连接数据库
    db, err := sql.Open(dc.config.Driver, dc.config.DSN)
    if err != nil {
        return err
    }
    
    db.SetMaxOpenConns(dc.config.MaxConns)
    dc.db = db
    
    // 注册数据库接口
    dc.RegisterInterface("database", dc.db)
    return nil
}

func (dc *DatabaseComponent) Stop() error {
    if dc.db != nil {
        dc.db.Close()
    }
    return dc.BaseComponent.Stop()
}

// 缓存组件
type CacheComponent struct {
    BaseComponent
    cache  map[string]interface{}
    config *CacheConfig
}

type CacheConfig struct {
    MaxSize int `json:"max_size"`
    TTL     int `json:"ttl"`
}

func NewCacheComponent() Component {
    return &CacheComponent{
        cache: make(map[string]interface{}),
    }
}

func (cc *CacheComponent) Initialize(config *ComponentConfig) error {
    if err := cc.BaseComponent.Initialize(config); err != nil {
        return err
    }
    
    // 解析缓存配置
    cacheConfig := &CacheConfig{}
    if configData, ok := config.Config["cache"].(map[string]interface{}); ok {
        cacheConfig.MaxSize = int(configData["max_size"].(float64))
        cacheConfig.TTL = int(configData["ttl"].(float64))
    }
    
    cc.config = cacheConfig
    return nil
}

func (cc *CacheComponent) Start() error {
    if err := cc.BaseComponent.Start(); err != nil {
        return err
    }
    
    // 注册缓存接口
    cc.RegisterInterface("cache", cc)
    return nil
}

func (cc *CacheComponent) Get(key string) (interface{}, bool) {
    value, exists := cc.cache[key]
    return value, exists
}

func (cc *CacheComponent) Set(key string, value interface{}) {
    if len(cc.cache) >= cc.config.MaxSize {
        // 简单的LRU策略：删除第一个元素
        for k := range cc.cache {
            delete(cc.cache, k)
            break
        }
    }
    cc.cache[key] = value
}

func (cc *CacheComponent) Delete(key string) {
    delete(cc.cache, key)
}
```

## 设计原则

### 1. 单一职责原则

- **明确职责**：每个组件只负责一个特定的功能领域
- **内聚性**：组件内部功能紧密相关
- **独立性**：组件可以独立开发和测试

### 2. 接口隔离原则

- **最小接口**：组件只暴露必要的接口
- **接口稳定**：保持接口的向后兼容性
- **版本管理**：支持接口的版本控制

### 3. 依赖倒置原则

- **抽象依赖**：组件依赖抽象而非具体实现
- **接口编程**：通过接口进行组件间通信
- **依赖注入**：通过配置注入组件依赖

### 4. 开闭原则

- **扩展开放**：支持新功能的扩展
- **修改封闭**：避免修改现有组件
- **插件化**：支持插件式扩展

## 实现示例

```go
func main() {
    // 创建组件管理器
    manager := NewComponentManager(&ManagerConfig{
        AutoStart:     true,
        HealthCheck:   true,
        CheckInterval: 30,
        LogLevel:      "INFO",
    })
    
    // 注册组件工厂
    manager.RegisterComponent("database", func(config *ComponentConfig) (Component, error) {
        component := NewDatabaseComponent()
        return component, component.Initialize(config)
    })
    
    manager.RegisterComponent("cache", func(config *ComponentConfig) (Component, error) {
        component := NewCacheComponent()
        return component, component.Initialize(config)
    })
    
    // 创建消息总线
    messageBus := NewMessageBus()
    messageBus.Start()
    
    // 创建数据库组件
    dbConfig := &ComponentConfig{
        ID:      "db-1",
        Name:    "Main Database",
        Version: "1.0.0",
        Config: map[string]interface{}{
            "database": map[string]interface{}{
                "driver":    "mysql",
                "dsn":      "user:password@tcp(localhost:3306)/testdb",
                "max_conns": 10,
            },
        },
    }
    
    dbComponent, err := manager.CreateComponent("database", dbConfig)
    if err != nil {
        log.Fatalf("Failed to create database component: %v", err)
    }
    
    // 创建缓存组件
    cacheConfig := &ComponentConfig{
        ID:      "cache-1",
        Name:    "Redis Cache",
        Version: "1.0.0",
        Config: map[string]interface{}{
            "cache": map[string]interface{}{
                "max_size": 1000,
                "ttl":      3600,
            },
        },
    }
    
    cacheComponent, err := manager.CreateComponent("cache", cacheConfig)
    if err != nil {
        log.Fatalf("Failed to create cache component: %v", err)
    }
    
    // 启动组件
    if err := manager.StartComponent("db-1"); err != nil {
        log.Fatalf("Failed to start database component: %v", err)
    }
    
    if err := manager.StartComponent("cache-1"); err != nil {
        log.Fatalf("Failed to start cache component: %v", err)
    }
    
    // 获取组件接口
    if dbInterface, err := dbComponent.GetInterface("database"); err == nil {
        db := dbInterface.(*sql.DB)
        log.Printf("Database connected: %v", db.Stats())
    }
    
    if cacheInterface, err := cacheComponent.GetInterface("cache"); err == nil {
        cache := cacheInterface.(*CacheComponent)
        cache.Set("key1", "value1")
        if value, exists := cache.Get("key1"); exists {
            log.Printf("Cache value: %v", value)
        }
    }
    
    // 发布消息
    message := &BaseMessage{
        id:        "msg-1",
        msgType:   "data_updated",
        source:    "db-1",
        target:    "cache-1",
        data:      map[string]interface{}{"table": "users", "action": "insert"},
        timestamp: time.Now().Unix(),
    }
    
    messageBus.Publish(message)
    
    // 等待一段时间后停止
    time.Sleep(5 * time.Second)
    
    manager.StopComponent("cache-1")
    manager.StopComponent("db-1")
    messageBus.Stop()
    
    log.Println("Component system stopped")
}
```

## 总结

组件架构通过模块化设计、标准化接口和松耦合原则，实现了系统的可维护性、可扩展性和可重用性。

### 关键要点

1. **模块化设计**：将系统分解为独立的组件
2. **接口标准化**：通过标准化接口进行组件通信
3. **生命周期管理**：统一的组件生命周期管理
4. **消息通信**：异步消息总线进行组件间通信
5. **依赖管理**：通过依赖注入管理组件依赖

### 发展趋势

- **微服务化**：组件向微服务架构演进
- **容器化部署**：支持容器化部署和编排
- **云原生**：与云原生技术深度集成
- **AI集成**：智能化的组件管理和优化
- **边缘计算**：支持边缘计算场景的组件架构
