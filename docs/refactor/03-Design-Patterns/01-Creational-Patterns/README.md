# 创建型设计模式 (Creational Design Patterns)

## 概述

创建型设计模式关注对象的创建过程，提供灵活的对象创建机制，使系统独立于对象创建、组合和表示的方式。本文档采用严格的数学形式化方法，结合最新的 Go 语言特性，对创建型模式进行系统性重构。

## 形式化定义

### 1. 创建型模式的形式化框架

**定义 1.1** (创建型模式)
创建型模式是一个三元组 $\mathcal{C} = (F, P, \phi)$，其中：

- $F$ 是工厂函数集合
- $P$ 是产品类型集合  
- $\phi: F \rightarrow P$ 是工厂到产品的映射函数

**公理 1.1** (创建型模式公理)
对于任意创建型模式 $\mathcal{C} = (F, P, \phi)$：

1. **存在性**: $\forall p \in P, \exists f \in F: \phi(f) = p$
2. **唯一性**: $\forall f_1, f_2 \in F: \phi(f_1) = \phi(f_2) \Rightarrow f_1 = f_2$
3. **封闭性**: $F$ 和 $P$ 在模式操作下封闭

### 2. 模式分类的形式化

**定义 1.2** (模式分类)
设 $\mathcal{P}$ 为所有设计模式的集合，创建型模式集合 $\mathcal{C}$ 满足：
$$\mathcal{C} = \{p \in \mathcal{P} \mid \text{type}(p) = \text{creational}\}$$

**定理 1.1** (模式分类的完备性)
创建型模式集合 $\mathcal{C}$ 是完备的，即：
$$\forall p \in \mathcal{P}: \text{type}(p) = \text{creational} \Leftrightarrow p \in \mathcal{C}$$

## 核心模式

### 1. 单例模式 (Singleton Pattern)

**定义 1.3** (单例模式)
单例模式是一个创建型模式，确保一个类只有一个实例，并提供全局访问点。

**形式化定义**:
设 $T$ 为类型，单例模式定义为函数：
$$S: \text{Unit} \rightarrow T$$
满足：
$$\forall x, y \in \text{Unit}: S(x) = S(y)$$

**Go 语言实现**:

```go
package singleton

import (
    "sync"
    "sync/atomic"
)

// Singleton 单例接口
type Singleton interface {
    GetID() string
    DoSomething() string
}

// singletonImpl 单例实现
type singletonImpl struct {
    id string
}

func (s *singletonImpl) GetID() string {
    return s.id
}

func (s *singletonImpl) DoSomething() string {
    return "Singleton instance: " + s.id
}

// 方法1: 使用 sync.Once (推荐)
type SingletonManager struct {
    instance *singletonImpl
    once     sync.Once
}

func (sm *SingletonManager) GetInstance() Singleton {
    sm.once.Do(func() {
        sm.instance = &singletonImpl{
            id: "singleton-" + generateID(),
        }
    })
    return sm.instance
}

// 方法2: 使用原子操作
var (
    instance     *singletonImpl
    instanceOnce sync.Once
)

func GetInstance() Singleton {
    instanceOnce.Do(func() {
        instance = &singletonImpl{
            id: "atomic-singleton-" + generateID(),
        }
    })
    return instance
}

// 方法3: 包级单例
var globalInstance = &singletonImpl{
    id: "global-singleton-" + generateID(),
}

func GetGlobalInstance() Singleton {
    return globalInstance
}

// 辅助函数
func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 线程安全的单例测试
func TestSingletonThreadSafety() {
    const numGoroutines = 1000
    var wg sync.WaitGroup
    instances := make([]Singleton, numGoroutines)
    
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
            panic("Singleton pattern violated!")
        }
    }
}
```

**形式化证明**:

**定理 1.2** (单例模式的唯一性)
对于任意单例模式实现，满足：
$$\forall x, y: \text{GetInstance}(x) = \text{GetInstance}(y)$$

**证明**:

1. 设 $x, y$ 为任意调用
2. 由于 `sync.Once` 的保证，`GetInstance` 最多执行一次初始化
3. 因此 $\text{GetInstance}(x) = \text{GetInstance}(y)$
4. 证毕。

### 2. 工厂方法模式 (Factory Method Pattern)

**定义 1.4** (工厂方法模式)
工厂方法模式定义一个用于创建对象的接口，让子类决定实例化哪一个类。

**形式化定义**:
设 $P$ 为产品类型集合，$C$ 为创建者类型集合，工厂方法模式定义为：
$$F: C \times \text{Config} \rightarrow P$$
其中 $\text{Config}$ 是配置参数集合。

**Go 语言实现**:

```go
package factory

import (
    "fmt"
    "time"
)

// Product 产品接口
type Product interface {
    GetName() string
    GetPrice() float64
    GetDescription() string
}

// ConcreteProductA 具体产品A
type ConcreteProductA struct {
    name        string
    price       float64
    description string
}

func (p *ConcreteProductA) GetName() string {
    return p.name
}

func (p *ConcreteProductA) GetPrice() float64 {
    return p.price
}

func (p *ConcreteProductA) GetDescription() string {
    return p.description
}

// ConcreteProductB 具体产品B
type ConcreteProductB struct {
    name        string
    price       float64
    description string
    features    []string
}

func (p *ConcreteProductB) GetName() string {
    return p.name
}

func (p *ConcreteProductB) GetPrice() float64 {
    return p.price
}

func (p *ConcreteProductB) GetDescription() string {
    return p.description
}

func (p *ConcreteProductB) GetFeatures() []string {
    return p.features
}

// Creator 创建者接口
type Creator interface {
    CreateProduct(config ProductConfig) Product
}

// ProductConfig 产品配置
type ProductConfig struct {
    Name        string
    Price       float64
    Description string
    Features    []string
}

// ConcreteCreatorA 具体创建者A
type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) CreateProduct(config ProductConfig) Product {
    return &ConcreteProductA{
        name:        config.Name,
        price:       config.Price,
        description: config.Description,
    }
}

// ConcreteCreatorB 具体创建者B
type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) CreateProduct(config ProductConfig) Product {
    return &ConcreteProductB{
        name:        config.Name,
        price:       config.Price,
        description: config.Description,
        features:    config.Features,
    }
}

// 泛型工厂方法
type GenericCreator[T Product] interface {
    Create(config ProductConfig) T
}

// 函数式工厂方法
type ProductFactory func(config ProductConfig) Product

func NewProductFactory(productType string) ProductFactory {
    switch productType {
    case "A":
        return func(config ProductConfig) Product {
            return &ConcreteProductA{
                name:        config.Name,
                price:       config.Price,
                description: config.Description,
            }
        }
    case "B":
        return func(config ProductConfig) Product {
            return &ConcreteProductB{
                name:        config.Name,
                price:       config.Price,
                description: config.Description,
                features:    config.Features,
            }
        }
    default:
        return func(config ProductConfig) Product {
            return &ConcreteProductA{
                name:        "Default Product",
                price:       0.0,
                description: "Default description",
            }
        }
    }
}

// 工厂方法使用示例
func ExampleFactoryMethod() {
    // 使用具体创建者
    creatorA := &ConcreteCreatorA{}
    creatorB := &ConcreteCreatorB{}
    
    configA := ProductConfig{
        Name:        "Product A",
        Price:       100.0,
        Description: "This is product A",
    }
    
    configB := ProductConfig{
        Name:        "Product B",
        Price:       200.0,
        Description: "This is product B",
        Features:    []string{"Feature 1", "Feature 2"},
    }
    
    productA := creatorA.CreateProduct(configA)
    productB := creatorB.CreateProduct(configB)
    
    fmt.Printf("Product A: %s, Price: %.2f\n", productA.GetName(), productA.GetPrice())
    fmt.Printf("Product B: %s, Price: %.2f\n", productB.GetName(), productB.GetPrice())
    
    // 使用函数式工厂
    factoryA := NewProductFactory("A")
    factoryB := NewProductFactory("B")
    
    productA2 := factoryA(configA)
    productB2 := factoryB(configB)
    
    fmt.Printf("Factory Product A: %s\n", productA2.GetName())
    fmt.Printf("Factory Product B: %s\n", productB2.GetName())
}
```

**形式化证明**:

**定理 1.3** (工厂方法的可扩展性)
工厂方法模式支持开闭原则：
$$\forall c \in C, \forall p \in P: \exists f \in F: f(c) = p$$

**证明**:

1. 设 $c$ 为任意创建者，$p$ 为任意产品
2. 通过实现新的具体创建者，可以创建新的产品类型
3. 无需修改现有代码，满足开闭原则
4. 证毕。

### 3. 抽象工厂模式 (Abstract Factory Pattern)

**定义 1.5** (抽象工厂模式)
抽象工厂模式提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们的具体类。

**形式化定义**:
设 $\mathcal{P} = \{P_1, P_2, \ldots, P_n\}$ 为产品族集合，抽象工厂模式定义为：
$$F: \text{Family} \rightarrow \prod_{i=1}^{n} P_i$$
其中 $\text{Family}$ 是产品族标识符集合。

**Go 语言实现**:

```go
package abstractfactory

import (
    "fmt"
)

// 产品族A的接口
type ProductA interface {
    GetName() string
    GetType() string
}

type ProductB interface {
    GetName() string
    GetCategory() string
}

// 具体产品族A的实现
type ConcreteProductA1 struct {
    name string
}

func (p *ConcreteProductA1) GetName() string {
    return p.name
}

func (p *ConcreteProductA1) GetType() string {
    return "Type A1"
}

type ConcreteProductB1 struct {
    name     string
    category string
}

func (p *ConcreteProductB1) GetName() string {
    return p.name
}

func (p *ConcreteProductB1) GetCategory() string {
    return p.category
}

// 具体产品族B的实现
type ConcreteProductA2 struct {
    name string
}

func (p *ConcreteProductA2) GetName() string {
    return p.name
}

func (p *ConcreteProductA2) GetType() string {
    return "Type A2"
}

type ConcreteProductB2 struct {
    name     string
    category string
}

func (p *ConcreteProductB2) GetName() string {
    return p.name
}

func (p *ConcreteProductB2) GetCategory() string {
    return p.category
}

// 抽象工厂接口
type AbstractFactory interface {
    CreateProductA() ProductA
    CreateProductB() ProductB
}

// 具体工厂1
type ConcreteFactory1 struct{}

func (f *ConcreteFactory1) CreateProductA() ProductA {
    return &ConcreteProductA1{
        name: "Product A1",
    }
}

func (f *ConcreteFactory1) CreateProductB() ProductB {
    return &ConcreteProductB1{
        name:     "Product B1",
        category: "Category 1",
    }
}

// 具体工厂2
type ConcreteFactory2 struct{}

func (f *ConcreteFactory2) CreateProductA() ProductA {
    return &ConcreteProductA2{
        name: "Product A2",
    }
}

func (f *ConcreteFactory2) CreateProductB() ProductB {
    return &ConcreteProductB2{
        name:     "Product B2",
        category: "Category 2",
    }
}

// 泛型抽象工厂
type GenericAbstractFactory[T1 ProductA, T2 ProductB] interface {
    CreateProductA() T1
    CreateProductB() T2
}

// 函数式抽象工厂
type ProductAFactory func() ProductA
type ProductBFactory func() ProductB

type FunctionalAbstractFactory struct {
    createA ProductAFactory
    createB ProductBFactory
}

func (f *FunctionalAbstractFactory) CreateProductA() ProductA {
    return f.createA()
}

func (f *FunctionalAbstractFactory) CreateProductB() ProductB {
    return f.createB()
}

// 工厂注册表
type FactoryRegistry struct {
    factories map[string]AbstractFactory
}

func NewFactoryRegistry() *FactoryRegistry {
    return &FactoryRegistry{
        factories: make(map[string]AbstractFactory),
    }
}

func (r *FactoryRegistry) Register(name string, factory AbstractFactory) {
    r.factories[name] = factory
}

func (r *FactoryRegistry) GetFactory(name string) (AbstractFactory, error) {
    factory, exists := r.factories[name]
    if !exists {
        return nil, fmt.Errorf("factory %s not found", name)
    }
    return factory, nil
}

// 使用示例
func ExampleAbstractFactory() {
    // 使用具体工厂
    factory1 := &ConcreteFactory1{}
    factory2 := &ConcreteFactory2{}
    
    productA1 := factory1.CreateProductA()
    productB1 := factory1.CreateProductB()
    
    productA2 := factory2.CreateProductA()
    productB2 := factory2.CreateProductB()
    
    fmt.Printf("Factory 1 - Product A: %s (%s)\n", productA1.GetName(), productA1.GetType())
    fmt.Printf("Factory 1 - Product B: %s (%s)\n", productB1.GetName(), productB1.GetCategory())
    fmt.Printf("Factory 2 - Product A: %s (%s)\n", productA2.GetName(), productA2.GetType())
    fmt.Printf("Factory 2 - Product B: %s (%s)\n", productB2.GetName(), productB2.GetCategory())
    
    // 使用工厂注册表
    registry := NewFactoryRegistry()
    registry.Register("factory1", factory1)
    registry.Register("factory2", factory2)
    
    if f1, err := registry.GetFactory("factory1"); err == nil {
        pa1 := f1.CreateProductA()
        pb1 := f1.CreateProductB()
        fmt.Printf("Registry Factory 1 - Product A: %s\n", pa1.GetName())
        fmt.Printf("Registry Factory 1 - Product B: %s\n", pb1.GetName())
    }
}
```

### 4. 建造者模式 (Builder Pattern)

**定义 1.6** (建造者模式)
建造者模式将一个复杂对象的构建与其表示分离，允许使用相同的构建过程创建不同的表示。

**形式化定义**:
设 $P$ 为产品类型，$B$ 为建造者类型，建造者模式定义为：
$$B \times \text{Config}_1 \times \text{Config}_2 \times \ldots \times \text{Config}_n \rightarrow P$$

**Go 语言实现**:

```go
package builder

import (
    "fmt"
    "strings"
)

// Product 产品
type Product struct {
    parts []string
}

func (p *Product) AddPart(part string) {
    p.parts = append(p.parts, part)
}

func (p *Product) Show() string {
    return strings.Join(p.parts, " + ")
}

// Builder 建造者接口
type Builder interface {
    BuildPartA()
    BuildPartB()
    BuildPartC()
    GetResult() *Product
}

// ConcreteBuilder1 具体建造者1
type ConcreteBuilder1 struct {
    product *Product
}

func NewConcreteBuilder1() *ConcreteBuilder1 {
    return &ConcreteBuilder1{
        product: &Product{},
    }
}

func (b *ConcreteBuilder1) BuildPartA() {
    b.product.AddPart("Part A1")
}

func (b *ConcreteBuilder1) BuildPartB() {
    b.product.AddPart("Part B1")
}

func (b *ConcreteBuilder1) BuildPartC() {
    b.product.AddPart("Part C1")
}

func (b *ConcreteBuilder1) GetResult() *Product {
    return b.product
}

// ConcreteBuilder2 具体建造者2
type ConcreteBuilder2 struct {
    product *Product
}

func NewConcreteBuilder2() *ConcreteBuilder2 {
    return &ConcreteBuilder2{
        product: &Product{},
    }
}

func (b *ConcreteBuilder2) BuildPartA() {
    b.product.AddPart("Part A2")
}

func (b *ConcreteBuilder2) BuildPartB() {
    b.product.AddPart("Part B2")
}

func (b *ConcreteBuilder2) BuildPartC() {
    b.product.AddPart("Part C2")
}

func (b *ConcreteBuilder2) GetResult() *Product {
    return b.product
}

// Director 指挥者
type Director struct {
    builder Builder
}

func NewDirector(builder Builder) *Director {
    return &Director{
        builder: builder,
    }
}

func (d *Director) Construct() *Product {
    d.builder.BuildPartA()
    d.builder.BuildPartB()
    d.builder.BuildPartC()
    return d.builder.GetResult()
}

// 函数式建造者
type ProductBuilder struct {
    product *Product
}

func NewProductBuilder() *ProductBuilder {
    return &ProductBuilder{
        product: &Product{},
    }
}

func (pb *ProductBuilder) WithPartA(part string) *ProductBuilder {
    pb.product.AddPart("Part A: " + part)
    return pb
}

func (pb *ProductBuilder) WithPartB(part string) *ProductBuilder {
    pb.product.AddPart("Part B: " + part)
    return pb
}

func (pb *ProductBuilder) WithPartC(part string) *ProductBuilder {
    pb.product.AddPart("Part C: " + part)
    return pb
}

func (pb *ProductBuilder) Build() *Product {
    return pb.product
}

// 泛型建造者
type GenericBuilder[T any] interface {
    Build() T
}

// 使用示例
func ExampleBuilder() {
    // 使用传统建造者模式
    builder1 := NewConcreteBuilder1()
    director1 := NewDirector(builder1)
    product1 := director1.Construct()
    fmt.Printf("Product 1: %s\n", product1.Show())
    
    builder2 := NewConcreteBuilder2()
    director2 := NewDirector(builder2)
    product2 := director2.Construct()
    fmt.Printf("Product 2: %s\n", product2.Show())
    
    // 使用函数式建造者
    product3 := NewProductBuilder().
        WithPartA("Custom A").
        WithPartB("Custom B").
        WithPartC("Custom C").
        Build()
    fmt.Printf("Product 3: %s\n", product3.Show())
}
```

### 5. 原型模式 (Prototype Pattern)

**定义 1.7** (原型模式)
原型模式通过复制现有的实例创建新的实例，而不是通过新建。

**形式化定义**:
设 $P$ 为原型类型，原型模式定义为：
$$\text{Clone}: P \rightarrow P$$
满足：
$$\forall p \in P: \text{Clone}(p) \neq p \land \text{DeepEqual}(\text{Clone}(p), p)$$

**Go 语言实现**:

```go
package prototype

import (
    "fmt"
    "time"
)

// Prototype 原型接口
type Prototype interface {
    Clone() Prototype
    GetName() string
    SetName(name string)
}

// ConcretePrototype 具体原型
type ConcretePrototype struct {
    name    string
    data    []int
    created time.Time
}

func NewConcretePrototype(name string) *ConcretePrototype {
    return &ConcretePrototype{
        name:    name,
        data:    make([]int, 0),
        created: time.Now(),
    }
}

func (p *ConcretePrototype) Clone() Prototype {
    // 深拷贝
    cloned := &ConcretePrototype{
        name:    p.name,
        data:    make([]int, len(p.data)),
        created: time.Now(), // 新创建时间
    }
    copy(cloned.data, p.data)
    return cloned
}

func (p *ConcretePrototype) GetName() string {
    return p.name
}

func (p *ConcretePrototype) SetName(name string) {
    p.name = name
}

func (p *ConcretePrototype) AddData(value int) {
    p.data = append(p.data, value)
}

func (p *ConcretePrototype) GetData() []int {
    return p.data
}

// 原型注册表
type PrototypeRegistry struct {
    prototypes map[string]Prototype
}

func NewPrototypeRegistry() *PrototypeRegistry {
    return &PrototypeRegistry{
        prototypes: make(map[string]Prototype),
    }
}

func (r *PrototypeRegistry) Register(name string, prototype Prototype) {
    r.prototypes[name] = prototype
}

func (r *PrototypeRegistry) Clone(name string) (Prototype, error) {
    prototype, exists := r.prototypes[name]
    if !exists {
        return nil, fmt.Errorf("prototype %s not found", name)
    }
    return prototype.Clone(), nil
}

// 使用示例
func ExamplePrototype() {
    // 创建原型
    original := NewConcretePrototype("Original")
    original.AddData(1)
    original.AddData(2)
    original.AddData(3)
    
    // 克隆原型
    cloned := original.Clone().(*ConcretePrototype)
    cloned.SetName("Cloned")
    cloned.AddData(4)
    
    fmt.Printf("Original: %s, Data: %v\n", original.GetName(), original.GetData())
    fmt.Printf("Cloned: %s, Data: %v\n", cloned.GetName(), cloned.GetData())
    
    // 使用原型注册表
    registry := NewPrototypeRegistry()
    registry.Register("template", original)
    
    if cloned2, err := registry.Clone("template"); err == nil {
        cloned2.SetName("From Registry")
        fmt.Printf("From Registry: %s\n", cloned2.GetName())
    }
}
```

## 模式关系分析

### 1. 模式间的数学关系

**定义 1.8** (模式关系)
设 $\mathcal{P}$ 为所有创建型模式的集合，模式关系 $R \subseteq \mathcal{P} \times \mathcal{P}$ 定义为：
$$(p_1, p_2) \in R \Leftrightarrow p_1 \text{ 可以组合或替代 } p_2$$

**定理 1.4** (模式组合的传递性)
模式组合关系是传递的：
$$\forall p_1, p_2, p_3 \in \mathcal{P}: (p_1, p_2) \in R \land (p_2, p_3) \in R \Rightarrow (p_1, p_3) \in R$$

### 2. 模式等价性

**定义 1.9** (模式等价)
两个创建型模式 $p_1, p_2$ 等价，记作 $p_1 \equiv p_2$，当且仅当：
$$\forall \text{context}: \text{result}(p_1, \text{context}) = \text{result}(p_2, \text{context})$$

## 性能分析

### 1. 时间复杂度分析

| 模式 | 创建时间复杂度 | 空间复杂度 | 适用场景 |
|------|----------------|------------|----------|
| 单例 | O(1) | O(1) | 全局状态管理 |
| 工厂方法 | O(1) | O(1) | 对象创建封装 |
| 抽象工厂 | O(n) | O(n) | 产品族创建 |
| 建造者 | O(n) | O(n) | 复杂对象构建 |
| 原型 | O(n) | O(n) | 对象复制 |

### 2. 内存使用分析

```go
// 内存使用分析工具
type MemoryAnalyzer struct{}

func (ma *MemoryAnalyzer) AnalyzePattern(pattern string, iterations int) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    before := m.Alloc
    
    // 执行模式
    for i := 0; i < iterations; i++ {
        switch pattern {
        case "singleton":
            _ = GetInstance()
        case "factory":
            factory := NewProductFactory("A")
            _ = factory(ProductConfig{})
        case "builder":
            _ = NewProductBuilder().WithPartA("test").Build()
        case "prototype":
            proto := NewConcretePrototype("test")
            _ = proto.Clone()
        }
    }
    
    runtime.ReadMemStats(&m)
    after := m.Alloc
    
    fmt.Printf("Pattern: %s, Memory delta: %d bytes\n", pattern, after-before)
}
```

## 最佳实践

### 1. 模式选择指南

**规则 1.1** (单例模式选择)
当且仅当满足以下条件时使用单例模式：

1. 全局状态管理需求
2. 资源池管理
3. 配置管理

**规则 1.2** (工厂模式选择)
当且仅当满足以下条件时使用工厂模式：

1. 对象创建逻辑复杂
2. 需要支持多种产品类型
3. 需要延迟对象创建

**规则 1.3** (建造者模式选择)
当且仅当满足以下条件时使用建造者模式：

1. 对象构建过程复杂
2. 需要支持不同的构建顺序
3. 需要构建不可变对象

### 2. 性能优化建议

1. **缓存机制**: 在工厂模式中使用对象池
2. **延迟初始化**: 在单例模式中使用懒加载
3. **内存复用**: 在原型模式中使用浅拷贝
4. **并发安全**: 使用原子操作和锁机制

## 持续构建状态

- [x] 单例模式 (100%)
- [x] 工厂方法模式 (100%)
- [x] 抽象工厂模式 (100%)
- [x] 建造者模式 (100%)
- [x] 原型模式 (100%)
- [ ] 对象池模式 (0%)
- [ ] 多例模式 (0%)

---

**构建原则**: 严格数学规范，形式化证明，Go语言实现！<(￣︶￣)↗[GO!]
