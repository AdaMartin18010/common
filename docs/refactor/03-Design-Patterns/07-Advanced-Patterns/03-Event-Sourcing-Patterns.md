# 03-事件溯源模式 (Event Sourcing Patterns)

## 概述

事件溯源模式是一种数据持久化模式，将系统的状态变化记录为一系列事件，而不是直接存储当前状态。本文档探讨在Go语言中实现事件溯源模式的方法和理论。

## 目录

1. [理论基础](#理论基础)
2. [核心概念](#核心概念)
3. [事件模型](#事件模型)
4. [事件存储](#事件存储)
5. [事件流处理](#事件流处理)
6. [快照机制](#快照机制)
7. [投影模式](#投影模式)
8. [并发控制](#并发控制)
9. [实际应用](#实际应用)
10. [最佳实践](#最佳实践)

## 理论基础

### 事件溯源的数学基础

事件溯源基于事件流的概念，可以形式化为：

```latex
S_t = S_0 \oplus e_1 \oplus e_2 \oplus \cdots \oplus e_t
```

其中：

- ```latex
$S_t$
``` 是时间 ```latex
$t$
``` 的状态
- ```latex
$S_0$
``` 是初始状态
- ```latex
$e_i$
``` 是第 ```latex
$i$
``` 个事件
- ```latex
$\oplus$
``` 是状态转换操作

### 事件序列的不可变性

事件序列一旦创建就不可修改，这保证了：

```latex
\forall i < j: e_i \prec e_j \Rightarrow \text{事件顺序不可变}
```

### 状态重建

从事件序列重建状态：

```latex
\text{rebuild}(S_0, [e_1, e_2, \ldots, e_n]) = S_0 \oplus e_1 \oplus e_2 \oplus \cdots \oplus e_n
```

## 核心概念

### 1. 事件 (Event)

事件是不可变的、描述状态变化的事实：

```go
// 事件接口
type Event interface {
    EventID() string
    EventType() string
    Timestamp() time.Time
    AggregateID() string
    Version() int
}

// 基础事件实现
type BaseEvent struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    Timestamp   time.Time `json:"timestamp"`
    AggregateID string    `json:"aggregate_id"`
    Version     int       `json:"version"`
}

func (e BaseEvent) EventID() string {
    return e.ID
}

func (e BaseEvent) EventType() string {
    return e.Type
}

func (e BaseEvent) Timestamp() time.Time {
    return e.Timestamp
}

func (e BaseEvent) AggregateID() string {
    return e.AggregateID
}

func (e BaseEvent) Version() int {
    return e.Version
}
```

### 2. 聚合根 (Aggregate Root)

聚合根是业务逻辑的入口点，负责产生事件：

```go
// 聚合根接口
type AggregateRoot interface {
    ID() string
    Version() int
    UncommittedEvents() []Event
    MarkEventsAsCommitted()
    Apply(event Event) error
}

// 基础聚合根实现
type BaseAggregate struct {
    id                string
    version           int
    uncommittedEvents []Event
}

func (a *BaseAggregate) ID() string {
    return a.id
}

func (a *BaseAggregate) Version() int {
    return a.version
}

func (a *BaseAggregate) UncommittedEvents() []Event {
    return a.uncommittedEvents
}

func (a *BaseAggregate) MarkEventsAsCommitted() {
    a.uncommittedEvents = nil
}

func (a *BaseAggregate) AddEvent(event Event) {
    a.uncommittedEvents = append(a.uncommittedEvents, event)
    a.version++
}
```

### 3. 事件存储 (Event Store)

事件存储负责持久化事件：

```go
// 事件存储接口
type EventStore interface {
    SaveEvents(aggregateID string, events []Event, expectedVersion int) error
    GetEvents(aggregateID string) ([]Event, error)
    GetEventsFromVersion(aggregateID string, fromVersion int) ([]Event, error)
    GetAllEvents() ([]Event, error)
}

// 内存事件存储实现
type InMemoryEventStore struct {
    events map[string][]Event
    mu     sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
    return &InMemoryEventStore{
        events: make(map[string][]Event),
    }
}

func (es *InMemoryEventStore) SaveEvents(aggregateID string, events []Event, expectedVersion int) error {
    es.mu.Lock()
    defer es.mu.Unlock()
    
    currentEvents := es.events[aggregateID]
    if len(currentEvents) != expectedVersion {
        return fmt.Errorf("concurrent modification detected")
    }
    
    es.events[aggregateID] = append(currentEvents, events...)
    return nil
}

func (es *InMemoryEventStore) GetEvents(aggregateID string) ([]Event, error) {
    es.mu.RLock()
    defer es.mu.RUnlock()
    
    events, exists := es.events[aggregateID]
    if !exists {
        return nil, fmt.Errorf("aggregate not found")
    }
    
    return events, nil
}
```

## 事件模型

### 1. 领域事件

```go
// 用户创建事件
type UserCreatedEvent struct {
    BaseEvent
    Username string `json:"username"`
    Email    string `json:"email"`
}

func NewUserCreatedEvent(aggregateID, username, email string) *UserCreatedEvent {
    return &UserCreatedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "UserCreated",
            Timestamp:   time.Now(),
            AggregateID: aggregateID,
            Version:     1,
        },
        Username: username,
        Email:    email,
    }
}

// 用户更新事件
type UserUpdatedEvent struct {
    BaseEvent
    Username string `json:"username"`
    Email    string `json:"email"`
}

func NewUserUpdatedEvent(aggregateID, username, email string, version int) *UserUpdatedEvent {
    return &UserUpdatedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "UserUpdated",
            Timestamp:   time.Now(),
            AggregateID: aggregateID,
            Version:     version,
        },
        Username: username,
        Email:    email,
    }
}
```

### 2. 聚合实现

```go
// 用户聚合
type User struct {
    BaseAggregate
    Username string
    Email    string
    Active   bool
}

func NewUser(id, username, email string) *User {
    user := &User{
        BaseAggregate: BaseAggregate{id: id},
        Username:      username,
        Email:         email,
        Active:        true,
    }
    
    event := NewUserCreatedEvent(id, username, email)
    user.AddEvent(event)
    
    return user
}

func (u *User) Update(username, email string) error {
    if username == "" || email == "" {
        return fmt.Errorf("username and email cannot be empty")
    }
    
    event := NewUserUpdatedEvent(u.ID(), username, email, u.Version()+1)
    u.AddEvent(event)
    
    return nil
}

func (u *User) Deactivate() {
    event := &UserDeactivatedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "UserDeactivated",
            Timestamp:   time.Now(),
            AggregateID: u.ID(),
            Version:     u.Version() + 1,
        },
    }
    u.AddEvent(event)
}

func (u *User) Apply(event Event) error {
    switch e := event.(type) {
    case *UserCreatedEvent:
        u.Username = e.Username
        u.Email = e.Email
        u.Active = true
    case *UserUpdatedEvent:
        u.Username = e.Username
        u.Email = e.Email
    case *UserDeactivatedEvent:
        u.Active = false
    default:
        return fmt.Errorf("unknown event type: %s", event.EventType())
    }
    
    u.version = event.Version()
    return nil
}
```

## 事件存储

### 1. 数据库事件存储

```go
// PostgreSQL事件存储
type PostgresEventStore struct {
    db *sql.DB
}

func NewPostgresEventStore(connStr string) (*PostgresEventStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // 创建事件表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS events (
            id VARCHAR(255) PRIMARY KEY,
            aggregate_id VARCHAR(255) NOT NULL,
            event_type VARCHAR(255) NOT NULL,
            event_data JSONB NOT NULL,
            version INTEGER NOT NULL,
            timestamp TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
        CREATE INDEX IF NOT EXISTS idx_events_aggregate_id ON events(aggregate_id);
        CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
    `)
    
    return &PostgresEventStore{db: db}, err
}

func (es *PostgresEventStore) SaveEvents(aggregateID string, events []Event, expectedVersion int) error {
    tx, err := es.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 检查版本
    var currentVersion int
    err = tx.QueryRow("SELECT COALESCE(MAX(version), 0) FROM events WHERE aggregate_id = $1", aggregateID).Scan(&currentVersion)
    if err != nil {
        return err
    }
    
    if currentVersion != expectedVersion {
        return fmt.Errorf("concurrent modification detected")
    }
    
    // 插入事件
    stmt, err := tx.Prepare(`
        INSERT INTO events (id, aggregate_id, event_type, event_data, version, timestamp)
        VALUES (```latex
$1, $
```2, ```latex
$3, $
```4, ```latex
$5, $
```6)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, event := range events {
        eventData, err := json.Marshal(event)
        if err != nil {
            return err
        }
        
        _, err = stmt.Exec(
            event.EventID(),
            event.AggregateID(),
            event.EventType(),
            eventData,
            event.Version(),
            event.Timestamp(),
        )
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (es *PostgresEventStore) GetEvents(aggregateID string) ([]Event, error) {
    rows, err := es.db.Query(`
        SELECT event_data FROM events 
        WHERE aggregate_id = $1 
        ORDER BY version ASC
    `, aggregateID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var events []Event
    for rows.Next() {
        var eventData []byte
        if err := rows.Scan(&eventData); err != nil {
            return nil, err
        }
        
        var eventMap map[string]interface{}
        if err := json.Unmarshal(eventData, &eventMap); err != nil {
            return nil, err
        }
        
        event, err := es.deserializeEvent(eventMap)
        if err != nil {
            return nil, err
        }
        
        events = append(events, event)
    }
    
    return events, nil
}
```

### 2. 事件序列化

```go
// 事件序列化器
type EventSerializer interface {
    Serialize(event Event) ([]byte, error)
    Deserialize(data []byte) (Event, error)
}

// JSON事件序列化器
type JSONEventSerializer struct {
    eventTypes map[string]reflect.Type
}

func NewJSONEventSerializer() *JSONEventSerializer {
    return &JSONEventSerializer{
        eventTypes: make(map[string]reflect.Type),
    }
}

func (s *JSONEventSerializer) RegisterEventType(eventType string, eventStruct interface{}) {
    s.eventTypes[eventType] = reflect.TypeOf(eventStruct)
}

func (s *JSONEventSerializer) Serialize(event Event) ([]byte, error) {
    return json.Marshal(event)
}

func (s *JSONEventSerializer) Deserialize(data []byte) (Event, error) {
    var eventMap map[string]interface{}
    if err := json.Unmarshal(data, &eventMap); err != nil {
        return nil, err
    }
    
    eventType, ok := eventMap["type"].(string)
    if !ok {
        return nil, fmt.Errorf("event type not found")
    }
    
    eventStructType, exists := s.eventTypes[eventType]
    if !exists {
        return nil, fmt.Errorf("unknown event type: %s", eventType)
    }
    
    eventStruct := reflect.New(eventStructType).Interface()
    if err := json.Unmarshal(data, eventStruct); err != nil {
        return nil, err
    }
    
    return eventStruct.(Event), nil
}
```

## 事件流处理

### 1. 事件处理器

```go
// 事件处理器接口
type EventHandler interface {
    Handle(event Event) error
}

// 用户事件处理器
type UserEventHandler struct {
    userRepo UserRepository
}

func NewUserEventHandler(userRepo UserRepository) *UserEventHandler {
    return &UserEventHandler{userRepo: userRepo}
}

func (h *UserEventHandler) Handle(event Event) error {
    switch e := event.(type) {
    case *UserCreatedEvent:
        return h.handleUserCreated(e)
    case *UserUpdatedEvent:
        return h.handleUserUpdated(e)
    case *UserDeactivatedEvent:
        return h.handleUserDeactivated(e)
    default:
        return fmt.Errorf("unknown event type: %s", event.EventType())
    }
}

func (h *UserEventHandler) handleUserCreated(event *UserCreatedEvent) error {
    user := &UserProjection{
        ID:       event.AggregateID(),
        Username: event.Username,
        Email:    event.Email,
        Active:   true,
        CreatedAt: event.Timestamp(),
    }
    
    return h.userRepo.Save(user)
}

func (h *UserEventHandler) handleUserUpdated(event *UserUpdatedEvent) error {
    user, err := h.userRepo.GetByID(event.AggregateID())
    if err != nil {
        return err
    }
    
    user.Username = event.Username
    user.Email = event.Email
    user.UpdatedAt = event.Timestamp()
    
    return h.userRepo.Save(user)
}

func (h *UserEventHandler) handleUserDeactivated(event *UserDeactivatedEvent) error {
    user, err := h.userRepo.GetByID(event.AggregateID())
    if err != nil {
        return err
    }
    
    user.Active = false
    user.UpdatedAt = event.Timestamp()
    
    return h.userRepo.Save(user)
}
```

### 2. 事件总线

```go
// 事件总线
type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Publish(event Event) error {
    eb.mu.RLock()
    handlers := eb.handlers[event.EventType()]
    eb.mu.RUnlock()
    
    for _, handler := range handlers {
        if err := handler.Handle(event); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 快照机制

### 1. 快照接口

```go
// 快照接口
type Snapshot interface {
    AggregateID() string
    Version() int
    State() interface{}
    Timestamp() time.Time
}

// 用户快照
type UserSnapshot struct {
    AggregateID string    `json:"aggregate_id"`
    Version     int       `json:"version"`
    State       UserState `json:"state"`
    Timestamp   time.Time `json:"timestamp"`
}

type UserState struct {
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (s UserSnapshot) AggregateID() string {
    return s.AggregateID
}

func (s UserSnapshot) Version() int {
    return s.Version
}

func (s UserSnapshot) State() interface{} {
    return s.State
}

func (s UserSnapshot) Timestamp() time.Time {
    return s.Timestamp
}
```

### 2. 快照存储

```go
// 快照存储接口
type SnapshotStore interface {
    Save(snapshot Snapshot) error
    Get(aggregateID string) (Snapshot, error)
    GetLatest(aggregateID string) (Snapshot, error)
}

// 内存快照存储
type InMemorySnapshotStore struct {
    snapshots map[string][]Snapshot
    mu        sync.RWMutex
}

func NewInMemorySnapshotStore() *InMemorySnapshotStore {
    return &InMemorySnapshotStore{
        snapshots: make(map[string][]Snapshot),
    }
}

func (ss *InMemorySnapshotStore) Save(snapshot Snapshot) error {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    
    ss.snapshots[snapshot.AggregateID()] = append(
        ss.snapshots[snapshot.AggregateID()], 
        snapshot,
    )
    
    return nil
}

func (ss *InMemorySnapshotStore) GetLatest(aggregateID string) (Snapshot, error) {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    
    snapshots, exists := ss.snapshots[aggregateID]
    if !exists || len(snapshots) == 0 {
        return nil, fmt.Errorf("no snapshot found")
    }
    
    return snapshots[len(snapshots)-1], nil
}
```

## 投影模式

### 1. 投影接口

```go
// 投影接口
type Projection interface {
    ID() string
    Handle(event Event) error
    GetState() interface{}
}

// 用户投影
type UserProjection struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (p *UserProjection) ID() string {
    return p.ID
}

func (p *UserProjection) Handle(event Event) error {
    switch e := event.(type) {
    case *UserCreatedEvent:
        p.ID = e.AggregateID()
        p.Username = e.Username
        p.Email = e.Email
        p.Active = true
        p.CreatedAt = e.Timestamp()
    case *UserUpdatedEvent:
        p.Username = e.Username
        p.Email = e.Email
        p.UpdatedAt = e.Timestamp()
    case *UserDeactivatedEvent:
        p.Active = false
        p.UpdatedAt = e.Timestamp()
    }
    
    return nil
}

func (p *UserProjection) GetState() interface{} {
    return p
}
```

### 2. 投影存储

```go
// 投影存储接口
type ProjectionStore interface {
    Save(projection Projection) error
    Get(id string) (Projection, error)
    GetAll() ([]Projection, error)
    Delete(id string) error
}

// 内存投影存储
type InMemoryProjectionStore struct {
    projections map[string]Projection
    mu          sync.RWMutex
}

func NewInMemoryProjectionStore() *InMemoryProjectionStore {
    return &InMemoryProjectionStore{
        projections: make(map[string]Projection),
    }
}

func (ps *InMemoryProjectionStore) Save(projection Projection) error {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    ps.projections[projection.ID()] = projection
    return nil
}

func (ps *InMemoryProjectionStore) Get(id string) (Projection, error) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    
    projection, exists := ps.projections[id]
    if !exists {
        return nil, fmt.Errorf("projection not found")
    }
    
    return projection, nil
}
```

## 并发控制

### 1. 乐观并发控制

```go
// 乐观锁实现
type OptimisticLock struct {
    version int
    mu      sync.Mutex
}

func (ol *OptimisticLock) GetVersion() int {
    ol.mu.Lock()
    defer ol.mu.Unlock()
    return ol.version
}

func (ol *OptimisticLock) IncrementVersion() {
    ol.mu.Lock()
    defer ol.mu.Unlock()
    ol.version++
}

// 带乐观锁的聚合根
type OptimisticAggregate struct {
    BaseAggregate
    lock OptimisticLock
}

func (oa *OptimisticAggregate) GetVersion() int {
    return oa.lock.GetVersion()
}

func (oa *OptimisticAggregate) AddEvent(event Event) {
    oa.BaseAggregate.AddEvent(event)
    oa.lock.IncrementVersion()
}
```

### 2. 事件版本控制

```go
// 事件版本检查器
type EventVersionChecker struct {
    eventStore EventStore
}

func NewEventVersionChecker(eventStore EventStore) *EventVersionChecker {
    return &EventVersionChecker{eventStore: eventStore}
}

func (evc *EventVersionChecker) CheckVersion(aggregateID string, expectedVersion int) error {
    events, err := evc.eventStore.GetEvents(aggregateID)
    if err != nil {
        return err
    }
    
    if len(events) != expectedVersion {
        return fmt.Errorf("version mismatch: expected %d, got %d", expectedVersion, len(events))
    }
    
    return nil
}
```

## 实际应用

### 1. 银行账户系统

```go
// 银行账户聚合
type BankAccount struct {
    BaseAggregate
    AccountNumber string
    Balance       decimal.Decimal
    Owner         string
    Active        bool
}

// 账户创建事件
type AccountCreatedEvent struct {
    BaseEvent
    AccountNumber string          `json:"account_number"`
    InitialBalance decimal.Decimal `json:"initial_balance"`
    Owner         string          `json:"owner"`
}

// 存款事件
type MoneyDepositedEvent struct {
    BaseEvent
    Amount decimal.Decimal `json:"amount"`
}

// 取款事件
type MoneyWithdrawnEvent struct {
    BaseEvent
    Amount decimal.Decimal `json:"amount"`
}

// 账户操作
func (ba *BankAccount) Deposit(amount decimal.Decimal) error {
    if amount.LessThanOrEqual(decimal.Zero) {
        return fmt.Errorf("deposit amount must be positive")
    }
    
    event := &MoneyDepositedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "MoneyDeposited",
            Timestamp:   time.Now(),
            AggregateID: ba.ID(),
            Version:     ba.Version() + 1,
        },
        Amount: amount,
    }
    
    ba.AddEvent(event)
    return nil
}

func (ba *BankAccount) Withdraw(amount decimal.Decimal) error {
    if amount.LessThanOrEqual(decimal.Zero) {
        return fmt.Errorf("withdrawal amount must be positive")
    }
    
    if ba.Balance.LessThan(amount) {
        return fmt.Errorf("insufficient funds")
    }
    
    event := &MoneyWithdrawnEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "MoneyWithdrawn",
            Timestamp:   time.Now(),
            AggregateID: ba.ID(),
            Version:     ba.Version() + 1,
        },
        Amount: amount,
    }
    
    ba.AddEvent(event)
    return nil
}

func (ba *BankAccount) Apply(event Event) error {
    switch e := event.(type) {
    case *AccountCreatedEvent:
        ba.AccountNumber = e.AccountNumber
        ba.Balance = e.InitialBalance
        ba.Owner = e.Owner
        ba.Active = true
    case *MoneyDepositedEvent:
        ba.Balance = ba.Balance.Add(e.Amount)
    case *MoneyWithdrawnEvent:
        ba.Balance = ba.Balance.Sub(e.Amount)
    }
    
    ba.version = event.Version()
    return nil
}
```

### 2. 订单管理系统

```go
// 订单聚合
type Order struct {
    BaseAggregate
    OrderID      string
    CustomerID   string
    Items        []OrderItem
    Status       OrderStatus
    TotalAmount  decimal.Decimal
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type OrderItem struct {
    ProductID string          `json:"product_id"`
    Quantity  int             `json:"quantity"`
    Price     decimal.Decimal `json:"price"`
}

type OrderStatus string

const (
    OrderStatusCreated   OrderStatus = "created"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
)

// 订单事件
type OrderCreatedEvent struct {
    BaseEvent
    CustomerID string      `json:"customer_id"`
    Items      []OrderItem `json:"items"`
    TotalAmount decimal.Decimal `json:"total_amount"`
}

type OrderConfirmedEvent struct {
    BaseEvent
    ConfirmedAt time.Time `json:"confirmed_at"`
}

type OrderShippedEvent struct {
    BaseEvent
    TrackingNumber string    `json:"tracking_number"`
    ShippedAt      time.Time `json:"shipped_at"`
}

// 订单操作
func (o *Order) Confirm() error {
    if o.Status != OrderStatusCreated {
        return fmt.Errorf("order cannot be confirmed in current status")
    }
    
    event := &OrderConfirmedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "OrderConfirmed",
            Timestamp:   time.Now(),
            AggregateID: o.ID(),
            Version:     o.Version() + 1,
        },
        ConfirmedAt: time.Now(),
    }
    
    o.AddEvent(event)
    return nil
}

func (o *Order) Ship(trackingNumber string) error {
    if o.Status != OrderStatusConfirmed {
        return fmt.Errorf("order cannot be shipped in current status")
    }
    
    event := &OrderShippedEvent{
        BaseEvent: BaseEvent{
            ID:          uuid.New().String(),
            Type:        "OrderShipped",
            Timestamp:   time.Now(),
            AggregateID: o.ID(),
            Version:     o.Version() + 1,
        },
        TrackingNumber: trackingNumber,
        ShippedAt:      time.Now(),
    }
    
    o.AddEvent(event)
    return nil
}
```

## 最佳实践

### 1. 事件设计原则

- **事件命名**：使用过去时态，描述已发生的事实
- **事件不可变性**：事件一旦创建就不能修改
- **事件粒度**：保持事件粒度适中，既不过细也不过粗
- **事件版本化**：为事件结构提供版本控制机制

### 2. 聚合设计原则

- **聚合边界**：明确定义聚合的边界和职责
- **业务规则**：在聚合中实现业务规则和约束
- **事件产生**：聚合负责产生领域事件
- **状态一致性**：确保聚合内部状态的一致性

### 3. 性能优化

```go
// 批量事件处理
type BatchEventProcessor struct {
    eventStore EventStore
    batchSize  int
    processor  EventHandler
}

func (bep *BatchEventProcessor) ProcessBatch() error {
    events, err := bep.eventStore.GetAllEvents()
    if err != nil {
        return err
    }
    
    for i := 0; i < len(events); i += bep.batchSize {
        end := i + bep.batchSize
        if end > len(events) {
            end = len(events)
        }
        
        batch := events[i:end]
        for _, event := range batch {
            if err := bep.processor.Handle(event); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### 4. 错误处理

```go
// 事件处理错误
type EventHandlingError struct {
    EventID   string
    EventType string
    Error     error
    Retryable bool
}

func (ehe EventHandlingError) Error() string {
    return fmt.Sprintf("failed to handle event %s (%s): %v", ehe.EventID, ehe.EventType, ehe.Error)
}

// 重试机制
type RetryableEventHandler struct {
    handler EventHandler
    maxRetries int
    backoff    time.Duration
}

func (reh *RetryableEventHandler) Handle(event Event) error {
    var lastErr error
    
    for i := 0; i <= reh.maxRetries; i++ {
        if err := reh.handler.Handle(event); err != nil {
            lastErr = err
            if i < reh.maxRetries {
                time.Sleep(reh.backoff * time.Duration(i+1))
                continue
            }
        } else {
            return nil
        }
    }
    
    return lastErr
}
```

## 总结

事件溯源模式提供了一种强大的数据持久化方法，通过记录事件序列来维护系统状态。这种模式具有以下优势：

- **审计能力**：完整的事件历史记录
- **时间旅行**：可以重建任意时间点的状态
- **解耦性**：事件发布者和订阅者解耦
- **可扩展性**：支持多种投影和查询模型

关键要点：

- 事件是不可变的事实记录
- 聚合根负责产生事件和应用事件
- 事件存储提供持久化能力
- 投影模式支持多种查询需求
- 快照机制优化性能
- 并发控制确保数据一致性

事件溯源模式特别适用于需要完整审计历史、复杂业务逻辑和多种查询模式的系统。

---

**相关链接**：

- [01-响应式模式](./01-Reactive-Patterns.md)
- [02-函数式编程模式](./02-Functional-Patterns.md)
- [04-CQRS模式](./04-CQRS-Patterns.md)
- [../README.md](../README.md)
