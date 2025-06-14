# 04-编排vs协同模式 (Orchestration vs Choreography Pattern)

## 目录

- [04-编排vs协同模式 (Orchestration vs Choreography Pattern)](#04-编排vs协同模式-orchestration-vs-choreography-pattern)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 定义](#11-定义)
    - [1.2 问题描述](#12-问题描述)
    - [1.3 设计目标](#13-设计目标)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 编排模式模型](#21-编排模式模型)
    - [2.2 协同模式模型](#22-协同模式模型)
    - [2.3 模式等价性](#23-模式等价性)
  - [3. 数学基础](#3-数学基础)
    - [3.1 状态机理论](#31-状态机理论)
    - [3.2 事件流理论](#32-事件流理论)
    - [3.3 复杂度分析](#33-复杂度分析)
  - [4. 编排模式](#4-编排模式)
    - [4.1 中央协调器](#41-中央协调器)
    - [4.2 编排执行](#42-编排执行)
    - [4.3 状态管理](#43-状态管理)
  - [5. 协同模式](#5-协同模式)
    - [5.1 事件驱动](#51-事件驱动)
    - [5.2 服务协调](#52-服务协调)
    - [5.3 消息传递](#53-消息传递)
  - [6. Go语言实现](#6-go语言实现)
    - [6.1 编排模式实现](#61-编排模式实现)
    - [6.2 协同模式实现](#62-协同模式实现)
    - [6.3 服务接口实现](#63-服务接口实现)
    - [6.4 事件实现](#64-事件实现)
    - [6.5 使用示例](#65-使用示例)
  - [7. 性能分析](#7-性能分析)
    - [7.1 时间复杂度](#71-时间复杂度)
    - [7.2 空间复杂度](#72-空间复杂度)
    - [7.3 扩展性分析](#73-扩展性分析)
  - [8. 应用场景](#8-应用场景)
    - [8.1 编排模式适用场景](#81-编排模式适用场景)
    - [8.2 协同模式适用场景](#82-协同模式适用场景)
    - [8.3 混合模式](#83-混合模式)
  - [9. 最佳实践](#9-最佳实践)
    - [9.1 模式选择](#91-模式选择)
    - [9.2 设计原则](#92-设计原则)
    - [9.3 错误处理](#93-错误处理)
    - [9.4 监控指标](#94-监控指标)
  - [10. 总结](#10-总结)
    - [10.1 关键要点](#101-关键要点)
    - [10.2 未来发展方向](#102-未来发展方向)

## 1. 概述

### 1.1 定义

编排vs协同模式是两种不同的分布式系统协调方式：

- **编排模式 (Orchestration)**: 中央协调器控制整个业务流程的执行
- **协同模式 (Choreography)**: 各个服务通过事件和消息自主协调

### 1.2 问题描述

在微服务架构中，服务间的协调面临以下挑战：

- **服务依赖**: 多个服务需要协调完成复杂业务
- **状态管理**: 分布式状态的一致性问题
- **故障处理**: 部分服务失败时的处理策略
- **扩展性**: 服务数量增加时的协调复杂度

### 1.3 设计目标

1. **服务解耦**: 减少服务间的直接依赖
2. **状态一致性**: 确保分布式状态的一致性
3. **故障隔离**: 单个服务故障不影响整体
4. **可扩展性**: 支持服务的动态扩展

## 2. 形式化定义

### 2.1 编排模式模型

**定义 2.1 (编排模式)**
编排模式是一个三元组 $O = (C, S, F)$，其中：

- $C$ 是中央协调器
- $S = \{s_1, s_2, ..., s_n\}$ 是服务集合
- $F: C \times S \rightarrow S$ 是协调函数

**定义 2.2 (编排执行)**
编排执行是一个序列：
$$\langle c_0, s_1, c_1, s_2, c_2, ..., s_n, c_n \rangle$$
其中 $c_i$ 是协调器状态，$s_i$ 是服务调用。

### 2.2 协同模式模型

**定义 2.3 (协同模式)**
协同模式是一个四元组 $H = (S, E, M, R)$，其中：

- $S = \{s_1, s_2, ..., s_n\}$ 是服务集合
- $E = \{e_1, e_2, ..., e_m\}$ 是事件集合
- $M: S \times E \rightarrow S$ 是消息传递函数
- $R: S \times E \rightarrow 2^E$ 是反应函数

**定义 2.4 (协同执行)**
协同执行是一个事件序列：
$$\langle e_1, e_2, ..., e_k \rangle$$
其中每个事件触发相应的服务反应。

### 2.3 模式等价性

**定理 2.1 (模式等价性)**
对于任何编排模式 $O$，存在等价的协同模式 $H$，反之亦然。

**证明**:

- **编排到协同**: 将协调器的决策转换为事件
- **协同到编排**: 将事件序列转换为协调器状态

## 3. 数学基础

### 3.1 状态机理论

**定义 3.1 (分布式状态机)**
分布式状态机是一个五元组 $M = (Q, \Sigma, \delta, q_0, F)$，其中：

- $Q$ 是状态集合
- $\Sigma$ 是输入字母表
- $\delta: Q \times \Sigma \rightarrow Q$ 是状态转换函数
- $q_0 \in Q$ 是初始状态
- $F \subseteq Q$ 是接受状态集合

**定理 3.1 (状态一致性)**
在异步网络中，分布式状态机的一致性需要满足：
$$\forall i, j: |s_i - s_j| \leq 1$$
其中 $s_i$ 是节点 $i$ 的状态。

### 3.2 事件流理论

**定义 3.2 (事件流)**
事件流是一个有序的事件序列：
$$E = \langle e_1, e_2, ..., e_n \rangle$$

**定义 3.3 (事件因果关系)**
事件 $e_i$ 因果先于事件 $e_j$，记作 $e_i \rightarrow e_j$，当且仅当：

1. $e_i$ 和 $e_j$ 在同一进程中，且 $i < j$
2. $e_i$ 发送消息，$e_j$ 接收该消息
3. 存在事件 $e_k$ 使得 $e_i \rightarrow e_k \rightarrow e_j$

**定理 3.2 (因果一致性)**
因果一致性要求：
$$\forall e_i, e_j: e_i \rightarrow e_j \Rightarrow \text{deliver}(e_i) < \text{deliver}(e_j)$$

### 3.3 复杂度分析

**定理 3.3 (编排复杂度)**
编排模式的时间复杂度为 $O(n)$，空间复杂度为 $O(n)$，其中 $n$ 是服务数量。

**定理 3.4 (协同复杂度)**
协同模式的时间复杂度为 $O(m)$，空间复杂度为 $O(m)$，其中 $m$ 是事件数量。

## 4. 编排模式

### 4.1 中央协调器

```go
// Orchestrator 编排器接口
type Orchestrator interface {
    ExecuteWorkflow(workflow *Workflow) error
    AddService(service Service) error
    RemoveService(serviceID string) error
    GetStatus(workflowID string) (*WorkflowStatus, error)
}

// Workflow 工作流定义
type Workflow struct {
    ID       string
    Steps    []*WorkflowStep
    Variables map[string]interface{}
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
    ID       string
    ServiceID string
    Action   string
    Input    map[string]interface{}
    Output   map[string]interface{}
    Next     []string
    Condition func(map[string]interface{}) bool
}
```

### 4.2 编排执行

**定义 4.1 (顺序执行)**
步骤按预定义顺序依次执行：
$$\forall i < j: \text{step}_i \prec \text{step}_j$$

**定义 4.2 (并行执行)**
满足条件的步骤可以并行执行：
$$\text{step}_i \parallel \text{step}_j \Leftrightarrow \text{not } (\text{step}_i \rightarrow \text{step}_j \lor \text{step}_j \rightarrow \text{step}_i)$$

**定义 4.3 (条件执行)**
步骤根据条件决定是否执行：
$$\text{execute}(\text{step}) \Leftrightarrow \text{condition}(\text{step}, \text{context})$$

### 4.3 状态管理

```go
// WorkflowState 工作流状态
type WorkflowState struct {
    WorkflowID string
    Status     string
    CurrentStep string
    Variables  map[string]interface{}
    History    []*ExecutionStep
    StartTime  time.Time
    EndTime    time.Time
}

// ExecutionStep 执行步骤
type ExecutionStep struct {
    StepID    string
    ServiceID string
    Status    string
    Input     map[string]interface{}
    Output    map[string]interface{}
    Error     error
    StartTime time.Time
    EndTime   time.Time
}
```

## 5. 协同模式

### 5.1 事件驱动

```go
// Event 事件接口
type Event interface {
    ID() string
    Type() string
    Source() string
    Payload() interface{}
    Timestamp() time.Time
    CorrelationID() string
}

// EventBus 事件总线
type EventBus interface {
    Publish(event Event) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string) error
}

// EventHandler 事件处理器
type EventHandler func(event Event) error
```

### 5.2 服务协调

**定义 5.1 (事件反应)**
服务对事件的反应是一个函数：
$$R: S \times E \rightarrow 2^E$$

**定义 5.2 (状态转换)**
事件导致的状态转换：
$$\delta(s, e) = s'$$

**定义 5.3 (因果链)**
事件的因果链是一个序列：
$$\langle e_1, e_2, ..., e_n \rangle \text{ where } e_i \rightarrow e_{i+1}$$

### 5.3 消息传递

```go
// Message 消息接口
type Message interface {
    ID() string
    Type() string
    Source() string
    Target() string
    Payload() interface{}
    Headers() map[string]string
    Timestamp() time.Time
}

// MessageBroker 消息代理
type MessageBroker interface {
    Send(message Message) error
    Receive(queue string) (Message, error)
    Subscribe(queue string, handler MessageHandler) error
    Unsubscribe(queue string) error
}

// MessageHandler 消息处理器
type MessageHandler func(message Message) error
```

## 6. Go语言实现

### 6.1 编排模式实现

```go
// OrchestratorImpl 编排器实现
type OrchestratorImpl struct {
    services   map[string]Service
    workflows  map[string]*WorkflowState
    mu         sync.RWMutex
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// NewOrchestrator 创建编排器
func NewOrchestrator() *OrchestratorImpl {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &OrchestratorImpl{
        services:  make(map[string]Service),
        workflows: make(map[string]*WorkflowState),
        ctx:       ctx,
        cancel:    cancel,
    }
}

// AddService 添加服务
func (o *OrchestratorImpl) AddService(service Service) error {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    o.services[service.ID()] = service
    log.Printf("Service %s added to orchestrator", service.ID())
    
    return nil
}

// RemoveService 移除服务
func (o *OrchestratorImpl) RemoveService(serviceID string) error {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    delete(o.services, serviceID)
    log.Printf("Service %s removed from orchestrator", serviceID)
    
    return nil
}

// ExecuteWorkflow 执行工作流
func (o *OrchestratorImpl) ExecuteWorkflow(workflow *Workflow) error {
    // 创建工作流状态
    state := &WorkflowState{
        WorkflowID: workflow.ID,
        Status:     "running",
        Variables:  make(map[string]interface{}),
        History:    make([]*ExecutionStep, 0),
        StartTime:  time.Now(),
    }
    
    // 复制初始变量
    for k, v := range workflow.Variables {
        state.Variables[k] = v
    }
    
    o.mu.Lock()
    o.workflows[workflow.ID] = state
    o.mu.Unlock()
    
    // 启动工作流执行
    o.wg.Add(1)
    go func() {
        defer o.wg.Done()
        o.executeWorkflowSteps(workflow, state)
    }()
    
    return nil
}

// executeWorkflowSteps 执行工作流步骤
func (o *OrchestratorImpl) executeWorkflowSteps(workflow *Workflow, state *WorkflowState) {
    defer func() {
        state.EndTime = time.Now()
        if state.Status == "running" {
            state.Status = "completed"
        }
    }()
    
    // 创建步骤映射
    stepMap := make(map[string]*WorkflowStep)
    for _, step := range workflow.Steps {
        stepMap[step.ID] = step
    }
    
    // 查找开始步骤
    var startSteps []*WorkflowStep
    for _, step := range workflow.Steps {
        if len(step.Next) == 0 || step.ID == "start" {
            startSteps = append(startSteps, step)
        }
    }
    
    // 执行开始步骤
    for _, step := range startSteps {
        o.executeStep(step, state)
    }
    
    // 继续执行后续步骤
    for state.Status == "running" {
        time.Sleep(100 * time.Millisecond)
    }
}

// executeStep 执行单个步骤
func (o *OrchestratorImpl) executeStep(step *WorkflowStep, state *WorkflowState) {
    // 创建执行步骤记录
    execStep := &ExecutionStep{
        StepID:    step.ID,
        ServiceID: step.ServiceID,
        Status:    "running",
        Input:     make(map[string]interface{}),
        Output:    make(map[string]interface{}),
        StartTime: time.Now(),
    }
    
    // 复制输入变量
    for k, v := range step.Input {
        if value, exists := state.Variables[k]; exists {
            execStep.Input[k] = value
        } else {
            execStep.Input[k] = v
        }
    }
    
    // 检查条件
    if step.Condition != nil && !step.Condition(state.Variables) {
        execStep.Status = "skipped"
        execStep.EndTime = time.Now()
        state.History = append(state.History, execStep)
        return
    }
    
    // 获取服务
    o.mu.RLock()
    service, exists := o.services[step.ServiceID]
    o.mu.RUnlock()
    
    if !exists {
        execStep.Status = "failed"
        execStep.Error = fmt.Errorf("service %s not found", step.ServiceID)
        execStep.EndTime = time.Now()
        state.History = append(state.History, execStep)
        return
    }
    
    // 执行服务调用
    err := service.Execute(step.Action, execStep.Input, execStep.Output)
    
    execStep.EndTime = time.Now()
    
    if err != nil {
        execStep.Status = "failed"
        execStep.Error = err
        log.Printf("Step %s failed: %v", step.ID, err)
    } else {
        execStep.Status = "completed"
        log.Printf("Step %s completed", step.ID)
        
        // 更新状态变量
        for k, v := range execStep.Output {
            state.Variables[k] = v
        }
    }
    
    state.History = append(state.History, execStep)
    
    // 执行后续步骤
    for _, nextStepID := range step.Next {
        if nextStep, exists := o.findStepByID(nextStepID); exists {
            go o.executeStep(nextStep, state)
        }
    }
}

// findStepByID 根据ID查找步骤
func (o *OrchestratorImpl) findStepByID(stepID string) (*WorkflowStep, bool) {
    // 这里简化实现，实际应该从工作流定义中查找
    return nil, false
}

// GetStatus 获取工作流状态
func (o *OrchestratorImpl) GetStatus(workflowID string) (*WorkflowStatus, error) {
    o.mu.RLock()
    state, exists := o.workflows[workflowID]
    o.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("workflow %s not found", workflowID)
    }
    
    return &WorkflowStatus{
        WorkflowID: state.WorkflowID,
        Status:     state.Status,
        CurrentStep: state.CurrentStep,
        Variables:  state.Variables,
        History:    state.History,
        StartTime:  state.StartTime,
        EndTime:    state.EndTime,
    }, nil
}

// WorkflowStatus 工作流状态
type WorkflowStatus struct {
    WorkflowID  string
    Status      string
    CurrentStep string
    Variables   map[string]interface{}
    History     []*ExecutionStep
    StartTime   time.Time
    EndTime     time.Time
}
```

### 6.2 协同模式实现

```go
// ChoreographyImpl 协同模式实现
type ChoreographyImpl struct {
    services   map[string]Service
    eventBus   EventBus
    mu         sync.RWMutex
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// NewChoreography 创建协同模式
func NewChoreography(eventBus EventBus) *ChoreographyImpl {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ChoreographyImpl{
        services: make(map[string]Service),
        eventBus: eventBus,
        ctx:      ctx,
        cancel:   cancel,
    }
}

// AddService 添加服务
func (c *ChoreographyImpl) AddService(service Service) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.services[service.ID()] = service
    
    // 注册服务的事件处理器
    for eventType, handler := range service.GetEventHandlers() {
        c.eventBus.Subscribe(eventType, handler)
    }
    
    log.Printf("Service %s added to choreography", service.ID())
    
    return nil
}

// RemoveService 移除服务
func (c *ChoreographyImpl) RemoveService(serviceID string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    service, exists := c.services[serviceID]
    if !exists {
        return fmt.Errorf("service %s not found", serviceID)
    }
    
    // 取消注册事件处理器
    for eventType := range service.GetEventHandlers() {
        c.eventBus.Unsubscribe(eventType)
    }
    
    delete(c.services, serviceID)
    log.Printf("Service %s removed from choreography", serviceID)
    
    return nil
}

// StartWorkflow 启动工作流
func (c *ChoreographyImpl) StartWorkflow(workflowID string, initialEvent Event) error {
    log.Printf("Starting workflow %s with event %s", workflowID, initialEvent.ID())
    
    // 发布初始事件
    return c.eventBus.Publish(initialEvent)
}

// EventBusImpl 事件总线实现
type EventBusImpl struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBusImpl {
    return &EventBusImpl{
        handlers: make(map[string][]EventHandler),
    }
}

// Publish 发布事件
func (eb *EventBusImpl) Publish(event Event) error {
    eb.mu.RLock()
    handlers := make([]EventHandler, len(eb.handlers[event.Type()]))
    copy(handlers, eb.handlers[event.Type()])
    eb.mu.RUnlock()
    
    // 异步处理事件
    go func() {
        for _, handler := range handlers {
            err := handler(event)
            if err != nil {
                log.Printf("Error handling event %s: %v", event.ID(), err)
            }
        }
    }()
    
    return nil
}

// Subscribe 订阅事件
func (eb *EventBusImpl) Subscribe(eventType string, handler EventHandler) error {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
    log.Printf("Handler subscribed to event type: %s", eventType)
    
    return nil
}

// Unsubscribe 取消订阅
func (eb *EventBusImpl) Unsubscribe(eventType string) error {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    delete(eb.handlers, eventType)
    log.Printf("Handlers unsubscribed from event type: %s", eventType)
    
    return nil
}
```

### 6.3 服务接口实现

```go
// Service 服务接口
type Service interface {
    ID() string
    Execute(action string, input map[string]interface{}, output map[string]interface{}) error
    GetEventHandlers() map[string]EventHandler
}

// BaseService 基础服务实现
type BaseService struct {
    id           string
    eventHandlers map[string]EventHandler
}

// NewBaseService 创建基础服务
func NewBaseService(id string) *BaseService {
    return &BaseService{
        id:           id,
        eventHandlers: make(map[string]EventHandler),
    }
}

// ID 获取服务ID
func (bs *BaseService) ID() string {
    return bs.id
}

// GetEventHandlers 获取事件处理器
func (bs *BaseService) GetEventHandlers() map[string]EventHandler {
    return bs.eventHandlers
}

// AddEventHandler 添加事件处理器
func (bs *BaseService) AddEventHandler(eventType string, handler EventHandler) {
    bs.eventHandlers[eventType] = handler
}

// OrderService 订单服务
type OrderService struct {
    *BaseService
}

// NewOrderService 创建订单服务
func NewOrderService() *OrderService {
    service := &OrderService{
        BaseService: NewBaseService("order-service"),
    }
    
    // 注册事件处理器
    service.AddEventHandler("order.created", service.handleOrderCreated)
    service.AddEventHandler("payment.completed", service.handlePaymentCompleted)
    service.AddEventHandler("inventory.reserved", service.handleInventoryReserved)
    
    return service
}

// Execute 执行服务操作
func (os *OrderService) Execute(action string, input map[string]interface{}, output map[string]interface{}) error {
    switch action {
    case "create_order":
        return os.createOrder(input, output)
    case "update_order":
        return os.updateOrder(input, output)
    case "cancel_order":
        return os.cancelOrder(input, output)
    default:
        return fmt.Errorf("unknown action: %s", action)
    }
}

// createOrder 创建订单
func (os *OrderService) createOrder(input map[string]interface{}, output map[string]interface{}) error {
    log.Printf("Order service creating order")
    
    // 模拟订单创建
    time.Sleep(100 * time.Millisecond)
    
    orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
    output["order_id"] = orderID
    output["status"] = "created"
    
    log.Printf("Order created: %s", orderID)
    
    return nil
}

// updateOrder 更新订单
func (os *OrderService) updateOrder(input map[string]interface{}, output map[string]interface{}) error {
    orderID := input["order_id"].(string)
    status := input["status"].(string)
    
    log.Printf("Order service updating order %s to status %s", orderID, status)
    
    // 模拟订单更新
    time.Sleep(50 * time.Millisecond)
    
    output["order_id"] = orderID
    output["status"] = status
    
    return nil
}

// cancelOrder 取消订单
func (os *OrderService) cancelOrder(input map[string]interface{}, output map[string]interface{}) error {
    orderID := input["order_id"].(string)
    
    log.Printf("Order service cancelling order %s", orderID)
    
    // 模拟订单取消
    time.Sleep(50 * time.Millisecond)
    
    output["order_id"] = orderID
    output["status"] = "cancelled"
    
    return nil
}

// handleOrderCreated 处理订单创建事件
func (os *OrderService) handleOrderCreated(event Event) error {
    log.Printf("Order service handling order created event: %s", event.ID())
    
    // 处理订单创建事件
    payload := event.Payload().(map[string]interface{})
    orderID := payload["order_id"].(string)
    
    // 发布支付事件
    paymentEvent := &BaseEvent{
        id:            fmt.Sprintf("payment-%d", time.Now().UnixNano()),
        eventType:     "payment.requested",
        source:        os.ID(),
        payload:       map[string]interface{}{"order_id": orderID, "amount": payload["amount"]},
        timestamp:     time.Now(),
        correlationID: event.CorrelationID(),
    }
    
    // 这里应该通过事件总线发布事件
    log.Printf("Payment event published: %s", paymentEvent.ID())
    
    return nil
}

// handlePaymentCompleted 处理支付完成事件
func (os *OrderService) handlePaymentCompleted(event Event) error {
    log.Printf("Order service handling payment completed event: %s", event.ID())
    
    // 处理支付完成事件
    payload := event.Payload().(map[string]interface{})
    orderID := payload["order_id"].(string)
    
    // 更新订单状态
    output := make(map[string]interface{})
    os.updateOrder(map[string]interface{}{"order_id": orderID, "status": "paid"}, output)
    
    // 发布库存预留事件
    inventoryEvent := &BaseEvent{
        id:            fmt.Sprintf("inventory-%d", time.Now().UnixNano()),
        eventType:     "inventory.reserve",
        source:        os.ID(),
        payload:       map[string]interface{}{"order_id": orderID, "items": payload["items"]},
        timestamp:     time.Now(),
        correlationID: event.CorrelationID(),
    }
    
    log.Printf("Inventory event published: %s", inventoryEvent.ID())
    
    return nil
}

// handleInventoryReserved 处理库存预留事件
func (os *OrderService) handleInventoryReserved(event Event) error {
    log.Printf("Order service handling inventory reserved event: %s", event.ID())
    
    // 处理库存预留事件
    payload := event.Payload().(map[string]interface{})
    orderID := payload["order_id"].(string)
    
    // 更新订单状态
    output := make(map[string]interface{})
    os.updateOrder(map[string]interface{}{"order_id": orderID, "status": "confirmed"}, output)
    
    return nil
}

// PaymentService 支付服务
type PaymentService struct {
    *BaseService
}

// NewPaymentService 创建支付服务
func NewPaymentService() *PaymentService {
    service := &PaymentService{
        BaseService: NewBaseService("payment-service"),
    }
    
    // 注册事件处理器
    service.AddEventHandler("payment.requested", service.handlePaymentRequested)
    
    return service
}

// Execute 执行服务操作
func (ps *PaymentService) Execute(action string, input map[string]interface{}, output map[string]interface{}) error {
    switch action {
    case "process_payment":
        return ps.processPayment(input, output)
    default:
        return fmt.Errorf("unknown action: %s", action)
    }
}

// processPayment 处理支付
func (ps *PaymentService) processPayment(input map[string]interface{}, output map[string]interface{}) error {
    orderID := input["order_id"].(string)
    amount := input["amount"].(float64)
    
    log.Printf("Payment service processing payment for order %s, amount: %.2f", orderID, amount)
    
    // 模拟支付处理
    time.Sleep(200 * time.Millisecond)
    
    output["order_id"] = orderID
    output["amount"] = amount
    output["status"] = "completed"
    output["transaction_id"] = fmt.Sprintf("txn-%d", time.Now().UnixNano())
    
    log.Printf("Payment completed for order %s", orderID)
    
    return nil
}

// handlePaymentRequested 处理支付请求事件
func (ps *PaymentService) handlePaymentRequested(event Event) error {
    log.Printf("Payment service handling payment requested event: %s", event.ID())
    
    // 处理支付请求事件
    payload := event.Payload().(map[string]interface{})
    
    // 处理支付
    output := make(map[string]interface{})
    ps.processPayment(payload, output)
    
    // 发布支付完成事件
    paymentCompletedEvent := &BaseEvent{
        id:            fmt.Sprintf("payment-completed-%d", time.Now().UnixNano()),
        eventType:     "payment.completed",
        source:        ps.ID(),
        payload:       output,
        timestamp:     time.Now(),
        correlationID: event.CorrelationID(),
    }
    
    log.Printf("Payment completed event published: %s", paymentCompletedEvent.ID())
    
    return nil
}

// InventoryService 库存服务
type InventoryService struct {
    *BaseService
}

// NewInventoryService 创建库存服务
func NewInventoryService() *InventoryService {
    service := &InventoryService{
        BaseService: NewBaseService("inventory-service"),
    }
    
    // 注册事件处理器
    service.AddEventHandler("inventory.reserve", service.handleInventoryReserve)
    
    return service
}

// Execute 执行服务操作
func (is *InventoryService) Execute(action string, input map[string]interface{}, output map[string]interface{}) error {
    switch action {
    case "reserve_inventory":
        return is.reserveInventory(input, output)
    default:
        return fmt.Errorf("unknown action: %s", action)
    }
}

// reserveInventory 预留库存
func (is *InventoryService) reserveInventory(input map[string]interface{}, output map[string]interface{}) error {
    orderID := input["order_id"].(string)
    items := input["items"].([]interface{})
    
    log.Printf("Inventory service reserving inventory for order %s", orderID)
    
    // 模拟库存预留
    time.Sleep(150 * time.Millisecond)
    
    output["order_id"] = orderID
    output["items"] = items
    output["status"] = "reserved"
    
    log.Printf("Inventory reserved for order %s", orderID)
    
    return nil
}

// handleInventoryReserve 处理库存预留事件
func (is *InventoryService) handleInventoryReserve(event Event) error {
    log.Printf("Inventory service handling inventory reserve event: %s", event.ID())
    
    // 处理库存预留事件
    payload := event.Payload().(map[string]interface{})
    
    // 预留库存
    output := make(map[string]interface{})
    is.reserveInventory(payload, output)
    
    // 发布库存预留完成事件
    inventoryReservedEvent := &BaseEvent{
        id:            fmt.Sprintf("inventory-reserved-%d", time.Now().UnixNano()),
        eventType:     "inventory.reserved",
        source:        is.ID(),
        payload:       output,
        timestamp:     time.Now(),
        correlationID: event.CorrelationID(),
    }
    
    log.Printf("Inventory reserved event published: %s", inventoryReservedEvent.ID())
    
    return nil
}
```

### 6.4 事件实现

```go
// BaseEvent 基础事件实现
type BaseEvent struct {
    id            string
    eventType     string
    source        string
    payload       interface{}
    timestamp     time.Time
    correlationID string
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(id, eventType, source string, payload interface{}, correlationID string) *BaseEvent {
    return &BaseEvent{
        id:            id,
        eventType:     eventType,
        source:        source,
        payload:       payload,
        timestamp:     time.Now(),
        correlationID: correlationID,
    }
}

// ID 获取事件ID
func (be *BaseEvent) ID() string {
    return be.id
}

// Type 获取事件类型
func (be *BaseEvent) Type() string {
    return be.eventType
}

// Source 获取事件源
func (be *BaseEvent) Source() string {
    return be.source
}

// Payload 获取事件载荷
func (be *BaseEvent) Payload() interface{} {
    return be.payload
}

// Timestamp 获取时间戳
func (be *BaseEvent) Timestamp() time.Time {
    return be.timestamp
}

// CorrelationID 获取关联ID
func (be *BaseEvent) CorrelationID() string {
    return be.correlationID
}
```

### 6.5 使用示例

```go
// main.go
func main() {
    // 创建事件总线
    eventBus := NewEventBus()
    
    // 创建编排器
    orchestrator := NewOrchestrator()
    
    // 创建协同模式
    choreography := NewChoreography(eventBus)
    
    // 创建服务
    orderService := NewOrderService()
    paymentService := NewPaymentService()
    inventoryService := NewInventoryService()
    
    // 添加到编排器
    orchestrator.AddService(orderService)
    orchestrator.AddService(paymentService)
    orchestrator.AddService(inventoryService)
    
    // 添加到协同模式
    choreography.AddService(orderService)
    choreography.AddService(paymentService)
    choreography.AddService(inventoryService)
    
    // 编排模式示例
    log.Println("=== Orchestration Pattern ===")
    
    workflow := &Workflow{
        ID: "order-processing",
        Steps: []*WorkflowStep{
            {
                ID:       "create_order",
                ServiceID: "order-service",
                Action:   "create_order",
                Input:    map[string]interface{}{"customer": "Alice", "amount": 100.0},
                Next:     []string{"process_payment"},
            },
            {
                ID:       "process_payment",
                ServiceID: "payment-service",
                Action:   "process_payment",
                Input:    map[string]interface{}{"order_id": "{{order_id}}", "amount": "{{amount}}"},
                Next:     []string{"reserve_inventory"},
            },
            {
                ID:       "reserve_inventory",
                ServiceID: "inventory-service",
                Action:   "reserve_inventory",
                Input:    map[string]interface{}{"order_id": "{{order_id}}", "items": []string{"item1", "item2"}},
                Next:     []string{},
            },
        },
    }
    
    err := orchestrator.ExecuteWorkflow(workflow)
    if err != nil {
        log.Printf("Orchestration failed: %v", err)
    }
    
    // 等待编排执行完成
    time.Sleep(2 * time.Second)
    
    // 协同模式示例
    log.Println("=== Choreography Pattern ===")
    
    // 创建初始事件
    initialEvent := NewBaseEvent(
        fmt.Sprintf("order-%d", time.Now().UnixNano()),
        "order.created",
        "client",
        map[string]interface{}{
            "customer": "Bob",
            "amount":   150.0,
            "items":    []string{"item3", "item4"},
        },
        fmt.Sprintf("correlation-%d", time.Now().UnixNano()),
    )
    
    err = choreography.StartWorkflow("order-processing-choreography", initialEvent)
    if err != nil {
        log.Printf("Choreography failed: %v", err)
    }
    
    // 等待协同执行完成
    time.Sleep(3 * time.Second)
    
    log.Println("=== Comparison ===")
    log.Println("Orchestration: Centralized control, easier to monitor and debug")
    log.Println("Choreography: Decentralized, better scalability and fault tolerance")
}
```

## 7. 性能分析

### 7.1 时间复杂度

**定理 7.1 (编排模式时间复杂度)**
编排模式的时间复杂度为 $O(n)$，其中 $n$ 是服务调用次数。

**定理 7.2 (协同模式时间复杂度)**
协同模式的时间复杂度为 $O(m)$，其中 $m$ 是事件数量。

### 7.2 空间复杂度

**定理 7.3 (编排模式空间复杂度)**
编排模式的空间复杂度为 $O(n)$，主要用于存储工作流状态。

**定理 7.4 (协同模式空间复杂度)**
协同模式的空间复杂度为 $O(m)$，主要用于存储事件和状态。

### 7.3 扩展性分析

**定理 7.5 (编排模式扩展性)**
编排模式的扩展性受中央协调器限制：
$$\text{Scalability} = O(1)$$

**定理 7.6 (协同模式扩展性)**
协同模式的扩展性随服务数量线性增长：
$$\text{Scalability} = O(n)$$

## 8. 应用场景

### 8.1 编排模式适用场景

- **简单流程**: 步骤较少、逻辑清晰的流程
- **集中控制**: 需要统一监控和管理的场景
- **调试友好**: 需要详细执行日志的场景
- **事务管理**: 需要强一致性的场景

### 8.2 协同模式适用场景

- **复杂流程**: 步骤较多、逻辑复杂的流程
- **高扩展性**: 需要动态扩展服务的场景
- **故障隔离**: 需要故障隔离的场景
- **松耦合**: 服务间需要松耦合的场景

### 8.3 混合模式

- **分层架构**: 不同层次使用不同模式
- **渐进迁移**: 从编排逐步迁移到协同
- **场景适配**: 根据具体场景选择合适模式

## 9. 最佳实践

### 9.1 模式选择

```go
// 模式选择指南
type PatternSelection struct {
    // 编排模式适合：
    // 1. 流程简单，步骤较少
    // 2. 需要集中监控
    // 3. 调试要求高
    // 4. 强一致性要求
    
    // 协同模式适合：
    // 1. 流程复杂，步骤较多
    // 2. 需要高扩展性
    // 3. 故障隔离要求
    // 4. 松耦合要求
}
```

### 9.2 设计原则

```go
// 设计原则
type DesignPrinciples struct {
    // 1. 单一职责原则
    // 2. 开闭原则
    // 3. 依赖倒置原则
    // 4. 接口隔离原则
}
```

### 9.3 错误处理

```go
// 错误处理策略
type ErrorHandling struct {
    // 1. 重试机制
    // 2. 补偿机制
    // 3. 超时处理
    // 4. 降级处理
}
```

### 9.4 监控指标

```go
// 监控指标
type MonitoringMetrics struct {
    ExecutionTime    time.Duration
    SuccessRate      float64
    ErrorRate        float64
    Throughput       float64
    ResourceUsage    float64
}
```

## 10. 总结

编排模式和协同模式是两种不同的分布式系统协调方式，各有优缺点，需要根据具体场景选择合适的模式。

### 10.1 关键要点

1. **编排模式**: 中央控制，易于监控，适合简单流程
2. **协同模式**: 分布式协调，高扩展性，适合复杂流程
3. **模式选择**: 根据业务需求和技术约束选择
4. **混合使用**: 可以在不同层次使用不同模式

### 10.2 未来发展方向

1. **智能协调**: 使用ML优化协调策略
2. **自适应模式**: 根据负载自动切换模式
3. **可视化工具**: 图形化流程设计和监控
4. **标准化接口**: 统一的协调接口标准

---

**参考文献**:

1. Hohpe, G., & Woolf, B. (2003). "Enterprise Integration Patterns"
2. Newman, S. (2021). "Building Microservices"
3. Richardson, C. (2018). "Microservices Patterns"

**相关链接**:

- [01-状态机模式](../01-State-Machine-Pattern.md)
- [02-工作流引擎模式](../02-Workflow-Engine-Pattern.md)
- [03-任务队列模式](../03-Task-Queue-Pattern.md)
