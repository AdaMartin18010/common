# 01-类型系统理论 (Type System Theory)

## 目录

- [01-类型系统理论 (Type System Theory)](#01-类型系统理论-type-system-theory)
  - [目录](#目录)
  - [概述](#概述)
  - [1. 类型基础 (Type Foundations)](#1-类型基础-type-foundations)
    - [1.1 类型定义](#11-类型定义)
    - [1.2 基本类型](#12-基本类型)
    - [1.3 类型关系](#13-类型关系)
  - [2. 类型推导 (Type Inference)](#2-类型推导-type-inference)
    - [2.1 类型推导规则](#21-类型推导规则)
    - [2.2 统一算法](#22-统一算法)
    - [2.3 Hindley-Milner类型系统](#23-hindley-milner类型系统)
  - [3. 类型安全 (Type Safety)](#3-类型安全-type-safety)
    - [3.1 类型安全定义](#31-类型安全定义)
    - [3.2 类型安全证明](#32-类型安全证明)
    - [3.3 运行时类型检查](#33-运行时类型检查)
  - [4. 高级类型系统 (Advanced Type Systems)](#4-高级类型系统-advanced-type-systems)
    - [4.1 依赖类型](#41-依赖类型)
    - [4.2 高阶类型](#42-高阶类型)
    - [4.3 类型族](#43-类型族)
  - [5. Go语言类型系统](#5-go语言类型系统)
    - [5.1 Go类型系统特征](#51-go类型系统特征)
    - [5.2 接口类型](#52-接口类型)
    - [5.3 泛型实现](#53-泛型实现)
  - [6. 形式化验证](#6-形式化验证)
    - [6.1 类型检查算法](#61-类型检查算法)
    - [6.2 类型安全证明](#62-类型安全证明)
  - [7. 性能分析](#7-性能分析)
    - [7.1 类型检查复杂度](#71-类型检查复杂度)
    - [7.2 类型推导复杂度](#72-类型推导复杂度)
    - [7.3 性能优化](#73-性能优化)
  - [8. 应用实例](#8-应用实例)
    - [8.1 类型安全的计算器](#81-类型安全的计算器)
    - [8.2 类型安全的数据库查询](#82-类型安全的数据库查询)
  - [参考文献](#参考文献)

## 概述

类型系统理论是编程语言理论的核心组成部分，研究类型的概念、类型检查算法、类型安全性和类型推导等问题。本章节从形式化角度分析类型系统，并结合Go语言的类型系统进行实践。

### 核心概念

- **类型**: 值的集合和操作的规范
- **类型检查**: 验证程序类型正确性的过程
- **类型推导**: 自动推断表达式类型的过程
- **类型安全**: 防止类型错误的保证

## 1. 类型基础 (Type Foundations)

### 1.1 类型定义

**定义 1.1** (类型): 类型 $T$ 是值的集合 $V_T$ 和操作集合 $O_T$ 的二元组：
```latex
T = (V_T, O_T)
```

其中：
- $V_T$ 是类型 $T$ 的值域
- $O_T$ 是类型 $T$ 支持的操作集合

### 1.2 基本类型

**定义 1.2** (基本类型): 基本类型是语言预定义的类型，包括：

- **布尔类型**: $\text{Bool} = (\{\text{true}, \text{false}\}, \{\land, \lor, \neg\})$
- **整数类型**: $\text{Int} = (\mathbb{Z}, \{+, -, \times, \div, \mod\})$
- **浮点类型**: $\text{Float} = (\mathbb{R}, \{+, -, \times, \div\})$
- **字符串类型**: $\text{String} = (\Sigma^*, \{\text{concat}, \text{length}, \text{substring}\})$

### 1.3 类型关系

**定义 1.3** (子类型关系): 类型 $S$ 是类型 $T$ 的子类型，记作 $S \leq T$，如果：
```latex
V_S \subseteq V_T \land O_S \supseteq O_T
```

**定义 1.4** (类型等价): 类型 $S$ 和 $T$ 等价，记作 $S \equiv T$，如果：
```latex
S \leq T \land T \leq S
```

### 1.4 类型环境

**定义 1.5** (类型环境): 类型环境 $\Gamma$ 是从变量到类型的映射：
```latex
\Gamma: \text{Var} \rightarrow \text{Type}
```

**定义 1.6** (类型判断): 类型判断的形式为：
```latex
\Gamma \vdash e: T
```

表示在环境 $\Gamma$ 下，表达式 $e$ 具有类型 $T$。

## 2. 类型推导 (Type Inference)

### 2.1 类型推导规则

**规则 2.1** (变量规则):
```latex
\frac{x: T \in \Gamma}{\Gamma \vdash x: T}
```

**规则 2.2** (函数应用规则):
```latex
\frac{\Gamma \vdash e_1: T_1 \rightarrow T_2 \quad \Gamma \vdash e_2: T_1}{\Gamma \vdash e_1(e_2): T_2}
```

**规则 2.3** (函数抽象规则):
```latex
\frac{\Gamma, x: T_1 \vdash e: T_2}{\Gamma \vdash \lambda x: T_1.e: T_1 \rightarrow T_2}
```

### 2.2 统一算法

**定义 2.1** (类型方程): 类型方程的形式为 $T_1 = T_2$，其中 $T_1$ 和 $T_2$ 是类型表达式。

**算法 2.1** (Robinson统一算法):

```go
// 类型表达式
type TypeExpr interface {
    Unify(other TypeExpr) (Substitution, error)
}

// 类型变量
type TypeVar struct {
    name string
}

// 函数类型
type FuncType struct {
    domain   TypeExpr
    codomain TypeExpr
}

// 基本类型
type BasicType struct {
    name string
}

// 替换
type Substitution map[string]TypeExpr

// 统一算法
func Unify(e1, e2 TypeExpr) (Substitution, error) {
    switch t1 := e1.(type) {
    case *TypeVar:
        return unifyVar(t1, e2)
    case *BasicType:
        if t2, ok := e2.(*BasicType); ok {
            if t1.name == t2.name {
                return Substitution{}, nil
            }
            return nil, fmt.Errorf("type mismatch: %s != %s", t1.name, t2.name)
        }
        return nil, fmt.Errorf("cannot unify basic type with %T", e2)
    case *FuncType:
        if t2, ok := e2.(*FuncType); ok {
            s1, err := Unify(t1.domain, t2.domain)
            if err != nil {
                return nil, err
            }
            
            s2, err := Unify(applySubstitution(t1.codomain, s1), 
                           applySubstitution(t2.codomain, s1))
            if err != nil {
                return nil, err
            }
            
            return composeSubstitutions(s1, s2), nil
        }
        return nil, fmt.Errorf("cannot unify function type with %T", e2)
    default:
        return nil, fmt.Errorf("unknown type expression: %T", e1)
    }
}

// 统一类型变量
func unifyVar(v *TypeVar, t TypeExpr) (Substitution, error) {
    if v2, ok := t.(*TypeVar); ok && v.name == v2.name {
        return Substitution{}, nil
    }
    
    if occursIn(v, t) {
        return nil, fmt.Errorf("circular type: %s occurs in %v", v.name, t)
    }
    
    return Substitution{v.name: t}, nil
}

// 检查类型变量是否出现在类型中
func occursIn(v *TypeVar, t TypeExpr) bool {
    switch t2 := t.(type) {
    case *TypeVar:
        return v.name == t2.name
    case *FuncType:
        return occursIn(v, t2.domain) || occursIn(v, t2.codomain)
    default:
        return false
    }
}
```

### 2.3 Hindley-Milner类型系统

**定义 2.2** (多态类型): 多态类型的形式为 $\forall \alpha. T$，其中 $\alpha$ 是类型变量。

**算法 2.2** (Hindley-Milner类型推导):

```go
// 多态类型
type PolyType struct {
    vars []string
    body TypeExpr
}

// 类型推导器
type TypeInferrer struct {
    env     map[string]TypeExpr
    counter int
}

func NewTypeInferrer() *TypeInferrer {
    return &TypeInferrer{
        env:     make(map[string]TypeExpr),
        counter: 0,
    }
}

// 生成新的类型变量
func (ti *TypeInferrer) freshVar() *TypeVar {
    ti.counter++
    return &TypeVar{name: fmt.Sprintf("α%d", ti.counter)}
}

// 推导表达式类型
func (ti *TypeInferrer) Infer(expr Expr) (TypeExpr, error) {
    switch e := expr.(type) {
    case *VarExpr:
        if t, ok := ti.env[e.name]; ok {
            return ti.instantiate(t), nil
        }
        return nil, fmt.Errorf("undefined variable: %s", e.name)
        
    case *AppExpr:
        t1, err := ti.Infer(e.func)
        if err != nil {
            return nil, err
        }
        
        t2, err := ti.Infer(e.arg)
        if err != nil {
            return nil, err
        }
        
        resultType := ti.freshVar()
        funcType := &FuncType{domain: t2, codomain: resultType}
        
        substitution, err := Unify(t1, funcType)
        if err != nil {
            return nil, err
        }
        
        return applySubstitution(resultType, substitution), nil
        
    case *LambdaExpr:
        paramType := ti.freshVar()
        oldEnv := ti.env
        ti.env[e.param] = paramType
        
        bodyType, err := ti.Infer(e.body)
        if err != nil {
            ti.env = oldEnv
            return nil, err
        }
        
        ti.env = oldEnv
        return &FuncType{domain: paramType, codomain: bodyType}, nil
        
    default:
        return nil, fmt.Errorf("unknown expression type: %T", expr)
    }
}

// 实例化多态类型
func (ti *TypeInferrer) instantiate(polyType TypeExpr) TypeExpr {
    if poly, ok := polyType.(*PolyType); ok {
        substitution := make(Substitution)
        for _, varName := range poly.vars {
            substitution[varName] = ti.freshVar()
        }
        return applySubstitution(poly.body, substitution)
    }
    return polyType
}
```

## 3. 类型安全 (Type Safety)

### 3.1 类型安全定义

**定义 3.1** (类型安全): 语言是类型安全的，如果所有类型正确的程序都不会产生运行时类型错误。

**定理 3.1** (进展定理): 如果 $\vdash e: T$ 且 $e$ 不是值，则存在 $e'$ 使得 $e \rightarrow e'$。

**定理 3.2** (保持定理): 如果 $\vdash e: T$ 且 $e \rightarrow e'$，则 $\vdash e': T$。

### 3.2 类型安全证明

**证明 3.1** (进展定理证明):

1. 对表达式 $e$ 的结构进行归纳
2. 对于每种表达式形式，证明要么是值，要么可以继续求值
3. 利用类型推导规则确保类型一致性

**证明 3.2** (保持定理证明):

1. 对求值规则进行归纳
2. 证明每个求值步骤保持类型
3. 利用类型推导规则验证类型保持

### 3.3 运行时类型检查

```go
// 运行时类型检查器
type RuntimeTypeChecker struct {
    typeMap map[interface{}]reflect.Type
}

func NewRuntimeTypeChecker() *RuntimeTypeChecker {
    return &RuntimeTypeChecker{
        typeMap: make(map[interface{}]reflect.Type),
    }
}

// 检查类型
func (rtc *RuntimeTypeChecker) CheckType(value interface{}, expectedType reflect.Type) error {
    actualType := reflect.TypeOf(value)
    
    if actualType != expectedType {
        return fmt.Errorf("type mismatch: expected %v, got %v", expectedType, actualType)
    }
    
    return nil
}

// 类型安全的函数调用
func (rtc *RuntimeTypeChecker) SafeCall(fn interface{}, args ...interface{}) ([]interface{}, error) {
    fnValue := reflect.ValueOf(fn)
    fnType := fnValue.Type()
    
    if fnType.Kind() != reflect.Func {
        return nil, fmt.Errorf("not a function: %v", fnType)
    }
    
    // 检查参数数量
    if fnType.NumIn() != len(args) {
        return nil, fmt.Errorf("argument count mismatch: expected %d, got %d", 
            fnType.NumIn(), len(args))
    }
    
    // 检查参数类型
    for i := 0; i < fnType.NumIn(); i++ {
        if err := rtc.CheckType(args[i], fnType.In(i)); err != nil {
            return nil, fmt.Errorf("argument %d: %v", i, err)
        }
    }
    
    // 调用函数
    argValues := make([]reflect.Value, len(args))
    for i, arg := range args {
        argValues[i] = reflect.ValueOf(arg)
    }
    
    results := fnValue.Call(argValues)
    
    // 转换结果
    resultValues := make([]interface{}, len(results))
    for i, result := range results {
        resultValues[i] = result.Interface()
    }
    
    return resultValues, nil
}
```

## 4. 高级类型系统 (Advanced Type Systems)

### 4.1 依赖类型

**定义 4.1** (依赖类型): 依赖类型是依赖于值的类型，形式为 $\Pi x: A. B(x)$。

```go
// 依赖类型系统
type DependentType struct {
    paramType TypeExpr
    bodyType  func(Value) TypeExpr
}

// 向量类型（长度依赖类型）
type VectorType struct {
    elementType TypeExpr
    length      int
}

// 类型安全的向量
type Vector[T any] struct {
    elements []T
    length   int
}

func NewVector[T any](length int) *Vector[T] {
    return &Vector[T]{
        elements: make([]T, length),
        length:   length,
    }
}

func (v *Vector[T]) Get(index int) T {
    if index < 0 || index >= v.length {
        panic("index out of bounds")
    }
    return v.elements[index]
}

func (v *Vector[T]) Set(index int, value T) {
    if index < 0 || index >= v.length {
        panic("index out of bounds")
    }
    v.elements[index] = value
}

// 类型安全的向量连接
func ConcatVectors[T any](v1, v2 *Vector[T]) *Vector[T] {
    newLength := v1.length + v2.length
    result := NewVector[T](newLength)
    
    for i := 0; i < v1.length; i++ {
        result.Set(i, v1.Get(i))
    }
    
    for i := 0; i < v2.length; i++ {
        result.Set(v1.length+i, v2.Get(i))
    }
    
    return result
}
```

### 4.2 高阶类型

**定义 4.2** (高阶类型): 高阶类型是接受类型作为参数的类型构造器。

```go
// 类型构造器
type TypeConstructor interface {
    Apply(args []TypeExpr) TypeExpr
}

// 函子类型
type Functor[T any] interface {
    Map[U any](f func(T) U) Functor[U]
}

// 单子类型
type Monad[T any] interface {
    Functor[T]
    Bind[U any](f func(T) Monad[U]) Monad[U]
    Return(value T) Monad[T]
}

// Maybe类型实现
type Maybe[T any] struct {
    value *T
}

func Just[T any](value T) Maybe[T] {
    return Maybe[T]{value: &value}
}

func Nothing[T any]() Maybe[T] {
    return Maybe[T]{value: nil}
}

func (m Maybe[T]) IsJust() bool {
    return m.value != nil
}

func (m Maybe[T]) IsNothing() bool {
    return m.value == nil
}

func (m Maybe[T]) FromJust() T {
    if m.value == nil {
        panic("fromJust: Nothing")
    }
    return *m.value
}

// Functor实现
func (m Maybe[T]) Map[U any](f func(T) U) Maybe[U] {
    if m.IsNothing() {
        return Nothing[U]()
    }
    return Just(f(m.FromJust()))
}

// Monad实现
func (m Maybe[T]) Bind[U any](f func(T) Maybe[U]) Maybe[U] {
    if m.IsNothing() {
        return Nothing[U]()
    }
    return f(m.FromJust())
}

func (m Maybe[T]) Return(value T) Maybe[T] {
    return Just(value)
}
```

### 4.3 类型族

**定义 4.3** (类型族): 类型族是相关类型的集合，通过类型函数定义。

```go
// 类型族定义
type TypeFamily interface {
    Instance(args []TypeExpr) TypeExpr
}

// 数字类型族
type NumberType interface {
    Add(other NumberType) NumberType
    Multiply(other NumberType) NumberType
    Zero() NumberType
}

// 整数类型
type IntType int

func (i IntType) Add(other NumberType) NumberType {
    if o, ok := other.(IntType); ok {
        return IntType(int(i) + int(o))
    }
    panic("type mismatch")
}

func (i IntType) Multiply(other NumberType) NumberType {
    if o, ok := other.(IntType); ok {
        return IntType(int(i) * int(o))
    }
    panic("type mismatch")
}

func (i IntType) Zero() NumberType {
    return IntType(0)
}

// 浮点类型
type FloatType float64

func (f FloatType) Add(other NumberType) NumberType {
    if o, ok := other.(FloatType); ok {
        return FloatType(float64(f) + float64(o))
    }
    panic("type mismatch")
}

func (f FloatType) Multiply(other NumberType) NumberType {
    if o, ok := other.(FloatType); ok {
        return FloatType(float64(f) * float64(o))
    }
    panic("type mismatch")
}

func (f FloatType) Zero() NumberType {
    return FloatType(0.0)
}

// 泛型数字运算
func Sum[T NumberType](values []T) T {
    result := values[0].Zero().(T)
    for _, value := range values {
        result = result.Add(value).(T)
    }
    return result
}
```

## 5. Go语言类型系统

### 5.1 Go类型系统特征

**特征 5.1** (Go类型系统):
- 静态类型系统
- 结构类型系统
- 接口类型系统
- 泛型支持（Go 1.18+）

### 5.2 接口类型

```go
// 接口定义
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 结构体实现
type Circle struct {
    radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.radius
}

type Rectangle struct {
    width, height float64
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.width + r.height)
}

// 接口使用
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %f, Perimeter: %f\n", s.Area(), s.Perimeter())
}
```

### 5.3 泛型实现

```go
// 泛型容器
type Container[T any] struct {
    elements []T
}

func NewContainer[T any]() *Container[T] {
    return &Container[T]{
        elements: make([]T, 0),
    }
}

func (c *Container[T]) Add(element T) {
    c.elements = append(c.elements, element)
}

func (c *Container[T]) Get(index int) T {
    if index < 0 || index >= len(c.elements) {
        panic("index out of bounds")
    }
    return c.elements[index]
}

func (c *Container[T]) Size() int {
    return len(c.elements)
}

// 泛型算法
func Map[T, U any](elements []T, f func(T) U) []U {
    result := make([]U, len(elements))
    for i, element := range elements {
        result[i] = f(element)
    }
    return result
}

func Filter[T any](elements []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, element := range elements {
        if predicate(element) {
            result = append(result, element)
        }
    }
    return result
}

func Reduce[T, U any](elements []T, initial U, f func(U, T) U) U {
    result := initial
    for _, element := range elements {
        result = f(result, element)
    }
    return result
}
```

## 6. 形式化验证

### 6.1 类型检查算法

```go
// 类型检查器
type TypeChecker struct {
    env map[string]TypeExpr
}

func NewTypeChecker() *TypeChecker {
    return &TypeChecker{
        env: make(map[string]TypeExpr),
    }
}

// 类型检查
func (tc *TypeChecker) Check(expr Expr) (TypeExpr, error) {
    switch e := expr.(type) {
    case *LiteralExpr:
        return tc.checkLiteral(e)
    case *VarExpr:
        return tc.checkVariable(e)
    case *BinaryExpr:
        return tc.checkBinary(e)
    case *FuncExpr:
        return tc.checkFunction(e)
    case *CallExpr:
        return tc.checkCall(e)
    default:
        return nil, fmt.Errorf("unknown expression type: %T", expr)
    }
}

// 检查字面量
func (tc *TypeChecker) checkLiteral(lit *LiteralExpr) (TypeExpr, error) {
    switch lit.value.(type) {
    case int:
        return &BasicType{name: "int"}, nil
    case float64:
        return &BasicType{name: "float"}, nil
    case string:
        return &BasicType{name: "string"}, nil
    case bool:
        return &BasicType{name: "bool"}, nil
    default:
        return nil, fmt.Errorf("unknown literal type: %T", lit.value)
    }
}

// 检查变量
func (tc *TypeChecker) checkVariable(varExpr *VarExpr) (TypeExpr, error) {
    if t, ok := tc.env[varExpr.name]; ok {
        return t, nil
    }
    return nil, fmt.Errorf("undefined variable: %s", varExpr.name)
}

// 检查二元表达式
func (tc *TypeChecker) checkBinary(bin *BinaryExpr) (TypeExpr, error) {
    leftType, err := tc.Check(bin.left)
    if err != nil {
        return nil, err
    }
    
    rightType, err := tc.Check(bin.right)
    if err != nil {
        return nil, err
    }
    
    // 检查类型兼容性
    if !tc.isCompatible(leftType, rightType) {
        return nil, fmt.Errorf("type mismatch: %v %s %v", leftType, bin.operator, rightType)
    }
    
    // 返回结果类型
    return tc.getResultType(bin.operator, leftType), nil
}
```

### 6.2 类型安全证明

```go
// 类型安全证明器
type TypeSafetyProver struct {
    checker *TypeChecker
}

func NewTypeSafetyProver() *TypeSafetyProver {
    return &TypeSafetyProver{
        checker: NewTypeChecker(),
    }
}

// 证明类型安全
func (tsp *TypeSafetyProver) ProveSafety(expr Expr) error {
    // 1. 类型检查
    _, err := tsp.checker.Check(expr)
    if err != nil {
        return fmt.Errorf("type check failed: %v", err)
    }
    
    // 2. 证明进展定理
    if err := tsp.proveProgress(expr); err != nil {
        return fmt.Errorf("progress theorem failed: %v", err)
    }
    
    // 3. 证明保持定理
    if err := tsp.provePreservation(expr); err != nil {
        return fmt.Errorf("preservation theorem failed: %v", err)
    }
    
    return nil
}

// 证明进展定理
func (tsp *TypeSafetyProver) proveProgress(expr Expr) error {
    // 检查表达式是否可以继续求值
    if tsp.isValue(expr) {
        return nil // 已经是值
    }
    
    // 检查是否可以应用求值规则
    if tsp.canStep(expr) {
        return nil // 可以继续求值
    }
    
    return fmt.Errorf("expression cannot progress: %v", expr)
}

// 证明保持定理
func (tsp *TypeSafetyProver) provePreservation(expr Expr) error {
    // 获取原始类型
    originalType, err := tsp.checker.Check(expr)
    if err != nil {
        return err
    }
    
    // 模拟一步求值
    if nextExpr, err := tsp.step(expr); err == nil {
        // 检查求值后的类型
        newType, err := tsp.checker.Check(nextExpr)
        if err != nil {
            return err
        }
        
        // 验证类型保持
        if !tsp.typesEqual(originalType, newType) {
            return fmt.Errorf("type not preserved: %v -> %v", originalType, newType)
        }
    }
    
    return nil
}
```

## 7. 性能分析

### 7.1 类型检查复杂度

**定理 7.1**: 简单类型检查的时间复杂度为 $O(n)$，其中 $n$ 是表达式的大小。

**证明**:
1. 每个节点最多被访问一次
2. 每个节点的类型检查操作是常数时间
3. 总体时间复杂度为 $O(n)$

### 7.2 类型推导复杂度

**定理 7.2**: Hindley-Milner类型推导的时间复杂度为 $O(n^3)$。

**证明**:
1. 统一算法的时间复杂度为 $O(n^2)$
2. 每个节点可能需要统一操作
3. 总体时间复杂度为 $O(n^3)$

### 7.3 性能优化

```go
// 缓存类型检查器
type CachedTypeChecker struct {
    checker *TypeChecker
    cache   map[string]TypeExpr
    mu      sync.RWMutex
}

func NewCachedTypeChecker() *CachedTypeChecker {
    return &CachedTypeChecker{
        checker: NewTypeChecker(),
        cache:   make(map[string]TypeExpr),
    }
}

func (ctc *CachedTypeChecker) Check(expr Expr) (TypeExpr, error) {
    // 生成缓存键
    key := ctc.generateKey(expr)
    
    // 检查缓存
    ctc.mu.RLock()
    if cached, ok := ctc.cache[key]; ok {
        ctc.mu.RUnlock()
        return cached, nil
    }
    ctc.mu.RUnlock()
    
    // 执行类型检查
    result, err := ctc.checker.Check(expr)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果
    ctc.mu.Lock()
    ctc.cache[key] = result
    ctc.mu.Unlock()
    
    return result, nil
}

func (ctc *CachedTypeChecker) generateKey(expr Expr) string {
    // 简化的键生成算法
    return fmt.Sprintf("%T-%v", expr, expr)
}
```

## 8. 应用实例

### 8.1 类型安全的计算器

```go
// 类型安全的计算器
type Calculator struct {
    checker *TypeChecker
}

func NewCalculator() *Calculator {
    return &Calculator{
        checker: NewTypeChecker(),
    }
}

// 表达式类型
type Expr interface {
    Eval() interface{}
}

type NumberExpr struct {
    value float64
}

func (n NumberExpr) Eval() interface{} {
    return n.value
}

type AddExpr struct {
    left, right Expr
}

func (a AddExpr) Eval() interface{} {
    left := a.left.Eval().(float64)
    right := a.right.Eval().(float64)
    return left + right
}

type MulExpr struct {
    left, right Expr
}

func (m MulExpr) Eval() interface{} {
    left := m.left.Eval().(float64)
    right := m.right.Eval().(float64)
    return left * right
}

// 类型安全的求值
func (c *Calculator) SafeEval(expr Expr) (interface{}, error) {
    // 类型检查
    _, err := c.checker.Check(expr)
    if err != nil {
        return nil, err
    }
    
    // 安全求值
    return expr.Eval(), nil
}
```

### 8.2 类型安全的数据库查询

```go
// 类型安全的查询构建器
type QueryBuilder[T any] struct {
    table   string
    fields  []string
    where   []Condition
    orderBy []OrderBy
    limit   *int
    offset  *int
}

type Condition struct {
    field    string
    operator string
    value    interface{}
}

type OrderBy struct {
    field     string
    direction string
}

func NewQueryBuilder[T any](table string) *QueryBuilder[T] {
    return &QueryBuilder[T]{
        table:  table,
        fields: make([]string, 0),
        where:  make([]Condition, 0),
        orderBy: make([]OrderBy, 0),
    }
}

func (qb *QueryBuilder[T]) Select(fields ...string) *QueryBuilder[T] {
    qb.fields = append(qb.fields, fields...)
    return qb
}

func (qb *QueryBuilder[T]) Where(field, operator string, value interface{}) *QueryBuilder[T] {
    qb.where = append(qb.where, Condition{
        field:    field,
        operator: operator,
        value:    value,
    })
    return qb
}

func (qb *QueryBuilder[T]) OrderBy(field, direction string) *QueryBuilder[T] {
    qb.orderBy = append(qb.orderBy, OrderBy{
        field:     field,
        direction: direction,
    })
    return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
    qb.limit = &limit
    return qb
}

func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
    qb.offset = &offset
    return qb
}

func (qb *QueryBuilder[T]) Build() (string, []interface{}, error) {
    // 构建SQL查询
    query := "SELECT "
    
    if len(qb.fields) == 0 {
        query += "*"
    } else {
        query += strings.Join(qb.fields, ", ")
    }
    
    query += " FROM " + qb.table
    
    args := make([]interface{}, 0)
    
    if len(qb.where) > 0 {
        query += " WHERE "
        conditions := make([]string, 0)
        for _, condition := range qb.where {
            conditions = append(conditions, fmt.Sprintf("%s %s ?", condition.field, condition.operator))
            args = append(args, condition.value)
        }
        query += strings.Join(conditions, " AND ")
    }
    
    if len(qb.orderBy) > 0 {
        query += " ORDER BY "
        orders := make([]string, 0)
        for _, order := range qb.orderBy {
            orders = append(orders, fmt.Sprintf("%s %s", order.field, order.direction))
        }
        query += strings.Join(orders, ", ")
    }
    
    if qb.limit != nil {
        query += fmt.Sprintf(" LIMIT %d", *qb.limit)
    }
    
    if qb.offset != nil {
        query += fmt.Sprintf(" OFFSET %d", *qb.offset)
    }
    
    return query, args, nil
}
```

## 参考文献

1. Pierce, B. C. (2002). *Types and Programming Languages*. MIT Press.
2. Cardelli, L., & Wegner, P. (1985). *On Understanding Types, Data Abstraction, and Polymorphism*. ACM Computing Surveys.
3. Milner, R. (1978). *A Theory of Type Polymorphism in Programming*. Journal of Computer and System Sciences.
4. Hindley, J. R. (1969). *The Principal Type-Scheme of an Object in Combinatory Logic*. Transactions of the American Mathematical Society.
5. Reynolds, J. C. (1974). *Towards a Theory of Type Structure*. Programming Symposium.

---

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **类型系统理论完成！** 🚀
