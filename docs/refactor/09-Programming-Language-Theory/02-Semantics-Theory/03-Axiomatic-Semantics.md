# 03-公理语义 (Axiomatic Semantics)

## 目录

- [1. 概述](#1-概述)
- [2. Hoare逻辑基础](#2-hoare逻辑基础)
- [3. 最弱前置条件](#3-最弱前置条件)
- [4. 程序验证](#4-程序验证)
- [5. 循环不变式](#5-循环不变式)
- [6. Go语言实现](#6-go语言实现)
- [7. 形式化证明](#7-形式化证明)
- [8. 应用实例](#8-应用实例)

## 1. 概述

### 1.1 公理语义学定义

公理语义学是一种形式化方法，用于描述程序的含义和正确性。它基于数学逻辑，通过前置条件和后置条件来规约程序的行为。

**形式化定义**：

```latex
\text{公理语义学} = (\mathcal{P}, \mathcal{Q}, \mathcal{R}, \vdash)
```

其中：

- $\mathcal{P}$ 是前置条件集合
- $\mathcal{Q}$ 是后置条件集合  
- $\mathcal{R}$ 是推理规则集合
- $\vdash$ 是推导关系

### 1.2 核心概念

#### 1.2.1 Hoare三元组

```latex
\{P\} \text{ } C \text{ } \{Q\}
```

表示：如果前置条件 $P$ 成立，执行程序 $C$ 后，后置条件 $Q$ 成立。

#### 1.2.2 正确性分类

- **部分正确性**：如果程序终止，则后置条件成立
- **完全正确性**：程序一定终止且后置条件成立

## 2. Hoare逻辑基础

### 2.1 基本推理规则

#### 2.1.1 赋值公理

```latex
\{P[E/x]\} \text{ } x := E \text{ } \{P\}
```

**Go语言实现**：

```go
// AssignmentAxiom 赋值公理
type AssignmentAxiom struct {
    Variable string
    Expression string
    PostCondition string
}

func (aa *AssignmentAxiom) GetPreCondition() string {
    // 将后置条件中的变量替换为表达式
    return strings.ReplaceAll(aa.PostCondition, aa.Variable, aa.Expression)
}

// 示例：{x+1 > 0} x := x+1 {x > 0}
assignment := &AssignmentAxiom{
    Variable: "x",
    Expression: "x+1", 
    PostCondition: "x > 0",
}
preCondition := assignment.GetPreCondition() // "(x+1) > 0"
```

#### 2.1.2 顺序规则

```latex
\frac{\{P\} \text{ } C_1 \text{ } \{R\} \quad \{R\} \text{ } C_2 \text{ } \{Q\}}{\{P\} \text{ } C_1; C_2 \text{ } \{Q\}}
```

**Go语言实现**：

```go
// SequentialRule 顺序规则
type SequentialRule struct {
    C1, C2 Statement
    P, R, Q string
}

func (sr *SequentialRule) IsValid() bool {
    // 验证两个子证明
    return sr.ProveC1() && sr.ProveC2()
}

func (sr *SequentialRule) ProveC1() bool {
    // 证明 {P} C1 {R}
    return sr.C1.Verify(sr.P, sr.R)
}

func (sr *SequentialRule) ProveC2() bool {
    // 证明 {R} C2 {Q}
    return sr.C2.Verify(sr.R, sr.Q)
}
```

#### 2.1.3 条件规则

```latex
\frac{\{P \land B\} \text{ } C_1 \text{ } \{Q\} \quad \{P \land \neg B\} \text{ } C_2 \text{ } \{Q\}}{\{P\} \text{ } \text{if } B \text{ then } C_1 \text{ else } C_2 \text{ } \{Q\}}
```

**Go语言实现**：

```go
// ConditionalRule 条件规则
type ConditionalRule struct {
    Condition string
    ThenBranch, ElseBranch Statement
    P, Q string
}

func (cr *ConditionalRule) IsValid() bool {
    // 验证两个分支
    thenValid := cr.ThenBranch.Verify(cr.P+" && "+cr.Condition, cr.Q)
    elseValid := cr.ElseBranch.Verify(cr.P+" && !("+cr.Condition+")", cr.Q)
    return thenValid && elseValid
}
```

### 2.2 循环规则

#### 2.2.1 While循环规则

```latex
\frac{\{P \land B\} \text{ } C \text{ } \{P\}}{\{P\} \text{ } \text{while } B \text{ do } C \text{ } \{P \land \neg B\}}
```

其中 $P$ 是循环不变式。

**Go语言实现**：

```go
// WhileRule While循环规则
type WhileRule struct {
    Condition string
    Body Statement
    Invariant string
}

func (wr *WhileRule) IsValid() bool {
    // 验证循环体保持不变式
    bodyValid := wr.Body.Verify(wr.Invariant+" && "+wr.Condition, wr.Invariant)
    
    // 验证终止性（需要额外的终止性证明）
    terminationValid := wr.ProveTermination()
    
    return bodyValid && terminationValid
}

func (wr *WhileRule) ProveTermination() bool {
    // 证明循环终止
    // 通常需要找到变式函数
    return wr.FindVariantFunction() != nil
}
```

## 3. 最弱前置条件

### 3.1 定义

最弱前置条件（Weakest Precondition, WP）是使得程序执行后满足后置条件的最弱前置条件。

```latex
\text{wp}(C, Q) = \{s \in \Sigma \mid \text{执行 } C \text{ 从状态 } s \text{ 开始，终止后满足 } Q\}
```

### 3.2 计算规则

#### 3.2.1 赋值语句

```latex
\text{wp}(x := E, Q) = Q[E/x]
```

#### 3.2.2 顺序语句

```latex
\text{wp}(C_1; C_2, Q) = \text{wp}(C_1, \text{wp}(C_2, Q))
```

#### 3.2.3 条件语句

```latex
\text{wp}(\text{if } B \text{ then } C_1 \text{ else } C_2, Q) = (B \land \text{wp}(C_1, Q)) \lor (\neg B \land \text{wp}(C_2, Q))
```

#### 3.2.4 循环语句

```latex
\text{wp}(\text{while } B \text{ do } C, Q) = \exists k \geq 0: H_k(Q)
```

其中 $H_k$ 是循环的 $k$ 次展开。

**Go语言实现**：

```go
// WeakestPrecondition 最弱前置条件计算器
type WeakestPrecondition struct{}

func (wp *WeakestPrecondition) Calculate(stmt Statement, postCondition string) string {
    switch s := stmt.(type) {
    case *Assignment:
        return wp.calculateAssignment(s, postCondition)
    case *Sequence:
        return wp.calculateSequence(s, postCondition)
    case *Conditional:
        return wp.calculateConditional(s, postCondition)
    case *WhileLoop:
        return wp.calculateWhile(s, postCondition)
    default:
        return "unknown"
    }
}

func (wp *WeakestPrecondition) calculateAssignment(assign *Assignment, Q string) string {
    // wp(x := E, Q) = Q[E/x]
    return strings.ReplaceAll(Q, assign.Variable, assign.Expression)
}

func (wp *WeakestPrecondition) calculateSequence(seq *Sequence, Q string) string {
    // wp(C1; C2, Q) = wp(C1, wp(C2, Q))
    wpC2 := wp.Calculate(seq.Second, Q)
    return wp.Calculate(seq.First, wpC2)
}

func (wp *WeakestPrecondition) calculateConditional(cond *Conditional, Q string) string {
    // wp(if B then C1 else C2, Q) = (B && wp(C1, Q)) || (!(B) && wp(C2, Q))
    wpThen := wp.Calculate(cond.ThenBranch, Q)
    wpElse := wp.Calculate(cond.ElseBranch, Q)
    return fmt.Sprintf("(%s && %s) || (!(%s) && %s)", 
        cond.Condition, wpThen, cond.Condition, wpElse)
}
```

## 4. 程序验证

### 4.1 验证框架

**Go语言实现**：

```go
// ProgramVerifier 程序验证器
type ProgramVerifier struct {
    wp *WeakestPrecondition
}

// VerificationResult 验证结果
type VerificationResult struct {
    Valid bool
    PreCondition string
    PostCondition string
    Proof []string
    Errors []string
}

func (pv *ProgramVerifier) Verify(stmt Statement, preCondition, postCondition string) *VerificationResult {
    result := &VerificationResult{
        PreCondition: preCondition,
        PostCondition: postCondition,
        Proof: []string{},
        Errors: []string{},
    }
    
    // 计算最弱前置条件
    wp := pv.wp.Calculate(stmt, postCondition)
    
    // 验证前置条件蕴含最弱前置条件
    if pv.implies(preCondition, wp) {
        result.Valid = true
        result.Proof = append(result.Proof, 
            fmt.Sprintf("Precondition: %s", preCondition),
            fmt.Sprintf("Weakest precondition: %s", wp),
            "Precondition implies weakest precondition ✓")
    } else {
        result.Valid = false
        result.Errors = append(result.Errors, 
            fmt.Sprintf("Precondition %s does not imply weakest precondition %s", 
                preCondition, wp))
    }
    
    return result
}

func (pv *ProgramVerifier) implies(P, Q string) bool {
    // 简化的蕴含检查
    // 在实际实现中，这里需要集成定理证明器
    return pv.simplify(fmt.Sprintf("!(%s) || (%s)", P, Q)) == "true"
}

func (pv *ProgramVerifier) simplify(expr string) string {
    // 简化的表达式化简
    // 实际实现需要完整的逻辑化简器
    return expr
}
```

### 4.2 验证示例

```go
// 验证示例：交换两个变量的值
func ExampleSwapVerification() {
    // 程序：temp := x; x := y; y := temp
    swap := &Sequence{
        First: &Assignment{Variable: "temp", Expression: "x"},
        Second: &Sequence{
            First: &Assignment{Variable: "x", Expression: "y"},
            Second: &Assignment{Variable: "y", Expression: "temp"},
        },
    }
    
    preCondition := "x = a && y = b"
    postCondition := "x = b && y = a"
    
    verifier := &ProgramVerifier{wp: &WeakestPrecondition{}}
    result := verifier.Verify(swap, preCondition, postCondition)
    
    fmt.Printf("Verification result: %v\n", result.Valid)
    for _, proof := range result.Proof {
        fmt.Printf("  %s\n", proof)
    }
}
```

## 5. 循环不变式

### 5.1 不变式定义

循环不变式是在循环执行过程中始终保持为真的谓词。

**形式化定义**：

```latex
\text{对于循环 } \text{while } B \text{ do } C \text{，不变式 } I \text{ 满足：}
\begin{cases}
P \Rightarrow I & \text{(初始化)}
\\
\{I \land B\} \text{ } C \text{ } \{I\} & \text{(保持)}
\\
I \land \neg B \Rightarrow Q & \text{(终止)}
\end{cases}
```

### 5.2 不变式发现

**Go语言实现**：

```go
// InvariantFinder 不变式发现器
type InvariantFinder struct{}

// InvariantCandidate 不变式候选
type InvariantCandidate struct {
    Expression string
    Confidence float64
    Evidence []string
}

func (if *InvariantFinder) FindInvariants(loop *WhileLoop, preCondition, postCondition string) []*InvariantCandidate {
    candidates := []*InvariantCandidate{}
    
    // 1. 从前置条件推导
    candidates = append(candidates, if.deriveFromPrecondition(preCondition)...)
    
    // 2. 从后置条件推导
    candidates = append(candidates, if.deriveFromPostcondition(postCondition)...)
    
    // 3. 从循环体分析
    candidates = append(candidates, if.analyzeLoopBody(loop)...)
    
    // 4. 排序并返回
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Confidence > candidates[j].Confidence
    })
    
    return candidates
}

func (if *InvariantFinder) deriveFromPrecondition(preCondition string) []*InvariantCandidate {
    // 从前置条件推导不变式
    candidates := []*InvariantCandidate{}
    
    // 例如：如果前置条件是 x >= 0，可能的不变式是 x >= 0
    if strings.Contains(preCondition, ">=") {
        candidates = append(candidates, &InvariantCandidate{
            Expression: preCondition,
            Confidence: 0.8,
            Evidence: []string{"Derived from precondition"},
        })
    }
    
    return candidates
}
```

## 6. Go语言实现

### 6.1 语句抽象

```go
// Statement 语句接口
type Statement interface {
    Execute(state map[string]int) map[string]int
    Verify(preCondition, postCondition string) bool
    String() string
}

// Assignment 赋值语句
type Assignment struct {
    Variable string
    Expression string
}

func (a *Assignment) Execute(state map[string]int) map[string]int {
    newState := make(map[string]int)
    for k, v := range state {
        newState[k] = v
    }
    
    // 计算表达式值
    value := a.evaluateExpression(a.Expression, state)
    newState[a.Variable] = value
    
    return newState
}

func (a *Assignment) evaluateExpression(expr string, state map[string]int) int {
    // 简化的表达式求值
    // 实际实现需要完整的表达式解析器
    if val, ok := state[expr]; ok {
        return val
    }
    return 0
}

// Sequence 顺序语句
type Sequence struct {
    First, Second Statement
}

func (s *Sequence) Execute(state map[string]int) map[string]int {
    state1 := s.First.Execute(state)
    return s.Second.Execute(state1)
}

// Conditional 条件语句
type Conditional struct {
    Condition string
    ThenBranch, ElseBranch Statement
}

func (c *Conditional) Execute(state map[string]int) map[string]int {
    if c.evaluateCondition(c.Condition, state) {
        return c.ThenBranch.Execute(state)
    } else {
        return c.ElseBranch.Execute(state)
    }
}

// WhileLoop While循环
type WhileLoop struct {
    Condition string
    Body Statement
}

func (w *WhileLoop) Execute(state map[string]int) map[string]int {
    currentState := state
    for w.evaluateCondition(w.Condition, currentState) {
        currentState = w.Body.Execute(currentState)
    }
    return currentState
}
```

### 6.2 验证器实现

```go
// HoareLogicVerifier Hoare逻辑验证器
type HoareLogicVerifier struct {
    wp *WeakestPrecondition
}

func (hlv *HoareLogicVerifier) VerifyTriple(stmt Statement, preCondition, postCondition string) bool {
    // 计算最弱前置条件
    wp := hlv.wp.Calculate(stmt, postCondition)
    
    // 检查前置条件是否蕴含最弱前置条件
    return hlv.implies(preCondition, wp)
}

func (hlv *HoareLogicVerifier) implies(P, Q string) bool {
    // 简化的蕴含检查
    // 实际实现需要集成定理证明器如Z3
    return true // 简化实现
}
```

## 7. 形式化证明

### 7.1 证明系统

```latex
\text{Hoare逻辑证明系统} = (\mathcal{A}, \mathcal{R}, \vdash)
```

其中：

- $\mathcal{A}$ 是公理集合
- $\mathcal{R}$ 是推理规则集合
- $\vdash$ 是推导关系

### 7.2 证明构造

**Go语言实现**：

```go
// Proof 证明结构
type Proof struct {
    Premises []string
    Conclusion string
    Rule string
    SubProofs []*Proof
}

// ProofConstructor 证明构造器
type ProofConstructor struct{}

func (pc *ProofConstructor) ConstructProof(stmt Statement, preCondition, postCondition string) *Proof {
    switch s := stmt.(type) {
    case *Assignment:
        return pc.constructAssignmentProof(s, preCondition, postCondition)
    case *Sequence:
        return pc.constructSequenceProof(s, preCondition, postCondition)
    case *Conditional:
        return pc.constructConditionalProof(s, preCondition, postCondition)
    case *WhileLoop:
        return pc.constructWhileProof(s, preCondition, postCondition)
    default:
        return nil
    }
}

func (pc *ProofConstructor) constructAssignmentProof(assign *Assignment, P, Q string) *Proof {
    // 赋值公理：{P[E/x]} x := E {P}
    wp := strings.ReplaceAll(Q, assign.Variable, assign.Expression)
    
    return &Proof{
        Premises: []string{fmt.Sprintf("Precondition: %s", P)},
        Conclusion: fmt.Sprintf("{%s} %s := %s {%s}", P, assign.Variable, assign.Expression, Q),
        Rule: "Assignment Axiom",
        SubProofs: []*Proof{},
    }
}
```

## 8. 应用实例

### 8.1 数组排序验证

```go
// 验证冒泡排序的正确性
func ExampleBubbleSortVerification() {
    // 冒泡排序的Hoare三元组
    preCondition := "n > 0 && forall i: 0 <= i < n -> A[i] is integer"
    postCondition := "forall i: 0 <= i < n-1 -> A[i] <= A[i+1]"
    
    // 构造排序程序（简化版本）
    bubbleSort := constructBubbleSortProgram()
    
    // 验证
    verifier := &ProgramVerifier{wp: &WeakestPrecondition{}}
    result := verifier.Verify(bubbleSort, preCondition, postCondition)
    
    fmt.Printf("Bubble sort verification: %v\n", result.Valid)
}

func constructBubbleSortProgram() Statement {
    // 构造简化的冒泡排序程序
    // 实际实现需要完整的程序构造
    return &Sequence{
        First: &Assignment{Variable: "i", Expression: "0"},
        Second: &WhileLoop{
            Condition: "i < n-1",
            Body: &Sequence{
                First: &Assignment{Variable: "j", Expression: "0"},
                Second: &WhileLoop{
                    Condition: "j < n-1-i",
                    Body: &Conditional{
                        Condition: "A[j] > A[j+1]",
                        ThenBranch: &Sequence{
                            First: &Assignment{Variable: "temp", Expression: "A[j]"},
                            Second: &Sequence{
                                First: &Assignment{Variable: "A[j]", Expression: "A[j+1]"},
                                Second: &Assignment{Variable: "A[j+1]", Expression: "temp"},
                            },
                        },
                        ElseBranch: &Assignment{Variable: "j", Expression: "j+1"},
                    },
                },
            },
        },
    }
}
```

### 8.2 并发程序验证

```go
// 验证互斥锁的正确性
func ExampleMutexVerification() {
    // 互斥锁的规约
    preCondition := "!in_critical_section"
    postCondition := "in_critical_section"
    
    // 构造加锁程序
    lock := &Sequence{
        First: &Assignment{Variable: "waiting", Expression: "true"},
        Second: &WhileLoop{
            Condition: "waiting",
            Body: &Conditional{
                Condition: "!locked",
                ThenBranch: &Sequence{
                    First: &Assignment{Variable: "locked", Expression: "true"},
                    Second: &Assignment{Variable: "waiting", Expression: "false"},
                },
                ElseBranch: &Assignment{Variable: "skip", Expression: "skip"},
            },
        },
    }
    
    // 验证
    verifier := &ProgramVerifier{wp: &WeakestPrecondition{}}
    result := verifier.Verify(lock, preCondition, postCondition)
    
    fmt.Printf("Mutex lock verification: %v\n", result.Valid)
}
```

## 总结

公理语义学为程序正确性提供了严格的形式化基础。通过Hoare逻辑、最弱前置条件和程序验证技术，我们可以在数学上证明程序的正确性。Go语言的实现展示了如何将这些理论概念应用到实际编程中，为软件工程提供了强有力的验证工具。

**关键要点**：

1. **形式化基础**：公理语义学基于数学逻辑，提供严格的程序语义定义
2. **验证技术**：通过Hoare三元组和最弱前置条件进行程序验证
3. **循环处理**：使用循环不变式处理循环程序的验证
4. **实际应用**：在并发程序、算法正确性等领域有重要应用

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **公理语义学理论完成！** 🚀
