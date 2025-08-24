# 实时系统 (Real-Time Systems)

## 1. 基本概念

### 1.1 实时系统定义

**实时系统 (Real-Time System)** 是一种必须在严格的时间约束内响应的计算机系统。系统的正确性不仅取决于计算结果的正确性，还取决于结果产生的时间。

### 1.2 实时系统分类

- **硬实时系统**: 必须在严格的时间限制内完成，否则会导致系统失败
- **软实时系统**: 可以容忍偶尔的时间违反，但会影响系统性能
- **固实时系统**: 介于硬实时和软实时之间，偶尔的时间违反是可接受的

### 1.3 时间约束

- **截止时间 (Deadline)**: 任务必须完成的最晚时间
- **响应时间 (Response Time)**: 从事件发生到系统响应的时间
- **抖动 (Jitter)**: 响应时间的变化
- **延迟 (Latency)**: 从输入到输出的时间间隔

## 2. 实时调度算法

### 2.1 速率单调调度 (Rate Monotonic)

```go
// 速率单调调度器
type RateMonotonicScheduler struct {
    tasks    []*Task
    current  *Task
    mu       sync.Mutex
}

type Task struct {
    ID           string
    Priority     int
    Period       time.Duration
    ExecutionTime time.Duration
    Deadline     time.Duration
    LastRelease  time.Time
    NextRelease  time.Time
    RemainingTime time.Duration
    State        TaskState
}

type TaskState int

const (
    Ready TaskState = iota
    Running
    Blocked
    Completed
)

func (rms *RateMonotonicScheduler) AddTask(task *Task) {
    rms.mu.Lock()
    defer rms.mu.Unlock()
    
    // 计算优先级（周期越短，优先级越高）
    task.Priority = int(1.0 / float64(task.Period.Milliseconds()) * 1000)
    
    rms.tasks = append(rms.tasks, task)
    
    // 按优先级排序
    sort.Slice(rms.tasks, func(i, j int) bool {
        return rms.tasks[i].Priority > rms.tasks[j].Priority
    })
}

func (rms *RateMonotonicScheduler) Schedule() *Task {
    rms.mu.Lock()
    defer rms.mu.Unlock()
    
    now := time.Now()
    
    // 检查是否有任务需要释放
    for _, task := range rms.tasks {
        if now.After(task.NextRelease) {
            task.State = Ready
            task.LastRelease = now
            task.NextRelease = now.Add(task.Period)
            task.RemainingTime = task.ExecutionTime
        }
    }
    
    // 选择最高优先级的就绪任务
    var highestPriorityTask *Task
    for _, task := range rms.tasks {
        if task.State == Ready && (highestPriorityTask == nil || task.Priority > highestPriorityTask.Priority) {
            highestPriorityTask = task
        }
    }
    
    if highestPriorityTask != nil {
        highestPriorityTask.State = Running
        rms.current = highestPriorityTask
    }
    
    return rms.current
}

func (rms *RateMonotonicScheduler) ExecuteTask(duration time.Duration) {
    rms.mu.Lock()
    defer rms.mu.Unlock()
    
    if rms.current == nil {
        return
    }
    
    // 执行任务
    if duration >= rms.current.RemainingTime {
        rms.current.RemainingTime = 0
        rms.current.State = Completed
        rms.current = nil
    } else {
        rms.current.RemainingTime -= duration
    }
}

func (rms *RateMonotonicScheduler) CheckDeadlines() []*Task {
    rms.mu.Lock()
    defer rms.mu.Unlock()
    
    now := time.Now()
    missedDeadlines := make([]*Task, 0)
    
    for _, task := range rms.tasks {
        if task.State == Ready || task.State == Running {
            deadline := task.LastRelease.Add(task.Deadline)
            if now.After(deadline) {
                missedDeadlines = append(missedDeadlines, task)
            }
        }
    }
    
    return missedDeadlines
}
```

### 2.2 最早截止时间优先 (Earliest Deadline First)

```go
// EDF调度器
type EDFScheduler struct {
    tasks    []*Task
    current  *Task
    mu       sync.Mutex
}

func (edf *EDFScheduler) AddTask(task *Task) {
    edf.mu.Lock()
    defer edf.mu.Unlock()
    
    edf.tasks = append(edf.tasks, task)
}

func (edf *EDFScheduler) Schedule() *Task {
    edf.mu.Lock()
    defer edf.mu.Unlock()
    
    now := time.Now()
    
    // 检查是否有任务需要释放
    for _, task := range edf.tasks {
        if now.After(task.NextRelease) {
            task.State = Ready
            task.LastRelease = now
            task.NextRelease = now.Add(task.Period)
            task.RemainingTime = task.ExecutionTime
        }
    }
    
    // 选择最早截止时间的就绪任务
    var earliestDeadlineTask *Task
    for _, task := range edf.tasks {
        if task.State == Ready {
            if earliestDeadlineTask == nil {
                earliestDeadlineTask = task
            } else {
                taskDeadline := task.LastRelease.Add(task.Deadline)
                currentDeadline := earliestDeadlineTask.LastRelease.Add(earliestDeadlineTask.Deadline)
                if taskDeadline.Before(currentDeadline) {
                    earliestDeadlineTask = task
                }
            }
        }
    }
    
    if earliestDeadlineTask != nil {
        earliestDeadlineTask.State = Running
        edf.current = earliestDeadlineTask
    }
    
    return edf.current
}

func (edf *EDFScheduler) ExecuteTask(duration time.Duration) {
    edf.mu.Lock()
    defer edf.mu.Unlock()
    
    if edf.current == nil {
        return
    }
    
    // 执行任务
    if duration >= edf.current.RemainingTime {
        edf.current.RemainingTime = 0
        edf.current.State = Completed
        edf.current = nil
    } else {
        edf.current.RemainingTime -= duration
    }
}
```

### 2.3 优先级继承协议

```go
// 优先级继承协议
type PriorityInheritanceProtocol struct {
    tasks    map[string]*Task
    resources map[string]*Resource
    mu       sync.Mutex
}

type Resource struct {
    ID       string
    Owner    *Task
    WaitList []*Task
    mu       sync.Mutex
}

func (pip *PriorityInheritanceProtocol) RequestResource(taskID, resourceID string) bool {
    pip.mu.Lock()
    defer pip.mu.Unlock()
    
    task := pip.tasks[taskID]
    resource := pip.resources[resourceID]
    
    if resource.Owner == nil {
        // 资源可用，直接分配
        resource.Owner = task
        return true
    }
    
    // 资源被占用，检查优先级
    if task.Priority > resource.Owner.Priority {
        // 优先级继承
        resource.Owner.Priority = task.Priority
    }
    
    // 添加到等待列表
    resource.WaitList = append(resource.WaitList, task)
    task.State = Blocked
    
    return false
}

func (pip *PriorityInheritanceProtocol) ReleaseResource(taskID, resourceID string) {
    pip.mu.Lock()
    defer pip.mu.Unlock()
    
    task := pip.tasks[taskID]
    resource := pip.resources[resourceID]
    
    if resource.Owner != task {
        return // 不是资源的所有者
    }
    
    // 恢复原始优先级
    task.Priority = task.OriginalPriority
    
    // 选择下一个等待的任务
    if len(resource.WaitList) > 0 {
        var highestPriorityTask *Task
        var highestPriorityIndex int
        
        for i, waitingTask := range resource.WaitList {
            if highestPriorityTask == nil || waitingTask.Priority > highestPriorityTask.Priority {
                highestPriorityTask = waitingTask
                highestPriorityIndex = i
            }
        }
        
        if highestPriorityTask != nil {
            resource.Owner = highestPriorityTask
            highestPriorityTask.State = Ready
            
            // 从等待列表中移除
            resource.WaitList = append(resource.WaitList[:highestPriorityIndex], resource.WaitList[highestPriorityIndex+1:]...)
        }
    } else {
        resource.Owner = nil
    }
}
```

## 3. 实时任务管理

### 3.1 周期性任务

```go
// 周期性任务管理器
type PeriodicTaskManager struct {
    tasks    map[string]*PeriodicTask
    scheduler Scheduler
    mu       sync.Mutex
}

type PeriodicTask struct {
    ID           string
    Period       time.Duration
    ExecutionTime time.Duration
    Deadline     time.Duration
    Priority     int
    Handler      TaskHandler
    State        TaskState
    LastRelease  time.Time
    NextRelease  time.Time
    MissedDeadlines int
}

type TaskHandler func() error

func (ptm *PeriodicTaskManager) AddTask(task *PeriodicTask) {
    ptm.mu.Lock()
    defer ptm.mu.Unlock()
    
    task.State = Ready
    task.LastRelease = time.Now()
    task.NextRelease = time.Now().Add(task.Period)
    
    ptm.tasks[task.ID] = task
    ptm.scheduler.AddTask(task)
}

func (ptm *PeriodicTaskManager) Start() {
    ticker := time.NewTicker(1 * time.Millisecond)
    defer ticker.Stop()
    
    for range ticker.C {
        ptm.schedule()
    }
}

func (ptm *PeriodicTaskManager) schedule() {
    ptm.mu.Lock()
    defer ptm.mu.Unlock()
    
    // 获取下一个要执行的任务
    task := ptm.scheduler.Schedule()
    if task == nil {
        return
    }
    
    // 执行任务
    go ptm.executeTask(task)
}

func (ptm *PeriodicTaskManager) executeTask(task *PeriodicTask) {
    startTime := time.Now()
    
    // 执行任务
    if err := task.Handler(); err != nil {
        log.Printf("Task %s execution failed: %v", task.ID, err)
    }
    
    executionTime := time.Since(startTime)
    
    // 检查是否超过执行时间
    if executionTime > task.ExecutionTime {
        log.Printf("Task %s exceeded execution time: %v > %v", task.ID, executionTime, task.ExecutionTime)
    }
    
    // 标记任务完成
    ptm.mu.Lock()
    task.State = Completed
    ptm.mu.Unlock()
}

func (ptm *PeriodicTaskManager) CheckDeadlines() {
    ptm.mu.Lock()
    defer ptm.mu.Unlock()
    
    now := time.Now()
    for _, task := range ptm.tasks {
        if task.State == Ready || task.State == Running {
            deadline := task.LastRelease.Add(task.Deadline)
            if now.After(deadline) {
                task.MissedDeadlines++
                log.Printf("Task %s missed deadline", task.ID)
            }
        }
    }
}
```

### 3.2 非周期性任务

```go
// 非周期性任务管理器
type AperiodicTaskManager struct {
    tasks    map[string]*AperiodicTask
    scheduler Scheduler
    mu       sync.Mutex
}

type AperiodicTask struct {
    ID           string
    ArrivalTime  time.Time
    ExecutionTime time.Duration
    Deadline     time.Duration
    Priority     int
    Handler      TaskHandler
    State        TaskState
}

func (atm *AperiodicTaskManager) SubmitTask(task *AperiodicTask) {
    atm.mu.Lock()
    defer atm.mu.Unlock()
    
    task.ArrivalTime = time.Now()
    task.State = Ready
    
    atm.tasks[task.ID] = task
    atm.scheduler.AddTask(task)
}

func (atm *AperiodicTaskManager) RemoveTask(taskID string) {
    atm.mu.Lock()
    defer atm.mu.Unlock()
    
    delete(atm.tasks, taskID)
}

func (atm *AperiodicTaskManager) GetTask(taskID string) (*AperiodicTask, error) {
    atm.mu.RLock()
    defer atm.mu.RUnlock()
    
    task, exists := atm.tasks[taskID]
    if !exists {
        return nil, fmt.Errorf("task %s not found", taskID)
    }
    
    return task, nil
}
```

## 4. 时间约束管理

### 4.1 截止时间监控

```go
// 截止时间监控器
type DeadlineMonitor struct {
    tasks    map[string]*Task
    alerts   chan DeadlineAlert
    mu       sync.Mutex
}

type DeadlineAlert struct {
    TaskID    string
    Type      string // "approaching", "missed"
    Message   string
    Timestamp time.Time
}

func (dm *DeadlineMonitor) AddTask(task *Task) {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    dm.tasks[task.ID] = task
}

func (dm *DeadlineMonitor) StartMonitoring() {
    ticker := time.NewTicker(10 * time.Millisecond)
    defer ticker.Stop()
    
    for range ticker.C {
        dm.checkDeadlines()
    }
}

func (dm *DeadlineMonitor) checkDeadlines() {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    now := time.Now()
    for _, task := range dm.tasks {
        if task.State == Ready || task.State == Running {
            deadline := task.LastRelease.Add(task.Deadline)
            remainingTime := deadline.Sub(now)
            
            if remainingTime < 0 {
                // 截止时间已过
                alert := DeadlineAlert{
                    TaskID:    task.ID,
                    Type:      "missed",
                    Message:   fmt.Sprintf("Task %s missed deadline", task.ID),
                    Timestamp: now,
                }
                dm.alerts <- alert
            } else if remainingTime < task.ExecutionTime {
                // 接近截止时间
                alert := DeadlineAlert{
                    TaskID:    task.ID,
                    Type:      "approaching",
                    Message:   fmt.Sprintf("Task %s approaching deadline", task.ID),
                    Timestamp: now,
                }
                dm.alerts <- alert
            }
        }
    }
}

func (dm *DeadlineMonitor) ProcessAlerts() {
    for alert := range dm.alerts {
        dm.handleAlert(alert)
    }
}

func (dm *DeadlineMonitor) handleAlert(alert DeadlineAlert) {
    log.Printf("Deadline Alert: %s - %s", alert.Type, alert.Message)
    
    switch alert.Type {
    case "missed":
        dm.handleMissedDeadline(alert)
    case "approaching":
        dm.handleApproachingDeadline(alert)
    }
}

func (dm *DeadlineMonitor) handleMissedDeadline(alert DeadlineAlert) {
    // 处理错过截止时间的情况
    log.Printf("Handling missed deadline for task: %s", alert.TaskID)
    
    // 可以采取的措施：
    // 1. 记录错误日志
    // 2. 发送通知
    // 3. 调整系统参数
    // 4. 重启任务
}

func (dm *DeadlineMonitor) handleApproachingDeadline(alert DeadlineAlert) {
    // 处理接近截止时间的情况
    log.Printf("Handling approaching deadline for task: %s", alert.TaskID)
    
    // 可以采取的措施：
    // 1. 提高任务优先级
    // 2. 预分配资源
    // 3. 发送警告
}
```

### 4.2 响应时间分析

```go
// 响应时间分析器
type ResponseTimeAnalyzer struct {
    tasks    map[string]*Task
    mu       sync.Mutex
}

type ResponseTimeStats struct {
    TaskID       string
    MinResponse  time.Duration
    MaxResponse  time.Duration
    AvgResponse  time.Duration
    Jitter       time.Duration
    SampleCount  int
}

func (rta *ResponseTimeAnalyzer) RecordResponseTime(taskID string, responseTime time.Duration) {
    rta.mu.Lock()
    defer rta.mu.Unlock()
    
    task := rta.tasks[taskID]
    if task == nil {
        return
    }
    
    // 更新统计信息
    if task.ResponseTimeStats.MinResponse == 0 || responseTime < task.ResponseTimeStats.MinResponse {
        task.ResponseTimeStats.MinResponse = responseTime
    }
    
    if responseTime > task.ResponseTimeStats.MaxResponse {
        task.ResponseTimeStats.MaxResponse = responseTime
    }
    
    // 计算平均响应时间
    task.ResponseTimeStats.SampleCount++
    total := task.ResponseTimeStats.AvgResponse * time.Duration(task.ResponseTimeStats.SampleCount-1)
    task.ResponseTimeStats.AvgResponse = (total + responseTime) / time.Duration(task.ResponseTimeStats.SampleCount)
    
    // 计算抖动
    task.ResponseTimeStats.Jitter = task.ResponseTimeStats.MaxResponse - task.ResponseTimeStats.MinResponse
}

func (rta *ResponseTimeAnalyzer) GetStats(taskID string) (*ResponseTimeStats, error) {
    rta.mu.RLock()
    defer rta.mu.RUnlock()
    
    task := rta.tasks[taskID]
    if task == nil {
        return nil, fmt.Errorf("task %s not found", taskID)
    }
    
    return &task.ResponseTimeStats, nil
}

func (rta *ResponseTimeAnalyzer) CheckResponseTimeConstraints(taskID string) bool {
    stats, err := rta.GetStats(taskID)
    if err != nil {
        return false
    }
    
    task := rta.tasks[taskID]
    
    // 检查平均响应时间是否满足约束
    if stats.AvgResponse > task.Deadline {
        return false
    }
    
    // 检查最大响应时间是否满足约束
    if stats.MaxResponse > task.Deadline {
        return false
    }
    
    return true
}
```

## 5. 实时内存管理

### 5.1 内存池管理

```go
// 实时内存池
type RealTimeMemoryPool struct {
    pools    map[string]*MemoryPool
    mu       sync.Mutex
}

type MemoryPool struct {
    ID       string
    Size     int
    BlockSize int
    Blocks   []*MemoryBlock
    FreeList []*MemoryBlock
    mu       sync.Mutex
}

type MemoryBlock struct {
    ID       int
    Address  uintptr
    Size     int
    IsUsed   bool
    TaskID   string
    AllocTime time.Time
}

func (rtmp *RealTimeMemoryPool) CreatePool(id string, size, blockSize int) error {
    rtmp.mu.Lock()
    defer rtmp.mu.Unlock()
    
    pool := &MemoryPool{
        ID:        id,
        Size:      size,
        BlockSize: blockSize,
        Blocks:    make([]*MemoryBlock, 0),
        FreeList:  make([]*MemoryBlock, 0),
    }
    
    // 分配内存块
    numBlocks := size / blockSize
    for i := 0; i < numBlocks; i++ {
        block := &MemoryBlock{
            ID:     i,
            Size:   blockSize,
            IsUsed: false,
        }
        pool.Blocks = append(pool.Blocks, block)
        pool.FreeList = append(pool.FreeList, block)
    }
    
    rtmp.pools[id] = pool
    return nil
}

func (rtmp *RealTimeMemoryPool) Allocate(poolID, taskID string) (*MemoryBlock, error) {
    rtmp.mu.Lock()
    defer rtmp.mu.Unlock()
    
    pool := rtmp.pools[poolID]
    if pool == nil {
        return nil, fmt.Errorf("pool %s not found", poolID)
    }
    
    pool.mu.Lock()
    defer pool.mu.Unlock()
    
    if len(pool.FreeList) == 0 {
        return nil, fmt.Errorf("no free blocks in pool %s", poolID)
    }
    
    // 获取第一个可用块
    block := pool.FreeList[0]
    pool.FreeList = pool.FreeList[1:]
    
    block.IsUsed = true
    block.TaskID = taskID
    block.AllocTime = time.Now()
    
    return block, nil
}

func (rtmp *RealTimeMemoryPool) Deallocate(poolID string, block *MemoryBlock) error {
    rtmp.mu.Lock()
    defer rtmp.mu.Unlock()
    
    pool := rtmp.pools[poolID]
    if pool == nil {
        return fmt.Errorf("pool %s not found", poolID)
    }
    
    pool.mu.Lock()
    defer pool.mu.Unlock()
    
    block.IsUsed = false
    block.TaskID = ""
    block.AllocTime = time.Time{}
    
    pool.FreeList = append(pool.FreeList, block)
    
    return nil
}

func (rtmp *RealTimeMemoryPool) GetPoolStats(poolID string) (*PoolStats, error) {
    rtmp.mu.RLock()
    defer rtmp.mu.RUnlock()
    
    pool := rtmp.pools[poolID]
    if pool == nil {
        return nil, fmt.Errorf("pool %s not found", poolID)
    }
    
    pool.mu.Lock()
    defer pool.mu.Unlock()
    
    usedBlocks := 0
    for _, block := range pool.Blocks {
        if block.IsUsed {
            usedBlocks++
        }
    }
    
    return &PoolStats{
        TotalBlocks: len(pool.Blocks),
        UsedBlocks:  usedBlocks,
        FreeBlocks:  len(pool.FreeList),
        Utilization: float64(usedBlocks) / float64(len(pool.Blocks)),
    }, nil
}

type PoolStats struct {
    TotalBlocks int
    UsedBlocks  int
    FreeBlocks  int
    Utilization float64
}
```

### 5.2 垃圾回收控制

```go
// 实时垃圾回收控制器
type RealTimeGCController struct {
    enabled  bool
    interval time.Duration
    mu       sync.Mutex
}

func (rtgc *RealTimeGCController) Enable() {
    rtgc.mu.Lock()
    defer rtgc.mu.Unlock()
    
    rtgc.enabled = true
    go rtgc.run()
}

func (rtgc *RealTimeGCController) Disable() {
    rtgc.mu.Lock()
    defer rtgc.mu.Unlock()
    
    rtgc.enabled = false
}

func (rtgc *RealTimeGCController) run() {
    ticker := time.NewTicker(rtgc.interval)
    defer ticker.Stop()
    
    for range ticker.C {
        rtgc.mu.Lock()
        if !rtgc.enabled {
            rtgc.mu.Unlock()
            return
        }
        rtgc.mu.Unlock()
        
        rtgc.performGC()
    }
}

func (rtgc *RealTimeGCController) performGC() {
    // 设置GC目标
    runtime.GC()
    
    // 等待GC完成
    runtime.Gosched()
}

func (rtgc *RealTimeGCController) SetInterval(interval time.Duration) {
    rtgc.mu.Lock()
    defer rtgc.mu.Unlock()
    
    rtgc.interval = interval
}

func (rtgc *RealTimeGCController) GetStats() *GCStats {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)
    
    return &GCStats{
        Alloc:      stats.Alloc,
        TotalAlloc: stats.TotalAlloc,
        Sys:        stats.Sys,
        NumGC:      stats.NumGC,
        PauseTotalNs: stats.PauseTotalNs,
    }
}

type GCStats struct {
    Alloc        uint64
    TotalAlloc   uint64
    Sys          uint64
    NumGC        uint32
    PauseTotalNs uint64
}
```

## 6. 实时通信

### 6.1 实时消息队列

```go
// 实时消息队列
type RealTimeMessageQueue struct {
    queues   map[string]*MessageQueue
    mu       sync.Mutex
}

type MessageQueue struct {
    ID       string
    Capacity int
    Messages []*Message
    Head     int
    Tail     int
    Size     int
    mu       sync.Mutex
}

type Message struct {
    ID        string
    Priority  int
    Data      interface{}
    Timestamp time.Time
    Deadline  time.Time
}

func (rtmq *RealTimeMessageQueue) CreateQueue(id string, capacity int) error {
    rtmq.mu.Lock()
    defer rtmq.mu.Unlock()
    
    queue := &MessageQueue{
        ID:       id,
        Capacity: capacity,
        Messages: make([]*Message, capacity),
    }
    
    rtmq.queues[id] = queue
    return nil
}

func (rtmq *RealTimeMessageQueue) SendMessage(queueID string, message *Message) error {
    rtmq.mu.Lock()
    defer rtmq.mu.Unlock()
    
    queue := rtmq.queues[queueID]
    if queue == nil {
        return fmt.Errorf("queue %s not found", queueID)
    }
    
    queue.mu.Lock()
    defer queue.mu.Unlock()
    
    if queue.Size >= queue.Capacity {
        return fmt.Errorf("queue %s is full", queueID)
    }
    
    // 插入消息（按优先级排序）
    insertIndex := queue.findInsertIndex(message.Priority)
    queue.insertMessage(insertIndex, message)
    
    return nil
}

func (mq *MessageQueue) findInsertIndex(priority int) int {
    for i := 0; i < mq.Size; i++ {
        index := (mq.Head + i) % mq.Capacity
        if mq.Messages[index].Priority < priority {
            return index
        }
    }
    return (mq.Head + mq.Size) % mq.Capacity
}

func (mq *MessageQueue) insertMessage(index int, message *Message) {
    // 移动现有消息
    for i := mq.Size - 1; i >= 0; i-- {
        currentIndex := (mq.Head + i) % mq.Capacity
        nextIndex := (mq.Head + i + 1) % mq.Capacity
        mq.Messages[nextIndex] = mq.Messages[currentIndex]
    }
    
    mq.Messages[index] = message
    mq.Size++
}

func (rtmq *RealTimeMessageQueue) ReceiveMessage(queueID string) (*Message, error) {
    rtmq.mu.Lock()
    defer rtmq.mu.Unlock()
    
    queue := rtmq.queues[queueID]
    if queue == nil {
        return nil, fmt.Errorf("queue %s not found", queueID)
    }
    
    queue.mu.Lock()
    defer queue.mu.Unlock()
    
    if queue.Size == 0 {
        return nil, fmt.Errorf("queue %s is empty", queueID)
    }
    
    // 获取最高优先级的消息
    message := queue.Messages[queue.Head]
    queue.Messages[queue.Head] = nil
    queue.Head = (queue.Head + 1) % queue.Capacity
    queue.Size--
    
    return message, nil
}

func (rtmq *RealTimeMessageQueue) CheckDeadlines() []*Message {
    rtmq.mu.Lock()
    defer rtmq.mu.Unlock()
    
    now := time.Now()
    expiredMessages := make([]*Message, 0)
    
    for _, queue := range rtmq.queues {
        queue.mu.Lock()
        for i := 0; i < queue.Size; i++ {
            index := (queue.Head + i) % queue.Capacity
            message := queue.Messages[index]
            if now.After(message.Deadline) {
                expiredMessages = append(expiredMessages, message)
            }
        }
        queue.mu.Unlock()
    }
    
    return expiredMessages
}
```

## 总结

实时系统是一个复杂的系统，需要在严格的时间约束下保证系统的正确性和可靠性。

**关键设计原则**：

1. 可预测性：系统行为必须是可预测的
2. 确定性：相同输入必须产生相同输出
3. 时间约束：必须在规定时间内完成
4. 资源管理：有效管理CPU、内存等资源
5. 错误处理：快速检测和处理错误

**常见挑战**：

1. 调度算法选择
2. 时间约束满足
3. 资源竞争处理
4. 优先级反转
5. 内存管理

实时系统的设计需要根据具体的应用场景和性能要求来选择合适的算法和技术，确保系统能够满足严格的时间约束。
