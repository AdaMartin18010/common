# 01-类型基础 (Type Foundations)

## 目录

- [01-类型基础 (Type Foundations)](#01-类型基础-type-foundations)
  - [目录](#目录)
  - [概述](#概述)
  - [基本概念](#基本概念)
    - [类型定义](#类型定义)
    - [类型关系](#类型关系)
    - [类型安全](#类型安全)
  - [形式化理论](#形式化理论)
    - [类型系统公理](#类型系统公理)
    - [类型推导规则](#类型推导规则)
    - [类型等价性](#类型等价性)
  - [Go语言类型系统](#go语言类型系统)
    - [基本类型](#基本类型)
    - [复合类型](#复合类型)
    - [接口类型](#接口类型)
    - [泛型类型](#泛型类型)
  - [类型推导算法](#类型推导算法)
    - [Hindley-Milner算法](#hindley-milner算法)
    - [Go语言类型推导](#go语言类型推导)
    - [类型约束求解](#类型约束求解)
  - [类型安全证明](#类型安全证明)
    - [进展定理](#进展定理)
    - [保持定理](#保持定理)
    - [类型健全性](#类型健全性)
  - [应用领域](#应用领域)
    - [编译器设计](#编译器设计)
    - [静态分析](#静态分析)
    - [程序验证](#程序验证)
    - [代码生成](#代码生成)
  - [相关链接](#相关链接)

## 概述

类型基础是编程语言理论的核心组成部分，研究类型系统的数学基础和形式化定义。类型系统为程序提供了抽象层次，确保程序的正确性和安全性。本模块基于 `/docs/model` 目录中的编程语言理论，结合 Go 语言的类型系统特性，建立完整的类型基础理论体系。

## 基本概念

### 类型定义

**定义 1 (类型)**: 类型 $T$ 是值的集合，定义了值的结构和操作。

**定义 2 (类型环境)**: 类型环境 $\Gamma$ 是从变量到类型的映射：
$$\Gamma: \text{Var} \rightarrow \text{Type}$$

**定义 3 (类型判断)**: 类型判断 $\Gamma \vdash e: T$ 表示在环境 $\Gamma$ 下，表达式 $e$ 具有类型 $T$。

**定义 4 (类型上下文)**: 类型上下文包含类型变量和约束：
$$\Delta = \{\alpha_1, \alpha_2, \ldots, \alpha_n\}$$

### 类型关系

**定义 5 (子类型关系)**: 类型 $S$ 是类型 $T$ 的子类型，记为 $S \leq T$，如果 $S$ 的值可以安全地用在期望 $T$ 的地方。

**定义 6 (类型等价)**: 类型 $S$ 和 $T$ 等价，记为 $S \equiv T$，如果 $S \leq T$ 且 $T \leq S$。

**定义 7 (类型包含)**: 类型 $S$ 包含类型 $T$，记为 $S \supseteq T$，如果 $T$ 的所有值都是 $S$ 的值。

### 类型安全

**定义 8 (类型安全)**: 程序是类型安全的，如果所有表达式都有正确的类型，且类型检查通过。

**定义 9 (类型错误)**: 类型错误是违反类型系统规则的错误。

**定义 10 (类型检查)**: 类型检查是验证程序类型安全性的过程。

## 形式化理论

### 类型系统公理

**公理 1 (变量规则)**: 如果 $x: T \in \Gamma$，则 $\Gamma \vdash x: T$

**公理 2 (函数应用)**: 如果 $\Gamma \vdash f: T_1 \rightarrow T_2$ 且 $\Gamma \vdash x: T_1$，则 $\Gamma \vdash f(x): T_2$

**公理 3 (函数抽象)**: 如果 $\Gamma, x: T_1 \vdash e: T_2$，则 $\Gamma \vdash \lambda x: T_1.e: T_1 \rightarrow T_2$

**公理 4 (子类型)**: 如果 $\Gamma \vdash e: S$ 且 $S \leq T$，则 $\Gamma \vdash e: T$

### 类型推导规则

**规则 1 (类型推导)**: 类型推导算法 $\mathcal{W}$ 计算表达式 $e$ 的类型：
$$\mathcal{W}(\Gamma, e) = (S, \theta)$$
其中 $S$ 是类型，$\theta$ 是替换。

**规则 2 (统一)**: 类型统一算法 $\mathcal{U}$ 求解类型方程：
$$\mathcal{U}(T_1, T_2) = \theta$$
其中 $\theta$ 是最一般统一子。

**规则 3 (泛化)**: 泛化算法 $\mathcal{G}$ 计算多态类型：
$$\mathcal{G}(\Gamma, T) = \forall \alpha_1, \alpha_2, \ldots, \alpha_n.T$$

### 类型等价性

**定理 1 (类型等价的自反性)**: $\forall T: T \equiv T$

**定理 2 (类型等价的对称性)**: 如果 $S \equiv T$，则 $T \equiv S$

**定理 3 (类型等价的传递性)**: 如果 $S \equiv T$ 且 $T \equiv U$，则 $S \equiv U$

**定理 4 (函数类型等价)**: $(S_1 \rightarrow T_1) \equiv (S_2 \rightarrow T_2)$ 当且仅当 $S_1 \equiv S_2$ 且 $T_1 \equiv T_2$

## Go语言类型系统

### 基本类型

```go
package typesystem

import (
    "fmt"
    "reflect"
)

// Type 类型接口
type Type interface {
    String() string
    IsAssignableTo(other Type) bool
    IsEquivalentTo(other Type) bool
    GetKind() reflect.Kind
}

// BasicType 基本类型
type BasicType struct {
    Name string
    Kind reflect.Kind
}

// NewBasicType 创建基本类型
func NewBasicType(name string, kind reflect.Kind) *BasicType {
    return &BasicType{
        Name: name,
        Kind: kind,
    }
}

// String 字符串表示
func (bt *BasicType) String() string {
    return bt.Name
}

// IsAssignableTo 检查是否可赋值
func (bt *BasicType) IsAssignableTo(other Type) bool {
    if otherBT, ok := other.(*BasicType); ok {
        return bt.Kind == otherBT.Kind
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (bt *BasicType) IsEquivalentTo(other Type) bool {
    return bt.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (bt *BasicType) GetKind() reflect.Kind {
    return bt.Kind
}

// 预定义基本类型
var (
    IntType    = NewBasicType("int", reflect.Int)
    FloatType  = NewBasicType("float64", reflect.Float64)
    StringType = NewBasicType("string", reflect.String)
    BoolType   = NewBasicType("bool", reflect.Bool)
)
```

### 复合类型

```go
// ArrayType 数组类型
type ArrayType struct {
    ElementType Type
    Length      int
}

// NewArrayType 创建数组类型
func NewArrayType(elementType Type, length int) *ArrayType {
    return &ArrayType{
        ElementType: elementType,
        Length:      length,
    }
}

// String 字符串表示
func (at *ArrayType) String() string {
    return fmt.Sprintf("[%d]%s", at.Length, at.ElementType.String())
}

// IsAssignableTo 检查是否可赋值
func (at *ArrayType) IsAssignableTo(other Type) bool {
    if otherAT, ok := other.(*ArrayType); ok {
        return at.Length == otherAT.Length && at.ElementType.IsAssignableTo(otherAT.ElementType)
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (at *ArrayType) IsEquivalentTo(other Type) bool {
    return at.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (at *ArrayType) GetKind() reflect.Kind {
    return reflect.Array
}

// SliceType 切片类型
type SliceType struct {
    ElementType Type
}

// NewSliceType 创建切片类型
func NewSliceType(elementType Type) *SliceType {
    return &SliceType{
        ElementType: elementType,
    }
}

// String 字符串表示
func (st *SliceType) String() string {
    return fmt.Sprintf("[]%s", st.ElementType.String())
}

// IsAssignableTo 检查是否可赋值
func (st *SliceType) IsAssignableTo(other Type) bool {
    if otherST, ok := other.(*SliceType); ok {
        return st.ElementType.IsAssignableTo(otherST.ElementType)
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (st *SliceType) IsEquivalentTo(other Type) bool {
    return st.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (st *SliceType) GetKind() reflect.Kind {
    return reflect.Slice
}

// MapType 映射类型
type MapType struct {
    KeyType   Type
    ValueType Type
}

// NewMapType 创建映射类型
func NewMapType(keyType, valueType Type) *MapType {
    return &MapType{
        KeyType:   keyType,
        ValueType: valueType,
    }
}

// String 字符串表示
func (mt *MapType) String() string {
    return fmt.Sprintf("map[%s]%s", mt.KeyType.String(), mt.ValueType.String())
}

// IsAssignableTo 检查是否可赋值
func (mt *MapType) IsAssignableTo(other Type) bool {
    if otherMT, ok := other.(*MapType); ok {
        return mt.KeyType.IsAssignableTo(otherMT.KeyType) && 
               mt.ValueType.IsAssignableTo(otherMT.ValueType)
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (mt *MapType) IsEquivalentTo(other Type) bool {
    return mt.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (mt *MapType) GetKind() reflect.Kind {
    return reflect.Map
}

// StructType 结构体类型
type StructType struct {
    Fields map[string]Type
}

// NewStructType 创建结构体类型
func NewStructType(fields map[string]Type) *StructType {
    return &StructType{
        Fields: fields,
    }
}

// String 字符串表示
func (st *StructType) String() string {
    fields := make([]string, 0, len(st.Fields))
    for name, typ := range st.Fields {
        fields = append(fields, fmt.Sprintf("%s %s", name, typ.String()))
    }
    return fmt.Sprintf("struct{%s}", strings.Join(fields, "; "))
}

// IsAssignableTo 检查是否可赋值
func (st *StructType) IsAssignableTo(other Type) bool {
    if otherST, ok := other.(*StructType); ok {
        if len(st.Fields) != len(otherST.Fields) {
            return false
        }
        for name, typ := range st.Fields {
            if otherTyp, exists := otherST.Fields[name]; !exists || !typ.IsAssignableTo(otherTyp) {
                return false
            }
        }
        return true
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (st *StructType) IsEquivalentTo(other Type) bool {
    return st.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (st *StructType) GetKind() reflect.Kind {
    return reflect.Struct
}
```

### 接口类型

```go
// InterfaceType 接口类型
type InterfaceType struct {
    Methods map[string]*FunctionType
}

// NewInterfaceType 创建接口类型
func NewInterfaceType(methods map[string]*FunctionType) *InterfaceType {
    return &InterfaceType{
        Methods: methods,
    }
}

// String 字符串表示
func (it *InterfaceType) String() string {
    methods := make([]string, 0, len(it.Methods))
    for name, method := range it.Methods {
        methods = append(methods, fmt.Sprintf("%s%s", name, method.String()))
    }
    return fmt.Sprintf("interface{%s}", strings.Join(methods, "; "))
}

// IsAssignableTo 检查是否可赋值
func (it *InterfaceType) IsAssignableTo(other Type) bool {
    if otherIT, ok := other.(*InterfaceType); ok {
        // 检查是否实现了所有方法
        for name, method := range otherIT.Methods {
            if thisMethod, exists := it.Methods[name]; !exists || !thisMethod.IsAssignableTo(method) {
                return false
            }
        }
        return true
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (it *InterfaceType) IsEquivalentTo(other Type) bool {
    return it.IsAssignableTo(other) && other.(*InterfaceType).IsAssignableTo(it)
}

// GetKind 获取类型种类
func (it *InterfaceType) GetKind() reflect.Kind {
    return reflect.Interface
}

// FunctionType 函数类型
type FunctionType struct {
    Parameters []Type
    ReturnType Type
}

// NewFunctionType 创建函数类型
func NewFunctionType(parameters []Type, returnType Type) *FunctionType {
    return &FunctionType{
        Parameters: parameters,
        ReturnType: returnType,
    }
}

// String 字符串表示
func (ft *FunctionType) String() string {
    params := make([]string, 0, len(ft.Parameters))
    for _, param := range ft.Parameters {
        params = append(params, param.String())
    }
    return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), ft.ReturnType.String())
}

// IsAssignableTo 检查是否可赋值
func (ft *FunctionType) IsAssignableTo(other Type) bool {
    if otherFT, ok := other.(*FunctionType); ok {
        if len(ft.Parameters) != len(otherFT.Parameters) {
            return false
        }
        
        // 检查参数类型（逆变）
        for i, param := range ft.Parameters {
            if !otherFT.Parameters[i].IsAssignableTo(param) {
                return false
            }
        }
        
        // 检查返回类型（协变）
        return ft.ReturnType.IsAssignableTo(otherFT.ReturnType)
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (ft *FunctionType) IsEquivalentTo(other Type) bool {
    return ft.IsAssignableTo(other) && other.(*FunctionType).IsAssignableTo(ft)
}

// GetKind 获取类型种类
func (ft *FunctionType) GetKind() reflect.Kind {
    return reflect.Func
}
```

### 泛型类型

```go
// GenericType 泛型类型
type GenericType struct {
    Name       string
    TypeParams []string
    BaseType   Type
}

// NewGenericType 创建泛型类型
func NewGenericType(name string, typeParams []string, baseType Type) *GenericType {
    return &GenericType{
        Name:       name,
        TypeParams: typeParams,
        BaseType:   baseType,
    }
}

// String 字符串表示
func (gt *GenericType) String() string {
    if len(gt.TypeParams) > 0 {
        return fmt.Sprintf("%s[%s]", gt.Name, strings.Join(gt.TypeParams, ", "))
    }
    return gt.Name
}

// IsAssignableTo 检查是否可赋值
func (gt *GenericType) IsAssignableTo(other Type) bool {
    if otherGT, ok := other.(*GenericType); ok {
        return gt.Name == otherGT.Name && len(gt.TypeParams) == len(otherGT.TypeParams)
    }
    return false
}

// IsEquivalentTo 检查是否等价
func (gt *GenericType) IsEquivalentTo(other Type) bool {
    return gt.IsAssignableTo(other)
}

// GetKind 获取类型种类
func (gt *GenericType) GetKind() reflect.Kind {
    return reflect.Invalid // 泛型类型没有具体的Kind
}

// TypeSubstitution 类型替换
type TypeSubstitution struct {
    Mapping map[string]Type
}

// NewTypeSubstitution 创建类型替换
func NewTypeSubstitution() *TypeSubstitution {
    return &TypeSubstitution{
        Mapping: make(map[string]Type),
    }
}

// Add 添加替换
func (ts *TypeSubstitution) Add(typeParam string, typ Type) {
    ts.Mapping[typeParam] = typ
}

// Apply 应用替换
func (ts *TypeSubstitution) Apply(typ Type) Type {
    switch t := typ.(type) {
    case *GenericType:
        if replacement, exists := ts.Mapping[t.Name]; exists {
            return replacement
        }
        return t
    case *FunctionType:
        newParams := make([]Type, len(t.Parameters))
        for i, param := range t.Parameters {
            newParams[i] = ts.Apply(param)
        }
        return NewFunctionType(newParams, ts.Apply(t.ReturnType))
    case *ArrayType:
        return NewArrayType(ts.Apply(t.ElementType), t.Length)
    case *SliceType:
        return NewSliceType(ts.Apply(t.ElementType))
    case *MapType:
        return NewMapType(ts.Apply(t.KeyType), ts.Apply(t.ValueType))
    case *StructType:
        newFields := make(map[string]Type)
        for name, fieldType := range t.Fields {
            newFields[name] = ts.Apply(fieldType)
        }
        return NewStructType(newFields)
    default:
        return typ
    }
}
```

## 类型推导算法

### Hindley-Milner算法

```go
// TypeInference 类型推导
type TypeInference struct {
    typeVarCounter int
    constraints    []TypeConstraint
}

// TypeConstraint 类型约束
type TypeConstraint struct {
    Left  Type
    Right Type
}

// NewTypeInference 创建类型推导器
func NewTypeInference() *TypeInference {
    return &TypeInference{
        typeVarCounter: 0,
        constraints:    make([]TypeConstraint, 0),
    }
}

// FreshTypeVar 生成新的类型变量
func (ti *TypeInference) FreshTypeVar() *TypeVariable {
    ti.typeVarCounter++
    return NewTypeVariable(fmt.Sprintf("α%d", ti.typeVarCounter))
}

// TypeVariable 类型变量
type TypeVariable struct {
    Name string
}

// NewTypeVariable 创建类型变量
func NewTypeVariable(name string) *TypeVariable {
    return &TypeVariable{Name: name}
}

// String 字符串表示
func (tv *TypeVariable) String() string {
    return tv.Name
}

// IsAssignableTo 检查是否可赋值
func (tv *TypeVariable) IsAssignableTo(other Type) bool {
    return true // 类型变量可以赋值给任何类型
}

// IsEquivalentTo 检查是否等价
func (tv *TypeVariable) IsEquivalentTo(other Type) bool {
    if otherTV, ok := other.(*TypeVariable); ok {
        return tv.Name == otherTV.Name
    }
    return false
}

// GetKind 获取类型种类
func (tv *TypeVariable) GetKind() reflect.Kind {
    return reflect.Invalid
}

// InferType 推导类型
func (ti *TypeInference) InferType(expr Expression, env TypeEnvironment) (Type, *TypeSubstitution, error) {
    switch e := expr.(type) {
    case *Variable:
        return ti.inferVariable(e, env)
    case *Application:
        return ti.inferApplication(e, env)
    case *Abstraction:
        return ti.inferAbstraction(e, env)
    case *Literal:
        return ti.inferLiteral(e, env)
    default:
        return nil, nil, fmt.Errorf("unsupported expression type")
    }
}

// inferVariable 推导变量类型
func (ti *TypeInference) inferVariable(v *Variable, env TypeEnvironment) (Type, *TypeSubstitution, error) {
    if typ, exists := env[v.Name]; exists {
        return typ, NewTypeSubstitution(), nil
    }
    return nil, nil, fmt.Errorf("unbound variable: %s", v.Name)
}

// inferApplication 推导函数应用类型
func (ti *TypeInference) inferApplication(app *Application, env TypeEnvironment) (Type, *TypeSubstitution, error) {
    // 推导函数类型
    funType, sub1, err := ti.InferType(app.Function, env)
    if err != nil {
        return nil, nil, err
    }
    
    // 推导参数类型
    argType, sub2, err := ti.InferType(app.Argument, env)
    if err != nil {
        return nil, nil, err
    }
    
    // 生成新的类型变量作为返回类型
    returnType := ti.FreshTypeVar()
    
    // 添加约束：函数类型 = 参数类型 -> 返回类型
    ti.constraints = append(ti.constraints, TypeConstraint{
        Left:  funType,
        Right: NewFunctionType([]Type{argType}, returnType),
    })
    
    // 合并替换
    combinedSub := ti.combineSubstitutions(sub1, sub2)
    
    return returnType, combinedSub, nil
}

// inferAbstraction 推导函数抽象类型
func (ti *TypeInference) inferAbstraction(abs *Abstraction, env TypeEnvironment) (Type, *TypeSubstitution, error) {
    // 生成参数类型变量
    paramType := ti.FreshTypeVar()
    
    // 扩展环境
    newEnv := make(TypeEnvironment)
    for k, v := range env {
        newEnv[k] = v
    }
    newEnv[abs.Parameter] = paramType
    
    // 推导函数体类型
    bodyType, sub, err := ti.InferType(abs.Body, newEnv)
    if err != nil {
        return nil, nil, err
    }
    
    // 返回函数类型
    return NewFunctionType([]Type{paramType}, bodyType), sub, nil
}

// inferLiteral 推导字面量类型
func (ti *TypeInference) inferLiteral(lit *Literal, env TypeEnvironment) (Type, *TypeSubstitution, error) {
    switch lit.Value.(type) {
    case int:
        return IntType, NewTypeSubstitution(), nil
    case float64:
        return FloatType, NewTypeSubstitution(), nil
    case string:
        return StringType, NewTypeSubstitution(), nil
    case bool:
        return BoolType, NewTypeSubstitution(), nil
    default:
        return nil, nil, fmt.Errorf("unsupported literal type")
    }
}

// combineSubstitutions 合并替换
func (ti *TypeInference) combineSubstitutions(sub1, sub2 *TypeSubstitution) *TypeSubstitution {
    combined := NewTypeSubstitution()
    
    // 添加sub1的映射
    for k, v := range sub1.Mapping {
        combined.Mapping[k] = v
    }
    
    // 添加sub2的映射
    for k, v := range sub2.Mapping {
        combined.Mapping[k] = v
    }
    
    return combined
}

// SolveConstraints 求解约束
func (ti *TypeInference) SolveConstraints() (*TypeSubstitution, error) {
    solution := NewTypeSubstitution()
    
    for _, constraint := range ti.constraints {
        sub, err := ti.unify(constraint.Left, constraint.Right)
        if err != nil {
            return nil, err
        }
        
        // 应用替换到所有约束
        for i := range ti.constraints {
            ti.constraints[i].Left = sub.Apply(ti.constraints[i].Left)
            ti.constraints[i].Right = sub.Apply(ti.constraints[i].Right)
        }
        
        // 合并替换
        solution = ti.combineSubstitutions(solution, sub)
    }
    
    return solution, nil
}

// unify 统一算法
func (ti *TypeInference) unify(t1, t2 Type) (*TypeSubstitution, error) {
    sub := NewTypeSubstitution()
    
    // 如果类型相同，返回空替换
    if t1.IsEquivalentTo(t2) {
        return sub, nil
    }
    
    // 如果t1是类型变量
    if tv1, ok := t1.(*TypeVariable); ok {
        if ti.occursIn(tv1, t2) {
            return nil, fmt.Errorf("circular type constraint")
        }
        sub.Add(tv1.Name, t2)
        return sub, nil
    }
    
    // 如果t2是类型变量
    if tv2, ok := t2.(*TypeVariable); ok {
        if ti.occursIn(tv2, t1) {
            return nil, fmt.Errorf("circular type constraint")
        }
        sub.Add(tv2.Name, t1)
        return sub, nil
    }
    
    // 如果都是函数类型
    if ft1, ok1 := t1.(*FunctionType); ok1 {
        if ft2, ok2 := t2.(*FunctionType); ok2 {
            if len(ft1.Parameters) != len(ft2.Parameters) {
                return nil, fmt.Errorf("function parameter count mismatch")
            }
            
            // 统一参数类型
            for i := range ft1.Parameters {
                sub1, err := ti.unify(ft1.Parameters[i], ft2.Parameters[i])
                if err != nil {
                    return nil, err
                }
                sub = ti.combineSubstitutions(sub, sub1)
            }
            
            // 统一返回类型
            sub1, err := ti.unify(ft1.ReturnType, ft2.ReturnType)
            if err != nil {
                return nil, err
            }
            sub = ti.combineSubstitutions(sub, sub1)
            
            return sub, nil
        }
    }
    
    return nil, fmt.Errorf("cannot unify types %s and %s", t1.String(), t2.String())
}

// occursIn 检查类型变量是否出现在类型中
func (ti *TypeInference) occursIn(tv *TypeVariable, typ Type) bool {
    switch t := typ.(type) {
    case *TypeVariable:
        return tv.Name == t.Name
    case *FunctionType:
        for _, param := range t.Parameters {
            if ti.occursIn(tv, param) {
                return true
            }
        }
        return ti.occursIn(tv, t.ReturnType)
    case *ArrayType:
        return ti.occursIn(tv, t.ElementType)
    case *SliceType:
        return ti.occursIn(tv, t.ElementType)
    case *MapType:
        return ti.occursIn(tv, t.KeyType) || ti.occursIn(tv, t.ValueType)
    case *StructType:
        for _, fieldType := range t.Fields {
            if ti.occursIn(tv, fieldType) {
                return true
            }
        }
        return false
    default:
        return false
    }
}
```

### Go语言类型推导

```go
// GoTypeInference Go语言类型推导
type GoTypeInference struct {
    typeInference *TypeInference
}

// NewGoTypeInference 创建Go语言类型推导器
func NewGoTypeInference() *GoTypeInference {
    return &GoTypeInference{
        typeInference: NewTypeInference(),
    }
}

// InferExpression 推导表达式类型
func (gti *GoTypeInference) InferExpression(expr string, env TypeEnvironment) (Type, error) {
    // 解析表达式
    parsedExpr, err := gti.parseExpression(expr)
    if err != nil {
        return nil, err
    }
    
    // 推导类型
    typ, sub, err := gti.typeInference.InferType(parsedExpr, env)
    if err != nil {
        return nil, err
    }
    
    // 求解约束
    solution, err := gti.typeInference.SolveConstraints()
    if err != nil {
        return nil, err
    }
    
    // 应用最终替换
    finalType := solution.Apply(typ)
    
    return finalType, nil
}

// parseExpression 解析表达式（简化版本）
func (gti *GoTypeInference) parseExpression(expr string) (Expression, error) {
    // 这里应该实现完整的表达式解析
    // 为了演示，我们返回一个简单的字面量
    return &Literal{Value: 42}, nil
}
```

### 类型约束求解

```go
// ConstraintSolver 约束求解器
type ConstraintSolver struct {
    constraints []TypeConstraint
}

// NewConstraintSolver 创建约束求解器
func NewConstraintSolver() *ConstraintSolver {
    return &ConstraintSolver{
        constraints: make([]TypeConstraint, 0),
    }
}

// AddConstraint 添加约束
func (cs *ConstraintSolver) AddConstraint(left, right Type) {
    cs.constraints = append(cs.constraints, TypeConstraint{
        Left:  left,
        Right: right,
    })
}

// Solve 求解约束
func (cs *ConstraintSolver) Solve() (*TypeSubstitution, error) {
    solution := NewTypeSubstitution()
    
    for _, constraint := range cs.constraints {
        sub, err := cs.unify(constraint.Left, constraint.Right)
        if err != nil {
            return nil, err
        }
        
        // 应用替换到所有约束
        for i := range cs.constraints {
            cs.constraints[i].Left = sub.Apply(cs.constraints[i].Left)
            cs.constraints[i].Right = sub.Apply(cs.constraints[i].Right)
        }
        
        // 合并替换
        solution = cs.combineSubstitutions(solution, sub)
    }
    
    return solution, nil
}

// unify 统一算法
func (cs *ConstraintSolver) unify(t1, t2 Type) (*TypeSubstitution, error) {
    // 实现与TypeInference.unify相同的逻辑
    return NewTypeSubstitution(), nil
}

// combineSubstitutions 合并替换
func (cs *ConstraintSolver) combineSubstitutions(sub1, sub2 *TypeSubstitution) *TypeSubstitution {
    // 实现与TypeInference.combineSubstitutions相同的逻辑
    return NewTypeSubstitution()
}
```

## 类型安全证明

### 进展定理

**定理 5 (进展定理)**: 如果 $\Gamma \vdash e: T$ 且 $e$ 是封闭的，则要么 $e$ 是值，要么存在 $e'$ 使得 $e \rightarrow e'$。

**证明**: 通过对表达式 $e$ 的结构进行归纳。

### 保持定理

**定理 6 (保持定理)**: 如果 $\Gamma \vdash e: T$ 且 $e \rightarrow e'$，则 $\Gamma \vdash e': T$。

**证明**: 通过对归约规则进行案例分析。

### 类型健全性

**定理 7 (类型健全性)**: 如果 $\Gamma \vdash e: T$ 且 $e$ 是封闭的，则 $e$ 不会产生类型错误。

**证明**: 结合进展定理和保持定理。

## 应用领域

### 编译器设计

类型基础在编译器设计中的应用：
- 类型检查
- 类型推导
- 代码生成
- 优化

### 静态分析

类型基础在静态分析中的应用：
- 类型推断
- 错误检测
- 程序验证
- 代码重构

### 程序验证

类型基础在程序验证中的应用：
- 类型安全证明
- 程序正确性
- 形式化验证
- 定理证明

### 代码生成

类型基础在代码生成中的应用：
- 类型特化
- 泛型实例化
- 优化代码
- 目标代码

## 相关链接

- [02-类型推导 (Type Inference)](../02-Type-Inference/README.md)
- [03-类型安全 (Type Safety)](../03-Type-Safety/README.md)
- [04-高级类型系统 (Advanced Type Systems)](../04-Advanced-Type-Systems/README.md)
- [02-语义学理论 (Semantics Theory)](../../02-Semantics-Theory/README.md)
- [03-编译原理 (Compiler Theory)](../../03-Compiler-Theory/README.md) 