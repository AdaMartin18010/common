# 4. 工作流集成

## 4.1 集成理论基础

### 4.1.1 集成模式定义

**定义 4.1** (工作流集成): 工作流集成是一个三元组 $\mathcal{I} = (W, S, A)$，其中：

- $W$ 是工作流系统
- $S$ 是外部系统集合
- $A$ 是适配器集合

**集成架构**:

```latex
\text{WorkflowIntegration} = \text{AdapterLayer} \times \text{ProtocolLayer} \times \text{DataLayer}
```

### 4.1.2 集成模式分类

```latex
\text{IntegrationPatterns} = \text{Synchronous} \cup \text{Asynchronous} \cup \text{EventDriven} \cup \text{MessageBased}
```

## 4.2 系统集成实现

### 4.2.1 集成适配器

```go
// IntegrationAdapter 集成适配器接口
type IntegrationAdapter interface {
    // 适配器标识
    GetID() string
    GetName() string
    GetType() AdapterType
    
    // 连接管理
    Connect() error
    Disconnect() error
    IsConnected() bool
    
    // 数据操作
    Send(data interface{}) error
    Receive() (<-chan interface{}, error)
    
    // 配置管理
    GetConfig() AdapterConfig
    UpdateConfig(config AdapterConfig) error
}

// AdapterType 适配器类型
type AdapterType string

const (
    AdapterTypeDatabase  AdapterType = "DATABASE"
    AdapterTypeMessage   AdapterType = "MESSAGE"
    AdapterTypeAPI       AdapterType = "API"
    AdapterTypeFile      AdapterType = "FILE"
    AdapterTypeLegacy    AdapterType = "LEGACY"
)

// AdapterConfig 适配器配置
type AdapterConfig struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        AdapterType            `json:"type"`
    Connection  ConnectionConfig       `json:"connection"`
    Mapping     DataMapping            `json:"mapping"`
    Parameters  map[string]interface{} `json:"parameters"`
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
    URL      string            `json:"url"`
    Username string            `json:"username"`
    Password string            `json:"password"`
    Timeout  time.Duration     `json:"timeout"`
    Headers  map[string]string `json:"headers"`
}

// DataMapping 数据映射
type DataMapping struct {
    InputMapping  map[string]string `json:"input_mapping"`
    OutputMapping map[string]string `json:"output_mapping"`
    Transformers  []Transformer     `json:"transformers"`
}

// Transformer 数据转换器
type Transformer struct {
    Type       string                 `json:"type"`
    Parameters map[string]interface{} `json:"parameters"`
}

// BaseAdapter 基础适配器实现
type BaseAdapter struct {
    id       string
    name     string
    adapterType AdapterType
    config   AdapterConfig
    connected bool
    mutex    sync.RWMutex
}

// NewBaseAdapter 创建基础适配器
func NewBaseAdapter(id, name string, adapterType AdapterType) *BaseAdapter {
    return &BaseAdapter{
        id:         id,
        name:       name,
        adapterType: adapterType,
    }
}

// GetID 获取适配器ID
func (ba *BaseAdapter) GetID() string {
    return ba.id
}

// GetName 获取适配器名称
func (ba *BaseAdapter) GetName() string {
    return ba.name
}

// GetType 获取适配器类型
func (ba *BaseAdapter) GetType() AdapterType {
    return ba.adapterType
}

// IsConnected 检查连接状态
func (ba *BaseAdapter) IsConnected() bool {
    ba.mutex.RLock()
    defer ba.mutex.RUnlock()
    return ba.connected
}

// GetConfig 获取适配器配置
func (ba *BaseAdapter) GetConfig() AdapterConfig {
    ba.mutex.RLock()
    defer ba.mutex.RUnlock()
    return ba.config
}

// UpdateConfig 更新适配器配置
func (ba *BaseAdapter) UpdateConfig(config AdapterConfig) error {
    ba.mutex.Lock()
    defer ba.mutex.Unlock()
    ba.config = config
    return nil
}
```

### 4.2.2 数据库适配器

```go
// DatabaseAdapter 数据库适配器
type DatabaseAdapter struct {
    *BaseAdapter
    db     *sql.DB
    driver string
}

// NewDatabaseAdapter 创建数据库适配器
func NewDatabaseAdapter(id, name, driver string) *DatabaseAdapter {
    return &DatabaseAdapter{
        BaseAdapter: NewBaseAdapter(id, name, AdapterTypeDatabase),
        driver:      driver,
    }
}

// Connect 连接数据库
func (da *DatabaseAdapter) Connect() error {
    da.mutex.Lock()
    defer da.mutex.Unlock()
    
    if da.connected {
        return nil
    }
    
    db, err := sql.Open(da.driver, da.config.Connection.URL)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    // 测试连接
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    da.db = db
    da.connected = true
    return nil
}

// Disconnect 断开数据库连接
func (da *DatabaseAdapter) Disconnect() error {
    da.mutex.Lock()
    defer da.mutex.Unlock()
    
    if !da.connected {
        return nil
    }
    
    if da.db != nil {
        if err := da.db.Close(); err != nil {
            return fmt.Errorf("failed to close database: %w", err)
        }
    }
    
    da.connected = false
    return nil
}

// Send 发送数据到数据库
func (da *DatabaseAdapter) Send(data interface{}) error {
    da.mutex.RLock()
    defer da.mutex.RUnlock()
    
    if !da.connected {
        return fmt.Errorf("not connected to database")
    }
    
    // 根据数据类型执行不同的操作
    switch v := data.(type) {
    case *InsertData:
        return da.insertData(v)
    case *UpdateData:
        return da.updateData(v)
    case *DeleteData:
        return da.deleteData(v)
    case *QueryData:
        return da.queryData(v)
    default:
        return fmt.Errorf("unsupported data type: %T", data)
    }
}

// InsertData 插入数据
type InsertData struct {
    Table   string                 `json:"table"`
    Columns []string               `json:"columns"`
    Values  []interface{}          `json:"values"`
    Metadata map[string]interface{} `json:"metadata"`
}

// insertData 插入数据
func (da *DatabaseAdapter) insertData(data *InsertData) error {
    // 构建INSERT语句
    placeholders := make([]string, len(data.Values))
    for i := range data.Values {
        placeholders[i] = "?"
    }
    
    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)",
        data.Table,
        strings.Join(data.Columns, ", "),
        strings.Join(placeholders, ", "),
    )
    
    // 执行插入
    _, err := da.db.Exec(query, data.Values...)
    if err != nil {
        return fmt.Errorf("failed to insert data: %w", err)
    }
    
    return nil
}

// UpdateData 更新数据
type UpdateData struct {
    Table   string                 `json:"table"`
    Set     map[string]interface{} `json:"set"`
    Where   map[string]interface{} `json:"where"`
    Metadata map[string]interface{} `json:"metadata"`
}

// updateData 更新数据
func (da *DatabaseAdapter) updateData(data *UpdateData) error {
    // 构建SET子句
    setClause := make([]string, 0, len(data.Set))
    values := make([]interface{}, 0, len(data.Set))
    
    for column, value := range data.Set {
        setClause = append(setClause, fmt.Sprintf("%s = ?", column))
        values = append(values, value)
    }
    
    // 构建WHERE子句
    whereClause := make([]string, 0, len(data.Where))
    for column, value := range data.Where {
        whereClause = append(whereClause, fmt.Sprintf("%s = ?", column))
        values = append(values, value)
    }
    
    query := fmt.Sprintf(
        "UPDATE %s SET %s WHERE %s",
        data.Table,
        strings.Join(setClause, ", "),
        strings.Join(whereClause, " AND "),
    )
    
    // 执行更新
    _, err := da.db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to update data: %w", err)
    }
    
    return nil
}

// DeleteData 删除数据
type DeleteData struct {
    Table   string                 `json:"table"`
    Where   map[string]interface{} `json:"where"`
    Metadata map[string]interface{} `json:"metadata"`
}

// deleteData 删除数据
func (da *DatabaseAdapter) deleteData(data *DeleteData) error {
    // 构建WHERE子句
    whereClause := make([]string, 0, len(data.Where))
    values := make([]interface{}, 0, len(data.Where))
    
    for column, value := range data.Where {
        whereClause = append(whereClause, fmt.Sprintf("%s = ?", column))
        values = append(values, value)
    }
    
    query := fmt.Sprintf(
        "DELETE FROM %s WHERE %s",
        data.Table,
        strings.Join(whereClause, " AND "),
    )
    
    // 执行删除
    _, err := da.db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to delete data: %w", err)
    }
    
    return nil
}

// QueryData 查询数据
type QueryData struct {
    Table   string                 `json:"table"`
    Columns []string               `json:"columns"`
    Where   map[string]interface{} `json:"where"`
    OrderBy string                 `json:"order_by"`
    Limit   int                    `json:"limit"`
    Metadata map[string]interface{} `json:"metadata"`
}

// queryData 查询数据
func (da *DatabaseAdapter) queryData(data *QueryData) error {
    // 构建SELECT语句
    columns := "*"
    if len(data.Columns) > 0 {
        columns = strings.Join(data.Columns, ", ")
    }
    
    query := fmt.Sprintf("SELECT %s FROM %s", columns, data.Table)
    values := make([]interface{}, 0)
    
    // 添加WHERE子句
    if len(data.Where) > 0 {
        whereClause := make([]string, 0, len(data.Where))
        for column, value := range data.Where {
            whereClause = append(whereClause, fmt.Sprintf("%s = ?", column))
            values = append(values, value)
        }
        query += " WHERE " + strings.Join(whereClause, " AND ")
    }
    
    // 添加ORDER BY子句
    if data.OrderBy != "" {
        query += " ORDER BY " + data.OrderBy
    }
    
    // 添加LIMIT子句
    if data.Limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", data.Limit)
    }
    
    // 执行查询
    rows, err := da.db.Query(query, values...)
    if err != nil {
        return fmt.Errorf("failed to query data: %w", err)
    }
    defer rows.Close()
    
    // 处理结果集
    // 这里可以根据需要返回结果数据
    return nil
}

// Receive 接收数据库事件
func (da *DatabaseAdapter) Receive() (<-chan interface{}, error) {
    // 数据库适配器通常不主动接收数据
    // 可以通过轮询或触发器实现
    return nil, fmt.Errorf("database adapter does not support receive")
}
```

### 4.2.3 消息队列适配器

```go
// MessageQueueAdapter 消息队列适配器
type MessageQueueAdapter struct {
    *BaseAdapter
    producer MessageProducer
    consumer MessageConsumer
    queue    string
}

// MessageProducer 消息生产者接口
type MessageProducer interface {
    SendMessage(topic string, message []byte) error
    Close() error
}

// MessageConsumer 消息消费者接口
type MessageConsumer interface {
    ConsumeMessages(topic string) (<-chan Message, error)
    Close() error
}

// Message 消息结构
type Message struct {
    ID       string                 `json:"id"`
    Topic    string                 `json:"topic"`
    Payload  []byte                 `json:"payload"`
    Headers  map[string]string      `json:"headers"`
    Timestamp time.Time             `json:"timestamp"`
    Metadata map[string]interface{} `json:"metadata"`
}

// NewMessageQueueAdapter 创建消息队列适配器
func NewMessageQueueAdapter(id, name, queue string) *MessageQueueAdapter {
    return &MessageQueueAdapter{
        BaseAdapter: NewBaseAdapter(id, name, AdapterTypeMessage),
        queue:       queue,
    }
}

// Connect 连接消息队列
func (mqa *MessageQueueAdapter) Connect() error {
    mqa.mutex.Lock()
    defer mqa.mutex.Unlock()
    
    if mqa.connected {
        return nil
    }
    
    // 根据配置创建生产者和消费者
    // 这里以Redis为例
    if err := mqa.createRedisConnection(); err != nil {
        return fmt.Errorf("failed to create Redis connection: %w", err)
    }
    
    mqa.connected = true
    return nil
}

// Disconnect 断开消息队列连接
func (mqa *MessageQueueAdapter) Disconnect() error {
    mqa.mutex.Lock()
    defer mqa.mutex.Unlock()
    
    if !mqa.connected {
        return nil
    }
    
    if mqa.producer != nil {
        if err := mqa.producer.Close(); err != nil {
            return fmt.Errorf("failed to close producer: %w", err)
        }
    }
    
    if mqa.consumer != nil {
        if err := mqa.consumer.Close(); err != nil {
            return fmt.Errorf("failed to close consumer: %w", err)
        }
    }
    
    mqa.connected = false
    return nil
}

// Send 发送消息
func (mqa *MessageQueueAdapter) Send(data interface{}) error {
    mqa.mutex.RLock()
    defer mqa.mutex.RUnlock()
    
    if !mqa.connected {
        return fmt.Errorf("not connected to message queue")
    }
    
    // 序列化数据
    payload, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }
    
    // 发送消息
    if err := mqa.producer.SendMessage(mqa.queue, payload); err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    
    return nil
}

// Receive 接收消息
func (mqa *MessageQueueAdapter) Receive() (<-chan interface{}, error) {
    mqa.mutex.RLock()
    defer mqa.mutex.RUnlock()
    
    if !mqa.connected {
        return nil, fmt.Errorf("not connected to message queue")
    }
    
    // 创建消息通道
    messageChan := make(chan interface{}, 100)
    
    // 启动消费者
    go func() {
        defer close(messageChan)
        
        messages, err := mqa.consumer.ConsumeMessages(mqa.queue)
        if err != nil {
            log.Printf("Failed to consume messages: %v", err)
            return
        }
        
        for msg := range messages {
            // 反序列化消息
            var data interface{}
            if err := json.Unmarshal(msg.Payload, &data); err != nil {
                log.Printf("Failed to unmarshal message: %v", err)
                continue
            }
            
            select {
            case messageChan <- data:
            default:
                log.Printf("Message channel full, dropping message")
            }
        }
    }()
    
    return messageChan, nil
}

// createRedisConnection 创建Redis连接
func (mqa *MessageQueueAdapter) createRedisConnection() error {
    // 这里实现Redis连接逻辑
    // 简化实现，实际项目中可以使用Redis客户端库
    return nil
}
```

## 4.3 工作流集成管理器

### 4.3.1 集成管理器实现

```go
// IntegrationManager 集成管理器
type IntegrationManager struct {
    adapters map[string]IntegrationAdapter
    mappings map[string]DataMapping
    mutex    sync.RWMutex
}

// NewIntegrationManager 创建集成管理器
func NewIntegrationManager() *IntegrationManager {
    return &IntegrationManager{
        adapters: make(map[string]IntegrationAdapter),
        mappings: make(map[string]DataMapping),
    }
}

// RegisterAdapter 注册适配器
func (im *IntegrationManager) RegisterAdapter(adapter IntegrationAdapter) error {
    im.mutex.Lock()
    defer im.mutex.Unlock()
    
    if _, exists := im.adapters[adapter.GetID()]; exists {
        return fmt.Errorf("adapter already registered: %s", adapter.GetID())
    }
    
    im.adapters[adapter.GetID()] = adapter
    return nil
}

// UnregisterAdapter 注销适配器
func (im *IntegrationManager) UnregisterAdapter(adapterID string) error {
    im.mutex.Lock()
    defer im.mutex.Unlock()
    
    adapter, exists := im.adapters[adapterID]
    if !exists {
        return fmt.Errorf("adapter not found: %s", adapterID)
    }
    
    // 断开连接
    if err := adapter.Disconnect(); err != nil {
        return fmt.Errorf("failed to disconnect adapter: %w", err)
    }
    
    delete(im.adapters, adapterID)
    return nil
}

// GetAdapter 获取适配器
func (im *IntegrationManager) GetAdapter(adapterID string) (IntegrationAdapter, error) {
    im.mutex.RLock()
    defer im.mutex.RUnlock()
    
    adapter, exists := im.adapters[adapterID]
    if !exists {
        return nil, fmt.Errorf("adapter not found: %s", adapterID)
    }
    
    return adapter, nil
}

// SendToAdapter 发送数据到适配器
func (im *IntegrationManager) SendToAdapter(adapterID string, data interface{}) error {
    adapter, err := im.GetAdapter(adapterID)
    if err != nil {
        return err
    }
    
    // 应用数据映射
    mappedData, err := im.applyMapping(adapterID, data)
    if err != nil {
        return fmt.Errorf("failed to apply mapping: %w", err)
    }
    
    return adapter.Send(mappedData)
}

// ReceiveFromAdapter 从适配器接收数据
func (im *IntegrationManager) ReceiveFromAdapter(adapterID string) (<-chan interface{}, error) {
    adapter, err := im.GetAdapter(adapterID)
    if err != nil {
        return nil, err
    }
    
    return adapter.Receive()
}

// applyMapping 应用数据映射
func (im *IntegrationManager) applyMapping(adapterID string, data interface{}) (interface{}, error) {
    im.mutex.RLock()
    mapping, exists := im.mappings[adapterID]
    im.mutex.RUnlock()
    
    if !exists {
        return data, nil
    }
    
    // 应用输入映射
    if len(mapping.InputMapping) > 0 {
        data = im.mapInputData(data, mapping.InputMapping)
    }
    
    // 应用转换器
    for _, transformer := range mapping.Transformers {
        data = im.applyTransformer(data, transformer)
    }
    
    return data, nil
}

// mapInputData 映射输入数据
func (im *IntegrationManager) mapInputData(data interface{}, mapping map[string]string) interface{} {
    // 实现数据映射逻辑
    // 这里可以根据映射规则转换数据结构
    return data
}

// applyTransformer 应用转换器
func (im *IntegrationManager) applyTransformer(data interface{}, transformer Transformer) interface{} {
    switch transformer.Type {
    case "json":
        return im.transformJSON(data, transformer.Parameters)
    case "xml":
        return im.transformXML(data, transformer.Parameters)
    case "csv":
        return im.transformCSV(data, transformer.Parameters)
    default:
        return data
    }
}

// transformJSON JSON转换
func (im *IntegrationManager) transformJSON(data interface{}, params map[string]interface{}) interface{} {
    // 实现JSON转换逻辑
    return data
}

// transformXML XML转换
func (im *IntegrationManager) transformXML(data interface{}, params map[string]interface{}) interface{} {
    // 实现XML转换逻辑
    return data
}

// transformCSV CSV转换
func (im *IntegrationManager) transformCSV(data interface{}, params map[string]interface{}) interface{} {
    // 实现CSV转换逻辑
    return data
}
```

## 4.4 总结

工作流集成模块涵盖了以下核心内容：

1. **集成理论基础**: 形式化定义工作流集成的数学模型
2. **适配器模式**: 统一的适配器接口和基础实现
3. **系统集成**: 数据库、消息队列等具体系统的集成方案
4. **集成管理**: 统一的集成管理器，支持多系统集成

这个设计提供了一个完整的工作流集成框架，支持工作流系统与各种外部系统的无缝集成。 