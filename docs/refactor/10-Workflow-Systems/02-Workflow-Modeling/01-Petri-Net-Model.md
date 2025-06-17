# 01-Petri网模型

(Petri Net Model)

## 目录

- [01-Petri网模型](#01-petri网模型)
  - [目录](#目录)
  - [1. Petri网基础理论](#1-petri网基础理论)
    - [1.1 基本定义](#11-基本定义)
    - [1.2 网结构](#12-网结构)
    - [1.3 标记和变迁](#13-标记和变迁)
  - [2. 工作流Petri网](#2-工作流petri网)
    - [2.1 工作流映射](#21-工作流映射)
    - [2.2 控制流建模](#22-控制流建模)
    - [2.3 数据流建模](#23-数据流建模)
  - [3. 形式化分析](#3-形式化分析)
    - [3.1 可达性分析](#31-可达性分析)
    - [3.2 活性分析](#32-活性分析)
    - [3.3 有界性分析](#33-有界性分析)
  - [4. 高级Petri网](#4-高级petri网)
    - [4.1 时间Petri网](#41-时间petri网)
    - [4.2 着色Petri网](#42-着色petri网)
    - [4.3 层次Petri网](#43-层次petri网)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 基础数据结构](#51-基础数据结构)
    - [5.2 网执行引擎](#52-网执行引擎)
    - [5.3 分析算法](#53-分析算法)
  - [6. 定理和证明](#6-定理和证明)
    - [6.1 可达性定理](#61-可达性定理)
    - [6.2 活性定理](#62-活性定理)
    - [6.3 有界性定理](#63-有界性定理)

## 1. Petri网基础理论

### 1.1 基本定义

**定义 1.1** (Petri网): Petri网是一个四元组 $N = (P, T, F, M_0)$，其中：
- $P$ 是库所(places)的有限集合
- $T$ 是变迁(transitions)的有限集合
- $F \subseteq (P \times T) \cup (T \times P)$ 是流关系
- $M_0: P \rightarrow \mathbb{N}$ 是初始标记

**定义 1.2** (前集和后集): 对于 $x \in P \cup T$：
- 前集：$^\bullet x = \{y \mid (y, x) \in F\}$
- 后集：$x^\bullet = \{y \mid (x, y) \in F\}$

**定义 1.3** (标记): 标记 $M: P \rightarrow \mathbb{N}$ 是库所到自然数的映射，表示每个库所中的托肯(token)数量。

### 1.2 网结构

**定义 1.4** (网结构): 网结构是三元组 $S = (P, T, F)$，其中：
- $P \cap T = \emptyset$ (库所和变迁不相交)
- $P \cup T \neq \emptyset$ (至少有一个元素)
- $F \subseteq (P \times T) \cup (T \times P)$ (流关系)

**性质 1.1** (网结构性质):
- 二元性：库所和变迁是对偶的
- 有向性：流关系是有向的
- 有限性：库所和变迁集合都是有限的

### 1.3 标记和变迁

**定义 1.5** (变迁使能): 变迁 $t \in T$ 在标记 $M$ 下使能，记作 $M[t\rangle$，当且仅当：
$$\forall p \in ^\bullet t: M(p) \geq F(p, t)$$

**定义 1.6** (变迁发生): 如果变迁 $t$ 在标记 $M$ 下使能，则 $t$ 可以发生，产生新标记 $M'$，记作 $M[t\rangle M'$，其中：
$$M'(p) = M(p) - F(p, t) + F(t, p)$$

**定义 1.7** (发生序列): 发生序列 $\sigma = t_1 t_2 \cdots t_n$ 是变迁序列，使得存在标记序列 $M_0, M_1, \ldots, M_n$ 满足：
$$M_0[t_1\rangle M_1[t_2\rangle \cdots [t_n\rangle M_n$$

## 2. 工作流Petri网

### 2.1 工作流映射

**定义 2.1** (工作流Petri网): 工作流Petri网 $WPN = (N, A, \lambda)$，其中：
- $N = (P, T, F, M_0)$ 是Petri网
- $A$ 是活动集合
- $\lambda: T \rightarrow A \cup \{\tau\}$ 是变迁到活动的映射，$\tau$ 表示静默变迁

**映射规则**:
1. 每个活动对应一个变迁
2. 活动之间的依赖关系对应流关系
3. 开始和结束状态对应特殊库所

**定义 2.2** (工作流结构): 工作流结构包含：
- 开始库所：$i \in P$，$^\bullet i = \emptyset$
- 结束库所：$o \in P$，$o^\bullet = \emptyset$
- 活动库所：$P_A = P \setminus \{i, o\}$

### 2.2 控制流建模

**定义 2.3** (顺序结构): 活动 $a_1$ 和 $a_2$ 的顺序结构建模为：

$$P = \{p_1, p_2, p_3\}, T = \{t_1, t_2\}, F = \{(p_1, t_1), (t_1, p_2), (p_2, t_2), (t_2, p_3)\}$$

**定义 2.4** (并行结构): 活动 $a_1$ 和 $a_2$ 的并行结构建模为：

$$P = \{p_1, p_2, p_3, p_4, p_5\}, T = \{t_1, t_2, t_3, t_4\}$$
$$F = \{(p_1, t_1), (t_1, p_2), (p_1, t_2), (t_2, p_3), (p_2, t_3), (p_3, t_3), (t_3, p_4), (p_4, t_4), (t_4, p_5)\}$$

**定义 2.5** (选择结构): 活动 $a_1$ 和 $a_2$ 的选择结构建模为：

$$P = \{p_1, p_2, p_3, p_4\}, T = \{t_1, t_2, t_3\}$$
$$F = \{(p_1, t_1), (p_1, t_2), (t_1, p_2), (t_2, p_3), (p_2, t_3), (p_3, t_3), (t_3, p_4)\}$$

### 2.3 数据流建模

**定义 2.6** (数据Petri网): 数据Petri网 $DPN = (N, D, \delta)$，其中：

- $N = (P, T, F, M_0)$ 是Petri网
- $D$ 是数据集合
- $\delta: T \rightarrow (D \rightarrow D)$ 是变迁到数据转换函数的映射

**定义 2.7** (数据流): 数据流是数据在活动之间的传递关系，建模为：

$$
\delta(t)(d) = \begin{cases}
d' & \text{if } t \text{ transforms } d \text{ to } d' \\
d & \text{otherwise}
\end{cases}
$$

## 3. 形式化分析

### 3.1 可达性分析

**定义 3.1** (可达性): 标记 $M$ 从 $M_0$ 可达，记作 $M_0 \rightarrow^* M$，如果存在发生序列 $\sigma$ 使得 $M_0[\sigma\rangle M$。

**定义 3.2** (可达图): 可达图 $RG(N) = (R(N), E)$，其中：

- $R(N)$ 是从 $M_0$ 可达的所有标记集合
- $E = \{(M, t, M') \mid M[t\rangle M'\}$

**算法 3.1** (可达性分析):
```go
func ReachabilityAnalysis(petriNet PetriNet) ReachabilityGraph {
    reachable := make(map[string]Marking)
    queue := []Marking{petriNet.InitialMarking}
    edges := []Edge{}
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        currentKey := current.String()
        if _, exists := reachable[currentKey]; exists {
            continue
        }
        
        reachable[currentKey] = current
        
        // 找到所有使能的变迁
        for _, transition := range petriNet.Transitions {
            if current.IsEnabled(transition) {
                next := current.Fire(transition)
                nextKey := next.String()
                
                edges = append(edges, Edge{
                    From: currentKey,
                    To:   nextKey,
                    Transition: transition,
                })
                
                if _, exists := reachable[nextKey]; !exists {
                    queue = append(queue, next)
                }
            }
        }
    }
    
    return ReachabilityGraph{
        Markings: reachable,
        Edges:    edges,
    }
}
```

### 3.2 活性分析

**定义 3.3** (活性): Petri网是活的，如果对于任意可达标记 $M$ 和任意变迁 $t$，存在从 $M$ 可达的标记 $M'$ 使得 $t$ 在 $M'$ 下使能。

**定义 3.4** (死锁): 标记 $M$ 是死锁，如果没有变迁在 $M$ 下使能。

**算法 3.2** (活性分析):
```go
func LivenessAnalysis(petriNet PetriNet) LivenessResult {
    reachabilityGraph := ReachabilityAnalysis(petriNet)
    deadlocks := []Marking{}
    liveTransitions := make(map[string]bool)
    
    // 检查死锁
    for _, marking := range reachabilityGraph.Markings {
        if !marking.HasEnabledTransitions(petriNet.Transitions) {
            deadlocks = append(deadlocks, marking)
        }
    }
    
    // 检查变迁活性
    for _, transition := range petriNet.Transitions {
        isLive := true
        for _, marking := range reachabilityGraph.Markings {
            if !marking.CanEnableTransition(transition, reachabilityGraph) {
                isLive = false
                break
            }
        }
        liveTransitions[transition.ID] = isLive
    }
    
    return LivenessResult{
        IsLive:           len(deadlocks) == 0,
        Deadlocks:        deadlocks,
        LiveTransitions:  liveTransitions,
    }
}
```

### 3.3 有界性分析

**定义 3.5** (有界性): Petri网是k-有界的，如果对于任意库所 $p$ 和任意可达标记 $M$，有 $M(p) \leq k$。

**定义 3.6** (安全性): Petri网是安全的，如果它是1-有界的。

**算法 3.3** (有界性分析):
```go
func BoundednessAnalysis(petriNet PetriNet) BoundednessResult {
    reachabilityGraph := ReachabilityAnalysis(petriNet)
    bounds := make(map[string]int)
    unboundedPlaces := []string{}
    
    for _, place := range petriNet.Places {
        maxTokens := 0
        for _, marking := range reachabilityGraph.Markings {
            if marking.Tokens[place.ID] > maxTokens {
                maxTokens = marking.Tokens[place.ID]
            }
        }
        
        if maxTokens == math.MaxInt32 {
            unboundedPlaces = append(unboundedPlaces, place.ID)
        } else {
            bounds[place.ID] = maxTokens
        }
    }
    
    return BoundednessResult{
        IsBounded:        len(unboundedPlaces) == 0,
        Bounds:           bounds,
        UnboundedPlaces:  unboundedPlaces,
    }
}
```

## 4. 高级Petri网

### 4.1 时间Petri网

**定义 4.1** (时间Petri网): 时间Petri网 $TPN = (N, I)$，其中：
- $N = (P, T, F, M_0)$ 是Petri网
- $I: T \rightarrow \mathbb{R}^+ \times \mathbb{R}^+$ 是时间间隔函数

**定义 4.2** (时间变迁): 时间变迁 $t$ 的时间间隔 $I(t) = [EFT(t), LFT(t)]$，其中：
- $EFT(t)$ 是最早使能时间
- $LFT(t)$ 是最晚使能时间

### 4.2 着色Petri网

**定义 4.3** (着色Petri网): 着色Petri网 $CPN = (N, \Sigma, C, G, E)$，其中：
- $N = (P, T, F, M_0)$ 是Petri网
- $\Sigma$ 是颜色集合
- $C: P \cup T \rightarrow \Sigma$ 是颜色函数
- $G: T \rightarrow \text{Guard}$ 是守卫函数
- $E: F \rightarrow \text{Expression}$ 是表达式函数

### 4.3 层次Petri网

**定义 4.4** (层次Petri网): 层次Petri网 $HPN = (N, H)$，其中：
- $N = (P, T, F, M_0)$ 是Petri网
- $H: T \rightarrow \text{Subnet}$ 是层次函数，将变迁映射到子网

## 5. Go语言实现

### 5.1 基础数据结构

```go
// PetriNet Petri网
type PetriNet struct {
    Places           []Place
    Transitions      []Transition
    FlowRelation     map[string]map[string]int
    InitialMarking   Marking
}

// Place 库所
type Place struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Capacity int    `json:"capacity"`
}

// Transition 变迁
type Transition struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Guard    Guard  `json:"guard"`
    Action   Action `json:"action"`
}

// Marking 标记
type Marking struct {
    Tokens map[string]int `json:"tokens"`
}

func (m Marking) String() string {
    tokens := make([]string, 0, len(m.Tokens))
    for place, count := range m.Tokens {
        tokens = append(tokens, fmt.Sprintf("%s:%d", place, count))
    }
    return fmt.Sprintf("{%s}", strings.Join(tokens, ", "))
}

func (m Marking) IsEnabled(transition Transition, petriNet PetriNet) bool {
    // 检查变迁是否使能
    for placeID, requiredTokens := range petriNet.getInputTokens(transition.ID) {
        if m.Tokens[placeID] < requiredTokens {
            return false
        }
    }
    return true
}

func (m Marking) Fire(transition Transition, petriNet PetriNet) Marking {
    newMarking := Marking{
        Tokens: make(map[string]int),
    }
    
    // 复制当前标记
    for place, tokens := range m.Tokens {
        newMarking.Tokens[place] = tokens
    }
    
    // 消耗输入托肯
    for placeID, tokens := range petriNet.getInputTokens(transition.ID) {
        newMarking.Tokens[placeID] -= tokens
    }
    
    // 产生输出托肯
    for placeID, tokens := range petriNet.getOutputTokens(transition.ID) {
        newMarking.Tokens[placeID] += tokens
    }
    
    return newMarking
}

// Guard 守卫条件
type Guard interface {
    Evaluate(ctx context.Context, data interface{}) (bool, error)
}

// Action 动作
type Action interface {
    Execute(ctx context.Context, data interface{}) (interface{}, error)
}

// ReachabilityGraph 可达图
type ReachabilityGraph struct {
    Markings map[string]Marking `json:"markings"`
    Edges    []Edge             `json:"edges"`
}

// Edge 边
type Edge struct {
    From       string     `json:"from"`
    To         string     `json:"to"`
    Transition Transition `json:"transition"`
}
```

### 5.2 网执行引擎

```go
// PetriNetEngine Petri网执行引擎
type PetriNetEngine struct {
    PetriNet     PetriNet
    CurrentState Marking
    Data         interface{}
}

func NewPetriNetEngine(petriNet PetriNet) *PetriNetEngine {
    return &PetriNetEngine{
        PetriNet:     petriNet,
        CurrentState: petriNet.InitialMarking,
        Data:         nil,
    }
}

func (pne *PetriNetEngine) Execute(ctx context.Context, data interface{}) error {
    pne.Data = data
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 找到所有使能的变迁
            enabledTransitions := pne.getEnabledTransitions()
            if len(enabledTransitions) == 0 {
                return fmt.Errorf("no enabled transitions")
            }
            
            // 选择下一个变迁（这里简化，实际应该有选择策略）
            transition := enabledTransitions[0]
            
            // 执行变迁
            if err := pne.fireTransition(ctx, transition); err != nil {
                return fmt.Errorf("failed to fire transition %s: %w", transition.ID, err)
            }
            
            // 检查是否到达终止状态
            if pne.isTerminalState() {
                break
            }
        }
    }
    
    return nil
}

func (pne *PetriNetEngine) getEnabledTransitions() []Transition {
    var enabled []Transition
    
    for _, transition := range pne.PetriNet.Transitions {
        if pne.CurrentState.IsEnabled(transition, pne.PetriNet) {
            // 检查守卫条件
            if transition.Guard != nil {
                if enabled, err := transition.Guard.Evaluate(context.Background(), pne.Data); err == nil && enabled {
                    enabled = append(enabled, transition)
                }
            } else {
                enabled = append(enabled, transition)
            }
        }
    }
    
    return enabled
}

func (pne *PetriNetEngine) fireTransition(ctx context.Context, transition Transition) error {
    // 执行动作
    if transition.Action != nil {
        newData, err := transition.Action.Execute(ctx, pne.Data)
        if err != nil {
            return err
        }
        pne.Data = newData
    }
    
    // 更新标记
    pne.CurrentState = pne.CurrentState.Fire(transition, pne.PetriNet)
    
    return nil
}

func (pne *PetriNetEngine) isTerminalState() bool {
    // 检查是否到达终止状态（所有输出库所都有托肯）
    for _, place := range pne.PetriNet.Places {
        if pne.isOutputPlace(place.ID) && pne.CurrentState.Tokens[place.ID] == 0 {
            return false
        }
    }
    return true
}

func (pne *PetriNetEngine) isOutputPlace(placeID string) bool {
    // 检查是否是输出库所（没有输出变迁）
    for _, transition := range pne.PetriNet.Transitions {
        if pne.PetriNet.hasFlow(placeID, transition.ID) {
            return false
        }
    }
    return true
}
```

### 5.3 分析算法

```go
// PetriNetAnalyzer Petri网分析器
type PetriNetAnalyzer struct {
    PetriNet PetriNet
}

func NewPetriNetAnalyzer(petriNet PetriNet) *PetriNetAnalyzer {
    return &PetriNetAnalyzer{
        PetriNet: petriNet,
    }
}

func (pna *PetriNetAnalyzer) Analyze() AnalysisResult {
    return AnalysisResult{
        Reachability: pna.analyzeReachability(),
        Liveness:     pna.analyzeLiveness(),
        Boundedness:  pna.analyzeBoundedness(),
        Safety:       pna.analyzeSafety(),
    }
}

func (pna *PetriNetAnalyzer) analyzeReachability() ReachabilityResult {
    reachabilityGraph := ReachabilityAnalysis(pna.PetriNet)
    
    return ReachabilityResult{
        TotalMarkings: len(reachabilityGraph.Markings),
        TotalEdges:    len(reachabilityGraph.Edges),
        Graph:         reachabilityGraph,
    }
}

func (pna *PetriNetAnalyzer) analyzeLiveness() LivenessResult {
    return LivenessAnalysis(pna.PetriNet)
}

func (pna *PetriNetAnalyzer) analyzeBoundedness() BoundednessResult {
    return BoundednessAnalysis(pna.PetriNet)
}

func (pna *PetriNetAnalyzer) analyzeSafety() SafetyResult {
    boundedness := pna.analyzeBoundedness()
    
    isSafe := true
    unsafePlaces := []string{}
    
    for placeID, bound := range boundedness.Bounds {
        if bound > 1 {
            isSafe = false
            unsafePlaces = append(unsafePlaces, placeID)
        }
    }
    
    return SafetyResult{
        IsSafe:        isSafe,
        UnsafePlaces:  unsafePlaces,
    }
}

// AnalysisResult 分析结果
type AnalysisResult struct {
    Reachability ReachabilityResult `json:"reachability"`
    Liveness     LivenessResult     `json:"liveness"`
    Boundedness  BoundednessResult  `json:"boundedness"`
    Safety       SafetyResult       `json:"safety"`
}

// ReachabilityResult 可达性分析结果
type ReachabilityResult struct {
    TotalMarkings int                 `json:"total_markings"`
    TotalEdges    int                 `json:"total_edges"`
    Graph         ReachabilityGraph   `json:"graph"`
}

// LivenessResult 活性分析结果
type LivenessResult struct {
    IsLive          bool                `json:"is_live"`
    Deadlocks       []Marking           `json:"deadlocks"`
    LiveTransitions map[string]bool     `json:"live_transitions"`
}

// BoundednessResult 有界性分析结果
type BoundednessResult struct {
    IsBounded       bool                `json:"is_bounded"`
    Bounds          map[string]int      `json:"bounds"`
    UnboundedPlaces []string            `json:"unbounded_places"`
}

// SafetyResult 安全性分析结果
type SafetyResult struct {
    IsSafe         bool     `json:"is_safe"`
    UnsafePlaces   []string `json:"unsafe_places"`
}
```

## 6. 定理和证明

### 6.1 可达性定理

**定理 6.1** (可达性判定): 标记 $M$ 从 $M_0$ 可达当且仅当存在非负整数向量 $x$ 使得：
$$M = M_0 + C \cdot x$$

其中 $C$ 是关联矩阵。

**证明**:
1. 必要性：如果 $M$ 可达，则存在发生序列 $\sigma$ 使得 $M_0[\sigma\rangle M$
2. 设 $x$ 是 $\sigma$ 的Parikh向量
3. 则 $M = M_0 + C \cdot x$
4. 充分性：如果存在 $x$ 使得 $M = M_0 + C \cdot x$，则存在发生序列 $\sigma$ 使得 $M_0[\sigma\rangle M$
5. 因此可达性判定成立。$\square$

### 6.2 活性定理

**定理 6.2** (活性判定): Petri网是活的当且仅当可达图中没有死锁状态。

**证明**:
1. 必要性：如果网是活的，则任意可达标记都能使能某个变迁
2. 因此没有死锁状态
3. 充分性：如果没有死锁状态，则任意可达标记都能使能某个变迁
4. 因此网是活的。$\square$

### 6.3 有界性定理

**定理 6.3** (有界性判定): Petri网是k-有界的当且仅当可达图中所有标记的托肯数量都不超过k。

**证明**:
1. 必要性：如果网是k-有界的，则任意可达标记的托肯数量都不超过k
2. 充分性：如果可达图中所有标记的托肯数量都不超过k，则网是k-有界的
3. 因此有界性判定成立。$\square$

---

**参考文献**:
1. Murata, T. (1989). Petri Nets: Properties, Analysis and Applications. Proceedings of the IEEE, 77(4), 541-580.
2. van der Aalst, W. M. P. (1998). The Application of Petri Nets to Workflow Management. The Journal of Circuits, Systems and Computers, 8(1), 21-66.
3. Reisig, W. (2013). Understanding Petri Nets: Modeling Techniques, Analysis Methods, Case Studies. Springer.
4. Jensen, K., & Kristensen, L. M. (2009). Coloured Petri Nets: Modelling and Validation of Concurrent Systems. Springer. 