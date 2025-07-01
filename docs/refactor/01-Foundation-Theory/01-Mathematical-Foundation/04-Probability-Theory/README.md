# 04-概率论 (Probability Theory)

## 目录

- [04-概率论 (Probability Theory)](#04-概率论-probability-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 概率空间](#11-概率空间)
    - [1.2 随机变量](#12-随机变量)
    - [1.3 概率分布](#13-概率分布)
  - [2. 概率公理](#2-概率公理)
    - [2.1 Kolmogorov公理](#21-kolmogorov公理)
    - [2.2 基本性质](#22-基本性质)
  - [3. 条件概率与独立性](#3-条件概率与独立性)
    - [3.1 条件概率](#31-条件概率)
    - [3.2 独立性](#32-独立性)
    - [3.3 贝叶斯定理](#33-贝叶斯定理)
  - [4. 随机变量](#4-随机变量)
    - [4.1 离散随机变量](#41-离散随机变量)
    - [4.2 连续随机变量](#42-连续随机变量)
    - [4.3 期望与方差](#43-期望与方差)
  - [5. 大数定律与中心极限定理](#5-大数定律与中心极限定理)
    - [5.1 大数定律](#51-大数定律)
    - [5.2 中心极限定理](#52-中心极限定理)
  - [6. Go语言实现](#6-go语言实现)
    - [6.1 基础概率结构](#61-基础概率结构)
    - [6.2 正态分布实现](#62-正态分布实现)
    - [6.3 统计函数](#63-统计函数)
    - [6.4 蒙特卡洛方法](#64-蒙特卡洛方法)
    - [6.5 使用示例](#65-使用示例)
  - [总结](#总结)
    - [关键要点](#关键要点)
    - [进一步研究方向](#进一步研究方向)

## 1. 基础概念

### 1.1 概率空间

**定义 1.1**: 概率空间
概率空间是一个三元组 $(\Omega, \mathcal{F}, P)$，其中：

- $\Omega$ 是样本空间（sample space），包含所有可能的结果
- $\mathcal{F}$ 是事件域（event space），是 $\Omega$ 的子集的 $\sigma$-代数
- $P$ 是概率测度（probability measure），$P: \mathcal{F} \rightarrow [0,1]$

**定义 1.2**: $\sigma$-代数
集合 $\mathcal{F}$ 是 $\sigma$-代数，如果满足：

1. $\Omega \in \mathcal{F}$
2. 如果 $A \in \mathcal{F}$，则 $A^c \in \mathcal{F}$
3. 如果 $A_1, A_2, \ldots \in \mathcal{F}$，则 $\bigcup_{i=1}^{\infty} A_i \in \mathcal{F}$

### 1.2 随机变量

**定义 1.3**: 随机变量
随机变量 $X$ 是从概率空间 $(\Omega, \mathcal{F}, P)$ 到实数集 $\mathbb{R}$ 的可测函数：
$$X: \Omega \rightarrow \mathbb{R}$$

**定义 1.4**: 分布函数
随机变量 $X$ 的分布函数 $F_X: \mathbb{R} \rightarrow [0,1]$ 定义为：
$$F_X(x) = P(X \leq x)$$

### 1.3 概率分布

**定义 1.5**: 离散分布
如果随机变量 $X$ 只取有限或可数无限个值，则称 $X$ 为离散随机变量。

**定义 1.6**: 连续分布
如果随机变量 $X$ 的分布函数 $F_X$ 是连续的，则称 $X$ 为连续随机变量。

## 2. 概率公理

### 2.1 Kolmogorov公理

**公理 2.1**: Kolmogorov概率公理
概率测度 $P$ 满足以下公理：

1. **非负性**: 对于所有 $A \in \mathcal{F}$，$P(A) \geq 0$
2. **规范性**: $P(\Omega) = 1$
3. **可列可加性**: 对于互不相容的事件 $A_1, A_2, \ldots$，
   $$P\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} P(A_i)$$

### 2.2 基本性质

**定理 2.1**: 概率的基本性质
对于任意事件 $A, B \in \mathcal{F}$：

1. $P(\emptyset) = 0$
2. $P(A^c) = 1 - P(A)$
3. 如果 $A \subseteq B$，则 $P(A) \leq P(B)$
4. $P(A \cup B) = P(A) + P(B) - P(A \cap B)$

**证明**:

1. 由可列可加性，$P(\emptyset) = P(\emptyset \cup \emptyset \cup \cdots) = P(\emptyset) + P(\emptyset) + \cdots$，因此 $P(\emptyset) = 0$
2. $1 = P(\Omega) = P(A \cup A^c) = P(A) + P(A^c)$，因此 $P(A^c) = 1 - P(A)$
3. 如果 $A \subseteq B$，则 $B = A \cup (B \setminus A)$，因此 $P(B) = P(A) + P(B \setminus A) \geq P(A)$
4. $A \cup B = A \cup (B \setminus A)$，且 $A$ 和 $B \setminus A$ 互不相容，因此 $P(A \cup B) = P(A) + P(B \setminus A) = P(A) + P(B) - P(A \cap B)$

## 3. 条件概率与独立性

### 3.1 条件概率

**定义 3.1**: 条件概率
对于事件 $A, B \in \mathcal{F}$，且 $P(B) > 0$，$A$ 在 $B$ 条件下的条件概率定义为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 3.1**: 乘法公式
对于事件 $A_1, A_2, \ldots, A_n$，如果 $P(A_1 \cap A_2 \cap \cdots \cap A_{n-1}) > 0$，则：
$$P(A_1 \cap A_2 \cap \cdots \cap A_n) = P(A_1) \cdot P(A_2|A_1) \cdot P(A_3|A_1 \cap A_2) \cdots P(A_n|A_1 \cap A_2 \cap \cdots \cap A_{n-1})$$

### 3.2 独立性

**定义 3.2**: 独立性
事件 $A$ 和 $B$ 是独立的，如果：
$$P(A \cap B) = P(A) \cdot P(B)$$

**定义 3.3**: 条件独立性
事件 $A$ 和 $B$ 在事件 $C$ 条件下是独立的，如果：
$$P(A \cap B|C) = P(A|C) \cdot P(B|C)$$

### 3.3 贝叶斯定理

**定理 3.2**: 贝叶斯定理
对于事件 $A, B \in \mathcal{F}$，且 $P(A) > 0, P(B) > 0$：
$$P(A|B) = \frac{P(B|A) \cdot P(A)}{P(B)}$$

**证明**:
由条件概率定义，$P(A|B) = \frac{P(A \cap B)}{P(B)}$ 和 $P(B|A) = \frac{P(A \cap B)}{P(A)}$

因此 $P(A \cap B) = P(B|A) \cdot P(A)$
代入第一个等式得到贝叶斯定理。

## 4. 随机变量

### 4.1 离散随机变量

**定义 4.1**: 概率质量函数
离散随机变量 $X$ 的概率质量函数 $p_X: \mathbb{R} \rightarrow [0,1]$ 定义为：
$$p_X(x) = P(X = x)$$

**性质**:

1. $p_X(x) \geq 0$ 对于所有 $x \in \mathbb{R}$
2. $\sum_{x \in \text{range}(X)} p_X(x) = 1$

**常见离散分布**:

1. **伯努利分布** $X \sim \text{Bernoulli}(p)$
   $$p_X(x) = \begin{cases}
   p & \text{if } x = 1 \\
   1-p & \text{if } x = 0 \\
   0 & \text{otherwise}
   \end{cases}$$

2. **二项分布** $X \sim \text{Binomial}(n, p)$
   $$p_X(x) = \binom{n}{x} p^x (1-p)^{n-x}, \quad x = 0, 1, \ldots, n$$

3. **泊松分布** $X \sim \text{Poisson}(\lambda)$
   $$p_X(x) = \frac{e^{-\lambda} \lambda^x}{x!}, \quad x = 0, 1, 2, \ldots$$

### 4.2 连续随机变量

**定义 4.2**: 概率密度函数
连续随机变量 $X$ 的概率密度函数 $f_X: \mathbb{R} \rightarrow [0, \infty)$ 满足：
$$F_X(x) = \int_{-\infty}^x f_X(t) dt$$

**性质**:

1. $f_X(x) \geq 0$ 对于所有 $x \in \mathbb{R}$
2. $\int_{-\infty}^{\infty} f_X(x) dx = 1$

**常见连续分布**:

1. **均匀分布** $X \sim \text{Uniform}(a, b)$
   $$f_X(x) = \begin{cases}
   \frac{1}{b-a} & \text{if } a \leq x \leq b \\
   0 & \text{otherwise}
   \end{cases}$$

2. **正态分布** $X \sim \text{Normal}(\mu, \sigma^2)$
   $$f_X(x) = \frac{1}{\sqrt{2\pi\sigma^2}} e^{-\frac{(x-\mu)^2}{2\sigma^2}}$$

3. **指数分布** $X \sim \text{Exponential}(\lambda)$
   $$f_X(x) = \begin{cases}
   \lambda e^{-\lambda x} & \text{if } x \geq 0 \\
   0 & \text{otherwise}
   \end{cases}$$

### 4.3 期望与方差

**定义 4.3**: 期望
随机变量 $X$ 的期望定义为：

- 离散情况：$E[X] = \sum_{x} x \cdot p_X(x)$
- 连续情况：$E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$

**定义 4.4**: 方差
随机变量 $X$ 的方差定义为：
$$\text{Var}(X) = E[(X - E[X])^2] = E[X^2] - (E[X])^2$$

**定理 4.1**: 期望的线性性质
对于随机变量 $X, Y$ 和常数 $a, b$：
$$E[aX + bY] = aE[X] + bE[Y]$$

**定理 4.2**: 方差的线性性质
对于独立随机变量 $X, Y$ 和常数 $a, b$：
$$\text{Var}(aX + bY) = a^2\text{Var}(X) + b^2\text{Var}(Y)$$

## 5. 大数定律与中心极限定理

### 5.1 大数定律

**定理 5.1**: 弱大数定律
设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，则对于任意 $\epsilon > 0$：
$$\lim_{n \rightarrow \infty} P\left(\left|\frac{1}{n}\sum_{i=1}^n X_i - \mu\right| > \epsilon\right) = 0$$

**定理 5.2**: 强大数定律
设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，则：
$$P\left(\lim_{n \rightarrow \infty} \frac{1}{n}\sum_{i=1}^n X_i = \mu\right) = 1$$

### 5.2 中心极限定理

**定理 5.3**: 中心极限定理
设 $X_1, X_2, \ldots$ 是独立同分布的随机变量，期望为 $\mu$，方差为 $\sigma^2$，则：
$$\frac{\sum_{i=1}^n X_i - n\mu}{\sqrt{n}\sigma} \xrightarrow{d} \text{Normal}(0, 1)$$

其中 $\xrightarrow{d}$ 表示依分布收敛。

## 6. Go语言实现

### 6.1 基础概率结构

```go
package probability

import (
    "math"
    "math/rand"
    "time"
)

// RandomVariable 表示随机变量接口
type RandomVariable interface {
    Sample() float64
    PDF(x float64) float64
    CDF(x float64) float64
    Mean() float64
    Variance() float64
}

// Bernoulli 伯努利分布
type Bernoulli struct {
    p float64 // 成功概率
}

func NewBernoulli(p float64) *Bernoulli {
    if p < 0 || p > 1 {
        panic("probability must be between 0 and 1")
    }
    return &Bernoulli{p: p}
}

func (b *Bernoulli) Sample() float64 {
    if rand.Float64() < b.p {
        return 1.0
    }
    return 0.0
}

func (b *Bernoulli) PDF(x float64) float64 {
    if x == 1 {
        return b.p
    } else if x == 0 {
        return 1 - b.p
    }
    return 0
}

func (b *Bernoulli) CDF(x float64) float64 {
    if x < 0 {
        return 0
    } else if x < 1 {
        return 1 - b.p
    }
    return 1
}

func (b *Bernoulli) Mean() float64 {
    return b.p
}

func (b *Bernoulli) Variance() float64 {
    return b.p * (1 - b.p)
}
```

### 6.2 正态分布实现

```go
// Normal 正态分布
type Normal struct {
    mu    float64 // 均值
    sigma float64 // 标准差
}

func NewNormal(mu, sigma float64) *Normal {
    if sigma <= 0 {
        panic("standard deviation must be positive")
    }
    return &Normal{mu: mu, sigma: sigma}
}

func (n *Normal) Sample() float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.mu + n.sigma*z0
}

func (n *Normal) PDF(x float64) float64 {
    exponent := -0.5 * math.Pow((x-n.mu)/n.sigma, 2)
    return (1.0 / (n.sigma * math.Sqrt(2*math.Pi))) * math.Exp(exponent)
}

func (n *Normal) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + erf(z/math.Sqrt(2)))
}

func (n *Normal) Mean() float64 {
    return n.mu
}

func (n *Normal) Variance() float64 {
    return n.sigma * n.sigma
}

// erf 误差函数近似
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
```

### 6.3 统计函数

```go
// Statistics 统计函数
type Statistics struct{}

// Mean 计算样本均值
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

// Variance 计算样本方差
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

// StandardDeviation 计算样本标准差
func (s *Statistics) StandardDeviation(data []float64) float64 {
    return math.Sqrt(s.Variance(data))
}

// ConfidenceInterval 计算置信区间
func (s *Statistics) ConfidenceInterval(data []float64, confidence float64) (float64, float64) {
    if len(data) == 0 {
        return 0, 0
    }
    
    mean := s.Mean(data)
    std := s.StandardDeviation(data)
    n := float64(len(data))
    
    // 使用正态分布近似（大样本）
    z := 1.96 // 95%置信区间对应的z值
    margin := z * std / math.Sqrt(n)
    
    return mean - margin, mean + margin
}
```

### 6.4 蒙特卡洛方法

```go
// MonteCarlo 蒙特卡洛方法
type MonteCarlo struct{}

// Integrate 使用蒙特卡洛方法计算定积分
func (mc *MonteCarlo) Integrate(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + rand.Float64()*(b-a)
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}

// EstimatePi 使用蒙特卡洛方法估计π
func (mc *MonteCarlo) EstimatePi(n int) float64 {
    inside := 0
    for i := 0; i < n; i++ {
        x := rand.Float64()
        y := rand.Float64()
        if x*x+y*y <= 1 {
            inside++
        }
    }
    return 4.0 * float64(inside) / float64(n)
}

// SimulateRandomWalk 模拟随机游走
func (mc *MonteCarlo) SimulateRandomWalk(steps int) []float64 {
    position := 0.0
    positions := make([]float64, steps+1)
    positions[0] = position
    
    for i := 1; i <= steps; i++ {
        if rand.Float64() < 0.5 {
            position += 1
        } else {
            position -= 1
        }
        positions[i] = position
    }
    
    return positions
}
```

### 6.5 使用示例

```go
func main() {
    rand.Seed(time.Now().UnixNano())
    
    // 创建随机变量
    bernoulli := NewBernoulli(0.3)
    normal := NewNormal(0, 1)
    
    // 生成样本
    bernoulliSamples := make([]float64, 1000)
    normalSamples := make([]float64, 1000)
    
    for i := 0; i < 1000; i++ {
        bernoulliSamples[i] = bernoulli.Sample()
        normalSamples[i] = normal.Sample()
    }
    
    // 计算统计量
    stats := &Statistics{}
    
    fmt.Printf("Bernoulli Distribution (p=0.3):\n")
    fmt.Printf("Sample Mean: %.4f (Expected: %.4f)\n", 
        stats.Mean(bernoulliSamples), bernoulli.Mean())
    fmt.Printf("Sample Variance: %.4f (Expected: %.4f)\n", 
        stats.Variance(bernoulliSamples), bernoulli.Variance())
    
    fmt.Printf("\nNormal Distribution (μ=0, σ=1):\n")
    fmt.Printf("Sample Mean: %.4f (Expected: %.4f)\n", 
        stats.Mean(normalSamples), normal.Mean())
    fmt.Printf("Sample Variance: %.4f (Expected: %.4f)\n", 
        stats.Variance(normalSamples), normal.Variance())
    
    // 蒙特卡洛方法
    mc := &MonteCarlo{}
    
    // 计算积分 ∫[0,1] x^2 dx
    f := func(x float64) float64 { return x * x }
    integral := mc.Integrate(f, 0, 1, 10000)
    fmt.Printf("\nMonte Carlo Integration: ∫[0,1] x^2 dx ≈ %.4f (Exact: 0.3333)\n", integral)
    
    // 估计π
    piEstimate := mc.EstimatePi(100000)
    fmt.Printf("Monte Carlo π estimation: %.4f (Exact: 3.1416)\n", piEstimate)
    
    // 随机游走
    walk := mc.SimulateRandomWalk(100)
    fmt.Printf("\nRandom Walk final position: %.2f\n", walk[len(walk)-1])
}
```

## 总结

概率论是数学和统计学的基础，在机器学习、数据科学、金融工程等领域有广泛应用。通过形式化定义和Go语言实现，我们建立了从理论到实践的完整框架。

### 关键要点

1. **理论基础**: 概率空间、随机变量、概率分布
2. **核心定理**: 大数定律、中心极限定理、贝叶斯定理
3. **实现技术**: 随机数生成、统计计算、蒙特卡洛方法
4. **应用场景**: 机器学习、金融建模、风险评估等

### 进一步研究方向

1. **随机过程**: 马尔可夫链、布朗运动、泊松过程
2. **贝叶斯统计**: 贝叶斯推断、马尔可夫链蒙特卡洛
3. **信息论**: 熵、互信息、信道容量
4. **随机优化**: 随机梯度下降、遗传算法
