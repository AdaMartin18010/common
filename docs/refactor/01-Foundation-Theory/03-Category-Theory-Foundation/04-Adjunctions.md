# 04-伴随函子 (Adjunctions)

## 目录

- [04-伴随函子 (Adjunctions)](#04-伴随函子-adjunctions)
  - [目录](#目录)
  - [1. 概念定义](#1-概念定义)
    - [1.1 核心思想](#11-核心思想)
    - [1.2 非形式化例子：自由幺半群](#12-非形式化例子自由幺半群)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 Hom-Set 定义](#21-hom-set-定义)
    - [2.2 单位和余单位定义 (Unit and Counit)](#22-单位和余单位定义-unit-and-counit)
  - [3. 伴随的性质](#3-伴随的性质)
    - [3.1 伴随函子保持极限/余极限](#31-伴随函子保持极限余极限)
    - [3.2 伴随的唯一性](#32-伴随的唯一性)
  - [4. Go语言实现思路](#4-go语言实现思路)
    - [4.1 Currying 和 Uncurrying](#41-currying-和-uncurrying)
    - [4.2 示例：Reader Monad](#42-示例reader-monad)
  - [5. 参考文献](#5-参考文献)

## 1. 概念定义

### 1.1 核心思想

**伴随 (Adjunction)** 是范畴论中一个深刻且强大的概念，它描述了两个范畴之间的一对函子（$F$ 和 $G$）的一种非常紧密的"对偶"关系。这种关系比同构要弱，但比单纯的函子映射要强得多。

一个伴随关系，记作 $F \dashv G$，意味着函子 $F$ 是函子 $G$ 的**左伴随 (Left Adjoint)**，而 $G$ 是 $F$ 的**右伴随 (Right Adjoint)**。

这种关系的核心在于，它在两个范畴的态射集合（Hom-sets）之间建立了一个自然的双射（bijection）。具体来说，从 $F$ 作用后的对象 $F(C)$ 到某个对象 $D$ 的态射，与从原始对象 $C$ 到 $G$ 作用后的对象 $G(D)$ 的态射是一一对应的。

### 1.2 非形式化例子：自由幺半群

这是一个经典的伴随例子，可以帮助建立直觉。

- 范畴 $\mathcal{C} = \textbf{Set}$ (集合范畴)。
- 范畴 $\mathcal{D} = \textbf{Mon}$ (幺半群范畴)。

我们有两个函子：

1. **遗忘函子 (Forgetful Functor)** $G: \textbf{Mon} \to \textbf{Set}$。
    这个函子"忘记"了幺半群的结构（二元操作和单位元），只保留其底层的集合。例如，它将幺半群 $(\mathbb{N}, +, 0)$ 映射到集合 $\mathbb{N}$。
2. **自由函子 (Free Functor)** $F: \textbf{Set} \to \textbf{Mon}$。
    这个函子接收一个集合（比如 `{'a', 'b'}`），并从中构造出最"自由"、最没有多余约束的幺半群。这个幺半群就是由该集合元素组成的所有列表（字符串），其操作是列表连接，单位元是空列表。例如，它将集合 `{'a', 'b'}` 映射到幺半群 `(["", "a", "b", "aa", "ab", "ba", "bb", ...], ++, "")`。

伴随关系 $F \dashv G$ 在这里意味着：
> 从一个集合 `S` 生成的自由幺半群 `F(S)` 到任何其他幺半群 `M` 的同态（保持幺半群结构的函数），与从集合 `S` 到 `M` 的底层集合 `G(M)` 的普通函数，是完全一一对应的。

换句话说，一旦你定义了如何将生成元（集合 `S` 中的元素）映射到目标幺半群 `M` 中，这个映射就**唯一地**决定了一个完整的幺半群同态。

## 2. 形式化定义

伴随关系 $F: \mathcal{C} \to \mathcal{D}$ 和 $G: \mathcal{D} \to \mathcal{C}$ 之间 ($F \dashv G$) 有几种等价的定义方式。

### 2.1 Hom-Set 定义

存在一个对所有对象 $C \in \text{Ob}(\mathcal{C})$ 和 $D \in \text{Ob}(\mathcal{D})$ 都自然的双射：
$$
\text{Hom}_\mathcal{D}(F(C), D) \cong \text{Hom}_\mathcal{C}(C, G(D))
$$
"自然"意味着这个双射与态射的组合是兼容的。这个双射通常被称为 **伴随同构 (adjunction isomorphism)**。

### 2.2 单位和余单位定义 (Unit and Counit)

存在两个自然变换：

1. **单位 (Unit)**: $\eta: \text{id}_\mathcal{C} \Rightarrow G \circ F$
2. **余单位 (Counit)**: $\epsilon: F \circ G \Rightarrow \text{id}_\mathcal{D}$

它们必须满足以下两个**三角恒等式 (Triangle Identities)**：

1. $(\epsilon F) \circ (F \eta) = \text{id}_F$
    （这里 $\text{id}_F$ 是函子 $F$ 的恒等自然变换）
    可视化：对于任何 $C \in \text{Ob}(\mathcal{C})$，下面的图是交换的：
    $$
    \begin{CD}
    F(C) @>F(\eta_C)>> F(G(F(C))) \\
    @| @VV{\epsilon_{F(C)}}V \\
    F(C) @= F(C)
    \end{CD}
    $$

2. $(G \epsilon) \circ (\eta G) = \text{id}_G$
    （这里 $\text{id}_G$ 是函子 $G$ 的恒等自然变换）
    可视化：对于任何 $D \in \text{Ob}(\mathcal{D})$，下面的图是交换的：
    $$
    \begin{CD}
    G(D) @>\eta_{G(D)}>> G(F(G(D))) \\
    @| @VV{G(\epsilon_D)}V \\
    G(D) @= G(D)
    \end{CD}
    $$

单位 $\eta$ 告诉我们如何将一个对象 $C$ 嵌入到经过"先 $F$ 后 $G$"变换后的结果中。余单位 $\epsilon$ 则告诉我们如何从经过"先 $G$ 后 $F$"变换后的结果中提取出原始对象 $D$。三角恒等式确保这两个过程是一致的。

## 3. 伴随的性质

### 3.1 伴随函子保持极限/余极限

这是伴随函子最重要的性质之一：

- **右伴随函子 ($G$) 保持所有存在的极限。**
  例如，如果 $D_1, D_2$ 在 $\mathcal{D}$ 中有乘积 $D_1 \times D_2$，那么 $G(D_1 \times D_2)$ 就是 $G(D_1)$ 和 $G(D_2)$ 在 $\mathcal{C}$ 中的乘积 $G(D_1) \times G(D_2)$。
- **左伴随函子 ($F$) 保持所有存在的余极限。**
  例如，如果 $C_1, C_2$ 在 $\mathcal{C}$ 中有余积 $C_1 + C_2$，那么 $F(C_1 + C_2)$ 就是 $F(C_1)$ 和 $F(C_2)$ 在 $\mathcal{D}$ 中的余积 $F(C_1) + F(C_2)$。

### 3.2 伴随的唯一性

一个函子的左伴随（如果存在）在自然同构的意义下是唯一的。同样，右伴随也是唯一的。

## 4. Go语言实现思路

在像 Go 这样的编程语言中，直接实现范畴和函子的形式化结构是比较复杂的。但我们可以从函数式编程的角度来理解伴随关系。

### 4.1 Currying 和 Uncurrying

在编程中，柯里化 (Currying) 的过程与伴随关系非常相似。
考虑三个类型 `C`, `X`, `D`。
一个函数 `f: (C, X) -> D` 可以被转换为一个高阶函数 `g: C -> (X -> D)`。

这种转换是可逆的，并且是一一对应的。这正是在笛卡尔闭范畴（大多数函数式语言的类型系统模型）中 `(C, -)` 函子和 `(C -> -)` 函子之间的伴随关系。
$$
\text{Hom}(C \times X, D) \cong \text{Hom}(X, C \to D)
$$
这里，`(- \times C)` 是左伴随，而 `(C \to -)` 是右伴随。

### 4.2 示例：Reader Monad

`Reader` monad (或者说 `Environment` monad) 的核心是 `(E -> A)` 这样一个函数类型，其中 `E` 是环境类型。
这个 `Reader` 函子 $R_E(A) = (E \to A)$ 是函子 $P_E(A) = E \times A$ 的右伴随。
这种伴随关系解释了为什么我们可以将一个需要环境的计算 `(E -> A)` 和一个简单的值 `A` 关联起来。

```go
package main

import "fmt"

// F(C) = (C, E) (Product)
func product[C, E any](c C, e E) struct { C; E } {
 return struct { C; E }{c, e}
}

// G(D) = (E -> D) (Function type, Reader)
type Reader[E, D any] func(e E) D

// Hom(F(C), D) -> Hom((C, E), D)
type Hom_F_D[C, E, D any] func(p struct { C; E }) D

// Hom(C, G(D)) -> Hom(C, (E -> D))
type Hom_C_G[C, E, D any] func(c C) Reader[E, D]

// to: Hom((C, E), D) -> Hom(C, (E -> D))
func to[C, E, D any](f Hom_F_D[C, E, D]) Hom_C_G[C, E, D] {
 return func(c C) Reader[E, D] {
  return func(e E) D {
   return f(struct { C; E }{c, e})
  }
 }
}

// from: Hom(C, (E -> D)) -> Hom((C, E), D)
func from[C, E, D any](g Hom_C_G[C, E, D]) Hom_F_D[C, E, D] {
 return func(p struct { C; E }) D {
  return g(p.C)(p.E)
 }
}

func main() {
 // Example function in Hom((C, E), D)
 addWithEnv := func(p struct{ int; int }) int {
  return p.int + p.int // Simplified, assume first is C, second is E
 }

 // Transform it to Hom(C, (E -> D))
 curriedAdd := to(addWithEnv)

 // Use the curried version
 add5 := curriedAdd(5) // c = 5
 result := add5(10)   // e = 10
 fmt.Println(result)  // 15

 // Transform it back
 originalAdd := from(curriedAdd)
 fmt.Println(originalAdd(struct{ int; int }{5, 10})) // 15
}
```

这个例子展示了 Hom-set 同构的核心思想，即两种不同形式的函数之间可以进行无损的双向转换。

## 5. 参考文献

1. Mac Lane, Saunders. *Categories for the Working Mathematician*. 2nd ed., Springer, 1998.
2. Awodey, Steve. *Category Theory*. 2nd ed., Oxford University Press, 2010.
3. Milewski, Bartosz. *Category Theory for Programmers*. 2019.
