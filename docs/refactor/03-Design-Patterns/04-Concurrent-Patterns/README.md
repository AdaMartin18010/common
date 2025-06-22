# 并发设计模式 (Concurrent Design Patterns)

## 概述

并发设计模式是专门处理并发和多线程环境下的设计模式，旨在解决资源共享、任务协调、同步和通信等问题。本文档采用严格的数学形式化方法，结合 Go 语言的并发特性，对并发模式进行系统性重构。

## 形式化定义

### 1. 并发模式的形式化框架

**定义 1.1** (并发模式)
并发模式是一个四元组 ```latex
\mathcal{CP} = (T, S, R, \phi)
```，其中：

- ```latex
T
``` 是线程/协程集合
- ```latex
S
``` 是共享状态集合
- ```latex
R
``` 是同步关系集合
- ```latex
\phi: T \times S \rightarrow R
``` 是线程到同步关系的映射

**公理 1.1** (并发模式公理)
对于任意并发模式 ```latex
\mathcal{CP} = (T, S, R, \phi)
```：

1. **安全性**: ```latex
\forall t_1, t_2 \in T: \text{race\_free}(t_1, t_2)
```
2. **活性**: ```latex
\forall t \in T: \text{eventually\_progress}(t)
```
3. **公平性**: ```latex
\forall t_1, t_2 \in T: \text{fair\_scheduling}(t_1, t_2)
```

### 2. 并发安全性的形式化

**定义 1.2** (数据竞争)
对于线程 ```latex
t_1, t_2
``` 和共享状态 ```latex
s
```，存在数据竞争当且仅当：
$```latex
\exists t_1, t_2 \in T, s \in S: \text{concurrent\_access}(t_1, t_2, s) \land \text{one\_write}(t_1, t_2, s)
```$

**定义 1.3** (死锁)
死锁是线程集合 ```latex
T' \subseteq T
``` 的状态，满足：
$```latex
\forall t \in T': \text{waiting\_for}(t) \in T' \land \text{circular\_wait}(T')
```$

## 核心模式

### 1. 线程池模式 (Thread Pool Pattern)

**定义 1.4** (线程池模式)
线程池模式预先创建一组线程，用于执行任务，避免频繁创建和销毁线程的开销。

**形式化定义**:
设 ```latex
W
``` 为工作线程集合，```latex
Q
``` 为任务队列，线程池模式定义为：
$```latex
\text{ThreadPool}: Q \rightarrow W \times \text{Result}
```$

**Go 语言实现**:

```go
package threadpool

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Task 任务接口
type Task interface {
    Execute() (interface{}, error)
    GetID() string
}

// SimpleTask 简单任务实现
type SimpleTask struct {
    id       string
    function func() (interface{}, error)
}

func NewSimpleTask(id string, fn func() (interface{}, error)) *SimpleTask {
    return &SimpleTask{
        id:       id,
        function: fn,
    }
}

func (t *SimpleTask) Execute() (interface{}, error) {
    return t.function()
}

func (t *SimpleTask) GetID() string {
    return t.id
}

// Worker 工作线程
type Worker struct {
    id       int
    taskChan <-chan Task
    resultChan chan<- TaskResult
    wg        *sync.WaitGroup
    ctx       context.Context
}

type TaskResult struct {
    TaskID string
    Result interface{}
    Error  error
}

func NewWorker(id int, taskChan <-chan Task, resultChan chan<- TaskResult, wg *sync.WaitGroup, ctx context.Context) *Worker {
    return &Worker{
        id:        id,
        taskChan:  taskChan,
        resultChan: resultChan,
        wg:        wg,
        ctx:       ctx,
    }
}

func (w *Worker) Start() {
    defer w.wg.Done()
    
    for {
        select {
        case task, ok := <-w.taskChan:
            if !ok {
                return // 通道已关闭
            }
            
            result, err := task.Execute()
            w.resultChan <- TaskResult{
                TaskID: task.GetID(),
                Result: result,
                Error:  err,
            }
            
        case <-w.ctx.Done():
            return // 上下文取消
        }
    }
}

// ThreadPool 线程池
type ThreadPool struct {
    workers    []*Worker
    taskChan   chan Task
    resultChan chan TaskResult
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewThreadPool(workerCount int, queueSize int) *ThreadPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    pool := &ThreadPool{
        workers:    make([]*Worker, workerCount),
        taskChan:   make(chan Task, queueSize),
        resultChan: make(chan TaskResult, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
    
    // 创建工作线程
    for i := 0; i < workerCount; i++ {
        worker := NewWorker(i, pool.taskChan, pool.resultChan, &pool.wg, pool.ctx)
        pool.workers[i] = worker
        pool.wg.Add(1)
        go worker.Start()
    }
    
    return pool
}

func (tp *ThreadPool) Submit(task Task) error {
    select {
    case tp.taskChan <- task:
        return nil
    case <-tp.ctx.Done():
        return fmt.Errorf("thread pool is shutdown")
    default:
        return fmt.Errorf("task queue is full")
    }
}

func (tp *ThreadPool) GetResult() (TaskResult, error) {
    select {
    case result := <-tp.resultChan:
        return result, nil
    case <-tp.ctx.Done():
        return TaskResult{}, fmt.Errorf("thread pool is shutdown")
    }
}

func (tp *ThreadPool) Shutdown() {
    tp.cancel()
    close(tp.taskChan)
    tp.wg.Wait()
    close(tp.resultChan)
}

// 泛型线程池
type GenericThreadPool[T any, R any] struct {
    workers    int
    taskChan   chan func() (R, error)
    resultChan chan Result[R]
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

type Result[R any] struct {
    Value R
    Error error
}

func NewGenericThreadPool[T any, R any](workerCount int) *GenericThreadPool[T, R] {
    ctx, cancel := context.WithCancel(context.Background())
    
    pool := &GenericThreadPool[T, R]{
        workers:    workerCount,
        taskChan:   make(chan func() (R, error), workerCount*2),
        resultChan: make(chan Result[R], workerCount*2),
        ctx:        ctx,
        cancel:     cancel,
    }
    
    pool.startWorkers()
    return pool
}

func (tp *GenericThreadPool[T, R]) startWorkers() {
    for i := 0; i < tp.workers; i++ {
        tp.wg.Add(1)
        go tp.worker(i)
    }
}

func (tp *GenericThreadPool[T, R]) worker(id int) {
    defer tp.wg.Done()
    
    for {
        select {
        case task, ok := <-tp.taskChan:
            if !ok {
                return
            }
            
            result, err := task()
            tp.resultChan <- Result[R]{
                Value: result,
                Error: err,
            }
            
        case <-tp.ctx.Done():
            return
        }
    }
}

func (tp *GenericThreadPool[T, R]) Submit(task func() (R, error)) error {
    select {
    case tp.taskChan <- task:
        return nil
    case <-tp.ctx.Done():
        return fmt.Errorf("pool is shutdown")
    default:
        return fmt.Errorf("task queue is full")
    }
}

func (tp *GenericThreadPool[T, R]) GetResult() (Result[R], error) {
    select {
    case result := <-tp.resultChan:
        return result, nil
    case <-tp.ctx.Done():
        return Result[R]{}, fmt.Errorf("pool is shutdown")
    }
}

func (tp *GenericThreadPool[T, R]) Shutdown() {
    tp.cancel()
    close(tp.taskChan)
    tp.wg.Wait()
    close(tp.resultChan)
}

// 使用示例
func ExampleThreadPool() {
    // 创建线程池
    pool := NewThreadPool(4, 100)
    defer pool.Shutdown()
    
    // 提交任务
    for i := 0; i < 10; i++ {
        task := NewSimpleTask(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
            time.Sleep(100 * time.Millisecond)
            return fmt.Sprintf("Result from task %d", i), nil
        })
        
        if err := pool.Submit(task); err != nil {
            fmt.Printf("Failed to submit task %d: %v\n", i, err)
        }
    }
    
    // 获取结果
    for i := 0; i < 10; i++ {
        result, err := pool.GetResult()
        if err != nil {
            fmt.Printf("Failed to get result: %v\n", err)
            continue
        }
        
        if result.Error != nil {
            fmt.Printf("Task %s failed: %v\n", result.TaskID, result.Error)
        } else {
            fmt.Printf("Task %s completed: %v\n", result.TaskID, result.Result)
        }
    }
    
    // 使用泛型线程池
    genericPool := NewGenericThreadPool[string, int](2)
    defer genericPool.Shutdown()
    
    // 提交计算任务
    for i := 0; i < 5; i++ {
        i := i // 捕获循环变量
        task := func() (int, error) {
            time.Sleep(50 * time.Millisecond)
            return i * i, nil
        }
        
        if err := genericPool.Submit(task); err != nil {
            fmt.Printf("Failed to submit generic task: %v\n", err)
        }
    }
    
    // 获取计算结果
    for i := 0; i < 5; i++ {
        result, err := genericPool.GetResult()
        if err != nil {
            fmt.Printf("Failed to get generic result: %v\n", err)
            continue
        }
        
        if result.Error != nil {
            fmt.Printf("Generic task failed: %v\n", result.Error)
        } else {
            fmt.Printf("Generic task completed: %d\n", result.Value)
        }
    }
}
```

**形式化证明**:

**定理 1.5** (线程池的安全性)
线程池模式保证任务执行的线程安全：
$```latex
\forall t_1, t_2 \in T, \forall s \in S: \text{no\_race\_condition}(t_1, t_2, s)
```$

**证明**:

1. 每个任务在独立的工作线程中执行
2. 任务队列使用通道，保证线程安全
3. 结果通道保证结果传递的线程安全
4. 证毕。

### 2. Future/Promise 模式

**定义 1.5** (Future/Promise 模式)
Future 表示一个尚未完成的异步计算的结果，Promise 用于设置异步操作的最终结果。

**形式化定义**:
设 ```latex
T
``` 为结果类型，Future/Promise 模式定义为：
$```latex
\text{Future}[T]: \text{Unit} \rightarrow T \times \text{Status}
```$
$```latex
\text{Promise}[T]: T \rightarrow \text{Unit}
```$

**Go 语言实现**:

```go
package future

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Future 表示异步计算的结果
type Future[T any] struct {
    result T
    err    error
    done   chan struct{}
    once   sync.Once
}

// Promise 用于设置异步操作的结果
type Promise[T any] struct {
    future *Future[T]
}

// NewPromise 创建新的 Promise
func NewPromise[T any]() *Promise[T] {
    return &Promise[T]{
        future: &Future[T]{
            done: make(chan struct{}),
        },
    }
}

// Get 获取 Future 的结果，阻塞直到结果可用
func (f *Future[T]) Get() (T, error) {
    <-f.done
    return f.result, f.err
}

// GetWithTimeout 带超时的获取结果
func (f *Future[T]) GetWithTimeout(timeout time.Duration) (T, error) {
    select {
    case <-f.done:
        return f.result, f.err
    case <-time.After(timeout):
        var zero T
        return zero, fmt.Errorf("timeout waiting for result")
    }
}

// IsDone 检查是否完成
func (f *Future[T]) IsDone() bool {
    select {
    case <-f.done:
        return true
    default:
        return false
    }
}

// Set 设置 Promise 的结果
func (p *Promise[T]) Set(result T, err error) {
    p.future.once.Do(func() {
        p.future.result = result
        p.future.err = err
        close(p.future.done)
    })
}

// GetFuture 获取对应的 Future
func (p *Promise[T]) GetFuture() *Future[T] {
    return p.future
}

// 异步任务执行器
type AsyncExecutor struct {
    workerPool chan chan func()
    taskQueue  chan func()
    quit       chan bool
}

func NewAsyncExecutor(workerCount int) *AsyncExecutor {
    executor := &AsyncExecutor{
        workerPool: make(chan chan func(), workerCount),
        taskQueue:  make(chan func(), 100),
        quit:       make(chan bool),
    }
    
    for i := 0; i < workerCount; i++ {
        worker := NewWorker(executor.workerPool, executor.quit)
        worker.Start()
    }
    
    go executor.dispatch()
    return executor
}

type Worker struct {
    workerPool chan chan func()
    taskChan   chan func()
    quit       chan bool
}

func NewWorker(workerPool chan chan func(), quit chan bool) *Worker {
    return &Worker{
        workerPool: workerPool,
        taskChan:   make(chan func()),
        quit:       quit,
    }
}

func (w *Worker) Start() {
    go func() {
        for {
            w.workerPool <- w.taskChan
            
            select {
            case task := <-w.taskChan:
                task()
            case <-w.quit:
                return
            }
        }
    }()
}

func (e *AsyncExecutor) dispatch() {
    for {
        select {
        case task := <-e.taskQueue:
            go func() {
                worker := <-e.workerPool
                worker <- task
            }()
        case <-e.quit:
            return
        }
    }
}

func (e *AsyncExecutor) Submit(task func()) {
    e.taskQueue <- task
}

func (e *AsyncExecutor) Stop() {
    close(e.quit)
}

// 使用示例
func ExampleFuturePromise() {
    // 创建 Promise
    promise := NewPromise[string]()
    future := promise.GetFuture()
    
    // 异步执行任务
    go func() {
        time.Sleep(1 * time.Second)
        promise.Set("Hello, Future!", nil)
    }()
    
    // 获取结果
    result, err := future.Get()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %s\n", result)
    }
    
    // 使用超时
    promise2 := NewPromise[int]()
    future2 := promise2.GetFuture()
    
    go func() {
        time.Sleep(2 * time.Second)
        promise2.Set(42, nil)
    }()
    
    result2, err2 := future2.GetWithTimeout(500 * time.Millisecond)
    if err2 != nil {
        fmt.Printf("Timeout: %v\n", err2)
    } else {
        fmt.Printf("Result: %d\n", result2)
    }
    
    // 使用异步执行器
    executor := NewAsyncExecutor(4)
    defer executor.Stop()
    
    promise3 := NewPromise[string]()
    future3 := promise3.GetFuture()
    
    executor.Submit(func() {
        time.Sleep(500 * time.Millisecond)
        promise3.Set("Async result", nil)
    })
    
    result3, err3 := future3.Get()
    if err3 != nil {
        fmt.Printf("Error: %v\n", err3)
    } else {
        fmt.Printf("Async result: %s\n", result3)
    }
}
```

### 3. Actor 模式

**定义 1.6** (Actor 模式)
Actor 模式中，每个 Actor 是一个并发执行的实体，拥有自己的状态和行为，通过消息传递与其他 Actor 通信。

**形式化定义**:
设 ```latex
A
``` 为 Actor 集合，```latex
M
``` 为消息集合，Actor 模式定义为：
$```latex
\text{Actor}: A \times M \rightarrow A \times \text{Response}
```$

**Go 语言实现**:

```go
package actor

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Message 消息接口
type Message interface {
    GetType() string
    GetData() interface{}
}

// SimpleMessage 简单消息实现
type SimpleMessage struct {
    msgType string
    data    interface{}
}

func NewSimpleMessage(msgType string, data interface{}) *SimpleMessage {
    return &SimpleMessage{
        msgType: msgType,
        data:    data,
    }
}

func (m *SimpleMessage) GetType() string {
    return m.msgType
}

func (m *SimpleMessage) GetData() interface{} {
    return m.data
}

// Response 响应
type Response struct {
    Data interface{}
    Err  error
}

// Actor 接口
type Actor interface {
    Receive(msg Message) Response
    Start()
    Stop()
}

// BaseActor 基础 Actor 实现
type BaseActor struct {
    id       string
    mailbox  chan Message
    response chan Response
    ctx      context.Context
    cancel   context.CancelFunc
    wg       sync.WaitGroup
    handler  func(Message) Response
}

func NewBaseActor(id string, handler func(Message) Response) *BaseActor {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &BaseActor{
        id:       id,
        mailbox:  make(chan Message, 100),
        response: make(chan Response, 100),
        ctx:      ctx,
        cancel:   cancel,
        handler:  handler,
    }
}

func (a *BaseActor) Start() {
    a.wg.Add(1)
    go a.process()
}

func (a *BaseActor) process() {
    defer a.wg.Done()
    
    for {
        select {
        case msg := <-a.mailbox:
            response := a.handler(msg)
            a.response <- response
            
        case <-a.ctx.Done():
            return
        }
    }
}

func (a *BaseActor) Stop() {
    a.cancel()
    a.wg.Wait()
    close(a.mailbox)
    close(a.response)
}

func (a *BaseActor) Send(msg Message) Response {
    select {
    case a.mailbox <- msg:
        select {
        case response := <-a.response:
            return response
        case <-a.ctx.Done():
            return Response{Err: fmt.Errorf("actor stopped")}
        }
    case <-a.ctx.Done():
        return Response{Err: fmt.Errorf("actor stopped")}
    }
}

func (a *BaseActor) SendAsync(msg Message) {
    select {
    case a.mailbox <- msg:
    case <-a.ctx.Done():
    }
}

// ActorSystem Actor 系统
type ActorSystem struct {
    actors map[string]Actor
    mu     sync.RWMutex
}

func NewActorSystem() *ActorSystem {
    return &ActorSystem{
        actors: make(map[string]Actor),
    }
}

func (as *ActorSystem) Register(id string, actor Actor) {
    as.mu.Lock()
    defer as.mu.Unlock()
    as.actors[id] = actor
    actor.Start()
}

func (as *ActorSystem) Send(id string, msg Message) (Response, error) {
    as.mu.RLock()
    actor, exists := as.actors[id]
    as.mu.RUnlock()
    
    if !exists {
        return Response{}, fmt.Errorf("actor %s not found", id)
    }
    
    return actor.(*BaseActor).Send(msg), nil
}

func (as *ActorSystem) SendAsync(id string, msg Message) error {
    as.mu.RLock()
    actor, exists := as.actors[id]
    as.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("actor %s not found", id)
    }
    
    actor.(*BaseActor).SendAsync(msg)
    return nil
}

func (as *ActorSystem) Stop() {
    as.mu.Lock()
    defer as.mu.Unlock()
    
    for _, actor := range as.actors {
        actor.Stop()
    }
    as.actors = make(map[string]Actor)
}

// 具体 Actor 实现示例
type CalculatorActor struct {
    *BaseActor
    state int
}

func NewCalculatorActor(id string) *CalculatorActor {
    actor := &CalculatorActor{
        state: 0,
    }
    
    actor.BaseActor = NewBaseActor(id, actor.handleMessage)
    return actor
}

func (ca *CalculatorActor) handleMessage(msg Message) Response {
    switch msg.GetType() {
    case "add":
        if value, ok := msg.GetData().(int); ok {
            ca.state += value
            return Response{Data: ca.state}
        }
        return Response{Err: fmt.Errorf("invalid data type for add")}
        
    case "subtract":
        if value, ok := msg.GetData().(int); ok {
            ca.state -= value
            return Response{Data: ca.state}
        }
        return Response{Err: fmt.Errorf("invalid data type for subtract")}
        
    case "get":
        return Response{Data: ca.state}
        
    default:
        return Response{Err: fmt.Errorf("unknown message type: %s", msg.GetType())}
    }
}

// 使用示例
func ExampleActor() {
    // 创建 Actor 系统
    system := NewActorSystem()
    defer system.Stop()
    
    // 创建计算器 Actor
    calculator := NewCalculatorActor("calculator")
    system.Register("calculator", calculator)
    
    // 发送消息
    addMsg := NewSimpleMessage("add", 10)
    response, err := system.Send("calculator", addMsg)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Add result: %v\n", response.Data)
    }
    
    subtractMsg := NewSimpleMessage("subtract", 3)
    response, err = system.Send("calculator", subtractMsg)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Subtract result: %v\n", response.Data)
    }
    
    getMsg := NewSimpleMessage("get", nil)
    response, err = system.Send("calculator", getMsg)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Current state: %v\n", response.Data)
    }
    
    // 异步发送消息
    for i := 0; i < 5; i++ {
        addMsg := NewSimpleMessage("add", i)
        if err := system.SendAsync("calculator", addMsg); err != nil {
            fmt.Printf("Async send error: %v\n", err)
        }
    }
    
    time.Sleep(100 * time.Millisecond)
    
    getMsg = NewSimpleMessage("get", nil)
    response, err = system.Send("calculator", getMsg)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Final state: %v\n", response.Data)
    }
}
```

### 4. 生产者-消费者模式

**定义 1.7** (生产者-消费者模式)
生产者-消费者模式通过队列协调生产者和消费者之间的数据流，实现解耦和缓冲。

**形式化定义**:
设 ```latex
P
``` 为生产者集合，```latex
C
``` 为消费者集合，```latex
Q
``` 为队列，模式定义为：
$```latex
P \times Q \rightarrow Q \times C \rightarrow \text{Result}
```$

**Go 语言实现**:

```go
package producerconsumer

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Item 数据项
type Item struct {
    ID   string
    Data interface{}
}

// Producer 生产者接口
type Producer interface {
    Produce() (Item, error)
    Stop()
}

// Consumer 消费者接口
type Consumer interface {
    Consume(item Item) error
    Stop()
}

// SimpleProducer 简单生产者
type SimpleProducer struct {
    id       string
    counter  int
    itemChan chan<- Item
    ctx      context.Context
    cancel   context.CancelFunc
    wg       sync.WaitGroup
}

func NewSimpleProducer(id string, itemChan chan<- Item) *SimpleProducer {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &SimpleProducer{
        id:       id,
        itemChan: itemChan,
        ctx:      ctx,
        cancel:   cancel,
    }
}

func (p *SimpleProducer) Start() {
    p.wg.Add(1)
    go p.produce()
}

func (p *SimpleProducer) produce() {
    defer p.wg.Done()
    
    for {
        select {
        case <-p.ctx.Done():
            return
            
        default:
            item := Item{
                ID:   fmt.Sprintf("%s-%d", p.id, p.counter),
                Data: fmt.Sprintf("Data from producer %s", p.id),
            }
            
            select {
            case p.itemChan <- item:
                p.counter++
                time.Sleep(100 * time.Millisecond)
            case <-p.ctx.Done():
                return
            }
        }
    }
}

func (p *SimpleProducer) Stop() {
    p.cancel()
    p.wg.Wait()
}

// SimpleConsumer 简单消费者
type SimpleConsumer struct {
    id       string
    itemChan <-chan Item
    ctx      context.Context
    cancel   context.CancelFunc
    wg       sync.WaitGroup
}

func NewSimpleConsumer(id string, itemChan <-chan Item) *SimpleConsumer {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &SimpleConsumer{
        id:       id,
        itemChan: itemChan,
        ctx:      ctx,
        cancel:   cancel,
    }
}

func (c *SimpleConsumer) Start() {
    c.wg.Add(1)
    go c.consume()
}

func (c *SimpleConsumer) consume() {
    defer c.wg.Done()
    
    for {
        select {
        case item, ok := <-c.itemChan:
            if !ok {
                return // 通道已关闭
            }
            
            if err := c.Consume(item); err != nil {
                fmt.Printf("Consumer %s failed to consume item %s: %v\n", c.id, item.ID, err)
            } else {
                fmt.Printf("Consumer %s consumed item %s: %v\n", c.id, item.ID, item.Data)
            }
            
        case <-c.ctx.Done():
            return
        }
    }
}

func (c *SimpleConsumer) Consume(item Item) error {
    // 模拟处理时间
    time.Sleep(50 * time.Millisecond)
    return nil
}

func (c *SimpleConsumer) Stop() {
    c.cancel()
    c.wg.Wait()
}

// ProducerConsumerSystem 生产者-消费者系统
type ProducerConsumerSystem struct {
    producers []*SimpleProducer
    consumers []*SimpleConsumer
    itemChan  chan Item
    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

func NewProducerConsumerSystem(producerCount, consumerCount, bufferSize int) *ProducerConsumerSystem {
    ctx, cancel := context.WithCancel(context.Background())
    
    system := &ProducerConsumerSystem{
        producers: make([]*SimpleProducer, producerCount),
        consumers: make([]*SimpleConsumer, consumerCount),
        itemChan:  make(chan Item, bufferSize),
        ctx:       ctx,
        cancel:    cancel,
    }
    
    // 创建生产者
    for i := 0; i < producerCount; i++ {
        producer := NewSimpleProducer(fmt.Sprintf("producer-%d", i), system.itemChan)
        system.producers[i] = producer
    }
    
    // 创建消费者
    for i := 0; i < consumerCount; i++ {
        consumer := NewSimpleConsumer(fmt.Sprintf("consumer-%d", i), system.itemChan)
        system.consumers[i] = consumer
    }
    
    return system
}

func (pcs *ProducerConsumerSystem) Start() {
    // 启动生产者
    for _, producer := range pcs.producers {
        producer.Start()
    }
    
    // 启动消费者
    for _, consumer := range pcs.consumers {
        consumer.Start()
    }
}

func (pcs *ProducerConsumerSystem) Stop() {
    pcs.cancel()
    
    // 停止生产者
    for _, producer := range pcs.producers {
        producer.Stop()
    }
    
    // 停止消费者
    for _, consumer := range pcs.consumers {
        consumer.Stop()
    }
    
    close(pcs.itemChan)
}

// 泛型生产者-消费者系统
type GenericProducerConsumer[T any] struct {
    producers []func() (T, error)
    consumers []func(T) error
    itemChan  chan T
    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

func NewGenericProducerConsumer[T any](bufferSize int) *GenericProducerConsumer[T] {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &GenericProducerConsumer[T]{
        producers: make([]func() (T, error), 0),
        consumers: make([]func(T) error, 0),
        itemChan:  make(chan T, bufferSize),
        ctx:       ctx,
        cancel:    cancel,
    }
}

func (gpc *GenericProducerConsumer[T]) AddProducer(producer func() (T, error)) {
    gpc.producers = append(gpc.producers, producer)
}

func (gpc *GenericProducerConsumer[T]) AddConsumer(consumer func(T) error) {
    gpc.consumers = append(gpc.consumers, consumer)
}

func (gpc *GenericProducerConsumer[T]) Start() {
    // 启动生产者
    for _, producer := range gpc.producers {
        gpc.wg.Add(1)
        go gpc.runProducer(producer)
    }
    
    // 启动消费者
    for _, consumer := range gpc.consumers {
        gpc.wg.Add(1)
        go gpc.runConsumer(consumer)
    }
}

func (gpc *GenericProducerConsumer[T]) runProducer(producer func() (T, error)) {
    defer gpc.wg.Done()
    
    for {
        select {
        case <-gpc.ctx.Done():
            return
            
        default:
            item, err := producer()
            if err != nil {
                fmt.Printf("Producer error: %v\n", err)
                continue
            }
            
            select {
            case gpc.itemChan <- item:
            case <-gpc.ctx.Done():
                return
            }
        }
    }
}

func (gpc *GenericProducerConsumer[T]) runConsumer(consumer func(T) error) {
    defer gpc.wg.Done()
    
    for {
        select {
        case item, ok := <-gpc.itemChan:
            if !ok {
                return
            }
            
            if err := consumer(item); err != nil {
                fmt.Printf("Consumer error: %v\n", err)
            }
            
        case <-gpc.ctx.Done():
            return
        }
    }
}

func (gpc *GenericProducerConsumer[T]) Stop() {
    gpc.cancel()
    gpc.wg.Wait()
    close(gpc.itemChan)
}

// 使用示例
func ExampleProducerConsumer() {
    // 使用传统生产者-消费者系统
    system := NewProducerConsumerSystem(2, 3, 10)
    system.Start()
    
    // 运行一段时间
    time.Sleep(2 * time.Second)
    system.Stop()
    
    // 使用泛型生产者-消费者系统
    genericSystem := NewGenericProducerConsumer[string](5)
    
    // 添加生产者
    counter := 0
    genericSystem.AddProducer(func() (string, error) {
        counter++
        return fmt.Sprintf("Item-%d", counter), nil
    })
    
    genericSystem.AddProducer(func() (string, error) {
        counter++
        return fmt.Sprintf("Data-%d", counter), nil
    })
    
    // 添加消费者
    genericSystem.AddConsumer(func(item string) error {
        fmt.Printf("Consumer 1 processed: %s\n", item)
        return nil
    })
    
    genericSystem.AddConsumer(func(item string) error {
        fmt.Printf("Consumer 2 processed: %s\n", item)
        return nil
    })
    
    genericSystem.Start()
    time.Sleep(1 * time.Second)
    genericSystem.Stop()
}
```

## 模式关系分析

### 1. 并发模式间的数学关系

**定义 1.8** (并发模式组合)
设 ```latex
\mathcal{CP}_1, \mathcal{CP}_2
``` 为两个并发模式，其组合定义为：
$```latex
\mathcal{CP}_1 \circ \mathcal{CP}_2 = (T_1 \cup T_2, S_1 \cup S_2, R_1 \cup R_2, \phi_1 \cup \phi_2)
```$

**定理 1.6** (组合的安全性)
如果 ```latex
\mathcal{CP}_1, \mathcal{CP}_2
``` 都是安全的，则其组合 ```latex
\mathcal{CP}_1 \circ \mathcal{CP}_2
``` 也是安全的。

### 2. 性能分析

| 模式 | 时间复杂度 | 空间复杂度 | 适用场景 |
|------|------------|------------|----------|
| 线程池 | O(1) | O(n) | 任务执行 |
| Future/Promise | O(1) | O(1) | 异步计算 |
| Actor | O(1) | O(n) | 消息传递 |
| 生产者-消费者 | O(1) | O(n) | 数据流处理 |

## 最佳实践

### 1. 模式选择指南

**规则 1.1** (线程池选择)
当且仅当满足以下条件时使用线程池：

1. 需要执行大量短期任务
2. 任务执行时间相对固定
3. 需要控制并发度

**规则 1.2** (Future/Promise 选择)
当且仅当满足以下条件时使用 Future/Promise：

1. 需要异步计算结果
2. 需要组合多个异步操作
3. 需要超时控制

**规则 1.3** (Actor 选择)
当且仅当满足以下条件时使用 Actor：

1. 需要隔离状态
2. 需要消息传递通信
3. 需要容错和恢复

### 2. 性能优化建议

1. **线程池调优**: 根据任务特性调整线程数量
2. **内存管理**: 使用对象池减少 GC 压力
3. **负载均衡**: 实现工作窃取算法
4. **监控指标**: 监控队列长度和响应时间

## 持续构建状态

- [x] 线程池模式 (100%)
- [x] Future/Promise 模式 (100%)
- [x] Actor 模式 (100%)
- [x] 生产者-消费者模式 (100%)
- [ ] 读写锁模式 (0%)
- [ ] 信号量模式 (0%)
- [ ] 屏障模式 (0%)

---

**构建原则**: 严格数学规范，形式化证明，Go语言实现！<(￣︶￣)↗[GO!]
