# 02-逻辑学 (Logic)

## 目录

- [02-逻辑学 (Logic)](#02-逻辑学-logic)
  - [目录](#目录)
  - [1. 命题逻辑](#1-命题逻辑)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 逻辑运算符](#12-逻辑运算符)
    - [1.3 真值表](#13-真值表)
  - [2. 谓词逻辑](#2-谓词逻辑)
    - [2.1 量词](#21-量词)
    - [2.2 谓词公式](#22-谓词公式)
  - [3. 形式化定义](#3-形式化定义)
    - [3.1 逻辑公理](#31-逻辑公理)
    - [3.2 推理规则](#32-推理规则)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 命题逻辑实现](#41-命题逻辑实现)
    - [4.2 谓词逻辑实现](#42-谓词逻辑实现)
    - [4.3 逻辑推理引擎](#43-逻辑推理引擎)
  - [5. 应用示例](#5-应用示例)
    - [5.1 规则引擎](#51-规则引擎)
    - [5.2 形式化验证](#52-形式化验证)
    - [5.3 专家系统](#53-专家系统)
  - [总结](#总结)

## 1. 命题逻辑

### 1.1 基本概念

**定义 1.1**: 命题是一个可以判断真假的陈述句。

**形式化表达**:

- 原子命题：$p, q, r, \ldots$
- 复合命题：由原子命题通过逻辑运算符构成
- 真值：$\text{true}$ 或 $\text{false}$

### 1.2 逻辑运算符

**定义 1.2**: 基本逻辑运算符

1. **否定** ($\neg$): $\neg p$ 表示"非 $p$"
2. **合取** ($\land$): $p \land q$ 表示"$p$ 且 $q$"
3. **析取** ($\lor$): $p \lor q$ 表示"$p$ 或 $q$"
4. **蕴含** ($\rightarrow$): $p \rightarrow q$ 表示"如果 $p$ 则 $q$"
5. **等价** ($\leftrightarrow$): $p \leftrightarrow q$ 表示"$p$ 当且仅当 $q$"

### 1.3 真值表

**定义 1.3**: 逻辑运算符的真值表

| $p$ | $q$ | $\neg p$ | $p \land q$ | $p \lor q$ | $p \rightarrow q$ | $p \leftrightarrow q$ |
|-----|-----|----------|-------------|------------|-------------------|----------------------|
| T   | T   | F        | T           | T          | T                 | T                    |
| T   | F   | F        | F           | T          | F                 | F                    |
| F   | T   | T        | F           | T          | T                 | F                    |
| F   | F   | T        | F           | F          | T                 | T                    |

## 2. 谓词逻辑

### 2.1 量词

**定义 2.1**: 量词

1. **全称量词** ($\forall$): $\forall x P(x)$ 表示"对所有 $x$，$P(x)$ 成立"
2. **存在量词** ($\exists$): $\exists x P(x)$ 表示"存在 $x$，使得 $P(x)$ 成立"

### 2.2 谓词公式

**定义 2.2**: 谓词公式的构成

- 原子公式：$P(x), Q(x,y), \ldots$
- 复合公式：由原子公式通过逻辑运算符和量词构成
- 自由变量：不被量词约束的变量
- 约束变量：被量词约束的变量

## 3. 形式化定义

### 3.1 逻辑公理

**公理 3.1** (同一律): $p \rightarrow p$

**公理 3.2** (排中律): $p \lor \neg p$

**公理 3.3** (矛盾律): $\neg(p \land \neg p)$

**公理 3.4** (双重否定): $\neg\neg p \leftrightarrow p$

### 3.2 推理规则

**规则 3.1** (假言推理): 从 $p \rightarrow q$ 和 $p$ 可以推出 $q$

$$\frac{p \rightarrow q \quad p}{q}$$

**规则 3.2** (全称实例化): 从 $\forall x P(x)$ 可以推出 $P(a)$

$$\frac{\forall x P(x)}{P(a)}$$

**规则 3.3** (存在概括): 从 $P(a)$ 可以推出 $\exists x P(x)$

$$\frac{P(a)}{\exists x P(x)}$$

## 4. Go语言实现

### 4.1 命题逻辑实现

```go
// 命题类型
type Proposition interface {
    Evaluate(assignment map[string]bool) bool
    String() string
}

// 原子命题
type Atom struct {
    Name string
}

func (a *Atom) Evaluate(assignment map[string]bool) bool {
    return assignment[a.Name]
}

// 否定命题
type Not struct {
    Operand Proposition
}

func (n *Not) Evaluate(assignment map[string]bool) bool {
    return !n.Operand.Evaluate(assignment)
}

// 合取命题
type And struct {
    Left, Right Proposition
}

func (a *And) Evaluate(assignment map[string]bool) bool {
    return a.Left.Evaluate(assignment) && a.Right.Evaluate(assignment)
}

// 析取命题
type Or struct {
    Left, Right Proposition
}

func (o *Or) Evaluate(assignment map[string]bool) bool {
    return o.Left.Evaluate(assignment) || o.Right.Evaluate(assignment)
}
```

### 4.2 谓词逻辑实现

```go
// 谓词类型
type Predicate interface {
    Evaluate(domain []interface{}, assignment map[string]interface{}) bool
    String() string
}

// 原子谓词
type AtomicPredicate struct {
    Name     string
    Arity    int
    Function func([]interface{}) bool
}

func (ap *AtomicPredicate) Evaluate(args []interface{}) bool {
    return ap.Function(args)
}

// 全称量词
type ForAll struct {
    Variable  string
    Domain   []interface{}
    Formula   Predicate
}

func (fa *ForAll) Evaluate(domain []interface{}, assignment map[string]interface{}) bool {
    for _, value := range domain {
        newAssignment := copyAssignment(assignment)
        newAssignment[fa.Variable] = value
        if !fa.Formula.Evaluate(domain, newAssignment) {
            return false
        }
    }
    return true
}
```

### 4.3 逻辑推理引擎

```go
// 推理引擎
type InferenceEngine struct {
    KnowledgeBase []Proposition
    Rules         []InferenceRule
}

// 推理规则
type InferenceRule struct {
    Premises    []Proposition
    Conclusion  Proposition
}

// 应用推理规则
func (ie *InferenceEngine) ApplyRule(rule InferenceRule, facts map[string]bool) bool {
    // 检查所有前提是否满足
    for _, premise := range rule.Premises {
        if !premise.Evaluate(facts) {
            return false
        }
    }
    // 添加结论到知识库
    ie.KnowledgeBase = append(ie.KnowledgeBase, rule.Conclusion)
    return true
}

// 前向链接推理
func (ie *InferenceEngine) ForwardChaining(facts map[string]bool) []Proposition {
    var conclusions []Proposition
    changed := true
    
    for changed {
        changed = false
        for _, rule := range ie.Rules {
            if ie.ApplyRule(rule, facts) {
                changed = true
                conclusions = append(conclusions, rule.Conclusion)
            }
        }
    }
    
    return conclusions
}
```

## 5. 应用示例

### 5.1 规则引擎

```go
// 业务规则引擎示例
type BusinessRuleEngine struct {
    InferenceEngine
    BusinessRules map[string]Rule
}

func (bre *BusinessRuleEngine) AddBusinessRule(name string, rule Rule) {
    bre.BusinessRules[name] = rule
    // 将业务规则转换为逻辑命题
    proposition := rule.ToProposition()
    bre.KnowledgeBase = append(bre.KnowledgeBase, proposition)
}

func (bre *BusinessRuleEngine) EvaluateRules(context map[string]interface{}) []string {
    var results []string
    facts := bre.extractFacts(context)
    
    conclusions := bre.ForwardChaining(facts)
    for _, conclusion := range conclusions {
        results = append(results, conclusion.String())
    }
    
    return results
}
```

### 5.2 形式化验证

```go
// 形式化验证器
type FormalVerifier struct {
    Specification Proposition
    Implementation Proposition
}

func (fv *FormalVerifier) Verify() bool {
    // 验证实现是否满足规范
    implication := &Implies{
        Left:  fv.Implementation,
        Right: fv.Specification,
    }
    
    // 使用SAT求解器检查蕴含式是否为永真式
    return fv.checkTautology(implication)
}

func (fv *FormalVerifier) checkTautology(prop Proposition) bool {
    // 使用真值表法检查永真式
    assignments := generateAllAssignments(prop.Variables())
    for _, assignment := range assignments {
        if !prop.Evaluate(assignment) {
            return false
        }
    }
    return true
}
```

### 5.3 专家系统

```go
// 专家系统
type ExpertSystem struct {
    InferenceEngine
    Facts        map[string]bool
    Conclusions  []string
}

func (es *ExpertSystem) AddFact(fact string, value bool) {
    es.Facts[fact] = value
}

func (es *ExpertSystem) Infer() []string {
    conclusions := es.ForwardChaining(es.Facts)
    for _, conclusion := range conclusions {
        es.Conclusions = append(es.Conclusions, conclusion.String())
    }
    return es.Conclusions
}

func (es *ExpertSystem) Explain(conclusion string) string {
    // 生成推理链的解释
    return es.generateExplanation(conclusion)
}
```

## 总结

逻辑学为计算机科学提供了形式化推理的基础，通过命题逻辑和谓词逻辑，我们可以：

1. 形式化描述问题和解决方案
2. 验证程序的正确性
3. 构建智能推理系统
4. 实现自动化决策系统

在实际应用中，逻辑学的原理被广泛应用于：

- 规则引擎
- 形式化验证
- 专家系统
- 自动定理证明
- 人工智能推理

通过Go语言的实现，我们可以将这些理论概念转化为实用的工程解决方案。 