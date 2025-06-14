# 04-概率论 (Probability Theory)

## 概述

概率论是数学的重要分支，为随机现象建模提供理论基础。本文档介绍概率论的基本概念、随机变量、分布函数以及在Go语言中的实现。

## 目录

1. [概率空间 (Probability Space)](#1-概率空间-probability-space)
2. [随机变量 (Random Variables)](#2-随机变量-random-variables)
3. [分布函数 (Distribution Functions)](#3-分布函数-distribution-functions)
4. [期望与方差 (Expectation and Variance)](#4-期望与方差-expectation-and-variance)
5. [大数定律与中心极限定理 (Law of Large Numbers and Central Limit Theorem)](#5-大数定律与中心极限定理-law-of-large-numbers-and-central-limit-theorem)
6. [Go语言中的概率实现](#6-go语言中的概率实现)

---

## 1. 概率空间 (Probability Space)

### 1.1 基本概念

**定义 1.1.1** (样本空间)
样本空间 $\Omega$ 是随机试验所有可能结果的集合。

**定义 1.1.2** (事件)
事件是样本空间的子集 $A \subseteq \Omega$。

**定义 1.1.3** (σ-代数)
σ-代数 $\mathcal{F}$ 是 $\Omega$ 的子集族，满足：
1. $\Omega \in \mathcal{F}$
2. 如果 $A \in \mathcal{F}$，则 $A^c \in \mathcal{F}$
3. 如果 $A_1, A_2, ... \in \mathcal{F}$，则 $\bigcup_{i=1}^{\infty} A_i \in \mathcal{F}$

**定义 1.1.4** (概率测度)
概率测度 $P: \mathcal{F} \rightarrow [0,1]$ 满足：
1. $P(\Omega) = 1$
2. 对于互斥事件 $A_1, A_2, ...$，$P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$

### 1.2 概率公理

**公理 1.2.1** (Kolmogorov公理)
1. 非负性: $P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
2. 规范性: $P(\Omega) = 1$
3. 可列可加性: 对于互斥事件 $A_1, A_2, ...$，$P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$

### 1.3 Go语言实现

```go
package probability

import (
    "math"
    "math/rand"
    "sort"
)

// SampleSpace 样本空间
type SampleSpace struct {
    Elements []interface{}
}

// Event 事件
type Event struct {
    Elements map[interface{}]bool
}

// NewEvent 创建事件
func NewEvent(elements ...interface{}) *Event {
    event := &Event{
        Elements: make(map[interface{}]bool),
    }
    for _, element := range elements {
        event.Elements[element] = true
    }
    return event
}

// ProbabilitySpace 概率空间
type ProbabilitySpace struct {
    SampleSpace *SampleSpace
    Events      map[string]*Event
    Probabilities map[string]float64
}

// NewProbabilitySpace 创建概率空间
func NewProbabilitySpace(elements []interface{}) *ProbabilitySpace {
    return &ProbabilitySpace{
        SampleSpace: &SampleSpace{Elements: elements},
        Events:      make(map[string]*Event),
        Probabilities: make(map[string]float64),
    }
}

// AddEvent 添加事件
func (ps *ProbabilitySpace) AddEvent(name string, event *Event, probability float64) {
    ps.Events[name] = event
    ps.Probabilities[name] = probability
}

// GetProbability 获取事件概率
func (ps *ProbabilitySpace) GetProbability(eventName string) float64 {
    return ps.Probabilities[eventName]
}

// Union 事件并集
func (ps *ProbabilitySpace) Union(event1, event2 *Event) *Event {
    union := &Event{
        Elements: make(map[interface{}]bool),
    }
    
    for element := range event1.Elements {
        union.Elements[element] = true
    }
    for element := range event2.Elements {
        union.Elements[element] = true
    }
    
    return union
}

// Intersection 事件交集
func (ps *ProbabilitySpace) Intersection(event1, event2 *Event) *Event {
    intersection := &Event{
        Elements: make(map[interface{}]bool),
    }
    
    for element := range event1.Elements {
        if event2.Elements[element] {
            intersection.Elements[element] = true
        }
    }
    
    return intersection
}

// Complement 事件补集
func (ps *ProbabilitySpace) Complement(event *Event) *Event {
    complement := &Event{
        Elements: make(map[interface{}]bool),
    }
    
    for _, element := range ps.SampleSpace.Elements {
        if !event.Elements[element] {
            complement.Elements[element] = true
        }
    }
    
    return complement
}

// ConditionalProbability 条件概率
func (ps *ProbabilitySpace) ConditionalProbability(eventA, eventB *Event) float64 {
    intersection := ps.Intersection(eventA, eventB)
    intersectionProb := float64(len(intersection.Elements)) / float64(len(ps.SampleSpace.Elements))
    eventBProb := float64(len(eventB.Elements)) / float64(len(ps.SampleSpace.Elements))
    
    if eventBProb == 0 {
        return 0
    }
    
    return intersectionProb / eventBProb
}
```

---

## 2. 随机变量 (Random Variables)

### 2.1 定义

**定义 2.1.1** (随机变量)
随机变量 $X: \Omega \rightarrow \mathbb{R}$ 是从样本空间到实数的可测函数。

**定义 2.1.2** (离散随机变量)
取有限或可列个值的随机变量。

**定义 2.1.3** (连续随机变量)
取值在连续区间上的随机变量。

### 2.2 概率质量函数 (PMF)

**定义 2.2.1** (PMF)
对于离散随机变量 $X$，概率质量函数 $p_X(x) = P(X = x)$。

**性质**:
1. $p_X(x) \geq 0$ 对所有 $x$
2. $\sum_x p_X(x) = 1$

### 2.3 概率密度函数 (PDF)

**定义 2.3.1** (PDF)
对于连续随机变量 $X$，概率密度函数 $f_X(x)$ 满足：
$P(a \leq X \leq b) = \int_a^b f_X(x) dx$

**性质**:
1. $f_X(x) \geq 0$ 对所有 $x$
2. $\int_{-\infty}^{\infty} f_X(x) dx = 1$

### 2.4 Go语言实现

```go
package random_variables

import (
    "math"
    "math/rand"
)

// RandomVariable 随机变量接口
type RandomVariable interface {
    Sample() float64
    PMF(x float64) float64
    CDF(x float64) float64
    Mean() float64
    Variance() float64
}

// DiscreteRandomVariable 离散随机变量
type DiscreteRandomVariable struct {
    Values     []float64
    Probabilities []float64
}

// NewDiscreteRandomVariable 创建离散随机变量
func NewDiscreteRandomVariable(values, probabilities []float64) *DiscreteRandomVariable {
    return &DiscreteRandomVariable{
        Values:       values,
        Probabilities: probabilities,
    }
}

func (drv *DiscreteRandomVariable) Sample() float64 {
    r := rand.Float64()
    cumulative := 0.0
    
    for i, prob := range drv.Probabilities {
        cumulative += prob
        if r <= cumulative {
            return drv.Values[i]
        }
    }
    
    return drv.Values[len(drv.Values)-1]
}

func (drv *DiscreteRandomVariable) PMF(x float64) float64 {
    for i, value := range drv.Values {
        if value == x {
            return drv.Probabilities[i]
        }
    }
    return 0
}

func (drv *DiscreteRandomVariable) CDF(x float64) float64 {
    cumulative := 0.0
    for i, value := range drv.Values {
        if value <= x {
            cumulative += drv.Probabilities[i]
        }
    }
    return cumulative
}

func (drv *DiscreteRandomVariable) Mean() float64 {
    mean := 0.0
    for i, value := range drv.Values {
        mean += value * drv.Probabilities[i]
    }
    return mean
}

func (drv *DiscreteRandomVariable) Variance() float64 {
    mean := drv.Mean()
    variance := 0.0
    for i, value := range drv.Values {
        variance += math.Pow(value-mean, 2) * drv.Probabilities[i]
    }
    return variance
}

// ContinuousRandomVariable 连续随机变量
type ContinuousRandomVariable struct {
    PDF func(float64) float64
    CDF func(float64) float64
    Mean float64
    Variance float64
}

// NewContinuousRandomVariable 创建连续随机变量
func NewContinuousRandomVariable(pdf func(float64) float64, mean, variance float64) *ContinuousRandomVariable {
    return &ContinuousRandomVariable{
        PDF: pdf,
        Mean: mean,
        Variance: variance,
    }
}

func (crv *ContinuousRandomVariable) Sample() float64 {
    // 使用逆变换采样
    r := rand.Float64()
    return crv.inverseCDF(r)
}

func (crv *ContinuousRandomVariable) PMF(x float64) float64 {
    // 连续随机变量的PMF为0
    return 0
}

func (crv *ContinuousRandomVariable) CDF(x float64) float64 {
    if crv.CDF != nil {
        return crv.CDF(x)
    }
    
    // 数值积分计算CDF
    return crv.numericalCDF(x)
}

func (crv *ContinuousRandomVariable) Mean() float64 {
    return crv.Mean
}

func (crv *ContinuousRandomVariable) Variance() float64 {
    return crv.Variance
}

// numericalCDF 数值计算CDF
func (crv *ContinuousRandomVariable) numericalCDF(x float64) float64 {
    const dx = 0.001
    integral := 0.0
    
    for t := -10.0; t <= x; t += dx {
        integral += crv.PDF(t) * dx
    }
    
    return integral
}

// inverseCDF 逆CDF函数（数值方法）
func (crv *ContinuousRandomVariable) inverseCDF(p float64) float64 {
    const tolerance = 1e-6
    const maxIterations = 100
    
    // 二分搜索
    left, right := -10.0, 10.0
    
    for i := 0; i < maxIterations; i++ {
        mid := (left + right) / 2
        cdf := crv.CDF(mid)
        
        if math.Abs(cdf-p) < tolerance {
            return mid
        }
        
        if cdf < p {
            left = mid
        } else {
            right = mid
        }
    }
    
    return (left + right) / 2
}
```

---

## 3. 分布函数 (Distribution Functions)

### 3.1 常见离散分布

#### 3.1.1 伯努利分布

**定义 3.1.1** (伯努利分布)
$X \sim Bernoulli(p)$ 的概率质量函数：
$P(X = k) = p^k(1-p)^{1-k}$，其中 $k \in \{0,1\}$

**期望**: $E[X] = p$
**方差**: $Var(X) = p(1-p)$

#### 3.1.2 二项分布

**定义 3.1.2** (二项分布)
$X \sim Binomial(n,p)$ 的概率质量函数：
$P(X = k) = \binom{n}{k} p^k(1-p)^{n-k}$，其中 $k \in \{0,1,...,n\}$

**期望**: $E[X] = np$
**方差**: $Var(X) = np(1-p)$

#### 3.1.3 泊松分布

**定义 3.1.3** (泊松分布)
$X \sim Poisson(\lambda)$ 的概率质量函数：
$P(X = k) = \frac{\lambda^k e^{-\lambda}}{k!}$，其中 $k \in \{0,1,2,...\}$

**期望**: $E[X] = \lambda$
**方差**: $Var(X) = \lambda$

### 3.2 常见连续分布

#### 3.2.1 均匀分布

**定义 3.2.1** (均匀分布)
$X \sim Uniform(a,b)$ 的概率密度函数：
$f_X(x) = \frac{1}{b-a}$，其中 $x \in [a,b]$

**期望**: $E[X] = \frac{a+b}{2}$
**方差**: $Var(X) = \frac{(b-a)^2}{12}$

#### 3.2.2 正态分布

**定义 3.2.2** (正态分布)
$X \sim Normal(\mu,\sigma^2)$ 的概率密度函数：
$f_X(x) = \frac{1}{\sqrt{2\pi\sigma^2}} e^{-\frac{(x-\mu)^2}{2\sigma^2}}$

**期望**: $E[X] = \mu$
**方差**: $Var(X) = \sigma^2$

#### 3.2.3 指数分布

**定义 3.2.3** (指数分布)
$X \sim Exponential(\lambda)$ 的概率密度函数：
$f_X(x) = \lambda e^{-\lambda x}$，其中 $x \geq 0$

**期望**: $E[X] = \frac{1}{\lambda}$
**方差**: $Var(X) = \frac{1}{\lambda^2}$

### 3.3 Go语言实现

```go
package distributions

import (
    "math"
    "math/rand"
)

// Bernoulli 伯努利分布
type Bernoulli struct {
    P float64
}

func NewBernoulli(p float64) *Bernoulli {
    return &Bernoulli{P: p}
}

func (b *Bernoulli) Sample() float64 {
    if rand.Float64() < b.P {
        return 1
    }
    return 0
}

func (b *Bernoulli) PMF(x float64) float64 {
    if x == 1 {
        return b.P
    } else if x == 0 {
        return 1 - b.P
    }
    return 0
}

func (b *Bernoulli) CDF(x float64) float64 {
    if x < 0 {
        return 0
    } else if x < 1 {
        return 1 - b.P
    }
    return 1
}

func (b *Bernoulli) Mean() float64 {
    return b.P
}

func (b *Bernoulli) Variance() float64 {
    return b.P * (1 - b.P)
}

// Binomial 二项分布
type Binomial struct {
    N int
    P float64
}

func NewBinomial(n int, p float64) *Binomial {
    return &Binomial{N: n, P: p}
}

func (b *Binomial) Sample() float64 {
    sum := 0
    for i := 0; i < b.N; i++ {
        if rand.Float64() < b.P {
            sum++
        }
    }
    return float64(sum)
}

func (b *Binomial) PMF(x float64) float64 {
    k := int(x)
    if k < 0 || k > b.N {
        return 0
    }
    
    // 计算组合数
    combination := b.combination(b.N, k)
    return combination * math.Pow(b.P, float64(k)) * math.Pow(1-b.P, float64(b.N-k))
}

func (b *Binomial) CDF(x float64) float64 {
    cumulative := 0.0
    for k := 0; k <= int(x) && k <= b.N; k++ {
        cumulative += b.PMF(float64(k))
    }
    return cumulative
}

func (b *Binomial) Mean() float64 {
    return float64(b.N) * b.P
}

func (b *Binomial) Variance() float64 {
    return float64(b.N) * b.P * (1 - b.P)
}

// combination 计算组合数
func (b *Binomial) combination(n, k int) float64 {
    if k > n-k {
        k = n - k
    }
    
    result := 1.0
    for i := 0; i < k; i++ {
        result *= float64(n-i) / float64(i+1)
    }
    return result
}

// Poisson 泊松分布
type Poisson struct {
    Lambda float64
}

func NewPoisson(lambda float64) *Poisson {
    return &Poisson{Lambda: lambda}
}

func (p *Poisson) Sample() float64 {
    // 使用逆变换采样
    r := rand.Float64()
    k := 0
    cumulative := math.Exp(-p.Lambda)
    
    for r > cumulative {
        k++
        cumulative += p.PMF(float64(k))
    }
    
    return float64(k)
}

func (p *Poisson) PMF(x float64) float64 {
    k := int(x)
    if k < 0 {
        return 0
    }
    
    return math.Pow(p.Lambda, float64(k)) * math.Exp(-p.Lambda) / float64(p.factorial(k))
}

func (p *Poisson) CDF(x float64) float64 {
    cumulative := 0.0
    for k := 0; k <= int(x); k++ {
        cumulative += p.PMF(float64(k))
    }
    return cumulative
}

func (p *Poisson) Mean() float64 {
    return p.Lambda
}

func (p *Poisson) Variance() float64 {
    return p.Lambda
}

// factorial 计算阶乘
func (p *Poisson) factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * p.factorial(n-1)
}

// Normal 正态分布
type Normal struct {
    Mu    float64
    Sigma float64
}

func NewNormal(mu, sigma float64) *Normal {
    return &Normal{Mu: mu, Sigma: sigma}
}

func (n *Normal) Sample() float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.Mu + n.Sigma*z0
}

func (n *Normal) PDF(x float64) float64 {
    exponent := -0.5 * math.Pow((x-n.Mu)/n.Sigma, 2)
    return math.Exp(exponent) / (n.Sigma * math.Sqrt(2*math.Pi))
}

func (n *Normal) CDF(x float64) float64 {
    z := (x - n.Mu) / n.Sigma
    return 0.5 * (1 + n.erf(z/math.Sqrt(2)))
}

func (n *Normal) Mean() float64 {
    return n.Mu
}

func (n *Normal) Variance() float64 {
    return n.Sigma * n.Sigma
}

// erf 误差函数近似
func (n *Normal) erf(x float64) float64 {
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

// Exponential 指数分布
type Exponential struct {
    Lambda float64
}

func NewExponential(lambda float64) *Exponential {
    return &Exponential{Lambda: lambda}
}

func (e *Exponential) Sample() float64 {
    return -math.Log(rand.Float64()) / e.Lambda
}

func (e *Exponential) PDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return e.Lambda * math.Exp(-e.Lambda*x)
}

func (e *Exponential) CDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return 1 - math.Exp(-e.Lambda*x)
}

func (e *Exponential) Mean() float64 {
    return 1 / e.Lambda
}

func (e *Exponential) Variance() float64 {
    return 1 / (e.Lambda * e.Lambda)
}
```

---

## 4. 期望与方差 (Expectation and Variance)

### 4.1 期望

**定义 4.1.1** (期望)
对于离散随机变量 $X$，期望 $E[X] = \sum_x x \cdot P(X = x)$
对于连续随机变量 $X$，期望 $E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$

**性质**:
1. 线性性: $E[aX + b] = aE[X] + b$
2. 可加性: $E[X + Y] = E[X] + E[Y]$

### 4.2 方差

**定义 4.2.1** (方差)
方差 $Var(X) = E[(X - E[X])^2] = E[X^2] - (E[X])^2$

**性质**:
1. $Var(aX + b) = a^2 Var(X)$
2. $Var(X + Y) = Var(X) + Var(Y)$ (如果 $X$ 和 $Y$ 独立)

### 4.3 Go语言实现

```go
package moments

import (
    "math"
)

// Moments 计算矩
type Moments struct {
    Mean     float64
    Variance float64
    Skewness float64
    Kurtosis float64
}

// CalculateMoments 计算样本矩
func CalculateMoments(samples []float64) *Moments {
    n := float64(len(samples))
    if n == 0 {
        return &Moments{}
    }
    
    // 计算均值
    mean := 0.0
    for _, x := range samples {
        mean += x
    }
    mean /= n
    
    // 计算方差
    variance := 0.0
    for _, x := range samples {
        variance += math.Pow(x-mean, 2)
    }
    variance /= n
    
    // 计算偏度
    skewness := 0.0
    for _, x := range samples {
        skewness += math.Pow((x-mean)/math.Sqrt(variance), 3)
    }
    skewness /= n
    
    // 计算峰度
    kurtosis := 0.0
    for _, x := range samples {
        kurtosis += math.Pow((x-mean)/math.Sqrt(variance), 4)
    }
    kurtosis /= n
    
    return &Moments{
        Mean:     mean,
        Variance: variance,
        Skewness: skewness,
        Kurtosis: kurtosis,
    }
}

// Covariance 计算协方差
func Covariance(x, y []float64) float64 {
    n := float64(len(x))
    if n != float64(len(y)) || n == 0 {
        return 0
    }
    
    meanX := 0.0
    meanY := 0.0
    for i := range x {
        meanX += x[i]
        meanY += y[i]
    }
    meanX /= n
    meanY /= n
    
    covariance := 0.0
    for i := range x {
        covariance += (x[i] - meanX) * (y[i] - meanY)
    }
    covariance /= n
    
    return covariance
}

// Correlation 计算相关系数
func Correlation(x, y []float64) float64 {
    covariance := Covariance(x, y)
    
    // 计算标准差
    momentsX := CalculateMoments(x)
    momentsY := CalculateMoments(y)
    
    stdX := math.Sqrt(momentsX.Variance)
    stdY := math.Sqrt(momentsY.Variance)
    
    if stdX == 0 || stdY == 0 {
        return 0
    }
    
    return covariance / (stdX * stdY)
}
```

---

## 5. 大数定律与中心极限定理 (Law of Large Numbers and Central Limit Theorem)

### 5.1 大数定律

**定理 5.1.1** (弱大数定律)
设 $X_1, X_2, ...$ 是独立同分布的随机变量，期望为 $\mu$，则：
$\lim_{n \to \infty} P(|\frac{1}{n}\sum_{i=1}^n X_i - \mu| > \epsilon) = 0$

**定理 5.1.2** (强大数定律)
设 $X_1, X_2, ...$ 是独立同分布的随机变量，期望为 $\mu$，则：
$P(\lim_{n \to \infty} \frac{1}{n}\sum_{i=1}^n X_i = \mu) = 1$

### 5.2 中心极限定理

**定理 5.2.1** (中心极限定理)
设 $X_1, X_2, ...$ 是独立同分布的随机变量，期望为 $\mu$，方差为 $\sigma^2$，则：
$\frac{\sum_{i=1}^n X_i - n\mu}{\sqrt{n}\sigma} \xrightarrow{d} N(0,1)$

### 5.3 Go语言实现

```go
package limit_theorems

import (
    "math"
    "math/rand"
)

// LawOfLargeNumbers 大数定律演示
func LawOfLargeNumbers(distribution RandomVariable, n int) []float64 {
    samples := make([]float64, n)
    runningMeans := make([]float64, n)
    
    cumulative := 0.0
    for i := 0; i < n; i++ {
        samples[i] = distribution.Sample()
        cumulative += samples[i]
        runningMeans[i] = cumulative / float64(i+1)
    }
    
    return runningMeans
}

// CentralLimitTheorem 中心极限定理演示
func CentralLimitTheorem(distribution RandomVariable, n, m int) []float64 {
    // 生成m个样本，每个样本包含n个随机变量
    sampleMeans := make([]float64, m)
    
    for i := 0; i < m; i++ {
        sum := 0.0
        for j := 0; j < n; j++ {
            sum += distribution.Sample()
        }
        sampleMeans[i] = sum / float64(n)
    }
    
    return sampleMeans
}

// MonteCarloIntegration 蒙特卡洛积分
func MonteCarloIntegration(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + rand.Float64()*(b-a)
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}

// Example: 计算π的蒙特卡洛方法
func MonteCarloPi(n int) float64 {
    inside := 0
    for i := 0; i < n; i++ {
        x := 2*rand.Float64() - 1
        y := 2*rand.Float64() - 1
        if x*x+y*y <= 1 {
            inside++
        }
    }
    return 4 * float64(inside) / float64(n)
}

// Bootstrap 自助法
func Bootstrap(samples []float64, n int, statistic func([]float64) float64) []float64 {
    bootstrapStats := make([]float64, n)
    
    for i := 0; i < n; i++ {
        // 有放回抽样
        bootstrapSample := make([]float64, len(samples))
        for j := range bootstrapSample {
            bootstrapSample[j] = samples[rand.Intn(len(samples))]
        }
        bootstrapStats[i] = statistic(bootstrapSample)
    }
    
    return bootstrapStats
}

// ConfidenceInterval 计算置信区间
func ConfidenceInterval(stats []float64, confidence float64) (float64, float64) {
    // 排序
    sorted := make([]float64, len(stats))
    copy(sorted, stats)
    sort.Float64s(sorted)
    
    alpha := 1 - confidence
    lowerIndex := int(alpha/2 * float64(len(sorted)))
    upperIndex := int((1-alpha/2) * float64(len(sorted)))
    
    return sorted[lowerIndex], sorted[upperIndex]
}
```

---

## 6. Go语言中的概率实现

### 6.1 随机数生成器

```go
package random

import (
    "crypto/rand"
    "encoding/binary"
    "math"
    "math/big"
)

// SecureRandom 安全随机数生成器
type SecureRandom struct{}

func NewSecureRandom() *SecureRandom {
    return &SecureRandom{}
}

// Int64 生成安全的int64随机数
func (sr *SecureRandom) Int64() int64 {
    var b [8]byte
    rand.Read(b[:])
    return int64(binary.BigEndian.Uint64(b[:]))
}

// Float64 生成安全的float64随机数
func (sr *SecureRandom) Float64() float64 {
    // 生成[0,1)范围内的随机数
    n, _ := rand.Int(rand.Reader, big.NewInt(1<<53))
    return float64(n.Int64()) / (1 << 53)
}

// Normal 生成正态分布随机数
func (sr *SecureRandom) Normal(mu, sigma float64) float64 {
    // Box-Muller变换
    u1 := sr.Float64()
    u2 := sr.Float64()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return mu + sigma*z0
}

// Exponential 生成指数分布随机数
func (sr *SecureRandom) Exponential(lambda float64) float64 {
    return -math.Log(sr.Float64()) / lambda
}

// Multinomial 多项分布采样
func (sr *SecureRandom) Multinomial(n int, probabilities []float64) []int {
    result := make([]int, len(probabilities))
    
    for i := 0; i < n; i++ {
        r := sr.Float64()
        cumulative := 0.0
        
        for j, prob := range probabilities {
            cumulative += prob
            if r <= cumulative {
                result[j]++
                break
            }
        }
    }
    
    return result
}
```

### 6.2 概率统计工具

```go
package statistics

import (
    "math"
    "sort"
)

// DescriptiveStatistics 描述性统计
type DescriptiveStatistics struct {
    Count     int
    Mean      float64
    Median    float64
    Mode      []float64
    Variance  float64
    StdDev    float64
    Min       float64
    Max       float64
    Quartiles [3]float64
}

// CalculateDescriptiveStatistics 计算描述性统计
func CalculateDescriptiveStatistics(data []float64) *DescriptiveStatistics {
    if len(data) == 0 {
        return &DescriptiveStatistics{}
    }
    
    // 排序
    sorted := make([]float64, len(data))
    copy(sorted, data)
    sort.Float64s(sorted)
    
    n := float64(len(sorted))
    
    // 基本统计量
    mean := 0.0
    for _, x := range sorted {
        mean += x
    }
    mean /= n
    
    variance := 0.0
    for _, x := range sorted {
        variance += math.Pow(x-mean, 2)
    }
    variance /= n
    
    // 中位数
    median := 0.0
    if int(n)%2 == 0 {
        median = (sorted[int(n)/2-1] + sorted[int(n)/2]) / 2
    } else {
        median = sorted[int(n)/2]
    }
    
    // 众数
    mode := calculateMode(sorted)
    
    // 四分位数
    quartiles := calculateQuartiles(sorted)
    
    return &DescriptiveStatistics{
        Count:     len(sorted),
        Mean:      mean,
        Median:    median,
        Mode:      mode,
        Variance:  variance,
        StdDev:    math.Sqrt(variance),
        Min:       sorted[0],
        Max:       sorted[len(sorted)-1],
        Quartiles: quartiles,
    }
}

// calculateMode 计算众数
func calculateMode(data []float64) []float64 {
    frequency := make(map[float64]int)
    maxFreq := 0
    
    for _, x := range data {
        frequency[x]++
        if frequency[x] > maxFreq {
            maxFreq = frequency[x]
        }
    }
    
    var modes []float64
    for value, freq := range frequency {
        if freq == maxFreq {
            modes = append(modes, value)
        }
    }
    
    return modes
}

// calculateQuartiles 计算四分位数
func calculateQuartiles(data []float64) [3]float64 {
    n := len(data)
    quartiles := [3]float64{}
    
    // Q1 (25%)
    q1Index := (n - 1) / 4
    quartiles[0] = data[q1Index]
    
    // Q2 (50%) - 中位数
    q2Index := (n - 1) / 2
    quartiles[1] = data[q2Index]
    
    // Q3 (75%)
    q3Index := 3 * (n - 1) / 4
    quartiles[2] = data[q3Index]
    
    return quartiles
}

// HypothesisTest 假设检验
type HypothesisTest struct {
    TestStatistic float64
    PValue        float64
    RejectNull    bool
}

// TTest 独立样本t检验
func TTest(sample1, sample2 []float64, alpha float64) *HypothesisTest {
    n1 := float64(len(sample1))
    n2 := float64(len(sample2))
    
    if n1 == 0 || n2 == 0 {
        return &HypothesisTest{}
    }
    
    // 计算样本均值和方差
    mean1 := 0.0
    for _, x := range sample1 {
        mean1 += x
    }
    mean1 /= n1
    
    mean2 := 0.0
    for _, x := range sample2 {
        mean2 += x
    }
    mean2 /= n2
    
    var1 := 0.0
    for _, x := range sample1 {
        var1 += math.Pow(x-mean1, 2)
    }
    var1 /= (n1 - 1)
    
    var2 := 0.0
    for _, x := range sample2 {
        var2 += math.Pow(x-mean2, 2)
    }
    var2 /= (n2 - 1)
    
    // 计算t统计量
    pooledVar := ((n1-1)*var1 + (n2-1)*var2) / (n1 + n2 - 2)
    tStat := (mean1 - mean2) / math.Sqrt(pooledVar*(1/n1+1/n2))
    
    // 计算p值（简化版本）
    df := n1 + n2 - 2
    pValue := 2 * (1 - tDistributionCDF(math.Abs(tStat), df))
    
    return &HypothesisTest{
        TestStatistic: tStat,
        PValue:        pValue,
        RejectNull:    pValue < alpha,
    }
}

// tDistributionCDF t分布累积分布函数（近似）
func tDistributionCDF(x float64, df float64) float64 {
    // 使用正态分布近似
    if df > 30 {
        return normalCDF(x)
    }
    
    // 简化的t分布CDF近似
    z := x / math.Sqrt(df/(df-2))
    return normalCDF(z)
}

// normalCDF 标准正态分布CDF
func normalCDF(x float64) float64 {
    return 0.5 * (1 + erf(x/math.Sqrt(2)))
}

// erf 误差函数
func erf(x float64) float64 {
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
```

---

## 总结

本文档介绍了概率论的基本概念和实现：

1. **概率空间** - 样本空间、事件、概率测度
2. **随机变量** - 离散和连续随机变量、PMF、PDF
3. **分布函数** - 常见分布及其性质
4. **期望与方差** - 矩的计算和性质
5. **极限定理** - 大数定律、中心极限定理
6. **Go语言实现** - 随机数生成、统计工具

这些概念为机器学习、统计学、金融建模等领域提供了理论基础。

---

**相关链接**:
- [01-集合论 (Set Theory)](01-Set-Theory.md)
- [02-逻辑学 (Logic)](02-Logic.md)
- [03-图论 (Graph Theory)](03-Graph-Theory.md)
- [02-形式化验证 (Formal Verification)](../02-Formal-Verification/README.md) 