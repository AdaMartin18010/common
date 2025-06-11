# Golang 设计模式详解

## 概述

设计模式是软件开发中常见问题的典型解决方案。在Golang中，设计模式的应用有其独特的特点，充分利用了Go语言的简洁性和并发特性。

## 设计模式分类

### 1. 创建型模式 (Creational Patterns)

#### 单例模式 (Singleton Pattern)

```go
package main

import (
    "fmt"
    "sync"
)

// 线程安全的单例模式
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "Initialized"}
        fmt.Println("Singleton initialized")
    })
    return instance
}

func (s *Singleton) GetData() string {
    return s.data
}

func (s *Singleton) SetData(data string) {
    s.data = data
}

func main() {
    // 多个goroutine同时获取实例
    for i := 0; i < 5; i++ {
        go func(id int) {
            instance := GetInstance()
            fmt.Printf("Goroutine %d got instance: %v\n", id, instance.GetData())
        }(i)
    }
    
    // 等待所有goroutine完成
    time.Sleep(1 * time.Second)
}
```

#### 工厂模式 (Factory Pattern)

```go
package main

import "fmt"

// 产品接口
type Product interface {
    Use() string
}

// 具体产品
type ConcreteProductA struct{}
type ConcreteProductB struct{}

func (p *ConcreteProductA) Use() string {
    return "Using Product A"
}

func (p *ConcreteProductB) Use() string {
    return "Using Product B"
}

// 工厂接口
type Factory interface {
    CreateProduct() Product
}

// 具体工厂
type ConcreteFactoryA struct{}
type ConcreteFactoryB struct{}

func (f *ConcreteFactoryA) CreateProduct() Product {
    return &ConcreteProductA{}
}

func (f *ConcreteFactoryB) CreateProduct() Product {
    return &ConcreteProductB{}
}

// 工厂方法
func NewFactory(factoryType string) Factory {
    switch factoryType {
    case "A":
        return &ConcreteFactoryA{}
    case "B":
        return &ConcreteFactoryB{}
    default:
        return nil
    }
}

func main() {
    factoryA := NewFactory("A")
    productA := factoryA.CreateProduct()
    fmt.Println(productA.Use())
    
    factoryB := NewFactory("B")
    productB := factoryB.CreateProduct()
    fmt.Println(productB.Use())
}
```

#### 建造者模式 (Builder Pattern)

```go
package main

import "fmt"

// 产品
type Computer struct {
    CPU    string
    Memory string
    Disk   string
    GPU    string
}

func (c *Computer) String() string {
    return fmt.Sprintf("Computer{CPU: %s, Memory: %s, Disk: %s, GPU: %s}", 
        c.CPU, c.Memory, c.Disk, c.GPU)
}

// 建造者接口
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetMemory(memory string) ComputerBuilder
    SetDisk(disk string) ComputerBuilder
    SetGPU(gpu string) ComputerBuilder
    Build() *Computer
}

// 具体建造者
type GamingComputerBuilder struct {
    computer *Computer
}

func NewGamingComputerBuilder() *GamingComputerBuilder {
    return &GamingComputerBuilder{
        computer: &Computer{},
    }
}

func (b *GamingComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    b.computer.CPU = cpu
    return b
}

func (b *GamingComputerBuilder) SetMemory(memory string) ComputerBuilder {
    b.computer.Memory = memory
    return b
}

func (b *GamingComputerBuilder) SetDisk(disk string) ComputerBuilder {
    b.computer.Disk = disk
    return b
}

func (b *GamingComputerBuilder) SetGPU(gpu string) ComputerBuilder {
    b.computer.GPU = gpu
    return b
}

func (b *GamingComputerBuilder) Build() *Computer {
    return b.computer
}

// 导演
type ComputerDirector struct {
    builder ComputerBuilder
}

func NewComputerDirector(builder ComputerBuilder) *ComputerDirector {
    return &ComputerDirector{builder: builder}
}

func (d *ComputerDirector) ConstructGamingComputer() *Computer {
    return d.builder.
        SetCPU("Intel i9-12900K").
        SetMemory("32GB DDR5").
        SetDisk("2TB NVMe SSD").
        SetGPU("RTX 4090").
        Build()
}

func main() {
    builder := NewGamingComputerBuilder()
    director := NewComputerDirector(builder)
    
    gamingComputer := director.ConstructGamingComputer()
    fmt.Printf("Built: %s\n", gamingComputer)
}
```

### 2. 结构型模式 (Structural Patterns)

#### 适配器模式 (Adapter Pattern)

```go
package main

import "fmt"

// 目标接口
type Target interface {
    Request() string
}

// 被适配的类
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
    return "Specific request from Adaptee"
}

// 适配器
type Adapter struct {
    adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
    return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}

func main() {
    adaptee := &Adaptee{}
    adapter := NewAdapter(adaptee)
    
    fmt.Println(adapter.Request())
}
```

#### 装饰器模式 (Decorator Pattern)

```go
package main

import "fmt"

// 组件接口
type Coffee interface {
    Cost() float64
    Description() string
}

// 具体组件
type SimpleCoffee struct{}

func (c *SimpleCoffee) Cost() float64 {
    return 2.0
}

func (c *SimpleCoffee) Description() string {
    return "Simple coffee"
}

// 装饰器基类
type CoffeeDecorator struct {
    coffee Coffee
}

func (d *CoffeeDecorator) Cost() float64 {
    return d.coffee.Cost()
}

func (d *CoffeeDecorator) Description() string {
    return d.coffee.Description()
}

// 具体装饰器
type MilkDecorator struct {
    *CoffeeDecorator
}

func NewMilkDecorator(coffee Coffee) *MilkDecorator {
    return &MilkDecorator{
        CoffeeDecorator: &CoffeeDecorator{coffee: coffee},
    }
}

func (d *MilkDecorator) Cost() float64 {
    return d.coffee.Cost() + 0.5
}

func (d *MilkDecorator) Description() string {
    return d.coffee.Description() + ", milk"
}

type SugarDecorator struct {
    *CoffeeDecorator
}

func NewSugarDecorator(coffee Coffee) *SugarDecorator {
    return &SugarDecorator{
        CoffeeDecorator: &CoffeeDecorator{coffee: coffee},
    }
}

func (d *SugarDecorator) Cost() float64 {
    return d.coffee.Cost() + 0.2
}

func (d *SugarDecorator) Description() string {
    return d.coffee.Description() + ", sugar"
}

func main() {
    coffee := &SimpleCoffee{}
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffee.Cost(), coffee.Description())
    
    // 添加牛奶
    coffeeWithMilk := NewMilkDecorator(coffee)
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffeeWithMilk.Cost(), coffeeWithMilk.Description())
    
    // 添加糖
    coffeeWithMilkAndSugar := NewSugarDecorator(coffeeWithMilk)
    fmt.Printf("Cost: $%.2f, Description: %s\n", coffeeWithMilkAndSugar.Cost(), coffeeWithMilkAndSugar.Description())
}
```

#### 代理模式 (Proxy Pattern)

```go
package main

import (
    "fmt"
    "time"
)

// 主题接口
type Subject interface {
    Request() string
}

// 真实主题
type RealSubject struct{}

func (r *RealSubject) Request() string {
    // 模拟耗时操作
    time.Sleep(2 * time.Second)
    return "Real subject response"
}

// 代理
type Proxy struct {
    realSubject *RealSubject
    cache       string
    cached      bool
}

func NewProxy() *Proxy {
    return &Proxy{
        realSubject: &RealSubject{},
    }
}

func (p *Proxy) Request() string {
    if !p.cached {
        fmt.Println("Proxy: Caching real subject response")
        p.cache = p.realSubject.Request()
        p.cached = true
    } else {
        fmt.Println("Proxy: Returning cached response")
    }
    return p.cache
}

func main() {
    proxy := NewProxy()
    
    // 第一次请求，会调用真实主题
    fmt.Println("First request:")
    result := proxy.Request()
    fmt.Printf("Result: %s\n", result)
    
    // 第二次请求，使用缓存
    fmt.Println("\nSecond request:")
    result = proxy.Request()
    fmt.Printf("Result: %s\n", result)
}
```

### 3. 行为型模式 (Behavioral Patterns)

#### 策略模式 (Strategy Pattern)

```go
package main

import "fmt"

// 策略接口
type PaymentStrategy interface {
    Pay(amount float64) string
}

// 具体策略
type CreditCardPayment struct{}

func (c *CreditCardPayment) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f using Credit Card", amount)
}

type PayPalPayment struct{}

func (p *PayPalPayment) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f using PayPal", amount)
}

type CryptoPayment struct{}

func (c *CryptoPayment) Pay(amount float64) string {
    return fmt.Sprintf("Paid $%.2f using Cryptocurrency", amount)
}

// 上下文
type ShoppingCart struct {
    paymentStrategy PaymentStrategy
}

func NewShoppingCart() *ShoppingCart {
    return &ShoppingCart{}
}

func (c *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
    c.paymentStrategy = strategy
}

func (c *ShoppingCart) Checkout(amount float64) string {
    if c.paymentStrategy == nil {
        return "No payment strategy set"
    }
    return c.paymentStrategy.Pay(amount)
}

func main() {
    cart := NewShoppingCart()
    
    // 使用信用卡支付
    cart.SetPaymentStrategy(&CreditCardPayment{})
    fmt.Println(cart.Checkout(100.50))
    
    // 使用PayPal支付
    cart.SetPaymentStrategy(&PayPalPayment{})
    fmt.Println(cart.Checkout(75.25))
    
    // 使用加密货币支付
    cart.SetPaymentStrategy(&CryptoPayment{})
    fmt.Println(cart.Checkout(200.00))
}
```

#### 观察者模式 (Observer Pattern)

```go
package main

import (
    "fmt"
    "sync"
)

// 观察者接口
type Observer interface {
    Update(data string)
}

// 主题接口
type Subject interface {
    Attach(observer Observer)
    Detach(observer Observer)
    Notify()
}

// 具体主题
type NewsAgency struct {
    observers []Observer
    news      string
    mu        sync.RWMutex
}

func NewNewsAgency() *NewsAgency {
    return &NewsAgency{
        observers: make([]Observer, 0),
    }
}

func (na *NewsAgency) Attach(observer Observer) {
    na.mu.Lock()
    defer na.mu.Unlock()
    na.observers = append(na.observers, observer)
}

func (na *NewsAgency) Detach(observer Observer) {
    na.mu.Lock()
    defer na.mu.Unlock()
    
    for i, obs := range na.observers {
        if obs == observer {
            na.observers = append(na.observers[:i], na.observers[i+1:]...)
            break
        }
    }
}

func (na *NewsAgency) Notify() {
    na.mu.RLock()
    defer na.mu.RUnlock()
    
    for _, observer := range na.observers {
        observer.Update(na.news)
    }
}

func (na *NewsAgency) SetNews(news string) {
    na.mu.Lock()
    na.news = news
    na.mu.Unlock()
    na.Notify()
}

// 具体观察者
type NewsChannel struct {
    name string
}

func NewNewsChannel(name string) *NewsChannel {
    return &NewsChannel{name: name}
}

func (nc *NewsChannel) Update(news string) {
    fmt.Printf("%s received news: %s\n", nc.name, news)
}

func main() {
    agency := NewNewsAgency()
    
    // 创建观察者
    channel1 := NewNewsChannel("CNN")
    channel2 := NewNewsChannel("BBC")
    channel3 := NewNewsChannel("Al Jazeera")
    
    // 注册观察者
    agency.Attach(channel1)
    agency.Attach(channel2)
    agency.Attach(channel3)
    
    // 发布新闻
    agency.SetNews("Breaking: Go 1.21 released!")
    
    // 移除一个观察者
    agency.Detach(channel2)
    
    // 发布另一条新闻
    agency.SetNews("Update: New features in Go 1.21")
}
```

#### 命令模式 (Command Pattern)

```go
package main

import "fmt"

// 命令接口
type Command interface {
    Execute()
    Undo()
}

// 接收者
type Light struct {
    location string
}

func NewLight(location string) *Light {
    return &Light{location: location}
}

func (l *Light) TurnOn() {
    fmt.Printf("%s light is ON\n", l.location)
}

func (l *Light) TurnOff() {
    fmt.Printf("%s light is OFF\n", l.location)
}

// 具体命令
type LightOnCommand struct {
    light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
    return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
    c.light.TurnOn()
}

func (c *LightOnCommand) Undo() {
    c.light.TurnOff()
}

type LightOffCommand struct {
    light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
    return &LightOffCommand{light: light}
}

func (c *LightOffCommand) Execute() {
    c.light.TurnOff()
}

func (c *LightOffCommand) Undo() {
    c.light.TurnOn()
}

// 调用者
type RemoteControl struct {
    onCommands  []Command
    offCommands []Command
    undoCommand Command
}

func NewRemoteControl() *RemoteControl {
    return &RemoteControl{
        onCommands:  make([]Command, 7),
        offCommands: make([]Command, 7),
    }
}

func (rc *RemoteControl) SetCommand(slot int, onCommand, offCommand Command) {
    rc.onCommands[slot] = onCommand
    rc.offCommands[slot] = offCommand
}

func (rc *RemoteControl) OnButtonWasPushed(slot int) {
    if rc.onCommands[slot] != nil {
        rc.onCommands[slot].Execute()
        rc.undoCommand = rc.onCommands[slot]
    }
}

func (rc *RemoteControl) OffButtonWasPushed(slot int) {
    if rc.offCommands[slot] != nil {
        rc.offCommands[slot].Execute()
        rc.undoCommand = rc.offCommands[slot]
    }
}

func (rc *RemoteControl) UndoButtonWasPushed() {
    if rc.undoCommand != nil {
        rc.undoCommand.Undo()
    }
}

func main() {
    remote := NewRemoteControl()
    
    // 创建接收者
    livingRoomLight := NewLight("Living Room")
    kitchenLight := NewLight("Kitchen")
    
    // 创建命令
    livingRoomLightOn := NewLightOnCommand(livingRoomLight)
    livingRoomLightOff := NewLightOffCommand(livingRoomLight)
    kitchenLightOn := NewLightOnCommand(kitchenLight)
    kitchenLightOff := NewLightOffCommand(kitchenLight)
    
    // 设置命令
    remote.SetCommand(0, livingRoomLightOn, livingRoomLightOff)
    remote.SetCommand(1, kitchenLightOn, kitchenLightOff)
    
    // 执行命令
    remote.OnButtonWasPushed(0)
    remote.OnButtonWasPushed(1)
    remote.OffButtonWasPushed(0)
    remote.UndoButtonWasPushed()
}
```

## Go语言特有的设计模式

### 1. 函数式编程模式

```go
package main

import "fmt"

// 高阶函数
type Operation func(int, int) int

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}

func ApplyOperation(op Operation, a, b int) int {
    return op(a, b)
}

// 函数装饰器
func LoggingDecorator(fn func(int, int) int) func(int, int) int {
    return func(a, b int) int {
        fmt.Printf("Calling function with %d and %d\n", a, b)
        result := fn(a, b)
        fmt.Printf("Result: %d\n", result)
        return result
    }
}

func main() {
    // 使用高阶函数
    result := ApplyOperation(Add, 5, 3)
    fmt.Printf("Add result: %d\n", result)
    
    result = ApplyOperation(Multiply, 4, 6)
    fmt.Printf("Multiply result: %d\n", result)
    
    // 使用装饰器
    loggedAdd := LoggingDecorator(Add)
    loggedAdd(10, 20)
}
```

### 2. 通道模式 (Channel Patterns)

```go
package main

import (
    "fmt"
    "time"
)

// 工作池模式
type WorkerPool struct {
    workers int
    tasks   chan func()
    wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    wp := &WorkerPool{
        workers: workers,
        tasks:   make(chan func()),
    }
    
    for i := 0; i < workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
    
    return wp
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    for task := range wp.tasks {
        task()
    }
}

func (wp *WorkerPool) Submit(task func()) {
    wp.tasks <- task
}

func (wp *WorkerPool) Close() {
    close(wp.tasks)
    wp.wg.Wait()
}

// 发布-订阅模式
type Publisher struct {
    subscribers map[string]chan interface{}
    mu          sync.RWMutex
}

func NewPublisher() *Publisher {
    return &Publisher{
        subscribers: make(map[string]chan interface{}),
    }
}

func (p *Publisher) Subscribe(topic string) chan interface{} {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    ch := make(chan interface{}, 10)
    p.subscribers[topic] = ch
    return ch
}

func (p *Publisher) Publish(topic string, data interface{}) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    if ch, exists := p.subscribers[topic]; exists {
        select {
        case ch <- data:
        default:
            // Channel已满，跳过
        }
    }
}

func main() {
    // 工作池示例
    pool := NewWorkerPool(3)
    
    for i := 0; i < 10; i++ {
        taskID := i
        pool.Submit(func() {
            fmt.Printf("Worker processing task %d\n", taskID)
            time.Sleep(100 * time.Millisecond)
        })
    }
    
    pool.Close()
    
    // 发布-订阅示例
    pub := NewPublisher()
    
    // 订阅者
    go func() {
        ch := pub.Subscribe("news")
        for data := range ch {
            fmt.Printf("Received news: %v\n", data)
        }
    }()
    
    // 发布者
    pub.Publish("news", "Breaking news!")
    pub.Publish("news", "Update available")
    
    time.Sleep(1 * time.Second)
}
```

### 3. 中间件模式 (Middleware Pattern)

```go
package main

import (
    "fmt"
    "log"
    "time"
)

// 处理器接口
type Handler interface {
    Handle(request string) string
}

// 基础处理器
type BaseHandler struct{}

func (h *BaseHandler) Handle(request string) string {
    return fmt.Sprintf("Handled: %s", request)
}

// 中间件接口
type Middleware func(Handler) Handler

// 日志中间件
func LoggingMiddleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        log.Printf("Request: %s", request)
        start := time.Now()
        
        response := next.Handle(request)
        
        log.Printf("Response: %s, Duration: %v", response, time.Since(start))
        return response
    })
}

// 认证中间件
func AuthMiddleware(next Handler) Handler {
    return HandlerFunc(func(request string) string {
        // 模拟认证检查
        if request == "unauthorized" {
            return "Unauthorized"
        }
        return next.Handle(request)
    })
}

// 处理器函数类型
type HandlerFunc func(string) string

func (f HandlerFunc) Handle(request string) string {
    return f(request)
}

// 中间件链
func Chain(handler Handler, middlewares ...Middleware) Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

func main() {
    baseHandler := &BaseHandler{}
    
    // 创建中间件链
    handler := Chain(
        baseHandler,
        LoggingMiddleware,
        AuthMiddleware,
    )
    
    // 测试
    fmt.Println(handler.Handle("authorized request"))
    fmt.Println(handler.Handle("unauthorized"))
}
```

## 最佳实践

### 1. 接口设计原则

```go
// 接口应该小而专注
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

### 2. 错误处理模式

```go
// 错误包装
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// 错误处理函数
func HandleError(err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        log.Printf("Application error: %s (code: %d)", appErr.Message, appErr.Code)
    } else {
        log.Printf("System error: %v", err)
    }
}
```

### 3. 配置管理模式

```go
// 配置结构
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type DatabaseConfig struct {
    DSN string `json:"dsn"`
}

type LoggingConfig struct {
    Level string `json:"level"`
}

// 配置加载器
type ConfigLoader interface {
    Load() (*Config, error)
}

type FileConfigLoader struct {
    path string
}

func (f *FileConfigLoader) Load() (*Config, error) {
    // 从文件加载配置
    return &Config{}, nil
}

type EnvConfigLoader struct{}

func (e *EnvConfigLoader) Load() (*Config, error) {
    // 从环境变量加载配置
    return &Config{}, nil
}
```

## 2025年改进

### 1. 泛型设计模式

```go
// 泛型容器
type Container[T any] struct {
    items []T
}

func NewContainer[T any]() *Container[T] {
    return &Container[T]{
        items: make([]T, 0),
    }
}

func (c *Container[T]) Add(item T) {
    c.items = append(c.items, item)
}

func (c *Container[T]) Get(index int) T {
    return c.items[index]
}

func (c *Container[T]) Filter(predicate func(T) bool) *Container[T] {
    result := NewContainer[T]()
    for _, item := range c.items {
        if predicate(item) {
            result.Add(item)
        }
    }
    return result
}

// 泛型工厂
type Factory[T any] interface {
    Create() T
}

type StringFactory struct{}

func (f *StringFactory) Create() string {
    return "default string"
}

type IntFactory struct{}

func (f *IntFactory) Create() int {
    return 42
}
```

### 2. 并发安全模式

```go
// 线程安全的缓存
type Cache[K comparable, V any] struct {
    data map[K]V
    mu   sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        data: make(map[K]V),
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.data, key)
}
```

## 总结

Golang的设计模式应用具有以下特点：

1. **简洁性** - 利用Go语言的简洁语法
2. **并发性** - 充分利用goroutines和channels
3. **接口性** - 基于接口的组合和扩展
4. **函数式** - 支持函数式编程模式
5. **泛型性** - 利用Go 1.18+的泛型特性

通过合理应用这些设计模式，可以构建出高质量、可维护的Go应用程序。

---

*最后更新时间: 2025年1月*
*文档版本: v1.0*
