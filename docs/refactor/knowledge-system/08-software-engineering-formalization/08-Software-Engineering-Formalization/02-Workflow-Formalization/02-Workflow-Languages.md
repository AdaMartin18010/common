# 02-工作流语言 (Workflow Languages)

## 概述

工作流语言是用于描述和定义工作流的专用领域特定语言(DSL)。本文档基于对 `/docs/model/Software/WorkFlow` 目录的深度分析，建立了完整的工作流语言形式化理论体系。

## 1. 工作流语言基础

### 1.1 语言语法定义

**定义 1.1** (工作流语言语法)
工作流语言的语法是一个四元组 ```latex
G = (V_N, V_T, P, S)
```，其中：

- ```latex
V_N
``` 是非终结符集合
- ```latex
V_T
``` 是终结符集合
- ```latex
P
``` 是产生式规则集合
- ```latex
S
``` 是开始符号

**BNF语法定义**:

```
workflow    ::= 'workflow' ID '{' workflow_body '}'
workflow_body ::= states events transitions metadata
states      ::= 'states' '{' state_def* '}'
state_def   ::= 'state' ID '{' state_props '}'
state_props ::= 'type:' state_type ',' 'actions:' action_list
state_type  ::= 'initial' | 'intermediate' | 'final' | 'error'
events      ::= 'events' '{' event_def* '}'
event_def   ::= 'event' ID '{' event_props '}'
event_props ::= 'type:' event_type ',' 'payload:' payload_def
transitions ::= 'transitions' '{' transition_def* '}'
transition_def ::= 'transition' '{' transition_props '}'
transition_props ::= 'from:' ID ',' 'to:' ID ',' 'event:' ID ',' 'condition:' expr?
metadata    ::= 'metadata' '{' key_value* '}'
```

### 1.2 语言语义模型

**定义 1.2** (工作流语言语义)
工作流语言的语义是一个三元组 ```latex
\mathcal{S} = (D, \mathcal{I}, \mathcal{E})
```，其中：

- ```latex
D
``` 是域集合
- ```latex
\mathcal{I}
``` 是解释函数
- ```latex
\mathcal{E}
``` 是求值函数

**语义规则**:
$```latex
\mathcal{E}[\![\text{workflow}]\!] = \mathcal{I}(\text{workflow})
```$
$```latex
\mathcal{E}[\![\text{state}]\!] = \mathcal{I}(\text{state})
```$
$```latex
\mathcal{E}[\![\text{transition}]\!] = \mathcal{I}(\text{transition})
```$

## 2. 工作流DSL设计

### 2.1 核心语言特性

```go
// 工作流DSL解析器
type WorkflowDSLParser struct {
    lexer    *WorkflowLexer
    current  Token
    position int
    tokens   []Token
}

// 词法分析器
type WorkflowLexer struct {
    input string
    pos   int
    tokens []Token
}

type Token struct {
    Type    TokenType
    Value   string
    Line    int
    Column  int
}

type TokenType int

const (
    TOKEN_WORKFLOW TokenType = iota
    TOKEN_STATE
    TOKEN_EVENT
    TOKEN_TRANSITION
    TOKEN_IDENTIFIER
    TOKEN_STRING
    TOKEN_NUMBER
    TOKEN_LBRACE
    TOKEN_RBRACE
    TOKEN_COMMA
    TOKEN_COLON
    TOKEN_EOF
)

// 抽象语法树节点
type ASTNode interface {
    Accept(visitor ASTVisitor) interface{}
}

type WorkflowAST struct {
    Name       string
    States     []StateAST
    Events     []EventAST
    Transitions []TransitionAST
    Metadata   map[string]interface{}
}

type StateAST struct {
    Name   string
    Type   StateType
    Actions []ActionAST
}

type EventAST struct {
    Name    string
    Type    EventType
    Payload PayloadAST
}

type TransitionAST struct {
    From      string
    To        string
    Event     string
    Condition *ExpressionAST
    Action    *ActionAST
}
```

### 2.2 DSL语法解析

```go
// 递归下降解析器实现
func (p *WorkflowDSLParser) ParseWorkflow() (*WorkflowAST, error) {
    p.nextToken()
    
    if p.current.Type != TOKEN_WORKFLOW {
        return nil, fmt.Errorf("expected 'workflow', got %s", p.current.Value)
    }
    
    p.nextToken()
    if p.current.Type != TOKEN_IDENTIFIER {
        return nil, fmt.Errorf("expected workflow name, got %s", p.current.Value)
    }
    
    name := p.current.Value
    p.nextToken()
    
    if p.current.Type != TOKEN_LBRACE {
        return nil, fmt.Errorf("expected '{', got %s", p.current.Value)
    }
    
    workflow := &WorkflowAST{
        Name:       name,
        States:     []StateAST{},
        Events:     []EventAST{},
        Transitions: []TransitionAST{},
        Metadata:   make(map[string]interface{}),
    }
    
    p.nextToken()
    
    // 解析工作流体
    for p.current.Type != TOKEN_RBRACE && p.current.Type != TOKEN_EOF {
        switch p.current.Type {
        case TOKEN_STATE:
            state, err := p.parseState()
            if err != nil {
                return nil, err
            }
            workflow.States = append(workflow.States, *state)
        case TOKEN_EVENT:
            event, err := p.parseEvent()
            if err != nil {
                return nil, err
            }
            workflow.Events = append(workflow.Events, *event)
        case TOKEN_TRANSITION:
            transition, err := p.parseTransition()
            if err != nil {
                return nil, err
            }
            workflow.Transitions = append(workflow.Transitions, *transition)
        default:
            return nil, fmt.Errorf("unexpected token: %s", p.current.Value)
        }
    }
    
    if p.current.Type != TOKEN_RBRACE {
        return nil, fmt.Errorf("expected '}', got %s", p.current.Value)
    }
    
    return workflow, nil
}

func (p *WorkflowDSLParser) parseState() (*StateAST, error) {
    p.nextToken() // 跳过 'state'
    
    if p.current.Type != TOKEN_IDENTIFIER {
        return nil, fmt.Errorf("expected state name, got %s", p.current.Value)
    }
    
    name := p.current.Value
    p.nextToken()
    
    if p.current.Type != TOKEN_LBRACE {
        return nil, fmt.Errorf("expected '{', got %s", p.current.Value)
    }
    
    state := &StateAST{
        Name:    name,
        Actions: []ActionAST{},
    }
    
    p.nextToken()
    
    // 解析状态属性
    for p.current.Type != TOKEN_RBRACE {
        switch p.current.Type {
        case TOKEN_IDENTIFIER:
            if p.current.Value == "type" {
                p.nextToken() // 跳过 ':'
                if p.current.Type != TOKEN_IDENTIFIER {
                    return nil, fmt.Errorf("expected state type, got %s", p.current.Value)
                }
                state.Type = StateType(p.current.Value)
                p.nextToken()
            }
        default:
            return nil, fmt.Errorf("unexpected token in state: %s", p.current.Value)
        }
    }
    
    p.nextToken() // 跳过 '}'
    return state, nil
}
```

## 3. 工作流语言实现

### 3.1 语言运行时

```go
// 工作流语言运行时
type WorkflowRuntime struct {
    workflows map[string]*WorkflowDefinition
    executor  *WorkflowExecutor
    registry  *ComponentRegistry
}

// 工作流执行器
type WorkflowExecutor struct {
    runtime    *WorkflowRuntime
    scheduler  *TaskScheduler
    dispatcher *EventDispatcher
}

// 任务调度器
type TaskScheduler struct {
    tasks    chan Task
    workers  int
    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc
}

type Task struct {
    ID       string
    Workflow string
    State    string
    Action   ActionAST
    Data     map[string]interface{}
    Priority int
}

func (ts *TaskScheduler) Start() {
    ts.ctx, ts.cancel = context.WithCancel(context.Background())
    
    for i := 0; i < ts.workers; i++ {
        ts.wg.Add(1)
        go ts.worker()
    }
}

func (ts *TaskScheduler) worker() {
    defer ts.wg.Done()
    
    for {
        select {
        case task := <-ts.tasks:
            ts.executeTask(task)
        case <-ts.ctx.Done():
            return
        }
    }
}

func (ts *TaskScheduler) executeTask(task Task) {
    // 执行任务逻辑
    if task.Action != nil {
        if err := task.Action.Execute(task.Data); err != nil {
            // 处理错误
            log.Printf("Task %s failed: %v", task.ID, err)
        }
    }
}

// 事件分发器
type EventDispatcher struct {
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

type EventHandler func(event Event) error

func (ed *EventDispatcher) RegisterHandler(eventType string, handler EventHandler) {
    ed.mutex.Lock()
    defer ed.mutex.Unlock()
    
    if ed.handlers[eventType] == nil {
        ed.handlers[eventType] = []EventHandler{}
    }
    ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

func (ed *EventDispatcher) Dispatch(event Event) error {
    ed.mutex.RLock()
    handlers := ed.handlers[event.Type]
    ed.mutex.RUnlock()
    
    for _, handler := range handlers {
        if err := handler(event); err != nil {
            return err
        }
    }
    return nil
}
```

### 3.2 语言扩展机制

```go
// 工作流语言扩展接口
type WorkflowExtension interface {
    Name() string
    Version() string
    Initialize(ctx context.Context) error
    Execute(ctx context.Context, data map[string]interface{}) error
    Cleanup() error
}

// 扩展注册表
type ExtensionRegistry struct {
    extensions map[string]WorkflowExtension
    mutex      sync.RWMutex
}

func (er *ExtensionRegistry) Register(ext WorkflowExtension) error {
    er.mutex.Lock()
    defer er.mutex.Unlock()
    
    if er.extensions[ext.Name()] != nil {
        return fmt.Errorf("extension %s already registered", ext.Name())
    }
    
    er.extensions[ext.Name()] = ext
    return nil
}

func (er *ExtensionRegistry) Get(name string) (WorkflowExtension, bool) {
    er.mutex.RLock()
    defer er.mutex.RUnlock()
    
    ext, exists := er.extensions[name]
    return ext, exists
}

// 示例扩展：HTTP请求扩展
type HTTPRequestExtension struct {
    client *http.Client
}

func (h *HTTPRequestExtension) Name() string {
    return "http_request"
}

func (h *HTTPRequestExtension) Version() string {
    return "1.0.0"
}

func (h *HTTPRequestExtension) Initialize(ctx context.Context) error {
    h.client = &http.Client{
        Timeout: 30 * time.Second,
    }
    return nil
}

func (h *HTTPRequestExtension) Execute(ctx context.Context, data map[string]interface{}) error {
    url, ok := data["url"].(string)
    if !ok {
        return fmt.Errorf("url is required")
    }
    
    method, ok := data["method"].(string)
    if !ok {
        method = "GET"
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, nil)
    if err != nil {
        return err
    }
    
    resp, err := h.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // 处理响应
    data["status_code"] = resp.StatusCode
    data["response_headers"] = resp.Header
    
    return nil
}

func (h *HTTPRequestExtension) Cleanup() error {
    return nil
}
```

## 4. 工作流语言优化

### 4.1 语言性能优化

**定理 4.1** (工作流语言性能)
对于工作流语言 ```latex
L
```，其执行时间复杂度为 ```latex
O(n \cdot m \cdot k)
```，其中：

- ```latex
n
``` 是工作流状态数量
- ```latex
m
``` 是事件数量
- ```latex
k
``` 是平均转移数量

**优化策略**:

1. **缓存优化**: 缓存频繁访问的语法树节点
2. **并行执行**: 利用Go的goroutine实现并行处理
3. **内存池**: 重用对象减少GC压力

```go
// 工作流语言优化器
type WorkflowOptimizer struct {
    ast       *WorkflowAST
    optimized *OptimizedWorkflow
}

type OptimizedWorkflow struct {
    States     map[string]*OptimizedState
    Events     map[string]*OptimizedEvent
    Transitions map[string][]*OptimizedTransition
    Cache      map[string]interface{}
}

type OptimizedState struct {
    Name       string
    Type       StateType
    Actions    []ActionAST
    CacheKey   string
    IsCached   bool
}

func (wo *WorkflowOptimizer) Optimize() *OptimizedWorkflow {
    optimized := &OptimizedWorkflow{
        States:     make(map[string]*OptimizedState),
        Events:     make(map[string]*OptimizedEvent),
        Transitions: make(map[string][]*OptimizedTransition),
        Cache:      make(map[string]interface{}),
    }
    
    // 优化状态
    for _, state := range wo.ast.States {
        optState := &OptimizedState{
            Name:     state.Name,
            Type:     state.Type,
            Actions:  state.Actions,
            CacheKey: fmt.Sprintf("state_%s", state.Name),
        }
        optimized.States[state.Name] = optState
    }
    
    // 优化转移
    for _, transition := range wo.ast.Transitions {
        optTrans := &OptimizedTransition{
            From:      transition.From,
            To:        transition.To,
            Event:     transition.Event,
            Condition: transition.Condition,
            Action:    transition.Action,
        }
        
        if optimized.Transitions[transition.From] == nil {
            optimized.Transitions[transition.From] = []*OptimizedTransition{}
        }
        optimized.Transitions[transition.From] = append(
            optimized.Transitions[transition.From], optTrans)
    }
    
    return optimized
}
```

### 4.2 语言安全性

**定义 4.1** (工作流语言安全性)
工作流语言 ```latex
L
``` 是安全的，当且仅当：

1. 类型安全：所有表达式都有正确的类型
2. 状态安全：所有状态转移都是有效的
3. 资源安全：所有资源使用都是安全的

```go
// 工作流语言安全检查器
type WorkflowSafetyChecker struct {
    ast     *WorkflowAST
    errors  []SafetyError
    warnings []SafetyWarning
}

type SafetyError struct {
    Type    string
    Message string
    Line    int
    Column  int
}

type SafetyWarning struct {
    Type    string
    Message string
    Line    int
    Column  int
}

func (wsc *WorkflowSafetyChecker) Check() ([]SafetyError, []SafetyWarning) {
    wsc.errors = []SafetyError{}
    wsc.warnings = []SafetyWarning{}
    
    // 检查状态完整性
    wsc.checkStateCompleteness()
    
    // 检查转移有效性
    wsc.checkTransitionValidity()
    
    // 检查类型安全
    wsc.checkTypeSafety()
    
    // 检查资源安全
    wsc.checkResourceSafety()
    
    return wsc.errors, wsc.warnings
}

func (wsc *WorkflowSafetyChecker) checkStateCompleteness() {
    stateMap := make(map[string]bool)
    
    // 收集所有状态
    for _, state := range wsc.ast.States {
        stateMap[state.Name] = true
    }
    
    // 检查转移中的状态是否存在
    for _, transition := range wsc.ast.Transitions {
        if !stateMap[transition.From] {
            wsc.errors = append(wsc.errors, SafetyError{
                Type:    "MissingState",
                Message: fmt.Sprintf("State '%s' not defined", transition.From),
                Line:    0, // 需要从AST中获取位置信息
                Column:  0,
            })
        }
        
        if !stateMap[transition.To] {
            wsc.errors = append(wsc.errors, SafetyError{
                Type:    "MissingState",
                Message: fmt.Sprintf("State '%s' not defined", transition.To),
                Line:    0,
                Column:  0,
            })
        }
    }
}
```

## 5. 工作流语言应用

### 5.1 业务规则引擎

```go
// 业务规则引擎
type BusinessRuleEngine struct {
    rules     map[string]Rule
    workflow  *WorkflowDefinition
    executor  *RuleExecutor
}

type Rule struct {
    ID          string
    Name        string
    Condition   Expression
    Action      Action
    Priority    int
    IsActive    bool
}

type RuleExecutor struct {
    engine *BusinessRuleEngine
    cache  map[string]interface{}
}

func (re *RuleExecutor) ExecuteRules(context map[string]interface{}) error {
    // 按优先级排序规则
    rules := re.sortRulesByPriority()
    
    for _, rule := range rules {
        if !rule.IsActive {
            continue
        }
        
        // 评估条件
        if rule.Condition.Evaluate(context) {
            // 执行动作
            if err := rule.Action.Execute(context); err != nil {
                return fmt.Errorf("rule %s failed: %w", rule.ID, err)
            }
        }
    }
    
    return nil
}
```

### 5.2 工作流可视化

```go
// 工作流可视化生成器
type WorkflowVisualizer struct {
    workflow *WorkflowDefinition
    renderer Renderer
}

type Renderer interface {
    RenderDot(workflow *WorkflowDefinition) string
    RenderMermaid(workflow *WorkflowDefinition) string
    RenderPlantUML(workflow *WorkflowDefinition) string
}

type DotRenderer struct{}

func (dr *DotRenderer) RenderDot(workflow *WorkflowDefinition) string {
    var builder strings.Builder
    
    builder.WriteString("digraph workflow {\n")
    builder.WriteString("  rankdir=LR;\n")
    builder.WriteString("  node [shape=box];\n\n")
    
    // 添加状态节点
    for _, state := range workflow.States {
        builder.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", state.ID, state.Name))
    }
    
    builder.WriteString("\n")
    
    // 添加转移边
    for _, transition := range workflow.Transitions {
        builder.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n", 
            transition.From, transition.To, transition.Event))
    }
    
    builder.WriteString("}\n")
    return builder.String()
}
```

## 6. 总结

工作流语言为工作流系统提供了强大的表达能力，通过形式化的语法和语义定义，结合Go语言的高性能特性，实现了高效、安全、可扩展的工作流执行环境。

### 关键特性

1. **形式化定义**: 基于BNF语法的严格语言定义
2. **类型安全**: 完整的类型检查和验证
3. **高性能**: 优化的执行引擎和缓存机制
4. **可扩展**: 插件化的扩展机制
5. **可视化**: 多种格式的工作流可视化

### 应用场景

1. **业务流程自动化**: 企业级工作流管理
2. **微服务编排**: 服务间协调和编排
3. **数据处理管道**: 复杂数据处理流程
4. **规则引擎**: 业务规则管理和执行

---

**相关链接**:

- [01-工作流模型](./01-Workflow-Models.md)
- [03-工作流验证](./03-Workflow-Verification.md)
- [04-工作流优化](./04-Workflow-Optimization.md)
