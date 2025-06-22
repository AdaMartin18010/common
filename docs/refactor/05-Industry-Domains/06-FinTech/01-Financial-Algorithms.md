# 1. 金融算法 (Financial Algorithms)

## 1.1 金融算法理论基础

### 1.1.1 金融数学基础

**定义 1.1** (金融算法): 金融算法是一个四元组 ```latex
$\mathcal{F} = (M, P, R, T)$
```，其中：

- ```latex
$M$
``` 是市场模型集合 (Market Models)
- ```latex
$P$
``` 是定价模型集合 (Pricing Models)
- ```latex
$R$
``` 是风险模型集合 (Risk Models)
- ```latex
$T$
``` 是交易模型集合 (Trading Models)

**金融计算核心**:

```latex
\text{FinancialCalculation} = \text{TimeValue} \times \text{RiskMetrics} \times \text{PricingModels} \times \text{PortfolioTheory}
```

### 1.1.2 时间价值理论

**定义 1.2** (现值): 现值 ```latex
$PV$
``` 是未来现金流 ```latex
$CF_t$
``` 在当前时点的价值：

$```latex
$PV = \sum_{t=1}^{n} \frac{CF_t}{(1 + r)^t}$
```$

其中 ```latex
$r$
``` 是折现率，```latex
$n$
``` 是期数。

**定义 1.3** (终值): 终值 ```latex
$FV$
``` 是当前投资在未来时点的价值：

$```latex
$FV = PV \times (1 + r)^n$
```$

## 1.2 Go语言金融算法实现

### 1.2.1 时间价值计算

```go
// FinancialCalculator 金融计算器
type FinancialCalculator struct {
    precision int
}

// NewFinancialCalculator 创建金融计算器
func NewFinancialCalculator(precision int) *FinancialCalculator {
    return &FinancialCalculator{
        precision: precision,
    }
}

// PresentValue 计算现值
func (fc *FinancialCalculator) PresentValue(cashFlows []float64, rate float64) float64 {
    var pv float64
    for i, cf := range cashFlows {
        pv += cf / math.Pow(1+rate, float64(i+1))
    }
    return fc.round(pv)
}

// FutureValue 计算终值
func (fc *FinancialCalculator) FutureValue(presentValue, rate float64, periods int) float64 {
    fv := presentValue * math.Pow(1+rate, float64(periods))
    return fc.round(fv)
}

// NetPresentValue 计算净现值
func (fc *FinancialCalculator) NetPresentValue(initialInvestment float64, cashFlows []float64, rate float64) float64 {
    pv := fc.PresentValue(cashFlows, rate)
    npv := pv - initialInvestment
    return fc.round(npv)
}

// InternalRateOfReturn 计算内部收益率
func (fc *FinancialCalculator) InternalRateOfReturn(initialInvestment float64, cashFlows []float64) float64 {
    // 使用牛顿法求解IRR
    const maxIterations = 100
    const tolerance = 1e-6
    
    // 初始猜测值
    rate := 0.1
    
    for i := 0; i < maxIterations; i++ {
        npv := fc.NetPresentValue(initialInvestment, cashFlows, rate)
        
        if math.Abs(npv) < tolerance {
            return fc.round(rate)
        }
        
        // 计算导数
        derivative := fc.calculateNPVDerivative(cashFlows, rate)
        if math.Abs(derivative) < tolerance {
            break
        }
        
        // 牛顿法迭代
        rate = rate - npv/derivative
    }
    
    return fc.round(rate)
}

// calculateNPVDerivative 计算NPV的导数
func (fc *FinancialCalculator) calculateNPVDerivative(cashFlows []float64, rate float64) float64 {
    var derivative float64
    for i, cf := range cashFlows {
        derivative -= float64(i+1) * cf / math.Pow(1+rate, float64(i+2))
    }
    return derivative
}

// Payment 计算等额本息还款
func (fc *FinancialCalculator) Payment(principal, rate float64, periods int) float64 {
    if rate == 0 {
        return principal / float64(periods)
    }
    
    monthlyRate := rate / 12
    payment := principal * (monthlyRate * math.Pow(1+monthlyRate, float64(periods))) / 
               (math.Pow(1+monthlyRate, float64(periods)) - 1)
    
    return fc.round(payment)
}

// AmortizationSchedule 生成还款计划
func (fc *FinancialCalculator) AmortizationSchedule(principal, rate float64, periods int) []AmortizationRow {
    payment := fc.Payment(principal, rate, periods)
    schedule := make([]AmortizationRow, periods)
    
    remainingBalance := principal
    monthlyRate := rate / 12
    
    for i := 0; i < periods; i++ {
        interest := remainingBalance * monthlyRate
        principal := payment - interest
        remainingBalance -= principal
        
        schedule[i] = AmortizationRow{
            Period:           i + 1,
            Payment:          fc.round(payment),
            Principal:        fc.round(principal),
            Interest:         fc.round(interest),
            RemainingBalance: fc.round(remainingBalance),
        }
    }
    
    return schedule
}

// AmortizationRow 还款计划行
type AmortizationRow struct {
    Period           int     `json:"period"`
    Payment          float64 `json:"payment"`
    Principal        float64 `json:"principal"`
    Interest         float64 `json:"interest"`
    RemainingBalance float64 `json:"remaining_balance"`
}

// round 四舍五入到指定精度
func (fc *FinancialCalculator) round(value float64) float64 {
    return math.Round(value*math.Pow(10, float64(fc.precision))) / math.Pow(10, float64(fc.precision))
}
```

### 1.2.2 债券定价

```go
// Bond 债券结构
type Bond struct {
    FaceValue    float64   `json:"face_value"`
    CouponRate   float64   `json:"coupon_rate"`
    MaturityDate time.Time `json:"maturity_date"`
    IssueDate    time.Time `json:"issue_date"`
    Frequency    int       `json:"frequency"` // 年付息次数
}

// BondCalculator 债券计算器
type BondCalculator struct {
    calculator *FinancialCalculator
}

// NewBondCalculator 创建债券计算器
func NewBondCalculator() *BondCalculator {
    return &BondCalculator{
        calculator: NewFinancialCalculator(6),
    }
}

// Price 计算债券价格
func (bc *BondCalculator) Price(bond *Bond, yield float64, settlementDate time.Time) float64 {
    // 计算到下一个付息日的天数
    nextCouponDate := bc.getNextCouponDate(bond, settlementDate)
    daysToNextCoupon := nextCouponDate.Sub(settlementDate).Hours() / 24
    
    // 计算付息周期
    couponPeriod := 365.0 / float64(bond.Frequency)
    
    // 计算应计利息
    accruedInterest := bc.calculateAccruedInterest(bond, settlementDate)
    
    // 计算现金流
    cashFlows := bc.calculateCashFlows(bond, settlementDate)
    
    // 计算现值
    presentValue := bc.calculator.PresentValue(cashFlows, yield/float64(bond.Frequency))
    
    // 债券价格 = 现值 - 应计利息
    price := presentValue - accruedInterest
    
    return bc.calculator.round(price)
}

// Yield 计算债券收益率
func (bc *BondCalculator) Yield(bond *Bond, price float64, settlementDate time.Time) float64 {
    // 使用二分法求解收益率
    const maxIterations = 100
    const tolerance = 1e-6
    
    low := -0.5
    high := 2.0
    
    for i := 0; i < maxIterations; i++ {
        mid := (low + high) / 2
        calculatedPrice := bc.Price(bond, mid, settlementDate)
        
        if math.Abs(calculatedPrice-price) < tolerance {
            return bc.calculator.round(mid)
        }
        
        if calculatedPrice > price {
            high = mid
        } else {
            low = mid
        }
    }
    
    return bc.calculator.round((low + high) / 2)
}

// Duration 计算久期
func (bc *BondCalculator) Duration(bond *Bond, yield float64, settlementDate time.Time) float64 {
    cashFlows := bc.calculateCashFlows(bond, settlementDate)
    price := bc.Price(bond, yield, settlementDate)
    
    var duration float64
    for i, cf := range cashFlows {
        timeToCashFlow := float64(i+1) / float64(bond.Frequency)
        pv := cf / math.Pow(1+yield/float64(bond.Frequency), float64(i+1))
        duration += timeToCashFlow * pv
    }
    
    duration /= price
    return bc.calculator.round(duration)
}

// ModifiedDuration 计算修正久期
func (bc *BondCalculator) ModifiedDuration(bond *Bond, yield float64, settlementDate time.Time) float64 {
    duration := bc.Duration(bond, yield, settlementDate)
    modifiedDuration := duration / (1 + yield/float64(bond.Frequency))
    return bc.calculator.round(modifiedDuration)
}

// getNextCouponDate 获取下一个付息日
func (bc *BondCalculator) getNextCouponDate(bond *Bond, settlementDate time.Time) time.Time {
    // 简化实现，实际需要考虑节假日等
    daysPerPeriod := 365 / bond.Frequency
    daysSinceIssue := int(settlementDate.Sub(bond.IssueDate).Hours() / 24)
    periodsSinceIssue := daysSinceIssue / daysPerPeriod
    
    nextPeriod := periodsSinceIssue + 1
    daysToNext := nextPeriod*daysPerPeriod - daysSinceIssue
    
    return settlementDate.AddDate(0, 0, daysToNext)
}

// calculateAccruedInterest 计算应计利息
func (bc *BondCalculator) calculateAccruedInterest(bond *Bond, settlementDate time.Time) float64 {
    lastCouponDate := bc.getLastCouponDate(bond, settlementDate)
    daysSinceLastCoupon := settlementDate.Sub(lastCouponDate).Hours() / 24
    daysPerPeriod := 365.0 / float64(bond.Frequency)
    
    couponPayment := bond.FaceValue * bond.CouponRate / float64(bond.Frequency)
    accruedInterest := couponPayment * daysSinceLastCoupon / daysPerPeriod
    
    return bc.calculator.round(accruedInterest)
}

// getLastCouponDate 获取上一个付息日
func (bc *BondCalculator) getLastCouponDate(bond *Bond, settlementDate time.Time) time.Time {
    nextCouponDate := bc.getNextCouponDate(bond, settlementDate)
    daysPerPeriod := 365 / bond.Frequency
    return nextCouponDate.AddDate(0, 0, -daysPerPeriod)
}

// calculateCashFlows 计算现金流
func (bc *BondCalculator) calculateCashFlows(bond *Bond, settlementDate time.Time) []float64 {
    maturityDate := bond.MaturityDate
    nextCouponDate := bc.getNextCouponDate(bond, settlementDate)
    
    var cashFlows []float64
    currentDate := nextCouponDate
    
    for currentDate.Before(maturityDate) || currentDate.Equal(maturityDate) {
        if currentDate.Equal(maturityDate) {
            // 最后一期包含本金
            cashFlows = append(cashFlows, bond.FaceValue*(1+bond.CouponRate/float64(bond.Frequency)))
        } else {
            // 付息
            cashFlows = append(cashFlows, bond.FaceValue*bond.CouponRate/float64(bond.Frequency))
        }
        
        // 移动到下一个付息日
        daysPerPeriod := 365 / bond.Frequency
        currentDate = currentDate.AddDate(0, 0, daysPerPeriod)
    }
    
    return cashFlows
}
```

### 1.2.3 期权定价

```go
// Option 期权结构
type Option struct {
    Type        OptionType `json:"type"`        // CALL or PUT
    StrikePrice float64    `json:"strike_price"`
    SpotPrice   float64    `json:"spot_price"`
    TimeToMaturity float64 `json:"time_to_maturity"`
    RiskFreeRate float64   `json:"risk_free_rate"`
    Volatility   float64   `json:"volatility"`
    DividendYield float64  `json:"dividend_yield"`
}

// OptionType 期权类型
type OptionType string

const (
    OptionTypeCall OptionType = "CALL"
    OptionTypePut  OptionType = "PUT"
)

// OptionCalculator 期权计算器
type OptionCalculator struct {
    calculator *FinancialCalculator
}

// NewOptionCalculator 创建期权计算器
func NewOptionCalculator() *OptionCalculator {
    return &OptionCalculator{
        calculator: NewFinancialCalculator(6),
    }
}

// BlackScholesPrice 使用Black-Scholes模型计算期权价格
func (oc *OptionCalculator) BlackScholesPrice(option *Option) float64 {
    d1 := oc.calculateD1(option)
    d2 := oc.calculateD2(option, d1)
    
    var price float64
    if option.Type == OptionTypeCall {
        price = option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*oc.normalCDF(d1) -
                option.StrikePrice*math.Exp(-option.RiskFreeRate*option.TimeToMaturity)*oc.normalCDF(d2)
    } else {
        price = option.StrikePrice*math.Exp(-option.RiskFreeRate*option.TimeToMaturity)*oc.normalCDF(-d2) -
                option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*oc.normalCDF(-d1)
    }
    
    return oc.calculator.round(price)
}

// calculateD1 计算d1参数
func (oc *OptionCalculator) calculateD1(option *Option) float64 {
    return (math.Log(option.SpotPrice/option.StrikePrice) +
            (option.RiskFreeRate-option.DividendYield+0.5*option.Volatility*option.Volatility)*
                option.TimeToMaturity) / (option.Volatility * math.Sqrt(option.TimeToMaturity))
}

// calculateD2 计算d2参数
func (oc *OptionCalculator) calculateD2(option *Option, d1 float64) float64 {
    return d1 - option.Volatility*math.Sqrt(option.TimeToMaturity)
}

// normalCDF 标准正态分布累积分布函数
func (oc *OptionCalculator) normalCDF(x float64) float64 {
    return 0.5 * (1 + math.Erf(x/math.Sqrt(2)))
}

// Delta 计算Delta
func (oc *OptionCalculator) Delta(option *Option) float64 {
    d1 := oc.calculateD1(option)
    
    var delta float64
    if option.Type == OptionTypeCall {
        delta = math.Exp(-option.DividendYield*option.TimeToMaturity) * oc.normalCDF(d1)
    } else {
        delta = math.Exp(-option.DividendYield*option.TimeToMaturity) * (oc.normalCDF(d1) - 1)
    }
    
    return oc.calculator.round(delta)
}

// Gamma 计算Gamma
func (oc *OptionCalculator) Gamma(option *Option) float64 {
    d1 := oc.calculateD1(option)
    gamma := math.Exp(-option.DividendYield*option.TimeToMaturity) *
             oc.normalPDF(d1) / (option.SpotPrice * option.Volatility * math.Sqrt(option.TimeToMaturity))
    
    return oc.calculator.round(gamma)
}

// Theta 计算Theta
func (oc *OptionCalculator) Theta(option *Option) float64 {
    d1 := oc.calculateD1(option)
    d2 := oc.calculateD2(option, d1)
    
    var theta float64
    if option.Type == OptionTypeCall {
        theta = -option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*
                oc.normalPDF(d1)*option.Volatility/(2*math.Sqrt(option.TimeToMaturity)) +
                option.RiskFreeRate*option.StrikePrice*math.Exp(-option.RiskFreeRate*option.TimeToMaturity)*
                    oc.normalCDF(d2) -
                option.DividendYield*option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*
                    oc.normalCDF(d1)
    } else {
        theta = -option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*
                oc.normalPDF(d1)*option.Volatility/(2*math.Sqrt(option.TimeToMaturity)) -
                option.RiskFreeRate*option.StrikePrice*math.Exp(-option.RiskFreeRate*option.TimeToMaturity)*
                    oc.normalCDF(-d2) +
                option.DividendYield*option.SpotPrice*math.Exp(-option.DividendYield*option.TimeToMaturity)*
                    oc.normalCDF(-d1)
    }
    
    return oc.calculator.round(theta)
}

// Vega 计算Vega
func (oc *OptionCalculator) Vega(option *Option) float64 {
    d1 := oc.calculateD1(option)
    vega := option.SpotPrice * math.Exp(-option.DividendYield*option.TimeToMaturity) *
            oc.normalPDF(d1) * math.Sqrt(option.TimeToMaturity)
    
    return oc.calculator.round(vega)
}

// normalPDF 标准正态分布概率密度函数
func (oc *OptionCalculator) normalPDF(x float64) float64 {
    return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi)
}

// ImpliedVolatility 计算隐含波动率
func (oc *OptionCalculator) ImpliedVolatility(option *Option, marketPrice float64) float64 {
    // 使用牛顿法求解隐含波动率
    const maxIterations = 100
    const tolerance = 1e-6
    
    // 初始猜测值
    volatility := 0.3
    
    for i := 0; i < maxIterations; i++ {
        option.Volatility = volatility
        calculatedPrice := oc.BlackScholesPrice(option)
        
        if math.Abs(calculatedPrice-marketPrice) < tolerance {
            return oc.calculator.round(volatility)
        }
        
        // 计算Vega
        vega := oc.Vega(option)
        if math.Abs(vega) < tolerance {
            break
        }
        
        // 牛顿法迭代
        volatility = volatility - (calculatedPrice-marketPrice)/vega
        
        // 确保波动率在合理范围内
        if volatility < 0.001 {
            volatility = 0.001
        } else if volatility > 5.0 {
            volatility = 5.0
        }
    }
    
    return oc.calculator.round(volatility)
}
```

## 1.3 风险管理算法

### 1.3.1 VaR计算

```go
// RiskManager 风险管理器
type RiskManager struct {
    calculator *FinancialCalculator
}

// NewRiskManager 创建风险管理器
func NewRiskManager() *RiskManager {
    return &RiskManager{
        calculator: NewFinancialCalculator(6),
    }
}

// ValueAtRisk 计算VaR
func (rm *RiskManager) ValueAtRisk(returns []float64, confidenceLevel float64) float64 {
    // 计算收益率统计
    mean := rm.calculateMean(returns)
    stdDev := rm.calculateStdDev(returns, mean)
    
    // 计算分位数
    zScore := rm.calculateZScore(confidenceLevel)
    
    // VaR = 均值 - z分数 * 标准差
    var := mean - zScore*stdDev
    
    return rm.calculator.round(var)
}

// ExpectedShortfall 计算期望损失
func (rm *RiskManager) ExpectedShortfall(returns []float64, confidenceLevel float64) float64 {
    // 排序收益率
    sortedReturns := make([]float64, len(returns))
    copy(sortedReturns, returns)
    sort.Float64s(sortedReturns)
    
    // 计算VaR
    var := rm.ValueAtRisk(returns, confidenceLevel)
    
    // 计算超过VaR的收益率的平均值
    var sum float64
    count := 0
    
    for _, ret := range sortedReturns {
        if ret <= var {
            sum += ret
            count++
        }
    }
    
    if count == 0 {
        return 0
    }
    
    return rm.calculator.round(sum / float64(count))
}

// calculateMean 计算均值
func (rm *RiskManager) calculateMean(data []float64) float64 {
    var sum float64
    for _, value := range data {
        sum += value
    }
    return sum / float64(len(data))
}

// calculateStdDev 计算标准差
func (rm *RiskManager) calculateStdDev(data []float64, mean float64) float64 {
    var sum float64
    for _, value := range data {
        diff := value - mean
        sum += diff * diff
    }
    return math.Sqrt(sum / float64(len(data)-1))
}

// calculateZScore 计算z分数
func (rm *RiskManager) calculateZScore(confidenceLevel float64) float64 {
    // 简化实现，实际应该使用正态分布表
    switch confidenceLevel {
    case 0.95:
        return 1.645
    case 0.99:
        return 2.326
    case 0.995:
        return 2.576
    default:
        return 1.96 // 默认95%置信水平
    }
}
```

## 1.4 总结

金融算法模块涵盖了以下核心内容：

1. **时间价值计算**: 现值、终值、NPV、IRR等基础金融计算
2. **债券定价**: 债券价格、收益率、久期等债券分析
3. **期权定价**: Black-Scholes模型、希腊字母等期权分析
4. **风险管理**: VaR、期望损失等风险度量

这个设计提供了一个完整的金融算法框架，支持各种金融产品的定价和风险管理。
