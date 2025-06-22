# 02-语法分析 (Syntax Analysis)

## 目录

- [1. 概述](#1-概述)
- [2. 上下文无关文法](#2-上下文无关文法)
- [3. LL(1)分析](#3-ll1分析)
- [4. LR分析](#4-lr分析)
- [5. Go语言实现](#5-go语言实现)
- [6. 应用实例](#6-应用实例)

## 1. 概述

### 1.1 语法分析定义

语法分析是将词法单元序列转换为抽象语法树（AST）的过程。

**形式化定义**：

```latex
\text{语法分析器} = (G, T, N, P, S, \vdash)
```

其中：

- ```latex
$G$
``` 是文法
- ```latex
$T$
``` 是终结符集合
- ```latex
$N$
``` 是非终结符集合
- ```latex
$P$
``` 是产生式集合
- ```latex
$S$
``` 是开始符号
- ```latex
$\vdash$
``` 是推导关系

### 1.2 核心概念

#### 1.2.1 抽象语法树

```go
// ASTNode 抽象语法树节点
type ASTNode interface {
    Type() string
    Children() []ASTNode
    Value() string
    String() string
}

// ProgramNode 程序节点
type ProgramNode struct {
    Declarations []ASTNode
}

func (pn *ProgramNode) Type() string { return "Program" }
func (pn *ProgramNode) Children() []ASTNode { return pn.Declarations }
func (pn *ProgramNode) Value() string { return "" }
func (pn *ProgramNode) String() string { return "Program" }

// ExpressionNode 表达式节点
type ExpressionNode struct {
    Operator string
    Left, Right ASTNode
}

func (en *ExpressionNode) Type() string { return "Expression" }
func (en *ExpressionNode) Children() []ASTNode { 
    return []ASTNode{en.Left, en.Right} 
}
func (en *ExpressionNode) Value() string { return en.Operator }
func (en *ExpressionNode) String() string { 
    return fmt.Sprintf("(%s %s %s)", en.Left, en.Operator, en.Right) 
}
```

## 2. 上下文无关文法

### 2.1 文法定义

```latex
G = (V, T, P, S)
```

其中：

- ```latex
$V$
``` 是变量集合（非终结符）
- ```latex
$T$
``` 是终结符集合
- ```latex
$P$
``` 是产生式集合
- ```latex
$S$
``` 是开始符号

**Go语言实现**：

```go
// Grammar 上下文无关文法
type Grammar struct {
    Variables    map[string]bool  // 非终结符
    Terminals    map[string]bool  // 终结符
    Productions  []Production     // 产生式
    StartSymbol  string           // 开始符号
}

// Production 产生式
type Production struct {
    Left  string   // 左部（非终结符）
    Right []string // 右部（符号序列）
}

func (g *Grammar) AddProduction(left string, right ...string) {
    g.Productions = append(g.Productions, Production{
        Left: left,
        Right: right,
    })
    g.Variables[left] = true
    for _, symbol := range right {
        if !g.Variables[symbol] {
            g.Terminals[symbol] = true
        }
    }
}

// 示例：简单算术表达式文法
func CreateArithmeticGrammar() *Grammar {
    grammar := &Grammar{
        Variables: make(map[string]bool),
        Terminals: make(map[string]bool),
        StartSymbol: "E",
    }
    
    // E -> E + T | E - T | T
    grammar.AddProduction("E", "E", "+", "T")
    grammar.AddProduction("E", "E", "-", "T")
    grammar.AddProduction("E", "T")
    
    // T -> T * F | T / F | F
    grammar.AddProduction("T", "T", "*", "F")
    grammar.AddProduction("T", "T", "/", "F")
    grammar.AddProduction("T", "F")
    
    // F -> (E) | id | num
    grammar.AddProduction("F", "(", "E", ")")
    grammar.AddProduction("F", "id")
    grammar.AddProduction("F", "num")
    
    return grammar
}
```

### 2.2 推导和归约

```go
// Derivation 推导
type Derivation struct {
    Steps []string
}

func (g *Grammar) Derive(input []string) *Derivation {
    derivation := &Derivation{
        Steps: []string{g.StartSymbol},
    }
    
    current := g.StartSymbol
    
    for i := 0; i < len(input); i++ {
        // 找到可以应用的产生式
        for _, prod := range g.Productions {
            if strings.Contains(current, prod.Left) {
                // 应用产生式
                current = strings.Replace(current, prod.Left, 
                    strings.Join(prod.Right, " "), 1)
                derivation.Steps = append(derivation.Steps, current)
                break
            }
        }
    }
    
    return derivation
}

// Reduction 归约
type Reduction struct {
    Steps []string
}

func (g *Grammar) Reduce(input []string) *Reduction {
    reduction := &Reduction{
        Steps: []string{strings.Join(input, " ")},
    }
    
    current := input
    
    for len(current) > 1 || current[0] != g.StartSymbol {
        // 找到可以归约的产生式
        for _, prod := range g.Productions {
            if g.canReduce(current, prod) {
                // 执行归约
                current = g.applyReduction(current, prod)
                reduction.Steps = append(reduction.Steps, 
                    strings.Join(current, " "))
                break
            }
        }
    }
    
    return reduction
}

func (g *Grammar) canReduce(input []string, prod Production) bool {
    if len(input) < len(prod.Right) {
        return false
    }
    
    for i, symbol := range prod.Right {
        if input[i] != symbol {
            return false
        }
    }
    
    return true
}

func (g *Grammar) applyReduction(input []string, prod Production) []string {
    result := []string{prod.Left}
    result = append(result, input[len(prod.Right):]...)
    return result
}
```

## 3. LL(1)分析

### 3.1 LL(1)分析表

```go
// LL1Parser LL(1)分析器
type LL1Parser struct {
    Grammar *Grammar
    First   map[string]map[string]bool
    Follow  map[string]map[string]bool
    Table   map[string]map[string]Production
}

func NewLL1Parser(grammar *Grammar) *LL1Parser {
    parser := &LL1Parser{
        Grammar: grammar,
        First:   make(map[string]map[string]bool),
        Follow:  make(map[string]map[string]bool),
        Table:   make(map[string]map[string]Production),
    }
    
    parser.computeFirst()
    parser.computeFollow()
    parser.buildTable()
    
    return parser
}

func (p *LL1Parser) computeFirst() {
    // 初始化First集
    for variable := range p.Grammar.Variables {
        p.First[variable] = make(map[string]bool)
    }
    for terminal := range p.Grammar.Terminals {
        p.First[terminal] = map[string]bool{terminal: true}
    }
    
    // 迭代计算First集
    changed := true
    for changed {
        changed = false
        for _, prod := range p.Grammar.Productions {
            if p.addFirstSet(prod.Left, prod.Right) {
                changed = true
            }
        }
    }
}

func (p *LL1Parser) addFirstSet(left string, right []string) bool {
    changed := false
    
    if len(right) == 0 {
        // 空产生式
        if !p.First[left]["ε"] {
            p.First[left]["ε"] = true
            changed = true
        }
        return changed
    }
    
    // 计算右部符号序列的First集
    for i, symbol := range right {
        if p.Grammar.Terminals[symbol] {
            // 终结符
            if !p.First[left][symbol] {
                p.First[left][symbol] = true
                changed = true
            }
            break
        } else {
            // 非终结符
            for terminal := range p.First[symbol] {
                if terminal != "ε" {
                    if !p.First[left][terminal] {
                        p.First[left][terminal] = true
                        changed = true
                    }
                }
            }
            
            // 如果当前符号可以推导出ε，继续检查下一个符号
            if !p.First[symbol]["ε"] {
                break
            }
            
            // 如果所有符号都可以推导出ε
            if i == len(right)-1 {
                if !p.First[left]["ε"] {
                    p.First[left]["ε"] = true
                    changed = true
                }
            }
        }
    }
    
    return changed
}

func (p *LL1Parser) computeFollow() {
    // 初始化Follow集
    for variable := range p.Grammar.Variables {
        p.Follow[variable] = make(map[string]bool)
    }
    
    // 开始符号的Follow集包含$
    p.Follow[p.Grammar.StartSymbol]["$"] = true
    
    // 迭代计算Follow集
    changed := true
    for changed {
        changed = false
        for _, prod := range p.Grammar.Productions {
            if p.addFollowSet(prod.Left, prod.Right) {
                changed = true
            }
        }
    }
}

func (p *LL1Parser) addFollowSet(left string, right []string) bool {
    changed := false
    
    for i, symbol := range right {
        if p.Grammar.Variables[symbol] {
            // 计算非终结符的Follow集
            
            // 情况1：A -> αBβ，将First(β) - {ε} 加入Follow(B)
            if i < len(right)-1 {
                beta := right[i+1:]
                firstBeta := p.computeFirstOfSequence(beta)
                
                for terminal := range firstBeta {
                    if terminal != "ε" {
                        if !p.Follow[symbol][terminal] {
                            p.Follow[symbol][terminal] = true
                            changed = true
                        }
                    }
                }
                
                // 如果β可以推导出ε，将Follow(A)加入Follow(B)
                if firstBeta["ε"] {
                    for terminal := range p.Follow[left] {
                        if !p.Follow[symbol][terminal] {
                            p.Follow[symbol][terminal] = true
                            changed = true
                        }
                    }
                }
            } else {
                // 情况2：A -> αB，将Follow(A)加入Follow(B)
                for terminal := range p.Follow[left] {
                    if !p.Follow[symbol][terminal] {
                        p.Follow[symbol][terminal] = true
                        changed = true
                    }
                }
            }
        }
    }
    
    return changed
}

func (p *LL1Parser) computeFirstOfSequence(sequence []string) map[string]bool {
    result := make(map[string]bool)
    
    for i, symbol := range sequence {
        if p.Grammar.Terminals[symbol] {
            result[symbol] = true
            break
        } else {
            for terminal := range p.First[symbol] {
                if terminal != "ε" {
                    result[terminal] = true
                }
            }
            
            if !p.First[symbol]["ε"] {
                break
            }
            
            if i == len(sequence)-1 {
                result["ε"] = true
            }
        }
    }
    
    return result
}

func (p *LL1Parser) buildTable() {
    for _, prod := range p.Grammar.Productions {
        first := p.computeFirstOfSequence(prod.Right)
        
        for terminal := range first {
            if terminal != "ε" {
                if p.Table[prod.Left] == nil {
                    p.Table[prod.Left] = make(map[string]Production)
                }
                p.Table[prod.Left][terminal] = prod
            }
        }
        
        if first["ε"] {
            for terminal := range p.Follow[prod.Left] {
                if p.Table[prod.Left] == nil {
                    p.Table[prod.Left] = make(map[string]Production)
                }
                p.Table[prod.Left][terminal] = prod
            }
        }
    }
}
```

### 3.2 LL(1)分析算法

```go
func (p *LL1Parser) Parse(input []string) (ASTNode, error) {
    stack := []string{"$", p.Grammar.StartSymbol}
    input = append(input, "$")
    inputPos := 0
    
    var root ASTNode
    
    for len(stack) > 1 {
        top := stack[len(stack)-1]
        current := input[inputPos]
        
        if p.Grammar.Terminals[top] {
            // 栈顶是终结符
            if top == current {
                stack = stack[:len(stack)-1]
                inputPos++
            } else {
                return nil, fmt.Errorf("syntax error: expected %s, got %s", top, current)
            }
        } else {
            // 栈顶是非终结符
            if p.Table[top] == nil || p.Table[top][current].Left == "" {
                return nil, fmt.Errorf("syntax error: no production for %s with %s", top, current)
            }
            
            prod := p.Table[top][current]
            stack = stack[:len(stack)-1]
            
            // 将产生式右部逆序压入栈
            for i := len(prod.Right) - 1; i >= 0; i-- {
                if prod.Right[i] != "ε" {
                    stack = append(stack, prod.Right[i])
                }
            }
        }
    }
    
    return root, nil
}
```

## 4. LR分析

### 4.1 LR(0)项目

```go
// LR0Item LR(0)项目
type LR0Item struct {
    Production Production
    DotPosition int
}

func (item *LR0Item) String() string {
    right := make([]string, len(item.Production.Right))
    copy(right, item.Production.Right)
    
    if item.DotPosition <= len(right) {
        right = append(right[:item.DotPosition], append([]string{"·"}, right[item.DotPosition:]...)...)
    }
    
    return fmt.Sprintf("%s → %s", item.Production.Left, strings.Join(right, " "))
}

func (item *LR0Item) IsComplete() bool {
    return item.DotPosition >= len(item.Production.Right)
}

func (item *LR0Item) NextSymbol() string {
    if item.DotPosition < len(item.Production.Right) {
        return item.Production.Right[item.DotPosition]
    }
    return ""
}

func (item *LR0Item) Advance() *LR0Item {
    return &LR0Item{
        Production: item.Production,
        DotPosition: item.DotPosition + 1,
    }
}
```

### 4.2 LR(0)分析器

```go
// LR0Parser LR(0)分析器
type LR0Parser struct {
    Grammar *Grammar
    States  []LR0State
    Actions map[int]map[string]Action
    Gotos   map[int]map[string]int
}

// LR0State LR(0)状态
type LR0State struct {
    ID    int
    Items []*LR0Item
}

// Action 分析动作
type Action struct {
    Type   string // "shift", "reduce", "accept"
    Value  int    // 状态号或产生式编号
}

func NewLR0Parser(grammar *Grammar) *LR0Parser {
    parser := &LR0Parser{
        Grammar: grammar,
        Actions: make(map[int]map[string]Action),
        Gotos:   make(map[int]map[string]int),
    }
    
    parser.buildStates()
    parser.buildActions()
    
    return parser
}

func (p *LR0Parser) buildStates() {
    // 初始状态
    initialItems := []*LR0Item{}
    for _, prod := range p.Grammar.Productions {
        if prod.Left == p.Grammar.StartSymbol {
            initialItems = append(initialItems, &LR0Item{
                Production: prod,
                DotPosition: 0,
            })
        }
    }
    
    initialState := &LR0State{
        ID: 0,
        Items: p.closure(initialItems),
    }
    
    p.States = []*LR0State{initialState}
    
    // 构建所有状态
    for i := 0; i < len(p.States); i++ {
        state := p.States[i]
        p.addSuccessorStates(state)
    }
}

func (p *LR0Parser) closure(items []*LR0Item) []*LR0Item {
    closure := make([]*LR0Item, len(items))
    copy(closure, items)
    
    changed := true
    for changed {
        changed = false
        
        for _, item := range closure {
            nextSymbol := item.NextSymbol()
            if p.Grammar.Variables[nextSymbol] {
                // 添加所有以nextSymbol为左部的产生式的项目
                for _, prod := range p.Grammar.Productions {
                    if prod.Left == nextSymbol {
                        newItem := &LR0Item{
                            Production: prod,
                            DotPosition: 0,
                        }
                        
                        // 检查是否已存在
                        exists := false
                        for _, existingItem := range closure {
                            if existingItem.Production.Left == newItem.Production.Left &&
                               existingItem.DotPosition == newItem.DotPosition {
                                exists = true
                                break
                            }
                        }
                        
                        if !exists {
                            closure = append(closure, newItem)
                            changed = true
                        }
                    }
                }
            }
        }
    }
    
    return closure
}

func (p *LR0Parser) addSuccessorStates(state *LR0State) {
    // 按符号分组项目
    groups := make(map[string][]*LR0Item)
    
    for _, item := range state.Items {
        nextSymbol := item.NextSymbol()
        if nextSymbol != "" {
            groups[nextSymbol] = append(groups[nextSymbol], item.Advance())
        }
    }
    
    // 为每个符号创建后继状态
    for symbol, items := range groups {
        successorState := &LR0State{
            ID: len(p.States),
            Items: p.closure(items),
        }
        
        // 检查是否已存在相同状态
        existingStateID := p.findExistingState(successorState)
        if existingStateID == -1 {
            p.States = append(p.States, successorState)
            existingStateID = successorState.ID
        }
        
        // 记录转移
        if p.Grammar.Terminals[symbol] {
            if p.Actions[state.ID] == nil {
                p.Actions[state.ID] = make(map[string]Action)
            }
            p.Actions[state.ID][symbol] = Action{
                Type: "shift",
                Value: existingStateID,
            }
        } else {
            if p.Gotos[state.ID] == nil {
                p.Gotos[state.ID] = make(map[string]int)
            }
            p.Gotos[state.ID][symbol] = existingStateID
        }
    }
}

func (p *LR0Parser) findExistingState(state *LR0State) int {
    for _, existingState := range p.States {
        if p.statesEqual(state, existingState) {
            return existingState.ID
        }
    }
    return -1
}

func (p *LR0Parser) statesEqual(state1, state2 *LR0State) bool {
    if len(state1.Items) != len(state2.Items) {
        return false
    }
    
    for _, item1 := range state1.Items {
        found := false
        for _, item2 := range state2.Items {
            if item1.Production.Left == item2.Production.Left &&
               item1.DotPosition == item2.DotPosition {
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

func (p *LR0Parser) buildActions() {
    for _, state := range p.States {
        for _, item := range state.Items {
            if item.IsComplete() {
                // 归约动作
                if p.Actions[state.ID] == nil {
                    p.Actions[state.ID] = make(map[string]Action)
                }
                
                // 为所有终结符添加归约动作
                for terminal := range p.Grammar.Terminals {
                    p.Actions[state.ID][terminal] = Action{
                        Type: "reduce",
                        Value: p.findProductionIndex(item.Production),
                    }
                }
            }
        }
    }
}

func (p *LR0Parser) findProductionIndex(prod Production) int {
    for i, production := range p.Grammar.Productions {
        if production.Left == prod.Left && 
           len(production.Right) == len(prod.Right) {
            match := true
            for j, symbol := range prod.Right {
                if production.Right[j] != symbol {
                    match = false
                    break
                }
            }
            if match {
                return i
            }
        }
    }
    return -1
}
```

### 4.3 LR分析算法

```go
func (p *LR0Parser) Parse(input []string) (ASTNode, error) {
    stack := []int{0} // 状态栈
    symbols := []string{"$"} // 符号栈
    input = append(input, "$")
    inputPos := 0
    
    for {
        currentState := stack[len(stack)-1]
        currentSymbol := input[inputPos]
        
        action, exists := p.Actions[currentState][currentSymbol]
        if !exists {
            return nil, fmt.Errorf("syntax error: no action for state %d with symbol %s", 
                currentState, currentSymbol)
        }
        
        switch action.Type {
        case "shift":
            // 移进
            stack = append(stack, action.Value)
            symbols = append(symbols, currentSymbol)
            inputPos++
            
        case "reduce":
            // 归约
            prod := p.Grammar.Productions[action.Value]
            
            // 弹出符号栈
            for i := 0; i < len(prod.Right); i++ {
                stack = stack[:len(stack)-1]
                symbols = symbols[:len(symbols)-1]
            }
            
            // 压入归约后的非终结符
            symbols = append(symbols, prod.Left)
            
            // 查找GOTO表
            gotoState, exists := p.Gotos[stack[len(stack)-1]][prod.Left]
            if !exists {
                return nil, fmt.Errorf("syntax error: no goto for state %d with symbol %s", 
                    stack[len(stack)-1], prod.Left)
            }
            stack = append(stack, gotoState)
            
        case "accept":
            // 接受
            return p.buildAST(symbols), nil
            
        default:
            return nil, fmt.Errorf("unknown action type: %s", action.Type)
        }
    }
}

func (p *LR0Parser) buildAST(symbols []string) ASTNode {
    // 简化的AST构建
    // 实际实现需要更复杂的树构建逻辑
    return &ProgramNode{
        Declarations: []ASTNode{},
    }
}
```

## 5. Go语言实现

### 5.1 完整的语法分析器

```go
// Parser 语法分析器接口
type Parser interface {
    Parse(tokens []Token) (ASTNode, error)
}

// GoParser Go语言语法分析器
type GoParser struct {
    lexer *GoLexer
    grammar *Grammar
}

func NewGoParser() *GoParser {
    parser := &GoParser{
        lexer: NewGoLexer(""),
        grammar: CreateGoGrammar(),
    }
    return parser
}

func CreateGoGrammar() *Grammar {
    grammar := &Grammar{
        Variables: make(map[string]bool),
        Terminals: make(map[string]bool),
        StartSymbol: "Program",
    }
    
    // 简化的Go语言文法
    grammar.AddProduction("Program", "PackageDecl", "ImportDecls", "TopLevelDecls")
    grammar.AddProduction("PackageDecl", "package", "PackageName")
    grammar.AddProduction("ImportDecls", "ImportDecl")
    grammar.AddProduction("ImportDecls", "ImportDecls", "ImportDecl")
    grammar.AddProduction("ImportDecl", "import", "ImportSpec")
    grammar.AddProduction("TopLevelDecls", "TopLevelDecl")
    grammar.AddProduction("TopLevelDecls", "TopLevelDecls", "TopLevelDecl")
    grammar.AddProduction("TopLevelDecl", "FunctionDecl")
    grammar.AddProduction("FunctionDecl", "func", "FunctionName", "Signature", "Block")
    grammar.AddProduction("Signature", "Parameters", "Result")
    grammar.AddProduction("Parameters", "(", "ParameterList", ")")
    grammar.AddProduction("ParameterList", "Parameter")
    grammar.AddProduction("ParameterList", "ParameterList", ",", "Parameter")
    grammar.AddProduction("Parameter", "ParameterName", "Type")
    grammar.AddProduction("Result", "Type")
    grammar.AddProduction("Block", "{", "StatementList", "}")
    grammar.AddProduction("StatementList", "Statement")
    grammar.AddProduction("StatementList", "StatementList", "Statement")
    grammar.AddProduction("Statement", "ExpressionStatement")
    grammar.AddProduction("ExpressionStatement", "Expression", ";")
    grammar.AddProduction("Expression", "PrimaryExpr")
    grammar.AddProduction("Expression", "Expression", "+", "Expression")
    grammar.AddProduction("Expression", "Expression", "-", "Expression")
    grammar.AddProduction("Expression", "Expression", "*", "Expression")
    grammar.AddProduction("Expression", "Expression", "/", "Expression")
    grammar.AddProduction("PrimaryExpr", "Operand")
    grammar.AddProduction("Operand", "Literal")
    grammar.AddProduction("Operand", "OperandName")
    grammar.AddProduction("Literal", "BasicLit")
    grammar.AddProduction("BasicLit", "INT_LIT")
    grammar.AddProduction("BasicLit", "FLOAT_LIT")
    grammar.AddProduction("BasicLit", "STRING_LIT")
    grammar.AddProduction("OperandName", "IDENTIFIER")
    grammar.AddProduction("FunctionName", "IDENTIFIER")
    grammar.AddProduction("ParameterName", "IDENTIFIER")
    grammar.AddProduction("PackageName", "IDENTIFIER")
    grammar.AddProduction("Type", "TypeName")
    grammar.AddProduction("TypeName", "IDENTIFIER")
    grammar.AddProduction("ImportSpec", "STRING_LIT")
    
    return grammar
}

func (gp *GoParser) Parse(input string) (ASTNode, error) {
    // 词法分析
    gp.lexer.Input = input
    tokens := gp.lexer.Tokenize()
    
    // 语法分析
    return gp.parseProgram(tokens)
}

func (gp *GoParser) parseProgram(tokens []Token) (ASTNode, int, error) {
    // 简化的程序解析
    declarations := []ASTNode{}
    
    i := 0
    for i < len(tokens) {
        if tokens[i].Type == TokenKeyword && tokens[i].Value == "func" {
            funcDecl, newI, err := gp.parseFunctionDecl(tokens, i)
            if err != nil {
                return nil, i, err
            }
            declarations = append(declarations, funcDecl)
            i = newI
        } else {
            i++
        }
    }
    
    return &ProgramNode{Declarations: declarations}, i, nil
}

func (gp *GoParser) parseFunctionDecl(tokens []Token, start int) (ASTNode, int, error) {
    i := start
    
    // 解析 func
    if i >= len(tokens) || tokens[i].Value != "func" {
        return nil, i, fmt.Errorf("expected 'func'")
    }
    i++
    
    // 解析函数名
    if i >= len(tokens) || tokens[i].Type != TokenIdentifier {
        return nil, i, fmt.Errorf("expected function name")
    }
    funcName := tokens[i].Value
    i++
    
    // 解析参数列表
    if i >= len(tokens) || tokens[i].Value != "(" {
        return nil, i, fmt.Errorf("expected '('")
    }
    i++
    
    // 简化的参数解析
    for i < len(tokens) && tokens[i].Value != ")" {
        i++
    }
    
    if i >= len(tokens) || tokens[i].Value != ")" {
        return nil, i, fmt.Errorf("expected ')'")
    }
    i++
    
    // 解析函数体
    if i >= len(tokens) || tokens[i].Value != "{" {
        return nil, i, fmt.Errorf("expected '{'")
    }
    i++
    
    // 简化的函数体解析
    for i < len(tokens) && tokens[i].Value != "}" {
        i++
    }
    
    if i >= len(tokens) || tokens[i].Value != "}" {
        return nil, i, fmt.Errorf("expected '}'")
    }
    i++
    
    return &FunctionNode{
        Name: funcName,
        Body: &BlockNode{},
    }, i, nil
}
```

## 6. 应用实例

### 6.1 算术表达式解析

```go
func ExampleArithmeticParsing() {
    input := "2 + 3 * 4"
    
    // 词法分析
    lexer := NewGoLexer(input)
    tokens := lexer.Tokenize()
    
    // 语法分析
    grammar := CreateArithmeticGrammar()
    parser := NewLL1Parser(grammar)
    
    ast, err := parser.Parse([]string{"2", "+", "3", "*", "4"})
    if err != nil {
        fmt.Printf("Parse error: %v\n", err)
        return
    }
    
    fmt.Printf("AST: %s\n", ast)
}
```

### 6.2 Go代码解析

```go
func ExampleGoCodeParsing() {
    input := `
package main

import "fmt"

func main() {
    x := 42
    fmt.Println("Hello, World!")
}
`
    
    parser := NewGoParser()
    ast, err := parser.Parse(input)
    if err != nil {
        fmt.Printf("Parse error: %v\n", err)
        return
    }
    
    fmt.Printf("AST: %s\n", ast)
}
```

## 总结

语法分析是编译器的核心阶段，负责将词法单元序列转换为抽象语法树。通过上下文无关文法、LL(1)分析和LR分析等理论，我们可以构建高效准确的语法分析器。

**关键要点**：

1. **理论基础**：基于上下文无关文法和自动机理论
2. **分析方法**：自顶向下（LL）和自底向上（LR）分析
3. **实际应用**：在编译器、解释器、代码分析工具中广泛应用
4. **Go语言特性**：充分利用Go语言的接口和数据结构

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **语法分析理论完成！** 🚀
