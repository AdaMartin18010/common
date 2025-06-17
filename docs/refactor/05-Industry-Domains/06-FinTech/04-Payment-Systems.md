# 04-支付系统 (Payment Systems)

## 目录

1. [概述](#1-概述)
2. [支付系统架构](#2-支付系统架构)
3. [支付处理](#3-支付处理)
4. [支付网关](#4-支付网关)
5. [清算结算](#5-清算结算)
6. [安全机制](#6-安全机制)
7. [监控告警](#7-监控告警)
8. [总结](#8-总结)

## 1. 概述

### 1.1 支付系统的重要性

支付系统是现代金融基础设施的核心，负责处理各种支付交易。Go语言的高并发和低延迟特性使其特别适合构建高性能的支付系统。

### 1.2 支付系统组件

```go
type PaymentSystem struct {
    Gateway      *PaymentGateway
    Processor    *PaymentProcessor
    Settlement   *SettlementEngine
    Security     *SecurityManager
    Monitor      *PaymentMonitor
}

type PaymentConfig struct {
    MaxAmount     float64
    MinAmount     float64
    Timeout       time.Duration
    RetryCount    int
    Currency      string
}
```

## 2. 支付系统架构

### 2.1 微服务架构

```go
type PaymentMicroservices struct {
    Services map[string]*PaymentService
}

type PaymentService struct {
    Name        string
    Port        int
    Handler     http.Handler
    Dependencies []string
}

func NewPaymentMicroservices() *PaymentMicroservices {
    pms := &PaymentMicroservices{
        Services: make(map[string]*PaymentService),
    }
    
    // 支付网关服务
    pms.Services["gateway"] = &PaymentService{
        Name:         "gateway",
        Port:         8081,
        Dependencies: []string{"processor", "security"},
        Handler:      createGatewayHandler(),
    }
    
    // 支付处理服务
    pms.Services["processor"] = &PaymentService{
        Name:         "processor",
        Port:         8082,
        Dependencies: []string{"settlement"},
        Handler:      createProcessorHandler(),
    }
    
    // 清算结算服务
    pms.Services["settlement"] = &PaymentService{
        Name:         "settlement",
        Port:         8083,
        Dependencies: []string{},
        Handler:      createSettlementHandler(),
    }
    
    return pms
}
```

### 2.2 事件驱动架构

```go
type PaymentEventBus struct {
    subscribers map[string][]chan PaymentEvent
    mu          sync.RWMutex
}

type PaymentEvent struct {
    Type      string
    Data      interface{}
    Timestamp time.Time
    Source    string
}

const (
    EventPaymentInitiated = "payment.initiated"
    EventPaymentProcessed = "payment.processed"
    EventPaymentSettled   = "payment.settled"
    EventPaymentFailed    = "payment.failed"
    EventPaymentReversed  = "payment.reversed"
)

func NewPaymentEventBus() *PaymentEventBus {
    return &PaymentEventBus{
        subscribers: make(map[string][]chan PaymentEvent),
    }
}

func (peb *PaymentEventBus) Publish(event PaymentEvent) {
    peb.mu.RLock()
    defer peb.mu.RUnlock()
    
    for _, ch := range peb.subscribers[event.Type] {
        select {
        case ch <- event:
        default:
            // 通道已满，跳过
        }
    }
}
```

## 3. 支付处理

### 3.1 支付交易

```go
type PaymentTransaction struct {
    ID              string
    MerchantID      string
    CustomerID      string
    Amount          float64
    Currency        string
    PaymentMethod   PaymentMethod
    Status          PaymentStatus
    CreatedAt       time.Time
    UpdatedAt       time.Time
    ProcessedAt     *time.Time
    SettledAt       *time.Time
    Reference       string
    Description     string
    Metadata        map[string]interface{}
}

type PaymentMethod string
const (
    PaymentMethodCard     PaymentMethod = "CARD"
    PaymentMethodBank     PaymentMethod = "BANK"
    PaymentMethodDigital  PaymentMethod = "DIGITAL"
    PaymentMethodCrypto   PaymentMethod = "CRYPTO"
)

type PaymentStatus string
const (
    PaymentStatusPending   PaymentStatus = "PENDING"
    PaymentStatusProcessing PaymentStatus = "PROCESSING"
    PaymentStatusCompleted PaymentStatus = "COMPLETED"
    PaymentStatusFailed    PaymentStatus = "FAILED"
    PaymentStatusReversed  PaymentStatus = "REVERSED"
)

// 支付处理器
type PaymentProcessor struct {
    transactions map[string]*PaymentTransaction
    eventBus     *PaymentEventBus
    security     *SecurityManager
    mu           sync.RWMutex
}

func NewPaymentProcessor(eventBus *PaymentEventBus, security *SecurityManager) *PaymentProcessor {
    return &PaymentProcessor{
        transactions: make(map[string]*PaymentTransaction),
        eventBus:     eventBus,
        security:     security,
    }
}

func (pp *PaymentProcessor) ProcessPayment(payment *PaymentTransaction) error {
    // 验证支付
    if err := pp.validatePayment(payment); err != nil {
        return err
    }
    
    // 安全检查
    if err := pp.security.CheckPayment(payment); err != nil {
        payment.Status = PaymentStatusFailed
        return err
    }
    
    // 保存交易
    pp.mu.Lock()
    pp.transactions[payment.ID] = payment
    pp.mu.Unlock()
    
    // 发布事件
    pp.eventBus.Publish(PaymentEvent{
        Type:      EventPaymentInitiated,
        Data:      payment,
        Timestamp: time.Now(),
        Source:    "payment-processor",
    })
    
    // 异步处理支付
    go pp.processPaymentAsync(payment)
    
    return nil
}

func (pp *PaymentProcessor) validatePayment(payment *PaymentTransaction) error {
    if payment.Amount <= 0 {
        return fmt.Errorf("invalid amount: %f", payment.Amount)
    }
    
    if payment.MerchantID == "" {
        return fmt.Errorf("merchant ID is required")
    }
    
    if payment.CustomerID == "" {
        return fmt.Errorf("customer ID is required")
    }
    
    return nil
}

func (pp *PaymentProcessor) processPaymentAsync(payment *PaymentTransaction) {
    // 更新状态为处理中
    payment.Status = PaymentStatusProcessing
    payment.UpdatedAt = time.Now()
    
    // 模拟支付处理
    time.Sleep(100 * time.Millisecond)
    
    // 随机成功或失败
    if rand.Float64() > 0.1 { // 90%成功率
        payment.Status = PaymentStatusCompleted
        now := time.Now()
        payment.ProcessedAt = &now
        payment.UpdatedAt = now
        
        pp.eventBus.Publish(PaymentEvent{
            Type:      EventPaymentProcessed,
            Data:      payment,
            Timestamp: now,
            Source:    "payment-processor",
        })
    } else {
        payment.Status = PaymentStatusFailed
        payment.UpdatedAt = time.Now()
        
        pp.eventBus.Publish(PaymentEvent{
            Type:      EventPaymentFailed,
            Data:      payment,
            Timestamp: time.Now(),
            Source:    "payment-processor",
        })
    }
}
```

### 3.2 支付路由

```go
type PaymentRouter struct {
    routes map[PaymentMethod]*PaymentRoute
    mu     sync.RWMutex
}

type PaymentRoute struct {
    Method    PaymentMethod
    Processor string
    Priority  int
    Enabled   bool
}

func NewPaymentRouter() *PaymentRouter {
    pr := &PaymentRouter{
        routes: make(map[PaymentMethod]*PaymentRoute),
    }
    
    // 配置支付路由
    pr.routes[PaymentMethodCard] = &PaymentRoute{
        Method:    PaymentMethodCard,
        Processor: "card-processor",
        Priority:  1,
        Enabled:   true,
    }
    
    pr.routes[PaymentMethodBank] = &PaymentRoute{
        Method:    PaymentMethodBank,
        Processor: "bank-processor",
        Priority:  2,
        Enabled:   true,
    }
    
    return pr
}

func (pr *PaymentRouter) RoutePayment(payment *PaymentTransaction) (*PaymentRoute, error) {
    pr.mu.RLock()
    defer pr.mu.RUnlock()
    
    route, exists := pr.routes[payment.PaymentMethod]
    if !exists {
        return nil, fmt.Errorf("no route for payment method: %s", payment.PaymentMethod)
    }
    
    if !route.Enabled {
        return nil, fmt.Errorf("payment route is disabled: %s", payment.PaymentMethod)
    }
    
    return route, nil
}
```

## 4. 支付网关

### 4.1 网关接口

```go
type PaymentGateway struct {
    router    *PaymentRouter
    processor *PaymentProcessor
    security  *SecurityManager
    monitor   *PaymentMonitor
}

func NewPaymentGateway(router *PaymentRouter, processor *PaymentProcessor, 
                      security *SecurityManager, monitor *PaymentMonitor) *PaymentGateway {
    return &PaymentGateway{
        router:    router,
        processor: processor,
        security:  security,
        monitor:   monitor,
    }
}

func (pg *PaymentGateway) CreatePayment(request *PaymentRequest) (*PaymentResponse, error) {
    // 记录开始时间
    start := time.Now()
    
    // 创建支付交易
    payment := &PaymentTransaction{
        ID:            generatePaymentID(),
        MerchantID:    request.MerchantID,
        CustomerID:    request.CustomerID,
        Amount:        request.Amount,
        Currency:      request.Currency,
        PaymentMethod: request.PaymentMethod,
        Status:        PaymentStatusPending,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
        Reference:     request.Reference,
        Description:   request.Description,
        Metadata:      request.Metadata,
    }
    
    // 路由支付
    route, err := pg.router.RoutePayment(payment)
    if err != nil {
        return nil, err
    }
    
    // 处理支付
    if err := pg.processor.ProcessPayment(payment); err != nil {
        return nil, err
    }
    
    // 记录监控指标
    pg.monitor.RecordPayment(payment, time.Since(start))
    
    return &PaymentResponse{
        PaymentID: payment.ID,
        Status:    payment.Status,
        Message:   "Payment initiated successfully",
    }, nil
}

type PaymentRequest struct {
    MerchantID    string                 `json:"merchant_id"`
    CustomerID    string                 `json:"customer_id"`
    Amount        float64                `json:"amount"`
    Currency      string                 `json:"currency"`
    PaymentMethod PaymentMethod          `json:"payment_method"`
    Reference     string                 `json:"reference"`
    Description   string                 `json:"description"`
    Metadata      map[string]interface{} `json:"metadata"`
}

type PaymentResponse struct {
    PaymentID string        `json:"payment_id"`
    Status    PaymentStatus `json:"status"`
    Message   string        `json:"message"`
}
```

### 4.2 网关API

```go
func createGatewayHandler() http.Handler {
    mux := http.NewServeMux()
    
    // 创建支付
    mux.HandleFunc("POST /payments", func(w http.ResponseWriter, r *http.Request) {
        var request PaymentRequest
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        response, err := gateway.CreatePayment(&request)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })
    
    // 查询支付状态
    mux.HandleFunc("GET /payments/{id}", func(w http.ResponseWriter, r *http.Request) {
        paymentID := chi.URLParam(r, "id")
        
        payment, err := processor.GetPayment(paymentID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(payment)
    })
    
    // 退款
    mux.HandleFunc("POST /payments/{id}/refund", func(w http.ResponseWriter, r *http.Request) {
        paymentID := chi.URLParam(r, "id")
        
        var request RefundRequest
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        response, err := gateway.RefundPayment(paymentID, &request)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })
    
    return mux
}

type RefundRequest struct {
    Amount   float64 `json:"amount"`
    Reason   string  `json:"reason"`
    Reference string `json:"reference"`
}

type RefundResponse struct {
    RefundID string `json:"refund_id"`
    Status   string `json:"status"`
    Message  string `json:"message"`
}
```

## 5. 清算结算

### 5.1 清算引擎

```go
type SettlementEngine struct {
    settlements map[string]*Settlement
    accounts    map[string]*Account
    eventBus    *PaymentEventBus
    mu          sync.RWMutex
}

type Settlement struct {
    ID            string
    PaymentID     string
    MerchantID    string
    Amount        float64
    Currency      string
    Status        SettlementStatus
    CreatedAt     time.Time
    SettledAt     *time.Time
    Fee           float64
    NetAmount     float64
}

type SettlementStatus string
const (
    SettlementStatusPending   SettlementStatus = "PENDING"
    SettlementStatusProcessing SettlementStatus = "PROCESSING"
    SettlementStatusCompleted SettlementStatus = "COMPLETED"
    SettlementStatusFailed    SettlementStatus = "FAILED"
)

type Account struct {
    ID       string
    Balance  float64
    Currency string
    UpdatedAt time.Time
}

func NewSettlementEngine(eventBus *PaymentEventBus) *SettlementEngine {
    se := &SettlementEngine{
        settlements: make(map[string]*Settlement),
        accounts:    make(map[string]*Account),
        eventBus:    eventBus,
    }
    
    // 监听支付完成事件
    go se.listenPaymentEvents()
    
    return se
}

func (se *SettlementEngine) listenPaymentEvents() {
    ch := make(chan PaymentEvent, 100)
    se.eventBus.Subscribe(EventPaymentProcessed, ch)
    
    for event := range ch {
        if payment, ok := event.Data.(*PaymentTransaction); ok {
            se.createSettlement(payment)
        }
    }
}

func (se *SettlementEngine) createSettlement(payment *PaymentTransaction) {
    settlement := &Settlement{
        ID:         generateSettlementID(),
        PaymentID:  payment.ID,
        MerchantID: payment.MerchantID,
        Amount:     payment.Amount,
        Currency:   payment.Currency,
        Status:     SettlementStatusPending,
        CreatedAt:  time.Now(),
        Fee:        se.calculateFee(payment.Amount),
    }
    
    settlement.NetAmount = settlement.Amount - settlement.Fee
    
    se.mu.Lock()
    se.settlements[settlement.ID] = settlement
    se.mu.Unlock()
    
    // 异步处理清算
    go se.processSettlement(settlement)
}

func (se *SettlementEngine) calculateFee(amount float64) float64 {
    // 简化的费用计算：2.9% + $0.30
    return amount*0.029 + 0.30
}

func (se *SettlementEngine) processSettlement(settlement *Settlement) {
    // 更新状态为处理中
    settlement.Status = SettlementStatusProcessing
    
    // 模拟清算处理
    time.Sleep(200 * time.Millisecond)
    
    // 更新商户账户
    se.updateMerchantAccount(settlement)
    
    // 标记为完成
    settlement.Status = SettlementStatusCompleted
    now := time.Now()
    settlement.SettledAt = &now
    
    // 发布清算完成事件
    se.eventBus.Publish(PaymentEvent{
        Type:      EventPaymentSettled,
        Data:      settlement,
        Timestamp: now,
        Source:    "settlement-engine",
    })
}

func (se *SettlementEngine) updateMerchantAccount(settlement *Settlement) {
    se.mu.Lock()
    defer se.mu.Unlock()
    
    account, exists := se.accounts[settlement.MerchantID]
    if !exists {
        account = &Account{
            ID:       settlement.MerchantID,
            Balance:  0,
            Currency: settlement.Currency,
        }
        se.accounts[settlement.MerchantID] = account
    }
    
    account.Balance += settlement.NetAmount
    account.UpdatedAt = time.Now()
}
```

## 6. 安全机制

### 6.1 安全管理器

```go
type SecurityManager struct {
    validators map[string]Validator
    encryptor  *Encryptor
    rateLimiter *RateLimiter
}

type Validator interface {
    Validate(data interface{}) error
}

type Encryptor struct {
    key []byte
}

type RateLimiter struct {
    limits map[string]*RateLimit
    mu     sync.RWMutex
}

type RateLimit struct {
    Key       string
    Limit     int
    Window    time.Duration
    Count     int
    ResetTime time.Time
}

func NewSecurityManager() *SecurityManager {
    sm := &SecurityManager{
        validators: make(map[string]Validator),
        encryptor:  &Encryptor{key: []byte("secret-key")},
        rateLimiter: &RateLimiter{limits: make(map[string]*RateLimit)},
    }
    
    // 注册验证器
    sm.validators["amount"] = &AmountValidator{}
    sm.validators["merchant"] = &MerchantValidator{}
    sm.validators["card"] = &CardValidator{}
    
    return sm
}

func (sm *SecurityManager) CheckPayment(payment *PaymentTransaction) error {
    // 验证金额
    if err := sm.validators["amount"].Validate(payment.Amount); err != nil {
        return err
    }
    
    // 验证商户
    if err := sm.validators["merchant"].Validate(payment.MerchantID); err != nil {
        return err
    }
    
    // 检查频率限制
    if err := sm.rateLimiter.CheckLimit(payment.MerchantID); err != nil {
        return err
    }
    
    return nil
}

// 金额验证器
type AmountValidator struct{}

func (av *AmountValidator) Validate(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    
    if amount > 10000 {
        return fmt.Errorf("amount exceeds limit: %f", amount)
    }
    
    return nil
}

// 商户验证器
type MerchantValidator struct{}

func (mv *MerchantValidator) Validate(merchantID string) error {
    if merchantID == "" {
        return fmt.Errorf("merchant ID is required")
    }
    
    // 检查商户是否有效
    if !isValidMerchant(merchantID) {
        return fmt.Errorf("invalid merchant: %s", merchantID)
    }
    
    return nil
}

// 频率限制器
func (rl *RateLimiter) CheckLimit(key string) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limit, exists := rl.limits[key]
    if !exists {
        limit = &RateLimit{
            Key:       key,
            Limit:     100,
            Window:    time.Minute,
            ResetTime: time.Now().Add(time.Minute),
        }
        rl.limits[key] = limit
    }
    
    // 检查是否需要重置
    if time.Now().After(limit.ResetTime) {
        limit.Count = 0
        limit.ResetTime = time.Now().Add(limit.Window)
    }
    
    // 检查限制
    if limit.Count >= limit.Limit {
        return fmt.Errorf("rate limit exceeded for key: %s", key)
    }
    
    limit.Count++
    return nil
}
```

## 7. 监控告警

### 7.1 支付监控

```go
type PaymentMonitor struct {
    metrics map[string]*PaymentMetric
    alerts  chan PaymentAlert
    mu      sync.RWMutex
}

type PaymentMetric struct {
    Name      string
    Value     float64
    Count     int64
    UpdatedAt time.Time
}

type PaymentAlert struct {
    ID        string
    Type      string
    Severity  string
    Message   string
    Timestamp time.Time
    Data      map[string]interface{}
}

func NewPaymentMonitor() *PaymentMonitor {
    pm := &PaymentMonitor{
        metrics: make(map[string]*PaymentMetric),
        alerts:  make(chan PaymentAlert, 100),
    }
    
    // 启动监控协程
    go pm.runMonitoring()
    
    return pm
}

func (pm *PaymentMonitor) RecordPayment(payment *PaymentTransaction, duration time.Duration) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    // 记录交易量
    pm.updateMetric("transaction_volume", payment.Amount)
    
    // 记录交易数量
    pm.updateMetric("transaction_count", 1)
    
    // 记录处理时间
    pm.updateMetric("processing_time", duration.Seconds())
    
    // 记录成功率
    if payment.Status == PaymentStatusCompleted {
        pm.updateMetric("success_rate", 1)
    } else {
        pm.updateMetric("success_rate", 0)
    }
}

func (pm *PaymentMonitor) updateMetric(name string, value float64) {
    metric, exists := pm.metrics[name]
    if !exists {
        metric = &PaymentMetric{Name: name}
        pm.metrics[name] = metric
    }
    
    metric.Value = value
    metric.Count++
    metric.UpdatedAt = time.Now()
}

func (pm *PaymentMonitor) runMonitoring() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        pm.checkMetrics()
    }
}

func (pm *PaymentMonitor) checkMetrics() {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    // 检查成功率
    if successRate, exists := pm.metrics["success_rate"]; exists {
        if successRate.Value < 0.95 { // 低于95%成功率
            pm.createAlert("low_success_rate", "HIGH", 
                fmt.Sprintf("Success rate is low: %.2f%%", successRate.Value*100))
        }
    }
    
    // 检查处理时间
    if processingTime, exists := pm.metrics["processing_time"]; exists {
        if processingTime.Value > 1.0 { // 超过1秒
            pm.createAlert("high_processing_time", "MEDIUM",
                fmt.Sprintf("Processing time is high: %.2fs", processingTime.Value))
        }
    }
}

func (pm *PaymentMonitor) createAlert(alertType, severity, message string) {
    alert := PaymentAlert{
        ID:        generateAlertID(),
        Type:      alertType,
        Severity:  severity,
        Message:   message,
        Timestamp: time.Now(),
        Data:      make(map[string]interface{}),
    }
    
    select {
    case pm.alerts <- alert:
    default:
        log.Printf("Alert channel full, dropping alert: %s", message)
    }
}

func (pm *PaymentMonitor) GetAlerts() <-chan PaymentAlert {
    return pm.alerts
}
```

## 8. 总结

### 8.1 支付系统优势

Go语言在构建支付系统方面的优势：

1. **高性能**: 低延迟的支付处理
2. **高并发**: 同时处理大量支付请求
3. **可靠性**: 强类型系统和错误处理
4. **可扩展性**: 微服务架构支持
5. **安全性**: 内置的安全特性

### 8.2 最佳实践

1. **异步处理**: 使用goroutine处理支付
2. **事件驱动**: 使用事件总线解耦组件
3. **监控告警**: 实时监控支付状态
4. **安全验证**: 多层安全验证机制
5. **错误处理**: 完善的错误处理和恢复

### 8.3 技术挑战

1. **一致性**: 确保支付数据一致性
2. **安全性**: 防止欺诈和攻击
3. **合规性**: 满足监管要求
4. **可扩展性**: 支持大规模交易
5. **可用性**: 高可用性要求

Go语言凭借其优秀的性能和并发特性，是构建现代支付系统的理想选择。

---

**相关链接**:
- [01-Financial-Algorithms](../01-Financial-Algorithms.md)
- [02-Trading-Systems](../02-Trading-Systems.md)
- [03-Risk-Management](../03-Risk-Management.md)
- [../README.md](../README.md) 