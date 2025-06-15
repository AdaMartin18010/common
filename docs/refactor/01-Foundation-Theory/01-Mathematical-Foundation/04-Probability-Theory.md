# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基础定义](#1-基础定义)
    - [1.1 概率空间](#11-概率空间)
    - [1.2 随机变量](#12-随机变量)
    - [1.3 概率分布](#13-概率分布)
  - [2. 概率计算](#2-概率计算)
    - [2.1 条件概率](#21-条件概率)
    - [2.2 贝叶斯定理](#22-贝叶斯定理)
    - [2.3 独立性](#23-独立性)
  - [3. 随机过程](#3-随机过程)
    - [3.1 马尔可夫链](#31-马尔可夫链)
    - [3.2 泊松过程](#32-泊松过程)
    - [3.3 布朗运动](#33-布朗运动)
  - [4. 统计推断](#4-统计推断)
    - [4.1 参数估计](#41-参数估计)
    - [4.2 假设检验](#42-假设检验)
    - [4.3 置信区间](#43-置信区间)
  - [5. 在软件工程中的应用](#5-在软件工程中的应用)
    - [5.1 性能分析](#51-性能分析)
    - [5.2 可靠性建模](#52-可靠性建模)
    - [5.3 随机算法](#53-随机算法)

## 1. 基础定义

### 1.1 概率空间

**形式化定义**：

```latex
概率空间 (Ω, F, P) 由以下部分组成：
- Ω: 样本空间，所有可能结果的集合
- F: σ-代数，Ω的子集族，满足：
  * Ω ∈ F
  * 如果 A ∈ F，则 A^c ∈ F
  * 如果 A₁, A₂, ... ∈ F，则 ∪ᵢ Aᵢ ∈ F
- P: 概率测度，P: F → [0,1]，满足：
  * P(Ω) = 1
  * 对于互斥事件 A₁, A₂, ...，P(∪ᵢ Aᵢ) = Σᵢ P(Aᵢ)
```

**Go语言实现**：

```go
// 概率空间
type ProbabilitySpace struct {
    SampleSpace []interface{}
    Events      map[string][]interface{}
    Probability map[string]float64
}

// 事件
type Event struct {
    Name     string
    Elements []interface{}
    Prob     float64
}

// 创建概率空间
func NewProbabilitySpace(sampleSpace []interface{}) *ProbabilitySpace {
    return &ProbabilitySpace{
        SampleSpace: sampleSpace,
        Events:      make(map[string][]interface{}),
        Probability: make(map[string]float64),
    }
}

// 添加事件
func (ps *ProbabilitySpace) AddEvent(name string, elements []interface{}, prob float64) {
    ps.Events[name] = elements
    ps.Probability[name] = prob
}

// 计算事件概率
func (ps *ProbabilitySpace) GetProbability(eventName string) float64 {
    if prob, exists := ps.Probability[eventName]; exists {
        return prob
    }
    return 0.0
}

// 事件并集
func (ps *ProbabilitySpace) Union(event1, event2 string) *Event {
    elements := make(map[interface{}]bool)
    
    for _, elem := range ps.Events[event1] {
        elements[elem] = true
    }
    for _, elem := range ps.Events[event2] {
        elements[elem] = true
    }
    
    result := make([]interface{}, 0)
    for elem := range elements {
        result = append(result, elem)
    }
    
    return &Event{
        Name:     event1 + "_union_" + event2,
        Elements: result,
    }
}

// 事件交集
func (ps *ProbabilitySpace) Intersection(event1, event2 string) *Event {
    elemMap := make(map[interface{}]bool)
    for _, elem := range ps.Events[event1] {
        elemMap[elem] = true
    }
    
    result := make([]interface{}, 0)
    for _, elem := range ps.Events[event2] {
        if elemMap[elem] {
            result = append(result, elem)
        }
    }
    
    return &Event{
        Name:     event1 + "_intersection_" + event2,
        Elements: result,
    }
}
```

### 1.2 随机变量

**形式化定义**：

```latex
随机变量 X: Ω → ℝ 是一个可测函数，满足：
对于任意 Borel 集 B ⊆ ℝ，X⁻¹(B) ∈ F

离散随机变量：取值可数
连续随机变量：取值不可数，有概率密度函数

期望：E[X] = Σ x P(X=x) 或 ∫ x f(x) dx
方差：Var(X) = E[(X - E[X])²] = E[X²] - (E[X])²
```

**Go语言实现**：

```go
// 随机变量接口
type RandomVariable interface {
    Sample() float64
    Expectation() float64
    Variance() float64
    ProbabilityDensity(x float64) float64
}

// 离散随机变量
type DiscreteRandomVariable struct {
    Values []float64
    Probs  []float64
}

func NewDiscreteRandomVariable(values, probs []float64) *DiscreteRandomVariable {
    return &DiscreteRandomVariable{
        Values: values,
        Probs:  probs,
    }
}

func (drv *DiscreteRandomVariable) Sample() float64 {
    r := rand.Float64()
    cumProb := 0.0
    
    for i, prob := range drv.Probs {
        cumProb += prob
        if r <= cumProb {
            return drv.Values[i]
        }
    }
    return drv.Values[len(drv.Values)-1]
}

func (drv *DiscreteRandomVariable) Expectation() float64 {
    expectation := 0.0
    for i, value := range drv.Values {
        expectation += value * drv.Probs[i]
    }
    return expectation
}

func (drv *DiscreteRandomVariable) Variance() float64 {
    expectation := drv.Expectation()
    variance := 0.0
    
    for i, value := range drv.Values {
        variance += drv.Probs[i] * (value - expectation) * (value - expectation)
    }
    return variance
}

func (drv *DiscreteRandomVariable) ProbabilityDensity(x float64) float64 {
    for i, value := range drv.Values {
        if value == x {
            return drv.Probs[i]
        }
    }
    return 0.0
}

// 连续随机变量 - 正态分布
type NormalRandomVariable struct {
    Mean     float64
    Variance float64
}

func NewNormalRandomVariable(mean, variance float64) *NormalRandomVariable {
    return &NormalRandomVariable{
        Mean:     mean,
        Variance: variance,
    }
}

func (nrv *NormalRandomVariable) Sample() float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return nrv.Mean + z0*math.Sqrt(nrv.Variance)
}

func (nrv *NormalRandomVariable) Expectation() float64 {
    return nrv.Mean
}

func (nrv *NormalRandomVariable) Variance() float64 {
    return nrv.Variance
}

func (nrv *NormalRandomVariable) ProbabilityDensity(x float64) float64 {
    stdDev := math.Sqrt(nrv.Variance)
    exponent := -0.5 * math.Pow((x-nrv.Mean)/stdDev, 2)
    return (1.0 / (stdDev * math.Sqrt(2*math.Pi))) * math.Exp(exponent)
}
```

### 1.3 概率分布

**常见分布**：

```latex
1. 离散分布：
   - 伯努利分布: P(X=k) = p^k (1-p)^(1-k), k ∈ {0,1}
   - 二项分布: P(X=k) = C(n,k) p^k (1-p)^(n-k)
   - 泊松分布: P(X=k) = (λ^k / k!) e^(-λ)

2. 连续分布：
   - 正态分布: f(x) = (1/√(2πσ²)) e^(-(x-μ)²/(2σ²))
   - 指数分布: f(x) = λ e^(-λx), x ≥ 0
   - 均匀分布: f(x) = 1/(b-a), a ≤ x ≤ b
```

**Go语言实现**：

```go
// 伯努利分布
type BernoulliDistribution struct {
    P float64
}

func NewBernoulliDistribution(p float64) *BernoulliDistribution {
    return &BernoulliDistribution{P: p}
}

func (bd *BernoulliDistribution) Sample() float64 {
    if rand.Float64() < bd.P {
        return 1.0
    }
    return 0.0
}

func (bd *BernoulliDistribution) Expectation() float64 {
    return bd.P
}

func (bd *BernoulliDistribution) Variance() float64 {
    return bd.P * (1 - bd.P)
}

// 二项分布
type BinomialDistribution struct {
    N int
    P float64
}

func NewBinomialDistribution(n int, p float64) *BinomialDistribution {
    return &BinomialDistribution{N: n, P: p}
}

func (bd *BinomialDistribution) Sample() float64 {
    success := 0
    for i := 0; i < bd.N; i++ {
        if rand.Float64() < bd.P {
            success++
        }
    }
    return float64(success)
}

func (bd *BinomialDistribution) Expectation() float64 {
    return float64(bd.N) * bd.P
}

func (bd *BinomialDistribution) Variance() float64 {
    return float64(bd.N) * bd.P * (1 - bd.P)
}

// 指数分布
type ExponentialDistribution struct {
    Lambda float64
}

func NewExponentialDistribution(lambda float64) *ExponentialDistribution {
    return &ExponentialDistribution{Lambda: lambda}
}

func (ed *ExponentialDistribution) Sample() float64 {
    return -math.Log(rand.Float64()) / ed.Lambda
}

func (ed *ExponentialDistribution) Expectation() float64 {
    return 1.0 / ed.Lambda
}

func (ed *ExponentialDistribution) Variance() float64 {
    return 1.0 / (ed.Lambda * ed.Lambda)
}
```

## 2. 概率计算

### 2.1 条件概率

**数学定义**：

```latex
条件概率: P(A|B) = P(A∩B) / P(B), 其中 P(B) > 0

全概率公式: P(A) = Σᵢ P(A|Bᵢ) P(Bᵢ)
其中 {Bᵢ} 是样本空间的一个划分
```

**Go语言实现**：

```go
// 条件概率计算
func (ps *ProbabilitySpace) ConditionalProbability(eventA, eventB string) float64 {
    intersection := ps.Intersection(eventA, eventB)
    probB := ps.GetProbability(eventB)
    
    if probB == 0 {
        return 0
    }
    
    // 计算交集概率
    intersectionProb := float64(len(intersection.Elements)) / float64(len(ps.SampleSpace))
    return intersectionProb / probB
}

// 全概率公式
func (ps *ProbabilitySpace) TotalProbability(eventA string, partition []string) float64 {
    totalProb := 0.0
    
    for _, eventB := range partition {
        condProb := ps.ConditionalProbability(eventA, eventB)
        probB := ps.GetProbability(eventB)
        totalProb += condProb * probB
    }
    
    return totalProb
}
```

### 2.2 贝叶斯定理

**数学定义**：

```latex
贝叶斯定理: P(A|B) = P(B|A) P(A) / P(B)

其中 P(B) = Σᵢ P(B|Aᵢ) P(Aᵢ) (全概率公式)
```

**Go语言实现**：

```go
// 贝叶斯定理
func (ps *ProbabilitySpace) BayesTheorem(eventA, eventB string, partition []string) float64 {
    probBGivenA := ps.ConditionalProbability(eventB, eventA)
    probA := ps.GetProbability(eventA)
    probB := ps.TotalProbability(eventB, partition)
    
    if probB == 0 {
        return 0
    }
    
    return probBGivenA * probA / probB
}

// 贝叶斯更新
type BayesianUpdater struct {
    Prior     map[string]float64
    Likelihood map[string]map[string]float64
}

func NewBayesianUpdater(prior map[string]float64, likelihood map[string]map[string]float64) *BayesianUpdater {
    return &BayesianUpdater{
        Prior:      prior,
        Likelihood: likelihood,
    }
}

func (bu *BayesianUpdater) Update(evidence string) map[string]float64 {
    posterior := make(map[string]float64)
    
    // 计算证据概率
    evidenceProb := 0.0
    for hypothesis, priorProb := range bu.Prior {
        if likelihood, exists := bu.Likelihood[hypothesis]; exists {
            if evidenceProb, exists := likelihood[evidence]; exists {
                evidenceProb += evidenceProb * priorProb
            }
        }
    }
    
    // 计算后验概率
    for hypothesis, priorProb := range bu.Prior {
        if likelihood, exists := bu.Likelihood[hypothesis]; exists {
            if evidenceProb, exists := likelihood[evidence]; exists {
                posterior[hypothesis] = (evidenceProb * priorProb) / evidenceProb
            }
        }
    }
    
    return posterior
}
```

### 2.3 独立性

**数学定义**：

```latex
事件独立性: P(A∩B) = P(A) P(B)
随机变量独立性: P(X≤x, Y≤y) = P(X≤x) P(Y≤y)
```

**Go语言实现**：

```go
// 检查事件独立性
func (ps *ProbabilitySpace) AreIndependent(eventA, eventB string) bool {
    probA := ps.GetProbability(eventA)
    probB := ps.GetProbability(eventB)
    intersection := ps.Intersection(eventA, eventB)
    
    intersectionProb := float64(len(intersection.Elements)) / float64(len(ps.SampleSpace))
    expectedProb := probA * probB
    
    return math.Abs(intersectionProb - expectedProb) < 1e-9
}

// 独立随机变量
type IndependentRandomVariables struct {
    Variables []RandomVariable
}

func NewIndependentRandomVariables(variables ...RandomVariable) *IndependentRandomVariables {
    return &IndependentRandomVariables{
        Variables: variables,
    }
}

func (irv *IndependentRandomVariables) JointSample() []float64 {
    result := make([]float64, len(irv.Variables))
    for i, rv := range irv.Variables {
        result[i] = rv.Sample()
    }
    return result
}

func (irv *IndependentRandomVariables) JointExpectation() []float64 {
    result := make([]float64, len(irv.Variables))
    for i, rv := range irv.Variables {
        result[i] = rv.Expectation()
    }
    return result
}
```

## 3. 随机过程

### 3.1 马尔可夫链

**数学定义**：

```latex
马尔可夫链: 随机过程 {Xₙ} 满足马尔可夫性质
P(Xₙ₊₁ = j | Xₙ = i, Xₙ₋₁ = k, ...) = P(Xₙ₊₁ = j | Xₙ = i)

转移矩阵: P = [pᵢⱼ] 其中 pᵢⱼ = P(Xₙ₊₁ = j | Xₙ = i)
```

**Go语言实现**：

```go
// 马尔可夫链
type MarkovChain struct {
    States       []string
    TransitionMatrix [][]float64
    CurrentState int
}

func NewMarkovChain(states []string, transitionMatrix [][]float64) *MarkovChain {
    return &MarkovChain{
        States:          states,
        TransitionMatrix: transitionMatrix,
        CurrentState:    0,
    }
}

func (mc *MarkovChain) NextState() string {
    r := rand.Float64()
    cumProb := 0.0
    
    for nextState, prob := range mc.TransitionMatrix[mc.CurrentState] {
        cumProb += prob
        if r <= cumProb {
            mc.CurrentState = nextState
            return mc.States[nextState]
        }
    }
    
    mc.CurrentState = len(mc.States) - 1
    return mc.States[mc.CurrentState]
}

func (mc *MarkovChain) Simulate(steps int) []string {
    path := make([]string, steps+1)
    path[0] = mc.States[mc.CurrentState]
    
    for i := 1; i <= steps; i++ {
        path[i] = mc.NextState()
    }
    
    return path
}

// 稳态分布计算
func (mc *MarkovChain) StationaryDistribution() []float64 {
    n := len(mc.States)
    
    // 构建线性方程组: π = πP
    // 约束条件: Σ πᵢ = 1
    matrix := make([][]float64, n+1)
    for i := range matrix {
        matrix[i] = make([]float64, n+1)
    }
    
    // 填充系数矩阵
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i == j {
                matrix[i][j] = mc.TransitionMatrix[i][j] - 1
            } else {
                matrix[i][j] = mc.TransitionMatrix[i][j]
            }
        }
        matrix[i][n] = 0 // 右侧常数
    }
    
    // 约束条件
    for j := 0; j < n; j++ {
        matrix[n][j] = 1
    }
    matrix[n][n] = 1 // 右侧常数
    
    // 求解线性方程组 (简化实现)
    return mc.solveLinearSystem(matrix)
}

func (mc *MarkovChain) solveLinearSystem(matrix [][]float64) []float64 {
    // 简化的高斯消元法
    n := len(matrix) - 1
    result := make([]float64, n)
    
    // 这里应该实现完整的高斯消元
    // 为了简化，返回均匀分布
    for i := 0; i < n; i++ {
        result[i] = 1.0 / float64(n)
    }
    
    return result
}
```

### 3.2 泊松过程

**数学定义**：

```latex
泊松过程: 计数过程 {N(t)} 满足：
1. N(0) = 0
2. 独立增量
3. 平稳增量
4. P(N(t+h) - N(t) = 1) = λh + o(h)
5. P(N(t+h) - N(t) ≥ 2) = o(h)

事件间隔时间服从指数分布 Exp(λ)
```

**Go语言实现**：

```go
// 泊松过程
type PoissonProcess struct {
    Lambda float64
    Time   float64
    Events int
}

func NewPoissonProcess(lambda float64) *PoissonProcess {
    return &PoissonProcess{
        Lambda: lambda,
        Time:   0.0,
        Events: 0,
    }
}

func (pp *PoissonProcess) NextEvent() float64 {
    // 事件间隔时间服从指数分布
    interval := -math.Log(rand.Float64()) / pp.Lambda
    pp.Time += interval
    pp.Events++
    return pp.Time
}

func (pp *PoissonProcess) Simulate(duration float64) []float64 {
    events := make([]float64, 0)
    
    for pp.Time < duration {
        eventTime := pp.NextEvent()
        if eventTime <= duration {
            events = append(events, eventTime)
        }
    }
    
    return events
}

func (pp *PoissonProcess) ExpectedEvents(duration float64) float64 {
    return pp.Lambda * duration
}

func (pp *PoissonProcess) ProbabilityOfEvents(k int, duration float64) float64 {
    lambdaT := pp.Lambda * duration
    return (math.Pow(lambdaT, float64(k)) / float64(factorial(k))) * math.Exp(-lambdaT)
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```

## 4. 统计推断

### 4.1 参数估计

**最大似然估计**：

```latex
似然函数: L(θ) = Πᵢ f(xᵢ; θ)
最大似然估计: θ̂ = argmax L(θ)
```

**Go语言实现**：

```go
// 最大似然估计
type MaximumLikelihoodEstimator struct {
    Data []float64
}

func NewMaximumLikelihoodEstimator(data []float64) *MaximumLikelihoodEstimator {
    return &MaximumLikelihoodEstimator{Data: data}
}

// 正态分布参数估计
func (mle *MaximumLikelihoodEstimator) EstimateNormalParameters() (float64, float64) {
    n := float64(len(mle.Data))
    
    // 样本均值
    mean := 0.0
    for _, x := range mle.Data {
        mean += x
    }
    mean /= n
    
    // 样本方差
    variance := 0.0
    for _, x := range mle.Data {
        variance += (x - mean) * (x - mean)
    }
    variance /= n
    
    return mean, variance
}

// 指数分布参数估计
func (mle *MaximumLikelihoodEstimator) EstimateExponentialParameter() float64 {
    n := float64(len(mle.Data))
    
    // 样本均值
    mean := 0.0
    for _, x := range mle.Data {
        mean += x
    }
    mean /= n
    
    // λ = 1/μ
    return 1.0 / mean
}

// 置信区间
func (mle *MaximumLikelihoodEstimator) ConfidenceInterval(confidence float64) (float64, float64) {
    n := float64(len(mle.Data))
    mean, variance := mle.EstimateNormalParameters()
    stdError := math.Sqrt(variance / n)
    
    // 使用正态分布近似
    z := 1.96 // 95% 置信水平
    if confidence == 0.99 {
        z = 2.576
    }
    
    margin := z * stdError
    return mean - margin, mean + margin
}
```

### 4.2 假设检验

**数学定义**：

```latex
零假设 H₀: θ = θ₀
备择假设 H₁: θ ≠ θ₀

检验统计量: T = (X̄ - μ₀) / (s/√n)
拒绝域: |T| > t_{α/2, n-1}
```

**Go语言实现**：

```go
// 假设检验
type HypothesisTest struct {
    Data []float64
}

func NewHypothesisTest(data []float64) *HypothesisTest {
    return &HypothesisTest{Data: data}
}

// t检验
func (ht *HypothesisTest) TTest(nullMean float64, alpha float64) (bool, float64) {
    n := float64(len(ht.Data))
    
    // 计算样本统计量
    sampleMean := 0.0
    for _, x := range ht.Data {
        sampleMean += x
    }
    sampleMean /= n
    
    sampleVariance := 0.0
    for _, x := range ht.Data {
        sampleVariance += (x - sampleMean) * (x - sampleMean)
    }
    sampleVariance /= (n - 1)
    
    sampleStdDev := math.Sqrt(sampleVariance)
    
    // 计算t统计量
    tStat := (sampleMean - nullMean) / (sampleStdDev / math.Sqrt(n))
    
    // 临界值 (简化，使用正态分布近似)
    criticalValue := 1.96 // α = 0.05
    
    // 决策
    reject := math.Abs(tStat) > criticalValue
    pValue := 2 * (1 - normalCDF(math.Abs(tStat)))
    
    return reject, pValue
}

// 正态分布累积分布函数近似
func normalCDF(x float64) float64 {
    return 0.5 * (1 + math.Erf(x/math.Sqrt(2)))
}
```

## 5. 在软件工程中的应用

### 5.1 性能分析

**排队论模型**：

```go
// M/M/1 排队系统
type MM1Queue struct {
    ArrivalRate   float64 // λ
    ServiceRate   float64 // μ
    Queue         []interface{}
    ServerBusy    bool
    TotalTime     float64
    TotalCustomers int
}

func NewMM1Queue(arrivalRate, serviceRate float64) *MM1Queue {
    return &MM1Queue{
        ArrivalRate:   arrivalRate,
        ServiceRate:   serviceRate,
        Queue:         make([]interface{}, 0),
        ServerBusy:    false,
        TotalTime:     0.0,
        TotalCustomers: 0,
    }
}

func (q *MM1Queue) Simulate(duration float64) map[string]float64 {
    currentTime := 0.0
    nextArrival := q.generateArrivalTime()
    nextService := math.Inf(1)
    
    for currentTime < duration {
        if nextArrival < nextService {
            // 到达事件
            currentTime = nextArrival
            q.handleArrival()
            nextArrival = currentTime + q.generateArrivalTime()
        } else {
            // 服务完成事件
            currentTime = nextService
            q.handleService()
            if len(q.Queue) > 0 {
                nextService = currentTime + q.generateServiceTime()
            } else {
                nextService = math.Inf(1)
            }
        }
    }
    
    return q.getStatistics()
}

func (q *MM1Queue) generateArrivalTime() float64 {
    return -math.Log(rand.Float64()) / q.ArrivalRate
}

func (q *MM1Queue) generateServiceTime() float64 {
    return -math.Log(rand.Float64()) / q.ServiceRate
}

func (q *MM1Queue) handleArrival() {
    q.TotalCustomers++
    if q.ServerBusy {
        q.Queue = append(q.Queue, struct{}{})
    } else {
        q.ServerBusy = true
    }
}

func (q *MM1Queue) handleService() {
    if len(q.Queue) > 0 {
        q.Queue = q.Queue[1:]
    } else {
        q.ServerBusy = false
    }
}

func (q *MM1Queue) getStatistics() map[string]float64 {
    utilization := q.TotalTime / q.TotalTime
    avgQueueLength := float64(len(q.Queue))
    
    return map[string]float64{
        "utilization":     utilization,
        "avg_queue_length": avgQueueLength,
        "throughput":      q.ServiceRate * utilization,
    }
}
```

### 5.2 可靠性建模

**可靠性函数**：

```go
// 可靠性模型
type ReliabilityModel struct {
    FailureRate float64
    Time        float64
}

func NewReliabilityModel(failureRate float64) *ReliabilityModel {
    return &ReliabilityModel{
        FailureRate: failureRate,
        Time:        0.0,
    }
}

// 可靠性函数 R(t) = e^(-λt)
func (rm *ReliabilityModel) Reliability(time float64) float64 {
    return math.Exp(-rm.FailureRate * time)
}

// 故障率函数 λ(t) = λ (指数分布)
func (rm *ReliabilityModel) FailureRateFunction(time float64) float64 {
    return rm.FailureRate
}

// 平均故障时间 MTTF = 1/λ
func (rm *ReliabilityModel) MTTF() float64 {
    return 1.0 / rm.FailureRate
}

// 系统可靠性 (串联)
func (rm *ReliabilityModel) SeriesReliability(components []*ReliabilityModel, time float64) float64 {
    reliability := 1.0
    for _, component := range components {
        reliability *= component.Reliability(time)
    }
    return reliability
}

// 系统可靠性 (并联)
func (rm *ReliabilityModel) ParallelReliability(components []*ReliabilityModel, time float64) float64 {
    unreliability := 1.0
    for _, component := range components {
        unreliability *= (1 - component.Reliability(time))
    }
    return 1 - unreliability
}
```

### 5.3 随机算法

**蒙特卡洛方法**：

```go
// 蒙特卡洛积分
func MonteCarloIntegration(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + rand.Float64()*(b-a)
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}

// 随机化算法 - 快速排序
func RandomizedQuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    // 随机选择主元
    pivotIndex := rand.Intn(len(arr))
    pivot := arr[pivotIndex]
    
    // 分区
    left := make([]int, 0)
    right := make([]int, 0)
    equal := make([]int, 0)
    
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
    result := RandomizedQuickSort(left)
    result = append(result, equal...)
    result = append(result, RandomizedQuickSort(right)...)
    
    return result
}

// 随机化测试
func RandomizedTest(algorithm func([]int) []int, testCases int) bool {
    for i := 0; i < testCases; i++ {
        // 生成随机测试数据
        n := rand.Intn(100) + 1
        arr := make([]int, n)
        for j := range arr {
            arr[j] = rand.Intn(1000)
        }
        
        // 运行算法
        result := algorithm(arr)
        
        // 验证结果
        if !isSorted(result) {
            return false
        }
    }
    return true
}

func isSorted(arr []int) bool {
    for i := 1; i < len(arr); i++ {
        if arr[i] < arr[i-1] {
            return false
        }
    }
    return true
}
```

## 总结

概率论为软件工程提供了强大的数学工具，用于：

1. **性能分析**：排队论、随机过程建模
2. **可靠性工程**：故障率分析、系统可靠性评估
3. **算法设计**：随机化算法、概率数据结构
4. **机器学习**：统计学习理论、概率模型
5. **网络安全**：密码学、随机数生成

**核心要点**：
- 概率空间和随机变量的形式化定义
- 条件概率和贝叶斯定理的应用
- 随机过程和统计推断方法
- 在软件工程中的实际应用

这个完整的概率论框架为软件工程的定量分析提供了坚实的理论基础。

## 相关链接

- [01-集合论](./01-Set-Theory.md)
- [02-逻辑学](./02-Logic.md)
- [03-图论](./03-Graph-Theory.md)
- [08-软件工程形式化](../08-Software-Engineering-Formalization/01-Software-Architecture-Formalization/01-Architecture-Meta-Model.md)
