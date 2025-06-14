# 01-模型训练平台

(Model Training Platform)

## 概述

模型训练平台是人工智能和机器学习系统的核心组件，负责管理、调度和执行机器学习模型的训练任务。本文档提供基于Go语言的模型训练平台架构设计和实现方案。

## 目录

- [01-模型训练平台](#01-模型训练平台)
  - [概述](#概述)
  - [目录](#目录)
  - [1. 形式化定义](#1-形式化定义)
    - [1.1 模型训练平台定义](#11-模型训练平台定义)
    - [1.2 训练任务定义](#12-训练任务定义)
    - [1.3 模型性能评估](#13-模型性能评估)
  - [2. 数学建模](#2-数学建模)
    - [2.1 梯度下降算法](#21-梯度下降算法)
    - [2.2 资源调度优化](#22-资源调度优化)
  - [3. 架构设计](#3-架构设计)
    - [3.1 系统架构图](#31-系统架构图)
    - [3.2 组件职责](#32-组件职责)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 训练任务模型](#41-训练任务模型)
    - [4.2 任务调度服务](#42-任务调度服务)
    - [4.3 资源管理服务](#43-资源管理服务)
    - [4.4 实验管理服务](#44-实验管理服务)
  - [5. 性能优化](#5-性能优化)
    - [5.1 分布式训练优化](#51-分布式训练优化)
    - [5.2 内存优化](#52-内存优化)
  - [6. 分布式训练](#6-分布式训练)
    - [6.1 通信协议](#61-通信协议)
  - [总结](#总结)

## 1. 形式化定义

### 1.1 模型训练平台定义

**定义 1.1** 模型训练平台 (Model Training Platform)
模型训练平台是一个六元组 $MTP = (M, D, T, R, S, P)$，其中：

- $M = \{m_1, m_2, ..., m_n\}$ 是模型集合
- $D = \{d_1, d_2, ..., d_k\}$ 是数据集集合
- $T = \{t_1, t_2, ..., t_l\}$ 是训练任务集合
- $R = \{r_1, r_2, ..., r_m\}$ 是计算资源集合
- $S = \{s_1, s_2, ..., s_o\}$ 是调度策略集合
- $P = \{p_1, p_2, ..., p_q\}$ 是性能指标集合

### 1.2 训练任务定义

**定义 1.2** 训练任务 (Training Task)
训练任务是一个五元组 $T = (m, d, h, r, c)$，其中：

- $m \in M$ 是要训练的模型
- $d \in D$ 是训练数据集
- $h$ 是超参数集合
- $r \in R$ 是分配的计算资源
- $c$ 是训练配置

### 1.3 模型性能评估

**定义 1.3** 模型性能函数
模型性能函数定义为：
$f: M \times D \rightarrow P$

其中 $f(m, d)$ 表示模型 $m$ 在数据集 $d$ 上的性能指标。

## 2. 数学建模

### 2.1 梯度下降算法

**定理 2.1** 梯度下降收敛性
对于损失函数 $L(\theta)$，如果满足Lipschitz条件：
$\|\nabla L(\theta_1) - \nabla L(\theta_2)\| \leq L\|\theta_1 - \theta_2\|$

则梯度下降算法以步长 $\eta = \frac{1}{L}$ 收敛到局部最优解。

**证明**：
设 $\theta_t$ 为第 $t$ 次迭代的参数，则：
$\theta_{t+1} = \theta_t - \eta \nabla L(\theta_t)$

根据Lipschitz条件：
$L(\theta_{t+1}) \leq L(\theta_t) + \nabla L(\theta_t)^T(\theta_{t+1} - \theta_t) + \frac{L}{2}\|\theta_{t+1} - \theta_t\|^2$

代入梯度下降更新规则：
$L(\theta_{t+1}) \leq L(\theta_t) - \eta\|\nabla L(\theta_t)\|^2 + \frac{L\eta^2}{2}\|\nabla L(\theta_t)\|^2$

当 $\eta = \frac{1}{L}$ 时：
$L(\theta_{t+1}) \leq L(\theta_t) - \frac{1}{2L}\|\nabla L(\theta_t)\|^2$

因此损失函数单调递减，算法收敛。

### 2.2 资源调度优化

**定义 2.2** 资源调度函数
资源调度函数 $S: T \times R \rightarrow R$ 定义为：

$$S(t, R) = \arg\max_{r \in R} \left( \frac{Performance(t, r)}{Cost(r)} \right)$$

其中 $Performance(t, r)$ 是任务 $t$ 在资源 $r$ 上的性能，$Cost(r)$ 是资源 $r$ 的成本。

## 3. 架构设计

### 3.1 系统架构图

```text
┌─────────────────────────────────────────────────────────────┐
│                    模型训练平台架构                           │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  任务调度   │  │  资源管理   │  │  模型管理   │         │
│  │  服务       │  │  服务       │  │  服务       │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  数据管理   │  │  实验管理   │  │  监控告警   │         │
│  │  服务       │  │  服务       │  │  服务       │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  分布式训练 │  │  模型版本   │  │  性能评估   │         │
│  │  引擎       │  │  管理       │  │  服务       │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 组件职责

| 组件 | 职责 | 技术栈 |
|------|------|--------|
| 任务调度服务 | 训练任务调度、优先级管理 | Go + Kubernetes + Redis |
| 资源管理服务 | 计算资源分配、监控 | Go + Prometheus + Grafana |
| 模型管理服务 | 模型版本控制、存储 | Go + MinIO + PostgreSQL |
| 数据管理服务 | 数据集管理、预处理 | Go + Apache Arrow + Parquet |
| 实验管理服务 | 实验跟踪、超参数调优 | Go + MLflow + Weights & Biases |
| 监控告警服务 | 训练监控、异常告警 | Go + Prometheus + AlertManager |
| 分布式训练引擎 | 多机多卡训练 | Go + Horovod + TensorFlow |
| 模型版本管理 | 模型版本控制、部署 | Go + DVC + Docker |
| 性能评估服务 | 模型评估、指标计算 | Go + scikit-learn + NumPy |

## 4. Go语言实现

### 4.1 训练任务模型

```go
// TrainingTask 训练任务模型
type TrainingTask struct {
    ID          string                 `json:"id" validate:"required"`
    Name        string                 `json:"name" validate:"required"`
    ModelID     string                 `json:"model_id" validate:"required"`
    DatasetID   string                 `json:"dataset_id" validate:"required"`
    HyperParams map[string]interface{} `json:"hyper_params"`
    Config      TrainingConfig         `json:"config"`
    Status      TaskStatus             `json:"status"`
    Progress    float64                `json:"progress"`
    CreatedAt   time.Time              `json:"created_at"`
    StartedAt   *time.Time             `json:"started_at"`
    CompletedAt *time.Time             `json:"completed_at"`
    Metrics     map[string]float64     `json:"metrics"`
    Logs        []LogEntry             `json:"logs"`
    Resources   ResourceAllocation     `json:"resources"`
}

// TaskStatus 任务状态枚举
type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusRunning    TaskStatus = "running"
    TaskStatusCompleted  TaskStatus = "completed"
    TaskStatusFailed     TaskStatus = "failed"
    TaskStatusCancelled  TaskStatus = "cancelled"
)

// TrainingConfig 训练配置
type TrainingConfig struct {
    Epochs           int     `json:"epochs" validate:"min=1"`
    BatchSize        int     `json:"batch_size" validate:"min=1"`
    LearningRate     float64 `json:"learning_rate" validate:"min=0"`
    Optimizer        string  `json:"optimizer" validate:"required"`
    LossFunction     string  `json:"loss_function" validate:"required"`
    ValidationSplit  float64 `json:"validation_split" validate:"min=0,max=1"`
    EarlyStopping    bool    `json:"early_stopping"`
    Patience         int     `json:"patience"`
    CheckpointEvery  int     `json:"checkpoint_every"`
    MaxRuntime       int     `json:"max_runtime"` // 秒
}

// ResourceAllocation 资源分配
type ResourceAllocation struct {
    CPU        float64 `json:"cpu" validate:"min=0"`
    Memory     int64   `json:"memory" validate:"min=0"` // MB
    GPU        int     `json:"gpu" validate:"min=0"`
    GPUMemory  int64   `json:"gpu_memory" validate:"min=0"` // MB
    Storage    int64   `json:"storage" validate:"min=0"` // GB
}

// LogEntry 日志条目
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Metrics   map[string]float64 `json:"metrics,omitempty"`
}
```

### 4.2 任务调度服务

```go
// TaskScheduler 任务调度器
type TaskScheduler struct {
    db          *gorm.DB
    redis       *redis.Client
    k8sClient   *kubernetes.Clientset
    resourceMgr *ResourceManager
    logger      *zap.Logger
    taskQueue   chan *TrainingTask
    stopChan    chan struct{}
    wg          sync.WaitGroup
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(db *gorm.DB, redis *redis.Client, k8sClient *kubernetes.Clientset, resourceMgr *ResourceManager) *TaskScheduler {
    return &TaskScheduler{
        db:          db,
        redis:       redis,
        k8sClient:   k8sClient,
        resourceMgr: resourceMgr,
        logger:      zap.L().Named("task_scheduler"),
        taskQueue:   make(chan *TrainingTask, 1000),
        stopChan:    make(chan struct{}),
    }
}

// Start 启动调度器
func (ts *TaskScheduler) Start() error {
    ts.logger.Info("starting task scheduler")

    // 启动任务处理协程
    ts.wg.Add(1)
    go ts.processTaskQueue()

    // 启动资源监控协程
    ts.wg.Add(1)
    go ts.monitorResources()

    // 启动任务状态同步协程
    ts.wg.Add(1)
    go ts.syncTaskStatus()

    return nil
}

// Stop 停止调度器
func (ts *TaskScheduler) Stop() error {
    ts.logger.Info("stopping task scheduler")
    close(ts.stopChan)
    ts.wg.Wait()
    return nil
}

// SubmitTask 提交训练任务
func (ts *TaskScheduler) SubmitTask(ctx context.Context, task *TrainingTask) error {
    // 验证任务
    if err := ts.validateTask(task); err != nil {
        return fmt.Errorf("invalid task: %w", err)
    }

    // 生成任务ID
    task.ID = ts.generateTaskID()

    // 设置初始状态
    task.Status = TaskStatusPending
    task.CreatedAt = time.Now()
    task.Progress = 0.0

    // 保存到数据库
    if err := ts.db.Create(task).Error; err != nil {
        return fmt.Errorf("failed to save task: %w", err)
    }

    // 添加到任务队列
    select {
    case ts.taskQueue <- task:
        ts.logger.Info("task submitted successfully",
            zap.String("task_id", task.ID),
            zap.String("task_name", task.Name))
    default:
        return fmt.Errorf("task queue is full")
    }

    return nil
}

// processTaskQueue 处理任务队列
func (ts *TaskScheduler) processTaskQueue() {
    defer ts.wg.Done()

    for {
        select {
        case <-ts.stopChan:
            return
        case task := <-ts.taskQueue:
            ts.processTask(task)
        }
    }
}

// processTask 处理单个任务
func (ts *TaskScheduler) processTask(task *TrainingTask) {
    ts.logger.Info("processing task",
        zap.String("task_id", task.ID),
        zap.String("task_name", task.Name))

    // 检查资源可用性
    if !ts.resourceMgr.HasAvailableResources(task.Resources) {
        ts.logger.Info("insufficient resources, re-queuing task",
            zap.String("task_id", task.ID))
        
        // 重新加入队列
        go func() {
            time.Sleep(30 * time.Second)
            ts.taskQueue <- task
        }()
        return
    }

    // 分配资源
    if err := ts.resourceMgr.AllocateResources(task.ID, task.Resources); err != nil {
        ts.logger.Error("failed to allocate resources",
            zap.String("task_id", task.ID),
            zap.Error(err))
        ts.updateTaskStatus(task.ID, TaskStatusFailed, err.Error())
        return
    }

    // 创建Kubernetes作业
    if err := ts.createTrainingJob(task); err != nil {
        ts.logger.Error("failed to create training job",
            zap.String("task_id", task.ID),
            zap.Error(err))
        ts.resourceMgr.ReleaseResources(task.ID)
        ts.updateTaskStatus(task.ID, TaskStatusFailed, err.Error())
        return
    }

    // 更新任务状态
    now := time.Now()
    task.StartedAt = &now
    task.Status = TaskStatusRunning
    ts.updateTask(task)
}

// createTrainingJob 创建训练作业
func (ts *TaskScheduler) createTrainingJob(task *TrainingTask) error {
    // 构建容器镜像
    imageName := ts.buildTrainingImage(task)

    // 创建Kubernetes作业
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("training-%s", task.ID),
            Namespace: "ml-training",
            Labels: map[string]string{
                "task-id":   task.ID,
                "task-name": task.Name,
            },
        },
        Spec: batchv1.JobSpec{
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "task-id":   task.ID,
                        "task-name": task.Name,
                    },
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "training",
                            Image: imageName,
                            Env: []corev1.EnvVar{
                                {Name: "TASK_ID", Value: task.ID},
                                {Name: "MODEL_ID", Value: task.ModelID},
                                {Name: "DATASET_ID", Value: task.DatasetID},
                                {Name: "HYPER_PARAMS", Value: ts.serializeHyperParams(task.HyperParams)},
                                {Name: "CONFIG", Value: ts.serializeConfig(task.Config)},
                            },
                            Resources: corev1.ResourceRequirements{
                                Requests: corev1.ResourceList{
                                    corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%.1f", task.Resources.CPU)),
                                    corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", task.Resources.Memory)),
                                },
                                Limits: corev1.ResourceList{
                                    corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%.1f", task.Resources.CPU)),
                                    corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", task.Resources.Memory)),
                                },
                            },
                        },
                    },
                    RestartPolicy: corev1.RestartPolicyNever,
                },
            },
            BackoffLimit: int32Ptr(3),
        },
    }

    // 添加GPU资源
    if task.Resources.GPU > 0 {
        job.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"] = resource.MustParse(fmt.Sprintf("%d", task.Resources.GPU))
        job.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"] = resource.MustParse(fmt.Sprintf("%d", task.Resources.GPU))
    }

    _, err := ts.k8sClient.BatchV1().Jobs("ml-training").Create(context.Background(), job, metav1.CreateOptions{})
    return err
}

// monitorResources 监控资源
func (ts *TaskScheduler) monitorResources() {
    defer ts.wg.Done()

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ts.stopChan:
            return
        case <-ticker.C:
            ts.checkResourceUsage()
        }
    }
}

// syncTaskStatus 同步任务状态
func (ts *TaskScheduler) syncTaskStatus() {
    defer ts.wg.Done()

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ts.stopChan:
            return
        case <-ticker.C:
            ts.syncRunningTasks()
        }
    }
}

// 辅助方法
func (ts *TaskScheduler) validateTask(task *TrainingTask) error {
    if task.Name == "" {
        return fmt.Errorf("task name is required")
    }
    if task.ModelID == "" {
        return fmt.Errorf("model ID is required")
    }
    if task.DatasetID == "" {
        return fmt.Errorf("dataset ID is required")
    }
    return nil
}

func (ts *TaskScheduler) generateTaskID() string {
    return fmt.Sprintf("task-%s", uuid.New().String()[:8])
}

func (ts *TaskScheduler) updateTaskStatus(taskID string, status TaskStatus, message string) {
    updates := map[string]interface{}{
        "status": status,
    }

    if status == TaskStatusCompleted || status == TaskStatusFailed {
        now := time.Now()
        updates["completed_at"] = &now
    }

    if err := ts.db.Model(&TrainingTask{}).Where("id = ?", taskID).Updates(updates).Error; err != nil {
        ts.logger.Error("failed to update task status", zap.Error(err))
    }
}

func (ts *TaskScheduler) updateTask(task *TrainingTask) {
    if err := ts.db.Save(task).Error; err != nil {
        ts.logger.Error("failed to update task", zap.Error(err))
    }
}

func (ts *TaskScheduler) buildTrainingImage(task *TrainingTask) string {
    // 实现镜像构建逻辑
    return fmt.Sprintf("ml-training:%s", task.ID)
}

func (ts *TaskScheduler) serializeHyperParams(params map[string]interface{}) string {
    data, _ := json.Marshal(params)
    return string(data)
}

func (ts *TaskScheduler) serializeConfig(config TrainingConfig) string {
    data, _ := json.Marshal(config)
    return string(data)
}

func int32Ptr(i int32) *int32 {
    return &i
}
```

### 4.3 资源管理服务

```go
// ResourceManager 资源管理器
type ResourceManager struct {
    k8sClient   *kubernetes.Clientset
    prometheus  *prometheus.Client
    logger      *zap.Logger
    mu          sync.RWMutex
    allocations map[string]ResourceAllocation
}

// NewResourceManager 创建资源管理器
func NewResourceManager(k8sClient *kubernetes.Clientset, prometheus *prometheus.Client) *ResourceManager {
    return &ResourceManager{
        k8sClient:   k8sClient,
        prometheus:  prometheus,
        logger:      zap.L().Named("resource_manager"),
        allocations: make(map[string]ResourceAllocation),
    }
}

// GetClusterResources 获取集群资源
func (rm *ResourceManager) GetClusterResources() (*ClusterResources, error) {
    nodes, err := rm.k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list nodes: %w", err)
    }

    var totalCPU, totalMemory, totalGPU int64
    var availableCPU, availableMemory, availableGPU int64

    for _, node := range nodes.Items {
        // 总资源
        cpu := node.Status.Capacity.Cpu().MilliValue()
        memory := node.Status.Capacity.Memory().Value()
        gpu := node.Status.Capacity["nvidia.com/gpu"].Value()

        totalCPU += cpu
        totalMemory += memory
        totalGPU += gpu

        // 可用资源
        allocatableCPU := node.Status.Allocatable.Cpu().MilliValue()
        allocatableMemory := node.Status.Allocatable.Memory().Value()
        allocatableGPU := node.Status.Allocatable["nvidia.com/gpu"].Value()

        availableCPU += allocatableCPU
        availableMemory += allocatableMemory
        availableGPU += allocatableGPU
    }

    return &ClusterResources{
        Total: ResourceAllocation{
            CPU:    float64(totalCPU) / 1000,
            Memory: totalMemory / 1024 / 1024, // MB
            GPU:    int(totalGPU),
        },
        Available: ResourceAllocation{
            CPU:    float64(availableCPU) / 1000,
            Memory: availableMemory / 1024 / 1024, // MB
            GPU:    int(availableGPU),
        },
    }, nil
}

// HasAvailableResources 检查是否有可用资源
func (rm *ResourceManager) HasAvailableResources(required ResourceAllocation) bool {
    cluster, err := rm.GetClusterResources()
    if err != nil {
        rm.logger.Error("failed to get cluster resources", zap.Error(err))
        return false
    }

    // 计算已分配资源
    rm.mu.RLock()
    var allocatedCPU, allocatedMemory float64
    var allocatedGPU int
    for _, allocation := range rm.allocations {
        allocatedCPU += allocation.CPU
        allocatedMemory += float64(allocation.Memory)
        allocatedGPU += allocation.GPU
    }
    rm.mu.RUnlock()

    // 检查资源是否足够
    availableCPU := cluster.Available.CPU - allocatedCPU
    availableMemory := cluster.Available.Memory - int64(allocatedMemory)
    availableGPU := cluster.Available.GPU - allocatedGPU

    return availableCPU >= required.CPU &&
           availableMemory >= required.Memory &&
           availableGPU >= required.GPU
}

// AllocateResources 分配资源
func (rm *ResourceManager) AllocateResources(taskID string, resources ResourceAllocation) error {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    rm.allocations[taskID] = resources
    rm.logger.Info("resources allocated",
        zap.String("task_id", taskID),
        zap.Float64("cpu", resources.CPU),
        zap.Int64("memory", resources.Memory),
        zap.Int("gpu", resources.GPU))

    return nil
}

// ReleaseResources 释放资源
func (rm *ResourceManager) ReleaseResources(taskID string) {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    if allocation, exists := rm.allocations[taskID]; exists {
        delete(rm.allocations, taskID)
        rm.logger.Info("resources released",
            zap.String("task_id", taskID),
            zap.Float64("cpu", allocation.CPU),
            zap.Int64("memory", allocation.Memory),
            zap.Int("gpu", allocation.GPU))
    }
}

// GetResourceUsage 获取资源使用情况
func (rm *ResourceManager) GetResourceUsage() (*ResourceUsage, error) {
    cluster, err := rm.GetClusterResources()
    if err != nil {
        return nil, err
    }

    rm.mu.RLock()
    var allocatedCPU, allocatedMemory float64
    var allocatedGPU int
    for _, allocation := range rm.allocations {
        allocatedCPU += allocation.CPU
        allocatedMemory += float64(allocation.Memory)
        allocatedGPU += allocation.GPU
    }
    rm.mu.RUnlock()

    return &ResourceUsage{
        Total: cluster.Total,
        Allocated: ResourceAllocation{
            CPU:    allocatedCPU,
            Memory: int64(allocatedMemory),
            GPU:    allocatedGPU,
        },
        Available: ResourceAllocation{
            CPU:    cluster.Available.CPU - allocatedCPU,
            Memory: cluster.Available.Memory - int64(allocatedMemory),
            GPU:    cluster.Available.GPU - allocatedGPU,
        },
        Utilization: ResourceUtilization{
            CPU:    allocatedCPU / cluster.Total.CPU * 100,
            Memory: allocatedMemory / float64(cluster.Total.Memory) * 100,
            GPU:    float64(allocatedGPU) / float64(cluster.Total.GPU) * 100,
        },
    }, nil
}
```

### 4.4 实验管理服务

```go
// ExperimentManager 实验管理器
type ExperimentManager struct {
    db          *gorm.DB
    mlflow      *MLflowClient
    logger      *zap.Logger
    mu          sync.RWMutex
    experiments map[string]*Experiment
}

// NewExperimentManager 创建实验管理器
func NewExperimentManager(db *gorm.DB, mlflow *MLflowClient) *ExperimentManager {
    return &ExperimentManager{
        db:          db,
        mlflow:      mlflow,
        logger:      zap.L().Named("experiment_manager"),
        experiments: make(map[string]*Experiment),
    }
}

// CreateExperiment 创建实验
func (em *ExperimentManager) CreateExperiment(ctx context.Context, req *CreateExperimentRequest) (*Experiment, error) {
    // 验证请求
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // 创建实验
    experiment := &Experiment{
        ID:          em.generateExperimentID(),
        Name:        req.Name,
        Description: req.Description,
        Tags:        req.Tags,
        CreatedAt:   time.Now(),
        Status:      ExperimentStatusActive,
    }

    // 保存到数据库
    if err := em.db.Create(experiment).Error; err != nil {
        return nil, fmt.Errorf("failed to save experiment: %w", err)
    }

    // 创建MLflow实验
    if err := em.mlflow.CreateExperiment(experiment.ID, experiment.Name); err != nil {
        em.logger.Warn("failed to create MLflow experiment", zap.Error(err))
    }

    // 更新内存缓存
    em.mu.Lock()
    em.experiments[experiment.ID] = experiment
    em.mu.Unlock()

    em.logger.Info("experiment created successfully",
        zap.String("experiment_id", experiment.ID),
        zap.String("experiment_name", experiment.Name))

    return experiment, nil
}

// LogMetrics 记录指标
func (em *ExperimentManager) LogMetrics(ctx context.Context, experimentID string, runID string, metrics map[string]float64) error {
    // 记录到MLflow
    if err := em.mlflow.LogMetrics(runID, metrics); err != nil {
        return fmt.Errorf("failed to log metrics to MLflow: %w", err)
    }

    // 更新数据库
    for name, value := range metrics {
        metric := &Metric{
            ExperimentID: experimentID,
            RunID:        runID,
            Name:         name,
            Value:        value,
            Timestamp:    time.Now(),
        }

        if err := em.db.Create(metric).Error; err != nil {
            em.logger.Error("failed to save metric", zap.Error(err))
        }
    }

    return nil
}

// LogParameters 记录参数
func (em *ExperimentManager) LogParameters(ctx context.Context, experimentID string, runID string, params map[string]interface{}) error {
    // 记录到MLflow
    if err := em.mlflow.LogParameters(runID, params); err != nil {
        return fmt.Errorf("failed to log parameters to MLflow: %w", err)
    }

    // 更新数据库
    for name, value := range params {
        param := &Parameter{
            ExperimentID: experimentID,
            RunID:        runID,
            Name:         name,
            Value:        fmt.Sprintf("%v", value),
            Timestamp:    time.Now(),
        }

        if err := em.db.Create(param).Error; err != nil {
            em.logger.Error("failed to save parameter", zap.Error(err))
        }
    }

    return nil
}

// GetExperimentResults 获取实验结果
func (em *ExperimentManager) GetExperimentResults(ctx context.Context, experimentID string) (*ExperimentResults, error) {
    // 获取实验信息
    var experiment Experiment
    if err := em.db.Where("id = ?", experimentID).First(&experiment).Error; err != nil {
        return nil, fmt.Errorf("experiment not found: %w", err)
    }

    // 获取所有运行
    var runs []Run
    if err := em.db.Where("experiment_id = ?", experimentID).Find(&runs).Error; err != nil {
        return nil, fmt.Errorf("failed to get runs: %w", err)
    }

    // 获取最佳运行
    var bestRun *Run
    var bestMetric float64
    for i := range runs {
        if runs[i].Metrics != nil {
            if metric, exists := runs[i].Metrics["accuracy"]; exists {
                if bestRun == nil || metric > bestMetric {
                    bestRun = &runs[i]
                    bestMetric = metric
                }
            }
        }
    }

    return &ExperimentResults{
        Experiment: experiment,
        Runs:       runs,
        BestRun:    bestRun,
        Summary:    em.calculateSummary(runs),
    }, nil
}

// 辅助方法
func (em *ExperimentManager) generateExperimentID() string {
    return fmt.Sprintf("exp-%s", uuid.New().String()[:8])
}

func (em *ExperimentManager) calculateSummary(runs []Run) *ExperimentSummary {
    if len(runs) == 0 {
        return &ExperimentSummary{}
    }

    var totalDuration time.Duration
    var successCount int
    var metrics map[string][]float64 = make(map[string][]float64)

    for _, run := range runs {
        if run.Status == RunStatusCompleted {
            successCount++
            if run.EndTime != nil && run.StartTime != nil {
                totalDuration += run.EndTime.Sub(*run.StartTime)
            }
        }

        if run.Metrics != nil {
            for name, value := range run.Metrics {
                metrics[name] = append(metrics[name], value)
            }
        }
    }

    // 计算统计信息
    summary := &ExperimentSummary{
        TotalRuns:    len(runs),
        SuccessRuns:  successCount,
        FailureRuns:  len(runs) - successCount,
        AvgDuration:  totalDuration / time.Duration(successCount),
        SuccessRate:  float64(successCount) / float64(len(runs)) * 100,
    }

    // 计算指标统计
    for name, values := range metrics {
        if len(values) > 0 {
            summary.MetricStats[name] = &MetricStats{
                Min:     em.min(values),
                Max:     em.max(values),
                Mean:    em.mean(values),
                Median:  em.median(values),
                StdDev:  em.stdDev(values),
            }
        }
    }

    return summary
}

func (em *ExperimentManager) min(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}

func (em *ExperimentManager) max(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    max := values[0]
    for _, v := range values[1:] {
        if v > max {
            max = v
        }
    }
    return max
}

func (em *ExperimentManager) mean(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (em *ExperimentManager) median(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    sorted := make([]float64, len(values))
    copy(sorted, values)
    sort.Float64s(sorted)
    
    n := len(sorted)
    if n%2 == 0 {
        return (sorted[n/2-1] + sorted[n/2]) / 2
    }
    return sorted[n/2]
}

func (em *ExperimentManager) stdDev(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    mean := em.mean(values)
    sum := 0.0
    for _, v := range values {
        sum += math.Pow(v-mean, 2)
    }
    return math.Sqrt(sum / float64(len(values)))
}
```

## 5. 性能优化

### 5.1 分布式训练优化

```go
// DistributedTrainer 分布式训练器
type DistributedTrainer struct {
    strategy    TrainingStrategy
    communicator Communicator
    logger      *zap.Logger
}

// TrainingStrategy 训练策略接口
type TrainingStrategy interface {
    DistributeData(data []float64, numWorkers int) [][]float64
    AggregateGradients(gradients [][]float64) []float64
    UpdateModel(model *Model, gradients []float64) error
}

// DataParallelStrategy 数据并行策略
type DataParallelStrategy struct{}

// DistributeData 数据分发
func (dps *DataParallelStrategy) DistributeData(data []float64, numWorkers int) [][]float64 {
    chunks := make([][]float64, numWorkers)
    chunkSize := len(data) / numWorkers
    
    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if i == numWorkers-1 {
            end = len(data)
        }
        chunks[i] = data[start:end]
    }
    
    return chunks
}

// AggregateGradients 梯度聚合
func (dps *DataParallelStrategy) AggregateGradients(gradients [][]float64) []float64 {
    if len(gradients) == 0 {
        return nil
    }
    
    result := make([]float64, len(gradients[0]))
    for _, grad := range gradients {
        for i, v := range grad {
            result[i] += v
        }
    }
    
    // 平均梯度
    for i := range result {
        result[i] /= float64(len(gradients))
    }
    
    return result
}
```

### 5.2 内存优化

```go
// MemoryOptimizer 内存优化器
type MemoryOptimizer struct {
    pool        *sync.Pool
    cache       *lru.Cache
    logger      *zap.Logger
}

// NewMemoryOptimizer 创建内存优化器
func NewMemoryOptimizer() *MemoryOptimizer {
    cache, _ := lru.New(1000)
    
    return &MemoryOptimizer{
        pool: &sync.Pool{
            New: func() interface{} {
                return make([]float64, 0, 1024)
            },
        },
        cache: cache,
        logger: zap.L().Named("memory_optimizer"),
    }
}

// GetBuffer 获取缓冲区
func (mo *MemoryOptimizer) GetBuffer() []float64 {
    return mo.pool.Get().([]float64)
}

// PutBuffer 归还缓冲区
func (mo *MemoryOptimizer) PutBuffer(buf []float64) {
    buf = buf[:0] // 清空但保留容量
    mo.pool.Put(buf)
}

// CacheResult 缓存结果
func (mo *MemoryOptimizer) CacheResult(key string, result interface{}) {
    mo.cache.Add(key, result)
}

// GetCachedResult 获取缓存结果
func (mo *MemoryOptimizer) GetCachedResult(key string) (interface{}, bool) {
    return mo.cache.Get(key)
}
```

## 6. 分布式训练

### 6.1 通信协议

```go
// Communicator 通信器接口
type Communicator interface {
    Send(rank int, data []byte) error
    Receive(rank int) ([]byte, error)
    Broadcast(rank int, data []byte) error
    AllReduce(data []float64, op ReduceOp) ([]float64, error)
}

// TCPCommunicator TCP通信器
type TCPCommunicator struct {
    rank        int
    worldSize   int
    connections map[int]net.Conn
    listener    net.Listener
    logger      *zap.Logger
}

// NewTCPCommunicator 创建TCP通信器
func NewTCPCommunicator(rank int, worldSize int, port int) (*TCPCommunicator, error) {
    tc := &TCPCommunicator{
        rank:        rank,
        worldSize:   worldSize,
        connections: make(map[int]net.Conn),
        logger:      zap.L().Named("tcp_communicator"),
    }
    
    // 启动监听
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return nil, err
    }
    tc.listener = listener
    
    // 启动接受连接协程
    go tc.acceptConnections()
    
    return tc, nil
}

// AllReduce 全归约操作
func (tc *TCPCommunicator) AllReduce(data []float64, op ReduceOp) ([]float64, error) {
    result := make([]float64, len(data))
    copy(result, data)
    
    // 执行全归约
    for i := 0; i < tc.worldSize; i++ {
        if i == tc.rank {
            continue
        }
        
        // 接收其他节点的数据
        received, err := tc.Receive(i)
        if err != nil {
            return nil, err
        }
        
        var otherData []float64
        if err := json.Unmarshal(received, &otherData); err != nil {
            return nil, err
        }
        
        // 执行归约操作
        for j := range result {
            switch op {
            case ReduceOpSum:
                result[j] += otherData[j]
            case ReduceOpMax:
                if otherData[j] > result[j] {
                    result[j] = otherData[j]
                }
            case ReduceOpMin:
                if otherData[j] < result[j] {
                    result[j] = otherData[j]
                }
            }
        }
    }
    
    return result, nil
}

// acceptConnections 接受连接
func (tc *TCPCommunicator) acceptConnections() {
    for {
        conn, err := tc.listener.Accept()
        if err != nil {
            tc.logger.Error("failed to accept connection", zap.Error(err))
            continue
        }
        
        // 处理连接
        go tc.handleConnection(conn)
    }
}

// handleConnection 处理连接
func (tc *TCPCommunicator) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // 读取消息
    buffer := make([]byte, 4096)
    n, err := conn.Read(buffer)
    if err != nil {
        tc.logger.Error("failed to read from connection", zap.Error(err))
        return
    }
    
    // 处理消息
    tc.logger.Info("received message", zap.Int("size", n))
}
```

## 总结

本文档提供了基于Go语言的AI/ML模型训练平台完整实现方案，包括：

1. **形式化定义**：使用数学符号严格定义模型训练平台的概念
2. **数学建模**：提供梯度下降算法和资源调度优化的数学证明
3. **架构设计**：清晰的系统架构图和组件职责划分
4. **Go语言实现**：完整的任务调度、资源管理、实验管理服务实现
5. **性能优化**：分布式训练和内存优化策略
6. **分布式训练**：通信协议和全归约操作实现

该实现方案具有高可扩展性、高可靠性和高性能，适用于大规模机器学习模型训练场景。
