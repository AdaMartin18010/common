# 02-Go并发编程 (Go Concurrency)

## 目录

1. [概述](#1-概述)
2. [形式化定义](#2-形式化定义)
3. [Goroutine](#3-goroutine)
4. [Channel](#4-channel)
5. [Sync包](#5-sync包)
6. [并发模式](#6-并发模式)
7. [并发安全](#7-并发安全)
8. [性能优化](#8-性能优化)

### 1. 概述

Go语言的并发模型基于CSP（Communicating Sequential Processes）理论，通过goroutine和channel提供简洁而强大的并发编程能力。

#### 1.1 核心概念

**Goroutine**：轻量级线程，由Go运行时管理
**Channel**：类型安全的通信机制
**CSP模型**：通过通信共享内存，而非通过共享内存通信

#### 1.2 并发模型

```go
// Go并发模型
type ConcurrencyModel struct {
    Goroutines []Goroutine
    Channels   []Channel
    SyncPrimitives []SyncPrimitive
}

// Goroutine
type Goroutine struct {
    ID       int
    Function func()
    Status   GoroutineStatus
}

// Channel
type Channel struct {
    Buffer   int
    Type     reflect.Type
    Direction ChannelDirection
}

// 同步原语
type SyncPrimitive interface {
    Lock()
    Unlock()
}
```

### 2. 形式化定义

#### 2.1 CSP模型形式化

**定义 2.1.1** (CSP进程)
CSP进程是一个三元组 ```latex
P = (S, A, T)
```，其中：

- ```latex
S
``` 是状态集合
- ```latex
A
``` 是动作集合
- ```latex
T \subseteq S \times A \times S
``` 是转换关系

**定义 2.1.2** (Channel通信)
Channel通信是一个四元组 ```latex
C = (P_1, P_2, M, \tau)
```，其中：

- ```latex
P_1
``` 是发送进程
- ```latex
P_2
``` 是接收进程
- ```latex
M
``` 是消息集合
- ```latex
\tau: P_1 \times M \rightarrow P_2
``` 是通信函数

**定理 2.1.1** (通信安全性)
对于任意channel ```latex
c
``` 和消息 ```latex
m
```，如果 ```latex
c
``` 是类型安全的，则 ```latex
\tau(p_1, m) = p_2
``` 当且仅当 ```latex
m
``` 的类型与 ```latex
c
``` 的类型匹配。

**证明**：
Go的channel是强类型的，编译时检查确保只有匹配类型的消息才能通过channel传输。
因此，通信安全性成立。```latex
\square
```

#### 2.2 并发安全形式化

**定义 2.2.1** (数据竞争)
数据竞争是一个三元组 ```latex
R = (T_1, T_2, V)
```，其中：

- ```latex
T_1, T_2
``` 是并发线程
- ```latex
V
``` 是共享变量
- ```latex
T_1
``` 和 ```latex
T_2
``` 同时访问 ```latex
V
```，且至少有一个是写操作

**算法 2.2.1** (数据竞争检测)

```text
输入: 程序 P, 执行轨迹 T
输出: 数据竞争集合 R

1. R ← ∅
2. for each 访问 a in T do
3.     for each 并发访问 b in T do
4.         if isDataRace(a, b) then
5.             R ← R ∪ {(a.thread, b.thread, a.variable)}
6.         end if
7.     end for
8. end for
9. return R
```

### 3. Goroutine

#### 3.1 Goroutine基础

```go
// Goroutine创建
func GoroutineBasics() {
    // 1. 基本goroutine
    go func() {
        fmt.Println("Hello from goroutine")
    }()
    
    // 2. 带参数的goroutine
    message := "Hello, World!"
    go func(msg string) {
        fmt.Println(msg)
    }(message)
    
    // 3. 函数调用goroutine
    go printMessage("Async message")
    
    // 4. 等待goroutine完成
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("Goroutine with WaitGroup")
    }()
    wg.Wait()
}

// 函数定义
func printMessage(msg string) {
    fmt.Println(msg)
}

// Goroutine生命周期管理
func GoroutineLifecycle() {
    // 1. 启动多个goroutine
    for i := 0; i < 5; i++ {
        go func(id int) {
            fmt.Printf("Goroutine %d started\n", id)
            time.Sleep(time.Second)
            fmt.Printf("Goroutine %d finished\n", id)
        }(i)
    }
    
    // 2. 等待所有goroutine完成
    time.Sleep(2 * time.Second)
}

// Goroutine池
type GoroutinePool struct {
    workers    int
    tasks      chan func()
    wg         sync.WaitGroup
}

// 创建goroutine池
func NewGoroutinePool(workers int) *GoroutinePool {
    pool := &GoroutinePool{
        workers: workers,
        tasks:   make(chan func(), 100),
    }
    
    // 启动工作goroutine
    for i := 0; i < workers; i++ {
        pool.wg.Add(1)
        go pool.worker()
    }
    
    return pool
}

// 工作goroutine
func (p *GoroutinePool) worker() {
    defer p.wg.Done()
    
    for task := range p.tasks {
        task()
    }
}

// 提交任务
func (p *GoroutinePool) Submit(task func()) {
    p.tasks <- task
}

// 关闭池
func (p *GoroutinePool) Close() {
    close(p.tasks)
    p.wg.Wait()
}
```

#### 3.2 Goroutine调度

```go
// 调度器信息
func SchedulerInfo() {
    // 获取GOMAXPROCS
    maxProcs := runtime.GOMAXPROCS(0)
    fmt.Printf("GOMAXPROCS: %d\n", maxProcs)
    
    // 获取goroutine数量
    numGoroutines := runtime.NumGoroutine()
    fmt.Printf("Number of goroutines: %d\n", numGoroutines)
    
    // 获取CPU数量
    numCPU := runtime.NumCPU()
    fmt.Printf("Number of CPUs: %d\n", numCPU)
}

// 手动调度
func ManualScheduling() {
    // 设置GOMAXPROCS
    runtime.GOMAXPROCS(4)
    
    // 让出CPU时间片
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 1000; j++ {
                if j%100 == 0 {
                    runtime.Gosched() // 让出CPU
                }
            }
            fmt.Printf("Goroutine %d completed\n", id)
        }(i)
    }
    
    time.Sleep(time.Second)
}

// 阻塞检测
func BlockingDetection() {
    // 启动监控goroutine
    go func() {
        for {
            numGoroutines := runtime.NumGoroutine()
            if numGoroutines > 1000 {
                fmt.Printf("Warning: High number of goroutines: %d\n", numGoroutines)
            }
            time.Sleep(time.Second)
        }
    }()
}
```

### 4. Channel

#### 4.1 Channel基础

```go
// Channel创建和使用
func ChannelBasics() {
    // 1. 无缓冲channel
    ch := make(chan int)
    
    // 发送数据
    go func() {
        ch <- 42
    }()
    
    // 接收数据
    value := <-ch
    fmt.Printf("Received: %d\n", value)
    
    // 2. 有缓冲channel
    bufferedCh := make(chan string, 3)
    
    // 发送数据（不会阻塞）
    bufferedCh <- "Hello"
    bufferedCh <- "World"
    bufferedCh <- "Go"
    
    // 接收数据
    for i := 0; i < 3; i++ {
        msg := <-bufferedCh
        fmt.Printf("Received: %s\n", msg)
    }
}

// Channel方向
func ChannelDirections() {
    ch := make(chan int, 5)
    
    // 只发送channel
    go sender(ch)
    
    // 只接收channel
    go receiver(ch)
    
    time.Sleep(time.Second)
}

// 发送函数
func sender(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
        fmt.Printf("Sent: %d\n", i)
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

// 接收函数
func receiver(ch <-chan int) {
    for value := range ch {
        fmt.Printf("Received: %d\n", value)
    }
}

// Channel操作
func ChannelOperations() {
    ch := make(chan int, 2)
    
    // 1. 发送操作
    ch <- 1
    ch <- 2
    
    // 2. 接收操作
    value1 := <-ch
    value2 := <-ch
    
    fmt.Printf("Values: %d, %d\n", value1, value2)
    
    // 3. 关闭channel
    close(ch)
    
    // 4. 检查channel是否关闭
    value, ok := <-ch
    if !ok {
        fmt.Println("Channel is closed")
    }
}
```

#### 4.2 Channel模式

```go
// 1. 生产者-消费者模式
func ProducerConsumer() {
    ch := make(chan int, 10)
    
    // 生产者
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
            fmt.Printf("Produced: %d\n", i)
            time.Sleep(100 * time.Millisecond)
        }
        close(ch)
    }()
    
    // 消费者
    go func() {
        for value := range ch {
            fmt.Printf("Consumed: %d\n", value)
            time.Sleep(200 * time.Millisecond)
        }
    }()
    
    time.Sleep(3 * time.Second)
}

// 2. 工作池模式
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    workers int
}

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Result string
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, 100),
        results: make(chan Result, 100),
        workers: workers,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    for job := range wp.jobs {
        result := wp.processJob(job)
        wp.results <- result
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // 模拟处理
    time.Sleep(100 * time.Millisecond)
    return Result{
        JobID:  job.ID,
        Result: fmt.Sprintf("Processed: %s", job.Data),
    }
}

// 3. 扇入扇出模式
func FanInFanOut() {
    // 扇出：一个channel分发到多个goroutine
    input := make(chan int, 10)
    
    // 启动多个处理器
    output1 := process(input)
    output2 := process(input)
    output3 := process(input)
    
    // 扇入：多个channel合并到一个channel
    merged := merge(output1, output2, output3)
    
    // 发送数据
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // 接收结果
    for result := range merged {
        fmt.Printf("Result: %d\n", result)
    }
}

func process(input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for value := range input {
            output <- value * 2
        }
    }()
    return output
}

func merge(channels ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for value := range c {
                output <- value
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

#### 4.3 Select语句

```go
// Select基础
func SelectBasics() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // 发送数据
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "from ch2"
    }()
    
    // 选择第一个可用的channel
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Printf("Received: %s\n", msg1)
        case msg2 := <-ch2:
            fmt.Printf("Received: %s\n", msg2)
        }
    }
}

// Select with timeout
func SelectWithTimeout() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case result := <-ch:
        fmt.Printf("Got result: %s\n", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}

// Select with default
func SelectWithDefault() {
    ch := make(chan string, 1)
    
    select {
    case ch <- "message":
        fmt.Println("Message sent")
    default:
        fmt.Println("Channel is full, message not sent")
    }
    
    select {
    case msg := <-ch:
        fmt.Printf("Received: %s\n", msg)
    default:
        fmt.Println("No message available")
    }
}

// 非阻塞channel操作
func NonBlockingChannelOps() {
    ch := make(chan int, 1)
    
    // 非阻塞发送
    select {
    case ch <- 1:
        fmt.Println("Sent successfully")
    default:
        fmt.Println("Send would block")
    }
    
    // 非阻塞接收
    select {
    case value := <-ch:
        fmt.Printf("Received: %d\n", value)
    default:
        fmt.Println("Receive would block")
    }
}
```

### 5. Sync包

#### 5.1 同步原语

```go
// Mutex
func MutexExample() {
    var mu sync.Mutex
    var counter int
    
    var wg sync.WaitGroup
    wg.Add(100)
    
    for i := 0; i < 100; i++ {
        go func() {
            defer wg.Done()
            
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", counter)
}

// RWMutex
func RWMutexExample() {
    var mu sync.RWMutex
    var data map[string]string = make(map[string]string)
    
    // 写操作
    go func() {
        for i := 0; i < 10; i++ {
            mu.Lock()
            data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
            mu.Unlock()
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // 读操作
    for i := 0; i < 5; i++ {
        go func(id int) {
            for j := 0; j < 10; j++ {
                mu.RLock()
                value := data[fmt.Sprintf("key%d", j)]
                mu.RUnlock()
                fmt.Printf("Reader %d: key%d = %s\n", id, j, value)
                time.Sleep(50 * time.Millisecond)
            }
        }(i)
    }
    
    time.Sleep(2 * time.Second)
}

// WaitGroup
func WaitGroupExample() {
    var wg sync.WaitGroup
    
    // 启动多个goroutine
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d started\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Goroutine %d finished\n", id)
        }(i)
    }
    
    // 等待所有goroutine完成
    wg.Wait()
    fmt.Println("All goroutines completed")
}

// Once
func OnceExample() {
    var once sync.Once
    var initialized bool
    
    // 多次调用，但只执行一次
    for i := 0; i < 5; i++ {
        go func(id int) {
            once.Do(func() {
                initialized = true
                fmt.Printf("Initialized by goroutine %d\n", id)
            })
        }(i)
    }
    
    time.Sleep(time.Second)
    fmt.Printf("Initialized: %t\n", initialized)
}

// Cond
func CondExample() {
    var mu sync.Mutex
    cond := sync.NewCond(&mu)
    var ready bool
    
    // 等待条件
    go func() {
        mu.Lock()
        for !ready {
            cond.Wait()
        }
        fmt.Println("Condition met!")
        mu.Unlock()
    }()
    
    // 设置条件
    time.Sleep(time.Second)
    mu.Lock()
    ready = true
    cond.Signal()
    mu.Unlock()
    
    time.Sleep(time.Second)
}
```

#### 5.2 原子操作

```go
// 原子操作
func AtomicOperations() {
    var counter int64
    
    var wg sync.WaitGroup
    wg.Add(100)
    
    for i := 0; i < 100; i++ {
        go func() {
            defer wg.Done()
            
            // 原子递增
            atomic.AddInt64(&counter, 1)
        }()
    }
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", atomic.LoadInt64(&counter))
}

// 原子比较和交换
func AtomicCAS() {
    var value int64 = 0
    
    // 比较并交换
    swapped := atomic.CompareAndSwapInt64(&value, 0, 1)
    fmt.Printf("CAS(0, 1): %t, value: %d\n", swapped, atomic.LoadInt64(&value))
    
    swapped = atomic.CompareAndSwapInt64(&value, 0, 2)
    fmt.Printf("CAS(0, 2): %t, value: %d\n", swapped, atomic.LoadInt64(&value))
    
    swapped = atomic.CompareAndSwapInt64(&value, 1, 2)
    fmt.Printf("CAS(1, 2): %t, value: %d\n", swapped, atomic.LoadInt64(&value))
}

// 原子指针
func AtomicPointer() {
    var ptr atomic.Value
    
    // 存储指针
    data := "hello"
    ptr.Store(&data)
    
    // 加载指针
    if loaded := ptr.Load(); loaded != nil {
        strPtr := loaded.(*string)
        fmt.Printf("Loaded: %s\n", *strPtr)
    }
}
```

### 6. 并发模式

#### 6.1 常见并发模式

```go
// 1. 管道模式
func Pipeline() {
    // 创建管道
    numbers := make(chan int, 10)
    squares := make(chan int, 10)
    results := make(chan int, 10)
    
    // 生成数字
    go func() {
        defer close(numbers)
        for i := 0; i < 10; i++ {
            numbers <- i
        }
    }()
    
    // 计算平方
    go func() {
        defer close(squares)
        for n := range numbers {
            squares <- n * n
        }
    }()
    
    // 过滤偶数
    go func() {
        defer close(results)
        for s := range squares {
            if s%2 == 0 {
                results <- s
            }
        }
    }()
    
    // 收集结果
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}

// 2. 工作窃取模式
type WorkStealingPool struct {
    workers []*Worker
    tasks   chan Task
}

type Worker struct {
    id       int
    tasks    []Task
    pool     *WorkStealingPool
    stealing bool
}

type Task struct {
    ID   int
    Work func()
}

func NewWorkStealingPool(workerCount int) *WorkStealingPool {
    pool := &WorkStealingPool{
        workers: make([]*Worker, workerCount),
        tasks:   make(chan Task, 1000),
    }
    
    for i := 0; i < workerCount; i++ {
        pool.workers[i] = &Worker{
            id:    i,
            tasks: make([]Task, 0),
            pool:  pool,
        }
        go pool.workers[i].run()
    }
    
    return pool
}

func (w *Worker) run() {
    for {
        // 处理自己的任务
        if len(w.tasks) > 0 {
            task := w.tasks[len(w.tasks)-1]
            w.tasks = w.tasks[:len(w.tasks)-1]
            task.Work()
            continue
        }
        
        // 窃取其他worker的任务
        w.steal()
        
        // 如果没有任务，等待
        time.Sleep(1 * time.Millisecond)
    }
}

func (w *Worker) steal() {
    for _, other := range w.pool.workers {
        if other.id == w.id || len(other.tasks) == 0 {
            continue
        }
        
        // 窃取一半的任务
        stealCount := len(other.tasks) / 2
        if stealCount > 0 {
            w.tasks = append(w.tasks, other.tasks[:stealCount]...)
            other.tasks = other.tasks[stealCount:]
        }
    }
}

// 3. 发布订阅模式
type Publisher struct {
    subscribers map[string][]chan interface{}
    mu          sync.RWMutex
}

func NewPublisher() *Publisher {
    return &Publisher{
        subscribers: make(map[string][]chan interface{}),
    }
}

func (p *Publisher) Subscribe(topic string) <-chan interface{} {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    ch := make(chan interface{}, 1)
    p.subscribers[topic] = append(p.subscribers[topic], ch)
    return ch
}

func (p *Publisher) Publish(topic string, message interface{}) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    for _, ch := range p.subscribers[topic] {
        select {
        case ch <- message:
        default:
            // Channel is full, skip
        }
    }
}

func (p *Publisher) Unsubscribe(topic string, ch <-chan interface{}) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    subscribers := p.subscribers[topic]
    for i, subscriber := range subscribers {
        if subscriber == ch {
            p.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
            close(subscriber)
            break
        }
    }
}
```

#### 6.2 高级并发模式

```go
// 1. 上下文模式
func ContextPattern() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // 启动工作goroutine
    go func() {
        select {
        case <-ctx.Done():
            fmt.Println("Context cancelled")
        case <-time.After(3 * time.Second):
            fmt.Println("Work completed")
        }
    }()
    
    // 等待上下文取消
    <-ctx.Done()
    fmt.Println("Main: context cancelled")
}

// 2. 错误组模式
func ErrorGroupPattern() {
    var eg errgroup.Group
    
    // 启动多个任务
    for i := 0; i < 3; i++ {
        i := i
        eg.Go(func() error {
            if i == 1 {
                return fmt.Errorf("task %d failed", i)
            }
            fmt.Printf("Task %d completed\n", i)
            return nil
        })
    }
    
    // 等待所有任务完成
    if err := eg.Wait(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

// 3. 限制并发数
func LimitConcurrency() {
    semaphore := make(chan struct{}, 3) // 最多3个并发
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            semaphore <- struct{}{} // 获取信号量
            defer func() { <-semaphore }() // 释放信号量
            
            fmt.Printf("Task %d started\n", id)
            time.Sleep(time.Second)
            fmt.Printf("Task %d finished\n", id)
        }(i)
    }
    
    wg.Wait()
}
```

### 7. 并发安全

#### 7.1 数据竞争检测

```go
// 数据竞争示例
func DataRaceExample() {
    var counter int
    var wg sync.WaitGroup
    wg.Add(2)
    
    // 两个goroutine同时访问counter
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            counter++ // 数据竞争！
        }
    }()
    
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            counter++ // 数据竞争！
        }
    }()
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", counter)
}

// 修复数据竞争
func FixedDataRaceExample() {
    var counter int
    var mu sync.Mutex
    var wg sync.WaitGroup
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    }()
    
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    }()
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", counter)
}

// 使用原子操作
func AtomicDataRaceExample() {
    var counter int64
    var wg sync.WaitGroup
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            atomic.AddInt64(&counter, 1)
        }
    }()
    
    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            atomic.AddInt64(&counter, 1)
        }
    }()
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", atomic.LoadInt64(&counter))
}
```

#### 7.2 死锁检测

```go
// 死锁示例
func DeadlockExample() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        time.Sleep(100 * time.Millisecond)
        mu2.Lock()
        mu2.Unlock()
        mu1.Unlock()
    }()
    
    go func() {
        mu2.Lock()
        time.Sleep(100 * time.Millisecond)
        mu1.Lock()
        mu1.Unlock()
        mu2.Unlock()
    }()
    
    time.Sleep(1 * time.Second)
}

// 避免死锁
func AvoidDeadlock() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        defer mu1.Unlock()
        time.Sleep(100 * time.Millisecond)
        mu2.Lock()
        defer mu2.Unlock()
        fmt.Println("Goroutine 1 completed")
    }()
    
    go func() {
        mu1.Lock()
        defer mu1.Unlock()
        time.Sleep(100 * time.Millisecond)
        mu2.Lock()
        defer mu2.Unlock()
        fmt.Println("Goroutine 2 completed")
    }()
    
    time.Sleep(1 * time.Second)
}
```

### 8. 性能优化

#### 8.1 并发性能优化

```go
// 1. 减少锁竞争
func ReduceLockContention() {
    // 使用分片锁
    type ShardedCounter struct {
        shards []*CounterShard
    }
    
    type CounterShard struct {
        mu      sync.Mutex
        counter int64
    }
    
    counter := &ShardedCounter{
        shards: make([]*CounterShard, 16),
    }
    
    for i := range counter.shards {
        counter.shards[i] = &CounterShard{}
    }
    
    // 根据key选择分片
    hash := func(key string) int {
        h := 0
        for _, c := range key {
            h = 31*h + int(c)
        }
        return h % len(counter.shards)
    }
    
    // 并发更新
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", id)
            shard := counter.shards[hash(key)]
            shard.mu.Lock()
            shard.counter++
            shard.mu.Unlock()
        }(i)
    }
    
    wg.Wait()
}

// 2. 内存池
func MemoryPool() {
    pool := sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            // 从池中获取buffer
            buffer := pool.Get().([]byte)
            defer pool.Put(buffer)
            
            // 使用buffer
            copy(buffer, []byte("hello"))
        }()
    }
    
    wg.Wait()
}

// 3. 批量处理
func BatchProcessing() {
    const batchSize = 100
    jobs := make(chan int, 1000)
    results := make(chan int, 1000)
    
    // 启动工作goroutine
    for i := 0; i < 4; i++ {
        go func() {
            batch := make([]int, 0, batchSize)
            
            for job := range jobs {
                batch = append(batch, job)
                
                if len(batch) >= batchSize {
                    // 批量处理
                    processBatch(batch, results)
                    batch = batch[:0]
                }
            }
            
            // 处理剩余的任务
            if len(batch) > 0 {
                processBatch(batch, results)
            }
        }()
    }
    
    // 发送任务
    go func() {
        for i := 0; i < 1000; i++ {
            jobs <- i
        }
        close(jobs)
    }()
    
    // 收集结果
    go func() {
        for i := 0; i < 1000; i++ {
            <-results
        }
        close(results)
    }()
}

func processBatch(batch []int, results chan<- int) {
    // 批量处理逻辑
    for _, item := range batch {
        results <- item * 2
    }
}
```

### 总结

本模块提供了完整的Go并发编程实现，包括：

1. **形式化定义**：基于CSP模型的并发系统定义和安全性证明
2. **Goroutine**：轻量级线程的创建、管理和调度
3. **Channel**：类型安全的通信机制和模式
4. **Sync包**：同步原语和原子操作
5. **并发模式**：常见和高级并发模式
6. **并发安全**：数据竞争检测和死锁避免
7. **性能优化**：并发性能优化技巧

该实现遵循了Go并发编程的最佳实践，提供了安全、高效、可维护的并发解决方案。
