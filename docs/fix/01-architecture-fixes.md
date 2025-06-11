# 架构修复方案

## 目录

1. [当前架构问题](#当前架构问题)
2. [修复目标](#修复目标)
3. [重构方案](#重构方案)
4. [具体实现](#具体实现)
5. [迁移策略](#迁移策略)
6. [测试验证](#测试验证)

## 当前架构问题

### 1.1 控制流复杂性问题

当前的控制流设计存在以下问题：

```go
// 问题代码：复杂的控制结构
type CtrlSt struct {
    ctx    context.Context
    cancel context.CancelFunc
    wg     *WorkerWG
}

type WorkerWG struct {
    wg sync.WaitGroup
    // 复杂的goroutine管理逻辑
}
```

**问题分析**：

- CtrlSt和WorkerWG交互复杂
- 上下文管理分散
- 错误传播不清晰
- 调试困难

### 1.2 组件耦合问题

```go
// 问题代码：组件间耦合度高
type CptMetaSt struct {
    id    string
    kind  string
    state atomic.Value
    // 直接依赖具体实现
}
```

**问题分析**：

- 组件间直接依赖
- 缺乏抽象层次
- 难以进行单元测试
- 扩展性差

### 1.3 配置管理问题

```go
// 问题代码：配置管理分散
func CompiledExectionFilePath() string {
    _, file, _, ok := runtime.Caller(1)
    if !ok {
        log.Default().Printf("Can not get current file info")
        return ""
    }
    // 硬编码的配置逻辑
}
```

**问题分析**：

- 配置逻辑分散
- 缺乏统一管理
- 难以动态配置
- 测试困难

## 修复目标

### 2.1 简化控制流

**目标**：

- 简化CtrlSt和WorkerWG的交互
- 统一上下文管理
- 清晰错误传播路径
- 提升调试能力

### 2.2 解耦组件系统

**目标**：

- 引入依赖注入机制
- 建立抽象层次
- 提升可测试性
- 增强扩展性

### 2.3 统一配置管理

**目标**：

- 建立统一配置抽象
- 支持多种配置源
- 实现动态配置
- 提升配置可测试性

## 重构方案

### 3.1 简化控制流设计

#### 3.1.1 新的控制流架构

```go
// 简化的控制流设计
type LifecycleManager struct {
    ctx       context.Context
    cancel    context.CancelFunc
    workers   *WorkerPool
    logger    *zap.Logger
    metrics   *MetricsCollector
}

type WorkerPool struct {
    workers map[string]*Worker
    mu      sync.RWMutex
    logger  *zap.Logger
}

type Worker struct {
    id       string
    name     string
    handler  WorkerHandler
    status   WorkerStatus
    ctx      context.Context
    cancel   context.CancelFunc
    logger   *zap.Logger
}

type WorkerHandler func(ctx context.Context) error
type WorkerStatus int

const (
    WorkerStatusStopped WorkerStatus = iota
    WorkerStatusStarting
    WorkerStatusRunning
    WorkerStatusStopping
    WorkerStatusError
)
```

#### 3.1.2 生命周期管理

```go
// 生命周期管理器实现
func NewLifecycleManager(ctx context.Context) *LifecycleManager {
    ctx, cancel := context.WithCancel(ctx)
    return &LifecycleManager{
        ctx:       ctx,
        cancel:    cancel,
        workers:   NewWorkerPool(),
        logger:    zap.L().Named("lifecycle-manager"),
        metrics:   NewMetricsCollector(),
    }
}

func (lm *LifecycleManager) Start() error {
    lm.logger.Info("starting lifecycle manager")
    
    // 启动所有worker
    if err := lm.workers.StartAll(lm.ctx); err != nil {
        lm.logger.Error("failed to start workers", zap.Error(err))
        return fmt.Errorf("failed to start workers: %w", err)
    }
    
    lm.logger.Info("lifecycle manager started")
    return nil
}

func (lm *LifecycleManager) Stop() error {
    lm.logger.Info("stopping lifecycle manager")
    
    // 取消上下文
    lm.cancel()
    
    // 等待所有worker停止
    if err := lm.workers.StopAll(); err != nil {
        lm.logger.Error("failed to stop workers", zap.Error(err))
        return fmt.Errorf("failed to stop workers: %w", err)
    }
    
    lm.logger.Info("lifecycle manager stopped")
    return nil
}

func (lm *LifecycleManager) AddWorker(name string, handler WorkerHandler) error {
    return lm.workers.AddWorker(name, handler)
}

func (lm *LifecycleManager) RemoveWorker(name string) error {
    return lm.workers.RemoveWorker(name)
}
```

### 3.2 依赖注入容器

#### 3.2.1 依赖注入接口

```go
// 依赖注入容器
type DependencyContainer struct {
    services map[string]interface{}
    factories map[string]ServiceFactory
    mu       sync.RWMutex
    logger   *zap.Logger
}

type ServiceFactory func() (interface{}, error)

func NewDependencyContainer() *DependencyContainer {
    return &DependencyContainer{
        services:  make(map[string]interface{}),
        factories: make(map[string]ServiceFactory),
        logger:    zap.L().Named("dependency-container"),
    }
}

func (dc *DependencyContainer) RegisterService(name string, service interface{}) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    dc.services[name] = service
    dc.logger.Info("service registered", zap.String("name", name))
}

func (dc *DependencyContainer) RegisterFactory(name string, factory ServiceFactory) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    dc.factories[name] = factory
    dc.logger.Info("service factory registered", zap.String("name", name))
}

func (dc *DependencyContainer) GetService(name string) (interface{}, error) {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    
    // 先检查已注册的服务
    if service, exists := dc.services[name]; exists {
        return service, nil
    }
    
    // 检查工厂
    if factory, exists := dc.factories[name]; exists {
        service, err := factory()
        if err != nil {
            return nil, fmt.Errorf("failed to create service %s: %w", name, err)
        }
        
        // 缓存服务
        dc.services[name] = service
        dc.logger.Info("service created", zap.String("name", name))
        
        return service, nil
    }
    
    return nil, fmt.Errorf("service %s not found", name)
}

func (dc *DependencyContainer) GetServiceTyped(name string, serviceType interface{}) error {
    service, err := dc.GetService(name)
    if err != nil {
        return err
    }
    
    // 使用反射设置值
    reflect.ValueOf(serviceType).Elem().Set(reflect.ValueOf(service))
    return nil
}
```

#### 3.2.2 组件依赖注入

```go
// 依赖注入组件
type InjectableComponent struct {
    id       string
    kind     string
    container *DependencyContainer
    logger   *zap.Logger
}

func NewInjectableComponent(id, kind string, container *DependencyContainer) *InjectableComponent {
    return &InjectableComponent{
        id:        id,
        kind:      kind,
        container: container,
        logger:    zap.L().Named("injectable-component"),
    }
}

func (ic *InjectableComponent) InjectDependencies() error {
    // 注入日志服务
    var logger *zap.Logger
    if err := ic.container.GetServiceTyped("logger", &logger); err != nil {
        return fmt.Errorf("failed to inject logger: %w", err)
    }
    ic.logger = logger.Named(ic.id)
    
    // 注入其他依赖
    // ...
    
    ic.logger.Info("dependencies injected")
    return nil
}

func (ic *InjectableComponent) ID() string {
    return ic.id
}

func (ic *InjectableComponent) Kind() string {
    return ic.kind
}
```

### 3.3 统一配置管理

#### 3.3.1 配置抽象层

```go
// 配置管理器
type ConfigManager struct {
    configs  map[string]interface{}
    sources  []ConfigSource
    viper    *viper.Viper
    logger   *zap.Logger
    watchers map[string][]ConfigWatcher
    mu       sync.RWMutex
}

type ConfigSource interface {
    Load() (map[string]interface{}, error)
    Watch(callback func(map[string]interface{}) error) error
    Name() string
}

type ConfigWatcher func(key string, oldValue, newValue interface{})

func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        configs:  make(map[string]interface{}),
        sources:  make([]ConfigSource, 0),
        viper:    viper.New(),
        logger:   zap.L().Named("config-manager"),
        watchers: make(map[string][]ConfigWatcher),
    }
}

func (cm *ConfigManager) AddSource(source ConfigSource) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.sources = append(cm.sources, source)
    cm.logger.Info("config source added", zap.String("source", source.Name()))
}

func (cm *ConfigManager) Load() error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    for _, source := range cm.sources {
        config, err := source.Load()
        if err != nil {
            cm.logger.Error("failed to load config from source", 
                zap.String("source", source.Name()),
                zap.Error(err))
            return fmt.Errorf("failed to load config from %s: %w", source.Name(), err)
        }
        
        // 合并配置
        for key, value := range config {
            cm.configs[key] = value
        }
        
        cm.logger.Info("config loaded from source", zap.String("source", source.Name()))
    }
    
    return nil
}

func (cm *ConfigManager) Get(key string) interface{} {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    return cm.configs[key]
}

func (cm *ConfigManager) GetString(key string) string {
    if value := cm.Get(key); value != nil {
        if str, ok := value.(string); ok {
            return str
        }
    }
    return ""
}

func (cm *ConfigManager) GetInt(key string) int {
    if value := cm.Get(key); value != nil {
        if i, ok := value.(int); ok {
            return i
        }
    }
    return 0
}

func (cm *ConfigManager) Watch(key string, watcher ConfigWatcher) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.watchers[key] = append(cm.watchers[key], watcher)
    cm.logger.Info("config watcher added", zap.String("key", key))
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    oldValue := cm.configs[key]
    cm.configs[key] = value
    
    // 通知观察者
    if watchers, exists := cm.watchers[key]; exists {
        for _, watcher := range watchers {
            go watcher(key, oldValue, value)
        }
    }
    
    cm.logger.Info("config updated", zap.String("key", key))
}
```

#### 3.3.2 配置源实现

```go
// 文件配置源
type FileConfigSource struct {
    path     string
    format   string
    logger   *zap.Logger
}

func NewFileConfigSource(path, format string) *FileConfigSource {
    return &FileConfigSource{
        path:   path,
        format: format,
        logger: zap.L().Named("file-config-source"),
    }
}

func (fcs *FileConfigSource) Load() (map[string]interface{}, error) {
    viper := viper.New()
    viper.SetConfigFile(fcs.path)
    viper.SetConfigType(fcs.format)
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    config := make(map[string]interface{})
    for _, key := range viper.AllKeys() {
        config[key] = viper.Get(key)
    }
    
    fcs.logger.Info("config loaded from file", zap.String("path", fcs.path))
    return config, nil
}

func (fcs *FileConfigSource) Watch(callback func(map[string]interface{}) error) error {
    viper := viper.New()
    viper.SetConfigFile(fcs.path)
    viper.SetConfigType(fcs.format)
    
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fcs.logger.Info("config file changed", zap.String("path", fcs.path))
        
        config, err := fcs.Load()
        if err != nil {
            fcs.logger.Error("failed to reload config", zap.Error(err))
            return
        }
        
        if err := callback(config); err != nil {
            fcs.logger.Error("failed to apply config change", zap.Error(err))
        }
    })
    
    return nil
}

func (fcs *FileConfigSource) Name() string {
    return fmt.Sprintf("file:%s", fcs.path)
}

// 环境变量配置源
type EnvConfigSource struct {
    prefix string
    logger *zap.Logger
}

func NewEnvConfigSource(prefix string) *EnvConfigSource {
    return &EnvConfigSource{
        prefix: prefix,
        logger: zap.L().Named("env-config-source"),
    }
}

func (ecs *EnvConfigSource) Load() (map[string]interface{}, error) {
    config := make(map[string]interface{})
    
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) != 2 {
            continue
        }
        
        key, value := pair[0], pair[1]
        if strings.HasPrefix(key, ecs.prefix) {
            // 移除前缀并转换为小写
            configKey := strings.ToLower(strings.TrimPrefix(key, ecs.prefix))
            config[configKey] = value
        }
    }
    
    ecs.logger.Info("config loaded from environment", zap.String("prefix", ecs.prefix))
    return config, nil
}

func (ecs *EnvConfigSource) Watch(callback func(map[string]interface{}) error) error {
    // 环境变量不支持动态监听，这里只是占位符
    return nil
}

func (ecs *EnvConfigSource) Name() string {
    return fmt.Sprintf("env:%s", ecs.prefix)
}
```

## 具体实现

### 4.1 重构后的组件系统

```go
// 重构后的组件接口
type Component interface {
    ID() string
    Kind() string
    Start(ctx context.Context) error
    Stop() error
    Status() ComponentStatus
    Dependencies() []string
}

type ComponentStatus int

const (
    ComponentStatusStopped ComponentStatus = iota
    ComponentStatusStarting
    ComponentStatusRunning
    ComponentStatusStopping
    ComponentStatusError
)

// 基础组件实现
type BaseComponent struct {
    id           string
    kind         string
    status       ComponentStatus
    dependencies []string
    container    *DependencyContainer
    lifecycle    *LifecycleManager
    logger       *zap.Logger
    mu           sync.RWMutex
}

func NewBaseComponent(id, kind string, dependencies []string, container *DependencyContainer) *BaseComponent {
    return &BaseComponent{
        id:           id,
        kind:         kind,
        status:       ComponentStatusStopped,
        dependencies: dependencies,
        container:    container,
        lifecycle:    NewLifecycleManager(context.Background()),
        logger:       zap.L().Named(fmt.Sprintf("component:%s", id)),
    }
}

func (bc *BaseComponent) ID() string {
    return bc.id
}

func (bc *BaseComponent) Kind() string {
    return bc.kind
}

func (bc *BaseComponent) Status() ComponentStatus {
    bc.mu.RLock()
    defer bc.mu.RUnlock()
    return bc.status
}

func (bc *BaseComponent) Dependencies() []string {
    return bc.dependencies
}

func (bc *BaseComponent) Start(ctx context.Context) error {
    bc.mu.Lock()
    defer bc.mu.Unlock()
    
    if bc.status != ComponentStatusStopped {
        return fmt.Errorf("component %s is not in stopped status", bc.id)
    }
    
    bc.status = ComponentStatusStarting
    bc.logger.Info("starting component")
    
    // 检查依赖
    if err := bc.checkDependencies(); err != nil {
        bc.status = ComponentStatusError
        return fmt.Errorf("dependency check failed: %w", err)
    }
    
    // 启动生命周期管理器
    if err := bc.lifecycle.Start(); err != nil {
        bc.status = ComponentStatusError
        return fmt.Errorf("failed to start lifecycle: %w", err)
    }
    
    bc.status = ComponentStatusRunning
    bc.logger.Info("component started")
    
    return nil
}

func (bc *BaseComponent) Stop() error {
    bc.mu.Lock()
    defer bc.mu.Unlock()
    
    if bc.status != ComponentStatusRunning {
        return fmt.Errorf("component %s is not running", bc.id)
    }
    
    bc.status = ComponentStatusStopping
    bc.logger.Info("stopping component")
    
    // 停止生命周期管理器
    if err := bc.lifecycle.Stop(); err != nil {
        bc.status = ComponentStatusError
        return fmt.Errorf("failed to stop lifecycle: %w", err)
    }
    
    bc.status = ComponentStatusStopped
    bc.logger.Info("component stopped")
    
    return nil
}

func (bc *BaseComponent) checkDependencies() error {
    for _, dep := range bc.dependencies {
        if _, err := bc.container.GetService(dep); err != nil {
            return fmt.Errorf("dependency %s not available: %w", dep, err)
        }
    }
    return nil
}
```

### 4.2 组件管理器

```go
// 组件管理器
type ComponentManager struct {
    components map[string]Component
    container  *DependencyContainer
    config     *ConfigManager
    logger     *zap.Logger
    mu         sync.RWMutex
}

func NewComponentManager(container *DependencyContainer, config *ConfigManager) *ComponentManager {
    return &ComponentManager{
        components: make(map[string]Component),
        container:  container,
        config:     config,
        logger:     zap.L().Named("component-manager"),
    }
}

func (cm *ComponentManager) RegisterComponent(component Component) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    if _, exists := cm.components[component.ID()]; exists {
        return fmt.Errorf("component %s already registered", component.ID())
    }
    
    cm.components[component.ID()] = component
    cm.logger.Info("component registered", 
        zap.String("id", component.ID()),
        zap.String("kind", component.Kind()))
    
    return nil
}

func (cm *ComponentManager) StartAll(ctx context.Context) error {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    // 按依赖顺序启动组件
    sortedComponents, err := cm.sortByDependencies()
    if err != nil {
        return fmt.Errorf("failed to sort components: %w", err)
    }
    
    for _, component := range sortedComponents {
        if err := component.Start(ctx); err != nil {
            cm.logger.Error("failed to start component", 
                zap.String("id", component.ID()),
                zap.Error(err))
            return fmt.Errorf("failed to start component %s: %w", component.ID(), err)
        }
    }
    
    cm.logger.Info("all components started")
    return nil
}

func (cm *ComponentManager) StopAll() error {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    // 按依赖逆序停止组件
    sortedComponents, err := cm.sortByDependencies()
    if err != nil {
        return fmt.Errorf("failed to sort components: %w", err)
    }
    
    // 逆序停止
    for i := len(sortedComponents) - 1; i >= 0; i-- {
        component := sortedComponents[i]
        if err := component.Stop(); err != nil {
            cm.logger.Error("failed to stop component", 
                zap.String("id", component.ID()),
                zap.Error(err))
            return fmt.Errorf("failed to stop component %s: %w", component.ID(), err)
        }
    }
    
    cm.logger.Info("all components stopped")
    return nil
}

func (cm *ComponentManager) sortByDependencies() ([]Component, error) {
    // 使用拓扑排序
    components := make([]Component, 0, len(cm.components))
    for _, component := range cm.components {
        components = append(components, component)
    }
    
    // 简单的拓扑排序实现
    sorted := make([]Component, 0)
    visited := make(map[string]bool)
    temp := make(map[string]bool)
    
    var visit func(component Component) error
    visit = func(component Component) error {
        if temp[component.ID()] {
            return fmt.Errorf("circular dependency detected")
        }
        
        if visited[component.ID()] {
            return nil
        }
        
        temp[component.ID()] = true
        
        // 先访问依赖
        for _, depID := range component.Dependencies() {
            dep, exists := cm.components[depID]
            if !exists {
                return fmt.Errorf("dependency %s not found", depID)
            }
            
            if err := visit(dep); err != nil {
                return err
            }
        }
        
        temp[component.ID()] = false
        visited[component.ID()] = true
        sorted = append(sorted, component)
        
        return nil
    }
    
    for _, component := range components {
        if !visited[component.ID()] {
            if err := visit(component); err != nil {
                return nil, err
            }
        }
    }
    
    return sorted, nil
}
```

## 迁移策略

### 5.1 渐进式迁移

#### 5.1.1 第一阶段：基础设施

**目标**：建立新的基础设施

**任务**：

1. 实现依赖注入容器
2. 建立配置管理系统
3. 创建生命周期管理器
4. 编写基础测试

**时间**：2-3周

#### 5.1.2 第二阶段：组件迁移

**目标**：逐步迁移现有组件

**任务**：

1. 创建新的组件接口
2. 实现基础组件类
3. 迁移现有组件
4. 保持向后兼容

**时间**：3-4周

#### 5.1.3 第三阶段：系统集成

**目标**：完成系统集成

**任务**：

1. 集成所有组件
2. 完善错误处理
3. 优化性能
4. 完善文档

**时间**：2-3周

### 5.2 兼容性保证

#### 5.2.1 接口兼容性

```go
// 兼容性适配器
type CompatibilityAdapter struct {
    oldComponent *CptMetaSt
    newComponent Component
}

func NewCompatibilityAdapter(oldComponent *CptMetaSt) *CompatibilityAdapter {
    // 创建新的组件
    newComponent := NewBaseComponent(
        oldComponent.ID(),
        oldComponent.Kind(),
        []string{}, // 需要分析依赖
        nil,        // 需要注入容器
    )
    
    return &CompatibilityAdapter{
        oldComponent: oldComponent,
        newComponent: newComponent,
    }
}

// 实现旧的接口以保持兼容性
func (ca *CompatibilityAdapter) ID() string {
    return ca.oldComponent.ID()
}

func (ca *CompatibilityAdapter) Kind() string {
    return ca.oldComponent.Kind()
}

func (ca *CompatibilityAdapter) Start() error {
    return ca.newComponent.Start(context.Background())
}

func (ca *CompatibilityAdapter) Stop() error {
    return ca.newComponent.Stop()
}
```

#### 5.2.2 配置兼容性

```go
// 配置兼容性适配器
type ConfigCompatibilityAdapter struct {
    oldConfig map[string]interface{}
    newConfig *ConfigManager
}

func NewConfigCompatibilityAdapter(oldConfig map[string]interface{}) *ConfigCompatibilityAdapter {
    newConfig := NewConfigManager()
    
    // 迁移旧配置
    for key, value := range oldConfig {
        newConfig.Set(key, value)
    }
    
    return &ConfigCompatibilityAdapter{
        oldConfig: oldConfig,
        newConfig: newConfig,
    }
}

// 提供兼容的配置访问方法
func (cca *ConfigCompatibilityAdapter) Get(key string) interface{} {
    return cca.newConfig.Get(key)
}

func (cca *ConfigCompatibilityAdapter) GetString(key string) string {
    return cca.newConfig.GetString(key)
}
```

## 测试验证

### 6.1 单元测试

```go
// 组件测试
func TestBaseComponent(t *testing.T) {
    container := NewDependencyContainer()
    component := NewBaseComponent("test", "test-kind", []string{}, container)
    
    // 测试初始状态
    assert.Equal(t, ComponentStatusStopped, component.Status())
    assert.Equal(t, "test", component.ID())
    assert.Equal(t, "test-kind", component.Kind())
    
    // 测试启动
    ctx := context.Background()
    err := component.Start(ctx)
    assert.NoError(t, err)
    assert.Equal(t, ComponentStatusRunning, component.Status())
    
    // 测试停止
    err = component.Stop()
    assert.NoError(t, err)
    assert.Equal(t, ComponentStatusStopped, component.Status())
}

// 依赖注入测试
func TestDependencyInjection(t *testing.T) {
    container := NewDependencyContainer()
    
    // 注册服务
    testService := "test-service"
    container.RegisterService("test", testService)
    
    // 获取服务
    service, err := container.GetService("test")
    assert.NoError(t, err)
    assert.Equal(t, testService, service)
    
    // 测试工厂
    container.RegisterFactory("factory", func() (interface{}, error) {
        return "factory-service", nil
    })
    
    factoryService, err := container.GetService("factory")
    assert.NoError(t, err)
    assert.Equal(t, "factory-service", factoryService)
}

// 配置管理测试
func TestConfigManager(t *testing.T) {
    config := NewConfigManager()
    
    // 测试设置和获取
    config.Set("test-key", "test-value")
    assert.Equal(t, "test-value", config.GetString("test-key"))
    
    // 测试观察者
    var observedKey string
    var observedValue interface{}
    
    config.Watch("test-key", func(key string, oldValue, newValue interface{}) {
        observedKey = key
        observedValue = newValue
    })
    
    config.Set("test-key", "new-value")
    time.Sleep(100 * time.Millisecond) // 等待goroutine执行
    
    assert.Equal(t, "test-key", observedKey)
    assert.Equal(t, "new-value", observedValue)
}
```

### 6.2 集成测试

```go
// 系统集成测试
func TestSystemIntegration(t *testing.T) {
    // 创建系统
    container := NewDependencyContainer()
    config := NewConfigManager()
    manager := NewComponentManager(container, config)
    
    // 注册组件
    component1 := NewBaseComponent("comp1", "test", []string{}, container)
    component2 := NewBaseComponent("comp2", "test", []string{"comp1"}, container)
    
    err := manager.RegisterComponent(component1)
    assert.NoError(t, err)
    
    err = manager.RegisterComponent(component2)
    assert.NoError(t, err)
    
    // 测试启动
    ctx := context.Background()
    err = manager.StartAll(ctx)
    assert.NoError(t, err)
    
    // 验证状态
    assert.Equal(t, ComponentStatusRunning, component1.Status())
    assert.Equal(t, ComponentStatusRunning, component2.Status())
    
    // 测试停止
    err = manager.StopAll()
    assert.NoError(t, err)
    
    // 验证状态
    assert.Equal(t, ComponentStatusStopped, component1.Status())
    assert.Equal(t, ComponentStatusStopped, component2.Status())
}
```

### 6.3 性能测试

```go
// 性能基准测试
func BenchmarkComponentStart(b *testing.B) {
    container := NewDependencyContainer()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        component := NewBaseComponent(fmt.Sprintf("comp-%d", i), "test", []string{}, container)
        component.Start(context.Background())
        component.Stop()
    }
}

func BenchmarkDependencyInjection(b *testing.B) {
    container := NewDependencyContainer()
    container.RegisterService("test", "test-value")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := container.GetService("test")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 总结

通过系统性的架构重构，我们实现了以下目标：

1. **简化控制流**: 通过生命周期管理器统一管理组件生命周期
2. **解耦组件系统**: 通过依赖注入容器实现组件解耦
3. **统一配置管理**: 通过配置管理器提供统一的配置访问
4. **提升可测试性**: 通过接口抽象和依赖注入提升测试能力
5. **增强扩展性**: 通过插件化架构支持功能扩展

关键改进包括：

- **架构清晰**: 明确的层次结构和职责分离
- **依赖管理**: 统一的依赖注入和生命周期管理
- **配置灵活**: 支持多种配置源和动态配置
- **测试友好**: 完整的测试覆盖和模拟支持
- **向后兼容**: 渐进式迁移和兼容性保证

这个架构修复方案为项目的长期发展提供了稳固的基础，确保系统能够持续演进和扩展。
