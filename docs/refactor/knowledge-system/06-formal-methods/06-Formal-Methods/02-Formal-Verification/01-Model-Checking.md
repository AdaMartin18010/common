# 01-模型检查 (Model Checking)

## 概述

模型检查是一种自动化的形式化验证技术，用于验证系统是否满足给定的规范。本文档介绍模型检查的基本概念、算法以及在Go语言中的实现。

## 目录

- [01-模型检查 (Model Checking)](#01-模型检查-model-checking)
  - [概述](#概述)
  - [目录](#目录)
  - [1. 状态机模型 (State Machine Models)](#1-状态机模型-state-machine-models)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 Go语言实现](#12-go语言实现)
  - [2. 时态逻辑 (Temporal Logic)](#2-时态逻辑-temporal-logic)
    - [2.1 线性时态逻辑 (LTL)](#21-线性时态逻辑-ltl)
    - [2.2 计算树逻辑 (CTL)](#22-计算树逻辑-ctl)
    - [2.3 Go语言实现](#23-go语言实现)
  - [3. 模型检查算法 (Model Checking Algorithms)](#3-模型检查算法-model-checking-algorithms)
    - [3.1 LTL模型检查](#31-ltl模型检查)
    - [3.2 CTL模型检查](#32-ctl模型检查)
    - [3.3 Go语言实现](#33-go语言实现)
  - [4. 符号模型检查 (Symbolic Model Checking)](#4-符号模型检查-symbolic-model-checking)
    - [4.1 二元决策图 (BDD)](#41-二元决策图-bdd)
    - [4.2 Go语言实现](#42-go语言实现)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 完整的模型检查系统](#51-完整的模型检查系统)
  - [总结](#总结)

---

## 1. 状态机模型 (State Machine Models)

### 1.1 基本概念

**定义 1.1.1** (状态机)
状态机是一个五元组 ```latex
M = (S, S_0, \Sigma, \delta, F)
```，其中：

- ```latex
S
``` 是状态集
- ```latex
S_0 \subseteq S
``` 是初始状态集
- ```latex
\Sigma
``` 是输入字母表
- ```latex
\delta: S \times \Sigma \rightarrow 2^S
``` 是转移函数
- ```latex
F \subseteq S
``` 是接受状态集

**定义 1.1.2** (Kripke结构)
Kripke结构是一个四元组 ```latex
K = (S, S_0, R, L)
```，其中：

- ```latex
S
``` 是状态集
- ```latex
S_0 \subseteq S
``` 是初始状态集
- ```latex
R \subseteq S \times S
``` 是转移关系
- ```latex
L: S \rightarrow 2^{AP}
``` 是标记函数，```latex
AP
``` 是原子命题集

### 1.2 Go语言实现

```go
package state_machine

import (
    "fmt"
    "strings"
)

// State 状态
type State struct {
    ID   string
    Data map[string]interface{}
}

// Transition 转移
type Transition struct {
    From     string
    To       string
    Label    string
    Guard    func(map[string]interface{}) bool
    Action   func(map[string]interface{}) map[string]interface{}
}

// StateMachine 状态机
type StateMachine struct {
    States       map[string]*State
    InitialState string
    Transitions  []*Transition
    CurrentState string
    Variables    map[string]interface{}
}

// NewStateMachine 创建状态机
func NewStateMachine(initialState string) *StateMachine {
    return &StateMachine{
        States:       make(map[string]*State),
        InitialState: initialState,
        Transitions:  make([]*Transition, 0),
        CurrentState: initialState,
        Variables:    make(map[string]interface{}),
    }
}

// AddState 添加状态
func (sm *StateMachine) AddState(id string, data map[string]interface{}) {
    sm.States[id] = &State{
        ID:   id,
        Data: data,
    }
}

// AddTransition 添加转移
func (sm *StateMachine) AddTransition(from, to, label string, 
    guard func(map[string]interface{}) bool,
    action func(map[string]interface{}) map[string]interface{}) {
    
    transition := &Transition{
        From:   from,
        To:     to,
        Label:   label,
        Guard:   guard,
        Action:  action,
    }
    sm.Transitions = append(sm.Transitions, transition)
}

// Step 执行一步
func (sm *StateMachine) Step(input string) error {
    currentState := sm.CurrentState
    
    // 查找可用的转移
    for _, transition := range sm.Transitions {
        if transition.From == currentState && transition.Label == input {
            // 检查守卫条件
            if transition.Guard != nil && !transition.Guard(sm.Variables) {
                continue
            }
            
            // 执行动作
            if transition.Action != nil {
                sm.Variables = transition.Action(sm.Variables)
            }
            
            // 更新状态
            sm.CurrentState = transition.To
            return nil
        }
    }
    
    return fmt.Errorf("no valid transition from state %s with input %s", currentState, input)
}

// GetCurrentState 获取当前状态
func (sm *StateMachine) GetCurrentState() *State {
    return sm.States[sm.CurrentState]
}

// Reset 重置状态机
func (sm *StateMachine) Reset() {
    sm.CurrentState = sm.InitialState
    sm.Variables = make(map[string]interface{})
}

// KripkeStructure Kripke结构
type KripkeStructure struct {
    States       map[string]*State
    InitialState string
    Transitions  map[string][]string
    Labels       map[string][]string // 原子命题
}

// NewKripkeStructure 创建Kripke结构
func NewKripkeStructure(initialState string) *KripkeStructure {
    return &KripkeStructure{
        States:       make(map[string]*State),
        InitialState: initialState,
        Transitions:  make(map[string][]string),
        Labels:       make(map[string][]string),
    }
}

// AddState 添加状态
func (ks *KripkeStructure) AddState(id string, labels []string) {
    ks.States[id] = &State{ID: id}
    ks.Labels[id] = labels
    if ks.Transitions[id] == nil {
        ks.Transitions[id] = make([]string, 0)
    }
}

// AddTransition 添加转移
func (ks *KripkeStructure) AddTransition(from, to string) {
    if ks.Transitions[from] == nil {
        ks.Transitions[from] = make([]string, 0)
    }
    ks.Transitions[from] = append(ks.Transitions[from], to)
}

// GetSuccessors 获取后继状态
func (ks *KripkeStructure) GetSuccessors(state string) []string {
    return ks.Transitions[state]
}

// GetLabels 获取状态标记
func (ks *KripkeStructure) GetLabels(state string) []string {
    return ks.Labels[state]
}
```

---

## 2. 时态逻辑 (Temporal Logic)

### 2.1 线性时态逻辑 (LTL)

**定义 2.1.1** (LTL语法)
LTL公式的语法：
```latex
\phi ::= p \mid \neg \phi \mid \phi \land \phi \mid \phi \lor \phi \mid \phi \rightarrow \phi \mid X \phi \mid F \phi \mid G \phi \mid \phi U \phi \mid \phi R \phi
```

其中：

- ```latex
X \phi
```: 下一个状态满足 ```latex
\phi
```
- ```latex
F \phi
```: 将来某个状态满足 ```latex
\phi
```
- ```latex
G \phi
```: 所有将来状态都满足 ```latex
\phi
```
- ```latex
\phi U \psi
```: ```latex
\phi
``` 直到 ```latex
\psi
``` 成立
- ```latex
\phi R \psi
```: ```latex
\phi
``` 释放 ```latex
\psi
```

### 2.2 计算树逻辑 (CTL)

**定义 2.2.1** (CTL语法)
CTL公式的语法：
```latex
\phi ::= p \mid \neg \phi \mid \phi \land \phi \mid \phi \lor \phi \mid \phi \rightarrow \phi \mid EX \phi \mid AX \phi \mid EF \phi \mid AF \phi \mid EG \phi \mid AG \phi \mid E[\phi U \psi] \mid A[\phi U \psi]
```

其中：

- ```latex
EX \phi
```: 存在一个后继状态满足 ```latex
\phi
```
- ```latex
AX \phi
```: 所有后继状态都满足 ```latex
\phi
```
- ```latex
EF \phi
```: 存在一条路径，将来某个状态满足 ```latex
\phi
```
- ```latex
AF \phi
```: 所有路径，将来某个状态都满足 ```latex
\phi
```

### 2.3 Go语言实现

```go
package temporal_logic

import (
    "fmt"
    "strings"
)

// LTLFormula LTL公式
type LTLFormula interface {
    Evaluate(path []map[string]bool, position int) bool
    String() string
}

// AtomicProposition 原子命题
type AtomicProposition struct {
    Name string
}

func (ap *AtomicProposition) Evaluate(path []map[string]bool, position int) bool {
    if position >= len(path) {
        return false
    }
    return path[position][ap.Name]
}

func (ap *AtomicProposition) String() string {
    return ap.Name
}

// Negation 否定
type Negation struct {
    Formula LTLFormula
}

func (n *Negation) Evaluate(path []map[string]bool, position int) bool {
    return !n.Formula.Evaluate(path, position)
}

func (n *Negation) String() string {
    return "¬(" + n.Formula.String() + ")"
}

// Conjunction 合取
type Conjunction struct {
    Left, Right LTLFormula
}

func (c *Conjunction) Evaluate(path []map[string]bool, position int) bool {
    return c.Left.Evaluate(path, position) && c.Right.Evaluate(path, position)
}

func (c *Conjunction) String() string {
    return "(" + c.Left.String() + " ∧ " + c.Right.String() + ")"
}

// Disjunction 析取
type Disjunction struct {
    Left, Right LTLFormula
}

func (d *Disjunction) Evaluate(path []map[string]bool, position int) bool {
    return d.Left.Evaluate(path, position) || d.Right.Evaluate(path, position)
}

func (d *Disjunction) String() string {
    return "(" + d.Left.String() + " ∨ " + d.Right.String() + ")"
}

// Next 下一个
type Next struct {
    Formula LTLFormula
}

func (n *Next) Evaluate(path []map[string]bool, position int) bool {
    return n.Formula.Evaluate(path, position+1)
}

func (n *Next) String() string {
    return "X(" + n.Formula.String() + ")"
}

// Finally 将来
type Finally struct {
    Formula LTLFormula
}

func (f *Finally) Evaluate(path []map[string]bool, position int) bool {
    for i := position; i < len(path); i++ {
        if f.Formula.Evaluate(path, i) {
            return true
        }
    }
    return false
}

func (f *Finally) String() string {
    return "F(" + f.Formula.String() + ")"
}

// Globally 全局
type Globally struct {
    Formula LTLFormula
}

func (g *Globally) Evaluate(path []map[string]bool, position int) bool {
    for i := position; i < len(path); i++ {
        if !g.Formula.Evaluate(path, i) {
            return false
        }
    }
    return true
}

func (g *Globally) String() string {
    return "G(" + g.Formula.String() + ")"
}

// Until 直到
type Until struct {
    Left, Right LTLFormula
}

func (u *Until) Evaluate(path []map[string]bool, position int) bool {
    for i := position; i < len(path); i++ {
        if u.Right.Evaluate(path, i) {
            return true
        }
        if !u.Left.Evaluate(path, i) {
            return false
        }
    }
    return false
}

func (u *Until) String() string {
    return "(" + u.Left.String() + " U " + u.Right.String() + ")"
}

// CTLFormula CTL公式
type CTLFormula interface {
    Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool
    String() string
}

// CTLAtomic 原子命题
type CTLAtomic struct {
    Name string
}

func (ca *CTLAtomic) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    labels := ks.GetLabels(state)
    for _, label := range labels {
        if label == ca.Name {
            return true
        }
    }
    return false
}

func (ca *CTLAtomic) String() string {
    return ca.Name
}

// CTLExistsNext 存在下一个
type CTLExistsNext struct {
    Formula CTLFormula
}

func (en *CTLExistsNext) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    successors := ks.GetSuccessors(state)
    for _, successor := range successors {
        if en.Formula.Evaluate(ks, successor, visited) {
            return true
        }
    }
    return false
}

func (en *CTLExistsNext) String() string {
    return "EX(" + en.Formula.String() + ")"
}

// CTLForAllNext 所有下一个
type CTLForAllNext struct {
    Formula CTLFormula
}

func (an *CTLForAllNext) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    successors := ks.GetSuccessors(state)
    for _, successor := range successors {
        if !an.Formula.Evaluate(ks, successor, visited) {
            return false
        }
    }
    return true
}

func (an *CTLForAllNext) String() string {
    return "AX(" + an.Formula.String() + ")"
}

// CTLExistsFinally 存在将来
type CTLExistsFinally struct {
    Formula CTLFormula
}

func (ef *CTLExistsFinally) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    if visited[state] {
        return false // 避免循环
    }
    
    visited[state] = true
    defer delete(visited, state)
    
    if ef.Formula.Evaluate(ks, state, visited) {
        return true
    }
    
    successors := ks.GetSuccessors(state)
    for _, successor := range successors {
        if ef.Evaluate(ks, successor, visited) {
            return true
        }
    }
    
    return false
}

func (ef *CTLExistsFinally) String() string {
    return "EF(" + ef.Formula.String() + ")"
}

// CTLForAllFinally 所有将来
type CTLForAllFinally struct {
    Formula CTLFormula
}

func (af *CTLForAllFinally) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    // 使用EG¬φ的否定
    notFormula := &CTLNegation{Formula: af.Formula}
    egNot := &CTLExistsGlobally{Formula: notFormula}
    return !egNot.Evaluate(ks, state, visited)
}

func (af *CTLForAllFinally) String() string {
    return "AF(" + af.Formula.String() + ")"
}

// CTLNegation 否定
type CTLNegation struct {
    Formula CTLFormula
}

func (cn *CTLNegation) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    return !cn.Formula.Evaluate(ks, state, visited)
}

func (cn *CTLNegation) String() string {
    return "¬(" + cn.Formula.String() + ")"
}

// CTLExistsGlobally 存在全局
type CTLExistsGlobally struct {
    Formula CTLFormula
}

func (eg *CTLExistsGlobally) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    if visited[state] {
        return true // 循环路径满足Gφ
    }
    
    visited[state] = true
    defer delete(visited, state)
    
    if !eg.Formula.Evaluate(ks, state, visited) {
        return false
    }
    
    successors := ks.GetSuccessors(state)
    for _, successor := range successors {
        if !eg.Evaluate(ks, successor, visited) {
            return false
        }
    }
    
    return true
}

func (eg *CTLExistsGlobally) String() string {
    return "EG(" + eg.Formula.String() + ")"
}
```

---

## 3. 模型检查算法 (Model Checking Algorithms)

### 3.1 LTL模型检查

**算法 3.1.1** (LTL模型检查)

1. 将LTL公式转换为Büchi自动机
2. 将系统模型转换为Büchi自动机
3. 计算两个自动机的乘积
4. 检查乘积自动机是否为空

### 3.2 CTL模型检查

**算法 3.2.1** (CTL模型检查)
使用标记算法，递归计算满足公式的状态集：

```python
def model_check(ks, formula):
    if formula is atomic:
        return {s | formula ∈ L(s)}
    elif formula is ¬φ:
        return S - model_check(ks, φ)
    elif formula is φ ∧ ψ:
        return model_check(ks, φ) ∩ model_check(ks, ψ)
    elif formula is EX φ:
        return {s | ∃t ∈ R(s) : t ∈ model_check(ks, φ)}
    elif formula is EG φ:
        return greatest_fixed_point(λX. model_check(ks, φ) ∩ pre(X))
    # ... 其他情况
```

### 3.3 Go语言实现

```go
package model_checking

import (
    "fmt"
    "strings"
)

// ModelChecker 模型检查器
type ModelChecker struct {
    ks *KripkeStructure
}

// NewModelChecker 创建模型检查器
func NewModelChecker(ks *KripkeStructure) *ModelChecker {
    return &ModelChecker{ks: ks}
}

// CheckCTL 检查CTL公式
func (mc *ModelChecker) CheckCTL(formula CTLFormula) map[string]bool {
    visited := make(map[string]bool)
    result := make(map[string]bool)
    
    for state := range mc.ks.States {
        result[state] = formula.Evaluate(mc.ks, state, visited)
    }
    
    return result
}

// CheckLTL 检查LTL公式（简化版本）
func (mc *ModelChecker) CheckLTL(formula LTLFormula) bool {
    // 生成所有可能的路径
    paths := mc.generatePaths()
    
    // 检查每条路径
    for _, path := range paths {
        if !formula.Evaluate(path, 0) {
            return false
        }
    }
    
    return true
}

// generatePaths 生成所有可能的路径
func (mc *ModelChecker) generatePaths() [][]map[string]bool {
    var paths [][]map[string]bool
    maxDepth := 10 // 限制路径深度
    
    var dfs func(state string, path []map[string]bool, depth int)
    dfs = func(state string, path []map[string]bool, depth int) {
        if depth >= maxDepth {
            return
        }
        
        // 添加当前状态到路径
        labels := mc.ks.GetLabels(state)
        stateLabels := make(map[string]bool)
        for _, label := range labels {
            stateLabels[label] = true
        }
        
        newPath := append(path, stateLabels)
        paths = append(paths, newPath)
        
        // 继续搜索后继状态
        successors := mc.ks.GetSuccessors(state)
        for _, successor := range successors {
            dfs(successor, newPath, depth+1)
        }
    }
    
    dfs(mc.ks.InitialState, nil, 0)
    return paths
}

// FixedPoint 不动点计算
type FixedPoint struct {
    ks *KripkeStructure
}

// GreatestFixedPoint 最大不动点
func (fp *FixedPoint) GreatestFixedPoint(initial map[string]bool, 
    transformer func(map[string]bool) map[string]bool) map[string]bool {
    
    current := initial
    for {
        next := transformer(current)
        
        // 检查是否收敛
        if fp.setsEqual(current, next) {
            break
        }
        
        current = next
    }
    
    return current
}

// LeastFixedPoint 最小不动点
func (fp *FixedPoint) LeastFixedPoint(initial map[string]bool,
    transformer func(map[string]bool) map[string]bool) map[string]bool {
    
    current := initial
    for {
        next := transformer(current)
        
        // 检查是否收敛
        if fp.setsEqual(current, next) {
            break
        }
        
        current = next
    }
    
    return current
}

// setsEqual 检查两个集合是否相等
func (fp *FixedPoint) setsEqual(set1, set2 map[string]bool) bool {
    if len(set1) != len(set2) {
        return false
    }
    
    for key := range set1 {
        if !set2[key] {
            return false
        }
    }
    
    return true
}

// Predecessor 前驱状态计算
func (fp *FixedPoint) Predecessor(states map[string]bool) map[string]bool {
    predecessors := make(map[string]bool)
    
    for state := range fp.ks.States {
        successors := fp.ks.GetSuccessors(state)
        for _, successor := range successors {
            if states[successor] {
                predecessors[state] = true
                break
            }
        }
    }
    
    return predecessors
}

// Example: 验证互斥协议
func ExampleMutualExclusion() {
    // 创建Kripke结构
    ks := NewKripkeStructure("idle")
    
    // 添加状态
    ks.AddState("idle", []string{"idle"})
    ks.AddState("trying1", []string{"trying1"})
    ks.AddState("trying2", []string{"trying2"})
    ks.AddState("critical1", []string{"critical1"})
    ks.AddState("critical2", []string{"critical2"})
    
    // 添加转移
    ks.AddTransition("idle", "trying1")
    ks.AddTransition("idle", "trying2")
    ks.AddTransition("trying1", "critical1")
    ks.AddTransition("trying2", "critical2")
    ks.AddTransition("critical1", "idle")
    ks.AddTransition("critical2", "idle")
    
    // 创建模型检查器
    mc := NewModelChecker(ks)
    
    // 检查互斥性：AG¬(critical1 ∧ critical2)
    notCritical1 := &CTLNegation{Formula: &CTLAtomic{Name: "critical1"}}
    notCritical2 := &CTLNegation{Formula: &CTLAtomic{Name: "critical2"}}
    notBothCritical := &CTLConjunction{Left: notCritical1, Right: notCritical2}
    globallyNotBoth := &CTLForAllGlobally{Formula: notBothCritical}
    
    result := mc.CheckCTL(globallyNotBoth)
    
    fmt.Println("Mutual exclusion property:")
    for state, satisfied := range result {
        fmt.Printf("State %s: %v\n", state, satisfied)
    }
}

// CTLConjunction 合取
type CTLConjunction struct {
    Left, Right CTLFormula
}

func (cc *CTLConjunction) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    return cc.Left.Evaluate(ks, state, visited) && cc.Right.Evaluate(ks, state, visited)
}

func (cc *CTLConjunction) String() string {
    return "(" + cc.Left.String() + " ∧ " + cc.Right.String() + ")"
}

// CTLForAllGlobally 所有全局
type CTLForAllGlobally struct {
    Formula CTLFormula
}

func (ag *CTLForAllGlobally) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    // 使用EG¬φ的否定
    notFormula := &CTLNegation{Formula: ag.Formula}
    egNot := &CTLExistsGlobally{Formula: notFormula}
    return !egNot.Evaluate(ks, state, visited)
}

func (ag *CTLForAllGlobally) String() string {
    return "AG(" + ag.Formula.String() + ")"
}
```

---

## 4. 符号模型检查 (Symbolic Model Checking)

### 4.1 二元决策图 (BDD)

**定义 4.1.1** (BDD)
二元决策图是一种表示布尔函数的压缩数据结构。

### 4.2 Go语言实现

```go
package symbolic

import (
    "fmt"
    "math"
)

// BDDNode BDD节点
type BDDNode struct {
    Variable int
    Low      *BDDNode
    High     *BDDNode
    ID       int
}

// BDD BDD结构
type BDD struct {
    Nodes     map[int]*BDDNode
    NextID    int
    Variables int
}

// NewBDD 创建BDD
func NewBDD(variables int) *BDD {
    return &BDD{
        Nodes:     make(map[int]*BDDNode),
        NextID:    0,
        Variables: variables,
    }
}

// CreateNode 创建节点
func (bdd *BDD) CreateNode(variable int, low, high *BDDNode) *BDDNode {
    node := &BDDNode{
        Variable: variable,
        Low:      low,
        High:     high,
        ID:       bdd.NextID,
    }
    bdd.Nodes[bdd.NextID] = node
    bdd.NextID++
    return node
}

// True 真值节点
func (bdd *BDD) True() *BDDNode {
    return bdd.CreateNode(bdd.Variables, nil, nil)
}

// False 假值节点
func (bdd *BDD) False() *BDDNode {
    return bdd.CreateNode(bdd.Variables+1, nil, nil)
}

// Variable 变量节点
func (bdd *BDD) Variable(varIndex int) *BDDNode {
    return bdd.CreateNode(varIndex, bdd.False(), bdd.True())
}

// And 与操作
func (bdd *BDD) And(node1, node2 *BDDNode) *BDDNode {
    if node1 == bdd.False() || node2 == bdd.False() {
        return bdd.False()
    }
    if node1 == bdd.True() {
        return node2
    }
    if node2 == bdd.True() {
        return node1
    }
    
    if node1.Variable < node2.Variable {
        low := bdd.And(node1.Low, node2)
        high := bdd.And(node1.High, node2)
        return bdd.CreateNode(node1.Variable, low, high)
    } else if node1.Variable > node2.Variable {
        low := bdd.And(node1, node2.Low)
        high := bdd.And(node1, node2.High)
        return bdd.CreateNode(node2.Variable, low, high)
    } else {
        low := bdd.And(node1.Low, node2.Low)
        high := bdd.And(node1.High, node2.High)
        return bdd.CreateNode(node1.Variable, low, high)
    }
}

// Or 或操作
func (bdd *BDD) Or(node1, node2 *BDDNode) *BDDNode {
    if node1 == bdd.True() || node2 == bdd.True() {
        return bdd.True()
    }
    if node1 == bdd.False() {
        return node2
    }
    if node2 == bdd.False() {
        return node1
    }
    
    if node1.Variable < node2.Variable {
        low := bdd.Or(node1.Low, node2)
        high := bdd.Or(node1.High, node2)
        return bdd.CreateNode(node1.Variable, low, high)
    } else if node1.Variable > node2.Variable {
        low := bdd.Or(node1, node2.Low)
        high := bdd.Or(node1, node2.High)
        return bdd.CreateNode(node2.Variable, low, high)
    } else {
        low := bdd.Or(node1.Low, node2.Low)
        high := bdd.Or(node1.High, node2.High)
        return bdd.CreateNode(node1.Variable, low, high)
    }
}

// Not 非操作
func (bdd *BDD) Not(node *BDDNode) *BDDNode {
    if node == bdd.True() {
        return bdd.False()
    }
    if node == bdd.False() {
        return bdd.True()
    }
    
    low := bdd.Not(node.Low)
    high := bdd.Not(node.High)
    return bdd.CreateNode(node.Variable, low, high)
}

// Evaluate 评估BDD
func (bdd *BDD) Evaluate(node *BDDNode, assignment []bool) bool {
    if node == bdd.True() {
        return true
    }
    if node == bdd.False() {
        return false
    }
    
    if assignment[node.Variable] {
        return bdd.Evaluate(node.High, assignment)
    } else {
        return bdd.Evaluate(node.Low, assignment)
    }
}

// CountSatisfyingAssignments 计算满足赋值的数量
func (bdd *BDD) CountSatisfyingAssignments(node *BDDNode) int {
    if node == bdd.True() {
        return 1
    }
    if node == bdd.False() {
        return 0
    }
    
    lowCount := bdd.CountSatisfyingAssignments(node.Low)
    highCount := bdd.CountSatisfyingAssignments(node.High)
    
    return lowCount + highCount
}
```

---

## 5. Go语言实现

### 5.1 完整的模型检查系统

```go
package model_checking_system

import (
    "fmt"
    "strings"
)

// ModelCheckingSystem 模型检查系统
type ModelCheckingSystem struct {
    ks  *KripkeStructure
    bdd *BDD
}

// NewModelCheckingSystem 创建模型检查系统
func NewModelCheckingSystem(ks *KripkeStructure) *ModelCheckingSystem {
    return &ModelCheckingSystem{
        ks:  ks,
        bdd: NewBDD(10), // 假设最多10个变量
    }
}

// VerifyProperty 验证属性
func (mcs *ModelCheckingSystem) VerifyProperty(formula CTLFormula) *VerificationResult {
    result := &VerificationResult{
        Satisfied: true,
        CounterExamples: make([]string, 0),
    }
    
    // 计算满足公式的状态
    satisfiedStates := make(map[string]bool)
    for state := range mcs.ks.States {
        visited := make(map[string]bool)
        if formula.Evaluate(mcs.ks, state, visited) {
            satisfiedStates[state] = true
        }
    }
    
    // 检查初始状态是否满足
    if !satisfiedStates[mcs.ks.InitialState] {
        result.Satisfied = false
        result.CounterExamples = append(result.CounterExamples, 
            fmt.Sprintf("Initial state %s does not satisfy the property", mcs.ks.InitialState))
    }
    
    return result
}

// VerificationResult 验证结果
type VerificationResult struct {
    Satisfied       bool
    CounterExamples []string
}

// String 字符串表示
func (vr *VerificationResult) String() string {
    var sb strings.Builder
    
    if vr.Satisfied {
        sb.WriteString("Property is satisfied.\n")
    } else {
        sb.WriteString("Property is violated.\n")
        sb.WriteString("Counter examples:\n")
        for _, example := range vr.CounterExamples {
            sb.WriteString("  - " + example + "\n")
        }
    }
    
    return sb.String()
}

// Example: 验证死锁自由性
func ExampleDeadlockFreedom() {
    // 创建系统模型
    ks := NewKripkeStructure("s0")
    
    // 添加状态
    ks.AddState("s0", []string{"ready"})
    ks.AddState("s1", []string{"running"})
    ks.AddState("s2", []string{"blocked"})
    
    // 添加转移
    ks.AddTransition("s0", "s1")
    ks.AddTransition("s1", "s0")
    ks.AddTransition("s1", "s2")
    ks.AddTransition("s2", "s0")
    
    // 创建模型检查系统
    mcs := NewModelCheckingSystem(ks)
    
    // 验证死锁自由性：AG(EX true)
    trueFormula := &CTLAtomic{Name: "true"}
    exTrue := &CTLExistsNext{Formula: trueFormula}
    agExTrue := &CTLForAllGlobally{Formula: exTrue}
    
    result := mcs.VerifyProperty(agExTrue)
    fmt.Println("Deadlock freedom verification:")
    fmt.Println(result)
}

// Example: 验证安全性
func ExampleSafety() {
    // 创建系统模型
    ks := NewKripkeStructure("safe")
    
    // 添加状态
    ks.AddState("safe", []string{"safe"})
    ks.AddState("unsafe", []string{"unsafe"})
    
    // 添加转移
    ks.AddTransition("safe", "safe")
    ks.AddTransition("safe", "unsafe")
    ks.AddTransition("unsafe", "unsafe")
    
    // 创建模型检查系统
    mcs := NewModelCheckingSystem(ks)
    
    // 验证安全性：AG¬unsafe
    unsafe := &CTLAtomic{Name: "unsafe"}
    notUnsafe := &CTLNegation{Formula: unsafe}
    agNotUnsafe := &CTLForAllGlobally{Formula: notUnsafe}
    
    result := mcs.VerifyProperty(agNotUnsafe)
    fmt.Println("Safety verification:")
    fmt.Println(result)
}

// Example: 验证活性
func ExampleLiveness() {
    // 创建系统模型
    ks := NewKripkeStructure("waiting")
    
    // 添加状态
    ks.AddState("waiting", []string{"waiting"})
    ks.AddState("served", []string{"served"})
    
    // 添加转移
    ks.AddTransition("waiting", "waiting")
    ks.AddTransition("waiting", "served")
    ks.AddTransition("served", "waiting")
    
    // 创建模型检查系统
    mcs := NewModelCheckingSystem(ks)
    
    // 验证活性：AG(waiting → AF served)
    waiting := &CTLAtomic{Name: "waiting"}
    served := &CTLAtomic{Name: "served"}
    afServed := &CTLForAllFinally{Formula: served}
    implication := &CTLImplication{Left: waiting, Right: afServed}
    agImplication := &CTLForAllGlobally{Formula: implication}
    
    result := mcs.VerifyProperty(agImplication)
    fmt.Println("Liveness verification:")
    fmt.Println(result)
}

// CTLImplication 蕴含
type CTLImplication struct {
    Left, Right CTLFormula
}

func (ci *CTLImplication) Evaluate(ks *KripkeStructure, state string, visited map[string]bool) bool {
    return !ci.Left.Evaluate(ks, state, visited) || ci.Right.Evaluate(ks, state, visited)
}

func (ci *CTLImplication) String() string {
    return "(" + ci.Left.String() + " → " + ci.Right.String() + ")"
}
```

---

## 总结

本文档介绍了模型检查的基本概念和实现：

1. **状态机模型** - 状态机、Kripke结构
2. **时态逻辑** - LTL、CTL语法和语义
3. **模型检查算法** - LTL和CTL模型检查
4. **符号模型检查** - BDD数据结构
5. **Go语言实现** - 完整的模型检查系统

模型检查是形式化验证的重要技术，广泛应用于硬件验证、协议验证、软件验证等领域。

---

**相关链接**:

- [02-定理证明 (Theorem Proving)](02-Theorem-Proving.md)
- [03-静态分析 (Static Analysis)](03-Static-Analysis.md)
- [04-动态分析 (Dynamic Analysis)](04-Dynamic-Analysis.md)
- [01-数学基础 (Mathematical Foundation)](../01-Mathematical-Foundation/README.md)
