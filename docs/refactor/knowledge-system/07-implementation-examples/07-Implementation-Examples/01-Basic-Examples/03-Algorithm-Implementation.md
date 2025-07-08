# 03-算法实现 (Algorithm Implementation)

## 1. 概述

### 1.1 算法基础

**算法** 是解决特定问题的有限步骤序列。

**形式化定义**：
设 ```latex
A
``` 为算法，```latex
I
``` 为输入集合，```latex
O
``` 为输出集合，则：
$```latex
A: I \rightarrow O
```$

### 1.2 复杂度分析

**时间复杂度**：```latex
T(n) = O(f(n))
```
**空间复杂度**：```latex
S(n) = O(g(n))
```

## 2. 排序算法

### 2.1 快速排序

#### 2.1.1 理论基础

**快速排序** 是一种分治算法，平均时间复杂度为 ```latex
O(n \log n)
```。

**分治策略**：

1. **分解**：选择基准元素，将数组分为两部分
2. **解决**：递归排序两个子数组
3. **合并**：无需合并，原地排序

**形式化定义**：
$\text{QuickSort}(A) = \begin{cases}
A & \text{if } |A| \leq 1 \\
\text{QuickSort}(L) \oplus [p] \oplus \text{QuickSort}(R) & \text{otherwise}
\end{cases}$

其中 ```latex
L = \{x \in A | x < p\}
```，```latex
R = \{x \in A | x > p\}
```，```latex
p
``` 为基准元素。

#### 2.1.2 Go语言实现

```go
package algorithms

import (
 "fmt"
 "math/rand"
 "time"
)

// QuickSort 快速排序
func QuickSort[T comparable](arr []T, less func(T, T) bool) []T {
 if len(arr) <= 1 {
  return arr
 }

 // 选择基准元素
 pivot := selectPivot(arr)

 // 分区
 left, right := partition(arr, pivot, less)

 // 递归排序
 left = QuickSort(left, less)
 right = QuickSort(right, less)

 // 合并结果
 return append(append(left, pivot), right...)
}

// selectPivot 选择基准元素
func selectPivot[T comparable](arr []T) T {
 // 使用三数取中法
 n := len(arr)
 if n >= 3 {
  mid := n / 2
  if less(arr[0], arr[mid]) {
   if less(arr[mid], arr[n-1]) {
    return arr[mid]
   } else if less(arr[0], arr[n-1]) {
    return arr[n-1]
   } else {
    return arr[0]
   }
  } else {
   if less(arr[0], arr[n-1]) {
    return arr[0]
   } else if less(arr[mid], arr[n-1]) {
    return arr[n-1]
   } else {
    return arr[mid]
   }
  }
 }
 return arr[0]
}

// partition 分区函数
func partition[T comparable](arr []T, pivot T, less func(T, T) bool) ([]T, []T) {
 var left, right []T

 for _, v := range arr {
  if less(v, pivot) {
   left = append(left, v)
  } else if v != pivot {
   right = append(right, v)
  }
 }

 return left, right
}

// less 比较函数
func less(a, b int) bool {
 return a < b
}

// 泛型实现
type Sortable[T any] interface {
 Less(other T) bool
}

// QuickSortGeneric 泛型快速排序
func QuickSortGeneric[T Sortable[T]](arr []T) []T {
 if len(arr) <= 1 {
  return arr
 }

 pivot := arr[0]
 var left, right []T

 for i := 1; i < len(arr); i++ {
  if arr[i].Less(pivot) {
   left = append(left, arr[i])
  } else {
   right = append(right, arr[i])
  }
 }

 left = QuickSortGeneric(left)
 right = QuickSortGeneric(right)

 return append(append(left, pivot), right...)
}

// 函数式实现
type SortFunc[T any] func([]T) []T

// ComposeSort 组合排序函数
func ComposeSort[T any](sorts ...SortFunc[T]) SortFunc[T] {
 return func(arr []T) []T {
  result := arr
  for _, sort := range sorts {
   result = sort(result)
  }
  return result
 }
}
```

### 2.2 归并排序

#### 2.2.1 理论基础

**归并排序** 是一种稳定的分治排序算法，时间复杂度为 ```latex
O(n \log n)
```。

**分治策略**：
1. **分解**：将数组分成两半
2. **解决**：递归排序两个子数组
3. **合并**：合并两个有序子数组

**形式化定义**：
$\text{MergeSort}(A) = \begin{cases}
A & \text{if } |A| \leq 1 \\
\text{Merge}(\text{MergeSort}(L), \text{MergeSort}(R)) & \text{otherwise}
\end{cases}$

#### 2.2.2 Go语言实现

```go
// MergeSort 归并排序
func MergeSort[T comparable](arr []T, less func(T, T) bool) []T {
 if len(arr) <= 1 {
  return arr
 }

 mid := len(arr) / 2
 left := MergeSort(arr[:mid], less)
 right := MergeSort(arr[mid:], less)

 return merge(left, right, less)
}

// merge 合并两个有序数组
func merge[T comparable](left, right []T, less func(T, T) bool) []T {
 result := make([]T, 0, len(left)+len(right))
 i, j := 0, 0

 for i < len(left) && j < len(right) {
  if less(left[i], right[j]) {
   result = append(result, left[i])
   i++
  } else {
   result = append(result, right[j])
   j++
  }
 }

 // 添加剩余元素
 result = append(result, left[i:]...)
 result = append(result, right[j:]...)

 return result
}

// 并发实现
func MergeSortConcurrent[T comparable](arr []T, less func(T, T) bool) []T {
 if len(arr) <= 1 {
  return arr
 }

 mid := len(arr) / 2

 // 使用channel进行并发
 leftChan := make(chan []T)
 rightChan := make(chan []T)

 go func() {
  leftChan <- MergeSortConcurrent(arr[:mid], less)
 }()

 go func() {
  rightChan <- MergeSortConcurrent(arr[mid:], less)
 }()

 left := <-leftChan
 right := <-rightChan

 return merge(left, right, less)
}
```

## 3. 搜索算法

### 3.1 二分搜索

#### 3.1.1 理论基础

**二分搜索** 在有序数组中查找目标值，时间复杂度为 ```latex
O(\log n)
```。

**形式化定义**：
$\text{BinarySearch}(A, target) = \begin{cases}
\text{mid} & \text{if } A[\text{mid}] = target \\
\text{BinarySearch}(A[\text{left}:\text{mid}], target) & \text{if } A[\text{mid}] > target \\
\text{BinarySearch}(A[\text{mid}+1:\text{right}], target) & \text{if } A[\text{mid}] < target \\
-1 & \text{if } \text{left} > \text{right}
\end{cases}$

#### 3.1.2 Go语言实现

```go
// BinarySearch 二分搜索
func BinarySearch[T comparable](arr []T, target T, less func(T, T) bool) int {
 left, right := 0, len(arr)-1

 for left <= right {
  mid := left + (right-left)/2
  
  if arr[mid] == target {
   return mid
  } else if less(arr[mid], target) {
   left = mid + 1
  } else {
   right = mid - 1
  }
 }

 return -1
}

// BinarySearchFirst 查找第一个等于目标值的元素
func BinarySearchFirst[T comparable](arr []T, target T, less func(T, T) bool) int {
 left, right := 0, len(arr)-1
 result := -1

 for left <= right {
  mid := left + (right-left)/2
  
  if arr[mid] == target {
   result = mid
   right = mid - 1 // 继续向左查找
  } else if less(arr[mid], target) {
   left = mid + 1
  } else {
   right = mid - 1
  }
 }

 return result
}

// BinarySearchLast 查找最后一个等于目标值的元素
func BinarySearchLast[T comparable](arr []T, target T, less func(T, T) bool) int {
 left, right := 0, len(arr)-1
 result := -1

 for left <= right {
  mid := left + (right-left)/2
  
  if arr[mid] == target {
   result = mid
   left = mid + 1 // 继续向右查找
  } else if less(arr[mid], target) {
   left = mid + 1
  } else {
   right = mid - 1
  }
 }

 return result
}
```

### 3.2 深度优先搜索

#### 3.2.1 理论基础

**深度优先搜索 (DFS)** 是一种图遍历算法，优先访问深层节点。

**形式化定义**：
$\text{DFS}(G, v) = \begin{cases}
\text{visit}(v) & \text{if } v \text{ is unvisited} \\
\text{DFS}(G, u) & \text{for each unvisited neighbor } u \text{ of } v
\end{cases}$

#### 3.2.2 Go语言实现

```go
// Graph 图结构
type Graph[T comparable] struct {
 adjacency map[T][]T
}

// NewGraph 创建新图
func NewGraph[T comparable]() *Graph[T] {
 return &Graph[T]{
  adjacency: make(map[T][]T),
 }
}

// AddEdge 添加边
func (g *Graph[T]) AddEdge(from, to T) {
 g.adjacency[from] = append(g.adjacency[from], to)
}

// DFS 深度优先搜索
func (g *Graph[T]) DFS(start T) []T {
 visited := make(map[T]bool)
 result := make([]T, 0)

 g.dfsHelper(start, visited, &result)
 return result
}

// dfsHelper DFS辅助函数
func (g *Graph[T]) dfsHelper(node T, visited map[T]bool, result *[]T) {
 if visited[node] {
  return
 }

 visited[node] = true
 *result = append(*result, node)

 for _, neighbor := range g.adjacency[node] {
  g.dfsHelper(neighbor, visited, result)
 }
}

// DFSIterative 迭代式DFS
func (g *Graph[T]) DFSIterative(start T) []T {
 visited := make(map[T]bool)
 result := make([]T, 0)
 stack := []T{start}

 for len(stack) > 0 {
  node := stack[len(stack)-1]
  stack = stack[:len(stack)-1]
  
  if visited[node] {
   continue
  }
  
  visited[node] = true
  result = append(result, node)
  
  // 将邻居压入栈中（逆序）
  for i := len(g.adjacency[node]) - 1; i >= 0; i-- {
   neighbor := g.adjacency[node][i]
   if !visited[neighbor] {
    stack = append(stack, neighbor)
   }
  }
 }

 return result
}
```

### 3.3 广度优先搜索

#### 3.3.1 理论基础

**广度优先搜索 (BFS)** 是一种图遍历算法，优先访问近邻节点。

**形式化定义**：
$```latex
\text{BFS}(G, v) = \text{visit}(v) \cup \bigcup_{u \in N(v)} \text{BFS}(G, u)
```$

其中 ```latex
N(v)
``` 是节点 ```latex
v
``` 的邻居集合。

#### 3.3.2 Go语言实现

```go
// BFS 广度优先搜索
func (g *Graph[T]) BFS(start T) []T {
 visited := make(map[T]bool)
 result := make([]T, 0)
 queue := []T{start}
 visited[start] = true

 for len(queue) > 0 {
  node := queue[0]
  queue = queue[1:]
  
  result = append(result, node)
  
  for _, neighbor := range g.adjacency[node] {
   if !visited[neighbor] {
    visited[neighbor] = true
    queue = append(queue, neighbor)
   }
  }
 }

 return result
}

// BFSWithLevel BFS带层级信息
func (g *Graph[T]) BFSWithLevel(start T) map[T]int {
 visited := make(map[T]bool)
 levels := make(map[T]int)
 queue := []T{start}
 visited[start] = true
 levels[start] = 0

 for len(queue) > 0 {
  node := queue[0]
  queue = queue[1:]
  currentLevel := levels[node]
  
  for _, neighbor := range g.adjacency[node] {
   if !visited[neighbor] {
    visited[neighbor] = true
    levels[neighbor] = currentLevel + 1
    queue = append(queue, neighbor)
   }
  }
 }

 return levels
}
```

## 4. 动态规划

### 4.1 斐波那契数列

#### 4.1.1 理论基础

**斐波那契数列** 是一个经典的动态规划问题：
$F(n) = \begin{cases}
0 & \text{if } n = 0 \\
1 & \text{if } n = 1 \\
F(n-1) + F(n-2) & \text{if } n > 1
\end{cases}$

#### 4.1.2 Go语言实现

```go
// Fibonacci 递归实现
func Fibonacci(n int) int {
 if n <= 1 {
  return n
 }
 return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciDP 动态规划实现
func FibonacciDP(n int) int {
 if n <= 1 {
  return n
 }

 dp := make([]int, n+1)
 dp[0] = 0
 dp[1] = 1

 for i := 2; i <= n; i++ {
  dp[i] = dp[i-1] + dp[i-2]
 }

 return dp[n]
}

// FibonacciOptimized 空间优化实现
func FibonacciOptimized(n int) int {
 if n <= 1 {
  return n
 }

 prev, curr := 0, 1
 for i := 2; i <= n; i++ {
  prev, curr = curr, prev+curr
 }

 return curr
}

// FibonacciMatrix 矩阵快速幂实现
func FibonacciMatrix(n int) int {
 if n <= 1 {
  return n
 }

 matrix := [2][2]int{{1, 1}, {1, 0}}
 result := matrixPower(matrix, n-1)
 return result[0][0]
}

// matrixPower 矩阵快速幂
func matrixPower(matrix [2][2]int, n int) [2][2]int {
 if n == 0 {
  return [2][2]int{{1, 0}, {0, 1}}
 }

 if n == 1 {
  return matrix
 }

 if n%2 == 0 {
  half := matrixPower(matrix, n/2)
  return matrixMultiply(half, half)
 } else {
  half := matrixPower(matrix, n/2)
  return matrixMultiply(matrixMultiply(half, half), matrix)
 }
}

// matrixMultiply 矩阵乘法
func matrixMultiply(a, b [2][2]int) [2][2]int {
 return [2][2]int{
  {a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
  {a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
 }
}
```

### 4.2 最长公共子序列

#### 4.2.1 理论基础

**最长公共子序列 (LCS)** 问题：
给定两个序列 ```latex
X
``` 和 ```latex
Y
```，找到它们的最长公共子序列。

**动态规划方程**：
$LCS[i][j] = \begin{cases}
0 & \text{if } i = 0 \text{ or } j = 0 \\
LCS[i-1][j-1] + 1 & \text{if } X[i-1] = Y[j-1] \\
\max(LCS[i-1][j], LCS[i][j-1]) & \text{otherwise}
\end{cases}$

#### 4.2.2 Go语言实现

```go
// LCS 最长公共子序列
func LCS(s1, s2 string) string {
 m, n := len(s1), len(s2)
 dp := make([][]int, m+1)
 for i := range dp {
  dp[i] = make([]int, n+1)
 }

 // 填充DP表
 for i := 1; i <= m; i++ {
  for j := 1; j <= n; j++ {
   if s1[i-1] == s2[j-1] {
    dp[i][j] = dp[i-1][j-1] + 1
   } else {
    dp[i][j] = max(dp[i-1][j], dp[i][j-1])
   }
  }
 }

 // 回溯构造结果
 return lcsBacktrack(dp, s1, s2, m, n)
}

// lcsBacktrack 回溯构造LCS
func lcsBacktrack(dp [][]int, s1, s2 string, i, j int) string {
 if i == 0 || j == 0 {
  return ""
 }

 if s1[i-1] == s2[j-1] {
  return lcsBacktrack(dp, s1, s2, i-1, j-1) + string(s1[i-1])
 }

 if dp[i-1][j] > dp[i][j-1] {
  return lcsBacktrack(dp, s1, s2, i-1, j)
 } else {
  return lcsBacktrack(dp, s1, s2, i, j-1)
 }
}

// max 最大值函数
func max(a, b int) int {
 if a > b {
  return a
 }
 return b
}

// LCSLength 只计算LCS长度
func LCSLength(s1, s2 string) int {
 m, n := len(s1), len(s2)
 dp := make([][]int, m+1)
 for i := range dp {
  dp[i] = make([]int, n+1)
 }

 for i := 1; i <= m; i++ {
  for j := 1; j <= n; j++ {
   if s1[i-1] == s2[j-1] {
    dp[i][j] = dp[i-1][j-1] + 1
   } else {
    dp[i][j] = max(dp[i-1][j], dp[i][j-1])
   }
  }
 }

 return dp[m][n]
}
```

## 5. 贪心算法

### 5.1 活动选择问题

#### 5.1.1 理论基础

**活动选择问题**：给定一组活动，每个活动有开始时间和结束时间，选择最多的不重叠活动。

**贪心策略**：按结束时间排序，选择结束时间最早的活动。

#### 5.1.2 Go语言实现

```go
// Activity 活动结构
type Activity struct {
 Start int
 End   int
}

// ActivitySelection 活动选择
func ActivitySelection(activities []Activity) []Activity {
 if len(activities) == 0 {
  return nil
 }

 // 按结束时间排序
 sort.Slice(activities, func(i, j int) bool {
  return activities[i].End < activities[j].End
 })

 result := []Activity{activities[0]}
 lastEnd := activities[0].End

 for i := 1; i < len(activities); i++ {
  if activities[i].Start >= lastEnd {
   result = append(result, activities[i])
   lastEnd = activities[i].End
  }
 }

 return result
}

// ActivitySelectionGeneric 泛型实现
type Selectable interface {
 GetStart() int
 GetEnd() int
}

func ActivitySelectionGeneric[T Selectable](activities []T) []T {
 if len(activities) == 0 {
  return nil
 }

 // 按结束时间排序
 sort.Slice(activities, func(i, j int) bool {
  return activities[i].GetEnd() < activities[j].GetEnd()
 })

 result := []T{activities[0]}
 lastEnd := activities[0].GetEnd()

 for i := 1; i < len(activities); i++ {
  if activities[i].GetStart() >= lastEnd {
   result = append(result, activities[i])
   lastEnd = activities[i].GetEnd()
  }
 }

 return result
}
```

## 6. 总结

### 6.1 算法复杂度对比

| 算法 | 时间复杂度 | 空间复杂度 | 稳定性 |
|------|------------|------------|--------|
| 快速排序 | ```latex
O(n \log n)
``` | ```latex
O(\log n)
``` | 不稳定 |
| 归并排序 | ```latex
O(n \log n)
``` | ```latex
O(n)
``` | 稳定 |
| 二分搜索 | ```latex
O(\log n)
``` | ```latex
O(1)
``` | - |
| DFS | ```latex
O(V + E)
``` | ```latex
O(V)
``` | - |
| BFS | ```latex
O(V + E)
``` | ```latex
O(V)
``` | - |

### 6.2 最佳实践

1. **选择合适的算法**：根据问题规模和要求选择
2. **优化实现**：使用适当的数据结构和优化技巧
3. **考虑边界情况**：处理空输入、单元素等特殊情况
4. **性能测试**：使用基准测试验证性能

---

**参考文献**：
1. Cormen, T. H., Leiserson, C. E., Rivest, R. L., & Stein, C. (2009). Introduction to Algorithms
2. Knuth, D. E. (1997). The Art of Computer Programming
3. Sedgewick, R., & Wayne, K. (2011). Algorithms
