# 01-Petri网模型 (Petri Net Model)

## 目录

- [01-Petri网模型 (Petri Net Model)](#01-petri网模型-petri-net-model)
  - [目录](#目录)
  - [1. Petri网基础](#1-petri网基础)
    - [1.1 基本定义](#11-基本定义)
    - [1.2 标记和变迁](#12-标记和变迁)
    - [1.3 可达性](#13-可达性)
  - [2. 工作流Petri网](#2-工作流petri网)
    - [2.1 工作流映射](#21-工作流映射)
    - [2.2 结构性质](#22-结构性质)
    - [2.3 行为性质](#23-行为性质)
  - [3. 高级Petri网](#3-高级petri网)
    - [3.1 时间Petri网](#31-时间petri网)
    - [3.2 颜色Petri网](#32-颜色petri网)
    - [3.3 层次Petri网](#33-层次petri网)
  - [4. 形式化分析](#4-形式化分析)
    - [4.1 可达性分析](#41-可达性分析)
    - [4.2 活性分析](#42-活性分析)
    - [4.3 有界性分析](#43-有界性分析)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 基本Petri网](#51-基本petri网)
    - [5.2 工作流Petri网](#52-工作流petri网)
    - [5.3 分析算法](#53-分析算法)
  - [6. 应用案例](#6-应用案例)
    - [6.1 业务流程建模](#61-业务流程建模)
    - [6.2 并发系统分析](#62-并发系统分析)
    - [6.3 死锁检测](#63-死锁检测)
  - [总结](#总结)

---

## 1. Petri网基础

### 1.1 基本定义

**定义 1.1** (Petri网): Petri网是一个四元组 $N = (P, T, F, M_0)$，其中：

- $P$ 是库所(places)的有限集合
- $T$ 是变迁(transitions)的有限集合，且 $P \cap T = \emptyset$
- $F \subseteq (P \times T) \cup (T \times P)$ 是流关系
- $M_0: P \rightarrow \mathbb{N}$ 是初始标记

**定义 1.2** (前集和后集): 对于 $x \in P \cup T$：

- 前集：$^\bullet x = \{y \mid (y, x) \in F\}$
- 后集：$x^\bullet = \{y \mid (x, y) \in F\}$

**定义 1.3** (标记): 标记 $M: P \rightarrow \mathbb{N}$ 为每个库所分配一个非负整数，表示该库所中的令牌数量。

### 1.2 标记和变迁

**定义 1.4** (变迁使能): 变迁 $t \in T$ 在标记 $M$ 下使能，当且仅当：

$$\forall p \in ^\bullet t: M(p) \geq F(p, t)$$

**定义 1.5** (变迁发生): 如果变迁 $t$ 在标记 $M$ 下使能，则它可以发生，产生新标记 $M'$：

$$M'(p) = M(p) - F(p, t) + F(t, p)$$

其中 $F(p, t)$ 表示从库所 $p$ 到变迁 $t$ 的弧的权重。

**定理 1.1** (变迁发生唯一性): 对于给定的标记 $M$ 和使能的变迁 $t$，新标记 $M'$ 是唯一确定的。

**证明**: 根据定义 1.5，$M'$ 的计算是确定性的，因此唯一性成立。

### 1.3 可达性

**定义 1.6** (可达性): 标记 $M'$ 从标记 $M$ 可达，如果存在变迁序列 $\sigma = t_1 t_2 \ldots t_n$ 使得：

$$M[t_1\rangle M_1[t_2\rangle M_2 \ldots [t_n\rangle M'$$

**定义 1.7** (可达集): 从初始标记 $M_0$ 可达的所有标记的集合称为可达集 $R(N, M_0)$。

**定理 1.2** (可达性传递性): 如果 $M_1$ 从 $M_0$ 可达，$M_2$ 从 $M_1$ 可达，则 $M_2$ 从 $M_0$ 可达。

**证明**: 通过连接两个变迁序列，可以构造从 $M_0$ 到 $M_2$ 的变迁序列。

## 2. 工作流Petri网

### 2.1 工作流映射

**定义 2.1** (工作流Petri网): 工作流Petri网 $WPN = (P, T, F, M_0, \lambda)$ 是一个扩展的Petri网，其中：

- $(P, T, F, M_0)$ 是基本Petri网
- $\lambda: T \rightarrow \Sigma$ 是标签函数，将变迁映射到活动类型

**映射规则 2.1** (工作流到Petri网):

1. **活动映射**: 每个活动 $A$ 映射为变迁 $t_A$
2. **状态映射**: 每个状态 $S$ 映射为库所 $p_S$
3. **转换映射**: 每个转换 $(S_1, A, S_2)$ 映射为弧 $(p_{S_1}, t_A)$ 和 $(t_A, p_{S_2})$
4. **初始状态**: 初始状态映射为初始标记

**示例 2.1**: 简单顺序工作流的Petri网表示：

```go
// 顺序工作流: A -> B -> C
// Petri网表示:
// p_start --t_A--> p_1 --t_B--> p_2 --t_C--> p_end
```

### 2.2 结构性质

**定义 2.2** (结构有界性): Petri网 $N$ 是结构有界的，如果对于任意初始标记 $M_0$，网都是有界的。

**定义 2.3** (结构活性): Petri网 $N$ 是结构活的，如果存在初始标记 $M_0$ 使得网是活的。

**定义 2.4** (可重复性): Petri网 $N$ 是可重复的，如果存在变迁序列 $\sigma$ 使得 $M_0[\sigma\rangle M_0$。

**定理 2.1** (工作流Petri网性质): 工作流Petri网满足以下性质：

1. **有界性**: 每个库所的令牌数量有上界
2. **活性**: 每个变迁最终都可以发生
3. **可重复性**: 可以从初始状态回到初始状态

**证明**:

- 有界性：工作流有有限状态，因此有界
- 活性：工作流设计确保所有活动都能执行
- 可重复性：工作流可以重新开始

### 2.3 行为性质

**定义 2.5** (安全性): 工作流Petri网是安全的，如果：

$$\forall M \in R(N, M_0), \forall p \in P: M(p) \leq 1$$

**定义 2.6** (活性): 工作流Petri网是活的，如果：

$$\forall t \in T, \forall M \in R(N, M_0), \exists M' \in R(N, M): M'[t\rangle$$

**定义 2.7** (公平性): 工作流Petri网是公平的，如果：

$$\forall t_1, t_2 \in T: \text{if } t_1 \text{ and } t_2 \text{ are enabled infinitely often, then both occur infinitely often}$$

## 3. 高级Petri网

### 3.1 时间Petri网

**定义 3.1** (时间Petri网): 时间Petri网 $TPN = (P, T, F, M_0, I)$ 是一个扩展的Petri网，其中：

- $(P, T, F, M_0)$ 是基本Petri网
- $I: T \rightarrow \mathbb{R}^+ \times \mathbb{R}^+$ 是时间间隔函数

**定义 3.2** (时间变迁发生): 变迁 $t$ 在时间 $\tau$ 发生，如果：

$$\tau \in I(t) \text{ and } t \text{ is enabled}$$

**定理 3.1** (时间可达性): 时间Petri网的可达性问题是PSPACE完全的。

### 3.2 颜色Petri网

**定义 3.3** (颜色Petri网): 颜色Petri网 $CPN = (P, T, F, M_0, C, V)$ 是一个扩展的Petri网，其中：

- $(P, T, F, M_0)$ 是基本Petri网
- $C: P \cup T \rightarrow \Sigma$ 是颜色函数
- $V: F \rightarrow \text{Expr}$ 是变量函数

**定义 3.4** (颜色标记): 颜色标记为每个库所分配带颜色的令牌。

### 3.3 层次Petri网

**定义 3.5** (层次Petri网): 层次Petri网允许变迁包含子网，形成层次结构。

## 4. 形式化分析

### 4.1 可达性分析

**算法 4.1** (可达性分析): 使用状态空间搜索分析可达性：

```go
func ReachabilityAnalysis(petriNet PetriNet) map[string]bool {
    reachable := make(map[string]bool)
    queue := []Marking{petriNet.InitialMarking}
    reachable[petriNet.InitialMarking.String()] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for _, transition := range petriNet.Transitions {
            if current.IsEnabled(transition) {
                next := current.Fire(transition)
                nextStr := next.String()
                
                if !reachable[nextStr] {
                    reachable[nextStr] = true
                    queue = append(queue, next)
                }
            }
        }
    }
    
    return reachable
}
```

**定理 4.1** (可达性复杂性): Petri网的可达性问题在EXPSPACE中。

### 4.2 活性分析

**算法 4.2** (活性分析): 检查每个变迁是否最终可以发生：

```go
func LivenessAnalysis(petriNet PetriNet) map[string]bool {
    liveness := make(map[string]bool)
    
    for _, transition := range petriNet.Transitions {
        liveness[transition.ID] = isTransitionLive(petriNet, transition)
    }
    
    return liveness
}

func isTransitionLive(petriNet PetriNet, transition Transition) bool {
    // 使用可达性分析检查变迁是否可以在任意可达标记中发生
    reachable := ReachabilityAnalysis(petriNet)
    
    for markingStr := range reachable {
        marking := parseMarking(markingStr)
        if marking.IsEnabled(transition) {
            return true
        }
    }
    
    return false
}
```

### 4.3 有界性分析

**算法 4.3** (有界性分析): 检查库所的令牌数量是否有上界：

```go
func BoundednessAnalysis(petriNet PetriNet) map[string]int {
    boundedness := make(map[string]int)
    
    for _, place := range petriNet.Places {
        boundedness[place.ID] = getPlaceBound(petriNet, place)
    }
    
    return boundedness
}

func getPlaceBound(petriNet PetriNet, place Place) int {
    reachable := ReachabilityAnalysis(petriNet)
    maxTokens := 0
    
    for markingStr := range reachable {
        marking := parseMarking(markingStr)
        tokens := marking.Tokens[place.ID]
        if tokens > maxTokens {
            maxTokens = tokens
        }
    }
    
    return maxTokens
}
```

## 5. Go语言实现

### 5.1 基本Petri网

```go
// PetriNet 基本Petri网
type PetriNet struct {
    Places         map[string]Place
    Transitions    map[string]Transition
    Flow           map[string]int // (place,transition) -> weight
    InitialMarking Marking
}

// Place 库所
type Place struct {
    ID       string
    Name     string
    Capacity int // -1表示无限制
}

// Transition 变迁
type Transition struct {
    ID       string
    Name     string
    Guard    func(Marking) bool // 守卫条件
    Action   func(Marking) Marking // 动作
}

// Marking 标记
type Marking struct {
    Tokens map[string]int // place -> token count
}

// NewPetriNet 创建新的Petri网
func NewPetriNet() *PetriNet {
    return &PetriNet{
        Places:      make(map[string]Place),
        Transitions: make(map[string]Transition),
        Flow:        make(map[string]int),
        InitialMarking: Marking{
            Tokens: make(map[string]int),
        },
    }
}

// AddPlace 添加库所
func (pn *PetriNet) AddPlace(id, name string, capacity int) {
    pn.Places[id] = Place{
        ID:       id,
        Name:     name,
        Capacity: capacity,
    }
}

// AddTransition 添加变迁
func (pn *PetriNet) AddTransition(id, name string, guard func(Marking) bool, action func(Marking) Marking) {
    pn.Transitions[id] = Transition{
        ID:     id,
        Name:   name,
        Guard:  guard,
        Action: action,
    }
}

// AddFlow 添加流关系
func (pn *PetriNet) AddFlow(fromPlace, toTransition string, weight int) {
    key := fmt.Sprintf("%s->%s", fromPlace, toTransition)
    pn.Flow[key] = weight
}

// AddFlowFromTransition 添加从变迁到库所的流
func (pn *PetriNet) AddFlowFromTransition(fromTransition, toPlace string, weight int) {
    key := fmt.Sprintf("%s->%s", fromTransition, toPlace)
    pn.Flow[key] = weight
}

// IsEnabled 检查变迁是否使能
func (pn *PetriNet) IsEnabled(marking Marking, transition Transition) bool {
    // 检查守卫条件
    if transition.Guard != nil && !transition.Guard(marking) {
        return false
    }
    
    // 检查输入库所的令牌数量
    for placeID := range pn.Places {
        key := fmt.Sprintf("%s->%s", placeID, transition.ID)
        if weight, exists := pn.Flow[key]; exists {
            if marking.Tokens[placeID] < weight {
                return false
            }
        }
    }
    
    return true
}

// Fire 执行变迁
func (pn *PetriNet) Fire(marking Marking, transition Transition) (Marking, error) {
    if !pn.IsEnabled(marking, transition) {
        return marking, fmt.Errorf("transition %s is not enabled", transition.ID)
    }
    
    // 创建新标记
    newMarking := Marking{
        Tokens: make(map[string]int),
    }
    
    // 复制当前标记
    for placeID, tokens := range marking.Tokens {
        newMarking.Tokens[placeID] = tokens
    }
    
    // 移除输入令牌
    for placeID := range pn.Places {
        key := fmt.Sprintf("%s->%s", placeID, transition.ID)
        if weight, exists := pn.Flow[key]; exists {
            newMarking.Tokens[placeID] -= weight
        }
    }
    
    // 添加输出令牌
    for placeID := range pn.Places {
        key := fmt.Sprintf("%s->%s", transition.ID, placeID)
        if weight, exists := pn.Flow[key]; exists {
            newMarking.Tokens[placeID] += weight
        }
    }
    
    // 执行动作
    if transition.Action != nil {
        newMarking = transition.Action(newMarking)
    }
    
    return newMarking, nil
}

// String 标记的字符串表示
func (m Marking) String() string {
    tokens := make([]string, 0)
    for placeID, count := range m.Tokens {
        tokens = append(tokens, fmt.Sprintf("%s:%d", placeID, count))
    }
    return strings.Join(tokens, ",")
}
```

### 5.2 工作流Petri网

```go
// WorkflowPetriNet 工作流Petri网
type WorkflowPetriNet struct {
    *PetriNet
    Labels map[string]string // transition -> activity label
}

// NewWorkflowPetriNet 创建新的工作流Petri网
func NewWorkflowPetriNet() *WorkflowPetriNet {
    return &WorkflowPetriNet{
        PetriNet: NewPetriNet(),
        Labels:   make(map[string]string),
    }
}

// AddActivity 添加活动
func (wpn *WorkflowPetriNet) AddActivity(activityID, activityName string) {
    // 创建前状态库所
    prePlaceID := fmt.Sprintf("pre_%s", activityID)
    wpn.AddPlace(prePlaceID, fmt.Sprintf("Pre-%s", activityName), 1)
    
    // 创建后状态库所
    postPlaceID := fmt.Sprintf("post_%s", activityID)
    wpn.AddPlace(postPlaceID, fmt.Sprintf("Post-%s", activityName), 1)
    
    // 创建变迁
    wpn.AddTransition(activityID, activityName, nil, nil)
    
    // 添加流关系
    wpn.AddFlow(prePlaceID, activityID, 1)
    wpn.AddFlowFromTransition(activityID, postPlaceID, 1)
    
    // 添加标签
    wpn.Labels[activityID] = activityName
}

// AddSequence 添加顺序关系
func (wpn *WorkflowPetriNet) AddSequence(fromActivity, toActivity string) {
    fromPostPlace := fmt.Sprintf("post_%s", fromActivity)
    toPrePlace := fmt.Sprintf("pre_%s", toActivity)
    
    // 添加连接库所
    connectionPlace := fmt.Sprintf("conn_%s_%s", fromActivity, toActivity)
    wpn.AddPlace(connectionPlace, fmt.Sprintf("Connection-%s-%s", fromActivity, toActivity), 1)
    
    // 添加流关系
    wpn.AddFlowFromTransition(fromActivity, connectionPlace, 1)
    wpn.AddFlow(connectionPlace, toActivity, 1)
}

// AddParallel 添加并行关系
func (wpn *WorkflowPetriNet) AddParallel(activities []string) {
    // 创建分叉库所
    forkPlace := fmt.Sprintf("fork_%s", strings.Join(activities, "_"))
    wpn.AddPlace(forkPlace, "Fork", 1)
    
    // 创建合并库所
    joinPlace := fmt.Sprintf("join_%s", strings.Join(activities, "_"))
    wpn.AddPlace(joinPlace, "Join", 1)
    
    // 添加流关系
    for _, activity := range activities {
        wpn.AddFlow(forkPlace, activity, 1)
        wpn.AddFlowFromTransition(activity, joinPlace, 1)
    }
}

// AddCondition 添加条件分支
func (wpn *WorkflowPetriNet) AddCondition(condition string, trueActivity, falseActivity string) {
    // 创建条件库所
    conditionPlace := fmt.Sprintf("cond_%s", condition)
    wpn.AddPlace(conditionPlace, fmt.Sprintf("Condition-%s", condition), 1)
    
    // 添加条件变迁
    trueTransitionID := fmt.Sprintf("true_%s", condition)
    falseTransitionID := fmt.Sprintf("false_%s", condition)
    
    wpn.AddTransition(trueTransitionID, fmt.Sprintf("True-%s", condition), nil, nil)
    wpn.AddTransition(falseTransitionID, fmt.Sprintf("False-%s", condition), nil, nil)
    
    // 添加流关系
    wpn.AddFlow(conditionPlace, trueTransitionID, 1)
    wpn.AddFlow(conditionPlace, falseTransitionID, 1)
    wpn.AddFlowFromTransition(trueTransitionID, fmt.Sprintf("pre_%s", trueActivity), 1)
    wpn.AddFlowFromTransition(falseTransitionID, fmt.Sprintf("pre_%s", falseActivity), 1)
}

// SetInitialActivity 设置初始活动
func (wpn *WorkflowPetriNet) SetInitialActivity(activityID string) {
    prePlaceID := fmt.Sprintf("pre_%s", activityID)
    wpn.InitialMarking.Tokens[prePlaceID] = 1
}

// SetFinalActivity 设置最终活动
func (wpn *WorkflowPetriNet) SetFinalActivity(activityID string) {
    postPlaceID := fmt.Sprintf("post_%s", activityID)
    // 最终活动不需要特殊处理，只需要确保没有后继活动
}
```

### 5.3 分析算法

```go
// PetriNetAnalyzer Petri网分析器
type PetriNetAnalyzer struct {
    petriNet *PetriNet
}

// NewPetriNetAnalyzer 创建新的分析器
func NewPetriNetAnalyzer(petriNet *PetriNet) *PetriNetAnalyzer {
    return &PetriNetAnalyzer{
        petriNet: petriNet,
    }
}

// AnalyzeReachability 分析可达性
func (pna *PetriNetAnalyzer) AnalyzeReachability() map[string]bool {
    reachable := make(map[string]bool)
    queue := []Marking{pna.petriNet.InitialMarking}
    reachable[pna.petriNet.InitialMarking.String()] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for _, transition := range pna.petriNet.Transitions {
            if pna.petriNet.IsEnabled(current, transition) {
                next, err := pna.petriNet.Fire(current, transition)
                if err != nil {
                    continue
                }
                
                nextStr := next.String()
                if !reachable[nextStr] {
                    reachable[nextStr] = true
                    queue = append(queue, next)
                }
            }
        }
    }
    
    return reachable
}

// AnalyzeLiveness 分析活性
func (pna *PetriNetAnalyzer) AnalyzeLiveness() map[string]bool {
    liveness := make(map[string]bool)
    reachable := pna.AnalyzeReachability()
    
    for transitionID := range pna.petriNet.Transitions {
        liveness[transitionID] = false
        
        for markingStr := range reachable {
            marking := pna.parseMarking(markingStr)
            if pna.petriNet.IsEnabled(marking, pna.petriNet.Transitions[transitionID]) {
                liveness[transitionID] = true
                break
            }
        }
    }
    
    return liveness
}

// AnalyzeBoundedness 分析有界性
func (pna *PetriNetAnalyzer) AnalyzeBoundedness() map[string]int {
    boundedness := make(map[string]int)
    reachable := pna.AnalyzeReachability()
    
    for placeID := range pna.petriNet.Places {
        maxTokens := 0
        
        for markingStr := range reachable {
            marking := pna.parseMarking(markingStr)
            tokens := marking.Tokens[placeID]
            if tokens > maxTokens {
                maxTokens = tokens
            }
        }
        
        boundedness[placeID] = maxTokens
    }
    
    return boundedness
}

// AnalyzeDeadlock 分析死锁
func (pna *PetriNetAnalyzer) AnalyzeDeadlock() []Marking {
    deadlocks := make([]Marking, 0)
    reachable := pna.AnalyzeReachability()
    
    for markingStr := range reachable {
        marking := pna.parseMarking(markingStr)
        isDeadlock := true
        
        for _, transition := range pna.petriNet.Transitions {
            if pna.petriNet.IsEnabled(marking, transition) {
                isDeadlock = false
                break
            }
        }
        
        if isDeadlock {
            deadlocks = append(deadlocks, marking)
        }
    }
    
    return deadlocks
}

// parseMarking 解析标记字符串
func (pna *PetriNetAnalyzer) parseMarking(markingStr string) Marking {
    marking := Marking{
        Tokens: make(map[string]int),
    }
    
    if markingStr == "" {
        return marking
    }
    
    parts := strings.Split(markingStr, ",")
    for _, part := range parts {
        if strings.Contains(part, ":") {
            placeTokens := strings.Split(part, ":")
            if len(placeTokens) == 2 {
                placeID := placeTokens[0]
                tokens, _ := strconv.Atoi(placeTokens[1])
                marking.Tokens[placeID] = tokens
            }
        }
    }
    
    return marking
}
```

## 6. 应用案例

### 6.1 业务流程建模

```go
// 示例：订单处理流程
func createOrderProcess() *WorkflowPetriNet {
    wpn := NewWorkflowPetriNet()
    
    // 添加活动
    wpn.AddActivity("receive_order", "接收订单")
    wpn.AddActivity("validate_order", "验证订单")
    wpn.AddActivity("check_inventory", "检查库存")
    wpn.AddActivity("process_payment", "处理支付")
    wpn.AddActivity("ship_order", "发货")
    wpn.AddActivity("send_notification", "发送通知")
    
    // 添加顺序关系
    wpn.AddSequence("receive_order", "validate_order")
    wpn.AddSequence("validate_order", "check_inventory")
    wpn.AddSequence("check_inventory", "process_payment")
    wpn.AddSequence("process_payment", "ship_order")
    wpn.AddSequence("ship_order", "send_notification")
    
    // 设置初始和最终活动
    wpn.SetInitialActivity("receive_order")
    wpn.SetFinalActivity("send_notification")
    
    return wpn
}
```

### 6.2 并发系统分析

```go
// 示例：生产者-消费者系统
func createProducerConsumerSystem() *PetriNet {
    pn := NewPetriNet()
    
    // 添加库所
    pn.AddPlace("producer_ready", "生产者就绪", 1)
    pn.AddPlace("buffer", "缓冲区", 5)
    pn.AddPlace("consumer_ready", "消费者就绪", 1)
    pn.AddPlace("produced", "已生产", -1)
    pn.AddPlace("consumed", "已消费", -1)
    
    // 添加变迁
    pn.AddTransition("produce", "生产", nil, nil)
    pn.AddTransition("consume", "消费", nil, nil)
    
    // 添加流关系
    pn.AddFlow("producer_ready", "produce", 1)
    pn.AddFlow("buffer", "consume", 1)
    pn.AddFlowFromTransition("produce", "buffer", 1)
    pn.AddFlowFromTransition("produce", "produced", 1)
    pn.AddFlowFromTransition("consume", "consumer_ready", 1)
    pn.AddFlowFromTransition("consume", "consumed", 1)
    
    // 设置初始标记
    pn.InitialMarking.Tokens["producer_ready"] = 1
    pn.InitialMarking.Tokens["consumer_ready"] = 1
    
    return pn
}
```

### 6.3 死锁检测

```go
// 死锁检测示例
func detectDeadlock() {
    // 创建一个可能导致死锁的Petri网
    pn := NewPetriNet()
    
    // 添加库所和变迁
    pn.AddPlace("p1", "Place1", 1)
    pn.AddPlace("p2", "Place2", 1)
    pn.AddTransition("t1", "Transition1", nil, nil)
    pn.AddTransition("t2", "Transition2", nil, nil)
    
    // 添加可能导致死锁的流关系
    pn.AddFlow("p1", "t1", 1)
    pn.AddFlow("p2", "t2", 1)
    pn.AddFlowFromTransition("t1", "p2", 1)
    pn.AddFlowFromTransition("t2", "p1", 1)
    
    // 设置初始标记
    pn.InitialMarking.Tokens["p1"] = 1
    
    // 分析死锁
    analyzer := NewPetriNetAnalyzer(pn)
    deadlocks := analyzer.AnalyzeDeadlock()
    
    if len(deadlocks) > 0 {
        fmt.Println("发现死锁状态:")
        for i, deadlock := range deadlocks {
            fmt.Printf("死锁 %d: %s\n", i+1, deadlock.String())
        }
    } else {
        fmt.Println("未发现死锁")
    }
}
```

## 总结

本文档详细介绍了Petri网模型在工作流系统中的应用，包括：

1. **基础理论**: Petri网的基本定义、标记和变迁、可达性
2. **工作流映射**: 将工作流映射到Petri网的方法
3. **高级模型**: 时间Petri网、颜色Petri网、层次Petri网
4. **形式化分析**: 可达性、活性、有界性分析
5. **Go语言实现**: 完整的Petri网和工作流Petri网实现
6. **应用案例**: 业务流程建模、并发系统分析、死锁检测

Petri网为工作流系统提供了强大的形式化建模和分析工具。
