# 01-Petri网模型 (Petri Net Model)

## 概述

Petri网是描述并发系统的经典形式化工具，特别适用于工作流建模。本模块基于Petri网理论，为工作流系统提供形式化的建模、分析和验证方法。

## 目录结构

### [01-Petri网基础](01-Petri-Net-Foundation/README.md)

- **01-Petri网定义** - 基本概念、形式化定义、图形表示
- **02-Petri网类型** - 基本Petri网、高级Petri网、有色Petri网
- **03-Petri网性质** - 可达性、活性、有界性、安全性
- **04-Petri网分析** - 状态空间分析、不变性分析、覆盖性分析

### [02-WF-net模型](02-WF-Net-Model/README.md)

- **01-WF-net定义** - 工作流Petri网、WF-net特性、健全性
- **02-WF-net性质** - 可达性、活性、有界性、健全性验证
- **03-WF-net分析** - 死锁检测、活锁检测、性能分析
- **04-WF-net优化** - 结构优化、性能优化、资源优化

### [03-工作流Petri网](03-Workflow-Petri-Net/README.md)

- **01-工作流映射** - 工作流到Petri网的映射规则
- **02-模式建模** - 顺序、并行、选择、循环模式建模
- **03-数据流建模** - 数据传递、数据转换、数据聚合建模
- **04-资源建模** - 资源分配、资源竞争、资源优化建模

### [04-Petri网分析](04-Petri-Net-Analysis/README.md)

- **01-可达性分析** - 可达性算法、状态空间构建、路径分析
- **02-死锁检测** - 死锁检测算法、死锁预防、死锁恢复
- **03-性能分析** - 执行时间分析、吞吐量分析、瓶颈识别
- **04-形式化验证** - 模型检查、性质验证、定理证明

## 核心概念

### 1. Petri网基础

#### 1.1 Petri网定义

Petri网是一个五元组 ```latex
$N = (P, T, F, W, M_0)$
```，其中：

- ```latex
$P = \{p_1, p_2, ..., p_n\}$
```：库所集合 (Places)
- ```latex
$T = \{t_1, t_2, ..., t_m\}$
```：变迁集合 (Transitions)
- ```latex
$F \subseteq (P \times T) \cup (T \times P)$
```：流关系 (Flow Relation)
- ```latex
$W: F \rightarrow \mathbb{N}^+$
```：权重函数 (Weight Function)
- ```latex
$M_0: P \rightarrow \mathbb{N}$
```：初始标识 (Initial Marking)

#### 1.2 图形表示

Petri网用图形表示：

- **库所**：用圆圈表示，包含令牌数量
- **变迁**：用矩形表示，表示事件或活动
- **弧**：用箭头表示，连接库所和变迁
- **令牌**：用小圆点表示，表示资源或状态

#### 1.3 激发规则

变迁 ```latex
$t$
``` 在标识 ```latex
$M$
``` 下可激发的条件：

1. ```latex
$\forall p \in \bullet t: M(p) \geq W(p, t)$
```
2. 激发后产生新标识 ```latex
$M'$
```：
   - ```latex
$M'(p) = M(p) - W(p, t)$
``` 如果 ```latex
$p \in \bullet t$
```
   - ```latex
$M'(p) = M(p) + W(t, p)$
``` 如果 ```latex
$p \in t \bullet$
```
   - ```latex
$M'(p) = M(p)$
``` 其他情况

### 2. WF-net模型

#### 2.1 WF-net定义

工作流Petri网 (WF-net) 是具有以下特性的Petri网：

1. **存在唯一的源库所i**：```latex
$\bullet i = \emptyset$
```
2. **存在唯一的汇库所o**：```latex
$o \bullet = \emptyset$
```
3. **网络中每个节点都在从i到o的路径上**

#### 2.2 WF-net健全性

WF-net ```latex
$N$
``` 是健全的当且仅当：

1. **可达性**：从初始标识可达终态标识
2. **活性**：不存在死锁
3. **有界性**：所有库所的令牌数有界
4. **完整性**：从任何可达标识都能到达终态

形式化定义：
```latex
$\text{sound}(N) \Leftrightarrow \text{live}(N) \land \text{bounded}(N) \land \text{complete}(N)$
```

#### 2.3 WF-net性质验证

- **可达性**：```latex
$M_0 \rightarrow^* M_f$
```
- **活性**：```latex
$\forall t \in T: \Box \Diamond \text{enabled}(t)$
```
- **有界性**：```latex
$\forall p \in P: \exists k \in \mathbb{N}: M(p) \leq k$
```
- **安全性**：```latex
$\Box \neg \text{deadlock}$
```

### 3. 工作流模式建模

#### 3.1 顺序模式

顺序模式 ```latex
$A \rightarrow B \rightarrow C$
``` 的Petri网表示：

```latex
p_1 \xrightarrow{t_A} p_2 \xrightarrow{t_B} p_3 \xrightarrow{t_C} p_4
```

其中：

- ```latex
$p_1$
```：开始库所
- ```latex
$p_4$
```：结束库所
- ```latex
$t_A, t_B, t_C$
```：对应活动A、B、C的变迁

#### 3.2 并行模式

并行模式 ```latex
$A \parallel B$
``` 的Petri网表示：

```latex
p_1 \xrightarrow{t_{split}} p_2 \parallel p_3
p_2 \xrightarrow{t_A} p_4
p_3 \xrightarrow{t_B} p_5
p_4 \parallel p_5 \xrightarrow{t_{join}} p_6
```

其中：

- ```latex
$t_{split}$
```：AND-split变迁
- ```latex
$t_{join}$
```：AND-join变迁
- ```latex
$p_2, p_3$
```：并行分支库所

#### 3.3 选择模式

选择模式 ```latex
$A | B$
``` 的Petri网表示：

```latex
p_1 \xrightarrow{t_{choice}} p_2 | p_3
p_2 \xrightarrow{t_A} p_4
p_3 \xrightarrow{t_B} p_5
p_4 | p_5 \xrightarrow{t_{merge}} p_6
```

其中：

- ```latex
$t_{choice}$
```：OR-split变迁
- ```latex
$t_{merge}$
```：OR-join变迁
- 选择基于条件或概率

#### 3.4 循环模式

循环模式 ```latex
$A \rightarrow (B \rightarrow A)^* \rightarrow C$
``` 的Petri网表示：

```latex
p_1 \xrightarrow{t_A} p_2 \xrightarrow{t_B} p_3
p_3 \xrightarrow{t_{loop}} p_2
p_3 \xrightarrow{t_{exit}} p_4 \xrightarrow{t_C} p_5
```

其中：

- ```latex
$t_{loop}$
```：循环变迁
- ```latex
$t_{exit}$
```：退出变迁
- 循环条件在 ```latex
$t_{loop}$
``` 中定义

### 4. Petri网分析

#### 4.1 可达性分析

可达性分析检查从初始标识是否可达目标标识：

```latex
\text{Reachability}: M_0 \rightarrow^* M
```

算法实现：

1. **状态空间构建**：从初始标识开始，探索所有可能的激发序列
2. **可达性图**：构建状态转换图
3. **可达性检查**：在图搜索中检查目标标识

#### 4.2 死锁检测

死锁检测检查是否存在无法继续执行的状态：

```latex
\text{Deadlock}: \exists M: M \text{ is reachable} \land \forall t \in T: \neg \text{enabled}(t, M)
```

算法实现：

1. **状态空间探索**：构建完整的可达性图
2. **死锁识别**：识别没有出边的状态
3. **死锁路径**：找到导致死锁的激发序列

#### 4.3 性能分析

性能分析评估Petri网的执行性能：

1. **执行时间分析**：
   - 关键路径分析
   - 平均执行时间计算
   - 最坏情况分析

2. **吞吐量分析**：
   - 稳态吞吐量计算
   - 瓶颈识别
   - 资源利用率分析

3. **资源分析**：
   - 资源需求分析
   - 资源竞争检测
   - 资源优化建议

## 技术栈

### Go语言实现

```go
// Petri网模型
type PetriNet struct {
    Places     map[string]*Place      `json:"places"`
    Transitions map[string]*Transition `json:"transitions"`
    Flow       map[string][]string    `json:"flow"`
    Weights    map[string]int         `json:"weights"`
    Marking    map[string]int         `json:"marking"`
    Initial    map[string]int         `json:"initial"`
}

// 库所
type Place struct {
    ID       string                 `json:"id"`
    Name     string                 `json:"name"`
    Tokens   int                    `json:"tokens"`
    Capacity int                    `json:"capacity"`
    Type     PlaceType              `json:"type"`
    Metadata map[string]interface{} `json:"metadata"`
}

// 变迁
type Transition struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Guard       string                 `json:"guard"`
    Action      string                 `json:"action"`
    Priority    int                    `json:"priority"`
    Timeout     time.Duration          `json:"timeout"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// WF-net工作流Petri网
type WFNet struct {
    PetriNet
    SourcePlace string `json:"source_place"`
    SinkPlace   string `json:"sink_place"`
}

// Petri网分析器
type PetriNetAnalyzer struct {
    net *PetriNet
}

// 可达性分析
func (pna *PetriNetAnalyzer) AnalyzeReachability() map[string]bool {
    reachable := make(map[string]bool)
    visited := make(map[string]bool)
    
    initialMarking := pna.net.markingToString(pna.net.Initial)
    pna.exploreReachability(initialMarking, reachable, visited)
    
    return reachable
}

// 探索可达性
func (pna *PetriNetAnalyzer) exploreReachability(marking string, reachable, visited map[string]bool) {
    if visited[marking] {
        return
    }
    
    visited[marking] = true
    reachable[marking] = true
    
    // 尝试所有可能的变迁
    for _, transition := range pna.net.Transitions {
        if pna.net.canFire(transition.ID, marking) {
            newMarking := pna.net.fireTransition(transition.ID, marking)
            pna.exploreReachability(newMarking, reachable, visited)
        }
    }
}

// 检查变迁是否可以激发
func (pn *PetriNet) canFire(transitionID string, marking string) bool {
    // 检查前置库所是否有足够的令牌
    for _, placeID := range pn.getPrePlaces(transitionID) {
        if pn.getTokenCount(placeID, marking) < pn.getWeight(placeID, transitionID) {
            return false
        }
    }
    return true
}

// 激发变迁
func (pn *PetriNet) fireTransition(transitionID string, marking string) string {
    newMarking := marking
    
    // 移除前置库所的令牌
    for _, placeID := range pn.getPrePlaces(transitionID) {
        weight := pn.getWeight(placeID, transitionID)
        newMarking = pn.removeTokens(placeID, weight, newMarking)
    }
    
    // 添加后置库所的令牌
    for _, placeID := range pn.getPostPlaces(transitionID) {
        weight := pn.getWeight(transitionID, placeID)
        newMarking = pn.addTokens(placeID, weight, newMarking)
    }
    
    return newMarking
}

// 死锁检测
func (pna *PetriNetAnalyzer) DetectDeadlocks() []string {
    deadlocks := make([]string, 0)
    reachable := pna.AnalyzeReachability()
    
    for marking := range reachable {
        if pna.isDeadlock(marking) {
            deadlocks = append(deadlocks, marking)
        }
    }
    
    return deadlocks
}

// 检查是否为死锁状态
func (pna *PetriNetAnalyzer) isDeadlock(marking string) bool {
    // 检查是否有任何变迁可以激发
    for _, transition := range pna.net.Transitions {
        if pna.net.canFire(transition.ID, marking) {
            return false
        }
    }
    return true
}

// WF-net健全性检查
func (wf *WFNet) CheckSoundness() SoundnessResult {
    result := SoundnessResult{}
    
    // 检查可达性
    analyzer := &PetriNetAnalyzer{net: &wf.PetriNet}
    reachable := analyzer.AnalyzeReachability()
    
    // 检查是否可达终态
    finalMarking := wf.getFinalMarking()
    result.Reachable = reachable[finalMarking]
    
    // 检查活性
    result.Live = wf.checkLiveness()
    
    // 检查有界性
    result.Bounded = wf.checkBoundedness()
    
    // 检查完整性
    result.Complete = wf.checkCompleteness()
    
    // 整体健全性
    result.Sound = result.Reachable && result.Live && result.Bounded && result.Complete
    
    return result
}

// 健全性检查结果
type SoundnessResult struct {
    Reachable bool `json:"reachable"`
    Live      bool `json:"live"`
    Bounded   bool `json:"bounded"`
    Complete  bool `json:"complete"`
    Sound     bool `json:"sound"`
}

// 工作流模式到Petri网的映射
type WorkflowToPetriNetMapper struct{}

// 映射顺序模式
func (wtpnm *WorkflowToPetriNetMapper) MapSequential(activities []Activity) *PetriNet {
    net := &PetriNet{
        Places:     make(map[string]*Place),
        Transitions: make(map[string]*Transition),
        Flow:       make(map[string][]string),
        Weights:    make(map[string]int),
    }
    
    // 创建库所和变迁
    for i, activity := range activities {
        placeID := fmt.Sprintf("p_%d", i)
        transitionID := fmt.Sprintf("t_%s", activity.ID)
        
        net.Places[placeID] = &Place{
            ID:     placeID,
            Name:   fmt.Sprintf("Place_%d", i),
            Tokens: 0,
        }
        
        net.Transitions[transitionID] = &Transition{
            ID:   transitionID,
            Name: activity.Name,
        }
        
        // 添加流关系
        if i > 0 {
            prevPlaceID := fmt.Sprintf("p_%d", i-1)
            net.Flow[prevPlaceID] = append(net.Flow[prevPlaceID], transitionID)
            net.Flow[transitionID] = append(net.Flow[transitionID], placeID)
            net.Weights[fmt.Sprintf("%s->%s", prevPlaceID, transitionID)] = 1
            net.Weights[fmt.Sprintf("%s->%s", transitionID, placeID)] = 1
        }
    }
    
    // 设置初始标识
    net.Initial[net.Places["p_0"].ID] = 1
    
    return net
}

// 映射并行模式
func (wtpnm *WorkflowToPetriNetMapper) MapParallel(activities []Activity) *PetriNet {
    net := &PetriNet{
        Places:     make(map[string]*Place),
        Transitions: make(map[string]*Transition),
        Flow:       make(map[string][]string),
        Weights:    make(map[string]int),
    }
    
    // 创建开始和结束库所
    startPlace := &Place{ID: "p_start", Name: "Start", Tokens: 1}
    endPlace := &Place{ID: "p_end", Name: "End", Tokens: 0}
    net.Places["p_start"] = startPlace
    net.Places["p_end"] = endPlace
    
    // 创建AND-split和AND-join变迁
    splitTransition := &Transition{ID: "t_split", Name: "AND-Split"}
    joinTransition := &Transition{ID: "t_join", Name: "AND-Join"}
    net.Transitions["t_split"] = splitTransition
    net.Transitions["t_join"] = joinTransition
    
    // 为每个活动创建库所和变迁
    for i, activity := range activities {
        placeID := fmt.Sprintf("p_%d", i)
        transitionID := fmt.Sprintf("t_%s", activity.ID)
        
        net.Places[placeID] = &Place{
            ID:     placeID,
            Name:   fmt.Sprintf("Place_%d", i),
            Tokens: 0,
        }
        
        net.Transitions[transitionID] = &Transition{
            ID:   transitionID,
            Name: activity.Name,
        }
        
        // 添加流关系
        net.Flow["p_start"] = append(net.Flow["p_start"], "t_split")
        net.Flow["t_split"] = append(net.Flow["t_split"], placeID)
        net.Flow[placeID] = append(net.Flow[placeID], transitionID)
        net.Flow[transitionID] = append(net.Flow[transitionID], "t_join")
        net.Flow["t_join"] = append(net.Flow["t_join"], "p_end")
        
        // 设置权重
        net.Weights["p_start->t_split"] = 1
        net.Weights["t_split->p_"+fmt.Sprintf("%d", i)] = 1
        net.Weights["p_"+fmt.Sprintf("%d", i)+"->t_"+activity.ID] = 1
        net.Weights["t_"+activity.ID+"->t_join"] = 1
        net.Weights["t_join->p_end"] = 1
    }
    
    // 设置初始标识
    net.Initial["p_start"] = 1
    
    return net
}
```

### 核心库

```go
import (
    "context"
    "time"
    "sync"
    "fmt"
    "encoding/json"
    
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "github.com/streadway/amqp"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/gorilla/websocket"
)
```

## 形式化规范

### 数学符号

使用LaTeX格式的数学公式：

```latex
\text{Petri网定义}: N = (P, T, F, W, M_0)

\text{库所集合}: P = \{p_1, p_2, ..., p_n\}

\text{变迁集合}: T = \{t_1, t_2, ..., t_m\}

\text{流关系}: F \subseteq (P \times T) \cup (T \times P)

\text{权重函数}: W: F \rightarrow \mathbb{N}^+

\text{初始标识}: M_0: P \rightarrow \mathbb{N}

\text{激发条件}: \forall p \in \bullet t: M(p) \geq W(p, t)

\text{激发后标识}: M'(p) = M(p) - W(p, t) + W(t, p)

\text{WF-net特性}: \bullet i = \emptyset \land o \bullet = \emptyset

\text{健全性}: \text{sound}(N) \Leftrightarrow \text{live}(N) \land \text{bounded}(N) \land \text{complete}(N)

\text{可达性}: M_0 \rightarrow^* M

\text{活性}: \forall t \in T: \Box \Diamond \text{enabled}(t)

\text{有界性}: \forall p \in P: \exists k \in \mathbb{N}: M(p) \leq k

\text{安全性}: \Box \neg \text{deadlock}

\text{顺序模式}: p_1 \xrightarrow{t_A} p_2 \xrightarrow{t_B} p_3 \xrightarrow{t_C} p_4

\text{并行模式}: p_1 \xrightarrow{t_{split}} p_2 \parallel p_3 \xrightarrow{t_{join}} p_4

\text{选择模式}: p_1 \xrightarrow{t_{choice}} p_2 | p_3 \xrightarrow{t_{merge}} p_4

\text{循环模式}: p_1 \xrightarrow{t_A} p_2 \xrightarrow{t_B} p_3 \xrightarrow{t_{loop}} p_2
```

### 算法分析

```go
// Petri网可达性分析算法
func (pna *PetriNetAnalyzer) AnalyzeReachability() map[string]bool {
    // 时间复杂度: O(2^{|P|})
    // 空间复杂度: O(2^{|P|})
    
    reachable := make(map[string]bool)
    visited := make(map[string]bool)
    
    initialMarking := pna.net.markingToString(pna.net.Initial)
    pna.exploreReachability(initialMarking, reachable, visited)
    
    return reachable
}

// 死锁检测算法
func (pna *PetriNetAnalyzer) DetectDeadlocks() []string {
    // 时间复杂度: O(2^{|P|} \times |T|)
    // 空间复杂度: O(2^{|P|})
    
    deadlocks := make([]string, 0)
    reachable := pna.AnalyzeReachability()
    
    for marking := range reachable {
        if pna.isDeadlock(marking) {
            deadlocks = append(deadlocks, marking)
        }
    }
    
    return deadlocks
}

// WF-net健全性检查算法
func (wf *WFNet) CheckSoundness() SoundnessResult {
    // 时间复杂度: O(2^{|P|} \times |T|)
    // 空间复杂度: O(2^{|P|})
    
    result := SoundnessResult{}
    
    // 可达性检查
    analyzer := &PetriNetAnalyzer{net: &wf.PetriNet}
    reachable := analyzer.AnalyzeReachability()
    finalMarking := wf.getFinalMarking()
    result.Reachable = reachable[finalMarking]
    
    // 活性检查
    result.Live = wf.checkLiveness()
    
    // 有界性检查
    result.Bounded = wf.checkBoundedness()
    
    // 完整性检查
    result.Complete = wf.checkCompleteness()
    
    result.Sound = result.Reachable && result.Live && result.Bounded && result.Complete
    
    return result
}

// 工作流模式映射算法
func (wtpnm *WorkflowToPetriNetMapper) MapWorkflow(workflow *Workflow) *PetriNet {
    // 时间复杂度: O(|A| + |T|)
    // 空间复杂度: O(|A| + |T|)
    
    net := &PetriNet{
        Places:     make(map[string]*Place),
        Transitions: make(map[string]*Transition),
        Flow:       make(map[string][]string),
        Weights:    make(map[string]int),
    }
    
    // 根据工作流类型选择映射策略
    switch workflow.Type {
    case WorkflowTypeSequential:
        return wtpnm.MapSequential(workflow.Activities)
    case WorkflowTypeParallel:
        return wtpnm.MapParallel(workflow.Activities)
    case WorkflowTypeChoice:
        return wtpnm.MapChoice(workflow.Activities)
    case WorkflowTypeIterative:
        return wtpnm.MapIterative(workflow.Activities)
    default:
        return wtpnm.MapSequential(workflow.Activities)
    }
}
```

## 质量保证

### 内容质量

- 不重复、分类严谨
- 与当前最新最成熟的哲科工程想法一致
- 符合学术要求
- 内容一致性、证明一致性、相关性一致性

### 结构质量

- 语义一致性
- 不交不空不漏的层次化分类
- 由理念到理性到形式化论证证明
- 有概念、定义的详细解释论证

## 本地跳转链接

- [返回工作流建模主目录](../README.md)
- [返回工作流系统主目录](../../README.md)
- [返回主目录](../../../../README.md)
- [01-基础理论层](../../../01-Foundation-Theory/README.md)
- [02-软件架构层](../../../02-Software-Architecture/README.md)
- [08-软件工程形式化](../../../08-Software-Engineering-Formalization/README.md)

---

**最后更新**: 2024年12月19日
**当前状态**: 🔄 第15轮重构进行中
**激情澎湃的持续构建** <(￣︶￣)↗[GO!] 🚀
