# 03-语义分析 (Semantic Analysis)

## 目录

- [03-语义分析 (Semantic Analysis)](#03-语义分析-semantic-analysis)
  - [目录](#目录)
  - [1. 语义分析基础](#1-语义分析基础)
    - [1.1 语义分析定义](#11-语义分析定义)
    - [1.2 语义检查类型](#12-语义检查类型)
    - [1.3 符号表管理](#13-符号表管理)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 语义域](#21-语义域)
    - [2.2 语义函数](#22-语义函数)
    - [2.3 类型系统](#23-类型系统)
  - [3. 类型检查](#3-类型检查)
    - [3.1 类型推导](#31-类型推导)
    - [3.2 类型统一](#32-类型统一)
    - [3.3 类型安全](#33-类型安全)
  - [4. 作用域分析](#4-作用域分析)
    - [4.1 静态作用域](#41-静态作用域)
    - [4.2 动态作用域](#42-动态作用域)
    - [4.3 闭包处理](#43-闭包处理)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 符号表](#51-符号表)
    - [5.2 类型检查器](#52-类型检查器)
    - [5.3 语义分析器](#53-语义分析器)
  - [6. 应用示例](#6-应用示例)
    - [6.1 简单语言](#61-简单语言)
    - [6.2 函数式语言](#62-函数式语言)
    - [6.3 面向对象语言](#63-面向对象语言)
  - [7. 数学证明](#7-数学证明)
    - [7.1 类型安全定理](#71-类型安全定理)
    - [7.2 类型推导定理](#72-类型推导定理)
    - [7.3 语义一致性定理](#73-语义一致性定理)
  - [总结](#总结)

---

## 1. 语义分析基础

### 1.1 语义分析定义

语义分析是编译过程中的一个重要阶段，它检查程序的语义正确性，包括类型检查、作用域分析、语义约束验证等。

**定义 1.1** (语义分析器): 语义分析器是一个四元组 ```latex
$(\mathcal{S}, \mathcal{T}, \mathcal{E}, \mathcal{C})$
```，其中：

- ```latex
$\mathcal{S}$
``` 是符号表
- ```latex
$\mathcal{T}$
``` 是类型系统
- ```latex
$\mathcal{E}$
``` 是语义环境
- ```latex
$\mathcal{C}$
``` 是语义检查函数

**定义 1.2** (语义正确性): 程序 ```latex
$P$
``` 在语义环境 ```latex
$\mathcal{E}$
``` 下是语义正确的，当且仅当：

$```latex
$\forall \sigma \in \mathcal{E}, \mathcal{C}(P, \sigma) = \text{true}$
```$

### 1.2 语义检查类型

**类型检查**: 验证表达式和语句的类型正确性
**作用域分析**: 检查变量和函数的可见性
**语义约束**: 验证语言特定的语义规则
**控制流分析**: 检查程序的控制流正确性

### 1.3 符号表管理

**定义 1.3** (符号表): 符号表是一个映射 ```latex
$\mathcal{S}: \text{Name} \rightarrow \text{Symbol}$
```，其中：

- ```latex
$\text{Name}$
``` 是标识符集合
- ```latex
$\text{Symbol}$
``` 是符号信息集合

---

## 2. 形式化定义

### 2.1 语义域

**定义 2.1** (语义域): 语义域 ```latex
$\mathcal{D}$
``` 是程序语义值的集合，定义为：

$```latex
$\mathcal{D} = \mathcal{D}_{\text{int}} \cup \mathcal{D}_{\text{bool}} \cup \mathcal{D}_{\text{string}} \cup \mathcal{D}_{\text{func}} \cup \mathcal{D}_{\text{ref}}$
```$

其中：

- ```latex
$\mathcal{D}_{\text{int}}$
``` 是整数域
- ```latex
$\mathcal{D}_{\text{bool}}$
``` 是布尔域
- ```latex
$\mathcal{D}_{\text{string}}$
``` 是字符串域
- ```latex
$\mathcal{D}_{\text{func}}$
``` 是函数域
- ```latex
$\mathcal{D}_{\text{ref}}$
``` 是引用域

**定义 2.2** (语义环境): 语义环境 ```latex
$\rho$
``` 是一个函数：

$```latex
$\rho: \text{Var} \rightarrow \mathcal{D}$
```$

### 2.2 语义函数

**定义 2.3** (表达式语义函数): 表达式 ```latex
$e$
``` 的语义函数 ```latex
$\mathcal{E}[\![e]\!]$
``` 定义为：

$```latex
$\mathcal{E}[\![e]\!]: \text{Env} \rightarrow \mathcal{D}$
```$

**定义 2.4** (语句语义函数): 语句 ```latex
$s$
``` 的语义函数 ```latex
$\mathcal{S}[\![s]\!]$
``` 定义为：

$```latex
$\mathcal{S}[\![s]\!]: \text{Env} \rightarrow \text{Env}$
```$

### 2.3 类型系统

**定义 2.5** (类型): 类型 ```latex
$\tau$
``` 递归定义为：

$```latex
$\tau ::= \text{int} \mid \text{bool} \mid \text{string} \mid \tau_1 \rightarrow \tau_2 \mid \tau_1 \times \tau_2 \mid \text{ref } \tau$
```$

**定义 2.6** (类型环境): 类型环境 ```latex
$\Gamma$
``` 是一个函数：

$```latex
$\Gamma: \text{Var} \rightarrow \text{Type}$
```$

---

## 3. 类型检查

### 3.1 类型推导

**定义 3.1** (类型推导关系): 类型推导关系 ```latex
$\Gamma \vdash e : \tau$
``` 表示在类型环境 ```latex
$\Gamma$
``` 下，表达式 ```latex
$e$
``` 具有类型 ```latex
$\tau$
```。

**类型推导规则**:

**变量规则**:
$```latex
$\frac{x : \tau \in \Gamma}{\Gamma \vdash x : \tau}$
```$

**常量规则**:
$```latex
$\frac{}{\Gamma \vdash n : \text{int}} \quad \frac{}{\Gamma \vdash \text{true} : \text{bool}} \quad \frac{}{\Gamma \vdash \text{false} : \text{bool}}$
```$

**函数应用规则**:
$```latex
$\frac{\Gamma \vdash e_1 : \tau_1 \rightarrow \tau_2 \quad \Gamma \vdash e_2 : \tau_1}{\Gamma \vdash e_1(e_2) : \tau_2}$
```$

**条件表达式规则**:
$```latex
$\frac{\Gamma \vdash e_1 : \text{bool} \quad \Gamma \vdash e_2 : \tau \quad \Gamma \vdash e_3 : \tau}{\Gamma \vdash \text{if } e_1 \text{ then } e_2 \text{ else } e_3 : \tau}$
```$

### 3.2 类型统一

**定义 3.2** (类型统一): 类型统一是寻找类型变量的替换，使得两个类型表达式相等的过程。

**统一算法**:

1. 如果两个类型都是基本类型且相等，则统一成功
2. 如果一个类型是类型变量，则用另一个类型替换它
3. 如果两个类型都是函数类型，则递归统一参数和返回值类型
4. 否则统一失败

**定理 3.1** (统一算法正确性): 如果类型统一算法成功，则存在一个替换使得两个类型相等。

### 3.3 类型安全

**定义 3.3** (类型安全): 程序是类型安全的，如果所有表达式都有正确的类型，且类型检查通过。

**定理 3.2** (类型安全定理): 如果程序 ```latex
$P$
``` 是类型安全的，则 ```latex
$P$
``` 不会产生类型错误。

---

## 4. 作用域分析

### 4.1 静态作用域

**定义 4.1** (静态作用域): 在静态作用域中，变量的可见性在编译时确定。

**作用域规则**:

- 内层作用域可以访问外层作用域的变量
- 内层作用域可以隐藏外层作用域的同名变量
- 变量的生命周期从其声明开始到作用域结束

### 4.2 动态作用域

**定义 4.2** (动态作用域): 在动态作用域中，变量的可见性在运行时确定。

**特点**:

- 变量的可见性取决于程序的执行路径
- 函数调用时，被调用函数可以访问调用者的局部变量
- 实现相对复杂，现代语言较少使用

### 4.3 闭包处理

**定义 4.3** (闭包): 闭包是一个函数和其捕获的环境的组合。

**闭包环境**: 闭包捕获的环境包括：

- 函数定义时的局部变量
- 外层作用域的变量
- 全局变量

---

## 5. Go语言实现

### 5.1 符号表

```go
package semanticanalysis

import (
    "fmt"
    "strings"
)

// Symbol 符号信息
type Symbol struct {
    Name     string
    Type     Type
    Kind     SymbolKind
    Scope    *Scope
    Line     int
    Column   int
}

// SymbolKind 符号类型
type SymbolKind int

const (
    KindVariable SymbolKind = iota
    KindFunction
    KindType
    KindConstant
)

// Scope 作用域
type Scope struct {
    Parent    *Scope
    Symbols   map[string]*Symbol
    Children  []*Scope
    Level     int
}

// NewScope 创建新作用域
func NewScope(parent *Scope) *Scope {
    level := 0
    if parent != nil {
        level = parent.Level + 1
    }
    
    return &Scope{
        Parent:   parent,
        Symbols:  make(map[string]*Symbol),
        Children: make([]*Scope, 0),
        Level:    level,
    }
}

// Define 定义符号
func (s *Scope) Define(name string, symbolType Type, kind SymbolKind, line, column int) *Symbol {
    symbol := &Symbol{
        Name:   name,
        Type:   symbolType,
        Kind:   kind,
        Scope:  s,
        Line:   line,
        Column: column,
    }
    
    s.Symbols[name] = symbol
    return symbol
}

// Resolve 解析符号
func (s *Scope) Resolve(name string) (*Symbol, error) {
    if symbol, exists := s.Symbols[name]; exists {
        return symbol, nil
    }
    
    if s.Parent != nil {
        return s.Parent.Resolve(name)
    }
    
    return nil, fmt.Errorf("undefined symbol: %s", name)
}

// SymbolTable 符号表
type SymbolTable struct {
    GlobalScope *Scope
    CurrentScope *Scope
}

// NewSymbolTable 创建符号表
func NewSymbolTable() *SymbolTable {
    globalScope := NewScope(nil)
    return &SymbolTable{
        GlobalScope:  globalScope,
        CurrentScope: globalScope,
    }
}

// EnterScope 进入作用域
func (st *SymbolTable) EnterScope() {
    newScope := NewScope(st.CurrentScope)
    st.CurrentScope.Children = append(st.CurrentScope.Children, newScope)
    st.CurrentScope = newScope
}

// ExitScope 退出作用域
func (st *SymbolTable) ExitScope() {
    if st.CurrentScope.Parent != nil {
        st.CurrentScope = st.CurrentScope.Parent
    }
}
```

### 5.2 类型检查器

```go
// Type 类型接口
type Type interface {
    String() string
    Equals(other Type) bool
    IsAssignableTo(target Type) bool
}

// BasicType 基本类型
type BasicType struct {
    Name string
}

func (bt *BasicType) String() string { return bt.Name }
func (bt *BasicType) Equals(other Type) bool {
    if otherBT, ok := other.(*BasicType); ok {
        return bt.Name == otherBT.Name
    }
    return false
}
func (bt *BasicType) IsAssignableTo(target Type) bool {
    return bt.Equals(target)
}

// FunctionType 函数类型
type FunctionType struct {
    Parameters []Type
    ReturnType Type
}

func (ft *FunctionType) String() string {
    paramStrs := make([]string, len(ft.Parameters))
    for i, param := range ft.Parameters {
        paramStrs[i] = param.String()
    }
    return fmt.Sprintf("(%s) -> %s", strings.Join(paramStrs, ", "), ft.ReturnType.String())
}

func (ft *FunctionType) Equals(other Type) bool {
    if otherFT, ok := other.(*FunctionType); ok {
        if !ft.ReturnType.Equals(otherFT.ReturnType) {
            return false
        }
        if len(ft.Parameters) != len(otherFT.Parameters) {
            return false
        }
        for i, param := range ft.Parameters {
            if !param.Equals(otherFT.Parameters[i]) {
                return false
            }
        }
        return true
    }
    return false
}

func (ft *FunctionType) IsAssignableTo(target Type) bool {
    return ft.Equals(target)
}

// TypeChecker 类型检查器
type TypeChecker struct {
    symbolTable *SymbolTable
    errors      []string
}

// NewTypeChecker 创建类型检查器
func NewTypeChecker(symbolTable *SymbolTable) *TypeChecker {
    return &TypeChecker{
        symbolTable: symbolTable,
        errors:      make([]string, 0),
    }
}

// CheckExpression 检查表达式类型
func (tc *TypeChecker) CheckExpression(expr Expression) Type {
    switch e := expr.(type) {
    case *LiteralExpression:
        return tc.checkLiteral(e)
    case *VariableExpression:
        return tc.checkVariable(e)
    case *BinaryExpression:
        return tc.checkBinary(e)
    case *FunctionCallExpression:
        return tc.checkFunctionCall(e)
    default:
        tc.addError(fmt.Sprintf("unknown expression type: %T", expr))
        return &BasicType{Name: "error"}
    }
}

// checkLiteral 检查字面量类型
func (tc *TypeChecker) checkLiteral(lit *LiteralExpression) Type {
    switch lit.Value.(type) {
    case int:
        return &BasicType{Name: "int"}
    case float64:
        return &BasicType{Name: "float"}
    case string:
        return &BasicType{Name: "string"}
    case bool:
        return &BasicType{Name: "bool"}
    default:
        tc.addError(fmt.Sprintf("unknown literal type: %T", lit.Value))
        return &BasicType{Name: "error"}
    }
}

// checkVariable 检查变量类型
func (tc *TypeChecker) checkVariable(varExpr *VariableExpression) Type {
    symbol, err := tc.symbolTable.CurrentScope.Resolve(varExpr.Name)
    if err != nil {
        tc.addError(err.Error())
        return &BasicType{Name: "error"}
    }
    return symbol.Type
}

// checkBinary 检查二元表达式类型
func (tc *TypeChecker) checkBinary(bin *BinaryExpression) Type {
    leftType := tc.CheckExpression(bin.Left)
    rightType := tc.CheckExpression(bin.Right)
    
    // 检查操作数类型兼容性
    if !tc.isCompatible(leftType, rightType, bin.Operator) {
        tc.addError(fmt.Sprintf("incompatible types for operator %s: %s and %s", 
            bin.Operator, leftType.String(), rightType.String()))
        return &BasicType{Name: "error"}
    }
    
    // 返回结果类型
    return tc.getResultType(leftType, rightType, bin.Operator)
}

// checkFunctionCall 检查函数调用类型
func (tc *TypeChecker) checkFunctionCall(call *FunctionCallExpression) Type {
    // 解析函数名
    symbol, err := tc.symbolTable.CurrentScope.Resolve(call.FunctionName)
    if err != nil {
        tc.addError(err.Error())
        return &BasicType{Name: "error"}
    }
    
    // 检查是否为函数
    if symbol.Kind != KindFunction {
        tc.addError(fmt.Sprintf("%s is not a function", call.FunctionName))
        return &BasicType{Name: "error"}
    }
    
    // 检查函数类型
    funcType, ok := symbol.Type.(*FunctionType)
    if !ok {
        tc.addError(fmt.Sprintf("%s is not a function", call.FunctionName))
        return &BasicType{Name: "error"}
    }
    
    // 检查参数数量和类型
    if len(call.Arguments) != len(funcType.Parameters) {
        tc.addError(fmt.Sprintf("function %s expects %d arguments, got %d", 
            call.FunctionName, len(funcType.Parameters), len(call.Arguments)))
        return &BasicType{Name: "error"}
    }
    
    for i, arg := range call.Arguments {
        argType := tc.CheckExpression(arg)
        if !argType.IsAssignableTo(funcType.Parameters[i]) {
            tc.addError(fmt.Sprintf("argument %d of function %s has wrong type: expected %s, got %s", 
                i+1, call.FunctionName, funcType.Parameters[i].String(), argType.String()))
        }
    }
    
    return funcType.ReturnType
}

// isCompatible 检查类型兼容性
func (tc *TypeChecker) isCompatible(left, right Type, operator string) bool {
    // 基本类型兼容性检查
    if left.Equals(right) {
        return true
    }
    
    // 数值类型兼容性
    if tc.isNumeric(left) && tc.isNumeric(right) {
        return true
    }
    
    return false
}

// isNumeric 检查是否为数值类型
func (tc *TypeChecker) isNumeric(t Type) bool {
    if bt, ok := t.(*BasicType); ok {
        return bt.Name == "int" || bt.Name == "float"
    }
    return false
}

// getResultType 获取二元运算结果类型
func (tc *TypeChecker) getResultType(left, right Type, operator string) Type {
    switch operator {
    case "+", "-", "*", "/":
        if tc.isNumeric(left) && tc.isNumeric(right) {
            // 如果有一个是float，结果就是float
            if left.String() == "float" || right.String() == "float" {
                return &BasicType{Name: "float"}
            }
            return &BasicType{Name: "int"}
        }
    case "==", "!=", "<", ">", "<=", ">=":
        return &BasicType{Name: "bool"}
    case "&&", "||":
        return &BasicType{Name: "bool"}
    }
    
    return &BasicType{Name: "error"}
}

// addError 添加错误
func (tc *TypeChecker) addError(message string) {
    tc.errors = append(tc.errors, message)
}

// GetErrors 获取错误列表
func (tc *TypeChecker) GetErrors() []string {
    return tc.errors
}
```

### 5.3 语义分析器

```go
// Expression 表达式接口
type Expression interface {
    Accept(visitor ExpressionVisitor) interface{}
}

// ExpressionVisitor 表达式访问者
type ExpressionVisitor interface {
    VisitLiteral(expr *LiteralExpression) interface{}
    VisitVariable(expr *VariableExpression) interface{}
    VisitBinary(expr *BinaryExpression) interface{}
    VisitFunctionCall(expr *FunctionCallExpression) interface{}
}

// LiteralExpression 字面量表达式
type LiteralExpression struct {
    Value interface{}
}

func (le *LiteralExpression) Accept(visitor ExpressionVisitor) interface{} {
    return visitor.VisitLiteral(le)
}

// VariableExpression 变量表达式
type VariableExpression struct {
    Name string
}

func (ve *VariableExpression) Accept(visitor ExpressionVisitor) interface{} {
    return visitor.VisitVariable(ve)
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
    Left     Expression
    Operator string
    Right    Expression
}

func (be *BinaryExpression) Accept(visitor ExpressionVisitor) interface{} {
    return visitor.VisitBinary(be)
}

// FunctionCallExpression 函数调用表达式
type FunctionCallExpression struct {
    FunctionName string
    Arguments    []Expression
}

func (fc *FunctionCallExpression) Accept(visitor ExpressionVisitor) interface{} {
    return visitor.VisitFunctionCall(fc)
}

// SemanticAnalyzer 语义分析器
type SemanticAnalyzer struct {
    symbolTable *SymbolTable
    typeChecker *TypeChecker
    errors      []string
}

// NewSemanticAnalyzer 创建语义分析器
func NewSemanticAnalyzer() *SemanticAnalyzer {
    symbolTable := NewSymbolTable()
    typeChecker := NewTypeChecker(symbolTable)
    
    return &SemanticAnalyzer{
        symbolTable: symbolTable,
        typeChecker: typeChecker,
        errors:      make([]string, 0),
    }
}

// AnalyzeProgram 分析程序
func (sa *SemanticAnalyzer) AnalyzeProgram(program *Program) error {
    // 第一遍：收集声明
    sa.collectDeclarations(program)
    
    // 第二遍：类型检查和语义分析
    sa.analyzeStatements(program.Statements)
    
    // 检查是否有错误
    if len(sa.errors) > 0 {
        return fmt.Errorf("semantic analysis failed:\n%s", strings.Join(sa.errors, "\n"))
    }
    
    return nil
}

// collectDeclarations 收集声明
func (sa *SemanticAnalyzer) collectDeclarations(program *Program) {
    for _, stmt := range program.Statements {
        if decl, ok := stmt.(*VariableDeclaration); ok {
            sa.collectVariableDeclaration(decl)
        } else if funcDecl, ok := stmt.(*FunctionDeclaration); ok {
            sa.collectFunctionDeclaration(funcDecl)
        }
    }
}

// collectVariableDeclaration 收集变量声明
func (sa *SemanticAnalyzer) collectVariableDeclaration(decl *VariableDeclaration) {
    // 检查变量是否已声明
    if _, err := sa.symbolTable.CurrentScope.Resolve(decl.Name); err == nil {
        sa.addError(fmt.Sprintf("variable %s already declared", decl.Name))
        return
    }
    
    // 确定变量类型
    var varType Type
    if decl.Type != nil {
        varType = decl.Type
    } else if decl.Initializer != nil {
        varType = sa.typeChecker.CheckExpression(decl.Initializer)
    } else {
        sa.addError(fmt.Sprintf("variable %s must have a type or initializer", decl.Name))
        return
    }
    
    // 定义变量
    sa.symbolTable.CurrentScope.Define(decl.Name, varType, KindVariable, decl.Line, decl.Column)
}

// collectFunctionDeclaration 收集函数声明
func (sa *SemanticAnalyzer) collectFunctionDeclaration(decl *FunctionDeclaration) {
    // 检查函数是否已声明
    if _, err := sa.symbolTable.CurrentScope.Resolve(decl.Name); err == nil {
        sa.addError(fmt.Sprintf("function %s already declared", decl.Name))
        return
    }
    
    // 创建函数类型
    paramTypes := make([]Type, len(decl.Parameters))
    for i, param := range decl.Parameters {
        paramTypes[i] = param.Type
    }
    
    funcType := &FunctionType{
        Parameters: paramTypes,
        ReturnType: decl.ReturnType,
    }
    
    // 定义函数
    sa.symbolTable.CurrentScope.Define(decl.Name, funcType, KindFunction, decl.Line, decl.Column)
}

// analyzeStatements 分析语句
func (sa *SemanticAnalyzer) analyzeStatements(statements []Statement) {
    for _, stmt := range statements {
        sa.analyzeStatement(stmt)
    }
}

// analyzeStatement 分析单个语句
func (sa *SemanticAnalyzer) analyzeStatement(stmt Statement) {
    switch s := stmt.(type) {
    case *VariableDeclaration:
        sa.analyzeVariableDeclaration(s)
    case *AssignmentStatement:
        sa.analyzeAssignment(s)
    case *ExpressionStatement:
        sa.analyzeExpressionStatement(s)
    case *IfStatement:
        sa.analyzeIfStatement(s)
    case *WhileStatement:
        sa.analyzeWhileStatement(s)
    case *ReturnStatement:
        sa.analyzeReturnStatement(s)
    default:
        sa.addError(fmt.Sprintf("unknown statement type: %T", stmt))
    }
}

// analyzeVariableDeclaration 分析变量声明
func (sa *SemanticAnalyzer) analyzeVariableDeclaration(decl *VariableDeclaration) {
    if decl.Initializer != nil {
        initType := sa.typeChecker.CheckExpression(decl.Initializer)
        if decl.Type != nil && !initType.IsAssignableTo(decl.Type) {
            sa.addError(fmt.Sprintf("cannot assign %s to variable %s of type %s", 
                initType.String(), decl.Name, decl.Type.String()))
        }
    }
}

// analyzeAssignment 分析赋值语句
func (sa *SemanticAnalyzer) analyzeAssignment(assign *AssignmentStatement) {
    // 检查变量是否存在
    symbol, err := sa.symbolTable.CurrentScope.Resolve(assign.VariableName)
    if err != nil {
        sa.addError(err.Error())
        return
    }
    
    // 检查变量类型
    if symbol.Kind != KindVariable {
        sa.addError(fmt.Sprintf("%s is not a variable", assign.VariableName))
        return
    }
    
    // 检查赋值类型兼容性
    valueType := sa.typeChecker.CheckExpression(assign.Value)
    if !valueType.IsAssignableTo(symbol.Type) {
        sa.addError(fmt.Sprintf("cannot assign %s to variable %s of type %s", 
            valueType.String(), assign.VariableName, symbol.Type.String()))
    }
}

// analyzeExpressionStatement 分析表达式语句
func (sa *SemanticAnalyzer) analyzeExpressionStatement(exprStmt *ExpressionStatement) {
    sa.typeChecker.CheckExpression(exprStmt.Expression)
}

// analyzeIfStatement 分析if语句
func (sa *SemanticAnalyzer) analyzeIfStatement(ifStmt *IfStatement) {
    // 检查条件表达式类型
    condType := sa.typeChecker.CheckExpression(ifStmt.Condition)
    if !condType.Equals(&BasicType{Name: "bool"}) {
        sa.addError("if condition must be boolean")
    }
    
    // 分析then分支
    sa.symbolTable.EnterScope()
    sa.analyzeStatements(ifStmt.ThenBranch)
    sa.symbolTable.ExitScope()
    
    // 分析else分支
    if ifStmt.ElseBranch != nil {
        sa.symbolTable.EnterScope()
        sa.analyzeStatements(ifStmt.ElseBranch)
        sa.symbolTable.ExitScope()
    }
}

// analyzeWhileStatement 分析while语句
func (sa *SemanticAnalyzer) analyzeWhileStatement(whileStmt *WhileStatement) {
    // 检查条件表达式类型
    condType := sa.typeChecker.CheckExpression(whileStmt.Condition)
    if !condType.Equals(&BasicType{Name: "bool"}) {
        sa.addError("while condition must be boolean")
    }
    
    // 分析循环体
    sa.symbolTable.EnterScope()
    sa.analyzeStatements(whileStmt.Body)
    sa.symbolTable.ExitScope()
}

// analyzeReturnStatement 分析return语句
func (sa *SemanticAnalyzer) analyzeReturnStatement(returnStmt *ReturnStatement) {
    if returnStmt.Value != nil {
        valueType := sa.typeChecker.CheckExpression(returnStmt.Value)
        // 这里需要检查返回值类型是否与函数返回类型匹配
        // 需要从上下文获取当前函数的返回类型
    }
}

// addError 添加错误
func (sa *SemanticAnalyzer) addError(message string) {
    sa.errors = append(sa.errors, message)
}

// GetErrors 获取错误列表
func (sa *SemanticAnalyzer) GetErrors() []string {
    return sa.errors
}
```

---

## 6. 应用示例

### 6.1 简单语言

```go
// Program 程序
type Program struct {
    Statements []Statement
}

// Statement 语句接口
type Statement interface {
    Accept(visitor StatementVisitor) interface{}
}

// StatementVisitor 语句访问者
type StatementVisitor interface {
    VisitVariableDeclaration(decl *VariableDeclaration) interface{}
    VisitAssignment(assign *AssignmentStatement) interface{}
    VisitExpression(expr *ExpressionStatement) interface{}
    VisitIf(ifStmt *IfStatement) interface{}
    VisitWhile(whileStmt *WhileStatement) interface{}
    VisitReturn(returnStmt *ReturnStatement) interface{}
}

// VariableDeclaration 变量声明
type VariableDeclaration struct {
    Name        string
    Type        Type
    Initializer Expression
    Line        int
    Column      int
}

func (vd *VariableDeclaration) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitVariableDeclaration(vd)
}

// AssignmentStatement 赋值语句
type AssignmentStatement struct {
    VariableName string
    Value        Expression
    Line         int
    Column       int
}

func (as *AssignmentStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitAssignment(as)
}

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
    Expression Expression
    Line       int
    Column     int
}

func (es *ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitExpression(es)
}

// IfStatement if语句
type IfStatement struct {
    Condition  Expression
    ThenBranch []Statement
    ElseBranch []Statement
    Line       int
    Column     int
}

func (is *IfStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitIf(is)
}

// WhileStatement while语句
type WhileStatement struct {
    Condition Expression
    Body      []Statement
    Line      int
    Column    int
}

func (ws *WhileStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitWhile(ws)
}

// ReturnStatement return语句
type ReturnStatement struct {
    Value Expression
    Line  int
    Column int
}

func (rs *ReturnStatement) Accept(visitor StatementVisitor) interface{} {
    return visitor.VisitReturn(rs)
}

// FunctionDeclaration 函数声明
type FunctionDeclaration struct {
    Name       string
    Parameters []*Parameter
    ReturnType Type
    Body       []Statement
    Line       int
    Column     int
}

// Parameter 参数
type Parameter struct {
    Name string
    Type Type
}

// 示例：分析简单程序
func ExampleSimpleProgram() {
    // 创建语义分析器
    analyzer := NewSemanticAnalyzer()
    
    // 创建简单程序
    program := &Program{
        Statements: []Statement{
            &VariableDeclaration{
                Name: "x",
                Type: &BasicType{Name: "int"},
                Initializer: &LiteralExpression{Value: 42},
            },
            &AssignmentStatement{
                VariableName: "x",
                Value: &BinaryExpression{
                    Left:     &VariableExpression{Name: "x"},
                    Operator: "+",
                    Right:    &LiteralExpression{Value: 1},
                },
            },
        },
    }
    
    // 分析程序
    err := analyzer.AnalyzeProgram(program)
    if err != nil {
        fmt.Printf("Semantic analysis failed: %v\n", err)
    } else {
        fmt.Println("Semantic analysis successful")
    }
}
```

### 6.2 函数式语言

```go
// 函数式语言的语义分析示例
func ExampleFunctionalLanguage() {
    analyzer := NewSemanticAnalyzer()
    
    // 定义函数类型
    intType := &BasicType{Name: "int"}
    funcType := &FunctionType{
        Parameters: []Type{intType, intType},
        ReturnType: intType,
    }
    
    // 创建函数声明
    functionDecl := &FunctionDeclaration{
        Name: "add",
        Parameters: []*Parameter{
            {Name: "a", Type: intType},
            {Name: "b", Type: intType},
        },
        ReturnType: intType,
        Body: []Statement{
            &ReturnStatement{
                Value: &BinaryExpression{
                    Left:     &VariableExpression{Name: "a"},
                    Operator: "+",
                    Right:    &VariableExpression{Name: "b"},
                },
            },
        },
    }
    
    // 创建函数调用
    functionCall := &FunctionCallExpression{
        FunctionName: "add",
        Arguments: []Expression{
            &LiteralExpression{Value: 1},
            &LiteralExpression{Value: 2},
        },
    }
    
    // 分析函数调用
    resultType := analyzer.typeChecker.CheckExpression(functionCall)
    fmt.Printf("Function call result type: %s\n", resultType.String())
}
```

### 6.3 面向对象语言

```go
// 面向对象语言的语义分析示例
func ExampleObjectOrientedLanguage() {
    analyzer := NewSemanticAnalyzer()
    
    // 定义类类型
    classType := &ClassType{
        Name: "Person",
        Fields: map[string]Type{
            "name": &BasicType{Name: "string"},
            "age":  &BasicType{Name: "int"},
        },
        Methods: map[string]*FunctionType{
            "getName": {
                Parameters: []Type{},
                ReturnType: &BasicType{Name: "string"},
            },
        },
    }
    
    // 创建对象实例
    objectExpr := &ObjectExpression{
        ClassName: "Person",
        Fields: map[string]Expression{
            "name": &LiteralExpression{Value: "John"},
            "age":  &LiteralExpression{Value: 30},
        },
    }
    
    // 分析对象表达式
    resultType := analyzer.typeChecker.CheckExpression(objectExpr)
    fmt.Printf("Object expression type: %s\n", resultType.String())
}

// ClassType 类类型
type ClassType struct {
    Name   string
    Fields map[string]Type
    Methods map[string]*FunctionType
}

func (ct *ClassType) String() string { return ct.Name }
func (ct *ClassType) Equals(other Type) bool {
    if otherCT, ok := other.(*ClassType); ok {
        return ct.Name == otherCT.Name
    }
    return false
}
func (ct *ClassType) IsAssignableTo(target Type) bool {
    return ct.Equals(target)
}

// ObjectExpression 对象表达式
type ObjectExpression struct {
    ClassName string
    Fields    map[string]Expression
}

func (oe *ObjectExpression) Accept(visitor ExpressionVisitor) interface{} {
    // 实现访问者模式
    return nil
}
```

---

## 7. 数学证明

### 7.1 类型安全定理

**定理 7.1** (类型安全定理): 如果程序 ```latex
$P$
``` 通过类型检查，则 ```latex
$P$
``` 在执行时不会产生类型错误。

**证明**: 使用结构归纳法证明。

**基础情况**: 对于基本表达式（字面量、变量），类型检查确保类型正确。

**归纳步骤**: 对于复合表达式，类型检查确保：

1. 子表达式的类型正确
2. 操作符的类型要求满足
3. 结果类型正确

因此，通过类型检查的程序在执行时不会产生类型错误。```latex
$\square$
```

### 7.2 类型推导定理

**定理 7.2** (类型推导定理): 如果表达式 ```latex
$e$
``` 有类型 ```latex
$\tau$
```，则存在类型推导 ```latex
$\Gamma \vdash e : \tau$
```。

**证明**: 使用结构归纳法。

**基础情况**: 对于基本表达式，类型推导规则直接给出类型。

**归纳步骤**: 对于复合表达式，使用相应的类型推导规则组合子表达式的类型推导。```latex
$\square$
```

### 7.3 语义一致性定理

**定理 7.3** (语义一致性定理): 如果两个表达式 ```latex
$e_1$
``` 和 ```latex
$e_2$
``` 类型相同且语义等价，则它们在所有上下文中可以互换。

**证明**: 使用语义等价的定义和类型安全定理。

由于 ```latex
$e_1$
``` 和 ```latex
$e_2$
``` 类型相同，它们在类型检查中表现相同。由于语义等价，它们在执行时产生相同的结果。因此，它们可以在所有上下文中互换。```latex
$\square$
```

---

## 总结

本文档提供了语义分析的完整形式化定义和Go语言实现。通过类型系统、作用域分析、语义检查等多重机制，确保了程序的语义正确性。

**关键特性**:

- 完整的类型系统
- 作用域管理
- 语义检查
- 错误报告
- 类型推导

**应用场景**:

- 编译器实现
- 静态分析工具
- 语言设计
- 程序验证
- 代码生成
