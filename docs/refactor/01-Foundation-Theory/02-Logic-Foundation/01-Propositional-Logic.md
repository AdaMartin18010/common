# 01-命题逻辑 (Propositional Logic)

## 概述

命题逻辑是形式逻辑的基础分支，研究命题之间的逻辑关系。它提供了形式化推理的数学工具，是计算机科学中逻辑编程、形式化验证、人工智能等领域的重要基础。

## 1. 基本概念

### 1.1 命题

**定义 1.1** (命题)
命题是一个具有确定真值的陈述句，要么为真（True），要么为假（False）。

**形式化定义**：

```latex
P \in \{T, F\}
```

**示例**：

- "2 + 2 = 4" 是一个真命题
- "地球是平的" 是一个假命题
- "x + 1 = 5" 不是命题（因为x未定义）

### 1.2 命题变量

**定义 1.2** (命题变量)
命题变量是表示命题的符号，通常用大写字母表示，如 $P, Q, R$ 等。

**性质**：

- 每个命题变量只能取真值 $T$ 或假值 $F$
- 命题变量可以组合形成复合命题

### 1.3 逻辑连接词

#### 1.3.1 否定 (Negation)

**定义 1.3** (否定)
命题 $P$ 的否定记作 $\neg P$，读作"非P"。

**真值表**：

| $P$ | $\neg P$ |
|-----|----------|
| T   | F        |
| F   | T        |

**性质**：

- $\neg(\neg P) \equiv P$ (双重否定律)

#### 1.3.2 合取 (Conjunction)

**定义 1.4** (合取)
命题 $P$ 和 $Q$ 的合取记作 $P \land Q$，读作"P且Q"。

**真值表**：

| $P$ | $Q$ | $P \land Q$ |
|-----|-----|-------------|
| T   | T   | T           |
| T   | F   | F           |
| F   | T   | F           |
| F   | F   | F           |

**性质**：

- $P \land Q \equiv Q \land P$ (交换律)
- $(P \land Q) \land R \equiv P \land (Q \land R)$ (结合律)
- $P \land P \equiv P$ (幂等律)

#### 1.3.3 析取 (Disjunction)

**定义 1.5** (析取)
命题 $P$ 和 $Q$ 的析取记作 $P \lor Q$，读作"P或Q"。

**真值表**：

| $P$ | $Q$ | $P \lor Q$ |
|-----|-----|------------|
| T   | T   | T          |
| T   | F   | T          |
| F   | T   | T          |
| F   | F   | F          |

**性质**：

- $P \lor Q \equiv Q \lor P$ (交换律)
- $(P \lor Q) \lor R \equiv P \lor (Q \lor R)$ (结合律)
- $P \lor P \equiv P$ (幂等律)

#### 1.3.4 蕴含 (Implication)

**定义 1.6** (蕴含)
命题 $P$ 蕴含 $Q$ 记作 $P \rightarrow Q$，读作"如果P，那么Q"。

**真值表**：

| $P$ | $Q$ | $P \rightarrow Q$ |
|-----|-----|-------------------|
| T   | T   | T                 |
| T   | F   | F                 |
| F   | T   | T                 |
| F   | F   | T                 |

**性质**：

- $P \rightarrow Q \equiv \neg P \lor Q$
- $\neg(P \rightarrow Q) \equiv P \land \neg Q$

#### 1.3.5 等价 (Equivalence)

**定义 1.7** (等价)
命题 $P$ 等价于 $Q$ 记作 $P \leftrightarrow Q$，读作"P当且仅当Q"。

**真值表**：

| $P$ | $Q$ | $P \leftrightarrow Q$ |
|-----|-----|----------------------|
| T   | T   | T                    |
| T   | F   | F                    |
| F   | T   | F                    |
| F   | F   | T                    |

**性质**：

- $P \leftrightarrow Q \equiv (P \rightarrow Q) \land (Q \rightarrow P)$

## 2. 逻辑等价和重言式

### 2.1 逻辑等价

**定义 2.1** (逻辑等价)
两个命题公式 $A$ 和 $B$ 是逻辑等价的，记作 $A \equiv B$，当且仅当它们在所有真值赋值下具有相同的真值。

**定理 2.1** (德摩根律)
对于任意命题 $P$ 和 $Q$：

1. $\neg(P \land Q) \equiv \neg P \lor \neg Q$
2. $\neg(P \lor Q) \equiv \neg P \land \neg Q$

**证明**：
通过真值表验证：

| $P$ | $Q$ | $P \land Q$ | $\neg(P \land Q)$ | $\neg P$ | $\neg Q$ | $\neg P \lor \neg Q$ |
|-----|-----|-------------|-------------------|----------|----------|---------------------|
| T   | T   | T           | F                 | F        | F        | F                   |
| T   | F   | F           | T                 | F        | T        | T                   |
| F   | T   | F           | T                 | T        | F        | T                   |
| F   | F   | F           | T                 | T        | T        | T                   |

### 2.2 重言式和矛盾式

**定义 2.2** (重言式)
一个命题公式是重言式，当且仅当它在所有真值赋值下都为真。

**定义 2.3** (矛盾式)
一个命题公式是矛盾式，当且仅当它在所有真值赋值下都为假。

**定义 2.4** (可满足式)
一个命题公式是可满足式，当且仅当它在至少一个真值赋值下为真。

**定理 2.2** (重言式性质)

1. 如果 $A$ 是重言式，则 $\neg A$ 是矛盾式
2. 如果 $A$ 是矛盾式，则 $\neg A$ 是重言式
3. 如果 $A$ 和 $B$ 都是重言式，则 $A \land B$ 是重言式

## 3. 推理规则

### 3.1 基本推理规则

#### 3.1.1 假言推理 (Modus Ponens)

**规则**：
$$\frac{P \rightarrow Q \quad P}{Q}$$

**含义**：如果 $P \rightarrow Q$ 为真且 $P$ 为真，则 $Q$ 为真。

#### 3.1.2 假言三段论 (Hypothetical Syllogism)

**规则**：
$$\frac{P \rightarrow Q \quad Q \rightarrow R}{P \rightarrow R}$$

**含义**：如果 $P \rightarrow Q$ 为真且 $Q \rightarrow R$ 为真，则 $P \rightarrow R$ 为真。

#### 3.1.3 析取三段论 (Disjunctive Syllogism)

**规则**：
$$\frac{P \lor Q \quad \neg P}{Q}$$

**含义**：如果 $P \lor Q$ 为真且 $\neg P$ 为真，则 $Q$ 为真。

### 3.2 证明方法

#### 3.2.1 直接证明

**方法**：从前提直接推导出结论。

**示例**：
证明：如果 $P \rightarrow Q$ 且 $Q \rightarrow R$，则 $P \rightarrow R$

**证明**：

1. $P \rightarrow Q$ (前提)
2. $Q \rightarrow R$ (前提)
3. $P \rightarrow R$ (假言三段论，从1和2)

#### 3.2.2 反证法

**方法**：假设结论为假，推导出矛盾。

**示例**：
证明：$\neg(P \land \neg P)$

**证明**：

1. 假设 $P \land \neg P$ 为真
2. 从1可得 $P$ 为真
3. 从1可得 $\neg P$ 为真
4. 2和3矛盾
5. 因此假设错误，$\neg(P \land \neg P)$ 为真

## 4. Go语言实现

### 4.1 命题逻辑基础结构

```go
package propositional

import (
    "fmt"
    "strings"
)

// TruthValue 真值类型
type TruthValue bool

const (
    False TruthValue = false
    True  TruthValue = true
)

// String 真值的字符串表示
func (tv TruthValue) String() string {
    if tv {
        return "T"
    }
    return "F"
}

// Proposition 命题接口
type Proposition interface {
    Evaluate(assignment map[string]TruthValue) TruthValue
    String() string
    GetVariables() map[string]bool
}

// Variable 命题变量
type Variable struct {
    Name string
}

// NewVariable 创建命题变量
func NewVariable(name string) *Variable {
    return &Variable{Name: name}
}

// Evaluate 计算变量在给定赋值下的真值
func (v *Variable) Evaluate(assignment map[string]TruthValue) TruthValue {
    if value, exists := assignment[v.Name]; exists {
        return value
    }
    return False // 默认值
}

// String 变量的字符串表示
func (v *Variable) String() string {
    return v.Name
}

// GetVariables 获取变量集合
func (v *Variable) GetVariables() map[string]bool {
    return map[string]bool{v.Name: true}
}

// Negation 否定
type Negation struct {
    Operand Proposition
}

// NewNegation 创建否定
func NewNegation(operand Proposition) *Negation {
    return &Negation{Operand: operand}
}

// Evaluate 计算否定的真值
func (n *Negation) Evaluate(assignment map[string]TruthValue) TruthValue {
    return !n.Operand.Evaluate(assignment)
}

// String 否定的字符串表示
func (n *Negation) String() string {
    return fmt.Sprintf("¬(%s)", n.Operand.String())
}

// GetVariables 获取变量集合
func (n *Negation) GetVariables() map[string]bool {
    return n.Operand.GetVariables()
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

// Evaluate 计算合取的真值
func (c *Conjunction) Evaluate(assignment map[string]TruthValue) TruthValue {
    return c.Left.Evaluate(assignment) && c.Right.Evaluate(assignment)
}

// String 合取的字符串表示
func (c *Conjunction) String() string {
    return fmt.Sprintf("(%s ∧ %s)", c.Left.String(), c.Right.String())
}

// GetVariables 获取变量集合
func (c *Conjunction) GetVariables() map[string]bool {
    vars := c.Left.GetVariables()
    for v := range c.Right.GetVariables() {
        vars[v] = true
    }
    return vars
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

// Evaluate 计算析取的真值
func (d *Disjunction) Evaluate(assignment map[string]TruthValue) TruthValue {
    return d.Left.Evaluate(assignment) || d.Right.Evaluate(assignment)
}

// String 析取的字符串表示
func (d *Disjunction) String() string {
    return fmt.Sprintf("(%s ∨ %s)", d.Left.String(), d.Right.String())
}

// GetVariables 获取变量集合
func (d *Disjunction) GetVariables() map[string]bool {
    vars := d.Left.GetVariables()
    for v := range d.Right.GetVariables() {
        vars[v] = true
    }
    return vars
}

// Implication 蕴含
type Implication struct {
    Antecedent Proposition
    Consequent Proposition
}

// NewImplication 创建蕴含
func NewImplication(antecedent, consequent Proposition) *Implication {
    return &Implication{
        Antecedent: antecedent,
        Consequent: consequent,
    }
}

// Evaluate 计算蕴含的真值
func (i *Implication) Evaluate(assignment map[string]TruthValue) TruthValue {
    antecedent := i.Antecedent.Evaluate(assignment)
    consequent := i.Consequent.Evaluate(assignment)
    
    // P → Q ≡ ¬P ∨ Q
    return !antecedent || consequent
}

// String 蕴含的字符串表示
func (i *Implication) String() string {
    return fmt.Sprintf("(%s → %s)", i.Antecedent.String(), i.Consequent.String())
}

// GetVariables 获取变量集合
func (i *Implication) GetVariables() map[string]bool {
    vars := i.Antecedent.GetVariables()
    for v := range i.Consequent.GetVariables() {
        vars[v] = true
    }
    return vars
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

// Evaluate 计算等价的真值
func (e *Equivalence) Evaluate(assignment map[string]TruthValue) TruthValue {
    left := e.Left.Evaluate(assignment)
    right := e.Right.Evaluate(assignment)
    
    // P ↔ Q ≡ (P → Q) ∧ (Q → P)
    return left == right
}

// String 等价的字符串表示
func (e *Equivalence) String() string {
    return fmt.Sprintf("(%s ↔ %s)", e.Left.String(), e.Right.String())
}

// GetVariables 获取变量集合
func (e *Equivalence) GetVariables() map[string]bool {
    vars := e.Left.GetVariables()
    for v := range e.Right.GetVariables() {
        vars[v] = true
    }
    return vars
}
```

### 4.2 真值表生成

```go
// TruthTable 真值表
type TruthTable struct {
    Formula    Proposition
    Variables  []string
    Assignments [][]TruthValue
    Results    []TruthValue
}

// NewTruthTable 创建真值表
func NewTruthTable(formula Proposition) *TruthTable {
    vars := formula.GetVariables()
    varNames := make([]string, 0, len(vars))
    for name := range vars {
        varNames = append(varNames, name)
    }
    
    return &TruthTable{
        Formula:   formula,
        Variables: varNames,
    }
}

// Generate 生成真值表
func (tt *TruthTable) Generate() {
    n := len(tt.Variables)
    numAssignments := 1 << n // 2^n
    
    tt.Assignments = make([][]TruthValue, numAssignments)
    tt.Results = make([]TruthValue, numAssignments)
    
    for i := 0; i < numAssignments; i++ {
        assignment := make(map[string]TruthValue)
        tt.Assignments[i] = make([]TruthValue, n)
        
        for j := 0; j < n; j++ {
            value := (i & (1 << j)) != 0
            assignment[tt.Variables[j]] = TruthValue(value)
            tt.Assignments[i][j] = TruthValue(value)
        }
        
        tt.Results[i] = tt.Formula.Evaluate(assignment)
    }
}

// Print 打印真值表
func (tt *TruthTable) Print() {
    if len(tt.Assignments) == 0 {
        tt.Generate()
    }
    
    // 打印表头
    for _, varName := range tt.Variables {
        fmt.Printf("%-8s", varName)
    }
    fmt.Printf("%-20s\n", tt.Formula.String())
    
    // 打印分隔线
    for i := 0; i < len(tt.Variables); i++ {
        fmt.Print("--------")
    }
    fmt.Println("--------------------")
    
    // 打印真值表
    for i, assignment := range tt.Assignments {
        for _, value := range assignment {
            fmt.Printf("%-8s", value.String())
        }
        fmt.Printf("%-20s\n", tt.Results[i].String())
    }
}

// IsTautology 判断是否为重言式
func (tt *TruthTable) IsTautology() bool {
    if len(tt.Results) == 0 {
        tt.Generate()
    }
    
    for _, result := range tt.Results {
        if !result {
            return false
        }
    }
    return true
}

// IsContradiction 判断是否为矛盾式
func (tt *TruthTable) IsContradiction() bool {
    if len(tt.Results) == 0 {
        tt.Generate()
    }
    
    for _, result := range tt.Results {
        if result {
            return false
        }
    }
    return true
}

// IsSatisfiable 判断是否为可满足式
func (tt *TruthTable) IsSatisfiable() bool {
    if len(tt.Results) == 0 {
        tt.Generate()
    }
    
    for _, result := range tt.Results {
        if result {
            return true
        }
    }
    return false
}
```

### 4.3 逻辑推理引擎

```go
// InferenceEngine 推理引擎
type InferenceEngine struct {
    Premises []Proposition
    Rules    []InferenceRule
}

// InferenceRule 推理规则
type InferenceRule struct {
    Name     string
    Premises []Proposition
    Conclusion Proposition
}

// NewInferenceEngine 创建推理引擎
func NewInferenceEngine() *InferenceEngine {
    return &InferenceEngine{
        Premises: make([]Proposition, 0),
        Rules:    make([]InferenceRule, 0),
    }
}

// AddPremise 添加前提
func (ie *InferenceEngine) AddPremise(premise Proposition) {
    ie.Premises = append(ie.Premises, premise)
}

// AddRule 添加推理规则
func (ie *InferenceEngine) AddRule(rule InferenceRule) {
    ie.Rules = append(ie.Rules, rule)
}

// ModusPonens 假言推理规则
func ModusPonens(implication, antecedent Proposition) Proposition {
    if imp, ok := implication.(*Implication); ok {
        // 检查前提是否匹配
        if imp.Antecedent.String() == antecedent.String() {
            return imp.Consequent
        }
    }
    return nil
}

// HypotheticalSyllogism 假言三段论规则
func HypotheticalSyllogism(imp1, imp2 Proposition) Proposition {
    if i1, ok1 := imp1.(*Implication); ok1 {
        if i2, ok2 := imp2.(*Implication); ok2 {
            // 检查中间项是否匹配
            if i1.Consequent.String() == i2.Antecedent.String() {
                return NewImplication(i1.Antecedent, i2.Consequent)
            }
        }
    }
    return nil
}

// DisjunctiveSyllogism 析取三段论规则
func DisjunctiveSyllogism(disjunction, negation Proposition) Proposition {
    if disj, ok := disjunction.(*Disjunction); ok {
        if neg, ok := negation.(*Negation); ok {
            // 检查否定项是否匹配
            if disj.Left.String() == neg.Operand.String() {
                return disj.Right
            }
            if disj.Right.String() == neg.Operand.String() {
                return disj.Left
            }
        }
    }
    return nil
}

// Prove 证明结论
func (ie *InferenceEngine) Prove(conclusion Proposition) bool {
    // 简化的证明算法
    // 在实际应用中，需要更复杂的证明策略
    
    // 检查结论是否已经在前提中
    for _, premise := range ie.Premises {
        if premise.String() == conclusion.String() {
            return true
        }
    }
    
    // 尝试应用推理规则
    for _, rule := range ie.Rules {
        // 检查规则的前提是否满足
        premisesSatisfied := true
        for _, rulePremise := range rule.Premises {
            found := false
            for _, premise := range ie.Premises {
                if premise.String() == rulePremise.String() {
                    found = true
                    break
                }
            }
            if !found {
                premisesSatisfied = false
                break
            }
        }
        
        if premisesSatisfied && rule.Conclusion.String() == conclusion.String() {
            return true
        }
    }
    
    return false
}
```

## 5. 应用实例

### 5.1 逻辑电路设计

```go
// LogicGate 逻辑门接口
type LogicGate interface {
    Evaluate(inputs []bool) bool
    GetInputCount() int
}

// ANDGate 与门
type ANDGate struct{}

// Evaluate 与门计算
func (ag *ANDGate) Evaluate(inputs []bool) bool {
    for _, input := range inputs {
        if !input {
            return false
        }
    }
    return true
}

// GetInputCount 获取输入数量
func (ag *ANDGate) GetInputCount() int {
    return 2 // 标准与门有2个输入
}

// ORGate 或门
type ORGate struct{}

// Evaluate 或门计算
func (og *ORGate) Evaluate(inputs []bool) bool {
    for _, input := range inputs {
        if input {
            return true
        }
    }
    return false
}

// GetInputCount 获取输入数量
func (og *ORGate) GetInputCount() int {
    return 2 // 标准或门有2个输入
}

// NOTGate 非门
type NOTGate struct{}

// Evaluate 非门计算
func (ng *NOTGate) Evaluate(inputs []bool) bool {
    if len(inputs) != 1 {
        panic("非门只能有一个输入")
    }
    return !inputs[0]
}

// GetInputCount 获取输入数量
func (ng *NOTGate) GetInputCount() int {
    return 1
}

// LogicCircuit 逻辑电路
type LogicCircuit struct {
    gates  []LogicGate
    inputs []string
    outputs []string
    connections map[string][]string
}

// NewLogicCircuit 创建逻辑电路
func NewLogicCircuit() *LogicCircuit {
    return &LogicCircuit{
        gates:       make([]LogicGate, 0),
        inputs:      make([]string, 0),
        outputs:     make([]string, 0),
        connections: make(map[string][]string),
    }
}

// AddGate 添加逻辑门
func (lc *LogicCircuit) AddGate(gate LogicGate, name string) {
    lc.gates = append(lc.gates, gate)
    lc.connections[name] = make([]string, 0)
}

// Connect 连接门
func (lc *LogicCircuit) Connect(from, to string) {
    if connections, exists := lc.connections[from]; exists {
        lc.connections[from] = append(connections, to)
    }
}

// Evaluate 计算电路输出
func (lc *LogicCircuit) Evaluate(inputs map[string]bool) map[string]bool {
    // 简化的电路评估
    // 在实际应用中，需要拓扑排序和更复杂的评估逻辑
    outputs := make(map[string]bool)
    
    // 复制输入
    for input := range inputs {
        outputs[input] = inputs[input]
    }
    
    return outputs
}
```

### 5.2 知识表示系统

```go
// KnowledgeBase 知识库
type KnowledgeBase struct {
    facts    []Proposition
    rules    []Implication
    engine   *InferenceEngine
}

// NewKnowledgeBase 创建知识库
func NewKnowledgeBase() *KnowledgeBase {
    return &KnowledgeBase{
        facts:  make([]Proposition, 0),
        rules:  make([]Implication, 0),
        engine: NewInferenceEngine(),
    }
}

// AddFact 添加事实
func (kb *KnowledgeBase) AddFact(fact Proposition) {
    kb.facts = append(kb.facts, fact)
    kb.engine.AddPremise(fact)
}

// AddRule 添加规则
func (kb *KnowledgeBase) AddRule(rule Implication) {
    kb.rules = append(kb.rules, rule)
    kb.engine.AddPremise(&rule)
}

// Query 查询
func (kb *KnowledgeBase) Query(query Proposition) bool {
    return kb.engine.Prove(query)
}

// ForwardChaining 前向推理
func (kb *KnowledgeBase) ForwardChaining() []Proposition {
    conclusions := make([]Proposition, 0)
    
    // 简化的前向推理算法
    for _, rule := range kb.rules {
        if kb.engine.Prove(rule.Antecedent) {
            conclusions = append(conclusions, rule.Consequent)
        }
    }
    
    return conclusions
}

// BackwardChaining 后向推理
func (kb *KnowledgeBase) BackwardChaining(goal Proposition) bool {
    // 简化的后向推理算法
    return kb.engine.Prove(goal)
}
```

## 6. 性能分析

### 6.1 算法复杂度

| 操作 | 时间复杂度 | 空间复杂度 |
|------|------------|------------|
| 真值表生成 | $O(2^n)$ | $O(2^n)$ |
| 命题求值 | $O(1)$ | $O(1)$ |
| 逻辑等价检查 | $O(2^n)$ | $O(2^n)$ |
| 推理证明 | $O(n!)$ | $O(n)$ |

### 6.2 优化策略

1. **符号计算**：使用符号计算避免枚举所有真值赋值
2. **SAT求解器**：使用高效的SAT求解器
3. **缓存机制**：缓存中间计算结果
4. **启发式搜索**：在推理中使用启发式搜索

## 7. 总结

命题逻辑为计算机科学提供了形式化推理的基础。通过Go语言的实现，我们可以：

1. **形式化表示**：用数据结构表示命题和逻辑公式
2. **真值计算**：计算命题在各种赋值下的真值
3. **推理证明**：实现基本的逻辑推理规则
4. **应用开发**：构建逻辑电路、知识库等应用

命题逻辑的理论基础和算法实现为构建智能系统和形式化验证工具提供了重要的数学支撑。

## 参考文献

1. Enderton, H. B. (2001). A Mathematical Introduction to Logic (2nd ed.). Academic Press.
2. Mendelson, E. (2015). Introduction to Mathematical Logic (6th ed.). CRC Press.
3. Boolos, G. S., Burgess, J. P., & Jeffrey, R. C. (2007). Computability and Logic (5th ed.). Cambridge University Press.
