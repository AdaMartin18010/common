# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [概述](#概述)
  - [基本概念](#基本概念)
    - [概率空间](#概率空间)
    - [随机变量](#随机变量)
    - [概率分布](#概率分布)
  - [形式化理论](#形式化理论)
    - [概率公理](#概率公理)
    - [条件概率](#条件概率)
    - [独立性](#独立性)
    - [大数定律](#大数定律)
  - [随机过程](#随机过程)
    - [马尔可夫链](#马尔可夫链)
    - [泊松过程](#泊松过程)
    - [布朗运动](#布朗运动)
  - [Go语言实现](#go语言实现)
    - [随机数生成](#随机数生成)
    - [概率分布](#概率分布)
    - [统计函数](#统计函数)
    - [随机过程模拟](#随机过程模拟)
  - [应用领域](#应用领域)
    - [机器学习](#机器学习)
    - [金融建模](#金融建模)
    - [网络分析](#网络分析)
    - [性能分析](#性能分析)
  - [相关链接](#相关链接)

## 概述

概率论是研究随机现象数量规律的数学分支，为统计学、机器学习、金融工程等领域提供理论基础。概率论通过严格的数学公理化体系，建立了处理不确定性的理论框架。

## 基本概念

### 概率空间

**定义 1 (概率空间)**: 概率空间 $(\Omega, \mathcal{F}, P)$ 由三部分组成：
- $\Omega$: 样本空间，所有可能结果的集合
- $\mathcal{F}$: 事件域，$\Omega$ 的子集的 $\sigma$-代数
- $P$: 概率测度，满足概率公理的函数

**定义 2 (事件)**: 事件是样本空间 $\Omega$ 的子集，属于事件域 $\mathcal{F}$。

**定义 3 (概率)**: 概率 $P: \mathcal{F} \rightarrow [0,1]$ 满足：
1. $P(\Omega) = 1$
2. 对于互斥事件序列 $\{A_i\}$，$P(\bigcup_i A_i) = \sum_i P(A_i)$

### 随机变量

**定义 4 (随机变量)**: 随机变量 $X: \Omega \rightarrow \mathbb{R}$ 是可测函数，满足：
$$\{\omega \in \Omega : X(\omega) \leq x\} \in \mathcal{F}, \forall x \in \mathbb{R}$$

**定义 5 (分布函数)**: 随机变量 $X$ 的分布函数 $F_X: \mathbb{R} \rightarrow [0,1]$ 定义为：
$$F_X(x) = P(X \leq x)$$

**定义 6 (概率密度函数)**: 连续随机变量 $X$ 的概率密度函数 $f_X$ 满足：
$$F_X(x) = \int_{-\infty}^x f_X(t) dt$$

### 概率分布

**定义 7 (离散分布)**: 离散随机变量的概率质量函数 $p_X$ 满足：
$$P(X = x) = p_X(x), \sum_x p_X(x) = 1$$

**定义 8 (连续分布)**: 连续随机变量的概率密度函数 $f_X$ 满足：
$$P(a \leq X \leq b) = \int_a^b f_X(x) dx$$

## 形式化理论

### 概率公理

**公理 1 (非负性)**: $P(A) \geq 0$ 对所有事件 $A \in \mathcal{F}$

**公理 2 (规范性)**: $P(\Omega) = 1$

**公理 3 (可列可加性)**: 对于互斥事件序列 $\{A_i\}_{i=1}^{\infty}$，
$$P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$$

**定理 1 (概率的基本性质)**:
1. $P(\emptyset) = 0$
2. $P(A^c) = 1 - P(A)$
3. $P(A \cup B) = P(A) + P(B) - P(A \cap B)$

**证明**: 
1. 由公理2和3，$P(\Omega) = P(\Omega \cup \emptyset) = P(\Omega) + P(\emptyset)$，因此 $P(\emptyset) = 0$
2. $1 = P(\Omega) = P(A \cup A^c) = P(A) + P(A^c)$，因此 $P(A^c) = 1 - P(A)$
3. $P(A \cup B) = P(A \cup (B \setminus A)) = P(A) + P(B \setminus A) = P(A) + P(B) - P(A \cap B)$

### 条件概率

**定义 9 (条件概率)**: 在事件 $B$ 发生的条件下，事件 $A$ 发生的概率为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}, \text{ if } P(B) > 0$$

**定理 2 (乘法公式)**: 对于事件 $A_1, A_2, \ldots, A_n$，
$$P(A_1 \cap A_2 \cap \cdots \cap A_n) = P(A_1) \cdot P(A_2|A_1) \cdot P(A_3|A_1 \cap A_2) \cdots P(A_n|A_1 \cap A_2 \cap \cdots \cap A_{n-1})$$

**定理 3 (全概率公式)**: 如果 $\{B_i\}$ 是样本空间的分割，则：
$$P(A) = \sum_i P(A|B_i) \cdot P(B_i)$$

**定理 4 (贝叶斯公式)**: 
$$P(B_i|A) = \frac{P(A|B_i) \cdot P(B_i)}{\sum_j P(A|B_j) \cdot P(B_j)}$$

### 独立性

**定义 10 (独立性)**: 事件 $A$ 和 $B$ 独立，如果：
$$P(A \cap B) = P(A) \cdot P(B)$$

**定义 11 (条件独立性)**: 在事件 $C$ 的条件下，事件 $A$ 和 $B$ 条件独立，如果：
$$P(A \cap B|C) = P(A|C) \cdot P(B|C)$$

**定理 5 (独立性的性质)**: 如果 $A$ 和 $B$ 独立，则：
1. $A$ 和 $B^c$ 独立
2. $A^c$ 和 $B$ 独立
3. $A^c$ 和 $B^c$ 独立

### 大数定律

**定理 6 (弱大数定律)**: 设 $\{X_i\}$ 是独立同分布的随机变量序列，期望为 $\mu$，则：
$$\frac{1}{n} \sum_{i=1}^n X_i \xrightarrow{P} \mu$$

**定理 7 (强大数定律)**: 在相同条件下，
$$\frac{1}{n} \sum_{i=1}^n X_i \xrightarrow{a.s.} \mu$$

**定理 8 (中心极限定理)**: 设 $\{X_i\}$ 是独立同分布的随机变量序列，期望为 $\mu$，方差为 $\sigma^2$，则：
$$\frac{\sum_{i=1}^n X_i - n\mu}{\sqrt{n}\sigma} \xrightarrow{d} N(0,1)$$

## 随机过程

### 马尔可夫链

**定义 12 (马尔可夫链)**: 随机过程 $\{X_n\}$ 是马尔可夫链，如果：
$$P(X_{n+1} = j | X_n = i, X_{n-1} = i_{n-1}, \ldots, X_0 = i_0) = P(X_{n+1} = j | X_n = i)$$

**定义 13 (转移矩阵)**: 马尔可夫链的转移矩阵 $P$ 满足：
$$P_{ij} = P(X_{n+1} = j | X_n = i)$$

**定理 9 (Chapman-Kolmogorov方程)**: 
$$P^{(n+m)}_{ij} = \sum_k P^{(n)}_{ik} P^{(m)}_{kj}$$

### 泊松过程

**定义 14 (泊松过程)**: 计数过程 $\{N(t), t \geq 0\}$ 是强度为 $\lambda$ 的泊松过程，如果：
1. $N(0) = 0$
2. 具有独立增量
3. 具有平稳增量
4. $P(N(t) = k) = \frac{(\lambda t)^k}{k!} e^{-\lambda t}$

**定理 10 (泊松过程的性质)**: 泊松过程的到达间隔时间服从指数分布：
$$P(T > t) = e^{-\lambda t}$$

### 布朗运动

**定义 15 (布朗运动)**: 随机过程 $\{B(t), t \geq 0\}$ 是标准布朗运动，如果：
1. $B(0) = 0$
2. 具有独立增量
3. 具有平稳增量
4. $B(t) - B(s) \sim N(0, t-s)$

## Go语言实现

### 随机数生成

```go
package probability

import (
    "math"
    "math/rand"
    "time"
)

// RandomGenerator 随机数生成器接口
type RandomGenerator interface {
    Next() float64
    NextInt(n int) int
    SetSeed(seed int64)
}

// StandardRandomGenerator 标准随机数生成器
type StandardRandomGenerator struct {
    rng *rand.Rand
}

// NewStandardRandomGenerator 创建标准随机数生成器
func NewStandardRandomGenerator() *StandardRandomGenerator {
    return &StandardRandomGenerator{
        rng: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// Next 生成 [0,1) 之间的随机数
func (r *StandardRandomGenerator) Next() float64 {
    return r.rng.Float64()
}

// NextInt 生成 [0,n) 之间的随机整数
func (r *StandardRandomGenerator) NextInt(n int) int {
    return r.rng.Intn(n)
}

// SetSeed 设置随机数种子
func (r *StandardRandomGenerator) SetSeed(seed int64) {
    r.rng.Seed(seed)
}

// ProbabilitySpace 概率空间
type ProbabilitySpace struct {
    SampleSpace []interface{}
    Events      map[string][]interface{}
    Probabilities map[string]float64
}

// NewProbabilitySpace 创建概率空间
func NewProbabilitySpace() *ProbabilitySpace {
    return &ProbabilitySpace{
        SampleSpace:   make([]interface{}, 0),
        Events:        make(map[string][]interface{}),
        Probabilities: make(map[string]float64),
    }
}

// AddSamplePoint 添加样本点
func (ps *ProbabilitySpace) AddSamplePoint(point interface{}) {
    ps.SampleSpace = append(ps.SampleSpace, point)
}

// AddEvent 添加事件
func (ps *ProbabilitySpace) AddEvent(name string, outcomes []interface{}) {
    ps.Events[name] = outcomes
}

// SetProbability 设置事件概率
func (ps *ProbabilitySpace) SetProbability(event string, prob float64) {
    if prob >= 0 && prob <= 1 {
        ps.Probabilities[event] = prob
    }
}

// GetProbability 获取事件概率
func (ps *ProbabilitySpace) GetProbability(event string) float64 {
    return ps.Probabilities[event]
}
```

### 概率分布

```go
// Distribution 概率分布接口
type Distribution interface {
    PDF(x float64) float64    // 概率密度函数
    CDF(x float64) float64    // 累积分布函数
    Mean() float64            // 期望
    Variance() float64        // 方差
    Sample(rng RandomGenerator) float64 // 采样
}

// UniformDistribution 均匀分布
type UniformDistribution struct {
    a, b float64 // 区间 [a,b]
}

// NewUniformDistribution 创建均匀分布
func NewUniformDistribution(a, b float64) *UniformDistribution {
    return &UniformDistribution{a: a, b: b}
}

// PDF 概率密度函数
func (u *UniformDistribution) PDF(x float64) float64 {
    if x >= u.a && x <= u.b {
        return 1.0 / (u.b - u.a)
    }
    return 0.0
}

// CDF 累积分布函数
func (u *UniformDistribution) CDF(x float64) float64 {
    if x < u.a {
        return 0.0
    } else if x >= u.b {
        return 1.0
    } else {
        return (x - u.a) / (u.b - u.a)
    }
}

// Mean 期望
func (u *UniformDistribution) Mean() float64 {
    return (u.a + u.b) / 2.0
}

// Variance 方差
func (u *UniformDistribution) Variance() float64 {
    return math.Pow(u.b-u.a, 2) / 12.0
}

// Sample 采样
func (u *UniformDistribution) Sample(rng RandomGenerator) float64 {
    return u.a + (u.b-u.a)*rng.Next()
}

// NormalDistribution 正态分布
type NormalDistribution struct {
    mu, sigma float64 // 均值，标准差
}

// NewNormalDistribution 创建正态分布
func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    return &NormalDistribution{mu: mu, sigma: sigma}
}

// PDF 概率密度函数
func (n *NormalDistribution) PDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return math.Exp(-0.5*z*z) / (n.sigma * math.Sqrt(2*math.Pi))
}

// CDF 累积分布函数（近似）
func (n *NormalDistribution) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + erf(z/math.Sqrt(2)))
}

// Mean 期望
func (n *NormalDistribution) Mean() float64 {
    return n.mu
}

// Variance 方差
func (n *NormalDistribution) Variance() float64 {
    return n.sigma * n.sigma
}

// Sample 采样（Box-Muller变换）
func (n *NormalDistribution) Sample(rng RandomGenerator) float64 {
    u1 := rng.Next()
    u2 := rng.Next()
    
    // Box-Muller变换
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.mu + n.sigma*z0
}

// erf 误差函数（近似）
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
        x = -x
    }
    
    t := 1.0 / (1.0 + p*x)
    y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)
    
    return sign * y
}

// ExponentialDistribution 指数分布
type ExponentialDistribution struct {
    lambda float64 // 参数
}

// NewExponentialDistribution 创建指数分布
func NewExponentialDistribution(lambda float64) *ExponentialDistribution {
    return &ExponentialDistribution{lambda: lambda}
}

// PDF 概率密度函数
func (e *ExponentialDistribution) PDF(x float64) float64 {
    if x >= 0 {
        return e.lambda * math.Exp(-e.lambda*x)
    }
    return 0.0
}

// CDF 累积分布函数
func (e *ExponentialDistribution) CDF(x float64) float64 {
    if x >= 0 {
        return 1.0 - math.Exp(-e.lambda*x)
    }
    return 0.0
}

// Mean 期望
func (e *ExponentialDistribution) Mean() float64 {
    return 1.0 / e.lambda
}

// Variance 方差
func (e *ExponentialDistribution) Variance() float64 {
    return 1.0 / (e.lambda * e.lambda)
}

// Sample 采样
func (e *ExponentialDistribution) Sample(rng RandomGenerator) float64 {
    return -math.Log(1-rng.Next()) / e.lambda
}
```

### 统计函数

```go
// Statistics 统计函数
type Statistics struct{}

// Mean 计算均值
func (s *Statistics) Mean(data []float64) float64 {
    if len(data) == 0 {
        return 0.0
    }
    
    sum := 0.0
    for _, x := range data {
        sum += x
    }
    return sum / float64(len(data))
}

// Variance 计算方差
func (s *Statistics) Variance(data []float64) float64 {
    if len(data) == 0 {
        return 0.0
    }
    
    mean := s.Mean(data)
    sum := 0.0
    for _, x := range data {
        sum += math.Pow(x-mean, 2)
    }
    return sum / float64(len(data))
}

// StandardDeviation 计算标准差
func (s *Statistics) StandardDeviation(data []float64) float64 {
    return math.Sqrt(s.Variance(data))
}

// Covariance 计算协方差
func (s *Statistics) Covariance(x, y []float64) float64 {
    if len(x) != len(y) || len(x) == 0 {
        return 0.0
    }
    
    meanX := s.Mean(x)
    meanY := s.Mean(y)
    
    sum := 0.0
    for i := 0; i < len(x); i++ {
        sum += (x[i] - meanX) * (y[i] - meanY)
    }
    return sum / float64(len(x))
}

// Correlation 计算相关系数
func (s *Statistics) Correlation(x, y []float64) float64 {
    if len(x) != len(y) || len(x) == 0 {
        return 0.0
    }
    
    cov := s.Covariance(x, y)
    stdX := s.StandardDeviation(x)
    stdY := s.StandardDeviation(y)
    
    if stdX == 0 || stdY == 0 {
        return 0.0
    }
    
    return cov / (stdX * stdY)
}

// ConfidenceInterval 计算置信区间
func (s *Statistics) ConfidenceInterval(data []float64, confidence float64) (float64, float64) {
    if len(data) == 0 {
        return 0.0, 0.0
    }
    
    mean := s.Mean(data)
    std := s.StandardDeviation(data)
    n := float64(len(data))
    
    // 使用正态分布近似（大样本）
    z := 1.96 // 95% 置信水平
    if confidence == 0.99 {
        z = 2.576
    } else if confidence == 0.90 {
        z = 1.645
    }
    
    margin := z * std / math.Sqrt(n)
    return mean - margin, mean + margin
}
```

### 随机过程模拟

```go
// MarkovChain 马尔可夫链
type MarkovChain struct {
    States     []int
    Transition [][]float64
    Current    int
    rng        RandomGenerator
}

// NewMarkovChain 创建马尔可夫链
func NewMarkovChain(states []int, transition [][]float64, initial int, rng RandomGenerator) *MarkovChain {
    return &MarkovChain{
        States:     states,
        Transition: transition,
        Current:    initial,
        rng:        rng,
    }
}

// NextState 转移到下一个状态
func (mc *MarkovChain) NextState() int {
    u := mc.rng.Next()
    cumulative := 0.0
    
    for nextState := 0; nextState < len(mc.States); nextState++ {
        cumulative += mc.Transition[mc.Current][nextState]
        if u <= cumulative {
            mc.Current = nextState
            return nextState
        }
    }
    
    mc.Current = len(mc.States) - 1
    return mc.Current
}

// Simulate 模拟马尔可夫链
func (mc *MarkovChain) Simulate(steps int) []int {
    path := make([]int, steps+1)
    path[0] = mc.Current
    
    for i := 1; i <= steps; i++ {
        path[i] = mc.NextState()
    }
    
    return path
}

// PoissonProcess 泊松过程
type PoissonProcess struct {
    lambda float64 // 强度
    rng    RandomGenerator
}

// NewPoissonProcess 创建泊松过程
func NewPoissonProcess(lambda float64, rng RandomGenerator) *PoissonProcess {
    return &PoissonProcess{
        lambda: lambda,
        rng:    rng,
    }
}

// NextArrival 下一个到达时间
func (pp *PoissonProcess) NextArrival() float64 {
    return -math.Log(1-pp.rng.Next()) / pp.lambda
}

// Simulate 模拟泊松过程
func (pp *PoissonProcess) Simulate(duration float64) []float64 {
    var arrivals []float64
    currentTime := 0.0
    
    for currentTime < duration {
        arrival := pp.NextArrival()
        currentTime += arrival
        if currentTime < duration {
            arrivals = append(arrivals, currentTime)
        }
    }
    
    return arrivals
}

// BrownianMotion 布朗运动
type BrownianMotion struct {
    rng RandomGenerator
}

// NewBrownianMotion 创建布朗运动
func NewBrownianMotion(rng RandomGenerator) *BrownianMotion {
    return &BrownianMotion{rng: rng}
}

// Simulate 模拟布朗运动
func (bm *BrownianMotion) Simulate(duration float64, steps int) []float64 {
    dt := duration / float64(steps)
    path := make([]float64, steps+1)
    path[0] = 0.0
    
    for i := 1; i <= steps; i++ {
        // 生成标准正态随机变量
        u1 := bm.rng.Next()
        u2 := bm.rng.Next()
        z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
        
        // 布朗运动增量
        increment := z * math.Sqrt(dt)
        path[i] = path[i-1] + increment
    }
    
    return path
}
```

## 应用领域

### 机器学习

概率论在机器学习中的应用：
- 贝叶斯分类器
- 概率图模型
- 期望最大化算法
- 变分推断

### 金融建模

概率论在金融建模中的应用：
- 期权定价
- 风险度量
- 投资组合优化
- 信用风险模型

### 网络分析

概率论在网络分析中的应用：
- 随机图模型
- 网络可靠性分析
- 流量建模
- 性能预测

### 性能分析

概率论在性能分析中的应用：
- 排队论
- 系统可靠性
- 负载均衡
- 容量规划

## 相关链接

- [01-集合论 (Set Theory)](../01-Set-Theory/README.md)
- [02-逻辑学 (Logic)](../02-Logic/README.md)
- [03-图论 (Graph Theory)](../03-Graph-Theory/README.md)
- [08-软件工程形式化 (Software Engineering Formalization)](../../08-Software-Engineering-Formalization/README.md)
- [09-编程语言理论 (Programming Language Theory)](../../09-Programming-Language-Theory/README.md)
