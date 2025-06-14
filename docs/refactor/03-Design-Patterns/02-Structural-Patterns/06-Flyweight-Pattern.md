# 06-享元模式 (Flyweight Pattern)

## 目录

1. [概述](#1-概述)
2. [形式化定义](#2-形式化定义)
3. [数学基础](#3-数学基础)
4. [系统架构](#4-系统架构)
5. [核心算法](#5-核心算法)
6. [Go语言实现](#6-go语言实现)
7. [性能优化](#7-性能优化)
8. [安全考虑](#8-安全考虑)
9. [总结](#9-总结)

## 1. 概述

### 1.1 定义

享元模式（Flyweight Pattern）是一种结构型设计模式，通过共享技术有效地支持大量细粒度对象的复用。

**形式化定义**：

```latex
F = (I, E, S, M, C)
```

其中：

- I：内部状态集合（Intrinsic State）
- E：外部状态集合（Extrinsic State）
- S：共享对象集合（Shared Objects）
- M：内存映射（Memory Mapping）
- C：上下文管理器（Context Manager）

### 1.2 核心概念

| 概念 | 定义 | 数学表示 |
|------|------|----------|
| 内部状态 | 对象固有的、可共享的状态 | I = {i₁, i₂, ..., iₙ} |
| 外部状态 | 对象特有的、不可共享的状态 | E = {e₁, e₂, ..., eₘ} |
| 享元工厂 | 创建和管理享元对象的工厂 | Factory: I → S |
| 对象池 | 存储共享对象的容器 | Pool ⊆ S |

### 1.3 模式结构

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client        │    │  Flyweight      │    │  Flyweight      │
│                 │    │  Factory        │    │                 │
│ - extrinsic     │───▶│                 │───▶│ - intrinsic     │
│   state         │    │ - getFlyweight  │    │ - operation     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Concrete       │
                       │  Flyweight      │
                       │                 │
                       │ - intrinsic     │
                       │ - operation     │
                       └─────────────────┘
```

## 2. 形式化定义

### 2.1 享元空间

**定义 2.1** 享元空间是一个五元组 (I, E, S, F, C)：

- I：内部状态集合，I = {i₁, i₂, ..., iₙ}
- E：外部状态集合，E = {e₁, e₂, ..., eₘ}
- S：共享对象集合，S = {s₁, s₂, ..., sₖ}
- F：工厂函数，F: I → S
- C：上下文集合，C = {c₁, c₂, ..., cₗ}

**公理 2.1** 内部状态不变性：

```latex
∀s ∈ S, ∀i ∈ I : s.intrinsic = i ⇒ s.intrinsic = i
```

**公理 2.2** 外部状态可变性：

```latex
∀s ∈ S, ∀e ∈ E : s.extrinsic = e ⇒ s.extrinsic ≠ e
```

### 2.2 享元函数

**定义 2.2** 享元函数 flyweight: I × E → O 满足：

1. **共享性**：∀i ∈ I : flyweight(i, e₁) = flyweight(i, e₂)
2. **唯一性**：∀i₁, i₂ ∈ I : i₁ ≠ i₂ ⇒ flyweight(i₁, e) ≠ flyweight(i₂, e)
3. **可逆性**：∀s ∈ S : ∃i ∈ I : s = flyweight(i, e)

### 2.3 内存映射

**定义 2.3** 内存映射函数 M: I → S 满足：

1. **单射性**：∀i₁, i₂ ∈ I : i₁ ≠ i₂ ⇒ M(i₁) ≠ M(i₂)
2. **满射性**：∀s ∈ S : ∃i ∈ I : M(i) = s
3. **双射性**：M是双射函数

**定理 2.1** 内存节省定理：

```latex
memory_saved = |I| × (object_size - reference_size)
```

**证明**：

```latex
设每个对象大小为 object_size，引用大小为 reference_size
共享前总内存：|I| × |E| × object_size
共享后总内存：|I| × object_size + |E| × reference_size
节省内存：|I| × |E| × object_size - (|I| × object_size + |E| × reference_size)
         = |I| × (|E| × object_size - object_size - |E| × reference_size / |I|)
         ≈ |I| × (object_size - reference_size) (当|E| >> 1时)
```

## 3. 数学基础

### 3.1 集合论基础

**定义 3.1** 享元集合的基数：

```latex
|S| = |I| ≤ |O|
```

其中：

- |S|：共享对象数量
- |I|：内部状态数量
- |O|：总对象数量

**定理 3.1** 内存效率定理：

```latex
efficiency = |S| / |O| = |I| / (|I| × |E|) = 1 / |E|
```

### 3.2 图论基础

**定义 3.2** 享元图 G = (V, E)：

- V：对象节点集合
- E：共享关系边集合

**定理 3.2** 享元图的连通性：

```latex
∀v₁, v₂ ∈ V : ∃path(v₁, v₂) ⇒ v₁.intrinsic = v₂.intrinsic
```

### 3.3 复杂度分析

**定理 3.3** 享元模式的时间复杂度：

```latex
T(n) = O(1)  // 获取享元对象
S(n) = O(|I|) // 空间复杂度
```

**证明**：

```latex
获取享元对象：哈希表查找 O(1)
创建享元对象：最多创建 |I| 个对象
总时间复杂度：O(1)
总空间复杂度：O(|I|)
```

## 4. 系统架构

### 4.1 分层架构

```latex
┌─────────────────────────────────────┐
│            Client Layer             │
├─────────────────────────────────────┤
│         Flyweight Factory           │
├─────────────────────────────────────┤
│         Flyweight Pool              │
├─────────────────────────────────────┤
│         Concrete Flyweight          │
├─────────────────────────────────────┤
│         Context Manager             │
└─────────────────────────────────────┘
```

### 4.2 组件设计

#### 4.2.1 享元工厂

```go
type FlyweightFactory struct {
    flyweights map[string]Flyweight
    mu         sync.RWMutex
}

type Flyweight interface {
    Operation(extrinsicState string)
    GetIntrinsicState() string
}
```

#### 4.2.2 享元池

```go
type FlyweightPool struct {
    pool map[string]Flyweight
    mu   sync.RWMutex
    max  int
}
```

## 5. 核心算法

### 5.1 享元获取算法

**算法 5.1** 享元对象获取：

```go
func (f *FlyweightFactory) GetFlyweight(intrinsicState string) Flyweight {
    f.mu.RLock()
    if flyweight, exists := f.flyweights[intrinsicState]; exists {
        f.mu.RUnlock()
        return flyweight
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查
    if flyweight, exists := f.flyweights[intrinsicState]; exists {
        return flyweight
    }
    
    // 创建新的享元对象
    flyweight := NewConcreteFlyweight(intrinsicState)
    f.flyweights[intrinsicState] = flyweight
    
    return flyweight
}
```

**复杂度分析**：

- 时间复杂度：O(1) 平均情况
- 空间复杂度：O(|I|)

### 5.2 享元池管理算法

**算法 5.2** LRU享元池管理：

```go
type LRUFlyweightPool struct {
    pool    map[string]*FlyweightNode
    head    *FlyweightNode
    tail    *FlyweightNode
    max     int
    mu      sync.RWMutex
}

type FlyweightNode struct {
    key      string
    flyweight Flyweight
    prev     *FlyweightNode
    next     *FlyweightNode
}

func (p *LRUFlyweightPool) Get(key string) Flyweight {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if node, exists := p.pool[key]; exists {
        // 移动到头部
        p.moveToHead(node)
        return node.flyweight
    }
    
    return nil
}

func (p *LRUFlyweightPool) Put(key string, flyweight Flyweight) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if node, exists := p.pool[key]; exists {
        node.flyweight = flyweight
        p.moveToHead(node)
        return
    }
    
    // 创建新节点
    node := &FlyweightNode{
        key:       key,
        flyweight: flyweight,
    }
    
    p.pool[key] = node
    p.addToHead(node)
    
    // 检查容量
    if len(p.pool) > p.max {
        p.removeTail()
    }
}
```

### 5.3 内存优化算法

**算法 5.3** 享元内存优化：

```go
func OptimizeMemoryUsage(flyweights []Flyweight) map[string]int {
    optimization := make(map[string]int)
    
    for _, flyweight := range flyweights {
        intrinsicState := flyweight.GetIntrinsicState()
        optimization[intrinsicState]++
    }
    
    return optimization
}
```

## 6. Go语言实现

### 6.1 基础数据结构

```go
package flyweight

import (
    "fmt"
    "sync"
)

// Flyweight 享元接口
type Flyweight interface {
    Operation(extrinsicState string)
    GetIntrinsicState() string
}

// IntrinsicState 内部状态
type IntrinsicState struct {
    Color    string `json:"color"`
    Size     int    `json:"size"`
    Font     string `json:"font"`
    Style    string `json:"style"`
}

// ExtrinsicState 外部状态
type ExtrinsicState struct {
    X        int    `json:"x"`
    Y        int    `json:"y"`
    Text     string `json:"text"`
    Visible  bool   `json:"visible"`
}

// ConcreteFlyweight 具体享元
type ConcreteFlyweight struct {
    intrinsicState IntrinsicState
    mu             sync.RWMutex
}

// NewConcreteFlyweight 创建具体享元
func NewConcreteFlyweight(intrinsicState IntrinsicState) *ConcreteFlyweight {
    return &ConcreteFlyweight{
        intrinsicState: intrinsicState,
    }
}

// Operation 操作
func (f *ConcreteFlyweight) Operation(extrinsicState string) {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    fmt.Printf("Flyweight with intrinsic state %+v, extrinsic state: %s\n", 
        f.intrinsicState, extrinsicState)
}

// GetIntrinsicState 获取内部状态
func (f *ConcreteFlyweight) GetIntrinsicState() string {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    return fmt.Sprintf("%s_%d_%s_%s", 
        f.intrinsicState.Color, 
        f.intrinsicState.Size, 
        f.intrinsicState.Font, 
        f.intrinsicState.Style)
}

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
    flyweights map[string]Flyweight
    mu         sync.RWMutex
}

// NewFlyweightFactory 创建享元工厂
func NewFlyweightFactory() *FlyweightFactory {
    return &FlyweightFactory{
        flyweights: make(map[string]Flyweight),
    }
}

// GetFlyweight 获取享元对象
func (f *FlyweightFactory) GetFlyweight(intrinsicState IntrinsicState) Flyweight {
    key := f.getKey(intrinsicState)
    
    f.mu.RLock()
    if flyweight, exists := f.flyweights[key]; exists {
        f.mu.RUnlock()
        return flyweight
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查锁定
    if flyweight, exists := f.flyweights[key]; exists {
        return flyweight
    }
    
    // 创建新的享元对象
    flyweight := NewConcreteFlyweight(intrinsicState)
    f.flyweights[key] = flyweight
    
    fmt.Printf("Created new flyweight: %s\n", key)
    
    return flyweight
}

// getKey 生成键
func (f *FlyweightFactory) getKey(intrinsicState IntrinsicState) string {
    return fmt.Sprintf("%s_%d_%s_%s", 
        intrinsicState.Color, 
        intrinsicState.Size, 
        intrinsicState.Font, 
        intrinsicState.Style)
}

// GetFlyweightCount 获取享元对象数量
func (f *FlyweightFactory) GetFlyweightCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.flyweights)
}

// ListFlyweights 列出所有享元对象
func (f *FlyweightFactory) ListFlyweights() []string {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    var keys []string
    for key := range f.flyweights {
        keys = append(keys, key)
    }
    
    return keys
}
```

### 6.2 享元池实现

```go
// FlyweightPool 享元池
type FlyweightPool struct {
    pool map[string]Flyweight
    mu   sync.RWMutex
    max  int
}

// NewFlyweightPool 创建享元池
func NewFlyweightPool(max int) *FlyweightPool {
    return &FlyweightPool{
        pool: make(map[string]Flyweight),
        max:  max,
    }
}

// Get 获取享元对象
func (p *FlyweightPool) Get(key string) Flyweight {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    if flyweight, exists := p.pool[key]; exists {
        return flyweight
    }
    
    return nil
}

// Put 放入享元对象
func (p *FlyweightPool) Put(key string, flyweight Flyweight) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if len(p.pool) >= p.max {
        return fmt.Errorf("pool is full")
    }
    
    p.pool[key] = flyweight
    return nil
}

// Remove 移除享元对象
func (p *FlyweightPool) Remove(key string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    delete(p.pool, key)
}

// Size 获取池大小
func (p *FlyweightPool) Size() int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return len(p.pool)
}

// Clear 清空池
func (p *FlyweightPool) Clear() {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.pool = make(map[string]Flyweight)
}

// LRUFlyweightPool LRU享元池
type LRUFlyweightPool struct {
    pool map[string]*FlyweightNode
    head *FlyweightNode
    tail *FlyweightNode
    max  int
    mu   sync.RWMutex
}

// FlyweightNode 享元节点
type FlyweightNode struct {
    key       string
    flyweight Flyweight
    prev      *FlyweightNode
    next      *FlyweightNode
}

// NewLRUFlyweightPool 创建LRU享元池
func NewLRUFlyweightPool(max int) *LRUFlyweightPool {
    pool := &LRUFlyweightPool{
        pool: make(map[string]*FlyweightNode),
        max:  max,
    }
    
    // 初始化头尾节点
    pool.head = &FlyweightNode{}
    pool.tail = &FlyweightNode{}
    pool.head.next = pool.tail
    pool.tail.prev = pool.head
    
    return pool
}

// Get 获取享元对象
func (p *LRUFlyweightPool) Get(key string) Flyweight {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if node, exists := p.pool[key]; exists {
        // 移动到头部
        p.moveToHead(node)
        return node.flyweight
    }
    
    return nil
}

// Put 放入享元对象
func (p *LRUFlyweightPool) Put(key string, flyweight Flyweight) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if node, exists := p.pool[key]; exists {
        node.flyweight = flyweight
        p.moveToHead(node)
        return
    }
    
    // 创建新节点
    node := &FlyweightNode{
        key:       key,
        flyweight: flyweight,
    }
    
    p.pool[key] = node
    p.addToHead(node)
    
    // 检查容量
    if len(p.pool) > p.max {
        p.removeTail()
    }
}

// moveToHead 移动到头部
func (p *LRUFlyweightPool) moveToHead(node *FlyweightNode) {
    p.removeNode(node)
    p.addToHead(node)
}

// addToHead 添加到头部
func (p *LRUFlyweightPool) addToHead(node *FlyweightNode) {
    node.prev = p.head
    node.next = p.head.next
    p.head.next.prev = node
    p.head.next = node
}

// removeNode 移除节点
func (p *LRUFlyweightPool) removeNode(node *FlyweightNode) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

// removeTail 移除尾部节点
func (p *LRUFlyweightPool) removeTail() {
    node := p.tail.prev
    p.removeNode(node)
    delete(p.pool, node.key)
}

// Size 获取池大小
func (p *LRUFlyweightPool) Size() int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    return len(p.pool)
}
```

### 6.3 上下文管理器

```go
// Context 上下文
type Context struct {
    extrinsicState ExtrinsicState
    flyweight      Flyweight
}

// NewContext 创建上下文
func NewContext(extrinsicState ExtrinsicState, flyweight Flyweight) *Context {
    return &Context{
        extrinsicState: extrinsicState,
        flyweight:      flyweight,
    }
}

// Execute 执行操作
func (c *Context) Execute() {
    extrinsicStateStr := fmt.Sprintf("x=%d,y=%d,text=%s,visible=%t", 
        c.extrinsicState.X, 
        c.extrinsicState.Y, 
        c.extrinsicState.Text, 
        c.extrinsicState.Visible)
    
    c.flyweight.Operation(extrinsicStateStr)
}

// ContextManager 上下文管理器
type ContextManager struct {
    contexts map[string]*Context
    mu       sync.RWMutex
}

// NewContextManager 创建上下文管理器
func NewContextManager() *ContextManager {
    return &ContextManager{
        contexts: make(map[string]*Context),
    }
}

// AddContext 添加上下文
func (m *ContextManager) AddContext(id string, context *Context) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.contexts[id] = context
}

// GetContext 获取上下文
func (m *ContextManager) GetContext(id string) *Context {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    if context, exists := m.contexts[id]; exists {
        return context
    }
    
    return nil
}

// RemoveContext 移除上下文
func (m *ContextManager) RemoveContext(id string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.contexts, id)
}

// ExecuteAll 执行所有上下文
func (m *ContextManager) ExecuteAll() {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for _, context := range m.contexts {
        context.Execute()
    }
}
```

### 6.4 应用示例

```go
// TextEditor 文本编辑器示例
type TextEditor struct {
    factory *FlyweightFactory
    manager *ContextManager
}

// NewTextEditor 创建文本编辑器
func NewTextEditor() *TextEditor {
    return &TextEditor{
        factory: NewFlyweightFactory(),
        manager: NewContextManager(),
    }
}

// AddCharacter 添加字符
func (e *TextEditor) AddCharacter(char rune, x, y int, color, font string, size int) {
    // 创建内部状态
    intrinsicState := IntrinsicState{
        Color: color,
        Size:  size,
        Font:  font,
        Style: "normal",
    }
    
    // 获取享元对象
    flyweight := e.factory.GetFlyweight(intrinsicState)
    
    // 创建外部状态
    extrinsicState := ExtrinsicState{
        X:       x,
        Y:       y,
        Text:    string(char),
        Visible: true,
    }
    
    // 创建上下文
    context := NewContext(extrinsicState, flyweight)
    
    // 添加到管理器
    id := fmt.Sprintf("%d_%d", x, y)
    e.manager.AddContext(id, context)
}

// Render 渲染
func (e *TextEditor) Render() {
    fmt.Printf("Rendering text editor with %d flyweights and %d contexts\n", 
        e.factory.GetFlyweightCount(), len(e.manager.contexts))
    
    e.manager.ExecuteAll()
}

// MemoryUsage 内存使用情况
func (e *TextEditor) MemoryUsage() {
    flyweightCount := e.factory.GetFlyweightCount()
    contextCount := len(e.manager.contexts)
    
    fmt.Printf("Memory usage:\n")
    fmt.Printf("  Flyweights: %d\n", flyweightCount)
    fmt.Printf("  Contexts: %d\n", contextCount)
    fmt.Printf("  Memory saved: %d objects\n", contextCount-flyweightCount)
}
```

## 7. 性能优化

### 7.1 并发优化

```go
// ConcurrentFlyweightFactory 并发享元工厂
type ConcurrentFlyweightFactory struct {
    factory *FlyweightFactory
    workers int
    jobQueue chan FlyweightJob
}

// FlyweightJob 享元任务
type FlyweightJob struct {
    IntrinsicState IntrinsicState
    Result         chan Flyweight
}

// NewConcurrentFlyweightFactory 创建并发享元工厂
func NewConcurrentFlyweightFactory(workers int) *ConcurrentFlyweightFactory {
    factory := &ConcurrentFlyweightFactory{
        factory:  NewFlyweightFactory(),
        workers:  workers,
        jobQueue: make(chan FlyweightJob, 1000),
    }
    
    // 启动工作协程
    for i := 0; i < workers; i++ {
        go factory.worker()
    }
    
    return factory
}

// worker 工作协程
func (f *ConcurrentFlyweightFactory) worker() {
    for job := range f.jobQueue {
        flyweight := f.factory.GetFlyweight(job.IntrinsicState)
        job.Result <- flyweight
    }
}

// GetFlyweightAsync 异步获取享元对象
func (f *ConcurrentFlyweightFactory) GetFlyweightAsync(intrinsicState IntrinsicState) <-chan Flyweight {
    result := make(chan Flyweight, 1)
    
    job := FlyweightJob{
        IntrinsicState: intrinsicState,
        Result:         result,
    }
    
    f.jobQueue <- job
    
    return result
}
```

### 7.2 内存优化

```go
// MemoryOptimizedFlyweightFactory 内存优化的享元工厂
type MemoryOptimizedFlyweightFactory struct {
    factory *FlyweightFactory
    pool    *sync.Pool
}

// NewMemoryOptimizedFlyweightFactory 创建内存优化的享元工厂
func NewMemoryOptimizedFlyweightFactory() *MemoryOptimizedFlyweightFactory {
    return &MemoryOptimizedFlyweightFactory{
        factory: NewFlyweightFactory(),
        pool: &sync.Pool{
            New: func() interface{} {
                return &IntrinsicState{}
            },
        },
    }
}

// GetFlyweight 获取享元对象（内存优化）
func (f *MemoryOptimizedFlyweightFactory) GetFlyweight(intrinsicState IntrinsicState) Flyweight {
    // 从对象池获取状态对象
    pooledState := f.pool.Get().(*IntrinsicState)
    defer f.pool.Put(pooledState)
    
    // 复制状态
    *pooledState = intrinsicState
    
    return f.factory.GetFlyweight(*pooledState)
}
```

## 8. 安全考虑

### 8.1 线程安全

```go
// ThreadSafeFlyweightFactory 线程安全的享元工厂
type ThreadSafeFlyweightFactory struct {
    factory *FlyweightFactory
    rwmu    sync.RWMutex
}

// NewThreadSafeFlyweightFactory 创建线程安全的享元工厂
func NewThreadSafeFlyweightFactory() *ThreadSafeFlyweightFactory {
    return &ThreadSafeFlyweightFactory{
        factory: NewFlyweightFactory(),
    }
}

// GetFlyweight 获取享元对象（线程安全）
func (f *ThreadSafeFlyweightFactory) GetFlyweight(intrinsicState IntrinsicState) Flyweight {
    f.rwmu.RLock()
    flyweight := f.factory.GetFlyweight(intrinsicState)
    f.rwmu.RUnlock()
    
    return flyweight
}
```

### 8.2 内存泄漏防护

```go
// MemoryLeakProtectedFlyweightFactory 内存泄漏防护的享元工厂
type MemoryLeakProtectedFlyweightFactory struct {
    factory *FlyweightFactory
    pool    *LRUFlyweightPool
    max     int
}

// NewMemoryLeakProtectedFlyweightFactory 创建内存泄漏防护的享元工厂
func NewMemoryLeakProtectedFlyweightFactory(max int) *MemoryLeakProtectedFlyweightFactory {
    return &MemoryLeakProtectedFlyweightFactory{
        factory: NewFlyweightFactory(),
        pool:    NewLRUFlyweightPool(max),
        max:     max,
    }
}

// GetFlyweight 获取享元对象（内存泄漏防护）
func (f *MemoryLeakProtectedFlyweightFactory) GetFlyweight(intrinsicState IntrinsicState) Flyweight {
    key := f.factory.getKey(intrinsicState)
    
    // 先从池中获取
    if flyweight := f.pool.Get(key); flyweight != nil {
        return flyweight
    }
    
    // 从工厂获取
    flyweight := f.factory.GetFlyweight(intrinsicState)
    
    // 放入池中
    f.pool.Put(key, flyweight)
    
    return flyweight
}
```

## 9. 总结

### 9.1 核心特性

1. **形式化定义**：基于数学公理的享元模式体系
2. **内存优化**：通过共享减少内存使用
3. **性能提升**：减少对象创建和销毁开销
4. **线程安全**：支持并发访问
5. **可扩展性**：支持多种享元池策略

### 9.2 应用场景

- **图形系统**：字符、图形对象共享
- **游戏开发**：游戏对象、纹理共享
- **文本处理**：字符、字体共享
- **缓存系统**：对象缓存、数据共享

### 9.3 扩展方向

1. **分布式享元**：跨进程、跨机器共享
2. **持久化享元**：磁盘存储、数据库缓存
3. **智能享元**：机器学习优化、自适应策略
4. **可视化工具**：享元使用情况监控

---

**相关链接**：

- [01-适配器模式](./01-Adapter-Pattern.md)
- [02-桥接模式](./02-Bridge-Pattern.md)
- [03-组合模式](./03-Composite-Pattern.md)
- [04-装饰器模式](./04-Decorator-Pattern.md)
- [05-外观模式](./05-Facade-Pattern.md)
- [07-代理模式](./07-Proxy-Pattern.md)
