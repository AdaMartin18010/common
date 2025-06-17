# 03-工作流基本术语 (Workflow Basic Terms)

## 目录

- [03-工作流基本术语 (Workflow Basic Terms)](#03-工作流基本术语-workflow-basic-terms)
  - [目录](#目录)
  - [1. 核心概念定义](#1-核心概念定义)
    - [1.1 工作流](#11-工作流)
    - [1.2 活动](#12-活动)
    - [1.3 转换](#13-转换)
    - [1.4 状态](#14-状态)
  - [2. 控制流概念](#2-控制流概念)
    - [2.1 顺序执行](#21-顺序执行)
    - [2.2 并行执行](#22-并行执行)
    - [2.3 条件分支](#23-条件分支)
    - [2.4 循环](#24-循环)
  - [3. 数据流概念](#3-数据流概念)
    - [3.1 数据传递](#31-数据传递)
    - [3.2 数据转换](#32-数据转换)
    - [3.3 数据聚合](#33-数据聚合)
  - [4. 异常处理概念](#4-异常处理概念)
    - [4.1 错误处理](#41-错误处理)
    - [4.2 补偿机制](#42-补偿机制)
    - [4.3 重试策略](#43-重试策略)
  - [5. 形式化定义](#5-形式化定义)
    - [5.1 数学符号](#51-数学符号)
    - [5.2 公理系统](#52-公理系统)
    - [5.3 推理规则](#53-推理规则)
  - [6. Go语言映射](#6-go语言映射)
    - [6.1 接口定义](#61-接口定义)
    - [6.2 类型系统](#62-类型系统)
    - [6.3 实现示例](#63-实现示例)
  - [总结](#总结)

---

## 1. 核心概念定义

### 1.1 工作流

**定义 1.1** (工作流): 工作流是一个有向图 $W = (N, E, \lambda)$，其中：

- $N$ 是节点集合，表示活动
- $E \subseteq N \times N$ 是边集合，表示控制流
- $\lambda: N \rightarrow \Sigma$ 是标签函数，将节点映射到活动类型

**定义 1.2** (工作流实例): 工作流实例是工作流的一个执行，表示为 $I = (W, s, \tau)$，其中：

- $W$ 是工作流定义
- $s: N \rightarrow \{active, completed, failed\}$ 是状态函数
- $\tau: N \rightarrow \mathbb{R}^+$ 是时间戳函数

### 1.2 活动

**定义 1.3** (活动): 活动是工作流中的基本执行单元，表示为 $A = (id, type, input, output, behavior)$，其中：

- $id$ 是唯一标识符
- $type \in \{task, subprocess, gateway\}$ 是活动类型
- $input$ 是输入数据模式
- $output$ 是输出数据模式
- $behavior$ 是执行行为描述

**活动类型分类**:

1. **任务活动 (Task)**: 原子执行单元
2. **子流程活动 (Subprocess)**: 包含其他工作流
3. **网关活动 (Gateway)**: 控制流决策点

### 1.3 转换

**定义 1.4** (转换): 转换是活动之间的连接，表示为 $T = (from, to, condition, action)$，其中：

- $from$ 是源活动
- $to$ 是目标活动
- $condition$ 是转换条件
- $action$ 是转换动作

**转换类型**:

1. **无条件转换**: 总是执行
2. **条件转换**: 基于条件执行
3. **默认转换**: 当其他条件都不满足时执行

### 1.4 状态

**定义 1.5** (活动状态): 活动状态是活动在某个时刻的执行状态，定义为：

$$\text{State} = \{ready, active, completed, failed, suspended, cancelled\}$$

**定义 1.6** (工作流状态): 工作流状态是所有活动状态的组合：

$$S_W = \prod_{A \in N} \text{State}_A$$

## 2. 控制流概念

### 2.1 顺序执行

**定义 2.1** (顺序执行): 活动 $A_1, A_2, \ldots, A_n$ 的顺序执行定义为：

$$\text{Sequence}(A_1, A_2, \ldots, A_n) = A_1 \circ A_2 \circ \cdots \circ A_n$$

**公理 2.1** (结合律): $(A_1 \circ A_2) \circ A_3 = A_1 \circ (A_2 \circ A_3)$

### 2.2 并行执行

**定义 2.2** (并行执行): 活动 $A_1, A_2, \ldots, A_n$ 的并行执行定义为：

$$\text{Parallel}(A_1, A_2, \ldots, A_n) = A_1 \parallel A_2 \parallel \cdots \parallel A_n$$

**公理 2.2** (交换律): $A_1 \parallel A_2 = A_2 \parallel A_1$

**公理 2.3** (结合律): $(A_1 \parallel A_2) \parallel A_3 = A_1 \parallel (A_2 \parallel A_3)$

### 2.3 条件分支

**定义 2.3** (条件分支): 基于条件 $c$ 在活动 $A_1$ 和 $A_2$ 之间选择：

$$\text{Choice}(c, A_1, A_2) = \text{if } c \text{ then } A_1 \text{ else } A_2$$

**公理 2.4** (选择幂等): $\text{Choice}(c, A, A) = A$

### 2.4 循环

**定义 2.4** (循环): 活动 $A$ 在条件 $c$ 下重复执行：

$$\text{Loop}(c, A) = \text{while } c \text{ do } A$$

**定理 2.1**: 循环可以表示为递归：

$$\text{Loop}(c, A) = \text{Choice}(c, A \circ \text{Loop}(c, A), \text{skip})$$

## 3. 数据流概念

### 3.1 数据传递

**定义 3.1** (数据传递): 从活动 $A_1$ 到活动 $A_2$ 的数据传递定义为：

$$\text{DataFlow}(A_1, A_2, \sigma) = \{(d_1, d_2) \mid d_2 = \sigma(d_1)\}$$

其中 $\sigma$ 是数据转换函数。

### 3.2 数据转换

**定义 3.2** (数据转换): 数据转换函数 $\sigma: D_1 \rightarrow D_2$ 满足：

$$\forall d_1 \in D_1, \exists d_2 \in D_2: \sigma(d_1) = d_2$$

### 3.3 数据聚合

**定义 3.3** (数据聚合): 多个数据源的聚合定义为：

$$\text{Aggregate}(\{d_1, d_2, \ldots, d_n\}, \oplus) = d_1 \oplus d_2 \oplus \cdots \oplus d_n$$

其中 $\oplus$ 是聚合操作符。

## 4. 异常处理概念

### 4.1 错误处理

**定义 4.1** (错误处理): 活动 $A$ 的错误处理定义为：

$$\text{ErrorHandler}(A, E) = A \oplus E$$

其中 $E$ 是错误处理活动，$\oplus$ 表示错误处理组合。

### 4.2 补偿机制

**定义 4.2** (补偿): 活动 $A$ 的补偿活动 $\bar{A}$ 满足：

$$A \circ \bar{A} \sim \text{skip}$$

其中 $\sim$ 表示语义等价。

### 4.3 重试策略

**定义 4.3** (重试): 活动 $A$ 的重试策略定义为：

$$\text{Retry}(A, n, \delta) = \text{repeat } A \text{ up to } n \text{ times with delay } \delta$$

## 5. 形式化定义

### 5.1 数学符号

| 符号 | 含义 | 定义 |
|------|------|------|
| $W$ | 工作流 | 有向图 $(N, E, \lambda)$ |
| $A$ | 活动 | 执行单元 |
| $T$ | 转换 | 活动间连接 |
| $S$ | 状态 | 执行状态 |
| $\circ$ | 顺序组合 | 活动顺序执行 |
| $\parallel$ | 并行组合 | 活动并行执行 |
| $\oplus$ | 错误处理组合 | 错误处理 |
| $\sim$ | 语义等价 | 行为等价 |

### 5.2 公理系统

**公理系统 $\mathcal{A}$**:

1. **结合律**: $(A_1 \circ A_2) \circ A_3 = A_1 \circ (A_2 \circ A_3)$
2. **交换律**: $A_1 \parallel A_2 = A_2 \parallel A_1$
3. **分配律**: $A_1 \circ (A_2 \parallel A_3) = (A_1 \circ A_2) \parallel (A_1 \circ A_3)$
4. **单位元**: $A \circ \text{skip} = \text{skip} \circ A = A$
5. **幂等律**: $A \parallel A = A$

### 5.3 推理规则

**推理规则 $\mathcal{R}$**:

1. **替换规则**: 如果 $A_1 = A_2$，则 $C[A_1] = C[A_2]$
2. **上下文规则**: 如果 $A_1 \sim A_2$，则 $C[A_1] \sim C[A_2]$
3. **组合规则**: 如果 $A_1 \sim A_2$ 且 $B_1 \sim B_2$，则 $A_1 \circ B_1 \sim A_2 \circ B_2$

## 6. Go语言映射

### 6.1 接口定义

```go
// Workflow 工作流接口
type Workflow interface {
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    GetState() WorkflowState
    GetActivities() []Activity
    GetTransitions() []Transition
}

// Activity 活动接口
type Activity interface {
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    GetID() string
    GetType() ActivityType
    GetState() ActivityState
}

// Transition 转换接口
type Transition interface {
    GetFrom() string
    GetTo() string
    GetCondition() func(interface{}) bool
    GetAction() func(interface{}) interface{}
}

// WorkflowState 工作流状态
type WorkflowState struct {
    ID        string
    Status    WorkflowStatus
    Activities map[string]ActivityState
    Data      map[string]interface{}
    Timestamp time.Time
}

// ActivityState 活动状态
type ActivityState struct {
    ID        string
    Status    ActivityStatus
    Input     interface{}
    Output    interface{}
    Error     error
    StartTime time.Time
    EndTime   time.Time
}

// 枚举类型
type WorkflowStatus int
type ActivityStatus int
type ActivityType int

const (
    WorkflowReady WorkflowStatus = iota
    WorkflowRunning
    WorkflowCompleted
    WorkflowFailed
    WorkflowSuspended
    WorkflowCancelled
)

const (
    ActivityReady ActivityStatus = iota
    ActivityRunning
    ActivityCompleted
    ActivityFailed
    ActivitySuspended
    ActivityCancelled
)

const (
    ActivityTask ActivityType = iota
    ActivitySubprocess
    ActivityGateway
)
```

### 6.2 类型系统

```go
// WorkflowBuilder 工作流构建器
type WorkflowBuilder struct {
    activities  map[string]Activity
    transitions []Transition
    startNode   string
    endNodes    []string
}

// NewWorkflowBuilder 创建新的工作流构建器
func NewWorkflowBuilder() *WorkflowBuilder {
    return &WorkflowBuilder{
        activities:  make(map[string]Activity),
        transitions: make([]Transition, 0),
        endNodes:    make([]string, 0),
    }
}

// AddActivity 添加活动
func (wb *WorkflowBuilder) AddActivity(activity Activity) *WorkflowBuilder {
    wb.activities[activity.GetID()] = activity
    return wb
}

// AddTransition 添加转换
func (wb *WorkflowBuilder) AddTransition(from, to string, condition func(interface{}) bool) *WorkflowBuilder {
    transition := &transition{
        from:     from,
        to:       to,
        condition: condition,
    }
    wb.transitions = append(wb.transitions, transition)
    return wb
}

// SetStartNode 设置开始节点
func (wb *WorkflowBuilder) SetStartNode(nodeID string) *WorkflowBuilder {
    wb.startNode = nodeID
    return wb
}

// AddEndNode 添加结束节点
func (wb *WorkflowBuilder) AddEndNode(nodeID string) *WorkflowBuilder {
    wb.endNodes = append(wb.endNodes, nodeID)
    return wb
}

// Build 构建工作流
func (wb *WorkflowBuilder) Build() (Workflow, error) {
    if wb.startNode == "" {
        return nil, fmt.Errorf("start node not set")
    }
    
    if len(wb.endNodes) == 0 {
        return nil, fmt.Errorf("no end nodes defined")
    }
    
    return &workflow{
        activities:  wb.activities,
        transitions: wb.transitions,
        startNode:   wb.startNode,
        endNodes:    wb.endNodes,
        state:       WorkflowState{Status: WorkflowReady},
    }, nil
}
```

### 6.3 实现示例

```go
// workflow 工作流实现
type workflow struct {
    activities  map[string]Activity
    transitions []Transition
    startNode   string
    endNodes    []string
    state       WorkflowState
    mu          sync.RWMutex
}

// Execute 执行工作流
func (w *workflow) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    w.mu.Lock()
    w.state.Status = WorkflowRunning
    w.state.Data = make(map[string]interface{})
    w.state.Data["input"] = input
    w.mu.Unlock()
    
    defer func() {
        w.mu.Lock()
        if w.state.Status == WorkflowRunning {
            w.state.Status = WorkflowCompleted
        }
        w.mu.Unlock()
    }()
    
    // 初始化活动状态
    for id, activity := range w.activities {
        w.state.Activities[id] = ActivityState{
            ID:     id,
            Status: ActivityReady,
        }
    }
    
    // 执行工作流
    return w.executeNode(ctx, w.startNode, input)
}

// executeNode 执行节点
func (w *workflow) executeNode(ctx context.Context, nodeID string, input interface{}) (interface{}, error) {
    activity, exists := w.activities[nodeID]
    if !exists {
        return nil, fmt.Errorf("activity not found: %s", nodeID)
    }
    
    // 更新活动状态
    w.mu.Lock()
    state := w.state.Activities[nodeID]
    state.Status = ActivityRunning
    state.Input = input
    state.StartTime = time.Now()
    w.state.Activities[nodeID] = state
    w.mu.Unlock()
    
    // 执行活动
    output, err := activity.Execute(ctx, input)
    
    // 更新活动状态
    w.mu.Lock()
    state = w.state.Activities[nodeID]
    state.Output = output
    state.EndTime = time.Now()
    if err != nil {
        state.Status = ActivityFailed
        state.Error = err
        w.state.Status = WorkflowFailed
    } else {
        state.Status = ActivityCompleted
    }
    w.state.Activities[nodeID] = state
    w.mu.Unlock()
    
    if err != nil {
        return nil, err
    }
    
    // 检查是否为结束节点
    for _, endNode := range w.endNodes {
        if nodeID == endNode {
            return output, nil
        }
    }
    
    // 查找下一个节点
    nextNode, err := w.findNextNode(nodeID, output)
    if err != nil {
        return nil, err
    }
    
    return w.executeNode(ctx, nextNode, output)
}

// findNextNode 查找下一个节点
func (w *workflow) findNextNode(currentNode string, data interface{}) (string, error) {
    for _, transition := range w.transitions {
        if transition.GetFrom() == currentNode {
            if transition.GetCondition() == nil || transition.GetCondition()(data) {
                return transition.GetTo(), nil
            }
        }
    }
    return "", fmt.Errorf("no valid transition from node: %s", currentNode)
}

// GetState 获取工作流状态
func (w *workflow) GetState() WorkflowState {
    w.mu.RLock()
    defer w.mu.RUnlock()
    return w.state
}

// GetActivities 获取活动列表
func (w *workflow) GetActivities() []Activity {
    activities := make([]Activity, 0, len(w.activities))
    for _, activity := range w.activities {
        activities = append(activities, activity)
    }
    return activities
}

// GetTransitions 获取转换列表
func (w *workflow) GetTransitions() []Transition {
    return w.transitions
}

// transition 转换实现
type transition struct {
    from     string
    to       string
    condition func(interface{}) bool
    action    func(interface{}) interface{}
}

func (t *transition) GetFrom() string { return t.from }
func (t *transition) GetTo() string { return t.to }
func (t *transition) GetCondition() func(interface{}) bool { return t.condition }
func (t *transition) GetAction() func(interface{}) interface{} { return t.action }
```

## 总结

本文档定义了工作流系统的基本术语和概念，包括：

1. **核心概念**: 工作流、活动、转换、状态
2. **控制流概念**: 顺序、并行、条件分支、循环
3. **数据流概念**: 数据传递、转换、聚合
4. **异常处理**: 错误处理、补偿、重试
5. **形式化定义**: 数学符号、公理系统、推理规则
6. **Go语言映射**: 接口定义、类型系统、实现示例

这些术语为工作流系统的设计和实现提供了统一的概念基础。
