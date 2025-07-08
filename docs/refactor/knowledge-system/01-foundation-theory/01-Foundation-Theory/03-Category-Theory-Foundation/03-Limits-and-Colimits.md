# 03-极限和余极限 (Limits and Colimits)

## 目录

- [03-极限和余极限 (Limits and Colimits)](#03-极限和余极限-limits-and-colimits)
  - [目录](#目录)
  - [1. 概念定义](#1-概念定义)
    - [1.1 核心思想](#11-核心思想)
    - [1.2 图示 (Diagram) 和锥 (Cone)](#12-图示-diagram-和锥-cone)
  - [2. 极限 (Limits)](#2-极限-limits)
    - [2.1 形式化定义](#21-形式化定义)
    - [2.2 常见极限类型](#22-常见极限类型)
      - [2.2.1 终对象 (Terminal Object)](#221-终对象-terminal-object)
      - [2.2.2 乘积 (Product)](#222-乘积-product)
      - [2.2.3 拉回 (Pullback)](#223-拉回-pullback)
      - [2.2.4 等化子 (Equalizer)](#224-等化子-equalizer)
  - [3. 余极限 (Colimits)](#3-余极限-colimits)
    - [3.1 形式化定义](#31-形式化定义)
    - [3.2 常见余极限类型](#32-常见余极限类型)
      - [3.2.1 始对象 (Initial Object)](#321-始对象-initial-object)
      - [3.2.2 余积 (Coproduct)](#322-余积-coproduct)
      - [3.2.3 推出 (Pushout)](#323-推出-pushout)
      - [3.2.4 余等化子 (Coequalizer)](#324-余等化子-coequalizer)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 乘积的实现](#41-乘积的实现)
    - [4.2 余积的实现](#42-余积的实现)
  - [5. 参考文献](#5-参考文献)

## 1. 概念定义

### 1.1 核心思想

在范畴论中，**极限 (Limit)** 和 **余极限 (Colimit)** 是两个核心的对偶概念，它们以一种统一和抽象的方式概括了数学和计算机科学中许多常见的构造。

- **极限** 可以被非形式地理解为一个"最优化"的解决方案，它将一族对象和它们之间的关系（一个"图示"）汇集到一个单一的对象中，同时保留了所有原始结构。它代表了图示中所有对象共有的"共同部分"或"约束"。
- **余极限** 是极限的对偶概念，可以理解为将一族对象"粘合"或"合并"在一起的最有效方式。

### 1.2 图示 (Diagram) 和锥 (Cone)

要理解极限，首先需要了解**图示 (Diagram)**和**锥 (Cone)**。

- **图示**: 在范畴 $\mathcal{C}$ 中的一个图示，是一个从一个小的索引范畴 $\mathcal{J}$ 到 $\mathcal{C}$ 的函子 $F: \mathcal{J} \to \mathcal{C}$。你可以把它想象成在 $\mathcal{C}$ 中选出了一组对象和态射，它们的形状由 $\mathcal{J}$ 决定。

- **锥 (Cone)**: 对于一个图示 $F: \mathcal{J} \to \mathcal{C}$，一个以顶点 (vertex) $N \in \text{Ob}(\mathcal{C})$ 为顶的锥，是一族态射 $\{\psi_X: N \to F(X)\}_{X \in \text{Ob}(\mathcal{J})}$，使得对于 $\mathcal{J}$ 中的每个态射 $f: X \to Y$，都有 $F(f) \circ \psi_X = \psi_Y$。

简单来说，一个锥就是从一个新对象 $N$ 出发，到图示中每个对象都有一条路径，并且这些路径与图示内部的路径是兼容的。

## 2. 极限 (Limits)

### 2.1 形式化定义

一个图示 $F: \mathcal{J} \to \mathcal{C}$ 的**极限**是一个特殊的锥，我们记作 $\lim F$。它是一个对象 $L \in \text{Ob}(\mathcal{C})$ 连同一族态射 $\{\pi_X: L \to F(X)\}_{X \in \text{Ob}(\mathcal{J})}$，满足以下**泛性质 (universal property)**：

对于任何其他以 $N$ 为顶的锥 $\{\psi_X: N \to F(X)\}_{X \in \text{Ob}(\mathcal{J})}$，存在一个**唯一的**态射 $u: N \to L$，使得对于 $\mathcal{J}$ 中所有的对象 $X$，都有 $\pi_X \circ u = \psi_X$。

$$
\forall X \in \text{Ob}(\mathcal{J}), \quad \pi_X \circ u = \psi_X
$$

这个唯一的态射 $u$ 确保了极限 $L$ 是所有锥中"最通用"或"最紧凑"的那个。

### 2.2 常见极限类型

不同的索引范畴 $\mathcal{J}$ 会产生不同类型的极限。

#### 2.2.1 终对象 (Terminal Object)

- **图示**: 空图示 (索引范畴 $\mathcal{J}$ 是空范畴)。
- **极限**: **终对象** $1$。对于范畴中任何对象 $A$，都存在唯一的态射 $!: A \to 1$。

#### 2.2.2 乘积 (Product)

- **图示**: 由两个（或多个）没有非恒等态射的对象组成的离散范畴。
- **极限**: **乘积** $A \times B$。它带有一对投影态射 $\pi_A: A \times B \to A$ 和 $\pi_B: A \times B \to B$。对于任何对象 $C$ 和一对态射 $f: C \to A, g: C \to B$，存在唯一的态射 $u: C \to A \times B$ 使得 $\pi_A \circ u = f$ 和 $\pi_B \circ u = g$。

#### 2.2.3 拉回 (Pullback)

- **图示**: $A \xrightarrow{f} C \xleftarrow{g} B$ (cospan)。
- **极限**: **拉回** $A \times_C B$。它是一个对象连同两个态射 $p_A: A \times_C B \to A$ 和 $p_B: A \times_C B \to B$，使得 $f \circ p_A = g \circ p_B$，并且对于任何其他满足此条件的构造，都存在一个唯一的态射到该拉回。

#### 2.2.4 等化子 (Equalizer)

- **图示**: $A \rightrightarrows B$ (两个平行箭头)。
- **极限**: **等化子**。对于两个态射 $f,g: A \to B$，它们的等化子是一个对象 $E$ 和一个态射 $e: E \to A$，使得 $f \circ e = g \circ e$，并且对于任何其他的 $e': E' \to A$ 满足 $f \circ e' = g \circ e'$，都存在一个唯一的态射 $u: E' \to E$ 使得 $e \circ u = e'$。

## 3. 余极限 (Colimits)

余极限是极限的对偶概念。我们只需将极限定义中的所有箭头反向即可。

### 3.1 形式化定义

- **余锥 (Cocone)**: 对于一个图示 $F: \mathcal{J} \to \mathcal{C}$，一个以顶点 $N$ 为顶的余锥，是一族态射 $\{\iota_X: F(X) \to N\}_{X \in \text{Ob}(\mathcal{J})}$，使得对于 $\mathcal{J}$ 中的每个态射 $f: X \to Y$，都有 $\iota_Y \circ F(f) = \iota_X$。

- **余极限 (Colimit)**: 图示 $F$ 的余极限 $\text{colim} F$ 是一个"通用"的余锥。它是一个对象 $C$ 连同一族态射 $\{\sigma_X: F(X) \to C\}_{X \in \text{Ob}(\mathcal{J})}$，使得对于任何其他余锥 $\{\iota_X: F(X) \to N\}_{X \in \text{Ob}(\mathcal{J})}$，存在一个**唯一的**态射 $u: C \to N$，使得 $u \circ \sigma_X = \iota_X$。

### 3.2 常见余极限类型

#### 3.2.1 始对象 (Initial Object)

- **对偶于**: 终对象。
- **余极限**: **始对象** $0$。对于范畴中任何对象 $A$，都存在唯一的态射 $!: 0 \to A$。

#### 3.2.2 余积 (Coproduct)

- **对偶于**: 乘积。
- **余极限**: **余积** $A + B$ (或 $A \oplus B$)。它带有一对内含态射 $\iota_A: A \to A+B$ 和 $\iota_B: B \to A+B$。在集合范畴 **Set** 中，这是不交并。

#### 3.2.3 推出 (Pushout)

- **对偶于**: 拉回。
- **图示**: $A \xleftarrow{f} C \xrightarrow{g} B$ (span)。
- **余极限**: **推出**。它将两个对象沿着一个共同的子对象"粘合"起来。

#### 3.2.4 余等化子 (Coequalizer)

- **对偶于**: 等化子。
- **图示**: $A \rightrightarrows B$ (两个平行箭头)。
- **余极限**: **余等化子**。它通过"等同"两个平行态射来构造一个新的商对象。

## 4. Go语言实现

### 4.1 乘积的实现

在Go中，我们可以用 `struct` 来表示类型（对象）的乘积。

```go
package main

import "fmt"

// 对象 A
type A struct {
	Value int
}

// 对象 B
type B struct {
	Name string
}

// 乘积 A × B
type ProductAB struct {
	P1 A // 投影到 A
	P2 B // 投影到 B
}

func main() {
	// 创建乘积实例
	p := ProductAB{
		P1: A{Value: 42},
		P2: B{Name: "Go"},
	}

	// 投影
	a := p.P1
	b := p.P2

	fmt.Printf("Projection to A: %v\n", a)
	fmt.Printf("Projection to B: %v\n", b)
}
```
泛性质的实现需要更高阶的类型系统，在Go中通常通过接口和反射来模拟，但会比较繁琐。

### 4.2 余积的实现

在Go中，可以使用接口或tagged union（如 `interface{}` 或特定接口类型）来模拟余积（不交并）。

```go
package main

import "fmt"

// 定义一个接口作为余积类型
type Coproduct interface{ isCoproduct() }

// 实现接口的对象 A
type ValA struct{ Value int }
func (ValA) isCoproduct() {}

// 实现接口的对象 B
type ValB struct{ Name string }
func (ValB) isCoproduct() {}

// 一个处理余积的函数
func processCoproduct(c Coproduct) {
	switch v := c.(type) {
	case ValA:
		fmt.Printf("Got A with value: %d\n", v.Value)
	case ValB:
		fmt.Printf("Got B with name: %s\n", v.Name)
	default:
		fmt.Println("Unknown type in coproduct")
	}
}

func main() {
	// 内含 A -> A+B
	ia := ValA{Value: 100}
	// 内含 B -> A+B
	ib := ValB{Name: "Category"}

	processCoproduct(ia)
	processCoproduct(ib)
}
```

## 5. 参考文献

1. Mac Lane, Saunders. *Categories for the Working Mathematician*. 2nd ed., Springer, 1998.
2. Awodey, Steve. *Category Theory*. 2nd ed., Oxford University Press, 2010.
3. Riehl, Emily. *Category Theory in Context*. Dover Publications, 2016. 