# 11.8.3 IoT边缘计算

## 11.8.3.1 概述

IoT边缘计算是将计算能力从云端推向网络边缘，在靠近数据源的地方进行数据处理和分析的技术。

### 11.8.3.1.1 基本概念

**定义 11.8.3.1** (边缘计算)
边缘计算是一种分布式计算范式，将计算任务从集中式数据中心转移到网络边缘，在靠近数据源的地方进行数据处理。

**定义 11.8.3.2** (边缘节点)
边缘节点是部署在网络边缘的计算设备，具有有限的计算、存储和网络资源。

### 11.8.3.1.2 边缘计算层次

```go
// 边缘计算层次
type EdgeLayer int

const (
    DeviceLayer EdgeLayer = iota    // 设备层
    EdgeLayer1                      // 边缘层1
    EdgeLayer2                      // 边缘层2
    CloudLayer                      // 云层
)

// 计算资源类型
type ResourceType int

const (
    CPU ResourceType = iota
    Memory
    Storage
    Network
)
```

## 11.8.3.2 边缘计算架构

### 11.8.3.2.1 架构模型

**定义 11.8.3.3** (边缘计算架构)
边缘计算架构定义了边缘节点、云中心和IoT设备之间的协作关系。

**定理 11.8.3.1** (延迟优化)
边缘计算可以将端到端延迟降低50-80%。

**证明**:
设云端处理延迟为 ```latex
$T_{cloud}$
```，网络传输延迟为 ```latex
$T_{network}$
```，边缘处理延迟为 ```latex
$T_{edge}$
```，则：
$```latex
$T_{total} = T_{edge} + T_{network\_edge} \ll T_{cloud} + T_{network\_cloud}$
```$

### 11.8.3.2.2 Go实现边缘架构

```go
// 边缘计算架构
type EdgeArchitecture struct {
    nodes       map[string]*EdgeNode
    scheduler   *TaskScheduler
    resourceMgr *ResourceManager
    dataMgr     *DataManager
}

// 边缘节点
type EdgeNode struct {
    ID          string
    Layer       EdgeLayer
    Resources   *NodeResources
    Tasks       []*Task
    Status      NodeStatus
    Location    *Location
    LastUpdate  time.Time
}

// 节点资源
type NodeResources struct {
    CPU     float64 // CPU使用率
    Memory  float64 // 内存使用率
    Storage float64 // 存储使用率
    Network float64 // 网络带宽使用率
}

// 节点状态
type NodeStatus int

const (
    Online NodeStatus = iota
    Offline
    Overloaded
    Maintenance
)

// 位置信息
type Location struct {
    Latitude  float64
    Longitude float64
    Region    string
}

// 任务
type Task struct {
    ID          string
    Name        string
    Priority    int
    Resources   *TaskResources
    Status      TaskStatus
    Created     time.Time
    Started     time.Time
    Completed   time.Time
    Result      interface{}
}

// 任务资源需求
type TaskResources struct {
    CPU     float64
    Memory  float64
    Storage float64
    Network float64
}

// 任务状态
type TaskStatus int

const (
    Pending TaskStatus = iota
    Running
    Completed
    Failed
    Cancelled
)

// 创建边缘架构
func NewEdgeArchitecture() *EdgeArchitecture {
    return &EdgeArchitecture{
        nodes:       make(map[string]*EdgeNode),
        scheduler:   NewTaskScheduler(),
        resourceMgr: NewResourceManager(),
        dataMgr:     NewDataManager(),
    }
}

// 添加边缘节点
func (ea *EdgeArchitecture) AddNode(node *EdgeNode) {
    ea.nodes[node.ID] = node
    ea.resourceMgr.RegisterNode(node)
}

// 获取可用节点
func (ea *EdgeArchitecture) GetAvailableNodes() []*EdgeNode {
    var available []*EdgeNode
    for _, node := range ea.nodes {
        if node.Status == Online && ea.resourceMgr.HasCapacity(node.ID) {
            available = append(available, node)
        }
    }
    return available
}

// 提交任务
func (ea *EdgeArchitecture) SubmitTask(task *Task) error {
    return ea.scheduler.ScheduleTask(task, ea.nodes)
}

// 获取架构状态
func (ea *EdgeArchitecture) GetStatus() *ArchitectureStatus {
    status := &ArchitectureStatus{
        TotalNodes:    len(ea.nodes),
        OnlineNodes:   0,
        TotalTasks:    0,
        RunningTasks:  0,
        ResourceUsage: make(map[string]float64),
    }
    
    for _, node := range ea.nodes {
        if node.Status == Online {
            status.OnlineNodes++
        }
        status.TotalTasks += len(node.Tasks)
        
        for _, task := range node.Tasks {
            if task.Status == Running {
                status.RunningTasks++
            }
        }
    }
    
    status.ResourceUsage = ea.resourceMgr.GetOverallUsage()
    
    return status
}

// 架构状态
type ArchitectureStatus struct {
    TotalNodes    int
    OnlineNodes   int
    TotalTasks    int
    RunningTasks  int
    ResourceUsage map[string]float64
}
```

## 11.8.3.3 资源管理

### 11.8.3.3.1 资源分配策略

**定义 11.8.3.4** (资源管理)
资源管理是优化边缘节点计算、存储和网络资源分配的过程。

**定理 11.8.3.2** (资源利用率)
最优资源分配可以最大化系统整体利用率。

### 11.8.3.3.2 Go实现资源管理

```go
// 资源管理器
type ResourceManager struct {
    nodes      map[string]*NodeResources
    policies   map[string]*ResourcePolicy
    monitoring *ResourceMonitor
}

// 资源策略
type ResourcePolicy struct {
    ID              string
    Name            string
    MaxCPUUsage     float64
    MaxMemoryUsage  float64
    MaxStorageUsage float64
    MaxNetworkUsage float64
    Priority        int
}

// 资源监控器
type ResourceMonitor struct {
    metrics map[string]*ResourceMetrics
    history []*ResourceSnapshot
}

// 资源指标
type ResourceMetrics struct {
    NodeID       string
    Timestamp    time.Time
    CPUUsage     float64
    MemoryUsage  float64
    StorageUsage float64
    NetworkUsage float64
}

// 创建资源管理器
func NewResourceManager() *ResourceManager {
    return &ResourceManager{
        nodes:      make(map[string]*NodeResources),
        policies:   make(map[string]*ResourcePolicy),
        monitoring: &ResourceMonitor{
            metrics: make(map[string]*ResourceMetrics),
            history: make([]*ResourceSnapshot, 0),
        },
    }
}

// 注册节点
func (rm *ResourceManager) RegisterNode(node *EdgeNode) {
    rm.nodes[node.ID] = node.Resources
}

// 检查资源容量
func (rm *ResourceManager) HasCapacity(nodeID string) bool {
    resources, exists := rm.nodes[nodeID]
    if !exists {
        return false
    }
    
    return resources.CPU < 0.9 && 
           resources.Memory < 0.9 && 
           resources.Storage < 0.9 && 
           resources.Network < 0.9
}

// 分配资源
func (rm *ResourceManager) AllocateResources(nodeID string, taskResources *TaskResources) bool {
    resources, exists := rm.nodes[nodeID]
    if !exists {
        return false
    }
    
    if resources.CPU+taskResources.CPU > 1.0 ||
       resources.Memory+taskResources.Memory > 1.0 ||
       resources.Storage+taskResources.Storage > 1.0 ||
       resources.Network+taskResources.Network > 1.0 {
        return false
    }
    
    resources.CPU += taskResources.CPU
    resources.Memory += taskResources.Memory
    resources.Storage += taskResources.Storage
    resources.Network += taskResources.Network
    
    return true
}

// 释放资源
func (rm *ResourceManager) ReleaseResources(nodeID string, taskResources *TaskResources) {
    resources, exists := rm.nodes[nodeID]
    if !exists {
        return
    }
    
    resources.CPU -= taskResources.CPU
    resources.Memory -= taskResources.Memory
    resources.Storage -= taskResources.Storage
    resources.Network -= taskResources.Network
    
    if resources.CPU < 0 {
        resources.CPU = 0
    }
    if resources.Memory < 0 {
        resources.Memory = 0
    }
    if resources.Storage < 0 {
        resources.Storage = 0
    }
    if resources.Network < 0 {
        resources.Network = 0
    }
}

// 获取整体使用率
func (rm *ResourceManager) GetOverallUsage() map[string]float64 {
    usage := make(map[string]float64)
    
    if len(rm.nodes) == 0 {
        return usage
    }
    
    totalCPU := 0.0
    totalMemory := 0.0
    totalStorage := 0.0
    totalNetwork := 0.0
    
    for _, resources := range rm.nodes {
        totalCPU += resources.CPU
        totalMemory += resources.Memory
        totalStorage += resources.Storage
        totalNetwork += resources.Network
    }
    
    nodeCount := float64(len(rm.nodes))
    usage["CPU"] = totalCPU / nodeCount
    usage["Memory"] = totalMemory / nodeCount
    usage["Storage"] = totalStorage / nodeCount
    usage["Network"] = totalNetwork / nodeCount
    
    return usage
}
```

## 11.8.3.4 任务调度

### 11.8.3.4.1 调度算法

**定义 11.8.3.5** (任务调度)
任务调度是将计算任务分配给合适的边缘节点的过程。

**定理 11.8.3.3** (调度优化)
最优调度可以最小化任务完成时间和资源消耗。

### 11.8.3.4.2 Go实现任务调度

```go
// 任务调度器
type TaskScheduler struct {
    queue           *TaskQueue
    strategies      map[string]SchedulingStrategy
    currentStrategy string
}

// 任务队列
type TaskQueue struct {
    tasks []*Task
    mutex sync.Mutex
}

// 调度策略接口
type SchedulingStrategy interface {
    Schedule(tasks []*Task, nodes map[string]*EdgeNode) map[string]*Task
    Name() string
}

// 先来先服务策略
type FIFOStrategy struct{}

func (f *FIFOStrategy) Schedule(tasks []*Task, nodes map[string]*EdgeNode) map[string]*Task {
    assignments := make(map[string]*Task)
    availableNodes := getAvailableNodes(nodes)
    
    if len(availableNodes) == 0 {
        return assignments
    }
    
    nodeIndex := 0
    for _, task := range tasks {
        if nodeIndex >= len(availableNodes) {
            break
        }
        
        node := availableNodes[nodeIndex]
        if canAssignTask(task, node) {
            assignments[node.ID] = task
            nodeIndex++
        }
    }
    
    return assignments
}

func (f *FIFOStrategy) Name() string {
    return "FIFO"
}

// 负载均衡策略
type LoadBalancingStrategy struct{}

func (l *LoadBalancingStrategy) Schedule(tasks []*Task, nodes map[string]*EdgeNode) map[string]*Task {
    assignments := make(map[string]*Task)
    availableNodes := getAvailableNodes(nodes)
    
    if len(availableNodes) == 0 {
        return assignments
    }
    
    sort.Slice(availableNodes, func(i, j int) bool {
        return getNodeLoad(availableNodes[i]) < getNodeLoad(availableNodes[j])
    })
    
    nodeIndex := 0
    for _, task := range tasks {
        if nodeIndex >= len(availableNodes) {
            break
        }
        
        node := availableNodes[nodeIndex]
        if canAssignTask(task, node) {
            assignments[node.ID] = task
            nodeIndex++
        }
    }
    
    return assignments
}

func (l *LoadBalancingStrategy) Name() string {
    return "LoadBalancing"
}

// 创建任务调度器
func NewTaskScheduler() *TaskScheduler {
    scheduler := &TaskScheduler{
        queue: &TaskQueue{
            tasks: make([]*Task, 0),
        },
        strategies: make(map[string]SchedulingStrategy),
    }
    
    scheduler.strategies["FIFO"] = &FIFOStrategy{}
    scheduler.strategies["LoadBalancing"] = &LoadBalancingStrategy{}
    scheduler.currentStrategy = "LoadBalancing"
    
    return scheduler
}

// 调度任务
func (ts *TaskScheduler) ScheduleTask(task *Task, nodes map[string]*EdgeNode) error {
    strategy, exists := ts.strategies[ts.currentStrategy]
    if !exists {
        return fmt.Errorf("current strategy not found")
    }
    
    ts.queue.mutex.Lock()
    pendingTasks := make([]*Task, len(ts.queue.tasks))
    copy(pendingTasks, ts.queue.tasks)
    ts.queue.mutex.Unlock()
    
    assignments := strategy.Schedule(pendingTasks, nodes)
    
    for nodeID, assignedTask := range assignments {
        if node, exists := nodes[nodeID]; exists {
            node.Tasks = append(node.Tasks, assignedTask)
            assignedTask.Status = Running
            assignedTask.Started = time.Now()
        }
    }
    
    return nil
}

// 辅助函数
func getAvailableNodes(nodes map[string]*EdgeNode) []*EdgeNode {
    var available []*EdgeNode
    for _, node := range nodes {
        if node.Status == Online {
            available = append(available, node)
        }
    }
    return available
}

func canAssignTask(task *Task, node *EdgeNode) bool {
    return node.Resources.CPU+task.Resources.CPU <= 1.0 &&
           node.Resources.Memory+task.Resources.Memory <= 1.0 &&
           node.Resources.Storage+task.Resources.Storage <= 1.0 &&
           node.Resources.Network+task.Resources.Network <= 1.0
}

func getNodeLoad(node *EdgeNode) float64 {
    return (node.Resources.CPU + node.Resources.Memory + 
            node.Resources.Storage + node.Resources.Network) / 4.0
}
```

## 11.8.3.5 数据管理

### 11.8.3.5.1 数据流管理

**定义 11.8.3.6** (数据管理)
数据管理是处理边缘计算中数据存储、传输和处理的策略。

### 11.8.3.5.2 Go实现数据管理

```go
// 数据管理器
type DataManager struct {
    storage  map[string]*DataStorage
    cache    *DataCache
    pipeline *DataPipeline
}

// 数据存储
type DataStorage struct {
    ID       string
    Type     StorageType
    Capacity int64
    Used     int64
    Data     map[string][]byte
}

// 存储类型
type StorageType int

const (
    LocalStorage StorageType = iota
    DistributedStorage
    CloudStorage
)

// 数据缓存
type DataCache struct {
    cache   map[string]*CacheEntry
    size    int
    maxSize int
}

// 缓存条目
type CacheEntry struct {
    Key       string
    Value     []byte
    Timestamp time.Time
    TTL       time.Duration
}

// 数据管道
type DataPipeline struct {
    stages []*PipelineStage
    data   chan *DataPacket
}

// 管道阶段
type PipelineStage struct {
    ID       string
    Name     string
    Function func(*DataPacket) *DataPacket
    Next     *PipelineStage
}

// 数据包
type DataPacket struct {
    ID       string
    Data     []byte
    Metadata map[string]interface{}
    Source   string
    Target   string
    Timestamp time.Time
}

// 创建数据管理器
func NewDataManager() *DataManager {
    return &DataManager{
        storage: make(map[string]*DataStorage),
        cache: &DataCache{
            cache:   make(map[string]*CacheEntry),
            maxSize: 1000,
        },
        pipeline: &DataPipeline{
            stages: make([]*PipelineStage, 0),
            data:   make(chan *DataPacket, 100),
        },
    }
}

// 存储数据
func (dm *DataManager) StoreData(storageID string, key string, data []byte) error {
    storage, exists := dm.storage[storageID]
    if !exists {
        return fmt.Errorf("storage %s not found", storageID)
    }
    
    if int64(len(data)) > storage.Capacity-storage.Used {
        return fmt.Errorf("insufficient storage space")
    }
    
    storage.Data[key] = data
    storage.Used += int64(len(data))
    
    return nil
}

// 获取数据
func (dm *DataManager) GetData(storageID string, key string) ([]byte, error) {
    if cached, exists := dm.cache.Get(key); exists {
        return cached, nil
    }
    
    storage, exists := dm.storage[storageID]
    if !exists {
        return nil, fmt.Errorf("storage %s not found", storageID)
    }
    
    data, exists := storage.Data[key]
    if !exists {
        return nil, fmt.Errorf("data not found")
    }
    
    dm.cache.Put(key, data)
    
    return data, nil
}

// 缓存操作
func (dc *DataCache) Get(key string) ([]byte, bool) {
    entry, exists := dc.cache[key]
    if !exists {
        return nil, false
    }
    
    if time.Since(entry.Timestamp) > entry.TTL {
        delete(dc.cache, key)
        dc.size--
        return nil, false
    }
    
    return entry.Value, true
}

func (dc *DataCache) Put(key string, value []byte) {
    if dc.size >= dc.maxSize {
        dc.evictOldest()
    }
    
    dc.cache[key] = &CacheEntry{
        Key:       key,
        Value:     value,
        Timestamp: time.Now(),
        TTL:       5 * time.Minute,
    }
    dc.size++
}

func (dc *DataCache) evictOldest() {
    var oldestKey string
    var oldestTime time.Time
    
    for key, entry := range dc.cache {
        if oldestKey == "" || entry.Timestamp.Before(oldestTime) {
            oldestKey = key
            oldestTime = entry.Timestamp
        }
    }
    
    if oldestKey != "" {
        delete(dc.cache, oldestKey)
        dc.size--
    }
}
```

## 11.8.3.6 总结

本章详细介绍了IoT边缘计算的核心概念和技术，包括：

1. **边缘计算架构**: 层次结构、组件设计
2. **资源管理**: 资源分配、监控优化
3. **任务调度**: 调度算法、策略实现
4. **数据管理**: 存储、缓存、管道处理

通过Go语言实现，展示了边缘计算技术的核心思想和实际应用。

---

**相关链接**:

- [11.8.1 IoT基础理论](../01-IoT-Fundamentals/README.md)
- [11.8.2 IoT安全](../02-IoT-Security/README.md)
- [11.8.4 IoT应用](../04-IoT-Applications/README.md)
- [11.9 人工智能](../09-Artificial-Intelligence/README.md)
