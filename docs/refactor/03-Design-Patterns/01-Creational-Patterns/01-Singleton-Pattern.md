# 01-单例模式 (Singleton Pattern)

## 目录

- [01-单例模式 (Singleton Pattern)](#01-单例模式-singleton-pattern)
  - [目录](#目录)
  - [概述](#概述)
  - [1. 形式化定义](#1-形式化定义)
  - [2. 数学证明](#2-数学证明)
  - [3. 实现方式](#3-实现方式)
  - [4. Go语言实现](#4-go语言实现)
  - [5. 性能分析](#5-性能分析)
  - [6. 应用场景](#6-应用场景)
  - [7. 优缺点分析](#7-优缺点分析)
  - [8. 相关模式](#8-相关模式)
  - [参考文献](#参考文献)

## 概述

单例模式是一种创建型设计模式，确保一个类只有一个实例，并提供一个全局访问点。在Go语言中，单例模式通过包级别的变量和sync.Once来实现线程安全的单例。

### 核心特征

- **唯一性**: 确保类只有一个实例
- **全局访问**: 提供全局访问点
- **延迟初始化**: 实例在首次使用时创建
- **线程安全**: 在多线程环境下安全使用

## 1. 形式化定义

### 1.1 基本定义

**定义 1.1** (单例模式): 单例模式是一个三元组 $(C, \text{getInstance}, \text{instance})$，其中：
- $C$ 是单例类
- $\text{getInstance}$ 是获取实例的方法
- $\text{instance}$ 是唯一的实例

### 1.2 形式化约束

**约束 1.1** (唯一性): $\forall x, y \in C: \text{getInstance}() = x \land \text{getInstance}() = y \Rightarrow x = y$

**约束 1.2** (存在性): $\exists x \in C: \text{getInstance}() = x$

**约束 1.3** (全局访问): $\forall \text{context}: \text{getInstance}() \text{ is accessible}$

### 1.3 状态机模型

**定义 1.2** (单例状态机): 单例模式的状态机 $M = (Q, \Sigma, \delta, q_0, F)$ 定义为：
- $Q = \{\text{Uninitialized}, \text{Initialized}\}$
- $\Sigma = \{\text{getInstance}\}$
- $\delta(\text{Uninitialized}, \text{getInstance}) = \text{Initialized}$
- $\delta(\text{Initialized}, \text{getInstance}) = \text{Initialized}$
- $q_0 = \text{Uninitialized}$
- $F = \{\text{Initialized}\}$

## 2. 数学证明

### 2.1 唯一性证明

**定理 2.1** (唯一性): 单例模式确保实例的唯一性。

**证明**:
1. 假设存在两个不同的实例 $x$ 和 $y$
2. 根据约束1.1，$\text{getInstance}() = x$ 且 $\text{getInstance}() = y$
3. 因此 $x = y$，与假设矛盾
4. 所以实例是唯一的

### 2.2 线程安全性证明

**定理 2.2** (线程安全): 使用sync.Once的单例模式是线程安全的。

**证明**:
1. sync.Once保证Do方法只执行一次
2. 实例创建在Do方法中执行
3. 因此实例只创建一次
4. 所有线程访问同一个实例

### 2.3 延迟初始化证明

**定理 2.3** (延迟初始化): 单例模式支持延迟初始化。

**证明**:
1. 实例在首次调用getInstance时创建
2. 在此之前，实例为nil
3. 满足延迟初始化的定义

## 3. 实现方式

### 3.1 饿汉式单例

```go
// 饿汉式单例 - 在包初始化时创建实例
type EagerSingleton struct {
    data string
}

var eagerInstance = &EagerSingleton{data: "Eager Singleton"}

func GetEagerInstance() *EagerSingleton {
    return eagerInstance
}
```

### 3.2 懒汉式单例

```go
// 懒汉式单例 - 延迟初始化
type LazySingleton struct {
    data string
}

var (
    lazyInstance *LazySingleton
    lazyOnce     sync.Once
)

func GetLazyInstance() *LazySingleton {
    lazyOnce.Do(func() {
        lazyInstance = &LazySingleton{data: "Lazy Singleton"}
    })
    return lazyInstance
}
```

### 3.3 双重检查锁定

```go
// 双重检查锁定单例
type DoubleCheckSingleton struct {
    data string
}

var (
    doubleCheckInstance *DoubleCheckSingleton
    doubleCheckMutex    sync.RWMutex
)

func GetDoubleCheckInstance() *DoubleCheckSingleton {
    if doubleCheckInstance == nil {
        doubleCheckMutex.Lock()
        defer doubleCheckMutex.Unlock()
        
        if doubleCheckInstance == nil {
            doubleCheckInstance = &DoubleCheckSingleton{data: "Double Check Singleton"}
        }
    }
    return doubleCheckInstance
}
```

## 4. Go语言实现

### 4.1 基础单例实现

```go
// 基础单例接口
type Singleton interface {
    GetData() string
    SetData(data string)
}

// 基础单例实现
type BaseSingleton struct {
    data string
    mu   sync.RWMutex
}

var (
    baseInstance *BaseSingleton
    baseOnce     sync.Once
)

func GetBaseInstance() Singleton {
    baseOnce.Do(func() {
        baseInstance = &BaseSingleton{
            data: "Default Data",
        }
    })
    return baseInstance
}

func (s *BaseSingleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

func (s *BaseSingleton) SetData(data string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data = data
}
```

### 4.2 泛型单例实现

```go
// 泛型单例管理器
type SingletonManager[T any] struct {
    instance T
    once     sync.Once
    factory  func() T
}

func NewSingletonManager[T any](factory func() T) *SingletonManager[T] {
    return &SingletonManager[T]{
        factory: factory,
    }
}

func (sm *SingletonManager[T]) GetInstance() T {
    sm.once.Do(func() {
        sm.instance = sm.factory()
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

var configManager = NewSingletonManager(NewConfig)

func GetConfig() Config {
    return configManager.GetInstance()
}
```

### 4.3 函数式单例实现

```go
// 函数式单例
type FunctionalSingleton struct {
    data string
}

var (
    functionalInstance *FunctionalSingleton
    functionalOnce     sync.Once
)

// 使用闭包实现
func NewFunctionalSingleton() func() *FunctionalSingleton {
    return func() *FunctionalSingleton {
        functionalOnce.Do(func() {
            functionalInstance = &FunctionalSingleton{
                data: "Functional Singleton",
            }
        })
        return functionalInstance
    }
}

// 全局函数
var GetFunctionalInstance = NewFunctionalSingleton()
```

### 4.4 带配置的单例

```go
// 配置选项
type SingletonOption func(*ConfigurableSingleton)

func WithData(data string) SingletonOption {
    return func(s *ConfigurableSingleton) {
        s.data = data
    }
}

func WithTimeout(timeout time.Duration) SingletonOption {
    return func(s *ConfigurableSingleton) {
        s.timeout = timeout
    }
}

// 可配置单例
type ConfigurableSingleton struct {
    data    string
    timeout time.Duration
    mu      sync.RWMutex
}

var (
    configurableInstance *ConfigurableSingleton
    configurableOnce     sync.Once
)

func GetConfigurableInstance(options ...SingletonOption) *ConfigurableSingleton {
    configurableOnce.Do(func() {
        configurableInstance = &ConfigurableSingleton{
            data:    "Default Data",
            timeout: 30 * time.Second,
        }
        
        // 应用配置选项
        for _, option := range options {
            option(configurableInstance)
        }
    })
    return configurableInstance
}

func (s *ConfigurableSingleton) GetData() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

func (s *ConfigurableSingleton) GetTimeout() time.Duration {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.timeout
}
```

## 5. 性能分析

### 5.1 时间复杂度

**定理 5.1**: 单例模式的getInstance操作时间复杂度为 $O(1)$。

**证明**:
1. 首次调用需要创建实例，时间复杂度 $O(1)$
2. 后续调用直接返回实例，时间复杂度 $O(1)$
3. 总体时间复杂度为 $O(1)$

### 5.2 空间复杂度

**定理 5.2**: 单例模式的空间复杂度为 $O(1)$。

**证明**:
1. 只存储一个实例
2. 实例大小固定
3. 空间复杂度为 $O(1)$

### 5.3 并发性能

```go
// 性能测试
func BenchmarkSingleton(b *testing.B) {
    b.Run("BaseSingleton", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = GetBaseInstance()
        }
    })
    
    b.Run("LazySingleton", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = GetLazyInstance()
        }
    })
    
    b.Run("DoubleCheckSingleton", func(b *testing.B) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _ = GetDoubleCheckInstance()
        }
    })
}
```

## 6. 应用场景

### 6.1 配置管理

```go
// 配置管理器单例
type ConfigManager struct {
    config map[string]interface{}
    mu     sync.RWMutex
}

var (
    configInstance *ConfigManager
    configOnce     sync.Once
)

func GetConfigManager() *ConfigManager {
    configOnce.Do(func() {
        configInstance = &ConfigManager{
            config: make(map[string]interface{}),
        }
    })
    return configInstance
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

### 6.2 日志管理器

```go
// 日志管理器单例
type LogManager struct {
    logger *log.Logger
    mu     sync.Mutex
}

var (
    logInstance *LogManager
    logOnce     sync.Once
)

func GetLogManager() *LogManager {
    logOnce.Do(func() {
        logInstance = &LogManager{
            logger: log.New(os.Stdout, "[APP] ", log.LstdFlags),
        }
    })
    return logInstance
}

func (lm *LogManager) Log(message string) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    lm.logger.Println(message)
}
```

### 6.3 数据库连接池

```go
// 数据库连接池单例
type DBConnectionPool struct {
    connections chan *sql.DB
    mu          sync.Mutex
}

var (
    dbInstance *DBConnectionPool
    dbOnce     sync.Once
)

func GetDBConnectionPool() *DBConnectionPool {
    dbOnce.Do(func() {
        dbInstance = &DBConnectionPool{
            connections: make(chan *sql.DB, 10),
        }
    })
    return dbInstance
}

func (db *DBConnectionPool) GetConnection() *sql.DB {
    select {
    case conn := <-db.connections:
        return conn
    default:
        // 创建新连接
        return nil
    }
}
```

## 7. 优缺点分析

### 7.1 优点

1. **内存效率**: 只创建一个实例，节省内存
2. **全局访问**: 提供全局访问点
3. **延迟初始化**: 按需创建实例
4. **线程安全**: 支持并发访问

### 7.2 缺点

1. **全局状态**: 引入全局状态，可能影响测试
2. **违反单一职责**: 类既要管理实例又要提供业务功能
3. **难以扩展**: 难以支持多个实例
4. **生命周期管理**: 实例生命周期难以控制

### 7.3 改进方案

```go
// 改进的单例 - 支持重置
type ImprovedSingleton struct {
    data string
    mu   sync.RWMutex
}

var (
    improvedInstance *ImprovedSingleton
    improvedOnce     sync.Once
    resetMutex       sync.Mutex
)

func GetImprovedInstance() *ImprovedSingleton {
    improvedOnce.Do(func() {
        improvedInstance = &ImprovedSingleton{
            data: "Improved Singleton",
        }
    })
    return improvedInstance
}

// 支持重置（用于测试）
func ResetImprovedInstance() {
    resetMutex.Lock()
    defer resetMutex.Unlock()
    
    // 重置sync.Once
    improvedOnce = sync.Once{}
    improvedInstance = nil
}
```

## 8. 相关模式

### 8.1 与工厂模式的关系

单例模式可以用于实现工厂模式，确保工厂实例的唯一性。

### 8.2 与抽象工厂模式的关系

抽象工厂可以使用单例模式管理工厂实例。

### 8.3 与建造者模式的关系

建造者模式可以使用单例模式管理构建过程。

## 参考文献

1. Gamma, E., Helm, R., Johnson, R., & Vlissides, J. (1994). *Design Patterns: Elements of Reusable Object-Oriented Software*. Addison-Wesley.
2. Freeman, E., Robson, E., Sierra, K., & Bates, B. (2004). *Head First Design Patterns*. O'Reilly Media.
3. Goetz, B., Peierls, T., Bloch, J., Bowbeer, J., Holmes, D., & Lea, D. (2006). *Java Concurrency in Practice*. Addison-Wesley.

---

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **单例模式完成！** 🚀
