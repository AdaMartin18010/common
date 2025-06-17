# 3. 工作流模式

## 3.1 工作流模式理论基础

### 3.1.1 模式分类体系

**定义 3.1** (工作流模式): 工作流模式是一个三元组 $\mathcal{P} = (S, T, C)$，其中：

- $S$ 是结构模式集合
- $T$ 是任务模式集合  
- $C$ 是控制模式集合

**模式分类**:

```latex
\text{WorkflowPatterns} = \text{Sequential} \cup \text{Parallel} \cup \text{Choice} \cup \text{Iteration} \cup \text{Exception}
```

### 3.1.2 模式形式化定义

**定义 3.2** (顺序模式): 顺序模式是一个二元组 $\text{Seq}(t_1, t_2)$，满足：

$$\text{Seq}(t_1, t_2) = \{s_1 \rightarrow s_2 \rightarrow s_3 \mid s_1 \xrightarrow{t_1} s_2 \xrightarrow{t_2} s_3\}$$

**定义 3.3** (并行模式): 并行模式是一个二元组 $\text{Par}(t_1, t_2)$，满足：

$$\text{Par}(t_1, t_2) = \{s_1 \rightarrow s_2 \mid s_1 \xrightarrow{t_1 \parallel t_2} s_2\}$$

**定义 3.4** (选择模式): 选择模式是一个三元组 $\text{Choice}(c, t_1, t_2)$，满足：

$$\text{Choice}(c, t_1, t_2) = \{s_1 \rightarrow s_2 \mid s_1 \xrightarrow{\text{if } c \text{ then } t_1 \text{ else } t_2} s_2\}$$

## 3.2 基础工作流模式

### 3.2.1 顺序执行模式

```go
// SequentialPattern 顺序执行模式
type SequentialPattern struct {
    tasks []Task
}

// NewSequentialPattern 创建顺序执行模式
func NewSequentialPattern(tasks ...Task) *SequentialPattern {
    return &SequentialPattern{
        tasks: tasks,
    }
}

// Execute 执行顺序模式
func (sp *SequentialPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    currentInput := input
    
    for i, task := range sp.tasks {
        // 执行任务
        result, err := task.Execute(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("task %d failed: %w", i, err)
        }
        
        // 合并结果
        currentInput = mergeContext(currentInput, result)
        
        // 检查是否需要继续
        if shouldStop, exists := currentInput["_stop"]; exists && shouldStop.(bool) {
            break
        }
    }
    
    return currentInput, nil
}

// Task 任务接口
type Task interface {
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
    GetID() string
    GetType() string
}

// BaseTask 基础任务实现
type BaseTask struct {
    ID       string
    Type     string
    Handler  TaskHandler
}

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)

// Execute 执行任务
func (bt *BaseTask) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return bt.Handler(ctx, input)
}

// GetID 获取任务ID
func (bt *BaseTask) GetID() string {
    return bt.ID
}

// GetType 获取任务类型
func (bt *BaseTask) GetType() string {
    return bt.Type
}
```

### 3.2.2 并行执行模式

```go
// ParallelPattern 并行执行模式
type ParallelPattern struct {
    tasks []Task
    mode  ParallelMode
}

// ParallelMode 并行模式类型
type ParallelMode string

const (
    ParallelModeAll     ParallelMode = "ALL"     // 等待所有任务完成
    ParallelModeAny     ParallelMode = "ANY"     // 任一任务完成即可
    ParallelModeN       ParallelMode = "N"       // N个任务完成即可
)

// NewParallelPattern 创建并行执行模式
func NewParallelPattern(mode ParallelMode, tasks ...Task) *ParallelPattern {
    return &ParallelPattern{
        tasks: tasks,
        mode:  mode,
    }
}

// Execute 执行并行模式
func (pp *ParallelPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    switch pp.mode {
    case ParallelModeAll:
        return pp.executeAll(ctx, input)
    case ParallelModeAny:
        return pp.executeAny(ctx, input)
    case ParallelModeN:
        return pp.executeN(ctx, input, 1) // 默认N=1
    default:
        return nil, fmt.Errorf("unknown parallel mode: %s", pp.mode)
    }
}

// executeAll 执行所有任务
func (pp *ParallelPattern) executeAll(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    var wg sync.WaitGroup
    results := make([]map[string]interface{}, len(pp.tasks))
    errors := make([]error, len(pp.tasks))
    
    // 启动所有任务
    for i, task := range pp.tasks {
        wg.Add(1)
        go func(index int, t Task) {
            defer wg.Done()
            result, err := t.Execute(ctx, input)
            results[index] = result
            errors[index] = err
        }(i, task)
    }
    
    // 等待所有任务完成
    wg.Wait()
    
    // 检查错误
    for i, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("task %d failed: %w", i, err)
        }
    }
    
    // 合并所有结果
    finalResult := make(map[string]interface{})
    for _, result := range results {
        finalResult = mergeContext(finalResult, result)
    }
    
    return finalResult, nil
}

// executeAny 执行任一任务
func (pp *ParallelPattern) executeAny(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    resultChan := make(chan map[string]interface{}, 1)
    errorChan := make(chan error, 1)
    
    // 启动所有任务
    for _, task := range pp.tasks {
        go func(t Task) {
            result, err := t.Execute(ctx, input)
            if err != nil {
                select {
                case errorChan <- err:
                default:
                }
                return
            }
            
            select {
            case resultChan <- result:
                cancel() // 取消其他任务
            default:
            }
        }(task)
    }
    
    // 等待第一个完成的任务
    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errorChan:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// executeN 执行N个任务
func (pp *ParallelPattern) executeN(ctx context.Context, input map[string]interface{}, n int) (map[string]interface{}, error) {
    if n <= 0 || n > len(pp.tasks) {
        return nil, fmt.Errorf("invalid N value: %d", n)
    }
    
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    resultChan := make(chan map[string]interface{}, n)
    errorChan := make(chan error, len(pp.tasks))
    
    // 启动所有任务
    for _, task := range pp.tasks {
        go func(t Task) {
            result, err := t.Execute(ctx, input)
            if err != nil {
                select {
                case errorChan <- err:
                default:
                }
                return
            }
            
            select {
            case resultChan <- result:
            default:
            }
        }(task)
    }
    
    // 收集N个结果
    results := make([]map[string]interface{}, 0, n)
    errors := make([]error, 0)
    
    for i := 0; i < len(pp.tasks); i++ {
        select {
        case result := <-resultChan:
            results = append(results, result)
            if len(results) >= n {
                // 合并结果
                finalResult := make(map[string]interface{})
                for _, r := range results {
                    finalResult = mergeContext(finalResult, r)
                }
                return finalResult, nil
            }
        case err := <-errorChan:
            errors = append(errors, err)
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    // 如果没有足够的成功结果，返回错误
    if len(errors) > 0 {
        return nil, fmt.Errorf("insufficient successful tasks: %v", errors)
    }
    
    return nil, fmt.Errorf("unexpected end of execution")
}
```

### 3.2.3 条件选择模式

```go
// ConditionalPattern 条件选择模式
type ConditionalPattern struct {
    conditions []Condition
    tasks      []Task
    defaultTask Task
}

// Condition 条件接口
type Condition interface {
    Evaluate(ctx context.Context, input map[string]interface{}) (bool, error)
    GetDescription() string
}

// SimpleCondition 简单条件
type SimpleCondition struct {
    expression string
    evaluator  ConditionEvaluator
}

// ConditionEvaluator 条件求值器
type ConditionEvaluator func(ctx context.Context, input map[string]interface{}, expression string) (bool, error)

// Evaluate 求值条件
func (sc *SimpleCondition) Evaluate(ctx context.Context, input map[string]interface{}) (bool, error) {
    return sc.evaluator(ctx, input, sc.expression)
}

// GetDescription 获取条件描述
func (sc *SimpleCondition) GetDescription() string {
    return sc.expression
}

// NewConditionalPattern 创建条件选择模式
func NewConditionalPattern() *ConditionalPattern {
    return &ConditionalPattern{
        conditions: make([]Condition, 0),
        tasks:      make([]Task, 0),
    }
}

// AddCondition 添加条件分支
func (cp *ConditionalPattern) AddCondition(condition Condition, task Task) {
    cp.conditions = append(cp.conditions, condition)
    cp.tasks = append(cp.tasks, task)
}

// SetDefault 设置默认任务
func (cp *ConditionalPattern) SetDefault(task Task) {
    cp.defaultTask = task
}

// Execute 执行条件选择模式
func (cp *ConditionalPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // 评估所有条件
    for i, condition := range cp.conditions {
        result, err := condition.Evaluate(ctx, input)
        if err != nil {
            return nil, fmt.Errorf("condition evaluation failed: %w", err)
        }
        
        if result {
            // 执行对应的任务
            return cp.tasks[i].Execute(ctx, input)
        }
    }
    
    // 如果没有条件满足，执行默认任务
    if cp.defaultTask != nil {
        return cp.defaultTask.Execute(ctx, input)
    }
    
    return input, nil
}

// ExpressionCondition 表达式条件
type ExpressionCondition struct {
    expression string
}

// NewExpressionCondition 创建表达式条件
func NewExpressionCondition(expression string) *ExpressionCondition {
    return &ExpressionCondition{
        expression: expression,
    }
}

// Evaluate 求值表达式条件
func (ec *ExpressionCondition) Evaluate(ctx context.Context, input map[string]interface{}) (bool, error) {
    // 这里可以使用表达式引擎，如govaluate
    // 简化实现，支持基本的比较操作
    return evaluateExpression(ec.expression, input)
}

// evaluateExpression 求值表达式
func evaluateExpression(expression string, input map[string]interface{}) (bool, error) {
    // 简单的表达式求值实现
    // 支持: ==, !=, >, <, >=, <=, &&, ||
    
    // 示例: "value > 10" 或 "status == 'active'"
    parts := strings.Split(expression, " ")
    if len(parts) != 3 {
        return false, fmt.Errorf("invalid expression format: %s", expression)
    }
    
    left := parts[0]
    operator := parts[1]
    right := parts[2]
    
    // 获取左值
    leftValue, exists := input[left]
    if !exists {
        return false, fmt.Errorf("variable not found: %s", left)
    }
    
    // 解析右值
    var rightValue interface{}
    if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
        // 字符串值
        rightValue = strings.Trim(right, "'")
    } else {
        // 数值或其他类型
        if num, err := strconv.ParseFloat(right, 64); err == nil {
            rightValue = num
        } else {
            rightValue = right
        }
    }
    
    // 执行比较
    switch operator {
    case "==":
        return reflect.DeepEqual(leftValue, rightValue), nil
    case "!=":
        return !reflect.DeepEqual(leftValue, rightValue), nil
    case ">":
        return compare(leftValue, rightValue) > 0, nil
    case "<":
        return compare(leftValue, rightValue) < 0, nil
    case ">=":
        return compare(leftValue, rightValue) >= 0, nil
    case "<=":
        return compare(leftValue, rightValue) <= 0, nil
    default:
        return false, fmt.Errorf("unsupported operator: %s", operator)
    }
}

// compare 比较两个值
func compare(a, b interface{}) int {
    switch va := a.(type) {
    case int:
        if vb, ok := b.(int); ok {
            if va < vb {
                return -1
            } else if va > vb {
                return 1
            }
            return 0
        }
    case float64:
        if vb, ok := b.(float64); ok {
            if va < vb {
                return -1
            } else if va > vb {
                return 1
            }
            return 0
        }
    case string:
        if vb, ok := b.(string); ok {
            return strings.Compare(va, vb)
        }
    }
    return 0
}
```

## 3.3 高级工作流模式

### 3.3.1 循环迭代模式

```go
// IterationPattern 循环迭代模式
type IterationPattern struct {
    task        Task
    condition   Condition
    maxIterations int
}

// NewIterationPattern 创建循环迭代模式
func NewIterationPattern(task Task, condition Condition, maxIterations int) *IterationPattern {
    return &IterationPattern{
        task:          task,
        condition:     condition,
        maxIterations: maxIterations,
    }
}

// Execute 执行循环迭代模式
func (ip *IterationPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    currentInput := input
    iteration := 0
    
    for {
        // 检查最大迭代次数
        if iteration >= ip.maxIterations {
            return currentInput, fmt.Errorf("max iterations reached: %d", ip.maxIterations)
        }
        
        // 检查循环条件
        shouldContinue, err := ip.condition.Evaluate(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("condition evaluation failed: %w", err)
        }
        
        if !shouldContinue {
            break
        }
        
        // 执行任务
        result, err := ip.task.Execute(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("iteration %d failed: %w", iteration, err)
        }
        
        // 更新输入
        currentInput = mergeContext(currentInput, result)
        currentInput["_iteration"] = iteration
        
        iteration++
    }
    
    return currentInput, nil
}

// WhilePattern While循环模式
type WhilePattern struct {
    condition Condition
    task      Task
    maxIterations int
}

// NewWhilePattern 创建While循环模式
func NewWhilePattern(condition Condition, task Task, maxIterations int) *WhilePattern {
    return &WhilePattern{
        condition:     condition,
        task:          task,
        maxIterations: maxIterations,
    }
}

// Execute 执行While循环模式
func (wp *WhilePattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    currentInput := input
    iteration := 0
    
    for {
        // 检查最大迭代次数
        if iteration >= wp.maxIterations {
            return currentInput, fmt.Errorf("max iterations reached: %d", wp.maxIterations)
        }
        
        // 检查循环条件
        shouldContinue, err := wp.condition.Evaluate(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("condition evaluation failed: %w", err)
        }
        
        if !shouldContinue {
            break
        }
        
        // 执行任务
        result, err := wp.task.Execute(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("iteration %d failed: %w", iteration, err)
        }
        
        // 更新输入
        currentInput = mergeContext(currentInput, result)
        currentInput["_iteration"] = iteration
        
        iteration++
    }
    
    return currentInput, nil
}

// ForEachPattern ForEach循环模式
type ForEachPattern struct {
    task      Task
    collection string
}

// NewForEachPattern 创建ForEach循环模式
func NewForEachPattern(task Task, collection string) *ForEachPattern {
    return &ForEachPattern{
        task:       task,
        collection: collection,
    }
}

// Execute 执行ForEach循环模式
func (fep *ForEachPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // 获取集合
    collectionValue, exists := input[fep.collection]
    if !exists {
        return nil, fmt.Errorf("collection not found: %s", fep.collection)
    }
    
    // 转换为切片
    var items []interface{}
    switch v := collectionValue.(type) {
    case []interface{}:
        items = v
    case []string:
        for _, item := range v {
            items = append(items, item)
        }
    case []int:
        for _, item := range v {
            items = append(items, item)
        }
    default:
        return nil, fmt.Errorf("unsupported collection type: %T", collectionValue)
    }
    
    results := make([]interface{}, 0, len(items))
    
    // 遍历集合
    for i, item := range items {
        // 创建当前项的上下文
        itemInput := make(map[string]interface{})
        for k, v := range input {
            itemInput[k] = v
        }
        itemInput["_item"] = item
        itemInput["_index"] = i
        
        // 执行任务
        result, err := fep.task.Execute(ctx, itemInput)
        if err != nil {
            return nil, fmt.Errorf("item %d failed: %w", i, err)
        }
        
        results = append(results, result)
    }
    
    // 合并结果
    finalResult := make(map[string]interface{})
    for k, v := range input {
        finalResult[k] = v
    }
    finalResult["_results"] = results
    
    return finalResult, nil
}
```

### 3.3.2 异常处理模式

```go
// ExceptionPattern 异常处理模式
type ExceptionPattern struct {
    mainTask     Task
    errorHandler ErrorHandler
    finallyTask  Task
}

// ErrorHandler 错误处理器
type ErrorHandler interface {
    Handle(ctx context.Context, err error, input map[string]interface{}) (map[string]interface{}, error)
}

// SimpleErrorHandler 简单错误处理器
type SimpleErrorHandler struct {
    handler func(ctx context.Context, err error, input map[string]interface{}) (map[string]interface{}, error)
}

// Handle 处理错误
func (seh *SimpleErrorHandler) Handle(ctx context.Context, err error, input map[string]interface{}) (map[string]interface{}, error) {
    return seh.handler(ctx, err, input)
}

// NewExceptionPattern 创建异常处理模式
func NewExceptionPattern(mainTask Task) *ExceptionPattern {
    return &ExceptionPattern{
        mainTask: mainTask,
    }
}

// SetErrorHandler 设置错误处理器
func (ep *ExceptionPattern) SetErrorHandler(handler ErrorHandler) {
    ep.errorHandler = handler
}

// SetFinallyTask 设置Finally任务
func (ep *ExceptionPattern) SetFinallyTask(task Task) {
    ep.finallyTask = task
}

// Execute 执行异常处理模式
func (ep *ExceptionPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    var result map[string]interface{}
    var err error
    
    // 执行主任务
    func() {
        defer func() {
            if r := recover(); r != nil {
                err = fmt.Errorf("panic recovered: %v", r)
            }
        }()
        
        result, err = ep.mainTask.Execute(ctx, input)
    }()
    
    // 如果有错误且设置了错误处理器
    if err != nil && ep.errorHandler != nil {
        result, err = ep.errorHandler.Handle(ctx, err, input)
    }
    
    // 执行Finally任务
    if ep.finallyTask != nil {
        finallyResult, finallyErr := ep.finallyTask.Execute(ctx, input)
        if finallyErr != nil {
            // Finally任务失败，记录但不影响主流程
            log.Printf("Finally task failed: %v", finallyErr)
        }
        
        // 合并Finally结果
        if finallyResult != nil {
            result = mergeContext(result, finallyResult)
        }
    }
    
    return result, err
}

// RetryPattern 重试模式
type RetryPattern struct {
    task         Task
    maxRetries   int
    backoff      BackoffStrategy
    retryCondition RetryCondition
}

// BackoffStrategy 退避策略
type BackoffStrategy interface {
    GetDelay(attempt int) time.Duration
}

// FixedBackoff 固定退避策略
type FixedBackoff struct {
    delay time.Duration
}

// GetDelay 获取延迟时间
func (fb *FixedBackoff) GetDelay(attempt int) time.Duration {
    return fb.delay
}

// ExponentialBackoff 指数退避策略
type ExponentialBackoff struct {
    initialDelay time.Duration
    maxDelay     time.Duration
    multiplier   float64
}

// GetDelay 获取延迟时间
func (eb *ExponentialBackoff) GetDelay(attempt int) time.Duration {
    delay := time.Duration(float64(eb.initialDelay) * math.Pow(eb.multiplier, float64(attempt)))
    if delay > eb.maxDelay {
        delay = eb.maxDelay
    }
    return delay
}

// RetryCondition 重试条件
type RetryCondition interface {
    ShouldRetry(err error) bool
}

// SimpleRetryCondition 简单重试条件
type SimpleRetryCondition struct {
    retryableErrors []string
}

// ShouldRetry 判断是否应该重试
func (src *SimpleRetryCondition) ShouldRetry(err error) bool {
    if err == nil {
        return false
    }
    
    errStr := err.Error()
    for _, retryableError := range src.retryableErrors {
        if strings.Contains(errStr, retryableError) {
            return true
        }
    }
    return false
}

// NewRetryPattern 创建重试模式
func NewRetryPattern(task Task, maxRetries int) *RetryPattern {
    return &RetryPattern{
        task:       task,
        maxRetries: maxRetries,
        backoff:    &FixedBackoff{delay: time.Second},
        retryCondition: &SimpleRetryCondition{
            retryableErrors: []string{"timeout", "connection", "temporary"},
        },
    }
}

// SetBackoff 设置退避策略
func (rp *RetryPattern) SetBackoff(backoff BackoffStrategy) {
    rp.backoff = backoff
}

// SetRetryCondition 设置重试条件
func (rp *RetryPattern) SetRetryCondition(condition RetryCondition) {
    rp.retryCondition = condition
}

// Execute 执行重试模式
func (rp *RetryPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    var lastErr error
    
    for attempt := 0; attempt <= rp.maxRetries; attempt++ {
        // 执行任务
        result, err := rp.task.Execute(ctx, input)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // 检查是否应该重试
        if !rp.retryCondition.ShouldRetry(err) {
            return nil, err
        }
        
        // 如果是最后一次尝试，不等待
        if attempt == rp.maxRetries {
            break
        }
        
        // 等待退避时间
        delay := rp.backoff.GetDelay(attempt)
        select {
        case <-time.After(delay):
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## 3.4 复合工作流模式

### 3.4.1 模式组合器

```go
// PatternComposer 模式组合器
type PatternComposer struct {
    patterns []WorkflowPattern
}

// WorkflowPattern 工作流模式接口
type WorkflowPattern interface {
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}

// NewPatternComposer 创建模式组合器
func NewPatternComposer() *PatternComposer {
    return &PatternComposer{
        patterns: make([]WorkflowPattern, 0),
    }
}

// AddPattern 添加模式
func (pc *PatternComposer) AddPattern(pattern WorkflowPattern) {
    pc.patterns = append(pc.patterns, pattern)
}

// Compose 组合模式
func (pc *PatternComposer) Compose() WorkflowPattern {
    if len(pc.patterns) == 0 {
        return &NoOpPattern{}
    }
    
    if len(pc.patterns) == 1 {
        return pc.patterns[0]
    }
    
    // 创建顺序组合
    return &SequentialPattern{
        tasks: pc.convertToTasks(pc.patterns),
    }
}

// convertToTasks 将模式转换为任务
func (pc *PatternComposer) convertToTasks(patterns []WorkflowPattern) []Task {
    tasks := make([]Task, len(patterns))
    for i, pattern := range patterns {
        tasks[i] = &PatternTask{pattern: pattern}
    }
    return tasks
}

// PatternTask 模式任务包装器
type PatternTask struct {
    pattern WorkflowPattern
}

// Execute 执行模式任务
func (pt *PatternTask) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return pt.pattern.Execute(ctx, input)
}

// GetID 获取任务ID
func (pt *PatternTask) GetID() string {
    return fmt.Sprintf("pattern_%p", pt.pattern)
}

// GetType 获取任务类型
func (pt *PatternTask) GetType() string {
    return "pattern"
}

// NoOpPattern 空操作模式
type NoOpPattern struct{}

// Execute 执行空操作
func (nop *NoOpPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    return input, nil
}
```

### 3.4.2 复杂业务模式

```go
// BusinessProcessPattern 业务流程模式
type BusinessProcessPattern struct {
    validationTask   Task
    processingTask   Task
    notificationTask Task
    errorHandler     ErrorHandler
}

// NewBusinessProcessPattern 创建业务流程模式
func NewBusinessProcessPattern(validationTask, processingTask, notificationTask Task) *BusinessProcessPattern {
    return &BusinessProcessPattern{
        validationTask:   validationTask,
        processingTask:   processingTask,
        notificationTask: notificationTask,
    }
}

// SetErrorHandler 设置错误处理器
func (bpp *BusinessProcessPattern) SetErrorHandler(handler ErrorHandler) {
    bpp.errorHandler = handler
}

// Execute 执行业务流程模式
func (bpp *BusinessProcessPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // 创建异常处理模式
    exceptionPattern := NewExceptionPattern(&SequentialPattern{
        tasks: []Task{
            bpp.validationTask,
            bpp.processingTask,
            bpp.notificationTask,
        },
    })
    
    // 设置错误处理器
    if bpp.errorHandler != nil {
        exceptionPattern.SetErrorHandler(bpp.errorHandler)
    }
    
    return exceptionPattern.Execute(ctx, input)
}

// DataProcessingPattern 数据处理模式
type DataProcessingPattern struct {
    extractTask   Task
    transformTask Task
    loadTask      Task
    batchSize     int
}

// NewDataProcessingPattern 创建数据处理模式
func NewDataProcessingPattern(extractTask, transformTask, loadTask Task, batchSize int) *DataProcessingPattern {
    return &DataProcessingPattern{
        extractTask:   extractTask,
        transformTask: transformTask,
        loadTask:      loadTask,
        batchSize:     batchSize,
    }
}

// Execute 执行数据处理模式
func (dpp *DataProcessingPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    // 提取数据
    extractResult, err := dpp.extractTask.Execute(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("extract failed: %w", err)
    }
    
    // 获取数据列表
    dataList, ok := extractResult["data"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid data format")
    }
    
    // 分批处理
    processedData := make([]interface{}, 0, len(dataList))
    
    for i := 0; i < len(dataList); i += dpp.batchSize {
        end := i + dpp.batchSize
        if end > len(dataList) {
            end = len(dataList)
        }
        
        batch := dataList[i:end]
        
        // 转换批次数据
        transformInput := make(map[string]interface{})
        for k, v := range input {
            transformInput[k] = v
        }
        transformInput["batch"] = batch
        
        transformResult, err := dpp.transformTask.Execute(ctx, transformInput)
        if err != nil {
            return nil, fmt.Errorf("transform batch %d failed: %w", i/dpp.batchSize, err)
        }
        
        // 加载批次数据
        loadInput := make(map[string]interface{})
        for k, v := range input {
            loadInput[k] = v
        }
        loadInput["transformed_batch"] = transformResult["transformed_batch"]
        
        loadResult, err := dpp.loadTask.Execute(ctx, loadInput)
        if err != nil {
            return nil, fmt.Errorf("load batch %d failed: %w", i/dpp.batchSize, err)
        }
        
        processedData = append(processedData, loadResult["processed_batch"])
    }
    
    // 返回处理结果
    result := make(map[string]interface{})
    for k, v := range input {
        result[k] = v
    }
    result["processed_data"] = processedData
    result["total_processed"] = len(processedData)
    
    return result, nil
}
```

## 3.5 模式验证和测试

### 3.5.1 模式验证器

```go
// PatternValidator 模式验证器
type PatternValidator struct {
    rules []ValidationRule
}

// ValidationRule 验证规则
type ValidationRule interface {
    Validate(pattern WorkflowPattern) error
}

// CircularDependencyRule 循环依赖验证规则
type CircularDependencyRule struct{}

// Validate 验证循环依赖
func (cdr *CircularDependencyRule) Validate(pattern WorkflowPattern) error {
    // 检测循环依赖
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    return cdr.detectCycle(pattern, visited, recStack)
}

// detectCycle 检测循环
func (cdr *CircularDependencyRule) detectCycle(pattern WorkflowPattern, visited, recStack map[string]bool) error {
    patternID := fmt.Sprintf("%p", pattern)
    
    if recStack[patternID] {
        return fmt.Errorf("circular dependency detected: %s", patternID)
    }
    
    if visited[patternID] {
        return nil
    }
    
    visited[patternID] = true
    recStack[patternID] = true
    
    // 检查子模式
    if composite, ok := pattern.(*SequentialPattern); ok {
        for _, task := range composite.tasks {
            if patternTask, ok := task.(*PatternTask); ok {
                if err := cdr.detectCycle(patternTask.pattern, visited, recStack); err != nil {
                    return err
                }
            }
        }
    }
    
    recStack[patternID] = false
    return nil
}

// NewPatternValidator 创建模式验证器
func NewPatternValidator() *PatternValidator {
    return &PatternValidator{
        rules: []ValidationRule{
            &CircularDependencyRule{},
        },
    }
}

// AddRule 添加验证规则
func (pv *PatternValidator) AddRule(rule ValidationRule) {
    pv.rules = append(pv.rules, rule)
}

// Validate 验证模式
func (pv *PatternValidator) Validate(pattern WorkflowPattern) error {
    for _, rule := range pv.rules {
        if err := rule.Validate(pattern); err != nil {
            return err
        }
    }
    return nil
}
```

### 3.5.2 模式测试框架

```go
// PatternTestSuite 模式测试套件
type PatternTestSuite struct {
    tests []PatternTest
}

// PatternTest 模式测试
type PatternTest struct {
    Name     string
    Pattern  WorkflowPattern
    Input    map[string]interface{}
    Expected map[string]interface{}
    ShouldError bool
}

// NewPatternTestSuite 创建模式测试套件
func NewPatternTestSuite() *PatternTestSuite {
    return &PatternTestSuite{
        tests: make([]PatternTest, 0),
    }
}

// AddTest 添加测试
func (pts *PatternTestSuite) AddTest(test PatternTest) {
    pts.tests = append(pts.tests, test)
}

// RunTests 运行测试
func (pts *PatternTestSuite) RunTests(ctx context.Context) []TestResult {
    results := make([]TestResult, 0, len(pts.tests))
    
    for _, test := range pts.tests {
        result := pts.runTest(ctx, test)
        results = append(results, result)
    }
    
    return results
}

// TestResult 测试结果
type TestResult struct {
    TestName string
    Passed   bool
    Error    error
    Actual   map[string]interface{}
    Expected map[string]interface{}
}

// runTest 运行单个测试
func (pts *PatternTestSuite) runTest(ctx context.Context, test PatternTest) TestResult {
    result, err := test.Pattern.Execute(ctx, test.Input)
    
    if test.ShouldError {
        if err == nil {
            return TestResult{
                TestName: test.Name,
                Passed:   false,
                Error:    fmt.Errorf("expected error but got none"),
                Actual:   result,
                Expected: test.Expected,
            }
        }
        return TestResult{
            TestName: test.Name,
            Passed:   true,
            Error:    err,
        }
    }
    
    if err != nil {
        return TestResult{
            TestName: test.Name,
            Passed:   false,
            Error:    err,
        }
    }
    
    // 比较结果
    if reflect.DeepEqual(result, test.Expected) {
        return TestResult{
            TestName: test.Name,
            Passed:   true,
            Actual:   result,
            Expected: test.Expected,
        }
    }
    
    return TestResult{
        TestName: test.Name,
        Passed:   false,
        Error:    fmt.Errorf("result mismatch"),
        Actual:   result,
        Expected: test.Expected,
    }
}
```

## 3.6 总结

工作流模式模块涵盖了以下核心内容：

1. **理论基础**: 形式化定义各种工作流模式
2. **基础模式**: 顺序、并行、条件选择等基础模式
3. **高级模式**: 循环迭代、异常处理等高级模式
4. **复合模式**: 模式组合和复杂业务模式
5. **验证测试**: 模式验证和测试框架

这些模式提供了构建复杂工作流的基础构件，支持灵活的业务流程建模和实现。
