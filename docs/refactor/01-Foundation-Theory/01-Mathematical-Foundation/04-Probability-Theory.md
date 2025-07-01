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
  - [总结](#总结)

## 1. 基础定义

### 1.1 概率空间

**定义 1.1** (概率空间)
概率空间 $(\Omega, \mathcal{F}, P)$ 由以下三个部分组成：

- $\Omega$：样本空间，包含所有可能的结果
- $\mathcal{F}$：事件域，是 $\Omega$ 的子集的 $\sigma$-代数
- $P$：概率测度，满足：
  1. $P(\Omega) = 1$
  2. $P(A) \geq 0$ 对所有 $A \in \mathcal{F}$
  3. 可列可加性：对于互不相交的事件序列 $\{A_i\}$，$P(\bigcup_{i=1}^{\infty} A_i) = \sum_{i=1}^{\infty} P(A_i)$

**定义 1.2** (事件)
事件是样本空间 $\Omega$ 的子集，属于事件域 $\mathcal{F}$。

**定义 1.3** (基本事件)
基本事件是样本空间中的单个元素。

### 1.2 随机变量

**定义 1.4** (随机变量)
随机变量 $X$ 是从概率空间 $(\Omega, \mathcal{F}, P)$ 到实数集 $\mathbb{R}$ 的可测函数，即对于任意 $a \in \mathbb{R}$，$\{\omega \in \Omega : X(\omega) \leq a\} \in \mathcal{F}$。

**定义 1.5** (离散随机变量)
如果随机变量 $X$ 只取有限或可列个值，则称 $X$ 为离散随机变量。

**定义 1.6** (连续随机变量)
如果随机变量 $X$ 的分布函数 $F_X(x) = P(X \leq x)$ 是连续函数，则称 $X$ 为连续随机变量。

### 1.3 概率分布

**定义 1.7** (分布函数)
随机变量 $X$ 的分布函数定义为：
$$F_X(x) = P(X \leq x), \quad x \in \mathbb{R}$$

**定义 1.8** (概率质量函数)
离散随机变量 $X$ 的概率质量函数定义为：
$$p_X(x) = P(X = x), \quad x \in \mathbb{R}$$

**定义 1.9** (概率密度函数)
连续随机变量 $X$ 的概率密度函数 $f_X(x)$ 满足：
$$F_X(x) = \int_{-\infty}^x f_X(t) dt, \quad x \in \mathbb{R}$$

## 2. 概率性质

### 2.1 条件概率

**定义 2.1** (条件概率)
对于事件 $A, B$，且 $P(B) > 0$，$A$ 在 $B$ 条件下的条件概率定义为：
$$P(A|B) = \frac{P(A \cap B)}{P(B)}$$

**定理 2.1** (乘法公式)
对于事件 $A_1, A_2, \ldots, A_n$，且 $P(A_1 \cap A_2 \cap \cdots \cap A_{n-1}) > 0$：
$$P(A_1 \cap A_2 \cap \cdots \cap A_n) = P(A_1) \cdot P(A_2|A_1) \cdot P(A_3|A_1 \cap A_2) \cdots P(A_n|A_1 \cap A_2 \cap \cdots \cap A_{n-1})$$

**证明**：
使用数学归纳法。对于 $n=2$，结论显然成立。假设对于 $n=k$ 成立，考虑 $n=k+1$：
$$P(A_1 \cap A_2 \cap \cdots \cap A_{k+1}) = P(A_1 \cap A_2 \cap \cdots \cap A_k) \cdot P(A_{k+1}|A_1 \cap A_2 \cap \cdots \cap A_k)$$
根据归纳假设，结论成立。

### 2.2 独立性

**定义 2.2** (事件独立性)
事件 $A$ 和 $B$ 是独立的，当且仅当：
$$P(A \cap B) = P(A) \cdot P(B)$$

**定义 2.3** (条件独立性)
事件 $A$ 和 $B$ 在事件 $C$ 条件下是独立的，当且仅当：
$$P(A \cap B|C) = P(A|C) \cdot P(B|C)$$

### 2.3 贝叶斯定理

**定理 2.2** (贝叶斯定理)
对于事件 $A$ 和 $B$，且 $P(B) > 0$：
$$P(A|B) = \frac{P(B|A)P(A)}{P(B)}$$

**推论 2.1** (全概率公式)
对于完备事件组 $\{A_i\}$：
$$P(B) = \sum_i P(B|A_i)P(A_i)$$

## 3. 随机过程

### 3.1 马尔可夫链

**定义 3.1** (马尔可夫链)
随机过程 $\{X_n\}$ 是马尔可夫链，如果对于任意 $n$ 和状态 $i_0, i_1, \ldots, i_n, j$：
$$P(X_{n+1}=j|X_n=i_n,X_{n-1}=i_{n-1},\ldots,X_0=i_0) = P(X_{n+1}=j|X_n=i_n)$$

**定义 3.2** (转移概率)
马尔可夫链的转移概率矩阵 $P$ 定义为：
$$P_{ij} = P(X_{n+1}=j|X_n=i)$$

### 3.2 泊松过程

**定义 3.3** (泊松过程)
计数过程 $\{N(t), t \geq 0\}$ 是泊松过程，如果：

1. $N(0) = 0$
2. 增量独立
3. 对于任意 $t \geq 0$ 和 $h > 0$：
   $$P(N(t+h)-N(t)=1) = \lambda h + o(h)$$
   $$P(N(t+h)-N(t)\geq 2) = o(h)$$

### 3.3 布朗运动

**定义 3.4** (布朗运动)
随机过程 $\{B(t), t \geq 0\}$ 是布朗运动，如果：

1. $B(0) = 0$
2. 增量独立
3. 对于任意 $0 \leq s < t$，$B(t)-B(s) \sim N(0, t-s)$

## 4. 形式化定义

### 4.1 概率测度

**定义 4.1** (概率测度空间)
概率测度空间 $(\Omega, \mathcal{F}, P)$ 满足：

1. $\mathcal{F}$ 是 $\sigma$-代数
2. $P$ 是概率测度
3. $P$ 满足可列可加性

### 4.2 期望和方差

**定义 4.2** (期望)
随机变量 $X$ 的期望定义为：

- 离散情况：$E[X] = \sum_x x P(X=x)$
- 连续情况：$E[X] = \int_{-\infty}^{\infty} x f_X(x) dx$

**定义 4.3** (方差)
随机变量 $X$ 的方差定义为：
$$Var(X) = E[(X-E[X])^2] = E[X^2] - (E[X])^2$$

### 4.3 大数定律和中心极限定理

**定理 4.1** (大数定律)
设 $\{X_n\}$ 是独立同分布的随机变量序列，$E[X_1] = \mu$，则：
$$\frac{1}{n}\sum_{i=1}^n X_i \xrightarrow{P} \mu$$

**定理 4.2** (中心极限定理)
设 $\{X_n\}$ 是独立同分布的随机变量序列，$E[X_1] = \mu$，$Var(X_1) = \sigma^2$，则：
$$\frac{\sum_{i=1}^n X_i - n\mu}{\sigma\sqrt{n}} \xrightarrow{d} N(0,1)$$

## 5. Go语言实现

### 5.1 随机数生成

```go
// 随机数生成器
type RandomGenerator struct {
    seed int64
    rng  *rand.Rand
}

func NewRandomGenerator(seed int64) *RandomGenerator {
    return &RandomGenerator{
        seed: seed,
        rng:  rand.New(rand.NewSource(seed)),
    }
}

// 生成标准正态分布随机数
func (rg *RandomGenerator) StandardNormal() float64 {
    return rg.rng.NormFloat64()
}

// 生成指定参数的正态分布随机数
func (rg *RandomGenerator) Normal(mean, stdDev float64) float64 {
    return mean + stdDev*rg.StandardNormal()
}

// 生成泊松分布随机数
func (rg *RandomGenerator) Poisson(lambda float64) int {
    L := math.Exp(-lambda)
    k := 0
    p := 1.0
    
    for p > L {
        k++
        p *= rg.rng.Float64()
    }
    
    return k - 1
}
```

### 5.2 概率分布

```go
// 概率分布接口
type Distribution interface {
    PDF(x float64) float64    // 概率密度函数
    CDF(x float64) float64    // 累积分布函数
    Mean() float64            // 期望
    Variance() float64        // 方差
    Sample() float64          // 生成随机样本
}

// 正态分布
type NormalDistribution struct {
    mean   float64
    stdDev float64
    rng    *RandomGenerator
}

func (nd *NormalDistribution) PDF(x float64) float64 {
    z := (x - nd.mean) / nd.stdDev
    return math.Exp(-z*z/2) / (nd.stdDev * math.Sqrt(2*math.Pi))
}

func (nd *NormalDistribution) CDF(x float64) float64 {
    z := (x - nd.mean) / nd.stdDev
    return 0.5 * (1 + math.Erf(z/math.Sqrt(2)))
}

// 指数分布
type ExponentialDistribution struct {
    rate float64
    rng  *RandomGenerator
}

func (ed *ExponentialDistribution) PDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return ed.rate * math.Exp(-ed.rate*x)
}

func (ed *ExponentialDistribution) CDF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return 1 - math.Exp(-ed.rate*x)
}
```

### 5.3 统计推断

```go
// 统计推断器
type StatisticalInference struct {
    data []float64
}

// 计算样本均值
func (si *StatisticalInference) SampleMean() float64 {
    sum := 0.0
    for _, x := range si.data {
        sum += x
    }
    return sum / float64(len(si.data))
}

// 计算样本方差
func (si *StatisticalInference) SampleVariance() float64 {
    mean := si.SampleMean()
    sumSquares := 0.0
    for _, x := range si.data {
        diff := x - mean
        sumSquares += diff * diff
    }
    return sumSquares / float64(len(si.data)-1)
}

// 执行假设检验
func (si *StatisticalInference) HypothesisTest(nullMean float64, alpha float64) bool {
    n := float64(len(si.data))
    sampleMean := si.SampleMean()
    sampleStdDev := math.Sqrt(si.SampleVariance())
    
    // 计算t统计量
    t := (sampleMean - nullMean) / (sampleStdDev / math.Sqrt(n))
    
    // 查找临界值
    criticalValue := si.tCriticalValue(n-1, alpha)
    
    return math.Abs(t) > criticalValue
}
```

## 6. 应用实例

### 6.1 蒙特卡洛方法

```go
// 蒙特卡洛积分器
type MonteCarloIntegrator struct {
    rng *RandomGenerator
}

// 计算定积分
func (mci *MonteCarloIntegrator) Integrate(f func(float64) float64, a, b float64, n int) float64 {
    sum := 0.0
    for i := 0; i < n; i++ {
        x := a + (b-a)*mci.rng.rng.Float64()
        sum += f(x)
    }
    return (b - a) * sum / float64(n)
}

// 估计π值
func (mci *MonteCarloIntegrator) EstimatePi(n int) float64 {
    inside := 0
    for i := 0; i < n; i++ {
        x := mci.rng.rng.Float64()
        y := mci.rng.rng.Float64()
        if x*x + y*y <= 1 {
            inside++
        }
    }
    return 4 * float64(inside) / float64(n)
}
```

### 6.2 随机算法

```go
// 随机快速排序
func RandomizedQuickSort(arr []int, rng *RandomGenerator) {
    if len(arr) <= 1 {
        return
    }
    
    // 随机选择枢轴
    pivotIndex := rng.rng.Intn(len(arr))
    arr[0], arr[pivotIndex] = arr[pivotIndex], arr[0]
    
    // 分区
    pivot := arr[0]
    i := 1
    for j := 1; j < len(arr); j++ {
        if arr[j] < pivot {
            arr[i], arr[j] = arr[j], arr[i]
            i++
        }
    }
    
    // 将枢轴放到正确位置
    arr[0], arr[i-1] = arr[i-1], arr[0]
    
    // 递归排序
    RandomizedQuickSort(arr[:i-1], rng)
    RandomizedQuickSort(arr[i:], rng)
}
```

### 6.3 概率模型

```go
// 隐马尔可夫模型
type HMM struct {
    states        []string
    observations  []string
    initialProb   map[string]float64
    transitionProb map[string]map[string]float64
    emissionProb  map[string]map[string]float64
}

// 前向算法
func (hmm *HMM) Forward(observations []string) float64 {
    T := len(observations)
    N := len(hmm.states)
    
    // 初始化前向概率矩阵
    alpha := make([]map[string]float64, T)
    for t := range alpha {
        alpha[t] = make(map[string]float64)
    }
    
    // 初始化第一个时刻
    for _, state := range hmm.states {
        alpha[0][state] = hmm.initialProb[state] * hmm.emissionProb[state][observations[0]]
    }
    
    // 递推
    for t := 1; t < T; t++ {
        for _, state := range hmm.states {
            sum := 0.0
            for _, prevState := range hmm.states {
                sum += alpha[t-1][prevState] * hmm.transitionProb[prevState][state]
            }
            alpha[t][state] = sum * hmm.emissionProb[state][observations[t]]
        }
    }
    
    // 计算最终概率
    prob := 0.0
    for _, state := range hmm.states {
        prob += alpha[T-1][state]
    }
    
    return prob
}
```

## 总结

概率论为计算机科学提供了处理不确定性的理论基础，特别是在以下方面：

1. **随机算法设计**：
   - 蒙特卡洛方法
   - 随机化算法
   - 概率数据结构

2. **机器学习**：
   - 概率模型
   - 统计推断
   - 贝叶斯方法

3. **应用领域**：
   - 密码学
   - 网络协议
   - 人工智能
   - 量子计算

通过Go语言的实现，我们可以将这些理论概念转化为实用的工程解决方案，为实际问题提供可靠的概率分析工具。 