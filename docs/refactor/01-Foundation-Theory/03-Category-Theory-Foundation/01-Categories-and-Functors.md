# 01-范畴和函子 (Categories and Functors)

## 目录

- [01-范畴和函子](#01-范畴和函子)
  - [1. 概念定义](#1-概念定义)
  - [2. 形式化定义](#2-形式化定义)
  - [3. 定理证明](#3-定理证明)
  - [4. Go语言实现](#4-go语言实现)
  - [5. 应用示例](#5-应用示例)
  - [6. 性能分析](#6-性能分析)
  - [7. 参考文献](#7-参考文献)

## 1. 概念定义

### 1.1 基本概念

**定义 1.1**: 范畴 (Category)
一个范畴 $\mathcal{C}$ 由以下数据组成：
- 对象集合 $\text{Ob}(\mathcal{C})$
- 对于每对对象 $A, B \in \text{Ob}(\mathcal{C})$，有一个态射集合 $\text{Hom}_{\mathcal{C}}(A, B)$
- 对于每个对象 $A$，有一个单位态射 $1_A: A \rightarrow A$
- 对于态射 $f: A \rightarrow B$ 和 $g: B \rightarrow C$，有一个复合态射 $g \circ f: A \rightarrow C$

满足以下公理：
1. **结合律**: $(h \circ g) \circ f = h \circ (g \circ f)$
2. **单位律**: $f \circ 1_A = f = 1_B \circ f$

### 1.2 核心思想

范畴论的核心思想是通过态射（箭头）来研究对象之间的关系，而不是直接研究对象本身。这种抽象方法在计算机科学中特别有用，因为它可以统一处理各种不同的数学结构。

## 2. 形式化定义

### 2.1 数学定义

**定义 2.1**: 范畴的严格定义
一个范畴 $\mathcal{C}$ 是一个六元组 $(\text{Ob}, \text{Hom}, \text{dom}, \text{cod}, \text{id}, \circ)$，其中：

- $\text{Ob}$ 是对象集合
- $\text{Hom}$ 是态射集合
- $\text{dom}: \text{Hom} \rightarrow \text{Ob}$ 是定义域函数
- $\text{cod}: \text{Hom} \rightarrow \text{Ob}$ 是陪域函数
- $\text{id}: \text{Ob} \rightarrow \text{Hom}$ 是单位态射函数
- $\circ: \text{Hom} \times \text{Hom} \rightarrow \text{Hom}$ 是复合函数

满足以下条件：
1. 对于态射 $f, g$，$g \circ f$ 定义当且仅当 $\text{cod}(f) = \text{dom}(g)$
2. $\text{dom}(\text{id}(A)) = A = \text{cod}(\text{id}(A))$
3. 结合律和单位律

**定义 2.2**: 函子 (Functor)
从范畴 $\mathcal{C}$ 到范畴 $\mathcal{D}$ 的函子 $F: \mathcal{C} \rightarrow \mathcal{D}$ 由以下数据组成：
- 对象函数 $F_{\text{Ob}}: \text{Ob}(\mathcal{C}) \rightarrow \text{Ob}(\mathcal{D})$
- 态射函数 $F_{\text{Hom}}: \text{Hom}(\mathcal{C}) \rightarrow \text{Hom}(\mathcal{D})$

满足：
1. $F(1_A) = 1_{F(A)}$
2. $F(g \circ f) = F(g) \circ F(f)$

### 2.2 类型定义

```go
// Object 表示范畴中的对象
type Object interface {
    ID() string
    String() string
}

// Morphism 表示范畴中的态射
type Morphism interface {
    ID() string
    Domain() Object
    Codomain() Object
    Compose(other Morphism) (Morphism, error)
    String() string
}

// Category 表示范畴
type Category struct {
    Name     string
    Objects  map[string]Object
    Morphisms map[string]Morphism
}

// Functor 表示函子
type Functor struct {
    Name     string
    Source   *Category
    Target   *Category
    ObjectMap map[string]Object
    MorphismMap map[string]Morphism
}

// ConcreteObject 具体对象实现
type ConcreteObject struct {
    id   string
    name string
}

// ConcreteMorphism 具体态射实现
type ConcreteMorphism struct {
    id       string
    domain   Object
    codomain Object
    compose  func(Morphism) (Morphism, error)
}
```

## 3. 定理证明

### 3.1 定理陈述

**定理 3.1**: 单位态射的唯一性
在任意范畴中，每个对象的单位态射是唯一的。

**定理 3.2**: 函子的保结构性
函子保持范畴的所有结构，包括单位态射和复合运算。

**定理 3.3**: 函子复合
如果 $F: \mathcal{C} \rightarrow \mathcal{D}$ 和 $G: \mathcal{D} \rightarrow \mathcal{E}$ 是函子，则 $G \circ F: \mathcal{C} \rightarrow \mathcal{E}$ 也是函子。

### 3.2 证明过程

**证明定理 3.1**: 单位态射的唯一性

设 $1_A$ 和 $1'_A$ 都是对象 $A$ 的单位态射。

根据单位律：
- $1_A \circ 1'_A = 1'_A$ (因为 $1_A$ 是单位态射)
- $1_A \circ 1'_A = 1_A$ (因为 $1'_A$ 是单位态射)

因此 $1_A = 1'_A$，证毕。

**证明定理 3.2**: 函子的保结构性

设 $F: \mathcal{C} \rightarrow \mathcal{D}$ 是函子。

对于单位态射：
$$F(1_A) = 1_{F(A)}$$

对于复合运算：
$$F(g \circ f) = F(g) \circ F(f)$$

这些正是函子定义中的条件，证毕。

## 4. Go语言实现

### 4.1 基础实现

```go
package categorytheory

import (
    "fmt"
    "sync"
)

// Object 表示范畴中的对象
type Object interface {
    ID() string
    String() string
}

// Morphism 表示范畴中的态射
type Morphism interface {
    ID() string
    Domain() Object
    Codomain() Object
    Compose(other Morphism) (Morphism, error)
    String() string
}

// Category 表示范畴
type Category struct {
    Name       string
    Objects    map[string]Object
    Morphisms  map[string]Morphism
    mu         sync.RWMutex
}

// NewCategory 创建新范畴
func NewCategory(name string) *Category {
    return &Category{
        Name:      name,
        Objects:   make(map[string]Object),
        Morphisms: make(map[string]Morphism),
    }
}

// AddObject 添加对象
func (c *Category) AddObject(obj Object) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if _, exists := c.Objects[obj.ID()]; exists {
        return fmt.Errorf("object %s already exists", obj.ID())
    }
    
    c.Objects[obj.ID()] = obj
    return nil
}

// AddMorphism 添加态射
func (c *Category) AddMorphism(morph Morphism) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if _, exists := c.Morphisms[morph.ID()]; exists {
        return fmt.Errorf("morphism %s already exists", morph.ID())
    }
    
    // 检查定义域和陪域是否存在
    if _, exists := c.Objects[morph.Domain().ID()]; !exists {
        return fmt.Errorf("domain object %s does not exist", morph.Domain().ID())
    }
    
    if _, exists := c.Objects[morph.Codomain().ID()]; !exists {
        return fmt.Errorf("codomain object %s does not exist", morph.Codomain().ID())
    }
    
    c.Morphisms[morph.ID()] = morph
    return nil
}

// GetObject 获取对象
func (c *Category) GetObject(id string) (Object, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    obj, exists := c.Objects[id]
    return obj, exists
}

// GetMorphism 获取态射
func (c *Category) GetMorphism(id string) (Morphism, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    morph, exists := c.Morphisms[id]
    return morph, exists
}

// ComposeMorphisms 复合态射
func (c *Category) ComposeMorphisms(f, g Morphism) (Morphism, error) {
    if f.Codomain().ID() != g.Domain().ID() {
        return nil, fmt.Errorf("cannot compose morphisms: codomain of f != domain of g")
    }
    
    return f.Compose(g)
}

// ConcreteObject 具体对象实现
type ConcreteObject struct {
    id   string
    name string
}

// NewConcreteObject 创建具体对象
func NewConcreteObject(id, name string) *ConcreteObject {
    return &ConcreteObject{
        id:   id,
        name: name,
    }
}

// ID 返回对象ID
func (o *ConcreteObject) ID() string {
    return o.id
}

// String 字符串表示
func (o *ConcreteObject) String() string {
    return fmt.Sprintf("Object(%s: %s)", o.id, o.name)
}

// ConcreteMorphism 具体态射实现
type ConcreteMorphism struct {
    id       string
    domain   Object
    codomain Object
    compose  func(Morphism) (Morphism, error)
}

// NewConcreteMorphism 创建具体态射
func NewConcreteMorphism(id string, domain, codomain Object) *ConcreteMorphism {
    return &ConcreteMorphism{
        id:       id,
        domain:   domain,
        codomain: codomain,
    }
}

// ID 返回态射ID
func (m *ConcreteMorphism) ID() string {
    return m.id
}

// Domain 返回定义域
func (m *ConcreteMorphism) Domain() Object {
    return m.domain
}

// Codomain 返回陪域
func (m *ConcreteMorphism) Codomain() Object {
    return m.codomain
}

// Compose 复合态射
func (m *ConcreteMorphism) Compose(other Morphism) (Morphism, error) {
    if m.codomain.ID() != other.Domain().ID() {
        return nil, fmt.Errorf("cannot compose: codomain != domain")
    }
    
    // 创建复合态射
    compositeID := fmt.Sprintf("%s_%s", m.id, other.ID())
    composite := NewConcreteMorphism(compositeID, m.domain, other.Codomain())
    
    return composite, nil
}

// String 字符串表示
func (m *ConcreteMorphism) String() string {
    return fmt.Sprintf("Morphism(%s: %s -> %s)", m.id, m.domain.ID(), m.codomain.ID())
}

// Functor 表示函子
type Functor struct {
    Name         string
    Source       *Category
    Target       *Category
    ObjectMap    map[string]Object
    MorphismMap  map[string]Morphism
    mu           sync.RWMutex
}

// NewFunctor 创建新函子
func NewFunctor(name string, source, target *Category) *Functor {
    return &Functor{
        Name:        name,
        Source:      source,
        Target:      target,
        ObjectMap:   make(map[string]Object),
        MorphismMap: make(map[string]Morphism),
    }
}

// MapObject 映射对象
func (f *Functor) MapObject(sourceObj Object, targetObj Object) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    
    f.ObjectMap[sourceObj.ID()] = targetObj
    return nil
}

// MapMorphism 映射态射
func (f *Functor) MapMorphism(sourceMorph Morphism, targetMorph Morphism) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    
    f.MorphismMap[sourceMorph.ID()] = targetMorph
    return nil
}

// ApplyToObject 应用到对象
func (f *Functor) ApplyToObject(sourceObj Object) (Object, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    if targetObj, exists := f.ObjectMap[sourceObj.ID()]; exists {
        return targetObj, nil
    }
    
    return nil, fmt.Errorf("no mapping for object %s", sourceObj.ID())
}

// ApplyToMorphism 应用到态射
func (f *Functor) ApplyToMorphism(sourceMorph Morphism) (Morphism, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    if targetMorph, exists := f.MorphismMap[sourceMorph.ID()]; exists {
        return targetMorph, nil
    }
    
    return nil, fmt.Errorf("no mapping for morphism %s", sourceMorph.ID())
}
```

### 4.2 泛型实现

```go
// GenericCategory 泛型范畴
type GenericCategory[T Object, U Morphism] struct {
    Name      string
    Objects   map[string]T
    Morphisms map[string]U
    mu        sync.RWMutex
}

// GenericFunctor 泛型函子
type GenericFunctor[T1, T2 Object, U1, U2 Morphism] struct {
    Name         string
    Source       *GenericCategory[T1, U1]
    Target       *GenericCategory[T2, U2]
    ObjectMap    map[string]T2
    MorphismMap  map[string]U2
    mu           sync.RWMutex
}

// NewGenericCategory 创建泛型范畴
func NewGenericCategory[T Object, U Morphism](name string) *GenericCategory[T, U] {
    return &GenericCategory[T, U]{
        Name:      name,
        Objects:   make(map[string]T),
        Morphisms: make(map[string]U),
    }
}

// AddObject 添加对象
func (gc *GenericCategory[T, U]) AddObject(obj T) error {
    gc.mu.Lock()
    defer gc.mu.Unlock()
    
    if _, exists := gc.Objects[obj.ID()]; exists {
        return fmt.Errorf("object %s already exists", obj.ID())
    }
    
    gc.Objects[obj.ID()] = obj
    return nil
}
```

### 4.3 并发实现

```go
// ConcurrentCategory 并发范畴
type ConcurrentCategory struct {
    Name       string
    Objects    sync.Map
    Morphisms  sync.Map
}

// NewConcurrentCategory 创建并发范畴
func NewConcurrentCategory(name string) *ConcurrentCategory {
    return &ConcurrentCategory{
        Name: name,
    }
}

// AddObject 线程安全添加对象
func (cc *ConcurrentCategory) AddObject(obj Object) error {
    if _, loaded := cc.Objects.LoadOrStore(obj.ID(), obj); loaded {
        return fmt.Errorf("object %s already exists", obj.ID())
    }
    return nil
}

// GetObject 线程安全获取对象
func (cc *ConcurrentCategory) GetObject(id string) (Object, bool) {
    if value, ok := cc.Objects.Load(id); ok {
        return value.(Object), true
    }
    return nil, false
}

// RangeObjects 遍历所有对象
func (cc *ConcurrentCategory) RangeObjects(f func(key, value interface{}) bool) {
    cc.Objects.Range(f)
}
```

## 5. 应用示例

### 5.1 基础示例

```go
// 示例：集合范畴
func SetCategoryExample() {
    // 创建集合范畴
    setCat := NewCategory("Set")
    
    // 创建对象（集合）
    emptySet := NewConcreteObject("empty", "Empty Set")
    singletonSet := NewConcreteObject("singleton", "Singleton Set")
    twoElementSet := NewConcreteObject("two", "Two Element Set")
    
    // 添加对象到范畴
    setCat.AddObject(emptySet)
    setCat.AddObject(singletonSet)
    setCat.AddObject(twoElementSet)
    
    // 创建态射（函数）
    emptyToSingleton := NewConcreteMorphism("empty_to_singleton", emptySet, singletonSet)
    singletonToTwo := NewConcreteMorphism("singleton_to_two", singletonSet, twoElementSet)
    
    // 添加态射到范畴
    setCat.AddMorphism(emptyToSingleton)
    setCat.AddMorphism(singletonToTwo)
    
    // 复合态射
    composite, err := setCat.ComposeMorphisms(emptyToSingleton, singletonToTwo)
    if err == nil {
        fmt.Printf("Composite morphism: %s\n", composite)
    }
}
```

### 5.2 高级示例

```go
// 示例：函子应用
func FunctorExample() {
    // 创建源范畴和目标范畴
    sourceCat := NewCategory("Source")
    targetCat := NewCategory("Target")
    
    // 创建对象
    sourceObj := NewConcreteObject("A", "Source Object A")
    targetObj := NewConcreteObject("F(A)", "Target Object F(A)")
    
    sourceCat.AddObject(sourceObj)
    targetCat.AddObject(targetObj)
    
    // 创建函子
    functor := NewFunctor("F", sourceCat, targetCat)
    
    // 定义对象映射
    functor.MapObject(sourceObj, targetObj)
    
    // 应用函子
    result, err := functor.ApplyToObject(sourceObj)
    if err == nil {
        fmt.Printf("Functor applied: %s -> %s\n", sourceObj, result)
    }
}

// 示例：程序状态范畴
func ProgramStateCategoryExample() {
    // 创建程序状态范畴
    progCat := NewCategory("ProgramState")
    
    // 状态对象
    initialState := NewConcreteObject("init", "Initial State")
    runningState := NewConcreteObject("running", "Running State")
    finalState := NewConcreteObject("final", "Final State")
    
    progCat.AddObject(initialState)
    progCat.AddObject(runningState)
    progCat.AddObject(finalState)
    
    // 状态转换态射
    initToRunning := NewConcreteMorphism("start", initialState, runningState)
    runningToFinal := NewConcreteMorphism("finish", runningState, finalState)
    
    progCat.AddMorphism(initToRunning)
    progCat.AddMorphism(runningToFinal)
    
    // 计算完整执行路径
    fullExecution, err := progCat.ComposeMorphisms(initToRunning, runningToFinal)
    if err == nil {
        fmt.Printf("Full execution path: %s\n", fullExecution)
    }
}
```

## 6. 性能分析

### 6.1 时间复杂度

- **对象查找**: $O(1)$ (使用哈希表)
- **态射查找**: $O(1)$ (使用哈希表)
- **态射复合**: $O(1)$ (直接操作)
- **函子应用**: $O(1)$ (映射查找)

### 6.2 空间复杂度

- **范畴存储**: $O(|Ob| + |Mor|)$
- **函子存储**: $O(|Ob| + |Mor|)$
- **态射复合**: $O(1)$

### 6.3 基准测试

```go
func BenchmarkCategoryOperations(b *testing.B) {
    cat := NewCategory("Benchmark")
    
    // 预创建对象和态射
    objects := make([]Object, 1000)
    morphisms := make([]Morphism, 1000)
    
    for i := 0; i < 1000; i++ {
        obj := NewConcreteObject(fmt.Sprintf("obj%d", i), fmt.Sprintf("Object %d", i))
        objects[i] = obj
        cat.AddObject(obj)
        
        if i > 0 {
            morph := NewConcreteMorphism(fmt.Sprintf("morph%d", i), objects[i-1], obj)
            morphisms[i] = morph
            cat.AddMorphism(morph)
        }
    }
    
    b.ResetTimer()
    
    b.Run("ObjectLookup", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            cat.GetObject("obj500")
        }
    })
    
    b.Run("MorphismLookup", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            cat.GetMorphism("morph500")
        }
    })
    
    b.Run("Composition", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            if i < len(morphisms)-1 {
                cat.ComposeMorphisms(morphisms[i], morphisms[i+1])
            }
        }
    })
}
```

## 7. 参考文献

1. Mac Lane, S. (1998). *Categories for the Working Mathematician*. Springer-Verlag.
2. Awodey, S. (2010). *Category Theory*. Oxford University Press.
3. Barr, M., & Wells, C. (1995). *Category Theory for Computing Science*. Prentice Hall.
4. Pierce, B. C. (1991). *Basic Category Theory for Computer Scientists*. MIT Press.
5. Spivak, D. I. (2014). *Category Theory for the Sciences*. MIT Press.
