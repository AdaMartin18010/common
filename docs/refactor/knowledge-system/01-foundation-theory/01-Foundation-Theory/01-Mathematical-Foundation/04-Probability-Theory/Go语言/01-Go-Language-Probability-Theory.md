# Go语言在概率论中的应用

## 概述

Go语言在概率论领域具有独特优势，其高性能、并发处理能力和数值计算支持使其成为实现概率分布、随机变量、统计推断和蒙特卡洛方法的理想选择。从基础的概率计算到复杂的随机过程，从统计推断到概率建模，Go语言为概率论研究和应用提供了高效、可靠的技术基础。

## 核心组件

### 1. 概率分布 (Probability Distributions)

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
)

// 概率分布接口
type ProbabilityDistribution interface {
    PDF(x float64) float64
    CDF(x float64) float64
    Mean() float64
    Variance() float64
    Sample() float64
    Name() string
}

// 正态分布
type NormalDistribution struct {
    mu    float64 // 均值
    sigma float64 // 标准差
}

// 创建正态分布
func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    return &NormalDistribution{
        mu:    mu,
        sigma: sigma,
    }
}

// 概率密度函数
func (n *NormalDistribution) PDF(x float64) float64 {
    coefficient := 1.0 / (n.sigma * math.Sqrt(2*math.Pi))
    exponent := -0.5 * math.Pow((x-n.mu)/n.sigma, 2)
    return coefficient * math.Exp(exponent)
}

// 累积分布函数
func (n *NormalDistribution) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + erf(z/math.Sqrt(2)))
}

// 均值
func (n *NormalDistribution) Mean() float64 {
    return n.mu
}

// 方差
func (n *NormalDistribution) Variance() float64 {
    return n.sigma * n.sigma
}

// 采样
func (n *NormalDistribution) Sample() float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.mu + n.sigma*z0
}

// 分布名称
func (n *NormalDistribution) Name() string {
    return fmt.Sprintf("Normal(μ=%.2f, σ=%.2f)", n.mu, n.sigma)
}

// 指数分布
type ExponentialDistribution struct {
    lambda float64 // 参数
}

// 创建指数分布
func NewExponentialDistribution(lambda float64) *ExponentialDistribution {
    return &ExponentialDistribution{
        lambda: lambda,
    }
}

// 概率密度函数
func (e *ExponentialDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return e.lambda * math.Exp(-e.lambda*x)
}

// 累积分布函数
func (e *ExponentialDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return 1 - math.Exp(-e.lambda*x)
}

// 均值
func (e *ExponentialDistribution) Mean() float64 {
    return 1.0 / e.lambda
}

// 方差
func (e *ExponentialDistribution) Variance() float64 {
    return 1.0 / (e.lambda * e.lambda)
}

// 采样
func (e *ExponentialDistribution) Sample() float64 {
    u := rand.Float64()
    return -math.Log(1-u) / e.lambda
}

// 分布名称
func (e *ExponentialDistribution) Name() string {
    return fmt.Sprintf("Exponential(λ=%.2f)", e.lambda)
}

// 泊松分布
type PoissonDistribution struct {
    lambda float64 // 参数
}

// 创建泊松分布
func NewPoissonDistribution(lambda float64) *PoissonDistribution {
    return &PoissonDistribution{
        lambda: lambda,
    }
}

// 概率质量函数
func (p *PoissonDistribution) PMF(k int) float64 {
    if k < 0 {
        return 0
    }
    return math.Pow(p.lambda, float64(k)) * math.Exp(-p.lambda) / float64(factorial(k))
}

// 累积分布函数
func (p *PoissonDistribution) CDF(k int) float64 {
    if k < 0 {
        return 0
    }
    
    sum := 0.0
    for i := 0; i <= k; i++ {
        sum += p.PMF(i)
    }
    return sum
}

// 均值
func (p *PoissonDistribution) Mean() float64 {
    return p.lambda
}

// 方差
func (p *PoissonDistribution) Variance() float64 {
    return p.lambda
}

// 采样
func (p *PoissonDistribution) Sample() int {
    // 使用拒绝采样
    L := math.Exp(-p.lambda)
    k := 0
    p_val := 1.0
    
    for {
        k++
        u := rand.Float64()
        p_val *= u
        if p_val <= L {
            return k - 1
        }
    }
}

// 分布名称
func (p *PoissonDistribution) Name() string {
    return fmt.Sprintf("Poisson(λ=%.2f)", p.lambda)
}

// 辅助函数：阶乘
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 辅助函数：误差函数
func erf(x float64) float64 {
    // 近似实现
    a1 := 0.254829592
    a2 := -0.284496736
    a3 := 1.421413741
    a4 := -1.453152027
    a5 := 1.061405429
    p := 0.3275911
    
    sign := 1.0
    if x < 0 {
        sign = -1
        x = -x
    }
    
    t := 1.0 / (1.0 + p*x)
    y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)
    
    return sign * y
}
```

### 2. 随机变量 (Random Variables)

```go
package main

import (
    "fmt"
    "math"
    "sort"
)

// 随机变量
type RandomVariable struct {
    distribution ProbabilityDistribution
    samples      []float64
    name         string
}

// 创建随机变量
func NewRandomVariable(dist ProbabilityDistribution, name string) *RandomVariable {
    return &RandomVariable{
        distribution: dist,
        samples:      make([]float64, 0),
        name:         name,
    }
}

// 生成样本
func (rv *RandomVariable) GenerateSamples(n int) {
    rv.samples = make([]float64, n)
    for i := 0; i < n; i++ {
        rv.samples[i] = rv.distribution.Sample()
    }
}

// 计算样本均值
func (rv *RandomVariable) SampleMean() float64 {
    if len(rv.samples) == 0 {
        return 0
    }
    
    sum := 0.0
    for _, sample := range rv.samples {
        sum += sample
    }
    return sum / float64(len(rv.samples))
}

// 计算样本方差
func (rv *RandomVariable) SampleVariance() float64 {
    if len(rv.samples) < 2 {
        return 0
    }
    
    mean := rv.SampleMean()
    sum := 0.0
    
    for _, sample := range rv.samples {
        sum += math.Pow(sample-mean, 2)
    }
    
    return sum / float64(len(rv.samples)-1)
}

// 计算样本标准差
func (rv *RandomVariable) SampleStdDev() float64 {
    return math.Sqrt(rv.SampleVariance())
}

// 计算分位数
func (rv *RandomVariable) Quantile(p float64) float64 {
    if len(rv.samples) == 0 || p < 0 || p > 1 {
        return 0
    }
    
    sorted := make([]float64, len(rv.samples))
    copy(sorted, rv.samples)
    sort.Float64s(sorted)
    
    n := len(sorted)
    index := p * float64(n-1)
    
    if index == float64(int(index)) {
        return sorted[int(index)]
    }
    
    lower := int(index)
    upper := lower + 1
    weight := index - float64(lower)
    
    return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// 计算理论值与样本值的比较
func (rv *RandomVariable) CompareTheoretical() map[string]interface{} {
    return map[string]interface{}{
        "name":           rv.name,
        "theoretical_mean": rv.distribution.Mean(),
        "sample_mean":      rv.SampleMean(),
        "theoretical_variance": rv.distribution.Variance(),
        "sample_variance":      rv.SampleVariance(),
        "sample_count":         len(rv.samples),
    }
}

// 多变量随机变量
type MultivariateRandomVariable struct {
    variables []*RandomVariable
    name      string
}

// 创建多变量随机变量
func NewMultivariateRandomVariable(variables []*RandomVariable, name string) *MultivariateRandomVariable {
    return &MultivariateRandomVariable{
        variables: variables,
        name:      name,
    }
}

// 生成多变量样本
func (mrv *MultivariateRandomVariable) GenerateSamples(n int) {
    for _, variable := range mrv.variables {
        variable.GenerateSamples(n)
    }
}

// 计算协方差矩阵
func (mrv *MultivariateRandomVariable) CovarianceMatrix() [][]float64 {
    n := len(mrv.variables)
    if n == 0 {
        return nil
    }
    
    matrix := make([][]float64, n)
    for i := range matrix {
        matrix[i] = make([]float64, n)
    }
    
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i == j {
                matrix[i][j] = mrv.variables[i].SampleVariance()
            } else {
                matrix[i][j] = mrv.covariance(mrv.variables[i], mrv.variables[j])
            }
        }
    }
    
    return matrix
}

// 计算两个随机变量的协方差
func (mrv *MultivariateRandomVariable) covariance(x, y *RandomVariable) float64 {
    if len(x.samples) != len(y.samples) || len(x.samples) == 0 {
        return 0
    }
    
    meanX := x.SampleMean()
    meanY := y.SampleMean()
    
    sum := 0.0
    for i := 0; i < len(x.samples); i++ {
        sum += (x.samples[i] - meanX) * (y.samples[i] - meanY)
    }
    
    return sum / float64(len(x.samples)-1)
}
```

### 3. 统计推断 (Statistical Inference)

```go
package main

import (
    "fmt"
    "math"
)

// 统计推断器
type StatisticalInference struct{}

// 置信区间
type ConfidenceInterval struct {
    Lower    float64
    Upper    float64
    Level    float64
    Mean     float64
    StdError float64
}

// 计算置信区间
func (si *StatisticalInference) ConfidenceInterval(samples []float64, confidence float64) *ConfidenceInterval {
    if len(samples) == 0 {
        return nil
    }
    
    n := float64(len(samples))
    mean := 0.0
    for _, sample := range samples {
        mean += sample
    }
    mean /= n
    
    variance := 0.0
    for _, sample := range samples {
        variance += math.Pow(sample-mean, 2)
    }
    variance /= (n - 1)
    
    stdDev := math.Sqrt(variance)
    stdError := stdDev / math.Sqrt(n)
    
    // 使用正态分布近似（大样本）
    z := si.zScore(confidence)
    margin := z * stdError
    
    return &ConfidenceInterval{
        Lower:    mean - margin,
        Upper:    mean + margin,
        Level:    confidence,
        Mean:     mean,
        StdError: stdError,
    }
}

// 假设检验
type HypothesisTest struct {
    NullHypothesis    string
    AlternativeHypothesis string
    TestStatistic    float64
    PValue           float64
    Significance     float64
    RejectNull       bool
}

// 单样本t检验
func (si *StatisticalInference) OneSampleTTest(samples []float64, hypothesizedMean, significance float64) *HypothesisTest {
    if len(samples) < 2 {
        return nil
    }
    
    n := float64(len(samples))
    sampleMean := 0.0
    for _, sample := range samples {
        sampleMean += sample
    }
    sampleMean /= n
    
    variance := 0.0
    for _, sample := range samples {
        variance += math.Pow(sample-sampleMean, 2)
    }
    variance /= (n - 1)
    
    stdDev := math.Sqrt(variance)
    stdError := stdDev / math.Sqrt(n)
    
    tStat := (sampleMean - hypothesizedMean) / stdError
    pValue := si.tTestPValue(tStat, int(n-1))
    
    return &HypothesisTest{
        NullHypothesis:        fmt.Sprintf("μ = %.2f", hypothesizedMean),
        AlternativeHypothesis: fmt.Sprintf("μ ≠ %.2f", hypothesizedMean),
        TestStatistic:         tStat,
        PValue:                pValue,
        Significance:          significance,
        RejectNull:            pValue < significance,
    }
}

// 双样本t检验
func (si *StatisticalInference) TwoSampleTTest(samples1, samples2 []float64, significance float64) *HypothesisTest {
    if len(samples1) < 2 || len(samples2) < 2 {
        return nil
    }
    
    n1 := float64(len(samples1))
    n2 := float64(len(samples2))
    
    mean1 := 0.0
    for _, sample := range samples1 {
        mean1 += sample
    }
    mean1 /= n1
    
    mean2 := 0.0
    for _, sample := range samples2 {
        mean2 += sample
    }
    mean2 /= n2
    
    variance1 := 0.0
    for _, sample := range samples1 {
        variance1 += math.Pow(sample-mean1, 2)
    }
    variance1 /= (n1 - 1)
    
    variance2 := 0.0
    for _, sample := range samples2 {
        variance2 += math.Pow(sample-mean2, 2)
    }
    variance2 /= (n2 - 1)
    
    pooledVariance := ((n1-1)*variance1 + (n2-1)*variance2) / (n1 + n2 - 2)
    stdError := math.Sqrt(pooledVariance*(1/n1+1/n2))
    
    tStat := (mean1 - mean2) / stdError
    pValue := si.tTestPValue(tStat, int(n1+n2-2))
    
    return &HypothesisTest{
        NullHypothesis:        "μ₁ = μ₂",
        AlternativeHypothesis: "μ₁ ≠ μ₂",
        TestStatistic:         tStat,
        PValue:                pValue,
        Significance:          significance,
        RejectNull:            pValue < significance,
    }
}

// 辅助函数：z分数
func (si *StatisticalInference) zScore(confidence float64) float64 {
    // 常用置信水平的z分数
    zScores := map[float64]float64{
        0.90: 1.645,
        0.95: 1.96,
        0.99: 2.576,
    }
    
    if z, exists := zScores[confidence]; exists {
        return z
    }
    
    // 近似计算
    return 1.96 // 默认95%置信水平
}

// 辅助函数：t检验p值（近似）
func (si *StatisticalInference) tTestPValue(tStat float64, df int) float64 {
    // 使用正态分布近似（大样本）
    if df > 30 {
        return 2 * (1 - si.normalCDF(math.Abs(tStat)))
    }
    
    // 简化实现，实际应用中应使用更精确的t分布
    return 2 * (1 - si.normalCDF(math.Abs(tStat)))
}

// 辅助函数：标准正态分布CDF
func (si *StatisticalInference) normalCDF(x float64) float64 {
    return 0.5 * (1 + erf(x/math.Sqrt(2)))
}
```

### 4. 蒙特卡洛方法 (Monte Carlo Methods)

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
)

// 蒙特卡洛积分器
type MonteCarloIntegrator struct {
    random *rand.Rand
}

// 创建蒙特卡洛积分器
func NewMonteCarloIntegrator() *MonteCarloIntegrator {
    return &MonteCarloIntegrator{
        random: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// 一维积分
func (mci *MonteCarloIntegrator) Integrate1D(f func(float64) float64, a, b float64, n int) float64 {
    if n <= 0 {
        return 0
    }
    
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + (b-a)*mci.random.Float64()
        sum += f(x)
    }
    
    return (b - a) * sum / float64(n)
}

// 二维积分
func (mci *MonteCarloIntegrator) Integrate2D(f func(float64, float64) float64, a, b, c, d float64, n int) float64 {
    if n <= 0 {
        return 0
    }
    
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + (b-a)*mci.random.Float64()
        y := c + (d-c)*mci.random.Float64()
        sum += f(x, y)
    }
    
    return (b - a) * (d - c) * sum / float64(n)
}

// 蒙特卡洛模拟器
type MonteCarloSimulator struct {
    random *rand.Rand
}

// 创建蒙特卡洛模拟器
func NewMonteCarloSimulator() *MonteCarloSimulator {
    return &MonteCarloSimulator{
        random: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// 随机游走
func (mcs *MonteCarloSimulator) RandomWalk(steps int, stepSize float64) []float64 {
    path := make([]float64, steps+1)
    path[0] = 0.0
    
    for i := 1; i <= steps; i++ {
        if mcs.random.Float64() < 0.5 {
            path[i] = path[i-1] + stepSize
        } else {
            path[i] = path[i-1] - stepSize
        }
    }
    
    return path
}

// 布朗运动
func (mcs *MonteCarloSimulator) BrownianMotion(steps int, dt float64) []float64 {
    path := make([]float64, steps+1)
    path[0] = 0.0
    
    for i := 1; i <= steps; i++ {
        // 生成标准正态随机数
        u1 := mcs.random.Float64()
        u2 := mcs.random.Float64()
        z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
        
        path[i] = path[i-1] + z*math.Sqrt(dt)
    }
    
    return path
}

// 期权定价（Black-Scholes模型）
type OptionPricing struct {
    simulator *MonteCarloSimulator
}

// 创建期权定价器
func NewOptionPricing() *OptionPricing {
    return &OptionPricing{
        simulator: NewMonteCarloSimulator(),
    }
}

// 欧式看涨期权定价
func (op *OptionPricing) EuropeanCallOption(S0, K, T, r, sigma float64, nSimulations int) float64 {
    if nSimulations <= 0 {
        return 0
    }
    
    sum := 0.0
    for i := 0; i < nSimulations; i++ {
        // 生成股票价格路径
        ST := S0 * math.Exp((r-0.5*sigma*sigma)*T + sigma*math.Sqrt(T)*op.simulator.random.NormFloat64())
        
        // 计算期权收益
        payoff := math.Max(ST-K, 0)
        sum += payoff
    }
    
    // 贴现
    return math.Exp(-r*T) * sum / float64(nSimulations)
}

// 欧式看跌期权定价
func (op *OptionPricing) EuropeanPutOption(S0, K, T, r, sigma float64, nSimulations int) float64 {
    if nSimulations <= 0 {
        return 0
    }
    
    sum := 0.0
    for i := 0; i < nSimulations; i++ {
        // 生成股票价格路径
        ST := S0 * math.Exp((r-0.5*sigma*sigma)*T + sigma*math.Sqrt(T)*op.simulator.random.NormFloat64())
        
        // 计算期权收益
        payoff := math.Max(K-ST, 0)
        sum += payoff
    }
    
    // 贴现
    return math.Exp(-r*T) * sum / float64(nSimulations)
}

// 风险分析
type RiskAnalysis struct {
    simulator *MonteCarloSimulator
}

// 创建风险分析器
func NewRiskAnalysis() *RiskAnalysis {
    return &RiskAnalysis{
        simulator: NewMonteCarloSimulator(),
    }
}

// 计算VaR（风险价值）
func (ra *RiskAnalysis) ValueAtRisk(returns []float64, confidence float64) float64 {
    if len(returns) == 0 {
        return 0
    }
    
    // 排序
    sorted := make([]float64, len(returns))
    copy(sorted, returns)
    sort.Float64s(sorted)
    
    // 计算分位数
    index := int((1 - confidence) * float64(len(sorted)))
    if index >= len(sorted) {
        index = len(sorted) - 1
    }
    
    return -sorted[index] // 负号表示损失
}

// 计算CVaR（条件风险价值）
func (ra *RiskAnalysis) ConditionalValueAtRisk(returns []float64, confidence float64) float64 {
    if len(returns) == 0 {
        return 0
    }
    
    // 排序
    sorted := make([]float64, len(returns))
    copy(sorted, returns)
    sort.Float64s(sorted)
    
    // 计算VaR
    varIndex := int((1 - confidence) * float64(len(sorted)))
    if varIndex >= len(sorted) {
        varIndex = len(sorted) - 1
    }
    
    // 计算超过VaR的损失的平均值
    sum := 0.0
    count := 0
    for i := 0; i <= varIndex; i++ {
        sum += sorted[i]
        count++
    }
    
    if count == 0 {
        return 0
    }
    
    return -sum / float64(count) // 负号表示损失
}
```

## 实践应用

### 概率论分析平台

```go
package main

import (
    "fmt"
    "log"
)

// 概率论分析平台
type ProbabilityTheoryPlatform struct {
    distributions map[string]ProbabilityDistribution
    integrator    *MonteCarloIntegrator
    simulator     *MonteCarloSimulator
    optionPricing *OptionPricing
    riskAnalysis  *RiskAnalysis
    inference     *StatisticalInference
}

// 创建概率论分析平台
func NewProbabilityTheoryPlatform() *ProbabilityTheoryPlatform {
    return &ProbabilityTheoryPlatform{
        distributions: make(map[string]ProbabilityDistribution),
        integrator:    NewMonteCarloIntegrator(),
        simulator:     NewMonteCarloSimulator(),
        optionPricing: NewOptionPricing(),
        riskAnalysis:  NewRiskAnalysis(),
        inference:     &StatisticalInference{},
    }
}

// 概率分布演示
func (ptp *ProbabilityTheoryPlatform) DistributionDemo() {
    fmt.Println("=== Probability Distribution Demo ===")
    
    // 正态分布
    normal := NewNormalDistribution(0, 1)
    ptp.distributions["normal"] = normal
    
    fmt.Printf("Normal Distribution: %s\n", normal.Name())
    fmt.Printf("PDF(0) = %.4f\n", normal.PDF(0))
    fmt.Printf("CDF(1) = %.4f\n", normal.CDF(1))
    fmt.Printf("Mean = %.2f\n", normal.Mean())
    fmt.Printf("Variance = %.2f\n", normal.Variance())
    
    // 指数分布
    exponential := NewExponentialDistribution(2.0)
    ptp.distributions["exponential"] = exponential
    
    fmt.Printf("\nExponential Distribution: %s\n", exponential.Name())
    fmt.Printf("PDF(1) = %.4f\n", exponential.PDF(1))
    fmt.Printf("CDF(1) = %.4f\n", exponential.CDF(1))
    fmt.Printf("Mean = %.2f\n", exponential.Mean())
    fmt.Printf("Variance = %.2f\n", exponential.Variance())
    
    // 泊松分布
    poisson := NewPoissonDistribution(3.0)
    ptp.distributions["poisson"] = poisson
    
    fmt.Printf("\nPoisson Distribution: %s\n", poisson.Name())
    fmt.Printf("PMF(2) = %.4f\n", poisson.PMF(2))
    fmt.Printf("CDF(2) = %.4f\n", poisson.CDF(2))
    fmt.Printf("Mean = %.2f\n", poisson.Mean())
    fmt.Printf("Variance = %.2f\n", poisson.Variance())
}

// 随机变量演示
func (ptp *ProbabilityTheoryPlatform) RandomVariableDemo() {
    fmt.Println("=== Random Variable Demo ===")
    
    // 创建正态随机变量
    normalRV := NewRandomVariable(ptp.distributions["normal"], "Normal RV")
    normalRV.GenerateSamples(1000)
    
    fmt.Printf("Sample Mean: %.4f\n", normalRV.SampleMean())
    fmt.Printf("Sample Variance: %.4f\n", normalRV.SampleVariance())
    fmt.Printf("Sample Std Dev: %.4f\n", normalRV.SampleStdDev())
    fmt.Printf("95th Percentile: %.4f\n", normalRV.Quantile(0.95))
    
    comparison := normalRV.CompareTheoretical()
    fmt.Printf("\nComparison:\n")
    for key, value := range comparison {
        fmt.Printf("  %s: %v\n", key, value)
    }
}

// 统计推断演示
func (ptp *ProbabilityTheoryPlatform) StatisticalInferenceDemo() {
    fmt.Println("=== Statistical Inference Demo ===")
    
    // 生成样本
    normalRV := NewRandomVariable(ptp.distributions["normal"], "Normal RV")
    normalRV.GenerateSamples(100)
    
    // 置信区间
    ci := ptp.inference.ConfidenceInterval(normalRV.samples, 0.95)
    fmt.Printf("95%% Confidence Interval: [%.4f, %.4f]\n", ci.Lower, ci.Upper)
    fmt.Printf("Sample Mean: %.4f\n", ci.Mean)
    fmt.Printf("Standard Error: %.4f\n", ci.StdError)
    
    // 假设检验
    test := ptp.inference.OneSampleTTest(normalRV.samples, 0.0, 0.05)
    fmt.Printf("\nOne-Sample T-Test:\n")
    fmt.Printf("  Null Hypothesis: %s\n", test.NullHypothesis)
    fmt.Printf("  Alternative Hypothesis: %s\n", test.AlternativeHypothesis)
    fmt.Printf("  Test Statistic: %.4f\n", test.TestStatistic)
    fmt.Printf("  P-Value: %.4f\n", test.PValue)
    fmt.Printf("  Reject Null: %t\n", test.RejectNull)
}

// 蒙特卡洛方法演示
func (ptp *ProbabilityTheoryPlatform) MonteCarloDemo() {
    fmt.Println("=== Monte Carlo Methods Demo ===")
    
    // 积分计算
    f := func(x float64) float64 {
        return x * x // f(x) = x²
    }
    
    integral := ptp.integrator.Integrate1D(f, 0, 1, 10000)
    exact := 1.0 / 3.0
    fmt.Printf("Monte Carlo Integration of x² from 0 to 1:\n")
    fmt.Printf("  Estimated: %.6f\n", integral)
    fmt.Printf("  Exact: %.6f\n", exact)
    fmt.Printf("  Error: %.6f\n", math.Abs(integral-exact))
    
    // 随机游走
    walk := ptp.simulator.RandomWalk(100, 0.1)
    fmt.Printf("\nRandom Walk (last 5 steps):\n")
    for i := len(walk) - 5; i < len(walk); i++ {
        fmt.Printf("  Step %d: %.4f\n", i, walk[i])
    }
    
    // 布朗运动
    brownian := ptp.simulator.BrownianMotion(100, 0.01)
    fmt.Printf("\nBrownian Motion (last 5 steps):\n")
    for i := len(brownian) - 5; i < len(brownian); i++ {
        fmt.Printf("  Step %d: %.4f\n", i, brownian[i])
    }
}

// 期权定价演示
func (ptp *ProbabilityTheoryPlatform) OptionPricingDemo() {
    fmt.Println("=== Option Pricing Demo ===")
    
    S0 := 100.0  // 当前股价
    K := 100.0   // 执行价格
    T := 1.0     // 到期时间（年）
    r := 0.05    // 无风险利率
    sigma := 0.2 // 波动率
    nSim := 10000 // 模拟次数
    
    // 欧式看涨期权
    callPrice := ptp.optionPricing.EuropeanCallOption(S0, K, T, r, sigma, nSim)
    fmt.Printf("European Call Option Price: %.4f\n", callPrice)
    
    // 欧式看跌期权
    putPrice := ptp.optionPricing.EuropeanPutOption(S0, K, T, r, sigma, nSim)
    fmt.Printf("European Put Option Price: %.4f\n", putPrice)
    
    // 验证看跌-看涨平价
    parity := callPrice - putPrice - S0 + K*math.Exp(-r*T)
    fmt.Printf("Put-Call Parity Check: %.6f (should be close to 0)\n", parity)
}

// 风险分析演示
func (ptp *ProbabilityTheoryPlatform) RiskAnalysisDemo() {
    fmt.Println("=== Risk Analysis Demo ===")
    
    // 生成收益率数据
    normalRV := NewRandomVariable(ptp.distributions["normal"], "Returns")
    normalRV.GenerateSamples(1000)
    
    // 计算VaR和CVaR
    var95 := ptp.riskAnalysis.ValueAtRisk(normalRV.samples, 0.95)
    cvar95 := ptp.riskAnalysis.ConditionalValueAtRisk(normalRV.samples, 0.95)
    
    fmt.Printf("95%% VaR: %.4f\n", var95)
    fmt.Printf("95%% CVaR: %.4f\n", cvar95)
    
    var99 := ptp.riskAnalysis.ValueAtRisk(normalRV.samples, 0.99)
    cvar99 := ptp.riskAnalysis.ConditionalValueAtRisk(normalRV.samples, 0.99)
    
    fmt.Printf("99%% VaR: %.4f\n", var99)
    fmt.Printf("99%% CVaR: %.4f\n", cvar99)
}

// 综合演示
func (ptp *ProbabilityTheoryPlatform) ComprehensiveDemo() {
    fmt.Println("=== Probability Theory Comprehensive Demo ===")
    
    ptp.DistributionDemo()
    fmt.Println()
    
    ptp.RandomVariableDemo()
    fmt.Println()
    
    ptp.StatisticalInferenceDemo()
    fmt.Println()
    
    ptp.MonteCarloDemo()
    fmt.Println()
    
    ptp.OptionPricingDemo()
    fmt.Println()
    
    ptp.RiskAnalysisDemo()
    
    fmt.Println("=== Demo Completed ===")
}
```

## 设计原则

### 1. 数学正确性 (Mathematical Correctness)

- **概率公理**: 基于概率论公理的设计
- **分布性质**: 确保分布的正确性
- **统计推断**: 正确的统计方法实现
- **数值精度**: 高精度的数值计算

### 2. 性能优化 (Performance Optimization)

- **算法选择**: 选择高效的随机算法
- **并行计算**: 利用并发进行大规模模拟
- **内存管理**: 优化内存使用
- **缓存策略**: 缓存中间计算结果

### 3. 可扩展性 (Scalability)

- **模块化设计**: 将概率组件分离
- **接口抽象**: 定义统一的概率接口
- **插件架构**: 支持自定义分布
- **分布式处理**: 支持大规模模拟

### 4. 易用性 (Usability)

- **简洁API**: 提供简单易用的接口
- **错误处理**: 完善的错误处理和提示
- **文档支持**: 详细的使用文档和示例
- **可视化**: 提供概率分布的可视化

## 总结

Go语言在概率论领域提供了强大的工具和框架，通过其高性能、并发处理能力和数值计算支持，能够构建高效、可靠的概率论应用。从基础的概率分布到复杂的随机过程，从统计推断到风险分析，Go语言为概率论研究和应用提供了完整的技术栈。

通过合理的设计原则和最佳实践，可以构建出数学正确、性能优化、可扩展、易用的概率论分析平台，满足各种概率论研究和应用需求。
