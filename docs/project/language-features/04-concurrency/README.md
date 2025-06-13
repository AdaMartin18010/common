# Golang 并发编程详解

## 概述

Golang的并发模型基于CSP（Communicating Sequential Processes）理论，通过goroutines和channels提供了简洁而强大的并发编程能力。这种设计使得并发编程变得简单、安全且高效。

## 核心概念

### 1. Goroutines

Goroutines是Go语言中的轻量级线程，由Go运行时管理。

#### 基本语法

```go
// 启动一个goroutine
go functionName()

// 匿名函数goroutine
go func() {
    // 执行代码
}()
```

#### 特点

- **轻量级**: 初始栈大小仅2KB，远小于传统线程
- **可扩展**: 可以轻松创建数百万个goroutine
- **自动调度**: 由Go运行时自动调度到系统线程
- **简单易用**: 使用`go`关键字即可启动

#### 示例代码

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 5; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // 启动两个goroutine
    go sayHello("Alice")
    go sayHello("Bob")
    
    // 等待goroutine完成
    time.Sleep(1 * time.Second)
}
```

### 2. Channels

Channels是goroutine之间通信的主要机制，遵循"不要通过共享内存来通信，而要通过通信来共享内存"的原则。

#### 基本语法

```go
// 创建channel
ch := make(chan Type, bufferSize)

// 发送数据
ch <- value

// 接收数据
value := <-ch
```

#### Channel类型

**无缓冲Channel**

```go
ch := make(chan int) // 无缓冲
```

**有缓冲Channel**

```go
ch := make(chan int, 10) // 缓冲大小为10
```

#### 示例代码

```go
package main

import (
    "fmt"
    "time"
)

func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        fmt.Printf("Producing: %d\n", i)
        ch <- i
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int) {
    for value := range ch {
        fmt.Printf("Consuming: %d\n", value)
        time.Sleep(200 * time.Millisecond)
    }
}

func main() {
    ch := make(chan int, 3)
    
    go producer(ch)
    go consumer(ch)
    
    time.Sleep(2 * time.Second)
}
```

### 3. Select语句

Select语句用于在多个channel操作中进行选择，类似于switch语句。

#### 基本语法

```go
select {
case <-ch1:
    // 处理ch1的数据
case ch2 <- value:
    // 向ch2发送数据
case <-time.After(timeout):
    // 超时处理
default:
    // 默认处理
}
```

#### 示例代码

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func generator(name string, ch chan<- string) {
    for {
        time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
        ch <- fmt.Sprintf("Message from %s", name)
    }
}

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go generator("A", ch1)
    go generator("B", ch2)
    
    for i := 0; i < 10; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        case <-time.After(500 * time.Millisecond):
            fmt.Println("Timeout!")
        }
    }
}
```

### 4. 同步原语

#### Mutex (互斥锁)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) GetCount() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func main() {
    counter := &Counter{}
    
    // 启动多个goroutine
    for i := 0; i < 10; i++ {
        go func() {
            for j := 0; j < 1000; j++ {
                counter.Increment()
            }
        }()
    }
    
    time.Sleep(2 * time.Second)
    fmt.Printf("Final count: %d\n", counter.GetCount())
}
```

#### WaitGroup

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
    fmt.Println("All workers completed")
}
```

#### Once

```go
package main

import (
    "fmt"
    "sync"
)

var (
    instance *Singleton
    once     sync.Once
)

type Singleton struct {
    data string
}

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "Initialized"}
        fmt.Println("Singleton initialized")
    })
    return instance
}

func main() {
    for i := 0; i < 5; i++ {
        go func(id int) {
            instance := GetInstance()
            fmt.Printf("Goroutine %d got instance: %v\n", id, instance)
        }(i)
    }
    
    // 等待所有goroutine完成
    time.Sleep(1 * time.Second)
}
```

## 并发模式

### 1. Worker Pool模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Millisecond * 500) // 模拟工作
        results <- job * 2
    }
}

func main() {
    const numJobs = 10
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    
    var wg sync.WaitGroup
    
    // 启动workers
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }
    
    // 发送jobs
    for i := 1; i <= numJobs; i++ {
        jobs <- i
    }
    close(jobs)
    
    // 等待所有workers完成
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集结果
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

### 2. Pipeline模式

```go
package main

import (
    "fmt"
    "math"
)

func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func filter(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n > 10 {
                out <- n
            }
        }
    }()
    return out
}

func main() {
    // 创建pipeline: generator -> square -> filter
    nums := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squares := square(nums)
    filtered := filter(squares)
    
    for result := range filtered {
        fmt.Printf("Filtered result: %d\n", result)
    }
}
```

### 3. Fan-out/Fan-in模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func producer() <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; i < 10; i++ {
            out <- i
            time.Sleep(100 * time.Millisecond)
        }
    }()
    return out
}

func worker(id int, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            fmt.Printf("Worker %d processing %d\n", id, n)
            time.Sleep(200 * time.Millisecond)
            out <- n * n
        }
    }()
    return out
}

func merge(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // 启动goroutine从每个channel读取
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    
    // 等待所有goroutine完成
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

func main() {
    // Fan-out: 一个producer分发到多个workers
    input := producer()
    
    // 创建多个workers
    worker1 := worker(1, input)
    worker2 := worker(2, input)
    worker3 := worker(3, input)
    
    // Fan-in: 合并多个workers的结果
    results := merge(worker1, worker2, worker3)
    
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

## Context包

Context包用于在goroutine之间传递请求范围的值、取消信号和截止时间。

### 基本用法

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %s cancelled: %v\n", name, ctx.Err())
            return
        default:
            fmt.Printf("Worker %s working...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // 创建带超时的context
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go worker(ctx, "A")
    go worker(ctx, "B")
    
    // 等待context取消
    <-ctx.Done()
    fmt.Println("Main: context cancelled")
}
```

### 传递值

```go
package main

import (
    "context"
    "fmt"
)

func processRequest(ctx context.Context) {
    // 从context获取值
    userID, ok := ctx.Value("userID").(string)
    if ok {
        fmt.Printf("Processing request for user: %s\n", userID)
    }
    
    // 传递值给子context
    childCtx := context.WithValue(ctx, "requestID", "req-123")
    processSubRequest(childCtx)
}

func processSubRequest(ctx context.Context) {
    userID := ctx.Value("userID")
    requestID := ctx.Value("requestID")
    fmt.Printf("Sub-request - User: %v, Request: %v\n", userID, requestID)
}

func main() {
    ctx := context.WithValue(context.Background(), "userID", "user-456")
    processRequest(ctx)
}
```

## 最佳实践

### 1. 避免Goroutine泄漏

```go
// 错误示例
func badExample() {
    go func() {
        for {
            // 无限循环，没有退出条件
            time.Sleep(time.Second)
        }
    }()
}

// 正确示例
func goodExample() {
    done := make(chan bool)
    go func() {
        defer close(done)
        for {
            select {
            case <-done:
                return
            default:
                time.Sleep(time.Second)
            }
        }
    }()
    
    // 在适当的时候关闭
    time.Sleep(5 * time.Second)
    close(done)
}
```

### 2. 合理使用Channel缓冲

```go
// 无缓冲channel - 同步通信
ch := make(chan int)

// 有缓冲channel - 异步通信
ch := make(chan int, 10)

// 根据实际需求选择
func example() {
    // 需要同步时使用无缓冲
    syncCh := make(chan int)
    
    // 需要异步时使用有缓冲
    asyncCh := make(chan int, 100)
}
```

### 3. 使用Context进行取消

```go
func longRunningTask(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 执行任务
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    go func() {
        if err := longRunningTask(ctx); err != nil {
            fmt.Printf("Task failed: %v\n", err)
        }
    }()
    
    // 等待任务完成或超时
    <-ctx.Done()
}
```

## 性能优化

### 1. Goroutine池

```go
type Pool struct {
    workers int
    tasks   chan func()
    wg      sync.WaitGroup
}

func NewPool(workers int) *Pool {
    p := &Pool{
        workers: workers,
        tasks:   make(chan func()),
    }
    
    for i := 0; i < workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
    
    return p
}

func (p *Pool) worker() {
    defer p.wg.Done()
    for task := range p.tasks {
        task()
    }
}

func (p *Pool) Submit(task func()) {
    p.tasks <- task
}

func (p *Pool) Close() {
    close(p.tasks)
    p.wg.Wait()
}
```

### 2. 批量处理

```go
func batchProcessor(items []int, batchSize int) []int {
    results := make([]int, 0, len(items))
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }
        
        batch := items[i:end]
        wg.Add(1)
        
        go func(batch []int) {
            defer wg.Done()
            
            // 处理批次
            batchResults := make([]int, len(batch))
            for j, item := range batch {
                batchResults[j] = item * 2 // 模拟处理
            }
            
            // 安全地添加结果
            mu.Lock()
            results = append(results, batchResults...)
            mu.Unlock()
        }(batch)
    }
    
    wg.Wait()
    return results
}
```

## 常见陷阱

### 1. 闭包变量捕获

```go
// 错误示例
func wrongClosure() {
    for i := 0; i < 3; i++ {
        go func() {
            fmt.Println(i) // 所有goroutine都会打印相同的值
        }()
    }
}

// 正确示例
func correctClosure() {
    for i := 0; i < 3; i++ {
        go func(id int) {
            fmt.Println(id) // 每个goroutine打印不同的值
        }(i)
    }
}
```

### 2. Channel关闭

```go
// 错误示例 - 重复关闭
func wrongClose() {
    ch := make(chan int)
    close(ch)
    close(ch) // panic: close of closed channel
}

// 正确示例 - 使用sync.Once
func correctClose() {
    ch := make(chan int)
    var once sync.Once
    
    closeFunc := func() {
        close(ch)
    }
    
    once.Do(closeFunc)
    once.Do(closeFunc) // 安全，不会重复关闭
}
```

## 2025年改进

### 1. 更高效的调度器

- 改进的GOMAXPROCS管理
- 更好的负载均衡
- 减少调度开销

### 2. 改进的Channel性能

- 更快的channel操作
- 更好的内存使用
- 优化的select性能

### 3. 更好的并发工具

- 新的同步原语
- 改进的context包
- 更好的调试工具

## 总结

Golang的并发编程模型简洁而强大，通过goroutines和channels提供了高效的并发解决方案。掌握这些核心概念和模式，可以构建出高性能、可维护的并发应用程序。

---

*最后更新时间: 2025年1月*
*文档版本: v1.0*
