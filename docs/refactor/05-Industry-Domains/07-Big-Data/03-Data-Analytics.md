# 03-数据分析 (Data Analytics)

## 概述

数据分析是大数据系统的核心价值体现，包括统计分析、机器学习、数据挖掘等技术。本章将介绍数据分析的理论基础、算法实现和Go语言应用。

## 目录

1. [理论基础](#1-理论基础)
2. [统计分析](#2-统计分析)
3. [机器学习](#3-机器学习)
4. [数据挖掘](#4-数据挖掘)
5. [预测分析](#5-预测分析)
6. [异常检测](#6-异常检测)
7. [模式识别](#7-模式识别)
8. [可视化分析](#8-可视化分析)

## 1. 理论基础

### 1.1 数据分析模型

**定义 1.1** (数据分析函数)
数据分析函数 ```latex
$f: D \rightarrow A$
``` 将数据集 ```latex
$D$
``` 映射到分析结果 ```latex
$A$
```：

$```latex
$f(D) = \text{analyze}(D, \text{method}, \text{params})$
```$

其中 ```latex
$\text{method}$
``` 是分析方法，```latex
$\text{params}$
``` 是参数集合。

**定理 1.1** (分析一致性)
对于数据集 ```latex
$D$
``` 和分析方法 ```latex
$m$
```，如果 ```latex
$D_1 \subseteq D_2$
```，则：

$```latex
$f(D_1, m) \subseteq f(D_2, m)$
```$

### 1.2 分析框架

```go
// 数据分析接口
type DataAnalyzer interface {
    Analyze(data []DataPoint) (AnalysisResult, error)
    Validate(data []DataPoint) bool
    GetMetrics() map[string]float64
}

// 分析结果
type AnalysisResult struct {
    Method    string                 `json:"method"`
    Metrics   map[string]float64     `json:"metrics"`
    Insights  []Insight              `json:"insights"`
    Model     interface{}            `json:"model,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

// 洞察
type Insight struct {
    Type        string      `json:"type"`
    Description string      `json:"description"`
    Confidence  float64     `json:"confidence"`
    Data        interface{} `json:"data"`
}

// 数据点
type DataPoint struct {
    ID        string                 `json:"id"`
    Features  map[string]interface{} `json:"features"`
    Target    interface{}            `json:"target,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
    Weight    float64                `json:"weight,omitempty"`
}

// 基础分析器
type BaseAnalyzer struct {
    name    string
    metrics map[string]float64
    config  map[string]interface{}
}

func (ba *BaseAnalyzer) Validate(data []DataPoint) bool {
    if len(data) == 0 {
        return false
    }
    
    // 检查数据一致性
    featureCount := len(data[0].Features)
    for _, point := range data {
        if len(point.Features) != featureCount {
            return false
        }
    }
    
    return true
}

func (ba *BaseAnalyzer) GetMetrics() map[string]float64 {
    return ba.metrics
}
```

## 2. 统计分析

### 2.1 描述性统计

**定义 2.1** (描述性统计)
描述性统计是对数据集基本特征的量化描述：

$```latex
$\text{DescriptiveStats}(D) = \{\text{mean}, \text{median}, \text{mode}, \text{variance}, \text{skewness}, \text{kurtosis}\}$
```$

```go
// 描述性统计分析器
type DescriptiveAnalyzer struct {
    BaseAnalyzer
}

func (da *DescriptiveAnalyzer) Analyze(data []DataPoint) (AnalysisResult, error) {
    if !da.Validate(data) {
        return AnalysisResult{}, errors.New("invalid data")
    }
    
    // 提取数值特征
    numericFeatures := da.extractNumericFeatures(data)
    
    // 计算统计量
    stats := make(map[string]map[string]float64)
    for feature, values := range numericFeatures {
        stats[feature] = da.calculateStats(values)
    }
    
    // 生成洞察
    insights := da.generateInsights(stats)
    
    return AnalysisResult{
        Method:    "descriptive_statistics",
        Metrics:   da.aggregateMetrics(stats),
        Insights:  insights,
        Timestamp: time.Now(),
    }, nil
}

func (da *DescriptiveAnalyzer) extractNumericFeatures(data []DataPoint) map[string][]float64 {
    features := make(map[string][]float64)
    
    for _, point := range data {
        for name, value := range point.Features {
            if num, ok := value.(float64); ok {
                features[name] = append(features[name], num)
            }
        }
    }
    
    return features
}

func (da *DescriptiveAnalyzer) calculateStats(values []float64) map[string]float64 {
    if len(values) == 0 {
        return make(map[string]float64)
    }
    
    stats := make(map[string]float64)
    
    // 基本统计量
    stats["count"] = float64(len(values))
    stats["sum"] = da.sum(values)
    stats["mean"] = stats["sum"] / stats["count"]
    stats["min"] = da.min(values)
    stats["max"] = da.max(values)
    
    // 中位数
    sorted := make([]float64, len(values))
    copy(sorted, values)
    sort.Float64s(sorted)
    stats["median"] = da.median(sorted)
    
    // 方差和标准差
    variance := da.variance(values, stats["mean"])
    stats["variance"] = variance
    stats["std_dev"] = math.Sqrt(variance)
    
    // 偏度和峰度
    stats["skewness"] = da.skewness(values, stats["mean"], stats["std_dev"])
    stats["kurtosis"] = da.kurtosis(values, stats["mean"], stats["std_dev"])
    
    return stats
}

func (da *DescriptiveAnalyzer) sum(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum
}

func (da *DescriptiveAnalyzer) min(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    min := values[0]
    for _, v := range values {
        if v < min {
            min = v
        }
    }
    return min
}

func (da *DescriptiveAnalyzer) max(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    max := values[0]
    for _, v := range values {
        if v > max {
            max = v
        }
    }
    return max
}

func (da *DescriptiveAnalyzer) median(sorted []float64) float64 {
    n := len(sorted)
    if n == 0 {
        return 0
    }
    
    if n%2 == 0 {
        return (sorted[n/2-1] + sorted[n/2]) / 2
    }
    return sorted[n/2]
}

func (da *DescriptiveAnalyzer) variance(values []float64, mean float64) float64 {
    if len(values) <= 1 {
        return 0
    }
    
    sumSquaredDiff := 0.0
    for _, v := range values {
        diff := v - mean
        sumSquaredDiff += diff * diff
    }
    
    return sumSquaredDiff / float64(len(values)-1)
}

func (da *DescriptiveAnalyzer) skewness(values []float64, mean, stdDev float64) float64 {
    if stdDev == 0 {
        return 0
    }
    
    n := float64(len(values))
    sum := 0.0
    
    for _, v := range values {
        z := (v - mean) / stdDev
        sum += z * z * z
    }
    
    return (sum / n) * math.Sqrt(n*(n-1)) / (n - 2)
}

func (da *DescriptiveAnalyzer) kurtosis(values []float64, mean, stdDev float64) float64 {
    if stdDev == 0 {
        return 0
    }
    
    n := float64(len(values))
    sum := 0.0
    
    for _, v := range values {
        z := (v - mean) / stdDev
        sum += z * z * z * z
    }
    
    return (sum/n - 3) * (n*n-1) / ((n-1)*(n-2)*(n-3))
}
```

### 2.2 推断性统计

```go
// 推断性统计分析器
type InferentialAnalyzer struct {
    BaseAnalyzer
    confidenceLevel float64
}

func (ia *InferentialAnalyzer) Analyze(data []DataPoint) (AnalysisResult, error) {
    if !ia.Validate(data) {
        return AnalysisResult{}, errors.New("invalid data")
    }
    
    numericFeatures := ia.extractNumericFeatures(data)
    results := make(map[string]interface{})
    
    for feature, values := range numericFeatures {
        // 置信区间
        ci := ia.confidenceInterval(values)
        
        // 假设检验
        test := ia.hypothesisTest(values)
        
        results[feature] = map[string]interface{}{
            "confidence_interval": ci,
            "hypothesis_test":     test,
        }
    }
    
    return AnalysisResult{
        Method:    "inferential_statistics",
        Metrics:   ia.calculateInferentialMetrics(results),
        Insights:  ia.generateInferentialInsights(results),
        Timestamp: time.Now(),
    }, nil
}

func (ia *InferentialAnalyzer) confidenceInterval(values []float64) map[string]float64 {
    if len(values) < 2 {
        return map[string]float64{"lower": 0, "upper": 0}
    }
    
    mean := ia.mean(values)
    stdDev := ia.stdDev(values)
    n := float64(len(values))
    
    // t分布临界值（简化，实际应查表）
    tValue := 1.96 // 95%置信水平
    
    margin := tValue * stdDev / math.Sqrt(n)
    
    return map[string]float64{
        "lower": mean - margin,
        "upper": mean + margin,
        "mean":  mean,
        "margin": margin,
    }
}

func (ia *InferentialAnalyzer) hypothesisTest(values []float64) map[string]interface{} {
    if len(values) < 2 {
        return map[string]interface{}{"p_value": 1.0, "significant": false}
    }
    
    // 单样本t检验（检验均值是否为0）
    mean := ia.mean(values)
    stdDev := ia.stdDev(values)
    n := float64(len(values))
    
    tStat := mean / (stdDev / math.Sqrt(n))
    pValue := ia.tTestPValue(tStat, n-1)
    
    return map[string]interface{}{
        "t_statistic": tStat,
        "p_value":     pValue,
        "significant": pValue < 0.05,
        "null_hypothesis": "mean = 0",
        "alternative_hypothesis": "mean ≠ 0",
    }
}

func (ia *InferentialAnalyzer) mean(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (ia *InferentialAnalyzer) stdDev(values []float64) float64 {
    mean := ia.mean(values)
    sumSquaredDiff := 0.0
    for _, v := range values {
        diff := v - mean
        sumSquaredDiff += diff * diff
    }
    return math.Sqrt(sumSquaredDiff / float64(len(values)-1))
}

func (ia *InferentialAnalyzer) tTestPValue(tStat, df float64) float64 {
    // 简化实现，实际应使用t分布表或数值方法
    // 这里使用正态分布近似
    return 2 * (1 - ia.normalCDF(math.Abs(tStat)))
}

func (ia *InferentialAnalyzer) normalCDF(x float64) float64 {
    // 标准正态分布累积分布函数近似
    return 0.5 * (1 + math.Erf(x/math.Sqrt(2)))
}
```

## 3. 机器学习

### 3.1 监督学习

**定义 3.1** (监督学习)
监督学习是从标记数据中学习映射函数 ```latex
$f: X \rightarrow Y$
```：

$```latex
$\text{SupervisedLearning}(D) = \arg\min_f \sum_{(x,y) \in D} L(f(x), y)$
```$

其中 ```latex
$L$
``` 是损失函数。

```go
// 监督学习接口
type SupervisedLearner interface {
    Train(data []DataPoint) error
    Predict(features map[string]interface{}) (interface{}, error)
    Evaluate(testData []DataPoint) (EvaluationResult, error)
}

// 线性回归
type LinearRegression struct {
    weights map[string]float64
    bias    float64
    learningRate float64
    epochs       int
}

func (lr *LinearRegression) Train(data []DataPoint) error {
    if len(data) == 0 {
        return errors.New("no training data")
    }
    
    // 初始化权重
    lr.initializeWeights(data[0].Features)
    
    // 梯度下降训练
    for epoch := 0; epoch < lr.epochs; epoch++ {
        for _, point := range data {
            lr.updateWeights(point)
        }
    }
    
    return nil
}

func (lr *LinearRegression) initializeWeights(features map[string]interface{}) {
    lr.weights = make(map[string]float64)
    for feature := range features {
        lr.weights[feature] = rand.Float64() * 0.1
    }
    lr.bias = rand.Float64() * 0.1
}

func (lr *LinearRegression) updateWeights(point DataPoint) {
    // 前向传播
    prediction := lr.predict(point.Features)
    target := point.Target.(float64)
    
    // 计算梯度
    error := prediction - target
    
    // 更新权重
    for feature, value := range point.Features {
        if num, ok := value.(float64); ok {
            gradient := error * num
            lr.weights[feature] -= lr.learningRate * gradient
        }
    }
    
    // 更新偏置
    lr.bias -= lr.learningRate * error
}

func (lr *LinearRegression) predict(features map[string]interface{}) float64 {
    prediction := lr.bias
    
    for feature, value := range features {
        if weight, exists := lr.weights[feature]; exists {
            if num, ok := value.(float64); ok {
                prediction += weight * num
            }
        }
    }
    
    return prediction
}

func (lr *LinearRegression) Predict(features map[string]interface{}) (interface{}, error) {
    return lr.predict(features), nil
}

func (lr *LinearRegression) Evaluate(testData []DataPoint) (EvaluationResult, error) {
    var predictions []float64
    var targets []float64
    
    for _, point := range testData {
        pred, err := lr.Predict(point.Features)
        if err != nil {
            return EvaluationResult{}, err
        }
        
        predictions = append(predictions, pred.(float64))
        targets = append(targets, point.Target.(float64))
    }
    
    return EvaluationResult{
        MSE:      lr.meanSquaredError(predictions, targets),
        MAE:      lr.meanAbsoluteError(predictions, targets),
        R2:       lr.rSquared(predictions, targets),
        Predictions: predictions,
        Targets:     targets,
    }, nil
}

type EvaluationResult struct {
    MSE         float64   `json:"mse"`
    MAE         float64   `json:"mae"`
    R2          float64   `json:"r2"`
    Predictions []float64 `json:"predictions"`
    Targets     []float64 `json:"targets"`
}

func (lr *LinearRegression) meanSquaredError(predictions, targets []float64) float64 {
    if len(predictions) != len(targets) {
        return 0
    }
    
    sum := 0.0
    for i := range predictions {
        diff := predictions[i] - targets[i]
        sum += diff * diff
    }
    
    return sum / float64(len(predictions))
}

func (lr *LinearRegression) meanAbsoluteError(predictions, targets []float64) float64 {
    if len(predictions) != len(targets) {
        return 0
    }
    
    sum := 0.0
    for i := range predictions {
        sum += math.Abs(predictions[i] - targets[i])
    }
    
    return sum / float64(len(predictions))
}

func (lr *LinearRegression) rSquared(predictions, targets []float64) float64 {
    if len(predictions) != len(targets) {
        return 0
    }
    
    // 计算目标均值
    targetMean := 0.0
    for _, target := range targets {
        targetMean += target
    }
    targetMean /= float64(len(targets))
    
    // 计算总平方和
    totalSS := 0.0
    for _, target := range targets {
        diff := target - targetMean
        totalSS += diff * diff
    }
    
    // 计算残差平方和
    residualSS := 0.0
    for i := range predictions {
        diff := targets[i] - predictions[i]
        residualSS += diff * diff
    }
    
    if totalSS == 0 {
        return 0
    }
    
    return 1 - (residualSS / totalSS)
}
```

### 3.2 无监督学习

```go
// 无监督学习接口
type UnsupervisedLearner interface {
    Train(data []DataPoint) error
    Cluster(data []DataPoint) ([]int, error)
    GetCentroids() []map[string]float64
}

// K-means聚类
type KMeans struct {
    k           int
    centroids   []map[string]float64
    maxIterations int
    tolerance   float64
}

func (km *KMeans) Train(data []DataPoint) error {
    if len(data) == 0 {
        return errors.New("no training data")
    }
    
    // 初始化聚类中心
    km.initializeCentroids(data)
    
    // 迭代优化
    for iteration := 0; iteration < km.maxIterations; iteration++ {
        // 分配数据点到最近的中心
        assignments := km.assignToClusters(data)
        
        // 更新聚类中心
        newCentroids := km.updateCentroids(data, assignments)
        
        // 检查收敛
        if km.hasConverged(newCentroids) {
            break
        }
        
        km.centroids = newCentroids
    }
    
    return nil
}

func (km *KMeans) initializeCentroids(data []DataPoint) {
    km.centroids = make([]map[string]float64, km.k)
    
    // 随机选择初始中心
    for i := 0; i < km.k; i++ {
        randomIndex := rand.Intn(len(data))
        point := data[randomIndex]
        
        km.centroids[i] = make(map[string]float64)
        for feature, value := range point.Features {
            if num, ok := value.(float64); ok {
                km.centroids[i][feature] = num
            }
        }
    }
}

func (km *KMeans) assignToClusters(data []DataPoint) []int {
    assignments := make([]int, len(data))
    
    for i, point := range data {
        minDistance := math.MaxFloat64
        bestCluster := 0
        
        for j, centroid := range km.centroids {
            distance := km.euclideanDistance(point.Features, centroid)
            if distance < minDistance {
                minDistance = distance
                bestCluster = j
            }
        }
        
        assignments[i] = bestCluster
    }
    
    return assignments
}

func (km *KMeans) euclideanDistance(features map[string]interface{}, centroid map[string]float64) float64 {
    sum := 0.0
    
    for feature, value := range features {
        if num, ok := value.(float64); ok {
            if centroidValue, exists := centroid[feature]; exists {
                diff := num - centroidValue
                sum += diff * diff
            }
        }
    }
    
    return math.Sqrt(sum)
}

func (km *KMeans) updateCentroids(data []DataPoint, assignments []int) []map[string]float64 {
    newCentroids := make([]map[string]float64, km.k)
    clusterSizes := make([]int, km.k)
    
    // 初始化新中心
    for i := range newCentroids {
        newCentroids[i] = make(map[string]float64)
    }
    
    // 累加特征值
    for i, point := range data {
        cluster := assignments[i]
        clusterSizes[cluster]++
        
        for feature, value := range point.Features {
            if num, ok := value.(float64); ok {
                newCentroids[cluster][feature] += num
            }
        }
    }
    
    // 计算平均值
    for i := range newCentroids {
        if clusterSizes[i] > 0 {
            for feature := range newCentroids[i] {
                newCentroids[i][feature] /= float64(clusterSizes[i])
            }
        }
    }
    
    return newCentroids
}

func (km *KMeans) hasConverged(newCentroids []map[string]float64) bool {
    for i, newCentroid := range newCentroids {
        if i >= len(km.centroids) {
            return false
        }
        
        oldCentroid := km.centroids[i]
        for feature, newValue := range newCentroid {
            if oldValue, exists := oldCentroid[feature]; exists {
                if math.Abs(newValue-oldValue) > km.tolerance {
                    return false
                }
            }
        }
    }
    
    return true
}

func (km *KMeans) Cluster(data []DataPoint) ([]int, error) {
    if len(km.centroids) == 0 {
        return nil, errors.New("model not trained")
    }
    
    return km.assignToClusters(data), nil
}

func (km *KMeans) GetCentroids() []map[string]float64 {
    return km.centroids
}
```

## 4. 数据挖掘

### 4.1 关联规则挖掘

**定义 4.1** (关联规则)
关联规则是形如 ```latex
$X \rightarrow Y$
``` 的规则，其中 ```latex
$X$
``` 和 ```latex
$Y$
``` 是项集：

$```latex
$\text{Support}(X) = \frac{|\{t \in T : X \subseteq t\}|}{|T|}$
```$

$```latex
$\text{Confidence}(X \rightarrow Y) = \frac{\text{Support}(X \cup Y)}{\text{Support}(X)}$
```$

```go
// 关联规则挖掘器
type AssociationRuleMiner struct {
    minSupport    float64
    minConfidence float64
    rules         []AssociationRule
}

type AssociationRule struct {
    Antecedent    []string  `json:"antecedent"`
    Consequent    []string  `json:"consequent"`
    Support       float64   `json:"support"`
    Confidence    float64   `json:"confidence"`
    Lift          float64   `json:"lift"`
}

type Transaction struct {
    ID   string   `json:"id"`
    Items []string `json:"items"`
}

func (arm *AssociationRuleMiner) MineRules(transactions []Transaction) ([]AssociationRule, error) {
    // 计算频繁项集
    frequentItemsets := arm.findFrequentItemsets(transactions)
    
    // 生成关联规则
    rules := make([]AssociationRule, 0)
    
    for itemset := range frequentItemsets {
        if len(itemset) < 2 {
            continue
        }
        
        // 生成所有可能的规则
        subsets := arm.generateSubsets(itemset)
        for _, antecedent := range subsets {
            if len(antecedent) == 0 || len(antecedent) == len(itemset) {
                continue
            }
            
            consequent := arm.subtract(itemset, antecedent)
            
            rule := arm.createRule(antecedent, consequent, transactions)
            if rule.Confidence >= arm.minConfidence {
                rules = append(rules, rule)
            }
        }
    }
    
    arm.rules = rules
    return rules, nil
}

func (arm *AssociationRuleMiner) findFrequentItemsets(transactions []Transaction) map[string]float64 {
    itemCounts := make(map[string]int)
    totalTransactions := len(transactions)
    
    // 计算单项支持度
    for _, transaction := range transactions {
        for _, item := range transaction.Items {
            itemCounts[item]++
        }
    }
    
    // 筛选频繁单项
    frequentItems := make([]string, 0)
    for item, count := range itemCounts {
        support := float64(count) / float64(totalTransactions)
        if support >= arm.minSupport {
            frequentItems = append(frequentItems, item)
        }
    }
    
    // 生成频繁项集（简化实现，只考虑2项集）
    frequentItemsets := make(map[string]float64)
    
    for i := 0; i < len(frequentItems); i++ {
        for j := i + 1; j < len(frequentItems); j++ {
            itemset := []string{frequentItems[i], frequentItems[j]}
            support := arm.calculateSupport(itemset, transactions)
            if support >= arm.minSupport {
                key := strings.Join(itemset, ",")
                frequentItemsets[key] = support
            }
        }
    }
    
    return frequentItemsets
}

func (arm *AssociationRuleMiner) calculateSupport(itemset []string, transactions []Transaction) float64 {
    count := 0
    total := len(transactions)
    
    for _, transaction := range transactions {
        if arm.containsAll(transaction.Items, itemset) {
            count++
        }
    }
    
    return float64(count) / float64(total)
}

func (arm *AssociationRuleMiner) containsAll(transactionItems, itemset []string) bool {
    for _, item := range itemset {
        found := false
        for _, transactionItem := range transactionItems {
            if item == transactionItem {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}

func (arm *AssociationRuleMiner) generateSubsets(itemset []string) [][]string {
    n := len(itemset)
    subsets := make([][]string, 0)
    
    // 生成所有子集
    for i := 1; i < (1 << n); i++ {
        subset := make([]string, 0)
        for j := 0; j < n; j++ {
            if i&(1<<j) != 0 {
                subset = append(subset, itemset[j])
            }
        }
        subsets = append(subsets, subset)
    }
    
    return subsets
}

func (arm *AssociationRuleMiner) subtract(itemset, antecedent []string) []string {
    result := make([]string, 0)
    
    for _, item := range itemset {
        found := false
        for _, antItem := range antecedent {
            if item == antItem {
                found = true
                break
            }
        }
        if !found {
            result = append(result, item)
        }
    }
    
    return result
}

func (arm *AssociationRuleMiner) createRule(antecedent, consequent []string, transactions []Transaction) AssociationRule {
    antecedentSupport := arm.calculateSupport(antecedent, transactions)
    consequentSupport := arm.calculateSupport(consequent, transactions)
    
    union := append(antecedent, consequent...)
    unionSupport := arm.calculateSupport(union, transactions)
    
    confidence := 0.0
    if antecedentSupport > 0 {
        confidence = unionSupport / antecedentSupport
    }
    
    lift := 0.0
    if antecedentSupport > 0 && consequentSupport > 0 {
        lift = unionSupport / (antecedentSupport * consequentSupport)
    }
    
    return AssociationRule{
        Antecedent: antecedent,
        Consequent: consequent,
        Support:    unionSupport,
        Confidence: confidence,
        Lift:       lift,
    }
}
```

## 5. 预测分析

### 5.1 时间序列预测

```go
// 时间序列预测器
type TimeSeriesPredictor struct {
    method    string
    window    int
    model     interface{}
}

// 移动平均预测
type MovingAveragePredictor struct {
    window int
    values []float64
}

func (map *MovingAveragePredictor) Train(data []DataPoint) error {
    if len(data) < map.window {
        return errors.New("insufficient data for window size")
    }
    
    map.values = make([]float64, len(data))
    for i, point := range data {
        if target, ok := point.Target.(float64); ok {
            map.values[i] = target
        }
    }
    
    return nil
}

func (map *MovingAveragePredictor) Predict(steps int) []float64 {
    if len(map.values) < map.window {
        return nil
    }
    
    predictions := make([]float64, steps)
    
    for i := 0; i < steps; i++ {
        // 计算移动平均
        sum := 0.0
        start := len(map.values) - map.window + i
        if start < 0 {
            start = 0
        }
        
        for j := start; j < len(map.values); j++ {
            sum += map.values[j]
        }
        
        predictions[i] = sum / float64(map.window)
    }
    
    return predictions
}

// 指数平滑预测
type ExponentialSmoothingPredictor struct {
    alpha float64
    values []float64
    lastForecast float64
}

func (esp *ExponentialSmoothingPredictor) Train(data []DataPoint) error {
    if len(data) == 0 {
        return errors.New("no training data")
    }
    
    esp.values = make([]float64, len(data))
    for i, point := range data {
        if target, ok := point.Target.(float64); ok {
            esp.values[i] = target
        }
    }
    
    // 初始化预测值
    if len(esp.values) > 0 {
        esp.lastForecast = esp.values[0]
    }
    
    return nil
}

func (esp *ExponentialSmoothingPredictor) Predict(steps int) []float64 {
    predictions := make([]float64, steps)
    
    currentForecast := esp.lastForecast
    
    for i := 0; i < steps; i++ {
        predictions[i] = currentForecast
        // 更新预测值（这里简化处理）
        currentForecast = esp.alpha * esp.values[len(esp.values)-1] + (1-esp.alpha) * currentForecast
    }
    
    return predictions
}
```

## 6. 异常检测

### 6.1 统计异常检测

```go
// 异常检测器
type AnomalyDetector interface {
    Train(data []DataPoint) error
    Detect(data []DataPoint) ([]bool, error)
    GetThreshold() float64
}

// Z-score异常检测
type ZScoreAnomalyDetector struct {
    mean     float64
    stdDev   float64
    threshold float64
}

func (zad *ZScoreAnomalyDetector) Train(data []DataPoint) error {
    if len(data) == 0 {
        return errors.New("no training data")
    }
    
    // 提取数值特征
    values := make([]float64, 0)
    for _, point := range data {
        for _, value := range point.Features {
            if num, ok := value.(float64); ok {
                values = append(values, num)
            }
        }
    }
    
    if len(values) == 0 {
        return errors.New("no numeric features")
    }
    
    // 计算均值和标准差
    zad.mean = zad.calculateMean(values)
    zad.stdDev = zad.calculateStdDev(values, zad.mean)
    
    return nil
}

func (zad *ZScoreAnomalyDetector) Detect(data []DataPoint) ([]bool, error) {
    if zad.stdDev == 0 {
        return nil, errors.New("model not trained or zero standard deviation")
    }
    
    anomalies := make([]bool, len(data))
    
    for i, point := range data {
        for _, value := range point.Features {
            if num, ok := value.(float64); ok {
                zScore := math.Abs((num - zad.mean) / zad.stdDev)
                if zScore > zad.threshold {
                    anomalies[i] = true
                    break
                }
            }
        }
    }
    
    return anomalies, nil
}

func (zad *ZScoreAnomalyDetector) calculateMean(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (zad *ZScoreAnomalyDetector) calculateStdDev(values []float64, mean float64) float64 {
    if len(values) <= 1 {
        return 0
    }
    
    sumSquaredDiff := 0.0
    for _, v := range values {
        diff := v - mean
        sumSquaredDiff += diff * diff
    }
    
    return math.Sqrt(sumSquaredDiff / float64(len(values)-1))
}

func (zad *ZScoreAnomalyDetector) GetThreshold() float64 {
    return zad.threshold
}
```

## 7. 模式识别

### 7.1 序列模式识别

```go
// 序列模式识别器
type SequencePatternRecognizer struct {
    patterns map[string]int
    minSupport float64
}

type Sequence struct {
    ID   string   `json:"id"`
    Items []string `json:"items"`
}

func (spr *SequencePatternRecognizer) FindPatterns(sequences []Sequence) (map[string]float64, error) {
    if len(sequences) == 0 {
        return nil, errors.New("no sequences provided")
    }
    
    // 生成所有可能的子序列
    allSubsequences := make(map[string]int)
    
    for _, sequence := range sequences {
        subsequences := spr.generateSubsequences(sequence.Items)
        for _, subsequence := range subsequences {
            key := strings.Join(subsequence, ",")
            allSubsequences[key]++
        }
    }
    
    // 筛选频繁模式
    frequentPatterns := make(map[string]float64)
    totalSequences := float64(len(sequences))
    
    for pattern, count := range allSubsequences {
        support := float64(count) / totalSequences
        if support >= spr.minSupport {
            frequentPatterns[pattern] = support
        }
    }
    
    return frequentPatterns, nil
}

func (spr *SequencePatternRecognizer) generateSubsequences(items []string) [][]string {
    subsequences := make([][]string, 0)
    
    // 生成所有长度的子序列
    for length := 1; length <= len(items); length++ {
        for start := 0; start <= len(items)-length; start++ {
            subsequence := items[start : start+length]
            subsequences = append(subsequences, subsequence)
        }
    }
    
    return subsequences
}
```

## 8. 可视化分析

### 8.1 数据可视化

```go
// 数据可视化器
type DataVisualizer interface {
    CreateChart(data []DataPoint, chartType string) (Chart, error)
    ExportChart(chart Chart, format string) ([]byte, error)
}

type Chart struct {
    Type    string                 `json:"type"`
    Data    interface{}            `json:"data"`
    Options map[string]interface{} `json:"options"`
}

// 基础可视化器
type BaseVisualizer struct {
    defaultOptions map[string]interface{}
}

func (bv *BaseVisualizer) CreateChart(data []DataPoint, chartType string) (Chart, error) {
    switch chartType {
    case "scatter":
        return bv.createScatterChart(data)
    case "line":
        return bv.createLineChart(data)
    case "bar":
        return bv.createBarChart(data)
    case "histogram":
        return bv.createHistogramChart(data)
    default:
        return Chart{}, errors.New("unsupported chart type")
    }
}

func (bv *BaseVisualizer) createScatterChart(data []DataPoint) (Chart, error) {
    points := make([]map[string]interface{}, 0)
    
    for _, point := range data {
        // 提取两个数值特征作为x和y坐标
        x, y := bv.extractXYCoordinates(point.Features)
        if x != nil && y != nil {
            points = append(points, map[string]interface{}{
                "x": x,
                "y": y,
                "label": point.ID,
            })
        }
    }
    
    return Chart{
        Type: "scatter",
        Data: points,
        Options: map[string]interface{}{
            "title": "Scatter Plot",
            "xLabel": "X Axis",
            "yLabel": "Y Axis",
        },
    }, nil
}

func (bv *BaseVisualizer) extractXYCoordinates(features map[string]interface{}) (x, y interface{}) {
    var xVal, yVal interface{}
    count := 0
    
    for _, value := range features {
        if count == 0 {
            xVal = value
            count++
        } else if count == 1 {
            yVal = value
            break
        }
    }
    
    return xVal, yVal
}

func (bv *BaseVisualizer) ExportChart(chart Chart, format string) ([]byte, error) {
    switch format {
    case "json":
        return json.Marshal(chart)
    case "csv":
        return bv.exportToCSV(chart)
    default:
        return nil, errors.New("unsupported export format")
    }
}

func (bv *BaseVisualizer) exportToCSV(chart Chart) ([]byte, error) {
    var buffer bytes.Buffer
    writer := csv.NewWriter(&buffer)
    
    // 写入标题
    if err := writer.Write([]string{"Type", "Data", "Options"}); err != nil {
        return nil, err
    }
    
    // 写入数据
    dataStr := fmt.Sprintf("%v", chart.Data)
    optionsStr := fmt.Sprintf("%v", chart.Options)
    if err := writer.Write([]string{chart.Type, dataStr, optionsStr}); err != nil {
        return nil, err
    }
    
    writer.Flush()
    return buffer.Bytes(), nil
}
```

## 总结

本章详细介绍了数据分析的核心技术，包括：

1. **理论基础**：数据分析函数、分析框架
2. **统计分析**：描述性统计、推断性统计
3. **机器学习**：监督学习、无监督学习
4. **数据挖掘**：关联规则挖掘
5. **预测分析**：时间序列预测
6. **异常检测**：统计异常检测
7. **模式识别**：序列模式识别
8. **可视化分析**：数据可视化

这些技术为构建智能、高效的数据分析系统提供了完整的理论基础和实现方案。

---

**相关链接**：

- [01-数据摄入](../01-Data-Ingestion.md)
- [02-数据处理](../02-Data-Processing.md)
- [04-数据存储](../04-Data-Storage.md)
- [05-数据可视化](../05-Data-Visualization.md)
- [06-机器学习](../06-Machine-Learning.md)
- [07-实时计算](../07-Real-Time-Computing.md)
- [08-数据治理](../08-Data-Governance.md)
