# 03-任务队列模式 (Task Queue Pattern)

## 目录

- [03-任务队列模式 (Task Queue Pattern)](#03-任务队列模式-task-queue-pattern)
  - [目录](#目录)
  - [1. 概念与定义](#1-概念与定义)
  - [2. 形式化定义](#2-形式化定义)
  - [3. 数学证明](#3-数学证明)
  - [4. 设计原则](#4-设计原则)
  - [5. Go语言实现](#5-go语言实现)
  - [6. 应用场景](#6-应用场景)
  - [7. 性能分析](#7-性能分析)
  - [8. 最佳实践](#8-最佳实践)
  - [9. 相关模式](#9-相关模式)

## 1. 概念与定义

### 1.1 基本概念

任务队列模式是一种用于管理异步任务执行的设计模式。它将任务提交和任务执行分离，通过队列机制实现任务的缓冲、调度和并发处理。

**定义**: 任务队列模式提供了一个异步任务处理框架，通过队列机制管理任务的提交、调度和执行，支持任务的优先级、重试、超时等特性。

### 1.2 核心组件

- **TaskQueue (任务队列)**: 核心队列组件，负责任务的存储和调度
- **Task (任务)**: 表示一个具体的执行任务
- **Worker (工作者)**: 负责从队列中获取任务并执行
- **TaskExecutor (任务执行器)**: 执行具体的任务逻辑
- **TaskScheduler (任务调度器)**: 负责任务的调度和优先级管理

### 1.3 模式结构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   TaskQueue     │    │      Task       │    │     Worker      │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + submit()      │◄──►│ + execute()     │◄──►│ + process()     │
│ + schedule()    │    │ + priority      │    │ + run()         │
│ + cancel()      │    │ + retryCount    │    │ + stop()        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ TaskScheduler   │    │ TaskExecutor    │    │ WorkerPool      │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ + prioritize()  │    │ + execute()     │    │ + workers       │
│ + balance()     │    │ + validate()    │    │ + distribute()  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 2. 形式化定义

### 2.1 任务队列数学模型

设 $Q = (T, W, S, P)$ 为一个任务队列系统，其中：

- $T = \{t_1, t_2, ..., t_n\}$ 是任务集合
- $W = \{w_1, w_2, ..., w_m\}$ 是工作者集合
- $S: T \rightarrow \mathbb{N}$ 是调度函数，$\mathbb{N}$ 是自然数集合
- $P: T \rightarrow \mathbb{R}$ 是优先级函数，$\mathbb{R}$ 是实数集合

### 2.2 任务执行函数

对于任务 $t \in T$，执行函数定义为：

$$execute(t, worker) = (result, status)$$

其中：
- $result$ 是执行结果
- $status \in \{success, failed, timeout, cancelled\}$ 是执行状态

### 2.3 队列操作函数

队列操作函数定义为：

$$enqueue(Q, t) = Q'$$
$$dequeue(Q) = (t, Q')$$

其中 $Q'$ 是操作后的队列状态。

## 3. 数学证明

### 3.1 队列公平性定理

**定理**: 如果任务队列使用FIFO（先进先出）策略，则任务按照提交顺序执行。

**证明**:
1. 设任务 $t_1, t_2, ..., t_n$ 按顺序提交到队列
2. FIFO策略确保 $dequeue(Q)$ 总是返回最早提交的任务
3. 因此，任务执行顺序与提交顺序一致

### 3.2 优先级调度定理

**定理**: 如果任务队列使用优先级调度，则高优先级任务优先执行。

**证明**:
1. 设任务 $t_i$ 和 $t_j$ 的优先级分别为 $P(t_i)$ 和 $P(t_j)$
2. 如果 $P(t_i) > P(t_j)$，则 $t_i$ 优先于 $t_j$ 执行
3. 调度函数 $S$ 确保高优先级任务优先被选择

### 3.3 负载均衡定理

**定理**: 如果工作者数量为 $m$，任务数量为 $n$，则平均每个工作者处理 $\frac{n}{m}$ 个任务。

**证明**:
1. 设工作者集合 $W = \{w_1, w_2, ..., w_m\}$
2. 任务集合 $T = \{t_1, t_2, ..., t_n\}$
3. 理想情况下，每个工作者处理 $\frac{n}{m}$ 个任务
4. 负载均衡算法确保任务均匀分布

## 4. 设计原则

### 4.1 单一职责原则

每个组件只负责特定的功能，队列负责存储，工作者负责执行，调度器负责调度。

### 4.2 开闭原则

可以通过扩展任务类型和工作者类型来支持新的业务逻辑，而不需要修改核心代码。

### 4.3 依赖倒置原则

任务队列依赖于抽象的任务接口，而不是具体的任务实现。

## 5. Go语言实现

### 5.1 基础实现

```go
package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task 任务接口
type Task interface {
	Execute(ctx context.Context) (interface{}, error)
	GetID() string
	GetPriority() int
	GetRetryCount() int
	IncrementRetryCount()
	GetMaxRetries() int
	GetTimeout() time.Duration
}

// TaskResult 任务结果
type TaskResult struct {
	TaskID   string
	Result   interface{}
	Error    error
	Duration time.Duration
}

// TaskQueue 任务队列
type TaskQueue struct {
	tasks    []Task
	mu       sync.RWMutex
	notEmpty chan struct{}
}

// NewTaskQueue 创建任务队列
func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks:    make([]Task, 0),
		notEmpty: make(chan struct{}, 1),
	}
}

// Submit 提交任务
func (q *TaskQueue) Submit(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	q.tasks = append(q.tasks, task)
	
	select {
	case q.notEmpty <- struct{}{}:
	default:
	}
}

// Dequeue 获取任务
func (q *TaskQueue) Dequeue() Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if len(q.tasks) == 0 {
		return nil
	}
	
	// 按优先级排序，高优先级在前
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	
	return task
}

// WaitForTask 等待任务
func (q *TaskQueue) WaitForTask(ctx context.Context) (Task, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		task := q.Dequeue()
		if task != nil {
			return task, nil
		}
		
		select {
		case <-q.notEmpty:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Size 获取队列大小
func (q *TaskQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}

// Worker 工作者
type Worker struct {
	id       string
	queue    *TaskQueue
	executor TaskExecutor
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewWorker 创建工作者
func NewWorker(id string, queue *TaskQueue, executor TaskExecutor) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		id:       id,
		queue:    queue,
		executor: executor,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动工作者
func (w *Worker) Start() {
	w.wg.Add(1)
	go w.run()
}

// Stop 停止工作者
func (w *Worker) Stop() {
	w.cancel()
	w.wg.Wait()
}

// run 运行工作者
func (w *Worker) run() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}
		
		task, err := w.queue.WaitForTask(w.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			continue
		}
		
		w.processTask(task)
	}
}

// processTask 处理任务
func (w *Worker) processTask(task Task) {
	ctx, cancel := context.WithTimeout(w.ctx, task.GetTimeout())
	defer cancel()
	
	start := time.Now()
	result, err := w.executor.Execute(ctx, task)
	duration := time.Since(start)
	
	if err != nil {
		if task.GetRetryCount() < task.GetMaxRetries() {
			task.IncrementRetryCount()
			w.queue.Submit(task)
			fmt.Printf("工作者 %s: 任务 %s 执行失败，将重试 (重试次数: %d)\n", 
				w.id, task.GetID(), task.GetRetryCount())
		} else {
			fmt.Printf("工作者 %s: 任务 %s 执行失败，已达到最大重试次数\n", 
				w.id, task.GetID())
		}
	} else {
		fmt.Printf("工作者 %s: 任务 %s 执行成功，耗时: %v\n", 
			w.id, task.GetID(), duration)
	}
}

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Execute(ctx context.Context, task Task) (interface{}, error)
}

// DefaultTaskExecutor 默认任务执行器
type DefaultTaskExecutor struct{}

func (e *DefaultTaskExecutor) Execute(ctx context.Context, task Task) (interface{}, error) {
	return task.Execute(ctx)
}
```

### 5.2 泛型实现

```go
package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task[T] 泛型任务接口
type Task[T any] interface {
	Execute(ctx context.Context) (T, error)
	GetID() string
	GetPriority() int
	GetRetryCount() int
	IncrementRetryCount()
	GetMaxRetries() int
	GetTimeout() time.Duration
}

// TaskResult[T] 泛型任务结果
type TaskResult[T any] struct {
	TaskID   string
	Result   T
	Error    error
	Duration time.Duration
}

// TaskQueue[T] 泛型任务队列
type TaskQueue[T any] struct {
	tasks    []Task[T]
	mu       sync.RWMutex
	notEmpty chan struct{}
}

// NewTaskQueue[T] 创建泛型任务队列
func NewTaskQueue[T any]() *TaskQueue[T] {
	return &TaskQueue[T]{
		tasks:    make([]Task[T], 0),
		notEmpty: make(chan struct{}, 1),
	}
}

// Submit[T] 提交泛型任务
func (q *TaskQueue[T]) Submit(task Task[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	q.tasks = append(q.tasks, task)
	
	select {
	case q.notEmpty <- struct{}{}:
	default:
	}
}

// Dequeue[T] 获取泛型任务
func (q *TaskQueue[T]) Dequeue() Task[T] {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if len(q.tasks) == 0 {
		return nil
	}
	
	// 按优先级排序，高优先级在前
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	
	return task
}

// WaitForTask[T] 等待泛型任务
func (q *TaskQueue[T]) WaitForTask(ctx context.Context) (Task[T], error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		task := q.Dequeue()
		if task != nil {
			return task, nil
		}
		
		select {
		case <-q.notEmpty:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Size[T] 获取泛型队列大小
func (q *TaskQueue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.tasks)
}

// Worker[T] 泛型工作者
type Worker[T any] struct {
	id       string
	queue    *TaskQueue[T]
	executor TaskExecutor[T]
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewWorker[T] 创建泛型工作者
func NewWorker[T any](id string, queue *TaskQueue[T], executor TaskExecutor[T]) *Worker[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker[T]{
		id:       id,
		queue:    queue,
		executor: executor,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start[T] 启动泛型工作者
func (w *Worker[T]) Start() {
	w.wg.Add(1)
	go w.run()
}

// Stop[T] 停止泛型工作者
func (w *Worker[T]) Stop() {
	w.cancel()
	w.wg.Wait()
}

// run[T] 运行泛型工作者
func (w *Worker[T]) run() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}
		
		task, err := w.queue.WaitForTask(w.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			continue
		}
		
		w.processTask(task)
	}
}

// processTask[T] 处理泛型任务
func (w *Worker[T]) processTask(task Task[T]) {
	ctx, cancel := context.WithTimeout(w.ctx, task.GetTimeout())
	defer cancel()
	
	start := time.Now()
	result, err := w.executor.Execute(ctx, task)
	duration := time.Since(start)
	
	if err != nil {
		if task.GetRetryCount() < task.GetMaxRetries() {
			task.IncrementRetryCount()
			w.queue.Submit(task)
			fmt.Printf("工作者 %s: 任务 %s 执行失败，将重试 (重试次数: %d)\n", 
				w.id, task.GetID(), task.GetRetryCount())
		} else {
			fmt.Printf("工作者 %s: 任务 %s 执行失败，已达到最大重试次数\n", 
				w.id, task.GetID())
		}
	} else {
		fmt.Printf("工作者 %s: 任务 %s 执行成功，耗时: %v\n", 
			w.id, task.GetID(), duration)
	}
}

// TaskExecutor[T] 泛型任务执行器接口
type TaskExecutor[T any] interface {
	Execute(ctx context.Context, task Task[T]) (T, error)
}

// DefaultTaskExecutor[T] 默认泛型任务执行器
type DefaultTaskExecutor[T any] struct{}

func (e *DefaultTaskExecutor[T]) Execute(ctx context.Context, task Task[T]) (T, error) {
	return task.Execute(ctx)
}
```

### 5.3 并发安全实现

```go
package taskqueue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ConcurrentTaskQueue 并发任务队列
type ConcurrentTaskQueue struct {
	tasks       []Task
	mu          sync.RWMutex
	notEmpty    chan struct{}
	workers     []*Worker
	workerPool  chan *Worker
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewConcurrentTaskQueue 创建并发任务队列
func NewConcurrentTaskQueue(workerCount int) *ConcurrentTaskQueue {
	ctx, cancel := context.WithCancel(context.Background())
	
	cq := &ConcurrentTaskQueue{
		tasks:      make([]Task, 0),
		notEmpty:   make(chan struct{}, 1),
		workers:    make([]*Worker, workerCount),
		workerPool: make(chan *Worker, workerCount),
		ctx:        ctx,
		cancel:     cancel,
	}
	
	// 创建工作协程
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(fmt.Sprintf("worker-%d", i), cq, &DefaultTaskExecutor{})
		cq.workers[i] = worker
		cq.workerPool <- worker
	}
	
	return cq
}

// Start 启动并发任务队列
func (cq *ConcurrentTaskQueue) Start() {
	for _, worker := range cq.workers {
		worker.Start()
	}
}

// Stop 停止并发任务队列
func (cq *ConcurrentTaskQueue) Stop() {
	cq.cancel()
	
	for _, worker := range cq.workers {
		worker.Stop()
	}
	
	cq.wg.Wait()
}

// Submit 提交任务
func (cq *ConcurrentTaskQueue) Submit(task Task) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	
	cq.tasks = append(cq.tasks, task)
	
	select {
	case cq.notEmpty <- struct{}{}:
	default:
	}
}

// Dequeue 获取任务
func (cq *ConcurrentTaskQueue) Dequeue() Task {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	
	if len(cq.tasks) == 0 {
		return nil
	}
	
	// 按优先级排序，高优先级在前
	task := cq.tasks[0]
	cq.tasks = cq.tasks[1:]
	
	return task
}

// WaitForTask 等待任务
func (cq *ConcurrentTaskQueue) WaitForTask(ctx context.Context) (Task, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		task := cq.Dequeue()
		if task != nil {
			return task, nil
		}
		
		select {
		case <-cq.notEmpty:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Size 获取队列大小
func (cq *ConcurrentTaskQueue) Size() int {
	cq.mu.RLock()
	defer cq.mu.RUnlock()
	return len(cq.tasks)
}

// GetWorker 获取工作者
func (cq *ConcurrentTaskQueue) GetWorker() *Worker {
	select {
	case worker := <-cq.workerPool:
		return worker
	case <-time.After(5 * time.Second):
		return nil
	}
}

// ReturnWorker 返回工作者
func (cq *ConcurrentTaskQueue) ReturnWorker(worker *Worker) {
	select {
	case cq.workerPool <- worker:
	default:
	}
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	queues map[string]*ConcurrentTaskQueue
	mu     sync.RWMutex
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		queues: make(map[string]*ConcurrentTaskQueue),
	}
}

// CreateQueue 创建队列
func (s *TaskScheduler) CreateQueue(name string, workerCount int) *ConcurrentTaskQueue {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	queue := NewConcurrentTaskQueue(workerCount)
	s.queues[name] = queue
	queue.Start()
	
	return queue
}

// GetQueue 获取队列
func (s *TaskScheduler) GetQueue(name string) *ConcurrentTaskQueue {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.queues[name]
}

// SubmitToQueue 提交任务到指定队列
func (s *TaskScheduler) SubmitToQueue(queueName string, task Task) error {
	queue := s.GetQueue(queueName)
	if queue == nil {
		return fmt.Errorf("队列 %s 不存在", queueName)
	}
	
	queue.Submit(task)
	return nil
}

// StopAll 停止所有队列
func (s *TaskScheduler) StopAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, queue := range s.queues {
		queue.Stop()
	}
}
```

## 6. 应用场景

### 6.1 图像处理任务队列

```go
package imageprocessing

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// ImageTask 图像处理任务
type ImageTask struct {
	ID           string
	InputPath    string
	OutputPath   string
	Operation    string
	Priority     int
	RetryCount   int
	MaxRetries   int
	Timeout      time.Duration
}

func (t *ImageTask) Execute(ctx context.Context) (interface{}, error) {
	fmt.Printf("处理图像: %s -> %s (操作: %s)\n", t.InputPath, t.OutputPath, t.Operation)
	
	// 模拟图像处理
	time.Sleep(100 * time.Millisecond)
	
	// 根据操作类型处理图像
	switch t.Operation {
	case "resize":
		return t.resizeImage()
	case "convert":
		return t.convertImage()
	case "compress":
		return t.compressImage()
	default:
		return nil, fmt.Errorf("不支持的操作: %s", t.Operation)
	}
}

func (t *ImageTask) GetID() string {
	return t.ID
}

func (t *ImageTask) GetPriority() int {
	return t.Priority
}

func (t *ImageTask) GetRetryCount() int {
	return t.RetryCount
}

func (t *ImageTask) IncrementRetryCount() {
	t.RetryCount++
}

func (t *ImageTask) GetMaxRetries() int {
	return t.MaxRetries
}

func (t *ImageTask) GetTimeout() time.Duration {
	return t.Timeout
}

func (t *ImageTask) resizeImage() (interface{}, error) {
	// 模拟图像缩放
	return map[string]interface{}{
		"operation": "resize",
		"result":    "success",
	}, nil
}

func (t *ImageTask) convertImage() (interface{}, error) {
	// 模拟图像格式转换
	return map[string]interface{}{
		"operation": "convert",
		"result":    "success",
	}, nil
}

func (t *ImageTask) compressImage() (interface{}, error) {
	// 模拟图像压缩
	return map[string]interface{}{
		"operation": "compress",
		"result":    "success",
	}, nil
}

// CreateImageProcessingQueue 创建图像处理队列
func CreateImageProcessingQueue(workerCount int) *ConcurrentTaskQueue {
	queue := NewConcurrentTaskQueue(workerCount)
	queue.Start()
	return queue
}
```

### 6.2 邮件发送任务队列

```go
package email

import (
	"context"
	"fmt"
	"time"
)

// EmailTask 邮件发送任务
type EmailTask struct {
	ID         string
	To         string
	Subject    string
	Body       string
	Priority   int
	RetryCount int
	MaxRetries int
	Timeout    time.Duration
}

func (t *EmailTask) Execute(ctx context.Context) (interface{}, error) {
	fmt.Printf("发送邮件: %s -> %s (主题: %s)\n", t.To, t.Subject, t.Subject)
	
	// 模拟邮件发送
	time.Sleep(200 * time.Millisecond)
	
	// 模拟发送失败的情况
	if t.To == "invalid@example.com" {
		return nil, fmt.Errorf("邮件地址无效")
	}
	
	return map[string]interface{}{
		"to":      t.To,
		"subject": t.Subject,
		"status":  "sent",
	}, nil
}

func (t *EmailTask) GetID() string {
	return t.ID
}

func (t *EmailTask) GetPriority() int {
	return t.Priority
}

func (t *EmailTask) GetRetryCount() int {
	return t.RetryCount
}

func (t *EmailTask) IncrementRetryCount() {
	t.RetryCount++
}

func (t *EmailTask) GetMaxRetries() int {
	return t.MaxRetries
}

func (t *EmailTask) GetTimeout() time.Duration {
	return t.Timeout
}

// CreateEmailQueue 创建邮件队列
func CreateEmailQueue(workerCount int) *ConcurrentTaskQueue {
	queue := NewConcurrentTaskQueue(workerCount)
	queue.Start()
	return queue
}
```

## 7. 性能分析

### 7.1 时间复杂度

- **任务提交**: $O(1)$
- **任务获取**: $O(1)$
- **任务执行**: $O(1)$ 每个任务
- **优先级排序**: $O(n \log n)$，其中 $n$ 是任务数量

### 7.2 空间复杂度

- **任务存储**: $O(n)$，其中 $n$ 是任务数量
- **工作者存储**: $O(m)$，其中 $m$ 是工作者数量
- **结果存储**: $O(k)$，其中 $k$ 是结果数量

### 7.3 并发性能

- **并行执行**: 支持多个工作者并行处理任务
- **负载均衡**: 自动分配任务给空闲的工作者
- **资源管理**: 控制并发工作者数量，避免资源耗尽

## 8. 最佳实践

### 8.1 任务设计原则

1. **任务粒度**: 设计合适粒度的任务，避免过大或过小
2. **幂等性**: 确保任务可以安全地重试
3. **超时设置**: 为每个任务设置合理的超时时间
4. **错误处理**: 明确定义错误处理和重试策略

### 8.2 性能优化

1. **工作者数量**: 根据系统资源调整工作者数量
2. **队列大小**: 设置合适的队列大小，避免内存溢出
3. **批量处理**: 对于小任务，考虑批量处理以提高效率

### 8.3 监控和调试

1. **任务监控**: 监控任务执行时间和成功率
2. **队列监控**: 监控队列大小和工作者状态
3. **错误追踪**: 记录详细的错误信息和重试历史

## 9. 相关模式

### 9.1 生产者-消费者模式

任务队列模式可以看作是生产者-消费者模式的扩展，提供了更复杂的任务管理功能。

### 9.2 线程池模式

工作者池可以看作是线程池模式的实现，提供了可重用的执行资源。

### 9.3 观察者模式

任务队列可以使用观察者模式来通知任务状态变化和完成事件。

---

**相关链接**:
- [01-状态机模式](../01-状态机模式/README.md)
- [02-工作流引擎模式](../02-工作流引擎模式/README.md)
- [04-编排vs协同模式](../04-编排vs协同模式/README.md)
- [返回上级目录](../../README.md) 