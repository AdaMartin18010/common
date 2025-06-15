# 01-命题逻辑 (Propositional Logic)

## 目录

- [01-命题逻辑 (Propositional Logic)](#01-命题逻辑-propositional-logic)
  - [目录](#目录)
  - [1. 形式化定义](#1-形式化定义)
    - [1.1 命题的基本定义](#11-命题的基本定义)
    - [1.2 逻辑连接词](#12-逻辑连接词)
  - [2. 基本概念](#2-基本概念)
    - [2.1 真值赋值](#21-真值赋值)
    - [2.2 重言式与矛盾式](#22-重言式与矛盾式)
  - [3. 逻辑运算](#3-逻辑运算)
    - [3.1 基本运算律](#31-基本运算律)
    - [3.2 蕴含的等价形式](#32-蕴含的等价形式)
  - [4. 真值表](#4-真值表)
    - [4.1 基本真值表](#41-基本真值表)
    - [4.2 复合公式的真值表](#42-复合公式的真值表)
  - [5. 逻辑等价](#5-逻辑等价)
    - [5.1 等价关系](#51-等价关系)
    - [5.2 重要的等价关系](#52-重要的等价关系)
  - [6. 推理规则](#6-推理规则)
    - [6.1 基本推理规则](#61-基本推理规则)
    - [6.2 证明方法](#62-证明方法)
  - [7. Go语言实现](#7-go语言实现)
    - [7.1 基本类型定义](#71-基本类型定义)
    - [7.2 真值表生成](#72-真值表生成)
    - [7.3 逻辑等价检查](#73-逻辑等价检查)
    - [7.4 推理系统](#74-推理系统)
  - [8. 应用示例](#8-应用示例)
    - [8.1 基本使用示例](#81-基本使用示例)
    - [8.2 逻辑推理示例](#82-逻辑推理示例)
  - [9. 定理证明](#9-定理证明)
    - [9.1 德摩根律的证明](#91-德摩根律的证明)
    - [9.2 蕴含的等价形式证明](#92-蕴含的等价形式证明)
  - [10. 总结](#10-总结)

## 1. 形式化定义

### 1.1 命题的基本定义

**定义 1.1** (命题)
命题是一个陈述句，具有确定的真值（真或假）。命题通常用大写字母 $P, Q, R, \ldots$ 表示。

**定义 1.2** (原子命题)
原子命题是最基本的命题，不能再分解为更简单的命题。

**定义 1.3** (复合命题)
复合命题是由原子命题通过逻辑连接词构成的命题。

### 1.2 逻辑连接词

**定义 1.4** (否定)
命题 $P$ 的否定记作 $\neg P$，读作"非 $P$"。

- 如果 $P$ 为真，则 $\neg P$ 为假
- 如果 $P$ 为假，则 $\neg P$ 为真

**定义 1.5** (合取)
命题 $P$ 和 $Q$ 的合取记作 $P \land Q$，读作"$P$ 且 $Q$"。

- $P \land Q$ 为真当且仅当 $P$ 和 $Q$ 都为真

**定义 1.6** (析取)
命题 $P$ 和 $Q$ 的析取记作 $P \lor Q$，读作"$P$ 或 $Q$"。

- $P \lor Q$ 为真当且仅当 $P$ 或 $Q$ 至少有一个为真

**定义 1.7** (蕴含)
命题 $P$ 蕴含 $Q$ 记作 $P \rightarrow Q$，读作"如果 $P$ 则 $Q$"。

- $P \rightarrow Q$ 为假当且仅当 $P$ 为真且 $Q$ 为假

**定义 1.8** (等价)
命题 $P$ 等价于 $Q$ 记作 $P \leftrightarrow Q$，读作"$P$ 当且仅当 $Q$"。

- $P \leftrightarrow Q$ 为真当且仅当 $P$ 和 $Q$ 具有相同的真值

## 2. 基本概念

### 2.1 真值赋值

**定义 2.1** (真值赋值)
真值赋值是一个函数 $v: \mathcal{P} \rightarrow \{T, F\}$，其中 $\mathcal{P}$ 是原子命题的集合。

**定义 2.2** (真值扩展)
给定真值赋值 $v$，可以扩展到所有命题公式：

1. $v(\neg P) = T$ 当且仅当 $v(P) = F$
2. $v(P \land Q) = T$ 当且仅当 $v(P) = T$ 且 $v(Q) = T$
3. $v(P \lor Q) = T$ 当且仅当 $v(P) = T$ 或 $v(Q) = T$
4. $v(P \rightarrow Q) = T$ 当且仅当 $v(P) = F$ 或 $v(Q) = T$
5. $v(P \leftrightarrow Q) = T$ 当且仅当 $v(P) = v(Q)$

### 2.2 重言式与矛盾式

**定义 2.3** (重言式)
命题公式 $A$ 是重言式，当且仅当对于所有真值赋值 $v$，都有 $v(A) = T$。

**定义 2.4** (矛盾式)
命题公式 $A$ 是矛盾式，当且仅当对于所有真值赋值 $v$，都有 $v(A) = F$。

**定义 2.5** (可满足式)
命题公式 $A$ 是可满足式，当且仅当存在真值赋值 $v$，使得 $v(A) = T$。

## 3. 逻辑运算

### 3.1 基本运算律

**定理 3.1** (双重否定律)
$$\neg \neg P \leftrightarrow P$$

**定理 3.2** (德摩根律)
$$\neg (P \land Q) \leftrightarrow \neg P \lor \neg Q$$
$$\neg (P \lor Q) \leftrightarrow \neg P \land \neg Q$$

**定理 3.3** (分配律)
$$P \land (Q \lor R) \leftrightarrow (P \land Q) \lor (P \land R)$$
$$P \lor (Q \land R) \leftrightarrow (P \lor Q) \land (P \lor R)$$

**定理 3.4** (结合律)
$$(P \land Q) \land R \leftrightarrow P \land (Q \land R)$$
$$(P \lor Q) \lor R \leftrightarrow P \lor (Q \lor R)$$

**定理 3.5** (交换律)
$$P \land Q \leftrightarrow Q \land P$$
$$P \lor Q \leftrightarrow Q \lor P$$

### 3.2 蕴含的等价形式

**定理 3.6** (蕴含的等价形式)
$$P \rightarrow Q \leftrightarrow \neg P \lor Q$$
$$P \rightarrow Q \leftrightarrow \neg Q \rightarrow \neg P$$

## 4. 真值表

### 4.1 基本真值表

| $P$ | $Q$ | $\neg P$ | $P \land Q$ | $P \lor Q$ | $P \rightarrow Q$ | $P \leftrightarrow Q$ |
|-----|-----|----------|-------------|------------|-------------------|----------------------|
| T   | T   | F        | T           | T          | T                 | T                    |
| T   | F   | F        | F           | T          | F                 | F                    |
| F   | T   | T        | F           | T          | T                 | F                    |
| F   | F   | T        | F           | F          | T                 | T                    |

### 4.2 复合公式的真值表

**示例 4.1** (德摩根律的真值表验证)

| $P$ | $Q$ | $\neg P$ | $\neg Q$ | $P \land Q$ | $\neg(P \land Q)$ | $\neg P \lor \neg Q$ |
|-----|-----|----------|----------|-------------|-------------------|----------------------|
| T   | T   | F        | F        | T           | F                 | F                    |
| T   | F   | F        | T        | F           | T                 | T                    |
| F   | T   | T        | F        | F           | T                 | T                    |
| F   | F   | T        | T        | F           | T                 | T                    |

## 5. 逻辑等价

### 5.1 等价关系

**定义 5.1** (逻辑等价)
两个命题公式 $A$ 和 $B$ 逻辑等价，记作 $A \equiv B$，当且仅当对于所有真值赋值 $v$，都有 $v(A) = v(B)$。

**定理 5.1** (等价的性质)

1. **自反性**: $A \equiv A$
2. **对称性**: 如果 $A \equiv B$，则 $B \equiv A$
3. **传递性**: 如果 $A \equiv B$ 且 $B \equiv C$，则 $A \equiv C$

### 5.2 重要的等价关系

**定理 5.2** (幂等律)
$$P \land P \equiv P$$
$$P \lor P \equiv P$$

**定理 5.3** (吸收律)
$$P \land (P \lor Q) \equiv P$$
$$P \lor (P \land Q) \equiv P$$

**定理 5.4** (对偶律)
$$P \land F \equiv F$$
$$P \lor T \equiv T$$
$$P \land T \equiv P$$
$$P \lor F \equiv P$$

## 6. 推理规则

### 6.1 基本推理规则

**规则 6.1** (假言推理)
$$\frac{P \rightarrow Q \quad P}{Q}$$

**规则 6.2** (否定后件)
$$\frac{P \rightarrow Q \quad \neg Q}{\neg P}$$

**规则 6.3** (合取引入)
$$\frac{P \quad Q}{P \land Q}$$

**规则 6.4** (合取消除)
$$\frac{P \land Q}{P} \quad \frac{P \land Q}{Q}$$

**规则 6.5** (析取引入)
$$\frac{P}{P \lor Q} \quad \frac{Q}{P \lor Q}$$

**规则 6.6** (析取消除)
$$\frac{P \lor Q \quad P \rightarrow R \quad Q \rightarrow R}{R}$$

### 6.2 证明方法

**方法 6.1** (直接证明)
要证明 $P \rightarrow Q$，假设 $P$ 为真，然后证明 $Q$ 为真。

**方法 6.2** (反证法)
要证明 $P$，假设 $\neg P$ 为真，然后导出矛盾。

**方法 6.3** (分情况证明)
要证明 $P \lor Q \rightarrow R$，分别证明 $P \rightarrow R$ 和 $Q \rightarrow R$。

## 7. Go语言实现

### 7.1 基本类型定义

```go
// Proposition 表示命题
type Proposition interface {
    Evaluate(assignment map[string]bool) bool
    GetVariables() []string
    String() string
}

// AtomicProposition 表示原子命题
type AtomicProposition struct {
    Name string
}

func (ap AtomicProposition) Evaluate(assignment map[string]bool) bool {
    return assignment[ap.Name]
}

func (ap AtomicProposition) GetVariables() []string {
    return []string{ap.Name}
}

func (ap AtomicProposition) String() string {
    return ap.Name
}

// Negation 表示否定
type Negation struct {
    Proposition Proposition
}

func (n Negation) Evaluate(assignment map[string]bool) bool {
    return !n.Proposition.Evaluate(assignment)
}

func (n Negation) GetVariables() []string {
    return n.Proposition.GetVariables()
}

func (n Negation) String() string {
    return "¬" + n.Proposition.String()
}

// Conjunction 表示合取
type Conjunction struct {
    Left  Proposition
    Right Proposition
}

func (c Conjunction) Evaluate(assignment map[string]bool) bool {
    return c.Left.Evaluate(assignment) && c.Right.Evaluate(assignment)
}

func (c Conjunction) GetVariables() []string {
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

func (c Conjunction) String() string {
    return "(" + c.Left.String() + " ∧ " + c.Right.String() + ")"
}

// Disjunction 表示析取
type Disjunction struct {
    Left  Proposition
    Right Proposition
}

func (d Disjunction) Evaluate(assignment map[string]bool) bool {
    return d.Left.Evaluate(assignment) || d.Right.Evaluate(assignment)
}

func (d Disjunction) GetVariables() []string {
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

func (d Disjunction) String() string {
    return "(" + d.Left.String() + " ∨ " + d.Right.String() + ")"
}

// Implication 表示蕴含
type Implication struct {
    Antecedent Proposition
    Consequent Proposition
}

func (i Implication) Evaluate(assignment map[string]bool) bool {
    return !i.Antecedent.Evaluate(assignment) || i.Consequent.Evaluate(assignment)
}

func (i Implication) GetVariables() []string {
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

func (i Implication) String() string {
    return "(" + i.Antecedent.String() + " → " + i.Consequent.String() + ")"
}

// Equivalence 表示等价
type Equivalence struct {
    Left  Proposition
    Right Proposition
}

func (e Equivalence) Evaluate(assignment map[string]bool) bool {
    return e.Left.Evaluate(assignment) == e.Right.Evaluate(assignment)
}

func (e Equivalence) GetVariables() []string {
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

func (e Equivalence) String() string {
    return "(" + e.Left.String() + " ↔ " + e.Right.String() + ")"
}
```

### 7.2 真值表生成

```go
// TruthTable 表示真值表
type TruthTable struct {
    Variables []string
    Rows      []TruthTableRow
}

type TruthTableRow struct {
    Assignment map[string]bool
    Result     bool
}

// GenerateTruthTable 生成命题的真值表
func GenerateTruthTable(prop Proposition) TruthTable {
    variables := prop.GetVariables()
    numVars := len(variables)
    numRows := 1 << numVars
    
    var rows []TruthTableRow
    
    for i := 0; i < numRows; i++ {
        assignment := make(map[string]bool)
        for j, varName := range variables {
            assignment[varName] = (i & (1 << j)) != 0
        }
        
        result := prop.Evaluate(assignment)
        rows = append(rows, TruthTableRow{
            Assignment: assignment,
            Result:     result,
        })
    }
    
    return TruthTable{
        Variables: variables,
        Rows:      rows,
    }
}

// IsTautology 检查命题是否是重言式
func IsTautology(prop Proposition) bool {
    table := GenerateTruthTable(prop)
    for _, row := range table.Rows {
        if !row.Result {
            return false
        }
    }
    return true
}

// IsContradiction 检查命题是否是矛盾式
func IsContradiction(prop Proposition) bool {
    table := GenerateTruthTable(prop)
    for _, row := range table.Rows {
        if row.Result {
            return false
        }
    }
    return true
}

// IsSatisfiable 检查命题是否是可满足式
func IsSatisfiable(prop Proposition) bool {
    table := GenerateTruthTable(prop)
    for _, row := range table.Rows {
        if row.Result {
            return true
        }
    }
    return false
}
```

### 7.3 逻辑等价检查

```go
// AreEquivalent 检查两个命题是否逻辑等价
func AreEquivalent(prop1, prop2 Proposition) bool {
    // 获取所有变量
    vars1 := prop1.GetVariables()
    vars2 := prop2.GetVariables()
    
    // 合并变量
    allVars := make(map[string]bool)
    for _, v := range vars1 {
        allVars[v] = true
    }
    for _, v := range vars2 {
        allVars[v] = true
    }
    
    variables := make([]string, 0, len(allVars))
    for v := range allVars {
        variables = append(variables, v)
    }
    
    numVars := len(variables)
    numRows := 1 << numVars
    
    // 检查所有可能的真值赋值
    for i := 0; i < numRows; i++ {
        assignment := make(map[string]bool)
        for j, varName := range variables {
            assignment[varName] = (i & (1 << j)) != 0
        }
        
        if prop1.Evaluate(assignment) != prop2.Evaluate(assignment) {
            return false
        }
    }
    
    return true
}

// Simplify 简化命题（基本实现）
func Simplify(prop Proposition) Proposition {
    // 这里可以实现各种简化规则
    // 例如：双重否定、德摩根律等
    return prop
}
```

### 7.4 推理系统

```go
// Proof 表示证明
type Proof struct {
    Premises []Proposition
    Conclusion Proposition
    Steps    []ProofStep
}

type ProofStep struct {
    Proposition Proposition
    Rule       string
    Dependencies []int // 依赖的步骤索引
}

// ValidProof 检查证明是否有效
func ValidProof(proof Proof) bool {
    // 检查每个步骤是否可以从前提和前面的步骤推导出
    for i, step := range proof.Steps {
        if !canDerive(step, proof.Steps[:i], proof.Premises) {
            return false
        }
    }
    
    // 检查结论是否在最后一步
    if len(proof.Steps) == 0 {
        return false
    }
    
    lastStep := proof.Steps[len(proof.Steps)-1]
    return AreEquivalent(lastStep.Proposition, proof.Conclusion)
}

func canDerive(step ProofStep, previousSteps []ProofStep, premises []Proposition) bool {
    // 这里可以实现具体的推理规则检查
    // 例如：假言推理、合取引入等
    return true // 简化实现
}
```

## 8. 应用示例

### 8.1 基本使用示例

```go
func main() {
    // 创建原子命题
    p := AtomicProposition{Name: "P"}
    q := AtomicProposition{Name: "Q"}
    
    // 创建复合命题
    notP := Negation{Proposition: p}
    pAndQ := Conjunction{Left: p, Right: q}
    pOrQ := Disjunction{Left: p, Right: q}
    pImpliesQ := Implication{Antecedent: p, Consequent: q}
    
    // 生成真值表
    fmt.Println("Truth table for P ∧ Q:")
    table := GenerateTruthTable(pAndQ)
    printTruthTable(table)
    
    // 检查重言式
    tautology := Implication{
        Antecedent: p,
        Consequent: Disjunction{Left: p, Right: q},
    }
    fmt.Printf("P → (P ∨ Q) is a tautology: %v\n", IsTautology(tautology))
    
    // 检查逻辑等价
    deMorgan1 := Negation{Proposition: Conjunction{Left: p, Right: q}}
    deMorgan2 := Disjunction{Left: Negation{Proposition: p}, Right: Negation{Proposition: q}}
    fmt.Printf("¬(P ∧ Q) ≡ (¬P ∨ ¬Q): %v\n", AreEquivalent(deMorgan1, deMorgan2))
}

func printTruthTable(table TruthTable) {
    // 打印表头
    for _, varName := range table.Variables {
        fmt.Printf("%s\t", varName)
    }
    fmt.Println("Result")
    
    // 打印分隔线
    for range table.Variables {
        fmt.Print("---\t")
    }
    fmt.Println("------")
    
    // 打印数据行
    for _, row := range table.Rows {
        for _, varName := range table.Variables {
            if row.Assignment[varName] {
                fmt.Print("T\t")
            } else {
                fmt.Print("F\t")
            }
        }
        if row.Result {
            fmt.Println("T")
        } else {
            fmt.Println("F")
        }
    }
}
```

### 8.2 逻辑推理示例

```go
// 验证假言推理
func verifyModusPonens() {
    p := AtomicProposition{Name: "P"}
    q := AtomicProposition{Name: "Q"}
    
    premise1 := Implication{Antecedent: p, Consequent: q}
    premise2 := p
    conclusion := q
    
    // 检查 (P → Q) ∧ P → Q 是否是重言式
    formula := Implication{
        Antecedent: Conjunction{Left: premise1, Right: premise2},
        Consequent: conclusion,
    }
    
    fmt.Printf("Modus Ponens is valid: %v\n", IsTautology(formula))
}

// 验证反证法
func verifyProofByContradiction() {
    p := AtomicProposition{Name: "P"}
    
    // 要证明 P，假设 ¬P 并导出矛盾
    notP := Negation{Proposition: p}
    contradiction := Conjunction{Left: p, Right: notP}
    
    // 检查 (¬P → (P ∧ ¬P)) → P 是否是重言式
    formula := Implication{
        Antecedent: Implication{Antecedent: notP, Consequent: contradiction},
        Consequent: p,
    }
    
    fmt.Printf("Proof by contradiction is valid: %v\n", IsTautology(formula))
}
```

## 9. 定理证明

### 9.1 德摩根律的证明

**定理 9.1** (德摩根律)
$$\neg (P \land Q) \equiv \neg P \lor \neg Q$$

**证明**:
我们需要证明对于所有真值赋值 $v$，都有：
$$v(\neg (P \land Q)) = v(\neg P \lor \neg Q)$$

分情况讨论：

1. **情况1**: $v(P) = T$ 且 $v(Q) = T$
   - $v(P \land Q) = T$
   - $v(\neg (P \land Q)) = F$
   - $v(\neg P) = F$, $v(\neg Q) = F$
   - $v(\neg P \lor \neg Q) = F$

2. **情况2**: $v(P) = T$ 且 $v(Q) = F$
   - $v(P \land Q) = F$
   - $v(\neg (P \land Q)) = T$
   - $v(\neg P) = F$, $v(\neg Q) = T$
   - $v(\neg P \lor \neg Q) = T$

3. **情况3**: $v(P) = F$ 且 $v(Q) = T$
   - $v(P \land Q) = F$
   - $v(\neg (P \land Q)) = T$
   - $v(\neg P) = T$, $v(\neg Q) = F$
   - $v(\neg P \lor \neg Q) = T$

4. **情况4**: $v(P) = F$ 且 $v(Q) = F$
   - $v(P \land Q) = F$
   - $v(\neg (P \land Q)) = T$
   - $v(\neg P) = T$, $v(\neg Q) = T$
   - $v(\neg P \lor \neg Q) = T$

在所有情况下，$v(\neg (P \land Q)) = v(\neg P \lor \neg Q)$，因此德摩根律成立。

### 9.2 蕴含的等价形式证明

**定理 9.2** (蕴含的等价形式)
$$P \rightarrow Q \equiv \neg P \lor Q$$

**证明**:
我们需要证明对于所有真值赋值 $v$，都有：
$$v(P \rightarrow Q) = v(\neg P \lor Q)$$

分情况讨论：

1. **情况1**: $v(P) = T$ 且 $v(Q) = T$
   - $v(P \rightarrow Q) = T$
   - $v(\neg P) = F$
   - $v(\neg P \lor Q) = T$

2. **情况2**: $v(P) = T$ 且 $v(Q) = F$
   - $v(P \rightarrow Q) = F$
   - $v(\neg P) = F$
   - $v(\neg P \lor Q) = F$

3. **情况3**: $v(P) = F$ 且 $v(Q) = T$
   - $v(P \rightarrow Q) = T$
   - $v(\neg P) = T$
   - $v(\neg P \lor Q) = T$

4. **情况4**: $v(P) = F$ 且 $v(Q) = F$
   - $v(P \rightarrow Q) = T$
   - $v(\neg P) = T$
   - $v(\neg P \lor Q) = T$

在所有情况下，$v(P \rightarrow Q) = v(\neg P \lor Q)$，因此蕴含的等价形式成立。

## 10. 总结

命题逻辑为数学和计算机科学提供了重要的逻辑基础：

1. **形式化基础**: 提供了严格的逻辑定义和推理规则
2. **运算体系**: 建立了完整的逻辑运算体系
3. **推理系统**: 为数学证明和逻辑推理提供了工具
4. **应用价值**: 在计算机科学、人工智能、电路设计等领域有广泛应用
5. **理论基础**: 为更高级的逻辑系统（如一阶逻辑）奠定了基础

通过Go语言的实现，我们展示了命题逻辑概念在编程中的具体应用，为后续的逻辑理论和算法实现奠定了基础。

---

**参考文献**:

1. Enderton, H. B. (2001). A Mathematical Introduction to Logic. Academic Press.
2. Mendelson, E. (2015). Introduction to Mathematical Logic. CRC Press.
3. Boolos, G. S., Burgess, J. P., & Jeffrey, R. C. (2007). Computability and Logic. Cambridge University Press.
