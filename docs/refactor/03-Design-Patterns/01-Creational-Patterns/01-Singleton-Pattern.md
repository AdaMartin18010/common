# 01-单例模式 (Singleton Pattern)

## 目录

- [01-单例模式 (Singleton Pattern)](#01-单例模式-singleton-pattern)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心特征](#12-核心特征)
    - [1.3 设计原则](#13-设计原则)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 集合论定义](#21-集合论定义)
    - [2.2 状态机定义](#22-状态机定义)
    - [2.3 形式化约束](#23-形式化约束)
  - [3. 数学证明](#3-数学证明)
    - [3.1 唯一性证明](#31-唯一性证明)
    - [3.2 线程安全性证明](#32-线程安全性证明)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 基础实现](#41-基础实现)
    - [4.2 泛型实现](#42-泛型实现)
    - [4.3 函数式实现](#43-函数式实现)
  - [5. 性能分析](#5-性能分析)
    - [5.1 时间复杂度](#51-时间复杂度)
    - [5.2 空间复杂度](#52-空间复杂度)
    - [5.3 性能对比](#53-性能对比)
  - [6. 应用场景](#6-应用场景)
    - [6.1 配置管理](#61-配置管理)
    - [6.2 日志记录器](#62-日志记录器)
    - [6.3 数据库连接池](#63-数据库连接池)
  - [7. 相关模式](#7-相关模式)
    - [7.1 与工厂模式的关系](#71-与工厂模式的关系)
    - [7.2 与享元模式的关系](#72-与享元模式的关系)
    - [7.3 与依赖注入的关系](#73-与依赖注入的关系)
  - [总结](#总结)

---

## 1. 概念与定义

### 1.1 基本概念

单例模式是一种创建型设计模式，确保一个类只有一个实例，并提供一个全局访问点。

### 1.2 核心特征

- **唯一性**: 类只有一个实例
- **全局访问**: 提供全局访问点
- **延迟初始化**: 实例在首次访问时创建
- **线程安全**: 在多线程环境下保证唯一性

### 1.3 设计原则

```go
// 设计原则：单一职责原则 + 开闭原则
type Singleton interface {
    DoSomething() string
    GetInstance() Singleton
}
```

---

## 2. 形式化定义

### 2.1 集合论定义

设 $S$ 为单例类，$I$ 为实例集合，则：

$$S = \{s \mid s \in I \land |I| = 1\}$$

其中 $|I|$ 表示集合 $I$ 的基数。

### 2.2 状态机定义

单例模式可以表示为有限状态机 $M = (Q, \Sigma, \delta, q_0, F)$：

- $Q = \{\text{未初始化}, \text{已初始化}\}$
- $\Sigma = \{\text{getInstance()}\}$
- $\delta: Q \times \Sigma \rightarrow Q$
- $q_0 = \text{未初始化}$
- $F = \{\text{已初始化}\}$

### 2.3 形式化约束

```go
// 形式化约束定义
type SingletonConstraints struct {
    Uniqueness    bool // ∀x,y ∈ S: x = y
    GlobalAccess  bool // ∃g: ∀s ∈ S: g() = s
    ThreadSafety  bool // ∀t1,t2: getInstance(t1) = getInstance(t2)
}
```

---

## 3. 数学证明

### 3.1 唯一性证明

**定理**: 单例模式保证实例唯一性

**证明**:

1. 假设存在两个实例 $s_1, s_2 \in S$
2. 根据单例约束：$s_1 = s_2$
3. 因此 $|S| = 1$，唯一性得证

### 3.2 线程安全性证明

**定理**: 使用互斥锁的单例模式是线程安全的

**证明**:

1. 设 $M$ 为互斥锁，$s$ 为单例实例
2. 对于任意线程 $t_1, t_2$：
   - $t_1$ 获取锁：$M.Lock()$
   - 检查实例：$\text{if } s == \text{nil}$
   - 创建实例：$s = \text{new}(S)$
   - 释放锁：$M.Unlock()$
3. 由于互斥锁的排他性，$t_2$ 必须等待 $t_1$ 完成
4. 因此线程安全性得证

---

## 4. Go语言实现

### 4.1 基础实现

```go
package singleton

import (
    "fmt"
    "sync"
    "time"
)

// Singleton 单例接口
type Singleton interface {
    DoSomething() string
    GetID() string
}

// singleton 具体单例实现
type singleton struct {
    id        string
    createdAt time.Time
}

var (
    instance *singleton
    once     sync.Once
    mu       sync.Mutex
)

// GetInstance 获取单例实例（线程安全）
func GetInstance() Singleton {
    once.Do(func() {
        instance = &singleton{
            id:        fmt.Sprintf("singleton-%d", time.Now().UnixNano()),
            createdAt: time.Now(),
        }
    })
    return instance
}

// DoSomething 业务方法
func (s *singleton) DoSomething() string {
    return fmt.Sprintf("Singleton[%s] doing work", s.id)
}

// GetID 获取实例ID
func (s *singleton) GetID() string {
    return s.id
}
```

### 4.2 泛型实现

```go
package singleton

import (
    "fmt"
    "sync"
)

// GenericSingleton 泛型单例接口
type GenericSingleton[T any] interface {
    GetValue() T
    SetValue(T)
}

// genericSingleton 泛型单例实现
type genericSingleton[T any] struct {
    value T
    mu    sync.RWMutex
}

var (
    genericInstances = make(map[string]interface{})
    genericMu        sync.RWMutex
)

// GetGenericInstance 获取泛型单例实例
func GetGenericInstance[T any](key string) GenericSingleton[T] {
    genericMu.RLock()
    if instance, exists := genericInstances[key]; exists {
        genericMu.RUnlock()
        return instance.(GenericSingleton[T])
    }
    genericMu.RUnlock()
    
    genericMu.Lock()
    defer genericMu.Unlock()
    
    // 双重检查
    if instance, exists := genericInstances[key]; exists {
        return instance.(GenericSingleton[T])
    }
    
    newInstance := &genericSingleton[T]{}
    genericInstances[key] = newInstance
    return newInstance
}

// GetValue 获取值
func (g *genericSingleton[T]) GetValue() T {
    g.mu.RLock()
    defer g.mu.RUnlock()
    return g.value
}

// SetValue 设置值
func (g *genericSingleton[T]) SetValue(value T) {
    g.mu.Lock()
    defer g.mu.Unlock()
    g.value = value
}
```

### 4.3 函数式实现

```go
package singleton

import (
    "fmt"
    "sync"
)

// FunctionalSingleton 函数式单例
type FunctionalSingleton struct {
    operations []func() string
    mu         sync.RWMutex
}

var (
    functionalInstance *FunctionalSingleton
    functionalOnce     sync.Once
)

// GetFunctionalInstance 获取函数式单例
func GetFunctionalInstance() *FunctionalSingleton {
    functionalOnce.Do(func() {
        functionalInstance = &FunctionalSingleton{
            operations: make([]func() string, 0),
        }
    })
    return functionalInstance
}

// AddOperation 添加操作
func (fs *FunctionalSingleton) AddOperation(op func() string) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    fs.operations = append(fs.operations, op)
}

// ExecuteOperations 执行所有操作
func (fs *FunctionalSingleton) ExecuteOperations() []string {
    fs.mu.RLock()
    defer fs.mu.RUnlock()
    
    results := make([]string, len(fs.operations))
    for i, op := range fs.operations {
        results[i] = op()
    }
    return results
}
```

---

## 5. 性能分析

### 5.1 时间复杂度

| 操作 | 时间复杂度 | 说明 |
|------|------------|------|
| 首次创建 | O(1) | 创建新实例 |
| 后续访问 | O(1) | 直接返回实例 |
| 线程安全检查 | O(1) | 互斥锁操作 |

### 5.2 空间复杂度

- **空间复杂度**: O(1)
- **内存占用**: 固定大小，不随访问次数增长

### 5.3 性能对比

```go
// 性能测试代码
func BenchmarkSingleton(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        GetInstance()
    }
}

func BenchmarkNewInstance(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = &singleton{
            id:        fmt.Sprintf("instance-%d", i),
            createdAt: time.Now(),
        }
    }
}
```

---

## 6. 应用场景

### 6.1 配置管理

```go
// 配置管理器单例
type ConfigManager struct {
    config map[string]interface{}
    mu     sync.RWMutex
}

func (cm *ConfigManager) Get(key string) interface{} {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.config[key]
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.config[key] = value
}
```

### 6.2 日志记录器

```go
// 日志记录器单例
type Logger struct {
    level  string
    output io.Writer
    mu     sync.Mutex
}

func (l *Logger) Log(level, message string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    fmt.Fprintf(l.output, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, message)
}
```

### 6.3 数据库连接池

```go
// 数据库连接池单例
type DBConnectionPool struct {
    connections chan *sql.DB
    maxConn     int
    mu          sync.Mutex
}

func (pool *DBConnectionPool) GetConnection() *sql.DB {
    select {
    case conn := <-pool.connections:
        return conn
    default:
        return pool.createConnection()
    }
}
```

---

## 7. 相关模式

### 7.1 与工厂模式的关系

- **单例模式**: 确保唯一实例
- **工厂模式**: 创建多个实例
- **组合使用**: 工厂方法返回单例实例

### 7.2 与享元模式的关系

- **单例模式**: 全局唯一实例
- **享元模式**: 共享多个相似实例
- **区别**: 单例强调唯一性，享元强调共享性

### 7.3 与依赖注入的关系

```go
// 依赖注入容器中的单例
type Container struct {
    singletons map[string]interface{}
    mu         sync.RWMutex
}

func (c *Container) RegisterSingleton(key string, factory func() interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.singletons[key] = factory()
}

func (c *Container) GetSingleton(key string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.singletons[key]
}
```

---

## 总结

单例模式通过数学形式化定义和Go语言实现，确保了实例的唯一性和全局访问性。通过互斥锁和sync.Once等机制，保证了线程安全性。该模式在配置管理、日志记录、连接池等场景中广泛应用，是软件工程中的重要设计模式。

**相关链接**:

- [02-工厂方法模式](../02-Factory-Method-Pattern.md)
- [03-抽象工厂模式](../03-Abstract-Factory-Pattern.md)
- [返回设计模式目录](../../README.md)
