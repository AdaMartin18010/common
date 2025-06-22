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

- 原子命题：```latex
p, q, r, \ldots
```
- 复合命题：由原子命题通过逻辑运算符构成
- 真值：```latex
\text{true}
``` 或 ```latex
\text{false}
```

### 1.2 逻辑运算符

**定义 1.2**: 基本逻辑运算符

1. **否定** (```latex
\neg
```): ```latex
\neg p
``` 表示"非 ```latex
p
```"
2. **合取** (```latex
\land
```): ```latex
p \land q
``` 表示"```latex
p
``` 且 ```latex
q
```"
3. **析取** (```latex
\lor
```): ```latex
p \lor q
``` 表示"```latex
p
``` 或 ```latex
q
```"
4. **蕴含** (```latex
\rightarrow
```): ```latex
p \rightarrow q
``` 表示"如果 ```latex
p
``` 则 ```latex
q
```"
5. **等价** (```latex
\leftrightarrow
```): ```latex
p \leftrightarrow q
``` 表示"```latex
p
``` 当且仅当 ```latex
q
```"

### 1.3 真值表

**定义 1.3**: 逻辑运算符的真值表

| ```latex
p
``` | ```latex
q
``` | ```latex
\neg p
``` | ```latex
p \land q
``` | ```latex
p \lor q
``` | ```latex
p \rightarrow q
``` | ```latex
p \leftrightarrow q
``` |
|-----|-----|----------|-------------|------------|-------------------|----------------------|
| T   | T   | F        | T           | T          | T                 | T                    |
| T   | F   | F        | F           | T          | F                 | F                    |
| F   | T   | T        | F           | T          | T                 | F                    |
| F   | F   | T        | F           | F          | T                 | T                    |

## 2. 谓词逻辑

### 2.1 量词

**定义 2.1**: 量词

1. **全称量词** (```latex
\forall
```): ```latex
\forall x P(x)
``` 表示"对所有 ```latex
x
```，```latex
P(x)
``` 成立"
2. **存在量词** (```latex
\exists
```): ```latex
\exists x P(x)
``` 表示"存在 ```latex
x
```，使得 ```latex
P(x)
``` 成立"

### 2.2 谓词公式

**定义 2.2**: 谓词公式的构成

- 原子公式：```latex
P(x), Q(x,y), \ldots
```
- 复合公式：由原子公式通过逻辑运算符和量词构成
- 自由变量：不被量词约束的变量
- 约束变量：被量词约束的变量

## 3. 形式化定义

### 3.1 逻辑公理

**公理 3.1** (同一律): ```latex
p \rightarrow p
```

**公理 3.2** (排中律): ```latex
p \lor \neg p
```

**公理 3.3** (矛盾律): ```latex
\neg(p \land \neg p)
```

**公理 3.4** (双重否定): ```latex
\neg\neg p \leftrightarrow p
```

### 3.2 推理规则

**规则 3.1** (假言推理): 从 ```latex
p \rightarrow q
``` 和 ```latex
p
``` 可以推出 ```latex
q
```

$```latex
\frac{p \rightarrow q \quad p}{q}
```$

**规则 3.2** (全称实例化): 从 ```latex
\forall x P(x)
``` 可以推出 ```latex
P(a)
```

$```latex
\frac{\forall x P(x)}{P(a)}
```$

**规则 3.3** (存在概括): 从 ```latex
P(a)
``` 可以推出 ```latex
\exists x P(x)
```

$```latex
\frac{P(a)}{\exists x P(x)}
```$

## 4. Go语言实现

### 4.1 命题逻辑实现

```go
// Proposition 表示一个命题
type Proposition interface {
    Evaluate(assignment map[string]bool) bool
    GetVariables() Set[string]
}

// AtomicProposition 原子命题
type AtomicProposition struct {
    Variable string
}

func (ap *AtomicProposition) Evaluate(assignment map[string]bool) bool {
    return assignment[ap.Variable]
}

func (ap *AtomicProposition) GetVariables() Set[string] {
    vars := NewSet[string]()
    vars.Add(ap.Variable)
    return vars
}

// Negation 否定
type Negation struct {
    Proposition Proposition
}

func (n *Negation) Evaluate(assignment map[string]bool) bool {
    return !n.Proposition.Evaluate(assignment)
}

func (n *Negation) GetVariables() Set[string] {
    return n.Proposition.GetVariables()
}

// Conjunction 合取
type Conjunction struct {
    Left  Proposition
    Right Proposition
}

func (c *Conjunction) Evaluate(assignment map[string]bool) bool {
    return c.Left.Evaluate(assignment) && c.Right.Evaluate(assignment)
}

func (c *Conjunction) GetVariables() Set[string] {
    vars := c.Left.GetVariables()
    rightVars := c.Right.GetVariables()
    return vars.Union(rightVars)
}

// Disjunction 析取
type Disjunction struct {
    Left  Proposition
    Right Proposition
}

func (d *Disjunction) Evaluate(assignment map[string]bool) bool {
    return d.Left.Evaluate(assignment) || d.Right.Evaluate(assignment)
}

func (d *Disjunction) GetVariables() Set[string] {
    vars := d.Left.GetVariables()
    rightVars := d.Right.GetVariables()
    return vars.Union(rightVars)
}

// Implication 蕴含
type Implication struct {
    Antecedent Proposition
    Consequent Proposition
}

func (i *Implication) Evaluate(assignment map[string]bool) bool {
    return !i.Antecedent.Evaluate(assignment) || i.Consequent.Evaluate(assignment)
}

func (i *Implication) GetVariables() Set[string] {
    vars := i.Antecedent.GetVariables()
    consequentVars := i.Consequent.GetVariables()
    return vars.Union(consequentVars)
}
```

### 4.2 谓词逻辑实现

```go
// Predicate 表示一个谓词
type Predicate interface {
    Evaluate(assignment map[string]interface{}) bool
    GetVariables() Set[string]
}

// AtomicPredicate 原子谓词
type AtomicPredicate struct {
    Name      string
    Arguments []string
    Function  func(args []interface{}) bool
}

func (ap *AtomicPredicate) Evaluate(assignment map[string]interface{}) bool {
    args := make([]interface{}, len(ap.Arguments))
    for i, arg := range ap.Arguments {
        args[i] = assignment[arg]
    }
    return ap.Function(args)
}

func (ap *AtomicPredicate) GetVariables() Set[string] {
    vars := NewSet[string]()
    for _, arg := range ap.Arguments {
        vars.Add(arg)
    }
    return vars
}

// UniversalQuantifier 全称量词
type UniversalQuantifier struct {
    Variable   string
    Predicate  Predicate
    Domain     []interface{}
}

func (uq *UniversalQuantifier) Evaluate(assignment map[string]interface{}) bool {
    for _, value := range uq.Domain {
        newAssignment := make(map[string]interface{})
        for k, v := range assignment {
            newAssignment[k] = v
        }
        newAssignment[uq.Variable] = value
        
        if !uq.Predicate.Evaluate(newAssignment) {
            return false
        }
    }
    return true
}

func (uq *UniversalQuantifier) GetVariables() Set[string] {
    vars := uq.Predicate.GetVariables()
    vars.Remove(uq.Variable) // 移除约束变量
    return vars
}

// ExistentialQuantifier 存在量词
type ExistentialQuantifier struct {
    Variable   string
    Predicate  Predicate
    Domain     []interface{}
}

func (eq *ExistentialQuantifier) Evaluate(assignment map[string]interface{}) bool {
    for _, value := range eq.Domain {
        newAssignment := make(map[string]interface{})
        for k, v := range assignment {
            newAssignment[k] = v
        }
        newAssignment[eq.Variable] = value
        
        if eq.Predicate.Evaluate(newAssignment) {
            return true
        }
    }
    return false
}

func (eq *ExistentialQuantifier) GetVariables() Set[string] {
    vars := eq.Predicate.GetVariables()
    vars.Remove(eq.Variable) // 移除约束变量
    return vars
}
```

### 4.3 逻辑推理引擎

```go
// LogicEngine 逻辑推理引擎
type LogicEngine struct {
    KnowledgeBase []Proposition
    Rules         []InferenceRule
}

// InferenceRule 推理规则
type InferenceRule struct {
    Name     string
    Premises []Proposition
    Conclusion Proposition
}

// NewLogicEngine 创建逻辑推理引擎
func NewLogicEngine() *LogicEngine {
    return &LogicEngine{
        KnowledgeBase: make([]Proposition, 0),
        Rules:         make([]InferenceRule, 0),
    }
}

// AddKnowledge 添加知识
func (le *LogicEngine) AddKnowledge(proposition Proposition) {
    le.KnowledgeBase = append(le.KnowledgeBase, proposition)
}

// AddRule 添加推理规则
func (le *LogicEngine) AddRule(rule InferenceRule) {
    le.Rules = append(le.Rules, rule)
}

// ForwardChaining 前向推理
func (le *LogicEngine) ForwardChaining() []Proposition {
    inferred := NewSet[Proposition]()
    agenda := make([]Proposition, len(le.KnowledgeBase))
    copy(agenda, le.KnowledgeBase)
    
    for len(agenda) > 0 {
        current := agenda[0]
        agenda = agenda[1:]
        
        if !inferred.Contains(current) {
            inferred.Add(current)
            
            // 应用推理规则
            for _, rule := range le.Rules {
                if le.canApplyRule(rule, inferred) {
                    conclusion := rule.Conclusion
                    if !inferred.Contains(conclusion) {
                        agenda = append(agenda, conclusion)
                    }
                }
            }
        }
    }
    
    return inferred.ToSlice()
}

// canApplyRule 检查是否可以应用规则
func (le *LogicEngine) canApplyRule(rule InferenceRule, knowledge Set[Proposition]) bool {
    for _, premise := range rule.Premises {
        if !knowledge.Contains(premise) {
            return false
        }
    }
    return true
}

// BackwardChaining 后向推理
func (le *LogicEngine) BackwardChaining(goal Proposition) bool {
    return le.backwardChainingHelper(goal, NewSet[Proposition]())
}

// backwardChainingHelper 后向推理辅助函数
func (le *LogicEngine) backwardChainingHelper(goal Proposition, visited Set[Proposition]) bool {
    if visited.Contains(goal) {
        return false // 避免循环
    }
    
    visited.Add(goal)
    
    // 检查目标是否在知识库中
    for _, knowledge := range le.KnowledgeBase {
        if knowledge == goal {
            return true
        }
    }
    
    // 尝试通过规则证明目标
    for _, rule := range le.Rules {
        if rule.Conclusion == goal {
            allPremisesProven := true
            for _, premise := range rule.Premises {
                if !le.backwardChainingHelper(premise, visited) {
                    allPremisesProven = false
                    break
                }
            }
            if allPremisesProven {
                return true
            }
        }
    }
    
    return false
}
```

## 5. 应用示例

### 5.1 规则引擎

```go
// BusinessRule 业务规则
type BusinessRule struct {
    Name        string
    Condition   Proposition
    Action      func() error
}

// RuleEngine 规则引擎
type RuleEngine struct {
    rules []BusinessRule
}

// NewRuleEngine 创建规则引擎
func NewRuleEngine() *RuleEngine {
    return &RuleEngine{
        rules: make([]BusinessRule, 0),
    }
}

// AddRule 添加规则
func (re *RuleEngine) AddRule(rule BusinessRule) {
    re.rules = append(re.rules, rule)
}

// Execute 执行规则
func (re *RuleEngine) Execute(context map[string]bool) error {
    for _, rule := range re.rules {
        if rule.Condition.Evaluate(context) {
            if err := rule.Action(); err != nil {
                return fmt.Errorf("rule %s failed: %w", rule.Name, err)
            }
        }
    }
    return nil
}

// 示例：贷款审批规则
func createLoanApprovalRules() *RuleEngine {
    engine := NewRuleEngine()
    
    // 规则1：收入大于50000且信用分数大于700
    income := &AtomicProposition{Variable: "income_gt_50k"}
    credit := &AtomicProposition{Variable: "credit_score_gt_700"}
    condition1 := &Conjunction{Left: income, Right: credit}
    
    rule1 := BusinessRule{
        Name:      "High Income High Credit",
        Condition: condition1,
        Action: func() error {
            fmt.Println("Approved: High income and high credit score")
            return nil
        },
    }
    
    engine.AddRule(rule1)
    
    return engine
}
```

### 5.2 形式化验证

```go
// SystemState 系统状态
type SystemState struct {
    Variables map[string]interface{}
}

// Invariant 不变式
type Invariant struct {
    Name      string
    Predicate Predicate
}

// SystemVerifier 系统验证器
type SystemVerifier struct {
    invariants []Invariant
    states     []SystemState
}

// NewSystemVerifier 创建系统验证器
func NewSystemVerifier() *SystemVerifier {
    return &SystemVerifier{
        invariants: make([]Invariant, 0),
        states:     make([]SystemState, 0),
    }
}

// AddInvariant 添加不变式
func (sv *SystemVerifier) AddInvariant(invariant Invariant) {
    sv.invariants = append(sv.invariants, invariant)
}

// AddState 添加状态
func (sv *SystemVerifier) AddState(state SystemState) {
    sv.states = append(sv.states, state)
}

// Verify 验证系统
func (sv *SystemVerifier) Verify() []string {
    violations := make([]string, 0)
    
    for _, state := range sv.states {
        for _, invariant := range sv.invariants {
            if !invariant.Predicate.Evaluate(state.Variables) {
                violation := fmt.Sprintf("Invariant '%s' violated in state", invariant.Name)
                violations = append(violations, violation)
            }
        }
    }
    
    return violations
}

// 示例：银行账户验证
func createBankAccountVerifier() *SystemVerifier {
    verifier := NewSystemVerifier()
    
    // 不变式：余额不能为负
    balancePredicate := &AtomicPredicate{
        Name:      "balance_non_negative",
        Arguments: []string{"balance"},
        Function: func(args []interface{}) bool {
            balance := args[0].(float64)
            return balance >= 0
        },
    }
    
    invariant := Invariant{
        Name:      "Non-negative Balance",
        Predicate: balancePredicate,
    }
    
    verifier.AddInvariant(invariant)
    
    // 添加测试状态
    state1 := SystemState{
        Variables: map[string]interface{}{
            "balance": 100.0,
        },
    }
    
    state2 := SystemState{
        Variables: map[string]interface{}{
            "balance": -50.0,
        },
    }
    
    verifier.AddState(state1)
    verifier.AddState(state2)
    
    return verifier
}
```

### 5.3 专家系统

```go
// ExpertSystem 专家系统
type ExpertSystem struct {
    facts     Set[string]
    rules     []ExpertRule
    engine    *LogicEngine
}

// ExpertRule 专家规则
type ExpertRule struct {
    Name     string
    If       []string
    Then     string
    Confidence float64
}

// NewExpertSystem 创建专家系统
func NewExpertSystem() *ExpertSystem {
    return &ExpertSystem{
        facts:  NewSet[string](),
        rules:  make([]ExpertRule, 0),
        engine: NewLogicEngine(),
    }
}

// AddFact 添加事实
func (es *ExpertSystem) AddFact(fact string) {
    es.facts.Add(fact)
}

// AddRule 添加规则
func (es *ExpertSystem) AddRule(rule ExpertRule) {
    es.rules = append(es.rules, rule)
}

// Infer 推理
func (es *ExpertSystem) Infer() Set[string] {
    conclusions := NewSet[string]()
    
    for _, rule := range es.rules {
        if es.canApplyRule(rule) {
            conclusions.Add(rule.Then)
        }
    }
    
    return conclusions
}

// canApplyRule 检查是否可以应用规则
func (es *ExpertSystem) canApplyRule(rule ExpertRule) bool {
    for _, condition := range rule.If {
        if !es.facts.Contains(condition) {
            return false
        }
    }
    return true
}

// 示例：医疗诊断专家系统
func createMedicalDiagnosisSystem() *ExpertSystem {
    system := NewExpertSystem()
    
    // 添加症状事实
    system.AddFact("fever")
    system.AddFact("cough")
    system.AddFact("fatigue")
    
    // 添加诊断规则
    rule1 := ExpertRule{
        Name: "Flu Diagnosis",
        If:   []string{"fever", "cough", "fatigue"},
        Then: "influenza",
        Confidence: 0.8,
    }
    
    rule2 := ExpertRule{
        Name: "Cold Diagnosis",
        If:   []string{"cough", "fatigue"},
        Then: "common_cold",
        Confidence: 0.6,
    }
    
    system.AddRule(rule1)
    system.AddRule(rule2)
    
    return system
}
```

## 总结

逻辑学为计算机科学提供了形式化推理的基础，通过Go语言的实现，我们可以构建强大的逻辑推理系统。
这些系统在规则引擎、形式化验证、专家系统等领域有广泛应用。

**关键特性**:

- 完整的命题逻辑和谓词逻辑实现
- 支持前向推理和后向推理
- 类型安全的泛型设计
- 实际应用场景的示例

**应用领域**:

- 人工智能和专家系统
- 形式化验证和模型检查
- 业务规则引擎
- 自然语言处理
- 数据库查询优化
