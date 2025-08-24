# 工作流系统 (Workflow Systems)

## 1. 基本概念

### 1.1 工作流系统定义

**工作流系统 (Workflow System)** 是一种用于自动化业务流程的软件系统，它能够定义、执行和监控一系列相互关联的任务，确保业务流程按照预定义的规则和顺序进行。

### 1.2 核心组件

- **工作流定义**: 描述业务流程的结构和规则
- **工作流引擎**: 执行和管理工作流的核心组件
- **任务执行器**: 执行具体任务的组件
- **状态管理器**: 管理工作流和任务的状态
- **监控系统**: 监控工作流的执行情况

### 1.3 工作流类型

- **顺序工作流**: 任务按顺序执行
- **并行工作流**: 多个任务同时执行
- **条件工作流**: 根据条件选择不同的执行路径
- **循环工作流**: 包含循环和迭代的工作流

## 2. 工作流引擎

### 2.1 基本工作流引擎

```go
// 工作流引擎
type WorkflowEngine struct {
    workflows map[string]*Workflow
    tasks     map[string]*Task
    executor  TaskExecutor
    mu        sync.RWMutex
}

type Workflow struct {
    ID          string
    Name        string
    Version     string
    Steps       []*Step
    State       WorkflowState
    Context     map[string]interface{}
    CreatedAt   time.Time
    StartedAt   time.Time
    CompletedAt time.Time
    mu          sync.Mutex
}

type Step struct {
    ID          string
    Name        string
    Type        StepType
    TaskID      string
    Condition   string
    NextSteps   []string
    State       StepState
    Input       map[string]interface{}
    Output      map[string]interface{}
    StartedAt   time.Time
    CompletedAt time.Time
}

type WorkflowState int

const (
    Created WorkflowState = iota
    Running
    Completed
    Failed
    Cancelled
)

type StepState int

const (
    Pending StepState = iota
    Running
    Completed
    Failed
    Skipped
)

type StepType int

const (
    Task StepType = iota
    Condition
    Parallel
    Loop
)

func (we *WorkflowEngine) CreateWorkflow(definition *WorkflowDefinition) (*Workflow, error) {
    we.mu.Lock()
    defer we.mu.Unlock()
    
    workflow := &Workflow{
        ID:        generateWorkflowID(),
        Name:      definition.Name,
        Version:   definition.Version,
        Steps:     make([]*Step, 0),
        State:     Created,
        Context:   make(map[string]interface{}),
        CreatedAt: time.Now(),
    }
    
    // 构建步骤
    for _, stepDef := range definition.Steps {
        step := &Step{
            ID:        stepDef.ID,
            Name:      stepDef.Name,
            Type:      stepDef.Type,
            TaskID:    stepDef.TaskID,
            Condition: stepDef.Condition,
            NextSteps: stepDef.NextSteps,
            State:     Pending,
            Input:     make(map[string]interface{}),
            Output:    make(map[string]interface{}),
        }
        workflow.Steps = append(workflow.Steps, step)
    }
    
    we.workflows[workflow.ID] = workflow
    return workflow, nil
}

func (we *WorkflowEngine) StartWorkflow(workflowID string, input map[string]interface{}) error {
    we.mu.Lock()
    defer we.mu.Unlock()
    
    workflow := we.workflows[workflowID]
    if workflow == nil {
        return fmt.Errorf("workflow %s not found", workflowID)
    }
    
    if workflow.State != Created {
        return fmt.Errorf("workflow %s is not in created state", workflowID)
    }
    
    // 设置输入上下文
    workflow.Context = input
    workflow.State = Running
    workflow.StartedAt = time.Now()
    
    // 启动工作流执行
    go we.executeWorkflow(workflow)
    
    return nil
}

func (we *WorkflowEngine) executeWorkflow(workflow *Workflow) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Workflow %s execution panicked: %v", workflow.ID, r)
            workflow.State = Failed
        }
    }()
    
    // 找到起始步骤
    startSteps := we.findStartSteps(workflow)
    
    // 执行起始步骤
    for _, step := range startSteps {
        we.executeStep(workflow, step)
    }
    
    // 检查工作流是否完成
    if we.isWorkflowCompleted(workflow) {
        workflow.State = Completed
        workflow.CompletedAt = time.Now()
    }
}

func (we *WorkflowEngine) executeStep(workflow *Workflow, step *Step) {
    step.mu.Lock()
    step.State = Running
    step.StartedAt = time.Now()
    step.mu.Unlock()
    
    defer func() {
        step.mu.Lock()
        step.CompletedAt = time.Now()
        step.mu.Unlock()
    }()
    
    switch step.Type {
    case Task:
        we.executeTaskStep(workflow, step)
    case Condition:
        we.executeConditionStep(workflow, step)
    case Parallel:
        we.executeParallelStep(workflow, step)
    case Loop:
        we.executeLoopStep(workflow, step)
    }
}

func (we *WorkflowEngine) executeTaskStep(workflow *Workflow, step *Step) {
    // 准备任务输入
    input := we.prepareTaskInput(workflow, step)
    
    // 执行任务
    output, err := we.executor.ExecuteTask(step.TaskID, input)
    if err != nil {
        step.State = Failed
        log.Printf("Task %s execution failed: %v", step.TaskID, err)
        return
    }
    
    // 保存任务输出
    step.Output = output
    step.State = Completed
    
    // 执行后续步骤
    we.executeNextSteps(workflow, step)
}

func (we *WorkflowEngine) executeConditionStep(workflow *Workflow, step *Step) {
    // 评估条件
    result, err := we.evaluateCondition(step.Condition, workflow.Context)
    if err != nil {
        step.State = Failed
        log.Printf("Condition evaluation failed: %v", err)
        return
    }
    
    step.Output["result"] = result
    step.State = Completed
    
    // 根据条件结果选择后续步骤
    if result.(bool) {
        // 条件为真，执行第一个分支
        if len(step.NextSteps) > 0 {
            nextStep := we.findStep(workflow, step.NextSteps[0])
            if nextStep != nil {
                we.executeStep(workflow, nextStep)
            }
        }
    } else {
        // 条件为假，执行第二个分支
        if len(step.NextSteps) > 1 {
            nextStep := we.findStep(workflow, step.NextSteps[1])
            if nextStep != nil {
                we.executeStep(workflow, nextStep)
            }
        }
    }
}

func (we *WorkflowEngine) executeParallelStep(workflow *Workflow, step *Step) {
    // 并行执行所有后续步骤
    var wg sync.WaitGroup
    for _, nextStepID := range step.NextSteps {
        nextStep := we.findStep(workflow, nextStepID)
        if nextStep != nil {
            wg.Add(1)
            go func(s *Step) {
                defer wg.Done()
                we.executeStep(workflow, s)
            }(nextStep)
        }
    }
    
    wg.Wait()
    step.State = Completed
}

func (we *WorkflowEngine) executeLoopStep(workflow *Workflow, step *Step) {
    // 获取循环条件
    condition := step.Condition
    maxIterations := 1000 // 防止无限循环
    
    for i := 0; i < maxIterations; i++ {
        // 评估循环条件
        shouldContinue, err := we.evaluateCondition(condition, workflow.Context)
        if err != nil {
            step.State = Failed
            log.Printf("Loop condition evaluation failed: %v", err)
            return
        }
        
        if !shouldContinue.(bool) {
            break
        }
        
        // 执行循环体
        for _, nextStepID := range step.NextSteps {
            nextStep := we.findStep(workflow, nextStepID)
            if nextStep != nil {
                we.executeStep(workflow, nextStep)
            }
        }
    }
    
    step.State = Completed
}

func (we *WorkflowEngine) findStartSteps(workflow *Workflow) []*Step {
    var startSteps []*Step
    
    // 找到没有前置步骤的步骤
    for _, step := range workflow.Steps {
        if !we.hasPredecessors(workflow, step) {
            startSteps = append(startSteps, step)
        }
    }
    
    return startSteps
}

func (we *WorkflowEngine) hasPredecessors(workflow *Workflow, step *Step) bool {
    for _, s := range workflow.Steps {
        for _, nextStepID := range s.NextSteps {
            if nextStepID == step.ID {
                return true
            }
        }
    }
    return false
}

func (we *WorkflowEngine) isWorkflowCompleted(workflow *Workflow) bool {
    for _, step := range workflow.Steps {
        if step.State != Completed && step.State != Skipped {
            return false
        }
    }
    return true
}

func (we *WorkflowEngine) findStep(workflow *Workflow, stepID string) *Step {
    for _, step := range workflow.Steps {
        if step.ID == stepID {
            return step
        }
    }
    return nil
}
```

### 2.2 任务执行器

```go
// 任务执行器
type TaskExecutor struct {
    tasks map[string]TaskHandler
    mu    sync.RWMutex
}

type TaskHandler func(input map[string]interface{}) (map[string]interface{}, error)

func (te *TaskExecutor) RegisterTask(taskID string, handler TaskHandler) {
    te.mu.Lock()
    defer te.mu.Unlock()
    
    te.tasks[taskID] = handler
}

func (te *TaskExecutor) ExecuteTask(taskID string, input map[string]interface{}) (map[string]interface{}, error) {
    te.mu.RLock()
    handler, exists := te.tasks[taskID]
    te.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("task %s not found", taskID)
    }
    
    return handler(input)
}

// 示例任务处理器
func (te *TaskExecutor) RegisterDefaultTasks() {
    // 数据处理任务
    te.RegisterTask("process_data", func(input map[string]interface{}) (map[string]interface{}, error) {
        data, ok := input["data"].(string)
        if !ok {
            return nil, fmt.Errorf("data field not found or not string")
        }
        
        // 处理数据
        processedData := strings.ToUpper(data)
        
        return map[string]interface{}{
            "processed_data": processedData,
            "length":         len(processedData),
        }, nil
    })
    
    // 验证任务
    te.RegisterTask("validate_data", func(input map[string]interface{}) (map[string]interface{}, error) {
        data, ok := input["data"].(string)
        if !ok {
            return nil, fmt.Errorf("data field not found or not string")
        }
        
        // 验证数据
        isValid := len(data) > 0 && len(data) < 1000
        
        return map[string]interface{}{
            "is_valid": isValid,
            "message":  "Data validation completed",
        }, nil
    })
    
    // 通知任务
    te.RegisterTask("send_notification", func(input map[string]interface{}) (map[string]interface{}, error) {
        message, ok := input["message"].(string)
        if !ok {
            return nil, fmt.Errorf("message field not found or not string")
        }
        
        // 发送通知（模拟）
        log.Printf("Sending notification: %s", message)
        
        return map[string]interface{}{
            "sent":     true,
            "timestamp": time.Now(),
        }, nil
    })
}
```

## 3. 状态管理

### 3.1 工作流状态管理器

```go
// 工作流状态管理器
type WorkflowStateManager struct {
    workflows map[string]*Workflow
    history   map[string][]*StateChange
    mu        sync.RWMutex
}

type StateChange struct {
    WorkflowID string
    StepID     string
    OldState   interface{}
    NewState   interface{}
    Timestamp  time.Time
    Reason     string
}

func (wsm *WorkflowStateManager) UpdateWorkflowState(workflowID string, newState WorkflowState, reason string) error {
    wsm.mu.Lock()
    defer wsm.mu.Unlock()
    
    workflow := wsm.workflows[workflowID]
    if workflow == nil {
        return fmt.Errorf("workflow %s not found", workflowID)
    }
    
    oldState := workflow.State
    workflow.State = newState
    
    // 记录状态变化
    change := &StateChange{
        WorkflowID: workflowID,
        OldState:   oldState,
        NewState:   newState,
        Timestamp:  time.Now(),
        Reason:     reason,
    }
    
    wsm.history[workflowID] = append(wsm.history[workflowID], change)
    
    return nil
}

func (wsm *WorkflowStateManager) UpdateStepState(workflowID, stepID string, newState StepState, reason string) error {
    wsm.mu.Lock()
    defer wsm.mu.Unlock()
    
    workflow := wsm.workflows[workflowID]
    if workflow == nil {
        return fmt.Errorf("workflow %s not found", workflowID)
    }
    
    step := wsm.findStep(workflow, stepID)
    if step == nil {
        return fmt.Errorf("step %s not found in workflow %s", stepID, workflowID)
    }
    
    oldState := step.State
    step.State = newState
    
    // 记录状态变化
    change := &StateChange{
        WorkflowID: workflowID,
        StepID:     stepID,
        OldState:   oldState,
        NewState:   newState,
        Timestamp:  time.Now(),
        Reason:     reason,
    }
    
    wsm.history[workflowID] = append(wsm.history[workflowID], change)
    
    return nil
}

func (wsm *WorkflowStateManager) GetWorkflowHistory(workflowID string) ([]*StateChange, error) {
    wsm.mu.RLock()
    defer wsm.mu.RUnlock()
    
    history, exists := wsm.history[workflowID]
    if !exists {
        return nil, fmt.Errorf("no history found for workflow %s", workflowID)
    }
    
    return history, nil
}

func (wsm *WorkflowStateManager) GetWorkflowState(workflowID string) (WorkflowState, error) {
    wsm.mu.RLock()
    defer wsm.mu.RUnlock()
    
    workflow := wsm.workflows[workflowID]
    if workflow == nil {
        return Created, fmt.Errorf("workflow %s not found", workflowID)
    }
    
    return workflow.State, nil
}

func (wsm *WorkflowStateManager) findStep(workflow *Workflow, stepID string) *Step {
    for _, step := range workflow.Steps {
        if step.ID == stepID {
            return step
        }
    }
    return nil
}
```

### 3.2 持久化存储

```go
// 工作流持久化存储
type WorkflowStorage struct {
    db *sql.DB
}

func (ws *WorkflowStorage) SaveWorkflow(workflow *Workflow) error {
    tx, err := ws.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 保存工作流基本信息
    _, err = tx.Exec(`
        INSERT INTO workflows (id, name, version, state, context, created_at, started_at, completed_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        state = VALUES(state),
        context = VALUES(context),
        started_at = VALUES(started_at),
        completed_at = VALUES(completed_at)
    `, workflow.ID, workflow.Name, workflow.Version, workflow.State, 
       marshalContext(workflow.Context), workflow.CreatedAt, workflow.StartedAt, workflow.CompletedAt)
    
    if err != nil {
        return err
    }
    
    // 保存步骤信息
    for _, step := range workflow.Steps {
        _, err = tx.Exec(`
            INSERT INTO workflow_steps (workflow_id, step_id, name, type, task_id, condition, next_steps, state, input, output, started_at, completed_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            ON DUPLICATE KEY UPDATE
            state = VALUES(state),
            input = VALUES(input),
            output = VALUES(output),
            started_at = VALUES(started_at),
            completed_at = VALUES(completed_at)
        `, workflow.ID, step.ID, step.Name, step.Type, step.TaskID, step.Condition,
           marshalNextSteps(step.NextSteps), step.State, marshalContext(step.Input),
           marshalContext(step.Output), step.StartedAt, step.CompletedAt)
        
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (ws *WorkflowStorage) LoadWorkflow(workflowID string) (*Workflow, error) {
    // 加载工作流基本信息
    var workflow Workflow
    err := ws.db.QueryRow(`
        SELECT id, name, version, state, context, created_at, started_at, completed_at
        FROM workflows WHERE id = ?
    `, workflowID).Scan(&workflow.ID, &workflow.Name, &workflow.Version, &workflow.State,
                       &workflow.Context, &workflow.CreatedAt, &workflow.StartedAt, &workflow.CompletedAt)
    
    if err != nil {
        return nil, err
    }
    
    // 加载步骤信息
    rows, err := ws.db.Query(`
        SELECT step_id, name, type, task_id, condition, next_steps, state, input, output, started_at, completed_at
        FROM workflow_steps WHERE workflow_id = ?
        ORDER BY step_id
    `, workflowID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var step Step
        var contextStr, nextStepsStr, inputStr, outputStr string
        
        err := rows.Scan(&step.ID, &step.Name, &step.Type, &step.TaskID, &step.Condition,
                        &nextStepsStr, &step.State, &inputStr, &outputStr,
                        &step.StartedAt, &step.CompletedAt)
        
        if err != nil {
            return nil, err
        }
        
        step.Context = unmarshalContext(contextStr)
        step.NextSteps = unmarshalNextSteps(nextStepsStr)
        step.Input = unmarshalContext(inputStr)
        step.Output = unmarshalContext(outputStr)
        
        workflow.Steps = append(workflow.Steps, &step)
    }
    
    return &workflow, nil
}

func (ws *WorkflowStorage) SaveStateChange(change *StateChange) error {
    _, err := ws.db.Exec(`
        INSERT INTO state_changes (workflow_id, step_id, old_state, new_state, timestamp, reason)
        VALUES (?, ?, ?, ?, ?, ?)
    `, change.WorkflowID, change.StepID, change.OldState, change.NewState,
       change.Timestamp, change.Reason)
    
    return err
}

func marshalContext(context map[string]interface{}) string {
    data, _ := json.Marshal(context)
    return string(data)
}

func unmarshalContext(data string) map[string]interface{} {
    var context map[string]interface{}
    json.Unmarshal([]byte(data), &context)
    return context
}

func marshalNextSteps(steps []string) string {
    data, _ := json.Marshal(steps)
    return string(data)
}

func unmarshalNextSteps(data string) []string {
    var steps []string
    json.Unmarshal([]byte(data), &steps)
    return steps
}
```

## 4. 监控和日志

### 4.1 工作流监控器

```go
// 工作流监控器
type WorkflowMonitor struct {
    workflows map[string]*Workflow
    metrics   *WorkflowMetrics
    alerts    chan WorkflowAlert
    mu        sync.RWMutex
}

type WorkflowMetrics struct {
    TotalWorkflows    int64
    RunningWorkflows  int64
    CompletedWorkflows int64
    FailedWorkflows   int64
    AvgExecutionTime  time.Duration
    mu                sync.Mutex
}

type WorkflowAlert struct {
    WorkflowID string
    Type       string
    Message    string
    Severity   string
    Timestamp  time.Time
}

func (wm *WorkflowMonitor) StartMonitoring() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        wm.collectMetrics()
        wm.checkAlerts()
    }
}

func (wm *WorkflowMonitor) collectMetrics() {
    wm.mu.RLock()
    defer wm.mu.RUnlock()
    
    wm.metrics.mu.Lock()
    defer wm.metrics.mu.Unlock()
    
    total := int64(len(wm.workflows))
    running := int64(0)
    completed := int64(0)
    failed := int64(0)
    totalExecutionTime := time.Duration(0)
    completedCount := int64(0)
    
    for _, workflow := range wm.workflows {
        switch workflow.State {
        case Running:
            running++
        case Completed:
            completed++
            if !workflow.CompletedAt.IsZero() && !workflow.StartedAt.IsZero() {
                executionTime := workflow.CompletedAt.Sub(workflow.StartedAt)
                totalExecutionTime += executionTime
                completedCount++
            }
        case Failed:
            failed++
        }
    }
    
    wm.metrics.TotalWorkflows = total
    wm.metrics.RunningWorkflows = running
    wm.metrics.CompletedWorkflows = completed
    wm.metrics.FailedWorkflows = failed
    
    if completedCount > 0 {
        wm.metrics.AvgExecutionTime = totalExecutionTime / completedCount
    }
}

func (wm *WorkflowMonitor) checkAlerts() {
    wm.mu.RLock()
    defer wm.mu.RUnlock()
    
    now := time.Now()
    
    for _, workflow := range wm.workflows {
        // 检查长时间运行的工作流
        if workflow.State == Running && !workflow.StartedAt.IsZero() {
            executionTime := now.Sub(workflow.StartedAt)
            if executionTime > 1*time.Hour {
                alert := WorkflowAlert{
                    WorkflowID: workflow.ID,
                    Type:       "long_running",
                    Message:    fmt.Sprintf("Workflow %s has been running for %v", workflow.ID, executionTime),
                    Severity:   "warning",
                    Timestamp:  now,
                }
                wm.alerts <- alert
            }
        }
        
        // 检查失败的工作流
        if workflow.State == Failed {
            alert := WorkflowAlert{
                WorkflowID: workflow.ID,
                Type:       "failed",
                Message:    fmt.Sprintf("Workflow %s has failed", workflow.ID),
                Severity:   "error",
                Timestamp:  now,
            }
            wm.alerts <- alert
        }
    }
}

func (wm *WorkflowMonitor) ProcessAlerts() {
    for alert := range wm.alerts {
        wm.handleAlert(alert)
    }
}

func (wm *WorkflowMonitor) handleAlert(alert WorkflowAlert) {
    log.Printf("Workflow Alert [%s]: %s - %s", alert.Severity, alert.Type, alert.Message)
    
    switch alert.Type {
    case "long_running":
        wm.handleLongRunningWorkflow(alert)
    case "failed":
        wm.handleFailedWorkflow(alert)
    }
}

func (wm *WorkflowMonitor) handleLongRunningWorkflow(alert WorkflowAlert) {
    // 处理长时间运行的工作流
    log.Printf("Handling long running workflow: %s", alert.WorkflowID)
    
    // 可以采取的措施：
    // 1. 发送通知给管理员
    // 2. 尝试重启工作流
    // 3. 记录详细日志
}

func (wm *WorkflowMonitor) handleFailedWorkflow(alert WorkflowAlert) {
    // 处理失败的工作流
    log.Printf("Handling failed workflow: %s", alert.WorkflowID)
    
    // 可以采取的措施：
    // 1. 发送通知给管理员
    // 2. 尝试重试工作流
    // 3. 记录错误详情
}

func (wm *WorkflowMonitor) GetMetrics() *WorkflowMetrics {
    wm.metrics.mu.Lock()
    defer wm.metrics.mu.Unlock()
    
    return &WorkflowMetrics{
        TotalWorkflows:     wm.metrics.TotalWorkflows,
        RunningWorkflows:   wm.metrics.RunningWorkflows,
        CompletedWorkflows: wm.metrics.CompletedWorkflows,
        FailedWorkflows:    wm.metrics.FailedWorkflows,
        AvgExecutionTime:   wm.metrics.AvgExecutionTime,
    }
}
```

### 4.2 工作流日志器

```go
// 工作流日志器
type WorkflowLogger struct {
    logger *log.Logger
    file   *os.File
}

func NewWorkflowLogger(filename string) (*WorkflowLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    logger := log.New(file, "", log.LstdFlags)
    
    return &WorkflowLogger{
        logger: logger,
        file:   file,
    }, nil
}

func (wl *WorkflowLogger) LogWorkflowEvent(workflowID string, event string, details map[string]interface{}) {
    data := map[string]interface{}{
        "workflow_id": workflowID,
        "event":       event,
        "timestamp":   time.Now(),
        "details":     details,
    }
    
    jsonData, _ := json.Marshal(data)
    wl.logger.Printf("WORKFLOW: %s", string(jsonData))
}

func (wl *WorkflowLogger) LogStepEvent(workflowID, stepID string, event string, details map[string]interface{}) {
    data := map[string]interface{}{
        "workflow_id": workflowID,
        "step_id":     stepID,
        "event":       event,
        "timestamp":   time.Now(),
        "details":     details,
    }
    
    jsonData, _ := json.Marshal(data)
    wl.logger.Printf("STEP: %s", string(jsonData))
}

func (wl *WorkflowLogger) LogError(workflowID string, error error, details map[string]interface{}) {
    data := map[string]interface{}{
        "workflow_id": workflowID,
        "error":       error.Error(),
        "timestamp":   time.Now(),
        "details":     details,
    }
    
    jsonData, _ := json.Marshal(data)
    wl.logger.Printf("ERROR: %s", string(jsonData))
}

func (wl *WorkflowLogger) Close() error {
    return wl.file.Close()
}
```

## 5. 工作流定义语言

### 5.1 JSON工作流定义

```go
// 工作流定义
type WorkflowDefinition struct {
    Name        string                 `json:"name"`
    Version     string                 `json:"version"`
    Description string                 `json:"description"`
    Steps       []*StepDefinition      `json:"steps"`
    Variables   map[string]interface{} `json:"variables"`
}

type StepDefinition struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        string                 `json:"type"`
    TaskID      string                 `json:"task_id,omitempty"`
    Condition   string                 `json:"condition,omitempty"`
    NextSteps   []string               `json:"next_steps,omitempty"`
    Input       map[string]interface{} `json:"input,omitempty"`
    Retry       *RetryConfig           `json:"retry,omitempty"`
    Timeout     time.Duration          `json:"timeout,omitempty"`
}

type RetryConfig struct {
    MaxAttempts int           `json:"max_attempts"`
    Delay       time.Duration `json:"delay"`
    Backoff     float64       `json:"backoff"`
}

// 工作流定义解析器
type WorkflowDefinitionParser struct{}

func (wdp *WorkflowDefinitionParser) ParseJSON(data []byte) (*WorkflowDefinition, error) {
    var definition WorkflowDefinition
    if err := json.Unmarshal(data, &definition); err != nil {
        return nil, err
    }
    
    // 验证工作流定义
    if err := wdp.validateDefinition(&definition); err != nil {
        return nil, err
    }
    
    return &definition, nil
}

func (wdp *WorkflowDefinitionParser) validateDefinition(definition *WorkflowDefinition) error {
    if definition.Name == "" {
        return fmt.Errorf("workflow name is required")
    }
    
    if len(definition.Steps) == 0 {
        return fmt.Errorf("workflow must have at least one step")
    }
    
    // 检查步骤ID的唯一性
    stepIDs := make(map[string]bool)
    for _, step := range definition.Steps {
        if step.ID == "" {
            return fmt.Errorf("step ID is required")
        }
        
        if stepIDs[step.ID] {
            return fmt.Errorf("duplicate step ID: %s", step.ID)
        }
        stepIDs[step.ID] = true
    }
    
    // 检查步骤引用的有效性
    for _, step := range definition.Steps {
        for _, nextStepID := range step.NextSteps {
            if !stepIDs[nextStepID] {
                return fmt.Errorf("step %s references non-existent step: %s", step.ID, nextStepID)
            }
        }
    }
    
    return nil
}

// 示例工作流定义
func (wdp *WorkflowDefinitionParser) GetExampleDefinition() string {
    return `{
        "name": "data_processing_workflow",
        "version": "1.0.0",
        "description": "A workflow for processing and validating data",
        "steps": [
            {
                "id": "validate",
                "name": "Validate Data",
                "type": "task",
                "task_id": "validate_data",
                "next_steps": ["process"],
                "input": {
                    "data": "{{.input.data}}"
                },
                "retry": {
                    "max_attempts": 3,
                    "delay": "5s",
                    "backoff": 2.0
                }
            },
            {
                "id": "process",
                "name": "Process Data",
                "type": "task",
                "task_id": "process_data",
                "next_steps": ["notify"],
                "input": {
                    "data": "{{.steps.validate.output.processed_data}}"
                }
            },
            {
                "id": "notify",
                "name": "Send Notification",
                "type": "task",
                "task_id": "send_notification",
                "input": {
                    "message": "Data processing completed successfully"
                }
            }
        ],
        "variables": {
            "max_retries": 3,
            "timeout": "30m"
        }
    }`
}
```

## 总结

工作流系统是一个复杂的系统，用于自动化和管理业务流程。成功的工作流系统需要：

**关键设计原则**：

1. 可扩展性：支持复杂的工作流定义
2. 可靠性：确保工作流的正确执行
3. 可监控性：提供全面的监控和日志
4. 容错性：处理执行过程中的错误
5. 性能：支持高并发的工作流执行

**常见挑战**：

1. 工作流定义的复杂性
2. 状态管理和持久化
3. 错误处理和重试机制
4. 性能和可扩展性
5. 监控和调试

工作流系统的设计需要根据具体的业务需求来选择合适的架构和技术，确保系统能够满足业务流程自动化的要求。
