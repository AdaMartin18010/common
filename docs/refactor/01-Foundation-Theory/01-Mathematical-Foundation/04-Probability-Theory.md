# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 概率空间](#11-概率空间)
    - [1.2 事件](#12-事件)
    - [1.3 概率公理](#13-概率公理)
  - [2. 随机变量](#2-随机变量)
    - [2.1 离散随机变量](#21-离散随机变量)
    - [2.2 连续随机变量](#22-连续随机变量)
    - [2.3 期望和方差](#23-期望和方差)
  - [3. 概率分布](#3-概率分布)
    - [3.1 常见离散分布](#31-常见离散分布)
    - [3.2 常见连续分布](#32-常见连续分布)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 概率空间实现](#41-概率空间实现)
    - [4.2 随机变量实现](#42-随机变量实现)
    - [4.3 概率分布实现](#43-概率分布实现)
  - [5. 应用示例](#5-应用示例)
    - [5.1 蒙特卡洛模拟](#51-蒙特卡洛模拟)
    - [5.2 随机算法](#52-随机算法)
    - [5.3 机器学习](#53-机器学习)

## 1. 基本概念

### 1.1 概率空间

**定义 1.1**: 概率空间 $(\Omega, \mathcal{F}, P)$ 由以下三个部分组成：

- **样本空间** $\Omega$: 所有可能结果的集合
- **事件域** $\mathcal{F}$: $\Omega$ 的子集族，满足 $\sigma$-代数性质
- **概率测度** $P$: $\mathcal{F} \rightarrow [0,1]$ 的函数

**形式化表达**:
- 样本空间：$\Omega = \{\omega_1, \omega_2, \ldots, \omega_n\}$
- 事件：$A \in \mathcal{F} \subseteq 2^\Omega$
- 概率：$P(A) \in [0,1]$

### 1.2 事件

**定义 1.2**: 事件的基本概念

1. **基本事件**: $\{\omega\}$ 其中 $\omega \in \Omega$
2. **复合事件**: 基本事件的并集
3. **必然事件**: $\Omega$
4. **不可能事件**: $\emptyset$

**定义 1.3**: 事件的运算

- **并集**: $A \cup B = \{\omega \mid \omega \in A \text{ 或 } \omega \in B\}$
- **交集**: $A \cap B = \{\omega \mid \omega \in A \text{ 且 } \omega \in B\}$
- **补集**: $A^c = \{\omega \mid \omega \notin A\}$

### 1.3 概率公理

**公理 1.1** (非负性): $P(A) \geq 0$ 对所有 $A \in \mathcal{F}$

**公理 1.2** (规范性): $P(\Omega) = 1$

**公理 1.3** (可列可加性): 对于互不相容的事件序列 $\{A_i\}$：

$$P\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} P(A_i)$$

## 2. 随机变量

### 2.1 离散随机变量

**定义 2.1**: 离散随机变量 $X$ 是样本空间到实数的函数：

$$X: \Omega \rightarrow \mathbb{R}$$

**定义 2.2**: 概率质量函数 (PMF)：

$$p_X(x) = P(X = x) = P(\{\omega \mid X(\omega) = x\})$$

**性质**:
- $p_X(x) \geq 0$ 对所有 $x$
- $\sum_x p_X(x) = 1$

### 2.2 连续随机变量

**定义 2.3**: 连续随机变量 $X$ 的概率密度函数 (PDF)：

$$f_X(x) = \frac{d}{dx} F_X(x)$$

其中 $F_X(x) = P(X \leq x)$ 是累积分布函数 (CDF)。

**性质**:
- $f_X(x) \geq 0$ 对所有 $x$
- $\int_{-\infty}^{\infty} f_X(x) dx = 1$

### 2.3 期望和方差

**定义 2.4**: 期望

对于离散随机变量：
$$E[X] = \sum_x x \cdot p_X(x)$$

对于连续随机变量：
$$E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$$

**定义 2.5**: 方差

$$\text{Var}(X) = E[(X - E[X])^2] = E[X^2] - (E[X])^2$$

## 3. 概率分布

### 3.1 常见离散分布

**定义 3.1** (伯努利分布): $X \sim \text{Bernoulli}(p)$

$$p_X(x) = \begin{cases}
p & \text{if } x = 1 \\
1-p & \text{if } x = 0
\end{cases}$$

**定义 3.2** (二项分布): $X \sim \text{Binomial}(n,p)$

$$p_X(x) = \binom{n}{x} p^x (1-p)^{n-x}$$

**定义 3.3** (泊松分布): $X \sim \text{Poisson}(\lambda)$

$$p_X(x) = \frac{e^{-\lambda} \lambda^x}{x!}$$

### 3.2 常见连续分布

**定义 3.4** (正态分布): $X \sim \text{Normal}(\mu, \sigma^2)$

$$f_X(x) = \frac{1}{\sqrt{2\pi\sigma^2}} e^{-\frac{(x-\mu)^2}{2\sigma^2}}$$

**定义 3.5** (指数分布): $X \sim \text{Exponential}(\lambda)$

$$f_X(x) = \lambda e^{-\lambda x} \text{ for } x \geq 0$$

## 4. Go语言实现

### 4.1 概率空间实现

```go
// SampleSpace 样本空间
type SampleSpace[T comparable] struct {
    Elements Set[T]
    Events   Set[Set[T]]
}

// ProbabilitySpace 概率空间
type ProbabilitySpace[T comparable] struct {
    SampleSpace SampleSpace[T]
    Probability map[Set[T]]float64
}

// NewProbabilitySpace 创建概率空间
func NewProbabilitySpace[T comparable]() *ProbabilitySpace[T] {
    return &ProbabilitySpace[T]{
        SampleSpace: SampleSpace[T]{
            Elements: NewSet[T](),
            Events:   NewSet[Set[T]](),
        },
        Probability: make(map[Set[T]]float64),
    }
}

// AddElement 添加样本点
func (ps *ProbabilitySpace[T]) AddElement(element T) {
    ps.SampleSpace.Elements.Add(element)
}

// AddEvent 添加事件
func (ps *ProbabilitySpace[T]) AddEvent(event Set[T], probability float64) {
    ps.SampleSpace.Events.Add(event)
    ps.Probability[event] = probability
}

// GetProbability 获取事件概率
func (ps *ProbabilitySpace[T]) GetProbability(event Set[T]) float64 {
    return ps.Probability[event]
}

// Union 事件并集
func (ps *ProbabilitySpace[T]) Union(event1, event2 Set[T]) Set[T] {
    return event1.Union(event2)
}

// Intersection 事件交集
func (ps *ProbabilitySpace[T]) Intersection(event1, event2 Set[T]) Set[T] {
    return event1.Intersection(event2)
}

// Complement 事件补集
func (ps *ProbabilitySpace[T]) Complement(event Set[T]) Set[T] {
    return ps.SampleSpace.Elements.Difference(event)
}

// ConditionalProbability 条件概率
func (ps *ProbabilitySpace[T]) ConditionalProbability(event, condition Set[T]) float64 {
    if ps.Probability[condition] == 0 {
        return 0
    }
    
    intersection := ps.Intersection(event, condition)
    return ps.Probability[intersection] / ps.Probability[condition]
}
```

### 4.2 随机变量实现

```go
// RandomVariable 随机变量接口
type RandomVariable[T comparable, R any] interface {
    Sample() R
    GetDistribution() Distribution[R]
}

// DiscreteRandomVariable 离散随机变量
type DiscreteRandomVariable[T comparable] struct {
    SampleSpace Set[T]
    PMF         map[T]float64
}

// NewDiscreteRandomVariable 创建离散随机变量
func NewDiscreteRandomVariable[T comparable]() *DiscreteRandomVariable[T] {
    return &DiscreteRandomVariable[T]{
        SampleSpace: NewSet[T](),
        PMF:         make(map[T]float64),
    }
}

// AddOutcome 添加结果
func (drv *DiscreteRandomVariable[T]) AddOutcome(outcome T, probability float64) {
    drv.SampleSpace.Add(outcome)
    drv.PMF[outcome] = probability
}

// Sample 采样
func (drv *DiscreteRandomVariable[T]) Sample() T {
    r := rand.Float64()
    cumulative := 0.0
    
    for outcome := range drv.SampleSpace {
        cumulative += drv.PMF[outcome]
        if r <= cumulative {
            return outcome
        }
    }
    
    // 返回最后一个结果（理论上不应该到达这里）
    var last T
    for outcome := range drv.SampleSpace {
        last = outcome
    }
    return last
}

// Expectation 期望
func (drv *DiscreteRandomVariable[T]) Expectation() float64 {
    expectation := 0.0
    for outcome := range drv.SampleSpace {
        expectation += float64(outcome) * drv.PMF[outcome]
    }
    return expectation
}

// Variance 方差
func (drv *DiscreteRandomVariable[T]) Variance() float64 {
    expectation := drv.Expectation()
    variance := 0.0
    
    for outcome := range drv.SampleSpace {
        diff := float64(outcome) - expectation
        variance += diff * diff * drv.PMF[outcome]
    }
    
    return variance
}

// ContinuousRandomVariable 连续随机变量
type ContinuousRandomVariable struct {
    PDF func(float64) float64
    CDF func(float64) float64
}

// NewContinuousRandomVariable 创建连续随机变量
func NewContinuousRandomVariable(pdf func(float64) float64, cdf func(float64) float64) *ContinuousRandomVariable {
    return &ContinuousRandomVariable{
        PDF: pdf,
        CDF: cdf,
    }
}

// Sample 采样（使用逆变换法）
func (crv *ContinuousRandomVariable) Sample() float64 {
    r := rand.Float64()
    
    // 使用二分查找找到CDF的逆
    left, right := -1000.0, 1000.0
    for right-left > 1e-6 {
        mid := (left + right) / 2
        if crv.CDF(mid) < r {
            left = mid
        } else {
            right = mid
        }
    }
    
    return (left + right) / 2
}

// Expectation 期望（数值积分）
func (crv *ContinuousRandomVariable) Expectation() float64 {
    return crv.integrate(func(x float64) float64 {
        return x * crv.PDF(x)
    })
}

// Variance 方差
func (crv *ContinuousRandomVariable) Variance() float64 {
    expectation := crv.Expectation()
    return crv.integrate(func(x float64) float64 {
        diff := x - expectation
        return diff * diff * crv.PDF(x)
    })
}

// integrate 数值积分
func (crv *ContinuousRandomVariable) integrate(f func(float64) float64) float64 {
    a, b := -10.0, 10.0
    n := 10000
    h := (b - a) / float64(n)
    
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + float64(i)*h
        sum += f(x)
    }
    
    return h * sum
}
```

### 4.3 概率分布实现

```go
// Distribution 分布接口
type Distribution[T any] interface {
    Sample() T
    PDF(x T) float64
    CDF(x T) float64
}

// BernoulliDistribution 伯努利分布
type BernoulliDistribution struct {
    P float64
}

// NewBernoulliDistribution 创建伯努利分布
func NewBernoulliDistribution(p float64) *BernoulliDistribution {
    return &BernoulliDistribution{P: p}
}

// Sample 采样
func (bd *BernoulliDistribution) Sample() bool {
    return rand.Float64() < bd.P
}

// PDF 概率质量函数
func (bd *BernoulliDistribution) PDF(x bool) float64 {
    if x {
        return bd.P
    }
    return 1 - bd.P
}

// CDF 累积分布函数
func (bd *BernoulliDistribution) CDF(x bool) float64 {
    if x {
        return 1.0
    }
    return 1 - bd.P
}

// BinomialDistribution 二项分布
type BinomialDistribution struct {
    N int
    P float64
}

// NewBinomialDistribution 创建二项分布
func NewBinomialDistribution(n int, p float64) *BinomialDistribution {
    return &BinomialDistribution{N: n, P: p}
}

// Sample 采样
func (bd *BinomialDistribution) Sample() int {
    count := 0
    for i := 0; i < bd.N; i++ {
        if rand.Float64() < bd.P {
            count++
        }
    }
    return count
}

// PDF 概率质量函数
func (bd *BinomialDistribution) PDF(x int) float64 {
    if x < 0 || x > bd.N {
        return 0
    }
    
    // 计算组合数
    comb := 1.0
    for i := 0; i < x; i++ {
        comb *= float64(bd.N-i) / float64(i+1)
    }
    
    return comb * math.Pow(bd.P, float64(x)) * math.Pow(1-bd.P, float64(bd.N-x))
}

// CDF 累积分布函数
func (bd *BinomialDistribution) CDF(x int) float64 {
    if x < 0 {
        return 0
    }
    if x >= bd.N {
        return 1
    }
    
    cdf := 0.0
    for i := 0; i <= x; i++ {
        cdf += bd.PDF(i)
    }
    return cdf
}

// NormalDistribution 正态分布
type NormalDistribution struct {
    Mu    float64
    Sigma float64
}

// NewNormalDistribution 创建正态分布
func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    return &NormalDistribution{Mu: mu, Sigma: sigma}
}

// Sample 采样（Box-Muller变换）
func (nd *NormalDistribution) Sample() float64 {
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return nd.Mu + nd.Sigma*z0
}

// PDF 概率密度函数
func (nd *NormalDistribution) PDF(x float64) float64 {
    exponent := -0.5 * math.Pow((x-nd.Mu)/nd.Sigma, 2)
    return math.Exp(exponent) / (nd.Sigma * math.Sqrt(2*math.Pi))
}

// CDF 累积分布函数
func (nd *NormalDistribution) CDF(x float64) float64 {
    z := (x - nd.Mu) / nd.Sigma
    return 0.5 * (1 + math.Erf(z/math.Sqrt(2)))
}

// ExponentialDistribution 指数分布
type ExponentialDistribution struct {
    Lambda float64
}

// NewExponentialDistribution 创建指数分布
func NewExponentialDistribution(lambda float64) *ExponentialDistribution {
    return &ExponentialDistribution{Lambda: lambda}
}

// Sample 采样
func (ed *ExponentialDistribution) Sample() float64 {
    return -math.Log(rand.Float64()) / ed.Lambda
}

// PDF 概率密度函数
func (ed *ExponentialDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return ed.Lambda * math.Exp(-ed.Lambda*x)
}

// CDF 累积分布函数
func (ed *ExponentialDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return 1 - math.Exp(-ed.Lambda*x)
}
```

## 5. 应用示例

### 5.1 蒙特卡洛模拟

```go
// MonteCarloSimulator 蒙特卡洛模拟器
type MonteCarloSimulator struct {
    random *rand.Rand
}

// NewMonteCarloSimulator 创建蒙特卡洛模拟器
func NewMonteCarloSimulator(seed int64) *MonteCarloSimulator {
    return &MonteCarloSimulator{
        random: rand.New(rand.NewSource(seed)),
    }
}

// EstimatePi 估算π值
func (mcs *MonteCarloSimulator) EstimatePi(iterations int) float64 {
    inside := 0
    
    for i := 0; i < iterations; i++ {
        x := mcs.random.Float64()
        y := mcs.random.Float64()
        
        if x*x+y*y <= 1 {
            inside++
        }
    }
    
    return 4.0 * float64(inside) / float64(iterations)
}

// EstimateIntegral 估算定积分
func (mcs *MonteCarloSimulator) EstimateIntegral(f func(float64) float64, a, b float64, iterations int) float64 {
    sum := 0.0
    
    for i := 0; i < iterations; i++ {
        x := a + mcs.random.Float64()*(b-a)
        sum += f(x)
    }
    
    return (b - a) * sum / float64(iterations)
}

// SimulateRandomWalk 模拟随机游走
func (mcs *MonteCarloSimulator) SimulateRandomWalk(steps int) []float64 {
    position := 0.0
    positions := make([]float64, steps+1)
    positions[0] = position
    
    for i := 0; i < steps; i++ {
        if mcs.random.Float64() < 0.5 {
            position += 1
        } else {
            position -= 1
        }
        positions[i+1] = position
    }
    
    return positions
}
```

### 5.2 随机算法

```go
// RandomizedAlgorithm 随机算法
type RandomizedAlgorithm struct {
    random *rand.Rand
}

// NewRandomizedAlgorithm 创建随机算法
func NewRandomizedAlgorithm(seed int64) *RandomizedAlgorithm {
    return &RandomizedAlgorithm{
        random: rand.New(rand.NewSource(seed)),
    }
}

// QuickSort 快速排序（随机化版本）
func (ra *RandomizedAlgorithm) QuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    // 随机选择主元
    pivotIndex := ra.random.Intn(len(arr))
    pivot := arr[pivotIndex]
    
    // 分区
    left := make([]int, 0)
    right := make([]int, 0)
    equal := make([]int, 0)
    
    for _, val := range arr {
        if val < pivot {
            left = append(left, val)
        } else if val > pivot {
            right = append(right, val)
        } else {
            equal = append(equal, val)
        }
    }
    
    // 递归排序
    left = ra.QuickSort(left)
    right = ra.QuickSort(right)
    
    // 合并结果
    result := append(left, equal...)
    result = append(result, right...)
    
    return result
}

// RandomizedSelection 随机化选择（找第k小元素）
func (ra *RandomizedAlgorithm) RandomizedSelection(arr []int, k int) int {
    if len(arr) == 1 {
        return arr[0]
    }
    
    // 随机选择主元
    pivotIndex := ra.random.Intn(len(arr))
    pivot := arr[pivotIndex]
    
    // 分区
    left := make([]int, 0)
    right := make([]int, 0)
    equal := make([]int, 0)
    
    for _, val := range arr {
        if val < pivot {
            left = append(left, val)
        } else if val > pivot {
            right = append(right, val)
        } else {
            equal = append(equal, val)
        }
    }
    
    // 递归选择
    if k <= len(left) {
        return ra.RandomizedSelection(left, k)
    } else if k <= len(left)+len(equal) {
        return pivot
    } else {
        return ra.RandomizedSelection(right, k-len(left)-len(equal))
    }
}

// RandomizedMatching 随机化匹配算法
func (ra *RandomizedAlgorithm) RandomizedMatching(graph [][]bool) []int {
    n := len(graph)
    matching := make([]int, n)
    for i := range matching {
        matching[i] = -1
    }
    
    // 随机化顶点顺序
    vertices := make([]int, n)
    for i := range vertices {
        vertices[i] = i
    }
    
    // Fisher-Yates洗牌
    for i := n - 1; i > 0; i-- {
        j := ra.random.Intn(i + 1)
        vertices[i], vertices[j] = vertices[j], vertices[i]
    }
    
    // 贪心匹配
    for _, v := range vertices {
        if matching[v] == -1 {
            for u := 0; u < n; u++ {
                if graph[v][u] && matching[u] == -1 {
                    matching[v] = u
                    matching[u] = v
                    break
                }
            }
        }
    }
    
    return matching
}
```

### 5.3 机器学习

```go
// BayesianClassifier 贝叶斯分类器
type BayesianClassifier struct {
    classes     map[string]float64
    features    map[string]map[string]map[string]float64
    smoothing   float64
}

// NewBayesianClassifier 创建贝叶斯分类器
func NewBayesianClassifier(smoothing float64) *BayesianClassifier {
    return &BayesianClassifier{
        classes:   make(map[string]float64),
        features:  make(map[string]map[string]map[string]float64),
        smoothing: smoothing,
    }
}

// Train 训练分类器
func (bc *BayesianClassifier) Train(data []TrainingExample) {
    // 计算类别先验概率
    classCounts := make(map[string]int)
    total := len(data)
    
    for _, example := range data {
        classCounts[example.Class]++
    }
    
    for class, count := range classCounts {
        bc.classes[class] = float64(count) / float64(total)
    }
    
    // 计算特征条件概率
    for _, example := range data {
        if bc.features[example.Class] == nil {
            bc.features[example.Class] = make(map[string]map[string]float64)
        }
        
        for feature, value := range example.Features {
            if bc.features[example.Class][feature] == nil {
                bc.features[example.Class][feature] = make(map[string]float64)
            }
            
            bc.features[example.Class][feature][value]++
        }
    }
    
    // 应用拉普拉斯平滑
    for class := range bc.features {
        for feature := range bc.features[class] {
            totalCount := 0
            for _, count := range bc.features[class][feature] {
                totalCount += int(count)
            }
            
            for value := range bc.features[class][feature] {
                bc.features[class][feature][value] = (bc.features[class][feature][value] + bc.smoothing) / float64(totalCount+int(bc.smoothing)*len(bc.features[class][feature]))
            }
        }
    }
}

// Predict 预测类别
func (bc *BayesianClassifier) Predict(features map[string]string) string {
    bestClass := ""
    bestScore := math.Inf(-1)
    
    for class, prior := range bc.classes {
        score := math.Log(prior)
        
        for feature, value := range features {
            if prob, exists := bc.features[class][feature][value]; exists {
                score += math.Log(prob)
            } else {
                // 使用平滑概率
                score += math.Log(bc.smoothing / float64(len(bc.features[class][feature])))
            }
        }
        
        if score > bestScore {
            bestScore = score
            bestClass = class
        }
    }
    
    return bestClass
}

// TrainingExample 训练样本
type TrainingExample struct {
    Features map[string]string
    Class    string
}

// GaussianMixtureModel 高斯混合模型
type GaussianMixtureModel struct {
    Components []GaussianComponent
    Weights    []float64
}

// GaussianComponent 高斯分量
type GaussianComponent struct {
    Mean   []float64
    Cov    [][]float64
}

// NewGaussianMixtureModel 创建高斯混合模型
func NewGaussianMixtureModel(nComponents int, nFeatures int) *GaussianMixtureModel {
    components := make([]GaussianComponent, nComponents)
    weights := make([]float64, nComponents)
    
    for i := range components {
        components[i] = GaussianComponent{
            Mean: make([]float64, nFeatures),
            Cov:  make([][]float64, nFeatures),
        }
        
        for j := range components[i].Cov {
            components[i].Cov[j] = make([]float64, nFeatures)
        }
        
        weights[i] = 1.0 / float64(nComponents)
    }
    
    return &GaussianMixtureModel{
        Components: components,
        Weights:    weights,
    }
}

// Fit 拟合模型
func (gmm *GaussianMixtureModel) Fit(data [][]float64, maxIterations int) {
    nSamples := len(data)
    nComponents := len(gmm.Components)
    nFeatures := len(data[0])
    
    // 初始化
    responsibilities := make([][]float64, nSamples)
    for i := range responsibilities {
        responsibilities[i] = make([]float64, nComponents)
    }
    
    for iteration := 0; iteration < maxIterations; iteration++ {
        // E步：计算后验概率
        for i, sample := range data {
            total := 0.0
            for j, component := range gmm.Components {
                prob := gmm.gaussianPDF(sample, component.Mean, component.Cov)
                responsibilities[i][j] = gmm.Weights[j] * prob
                total += responsibilities[i][j]
            }
            
            // 归一化
            for j := range responsibilities[i] {
                responsibilities[i][j] /= total
            }
        }
        
        // M步：更新参数
        for j := range gmm.Components {
            // 更新权重
            sum := 0.0
            for i := range data {
                sum += responsibilities[i][j]
            }
            gmm.Weights[j] = sum / float64(nSamples)
            
            // 更新均值
            for k := range gmm.Components[j].Mean {
                gmm.Components[j].Mean[k] = 0
                for i := range data {
                    gmm.Components[j].Mean[k] += responsibilities[i][j] * data[i][k]
                }
                gmm.Components[j].Mean[k] /= sum
            }
            
            // 更新协方差
            for k := range gmm.Components[j].Cov {
                for l := range gmm.Components[j].Cov[k] {
                    gmm.Components[j].Cov[k][l] = 0
                    for i := range data {
                        diff1 := data[i][k] - gmm.Components[j].Mean[k]
                        diff2 := data[i][l] - gmm.Components[j].Mean[l]
                        gmm.Components[j].Cov[k][l] += responsibilities[i][j] * diff1 * diff2
                    }
                    gmm.Components[j].Cov[k][l] /= sum
                }
            }
        }
    }
}

// gaussianPDF 高斯概率密度函数
func (gmm *GaussianMixtureModel) gaussianPDF(x, mean []float64, cov [][]float64) float64 {
    // 简化实现，假设协方差矩阵是对角矩阵
    prob := 1.0
    for i := range x {
        variance := cov[i][i]
        if variance > 0 {
            prob *= math.Exp(-0.5*math.Pow(x[i]-mean[i], 2)/variance) / math.Sqrt(2*math.Pi*variance)
        }
    }
    return prob
}

// Predict 预测
func (gmm *GaussianMixtureModel) Predict(x []float64) int {
    bestComponent := 0
    bestScore := math.Inf(-1)
    
    for j, component := range gmm.Components {
        score := math.Log(gmm.Weights[j]) + math.Log(gmm.gaussianPDF(x, component.Mean, component.Cov))
        if score > bestScore {
            bestScore = score
            bestComponent = j
        }
    }
    
    return bestComponent
}
```

## 总结

概率论为计算机科学提供了强大的数学基础，通过Go语言的实现，我们可以构建高效的概率计算和随机算法库。这些工具在机器学习、蒙特卡洛模拟、随机算法等领域有广泛应用。

**关键特性**:
- 完整的概率空间和随机变量实现
- 支持离散和连续分布
- 蒙特卡洛模拟和随机算法
- 贝叶斯分类和高斯混合模型
- 实际应用场景的示例

**应用领域**:
- 机器学习和人工智能
- 金融风险建模
- 网络性能分析
- 游戏开发
- 科学计算
- 密码学 