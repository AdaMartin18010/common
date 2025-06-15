# 04-时态逻辑 (Temporal Logic)

## 目录

- [04-时态逻辑 (Temporal Logic)](#04-时态逻辑-temporal-logic)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 时态算子](#11-时态算子)
    - [1.2 时间模型](#12-时间模型)
    - [1.3 时态系统](#13-时态系统)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 线性时态逻辑](#21-线性时态逻辑)
    - [2.2 分支时态逻辑](#22-分支时态逻辑)
    - [2.3 区间时态逻辑](#23-区间时态逻辑)
  - [3. 推理系统](#3-推理系统)
    - [3.1 公理系统](#31-公理系统)
    - [3.2 表推演](#32-表推演)
    - [3.3 模型检查](#33-模型检查)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 时态逻辑数据结构](#41-时态逻辑数据结构)
    - [4.2 时间模型实现](#42-时间模型实现)
    - [4.3 推理引擎](#43-推理引擎)
  - [5. 应用示例](#5-应用示例)
    - [5.1 程序验证](#51-程序验证)
    - [5.2 协议验证](#52-协议验证)
    - [5.3 实时系统](#53-实时系统)
  - [总结](#总结)

## 1. 基本概念

### 1.1 时态算子

**定义 1.1**: 时态算子用于表达时间相关的概念。

**基本时态算子**:

1. **G (Globally)**: "总是" - $G \phi$ 表示"在所有未来时刻都成立"
2. **F (Finally)**: "最终" - $F \phi$ 表示"在某个未来时刻成立"
3. **X (Next)**: "下一个" - $X \phi$ 表示"在下一个时刻成立"
4. **U (Until)**: "直到" - $\phi U \psi$ 表示"$\phi$ 成立直到 $\psi$ 成立"

**派生算子**:

- **R (Release)**: $\phi R \psi \equiv \neg(\neg \phi U \neg \psi)$
- **W (Weak Until)**: $\phi W \psi \equiv (\phi U \psi) \lor G \phi$

### 1.2 时间模型

**定义 1.2**: 时间结构

**线性时间结构**: $(T, <)$ 其中 $T$ 是时间点集合，$<$ 是严格全序关系。

**分支时间结构**: $(T, <)$ 其中 $T$ 是时间点集合，$<$ 是严格偏序关系。

**定义 1.3**: 路径

**路径** $\pi = t_0, t_1, t_2, \ldots$ 是时间点的无限序列，满足：

- $t_i < t_{i+1}$ 对所有 $i \geq 0$
- 如果存在 $t$ 使得 $t_i < t < t_{i+1}$，则 $t = t_i$ 或 $t = t_{i+1}$

### 1.3 时态系统

**定义 1.4**: 常见时态系统

1. **LTL (Linear Temporal Logic)**: 线性时态逻辑
2. **CTL (Computation Tree Logic)**: 计算树逻辑
3. **CTL* (CTL Star)**: 统一时态逻辑
4. **MTL (Metric Temporal Logic)**: 度量时态逻辑

## 2. 形式化定义

### 2.1 线性时态逻辑

**定义 2.1**: LTL语法

LTL公式的递归定义：

1. **基础**: 每个原子命题 $p \in \mathcal{P}$ 是公式
2. **归纳**: 如果 $\phi$ 和 $\psi$ 是公式，则：
   - $\neg \phi$ 是公式
   - $(\phi \land \psi)$ 是公式
   - $(\phi \lor \psi)$ 是公式
   - $(\phi \rightarrow \psi)$ 是公式
   - $X \phi$ 是公式
   - $F \phi$ 是公式
   - $G \phi$ 是公式
   - $(\phi U \psi)$ 是公式
3. **闭包**: 只有通过有限次应用上述规则得到的才是公式

**定义 2.2**: LTL语义

对于路径 $\pi = t_0, t_1, t_2, \ldots$ 和位置 $i \geq 0$：

- $\pi, i \models p$ 当且仅当 $p \in V(t_i)$
- $\pi, i \models \neg \phi$ 当且仅当 $\pi, i \not\models \phi$
- $\pi, i \models \phi \land \psi$ 当且仅当 $\pi, i \models \phi$ 且 $\pi, i \models \psi$
- $\pi, i \models X \phi$ 当且仅当 $\pi, i+1 \models \phi$
- $\pi, i \models F \phi$ 当且仅当存在 $j \geq i$，$\pi, j \models \phi$
- $\pi, i \models G \phi$ 当且仅当对所有 $j \geq i$，$\pi, j \models \phi$
- $\pi, i \models \phi U \psi$ 当且仅当存在 $j \geq i$，$\pi, j \models \psi$ 且对所有 $k$，$i \leq k < j$，$\pi, k \models \phi$

### 2.2 分支时态逻辑

**定义 2.3**: CTL语法

CTL公式的递归定义：

1. **状态公式**:
   - 原子命题 $p \in \mathcal{P}$ 是状态公式
   - 如果 $\phi$ 和 $\psi$ 是状态公式，则 $\neg \phi$，$(\phi \land \psi)$ 等是状态公式
   - 如果 $\phi$ 是路径公式，则 $A \phi$ 和 $E \phi$ 是状态公式

2. **路径公式**:
   - 如果 $\phi$ 是状态公式，则 $\phi$ 是路径公式
   - 如果 $\phi$ 和 $\psi$ 是路径公式，则 $X \phi$，$F \phi$，$G \phi$，$(\phi U \psi)$ 是路径公式

**定义 2.4**: CTL语义

- $s \models A \phi$ 当且仅当对所有从 $s$ 开始的路径 $\pi$，$\pi \models \phi$
- $s \models E \phi$ 当且仅当存在从 $s$ 开始的路径 $\pi$，$\pi \models \phi$

### 2.3 区间时态逻辑

**定义 2.5**: 区间时态逻辑

区间时态逻辑允许表达时间区间上的性质：

- **区间算子**: $[a,b] \phi$ 表示"在时间区间 $[a,b]$ 内 $\phi$ 成立"
- **持续时间**: $\phi \text{ for } d$ 表示"$\phi$ 持续 $d$ 个时间单位"

## 3. 推理系统

### 3.1 公理系统

**定义 3.1**: LTL公理系统

**公理**:

1. **命题公理**: 所有命题逻辑重言式
2. **时态公理**:
   - $G \phi \leftrightarrow \neg F \neg \phi$
   - $F \phi \leftrightarrow \text{true} U \phi$
   - $G(\phi \rightarrow \psi) \rightarrow (G \phi \rightarrow G \psi)$
   - $G \phi \rightarrow \phi \land X G \phi$
   - $\phi U \psi \leftrightarrow \psi \lor (\phi \land X(\phi U \psi))$

**推理规则**:

- **假言推理**: 从 $\phi$ 和 $\phi \rightarrow \psi$ 推出 $\psi$
- **必然化**: 从 $\phi$ 推出 $G \phi$

### 3.2 表推演

**定义 3.2**: 时态表推演规则

对于时态公式：

- **$X$公式**: 将 $X \phi$ 分解为 $\phi$，但移到下一个状态
- **$F$公式**: 将 $F \phi$ 分解为 $\phi$ 或 $X F \phi$
- **$G$公式**: 将 $G \phi$ 分解为 $\phi \land X G \phi$
- **$U$公式**: 将 $\phi U \psi$ 分解为 $\psi \lor (\phi \land X(\phi U \psi))$

### 3.3 模型检查

**定义 3.3**: 时态模型检查算法

```pseudocode
ModelCheck(φ, M, s):
    case φ of
        p: return s ∈ V(p)
        ¬ψ: return not ModelCheck(ψ, M, s)
        ψ ∧ χ: return ModelCheck(ψ, M, s) and ModelCheck(χ, M, s)
        Xψ: return for all s' such that s → s': ModelCheck(ψ, M, s')
        Fψ: return CheckEventually(ψ, M, s, {})
        Gψ: return CheckAlways(ψ, M, s, {})
        ψ U χ: return CheckUntil(ψ, χ, M, s, {})

CheckEventually(ψ, M, s, visited):
    if s ∈ visited: return false
    if ModelCheck(ψ, M, s): return true
    visited = visited ∪ {s}
    return for some s' such that s → s': CheckEventually(ψ, M, s', visited)

CheckAlways(ψ, M, s, visited):
    if s ∈ visited: return true
    if not ModelCheck(ψ, M, s): return false
    visited = visited ∪ {s}
    return for all s' such that s → s': CheckAlways(ψ, M, s', visited)
```

## 4. Go语言实现

### 4.1 时态逻辑数据结构

```go
// TemporalOperator 时态算子
type TemporalOperator int

const (
    Next TemporalOperator = iota
    Finally
    Globally
    Until
    Release
    WeakUntil
)

// TemporalFormula 时态逻辑公式
type TemporalFormula struct {
    IsAtom      bool
    IsNegation  bool
    IsBinary    bool
    IsTemporal  bool
    
    // 原子公式
    Proposition string
    
    // 连接词
    Connective string
    Left       *TemporalFormula
    Right      *TemporalFormula
    
    // 时态算子
    Operator TemporalOperator
    Body     *TemporalFormula
    UntilLeft *TemporalFormula // 用于Until算子
}

// NewAtomFormula 创建原子公式
func NewAtomFormula(proposition string) *TemporalFormula {
    return &TemporalFormula{
        IsAtom:      true,
        Proposition: proposition,
    }
}

// NewNegation 创建否定公式
func NewNegation(formula *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsNegation: true,
        Left:       formula,
    }
}

// NewBinaryFormula 创建二元连接词公式
func NewBinaryFormula(connective string, left, right *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsBinary:   true,
        Connective: connective,
        Left:       left,
        Right:      right,
    }
}

// NewNext 创建下一个公式
func NewNext(body *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsTemporal: true,
        Operator:   Next,
        Body:       body,
    }
}

// NewFinally 创建最终公式
func NewFinally(body *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsTemporal: true,
        Operator:   Finally,
        Body:       body,
    }
}

// NewGlobally 创建全局公式
func NewGlobally(body *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsTemporal: true,
        Operator:   Globally,
        Body:       body,
    }
}

// NewUntil 创建直到公式
func NewUntil(left, right *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsTemporal: true,
        Operator:   Until,
        UntilLeft:  left,
        Body:       right,
    }
}

// NewRelease 创建释放公式
func NewRelease(left, right *TemporalFormula) *TemporalFormula {
    return &TemporalFormula{
        IsTemporal: true,
        Operator:   Release,
        UntilLeft:  left,
        Body:       right,
    }
}
```

### 4.2 时间模型实现

```go
// State 状态
type State struct {
    ID   string
    Name string
}

// NewState 创建状态
func NewState(id, name string) *State {
    return &State{
        ID:   id,
        Name: name,
    }
}

// TransitionSystem 转移系统
type TransitionSystem struct {
    States       map[string]*State
    Transitions  map[string]map[string]bool // s1 -> s2
    InitialState string
    Valuation    map[string]map[string]bool // V(p, s) = true 表示命题p在状态s中为真
}

// NewTransitionSystem 创建转移系统
func NewTransitionSystem() *TransitionSystem {
    return &TransitionSystem{
        States:      make(map[string]*State),
        Transitions: make(map[string]map[string]bool),
        Valuation:   make(map[string]map[string]bool),
    }
}

// AddState 添加状态
func (ts *TransitionSystem) AddState(state *State) {
    ts.States[state.ID] = state
    ts.Transitions[state.ID] = make(map[string]bool)
}

// AddTransition 添加转移
func (ts *TransitionSystem) AddTransition(from, to string) {
    if ts.Transitions[from] == nil {
        ts.Transitions[from] = make(map[string]bool)
    }
    ts.Transitions[from][to] = true
}

// SetValuation 设置赋值
func (ts *TransitionSystem) SetValuation(proposition, state string, value bool) {
    if ts.Valuation[proposition] == nil {
        ts.Valuation[proposition] = make(map[string]bool)
    }
    ts.Valuation[proposition][state] = value
}

// GetValuation 获取赋值
func (ts *TransitionSystem) GetValuation(proposition, state string) bool {
    if ts.Valuation[proposition] == nil {
        return false
    }
    return ts.Valuation[proposition][state]
}

// GetSuccessors 获取后继状态
func (ts *TransitionSystem) GetSuccessors(state string) []string {
    successors := []string{}
    if ts.Transitions[state] != nil {
        for successor := range ts.Transitions[state] {
            successors = append(successors, successor)
        }
    }
    return successors
}

// Evaluate 计算公式在给定状态中的真值
func (ts *TransitionSystem) Evaluate(formula *TemporalFormula, state string) bool {
    if formula.IsAtom {
        return ts.GetValuation(formula.Proposition, state)
    }
    
    if formula.IsNegation {
        return !ts.Evaluate(formula.Left, state)
    }
    
    if formula.IsBinary {
        left := ts.Evaluate(formula.Left, state)
        right := ts.Evaluate(formula.Right, state)
        
        switch formula.Connective {
        case "∧":
            return left && right
        case "∨":
            return left || right
        case "→":
            return !left || right
        case "↔":
            return left == right
        }
    }
    
    if formula.IsTemporal {
        return ts.evaluateTemporal(formula, state)
    }
    
    return false
}

// evaluateTemporal 计算时态公式的真值
func (ts *TransitionSystem) evaluateTemporal(formula *TemporalFormula, state string) bool {
    switch formula.Operator {
    case Next:
        return ts.evaluateNext(formula.Body, state)
    case Finally:
        return ts.evaluateFinally(formula.Body, state, make(map[string]bool))
    case Globally:
        return ts.evaluateGlobally(formula.Body, state, make(map[string]bool))
    case Until:
        return ts.evaluateUntil(formula.UntilLeft, formula.Body, state, make(map[string]bool))
    case Release:
        return ts.evaluateRelease(formula.UntilLeft, formula.Body, state, make(map[string]bool))
    }
    return false
}

// evaluateNext 计算下一个公式
func (ts *TransitionSystem) evaluateNext(formula *TemporalFormula, state string) bool {
    successors := ts.GetSuccessors(state)
    if len(successors) == 0 {
        return false
    }
    
    // 对于所有后继状态，公式必须为真
    for _, successor := range successors {
        if !ts.Evaluate(formula, successor) {
            return false
        }
    }
    return true
}

// evaluateFinally 计算最终公式
func (ts *TransitionSystem) evaluateFinally(formula *TemporalFormula, state string, visited map[string]bool) bool {
    if visited[state] {
        return false
    }
    
    if ts.Evaluate(formula, state) {
        return true
    }
    
    visited[state] = true
    successors := ts.GetSuccessors(state)
    
    for _, successor := range successors {
        if ts.evaluateFinally(formula, successor, visited) {
            return true
        }
    }
    
    return false
}

// evaluateGlobally 计算全局公式
func (ts *TransitionSystem) evaluateGlobally(formula *TemporalFormula, state string, visited map[string]bool) bool {
    if visited[state] {
        return true
    }
    
    if !ts.Evaluate(formula, state) {
        return false
    }
    
    visited[state] = true
    successors := ts.GetSuccessors(state)
    
    for _, successor := range successors {
        if !ts.evaluateGlobally(formula, successor, visited) {
            return false
        }
    }
    
    return true
}

// evaluateUntil 计算直到公式
func (ts *TransitionSystem) evaluateUntil(left, right *TemporalFormula, state string, visited map[string]bool) bool {
    if visited[state] {
        return false
    }
    
    if ts.Evaluate(right, state) {
        return true
    }
    
    if !ts.Evaluate(left, state) {
        return false
    }
    
    visited[state] = true
    successors := ts.GetSuccessors(state)
    
    for _, successor := range successors {
        if ts.evaluateUntil(left, right, successor, visited) {
            return true
        }
    }
    
    return false
}

// evaluateRelease 计算释放公式
func (ts *TransitionSystem) evaluateRelease(left, right *TemporalFormula, state string, visited map[string]bool) bool {
    // φ R ψ ≡ ¬(¬φ U ¬ψ)
    notLeft := NewNegation(left)
    notRight := NewNegation(right)
    notUntil := NewNegation(NewUntil(notLeft, notRight))
    
    return ts.Evaluate(notUntil, state)
}
```

### 4.3 推理引擎

```go
// TemporalLogicEngine 时态逻辑推理引擎
type TemporalLogicEngine struct {
    system *TransitionSystem
}

// NewTemporalLogicEngine 创建时态逻辑推理引擎
func NewTemporalLogicEngine() *TemporalLogicEngine {
    return &TemporalLogicEngine{
        system: NewTransitionSystem(),
    }
}

// SetupExampleSystem 设置示例系统
func (e *TemporalLogicEngine) SetupExampleSystem() {
    // 创建一个简单的转移系统
    s1 := NewState("s1", "状态1")
    s2 := NewState("s2", "状态2")
    s3 := NewState("s3", "状态3")
    
    e.system.AddState(s1)
    e.system.AddState(s2)
    e.system.AddState(s3)
    
    // 添加转移
    e.system.AddTransition("s1", "s2")
    e.system.AddTransition("s2", "s3")
    e.system.AddTransition("s3", "s3") // 自环
    
    // 设置赋值
    e.system.SetValuation("p", "s1", true)
    e.system.SetValuation("p", "s2", false)
    e.system.SetValuation("p", "s3", true)
    
    e.system.SetValuation("q", "s1", false)
    e.system.SetValuation("q", "s2", true)
    e.system.SetValuation("q", "s3", false)
    
    e.system.InitialState = "s1"
}

// ModelChecking 模型检查
func (e *TemporalLogicEngine) ModelChecking(formula *TemporalFormula) {
    fmt.Printf("模型检查公式: %s\n", formula.String())
    
    for stateID := range e.system.States {
        value := e.system.Evaluate(formula, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
}

// ProveTemporalEquivalence 证明时态等价
func (e *TemporalLogicEngine) ProveTemporalEquivalence() {
    // 证明 Fφ ≡ true U φ
    
    p := NewAtomFormula("p")
    trueProp := NewAtomFormula("true")
    
    finallyP := NewFinally(p)
    untilTrueP := NewUntil(trueProp, p)
    
    fmt.Println("证明 Fφ ≡ true U φ")
    
    for stateID := range e.system.States {
        finallyValue := e.system.Evaluate(finallyP, stateID)
        untilValue := e.system.Evaluate(untilTrueP, stateID)
        
        fmt.Printf("状态 %s: Fp = %v, true U p = %v, 等价 = %v\n", 
            stateID, finallyValue, untilValue, finallyValue == untilValue)
    }
}

// ProveTemporalAxioms 证明时态公理
func (e *TemporalLogicEngine) ProveTemporalAxioms() {
    // 证明 Gφ → φ
    
    p := NewAtomFormula("p")
    globallyP := NewGlobally(p)
    axiom := NewBinaryFormula("→", globallyP, p)
    
    fmt.Println("证明公理: Gφ → φ")
    
    for stateID := range e.system.States {
        value := e.system.Evaluate(axiom, stateID)
        fmt.Printf("状态 %s: Gp → p = %v\n", stateID, value)
    }
}
```

## 5. 应用示例

### 5.1 程序验证

```go
// ProgramVerification 程序验证示例
func ProgramVerification() {
    // 验证一个简单的程序：while (x > 0) { x = x - 1; }
    
    system := NewTransitionSystem()
    
    // 状态：s0(x=3), s1(x=2), s2(x=1), s3(x=0)
    s0 := NewState("s0", "x=3")
    s1 := NewState("s1", "x=2")
    s2 := NewState("s2", "x=1")
    s3 := NewState("s3", "x=0")
    
    system.AddState(s0)
    system.AddState(s1)
    system.AddState(s2)
    system.AddState(s3)
    
    // 转移关系
    system.AddTransition("s0", "s1")
    system.AddTransition("s1", "s2")
    system.AddTransition("s2", "s3")
    system.AddTransition("s3", "s3") // 终止状态
    
    // 设置赋值
    system.SetValuation("x_positive", "s0", true)
    system.SetValuation("x_positive", "s1", true)
    system.SetValuation("x_positive", "s2", true)
    system.SetValuation("x_positive", "s3", false)
    
    system.SetValuation("terminated", "s0", false)
    system.SetValuation("terminated", "s1", false)
    system.SetValuation("terminated", "s2", false)
    system.SetValuation("terminated", "s3", true)
    
    // 验证性质
    xPositive := NewAtomFormula("x_positive")
    terminated := NewAtomFormula("terminated")
    
    // 性质1: 程序最终会终止
    property1 := NewFinally(terminated)
    
    // 性质2: 在终止之前，x始终为正数
    property2 := NewUntil(xPositive, terminated)
    
    fmt.Println("程序验证示例")
    fmt.Println("性质1: 程序最终会终止")
    for stateID := range system.States {
        value := system.Evaluate(property1, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
    
    fmt.Println("性质2: 在终止之前，x始终为正数")
    for stateID := range system.States {
        value := system.Evaluate(property2, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
}
```

### 5.2 协议验证

```go
// ProtocolVerification 协议验证示例
func ProtocolVerification() {
    // 验证互斥协议：两个进程不能同时进入临界区
    
    system := NewTransitionSystem()
    
    // 状态：idle-idle, trying-idle, critical-idle, idle-trying, idle-critical
    s1 := NewState("s1", "idle-idle")
    s2 := NewState("s2", "trying-idle")
    s3 := NewState("s3", "critical-idle")
    s4 := NewState("s4", "idle-trying")
    s5 := NewState("s5", "idle-critical")
    
    system.AddState(s1)
    system.AddState(s2)
    system.AddState(s3)
    system.AddState(s4)
    system.AddState(s5)
    
    // 转移关系
    system.AddTransition("s1", "s2") // 进程1开始尝试
    system.AddTransition("s1", "s4") // 进程2开始尝试
    system.AddTransition("s2", "s3") // 进程1进入临界区
    system.AddTransition("s4", "s5") // 进程2进入临界区
    system.AddTransition("s3", "s1") // 进程1离开临界区
    system.AddTransition("s5", "s1") // 进程2离开临界区
    
    // 设置赋值
    system.SetValuation("p1_critical", "s1", false)
    system.SetValuation("p1_critical", "s2", false)
    system.SetValuation("p1_critical", "s3", true)
    system.SetValuation("p1_critical", "s4", false)
    system.SetValuation("p1_critical", "s5", false)
    
    system.SetValuation("p2_critical", "s1", false)
    system.SetValuation("p2_critical", "s2", false)
    system.SetValuation("p2_critical", "s3", false)
    system.SetValuation("p2_critical", "s4", false)
    system.SetValuation("p2_critical", "s5", true)
    
    // 验证性质
    p1Critical := NewAtomFormula("p1_critical")
    p2Critical := NewAtomFormula("p2_critical")
    
    // 性质1: 两个进程永远不会同时进入临界区
    bothCritical := NewBinaryFormula("∧", p1Critical, p2Critical)
    property1 := NewGlobally(NewNegation(bothCritical))
    
    // 性质2: 每个进程最终都能进入临界区
    property2 := NewBinaryFormula("∧",
        NewFinally(p1Critical),
        NewFinally(p2Critical))
    
    fmt.Println("协议验证示例")
    fmt.Println("性质1: 两个进程永远不会同时进入临界区")
    for stateID := range system.States {
        value := system.Evaluate(property1, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
    
    fmt.Println("性质2: 每个进程最终都能进入临界区")
    for stateID := range system.States {
        value := system.Evaluate(property2, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
}
```

### 5.3 实时系统

```go
// RealTimeSystem 实时系统示例
func RealTimeSystem() {
    // 验证实时系统：任务必须在截止时间内完成
    
    system := NewTransitionSystem()
    
    // 状态：idle, running, completed, timeout
    s1 := NewState("s1", "idle")
    s2 := NewState("s2", "running")
    s3 := NewState("s3", "completed")
    s4 := NewState("s4", "timeout")
    
    system.AddState(s1)
    system.AddState(s2)
    system.AddState(s3)
    system.AddState(s4)
    
    // 转移关系
    system.AddTransition("s1", "s2") // 开始执行
    system.AddTransition("s2", "s3") // 正常完成
    system.AddTransition("s2", "s4") // 超时
    system.AddTransition("s3", "s1") // 回到空闲
    system.AddTransition("s4", "s1") // 回到空闲
    
    // 设置赋值
    system.SetValuation("running", "s1", false)
    system.SetValuation("running", "s2", true)
    system.SetValuation("running", "s3", false)
    system.SetValuation("running", "s4", false)
    
    system.SetValuation("completed", "s1", false)
    system.SetValuation("completed", "s2", false)
    system.SetValuation("completed", "s3", true)
    system.SetValuation("completed", "s4", false)
    
    system.SetValuation("timeout", "s1", false)
    system.SetValuation("timeout", "s2", false)
    system.SetValuation("timeout", "s3", false)
    system.SetValuation("timeout", "s4", true)
    
    // 验证性质
    running := NewAtomFormula("running")
    completed := NewAtomFormula("completed")
    timeout := NewAtomFormula("timeout")
    
    // 性质1: 任务最终会完成或超时
    property1 := NewFinally(NewBinaryFormula("∨", completed, timeout))
    
    // 性质2: 任务不会同时完成和超时
    both := NewBinaryFormula("∧", completed, timeout)
    property2 := NewGlobally(NewNegation(both))
    
    // 性质3: 如果任务正在运行，它最终会完成或超时
    property3 := NewBinaryFormula("→", running, NewFinally(NewBinaryFormula("∨", completed, timeout)))
    
    fmt.Println("实时系统验证示例")
    fmt.Println("性质1: 任务最终会完成或超时")
    for stateID := range system.States {
        value := system.Evaluate(property1, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
    
    fmt.Println("性质2: 任务不会同时完成和超时")
    for stateID := range system.States {
        value := system.Evaluate(property2, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
    
    fmt.Println("性质3: 如果任务正在运行，它最终会完成或超时")
    for stateID := range system.States {
        value := system.Evaluate(property3, stateID)
        fmt.Printf("状态 %s: %v\n", stateID, value)
    }
}
```

## 总结

时态逻辑是形式化验证的重要工具，提供了：

1. **时间表达能力**: 可以表达"总是"、"最终"、"直到"等时间概念
2. **多种时态系统**: LTL、CTL、CTL*等不同表达能力的系统
3. **模型检查**: 自动验证系统是否满足时态性质
4. **广泛应用**: 在程序验证、协议验证、实时系统等领域有重要应用

通过Go语言的实现，我们展示了：

- 时态逻辑公式的数据结构表示
- 转移系统的实现
- 语义解释和模型检查算法
- 程序验证、协议验证、实时系统等应用

这为后续的动态逻辑、混合逻辑等更高级的时态系统奠定了基础。
