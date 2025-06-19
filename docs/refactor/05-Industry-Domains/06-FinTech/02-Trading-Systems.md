# 2. 交易系统 (Trading Systems)

## 2.1 概述

### 2.1.1 交易系统定义

交易系统是金融科技领域的核心组件，负责执行金融资产的买卖操作，包括：

- 订单管理
- 价格发现
- 风险控制
- 清算结算

### 2.1.2 系统架构层次

```text
交易系统架构层次：
L1: 用户接口层 (User Interface Layer)
L2: 业务逻辑层 (Business Logic Layer)  
L3: 数据访问层 (Data Access Layer)
L4: 基础设施层 (Infrastructure Layer)
```

## 2.2 核心组件设计

### 2.2.1 订单管理系统

#### 2.2.1.1 订单数据结构

```go
// Order 订单基础结构
type Order struct {
    ID          string    `json:"id"`
    Symbol      string    `json:"symbol"`
    Side        OrderSide `json:"side"`
    Type        OrderType `json:"type"`
    Quantity    decimal.Decimal `json:"quantity"`
    Price       decimal.Decimal `json:"price"`
    Status      OrderStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// OrderSide 订单方向
type OrderSide string

const (
    Buy  OrderSide = "BUY"
    Sell OrderSide = "SELL"
)

// OrderType 订单类型
type OrderType string

const (
    Market OrderType = "MARKET"
    Limit  OrderType = "LIMIT"
    Stop   OrderType = "STOP"
)

// OrderStatus 订单状态
type OrderStatus string

const (
    Pending   OrderStatus = "PENDING"
    Filled    OrderStatus = "FILLED"
    Cancelled OrderStatus = "CANCELLED"
    Rejected  OrderStatus = "REJECTED"
)
```

#### 2.2.1.2 订单管理器

```go
// OrderManager 订单管理器
type OrderManager struct {
    orders    map[string]*Order
    orderBook *OrderBook
    mutex     sync.RWMutex
}

// NewOrderManager 创建订单管理器
func NewOrderManager() *OrderManager {
    return &OrderManager{
        orders:    make(map[string]*Order),
        orderBook: NewOrderBook(),
    }
}

// PlaceOrder 下单
func (om *OrderManager) PlaceOrder(order *Order) error {
    om.mutex.Lock()
    defer om.mutex.Unlock()
    
    // 验证订单
    if err := om.validateOrder(order); err != nil {
        return err
    }
    
    // 生成订单ID
    order.ID = om.generateOrderID()
    order.Status = Pending
    order.CreatedAt = time.Now()
    
    // 添加到订单簿
    om.orderBook.AddOrder(order)
    om.orders[order.ID] = order
    
    return nil
}

// CancelOrder 取消订单
func (om *OrderManager) CancelOrder(orderID string) error {
    om.mutex.Lock()
    defer om.mutex.Unlock()
    
    order, exists := om.orders[orderID]
    if !exists {
        return errors.New("order not found")
    }
    
    if order.Status != Pending {
        return errors.New("order cannot be cancelled")
    }
    
    order.Status = Cancelled
    order.UpdatedAt = time.Now()
    om.orderBook.RemoveOrder(order)
    
    return nil
}
```

### 2.2.2 订单簿系统

#### 2.2.2.1 订单簿数据结构

```go
// OrderBook 订单簿
type OrderBook struct {
    symbol    string
    bids      *OrderQueue // 买单队列
    asks      *OrderQueue // 卖单队列
    lastPrice decimal.Decimal
    mutex     sync.RWMutex
}

// OrderQueue 订单队列
type OrderQueue struct {
    orders []*Order
    mutex  sync.RWMutex
}

// NewOrderBook 创建订单簿
func NewOrderBook() *OrderBook {
    return &OrderBook{
        bids: &OrderQueue{orders: make([]*Order, 0)},
        asks: &OrderQueue{orders: make([]*Order, 0)},
    }
}
```

#### 2.2.2.2 价格匹配算法

```go
// MatchOrders 订单匹配
func (ob *OrderBook) MatchOrders(newOrder *Order) []*Trade {
    ob.mutex.Lock()
    defer ob.mutex.Unlock()
    
    var trades []*Trade
    
    switch newOrder.Side {
    case Buy:
        trades = ob.matchBuyOrder(newOrder)
    case Sell:
        trades = ob.matchSellOrder(newOrder)
    }
    
    return trades
}

// matchBuyOrder 匹配买单
func (ob *OrderBook) matchBuyOrder(buyOrder *Order) []*Trade {
    var trades []*Trade
    remainingQty := buyOrder.Quantity
    
    for remainingQty.GreaterThan(decimal.Zero) && len(ob.asks.orders) > 0 {
        askOrder := ob.asks.orders[0]
        
        // 检查价格匹配
        if buyOrder.Type == Limit && buyOrder.Price.LessThan(askOrder.Price) {
            break
        }
        
        // 计算成交数量
        tradeQty := decimal.Min(remainingQty, askOrder.Quantity)
        tradePrice := askOrder.Price
        
        // 创建成交记录
        trade := &Trade{
            ID:        ob.generateTradeID(),
            BuyOrder:  buyOrder.ID,
            SellOrder: askOrder.ID,
            Symbol:    buyOrder.Symbol,
            Quantity:  tradeQty,
            Price:     tradePrice,
            Timestamp: time.Now(),
        }
        trades = append(trades, trade)
        
        // 更新订单数量
        remainingQty = remainingQty.Sub(tradeQty)
        askOrder.Quantity = askOrder.Quantity.Sub(tradeQty)
        
        // 如果卖单完全成交，移除订单
        if askOrder.Quantity.Equal(decimal.Zero) {
            ob.asks.RemoveOrder(0)
            askOrder.Status = Filled
        }
    }
    
    // 更新买单剩余数量
    buyOrder.Quantity = remainingQty
    if remainingQty.Equal(decimal.Zero) {
        buyOrder.Status = Filled
    }
    
    return trades
}
```

### 2.2.3 风险管理系统

#### 2.2.3.1 风险控制规则

```go
// RiskManager 风险管理器
type RiskManager struct {
    rules []RiskRule
    limits map[string]decimal.Decimal
}

// RiskRule 风险规则接口
type RiskRule interface {
    Check(order *Order, account *Account) error
}

// PositionLimitRule 持仓限制规则
type PositionLimitRule struct {
    maxPosition decimal.Decimal
}

func (r *PositionLimitRule) Check(order *Order, account *Account) error {
    currentPosition := account.GetPosition(order.Symbol)
    newPosition := currentPosition.Add(order.Quantity)
    
    if newPosition.Abs().GreaterThan(r.maxPosition) {
        return errors.New("position limit exceeded")
    }
    
    return nil
}

// OrderSizeRule 订单大小规则
type OrderSizeRule struct {
    maxOrderSize decimal.Decimal
}

func (r *OrderSizeRule) Check(order *Order, account *Account) error {
    if order.Quantity.GreaterThan(r.maxOrderSize) {
        return errors.New("order size too large")
    }
    
    return nil
}
```

#### 2.2.3.2 实时风险监控

```go
// RiskMonitor 风险监控器
type RiskMonitor struct {
    riskManager *RiskManager
    alerts      chan RiskAlert
}

// RiskAlert 风险警报
type RiskAlert struct {
    Level     AlertLevel
    Message   string
    OrderID   string
    Timestamp time.Time
}

// MonitorOrder 监控订单风险
func (rm *RiskMonitor) MonitorOrder(order *Order, account *Account) error {
    for _, rule := range rm.riskManager.rules {
        if err := rule.Check(order, account); err != nil {
            alert := RiskAlert{
                Level:     High,
                Message:   err.Error(),
                OrderID:   order.ID,
                Timestamp: time.Now(),
            }
            rm.alerts <- alert
            return err
        }
    }
    
    return nil
}
```

## 2.3 性能优化策略

### 2.3.1 内存池优化

```go
// OrderPool 订单对象池
var orderPool = sync.Pool{
    New: func() interface{} {
        return &Order{}
    },
}

// GetOrder 从池中获取订单对象
func GetOrder() *Order {
    return orderPool.Get().(*Order)
}

// PutOrder 归还订单对象到池
func PutOrder(order *Order) {
    order.Reset()
    orderPool.Put(order)
}
```

### 2.3.2 并发处理优化

```go
// TradingEngine 交易引擎
type TradingEngine struct {
    orderManager *OrderManager
    riskMonitor  *RiskMonitor
    workers      []*Worker
    orderChan    chan *Order
    resultChan   chan *OrderResult
}

// Worker 工作协程
type Worker struct {
    id       int
    engine   *TradingEngine
    orderChan chan *Order
}

// ProcessOrder 处理订单
func (w *Worker) ProcessOrder(order *Order) {
    // 风险检查
    if err := w.engine.riskMonitor.MonitorOrder(order, w.getAccount(order.UserID)); err != nil {
        w.sendResult(order, err)
        return
    }
    
    // 下单处理
    if err := w.engine.orderManager.PlaceOrder(order); err != nil {
        w.sendResult(order, err)
        return
    }
    
    w.sendResult(order, nil)
}
```

## 2.4 数据一致性保证

### 2.4.1 事务管理

```go
// TradingTransaction 交易事务
type TradingTransaction struct {
    db *sql.DB
    tx *sql.Tx
}

// Begin 开始事务
func (tt *TradingTransaction) Begin() error {
    tx, err := tt.db.Begin()
    if err != nil {
        return err
    }
    tt.tx = tx
    return nil
}

// Commit 提交事务
func (tt *TradingTransaction) Commit() error {
    return tt.tx.Commit()
}

// Rollback 回滚事务
func (tt *TradingTransaction) Rollback() error {
    return tt.tx.Rollback()
}

// SaveOrder 保存订单
func (tt *TradingTransaction) SaveOrder(order *Order) error {
    query := `
        INSERT INTO orders (id, symbol, side, type, quantity, price, status, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := tt.tx.Exec(query, order.ID, order.Symbol, order.Side, order.Type,
        order.Quantity, order.Price, order.Status, order.CreatedAt)
    return err
}
```

### 2.4.2 事件溯源

```go
// TradingEvent 交易事件
type TradingEvent struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    Version   int64                  `json:"version"`
}

// EventStore 事件存储
type EventStore struct {
    events []*TradingEvent
    mutex  sync.RWMutex
}

// AppendEvent 追加事件
func (es *EventStore) AppendEvent(event *TradingEvent) error {
    es.mutex.Lock()
    defer es.mutex.Unlock()
    
    event.ID = es.generateEventID()
    event.Timestamp = time.Now()
    event.Version = int64(len(es.events)) + 1
    
    es.events = append(es.events, event)
    return nil
}

// GetEvents 获取事件
func (es *EventStore) GetEvents(aggregateID string) []*TradingEvent {
    es.mutex.RLock()
    defer es.mutex.RUnlock()
    
    var result []*TradingEvent
    for _, event := range es.events {
        if event.Data["aggregate_id"] == aggregateID {
            result = append(result, event)
        }
    }
    
    return result
}
```

## 2.5 监控和日志

### 2.5.1 性能监控

```go
// TradingMetrics 交易指标
type TradingMetrics struct {
    ordersPerSecond    int64
    latency            time.Duration
    errorRate          float64
    activeOrders       int64
    totalVolume        decimal.Decimal
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
    metrics *TradingMetrics
    mutex   sync.RWMutex
}

// RecordOrder 记录订单指标
func (mc *MetricsCollector) RecordOrder(latency time.Duration, success bool) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    mc.metrics.ordersPerSecond++
    mc.metrics.latency = latency
    
    if !success {
        mc.metrics.errorRate++
    }
}

// GetMetrics 获取指标
func (mc *MetricsCollector) GetMetrics() *TradingMetrics {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    return &TradingMetrics{
        ordersPerSecond: mc.metrics.ordersPerSecond,
        latency:         mc.metrics.latency,
        errorRate:       mc.metrics.errorRate,
        activeOrders:    mc.metrics.activeOrders,
        totalVolume:     mc.metrics.totalVolume,
    }
}
```

### 2.5.2 结构化日志

```go
// TradingLogger 交易日志器
type TradingLogger struct {
    logger *log.Logger
}

// LogOrder 记录订单日志
func (tl *TradingLogger) LogOrder(order *Order, action string) {
    tl.logger.Printf("[ORDER] %s | ID: %s | Symbol: %s | Side: %s | Quantity: %s | Price: %s",
        action, order.ID, order.Symbol, order.Side, order.Quantity, order.Price)
}

// LogTrade 记录成交日志
func (tl *TradingLogger) LogTrade(trade *Trade) {
    tl.logger.Printf("[TRADE] ID: %s | Symbol: %s | Quantity: %s | Price: %s | BuyOrder: %s | SellOrder: %s",
        trade.ID, trade.Symbol, trade.Quantity, trade.Price, trade.BuyOrder, trade.SellOrder)
}

// LogError 记录错误日志
func (tl *TradingLogger) LogError(err error, context string) {
    tl.logger.Printf("[ERROR] %s | %v", context, err)
}
```

## 2.6 测试策略

### 2.6.1 单元测试

```go
// TestOrderManager 测试订单管理器
func TestOrderManager(t *testing.T) {
    om := NewOrderManager()
    
    // 测试下单
    order := &Order{
        Symbol:   "BTC/USDT",
        Side:     Buy,
        Type:     Limit,
        Quantity: decimal.NewFromFloat(1.0),
        Price:    decimal.NewFromFloat(50000.0),
    }
    
    err := om.PlaceOrder(order)
    assert.NoError(t, err)
    assert.Equal(t, Pending, order.Status)
    
    // 测试取消订单
    err = om.CancelOrder(order.ID)
    assert.NoError(t, err)
    assert.Equal(t, Cancelled, order.Status)
}
```

### 2.6.2 性能测试

```go
// BenchmarkOrderMatching 订单匹配性能测试
func BenchmarkOrderMatching(b *testing.B) {
    ob := NewOrderBook()
    
    // 准备测试数据
    orders := generateTestOrders(1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, order := range orders {
            ob.MatchOrders(order)
        }
    }
}

// BenchmarkConcurrentOrders 并发订单处理测试
func BenchmarkConcurrentOrders(b *testing.B) {
    engine := NewTradingEngine()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            order := generateRandomOrder()
            engine.ProcessOrder(order)
        }
    })
}
```

## 2.7 部署和运维

### 2.7.1 容器化部署

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o trading-system ./cmd/trading-system

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/trading-system .
CMD ["./trading-system"]
```

### 2.7.2 配置管理

```yaml
# config.yaml
trading:
  order_book:
    max_orders: 10000
    price_precision: 8
    quantity_precision: 8
  
  risk_management:
    max_position: 1000000
    max_order_size: 10000
    max_daily_loss: 100000
  
  performance:
    workers: 8
    queue_size: 10000
    timeout: 5s
  
  database:
    host: localhost
    port: 5432
    name: trading_system
    user: trading_user
    password: trading_pass
  
  redis:
    host: localhost
    port: 6379
    db: 0
```

## 2.8 总结

交易系统是金融科技的核心组件，需要：

1. **高性能设计** - 使用内存池、并发处理、缓存优化
2. **风险控制** - 实时监控、规则引擎、限额管理
3. **数据一致性** - 事务管理、事件溯源、状态同步
4. **可观测性** - 指标监控、结构化日志、链路追踪
5. **可扩展性** - 微服务架构、水平扩展、负载均衡

通过Go语言的高并发特性和丰富的生态系统，可以构建出高性能、高可靠性的交易系统。
