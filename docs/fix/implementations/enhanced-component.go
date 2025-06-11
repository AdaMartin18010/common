package implementations

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 增强组件接口
type EnhancedComponent interface {
	Component
	Health() HealthStatus
	Metrics() ComponentMetrics
	Configuration() ComponentConfig
	Dependencies() []string
	Validate() error
}

// 基础组件接口
type Component interface {
	ID() string
	Name() string
	Version() string
	Status() ComponentStatus
	Start() error
	Stop() error
}

// 组件状态
type ComponentStatus int

const (
	StatusCreated ComponentStatus = iota
	StatusInitialized
	StatusStarted
	StatusRunning
	StatusStopping
	StatusStopped
	StatusError
)

func (cs ComponentStatus) String() string {
	switch cs {
	case StatusCreated:
		return "created"
	case StatusInitialized:
		return "initialized"
	case StatusStarted:
		return "started"
	case StatusRunning:
		return "running"
	case StatusStopping:
		return "stopping"
	case StatusStopped:
		return "stopped"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
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

// 组件指标
type ComponentMetrics struct {
	StartTime   int64
	StopTime    int64
	ErrorCount  int64
	StatusGauge int
	Duration    time.Duration
}

// 基础组件实现
type BaseComponent struct {
	id      string
	name    string
	version string
	status  ComponentStatus
	logger  *zap.Logger
	mu      sync.RWMutex
}

func NewBaseComponent(name, version string) *BaseComponent {
	return &BaseComponent{
		id:      uuid.New().String(),
		name:    name,
		version: version,
		status:  StatusCreated,
		logger:  zap.L().Named(name),
	}
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

func (bc *BaseComponent) Status() ComponentStatus {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.status
}

func (bc *BaseComponent) setStatus(status ComponentStatus) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.status = status
	bc.logger.Debug("status changed", zap.String("status", status.String()))
}

func (bc *BaseComponent) Start() error {
	bc.setStatus(StatusStarted)
	bc.logger.Info("base component started")
	return nil
}

func (bc *BaseComponent) Stop() error {
	bc.setStatus(StatusStopping)
	bc.setStatus(StatusStopped)
	bc.logger.Info("base component stopped")
	return nil
}

// 增强组件实现
type EnhancedComponentImpl struct {
	*BaseComponent
	health       HealthStatus
	metrics      *ComponentMetrics
	config       ComponentConfig
	dependencies []string
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewEnhancedComponent(config ComponentConfig) *EnhancedComponentImpl {
	ctx, cancel := context.WithCancel(context.Background())

	return &EnhancedComponentImpl{
		BaseComponent: NewBaseComponent(config.Name, config.Version),
		config:        config,
		dependencies:  config.Dependencies,
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

func (ec *EnhancedComponentImpl) Metrics() ComponentMetrics {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	if ec.metrics == nil {
		return ComponentMetrics{}
	}
	return *ec.metrics
}

func (ec *EnhancedComponentImpl) Configuration() ComponentConfig {
	return ec.config
}

func (ec *EnhancedComponentImpl) Dependencies() []string {
	return ec.dependencies
}

func (ec *EnhancedComponentImpl) Validate() error {
	if ec.config.Name == "" {
		return errors.New("component name is required")
	}
	if ec.config.Version == "" {
		return errors.New("component version is required")
	}
	if ec.config.Type == "" {
		return errors.New("component type is required")
	}
	if ec.config.Timeout < 0 {
		return errors.New("timeout must be non-negative")
	}
	if ec.config.Retries < 0 {
		return errors.New("retries must be non-negative")
	}
	return nil
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
			ec.metrics.StatusGauge = 1
		case "error":
			ec.metrics.StatusGauge = 0
			ec.metrics.ErrorCount++
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
	ec.metrics = &ComponentMetrics{
		StartTime: time.Now().Unix(),
	}
	return nil
}

func (ec *EnhancedComponentImpl) stopMetrics() {
	if ec.metrics != nil {
		ec.metrics.StopTime = time.Now().Unix()
		ec.metrics.Duration = time.Duration(ec.metrics.StopTime-ec.metrics.StartTime) * time.Second
	}
}
