package enhanced

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Component 组件接口
type Component interface {
	ID() string
	Kind() string
	Start(ctx context.Context) error
	Stop() error
	Status() ComponentStatus
	Dependencies() []string
	Health() HealthStatus
	Metrics() ComponentMetrics
}

// ComponentStatus 组件状态
type ComponentStatus int

const (
	ComponentStatusStopped ComponentStatus = iota
	ComponentStatusStarting
	ComponentStatusRunning
	ComponentStatusStopping
	ComponentStatusError
)

func (cs ComponentStatus) String() string {
	switch cs {
	case ComponentStatusStopped:
		return "Stopped"
	case ComponentStatusStarting:
		return "Starting"
	case ComponentStatusRunning:
		return "Running"
	case ComponentStatusStopping:
		return "Stopping"
	case ComponentStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

// ComponentMetrics 组件指标
type ComponentMetrics struct {
	StartTime    time.Time     `json:"start_time"`
	StopTime     time.Time     `json:"stop_time"`
	Uptime       time.Duration `json:"uptime"`
	RestartCount int64         `json:"restart_count"`
	ErrorCount   int64         `json:"error_count"`
	RequestCount int64         `json:"request_count"`
	ResponseTime time.Duration `json:"response_time"`
}

// BaseComponent 基础组件实现
type BaseComponent struct {
	id           string
	kind         string
	status       int32 // atomic
	dependencies []string
	container    *DependencyContainer
	lifecycle    *LifecycleManager
	logger       *zap.Logger
	health       HealthStatus
	metrics      ComponentMetrics
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBaseComponent 创建基础组件
func NewBaseComponent(id, kind string, dependencies []string, container *DependencyContainer) *BaseComponent {
	ctx, cancel := context.WithCancel(context.Background())

	return &BaseComponent{
		id:           id,
		kind:         kind,
		status:       int32(ComponentStatusStopped),
		dependencies: dependencies,
		container:    container,
		lifecycle:    NewLifecycleManager(ctx),
		logger:       zap.L().Named(fmt.Sprintf("component:%s", id)),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// ID 获取组件ID
func (bc *BaseComponent) ID() string {
	return bc.id
}

// Kind 获取组件类型
func (bc *BaseComponent) Kind() string {
	return bc.kind
}

// Status 获取组件状态
func (bc *BaseComponent) Status() ComponentStatus {
	return ComponentStatus(atomic.LoadInt32(&bc.status))
}

// Dependencies 获取依赖列表
func (bc *BaseComponent) Dependencies() []string {
	return bc.dependencies
}

// Health 获取健康状态
func (bc *BaseComponent) Health() HealthStatus {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.health
}

// Metrics 获取组件指标
func (bc *BaseComponent) Metrics() ComponentMetrics {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	metrics := bc.metrics
	if bc.Status() == ComponentStatusRunning {
		metrics.Uptime = time.Since(metrics.StartTime)
	}

	return metrics
}

// Start 启动组件
func (bc *BaseComponent) Start(ctx context.Context) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	currentStatus := ComponentStatus(atomic.LoadInt32(&bc.status))
	if currentStatus != ComponentStatusStopped {
		return fmt.Errorf("component %s is not in stopped status, current: %s", bc.id, currentStatus)
	}

	atomic.StoreInt32(&bc.status, int32(ComponentStatusStarting))
	bc.logger.Info("starting component")

	// 检查依赖
	if err := bc.checkDependencies(); err != nil {
		atomic.StoreInt32(&bc.status, int32(ComponentStatusError))
		bc.updateHealth("ERROR", fmt.Sprintf("Dependency check failed: %v", err), nil)
		return fmt.Errorf("dependency check failed: %w", err)
	}

	// 启动生命周期管理器
	if err := bc.lifecycle.Start(); err != nil {
		atomic.StoreInt32(&bc.status, int32(ComponentStatusError))
		bc.updateHealth("ERROR", fmt.Sprintf("Lifecycle start failed: %v", err), nil)
		return fmt.Errorf("failed to start lifecycle: %w", err)
	}

	// 更新指标
	bc.metrics.StartTime = time.Now()
	atomic.AddInt64(&bc.metrics.RestartCount, 1)

	atomic.StoreInt32(&bc.status, int32(ComponentStatusRunning))
	bc.updateHealth("HEALTHY", "Component is running", nil)
	bc.logger.Info("component started")

	return nil
}

// Stop 停止组件
func (bc *BaseComponent) Stop() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	currentStatus := ComponentStatus(atomic.LoadInt32(&bc.status))
	if currentStatus != ComponentStatusRunning {
		return fmt.Errorf("component %s is not running, current: %s", bc.id, currentStatus)
	}

	atomic.StoreInt32(&bc.status, int32(ComponentStatusStopping))
	bc.logger.Info("stopping component")

	// 停止生命周期管理器
	if err := bc.lifecycle.Stop(); err != nil {
		atomic.StoreInt32(&bc.status, int32(ComponentStatusError))
		bc.updateHealth("ERROR", fmt.Sprintf("Lifecycle stop failed: %v", err), nil)
		return fmt.Errorf("failed to stop lifecycle: %w", err)
	}

	// 更新指标
	bc.metrics.StopTime = time.Now()
	bc.metrics.Uptime = bc.metrics.StopTime.Sub(bc.metrics.StartTime)

	atomic.StoreInt32(&bc.status, int32(ComponentStatusStopped))
	bc.updateHealth("STOPPED", "Component is stopped", nil)
	bc.logger.Info("component stopped")

	return nil
}

// checkDependencies 检查依赖
func (bc *BaseComponent) checkDependencies() error {
	for _, dep := range bc.dependencies {
		if _, err := bc.container.GetService(dep); err != nil {
			return fmt.Errorf("dependency %s not available: %w", dep, err)
		}
	}
	return nil
}

// updateHealth 更新健康状态
func (bc *BaseComponent) updateHealth(status, message string, details map[string]interface{}) {
	bc.health = HealthStatus{
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
		Details:   details,
	}
}

// RecordError 记录错误
func (bc *BaseComponent) RecordError(err error) {
	atomic.AddInt64(&bc.metrics.ErrorCount, 1)
	bc.logger.Error("component error", zap.Error(err))

	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.updateHealth("ERROR", err.Error(), map[string]interface{}{
		"error_count": bc.metrics.ErrorCount,
	})
}

// RecordRequest 记录请求
func (bc *BaseComponent) RecordRequest(duration time.Duration) {
	atomic.AddInt64(&bc.metrics.RequestCount, 1)

	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.metrics.ResponseTime = duration
}

// AddWorker 添加工作器
func (bc *BaseComponent) AddWorker(name string, handler WorkerHandler) error {
	return bc.lifecycle.AddWorker(name, handler)
}

// RemoveWorker 移除工作器
func (bc *BaseComponent) RemoveWorker(name string) error {
	return bc.lifecycle.RemoveWorker(name)
}

// Context 获取组件上下文
func (bc *BaseComponent) Context() context.Context {
	return bc.ctx
}

// Cancel 取消组件上下文
func (bc *BaseComponent) Cancel() {
	bc.cancel()
}

// ComponentManager 组件管理器
type ComponentManager struct {
	components map[string]Component
	container  *DependencyContainer
	config     *ConfigManager
	logger     *zap.Logger
	mu         sync.RWMutex
}

// NewComponentManager 创建组件管理器
func NewComponentManager(container *DependencyContainer, config *ConfigManager) *ComponentManager {
	return &ComponentManager{
		components: make(map[string]Component),
		container:  container,
		config:     config,
		logger:     zap.L().Named("component-manager"),
	}
}

// RegisterComponent 注册组件
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

// UnregisterComponent 注销组件
func (cm *ComponentManager) UnregisterComponent(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.components[id]; !exists {
		return fmt.Errorf("component %s not found", id)
	}

	delete(cm.components, id)
	cm.logger.Info("component unregistered", zap.String("id", id))

	return nil
}

// GetComponent 获取组件
func (cm *ComponentManager) GetComponent(id string) (Component, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	component, exists := cm.components[id]
	if !exists {
		return nil, fmt.Errorf("component %s not found", id)
	}

	return component, nil
}

// GetComponentsByKind 按类型获取组件
func (cm *ComponentManager) GetComponentsByKind(kind string) []Component {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var components []Component
	for _, component := range cm.components {
		if component.Kind() == kind {
			components = append(components, component)
		}
	}

	return components
}

// ListComponents 列出所有组件
func (cm *ComponentManager) ListComponents() []Component {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	components := make([]Component, 0, len(cm.components))
	for _, component := range cm.components {
		components = append(components, component)
	}

	return components
}

// StartAll 启动所有组件
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

// StopAll 停止所有组件
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

// HealthCheck 健康检查
func (cm *ComponentManager) HealthCheck() map[string]HealthStatus {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	health := make(map[string]HealthStatus)
	for id, component := range cm.components {
		health[id] = component.Health()
	}

	return health
}

// Metrics 获取所有组件指标
func (cm *ComponentManager) Metrics() map[string]ComponentMetrics {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	metrics := make(map[string]ComponentMetrics)
	for id, component := range cm.components {
		metrics[id] = component.Metrics()
	}

	return metrics
}

// sortByDependencies 按依赖关系排序
func (cm *ComponentManager) sortByDependencies() ([]Component, error) {
	components := make([]Component, 0, len(cm.components))
	for _, component := range cm.components {
		components = append(components, component)
	}

	// 拓扑排序
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
