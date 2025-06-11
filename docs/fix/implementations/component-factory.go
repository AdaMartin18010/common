package implementations

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

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

	// 注册数据库组件创建器
	cf.RegisterCreator("database", func(config ComponentConfig) (Component, error) {
		return NewDatabaseComponent(config), nil
	})

	// 注册缓存组件创建器
	cf.RegisterCreator("cache", func(config ComponentConfig) (Component, error) {
		return NewCacheComponent(config), nil
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

	// 服务验证器
	cf.RegisterValidator("service", func(config ComponentConfig) error {
		if err := cf.validators["enhanced"](config); err != nil {
			return err
		}

		// 检查服务特定配置
		if port, ok := config.Settings["port"].(int); ok {
			if port <= 0 || port > 65535 {
				return errors.New("invalid port number")
			}
		}

		return nil
	})

	// 数据库验证器
	cf.RegisterValidator("database", func(config ComponentConfig) error {
		if err := cf.validators["enhanced"](config); err != nil {
			return err
		}

		// 检查数据库连接配置
		if dsn, ok := config.Settings["dsn"].(string); ok {
			if dsn == "" {
				return errors.New("database DSN is required")
			}
		} else {
			return errors.New("database DSN is required")
		}

		return nil
	})
}

// 服务组件
type ServiceComponent struct {
	*EnhancedComponentImpl
	port   int
	server interface{} // 实际的HTTP服务器
}

func NewServiceComponent(config ComponentConfig) *ServiceComponent {
	port := 8080
	if p, ok := config.Settings["port"].(int); ok {
		port = p
	}

	return &ServiceComponent{
		EnhancedComponentImpl: NewEnhancedComponent(config),
		port:                  port,
	}
}

func (sc *ServiceComponent) Start() error {
	// 启动HTTP服务器
	sc.logger.Info("starting service component", zap.Int("port", sc.port))

	// 这里应该启动实际的HTTP服务器
	// 为了示例，我们只是记录日志

	return sc.EnhancedComponentImpl.Start()
}

func (sc *ServiceComponent) Stop() error {
	// 停止HTTP服务器
	sc.logger.Info("stopping service component")

	// 这里应该停止实际的HTTP服务器

	return sc.EnhancedComponentImpl.Stop()
}

// 数据库组件
type DatabaseComponent struct {
	*EnhancedComponentImpl
	dsn string
	db  interface{} // 实际的数据库连接
}

func NewDatabaseComponent(config ComponentConfig) *DatabaseComponent {
	dsn := ""
	if d, ok := config.Settings["dsn"].(string); ok {
		dsn = d
	}

	return &DatabaseComponent{
		EnhancedComponentImpl: NewEnhancedComponent(config),
		dsn:                   dsn,
	}
}

func (dc *DatabaseComponent) Start() error {
	// 连接数据库
	dc.logger.Info("connecting to database", zap.String("dsn", dc.dsn))

	// 这里应该建立实际的数据库连接
	// 为了示例，我们只是记录日志

	return dc.EnhancedComponentImpl.Start()
}

func (dc *DatabaseComponent) Stop() error {
	// 关闭数据库连接
	dc.logger.Info("closing database connection")

	// 这里应该关闭实际的数据库连接

	return dc.EnhancedComponentImpl.Stop()
}

// 缓存组件
type CacheComponent struct {
	*EnhancedComponentImpl
	redisAddr string
	cache     interface{} // 实际的缓存客户端
}

func NewCacheComponent(config ComponentConfig) *CacheComponent {
	redisAddr := "localhost:6379"
	if addr, ok := config.Settings["redis_addr"].(string); ok {
		redisAddr = addr
	}

	return &CacheComponent{
		EnhancedComponentImpl: NewEnhancedComponent(config),
		redisAddr:             redisAddr,
	}
}

func (cc *CacheComponent) Start() error {
	// 连接缓存
	cc.logger.Info("connecting to cache", zap.String("redis_addr", cc.redisAddr))

	// 这里应该建立实际的缓存连接
	// 为了示例，我们只是记录日志

	return cc.EnhancedComponentImpl.Start()
}

func (cc *CacheComponent) Stop() error {
	// 关闭缓存连接
	cc.logger.Info("closing cache connection")

	// 这里应该关闭实际的缓存连接

	return cc.EnhancedComponentImpl.Stop()
}

// 组件构建器
type ComponentBuilder struct {
	config ComponentConfig
}

func NewComponentBuilder() *ComponentBuilder {
	return &ComponentBuilder{
		config: ComponentConfig{
			Settings: make(map[string]interface{}),
		},
	}
}

func (cb *ComponentBuilder) SetName(name string) *ComponentBuilder {
	cb.config.Name = name
	return cb
}

func (cb *ComponentBuilder) SetVersion(version string) *ComponentBuilder {
	cb.config.Version = version
	return cb
}

func (cb *ComponentBuilder) SetType(componentType string) *ComponentBuilder {
	cb.config.Type = componentType
	return cb
}

func (cb *ComponentBuilder) SetDependencies(dependencies []string) *ComponentBuilder {
	cb.config.Dependencies = dependencies
	return cb
}

func (cb *ComponentBuilder) SetTimeout(timeout int) *ComponentBuilder {
	cb.config.Timeout = time.Duration(timeout) * time.Second
	return cb
}

func (cb *ComponentBuilder) SetRetries(retries int) *ComponentBuilder {
	cb.config.Retries = retries
	return cb
}

func (cb *ComponentBuilder) SetSetting(key string, value interface{}) *ComponentBuilder {
	cb.config.Settings[key] = value
	return cb
}

func (cb *ComponentBuilder) Build() ComponentConfig {
	return cb.config
}

// 组件管理器
type ComponentManager struct {
	factory    *ComponentFactory
	components map[string]Component
	logger     *zap.Logger
	mu         sync.RWMutex
}

func NewComponentManager() *ComponentManager {
	factory := NewComponentFactory()
	factory.RegisterDefaultCreators()
	factory.RegisterDefaultValidators()

	return &ComponentManager{
		factory:    factory,
		components: make(map[string]Component),
		logger:     zap.L().Named("component-manager"),
	}
}

func (cm *ComponentManager) CreateComponent(config ComponentConfig) (Component, error) {
	component, err := cm.factory.CreateComponent(config)
	if err != nil {
		return nil, err
	}

	cm.mu.Lock()
	cm.components[component.ID()] = component
	cm.mu.Unlock()

	cm.logger.Info("component created and registered",
		zap.String("id", component.ID()),
		zap.String("name", component.Name()))

	return component, nil
}

func (cm *ComponentManager) GetComponent(id string) (Component, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

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

func (cm *ComponentManager) StartAll() error {
	cm.mu.RLock()
	components := make([]Component, 0, len(cm.components))
	for _, component := range cm.components {
		components = append(components, component)
	}
	cm.mu.RUnlock()

	var errors []error
	for _, component := range components {
		if err := component.Start(); err != nil {
			errors = append(errors, fmt.Errorf("failed to start component %s: %w", component.ID(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start some components: %v", errors)
	}

	return nil
}

func (cm *ComponentManager) StopAll() error {
	cm.mu.RLock()
	components := make([]Component, 0, len(cm.components))
	for _, component := range cm.components {
		components = append(components, component)
	}
	cm.mu.RUnlock()

	var errors []error
	for _, component := range components {
		if err := component.Stop(); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop component %s: %w", component.ID(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some components: %v", errors)
	}

	return nil
}

func (cm *ComponentManager) ListComponents() []Component {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	components := make([]Component, 0, len(cm.components))
	for _, component := range cm.components {
		components = append(components, component)
	}

	return components
}

func (cm *ComponentManager) RemoveComponent(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.components[id]; !exists {
		return fmt.Errorf("component %s not found", id)
	}

	delete(cm.components, id)
	cm.logger.Info("component removed", zap.String("id", id))

	return nil
}
