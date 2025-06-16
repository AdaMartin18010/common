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
    - [2.2 函数式定义](#22-函数式定义)
    - [2.3 状态机定义](#23-状态机定义)
  - [3. 数学证明](#3-数学证明)
    - [3.1 唯一性证明](#31-唯一性证明)
    - [3.2 线程安全性证明](#32-线程安全性证明)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 基础实现](#41-基础实现)
    - [4.2 泛型实现](#42-泛型实现)
    - [4.3 函数式实现](#43-函数式实现)
    - [4.4 测试代码](#44-测试代码)
  - [5. 性能分析](#5-性能分析)
    - [5.1 时间复杂度](#51-时间复杂度)
    - [5.2 空间复杂度](#52-空间复杂度)
    - [5.3 性能优化](#53-性能优化)
  - [6. 应用场景](#6-应用场景)
    - [6.1 配置管理](#61-配置管理)
    - [6.2 日志记录器](#62-日志记录器)
    - [6.3 数据库连接池](#63-数据库连接池)
  - [7. 相关模式](#7-相关模式)
    - [7.1 与工厂模式的关系](#71-与工厂模式的关系)
    - [7.2 与享元模式的关系](#72-与享元模式的关系)
    - [7.3 与注册表模式的关系](#73-与注册表模式的关系)
  - [总结](#总结)

---

## 1. 概念与定义

### 1.1 基本概念

单例模式是一种创建型设计模式，确保一个类只有一个实例，并提供一个全局访问点来访问该实例。

### 1.2 核心特征

- **唯一性**: 确保系统中只有一个实例存在
- **全局访问**: 提供全局访问点
- **延迟初始化**: 实例在第一次使用时才创建
- **线程安全**: 在多线程环境下保证唯一性

### 1.3 设计原则

- **单一职责原则**: 类只负责创建和管理自己的唯一实例
- **开闭原则**: 对扩展开放，对修改封闭
- **依赖倒置原则**: 依赖于抽象而不是具体实现

---

## 2. 形式化定义

### 2.1 集合论定义

设 $S$ 为单例类，$I$ 为实例集合，则单例模式满足：

$$\forall s_1, s_2 \in S : s_1 = s_2$$

其中 $s_1, s_2$ 为 $S$ 的任意两个实例。

### 2.2 函数式定义

定义单例函数 $f: \emptyset \rightarrow S$，满足：

$$
f() = \begin{cases}
\text{existing instance} & \text{if exists} \\
\text{new instance} & \text{otherwise}
\end{cases}
$$

### 2.3 状态机定义

单例模式可以表示为状态机 $M = (Q, \Sigma, \delta, q_0, F)$：

- $Q = \{\text{Uninitialized}, \text{Initialized}\}$
- $\Sigma = \{\text{getInstance}\}$
- $\delta: Q \times \Sigma \rightarrow Q$
- $q_0 = \text{Uninitialized}$
- $F = \{\text{Initialized}\}$

---

## 3. 数学证明

### 3.1 唯一性证明

**定理**: 单例模式保证实例的唯一性。

**证明**:

1. 假设存在两个不同的实例 $s_1$ 和 $s_2$
2. 根据单例模式的实现，所有 `getInstance()` 调用都返回同一个实例
3. 因此 $s_1 = s_2$，与假设矛盾
4. 故单例模式保证唯一性

### 3.2 线程安全性证明

**定理**: 使用互斥锁的单例模式是线程安全的。

**证明**:

1. 设 $L$ 为互斥锁，$C$ 为临界区
2. 对于任意两个线程 $T_1, T_2$，访问 $C$ 时：
   - $T_1$ 获得锁 $L$，进入 $C$
   - $T_2$ 等待锁 $L$ 释放
   - $T_1$ 完成操作，释放锁 $L$
   - $T_2$ 获得锁 $L$，进入 $C$
3. 因此保证了串行访问，线程安全

---

## 4. Go语言实现

### 4.1 基础实现

```go
package singleton

import (
    "fmt"
    "sync"
)

// Singleton 单例结构体
type Singleton struct {
    name string
}

var (
    instance *Singleton
    once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            name: "Default Singleton",
        }
        fmt.Println("Creating singleton instance")
    })
    return instance
}

// GetName 获取实例名称
func (s *Singleton) GetName() string {
    return s.name
}

// SetName 设置实例名称
func (s *Singleton) SetName(name string) {
    s.name = name
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
    GetInstance() T
}

// SingletonManager 单例管理器
type SingletonManager[T any] struct {
    instance T
    once     sync.Once
    factory  func() T
}

// NewSingletonManager 创建单例管理器
func NewSingletonManager[T any](factory func() T) *SingletonManager[T] {
    return &SingletonManager[T]{
        factory: factory,
    }
}

// GetInstance 获取泛型单例实例
func (sm *SingletonManager[T]) GetInstance() T {
    sm.once.Do(func() {
        sm.instance = sm.factory()
        fmt.Printf("Creating singleton instance of type %T\n", sm.instance)
    })
    return sm.instance
}

// 使用示例
type Config struct {
    DatabaseURL string
    Port        int
}

func NewConfig() Config {
    return Config{
        DatabaseURL: "localhost:5432",
        Port:        8080,
    }
}

// 全局单例管理器
var configManager = NewSingletonManager(NewConfig)

// GetConfig 获取配置单例
func GetConfig() Config {
    return configManager.GetInstance()
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
    operations []func()
    mutex      sync.RWMutex
}

var (
    functionalInstance *FunctionalSingleton
    functionalOnce     sync.Once
)

// GetFunctionalInstance 获取函数式单例实例
func GetFunctionalInstance() *FunctionalSingleton {
    functionalOnce.Do(func() {
        functionalInstance = &FunctionalSingleton{
            operations: make([]func(), 0),
        }
        fmt.Println("Creating functional singleton instance")
    })
    return functionalInstance
}

// AddOperation 添加操作
func (fs *FunctionalSingleton) AddOperation(operation func()) {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()
    fs.operations = append(fs.operations, operation)
}

// ExecuteOperations 执行所有操作
func (fs *FunctionalSingleton) ExecuteOperations() {
    fs.mutex.RLock()
    defer fs.mutex.RUnlock()

    for i, operation := range fs.operations {
        fmt.Printf("Executing operation %d\n", i+1)
        operation()
    }
}

// ClearOperations 清空操作列表
func (fs *FunctionalSingleton) ClearOperations() {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()
    fs.operations = make([]func(), 0)
}
```

### 4.4 测试代码

```go
package singleton

import (
    "fmt"
    "sync"
    "testing"
    "time"
)

// TestSingletonUniqueness 测试单例唯一性
func TestSingletonUniqueness(t *testing.T) {
    instance1 := GetInstance()
    instance2 := GetInstance()

    if instance1 != instance2 {
        t.Errorf("Singleton instances are not the same")
    }
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
    const numGoroutines = 100
    var wg sync.WaitGroup
    instances := make([]*Singleton, numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            instances[index] = GetInstance()
        }(i)
    }

    wg.Wait()

    // 验证所有实例都是同一个
    firstInstance := instances[0]
    for i := 1; i < numGoroutines; i++ {
        if instances[i] != firstInstance {
            t.Errorf("Instance %d is not the same as first instance", i)
        }
    }
}

// TestGenericSingleton 测试泛型单例
func TestGenericSingleton(t *testing.T) {
    config1 := GetConfig()
    config2 := GetConfig()

    if config1.DatabaseURL != config2.DatabaseURL {
        t.Errorf("Config instances are not the same")
    }
}

// BenchmarkSingleton 性能基准测试
func BenchmarkSingleton(b *testing.B) {
    for i := 0; i < b.N; i++ {
        GetInstance()
    }
}
```

---

## 5. 性能分析

### 5.1 时间复杂度

- **获取实例**: $O(1)$
- **初始化**: $O(1)$
- **内存访问**: $O(1)$

### 5.2 空间复杂度

- **内存占用**: $O(1)$
- **额外开销**: 互斥锁开销

### 5.3 性能优化

```go
// 双重检查锁定模式
type OptimizedSingleton struct {
    name string
}

var (
    optimizedInstance *OptimizedSingleton
    optimizedMutex    sync.Mutex
)

func GetOptimizedInstance() *OptimizedSingleton {
    if optimizedInstance == nil {
        optimizedMutex.Lock()
        defer optimizedMutex.Unlock()

        if optimizedInstance == nil {
            optimizedInstance = &OptimizedSingleton{
                name: "Optimized Singleton",
            }
        }
    }
    return optimizedInstance
}
```

---

## 6. 应用场景

### 6.1 配置管理

```go
// 全局配置单例
type GlobalConfig struct {
    DatabaseURL string
    RedisURL    string
    LogLevel    string
    MaxConnections int
}

var configInstance *GlobalConfig
var configOnce sync.Once

func GetGlobalConfig() *GlobalConfig {
    configOnce.Do(func() {
        configInstance = &GlobalConfig{
            DatabaseURL:    "localhost:5432",
            RedisURL:       "localhost:6379",
            LogLevel:       "INFO",
            MaxConnections: 100,
        }
    })
    return configInstance
}
```

### 6.2 日志记录器

```go
// 日志记录器单例
type Logger struct {
    level string
    file  *os.File
}

var loggerInstance *Logger
var loggerOnce sync.Once

func GetLogger() *Logger {
    loggerOnce.Do(func() {
        file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        loggerInstance = &Logger{
            level: "INFO",
            file:  file,
        }
    })
    return loggerInstance
}

func (l *Logger) Log(message string) {
    fmt.Fprintf(l.file, "[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), l.level, message)
}
```

### 6.3 数据库连接池

```go
// 数据库连接池单例
type DatabasePool struct {
    connections chan *sql.DB
    maxConnections int
}

var poolInstance *DatabasePool
var poolOnce sync.Once

func GetDatabasePool() *DatabasePool {
    poolOnce.Do(func() {
        poolInstance = &DatabasePool{
            connections:    make(chan *sql.DB, 10),
            maxConnections: 10,
        }
        // 初始化连接池
        for i := 0; i < poolInstance.maxConnections; i++ {
            db, _ := sql.Open("postgres", "connection_string")
            poolInstance.connections <- db
        }
    })
    return poolInstance
}

func (p *DatabasePool) GetConnection() *sql.DB {
    return <-p.connections
}

func (p *DatabasePool) ReturnConnection(db *sql.DB) {
    p.connections <- db
}
```

---

## 7. 相关模式

### 7.1 与工厂模式的关系

- **单例模式**: 确保只有一个实例
- **工厂模式**: 创建多个不同类型的实例

### 7.2 与享元模式的关系

- **单例模式**: 全局唯一实例
- **享元模式**: 共享多个相似实例

### 7.3 与注册表模式的关系

- **单例模式**: 单一访问点
- **注册表模式**: 多个实例的注册和查找

---

## 总结

单例模式是设计模式中最简单但最常用的模式之一。它通过确保类只有一个实例，提供了全局访问点，适用于需要全局状态管理的场景。在Go语言中，使用 `sync.Once` 可以优雅地实现线程安全的单例模式。

**关键要点**:
- 使用 `sync.Once` 保证线程安全
- 考虑延迟初始化的性能优势
- 注意内存泄漏和资源管理
- 合理选择实现方式（基础、泛型、函数式）

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **单例模式完成！** 🚀

**相关链接**:

- [02-工厂方法模式](../02-Factory-Method-Pattern.md)
- [03-抽象工厂模式](../03-Abstract-Factory-Pattern.md)
- [返回设计模式目录](../../README.md)
