# 03-抽象工厂模式 (Abstract Factory Pattern)

## 目录

- [03-抽象工厂模式 (Abstract Factory Pattern)](#03-抽象工厂模式-abstract-factory-pattern)
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
    - [3.1 一致性证明](#31-一致性证明)
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
    - [6.1 GUI组件工厂](#61-gui组件工厂)
    - [6.2 数据库抽象工厂](#62-数据库抽象工厂)
    - [6.3 支付系统抽象工厂](#63-支付系统抽象工厂)
  - [7. 相关模式](#7-相关模式)
    - [7.1 与工厂方法模式的关系](#71-与工厂方法模式的关系)
    - [7.2 与建造者模式的关系](#72-与建造者模式的关系)
    - [7.3 与单例模式的关系](#73-与单例模式的关系)
  - [总结](#总结)

---

## 1. 概念与定义

### 1.1 基本概念

抽象工厂模式是一种创建型设计模式，提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们具体的类。抽象工厂模式围绕一个超级工厂创建其他工厂。

### 1.2 核心特征

- **产品族**: 创建一系列相关的产品
- **一致性**: 确保产品之间的兼容性
- **封装性**: 隐藏产品创建的复杂性
- **扩展性**: 支持添加新的产品族

### 1.3 设计原则

- **开闭原则**: 对扩展开放，对修改封闭
- **依赖倒置原则**: 依赖于抽象而不是具体实现
- **单一职责原则**: 每个工厂只负责创建特定产品族

---

## 2. 形式化定义

### 2.1 集合论定义

设 $P$ 为产品集合，$F$ 为工厂集合，$PF$ 为产品族集合，则抽象工厂模式满足：

$$\forall pf \in PF, \exists f \in F : pf = \{p \mid p \in P \land \text{createProduct}(f, p)\}$$

其中 $\text{createProduct}: F \times P \rightarrow P$ 为抽象工厂方法。

### 2.2 函数式定义

定义抽象工厂函数族 $\mathcal{AF} = \{af_f : PF \rightarrow P^+ \mid f \in F\}$，满足：

$$af_f(pf) = \{p \mid p \in pf \land \text{compatible}(p, pf)\}$$

其中 $\text{compatible}(p, pf)$ 表示产品 $p$ 与产品族 $pf$ 兼容。

### 2.3 类型论定义

在类型论中，抽象工厂模式可以表示为：

$$\text{AbstractFactory} = \Pi_{f:F} \Sigma_{pf:PF} \Pi_{p:pf} \text{Product}(p)$$

其中 $\text{Product}(p)$ 表示产品 $p$ 的类型。

---

## 3. 数学证明

### 3.1 一致性证明

**定理**: 抽象工厂模式保证产品族的一致性。

**证明**:

1. 设 $pf_1, pf_2$ 为两个不同的产品族
2. 对于任意产品 $p_1 \in pf_1, p_2 \in pf_2$
3. 抽象工厂确保 $\text{compatible}(p_1, pf_1) \land \text{compatible}(p_2, pf_2)$
4. 因此产品族内部一致性得证

### 3.2 扩展性证明

**定理**: 抽象工厂模式支持开闭原则。

**证明**:

1. 设现有工厂集合 $F = \{f_1, f_2, \ldots, f_n\}$
2. 添加新工厂 $f_{n+1}$ 时，只需实现相同的抽象接口
3. 不需要修改现有代码
4. 因此满足开闭原则

---

## 4. Go语言实现

### 4.1 基础实现

```go
package abstractfactory

import "fmt"

// AbstractProductA 抽象产品A接口
type AbstractProductA interface {
    UsefulFunctionA() string
}

// AbstractProductB 抽象产品B接口
type AbstractProductB interface {
    UsefulFunctionB() string
    AnotherUsefulFunctionB(collaborator AbstractProductA) string
}

// ConcreteProductA1 具体产品A1
type ConcreteProductA1 struct{}

func (p *ConcreteProductA1) UsefulFunctionA() string {
    return "The result of the product A1."
}

// ConcreteProductA2 具体产品A2
type ConcreteProductA2 struct{}

func (p *ConcreteProductA2) UsefulFunctionA() string {
    return "The result of the product A2."
}

// ConcreteProductB1 具体产品B1
type ConcreteProductB1 struct{}

func (p *ConcreteProductB1) UsefulFunctionB() string {
    return "The result of the product B1."
}

func (p *ConcreteProductB1) AnotherUsefulFunctionB(collaborator AbstractProductA) string {
    result := collaborator.UsefulFunctionA()
    return fmt.Sprintf("The result of the B1 collaborating with the (%s)", result)
}

// ConcreteProductB2 具体产品B2
type ConcreteProductB2 struct{}

func (p *ConcreteProductB2) UsefulFunctionB() string {
    return "The result of the product B2."
}

func (p *ConcreteProductB2) AnotherUsefulFunctionB(collaborator AbstractProductA) string {
    result := collaborator.UsefulFunctionA()
    return fmt.Sprintf("The result of the B2 collaborating with the (%s)", result)
}

// AbstractFactory 抽象工厂接口
type AbstractFactory interface {
    CreateProductA() AbstractProductA
    CreateProductB() AbstractProductB
}

// ConcreteFactory1 具体工厂1
type ConcreteFactory1 struct{}

func NewConcreteFactory1() *ConcreteFactory1 {
    return &ConcreteFactory1{}
}

func (f *ConcreteFactory1) CreateProductA() AbstractProductA {
    return &ConcreteProductA1{}
}

func (f *ConcreteFactory1) CreateProductB() AbstractProductB {
    return &ConcreteProductB1{}
}

// ConcreteFactory2 具体工厂2
type ConcreteFactory2 struct{}

func NewConcreteFactory2() *ConcreteFactory2 {
    return &ConcreteFactory2{}
}

func (f *ConcreteFactory2) CreateProductA() AbstractProductA {
    return &ConcreteProductA2{}
}

func (f *ConcreteFactory2) CreateProductB() AbstractProductB {
    return &ConcreteProductB2{}
}

// ClientCode 客户端代码
func ClientCode(factory AbstractFactory) {
    productA := factory.CreateProductA()
    productB := factory.CreateProductB()

    fmt.Println(productB.UsefulFunctionB())
    fmt.Println(productB.AnotherUsefulFunctionB(productA))
}
```

### 4.2 泛型实现

```go
package abstractfactory

import (
    "fmt"
    "reflect"
)

// GenericProduct 泛型产品接口
type GenericProduct[T any] interface {
    Operation() T
    GetType() string
}

// GenericFactory 泛型工厂接口
type GenericFactory[T any] interface {
    CreateProductA() GenericProduct[T]
    CreateProductB() GenericProduct[T]
    GetFactoryType() string
}

// GenericAbstractFactory 泛型抽象工厂
type GenericAbstractFactory[T any] struct {
    factoryType reflect.Type
    factoryFunc func() GenericFactory[T]
}

// NewGenericAbstractFactory 创建泛型抽象工厂
func NewGenericAbstractFactory[T any](factoryFunc func() GenericFactory[T]) *GenericAbstractFactory[T] {
    return &GenericAbstractFactory[T]{
        factoryType: reflect.TypeOf((*T)(nil)).Elem(),
        factoryFunc: factoryFunc,
    }
}

// CreateFactory 创建工厂
func (af *GenericAbstractFactory[T]) CreateFactory() GenericFactory[T] {
    return af.factoryFunc()
}

// GetFactoryType 获取工厂类型
func (af *GenericAbstractFactory[T]) GetFactoryType() string {
    return af.factoryType.String()
}

// 使用示例
type StringProductA struct {
    value string
}

func NewStringProductA() GenericProduct[string] {
    return &StringProductA{
        value: "String Product A",
    }
}

func (p *StringProductA) Operation() string {
    return p.value
}

func (p *StringProductA) GetType() string {
    return "StringProductA"
}

type StringProductB struct {
    value string
}

func NewStringProductB() GenericProduct[string] {
    return &StringProductB{
        value: "String Product B",
    }
}

func (p *StringProductB) Operation() string {
    return p.value
}

func (p *StringProductB) GetType() string {
    return "StringProductB"
}

type StringFactory struct{}

func NewStringFactory() GenericFactory[string] {
    return &StringFactory{}
}

func (f *StringFactory) CreateProductA() GenericProduct[string] {
    return NewStringProductA()
}

func (f *StringFactory) CreateProductB() GenericProduct[string] {
    return NewStringProductB()
}

func (f *StringFactory) GetFactoryType() string {
    return "StringFactory"
}

// 全局抽象工厂实例
var stringAbstractFactory = NewGenericAbstractFactory(NewStringFactory)
```

### 4.3 函数式实现

```go
package abstractfactory

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
    products map[string]*FunctionalProduct
    mutex    sync.RWMutex
}

// NewFunctionalFactory 创建函数式工厂
func NewFunctionalFactory() *FunctionalFactory {
    return &FunctionalFactory{
        products: make(map[string]*FunctionalProduct),
    }
}

// RegisterProduct 注册产品
func (f *FunctionalFactory) RegisterProduct(name string, product *FunctionalProduct) {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    f.products[name] = product
}

// CreateProduct 创建产品
func (f *FunctionalFactory) CreateProduct(name string) (*FunctionalProduct, error) {
    f.mutex.RLock()
    defer f.mutex.RUnlock()
    
    product, exists := f.products[name]
    if !exists {
        return nil, fmt.Errorf("product not found: %s", name)
    }
    
    return product, nil
}

// GetAvailableProducts 获取可用的产品列表
func (f *FunctionalFactory) GetAvailableProducts() []string {
    f.mutex.RLock()
    defer f.mutex.RUnlock()
    
    names := make([]string, 0, len(f.products))
    for name := range f.products {
        names = append(names, name)
    }
    return names
}

// FunctionalAbstractFactory 函数式抽象工厂
type FunctionalAbstractFactory struct {
    factories map[string]*FunctionalFactory
    mutex     sync.RWMutex
}

// NewFunctionalAbstractFactory 创建函数式抽象工厂
func NewFunctionalAbstractFactory() *FunctionalAbstractFactory {
    return &FunctionalAbstractFactory{
        factories: make(map[string]*FunctionalFactory),
    }
}

// RegisterFactory 注册工厂
func (af *FunctionalAbstractFactory) RegisterFactory(name string, factory *FunctionalFactory) {
    af.mutex.Lock()
    defer af.mutex.Unlock()
    af.factories[name] = factory
}

// CreateFactory 创建工厂
func (af *FunctionalAbstractFactory) CreateFactory(name string) (*FunctionalFactory, error) {
    af.mutex.RLock()
    defer af.mutex.RUnlock()
    
    factory, exists := af.factories[name]
    if !exists {
        return nil, fmt.Errorf("factory not found: %s", name)
    }
    
    return factory, nil
}

// GetAvailableFactories 获取可用的工厂列表
func (af *FunctionalAbstractFactory) GetAvailableFactories() []string {
    af.mutex.RLock()
    defer af.mutex.RUnlock()
    
    names := make([]string, 0, len(af.factories))
    for name := range af.factories {
        names = append(names, name)
    }
    return names
}
```

### 4.4 测试代码

```go
package abstractfactory

import (
    "testing"
)

// TestAbstractFactory 测试抽象工厂
func TestAbstractFactory(t *testing.T) {
    factory1 := NewConcreteFactory1()
    factory2 := NewConcreteFactory2()
    
    // 测试工厂1
    productA1 := factory1.CreateProductA()
    productB1 := factory1.CreateProductB()
    
    if productA1.UsefulFunctionA() != "The result of the product A1." {
        t.Errorf("Expected A1 result, got %s", productA1.UsefulFunctionA())
    }
    
    if productB1.UsefulFunctionB() != "The result of the product B1." {
        t.Errorf("Expected B1 result, got %s", productB1.UsefulFunctionB())
    }
    
    // 测试工厂2
    productA2 := factory2.CreateProductA()
    productB2 := factory2.CreateProductB()
    
    if productA2.UsefulFunctionA() != "The result of the product A2." {
        t.Errorf("Expected A2 result, got %s", productA2.UsefulFunctionA())
    }
    
    if productB2.UsefulFunctionB() != "The result of the product B2." {
        t.Errorf("Expected B2 result, got %s", productB2.UsefulFunctionB())
    }
    
    // 测试产品协作
    collaboration1 := productB1.AnotherUsefulFunctionB(productA1)
    collaboration2 := productB2.AnotherUsefulFunctionB(productA2)
    
    if collaboration1 == collaboration2 {
        t.Error("Different factories should produce different collaborations")
    }
}

// TestGenericAbstractFactory 测试泛型抽象工厂
func TestGenericAbstractFactory(t *testing.T) {
    factory := stringAbstractFactory.CreateFactory()
    
    if factory.GetFactoryType() != "StringFactory" {
        t.Errorf("Expected StringFactory, got %s", factory.GetFactoryType())
    }
    
    productA := factory.CreateProductA()
    productB := factory.CreateProductB()
    
    if productA.GetType() != "StringProductA" {
        t.Errorf("Expected StringProductA, got %s", productA.GetType())
    }
    
    if productB.GetType() != "StringProductB" {
        t.Errorf("Expected StringProductB, got %s", productB.GetType())
    }
    
    stringResultA := productA.Operation()
    stringResultB := productB.Operation()
    
    if stringResultA != "String Product A" {
        t.Errorf("Expected 'String Product A', got %s", stringResultA)
    }
    
    if stringResultB != "String Product B" {
        t.Errorf("Expected 'String Product B', got %s", stringResultB)
    }
}

// TestFunctionalAbstractFactory 测试函数式抽象工厂
func TestFunctionalAbstractFactory(t *testing.T) {
    abstractFactory := NewFunctionalAbstractFactory()
    
    // 创建工厂1
    factory1 := NewFunctionalFactory()
    factory1.RegisterProduct("A1", NewFunctionalProduct("A1", func() string {
        return "Product A1"
    }))
    factory1.RegisterProduct("B1", NewFunctionalProduct("B1", func() string {
        return "Product B1"
    }))
    
    // 创建工厂2
    factory2 := NewFunctionalFactory()
    factory2.RegisterProduct("A2", NewFunctionalProduct("A2", func() string {
        return "Product A2"
    }))
    factory2.RegisterProduct("B2", NewFunctionalProduct("B2", func() string {
        return "Product B2"
    }))
    
    // 注册工厂
    abstractFactory.RegisterFactory("Factory1", factory1)
    abstractFactory.RegisterFactory("Factory2", factory2)
    
    // 测试工厂1
    createdFactory1, err := abstractFactory.CreateFactory("Factory1")
    if err != nil {
        t.Errorf("Failed to create Factory1: %v", err)
    }
    
    productA1, err := createdFactory1.CreateProduct("A1")
    if err != nil {
        t.Errorf("Failed to create product A1: %v", err)
    }
    
    if productA1.Execute() != "Product A1" {
        t.Errorf("Expected 'Product A1', got %s", productA1.Execute())
    }
    
    // 测试工厂2
    createdFactory2, err := abstractFactory.CreateFactory("Factory2")
    if err != nil {
        t.Errorf("Failed to create Factory2: %v", err)
    }
    
    productA2, err := createdFactory2.CreateProduct("A2")
    if err != nil {
        t.Errorf("Failed to create product A2: %v", err)
    }
    
    if productA2.Execute() != "Product A2" {
        t.Errorf("Expected 'Product A2', got %s", productA2.Execute())
    }
    
    // 测试不存在的工厂
    _, err = abstractFactory.CreateFactory("NonexistentFactory")
    if err == nil {
        t.Error("Expected error for nonexistent factory")
    }
}

// BenchmarkAbstractFactory 性能基准测试
func BenchmarkAbstractFactory(b *testing.B) {
    factory := NewConcreteFactory1()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        factory.CreateProductA()
        factory.CreateProductB()
    }
}
```

---

## 5. 性能分析

### 5.1 时间复杂度

- **创建产品**: $O(1)$
- **工厂创建**: $O(1)$
- **产品查找**: $O(1)$ (使用map)

### 5.2 空间复杂度

- **工厂存储**: $O(n)$ (n为工厂数量)
- **产品存储**: $O(m)$ (m为产品数量)

### 5.3 性能优化

```go
// 缓存抽象工厂实现
type CachedAbstractFactory struct {
    cache  map[string]AbstractFactory
    mutex  sync.RWMutex
    factory func(string) AbstractFactory
}

func NewCachedAbstractFactory(factory func(string) AbstractFactory) *CachedAbstractFactory {
    return &CachedAbstractFactory{
        cache:   make(map[string]AbstractFactory),
        factory: factory,
    }
}

func (caf *CachedAbstractFactory) CreateFactory(key string) AbstractFactory {
    // 先检查缓存
    caf.mutex.RLock()
    if factory, exists := caf.cache[key]; exists {
        caf.mutex.RUnlock()
        return factory
    }
    caf.mutex.RUnlock()
    
    // 创建新工厂
    caf.mutex.Lock()
    defer caf.mutex.Unlock()
    
    // 双重检查
    if factory, exists := caf.cache[key]; exists {
        return factory
    }
    
    factory := caf.factory(key)
    caf.cache[key] = factory
    return factory
}
```

---

## 6. 应用场景

### 6.1 GUI组件工厂

```go
// GUI组件抽象工厂
type Button interface {
    Render() string
    Click() string
}

type Checkbox interface {
    Render() string
    Check() string
}

type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

// Windows GUI组件
type WindowsButton struct{}

func (w *WindowsButton) Render() string {
    return "Windows Button"
}

func (w *WindowsButton) Click() string {
    return "Windows Button Clicked"
}

type WindowsCheckbox struct{}

func (w *WindowsCheckbox) Render() string {
    return "Windows Checkbox"
}

func (w *WindowsCheckbox) Check() string {
    return "Windows Checkbox Checked"
}

type WindowsFactory struct{}

func NewWindowsFactory() *WindowsFactory {
    return &WindowsFactory{}
}

func (w *WindowsFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (w *WindowsFactory) CreateCheckbox() Checkbox {
    return &WindowsCheckbox{}
}

// macOS GUI组件
type MacOSButton struct{}

func (m *MacOSButton) Render() string {
    return "macOS Button"
}

func (m *MacOSButton) Click() string {
    return "macOS Button Clicked"
}

type MacOSCheckbox struct{}

func (m *MacOSCheckbox) Render() string {
    return "macOS Checkbox"
}

func (m *MacOSCheckbox) Check() string {
    return "macOS Checkbox Checked"
}

type MacOSFactory struct{}

func NewMacOSFactory() *MacOSFactory {
    return &MacOSFactory{}
}

func (m *MacOSFactory) CreateButton() Button {
    return &MacOSButton{}
}

func (m *MacOSFactory) CreateCheckbox() Checkbox {
    return &MacOSCheckbox{}
}

// GUI应用
type Application struct {
    factory GUIFactory
}

func NewApplication(factory GUIFactory) *Application {
    return &Application{
        factory: factory,
    }
}

func (a *Application) CreateUI() {
    button := a.factory.CreateButton()
    checkbox := a.factory.CreateCheckbox()
    
    fmt.Println(button.Render())
    fmt.Println(checkbox.Render())
}
```

### 6.2 数据库抽象工厂

```go
// 数据库抽象工厂
type Connection interface {
    Connect() error
    Disconnect() error
    Execute(query string) (interface{}, error)
}

type Query interface {
    Build() string
    Execute(conn Connection) (interface{}, error)
}

type DatabaseFactory interface {
    CreateConnection() Connection
    CreateQuery() Query
}

// MySQL实现
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

type MySQLQuery struct {
    sql string
}

func NewMySQLQuery(sql string) *MySQLQuery {
    return &MySQLQuery{
        sql: sql,
    }
}

func (m *MySQLQuery) Build() string {
    return fmt.Sprintf("MySQL: %s", m.sql)
}

func (m *MySQLQuery) Execute(conn Connection) (interface{}, error) {
    return conn.Execute(m.sql)
}

type MySQLFactory struct {
    config map[string]interface{}
}

func NewMySQLFactory(config map[string]interface{}) *MySQLFactory {
    return &MySQLFactory{
        config: config,
    }
}

func (m *MySQLFactory) CreateConnection() Connection {
    return NewMySQLConnection(m.config)
}

func (m *MySQLFactory) CreateQuery() Query {
    return NewMySQLQuery("SELECT * FROM table")
}

// PostgreSQL实现
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

type PostgreSQLQuery struct {
    sql string
}

func NewPostgreSQLQuery(sql string) *PostgreSQLQuery {
    return &PostgreSQLQuery{
        sql: sql,
    }
}

func (p *PostgreSQLQuery) Build() string {
    return fmt.Sprintf("PostgreSQL: %s", p.sql)
}

func (p *PostgreSQLQuery) Execute(conn Connection) (interface{}, error) {
    return conn.Execute(p.sql)
}

type PostgreSQLFactory struct {
    config map[string]interface{}
}

func NewPostgreSQLFactory(config map[string]interface{}) *PostgreSQLFactory {
    return &PostgreSQLFactory{
        config: config,
    }
}

func (p *PostgreSQLFactory) CreateConnection() Connection {
    return NewPostgreSQLConnection(p.config)
}

func (p *PostgreSQLFactory) CreateQuery() Query {
    return NewPostgreSQLQuery("SELECT * FROM table")
}
```

### 6.3 支付系统抽象工厂

```go
// 支付系统抽象工厂
type PaymentGateway interface {
    ProcessPayment(amount float64, currency string) error
    RefundPayment(transactionID string) error
}

type PaymentValidator interface {
    ValidateCard(cardNumber string) bool
    ValidateAmount(amount float64) bool
}

type PaymentFactory interface {
    CreateGateway() PaymentGateway
    CreateValidator() PaymentValidator
}

// Stripe实现
type StripeGateway struct {
    apiKey string
}

func NewStripeGateway(apiKey string) *StripeGateway {
    return &StripeGateway{
        apiKey: apiKey,
    }
}

func (s *StripeGateway) ProcessPayment(amount float64, currency string) error {
    // Stripe支付处理逻辑
    return nil
}

func (s *StripeGateway) RefundPayment(transactionID string) error {
    // Stripe退款逻辑
    return nil
}

type StripeValidator struct{}

func NewStripeValidator() *StripeValidator {
    return &StripeValidator{}
}

func (s *StripeValidator) ValidateCard(cardNumber string) bool {
    // Stripe卡号验证逻辑
    return len(cardNumber) >= 13 && len(cardNumber) <= 19
}

func (s *StripeValidator) ValidateAmount(amount float64) bool {
    // Stripe金额验证逻辑
    return amount > 0 && amount <= 1000000
}

type StripeFactory struct {
    apiKey string
}

func NewStripeFactory(apiKey string) *StripeFactory {
    return &StripeFactory{
        apiKey: apiKey,
    }
}

func (s *StripeFactory) CreateGateway() PaymentGateway {
    return NewStripeGateway(s.apiKey)
}

func (s *StripeFactory) CreateValidator() PaymentValidator {
    return NewStripeValidator()
}

// PayPal实现
type PayPalGateway struct {
    clientID     string
    clientSecret string
}

func NewPayPalGateway(clientID, clientSecret string) *PayPalGateway {
    return &PayPalGateway{
        clientID:     clientID,
        clientSecret: clientSecret,
    }
}

func (p *PayPalGateway) ProcessPayment(amount float64, currency string) error {
    // PayPal支付处理逻辑
    return nil
}

func (p *PayPalGateway) RefundPayment(transactionID string) error {
    // PayPal退款逻辑
    return nil
}

type PayPalValidator struct{}

func NewPayPalValidator() *PayPalValidator {
    return &PayPalValidator{}
}

func (p *PayPalValidator) ValidateCard(cardNumber string) bool {
    // PayPal卡号验证逻辑
    return len(cardNumber) >= 13 && len(cardNumber) <= 19
}

func (p *PayPalValidator) ValidateAmount(amount float64) bool {
    // PayPal金额验证逻辑
    return amount > 0 && amount <= 500000
}

type PayPalFactory struct {
    clientID     string
    clientSecret string
}

func NewPayPalFactory(clientID, clientSecret string) *PayPalFactory {
    return &PayPalFactory{
        clientID:     clientID,
        clientSecret: clientSecret,
    }
}

func (p *PayPalFactory) CreateGateway() PaymentGateway {
    return NewPayPalGateway(p.clientID, p.clientSecret)
}

func (p *PayPalFactory) CreateValidator() PaymentValidator {
    return NewPayPalValidator()
}
```

---

## 7. 相关模式

### 7.1 与工厂方法模式的关系

- **抽象工厂模式**: 创建产品族
- **工厂方法模式**: 创建单个产品

### 7.2 与建造者模式的关系

- **抽象工厂模式**: 创建相关对象
- **建造者模式**: 创建复杂对象

### 7.3 与单例模式的关系

- **抽象工厂模式**: 可以结合单例模式管理工厂实例
- **单例模式**: 确保工厂的唯一性

---

## 总结

抽象工厂模式通过创建产品族，确保了产品之间的兼容性和一致性。它支持开闭原则，易于扩展新的产品族，是面向对象设计中重要的创建型模式。

**关键要点**:

- 创建相关产品族，确保兼容性
- 支持开闭原则，易于扩展
- 隐藏产品创建的复杂性
- 合理选择实现方式（基础、泛型、函数式）

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **抽象工厂模式完成！** 🚀

---

**下一模式**: [04-建造者模式](./04-Builder-Pattern.md)

**返回**: [创建型模式目录](./README.md)
