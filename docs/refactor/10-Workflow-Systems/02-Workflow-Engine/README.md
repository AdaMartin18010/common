# 02. 工作流引擎设计 (Workflow Engine Design)

## 概述

工作流引擎是工作流系统的核心组件，负责工作流的定义、执行、监控和管理。本模块基于基础理论，设计高性能、可扩展的工作流引擎架构。

## 目录

- [01-引擎架构设计](./01-Engine-Architecture.md)
- [02-执行引擎实现](./02-Execution-Engine.md)
- [03-调度器设计](./03-Scheduler-Design.md)
- [04-持久化机制](./04-Persistence-Mechanism.md)

## 核心架构

### 1. 引擎整体架构

```go
package workflow

import (
    "context"
    "sync"
    "time"
    
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "github.com/streadway/amqp"
)

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
    config       *EngineConfig
    scheduler    *Scheduler
    executor     *Executor
    storage      *Storage
    eventBus     *EventBus
    metrics      *Metrics
    mutex        sync.RWMutex
    running      bool
}

// EngineConfig 引擎配置
type EngineConfig struct {
    MaxConcurrency    int           `json:"max_concurrency"`
    WorkerPoolSize    int           `json:"worker_pool_size"`
    TaskTimeout       time.Duration `json:"task_timeout"`
    RetryAttempts     int           `json:"retry_attempts"`
    RetryDelay        time.Duration `json:"retry_delay"`
    PersistenceMode   string        `json:"persistence_mode"`
    EventBusType      string        `json:"event_bus_type"`
    MetricsEnabled    bool          `json:"metrics_enabled"`
}

// NewWorkflowEngine 创建新的工作流引擎
func NewWorkflowEngine(config *EngineConfig) *WorkflowEngine {
    engine := &WorkflowEngine{
        config: config,
    }
    
    // 初始化组件
    engine.scheduler = NewScheduler(config)
    engine.executor = NewExecutor(config)
    engine.storage = NewStorage(config)
    engine.eventBus = NewEventBus(config)
    engine.metrics = NewMetrics(config)
    
    return engine
}

// Start 启动引擎
func (we *WorkflowEngine) Start(ctx context.Context) error {
    we.mutex.Lock()
    defer we.mutex.Unlock()
    
    if we.running {
        return fmt.Errorf("engine is already running")
    }
    
    // 启动各个组件
    if err := we.scheduler.Start(ctx); err != nil {
        return fmt.Errorf("failed to start scheduler: %w", err)
    }
    
    if err := we.executor.Start(ctx); err != nil {
        return fmt.Errorf("failed to start executor: %w", err)
    }
    
    if err := we.eventBus.Start(ctx); err != nil {
        return fmt.Errorf("failed to start event bus: %w", err)
    }
    
    we.running = true
    return nil
}

// Stop 停止引擎
func (we *WorkflowEngine) Stop(ctx context.Context) error {
    we.mutex.Lock()
    defer we.mutex.Unlock()
    
    if !we.running {
        return nil
    }
    
    // 停止各个组件
    we.scheduler.Stop(ctx)
    we.executor.Stop(ctx)
    we.eventBus.Stop(ctx)
    
    we.running = false
    return nil
}
```

### 2. 调度器设计

#### 2.1 调度器架构

```go
// Scheduler 调度器
type Scheduler struct {
    config       *EngineConfig
    taskQueue    chan *Task
    workers      []*Worker
    taskStore    TaskStore
    eventBus     *EventBus
    metrics      *Metrics
    mutex        sync.RWMutex
    running      bool
    ctx          context.Context
    cancel       context.CancelFunc
}

// Task 任务定义
type Task struct {
    ID           string                 `json:"id"`
    WorkflowID   string                 `json:"workflow_id"`
    ActivityID   string                 `json:"activity_id"`
    Type         TaskType               `json:"type"`
    Data         map[string]interface{} `json:"data"`
    Priority     int                    `json:"priority"`
    Timeout      time.Duration          `json:"timeout"`
    RetryCount   int                    `json:"retry_count"`
    MaxRetries   int                    `json:"max_retries"`
    CreatedAt    time.Time              `json:"created_at"`
    ScheduledAt  *time.Time             `json:"scheduled_at"`
    StartedAt    *time.Time             `json:"started_at"`
    CompletedAt  *time.Time             `json:"completed_at"`
    Status       TaskStatus             `json:"status"`
    Result       map[string]interface{} `json:"result"`
    Error        string                 `json:"error"`
}

// TaskType 任务类型
type TaskType int

const (
    TaskTypeActivity TaskType = iota
    TaskTypeGateway
    TaskTypeEvent
    TaskTypeSubprocess
)

// TaskStatus 任务状态
type TaskStatus int

const (
    TaskStatusPending TaskStatus = iota
    TaskStatusScheduled
    TaskStatusRunning
    TaskStatusCompleted
    TaskStatusFailed
    TaskStatusCancelled
)

// NewScheduler 创建调度器
func NewScheduler(config *EngineConfig) *Scheduler {
    ctx, cancel := context.WithCancel(context.Background())
    
    scheduler := &Scheduler{
        config:    config,
        taskQueue: make(chan *Task, config.MaxConcurrency*2),
        workers:   make([]*Worker, config.WorkerPoolSize),
        ctx:       ctx,
        cancel:    cancel,
    }
    
    // 创建工作线程
    for i := 0; i < config.WorkerPoolSize; i++ {
        scheduler.workers[i] = NewWorker(i, scheduler)
    }
    
    return scheduler
}

// Start 启动调度器
func (s *Scheduler) Start(ctx context.Context) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.running {
        return fmt.Errorf("scheduler is already running")
    }
    
    // 启动工作线程
    for _, worker := range s.workers {
        go worker.Start(s.ctx)
    }
    
    // 启动调度循环
    go s.scheduleLoop()
    
    s.running = true
    return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop(ctx context.Context) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if !s.running {
        return
    }
    
    s.cancel()
    s.running = false
}

// Schedule 调度任务
func (s *Scheduler) Schedule(task *Task) error {
    if !s.running {
        return fmt.Errorf("scheduler is not running")
    }
    
    // 设置调度时间
    now := time.Now()
    task.ScheduledAt = &now
    task.Status = TaskStatusScheduled
    
    // 保存任务状态
    if err := s.taskStore.SaveTask(task); err != nil {
        return fmt.Errorf("failed to save task: %w", err)
    }
    
    // 发送到任务队列
    select {
    case s.taskQueue <- task:
        s.metrics.RecordTaskScheduled(task)
        return nil
    case <-time.After(5 * time.Second):
        return fmt.Errorf("task queue is full")
    }
}

// scheduleLoop 调度循环
func (s *Scheduler) scheduleLoop() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-s.ctx.Done():
            return
        case task := <-s.taskQueue:
            s.assignTask(task)
        case <-ticker.C:
            s.checkTimeouts()
        }
    }
}

// assignTask 分配任务
func (s *Scheduler) assignTask(task *Task) {
    // 找到空闲的工作线程
    for _, worker := range s.workers {
        if worker.IsIdle() {
            worker.AssignTask(task)
            return
        }
    }
    
    // 如果没有空闲线程，重新入队
    go func() {
        time.Sleep(100 * time.Millisecond)
        s.Schedule(task)
    }()
}

// checkTimeouts 检查超时任务
func (s *Scheduler) checkTimeouts() {
    // 实现超时检查逻辑
}
```

### 3. 执行器设计

#### 3.1 执行器架构

```go
// Executor 执行器
type Executor struct {
    config       *EngineConfig
    activityRegistry ActivityRegistry
    taskStore    TaskStore
    eventBus     *EventBus
    metrics      *Metrics
    mutex        sync.RWMutex
    running      bool
}

// ActivityRegistry 活动注册表
type ActivityRegistry interface {
    Register(name string, handler ActivityHandler) error
    Get(name string) (ActivityHandler, bool)
    List() []string
}

// ActivityHandler 活动处理器
type ActivityHandler func(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)

// NewExecutor 创建执行器
func NewExecutor(config *EngineConfig) *Executor {
    return &Executor{
        config:           config,
        activityRegistry: NewActivityRegistry(),
    }
}

// Start 启动执行器
func (e *Executor) Start(ctx context.Context) error {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    
    if e.running {
        return fmt.Errorf("executor is already running")
    }
    
    e.running = true
    return nil
}

// Stop 停止执行器
func (e *Executor) Stop(ctx context.Context) {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    
    e.running = false
}

// ExecuteTask 执行任务
func (e *Executor) ExecuteTask(ctx context.Context, task *Task) error {
    if !e.running {
        return fmt.Errorf("executor is not running")
    }
    
    // 更新任务状态
    now := time.Now()
    task.StartedAt = &now
    task.Status = TaskStatusRunning
    
    if err := e.taskStore.SaveTask(task); err != nil {
        return fmt.Errorf("failed to save task status: %w", err)
    }
    
    // 创建执行上下文
    execCtx, cancel := context.WithTimeout(ctx, task.Timeout)
    defer cancel()
    
    // 执行活动
    var result map[string]interface{}
    var err error
    
    switch task.Type {
    case TaskTypeActivity:
        result, err = e.executeActivity(execCtx, task)
    case TaskTypeGateway:
        result, err = e.executeGateway(execCtx, task)
    case TaskTypeEvent:
        result, err = e.executeEvent(execCtx, task)
    case TaskTypeSubprocess:
        result, err = e.executeSubprocess(execCtx, task)
    default:
        err = fmt.Errorf("unknown task type: %d", task.Type)
    }
    
    // 更新任务结果
    now = time.Now()
    task.CompletedAt = &now
    
    if err != nil {
        task.Status = TaskStatusFailed
        task.Error = err.Error()
        
        // 重试逻辑
        if task.RetryCount < task.MaxRetries {
            task.RetryCount++
            task.Status = TaskStatusPending
            task.StartedAt = nil
            task.CompletedAt = nil
            task.Error = ""
            
            // 延迟重试
            go func() {
                time.Sleep(e.config.RetryDelay)
                e.eventBus.Publish("task.retry", task)
            }()
        }
    } else {
        task.Status = TaskStatusCompleted
        task.Result = result
    }
    
    // 保存任务状态
    if err := e.taskStore.SaveTask(task); err != nil {
        return fmt.Errorf("failed to save task result: %w", err)
    }
    
    // 发布事件
    eventType := "task.completed"
    if task.Status == TaskStatusFailed {
        eventType = "task.failed"
    }
    e.eventBus.Publish(eventType, task)
    
    // 记录指标
    e.metrics.RecordTaskExecution(task)
    
    return nil
}

// executeActivity 执行活动
func (e *Executor) executeActivity(ctx context.Context, task *Task) (map[string]interface{}, error) {
    handler, exists := e.activityRegistry.Get(task.ActivityID)
    if !exists {
        return nil, fmt.Errorf("activity handler not found: %s", task.ActivityID)
    }
    
    return handler(ctx, task.Data)
}

// executeGateway 执行网关
func (e *Executor) executeGateway(ctx context.Context, task *Task) (map[string]interface{}, error) {
    // 实现网关逻辑
    return map[string]interface{}{
        "decision": "route_a",
    }, nil
}

// executeEvent 执行事件
func (e *Executor) executeEvent(ctx context.Context, task *Task) (map[string]interface{}, error) {
    // 实现事件处理逻辑
    return map[string]interface{}{
        "event_processed": true,
    }, nil
}

// executeSubprocess 执行子流程
func (e *Executor) executeSubprocess(ctx context.Context, task *Task) (map[string]interface{}, error) {
    // 实现子流程执行逻辑
    return map[string]interface{}{
        "subprocess_completed": true,
    }, nil
}
```

### 4. 工作线程设计

#### 4.1 工作线程实现

```go
// Worker 工作线程
type Worker struct {
    id       int
    scheduler *Scheduler
    taskChan chan *Task
    running  bool
    currentTask *Task
    mutex    sync.RWMutex
}

// NewWorker 创建工作线程
func NewWorker(id int, scheduler *Scheduler) *Worker {
    return &Worker{
        id:        id,
        scheduler: scheduler,
        taskChan:  make(chan *Task, 1),
    }
}

// Start 启动工作线程
func (w *Worker) Start(ctx context.Context) {
    w.mutex.Lock()
    w.running = true
    w.mutex.Unlock()
    
    for {
        select {
        case <-ctx.Done():
            w.mutex.Lock()
            w.running = false
            w.mutex.Unlock()
            return
        case task := <-w.taskChan:
            w.executeTask(ctx, task)
        }
    }
}

// AssignTask 分配任务
func (w *Worker) AssignTask(task *Task) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    if !w.running {
        return
    }
    
    select {
    case w.taskChan <- task:
        w.currentTask = task
    default:
        // 任务通道已满，拒绝任务
    }
}

// IsIdle 检查是否空闲
func (w *Worker) IsIdle() bool {
    w.mutex.RLock()
    defer w.mutex.RUnlock()
    
    return w.running && w.currentTask == nil
}

// executeTask 执行任务
func (w *Worker) executeTask(ctx context.Context, task *Task) {
    defer func() {
        w.mutex.Lock()
        w.currentTask = nil
        w.mutex.Unlock()
    }()
    
    // 执行任务
    if err := w.scheduler.executor.ExecuteTask(ctx, task); err != nil {
        log.Printf("Worker %d failed to execute task %s: %v", w.id, task.ID, err)
    }
}
```

### 5. 存储层设计

#### 5.1 存储接口

```go
// TaskStore 任务存储接口
type TaskStore interface {
    SaveTask(task *Task) error
    GetTask(id string) (*Task, error)
    GetTasksByWorkflow(workflowID string) ([]*Task, error)
    GetPendingTasks() ([]*Task, error)
    DeleteTask(id string) error
}

// WorkflowStore 工作流存储接口
type WorkflowStore interface {
    SaveWorkflow(workflow *Workflow) error
    GetWorkflow(id string) (*Workflow, error)
    ListWorkflows() ([]*Workflow, error)
    DeleteWorkflow(id string) error
}

// Storage 存储实现
type Storage struct {
    config       *EngineConfig
    taskStore    TaskStore
    workflowStore WorkflowStore
    db           *gorm.DB
    redis        *redis.Client
}

// NewStorage 创建存储
func NewStorage(config *EngineConfig) *Storage {
    storage := &Storage{
        config: config,
    }
    
    // 根据配置选择存储实现
    switch config.PersistenceMode {
    case "database":
        storage.taskStore = NewDatabaseTaskStore(storage.db)
        storage.workflowStore = NewDatabaseWorkflowStore(storage.db)
    case "redis":
        storage.taskStore = NewRedisTaskStore(storage.redis)
        storage.workflowStore = NewRedisWorkflowStore(storage.redis)
    case "hybrid":
        storage.taskStore = NewHybridTaskStore(storage.db, storage.redis)
        storage.workflowStore = NewHybridWorkflowStore(storage.db, storage.redis)
    default:
        storage.taskStore = NewMemoryTaskStore()
        storage.workflowStore = NewMemoryWorkflowStore()
    }
    
    return storage
}
```

### 6. 事件总线设计

#### 6.1 事件总线实现

```go
// EventBus 事件总线
type EventBus struct {
    config       *EngineConfig
    handlers     map[string][]EventHandler
    publisher    EventPublisher
    subscriber   EventSubscriber
    mutex        sync.RWMutex
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event *Event) error

// Event 事件定义
type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Source    string                 `json:"source"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}

// EventPublisher 事件发布者接口
type EventPublisher interface {
    Publish(topic string, event *Event) error
}

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
    Subscribe(topic string, handler EventHandler) error
    Unsubscribe(topic string, handler EventHandler) error
}

// NewEventBus 创建事件总线
func NewEventBus(config *EngineConfig) *EventBus {
    eventBus := &EventBus{
        config:   config,
        handlers: make(map[string][]EventHandler),
    }
    
    // 根据配置选择事件总线实现
    switch config.EventBusType {
    case "redis":
        eventBus.publisher = NewRedisEventPublisher(config)
        eventBus.subscriber = NewRedisEventSubscriber(config)
    case "rabbitmq":
        eventBus.publisher = NewRabbitMQEventPublisher(config)
        eventBus.subscriber = NewRabbitMQEventSubscriber(config)
    case "kafka":
        eventBus.publisher = NewKafkaEventPublisher(config)
        eventBus.subscriber = NewKafkaEventSubscriber(config)
    default:
        eventBus.publisher = NewMemoryEventPublisher()
        eventBus.subscriber = NewMemoryEventSubscriber()
    }
    
    return eventBus
}

// Start 启动事件总线
func (e *EventBus) Start(ctx context.Context) error {
    // 启动事件订阅
    return e.subscriber.Start(ctx)
}

// Stop 停止事件总线
func (e *EventBus) Stop(ctx context.Context) {
    e.subscriber.Stop(ctx)
}

// Publish 发布事件
func (e *EventBus) Publish(topic string, data interface{}) error {
    event := &Event{
        ID:        uuid.New().String(),
        Type:      topic,
        Source:    "workflow-engine",
        Data:      data.(map[string]interface{}),
        Timestamp: time.Now(),
    }
    
    return e.publisher.Publish(topic, event)
}

// Subscribe 订阅事件
func (e *EventBus) Subscribe(topic string, handler EventHandler) error {
    e.mutex.Lock()
    defer e.mutex.Unlock()
    
    e.handlers[topic] = append(e.handlers[topic], handler)
    return e.subscriber.Subscribe(topic, handler)
}
```

## 总结

工作流引擎设计模块提供了完整的工作流执行框架，包括：

1. **引擎架构**: 模块化设计，支持高并发和高可用
2. **调度器**: 智能任务调度和负载均衡
3. **执行器**: 灵活的活动执行和错误处理
4. **工作线程**: 高效的任务执行池
5. **存储层**: 多种持久化策略
6. **事件总线**: 松耦合的事件驱动架构

这些组件协同工作，为工作流系统提供了强大的执行能力。

---

**下一步**: 继续完善工作流模式模块，提供常用的工作流设计模式。
