# 数据库架构 (Database Architecture)

## 概述

数据库架构是数据库系统的整体设计，包括数据存储、访问、管理和分布的方式。它决定了数据库的性能、可扩展性、可用性和一致性特征，是现代应用程序的核心基础设施。

## 基本概念

### 核心特征

- **数据持久化**：可靠的数据存储和恢复
- **事务支持**：ACID特性的保证
- **并发控制**：多用户安全访问
- **查询优化**：高效的查询执行
- **索引管理**：快速数据检索
- **备份恢复**：数据安全和可用性

### 应用场景

- **关系型数据库**：结构化数据存储
- **NoSQL数据库**：非结构化数据存储
- **分布式数据库**：大规模数据分布
- **内存数据库**：高性能数据访问
- **时序数据库**：时间序列数据存储

## 核心组件

### 存储引擎 (Storage Engine)

```go
type Page struct {
    ID      uint32
    Data    []byte
    IsDirty bool
}

type BufferPool struct {
    pages    map[uint32]*Page
    capacity int
}

func (bp *BufferPool) GetPage(pageID uint32) (*Page, error) {
    if page, exists := bp.pages[pageID]; exists {
        return page, nil
    }
    
    // 从磁盘读取页面
    page, err := bp.readPageFromDisk(pageID)
    if err != nil {
        return nil, err
    }
    
    if len(bp.pages) >= bp.capacity {
        bp.evictPage()
    }
    
    bp.pages[pageID] = page
    return page, nil
}

func (bp *BufferPool) FlushPage(pageID uint32) error {
    page, exists := bp.pages[pageID]
    if !exists || !page.IsDirty {
        return nil
    }
    
    return bp.writePageToDisk(page)
}
```

### 查询处理器 (Query Processor)

```go
type QueryPlan struct {
    Type    string
    Table   string
    Columns []string
    Where   *Condition
}

type Condition struct {
    Column   string
    Operator string
    Value    interface{}
}

type QueryProcessor struct {
    storageEngine *StorageEngine
}

func (qp *QueryProcessor) ParseQuery(sql string) (*QueryPlan, error) {
    plan := &QueryPlan{}
    
    if strings.HasPrefix(strings.ToUpper(sql), "SELECT") {
        plan.Type = "SELECT"
        
        // 解析列
        columnsStart := strings.Index(sql, "SELECT") + 6
        columnsEnd := strings.Index(sql, "FROM")
        columnsStr := strings.TrimSpace(sql[columnsStart:columnsEnd])
        
        if columnsStr == "*" {
            plan.Columns = []string{"*"}
        } else {
            plan.Columns = strings.Split(columnsStr, ",")
        }
        
        // 解析表名
        fromStart := columnsEnd + 4
        whereIndex := strings.Index(strings.ToUpper(sql), "WHERE")
        if whereIndex == -1 {
            plan.Table = strings.TrimSpace(sql[fromStart:])
        } else {
            plan.Table = strings.TrimSpace(sql[fromStart:whereIndex])
            plan.Where = qp.parseWhereCondition(sql[whereIndex+5:])
        }
    }
    
    return plan, nil
}

func (qp *QueryProcessor) ExecuteQuery(plan *QueryPlan) (*ResultSet, error) {
    tableData, err := qp.storageEngine.ReadTable(plan.Table)
    if err != nil {
        return nil, err
    }
    
    result := &ResultSet{
        Columns: plan.Columns,
        Rows:    make([][]interface{}, 0),
    }
    
    for _, row := range tableData {
        if qp.evaluateCondition(row, plan.Where) {
            selectedRow := qp.selectColumns(row, plan.Columns)
            result.Rows = append(result.Rows, selectedRow)
        }
    }
    
    return result, nil
}
```

### 事务管理器 (Transaction Manager)

```go
type Transaction struct {
    ID        uint64
    Status    string
    StartTime time.Time
    Locks     map[string]*Lock
    Log       []*LogRecord
}

type Lock struct {
    Resource string
    Type     string
    TxnID    uint64
}

type TransactionManager struct {
    transactions map[uint64]*Transaction
    locks        map[string]*Lock
    nextTxnID    uint64
}

func (tm *TransactionManager) BeginTransaction() *Transaction {
    txn := &Transaction{
        ID:        tm.nextTxnID,
        Status:    "ACTIVE",
        StartTime: time.Now(),
        Locks:     make(map[string]*Lock),
        Log:       make([]*LogRecord, 0),
    }
    
    tm.nextTxnID++
    tm.transactions[txn.ID] = txn
    
    return txn
}

func (tm *TransactionManager) CommitTransaction(txnID uint64) error {
    txn, exists := tm.transactions[txnID]
    if !exists {
        return fmt.Errorf("transaction not found")
    }
    
    if txn.Status != "ACTIVE" {
        return fmt.Errorf("transaction is not active")
    }
    
    // 释放锁
    for resource := range txn.Locks {
        delete(tm.locks, resource)
    }
    
    txn.Status = "COMMITTED"
    return nil
}

func (tm *TransactionManager) AcquireLock(txnID uint64, resource, lockType string) error {
    txn, exists := tm.transactions[txnID]
    if !exists {
        return fmt.Errorf("transaction not found")
    }
    
    if existingLock, exists := tm.locks[resource]; exists {
        if existingLock.TxnID == txnID {
            return nil
        }
        
        if lockType == "EXCLUSIVE" || existingLock.Type == "EXCLUSIVE" {
            return fmt.Errorf("lock conflict")
        }
    }
    
    lock := &Lock{
        Resource: resource,
        Type:     lockType,
        TxnID:    txnID,
    }
    
    tm.locks[resource] = lock
    txn.Locks[resource] = lock
    
    return nil
}
```

### 索引管理器 (Index Manager)

```go
type BTreeNode struct {
    IsLeaf   bool
    Keys     []interface{}
    Values   []interface{}
    Children []*BTreeNode
}

type BTreeIndex struct {
    Root  *BTreeNode
    Order int
}

type IndexManager struct {
    indexes map[string]*BTreeIndex
}

func (im *IndexManager) CreateIndex(tableName, columnName string) error {
    indexName := fmt.Sprintf("%s_%s_idx", tableName, columnName)
    
    index := &BTreeIndex{
        Root:  &BTreeNode{IsLeaf: true, Keys: make([]interface{}, 0), Values: make([]interface{}, 0)},
        Order: 4,
    }
    
    im.indexes[indexName] = index
    return nil
}

func (im *IndexManager) Insert(indexName string, key interface{}, rowID string) error {
    index, exists := im.indexes[indexName]
    if !exists {
        return fmt.Errorf("index not found")
    }
    
    return im.insertIntoNode(index.Root, key, rowID)
}

func (im *IndexManager) insertIntoNode(node *BTreeNode, key interface{}, rowID string) error {
    if node.IsLeaf {
        insertIndex := im.findInsertIndex(node.Keys, key)
        
        node.Keys = append(node.Keys, nil)
        copy(node.Keys[insertIndex+1:], node.Keys[insertIndex:])
        node.Keys[insertIndex] = key
        
        node.Values = append(node.Values, nil)
        copy(node.Values[insertIndex+1:], node.Values[insertIndex:])
        node.Values[insertIndex] = rowID
        
        if len(node.Keys) > im.indexes[""].Order {
            im.splitLeafNode(node)
        }
    } else {
        childIndex := im.findChildIndex(node.Keys, key)
        return im.insertIntoNode(node.Children[childIndex], key, rowID)
    }
    
    return nil
}

func (im *IndexManager) Search(indexName string, key interface{}) ([]string, error) {
    index, exists := im.indexes[indexName]
    if !exists {
        return nil, fmt.Errorf("index not found")
    }
    
    return im.searchInNode(index.Root, key)
}

func (im *IndexManager) searchInNode(node *BTreeNode, key interface{}) ([]string, error) {
    if node.IsLeaf {
        for i, k := range node.Keys {
            if im.compareKeys(key, k) == 0 {
                return []string{node.Values[i].(string)}, nil
            }
        }
        return []string{}, nil
    } else {
        childIndex := im.findChildIndex(node.Keys, key)
        return im.searchInNode(node.Children[childIndex], key)
    }
}
```

### 数据库引擎 (Database Engine)

```go
type Table struct {
    Name    string
    Columns []*Column
    Data    []map[string]interface{}
}

type Column struct {
    Name     string
    Type     string
    Nullable bool
    Primary  bool
}

type DatabaseEngine struct {
    storageEngine      *StorageEngine
    queryProcessor     *QueryProcessor
    transactionManager *TransactionManager
    indexManager       *IndexManager
}

func NewDatabaseEngine() *DatabaseEngine {
    storageEngine := &StorageEngine{tables: make(map[string]*Table)}
    queryProcessor := NewQueryProcessor(storageEngine)
    transactionManager := NewTransactionManager()
    indexManager := NewIndexManager()
    
    return &DatabaseEngine{
        storageEngine:      storageEngine,
        queryProcessor:     queryProcessor,
        transactionManager: transactionManager,
        indexManager:       indexManager,
    }
}

func (de *DatabaseEngine) CreateTable(name string, columns []*Column) error {
    table := &Table{
        Name:    name,
        Columns: columns,
        Data:    make([]map[string]interface{}, 0),
    }
    
    de.storageEngine.tables[name] = table
    
    // 为主键创建索引
    for _, column := range columns {
        if column.Primary {
            de.indexManager.CreateIndex(name, column.Name)
            break
        }
    }
    
    return nil
}

func (de *DatabaseEngine) Insert(tableName string, data map[string]interface{}) error {
    table, exists := de.storageEngine.tables[tableName]
    if !exists {
        return fmt.Errorf("table not found")
    }
    
    if err := de.validateData(table, data); err != nil {
        return err
    }
    
    rowID := fmt.Sprintf("%s_%d", tableName, len(table.Data))
    table.Data = append(table.Data, data)
    
    // 更新索引
    for _, column := range table.Columns {
        if column.Primary {
            indexName := fmt.Sprintf("%s_%s_idx", tableName, column.Name)
            de.indexManager.Insert(indexName, data[column.Name], rowID)
            break
        }
    }
    
    return nil
}

func (de *DatabaseEngine) Select(sql string) (*ResultSet, error) {
    plan, err := de.queryProcessor.ParseQuery(sql)
    if err != nil {
        return nil, err
    }
    
    return de.queryProcessor.ExecuteQuery(plan)
}

func (de *DatabaseEngine) Update(tableName string, where map[string]interface{}, data map[string]interface{}) error {
    table, exists := de.storageEngine.tables[tableName]
    if !exists {
        return fmt.Errorf("table not found")
    }
    
    for i, row := range table.Data {
        if de.matchesCondition(row, where) {
            for key, value := range data {
                row[key] = value
            }
            
            // 更新索引
            for key, value := range data {
                if de.isIndexedColumn(table, key) {
                    indexName := fmt.Sprintf("%s_%s_idx", tableName, key)
                    rowID := fmt.Sprintf("%s_%d", tableName, i)
                    de.indexManager.Insert(indexName, value, rowID)
                }
            }
        }
    }
    
    return nil
}
```

## 设计原则

### 1. ACID特性设计

- **原子性**：事务的所有操作要么全部成功，要么全部失败
- **一致性**：事务执行前后数据库保持一致状态
- **隔离性**：并发事务之间相互隔离
- **持久性**：事务提交后数据永久保存

### 2. 性能优化设计

- **索引优化**：合理使用索引提高查询性能
- **缓冲池管理**：减少磁盘I/O操作
- **查询优化**：优化查询执行计划
- **并发控制**：提高并发访问性能

### 3. 可扩展性设计

- **分片策略**：水平分片支持大规模数据
- **读写分离**：主从复制提高读取性能
- **集群部署**：分布式数据库架构
- **缓存策略**：多级缓存提高访问速度

### 4. 高可用性设计

- **备份恢复**：定期备份和快速恢复
- **故障转移**：自动故障检测和转移
- **数据复制**：多副本保证数据安全
- **监控告警**：实时监控和异常告警

## 实现示例

```go
func main() {
    // 创建数据库引擎
    db := NewDatabaseEngine()
    
    // 创建表
    columns := []*Column{
        {Name: "id", Type: "int", Nullable: false, Primary: true},
        {Name: "name", Type: "string", Nullable: false, Primary: false},
        {Name: "age", Type: "int", Nullable: true, Primary: false},
    }
    
    db.CreateTable("users", columns)
    
    // 插入数据
    data1 := map[string]interface{}{
        "id":   1,
        "name": "Alice",
        "age":  25,
    }
    
    data2 := map[string]interface{}{
        "id":   2,
        "name": "Bob",
        "age":  30,
    }
    
    db.Insert("users", data1)
    db.Insert("users", data2)
    
    // 查询数据
    result, err := db.Select("SELECT * FROM users WHERE age > 25")
    if err != nil {
        log.Printf("Failed to query: %v", err)
    } else {
        fmt.Printf("Query result: %+v\n", result)
    }
    
    // 更新数据
    where := map[string]interface{}{"id": 1}
    updateData := map[string]interface{}{"age": 26}
    db.Update("users", where, updateData)
}
```

## 总结

数据库架构通过存储引擎、查询处理器、事务管理器和索引管理器等核心组件，实现了高效、可靠的数据存储和访问。

### 关键要点

1. **存储管理**：高效的缓冲池和页面管理
2. **查询处理**：智能的查询解析和优化
3. **事务管理**：ACID特性的保证
4. **索引管理**：快速的数据检索
5. **并发控制**：多用户安全访问

### 发展趋势

- **分布式数据库**：支持大规模分布式部署
- **内存数据库**：高性能内存存储
- **NoSQL数据库**：非关系型数据存储
- **时序数据库**：专门的时间序列数据存储
- **图数据库**：复杂关系数据存储
