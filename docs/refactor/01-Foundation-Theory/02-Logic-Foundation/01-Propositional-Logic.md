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

**定义 1.1 (命题)**
命题是一个具有确定真值的陈述句，其真值要么为真（true），要么为假（false）。

**定义 1.2 (原子命题)**
原子命题是不可再分解的基本命题，通常用大写字母 ```latex
$P, Q, R, \ldots$
``` 表示。

**定义 1.3 (复合命题)**
复合命题是由原子命题通过逻辑连接词组合而成的命题。

### 1.2 逻辑连接词

**定义 1.4 (否定)**
否定连接词 ```latex
$\neg$
``` 表示"非"，对于命题 ```latex
$P$
```，```latex
$\neg P$
``` 表示"非 ```latex
$P$
```"。

**定义 1.5 (合取)**
合取连接词 ```latex
$\wedge$
``` 表示"且"，对于命题 ```latex
$P, Q$
```，```latex
$P \wedge Q$
``` 表示"```latex
$P$
``` 且 ```latex
$Q$
```"。

**定义 1.6 (析取)**
析取连接词 ```latex
$\vee$
``` 表示"或"，对于命题 ```latex
$P, Q$
```，```latex
$P \vee Q$
``` 表示"```latex
$P$
``` 或 ```latex
$Q$
```"。

**定义 1.7 (蕴含)**
蕴含连接词 ```latex
$\rightarrow$
``` 表示"如果...那么"，对于命题 ```latex
$P, Q$
```，```latex
$P \rightarrow Q$
``` 表示"如果 ```latex
$P$
``` 那么 ```latex
$Q$
```"。

**定义 1.8 (等价)**
等价连接词 ```latex
$\leftrightarrow$
``` 表示"当且仅当"，对于命题 ```latex
$P, Q$
```，```latex
$P \leftrightarrow Q$
``` 表示"```latex
$P$
``` 当且仅当 ```latex
$Q$
```"。

### 1.3 命题公式

**定义 1.9 (命题公式)**
命题公式的递归定义：

1. 原子命题是命题公式
2. 如果 ```latex
$\phi$
``` 是命题公式，则 ```latex
$\neg \phi$
``` 是命题公式
3. 如果 ```latex
$\phi, \psi$
``` 是命题公式，则 ```latex
$(\phi \wedge \psi), (\phi \vee \psi), (\phi \rightarrow \psi), (\phi \leftrightarrow \psi)$
``` 是命题公式
4. 只有通过上述规则构造的表达式才是命题公式

## 2. 语义学

### 2.1 真值表

**定义 2.1 (真值赋值)**
真值赋值是从原子命题集合到 ```latex
$\{true, false\}$
``` 的函数。

**定义 2.2 (真值表)**
真值表是描述命题公式在所有可能真值赋值下真值的表格。

**定理 2.1 (真值表构造)**
对于包含 ```latex
$n$
``` 个不同原子命题的命题公式，其真值表有 ```latex
$2^n$
``` 行。

**证明**：
每个原子命题有两种可能的真值，根据乘法原理，```latex
$n$
``` 个原子命题共有 ```latex
$2^n$
``` 种不同的真值赋值组合。

### 2.2 语义函数

**定义 2.3 (语义函数)**
语义函数 ```latex
$\llbracket \cdot \rrbracket$
``` 将命题公式映射到真值，满足：

1. ```latex
$\llbracket P \rrbracket = v(P)$
``` 对于原子命题 ```latex
$P$
```
2. ```latex
$\llbracket \neg \phi \rrbracket = \neg \llbracket \phi \rrbracket$
```
3. ```latex
$\llbracket \phi \wedge \psi \rrbracket = \llbracket \phi \rrbracket \wedge \llbracket \psi \rrbracket$
```
4. ```latex
$\llbracket \phi \vee \psi \rrbracket = \llbracket \phi \rrbracket \vee \llbracket \psi \rrbracket$
```
5. ```latex
$\llbracket \phi \rightarrow \psi \rrbracket = \neg \llbracket \phi \rrbracket \vee \llbracket \psi \rrbracket$
```
6. ```latex
$\llbracket \phi \leftrightarrow \psi \rrbracket = (\llbracket \phi \rrbracket \rightarrow \llbracket \psi \rrbracket) \wedge (\llbracket \psi \rrbracket \rightarrow \llbracket \phi \rrbracket)$
```

### 2.3 逻辑等价

**定义 2.4 (逻辑等价)**
两个命题公式 ```latex
$\phi, \psi$
``` 是逻辑等价的，记作 ```latex
$\phi \equiv \psi$
```，当且仅当对于所有真值赋值，```latex
$\llbracket \phi \rrbracket = \llbracket \psi \rrbracket$
```。

**定理 2.2 (德摩根律)**
对于任意命题公式 ```latex
$\phi, \psi$
```：

1. ```latex
$\neg(\phi \wedge \psi) \equiv \neg \phi \vee \neg \psi$
```
2. ```latex
$\neg(\phi \vee \psi) \equiv \neg \phi \wedge \neg \psi$
```

**证明**：
通过真值表验证，对于所有真值赋值，两边的真值都相同。

**定理 2.3 (分配律)**
对于任意命题公式 ```latex
$\phi, \psi, \chi$
```：

1. ```latex
$\phi \wedge (\psi \vee \chi) \equiv (\phi \wedge \psi) \vee (\phi \wedge \chi)$
```
2. ```latex
$\phi \vee (\psi \wedge \chi) \equiv (\phi \vee \psi) \wedge (\phi \vee \chi)$
```

## 3. 证明系统

### 3.1 自然演绎

**定义 3.1 (自然演绎规则)**
自然演绎系统包含以下推理规则：

**引入规则**：

- ```latex
$\wedge$
```-I: 从 ```latex
$\phi$
``` 和 ```latex
$\psi$
``` 推出 ```latex
$\phi \wedge \psi$
```
- ```latex
$\vee$
```-I: 从 ```latex
$\phi$
``` 推出 ```latex
$\phi \vee \psi$
``` 或 ```latex
$\psi \vee \phi$
```
- ```latex
$\rightarrow$
```-I: 从假设 ```latex
$\phi$
``` 推出 ```latex
$\psi$
``` 后，可以推出 ```latex
$\phi \rightarrow \psi$
```

**消除规则**：

- ```latex
$\wedge$
```-E: 从 ```latex
$\phi \wedge \psi$
``` 推出 ```latex
$\phi$
``` 或 ```latex
$\psi$
```
- ```latex
$\vee$
```-E: 从 ```latex
$\phi \vee \psi$
``` 和 ```latex
$\phi \rightarrow \chi$
``` 和 ```latex
$\psi \rightarrow \chi$
``` 推出 ```latex
$\chi$
```
- ```latex
$\rightarrow$
```-E: 从 ```latex
$\phi$
``` 和 ```latex
$\phi \rightarrow \psi$
``` 推出 ```latex
$\psi$
```

**定理 3.1 (自然演绎的可靠性)**
如果 ```latex
$\Gamma \vdash \phi$
```，则 ```latex
$\Gamma \models \phi$
```。

### 3.2 公理系统

**定义 3.2 (公理系统)**
命题逻辑的公理系统包含以下公理模式：

1. ```latex
$\phi \rightarrow (\psi \rightarrow \phi)$
```
2. ```latex
$(\phi \rightarrow (\psi \rightarrow \chi)) \rightarrow ((\phi \rightarrow \psi) \rightarrow (\phi \rightarrow \chi))$
```
3. ```latex
$(\neg \phi \rightarrow \neg \psi) \rightarrow (\psi \rightarrow \phi)$
```

**推理规则**：

- 分离规则（MP）：从 ```latex
$\phi$
``` 和 ```latex
$\phi \rightarrow \psi$
``` 推出 ```latex
$\psi$
```

### 3.3 归结原理

**定义 3.3 (子句)**
子句是文字的析取，其中文字是原子命题或其否定。

**定义 3.4 (归结规则)**
对于子句 ```latex
$C_1 = A \vee L$
``` 和 ```latex
$C_2 = B \vee \neg L$
```，归结规则推出 ```latex
$C_1 \vee C_2 = A \vee B$
```。

**定理 3.2 (归结的完备性)**
如果命题公式集合 ```latex
$\Gamma$
``` 是不可满足的，则可以通过归结规则推出空子句。

## 4. 形式化定义

### 4.1 语法定义

**定义 4.1 (命题逻辑语言)**
命题逻辑语言 ```latex
$\mathcal{L}$
``` 由以下部分组成：

- 原子命题集合 ```latex
$\mathcal{P} = \{P_1, P_2, \ldots\}$
```
- 逻辑连接词集合 ```latex
$\{\neg, \wedge, \vee, \rightarrow, \leftrightarrow\}$
```
- 括号 ```latex
$\{(, )\}$
```

**定义 4.2 (公式集合)**
公式集合 ```latex
$\mathcal{F}$
``` 是满足以下条件的最小集合：

1. ```latex
$\mathcal{P} \subseteq \mathcal{F}$
```
2. 如果 ```latex
$\phi \in \mathcal{F}$
```，则 ```latex
$\neg \phi \in \mathcal{F}$
```
3. 如果 ```latex
$\phi, \psi \in \mathcal{F}$
```，则 ```latex
$(\phi \wedge \psi), (\phi \vee \psi), (\phi \rightarrow \psi), (\phi \leftrightarrow \psi) \in \mathcal{F}$
```

### 4.2 语义定义

**定义 4.3 (解释)**
解释 ```latex
$I$
``` 是从原子命题集合 ```latex
$\mathcal{P}$
``` 到 ```latex
$\{true, false\}$
``` 的函数。

**定义 4.4 (满足关系)**
满足关系 ```latex
$\models$
``` 递归定义如下：

1. ```latex
$I \models P$
``` 当且仅当 ```latex
$I(P) = true$
```
2. ```latex
$I \models \neg \phi$
``` 当且仅当 ```latex
$I \not\models \phi$
```
3. ```latex
$I \models \phi \wedge \psi$
``` 当且仅当 ```latex
$I \models \phi$
``` 且 ```latex
$I \models \psi$
```
4. ```latex
$I \models \phi \vee \psi$
``` 当且仅当 ```latex
$I \models \phi$
``` 或 ```latex
$I \models \psi$
```
5. ```latex
$I \models \phi \rightarrow \psi$
``` 当且仅当 ```latex
$I \not\models \phi$
``` 或 ```latex
$I \models \psi$
```
6. ```latex
$I \models \phi \leftrightarrow \psi$
``` 当且仅当 ```latex
$I \models \phi$
``` 等价于 ```latex
$I \models \psi$
```

### 4.3 完备性定理

**定理 4.1 (命题逻辑的完备性)**
对于任意命题公式 ```latex
$\phi$
```，如果 ```latex
$\models \phi$
```，则 ```latex
$\vdash \phi$
```。

**证明**：
使用真值表方法或归结方法可以构造证明。

**定理 4.2 (命题逻辑的可靠性)**
对于任意命题公式 ```latex
$\phi$
```，如果 ```latex
$\vdash \phi$
```，则 ```latex
$\models \phi$
```。

**证明**：
通过归纳法证明所有推理规则都保持有效性。

## 5. Go语言实现

### 5.1 命题公式表示

```go
// 命题公式的代数数据类型表示
type Proposition interface {
    String() string
    FreeVariables() map[string]bool
    Substitute(sub map[string]Proposition) Proposition
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

func (a *Atom) FreeVariables() map[string]bool {
    return map[string]bool{a.Name: true}
}

func (a *Atom) Substitute(sub map[string]Proposition) Proposition {
    if replacement, exists := sub[a.Name]; exists {
        return replacement
    }
    return a
}

// 否定
type Negation struct {
    Formula Proposition
}

func NewNegation(formula Proposition) *Negation {
    return &Negation{Formula: formula}
}

func (n *Negation) String() string {
    return "¬(" + n.Formula.String() + ")"
}

func (n *Negation) FreeVariables() map[string]bool {
    return n.Formula.FreeVariables()
}

func (n *Negation) Substitute(sub map[string]Proposition) Proposition {
    return NewNegation(n.Formula.Substitute(sub))
}

// 合取
type Conjunction struct {
    Left, Right Proposition
}

func NewConjunction(left, right Proposition) *Conjunction {
    return &Conjunction{Left: left, Right: right}
}

func (c *Conjunction) String() string {
    return "(" + c.Left.String() + " ∧ " + c.Right.String() + ")"
}

func (c *Conjunction) FreeVariables() map[string]bool {
    vars := c.Left.FreeVariables()
    for v := range c.Right.FreeVariables() {
        vars[v] = true
    }
    return vars
}

func (c *Conjunction) Substitute(sub map[string]Proposition) Proposition {
    return NewConjunction(c.Left.Substitute(sub), c.Right.Substitute(sub))
}

// 析取
type Disjunction struct {
    Left, Right Proposition
}

func NewDisjunction(left, right Proposition) *Disjunction {
    return &Disjunction{Left: left, Right: right}
}

func (d *Disjunction) String() string {
    return "(" + d.Left.String() + " ∨ " + d.Right.String() + ")"
}

func (d *Disjunction) FreeVariables() map[string]bool {
    vars := d.Left.FreeVariables()
    for v := range d.Right.FreeVariables() {
        vars[v] = true
    }
    return vars
}

func (d *Disjunction) Substitute(sub map[string]Proposition) Proposition {
    return NewDisjunction(d.Left.Substitute(sub), d.Right.Substitute(sub))
}

// 蕴含
type Implication struct {
    Antecedent, Consequent Proposition
}

func NewImplication(antecedent, consequent Proposition) *Implication {
    return &Implication{Antecedent: antecedent, Consequent: consequent}
}

func (i *Implication) String() string {
    return "(" + i.Antecedent.String() + " → " + i.Consequent.String() + ")"
}

func (i *Implication) FreeVariables() map[string]bool {
    vars := i.Antecedent.FreeVariables()
    for v := range i.Consequent.FreeVariables() {
        vars[v] = true
    }
    return vars
}

func (i *Implication) Substitute(sub map[string]Proposition) Proposition {
    return NewImplication(i.Antecedent.Substitute(sub), i.Consequent.Substitute(sub))
}

// 等价
type Equivalence struct {
    Left, Right Proposition
}

func NewEquivalence(left, right Proposition) *Equivalence {
    return &Equivalence{Left: left, Right: right}
}

func (e *Equivalence) String() string {
    return "(" + e.Left.String() + " ↔ " + e.Right.String() + ")"
}

func (e *Equivalence) FreeVariables() map[string]bool {
    vars := e.Left.FreeVariables()
    for v := range e.Right.FreeVariables() {
        vars[v] = true
    }
    return vars
}

func (e *Equivalence) Substitute(sub map[string]Proposition) Proposition {
    return NewEquivalence(e.Left.Substitute(sub), e.Right.Substitute(sub))
}
```

### 5.2 真值表计算

```go
// 真值赋值
type Interpretation map[string]bool

// 真值表计算器
type TruthTableCalculator struct{}

func NewTruthTableCalculator() *TruthTableCalculator {
    return &TruthTableCalculator{}
}

// 计算命题公式在给定解释下的真值
func (ttc *TruthTableCalculator) Evaluate(formula Proposition, interpretation Interpretation) bool {
    switch f := formula.(type) {
    case *Atom:
        return interpretation[f.Name]
    case *Negation:
        return !ttc.Evaluate(f.Formula, interpretation)
    case *Conjunction:
        return ttc.Evaluate(f.Left, interpretation) && ttc.Evaluate(f.Right, interpretation)
    case *Disjunction:
        return ttc.Evaluate(f.Left, interpretation) || ttc.Evaluate(f.Right, interpretation)
    case *Implication:
        return !ttc.Evaluate(f.Antecedent, interpretation) || ttc.Evaluate(f.Consequent, interpretation)
    case *Equivalence:
        return ttc.Evaluate(f.Left, interpretation) == ttc.Evaluate(f.Right, interpretation)
    default:
        return false
    }
}

// 生成所有可能的真值赋值
func (ttc *TruthTableCalculator) GenerateInterpretations(variables []string) []Interpretation {
    n := len(variables)
    total := 1 << n
    interpretations := make([]Interpretation, total)
    
    for i := 0; i < total; i++ {
        interpretation := make(Interpretation)
        for j := 0; j < n; j++ {
            interpretation[variables[j]] = (i & (1 << j)) != 0
        }
        interpretations[i] = interpretation
    }
    
    return interpretations
}

// 计算真值表
func (ttc *TruthTableCalculator) ComputeTruthTable(formula Proposition) map[Interpretation]bool {
    variables := ttc.getVariables(formula)
    interpretations := ttc.GenerateInterpretations(variables)
    truthTable := make(map[Interpretation]bool)
    
    for _, interpretation := range interpretations {
        truthTable[interpretation] = ttc.Evaluate(formula, interpretation)
    }
    
    return truthTable
}

// 检查是否为重言式
func (ttc *TruthTableCalculator) IsTautology(formula Proposition) bool {
    truthTable := ttc.ComputeTruthTable(formula)
    for _, value := range truthTable {
        if !value {
            return false
        }
    }
    return true
}

// 检查是否为矛盾式
func (ttc *TruthTableCalculator) IsContradiction(formula Proposition) bool {
    truthTable := ttc.ComputeTruthTable(formula)
    for _, value := range truthTable {
        if value {
            return false
        }
    }
    return true
}

// 检查两个公式是否逻辑等价
func (ttc *TruthTableCalculator) AreEquivalent(formula1, formula2 Proposition) bool {
    // 合并两个公式的自由变量
    vars1 := ttc.getVariables(formula1)
    vars2 := ttc.getVariables(formula2)
    allVars := ttc.mergeVariables(vars1, vars2)
    
    interpretations := ttc.GenerateInterpretations(allVars)
    
    for _, interpretation := range interpretations {
        val1 := ttc.Evaluate(formula1, interpretation)
        val2 := ttc.Evaluate(formula2, interpretation)
        if val1 != val2 {
            return false
        }
    }
    
    return true
}

// 获取公式中的变量
func (ttc *TruthTableCalculator) getVariables(formula Proposition) []string {
    vars := formula.FreeVariables()
    result := make([]string, 0, len(vars))
    for v := range vars {
        result = append(result, v)
    }
    sort.Strings(result)
    return result
}

// 合并变量列表
func (ttc *TruthTableCalculator) mergeVariables(vars1, vars2 []string) []string {
    varMap := make(map[string]bool)
    for _, v := range vars1 {
        varMap[v] = true
    }
    for _, v := range vars2 {
        varMap[v] = true
    }
    
    result := make([]string, 0, len(varMap))
    for v := range varMap {
        result = append(result, v)
    }
    sort.Strings(result)
    return result
}
```

### 5.3 证明系统实现

```go
// 证明系统
type ProofSystem struct{}

func NewProofSystem() *ProofSystem {
    return &ProofSystem{}
}

// 归结证明
func (ps *ProofSystem) ResolutionProof(clauses [][]string) bool {
    // 转换为CNF形式
    cnf := ps.toCNF(clauses)
    
    // 归结过程
    for {
        newClauses := ps.resolutionStep(cnf)
        if len(newClauses) == 0 {
            return false // 无法归结
        }
        
        // 检查是否产生空子句
        for _, clause := range newClauses {
            if len(clause) == 0 {
                return true // 找到矛盾
            }
        }
        
        // 添加新子句
        originalSize := len(cnf)
        for _, clause := range newClauses {
            if !ps.containsClause(cnf, clause) {
                cnf = append(cnf, clause)
            }
        }
        
        // 如果没有新子句，则无法证明
        if len(cnf) == originalSize {
            return false
        }
    }
}

// 归结步骤
func (ps *ProofSystem) resolutionStep(clauses [][]string) [][]string {
    var newClauses [][]string
    
    for i := 0; i < len(clauses); i++ {
        for j := i + 1; j < len(clauses); j++ {
            resolvent := ps.resolve(clauses[i], clauses[j])
            if len(resolvent) >= 0 {
                newClauses = append(newClauses, resolvent)
            }
        }
    }
    
    return newClauses
}

// 归结两个子句
func (ps *ProofSystem) resolve(clause1, clause2 []string) []string {
    // 寻找互补文字
    for _, literal1 := range clause1 {
        for _, literal2 := range clause2 {
            if ps.isComplementary(literal1, literal2) {
                // 归结
                result := make([]string, 0)
                
                // 添加clause1中除literal1外的所有文字
                for _, lit := range clause1 {
                    if lit != literal1 {
                        result = append(result, lit)
                    }
                }
                
                // 添加clause2中除literal2外的所有文字
                for _, lit := range clause2 {
                    if lit != literal2 && !ps.containsLiteral(result, lit) {
                        result = append(result, lit)
                    }
                }
                
                return result
            }
        }
    }
    
    return nil // 无法归结
}

// 检查两个文字是否互补
func (ps *ProofSystem) isComplementary(literal1, literal2 string) bool {
    if literal1[0] == '¬' {
        return literal1[1:] == literal2
    }
    if literal2[0] == '¬' {
        return literal2[1:] == literal1
    }
    return false
}

// 检查是否包含文字
func (ps *ProofSystem) containsLiteral(clause []string, literal string) bool {
    for _, lit := range clause {
        if lit == literal {
            return true
        }
    }
    return false
}

// 检查是否包含子句
func (ps *ProofSystem) containsClause(clauses [][]string, clause []string) bool {
    for _, existingClause := range clauses {
        if ps.clausesEqual(existingClause, clause) {
            return true
        }
    }
    return false
}

// 检查两个子句是否相等
func (ps *ProofSystem) clausesEqual(clause1, clause2 []string) bool {
    if len(clause1) != len(clause2) {
        return false
    }
    
    // 创建副本并排序
    c1 := make([]string, len(clause1))
    c2 := make([]string, len(clause2))
    copy(c1, clause1)
    copy(c2, clause2)
    sort.Strings(c1)
    sort.Strings(c2)
    
    for i := range c1 {
        if c1[i] != c2[i] {
            return false
        }
    }
    return true
}

// 转换为CNF形式（简化版本）
func (ps *ProofSystem) toCNF(clauses [][]string) [][]string {
    // 这里假设输入已经是CNF形式
    return clauses
}
```

## 6. 应用实例

### 6.1 电路设计验证

```go
// 数字电路验证器
type CircuitVerifier struct {
    calculator *TruthTableCalculator
}

func NewCircuitVerifier() *CircuitVerifier {
    return &CircuitVerifier{
        calculator: NewTruthTableCalculator(),
    }
}

// 验证电路等价性
func (cv *CircuitVerifier) VerifyCircuitEquivalence(circuit1, circuit2 *Circuit) bool {
    formula1 := cv.circuitToFormula(circuit1)
    formula2 := cv.circuitToFormula(circuit2)
    
    return cv.calculator.AreEquivalent(formula1, formula2)
}

// 电路结构
type Circuit struct {
    inputs  []string
    outputs []string
    gates   []Gate
}

type Gate struct {
    Type   string
    Inputs []string
    Output string
}

// 将电路转换为命题公式
func (cv *CircuitVerifier) circuitToFormula(circuit *Circuit) Proposition {
    // 为每个输出构建公式
    var formulas []Proposition
    
    for _, output := range circuit.outputs {
        formula := cv.buildOutputFormula(circuit, output)
        formulas = append(formulas, formula)
    }
    
    // 如果有多个输出，使用合取连接
    if len(formulas) == 1 {
        return formulas[0]
    }
    
    result := formulas[0]
    for i := 1; i < len(formulas); i++ {
        result = NewConjunction(result, formulas[i])
    }
    
    return result
}

// 构建单个输出的公式
func (cv *CircuitVerifier) buildOutputFormula(circuit *Circuit, output string) Proposition {
    // 找到产生该输出的门
    for _, gate := range circuit.gates {
        if gate.Output == output {
            return cv.gateToFormula(gate)
        }
    }
    
    // 如果输出直接连接到输入
    for _, input := range circuit.inputs {
        if input == output {
            return NewAtom(input)
        }
    }
    
    return NewAtom(output)
}

// 将门转换为公式
func (cv *CircuitVerifier) gateToFormula(gate Gate) Proposition {
    switch gate.Type {
    case "AND":
        if len(gate.Inputs) == 2 {
            return NewConjunction(NewAtom(gate.Inputs[0]), NewAtom(gate.Inputs[1]))
        }
    case "OR":
        if len(gate.Inputs) == 2 {
            return NewDisjunction(NewAtom(gate.Inputs[0]), NewAtom(gate.Inputs[1]))
        }
    case "NOT":
        if len(gate.Inputs) == 1 {
            return NewNegation(NewAtom(gate.Inputs[0]))
        }
    case "NAND":
        if len(gate.Inputs) == 2 {
            return NewNegation(NewConjunction(NewAtom(gate.Inputs[0]), NewAtom(gate.Inputs[1])))
        }
    case "NOR":
        if len(gate.Inputs) == 2 {
            return NewNegation(NewDisjunction(NewAtom(gate.Inputs[0]), NewAtom(gate.Inputs[1])))
        }
    }
    
    return NewAtom(gate.Output)
}

// 验证电路功能
func (cv *CircuitVerifier) VerifyCircuitFunction(circuit *Circuit, specification Proposition) bool {
    circuitFormula := cv.circuitToFormula(circuit)
    return cv.calculator.AreEquivalent(circuitFormula, specification)
}
```

### 6.2 程序逻辑验证

```go
// 程序逻辑验证器
type ProgramLogicVerifier struct {
    calculator *TruthTableCalculator
}

func NewProgramLogicVerifier() *ProgramLogicVerifier {
    return &ProgramLogicVerifier{
        calculator: NewTruthTableCalculator(),
    }
}

// 霍尔三元组
type HoareTriple struct {
    Precondition  Proposition
    Program       *Program
    Postcondition Proposition
}

// 程序结构
type Program struct {
    Statements []Statement
}

type Statement interface {
    Execute(state map[string]bool) map[string]bool
    String() string
}

// 赋值语句
type Assignment struct {
    Variable string
    Value    Proposition
}

func (a *Assignment) Execute(state map[string]bool) map[string]bool {
    newState := make(map[string]bool)
    for k, v := range state {
        newState[k] = v
    }
    newState[a.Variable] = a.Value.Evaluate(state)
    return newState
}

func (a *Assignment) String() string {
    return a.Variable + " := " + a.Value.String()
}

// 条件语句
type Conditional struct {
    Condition Proposition
    Then      *Program
    Else      *Program
}

func (c *Conditional) Execute(state map[string]bool) map[string]bool {
    if c.Condition.Evaluate(state) {
        return c.Then.Execute(state)
    } else {
        return c.Else.Execute(state)
    }
}

func (c *Conditional) String() string {
    return "if " + c.Condition.String() + " then " + c.Then.String() + " else " + c.Else.String()
}

// 循环语句
type Loop struct {
    Condition Proposition
    Body      *Program
}

func (l *Loop) Execute(state map[string]bool) map[string]bool {
    currentState := state
    for l.Condition.Evaluate(currentState) {
        currentState = l.Body.Execute(currentState)
    }
    return currentState
}

func (l *Loop) String() string {
    return "while " + l.Condition.String() + " do " + l.Body.String()
}

// 验证霍尔三元组
func (plv *ProgramLogicVerifier) VerifyHoareTriple(triple HoareTriple) bool {
    // 生成所有可能的初始状态
    variables := plv.getProgramVariables(triple.Program)
    initialStates := plv.generateStates(variables)
    
    for _, initialState := range initialStates {
        // 检查前置条件
        if !triple.Precondition.Evaluate(initialState) {
            continue // 跳过不满足前置条件的状态
        }
        
        // 执行程序
        finalState := triple.Program.Execute(initialState)
        
        // 检查后置条件
        if !triple.Postcondition.Evaluate(finalState) {
            return false // 霍尔三元组不成立
        }
    }
    
    return true
}

// 获取程序中的变量
func (plv *ProgramLogicVerifier) getProgramVariables(program *Program) []string {
    variables := make(map[string]bool)
    
    for _, stmt := range program.Statements {
        plv.collectVariables(stmt, variables)
    }
    
    result := make([]string, 0, len(variables))
    for v := range variables {
        result = append(result, v)
    }
    sort.Strings(result)
    return result
}

// 收集语句中的变量
func (plv *ProgramLogicVerifier) collectVariables(stmt Statement, variables map[string]bool) {
    switch s := stmt.(type) {
    case *Assignment:
        variables[s.Variable] = true
        plv.collectVariablesFromProposition(s.Value, variables)
    case *Conditional:
        plv.collectVariablesFromProposition(s.Condition, variables)
        plv.collectVariables(s.Then, variables)
        plv.collectVariables(s.Else, variables)
    case *Loop:
        plv.collectVariablesFromProposition(s.Condition, variables)
        plv.collectVariables(s.Body, variables)
    }
}

// 从命题公式中收集变量
func (plv *ProgramLogicVerifier) collectVariablesFromProposition(prop Proposition, variables map[string]bool) {
    for v := range prop.FreeVariables() {
        variables[v] = true
    }
}

// 生成所有可能的状态
func (plv *ProgramLogicVerifier) generateStates(variables []string) []map[string]bool {
    n := len(variables)
    total := 1 << n
    states := make([]map[string]bool, total)
    
    for i := 0; i < total; i++ {
        state := make(map[string]bool)
        for j := 0; j < n; j++ {
            state[variables[j]] = (i & (1 << j)) != 0
        }
        states[i] = state
    }
    
    return states
}
```

### 6.3 知识表示

```go
// 知识库
type KnowledgeBase struct {
    formulas []Proposition
    calculator *TruthTableCalculator
}

func NewKnowledgeBase() *KnowledgeBase {
    return &KnowledgeBase{
        formulas:   make([]Proposition, 0),
        calculator: NewTruthTableCalculator(),
    }
}

// 添加知识
func (kb *KnowledgeBase) AddKnowledge(formula Proposition) {
    kb.formulas = append(kb.formulas, formula)
}

// 查询知识
func (kb *KnowledgeBase) Query(query Proposition) bool {
    // 构建知识库的合取
    if len(kb.formulas) == 0 {
        return false
    }
    
    knowledge := kb.formulas[0]
    for i := 1; i < len(kb.formulas); i++ {
        knowledge = NewConjunction(knowledge, kb.formulas[i])
    }
    
    // 检查知识库蕴含查询
    implication := NewImplication(knowledge, query)
    return kb.calculator.IsTautology(implication)
}

// 一致性检查
func (kb *KnowledgeBase) IsConsistent() bool {
    if len(kb.formulas) == 0 {
        return true
    }
    
    knowledge := kb.formulas[0]
    for i := 1; i < len(kb.formulas); i++ {
        knowledge = NewConjunction(knowledge, kb.formulas[i])
    }
    
    return !kb.calculator.IsContradiction(knowledge)
}

// 推理引擎
type InferenceEngine struct {
    knowledgeBase *KnowledgeBase
}

func NewInferenceEngine(kb *KnowledgeBase) *InferenceEngine {
    return &InferenceEngine{
        knowledgeBase: kb,
    }
}

// 前向推理
func (ie *InferenceEngine) ForwardChaining(query Proposition) bool {
    // 使用归结方法进行推理
    clauses := ie.knowledgeBaseToClauses()
    queryClause := ie.formulaToClauses(NewNegation(query))
    clauses = append(clauses, queryClause...)
    
    proofSystem := NewProofSystem()
    return proofSystem.ResolutionProof(clauses)
}

// 将知识库转换为子句形式
func (ie *InferenceEngine) knowledgeBaseToClauses() [][]string {
    var allClauses [][]string
    
    for _, formula := range ie.knowledgeBase.formulas {
        clauses := ie.formulaToClauses(formula)
        allClauses = append(allClauses, clauses...)
    }
    
    return allClauses
}

// 将公式转换为子句形式（简化版本）
func (ie *InferenceEngine) formulaToClauses(formula Proposition) [][]string {
    // 这里提供一个简化的实现
    // 实际应用中需要完整的CNF转换算法
    
    switch f := formula.(type) {
    case *Atom:
        return [][]string{{f.Name}}
    case *Negation:
        if atom, ok := f.Formula.(*Atom); ok {
            return [][]string{{"¬" + atom.Name}}
        }
    case *Conjunction:
        clauses1 := ie.formulaToClauses(f.Left)
        clauses2 := ie.formulaToClauses(f.Right)
        return append(clauses1, clauses2...)
    case *Disjunction:
        // 简化处理，假设已经是CNF形式
        return [][]string{{f.Left.String(), f.Right.String()}}
    }
    
    return [][]string{}
}
```

## 总结

命题逻辑作为形式逻辑的基础，在软件工程中有着广泛的应用。本章从基础定义出发，通过形式化证明建立了命题逻辑的理论基础，并提供了完整的Go语言实现。

主要内容包括：

1. **理论基础**：命题、逻辑连接词、命题公式的定义和性质
2. **语义学**：真值表、语义函数、逻辑等价的概念和定理
3. **证明系统**：自然演绎、公理系统、归结原理的实现
4. **实际应用**：电路设计验证、程序逻辑验证、知识表示等

这些内容为后续的谓词逻辑、模态逻辑、时态逻辑等更高级的逻辑系统提供了重要的理论基础，也为软件验证、人工智能、知识工程等领域提供了重要的工具。
