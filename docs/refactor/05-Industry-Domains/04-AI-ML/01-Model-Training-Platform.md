# 01-模型训练平台 (Model Training Platform)

## 目录

1. [理论基础](#1-理论基础)
2. [形式化定义](#2-形式化定义)
3. [架构设计](#3-架构设计)
4. [Go语言实现](#4-go语言实现)
5. [分布式训练](#5-分布式训练)

## 1. 理论基础

### 1.1 模型训练平台定义

模型训练平台是机器学习系统的核心组件，负责数据预处理、模型训练、超参数调优和模型评估。

**形式化定义**：

```math
模型训练平台定义为六元组：
MTP = (D, M, T, H, E, V)

其中：
- D: 数据集集合，D = \{d_1, d_2, ..., d_n\}
- M: 模型集合，M = \{m_1, m_2, ..., m_k\}
- T: 训练函数，T: M \times D \times H \rightarrow M'
- H: 超参数空间，H = \{h_1, h_2, ..., h_m\}
- E: 评估函数，E: M \times D \rightarrow \mathbb{R}
- V: 验证函数，V: M \times D \rightarrow \mathbb{B}
```

### 1.2 核心功能

1. **数据管理**: 数据集版本控制和预处理
2. **模型训练**: 分布式训练和超参数调优
3. **实验管理**: 实验跟踪和结果比较
4. **模型评估**: 性能评估和模型选择
5. **模型部署**: 模型打包和部署

## 2. 形式化定义

### 2.1 训练过程

```math
训练过程定义为：
Train(model, data, hyperparams) = 
  \text{for } epoch \in \{1, 2, ..., epochs\} \text{ do}
    \text{for } batch \in data \text{ do}
      loss = Forward(model, batch)
      gradients = Backward(loss)
      model = Update(model, gradients, hyperparams)
    \text{end for}
    \text{if } epoch \% validation\_interval == 0 \text{ then}
      validation\_score = Evaluate(model, validation\_data)
    \text{end if}
  \text{end for}
  \text{return } model
```

### 2.2 损失函数

```math
损失函数定义为：
L(\theta, x, y) = \frac{1}{N} \sum_{i=1}^{N} l(f_\theta(x_i), y_i)

其中：
- \theta: 模型参数
- x: 输入数据
- y: 真实标签
- f_\theta: 模型函数
- l: 单个样本损失函数
```

## 3. 架构设计

### 3.1 分层架构

```go
// ModelTrainingPlatform 模型训练平台
type ModelTrainingPlatform struct {
    // 数据层 - 数据管理和预处理
    DataLayer *DataLayer
    // 模型层 - 模型定义和管理
    ModelLayer *ModelLayer
    // 训练层 - 训练执行和调度
    TrainingLayer *TrainingLayer
    // 实验层 - 实验管理和跟踪
    ExperimentLayer *ExperimentLayer
    // 部署层 - 模型部署和服务
    DeploymentLayer *DeploymentLayer
}

// DataLayer 数据层
type DataLayer struct {
    DatasetManager    *DatasetManager
    DataPreprocessor  *DataPreprocessor
    FeatureEngineer   *FeatureEngineer
    DataValidator     *DataValidator
}

// ModelLayer 模型层
type ModelLayer struct {
    ModelRegistry     *ModelRegistry
    ModelBuilder      *ModelBuilder
    ModelValidator    *ModelValidator
    ModelVersioning   *ModelVersioning
}

// TrainingLayer 训练层
type TrainingLayer struct {
    TrainingScheduler *TrainingScheduler
    HyperparameterOptimizer *HyperparameterOptimizer
    DistributedTrainer *DistributedTrainer
    TrainingMonitor   *TrainingMonitor
}

// ExperimentLayer 实验层
type ExperimentLayer struct {
    ExperimentManager *ExperimentManager
    MetricsTracker    *MetricsTracker
    ArtifactManager   *ArtifactManager
    ExperimentComparator *ExperimentComparator
}

// DeploymentLayer 部署层
type DeploymentLayer struct {
    ModelPackager     *ModelPackager
    ModelDeployer     *ModelDeployer
    ModelServing      *ModelServing
    ModelMonitoring   *ModelMonitoring
}
```

## 4. Go语言实现

### 4.1 数据集管理

```go
// Dataset 数据集定义
type Dataset struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Schema      *DataSchema       `json:"schema"`
    Statistics  *DataStatistics   `json:"statistics"`
    Storage     *StorageInfo      `json:"storage"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

// DataSchema 数据模式
type DataSchema struct {
    Features    []Feature         `json:"features"`
    Target      *Feature          `json:"target"`
    Metadata    map[string]string `json:"metadata"`
}

// Feature 特征定义
type Feature struct {
    Name        string            `json:"name"`
    Type        FeatureType       `json:"type"`
    Description string            `json:"description"`
    Constraints []Constraint      `json:"constraints"`
}

// FeatureType 特征类型
type FeatureType int

const (
    FeatureTypeNumeric FeatureType = iota
    FeatureTypeCategorical
    FeatureTypeText
    FeatureTypeImage
    FeatureTypeTimeSeries
)

// DataStatistics 数据统计
type DataStatistics struct {
    RowCount    int64             `json:"row_count"`
    ColumnCount int               `json:"column_count"`
    MissingValues map[string]int64 `json:"missing_values"`
    Distributions map[string]interface{} `json:"distributions"`
}

// DatasetManager 数据集管理器
type DatasetManager struct {
    db          *sql.DB
    storage     *StorageManager
    validator   *DatasetValidator
    processor   *DataPreprocessor
    mu          sync.RWMutex
}

// NewDatasetManager 创建数据集管理器
func NewDatasetManager(db *sql.DB, storage *StorageManager) *DatasetManager {
    return &DatasetManager{
        db:        db,
        storage:   storage,
        validator: NewDatasetValidator(),
        processor: NewDataPreprocessor(),
    }
}

// CreateDataset 创建数据集
func (dm *DatasetManager) CreateDataset(dataset *Dataset, data []byte) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    // 验证数据集
    if err := dm.validator.ValidateDataset(dataset); err != nil {
        return err
    }
    
    // 生成数据集ID
    dataset.ID = dm.generateDatasetID()
    dataset.CreatedAt = time.Now()
    dataset.UpdatedAt = time.Now()
    
    // 计算数据统计
    if err := dm.calculateStatistics(dataset, data); err != nil {
        return err
    }
    
    // 保存到存储
    if err := dm.storage.SaveDataset(dataset.ID, data); err != nil {
        return err
    }
    
    // 保存到数据库
    if err := dm.saveDataset(dataset); err != nil {
        return err
    }
    
    return nil
}

// GetDataset 获取数据集
func (dm *DatasetManager) GetDataset(datasetID string) (*Dataset, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    // 从数据库获取数据集信息
    dataset, err := dm.loadDataset(datasetID)
    if err != nil {
        return nil, err
    }
    
    return dataset, nil
}

// calculateStatistics 计算数据统计
func (dm *DatasetManager) calculateStatistics(dataset *Dataset, data []byte) error {
    // 解析数据
    records, err := dm.parseData(data, dataset.Schema)
    if err != nil {
        return err
    }
    
    // 计算统计信息
    stats := &DataStatistics{
        RowCount:    int64(len(records)),
        ColumnCount: len(dataset.Schema.Features),
        MissingValues: make(map[string]int64),
        Distributions: make(map[string]interface{}),
    }
    
    // 计算缺失值
    for _, feature := range dataset.Schema.Features {
        missingCount := int64(0)
        for _, record := range records {
            if record[feature.Name] == nil {
                missingCount++
            }
        }
        stats.MissingValues[feature.Name] = missingCount
    }
    
    // 计算分布
    for _, feature := range dataset.Schema.Features {
        distribution := dm.calculateDistribution(records, feature)
        stats.Distributions[feature.Name] = distribution
    }
    
    dataset.Statistics = stats
    return nil
}

// saveDataset 保存数据集到数据库
func (dm *DatasetManager) saveDataset(dataset *Dataset) error {
    query := `
        INSERT INTO datasets (id, name, version, description, schema, statistics, storage, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        name = VALUES(name),
        version = VALUES(version),
        description = VALUES(description),
        schema = VALUES(schema),
        statistics = VALUES(statistics),
        storage = VALUES(storage),
        updated_at = VALUES(updated_at)
    `
    
    schemaJSON, err := json.Marshal(dataset.Schema)
    if err != nil {
        return err
    }
    
    statisticsJSON, err := json.Marshal(dataset.Statistics)
    if err != nil {
        return err
    }
    
    storageJSON, err := json.Marshal(dataset.Storage)
    if err != nil {
        return err
    }
    
    _, err = dm.db.Exec(query,
        dataset.ID,
        dataset.Name,
        dataset.Version,
        dataset.Description,
        schemaJSON,
        statisticsJSON,
        storageJSON,
        dataset.CreatedAt,
        dataset.UpdatedAt,
    )
    
    return err
}

// loadDataset 从数据库加载数据集
func (dm *DatasetManager) loadDataset(datasetID string) (*Dataset, error) {
    query := `
        SELECT id, name, version, description, schema, statistics, storage, created_at, updated_at
        FROM datasets WHERE id = ?
    `
    
    var dataset Dataset
    var schemaJSON, statisticsJSON, storageJSON []byte
    
    err := dm.db.QueryRow(query, datasetID).Scan(
        &dataset.ID,
        &dataset.Name,
        &dataset.Version,
        &dataset.Description,
        &schemaJSON,
        &statisticsJSON,
        &storageJSON,
        &dataset.CreatedAt,
        &dataset.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    // 解析模式
    if err := json.Unmarshal(schemaJSON, &dataset.Schema); err != nil {
        return nil, err
    }
    
    // 解析统计信息
    if err := json.Unmarshal(statisticsJSON, &dataset.Statistics); err != nil {
        return nil, err
    }
    
    // 解析存储信息
    if err := json.Unmarshal(storageJSON, &dataset.Storage); err != nil {
        return nil, err
    }
    
    return &dataset, nil
}
```

### 4.2 模型管理

```go
// Model 模型定义
type Model struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Type        ModelType         `json:"type"`
    Architecture *ModelArchitecture `json:"architecture"`
    Hyperparameters map[string]interface{} `json:"hyperparameters"`
    Weights     *WeightsInfo      `json:"weights"`
    Metrics     *ModelMetrics     `json:"metrics"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

// ModelType 模型类型
type ModelType int

const (
    ModelTypeLinear ModelType = iota
    ModelTypeLogistic
    ModelTypeRandomForest
    ModelTypeGradientBoosting
    ModelTypeNeuralNetwork
    ModelTypeCNN
    ModelTypeRNN
    ModelTypeTransformer
)

// ModelArchitecture 模型架构
type ModelArchitecture struct {
    Layers      []Layer           `json:"layers"`
    InputShape  []int             `json:"input_shape"`
    OutputShape []int             `json:"output_shape"`
    Parameters  int64             `json:"parameters"`
}

// Layer 层定义
type Layer struct {
    Type        string            `json:"type"`
    Units       int               `json:"units"`
    Activation  string            `json:"activation"`
    Parameters  map[string]interface{} `json:"parameters"`
}

// WeightsInfo 权重信息
type WeightsInfo struct {
    FilePath    string            `json:"file_path"`
    Size        int64             `json:"size"`
    Checksum    string            `json:"checksum"`
    Format      string            `json:"format"`
}

// ModelMetrics 模型指标
type ModelMetrics struct {
    Accuracy    float64           `json:"accuracy"`
    Precision   float64           `json:"precision"`
    Recall      float64           `json:"recall"`
    F1Score     float64           `json:"f1_score"`
    Loss        float64           `json:"loss"`
    AUC         float64           `json:"auc"`
}

// ModelRegistry 模型注册表
type ModelRegistry struct {
    db          *sql.DB
    storage     *StorageManager
    validator   *ModelValidator
    builder     *ModelBuilder
    mu          sync.RWMutex
}

// NewModelRegistry 创建模型注册表
func NewModelRegistry(db *sql.DB, storage *StorageManager) *ModelRegistry {
    return &ModelRegistry{
        db:        db,
        storage:   storage,
        validator: NewModelValidator(),
        builder:   NewModelBuilder(),
    }
}

// RegisterModel 注册模型
func (mr *ModelRegistry) RegisterModel(model *Model) error {
    mr.mu.Lock()
    defer mr.mu.Unlock()
    
    // 验证模型
    if err := mr.validator.ValidateModel(model); err != nil {
        return err
    }
    
    // 生成模型ID
    model.ID = mr.generateModelID()
    model.CreatedAt = time.Now()
    model.UpdatedAt = time.Now()
    
    // 保存到数据库
    if err := mr.saveModel(model); err != nil {
        return err
    }
    
    return nil
}

// GetModel 获取模型
func (mr *ModelRegistry) GetModel(modelID string) (*Model, error) {
    mr.mu.RLock()
    defer mr.mu.RUnlock()
    
    // 从数据库获取模型
    model, err := mr.loadModel(modelID)
    if err != nil {
        return nil, err
    }
    
    return model, nil
}

// saveModel 保存模型到数据库
func (mr *ModelRegistry) saveModel(model *Model) error {
    query := `
        INSERT INTO models (id, name, version, type, architecture, hyperparameters, weights, metrics, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        name = VALUES(name),
        version = VALUES(version),
        type = VALUES(type),
        architecture = VALUES(architecture),
        hyperparameters = VALUES(hyperparameters),
        weights = VALUES(weights),
        metrics = VALUES(metrics),
        updated_at = VALUES(updated_at)
    `
    
    architectureJSON, err := json.Marshal(model.Architecture)
    if err != nil {
        return err
    }
    
    hyperparametersJSON, err := json.Marshal(model.Hyperparameters)
    if err != nil {
        return err
    }
    
    weightsJSON, err := json.Marshal(model.Weights)
    if err != nil {
        return err
    }
    
    metricsJSON, err := json.Marshal(model.Metrics)
    if err != nil {
        return err
    }
    
    _, err = mr.db.Exec(query,
        model.ID,
        model.Name,
        model.Version,
        model.Type,
        architectureJSON,
        hyperparametersJSON,
        weightsJSON,
        metricsJSON,
        model.CreatedAt,
        model.UpdatedAt,
    )
    
    return err
}

// loadModel 从数据库加载模型
func (mr *ModelRegistry) loadModel(modelID string) (*Model, error) {
    query := `
        SELECT id, name, version, type, architecture, hyperparameters, weights, metrics, created_at, updated_at
        FROM models WHERE id = ?
    `
    
    var model Model
    var architectureJSON, hyperparametersJSON, weightsJSON, metricsJSON []byte
    
    err := mr.db.QueryRow(query, modelID).Scan(
        &model.ID,
        &model.Name,
        &model.Version,
        &model.Type,
        &architectureJSON,
        &hyperparametersJSON,
        &weightsJSON,
        &metricsJSON,
        &model.CreatedAt,
        &model.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    // 解析架构
    if err := json.Unmarshal(architectureJSON, &model.Architecture); err != nil {
        return nil, err
    }
    
    // 解析超参数
    if err := json.Unmarshal(hyperparametersJSON, &model.Hyperparameters); err != nil {
        return nil, err
    }
    
    // 解析权重信息
    if err := json.Unmarshal(weightsJSON, &model.Weights); err != nil {
        return nil, err
    }
    
    // 解析指标
    if err := json.Unmarshal(metricsJSON, &model.Metrics); err != nil {
        return nil, err
    }
    
    return &model, nil
}
```

### 4.3 训练调度器

```go
// TrainingJob 训练任务
type TrainingJob struct {
    ID              string            `json:"id"`
    ModelID         string            `json:"model_id"`
    DatasetID       string            `json:"dataset_id"`
    Hyperparameters map[string]interface{} `json:"hyperparameters"`
    Status          JobStatus         `json:"status"`
    Progress        float64           `json:"progress"`
    Metrics         *TrainingMetrics  `json:"metrics"`
    CreatedAt       time.Time         `json:"created_at"`
    StartedAt       *time.Time        `json:"started_at"`
    CompletedAt     *time.Time        `json:"completed_at"`
}

// JobStatus 任务状态
type JobStatus int

const (
    JobStatusPending JobStatus = iota
    JobStatusRunning
    JobStatusCompleted
    JobStatusFailed
    JobStatusCancelled
)

// TrainingMetrics 训练指标
type TrainingMetrics struct {
    TrainLoss    []float64         `json:"train_loss"`
    ValLoss      []float64         `json:"val_loss"`
    TrainAccuracy []float64        `json:"train_accuracy"`
    ValAccuracy  []float64         `json:"val_accuracy"`
    LearningRate []float64         `json:"learning_rate"`
}

// TrainingScheduler 训练调度器
type TrainingScheduler struct {
    jobs         map[string]*TrainingJob
    workers      []*TrainingWorker
    queue        chan *TrainingJob
    modelRegistry *ModelRegistry
    datasetManager *DatasetManager
    mu           sync.RWMutex
}

// NewTrainingScheduler 创建训练调度器
func NewTrainingScheduler(workerCount int, modelRegistry *ModelRegistry, datasetManager *DatasetManager) *TrainingScheduler {
    scheduler := &TrainingScheduler{
        jobs:           make(map[string]*TrainingJob),
        workers:        make([]*TrainingWorker, workerCount),
        queue:          make(chan *TrainingJob, 100),
        modelRegistry:  modelRegistry,
        datasetManager: datasetManager,
    }
    
    // 启动工作线程
    for i := 0; i < workerCount; i++ {
        scheduler.workers[i] = NewTrainingWorker(scheduler.queue)
        go scheduler.workers[i].Start()
    }
    
    return scheduler
}

// SubmitJob 提交训练任务
func (ts *TrainingScheduler) SubmitJob(job *TrainingJob) error {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    
    // 生成任务ID
    job.ID = ts.generateJobID()
    job.Status = JobStatusPending
    job.CreatedAt = time.Now()
    
    // 保存任务
    ts.jobs[job.ID] = job
    
    // 提交到队列
    select {
    case ts.queue <- job:
        return nil
    default:
        return errors.New("training queue is full")
    }
}

// GetJob 获取任务
func (ts *TrainingScheduler) GetJob(jobID string) (*TrainingJob, error) {
    ts.mu.RLock()
    defer ts.mu.RUnlock()
    
    job, exists := ts.jobs[jobID]
    if !exists {
        return nil, errors.New("job not found")
    }
    
    return job, nil
}

// CancelJob 取消任务
func (ts *TrainingScheduler) CancelJob(jobID string) error {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    
    job, exists := ts.jobs[jobID]
    if !exists {
        return errors.New("job not found")
    }
    
    if job.Status == JobStatusRunning {
        // 通知工作线程取消任务
        for _, worker := range ts.workers {
            worker.CancelJob(jobID)
        }
    }
    
    job.Status = JobStatusCancelled
    return nil
}

// TrainingWorker 训练工作线程
type TrainingWorker struct {
    queue        chan *TrainingJob
    cancelChan   chan string
    running      bool
}

// NewTrainingWorker 创建训练工作线程
func NewTrainingWorker(queue chan *TrainingJob) *TrainingWorker {
    return &TrainingWorker{
        queue:      queue,
        cancelChan: make(chan string, 10),
        running:    false,
    }
}

// Start 启动工作线程
func (tw *TrainingWorker) Start() {
    tw.running = true
    
    for tw.running {
        select {
        case job := <-tw.queue:
            tw.executeJob(job)
        case jobID := <-tw.cancelChan:
            tw.cancelJob(jobID)
        }
    }
}

// executeJob 执行训练任务
func (tw *TrainingWorker) executeJob(job *TrainingJob) {
    // 更新任务状态
    job.Status = JobStatusRunning
    now := time.Now()
    job.StartedAt = &now
    
    // 执行训练
    err := tw.trainModel(job)
    
    if err != nil {
        job.Status = JobStatusFailed
    } else {
        job.Status = JobStatusCompleted
        job.Progress = 1.0
    }
    
    completedAt := time.Now()
    job.CompletedAt = &completedAt
}

// trainModel 训练模型
func (tw *TrainingWorker) trainModel(job *TrainingJob) error {
    // 获取模型
    model, err := tw.modelRegistry.GetModel(job.ModelID)
    if err != nil {
        return err
    }
    
    // 获取数据集
    dataset, err := tw.datasetManager.GetDataset(job.DatasetID)
    if err != nil {
        return err
    }
    
    // 创建训练器
    trainer := NewModelTrainer(model, dataset, job.Hyperparameters)
    
    // 开始训练
    return trainer.Train(func(epoch int, metrics *TrainingMetrics) {
        // 更新进度
        job.Progress = float64(epoch) / float64(trainer.GetTotalEpochs())
        job.Metrics = metrics
    })
}

// CancelJob 取消任务
func (tw *TrainingWorker) CancelJob(jobID string) {
    select {
    case tw.cancelChan <- jobID:
    default:
        // 取消通道已满
    }
}
```

## 5. 分布式训练

### 5.1 分布式训练器

```go
// DistributedTrainer 分布式训练器
type DistributedTrainer struct {
    nodes        []*TrainingNode
    coordinator  *TrainingCoordinator
    modelRegistry *ModelRegistry
    datasetManager *DatasetManager
}

// NewDistributedTrainer 创建分布式训练器
func NewDistributedTrainer(nodeCount int) *DistributedTrainer {
    trainer := &DistributedTrainer{
        nodes: make([]*TrainingNode, nodeCount),
    }
    
    // 创建训练节点
    for i := 0; i < nodeCount; i++ {
        trainer.nodes[i] = NewTrainingNode(i)
    }
    
    // 创建协调器
    trainer.coordinator = NewTrainingCoordinator(trainer.nodes)
    
    return trainer
}

// Train 分布式训练
func (dt *DistributedTrainer) Train(job *TrainingJob) error {
    // 分发数据到各个节点
    if err := dt.distributeData(job); err != nil {
        return err
    }
    
    // 分发模型到各个节点
    if err := dt.distributeModel(job); err != nil {
        return err
    }
    
    // 开始分布式训练
    return dt.coordinator.StartTraining(job)
}

// distributeData 分发数据
func (dt *DistributedTrainer) distributeData(job *TrainingJob) error {
    // 获取数据集
    dataset, err := dt.datasetManager.GetDataset(job.DatasetID)
    if err != nil {
        return err
    }
    
    // 分割数据
    dataShards := dt.splitData(dataset, len(dt.nodes))
    
    // 分发到各个节点
    for i, node := range dt.nodes {
        if err := node.SetData(dataShards[i]); err != nil {
            return err
        }
    }
    
    return nil
}

// distributeModel 分发模型
func (dt *DistributedTrainer) distributeModel(job *TrainingJob) error {
    // 获取模型
    model, err := dt.modelRegistry.GetModel(job.ModelID)
    if err != nil {
        return err
    }
    
    // 分发到各个节点
    for _, node := range dt.nodes {
        if err := node.SetModel(model); err != nil {
            return err
        }
    }
    
    return nil
}

// splitData 分割数据
func (dt *DistributedTrainer) splitData(dataset *Dataset, nodeCount int) [][]Record {
    // 实现数据分割逻辑
    // 这里简化处理，实际应该考虑数据分布和负载均衡
    return nil
}
```

## 总结

模型训练平台是机器学习系统的核心，提供完整的数据管理、模型训练、实验跟踪和模型部署功能。本文档提供了完整的理论基础、形式化定义和Go语言实现。

### 关键要点

1. **数据管理**: 数据集版本控制和预处理
2. **模型管理**: 模型注册和版本控制
3. **训练调度**: 分布式训练和资源管理
4. **实验跟踪**: 实验管理和指标跟踪
5. **模型部署**: 模型打包和服务化

### 扩展阅读

- [推理服务](./02-Inference-Service.md)
- [数据处理管道](./03-Data-Processing-Pipeline.md)
- [特征工程](./04-Feature-Engineering.md)
