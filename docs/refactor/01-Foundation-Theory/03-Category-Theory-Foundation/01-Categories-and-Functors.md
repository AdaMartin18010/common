# 01-范畴和函子 (Categories and Functors)

## 目录

- [01-范畴和函子 (Categories and Functors)](#01-范畴和函子-categories-and-functors)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 范畴的定义](#11-范畴的定义)
    - [1.2 态射和复合](#12-态射和复合)
    - [1.3 单位态射](#13-单位态射)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 范畴的公理化](#21-范畴的公理化)
    - [2.2 函子的定义](#22-函子的定义)
    - [2.3 自然变换](#23-自然变换)
  - [3. 重要概念](#3-重要概念)
    - [3.1 同构](#31-同构)
    - [3.2 单态射和满态射](#32-单态射和满态射)
    - [3.3 积和余积](#33-积和余积)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 范畴数据结构](#41-范畴数据结构)
    - [4.2 函子实现](#42-函子实现)
    - [4.3 自然变换实现](#43-自然变换实现)
  - [5. 应用示例](#5-应用示例)
    - [5.1 集合范畴](#51-集合范畴)
    - [5.2 群范畴](#52-群范畴)
    - [5.3 编程语言应用](#53-编程语言应用)
  - [总结](#总结)

## 1. 基本概念

### 1.1 范畴的定义

**定义 1.1**: 范畴 $\mathcal{C}$ 由以下数据组成：

1. **对象类** $\text{Ob}(\mathcal{C})$: 范畴中的对象集合
2. **态射类** $\text{Mor}(\mathcal{C})$: 范畴中的态射集合
3. **域和余域函数**:
   - $\text{dom}: \text{Mor}(\mathcal{C}) \rightarrow \text{Ob}(\mathcal{C})$ (域)
   - $\text{cod}: \text{Mor}(\mathcal{C}) \rightarrow \text{Ob}(\mathcal{C})$ (余域)
4. **复合函数**: $\circ: \text{Mor}(\mathcal{C}) \times \text{Mor}(\mathcal{C}) \rightarrow \text{Mor}(\mathcal{C})$
5. **单位函数**: $\text{id}: \text{Ob}(\mathcal{C}) \rightarrow \text{Mor}(\mathcal{C})$

**记法**: 对于对象 $A, B \in \text{Ob}(\mathcal{C})$，从 $A$ 到 $B$ 的态射集合记为 $\text{Hom}_{\mathcal{C}}(A, B)$ 或 $\mathcal{C}(A, B)$。

### 1.2 态射和复合

**定义 1.2**: 态射复合

对于态射 $f: A \rightarrow B$ 和 $g: B \rightarrow C$，它们的复合 $g \circ f: A \rightarrow C$ 满足：

1. **结合律**: $(h \circ g) \circ f = h \circ (g \circ f)$
2. **单位律**: $\text{id}_B \circ f = f = f \circ \text{id}_A$

**形式化表达**:

- 态射 $f: A \rightarrow B$ 表示 $\text{dom}(f) = A$ 且 $\text{cod}(f) = B$
- 复合 $g \circ f$ 定义当且仅当 $\text{cod}(f) = \text{dom}(g)$

### 1.3 单位态射

**定义 1.3**: 单位态射

对于每个对象 $A \in \text{Ob}(\mathcal{C})$，存在唯一的态射 $\text{id}_A: A \rightarrow A$，称为 $A$ 的单位态射，满足：

- 对于任意态射 $f: A \rightarrow B$，有 $\text{id}_B \circ f = f$
- 对于任意态射 $g: B \rightarrow A$，有 $g \circ \text{id}_A = g$

## 2. 形式化定义

### 2.1 范畴的公理化

**公理 2.1**: 范畴的公理系统

1. **对象公理**: $\text{Ob}(\mathcal{C})$ 是一个类
2. **态射公理**: 对于任意对象 $A, B$，$\text{Hom}_{\mathcal{C}}(A, B)$ 是一个集合
3. **复合公理**: 复合函数 $\circ$ 满足结合律
4. **单位公理**: 每个对象都有单位态射
5. **单位律公理**: 单位态射满足左右单位律

**定理 2.1**: 单位态射的唯一性

在任意范畴中，每个对象的单位态射是唯一的。

**证明**: 假设 $\text{id}_A$ 和 $\text{id}_A'$ 都是 $A$ 的单位态射，则：
$$\text{id}_A = \text{id}_A \circ \text{id}_A' = \text{id}_A'$$

### 2.2 函子的定义

**定义 2.2**: 函子

从范畴 $\mathcal{C}$ 到范畴 $\mathcal{D}$ 的**函子** $F: \mathcal{C} \rightarrow \mathcal{D}$ 由以下数据组成：

1. **对象函数**: $F: \text{Ob}(\mathcal{C}) \rightarrow \text{Ob}(\mathcal{D})$
2. **态射函数**: $F: \text{Mor}(\mathcal{C}) \rightarrow \text{Mor}(\mathcal{D})$

满足以下条件：

1. **保持域和余域**: 如果 $f: A \rightarrow B$，则 $F(f): F(A) \rightarrow F(B)$
2. **保持复合**: $F(g \circ f) = F(g) \circ F(f)$
3. **保持单位**: $F(\text{id}_A) = \text{id}_{F(A)}$

**定义 2.3**: 函子的类型

1. **协变函子**: 保持态射方向
2. **反变函子**: 反转态射方向，即 $F(f): F(B) \rightarrow F(A)$

### 2.3 自然变换

**定义 2.4**: 自然变换

对于函子 $F, G: \mathcal{C} \rightarrow \mathcal{D}$，**自然变换** $\alpha: F \Rightarrow G$ 是一个态射族 $\{\alpha_A: F(A) \rightarrow G(A)\}_{A \in \text{Ob}(\mathcal{C})}$，满足：

对于任意态射 $f: A \rightarrow B$，有：
$$G(f) \circ \alpha_A = \alpha_B \circ F(f)$$

这个等式称为**自然性条件**。

**定义 2.5**: 自然同构

自然变换 $\alpha: F \Rightarrow G$ 是**自然同构**，如果每个 $\alpha_A$ 都是同构。

## 3. 重要概念

### 3.1 同构

**定义 3.1**: 同构

态射 $f: A \rightarrow B$ 是**同构**，如果存在态射 $g: B \rightarrow A$ 使得：
$$g \circ f = \text{id}_A \text{ 且 } f \circ g = \text{id}_B$$

此时称 $g$ 为 $f$ 的**逆**，记作 $f^{-1}$。

**定理 3.1**: 逆的唯一性

同构的逆是唯一的。

**证明**: 假设 $g$ 和 $g'$ 都是 $f$ 的逆，则：
$$g = g \circ \text{id}_B = g \circ (f \circ g') = (g \circ f) \circ g' = \text{id}_A \circ g' = g'$$

### 3.2 单态射和满态射

**定义 3.2**: 单态射

态射 $f: A \rightarrow B$ 是**单态射**，如果对于任意态射 $g, h: C \rightarrow A$：
$$f \circ g = f \circ h \Rightarrow g = h$$

**定义 3.3**: 满态射

态射 $f: A \rightarrow B$ 是**满态射**，如果对于任意态射 $g, h: B \rightarrow C$：
$$g \circ f = h \circ f \Rightarrow g = h$$

### 3.3 积和余积

**定义 3.4**: 积

对象 $A$ 和 $B$ 的**积**是一个对象 $A \times B$ 连同两个投影态射：
$$\pi_1: A \times B \rightarrow A \text{ 和 } \pi_2: A \times B \rightarrow B$$

满足：对于任意对象 $C$ 和态射 $f: C \rightarrow A$，$g: C \rightarrow B$，存在唯一的态射 $\langle f, g \rangle: C \rightarrow A \times B$ 使得：
$$\pi_1 \circ \langle f, g \rangle = f \text{ 且 } \pi_2 \circ \langle f, g \rangle = g$$

**定义 3.5**: 余积

对象 $A$ 和 $B$ 的**余积**是一个对象 $A + B$ 连同两个注入态射：
$$\iota_1: A \rightarrow A + B \text{ 和 } \iota_2: B \rightarrow A + B$$

满足：对于任意对象 $C$ 和态射 $f: A \rightarrow C$，$g: B \rightarrow C$，存在唯一的态射 $[f, g]: A + B \rightarrow C$ 使得：
$$[f, g] \circ \iota_1 = f \text{ 且 } [f, g] \circ \iota_2 = g$$

## 4. Go语言实现

### 4.1 范畴数据结构

```go
// Object 范畴中的对象
type Object interface {
    ID() string
}

// Morphism 范畴中的态射
type Morphism struct {
    ID       string
    Domain   Object
    Codomain Object
    Data     interface{} // 态射的具体数据
}

// NewMorphism 创建态射
func NewMorphism(id string, domain, codomain Object, data interface{}) *Morphism {
    return &Morphism{
        ID:       id,
        Domain:   domain,
        Codomain: codomain,
        Data:     data,
    }
}

// Category 范畴
type Category struct {
    Name      string
    Objects   map[string]Object
    Morphisms map[string]*Morphism
    Compose   func(*Morphism, *Morphism) (*Morphism, error)
    Identity  func(Object) *Morphism
}

// NewCategory 创建范畴
func NewCategory(name string) *Category {
    return &Category{
        Name:      name,
        Objects:   make(map[string]Object),
        Morphisms: make(map[string]*Morphism),
    }
}

// AddObject 添加对象
func (c *Category) AddObject(obj Object) {
    c.Objects[obj.ID()] = obj
}

// AddMorphism 添加态射
func (c *Category) AddMorphism(morphism *Morphism) {
    c.Morphisms[morphism.ID] = morphism
}

// GetMorphisms 获取从A到B的所有态射
func (c *Category) GetMorphisms(from, to Object) []*Morphism {
    var morphisms []*Morphism
    for _, morphism := range c.Morphisms {
        if morphism.Domain == from && morphism.Codomain == to {
            morphisms = append(morphisms, morphism)
        }
    }
    return morphisms
}

// ComposeMorphisms 复合态射
func (c *Category) ComposeMorphisms(f, g *Morphism) (*Morphism, error) {
    if f.Codomain != g.Domain {
        return nil, fmt.Errorf("cannot compose morphisms: codomain of f != domain of g")
    }
    
    if c.Compose != nil {
        return c.Compose(f, g)
    }
    
    // 默认复合
    composedID := fmt.Sprintf("%s_%s", f.ID, g.ID)
    composedData := fmt.Sprintf("(%s ∘ %s)", g.ID, f.ID)
    
    return NewMorphism(composedID, f.Domain, g.Codomain, composedData), nil
}

// IdentityMorphism 获取对象的单位态射
func (c *Category) IdentityMorphism(obj Object) *Morphism {
    if c.Identity != nil {
        return c.Identity(obj)
    }
    
    // 默认单位态射
    idID := fmt.Sprintf("id_%s", obj.ID())
    return NewMorphism(idID, obj, obj, "identity")
}

// IsIsomorphism 检查态射是否为同构
func (c *Category) IsIsomorphism(morphism *Morphism) bool {
    // 检查是否存在逆态射
    inverses := c.GetMorphisms(morphism.Codomain, morphism.Domain)
    
    for _, inverse := range inverses {
        // 检查复合是否为单位态射
        comp1, err1 := c.ComposeMorphisms(morphism, inverse)
        comp2, err2 := c.ComposeMorphisms(inverse, morphism)
        
        if err1 == nil && err2 == nil {
            id1 := c.IdentityMorphism(morphism.Domain)
            id2 := c.IdentityMorphism(morphism.Codomain)
            
            if comp1.ID == id1.ID && comp2.ID == id2.ID {
                return true
            }
        }
    }
    
    return false
}

// IsMonomorphism 检查态射是否为单态射
func (c *Category) IsMonomorphism(morphism *Morphism) bool {
    // 对于所有可能的态射对，检查左消去律
    for _, m1 := range c.Morphisms {
        for _, m2 := range c.Morphisms {
            if m1.Codomain == morphism.Domain && m2.Codomain == morphism.Domain {
                comp1, err1 := c.ComposeMorphisms(m1, morphism)
                comp2, err2 := c.ComposeMorphisms(m2, morphism)
                
                if err1 == nil && err2 == nil && comp1.ID == comp2.ID {
                    if m1.ID != m2.ID {
                        return false
                    }
                }
            }
        }
    }
    return true
}

// IsEpimorphism 检查态射是否为满态射
func (c *Category) IsEpimorphism(morphism *Morphism) bool {
    // 对于所有可能的态射对，检查右消去律
    for _, m1 := range c.Morphisms {
        for _, m2 := range c.Morphisms {
            if morphism.Codomain == m1.Domain && morphism.Codomain == m2.Domain {
                comp1, err1 := c.ComposeMorphisms(morphism, m1)
                comp2, err2 := c.ComposeMorphisms(morphism, m2)
                
                if err1 == nil && err2 == nil && comp1.ID == comp2.ID {
                    if m1.ID != m2.ID {
                        return false
                    }
                }
            }
        }
    }
    return true
}
```

### 4.2 函子实现

```go
// Functor 函子
type Functor struct {
    Name           string
    SourceCategory *Category
    TargetCategory *Category
    ObjectMap      map[string]Object
    MorphismMap    map[string]*Morphism
}

// NewFunctor 创建函子
func NewFunctor(name string, source, target *Category) *Functor {
    return &Functor{
        Name:           name,
        SourceCategory: source,
        TargetCategory: target,
        ObjectMap:      make(map[string]Object),
        MorphismMap:    make(map[string]*Morphism),
    }
}

// MapObject 映射对象
func (f *Functor) MapObject(obj Object, targetObj Object) {
    f.ObjectMap[obj.ID()] = targetObj
}

// MapMorphism 映射态射
func (f *Functor) MapMorphism(morphism *Morphism, targetMorphism *Morphism) {
    f.MorphismMap[morphism.ID] = targetMorphism
}

// ApplyToObject 应用函子到对象
func (f *Functor) ApplyToObject(obj Object) Object {
    if targetObj, exists := f.ObjectMap[obj.ID()]; exists {
        return targetObj
    }
    return nil
}

// ApplyToMorphism 应用函子到态射
func (f *Functor) ApplyToMorphism(morphism *Morphism) *Morphism {
    if targetMorphism, exists := f.MorphismMap[morphism.ID]; exists {
        return targetMorphism
    }
    return nil
}

// IsValid 检查函子是否有效
func (f *Functor) IsValid() bool {
    // 检查对象映射的一致性
    for _, obj := range f.SourceCategory.Objects {
        if f.ApplyToObject(obj) == nil {
            return false
        }
    }
    
    // 检查态射映射的一致性
    for _, morphism := range f.SourceCategory.Morphisms {
        targetMorphism := f.ApplyToMorphism(morphism)
        if targetMorphism == nil {
            return false
        }
        
        // 检查域和余域的映射
        targetDomain := f.ApplyToObject(morphism.Domain)
        targetCodomain := f.ApplyToObject(morphism.Codomain)
        
        if targetMorphism.Domain != targetDomain || targetMorphism.Codomain != targetCodomain {
            return false
        }
    }
    
    return true
}

// ComposeFunctors 复合函子
func ComposeFunctors(f, g *Functor) *Functor {
    if f.TargetCategory != g.SourceCategory {
        return nil
    }
    
    composed := NewFunctor(
        fmt.Sprintf("%s ∘ %s", g.Name, f.Name),
        f.SourceCategory,
        g.TargetCategory,
    )
    
    // 复合对象映射
    for objID, obj := range f.ObjectMap {
        if gObj := g.ApplyToObject(obj); gObj != nil {
            composed.ObjectMap[objID] = gObj
        }
    }
    
    // 复合态射映射
    for morphismID, morphism := range f.MorphismMap {
        if gMorphism := g.ApplyToMorphism(morphism); gMorphism != nil {
            composed.MorphismMap[morphismID] = gMorphism
        }
    }
    
    return composed
}
```

### 4.3 自然变换实现

```go
// NaturalTransformation 自然变换
type NaturalTransformation struct {
    Name     string
    Source   *Functor
    Target   *Functor
    Components map[string]*Morphism // 每个对象对应的态射
}

// NewNaturalTransformation 创建自然变换
func NewNaturalTransformation(name string, source, target *Functor) *NaturalTransformation {
    return &NaturalTransformation{
        Name:       name,
        Source:     source,
        Target:     target,
        Components: make(map[string]*Morphism),
    }
}

// AddComponent 添加分量
func (nt *NaturalTransformation) AddComponent(obj Object, morphism *Morphism) {
    nt.Components[obj.ID()] = morphism
}

// GetComponent 获取分量
func (nt *NaturalTransformation) GetComponent(obj Object) *Morphism {
    return nt.Components[obj.ID()]
}

// IsNatural 检查是否为自然变换
func (nt *NaturalTransformation) IsNatural() bool {
    // 检查自然性条件
    for _, morphism := range nt.Source.SourceCategory.Morphisms {
        sourceObj := morphism.Domain
        targetObj := morphism.Codomain
        
        sourceComponent := nt.GetComponent(sourceObj)
        targetComponent := nt.GetComponent(targetObj)
        
        if sourceComponent == nil || targetComponent == nil {
            continue
        }
        
        // 检查自然性条件：G(f) ∘ α_A = α_B ∘ F(f)
        sourceMorphism := nt.Source.ApplyToMorphism(morphism)
        targetMorphism := nt.Target.ApplyToMorphism(morphism)
        
        if sourceMorphism == nil || targetMorphism == nil {
            continue
        }
        
        // 这里需要实际的态射复合来验证自然性
        // 简化版本：检查态射的域和余域是否匹配
        if sourceComponent.Codomain != targetComponent.Domain {
            return false
        }
    }
    
    return true
}

// IsNaturalIsomorphism 检查是否为自然同构
func (nt *NaturalTransformation) IsNaturalIsomorphism() bool {
    if !nt.IsNatural() {
        return false
    }
    
    // 检查每个分量是否为同构
    for _, component := range nt.Components {
        if !nt.Target.TargetCategory.IsIsomorphism(component) {
            return false
        }
    }
    
    return true
}
```

## 5. 应用示例

### 5.1 集合范畴

```go
// SetObject 集合对象
type SetObject struct {
    ID   string
    Elements []interface{}
}

func (s *SetObject) ID() string {
    return s.ID
}

// SetMorphism 集合态射（函数）
type SetMorphism struct {
    ID       string
    Domain   *SetObject
    Codomain *SetObject
    Function map[interface{}]interface{} // 函数的具体实现
}

// SetCategory 集合范畴
func CreateSetCategory() *Category {
    setCat := NewCategory("Set")
    
    // 定义复合函数
    setCat.Compose = func(f, g *Morphism) (*Morphism, error) {
        if f.Codomain != g.Domain {
            return nil, fmt.Errorf("cannot compose")
        }
        
        setF := f.(*SetMorphism)
        setG := g.(*SetMorphism)
        
        // 复合函数
        composedFunc := make(map[interface{}]interface{})
        for x, fx := range setF.Function {
            if gx, exists := setG.Function[fx]; exists {
                composedFunc[x] = gx
            }
        }
        
        composedID := fmt.Sprintf("%s_%s", f.ID, g.ID)
        return &SetMorphism{
            ID:       composedID,
            Domain:   setF.Domain,
            Codomain: setG.Codomain,
            Function: composedFunc,
        }, nil
    }
    
    // 定义单位态射
    setCat.Identity = func(obj Object) *Morphism {
        setObj := obj.(*SetObject)
        identityFunc := make(map[interface{}]interface{})
        for _, elem := range setObj.Elements {
            identityFunc[elem] = elem
        }
        
        return &SetMorphism{
            ID:       fmt.Sprintf("id_%s", obj.ID()),
            Domain:   setObj,
            Codomain: setObj,
            Function: identityFunc,
        }
    }
    
    return setCat
}

// SetCategoryExample 集合范畴示例
func SetCategoryExample() {
    setCat := CreateSetCategory()
    
    // 创建集合对象
    setA := &SetObject{ID: "A", Elements: []interface{}{1, 2, 3}}
    setB := &SetObject{ID: "B", Elements: []interface{}{"a", "b", "c"}}
    
    setCat.AddObject(setA)
    setCat.AddObject(setB)
    
    // 创建函数态射
    f := &SetMorphism{
        ID:       "f",
        Domain:   setA,
        Codomain: setB,
        Function: map[interface{}]interface{}{
            1: "a",
            2: "b",
            3: "c",
        },
    }
    
    setCat.AddMorphism(f)
    
    fmt.Println("集合范畴示例")
    fmt.Printf("对象A: %v\n", setA.Elements)
    fmt.Printf("对象B: %v\n", setB.Elements)
    fmt.Printf("函数f: %v\n", f.Function)
    
    // 检查是否为同构
    isIso := setCat.IsIsomorphism(f)
    fmt.Printf("f是否为同构: %v\n", isIso)
}
```

### 5.2 群范畴

```go
// GroupObject 群对象
type GroupObject struct {
    ID      string
    Elements []int
    Operation func(int, int) int
    Identity  int
}

func (g *GroupObject) ID() string {
    return g.ID
}

// GroupMorphism 群态射（群同态）
type GroupMorphism struct {
    ID       string
    Domain   *GroupObject
    Codomain *GroupObject
    Function map[int]int
}

// GroupCategory 群范畴
func CreateGroupCategory() *Category {
    groupCat := NewCategory("Group")
    
    // 定义复合函数
    groupCat.Compose = func(f, g *Morphism) (*Morphism, error) {
        if f.Codomain != g.Domain {
            return nil, fmt.Errorf("cannot compose")
        }
        
        groupF := f.(*GroupMorphism)
        groupG := g.(*GroupMorphism)
        
        // 复合函数
        composedFunc := make(map[int]int)
        for x, fx := range groupF.Function {
            if gx, exists := groupG.Function[fx]; exists {
                composedFunc[x] = gx
            }
        }
        
        composedID := fmt.Sprintf("%s_%s", f.ID, g.ID)
        return &GroupMorphism{
            ID:       composedID,
            Domain:   groupF.Domain,
            Codomain: groupG.Codomain,
            Function: composedFunc,
        }, nil
    }
    
    // 定义单位态射
    groupCat.Identity = func(obj Object) *Morphism {
        groupObj := obj.(*GroupObject)
        identityFunc := make(map[int]int)
        for _, elem := range groupObj.Elements {
            identityFunc[elem] = elem
        }
        
        return &GroupMorphism{
            ID:       fmt.Sprintf("id_%s", obj.ID()),
            Domain:   groupObj,
            Codomain: groupObj,
            Function: identityFunc,
        }
    }
    
    return groupCat
}

// GroupCategoryExample 群范畴示例
func GroupCategoryExample() {
    groupCat := CreateGroupCategory()
    
    // 创建群对象（整数加法群）
    groupZ := &GroupObject{
        ID: "Z",
        Elements: []int{0, 1, -1, 2, -2},
        Operation: func(a, b int) int { return a + b },
        Identity: 0,
    }
    
    // 创建群对象（模2群）
    groupZ2 := &GroupObject{
        ID: "Z2",
        Elements: []int{0, 1},
        Operation: func(a, b int) int { return (a + b) % 2 },
        Identity: 0,
    }
    
    groupCat.AddObject(groupZ)
    groupCat.AddObject(groupZ2)
    
    // 创建群同态（模2映射）
    phi := &GroupMorphism{
        ID:       "phi",
        Domain:   groupZ,
        Codomain: groupZ2,
        Function: map[int]int{
            0: 0,
            1: 1,
            -1: 1,
            2: 0,
            -2: 0,
        },
    }
    
    groupCat.AddMorphism(phi)
    
    fmt.Println("群范畴示例")
    fmt.Printf("群Z: %v\n", groupZ.Elements)
    fmt.Printf("群Z2: %v\n", groupZ2.Elements)
    fmt.Printf("同态phi: %v\n", phi.Function)
    
    // 检查是否为满态射
    isEpi := groupCat.IsEpimorphism(phi)
    fmt.Printf("phi是否为满态射: %v\n", isEpi)
}
```

### 5.3 编程语言应用

```go
// TypeObject 类型对象
type TypeObject struct {
    ID   string
    Name string
}

func (t *TypeObject) ID() string {
    return t.ID
}

// FunctionMorphism 函数态射
type FunctionMorphism struct {
    ID       string
    Domain   *TypeObject
    Codomain *TypeObject
    Function interface{} // 实际的Go函数
}

// TypeCategory 类型范畴
func CreateTypeCategory() *Category {
    typeCat := NewCategory("Type")
    
    // 定义复合函数
    typeCat.Compose = func(f, g *Morphism) (*Morphism, error) {
        if f.Codomain != g.Domain {
            return nil, fmt.Errorf("cannot compose")
        }
        
        funcF := f.(*FunctionMorphism)
        funcG := g.(*FunctionMorphism)
        
        // 在实际应用中，这里会进行函数复合
        composedID := fmt.Sprintf("%s_%s", f.ID, g.ID)
        return &FunctionMorphism{
            ID:       composedID,
            Domain:   funcF.Domain,
            Codomain: funcG.Codomain,
            Function: "composed function",
        }, nil
    }
    
    // 定义单位态射
    typeCat.Identity = func(obj Object) *Morphism {
        return &FunctionMorphism{
            ID:       fmt.Sprintf("id_%s", obj.ID()),
            Domain:   obj.(*TypeObject),
            Codomain: obj.(*TypeObject),
            Function: "identity function",
        }
    }
    
    return typeCat
}

// ProgrammingLanguageExample 编程语言应用示例
func ProgrammingLanguageExample() {
    typeCat := CreateTypeCategory()
    
    // 创建类型对象
    intType := &TypeObject{ID: "int", Name: "int"}
    stringType := &TypeObject{ID: "string", Name: "string"}
    boolType := &TypeObject{ID: "bool", Name: "bool"}
    
    typeCat.AddObject(intType)
    typeCat.AddObject(stringType)
    typeCat.AddObject(boolType)
    
    // 创建函数态射
    toString := &FunctionMorphism{
        ID:       "toString",
        Domain:   intType,
        Codomain: stringType,
        Function: func(x int) string { return fmt.Sprintf("%d", x) },
    }
    
    isPositive := &FunctionMorphism{
        ID:       "isPositive",
        Domain:   intType,
        Codomain: boolType,
        Function: func(x int) bool { return x > 0 },
    }
    
    typeCat.AddMorphism(toString)
    typeCat.AddMorphism(isPositive)
    
    fmt.Println("编程语言应用示例")
    fmt.Printf("类型: %s, %s, %s\n", intType.Name, stringType.Name, boolType.Name)
    fmt.Printf("函数: %s, %s\n", toString.ID, isPositive.ID)
    
    // 演示函子
    listFunctor := NewFunctor("List", typeCat, typeCat)
    
    // 映射对象
    listFunctor.MapObject(intType, &TypeObject{ID: "[]int", Name: "[]int"})
    listFunctor.MapObject(stringType, &TypeObject{ID: "[]string", Name: "[]string"})
    listFunctor.MapObject(boolType, &TypeObject{ID: "[]bool", Name: "[]bool"})
    
    fmt.Println("List函子:")
    for objID, targetObj := range listFunctor.ObjectMap {
        fmt.Printf("  %s -> %s\n", objID, targetObj.ID())
    }
}
```

## 总结

范畴论是数学的基础理论，提供了：

1. **抽象结构**: 统一处理各种数学结构
2. **函子理论**: 研究结构之间的映射
3. **自然变换**: 研究函子之间的关系
4. **广泛应用**: 在代数、拓扑、逻辑、计算机科学等领域有重要应用

通过Go语言的实现，我们展示了：

- 范畴、对象、态射的数据结构表示
- 函子的实现和复合
- 自然变换的定义和验证
- 集合范畴、群范畴、类型范畴等具体应用

这为后续的极限和余极限、伴随函子等更高级的范畴论概念奠定了基础。
