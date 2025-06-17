# 02-函数式编程模式 (Functional Patterns)

## 概述

函数式编程模式是现代软件设计中的重要范式，强调不可变性、纯函数和高阶函数的使用。本文档探讨在Go语言中实现函数式编程模式的方法和理论。

## 目录

1. [理论基础](#理论基础)
2. [核心概念](#核心概念)
3. [函数式模式实现](#函数式模式实现)
4. [高阶函数](#高阶函数)
5. [不可变性模式](#不可变性模式)
6. [函数组合](#函数组合)
7. [Monad模式](#monad模式)
8. [实际应用](#实际应用)
9. [性能考虑](#性能考虑)
10. [最佳实践](#最佳实践)

## 理论基础

### 函数式编程的数学基础

函数式编程基于λ演算（Lambda Calculus）的数学理论：

```latex
\lambda x.e \quad \text{(Lambda abstraction)}
```

其中：

- $x$ 是参数
- $e$ 是表达式
- $\lambda$ 表示函数抽象

### 纯函数定义

纯函数满足以下数学性质：

```latex
f: A \rightarrow B
```

对于任意输入 $a \in A$，输出 $f(a) \in B$ 是确定的，且没有副作用。

### 函数组合

函数组合的数学定义：

```latex
(f \circ g)(x) = f(g(x))
```

其中 $f$ 和 $g$ 是可组合的函数。

## 核心概念

### 1. 纯函数 (Pure Functions)

纯函数是函数式编程的核心概念，具有以下特性：

- **确定性**：相同输入总是产生相同输出
- **无副作用**：不修改外部状态
- **引用透明性**：函数调用可以被其返回值替换

```go
// 纯函数示例
func add(a, b int) int {
    return a + b
}

// 非纯函数示例
var counter int
func increment() int {
    counter++
    return counter
}
```

### 2. 不可变性 (Immutability)

不可变性确保数据一旦创建就不能被修改：

```go
// 不可变数据结构
type ImmutablePoint struct {
    x, y int
}

func (p ImmutablePoint) Move(dx, dy int) ImmutablePoint {
    return ImmutablePoint{x: p.x + dx, y: p.y + dy}
}
```

### 3. 高阶函数 (Higher-Order Functions)

高阶函数接受函数作为参数或返回函数：

```go
// 接受函数作为参数
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// 返回函数
func Adder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}
```

## 函数式模式实现

### 1. Map模式

Map模式将函数应用到集合的每个元素：

```go
// 通用Map实现
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// 使用示例
numbers := []int{1, 2, 3, 4, 5}
squares := Map(numbers, func(n int) int {
    return n * n
})
// 结果: [1, 4, 9, 16, 25]
```

### 2. Filter模式

Filter模式根据谓词函数过滤集合：

```go
// 通用Filter实现
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// 使用示例
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
evens := Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// 结果: [2, 4, 6, 8, 10]
```

### 3. Reduce模式

Reduce模式将集合归约为单个值：

```go
// 通用Reduce实现
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

// 使用示例
numbers := []int{1, 2, 3, 4, 5}
sum := Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
})
// 结果: 15
```

## 高阶函数

### 1. 柯里化 (Currying)

柯里化将接受多个参数的函数转换为接受单个参数的函数链：

```go
// 柯里化实现
func Curry2[A, B, C any](fn func(A, B) C) func(A) func(B) C {
    return func(a A) func(B) C {
        return func(b B) C {
            return fn(a, b)
        }
    }
}

// 使用示例
add := func(a, b int) int { return a + b }
curriedAdd := Curry2(add)
add5 := curriedAdd(5)
result := add5(3) // 结果: 8
```

### 2. 部分应用 (Partial Application)

部分应用固定函数的部分参数：

```go
// 部分应用实现
func Partial1[A, B, C any](fn func(A, B) C, a A) func(B) C {
    return func(b B) C {
        return fn(a, b)
    }
}

// 使用示例
multiply := func(a, b int) int { return a * b }
multiplyBy2 := Partial1(multiply, 2)
result := multiplyBy2(5) // 结果: 10
```

## 不可变性模式

### 1. 不可变数据结构

```go
// 不可变链表
type ImmutableList[T any] struct {
    head T
    tail *ImmutableList[T]
}

func NewImmutableList[T any](head T) *ImmutableList[T] {
    return &ImmutableList[T]{head: head}
}

func (l *ImmutableList[T]) Prepend(item T) *ImmutableList[T] {
    return &ImmutableList[T]{
        head: item,
        tail: l,
    }
}

func (l *ImmutableList[T]) Head() T {
    return l.head
}

func (l *ImmutableList[T]) Tail() *ImmutableList[T] {
    return l.tail
}
```

### 2. 不可变Map

```go
// 不可变Map实现
type ImmutableMap[K comparable, V any] struct {
    data map[K]V
}

func NewImmutableMap[K comparable, V any]() *ImmutableMap[K, V] {
    return &ImmutableMap[K, V]{
        data: make(map[K]V),
    }
}

func (m *ImmutableMap[K, V]) Set(key K, value V) *ImmutableMap[K, V] {
    newData := make(map[K]V)
    for k, v := range m.data {
        newData[k] = v
    }
    newData[key] = value
    
    return &ImmutableMap[K, V]{data: newData}
}

func (m *ImmutableMap[K, V]) Get(key K) (V, bool) {
    value, exists := m.data[key]
    return value, exists
}
```

## 函数组合

### 1. 函数管道 (Function Pipeline)

```go
// 函数管道实现
type Pipeline[T any] struct {
    functions []func(T) T
}

func NewPipeline[T any]() *Pipeline[T] {
    return &Pipeline[T]{}
}

func (p *Pipeline[T]) Add(fn func(T) T) *Pipeline[T] {
    p.functions = append(p.functions, fn)
    return p
}

func (p *Pipeline[T]) Execute(input T) T {
    result := input
    for _, fn := range p.functions {
        result = fn(result)
    }
    return result
}

// 使用示例
pipeline := NewPipeline[int]().
    Add(func(x int) int { return x * 2 }).
    Add(func(x int) int { return x + 1 }).
    Add(func(x int) int { return x * x })

result := pipeline.Execute(5) // 结果: 121 ((5*2+1)^2)
```

### 2. 函数组合器

```go
// 函数组合器
func Compose[T any](functions ...func(T) T) func(T) T {
    return func(input T) T {
        result := input
        for i := len(functions) - 1; i >= 0; i-- {
            result = functions[i](result)
        }
        return result
    }
}

// 使用示例
double := func(x int) int { return x * 2 }
addOne := func(x int) int { return x + 1 }
square := func(x int) int { return x * x }

composed := Compose(square, addOne, double)
result := composed(5) // 结果: 121
```

## Monad模式

### 1. Option Monad

Option Monad处理可能为空的值：

```go
// Option类型
type Option[T any] struct {
    value T
    hasValue bool
}

func Some[T any](value T) Option[T] {
    return Option[T]{value: value, hasValue: true}
}

func None[T any]() Option[T] {
    return Option[T]{hasValue: false}
}

func (o Option[T]) IsSome() bool {
    return o.hasValue
}

func (o Option[T]) IsNone() bool {
    return !o.hasValue
}

func (o Option[T]) Unwrap() T {
    if !o.hasValue {
        panic("attempting to unwrap None")
    }
    return o.value
}

func (o Option[T]) UnwrapOr(defaultValue T) T {
    if o.hasValue {
        return o.value
    }
    return defaultValue
}

// Monad操作
func (o Option[T]) Map[U any](fn func(T) U) Option[U] {
    if o.hasValue {
        return Some(fn(o.value))
    }
    return None[U]()
}

func (o Option[T]) FlatMap[U any](fn func(T) Option[U]) Option[U] {
    if o.hasValue {
        return fn(o.value)
    }
    return None[U]()
}
```

### 2. Result Monad

Result Monad处理可能失败的操作：

```go
// Result类型
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](value T) Result[T] {
    return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
    return Result[T]{err: err}
}

func (r Result[T]) IsOk() bool {
    return r.err == nil
}

func (r Result[T]) IsErr() bool {
    return r.err != nil
}

func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}

// Monad操作
func (r Result[T]) Map[U any](fn func(T) U) Result[U] {
    if r.err != nil {
        return Err[U](r.err)
    }
    return Ok(fn(r.value))
}

func (r Result[T]) FlatMap[U any](fn func(T) Result[U]) Result[U] {
    if r.err != nil {
        return Err[U](r.err)
    }
    return fn(r.value)
}
```

## 实际应用

### 1. 数据处理管道

```go
// 数据处理管道示例
type DataProcessor struct{}

func (dp *DataProcessor) ProcessData(data []int) []int {
    pipeline := NewPipeline[[]int]().
        Add(dp.filterValid).
        Add(dp.transform).
        Add(dp.sort).
        Add(dp.deduplicate)
    
    return pipeline.Execute(data)
}

func (dp *DataProcessor) filterValid(data []int) []int {
    return Filter(data, func(n int) bool {
        return n > 0 && n < 1000
    })
}

func (dp *DataProcessor) transform(data []int) []int {
    return Map(data, func(n int) int {
        return n * 2 + 1
    })
}

func (dp *DataProcessor) sort(data []int) []int {
    sorted := make([]int, len(data))
    copy(sorted, data)
    sort.Ints(sorted)
    return sorted
}

func (dp *DataProcessor) deduplicate(data []int) []int {
    seen := make(map[int]bool)
    return Filter(data, func(n int) bool {
        if seen[n] {
            return false
        }
        seen[n] = true
        return true
    })
}
```

### 2. 配置管理

```go
// 函数式配置管理
type Config struct {
    DatabaseURL string
    Port        int
    Debug       bool
}

type ConfigBuilder struct {
    config Config
}

func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{}
}

func (cb *ConfigBuilder) WithDatabaseURL(url string) *ConfigBuilder {
    cb.config.DatabaseURL = url
    return cb
}

func (cb *ConfigBuilder) WithPort(port int) *ConfigBuilder {
    cb.config.Port = port
    return cb
}

func (cb *ConfigBuilder) WithDebug(debug bool) *ConfigBuilder {
    cb.config.Debug = debug
    return cb
}

func (cb *ConfigBuilder) Build() Config {
    return cb.config
}

// 使用示例
config := NewConfigBuilder().
    WithDatabaseURL("localhost:5432").
    WithPort(8080).
    WithDebug(true).
    Build()
```

## 性能考虑

### 1. 内存分配优化

```go
// 预分配切片以减少内存分配
func MapOptimized[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, 0, len(slice)) // 预分配容量
    for _, v := range slice {
        result = append(result, fn(v))
    }
    return result
}
```

### 2. 惰性求值

```go
// 惰性求值实现
type Lazy[T any] struct {
    compute func() T
    value   *T
}

func NewLazy[T any](compute func() T) *Lazy[T] {
    return &Lazy[T]{compute: compute}
}

func (l *Lazy[T]) Get() T {
    if l.value == nil {
        val := l.compute()
        l.value = &val
    }
    return *l.value
}
```

## 最佳实践

### 1. 函数设计原则

- **单一职责**：每个函数只做一件事
- **纯函数优先**：尽可能使用纯函数
- **不可变性**：避免修改输入参数
- **组合优于继承**：使用函数组合而不是复杂的继承结构

### 2. 错误处理

```go
// 函数式错误处理
func SafeDivide(a, b float64) Option[float64] {
    if b == 0 {
        return None[float64]()
    }
    return Some(a / b)
}

// 链式错误处理
func ProcessData(data []int) Result[[]int] {
    return Ok(data).
        FlatMap(validateData).
        FlatMap(transformData).
        FlatMap(saveData)
}
```

### 3. 测试策略

```go
// 函数式代码测试
func TestMap(t *testing.T) {
    input := []int{1, 2, 3, 4, 5}
    expected := []int{2, 4, 6, 8, 10}
    
    result := Map(input, func(n int) int {
        return n * 2
    })
    
    if !reflect.DeepEqual(result, expected) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## 总结

函数式编程模式在Go语言中提供了强大的抽象能力和代码组织方式。通过使用纯函数、不可变数据结构和高阶函数，我们可以构建更加可读、可测试和可维护的代码。

关键要点：

- 优先使用纯函数
- 利用不可变性避免副作用
- 使用高阶函数进行抽象
- 通过函数组合构建复杂逻辑
- 使用Monad模式处理副作用

这些模式不仅提高了代码质量，还为并发编程和错误处理提供了更好的基础。

---

**相关链接**：

- [01-响应式模式](./01-Reactive-Patterns.md)
- [03-事件溯源模式](./03-Event-Sourcing-Patterns.md)
- [04-CQRS模式](./04-CQRS-Patterns.md)
- [../README.md](../README.md)
