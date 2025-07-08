# 04-算法分析 (Algorithm Analysis)

## 目录

- [04-算法分析 (Algorithm Analysis)](#04-算法分析-algorithm-analysis)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 算法分析概述](#11-算法分析概述)
    - [1.2 基本定义](#12-基本定义)
  - [2. 算法正确性](#2-算法正确性)
    - [2.1 循环不变量](#21-循环不变量)
    - [2.2 递归正确性](#22-递归正确性)
  - [3. 算法效率](#3-算法效率)
    - [3.1 渐近分析](#31-渐近分析)
    - [3.2 主定理 (Master Theorem)](#32-主定理-master-theorem)
  - [4. 算法设计与优化](#4-算法设计与优化)
    - [4.1 分治策略](#41-分治策略)
    - [4.2 动态规划](#42-动态规划)
    - [4.3 贪心算法](#43-贪心算法)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 算法分析框架](#51-算法分析框架)
    - [5.2 排序算法分析](#52-排序算法分析)
    - [5.3 搜索算法分析](#53-搜索算法分析)
  - [6. 高级主题](#6-高级主题)
    - [6.1 摊还分析 (Amortized Analysis)](#61-摊还分析-amortized-analysis)
    - [6.2 概率算法分析](#62-概率算法分析)
  - [7. 参考文献](#7-参考文献)

## 1. 基础概念

### 1.1 算法分析概述

算法分析是研究算法性能的科学，包括：

- **正确性分析**：算法是否为所有合法输入产生正确输出。
- **效率分析**：算法的时间和空间资源消耗。
- **最优性分析**：算法是否达到了特定问题复杂度的下界。
- **稳定性分析**：针对排序等算法，相等的元素是否保持原有顺序。

### 1.2 基本定义

**定义 1.1** (算法)
算法是一个有限的、确定性的、有效的计算过程，它接受输入并产生输出。
- **有限性**：算法必须在有限步后终止。
- **确定性**：每个步骤都有明确无歧义的定义。
- **有效性**：每个操作都必须是可执行的。

**定义 1.2** (算法正确性)
算法 $A$ 对于问题 $P$ 是正确的，如果对于问题 $P$ 的每一个输入实例 $x$，算法 $A(x)$ 都能终止并给出 $P(x)$ 作为输出。

**定义 1.3** (算法复杂度)
设 $A$ 是一个算法， $n$ 是输入的大小：
- **时间复杂度 $T_A(n)$**: 在最坏情况下，算法 $A$ 在大小为 $n$ 的任何输入上所需的基本操作次数。
- **空间复杂度 $S_A(n)$**: 在最坏情况下，算法 $A$ 在大小为 $n$ 的任何输入上所需的最大存储空间。

## 2. 算法正确性

### 2.1 循环不变量

**定义 2.1** (循环不变量)
循环不变量是一个在循环的每次迭代之前和之后都保持为真的断言。证明循环不变量的正确性通常使用数学归纳法，包含三个步骤：
1.  **初始化 (Initialization)**：在循环第一次迭代开始前，不变量为真。
2.  **保持 (Maintenance)**：如果在某次迭代开始前不变量为真，那么在下一次迭代开始前它仍然为真。
3.  **终止 (Termination)**：当循环终止时，不变量（通常与终止条件结合）能帮助我们证明算法的正确性。

**示例**: 插入排序的循环不变量
> 在外层 for 循环的每次迭代开始时，子数组 `A[0...j-1]` 包含了原来 `A[0...j-1]` 中的元素，但已排好序。

### 2.2 递归正确性

递归算法的正确性通常通过数学归纳法来证明。
1.  **基础情况 (Base Case)**：证明算法对于最简单的输入（递归的终点）是正确的。
2.  **归纳假设 (Inductive Hypothesis)**：假设算法对于所有小于当前规模的输入都是正确的。
3.  **归纳步骤 (Inductive Step)**：证明在归纳假设下，算法对于当前规模的输入也是正确的。

## 3. 算法效率

### 3.1 渐近分析

**定义 3.1** (渐近上界 $O$-notation)
$f(n) = O(g(n))$ 表示存在正常数 $c$ 和 $n_0$，使得对所有 $n \geq n_0$，有 $0 \leq f(n) \leq c g(n)$。

**定义 3.2** (渐近下界 $\Omega$-notation)
$f(n) = \Omega(g(n))$ 表示存在正常数 $c$ 和 $n_0$，使得对所有 $n \geq n_0$，有 $0 \leq c g(n) \leq f(n)$。

**定义 3.3** (紧确界 $\Theta$-notation)
$f(n) = \Theta(g(n))$ 表示存在正常数 $c_1, c_2$ 和 $n_0$，使得对所有 $n \geq n_0$，有 $0 \leq c_1 g(n) \leq f(n) \leq c_2 g(n)$。

### 3.2 主定理 (Master Theorem)

主定理为求解形如 $T(n) = aT(n/b) + f(n)$ 的递归式提供了一种"菜谱式"的解决方案，其中 $a \geq 1, b > 1$ 是常数，$f(n)$ 是渐近正函数。

1.  如果 $f(n) = O(n^{\log_b a - \epsilon})$ 对某个常数 $\epsilon > 0$ 成立，则 $T(n) = \Theta(n^{\log_b a})$。
2.  如果 $f(n) = \Theta(n^{\log_b a})$，则 $T(n) = \Theta(n^{\log_b a} \log n)$。
3.  如果 $f(n) = \Omega(n^{\log_b a + \epsilon})$ 对某个常数 $\epsilon > 0$ 成立，并且对某个常数 $c < 1$ 和所有足够大的 $n$ 有 $a f(n/b) \leq c f(n)$，则 $T(n) = \Theta(f(n))$。

## 4. 算法设计与优化

### 4.1 分治策略

分治算法将问题分解为若干个规模较小的相同问题，递归地解决这些子问题，然后合并其结果。
- **分解 (Divide)**: 将问题划分为子问题。
- **解决 (Conquer)**: 递归解决子问题。
- **合并 (Combine)**: 合并子问题的解。

**示例**: 归并排序, 快速排序。

### 4.2 动态规划

动态规划通常用于求解最优化问题，它通过将问题分解为重叠的子问题，并存储子问题的解来避免重复计算。
- **最优子结构**: 问题的最优解包含了其子问题的最优解。
- **重叠子问题**: 在求解过程中，许多子问题被反复计算。

**方法**: 通常采用自底向上的方法（表格法）或带备忘的自顶向下方法。

### 4.3 贪心算法

贪心算法在每一步选择中都采取在当前状态下最好或最优的选择，从而希望能导致全局最好或最优的解。
- **贪心选择性质**: 可以通过局部最优选择来达到全局最优。
- **最优子结构**: 一个问题的最优解包含其子问题的最优解。

**示例**: 活动选择问题，霍夫曼编码。

## 5. Go语言实现

### 5.1 算法分析框架

```go
package algorithmanalysis

import (
	"time"
)

// Algorithm 定义了可分析算法的接口
type Algorithm interface {
	Name() string
	Execute(data interface{}) (result interface{}, steps int)
}

// AnalysisResult 存储单次运行的分析结果
type AnalysisResult struct {
	InputSize    int
	Duration     time.Duration
	Steps        int
	MemoryUsage  uint64 // in bytes
}

// RunAnalysis 运行分析并返回结果
func RunAnalysis(algo Algorithm, data interface{}, inputSize int) AnalysisResult {
	// ... 实现内存使用快照 ...
	startTime := time.Now()
	_, steps := algo.Execute(data)
	duration := time.Since(startTime)
	// ... 实现内存使用快照比较 ...

	return AnalysisResult{
		InputSize: inputSize,
		Duration:  duration,
		Steps:     steps,
	}
}
```

### 5.2 排序算法分析

**插入排序**
- **时间复杂度**: $O(n^2)$
- **空间复杂度**: $O(1)$
- **实现**:
```go
package sorting

func InsertionSort(arr []int) (sortedArr []int, steps int) {
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	
	for i := 1; i < n; i++ {
		key := a[i]
		j := i - 1
		steps++ // for the comparison
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
			steps += 2 // for comparison and assignment
		}
		a[j+1] = key
		steps++ // for the assignment
	}
	return a, steps
}
```

### 5.3 搜索算法分析

**二分搜索**
- **时间复杂度**: $O(\log n)$
- **空间复杂度**: $O(1)$ (迭代实现)
- **实现**:
```go
package searching

func BinarySearch(arr []int, target int) (index int, steps int) {
	low, high := 0, len(arr)-1
	
	for low <= high {
		steps++ // for the loop condition check
		mid := low + (high-low)/2
		steps++ // for calculating mid

		if arr[mid] == target {
			return mid, steps
		}
		steps++ // for the comparison

		if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
		steps++ // for the assignment
	}
	
	return -1, steps // not found
}
```

## 6. 高级主题

### 6.1 摊还分析 (Amortized Analysis)

摊还分析用于评估一个操作序列的平均时间。它关心的是整个序列的总成本，而不是单个最坏操作的成本。
- **聚合分析 (Aggregate analysis)**: 确定 $n$ 个操作序列的总成本上界 $T(n)$，则每个操作的摊还成本为 $T(n)/n$。
- **记账方法 (Accounting method)**: 对不同操作收取不同的费用（摊还成本），多收的费用作为"信用"存起来，用于支付未来成本低于摊还成本的操作。
- **势能方法 (Potential method)**: 类似于记账方法，但将预付的代价存储为"势能"，势能可以释放以支付未来的操作。

### 6.2 概率算法分析

用于分析那些行为部分依赖于随机数的算法。我们通常分析其期望运行时间。
**示例**: 随机化快速排序的期望运行时间为 $O(n \log n)$。

## 7. 参考文献

1. Cormen, Thomas H., et al. *Introduction to Algorithms*. 3rd ed., MIT Press, 2009.
2. Sedgewick, Robert, and Kevin Wayne. *Algorithms*. 4th ed., Addison-Wesley, 2011.
3. Kleinberg, Jon, and Éva Tardos. *Algorithm Design*. Addison-Wesley, 2005. 