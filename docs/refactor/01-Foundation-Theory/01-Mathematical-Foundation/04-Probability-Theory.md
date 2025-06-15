# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基础定义](#1-基础定义)
    - [1.1 概率空间](#11-概率空间)
    - [1.2 随机变量](#12-随机变量)
    - [1.3 概率分布](#13-概率分布)
  - [2. 概率性质](#2-概率性质)
    - [2.1 条件概率](#21-条件概率)
    - [2.2 独立性](#22-独立性)
    - [2.3 贝叶斯定理](#23-贝叶斯定理)
  - [3. 随机过程](#3-随机过程)
    - [3.1 马尔可夫链](#31-马尔可夫链)
    - [3.2 泊松过程](#32-泊松过程)
    - [3.3 布朗运动](#33-布朗运动)
  - [4. 形式化定义](#4-形式化定义)
    - [4.1 概率测度](#41-概率测度)
    - [4.2 期望和方差](#42-期望和方差)
    - [4.3 大数定律和中心极限定理](#43-大数定律和中心极限定理)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 随机数生成](#51-随机数生成)
    - [5.2 概率分布](#52-概率分布)
    - [5.3 统计推断](#53-统计推断)
  - [6. 应用实例](#6-应用实例)
    - [6.1 蒙特卡洛方法](#61-蒙特卡洛方法)
    - [6.2 随机算法](#62-随机算法)
    - [6.3 概率模型](#63-概率模型)

## 1. 基础定义

### 1.1 概率空间

**定义 1.1 (概率空间)**
概率空间 $(\Omega, \mathcal{F}, P)$ 由以下三个部分组成：
- $\Omega$：样本空间，包含所有可能的结果
- $\mathcal{F}$：事件域，是 $\Omega$ 的子集的 $\sigma$-代数
- $P$：概率测度，满足：
  1. $P(\Omega) = 1$
  2. $P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
  3. 可列可加性：对于互不相交的事件序列 $\{A_i\}$，$P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$

**定义 1.2 (事件)**
事件是样本空间 $\Omega$ 的子集，属于事件域 $\mathcal{F}$。

**定义 1.3 (基本事件)**
基本事件是样本空间中的单个元素。

### 1.2 随机变量

**定义 1.4 (随机变量)**
随机变量 $X$ 是从概率空间 $(\Omega, \mathcal{F}, P)$ 到实数集 $\mathbb{R}$ 的可测函数，即对于任意 $a \in \mathbb{R}$，$\{\omega \in \Omega : X(\omega) \leq a\} \in \mathcal{F}$。

**定义 1.5 (离散随机变量)**
如果随机变量 $X$ 只取有限或可列个值，则称 $X$ 为离散随机变量。

**定义 1.6 (连续随机变量)**
如果随机变量 $X$ 的分布函数 $F_X(x) = P(X \leq x)$ 是连续函数，则称 $X$ 为连续随机变量。

### 1.3 概率分布

**定义 1.7 (分布函数)**
随机变量 $X$ 的分布函数定义为：
$$F_X(x) = P(X \leq x), \quad x \in \mathbb{R}$$

**定义 1.8 (概率质量函数)**
离散随机变量 $X$ 的概率质量函数定义为：
$$p_X(x) = P(X = x), \quad x \in \mathbb{R}$$

**定义 1.9 (概率密度函数)**
连续随机变量 $X$ 的概率密度函数 $f_X(x)$ 满足：
$$F_X(x) = \int_{-\infty}^x f_X(t) dt, \quad x \in \mathbb{R}$$

## 2. 概率性质

### 2.1 条件概率

**定义 2.1 (条件概率)**
对于事件 $A, B$，且 $P(B) > 0$，$A$ 在 $B$ 条件下的条件概率定义为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 2.1 (乘法公式)**
对于事件 $A_1, A_2, \ldots, A_n$，且 $P(A_1 \cap A_2 \cap \cdots \cap A_{n-1}) > 0$：
$$P(A_1 \cap A_2 \cap \cdots \cap A_n) = P(A_1) \cdot P(A_2|A_1) \cdot P(A_3|A_1 \cap A_2) \cdots P(A_n|A_1 \cap A_2 \cap \cdots \cap A_{n-1})$$

**证明**：
使用数学归纳法。对于 $n=2$，结论显然成立。假设对于 $n=k$ 成立，考虑 $n=k+1$：
$$P(A_1 \cap A_2 \cap \cdots \cap A_{k+1}) = P(A_1 \cap A_2 \cap \cdots \cap A_k) \cdot P(A_{k+1}|A_1 \cap A_2 \cap \cdots \cap A_k)$$
根据归纳假设，结论成立。

### 2.2 独立性

**定义 2.2 (事件独立性)**
事件 $A$ 和 $B$ 是独立的，当且仅当：
$$P(A \cap B) = P(A) \cdot P(B)$$

**定义 2.3 (条件独立性)**
事件 $A$ 和 $B$ 在事件 $C$ 条件下是独立的，当且仅当：
$$P(A \cap B|C) = P(A|C) \cdot P(B|C)$$

**定理 2.2 (独立性的等价条件)**
事件 $A$ 和 $B$ 是独立的，当且仅当 $P(A|B) = P(A)$ 或 $P(B|A) = P(B)$。

**证明**：
如果 $A$ 和 $B$ 独立，则：
$$P(A|B) = \frac{P(A \cap B)}{P(B)} = \frac{P(A) \cdot P(B)}{P(B)} = P(A)$$
反之，如果 $P(A|B) = P(A)$，则：
$$P(A \cap B) = P(A|B) \cdot P(B) = P(A) \cdot P(B)$$

### 2.3 贝叶斯定理

**定理 2.3 (贝叶斯定理)**
对于事件 $A, B$，且 $P(B) > 0$：
$$P(A|B) = \frac{P(B|A) \cdot P(A)}{P(B)}$$

**证明**：
根据条件概率定义：
$$P(A|B) = \frac{P(A \cap B)}{P(B)} = \frac{P(B|A) \cdot P(A)}{P(B)}$$

**定理 2.4 (全概率公式)**
如果事件 $B_1, B_2, \ldots, B_n$ 构成样本空间的一个划分，且 $P(B_i) > 0$ 对所有 $i$，则：
$$P(A) = \sum_{i=1}^n P(A|B_i) \cdot P(B_i)$$

**证明**：
由于 $\{B_i\}$ 构成划分，$A = A \cap \Omega = A \cap (\bigcup_{i=1}^n B_i) = \bigcup_{i=1}^n (A \cap B_i)$
由于 $B_i$ 互不相交，$A \cap B_i$ 也互不相交，因此：
$$P(A) = \sum_{i=1}^n P(A \cap B_i) = \sum_{i=1}^n P(A|B_i) \cdot P(B_i)$$

## 3. 随机过程

### 3.1 马尔可夫链

**定义 3.1 (马尔可夫链)**
随机过程 $\{X_n\}_{n=0}^{\infty}$ 是马尔可夫链，如果对于任意 $n \geq 0$ 和状态 $i_0, i_1, \ldots, i_{n-1}, i, j$：
$$P(X_{n+1} = j|X_0 = i_0, X_1 = i_1, \ldots, X_{n-1} = i_{n-1}, X_n = i) = P(X_{n+1} = j|X_n = i)$$

**定义 3.2 (转移概率)**
马尔可夫链的转移概率定义为：
$$P_{ij}(n) = P(X_{n+1} = j|X_n = i)$$

**定义 3.3 (转移矩阵)**
转移矩阵 $P = [P_{ij}]$ 满足：
1. $P_{ij} \geq 0$ 对所有 $i, j$
2. $\sum_j P_{ij} = 1$ 对所有 $i$

**定理 3.1 (Chapman-Kolmogorov方程)**
对于马尔可夫链，$P_{ij}^{(n+m)} = \sum_k P_{ik}^{(n)} P_{kj}^{(m)}$。

**证明**：
$$P_{ij}^{(n+m)} = P(X_{n+m} = j|X_0 = i) = \sum_k P(X_{n+m} = j, X_n = k|X_0 = i)$$
$$= \sum_k P(X_n = k|X_0 = i) \cdot P(X_{n+m} = j|X_n = k) = \sum_k P_{ik}^{(n)} P_{kj}^{(m)}$$

### 3.2 泊松过程

**定义 3.4 (泊松过程)**
计数过程 $\{N(t)\}_{t \geq 0}$ 是强度为 $\lambda$ 的泊松过程，如果：
1. $N(0) = 0$
2. 具有独立增量
3. 具有平稳增量
4. $P(N(t+h) - N(t) = 1) = \lambda h + o(h)$
5. $P(N(t+h) - N(t) \geq 2) = o(h)$

**定理 3.2 (泊松分布)**
对于泊松过程 $\{N(t)\}$，$N(t)$ 服从参数为 $\lambda t$ 的泊松分布：
$$P(N(t) = k) = \frac{(\lambda t)^k}{k!} e^{-\lambda t}, \quad k = 0, 1, 2, \ldots$$

**证明**：
通过求解微分方程或使用生成函数方法可以证明。

### 3.3 布朗运动

**定义 3.5 (布朗运动)**
随机过程 $\{B(t)\}_{t \geq 0}$ 是标准布朗运动，如果：
1. $B(0) = 0$
2. 具有独立增量
3. 具有平稳增量
4. $B(t) - B(s) \sim N(0, t-s)$ 对所有 $0 \leq s < t$

**定理 3.3 (布朗运动性质)**
布朗运动具有以下性质：
1. $E[B(t)] = 0$
2. $Var[B(t)] = t$
3. $Cov[B(s), B(t)] = \min(s, t)$

## 4. 形式化定义

### 4.1 概率测度

**定义 4.1 (概率测度的连续性)**
概率测度 $P$ 是连续的，即对于递减的事件序列 $\{A_n\}$：
$$P(\bigcap_{n=1}^{\infty} A_n) = \lim_{n \to \infty} P(A_n)$$

**定理 4.1 (概率测度的性质)**
概率测度 $P$ 满足：
1. $P(\emptyset) = 0$
2. $P(A^c) = 1 - P(A)$
3. $P(A \cup B) = P(A) + P(B) - P(A \cap B)$
4. 如果 $A \subseteq B$，则 $P(A) \leq P(B)$

**证明**：
1. 由于 $\emptyset = \emptyset \cup \emptyset$ 且 $\emptyset \cap \emptyset = \emptyset$，$P(\emptyset) = P(\emptyset) + P(\emptyset)$，因此 $P(\emptyset) = 0$
2. $1 = P(\Omega) = P(A \cup A^c) = P(A) + P(A^c)$，因此 $P(A^c) = 1 - P(A)$
3. $P(A \cup B) = P(A \cup (B \setminus A)) = P(A) + P(B \setminus A) = P(A) + P(B) - P(A \cap B)$
4. $B = A \cup (B \setminus A)$，因此 $P(B) = P(A) + P(B \setminus A) \geq P(A)$

### 4.2 期望和方差

**定义 4.2 (期望)**
离散随机变量 $X$ 的期望定义为：
$$E[X] = \sum_x x \cdot P(X = x)$$

连续随机变量 $X$ 的期望定义为：
$$E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$$

**定义 4.3 (方差)**
随机变量 $X$ 的方差定义为：
$$Var[X] = E[(X - E[X])^2] = E[X^2] - (E[X])^2$$

**定理 4.2 (期望的线性性)**
对于随机变量 $X, Y$ 和常数 $a, b$：
$$E[aX + bY] = aE[X] + bE[Y]$$

**证明**：
对于离散情况：
$$E[aX + bY] = \sum_{x,y} (ax + by) \cdot P(X = x, Y = y)$$
$$= a\sum_{x,y} x \cdot P(X = x, Y = y) + b\sum_{x,y} y \cdot P(X = x, Y = y)$$
$$= aE[X] + bE[Y]$$

**定理 4.3 (方差的线性变换)**
对于随机变量 $X$ 和常数 $a, b$：
$$Var[aX + b] = a^2 Var[X]$$

**证明**：
$$Var[aX + b] = E[(aX + b - E[aX + b])^2] = E[(aX + b - aE[X] - b)^2]$$
$$= E[(aX - aE[X])^2] = E[a^2(X - E[X])^2] = a^2 Var[X]$$

### 4.3 大数定律和中心极限定理

**定理 4.4 (弱大数定律)**
设 $\{X_n\}$ 是独立同分布的随机变量序列，$E[X_1] = \mu$，则：
$$\frac{1}{n} \sum_{i=1}^n X_i \xrightarrow{P} \mu$$

**定理 4.5 (强大数定律)**
设 $\{X_n\}$ 是独立同分布的随机变量序列，$E[X_1] = \mu$，则：
$$\frac{1}{n} \sum_{i=1}^n X_i \xrightarrow{a.s.} \mu$$

**定理 4.6 (中心极限定理)**
设 $\{X_n\}$ 是独立同分布的随机变量序列，$E[X_1] = \mu$，$Var[X_1] = \sigma^2$，则：
$$\frac{\sum_{i=1}^n X_i - n\mu}{\sqrt{n}\sigma} \xrightarrow{d} N(0, 1)$$

## 5. Go语言实现

### 5.1 随机数生成

```go
import (
    "crypto/rand"
    "math"
    "math/big"
)

// 随机数生成器接口
type RandomGenerator interface {
    NextInt() int
    NextFloat() float64
    NextIntRange(min, max int) int
    NextFloatRange(min, max float64) float64
}

// 标准随机数生成器
type StandardRandomGenerator struct {
    source *rand.Rand
}

func NewStandardRandomGenerator() *StandardRandomGenerator {
    return &StandardRandomGenerator{
        source: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (r *StandardRandomGenerator) NextInt() int {
    return r.source.Int()
}

func (r *StandardRandomGenerator) NextFloat() float64 {
    return r.source.Float64()
}

func (r *StandardRandomGenerator) NextIntRange(min, max int) int {
    return min + r.source.Intn(max-min+1)
}

func (r *StandardRandomGenerator) NextFloatRange(min, max float64) float64 {
    return min + r.source.Float64()*(max-min)
}

// 加密安全的随机数生成器
type CryptoRandomGenerator struct{}

func NewCryptoRandomGenerator() *CryptoRandomGenerator {
    return &CryptoRandomGenerator{}
}

func (r *CryptoRandomGenerator) NextInt() int {
    n, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
    return int(n.Int64())
}

func (r *CryptoRandomGenerator) NextFloat() float64 {
    n, _ := rand.Int(rand.Reader, big.NewInt(1<<53))
    return float64(n.Int64()) / (1 << 53)
}

func (r *CryptoRandomGenerator) NextIntRange(min, max int) int {
    delta := max - min + 1
    n, _ := rand.Int(rand.Reader, big.NewInt(int64(delta)))
    return min + int(n.Int64())
}

func (r *CryptoRandomGenerator) NextFloatRange(min, max float64) float64 {
    return min + r.NextFloat()*(max-min)
}
```

### 5.2 概率分布

```go
// 概率分布接口
type ProbabilityDistribution interface {
    Probability(x float64) float64
    Cumulative(x float64) float64
    Mean() float64
    Variance() float64
    Sample() float64
}

// 离散均匀分布
type DiscreteUniformDistribution struct {
    min, max int
    rng      RandomGenerator
}

func NewDiscreteUniformDistribution(min, max int, rng RandomGenerator) *DiscreteUniformDistribution {
    return &DiscreteUniformDistribution{
        min: min,
        max: max,
        rng: rng,
    }
}

func (d *DiscreteUniformDistribution) Probability(x float64) float64 {
    xi := int(x)
    if xi >= d.min && xi <= d.max {
        return 1.0 / float64(d.max-d.min+1)
    }
    return 0.0
}

func (d *DiscreteUniformDistribution) Cumulative(x float64) float64 {
    xi := int(x)
    if xi < d.min {
        return 0.0
    }
    if xi >= d.max {
        return 1.0
    }
    return float64(xi-d.min+1) / float64(d.max-d.min+1)
}

func (d *DiscreteUniformDistribution) Mean() float64 {
    return float64(d.min+d.max) / 2.0
}

func (d *DiscreteUniformDistribution) Variance() float64 {
    n := d.max - d.min + 1
    return float64(n*n-1) / 12.0
}

func (d *DiscreteUniformDistribution) Sample() float64 {
    return float64(d.rng.NextIntRange(d.min, d.max))
}

// 连续均匀分布
type ContinuousUniformDistribution struct {
    min, max float64
    rng      RandomGenerator
}

func NewContinuousUniformDistribution(min, max float64, rng RandomGenerator) *ContinuousUniformDistribution {
    return &ContinuousUniformDistribution{
        min: min,
        max: max,
        rng: rng,
    }
}

func (d *ContinuousUniformDistribution) Probability(x float64) float64 {
    if x >= d.min && x <= d.max {
        return 1.0 / (d.max - d.min)
    }
    return 0.0
}

func (d *ContinuousUniformDistribution) Cumulative(x float64) float64 {
    if x < d.min {
        return 0.0
    }
    if x >= d.max {
        return 1.0
    }
    return (x - d.min) / (d.max - d.min)
}

func (d *ContinuousUniformDistribution) Mean() float64 {
    return (d.min + d.max) / 2.0
}

func (d *ContinuousUniformDistribution) Variance() float64 {
    return (d.max - d.min) * (d.max - d.min) / 12.0
}

func (d *ContinuousUniformDistribution) Sample() float64 {
    return d.rng.NextFloatRange(d.min, d.max)
}

// 正态分布
type NormalDistribution struct {
    mean, stddev float64
    rng          RandomGenerator
}

func NewNormalDistribution(mean, stddev float64, rng RandomGenerator) *NormalDistribution {
    return &NormalDistribution{
        mean:   mean,
        stddev: stddev,
        rng:    rng,
    }
}

func (d *NormalDistribution) Probability(x float64) float64 {
    z := (x - d.mean) / d.stddev
    return math.Exp(-z*z/2.0) / (d.stddev * math.Sqrt(2*math.Pi))
}

func (d *NormalDistribution) Cumulative(x float64) float64 {
    z := (x - d.mean) / d.stddev
    return 0.5 * (1.0 + math.Erf(z/math.Sqrt(2.0)))
}

func (d *NormalDistribution) Mean() float64 {
    return d.mean
}

func (d *NormalDistribution) Variance() float64 {
    return d.stddev * d.stddev
}

func (d *NormalDistribution) Sample() float64 {
    // Box-Muller变换
    u1 := d.rng.NextFloat()
    u2 := d.rng.NextFloat()
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return d.mean + d.stddev*z0
}

// 泊松分布
type PoissonDistribution struct {
    lambda float64
    rng    RandomGenerator
}

func NewPoissonDistribution(lambda float64, rng RandomGenerator) *PoissonDistribution {
    return &PoissonDistribution{
        lambda: lambda,
        rng:    rng,
    }
}

func (d *PoissonDistribution) Probability(k float64) float64 {
    ki := int(k)
    if ki < 0 {
        return 0.0
    }
    return math.Pow(d.lambda, float64(ki)) * math.Exp(-d.lambda) / float64(factorial(ki))
}

func (d *PoissonDistribution) Cumulative(k float64) float64 {
    ki := int(k)
    if ki < 0 {
        return 0.0
    }
    sum := 0.0
    for i := 0; i <= ki; i++ {
        sum += d.Probability(float64(i))
    }
    return sum
}

func (d *PoissonDistribution) Mean() float64 {
    return d.lambda
}

func (d *PoissonDistribution) Variance() float64 {
    return d.lambda
}

func (d *PoissonDistribution) Sample() float64 {
    // 使用逆变换方法
    u := d.rng.NextFloat()
    k := 0
    p := math.Exp(-d.lambda)
    sum := p
    
    for u > sum {
        k++
        p *= d.lambda / float64(k)
        sum += p
    }
    
    return float64(k)
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```

### 5.3 统计推断

```go
// 统计推断器
type StatisticalInference struct{}

func NewStatisticalInference() *StatisticalInference {
    return &StatisticalInference{}
}

// 计算样本均值
func (si *StatisticalInference) SampleMean(data []float64) float64 {
    if len(data) == 0 {
        return 0.0
    }
    
    sum := 0.0
    for _, x := range data {
        sum += x
    }
    return sum / float64(len(data))
}

// 计算样本方差
func (si *StatisticalInference) SampleVariance(data []float64) float64 {
    if len(data) <= 1 {
        return 0.0
    }
    
    mean := si.SampleMean(data)
    sum := 0.0
    for _, x := range data {
        sum += (x - mean) * (x - mean)
    }
    return sum / float64(len(data)-1)
}

// 计算样本标准差
func (si *StatisticalInference) SampleStdDev(data []float64) float64 {
    return math.Sqrt(si.SampleVariance(data))
}

// 计算置信区间
func (si *StatisticalInference) ConfidenceInterval(data []float64, confidence float64) (float64, float64) {
    if len(data) == 0 {
        return 0.0, 0.0
    }
    
    mean := si.SampleMean(data)
    stddev := si.SampleStdDev(data)
    n := float64(len(data))
    
    // 使用t分布（对于小样本）或正态分布（对于大样本）
    var z float64
    if n > 30 {
        // 大样本，使用正态分布
        z = si.normalQuantile((1 + confidence) / 2)
    } else {
        // 小样本，使用t分布（简化版本）
        z = si.tQuantile((1+confidence)/2, int(n)-1)
    }
    
    margin := z * stddev / math.Sqrt(n)
    return mean - margin, mean + margin
}

// 正态分布分位数（简化版本）
func (si *StatisticalInference) normalQuantile(p float64) float64 {
    // 使用近似公式
    if p < 0.5 {
        return -si.normalQuantile(1 - p)
    }
    
    t := math.Sqrt(-2 * math.Log(1 - p))
    c0 := 2.515517
    c1 := 0.802853
    c2 := 0.010328
    d1 := 1.432788
    d2 := 0.189269
    d3 := 0.001308
    
    return t - (c0 + c1*t + c2*t*t) / (1 + d1*t + d2*t*t + d3*t*t*t)
}

// t分布分位数（简化版本）
func (si *StatisticalInference) tQuantile(p float64, df int) float64 {
    // 对于大自由度，t分布接近正态分布
    if df > 30 {
        return si.normalQuantile(p)
    }
    
    // 简化版本，使用近似
    z := si.normalQuantile(p)
    return z * (1 + z*z/(4*float64(df)))
}

// 假设检验
func (si *StatisticalInference) TTest(data []float64, mu0 float64, alpha float64) (bool, float64) {
    if len(data) <= 1 {
        return false, 0.0
    }
    
    mean := si.SampleMean(data)
    stddev := si.SampleStdDev(data)
    n := float64(len(data))
    
    t := (mean - mu0) / (stddev / math.Sqrt(n))
    df := int(n) - 1
    
    // 计算p值（双尾检验）
    pValue := 2 * (1 - si.tCDF(math.Abs(t), df))
    
    return pValue < alpha, pValue
}

// t分布累积分布函数（简化版本）
func (si *StatisticalInference) tCDF(t float64, df int) float64 {
    // 使用数值积分近似
    if df > 30 {
        return si.normalCDF(t)
    }
    
    // 简化版本
    x := t / math.Sqrt(float64(df))
    return 0.5 * (1 + math.Erf(x/math.Sqrt(2.0)))
}

// 正态分布累积分布函数
func (si *StatisticalInference) normalCDF(x float64) float64 {
    return 0.5 * (1 + math.Erf(x/math.Sqrt(2.0)))
}
```

## 6. 应用实例

### 6.1 蒙特卡洛方法

```go
// 蒙特卡洛积分器
type MonteCarloIntegrator struct {
    rng RandomGenerator
}

func NewMonteCarloIntegrator(rng RandomGenerator) *MonteCarloIntegrator {
    return &MonteCarloIntegrator{
        rng: rng,
    }
}

// 一维积分
func (mci *MonteCarloIntegrator) Integrate1D(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := mci.rng.NextFloatRange(a, b)
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}

// 二维积分
func (mci *MonteCarloIntegrator) Integrate2D(f func(float64, float64) float64, a, b, c, d float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := mci.rng.NextFloatRange(a, b)
        y := mci.rng.NextFloatRange(c, d)
        sum += f(x, y)
    }
    return (b - a) * (d - c) * sum / float64(n)
}

// 计算π值
func (mci *MonteCarloIntegrator) EstimatePi(n int) float64 {
    inside := 0
    for i := 0; i < n; i++ {
        x := mci.rng.NextFloatRange(-1, 1)
        y := mci.rng.NextFloatRange(-1, 1)
        if x*x+y*y <= 1 {
            inside++
        }
    }
    return 4.0 * float64(inside) / float64(n)
}

// 风险评估
func (mci *MonteCarloIntegrator) RiskAssessment(returns []float64, n int) (float64, float64) {
    rng := NewStandardRandomGenerator()
    var results []float64
    
    for i := 0; i < n; i++ {
        // 随机选择历史收益率
        idx := rng.NextIntRange(0, len(returns)-1)
        results = append(results, returns[idx])
    }
    
    mean := 0.0
    for _, r := range results {
        mean += r
    }
    mean /= float64(len(results))
    
    variance := 0.0
    for _, r := range results {
        variance += (r - mean) * (r - mean)
    }
    variance /= float64(len(results) - 1)
    
    return mean, math.Sqrt(variance)
}
```

### 6.2 随机算法

```go
// 随机快速排序
func RandomizedQuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    rng := NewStandardRandomGenerator()
    pivotIdx := rng.NextIntRange(0, len(arr)-1)
    pivot := arr[pivotIdx]
    
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
    
    result := RandomizedQuickSort(left)
    result = append(result, equal...)
    result = append(result, RandomizedQuickSort(right)...)
    
    return result
}

// 随机化最小割算法
func RandomizedMinCut(graph map[int][]int) int {
    if len(graph) <= 2 {
        return 0
    }
    
    rng := NewStandardRandomGenerator()
    vertices := make([]int, 0, len(graph))
    for v := range graph {
        vertices = append(vertices, v)
    }
    
    // 随机选择两个顶点进行收缩
    v1 := vertices[rng.NextIntRange(0, len(vertices)-1)]
    v2 := vertices[rng.NextIntRange(0, len(vertices)-1)]
    for v2 == v1 {
        v2 = vertices[rng.NextIntRange(0, len(vertices)-1)]
    }
    
    // 收缩边
    newGraph := make(map[int][]int)
    for v, neighbors := range graph {
        if v != v1 && v != v2 {
            newNeighbors := make([]int, 0)
            for _, neighbor := range neighbors {
                if neighbor == v1 || neighbor == v2 {
                    newNeighbors = append(newNeighbors, v1)
                } else {
                    newNeighbors = append(newNeighbors, neighbor)
                }
            }
            newGraph[v] = newNeighbors
        }
    }
    
    // 合并v1和v2的邻居
    v1Neighbors := make(map[int]bool)
    for _, neighbor := range graph[v1] {
        if neighbor != v2 {
            v1Neighbors[neighbor] = true
        }
    }
    for _, neighbor := range graph[v2] {
        if neighbor != v1 {
            v1Neighbors[neighbor] = true
        }
    }
    
    v1NewNeighbors := make([]int, 0, len(v1Neighbors))
    for neighbor := range v1Neighbors {
        v1NewNeighbors = append(v1NewNeighbors, neighbor)
    }
    newGraph[v1] = v1NewNeighbors
    
    return RandomizedMinCut(newGraph)
}
```

### 6.3 概率模型

```go
// 马尔可夫链模型
type MarkovChain struct {
    states     []int
    transition map[int]map[int]float64
    initial    map[int]float64
    rng        RandomGenerator
}

func NewMarkovChain(states []int, rng RandomGenerator) *MarkovChain {
    mc := &MarkovChain{
        states:     states,
        transition: make(map[int]map[int]float64),
        initial:    make(map[int]float64),
        rng:        rng,
    }
    
    // 初始化转移矩阵
    for _, state := range states {
        mc.transition[state] = make(map[int]float64)
        for _, nextState := range states {
            mc.transition[state][nextState] = 0.0
        }
    }
    
    // 初始化概率分布
    for _, state := range states {
        mc.initial[state] = 1.0 / float64(len(states))
    }
    
    return mc
}

func (mc *MarkovChain) SetTransition(from, to int, prob float64) {
    mc.transition[from][to] = prob
}

func (mc *MarkovChain) SetInitial(state int, prob float64) {
    mc.initial[state] = prob
}

func (mc *MarkovChain) Simulate(steps int) []int {
    // 选择初始状态
    current := mc.sampleFromDistribution(mc.initial)
    path := []int{current}
    
    for i := 1; i < steps; i++ {
        // 根据转移概率选择下一个状态
        current = mc.sampleFromDistribution(mc.transition[current])
        path = append(path, current)
    }
    
    return path
}

func (mc *MarkovChain) sampleFromDistribution(dist map[int]float64) int {
    u := mc.rng.NextFloat()
    cumulative := 0.0
    
    for state, prob := range dist {
        cumulative += prob
        if u <= cumulative {
            return state
        }
    }
    
    // 如果由于浮点误差没有选择任何状态，返回最后一个
    return mc.states[len(mc.states)-1]
}

func (mc *MarkovChain) StationaryDistribution() map[int]float64 {
    // 使用幂迭代方法计算平稳分布
    n := len(mc.states)
    pi := make([]float64, n)
    pi[0] = 1.0 // 初始分布
    
    for iter := 0; iter < 100; iter++ {
        newPi := make([]float64, n)
        for i, state := range mc.states {
            for j, nextState := range mc.states {
                newPi[j] += pi[i] * mc.transition[state][nextState]
            }
        }
        
        // 归一化
        sum := 0.0
        for _, p := range newPi {
            sum += p
        }
        for i := range newPi {
            newPi[i] /= sum
        }
        
        pi = newPi
    }
    
    result := make(map[int]float64)
    for i, state := range mc.states {
        result[state] = pi[i]
    }
    
    return result
}

// 隐马尔可夫模型
type HiddenMarkovModel struct {
    states     []int
    emissions  []string
    transition map[int]map[int]float64
    emission   map[int]map[string]float64
    initial    map[int]float64
    rng        RandomGenerator
}

func NewHiddenMarkovModel(states []int, emissions []string, rng RandomGenerator) *HiddenMarkovModel {
    hmm := &HiddenMarkovModel{
        states:     states,
        emissions:  emissions,
        transition: make(map[int]map[int]float64),
        emission:   make(map[int]map[string]float64),
        initial:    make(map[int]float64),
        rng:        rng,
    }
    
    // 初始化转移矩阵
    for _, state := range states {
        hmm.transition[state] = make(map[int]float64)
        for _, nextState := range states {
            hmm.transition[state][nextState] = 0.0
        }
    }
    
    // 初始化发射矩阵
    for _, state := range states {
        hmm.emission[state] = make(map[string]float64)
        for _, emission := range emissions {
            hmm.emission[state][emission] = 0.0
        }
    }
    
    // 初始化概率分布
    for _, state := range states {
        hmm.initial[state] = 1.0 / float64(len(states))
    }
    
    return hmm
}

func (hmm *HiddenMarkovModel) SetTransition(from, to int, prob float64) {
    hmm.transition[from][to] = prob
}

func (hmm *HiddenMarkovModel) SetEmission(state int, emission string, prob float64) {
    hmm.emission[state][emission] = prob
}

func (hmm *HiddenMarkovModel) SetInitial(state int, prob float64) {
    hmm.initial[state] = prob
}

func (hmm *HiddenMarkovModel) Simulate(steps int) ([]int, []string) {
    // 选择初始状态
    currentState := hmm.sampleFromDistribution(hmm.initial)
    states := []int{currentState}
    emissions := []string{hmm.sampleEmission(currentState)}
    
    for i := 1; i < steps; i++ {
        // 转移状态
        currentState = hmm.sampleFromDistribution(hmm.transition[currentState])
        states = append(states, currentState)
        
        // 发射观测
        emission := hmm.sampleEmission(currentState)
        emissions = append(emissions, emission)
    }
    
    return states, emissions
}

func (hmm *HiddenMarkovModel) sampleEmission(state int) string {
    u := hmm.rng.NextFloat()
    cumulative := 0.0
    
    for emission, prob := range hmm.emission[state] {
        cumulative += prob
        if u <= cumulative {
            return emission
        }
    }
    
    return hmm.emissions[len(hmm.emissions)-1]
}

func (hmm *HiddenMarkovModel) sampleFromDistribution(dist map[int]float64) int {
    u := hmm.rng.NextFloat()
    cumulative := 0.0
    
    for state, prob := range dist {
        cumulative += prob
        if u <= cumulative {
            return state
        }
    }
    
    return hmm.states[len(hmm.states)-1]
}
```

## 总结

概率论作为数学的重要分支，在软件工程中有着广泛的应用。本章从基础定义出发，通过形式化证明建立了概率论的理论基础，并提供了完整的Go语言实现。

主要内容包括：

1. **理论基础**：概率空间、随机变量、概率分布的定义和性质
2. **随机过程**：马尔可夫链、泊松过程、布朗运动等
3. **统计推断**：期望、方差、大数定律、中心极限定理
4. **实际应用**：蒙特卡洛方法、随机算法、概率模型等

这些内容为后续的机器学习、数据分析、随机算法设计等提供了重要的理论基础和实践指导。

---

**参考文献**：
1. Ross, S. M. "Introduction to Probability Models." Academic Press, 2014.
2. Casella, G., & Berger, R. L. "Statistical Inference." Cengage Learning, 2002.
3. Wasserman, L. "All of Statistics." Springer, 2013.
