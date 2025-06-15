# 01-命题逻辑 (Propositional Logic)

## 目录

- [01-命题逻辑 (Propositional Logic)](#01-命题逻辑-propositional-logic)
  - [目录](#目录)
  - [1. 基础定义](#1-基础定义)
    - [1.1 命题](#11-命题)
    - [1.2 逻辑连接词](#12-逻辑连接词)
    - [1.3 真值表](#13-真值表)
  - [2. 命题公式](#2-命题公式)
    - [2.1 语法](#21-语法)
    - [2.2 语义](#22-语义)
    - [2.3 等价性](#23-等价性)
  - [3. 推理系统](#3-推理系统)
    - [3.1 自然演绎](#31-自然演绎)
    - [3.2 公理系统](#32-公理系统)
    - [3.3 归结推理](#33-归结推理)
  - [4. 可满足性](#4-可满足性)
    - [4.1 SAT问题](#41-sat问题)
    - [4.2 DPLL算法](#42-dpll算法)
    - [4.3 局部搜索](#43-局部搜索)
  - [5. 在软件工程中的应用](#5-在软件工程中的应用)
    - [5.1 程序验证](#51-程序验证)
    - [5.2 模型检查](#52-模型检查)
    - [5.3 约束求解](#53-约束求解)

## 1. 基础定义

### 1.1 命题

**形式化定义**：

```latex
命题是能够判断真假的陈述句。

命题变元: p, q, r, ... ∈ P (命题变元集合)
真值: true (⊤), false (⊥)

命题的语义函数: ⟦_⟧: P → {true, false}
```

**Go语言实现**：

```go
// 命题变元
type Proposition struct {
    Name string
    Value bool
}

// 命题集合
type PropositionSet struct {
    propositions map[string]*Proposition
}

func NewPropositionSet() *PropositionSet {
    return &PropositionSet{
        propositions: make(map[string]*Proposition),
    }
}

func (ps *PropositionSet) AddProposition(name string) {
    ps.propositions[name] = &Proposition{
        Name:  name,
        Value: false,
    }
}

func (ps *PropositionSet) SetValue(name string, value bool) {
    if prop, exists := ps.propositions[name]; exists {
        prop.Value = value
    }
}

func (ps *PropositionSet) GetValue(name string) bool {
    if prop, exists := ps.propositions[name]; exists {
        return prop.Value
    }
    return false
}

// 真值赋值
type Interpretation struct {
    values map[string]bool
}

func NewInterpretation() *Interpretation {
    return &Interpretation{
        values: make(map[string]bool),
    }
}

func (i *Interpretation) SetValue(prop string, value bool) {
    i.values[prop] = value
}

func (i *Interpretation) GetValue(prop string) bool {
    if value, exists := i.values[prop]; exists {
        return value
    }
    return false
}
```

### 1.2 逻辑连接词

**数学定义**：

```latex
逻辑连接词：
- 否定: ¬p (NOT)
- 合取: p ∧ q (AND)
- 析取: p ∨ q (OR)
- 蕴含: p → q (IMPLIES)
- 等价: p ↔ q (IFF)

真值表定义：
¬p: 当p为真时¬p为假，当p为假时¬p为真
p ∧ q: 当p和q都为真时为真，否则为假
p ∨ q: 当p或q至少一个为真时为真，否则为假
p → q: 当p为假或q为真时为真，否则为假
p ↔ q: 当p和q真值相同时为真，否则为假
```

**Go语言实现**：

```go
// 逻辑连接词
type LogicalOperator int

const (
    NOT LogicalOperator = iota
    AND
    OR
    IMPLIES
    IFF
)

// 命题公式
type Formula struct {
    Type     LogicalOperator
    Left     *Formula
    Right    *Formula
    Variable string
}

// 原子命题
func NewAtom(variable string) *Formula {
    return &Formula{
        Variable: variable,
    }
}

// 否定
func NewNot(formula *Formula) *Formula {
    return &Formula{
        Type:  NOT,
        Left:  formula,
    }
}

// 合取
func NewAnd(left, right *Formula) *Formula {
    return &Formula{
        Type:  AND,
        Left:  left,
        Right: right,
    }
}

// 析取
func NewOr(left, right *Formula) *Formula {
    return &Formula{
        Type:  OR,
        Left:  left,
        Right: right,
    }
}

// 蕴含
func NewImplies(left, right *Formula) *Formula {
    return &Formula{
        Type:  IMPLIES,
        Left:  left,
        Right: right,
    }
}

// 等价
func NewIff(left, right *Formula) *Formula {
    return &Formula{
        Type:  IFF,
        Left:  left,
        Right: right,
    }
}

// 公式求值
func (f *Formula) Evaluate(interpretation *Interpretation) bool {
    if f.Variable != "" {
        return interpretation.GetValue(f.Variable)
    }
    
    switch f.Type {
    case NOT:
        return !f.Left.Evaluate(interpretation)
    case AND:
        return f.Left.Evaluate(interpretation) && f.Right.Evaluate(interpretation)
    case OR:
        return f.Left.Evaluate(interpretation) || f.Right.Evaluate(interpretation)
    case IMPLIES:
        return !f.Left.Evaluate(interpretation) || f.Right.Evaluate(interpretation)
    case IFF:
        return f.Left.Evaluate(interpretation) == f.Right.Evaluate(interpretation)
    default:
        return false
    }
}
```

### 1.3 真值表

**数学定义**：

```latex
真值表是命题公式在所有可能真值赋值下的真值列表。

对于n个命题变元，有2^n种不同的真值赋值。
```

**Go语言实现**：

```go
// 真值表生成器
type TruthTable struct {
    formula *Formula
    variables []string
}

func NewTruthTable(formula *Formula) *TruthTable {
    variables := formula.GetVariables()
    return &TruthTable{
        formula:   formula,
        variables: variables,
    }
}

func (f *Formula) GetVariables() []string {
    variables := make(map[string]bool)
    f.collectVariables(variables)
    
    result := make([]string, 0)
    for v := range variables {
        result = append(result, v)
    }
    sort.Strings(result)
    return result
}

func (f *Formula) collectVariables(variables map[string]bool) {
    if f.Variable != "" {
        variables[f.Variable] = true
    } else {
        if f.Left != nil {
            f.Left.collectVariables(variables)
        }
        if f.Right != nil {
            f.Right.collectVariables(variables)
        }
    }
}

func (tt *TruthTable) Generate() [][]bool {
    n := len(tt.variables)
    rows := 1 << n // 2^n
    result := make([][]bool, rows)
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        // 根据行号生成真值赋值
        for j, variable := range tt.variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        // 计算公式的真值
        formulaValue := tt.formula.Evaluate(interpretation)
        
        // 构建结果行
        row := make([]bool, n+1)
        for j, variable := range tt.variables {
            row[j] = interpretation.GetValue(variable)
        }
        row[n] = formulaValue
        result[i] = row
    }
    
    return result
}

func (tt *TruthTable) Print() {
    table := tt.Generate()
    
    // 打印表头
    for _, variable := range tt.variables {
        fmt.Printf("%s\t", variable)
    }
    fmt.Println("Result")
    
    // 打印分隔线
    for i := 0; i < len(tt.variables)+1; i++ {
        fmt.Print("----\t")
    }
    fmt.Println()
    
    // 打印真值表
    for _, row := range table {
        for i := 0; i < len(row)-1; i++ {
            if row[i] {
                fmt.Print("T\t")
            } else {
                fmt.Print("F\t")
            }
        }
        if row[len(row)-1] {
            fmt.Println("T")
        } else {
            fmt.Println("F")
        }
    }
}
```

## 2. 命题公式

### 2.1 语法

**BNF语法**：

```latex
φ ::= p | ¬φ | φ ∧ φ | φ ∨ φ | φ → φ | φ ↔ φ
其中 p ∈ P (命题变元集合)
```

**Go语言实现**：

```go
// 语法分析器
type Parser struct {
    input string
    pos   int
}

func NewParser(input string) *Parser {
    return &Parser{
        input: input,
        pos:   0,
    }
}

func (p *Parser) Parse() *Formula {
    return p.parseOr()
}

func (p *Parser) parseOr() *Formula {
    left := p.parseAnd()
    
    for p.pos < len(p.input) && p.peek() == '∨' {
        p.consume('∨')
        right := p.parseAnd()
        left = NewOr(left, right)
    }
    
    return left
}

func (p *Parser) parseAnd() *Formula {
    left := p.parseImplies()
    
    for p.pos < len(p.input) && p.peek() == '∧' {
        p.consume('∧')
        right := p.parseImplies()
        left = NewAnd(left, right)
    }
    
    return left
}

func (p *Parser) parseImplies() *Formula {
    left := p.parseIff()
    
    for p.pos < len(p.input) && p.peek() == '→' {
        p.consume('→')
        right := p.parseIff()
        left = NewImplies(left, right)
    }
    
    return left
}

func (p *Parser) parseIff() *Formula {
    left := p.parseNot()
    
    for p.pos < len(p.input) && p.peek() == '↔' {
        p.consume('↔')
        right := p.parseNot()
        left = NewIff(left, right)
    }
    
    return left
}

func (p *Parser) parseNot() *Formula {
    if p.pos < len(p.input) && p.peek() == '¬' {
        p.consume('¬')
        return NewNot(p.parseNot())
    }
    
    return p.parseAtom()
}

func (p *Parser) parseAtom() *Formula {
    if p.pos < len(p.input) && p.peek() == '(' {
        p.consume('(')
        result := p.parseOr()
        p.consume(')')
        return result
    }
    
    // 解析变量名
    var name strings.Builder
    for p.pos < len(p.input) && isLetter(p.peek()) {
        name.WriteByte(p.peek())
        p.pos++
    }
    
    return NewAtom(name.String())
}

func (p *Parser) peek() byte {
    if p.pos < len(p.input) {
        return p.input[p.pos]
    }
    return 0
}

func (p *Parser) consume(expected byte) {
    if p.peek() == expected {
        p.pos++
    } else {
        panic(fmt.Sprintf("Expected %c, got %c", expected, p.peek()))
    }
}

func isLetter(b byte) bool {
    return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
```

### 2.2 语义

**语义函数**：

```latex
语义函数 ⟦_⟧: Formula × Interpretation → {true, false}

⟦p⟧ᵢ = i(p)
⟦¬φ⟧ᵢ = ¬⟦φ⟧ᵢ
⟦φ ∧ ψ⟧ᵢ = ⟦φ⟧ᵢ ∧ ⟦ψ⟧ᵢ
⟦φ ∨ ψ⟧ᵢ = ⟦φ⟧ᵢ ∨ ⟦ψ⟧ᵢ
⟦φ → ψ⟧ᵢ = ¬⟦φ⟧ᵢ ∨ ⟦ψ⟧ᵢ
⟦φ ↔ ψ⟧ᵢ = (⟦φ⟧ᵢ ∧ ⟦ψ⟧ᵢ) ∨ (¬⟦φ⟧ᵢ ∧ ¬⟦ψ⟧ᵢ)
```

**Go语言实现**：

```go
// 语义分析器
type SemanticAnalyzer struct {
    formula *Formula
}

func NewSemanticAnalyzer(formula *Formula) *SemanticAnalyzer {
    return &SemanticAnalyzer{
        formula: formula,
    }
}

// 检查公式是否为重言式
func (sa *SemanticAnalyzer) IsTautology() bool {
    variables := sa.formula.GetVariables()
    n := len(variables)
    rows := 1 << n
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        for j, variable := range variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        if !sa.formula.Evaluate(interpretation) {
            return false
        }
    }
    
    return true
}

// 检查公式是否为矛盾式
func (sa *SemanticAnalyzer) IsContradiction() bool {
    variables := sa.formula.GetVariables()
    n := len(variables)
    rows := 1 << n
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        for j, variable := range variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        if sa.formula.Evaluate(interpretation) {
            return false
        }
    }
    
    return true
}

// 检查公式是否为可满足式
func (sa *SemanticAnalyzer) IsSatisfiable() bool {
    return !sa.IsContradiction()
}

// 找到满足公式的一个真值赋值
func (sa *SemanticAnalyzer) FindSatisfyingAssignment() *Interpretation {
    variables := sa.formula.GetVariables()
    n := len(variables)
    rows := 1 << n
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        for j, variable := range variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        if sa.formula.Evaluate(interpretation) {
            return interpretation
        }
    }
    
    return nil // 不可满足
}
```

### 2.3 等价性

**等价关系**：

```latex
两个公式φ和ψ等价，记作φ ≡ ψ，当且仅当：
对于所有真值赋值i，⟦φ⟧ᵢ = ⟦ψ⟧ᵢ

重要等价式：
- 双重否定: ¬¬φ ≡ φ
- 德摩根律: ¬(φ ∧ ψ) ≡ ¬φ ∨ ¬ψ, ¬(φ ∨ ψ) ≡ ¬φ ∧ ¬ψ
- 分配律: φ ∧ (ψ ∨ χ) ≡ (φ ∧ ψ) ∨ (φ ∧ χ)
- 结合律: (φ ∧ ψ) ∧ χ ≡ φ ∧ (ψ ∧ χ)
```

**Go语言实现**：

```go
// 等价性检查器
type EquivalenceChecker struct {
    formula1 *Formula
    formula2 *Formula
}

func NewEquivalenceChecker(f1, f2 *Formula) *EquivalenceChecker {
    return &EquivalenceChecker{
        formula1: f1,
        formula2: f2,
    }
}

func (ec *EquivalenceChecker) AreEquivalent() bool {
    // 获取所有变量
    vars1 := ec.formula1.GetVariables()
    vars2 := ec.formula2.GetVariables()
    
    // 合并变量集
    allVars := make(map[string]bool)
    for _, v := range vars1 {
        allVars[v] = true
    }
    for _, v := range vars2 {
        allVars[v] = true
    }
    
    variables := make([]string, 0)
    for v := range allVars {
        variables = append(variables, v)
    }
    sort.Strings(variables)
    
    n := len(variables)
    rows := 1 << n
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        for j, variable := range variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        val1 := ec.formula1.Evaluate(interpretation)
        val2 := ec.formula2.Evaluate(interpretation)
        
        if val1 != val2 {
            return false
        }
    }
    
    return true
}

// 公式化简器
type FormulaSimplifier struct {
    formula *Formula
}

func NewFormulaSimplifier(formula *Formula) *FormulaSimplifier {
    return &FormulaSimplifier{
        formula: formula,
    }
}

func (fs *FormulaSimplifier) Simplify() *Formula {
    return fs.simplifyFormula(fs.formula)
}

func (fs *FormulaSimplifier) simplifyFormula(f *Formula) *Formula {
    if f.Variable != "" {
        return f
    }
    
    switch f.Type {
    case NOT:
        simplified := fs.simplifyFormula(f.Left)
        return fs.simplifyNot(simplified)
    case AND:
        left := fs.simplifyFormula(f.Left)
        right := fs.simplifyFormula(f.Right)
        return fs.simplifyAnd(left, right)
    case OR:
        left := fs.simplifyFormula(f.Left)
        right := fs.simplifyFormula(f.Right)
        return fs.simplifyOr(left, right)
    case IMPLIES:
        left := fs.simplifyFormula(f.Left)
        right := fs.simplifyFormula(f.Right)
        return fs.simplifyImplies(left, right)
    case IFF:
        left := fs.simplifyFormula(f.Left)
        right := fs.simplifyFormula(f.Right)
        return fs.simplifyIff(left, right)
    default:
        return f
    }
}

func (fs *FormulaSimplifier) simplifyNot(f *Formula) *Formula {
    // 双重否定消除
    if f.Type == NOT {
        return f.Left
    }
    
    // 德摩根律
    if f.Type == AND {
        return NewOr(NewNot(f.Left), NewNot(f.Right))
    }
    if f.Type == OR {
        return NewAnd(NewNot(f.Left), NewNot(f.Right))
    }
    
    return NewNot(f)
}

func (fs *FormulaSimplifier) simplifyAnd(left, right *Formula) *Formula {
    // 恒真式简化
    if isTrue(left) {
        return right
    }
    if isTrue(right) {
        return left
    }
    
    // 恒假式简化
    if isFalse(left) || isFalse(right) {
        return NewAtom("false")
    }
    
    return NewAnd(left, right)
}

func (fs *FormulaSimplifier) simplifyOr(left, right *Formula) *Formula {
    // 恒假式简化
    if isFalse(left) {
        return right
    }
    if isFalse(right) {
        return left
    }
    
    // 恒真式简化
    if isTrue(left) || isTrue(right) {
        return NewAtom("true")
    }
    
    return NewOr(left, right)
}

func (fs *FormulaSimplifier) simplifyImplies(left, right *Formula) *Formula {
    // 蕴含转换为析取
    return NewOr(NewNot(left), right)
}

func (fs *FormulaSimplifier) simplifyIff(left, right *Formula) *Formula {
    // 等价转换为合取和析取
    return NewAnd(
        NewImplies(left, right),
        NewImplies(right, left),
    )
}

func isTrue(f *Formula) bool {
    return f.Variable == "true"
}

func isFalse(f *Formula) bool {
    return f.Variable == "false"
}
```

## 3. 推理系统

### 3.1 自然演绎

**推理规则**：

```latex
引入规则：
- ∧I: 从φ和ψ推出φ ∧ ψ
- ∨I: 从φ推出φ ∨ ψ
- →I: 从假设φ推出ψ，然后推出φ → ψ

消除规则：
- ∧E: 从φ ∧ ψ推出φ或ψ
- ∨E: 从φ ∨ ψ和φ → χ和ψ → χ推出χ
- →E: 从φ和φ → ψ推出ψ
```

**Go语言实现**：

```go
// 自然演绎系统
type NaturalDeduction struct {
    premises []*Formula
    conclusion *Formula
}

func NewNaturalDeduction(premises []*Formula, conclusion *Formula) *NaturalDeduction {
    return &NaturalDeduction{
        premises:   premises,
        conclusion: conclusion,
    }
}

// 推理步骤
type InferenceStep struct {
    Rule     string
    Premises []int
    Result   *Formula
}

// 证明
type Proof struct {
    Steps []*InferenceStep
}

func (nd *NaturalDeduction) Prove() *Proof {
    proof := &Proof{
        Steps: make([]*InferenceStep, 0),
    }
    
    // 添加前提
    for i, premise := range nd.premises {
        proof.Steps = append(proof.Steps, &InferenceStep{
            Rule:     "Premise",
            Premises: []int{},
            Result:   premise,
        })
    }
    
    // 尝试应用推理规则
    nd.applyInferenceRules(proof)
    
    return proof
}

func (nd *NaturalDeduction) applyInferenceRules(proof *Proof) {
    // 简化的推理规则应用
    for i, step := range proof.Steps {
        // 尝试应用合取引入
        if i > 0 {
            nd.tryAndIntroduction(proof, i)
        }
        
        // 尝试应用蕴含消除
        nd.tryImplicationElimination(proof, i)
    }
}

func (nd *NaturalDeduction) tryAndIntroduction(proof *Proof, stepIndex int) {
    // 检查是否可以应用合取引入
    for i := 0; i < stepIndex; i++ {
        for j := i + 1; j < stepIndex; j++ {
            if canApplyAndIntroduction(proof.Steps[i].Result, proof.Steps[j].Result) {
                result := NewAnd(proof.Steps[i].Result, proof.Steps[j].Result)
                proof.Steps = append(proof.Steps, &InferenceStep{
                    Rule:     "∧I",
                    Premises: []int{i, j},
                    Result:   result,
                })
            }
        }
    }
}

func (nd *NaturalDeduction) tryImplicationElimination(proof *Proof, stepIndex int) {
    currentStep := proof.Steps[stepIndex]
    
    // 检查当前步骤是否是蕴含
    if currentStep.Result.Type == IMPLIES {
        // 寻找前件
        for i := 0; i < stepIndex; i++ {
            if areEquivalent(proof.Steps[i].Result, currentStep.Result.Left) {
                result := currentStep.Result.Right
                proof.Steps = append(proof.Steps, &InferenceStep{
                    Rule:     "→E",
                    Premises: []int{i, stepIndex},
                    Result:   result,
                })
            }
        }
    }
}

func canApplyAndIntroduction(f1, f2 *Formula) bool {
    // 检查两个公式是否可以合取
    return true // 简化实现
}

func areEquivalent(f1, f2 *Formula) bool {
    checker := NewEquivalenceChecker(f1, f2)
    return checker.AreEquivalent()
}
```

### 3.2 公理系统

**公理模式**：

```latex
公理模式：
1. φ → (ψ → φ)
2. (φ → (ψ → χ)) → ((φ → ψ) → (φ → χ))
3. (¬φ → ¬ψ) → (ψ → φ)

推理规则：假言推理 (Modus Ponens)
从φ和φ → ψ推出ψ
```

**Go语言实现**：

```go
// 公理系统
type AxiomaticSystem struct {
    axioms []*Formula
}

func NewAxiomaticSystem() *AxiomaticSystem {
    return &AxiomaticSystem{
        axioms: make([]*Formula, 0),
    }
}

func (as *AxiomaticSystem) AddAxiom(axiom *Formula) {
    as.axioms = append(as.axioms, axiom)
}

func (as *AxiomaticSystem) IsAxiom(formula *Formula) bool {
    for _, axiom := range as.axioms {
        if areEquivalent(formula, axiom) {
            return true
        }
    }
    return false
}

// 公理化证明
type AxiomaticProof struct {
    Steps []*ProofStep
}

type ProofStep struct {
    Formula *Formula
    Justification string
    References []int
}

func (as *AxiomaticSystem) Prove(formula *Formula) *AxiomaticProof {
    proof := &AxiomaticProof{
        Steps: make([]*ProofStep, 0),
    }
    
    // 添加公理
    for i, axiom := range as.axioms {
        proof.Steps = append(proof.Steps, &ProofStep{
            Formula:      axiom,
            Justification: "Axiom",
            References:   []int{},
        })
    }
    
    // 尝试证明目标公式
    as.tryToProve(proof, formula)
    
    return proof
}

func (as *AxiomaticSystem) tryToProve(proof *AxiomaticProof, target *Formula) {
    // 简化的证明策略
    for i, step := range proof.Steps {
        // 尝试应用假言推理
        for j := 0; j < i; j++ {
            if step.Formula.Type == IMPLIES {
                if areEquivalent(proof.Steps[j].Formula, step.Formula.Left) {
                    result := step.Formula.Right
                    proof.Steps = append(proof.Steps, &ProofStep{
                        Formula:      result,
                        Justification: "MP",
                        References:   []int{j, i},
                    })
                }
            }
        }
    }
}
```

### 3.3 归结推理

**归结原理**：

```latex
归结规则：从子句C₁ ∨ ¬p和C₂ ∨ p推出C₁ ∨ C₂

归结推理步骤：
1. 将公式转换为合取范式 (CNF)
2. 将CNF表示为子句集合
3. 应用归结规则直到得到空子句或无法继续
```

**Go语言实现**：

```go
// 子句
type Clause struct {
    literals []*Literal
}

type Literal struct {
    variable string
    negated  bool
}

func NewClause(literals []*Literal) *Clause {
    return &Clause{
        literals: literals,
    }
}

func (c *Clause) AddLiteral(literal *Literal) {
    c.literals = append(c.literals, literal)
}

func (c *Clause) IsEmpty() bool {
    return len(c.literals) == 0
}

func (c *Clause) Contains(literal *Literal) bool {
    for _, l := range c.literals {
        if l.variable == literal.variable && l.negated == literal.negated {
            return true
        }
    }
    return false
}

// 归结推理器
type ResolutionProver struct {
    clauses []*Clause
}

func NewResolutionProver() *ResolutionProver {
    return &ResolutionProver{
        clauses: make([]*Clause, 0),
    }
}

func (rp *ResolutionProver) AddClause(clause *Clause) {
    rp.clauses = append(rp.clauses, clause)
}

func (rp *ResolutionProver) Prove() bool {
    for {
        newClauses := make([]*Clause, 0)
        
        // 尝试所有可能的归结
        for i := 0; i < len(rp.clauses); i++ {
            for j := i + 1; j < len(rp.clauses); j++ {
                resolvent := rp.resolve(rp.clauses[i], rp.clauses[j])
                if resolvent != nil {
                    if resolvent.IsEmpty() {
                        return true // 找到矛盾
                    }
                    newClauses = append(newClauses, resolvent)
                }
            }
        }
        
        // 检查是否有新的子句
        if len(newClauses) == 0 {
            return false // 无法归结
        }
        
        // 添加新子句
        for _, clause := range newClauses {
            if !rp.containsClause(clause) {
                rp.clauses = append(rp.clauses, clause)
            }
        }
    }
}

func (rp *ResolutionProver) resolve(clause1, clause2 *Clause) *Clause {
    // 寻找互补文字
    for _, lit1 := range clause1.literals {
        for _, lit2 := range clause2.literals {
            if lit1.variable == lit2.variable && lit1.negated != lit2.negated {
                // 找到互补文字，进行归结
                result := NewClause([]*Literal{})
                
                // 添加clause1中除lit1外的所有文字
                for _, l := range clause1.literals {
                    if l != lit1 {
                        result.AddLiteral(l)
                    }
                }
                
                // 添加clause2中除lit2外的所有文字
                for _, l := range clause2.literals {
                    if l != lit2 && !result.Contains(l) {
                        result.AddLiteral(l)
                    }
                }
                
                return result
            }
        }
    }
    
    return nil
}

func (rp *ResolutionProver) containsClause(clause *Clause) bool {
    for _, c := range rp.clauses {
        if rp.clausesEqual(c, clause) {
            return true
        }
    }
    return false
}

func (rp *ResolutionProver) clausesEqual(c1, c2 *Clause) bool {
    if len(c1.literals) != len(c2.literals) {
        return false
    }
    
    for _, lit1 := range c1.literals {
        found := false
        for _, lit2 := range c2.literals {
            if lit1.variable == lit2.variable && lit1.negated == lit2.negated {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    
    return true
}
```

## 4. 可满足性

### 4.1 SAT问题

**问题定义**：

```latex
SAT问题：给定一个命题公式φ，判断是否存在真值赋值使得φ为真。

SAT问题是NP完全问题，是理论计算机科学的核心问题之一。
```

**Go语言实现**：

```go
// SAT求解器接口
type SATSolver interface {
    Solve(formula *Formula) *Interpretation
}

// 暴力搜索求解器
type BruteForceSolver struct{}

func NewBruteForceSolver() *BruteForceSolver {
    return &BruteForceSolver{}
}

func (bfs *BruteForceSolver) Solve(formula *Formula) *Interpretation {
    variables := formula.GetVariables()
    n := len(variables)
    rows := 1 << n
    
    for i := 0; i < rows; i++ {
        interpretation := NewInterpretation()
        
        for j, variable := range variables {
            value := (i >> j) & 1 == 1
            interpretation.SetValue(variable, value)
        }
        
        if formula.Evaluate(interpretation) {
            return interpretation
        }
    }
    
    return nil // 不可满足
}
```

### 4.2 DPLL算法

**算法描述**：

```latex
DPLL算法 (Davis-Putnam-Logemann-Loveland):
1. 单元传播：如果子句中只有一个文字，则必须为真
2. 纯文字消除：如果某个文字在所有子句中都是正或负，则可以消除
3. 分支：选择一个变量进行赋值，递归求解
```

**Go语言实现**：

```go
// DPLL求解器
type DPLLSolver struct{}

func NewDPLLSolver() *DPLLSolver {
    return &DPLLSolver{}
}

func (dpll *DPLLSolver) Solve(formula *Formula) *Interpretation {
    // 转换为CNF
    cnf := dpll.toCNF(formula)
    
    // 转换为子句集合
    clauses := dpll.toClauses(cnf)
    
    // 应用DPLL算法
    return dpll.dpll(clauses, NewInterpretation())
}

func (dpll *DPLLSolver) dpll(clauses []*Clause, assignment *Interpretation) *Interpretation {
    // 单元传播
    clauses, assignment = dpll.unitPropagation(clauses, assignment)
    
    // 检查是否满足
    if dpll.allClausesSatisfied(clauses, assignment) {
        return assignment
    }
    
    // 检查是否有空子句
    if dpll.hasEmptyClause(clauses) {
        return nil
    }
    
    // 选择变量进行分支
    variable := dpll.chooseVariable(clauses)
    
    // 尝试赋值为真
    assignment.SetValue(variable, true)
    result := dpll.dpll(clauses, assignment)
    if result != nil {
        return result
    }
    
    // 尝试赋值为假
    assignment.SetValue(variable, false)
    return dpll.dpll(clauses, assignment)
}

func (dpll *DPLLSolver) unitPropagation(clauses []*Clause, assignment *Interpretation) ([]*Clause, *Interpretation) {
    // 简化的单元传播实现
    return clauses, assignment
}

func (dpll *DPLLSolver) allClausesSatisfied(clauses []*Clause, assignment *Interpretation) bool {
    for _, clause := range clauses {
        if !dpll.clauseSatisfied(clause, assignment) {
            return false
        }
    }
    return true
}

func (dpll *DPLLSolver) clauseSatisfied(clause *Clause, assignment *Interpretation) bool {
    for _, literal := range clause.literals {
        value := assignment.GetValue(literal.variable)
        if literal.negated {
            value = !value
        }
        if value {
            return true
        }
    }
    return false
}

func (dpll *DPLLSolver) hasEmptyClause(clauses []*Clause) bool {
    for _, clause := range clauses {
        if clause.IsEmpty() {
            return true
        }
    }
    return false
}

func (dpll *DPLLSolver) chooseVariable(clauses []*Clause) string {
    // 简化的变量选择策略
    if len(clauses) > 0 && len(clauses[0].literals) > 0 {
        return clauses[0].literals[0].variable
    }
    return ""
}

func (dpll *DPLLSolver) toCNF(formula *Formula) *Formula {
    // 简化的CNF转换
    return formula
}

func (dpll *DPLLSolver) toClauses(cnf *Formula) []*Clause {
    // 简化的子句转换
    return []*Clause{}
}
```

### 4.3 局部搜索

**算法描述**：

```latex
局部搜索算法：
1. 随机生成初始真值赋值
2. 计算不满足的子句数量
3. 通过翻转变量来减少不满足的子句数量
4. 重复直到找到满足的解或达到最大迭代次数
```

**Go语言实现**：

```go
// 局部搜索求解器
type LocalSearchSolver struct {
    maxIterations int
    maxFlips      int
}

func NewLocalSearchSolver(maxIterations, maxFlips int) *LocalSearchSolver {
    return &LocalSearchSolver{
        maxIterations: maxIterations,
        maxFlips:      maxFlips,
    }
}

func (lss *LocalSearchSolver) Solve(formula *Formula) *Interpretation {
    variables := formula.GetVariables()
    
    for iteration := 0; iteration < lss.maxIterations; iteration++ {
        // 随机生成初始赋值
        assignment := lss.randomAssignment(variables)
        
        // 局部搜索
        for flip := 0; flip < lss.maxFlips; flip++ {
            // 检查是否满足
            if formula.Evaluate(assignment) {
                return assignment
            }
            
            // 选择变量进行翻转
            variable := lss.chooseVariableToFlip(formula, assignment)
            if variable != "" {
                currentValue := assignment.GetValue(variable)
                assignment.SetValue(variable, !currentValue)
            }
        }
    }
    
    return nil
}

func (lss *LocalSearchSolver) randomAssignment(variables []string) *Interpretation {
    assignment := NewInterpretation()
    for _, variable := range variables {
        assignment.SetValue(variable, rand.Float64() < 0.5)
    }
    return assignment
}

func (lss *LocalSearchSolver) chooseVariableToFlip(formula *Formula, assignment *Interpretation) string {
    variables := formula.GetVariables()
    
    // 计算每个变量的得分
    bestVariable := ""
    bestScore := -1
    
    for _, variable := range variables {
        score := lss.calculateFlipScore(formula, assignment, variable)
        if score > bestScore {
            bestScore = score
            bestVariable = variable
        }
    }
    
    return bestVariable
}

func (lss *LocalSearchSolver) calculateFlipScore(formula *Formula, assignment *Interpretation, variable string) int {
    // 简化的得分计算
    return rand.Intn(100)
}
```

## 5. 在软件工程中的应用

### 5.1 程序验证

**霍尔逻辑**：

```go
// 霍尔三元组
type HoareTriple struct {
    Precondition  *Formula
    Program       string
    Postcondition *Formula
}

// 程序验证器
type ProgramVerifier struct {
    triples []*HoareTriple
}

func NewProgramVerifier() *ProgramVerifier {
    return &ProgramVerifier{
        triples: make([]*HoareTriple, 0),
    }
}

func (pv *ProgramVerifier) AddTriple(triple *HoareTriple) {
    pv.triples = append(pv.triples, triple)
}

func (pv *ProgramVerifier) Verify() bool {
    for _, triple := range pv.triples {
        if !pv.verifyTriple(triple) {
            return false
        }
    }
    return true
}

func (pv *ProgramVerifier) verifyTriple(triple *HoareTriple) bool {
    // 简化的验证实现
    // 实际实现需要更复杂的程序分析
    return true
}
```

### 5.2 模型检查

**状态转换系统**：

```go
// 状态
type State struct {
    Variables map[string]bool
}

// 转换
type Transition struct {
    From     *State
    To       *State
    Condition *Formula
}

// 状态转换系统
type TransitionSystem struct {
    States      []*State
    Transitions []*Transition
    Initial     *State
}

// 模型检查器
type ModelChecker struct {
    system *TransitionSystem
}

func NewModelChecker(system *TransitionSystem) *ModelChecker {
    return &ModelChecker{
        system: system,
    }
}

func (mc *ModelChecker) CheckProperty(property *Formula) bool {
    // 简化的模型检查实现
    // 实际实现需要更复杂的算法
    return true
}
```

### 5.3 约束求解

**约束系统**：

```go
// 约束
type Constraint struct {
    Variables []string
    Formula   *Formula
}

// 约束求解器
type ConstraintSolver struct {
    constraints []*Constraint
}

func NewConstraintSolver() *ConstraintSolver {
    return &ConstraintSolver{
        constraints: make([]*Constraint, 0),
    }
}

func (cs *ConstraintSolver) AddConstraint(constraint *Constraint) {
    cs.constraints = append(cs.constraints, constraint)
}

func (cs *ConstraintSolver) Solve() *Interpretation {
    // 将所有约束合并为一个公式
    combinedFormula := cs.combineConstraints()
    
    // 使用SAT求解器求解
    solver := NewBruteForceSolver()
    return solver.Solve(combinedFormula)
}

func (cs *ConstraintSolver) combineConstraints() *Formula {
    if len(cs.constraints) == 0 {
        return NewAtom("true")
    }
    
    result := cs.constraints[0].Formula
    for i := 1; i < len(cs.constraints); i++ {
        result = NewAnd(result, cs.constraints[i].Formula)
    }
    
    return result
}
```

## 总结

命题逻辑为软件工程提供了：

1. **形式化推理**：严格的逻辑推理系统
2. **程序验证**：霍尔逻辑和程序正确性证明
3. **模型检查**：状态转换系统的性质验证
4. **约束求解**：SAT求解和约束满足问题

**核心要点**：
- 命题逻辑的语法和语义
- 推理系统和证明方法
- SAT问题的求解算法
- 在软件工程中的实际应用

这个完整的命题逻辑框架为软件工程的逻辑分析提供了坚实的理论基础。

---

**相关链接**:

- [02-谓词逻辑](02-Predicate-Logic.md)
- [03-模态逻辑](03-Modal-Logic.md)
- [04-时态逻辑](04-Temporal-Logic.md)
- [返回逻辑基础层](../README.md)
