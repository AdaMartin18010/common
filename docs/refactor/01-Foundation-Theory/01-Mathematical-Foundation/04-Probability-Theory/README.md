# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 概率空间](#11-概率空间)
    - [1.2 随机变量](#12-随机变量)
    - [1.3 概率分布](#13-概率分布)
      - [1.3.1 离散分布](#131-离散分布)
      - [1.3.2 连续分布](#132-连续分布)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 测度论基础](#21-测度论基础)
    - [2.2 概率公理](#22-概率公理)
    - [2.3 条件概率](#23-条件概率)
  - [3. 重要定理](#3-重要定理)
    - [3.1 大数定律](#31-大数定律)
    - [3.2 中心极限定理](#32-中心极限定理)
    - [3.3 贝叶斯定理](#33-贝叶斯定理)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 随机数生成](#41-随机数生成)
    - [4.2 概率分布](#42-概率分布)
    - [4.3 统计计算](#43-统计计算)
  - [5. 应用场景](#5-应用场景)
    - [5.1 机器学习](#51-机器学习)
    - [5.2 金融建模](#52-金融建模)
    - [5.3 质量控制](#53-质量控制)
    - [5.4 生物统计学](#54-生物统计学)
  - [6. 总结](#6-总结)

## 1. 基础概念

### 1.1 概率空间

**定义 1.1** (概率空间): 概率空间 $(\Omega, \mathcal{F}, P)$ 由以下三个部分组成：

- $\Omega$: 样本空间，包含所有可能的结果
- $\mathcal{F}$: 事件域，是 $\Omega$ 的子集的 $\sigma$-代数
- $P$: 概率测度，满足概率公理

**定义 1.2** ($\sigma$-代数): 集合族 $\mathcal{F}$ 是 $\sigma$-代数，如果：

1. $\Omega \in \mathcal{F}$
2. 如果 $A \in \mathcal{F}$，则 $A^c \in \mathcal{F}$
3. 如果 $A_1, A_2, \ldots \in \mathcal{F}$，则 $\bigcup_{i=1}^{\infty} A_i \in \mathcal{F}$

### 1.2 随机变量

**定义 1.3** (随机变量): 随机变量 $X$ 是从样本空间 $\Omega$ 到实数集 $\mathbb{R}$ 的可测函数，即对于任意 $a \in \mathbb{R}$，$\{X \leq a\} \in \mathcal{F}$。

**定义 1.4** (分布函数): 随机变量 $X$ 的分布函数 $F_X(x)$ 定义为：
$$F_X(x) = P(X \leq x)$$

### 1.3 概率分布

#### 1.3.1 离散分布

**定义 1.5** (概率质量函数): 离散随机变量 $X$ 的概率质量函数 $p_X(x)$ 定义为：
$$p_X(x) = P(X = x)$$

#### 1.3.2 连续分布

**定义 1.6** (概率密度函数): 连续随机变量 $X$ 的概率密度函数 $f_X(x)$ 满足：
$$F_X(x) = \int_{-\infty}^{x} f_X(t) dt$$

## 2. 形式化定义

### 2.1 测度论基础

**定义 2.1** (测度): 测度 $\mu$ 是定义在可测空间 $(\Omega, \mathcal{F})$ 上的函数，满足：

1. 非负性: $\mu(A) \geq 0$ 对所有 $A \in \mathcal{F}$
2. 空集测度: $\mu(\emptyset) = 0$
3. 可数可加性: 对于互不相交的可数集族 $\{A_i\}$，
   $$\mu\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} \mu(A_i)$$

**定理 2.1** (测度的单调性): 如果 $A \subseteq B$，则 $\mu(A) \leq \mu(B)$。

**证明**: $B = A \cup (B \setminus A)$，由可数可加性得：
$$\mu(B) = \mu(A) + \mu(B \setminus A) \geq \mu(A)$$

### 2.2 概率公理

**公理 2.1** (Kolmogorov概率公理):

1. 非负性: $P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
2. 规范性: $P(\Omega) = 1$
3. 可数可加性: 对于互不相交的可数集族 $\{A_i\}$，
   $$P\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} P(A_i)$$

**定理 2.2** (概率的基本性质):

1. $P(A^c) = 1 - P(A)$
2. $P(A \cup B) = P(A) + P(B) - P(A \cap B)$
3. 如果 $A \subseteq B$，则 $P(A) \leq P(B)$

### 2.3 条件概率

**定义 2.2** (条件概率): 给定事件 $B$ 的条件下事件 $A$ 的概率定义为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 2.3** (乘法公式): $P(A \cap B) = P(A|B)P(B) = P(B|A)P(A)$

**定理 2.4** (全概率公式): 如果 $\{B_i\}$ 是样本空间的一个划分，则：
$$P(A) = \sum_{i=1}^{\infty} P(A|B_i)P(B_i)$$

## 3. 重要定理

### 3.1 大数定律

**定理 3.1** (弱大数定律): 设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，则：
$$\frac{1}{n}\sum_{i=1}^{n} X_i \xrightarrow{P} \mu$$

**定理 3.2** (强大数定律): 在相同条件下：
$$\frac{1}{n}\sum_{i=1}^{n} X_i \xrightarrow{a.s.} \mu$$

### 3.2 中心极限定理

**定理 3.3** (中心极限定理): 设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，方差为 $\sigma^2$，则：
$$\frac{\sum_{i=1}^{n} X_i - n\mu}{\sigma\sqrt{n}} \xrightarrow{d} N(0,1)$$

### 3.3 贝叶斯定理

**定理 3.4** (贝叶斯定理): 对于事件 $A$ 和 $B$：
$$P(A|B) = \frac{P(B|A)P(A)}{P(B)}$$

## 4. Go语言实现

### 4.1 随机数生成

```go
package probability

import (
    "crypto/rand"
    "math"
    "math/big"
)

// 随机数生成器接口
type RandomGenerator interface {
    Next() float64
    NextInt(n int) int
    NextNormal() float64
    NextExponential(lambda float64) float64
}

// 加密安全的随机数生成器
type CryptoRandomGenerator struct{}

func NewCryptoRandomGenerator() *CryptoRandomGenerator {
    return &CryptoRandomGenerator{}
}

// 生成 [0,1) 之间的随机数
func (r *CryptoRandomGenerator) Next() float64 {
    n, _ := rand.Int(rand.Reader, big.NewInt(1<<53))
    return float64(n.Int64()) / (1 << 53)
}

// 生成 [0,n) 之间的随机整数
func (r *CryptoRandomGenerator) NextInt(n int) int {
    if n <= 0 {
        return 0
    }
    bigN := big.NewInt(int64(n))
    result, _ := rand.Int(rand.Reader, bigN)
    return int(result.Int64())
}

// Box-Muller变换生成正态分布随机数
func (r *CryptoRandomGenerator) NextNormal() float64 {
    u1 := r.Next()
    u2 := r.Next()
    
    // Box-Muller变换
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return z0
}

// 生成指数分布随机数
func (r *CryptoRandomGenerator) NextExponential(lambda float64) float64 {
    u := r.Next()
    return -math.Log(1-u) / lambda
}
```

### 4.2 概率分布

```go
// 概率分布接口
type Distribution interface {
    PDF(x float64) float64
    CDF(x float64) float64
    Mean() float64
    Variance() float64
    Sample(rng RandomGenerator) float64
}

// 正态分布
type NormalDistribution struct {
    mu    float64 // 均值
    sigma float64 // 标准差
}

func NewNormalDistribution(mu, sigma float64) *NormalDistribution {
    return &NormalDistribution{
        mu:    mu,
        sigma: sigma,
    }
}

// 概率密度函数
func (n *NormalDistribution) PDF(x float64) float64 {
    exponent := -0.5 * math.Pow((x-n.mu)/n.sigma, 2)
    return math.Exp(exponent) / (n.sigma * math.Sqrt(2*math.Pi))
}

// 累积分布函数（使用误差函数近似）
func (n *NormalDistribution) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + math.Erf(z/math.Sqrt(2)))
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
func (n *NormalDistribution) Sample(rng RandomGenerator) float64 {
    return n.mu + n.sigma*rng.NextNormal()
}

// 二项分布
type BinomialDistribution struct {
    n int     // 试验次数
    p float64 // 成功概率
}

func NewBinomialDistribution(n int, p float64) *BinomialDistribution {
    return &BinomialDistribution{
        n: n,
        p: p,
    }
}

// 概率质量函数
func (b *BinomialDistribution) PMF(k int) float64 {
    if k < 0 || k > b.n {
        return 0
    }
    
    // 使用对数避免数值溢出
    logPmf := b.logBinomial(b.n, k) + float64(k)*math.Log(b.p) + 
              float64(b.n-k)*math.Log(1-b.p)
    return math.Exp(logPmf)
}

// 累积分布函数
func (b *BinomialDistribution) CDF(k int) float64 {
    if k < 0 {
        return 0
    }
    if k >= b.n {
        return 1
    }
    
    sum := 0.0
    for i := 0; i <= k; i++ {
        sum += b.PMF(i)
    }
    return sum
}

// 对数二项式系数
func (b *BinomialDistribution) logBinomial(n, k int) float64 {
    if k == 0 || k == n {
        return 0
    }
    
    // 使用Stirling公式近似
    return b.logFactorial(n) - b.logFactorial(k) - b.logFactorial(n-k)
}

// 对数阶乘
func (b *BinomialDistribution) logFactorial(n int) float64 {
    if n <= 1 {
        return 0
    }
    
    sum := 0.0
    for i := 2; i <= n; i++ {
        sum += math.Log(float64(i))
    }
    return sum
}

// 均值
func (b *BinomialDistribution) Mean() float64 {
    return float64(b.n) * b.p
}

// 方差
func (b *BinomialDistribution) Variance() float64 {
    return float64(b.n) * b.p * (1 - b.p)
}

// 采样（使用逆变换法）
func (b *BinomialDistribution) Sample(rng RandomGenerator) float64 {
    u := rng.Next()
    cdf := 0.0
    
    for k := 0; k <= b.n; k++ {
        cdf += b.PMF(k)
        if u <= cdf {
            return float64(k)
        }
    }
    
    return float64(b.n)
}
```

### 4.3 统计计算

```go
// 统计工具
type Statistics struct{}

// 计算样本均值
func (s *Statistics) Mean(data []float64) float64 {
    if len(data) == 0 {
        return 0
    }
    
    sum := 0.0
    for _, x := range data {
        sum += x
    }
    return sum / float64(len(data))
}

// 计算样本方差
func (s *Statistics) Variance(data []float64) float64 {
    if len(data) <= 1 {
        return 0
    }
    
    mean := s.Mean(data)
    sum := 0.0
    
    for _, x := range data {
        sum += math.Pow(x-mean, 2)
    }
    
    return sum / float64(len(data)-1) // 无偏估计
}

// 计算样本标准差
func (s *Statistics) StandardDeviation(data []float64) float64 {
    return math.Sqrt(s.Variance(data))
}

// 计算协方差
func (s *Statistics) Covariance(x, y []float64) float64 {
    if len(x) != len(y) || len(x) == 0 {
        return 0
    }
    
    meanX := s.Mean(x)
    meanY := s.Mean(y)
    
    sum := 0.0
    for i := 0; i < len(x); i++ {
        sum += (x[i] - meanX) * (y[i] - meanY)
    }
    
    return sum / float64(len(x)-1)
}

// 计算相关系数
func (s *Statistics) Correlation(x, y []float64) float64 {
    cov := s.Covariance(x, y)
    stdX := s.StandardDeviation(x)
    stdY := s.StandardDeviation(y)
    
    if stdX == 0 || stdY == 0 {
        return 0
    }
    
    return cov / (stdX * stdY)
}

// 置信区间计算
type ConfidenceInterval struct {
    Lower float64
    Upper float64
    Level float64 // 置信水平
}

// 计算正态分布的置信区间
func (s *Statistics) NormalConfidenceInterval(data []float64, confidence float64) ConfidenceInterval {
    n := len(data)
    if n == 0 {
        return ConfidenceInterval{}
    }
    
    mean := s.Mean(data)
    std := s.StandardDeviation(data)
    
    // 计算临界值（使用t分布近似）
    alpha := 1 - confidence
    tCritical := s.tCriticalValue(n-1, alpha/2)
    
    margin := tCritical * std / math.Sqrt(float64(n))
    
    return ConfidenceInterval{
        Lower: mean - margin,
        Upper: mean + margin,
        Level: confidence,
    }
}

// t分布临界值（简化实现）
func (s *Statistics) tCriticalValue(df int, alpha float64) float64 {
    // 这里使用正态分布近似，实际应用中应使用t分布表
    if alpha == 0.025 {
        return 1.96 // 95%置信区间的临界值
    }
    return 1.645 // 90%置信区间的临界值
}

// 假设检验
type HypothesisTest struct {
    TestStatistic float64
    PValue        float64
    RejectNull    bool
}

// 单样本t检验
func (s *Statistics) OneSampleTTest(data []float64, mu0 float64) HypothesisTest {
    n := len(data)
    if n <= 1 {
        return HypothesisTest{}
    }
    
    mean := s.Mean(data)
    std := s.StandardDeviation(data)
    
    // 计算t统计量
    t := (mean - mu0) / (std / math.Sqrt(float64(n)))
    
    // 计算p值（简化实现）
    pValue := s.tTestPValue(t, n-1)
    
    return HypothesisTest{
        TestStatistic: t,
        PValue:        pValue,
        RejectNull:    pValue < 0.05,
    }
}

// t检验p值计算（简化实现）
func (s *Statistics) tTestPValue(t float64, df int) float64 {
    // 这里使用正态分布近似
    if math.Abs(t) > 1.96 {
        return 0.05
    }
    return 0.1
}
```

## 5. 应用场景

### 5.1 机器学习

- 贝叶斯分类器
- 概率图模型
- 随机梯度下降

### 5.2 金融建模

- 风险度量
- 期权定价
- 投资组合优化

### 5.3 质量控制

- 统计过程控制
- 抽样检验
- 可靠性分析

### 5.4 生物统计学

- 临床试验分析
- 流行病学研究
- 基因组学分析

## 6. 总结

概率论作为现代数学的重要分支，在计算机科学、统计学、金融学等领域有广泛应用。通过Go语言的实现，我们可以看到：

1. **理论实现**: 概率论的基本概念可以转化为高效的代码
2. **数值计算**: 通过数值方法实现复杂的概率计算
3. **实际应用**: 概率论在多个领域都有重要应用

概率论的研究不仅有助于理解随机现象的本质，也为数据分析和决策制定提供了科学的基础。

---

**相关链接**:

- [01-集合论](../01-Set-Theory/README.md)
- [02-逻辑学](../02-Logic/README.md)
- [03-图论](../03-Graph-Theory/README.md)
- [03-设计模式](../../03-Design-Patterns/README.md)
- [02-软件架构](../../02-Software-Architecture/README.md)
