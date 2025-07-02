# 02-工厂方法模式 (Factory Method Pattern)

## 目录

- [02-工厂方法模式 (Factory Method Pattern)](#02-工厂方法模式-factory-method-pattern)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心特征](#12-核心特征)
    - [1.3 设计原则](#13-设计原则)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 集合论定义](#21-集合论定义)
    - [2.2 函数式定义](#22-函数式定义)
    - [2.3 类型论定义](#23-类型论定义)
  - [3. 数学证明](#3-数学证明)
    - [3.1 多态性证明](#31-多态性证明)
    - [3.2 扩展性证明](#32-扩展性证明)
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
    - [6.1 数据库连接工厂](#61-数据库连接工厂)
    - [6.2 日志记录器工厂](#62-日志记录器工厂)
    - [6.3 支付处理器工厂](#63-支付处理器工厂)
  - [7. 相关模式](#7-相关模式)
    - [7.1 与抽象工厂模式的关系](#71-与抽象工厂模式的关系)
    - [7.2 与简单工厂模式的关系](#72-与简单工厂模式的关系)
    - [7.3 与建造者模式的关系](#73-与建造者模式的关系)
  - [总结](#总结)

---

## 1. 概念与定义

### 1.1 基本概念

工厂方法模式是一种创建型设计模式，定义一个用于创建对象的接口，让子类决定实例化哪一个类。工厂方法使一个类的实例化延迟到其子类。

### 1.2 核心特征

- **抽象化**: 将对象创建过程抽象化
- **多态性**: 通过接口实现多态创建
- **扩展性**: 易于添加新的产品类型
- **封装性**: 隐藏对象创建的复杂性

### 1.3 设计原则

- **开闭原则**: 对扩展开放，对修改封闭
- **依赖倒置原则**: 依赖于抽象而不是具体实现
- **单一职责原则**: 每个工厂只负责创建特定类型的产品

---

## 2. 形式化定义

### 2.1 集合论定义

设 $$ P $$ 为产品集合，$$ F $$ 为工厂集合，$$ C $$ 为创建者集合，则工厂方法模式满足：

$$
\forall f \in F, \exists c \in C : f = \text{createProduct}(c)
$$

其中 $$ \text{createProduct}: C \rightarrow P $$ 为工厂方法。

### 2.2 函数式定义

定义工厂方法函数族 $$ \mathcal{F} = \{f_c : \emptyset \rightarrow P \mid c \in C\} $$，满足：

$$
f_c() = \text{createProduct}(c)
$$

### 2.3 类型论定义

在类型论中，工厂方法模式可以表示为：

$$
\text{Factory} = \Pi_{c:C} \Sigma_{p:P} \text{Product}(p)
$$

其中 $$ \text{Product}(p) $$ 表示产品 $$ p $$ 的类型。

---

## 3. 数学证明

### 3.1 多态性证明

**定理**: 工厂方法模式支持多态创建。

**证明**:

1. 设 $$ P_1, P_2 $$ 为两个不同的产品类型
2. 存在工厂方法 $$ f_1, f_2 $$ 分别创建 $$ P_1, P_2 $$
3. 通过接口 $$ I $$，$$ f_1, f_2 $$ 都实现相同的签名
4. 因此支持多态调用：$$ \text{create}(f_1) \neq \text{create}(f_2) $$
5. 多态性得证

### 3.2 扩展性证明

**定理**: 工厂方法模式支持开闭原则。

**证明**:

1. 设现有工厂集合 $$ F = \{f_1, f_2, \ldots, f_n\} $$
2. 添加新工厂 $$ f_{n+1} $$ 时，只需实现相同的接口
3. 不需要修改现有代码
4. 因此满足开闭原则

---

## 4. Go语言实现

### 4.1 基础实现

```go
package factory

import "fmt"

// Product 产品接口
type Product interface {
    Operation() string
    GetName() string
}

// ConcreteProductA 具体产品A
type ConcreteProductA struct {
    name string
}

func NewConcreteProductA() *ConcreteProductA {
    return &ConcreteProductA{
        name: "Product A",
    }
}

func (p *ConcreteProductA) Operation() string {
    return "Result of ConcreteProductA"
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

// ConcreteProductB 具体产品B
type ConcreteProductB struct {
    name string
}

func NewConcreteProductB() *ConcreteProductB {
    return &ConcreteProductB{
        name: "Product B",
    }
}

func (p *ConcreteProductB) Operation() string {
    return "Result of ConcreteProductB"
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

// Creator 创建者接口
type Creator interface {
    FactoryMethod() Product
    SomeOperation() string
}

// ConcreteCreatorA 具体创建者A
type ConcreteCreatorA struct{}

func NewConcreteCreatorA() *ConcreteCreatorA {
    return &ConcreteCreatorA{}
}

func (c *ConcreteCreatorA) FactoryMethod() Product {
    return NewConcreteProductA()
}

func (c *ConcreteCreatorA) SomeOperation() string {
    product := c.FactoryMethod()
    return fmt.Sprintf("Creator A: %s", product.Operation())
}

// ConcreteCreatorB 具体创建者B
type ConcreteCreatorB struct{}

func NewConcreteCreatorB() *ConcreteCreatorB {
    return &ConcreteCreatorB{}
}

func (c *ConcreteCreatorB) FactoryMethod() Product {
    return NewConcreteProductB()
}

func (c *ConcreteCreatorB) SomeOperation() string {
    product := c.FactoryMethod()
    return fmt.Sprintf("Creator B: %s", product.Operation())
}

// ClientCode 客户端代码
func ClientCode(creator Creator) {
    fmt.Println(creator.SomeOperation())
}
```

### 4.2 泛型实现

```go
package factory

import (
    "fmt"
    "reflect"
)

// GenericProduct 泛型产品接口
type GenericProduct[T any] interface {
    Operation() T
    GetType() string
}

// GenericCreator 泛型创建者接口
type GenericCreator[T any] interface {
    FactoryMethod() GenericProduct[T]
    SomeOperation() T
}

// GenericFactory 泛型工厂
type GenericFactory[T any] struct {
    productType reflect.Type
    factoryFunc func() GenericProduct[T]
}

// NewGenericFactory 创建泛型工厂
func NewGenericFactory[T any](factoryFunc func() GenericProduct[T]) *GenericFactory[T] {
    return &GenericFactory[T]{
        productType: reflect.TypeOf((*T)(nil)).Elem(),
        factoryFunc: factoryFunc,
    }
}

// CreateProduct 创建产品
func (f *GenericFactory[T]) CreateProduct() GenericProduct[T] {
    return f.factoryFunc()
}

// GetProductType 获取产品类型
func (f *GenericFactory[T]) GetProductType() string {
    return f.productType.String()
}

// 使用示例
type StringProduct struct {
    value string
}

func NewStringProduct() GenericProduct[string] {
    return &StringProduct{
        value: "String Product",
    }
}

func (p *StringProduct) Operation() string {
    return p.value
}

func (p *StringProduct) GetType() string {
    return "StringProduct"
}

type IntProduct struct {
    value int
}

func NewIntProduct() GenericProduct[int] {
    return &IntProduct{
        value: 42,
    }
}

func (p *IntProduct) Operation() int {
    return p.value
}

func (p *IntProduct) GetType() string {
    return "IntProduct"
}

// 全局工厂实例
var (
    stringFactory = NewGenericFactory(NewStringProduct)
    intFactory    = NewGenericFactory(NewIntProduct)
)
```

### 4.3 函数式实现

```go
package factory

import (
    "fmt"
    "sync"
)

// FunctionalProduct 函数式产品
type FunctionalProduct struct {
    operation func() string
    name      string
}

// NewFunctionalProduct 创建函数式产品
func NewFunctionalProduct(name string, operation func() string) *FunctionalProduct {
    return &FunctionalProduct{
        operation: operation,
        name:      name,
    }
}

func (p *FunctionalProduct) Execute() string {
    return p.operation()
}

func (p *FunctionalProduct) GetName() string {
    return p.name
}

// FunctionalFactory 函数式工厂
type FunctionalFactory struct {
    factories map[string]func() *FunctionalProduct
    mutex     sync.RWMutex
}

// NewFunctionalFactory 创建函数式工厂
func NewFunctionalFactory() *FunctionalFactory {
    return &FunctionalFactory{
        factories: make(map[string]func() *FunctionalProduct),
    }
}

// RegisterFactory 注册工厂方法
func (f *FunctionalFactory) RegisterFactory(name string, factory func() *FunctionalProduct) {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    f.factories[name] = factory
}

// CreateProduct 创建产品
func (f *FunctionalFactory) CreateProduct(name string) (*FunctionalProduct, error) {
    f.mutex.RLock()
    defer f.mutex.RUnlock()
    
    factory, exists := f.factories[name]
    if !exists {
        return nil, fmt.Errorf("factory not found: %s", name)
    }
    
    return factory(), nil
}

// GetAvailableFactories 获取可用的工厂列表
func (f *FunctionalFactory) GetAvailableFactories() []string {
    f.mutex.RLock()
    defer f.mutex.RUnlock()
    
    names := make([]string, 0, len(f.factories))
    for name := range f.factories {
        names = append(names, name)
    }
    return names
}
```

### 4.4 测试代码

```go
package factory

import (
    "testing"
)

// TestFactoryMethod 测试工厂方法
func TestFactoryMethod(t *testing.T) {
    creatorA := NewConcreteCreatorA()
    creatorB := NewConcreteCreatorB()
    
    productA := creatorA.FactoryMethod()
    productB := creatorB.FactoryMethod()
    
    if productA.GetName() != "Product A" {
        t.Errorf("Expected Product A, got %s", productA.GetName())
    }
    
    if productB.GetName() != "Product B" {
        t.Errorf("Expected Product B, got %s", productB.GetName())
    }
    
    if productA.Operation() == productB.Operation() {
        t.Error("Different products should have different operations")
    }
}

// TestGenericFactory 测试泛型工厂
func TestGenericFactory(t *testing.T) {
    stringProduct := stringFactory.CreateProduct()
    intProduct := intFactory.CreateProduct()
    
    if stringProduct.GetType() != "StringProduct" {
        t.Errorf("Expected StringProduct, got %s", stringProduct.GetType())
    }
    
    if intProduct.GetType() != "IntProduct" {
        t.Errorf("Expected IntProduct, got %s", intProduct.GetType())
    }
    
    stringResult := stringProduct.Operation()
    intResult := intProduct.Operation()
    
    if stringResult != "String Product" {
        t.Errorf("Expected 'String Product', got %s", stringResult)
    }
    
    if intResult != 42 {
        t.Errorf("Expected 42, got %d", intResult)
    }
}

// TestFunctionalFactory 测试函数式工厂
func TestFunctionalFactory(t *testing.T) {
    factory := NewFunctionalFactory()
    
    // 注册工厂方法
    factory.RegisterFactory("greeting", func() *FunctionalProduct {
        return NewFunctionalProduct("Greeting", func() string {
            return "Hello, World!"
        })
    })
    
    factory.RegisterFactory("calculation", func() *FunctionalProduct {
        return NewFunctionalProduct("Calculation", func() string {
            return "2 + 2 = 4"
        })
    })
    
    // 创建产品
    greeting, err := factory.CreateProduct("greeting")
    if err != nil {
        t.Errorf("Failed to create greeting product: %v", err)
    }
    
    calculation, err := factory.CreateProduct("calculation")
    if err != nil {
        t.Errorf("Failed to create calculation product: %v", err)
    }
    
    if greeting.Execute() != "Hello, World!" {
        t.Errorf("Expected 'Hello, World!', got %s", greeting.Execute())
    }
    
    if calculation.Execute() != "2 + 2 = 4" {
        t.Errorf("Expected '2 + 2 = 4', got %s", calculation.Execute())
    }
    
    // 测试不存在的工厂
    _, err = factory.CreateProduct("nonexistent")
    if err == nil {
        t.Error("Expected error for nonexistent factory")
    }
}

// BenchmarkFactoryMethod 性能基准测试
func BenchmarkFactoryMethod(b *testing.B) {
    creator := NewConcreteCreatorA()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        creator.FactoryMethod()
    }
}
```

---

## 5. 性能分析

### 5.1 时间复杂度

- **创建产品**: $$ O(1) $$
- **工厂注册**: $$ O(1) $$
- **产品查找**: $$ O(1) $$ (使用map)

### 5.2 空间复杂度

- **工厂存储**: $$ O(n) $$ (n为工厂数量)
- **产品实例**: $$ O(1) $$ 每个产品

### 5.3 性能优化

```go
// 缓存工厂实现
type CachedFactory struct {
    cache  map[string]Product
    mutex  sync.RWMutex
    factory func(string) Product
}

func NewCachedFactory(factory func(string) Product) *CachedFactory {
    return &CachedFactory{
        cache:   make(map[string]Product),
        factory: factory,
    }
}

func (cf *CachedFactory) CreateProduct(key string) Product {
    // 先检查缓存
    cf.mutex.RLock()
    if product, exists := cf.cache[key]; exists {
        cf.mutex.RUnlock()
        return product
    }
    cf.mutex.RUnlock()
    
    // 创建新产品
    cf.mutex.Lock()
    defer cf.mutex.Unlock()
    
    // 双重检查
    if product, exists := cf.cache[key]; exists {
        return product
    }
    
    product := cf.factory(key)
    cf.cache[key] = product
    return product
}
```

---

## 6. 应用场景

### 6.1 数据库连接工厂

```go
// 数据库连接工厂
type DatabaseConnection interface {
    Connect() error
    Disconnect() error
    Execute(query string) (interface{}, error)
}

type MySQLConnection struct {
    host     string
    port     int
    database string
    username string
    password string
}

func NewMySQLConnection(config map[string]interface{}) *MySQLConnection {
    return &MySQLConnection{
        host:     config["host"].(string),
        port:     config["port"].(int),
        database: config["database"].(string),
        username: config["username"].(string),
        password: config["password"].(string),
    }
}

func (m *MySQLConnection) Connect() error {
    // MySQL连接逻辑
    return nil
}

func (m *MySQLConnection) Disconnect() error {
    // MySQL断开连接逻辑
    return nil
}

func (m *MySQLConnection) Execute(query string) (interface{}, error) {
    // MySQL执行查询逻辑
    return nil, nil
}

type PostgreSQLConnection struct {
    host     string
    port     int
    database string
    username string
    password string
}

func NewPostgreSQLConnection(config map[string]interface{}) *PostgreSQLConnection {
    return &PostgreSQLConnection{
        host:     config["host"].(string),
        port:     config["port"].(int),
        database: config["database"].(string),
        username: config["username"].(string),
        password: config["password"].(string),
    }
}

func (p *PostgreSQLConnection) Connect() error {
    // PostgreSQL连接逻辑
    return nil
}

func (p *PostgreSQLConnection) Disconnect() error {
    // PostgreSQL断开连接逻辑
    return nil
}

func (p *PostgreSQLConnection) Execute(query string) (interface{}, error) {
    // PostgreSQL执行查询逻辑
    return nil, nil
}

// 数据库工厂
type DatabaseFactory struct{}

func NewDatabaseFactory() *DatabaseFactory {
    return &DatabaseFactory{}
}

func (df *DatabaseFactory) CreateConnection(dbType string, config map[string]interface{}) (DatabaseConnection, error) {
    switch dbType {
    case "mysql":
        return NewMySQLConnection(config), nil
    case "postgresql":
        return NewPostgreSQLConnection(config), nil
    default:
        return nil, fmt.Errorf("unsupported database type: %s", dbType)
    }
}
```

### 6.2 日志记录器工厂

```go
// 日志记录器工厂
type Logger interface {
    Log(level, message string) error
    SetLevel(level string)
    GetLevel() string
}

type ConsoleLogger struct {
    level string
}

func NewConsoleLogger() *ConsoleLogger {
    return &ConsoleLogger{
        level: "INFO",
    }
}

func (c *ConsoleLogger) Log(level, message string) error {
    fmt.Printf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, message)
    return nil
}

func (c *ConsoleLogger) SetLevel(level string) {
    c.level = level
}

func (c *ConsoleLogger) GetLevel() string {
    return c.level
}

type FileLogger struct {
    level string
    file  *os.File
}

func NewFileLogger(filename string) (*FileLogger, error) {
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }
    
    return &FileLogger{
        level: "INFO",
        file:  file,
    }, nil
}

func (f *FileLogger) Log(level, message string) error {
    _, err := fmt.Fprintf(f.file, "[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, message)
    return err
}

func (f *FileLogger) SetLevel(level string) {
    f.level = level
}

func (f *FileLogger) GetLevel() string {
    return f.level
}

// 日志工厂
type LoggerFactory struct{}

func NewLoggerFactory() *LoggerFactory {
    return &LoggerFactory{}
}

func (lf *LoggerFactory) CreateLogger(loggerType string, config map[string]interface{}) (Logger, error) {
    switch loggerType {
    case "console":
        return NewConsoleLogger(), nil
    case "file":
        filename, ok := config["filename"].(string)
        if !ok {
            return nil, fmt.Errorf("filename is required for file logger")
        }
        return NewFileLogger(filename)
    default:
        return nil, fmt.Errorf("unsupported logger type: %s", loggerType)
    }
}
```

### 6.3 支付处理器工厂

```go
// 支付处理器工厂
type PaymentProcessor interface {
    ProcessPayment(amount float64, currency string) error
    RefundPayment(transactionID string) error
    GetSupportedCurrencies() []string
}

type CreditCardProcessor struct {
    apiKey string
}

func NewCreditCardProcessor(apiKey string) *CreditCardProcessor {
    return &CreditCardProcessor{
        apiKey: apiKey,
    }
}

func (c *CreditCardProcessor) ProcessPayment(amount float64, currency string) error {
    // 信用卡支付处理逻辑
    return nil
}

func (c *CreditCardProcessor) RefundPayment(transactionID string) error {
    // 信用卡退款逻辑
    return nil
}

func (c *CreditCardProcessor) GetSupportedCurrencies() []string {
    return []string{"USD", "EUR", "GBP", "JPY"}
}

type PayPalProcessor struct {
    clientID     string
    clientSecret string
}

func NewPayPalProcessor(clientID, clientSecret string) *PayPalProcessor {
    return &PayPalProcessor{
        clientID:     clientID,
        clientSecret: clientSecret,
    }
}

func (p *PayPalProcessor) ProcessPayment(amount float64, currency string) error {
    // PayPal支付处理逻辑
    return nil
}

func (p *PayPalProcessor) RefundPayment(transactionID string) error {
    // PayPal退款逻辑
    return nil
}

func (p *PayPalProcessor) GetSupportedCurrencies() []string {
    return []string{"USD", "EUR", "GBP", "CAD", "AUD"}
}

// 支付工厂
type PaymentFactory struct{}

func NewPaymentFactory() *PaymentFactory {
    return &PaymentFactory{}
}

func (pf *PaymentFactory) CreateProcessor(processorType string, config map[string]interface{}) (PaymentProcessor, error) {
    switch processorType {
    case "creditcard":
        apiKey, ok := config["apiKey"].(string)
        if !ok {
            return nil, fmt.Errorf("apiKey is required for credit card processor")
        }
        return NewCreditCardProcessor(apiKey), nil
    case "paypal":
        clientID, ok1 := config["clientID"].(string)
        clientSecret, ok2 := config["clientSecret"].(string)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("clientID and clientSecret are required for PayPal processor")
        }
        return NewPayPalProcessor(clientID, clientSecret), nil
    default:
        return nil, fmt.Errorf("unsupported payment processor type: %s", processorType)
    }
}
```

---

## 7. 相关模式

### 7.1 与抽象工厂模式的关系

- **工厂方法模式**: 创建单个产品
- **抽象工厂模式**: 创建产品族

### 7.2 与简单工厂模式的关系

- **工厂方法模式**: 多态创建，支持扩展
- **简单工厂模式**: 静态创建，不易扩展

### 7.3 与建造者模式的关系

- **工厂方法模式**: 创建简单对象
- **建造者模式**: 创建复杂对象

---

## 总结

工厂方法模式通过抽象化对象创建过程，实现了创建者与产品的解耦。它支持多态创建，易于扩展，是面向对象设计中重要的创建型模式。

**关键要点**:

- 使用接口实现多态创建
- 支持开闭原则，易于扩展
- 隐藏对象创建的复杂性
- 合理选择实现方式（基础、泛型、函数式）

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **工厂方法模式完成！** 🚀

**相关链接**:

- [01-单例模式](../01-Singleton-Pattern.md)
- [03-抽象工厂模式](../03-Abstract-Factory-Pattern.md)
- [返回设计模式目录](../../README.md)
