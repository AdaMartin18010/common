# 03-任务队列模式 (Task Queue Pattern)

## 目录

- [03-任务队列模式 (Task Queue Pattern)](#03-任务队列模式-task-queue-pattern)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 定义](#11-定义)
    - [1.2 问题描述](#12-问题描述)
    - [1.3 设计目标](#13-设计目标)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 任务队列模型](#21-任务队列模型)
    - [2.2 调度语义](#22-调度语义)
    - [2.3 队列正确性](#23-队列正确性)
  - [3. 数学基础](#3-数学基础)
    - [3.1 队列理论](#31-队列理论)
    - [3.2 优先级调度](#32-优先级调度)
    - [3.3 负载均衡](#33-负载均衡)
  - [4. 队列模型](#4-队列模型)
    - [4.1 队列类型](#41-队列类型)
    - [4.2 任务状态](#42-任务状态)
    - [4.3 工作者模型](#43-工作者模型)
  - [5. 调度策略](#5-调度策略)
    - [5.1 先进先出 (FIFO)](#51-先进先出-fifo)
    - [5.2 优先级调度](#52-优先级调度)
    - [5.3 轮询调度](#53-轮询调度)
    - [5.4 最少连接调度](#54-最少连接调度)
  - [6. Go语言实现](#6-go语言实现)
    - [6.1 基础结构定义](#61-基础结构定义)
    - [6.2 FIFO队列实现](#62-fifo队列实现)
    - [6.3 优先级队列实现](#63-优先级队列实现)
    - [6.4 延迟队列实现](#64-延迟队列实现)
    - [6.5 任务实现](#65-任务实现)
    - [6.6 工作者实现](#66-工作者实现)
    - [6.7 任务队列管理器实现](#67-任务队列管理器实现)
    - [6.8 使用示例](#68-使用示例)
  - [7. 性能分析](#7-性能分析)
    - [7.1 时间复杂度](#71-时间复杂度)
    - [7.2 空间复杂度](#72-空间复杂度)
    - [7.3 吞吐量分析](#73-吞吐量分析)
  - [8. 应用场景](#8-应用场景)
    - [8.1 异步处理](#81-异步处理)
    - [8.2 负载均衡](#82-负载均衡)
    - [8.3 批处理](#83-批处理)
    - [8.4 实时处理](#84-实时处理)
  - [9. 最佳实践](#9-最佳实践)
    - [9.1 队列设计](#91-队列设计)
    - [9.2 工作者管理](#92-工作者管理)
    - [9.3 错误处理](#93-错误处理)
    - [9.4 性能优化](#94-性能优化)
  - [10. 总结](#10-总结)
    - [10.1 关键要点](#101-关键要点)
    - [10.2 未来发展方向](#102-未来发展方向)

## 1. 概述

### 1.1 定义

任务队列模式是一种用于管理异步任务执行的设计模式。它将任务提交到队列中，由后台工作者从队列中获取任务并执行，实现任务的异步处理和负载均衡。

### 1.2 问题描述

在需要处理大量异步任务的系统中，直接同步处理面临以下挑战：

- **性能瓶颈**: 同步处理导致响应延迟
- **资源竞争**: 并发任务竞争系统资源
- **故障传播**: 单个任务失败影响整个系统
- **扩展困难**: 难以根据负载动态调整处理能力

### 1.3 设计目标

1. **异步处理**: 任务提交后立即返回，后台异步执行
2. **负载均衡**: 多个工作者并行处理任务
3. **故障隔离**: 单个任务失败不影响其他任务
4. **动态扩展**: 根据负载动态调整工作者数量

## 2. 形式化定义

### 2.1 任务队列模型

**定义 2.1 (任务队列)**
任务队列是一个三元组 ```latex
Q = (T, W, S)
```，其中：

- ```latex
T = \{t_1, t_2, ..., t_n\}
``` 是任务集合
- ```latex
W = \{w_1, w_2, ..., w_m\}
``` 是工作者集合
- ```latex
S: T \rightarrow W
``` 是调度函数

**定义 2.2 (任务)**
任务是一个四元组 ```latex
t = (id, payload, priority, deadline)
```，其中：

- ```latex
id
``` 是任务唯一标识符
- ```latex
payload
``` 是任务执行内容
- ```latex
priority
``` 是任务优先级
- ```latex
deadline
``` 是任务截止时间

### 2.2 调度语义

**定义 2.3 (调度策略)**
调度策略是一个函数 ```latex
f: T \times W \rightarrow \mathbb{R}
```，用于计算任务与工作者的匹配度。

**定义 2.4 (最优调度)**
最优调度是使总执行时间最小的调度：
$```latex
\text{minimize } \sum_{t \in T} \text{execution\_time}(t, S(t))
```$

### 2.3 队列正确性

**定理 2.1 (队列正确性)**
任务队列是正确的，当且仅当：

1. **完整性**: 所有任务最终被执行
2. **顺序性**: 优先级高的任务优先执行
3. **公平性**: 相同优先级的任务按FIFO顺序执行

**证明**:

- **完整性**: 通过工作者持续处理保证
- **顺序性**: 通过优先级队列保证
- **公平性**: 通过FIFO队列保证

## 3. 数学基础

### 3.1 队列理论

**定义 3.1 (M/M/c队列)**
M/M/c队列是具有泊松到达、指数服务时间和c个服务器的队列模型。

-**定理 3.1 (队列长度分布)**

```latex
M/M/c队列中任务数的稳态分布为：
$P(N = n) = \begin{cases}
\frac{\rho^n}{n!}P_0 & \text{if } n < c \\
\frac{\rho^n}{c!c^{n-c}}P_0 & \text{if } n \geq c
\end{cases}$
其中 ```latex
\rho = \frac{\lambda}{c\mu}
``` 是系统利用率。
```

**定理 3.2 (等待时间)**
M/M/c队列的平均等待时间为：
$```latex
W_q = \frac{P_c}{\mu c(1-\rho)}
```$
其中 ```latex
P_c
``` 是所有服务器都忙的概率。

### 3.2 优先级调度

**定义 3.2 (优先级队列)**
优先级队列根据任务优先级排序：
$```latex
\forall t_1, t_2 \in T: \text{priority}(t_1) > \text{priority}(t_2) \Rightarrow t_1 \prec t_2
```$

**定理 3.3 (优先级调度性能)**
优先级调度的平均等待时间为：
$```latex
W_q = \sum_{i=1}^{k} \frac{\lambda_i W_{q,i}}{\lambda}
```$
其中 ```latex
k
``` 是优先级级别数，```latex
\lambda_i
``` 是第 ```latex
i
``` 级任务的到达率。

### 3.3 负载均衡

**定义 3.3 (负载均衡度)**
负载均衡度定义为：
$```latex
\text{Balance} = 1 - \frac{\max_{w \in W} L(w) - \min_{w \in W} L(w)}{\max_{w \in W} L(w)}
```$
其中 ```latex
L(w)
``` 是工作者 ```latex
w
``` 的负载。

## 4. 队列模型

### 4.1 队列类型

```go
// QueueType 队列类型
type QueueType string

const (
    FIFOQueue      QueueType = "fifo"
    PriorityQueue  QueueType = "priority"
    LIFOQueue      QueueType = "lifo"
    DelayQueue     QueueType = "delay"
)

// Task 任务接口
type Task interface {
    ID() string
    Priority() int
    Deadline() time.Time
    Execute() error
    RetryCount() int
    MaxRetries() int
    CanRetry() bool
}
```

### 4.2 任务状态

**定义 4.1 (任务状态)**
任务状态是一个有限状态机：

- **Pending**: 任务已提交，等待执行
- **Running**: 任务正在执行
- **Completed**: 任务执行成功
- **Failed**: 任务执行失败
- **Cancelled**: 任务被取消

### 4.3 工作者模型

```go
// Worker 工作者接口
type Worker interface {
    ID() string
    Status() WorkerStatus
    ProcessTask(task Task) error
    Stop() error
}

// WorkerStatus 工作者状态
type WorkerStatus string

const (
    IdleStatus     WorkerStatus = "idle"
    BusyStatus     WorkerStatus = "busy"
    StoppedStatus  WorkerStatus = "stopped"
)
```

## 5. 调度策略

### 5.1 先进先出 (FIFO)

**算法描述**:

1. 任务按提交顺序排队
2. 工作者按顺序获取任务
3. 简单公平，适合所有任务优先级相同的情况

**时间复杂度**: ```latex
O(1)
```
**空间复杂度**: ```latex
O(n)
```

### 5.2 优先级调度

**算法描述**:

1. 任务按优先级排序
2. 高优先级任务优先执行
3. 相同优先级按FIFO顺序

**时间复杂度**: ```latex
O(\log n)
```
**空间复杂度**: ```latex
O(n)
```

### 5.3 轮询调度

**算法描述**:

1. 工作者轮流获取任务
2. 确保负载均衡
3. 适合任务执行时间相近的情况

**时间复杂度**: ```latex
O(1)
```
**空间复杂度**: ```latex
O(n)
```

### 5.4 最少连接调度

**算法描述**:

1. 选择当前负载最少的工作者
2. 动态负载均衡
3. 适合任务执行时间差异较大的情况

**时间复杂度**: ```latex
O(m)
```
**空间复杂度**: ```latex
O(n)
```

## 6. Go语言实现

### 6.1 基础结构定义

```go
// TaskQueue 任务队列接口
type TaskQueue interface {
    Enqueue(task Task) error
    Dequeue() (Task, error)
    Size() int
    IsEmpty() bool
    Clear() error
}

// TaskQueueManager 任务队列管理器
type TaskQueueManager struct {
    queue       TaskQueue
    workers     []Worker
    config      *QueueConfig
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

// QueueConfig 队列配置
type QueueConfig struct {
    QueueType     QueueType
    WorkerCount   int
    BufferSize    int
    RetryDelay    time.Duration
    MaxRetries    int
    TaskTimeout   time.Duration
}
```

### 6.2 FIFO队列实现

```go
// FIFOQueue FIFO队列实现
type FIFOQueue struct {
    tasks chan Task
    mu    sync.RWMutex
}

// NewFIFOQueue 创建FIFO队列
func NewFIFOQueue(bufferSize int) *FIFOQueue {
    return &FIFOQueue{
        tasks: make(chan Task, bufferSize),
    }
}

// Enqueue 入队
func (fq *FIFOQueue) Enqueue(task Task) error {
    select {
    case fq.tasks <- task:
        return nil
    default:
        return fmt.Errorf("queue is full")
    }
}

// Dequeue 出队
func (fq *FIFOQueue) Dequeue() (Task, error) {
    select {
    case task := <-fq.tasks:
        return task, nil
    default:
        return nil, fmt.Errorf("queue is empty")
    }
}

// Size 获取队列大小
func (fq *FIFOQueue) Size() int {
    return len(fq.tasks)
}

// IsEmpty 检查队列是否为空
func (fq *FIFOQueue) IsEmpty() bool {
    return len(fq.tasks) == 0
}

// Clear 清空队列
func (fq *FIFOQueue) Clear() error {
    fq.mu.Lock()
    defer fq.mu.Unlock()

    for len(fq.tasks) > 0 {
        <-fq.tasks
    }

    return nil
}
```

### 6.3 优先级队列实现

```go
// PriorityQueue 优先级队列实现
type PriorityQueue struct {
    tasks []Task
    mu    sync.RWMutex
}

// NewPriorityQueue 创建优先级队列
func NewPriorityQueue() *PriorityQueue {
    return &PriorityQueue{
        tasks: make([]Task, 0),
    }
}

// Enqueue 入队
func (pq *PriorityQueue) Enqueue(task Task) error {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    pq.tasks = append(pq.tasks, task)
    pq.heapifyUp(len(pq.tasks) - 1)

    return nil
}

// Dequeue 出队
func (pq *PriorityQueue) Dequeue() (Task, error) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if len(pq.tasks) == 0 {
        return nil, fmt.Errorf("queue is empty")
    }

    task := pq.tasks[0]
    pq.tasks[0] = pq.tasks[len(pq.tasks)-1]
    pq.tasks = pq.tasks[:len(pq.tasks)-1]

    if len(pq.tasks) > 0 {
        pq.heapifyDown(0)
    }

    return task, nil
}

// heapifyUp 向上堆化
func (pq *PriorityQueue) heapifyUp(index int) {
    for index > 0 {
        parent := (index - 1) / 2
        if pq.tasks[index].Priority() > pq.tasks[parent].Priority() {
            pq.tasks[index], pq.tasks[parent] = pq.tasks[parent], pq.tasks[index]
            index = parent
        } else {
            break
        }
    }
}

// heapifyDown 向下堆化
func (pq *PriorityQueue) heapifyDown(index int) {
    for {
        left := 2*index + 1
        right := 2*index + 2
        largest := index

        if left < len(pq.tasks) && pq.tasks[left].Priority() > pq.tasks[largest].Priority() {
            largest = left
        }

        if right < len(pq.tasks) && pq.tasks[right].Priority() > pq.tasks[largest].Priority() {
            largest = right
        }

        if largest == index {
            break
        }

        pq.tasks[index], pq.tasks[largest] = pq.tasks[largest], pq.tasks[index]
        index = largest
    }
}

// Size 获取队列大小
func (pq *PriorityQueue) Size() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()

    return len(pq.tasks)
}

// IsEmpty 检查队列是否为空
func (pq *PriorityQueue) IsEmpty() bool {
    pq.mu.RLock()
    defer pq.mu.RUnlock()

    return len(pq.tasks) == 0
}

// Clear 清空队列
func (pq *PriorityQueue) Clear() error {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    pq.tasks = make([]Task, 0)
    return nil
}
```

### 6.4 延迟队列实现

```go
// DelayQueue 延迟队列实现
type DelayQueue struct {
    tasks []*DelayedTask
    mu    sync.RWMutex
}

// DelayedTask 延迟任务
type DelayedTask struct {
    Task      Task
    ExecuteAt time.Time
}

// NewDelayQueue 创建延迟队列
func NewDelayQueue() *DelayQueue {
    return &DelayQueue{
        tasks: make([]*DelayedTask, 0),
    }
}

// Enqueue 入队
func (dq *DelayQueue) Enqueue(task Task) error {
    return dq.EnqueueAt(task, time.Now())
}

// EnqueueAt 在指定时间入队
func (dq *DelayQueue) EnqueueAt(task Task, executeAt time.Time) error {
    dq.mu.Lock()
    defer dq.mu.Unlock()

    delayedTask := &DelayedTask{
        Task:      task,
        ExecuteAt: executeAt,
    }

    dq.tasks = append(dq.tasks, delayedTask)
    dq.sortTasks()

    return nil
}

// Dequeue 出队
func (dq *DelayQueue) Dequeue() (Task, error) {
    dq.mu.Lock()
    defer dq.mu.Unlock()

    if len(dq.tasks) == 0 {
        return nil, fmt.Errorf("queue is empty")
    }

    // 检查第一个任务是否可以执行
    if time.Now().Before(dq.tasks[0].ExecuteAt) {
        return nil, fmt.Errorf("no task ready for execution")
    }

    task := dq.tasks[0].Task
    dq.tasks = dq.tasks[1:]

    return task, nil
}

// sortTasks 排序任务
func (dq *DelayQueue) sortTasks() {
    sort.Slice(dq.tasks, func(i, j int) bool {
        return dq.tasks[i].ExecuteAt.Before(dq.tasks[j].ExecuteAt)
    })
}

// Size 获取队列大小
func (dq *DelayQueue) Size() int {
    dq.mu.RLock()
    defer dq.mu.RUnlock()

    return len(dq.tasks)
}

// IsEmpty 检查队列是否为空
func (dq *DelayQueue) IsEmpty() bool {
    dq.mu.RLock()
    defer dq.mu.RUnlock()

    return len(dq.tasks) == 0
}

// Clear 清空队列
func (dq *DelayQueue) Clear() error {
    dq.mu.Lock()
    defer dq.mu.Unlock()

    dq.tasks = make([]*DelayedTask, 0)
    return nil
}
```

### 6.5 任务实现

```go
// BaseTask 基础任务实现
type BaseTask struct {
    id         string
    priority   int
    deadline   time.Time
    retryCount int
    maxRetries int
    payload    interface{}
    executeFn  func() error
}

// NewBaseTask 创建基础任务
func NewBaseTask(id string, priority int, payload interface{}, executeFn func() error) *BaseTask {
    return &BaseTask{
        id:         id,
        priority:   priority,
        deadline:   time.Now().Add(24 * time.Hour), // 默认24小时
        retryCount: 0,
        maxRetries: 3,
        payload:    payload,
        executeFn:  executeFn,
    }
}

// ID 获取任务ID
func (bt *BaseTask) ID() string {
    return bt.id
}

// Priority 获取优先级
func (bt *BaseTask) Priority() int {
    return bt.priority
}

// Deadline 获取截止时间
func (bt *BaseTask) Deadline() time.Time {
    return bt.deadline
}

// Execute 执行任务
func (bt *BaseTask) Execute() error {
    if bt.executeFn != nil {
        return bt.executeFn()
    }
    return fmt.Errorf("no execute function defined")
}

// RetryCount 获取重试次数
func (bt *BaseTask) RetryCount() int {
    return bt.retryCount
}

// MaxRetries 获取最大重试次数
func (bt *BaseTask) MaxRetries() int {
    return bt.maxRetries
}

// CanRetry 检查是否可以重试
func (bt *BaseTask) CanRetry() bool {
    return bt.retryCount < bt.maxRetries
}

// IncrementRetry 增加重试次数
func (bt *BaseTask) IncrementRetry() {
    bt.retryCount++
}

// SetDeadline 设置截止时间
func (bt *BaseTask) SetDeadline(deadline time.Time) {
    bt.deadline = deadline
}

// SetMaxRetries 设置最大重试次数
func (bt *BaseTask) SetMaxRetries(maxRetries int) {
    bt.maxRetries = maxRetries
}
```

### 6.6 工作者实现

```go
// WorkerImpl 工作者实现
type WorkerImpl struct {
    id       string
    status   WorkerStatus
    queue    TaskQueue
    config   *QueueConfig
    ctx      context.Context
    cancel   context.CancelFunc
    wg       sync.WaitGroup
}

// NewWorker 创建工作者
func NewWorker(id string, queue TaskQueue, config *QueueConfig) *WorkerImpl {
    ctx, cancel := context.WithCancel(context.Background())

    return &WorkerImpl{
        id:     id,
        status: IdleStatus,
        queue:  queue,
        config: config,
        ctx:    ctx,
        cancel: cancel,
    }
}

// ID 获取工作者ID
func (w *WorkerImpl) ID() string {
    return w.id
}

// Status 获取工作者状态
func (w *WorkerImpl) Status() WorkerStatus {
    return w.status
}

// ProcessTask 处理任务
func (w *WorkerImpl) ProcessTask(task Task) error {
    w.status = BusyStatus
    defer func() {
        w.status = IdleStatus
    }()

    log.Printf("Worker %s processing task %s", w.id, task.ID())

    // 检查任务是否超时
    if time.Now().After(task.Deadline()) {
        return fmt.Errorf("task %s exceeded deadline", task.ID())
    }

    // 执行任务
    err := task.Execute()
    if err != nil {
        log.Printf("Worker %s failed to process task %s: %v", w.id, task.ID(), err)

        // 检查是否可以重试
        if task.CanRetry() {
            task.(*BaseTask).IncrementRetry()
            log.Printf("Task %s will be retried (attempt %d/%d)", task.ID(), task.RetryCount(), task.MaxRetries())

            // 延迟重试
            time.Sleep(w.config.RetryDelay)
            return w.queue.Enqueue(task)
        }

        return err
    }

    log.Printf("Worker %s completed task %s", w.id, task.ID())
    return nil
}

// Start 启动工作者
func (w *WorkerImpl) Start() {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        w.work()
    }()
}

// work 工作循环
func (w *WorkerImpl) work() {
    for {
        select {
        case <-w.ctx.Done():
            return
        default:
            // 尝试获取任务
            task, err := w.queue.Dequeue()
            if err != nil {
                // 队列为空，等待一段时间
                time.Sleep(100 * time.Millisecond)
                continue
            }

            // 处理任务
            w.ProcessTask(task)
        }
    }
}

// Stop 停止工作者
func (w *WorkerImpl) Stop() error {
    w.cancel()
    w.wg.Wait()
    w.status = StoppedStatus
    return nil
}
```

### 6.7 任务队列管理器实现

```go
// NewTaskQueueManager 创建任务队列管理器
func NewTaskQueueManager(config *QueueConfig) *TaskQueueManager {
    ctx, cancel := context.WithCancel(context.Background())

    var queue TaskQueue
    switch config.QueueType {
    case FIFOQueue:
        queue = NewFIFOQueue(config.BufferSize)
    case PriorityQueue:
        queue = NewPriorityQueue()
    case DelayQueue:
        queue = NewDelayQueue()
    default:
        queue = NewFIFOQueue(config.BufferSize)
    }

    return &TaskQueueManager{
        queue:  queue,
        workers: make([]Worker, 0),
        config: config,
        ctx:    ctx,
        cancel: cancel,
    }
}

// Start 启动任务队列管理器
func (tqm *TaskQueueManager) Start() error {
    // 创建工作者
    for i := 0; i < tqm.config.WorkerCount; i++ {
        worker := NewWorker(fmt.Sprintf("worker-%d", i), tqm.queue, tqm.config)
        tqm.workers = append(tqm.workers, worker)
        worker.Start()
    }

    log.Printf("Task queue manager started with %d workers", tqm.config.WorkerCount)
    return nil
}

// Stop 停止任务队列管理器
func (tqm *TaskQueueManager) Stop() error {
    tqm.cancel()

    // 停止所有工作者
    for _, worker := range tqm.workers {
        worker.Stop()
    }

    tqm.wg.Wait()

    log.Printf("Task queue manager stopped")
    return nil
}

// SubmitTask 提交任务
func (tqm *TaskQueueManager) SubmitTask(task Task) error {
    return tqm.queue.Enqueue(task)
}

// GetQueueSize 获取队列大小
func (tqm *TaskQueueManager) GetQueueSize() int {
    return tqm.queue.Size()
}

// GetWorkerStatus 获取工作者状态
func (tqm *TaskQueueManager) GetWorkerStatus() map[string]WorkerStatus {
    tqm.mu.RLock()
    defer tqm.mu.RUnlock()

    status := make(map[string]WorkerStatus)
    for _, worker := range tqm.workers {
        status[worker.ID()] = worker.Status()
    }

    return status
}

// AddWorker 添加工作者
func (tqm *TaskQueueManager) AddWorker() error {
    tqm.mu.Lock()
    defer tqm.mu.Unlock()

    worker := NewWorker(fmt.Sprintf("worker-%d", len(tqm.workers)), tqm.queue, tqm.config)
    tqm.workers = append(tqm.workers, worker)
    worker.Start()

    log.Printf("Added worker %s", worker.ID())
    return nil
}

// RemoveWorker 移除工作者
func (tqm *TaskQueueManager) RemoveWorker() error {
    tqm.mu.Lock()
    defer tqm.mu.Unlock()

    if len(tqm.workers) == 0 {
        return fmt.Errorf("no workers to remove")
    }

    worker := tqm.workers[len(tqm.workers)-1]
    tqm.workers = tqm.workers[:len(tqm.workers)-1]

    worker.Stop()

    log.Printf("Removed worker %s", worker.ID())
    return nil
}
```

### 6.8 使用示例

```go
// main.go
func main() {
    // 创建队列配置
    config := &QueueConfig{
        QueueType:   PriorityQueue,
        WorkerCount: 3,
        BufferSize:  1000,
        RetryDelay:  time.Second,
        MaxRetries:  3,
        TaskTimeout: 30 * time.Second,
    }

    // 创建任务队列管理器
    manager := NewTaskQueueManager(config)

    // 启动管理器
    err := manager.Start()
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Stop()

    // 创建任务
    tasks := []Task{
        NewBaseTask("task-1", 1, "Low priority task", func() error {
            time.Sleep(100 * time.Millisecond)
            log.Printf("Executed low priority task")
            return nil
        }),
        NewBaseTask("task-2", 5, "High priority task", func() error {
            time.Sleep(50 * time.Millisecond)
            log.Printf("Executed high priority task")
            return nil
        }),
        NewBaseTask("task-3", 3, "Medium priority task", func() error {
            time.Sleep(75 * time.Millisecond)
            log.Printf("Executed medium priority task")
            return nil
        }),
        NewBaseTask("task-4", 2, "Failing task", func() error {
            return fmt.Errorf("task failed")
        }),
    }

    // 设置任务属性
    tasks[3].(*BaseTask).SetMaxRetries(2)

    // 提交任务
    for _, task := range tasks {
        err := manager.SubmitTask(task)
        if err != nil {
            log.Printf("Failed to submit task %s: %v", task.ID(), err)
        }
    }

    // 监控队列状态
    go func() {
        for {
            queueSize := manager.GetQueueSize()
            workerStatus := manager.GetWorkerStatus()

            log.Printf("Queue size: %d", queueSize)
            log.Printf("Worker status: %v", workerStatus)

            if queueSize == 0 {
                break
            }

            time.Sleep(500 * time.Millisecond)
        }
    }()

    // 等待任务完成
    time.Sleep(5 * time.Second)

    // 动态调整工作者数量
    log.Printf("Adding worker...")
    manager.AddWorker()

    time.Sleep(2 * time.Second)

    log.Printf("Removing worker...")
    manager.RemoveWorker()

    time.Sleep(2 * time.Second)
}
```

## 7. 性能分析

### 7.1 时间复杂度

**定理 7.1 (任务提交时间复杂度)**
任务提交的时间复杂度为 ```latex
O(1)
```。

**定理 7.2 (任务调度时间复杂度)**
任务调度的时间复杂度为 ```latex
O(\log n)
```，其中 ```latex
n
``` 是队列中的任务数量。

**定理 7.3 (工作者调度时间复杂度)**
工作者调度的时间复杂度为 ```latex
O(m)
```，其中 ```latex
m
``` 是工作者数量。

### 7.2 空间复杂度

**定理 7.4 (任务队列空间复杂度)**
任务队列的空间复杂度为 ```latex
O(n)
```，其中 ```latex
n
``` 是任务数量。

### 7.3 吞吐量分析

**定理 7.5 (系统吞吐量)**
系统吞吐量为：
$```latex
\text{Throughput} = \min(\lambda, c\mu)
```$
其中 ```latex
\lambda
``` 是任务到达率，```latex
c
``` 是工作者数量，```latex
\mu
``` 是服务率。

## 8. 应用场景

### 8.1 异步处理

- **邮件发送**: 批量邮件处理
- **文件处理**: 大文件上传下载
- **数据处理**: 批量数据处理
- **通知推送**: 消息推送服务

### 8.2 负载均衡

- **Web服务器**: 请求分发
- **API网关**: 请求路由
- **微服务**: 服务调用
- **数据库**: 查询分发

### 8.3 批处理

- **ETL作业**: 数据提取转换加载
- **报表生成**: 定期报表生成
- **数据同步**: 跨系统数据同步
- **备份作业**: 定期备份任务

### 8.4 实时处理

- **事件处理**: 实时事件处理
- **流处理**: 数据流处理
- **监控告警**: 实时监控
- **日志处理**: 日志分析

## 9. 最佳实践

### 9.1 队列设计

```go
// 队列设计原则
type QueueDesign struct {
    // 1. 选择合适的队列类型
    // 2. 合理设置缓冲区大小
    // 3. 考虑任务优先级
    // 4. 实现任务超时机制
}
```

### 9.2 工作者管理

```go
// 工作者管理策略
type WorkerManagement struct {
    // 1. 动态调整工作者数量
    // 2. 监控工作者健康状态
    // 3. 实现工作者故障恢复
    // 4. 负载均衡策略
}
```

### 9.3 错误处理

```go
// 错误处理策略
type ErrorHandling struct {
    // 1. 任务重试机制
    // 2. 死信队列处理
    // 3. 错误监控告警
    // 4. 降级处理策略
}
```

### 9.4 性能优化

```go
// 性能优化建议
const (
    DefaultWorkerCount = 4
    DefaultBufferSize  = 1000
    DefaultRetryDelay  = time.Second
    DefaultMaxRetries  = 3
    DefaultTaskTimeout = 30 * time.Second
)
```

## 10. 总结

任务队列模式是处理异步任务和负载均衡的重要工具，通过合理设计队列和工作者，可以构建高性能、可扩展的异步处理系统。

### 10.1 关键要点

1. **异步处理**: 任务提交后立即返回
2. **负载均衡**: 多个工作者并行处理
3. **错误处理**: 健壮的重试和恢复机制
4. **动态扩展**: 根据负载调整处理能力

### 10.2 未来发展方向

1. **智能调度**: 使用ML优化任务调度
2. **分布式队列**: 支持跨节点任务分发
3. **实时监控**: 实时队列状态监控
4. **自适应调整**: 自动调整队列参数

---

**参考文献**:

1. Little, J. D. C. (1961). "A proof for the queuing formula: L = λW"
2. Kleinrock, L. (1975). "Queueing Systems, Volume 1: Theory"
3. Harchol-Balter, M. (2013). "Performance Modeling and Design of Computer Systems"

**相关链接**:

- [01-状态机模式](../01-State-Machine-Pattern.md)
- [02-工作流引擎模式](../02-Workflow-Engine-Pattern.md)
- [04-编排vs协同模式](../04-Orchestration-vs-Choreography-Pattern.md)
