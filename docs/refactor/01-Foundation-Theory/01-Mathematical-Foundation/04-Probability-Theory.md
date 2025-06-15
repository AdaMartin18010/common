# 04-概率论 (Probability Theory)

## 概述

概率论是数学的一个分支，研究随机现象的数学规律。它提供了描述不确定性、随机性和随机过程的数学工具。在计算机科学中，概率论广泛应用于算法分析、机器学习、网络建模、密码学等领域。

## 1. 基本概念

### 1.1 样本空间和事件

**定义 1.1** (样本空间)
样本空间 $\Omega$ 是所有可能结果的集合。

**定义 1.2** (事件)
事件 $A$ 是样本空间 $\Omega$ 的子集，即 $A \subseteq \Omega$。

**定义 1.3** (事件代数)
事件代数 $\mathcal{F}$ 是 $\Omega$ 的子集族，满足：

1. $\Omega \in \mathcal{F}$
2. 如果 $A \in \mathcal{F}$，则 $A^c \in \mathcal{F}$
3. 如果 $A_1, A_2, \ldots \in \mathcal{F}$，则 $\bigcup_{i=1}^{\infty} A_i \in \mathcal{F}$

### 1.2 概率测度

**定义 1.4** (概率测度)
概率测度 $P$ 是定义在事件代数 $\mathcal{F}$ 上的函数，满足：

1. **非负性**：$P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
2. **规范性**：$P(\Omega) = 1$
3. **可列可加性**：对于互不相容的事件 $A_1, A_2, \ldots$，
   $$P\left(\bigcup_{i=1}^{\infty} A_i\right) = \sum_{i=1}^{\infty} P(A_i)$$

**定理 1.1** (概率的基本性质)
对于任意事件 $A, B \in \mathcal{F}$：

1. $P(A^c) = 1 - P(A)$
2. $P(A \cup B) = P(A) + P(B) - P(A \cap B)$
3. 如果 $A \subseteq B$，则 $P(A) \leq P(B)$

**证明**：

1. 由于 $A \cup A^c = \Omega$ 且 $A \cap A^c = \emptyset$，由可列可加性：
   $$1 = P(\Omega) = P(A \cup A^c) = P(A) + P(A^c)$$
   因此 $P(A^c) = 1 - P(A)$

2. 由于 $A \cup B = A \cup (B \setminus A)$ 且 $A \cap (B \setminus A) = \emptyset$：
   $$P(A \cup B) = P(A) + P(B \setminus A)$$
   又因为 $B = (B \setminus A) \cup (A \cap B)$ 且 $(B \setminus A) \cap (A \cap B) = \emptyset$：
   $$P(B) = P(B \setminus A) + P(A \cap B)$$
   因此 $P(B \setminus A) = P(B) - P(A \cap B)$，代入得：
   $$P(A \cup B) = P(A) + P(B) - P(A \cap B)$$

## 2. 条件概率和独立性

### 2.1 条件概率

**定义 2.1** (条件概率)
对于事件 $A, B \in \mathcal{F}$，且 $P(B) > 0$，事件 $A$ 在事件 $B$ 发生的条件下的条件概率为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 2.1** (乘法公式)
对于事件 $A_1, A_2, \ldots, A_n$：
$$P(A_1 \cap A_2 \cap \cdots \cap A_n) = P(A_1) \cdot P(A_2|A_1) \cdot P(A_3|A_1 \cap A_2) \cdots P(A_n|A_1 \cap A_2 \cap \cdots \cap A_{n-1})$$

### 2.2 独立性

**定义 2.2** (独立性)
事件 $A$ 和 $B$ 是独立的，当且仅当：
$$P(A \cap B) = P(A) \cdot P(B)$$

**定义 2.3** (条件独立性)
事件 $A$ 和 $B$ 在事件 $C$ 的条件下是独立的，当且仅当：
$$P(A \cap B|C) = P(A|C) \cdot P(B|C)$$

### 2.3 全概率公式和贝叶斯公式

**定理 2.2** (全概率公式)
如果事件 $B_1, B_2, \ldots, B_n$ 构成样本空间的一个划分，即：

1. $B_i \cap B_j = \emptyset$ 对所有 $i \neq j$
2. $\bigcup_{i=1}^{n} B_i = \Omega$

则对于任意事件 $A$：
$$P(A) = \sum_{i=1}^{n} P(A|B_i) \cdot P(B_i)$$

**定理 2.3** (贝叶斯公式)
在相同条件下：
$$P(B_i|A) = \frac{P(A|B_i) \cdot P(B_i)}{\sum_{j=1}^{n} P(A|B_j) \cdot P(B_j)}$$

## 3. 随机变量

### 3.1 随机变量的定义

**定义 3.1** (随机变量)
随机变量 $X$ 是从样本空间 $\Omega$ 到实数集 $\mathbb{R}$ 的可测函数。

**定义 3.2** (分布函数)
随机变量 $X$ 的分布函数 $F_X(x)$ 定义为：
$$F_X(x) = P(X \leq x)$$

**性质**：

1. $F_X(x)$ 是非减函数
2. $\lim_{x \to -\infty} F_X(x) = 0$
3. $\lim_{x \to +\infty} F_X(x) = 1$
4. $F_X(x)$ 是右连续的

### 3.2 离散随机变量

**定义 3.3** (概率质量函数)
离散随机变量 $X$ 的概率质量函数 $p_X(x)$ 定义为：
$$p_X(x) = P(X = x)$$

**性质**：

1. $p_X(x) \geq 0$ 对所有 $x$
2. $\sum_{x} p_X(x) = 1$

### 3.3 连续随机变量

**定义 3.4** (概率密度函数)
连续随机变量 $X$ 的概率密度函数 $f_X(x)$ 定义为：
$$f_X(x) = \frac{d}{dx} F_X(x)$$

**性质**：

1. $f_X(x) \geq 0$ 对所有 $x$
2. $\int_{-\infty}^{\infty} f_X(x) dx = 1$
3. $P(a \leq X \leq b) = \int_{a}^{b} f_X(x) dx$

## 4. 期望和方差

### 4.1 数学期望

**定义 4.1** (数学期望)
随机变量 $X$ 的数学期望 $E[X]$ 定义为：

对于离散随机变量：
$$E[X] = \sum_{x} x \cdot p_X(x)$$

对于连续随机变量：
$$E[X] = \int_{-\infty}^{\infty} x \cdot f_X(x) dx$$

**定理 4.1** (期望的线性性质)
对于任意常数 $a, b$ 和随机变量 $X, Y$：
$$E[aX + bY] = aE[X] + bE[Y]$$

### 4.2 方差

**定义 4.2** (方差)
随机变量 $X$ 的方差 $Var(X)$ 定义为：
$$Var(X) = E[(X - E[X])^2] = E[X^2] - (E[X])^2$$

**定理 4.2** (方差的性质)
对于任意常数 $a, b$ 和随机变量 $X$：
$$Var(aX + b) = a^2 Var(X)$$

## 5. Go语言实现

### 5.1 概率分布接口

```go
package probability

import (
    "math"
    "math/rand"
    "time"
)

// Distribution 概率分布接口
type Distribution interface {
    // PDF 概率密度函数（连续）或概率质量函数（离散）
    PDF(x float64) float64
    
    // CDF 累积分布函数
    CDF(x float64) float64
    
    // Mean 期望
    Mean() float64
    
    // Variance 方差
    Variance() float64
    
    // Sample 生成随机样本
    Sample() float64
}

// RandomGenerator 随机数生成器
type RandomGenerator struct {
    rng *rand.Rand
}

// NewRandomGenerator 创建随机数生成器
func NewRandomGenerator(seed int64) *RandomGenerator {
    if seed == 0 {
        seed = time.Now().UnixNano()
    }
    return &RandomGenerator{
        rng: rand.New(rand.NewSource(seed)),
    }
}

// Uniform 生成[0,1)均匀分布随机数
func (rg *RandomGenerator) Uniform() float64 {
    return rg.rng.Float64()
}

// UniformRange 生成[a,b)均匀分布随机数
func (rg *RandomGenerator) UniformRange(a, b float64) float64 {
    return a + (b-a)*rg.Uniform()
}
```

### 5.2 离散分布

#### 5.2.1 伯努利分布

```go
// BernoulliDistribution 伯努利分布
type BernoulliDistribution struct {
    p float64 // 成功概率
    rg *RandomGenerator
}

// NewBernoulli 创建伯努利分布
func NewBernoulli(p float64, seed int64) *BernoulliDistribution {
    if p < 0 || p > 1 {
        panic("概率必须在[0,1]范围内")
    }
    return &BernoulliDistribution{
        p:  p,
        rg: NewRandomGenerator(seed),
    }
}

// PDF 概率质量函数
func (b *BernoulliDistribution) PDF(x float64) float64 {
    if x == 0 {
        return 1 - b.p
    } else if x == 1 {
        return b.p
    }
    return 0
}

// CDF 累积分布函数
func (b *BernoulliDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0
    } else if x < 1 {
        return 1 - b.p
    }
    return 1
}

// Mean 期望
func (b *BernoulliDistribution) Mean() float64 {
    return b.p
}

// Variance 方差
func (b *BernoulliDistribution) Variance() float64 {
    return b.p * (1 - b.p)
}

// Sample 生成随机样本
func (b *BernoulliDistribution) Sample() float64 {
    if b.rg.Uniform() < b.p {
        return 1
    }
    return 0
}
```

#### 5.2.2 二项分布

```go
// BinomialDistribution 二项分布
type BinomialDistribution struct {
    n int     // 试验次数
    p float64 // 成功概率
    rg *RandomGenerator
}

// NewBinomial 创建二项分布
func NewBinomial(n int, p float64, seed int64) *BinomialDistribution {
    if n < 0 {
        panic("试验次数必须非负")
    }
    if p < 0 || p > 1 {
        panic("概率必须在[0,1]范围内")
    }
    return &BinomialDistribution{
        n:  n,
        p:  p,
        rg: NewRandomGenerator(seed),
    }
}

// factorial 计算阶乘
func factorial(n int) int64 {
    if n <= 1 {
        return 1
    }
    result := int64(1)
    for i := 2; i <= n; i++ {
        result *= int64(i)
    }
    return result
}

// combination 计算组合数 C(n,k)
func combination(n, k int) int64 {
    if k > n {
        return 0
    }
    if k > n/2 {
        k = n - k
    }
    result := int64(1)
    for i := 0; i < k; i++ {
        result = result * int64(n-i) / int64(i+1)
    }
    return result
}

// PDF 概率质量函数
func (b *BinomialDistribution) PDF(x float64) float64 {
    k := int(x)
    if k < 0 || k > b.n {
        return 0
    }
    
    c := combination(b.n, k)
    prob := float64(c) * math.Pow(b.p, float64(k)) * math.Pow(1-b.p, float64(b.n-k))
    return prob
}

// CDF 累积分布函数
func (b *BinomialDistribution) CDF(x float64) float64 {
    k := int(x)
    if k < 0 {
        return 0
    }
    if k >= b.n {
        return 1
    }
    
    sum := 0.0
    for i := 0; i <= k; i++ {
        sum += b.PDF(float64(i))
    }
    return sum
}

// Mean 期望
func (b *BinomialDistribution) Mean() float64 {
    return float64(b.n) * b.p
}

// Variance 方差
func (b *BinomialDistribution) Variance() float64 {
    return float64(b.n) * b.p * (1 - b.p)
}

// Sample 生成随机样本
func (b *BinomialDistribution) Sample() float64 {
    sum := 0
    for i := 0; i < b.n; i++ {
        if b.rg.Uniform() < b.p {
            sum++
        }
    }
    return float64(sum)
}
```

### 5.3 连续分布

#### 5.3.1 正态分布

```go
// NormalDistribution 正态分布
type NormalDistribution struct {
    mu    float64 // 均值
    sigma float64 // 标准差
    rg    *RandomGenerator
}

// NewNormal 创建正态分布
func NewNormal(mu, sigma float64, seed int64) *NormalDistribution {
    if sigma <= 0 {
        panic("标准差必须为正数")
    }
    return &NormalDistribution{
        mu:    mu,
        sigma: sigma,
        rg:    NewRandomGenerator(seed),
    }
}

// PDF 概率密度函数
func (n *NormalDistribution) PDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return math.Exp(-0.5*z*z) / (n.sigma * math.Sqrt(2*math.Pi))
}

// CDF 累积分布函数（近似计算）
func (n *NormalDistribution) CDF(x float64) float64 {
    z := (x - n.mu) / n.sigma
    return 0.5 * (1 + erf(z/math.Sqrt(2)))
}

// erf 误差函数（近似）
func erf(x float64) float64 {
    // 使用近似公式
    if x < 0 {
        return -erf(-x)
    }
    
    a1 := 0.254829592
    a2 := -0.284496736
    a3 := 1.421413741
    a4 := -1.453152027
    a5 := 1.061405429
    p := 0.3275911
    
    t := 1.0 / (1.0 + p*x)
    y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)
    return y
}

// Mean 期望
func (n *NormalDistribution) Mean() float64 {
    return n.mu
}

// Variance 方差
func (n *NormalDistribution) Variance() float64 {
    return n.sigma * n.sigma
}

// Sample 生成随机样本（Box-Muller变换）
func (n *NormalDistribution) Sample() float64 {
    // Box-Muller变换
    u1 := n.rg.Uniform()
    u2 := n.rg.Uniform()
    
    z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return n.mu + n.sigma*z0
}
```

#### 5.3.2 指数分布

```go
// ExponentialDistribution 指数分布
type ExponentialDistribution struct {
    lambda float64 // 参数λ
    rg     *RandomGenerator
}

// NewExponential 创建指数分布
func NewExponential(lambda float64, seed int64) *ExponentialDistribution {
    if lambda <= 0 {
        panic("参数λ必须为正数")
    }
    return &ExponentialDistribution{
        lambda: lambda,
        rg:     NewRandomGenerator(seed),
    }
}

// PDF 概率密度函数
func (e *ExponentialDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return e.lambda * math.Exp(-e.lambda*x)
}

// CDF 累积分布函数
func (e *ExponentialDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return 1 - math.Exp(-e.lambda*x)
}

// Mean 期望
func (e *ExponentialDistribution) Mean() float64 {
    return 1.0 / e.lambda
}

// Variance 方差
func (e *ExponentialDistribution) Variance() float64 {
    return 1.0 / (e.lambda * e.lambda)
}

// Sample 生成随机样本
func (e *ExponentialDistribution) Sample() float64 {
    u := e.rg.Uniform()
    return -math.Log(1-u) / e.lambda
}
```

### 5.4 统计推断

```go
// Statistics 统计工具
type Statistics struct{}

// NewStatistics 创建统计工具
func NewStatistics() *Statistics {
    return &Statistics{}
}

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
        diff := x - mean
        sum += diff * diff
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
    z := 1.96 // 95%置信水平
    margin := z * std / math.Sqrt(n)
    
    return mean - margin, mean + margin
}

// Correlation 计算相关系数
func (s *Statistics) Correlation(x, y []float64) float64 {
    if len(x) != len(y) || len(x) == 0 {
        return 0
    }
    
    meanX := s.Mean(x)
    meanY := s.Mean(y)
    
    numerator := 0.0
    sumXSquared := 0.0
    sumYSquared := 0.0
    
    for i := 0; i < len(x); i++ {
        diffX := x[i] - meanX
        diffY := y[i] - meanY
        numerator += diffX * diffY
        sumXSquared += diffX * diffX
        sumYSquared += diffY * diffY
    }
    
    denominator := math.Sqrt(sumXSquared * sumYSquared)
    if denominator == 0 {
        return 0
    }
    
    return numerator / denominator
}
```

## 6. 应用实例

### 6.1 蒙特卡洛方法

```go
// MonteCarlo 蒙特卡洛方法
type MonteCarlo struct {
    rg *RandomGenerator
}

// NewMonteCarlo 创建蒙特卡洛模拟器
func NewMonteCarlo(seed int64) *MonteCarlo {
    return &MonteCarlo{
        rg: NewRandomGenerator(seed),
    }
}

// EstimatePi 估算π值
func (mc *MonteCarlo) EstimatePi(n int) float64 {
    inside := 0
    for i := 0; i < n; i++ {
        x := mc.rg.UniformRange(-1, 1)
        y := mc.rg.UniformRange(-1, 1)
        if x*x+y*y <= 1 {
            inside++
        }
    }
    return 4.0 * float64(inside) / float64(n)
}

// EstimateIntegral 估算定积分
func (mc *MonteCarlo) EstimateIntegral(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := mc.rg.UniformRange(a, b)
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}
```

### 6.2 随机游走

```go
// RandomWalk 随机游走
type RandomWalk struct {
    position float64
    step     float64
    rg       *RandomGenerator
}

// NewRandomWalk 创建随机游走
func NewRandomWalk(initialPosition, stepSize float64, seed int64) *RandomWalk {
    return &RandomWalk{
        position: initialPosition,
        step:     stepSize,
        rg:       NewRandomGenerator(seed),
    }
}

// Step 执行一步随机游走
func (rw *RandomWalk) Step() float64 {
    if rw.rg.Uniform() < 0.5 {
        rw.position += rw.step
    } else {
        rw.position -= rw.step
    }
    return rw.position
}

// Simulate 模拟n步随机游走
func (rw *RandomWalk) Simulate(n int) []float64 {
    positions := make([]float64, n+1)
    positions[0] = rw.position
    
    for i := 1; i <= n; i++ {
        positions[i] = rw.Step()
    }
    
    return positions
}
```

## 7. 性能分析

### 7.1 算法复杂度

| 操作 | 时间复杂度 | 空间复杂度 |
|------|------------|------------|
| 生成随机数 | $O(1)$ | $O(1)$ |
| 计算PDF | $O(1)$ | $O(1)$ |
| 计算CDF | $O(1)$ | $O(1)$ |
| 生成样本 | $O(1)$ | $O(1)$ |
| 统计计算 | $O(n)$ | $O(1)$ |

### 7.2 数值精度

- **浮点数精度**：使用 `float64` 类型保证数值精度
- **数值稳定性**：避免数值溢出和下溢
- **随机数质量**：使用高质量的随机数生成器

## 8. 总结

概率论为计算机科学提供了处理不确定性和随机性的数学基础。通过Go语言的实现，我们可以：

1. **概率分布**：实现各种离散和连续概率分布
2. **随机数生成**：生成高质量的随机数
3. **统计推断**：进行数据分析和统计推断
4. **蒙特卡洛方法**：解决复杂数值计算问题
5. **随机过程**：模拟随机现象和过程

概率论的理论基础和算法实现为构建可靠的软件系统提供了重要的数学支撑。

## 参考文献

1. Ross, S. M. (2014). Introduction to Probability Models (11th ed.). Academic Press.
2. Casella, G., & Berger, R. L. (2002). Statistical Inference (2nd ed.). Duxbury Press.
3. Feller, W. (1968). An Introduction to Probability Theory and Its Applications (3rd ed.). Wiley.
