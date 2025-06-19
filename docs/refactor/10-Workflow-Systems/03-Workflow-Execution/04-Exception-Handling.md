# 04-工作流异常处理 (Workflow Exception Handling)

## 概述

工作流异常处理是确保工作流系统可靠性和健壮性的关键技术，通过形式化的异常模型、异常检测算法和恢复策略，实现对工作流执行过程中各种异常情况的处理和管理。

## 目录

1. [异常处理基础理论](#1-异常处理基础理论)
2. [异常分类与建模](#2-异常分类与建模)
3. [异常检测算法](#3-异常检测算法)
4. [异常恢复策略](#4-异常恢复策略)
5. [Go语言实现](#5-go语言实现)
6. [应用案例](#6-应用案例)

## 1. 异常处理基础理论

### 1.1 异常定义

#### 1.1.1 异常基本概念

**定义 1.1.1** (工作流异常)
工作流异常 $E$ 是一个三元组 $(t, \tau, \sigma)$，其中：

- $t$: 异常发生的时间点
- $\tau$: 异常类型
- $\sigma$: 异常状态信息

**定义 1.1.2** (异常空间)
异常空间 $\mathcal{E}$ 是所有可能异常的集合：

```latex
$$\mathcal{E} = \mathcal{T} \times \mathcal{T} \times \mathcal{S}$$
```

其中：

- $\mathcal{T}$: 异常类型集合
- $\mathcal{S}$: 异常状态集合

#### 1.1.2 异常传播模型

**定义 1.1.3** (异常传播)
异常传播函数 $P: \mathcal{E} \times \mathcal{W} \to \mathcal{E}^*$ 定义异常在工作流中的传播：

```latex
$$P(E, W) = \{E_1, E_2, ..., E_n\}$$
```

其中：

- $W$: 工作流实例
- $E_i$: 传播的异常

**定理 1.1.1** (异常传播单调性)

对于异常 $E_1 \subseteq E_2$，其传播满足单调性：

```latex
$$P(E_1, W) \subseteq P(E_2, W)$$
```

### 1.2 异常处理语义

#### 1.2.1 异常处理操作

**定义 1.2.1** (异常处理操作)

异常处理操作 $H$ 是一个四元组 $(D, R, C, F)$，其中：

- $D$: 异常检测函数
- $R$: 异常恢复函数
- $C$: 异常补偿函数
- $F$: 异常失败处理函数

**定义 1.2.2** (异常处理语义)

异常处理语义 $\llbracket H \rrbracket$ 定义为：

```latex
$$\llbracket H \rrbracket = \lambda E. \begin{cases}
R(E) & \text{if } D(E) = \text{recoverable} \\
C(E) & \text{if } D(E) = \text{compensatable} \\
F(E) & \text{if } D(E) = \text{fatal}
\end{cases}$$
```

## 2. 异常分类与建模

### 2.1 异常分类体系

#### 2.1.1 按异常来源分类

**定义 2.1.1** (系统异常)
系统异常 $E_{sys}$ 是由系统资源或环境问题引起的异常：

```latex
$$E_{sys} \in \{E_{cpu}, E_{memory}, E_{network}, E_{disk}\}$$
```

**定义 2.1.2** (业务异常)
业务异常 $E_{biz}$ 是由业务逻辑或数据问题引起的异常：

```latex
$$E_{biz} \in \{E_{validation}, E_{authorization}, E_{data}, E_{logic}\}$$
```

**定义 2.1.3** (外部异常)
外部异常 $E_{ext}$ 是由外部服务或系统引起的异常：

```latex
$$E_{ext} \in \{E_{timeout}, E_{unavailable}, E_{invalid}, E_{rate\_limit}\}$$
```

#### 2.1.2 按异常严重程度分类

**定义 2.1.4** (异常严重程度)
异常严重程度 $S(E)$ 定义为：

```latex
$$S(E) = \begin{cases}
\text{low} & \text{if } E \text{ can be ignored} \\
\text{medium} & \text{if } E \text{ needs attention} \\
\text{high} & \text{if } E \text{ requires immediate action} \\
\text{critical} & \text{if } E \text{ causes system failure}
\end{cases}$$
```

### 2.2 异常建模

#### 2.2.1 异常状态机

**定义 2.2.1** (异常状态机)
异常状态机 $ASM = (Q, \Sigma, \delta, q_0, F)$ 其中：

- $Q$: 状态集合 $\{normal, detected, recovering, compensated, failed\}$
- $\Sigma$: 输入字母表（异常事件）
- $\delta: Q \times \Sigma \to Q$: 状态转移函数
- $q_0$: 初始状态（normal）
- $F$: 接受状态集合

**定理 2.2.1** (异常状态可达性)
对于异常状态机，从正常状态到任何异常状态都是可达的。

**证明**:
通过构造性证明，对于任意异常状态 $q \in Q$，存在输入序列 $\sigma_1, \sigma_2, ..., \sigma_n$ 使得：

```latex
$$\delta^*(q_0, \sigma_1\sigma_2...\sigma_n) = q$$
```

#### 2.2.2 异常概率模型

**定义 2.2.2** (异常概率)
异常概率 $P(E|W)$ 定义为在工作流 $W$ 执行过程中发生异常 $E$ 的概率。

**定义 2.2.3** (异常分布)
异常分布 $D_E$ 定义为异常发生时间的概率分布：

```latex
$$D_E(t) = P(T_E \leq t)$$
```

其中 $T_E$ 是异常发生时间。

## 3. 异常检测算法

### 3.1 基于规则的异常检测

#### 3.1.1 规则引擎

**算法 3.1.1** (规则引擎异常检测)

```go
// 规则引擎异常检测器
type RuleEngine struct {
    rules []DetectionRule
}

type DetectionRule struct {
    ID          string
    Condition   func(*WorkflowState) bool
    Action      func(*WorkflowState) *Exception
    Priority    int
    Enabled     bool
}

func (re *RuleEngine) DetectExceptions(state *WorkflowState) []*Exception {
    var exceptions []*Exception

    // 按优先级排序规则
    sortedRules := re.sortRulesByPriority()

    for _, rule := range sortedRules {
        if !rule.Enabled {
            continue
        }

        if rule.Condition(state) {
            exception := rule.Action(state)
            if exception != nil {
                exceptions = append(exceptions, exception)
            }
        }
    }

    return exceptions
}

func (re *RuleEngine) sortRulesByPriority() []DetectionRule {
    sorted := make([]DetectionRule, len(re.rules))
    copy(sorted, re.rules)

    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Priority > sorted[j].Priority
    })

    return sorted
}
```

#### 3.1.2 预定义规则

```go
// 预定义异常检测规则
func createPredefinedRules() []DetectionRule {
    return []DetectionRule{
        // 超时规则
        {
            ID: "timeout_rule",
            Condition: func(state *WorkflowState) bool {
                return time.Since(state.StartTime) > state.Timeout
            },
            Action: func(state *WorkflowState) *Exception {
                return &Exception{
                    Type:    TimeoutException,
                    Message: "Workflow execution timeout",
                    Time:    time.Now(),
                    State:   state,
                }
            },
            Priority: 1,
            Enabled:  true,
        },

        // 资源耗尽规则
        {
            ID: "resource_exhaustion_rule",
            Condition: func(state *WorkflowState) bool {
                return state.CPUUsage > 90.0 || state.MemoryUsage > 90.0
            },
            Action: func(state *WorkflowState) *Exception {
                return &Exception{
                    Type:    ResourceException,
                    Message: "Resource exhaustion detected",
                    Time:    time.Now(),
                    State:   state,
                }
            },
            Priority: 2,
            Enabled:  true,
        },

        // 数据一致性规则
        {
            ID: "data_consistency_rule",
            Condition: func(state *WorkflowState) bool {
                return !state.DataConsistency
            },
            Action: func(state *WorkflowState) *Exception {
                return &Exception{
                    Type:    DataException,
                    Message: "Data consistency violation",
                    Time:    time.Now(),
                    State:   state,
                }
            },
            Priority: 3,
            Enabled:  true,
        },
    }
}
```

### 3.2 基于统计的异常检测

#### 3.2.1 统计异常检测

**算法 3.2.1** (统计异常检测)

```go
// 统计异常检测器
type StatisticalDetector struct {
    baseline    *BaselineModel
    threshold   float64
    windowSize  int
    history     []float64
}

type BaselineModel struct {
    Mean     float64
    StdDev   float64
    Min      float64
    Max      float64
}

func (sd *StatisticalDetector) DetectAnomaly(value float64) bool {
    // 更新历史数据
    sd.history = append(sd.history, value)
    if len(sd.history) > sd.windowSize {
        sd.history = sd.history[1:]
    }

    // 计算当前统计量
    currentMean := sd.calculateMean(sd.history)
    currentStdDev := sd.calculateStdDev(sd.history, currentMean)

    // 计算z-score
    zScore := math.Abs((value - currentMean) / currentStdDev)

    // 检查是否超过阈值
    return zScore > sd.threshold
}

func (sd *StatisticalDetector) calculateMean(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }

    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (sd *StatisticalDetector) calculateStdDev(values []float64, mean float64) float64 {
    if len(values) == 0 {
        return 0
    }

    sum := 0.0
    for _, v := range values {
        sum += math.Pow(v-mean, 2)
    }
    return math.Sqrt(sum / float64(len(values)))
}
```

#### 3.2.2 机器学习异常检测

**算法 3.2.2** (机器学习异常检测)

```go
// 机器学习异常检测器
type MLDetector struct {
    model       AnomalyDetectionModel
    features    []string
    scaler      *StandardScaler
}

type AnomalyDetectionModel interface {
    Predict(features []float64) float64
    Train(data [][]float64) error
}

// 隔离森林异常检测
type IsolationForest struct {
    trees       []*IsolationTree
    nTrees      int
    maxDepth    int
}

type IsolationTree struct {
    root        *Node
    maxDepth    int
}

type Node struct {
    Feature     int
    Threshold   float64
    Left        *Node
    Right       *Node
    IsLeaf      bool
    Depth       int
}

func (iforest *IsolationForest) Train(data [][]float64) error {
    iforest.trees = make([]*IsolationTree, iforest.nTrees)

    for i := 0; i < iforest.nTrees; i++ {
        tree := &IsolationTree{
            maxDepth: iforest.maxDepth,
        }

        // 随机采样数据
        sampleSize := int(math.Min(float64(len(data)), 256.0))
        sample := iforest.randomSample(data, sampleSize)

        // 构建隔离树
        tree.root = tree.buildTree(sample, 0)
        iforest.trees[i] = tree
    }

    return nil
}

func (iforest *IsolationForest) Predict(features []float64) float64 {
    scores := make([]float64, iforest.nTrees)

    for i, tree := range iforest.trees {
        scores[i] = tree.pathLength(features)
    }

    // 计算异常分数
    avgPathLength := iforest.calculateMean(scores)
    return math.Pow(2, -avgPathLength/iforest.c(256))
}

func (iforest *IsolationForest) c(n int) float64 {
    if n > 2 {
        return 2*(math.Log(float64(n-1))+0.5772156649) - 2*float64(n-1)/float64(n)
    }
    if n == 2 {
        return 1
    }
    return 0
}

func (iforest *IsolationForest) calculateMean(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}
```

## 4. 异常恢复策略

### 4.1 重试策略

#### 4.1.1 指数退避重试

**算法 4.1.1** (指数退避重试)

```go
// 指数退避重试策略
type ExponentialBackoffRetry struct {
    maxRetries  int
    baseDelay   time.Duration
    maxDelay    time.Duration
    multiplier  float64
}

func (ebr *ExponentialBackoffRetry) Execute(operation func() error) error {
    var lastErr error

    for attempt := 0; attempt <= ebr.maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }

        lastErr = err

        if attempt == ebr.maxRetries {
            break
        }

        // 计算延迟时间
        delay := ebr.calculateDelay(attempt)
        time.Sleep(delay)
    }

    return fmt.Errorf("operation failed after %d retries: %v", ebr.maxRetries, lastErr)
}

func (ebr *ExponentialBackoffRetry) calculateDelay(attempt int) time.Duration {
    delay := float64(ebr.baseDelay) * math.Pow(ebr.multiplier, float64(attempt))

    // 添加随机抖动
    jitter := delay * 0.1 * (rand.Float64() - 0.5)
    delay += jitter

    // 限制最大延迟
    if delay > float64(ebr.maxDelay) {
        delay = float64(ebr.maxDelay)
    }

    return time.Duration(delay)
}
```

#### 4.1.2 自适应重试

**算法 4.1.2** (自适应重试)

```go
// 自适应重试策略
type AdaptiveRetry struct {
    maxRetries      int
    successRate     float64
    failureRate     float64
    currentDelay    time.Duration
    minDelay        time.Duration
    maxDelay        time.Duration
}

func (ar *AdaptiveRetry) Execute(operation func() error) error {
    var lastErr error
    consecutiveFailures := 0
    consecutiveSuccesses := 0

    for attempt := 0; attempt <= ar.maxRetries; attempt++ {
        err := operation()
        if err == nil {
            consecutiveSuccesses++
            consecutiveFailures = 0

            // 成功时减少延迟
            ar.currentDelay = time.Duration(float64(ar.currentDelay) * 0.9)
            if ar.currentDelay < ar.minDelay {
                ar.currentDelay = ar.minDelay
            }

            return nil
        }

        lastErr = err
        consecutiveFailures++
        consecutiveSuccesses = 0

        if attempt == ar.maxRetries {
            break
        }

        // 失败时增加延迟
        ar.currentDelay = time.Duration(float64(ar.currentDelay) * 1.1)
        if ar.currentDelay > ar.maxDelay {
            ar.currentDelay = ar.maxDelay
        }

        time.Sleep(ar.currentDelay)
    }

    return fmt.Errorf("operation failed after %d retries: %v", ar.maxRetries, lastErr)
}
```

### 4.2 补偿策略

#### 4.2.1 Saga模式补偿

**算法 4.2.1** (Saga补偿)

```go
// Saga补偿策略
type SagaCompensation struct {
    steps       []SagaStep
    compensations []CompensationStep
}

type SagaStep struct {
    ID          string
    Execute     func() error
    Compensate  func() error
    Completed   bool
}

type CompensationStep struct {
    StepID      string
    Compensate  func() error
    Executed    bool
}

func (saga *SagaCompensation) Execute() error {
    for i, step := range saga.steps {
        err := step.Execute()
        if err != nil {
            // 执行失败，开始补偿
            return saga.compensate(i)
        }

        step.Completed = true
        saga.compensations = append(saga.compensations, CompensationStep{
            StepID:     step.ID,
            Compensate: step.Compensate,
            Executed:   false,
        })
    }

    return nil
}

func (saga *SagaCompensation) compensate(failedIndex int) error {
    // 从失败步骤开始，反向执行补偿
    for i := failedIndex - 1; i >= 0; i-- {
        if saga.steps[i].Completed {
            err := saga.steps[i].Compensate()
            if err != nil {
                return fmt.Errorf("compensation failed for step %s: %v", saga.steps[i].ID, err)
            }
        }
    }

    return fmt.Errorf("saga execution failed at step %d", failedIndex)
}
```

#### 4.2.2 事务补偿

**算法 4.2.2** (事务补偿)

```go
// 事务补偿策略
type TransactionCompensation struct {
    operations  []CompensableOperation
    log         *CompensationLog
}

type CompensableOperation struct {
    ID          string
    Execute     func() error
    Compensate  func() error
    Completed   bool
}

type CompensationLog struct {
    operations  []string
    mutex       sync.Mutex
}

func (tc *TransactionCompensation) Execute() error {
    for i, op := range tc.operations {
        err := op.Execute()
        if err != nil {
            // 执行失败，开始补偿
            return tc.compensate(i)
        }

        op.Completed = true
        tc.log.AddOperation(op.ID)
    }

    return nil
}

func (tc *TransactionCompensation) compensate(failedIndex int) error {
    // 从失败操作开始，反向执行补偿
    for i := failedIndex - 1; i >= 0; i-- {
        if tc.operations[i].Completed {
            err := tc.operations[i].Compensate()
            if err != nil {
                return fmt.Errorf("compensation failed for operation %s: %v", tc.operations[i].ID, err)
            }

            tc.log.RemoveOperation(tc.operations[i].ID)
        }
    }

    return fmt.Errorf("transaction failed at operation %d", failedIndex)
}

func (cl *CompensationLog) AddOperation(opID string) {
    cl.mutex.Lock()
    defer cl.mutex.Unlock()
    cl.operations = append(cl.operations, opID)
}

func (cl *CompensationLog) RemoveOperation(opID string) {
    cl.mutex.Lock()
    defer cl.mutex.Unlock()

    for i, id := range cl.operations {
        if id == opID {
            cl.operations = append(cl.operations[:i], cl.operations[i+1:]...)
            break
        }
    }
}
```

## 5. Go语言实现

### 5.1 异常处理框架

#### 5.1.1 异常处理器

```go
// 异常处理器
type ExceptionHandler struct {
    detectors   []ExceptionDetector
    strategies  map[ExceptionType]RecoveryStrategy
    logger      *Logger
}

type ExceptionDetector interface {
    Detect(state *WorkflowState) []*Exception
}

type RecoveryStrategy interface {
    Recover(exception *Exception) error
    CanHandle(exception *Exception) bool
}

func NewExceptionHandler() *ExceptionHandler {
    return &ExceptionHandler{
        detectors:  make([]ExceptionDetector, 0),
        strategies: make(map[ExceptionType]RecoveryStrategy),
        logger:     NewLogger(),
    }
}

func (eh *ExceptionHandler) RegisterDetector(detector ExceptionDetector) {
    eh.detectors = append(eh.detectors, detector)
}

func (eh *ExceptionHandler) RegisterStrategy(exceptionType ExceptionType, strategy RecoveryStrategy) {
    eh.strategies[exceptionType] = strategy
}

func (eh *ExceptionHandler) Handle(state *WorkflowState) error {
    // 检测异常
    var allExceptions []*Exception
    for _, detector := range eh.detectors {
        exceptions := detector.Detect(state)
        allExceptions = append(allExceptions, exceptions...)
    }

    if len(allExceptions) == 0 {
        return nil
    }

    // 处理异常
    for _, exception := range allExceptions {
        strategy, exists := eh.strategies[exception.Type]
        if !exists {
            eh.logger.Error("No recovery strategy found for exception type: %v", exception.Type)
            continue
        }

        if !strategy.CanHandle(exception) {
            eh.logger.Error("Strategy cannot handle exception: %v", exception)
            continue
        }

        err := strategy.Recover(exception)
        if err != nil {
            eh.logger.Error("Recovery failed for exception: %v, error: %v", exception, err)
            return err
        }

        eh.logger.Info("Successfully recovered from exception: %v", exception)
    }

    return nil
}
```

#### 5.1.2 异常监控器

```go
// 异常监控器
type ExceptionMonitor struct {
    handlers    map[string]*ExceptionHandler
    metrics     *ExceptionMetrics
    alerting    *AlertingSystem
}

type ExceptionMetrics struct {
    totalExceptions    int64
    handledExceptions  int64
    recoveryTime       time.Duration
    mutex              sync.RWMutex
}

type AlertingSystem struct {
    rules      []AlertRule
    channels   []AlertChannel
}

type AlertRule struct {
    Condition  func(*ExceptionMetrics) bool
    Message    string
    Severity   AlertSeverity
}

type AlertChannel interface {
    Send(alert *Alert) error
}

type Alert struct {
    Message   string
    Severity  AlertSeverity
    Time      time.Time
    Data      map[string]interface{}
}

func NewExceptionMonitor() *ExceptionMonitor {
    return &ExceptionMonitor{
        handlers: make(map[string]*ExceptionHandler),
        metrics:  &ExceptionMetrics{},
        alerting: &AlertingSystem{
            rules:    make([]AlertRule, 0),
            channels: make([]AlertChannel, 0),
        },
    }
}

func (em *ExceptionMonitor) Monitor(workflowID string, handler *ExceptionHandler) {
    em.handlers[workflowID] = handler

    // 启动监控协程
    go em.monitorWorkflow(workflowID, handler)
}

func (em *ExceptionMonitor) monitorWorkflow(workflowID string, handler *ExceptionHandler) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        // 检查工作流状态
        state := em.getWorkflowState(workflowID)
        if state == nil {
            continue
        }

        // 处理异常
        start := time.Now()
        err := handler.Handle(state)
        recoveryTime := time.Since(start)

        // 更新指标
        em.updateMetrics(err, recoveryTime)

        // 检查告警规则
        em.checkAlertRules()
    }
}

func (em *ExceptionMonitor) updateMetrics(err error, recoveryTime time.Duration) {
    em.metrics.mutex.Lock()
    defer em.metrics.mutex.Unlock()

    em.metrics.totalExceptions++
    if err == nil {
        em.metrics.handledExceptions++
    }
    em.metrics.recoveryTime = recoveryTime
}

func (em *ExceptionMonitor) checkAlertRules() {
    em.metrics.mutex.RLock()
    metrics := *em.metrics
    em.metrics.mutex.RUnlock()

    for _, rule := range em.alerting.rules {
        if rule.Condition(&metrics) {
            alert := &Alert{
                Message:  rule.Message,
                Severity: rule.Severity,
                Time:     time.Now(),
                Data:     map[string]interface{}{"metrics": metrics},
            }

            em.sendAlert(alert)
        }
    }
}

func (em *ExceptionMonitor) sendAlert(alert *Alert) {
    for _, channel := range em.alerting.channels {
        go func(ch AlertChannel) {
            if err := ch.Send(alert); err != nil {
                log.Printf("Failed to send alert: %v", err)
            }
        }(channel)
    }
}
```

### 5.2 异常恢复框架

#### 5.2.1 恢复策略管理器

```go
// 恢复策略管理器
type RecoveryStrategyManager struct {
    strategies  map[ExceptionType][]RecoveryStrategy
    selector    StrategySelector
}

type StrategySelector interface {
    Select(exception *Exception, strategies []RecoveryStrategy) RecoveryStrategy
}

// 基于优先级的策略选择器
type PriorityBasedSelector struct{}

func (pbs *PriorityBasedSelector) Select(exception *Exception, strategies []RecoveryStrategy) RecoveryStrategy {
    var bestStrategy RecoveryStrategy
    bestPriority := -1

    for _, strategy := range strategies {
        if !strategy.CanHandle(exception) {
            continue
        }

        priority := pbs.getPriority(strategy)
        if priority > bestPriority {
            bestPriority = priority
            bestStrategy = strategy
        }
    }

    return bestStrategy
}

func (pbs *PriorityBasedSelector) getPriority(strategy RecoveryStrategy) int {
    // 基于策略类型返回优先级
    switch strategy.(type) {
    case *RetryStrategy:
        return 1
    case *CompensationStrategy:
        return 2
    case *FallbackStrategy:
        return 3
    default:
        return 0
    }
}

// 重试策略
type RetryStrategy struct {
    maxRetries  int
    backoff     BackoffStrategy
}

func (rs *RetryStrategy) CanHandle(exception *Exception) bool {
    return exception.Type == TimeoutException || exception.Type == NetworkException
}

func (rs *RetryStrategy) Recover(exception *Exception) error {
    for attempt := 0; attempt <= rs.maxRetries; attempt++ {
        err := rs.retryOperation(exception)
        if err == nil {
            return nil
        }

        if attempt == rs.maxRetries {
            return fmt.Errorf("retry failed after %d attempts: %v", rs.maxRetries, err)
        }

        delay := rs.backoff.CalculateDelay(attempt)
        time.Sleep(delay)
    }

    return nil
}

// 补偿策略
type CompensationStrategy struct {
    compensations map[string]func() error
}

func (cs *CompensationStrategy) CanHandle(exception *Exception) bool {
    return exception.Type == DataException || exception.Type == BusinessException
}

func (cs *CompensationStrategy) Recover(exception *Exception) error {
    // 执行补偿操作
    for stepID, compensate := range cs.compensations {
        if err := compensate(); err != nil {
            return fmt.Errorf("compensation failed for step %s: %v", stepID, err)
        }
    }

    return nil
}

// 降级策略
type FallbackStrategy struct {
    fallback func() error
}

func (fs *FallbackStrategy) CanHandle(exception *Exception) bool {
    return true // 降级策略可以处理所有异常
}

func (fs *FallbackStrategy) Recover(exception *Exception) error {
    return fs.fallback()
}
```

## 6. 应用案例

### 6.1 电商订单处理异常处理

#### 6.1.1 订单处理异常处理案例

```go
// 电商订单处理异常处理
func EcommerceOrderExceptionHandling() {
    // 创建异常处理器
    handler := NewExceptionHandler()

    // 注册检测器
    ruleEngine := &RuleEngine{
        rules: createPredefinedRules(),
    }
    handler.RegisterDetector(ruleEngine)

    // 注册恢复策略
    retryStrategy := &RetryStrategy{
        maxRetries: 3,
        backoff:     &ExponentialBackoffRetry{},
    }
    handler.RegisterStrategy(TimeoutException, retryStrategy)

    compensationStrategy := &CompensationStrategy{
        compensations: map[string]func() error{
            "payment": func() error {
                // 回滚支付
                return nil
            },
            "inventory": func() error {
                // 恢复库存
                return nil
            },
        },
    }
    handler.RegisterStrategy(DataException, compensationStrategy)

    // 创建监控器
    monitor := NewExceptionMonitor()
    monitor.Monitor("order_processing", handler)

    // 模拟订单处理
    order := createOrder()
    err := processOrder(order)
    if err != nil {
        log.Printf("Order processing failed: %v", err)
    }
}

func createOrder() *Order {
    return &Order{
        ID:       "order_123",
        Customer: "customer_456",
        Items:    []OrderItem{{ProductID: "prod_789", Quantity: 2}},
        Total:    100.0,
    }
}

func processOrder(order *Order) error {
    // 模拟订单处理流程
    steps := []func() error{
        func() error { return validateOrder(order) },
        func() error { return checkInventory(order) },
        func() error { return processPayment(order) },
        func() error { return updateInventory(order) },
        func() error { return sendNotification(order) },
    }

    for _, step := range steps {
        if err := step(); err != nil {
            return err
        }
    }

    return nil
}
```

### 6.2 微服务异常处理

#### 6.2.1 服务间异常处理

```go
// 微服务异常处理
func MicroserviceExceptionHandling() {
    // 创建熔断器
    circuitBreaker := NewCircuitBreaker(5, 10*time.Second, 0.5)

    // 创建重试策略
    retryStrategy := &AdaptiveRetry{
        maxRetries:  3,
        minDelay:    100 * time.Millisecond,
        maxDelay:    5 * time.Second,
        currentDelay: 1 * time.Second,
    }

    // 创建降级策略
    fallbackStrategy := &FallbackStrategy{
        fallback: func() error {
            // 返回缓存数据或默认值
            return nil
        },
    }

    // 服务调用包装器
    serviceCaller := &ServiceCaller{
        circuitBreaker: circuitBreaker,
        retryStrategy:  retryStrategy,
        fallback:       fallbackStrategy,
    }

    // 模拟服务调用
    result, err := serviceCaller.Call("user-service", func() (interface{}, error) {
        return callUserService()
    })

    if err != nil {
        log.Printf("Service call failed: %v", err)
    } else {
        log.Printf("Service call succeeded: %v", result)
    }
}

type ServiceCaller struct {
    circuitBreaker *CircuitBreaker
    retryStrategy  *AdaptiveRetry
    fallback       *FallbackStrategy
}

func (sc *ServiceCaller) Call(serviceName string, operation func() (interface{}, error)) (interface{}, error) {
    // 检查熔断器状态
    if !sc.circuitBreaker.Ready() {
        return sc.fallback.Recover(nil)
    }

    // 执行操作
    var result interface{}
    err := sc.retryStrategy.Execute(func() error {
        var callErr error
        result, callErr = operation()
        return callErr
    })

    // 更新熔断器状态
    sc.circuitBreaker.RecordResult(err == nil)

    if err != nil {
        return sc.fallback.Recover(&Exception{Type: ServiceException, Message: err.Error()})
    }

    return result, nil
}

func callUserService() (interface{}, error) {
    // 模拟服务调用
    if rand.Float64() < 0.1 { // 10% 失败率
        return nil, fmt.Errorf("service unavailable")
    }

    return map[string]interface{}{
        "user_id": "123",
        "name":    "John Doe",
        "email":   "john@example.com",
    }, nil
}
```

## 总结

工作流异常处理是一个复杂的系统工程，需要结合形式化理论、算法设计和工程实践。通过建立完整的异常检测、恢复和监控体系，可以显著提升工作流系统的可靠性和健壮性。

关键要点：

1. **异常建模**: 使用状态机、概率模型等数学工具建立异常模型
2. **检测算法**: 开发基于规则、统计和机器学习的异常检测算法
3. **恢复策略**: 实施重试、补偿、降级等恢复策略
4. **工程实践**: 使用Go语言实现高性能的异常处理系统

通过持续的异常监控、分析和处理，可以确保工作流系统在各种异常情况下都能保持稳定可靠的运行。
