# 01-命题逻辑 (Propositional Logic)

## 目录

- [01-命题逻辑 (Propositional Logic)](#01-命题逻辑-propositional-logic)
  - [目录](#目录)
  - [概述](#概述)
  - [基本概念](#基本概念)
    - [命题](#命题)
    - [逻辑连接词](#逻辑连接词)
    - [真值表](#真值表)
  - [形式化理论](#形式化理论)
    - [语法](#语法)
    - [语义](#语义)
    - [推理规则](#推理规则)
    - [公理系统](#公理系统)
  - [逻辑等价](#逻辑等价)
    - [基本等价律](#基本等价律)
    - [德摩根律](#德摩根律)
    - [分配律](#分配律)
  - [范式](#范式)
    - [合取范式](#合取范式)
    - [析取范式](#析取范式)
    - [主范式](#主范式)
  - [Go语言实现](#go语言实现)
    - [命题数据结构](#命题数据结构)
    - [真值表计算](#真值表计算)
    - [推理引擎](#推理引擎)
    - [范式转换](#范式转换)
  - [应用领域](#应用领域)
    - [数字电路设计](#数字电路设计)
    - [程序验证](#程序验证)
    - [知识表示](#知识表示)
    - [自动推理](#自动推理)
  - [相关链接](#相关链接)

## 概述

命题逻辑是研究命题之间逻辑关系的数学分支，是形式逻辑的基础。命题逻辑通过严格的语法和语义定义，建立了处理真值判断的形式化系统。

## 基本概念

### 命题

**定义 1 (命题)**: 命题是一个具有确定真值的陈述句。

**定义 2 (原子命题)**: 原子命题是不可再分解的基本命题，用命题变量表示。

**定义 3 (复合命题)**: 复合命题是由原子命题通过逻辑连接词构成的命题。

### 逻辑连接词

**定义 4 (否定)**: 命题 $p$ 的否定记为 $\neg p$，真值表为：
| $p$ | $\neg p$ |
|-----|----------|
| T   | F        |
| F   | T        |

**定义 5 (合取)**: 命题 $p$ 和 $q$ 的合取记为 $p \wedge q$，真值表为：
| $p$ | $q$ | $p \wedge q$ |
|-----|-----|--------------|
| T   | T   | T            |
| T   | F   | F            |
| F   | T   | F            |
| F   | F   | F            |

**定义 6 (析取)**: 命题 $p$ 和 $q$ 的析取记为 $p \vee q$，真值表为：
| $p$ | $q$ | $p \vee q$ |
|-----|-----|------------|
| T   | T   | T          |
| T   | F   | T          |
| F   | T   | T          |
| F   | F   | F          |

**定义 7 (蕴含)**: 命题 $p$ 蕴含 $q$ 记为 $p \rightarrow q$，真值表为：
| $p$ | $q$ | $p \rightarrow q$ |
|-----|-----|-------------------|
| T   | T   | T                 |
| T   | F   | F                 |
| F   | T   | T                 |
| F   | F   | T                 |

**定义 8 (等价)**: 命题 $p$ 等价于 $q$ 记为 $p \leftrightarrow q$，真值表为：
| $p$ | $q$ | $p \leftrightarrow q$ |
|-----|-----|----------------------|
| T   | T   | T                    |
| T   | F   | F                    |
| F   | T   | F                    |
| F   | F   | T                    |

### 真值表

**定义 9 (真值表)**: 真值表是列出命题所有可能真值组合的表格。

**定理 1 (真值表的完备性)**: 任何复合命题的真值都可以通过真值表确定。

## 形式化理论

### 语法

**定义 10 (命题公式)**: 命题公式按以下规则递归定义：
1. 命题变量是公式
2. 如果 $\phi$ 是公式，则 $\neg \phi$ 是公式
3. 如果 $\phi$ 和 $\psi$ 是公式，则 $(\phi \wedge \psi)$、$(\phi \vee \psi)$、$(\phi \rightarrow \psi)$、$(\phi \leftrightarrow \psi)$ 是公式

**定义 11 (子公式)**: 公式 $\phi$ 的子公式是 $\phi$ 的组成部分。

### 语义

**定义 12 (解释)**: 解释 $I$ 是从命题变量到真值的映射。

**定义 13 (满足关系)**: 解释 $I$ 满足公式 $\phi$，记为 $I \models \phi$，按以下规则定义：
1. $I \models p$ 当且仅当 $I(p) = \text{true}$
2. $I \models \neg \phi$ 当且仅当 $I \not\models \phi$
3. $I \models \phi \wedge \psi$ 当且仅当 $I \models \phi$ 且 $I \models \psi$
4. $I \models \phi \vee \psi$ 当且仅当 $I \models \phi$ 或 $I \models \psi$
5. $I \models \phi \rightarrow \psi$ 当且仅当 $I \not\models \phi$ 或 $I \models \psi$
6. $I \models \phi \leftrightarrow \psi$ 当且仅当 $I \models \phi$ 等价于 $I \models \psi$

**定义 14 (永真式)**: 公式 $\phi$ 是永真式，如果对所有解释 $I$，$I \models \phi$。

**定义 15 (永假式)**: 公式 $\phi$ 是永假式，如果对所有解释 $I$，$I \not\models \phi$。

**定义 16 (可满足式)**: 公式 $\phi$ 是可满足式，如果存在解释 $I$，$I \models \phi$。

### 推理规则

**定理 2 (假言推理)**: 如果 $\phi \rightarrow \psi$ 和 $\phi$ 都为真，则 $\psi$ 为真。

**定理 3 (拒取式)**: 如果 $\phi \rightarrow \psi$ 和 $\neg \psi$ 都为真，则 $\neg \phi$ 为真。

**定理 4 (析取三段论)**: 如果 $\phi \vee \psi$ 和 $\neg \phi$ 都为真，则 $\psi$ 为真。

**定理 5 (构造性二难)**: 如果 $\phi \rightarrow \psi$、$\chi \rightarrow \delta$ 和 $\phi \vee \chi$ 都为真，则 $\psi \vee \delta$ 为真。

### 公理系统

**公理 1**: $\phi \rightarrow (\psi \rightarrow \phi)$

**公理 2**: $(\phi \rightarrow (\psi \rightarrow \chi)) \rightarrow ((\phi \rightarrow \psi) \rightarrow (\phi \rightarrow \chi))$

**公理 3**: $(\neg \phi \rightarrow \neg \psi) \rightarrow (\psi \rightarrow \phi)$

**推理规则 (分离规则)**: 从 $\phi$ 和 $\phi \rightarrow \psi$ 可以推出 $\psi$。

## 逻辑等价

### 基本等价律

**定理 6 (双重否定)**: $\neg \neg \phi \equiv \phi$

**定理 7 (幂等律)**: 
- $\phi \wedge \phi \equiv \phi$
- $\phi \vee \phi \equiv \phi$

**定理 8 (交换律)**:
- $\phi \wedge \psi \equiv \psi \wedge \phi$
- $\phi \vee \psi \equiv \psi \vee \phi$

**定理 9 (结合律)**:
- $(\phi \wedge \psi) \wedge \chi \equiv \phi \wedge (\psi \wedge \chi)$
- $(\phi \vee \psi) \vee \chi \equiv \phi \vee (\psi \vee \chi)$

### 德摩根律

**定理 10 (德摩根律)**:
- $\neg (\phi \wedge \psi) \equiv \neg \phi \vee \neg \psi$
- $\neg (\phi \vee \psi) \equiv \neg \phi \wedge \neg \psi$

### 分配律

**定理 11 (分配律)**:
- $\phi \wedge (\psi \vee \chi) \equiv (\phi \wedge \psi) \vee (\phi \wedge \chi)$
- $\phi \vee (\psi \wedge \chi) \equiv (\phi \vee \psi) \wedge (\phi \vee \chi)$

## 范式

### 合取范式

**定义 17 (文字)**: 文字是命题变量或其否定。

**定义 18 (子句)**: 子句是文字的析取。

**定义 19 (合取范式)**: 合取范式是子句的合取。

**定理 12 (合取范式存在性)**: 任何命题公式都可以转换为等价的合取范式。

### 析取范式

**定义 20 (短语)**: 短语是文字的合取。

**定义 21 (析取范式)**: 析取范式是短语的析取。

**定理 13 (析取范式存在性)**: 任何命题公式都可以转换为等价的析取范式。

### 主范式

**定义 22 (最小项)**: 最小项是包含所有命题变量或其否定的短语。

**定义 23 (最大项)**: 最大项是包含所有命题变量或其否定的子句。

**定义 24 (主析取范式)**: 主析取范式是最小项的析取。

**定义 25 (主合取范式)**: 主合取范式是最大项的合取。

## Go语言实现

### 命题数据结构

```go
package propositional

import (
    "fmt"
    "strings"
)

// Proposition 命题类型
type Proposition interface {
    Evaluate(interpretation map[string]bool) bool
    GetVariables() []string
    String() string
}

// Variable 命题变量
type Variable struct {
    Name string
}

// NewVariable 创建命题变量
func NewVariable(name string) *Variable {
    return &Variable{Name: name}
}

// Evaluate 求值
func (v *Variable) Evaluate(interpretation map[string]bool) bool {
    return interpretation[v.Name]
}

// GetVariables 获取变量
func (v *Variable) GetVariables() []string {
    return []string{v.Name}
}

// String 字符串表示
func (v *Variable) String() string {
    return v.Name
}

// Negation 否定
type Negation struct {
    Operand Proposition
}

// NewNegation 创建否定
func NewNegation(operand Proposition) *Negation {
    return &Negation{Operand: operand}
}

// Evaluate 求值
func (n *Negation) Evaluate(interpretation map[string]bool) bool {
    return !n.Operand.Evaluate(interpretation)
}

// GetVariables 获取变量
func (n *Negation) GetVariables() []string {
    return n.Operand.GetVariables()
}

// String 字符串表示
func (n *Negation) String() string {
    return fmt.Sprintf("¬(%s)", n.Operand.String())
}

// Conjunction 合取
type Conjunction struct {
    Left  Proposition
    Right Proposition
}

// NewConjunction 创建合取
func NewConjunction(left, right Proposition) *Conjunction {
    return &Conjunction{Left: left, Right: right}
}

// Evaluate 求值
func (c *Conjunction) Evaluate(interpretation map[string]bool) bool {
    return c.Left.Evaluate(interpretation) && c.Right.Evaluate(interpretation)
}

// GetVariables 获取变量
func (c *Conjunction) GetVariables() []string {
    vars := make(map[string]bool)
    for _, v := range c.Left.GetVariables() {
        vars[v] = true
    }
    for _, v := range c.Right.GetVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// String 字符串表示
func (c *Conjunction) String() string {
    return fmt.Sprintf("(%s ∧ %s)", c.Left.String(), c.Right.String())
}

// Disjunction 析取
type Disjunction struct {
    Left  Proposition
    Right Proposition
}

// NewDisjunction 创建析取
func NewDisjunction(left, right Proposition) *Disjunction {
    return &Disjunction{Left: left, Right: right}
}

// Evaluate 求值
func (d *Disjunction) Evaluate(interpretation map[string]bool) bool {
    return d.Left.Evaluate(interpretation) || d.Right.Evaluate(interpretation)
}

// GetVariables 获取变量
func (d *Disjunction) GetVariables() []string {
    vars := make(map[string]bool)
    for _, v := range d.Left.GetVariables() {
        vars[v] = true
    }
    for _, v := range d.Right.GetVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// String 字符串表示
func (d *Disjunction) String() string {
    return fmt.Sprintf("(%s ∨ %s)", d.Left.String(), d.Right.String())
}

// Implication 蕴含
type Implication struct {
    Antecedent Proposition
    Consequent Proposition
}

// NewImplication 创建蕴含
func NewImplication(antecedent, consequent Proposition) *Implication {
    return &Implication{Antecedent: antecedent, Consequent: consequent}
}

// Evaluate 求值
func (i *Implication) Evaluate(interpretation map[string]bool) bool {
    return !i.Antecedent.Evaluate(interpretation) || i.Consequent.Evaluate(interpretation)
}

// GetVariables 获取变量
func (i *Implication) GetVariables() []string {
    vars := make(map[string]bool)
    for _, v := range i.Antecedent.GetVariables() {
        vars[v] = true
    }
    for _, v := range i.Consequent.GetVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// String 字符串表示
func (i *Implication) String() string {
    return fmt.Sprintf("(%s → %s)", i.Antecedent.String(), i.Consequent.String())
}

// Equivalence 等价
type Equivalence struct {
    Left  Proposition
    Right Proposition
}

// NewEquivalence 创建等价
func NewEquivalence(left, right Proposition) *Equivalence {
    return &Equivalence{Left: left, Right: right}
}

// Evaluate 求值
func (e *Equivalence) Evaluate(interpretation map[string]bool) bool {
    return e.Left.Evaluate(interpretation) == e.Right.Evaluate(interpretation)
}

// GetVariables 获取变量
func (e *Equivalence) GetVariables() []string {
    vars := make(map[string]bool)
    for _, v := range e.Left.GetVariables() {
        vars[v] = true
    }
    for _, v := range e.Right.GetVariables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

// String 字符串表示
func (e *Equivalence) String() string {
    return fmt.Sprintf("(%s ↔ %s)", e.Left.String(), e.Right.String())
}
```

### 真值表计算

```go
// TruthTable 真值表
type TruthTable struct {
    Variables []string
    Formulas  []Proposition
    Rows      [][]bool
}

// NewTruthTable 创建真值表
func NewTruthTable(variables []string, formulas ...Proposition) *TruthTable {
    return &TruthTable{
        Variables: variables,
        Formulas:  formulas,
        Rows:      make([][]bool, 0),
    }
}

// Generate 生成真值表
func (tt *TruthTable) Generate() {
    tt.Rows = make([][]bool, 0)
    n := len(tt.Variables)
    
    // 生成所有可能的真值组合
    for i := 0; i < (1 << n); i++ {
        interpretation := make(map[string]bool)
        row := make([]bool, 0)
        
        // 设置变量值
        for j := 0; j < n; j++ {
            value := (i & (1 << j)) != 0
            interpretation[tt.Variables[j]] = value
            row = append(row, value)
        }
        
        // 计算公式值
        for _, formula := range tt.Formulas {
            row = append(row, formula.Evaluate(interpretation))
        }
        
        tt.Rows = append(tt.Rows, row)
    }
}

// Print 打印真值表
func (tt *TruthTable) Print() {
    // 打印表头
    header := make([]string, 0)
    for _, v := range tt.Variables {
        header = append(header, v)
    }
    for i, formula := range tt.Formulas {
        header = append(header, fmt.Sprintf("F%d", i+1))
    }
    
    fmt.Println(strings.Join(header, " | "))
    fmt.Println(strings.Repeat("-", len(header)*4))
    
    // 打印行
    for _, row := range tt.Rows {
        values := make([]string, 0)
        for _, value := range row {
            if value {
                values = append(values, "T")
            } else {
                values = append(values, "F")
            }
        }
        fmt.Println(strings.Join(values, " | "))
    }
}

// IsTautology 检查是否为永真式
func (tt *TruthTable) IsTautology(formulaIndex int) bool {
    if len(tt.Rows) == 0 {
        tt.Generate()
    }
    
    for _, row := range tt.Rows {
        if !row[len(tt.Variables)+formulaIndex] {
            return false
        }
    }
    return true
}

// IsContradiction 检查是否为永假式
func (tt *TruthTable) IsContradiction(formulaIndex int) bool {
    if len(tt.Rows) == 0 {
        tt.Generate()
    }
    
    for _, row := range tt.Rows {
        if row[len(tt.Variables)+formulaIndex] {
            return false
        }
    }
    return true
}

// IsSatisfiable 检查是否为可满足式
func (tt *TruthTable) IsSatisfiable(formulaIndex int) bool {
    if len(tt.Rows) == 0 {
        tt.Generate()
    }
    
    for _, row := range tt.Rows {
        if row[len(tt.Variables)+formulaIndex] {
            return true
        }
    }
    return false
}
```

### 推理引擎

```go
// InferenceEngine 推理引擎
type InferenceEngine struct{}

// ModusPonens 假言推理
func (ie *InferenceEngine) ModusPonens(premise1, premise2 Proposition) (Proposition, error) {
    // 检查前提1是否为蕴含
    if imp, ok := premise1.(*Implication); ok {
        // 检查前提2是否与蕴含的前件相等
        if imp.Antecedent.String() == premise2.String() {
            return imp.Consequent, nil
        }
    }
    return nil, fmt.Errorf("invalid modus ponens")
}

// ModusTollens 拒取式
func (ie *InferenceEngine) ModusTollens(premise1, premise2 Proposition) (Proposition, error) {
    // 检查前提1是否为蕴含
    if imp, ok := premise1.(*Implication); ok {
        // 检查前提2是否为蕴含后件的否定
        if neg, ok := premise2.(*Negation); ok {
            if neg.Operand.String() == imp.Consequent.String() {
                return NewNegation(imp.Antecedent), nil
            }
        }
    }
    return nil, fmt.Errorf("invalid modus tollens")
}

// DisjunctiveSyllogism 析取三段论
func (ie *InferenceEngine) DisjunctiveSyllogism(premise1, premise2 Proposition) (Proposition, error) {
    // 检查前提1是否为析取
    if disj, ok := premise1.(*Disjunction); ok {
        // 检查前提2是否为析取左项的否定
        if neg, ok := premise2.(*Negation); ok {
            if neg.Operand.String() == disj.Left.String() {
                return disj.Right, nil
            }
        }
        // 检查前提2是否为析取右项的否定
        if neg, ok := premise2.(*Negation); ok {
            if neg.Operand.String() == disj.Right.String() {
                return disj.Left, nil
            }
        }
    }
    return nil, fmt.Errorf("invalid disjunctive syllogism")
}

// ConstructiveDilemma 构造性二难
func (ie *InferenceEngine) ConstructiveDilemma(premise1, premise2, premise3 Proposition) (Proposition, error) {
    // 检查前提1和2是否为蕴含
    if imp1, ok := premise1.(*Implication); ok {
        if imp2, ok := premise2.(*Implication); ok {
            // 检查前提3是否为蕴含前件的析取
            if disj, ok := premise3.(*Disjunction); ok {
                if (disj.Left.String() == imp1.Antecedent.String() && 
                    disj.Right.String() == imp2.Antecedent.String()) ||
                   (disj.Left.String() == imp2.Antecedent.String() && 
                    disj.Right.String() == imp1.Antecedent.String()) {
                    return NewDisjunction(imp1.Consequent, imp2.Consequent), nil
                }
            }
        }
    }
    return nil, fmt.Errorf("invalid constructive dilemma")
}
```

### 范式转换

```go
// NormalFormConverter 范式转换器
type NormalFormConverter struct{}

// ToCNF 转换为合取范式
func (nfc *NormalFormConverter) ToCNF(proposition Proposition) Proposition {
    // 1. 消除蕴含和等价
    step1 := nfc.eliminateImplications(proposition)
    
    // 2. 将否定内移
    step2 := nfc.moveNegationsInward(step1)
    
    // 3. 分配析取
    step3 := nfc.distributeDisjunctions(step2)
    
    return step3
}

// ToDNF 转换为析取范式
func (nfc *NormalFormConverter) ToDNF(proposition Proposition) Proposition {
    // 1. 消除蕴含和等价
    step1 := nfc.eliminateImplications(proposition)
    
    // 2. 将否定内移
    step2 := nfc.moveNegationsInward(step1)
    
    // 3. 分配合取
    step3 := nfc.distributeConjunctions(step2)
    
    return step3
}

// eliminateImplications 消除蕴含和等价
func (nfc *NormalFormConverter) eliminateImplications(proposition Proposition) Proposition {
    switch p := proposition.(type) {
    case *Implication:
        return NewDisjunction(
            NewNegation(p.Antecedent),
            p.Consequent,
        )
    case *Equivalence:
        return NewConjunction(
            NewImplication(p.Left, p.Right),
            NewImplication(p.Right, p.Left),
        )
    case *Negation:
        return NewNegation(nfc.eliminateImplications(p.Operand))
    case *Conjunction:
        return NewConjunction(
            nfc.eliminateImplications(p.Left),
            nfc.eliminateImplications(p.Right),
        )
    case *Disjunction:
        return NewDisjunction(
            nfc.eliminateImplications(p.Left),
            nfc.eliminateImplications(p.Right),
        )
    default:
        return proposition
    }
}

// moveNegationsInward 将否定内移
func (nfc *NormalFormConverter) moveNegationsInward(proposition Proposition) Proposition {
    switch p := proposition.(type) {
    case *Negation:
        switch operand := p.Operand.(type) {
        case *Negation:
            return operand.Operand // 双重否定
        case *Conjunction:
            return NewDisjunction(
                NewNegation(operand.Left),
                NewNegation(operand.Right),
            )
        case *Disjunction:
            return NewConjunction(
                NewNegation(operand.Left),
                NewNegation(operand.Right),
            )
        default:
            return proposition
        }
    case *Conjunction:
        return NewConjunction(
            nfc.moveNegationsInward(p.Left),
            nfc.moveNegationsInward(p.Right),
        )
    case *Disjunction:
        return NewDisjunction(
            nfc.moveNegationsInward(p.Left),
            nfc.moveNegationsInward(p.Right),
        )
    default:
        return proposition
    }
}

// distributeDisjunctions 分配析取
func (nfc *NormalFormConverter) distributeDisjunctions(proposition Proposition) Proposition {
    switch p := proposition.(type) {
    case *Disjunction:
        left := nfc.distributeDisjunctions(p.Left)
        right := nfc.distributeDisjunctions(p.Right)
        
        // 如果左项是合取，分配
        if conj, ok := left.(*Conjunction); ok {
            return NewConjunction(
                nfc.distributeDisjunctions(NewDisjunction(conj.Left, right)),
                nfc.distributeDisjunctions(NewDisjunction(conj.Right, right)),
            )
        }
        
        // 如果右项是合取，分配
        if conj, ok := right.(*Conjunction); ok {
            return NewConjunction(
                nfc.distributeDisjunctions(NewDisjunction(left, conj.Left)),
                nfc.distributeDisjunctions(NewDisjunction(left, conj.Right)),
            )
        }
        
        return NewDisjunction(left, right)
    case *Conjunction:
        return NewConjunction(
            nfc.distributeDisjunctions(p.Left),
            nfc.distributeDisjunctions(p.Right),
        )
    default:
        return proposition
    }
}

// distributeConjunctions 分配合取
func (nfc *NormalFormConverter) distributeConjunctions(proposition Proposition) Proposition {
    switch p := proposition.(type) {
    case *Conjunction:
        left := nfc.distributeConjunctions(p.Left)
        right := nfc.distributeConjunctions(p.Right)
        
        // 如果左项是析取，分配
        if disj, ok := left.(*Disjunction); ok {
            return NewDisjunction(
                nfc.distributeConjunctions(NewConjunction(disj.Left, right)),
                nfc.distributeConjunctions(NewConjunction(disj.Right, right)),
            )
        }
        
        // 如果右项是析取，分配
        if disj, ok := right.(*Disjunction); ok {
            return NewDisjunction(
                nfc.distributeConjunctions(NewConjunction(left, disj.Left)),
                nfc.distributeConjunctions(NewConjunction(left, disj.Right)),
            )
        }
        
        return NewConjunction(left, right)
    case *Disjunction:
        return NewDisjunction(
            nfc.distributeConjunctions(p.Left),
            nfc.distributeConjunctions(p.Right),
        )
    default:
        return proposition
    }
}
```

## 应用领域

### 数字电路设计

命题逻辑在数字电路设计中的应用：
- 逻辑门设计
- 电路优化
- 故障检测
- 时序分析

### 程序验证

命题逻辑在程序验证中的应用：
- 程序正确性证明
- 模型检查
- 静态分析
- 形式化验证

### 知识表示

命题逻辑在知识表示中的应用：
- 专家系统
- 知识库
- 推理系统
- 决策支持

### 自动推理

命题逻辑在自动推理中的应用：
- 定理证明
- 约束求解
- 规划系统
- 智能代理

## 相关链接

- [02-谓词逻辑 (Predicate Logic)](../02-Predicate-Logic/README.md)
- [03-模态逻辑 (Modal Logic)](../03-Modal-Logic/README.md)
- [04-时态逻辑 (Temporal Logic)](../04-Temporal-Logic/README.md)
- [08-软件工程形式化 (Software Engineering Formalization)](../../08-Software-Engineering-Formalization/README.md)
- [09-编程语言理论 (Programming Language Theory)](../../09-Programming-Language-Theory/README.md) 