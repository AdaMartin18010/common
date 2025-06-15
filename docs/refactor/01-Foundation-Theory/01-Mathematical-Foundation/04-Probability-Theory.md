# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 样本空间与事件](#11-样本空间与事件)
    - [1.2 概率公理](#12-概率公理)
    - [1.3 条件概率](#13-条件概率)
  - [2. 随机变量](#2-随机变量)
    - [2.1 离散随机变量](#21-离散随机变量)
    - [2.2 连续随机变量](#22-连续随机变量)
    - [2.3 期望与方差](#23-期望与方差)
  - [3. 概率分布](#3-概率分布)
    - [3.1 离散分布](#31-离散分布)
    - [3.2 连续分布](#32-连续分布)
    - [3.3 多维分布](#33-多维分布)
  - [4. 大数定律与中心极限定理](#4-大数定律与中心极限定理)
    - [4.1 大数定律](#41-大数定律)
    - [4.2 中心极限定理](#42-中心极限定理)
  - [5. 应用领域](#5-应用领域)
    - [5.1 机器学习](#51-机器学习)
    - [5.2 风险评估](#52-风险评估)
    - [5.3 性能分析](#53-性能分析)
  - [6. 总结](#6-总结)

## 1. 基础概念

### 1.1 样本空间与事件

**定义 1.1** (样本空间): 样本空间 $\Omega$ 是随机试验所有可能结果的集合。

**定义 1.2** (事件): 事件 $A$ 是样本空间 $\Omega$ 的子集，即 $A \subseteq \Omega$。

**定义 1.3** (事件代数): 事件代数 $\mathcal{F}$ 是样本空间子集的集合，满足：
1. $\Omega \in \mathcal{F}$
2. 如果 $A \in \mathcal{F}$，则 $A^c \in \mathcal{F}$
3. 如果 $A_1, A_2, \ldots \in \mathcal{F}$，则 $\bigcup_{i=1}^{\infty} A_i \in \mathcal{F}$

```go
// 样本空间和事件的基本定义
type SampleSpace struct {
    elements []interface{}
}

type Event struct {
    elements map[interface{}]bool
    space    *SampleSpace
}

// 创建样本空间
func NewSampleSpace(elements ...interface{}) *SampleSpace {
    return &SampleSpace{elements: elements}
}

// 创建事件
func NewEvent(space *SampleSpace, elements ...interface{}) *Event {
    event := &Event{
        elements: make(map[interface{}]bool),
        space:    space,
    }
    
    for _, element := range elements {
        event.elements[element] = true
    }
    
    return event
}

// 事件的基本操作
func (e *Event) Contains(element interface{}) bool {
    return e.elements[element]
}

func (e *Event) Complement() *Event {
    complement := &Event{
        elements: make(map[interface{}]bool),
        space:    e.space,
    }
    
    for _, element := range e.space.elements {
        if !e.elements[element] {
            complement.elements[element] = true
        }
    }
    
    return complement
}

func (e *Event) Union(other *Event) *Event {
    union := &Event{
        elements: make(map[interface{}]bool),
        space:    e.space,
    }
    
    // 复制当前事件的元素
    for element := range e.elements {
        union.elements[element] = true
    }
    
    // 添加另一个事件的元素
    for element := range other.elements {
        union.elements[element] = true
    }
    
    return union
}

func (e *Event) Intersection(other *Event) *Event {
    intersection := &Event{
        elements: make(map[interface{}]bool),
        space:    e.space,
    }
    
    for element := range e.elements {
        if other.elements[element] {
            intersection.elements[element] = true
        }
    }
    
    return intersection
}

func (e *Event) Size() int {
    return len(e.elements)
}
```

### 1.2 概率公理

**公理 1.1** (Kolmogorov公理): 概率函数 $P: \mathcal{F} \rightarrow [0,1]$ 满足：
1. 非负性：$P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
2. 规范性：$P(\Omega) = 1$
3. 可列可加性：对于互斥事件 $A_1, A_2, \ldots$，
   $$P\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} P(A_i)$$

```go
// 概率空间
type ProbabilitySpace struct {
    sampleSpace *SampleSpace
    events      map[string]*Event
    probabilities map[string]float64
}

func NewProbabilitySpace(space *SampleSpace) *ProbabilitySpace {
    return &ProbabilitySpace{
        sampleSpace:   space,
        events:        make(map[string]*Event),
        probabilities: make(map[string]float64),
    }
}

// 设置事件概率
func (ps *ProbabilitySpace) SetEventProbability(eventName string, event *Event, probability float64) {
    if probability < 0 || probability > 1 {
        panic("概率必须在[0,1]范围内")
    }
    
    ps.events[eventName] = event
    ps.probabilities[eventName] = probability
}

// 获取事件概率
func (ps *ProbabilitySpace) GetProbability(eventName string) float64 {
    if prob, exists := ps.probabilities[eventName]; exists {
        return prob
    }
    return 0.0
}

// 计算事件概率 (基于元素计数)
func (ps *ProbabilitySpace) CalculateProbability(event *Event) float64 {
    if len(ps.sampleSpace.elements) == 0 {
        return 0.0
    }
    
    return float64(event.Size()) / float64(len(ps.sampleSpace.elements))
}

// 验证概率公理
func (ps *ProbabilitySpace) ValidateAxioms() bool {
    // 检查规范性
    if ps.GetProbability("sample_space") != 1.0 {
        return false
    }
    
    // 检查非负性
    for _, prob := range ps.probabilities {
        if prob < 0 {
            return false
        }
    }
    
    return true
}
```

### 1.3 条件概率

**定义 1.4** (条件概率): 对于事件 $A$ 和 $B$，在 $B$ 发生的条件下 $A$ 发生的概率为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 1.1** (乘法公式): $P(A \cap B) = P(A|B) \cdot P(B) = P(B|A) \cdot P(A)$

**定理 1.2** (全概率公式): 如果 $B_1, B_2, \ldots, B_n$ 是样本空间的一个划分，则：
$$P(A) = \sum_{i=1}^{n} P(A|B_i) \cdot P(B_i)$$

```go
// 条件概率计算
func (ps *ProbabilitySpace) ConditionalProbability(eventA, eventB *Event) float64 {
    intersection := eventA.Intersection(eventB)
    probB := ps.CalculateProbability(eventB)
    
    if probB == 0 {
        return 0.0
    }
    
    probIntersection := ps.CalculateProbability(intersection)
    return probIntersection / probB
}

// 贝叶斯定理
func (ps *ProbabilitySpace) BayesTheorem(eventA, eventB *Event) float64 {
    probA := ps.CalculateProbability(eventA)
    probB := ps.CalculateProbability(eventB)
    probBGivenA := ps.ConditionalProbability(eventB, eventA)
    
    if probB == 0 {
        return 0.0
    }
    
    return (probBGivenA * probA) / probB
}

// 独立性检验
func (ps *ProbabilitySpace) AreIndependent(eventA, eventB *Event) bool {
    probA := ps.CalculateProbability(eventA)
    probB := ps.CalculateProbability(eventB)
    probAAndB := ps.CalculateProbability(eventA.Intersection(eventB))
    
    return math.Abs(probAAndB - probA*probB) < 1e-10
}
```

## 2. 随机变量

### 2.1 离散随机变量

**定义 2.1** (离散随机变量): 离散随机变量 $X$ 是定义在样本空间上的函数，其取值是有限或可数无限个实数。

**定义 2.2** (概率质量函数): 离散随机变量 $X$ 的概率质量函数 $p_X(x)$ 定义为：
$$p_X(x) = P(X = x)$$

```go
// 离散随机变量
type DiscreteRandomVariable struct {
    values    []float64
    pmf       map[float64]float64
    cdf       map[float64]float64
}

func NewDiscreteRandomVariable(values []float64, probabilities []float64) *DiscreteRandomVariable {
    if len(values) != len(probabilities) {
        panic("值和概率数组长度必须相同")
    }
    
    rv := &DiscreteRandomVariable{
        values: values,
        pmf:    make(map[float64]float64),
        cdf:    make(map[float64]float64),
    }
    
    // 设置概率质量函数
    for i, value := range values {
        rv.pmf[value] = probabilities[i]
    }
    
    // 计算累积分布函数
    rv.calculateCDF()
    
    return rv
}

func (rv *DiscreteRandomVariable) calculateCDF() {
    cumulative := 0.0
    for _, value := range rv.values {
        cumulative += rv.pmf[value]
        rv.cdf[value] = cumulative
    }
}

func (rv *DiscreteRandomVariable) PMF(x float64) float64 {
    return rv.pmf[x]
}

func (rv *DiscreteRandomVariable) CDF(x float64) float64 {
    return rv.cdf[x]
}

// 生成随机样本
func (rv *DiscreteRandomVariable) Sample() float64 {
    r := rand.Float64()
    
    for _, value := range rv.values {
        if r <= rv.cdf[value] {
            return value
        }
    }
    
    return rv.values[len(rv.values)-1]
}

// 生成多个样本
func (rv *DiscreteRandomVariable) SampleN(n int) []float64 {
    samples := make([]float64, n)
    for i := 0; i < n; i++ {
        samples[i] = rv.Sample()
    }
    return samples
}
```

### 2.2 连续随机变量

**定义 2.3** (连续随机变量): 连续随机变量 $X$ 是定义在样本空间上的函数，其取值是连续的实数。

**定义 2.4** (概率密度函数): 连续随机变量 $X$ 的概率密度函数 $f_X(x)$ 满足：
$$P(a \leq X \leq b) = \int_a^b f_X(x) dx$$

```go
// 连续随机变量接口
type ContinuousRandomVariable interface {
    PDF(x float64) float64
    CDF(x float64) float64
    Sample() float64
    SampleN(n int) []float64
}

// 均匀分布
type UniformDistribution struct {
    a, b float64 // 区间[a,b]
}

func NewUniformDistribution(a, b float64) *UniformDistribution {
    if a >= b {
        panic("a必须小于b")
    }
    return &UniformDistribution{a: a, b: b}
}

func (u *UniformDistribution) PDF(x float64) float64 {
    if x >= u.a && x <= u.b {
        return 1.0 / (u.b - u.a)
    }
    return 0.0
}

func (u *UniformDistribution) CDF(x float64) float64 {
    if x < u.a {
        return 0.0
    } else if x >= u.b {
        return 1.0
    } else {
        return (x - u.a) / (u.b - u.a)
    }
}

func (u *UniformDistribution) Sample() float64 {
    return u.a + rand.Float64()*(u.b-u.a)
}

func (u *UniformDistribution) SampleN(n int) []float64 {
    samples := make([]float64, n)
    for i := 0; i < n; i++ {
        samples[i] = u.Sample()
    }
    return samples
}

// 正态分布
type NormalDistribution struct {
    mu    float64 // 均值
    sigma float64 // 标准差
}

func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    if sigma <= 0 {
        panic("标准差必须为正数")
    }
    return &NormalDistribution{mu: mu, sigma: sigma}
}

func (n *NormalDistribution) PDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return (1.0 / (n.sigma * math.Sqrt(2*math.Pi))) * math.Exp(-0.5*z*z)
}

func (n *NormalDistribution) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + math.Erf(z/math.Sqrt(2)))
}

func (n *NormalDistribution) Sample() float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.mu + n.sigma*z0
}

func (n *NormalDistribution) SampleN(n_samples int) []float64 {
    samples := make([]float64, n_samples)
    for i := 0; i < n_samples; i++ {
        samples[i] = n.Sample()
    }
    return samples
}
```

### 2.3 期望与方差

**定义 2.5** (期望): 离散随机变量 $X$ 的期望定义为：
$$E[X] = \sum_x x \cdot p_X(x)$$

连续随机变量 $X$ 的期望定义为：
$$E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$$

**定义 2.6** (方差): 随机变量 $X$ 的方差定义为：
$$\text{Var}(X) = E[(X - E[X])^2] = E[X^2] - (E[X])^2$$

```go
// 期望和方差计算
func (rv *DiscreteRandomVariable) Expectation() float64 {
    expectation := 0.0
    for value, prob := range rv.pmf {
        expectation += value * prob
    }
    return expectation
}

func (rv *DiscreteRandomVariable) Variance() float64 {
    expectation := rv.Expectation()
    variance := 0.0
    
    for value, prob := range rv.pmf {
        variance += prob * (value - expectation) * (value - expectation)
    }
    
    return variance
}

func (rv *DiscreteRandomVariable) StandardDeviation() float64 {
    return math.Sqrt(rv.Variance())
}

// 连续随机变量的数值积分
func (rv *ContinuousRandomVariable) Expectation(a, b float64, n int) float64 {
    dx := (b - a) / float64(n)
    expectation := 0.0
    
    for i := 0; i < n; i++ {
        x := a + float64(i)*dx
        expectation += x * rv.PDF(x) * dx
    }
    
    return expectation
}

func (rv *ContinuousRandomVariable) Variance(a, b float64, n int) float64 {
    expectation := rv.Expectation(a, b, n)
    dx := (b - a) / float64(n)
    variance := 0.0
    
    for i := 0; i < n; i++ {
        x := a + float64(i)*dx
        variance += rv.PDF(x) * (x-expectation) * (x-expectation) * dx
    }
    
    return variance
}
```

## 3. 概率分布

### 3.1 离散分布

```go
// 二项分布
type BinomialDistribution struct {
    n int     // 试验次数
    p float64 // 成功概率
}

func NewBinomialDistribution(n int, p float64) *BinomialDistribution {
    if p < 0 || p > 1 {
        panic("概率必须在[0,1]范围内")
    }
    return &BinomialDistribution{n: n, p: p}
}

func (b *BinomialDistribution) PMF(k int) float64 {
    if k < 0 || k > b.n {
        return 0.0
    }
    
    // 计算组合数 C(n,k)
    combination := b.combination(b.n, k)
    return combination * math.Pow(b.p, float64(k)) * math.Pow(1-b.p, float64(b.n-k))
}

func (b *BinomialDistribution) combination(n, k int) float64 {
    if k > n {
        return 0
    }
    if k == 0 || k == n {
        return 1
    }
    
    result := 1.0
    for i := 1; i <= k; i++ {
        result *= float64(n-i+1) / float64(i)
    }
    return result
}

func (b *BinomialDistribution) Expectation() float64 {
    return float64(b.n) * b.p
}

func (b *BinomialDistribution) Variance() float64 {
    return float64(b.n) * b.p * (1 - b.p)
}

// 泊松分布
type PoissonDistribution struct {
    lambda float64 // 参数λ
}

func NewPoissonDistribution(lambda float64) *PoissonDistribution {
    if lambda <= 0 {
        panic("λ必须为正数")
    }
    return &PoissonDistribution{lambda: lambda}
}

func (p *PoissonDistribution) PMF(k int) float64 {
    if k < 0 {
        return 0.0
    }
    
    return (math.Pow(p.lambda, float64(k)) * math.Exp(-p.lambda)) / float64(factorial(k))
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

func (p *PoissonDistribution) Expectation() float64 {
    return p.lambda
}

func (p *PoissonDistribution) Variance() float64 {
    return p.lambda
}
```

### 3.2 连续分布

```go
// 指数分布
type ExponentialDistribution struct {
    lambda float64 // 参数λ
}

func NewExponentialDistribution(lambda float64) *ExponentialDistribution {
    if lambda <= 0 {
        panic("λ必须为正数")
    }
    return &ExponentialDistribution{lambda: lambda}
}

func (e *ExponentialDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0.0
    }
    return e.lambda * math.Exp(-e.lambda*x)
}

func (e *ExponentialDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0.0
    }
    return 1.0 - math.Exp(-e.lambda*x)
}

func (e *ExponentialDistribution) Sample() float64 {
    return -math.Log(1-rand.Float64()) / e.lambda
}

func (e *ExponentialDistribution) Expectation() float64 {
    return 1.0 / e.lambda
}

func (e *ExponentialDistribution) Variance() float64 {
    return 1.0 / (e.lambda * e.lambda)
}

// 伽马分布
type GammaDistribution struct {
    alpha float64 // 形状参数
    beta  float64 // 尺度参数
}

func NewGammaDistribution(alpha, beta float64) *GammaDistribution {
    if alpha <= 0 || beta <= 0 {
        panic("参数必须为正数")
    }
    return &GammaDistribution{alpha: alpha, beta: beta}
}

func (g *GammaDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0.0
    }
    
    // 简化实现，使用数值方法
    return (math.Pow(x, g.alpha-1) * math.Exp(-x/g.beta)) / 
           (math.Pow(g.beta, g.alpha) * gamma(g.alpha))
}

func gamma(x float64) float64 {
    // 简化实现，使用Stirling公式
    if x <= 0 {
        return math.Inf(1)
    }
    return math.Sqrt(2*math.Pi/x) * math.Pow(x/math.E, x)
}

func (g *GammaDistribution) Expectation() float64 {
    return g.alpha * g.beta
}

func (g *GammaDistribution) Variance() float64 {
    return g.alpha * g.beta * g.beta
}
```

### 3.3 多维分布

```go
// 二维正态分布
type BivariateNormalDistribution struct {
    mu1, mu2    float64 // 均值
    sigma1, sigma2 float64 // 标准差
    rho         float64 // 相关系数
}

func NewBivariateNormalDistribution(mu1, mu2, sigma1, sigma2, rho float64) *BivariateNormalDistribution {
    if sigma1 <= 0 || sigma2 <= 0 || rho < -1 || rho > 1 {
        panic("参数无效")
    }
    return &BivariateNormalDistribution{
        mu1: mu1, mu2: mu2,
        sigma1: sigma1, sigma2: sigma2,
        rho: rho,
    }
}

func (b *BivariateNormalDistribution) PDF(x, y float64) float64 {
    z1 := (x - b.mu1) / b.sigma1
    z2 := (y - b.mu2) / b.sigma2
    
    exponent := -(z1*z1 - 2*b.rho*z1*z2 + z2*z2) / (2 * (1 - b.rho*b.rho))
    denominator := 2 * math.Pi * b.sigma1 * b.sigma2 * math.Sqrt(1-b.rho*b.rho)
    
    return math.Exp(exponent) / denominator
}

func (b *BivariateNormalDistribution) Sample() (float64, float64) {
    // Box-Muller变换生成两个独立的标准正态随机变量
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z1 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    z2 := math.Sqrt(-2*math.Log(u1)) * math.Sin(2*math.Pi*u2)
    
    // 转换为相关的正态随机变量
    x := b.mu1 + b.sigma1*z1
    y := b.mu2 + b.sigma2*(b.rho*z1 + math.Sqrt(1-b.rho*b.rho)*z2)
    
    return x, y
}
```

## 4. 大数定律与中心极限定理

### 4.1 大数定律

**定理 4.1** (弱大数定律): 设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，则：
$$\frac{1}{n} \sum_{i=1}^{n} X_i \xrightarrow{P} \mu$$

**定理 4.2** (强大数定律): 在相同条件下：
$$\frac{1}{n} \sum_{i=1}^{n} X_i \xrightarrow{a.s.} \mu$$

```go
// 大数定律验证
func LawOfLargeNumbers(distribution ContinuousRandomVariable, n int) float64 {
    samples := distribution.SampleN(n)
    sum := 0.0
    
    for _, sample := range samples {
        sum += sample
    }
    
    return sum / float64(n)
}

// 多次实验验证大数定律
func VerifyLawOfLargeNumbers(distribution ContinuousRandomVariable, n, experiments int) []float64 {
    results := make([]float64, experiments)
    
    for i := 0; i < experiments; i++ {
        results[i] = LawOfLargeNumbers(distribution, n)
    }
    
    return results
}

// 计算样本均值的收敛性
func ConvergenceAnalysis(distribution ContinuousRandomVariable, maxN int) []float64 {
    convergence := make([]float64, maxN)
    
    for n := 1; n <= maxN; n++ {
        convergence[n-1] = LawOfLargeNumbers(distribution, n)
    }
    
    return convergence
}
```

### 4.2 中心极限定理

**定理 4.3** (中心极限定理): 设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，方差为 $\sigma^2$，则：
$$\frac{\sum_{i=1}^{n} X_i - n\mu}{\sigma\sqrt{n}} \xrightarrow{d} N(0,1)$$

```go
// 中心极限定理验证
func CentralLimitTheorem(distribution ContinuousRandomVariable, n, samples int) []float64 {
    standardizedSums := make([]float64, samples)
    
    for i := 0; i < samples; i++ {
        // 生成n个样本
        sampleSet := distribution.SampleN(n)
        
        // 计算样本和
        sum := 0.0
        for _, sample := range sampleSet {
            sum += sample
        }
        
        // 标准化
        mean := distribution.Expectation(-10, 10, 1000) // 近似期望
        variance := distribution.Variance(-10, 10, 1000) // 近似方差
        stdDev := math.Sqrt(variance)
        
        standardizedSums[i] = (sum - float64(n)*mean) / (stdDev * math.Sqrt(float64(n)))
    }
    
    return standardizedSums
}

// 正态性检验 (简化版)
func NormalityTest(data []float64) bool {
    // 使用偏度和峰度进行简单的正态性检验
    skewness := calculateSkewness(data)
    kurtosis := calculateKurtosis(data)
    
    // 正态分布的偏度接近0，峰度接近3
    return math.Abs(skewness) < 0.5 && math.Abs(kurtosis-3) < 1.0
}

func calculateSkewness(data []float64) float64 {
    mean := calculateMean(data)
    variance := calculateVariance(data, mean)
    stdDev := math.Sqrt(variance)
    
    skewness := 0.0
    for _, value := range data {
        z := (value - mean) / stdDev
        skewness += z * z * z
    }
    
    return skewness / float64(len(data))
}

func calculateKurtosis(data []float64) float64 {
    mean := calculateMean(data)
    variance := calculateVariance(data, mean)
    stdDev := math.Sqrt(variance)
    
    kurtosis := 0.0
    for _, value := range data {
        z := (value - mean) / stdDev
        kurtosis += z * z * z * z
    }
    
    return kurtosis / float64(len(data))
}

func calculateMean(data []float64) float64 {
    sum := 0.0
    for _, value := range data {
        sum += value
    }
    return sum / float64(len(data))
}

func calculateVariance(data []float64, mean float64) float64 {
    variance := 0.0
    for _, value := range data {
        variance += (value - mean) * (value - mean)
    }
    return variance / float64(len(data))
}
```

## 5. 应用领域

### 5.1 机器学习

```go
// 朴素贝叶斯分类器
type NaiveBayesClassifier struct {
    classes     map[string]float64
    features    map[string]map[string]map[string]float64
    classCounts map[string]int
    totalCount  int
}

func NewNaiveBayesClassifier() *NaiveBayesClassifier {
    return &NaiveBayesClassifier{
        classes:     make(map[string]float64),
        features:    make(map[string]map[string]map[string]float64),
        classCounts: make(map[string]int),
        totalCount:  0,
    }
}

func (nb *NaiveBayesClassifier) Train(data []TrainingExample) {
    // 计算类别先验概率
    for _, example := range data {
        nb.classCounts[example.Class]++
        nb.totalCount++
    }
    
    for class, count := range nb.classCounts {
        nb.classes[class] = float64(count) / float64(nb.totalCount)
    }
    
    // 计算特征条件概率
    for _, example := range data {
        for feature, value := range example.Features {
            if nb.features[feature] == nil {
                nb.features[feature] = make(map[string]map[string]float64)
            }
            if nb.features[feature][example.Class] == nil {
                nb.features[feature][example.Class] = make(map[string]float64)
            }
            nb.features[feature][example.Class][value]++
        }
    }
    
    // 归一化条件概率
    for feature := range nb.features {
        for class := range nb.features[feature] {
            total := 0.0
            for _, count := range nb.features[feature][class] {
                total += count
            }
            for value := range nb.features[feature][class] {
                nb.features[feature][class][value] /= total
            }
        }
    }
}

func (nb *NaiveBayesClassifier) Predict(features map[string]string) string {
    bestClass := ""
    bestScore := math.Inf(-1)
    
    for class, prior := range nb.classes {
        score := math.Log(prior)
        
        for feature, value := range features {
            if prob, exists := nb.features[feature][class][value]; exists {
                score += math.Log(prob)
            } else {
                // 拉普拉斯平滑
                score += math.Log(0.001)
            }
        }
        
        if score > bestScore {
            bestScore = score
            bestClass = class
        }
    }
    
    return bestClass
}

type TrainingExample struct {
    Features map[string]string
    Class    string
}
```

### 5.2 风险评估

```go
// 风险价值 (VaR) 计算
type RiskAnalyzer struct {
    returns []float64
}

func NewRiskAnalyzer(returns []float64) *RiskAnalyzer {
    return &RiskAnalyzer{returns: returns}
}

func (ra *RiskAnalyzer) VaR(confidence float64) float64 {
    // 使用历史模拟法计算VaR
    sortedReturns := make([]float64, len(ra.returns))
    copy(sortedReturns, ra.returns)
    sort.Float64s(sortedReturns)
    
    index := int((1 - confidence) * float64(len(sortedReturns)))
    if index >= len(sortedReturns) {
        index = len(sortedReturns) - 1
    }
    
    return -sortedReturns[index] // 负号表示损失
}

func (ra *RiskAnalyzer) ExpectedShortfall(confidence float64) float64 {
    varValue := ra.VaR(confidence)
    
    sum := 0.0
    count := 0
    
    for _, ret := range ra.returns {
        if -ret >= varValue {
            sum += -ret
            count++
        }
    }
    
    if count == 0 {
        return varValue
    }
    
    return sum / float64(count)
}

// 蒙特卡洛风险模拟
func MonteCarloRiskSimulation(initialValue, drift, volatility float64, timeSteps, simulations int) []float64 {
    results := make([]float64, simulations)
    
    for i := 0; i < simulations; i++ {
        value := initialValue
        
        for t := 0; t < timeSteps; t++ {
            // 几何布朗运动
            dt := 1.0 / 252.0 // 假设日度数据
            z := rand.NormFloat64()
            value *= math.Exp((drift-0.5*volatility*volatility)*dt + volatility*math.Sqrt(dt)*z)
        }
        
        results[i] = value
    }
    
    return results
}
```

### 5.3 性能分析

```go
// 系统性能分析
type PerformanceAnalyzer struct {
    responseTimes []float64
}

func NewPerformanceAnalyzer(responseTimes []float64) *PerformanceAnalyzer {
    return &PerformanceAnalyzer{responseTimes: responseTimes}
}

func (pa *PerformanceAnalyzer) Percentile(p float64) float64 {
    sorted := make([]float64, len(pa.responseTimes))
    copy(sorted, pa.responseTimes)
    sort.Float64s(sorted)
    
    index := int(p * float64(len(sorted)-1))
    return sorted[index]
}

func (pa *PerformanceAnalyzer) ConfidenceInterval(confidence float64) (float64, float64) {
    mean := pa.Mean()
    stdDev := pa.StandardDeviation()
    n := float64(len(pa.responseTimes))
    
    // 使用正态分布近似
    z := 1.96 // 95%置信区间对应的z值
    margin := z * stdDev / math.Sqrt(n)
    
    return mean - margin, mean + margin
}

func (pa *PerformanceAnalyzer) Mean() float64 {
    sum := 0.0
    for _, rt := range pa.responseTimes {
        sum += rt
    }
    return sum / float64(len(pa.responseTimes))
}

func (pa *PerformanceAnalyzer) StandardDeviation() float64 {
    mean := pa.Mean()
    variance := 0.0
    
    for _, rt := range pa.responseTimes {
        variance += (rt - mean) * (rt - mean)
    }
    
    return math.Sqrt(variance / float64(len(pa.responseTimes)-1))
}

// 负载测试结果分析
func AnalyzeLoadTest(results []LoadTestResult) PerformanceReport {
    responseTimes := make([]float64, len(results))
    throughputs := make([]float64, len(results))
    
    for i, result := range results {
        responseTimes[i] = result.AverageResponseTime
        throughputs[i] = result.Throughput
    }
    
    analyzer := NewPerformanceAnalyzer(responseTimes)
    
    return PerformanceReport{
        MeanResponseTime:    analyzer.Mean(),
        P95ResponseTime:     analyzer.Percentile(0.95),
        P99ResponseTime:     analyzer.Percentile(0.99),
        StandardDeviation:   analyzer.StandardDeviation(),
        AverageThroughput:   calculateMean(throughputs),
        ConfidenceInterval:  analyzer.ConfidenceInterval(0.95),
    }
}

type LoadTestResult struct {
    AverageResponseTime float64
    Throughput          float64
}

type PerformanceReport struct {
    MeanResponseTime   float64
    P95ResponseTime    float64
    P99ResponseTime    float64
    StandardDeviation  float64
    AverageThroughput  float64
    ConfidenceInterval [2]float64
}
```

## 6. 总结

概率论为软件工程提供了重要的理论基础：

1. **理论基础**: 为随机性建模和不确定性分析提供数学基础
2. **算法设计**: 为随机算法、蒙特卡洛方法等提供理论支持
3. **数据分析**: 为统计推断、机器学习等提供概率框架
4. **风险评估**: 为系统可靠性、性能分析等提供量化方法

通过Go语言的实现，我们可以：
- 高效地生成和操作随机变量
- 实现各种概率分布
- 验证概率论定理
- 解决实际工程问题

概率论在现代软件工程中扮演着越来越重要的角色，特别是在大数据、机器学习、系统可靠性等领域有广泛应用。
