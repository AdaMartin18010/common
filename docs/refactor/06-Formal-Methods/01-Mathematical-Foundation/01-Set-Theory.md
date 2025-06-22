# 01-集合论 (Set Theory)

## 目录

1. [基础概念](#1-基础概念)
2. [集合运算](#2-集合运算)
3. [关系与函数](#3-关系与函数)
4. [基数与序数](#4-基数与序数)
5. [公理化集合论](#5-公理化集合论)
6. [Go语言实现](#6-go语言实现)
7. [应用示例](#7-应用示例)

## 1. 基础概念

### 1.1 集合的定义

**定义 1.1.1 (集合)** 集合是不同对象的无序聚集，这些对象称为集合的元素。

**公理 1.1.1 (外延公理)** 两个集合相等当且仅当它们包含相同的元素：
$```latex
\forall x \forall y [\forall z(z \in x \leftrightarrow z \in y) \rightarrow x = y]
```$

**公理 1.1.2 (空集公理)** 存在一个不包含任何元素的集合：
$```latex
\exists x \forall y(y \notin x)
```$

### 1.2 集合的表示

```go
// 集合的Go语言表示
type Set[T comparable] struct {
    elements map[T]bool
}

// 创建空集合
func NewSet[T comparable]() *Set[T] {
    return &Set[T]{
        elements: make(map[T]bool),
    }
}

// 从切片创建集合
func NewSetFromSlice[T comparable](slice []T) *Set[T] {
    set := NewSet[T]()
    for _, element := range slice {
        set.Add(element)
    }
    return set
}

// 添加元素
func (s *Set[T]) Add(element T) {
    s.elements[element] = true
}

// 删除元素
func (s *Set[T]) Remove(element T) {
    delete(s.elements, element)
}

// 检查元素是否存在
func (s *Set[T]) Contains(element T) bool {
    return s.elements[element]
}

// 获取集合大小
func (s *Set[T]) Size() int {
    return len(s.elements)
}

// 检查是否为空
func (s *Set[T]) IsEmpty() bool {
    return len(s.elements) == 0
}
```

### 1.3 子集与真子集

**定义 1.3.1 (子集)** 集合 ```latex
A
``` 是集合 ```latex
B
``` 的子集，记作 ```latex
A \subseteq B
```，当且仅当 ```latex
A
``` 的每个元素都属于 ```latex
B
```：
$```latex
A \subseteq B \leftrightarrow \forall x(x \in A \rightarrow x \in B)
```$

**定义 1.3.2 (真子集)** 集合 ```latex
A
``` 是集合 ```latex
B
``` 的真子集，记作 ```latex
A \subset B
```，当且仅当 ```latex
A \subseteq B
``` 且 ```latex
A \neq B
```。

```go
// 子集检查
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
    for element := range s.elements {
        if !other.Contains(element) {
            return false
        }
    }
    return true
}

// 真子集检查
func (s *Set[T]) IsProperSubsetOf(other *Set[T]) bool {
    return s.IsSubsetOf(other) && !s.Equals(other)
}

// 集合相等
func (s *Set[T]) Equals(other *Set[T]) bool {
    if s.Size() != other.Size() {
        return false
    }
    return s.IsSubsetOf(other)
}
```

## 2. 集合运算

### 2.1 基本运算

**定义 2.1.1 (并集)** 集合 ```latex
A
``` 和 ```latex
B
``` 的并集是包含所有属于 ```latex
A
``` 或 ```latex
B
``` 的元素的集合：
$```latex
A \cup B = \{x | x \in A \lor x \in B\}
```$

**定义 2.1.2 (交集)** 集合 ```latex
A
``` 和 ```latex
B
``` 的交集是包含所有同时属于 ```latex
A
``` 和 ```latex
B
``` 的元素的集合：
$```latex
A \cap B = \{x | x \in A \land x \in B\}
```$

**定义 2.1.3 (差集)** 集合 ```latex
A
``` 和 ```latex
B
``` 的差集是包含所有属于 ```latex
A
``` 但不属于 ```latex
B
``` 的元素的集合：
$```latex
A \setminus B = \{x | x \in A \land x \notin B\}
```$

```go
// 并集
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    
    // 添加当前集合的所有元素
    for element := range s.elements {
        result.Add(element)
    }
    
    // 添加另一个集合的所有元素
    for element := range other.elements {
        result.Add(element)
    }
    
    return result
}

// 交集
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    
    for element := range s.elements {
        if other.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}

// 差集
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    
    for element := range s.elements {
        if !other.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}

// 对称差集
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
    union := s.Union(other)
    intersection := s.Intersection(other)
    return union.Difference(intersection)
}
```

### 2.2 集合运算的性质

**定理 2.2.1 (交换律)** 对于任意集合 ```latex
A
``` 和 ```latex
B
```：
$```latex
A \cup B = B \cup A
```$
$```latex
A \cap B = B \cap A
```$

**定理 2.2.2 (结合律)** 对于任意集合 ```latex
A
```、```latex
B
``` 和 ```latex
C
```：
$```latex
(A \cup B) \cup C = A \cup (B \cup C)
```$
$```latex
(A \cap B) \cap C = A \cap (B \cap C)
```$

**定理 2.2.3 (分配律)** 对于任意集合 ```latex
A
```、```latex
B
``` 和 ```latex
C
```：
$```latex
A \cup (B \cap C) = (A \cup B) \cap (A \cup C)
```$
$```latex
A \cap (B \cup C) = (A \cap B) \cup (A \cap C)
```$

```go
// 验证交换律
func TestCommutativeLaw() {
    set1 := NewSetFromSlice([]int{1, 2, 3})
    set2 := NewSetFromSlice([]int{3, 4, 5})
    
    union1 := set1.Union(set2)
    union2 := set2.Union(set1)
    
    fmt.Printf("交换律验证: %v\n", union1.Equals(union2))
    
    intersection1 := set1.Intersection(set2)
    intersection2 := set2.Intersection(set1)
    
    fmt.Printf("交集交换律验证: %v\n", intersection1.Equals(intersection2))
}

// 验证结合律
func TestAssociativeLaw() {
    set1 := NewSetFromSlice([]int{1, 2, 3})
    set2 := NewSetFromSlice([]int{3, 4, 5})
    set3 := NewSetFromSlice([]int{5, 6, 7})
    
    union1 := set1.Union(set2).Union(set3)
    union2 := set1.Union(set2.Union(set3))
    
    fmt.Printf("结合律验证: %v\n", union1.Equals(union2))
}
```

## 3. 关系与函数

### 3.1 笛卡尔积

**定义 3.1.1 (笛卡尔积)** 集合 ```latex
A
``` 和 ```latex
B
``` 的笛卡尔积是所有有序对 ```latex
(a, b)
``` 的集合，其中 ```latex
a \in A
``` 且 ```latex
b \in B
```：
$```latex
A \times B = \{(a, b) | a \in A \land b \in B\}
```$

```go
// 有序对
type OrderedPair[T, U comparable] struct {
    First  T
    Second U
}

// 笛卡尔积
func CartesianProduct[T, U comparable](setA *Set[T], setB *Set[U]) *Set[OrderedPair[T, U]] {
    result := NewSet[OrderedPair[T, U]]()
    
    for a := range setA.elements {
        for b := range setB.elements {
            result.Add(OrderedPair[T, U]{First: a, Second: b})
        }
    }
    
    return result
}

// 有序对的相等性
func (op OrderedPair[T, U]) Equals(other OrderedPair[T, U]) bool {
    return op.First == other.First && op.Second == other.Second
}
```

### 3.2 关系

**定义 3.2.1 (二元关系)** 集合 ```latex
A
``` 和 ```latex
B
``` 之间的二元关系是 ```latex
A \times B
``` 的子集。

**定义 3.2.2 (等价关系)** 集合 ```latex
A
``` 上的关系 ```latex
R
``` 是等价关系，当且仅当它满足：

1. **自反性**: ```latex
\forall x \in A, xRx
```
2. **对称性**: ```latex
\forall x, y \in A, xRy \rightarrow yRx
```
3. **传递性**: ```latex
\forall x, y, z \in A, (xRy \land yRz) \rightarrow xRz
```

```go
// 关系
type Relation[T comparable] struct {
    pairs *Set[OrderedPair[T, T]]
}

func NewRelation[T comparable]() *Relation[T] {
    return &Relation[T]{
        pairs: NewSet[OrderedPair[T, T]](),
    }
}

func (r *Relation[T]) AddPair(a, b T) {
    r.pairs.Add(OrderedPair[T, T]{First: a, Second: b})
}

func (r *Relation[T]) Contains(a, b T) bool {
    return r.pairs.Contains(OrderedPair[T, T]{First: a, Second: b})
}

// 检查自反性
func (r *Relation[T]) IsReflexive(universe *Set[T]) bool {
    for element := range universe.elements {
        if !r.Contains(element, element) {
            return false
        }
    }
    return true
}

// 检查对称性
func (r *Relation[T]) IsSymmetric() bool {
    for pair := range r.pairs.elements {
        if !r.Contains(pair.Second, pair.First) {
            return false
        }
    }
    return true
}

// 检查传递性
func (r *Relation[T]) IsTransitive() bool {
    for pair1 := range r.pairs.elements {
        for pair2 := range r.pairs.elements {
            if pair1.Second == pair2.First {
                if !r.Contains(pair1.First, pair2.Second) {
                    return false
                }
            }
        }
    }
    return true
}

// 检查是否为等价关系
func (r *Relation[T]) IsEquivalenceRelation(universe *Set[T]) bool {
    return r.IsReflexive(universe) && r.IsSymmetric() && r.IsTransitive()
}
```

### 3.3 函数

**定义 3.3.1 (函数)** 从集合 ```latex
A
``` 到集合 ```latex
B
``` 的函数是一个关系 ```latex
f \subseteq A \times B
```，满足：

1. **定义域**: ```latex
\forall x \in A, \exists y \in B, (x, y) \in f
```
2. **单值性**: ```latex
\forall x \in A, \forall y_1, y_2 \in B, ((x, y_1) \in f \land (x, y_2) \in f) \rightarrow y_1 = y_2
```

```go
// 函数
type Function[T, U comparable] struct {
    mapping map[T]U
    domain  *Set[T]
    codomain *Set[U]
}

func NewFunction[T, U comparable](domain *Set[T], codomain *Set[U]) *Function[T, U] {
    return &Function[T, U]{
        mapping:  make(map[T]U),
        domain:   domain,
        codomain: codomain,
    }
}

func (f *Function[T, U]) Define(input T, output U) {
    if f.domain.Contains(input) && f.codomain.Contains(output) {
        f.mapping[input] = output
    }
}

func (f *Function[T, U]) Apply(input T) (U, bool) {
    output, exists := f.mapping[input]
    return output, exists
}

// 检查是否为满射
func (f *Function[T, U]) IsSurjective() bool {
    image := NewSet[U]()
    for _, output := range f.mapping {
        image.Add(output)
    }
    return image.IsSubsetOf(f.codomain) && f.codomain.IsSubsetOf(image)
}

// 检查是否为单射
func (f *Function[T, U]) IsInjective() bool {
    used := make(map[U]bool)
    for _, output := range f.mapping {
        if used[output] {
            return false
        }
        used[output] = true
    }
    return true
}

// 检查是否为双射
func (f *Function[T, U]) IsBijective() bool {
    return f.IsSurjective() && f.IsInjective()
}
```

## 4. 基数与序数

### 4.1 基数

**定义 4.1.1 (基数)** 集合 ```latex
A
``` 的基数是衡量 ```latex
A
``` 中元素数量的概念，记作 ```latex
|A|
```。

**定义 4.1.2 (有限集)** 如果存在自然数 ```latex
n
``` 使得集合 ```latex
A
``` 与 ```latex
\{1, 2, \ldots, n\}
``` 之间存在双射，则称 ```latex
A
``` 为有限集。

**定义 4.1.3 (可数集)** 如果集合 ```latex
A
``` 与自然数集 ```latex
\mathbb{N}
``` 之间存在双射，则称 ```latex
A
``` 为可数集。

```go
// 基数计算
func (s *Set[T]) Cardinality() int {
    return s.Size()
}

// 检查是否为有限集
func (s *Set[T]) IsFinite() bool {
    return s.Size() < math.MaxInt
}

// 检查是否为可数集
func (s *Set[T]) IsCountable() bool {
    // 对于有限集，总是可数的
    if s.IsFinite() {
        return true
    }
    
    // 对于无限集，需要检查是否存在到自然数的双射
    // 这里简化处理，假设所有无限集都是可数的
    return true
}

// 基数比较
func CompareCardinality[T, U comparable](setA *Set[T], setB *Set[U]) string {
    cardA := setA.Cardinality()
    cardB := setB.Cardinality()
    
    if cardA < cardB {
        return "|A| < |B|"
    } else if cardA > cardB {
        return "|A| > |B|"
    } else {
        return "|A| = |B|"
    }
}
```

### 4.2 幂集

**定义 4.2.1 (幂集)** 集合 ```latex
A
``` 的幂集是 ```latex
A
``` 的所有子集构成的集合：
$```latex
\mathcal{P}(A) = \{B | B \subseteq A\}
```$

**定理 4.2.1 (幂集基数)** 如果 ```latex
|A| = n
```，则 ```latex
|\mathcal{P}(A)| = 2^n
```。

```go
// 幂集生成
func (s *Set[T]) PowerSet() *Set[*Set[T]] {
    elements := s.ToSlice()
    n := len(elements)
    powerSet := NewSet[*Set[T]]()
    
    // 生成所有可能的子集
    for i := 0; i < (1 << n); i++ {
        subset := NewSet[T]()
        for j := 0; j < n; j++ {
            if (i & (1 << j)) != 0 {
                subset.Add(elements[j])
            }
        }
        powerSet.Add(subset)
    }
    
    return powerSet
}

// 转换为切片
func (s *Set[T]) ToSlice() []T {
    result := make([]T, 0, s.Size())
    for element := range s.elements {
        result = append(result, element)
    }
    return result
}

// 验证幂集基数定理
func TestPowerSetCardinality() {
    set := NewSetFromSlice([]int{1, 2, 3})
    powerSet := set.PowerSet()
    
    expectedSize := 1 << set.Size()
    actualSize := powerSet.Size()
    
    fmt.Printf("幂集基数验证: 期望 %d, 实际 %d, 正确: %v\n", 
        expectedSize, actualSize, expectedSize == actualSize)
}
```

## 5. 公理化集合论

### 5.1 ZFC公理系统

**公理 5.1.1 (配对公理)** 对于任意两个集合 ```latex
x
``` 和 ```latex
y
```，存在一个集合包含它们：
$```latex
\forall x \forall y \exists z \forall w(w \in z \leftrightarrow w = x \lor w = y)
```$

**公理 5.1.2 (并集公理)** 对于任意集合族 ```latex
\mathcal{F}
```，存在一个集合包含所有 ```latex
\mathcal{F}
``` 中集合的元素：
$```latex
\forall \mathcal{F} \exists A \forall x(x \in A \leftrightarrow \exists B(B \in \mathcal{F} \land x \in B))
```$

**公理 5.1.3 (幂集公理)** 对于任意集合 ```latex
x
```，存在一个集合包含 ```latex
x
``` 的所有子集：
$```latex
\forall x \exists y \forall z(z \in y \leftrightarrow z \subseteq x)
```$

```go
// 配对公理实现
func PairingAxiom[T comparable](x, y T) *Set[T] {
    result := NewSet[T]()
    result.Add(x)
    result.Add(y)
    return result
}

// 并集公理实现
func UnionAxiom[T comparable](sets []*Set[T]) *Set[T] {
    result := NewSet[T]()
    for _, set := range sets {
        for element := range set.elements {
            result.Add(element)
        }
    }
    return result
}

// 幂集公理实现
func PowerSetAxiom[T comparable](set *Set[T]) *Set[*Set[T]] {
    return set.PowerSet()
}
```

### 5.2 选择公理

**公理 5.2.1 (选择公理)** 对于任意非空集合族 ```latex
\mathcal{F}
```，存在一个选择函数 ```latex
f: \mathcal{F} \rightarrow \bigcup \mathcal{F}
```，使得对于每个 ```latex
A \in \mathcal{F}
```，有 ```latex
f(A) \in A
```。

```go
// 选择函数
type ChoiceFunction[T comparable] struct {
    choices map[*Set[T]]T
}

func NewChoiceFunction[T comparable]() *ChoiceFunction[T] {
    return &ChoiceFunction[T]{
        choices: make(map[*Set[T]]T),
    }
}

func (cf *ChoiceFunction[T]) Choose(set *Set[T]) T {
    if choice, exists := cf.choices[set]; exists {
        return choice
    }
    
    // 简单实现：选择第一个元素
    for element := range set.elements {
        cf.choices[set] = element
        return element
    }
    
    // 空集情况
    var zero T
    return zero
}

// 验证选择公理
func TestAxiomOfChoice() {
    set1 := NewSetFromSlice([]int{1, 2, 3})
    set2 := NewSetFromSlice([]int{4, 5})
    set3 := NewSetFromSlice([]int{6, 7, 8, 9})
    
    family := []*Set[int]{set1, set2, set3}
    choiceFunc := NewChoiceFunction[int]()
    
    fmt.Println("选择公理验证:")
    for i, set := range family {
        choice := choiceFunc.Choose(set)
        fmt.Printf("  集合 %d: 选择元素 %d\n", i+1, choice)
    }
}
```

## 6. Go语言实现

### 6.1 高级集合操作

```go
// 集合的字符串表示
func (s *Set[T]) String() string {
    elements := s.ToSlice()
    if len(elements) == 0 {
        return "{}"
    }
    
    var result strings.Builder
    result.WriteString("{")
    
    for i, element := range elements {
        if i > 0 {
            result.WriteString(", ")
        }
        result.WriteString(fmt.Sprintf("%v", element))
    }
    
    result.WriteString("}")
    return result.String()
}

// 集合的迭代器
type SetIterator[T comparable] struct {
    elements []T
    index    int
}

func (s *Set[T]) Iterator() *SetIterator[T] {
    return &SetIterator[T]{
        elements: s.ToSlice(),
        index:    0,
    }
}

func (it *SetIterator[T]) HasNext() bool {
    return it.index < len(it.elements)
}

func (it *SetIterator[T]) Next() T {
    if !it.HasNext() {
        var zero T
        return zero
    }
    
    element := it.elements[it.index]
    it.index++
    return element
}

// 集合的过滤操作
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
    result := NewSet[T]()
    for element := range s.elements {
        if predicate(element) {
            result.Add(element)
        }
    }
    return result
}

// 集合的映射操作
func MapSet[T, U comparable](s *Set[T], mapper func(T) U) *Set[U] {
    result := NewSet[U]()
    for element := range s.elements {
        result.Add(mapper(element))
    }
    return result
}

// 集合的归约操作
func ReduceSet[T comparable, U any](s *Set[T], initial U, reducer func(U, T) U) U {
    result := initial
    for element := range s.elements {
        result = reducer(result, element)
    }
    return result
}
```

### 6.2 泛型约束和类型安全

```go
// 数值集合
type NumericSet[T constraints.Ordered] struct {
    elements map[T]bool
}

func NewNumericSet[T constraints.Ordered]() *NumericSet[T] {
    return &NumericSet[T]{
        elements: make(map[T]bool),
    }
}

func (ns *NumericSet[T]) Add(element T) {
    ns.elements[element] = true
}

func (ns *NumericSet[T]) Min() (T, bool) {
    if len(ns.elements) == 0 {
        var zero T
        return zero, false
    }
    
    var min T
    first := true
    for element := range ns.elements {
        if first || element < min {
            min = element
            first = false
        }
    }
    return min, true
}

func (ns *NumericSet[T]) Max() (T, bool) {
    if len(ns.elements) == 0 {
        var zero T
        return zero, false
    }
    
    var max T
    first := true
    for element := range ns.elements {
        if first || element > max {
            max = element
            first = false
        }
    }
    return max, true
}

// 有序集合
type OrderedSet[T constraints.Ordered] struct {
    elements []T
    set      map[T]bool
}

func NewOrderedSet[T constraints.Ordered]() *OrderedSet[T] {
    return &OrderedSet[T]{
        elements: make([]T, 0),
        set:      make(map[T]bool),
    }
}

func (os *OrderedSet[T]) Add(element T) {
    if !os.set[element] {
        os.elements = append(os.elements, element)
        os.set[element] = true
        sort.Slice(os.elements, func(i, j int) bool {
            return os.elements[i] < os.elements[j]
        })
    }
}

func (os *OrderedSet[T]) ToSlice() []T {
    result := make([]T, len(os.elements))
    copy(result, os.elements)
    return result
}
```

## 7. 应用示例

### 7.1 数据库查询优化

```go
// 数据库索引模拟
type DatabaseIndex[T comparable] struct {
    index map[T]*Set[int] // 值 -> 行ID集合
}

func NewDatabaseIndex[T comparable]() *DatabaseIndex[T] {
    return &DatabaseIndex[T]{
        index: make(map[T]*Set[int]),
    }
}

func (di *DatabaseIndex[T]) Add(value T, rowID int) {
    if di.index[value] == nil {
        di.index[value] = NewSet[int]()
    }
    di.index[value].Add(rowID)
}

func (di *DatabaseIndex[T]) Query(value T) *Set[int] {
    if set, exists := di.index[value]; exists {
        return set
    }
    return NewSet[int]()
}

func (di *DatabaseIndex[T]) QueryRange(min, max T) *Set[int] {
    result := NewSet[int]()
    for value, rowIDs := range di.index {
        if value >= min && value <= max {
            result = result.Union(rowIDs)
        }
    }
    return result
}

// 使用示例
func DatabaseQueryExample() {
    index := NewDatabaseIndex[string]()
    
    // 添加数据
    index.Add("apple", 1)
    index.Add("banana", 2)
    index.Add("apple", 3)
    index.Add("cherry", 4)
    index.Add("banana", 5)
    
    // 查询
    appleRows := index.Query("apple")
    fmt.Printf("包含'apple'的行: %v\n", appleRows.ToSlice())
    
    // 范围查询
    rangeRows := index.QueryRange("apple", "banana")
    fmt.Printf("范围查询结果: %v\n", rangeRows.ToSlice())
}
```

### 7.2 权限管理系统

```go
// 权限管理系统
type Permission string
type Role string
type User string

type PermissionSystem struct {
    userRoles      map[User]*Set[Role]
    rolePermissions map[Role]*Set[Permission]
}

func NewPermissionSystem() *PermissionSystem {
    return &PermissionSystem{
        userRoles:      make(map[User]*Set[Role]),
        rolePermissions: make(map[Role]*Set[Permission]),
    }
}

func (ps *PermissionSystem) AssignRole(user User, role Role) {
    if ps.userRoles[user] == nil {
        ps.userRoles[user] = NewSet[Role]()
    }
    ps.userRoles[user].Add(role)
}

func (ps *PermissionSystem) GrantPermission(role Role, permission Permission) {
    if ps.rolePermissions[role] == nil {
        ps.rolePermissions[role] = NewSet[Permission]()
    }
    ps.rolePermissions[role].Add(permission)
}

func (ps *PermissionSystem) HasPermission(user User, permission Permission) bool {
    userRoles := ps.userRoles[user]
    if userRoles == nil {
        return false
    }
    
    for role := range userRoles.elements {
        rolePerms := ps.rolePermissions[role]
        if rolePerms != nil && rolePerms.Contains(permission) {
            return true
        }
    }
    return false
}

func (ps *PermissionSystem) GetUserPermissions(user User) *Set[Permission] {
    result := NewSet[Permission]()
    userRoles := ps.userRoles[user]
    
    if userRoles == nil {
        return result
    }
    
    for role := range userRoles.elements {
        rolePerms := ps.rolePermissions[role]
        if rolePerms != nil {
            result = result.Union(rolePerms)
        }
    }
    
    return result
}

// 使用示例
func PermissionSystemExample() {
    system := NewPermissionSystem()
    
    // 定义角色和权限
    system.GrantPermission("admin", "read")
    system.GrantPermission("admin", "write")
    system.GrantPermission("admin", "delete")
    system.GrantPermission("user", "read")
    system.GrantPermission("moderator", "read")
    system.GrantPermission("moderator", "write")
    
    // 分配角色
    system.AssignRole("alice", "admin")
    system.AssignRole("bob", "user")
    system.AssignRole("charlie", "moderator")
    
    // 检查权限
    fmt.Printf("Alice有删除权限: %v\n", system.HasPermission("alice", "delete"))
    fmt.Printf("Bob有写入权限: %v\n", system.HasPermission("bob", "write"))
    fmt.Printf("Charlie有写入权限: %v\n", system.HasPermission("charlie", "write"))
    
    // 获取用户所有权限
    alicePerms := system.GetUserPermissions("alice")
    fmt.Printf("Alice的所有权限: %v\n", alicePerms.ToSlice())
}
```

## 总结

集合论为计算机科学提供了重要的理论基础，通过形式化的定义和公理化系统，我们可以：

1. **严格定义数据结构** - 集合、关系、函数等基本概念
2. **证明算法正确性** - 使用数学方法验证程序逻辑
3. **优化系统设计** - 利用集合运算优化查询和计算
4. **构建类型系统** - 基于集合论构建安全的类型系统

Go语言的泛型特性使得我们能够实现类型安全的集合操作，同时保持代码的可读性和性能。通过结合数学理论和工程实践，我们可以构建更加可靠和高效的软件系统。

---

**相关链接**:

- [02-逻辑学](./02-Logic.md)
- [03-图论](./03-Graph-Theory.md)
- [04-概率论](./04-Probability-Theory.md)
- [返回数学基础层](../README.md)
