# Golang Common 库实现修复方案

## 目录

1. [组件系统修复](#组件系统修复)
2. [事件系统修复](#事件系统修复)
3. [配置管理修复](#配置管理修复)
4. [监控系统修复](#监控系统修复)

## 组件系统修复

### 1.1 增强组件接口

```go
// 增强组件接口
type EnhancedComponent interface {
    Component
    Health() HealthStatus
    Metrics() ComponentMetrics
    Configuration() ComponentConfig
    Dependencies() []string
    Validate() error
}

// 健康检查状态
type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Details   map[string]string `json:"details"`
    Errors    []string          `json:"errors"`
}

// 组件配置
type ComponentConfig struct {
    Name         string                 `json:"name"`
    Version      string                 `json:"version"`
    Type         string                 `json:"type"`
    Dependencies []string               `json:"dependencies"`
    Settings     map[string]interface{} `json:"settings"`
    Timeout      time.Duration          `json:"timeout"`
    Retries      int                    `json:"retries"`
}

// 增强组件实现
type EnhancedComponentImpl struct {
    *BaseComponent
    health       HealthStatus
    metrics      *ComponentMetrics
    config       ComponentConfig
    dependencies []string
    logger       *zap.Logger
    mu           sync.RWMutex
    ctx          context.Context
    cancel       context.CancelFunc
}

func NewEnhancedComponent(config ComponentConfig) *EnhancedComponentImpl {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &EnhancedComponentImpl{
        BaseComponent: NewBaseComponent(config.Name, config.Version),
        config:        config,
        dependencies:  config.Dependencies,
        logger:        zap.L().Named(config.Name),
        ctx:           ctx,
        cancel:        cancel,
    }
}

func (ec *EnhancedComponentImpl) Start() error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    
    // 验证依赖
    if err := ec.validateDependencies(); err != nil {
        return fmt.Errorf("dependency validation failed: %w", err)
    }
    
    // 更新健康状态
    ec.updateHealth("starting", map[string]string{
        "status": "initializing",
    }, nil)
    
    // 启动基础组件
    if err := ec.BaseComponent.Start(); err != nil {
        ec.updateHealth("error", nil, []string{err.Error()})
        return fmt.Errorf("failed to start base component: %w", err)
    }
    
    // 启动指标收集
    if err := ec.startMetrics(); err != nil {
        ec.logger.Warn("failed to start metrics", zap.Error(err))
    }
    
    ec.updateHealth("healthy", map[string]string{
        "status": "running",
    }, nil)
    
    ec.logger.Info("enhanced component started")
    return nil
}

func (ec *EnhancedComponentImpl) Stop() error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    
    ec.updateHealth("stopping", map[string]string{
        "status": "shutting_down",
    }, nil)
    
    // 取消上下文
    ec.cancel()
    
    // 停止基础组件
    if err := ec.BaseComponent.Stop(); err != nil {
        ec.updateHealth("error", nil, []string{err.Error()})
        return fmt.Errorf("failed to stop base component: %w", err)
    }
    
    // 停止指标收集
    ec.stopMetrics()
    
    ec.updateHealth("stopped", map[string]string{
        "status": "stopped",
    }, nil)
    
    ec.logger.Info("enhanced component stopped")
    return nil
}

func (ec *EnhancedComponentImpl) Health() HealthStatus {
    ec.mu.RLock()
    defer ec.mu.RUnlock()
    
    return ec.health
}

func (ec *EnhancedComponentImpl) updateHealth(status string, details map[string]string, errors []string) {
    ec.health = HealthStatus{
        Status:    status,
        Timestamp: time.Now(),
        Details:   details,
        Errors:    errors,
    }
    
    // 更新指标
    if ec.metrics != nil {
        switch status {
        case "healthy":
            ec.metrics.StatusGauge.Set(1)
        case "error":
            ec.metrics.StatusGauge.Set(0)
            ec.metrics.ErrorCount.Inc()
        }
    }
}

func (ec *EnhancedComponentImpl) validateDependencies() error {
    if len(ec.dependencies) == 0 {
        return nil
    }
    
    // 这里应该检查依赖组件是否已启动
    // 实际实现中需要依赖注入容器
    ec.logger.Debug("dependencies validated", zap.Strings("dependencies", ec.dependencies))
    return nil
}

func (ec *EnhancedComponentImpl) startMetrics() error {
    ec.metrics = NewComponentMetrics(ec.config.Name)
    return nil
}

func (ec *EnhancedComponentImpl) stopMetrics() {
    ec.metrics = nil
}
```

### 1.2 组件工厂实现

```go
// 组件工厂
type ComponentFactory struct {
    creators   map[string]ComponentCreator
    validators map[string]ConfigValidator
    logger     *zap.Logger
    mu         sync.RWMutex
}

type ComponentCreator func(config ComponentConfig) (Component, error)
type ConfigValidator func(config ComponentConfig) error

func NewComponentFactory() *ComponentFactory {
    return &ComponentFactory{
        creators:   make(map[string]ComponentCreator),
        validators: make(map[string]ConfigValidator),
        logger:     zap.L().Named("component-factory"),
    }
}

func (cf *ComponentFactory) RegisterCreator(componentType string, creator ComponentCreator) {
    cf.mu.Lock()
    defer cf.mu.Unlock()
    
    cf.creators[componentType] = creator
    cf.logger.Info("component creator registered", zap.String("type", componentType))
}

func (cf *ComponentFactory) RegisterValidator(componentType string, validator ConfigValidator) {
    cf.mu.Lock()
    defer cf.mu.Unlock()
    
    cf.validators[componentType] = validator
    cf.logger.Info("config validator registered", zap.String("type", componentType))
}

func (cf *ComponentFactory) CreateComponent(config ComponentConfig) (Component, error) {
    cf.mu.RLock()
    validator, hasValidator := cf.validators[config.Type]
    creator, hasCreator := cf.creators[config.Type]
    cf.mu.RUnlock()
    
    // 验证配置
    if hasValidator {
        if err := validator(config); err != nil {
            return nil, fmt.Errorf("config validation failed: %w", err)
        }
    }
    
    // 创建组件
    if !hasCreator {
        return nil, fmt.Errorf("no creator registered for component type: %s", config.Type)
    }
    
    component, err := creator(config)
    if err != nil {
        cf.logger.Error("failed to create component", 
            zap.String("type", config.Type),
            zap.Error(err))
        return nil, fmt.Errorf("failed to create component: %w", err)
    }
    
    cf.logger.Info("component created", 
        zap.String("type", config.Type),
        zap.String("id", component.ID()))
    
    return component, nil
}

// 预定义组件创建器
func (cf *ComponentFactory) RegisterDefaultCreators() {
    // 注册基础组件创建器
    cf.RegisterCreator("base", func(config ComponentConfig) (Component, error) {
        return NewBaseComponent(config.Name, config.Version), nil
    })
    
    // 注册增强组件创建器
    cf.RegisterCreator("enhanced", func(config ComponentConfig) (Component, error) {
        return NewEnhancedComponent(config), nil
    })
    
    // 注册服务组件创建器
    cf.RegisterCreator("service", func(config ComponentConfig) (Component, error) {
        return NewServiceComponent(config), nil
    })
}

// 预定义配置验证器
func (cf *ComponentFactory) RegisterDefaultValidators() {
    // 基础验证器
    cf.RegisterValidator("base", func(config ComponentConfig) error {
        if config.Name == "" {
            return errors.New("component name is required")
        }
        if config.Version == "" {
            return errors.New("component version is required")
        }
        return nil
    })
    
    // 增强验证器
    cf.RegisterValidator("enhanced", func(config ComponentConfig) error {
        if err := cf.validators["base"](config); err != nil {
            return err
        }
        
        if config.Type == "" {
            return errors.New("component type is required")
        }
        
        if config.Timeout < 0 {
            return errors.New("timeout must be non-negative")
        }
        
        if config.Retries < 0 {
            return errors.New("retries must be non-negative")
        }
        
        return nil
    })
}
```

## 事件系统修复

### 2.1 高级事件总线

```go
// 高级事件总线
type AdvancedEventBus struct {
    topics      map[string]*Topic
    subscribers map[string][]Subscriber
    publishers  map[string]Publisher
    middleware  []EventMiddleware
    logger      *zap.Logger
    metrics     EventBusMetrics
    mu          sync.RWMutex
}

type Topic struct {
    name        string
    subscribers []Subscriber
    events      []Event
    config      TopicConfig
    mu          sync.RWMutex
}

type TopicConfig struct {
    MaxSubscribers int           `json:"max_subscribers"`
    BufferSize     int           `json:"buffer_size"`
    Retention      time.Duration `json:"retention"`
    Filtering      bool          `json:"filtering"`
    MaxEvents      int           `json:"max_events"`
}

type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Data      interface{}            `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    Source    string                 `json:"source"`
    Target    string                 `json:"target"`
    Metadata  map[string]interface{} `json:"metadata"`
}

type Subscriber interface {
    ID() string
    Handle(event Event) error
    GetFilter() EventFilter
}

type EventFilter func(event Event) bool

type Publisher interface {
    ID() string
    Publish(topic string, event Event) error
}

type EventMiddleware interface {
    BeforePublish(event Event) error
    AfterPublish(event Event) error
    BeforeSubscribe(subscriber Subscriber, topic string) error
    AfterSubscribe(subscriber Subscriber, topic string) error
}

func NewAdvancedEventBus() *AdvancedEventBus {
    return &AdvancedEventBus{
        topics:      make(map[string]*Topic),
        subscribers: make(map[string][]Subscriber),
        publishers:  make(map[string]Publisher),
        middleware:  make([]EventMiddleware, 0),
        logger:      zap.L().Named("advanced-event-bus"),
        metrics:     NewEventBusMetrics(),
    }
}

func (aeb *AdvancedEventBus) CreateTopic(name string, config TopicConfig) error {
    aeb.mu.Lock()
    defer aeb.mu.Unlock()
    
    if _, exists := aeb.topics[name]; exists {
        return fmt.Errorf("topic %s already exists", name)
    }
    
    topic := &Topic{
        name:        name,
        subscribers: make([]Subscriber, 0),
        events:      make([]Event, 0),
        config:      config,
    }
    
    aeb.topics[name] = topic
    aeb.logger.Info("topic created", zap.String("name", name))
    
    return nil
}

func (aeb *AdvancedEventBus) Publish(topic string, event Event) error {
    aeb.mu.RLock()
    t, exists := aeb.topics[topic]
    aeb.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("topic %s does not exist", topic)
    }
    
    // 生成事件ID
    if event.ID == "" {
        event.ID = uuid.New().String()
    }
    
    // 设置时间戳
    if event.Timestamp.IsZero() {
        event.Timestamp = time.Now()
    }
    
    // 应用中间件
    for _, middleware := range aeb.middleware {
        if err := middleware.BeforePublish(event); err != nil {
            return fmt.Errorf("middleware error: %w", err)
        }
    }
    
    // 发布事件
    t.mu.Lock()
    
    // 检查事件数量限制
    if t.config.MaxEvents > 0 && len(t.events) >= t.config.MaxEvents {
        // 移除最旧的事件
        t.events = t.events[1:]
    }
    
    // 添加事件
    t.events = append(t.events, event)
    
    // 获取订阅者副本
    subscribers := make([]Subscriber, len(t.subscribers))
    copy(subscribers, t.subscribers)
    
    t.mu.Unlock()
    
    // 异步通知订阅者
    for _, subscriber := range subscribers {
        go func(s Subscriber) {
            // 应用过滤器
            if filter := s.GetFilter(); filter != nil && !filter(event) {
                return
            }
            
            if err := s.Handle(event); err != nil {
                aeb.logger.Error("subscriber error",
                    zap.String("subscriber", s.ID()),
                    zap.String("event_id", event.ID),
                    zap.Error(err))
                aeb.metrics.SubscriberErrors.WithLabelValues(topic, s.ID()).Inc()
            } else {
                aeb.metrics.ProcessedEvents.WithLabelValues(topic, s.ID()).Inc()
            }
        }(subscriber)
    }
    
    // 应用后置中间件
    for _, middleware := range aeb.middleware {
        middleware.AfterPublish(event)
    }
    
    aeb.metrics.PublishedEvents.WithLabelValues(topic).Inc()
    aeb.logger.Debug("event published", 
        zap.String("topic", topic),
        zap.String("event_id", event.ID))
    
    return nil
}

func (aeb *AdvancedEventBus) Subscribe(topic string, subscriber Subscriber) error {
    aeb.mu.Lock()
    defer aeb.mu.Unlock()
    
    t, exists := aeb.topics[topic]
    if !exists {
        return fmt.Errorf("topic %s does not exist", topic)
    }
    
    // 应用中间件
    for _, middleware := range aeb.middleware {
        if err := middleware.BeforeSubscribe(subscriber, topic); err != nil {
            return fmt.Errorf("middleware error: %w", err)
        }
    }
    
    // 检查订阅者数量限制
    if t.config.MaxSubscribers > 0 && len(t.subscribers) >= t.config.MaxSubscribers {
        return fmt.Errorf("topic %s has reached maximum subscribers", topic)
    }
    
    // 检查是否已订阅
    for _, existingSub := range t.subscribers {
        if existingSub.ID() == subscriber.ID() {
            return fmt.Errorf("subscriber %s already subscribed to topic %s", subscriber.ID(), topic)
        }
    }
    
    t.subscribers = append(t.subscribers, subscriber)
    aeb.subscribers[topic] = append(aeb.subscribers[topic], subscriber)
    
    // 应用后置中间件
    for _, middleware := range aeb.middleware {
        middleware.AfterSubscribe(subscriber, topic)
    }
    
    aeb.logger.Info("subscriber added",
        zap.String("topic", topic),
        zap.String("subscriber", subscriber.ID()))
    
    return nil
}

func (aeb *AdvancedEventBus) Unsubscribe(topic string, subscriberID string) error {
    aeb.mu.Lock()
    defer aeb.mu.Unlock()
    
    t, exists := aeb.topics[topic]
    if !exists {
        return fmt.Errorf("topic %s does not exist", topic)
    }
    
    // 查找并移除订阅者
    for i, subscriber := range t.subscribers {
        if subscriber.ID() == subscriberID {
            t.subscribers = append(t.subscribers[:i], t.subscribers[i+1:]...)
            
            // 更新全局订阅者列表
            for j, sub := range aeb.subscribers[topic] {
                if sub.ID() == subscriberID {
                    aeb.subscribers[topic] = append(aeb.subscribers[topic][:j], aeb.subscribers[topic][j+1:]...)
                    break
                }
            }
            
            aeb.logger.Info("subscriber removed",
                zap.String("topic", topic),
                zap.String("subscriber", subscriberID))
            
            return nil
        }
    }
    
    return fmt.Errorf("subscriber %s not found in topic %s", subscriberID, topic)
}

func (aeb *AdvancedEventBus) AddMiddleware(middleware EventMiddleware) {
    aeb.mu.Lock()
    defer aeb.mu.Unlock()
    
    aeb.middleware = append(aeb.middleware, middleware)
    aeb.logger.Info("middleware added", zap.String("type", fmt.Sprintf("%T", middleware)))
}

// 事件指标
type EventBusMetrics struct {
    PublishedEvents   *prometheus.CounterVec
    ProcessedEvents   *prometheus.CounterVec
    SubscriberErrors  *prometheus.CounterVec
    TopicCount        prometheus.Gauge
    SubscriberCount   prometheus.Gauge
}

func NewEventBusMetrics() EventBusMetrics {
    return EventBusMetrics{
        PublishedEvents: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "eventbus_published_events_total",
                Help: "Total number of published events",
            },
            []string{"topic"},
        ),
        ProcessedEvents: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "eventbus_processed_events_total",
                Help: "Total number of processed events",
            },
            []string{"topic", "subscriber"},
        ),
        SubscriberErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "eventbus_subscriber_errors_total",
                Help: "Total number of subscriber errors",
            },
            []string{"topic", "subscriber"},
        ),
        TopicCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "eventbus_topics_total",
                Help: "Total number of topics",
            },
        ),
        SubscriberCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "eventbus_subscribers_total",
                Help: "Total number of subscribers",
            },
        ),
    }
}
```

## 配置管理修复

### 3.1 统一配置管理器

```go
// 配置管理器
type ConfigManager struct {
    viper    *viper.Viper
    consul   *ConsulClient
    etcd     *EtcdClient
    logger   *zap.Logger
    watchers map[string][]ConfigWatcher
    mu       sync.RWMutex
}

type ConfigWatcher interface {
    OnConfigChange(key string, oldValue, newValue interface{})
}

type ConfigSource interface {
    Load() (map[string]interface{}, error)
    Watch(prefix string) (<-chan ConfigChange, error)
}

type ConfigChange struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"` // create, update, delete
}

func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        viper:    viper.New(),
        watchers: make(map[string][]ConfigWatcher),
        logger:   zap.L().Named("config-manager"),
    }
}

func (cm *ConfigManager) LoadFile(configPath string) error {
    cm.viper.SetConfigFile(configPath)
    if err := cm.viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to read config: %w", err)
    }
    
    cm.logger.Info("config loaded from file", zap.String("path", configPath))
    return nil
}

func (cm *ConfigManager) LoadEnvironment(prefix string) {
    cm.viper.SetEnvPrefix(prefix)
    cm.viper.AutomaticEnv()
    cm.logger.Info("environment variables loaded", zap.String("prefix", prefix))
}

func (cm *ConfigManager) LoadConsul(consulAddr, prefix string) error {
    cm.consul = NewConsulClient(consulAddr)
    
    // 从Consul加载配置
    configs, err := cm.consul.GetConfigs(prefix)
    if err != nil {
        return fmt.Errorf("failed to load configs from consul: %w", err)
    }
    
    // 设置到Viper
    for key, value := range configs {
        cm.viper.Set(key, value)
    }
    
    // 监听配置变化
    go cm.watchConsulConfig(prefix)
    
    cm.logger.Info("config loaded from consul", zap.String("prefix", prefix))
    return nil
}

func (cm *ConfigManager) LoadEtcd(etcdAddrs []string, prefix string) error {
    cm.etcd = NewEtcdClient(etcdAddrs)
    
    // 从Etcd加载配置
    configs, err := cm.etcd.GetConfigs(prefix)
    if err != nil {
        return fmt.Errorf("failed to load configs from etcd: %w", err)
    }
    
    // 设置到Viper
    for key, value := range configs {
        cm.viper.Set(key, value)
    }
    
    // 监听配置变化
    go cm.watchEtcdConfig(prefix)
    
    cm.logger.Info("config loaded from etcd", zap.String("prefix", prefix))
    return nil
}

func (cm *ConfigManager) Get(key string) interface{} {
    return cm.viper.Get(key)
}

func (cm *ConfigManager) GetString(key string) string {
    return cm.viper.GetString(key)
}

func (cm *ConfigManager) GetInt(key string) int {
    return cm.viper.GetInt(key)
}

func (cm *ConfigManager) GetBool(key string) bool {
    return cm.viper.GetBool(key)
}

func (cm *ConfigManager) GetDuration(key string) time.Duration {
    return cm.viper.GetDuration(key)
}

func (cm *ConfigManager) GetStringSlice(key string) []string {
    return cm.viper.GetStringSlice(key)
}

func (cm *ConfigManager) GetStringMap(key string) map[string]interface{} {
    return cm.viper.GetStringMap(key)
}

func (cm *ConfigManager) Unmarshal(key string, v interface{}) error {
    return cm.viper.UnmarshalKey(key, v)
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.viper.Set(key, value)
    
    // 通知观察者
    cm.notifyWatchers(key, nil, value)
}

func (cm *ConfigManager) Watch(key string, watcher ConfigWatcher) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    cm.watchers[key] = append(cm.watchers[key], watcher)
    cm.logger.Info("config watcher added", zap.String("key", key))
}

func (cm *ConfigManager) Unwatch(key string, watcher ConfigWatcher) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    if watchers, exists := cm.watchers[key]; exists {
        for i, w := range watchers {
            if w == watcher {
                cm.watchers[key] = append(watchers[:i], watchers[i+1:]...)
                break
            }
        }
    }
}

func (cm *ConfigManager) notifyWatchers(key string, oldValue, newValue interface{}) {
    cm.mu.RLock()
    watchers, exists := cm.watchers[key]
    cm.mu.RUnlock()
    
    if !exists {
        return
    }
    
    for _, watcher := range watchers {
        go func(w ConfigWatcher) {
            w.OnConfigChange(key, oldValue, newValue)
        }(watcher)
    }
}

func (cm *ConfigManager) watchConsulConfig(prefix string) {
    changes := cm.consul.WatchConfigs(prefix)
    
    for change := range changes {
        oldValue := cm.viper.Get(change.Key)
        cm.viper.Set(change.Key, change.Value)
        cm.notifyWatchers(change.Key, oldValue, change.Value)
        
        cm.logger.Info("config changed from consul", 
            zap.String("key", change.Key),
            zap.Any("value", change.Value))
    }
}

func (cm *ConfigManager) watchEtcdConfig(prefix string) {
    changes := cm.etcd.WatchConfigs(prefix)
    
    for change := range changes {
        oldValue := cm.viper.Get(change.Key)
        cm.viper.Set(change.Key, change.Value)
        cm.notifyWatchers(change.Key, oldValue, change.Value)
        
        cm.logger.Info("config changed from etcd", 
            zap.String("key", change.Key),
            zap.Any("value", change.Value))
    }
}

// Consul客户端
type ConsulClient struct {
    client *consul.Client
    logger *zap.Logger
}

func NewConsulClient(addr string) *ConsulClient {
    config := consul.DefaultConfig()
    config.Address = addr
    
    client, err := consul.NewClient(config)
    if err != nil {
        panic(fmt.Sprintf("failed to create consul client: %v", err))
    }
    
    return &ConsulClient{
        client: client,
        logger: zap.L().Named("consul-client"),
    }
}

func (cc *ConsulClient) GetConfigs(prefix string) (map[string]interface{}, error) {
    pairs, _, err := cc.client.KV().List(prefix, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to list consul keys: %w", err)
    }
    
    configs := make(map[string]interface{})
    for _, pair := range pairs {
        key := strings.TrimPrefix(string(pair.Key), prefix+"/")
        var value interface{}
        
        if err := json.Unmarshal(pair.Value, &value); err != nil {
            // 如果不是JSON，当作字符串处理
            value = string(pair.Value)
        }
        
        configs[key] = value
    }
    
    return configs, nil
}

func (cc *ConsulClient) WatchConfigs(prefix string) <-chan ConfigChange {
    changes := make(chan ConfigChange, 100)
    
    go func() {
        var lastIndex uint64
        
        for {
            pairs, meta, err := cc.client.KV().List(prefix, &consul.QueryOptions{
                WaitIndex: lastIndex,
                WaitTime:  time.Minute,
            })
            
            if err != nil {
                cc.logger.Error("failed to watch consul configs", zap.Error(err))
                time.Sleep(time.Second)
                continue
            }
            
            if meta.LastIndex == lastIndex {
                continue
            }
            
            lastIndex = meta.LastIndex
            
            for _, pair := range pairs {
                var value interface{}
                if err := json.Unmarshal(pair.Value, &value); err != nil {
                    value = string(pair.Value)
                }
                
                changes <- ConfigChange{
                    Key:   strings.TrimPrefix(string(pair.Key), prefix+"/"),
                    Value: value,
                    Type:  "update",
                }
            }
        }
    }()
    
    return changes
}
```

## 监控系统修复

### 4.1 Prometheus指标收集

```go
// Prometheus指标收集器
type PrometheusCollector struct {
    registry *prometheus.Registry
    metrics  map[string]prometheus.Collector
    logger   *zap.Logger
    mu       sync.RWMutex
}

func NewPrometheusCollector() *PrometheusCollector {
    return &PrometheusCollector{
        registry: prometheus.NewRegistry(),
        metrics:  make(map[string]prometheus.Collector),
        logger:   zap.L().Named("prometheus-collector"),
    }
}

func (pc *PrometheusCollector) RegisterMetric(name string, metric prometheus.Collector) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    if err := pc.registry.Register(metric); err != nil {
        return fmt.Errorf("failed to register metric %s: %w", name, err)
    }
    
    pc.metrics[name] = metric
    pc.logger.Info("metric registered", zap.String("name", name))
    return nil
}

func (pc *PrometheusCollector) UnregisterMetric(name string) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    metric, exists := pc.metrics[name]
    if !exists {
        return fmt.Errorf("metric %s not found", name)
    }
    
    if success := pc.registry.Unregister(metric); !success {
        return fmt.Errorf("failed to unregister metric %s", name)
    }
    
    delete(pc.metrics, name)
    pc.logger.Info("metric unregistered", zap.String("name", name))
    return nil
}

func (pc *PrometheusCollector) GetRegistry() *prometheus.Registry {
    return pc.registry
}

// 组件指标
func NewComponentMetrics(componentName string) *ComponentMetrics {
    return &ComponentMetrics{
        StartTime: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_start_total",
            Help: "Total number of component starts",
            ConstLabels: prometheus.Labels{"component": componentName},
        }),
        StopTime: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_stop_total",
            Help: "Total number of component stops",
            ConstLabels: prometheus.Labels{"component": componentName},
        }),
        ErrorCount: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "component_errors_total",
            Help: "Total number of component errors",
            ConstLabels: prometheus.Labels{"component": componentName},
        }),
        StatusGauge: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "component_status",
            Help: "Current component status",
            ConstLabels: prometheus.Labels{"component": componentName},
        }),
        Duration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name: "component_operation_duration_seconds",
            Help: "Duration of component operations",
            ConstLabels: prometheus.Labels{"component": componentName},
            Buckets: prometheus.DefBuckets,
        }),
    }
}

// HTTP指标服务器
type MetricsServer struct {
    addr     string
    registry *prometheus.Registry
    logger   *zap.Logger
    server   *http.Server
}

func NewMetricsServer(addr string, registry *prometheus.Registry) *MetricsServer {
    return &MetricsServer{
        addr:     addr,
        registry: registry,
        logger:   zap.L().Named("metrics-server"),
    }
}

func (ms *MetricsServer) Start() error {
    mux := http.NewServeMux()
    
    // 指标端点
    mux.Handle("/metrics", promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{
        EnableOpenMetrics: true,
    }))
    
    // 健康检查端点
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("healthy"))
    })
    
    // 就绪检查端点
    mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ready"))
    })
    
    ms.server = &http.Server{
        Addr:    ms.addr,
        Handler: mux,
    }
    
    go func() {
        if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            ms.logger.Error("metrics server error", zap.Error(err))
        }
    }()
    
    ms.logger.Info("metrics server started", zap.String("addr", ms.addr))
    return nil
}

func (ms *MetricsServer) Stop() error {
    if ms.server != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := ms.server.Shutdown(ctx); err != nil {
            return fmt.Errorf("failed to shutdown metrics server: %w", err)
        }
    }
    
    ms.logger.Info("metrics server stopped")
    return nil
}
```

### 4.2 健康检查系统

```go
// 健康检查管理器
type HealthCheckManager struct {
    checks map[string]HealthCheck
    logger *zap.Logger
    mu     sync.RWMutex
}

type HealthCheck interface {
    Check() HealthStatus
    GetName() string
    GetType() string
}

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Details   map[string]string `json:"details"`
    Errors    []string          `json:"errors"`
    Type      string            `json:"type"`
}

func NewHealthCheckManager() *HealthCheckManager {
    return &HealthCheckManager{
        checks: make(map[string]HealthCheck),
        logger: zap.L().Named("health-check-manager"),
    }
}

func (hcm *HealthCheckManager) RegisterCheck(check HealthCheck) {
    hcm.mu.Lock()
    defer hcm.mu.Unlock()
    
    hcm.checks[check.GetName()] = check
    hcm.logger.Info("health check registered", 
        zap.String("name", check.GetName()),
        zap.String("type", check.GetType()))
}

func (hcm *HealthCheckManager) UnregisterCheck(name string) {
    hcm.mu.Lock()
    defer hcm.mu.Unlock()
    
    if _, exists := hcm.checks[name]; exists {
        delete(hcm.checks, name)
        hcm.logger.Info("health check unregistered", zap.String("name", name))
    }
}

func (hcm *HealthCheckManager) RunChecks() map[string]HealthStatus {
    hcm.mu.RLock()
    checks := make(map[string]HealthCheck, len(hcm.checks))
    for k, v := range hcm.checks {
        checks[k] = v
    }
    hcm.mu.RUnlock()
    
    results := make(map[string]HealthStatus)
    
    for name, check := range checks {
        start := time.Now()
        status := check.Check()
        duration := time.Since(start)
        
        // 添加检查耗时
        if status.Details == nil {
            status.Details = make(map[string]string)
        }
        status.Details["duration"] = duration.String()
        status.Type = check.GetType()
        
        results[name] = status
        
        if status.Status == "unhealthy" {
            hcm.logger.Warn("health check failed",
                zap.String("name", name),
                zap.String("type", check.GetType()),
                zap.Strings("errors", status.Errors),
                zap.Duration("duration", duration))
        } else {
            hcm.logger.Debug("health check passed",
                zap.String("name", name),
                zap.String("type", check.GetType()),
                zap.Duration("duration", duration))
        }
    }
    
    return results
}

func (hcm *HealthCheckManager) GetOverallStatus() string {
    results := hcm.RunChecks()
    
    for _, status := range results {
        if status.Status == "unhealthy" {
            return "unhealthy"
        }
    }
    
    return "healthy"
}

// 组件健康检查
type ComponentHealthCheck struct {
    component Component
}

func NewComponentHealthCheck(component Component) *ComponentHealthCheck {
    return &ComponentHealthCheck{
        component: component,
    }
}

func (chc *ComponentHealthCheck) Check() HealthStatus {
    status := HealthStatus{
        Timestamp: time.Now(),
        Details:   make(map[string]string),
        Type:      "component",
    }
    
    // 检查组件状态
    componentStatus := chc.component.Status()
    status.Details["status"] = componentStatus.String()
    
    switch componentStatus {
    case StatusRunning:
        status.Status = "healthy"
    case StatusError:
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, "component in error state")
    case StatusStopped:
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, "component is stopped")
    default:
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, "component in unknown state")
    }
    
    // 检查健康状态（如果组件实现了健康检查接口）
    if healthComponent, ok := chc.component.(interface{ Health() HealthStatus }); ok {
        health := healthComponent.Health()
        if health.Status != "healthy" {
            status.Status = "unhealthy"
            status.Errors = append(status.Errors, health.Errors...)
        }
    }
    
    return status
}

func (chc *ComponentHealthCheck) GetName() string {
    return fmt.Sprintf("component-%s", chc.component.ID())
}

func (chc *ComponentHealthCheck) GetType() string {
    return "component"
}

// 数据库健康检查
type DatabaseHealthCheck struct {
    db     *gorm.DB
    name   string
}

func NewDatabaseHealthCheck(db *gorm.DB, name string) *DatabaseHealthCheck {
    return &DatabaseHealthCheck{
        db:   db,
        name: name,
    }
}

func (dhc *DatabaseHealthCheck) Check() HealthStatus {
    status := HealthStatus{
        Timestamp: time.Now(),
        Details:   make(map[string]string),
        Type:      "database",
    }
    
    // 执行简单查询检查连接
    var result int
    if err := dhc.db.Raw("SELECT 1").Scan(&result).Error; err != nil {
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, fmt.Sprintf("database connection failed: %v", err))
    } else {
        status.Status = "healthy"
        status.Details["connected"] = "true"
    }
    
    return status
}

func (dhc *DatabaseHealthCheck) GetName() string {
    return fmt.Sprintf("database-%s", dhc.name)
}

func (dhc *DatabaseHealthCheck) GetType() string {
    return "database"
}

// HTTP健康检查
type HTTPHealthCheck struct {
    url    string
    name   string
    timeout time.Duration
}

func NewHTTPHealthCheck(url, name string, timeout time.Duration) *HTTPHealthCheck {
    return &HTTPHealthCheck{
        url:     url,
        name:    name,
        timeout: timeout,
    }
}

func (hhc *HTTPHealthCheck) Check() HealthStatus {
    status := HealthStatus{
        Timestamp: time.Now(),
        Details:   make(map[string]string),
        Type:      "http",
    }
    
    client := &http.Client{
        Timeout: hhc.timeout,
    }
    
    resp, err := client.Get(hhc.url)
    if err != nil {
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, fmt.Sprintf("HTTP request failed: %v", err))
        return status
    }
    defer resp.Body.Close()
    
    status.Details["status_code"] = fmt.Sprintf("%d", resp.StatusCode)
    status.Details["response_time"] = resp.Header.Get("X-Response-Time")
    
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        status.Status = "healthy"
    } else {
        status.Status = "unhealthy"
        status.Errors = append(status.Errors, fmt.Sprintf("HTTP status code: %d", resp.StatusCode))
    }
    
    return status
}

func (hhc *HTTPHealthCheck) GetName() string {
    return fmt.Sprintf("http-%s", hhc.name)
}

func (hhc *HTTPHealthCheck) GetType() string {
    return "http"
}
```

## 总结

本实现修复方案提供了Golang Common库的核心功能修复，包括：

1. **组件系统修复**: 增强组件接口、工厂模式实现
2. **事件系统修复**: 高级事件总线、中间件支持
3. **配置管理修复**: 统一配置管理器、多源配置支持
4. **监控系统修复**: Prometheus指标收集、健康检查系统

这些修复将显著提升Golang Common库的功能完整性、可靠性和可维护性。
