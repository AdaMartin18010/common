# 02-逻辑学 (Logic)

## 概述

逻辑学是形式化方法的基础，为软件工程和计算科学提供严格的推理框架。本文档介绍逻辑学的基本概念、形式化系统以及在Go语言中的应用。

## 目录

- [02-逻辑学 (Logic)](#02-逻辑学-logic)
  - [概述](#概述)
  - [目录](#目录)
  - [1. 命题逻辑 (Propositional Logic)](#1-命题逻辑-propositional-logic)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 真值表](#12-真值表)
    - [1.3 Go语言实现](#13-go语言实现)
    - [1.4 逻辑等价律](#14-逻辑等价律)
  - [2. 谓词逻辑 (Predicate Logic)](#2-谓词逻辑-predicate-logic)
    - [2.1 基本概念](#21-基本概念)
    - [2.2 Go语言实现](#22-go语言实现)
    - [2.3 量词等价律](#23-量词等价律)
  - [3. 形式化证明 (Formal Proofs)](#3-形式化证明-formal-proofs)
    - [3.1 自然演绎系统](#31-自然演绎系统)
    - [3.2 Go语言实现](#32-go语言实现)
  - [4. Go语言中的逻辑实现](#4-go语言中的逻辑实现)
    - [4.1 类型安全的逻辑编程](#41-类型安全的逻辑编程)
    - [4.2 约束逻辑编程](#42-约束逻辑编程)
  - [总结](#总结)

---

## 1. 命题逻辑 (Propositional Logic)

### 1.1 基本概念

**定义 1.1.1** (命题)
命题是一个具有确定真值的陈述句，用符号 $P, Q, R$ 等表示。

**定义 1.1.2** (逻辑连接词)

- 否定 (Negation): $\neg P$
- 合取 (Conjunction): $P \land Q$
- 析取 (Disjunction): $P \lor Q$
- 蕴含 (Implication): $P \rightarrow Q$
- 等价 (Equivalence): $P \leftrightarrow Q$

### 1.2 真值表

| P | Q | ¬P | P∧Q | P∨Q | P→Q | P↔Q |
|---|---|----|-----|-----|-----|-----|
| T | T | F  | T   | T   | T   | T   |
| T | F | F  | F   | T   | F   | F   |
| F | T | T  | F   | T   | T   | F   |
| F | F | T  | F   | F   | T   | T   |

### 1.3 Go语言实现

```go
package logic

// Proposition 表示一个命题
type Proposition interface {
    Evaluate(assignment map[string]bool) bool
    String() string
}

// AtomicProposition 原子命题
type AtomicProposition struct {
    Name string
}

func (p *AtomicProposition) Evaluate(assignment map[string]bool) bool {
    return assignment[p.Name]
}

func (p *AtomicProposition) String() string {
    return p.Name
}

// Negation 否定
type Negation struct {
    Prop Proposition
}

func (n *Negation) Evaluate(assignment map[string]bool) bool {
    return !n.Prop.Evaluate(assignment)
}

func (n *Negation) String() string {
    return "¬(" + n.Prop.String() + ")"
}

// Conjunction 合取
type Conjunction struct {
    Left, Right Proposition
}

func (c *Conjunction) Evaluate(assignment map[string]bool) bool {
    return c.Left.Evaluate(assignment) && c.Right.Evaluate(assignment)
}

func (c *Conjunction) String() string {
    return "(" + c.Left.String() + " ∧ " + c.Right.String() + ")"
}

// Disjunction 析取
type Disjunction struct {
    Left, Right Proposition
}

func (d *Disjunction) Evaluate(assignment map[string]bool) bool {
    return d.Left.Evaluate(assignment) || d.Right.Evaluate(assignment)
}

func (d *Disjunction) String() string {
    return "(" + d.Left.String() + " ∨ " + d.Right.String() + ")"
}

// Implication 蕴含
type Implication struct {
    Antecedent, Consequent Proposition
}

func (i *Implication) Evaluate(assignment map[string]bool) bool {
    return !i.Antecedent.Evaluate(assignment) || i.Consequent.Evaluate(assignment)
}

func (i *Implication) String() string {
    return "(" + i.Antecedent.String() + " → " + i.Consequent.String() + ")"
}

// TruthTable 生成真值表
func TruthTable(prop Proposition, variables []string) [][]bool {
    rows := 1 << len(variables)
    result := make([][]bool, rows)
    
    for i := 0; i < rows; i++ {
        assignment := make(map[string]bool)
        for j, varName := range variables {
            assignment[varName] = (i & (1 << j)) != 0
        }
        
        row := make([]bool, len(variables)+1)
        for j, varName := range variables {
            row[j] = assignment[varName]
        }
        row[len(variables)] = prop.Evaluate(assignment)
        result[i] = row
    }
    
    return result
}
```

### 1.4 逻辑等价律

**定理 1.4.1** (德摩根律)
$\neg(P \land Q) \equiv \neg P \lor \neg Q$
$\neg(P \lor Q) \equiv \neg P \land \neg Q$

**证明**:
通过真值表验证：

| P | Q | ¬(P∧Q) | ¬P∨¬Q | ¬(P∨Q) | ¬P∧¬Q |
|---|---|--------|-------|--------|-------|
| T | T | F      | F     | F      | F     |
| T | F | T      | T     | F      | F     |
| F | T | T      | T     | F      | F     |
| F | F | T      | T     | T      | T     |

---

## 2. 谓词逻辑 (Predicate Logic)

### 2.1 基本概念

**定义 2.1.1** (谓词)
谓词是描述对象性质的函数，用 $P(x), Q(x,y)$ 等表示。

**定义 2.1.2** (量词)

- 全称量词 (Universal): $\forall x P(x)$
- 存在量词 (Existential): $\exists x P(x)$

### 2.2 Go语言实现

```go
package predicate

import "reflect"

// Predicate 表示一个谓词
type Predicate[T any] func(T) bool

// Universal 全称量词
func Universal[T any](pred Predicate[T], domain []T) bool {
    for _, item := range domain {
        if !pred(item) {
            return false
        }
    }
    return true
}

// Existential 存在量词
func Existential[T any](pred Predicate[T], domain []T) bool {
    for _, item := range domain {
        if pred(item) {
            return true
        }
    }
    return false
}

// And 谓词合取
func And[T any](pred1, pred2 Predicate[T]) Predicate[T] {
    return func(x T) bool {
        return pred1(x) && pred2(x)
    }
}

// Or 谓词析取
func Or[T any](pred1, pred2 Predicate[T]) Predicate[T] {
    return func(x T) bool {
        return pred1(x) || pred2(x)
    }
}

// Not 谓词否定
func Not[T any](pred Predicate[T]) Predicate[T] {
    return func(x T) bool {
        return !pred(x)
    }
}

// Compose 谓词复合
func Compose[T, U any](f func(T) U, g Predicate[U]) Predicate[T] {
    return func(x T) bool {
        return g(f(x))
    }
}

// Example: 数学谓词
func IsEven(n int) bool {
    return n%2 == 0
}

func IsPositive(n int) bool {
    return n > 0
}

func IsPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}
```

### 2.3 量词等价律

**定理 2.3.1** (量词否定律)
$\neg \forall x P(x) \equiv \exists x \neg P(x)$
$\neg \exists x P(x) \equiv \forall x \neg P(x)$

**证明**:
设论域为 $D = \{a_1, a_2, ..., a_n\}$

$\neg \forall x P(x) \equiv \neg (P(a_1) \land P(a_2) \land ... \land P(a_n))$
$\equiv \neg P(a_1) \lor \neg P(a_2) \lor ... \lor \neg P(a_n)$
$\equiv \exists x \neg P(x)$

---

## 3. 形式化证明 (Formal Proofs)

### 3.1 自然演绎系统

**定义 3.1.1** (推理规则)

1. **引入规则** (Introduction Rules)
   - $\land$-I: 从 $A$ 和 $B$ 推出 $A \land B$
   - $\lor$-I: 从 $A$ 推出 $A \lor B$
   - $\rightarrow$-I: 从假设 $A$ 推出 $B$ 得到 $A \rightarrow B$

2. **消除规则** (Elimination Rules)
   - $\land$-E: 从 $A \land B$ 推出 $A$ 或 $B$
   - $\lor$-E: 从 $A \lor B$ 和 $A \rightarrow C$ 和 $B \rightarrow C$ 推出 $C$
   - $\rightarrow$-E: 从 $A \rightarrow B$ 和 $A$ 推出 $B$

### 3.2 Go语言实现

```go
package proof

import (
    "fmt"
    "strings"
)

// Formula 表示逻辑公式
type Formula interface {
    String() string
    IsAtomic() bool
}

// Proof 表示证明
type Proof struct {
    Premises []Formula
    Conclusion Formula
    Steps []ProofStep
}

// ProofStep 证明步骤
type ProofStep struct {
    Number int
    Formula Formula
    Justification string
    Dependencies []int
}

// NaturalDeduction 自然演绎系统
type NaturalDeduction struct {
    proofs []Proof
}

// NewNaturalDeduction 创建自然演绎系统
func NewNaturalDeduction() *NaturalDeduction {
    return &NaturalDeduction{
        proofs: make([]Proof, 0),
    }
}

// Prove 证明一个公式
func (nd *NaturalDeduction) Prove(premises []Formula, conclusion Formula) *Proof {
    proof := &Proof{
        Premises: premises,
        Conclusion: conclusion,
        Steps: make([]ProofStep, 0),
    }
    
    // 添加前提
    for i, premise := range premises {
        proof.Steps = append(proof.Steps, ProofStep{
            Number: i + 1,
            Formula: premise,
            Justification: "Premise",
            Dependencies: []int{},
        })
    }
    
    return proof
}

// AddStep 添加证明步骤
func (p *Proof) AddStep(formula Formula, justification string, deps []int) {
    step := ProofStep{
        Number: len(p.Steps) + 1,
        Formula: formula,
        Justification: justification,
        Dependencies: deps,
    }
    p.Steps = append(p.Steps, step)
}

// PrintProof 打印证明
func (p *Proof) PrintProof() string {
    var sb strings.Builder
    
    sb.WriteString("Proof:\n")
    sb.WriteString("Premises:\n")
    for i, premise := range p.Premises {
        sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, premise.String()))
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

// Example: 证明 P∧Q → Q∧P
func ExampleCommutativity() {
    nd := NewNaturalDeduction()
    
    // 创建原子命题
    p := &AtomicProposition{Name: "P"}
    q := &AtomicProposition{Name: "Q"}
    
    // 前提: P∧Q
    premise := &Conjunction{Left: p, Right: q}
    
    // 结论: Q∧P
    conclusion := &Conjunction{Left: q, Right: p}
    
    // 证明
    proof := nd.Prove([]Formula{premise}, &Implication{
        Antecedent: premise,
        Consequent: conclusion,
    })
    
    // 添加证明步骤
    proof.AddStep(premise, "Premise", []int{})
    proof.AddStep(q, "∧-E", []int{1})
    proof.AddStep(p, "∧-E", []int{1})
    proof.AddStep(&Conjunction{Left: q, Right: p}, "∧-I", []int{2, 3})
    
    fmt.Println(proof.PrintProof())
}
```

---

## 4. Go语言中的逻辑实现

### 4.1 类型安全的逻辑编程

```go
package logic_programming

import (
    "fmt"
    "reflect"
)

// LogicVar 逻辑变量
type LogicVar[T any] struct {
    value T
    bound bool
}

// NewLogicVar 创建逻辑变量
func NewLogicVar[T any]() *LogicVar[T] {
    return &LogicVar[T]{bound: false}
}

// Bind 绑定值
func (lv *LogicVar[T]) Bind(value T) bool {
    if lv.bound {
        return reflect.DeepEqual(lv.value, value)
    }
    lv.value = value
    lv.bound = true
    return true
}

// Value 获取值
func (lv *LogicVar[T]) Value() (T, bool) {
    return lv.value, lv.bound
}

// Unify 统一两个逻辑变量
func Unify[T comparable](v1, v2 *LogicVar[T]) bool {
    val1, bound1 := v1.Value()
    val2, bound2 := v2.Value()
    
    if bound1 && bound2 {
        return reflect.DeepEqual(val1, val2)
    } else if bound1 {
        return v2.Bind(val1)
    } else if bound2 {
        return v1.Bind(val2)
    }
    return true
}

// LogicProgram 逻辑程序
type LogicProgram struct {
    rules []Rule
}

// Rule 逻辑规则
type Rule struct {
    Head Predicate
    Body []Predicate
}

// Predicate 谓词
type Predicate struct {
    Name string
    Args []interface{}
}

// Query 查询
func (lp *LogicProgram) Query(goal Predicate) []map[string]interface{} {
    solutions := make([]map[string]interface{}, 0)
    
    // 简单的回溯搜索实现
    var backtrack func(Predicate, map[string]interface{}) bool
    backtrack = func(current Predicate, env map[string]interface{}) bool {
        // 查找匹配的规则
        for _, rule := range lp.rules {
            if rule.Head.Name == current.Name {
                // 尝试统一参数
                if lp.unify(rule.Head.Args, current.Args, env) {
                    // 如果规则体为空，找到解
                    if len(rule.Body) == 0 {
                        solution := make(map[string]interface{})
                        for k, v := range env {
                            solution[k] = v
                        }
                        solutions = append(solutions, solution)
                        return true
                    }
                    
                    // 递归处理规则体
                    success := true
                    for _, bodyPred := range rule.Body {
                        if !backtrack(bodyPred, env) {
                            success = false
                            break
                        }
                    }
                    if success {
                        return true
                    }
                }
            }
        }
        return false
    }
    
    backtrack(goal, make(map[string]interface{}))
    return solutions
}

// unify 统一算法
func (lp *LogicProgram) unify(args1, args2 []interface{}, env map[string]interface{}) bool {
    if len(args1) != len(args2) {
        return false
    }
    
    for i, arg1 := range args1 {
        arg2 := args2[i]
        
        // 处理变量
        if varName, ok := arg1.(string); ok && varName[0] == 'X' {
            if val, exists := env[varName]; exists {
                if !reflect.DeepEqual(val, arg2) {
                    return false
                }
            } else {
                env[varName] = arg2
            }
        } else if varName, ok := arg2.(string); ok && varName[0] == 'X' {
            if val, exists := env[varName]; exists {
                if !reflect.DeepEqual(val, arg1) {
                    return false
                }
            } else {
                env[varName] = arg1
            }
        } else {
            // 常量比较
            if !reflect.DeepEqual(arg1, arg2) {
                return false
            }
        }
    }
    
    return true
}

// Example: 家族关系逻辑程序
func ExampleFamilyLogic() {
    lp := &LogicProgram{
        rules: []Rule{
            {
                Head: Predicate{Name: "parent", Args: []interface{}{"X1", "X2"}},
                Body: []Predicate{
                    {Name: "father", Args: []interface{}{"X1", "X2"}},
                },
            },
            {
                Head: Predicate{Name: "parent", Args: []interface{}{"X1", "X2"}},
                Body: []Predicate{
                    {Name: "mother", Args: []interface{}{"X1", "X2"}},
                },
            },
            {
                Head: Predicate{Name: "grandparent", Args: []interface{}{"X1", "X2"}},
                Body: []Predicate{
                    {Name: "parent", Args: []interface{}{"X1", "X3"}},
                    {Name: "parent", Args: []interface{}{"X3", "X2"}},
                },
            },
        },
    }
    
    // 查询祖父关系
    solutions := lp.Query(Predicate{
        Name: "grandparent",
        Args: []interface{}{"X1", "X2"},
    })
    
    fmt.Printf("Found %d solutions for grandparent(X1, X2)\n", len(solutions))
    for i, solution := range solutions {
        fmt.Printf("Solution %d: %v\n", i+1, solution)
    }
}
```

### 4.2 约束逻辑编程

```go
package constraint_logic

import (
    "fmt"
    "math"
)

// Constraint 约束
type Constraint interface {
    Satisfied(assignment map[string]float64) bool
    Variables() []string
}

// LinearConstraint 线性约束
type LinearConstraint struct {
    Coefficients map[string]float64
    Operator     string // "==", "<=", ">="
    Constant     float64
}

func (lc *LinearConstraint) Satisfied(assignment map[string]float64) bool {
    sum := 0.0
    for varName, coeff := range lc.Coefficients {
        if val, exists := assignment[varName]; exists {
            sum += coeff * val
        } else {
            return false // 变量未赋值
        }
    }
    
    switch lc.Operator {
    case "==":
        return math.Abs(sum-lc.Constant) < 1e-10
    case "<=":
        return sum <= lc.Constant
    case ">=":
        return sum >= lc.Constant
    default:
        return false
    }
}

func (lc *LinearConstraint) Variables() []string {
    vars := make([]string, 0, len(lc.Coefficients))
    for varName := range lc.Coefficients {
        vars = append(vars, varName)
    }
    return vars
}

// ConstraintSolver 约束求解器
type ConstraintSolver struct {
    constraints []Constraint
    variables   []string
}

// NewConstraintSolver 创建约束求解器
func NewConstraintSolver() *ConstraintSolver {
    return &ConstraintSolver{
        constraints: make([]Constraint, 0),
        variables:   make([]string, 0),
    }
}

// AddConstraint 添加约束
func (cs *ConstraintSolver) AddConstraint(constraint Constraint) {
    cs.constraints = append(cs.constraints, constraint)
    
    // 收集变量
    for _, varName := range constraint.Variables() {
        found := false
        for _, existing := range cs.variables {
            if existing == varName {
                found = true
                break
            }
        }
        if !found {
            cs.variables = append(cs.variables, varName)
        }
    }
}

// Solve 求解约束
func (cs *ConstraintSolver) Solve() []map[string]float64 {
    solutions := make([]map[string]float64, 0)
    
    // 简单的回溯搜索
    var backtrack func(int, map[string]float64)
    backtrack = func(index int, assignment map[string]float64) {
        if index == len(cs.variables) {
            // 检查所有约束
            valid := true
            for _, constraint := range cs.constraints {
                if !constraint.Satisfied(assignment) {
                    valid = false
                    break
                }
            }
            if valid {
                solution := make(map[string]float64)
                for k, v := range assignment {
                    solution[k] = v
                }
                solutions = append(solutions, solution)
            }
            return
        }
        
        varName := cs.variables[index]
        
        // 尝试不同的值
        for val := 0.0; val <= 10.0; val += 1.0 {
            assignment[varName] = val
            
            // 检查部分约束
            partialValid := true
            for _, constraint := range cs.constraints {
                if !constraint.Satisfied(assignment) {
                    partialValid = false
                    break
                }
            }
            
            if partialValid {
                backtrack(index+1, assignment)
            }
        }
    }
    
    backtrack(0, make(map[string]float64))
    return solutions
}

// Example: 线性规划问题
func ExampleLinearProgramming() {
    solver := NewConstraintSolver()
    
    // 添加约束: x + y <= 10
    solver.AddConstraint(&LinearConstraint{
        Coefficients: map[string]float64{"x": 1, "y": 1},
        Operator:     "<=",
        Constant:     10,
    })
    
    // 添加约束: 2x + y <= 15
    solver.AddConstraint(&LinearConstraint{
        Coefficients: map[string]float64{"x": 2, "y": 1},
        Operator:     "<=",
        Constant:     15,
    })
    
    // 添加约束: x >= 0, y >= 0
    solver.AddConstraint(&LinearConstraint{
        Coefficients: map[string]float64{"x": 1},
        Operator:     ">=",
        Constant:     0,
    })
    solver.AddConstraint(&LinearConstraint{
        Coefficients: map[string]float64{"y": 1},
        Operator:     ">=",
        Constant:     0,
    })
    
    solutions := solver.Solve()
    fmt.Printf("Found %d solutions\n", len(solutions))
    for i, solution := range solutions {
        fmt.Printf("Solution %d: x=%.1f, y=%.1f\n", i+1, solution["x"], solution["y"])
    }
}
```

---

## 总结

本文档介绍了逻辑学的基本概念和在Go语言中的实现：

1. **命题逻辑** - 基本连接词、真值表、逻辑等价律
2. **谓词逻辑** - 量词、谓词函数、量词等价律
3. **形式化证明** - 自然演绎系统、推理规则
4. **Go语言实现** - 类型安全的逻辑编程、约束逻辑编程

这些形式化方法为软件工程提供了严格的推理基础，确保程序的正确性和可靠性。

---

**相关链接**:

- [01-集合论 (Set Theory)](01-Set-Theory.md)
- [03-图论 (Graph Theory)](03-Graph-Theory.md)
- [04-概率论 (Probability Theory)](04-Probability-Theory.md)
- [02-形式化验证 (Formal Verification)](../02-Formal-Verification/README.md)
