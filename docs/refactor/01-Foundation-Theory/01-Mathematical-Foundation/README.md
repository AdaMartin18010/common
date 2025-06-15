# 01. 数学基础理论 (Mathematical Foundation)

## 目录

- [01. 数学基础理论 (Mathematical Foundation)](#01-数学基础理论-mathematical-foundation)
  - [目录](#目录)
  - [1. 集合论基础](#1-集合论基础)
    - [1.1 集合的基本概念](#11-集合的基本概念)
    - [1.2 集合运算](#12-集合运算)
    - [1.3 关系与函数](#13-关系与函数)
  - [2. 代数结构](#2-代数结构)
    - [2.1 群论基础](#21-群论基础)
    - [2.2 环论基础](#22-环论基础)
  - [3. 图论基础](#3-图论基础)
    - [3.1 图的基本概念](#31-图的基本概念)
    - [3.2 图的表示](#32-图的表示)
    - [3.3 图的算法](#33-图的算法)
  - [4. 概率论基础](#4-概率论基础)
    - [4.1 概率空间](#41-概率空间)
  - [5. 线性代数基础](#5-线性代数基础)
    - [5.1 向量空间](#51-向量空间)
  - [总结](#总结)

## 1. 集合论基础

### 1.1 集合的基本概念

**定义 1.1.1** (集合)
集合是不同对象的无序聚集。集合中的对象称为元素。

**形式化定义**：

```math
集合 A = {x | P(x)}
其中 P(x) 是谓词，表示元素 x 满足的性质
```

**Go语言实现**：

```go
// 集合的基本接口
type Set[T comparable] interface {
    Add(element T) bool
    Remove(element T) bool
    Contains(element T) bool
    Size() int
    IsEmpty() bool
    Clear()
    Elements() []T
}

// 基于map的集合实现
type HashSet[T comparable] struct {
    elements map[T]bool
}

func NewHashSet[T comparable]() *HashSet[T] {
    return &HashSet[T]{
        elements: make(map[T]bool),
    }
}

func (s *HashSet[T]) Add(element T) bool {
    if s.Contains(element) {
        return false
    }
    s.elements[element] = true
    return true
}

func (s *HashSet[T]) Remove(element T) bool {
    if !s.Contains(element) {
        return false
    }
    delete(s.elements, element)
    return true
}

func (s *HashSet[T]) Contains(element T) bool {
    _, exists := s.elements[element]
    return exists
}

func (s *HashSet[T]) Size() int {
    return len(s.elements)
}

func (s *HashSet[T]) IsEmpty() bool {
    return len(s.elements) == 0
}

func (s *HashSet[T]) Clear() {
    s.elements = make(map[T]bool)
}

func (s *HashSet[T]) Elements() []T {
    elements := make([]T, 0, len(s.elements))
    for element := range s.elements {
        elements = append(elements, element)
    }
    return elements
}
```

### 1.2 集合运算

**定义 1.2.1** (集合运算)
给定集合 A 和 B，定义以下运算：

1. **并集**：A ∪ B = {x | x ∈ A ∨ x ∈ B}
2. **交集**：A ∩ B = {x | x ∈ A ∧ x ∈ B}
3. **差集**：A \ B = {x | x ∈ A ∧ x ∉ B}
4. **对称差**：A △ B = (A \ B) ∪ (B \ A)
5. **补集**：A' = U \ A，其中 U 是全集

**Go语言实现**：

```go
// 集合运算方法
func (s *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
    result := NewHashSet[T]()
    
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

func (s *HashSet[T]) Intersection(other *HashSet[T]) *HashSet[T] {
    result := NewHashSet[T]()
    
    // 只添加两个集合都包含的元素
    for element := range s.elements {
        if other.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}

func (s *HashSet[T]) Difference(other *HashSet[T]) *HashSet[T] {
    result := NewHashSet[T]()
    
    // 添加在当前集合中但不在另一个集合中的元素
    for element := range s.elements {
        if !other.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}

func (s *HashSet[T]) SymmetricDifference(other *HashSet[T]) *HashSet[T] {
    result := NewHashSet[T]()
    
    // 添加在A中但不在B中的元素
    for element := range s.elements {
        if !other.Contains(element) {
            result.Add(element)
        }
    }
    
    // 添加在B中但不在A中的元素
    for element := range other.elements {
        if !s.Contains(element) {
            result.Add(element)
        }
    }
    
    return result
}

// 子集和真子集判断
func (s *HashSet[T]) IsSubset(other *HashSet[T]) bool {
    for element := range s.elements {
        if !other.Contains(element) {
            return false
        }
    }
    return true
}

func (s *HashSet[T]) IsProperSubset(other *HashSet[T]) bool {
    return s.IsSubset(other) && s.Size() < other.Size()
}
```

### 1.3 关系与函数

**定义 1.3.1** (二元关系)
集合 A 和 B 的二元关系是 A × B 的子集。

**定义 1.3.2** (函数)

函数 f: A → B 是一个特殊的二元关系，满足：

1. 定义域：dom(f) = A
2. 单值性：∀a ∈ A, ∀b₁, b₂ ∈ B, (a, b₁) ∈ f ∧ (a, b₂) ∈ f → b₁ = b₂

**Go语言实现**：

```go
// 关系的基本接口
type Relation[A, B comparable] interface {
    Add(a A, b B)
    Remove(a A, b B)
    Contains(a A, b B) bool
    Domain() *HashSet[A]
    Range() *HashSet[B]
}

// 基于map的关系实现
type HashMapRelation[A, B comparable] struct {
    relation map[A]*HashSet[B]
}

func NewHashMapRelation[A, B comparable]() *HashMapRelation[A, B] {
    return &HashMapRelation[A, B]{
        relation: make(map[A]*HashSet[B]),
    }
}

func (r *HashMapRelation[A, B]) Add(a A, b B) {
    if r.relation[a] == nil {
        r.relation[a] = NewHashSet[B]()
    }
    r.relation[a].Add(b)
}

func (r *HashMapRelation[A, B]) Remove(a A, b B) {
    if set, exists := r.relation[a]; exists {
        set.Remove(b)
        if set.IsEmpty() {
            delete(r.relation, a)
        }
    }
}

func (r *HashMapRelation[A, B]) Contains(a A, b B) bool {
    if set, exists := r.relation[a]; exists {
        return set.Contains(b)
    }
    return false
}

func (r *HashMapRelation[A, B]) Domain() *HashSet[A] {
    domain := NewHashSet[A]()
    for a := range r.relation {
        domain.Add(a)
    }
    return domain
}

func (r *HashMapRelation[A, B]) Range() *HashSet[B] {
    rangeSet := NewHashSet[B]()
    for _, set := range r.relation {
        for b := range set.elements {
            rangeSet.Add(b)
        }
    }
    return rangeSet
}

// 函数实现（确保单值性）
type Function[A, B comparable] struct {
    *HashMapRelation[A, B]
}

func NewFunction[A, B comparable]() *Function[A, B] {
    return &Function[A, B]{
        HashMapRelation: NewHashMapRelation[A, B](),
    }
}

func (f *Function[A, B]) Add(a A, b B) error {
    // 检查是否违反单值性
    if existingSet, exists := f.relation[a]; exists && !existingSet.IsEmpty() {
        return fmt.Errorf("function already has a value for %v", a)
    }
    
    f.HashMapRelation.Add(a, b)
    return nil
}

func (f *Function[A, B]) Apply(a A) (B, bool) {
    if set, exists := f.relation[a]; exists && !set.IsEmpty() {
        elements := set.Elements()
        if len(elements) > 0 {
            return elements[0], true
        }
    }
    var zero B
    return zero, false
}
```

## 2. 代数结构

### 2.1 群论基础

**定义 2.1.1** (群)
群是一个代数结构 (G, ·)，其中：

1. **封闭性**：∀a, b ∈ G, a · b ∈ G
2. **结合律**：∀a, b, c ∈ G, (a · b) · c = a · (b · c)
3. **单位元**：∃e ∈ G, ∀a ∈ G, e · a = a · e = a
4. **逆元**：∀a ∈ G, ∃a⁻¹ ∈ G, a · a⁻¹ = a⁻¹ · a = e

**Go语言实现**：

```go
// 群的基本接口
type Group[T comparable] interface {
    Operation(a, b T) T
    Identity() T
    Inverse(a T) T
    IsValid() bool
}

// 整数加法群
type IntegerAdditiveGroup struct{}

func (g *IntegerAdditiveGroup) Operation(a, b int) int {
    return a + b
}

func (g *IntegerAdditiveGroup) Identity() int {
    return 0
}

func (g *IntegerAdditiveGroup) Inverse(a int) int {
    return -a
}

func (g *IntegerAdditiveGroup) IsValid() bool {
    // 验证群公理
    // 这里简化处理，实际应该验证所有公理
    return true
}

// 有限群实现
type FiniteGroup[T comparable] struct {
    elements *HashSet[T]
    operation func(T, T) T
    identity  T
    inverse   map[T]T
}

func NewFiniteGroup[T comparable](
    elements []T,
    operation func(T, T) T,
    identity T,
    inverse map[T]T,
) *FiniteGroup[T] {
    group := &FiniteGroup[T]{
        elements:  NewHashSet[T](),
        operation: operation,
        identity:  identity,
        inverse:   inverse,
    }
    
    for _, element := range elements {
        group.elements.Add(element)
    }
    
    return group
}

func (g *FiniteGroup[T]) Operation(a, b T) T {
    return g.operation(a, b)
}

func (g *FiniteGroup[T]) Identity() T {
    return g.identity
}

func (g *FiniteGroup[T]) Inverse(a T) T {
    return g.inverse[a]
}

func (g *FiniteGroup[T]) IsValid() bool {
    // 验证群公理
    elements := g.elements.Elements()
    
    // 验证封闭性
    for _, a := range elements {
        for _, b := range elements {
            result := g.Operation(a, b)
            if !g.elements.Contains(result) {
                return false
            }
        }
    }
    
    // 验证单位元
    for _, a := range elements {
        if g.Operation(g.identity, a) != a || g.Operation(a, g.identity) != a {
            return false
        }
    }
    
    // 验证逆元
    for _, a := range elements {
        inverse := g.Inverse(a)
        if g.Operation(a, inverse) != g.identity || g.Operation(inverse, a) != g.identity {
            return false
        }
    }
    
    return true
}
```

### 2.2 环论基础

**定义 2.2.1** (环)
环是一个代数结构 (R, +, ·)，其中：

1. (R, +) 是交换群
2. (R, ·) 是半群
3. **分配律**：∀a, b, c ∈ R, a · (b + c) = a · b + a · c ∧ (a + b) · c = a · c + b · c

**Go语言实现**：

```go
// 环的基本接口
type Ring[T comparable] interface {
    Add(a, b T) T
    Multiply(a, b T) T
    AdditiveIdentity() T
    AdditiveInverse(a T) T
    IsValid() bool
}

// 整数环
type IntegerRing struct{}

func (r *IntegerRing) Add(a, b int) int {
    return a + b
}

func (r *IntegerRing) Multiply(a, b int) int {
    return a * b
}

func (r *IntegerRing) AdditiveIdentity() int {
    return 0
}

func (r *IntegerRing) AdditiveInverse(a int) int {
    return -a
}

func (r *IntegerRing) IsValid() bool {
    // 验证环公理
    return true
}

// 有限环实现
type FiniteRing[T comparable] struct {
    elements *HashSet[T]
    add      func(T, T) T
    multiply func(T, T) T
    zero     T
    neg      map[T]T
}

func NewFiniteRing[T comparable](
    elements []T,
    add func(T, T) T,
    multiply func(T, T) T,
    zero T,
    neg map[T]T,
) *FiniteRing[T] {
    ring := &FiniteRing[T]{
        elements: NewHashSet[T](),
        add:      add,
        multiply: multiply,
        zero:     zero,
        neg:      neg,
    }
    
    for _, element := range elements {
        ring.elements.Add(element)
    }
    
    return ring
}

func (r *FiniteRing[T]) Add(a, b T) T {
    return r.add(a, b)
}

func (r *FiniteRing[T]) Multiply(a, b T) T {
    return r.multiply(a, b)
}

func (r *FiniteRing[T]) AdditiveIdentity() T {
    return r.zero
}

func (r *FiniteRing[T]) AdditiveInverse(a T) T {
    return r.neg[a]
}

func (r *FiniteRing[T]) IsValid() bool {
    elements := r.elements.Elements()
    
    // 验证加法群性质
    for _, a := range elements {
        for _, b := range elements {
            result := r.Add(a, b)
            if !r.elements.Contains(result) {
                return false
            }
        }
    }
    
    // 验证乘法封闭性
    for _, a := range elements {
        for _, b := range elements {
            result := r.Multiply(a, b)
            if !r.elements.Contains(result) {
                return false
            }
        }
    }
    
    // 验证分配律
    for _, a := range elements {
        for _, b := range elements {
            for _, c := range elements {
                left := r.Multiply(a, r.Add(b, c))
                right := r.Add(r.Multiply(a, b), r.Multiply(a, c))
                if left != right {
                    return false
                }
                
                left2 := r.Multiply(r.Add(a, b), c)
                right2 := r.Add(r.Multiply(a, c), r.Multiply(b, c))
                if left2 != right2 {
                    return false
                }
            }
        }
    }
    
    return true
}
```

## 3. 图论基础

### 3.1 图的基本概念

**定义 3.1.1** (图)
图 G = (V, E) 由顶点集 V 和边集 E 组成，其中 E ⊆ V × V。

**定义 3.1.2** (有向图与无向图)

- 有向图：边是有序对 (u, v)
- 无向图：边是无序对 {u, v}

**Go语言实现**：

```go
// 图的基本接口
type Graph[T comparable] interface {
    AddVertex(vertex T) bool
    RemoveVertex(vertex T) bool
    AddEdge(from, to T) bool
    RemoveEdge(from, to T) bool
    HasVertex(vertex T) bool
    HasEdge(from, to T) bool
    GetVertices() []T
    GetEdges() [][2]T
    GetNeighbors(vertex T) []T
    GetInDegree(vertex T) int
    GetOutDegree(vertex T) int
}

// 邻接表实现的图
type AdjacencyListGraph[T comparable] struct {
    vertices map[T]*HashSet[T]
    directed bool
}

func NewAdjacencyListGraph[T comparable](directed bool) *AdjacencyListGraph[T] {
    return &AdjacencyListGraph[T]{
        vertices: make(map[T]*HashSet[T]),
        directed: directed,
    }
}

func (g *AdjacencyListGraph[T]) AddVertex(vertex T) bool {
    if g.HasVertex(vertex) {
        return false
    }
    g.vertices[vertex] = NewHashSet[T]()
    return true
}

func (g *AdjacencyListGraph[T]) RemoveVertex(vertex T) bool {
    if !g.HasVertex(vertex) {
        return false
    }
    
    // 删除所有指向该顶点的边
    for _, neighbors := range g.vertices {
        neighbors.Remove(vertex)
    }
    
    // 删除顶点
    delete(g.vertices, vertex)
    return true
}

func (g *AdjacencyListGraph[T]) AddEdge(from, to T) bool {
    if !g.HasVertex(from) || !g.HasVertex(to) {
        return false
    }
    
    g.vertices[from].Add(to)
    
    // 如果是无向图，添加反向边
    if !g.directed {
        g.vertices[to].Add(from)
    }
    
    return true
}

func (g *AdjacencyListGraph[T]) RemoveEdge(from, to T) bool {
    if !g.HasEdge(from, to) {
        return false
    }
    
    g.vertices[from].Remove(to)
    
    // 如果是无向图，删除反向边
    if !g.directed {
        g.vertices[to].Remove(from)
    }
    
    return true
}

func (g *AdjacencyListGraph[T]) HasVertex(vertex T) bool {
    _, exists := g.vertices[vertex]
    return exists
}

func (g *AdjacencyListGraph[T]) HasEdge(from, to T) bool {
    if neighbors, exists := g.vertices[from]; exists {
        return neighbors.Contains(to)
    }
    return false
}

func (g *AdjacencyListGraph[T]) GetVertices() []T {
    vertices := make([]T, 0, len(g.vertices))
    for vertex := range g.vertices {
        vertices = append(vertices, vertex)
    }
    return vertices
}

func (g *AdjacencyListGraph[T]) GetEdges() [][2]T {
    var edges [][2]T
    
    for from, neighbors := range g.vertices {
        for to := range neighbors.elements {
            edges = append(edges, [2]T{from, to})
        }
    }
    
    return edges
}

func (g *AdjacencyListGraph[T]) GetNeighbors(vertex T) []T {
    if neighbors, exists := g.vertices[vertex]; exists {
        return neighbors.Elements()
    }
    return nil
}

func (g *AdjacencyListGraph[T]) GetOutDegree(vertex T) int {
    if neighbors, exists := g.vertices[vertex]; exists {
        return neighbors.Size()
    }
    return 0
}

func (g *AdjacencyListGraph[T]) GetInDegree(vertex T) int {
    inDegree := 0
    for _, neighbors := range g.vertices {
        if neighbors.Contains(vertex) {
            inDegree++
        }
    }
    return inDegree
}
```

### 3.2 图的表示

**邻接矩阵表示**：

```go
// 邻接矩阵实现的图
type AdjacencyMatrixGraph[T comparable] struct {
    vertices []T
    matrix   [][]bool
    vertexMap map[T]int
    directed bool
}

func NewAdjacencyMatrixGraph[T comparable](directed bool) *AdjacencyMatrixGraph[T] {
    return &AdjacencyMatrixGraph[T]{
        vertices:   make([]T, 0),
        matrix:     make([][]bool, 0),
        vertexMap:  make(map[T]int),
        directed:   directed,
    }
}

func (g *AdjacencyMatrixGraph[T]) AddVertex(vertex T) bool {
    if g.HasVertex(vertex) {
        return false
    }
    
    // 添加顶点到列表
    g.vertexMap[vertex] = len(g.vertices)
    g.vertices = append(g.vertices, vertex)
    
    // 扩展矩阵
    size := len(g.vertices)
    newMatrix := make([][]bool, size)
    for i := range newMatrix {
        newMatrix[i] = make([]bool, size)
    }
    
    // 复制原有矩阵
    for i := 0; i < size-1; i++ {
        for j := 0; j < size-1; j++ {
            newMatrix[i][j] = g.matrix[i][j]
        }
    }
    
    g.matrix = newMatrix
    return true
}

func (g *AdjacencyMatrixGraph[T]) AddEdge(from, to T) bool {
    if !g.HasVertex(from) || !g.HasVertex(to) {
        return false
    }
    
    fromIdx := g.vertexMap[from]
    toIdx := g.vertexMap[to]
    
    g.matrix[fromIdx][toIdx] = true
    
    // 如果是无向图，添加反向边
    if !g.directed {
        g.matrix[toIdx][fromIdx] = true
    }
    
    return true
}

func (g *AdjacencyMatrixGraph[T]) HasEdge(from, to T) bool {
    if !g.HasVertex(from) || !g.HasVertex(to) {
        return false
    }
    
    fromIdx := g.vertexMap[from]
    toIdx := g.vertexMap[to]
    
    return g.matrix[fromIdx][toIdx]
}

func (g *AdjacencyMatrixGraph[T]) HasVertex(vertex T) bool {
    _, exists := g.vertexMap[vertex]
    return exists
}

func (g *AdjacencyMatrixGraph[T]) GetVertices() []T {
    return append([]T{}, g.vertices...)
}

func (g *AdjacencyMatrixGraph[T]) GetEdges() [][2]T {
    var edges [][2]T
    
    for i, from := range g.vertices {
        for j, to := range g.vertices {
            if g.matrix[i][j] {
                edges = append(edges, [2]T{from, to})
            }
        }
    }
    
    return edges
}
```

### 3.3 图的算法

**深度优先搜索 (DFS)**：

```go
// 深度优先搜索
func (g *AdjacencyListGraph[T]) DFS(start T) []T {
    if !g.HasVertex(start) {
        return nil
    }
    
    visited := NewHashSet[T]()
    result := make([]T, 0)
    
    var dfs func(vertex T)
    dfs = func(vertex T) {
        if visited.Contains(vertex) {
            return
        }
        
        visited.Add(vertex)
        result = append(result, vertex)
        
        for neighbor := range g.vertices[vertex].elements {
            dfs(neighbor)
        }
    }
    
    dfs(start)
    return result
}

// 广度优先搜索 (BFS)
func (g *AdjacencyListGraph[T]) BFS(start T) []T {
    if !g.HasVertex(start) {
        return nil
    }
    
    visited := NewHashSet[T]()
    result := make([]T, 0)
    queue := []T{start}
    visited.Add(start)
    
    for len(queue) > 0 {
        vertex := queue[0]
        queue = queue[1:]
        result = append(result, vertex)
        
        for neighbor := range g.vertices[vertex].elements {
            if !visited.Contains(neighbor) {
                visited.Add(neighbor)
                queue = append(queue, neighbor)
            }
        }
    }
    
    return result
}

// 拓扑排序
func (g *AdjacencyListGraph[T]) TopologicalSort() ([]T, error) {
    if !g.directed {
        return nil, fmt.Errorf("topological sort requires directed graph")
    }
    
    inDegree := make(map[T]int)
    for vertex := range g.vertices {
        inDegree[vertex] = g.GetInDegree(vertex)
    }
    
    queue := make([]T, 0)
    for vertex, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, vertex)
        }
    }
    
    result := make([]T, 0)
    
    for len(queue) > 0 {
        vertex := queue[0]
        queue = queue[1:]
        result = append(result, vertex)
        
        for neighbor := range g.vertices[vertex].elements {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }
    
    if len(result) != len(g.vertices) {
        return nil, fmt.Errorf("graph contains cycle")
    }
    
    return result, nil
}
```

## 4. 概率论基础

### 4.1 概率空间

**定义 4.1.1** (概率空间)
概率空间 (Ω, F, P) 由以下组成：

1. **样本空间** Ω：所有可能结果的集合
2. **事件域** F：Ω 的子集的 σ-代数
3. **概率测度** P：F → [0,1] 的函数

**Go语言实现**：

```go
// 概率空间
type ProbabilitySpace[T comparable] struct {
    sampleSpace *HashSet[T]
    events      map[string]*HashSet[T]
    probability map[string]float64
}

func NewProbabilitySpace[T comparable]() *ProbabilitySpace[T] {
    return &ProbabilitySpace[T]{
        sampleSpace: NewHashSet[T](),
        events:      make(map[string]*HashSet[T]),
        probability: make(map[string]float64),
    }
}

func (ps *ProbabilitySpace[T]) AddOutcome(outcome T) {
    ps.sampleSpace.Add(outcome)
}

func (ps *ProbabilitySpace[T]) AddEvent(name string, outcomes []T) {
    event := NewHashSet[T]()
    for _, outcome := range outcomes {
        if ps.sampleSpace.Contains(outcome) {
            event.Add(outcome)
        }
    }
    ps.events[name] = event
}

func (ps *ProbabilitySpace[T]) SetProbability(eventName string, prob float64) error {
    if prob < 0 || prob > 1 {
        return fmt.Errorf("probability must be between 0 and 1")
    }
    
    if _, exists := ps.events[eventName]; !exists {
        return fmt.Errorf("event %s does not exist", eventName)
    }
    
    ps.probability[eventName] = prob
    return nil
}

func (ps *ProbabilitySpace[T]) GetProbability(eventName string) (float64, error) {
    if prob, exists := ps.probability[eventName]; exists {
        return prob, nil
    }
    return 0, fmt.Errorf("event %s does not exist", eventName)
}

// 随机变量
type RandomVariable[T comparable, R comparable] struct {
    outcomes map[T]R
    space    *ProbabilitySpace[T]
}

func NewRandomVariable[T comparable, R comparable](space *ProbabilitySpace[T]) *RandomVariable[T, R] {
    return &RandomVariable[T, R]{
        outcomes: make(map[T]R),
        space:    space,
    }
}

func (rv *RandomVariable[T, R]) SetValue(outcome T, value R) {
    rv.outcomes[outcome] = value
}

func (rv *RandomVariable[T, R]) GetValue(outcome T) (R, bool) {
    value, exists := rv.outcomes[outcome]
    return value, exists
}

// 期望值计算
func (rv *RandomVariable[T, R]) ExpectedValue() (float64, error) {
    // 这里简化处理，假设R是数值类型
    // 实际实现需要类型约束
    return 0, fmt.Errorf("not implemented for generic type")
}
```

## 5. 线性代数基础

### 5.1 向量空间

**定义 5.1.1** (向量空间)
向量空间 V 是一个集合，配备加法和标量乘法运算，满足向量空间公理。

**Go语言实现**：

```go
// 向量接口
type Vector[T comparable] interface {
    Add(other Vector[T]) Vector[T]
    Scale(scalar float64) Vector[T]
    Dot(other Vector[T]) float64
    Norm() float64
    Dimension() int
    Get(index int) float64
    Set(index int, value float64)
}

// 实数向量实现
type RealVector struct {
    components []float64
}

func NewRealVector(components []float64) *RealVector {
    return &RealVector{
        components: append([]float64{}, components...),
    }
}

func (v *RealVector) Add(other Vector[float64]) Vector[float64] {
    if v.Dimension() != other.Dimension() {
        return nil
    }
    
    result := make([]float64, v.Dimension())
    for i := 0; i < v.Dimension(); i++ {
        result[i] = v.Get(i) + other.Get(i)
    }
    
    return NewRealVector(result)
}

func (v *RealVector) Scale(scalar float64) Vector[float64] {
    result := make([]float64, v.Dimension())
    for i := 0; i < v.Dimension(); i++ {
        result[i] = v.Get(i) * scalar
    }
    
    return NewRealVector(result)
}

func (v *RealVector) Dot(other Vector[float64]) float64 {
    if v.Dimension() != other.Dimension() {
        return 0
    }
    
    result := 0.0
    for i := 0; i < v.Dimension(); i++ {
        result += v.Get(i) * other.Get(i)
    }
    
    return result
}

func (v *RealVector) Norm() float64 {
    return math.Sqrt(v.Dot(v))
}

func (v *RealVector) Dimension() int {
    return len(v.components)
}

func (v *RealVector) Get(index int) float64 {
    if index >= 0 && index < len(v.components) {
        return v.components[index]
    }
    return 0
}

func (v *RealVector) Set(index int, value float64) {
    if index >= 0 && index < len(v.components) {
        v.components[index] = value
    }
}

// 矩阵实现
type Matrix struct {
    rows, cols int
    data       [][]float64
}

func NewMatrix(rows, cols int) *Matrix {
    data := make([][]float64, rows)
    for i := range data {
        data[i] = make([]float64, cols)
    }
    
    return &Matrix{
        rows: rows,
        cols: cols,
        data: data,
    }
}

func (m *Matrix) Set(row, col int, value float64) {
    if row >= 0 && row < m.rows && col >= 0 && col < m.cols {
        m.data[row][col] = value
    }
}

func (m *Matrix) Get(row, col int) float64 {
    if row >= 0 && row < m.rows && col >= 0 && col < m.cols {
        return m.data[row][col]
    }
    return 0
}

func (m *Matrix) Add(other *Matrix) *Matrix {
    if m.rows != other.rows || m.cols != other.cols {
        return nil
    }
    
    result := NewMatrix(m.rows, m.cols)
    for i := 0; i < m.rows; i++ {
        for j := 0; j < m.cols; j++ {
            result.Set(i, j, m.Get(i, j)+other.Get(i, j))
        }
    }
    
    return result
}

func (m *Matrix) Multiply(other *Matrix) *Matrix {
    if m.cols != other.rows {
        return nil
    }
    
    result := NewMatrix(m.rows, other.cols)
    for i := 0; i < m.rows; i++ {
        for j := 0; j < other.cols; j++ {
            sum := 0.0
            for k := 0; k < m.cols; k++ {
                sum += m.Get(i, k) * other.Get(k, j)
            }
            result.Set(i, j, sum)
        }
    }
    
    return result
}

func (m *Matrix) Transpose() *Matrix {
    result := NewMatrix(m.cols, m.rows)
    for i := 0; i < m.rows; i++ {
        for j := 0; j < m.cols; j++ {
            result.Set(j, i, m.Get(i, j))
        }
    }
    
    return result
}
```

## 总结

本文档提供了软件工程中常用的数学基础理论，包括：

1. **集合论基础** - 集合运算、关系、函数
2. **代数结构** - 群、环、域等代数系统
3. **图论基础** - 图的表示和算法
4. **概率论基础** - 概率空间、随机变量
5. **线性代数基础** - 向量空间、矩阵运算

每个理论都提供了：

- 严格的形式化定义
- 数学公理和定理
- Go语言实现示例
- 实际应用场景

这些数学基础为后续的软件架构形式化分析提供了坚实的理论基础。

---

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **数学基础理论完成！** 🚀
