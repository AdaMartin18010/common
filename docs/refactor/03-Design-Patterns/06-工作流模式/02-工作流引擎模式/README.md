# 02-工作流引擎模式 (Workflow Engine Pattern)

## 目录

- [02-工作流引擎模式 (Workflow Engine Pattern)](#02-工作流引擎模式-workflow-engine-pattern)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心组件](#12-核心组件)
    - [1.3 模式结构](#13-模式结构)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 工作流数学模型](#21-工作流数学模型)
    - [2.2 工作流执行函数](#22-工作流执行函数)
    - [2.3 工作流状态转换](#23-工作流状态转换)
  - [3. 数学证明](#3-数学证明)
    - [3.1 工作流终止性定理](#31-工作流终止性定理)
    - [3.2 工作流一致性定理](#32-工作流一致性定理)
  - [4. 设计原则](#4-设计原则)
    - [4.1 单一职责原则](#41-单一职责原则)
    - [4.2 开闭原则](#42-开闭原则)
    - [4.3 依赖倒置原则](#43-依赖倒置原则)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 基础实现](#51-基础实现)
    - [5.2 泛型实现](#52-泛型实现)
    - [5.3 并发安全实现](#53-并发安全实现)
  - [6. 应用场景](#6-应用场景)
    - [6.1 订单处理工作流](#61-订单处理工作流)
    - [6.2 审批工作流](#62-审批工作流)
  - [7. 性能分析](#7-性能分析)
    - [7.1 时间复杂度](#71-时间复杂度)
    - [7.2 空间复杂度](#72-空间复杂度)
    - [7.3 并发性能](#73-并发性能)
  - [8. 最佳实践](#8-最佳实践)
    - [8.1 工作流设计原则](#81-工作流设计原则)
    - [8.2 性能优化](#82-性能优化)
    - [8.3 监控和调试](#83-监控和调试)
  - [9. 相关模式](#9-相关模式)
    - [9.1 状态机模式](#91-状态机模式)
    - [9.2 命令模式](#92-命令模式)
    - [9.3 观察者模式](#93-观察者模式)

## 1. 概念与定义

### 1.1 基本概念

工作流引擎模式是一种用于管理和执行复杂业务流程的设计模式。它提供了一个框架来定义、执行和监控工作流，支持条件分支、并行执行、错误处理等复杂场景。

**定义**: 工作流引擎模式提供了一个可配置的框架来定义和执行复杂的业务流程，支持工作流的定义、执行、监控和管理。

### 1.2 核心组件

- **WorkflowEngine (工作流引擎)**: 核心引擎，负责工作流的执行和管理
- **WorkflowDefinition (工作流定义)**: 定义工作流的结构和规则
- **Activity (活动)**: 工作流中的具体执行单元
- **WorkflowInstance (工作流实例)**: 工作流的具体执行实例
- **WorkflowContext (工作流上下文)**: 存储工作流执行过程中的数据

### 1.3 模式结构

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ WorkflowEngine  │    │WorkflowDefinition│    │    Activity     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + execute()     │◄──►│ + activities    │◄──►│ + execute()     │
│ + pause()       │    │ + transitions   │    │ + validate()    │
│ + resume()      │    │ + conditions    │    │ + rollback()    │
│ + cancel()      │    └─────────────────┘    └─────────────────┘
└─────────────────┘              ▲
                                 │
                    ┌─────────────────┐
                    │WorkflowInstance │
                    ├─────────────────┤
                    │ + state         │
                    │ + context       │
                    │ + history       │
                    └─────────────────┘
```

## 2. 形式化定义

### 2.1 工作流数学模型

设 $W = (A, T, C, I, F)$ 为一个工作流，其中：

- $A = \{a_1, a_2, ..., a_n\}$ 是活动集合
- $T = \{t_1, t_2, ..., t_m\}$ 是转换集合
- $C: T \rightarrow \mathbb{B}$ 是条件函数，$\mathbb{B} = \{true, false\}$
- $I \in A$ 是初始活动
- $F \subseteq A$ 是终止活动集合

### 2.2 工作流执行函数

对于活动 $a \in A$，执行函数定义为：

$$execute(a, context) = (result, next\_activities)$$

其中：

- $result$ 是执行结果
- $next\_activities$ 是下一个要执行的活动集合

### 2.3 工作流状态转换

工作流状态转换函数定义为：

$$\delta: (A \times Context) \rightarrow (A' \times Context')$$

其中 $A' \subseteq A$ 是下一个活动集合，$Context'$ 是更新后的上下文。

## 3. 数学证明

### 3.1 工作流终止性定理

**定理**: 如果工作流 $W$ 是有限且无环的，则工作流执行必然终止。

**证明**:

1. 设工作流 $W$ 包含 $n$ 个活动
2. 每次执行最多访问 $n$ 个活动
3. 由于无环，每个活动最多被访问一次
4. 因此，工作流执行最多进行 $n$ 步
5. 结论：工作流执行必然终止

### 3.2 工作流一致性定理

**定理**: 如果工作流定义是有效的，则所有执行路径都会产生一致的结果。

**证明**:

1. 设工作流定义 $W$ 是有效的
2. 对于任意两个执行路径 $P_1$ 和 $P_2$
3. 由于工作流定义的一致性约束
4. $P_1$ 和 $P_2$ 必须满足相同的业务规则
5. 因此，两个路径产生一致的结果

## 4. 设计原则

### 4.1 单一职责原则

每个组件只负责特定的功能，工作流引擎负责执行，活动负责具体业务逻辑。

### 4.2 开闭原则

可以通过扩展活动类型来支持新的业务逻辑，而不需要修改引擎核心代码。

### 4.3 依赖倒置原则

工作流引擎依赖于抽象的活动接口，而不是具体的活动实现。

## 5. Go语言实现

### 5.1 基础实现

```go
package workflow

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Activity 活动接口
type Activity interface {
 Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
 GetName() string
 GetType() string
}

// Transition 转换定义
type Transition struct {
 From       string
 To         string
 Condition  func(data map[string]interface{}) bool
}

// WorkflowDefinition 工作流定义
type WorkflowDefinition struct {
 Name        string
 Activities  map[string]Activity
 Transitions []Transition
 Initial     string
 Final       []string
}

// WorkflowInstance 工作流实例
type WorkflowInstance struct {
 ID           string
 Definition   *WorkflowDefinition
 CurrentState string
 Context      map[string]interface{}
 History      []ExecutionStep
 Status       string
 mu           sync.RWMutex
}

// ExecutionStep 执行步骤
type ExecutionStep struct {
 ActivityName string
 StartTime    time.Time
 EndTime      time.Time
 Result       map[string]interface{}
 Error        error
}

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
 instances map[string]*WorkflowInstance
 mu        sync.RWMutex
}

// NewWorkflowEngine 创建工作流引擎
func NewWorkflowEngine() *WorkflowEngine {
 return &WorkflowEngine{
  instances: make(map[string]*WorkflowInstance),
 }
}

// CreateInstance 创建工作流实例
func (e *WorkflowEngine) CreateInstance(definition *WorkflowDefinition, initialData map[string]interface{}) (*WorkflowInstance, error) {
 instance := &WorkflowInstance{
  ID:           generateID(),
  Definition:   definition,
  CurrentState: definition.Initial,
  Context:      initialData,
  History:      make([]ExecutionStep, 0),
  Status:       "created",
 }
 
 e.mu.Lock()
 e.instances[instance.ID] = instance
 e.mu.Unlock()
 
 return instance, nil
}

// Execute 执行工作流
func (e *WorkflowEngine) Execute(ctx context.Context, instanceID string) error {
 e.mu.RLock()
 instance, exists := e.instances[instanceID]
 e.mu.RUnlock()
 
 if !exists {
  return fmt.Errorf("工作流实例 %s 不存在", instanceID)
 }
 
 instance.mu.Lock()
 instance.Status = "running"
 instance.mu.Unlock()
 
 defer func() {
  instance.mu.Lock()
  if instance.Status == "running" {
   instance.Status = "completed"
  }
  instance.mu.Unlock()
 }()
 
 for {
  select {
  case <-ctx.Done():
   instance.mu.Lock()
   instance.Status = "cancelled"
   instance.mu.Unlock()
   return ctx.Err()
  default:
  }
  
  instance.mu.RLock()
  currentState := instance.CurrentState
  context := instance.Context
  instance.mu.RUnlock()
  
  // 检查是否到达终止状态
  if e.isFinalState(instance, currentState) {
   break
  }
  
  // 执行当前活动
  activity, exists := instance.Definition.Activities[currentState]
  if !exists {
   return fmt.Errorf("活动 %s 不存在", currentState)
  }
  
  step := ExecutionStep{
   ActivityName: currentState,
   StartTime:    time.Now(),
  }
  
  result, err := activity.Execute(ctx, context)
  step.EndTime = time.Now()
  step.Result = result
  step.Error = err
  
  instance.mu.Lock()
  instance.History = append(instance.History, step)
  instance.Context = result
  instance.mu.Unlock()
  
  if err != nil {
   instance.mu.Lock()
   instance.Status = "failed"
   instance.mu.Unlock()
   return fmt.Errorf("活动 %s 执行失败: %v", currentState, err)
  }
  
  // 确定下一个状态
  nextState := e.determineNextState(instance, currentState, result)
  if nextState == "" {
   instance.mu.Lock()
   instance.Status = "completed"
   instance.mu.Unlock()
   break
  }
  
  instance.mu.Lock()
  instance.CurrentState = nextState
  instance.mu.Unlock()
 }
 
 return nil
}

// isFinalState 检查是否为终止状态
func (e *WorkflowEngine) isFinalState(instance *WorkflowInstance, state string) bool {
 for _, final := range instance.Definition.Final {
  if final == state {
   return true
  }
 }
 return false
}

// determineNextState 确定下一个状态
func (e *WorkflowEngine) determineNextState(instance *WorkflowInstance, currentState string, data map[string]interface{}) string {
 for _, transition := range instance.Definition.Transitions {
  if transition.From == currentState {
   if transition.Condition == nil || transition.Condition(data) {
    return transition.To
   }
  }
 }
 return ""
}

// generateID 生成唯一ID
func generateID() string {
 return fmt.Sprintf("wf_%d", time.Now().UnixNano())
}
```

### 5.2 泛型实现

```go
package workflow

import (
 "context"
 "fmt"
 "sync"
)

// WorkflowEngine[T] 泛型工作流引擎
type WorkflowEngine[T any] struct {
 instances map[string]*WorkflowInstance[T]
 mu        sync.RWMutex
}

// WorkflowInstance[T] 泛型工作流实例
type WorkflowInstance[T any] struct {
 ID           string
 Definition   *WorkflowDefinition[T]
 CurrentState string
 Context      T
 History      []ExecutionStep[T]
 Status       string
 mu           sync.RWMutex
}

// ExecutionStep[T] 泛型执行步骤
type ExecutionStep[T any] struct {
 ActivityName string
 StartTime    time.Time
 EndTime      time.Time
 Result       T
 Error        error
}

// Activity[T] 泛型活动接口
type Activity[T any] interface {
 Execute(ctx context.Context, data T) (T, error)
 GetName() string
 GetType() string
}

// WorkflowDefinition[T] 泛型工作流定义
type WorkflowDefinition[T any] struct {
 Name        string
 Activities  map[string]Activity[T]
 Transitions []Transition[T]
 Initial     string
 Final       []string
}

// Transition[T] 泛型转换定义
type Transition[T any] struct {
 From       string
 To         string
 Condition  func(data T) bool
}

// NewWorkflowEngine[T] 创建泛型工作流引擎
func NewWorkflowEngine[T any]() *WorkflowEngine[T] {
 return &WorkflowEngine[T]{
  instances: make(map[string]*WorkflowInstance[T]),
 }
}

// CreateInstance[T] 创建泛型工作流实例
func (e *WorkflowEngine[T]) CreateInstance(definition *WorkflowDefinition[T], initialData T) (*WorkflowInstance[T], error) {
 instance := &WorkflowInstance[T]{
  ID:           generateID(),
  Definition:   definition,
  CurrentState: definition.Initial,
  Context:      initialData,
  History:      make([]ExecutionStep[T], 0),
  Status:       "created",
 }
 
 e.mu.Lock()
 e.instances[instance.ID] = instance
 e.mu.Unlock()
 
 return instance, nil
}

// Execute[T] 执行泛型工作流
func (e *WorkflowEngine[T]) Execute(ctx context.Context, instanceID string) error {
 e.mu.RLock()
 instance, exists := e.instances[instanceID]
 e.mu.RUnlock()
 
 if !exists {
  return fmt.Errorf("工作流实例 %s 不存在", instanceID)
 }
 
 instance.mu.Lock()
 instance.Status = "running"
 instance.mu.Unlock()
 
 defer func() {
  instance.mu.Lock()
  if instance.Status == "running" {
   instance.Status = "completed"
  }
  instance.mu.Unlock()
 }()
 
 for {
  select {
  case <-ctx.Done():
   instance.mu.Lock()
   instance.Status = "cancelled"
   instance.mu.Unlock()
   return ctx.Err()
  default:
  }
  
  instance.mu.RLock()
  currentState := instance.CurrentState
  context := instance.Context
  instance.mu.RUnlock()
  
  // 检查是否到达终止状态
  if e.isFinalState(instance, currentState) {
   break
  }
  
  // 执行当前活动
  activity, exists := instance.Definition.Activities[currentState]
  if !exists {
   return fmt.Errorf("活动 %s 不存在", currentState)
  }
  
  step := ExecutionStep[T]{
   ActivityName: currentState,
   StartTime:    time.Now(),
  }
  
  result, err := activity.Execute(ctx, context)
  step.EndTime = time.Now()
  step.Result = result
  step.Error = err
  
  instance.mu.Lock()
  instance.History = append(instance.History, step)
  instance.Context = result
  instance.mu.Unlock()
  
  if err != nil {
   instance.mu.Lock()
   instance.Status = "failed"
   instance.mu.Unlock()
   return fmt.Errorf("活动 %s 执行失败: %v", currentState, err)
  }
  
  // 确定下一个状态
  nextState := e.determineNextState(instance, currentState, result)
  if nextState == "" {
   instance.mu.Lock()
   instance.Status = "completed"
   instance.mu.Unlock()
   break
  }
  
  instance.mu.Lock()
  instance.CurrentState = nextState
  instance.mu.Unlock()
 }
 
 return nil
}

// isFinalState[T] 检查是否为终止状态
func (e *WorkflowEngine[T]) isFinalState(instance *WorkflowInstance[T], state string) bool {
 for _, final := range instance.Definition.Final {
  if final == state {
   return true
  }
 }
 return false
}

// determineNextState[T] 确定下一个状态
func (e *WorkflowEngine[T]) determineNextState(instance *WorkflowInstance[T], currentState string, data T) string {
 for _, transition := range instance.Definition.Transitions {
  if transition.From == currentState {
   if transition.Condition == nil || transition.Condition(data) {
    return transition.To
   }
  }
 }
 return ""
}
```

### 5.3 并发安全实现

```go
package workflow

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// ConcurrentWorkflowEngine 并发工作流引擎
type ConcurrentWorkflowEngine struct {
 instances map[string]*WorkflowInstance
 mu        sync.RWMutex
 workers   chan struct{}
 eventChan chan WorkflowEvent
 stopChan  chan struct{}
 wg        sync.WaitGroup
}

// WorkflowEvent 工作流事件
type WorkflowEvent struct {
 Type       string
 InstanceID string
 Data       interface{}
 Response   chan WorkflowResponse
}

// WorkflowResponse 工作流响应
type WorkflowResponse struct {
 Result interface{}
 Error  error
}

// NewConcurrentWorkflowEngine 创建并发工作流引擎
func NewConcurrentWorkflowEngine(maxWorkers int) *ConcurrentWorkflowEngine {
 cwe := &ConcurrentWorkflowEngine{
  instances: make(map[string]*WorkflowInstance),
  workers:   make(chan struct{}, maxWorkers),
  eventChan: make(chan WorkflowEvent, 100),
  stopChan:  make(chan struct{}),
 }
 
 cwe.wg.Add(1)
 go cwe.eventLoop()
 
 return cwe
}

// eventLoop 事件循环
func (cwe *ConcurrentWorkflowEngine) eventLoop() {
 defer cwe.wg.Done()
 
 for {
  select {
  case event := <-cwe.eventChan:
   cwe.handleEvent(event)
  case <-cwe.stopChan:
   return
  }
 }
}

// handleEvent 处理事件
func (cwe *ConcurrentWorkflowEngine) handleEvent(event WorkflowEvent) {
 switch event.Type {
 case "execute":
  cwe.workers <- struct{}{} // 获取工作协程
  go func() {
   defer func() { <-cwe.workers }() // 释放工作协程
   cwe.executeWorkflow(event)
  }()
 case "pause":
  cwe.pauseWorkflow(event)
 case "resume":
  cwe.resumeWorkflow(event)
 case "cancel":
  cwe.cancelWorkflow(event)
 }
}

// executeWorkflow 执行工作流
func (cwe *ConcurrentWorkflowEngine) executeWorkflow(event WorkflowEvent) {
 cwe.mu.RLock()
 instance, exists := cwe.instances[event.InstanceID]
 cwe.mu.RUnlock()
 
 if !exists {
  response := WorkflowResponse{
   Error: fmt.Errorf("工作流实例 %s 不存在", event.InstanceID),
  }
  event.Response <- response
  return
 }
 
 ctx := context.Background()
 err := cwe.executeInstance(ctx, instance)
 
 response := WorkflowResponse{
  Error: err,
 }
 event.Response <- response
}

// executeInstance 执行实例
func (cwe *ConcurrentWorkflowEngine) executeInstance(ctx context.Context, instance *WorkflowInstance) error {
 instance.mu.Lock()
 instance.Status = "running"
 instance.mu.Unlock()
 
 defer func() {
  instance.mu.Lock()
  if instance.Status == "running" {
   instance.Status = "completed"
  }
  instance.mu.Unlock()
 }()
 
 for {
  select {
  case <-ctx.Done():
   instance.mu.Lock()
   instance.Status = "cancelled"
   instance.mu.Unlock()
   return ctx.Err()
  default:
  }
  
  instance.mu.RLock()
  currentState := instance.CurrentState
  context := instance.Context
  instance.mu.RUnlock()
  
  // 检查是否到达终止状态
  if cwe.isFinalState(instance, currentState) {
   break
  }
  
  // 执行当前活动
  activity, exists := instance.Definition.Activities[currentState]
  if !exists {
   return fmt.Errorf("活动 %s 不存在", currentState)
  }
  
  step := ExecutionStep{
   ActivityName: currentState,
   StartTime:    time.Now(),
  }
  
  result, err := activity.Execute(ctx, context)
  step.EndTime = time.Now()
  step.Result = result
  step.Error = err
  
  instance.mu.Lock()
  instance.History = append(instance.History, step)
  instance.Context = result
  instance.mu.Unlock()
  
  if err != nil {
   instance.mu.Lock()
   instance.Status = "failed"
   instance.mu.Unlock()
   return fmt.Errorf("活动 %s 执行失败: %v", currentState, err)
  }
  
  // 确定下一个状态
  nextState := cwe.determineNextState(instance, currentState, result)
  if nextState == "" {
   instance.mu.Lock()
   instance.Status = "completed"
   instance.mu.Unlock()
   break
  }
  
  instance.mu.Lock()
  instance.CurrentState = nextState
  instance.mu.Unlock()
 }
 
 return nil
}

// isFinalState 检查是否为终止状态
func (cwe *ConcurrentWorkflowEngine) isFinalState(instance *WorkflowInstance, state string) bool {
 for _, final := range instance.Definition.Final {
  if final == state {
   return true
  }
 }
 return false
}

// determineNextState 确定下一个状态
func (cwe *ConcurrentWorkflowEngine) determineNextState(instance *WorkflowInstance, currentState string, data map[string]interface{}) string {
 for _, transition := range instance.Definition.Transitions {
  if transition.From == currentState {
   if transition.Condition == nil || transition.Condition(data) {
    return transition.To
   }
  }
 }
 return ""
}

// pauseWorkflow 暂停工作流
func (cwe *ConcurrentWorkflowEngine) pauseWorkflow(event WorkflowEvent) {
 cwe.mu.RLock()
 instance, exists := cwe.instances[event.InstanceID]
 cwe.mu.RUnlock()
 
 if !exists {
  response := WorkflowResponse{
   Error: fmt.Errorf("工作流实例 %s 不存在", event.InstanceID),
  }
  event.Response <- response
  return
 }
 
 instance.mu.Lock()
 instance.Status = "paused"
 instance.mu.Unlock()
 
 response := WorkflowResponse{}
 event.Response <- response
}

// resumeWorkflow 恢复工作流
func (cwe *ConcurrentWorkflowEngine) resumeWorkflow(event WorkflowEvent) {
 cwe.mu.RLock()
 instance, exists := cwe.instances[event.InstanceID]
 cwe.mu.RUnlock()
 
 if !exists {
  response := WorkflowResponse{
   Error: fmt.Errorf("工作流实例 %s 不存在", event.InstanceID),
  }
  event.Response <- response
  return
 }
 
 instance.mu.Lock()
 instance.Status = "running"
 instance.mu.Unlock()
 
 response := WorkflowResponse{}
 event.Response <- response
}

// cancelWorkflow 取消工作流
func (cwe *ConcurrentWorkflowEngine) cancelWorkflow(event WorkflowEvent) {
 cwe.mu.RLock()
 instance, exists := cwe.instances[event.InstanceID]
 cwe.mu.RUnlock()
 
 if !exists {
  response := WorkflowResponse{
   Error: fmt.Errorf("工作流实例 %s 不存在", event.InstanceID),
  }
  event.Response <- response
  return
 }
 
 instance.mu.Lock()
 instance.Status = "cancelled"
 instance.mu.Unlock()
 
 response := WorkflowResponse{}
 event.Response <- response
}

// SendEvent 发送事件
func (cwe *ConcurrentWorkflowEngine) SendEvent(eventType, instanceID string, data interface{}) error {
 responseChan := make(chan WorkflowResponse, 1)
 
 event := WorkflowEvent{
  Type:       eventType,
  InstanceID: instanceID,
  Data:       data,
  Response:   responseChan,
 }
 
 select {
 case cwe.eventChan <- event:
 case <-time.After(5 * time.Second):
  return fmt.Errorf("发送事件超时")
 }
 
 select {
 case response := <-responseChan:
  return response.Error
 case <-time.After(10 * time.Second):
  return fmt.Errorf("等待响应超时")
 }
}

// Stop 停止引擎
func (cwe *ConcurrentWorkflowEngine) Stop() {
 close(cwe.stopChan)
 cwe.wg.Wait()
}
```

## 6. 应用场景

### 6.1 订单处理工作流

```go
package order

import (
 "context"
 "fmt"
 "time"
)

// OrderData 订单数据
type OrderData struct {
 OrderID     string
 CustomerID  string
 Amount      float64
 Status      string
 CreatedAt   time.Time
 UpdatedAt   time.Time
}

// ValidateOrderActivity 订单验证活动
type ValidateOrderActivity struct{}

func (a *ValidateOrderActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 orderID := data["orderID"].(string)
 fmt.Printf("验证订单: %s\n", orderID)
 
 // 模拟验证逻辑
 if orderID == "" {
  return data, fmt.Errorf("订单ID不能为空")
 }
 
 data["validated"] = true
 return data, nil
}

func (a *ValidateOrderActivity) GetName() string {
 return "validate_order"
}

func (a *ValidateOrderActivity) GetType() string {
 return "validation"
}

// ProcessPaymentActivity 支付处理活动
type ProcessPaymentActivity struct{}

func (a *ProcessPaymentActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 amount := data["amount"].(float64)
 fmt.Printf("处理支付: %.2f\n", amount)
 
 // 模拟支付处理
 if amount > 1000 {
  return data, fmt.Errorf("金额过大，需要人工审核")
 }
 
 data["paymentProcessed"] = true
 return data, nil
}

func (a *ProcessPaymentActivity) GetName() string {
 return "process_payment"
}

func (a *ProcessPaymentActivity) GetType() string {
 return "payment"
}

// ShipOrderActivity 发货活动
type ShipOrderActivity struct{}

func (a *ShipOrderActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 orderID := data["orderID"].(string)
 fmt.Printf("发货订单: %s\n", orderID)
 
 data["shipped"] = true
 data["trackingNumber"] = fmt.Sprintf("TRK_%s", orderID)
 return data, nil
}

func (a *ShipOrderActivity) GetName() string {
 return "ship_order"
}

func (a *ShipOrderActivity) GetType() string {
 return "shipping"
}

// CreateOrderWorkflow 创建订单工作流
func CreateOrderWorkflow() *WorkflowDefinition {
 definition := &WorkflowDefinition{
  Name: "order_processing",
  Activities: map[string]Activity{
   "validate_order":   &ValidateOrderActivity{},
   "process_payment":  &ProcessPaymentActivity{},
   "ship_order":       &ShipOrderActivity{},
  },
  Transitions: []Transition{
   {
    From: "validate_order",
    To:   "process_payment",
    Condition: func(data map[string]interface{}) bool {
     return data["validated"] == true
    },
   },
   {
    From: "process_payment",
    To:   "ship_order",
    Condition: func(data map[string]interface{}) bool {
     return data["paymentProcessed"] == true
    },
   },
  },
  Initial: "validate_order",
  Final:   []string{"ship_order"},
 }
 
 return definition
}
```

### 6.2 审批工作流

```go
package approval

import (
 "context"
 "fmt"
)

// ApprovalData 审批数据
type ApprovalData struct {
 RequestID   string
 RequesterID string
 Amount      float64
 Type        string
 Status      string
 Approvers   []string
 CurrentStep int
}

// SubmitRequestActivity 提交申请活动
type SubmitRequestActivity struct{}

func (a *SubmitRequestActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 requestID := data["requestID"].(string)
 fmt.Printf("提交申请: %s\n", requestID)
 
 data["submitted"] = true
 data["currentStep"] = 0
 return data, nil
}

func (a *SubmitRequestActivity) GetName() string {
 return "submit_request"
}

func (a *SubmitRequestActivity) GetType() string {
 return "submission"
}

// ManagerApprovalActivity 经理审批活动
type ManagerApprovalActivity struct{}

func (a *ManagerApprovalActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 requestID := data["requestID"].(string)
 amount := data["amount"].(float64)
 fmt.Printf("经理审批申请: %s, 金额: %.2f\n", requestID, amount)
 
 // 模拟审批逻辑
 if amount > 5000 {
  data["managerApproved"] = false
  data["rejectionReason"] = "金额超过经理审批权限"
 } else {
  data["managerApproved"] = true
 }
 
 return data, nil
}

func (a *ManagerApprovalActivity) GetName() string {
 return "manager_approval"
}

func (a *ManagerApprovalActivity) GetType() string {
 return "approval"
}

// DirectorApprovalActivity 总监审批活动
type DirectorApprovalActivity struct{}

func (a *DirectorApprovalActivity) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
 requestID := data["requestID"].(string)
 amount := data["amount"].(float64)
 fmt.Printf("总监审批申请: %s, 金额: %.2f\n", requestID, amount)
 
 // 模拟审批逻辑
 if amount > 50000 {
  data["directorApproved"] = false
  data["rejectionReason"] = "金额超过总监审批权限"
 } else {
  data["directorApproved"] = true
 }
 
 return data, nil
}

func (a *DirectorApprovalActivity) GetName() string {
 return "director_approval"
}

func (a *DirectorApprovalActivity) GetType() string {
 return "approval"
}

// CreateApprovalWorkflow 创建审批工作流
func CreateApprovalWorkflow() *WorkflowDefinition {
 definition := &WorkflowDefinition{
  Name: "approval_workflow",
  Activities: map[string]Activity{
   "submit_request":     &SubmitRequestActivity{},
   "manager_approval":   &ManagerApprovalActivity{},
   "director_approval":  &DirectorApprovalActivity{},
  },
  Transitions: []Transition{
   {
    From: "submit_request",
    To:   "manager_approval",
    Condition: func(data map[string]interface{}) bool {
     return data["submitted"] == true
    },
   },
   {
    From: "manager_approval",
    To:   "director_approval",
    Condition: func(data map[string]interface{}) bool {
     amount := data["amount"].(float64)
     return data["managerApproved"] == true && amount > 5000
    },
   },
  },
  Initial: "submit_request",
  Final:   []string{"manager_approval", "director_approval"},
 }
 
 return definition
}
```

## 7. 性能分析

### 7.1 时间复杂度

- **工作流创建**: $O(1)$
- **活动执行**: $O(1)$ 每个活动
- **状态转换**: $O(n)$，其中 $n$ 是转换数量
- **工作流执行**: $O(m \times n)$，其中 $m$ 是活动数量，$n$ 是转换数量

### 7.2 空间复杂度

- **工作流定义**: $O(a + t)$，其中 $a$ 是活动数量，$t$ 是转换数量
- **工作流实例**: $O(h)$，其中 $h$ 是历史记录数量
- **上下文数据**: $O(d)$，其中 $d$ 是数据大小

### 7.3 并发性能

- **并行执行**: 支持多个工作流实例并行执行
- **资源管理**: 使用工作协程池控制并发数量
- **事件驱动**: 异步事件处理，避免阻塞

## 8. 最佳实践

### 8.1 工作流设计原则

1. **单一职责**: 每个活动只负责一个特定的业务功能
2. **可重用性**: 设计可重用的活动组件
3. **可测试性**: 每个活动都应该可以独立测试
4. **错误处理**: 明确定义错误处理和回滚策略

### 8.2 性能优化

1. **活动缓存**: 缓存频繁使用的活动对象
2. **并行执行**: 利用并行活动提高执行效率
3. **资源池**: 使用连接池和对象池减少资源开销

### 8.3 监控和调试

1. **执行历史**: 记录详细的执行历史
2. **性能指标**: 监控执行时间和资源使用
3. **错误追踪**: 提供详细的错误信息和堆栈跟踪

## 9. 相关模式

### 9.1 状态机模式

工作流引擎可以看作是状态机模式的扩展，提供了更复杂的状态转换和条件判断。

### 9.2 命令模式

活动可以看作是命令模式的实现，每个活动封装了一个具体的业务操作。

### 9.3 观察者模式

工作流引擎可以使用观察者模式来通知状态变化和事件发生。

---

**相关链接**:

- [01-状态机模式](../01-状态机模式/README.md)
- [03-任务队列模式](../03-任务队列模式/README.md)
- [04-编排vs协同模式](../04-编排vs协同模式/README.md)
- [返回上级目录](../../README.md)
