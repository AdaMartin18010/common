# 04-概率论 (Probability Theory)

## 概述

概率论是数学的一个分支，研究随机现象的数学规律。它为统计学、机器学习、算法分析等领域提供了理论基础，在计算机科学中广泛应用于随机算法、性能分析、风险评估等。

## 1. 基本概念

### 1.1 样本空间与事件

**定义 1.1** (样本空间)
样本空间 $\Omega$ 是随机试验所有可能结果的集合。

**定义 1.2** (事件)
事件是样本空间的子集，即 $A \subseteq \Omega$。

```go
// 样本空间和事件的基本结构
type SampleSpace[T comparable] struct {
    Elements map[T]bool
}

type Event[T comparable] struct {
    Elements map[T]bool
}

// 创建样本空间
func NewSampleSpace[T comparable](elements ...T) *SampleSpace[T] {
    ss := &SampleSpace[T]{
        Elements: make(map[T]bool),
    }
    for _, element := range elements {
        ss.Elements[element] = true
    }
    return ss
}

// 创建事件
func NewEvent[T comparable](elements ...T) *Event[T] {
    event := &Event[T]{
        Elements: make(map[T]bool),
    }
    for _, element := range elements {
        event.Elements[element] = true
    }
    return event
}
```

### 1.2 概率测度

**定义 1.3** (概率测度)
概率测度 $P$ 是定义在事件集合上的函数，满足：
1. $P(A) \geq 0$ 对所有事件 $A$
2. $P(\Omega) = 1$
3. 对于互斥事件 $A_1, A_2, \ldots$，有 $P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$

```go
// 概率测度
type ProbabilityMeasure[T comparable] struct {
    SampleSpace *SampleSpace[T]
    Probabilities map[T]float64
}

func NewProbabilityMeasure[T comparable](ss *SampleSpace[T]) *ProbabilityMeasure[T] {
    return &ProbabilityMeasure[T]{
        SampleSpace:  ss,
        Probabilities: make(map[T]float64),
    }
}

// 设置基本事件的概率
func (pm *ProbabilityMeasure[T]) SetProbability(element T, prob float64) {
    if pm.SampleSpace.Elements[element] {
        pm.Probabilities[element] = prob
    }
}

// 计算事件的概率
func (pm *ProbabilityMeasure[T]) Probability(event *Event[T]) float64 {
    total := 0.0
    for element := range event.Elements {
        if pm.SampleSpace.Elements[element] {
            total += pm.Probabilities[element]
        }
    }
    return total
}
```

## 2. 随机变量

### 2.1 随机变量定义

**定义 2.1** (随机变量)
随机变量 $X$ 是从样本空间 $\Omega$ 到实数集 $\mathbb{R}$ 的函数：$X: \Omega \rightarrow \mathbb{R}$

```go
// 随机变量
type RandomVariable[T comparable] struct {
    SampleSpace *SampleSpace[T]
    Function    func(T) float64
}

func NewRandomVariable[T comparable](ss *SampleSpace[T], f func(T) float64) *RandomVariable[T] {
    return &RandomVariable[T]{
        SampleSpace: ss,
        Function:    f,
    }
}

// 计算随机变量的值
func (rv *RandomVariable[T]) Value(omega T) float64 {
    return rv.Function(omega)
}

// 随机变量的期望
func (rv *RandomVariable[T]) Expectation(pm *ProbabilityMeasure[T]) float64 {
    expectation := 0.0
    for element := range rv.SampleSpace.Elements {
        expectation += rv.Function(element) * pm.Probabilities[element]
    }
    return expectation
}

// 随机变量的方差
func (rv *RandomVariable[T]) Variance(pm *ProbabilityMeasure[T]) float64 {
    expectation := rv.Expectation(pm)
    variance := 0.0
    
    for element := range rv.SampleSpace.Elements {
        diff := rv.Function(element) - expectation
        variance += diff * diff * pm.Probabilities[element]
    }
    
    return variance
}
```

### 2.2 离散随机变量

```go
// 离散随机变量
type DiscreteRandomVariable struct {
    Values       []float64
    Probabilities []float64
}

func NewDiscreteRandomVariable(values, probabilities []float64) *DiscreteRandomVariable {
    if len(values) != len(probabilities) {
        panic("Values and probabilities must have the same length")
    }
    
    return &DiscreteRandomVariable{
        Values:       values,
        Probabilities: probabilities,
    }
}

// 期望
func (drv *DiscreteRandomVariable) Expectation() float64 {
    expectation := 0.0
    for i, value := range drv.Values {
        expectation += value * drv.Probabilities[i]
    }
    return expectation
}

// 方差
func (drv *DiscreteRandomVariable) Variance() float64 {
    expectation := drv.Expectation()
    variance := 0.0
    
    for i, value := range drv.Values {
        diff := value - expectation
        variance += diff * diff * drv.Probabilities[i]
    }
    
    return variance
}

// 概率质量函数
func (drv *DiscreteRandomVariable) PMF(x float64) float64 {
    for i, value := range drv.Values {
        if value == x {
            return drv.Probabilities[i]
        }
    }
    return 0.0
}
```

## 3. 常见概率分布

### 3.1 伯努利分布

**定义 3.1** (伯努利分布)
伯努利分布是参数为 $p$ 的离散分布，随机变量 $X$ 取值为：
$$P(X = 1) = p, \quad P(X = 0) = 1-p$$

```go
// 伯努利分布
type BernoulliDistribution struct {
    P float64 // 成功概率
}

func NewBernoulliDistribution(p float64) *BernoulliDistribution {
    if p < 0 || p > 1 {
        panic("Probability must be between 0 and 1")
    }
    return &BernoulliDistribution{P: p}
}

// 概率质量函数
func (bd *BernoulliDistribution) PMF(x int) float64 {
    if x == 1 {
        return bd.P
    } else if x == 0 {
        return 1 - bd.P
    }
    return 0.0
}

// 期望
func (bd *BernoulliDistribution) Expectation() float64 {
    return bd.P
}

// 方差
func (bd *BernoulliDistribution) Variance() float64 {
    return bd.P * (1 - bd.P)
}

// 生成随机样本
func (bd *BernoulliDistribution) Sample() int {
    if rand.Float64() < bd.P {
        return 1
    }
    return 0
}
```

### 3.2 二项分布

**定义 3.2** (二项分布)
二项分布 $B(n,p)$ 是 $n$ 次独立伯努利试验中成功次数的分布：
$$P(X = k) = \binom{n}{k} p^k (1-p)^{n-k}$$

```go
// 二项分布
type BinomialDistribution struct {
    N int     // 试验次数
    P float64 // 成功概率
}

func NewBinomialDistribution(n int, p float64) *BinomialDistribution {
    if n < 0 || p < 0 || p > 1 {
        panic("Invalid parameters")
    }
    return &BinomialDistribution{N: n, P: p}
}

// 组合数计算
func combination(n, k int) int {
    if k > n {
        return 0
    }
    if k > n/2 {
        k = n - k
    }
    
    result := 1
    for i := 0; i < k; i++ {
        result = result * (n - i) / (i + 1)
    }
    return result
}

// 概率质量函数
func (bd *BinomialDistribution) PMF(k int) float64 {
    if k < 0 || k > bd.N {
        return 0.0
    }
    
    c := combination(bd.N, k)
    return float64(c) * math.Pow(bd.P, float64(k)) * math.Pow(1-bd.P, float64(bd.N-k))
}

// 期望
func (bd *BinomialDistribution) Expectation() float64 {
    return float64(bd.N) * bd.P
}

// 方差
func (bd *BinomialDistribution) Variance() float64 {
    return float64(bd.N) * bd.P * (1 - bd.P)
}

// 生成随机样本
func (bd *BinomialDistribution) Sample() int {
    successes := 0
    for i := 0; i < bd.N; i++ {
        if rand.Float64() < bd.P {
            successes++
        }
    }
    return successes
}
```

### 3.3 泊松分布

**定义 3.3** (泊松分布)
泊松分布 $P(\lambda)$ 是参数为 $\lambda$ 的离散分布：
$$P(X = k) = \frac{\lambda^k e^{-\lambda}}{k!}$$

```go
// 泊松分布
type PoissonDistribution struct {
    Lambda float64 // 参数λ
}

func NewPoissonDistribution(lambda float64) *PoissonDistribution {
    if lambda < 0 {
        panic("Lambda must be non-negative")
    }
    return &PoissonDistribution{Lambda: lambda}
}

// 阶乘计算
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}

// 概率质量函数
func (pd *PoissonDistribution) PMF(k int) float64 {
    if k < 0 {
        return 0.0
    }
    
    numerator := math.Pow(pd.Lambda, float64(k)) * math.Exp(-pd.Lambda)
    denominator := float64(factorial(k))
    
    return numerator / denominator
}

// 期望
func (pd *PoissonDistribution) Expectation() float64 {
    return pd.Lambda
}

// 方差
func (pd *PoissonDistribution) Variance() float64 {
    return pd.Lambda
}

// 生成随机样本（使用反变换法）
func (pd *PoissonDistribution) Sample() int {
    u := rand.Float64()
    k := 0
    p := math.Exp(-pd.Lambda)
    f := p
    
    for u > f {
        k++
        p *= pd.Lambda / float64(k)
        f += p
    }
    
    return k
}
```

### 3.4 正态分布

**定义 3.4** (正态分布)
正态分布 $N(\mu, \sigma^2)$ 是参数为 $\mu$ 和 $\sigma^2$ 的连续分布：
$$f(x) = \frac{1}{\sqrt{2\pi\sigma^2}} e^{-\frac{(x-\mu)^2}{2\sigma^2}}$$

```go
// 正态分布
type NormalDistribution struct {
    Mu    float64 // 均值μ
    Sigma float64 // 标准差σ
}

func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    if sigma <= 0 {
        panic("Sigma must be positive")
    }
    return &NormalDistribution{Mu: mu, Sigma: sigma}
}

// 概率密度函数
func (nd *NormalDistribution) PDF(x float64) float64 {
    exponent := -0.5 * math.Pow((x-nd.Mu)/nd.Sigma, 2)
    return (1.0 / (nd.Sigma * math.Sqrt(2*math.Pi))) * math.Exp(exponent)
}

// 累积分布函数（近似）
func (nd *NormalDistribution) CDF(x float64) float64 {
    z := (x - nd.Mu) / nd.Sigma
    return 0.5 * (1 + erf(z/math.Sqrt(2)))
}

// 误差函数近似
func erf(x float64) float64 {
    // 使用近似公式
    a1 := 0.254829592
    a2 := -0.284496736
    a3 := 1.421413741
    a4 := -1.453152027
    a5 := 1.061405429
    p := 0.3275911
    
    sign := 1.0
    if x < 0 {
        sign = -1.0
    }
    x = math.Abs(x)
    
    t := 1.0 / (1.0 + p*x)
    y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)
    
    return sign * y
}

// 期望
func (nd *NormalDistribution) Expectation() float64 {
    return nd.Mu
}

// 方差
func (nd *NormalDistribution) Variance() float64 {
    return nd.Sigma * nd.Sigma
}

// 生成随机样本（Box-Muller变换）
func (nd *NormalDistribution) Sample() float64 {
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    // Box-Muller变换
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    
    return nd.Mu + nd.Sigma*z0
}
```

## 4. 条件概率与独立性

### 4.1 条件概率

**定义 4.1** (条件概率)
在事件 $B$ 发生的条件下，事件 $A$ 发生的条件概率为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

```go
// 条件概率计算
func ConditionalProbability[T comparable](A, B *Event[T], pm *ProbabilityMeasure[T]) float64 {
    // 计算A ∩ B
    intersection := NewEvent[T]()
    for element := range A.Elements {
        if B.Elements[element] {
            intersection.Elements[element] = true
        }
    }
    
    pAB := pm.Probability(intersection)
    pB := pm.Probability(B)
    
    if pB == 0 {
        return 0
    }
    
    return pAB / pB
}

// 贝叶斯定理
func BayesTheorem[T comparable](A, B *Event[T], pm *ProbabilityMeasure[T]) float64 {
    pA := pm.Probability(A)
    pB := pm.Probability(B)
    pBA := ConditionalProbability(B, A, pm)
    
    if pB == 0 {
        return 0
    }
    
    return (pBA * pA) / pB
}
```

### 4.2 独立性

**定义 4.2** (独立性)
事件 $A$ 和 $B$ 独立当且仅当：
$$P(A \cap B) = P(A) \cdot P(B)$$

```go
// 检查事件独立性
func AreIndependent[T comparable](A, B *Event[T], pm *ProbabilityMeasure[T]) bool {
    pA := pm.Probability(A)
    pB := pm.Probability(B)
    
    // 计算A ∩ B
    intersection := NewEvent[T]()
    for element := range A.Elements {
        if B.Elements[element] {
            intersection.Elements[element] = true
        }
    }
    
    pAB := pm.Probability(intersection)
    
    // 检查是否独立
    return math.Abs(pAB-pA*pB) < 1e-10
}
```

## 5. 大数定律与中心极限定理

### 5.1 大数定律

**定理 5.1** (大数定律)
设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，则：
$$\lim_{n \to \infty} \frac{1}{n} \sum_{i=1}^{n} X_i = \mu \quad \text{a.s.}$$

```go
// 大数定律模拟
func LawOfLargeNumbers(distribution interface{ Sample() float64 }, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        sum += distribution.Sample()
    }
    return sum / float64(n)
}

// 验证大数定律
func VerifyLawOfLargeNumbers() {
    normal := NewNormalDistribution(5.0, 2.0)
    expected := normal.Expectation()
    
    fmt.Printf("Expected value: %.2f\n", expected)
    
    for n := 10; n <= 10000; n *= 10 {
        sampleMean := LawOfLargeNumbers(normal, n)
        fmt.Printf("n=%d, sample mean=%.2f, error=%.4f\n", 
                   n, sampleMean, math.Abs(sampleMean-expected))
    }
}
```

### 5.2 中心极限定理

**定理 5.2** (中心极限定理)
设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，方差为 $\sigma^2$，则：
$$\frac{\sum_{i=1}^{n} X_i - n\mu}{\sqrt{n}\sigma} \xrightarrow{d} N(0,1)$$

```go
// 中心极限定理模拟
func CentralLimitTheorem(distribution interface{ 
    Sample() float64 
    Expectation() float64 
    Variance() float64 
}, n, samples int) []float64 {
    mu := distribution.Expectation()
    sigma := math.Sqrt(distribution.Variance())
    
    standardizedMeans := make([]float64, samples)
    
    for i := 0; i < samples; i++ {
        sum := 0.0
        for j := 0; j < n; j++ {
            sum += distribution.Sample()
        }
        
        sampleMean := sum / float64(n)
        standardizedMeans[i] = (sampleMean - mu) / (sigma / math.Sqrt(float64(n)))
    }
    
    return standardizedMeans
}
```

## 6. 随机算法应用

### 6.1 随机快速排序

```go
// 随机快速排序
func RandomizedQuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    // 随机选择主元
    pivotIndex := rand.Intn(len(arr))
    pivot := arr[pivotIndex]
    
    var left, right, equal []int
    
    for _, x := range arr {
        if x < pivot {
            left = append(left, x)
        } else if x > pivot {
            right = append(right, x)
        } else {
            equal = append(equal, x)
        }
    }
    
    // 递归排序
    left = RandomizedQuickSort(left)
    right = RandomizedQuickSort(right)
    
    // 合并结果
    result := append(left, equal...)
    result = append(result, right...)
    
    return result
}

// 分析随机快速排序的期望时间复杂度
func AnalyzeRandomizedQuickSort(n int) float64 {
    // 期望比较次数约为 2n ln n
    return 2 * float64(n) * math.Log(float64(n))
}
```

### 6.2 随机化算法性能分析

```go
// 性能分析工具
type PerformanceAnalyzer struct {
    Trials int
}

func NewPerformanceAnalyzer(trials int) *PerformanceAnalyzer {
    return &PerformanceAnalyzer{Trials: trials}
}

// 分析算法性能
func (pa *PerformanceAnalyzer) AnalyzePerformance(algorithm func([]int) []int, sizes []int) map[int]float64 {
    results := make(map[int]float64)
    
    for _, size := range sizes {
        totalTime := 0.0
        
        for trial := 0; trial < pa.Trials; trial++ {
            // 生成随机数组
            arr := make([]int, size)
            for i := range arr {
                arr[i] = rand.Intn(1000)
            }
            
            // 测量执行时间
            start := time.Now()
            algorithm(arr)
            duration := time.Since(start)
            
            totalTime += float64(duration.Nanoseconds())
        }
        
        results[size] = totalTime / float64(pa.Trials)
    }
    
    return results
}
```

## 7. 概率论在机器学习中的应用

### 7.1 朴素贝叶斯分类器

```go
// 朴素贝叶斯分类器
type NaiveBayesClassifier struct {
    classProbabilities map[string]float64
    featureProbabilities map[string]map[string]float64
    classes []string
    features []string
}

func NewNaiveBayesClassifier() *NaiveBayesClassifier {
    return &NaiveBayesClassifier{
        classProbabilities: make(map[string]float64),
        featureProbabilities: make(map[string]map[string]float64),
    }
}

// 训练分类器
func (nb *NaiveBayesClassifier) Train(data [][]string, labels []string) {
    // 计算类别概率
    classCounts := make(map[string]int)
    for _, label := range labels {
        classCounts[label]++
    }
    
    total := len(labels)
    for class, count := range classCounts {
        nb.classProbabilities[class] = float64(count) / float64(total)
    }
    
    // 计算特征条件概率
    for featureIndex := range data[0] {
        featureName := fmt.Sprintf("feature_%d", featureIndex)
        nb.featureProbabilities[featureName] = make(map[string]float64)
        
        for class := range classCounts {
            // 计算在给定类别下该特征的概率
            featureCount := 0
            classDataCount := 0
            
            for i, label := range labels {
                if label == class {
                    classDataCount++
                    if data[i][featureIndex] == "1" {
                        featureCount++
                    }
                }
            }
            
            if classDataCount > 0 {
                nb.featureProbabilities[featureName][class] = float64(featureCount) / float64(classDataCount)
            }
        }
    }
}

// 预测
func (nb *NaiveBayesClassifier) Predict(features []string) string {
    bestClass := ""
    bestProb := -1.0
    
    for class, classProb := range nb.classProbabilities {
        prob := math.Log(classProb)
        
        for featureIndex, featureValue := range features {
            featureName := fmt.Sprintf("feature_%d", featureIndex)
            if featureProbs, exists := nb.featureProbabilities[featureName]; exists {
                if featureProb, exists := featureProbs[class]; exists {
                    if featureValue == "1" {
                        prob += math.Log(featureProb)
                    } else {
                        prob += math.Log(1 - featureProb)
                    }
                }
            }
        }
        
        if prob > bestProb {
            bestProb = prob
            bestClass = class
        }
    }
    
    return bestClass
}
```

## 总结

概率论为计算机科学提供了强大的数学工具，通过Go语言的实现，我们可以：

1. **理论实现**: 实现各种概率分布和统计方法
2. **算法分析**: 分析随机算法的期望性能
3. **机器学习**: 构建基于概率的机器学习模型
4. **性能评估**: 通过概率方法评估系统性能

概率论的理论基础与Go语言的实践相结合，为构建可靠的随机算法和机器学习系统提供了坚实的基础。
