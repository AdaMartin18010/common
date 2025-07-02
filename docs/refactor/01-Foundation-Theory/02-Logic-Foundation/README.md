# 02-逻辑基础 (Logic Foundation)

## 目录

- [02-逻辑基础 (Logic Foundation)](#02-逻辑基础-logic-foundation)
  - [目录](#目录)
  - [1. 命题逻辑 (Propositional Logic)](#1-命题逻辑-propositional-logic)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 语法](#12-语法)
    - [1.3 语义](#13-语义)
    - [1.4 推理规则](#14-推理规则)
    - [2.4 推理规则](#24-推理规则)
    - [3.4 系统](#34-系统)
    - [4.3 分支时态逻辑](#43-分支时态逻辑)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 命题逻辑实现](#51-命题逻辑实现)
    - [5.2 真值表生成](#52-真值表生成)
    - [5.3 谓词逻辑实现](#53-谓词逻辑实现)
    - [5.4 时态逻辑实现](#54-时态逻辑实现)
    - [5.5 使用示例](#55-使用示例)
  - [6. 应用场景](#6-应用场景)
    - [6.1 程序验证](#61-程序验证)
    - [6.2 人工智能](#62-人工智能)
    - [6.3 数据库](#63-数据库)
    - [6.4 硬件设计](#64-硬件设计)
  - [7. 总结](#7-总结)
    - [关键要点](#关键要点)
    - [进一步研究方向](#进一步研究方向)

## 1. 命题逻辑 (Propositional Logic)

### 1.1 基本概念

**定义 1.1**: 命题
命题是一个可以判断真假的陈述句。

**定义 1.2**: 原子命题
原子命题是最基本的命题，不能再分解为更简单的命题。

**定义 1.3**: 复合命题
复合命题是由原子命题通过逻辑连接词构成的命题。

### 1.2 语法

**定义 1.4**: 命题逻辑语言
命题逻辑语言由以下部分组成：

- 命题变元集合 $P = \{p, q, r, \ldots\}$
- 逻辑连接词：$\neg$（否定）、$\wedge$（合取）、$\vee$（析取）、$\rightarrow$（蕴含）、$\leftrightarrow$（等价）
- 辅助符号：$($ 和 $)$

**定义 1.5**: 合式公式
合式公式（well-formed formula）递归定义如下：

1. 每个命题变元 $p \in P$ 是合式公式
2. 如果 $\phi$ 是合式公式，则 $\neg\phi$ 是合式公式
3. 如果 $\phi$ 和 $\psi$ 是合式公式，则 $(\phi \wedge \psi)$、$(\phi \vee \psi)$、$(\phi \rightarrow \psi)$、$(\phi \leftrightarrow \psi)$ 是合式公式
4. 只有通过上述规则构造的表达式才是合式公式

### 1.3 语义

**定义 1.6**: 真值赋值
真值赋值是一个函数 $v: P \rightarrow \{T, F\}$，其中 $T$ 表示真，$F$ 表示假。

**定义 1.7**: 真值函数
真值函数 $\overline{v}$ 递归定义如下：

1. $\overline{v}(p) = v(p)$ 对于 $p \in P$
2. $\overline{v}(\neg\phi) = T$ 当且仅当 $\overline{v}(\phi) = F$
3. $\overline{v}(\phi \wedge \psi) = T$ 当且仅当 $\overline{v}(\phi) = T$ 且 $\overline{v}(\psi) = T$
4. $\overline{v}(\phi \vee \psi) = T$ 当且仅当 $\overline{v}(\phi) = T$ 或 $\overline{v}(\psi) = T$
5. $\overline{v}(\phi \rightarrow \psi) = T$ 当且仅当 $\overline{v}(\phi) = F$ 或 $\overline{v}(\psi) = T$
6. $\overline{v}(\phi \leftrightarrow \psi) = T$ 当且仅当 $\overline{v}(\phi) = \overline{v}(\psi)$

**定义 1.8**: 重言式、矛盾式和可满足式

- 公式 $\phi$ 是重言式（tautology），如果对于所有真值赋值 $v$，$\overline{v}(\phi) = T$
- 公式 $\phi$ 是矛盾式（contradiction），如果对于所有真值赋值 $v$，$\overline{v}(\phi) = F$
- 公式 $\phi$ 是可满足式（satisfiable），如果存在真值赋值 $v$ 使得 $\overline{v}(\phi) = T$

### 1.4 推理规则

**定义 1.9**: 推理规则
常用的推理规则包括：

1. **假言推理（Modus Ponens）**：
   $$\frac{\phi \rightarrow \psi \quad \phi}{\psi}$$

2. **假言三段论**：
   $$\frac{\phi \rightarrow \psi \quad \psi \rightarrow \chi}{\phi \rightarrow \chi}$$

3. **合取引入**：
   $$\frac{\phi \quad \psi}{\phi \wedge \psi}$$$

4. **合取消除**：
   $```latex
\frac{\phi \wedge \psi}{\phi} \quad \frac{\phi \wedge \psi}{\psi}

```$

5. **析取引入**：
   $```latex
\frac{\phi}{\phi \vee \psi} \quad \frac{\psi}{\phi \vee \psi}
```$

**定理 1.1**: 德摩根律
对于任意公式 ```latex
\phi
``` 和 ```latex
\psi
```：

1. ```latex
\neg(\phi \wedge \psi) \equiv \neg\phi \vee \neg\psi
```

2. ```latex

\neg(\phi \vee \psi) \equiv \neg\phi \wedge \neg\psi

```

## 2. 谓词逻辑 (Predicate Logic)

### 2.1 基本概念

**定义 2.1**: 谓词
谓词是描述对象性质或关系的符号。

**定义 2.2**: 量词

- 全称量词 ```latex
\forall
```：表示"对于所有"
- 存在量词 ```latex
\exists
```：表示"存在"

**定义 2.3**: 项
项递归定义如下：

1. 变量和常量是项
2. 如果 ```latex
f
``` 是 ```latex
n
``` 元函数符号，```latex
t_1, \ldots, t_n
``` 是项，则 ```latex
f(t_1, \ldots, t_n)
``` 是项

### 2.2 语法

**定义 2.4**: 一阶逻辑语言
一阶逻辑语言由以下部分组成：

- 变量集合 ```latex
V = \{x, y, z, \ldots\}
```

- 常量集合 ```latex
C = \{a, b, c, \ldots\}

```
- 函数符号集合 ```latex
F = \{f, g, h, \ldots\}
```

- 谓词符号集合 ```latex
P = \{P, Q, R, \ldots\}

```
- 逻辑连接词：```latex
\neg, \wedge, \vee, \rightarrow, \leftrightarrow
```

- 量词：```latex
\forall, \exists

```
- 辅助符号：```latex
(
``` 和 ```latex
)
```

**定义 2.5**: 原子公式
如果 ```latex
P
``` 是 ```latex
n
``` 元谓词符号，```latex
t_1, \ldots, t_n
``` 是项，则 ```latex
P(t_1, \ldots, t_n)

``` 是原子公式。

**定义 2.6**: 合式公式
合式公式递归定义如下：

1. 每个原子公式是合式公式
2. 如果 ```latex
\phi
``` 是合式公式，则 ```latex
\neg\phi
``` 是合式公式
3. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是合式公式，则 ```latex
(\phi \wedge \psi)
```、```latex
(\phi \vee \psi)
```、```latex
(\phi \rightarrow \psi)
```、```latex
(\phi \leftrightarrow \psi)
``` 是合式公式
4. 如果 ```latex
\phi
``` 是合式公式，```latex
x
``` 是变量，则 ```latex
\forall x \phi
``` 和 ```latex
\exists x \phi
``` 是合式公式

### 2.3 语义

**定义 2.7**: 解释
解释 ```latex
\mathcal{I} = (D, \cdot^{\mathcal{I}})
``` 由以下部分组成：

- 论域 ```latex
D
```（非空集合）
- 解释函数 ```latex
\cdot^{\mathcal{I}}
```，将常量、函数符号和谓词符号映射到论域中的对象

**定义 2.8**: 赋值
赋值是一个函数 ```latex
\sigma: V \rightarrow D
```，将变量映射到论域中的对象。

**定义 2.9**: 满足关系
满足关系 ```latex
\models
``` 递归定义如下：

1. 

```latex
\mathcal{I}, \sigma \models P(t_1, \ldots, t_n)
```

当且仅当

```latex
(t_1^{\mathcal{I},\sigma}, \ldots, t_n^{\mathcal{I},\sigma}) \in P^{\mathcal{I}}
```

2.

```latex
\mathcal{I}, \sigma \models \neg\phi
```

当且仅当

```latex
\mathcal{I}, \sigma \not\models \phi
```

3.

```latex
\mathcal{I}, \sigma \models \phi \wedge \psi
```

当且仅当

```latex
\mathcal{I}, \sigma \models \phi
```

且

```latex
\mathcal{I}, \sigma \models \psi
```

4. ```latex

\mathcal{I}, \sigma \models \forall x \phi
``` 当且仅当对于所有 ```latex
d \in D
```，```latex
\mathcal{I}, \sigma[x \mapsto d] \models \phi

```
5. ```latex
\mathcal{I}, \sigma \models \exists x \phi
``` 当且仅当存在 ```latex
d \in D
``` 使得 ```latex
\mathcal{I}, \sigma[x \mapsto d] \models \phi
```

### 2.4 推理规则

**定义 2.10**: 谓词逻辑推理规则
除了命题逻辑的推理规则外，还有：

1. **全称消除**：
   $```latex
\frac{\forall x \phi}{\phi[t/x]}

```$
   其中 ```latex
t
``` 是项，```latex
\phi[t/x]
``` 表示将 ```latex
\phi
``` 中的 ```latex
x
``` 替换为 ```latex
t
```

2. **全称引入**：
   $```latex
\frac{\phi}{\forall x \phi}

```$
   其中 ```latex
x
``` 不在 ```latex
\phi
``` 的自由变量中出现

3. **存在引入**：
   $```latex
\frac{\phi[t/x]}{\exists x \phi}
```$

4. **存在消除**：
   $```latex
\frac{\exists x \phi \quad \phi[y/x] \vdash \psi}{\psi}
```$
   其中 ```latex
y
``` 是新的变量

## 3. 模态逻辑 (Modal Logic)

### 3.1 基本概念

**定义 3.1**: 模态算子

- ```latex
\Box
```：必然算子（necessarily）
- ```latex
\Diamond
```：可能算子（possibly）

**定义 3.2**: 模态公式
模态公式在命题逻辑基础上增加：

- 如果 ```latex
\phi
``` 是模态公式，则 ```latex
\Box\phi
``` 和 ```latex
\Diamond\phi
``` 是模态公式

### 3.2 语法

**定义 3.3**: 模态逻辑语言
模态逻辑语言在命题逻辑基础上增加模态算子 ```latex
\Box
``` 和 ```latex
\Diamond
```。

**定义 3.4**: 模态公式
模态公式递归定义如下：

1. 每个命题变元是模态公式
2. 如果 ```latex
\phi
``` 是模态公式，则 ```latex
\neg\phi
```、```latex
\Box\phi
```、```latex
\Diamond\phi
``` 是模态公式
3. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是模态公式，则 ```latex
(\phi \wedge \psi)
```、```latex
(\phi \vee \psi)
```、```latex
(\phi \rightarrow \psi)
```、```latex
(\phi \leftrightarrow \psi)
``` 是模态公式

### 3.3 语义

**定义 3.5**: 克里普克模型
克里普克模型是一个三元组 ```latex
\mathcal{M} = (W, R, V)
```，其中：

- ```latex
W
``` 是可能世界集合
- ```latex
R \subseteq W \times W
``` 是可达关系
- ```latex
V: W \times P \rightarrow \{T, F\}
``` 是赋值函数

**定义 3.6**: 模态逻辑满足关系
满足关系 ```latex
\models
``` 递归定义如下：

1. ```latex
\mathcal{M}, w \models p
``` 当且仅当 ```latex
V(w, p) = T
```

2. ```latex

\mathcal{M}, w \models \neg\phi
``` 当且仅当 ```latex
\mathcal{M}, w \not\models \phi

```
3. ```latex
\mathcal{M}, w \models \phi \wedge \psi
``` 当且仅当 ```latex
\mathcal{M}, w \models \phi
``` 且 ```latex
\mathcal{M}, w \models \psi
```

4. ```latex

\mathcal{M}, w \models \Box\phi
``` 当且仅当对于所有 ```latex
w'
``` 使得 ```latex
wRw'
```，```latex
\mathcal{M}, w' \models \phi

```
5. ```latex
\mathcal{M}, w \models \Diamond\phi
``` 当且仅当存在 ```latex
w'
``` 使得 ```latex
wRw'
``` 且 ```latex
\mathcal{M}, w' \models \phi
```

### 3.4 系统

**定义 3.7**: 模态逻辑系统
常见的模态逻辑系统包括：

1. **K系统**：最基本的模态逻辑系统
2. **T系统**：K + ```latex
\Box\phi \rightarrow \phi

```（自反性）
3. **S4系统**：T + ```latex
\Box\phi \rightarrow \Box\Box\phi
```（传递性）
4. **S5系统**：S4 + ```latex
\Diamond\phi \rightarrow \Box\Diamond\phi
```（欧几里得性）

## 4. 时态逻辑 (Temporal Logic)

### 4.1 基本概念

**定义 4.1**: 时态算子

- ```latex
G
```：全局算子（always）
- ```latex
F
```：未来算子（eventually）
- ```latex
X
```：下一个算子（next）
- ```latex
U
```：直到算子（until）

### 4.2 线性时态逻辑

**定义 4.2**: LTL语法
线性时态逻辑（Linear Temporal Logic, LTL）公式递归定义如下：

1. 每个命题变元是LTL公式
2. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是LTL公式，则 ```latex
\neg\phi
```、```latex
\phi \wedge \psi
```、```latex
\phi \vee \psi
```、```latex
\phi \rightarrow \psi
``` 是LTL公式
3. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是LTL公式，则 ```latex
X\phi
```、```latex
G\phi
```、```latex
F\phi
```、```latex
\phi U\psi
``` 是LTL公式

**定义 4.3**: LTL语义
LTL公式在无限序列 ```latex
\pi = \pi_0\pi_1\pi_2\ldots
``` 上的满足关系：

1. ```latex
\pi \models p
``` 当且仅当 ```latex
p \in \pi_0
```

2. ```latex

\pi \models \neg\phi
``` 当且仅当 ```latex
\pi \not\models \phi

```
3. ```latex
\pi \models \phi \wedge \psi
``` 当且仅当 ```latex
\pi \models \phi
``` 且 ```latex
\pi \models \psi
```

4. ```latex

\pi \models X\phi
``` 当且仅当 ```latex
\pi^1 \models \phi

```
5. ```latex
\pi \models G\phi
``` 当且仅当对于所有 ```latex
i \geq 0
```，```latex
\pi^i \models \phi
```

6. ```latex

\pi \models F\phi
``` 当且仅当存在 ```latex
i \geq 0
``` 使得 ```latex
\pi^i \models \phi

```
7. ```latex
\pi \models \phi U\psi
``` 当且仅当存在 ```latex
i \geq 0
``` 使得 ```latex
\pi^i \models \psi
``` 且对于所有 ```latex
0 \leq j < i
```，```latex
\pi^j \models \phi
```

### 4.3 分支时态逻辑

**定义 4.4**: CTL语法
计算树逻辑（Computation Tree Logic, CTL）公式递归定义如下：

1. 每个命题变元是CTL公式
2. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是CTL公式，则 ```latex
\neg\phi
```、```latex
\phi \wedge \psi
```、```latex
\phi \vee \psi
```、```latex
\phi \rightarrow \psi

``` 是CTL公式
3. 如果 ```latex
\phi
``` 和 ```latex
\psi
``` 是CTL公式，则 ```latex
AX\phi
```、```latex
EX\phi
```、```latex
AG\phi
```、```latex
EG\phi
```、```latex
AF\phi
```、```latex
EF\phi
```、```latex
A[\phi U\psi]
```、```latex
E[\phi U\psi]
``` 是CTL公式

### 4.4 模型检验

**定义 4.5**: 模型检验问题
给定一个系统模型 ```latex
\mathcal{M}
``` 和一个时态逻辑公式 ```latex
\phi
```，判断是否 ```latex
\mathcal{M} \models \phi
```。

**算法 4.1**: CTL模型检验算法

```go
func ModelCheckCTL(model *KripkeModel, formula CTLFormula) bool {
    // 递归计算满足公式的状态集合
    switch f := formula.(type) {
    case *Atomic:
        return model.SatisfyAtomic(f)
    case *Not:
        return !ModelCheckCTL(model, f.Formula)
    case *And:
        return ModelCheckCTL(model, f.Left) && ModelCheckCTL(model, f.Right)
    case *AX:
        return model.SatisfyAX(f.Formula)
    case *EX:
        return model.SatisfyEX(f.Formula)
    case *AG:
        return model.SatisfyAG(f.Formula)
    case *EG:
        return model.SatisfyEG(f.Formula)
    case *AF:
        return model.SatisfyAF(f.Formula)
    case *EF:
        return model.SatisfyEF(f.Formula)
    case *AU:
        return model.SatisfyAU(f.Left, f.Right)
    case *EU:
        return model.SatisfyEU(f.Left, f.Right)
    }
    return false
}
```

## 5. Go语言实现

### 5.1 命题逻辑实现

```go
package logic

import (
    "fmt"
    "strings"
)

// PropositionalFormula 表示命题逻辑公式
type PropositionalFormula interface {
    Evaluate(assignment map[string]bool) bool
    String() string
}

// Atomic 原子命题
type Atomic struct {
    Name string
}

func (a *Atomic) Evaluate(assignment map[string]bool) bool {
    return assignment[a.Name]
}

func (a *Atomic) String() string {
    return a.Name
}

// Not 否定
type Not struct {
    Formula PropositionalFormula
}

func (n *Not) Evaluate(assignment map[string]bool) bool {
    return !n.Formula.Evaluate(assignment)
}

func (n *Not) String() string {
    return fmt.Sprintf("¬(%s)", n.Formula)
}

// And 合取
type And struct {
    Left, Right PropositionalFormula
}

func (a *And) Evaluate(assignment map[string]bool) bool {
    return a.Left.Evaluate(assignment) && a.Right.Evaluate(assignment)
}

func (a *And) String() string {
    return fmt.Sprintf("(%s ∧ %s)", a.Left, a.Right)
}

// Or 析取
type Or struct {
    Left, Right PropositionalFormula
}

func (o *Or) Evaluate(assignment map[string]bool) bool {
    return o.Left.Evaluate(assignment) || o.Right.Evaluate(assignment)
}

func (o *Or) String() string {
    return fmt.Sprintf("(%s ∨ %s)", o.Left, o.Right)
}

// Implies 蕴含
type Implies struct {
    Left, Right PropositionalFormula
}

func (i *Implies) Evaluate(assignment map[string]bool) bool {
    return !i.Left.Evaluate(assignment) || i.Right.Evaluate(assignment)
}

func (i *Implies) String() string {
    return fmt.Sprintf("(%s → %s)", i.Left, i.Right)
}
```

### 5.2 真值表生成

```go
// TruthTable 生成真值表
func TruthTable(formula PropositionalFormula, variables []string) {
    n := len(variables)
    fmt.Printf("Truth Table for: %s\n", formula)
    fmt.Println(strings.Repeat("-", 50))
    
    // 打印表头
    for _, v := range variables {
        fmt.Printf("%s ", v)
    }
    fmt.Printf("| Result\n")
    fmt.Println(strings.Repeat("-", 50))
    
    // 生成所有可能的赋值
    for i := 0; i < (1 << n); i++ {
        assignment := make(map[string]bool)
        for j, v := range variables {
            assignment[v] = (i & (1 << j)) != 0
        }
        
        // 打印赋值
        for _, v := range variables {
            if assignment[v] {
                fmt.Printf("T ")
            } else {
                fmt.Printf("F ")
            }
        }
        
        // 计算结果
        result := formula.Evaluate(assignment)
        if result {
            fmt.Printf("| T\n")
        } else {
            fmt.Printf("| F\n")
        }
    }
}

// IsTautology 检查是否为重言式
func IsTautology(formula PropositionalFormula, variables []string) bool {
    n := len(variables)
    for i := 0; i < (1 << n); i++ {
        assignment := make(map[string]bool)
        for j, v := range variables {
            assignment[v] = (i & (1 << j)) != 0
        }
        if !formula.Evaluate(assignment) {
            return false
        }
    }
    return true
}

// IsSatisfiable 检查是否为可满足式
func IsSatisfiable(formula PropositionalFormula, variables []string) bool {
    n := len(variables)
    for i := 0; i < (1 << n); i++ {
        assignment := make(map[string]bool)
        for j, v := range variables {
            assignment[v] = (i & (1 << j)) != 0
        }
        if formula.Evaluate(assignment) {
            return true
        }
    }
    return false
}
```

### 5.3 谓词逻辑实现

```go
// PredicateFormula 表示谓词逻辑公式
type PredicateFormula interface {
    Evaluate(interpretation *Interpretation, assignment map[string]interface{}) bool
    String() string
}

// Interpretation 解释
type Interpretation struct {
    Domain []interface{}
    Predicates map[string]func([]interface{}) bool
    Functions map[string]func([]interface{}) interface{}
    Constants map[string]interface{}
}

// AtomicPredicate 原子谓词
type AtomicPredicate struct {
    Name string
    Terms []Term
}

func (ap *AtomicPredicate) Evaluate(interpretation *Interpretation, assignment map[string]interface{}) bool {
    args := make([]interface{}, len(ap.Terms))
    for i, term := range ap.Terms {
        args[i] = term.Evaluate(interpretation, assignment)
    }
    return interpretation.Predicates[ap.Name](args)
}

func (ap *AtomicPredicate) String() string {
    terms := make([]string, len(ap.Terms))
    for i, term := range ap.Terms {
        terms[i] = term.String()
    }
    return fmt.Sprintf("%s(%s)", ap.Name, strings.Join(terms, ", "))
}

// Term 项
type Term interface {
    Evaluate(interpretation *Interpretation, assignment map[string]interface{}) interface{}
    String() string
}

// Variable 变量
type Variable struct {
    Name string
}

func (v *Variable) Evaluate(interpretation *Interpretation, assignment map[string]interface{}) interface{} {
    return assignment[v.Name]
}

func (v *Variable) String() string {
    return v.Name
}

// Constant 常量
type Constant struct {
    Name string
}

func (c *Constant) Evaluate(interpretation *Interpretation, assignment map[string]interface{}) interface{} {
    return interpretation.Constants[c.Name]
}

func (c *Constant) String() string {
    return c.Name
}

// UniversalQuantifier 全称量词
type UniversalQuantifier struct {
    Variable string
    Formula  PredicateFormula
}

func (uq *UniversalQuantifier) Evaluate(interpretation *Interpretation, assignment map[string]interface{}) bool {
    for _, element := range interpretation.Domain {
        newAssignment := make(map[string]interface{})
        for k, v := range assignment {
            newAssignment[k] = v
        }
        newAssignment[uq.Variable] = element
        if !uq.Formula.Evaluate(interpretation, newAssignment) {
            return false
        }
    }
    return true
}

func (uq *UniversalQuantifier) String() string {
    return fmt.Sprintf("∀%s(%s)", uq.Variable, uq.Formula)
}

// ExistentialQuantifier 存在量词
type ExistentialQuantifier struct {
    Variable string
    Formula  PredicateFormula
}

func (eq *ExistentialQuantifier) Evaluate(interpretation *Interpretation, assignment map[string]interface{}) bool {
    for _, element := range interpretation.Domain {
        newAssignment := make(map[string]interface{})
        for k, v := range assignment {
            newAssignment[k] = v
        }
        newAssignment[eq.Variable] = element
        if eq.Formula.Evaluate(interpretation, newAssignment) {
            return true
        }
    }
    return false
}

func (eq *ExistentialQuantifier) String() string {
    return fmt.Sprintf("∃%s(%s)", eq.Variable, eq.Formula)
}
```

### 5.4 时态逻辑实现

```go
// TemporalFormula 表示时态逻辑公式
type TemporalFormula interface {
    Evaluate(model *KripkeModel, state int) bool
    String() string
}

// KripkeModel 克里普克模型
type KripkeModel struct {
    States []map[string]bool
    Transitions [][]int
}

// LTLFormula LTL公式
type LTLFormula interface {
    TemporalFormula
}

// G 全局算子
type G struct {
    Formula LTLFormula
}

func (g *G) Evaluate(model *KripkeModel, state int) bool {
    visited := make(map[int]bool)
    return g.evaluateRecursive(model, state, visited)
}

func (g *G) evaluateRecursive(model *KripkeModel, state int, visited map[int]bool) bool {
    if visited[state] {
        return true // 避免循环
    }
    visited[state] = true
    
    if !g.Formula.Evaluate(model, state) {
        return false
    }
    
    for _, nextState := range model.Transitions[state] {
        if !g.evaluateRecursive(model, nextState, visited) {
            return false
        }
    }
    return true
}

func (g *G) String() string {
    return fmt.Sprintf("G(%s)", g.Formula)
}

// F 未来算子
type F struct {
    Formula LTLFormula
}

func (f *F) Evaluate(model *KripkeModel, state int) bool {
    visited := make(map[int]bool)
    return f.evaluateRecursive(model, state, visited)
}

func (f *F) evaluateRecursive(model *KripkeModel, state int, visited map[int]bool) bool {
    if visited[state] {
        return false // 避免循环
    }
    visited[state] = true
    
    if f.Formula.Evaluate(model, state) {
        return true
    }
    
    for _, nextState := range model.Transitions[state] {
        if f.evaluateRecursive(model, nextState, visited) {
            return true
        }
    }
    return false
}

func (f *F) String() string {
    return fmt.Sprintf("F(%s)", f.Formula)
}

// X 下一个算子
type X struct {
    Formula LTLFormula
}

func (x *X) Evaluate(model *KripkeModel, state int) bool {
    if len(model.Transitions[state]) == 0 {
        return false
    }
    // 假设只有一个后继状态
    nextState := model.Transitions[state][0]
    return x.Formula.Evaluate(model, nextState)
}

func (x *X) String() string {
    return fmt.Sprintf("X(%s)", x.Formula)
}

// U 直到算子
type U struct {
    Left, Right LTLFormula
}

func (u *U) Evaluate(model *KripkeModel, state int) bool {
    visited := make(map[int]bool)
    return u.evaluateRecursive(model, state, visited)
}

func (u *U) evaluateRecursive(model *KripkeModel, state int, visited map[int]bool) bool {
    if visited[state] {
        return false // 避免循环
    }
    visited[state] = true
    
    if u.Right.Evaluate(model, state) {
        return true
    }
    
    if !u.Left.Evaluate(model, state) {
        return false
    }
    
    for _, nextState := range model.Transitions[state] {
        if u.evaluateRecursive(model, nextState, visited) {
            return true
        }
    }
    return false
}

func (u *U) String() string {
    return fmt.Sprintf("(%s U %s)", u.Left, u.Right)
}
```

### 5.5 使用示例

```go
func main() {
    // 命题逻辑示例
    fmt.Println("=== Propositional Logic ===")
    
    // 创建公式: (p ∧ q) → r
    p := &Atomic{Name: "p"}
    q := &Atomic{Name: "q"}
    r := &Atomic{Name: "r"}
    
    pAndQ := &And{Left: p, Right: q}
    formula := &Implies{Left: pAndQ, Right: r}
    
    variables := []string{"p", "q", "r"}
    TruthTable(formula, variables)
    
    fmt.Printf("Is tautology: %t\n", IsTautology(formula, variables))
    fmt.Printf("Is satisfiable: %t\n", IsSatisfiable(formula, variables))
    
    // 谓词逻辑示例
    fmt.Println("\n=== Predicate Logic ===")
    
    // 创建解释
    interpretation := &Interpretation{
        Domain: []interface{}{1, 2, 3, 4, 5},
        Predicates: map[string]func([]interface{}) bool{
            "Even": func(args []interface{}) bool {
                n := args[0].(int)
                return n%2 == 0
            },
            "Greater": func(args []interface{}) bool {
                a, b := args[0].(int), args[1].(int)
                return a > b
            },
        },
        Constants: map[string]interface{}{
            "zero": 0,
        },
    }
    
    // 创建公式: ∀x(Even(x) → x > 0)
    x := &Variable{Name: "x"}
    evenX := &AtomicPredicate{Name: "Even", Terms: []Term{x}}
    zero := &Constant{Name: "zero"}
    greaterXZero := &AtomicPredicate{Name: "Greater", Terms: []Term{x, zero}}
    
    implies := &Implies{Left: evenX, Right: greaterXZero}
    universalFormula := &UniversalQuantifier{Variable: "x", Formula: implies}
    
    assignment := make(map[string]interface{})
    result := universalFormula.Evaluate(interpretation, assignment)
    fmt.Printf("Formula: %s\n", universalFormula)
    fmt.Printf("Result: %t\n", result)
    
    // 时态逻辑示例
    fmt.Println("\n=== Temporal Logic ===")
    
    // 创建克里普克模型
    model := &KripkeModel{
        States: []map[string]bool{
            {"p": true, "q": false},
            {"p": false, "q": true},
            {"p": true, "q": true},
        },
        Transitions: [][]int{
            {1, 2},
            {0, 2},
            {0, 1},
        },
    }
    
    // 创建LTL公式: G(p → Fq)
    atomicP := &Atomic{Name: "p"}
    atomicQ := &Atomic{Name: "q"}
    
    fQ := &F{Formula: atomicQ}
    impliesPQ := &Implies{Left: atomicP, Right: fQ}
    gFormula := &G{Formula: impliesPQ}
    
    fmt.Printf("Formula: %s\n", gFormula)
    for i := 0; i < len(model.States); i++ {
        result := gFormula.Evaluate(model, i)
        fmt.Printf("State %d: %t\n", i, result)
    }
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

逻辑学是计算机科学的基础理论，在程序验证、人工智能、数据库理论等领域有广泛应用。通过形式化定义和Go语言实现，我们建立了从理论到实践的完整框架。

### 关键要点

1. **理论基础**: 命题逻辑、谓词逻辑、模态逻辑、时态逻辑
2. **推理系统**: 公理化系统、自然演绎、模型检验
3. **实现技术**: 真值表、解释、克里普克模型
4. **应用场景**: 程序验证、人工智能、数据库查询

### 进一步研究方向

1. **高阶逻辑**: 二阶逻辑、类型论
2. **非经典逻辑**: 直觉逻辑、模糊逻辑、多值逻辑
3. **逻辑编程**: Prolog、约束逻辑编程
4. **自动定理证明**: 归结、表方法、模型检验

---

**相关链接**:

- [01-数学基础](../01-Mathematical-Foundation/README.md)
- [03-范畴论基础](../03-Category-Theory-Foundation/README.md)
- [04-计算理论基础](../04-Computational-Theory-Foundation/README.md)
- [03-设计模式](../../03-Design-Patterns/README.md)
- [02-软件架构](../../02-Software-Architecture/README.md)
