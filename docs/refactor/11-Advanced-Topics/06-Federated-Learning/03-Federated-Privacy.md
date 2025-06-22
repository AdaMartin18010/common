# 3. 联邦学习隐私保护

## 概述

联邦学习隐私保护（Federated Learning Privacy）是确保在分布式训练过程中保护用户数据隐私的关键技术，包括差分隐私、安全多方计算和同态加密等方法。

## 3.1 差分隐私

### 3.1.1 差分隐私定义

对于相邻数据集 ```latex
$D$
``` 和 ```latex
$D'$
```，算法 ```latex
$A$
``` 满足 ```latex
$(\epsilon, \delta)$
```-差分隐私：

```latex
$```latex
$\Pr[A(D) \in S] \leq e^\epsilon \Pr[A(D') \in S] + \delta$
```$
```

其中 ```latex
$\epsilon$
``` 是隐私预算，```latex
$\delta$
``` 是失败概率。

### 3.1.2 拉普拉斯机制

拉普拉斯机制添加噪声：

```latex
$```latex
$A(D) = f(D) + \text{Lap}\left(\frac{\Delta f}{\epsilon}\right)$
```$
```

其中 ```latex
$\Delta f$
``` 是函数 ```latex
$f$
``` 的敏感度。

### 3.1.3 高斯机制

高斯机制添加高斯噪声：

```latex
$```latex
$A(D) = f(D) + \mathcal{N}\left(0, \frac{\Delta f^2 \log(1/\delta)}{2\epsilon^2}\right)$
```$
```

## 3.2 联邦学习中的差分隐私

### 3.2.1 客户端级差分隐私

客户端级差分隐私保护整个客户端的数据：

```latex
$```latex
$\tilde{w}_i = w_i + \mathcal{N}\left(0, \frac{c^2 \log(1/\delta)}{2\epsilon^2} I\right)$
```$
```

其中 ```latex
$c$
``` 是裁剪范数。

### 3.2.2 样本级差分隐私

样本级差分隐私保护单个样本：

```latex
$```latex
$\tilde{g}_i = \text{Clip}(g_i, c) + \mathcal{N}\left(0, \frac{c^2 \log(1/\delta)}{2\epsilon^2} I\right)$
```$
```

其中 ```latex
$\text{Clip}(g, c) = g \cdot \min(1, c/\|g\|)$
```。

### 3.2.3 隐私预算管理

隐私预算的组成：

```latex
$```latex
$\epsilon_{total} = \sum_{t=1}^T \epsilon_t$
```$
```

其中 ```latex
$T$
``` 是通信轮数。

## 3.3 安全多方计算

### 3.3.1 秘密共享

```latex
$(t, n)$
``` 秘密共享将秘密 ```latex
$s$
``` 分割为 ```latex
$n$
``` 个份额：

```latex
$```latex
$s = \sum_{i=1}^t s_i \pmod{p}$
```$
```

其中任意 ```latex
$t$
``` 个份额可以重构秘密。

### 3.3.2 安全聚合

安全聚合协议：

```latex
$```latex
$y = \sum_{i=1}^n x_i = \sum_{i=1}^n \left(\sum_{j=1}^n s_{i,j}\right) = \sum_{j=1}^n \left(\sum_{i=1}^n s_{i,j}\right)$
```$
```

其中 ```latex
$s_{i,j}$
``` 是客户端 ```latex
$i$
``` 发送给客户端 ```latex
$j$
``` 的份额。

### 3.3.3 同态加密

同态加密支持在密文上进行计算：

```latex
$```latex
$\text{Enc}(m_1) \oplus \text{Enc}(m_2) = \text{Enc}(m_1 + m_2)$
```$
$```latex
$\text{Enc}(m_1) \otimes \text{Enc}(m_2) = \text{Enc}(m_1 \times m_2)$
```$
```

## 3.4 联邦学习隐私攻击

### 3.4.1 成员推理攻击

成员推理攻击判断样本是否在训练集中：

```latex
$```latex
$P(\text{member}|x) = \sigma(f(x; \theta))$
```$
```

其中 ```latex
$f$
``` 是攻击模型，```latex
$\sigma$
``` 是sigmoid函数。

### 3.4.2 模型反演攻击

模型反演攻击重构训练数据：

```latex
$```latex
$\hat{x} = \arg\min_x \mathcal{L}(f(x; \theta), y) + \lambda R(x)$
```$
```

其中 ```latex
$R(x)$
``` 是正则化项。

### 3.4.3 属性推理攻击

属性推理攻击推断敏感属性：

```latex
$```latex
$P(a|x) = \sigma(g(f(x; \theta)))$
```$
```

其中 ```latex
$g$
``` 是属性推理模型。

## 3.5 Go语言实现

### 3.5.1 差分隐私实现

```go
package federatedprivacy

import (
    "math"
    "math/rand"
)

// DifferentialPrivacy 差分隐私
type DifferentialPrivacy struct {
    Epsilon float64
    Delta   float64
}

// NewDifferentialPrivacy 创建差分隐私
func NewDifferentialPrivacy(epsilon, delta float64) *DifferentialPrivacy {
    return &DifferentialPrivacy{
        Epsilon: epsilon,
        Delta:   delta,
    }
}

// AddLaplaceNoise 添加拉普拉斯噪声
func (dp *DifferentialPrivacy) AddLaplaceNoise(value float64, sensitivity float64) float64 {
    scale := sensitivity / dp.Epsilon
    noise := dp.sampleLaplace(0, scale)
    return value + noise
}

// AddGaussianNoise 添加高斯噪声
func (dp *DifferentialPrivacy) AddGaussianNoise(value float64, sensitivity float64) float64 {
    variance := (sensitivity * sensitivity * math.Log(1/dp.Delta)) / (2 * dp.Epsilon * dp.Epsilon)
    noise := dp.sampleGaussian(0, math.Sqrt(variance))
    return value + noise
}

// sampleLaplace 采样拉普拉斯分布
func (dp *DifferentialPrivacy) sampleLaplace(loc, scale float64) float64 {
    u := rand.Float64() - 0.5
    if u < 0 {
        return loc + scale*math.Log(1+2*u)
    } else {
        return loc - scale*math.Log(1-2*u)
    }
}

// sampleGaussian 采样高斯分布
func (dp *DifferentialPrivacy) sampleGaussian(mean, stddev float64) float64 {
    // Box-Muller变换
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z0 := mean + stddev*math.Sqrt(-2*math.Log(u1))*math.Cos(2*math.Pi*u2)
    return z0
}

// ClipGradient 裁剪梯度
func (dp *DifferentialPrivacy) ClipGradient(gradient []float64, clipNorm float64) []float64 {
    clipped := make([]float64, len(gradient))
    copy(clipped, gradient)
    
    // 计算L2范数
    norm := 0.0
    for _, val := range gradient {
        norm += val * val
    }
    norm = math.Sqrt(norm)
    
    // 裁剪
    if norm > clipNorm {
        scale := clipNorm / norm
        for i := range clipped {
            clipped[i] *= scale
        }
    }
    
    return clipped
}

// AddNoiseToGradient 向梯度添加噪声
func (dp *DifferentialPrivacy) AddNoiseToGradient(gradient []float64, clipNorm float64) []float64 {
    // 裁剪梯度
    clipped := dp.ClipGradient(gradient, clipNorm)
    
    // 添加高斯噪声
    noisy := make([]float64, len(clipped))
    for i, val := range clipped {
        noisy[i] = dp.AddGaussianNoise(val, clipNorm)
    }
    
    return noisy
}
```

### 3.5.2 安全多方计算

```go
// SecretSharing 秘密共享
type SecretSharing struct {
    Threshold int
    Parties   int
    Prime     int64
}

// NewSecretSharing 创建秘密共享
func NewSecretSharing(threshold, parties int) *SecretSharing {
    return &SecretSharing{
        Threshold: threshold,
        Parties:   parties,
        Prime:     1000000007, // 大素数
    }
}

// Share 生成份额
func (ss *SecretSharing) Share(secret int64) []int64 {
    shares := make([]int64, ss.Parties)
    
    // 生成随机系数
    coefficients := make([]int64, ss.Threshold-1)
    for i := 0; i < ss.Threshold-1; i++ {
        coefficients[i] = rand.Int63n(ss.Prime)
    }
    
    // 计算份额
    for i := 0; i < ss.Parties; i++ {
        x := int64(i + 1)
        share := secret
        
        for j := 0; j < ss.Threshold-1; j++ {
            share = (share + coefficients[j]*pow(x, j+1, ss.Prime)) % ss.Prime
        }
        
        shares[i] = share
    }
    
    return shares
}

// Reconstruct 重构秘密
func (ss *SecretSharing) Reconstruct(shares []int64, indices []int) int64 {
    if len(shares) < ss.Threshold {
        return 0
    }
    
    // 拉格朗日插值
    secret := int64(0)
    for i := 0; i < ss.Threshold; i++ {
        numerator := int64(1)
        denominator := int64(1)
        
        for j := 0; j < ss.Threshold; j++ {
            if i != j {
                numerator = (numerator * int64(-indices[j]-1)) % ss.Prime
                denominator = (denominator * int64(indices[i]-indices[j])) % ss.Prime
            }
        }
        
        if denominator < 0 {
            denominator += ss.Prime
        }
        
        inverse := modInverse(denominator, ss.Prime)
        term := (shares[i] * numerator * inverse) % ss.Prime
        secret = (secret + term) % ss.Prime
    }
    
    return secret
}

// pow 模幂运算
func pow(base, exp int64, mod int64) int64 {
    result := int64(1)
    base = base % mod
    
    for exp > 0 {
        if exp%2 == 1 {
            result = (result * base) % mod
        }
        exp = exp >> 1
        base = (base * base) % mod
    }
    
    return result
}

// modInverse 模逆元
func modInverse(a, m int64) int64 {
    m0 := m
    y := int64(0)
    x := int64(1)
    
    if m == 1 {
        return 0
    }
    
    for a > 1 {
        q := a / m
        t := m
        
        m = int(a % m)
        a = t
        t = y
        
        y = x - q*y
        x = t
    }
    
    if x < 0 {
        x += m0
    }
    
    return x
}

// SecureAggregation 安全聚合
type SecureAggregation struct {
    SecretSharing *SecretSharing
}

// NewSecureAggregation 创建安全聚合
func NewSecureAggregation(threshold, parties int) *SecureAggregation {
    return &SecureAggregation{
        SecretSharing: NewSecretSharing(threshold, parties),
    }
}

// Aggregate 安全聚合
func (sa *SecureAggregation) Aggregate(values []int64) int64 {
    // 每个客户端生成份额
    shares := make([][]int64, len(values))
    for i, value := range values {
        shares[i] = sa.SecretSharing.Share(value)
    }
    
    // 计算聚合份额
    aggregatedShares := make([]int64, sa.SecretSharing.Parties)
    for i := 0; i < sa.SecretSharing.Parties; i++ {
        sum := int64(0)
        for j := 0; j < len(values); j++ {
            sum = (sum + shares[j][i]) % sa.SecretSharing.Prime
        }
        aggregatedShares[i] = sum
    }
    
    // 重构聚合结果
    indices := make([]int, sa.SecretSharing.Threshold)
    for i := 0; i < sa.SecretSharing.Threshold; i++ {
        indices[i] = i
    }
    
    return sa.SecretSharing.Reconstruct(aggregatedShares[:sa.SecretSharing.Threshold], indices)
}
```

### 3.5.3 隐私攻击检测

```go
// PrivacyAttack 隐私攻击接口
type PrivacyAttack interface {
    Attack(model *Model, target interface{}) interface{}
}

// MembershipInferenceAttack 成员推理攻击
type MembershipInferenceAttack struct {
    AttackModel *Model
    Threshold   float64
}

// NewMembershipInferenceAttack 创建成员推理攻击
func NewMembershipInferenceAttack() *MembershipInferenceAttack {
    return &MembershipInferenceAttack{
        AttackModel: NewModel(),
        Threshold:   0.5,
    }
}

// Attack 执行攻击
func (mia *MembershipInferenceAttack) Attack(model *Model, target interface{}) interface{} {
    // 简化实现：基于模型置信度判断成员关系
    if prediction, ok := target.([]float64); ok {
        confidence := mia.computeConfidence(prediction)
        return confidence > mia.Threshold
    }
    return false
}

// computeConfidence 计算置信度
func (mia *MembershipInferenceAttack) computeConfidence(prediction []float64) float64 {
    maxProb := 0.0
    for _, prob := range prediction {
        if prob > maxProb {
            maxProb = prob
        }
    }
    return maxProb
}

// ModelInversionAttack 模型反演攻击
type ModelInversionAttack struct {
    LearningRate float64
    MaxIter      int
}

// NewModelInversionAttack 创建模型反演攻击
func NewModelInversionAttack() *ModelInversionAttack {
    return &ModelInversionAttack{
        LearningRate: 0.01,
        MaxIter:      1000,
    }
}

// Attack 执行攻击
func (mia *ModelInversionAttack) Attack(model *Model, target interface{}) interface{} {
    if label, ok := target.(int); ok {
        return mia.invertModel(model, label)
    }
    return nil
}

// invertModel 反演模型
func (mia *ModelInversionAttack) invertModel(model *Model, label int) []float64 {
    // 初始化随机输入
    input := make([]float64, 784) // MNIST输入维度
    for i := range input {
        input[i] = rand.Float64()
    }
    
    // 梯度下降优化
    for iter := 0; iter < mia.MaxIter; iter++ {
        // 前向传播
        output := model.Forward(input)
        
        // 计算损失
        loss := mia.computeLoss(output, label)
        
        // 计算梯度
        gradients := mia.computeGradients(model, input, label)
        
        // 更新输入
        for i := range input {
            input[i] -= mia.LearningRate * gradients[i]
        }
    }
    
    return input
}

// computeLoss 计算损失
func (mia *ModelInversionAttack) computeLoss(output []float64, label int) float64 {
    target := make([]float64, len(output))
    target[label] = 1.0
    
    loss := 0.0
    for i := 0; i < len(output); i++ {
        diff := output[i] - target[i]
        loss += diff * diff
    }
    
    return loss
}

// computeGradients 计算梯度
func (mia *ModelInversionAttack) computeGradients(model *Model, input []float64, label int) []float64 {
    // 简化实现：数值梯度
    epsilon := 1e-6
    gradients := make([]float64, len(input))
    
    for i := range input {
        // 前向扰动
        input[i] += epsilon
        outputPlus := model.Forward(input)
        lossPlus := mia.computeLoss(outputPlus, label)
        
        // 后向扰动
        input[i] -= 2 * epsilon
        outputMinus := model.Forward(input)
        lossMinus := mia.computeLoss(outputMinus, label)
        
        // 恢复
        input[i] += epsilon
        
        // 计算梯度
        gradients[i] = (lossPlus - lossMinus) / (2 * epsilon)
    }
    
    return gradients
}
```

### 3.5.4 隐私保护联邦学习

```go
// PrivacyPreservingFL 隐私保护联邦学习
type PrivacyPreservingFL struct {
    DifferentialPrivacy *DifferentialPrivacy
    SecureAggregation   *SecureAggregation
    Config              *PrivacyConfig
}

// PrivacyConfig 隐私配置
type PrivacyConfig struct {
    Epsilon       float64
    Delta         float64
    ClipNorm      float64
    UseDP         bool
    UseSecureAgg  bool
    Threshold     int
    Parties       int
}

// NewPrivacyPreservingFL 创建隐私保护联邦学习
func NewPrivacyPreservingFL(config *PrivacyConfig) *PrivacyPreservingFL {
    var dp *DifferentialPrivacy
    var sa *SecureAggregation
    
    if config.UseDP {
        dp = NewDifferentialPrivacy(config.Epsilon, config.Delta)
    }
    
    if config.UseSecureAgg {
        sa = NewSecureAggregation(config.Threshold, config.Parties)
    }
    
    return &PrivacyPreservingFL{
        DifferentialPrivacy: dp,
        SecureAggregation:   sa,
        Config:              config,
    }
}

// TrainWithPrivacy 隐私保护训练
func (ppfl *PrivacyPreservingFL) TrainWithPrivacy(globalModel *Model, clientModels []*Model, clientWeights []float64) *Model {
    if ppfl.Config.UseSecureAgg {
        return ppfl.secureAggregation(globalModel, clientModels, clientWeights)
    } else {
        return ppfl.differentialPrivacyAggregation(globalModel, clientModels, clientWeights)
    }
}

// differentialPrivacyAggregation 差分隐私聚合
func (ppfl *PrivacyPreservingFL) differentialPrivacyAggregation(globalModel *Model, clientModels []*Model, clientWeights []float64) *Model {
    if ppfl.DifferentialPrivacy == nil {
        return ppfl.standardAggregation(globalModel, clientModels, clientWeights)
    }
    
    aggregated := NewModel()
    aggregated.Weights = make([]float64, len(clientModels[0].Weights))
    aggregated.Bias = make([]float64, len(clientModels[0].Bias))
    
    // 聚合权重
    for i := 0; i < len(clientModels[0].Weights); i++ {
        sum := 0.0
        for j := 0; j < len(clientModels); j++ {
            sum += clientWeights[j] * clientModels[j].Weights[i]
        }
        
        // 添加差分隐私噪声
        noisySum := ppfl.DifferentialPrivacy.AddGaussianNoise(sum, ppfl.Config.ClipNorm)
        aggregated.Weights[i] = noisySum
    }
    
    // 聚合偏置
    for i := 0; i < len(clientModels[0].Bias); i++ {
        sum := 0.0
        for j := 0; j < len(clientModels); j++ {
            sum += clientWeights[j] * clientModels[j].Bias[i]
        }
        
        // 添加差分隐私噪声
        noisySum := ppfl.DifferentialPrivacy.AddGaussianNoise(sum, ppfl.Config.ClipNorm)
        aggregated.Bias[i] = noisySum
    }
    
    return aggregated
}

// secureAggregation 安全聚合
func (ppfl *PrivacyPreservingFL) secureAggregation(globalModel *Model, clientModels []*Model, clientWeights []float64) *Model {
    if ppfl.SecureAggregation == nil {
        return ppfl.standardAggregation(globalModel, clientModels, clientWeights)
    }
    
    // 将模型参数转换为整数进行安全聚合
    // 简化实现：只聚合第一个权重参数
    values := make([]int64, len(clientModels))
    for i, model := range clientModels {
        values[i] = int64(model.Weights[0] * 1000000) // 放大为整数
    }
    
    aggregatedValue := ppfl.SecureAggregation.Aggregate(values)
    
    // 转换回模型
    aggregated := globalModel.Clone()
    aggregated.Weights[0] = float64(aggregatedValue) / 1000000.0
    
    return aggregated
}

// standardAggregation 标准聚合
func (ppfl *PrivacyPreservingFL) standardAggregation(globalModel *Model, clientModels []*Model, clientWeights []float64) *Model {
    aggregated := NewModel()
    aggregated.Weights = make([]float64, len(clientModels[0].Weights))
    aggregated.Bias = make([]float64, len(clientModels[0].Bias))
    
    // 加权平均
    for i := 0; i < len(clientModels[0].Weights); i++ {
        sum := 0.0
        for j := 0; j < len(clientModels); j++ {
            sum += clientWeights[j] * clientModels[j].Weights[i]
        }
        aggregated.Weights[i] = sum
    }
    
    for i := 0; i < len(clientModels[0].Bias); i++ {
        sum := 0.0
        for j := 0; j < len(clientModels); j++ {
            sum += clientWeights[j] * clientModels[j].Bias[i]
        }
        aggregated.Bias[i] = sum
    }
    
    return aggregated
}

// EvaluatePrivacy 评估隐私保护效果
func (ppfl *PrivacyPreservingFL) EvaluatePrivacy(model *Model, testData [][]float64) *PrivacyMetrics {
    metrics := &PrivacyMetrics{}
    
    // 成员推理攻击评估
    mia := NewMembershipInferenceAttack()
    membershipAccuracy := 0.0
    
    for _, data := range testData {
        output := model.Forward(data)
        isMember := mia.Attack(model, output).(bool)
        
        // 简化评估：假设50%是成员
        if isMember {
            membershipAccuracy += 0.5
        } else {
            membershipAccuracy += 0.5
        }
    }
    
    metrics.MembershipInferenceAccuracy = membershipAccuracy / float64(len(testData))
    
    // 计算隐私保护强度
    if ppfl.Config.UseDP {
        metrics.PrivacyStrength = ppfl.Config.Epsilon
    } else {
        metrics.PrivacyStrength = math.Inf(1)
    }
    
    return metrics
}

// PrivacyMetrics 隐私指标
type PrivacyMetrics struct {
    MembershipInferenceAccuracy float64
    PrivacyStrength             float64
}
```

## 3.6 应用示例

### 3.6.1 隐私保护训练示例

```go
// PrivacyPreservingTrainingExample 隐私保护训练示例
func PrivacyPreservingTrainingExample() {
    // 创建隐私配置
    config := &PrivacyConfig{
        Epsilon:      1.0,
        Delta:        1e-5,
        ClipNorm:     1.0,
        UseDP:        true,
        UseSecureAgg: false,
        Threshold:    3,
        Parties:      5,
    }
    
    // 创建隐私保护联邦学习
    ppfl := NewPrivacyPreservingFL(config)
    
    // 创建模型和客户端
    globalModel := NewModel()
    globalModel.InitializeRandom()
    
    clientModels := make([]*Model, 5)
    clientWeights := make([]float64, 5)
    
    for i := 0; i < 5; i++ {
        clientModels[i] = NewModel()
        clientModels[i].InitializeRandom()
        clientWeights[i] = 0.2
    }
    
    // 隐私保护聚合
    protectedModel := ppfl.TrainWithPrivacy(globalModel, clientModels, clientWeights)
    
    fmt.Printf("Privacy-protected model: %v\n", protectedModel)
    
    // 评估隐私保护效果
    testData := make([][]float64, 100)
    for i := 0; i < 100; i++ {
        testData[i] = make([]float64, 784)
        for j := 0; j < 784; j++ {
            testData[i][j] = rand.Float64()
        }
    }
    
    metrics := ppfl.EvaluatePrivacy(protectedModel, testData)
    fmt.Printf("Privacy metrics: %+v\n", metrics)
}
```

## 3.7 理论证明

### 3.7.1 差分隐私保护

**定理 3.1** (差分隐私保护)
高斯机制提供 ```latex
$(\epsilon, \delta)$
```-差分隐私保护。

**证明**：
通过分析高斯噪声的分布性质，可以证明差分隐私的数学定义。

### 3.7.2 安全聚合正确性

**定理 3.2** (安全聚合正确性)
秘密共享协议在任意 ```latex
$t$
``` 个参与方的情况下可以正确重构秘密。

**证明**：
通过拉格朗日插值定理，可以证明秘密共享的正确性。

## 3.8 总结

联邦学习隐私保护通过差分隐私、安全多方计算等技术，在保护用户数据隐私的同时实现模型训练。这些技术为联邦学习的实际应用提供了安全保障。

---

**参考文献**：

1. Dwork, C. (2006). Differential privacy.
2. McMahan, H. B., Ramage, D., Talwar, K., Zhang, L. (2018). Learning differentially private recurrent language models.
3. Bonawitz, K., Ivanov, V., Kreuter, B., Marcedone, A., McMahan, H. B., Patel, S., ... & Seth, K. (2017). Practical secure aggregation for privacy-preserving machine learning.
