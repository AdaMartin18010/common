# 01-命题逻辑 (Propositional Logic)

## 目录

- [01-命题逻辑 (Propositional Logic)](#01-命题逻辑-propositional-logic)
  - [目录](#目录)
  - [1. 基础定义](#1-基础定义)
    - [1.1 命题](#11-命题)
    - [1.2 逻辑连接词](#12-逻辑连接词)
    - [1.3 命题公式](#13-命题公式)
  - [2. 语义学](#2-语义学)
    - [2.1 真值表](#21-真值表)
    - [2.2 语义函数](#22-语义函数)
    - [2.3 逻辑等价](#23-逻辑等价)
  - [3. 证明系统](#3-证明系统)
    - [3.1 自然演绎](#31-自然演绎)
    - [3.2 公理系统](#32-公理系统)
    - [3.3 归结原理](#33-归结原理)
  - [4. 形式化定义](#4-形式化定义)
    - [4.1 语法定义](#41-语法定义)
    - [4.2 语义定义](#42-语义定义)
    - [4.3 完备性定理](#43-完备性定理)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 命题公式表示](#51-命题公式表示)
    - [5.2 真值表计算](#52-真值表计算)
    - [5.3 证明系统实现](#53-证明系统实现)
  - [6. 应用实例](#6-应用实例)
    - [6.1 电路设计验证](#61-电路设计验证)
    - [6.2 程序逻辑验证](#62-程序逻辑验证)
    - [6.3 知识表示](#63-知识表示)
  - [总结](#总结)

## 1. 基础定义

### 1.1 命题

**定义 1.1** (命题)
命题是一个具有确定真值的陈述句，其真值要么为真（true），要么为假（false）。

**定义 1.2** (原子命题)
原子命题是不可再分解的基本命题，通常用大写字母 $P, Q, R, \ldots$ 表示。

**定义 1.3** (复合命题)
复合命题是由原子命题通过逻辑连接词组合而成的命题。

### 1.2 逻辑连接词

**定义 1.4** (否定)
否定连接词 $\neg$ 表示"非"，对于命题 $P$，$\neg P$ 表示"非 $P$"。

**定义 1.5** (合取)
合取连接词 $\wedge$ 表示"且"，对于命题 $P, Q$，$P \wedge Q$ 表示"$P$ 且 $Q$"。

**定义 1.6** (析取)
析取连接词 $\vee$ 表示"或"，对于命题 $P, Q$，$P \vee Q$ 表示"$P$ 或 $Q$"。

**定义 1.7** (蕴含)
蕴含连接词 $\rightarrow$ 表示"如果...那么"，对于命题 $P, Q$，$P \rightarrow Q$ 表示"如果 $P$ 那么 $Q$"。

**定义 1.8** (等价)
等价连接词 $\leftrightarrow$ 表示"当且仅当"，对于命题 $P, Q$，$P \leftrightarrow Q$ 表示"$P$ 当且仅当 $Q$"。

### 1.3 命题公式

**定义 1.9** (命题公式)
命题公式的递归定义：

1. 原子命题是命题公式
2. 如果 $\phi$ 是命题公式，则 $\neg \phi$ 是命题公式
3. 如果 $\phi, \psi$ 是命题公式，则 $(\phi \wedge \psi), (\phi \vee \psi), (\phi \rightarrow \psi), (\phi \leftrightarrow \psi)$ 是命题公式
4. 只有通过上述规则构造的表达式才是命题公式

## 2. 语义学

### 2.1 真值表

**定义 2.1** (真值赋值)
真值赋值是从原子命题集合到 $\{true, false\}$ 的函数。

**定义 2.2** (真值表)
真值表是描述命题公式在所有可能真值赋值下真值的表格。

**定理 2.1** (真值表构造)
对于包含 $n$ 个不同原子命题的命题公式，其真值表有 $2^n$ 行。

**证明**：
每个原子命题有两种可能的真值，根据乘法原理，$n$ 个原子命题共有 $2^n$ 种不同的真值赋值组合。

### 2.2 语义函数

**定义 2.3** (语义函数)
语义函数 $\llbracket \cdot \rrbracket$ 将命题公式映射到真值，满足：

1. $\llbracket P \rrbracket = v(P)$ 对于原子命题 $P$
2. $\llbracket \neg \phi \rrbracket = \neg \llbracket \phi \rrbracket$
3. $\llbracket \phi \wedge \psi \rrbracket = \llbracket \phi \rrbracket \wedge \llbracket \psi \rrbracket$
4. $\llbracket \phi \vee \psi \rrbracket = \llbracket \phi \rrbracket \vee \llbracket \psi \rrbracket$
5. $\llbracket \phi \rightarrow \psi \rrbracket = \neg \llbracket \phi \rrbracket \vee \llbracket \psi \rrbracket$
6. $\llbracket \phi \leftrightarrow \psi \rrbracket = (\llbracket \phi \rrbracket \rightarrow \llbracket \psi \rrbracket) \wedge (\llbracket \psi \rrbracket \rightarrow \llbracket \phi \rrbracket)$

### 2.3 逻辑等价

**定义 2.4** (逻辑等价)
两个命题公式 $\phi, \psi$ 是逻辑等价的，记作 $\phi \equiv \psi$，当且仅当对于所有真值赋值，$\llbracket \phi \rrbracket = \llbracket \psi \rrbracket$。

**定理 2.2** (德摩根律)
对于任意命题公式 $\phi, \psi$：

$$\neg(\phi \wedge \psi) \equiv \neg\phi \vee \neg\psi$$
$$\neg(\phi \vee \psi) \equiv \neg\phi \wedge \neg\psi$$

**定理 2.3** (分配律)
对于任意命题公式 $\phi, \psi, \chi$：

$$\phi \wedge (\psi \vee \chi) \equiv (\phi \wedge \psi) \vee (\phi \wedge \chi)$$
$$\phi \vee (\psi \wedge \chi) \equiv (\phi \vee \psi) \wedge (\phi \vee \chi)$$

## 3. 证明系统

### 3.1 自然演绎

**定义 3.1** (自然演绎)
自然演绎是一种基于推理规则的证明系统，包括引入和消除规则。

**规则 3.1** (合取引入)
$$\frac{\phi \quad \psi}{\phi \wedge \psi}$$

**规则 3.2** (合取消除)
$$\frac{\phi \wedge \psi}{\phi} \quad \frac{\phi \wedge \psi}{\psi}$$

**规则 3.3** (析取引入)
$$\frac{\phi}{\phi \vee \psi} \quad \frac{\psi}{\phi \vee \psi}$$

**规则 3.4** (蕴含引入)
$$\frac{[\phi] \quad \psi}{\phi \rightarrow \psi}$$

### 3.2 公理系统

**定义 3.2** (公理系统)
公理系统由一组公理和推理规则组成。

**公理 3.1** (同一律)
$$\phi \rightarrow \phi$$

**公理 3.2** (排中律)
$$\phi \vee \neg\phi$$

**公理 3.3** (矛盾律)
$$\neg(\phi \wedge \neg\phi)$$

### 3.3 归结原理

**定义 3.3** (归结原理)
归结原理是一种基于反证的证明方法。

**定理 3.1** (归结原理)
如果 $\Gamma \cup \{\neg\phi\}$ 不可满足，则 $\Gamma \models \phi$。

## 4. 形式化定义

### 4.1 语法定义

**定义 4.1** (命题语言)
命题语言 $\mathcal{L}$ 由以下组成：

1. 原子命题集合 $\mathcal{P}$
2. 逻辑连接词 $\{\neg, \wedge, \vee, \rightarrow, \leftrightarrow\}$
3. 括号 $\{(, )\}$

**定义 4.2** (命题公式的BNF语法)
```
φ ::= P | ¬φ | (φ ∧ φ) | (φ ∨ φ) | (φ → φ) | (φ ↔ φ)
```

### 4.2 语义定义

**定义 4.3** (解释)
解释是从原子命题集合到 $\{0,1\}$ 的函数。

**定义 4.4** (满足关系)
满足关系 $\models$ 递归定义如下：

1. $v \models P$ 当且仅当 $v(P) = 1$
2. $v \models \neg\phi$ 当且仅当 $v \not\models \phi$
3. $v \models \phi \wedge \psi$ 当且仅当 $v \models \phi$ 且 $v \models \psi$
4. $v \models \phi \vee \psi$ 当且仅当 $v \models \phi$ 或 $v \models \psi$
5. $v \models \phi \rightarrow \psi$ 当且仅当 $v \not\models \phi$ 或 $v \models \psi$

### 4.3 完备性定理

**定理 4.1** (完备性定理)
对于任意命题公式 $\phi$，如果 $\phi$ 是永真式，则 $\phi$ 是可证明的。

**定理 4.2** (可靠性定理)
对于任意命题公式 $\phi$，如果 $\phi$ 是可证明的，则 $\phi$ 是永真式。

## 5. Go语言实现

### 5.1 命题公式表示

```go
// 命题公式接口
type Formula interface {
    Evaluate(assignment map[string]bool) bool
    Variables() []string
    String() string
}

// 原子命题
type Atom struct {
    Name string
}

func (a *Atom) Evaluate(assignment map[string]bool) bool {
    return assignment[a.Name]
}

func (a *Atom) Variables() []string {
    return []string{a.Name}
}

func (a *Atom) String() string {
    return a.Name
}

// 否定公式
type Not struct {
    Operand Formula
}

func (n *Not) Evaluate(assignment map[string]bool) bool {
    return !n.Operand.Evaluate(assignment)
}

func (n *Not) Variables() []string {
    return n.Operand.Variables()
}

func (n *Not) String() string {
    return "¬" + n.Operand.String()
}

// 合取公式
type And struct {
    Left, Right Formula
}

func (a *And) Evaluate(assignment map[string]bool) bool {
    return a.Left.Evaluate(assignment) && a.Right.Evaluate(assignment)
}

func (a *And) Variables() []string {
    vars := make(map[string]bool)
    for _, v := range a.Left.Variables() {
        vars[v] = true
    }
    for _, v := range a.Right.Variables() {
        vars[v] = true
    }
    
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    return result
}

func (a *And) String() string {
    return "(" + a.Left.String() + " ∧ " + a.Right.String() + ")"
}
```

### 5.2 真值表计算

```go
// 真值表生成器
type TruthTable struct {
    formula Formula
}

func NewTruthTable(formula Formula) *TruthTable {
    return &TruthTable{formula: formula}
}

// 生成所有可能的真值赋值
func (tt *TruthTable) GenerateAssignments() []map[string]bool {
    variables := tt.formula.Variables()
    n := len(variables)
    assignments := make([]map[string]bool, 0, 1<<n)
    
    for i := 0; i < 1<<n; i++ {
        assignment := make(map[string]bool)
        for j, var := range variables {
            assignment[var] = (i>>j)&1 == 1
        }
        assignments = append(assignments, assignment)
    }
    
    return assignments
}

// 计算真值表
func (tt *TruthTable) Compute() []TruthTableRow {
    assignments := tt.GenerateAssignments()
    variables := tt.formula.Variables()
    rows := make([]TruthTableRow, 0, len(assignments))
    
    for _, assignment := range assignments {
        row := TruthTableRow{
            Assignment: assignment,
            Result:     tt.formula.Evaluate(assignment),
        }
        rows = append(rows, row)
    }
    
    return rows
}

type TruthTableRow struct {
    Assignment map[string]bool
    Result     bool
}

// 检查是否为永真式
func (tt *TruthTable) IsTautology() bool {
    rows := tt.Compute()
    for _, row := range rows {
        if !row.Result {
            return false
        }
    }
    return true
}

// 检查是否为矛盾式
func (tt *TruthTable) IsContradiction() bool {
    rows := tt.Compute()
    for _, row := range rows {
        if row.Result {
            return false
        }
    }
    return true
}
```

### 5.3 证明系统实现

```go
// 证明系统
type ProofSystem struct {
    axioms    []Formula
    rules     []InferenceRule
}

type InferenceRule struct {
    Premises    []Formula
    Conclusion  Formula
    Name        string
}

// 自然演绎系统
func NewNaturalDeduction() *ProofSystem {
    ps := &ProofSystem{}
    
    // 添加基本推理规则
    ps.rules = append(ps.rules, InferenceRule{
        Name: "合取引入",
        Premises: []Formula{},
        Conclusion: nil, // 需要根据具体公式确定
    })
    
    return ps
}

// 证明步骤
type ProofStep struct {
    Formula    Formula
    Justification string
    Dependencies []int
}

// 证明
type Proof struct {
    Steps []ProofStep
}

func (p *Proof) AddStep(formula Formula, justification string, deps ...int) {
    step := ProofStep{
        Formula:       formula,
        Justification: justification,
        Dependencies:  deps,
    }
    p.Steps = append(p.Steps, step)
}

// 验证证明
func (p *Proof) Validate() bool {
    for i, step := range p.Steps {
        if !p.validateStep(i, step) {
            return false
        }
    }
    return true
}

func (p *Proof) validateStep(index int, step ProofStep) bool {
    // 实现证明步骤验证逻辑
    return true
}
```

## 6. 应用实例

### 6.1 电路设计验证

```go
// 数字电路验证器
type CircuitVerifier struct {
    specification Formula
    implementation Formula
}

func (cv *CircuitVerifier) Verify() bool {
    // 验证实现是否满足规范
    implication := &Implies{
        Left:  cv.implementation,
        Right: cv.specification,
    }
    
    tt := NewTruthTable(implication)
    return tt.IsTautology()
}

// 电路组件
type CircuitComponent interface {
    Formula() Formula
}

// 与门
type ANDGate struct {
    Input1, Input2 string
    Output         string
}

func (ag *ANDGate) Formula() Formula {
    return &And{
        Left:  &Atom{Name: ag.Input1},
        Right: &Atom{Name: ag.Input2},
    }
}

// 或门
type ORGate struct {
    Input1, Input2 string
    Output         string
}

func (og *ORGate) Formula() Formula {
    return &Or{
        Left:  &Atom{Name: og.Input1},
        Right: &Atom{Name: og.Input2},
    }
}
```

### 6.2 程序逻辑验证

```go
// 程序逻辑验证器
type ProgramLogicVerifier struct {
    precondition  Formula
    postcondition Formula
    program       Program
}

type Program interface {
    Execute(state map[string]interface{}) map[string]interface{}
}

// 验证霍尔逻辑
func (plv *ProgramLogicVerifier) VerifyHoareLogic() bool {
    // 实现霍尔逻辑验证
    return true
}

// 程序状态
type ProgramState struct {
    Variables map[string]interface{}
}

// 条件语句
type IfStatement struct {
    Condition Formula
    ThenBranch Program
    ElseBranch Program
}

func (is *IfStatement) Execute(state map[string]interface{}) map[string]interface{} {
    // 将状态转换为真值赋值
    assignment := make(map[string]bool)
    for k, v := range state {
        if b, ok := v.(bool); ok {
            assignment[k] = b
        }
    }
    
    if is.Condition.Evaluate(assignment) {
        return is.ThenBranch.Execute(state)
    } else {
        return is.ElseBranch.Execute(state)
    }
}
```

### 6.3 知识表示

```go
// 知识库
type KnowledgeBase struct {
    formulas []Formula
}

func (kb *KnowledgeBase) AddFormula(formula Formula) {
    kb.formulas = append(kb.formulas, formula)
}

// 查询
func (kb *KnowledgeBase) Query(query Formula) bool {
    // 检查查询是否可以从知识库推导出
    return kb.entails(query)
}

func (kb *KnowledgeBase) entails(query Formula) bool {
    // 实现逻辑推理
    return true
}

// 专家系统
type ExpertSystem struct {
    knowledgeBase KnowledgeBase
    rules         []Rule
}

type Rule struct {
    Antecedent Formula
    Consequent Formula
}

func (es *ExpertSystem) Infer(facts []Formula) []Formula {
    // 实现前向推理
    return nil
}
```

## 总结

命题逻辑为计算机科学提供了形式化推理的基础，通过严格的数学定义和证明系统，我们可以：

1. **形式化建模**: 将复杂问题抽象为逻辑公式
2. **自动推理**: 使用算法进行逻辑推理和证明
3. **系统验证**: 验证软件和硬件系统的正确性
4. **知识表示**: 构建智能系统的知识库

在实际应用中，命题逻辑被广泛应用于：

- 数字电路设计验证
- 程序正确性证明
- 专家系统推理
- 人工智能知识表示
- 形式化方法

通过Go语言的实现，我们可以将这些理论概念转化为实用的工程工具，为软件工程提供可靠的形式化基础。 