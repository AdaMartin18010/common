# 04-并发编程 (Concurrent Programming)

## 1. 概述

### 1.1 并发基础

**并发编程** 是同时执行多个计算任务的技术，Go语言通过goroutine和channel提供优雅的并发支持。

**形式化定义**：
设 ```latex
$T_1, T_2, ..., T_n$
``` 为任务集合，并发执行函数：
$```latex
$\text{Concurrent}(T_1, T_2, ..., T_n) = \text{parallel}(T_1) \parallel \text{parallel}(T_2) \parallel ... \parallel \text{parallel}(T_n)$
```$

### 1.2 Go并发模型

**CSP (Communicating Sequential Processes)** 模型：

- 通过通信共享内存，而不是通过共享内存通信
- goroutine作为轻量级线程
- channel作为通信机制

## 2. Goroutine

### 2.1 理论基础

**Goroutine** 是Go语言的轻量级线程，由Go运行时管理。

**特点**：

- 初始栈大小：2KB
- 可创建数量：理论上无限制
- 调度：由Go运行时调度器管理

### 2.2 基础实现

```go
package concurrency

import (
 "fmt"
 "sync"
 "time"
)

// BasicGoroutine 基础goroutine示例
func BasicGoroutine() {
 fmt.Println("主goroutine开始")
 
 // 启动goroutine
 go func() {
  fmt.Println("子goroutine执行")
  time.Sleep(100 * time.Millisecond)
  fmt.Println("子goroutine完成")
 }()
 
 // 等待一段时间让子goroutine执行
 time.Sleep(200 * time.Millisecond)
 fmt.Println("主goroutine结束")
}

// GoroutineWithParameter 带参数的goroutine
func GoroutineWithParameter(id int, message string) {
 go func(id int, msg string) {
  fmt.Printf("Goroutine %d: %s\n", id, msg)
  time.Sleep(50 * time.Millisecond)
 }(id, message)
}

// MultipleGoroutines 多个goroutine
func MultipleGoroutines(count int) {
 for i := 0; i < count; i++ {
  go func(id int) {
   fmt.Printf("Goroutine %d 开始\n", id)
   time.Sleep(time.Duration(id*10) * time.Millisecond)
   fmt.Printf("Goroutine %d 结束\n", id)
  }(i)
 }
 
 // 等待所有goroutine完成
 time.Sleep(time.Duration(count*10+100) * time.Millisecond)
}

// GoroutineWithReturn 带返回值的goroutine
func GoroutineWithReturn() int {
 result := make(chan int, 1)
 
 go func() {
  // 模拟计算
  time.Sleep(100 * time.Millisecond)
  result <- 42
 }()
 
 return <-result
}
```

### 2.3 高级goroutine模式

```go
// WorkerPool 工作池模式
type WorkerPool struct {
 workers    int
 taskQueue  chan func()
 wg         sync.WaitGroup
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int) *WorkerPool {
 pool := &WorkerPool{
  workers:   workers,
  taskQueue: make(chan func(), workers*2),
 }
 
 // 启动工作goroutine
 for i := 0; i < workers; i++ {
  pool.wg.Add(1)
  go pool.worker(i)
 }
 
 return pool
}

// worker 工作goroutine
func (wp *WorkerPool) worker(id int) {
 defer wp.wg.Done()
 
 for task := range wp.taskQueue {
  fmt.Printf("Worker %d 执行任务\n", id)
  task()
 }
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task func()) {
 wp.taskQueue <- task
}

// Close 关闭工作池
func (wp *WorkerPool) Close() {
 close(wp.taskQueue)
 wp.wg.Wait()
}

// Pipeline 管道模式
type Pipeline struct {
 stages []chan interface{}
}

// NewPipeline 创建管道
func NewPipeline(stageCount int) *Pipeline {
 pipeline := &Pipeline{
  stages: make([]chan interface{}, stageCount),
 }
 
 for i := 0; i < stageCount; i++ {
  pipeline.stages[i] = make(chan interface{}, 10)
 }
 
 return pipeline
}

// AddStage 添加处理阶段
func (p *Pipeline) AddStage(stageID int, processor func(interface{}) interface{}) {
 if stageID >= len(p.stages) {
  return
 }
 
 go func() {
  for data := range p.stages[stageID] {
   result := processor(data)
   if stageID+1 < len(p.stages) {
    p.stages[stageID+1] <- result
   }
  }
  if stageID+1 < len(p.stages) {
   close(p.stages[stageID+1])
  }
 }()
}

// Send 发送数据到管道
func (p *Pipeline) Send(data interface{}) {
 if len(p.stages) > 0 {
  p.stages[0] <- data
 }
}

// Close 关闭管道
func (p *Pipeline) Close() {
 if len(p.stages) > 0 {
  close(p.stages[0])
 }
}
```

## 3. Channel

### 3.1 理论基础

**Channel** 是Go语言的通信机制，提供goroutine之间的同步和通信。

**类型**：

- 无缓冲channel：同步通信
- 有缓冲channel：异步通信

### 3.2 基础channel操作

```go
// BasicChannel 基础channel操作
func BasicChannel() {
 // 创建无缓冲channel
 ch := make(chan int)
 
 // 发送数据
 go func() {
  fmt.Println("发送数据: 42")
  ch <- 42
 }()
 
 // 接收数据
 data := <-ch
 fmt.Printf("接收数据: %d\n", data)
}

// BufferedChannel 有缓冲channel
func BufferedChannel() {
 // 创建有缓冲channel
 ch := make(chan int, 3)
 
 // 发送数据（不会阻塞）
 ch <- 1
 ch <- 2
 ch <- 3
 
 // 接收数据
 fmt.Println(<-ch) // 1
 fmt.Println(<-ch) // 2
 fmt.Println(<-ch) // 3
}

// ChannelDirection 单向channel
func ChannelDirection() {
 // 双向channel
 ch := make(chan int)
 
 // 发送专用channel
 sendOnly := (chan<- int)(ch)
 
 // 接收专用channel
 receiveOnly := (<-chan int)(ch)
 
 go func() {
  sendOnly <- 42
  close(sendOnly)
 }()
 
 data := <-receiveOnly
 fmt.Printf("接收数据: %d\n", data)
}

// SelectStatement select语句
func SelectStatement() {
 ch1 := make(chan int)
 ch2 := make(chan string)
 
 go func() {
  time.Sleep(100 * time.Millisecond)
  ch1 <- 42
 }()
 
 go func() {
  time.Sleep(50 * time.Millisecond)
  ch2 <- "hello"
 }()
 
 // 使用select处理多个channel
 for i := 0; i < 2; i++ {
  select {
  case data := <-ch1:
   fmt.Printf("从ch1接收: %d\n", data)
  case data := <-ch2:
   fmt.Printf("从ch2接收: %s\n", data)
  case <-time.After(200 * time.Millisecond):
   fmt.Println("超时")
  }
 }
}
```

### 3.3 高级channel模式

```go
// FanOut 扇出模式
func FanOut(input <-chan int, workers int) []<-chan int {
 outputs := make([]<-chan int, workers)
 
 for i := 0; i < workers; i++ {
  output := make(chan int)
  outputs[i] = output
  
  go func(out chan<- int) {
   defer close(out)
   for data := range input {
    // 模拟处理
    time.Sleep(10 * time.Millisecond)
    out <- data * 2
   }
  }(output)
 }
 
 return outputs
}

// FanIn 扇入模式
func FanIn(inputs ...<-chan int) <-chan int {
 output := make(chan int)
 
 var wg sync.WaitGroup
 wg.Add(len(inputs))
 
 for _, input := range inputs {
  go func(in <-chan int) {
   defer wg.Done()
   for data := range in {
    output <- data
   }
  }(input)
 }
 
 go func() {
  wg.Wait()
  close(output)
 }()
 
 return output
}

// Generator 生成器模式
func Generator(values ...int) <-chan int {
 ch := make(chan int)
 
 go func() {
  defer close(ch)
  for _, v := range values {
   ch <- v
  }
 }()
 
 return ch
}

// Filter 过滤器模式
func Filter(input <-chan int, predicate func(int) bool) <-chan int {
 output := make(chan int)
 
 go func() {
  defer close(output)
  for data := range input {
   if predicate(data) {
    output <- data
   }
  }
 }()
 
 return output
}

// Map 映射模式
func Map(input <-chan int, mapper func(int) int) <-chan int {
 output := make(chan int)
 
 go func() {
  defer close(output)
  for data := range input {
   output <- mapper(data)
  }
 }()
 
 return output
}
```

## 4. 同步原语

### 4.1 Mutex

```go
// SafeCounter 线程安全计数器
type SafeCounter struct {
 mu    sync.Mutex
 count int
}

// Increment 增加计数
func (sc *SafeCounter) Increment() {
 sc.mu.Lock()
 defer sc.mu.Unlock()
 sc.count++
}

// Get 获取计数
func (sc *SafeCounter) Get() int {
 sc.mu.Lock()
 defer sc.mu.Unlock()
 return sc.count
}

// RWMutex 读写锁
type SafeMap struct {
 mu   sync.RWMutex
 data map[string]interface{}
}

// NewSafeMap 创建安全map
func NewSafeMap() *SafeMap {
 return &SafeMap{
  data: make(map[string]interface{}),
 }
}

// Set 设置值
func (sm *SafeMap) Set(key string, value interface{}) {
 sm.mu.Lock()
 defer sm.mu.Unlock()
 sm.data[key] = value
}

// Get 获取值
func (sm *SafeMap) Get(key string) (interface{}, bool) {
 sm.mu.RLock()
 defer sm.mu.RUnlock()
 value, exists := sm.data[key]
 return value, exists
}
```

### 4.2 WaitGroup

```go
// WaitGroupExample WaitGroup示例
func WaitGroupExample() {
 var wg sync.WaitGroup
 
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   fmt.Printf("Goroutine %d 开始\n", id)
   time.Sleep(time.Duration(id*100) * time.Millisecond)
   fmt.Printf("Goroutine %d 结束\n", id)
  }(i)
 }
 
 wg.Wait()
 fmt.Println("所有goroutine完成")
}

// WaitGroupWithResult 带结果的WaitGroup
func WaitGroupWithResult() []int {
 var wg sync.WaitGroup
 results := make([]int, 5)
 
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   results[id] = id * 2
  }(i)
 }
 
 wg.Wait()
 return results
}
```

### 4.3 Cond

```go
// ConditionVariable 条件变量示例
type ConditionVariable struct {
 mu    sync.Mutex
 cond  *sync.Cond
 ready bool
}

// NewConditionVariable 创建条件变量
func NewConditionVariable() *ConditionVariable {
 cv := &ConditionVariable{}
 cv.cond = sync.NewCond(&cv.mu)
 return cv
}

// Wait 等待条件
func (cv *ConditionVariable) Wait() {
 cv.mu.Lock()
 defer cv.mu.Unlock()
 
 for !cv.ready {
  cv.cond.Wait()
 }
}

// Signal 发送信号
func (cv *ConditionVariable) Signal() {
 cv.mu.Lock()
 defer cv.mu.Unlock()
 
 cv.ready = true
 cv.cond.Signal()
}

// Broadcast 广播信号
func (cv *ConditionVariable) Broadcast() {
 cv.mu.Lock()
 defer cv.mu.Unlock()
 
 cv.ready = true
 cv.cond.Broadcast()
}
```

### 4.4 Once

```go
// Singleton 单例模式
type Singleton struct {
 data string
}

var (
 instance *Singleton
 once     sync.Once
)

// GetInstance 获取单例实例
func GetInstance() *Singleton {
 once.Do(func() {
  instance = &Singleton{
   data: "单例数据",
  }
 })
 return instance
}

// LazyInitialization 延迟初始化
type LazyInitialization struct {
 mu     sync.Mutex
 once   sync.Once
 value  interface{}
}

// GetValue 获取值（延迟初始化）
func (li *LazyInitialization) GetValue() interface{} {
 li.once.Do(func() {
  li.mu.Lock()
  defer li.mu.Unlock()
  
  // 模拟昂贵的初始化
  time.Sleep(100 * time.Millisecond)
  li.value = "延迟初始化的值"
 })
 
 return li.value
}
```

## 5. 原子操作

### 5.1 基础原子操作

```go
import "sync/atomic"

// AtomicCounter 原子计数器
type AtomicCounter struct {
 value int64
}

// Increment 原子增加
func (ac *AtomicCounter) Increment() {
 atomic.AddInt64(&ac.value, 1)
}

// Get 原子获取值
func (ac *AtomicCounter) Get() int64 {
 return atomic.LoadInt64(&ac.value)
}

// CompareAndSwap 比较并交换
func (ac *AtomicCounter) CompareAndSwap(old, new int64) bool {
 return atomic.CompareAndSwapInt64(&ac.value, old, new)
}

// AtomicPointer 原子指针
type AtomicPointer struct {
 ptr atomic.Value
}

// Store 存储指针
func (ap *AtomicPointer) Store(value interface{}) {
 ap.ptr.Store(value)
}

// Load 加载指针
func (ap *AtomicPointer) Load() interface{} {
 return ap.ptr.Load()
}
```

## 6. 并发模式

### 6.1 生产者-消费者模式

```go
// ProducerConsumer 生产者-消费者模式
func ProducerConsumer() {
 ch := make(chan int, 10)
 
 // 生产者
 go func() {
  for i := 0; i < 10; i++ {
   fmt.Printf("生产: %d\n", i)
   ch <- i
   time.Sleep(50 * time.Millisecond)
  }
  close(ch)
 }()
 
 // 消费者
 go func() {
  for data := range ch {
   fmt.Printf("消费: %d\n", data)
   time.Sleep(100 * time.Millisecond)
  }
 }()
 
 time.Sleep(2 * time.Second)
}
```

### 6.2 工作窃取模式

```go
// WorkStealing 工作窃取模式
type WorkStealing struct {
 queues []*Deque
 workers int
}

// Deque 双端队列
type Deque struct {
 mu    sync.Mutex
 items []interface{}
}

// PushBack 从后端推入
func (d *Deque) PushBack(item interface{}) {
 d.mu.Lock()
 defer d.mu.Unlock()
 d.items = append(d.items, item)
}

// PopBack 从后端弹出
func (d *Deque) PopBack() (interface{}, bool) {
 d.mu.Lock()
 defer d.mu.Unlock()
 
 if len(d.items) == 0 {
  return nil, false
 }
 
 item := d.items[len(d.items)-1]
 d.items = d.items[:len(d.items)-1]
 return item, true
}

// PopFront 从前端弹出
func (d *Deque) PopFront() (interface{}, bool) {
 d.mu.Lock()
 defer d.mu.Unlock()
 
 if len(d.items) == 0 {
  return nil, false
 }
 
 item := d.items[0]
 d.items = d.items[1:]
 return item, true
}

// NewWorkStealing 创建工作窃取调度器
func NewWorkStealing(workers int) *WorkStealing {
 ws := &WorkStealing{
  queues:  make([]*Deque, workers),
  workers: workers,
 }
 
 for i := 0; i < workers; i++ {
  ws.queues[i] = &Deque{}
 }
 
 return ws
}

// Submit 提交任务
func (ws *WorkStealing) Submit(workerID int, task interface{}) {
 ws.queues[workerID].PushBack(task)
}

// Steal 窃取任务
func (ws *WorkStealing) Steal(workerID int) (interface{}, bool) {
 // 从其他队列窃取任务
 for i := 0; i < ws.workers; i++ {
  if i != workerID {
   if task, ok := ws.queues[i].PopFront(); ok {
    return task, true
   }
  }
 }
 return nil, false
}
```

## 7. 总结

### 7.1 并发编程最佳实践

1. **使用channel进行通信**：避免共享内存
2. **合理使用goroutine**：控制goroutine数量
3. **正确处理同步**：使用适当的同步原语
4. **避免竞态条件**：使用原子操作或锁
5. **资源管理**：及时关闭channel和goroutine

### 7.2 性能考虑

1. **goroutine开销**：每个goroutine约2KB内存
2. **channel性能**：无缓冲channel性能更好
3. **锁竞争**：减少锁的持有时间
4. **内存分配**：减少不必要的内存分配

### 7.3 调试技巧

1. **使用race detector**：`go run -race`
2. **使用pprof**：分析goroutine和内存使用
3. **使用trace**：分析程序执行轨迹
4. **日志记录**：记录goroutine生命周期

---

**参考文献**：

1. Pike, R. (2012). Concurrency is not parallelism
2. Cox, R. (2014). Go concurrency patterns
3. Donovan, A. A. A., & Kernighan, B. W. (2015). The Go Programming Language
