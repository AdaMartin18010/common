# 02-逻辑基础 (Logic Foundation)

## 目录

- [02-逻辑基础 (Logic Foundation)](#02-逻辑基础-logic-foundation)
  - [目录](#目录)
  - [1. 命题逻辑](#1-命题逻辑)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 形式化定义](#12-形式化定义)
    - [1.3 推理规则](#13-推理规则)
  - [2. 谓词逻辑](#2-谓词逻辑)
    - [2.1 基本概念](#21-基本概念)
    - [2.2 形式化定义](#22-形式化定义)
    - [2.3 量词](#23-量词)
  - [3. 模态逻辑](#3-模态逻辑)
    - [3.1 基本概念](#31-基本概念)
    - [3.2 可能世界语义](#32-可能世界语义)
    - [3.3 模态系统](#33-模态系统)
  - [4. 时态逻辑](#4-时态逻辑)
    - [4.1 基本概念](#41-基本概念)
    - [4.2 线性时态逻辑](#42-线性时态逻辑)
    - [4.3 分支时态逻辑](#43-分支时态逻辑)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 逻辑表达式](#51-逻辑表达式)
    - [5.2 推理引擎](#52-推理引擎)
    - [5.3 模型检查](#53-模型检查)
  - [6. 应用场景](#6-应用场景)
    - [6.1 程序验证](#61-程序验证)
    - [6.2 人工智能](#62-人工智能)
    - [6.3 数据库](#63-数据库)
    - [6.4 硬件设计](#64-硬件设计)
  - [7. 总结](#7-总结)

## 1. 命题逻辑

### 1.1 基本概念

**定义 1.1** (命题): 命题是一个具有确定真值的陈述句。

**定义 1.2** (原子命题): 原子命题是不可再分解的基本命题，用命题符号 $p, q, r, \ldots$ 表示。

**定义 1.3** (复合命题): 复合命题是由原子命题通过逻辑连接词构成的命题。

### 1.2 形式化定义

**定义 1.4** (命题逻辑语言): 命题逻辑语言 $\mathcal{L}$ 由以下部分组成：

- 命题符号集 $\mathcal{P} = \{p, q, r, \ldots\}$
- 逻辑连接词: $\neg$ (否定), $\wedge$ (合取), $\vee$ (析取), $\rightarrow$ (蕴含), $\leftrightarrow$ (等价)
- 辅助符号: $(, )$

**定义 1.5** (合式公式): 合式公式递归定义如下：

1. 每个命题符号 $p \in \mathcal{P}$ 是合式公式
2. 如果 $\phi$ 是合式公式，则 $\neg \phi$ 是合式公式
3. 如果 $\phi$ 和 $\psi$ 是合式公式，则 $(\phi \wedge \psi)$, $(\phi \vee \psi)$, $(\phi \rightarrow \psi)$, $(\phi \leftrightarrow \psi)$ 是合式公式

### 1.3 推理规则

**公理 1.1** (命题逻辑公理):

1. $\phi \rightarrow (\psi \rightarrow \phi)$
2. $(\phi \rightarrow (\psi \rightarrow \chi)) \rightarrow ((\phi \rightarrow \psi) \rightarrow (\phi \rightarrow \chi))$
3. $(\neg \phi \rightarrow \neg \psi) \rightarrow (\psi \rightarrow \phi)$

**推理规则 1.1** (分离规则): 从 $\phi$ 和 $\phi \rightarrow \psi$ 可以推出 $\psi$。

## 2. 谓词逻辑

### 2.1 基本概念

**定义 2.1** (谓词): 谓词是描述对象性质或关系的符号。

**定义 2.2** (个体): 个体是论域中的对象，用个体常项或个体变项表示。

**定义 2.3** (量词): 量词包括全称量词 $\forall$ 和存在量词 $\exists$。

### 2.2 形式化定义

**定义 2.4** (一阶逻辑语言): 一阶逻辑语言 $\mathcal{L}$ 包含：

- 个体常项集 $\mathcal{C}$
- 个体变项集 $\mathcal{V}$
- 谓词符号集 $\mathcal{P}$
- 函数符号集 $\mathcal{F}$
- 逻辑连接词和量词

**定义 2.5** (项): 项递归定义如下：

1. 每个个体常项和个体变项是项
2. 如果 $f$ 是 $n$ 元函数符号，$t_1, \ldots, t_n$ 是项，则 $f(t_1, \ldots, t_n)$ 是项

**定义 2.6** (原子公式): 如果 $P$ 是 $n$ 元谓词符号，$t_1, \ldots, t_n$ 是项，则 $P(t_1, \ldots, t_n)$ 是原子公式。

### 2.3 量词

**定义 2.7** (全称量词): $\forall x \phi$ 表示"对所有 $x$，$\phi$ 成立"。

**定义 2.8** (存在量词): $\exists x \phi$ 表示"存在 $x$，使得 $\phi$ 成立"。

**定理 2.1** (量词对偶性): $\neg \forall x \phi \equiv \exists x \neg \phi$ 和 $\neg \exists x \phi \equiv \forall x \neg \phi$。

## 3. 模态逻辑

### 3.1 基本概念

**定义 3.1** (模态算子): 模态算子包括 $\Box$ (必然) 和 $\Diamond$ (可能)。

**定义 3.2** (模态公式): 模态公式在命题逻辑基础上增加：

- 如果 $\phi$ 是模态公式，则 $\Box \phi$ 和 $\Diamond \phi$ 是模态公式

### 3.2 可能世界语义

**定义 3.3** (克里普克模型): 克里普克模型 $\mathcal{M} = (W, R, V)$ 包含：

- $W$: 可能世界集
- $R \subseteq W \times W$: 可达关系
- $V: W \times \mathcal{P} \rightarrow \{true, false\}$: 赋值函数

**定义 3.4** (模态公式的真值): 在可能世界 $w$ 中：

- $\mathcal{M}, w \models \Box \phi$ 当且仅当对所有 $v$ 使得 $wRv$，有 $\mathcal{M}, v \models \phi$
- $\mathcal{M}, w \models \Diamond \phi$ 当且仅当存在 $v$ 使得 $wRv$ 且 $\mathcal{M}, v \models \phi$

### 3.3 模态系统

**公理 3.1** (K公理): $\Box(\phi \rightarrow \psi) \rightarrow (\Box \phi \rightarrow \Box \psi)$

**公理 3.2** (T公理): $\Box \phi \rightarrow \phi$

**公理 3.3** (4公理): $\Box \phi \rightarrow \Box \Box \phi$

**公理 3.4** (5公理): $\Diamond \phi \rightarrow \Box \Diamond \phi$

## 4. 时态逻辑

### 4.1 基本概念

**定义 4.1** (时态算子): 时态算子包括：

- $G$ (总是), $F$ (将来), $X$ (下一个), $U$ (直到)

**定义 4.2** (时态公式): 时态公式在命题逻辑基础上增加：

- $G \phi$: $\phi$ 总是为真
- $F \phi$: $\phi$ 将来为真
- $X \phi$: $\phi$ 下一个时刻为真
- $\phi U \psi$: $\phi$ 为真直到 $\psi$ 为真

### 4.2 线性时态逻辑

**定义 4.3** (线性时态结构): 线性时态结构是序列 $\sigma = s_0, s_1, s_2, \ldots$，其中每个 $s_i$ 是状态。

**定义 4.4** (LTL语义): 在位置 $i$ 上：

- $\sigma, i \models G \phi$ 当且仅当对所有 $j \geq i$，$\sigma, j \models \phi$
- $\sigma, i \models F \phi$ 当且仅当存在 $j \geq i$，$\sigma, j \models \phi$
- $\sigma, i \models X \phi$ 当且仅当 $\sigma, i+1 \models \phi$
- $\sigma, i \models \phi U \psi$ 当且仅当存在 $j \geq i$ 使得 $\sigma, j \models \psi$ 且对所有 $k$ 满足 $i \leq k < j$，$\sigma, k \models \phi$

### 4.3 分支时态逻辑

**定义 4.5** (计算树逻辑): CTL公式包含路径量词 $A$ (对所有路径) 和 $E$ (存在路径)。

**定义 4.6** (CTL语义):

- $A \phi$: 在所有路径上 $\phi$ 为真
- $E \phi$: 存在路径使得 $\phi$ 为真

## 5. Go语言实现

### 5.1 逻辑表达式

```go
package logic

import (
    "fmt"
    "strings"
)

// 逻辑表达式接口
type Formula interface {
    String() string
    Evaluate(valuation map[string]bool) bool
    GetVariables() map[string]bool
}

// 原子命题
type Atom struct {
    Name string
}

func NewAtom(name string) *Atom {
    return &Atom{Name: name}
}

func (a *Atom) String() string {
    return a.Name
}

func (a *Atom) Evaluate(valuation map[string]bool) bool {
    return valuation[a.Name]
}

func (a *Atom) GetVariables() map[string]bool {
    return map[string]bool{a.Name: true}
}

// 否定
type Negation struct {
    Formula Formula
}

func NewNegation(formula Formula) *Negation {
    return &Negation{Formula: formula}
}

func (n *Negation) String() string {
    return fmt.Sprintf("¬(%s)", n.Formula.String())
}

func (n *Negation) Evaluate(valuation map[string]bool) bool {
    return !n.Formula.Evaluate(valuation)
}

func (n *Negation) GetVariables() map[string]bool {
    return n.Formula.GetVariables()
}

// 合取
type Conjunction struct {
    Left  Formula
    Right Formula
}

func NewConjunction(left, right Formula) *Conjunction {
    return &Conjunction{Left: left, Right: right}
}

func (c *Conjunction) String() string {
    return fmt.Sprintf("(%s ∧ %s)", c.Left.String(), c.Right.String())
}

func (c *Conjunction) Evaluate(valuation map[string]bool) bool {
    return c.Left.Evaluate(valuation) && c.Right.Evaluate(valuation)
}

func (c *Conjunction) GetVariables() map[string]bool {
    vars := c.Left.GetVariables()
    for v := range c.Right.GetVariables() {
        vars[v] = true
    }
    return vars
}

// 析取
type Disjunction struct {
    Left  Formula
    Right Formula
}

func NewDisjunction(left, right Formula) *Disjunction {
    return &Disjunction{Left: left, Right: right}
}

func (d *Disjunction) String() string {
    return fmt.Sprintf("(%s ∨ %s)", d.Left.String(), d.Right.String())
}

func (d *Disjunction) Evaluate(valuation map[string]bool) bool {
    return d.Left.Evaluate(valuation) || d.Right.Evaluate(valuation)
}

func (d *Disjunction) GetVariables() map[string]bool {
    vars := d.Left.GetVariables()
    for v := range d.Right.GetVariables() {
        vars[v] = true
    }
    return vars
}

// 蕴含
type Implication struct {
    Antecedent Formula
    Consequent Formula
}

func NewImplication(antecedent, consequent Formula) *Implication {
    return &Implication{
        Antecedent: antecedent,
        Consequent: consequent,
    }
}

func (i *Implication) String() string {
    return fmt.Sprintf("(%s → %s)", i.Antecedent.String(), i.Consequent.String())
}

func (i *Implication) Evaluate(valuation map[string]bool) bool {
    return !i.Antecedent.Evaluate(valuation) || i.Consequent.Evaluate(valuation)
}

func (i *Implication) GetVariables() map[string]bool {
    vars := i.Antecedent.GetVariables()
    for v := range i.Consequent.GetVariables() {
        vars[v] = true
    }
    return vars
}
```

### 5.2 推理引擎

```go
// 推理引擎
type InferenceEngine struct {
    axioms    []Formula
    theorems  []Formula
}

func NewInferenceEngine() *InferenceEngine {
    return &InferenceEngine{
        axioms:   make([]Formula, 0),
        theorems: make([]Formula, 0),
    }
}

// 添加公理
func (ie *InferenceEngine) AddAxiom(axiom Formula) {
    ie.axioms = append(ie.axioms, axiom)
}

// 分离规则
func (ie *InferenceEngine) ModusPonens(premise1, premise2 Formula) (Formula, error) {
    // 检查 premise1 是否为 premise2 → conclusion 的形式
    if imp, ok := premise2.(*Implication); ok {
        if imp.Antecedent.String() == premise1.String() {
            return imp.Consequent, nil
        }
    }
    return nil, fmt.Errorf("modus ponens cannot be applied")
}

// 证明检查
func (ie *InferenceEngine) IsProvable(formula Formula) bool {
    // 简化实现：检查是否为重言式
    return ie.IsTautology(formula)
}

// 重言式检查
func (ie *InferenceEngine) IsTautology(formula Formula) bool {
    variables := formula.GetVariables()
    varNames := make([]string, 0, len(variables))
    for v := range variables {
        varNames = append(varNames, v)
    }
    
    // 生成所有可能的赋值
    n := len(varNames)
    for i := 0; i < (1 << n); i++ {
        valuation := make(map[string]bool)
        for j, name := range varNames {
            valuation[name] = (i & (1 << j)) != 0
        }
        
        if !formula.Evaluate(valuation) {
            return false
        }
    }
    
    return true
}

// 矛盾检查
func (ie *InferenceEngine) IsContradiction(formula Formula) bool {
    variables := formula.GetVariables()
    varNames := make([]string, 0, len(variables))
    for v := range variables {
        varNames = append(varNames, v)
    }
    
    n := len(varNames)
    for i := 0; i < (1 << n); i++ {
        valuation := make(map[string]bool)
        for j, name := range varNames {
            valuation[name] = (i & (1 << j)) != 0
        }
        
        if formula.Evaluate(valuation) {
            return false
        }
    }
    
    return true
}

// 等价性检查
func (ie *InferenceEngine) AreEquivalent(formula1, formula2 Formula) bool {
    variables := make(map[string]bool)
    for v := range formula1.GetVariables() {
        variables[v] = true
    }
    for v := range formula2.GetVariables() {
        variables[v] = true
    }
    
    varNames := make([]string, 0, len(variables))
    for v := range variables {
        varNames = append(varNames, v)
    }
    
    n := len(varNames)
    for i := 0; i < (1 << n); i++ {
        valuation := make(map[string]bool)
        for j, name := range varNames {
            valuation[name] = (i & (1 << j)) != 0
        }
        
        if formula1.Evaluate(valuation) != formula2.Evaluate(valuation) {
            return false
        }
    }
    
    return true
}
```

### 5.3 模型检查

```go
// 时态逻辑公式
type TemporalFormula interface {
    Formula
    IsTemporal() bool
}

// 总是算子
type Always struct {
    Formula Formula
}

func NewAlways(formula Formula) *Always {
    return &Always{Formula: formula}
}

func (a *Always) String() string {
    return fmt.Sprintf("G(%s)", a.Formula.String())
}

func (a *Always) IsTemporal() bool {
    return true
}

func (a *Always) Evaluate(valuation map[string]bool) bool {
    // 简化实现：在当前状态下检查
    return a.Formula.Evaluate(valuation)
}

func (a *Always) GetVariables() map[string]bool {
    return a.Formula.GetVariables()
}

// 将来算子
type Eventually struct {
    Formula Formula
}

func NewEventually(formula Formula) *Eventually {
    return &Eventually{Formula: formula}
}

func (e *Eventually) String() string {
    return fmt.Sprintf("F(%s)", e.Formula.String())
}

func (e *Eventually) IsTemporal() bool {
    return true
}

func (e *Eventually) Evaluate(valuation map[string]bool) bool {
    // 简化实现：在当前状态下检查
    return e.Formula.Evaluate(valuation)
}

func (e *Eventually) GetVariables() map[string]bool {
    return e.Formula.GetVariables()
}

// 模型检查器
type ModelChecker struct {
    states     []map[string]bool
    transitions [][]int
}

func NewModelChecker() *ModelChecker {
    return &ModelChecker{
        states:      make([]map[string]bool, 0),
        transitions: make([][]int, 0),
    }
}

// 添加状态
func (mc *ModelChecker) AddState(valuation map[string]bool) int {
    state := make(map[string]bool)
    for k, v := range valuation {
        state[k] = v
    }
    mc.states = append(mc.states, state)
    mc.transitions = append(mc.transitions, make([]int, 0))
    return len(mc.states) - 1
}

// 添加转换
func (mc *ModelChecker) AddTransition(from, to int) {
    if from >= 0 && from < len(mc.transitions) {
        mc.transitions[from] = append(mc.transitions[from], to)
    }
}

// 检查LTL公式
func (mc *ModelChecker) CheckLTL(formula TemporalFormula, state int) bool {
    if state < 0 || state >= len(mc.states) {
        return false
    }
    
    switch f := formula.(type) {
    case *Always:
        return mc.checkAlways(f.Formula, state)
    case *Eventually:
        return mc.checkEventually(f.Formula, state)
    default:
        return formula.Evaluate(mc.states[state])
    }
}

// 检查总是算子
func (mc *ModelChecker) checkAlways(formula Formula, state int) bool {
    visited := make(map[int]bool)
    return mc.checkAlwaysRecursive(formula, state, visited)
}

func (mc *ModelChecker) checkAlwaysRecursive(formula Formula, state int, visited map[int]bool) bool {
    if visited[state] {
        return true // 避免循环
    }
    
    visited[state] = true
    
    // 检查当前状态
    if !formula.Evaluate(mc.states[state]) {
        return false
    }
    
    // 检查所有后继状态
    for _, nextState := range mc.transitions[state] {
        if !mc.checkAlwaysRecursive(formula, nextState, visited) {
            return false
        }
    }
    
    return true
}

// 检查将来算子
func (mc *ModelChecker) checkEventually(formula Formula, state int) bool {
    visited := make(map[int]bool)
    return mc.checkEventuallyRecursive(formula, state, visited)
}

func (mc *ModelChecker) checkEventuallyRecursive(formula Formula, state int, visited map[int]bool) bool {
    if visited[state] {
        return false // 避免循环
    }
    
    visited[state] = true
    
    // 检查当前状态
    if formula.Evaluate(mc.states[state]) {
        return true
    }
    
    // 检查后继状态
    for _, nextState := range mc.transitions[state] {
        if mc.checkEventuallyRecursive(formula, nextState, visited) {
            return true
        }
    }
    
    return false
}
```

## 6. 应用场景

### 6.1 程序验证

- 程序正确性证明
- 模型检查
- 静态分析

### 6.2 人工智能

- 知识表示
- 自动推理
- 专家系统

### 6.3 数据库

- 查询语言
- 约束检查
- 完整性验证

### 6.4 硬件设计

- 电路验证
- 时序分析
- 协议验证

## 7. 总结

逻辑基础作为计算机科学的重要理论基础，在程序验证、人工智能、数据库等领域有广泛应用。通过Go语言的实现，我们可以看到：

1. **理论实现**: 逻辑学的基本概念可以转化为高效的代码
2. **推理引擎**: 通过算法实现自动推理和证明
3. **模型检查**: 用于验证系统性质和行为

逻辑学的研究不仅有助于理解推理的本质，也为构建可靠的软件系统提供了理论基础。

---

**相关链接**:

- [01-数学基础](../01-Mathematical-Foundation/README.md)
- [03-范畴论基础](../03-Category-Theory-Foundation/README.md)
- [04-计算理论基础](../04-Computational-Theory-Foundation/README.md)
- [03-设计模式](../../03-Design-Patterns/README.md)
- [02-软件架构](../../02-Software-Architecture/README.md)
