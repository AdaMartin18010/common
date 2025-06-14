# 04-编排vs协同模式

 (Orchestration vs Choreography Pattern)

## 目录

- [04-编排vs协同模式](#04-编排vs协同模式)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心组件](#12-核心组件)
      - [编排模式组件](#编排模式组件)
      - [协同模式组件](#协同模式组件)
    - [1.3 模式结构](#13-模式结构)
      - [编排模式结构](#编排模式结构)
      - [协同模式结构](#协同模式结构)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 编排模式数学模型](#21-编排模式数学模型)
    - [2.2 协同模式数学模型](#22-协同模式数学模型)
    - [2.3 执行函数定义](#23-执行函数定义)
      - [编排模式执行函数](#编排模式执行函数)
      - [协同模式执行函数](#协同模式执行函数)
  - [3. 数学证明](#3-数学证明)
    - [3.1 编排模式确定性定理](#31-编排模式确定性定理)
    - [3.2 协同模式非确定性定理](#32-协同模式非确定性定理)
    - [3.3 复杂度比较定理](#33-复杂度比较定理)
  - [4. 设计原则](#4-设计原则)
    - [4.1 编排模式原则](#41-编排模式原则)
    - [4.2 协同模式原则](#42-协同模式原则)
    - [4.3 选择原则](#43-选择原则)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 编排模式实现](#51-编排模式实现)
    - [5.2 协同模式实现](#52-协同模式实现)
    - [5.3 混合模式实现](#53-混合模式实现)
  - [6. 应用场景](#6-应用场景)
    - [6.1 订单处理系统](#61-订单处理系统)
      - [编排模式实现](#编排模式实现)
      - [协同模式实现](#协同模式实现)
  - [7. 性能分析](#7-性能分析)
    - [7.1 时间复杂度](#71-时间复杂度)
      - [编排模式](#编排模式)
      - [协同模式](#协同模式)
    - [7.2 空间复杂度](#72-空间复杂度)
      - [编排模式7](#编排模式7)
      - [协同模式7](#协同模式7)
    - [7.3 并发性能](#73-并发性能)
      - [编排模式73](#编排模式73)
      - [协同模式73](#协同模式73)
  - [8. 最佳实践](#8-最佳实践)
    - [8.1 编排模式最佳实践](#81-编排模式最佳实践)
    - [8.2 协同模式最佳实践](#82-协同模式最佳实践)
    - [8.3 选择指南](#83-选择指南)
  - [9. 相关模式](#9-相关模式)
    - [9.1 状态机模式](#91-状态机模式)
    - [9.2 观察者模式](#92-观察者模式)
    - [9.3 命令模式](#93-命令模式)

## 1. 概念与定义

### 1.1 基本概念

编排vs协同模式是两种不同的分布式系统协调方式。编排模式通过中央控制器来协调各个服务的执行，而协同模式通过事件驱动的方式让各个服务自主协调。

**编排模式定义**: 通过中央协调器（Orchestrator）来控制和管理整个业务流程的执行，协调器负责任务的调度、状态管理和错误处理。

**协同模式定义**: 通过事件驱动的方式，让各个服务自主协调，每个服务根据接收到的事件自主决定下一步行动。

### 1.2 核心组件

#### 编排模式组件

- **Orchestrator (协调器)**: 中央控制器，负责任务调度和状态管理
- **Service (服务)**: 具体的业务服务
- **Workflow (工作流)**: 业务流程定义
- **StateManager (状态管理器)**: 管理执行状态

#### 协同模式组件

- **EventBus (事件总线)**: 事件发布和订阅的中间件
- **Service (服务)**: 自主的业务服务
- **Event (事件)**: 服务间通信的消息
- **EventHandler (事件处理器)**: 处理事件的逻辑

### 1.3 模式结构

#### 编排模式结构

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Orchestrator  │    │    Service A    │    │    Service B    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + coordinate()  │◄──►│ + execute()     │◄──►│ + execute()     │
│ + schedule()    │    │ + getStatus()   │    │ + getStatus()   │
│ + monitor()     │    │ + rollback()    │    │ + rollback()    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Workflow      │    │ StateManager    │    │    Service C    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + definition    │    │ + state         │    │ + execute()     │
│ + steps         │    │ + transitions   │    │ + getStatus()   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

#### 协同模式结构

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Service A    │    │    EventBus     │    │    Service B    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + publish()     │◄──►│ + subscribe()   │◄──►│ + publish()     │
│ + subscribe()   │    │ + route()       │    │ + subscribe()   │
│ + handle()      │    │ + broadcast()   │    │ + handle()      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Event A       │    │   Event B       │    │    Service C    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + type          │    │ + type          │    │ + publish()     │
│ + data          │    │ + data          │    │ + subscribe()   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 2. 形式化定义

### 2.1 编排模式数学模型

设 $O = (S, W, T, I, F)$ 为一个编排系统，其中：

- $S = \{s_1, s_2, ..., s_n\}$ 是服务集合
- $W = \{w_1, w_2, ..., w_m\}$ 是工作流步骤集合
- $T: W \times S \rightarrow S$ 是转换函数
- $I \in W$ 是初始步骤
- $F \subseteq W$ 是终止步骤集合

### 2.2 协同模式数学模型

设 $C = (S, E, P, H)$ 为一个协同系统，其中：

- $S = \{s_1, s_2, ..., s_n\}$ 是服务集合
- $E = \{e_1, e_2, ..., e_m\}$ 是事件集合
- $P: S \times E \rightarrow S$ 是发布函数
- $H: S \times E \rightarrow S$ 是处理函数

### 2.3 执行函数定义

#### 编排模式执行函数

$$orchestrate(w, context) = (result, next\_step)$$

其中：

- $w \in W$ 是当前工作流步骤
- $context$ 是执行上下文
- $result$ 是执行结果
- $next\_step \in W$ 是下一个步骤

#### 协同模式执行函数

$$choreograph(e, service) = (result, next\_events)$$

其中：

- $e \in E$ 是当前事件
- $service \in S$ 是处理服务的服务
- $result$ 是处理结果
- $next\_events \subseteq E$ 是下一个事件集合

## 3. 数学证明

### 3.1 编排模式确定性定理

**定理**: 在编排模式下，给定相同的初始状态和输入，执行路径是确定的。

**证明**:

1. 设编排系统 $O = (S, W, T, I, F)$
2. 对于任意步骤 $w \in W$，转换函数 $T$ 是确定性的
3. 因此，给定当前步骤和状态，下一个步骤是唯一确定的
4. 结论：编排模式的执行路径是确定的

### 3.2 协同模式非确定性定理

**定理**: 在协同模式下，给定相同的初始状态和输入，执行路径可能不是唯一的。

**证明**:

1. 设协同系统 $C = (S, E, P, H)$
2. 事件的处理可能依赖于服务的内部状态
3. 不同服务可能以不同顺序处理相同的事件
4. 因此，执行路径可能不是唯一的

### 3.3 复杂度比较定理

**定理**: 编排模式的复杂度为 $O(n)$，协同模式的复杂度为 $O(n^2)$，其中 $n$ 是服务数量。

**证明**:

1. **编排模式**: 中央协调器需要与每个服务通信，复杂度为 $O(n)$
2. **协同模式**: 每个服务可能需要与其他所有服务通信，复杂度为 $O(n^2)$

## 4. 设计原则

### 4.1 编排模式原则

1. **中央控制**: 通过中央协调器统一管理业务流程
2. **状态集中**: 将状态管理集中在协调器中
3. **错误处理**: 在协调器中统一处理错误和异常

### 4.2 协同模式原则

1. **服务自治**: 每个服务自主决定如何处理事件
2. **松耦合**: 服务之间通过事件进行松耦合通信
3. **事件驱动**: 基于事件驱动架构设计系统

### 4.3 选择原则

1. **简单流程**: 对于简单的线性流程，选择编排模式
2. **复杂流程**: 对于复杂的非线性流程，选择协同模式
3. **团队结构**: 根据团队的组织结构选择合适的模式

## 5. Go语言实现

### 5.1 编排模式实现

```go
package orchestration

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Service 服务接口
type Service interface {
 Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
 GetName() string
 GetStatus() string
 Rollback(ctx context.Context, data map[string]interface{}) error
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
 Name        string
 Service     Service
 Condition   func(data map[string]interface{}) bool
 NextSteps   []string
 Timeout     time.Duration
 RetryCount  int
}

// Workflow 工作流定义
type Workflow struct {
 Name  string
 Steps map[string]*WorkflowStep
 Start string
 End   []string
}

// Orchestrator 协调器
type Orchestrator struct {
 workflow    *Workflow
 state       map[string]interface{}
 history     []ExecutionRecord
 mu          sync.RWMutex
 ctx         context.Context
 cancel      context.CancelFunc
}

// ExecutionRecord 执行记录
type ExecutionRecord struct {
 StepName    string
 StartTime   time.Time
 EndTime     time.Time
 Result      map[string]interface{}
 Error       error
 Status      string
}

// NewOrchestrator 创建协调器
func NewOrchestrator(workflow *Workflow) *Orchestrator {
 ctx, cancel := context.WithCancel(context.Background())
 return &Orchestrator{
  workflow: workflow,
  state:    make(map[string]interface{}),
  history:  make([]ExecutionRecord, 0),
  ctx:      ctx,
  cancel:   cancel,
 }
}

// Execute 执行工作流
func (o *Orchestrator) Execute(ctx context.Context, initialData map[string]interface{}) error {
 o.mu.Lock()
 o.state = initialData
 o.mu.Unlock()
 
 currentStep := o.workflow.Start
 
 for {
  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }
  
  // 检查是否到达终止步骤
  if o.isEndStep(currentStep) {
   break
  }
  
  // 执行当前步骤
  step, exists := o.workflow.Steps[currentStep]
  if !exists {
   return fmt.Errorf("步骤 %s 不存在", currentStep)
  }
  
  record := ExecutionRecord{
   StepName:  step.Name,
   StartTime: time.Now(),
  }
  
  // 检查执行条件
  if step.Condition != nil && !step.Condition(o.state) {
   record.Status = "skipped"
   record.EndTime = time.Now()
   o.addHistory(record)
   
   // 选择下一个步骤
   nextStep := o.selectNextStep(step)
   if nextStep == "" {
    break
   }
   currentStep = nextStep
   continue
  }
  
  // 执行服务
  stepCtx, cancel := context.WithTimeout(ctx, step.Timeout)
  result, err := step.Service.Execute(stepCtx, o.state)
  cancel()
  
  record.EndTime = time.Now()
  record.Result = result
  record.Error = err
  
  if err != nil {
   record.Status = "failed"
   o.addHistory(record)
   
   // 重试逻辑
   if step.RetryCount > 0 {
    step.RetryCount--
    continue
   }
   
   // 回滚逻辑
   o.rollback(currentStep)
   return fmt.Errorf("步骤 %s 执行失败: %v", currentStep, err)
  }
  
  record.Status = "completed"
  o.addHistory(record)
  
  // 更新状态
  o.mu.Lock()
  for k, v := range result {
   o.state[k] = v
  }
  o.mu.Unlock()
  
  // 选择下一个步骤
  nextStep := o.selectNextStep(step)
  if nextStep == "" {
   break
  }
  currentStep = nextStep
 }
 
 return nil
}

// isEndStep 检查是否为终止步骤
func (o *Orchestrator) isEndStep(stepName string) bool {
 for _, end := range o.workflow.End {
  if end == stepName {
   return true
  }
 }
 return false
}

// selectNextStep 选择下一个步骤
func (o *Orchestrator) selectNextStep(step *WorkflowStep) string {
 if len(step.NextSteps) == 0 {
  return ""
 }
 
 // 简单实现：选择第一个下一步骤
 // 实际应用中可能需要更复杂的逻辑
 return step.NextSteps[0]
}

// rollback 回滚
func (o *Orchestrator) rollback(currentStep string) {
 // 从历史记录中找到需要回滚的步骤
 for i := len(o.history) - 1; i >= 0; i-- {
  record := o.history[i]
  if record.Status == "completed" {
   step, exists := o.workflow.Steps[record.StepName]
   if exists {
    step.Service.Rollback(o.ctx, record.Result)
   }
  }
 }
}

// addHistory 添加历史记录
func (o *Orchestrator) addHistory(record ExecutionRecord) {
 o.mu.Lock()
 defer o.mu.Unlock()
 o.history = append(o.history, record)
}

// GetState 获取当前状态
func (o *Orchestrator) GetState() map[string]interface{} {
 o.mu.RLock()
 defer o.mu.RUnlock()
 
 result := make(map[string]interface{})
 for k, v := range o.state {
  result[k] = v
 }
 return result
}

// GetHistory 获取执行历史
func (o *Orchestrator) GetHistory() []ExecutionRecord {
 o.mu.RLock()
 defer o.mu.RUnlock()
 
 result := make([]ExecutionRecord, len(o.history))
 copy(result, o.history)
 return result
}
```

### 5.2 协同模式实现

```go
package choreography

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Event 事件接口
type Event interface {
 GetType() string
 GetData() map[string]interface{}
 GetSource() string
 GetTimestamp() time.Time
}

// BaseEvent 基础事件
type BaseEvent struct {
 Type      string
 Data      map[string]interface{}
 Source    string
 Timestamp time.Time
}

func (e *BaseEvent) GetType() string {
 return e.Type
}

func (e *BaseEvent) GetData() map[string]interface{} {
 return e.Data
}

func (e *BaseEvent) GetSource() string {
 return e.Source
}

func (e *BaseEvent) GetTimestamp() time.Time {
 return e.Timestamp
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) ([]Event, error)

// Service 服务接口
type Service interface {
 GetName() string
 HandleEvent(ctx context.Context, event Event) ([]Event, error)
 Subscribe(eventType string, handler EventHandler)
 Unsubscribe(eventType string)
}

// EventBus 事件总线
type EventBus struct {
 subscribers map[string][]Service
 events      chan Event
 mu          sync.RWMutex
 ctx         context.Context
 cancel      context.CancelFunc
 wg          sync.WaitGroup
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
 ctx, cancel := context.WithCancel(context.Background())
 eb := &EventBus{
  subscribers: make(map[string][]Service),
  events:      make(chan Event, 1000),
  ctx:         ctx,
  cancel:      cancel,
 }
 
 eb.wg.Add(1)
 go eb.processEvents()
 
 return eb
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType string, service Service) {
 eb.mu.Lock()
 defer eb.mu.Unlock()
 
 if eb.subscribers[eventType] == nil {
  eb.subscribers[eventType] = make([]Service, 0)
 }
 eb.subscribers[eventType] = append(eb.subscribers[eventType], service)
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(eventType string, service Service) {
 eb.mu.Lock()
 defer eb.mu.Unlock()
 
 subscribers := eb.subscribers[eventType]
 for i, sub := range subscribers {
  if sub.GetName() == service.GetName() {
   eb.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
   break
  }
 }
}

// Publish 发布事件
func (eb *EventBus) Publish(event Event) error {
 select {
 case eb.events <- event:
  return nil
 case <-time.After(5 * time.Second):
  return fmt.Errorf("发布事件超时")
 }
}

// processEvents 处理事件
func (eb *EventBus) processEvents() {
 defer eb.wg.Done()
 
 for {
  select {
  case event := <-eb.events:
   eb.handleEvent(event)
  case <-eb.ctx.Done():
   return
  }
 }
}

// handleEvent 处理单个事件
func (eb *EventBus) handleEvent(event Event) {
 eb.mu.RLock()
 subscribers := eb.subscribers[event.GetType()]
 eb.mu.RUnlock()
 
 for _, service := range subscribers {
  go func(s Service, e Event) {
   ctx, cancel := context.WithTimeout(eb.ctx, 30*time.Second)
   defer cancel()
   
   nextEvents, err := s.HandleEvent(ctx, e)
   if err != nil {
    fmt.Printf("服务 %s 处理事件 %s 失败: %v\n", s.GetName(), e.GetType(), err)
    return
   }
   
   // 发布后续事件
   for _, nextEvent := range nextEvents {
    eb.Publish(nextEvent)
   }
  }(service, event)
 }
}

// Stop 停止事件总线
func (eb *EventBus) Stop() {
 eb.cancel()
 eb.wg.Wait()
}

// BaseService 基础服务
type BaseService struct {
 name       string
 handlers   map[string]EventHandler
 eventBus   *EventBus
 mu         sync.RWMutex
}

// NewBaseService 创建基础服务
func NewBaseService(name string, eventBus *EventBus) *BaseService {
 return &BaseService{
  name:     name,
  handlers: make(map[string]EventHandler),
  eventBus: eventBus,
 }
}

// GetName 获取服务名称
func (s *BaseService) GetName() string {
 return s.name
}

// HandleEvent 处理事件
func (s *BaseService) HandleEvent(ctx context.Context, event Event) ([]Event, error) {
 s.mu.RLock()
 handler, exists := s.handlers[event.GetType()]
 s.mu.RUnlock()
 
 if !exists {
  return nil, fmt.Errorf("未找到事件处理器: %s", event.GetType())
 }
 
 return handler(ctx, event)
}

// Subscribe 订阅事件
func (s *BaseService) Subscribe(eventType string, handler EventHandler) {
 s.mu.Lock()
 s.handlers[eventType] = handler
 s.mu.Unlock()
 
 s.eventBus.Subscribe(eventType, s)
}

// Unsubscribe 取消订阅
func (s *BaseService) Unsubscribe(eventType string) {
 s.mu.Lock()
 delete(s.handlers, eventType)
 s.mu.Unlock()
 
 s.eventBus.Unsubscribe(eventType, s)
}

// Publish 发布事件
func (s *BaseService) Publish(event Event) error {
 return s.eventBus.Publish(event)
}
```

### 5.3 混合模式实现

```go
package hybrid

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// HybridOrchestrator 混合协调器
type HybridOrchestrator struct {
 orchestrator *orchestration.Orchestrator
 eventBus     *choreography.EventBus
 services     map[string]choreography.Service
 mu           sync.RWMutex
}

// NewHybridOrchestrator 创建混合协调器
func NewHybridOrchestrator(workflow *orchestration.Workflow) *HybridOrchestrator {
 return &HybridOrchestrator{
  orchestrator: orchestration.NewOrchestrator(workflow),
  eventBus:     choreography.NewEventBus(),
  services:     make(map[string]choreography.Service),
 }
}

// AddService 添加服务
func (h *HybridOrchestrator) AddService(service choreography.Service) {
 h.mu.Lock()
 defer h.mu.Unlock()
 h.services[service.GetName()] = service
}

// ExecuteOrchestration 执行编排模式
func (h *HybridOrchestrator) ExecuteOrchestration(ctx context.Context, initialData map[string]interface{}) error {
 return h.orchestrator.Execute(ctx, initialData)
}

// ExecuteChoreography 执行协同模式
func (h *HybridOrchestrator) ExecuteChoreography(ctx context.Context, initialEvent choreography.Event) error {
 return h.eventBus.Publish(initialEvent)
}

// BridgeOrchestrationToChoreography 从编排桥接到协同
func (h *HybridOrchestrator) BridgeOrchestrationToChoreography(orchestrationStep string, choreographyEvent choreography.Event) {
 // 当编排模式执行到特定步骤时，触发协同模式
 h.orchestrator.Subscribe(orchestrationStep, func(ctx context.Context, data map[string]interface{}) error {
  return h.eventBus.Publish(choreographyEvent)
 })
}

// BridgeChoreographyToOrchestration 从协同桥接到编排
func (h *HybridOrchestrator) BridgeChoreographyToOrchestration(eventType string, orchestrationStep string) {
 // 当协同模式接收到特定事件时，触发编排模式
 h.eventBus.Subscribe(eventType, &BridgeService{
  orchestrator: h.orchestrator,
  step:         orchestrationStep,
 })
}

// BridgeService 桥接服务
type BridgeService struct {
 orchestrator *orchestration.Orchestrator
 step         string
}

func (b *BridgeService) GetName() string {
 return "bridge_service"
}

func (b *BridgeService) HandleEvent(ctx context.Context, event choreography.Event) ([]choreography.Event, error) {
 // 触发编排模式的特定步骤
 data := event.GetData()
 // 这里需要根据具体业务逻辑来实现
 return nil, nil
}
```

## 6. 应用场景

### 6.1 订单处理系统

#### 编排模式实现

```go
package order

import (
 "context"
 "fmt"
 "time"
)

// OrderService 订单服务
type OrderService struct {
 name string
}

func (s *OrderService) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 orderID := data["orderID"].(string)
 fmt.Printf("处理订单: %s\n", orderID)
 
 data["orderProcessed"] = true
 return data, nil
}

func (s *OrderService) GetName() string {
 return s.name
}

func (s *OrderService) GetStatus() string {
 return "ready"
}

func (s *OrderService) Rollback(ctx context.Context, data map[string]interface{}) error {
 fmt.Printf("回滚订单处理: %s\n", data["orderID"])
 return nil
}

// PaymentService 支付服务
type PaymentService struct {
 name string
}

func (s *PaymentService) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 amount := data["amount"].(float64)
 fmt.Printf("处理支付: %.2f\n", amount)
 
 data["paymentProcessed"] = true
 return data, nil
}

func (s *PaymentService) GetName() string {
 return s.name
}

func (s *PaymentService) GetStatus() string {
 return "ready"
}

func (s *PaymentService) Rollback(ctx context.Context, data map[string]interface{}) error {
 fmt.Printf("回滚支付处理: %.2f\n", data["amount"])
 return nil
}

// ShippingService 发货服务
type ShippingService struct {
 name string
}

func (s *ShippingService) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 orderID := data["orderID"].(string)
 fmt.Printf("处理发货: %s\n", orderID)
 
 data["shipped"] = true
 return data, nil
}

func (s *ShippingService) GetName() string {
 return s.name
}

func (s *ShippingService) GetStatus() string {
 return "ready"
}

func (s *ShippingService) Rollback(ctx context.Context, data map[string]interface{}) error {
 fmt.Printf("回滚发货处理: %s\n", data["orderID"])
 return nil
}

// CreateOrderWorkflow 创建订单工作流
func CreateOrderWorkflow() *orchestration.Workflow {
 return &orchestration.Workflow{
  Name: "order_processing",
  Steps: map[string]*orchestration.WorkflowStep{
   "process_order": {
    Name:    "process_order",
    Service: &OrderService{name: "order_service"},
    Timeout: 30 * time.Second,
   },
   "process_payment": {
    Name:    "process_payment",
    Service: &PaymentService{name: "payment_service"},
    Timeout: 60 * time.Second,
   },
   "ship_order": {
    Name:    "ship_order",
    Service: &ShippingService{name: "shipping_service"},
    Timeout: 120 * time.Second,
   },
  },
  Start: "process_order",
  End:   []string{"ship_order"},
 }
}
```

#### 协同模式实现

```go
package order

import (
 "context"
 "fmt"
 "time"
)

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
 choreography.BaseEvent
}

// PaymentProcessedEvent 支付处理事件
type PaymentProcessedEvent struct {
 choreography.BaseEvent
}

// OrderShippedEvent 订单发货事件
type OrderShippedEvent struct {
 choreography.BaseEvent
}

// OrderService 订单服务（协同模式）
type OrderService struct {
 choreography.BaseService
}

// NewOrderService 创建订单服务
func NewOrderService(eventBus *choreography.EventBus) *OrderService {
 service := &OrderService{
  BaseService: *choreography.NewBaseService("order_service", eventBus),
 }
 
 service.Subscribe("order_created", service.handleOrderCreated)
 return service
}

// handleOrderCreated 处理订单创建事件
func (s *OrderService) handleOrderCreated(ctx context.Context, event choreography.Event) ([]choreography.Event, error) {
 orderID := event.GetData()["orderID"].(string)
 fmt.Printf("订单服务处理订单创建: %s\n", orderID)
 
 // 发布支付处理事件
 paymentEvent := &PaymentProcessedEvent{
  BaseEvent: choreography.BaseEvent{
   Type:      "payment_processed",
   Data:      event.GetData(),
   Source:    s.GetName(),
   Timestamp: time.Now(),
  },
 }
 
 return []choreography.Event{paymentEvent}, nil
}

// PaymentService 支付服务（协同模式）
type PaymentService struct {
 choreography.BaseService
}

// NewPaymentService 创建支付服务
func NewPaymentService(eventBus *choreography.EventBus) *PaymentService {
 service := &PaymentService{
  BaseService: *choreography.NewBaseService("payment_service", eventBus),
 }
 
 service.Subscribe("payment_processed", service.handlePaymentProcessed)
 return service
}

// handlePaymentProcessed 处理支付处理事件
func (s *PaymentService) handlePaymentProcessed(ctx context.Context, event choreography.Event) ([]choreography.Event, error) {
 amount := event.GetData()["amount"].(float64)
 fmt.Printf("支付服务处理支付: %.2f\n", amount)
 
 // 发布发货事件
 shippingEvent := &OrderShippedEvent{
  BaseEvent: choreography.BaseEvent{
   Type:      "order_shipped",
   Data:      event.GetData(),
   Source:    s.GetName(),
   Timestamp: time.Now(),
  },
 }
 
 return []choreography.Event{shippingEvent}, nil
}

// ShippingService 发货服务（协同模式）
type ShippingService struct {
 choreography.BaseService
}

// NewShippingService 创建发货服务
func NewShippingService(eventBus *choreography.EventBus) *ShippingService {
 service := &ShippingService{
  BaseService: *choreography.NewBaseService("shipping_service", eventBus),
 }
 
 service.Subscribe("order_shipped", service.handleOrderShipped)
 return service
}

// handleOrderShipped 处理订单发货事件
func (s *ShippingService) handleOrderShipped(ctx context.Context, event choreography.Event) ([]choreography.Event, error) {
 orderID := event.GetData()["orderID"].(string)
 fmt.Printf("发货服务处理发货: %s\n", orderID)
 
 return nil, nil
}
```

## 7. 性能分析

### 7.1 时间复杂度

#### 编排模式

- **任务调度**: $O(1)$
- **状态管理**: $O(1)$
- **错误处理**: $O(n)$，其中 $n$ 是步骤数量

#### 协同模式

- **事件发布**: $O(1)$
- **事件路由**: $O(m)$，其中 $m$ 是订阅者数量
- **事件处理**: $O(1)$ 每个事件

### 7.2 空间复杂度

#### 编排模式7

- **状态存储**: $O(n)$，其中 $n$ 是状态数量
- **历史记录**: $O(s)$，其中 $s$ 是步骤数量

#### 协同模式7

- **事件存储**: $O(e)$，其中 $e$ 是事件数量
- **订阅者存储**: $O(m)$，其中 $m$ 是订阅者数量

### 7.3 并发性能

#### 编排模式73

- **串行执行**: 步骤按顺序执行，并发性较低
- **并行步骤**: 支持并行执行的步骤

#### 协同模式73

- **并行处理**: 多个服务可以并行处理事件
- **异步通信**: 基于事件的异步通信

## 8. 最佳实践

### 8.1 编排模式最佳实践

1. **简单流程**: 适用于线性或简单的业务流程
2. **状态管理**: 集中管理状态，便于监控和调试
3. **错误处理**: 在协调器中统一处理错误和回滚

### 8.2 协同模式最佳实践

1. **复杂流程**: 适用于复杂的非线性业务流程
2. **服务自治**: 每个服务自主决定如何处理事件
3. **事件设计**: 设计清晰的事件类型和数据结构

### 8.3 选择指南

1. **团队结构**: 根据团队的组织结构选择合适的模式
2. **业务复杂度**: 根据业务流程的复杂度选择模式
3. **性能要求**: 根据性能要求选择合适的模式

## 9. 相关模式

### 9.1 状态机模式

编排模式可以看作是状态机模式的扩展，提供了更复杂的状态转换和条件判断。

### 9.2 观察者模式

协同模式可以看作是观察者模式的扩展，提供了更复杂的事件处理机制。

### 9.3 命令模式

编排模式可以使用命令模式来封装具体的业务操作。

---

**相关链接**:

- [01-状态机模式](../01-状态机模式/README.md)
- [02-工作流引擎模式](../02-工作流引擎模式/README.md)
- [03-任务队列模式](../03-任务队列模式/README.md)
- [返回上级目录](../../README.md)
