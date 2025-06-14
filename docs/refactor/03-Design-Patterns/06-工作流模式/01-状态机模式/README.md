# 01-状态机模式 (State Machine Pattern)

## 目录

- [01-状态机模式 (State Machine Pattern)](#01-状态机模式-state-machine-pattern)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心组件](#12-核心组件)
    - [1.3 模式结构](#13-模式结构)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 状态机数学模型](#21-状态机数学模型)
    - [2.2 状态转移函数](#22-状态转移函数)
    - [2.3 扩展转移函数](#23-扩展转移函数)
  - [3. 数学证明](#3-数学证明)
    - [3.1 状态机确定性定理](#31-状态机确定性定理)
    - [3.2 状态机可达性定理](#32-状态机可达性定理)
  - [4. 设计原则](#4-设计原则)
    - [4.1 单一职责原则](#41-单一职责原则)
    - [4.2 开闭原则](#42-开闭原则)
    - [4.3 里氏替换原则](#43-里氏替换原则)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 基础实现](#51-基础实现)
    - [5.2 泛型实现](#52-泛型实现)
    - [5.3 函数式实现](#53-函数式实现)
    - [5.4 并发安全实现](#54-并发安全实现)
  - [6. 应用场景](#6-应用场景)
    - [6.1 订单处理系统](#61-订单处理系统)
    - [6.2 游戏状态管理](#62-游戏状态管理)
  - [7. 性能分析](#7-性能分析)
    - [7.1 时间复杂度](#71-时间复杂度)
    - [7.2 空间复杂度](#72-空间复杂度)
    - [7.3 并发性能](#73-并发性能)
  - [8. 最佳实践](#8-最佳实践)
    - [8.1 状态设计原则](#81-状态设计原则)
    - [8.2 性能优化](#82-性能优化)
    - [8.3 错误处理](#83-错误处理)
  - [9. 相关模式](#9-相关模式)
    - [9.1 策略模式](#91-策略模式)
    - [9.2 命令模式](#92-命令模式)
    - [9.3 观察者模式](#93-观察者模式)

## 1. 概念与定义

### 1.1 基本概念

状态机模式是一种行为型设计模式，它允许对象在内部状态改变时改变其行为。对象看起来似乎修改了它的类。

**定义**: 状态机模式使对象能够在其内部状态改变时改变其行为，对象看起来好像修改了它的类。

### 1.2 核心组件

- **Context (上下文)**: 维护当前状态，并将状态相关的请求委托给当前状态对象
- **State (状态)**: 定义状态接口，封装与特定状态相关的行为
- **ConcreteState (具体状态)**: 实现状态接口，定义特定状态下的行为

### 1.3 模式结构

```text
┌─────────────────┐    ┌─────────────────┐
│     Context     │    │      State      │
├─────────────────┤    ├─────────────────┤
│ - state: State  │◄──►│ + handle()      │
├─────────────────┤    └─────────────────┘
│ + request()     │              ▲
│ + setState()    │              │
└─────────────────┘    ┌─────────────────┐
                       │ ConcreteStateA  │
                       ├─────────────────┤
                       │ + handle()      │
                       └─────────────────┘
```

## 2. 形式化定义

### 2.1 状态机数学模型

设 $M = (Q, \Sigma, \delta, q_0, F)$ 为一个有限状态机，其中：

- $Q$ 是有限状态集合
- $\Sigma$ 是输入字母表
- $\delta: Q \times \Sigma \rightarrow Q$ 是状态转移函数
- $q_0 \in Q$ 是初始状态
- $F \subseteq Q$ 是接受状态集合

### 2.2 状态转移函数

对于任意状态 $q \in Q$ 和输入 $a \in \Sigma$，状态转移函数定义为：

$$\delta(q, a) = q'$$

其中 $q' \in Q$ 是转移后的状态。

### 2.3 扩展转移函数

对于字符串 $w = a_1a_2...a_n$，扩展转移函数定义为：

$$\hat{\delta}(q, w) = \delta(\delta(...\delta(\delta(q, a_1), a_2)...), a_n)$$

## 3. 数学证明

### 3.1 状态机确定性定理

**定理**: 对于确定性有限状态机，给定当前状态和输入，下一个状态是唯一确定的。

**证明**:

1. 设 $M = (Q, \Sigma, \delta, q_0, F)$ 为确定性有限状态机
2. 对于任意 $q \in Q$ 和 $a \in \Sigma$，$\delta(q, a)$ 返回唯一的状态 $q' \in Q$
3. 因此，状态转移是确定性的

### 3.2 状态机可达性定理

**定理**: 从初始状态 $q_0$ 出发，通过有限次状态转移可以到达状态 $q$，当且仅当存在字符串 $w \in \Sigma^*$ 使得 $\hat{\delta}(q_0, w) = q$。

**证明**:

1. **必要性**: 如果状态 $q$ 可达，则存在转移序列 $q_0 \xrightarrow{a_1} q_1 \xrightarrow{a_2} ... \xrightarrow{a_n} q$
2. **充分性**: 如果存在字符串 $w = a_1a_2...a_n$ 使得 $\hat{\delta}(q_0, w) = q$，则状态 $q$ 可达

## 4. 设计原则

### 4.1 单一职责原则

每个状态类只负责处理特定状态下的行为，符合单一职责原则。

### 4.2 开闭原则

添加新状态时不需要修改现有代码，符合开闭原则。

### 4.3 里氏替换原则

所有具体状态类都可以替换状态接口，符合里氏替换原则。

## 5. Go语言实现

### 5.1 基础实现

```go
package statemachine

import (
 "fmt"
 "sync"
)

// State 定义状态接口
type State interface {
 Handle(context *Context) error
 GetName() string
}

// Context 上下文，维护当前状态
type Context struct {
 state State
 mu    sync.RWMutex
}

// NewContext 创建新的上下文
func NewContext(initialState State) *Context {
 return &Context{
  state: initialState,
 }
}

// SetState 设置状态
func (c *Context) SetState(state State) {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.state = state
}

// GetState 获取当前状态
func (c *Context) GetState() State {
 c.mu.RLock()
 defer c.mu.RUnlock()
 return c.state
}

// Request 处理请求
func (c *Context) Request() error {
 c.mu.RLock()
 state := c.state
 c.mu.RUnlock()
 return state.Handle(c)
}

// ConcreteStateA 具体状态A
type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle(context *Context) error {
 fmt.Println("状态A处理请求")
 // 状态转换逻辑
 context.SetState(&ConcreteStateB{})
 return nil
}

func (s *ConcreteStateA) GetName() string {
 return "StateA"
}

// ConcreteStateB 具体状态B
type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle(context *Context) error {
 fmt.Println("状态B处理请求")
 // 状态转换逻辑
 context.SetState(&ConcreteStateC{})
 return nil
}

func (s *ConcreteStateB) GetName() string {
 return "StateB"
}

// ConcreteStateC 具体状态C
type ConcreteStateC struct{}

func (s *ConcreteStateC) Handle(context *Context) error {
 fmt.Println("状态C处理请求")
 // 状态转换逻辑
 context.SetState(&ConcreteStateA{})
 return nil
}

func (s *ConcreteStateC) GetName() string {
 return "StateC"
}
```

### 5.2 泛型实现

```go
package statemachine

import (
 "context"
 "fmt"
 "sync"
)

// StateMachine 泛型状态机
type StateMachine[T any] struct {
 currentState State[T]
 states       map[string]State[T]
 mu           sync.RWMutex
}

// State 泛型状态接口
type State[T any] interface {
 Handle(ctx context.Context, data T) (T, error)
 GetName() string
 CanTransitionTo(to string) bool
}

// NewStateMachine 创建新的状态机
func NewStateMachine[T any](initialState State[T]) *StateMachine[T] {
 sm := &StateMachine[T]{
  currentState: initialState,
  states:       make(map[string]State[T]),
 }
 sm.states[initialState.GetName()] = initialState
 return sm
}

// AddState 添加状态
func (sm *StateMachine[T]) AddState(state State[T]) {
 sm.mu.Lock()
 defer sm.mu.Unlock()
 sm.states[state.GetName()] = state
}

// SetState 设置当前状态
func (sm *StateMachine[T]) SetState(stateName string) error {
 sm.mu.Lock()
 defer sm.mu.Unlock()
 
 state, exists := sm.states[stateName]
 if !exists {
  return fmt.Errorf("状态 %s 不存在", stateName)
 }
 
 if !sm.currentState.CanTransitionTo(stateName) {
  return fmt.Errorf("不能从状态 %s 转换到状态 %s", 
   sm.currentState.GetName(), stateName)
 }
 
 sm.currentState = state
 return nil
}

// Process 处理数据
func (sm *StateMachine[T]) Process(ctx context.Context, data T) (T, error) {
 sm.mu.RLock()
 state := sm.currentState
 sm.mu.RUnlock()
 
 return state.Handle(ctx, data)
}

// GetCurrentState 获取当前状态
func (sm *StateMachine[T]) GetCurrentState() string {
 sm.mu.RLock()
 defer sm.mu.RUnlock()
 return sm.currentState.GetName()
}
```

### 5.3 函数式实现

```go
package statemachine

import (
 "context"
 "fmt"
)

// StateFunc 状态函数类型
type StateFunc[T any] func(ctx context.Context, data T) (T, StateFunc[T], error)

// FunctionalStateMachine 函数式状态机
type FunctionalStateMachine[T any] struct {
 currentState StateFunc[T]
 transitions  map[string][]string
}

// NewFunctionalStateMachine 创建函数式状态机
func NewFunctionalStateMachine[T any](initialState StateFunc[T]) *FunctionalStateMachine[T] {
 return &FunctionalStateMachine[T]{
  currentState: initialState,
  transitions:  make(map[string][]string),
 }
}

// AddTransition 添加状态转换规则
func (sm *FunctionalStateMachine[T]) AddTransition(from, to string) {
 if sm.transitions[from] == nil {
  sm.transitions[from] = make([]string, 0)
 }
 sm.transitions[from] = append(sm.transitions[from], to)
}

// Process 处理数据
func (sm *FunctionalStateMachine[T]) Process(ctx context.Context, data T) (T, error) {
 var err error
 for i := 0; i < 100; i++ { // 防止无限循环
  data, sm.currentState, err = sm.currentState(ctx, data)
  if err != nil {
   return data, err
  }
  if sm.currentState == nil {
   break
  }
 }
 return data, nil
}
```

### 5.4 并发安全实现

```go
package statemachine

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// ConcurrentStateMachine 并发安全状态机
type ConcurrentStateMachine[T any] struct {
 currentState State[T]
 states       map[string]State[T]
 mu           sync.RWMutex
 eventChan    chan StateEvent[T]
 stopChan     chan struct{}
 wg           sync.WaitGroup
}

// StateEvent 状态事件
type StateEvent[T any] struct {
 EventType string
 Data      T
 Response  chan StateEventResponse[T]
}

// StateEventResponse 状态事件响应
type StateEventResponse[T any] struct {
 Data T
 Err  error
}

// NewConcurrentStateMachine 创建并发状态机
func NewConcurrentStateMachine[T any](initialState State[T]) *ConcurrentStateMachine[T] {
 csm := &ConcurrentStateMachine[T]{
  currentState: initialState,
  states:       make(map[string]State[T]),
  eventChan:    make(chan StateEvent[T], 100),
  stopChan:     make(chan struct{}),
 }
 csm.states[initialState.GetName()] = initialState
 
 csm.wg.Add(1)
 go csm.eventLoop()
 
 return csm
}

// eventLoop 事件循环
func (csm *ConcurrentStateMachine[T]) eventLoop() {
 defer csm.wg.Done()
 
 for {
  select {
  case event := <-csm.eventChan:
   csm.handleEvent(event)
  case <-csm.stopChan:
   return
  }
 }
}

// handleEvent 处理事件
func (csm *ConcurrentStateMachine[T]) handleEvent(event StateEvent[T]) {
 csm.mu.Lock()
 defer csm.mu.Unlock()
 
 ctx := context.Background()
 data, err := csm.currentState.Handle(ctx, event.Data)
 
 response := StateEventResponse[T]{
  Data: data,
  Err:  err,
 }
 
 select {
 case event.Response <- response:
 default:
 }
}

// SendEvent 发送事件
func (csm *ConcurrentStateMachine[T]) SendEvent(eventType string, data T) (T, error) {
 responseChan := make(chan StateEventResponse[T], 1)
 
 event := StateEvent[T]{
  EventType: eventType,
  Data:      data,
  Response:  responseChan,
 }
 
 select {
 case csm.eventChan <- event:
 case <-time.After(5 * time.Second):
  return data, fmt.Errorf("发送事件超时")
 }
 
 select {
 case response := <-responseChan:
  return response.Data, response.Err
 case <-time.After(10 * time.Second):
  return data, fmt.Errorf("等待响应超时")
 }
}

// Stop 停止状态机
func (csm *ConcurrentStateMachine[T]) Stop() {
 close(csm.stopChan)
 csm.wg.Wait()
}
```

## 6. 应用场景

### 6.1 订单处理系统

```go
package order

import (
 "context"
 "fmt"
 "time"
)

// OrderState 订单状态接口
type OrderState interface {
 Process(ctx context.Context, order *Order) error
 GetName() string
}

// Order 订单
type Order struct {
 ID        string
 Status    string
 Amount    float64
 CreatedAt time.Time
 UpdatedAt time.Time
}

// PendingState 待处理状态
type PendingState struct{}

func (s *PendingState) Process(ctx context.Context, order *Order) error {
 fmt.Printf("订单 %s 处于待处理状态\n", order.ID)
 order.Status = "processing"
 order.UpdatedAt = time.Now()
 return nil
}

func (s *PendingState) GetName() string {
 return "pending"
}

// ProcessingState 处理中状态
type ProcessingState struct{}

func (s *ProcessingState) Process(ctx context.Context, order *Order) error {
 fmt.Printf("订单 %s 正在处理中\n", order.ID)
 order.Status = "completed"
 order.UpdatedAt = time.Now()
 return nil
}

func (s *ProcessingState) GetName() string {
 return "processing"
}

// CompletedState 已完成状态
type CompletedState struct{}

func (s *CompletedState) Process(ctx context.Context, order *Order) error {
 fmt.Printf("订单 %s 已完成\n", order.ID)
 return nil
}

func (s *CompletedState) GetName() string {
 return "completed"
}
```

### 6.2 游戏状态管理

```go
package game

import (
 "context"
 "fmt"
)

// GameState 游戏状态接口
type GameState interface {
 Update(ctx context.Context, game *Game) error
 HandleInput(ctx context.Context, game *Game, input string) error
 GetName() string
}

// Game 游戏
type Game struct {
 PlayerHealth int
 Score        int
 Level        int
 State        GameState
}

// MenuState 菜单状态
type MenuState struct{}

func (s *MenuState) Update(ctx context.Context, game *Game) error {
 fmt.Println("游戏菜单 - 按 's' 开始游戏，按 'q' 退出")
 return nil
}

func (s *MenuState) HandleInput(ctx context.Context, game *Game, input string) error {
 switch input {
 case "s":
  game.State = &PlayingState{}
  fmt.Println("游戏开始！")
 case "q":
  fmt.Println("退出游戏")
  return fmt.Errorf("游戏退出")
 }
 return nil
}

func (s *MenuState) GetName() string {
 return "menu"
}

// PlayingState 游戏进行状态
type PlayingState struct{}

func (s *PlayingState) Update(ctx context.Context, game *Game) error {
 fmt.Printf("游戏进行中 - 生命值: %d, 分数: %d, 等级: %d\n", 
  game.PlayerHealth, game.Score, game.Level)
 return nil
}

func (s *PlayingState) HandleInput(ctx context.Context, game *Game, input string) error {
 switch input {
 case "p":
  game.State = &PausedState{}
  fmt.Println("游戏暂停")
 case "q":
  game.State = &MenuState{}
  fmt.Println("返回菜单")
 }
 return nil
}

func (s *PlayingState) GetName() string {
 return "playing"
}

// PausedState 暂停状态
type PausedState struct{}

func (s *PausedState) Update(ctx context.Context, game *Game) error {
 fmt.Println("游戏暂停 - 按 'r' 继续游戏")
 return nil
}

func (s *PausedState) HandleInput(ctx context.Context, game *Game, input string) error {
 if input == "r" {
  game.State = &PlayingState{}
  fmt.Println("游戏继续")
 }
 return nil
}

func (s *PausedState) GetName() string {
 return "paused"
}
```

## 7. 性能分析

### 7.1 时间复杂度

- **状态转换**: $O(1)$
- **状态查找**: $O(1)$ (使用哈希表)
- **事件处理**: $O(1)$

### 7.2 空间复杂度

- **状态存储**: $O(n)$，其中 $n$ 是状态数量
- **上下文存储**: $O(1)$

### 7.3 并发性能

- **读操作**: 支持并发读取
- **写操作**: 需要互斥锁保护
- **事件处理**: 异步处理，避免阻塞

## 8. 最佳实践

### 8.1 状态设计原则

1. **状态单一职责**: 每个状态只负责处理特定状态下的行为
2. **状态转换明确**: 明确定义状态转换条件和规则
3. **状态可测试**: 每个状态都应该可以独立测试

### 8.2 性能优化

1. **状态缓存**: 缓存频繁使用的状态对象
2. **事件批处理**: 批量处理相同类型的事件
3. **内存池**: 使用对象池减少内存分配

### 8.3 错误处理

1. **状态转换验证**: 验证状态转换的合法性
2. **异常状态处理**: 处理异常状态和错误恢复
3. **超时机制**: 设置状态转换超时时间

## 9. 相关模式

### 9.1 策略模式

状态模式可以看作是策略模式的扩展，策略模式关注算法的选择，而状态模式关注状态的变化。

### 9.2 命令模式

命令模式可以用于封装状态转换操作，使状态转换更加灵活。

### 9.3 观察者模式

观察者模式可以用于通知状态变化，实现状态变化的监听和响应。

---

**相关链接**:

- [02-工作流引擎模式](../02-工作流引擎模式/README.md)
- [03-任务队列模式](../03-任务队列模式/README.md)
- [04-编排vs协同模式](../04-编排vs协同模式/README.md)
- [返回上级目录](../../README.md)
