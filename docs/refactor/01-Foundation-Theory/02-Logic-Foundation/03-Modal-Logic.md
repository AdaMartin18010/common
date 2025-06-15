# 03-模态逻辑 (Modal Logic)

## 目录

- [03-模态逻辑 (Modal Logic)](#03-模态逻辑-modal-logic)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 模态算子](#11-模态算子)
    - [1.2 可能世界语义](#12-可能世界语义)
    - [1.3 模态系统](#13-模态系统)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 模态语言](#21-模态语言)
    - [2.2 Kripke语义](#22-kripke语义)
    - [2.3 公理系统](#23-公理系统)
  - [3. 推理系统](#3-推理系统)
    - [3.1 自然演绎](#31-自然演绎)
    - [3.2 表推演](#32-表推演)
    - [3.3 模型检查](#33-模型检查)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 模态逻辑数据结构](#41-模态逻辑数据结构)
    - [4.2 Kripke模型实现](#42-kripke模型实现)
    - [4.3 推理引擎](#43-推理引擎)
  - [5. 应用示例](#5-应用示例)
    - [5.1 知识逻辑](#51-知识逻辑)
    - [5.2 信念逻辑](#52-信念逻辑)
    - [5.3 时态逻辑](#53-时态逻辑)
  - [总结](#总结)

## 1. 基本概念

### 1.1 模态算子

**定义 1.1**: 模态算子用于表达"必然"和"可能"的概念。

1. **必然算子** ($\Box$): "必然地"
2. **可能算子** ($\Diamond$): "可能地"

**关系**: $\Diamond \phi \equiv \neg \Box \neg \phi$

**示例**:

- $\Box p$: "必然p"
- $\Diamond p$: "可能p"
- $\Box (p \rightarrow q)$: "必然地，如果p则q"
- $\Diamond p \land \Diamond q$: "可能p且可能q"

### 1.2 可能世界语义

**定义 1.2**: 可能世界语义是模态逻辑的标准语义。

- **可能世界**: 描述一种可能的状态或情况
- **可达关系**: 世界之间的可及性关系
- **真值**: 公式在特定世界中的真值

**直观理解**:

- $\Box \phi$ 在世界 $w$ 中为真，当且仅当 $\phi$ 在所有从 $w$ 可达的世界中为真
- $\Diamond \phi$ 在世界 $w$ 中为真，当且仅当 $\phi$ 在某个从 $w$ 可达的世界中为真

### 1.3 模态系统

**定义 1.3**: 常见的模态系统

1. **系统K**: 最基本的模态系统
2. **系统T**: K + $\Box \phi \rightarrow \phi$ (自反性)
3. **系统S4**: T + $\Box \phi \rightarrow \Box \Box \phi$ (传递性)
4. **系统S5**: S4 + $\Diamond \phi \rightarrow \Box \Diamond \phi$ (欧几里得性)

## 2. 形式化定义

### 2.1 模态语言

**定义 2.1**: 模态语言 $\mathcal{L}_\Box$ 的递归定义

1. **基础**: 每个原子命题 $p \in \mathcal{P}$ 是公式
2. **归纳**: 如果 $\phi$ 和 $\psi$ 是公式，则：
   - $\neg \phi$ 是公式
   - $(\phi \land \psi)$ 是公式
   - $(\phi \lor \psi)$ 是公式
   - $(\phi \rightarrow \psi)$ 是公式
   - $(\phi \leftrightarrow \psi)$ 是公式
   - $\Box \phi$ 是公式
   - $\Diamond \phi$ 是公式
3. **闭包**: 只有通过有限次应用上述规则得到的才是公式

**BNF语法**:

```
φ ::= p | ¬φ | (φ ∧ φ) | (φ ∨ φ) | (φ → φ) | (φ ↔ φ) | □φ | ◇φ
```

### 2.2 Kripke语义

**定义 2.2**: Kripke框架

**Kripke框架** $\mathcal{F} = (W, R)$ 由以下部分组成：

1. **世界集** $W$: 非空集合
2. **可达关系** $R \subseteq W \times W$: 二元关系

**定义 2.3**: Kripke模型

**Kripke模型** $\mathcal{M} = (W, R, V)$ 由以下部分组成：

1. **框架** $(W, R)$
2. **赋值函数** $V: \mathcal{P} \rightarrow 2^W$: 为每个原子命题指定在其中为真的世界集

**定义 2.4**: 语义函数

语义函数 $\llbracket \cdot \rrbracket_{\mathcal{M}, w}$ 递归定义：

1. **原子公式**: $\llbracket p \rrbracket_{\mathcal{M}, w} = \text{true}$ 当且仅当 $w \in V(p)$

2. **连接词**: 与命题逻辑相同

3. **模态算子**:
   - $\llbracket \Box \phi \rrbracket_{\mathcal{M}, w} = \text{true}$ 当且仅当对所有 $v \in W$，如果 $wRv$ 则 $\llbracket \phi \rrbracket_{\mathcal{M}, v} = \text{true}$
   - $\llbracket \Diamond \phi \rrbracket_{\mathcal{M}, w} = \text{true}$ 当且仅当存在 $v \in W$，$wRv$ 且 $\llbracket \phi \rrbracket_{\mathcal{M}, v} = \text{true}$

### 2.3 公理系统

**定义 2.5**: 系统K的公理

1. **命题公理**: 所有命题逻辑重言式
2. **分配公理**: $\Box(\phi \rightarrow \psi) \rightarrow (\Box \phi \rightarrow \Box \psi)$
3. **推理规则**:
   - **假言推理**: 从 $\phi$ 和 $\phi \rightarrow \psi$ 推出 $\psi$
   - **必然化**: 从 $\phi$ 推出 $\Box \phi$

**定义 2.6**: 其他系统的公理

- **系统T**: K + $\Box \phi \rightarrow \phi$
- **系统S4**: T + $\Box \phi \rightarrow \Box \Box \phi$
- **系统S5**: S4 + $\Diamond \phi \rightarrow \Box \Diamond \phi$

## 3. 推理系统

### 3.1 自然演绎

**定义 3.1**: 模态逻辑自然演绎规则

**模态规则**:

- **$\Box$-引入**: 如果从假设 $\phi$ 可以推出 $\psi$，且 $\phi$ 是必然的，则可以从 $\phi$ 推出 $\Box \psi$
- **$\Box$-消除**: 从 $\Box \phi$ 可以推出 $\phi$
- **$\Diamond$-引入**: 从 $\phi$ 可以推出 $\Diamond \phi$
- **$\Diamond$-消除**: 从 $\Diamond \phi$ 和 $\phi \rightarrow \psi$ 可以推出 $\Diamond \psi$

### 3.2 表推演

**定义 3.2**: 模态表推演规则

对于模态公式：

- **$\Box$公式**: 将 $\Box \phi$ 分解为 $\phi$，但只在可达的世界中
- **$\Diamond$公式**: 将 $\Diamond \phi$ 分解为 $\phi$，并创建新的可达世界

### 3.3 模型检查

**定义 3.3**: 模型检查算法

模型检查是验证公式在给定模型中是否为真的算法：

```pseudocode
ModelCheck(φ, M, w):
    case φ of
        p: return w ∈ V(p)
        ¬ψ: return not ModelCheck(ψ, M, w)
        ψ ∧ χ: return ModelCheck(ψ, M, w) and ModelCheck(χ, M, w)
        ψ ∨ χ: return ModelCheck(ψ, M, w) or ModelCheck(χ, M, w)
        ψ → χ: return not ModelCheck(ψ, M, w) or ModelCheck(χ, M, w)
        □ψ: return for all v such that wRv: ModelCheck(ψ, M, v)
        ◇ψ: return exists v such that wRv: ModelCheck(ψ, M, v)
```

## 4. Go语言实现

### 4.1 模态逻辑数据结构

```go
// World 可能世界
type World struct {
    ID   string
    Name string
}

// NewWorld 创建可能世界
func NewWorld(id, name string) *World {
    return &World{
        ID:   id,
        Name: name,
    }
}

// ModalFormula 模态逻辑公式
type ModalFormula struct {
    IsAtom      bool
    IsNegation  bool
    IsBinary    bool
    IsModal     bool
    
    // 原子公式
    Proposition string
    
    // 连接词
    Connective string
    Left       *ModalFormula
    Right      *ModalFormula
    
    // 模态算子
    ModalOperator string // "□" 或 "◇"
    Body          *ModalFormula
}

// NewAtomFormula 创建原子公式
func NewAtomFormula(proposition string) *ModalFormula {
    return &ModalFormula{
        IsAtom:      true,
        Proposition: proposition,
    }
}

// NewNegation 创建否定公式
func NewNegation(formula *ModalFormula) *ModalFormula {
    return &ModalFormula{
        IsNegation: true,
        Left:       formula,
    }
}

// NewBinaryFormula 创建二元连接词公式
func NewBinaryFormula(connective string, left, right *ModalFormula) *ModalFormula {
    return &ModalFormula{
        IsBinary:   true,
        Connective: connective,
        Left:       left,
        Right:      right,
    }
}

// NewNecessity 创建必然公式
func NewNecessity(body *ModalFormula) *ModalFormula {
    return &ModalFormula{
        IsModal:       true,
        ModalOperator: "□",
        Body:          body,
    }
}

// NewPossibility 创建可能公式
func NewPossibility(body *ModalFormula) *ModalFormula {
    return &ModalFormula{
        IsModal:       true,
        ModalOperator: "◇",
        Body:          body,
    }
}
```

### 4.2 Kripke模型实现

```go
// KripkeFrame Kripke框架
type KripkeFrame struct {
    Worlds map[string]*World
    Relation map[string]map[string]bool // R(w1, w2) = true 表示 w1 可达 w2
}

// NewKripkeFrame 创建Kripke框架
func NewKripkeFrame() *KripkeFrame {
    return &KripkeFrame{
        Worlds:   make(map[string]*World),
        Relation: make(map[string]map[string]bool),
    }
}

// AddWorld 添加世界
func (kf *KripkeFrame) AddWorld(world *World) {
    kf.Worlds[world.ID] = world
    kf.Relation[world.ID] = make(map[string]bool)
}

// AddRelation 添加可达关系
func (kf *KripkeFrame) AddRelation(from, to string) {
    if kf.Relation[from] == nil {
        kf.Relation[from] = make(map[string]bool)
    }
    kf.Relation[from][to] = true
}

// IsAccessible 检查可达性
func (kf *KripkeFrame) IsAccessible(from, to string) bool {
    if kf.Relation[from] == nil {
        return false
    }
    return kf.Relation[from][to]
}

// KripkeModel Kripke模型
type KripkeModel struct {
    Frame     *KripkeFrame
    Valuation map[string]map[string]bool // V(p, w) = true 表示命题p在世界w中为真
}

// NewKripkeModel 创建Kripke模型
func NewKripkeModel(frame *KripkeFrame) *KripkeModel {
    return &KripkeModel{
        Frame:     frame,
        Valuation: make(map[string]map[string]bool),
    }
}

// SetValuation 设置赋值
func (km *KripkeModel) SetValuation(proposition, world string, value bool) {
    if km.Valuation[proposition] == nil {
        km.Valuation[proposition] = make(map[string]bool)
    }
    km.Valuation[proposition][world] = value
}

// GetValuation 获取赋值
func (km *KripkeModel) GetValuation(proposition, world string) bool {
    if km.Valuation[proposition] == nil {
        return false
    }
    return km.Valuation[proposition][world]
}

// Evaluate 计算公式在给定世界中的真值
func (km *KripkeModel) Evaluate(formula *ModalFormula, world string) bool {
    if formula.IsAtom {
        return km.GetValuation(formula.Proposition, world)
    }
    
    if formula.IsNegation {
        return !km.Evaluate(formula.Left, world)
    }
    
    if formula.IsBinary {
        left := km.Evaluate(formula.Left, world)
        right := km.Evaluate(formula.Right, world)
        
        switch formula.Connective {
        case "∧":
            return left && right
        case "∨":
            return left || right
        case "→":
            return !left || right
        case "↔":
            return left == right
        }
    }
    
    if formula.IsModal {
        return km.evaluateModal(formula, world)
    }
    
    return false
}

// evaluateModal 计算模态公式的真值
func (km *KripkeModel) evaluateModal(formula *ModalFormula, world string) bool {
    if formula.ModalOperator == "□" {
        // 必然公式：在所有可达世界中为真
        for targetWorld := range km.Frame.Worlds {
            if km.Frame.IsAccessible(world, targetWorld) {
                if !km.Evaluate(formula.Body, targetWorld) {
                    return false
                }
            }
        }
        return true
    } else if formula.ModalOperator == "◇" {
        // 可能公式：在某个可达世界中为真
        for targetWorld := range km.Frame.Worlds {
            if km.Frame.IsAccessible(world, targetWorld) {
                if km.Evaluate(formula.Body, targetWorld) {
                    return true
                }
            }
        }
        return false
    }
    
    return false
}

// IsValid 检查公式在模型中是否有效
func (km *KripkeModel) IsValid(formula *ModalFormula) bool {
    for world := range km.Frame.Worlds {
        if !km.Evaluate(formula, world) {
            return false
        }
    }
    return true
}

// IsSatisfiable 检查公式在模型中是否可满足
func (km *KripkeModel) IsSatisfiable(formula *ModalFormula) bool {
    for world := range km.Frame.Worlds {
        if km.Evaluate(formula, world) {
            return true
        }
    }
    return false
}
```

### 4.3 推理引擎

```go
// ModalLogicEngine 模态逻辑推理引擎
type ModalLogicEngine struct {
    model *KripkeModel
}

// NewModalLogicEngine 创建模态逻辑推理引擎
func NewModalLogicEngine() *ModalLogicEngine {
    return &ModalLogicEngine{}
}

// SetupExampleModel 设置示例模型
func (e *ModalLogicEngine) SetupExampleModel() {
    frame := NewKripkeFrame()
    
    // 添加世界
    w1 := NewWorld("w1", "世界1")
    w2 := NewWorld("w2", "世界2")
    w3 := NewWorld("w3", "世界3")
    
    frame.AddWorld(w1)
    frame.AddWorld(w2)
    frame.AddWorld(w3)
    
    // 添加可达关系
    frame.AddRelation("w1", "w1")
    frame.AddRelation("w1", "w2")
    frame.AddRelation("w2", "w2")
    frame.AddRelation("w2", "w3")
    frame.AddRelation("w3", "w3")
    
    e.model = NewKripkeModel(frame)
    
    // 设置赋值
    e.model.SetValuation("p", "w1", true)
    e.model.SetValuation("p", "w2", false)
    e.model.SetValuation("p", "w3", true)
    
    e.model.SetValuation("q", "w1", false)
    e.model.SetValuation("q", "w2", true)
    e.model.SetValuation("q", "w3", false)
}

// ProveModalEquivalence 证明模态等价
func (e *ModalLogicEngine) ProveModalEquivalence() {
    // 证明 ◇p ≡ ¬□¬p
    
    p := NewAtomFormula("p")
    notP := NewNegation(p)
    boxNotP := NewNecessity(notP)
    notBoxNotP := NewNegation(boxNotP)
    diamondP := NewPossibility(p)
    
    // 检查等价性
    fmt.Println("证明 ◇p ≡ ¬□¬p")
    
    for world := range e.model.Frame.Worlds {
        diamondValue := e.model.Evaluate(diamondP, world)
        notBoxValue := e.model.Evaluate(notBoxNotP, world)
        
        fmt.Printf("世界 %s: ◇p = %v, ¬□¬p = %v, 等价 = %v\n", 
            world, diamondValue, notBoxValue, diamondValue == notBoxValue)
    }
}

// ProveModalAxioms 证明模态公理
func (e *ModalLogicEngine) ProveModalAxioms() {
    // 证明系统T的公理：□p → p
    
    p := NewAtomFormula("p")
    boxP := NewNecessity(p)
    axiomT := NewBinaryFormula("→", boxP, p)
    
    fmt.Println("证明系统T公理：□p → p")
    
    for world := range e.model.Frame.Worlds {
        value := e.model.Evaluate(axiomT, world)
        fmt.Printf("世界 %s: □p → p = %v\n", world, value)
    }
    
    // 检查是否在所有世界中为真
    if e.model.IsValid(axiomT) {
        fmt.Println("公理T在所有世界中都成立")
    } else {
        fmt.Println("公理T在某些世界中不成立")
    }
}

// ModelChecking 模型检查
func (e *ModalLogicEngine) ModelChecking(formula *ModalFormula) {
    fmt.Printf("模型检查公式: %s\n", formula.String())
    
    for world := range e.model.Frame.Worlds {
        value := e.model.Evaluate(formula, world)
        fmt.Printf("世界 %s: %v\n", world, value)
    }
    
    if e.model.IsValid(formula) {
        fmt.Println("公式在所有世界中都成立")
    } else if e.model.IsSatisfiable(formula) {
        fmt.Println("公式在某些世界中成立")
    } else {
        fmt.Println("公式在所有世界中都不成立")
    }
}
```

## 5. 应用示例

### 5.1 知识逻辑

```go
// EpistemicLogic 知识逻辑示例
func EpistemicLogic() {
    // 知识逻辑：□_i φ 表示"智能体i知道φ"
    
    frame := NewKripkeFrame()
    
    // 添加世界（可能的状态）
    w1 := NewWorld("w1", "状态1")
    w2 := NewWorld("w2", "状态2")
    
    frame.AddWorld(w1)
    frame.AddWorld(w2)
    
    // 智能体1的知识关系（等价关系）
    frame.AddRelation("w1", "w1")
    frame.AddRelation("w1", "w2")
    frame.AddRelation("w2", "w1")
    frame.AddRelation("w2", "w2")
    
    model := NewKripkeModel(frame)
    
    // 设置赋值
    model.SetValuation("p", "w1", true)  // 在状态1中p为真
    model.SetValuation("p", "w2", false) // 在状态2中p为假
    
    // 构建知识公式：智能体1知道p
    p := NewAtomFormula("p")
    knowsP := NewNecessity(p) // □_1 p
    
    fmt.Println("知识逻辑示例")
    fmt.Println("智能体1知道p吗？")
    
    for world := range model.Frame.Worlds {
        value := model.Evaluate(knowsP, world)
        fmt.Printf("在状态 %s 中: %v\n", world, value)
    }
    
    // 解释：智能体1不知道p，因为p在状态2中为假
    // 而智能体1认为状态2是可能的
}
```

### 5.2 信念逻辑

```go
// DoxasticLogic 信念逻辑示例
func DoxasticLogic() {
    // 信念逻辑：□_i φ 表示"智能体i相信φ"
    // 信念关系不需要是自反的（可能相信错误的事情）
    
    frame := NewKripkeFrame()
    
    // 添加世界
    w1 := NewWorld("w1", "真实世界")
    w2 := NewWorld("w2", "错误信念世界")
    
    frame.AddWorld(w1)
    frame.AddWorld(w2)
    
    // 智能体的信念关系（非自反）
    frame.AddRelation("w1", "w2") // 在真实世界中，智能体认为w2是可能的
    frame.AddRelation("w2", "w2") // 在错误信念中，智能体仍然认为w2是可能的
    
    model := NewKripkeModel(frame)
    
    // 设置赋值
    model.SetValuation("p", "w1", true)  // 在真实世界中p为真
    model.SetValuation("p", "w2", false) // 在错误信念中p为假
    
    // 构建信念公式：智能体相信p
    p := NewAtomFormula("p")
    believesP := NewNecessity(p) // □_1 p
    
    fmt.Println("信念逻辑示例")
    fmt.Println("智能体相信p吗？")
    
    for world := range model.Frame.Worlds {
        value := model.Evaluate(believesP, world)
        fmt.Printf("在状态 %s 中: %v\n", world, value)
    }
    
    // 解释：智能体不相信p，因为p在w2中为假
    // 而智能体认为w2是可能的
}
```

### 5.3 时态逻辑

```go
// TemporalLogic 时态逻辑示例
func TemporalLogic() {
    // 时态逻辑：□φ 表示"总是φ"，◇φ 表示"有时φ"
    
    frame := NewKripkeFrame()
    
    // 添加时间点
    t1 := NewWorld("t1", "时间1")
    t2 := NewWorld("t2", "时间2")
    t3 := NewWorld("t3", "时间3")
    
    frame.AddWorld(t1)
    frame.AddWorld(t2)
    frame.AddWorld(t3)
    
    // 时间关系（线性时间）
    frame.AddRelation("t1", "t1")
    frame.AddRelation("t1", "t2")
    frame.AddRelation("t1", "t3")
    frame.AddRelation("t2", "t2")
    frame.AddRelation("t2", "t3")
    frame.AddRelation("t3", "t3")
    
    model := NewKripkeModel(frame)
    
    // 设置赋值
    model.SetValuation("p", "t1", true)  // 在时间1中p为真
    model.SetValuation("p", "t2", false) // 在时间2中p为假
    model.SetValuation("p", "t3", true)  // 在时间3中p为真
    
    // 构建时态公式
    p := NewAtomFormula("p")
    alwaysP := NewNecessity(p)    // □p (总是p)
    sometimesP := NewPossibility(p) // ◇p (有时p)
    
    fmt.Println("时态逻辑示例")
    
    for world := range model.Frame.Worlds {
        alwaysValue := model.Evaluate(alwaysP, world)
        sometimesValue := model.Evaluate(sometimesP, world)
        
        fmt.Printf("在时间 %s 中: 总是p = %v, 有时p = %v\n", 
            world, alwaysValue, sometimesValue)
    }
    
    // 解释：
    // - 在t1中：总是p为假（因为p在t2中为假），有时p为真
    // - 在t2中：总是p为假（因为p在t2中为假），有时p为真
    // - 在t3中：总是p为真（因为p在t3及以后都为真），有时p为真
}
```

## 总结

模态逻辑是经典逻辑的重要扩展，提供了：

1. **模态表达能力**: 可以表达必然性、可能性、知识、信念等概念
2. **可能世界语义**: 基于Kripke模型的直观语义
3. **多种模态系统**: K、T、S4、S5等不同强度的系统
4. **广泛应用**: 在哲学、计算机科学、人工智能等领域有重要应用

通过Go语言的实现，我们展示了：

- 模态逻辑公式的数据结构表示
- Kripke模型的实现
- 语义解释和模型检查
- 知识逻辑、信念逻辑、时态逻辑等应用

这为后续的时态逻辑、动态逻辑等更高级的模态系统奠定了基础。
