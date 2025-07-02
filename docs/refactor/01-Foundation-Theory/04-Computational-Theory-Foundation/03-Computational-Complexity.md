# 03-计算复杂性 (Computational Complexity)

## 目录

- [03-计算复杂性 (Computational Complexity)](#03-计算复杂性-computational-complexity)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 计算复杂性概述](#11-计算复杂性概述)
    - [1.2 基本定义](#12-基本定义)
  - [2. 时间复杂度](#2-时间复杂度)
    - [2.1 渐近记号](#21-渐近记号)
    - [2.2 常见时间复杂度类](#22-常见时间复杂度类)
  - [3. 空间复杂度](#3-空间复杂度)
    - [3.1 空间复杂性类](#31-空间复杂性类)
    - [3.2 空间层次定理](#32-空间层次定理)
  - [4. 复杂性类](#4-复杂性类)
    - [4.1 基本复杂性类](#41-基本复杂性类)
    - [4.2 NP完全性](#42-np完全性)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 复杂性分析框架](#51-复杂性分析框架)
    - [5.2 时间复杂度分析](#52-时间复杂度分析)
    - [5.3 空间复杂度分析](#53-空间复杂度分析)
  - [6. 定理与证明](#6-定理与证明)
    - [6.1 时间层次定理](#61-时间层次定理)
    - [6.2 Cook-Levin定理](#62-cook-levin定理)
  - [7. 应用示例](#7-应用示例)
    - [7.1 SAT问题](#71-sat问题)
    - [7.2 旅行商问题 (TSP)](#72-旅行商问题-tsp)
  - [8. 参考文献](#8-参考文献)

## 1. 基础概念

### 1.1 计算复杂性概述

计算复杂性理论研究算法和问题的计算资源需求：

- **时间复杂度**：算法执行所需的时间
- **空间复杂度**：算法执行所需的存储空间
- **复杂性类**：具有相似复杂性的问题集合
- **应用领域**：算法设计、密码学、人工智能

### 1.2 基本定义

**定义 1.1** (计算模型 - 图灵机)

$$
M = (Q, \Sigma, \Gamma, \delta, q_0, q_{accept}, q_{reject})
$$

其中：
- $Q$: 有限状态集
- $\Sigma$: 输入字母表
- $\Gamma$: 带字母表 ($\Sigma \subseteq \Gamma$)
- $\delta: Q \times \Gamma \rightarrow Q \times \Gamma \times \{L, R\}$: 转移函数
- $q_0$: 初始状态
- $q_{accept}$: 接受状态
- $q_{reject}$: 拒绝状态

**定义 1.2** (时间复杂度)

对于图灵机 $M$ 和输入 $w$， $M$ 在 $w$ 上的运行时间 $t_M(w)$ 是 $M$ 停机前执行的步数。
对于函数 $f: \mathbb{N} \to \mathbb{N}$，我们说 $M$ 的时间复杂度是 $O(f(n))$，如果存在常数 $c > 0$，使得对于所有长度为 $n$ 的输入 $w$，有：
$$
t_M(w) \leq c \cdot f(n)
$$

**定义 1.3** (空间复杂度)

对于图灵机 $M$ 和输入 $w$，$M$ 在 $w$ 上使用的空间 $s_M(w)$ 是 $M$ 停机前访问的带方格数。
对于函数 $f: \mathbb{N} \to \mathbb{N}$，我们说 $M$ 的空间复杂度是 $O(f(n))$，如果存在常数 $c > 0$，使得对于所有长度为 $n$ 的输入 $w$，有：
$$
s_M(w) \leq c \cdot f(n)
$$

## 2. 时间复杂度

### 2.1 渐近记号

**定义 2.1** (大O记号)
对于函数 $f, g: \mathbb{N} \to \mathbb{N}$，我们说 $f(n) = O(g(n))$，如果存在常数 $c > 0$ 和 $n_0 \in \mathbb{N}$，使得对于所有 $n \geq n_0$，有：
$$
f(n) \leq c \cdot g(n)
$$

**定义 2.2** (大Ω记号)
对于函数 $f, g: \mathbb{N} \to \mathbb{N}$，我们说 $f(n) = \Omega(g(n))$，如果存在常数 $c > 0$ 和 $n_0 \in \mathbb{N}$，使得对于所有 $n \geq n_0$，有：
$$
f(n) \geq c \cdot g(n)
$$

**定义 2.3** (大Θ记号)
对于函数 $f, g: \mathbb{N} \to \mathbb{N}$，我们说 $f(n) = \Theta(g(n))$，如果 $f(n) = O(g(n))$ 且 $f(n) = \Omega(g(n))$。

### 2.2 常见时间复杂度类

**定义 2.4** (多项式时间)
$$
P = \bigcup_{k \in \mathbb{N}} \text{TIME}(n^k)
$$

**定义 2.5** (非确定性多项式时间)
$$
NP = \bigcup_{k \in \mathbb{N}} \text{NTIME}(n^k)
$$

**定义 2.6** (指数时间)
$$
EXP = \bigcup_{k \in \mathbb{N}} \text{TIME}(2^{n^k})
$$

## 3. 空间复杂度

### 3.1 空间复杂性类

**定义 3.1** (对数空间)
$$
L = \text{SPACE}(\log n)
$$

**定义 3.2** (多项式空间)
$$
PSPACE = \bigcup_{k \in \mathbb{N}} \text{SPACE}(n^k)
$$

**定义 3.3** (非确定性对数空间)
$$
NL = \text{NSPACE}(\log n)
$$

### 3.2 空间层次定理

**定理 3.1** (空间层次定理)
对于空间可构造函数 $f, g: \mathbb{N} \to \mathbb{N}$，如果 $f(n) = o(g(n))$，则：
$$
\text{SPACE}(f(n)) \subsetneq \text{SPACE}(g(n))
$$

## 4. 复杂性类

### 4.1 基本复杂性类

**定义 4.1** (复杂性类层次)
$$
L \subseteq NL \subseteq P \subseteq NP \subseteq PSPACE \subseteq EXP
$$
已知 $L \neq PSPACE$ 且 $NL \neq PSPACE$ 且 $P \neq EXP$。

**定义 4.2** (完全性问题)
对于复杂性类 $C$ 和语言 $L$，我们说 $L$ 是 $C$-完全的 (C-complete)，如果：
1. $L \in C$
2. 对于所有 $L' \in C$，有 $L' \leq_p L$ (L' 可在多项式时间内归约到 L)

### 4.2 NP完全性

**定义 4.3** (NP完全性)
语言 $L$ 是NP完全的 (NP-complete)，如果：
1. $L \in NP$
2. 对于所有 $L' \in NP$，有 $L' \leq_p L$

**定理 4.1** (Cook-Levin定理)
布尔可满足性问题 (SAT) 是NP完全的。

## 5. Go语言实现

### 5.1 复杂性分析框架

```go
package complexity

import (
	"time"
)

// Algorithm 定义算法接口
type Algorithm interface {
	Execute(input Input) Result
	Name() string
}

// Input 定义输入接口
type Input interface {
	Size() int
	Data() interface{}
}

// Result 定义结果接口
type Result interface {
	Output() interface{}
	Time() time.Duration
	Space() int64 // 使用 int64 以防大空间占用
}

// TimeComplexityAnalysis 分析时间复杂度
func TimeComplexityAnalysis(algo Algorithm, inputs []Input) {
	// ... 实现分析逻辑 ...
}

// SpaceComplexityAnalysis 分析空间复杂度
func SpaceComplexityAnalysis(algo Algorithm, inputs []Input) {
	// ... 实现分析逻辑 ...
}
```

### 5.2 时间复杂度分析

```go
package complexity

import (
	"fmt"
	"time"
)

// BasicResult 基础结果实现
type BasicResult struct {
	output interface{}
	time   time.Duration
	space  int64
}

func (r *BasicResult) Output() interface{}   { return r.output }
func (r *BasicResult) Time() time.Duration   { return r.time }
func (r *BasicResult) Space() int64          { return r.space }

// MeasureTime 执行算法并测量时间
func MeasureTime(algo Algorithm, input Input) Result {
	startTime := time.Now()
	// 实际空间分析需要更复杂的工具，如 `pprof`
	output := algo.Execute(input)
	duration := time.Since(startTime)

	return &BasicResult{
		output: output,
		time:   duration,
		space:  -1, // 空间分析占位符
	}
}

func RunTimeAnalysis(algo Algorithm, inputs []Input) {
    fmt.Printf("--- Time Complexity Analysis for %s ---\n", algo.Name())
    for _, input := range inputs {
        result := MeasureTime(algo, input)
        fmt.Printf("Input Size: %d, Time Taken: %s\n", input.Size(), result.Time())
    }
}
```

### 5.3 空间复杂度分析

空间复杂度的精确测量在Go中较为困难，通常通过分析数据结构和算法逻辑来估算。
`pprof` 等工具可以用来分析内存分配，但精确的峰值空间使用量难以在运行时简单获取。

```go
// 估算递归斐波那契函数的空间复杂度
// Space: O(n) due to recursion depth
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}
```

## 6. 定理与证明

### 6.1 时间层次定理

**定理 6.1**
若 $f(n)$ 是时间可构造函数，则：
$$
\text{TIME}(f(n)) \subsetneq \text{TIME}(f(n) \log^2 f(n))
$$
**证明概要**:
通过对角线方法构造一个语言 $L$，它属于 $\text{TIME}(f(n) \log^2 f(n))$ 但不属于 $\text{TIME}(f(n))$。
构造一个图灵机 $D$，在输入 $M$ (一个图灵机的编码) 时，模拟 $M$ 在输入 $M$ 上的运行，但时间限制为 $f(|M|) \log |M|$。$D$ 的输出与被模拟的 $M$ 的输出相反。

### 6.2 Cook-Levin定理

**定理 6.2**
SAT 是 NP-完全的。
**证明概要**:
1. **SAT $\in$ NP**: 给定一个布尔表达式的赋值，可以在多项式时间内验证该赋值是否使其为真。
2. **所有 NP 问题 $\leq_p$ SAT**: 对于任何一个 NP 问题，其判定过程可以由一个非确定性图灵机在多项式时间内完成。该图灵机的计算历史 (computation tableau) 可以被编码成一个巨大的布尔表达式。该表达式为真的充要条件是，存在一个接受该输入的计算路径。这个转换过程可以在多项式时间内完成。

## 7. 应用示例

### 7.1 SAT问题

**问题**: 给定一个布尔合取范式 (CNF) 公式，是否存在一个变量赋值使其为真？

```go
// 这是一个简化的 SAT 求解器示例，仅用于说明
// 它采用暴力破解，时间复杂度为 O(2^n * m)，n是变量数，m是子句数
func solveSAT(clauses [][]int, numVars int) bool {
    // 尝试所有 2^n 种可能的赋值
    for i := 0; i < (1 << numVars); i++ {
        assignment := make([]bool, numVars+1)
        // 生成当前赋值
        for j := 0; j < numVars; j++ {
            if (i>>j)&1 == 1 {
                assignment[j+1] = true
            }
        }

        // 检查是否满足所有子句
        satisfied := true
        for _, clause := range clauses {
            clauseSatisfied := false
            for _, literal := range clause {
                variable := literal
                negated := false
                if variable < 0 {
                    variable = -variable
                    negated = true
                }
                if assignment[variable] != negated {
                    clauseSatisfied = true
                    break
                }
            }
            if !clauseSatisfied {
                satisfied = false
                break
            }
        }
        if satisfied {
            return true
        }
    }
    return false
}
```

### 7.2 旅行商问题 (TSP)

**问题**: 给定一个城市列表和每对城市之间的距离，找出访问每个城市一次并返回起点的最短可能路线。TSP 是 NP-hard 问题。

## 8. 参考文献

1. Sipser, Michael. *Introduction to the Theory of Computation*. 3rd ed., Cengage Learning, 2012.
2. Arora, Sanjeev, and Boaz Barak. *Computational Complexity: A Modern Approach*. Cambridge University Press, 2009.
3. Papadimitriou, Christos H. *Computational Complexity*. Addison-Wesley, 1994. 