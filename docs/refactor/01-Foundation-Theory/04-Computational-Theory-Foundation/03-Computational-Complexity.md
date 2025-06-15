# 03-计算复杂性 (Computational Complexity)

## 目录

1. [基础概念](#1-基础概念)
2. [时间复杂度](#2-时间复杂度)
3. [空间复杂度](#3-空间复杂度)
4. [复杂性类](#4-复杂性类)
5. [Go语言实现](#5-go语言实现)
6. [定理证明](#6-定理证明)
7. [应用示例](#7-应用示例)

## 1. 基础概念

### 1.1 计算复杂性概述

计算复杂性理论研究算法和问题的计算资源需求：

- **时间复杂度**：算法执行所需的时间
- **空间复杂度**：算法执行所需的存储空间
- **复杂性类**：具有相似复杂性的问题集合
- **应用领域**：算法设计、密码学、人工智能

### 1.2 基本定义

**定义 1.1** (计算模型)

```latex
图灵机 M = (Q, Σ, Γ, δ, q₀, q_accept, q_reject) 其中：

Q: 有限状态集
Σ: 输入字母表
Γ: 带字母表（Σ ⊆ Γ）
δ: Q × Γ → Q × Γ × {L, R} 转移函数
q₀: 初始状态
q_accept: 接受状态
q_reject: 拒绝状态
```

**定义 1.2** (时间复杂度)

```latex
对于图灵机 M 和输入 w，M 在 w 上的运行时间 t_M(w) 是 M 停机前执行的步数。

对于函数 f: N → N，我们说 M 的时间复杂度是 O(f(n))，如果存在常数 c > 0，使得对于所有长度为 n 的输入 w，有：
t_M(w) ≤ c · f(n)
```

**定义 1.3** (空间复杂度)

```latex
对于图灵机 M 和输入 w，M 在 w 上使用的空间 s_M(w) 是 M 停机前访问的带方格数。

对于函数 f: N → N，我们说 M 的空间复杂度是 O(f(n))，如果存在常数 c > 0，使得对于所有长度为 n 的输入 w，有：
s_M(w) ≤ c · f(n)
```

## 2. 时间复杂度

### 2.1 渐近记号

**定义 2.1** (大O记号)

```latex
对于函数 f, g: N → N，我们说 f(n) = O(g(n))，如果存在常数 c > 0 和 n₀ ∈ N，使得对于所有 n ≥ n₀，有：
f(n) ≤ c · g(n)
```

**定义 2.2** (大Ω记号)

```latex
对于函数 f, g: N → N，我们说 f(n) = Ω(g(n))，如果存在常数 c > 0 和 n₀ ∈ N，使得对于所有 n ≥ n₀，有：
f(n) ≥ c · g(n)
```

**定义 2.3** (大Θ记号)

```latex
对于函数 f, g: N → N，我们说 f(n) = Θ(g(n))，如果：
f(n) = O(g(n)) 且 f(n) = Ω(g(n))
```

### 2.2 常见时间复杂度类

**定义 2.4** (多项式时间)

```latex
P = {L | 存在确定性图灵机 M 和多项式 p，使得 M 在时间 O(p(n)) 内判定 L}
```

**定义 2.5** (非确定性多项式时间)

```latex
NP = {L | 存在非确定性图灵机 M 和多项式 p，使得 M 在时间 O(p(n)) 内判定 L}
```

**定义 2.6** (指数时间)

```latex
EXP = {L | 存在确定性图灵机 M 和多项式 p，使得 M 在时间 O(2^p(n)) 内判定 L}
```

## 3. 空间复杂度

### 3.1 空间复杂性类

**定义 3.1** (对数空间)

```latex
L = {L | 存在确定性图灵机 M，使得 M 在空间 O(log n) 内判定 L}
```

**定义 3.2** (多项式空间)

```latex
PSPACE = {L | 存在确定性图灵机 M 和多项式 p，使得 M 在空间 O(p(n)) 内判定 L}
```

**定义 3.3** (非确定性对数空间)

```latex
NL = {L | 存在非确定性图灵机 M，使得 M 在空间 O(log n) 内判定 L}
```

### 3.2 空间层次定理

**定理 3.1** (空间层次定理)

```latex
对于空间可构造函数 f, g: N → N，如果 f(n) = o(g(n))，则：
SPACE(f(n)) ⊊ SPACE(g(n))
```

## 4. 复杂性类

### 4.1 基本复杂性类

**定义 4.1** (复杂性类层次)

```latex
L ⊆ NL ⊆ P ⊆ NP ⊆ PSPACE ⊆ EXP ⊆ EXPSPACE
```

**定义 4.2** (完全性问题)

```latex
对于复杂性类 C 和语言 L ∈ C，我们说 L 是 C-完全的，如果：
1. L ∈ C
2. 对于所有 L' ∈ C，有 L' ≤_p L（多项式时间归约）
```

### 4.2 NP完全性

**定义 4.3** (NP完全性)

```latex
语言 L 是NP完全的，如果：
1. L ∈ NP
2. 对于所有 L' ∈ NP，有 L' ≤_p L
```

**定理 4.1** (Cook-Levin定理)

```latex
SAT 是NP完全的
```

## 5. Go语言实现

### 5.1 复杂性分析框架

```go
package complexity

import (
 "fmt"
 "math"
 "time"
)

// ComplexityAnalyzer 复杂性分析器
type ComplexityAnalyzer struct {
 Algorithm Algorithm
 Inputs    []Input
}

// Algorithm 算法接口
type Algorithm interface {
 Execute(input Input) Result
 Name() string
}

// Input 输入接口
type Input interface {
 Size() int
 Data() interface{}
}

// Result 结果接口
type Result interface {
 Output() interface{}
 Time() time.Duration
 Space() int
}

// TimeComplexity 时间复杂度分析
type TimeComplexity struct {
 Algorithm Algorithm
 Inputs    []Input
 Results   []TimeResult
}

// TimeResult 时间结果
type TimeResult struct {
 InputSize int
 Time      time.Duration
 Steps     int
}

// SpaceComplexity 空间复杂度分析
type SpaceComplexity struct {
 Algorithm Algorithm
 Inputs    []Input
 Results   []SpaceResult
}

// SpaceResult 空间结果
type SpaceResult struct {
 InputSize int
 Space     int
 PeakSpace int
}

// AsymptoticNotation 渐近记号
type AsymptoticNotation struct {
 Type string // "O", "Ω", "Θ"
 Function string
}

// ComplexityClass 复杂性类
type ComplexityClass struct {
 Name        string
 Description string
 Algorithms  []Algorithm
}
```

### 5.2 时间复杂度分析

```go
// AnalyzeTimeComplexity 分析时间复杂度
func (tca *TimeComplexity) AnalyzeTimeComplexity() *AsymptoticNotation {
 // 收集数据点
 var sizes []int
 var times []float64
 
 for _, result := range tca.Results {
  sizes = append(sizes, result.InputSize)
  times = append(times, float64(result.Time.Nanoseconds()))
 }
 
 // 分析渐近行为
 notation := tca.analyzeAsymptoticBehavior(sizes, times)
 return notation
}

// analyzeAsymptoticBehavior 分析渐近行为
func (tca *TimeComplexity) analyzeAsymptoticBehavior(sizes []int, times []float64) *AsymptoticNotation {
 // 尝试不同的函数类型
 functions := []struct {
  name string
  fn   func(int) float64
 }{
  {"O(1)", func(n int) float64 { return 1 }},
  {"O(log n)", func(n int) float64 { return math.Log(float64(n)) }},
  {"O(n)", func(n int) float64 { return float64(n) }},
  {"O(n log n)", func(n int) float64 { return float64(n) * math.Log(float64(n)) }},
  {"O(n²)", func(n int) float64 { return float64(n * n) }},
  {"O(2ⁿ)", func(n int) float64 { return math.Pow(2, float64(n)) }},
 }
 
 bestFit := ""
 bestError := math.MaxFloat64
 
 for _, fn := range functions {
  error := tca.calculateError(sizes, times, fn.fn)
  if error < bestError {
   bestError = error
   bestFit = fn.name
  }
 }
 
 return &AsymptoticNotation{
  Type:     "O",
  Function: bestFit,
 }
}

// calculateError 计算拟合误差
func (tca *TimeComplexity) calculateError(sizes []int, times []float64, fn func(int) float64) float64 {
 var totalError float64
 
 for i, size := range sizes {
  predicted := fn(size)
  actual := times[i]
  
  // 归一化误差
  if actual > 0 {
   error := math.Abs(predicted - actual) / actual
   totalError += error
  }
 }
 
 return totalError / float64(len(sizes))
}

// BenchmarkAlgorithm 基准测试算法
func (tca *TimeComplexity) BenchmarkAlgorithm() {
 tca.Results = make([]TimeResult, 0)
 
 for _, input := range tca.Inputs {
  start := time.Now()
  
  // 执行算法
  result := tca.Algorithm.Execute(input)
  
  duration := time.Since(start)
  
  timeResult := TimeResult{
   InputSize: input.Size(),
   Time:      duration,
   Steps:     tca.estimateSteps(input.Size()),
  }
  
  tca.Results = append(tca.Results, timeResult)
 }
}

// estimateSteps 估计步数
func (tca *TimeComplexity) estimateSteps(inputSize int) int {
 // 基于算法名称估计步数
 switch tca.Algorithm.Name() {
 case "LinearSearch":
  return inputSize
 case "BinarySearch":
  return int(math.Log2(float64(inputSize)))
 case "BubbleSort":
  return inputSize * inputSize
 case "QuickSort":
  return inputSize * int(math.Log2(float64(inputSize)))
 default:
  return inputSize
 }
}
```

### 5.3 空间复杂度分析

```go
// AnalyzeSpaceComplexity 分析空间复杂度
func (sca *SpaceComplexity) AnalyzeSpaceComplexity() *AsymptoticNotation {
 // 收集数据点
 var sizes []int
 var spaces []float64
 
 for _, result := range sca.Results {
  sizes = append(sizes, result.InputSize)
  spaces = append(spaces, float64(result.PeakSpace))
 }
 
 // 分析渐近行为
 notation := sca.analyzeAsymptoticBehavior(sizes, spaces)
 return notation
}

// analyzeAsymptoticBehavior 分析渐近行为
func (sca *SpaceComplexity) analyzeAsymptoticBehavior(sizes []int, spaces []float64) *AsymptoticNotation {
 // 尝试不同的函数类型
 functions := []struct {
  name string
  fn   func(int) float64
 }{
  {"O(1)", func(n int) float64 { return 1 }},
  {"O(log n)", func(n int) float64 { return math.Log(float64(n)) }},
  {"O(n)", func(n int) float64 { return float64(n) }},
  {"O(n log n)", func(n int) float64 { return float64(n) * math.Log(float64(n)) }},
  {"O(n²)", func(n int) float64 { return float64(n * n) }},
 }
 
 bestFit := ""
 bestError := math.MaxFloat64
 
 for _, fn := range functions {
  error := sca.calculateError(sizes, spaces, fn.fn)
  if error < bestError {
   bestError = error
   bestFit = fn.name
  }
 }
 
 return &AsymptoticNotation{
  Type:     "O",
  Function: bestFit,
 }
}

// calculateError 计算拟合误差
func (sca *SpaceComplexity) calculateError(sizes []int, spaces []float64, fn func(int) float64) float64 {
 var totalError float64
 
 for i, size := range sizes {
  predicted := fn(size)
  actual := spaces[i]
  
  // 归一化误差
  if actual > 0 {
   error := math.Abs(predicted - actual) / actual
   totalError += error
  }
 }
 
 return totalError / float64(len(sizes))
}

// MonitorSpaceUsage 监控空间使用
func (sca *SpaceComplexity) MonitorSpaceUsage() {
 sca.Results = make([]SpaceResult, 0)
 
 for _, input := range sca.Inputs {
  // 记录初始内存使用
  initialSpace := sca.getCurrentMemoryUsage()
  
  // 执行算法
  result := sca.Algorithm.Execute(input)
  
  // 记录峰值内存使用
  peakSpace := sca.getPeakMemoryUsage()
  finalSpace := sca.getCurrentMemoryUsage()
  
  spaceResult := SpaceResult{
   InputSize: input.Size(),
   Space:     finalSpace - initialSpace,
   PeakSpace: peakSpace - initialSpace,
  }
  
  sca.Results = append(sca.Results, spaceResult)
 }
}

// getCurrentMemoryUsage 获取当前内存使用
func (sca *SpaceComplexity) getCurrentMemoryUsage() int {
 // 简化的内存使用估算
 // 实际实现需要使用 runtime.ReadMemStats
 return 1024 // 1KB 基础使用
}

// getPeakMemoryUsage 获取峰值内存使用
func (sca *SpaceComplexity) getPeakMemoryUsage() int {
 // 简化的峰值内存估算
 return 2048 // 2KB 峰值使用
}
```

### 5.4 算法实现示例

```go
// LinearSearch 线性搜索算法
type LinearSearch struct{}

func (ls *LinearSearch) Execute(input Input) Result {
 data := input.Data().([]int)
 target := data[len(data)-1] // 假设最后一个元素是目标
 
 start := time.Now()
 
 // 线性搜索
 found := false
 for i, value := range data[:len(data)-1] {
  if value == target {
   found = true
   break
  }
 }
 
 duration := time.Since(start)
 
 return &SearchResult{
  output: found,
  time:   duration,
  space:  len(data) * 8, // 假设每个int占8字节
 }
}

func (ls *LinearSearch) Name() string {
 return "LinearSearch"
}

// BinarySearch 二分搜索算法
type BinarySearch struct{}

func (bs *BinarySearch) Execute(input Input) Result {
 data := input.Data().([]int)
 target := data[len(data)-1] // 假设最后一个元素是目标
 
 start := time.Now()
 
 // 二分搜索（假设数据已排序）
 left, right := 0, len(data)-2
 found := false
 
 for left <= right {
  mid := (left + right) / 2
  if data[mid] == target {
   found = true
   break
  } else if data[mid] < target {
   left = mid + 1
  } else {
   right = mid - 1
  }
 }
 
 duration := time.Since(start)
 
 return &SearchResult{
  output: found,
  time:   duration,
  space:  len(data) * 8, // 假设每个int占8字节
 }
}

func (bs *BinarySearch) Name() string {
 return "BinarySearch"
}

// BubbleSort 冒泡排序算法
type BubbleSort struct{}

func (bbs *BubbleSort) Execute(input Input) Result {
 data := input.Data().([]int)
 
 start := time.Now()
 
 // 冒泡排序
 n := len(data)
 for i := 0; i < n-1; i++ {
  for j := 0; j < n-i-1; j++ {
   if data[j] > data[j+1] {
    data[j], data[j+1] = data[j+1], data[j]
   }
  }
 }
 
 duration := time.Since(start)
 
 return &SortResult{
  output: data,
  time:   duration,
  space:  len(data) * 8, // 假设每个int占8字节
 }
}

func (bbs *BubbleSort) Name() string {
 return "BubbleSort"
}

// QuickSort 快速排序算法
type QuickSort struct{}

func (qs *QuickSort) Execute(input Input) Result {
 data := input.Data().([]int)
 
 start := time.Now()
 
 // 快速排序
 qs.quickSort(data, 0, len(data)-1)
 
 duration := time.Since(start)
 
 return &SortResult{
  output: data,
  time:   duration,
  space:  len(data) * 8, // 假设每个int占8字节
 }
}

func (qs *QuickSort) quickSort(arr []int, low, high int) {
 if low < high {
  pi := qs.partition(arr, low, high)
  qs.quickSort(arr, low, pi-1)
  qs.quickSort(arr, pi+1, high)
 }
}

func (qs *QuickSort) partition(arr []int, low, high int) int {
 pivot := arr[high]
 i := low - 1
 
 for j := low; j < high; j++ {
  if arr[j] < pivot {
   i++
   arr[i], arr[j] = arr[j], arr[i]
  }
 }
 
 arr[i+1], arr[high] = arr[high], arr[i+1]
 return i + 1
}

func (qs *QuickSort) Name() string {
 return "QuickSort"
}

// 结果类型
type SearchResult struct {
 output interface{}
 time   time.Duration
 space  int
}

func (sr *SearchResult) Output() interface{} { return sr.output }
func (sr *SearchResult) Time() time.Duration { return sr.time }
func (sr *SearchResult) Space() int          { return sr.space }

type SortResult struct {
 output interface{}
 time   time.Duration
 space  int
}

func (sr *SortResult) Output() interface{} { return sr.output }
func (sr *SortResult) Time() time.Duration { return sr.time }
func (sr *SortResult) Space() int          { return sr.space }
```

## 6. 定理证明

### 6.1 时间层次定理

**定理 6.1** (时间层次定理)

```latex
对于时间可构造函数 f, g: N → N，如果 f(n) = o(g(n))，则：
TIME(f(n)) ⊊ TIME(g(n))
```

**证明**：

```latex
使用对角线化方法构造语言 L，使得：
L ∈ TIME(g(n)) 但 L ∉ TIME(f(n))

构造图灵机 M，对于输入 w：
1. 模拟图灵机 M_w 在输入 w 上运行 g(|w|) 步
2. 如果 M_w 接受，则拒绝
3. 如果 M_w 拒绝或未停机，则接受

L = {w | M 接受 w}

显然 L ∈ TIME(g(n))
假设 L ∈ TIME(f(n))，则存在图灵机 M' 在时间 O(f(n)) 内判定 L
但 M' 在输入 M' 上的行为与 M 的定义矛盾
因此 L ∉ TIME(f(n))
```

### 6.2 P ≠ EXP

**定理 6.2** (P ≠ EXP)

```latex
P ⊊ EXP
```

**证明**：

```latex
由时间层次定理，对于任何多项式 p(n)，有：
TIME(p(n)) ⊊ TIME(2^p(n))

因此：
P = ∪_p TIME(p(n)) ⊊ ∪_p TIME(2^p(n)) = EXP
```

### 6.3 NP完全性传递性

**定理 6.3** (NP完全性传递性)

```latex
如果 L₁ 是NP完全的，L₁ ≤_p L₂，且 L₂ ∈ NP，则 L₂ 是NP完全的
```

**证明**：

```latex
对于任意 L ∈ NP，有：
L ≤_p L₁ ≤_p L₂

由于多项式时间归约的传递性，L ≤_p L₂
因此 L₂ 是NP完全的
```

## 7. 应用示例

### 7.1 算法性能分析

```go
// AlgorithmPerformanceAnalysis 算法性能分析
func AlgorithmPerformanceAnalysis() {
 // 创建测试输入
 inputs := []Input{
  &IntArrayInput{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, target: 5},
  &IntArrayInput{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, target: 8},
  &IntArrayInput{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, target: 12},
 }
 
 // 分析线性搜索
 linearSearch := &LinearSearch{}
 linearAnalyzer := &TimeComplexity{
  Algorithm: linearSearch,
  Inputs:    inputs,
 }
 
 linearAnalyzer.BenchmarkAlgorithm()
 linearNotation := linearAnalyzer.AnalyzeTimeComplexity()
 fmt.Printf("线性搜索时间复杂度: %s\n", linearNotation.Function)
 
 // 分析二分搜索
 binarySearch := &BinarySearch{}
 binaryAnalyzer := &TimeComplexity{
  Algorithm: binarySearch,
  Inputs:    inputs,
 }
 
 binaryAnalyzer.BenchmarkAlgorithm()
 binaryNotation := binaryAnalyzer.AnalyzeTimeComplexity()
 fmt.Printf("二分搜索时间复杂度: %s\n", binaryNotation.Function)
}
```

### 7.2 排序算法比较

```go
// SortingAlgorithmComparison 排序算法比较
func SortingAlgorithmComparison() {
 // 创建测试输入
 inputs := []Input{
  &IntArrayInput{data: []int{5, 2, 8, 1, 9, 3, 7, 4, 6}},
  &IntArrayInput{data: []int{10, 5, 2, 8, 1, 9, 3, 7, 4, 6, 11, 12, 13, 14, 15}},
  &IntArrayInput{data: []int{20, 10, 5, 2, 8, 1, 9, 3, 7, 4, 6, 11, 12, 13, 14, 15, 16, 17, 18, 19}},
 }
 
 // 分析冒泡排序
 bubbleSort := &BubbleSort{}
 bubbleAnalyzer := &TimeComplexity{
  Algorithm: bubbleSort,
  Inputs:    inputs,
 }
 
 bubbleAnalyzer.BenchmarkAlgorithm()
 bubbleNotation := bubbleAnalyzer.AnalyzeTimeComplexity()
 fmt.Printf("冒泡排序时间复杂度: %s\n", bubbleNotation.Function)
 
 // 分析快速排序
 quickSort := &QuickSort{}
 quickAnalyzer := &TimeComplexity{
  Algorithm: quickSort,
  Inputs:    inputs,
 }
 
 quickAnalyzer.BenchmarkAlgorithm()
 quickNotation := quickAnalyzer.AnalyzeTimeComplexity()
 fmt.Printf("快速排序时间复杂度: %s\n", quickNotation.Function)
}
```

### 7.3 空间复杂度分析

```go
// SpaceComplexityAnalysis 空间复杂度分析
func SpaceComplexityAnalysis() {
 // 创建测试输入
 inputs := []Input{
  &IntArrayInput{data: []int{1, 2, 3, 4, 5}},
  &IntArrayInput{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
  &IntArrayInput{data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
 }
 
 // 分析空间复杂度
 spaceAnalyzer := &SpaceComplexity{
  Algorithm: &QuickSort{},
  Inputs:    inputs,
 }
 
 spaceAnalyzer.MonitorSpaceUsage()
 spaceNotation := spaceAnalyzer.AnalyzeSpaceComplexity()
 fmt.Printf("快速排序空间复杂度: %s\n", spaceNotation.Function)
}

// IntArrayInput 整数数组输入
type IntArrayInput struct {
 data   []int
 target int
}

func (iai *IntArrayInput) Size() int {
 return len(iai.data)
}

func (iai *IntArrayInput) Data() interface{} {
 return iai.data
}
```

## 总结

计算复杂性理论为软件工程提供了重要的理论基础，能够：

1. **算法设计**：指导高效算法的设计
2. **性能分析**：提供算法性能的理论保证
3. **问题分类**：将问题按复杂性进行分类
4. **优化决策**：帮助选择最优的算法和数据结构

通过Go语言的实现，我们可以将计算复杂性理论应用到实际的软件工程问题中，提供性能分析和优化指导。
