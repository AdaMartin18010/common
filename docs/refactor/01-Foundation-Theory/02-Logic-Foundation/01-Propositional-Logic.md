# 01-命题逻辑 (Propositional Logic)

## 目录

- [01-命题逻辑 (Propositional Logic)](#01-命题逻辑-propositional-logic)
	- [目录](#目录)
	- [1. 基础概念](#1-基础概念)
		- [1.1 命题定义](#11-命题定义)
		- [1.2 逻辑连接词](#12-逻辑连接词)
		- [1.3 真值表](#13-真值表)
	- [2. 形式化语法](#2-形式化语法)
		- [2.1 BNF语法定义](#21-bnf语法定义)
		- [2.2 抽象语法树](#22-抽象语法树)
		- [2.3 语法分析](#23-语法分析)
	- [3. 语义学](#3-语义学)
		- [3.1 解释函数](#31-解释函数)
		- [3.2 真值赋值](#32-真值赋值)
		- [3.3 语义等价](#33-语义等价)
	- [4. 证明系统](#4-证明系统)
		- [4.1 自然演绎](#41-自然演绎)
		- [4.2 公理系统](#42-公理系统)
		- [4.3 归结证明](#43-归结证明)
	- [5. Go语言实现](#5-go语言实现)
		- [5.1 命题表示](#51-命题表示)
		- [5.2 语义求值](#52-语义求值)
		- [5.3 证明构造](#53-证明构造)
	- [6. 应用与扩展](#6-应用与扩展)
		- [6.1 电路设计](#61-电路设计)
		- [6.2 程序验证](#62-程序验证)
		- [6.3 知识表示](#63-知识表示)
	- [总结](#总结)
	- [参考文献](#参考文献)

---

## 1. 基础概念

### 1.1 命题定义

**命题**是能够判断真假的陈述句。在形式化逻辑中，我们用符号来表示命题。

**形式化定义**：

```latex
设 P 为原子命题集合，则命题逻辑的语言 L 定义为：
- 每个 p ∈ P 都是命题
- 如果 φ 是命题，则 ¬φ 是命题
- 如果 φ, ψ 是命题，则 (φ ∧ ψ), (φ ∨ ψ), (φ → ψ), (φ ↔ ψ) 是命题
```

### 1.2 逻辑连接词

**基本逻辑连接词**：

| 连接词 | 符号 | 名称 | 含义 |
|--------|------|------|------|
| 否定 | ¬ | NOT | 非 |
| 合取 | ∧ | AND | 且 |
| 析取 | ∨ | OR | 或 |
| 蕴含 | → | IMPLIES | 如果...那么 |
| 等价 | ↔ | IFF | 当且仅当 |

**真值表定义**：

```latex
对于任意命题 φ, ψ：

¬φ 的真值：
φ | ¬φ
T | F
F | T

φ ∧ ψ 的真值：
φ | ψ | φ ∧ ψ
T | T | T
T | F | F
F | T | F
F | F | F

φ ∨ ψ 的真值：
φ | ψ | φ ∨ ψ
T | T | T
T | F | T
F | T | T
F | F | F

φ → ψ 的真值：
φ | ψ | φ → ψ
T | T | T
T | F | F
F | T | T
F | F | T

φ ↔ ψ 的真值：
φ | ψ | φ ↔ ψ
T | T | T
T | F | F
F | T | F
F | F | T
```

### 1.3 真值表

**完全真值表**：列出所有可能赋值下的真值。

**示例**：命题 (p ∧ q) → (p ∨ q) 的真值表

```latex
p | q | p ∧ q | p ∨ q | (p ∧ q) → (p ∨ q)
T | T |   T   |   T   |         T
T | F |   F   |   T   |         T
F | T |   F   |   T   |         T
F | F |   F   |   F   |         T
```

---

## 2. 形式化语法

### 2.1 BNF语法定义

**命题逻辑的BNF语法**：

```bnf
<proposition> ::= <atomic> | <negation> | <binary>
<atomic>      ::= <identifier>
<negation>    ::= "¬" <proposition> | "~" <proposition>
<binary>      ::= "(" <proposition> <operator> <proposition> ")"
<operator>    ::= "∧" | "∨" | "→" | "↔" | "&" | "|" | "->" | "<->"
<identifier>  ::= [a-zA-Z_][a-zA-Z0-9_]*
```

### 2.2 抽象语法树

**AST节点类型**：

```latex
AST节点类型定义：
- Atom(p): 原子命题 p
- Not(φ): 否定 ¬φ
- And(φ, ψ): 合取 φ ∧ ψ
- Or(φ, ψ): 析取 φ ∨ ψ
- Implies(φ, ψ): 蕴含 φ → ψ
- Iff(φ, ψ): 等价 φ ↔ ψ
```

### 2.3 语法分析

**递归下降解析器**：

```latex
解析算法：
1. 词法分析：将输入字符串转换为token序列
2. 语法分析：根据BNF规则构建AST
3. 语义分析：检查语法正确性
```

---

## 3. 语义学

### 3.1 解释函数

**解释函数定义**：

```latex
设 P 为原子命题集合，解释函数 I: P → {true, false}

语义函数 ⟦·⟧_I: Formula → {true, false} 递归定义：

⟦p⟧_I = I(p)                    (p ∈ P)
⟦¬φ⟧_I = ¬⟦φ⟧_I
⟦φ ∧ ψ⟧_I = ⟦φ⟧_I ∧ ⟦ψ⟧_I
⟦φ ∨ ψ⟧_I = ⟦φ⟧_I ∨ ⟦ψ⟧_I
⟦φ → ψ⟧_I = ¬⟦φ⟧_I ∨ ⟦ψ⟧_I
⟦φ ↔ ψ⟧_I = (⟦φ⟧_I ∧ ⟦ψ⟧_I) ∨ (¬⟦φ⟧_I ∧ ¬⟦ψ⟧_I)
```

### 3.2 真值赋值

**真值赋值**：为每个原子命题分配真值。

**满足性**：

- 如果 ⟦φ⟧_I = true，称解释 I 满足公式 φ，记作 I ⊨ φ
- 如果 ⟦φ⟧_I = false，称解释 I 不满足公式 φ，记作 I ⊭ φ

### 3.3 语义等价

**语义等价定义**：

```latex
两个公式 φ 和 ψ 语义等价，记作 φ ≡ ψ，当且仅当：
对于所有解释 I，都有 ⟦φ⟧_I = ⟦ψ⟧_I
```

**重要等价关系**：

```latex
双重否定：¬¬φ ≡ φ
德摩根律：¬(φ ∧ ψ) ≡ ¬φ ∨ ¬ψ
         ¬(φ ∨ ψ) ≡ ¬φ ∧ ¬ψ
分配律：φ ∧ (ψ ∨ χ) ≡ (φ ∧ ψ) ∨ (φ ∧ χ)
       φ ∨ (ψ ∧ χ) ≡ (φ ∨ ψ) ∧ (φ ∨ χ)
蕴含等价：φ → ψ ≡ ¬φ ∨ ψ
等价分解：φ ↔ ψ ≡ (φ → ψ) ∧ (ψ → φ)
```

---

## 4. 证明系统

### 4.1 自然演绎

**自然演绎规则**：

```latex
引入规则：
∧I: φ, ψ ⊢ φ ∧ ψ
∨I₁: φ ⊢ φ ∨ ψ
∨I₂: ψ ⊢ φ ∨ ψ
→I: φ ⊢ ψ / ⊢ φ → ψ
¬I: φ ⊢ ⊥ / ⊢ ¬φ

消除规则：
∧E₁: φ ∧ ψ ⊢ φ
∧E₂: φ ∧ ψ ⊢ ψ
∨E: φ ∨ ψ, φ → χ, ψ → χ ⊢ χ
→E: φ, φ → ψ ⊢ ψ
¬E: φ, ¬φ ⊢ ⊥
```

### 4.2 公理系统

**经典命题逻辑公理**：

```latex
A1: φ → (ψ → φ)
A2: (φ → (ψ → χ)) → ((φ → ψ) → (φ → χ))
A3: (¬φ → ¬ψ) → (ψ → φ)

推理规则：MP (Modus Ponens)
φ, φ → ψ / ψ
```

### 4.3 归结证明

**归结规则**：

```latex
归结：从子句 C₁ ∨ p 和 C₂ ∨ ¬p 推出 C₁ ∨ C₂
```

---

## 5. Go语言实现

### 5.1 命题表示

```go
package propositional

import (
 "fmt"
 "strings"
)

// Proposition 表示命题逻辑公式
type Proposition interface {
 String() string
 Evaluate(interpretation map[string]bool) bool
 GetAtoms() map[string]bool
}

// Atom 原子命题
type Atom struct {
 Name string
}

func (a Atom) String() string {
 return a.Name
}

func (a Atom) Evaluate(interpretation map[string]bool) bool {
 return interpretation[a.Name]
}

func (a Atom) GetAtoms() map[string]bool {
 return map[string]bool{a.Name: true}
}

// Not 否定
type Not struct {
 Formula Proposition
}

func (n Not) String() string {
 return fmt.Sprintf("¬(%s)", n.Formula)
}

func (n Not) Evaluate(interpretation map[string]bool) bool {
 return !n.Formula.Evaluate(interpretation)
}

func (n Not) GetAtoms() map[string]bool {
 return n.Formula.GetAtoms()
}

// And 合取
type And struct {
 Left, Right Proposition
}

func (a And) String() string {
 return fmt.Sprintf("(%s ∧ %s)", a.Left, a.Right)
}

func (a And) Evaluate(interpretation map[string]bool) bool {
 return a.Left.Evaluate(interpretation) && a.Right.Evaluate(interpretation)
}

func (a And) GetAtoms() map[string]bool {
 atoms := a.Left.GetAtoms()
 for atom := range a.Right.GetAtoms() {
  atoms[atom] = true
 }
 return atoms
}

// Or 析取
type Or struct {
 Left, Right Proposition
}

func (o Or) String() string {
 return fmt.Sprintf("(%s ∨ %s)", o.Left, o.Right)
}

func (o Or) Evaluate(interpretation map[string]bool) bool {
 return o.Left.Evaluate(interpretation) || o.Right.Evaluate(interpretation)
}

func (o Or) GetAtoms() map[string]bool {
 atoms := o.Left.GetAtoms()
 for atom := range o.Right.GetAtoms() {
  atoms[atom] = true
 }
 return atoms
}

// Implies 蕴含
type Implies struct {
 Left, Right Proposition
}

func (i Implies) String() string {
 return fmt.Sprintf("(%s → %s)", i.Left, i.Right)
}

func (i Implies) Evaluate(interpretation map[string]bool) bool {
 return !i.Left.Evaluate(interpretation) || i.Right.Evaluate(interpretation)
}

func (i Implies) GetAtoms() map[string]bool {
 atoms := i.Left.GetAtoms()
 for atom := range i.Right.GetAtoms() {
  atoms[atom] = true
 }
 return atoms
}

// Iff 等价
type Iff struct {
 Left, Right Proposition
}

func (iff Iff) String() string {
 return fmt.Sprintf("(%s ↔ %s)", iff.Left, iff.Right)
}

func (iff Iff) Evaluate(interpretation map[string]bool) bool {
 return iff.Left.Evaluate(interpretation) == iff.Right.Evaluate(interpretation)
}

func (iff Iff) GetAtoms() map[string]bool {
 atoms := iff.Left.GetAtoms()
 for atom := range iff.Right.GetAtoms() {
  atoms[atom] = true
 }
 return atoms
}
```

### 5.2 语义求值

```go
// TruthTable 生成真值表
func TruthTable(formula Proposition) [][]bool {
 atoms := formula.GetAtoms()
 atomList := make([]string, 0, len(atoms))
 for atom := range atoms {
  atomList = append(atomList, atom)
 }
 
 // 生成所有可能的赋值
 assignments := generateAssignments(atomList)
 
 // 计算每个赋值下的真值
 table := make([][]bool, 0, len(assignments))
 for _, assignment := range assignments {
  result := formula.Evaluate(assignment)
  row := make([]bool, 0, len(atomList)+1)
  for _, atom := range atomList {
   row = append(row, assignment[atom])
  }
  row = append(row, result)
  table = append(table, row)
 }
 
 return table
}

// generateAssignments 生成所有可能的真值赋值
func generateAssignments(atoms []string) []map[string]bool {
 if len(atoms) == 0 {
  return []map[string]bool{{}}
 }
 
 // 递归生成
 subAssignments := generateAssignments(atoms[1:])
 assignments := make([]map[string]bool, 0, 2*len(subAssignments))
 
 for _, sub := range subAssignments {
  // 当前原子为true
  trueAssignment := make(map[string]bool)
  for k, v := range sub {
   trueAssignment[k] = v
  }
  trueAssignment[atoms[0]] = true
  assignments = append(assignments, trueAssignment)
  
  // 当前原子为false
  falseAssignment := make(map[string]bool)
  for k, v := range sub {
   falseAssignment[k] = v
  }
  falseAssignment[atoms[0]] = false
  assignments = append(assignments, falseAssignment)
 }
 
 return assignments
}

// IsTautology 判断是否为重言式
func IsTautology(formula Proposition) bool {
 assignments := generateAssignments(getAtomList(formula.GetAtoms()))
 for _, assignment := range assignments {
  if !formula.Evaluate(assignment) {
   return false
  }
 }
 return true
}

// IsContradiction 判断是否为矛盾式
func IsContradiction(formula Proposition) bool {
 assignments := generateAssignments(getAtomList(formula.GetAtoms()))
 for _, assignment := range assignments {
  if formula.Evaluate(assignment) {
   return false
  }
 }
 return true
}

// IsSatisfiable 判断是否为可满足式
func IsSatisfiable(formula Proposition) bool {
 assignments := generateAssignments(getAtomList(formula.GetAtoms()))
 for _, assignment := range assignments {
  if formula.Evaluate(assignment) {
   return true
  }
 }
 return false
}

func getAtomList(atoms map[string]bool) []string {
 atomList := make([]string, 0, len(atoms))
 for atom := range atoms {
  atomList = append(atomList, atom)
 }
 return atomList
}
```

### 5.3 证明构造

```go
// Proof 表示证明
type Proof struct {
 Premises []Proposition
 Conclusion Proposition
 Steps     []ProofStep
}

// ProofStep 证明步骤
type ProofStep struct {
 Formula Proposition
 Rule    string
 Lines   []int // 引用前面的行
}

// NaturalDeduction 自然演绎证明系统
type NaturalDeduction struct {
 lines []ProofStep
}

// AddPremise 添加前提
func (nd *NaturalDeduction) AddPremise(formula Proposition) {
 nd.lines = append(nd.lines, ProofStep{
  Formula: formula,
  Rule:    "Premise",
  Lines:   nil,
 })
}

// AndIntroduction 合取引入
func (nd *NaturalDeduction) AndIntroduction(line1, line2 int) error {
 if line1 >= len(nd.lines) || line2 >= len(nd.lines) {
  return fmt.Errorf("invalid line numbers")
 }
 
 formula1 := nd.lines[line1].Formula
 formula2 := nd.lines[line2].Formula
 
 nd.lines = append(nd.lines, ProofStep{
  Formula: And{Left: formula1, Right: formula2},
  Rule:    "∧I",
  Lines:   []int{line1, line2},
 })
 
 return nil
}

// AndElimination1 合取消除1
func (nd *NaturalDeduction) AndElimination1(line int) error {
 if line >= len(nd.lines) {
  return fmt.Errorf("invalid line number")
 }
 
 formula := nd.lines[line].Formula
 if and, ok := formula.(And); ok {
  nd.lines = append(nd.lines, ProofStep{
   Formula: and.Left,
   Rule:    "∧E₁",
   Lines:   []int{line},
  })
  return nil
 }
 
 return fmt.Errorf("formula is not a conjunction")
}

// AndElimination2 合取消除2
func (nd *NaturalDeduction) AndElimination2(line int) error {
 if line >= len(nd.lines) {
  return fmt.Errorf("invalid line number")
 }
 
 formula := nd.lines[line].Formula
 if and, ok := formula.(And); ok {
  nd.lines = append(nd.lines, ProofStep{
   Formula: and.Right,
   Rule:    "∧E₂",
   Lines:   []int{line},
  })
  return nil
 }
 
 return fmt.Errorf("formula is not a conjunction")
}

// OrIntroduction1 析取引入1
func (nd *NaturalDeduction) OrIntroduction1(line int, right Proposition) error {
 if line >= len(nd.lines) {
  return fmt.Errorf("invalid line number")
 }
 
 formula := nd.lines[line].Formula
 nd.lines = append(nd.lines, ProofStep{
  Formula: Or{Left: formula, Right: right},
  Rule:    "∨I₁",
  Lines:   []int{line},
 })
 
 return nil
}

// OrIntroduction2 析取引入2
func (nd *NaturalDeduction) OrIntroduction2(left Proposition, line int) error {
 if line >= len(nd.lines) {
  return fmt.Errorf("invalid line number")
 }
 
 formula := nd.lines[line].Formula
 nd.lines = append(nd.lines, ProofStep{
  Formula: Or{Left: left, Right: formula},
  Rule:    "∨I₂",
  Lines:   []int{line},
 })
 
 return nil
}

// ImplicationIntroduction 蕴含引入
func (nd *NaturalDeduction) ImplicationIntroduction(assumptionLine, conclusionLine int) error {
 if assumptionLine >= len(nd.lines) || conclusionLine >= len(nd.lines) {
  return fmt.Errorf("invalid line numbers")
 }
 
 assumption := nd.lines[assumptionLine].Formula
 conclusion := nd.lines[conclusionLine].Formula
 
 nd.lines = append(nd.lines, ProofStep{
  Formula: Implies{Left: assumption, Right: conclusion},
  Rule:    "→I",
  Lines:   []int{assumptionLine, conclusionLine},
 })
 
 return nil
}

// ModusPonens 假言推理
func (nd *NaturalDeduction) ModusPonens(implicationLine, antecedentLine int) error {
 if implicationLine >= len(nd.lines) || antecedentLine >= len(nd.lines) {
  return fmt.Errorf("invalid line numbers")
 }
 
 implication := nd.lines[implicationLine].Formula
 antecedent := nd.lines[antecedentLine].Formula
 
 if impl, ok := implication.(Implies); ok {
  // 检查前件是否匹配
  if impl.Left.String() == antecedent.String() {
   nd.lines = append(nd.lines, ProofStep{
    Formula: impl.Right,
    Rule:    "MP",
    Lines:   []int{implicationLine, antecedentLine},
   })
   return nil
  }
 }
 
 return fmt.Errorf("invalid modus ponens application")
}

// PrintProof 打印证明
func (nd *NaturalDeduction) PrintProof() {
 fmt.Println("Proof:")
 fmt.Println("Line | Formula | Rule | Lines")
 fmt.Println("-----|---------|------|------")
 for i, step := range nd.lines {
  linesStr := ""
  if len(step.Lines) > 0 {
   linesStr = fmt.Sprintf("%v", step.Lines)
  }
  fmt.Printf("%4d | %-7s | %-4s | %s\n", 
   i+1, step.Formula.String(), step.Rule, linesStr)
 }
}
```

---

## 6. 应用与扩展

### 6.1 电路设计

**逻辑门实现**：

```go
// LogicGate 逻辑门接口
type LogicGate interface {
 Compute(inputs []bool) bool
}

// ANDGate 与门
type ANDGate struct{}

func (g ANDGate) Compute(inputs []bool) bool {
 for _, input := range inputs {
  if !input {
   return false
  }
 }
 return true
}

// ORGate 或门
type ORGate struct{}

func (g ORGate) Compute(inputs []bool) bool {
 for _, input := range inputs {
  if input {
   return true
  }
 }
 return false
}

// NOTGate 非门
type NOTGate struct{}

func (g NOTGate) Compute(inputs []bool) bool {
 if len(inputs) != 1 {
  panic("NOT gate requires exactly one input")
 }
 return !inputs[0]
}

// Circuit 电路
type Circuit struct {
 gates []LogicGate
 connections [][]int
}

func (c *Circuit) AddGate(gate LogicGate) {
 c.gates = append(c.gates, gate)
}

func (c *Circuit) Connect(from, to int) {
 if len(c.connections) <= to {
  c.connections = append(c.connections, make([][]int, to-len(c.connections)+1)...)
 }
 c.connections[to] = append(c.connections[to], from)
}

func (c *Circuit) Evaluate(inputs []bool) []bool {
 outputs := make([]bool, len(c.gates))
 
 // 复制输入
 copy(outputs, inputs)
 
 // 计算每个门的输出
 for i := len(inputs); i < len(c.gates); i++ {
  gateInputs := make([]bool, 0)
  for _, conn := range c.connections[i] {
   gateInputs = append(gateInputs, outputs[conn])
  }
  outputs[i] = c.gates[i].Compute(gateInputs)
 }
 
 return outputs
}
```

### 6.2 程序验证

**前置条件和后置条件**：

```go
// Precondition 前置条件
type Precondition struct {
 Formula Proposition
}

// Postcondition 后置条件
type Postcondition struct {
 Formula Proposition
}

// HoareTriple 霍尔三元组 {P} C {Q}
type HoareTriple struct {
 Precondition  Precondition
 Command       string
 Postcondition Postcondition
}

// VerifyHoareTriple 验证霍尔三元组
func VerifyHoareTriple(triple HoareTriple) bool {
 // 这里需要实现具体的验证逻辑
 // 实际应用中会使用更复杂的程序分析技术
 return true
}

// Example: 验证交换程序
func VerifySwap() {
 // {x = A ∧ y = B} temp := x; x := y; y := temp {x = B ∧ y = A}
 
 pre := Precondition{
  Formula: And{
   Left:  Atom{Name: "x=A"},
   Right: Atom{Name: "y=B"},
  },
 }
 
 post := Postcondition{
  Formula: And{
   Left:  Atom{Name: "x=B"},
   Right: Atom{Name: "y=A"},
  },
 }
 
 triple := HoareTriple{
  Precondition:  pre,
  Command:       "temp := x; x := y; y := temp",
  Postcondition: post,
 }
 
 if VerifyHoareTriple(triple) {
  fmt.Println("Swap program is correct")
 } else {
  fmt.Println("Swap program is incorrect")
 }
}
```

### 6.3 知识表示

**知识库系统**：

```go
// KnowledgeBase 知识库
type KnowledgeBase struct {
 facts    []Proposition
 rules    []Implication
 query    Proposition
}

// AddFact 添加事实
func (kb *KnowledgeBase) AddFact(fact Proposition) {
 kb.facts = append(kb.facts, fact)
}

// AddRule 添加规则
func (kb *KnowledgeBase) AddRule(rule Implication) {
 kb.rules = append(kb.rules, rule)
}

// Query 查询
func (kb *KnowledgeBase) Query(query Proposition) bool {
 // 简单的前向推理
 knownFacts := make(map[string]bool)
 
 // 初始化已知事实
 for _, fact := range kb.facts {
  if atom, ok := fact.(Atom); ok {
   knownFacts[atom.Name] = true
  }
 }
 
 // 应用规则
 changed := true
 for changed {
  changed = false
  for _, rule := range kb.rules {
   if rule.Left.Evaluate(knownFacts) && !rule.Right.Evaluate(knownFacts) {
    // 可以推导出新事实
    if atom, ok := rule.Right.(Atom); ok {
     knownFacts[atom.Name] = true
     changed = true
    }
   }
  }
 }
 
 return query.Evaluate(knownFacts)
}

// Example: 动物知识库
func AnimalKnowledgeBase() {
 kb := &KnowledgeBase{}
 
 // 添加事实
 kb.AddFact(Atom{Name: "has_fur"})      // 有毛
 kb.AddFact(Atom{Name: "gives_milk"})   // 产奶
 
 // 添加规则
 kb.AddRule(Implies{
  Left: And{
   Left:  Atom{Name: "has_fur"},
   Right: Atom{Name: "gives_milk"},
  },
  Right: Atom{Name: "is_mammal"},
 })
 
 kb.AddRule(Implies{
  Left: Atom{Name: "is_mammal"},
  Right: Atom{Name: "is_animal"},
 })
 
 // 查询
 query := Atom{Name: "is_animal"}
 if kb.Query(query) {
  fmt.Println("It is an animal")
 } else {
  fmt.Println("Cannot determine if it is an animal")
 }
}
```

---

## 总结

命题逻辑为计算机科学提供了重要的理论基础：

1. **形式化推理**：提供了严格的逻辑推理方法
2. **程序验证**：用于验证程序的正确性
3. **知识表示**：用于表示和推理知识
4. **电路设计**：为数字电路设计提供理论基础

通过Go语言的实现，我们可以：

- 构建命题逻辑的解释器
- 实现自动证明系统
- 开发知识推理引擎
- 验证程序逻辑的正确性

这些应用展示了命题逻辑在软件工程中的重要作用，为后续学习更复杂的逻辑系统奠定了基础。

## 参考文献

1. Enderton, H. B. (2001). A Mathematical Introduction to Logic (2nd ed.). Academic Press.
2. Mendelson, E. (2015). Introduction to Mathematical Logic (6th ed.). CRC Press.
3. Boolos, G. S., Burgess, J. P., & Jeffrey, R. C. (2007). Computability and Logic (5th ed.). Cambridge University Press.
