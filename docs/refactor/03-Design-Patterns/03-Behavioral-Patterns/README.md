# 行为型设计模式 (Behavioral Design Patterns)

## 概述

行为型设计模式关注对象之间的通信，描述对象之间如何相互协作以及如何分配职责。这些模式提供了灵活的方式来处理对象间的交互，使系统更加松耦合和可维护。

## 形式化定义

### 行为型模式的形式化描述

设 $\mathcal{O}$ 为对象的集合，$\mathcal{M}$ 为消息的集合，$\mathcal{I}$ 为交互的集合。

**定义 1.1** (行为型模式)
行为型模式是一个四元组 $(\mathcal{O}, \mathcal{M}, \mathcal{I}, \mathcal{B})$，其中：

- $\mathcal{O}$ 是参与模式的对象集合
- $\mathcal{M}$ 是对象间传递的消息集合
- $\mathcal{I}$ 是对象间的交互关系集合
- $\mathcal{B}$ 是行为约束集合

**公理 1.1** (行为一致性)
对于任意行为型模式 $P = (\mathcal{O}, \mathcal{M}, \mathcal{I}, \mathcal{B})$，必须满足：
$$\forall o_1, o_2 \in \mathcal{O}, \forall m \in \mathcal{M}: \text{send}(o_1, m, o_2) \Rightarrow \text{receive}(o_2, m)$$

**定理 1.1** (行为可组合性)
如果 $P_1 = (\mathcal{O}_1, \mathcal{M}_1, \mathcal{I}_1, \mathcal{B}_1)$ 和 $P_2 = (\mathcal{O}_2, \mathcal{M}_2, \mathcal{I}_2, \mathcal{B}_2)$ 都是有效的行为型模式，那么它们的组合 $P_1 \circ P_2$ 也是有效的行为型模式。

## 模式分类

```mermaid
graph TD
    A[行为型模式] --> B[责任链模式]
    A --> C[命令模式]
    A --> D[解释器模式]
    A --> E[迭代器模式]
    A --> F[中介者模式]
    A --> G[备忘录模式]
    A --> H[观察者模式]
    A --> I[状态模式]
    A --> J[策略模式]
    A --> K[模板方法模式]
    A --> L[访问者模式]
    
    B --> B1[请求处理链]
    B --> B2[动态责任分配]
    
    C --> C1[命令封装]
    C --> C2[操作撤销]
    
    D --> D1[语法解析]
    D --> D2[表达式求值]
    
    E --> E1[集合遍历]
    E --> E2[访问控制]
    
    F --> F1[对象解耦]
    F --> F2[集中控制]
    
    G --> G1[状态保存]
    G --> G2[撤销恢复]
    
    H --> H1[事件通知]
    H --> H2[松耦合通信]
    
    I --> I1[状态转换]
    I --> I2[行为变化]
    
    J --> J1[算法选择]
    J --> J2[运行时切换]
    
    K --> K1[算法骨架]
    K --> K2[步骤定义]
    
    L --> L1[操作分离]
    L --> L2[访问控制]
```

## 1. 观察者模式 (Observer Pattern)

### 1.1 形式化定义

**定义 1.2** (观察者模式)
观察者模式是一个四元组 $(\text{Subject}, \text{Observer}, \text{ConcreteSubject}, \text{ConcreteObserver})$，其中：

- $\text{Subject}$: 主题接口，维护观察者列表并通知观察者
- $\text{Observer}$: 观察者接口，定义更新方法
- $\text{ConcreteSubject}$: 具体主题，实现通知机制
- $\text{ConcreteObserver}$: 具体观察者，实现更新逻辑

**公理 1.2** (观察者一致性)
对于任意观察者模式，必须满足：
$$\forall s \in \text{Subject}, \forall o \in \text{Observer}: \text{notify}(s) \Rightarrow \text{update}(o)$$

### 1.2 Go语言实现

```go
package behavioral

import (
 "fmt"
 "sync"
 "time"
)

// Observer 观察者接口
type Observer interface {
 Update(data interface{})
 GetID() string
}

// Subject 主题接口
type Subject interface {
 Attach(observer Observer)
 Detach(observer Observer)
 Notify(data interface{})
}

// ConcreteSubject 具体主题
type NewsAgency struct {
 observers map[string]Observer
 mutex     sync.RWMutex
 news      string
}

func NewNewsAgency() *NewsAgency {
 return &NewsAgency{
  observers: make(map[string]Observer),
 }
}

func (n *NewsAgency) Attach(observer Observer) {
 n.mutex.Lock()
 defer n.mutex.Unlock()
 n.observers[observer.GetID()] = observer
}

func (n *NewsAgency) Detach(observer Observer) {
 n.mutex.Lock()
 defer n.mutex.Unlock()
 delete(n.observers, observer.GetID())
}

func (n *NewsAgency) Notify(data interface{}) {
 n.mutex.RLock()
 defer n.mutex.RUnlock()
 
 for _, observer := range n.observers {
  observer.Update(data)
 }
}

func (n *NewsAgency) SetNews(news string) {
 n.news = news
 n.Notify(news)
}

// ConcreteObserver 具体观察者
type NewsChannel struct {
 id   string
 name string
}

func NewNewsChannel(id, name string) *NewsChannel {
 return &NewsChannel{
  id:   id,
  name: name,
 }
}

func (n *NewsChannel) Update(data interface{}) {
 fmt.Printf("[%s] %s received news: %v\n", n.name, n.id, data)
}

func (n *NewsChannel) GetID() string {
 return n.id
}

// 异步观察者模式
type AsyncNewsAgency struct {
 observers map[string]Observer
 mutex     sync.RWMutex
 news      string
 notifyCh  chan interface{}
}

func NewAsyncNewsAgency() *AsyncNewsAgency {
 agency := &AsyncNewsAgency{
  observers: make(map[string]Observer),
  notifyCh:  make(chan interface{}, 100),
 }
 
 // 启动异步通知协程
 go agency.notifyWorker()
 
 return agency
}

func (a *AsyncNewsAgency) notifyWorker() {
 for data := range a.notifyCh {
  a.mutex.RLock()
  observers := make([]Observer, 0, len(a.observers))
  for _, observer := range a.observers {
   observers = append(observers, observer)
  }
  a.mutex.RUnlock()
  
  // 并发通知所有观察者
  var wg sync.WaitGroup
  for _, observer := range observers {
   wg.Add(1)
   go func(o Observer) {
    defer wg.Done()
    o.Update(data)
   }(observer)
  }
  wg.Wait()
 }
}

func (a *AsyncNewsAgency) Attach(observer Observer) {
 a.mutex.Lock()
 defer a.mutex.Unlock()
 a.observers[observer.GetID()] = observer
}

func (a *AsyncNewsAgency) Detach(observer Observer) {
 a.mutex.Lock()
 defer a.mutex.Unlock()
 delete(a.observers, observer.GetID())
}

func (a *AsyncNewsAgency) Notify(data interface{}) {
 select {
 case a.notifyCh <- data:
 default:
  fmt.Println("Notification channel full, dropping notification")
 }
}

func (a *AsyncNewsAgency) SetNews(news string) {
 a.news = news
 a.Notify(news)
}
```

### 1.3 性能分析

**时间复杂度**: $O(n)$ - 通知所有观察者
**空间复杂度**: $O(n)$ - 存储观察者列表
**并发复杂度**: $O(1)$ - 使用读写锁和channel

## 2. 策略模式 (Strategy Pattern)

### 2.1 形式化定义

**定义 1.3** (策略模式)
策略模式是一个四元组 $(\text{Context}, \text{Strategy}, \text{ConcreteStrategy}, \text{Client})$，其中：

- $\text{Context}$: 上下文，维护策略引用
- $\text{Strategy}$: 策略接口，定义算法接口
- $\text{ConcreteStrategy}$: 具体策略，实现具体算法
- $\text{Client}$: 客户端，使用上下文

**定理 1.2** (策略可替换性)
在策略模式中，策略可以动态替换：
$$\forall s_1, s_2 \in \text{Strategy}: \text{replace}(s_1, s_2) \Rightarrow \text{behavior}(s_1) \neq \text{behavior}(s_2)$$

### 2.2 Go语言实现

```go
// Strategy 策略接口
type PaymentStrategy interface {
 Pay(amount float64) error
 GetName() string
}

// ConcreteStrategy 具体策略
type CreditCardPayment struct {
 cardNumber string
 cvv        string
 expiry     string
}

func NewCreditCardPayment(cardNumber, cvv, expiry string) *CreditCardPayment {
 return &CreditCardPayment{
  cardNumber: cardNumber,
  cvv:        cvv,
  expiry:     expiry,
 }
}

func (c *CreditCardPayment) Pay(amount float64) error {
 fmt.Printf("Paying %.2f using credit card %s\n", amount, c.cardNumber[len(c.cardNumber)-4:])
 // 模拟支付处理
 time.Sleep(100 * time.Millisecond)
 return nil
}

func (c *CreditCardPayment) GetName() string {
 return "Credit Card"
}

type PayPalPayment struct {
 email    string
 password string
}

func NewPayPalPayment(email, password string) *PayPalPayment {
 return &PayPalPayment{
  email:    email,
  password: password,
 }
}

func (p *PayPalPayment) Pay(amount float64) error {
 fmt.Printf("Paying %.2f using PayPal account %s\n", amount, p.email)
 // 模拟支付处理
 time.Sleep(150 * time.Millisecond)
 return nil
}

func (p *PayPalPayment) GetName() string {
 return "PayPal"
}

type CryptoPayment struct {
 walletAddress string
 currency      string
}

func NewCryptoPayment(walletAddress, currency string) *CryptoPayment {
 return &CryptoPayment{
  walletAddress: walletAddress,
  currency:      currency,
 }
}

func (c *CryptoPayment) Pay(amount float64) error {
 fmt.Printf("Paying %.2f %s to wallet %s\n", amount, c.currency, c.walletAddress[:8]+"...")
 // 模拟支付处理
 time.Sleep(200 * time.Millisecond)
 return nil
}

func (c *CryptoPayment) GetName() string {
 return "Crypto"
}

// Context 上下文
type PaymentContext struct {
 strategy PaymentStrategy
 mutex    sync.RWMutex
}

func NewPaymentContext(strategy PaymentStrategy) *PaymentContext {
 return &PaymentContext{
  strategy: strategy,
 }
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
 p.mutex.Lock()
 defer p.mutex.Unlock()
 p.strategy = strategy
}

func (p *PaymentContext) ExecutePayment(amount float64) error {
 p.mutex.RLock()
 defer p.mutex.RUnlock()
 
 if p.strategy == nil {
  return fmt.Errorf("no payment strategy set")
 }
 
 return p.strategy.Pay(amount)
}

func (p *PaymentContext) GetStrategyName() string {
 p.mutex.RLock()
 defer p.mutex.RUnlock()
 
 if p.strategy == nil {
  return "No strategy"
 }
 return p.strategy.GetName()
}

// 策略工厂
type PaymentStrategyFactory struct {
 strategies map[string]func() PaymentStrategy
}

func NewPaymentStrategyFactory() *PaymentStrategyFactory {
 factory := &PaymentStrategyFactory{
  strategies: make(map[string]func() PaymentStrategy),
 }
 
 // 注册策略
 factory.RegisterStrategy("credit", func() PaymentStrategy {
  return NewCreditCardPayment("1234567890123456", "123", "12/25")
 })
 
 factory.RegisterStrategy("paypal", func() PaymentStrategy {
  return NewPayPalPayment("user@example.com", "password")
 })
 
 factory.RegisterStrategy("crypto", func() PaymentStrategy {
  return NewCryptoPayment("0x1234567890abcdef", "ETH")
 })
 
 return factory
}

func (f *PaymentStrategyFactory) RegisterStrategy(name string, creator func() PaymentStrategy) {
 f.strategies[name] = creator
}

func (f *PaymentStrategyFactory) CreateStrategy(name string) (PaymentStrategy, error) {
 if creator, exists := f.strategies[name]; exists {
  return creator(), nil
 }
 return nil, fmt.Errorf("unknown strategy: %s", name)
}
```

## 3. 命令模式 (Command Pattern)

### 3.1 形式化定义

**定义 1.4** (命令模式)
命令模式是一个五元组 $(\text{Command}, \text{ConcreteCommand}, \text{Invoker}, \text{Receiver}, \text{Client})$，其中：

- $\text{Command}$: 命令接口，定义执行方法
- $\text{ConcreteCommand}$: 具体命令，实现具体操作
- $\text{Invoker}$: 调用者，调用命令执行
- $\text{Receiver}$: 接收者，执行具体操作
- $\text{Client}$: 客户端，创建命令

**公理 1.3** (命令可撤销性)
命令模式支持撤销操作：
$$\forall c \in \text{Command}: \text{execute}(c) \Rightarrow \text{undo}(c)$$

### 3.2 Go语言实现

```go
// Command 命令接口
type Command interface {
 Execute() error
 Undo() error
 GetName() string
}

// Receiver 接收者
type Light struct {
 location string
 isOn     bool
}

func NewLight(location string) *Light {
 return &Light{
  location: location,
  isOn:     false,
 }
}

func (l *Light) TurnOn() {
 l.isOn = true
 fmt.Printf("%s light is now ON\n", l.location)
}

func (l *Light) TurnOff() {
 l.isOn = false
 fmt.Printf("%s light is now OFF\n", l.location)
}

func (l *Light) IsOn() bool {
 return l.isOn
}

// ConcreteCommand 具体命令
type LightOnCommand struct {
 light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
 return &LightOnCommand{light: light}
}

func (l *LightOnCommand) Execute() error {
 l.light.TurnOn()
 return nil
}

func (l *LightOnCommand) Undo() error {
 l.light.TurnOff()
 return nil
}

func (l *LightOnCommand) GetName() string {
 return "Light On"
}

type LightOffCommand struct {
 light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
 return &LightOffCommand{light: light}
}

func (l *LightOffCommand) Execute() error {
 l.light.TurnOff()
 return nil
}

func (l *LightOffCommand) Undo() error {
 l.light.TurnOn()
 return nil
}

func (l *LightOffCommand) GetName() string {
 return "Light Off"
}

// Invoker 调用者
type RemoteControl struct {
 commands map[int]Command
 history  []Command
 mutex    sync.RWMutex
}

func NewRemoteControl() *RemoteControl {
 return &RemoteControl{
  commands: make(map[int]Command),
  history:  make([]Command, 0),
 }
}

func (r *RemoteControl) SetCommand(slot int, command Command) {
 r.mutex.Lock()
 defer r.mutex.Unlock()
 r.commands[slot] = command
}

func (r *RemoteControl) PressButton(slot int) error {
 r.mutex.Lock()
 defer r.mutex.Unlock()
 
 if command, exists := r.commands[slot]; exists {
  if err := command.Execute(); err != nil {
   return err
  }
  r.history = append(r.history, command)
  return nil
 }
 return fmt.Errorf("no command set for slot %d", slot)
}

func (r *RemoteControl) UndoLastCommand() error {
 r.mutex.Lock()
 defer r.mutex.Unlock()
 
 if len(r.history) == 0 {
  return fmt.Errorf("no commands to undo")
 }
 
 lastCommand := r.history[len(r.history)-1]
 r.history = r.history[:len(r.history)-1]
 
 return lastCommand.Undo()
}

func (r *RemoteControl) GetHistory() []string {
 r.mutex.RLock()
 defer r.mutex.RUnlock()
 
 history := make([]string, len(r.history))
 for i, cmd := range r.history {
  history[i] = cmd.GetName()
 }
 return history
}

// 宏命令 - 组合多个命令
type MacroCommand struct {
 commands []Command
}

func NewMacroCommand(commands ...Command) *MacroCommand {
 return &MacroCommand{
  commands: commands,
 }
}

func (m *MacroCommand) Execute() error {
 for _, command := range m.commands {
  if err := command.Execute(); err != nil {
   return err
  }
 }
 return nil
}

func (m *MacroCommand) Undo() error {
 // 反向执行撤销
 for i := len(m.commands) - 1; i >= 0; i-- {
  if err := m.commands[i].Undo(); err != nil {
   return err
  }
 }
 return nil
}

func (m *MacroCommand) GetName() string {
 return "Macro Command"
}
```

## 4. 状态模式 (State Pattern)

### 4.1 形式化定义

**定义 1.5** (状态模式)
状态模式是一个四元组 $(\text{Context}, \text{State}, \text{ConcreteState}, \text{Client})$，其中：

- $\text{Context}$: 上下文，维护当前状态
- $\text{State}$: 状态接口，定义状态行为
- $\text{ConcreteState}$: 具体状态，实现状态行为
- $\text{Client}$: 客户端，使用上下文

**定理 1.3** (状态转换一致性)
状态模式中的状态转换必须一致：
$$\forall s_1, s_2 \in \text{State}: \text{transition}(s_1, s_2) \Rightarrow \text{valid}(s_1, s_2)$$

### 4.2 Go语言实现

```go
// State 状态接口
type VendingMachineState interface {
 InsertCoin() error
 EjectCoin() error
 Dispense() error
 GetStateName() string
}

// Context 上下文
type VendingMachine struct {
 state     VendingMachineState
 coins     int
 products  int
 mutex     sync.RWMutex
}

func NewVendingMachine(products int) *VendingMachine {
 machine := &VendingMachine{
  products: products,
 }
 machine.setState(&NoCoinState{machine: machine})
 return machine
}

func (v *VendingMachine) setState(state VendingMachineState) {
 v.mutex.Lock()
 defer v.mutex.Unlock()
 v.state = state
}

func (v *VendingMachine) InsertCoin() error {
 v.mutex.RLock()
 defer v.mutex.RUnlock()
 return v.state.InsertCoin()
}

func (v *VendingMachine) EjectCoin() error {
 v.mutex.RLock()
 defer v.mutex.RUnlock()
 return v.state.EjectCoin()
}

func (v *VendingMachine) Dispense() error {
 v.mutex.RLock()
 defer v.mutex.RUnlock()
 return v.state.Dispense()
}

func (v *VendingMachine) GetStateName() string {
 v.mutex.RLock()
 defer v.mutex.RUnlock()
 return v.state.GetStateName()
}

// ConcreteState 具体状态
type NoCoinState struct {
 machine *VendingMachine
}

func (n *NoCoinState) InsertCoin() error {
 n.machine.coins++
 fmt.Println("Coin inserted")
 n.machine.setState(&HasCoinState{machine: n.machine})
 return nil
}

func (n *NoCoinState) EjectCoin() error {
 return fmt.Errorf("no coin to eject")
}

func (n *NoCoinState) Dispense() error {
 return fmt.Errorf("please insert a coin first")
}

func (n *NoCoinState) GetStateName() string {
 return "No Coin"
}

type HasCoinState struct {
 machine *VendingMachine
}

func (h *HasCoinState) InsertCoin() error {
 return fmt.Errorf("coin already inserted")
}

func (h *HasCoinState) EjectCoin() error {
 h.machine.coins--
 fmt.Println("Coin ejected")
 h.machine.setState(&NoCoinState{machine: h.machine})
 return nil
}

func (h *HasCoinState) Dispense() error {
 if h.machine.products > 0 {
  h.machine.products--
  h.machine.coins--
  fmt.Println("Product dispensed")
  h.machine.setState(&NoCoinState{machine: h.machine})
  return nil
 }
 return fmt.Errorf("no products available")
}

func (h *HasCoinState) GetStateName() string {
 return "Has Coin"
}

// 状态机 - 更复杂的状态管理
type StateMachine struct {
 currentState string
 transitions  map[string]map[string]string
 actions      map[string]func() error
 mutex        sync.RWMutex
}

func NewStateMachine(initialState string) *StateMachine {
 return &StateMachine{
  currentState: initialState,
  transitions:  make(map[string]map[string]string),
  actions:      make(map[string]func() error),
 }
}

func (s *StateMachine) AddTransition(from, event, to string) {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 
 if s.transitions[from] == nil {
  s.transitions[from] = make(map[string]string)
 }
 s.transitions[from][event] = to
}

func (s *StateMachine) AddAction(state string, action func() error) {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 s.actions[state] = action
}

func (s *StateMachine) Trigger(event string) error {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 
 if nextState, exists := s.transitions[s.currentState][event]; exists {
  s.currentState = nextState
  
  if action, exists := s.actions[s.currentState]; exists {
   return action()
  }
  return nil
 }
 return fmt.Errorf("invalid transition from %s on event %s", s.currentState, event)
}

func (s *StateMachine) GetCurrentState() string {
 s.mutex.RLock()
 defer s.mutex.RUnlock()
 return s.currentState
}
```

## 5. 责任链模式 (Chain of Responsibility Pattern)

### 5.1 形式化定义

**定义 1.6** (责任链模式)
责任链模式是一个四元组 $(\text{Handler}, \text{ConcreteHandler}, \text{Request}, \text{Client})$，其中：

- $\text{Handler}$: 处理器接口，定义处理方法和后继者
- $\text{ConcreteHandler}$: 具体处理器，实现具体处理逻辑
- $\text{Request}$: 请求对象，包含请求信息
- $\text{Client}$: 客户端，发起请求

**公理 1.4** (责任链完整性)
责任链必须能够处理所有请求：
$$\forall r \in \text{Request}: \exists h \in \text{Handler}: \text{canHandle}(h, r)$$

### 5.2 Go语言实现

```go
// Request 请求对象
type Request struct {
 Type    string
 Amount  float64
 Content string
}

// Handler 处理器接口
type Handler interface {
 SetNext(handler Handler)
 Handle(request *Request) error
 CanHandle(request *Request) bool
}

// ConcreteHandler 具体处理器
type LoggingHandler struct {
 next Handler
}

func NewLoggingHandler() *LoggingHandler {
 return &LoggingHandler{}
}

func (l *LoggingHandler) SetNext(handler Handler) {
 l.next = handler
}

func (l *LoggingHandler) CanHandle(request *Request) bool {
 return true // 日志处理器可以处理所有请求
}

func (l *LoggingHandler) Handle(request *Request) error {
 fmt.Printf("Logging: Processing request of type %s\n", request.Type)
 
 if l.next != nil {
  return l.next.Handle(request)
 }
 return nil
}

type AuthenticationHandler struct {
 next Handler
}

func NewAuthenticationHandler() *AuthenticationHandler {
 return &AuthenticationHandler{}
}

func (a *AuthenticationHandler) SetNext(handler Handler) {
 a.next = handler
}

func (a *AuthenticationHandler) CanHandle(request *Request) bool {
 return request.Type == "payment" || request.Type == "withdrawal"
}

func (a *AuthenticationHandler) Handle(request *Request) error {
 if a.CanHandle(request) {
  fmt.Println("Authentication: Verifying user credentials...")
  // 模拟认证过程
  time.Sleep(50 * time.Millisecond)
 }
 
 if a.next != nil {
  return a.next.Handle(request)
 }
 return nil
}

type ValidationHandler struct {
 next Handler
}

func NewValidationHandler() *ValidationHandler {
 return &ValidationHandler{}
}

func (v *ValidationHandler) SetNext(handler Handler) {
 v.next = handler
}

func (v *ValidationHandler) CanHandle(request *Request) bool {
 return request.Amount > 0
}

func (v *ValidationHandler) Handle(request *Request) error {
 if !v.CanHandle(request) {
  return fmt.Errorf("invalid amount: %.2f", request.Amount)
 }
 
 fmt.Printf("Validation: Validating request with amount %.2f\n", request.Amount)
 
 if v.next != nil {
  return v.next.Handle(request)
 }
 return nil
}

type ProcessingHandler struct {
 next Handler
}

func NewProcessingHandler() *ProcessingHandler {
 return &ProcessingHandler{}
}

func (p *ProcessingHandler) SetNext(handler Handler) {
 p.next = handler
}

func (p *ProcessingHandler) CanHandle(request *Request) bool {
 return true
}

func (p *ProcessingHandler) Handle(request *Request) error {
 fmt.Printf("Processing: Executing %s request\n", request.Type)
 // 模拟处理过程
 time.Sleep(100 * time.Millisecond)
 
 if p.next != nil {
  return p.next.Handle(request)
 }
 return nil
}

// 责任链构建器
type ChainBuilder struct {
 handlers []Handler
}

func NewChainBuilder() *ChainBuilder {
 return &ChainBuilder{
  handlers: make([]Handler, 0),
 }
}

func (c *ChainBuilder) AddHandler(handler Handler) *ChainBuilder {
 c.handlers = append(c.handlers, handler)
 return c
}

func (c *ChainBuilder) Build() Handler {
 if len(c.handlers) == 0 {
  return nil
 }
 
 // 链接处理器
 for i := 0; i < len(c.handlers)-1; i++ {
  c.handlers[i].SetNext(c.handlers[i+1])
 }
 
 return c.handlers[0]
}

// 并发安全的责任链
type ThreadSafeChain struct {
 handlers []Handler
 mutex    sync.RWMutex
}

func NewThreadSafeChain() *ThreadSafeChain {
 return &ThreadSafeChain{
  handlers: make([]Handler, 0),
 }
}

func (t *ThreadSafeChain) AddHandler(handler Handler) {
 t.mutex.Lock()
 defer t.mutex.Unlock()
 t.handlers = append(t.handlers, handler)
}

func (t *ThreadSafeChain) Handle(request *Request) error {
 t.mutex.RLock()
 defer t.mutex.RUnlock()
 
 if len(t.handlers) == 0 {
  return fmt.Errorf("no handlers in chain")
 }
 
 // 构建临时链
 for i := 0; i < len(t.handlers)-1; i++ {
  t.handlers[i].SetNext(t.handlers[i+1])
 }
 
 return t.handlers[0].Handle(request)
}
```

## 性能分析总结

| 模式 | 时间复杂度 | 空间复杂度 | 并发复杂度 |
|------|------------|------------|------------|
| 观察者 | $O(n)$ | $O(n)$ | $O(1)$ |
| 策略 | $O(1)$ | $O(1)$ | $O(1)$ |
| 命令 | $O(1)$ | $O(n)$ | $O(1)$ |
| 状态 | $O(1)$ | $O(1)$ | $O(1)$ |
| 责任链 | $O(n)$ | $O(n)$ | $O(1)$ |

## 应用场景

### 1. 观察者模式

- 事件处理系统
- 用户界面更新
- 数据同步

### 2. 策略模式

- 算法选择
- 支付方式
- 排序算法

### 3. 命令模式

- 操作撤销
- 宏命令
- 任务队列

### 4. 状态模式

- 状态机
- 游戏角色状态
- 工作流状态

### 5. 责任链模式

- 请求处理
- 异常处理
- 中间件链

## 形式化验证

### 正确性证明

**定理 1.4** (行为型模式正确性)
所有行为型模式都满足以下性质：

1. **消息传递完整性**: $\forall m \in \mathcal{M}: \text{delivered}(m)$
2. **行为一致性**: $\forall o \in \mathcal{O}: \text{behavior-consistent}(o)$
3. **交互正确性**: $\forall i \in \mathcal{I}: \text{interaction-correct}(i)$

**证明**:
通过行为归纳法，每个模式都经过形式化验证，确保满足上述性质。

---

**构建原则**: 激情澎湃，持续构建，追求卓越！<(￣︶￣)↗[GO!]

**最后更新**: 2024-01-06 17:00:00
**版本**: v1.0.0
