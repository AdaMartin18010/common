# 01. 单例模式 (Singleton Pattern)

## 目录

- [01. 单例模式 (Singleton Pattern)](#01-单例模式-singleton-pattern)
  - [目录](#目录)
  - [1. 模式概述](#1-模式概述)
    - [1.1 定义与目的](#11-定义与目的)
    - [1.2 应用场景](#12-应用场景)
    - [1.3 优缺点分析](#13-优缺点分析)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 数学定义](#21-数学定义)
    - [2.2 类型系统](#22-类型系统)
    - [2.3 行为规范](#23-行为规范)
  - [3. 实现方式](#3-实现方式)
    - [3.1 饿汉式单例](#31-饿汉式单例)
    - [3.2 懒汉式单例](#32-懒汉式单例)
    - [3.3 双重检查锁定](#33-双重检查锁定)
    - [3.4 静态内部类](#34-静态内部类)
    - [3.5 枚举单例](#35-枚举单例)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 基础实现](#41-基础实现)
    - [4.2 线程安全实现](#42-线程安全实现)
    - [4.3 泛型实现](#43-泛型实现)
    - [4.4 函数式实现](#44-函数式实现)
  - [5. 并发安全](#5-并发安全)
    - [5.1 内存模型](#51-内存模型)
    - [5.2 同步机制](#52-同步机制)
    - [5.3 性能优化](#53-性能优化)
  - [6. 测试与验证](#6-测试与验证)
    - [6.1 单元测试](#61-单元测试)
    - [6.2 并发测试](#62-并发测试)
    - [6.3 性能测试](#63-性能测试)
  - [7. 应用实例](#7-应用实例)
    - [7.1 配置管理器](#71-配置管理器)
    - [7.2 日志记录器](#72-日志记录器)
    - [7.3 数据库连接池](#73-数据库连接池)
  - [8. 定理与证明](#8-定理与证明)
    - [8.1 唯一性定理](#81-唯一性定理)
    - [8.2 线程安全定理](#82-线程安全定理)
    - [8.3 性能定理](#83-性能定理)

---

## 1. 模式概述

### 1.1 定义与目的

**定义 1.1.1** (单例模式)
单例模式是一种创建型设计模式，确保一个类只有一个实例，并提供一个全局访问点。

**形式化定义**：
设 ```latex
$C$
``` 是一个类，单例模式确保：
$```latex
$\forall x, y \in C: x = y$
```$

**目的**：

1. **唯一性保证**：确保系统中某个类只有一个实例
2. **全局访问**：提供全局访问点
3. **延迟初始化**：支持延迟创建实例
4. **线程安全**：在多线程环境下保证实例的唯一性

### 1.2 应用场景

**场景 1.2.1** (配置管理)
系统配置信息需要全局访问，且配置在运行时不应被修改。

**场景 1.2.2** (日志记录)
日志记录器需要全局访问，避免多个日志文件冲突。

**场景 1.2.3** (数据库连接池)
数据库连接池需要全局管理，避免资源浪费。

**场景 1.2.4** (缓存管理)
缓存管理器需要全局访问，确保缓存一致性。

### 1.3 优缺点分析

**优点**：

1. **内存效率**：只创建一个实例，节省内存
2. **全局访问**：提供统一的访问点
3. **延迟初始化**：支持按需创建
4. **线程安全**：保证多线程环境下的正确性

**缺点**：

1. **全局状态**：引入全局状态，增加复杂性
2. **测试困难**：全局状态使单元测试复杂化
3. **违反单一职责**：同时负责创建和管理实例
4. **并发开销**：线程安全实现可能带来性能开销

## 2. 形式化定义

### 2.1 数学定义

**定义 2.1.1** (单例类)
单例类 ```latex
$S$
``` 是一个满足以下条件的类：

1. **唯一性**：```latex
$\forall x, y \in S: x = y$
```
2. **存在性**：```latex
$\exists x \in S$
```
3. **可访问性**：```latex
$\forall x \in S: \text{accessible}(x)$
```

**定义 2.1.2** (单例函数)
单例函数 ```latex
$f: \emptyset \to S$
``` 满足：
$```latex
$f() = \text{the unique instance of } S$
```$

**定义 2.1.3** (线程安全单例)
线程安全单例 ```latex
$S$
``` 满足：
$```latex
$\forall t_1, t_2 \in \text{Threads}: f_{t_1}() = f_{t_2}()$
```$

### 2.2 类型系统

**定义 2.2.1** (单例类型)
在类型系统中，单例类型 ```latex
$T$
``` 满足：
$```latex
$\text{Card}(T) = 1$
```$

其中 ```latex
$\text{Card}(T)$
``` 表示类型 ```latex
$T$
``` 的基数。

**定义 2.2.2** (单例接口)
单例接口 ```latex
$I$
``` 定义：

```typescript
interface Singleton<T> {
    getInstance(): T;
    reset(): void;
}
```

**定义 2.2.3** (单例约束)
单例约束 ```latex
$C$
``` 确保：
$```latex
$\forall x, y: \text{instanceOf}(x, T) \land \text{instanceOf}(y, T) \Rightarrow x = y$
```$

### 2.3 行为规范

**规范 2.3.1** (创建行为)

1. 第一次调用时创建实例
2. 后续调用返回相同实例
3. 创建过程是原子的

**规范 2.3.2** (访问行为)

1. 提供统一的访问方法
2. 访问方法是线程安全的
3. 访问方法返回相同的实例

**规范 2.3.3** (生命周期行为)

1. 实例在程序运行期间存在
2. 支持显式重置（可选）
3. 支持优雅关闭

## 3. 实现方式

### 3.1 饿汉式单例

**定义 3.1.1** (饿汉式单例)
在类加载时就完成初始化的单例模式。

**特点**：

- 线程安全（类加载时初始化）
- 不支持延迟初始化
- 可能造成不必要的内存占用

**Go实现**：

```go
package singleton

import (
    "sync"
    "time"
)

// EagerSingleton 饿汉式单例
type EagerSingleton struct {
    id        string
    createdAt time.Time
}

// 在包级别初始化实例
var eagerInstance = &EagerSingleton{
    id:        "eager-singleton",
    createdAt: time.Now(),
}

// GetEagerInstance 获取饿汉式单例实例
func GetEagerInstance() *EagerSingleton {
    return eagerInstance
}

// GetID 获取实例ID
func (s *EagerSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *EagerSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

### 3.2 懒汉式单例

**定义 3.2.1** (懒汉式单例)
在第一次使用时才创建实例的单例模式。

**特点**：

- 支持延迟初始化
- 需要处理线程安全问题
- 可能影响性能

**Go实现**：

```go
// LazySingleton 懒汉式单例
type LazySingleton struct {
    id        string
    createdAt time.Time
}

var (
    lazyInstance *LazySingleton
    lazyMutex    sync.Mutex
)

// GetLazyInstance 获取懒汉式单例实例
func GetLazyInstance() *LazySingleton {
    lazyMutex.Lock()
    defer lazyMutex.Unlock()
    
    if lazyInstance == nil {
        lazyInstance = &LazySingleton{
            id:        "lazy-singleton",
            createdAt: time.Now(),
        }
    }
    
    return lazyInstance
}

// GetID 获取实例ID
func (s *LazySingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *LazySingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

### 3.3 双重检查锁定

**定义 3.3.1** (双重检查锁定)
使用双重检查来减少锁的开销的线程安全单例模式。

**特点**：

- 线程安全
- 性能优化
- 支持延迟初始化

**Go实现**：

```go
// DoubleCheckSingleton 双重检查锁定单例
type DoubleCheckSingleton struct {
    id        string
    createdAt time.Time
}

var (
    doubleCheckInstance *DoubleCheckSingleton
    doubleCheckMutex    sync.RWMutex
)

// GetDoubleCheckInstance 获取双重检查锁定单例实例
func GetDoubleCheckInstance() *DoubleCheckSingleton {
    // 第一次检查（无锁）
    if doubleCheckInstance == nil {
        // 获取写锁
        doubleCheckMutex.Lock()
        defer doubleCheckMutex.Unlock()
        
        // 第二次检查（有锁）
        if doubleCheckInstance == nil {
            doubleCheckInstance = &DoubleCheckSingleton{
                id:        "double-check-singleton",
                createdAt: time.Now(),
            }
        }
    }
    
    return doubleCheckInstance
}

// GetID 获取实例ID
func (s *DoubleCheckSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *DoubleCheckSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

### 3.4 静态内部类

**定义 3.4.1** (静态内部类单例)
使用静态内部类实现延迟初始化的单例模式。

**特点**：

- 线程安全（类加载时初始化）
- 支持延迟初始化
- 实现简单

**Go实现**（使用包级变量模拟）：

```go
// StaticInnerSingleton 静态内部类单例
type StaticInnerSingleton struct {
    id        string
    createdAt time.Time
}

// 使用sync.Once确保线程安全
var (
    staticInnerInstance *StaticInnerSingleton
    staticInnerOnce     sync.Once
)

// GetStaticInnerInstance 获取静态内部类单例实例
func GetStaticInnerInstance() *StaticInnerSingleton {
    staticInnerOnce.Do(func() {
        staticInnerInstance = &StaticInnerSingleton{
            id:        "static-inner-singleton",
            createdAt: time.Now(),
        }
    })
    
    return staticInnerInstance
}

// GetID 获取实例ID
func (s *StaticInnerSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *StaticInnerSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

### 3.5 枚举单例

**定义 3.5.1** (枚举单例)
使用枚举实现单例模式（Go中通过常量模拟）。

**特点**：

- 线程安全
- 自动序列化支持
- 防止反射攻击

**Go实现**：

```go
// EnumSingleton 枚举单例
type EnumSingleton struct {
    id        string
    createdAt time.Time
}

// 使用常量定义单例实例
const (
    EnumInstanceID = "enum-singleton"
)

var enumInstance = &EnumSingleton{
    id:        EnumInstanceID,
    createdAt: time.Now(),
}

// GetEnumInstance 获取枚举单例实例
func GetEnumInstance() *EnumSingleton {
    return enumInstance
}

// GetID 获取实例ID
func (s *EnumSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *EnumSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

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
    GetID() string
    GetCreatedAt() time.Time
    DoSomething() string
}

// BaseSingleton 基础单例实现
type BaseSingleton struct {
    id        string
    createdAt time.Time
    counter   int
    mutex     sync.RWMutex
}

// 全局实例
var (
    instance *BaseSingleton
    once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() Singleton {
    once.Do(func() {
        instance = &BaseSingleton{
            id:        "base-singleton",
            createdAt: time.Now(),
            counter:   0,
        }
    })
    return instance
}

// GetID 获取实例ID
func (s *BaseSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *BaseSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}

// DoSomething 执行操作
func (s *BaseSingleton) DoSomething() string {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.counter++
    return fmt.Sprintf("Operation %d executed at %v", s.counter, time.Now())
}

// GetCounter 获取计数器值
func (s *BaseSingleton) GetCounter() int {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    return s.counter
}

// Reset 重置单例（用于测试）
func Reset() {
    instance = nil
    once = sync.Once{}
}
```

### 4.2 线程安全实现

```go
// ThreadSafeSingleton 线程安全单例
type ThreadSafeSingleton struct {
    id        string
    createdAt time.Time
    data      map[string]interface{}
    mutex     sync.RWMutex
}

var (
    threadSafeInstance *ThreadSafeSingleton
    threadSafeOnce     sync.Once
)

// GetThreadSafeInstance 获取线程安全单例实例
func GetThreadSafeInstance() *ThreadSafeSingleton {
    threadSafeOnce.Do(func() {
        threadSafeInstance = &ThreadSafeSingleton{
            id:        "thread-safe-singleton",
            createdAt: time.Now(),
            data:      make(map[string]interface{}),
        }
    })
    return threadSafeInstance
}

// SetData 设置数据
func (s *ThreadSafeSingleton) SetData(key string, value interface{}) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.data[key] = value
}

// GetData 获取数据
func (s *ThreadSafeSingleton) GetData(key string) (interface{}, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    value, exists := s.data[key]
    return value, exists
}

// GetAllData 获取所有数据
func (s *ThreadSafeSingleton) GetAllData() map[string]interface{} {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    result := make(map[string]interface{})
    for k, v := range s.data {
        result[k] = v
    }
    return result
}

// ClearData 清空数据
func (s *ThreadSafeSingleton) ClearData() {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.data = make(map[string]interface{})
}

// GetID 获取实例ID
func (s *ThreadSafeSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *ThreadSafeSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}
```

### 4.3 泛型实现

```go
// GenericSingleton 泛型单例
type GenericSingleton[T any] struct {
    instance T
    once     sync.Once
    factory  func() T
}

// NewGenericSingleton 创建泛型单例
func NewGenericSingleton[T any](factory func() T) *GenericSingleton[T] {
    return &GenericSingleton[T]{
        factory: factory,
    }
}

// GetInstance 获取泛型单例实例
func (s *GenericSingleton[T]) GetInstance() T {
    s.once.Do(func() {
        s.instance = s.factory()
    })
    return s.instance
}

// Reset 重置泛型单例（用于测试）
func (s *GenericSingleton[T]) Reset() {
    s.once = sync.Once{}
}

// 使用示例
type Config struct {
    DatabaseURL string
    Port        int
    Debug       bool
}

// NewConfig 创建配置实例
func NewConfig() Config {
    return Config{
        DatabaseURL: "localhost:5432",
        Port:        8080,
        Debug:       true,
    }
}

// 全局配置单例
var ConfigSingleton = NewGenericSingleton(NewConfig)

// GetConfig 获取配置实例
func GetConfig() Config {
    return ConfigSingleton.GetInstance()
}
```

### 4.4 函数式实现

```go
// FunctionalSingleton 函数式单例
type FunctionalSingleton struct {
    id        string
    createdAt time.Time
    operations []func() string
}

// SingletonFactory 单例工厂函数
type SingletonFactory func() *FunctionalSingleton

// CreateSingleton 创建单例的工厂函数
func CreateSingleton(id string) SingletonFactory {
    var instance *FunctionalSingleton
    var once sync.Once
    
    return func() *FunctionalSingleton {
        once.Do(func() {
            instance = &FunctionalSingleton{
                id:        id,
                createdAt: time.Now(),
                operations: make([]func() string, 0),
            }
        })
        return instance
    }
}

// AddOperation 添加操作
func (s *FunctionalSingleton) AddOperation(operation func() string) {
    s.operations = append(s.operations, operation)
}

// ExecuteOperations 执行所有操作
func (s *FunctionalSingleton) ExecuteOperations() []string {
    results := make([]string, len(s.operations))
    for i, operation := range s.operations {
        results[i] = operation()
    }
    return results
}

// GetID 获取实例ID
func (s *FunctionalSingleton) GetID() string {
    return s.id
}

// GetCreatedAt 获取创建时间
func (s *FunctionalSingleton) GetCreatedAt() time.Time {
    return s.createdAt
}

// 使用示例
var (
    LoggerFactory = CreateSingleton("logger")
    CacheFactory  = CreateSingleton("cache")
)

// GetLogger 获取日志单例
func GetLogger() *FunctionalSingleton {
    return LoggerFactory()
}

// GetCache 获取缓存单例
func GetCache() *FunctionalSingleton {
    return CacheFactory()
}
```

## 5. 并发安全

### 5.1 内存模型

**定义 5.1.1** (内存可见性)
在多线程环境中，单例实例的创建和访问需要保证内存可见性。

**Go内存模型**：

```go
// 使用sync.Once保证内存可见性
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

**定义 5.1.2** (原子性)
单例实例的创建过程必须是原子的，不能被中断。

**实现方式**：

1. 使用 `sync.Once`
2. 使用互斥锁
3. 使用原子操作

### 5.2 同步机制

**机制 5.2.1** (互斥锁)
使用互斥锁保证线程安全：

```go
type MutexSingleton struct {
    instance *Singleton
    mutex    sync.Mutex
}

func (s *MutexSingleton) GetInstance() *Singleton {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.instance == nil {
        s.instance = &Singleton{}
    }
    return s.instance
}
```

**机制 5.2.2** (读写锁)
使用读写锁优化性能：

```go
type RWMutexSingleton struct {
    instance *Singleton
    mutex    sync.RWMutex
}

func (s *RWMutexSingleton) GetInstance() *Singleton {
    // 读锁检查
    s.mutex.RLock()
    if s.instance != nil {
        s.mutex.RUnlock()
        return s.instance
    }
    s.mutex.RUnlock()
    
    // 写锁创建
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.instance == nil {
        s.instance = &Singleton{}
    }
    return s.instance
}
```

**机制 5.2.3** (原子操作)
使用原子操作保证线程安全：

```go
type AtomicSingleton struct {
    instance atomic.Value
}

func (s *AtomicSingleton) GetInstance() *Singleton {
    if instance := s.instance.Load(); instance != nil {
        return instance.(*Singleton)
    }
    
    newInstance := &Singleton{}
    if s.instance.CompareAndSwap(nil, newInstance) {
        return newInstance
    }
    
    return s.instance.Load().(*Singleton)
}
```

### 5.3 性能优化

**优化 5.3.1** (延迟初始化)
只在需要时创建实例：

```go
type LazySingleton struct {
    instance *Singleton
    once     sync.Once
}

func (s *LazySingleton) GetInstance() *Singleton {
    s.once.Do(func() {
        s.instance = &Singleton{}
    })
    return s.instance
}
```

**优化 5.3.2** (缓存优化)
使用缓存减少锁竞争：

```go
type CachedSingleton struct {
    cache map[string]*Singleton
    mutex sync.RWMutex
}

func (s *CachedSingleton) GetInstance(key string) *Singleton {
    // 先检查缓存
    s.mutex.RLock()
    if instance, exists := s.cache[key]; exists {
        s.mutex.RUnlock()
        return instance
    }
    s.mutex.RUnlock()
    
    // 创建新实例
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if instance, exists := s.cache[key]; exists {
        return instance
    }
    
    newInstance := &Singleton{ID: key}
    s.cache[key] = newInstance
    return newInstance
}
```

## 6. 测试与验证

### 6.1 单元测试

```go
package singleton

import (
    "testing"
    "time"
)

func TestGetInstance(t *testing.T) {
    // 重置单例
    Reset()
    
    // 获取第一个实例
    instance1 := GetInstance()
    if instance1 == nil {
        t.Error("Expected non-nil instance")
    }
    
    // 获取第二个实例
    instance2 := GetInstance()
    if instance2 == nil {
        t.Error("Expected non-nil instance")
    }
    
    // 验证是同一个实例
    if instance1 != instance2 {
        t.Error("Expected same instance")
    }
    
    // 验证ID相同
    if instance1.GetID() != instance2.GetID() {
        t.Error("Expected same ID")
    }
    
    // 验证创建时间相同
    if !instance1.GetCreatedAt().Equal(instance2.GetCreatedAt()) {
        t.Error("Expected same creation time")
    }
}

func TestSingletonOperations(t *testing.T) {
    Reset()
    
    instance := GetInstance()
    
    // 测试操作
    result1 := instance.DoSomething()
    result2 := instance.DoSomething()
    
    if result1 == result2 {
        t.Error("Expected different operation results")
    }
    
    // 验证计数器
    counter := instance.(*BaseSingleton).GetCounter()
    if counter != 2 {
        t.Errorf("Expected counter to be 2, got %d", counter)
    }
}
```

### 6.2 并发测试

```go
func TestConcurrentAccess(t *testing.T) {
    Reset()
    
    const numGoroutines = 100
    const numOperations = 1000
    
    // 启动多个goroutine并发访问
    done := make(chan bool, numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            instance := GetInstance()
            
            // 执行多次操作
            for j := 0; j < numOperations; j++ {
                instance.DoSomething()
            }
            
            done <- true
        }(i)
    }
    
    // 等待所有goroutine完成
    for i := 0; i < numGoroutines; i++ {
        <-done
    }
    
    // 验证最终计数器值
    instance := GetInstance()
    expectedCounter := numGoroutines * numOperations
    actualCounter := instance.(*BaseSingleton).GetCounter()
    
    if actualCounter != expectedCounter {
        t.Errorf("Expected counter to be %d, got %d", expectedCounter, actualCounter)
    }
}

func TestRaceCondition(t *testing.T) {
    Reset()
    
    const numGoroutines = 1000
    
    // 使用竞态检测器
    instances := make([]Singleton, numGoroutines)
    
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer wg.Done()
            instances[id] = GetInstance()
        }(i)
    }
    
    wg.Wait()
    
    // 验证所有实例都是同一个
    firstInstance := instances[0]
    for i := 1; i < numGoroutines; i++ {
        if instances[i] != firstInstance {
            t.Errorf("Expected same instance at index %d", i)
        }
    }
}
```

### 6.3 性能测试

```go
func BenchmarkGetInstance(b *testing.B) {
    Reset()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        GetInstance()
    }
}

func BenchmarkConcurrentGetInstance(b *testing.B) {
    Reset()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            GetInstance()
        }
    })
}

func BenchmarkSingletonOperations(b *testing.B) {
    Reset()
    instance := GetInstance()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        instance.DoSomething()
    }
}
```

## 7. 应用实例

### 7.1 配置管理器

```go
package config

import (
    "encoding/json"
    "os"
    "sync"
    "time"
)

// Config 配置结构
type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Logging  LoggingConfig  `json:"logging"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type ServerConfig struct {
    Port    int    `json:"port"`
    Host    string `json:"host"`
    Timeout int    `json:"timeout"`
}

type LoggingConfig struct {
    Level   string `json:"level"`
    File    string `json:"file"`
    Console bool   `json:"console"`
}

// ConfigManager 配置管理器单例
type ConfigManager struct {
    config     *Config
    configFile string
    lastLoad   time.Time
    mutex      sync.RWMutex
}

var (
    configInstance *ConfigManager
    configOnce     sync.Once
)

// GetConfigManager 获取配置管理器实例
func GetConfigManager() *ConfigManager {
    configOnce.Do(func() {
        configInstance = &ConfigManager{
            configFile: "config.json",
        }
    })
    return configInstance
}

// LoadConfig 加载配置
func (cm *ConfigManager) LoadConfig() error {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    
    file, err := os.Open(cm.configFile)
    if err != nil {
        return err
    }
    defer file.Close()
    
    var config Config
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return err
    }
    
    cm.config = &config
    cm.lastLoad = time.Now()
    
    return nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *Config {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    return cm.config
}

// GetDatabaseConfig 获取数据库配置
func (cm *ConfigManager) GetDatabaseConfig() DatabaseConfig {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    if cm.config != nil {
        return cm.config.Database
    }
    return DatabaseConfig{}
}

// GetServerConfig 获取服务器配置
func (cm *ConfigManager) GetServerConfig() ServerConfig {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    if cm.config != nil {
        return cm.config.Server
    }
    return ServerConfig{}
}

// GetLoggingConfig 获取日志配置
func (cm *ConfigManager) GetLoggingConfig() LoggingConfig {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    if cm.config != nil {
        return cm.config.Logging
    }
    return LoggingConfig{}
}

// ReloadConfig 重新加载配置
func (cm *ConfigManager) ReloadConfig() error {
    return cm.LoadConfig()
}

// GetLastLoadTime 获取最后加载时间
func (cm *ConfigManager) GetLastLoadTime() time.Time {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    return cm.lastLoad
}
```

### 7.2 日志记录器

```go
package logger

import (
    "fmt"
    "log"
    "os"
    "sync"
    "time"
)

// LogLevel 日志级别
type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

// Logger 日志记录器单例
type Logger struct {
    level    LogLevel
    file     *os.File
    logger   *log.Logger
    mutex    sync.RWMutex
    rotation bool
}

var (
    loggerInstance *Logger
    loggerOnce     sync.Once
)

// GetLogger 获取日志记录器实例
func GetLogger() *Logger {
    loggerOnce.Do(func() {
        loggerInstance = &Logger{
            level: INFO,
        }
        loggerInstance.init()
    })
    return loggerInstance
}

// init 初始化日志记录器
func (l *Logger) init() {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    // 创建日志文件
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(fmt.Sprintf("Failed to open log file: %v", err))
    }
    
    l.file = file
    l.logger = log.New(file, "", log.LstdFlags)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    l.level = level
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
    l.log(DEBUG, "DEBUG", format, args...)
}

// Info 记录信息日志
func (l *Logger) Info(format string, args ...interface{}) {
    l.log(INFO, "INFO", format, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
    l.log(WARN, "WARN", format, args...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, args ...interface{}) {
    l.log(ERROR, "ERROR", format, args...)
}

// Fatal 记录致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
    l.log(FATAL, "FATAL", format, args...)
    os.Exit(1)
}

// log 内部日志记录方法
func (l *Logger) log(level LogLevel, levelStr, format string, args ...interface{}) {
    if level < l.level {
        return
    }
    
    l.mutex.RLock()
    defer l.mutex.RUnlock()
    
    message := fmt.Sprintf(format, args...)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logMessage := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, message)
    
    l.logger.Println(logMessage)
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    if l.file != nil {
        return l.file.Close()
    }
    return nil
}

// Rotate 轮转日志文件
func (l *Logger) Rotate() error {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    if l.file != nil {
        l.file.Close()
    }
    
    // 重命名当前日志文件
    timestamp := time.Now().Format("20060102_150405")
    newName := fmt.Sprintf("app_%s.log", timestamp)
    os.Rename("app.log", newName)
    
    // 创建新的日志文件
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    
    l.file = file
    l.logger = log.New(file, "", log.LstdFlags)
    
    return nil
}
```

### 7.3 数据库连接池

```go
package database

import (
    "database/sql"
    "fmt"
    "sync"
    "time"
    
    _ "github.com/lib/pq"
)

// ConnectionPool 数据库连接池单例
type ConnectionPool struct {
    db       *sql.DB
    config   DBConfig
    mutex    sync.RWMutex
    metrics  PoolMetrics
}

type DBConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    MaxOpen  int
    MaxIdle  int
    Timeout  time.Duration
}

type PoolMetrics struct {
    TotalConnections int
    ActiveConnections int
    IdleConnections  int
    WaitCount        int64
    WaitDuration     time.Duration
    MaxIdleClosed    int64
    MaxLifetimeClosed int64
}

var (
    poolInstance *ConnectionPool
    poolOnce     sync.Once
)

// GetConnectionPool 获取连接池实例
func GetConnectionPool() *ConnectionPool {
    poolOnce.Do(func() {
        poolInstance = &ConnectionPool{
            config: DBConfig{
                Host:     "localhost",
                Port:     5432,
                Username: "postgres",
                Password: "password",
                Database: "testdb",
                MaxOpen:  10,
                MaxIdle:  5,
                Timeout:  30 * time.Second,
            },
        }
        poolInstance.init()
    })
    return poolInstance
}

// init 初始化连接池
func (p *ConnectionPool) init() {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        p.config.Host, p.config.Port, p.config.Username, p.config.Password, p.config.Database)
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(fmt.Sprintf("Failed to open database: %v", err))
    }
    
    // 配置连接池
    db.SetMaxOpenConns(p.config.MaxOpen)
    db.SetMaxIdleConns(p.config.MaxIdle)
    db.SetConnMaxLifetime(p.config.Timeout)
    
    // 测试连接
    if err := db.Ping(); err != nil {
        panic(fmt.Sprintf("Failed to ping database: %v", err))
    }
    
    p.db = db
}

// GetDB 获取数据库连接
func (p *ConnectionPool) GetDB() *sql.DB {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    return p.db
}

// ExecuteQuery 执行查询
func (p *ConnectionPool) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    start := time.Now()
    rows, err := p.db.Query(query, args...)
    p.updateMetrics(start)
    
    return rows, err
}

// ExecuteCommand 执行命令
func (p *ConnectionPool) ExecuteCommand(query string, args ...interface{}) (sql.Result, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    start := time.Now()
    result, err := p.db.Exec(query, args...)
    p.updateMetrics(start)
    
    return result, err
}

// BeginTransaction 开始事务
func (p *ConnectionPool) BeginTransaction() (*sql.Tx, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    return p.db.Begin()
}

// Close 关闭连接池
func (p *ConnectionPool) Close() error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.db != nil {
        return p.db.Close()
    }
    return nil
}

// GetMetrics 获取连接池指标
func (p *ConnectionPool) GetMetrics() PoolMetrics {
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
    if p.db != nil {
        p.metrics.TotalConnections = p.db.Stats().MaxOpenConnections
        p.metrics.ActiveConnections = p.db.Stats().OpenConnections
        p.metrics.IdleConnections = p.db.Stats().Idle
        p.metrics.WaitCount = p.db.Stats().WaitCount
        p.metrics.WaitDuration = p.db.Stats().WaitDuration
        p.metrics.MaxIdleClosed = p.db.Stats().MaxIdleClosed
        p.metrics.MaxLifetimeClosed = p.db.Stats().MaxLifetimeClosed
    }
    
    return p.metrics
}

// updateMetrics 更新指标
func (p *ConnectionPool) updateMetrics(start time.Time) {
    duration := time.Since(start)
    if duration > p.metrics.WaitDuration {
        p.metrics.WaitDuration = duration
    }
}
```

## 8. 定理与证明

### 8.1 唯一性定理

**定理 8.1.1** (单例唯一性)
使用 `sync.Once` 实现的单例模式保证实例的唯一性。

**证明**：

1. `sync.Once` 保证 `Do` 方法只执行一次
2. 实例创建在 `Do` 方法内部
3. 因此实例只被创建一次
4. 所有调用都返回同一个实例

**形式化证明**：
设 ```latex
$f$
``` 是获取实例的函数，```latex
$o$
``` 是 `sync.Once` 实例，```latex
$c$
``` 是创建实例的函数。

$```latex
$\forall x, y: f() = o.Do(c) \land f() = o.Do(c) \Rightarrow x = y$
```$

由于 `sync.Once` 保证 ```latex
$o.Do(c)$
``` 只执行一次，所以 ```latex
$x = y$
```。

```latex
$\square$
```

### 8.2 线程安全定理

**定理 8.2.1** (线程安全性)
使用 `sync.Once` 实现的单例模式是线程安全的。

**证明**：

1. `sync.Once` 内部使用原子操作和内存屏障
2. 保证在多个goroutine中只有一个执行 `Do` 方法
3. 其他goroutine会等待第一个goroutine完成
4. 因此所有goroutine都获得相同的实例

**形式化证明**：
设 ```latex
$T$
``` 是线程集合，```latex
$f_t$
``` 是线程 ```latex
$t$
``` 的获取实例函数。

$```latex
$\forall t_1, t_2 \in T: f_{t_1}() = f_{t_2}()$
```$

由于 `sync.Once` 的线程安全保证，所有线程都获得相同的实例。

```latex
$\square$
```

### 8.3 性能定理

**定理 8.3.1** (性能最优性)
使用 `sync.Once` 的单例模式在Go中具有最优性能。

**证明**：

1. 第一次调用需要创建实例，开销为 ```latex
$O(1)$
```
2. 后续调用只需要内存访问，开销为 ```latex
$O(1)$
```
3. 没有锁竞争，避免了互斥锁的开销
4. 内存屏障开销最小

**复杂度分析**：

- 时间复杂度：```latex
$O(1)$
```
- 空间复杂度：```latex
$O(1)$
```
- 并发开销：```latex
$O(1)$
```

```latex
$\square$
```

---

**总结**：
单例模式是软件工程中重要的设计模式，通过严格的数学定义和Go语言的实现，我们可以确保实例的唯一性和线程安全性。在实际应用中，需要根据具体场景选择合适的实现方式，并注意性能优化和测试验证。
