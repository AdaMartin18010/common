# 01-类型系统理论 (Type System Theory)

## 目录

- [01-类型系统理论 (Type System Theory)](#01-类型系统理论-type-system-theory)
  - [目录](#目录)
  - [1. 类型基础 (Type Foundations)](#1-类型基础-type-foundations)
    - [1.1 类型概念](#11-类型概念)
    - [1.2 类型分类](#12-类型分类)
    - [1.3 类型关系](#13-类型关系)
  - [2. 类型推导 (Type Inference)](#2-类型推导-type-inference)
    - [2.1 类型推导算法](#21-类型推导算法)
    - [2.2 统一算法](#22-统一算法)
    - [2.3 约束求解](#23-约束求解)
  - [3. 类型安全 (Type Safety)](#3-类型安全-type-safety)
    - [3.1 类型安全定义](#31-类型安全定义)
    - [3.2 进展和保持](#32-进展和保持)
    - [3.3 类型错误处理](#33-类型错误处理)
  - [4. 高级类型系统 (Advanced Type Systems)](#4-高级类型系统-advanced-type-systems)
    - [4.1 多态类型](#41-多态类型)
    - [4.2 依赖类型](#42-依赖类型)
    - [4.3 线性类型](#43-线性类型)
  - [5. Go语言类型系统](#5-go语言类型系统)
    - [5.1 基础类型](#51-基础类型)
    - [5.2 接口类型](#52-接口类型)
    - [5.3 泛型类型](#53-泛型类型)

## 1. 类型基础 (Type Foundations)

### 1.1 类型概念

**定义 1.1**: 类型
类型是一个集合，表示程序中值的集合。类型系统是用于防止程序错误的静态分析工具。

**定义 1.2**: 类型环境
类型环境 $\Gamma$ 是一个从变量到类型的映射：
$$\Gamma: Var \rightarrow Type$$

**定义 1.3**: 类型判断
类型判断是一个三元组 $\Gamma \vdash e: \tau$，表示在环境 $\Gamma$ 下，表达式 $e$ 具有类型 $\tau$。

**定义 1.4**: 类型规则
类型规则是推导类型判断的规则，通常采用自然演绎的形式：

$$\frac{\text{premises}}{\text{conclusion}}$$

### 1.2 类型分类

**定义 1.5**: 基本类型
基本类型包括：
- $\text{Bool}$：布尔类型
- $\text{Int}$：整数类型
- $\text{Float}$：浮点类型
- $\text{String}$：字符串类型
- $\text{Unit}$：单位类型

**定义 1.6**: 复合类型
复合类型包括：
- **函数类型**：$\tau_1 \rightarrow \tau_2$
- **积类型**：$\tau_1 \times \tau_2$
- **和类型**：$\tau_1 + \tau_2$
- **列表类型**：$\text{List}[\tau]$
- **引用类型**：$\text{Ref}[\tau]$

**定义 1.7**: 类型构造器
类型构造器是高阶类型，如：
- $\text{List}$：列表构造器
- $\text{Option}$：可选构造器
- $\text{Either}$：或构造器
- $\text{Map}$：映射构造器

### 1.3 类型关系

**定义 1.8**: 子类型关系
子类型关系 $\tau_1 \leq \tau_2$ 表示 $\tau_1$ 是 $\tau_2$ 的子类型。

**定义 1.9**: 类型等价
类型等价 $\tau_1 \equiv \tau_2$ 表示 $\tau_1$ 和 $\tau_2$ 是等价的。

**定理 1.1**: 子类型的基本性质
1. **自反性**：$\tau \leq \tau$
2. **传递性**：$\tau_1 \leq \tau_2 \wedge \tau_2 \leq \tau_3 \Rightarrow \tau_1 \leq \tau_3$
3. **反对称性**：$\tau_1 \leq \tau_2 \wedge \tau_2 \leq \tau_1 \Rightarrow \tau_1 \equiv \tau_2$

**定义 1.10**: 协变和逆变
- 类型构造器 $F$ 是协变的，如果 $\tau_1 \leq \tau_2 \Rightarrow F[\tau_1] \leq F[\tau_2]$
- 类型构造器 $F$ 是逆变的，如果 $\tau_1 \leq \tau_2 \Rightarrow F[\tau_2] \leq F[\tau_1]$
- 类型构造器 $F$ 是不变的，如果既不是协变也不是逆变

## 2. 类型推导 (Type Inference)

### 2.1 类型推导算法

**算法 2.1**: Hindley-Milner类型推导
Hindley-Milner类型推导算法是函数式编程语言中最常用的类型推导算法。

**定义 2.1**: 类型变量
类型变量 $\alpha, \beta, \gamma, \ldots$ 表示未知类型。

**定义 2.2**: 类型模式
类型模式是包含类型变量的类型表达式。

**算法 2.2**: 类型推导规则
```go
type TypeVariable struct {
    ID string
}

type Type struct {
    Kind TypeKind
    Args []Type
}

type TypeKind int

const (
    BasicType TypeKind = iota
    FunctionType
    ProductType
    SumType
    VariableType
    ConstructorType
)

type TypeEnvironment struct {
    Bindings map[string]Type
}

// InferType 类型推导
func InferType(expr Expression, env TypeEnvironment) (Type, error) {
    switch e := expr.(type) {
    case *Variable:
        return env.Lookup(e.Name)
    case *Application:
        return inferApplication(e, env)
    case *Abstraction:
        return inferAbstraction(e, env)
    case *Let:
        return inferLet(e, env)
    default:
        return nil, fmt.Errorf("unsupported expression type")
    }
}

// inferApplication 推导函数应用类型
func inferApplication(app *Application, env TypeEnvironment) (Type, error) {
    funType, err := InferType(app.Function, env)
    if err != nil {
        return nil, err
    }
    
    argType, err := InferType(app.Argument, env)
    if err != nil {
        return nil, err
    }
    
    // 检查函数类型是否为函数类型
    if funType.Kind != FunctionType {
        return nil, fmt.Errorf("expected function type, got %v", funType)
    }
    
    // 统一参数类型
    substitution, err := Unify(funType.Args[0], argType)
    if err != nil {
        return nil, fmt.Errorf("type mismatch: %v", err)
    }
    
    // 应用替换
    return ApplySubstitution(funType.Args[1], substitution), nil
}

// inferAbstraction 推导lambda抽象类型
func inferAbstraction(abs *Abstraction, env TypeEnvironment) (Type, error) {
    // 创建新的类型变量
    paramType := NewTypeVariable()
    
    // 扩展环境
    newEnv := env.Extend(abs.Parameter, paramType)
    
    // 推导函数体类型
    bodyType, err := InferType(abs.Body, newEnv)
    if err != nil {
        return nil, err
    }
    
    // 返回函数类型
    return &Type{
        Kind: FunctionType,
        Args: []Type{paramType, bodyType},
    }, nil
}
```

### 2.2 统一算法

**定义 2.3**: 类型统一
类型统一是找到类型变量的替换，使得两个类型表达式相等。

**算法 2.3**: Robinson统一算法
```go
type Substitution map[string]Type

// Unify 统一两个类型
func Unify(t1, t2 Type) (Substitution, error) {
    switch {
    case t1.Kind == VariableType:
        return unifyVariable(t1, t2)
    case t2.Kind == VariableType:
        return unifyVariable(t2, t1)
    case t1.Kind == FunctionType && t2.Kind == FunctionType:
        return unifyFunction(t1, t2)
    case t1.Kind == ConstructorType && t2.Kind == ConstructorType:
        return unifyConstructor(t1, t2)
    default:
        if t1.Kind == t2.Kind {
            return unifySameKind(t1, t2)
        }
        return nil, fmt.Errorf("cannot unify %v and %v", t1, t2)
    }
}

// unifyVariable 统一类型变量
func unifyVariable(variable Type, other Type) (Substitution, error) {
    // 检查变量是否出现在其他类型中
    if occursIn(variable, other) {
        return nil, fmt.Errorf("occurs check failed")
    }
    
    substitution := make(Substitution)
    substitution[variable.ID] = other
    return substitution, nil
}

// unifyFunction 统一函数类型
func unifyFunction(f1, f2 Type) (Substitution, error) {
    if len(f1.Args) != 2 || len(f2.Args) != 2 {
        return nil, fmt.Errorf("invalid function type")
    }
    
    // 统一参数类型
    s1, err := Unify(f1.Args[0], f2.Args[0])
    if err != nil {
        return nil, err
    }
    
    // 应用替换到返回类型
    t1 := ApplySubstitution(f1.Args[1], s1)
    t2 := ApplySubstitution(f2.Args[1], s1)
    
    // 统一返回类型
    s2, err := Unify(t1, t2)
    if err != nil {
        return nil, err
    }
    
    // 组合替换
    return ComposeSubstitutions(s1, s2)
}

// occursIn 检查变量是否出现在类型中
func occursIn(variable Type, other Type) bool {
    switch other.Kind {
    case VariableType:
        return variable.ID == other.ID
    case FunctionType, ConstructorType:
        for _, arg := range other.Args {
            if occursIn(variable, arg) {
                return true
            }
        }
    }
    return false
}

// ApplySubstitution 应用替换
func ApplySubstitution(t Type, s Substitution) Type {
    switch t.Kind {
    case VariableType:
        if replacement, exists := s[t.ID]; exists {
            return ApplySubstitution(replacement, s)
        }
        return t
    case FunctionType, ConstructorType:
        newArgs := make([]Type, len(t.Args))
        for i, arg := range t.Args {
            newArgs[i] = ApplySubstitution(arg, s)
        }
        return &Type{
            Kind: t.Kind,
            Args: newArgs,
        }
    default:
        return t
    }
}
```

### 2.3 约束求解

**定义 2.4**: 类型约束
类型约束是一个等式 $t_1 = t_2$，其中 $t_1$ 和 $t_2$ 是类型表达式。

**定义 2.5**: 约束系统
约束系统是类型约束的集合。

**算法 2.4**: 约束求解算法
```go
type Constraint struct {
    Left  Type
    Right Type
}

type ConstraintSystem struct {
    Constraints []Constraint
}

// SolveConstraints 求解约束系统
func (cs *ConstraintSystem) SolveConstraints() (Substitution, error) {
    substitution := make(Substitution)
    
    for len(cs.Constraints) > 0 {
        constraint := cs.Constraints[0]
        cs.Constraints = cs.Constraints[1:]
        
        newSubstitution, err := Unify(constraint.Left, constraint.Right)
        if err != nil {
            return nil, err
        }
        
        // 应用新替换到所有约束
        cs.applySubstitution(newSubstitution)
        
        // 组合替换
        substitution = ComposeSubstitutions(substitution, newSubstitution)
    }
    
    return substitution, nil
}

// applySubstitution 应用替换到约束系统
func (cs *ConstraintSystem) applySubstitution(s Substitution) {
    for i := range cs.Constraints {
        cs.Constraints[i].Left = ApplySubstitution(cs.Constraints[i].Left, s)
        cs.Constraints[i].Right = ApplySubstitution(cs.Constraints[i].Right, s)
    }
}
```

## 3. 类型安全 (Type Safety)

### 3.1 类型安全定义

**定义 3.1**: 类型安全
类型安全是指程序在运行时不会出现类型错误。

**定义 3.2**: 类型错误
类型错误包括：
- 类型不匹配
- 未定义变量
- 函数参数类型错误
- 返回值类型错误

**定义 3.3**: 类型安全定理
如果 $\Gamma \vdash e: \tau$ 且 $e \rightarrow e'$，则 $\Gamma \vdash e': \tau$。

### 3.2 进展和保持

**定理 3.1**: 进展定理 (Progress)
如果 $\vdash e: \tau$，则要么 $e$ 是值，要么存在 $e'$ 使得 $e \rightarrow e'$。

**定理 3.2**: 保持定理 (Preservation)
如果 $\Gamma \vdash e: \tau$ 且 $e \rightarrow e'$，则 $\Gamma \vdash e': \tau$。

**定理 3.3**: 类型安全定理
如果 $\vdash e: \tau$，则 $e$ 不会陷入错误状态。

**证明**:
1. 由进展定理，$e$ 要么是值，要么可以继续求值
2. 由保持定理，求值过程中类型保持不变
3. 因此 $e$ 不会出现类型错误

### 3.3 类型错误处理

**算法 3.1**: 类型错误检测
```go
type TypeError struct {
    Message string
    Location Position
}

type TypeChecker struct {
    errors []TypeError
}

// CheckType 类型检查
func (tc *TypeChecker) CheckType(expr Expression, env TypeEnvironment) (Type, error) {
    switch e := expr.(type) {
    case *Variable:
        return tc.checkVariable(e, env)
    case *Application:
        return tc.checkApplication(e, env)
    case *Abstraction:
        return tc.checkAbstraction(e, env)
    case *Let:
        return tc.checkLet(e, env)
    default:
        return nil, fmt.Errorf("unsupported expression type")
    }
}

// checkVariable 检查变量
func (tc *TypeChecker) checkVariable(v *Variable, env TypeEnvironment) (Type, error) {
    if t, exists := env.Bindings[v.Name]; exists {
        return t, nil
    }
    
    tc.addError(TypeError{
        Message:  fmt.Sprintf("undefined variable: %s", v.Name),
        Location: v.Position,
    })
    return nil, fmt.Errorf("undefined variable: %s", v.Name)
}

// checkApplication 检查函数应用
func (tc *TypeChecker) checkApplication(app *Application, env TypeEnvironment) (Type, error) {
    funType, err := tc.CheckType(app.Function, env)
    if err != nil {
        return nil, err
    }
    
    argType, err := tc.CheckType(app.Argument, env)
    if err != nil {
        return nil, err
    }
    
    if funType.Kind != FunctionType {
        tc.addError(TypeError{
            Message:  fmt.Sprintf("expected function type, got %v", funType),
            Location: app.Position,
        })
        return nil, fmt.Errorf("expected function type, got %v", funType)
    }
    
    if !tc.isSubtype(argType, funType.Args[0]) {
        tc.addError(TypeError{
            Message:  fmt.Sprintf("type mismatch: expected %v, got %v", funType.Args[0], argType),
            Location: app.Position,
        })
        return nil, fmt.Errorf("type mismatch")
    }
    
    return funType.Args[1], nil
}

// addError 添加类型错误
func (tc *TypeChecker) addError(err TypeError) {
    tc.errors = append(tc.errors, err)
}

// GetErrors 获取所有类型错误
func (tc *TypeChecker) GetErrors() []TypeError {
    return tc.errors
}
```

## 4. 高级类型系统 (Advanced Type Systems)

### 4.1 多态类型

**定义 4.1**: 参数多态
参数多态允许类型包含类型变量，如 $\forall \alpha. \alpha \rightarrow \alpha$。

**定义 4.2**: 特设多态
特设多态允许同一函数名用于不同类型，如函数重载。

**算法 4.1**: 多态类型推导
```go
type PolymorphicType struct {
    Variables []string
    Type      Type
}

// Generalize 泛化类型
func Generalize(t Type, env TypeEnvironment) PolymorphicType {
    freeVars := freeTypeVariables(t)
    envVars := env.FreeTypeVariables()
    
    // 计算可以泛化的变量
    generalizable := make([]string, 0)
    for _, v := range freeVars {
        if !contains(envVars, v) {
            generalizable = append(generalizable, v)
        }
    }
    
    return PolymorphicType{
        Variables: generalizable,
        Type:      t,
    }
}

// Instantiate 实例化多态类型
func Instantiate(poly PolymorphicType) Type {
    substitution := make(Substitution)
    for _, v := range poly.Variables {
        substitution[v] = NewTypeVariable()
    }
    
    return ApplySubstitution(poly.Type, substitution)
}
```

### 4.2 依赖类型

**定义 4.3**: 依赖类型
依赖类型允许类型依赖于值，如 $\Pi x: \text{Nat}. \text{Vec}[x]$。

**定义 4.4**: 依赖函数类型
依赖函数类型 $\Pi x: A. B[x]$ 表示对于所有 $x: A$，返回类型为 $B[x]$。

**算法 4.2**: 依赖类型检查
```go
type DependentType struct {
    Parameter string
    Domain    Type
    Codomain  func(Type) Type
}

// CheckDependentFunction 检查依赖函数
func CheckDependentFunction(abs *Abstraction, expectedType Type) error {
    if expectedType.Kind != DependentType {
        return fmt.Errorf("expected dependent function type")
    }
    
    // 检查参数类型
    if !isSubtype(abs.ParameterType, expectedType.Domain) {
        return fmt.Errorf("parameter type mismatch")
    }
    
    // 扩展环境
    newEnv := env.Extend(abs.Parameter, abs.ParameterType)
    
    // 检查函数体类型
    expectedBodyType := expectedType.Codomain(abs.ParameterType)
    bodyType, err := CheckType(abs.Body, newEnv)
    if err != nil {
        return err
    }
    
    if !isSubtype(bodyType, expectedBodyType) {
        return fmt.Errorf("body type mismatch")
    }
    
    return nil
}
```

### 4.3 线性类型

**定义 4.5**: 线性类型
线性类型确保值被使用且仅使用一次。

**定义 4.6**: 线性函数类型
线性函数类型 $A \multimap B$ 表示函数使用参数一次并返回结果。

**算法 4.3**: 线性类型检查
```go
type LinearType struct {
    Type     Type
    Usage    Usage
}

type Usage int

const (
    Unused Usage = iota
    Used
    Consumed
)

// CheckLinearType 检查线性类型
func CheckLinearType(expr Expression, env LinearEnvironment) (LinearType, error) {
    switch e := expr.(type) {
    case *Variable:
        return checkLinearVariable(e, env)
    case *Application:
        return checkLinearApplication(e, env)
    case *Abstraction:
        return checkLinearAbstraction(e, env)
    default:
        return LinearType{}, fmt.Errorf("unsupported expression")
    }
}

// checkLinearVariable 检查线性变量
func checkLinearVariable(v *Variable, env LinearEnvironment) (LinearType, error) {
    linearType, exists := env.Bindings[v.Name]
    if !exists {
        return LinearType{}, fmt.Errorf("undefined variable: %s", v.Name)
    }
    
    if linearType.Usage == Consumed {
        return LinearType{}, fmt.Errorf("variable %s already consumed", v.Name)
    }
    
    // 标记为已使用
    env.Bindings[v.Name] = LinearType{
        Type:  linearType.Type,
        Usage: Consumed,
    }
    
    return linearType, nil
}
```

## 5. Go语言类型系统

### 5.1 基础类型

**定义 5.1**: Go基础类型
Go语言的基础类型包括：
- `bool`：布尔类型
- `int`, `int8`, `int16`, `int32`, `int64`：整数类型
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`：无符号整数类型
- `float32`, `float64`：浮点类型
- `complex64`, `complex128`：复数类型
- `string`：字符串类型
- `byte`：字节类型（`uint8`的别名）
- `rune`：Unicode字符类型（`int32`的别名）

**算法 5.1**: Go类型检查
```go
package types

import (
    "go/ast"
    "go/token"
)

type GoTypeChecker struct {
    errors []TypeError
    env    TypeEnvironment
}

type TypeEnvironment struct {
    variables map[string]Type
    functions map[string]FunctionType
    types     map[string]Type
}

type Type interface {
    String() string
    IsAssignableTo(other Type) bool
}

type BasicType struct {
    Name string
}

type FunctionType struct {
    Params []Type
    Result Type
}

type StructType struct {
    Fields map[string]Type
}

type InterfaceType struct {
    Methods map[string]FunctionType
}

type SliceType struct {
    Element Type
}

type MapType struct {
    Key   Type
    Value Type
}

type PointerType struct {
    Base Type
}

type ChannelType struct {
    Element Type
    Dir     ChanDir
}

type ChanDir int

const (
    SendRecv ChanDir = iota
    SendOnly
    RecvOnly
)

// CheckExpression 检查表达式类型
func (tc *GoTypeChecker) CheckExpression(expr ast.Expr) (Type, error) {
    switch e := expr.(type) {
    case *ast.Ident:
        return tc.checkIdentifier(e)
    case *ast.BasicLit:
        return tc.checkBasicLiteral(e)
    case *ast.BinaryExpr:
        return tc.checkBinaryExpression(e)
    case *ast.CallExpr:
        return tc.checkCallExpression(e)
    case *ast.SelectorExpr:
        return tc.checkSelectorExpression(e)
    case *ast.IndexExpr:
        return tc.checkIndexExpression(e)
    case *ast.SliceExpr:
        return tc.checkSliceExpression(e)
    case *ast.TypeAssertExpr:
        return tc.checkTypeAssertion(e)
    case *ast.UnaryExpr:
        return tc.checkUnaryExpression(e)
    default:
        return nil, fmt.Errorf("unsupported expression type: %T", expr)
    }
}

// checkIdentifier 检查标识符
func (tc *GoTypeChecker) checkIdentifier(ident *ast.Ident) (Type, error) {
    // 检查变量
    if t, exists := tc.env.variables[ident.Name]; exists {
        return t, nil
    }
    
    // 检查函数
    if t, exists := tc.env.functions[ident.Name]; exists {
        return t, nil
    }
    
    // 检查类型
    if t, exists := tc.env.types[ident.Name]; exists {
        return t, nil
    }
    
    tc.addError(TypeError{
        Message: fmt.Sprintf("undefined identifier: %s", ident.Name),
        Pos:     ident.Pos(),
    })
    return nil, fmt.Errorf("undefined identifier: %s", ident.Name)
}

// checkBinaryExpression 检查二元表达式
func (tc *GoTypeChecker) checkBinaryExpression(expr *ast.BinaryExpr) (Type, error) {
    leftType, err := tc.CheckExpression(expr.X)
    if err != nil {
        return nil, err
    }
    
    rightType, err := tc.CheckExpression(expr.Y)
    if err != nil {
        return nil, err
    }
    
    // 检查操作符兼容性
    resultType, err := tc.checkBinaryOperator(expr.Op, leftType, rightType)
    if err != nil {
        tc.addError(TypeError{
            Message: err.Error(),
            Pos:     expr.Pos(),
        })
        return nil, err
    }
    
    return resultType, nil
}

// checkBinaryOperator 检查二元操作符
func (tc *GoTypeChecker) checkBinaryOperator(op token.Token, left, right Type) (Type, error) {
    switch op {
    case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
        return tc.checkArithmeticOperator(op, left, right)
    case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
        return tc.checkComparisonOperator(op, left, right)
    case token.LAND, token.LOR:
        return tc.checkLogicalOperator(op, left, right)
    default:
        return nil, fmt.Errorf("unsupported binary operator: %s", op)
    }
}

// checkArithmeticOperator 检查算术操作符
func (tc *GoTypeChecker) checkArithmeticOperator(op token.Token, left, right Type) (Type, error) {
    // 检查类型兼容性
    if !tc.isNumericType(left) || !tc.isNumericType(right) {
        return nil, fmt.Errorf("arithmetic operator requires numeric types")
    }
    
    // 类型提升规则
    resultType := tc.promoteTypes(left, right)
    
    // 对于除法，结果类型可能是浮点
    if op == token.QUO && tc.isIntegerType(left) && tc.isIntegerType(right) {
        if tc.isIntegerDivision(expr) {
            return resultType, nil
        }
        return tc.promoteToFloat(resultType), nil
    }
    
    return resultType, nil
}

// isNumericType 检查是否为数值类型
func (tc *GoTypeChecker) isNumericType(t Type) bool {
    switch t.(type) {
    case *BasicType:
        name := t.(*BasicType).Name
        return name == "int" || name == "int8" || name == "int16" || name == "int32" || name == "int64" ||
               name == "uint" || name == "uint8" || name == "uint16" || name == "uint32" || name == "uint64" ||
               name == "float32" || name == "float64" || name == "complex64" || name == "complex128"
    }
    return false
}

// promoteTypes 类型提升
func (tc *GoTypeChecker) promoteTypes(t1, t2 Type) Type {
    // 实现Go语言的类型提升规则
    // 1. 如果两个操作数都是无类型常量，结果是无类型常量
    // 2. 如果两个操作数都是无类型浮点常量，结果是float64
    // 3. 如果两个操作数都是无类型复数常量，结果是complex128
    // 4. 否则，结果类型是操作数的公共类型
    
    // 简化实现
    return t1
}
```

### 5.2 接口类型

**定义 5.2**: Go接口类型
Go接口类型定义了一组方法签名，任何实现了这些方法的类型都实现了该接口。

**算法 5.2**: 接口类型检查
```go
// checkInterfaceImplementation 检查接口实现
func (tc *GoTypeChecker) checkInterfaceImplementation(typ Type, iface *InterfaceType) error {
    for methodName, methodType := range iface.Methods {
        // 检查类型是否有对应的方法
        if !tc.hasMethod(typ, methodName, methodType) {
            return fmt.Errorf("type %v does not implement method %s", typ, methodName)
        }
    }
    return nil
}

// hasMethod 检查类型是否有指定方法
func (tc *GoTypeChecker) hasMethod(typ Type, methodName string, expectedType FunctionType) bool {
    switch t := typ.(type) {
    case *StructType:
        // 检查结构体方法
        return tc.checkStructMethod(t, methodName, expectedType)
    case *BasicType:
        // 检查基本类型方法
        return tc.checkBasicTypeMethod(t, methodName, expectedType)
    case *InterfaceType:
        // 检查接口方法
        return tc.checkInterfaceMethod(t, methodName, expectedType)
    default:
        return false
    }
}

// checkStructMethod 检查结构体方法
func (tc *GoTypeChecker) checkStructMethod(structType *StructType, methodName string, expectedType FunctionType) bool {
    // 查找方法
    method, exists := tc.findMethod(structType, methodName)
    if !exists {
        return false
    }
    
    // 检查方法签名
    return tc.checkMethodSignature(method, expectedType)
}

// checkMethodSignature 检查方法签名
func (tc *GoTypeChecker) checkMethodSignature(actual, expected FunctionType) bool {
    // 检查参数数量
    if len(actual.Params) != len(expected.Params) {
        return false
    }
    
    // 检查参数类型
    for i, expectedParam := range expected.Params {
        if !actual.Params[i].IsAssignableTo(expectedParam) {
            return false
        }
    }
    
    // 检查返回类型
    return expected.Result.IsAssignableTo(actual.Result)
}
```

### 5.3 泛型类型

**定义 5.3**: Go泛型类型
Go 1.18引入的泛型支持参数化类型，使用类型参数来定义通用类型。

**算法 5.3**: 泛型类型检查
```go
type GenericType struct {
    Name       string
    TypeParams []TypeParameter
    BaseType   Type
}

type TypeParameter struct {
    Name string
    Constraint Type
}

// CheckGenericFunction 检查泛型函数
func (tc *GoTypeChecker) CheckGenericFunction(fun *ast.FuncDecl) error {
    // 解析类型参数
    typeParams, err := tc.parseTypeParameters(fun.Type.TypeParams)
    if err != nil {
        return err
    }
    
    // 创建泛型环境
    genericEnv := tc.createGenericEnvironment(typeParams)
    
    // 检查函数体
    return tc.checkFunctionBody(fun.Body, genericEnv)
}

// parseTypeParameters 解析类型参数
func (tc *GoTypeChecker) parseTypeParameters(typeParams *ast.FieldList) ([]TypeParameter, error) {
    if typeParams == nil {
        return nil, nil
    }
    
    params := make([]TypeParameter, 0)
    for _, field := range typeParams.List {
        for _, name := range field.Names {
            var constraint Type
            if field.Type != nil {
                var err error
                constraint, err = tc.CheckExpression(field.Type)
                if err != nil {
                    return nil, err
                }
            }
            
            params = append(params, TypeParameter{
                Name:      name.Name,
                Constraint: constraint,
            })
        }
    }
    
    return params, nil
}

// createGenericEnvironment 创建泛型环境
func (tc *GoTypeChecker) createGenericEnvironment(typeParams []TypeParameter) TypeEnvironment {
    env := tc.env
    
    for _, param := range typeParams {
        env.types[param.Name] = &GenericType{
            Name:       param.Name,
            TypeParams: []TypeParameter{param},
        }
    }
    
    return env
}

// InstantiateGeneric 实例化泛型
func (tc *GoTypeChecker) InstantiateGeneric(genericType *GenericType, typeArgs []Type) (Type, error) {
    if len(typeArgs) != len(genericType.TypeParams) {
        return nil, fmt.Errorf("wrong number of type arguments")
    }
    
    // 检查类型参数约束
    for i, arg := range typeArgs {
        param := genericType.TypeParams[i]
        if param.Constraint != nil {
            if !tc.satisfiesConstraint(arg, param.Constraint) {
                return nil, fmt.Errorf("type argument %v does not satisfy constraint %v", arg, param.Constraint)
            }
        }
    }
    
    // 替换类型参数
    return tc.substituteType(genericType.BaseType, genericType.TypeParams, typeArgs)
}

// substituteType 替换类型中的类型参数
func (tc *GoTypeChecker) substituteType(typ Type, params []TypeParameter, args []Type) (Type, error) {
    // 创建替换映射
    substitution := make(map[string]Type)
    for i, param := range params {
        substitution[param.Name] = args[i]
    }
    
    return tc.applySubstitution(typ, substitution)
}

// applySubstitution 应用替换
func (tc *GoTypeChecker) applySubstitution(typ Type, substitution map[string]Type) (Type, error) {
    switch t := typ.(type) {
    case *GenericType:
        if replacement, exists := substitution[t.Name]; exists {
            return replacement, nil
        }
        return t, nil
    case *FunctionType:
        newParams := make([]Type, len(t.Params))
        for i, param := range t.Params {
            newParam, err := tc.applySubstitution(param, substitution)
            if err != nil {
                return nil, err
            }
            newParams[i] = newParam
        }
        
        newResult, err := tc.applySubstitution(t.Result, substitution)
        if err != nil {
            return nil, err
        }
        
        return &FunctionType{
            Params: newParams,
            Result: newResult,
        }, nil
    default:
        return t, nil
    }
}
```

## 总结

类型系统理论是编程语言理论的核心，通过形式化定义和Go语言实现，我们建立了从理论到实践的完整框架。

### 关键要点

1. **理论基础**: 类型定义、类型关系、类型推导
2. **核心算法**: Hindley-Milner算法、统一算法、约束求解
3. **类型安全**: 进展定理、保持定理、错误处理
4. **高级特性**: 多态类型、依赖类型、线性类型

### 进一步研究方向

1. **类型系统设计**: 新类型系统、类型系统扩展
2. **类型推导优化**: 高效算法、增量推导
3. **类型安全证明**: 形式化证明、自动化验证
4. **实际应用**: 编译器实现、IDE支持

---

**激情澎湃的持续构建** <(￣︶￣)↗[GO!]
