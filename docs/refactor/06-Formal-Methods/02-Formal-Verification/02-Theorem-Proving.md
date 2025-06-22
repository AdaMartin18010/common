# 02-定理证明 (Theorem Proving)

## 概述

定理证明是形式化验证的核心技术，通过严格的逻辑推理验证数学命题的正确性。本文档介绍定理证明的基本概念、证明系统以及在Go语言中的实现。

## 目录

1. [证明系统 (Proof Systems)](#1-证明系统-proof-systems)
2. [自然演绎 (Natural Deduction)](#2-自然演绎-natural-deduction)
3. [公理化系统 (Axiomatic Systems)](#3-公理化系统-axiomatic-systems)
4. [证明策略 (Proof Strategies)](#4-证明策略-proof-strategies)
5. [Go语言实现](#5-go语言实现)

---

## 1. 证明系统 (Proof Systems)

### 1.1 基本概念

**定义 1.1.1** (证明)
证明是从公理和假设出发，通过推理规则得到结论的有限步骤序列。

**定义 1.1.2** (推理规则)
推理规则是形如 ```latex
\frac{P_1, P_2, ..., P_n}{C}
``` 的规则，表示从前提 ```latex
P_1, P_2, ..., P_n
``` 可以推出结论 ```latex
C
```。

**定义 1.1.3** (证明树)
证明树是表示证明过程的树形结构，根节点是结论，叶节点是公理或假设。

### 1.2 证明系统的性质

**定义 1.1.4** (可靠性)
如果 ```latex
\Gamma \vdash \phi
```，则 ```latex
\Gamma \models \phi
```（如果可证明，则语义有效）。

**定义 1.1.5** (完备性)
如果 ```latex
\Gamma \models \phi
```，则 ```latex
\Gamma \vdash \phi
```（如果语义有效，则可证明）。

### 1.3 Go语言实现

```go
package proof_system

import (
    "fmt"
    "strings"
)

// Formula 逻辑公式
type Formula interface {
    String() string
    IsAtomic() bool
    FreeVariables() []string
}

// AtomicFormula 原子公式
type AtomicFormula struct {
    Predicate string
    Arguments []string
}

func (af *AtomicFormula) String() string {
    if len(af.Arguments) == 0 {
        return af.Predicate
    }
    return fmt.Sprintf("%s(%s)", af.Predicate, strings.Join(af.Arguments, ", "))
}

func (af *AtomicFormula) IsAtomic() bool {
    return true
}

func (af *AtomicFormula) FreeVariables() []string {
    return af.Arguments
}

// Negation 否定
type Negation struct {
    Formula Formula
}

func (n *Negation) String() string {
    return "¬(" + n.Formula.String() + ")"
}

func (n *Negation) IsAtomic() bool {
    return false
}

func (n *Negation) FreeVariables() []string {
    return n.Formula.FreeVariables()
}

// Conjunction 合取
type Conjunction struct {
    Left, Right Formula
}

func (c *Conjunction) String() string {
    return "(" + c.Left.String() + " ∧ " + c.Right.String() + ")"
}

func (c *Conjunction) IsAtomic() bool {
    return false
}

func (c *Conjunction) FreeVariables() []string {
    vars := make(map[string]bool)
    for _, v := range c.Left.FreeVariables() {
        vars[v] = true
    }
    for _, v := range c.Right.FreeVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// Disjunction 析取
type Disjunction struct {
    Left, Right Formula
}

func (d *Disjunction) String() string {
    return "(" + d.Left.String() + " ∨ " + d.Right.String() + ")"
}

func (d *Disjunction) IsAtomic() bool {
    return false
}

func (d *Disjunction) FreeVariables() []string {
    vars := make(map[string]bool)
    for _, v := range d.Left.FreeVariables() {
        vars[v] = true
    }
    for _, v := range d.Right.FreeVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// Implication 蕴含
type Implication struct {
    Antecedent, Consequent Formula
}

func (i *Implication) String() string {
    return "(" + i.Antecedent.String() + " → " + i.Consequent.String() + ")"
}

func (i *Implication) IsAtomic() bool {
    return false
}

func (i *Implication) FreeVariables() []string {
    vars := make(map[string]bool)
    for _, v := range i.Antecedent.FreeVariables() {
        vars[v] = true
    }
    for _, v := range i.Consequent.FreeVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// ProofStep 证明步骤
type ProofStep struct {
    Number       int
    Formula      Formula
    Justification string
    Dependencies []int
    Assumptions  []Formula
}

// Proof 证明
type Proof struct {
    Steps    []*ProofStep
    Assumptions []Formula
    Conclusion Formula
}

// NewProof 创建证明
func NewProof(assumptions []Formula, conclusion Formula) *Proof {
    return &Proof{
        Steps:       make([]*ProofStep, 0),
        Assumptions: assumptions,
        Conclusion:  conclusion,
    }
}

// AddStep 添加证明步骤
func (p *Proof) AddStep(formula Formula, justification string, deps []int, assumptions []Formula) {
    step := &ProofStep{
        Number:       len(p.Steps) + 1,
        Formula:      formula,
        Justification: justification,
        Dependencies: deps,
        Assumptions:  assumptions,
    }
    p.Steps = append(p.Steps, step)
}

// PrintProof 打印证明
func (p *Proof) PrintProof() string {
    var sb strings.Builder
    
    sb.WriteString("Proof:\n")
    sb.WriteString("Assumptions:\n")
    for i, assumption := range p.Assumptions {
        sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, assumption.String()))
    }
    
    sb.WriteString("\nSteps:\n")
    for _, step := range p.Steps {
        sb.WriteString(fmt.Sprintf("%d. %s [%s", step.Number, step.Formula.String(), step.Justification))
        if len(step.Dependencies) > 0 {
            sb.WriteString(fmt.Sprintf(", %v", step.Dependencies))
        }
        sb.WriteString("]\n")
    }
    
    sb.WriteString(fmt.Sprintf("\nConclusion: %s\n", p.Conclusion.String()))
    
    return sb.String()
}
```

---

## 2. 自然演绎 (Natural Deduction)

### 2.1 引入规则

**规则 2.1.1** (∧-I: 合取引入)
```latex
\frac{\Gamma \vdash A \quad \Delta \vdash B}{\Gamma, \Delta \vdash A \land B}
```

**规则 2.1.2** (∨-I: 析取引入)
```latex
\frac{\Gamma \vdash A}{\Gamma \vdash A \lor B}
``` 和 ```latex
\frac{\Gamma \vdash B}{\Gamma \vdash A \lor B}
```

**规则 2.1.3** (→-I: 蕴含引入)
```latex
\frac{\Gamma, A \vdash B}{\Gamma \vdash A \rightarrow B}
```

### 2.2 消除规则

**规则 2.2.1** (∧-E: 合取消除)
```latex
\frac{\Gamma \vdash A \land B}{\Gamma \vdash A}
``` 和 ```latex
\frac{\Gamma \vdash A \land B}{\Gamma \vdash B}
```

**规则 2.2.2** (∨-E: 析取消除)
```latex
\frac{\Gamma \vdash A \lor B \quad \Delta, A \vdash C \quad \Sigma, B \vdash C}{\Gamma, \Delta, \Sigma \vdash C}
```

**规则 2.2.3** (→-E: 蕴含消除)
```latex
\frac{\Gamma \vdash A \rightarrow B \quad \Delta \vdash A}{\Gamma, \Delta \vdash B}
```

### 2.3 Go语言实现

```go
package natural_deduction

import (
    "fmt"
    "strings"
)

// NaturalDeduction 自然演绎系统
type NaturalDeduction struct {
    proofs []*Proof
}

// NewNaturalDeduction 创建自然演绎系统
func NewNaturalDeduction() *NaturalDeduction {
    return &NaturalDeduction{
        proofs: make([]*Proof, 0),
    }
}

// Prove 证明一个公式
func (nd *NaturalDeduction) Prove(assumptions []Formula, conclusion Formula) *Proof {
    proof := NewProof(assumptions, conclusion)
    
    // 添加假设
    for i, assumption := range assumptions {
        proof.AddStep(assumption, "Assumption", []int{}, []Formula{assumption})
    }
    
    return proof
}

// ConjunctionIntroduction 合取引入
func (nd *NaturalDeduction) ConjunctionIntroduction(proof *Proof, leftStep, rightStep int) {
    leftFormula := proof.Steps[leftStep-1].Formula
    rightFormula := proof.Steps[rightStep-1].Formula
    conjunction := &Conjunction{Left: leftFormula, Right: rightFormula}
    
    proof.AddStep(conjunction, "∧-I", []int{leftStep, rightStep}, nil)
}

// ConjunctionElimination 合取消除
func (nd *NaturalDeduction) ConjunctionElimination(proof *Proof, conjunctionStep int, left bool) {
    conjunction := proof.Steps[conjunctionStep-1].Formula
    if conj, ok := conjunction.(*Conjunction); ok {
        var result Formula
        if left {
            result = conj.Left
        } else {
            result = conj.Right
        }
        proof.AddStep(result, "∧-E", []int{conjunctionStep}, nil)
    }
}

// DisjunctionIntroduction 析取引入
func (nd *NaturalDeduction) DisjunctionIntroduction(proof *Proof, premiseStep int, otherFormula Formula, left bool) {
    premise := proof.Steps[premiseStep-1].Formula
    var disjunction Formula
    if left {
        disjunction = &Disjunction{Left: premise, Right: otherFormula}
    } else {
        disjunction = &Disjunction{Left: otherFormula, Right: premise}
    }
    
    proof.AddStep(disjunction, "∨-I", []int{premiseStep}, nil)
}

// ImplicationIntroduction 蕴含引入
func (nd *NaturalDeduction) ImplicationIntroduction(proof *Proof, assumption Formula, conclusionStep int) {
    conclusion := proof.Steps[conclusionStep-1].Formula
    implication := &Implication{Antecedent: assumption, Consequent: conclusion}
    
    // 移除假设
    proof.AddStep(implication, "→-I", []int{conclusionStep}, []Formula{assumption})
}

// ImplicationElimination 蕴含消除
func (nd *NaturalDeduction) ImplicationElimination(proof *Proof, implicationStep, antecedentStep int) {
    implication := proof.Steps[implicationStep-1].Formula
    antecedent := proof.Steps[antecedentStep-1].Formula
    
    if impl, ok := implication.(*Implication); ok {
        if impl.Antecedent.String() == antecedent.String() {
            proof.AddStep(impl.Consequent, "→-E", []int{implicationStep, antecedentStep}, nil)
        }
    }
}

// NegationIntroduction 否定引入
func (nd *NaturalDeduction) NegationIntroduction(proof *Proof, assumption Formula, contradictionStep int) {
    negation := &Negation{Formula: assumption}
    proof.AddStep(negation, "¬-I", []int{contradictionStep}, []Formula{assumption})
}

// NegationElimination 否定消除
func (nd *NaturalDeduction) NegationElimination(proof *Proof, formulaStep, negationStep int) {
    formula := proof.Steps[formulaStep-1].Formula
    negation := proof.Steps[negationStep-1].Formula
    
    if neg, ok := negation.(*Negation); ok {
        if neg.Formula.String() == formula.String() {
            // 引入矛盾
            contradiction := &AtomicFormula{Predicate: "⊥"}
            proof.AddStep(contradiction, "¬-E", []int{formulaStep, negationStep}, nil)
        }
    }
}

// Example: 证明 A ∧ B → B ∧ A
func ExampleCommutativity() {
    nd := NewNaturalDeduction()
    
    // 创建原子公式
    a := &AtomicFormula{Predicate: "A"}
    b := &AtomicFormula{Predicate: "B"}
    
    // 假设 A ∧ B
    assumption := &Conjunction{Left: a, Right: b}
    
    // 结论 B ∧ A
    conclusion := &Conjunction{Left: b, Right: a}
    
    // 证明 A ∧ B → B ∧ A
    proof := nd.Prove([]Formula{assumption}, &Implication{
        Antecedent: assumption,
        Consequent: conclusion,
    })
    
    // 证明步骤
    // 1. A ∧ B (假设)
    // 2. A (∧-E, 1)
    nd.ConjunctionElimination(proof, 1, true)
    
    // 3. B (∧-E, 1)
    nd.ConjunctionElimination(proof, 1, false)
    
    // 4. B ∧ A (∧-I, 3, 2)
    nd.ConjunctionIntroduction(proof, 3, 2)
    
    // 5. A ∧ B → B ∧ A (→-I, 1, 4)
    nd.ImplicationIntroduction(proof, assumption, 4)
    
    fmt.Println(proof.PrintProof())
}
```

---

## 3. 公理化系统 (Axiomatic Systems)

### 3.1 命题逻辑公理

**公理 3.1.1** (A1)
```latex
A \rightarrow (B \rightarrow A)
```

**公理 3.1.2** (A2)
```latex
(A \rightarrow (B \rightarrow C)) \rightarrow ((A \rightarrow B) \rightarrow (A \rightarrow C))
```

**公理 3.1.3** (A3)
```latex
(\neg A \rightarrow \neg B) \rightarrow (B \rightarrow A)
```

### 3.2 推理规则

**规则 3.2.1** (MP: 假言推理)
```latex
\frac{A \rightarrow B \quad A}{B}
```

### 3.3 Go语言实现

```go
package axiomatic_system

import (
    "fmt"
    "strings"
)

// Axiom 公理
type Axiom struct {
    Name     string
    Formula  Formula
}

// AxiomaticSystem 公理化系统
type AxiomaticSystem struct {
    Axioms []*Axiom
    Rules  []string
}

// NewAxiomaticSystem 创建公理化系统
func NewAxiomaticSystem() *AxiomaticSystem {
    system := &AxiomaticSystem{
        Axioms: make([]*Axiom, 0),
        Rules:  make([]string, 0),
    }
    
    // 添加命题逻辑公理
    system.addPropositionalAxioms()
    
    return system
}

// addPropositionalAxioms 添加命题逻辑公理
func (as *AxiomaticSystem) addPropositionalAxioms() {
    // 创建变量
    a := &AtomicFormula{Predicate: "A"}
    b := &AtomicFormula{Predicate: "B"}
    c := &AtomicFormula{Predicate: "C"}
    
    // A1: A → (B → A)
    axiom1 := &Axiom{
        Name: "A1",
        Formula: &Implication{
            Antecedent: a,
            Consequent: &Implication{Antecedent: b, Consequent: a},
        },
    }
    
    // A2: (A → (B → C)) → ((A → B) → (A → C))
    axiom2 := &Axiom{
        Name: "A2",
        Formula: &Implication{
            Antecedent: &Implication{
                Antecedent: a,
                Consequent: &Implication{Antecedent: b, Consequent: c},
            },
            Consequent: &Implication{
                Antecedent: &Implication{Antecedent: a, Consequent: b},
                Consequent: &Implication{Antecedent: a, Consequent: c},
            },
        },
    }
    
    // A3: (¬A → ¬B) → (B → A)
    axiom3 := &Axiom{
        Name: "A3",
        Formula: &Implication{
            Antecedent: &Implication{
                Antecedent: &Negation{Formula: a},
                Consequent: &Negation{Formula: b},
            },
            Consequent: &Implication{Antecedent: b, Consequent: a},
        },
    }
    
    as.Axioms = append(as.Axioms, axiom1, axiom2, axiom3)
    as.Rules = append(as.Rules, "MP")
}

// IsAxiom 检查是否为公理
func (as *AxiomaticSystem) IsAxiom(formula Formula) (bool, string) {
    for _, axiom := range as.Axioms {
        if as.formulasEqual(formula, axiom.Formula) {
            return true, axiom.Name
        }
    }
    return false, ""
}

// formulasEqual 检查两个公式是否相等
func (as *AxiomaticSystem) formulasEqual(f1, f2 Formula) bool {
    return f1.String() == f2.String()
}

// ModusPonens 假言推理
func (as *AxiomaticSystem) ModusPonens(proof *Proof, implicationStep, antecedentStep int) bool {
    implication := proof.Steps[implicationStep-1].Formula
    antecedent := proof.Steps[antecedentStep-1].Formula
    
    if impl, ok := implication.(*Implication); ok {
        if as.formulasEqual(impl.Antecedent, antecedent) {
            proof.AddStep(impl.Consequent, "MP", []int{implicationStep, antecedentStep}, nil)
            return true
        }
    }
    return false
}

// ProveWithAxioms 使用公理证明
func (as *AxiomaticSystem) ProveWithAxioms(assumptions []Formula, conclusion Formula) *Proof {
    proof := NewProof(assumptions, conclusion)
    
    // 添加假设
    for i, assumption := range assumptions {
        proof.AddStep(assumption, "Assumption", []int{}, []Formula{assumption})
    }
    
    // 添加公理
    for _, axiom := range as.Axioms {
        proof.AddStep(axiom.Formula, "Axiom "+axiom.Name, []int{}, nil)
    }
    
    return proof
}

// Example: 证明 A → A
func ExampleIdentity() {
    system := NewAxiomaticSystem()
    
    a := &AtomicFormula{Predicate: "A"}
    conclusion := &Implication{Antecedent: a, Consequent: a}
    
    proof := system.ProveWithAxioms(nil, conclusion)
    
    // 使用A1: A → (A → A)
    a1 := system.Axioms[0].Formula
    proof.AddStep(a1, "A1", []int{}, nil)
    
    // 使用A2: (A → ((A → A) → A)) → ((A → (A → A)) → (A → A))
    a2 := system.Axioms[1].Formula
    proof.AddStep(a2, "A2", []int{}, nil)
    
    // 使用MP从A1和A2推导
    system.ModusPonens(proof, len(proof.Steps), len(proof.Steps)-1)
    
    fmt.Println("Proof of A → A:")
    fmt.Println(proof.PrintProof())
}
```

---

## 4. 证明策略 (Proof Strategies)

### 4.1 前向推理

**策略 4.1.1** (前向推理)
从已知的公理和假设出发，逐步应用推理规则，直到得到目标结论。

### 4.2 后向推理

**策略 4.2.1** (后向推理)
从目标结论出发，寻找能够推导出该结论的前提，逐步回溯到已知的公理和假设。

### 4.3 证明搜索

**策略 4.3.1** (深度优先搜索)
使用深度优先搜索算法探索证明空间。

**策略 4.3.2** (广度优先搜索)
使用广度优先搜索算法探索证明空间。

### 4.4 Go语言实现

```go
package proof_strategies

import (
    "container/list"
    "fmt"
)

// ProofState 证明状态
type ProofState struct {
    Assumptions []Formula
    Goals       []Formula
    Steps       []*ProofStep
    Depth       int
}

// NewProofState 创建证明状态
func NewProofState(assumptions []Formula, goals []Formula) *ProofState {
    return &ProofState{
        Assumptions: assumptions,
        Goals:       goals,
        Steps:       make([]*ProofStep, 0),
        Depth:       0,
    }
}

// ForwardReasoning 前向推理
type ForwardReasoning struct {
    system *AxiomaticSystem
}

// NewForwardReasoning 创建前向推理器
func NewForwardReasoning(system *AxiomaticSystem) *ForwardReasoning {
    return &ForwardReasoning{system: system}
}

// Prove 前向推理证明
func (fr *ForwardReasoning) Prove(assumptions []Formula, goal Formula) *Proof {
    proof := NewProof(assumptions, goal)
    
    // 添加假设
    for i, assumption := range assumptions {
        proof.AddStep(assumption, "Assumption", []int{}, []Formula{assumption})
    }
    
    // 添加公理
    for _, axiom := range fr.system.Axioms {
        proof.AddStep(axiom.Formula, "Axiom "+axiom.Name, []int{}, nil)
    }
    
    // 前向推理
    fr.forwardStep(proof, goal)
    
    return proof
}

// forwardStep 前向推理步骤
func (fr *ForwardReasoning) forwardStep(proof *Proof, goal Formula) bool {
    // 检查是否已经证明
    for _, step := range proof.Steps {
        if step.Formula.String() == goal.String() {
            return true
        }
    }
    
    // 尝试应用推理规则
    for i := 0; i < len(proof.Steps); i++ {
        for j := i + 1; j < len(proof.Steps); j++ {
            if fr.system.ModusPonens(proof, i+1, j+1) {
                if fr.forwardStep(proof, goal) {
                    return true
                }
                // 回溯
                proof.Steps = proof.Steps[:len(proof.Steps)-1]
            }
        }
    }
    
    return false
}

// BackwardReasoning 后向推理
type BackwardReasoning struct {
    system *AxiomaticSystem
}

// NewBackwardReasoning 创建后向推理器
func NewBackwardReasoning(system *AxiomaticSystem) *BackwardReasoning {
    return &BackwardReasoning{system: system}
}

// Prove 后向推理证明
func (br *BackwardReasoning) Prove(assumptions []Formula, goal Formula) *Proof {
    proof := NewProof(assumptions, goal)
    
    // 添加假设
    for i, assumption := range assumptions {
        proof.AddStep(assumption, "Assumption", []int{}, []Formula{assumption})
    }
    
    // 添加公理
    for _, axiom := range br.system.Axioms {
        proof.AddStep(axiom.Formula, "Axiom "+axiom.Name, []int{}, nil)
    }
    
    // 后向推理
    br.backwardStep(proof, goal)
    
    return proof
}

// backwardStep 后向推理步骤
func (br *BackwardReasoning) backwardStep(proof *Proof, goal Formula) bool {
    // 检查是否已经证明
    for _, step := range proof.Steps {
        if step.Formula.String() == goal.String() {
            return true
        }
    }
    
    // 检查是否为公理
    if isAxiom, _ := br.system.IsAxiom(goal); isAxiom {
        proof.AddStep(goal, "Axiom", []int{}, nil)
        return true
    }
    
    // 尝试分解目标
    if impl, ok := goal.(*Implication); ok {
        // 目标: A → B，尝试证明 B（假设 A）
        proof.AddStep(impl.Antecedent, "Assumption", []int{}, []Formula{impl.Antecedent})
        if br.backwardStep(proof, impl.Consequent) {
            proof.AddStep(goal, "→-I", []int{len(proof.Steps)}, []Formula{impl.Antecedent})
            return true
        }
    }
    
    return false
}

// ProofSearch 证明搜索
type ProofSearch struct {
    system *AxiomaticSystem
}

// NewProofSearch 创建证明搜索器
func NewProofSearch(system *AxiomaticSystem) *ProofSearch {
    return &ProofSearch{system: system}
}

// BFS 广度优先搜索
func (ps *ProofSearch) BFS(assumptions []Formula, goal Formula) *Proof {
    queue := list.New()
    initialState := NewProofState(assumptions, []Formula{goal})
    queue.PushBack(initialState)
    
    visited := make(map[string]bool)
    
    for queue.Len() > 0 {
        element := queue.Front()
        queue.Remove(element)
        state := element.Value.(*ProofState)
        
        // 检查是否达到目标
        if ps.isGoalReached(state, goal) {
            return ps.buildProof(state)
        }
        
        // 生成后继状态
        successors := ps.generateSuccessors(state)
        for _, successor := range successors {
            stateKey := ps.stateKey(successor)
            if !visited[stateKey] {
                visited[stateKey] = true
                queue.PushBack(successor)
            }
        }
    }
    
    return nil // 未找到证明
}

// isGoalReached 检查是否达到目标
func (ps *ProofSearch) isGoalReached(state *ProofState, goal Formula) bool {
    for _, assumption := range state.Assumptions {
        if assumption.String() == goal.String() {
            return true
        }
    }
    return false
}

// generateSuccessors 生成后继状态
func (ps *ProofSearch) generateSuccessors(state *ProofState) []*ProofState {
    successors := make([]*ProofState, 0)
    
    // 应用推理规则生成后继状态
    for i := 0; i < len(state.Assumptions); i++ {
        for j := i + 1; j < len(state.Assumptions); j++ {
            // 尝试应用MP规则
            if impl, ok := state.Assumptions[i].(*Implication); ok {
                if impl.Antecedent.String() == state.Assumptions[j].String() {
                    newState := &ProofState{
                        Assumptions: make([]Formula, len(state.Assumptions)),
                        Goals:       state.Goals,
                        Steps:       state.Steps,
                        Depth:       state.Depth + 1,
                    }
                    copy(newState.Assumptions, state.Assumptions)
                    newState.Assumptions = append(newState.Assumptions, impl.Consequent)
                    successors = append(successors, newState)
                }
            }
        }
    }
    
    return successors
}

// stateKey 生成状态键
func (ps *ProofSearch) stateKey(state *ProofState) string {
    var sb strings.Builder
    for _, assumption := range state.Assumptions {
        sb.WriteString(assumption.String())
        sb.WriteString(";")
    }
    return sb.String()
}

// buildProof 构建证明
func (ps *ProofSearch) buildProof(state *ProofState) *Proof {
    proof := NewProof(nil, nil)
    for _, step := range state.Steps {
        proof.Steps = append(proof.Steps, step)
    }
    return proof
}
```

---

## 5. Go语言实现

### 5.1 完整的定理证明系统

```go
package theorem_proving_system

import (
    "fmt"
    "strings"
)

// TheoremProvingSystem 定理证明系统
type TheoremProvingSystem struct {
    axiomaticSystem *AxiomaticSystem
    naturalDeduction *NaturalDeduction
    forwardReasoning *ForwardReasoning
    backwardReasoning *BackwardReasoning
    proofSearch *ProofSearch
}

// NewTheoremProvingSystem 创建定理证明系统
func NewTheoremProvingSystem() *TheoremProvingSystem {
    axiomaticSystem := NewAxiomaticSystem()
    
    return &TheoremProvingSystem{
        axiomaticSystem: axiomaticSystem,
        naturalDeduction: NewNaturalDeduction(),
        forwardReasoning: NewForwardReasoning(axiomaticSystem),
        backwardReasoning: NewBackwardReasoning(axiomaticSystem),
        proofSearch: NewProofSearch(axiomaticSystem),
    }
}

// Prove 证明定理
func (tps *TheoremProvingSystem) Prove(assumptions []Formula, conclusion Formula, method string) *Proof {
    switch method {
    case "natural":
        return tps.naturalDeduction.Prove(assumptions, conclusion)
    case "axiomatic":
        return tps.axiomaticSystem.ProveWithAxioms(assumptions, conclusion)
    case "forward":
        return tps.forwardReasoning.Prove(assumptions, conclusion)
    case "backward":
        return tps.backwardReasoning.Prove(assumptions, conclusion)
    case "search":
        return tps.proofSearch.BFS(assumptions, conclusion)
    default:
        return tps.naturalDeduction.Prove(assumptions, conclusion)
    }
}

// VerifyProof 验证证明
func (tps *TheoremProvingSystem) VerifyProof(proof *Proof) *VerificationResult {
    result := &VerificationResult{
        Valid: true,
        Errors: make([]string, 0),
    }
    
    for i, step := range proof.Steps {
        if !tps.verifyStep(step, proof.Steps[:i+1]) {
            result.Valid = false
            result.Errors = append(result.Errors, 
                fmt.Sprintf("Step %d: Invalid justification", step.Number))
        }
    }
    
    return result
}

// verifyStep 验证证明步骤
func (tps *TheoremProvingSystem) verifyStep(step *ProofStep, previousSteps []*ProofStep) bool {
    switch step.Justification {
    case "Assumption":
        return true
    case "Axiom":
        _, axiomName := tps.axiomaticSystem.IsAxiom(step.Formula)
        return axiomName != ""
    case "MP":
        if len(step.Dependencies) != 2 {
            return false
        }
        i, j := step.Dependencies[0]-1, step.Dependencies[1]-1
        if i >= len(previousSteps) || j >= len(previousSteps) {
            return false
        }
        
        implication := previousSteps[i].Formula
        antecedent := previousSteps[j].Formula
        
        if impl, ok := implication.(*Implication); ok {
            return impl.Antecedent.String() == antecedent.String() &&
                   impl.Consequent.String() == step.Formula.String()
        }
        return false
    default:
        return true
    }
}

// VerificationResult 验证结果
type VerificationResult struct {
    Valid  bool
    Errors []string
}

// String 字符串表示
func (vr *VerificationResult) String() string {
    var sb strings.Builder
    
    if vr.Valid {
        sb.WriteString("Proof is valid.\n")
    } else {
        sb.WriteString("Proof is invalid.\n")
        sb.WriteString("Errors:\n")
        for _, error := range vr.Errors {
            sb.WriteString("  - " + error + "\n")
        }
    }
    
    return sb.String()
}

// Example: 证明德摩根律
func ExampleDeMorgan() {
    system := NewTheoremProvingSystem()
    
    // 创建原子公式
    a := &AtomicFormula{Predicate: "A"}
    b := &AtomicFormula{Predicate: "B"}
    
    // 证明 ¬(A ∧ B) → (¬A ∨ ¬B)
    leftSide := &Negation{Formula: &Conjunction{Left: a, Right: b}}
    rightSide := &Disjunction{
        Left:  &Negation{Formula: a},
        Right: &Negation{Formula: b},
    }
    
    conclusion := &Implication{Antecedent: leftSide, Consequent: rightSide}
    
    proof := system.Prove(nil, conclusion, "natural")
    
    fmt.Println("Proof of De Morgan's Law:")
    fmt.Println(proof.PrintProof())
    
    // 验证证明
    result := system.VerifyProof(proof)
    fmt.Println(result)
}

// Example: 证明排中律
func ExampleLawOfExcludedMiddle() {
    system := NewTheoremProvingSystem()
    
    a := &AtomicFormula{Predicate: "A"}
    conclusion := &Disjunction{Left: a, Right: &Negation{Formula: a}}
    
    proof := system.Prove(nil, conclusion, "axiomatic")
    
    fmt.Println("Proof of Law of Excluded Middle:")
    fmt.Println(proof.PrintProof())
    
    result := system.VerifyProof(proof)
    fmt.Println(result)
}

// Example: 证明双重否定
func ExampleDoubleNegation() {
    system := NewTheoremProvingSystem()
    
    a := &AtomicFormula{Predicate: "A"}
    conclusion := &Implication{
        Antecedent: &Negation{Formula: &Negation{Formula: a}},
        Consequent: a,
    }
    
    proof := system.Prove(nil, conclusion, "forward")
    
    fmt.Println("Proof of Double Negation:")
    fmt.Println(proof.PrintProof())
    
    result := system.VerifyProof(proof)
    fmt.Println(result)
}
```

---

## 总结

本文档介绍了定理证明的基本概念和实现：

1. **证明系统** - 基本概念、推理规则、证明树
2. **自然演绎** - 引入规则、消除规则、证明构造
3. **公理化系统** - 公理、推理规则、形式化证明
4. **证明策略** - 前向推理、后向推理、证明搜索
5. **Go语言实现** - 完整的定理证明系统

定理证明是形式化验证的重要技术，广泛应用于数学证明、程序验证、人工智能等领域。

---

**相关链接**:

- [01-模型检查 (Model Checking)](01-Model-Checking.md)
- [03-静态分析 (Static Analysis)](03-Static-Analysis.md)
- [04-动态分析 (Dynamic Analysis)](04-Dynamic-Analysis.md)
- [01-数学基础 (Mathematical Foundation)](../01-Mathematical-Foundation/README.md)
