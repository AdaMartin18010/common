# 04-数据存储 (Data Storage)

## 概述

数据存储是大数据系统的基础设施，包括分布式存储、数据仓库、数据湖等技术。本章将介绍数据存储的理论基础、架构设计和Go语言实现。

## 目录

1. [理论基础](#1-理论基础)
2. [分布式存储](#2-分布式存储)
3. [数据仓库](#3-数据仓库)
4. [数据湖](#4-数据湖)
5. [缓存系统](#5-缓存系统)
6. [数据压缩](#6-数据压缩)
7. [数据备份](#7-数据备份)
8. [存储优化](#8-存储优化)

## 1. 理论基础

### 1.1 存储模型

**定义 1.1** (存储函数)
存储函数 $S: D \times K \rightarrow V$ 将数据 $D$ 以键 $K$ 存储到值 $V$：

$$S(d, k) = v \text{ where } v = \text{encode}(d)$$

**定理 1.1** (存储一致性)
对于存储函数 $S$，如果 $d_1 = d_2$，则：

$$S(d_1, k) = S(d_2, k)$$

### 1.2 存储接口

```go
// 存储接口
type Storage interface {
    Put(key string, value []byte) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    List(prefix string) ([]string, error)
    Close() error
}

// 存储配置
type StorageConfig struct {
    Type        string            `json:"type"`
    Endpoint    string            `json:"endpoint"`
    Credentials map[string]string `json:"credentials"`
    Options     map[string]interface{} `json:"options"`
}

// 存储管理器
type StorageManager struct {
    storages map[string]Storage
    config   StorageConfig
}

func (sm *StorageManager) GetStorage(name string) (Storage, error) {
    if storage, exists := sm.storages[name]; exists {
        return storage, nil
    }
    return nil, errors.New("storage not found")
}

func (sm *StorageManager) RegisterStorage(name string, storage Storage) {
    sm.storages[name] = storage
}
```

## 2. 分布式存储

### 2.1 分布式存储架构

**定义 2.1** (分布式存储)
分布式存储将数据分散存储在多个节点上：

$$\text{DistributedStorage}(D) = \bigcup_{i=1}^{n} S_i(D_i)$$

其中 $S_i$ 是第 $i$ 个存储节点，$D_i$ 是分配给该节点的数据。

```go
// 分布式存储节点
type StorageNode struct {
    ID       string            `json:"id"`
    Address  string            `json:"address"`
    Capacity int64             `json:"capacity"`
    Used     int64             `json:"used"`
    Status   NodeStatus        `json:"status"`
    Data     map[string][]byte `json:"-"`
    mu       sync.RWMutex      `json:"-"`
}

type NodeStatus int

const (
    NodeOnline NodeStatus = iota
    NodeOffline
    NodeMaintenance
)

// 分布式存储集群
type DistributedStorage struct {
    nodes    []*StorageNode
    hashRing *ConsistentHashRing
    replicas int
}

// 一致性哈希环
type ConsistentHashRing struct {
    nodes    map[uint32]string
    sorted   []uint32
    replicas int
}

func (chr *ConsistentHashRing) AddNode(nodeID string) {
    for i := 0; i < chr.replicas; i++ {
        hash := chr.hash(fmt.Sprintf("%s:%d", nodeID, i))
        chr.nodes[hash] = nodeID
        chr.sorted = append(chr.sorted, hash)
    }
    sort.Slice(chr.sorted, func(i, j int) bool {
        return chr.sorted[i] < chr.sorted[j]
    })
}

func (chr *ConsistentHashRing) GetNode(key string) string {
    if len(chr.nodes) == 0 {
        return ""
    }
    
    hash := chr.hash(key)
    
    // 查找第一个大于等于hash的节点
    idx := sort.Search(len(chr.sorted), func(i int) bool {
        return chr.sorted[i] >= hash
    })
    
    if idx == len(chr.sorted) {
        idx = 0
    }
    
    return chr.nodes[chr.sorted[idx]]
}

func (chr *ConsistentHashRing) hash(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}

// 分布式存储操作
func (ds *DistributedStorage) Put(key string, value []byte) error {
    // 确定主节点
    primaryNode := ds.hashRing.GetNode(key)
    if primaryNode == "" {
        return errors.New("no available nodes")
    }
    
    // 写入主节点
    if err := ds.writeToNode(primaryNode, key, value); err != nil {
        return err
    }
    
    // 写入副本节点
    for i := 1; i < ds.replicas; i++ {
        replicaKey := fmt.Sprintf("%s:replica:%d", key, i)
        replicaNode := ds.hashRing.GetNode(replicaKey)
        if replicaNode != "" && replicaNode != primaryNode {
            ds.writeToNode(replicaNode, key, value)
        }
    }
    
    return nil
}

func (ds *DistributedStorage) Get(key string) ([]byte, error) {
    // 尝试从主节点读取
    primaryNode := ds.hashRing.GetNode(key)
    if primaryNode == "" {
        return nil, errors.New("no available nodes")
    }
    
    value, err := ds.readFromNode(primaryNode, key)
    if err == nil {
        return value, nil
    }
    
    // 主节点失败，尝试副本节点
    for i := 1; i < ds.replicas; i++ {
        replicaKey := fmt.Sprintf("%s:replica:%d", key, i)
        replicaNode := ds.hashRing.GetNode(replicaKey)
        if replicaNode != "" && replicaNode != primaryNode {
            if value, err := ds.readFromNode(replicaNode, key); err == nil {
                return value, nil
            }
        }
    }
    
    return nil, errors.New("data not found")
}

func (ds *DistributedStorage) writeToNode(nodeID, key string, value []byte) error {
    for _, node := range ds.nodes {
        if node.ID == nodeID && node.Status == NodeOnline {
            node.mu.Lock()
            node.Data[key] = value
            node.Used += int64(len(value))
            node.mu.Unlock()
            return nil
        }
    }
    return errors.New("node not found or offline")
}

func (ds *DistributedStorage) readFromNode(nodeID, key string) ([]byte, error) {
    for _, node := range ds.nodes {
        if node.ID == nodeID && node.Status == NodeOnline {
            node.mu.RLock()
            value, exists := node.Data[key]
            node.mu.RUnlock()
            if exists {
                return value, nil
            }
        }
    }
    return nil, errors.New("data not found")
}
```

### 2.2 数据分片

```go
// 数据分片器
type DataSharder struct {
    shardCount int
    hashFunc   func(string) int
}

func (ds *DataSharder) GetShard(key string) int {
    return ds.hashFunc(key) % ds.shardCount
}

// 范围分片
type RangeSharder struct {
    ranges []Range
}

type Range struct {
    Start string
    End   string
    Shard int
}

func (rs *RangeSharder) GetShard(key string) int {
    for _, r := range rs.ranges {
        if key >= r.Start && key <= r.End {
            return r.Shard
        }
    }
    return 0 // 默认分片
}

// 列表分片
type ListSharder struct {
    shardMap map[string]int
}

func (ls *ListSharder) GetShard(key string) int {
    if shard, exists := ls.shardMap[key]; exists {
        return shard
    }
    return 0 // 默认分片
}
```

## 3. 数据仓库

### 3.1 数据仓库架构

**定义 3.1** (数据仓库)
数据仓库是面向主题的、集成的、相对稳定的、反映历史变化的数据集合：

$$\text{DataWarehouse} = \{\text{ODS}, \text{DW}, \text{DM}\}$$

其中 ODS 是操作数据存储，DW 是数据仓库，DM 是数据集市。

```go
// 数据仓库架构
type DataWarehouse struct {
    ods        *OperationalDataStore
    dw         *DataWarehouseCore
    dm         *DataMart
    etl        *ETLProcessor
    scheduler  *JobScheduler
}

// 操作数据存储
type OperationalDataStore struct {
    storage Storage
    tables  map[string]*Table
}

type Table struct {
    Name    string
    Schema  []Column
    Data    [][]interface{}
    Indexes map[string]*Index
}

type Column struct {
    Name     string
    Type     string
    Nullable bool
    Default  interface{}
}

type Index struct {
    Name    string
    Columns []string
    Type    IndexType
}

type IndexType int

const (
    BTreeIndex IndexType = iota
    HashIndex
    BitmapIndex
)

// 数据仓库核心
type DataWarehouseCore struct {
    factTables    map[string]*FactTable
    dimensionTables map[string]*DimensionTable
    storage       Storage
}

type FactTable struct {
    Name       string
    Measures   []Measure
    Dimensions []string
    Data       [][]interface{}
}

type Measure struct {
    Name     string
    Type     MeasureType
    Function AggregationFunction
}

type MeasureType int

const (
    Additive MeasureType = iota
    SemiAdditive
    NonAdditive
)

type AggregationFunction int

const (
    Sum AggregationFunction = iota
    Avg
    Count
    Min
    Max
)

type DimensionTable struct {
    Name       string
    Attributes []Attribute
    Data       [][]interface{}
}

type Attribute struct {
    Name     string
    Type     string
    Hierarchical bool
    Levels    []string
}
```

### 3.2 ETL处理

```go
// ETL处理器
type ETLProcessor struct {
    extractors map[string]Extractor
    transformers map[string]Transformer
    loaders    map[string]Loader
}

type Extractor interface {
    Extract(source string) ([]Record, error)
    GetSchema() []Column
}

type Transformer interface {
    Transform(records []Record) ([]Record, error)
    Validate(records []Record) error
}

type Loader interface {
    Load(target string, records []Record) error
    GetTargetSchema() []Column
}

type Record struct {
    Data   map[string]interface{}
    Schema []Column
}

// 数据库提取器
type DatabaseExtractor struct {
    connection *sql.DB
    query      string
    schema     []Column
}

func (de *DatabaseExtractor) Extract(source string) ([]Record, error) {
    rows, err := de.connection.Query(de.query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    records := make([]Record, 0)
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, err
        }
        
        record := Record{
            Data:   make(map[string]interface{}),
            Schema: de.schema,
        }
        
        for i, column := range columns {
            record.Data[column] = values[i]
        }
        
        records = append(records, record)
    }
    
    return records, nil
}

// 数据转换器
type DataTransformer struct {
    rules []TransformRule
}

type TransformRule struct {
    SourceColumn string
    TargetColumn string
    Function     TransformFunction
    Parameters   map[string]interface{}
}

type TransformFunction int

const (
    Copy TransformFunction = iota
    Convert
    Calculate
    Lookup
    Filter
)

func (dt *DataTransformer) Transform(records []Record) ([]Record, error) {
    transformed := make([]Record, 0)
    
    for _, record := range records {
        newRecord := Record{
            Data:   make(map[string]interface{}),
            Schema: record.Schema,
        }
        
        // 应用转换规则
        for _, rule := range dt.rules {
            value, err := dt.applyRule(rule, record)
            if err != nil {
                return nil, err
            }
            newRecord.Data[rule.TargetColumn] = value
        }
        
        transformed = append(transformed, newRecord)
    }
    
    return transformed, nil
}

func (dt *DataTransformer) applyRule(rule TransformRule, record Record) (interface{}, error) {
    switch rule.Function {
    case Copy:
        return record.Data[rule.SourceColumn], nil
    case Convert:
        return dt.convertValue(record.Data[rule.SourceColumn], rule.Parameters)
    case Calculate:
        return dt.calculateValue(record.Data, rule.Parameters)
    case Lookup:
        return dt.lookupValue(record.Data[rule.SourceColumn], rule.Parameters)
    case Filter:
        if dt.shouldFilter(record.Data, rule.Parameters) {
            return nil, errors.New("record filtered out")
        }
        return record.Data[rule.SourceColumn], nil
    default:
        return nil, errors.New("unknown transform function")
    }
}

// 数据加载器
type DatabaseLoader struct {
    connection *sql.DB
    table      string
    schema     []Column
}

func (dl *DatabaseLoader) Load(target string, records []Record) error {
    if len(records) == 0 {
        return nil
    }
    
    // 构建INSERT语句
    columns := make([]string, 0)
    for _, col := range dl.schema {
        columns = append(columns, col.Name)
    }
    
    placeholders := make([]string, len(columns))
    for i := range placeholders {
        placeholders[i] = "?"
    }
    
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        dl.table,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "))
    
    // 执行批量插入
    tx, err := dl.connection.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, record := range records {
        values := make([]interface{}, len(columns))
        for i, col := range columns {
            values[i] = record.Data[col]
        }
        
        if _, err := stmt.Exec(values...); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

## 4. 数据湖

### 4.1 数据湖架构

**定义 4.2** (数据湖)
数据湖是存储原始数据的集中式存储库，支持多种数据格式：

$$\text{DataLake} = \{\text{RawData}, \text{ProcessedData}, \text{CuratedData}\}$$

```go
// 数据湖
type DataLake struct {
    rawZone      *RawZone
    processedZone *ProcessedZone
    curatedZone  *CuratedZone
    catalog      *DataCatalog
    governance   *DataGovernance
}

// 原始数据区
type RawZone struct {
    storage Storage
    format  string
    schema  *Schema
}

// 处理数据区
type ProcessedZone struct {
    storage Storage
    format  string
    schema  *Schema
    quality *DataQuality
}

// 治理数据区
type CuratedZone struct {
    storage Storage
    format  string
    schema  *Schema
    lineage *DataLineage
}

// 数据目录
type DataCatalog struct {
    tables    map[string]*TableMetadata
    schemas   map[string]*Schema
    lineage   map[string]*LineageInfo
    policies  map[string]*AccessPolicy
}

type TableMetadata struct {
    Name        string
    Location    string
    Format      string
    Schema      *Schema
    PartitionBy []string
    Created     time.Time
    Updated     time.Time
    Size        int64
    RowCount    int64
}

type Schema struct {
    Fields []Field
    Type   string
}

type Field struct {
    Name     string
    Type     string
    Nullable bool
    Metadata map[string]string
}

// 数据治理
type DataGovernance struct {
    policies  map[string]*Policy
    rules     map[string]*Rule
    audits    []*AuditLog
}

type Policy struct {
    Name        string
    Type        PolicyType
    Conditions  []Condition
    Actions     []Action
    Priority    int
}

type PolicyType int

const (
    AccessPolicy PolicyType = iota
    RetentionPolicy
    QualityPolicy
    SecurityPolicy
)

type Condition struct {
    Field    string
    Operator string
    Value    interface{}
}

type Action struct {
    Type   string
    Params map[string]interface{}
}
```

### 4.2 数据湖操作

```go
// 数据湖操作器
type DataLakeOperator struct {
    lake *DataLake
}

func (dlo *DataLakeOperator) Ingest(data []byte, metadata *IngestMetadata) error {
    // 写入原始数据区
    rawPath := fmt.Sprintf("raw/%s/%s", metadata.Source, metadata.Timestamp.Format("2006/01/02"))
    if err := dlo.lake.rawZone.storage.Put(rawPath, data); err != nil {
        return err
    }
    
    // 更新目录
    dlo.lake.catalog.tables[metadata.TableName] = &TableMetadata{
        Name:     metadata.TableName,
        Location: rawPath,
        Format:   metadata.Format,
        Schema:   metadata.Schema,
        Created:  time.Now(),
        Size:     int64(len(data)),
    }
    
    return nil
}

func (dlo *DataLakeOperator) Process(tableName string, processor DataProcessor) error {
    // 读取原始数据
    metadata := dlo.lake.catalog.tables[tableName]
    if metadata == nil {
        return errors.New("table not found")
    }
    
    data, err := dlo.lake.rawZone.storage.Get(metadata.Location)
    if err != nil {
        return err
    }
    
    // 处理数据
    processedData, err := processor.Process(data)
    if err != nil {
        return err
    }
    
    // 写入处理数据区
    processedPath := fmt.Sprintf("processed/%s/%s", tableName, time.Now().Format("2006/01/02"))
    if err := dlo.lake.processedZone.storage.Put(processedPath, processedData); err != nil {
        return err
    }
    
    // 更新目录
    dlo.lake.catalog.tables[tableName].Location = processedPath
    dlo.lake.catalog.tables[tableName].Updated = time.Now()
    
    return nil
}

func (dlo *DataLakeOperator) Query(query string) ([]Record, error) {
    // 解析查询
    parsedQuery, err := dlo.parseQuery(query)
    if err != nil {
        return nil, err
    }
    
    // 确定数据源
    source := dlo.determineSource(parsedQuery)
    
    // 执行查询
    return dlo.executeQuery(source, parsedQuery)
}

func (dlo *DataLakeOperator) parseQuery(query string) (*Query, error) {
    // 简化的查询解析
    return &Query{
        SQL: query,
    }, nil
}

func (dlo *DataLakeOperator) determineSource(query *Query) string {
    // 根据查询确定最佳数据源
    // 优先使用治理数据区，然后是处理数据区，最后是原始数据区
    return "curated"
}

func (dlo *DataLakeOperator) executeQuery(source string, query *Query) ([]Record, error) {
    // 执行查询（简化实现）
    return []Record{}, nil
}

type IngestMetadata struct {
    TableName string
    Source    string
    Format    string
    Schema    *Schema
    Timestamp time.Time
}

type Query struct {
    SQL string
}
```

## 5. 缓存系统

### 5.1 缓存架构

```go
// 缓存接口
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
    Size() int
}

// 内存缓存
type MemoryCache struct {
    data    map[string]*CacheEntry
    maxSize int
    mu      sync.RWMutex
}

type CacheEntry struct {
    Value      interface{}
    Expiration time.Time
    AccessTime time.Time
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
    mc.mu.RLock()
    defer mc.mu.RUnlock()
    
    entry, exists := mc.data[key]
    if !exists {
        return nil, false
    }
    
    // 检查过期
    if time.Now().After(entry.Expiration) {
        delete(mc.data, key)
        return nil, false
    }
    
    // 更新访问时间
    entry.AccessTime = time.Now()
    return entry.Value, true
}

func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    // 检查容量
    if len(mc.data) >= mc.maxSize {
        mc.evictLRU()
    }
    
    mc.data[key] = &CacheEntry{
        Value:      value,
        Expiration: time.Now().Add(ttl),
        AccessTime: time.Now(),
    }
    
    return nil
}

func (mc *MemoryCache) evictLRU() {
    var oldestKey string
    var oldestTime time.Time
    
    for key, entry := range mc.data {
        if oldestKey == "" || entry.AccessTime.Before(oldestTime) {
            oldestKey = key
            oldestTime = entry.AccessTime
        }
    }
    
    if oldestKey != "" {
        delete(mc.data, oldestKey)
    }
}

// 分布式缓存
type DistributedCache struct {
    nodes    []*CacheNode
    hashRing *ConsistentHashRing
}

type CacheNode struct {
    ID      string
    Address string
    Cache   Cache
}

func (dc *DistributedCache) Get(key string) (interface{}, bool) {
    nodeID := dc.hashRing.GetNode(key)
    if nodeID == "" {
        return nil, false
    }
    
    for _, node := range dc.nodes {
        if node.ID == nodeID {
            return node.Cache.Get(key)
        }
    }
    
    return nil, false
}

func (dc *DistributedCache) Set(key string, value interface{}, ttl time.Duration) error {
    nodeID := dc.hashRing.GetNode(key)
    if nodeID == "" {
        return errors.New("no available nodes")
    }
    
    for _, node := range dc.nodes {
        if node.ID == nodeID {
            return node.Cache.Set(key, value, ttl)
        }
    }
    
    return errors.New("node not found")
}
```

## 6. 数据压缩

### 6.1 压缩算法

```go
// 压缩器接口
type Compressor interface {
    Compress(data []byte) ([]byte, error)
    Decompress(data []byte) ([]byte, error)
    GetCompressionRatio() float64
}

// GZIP压缩器
type GZIPCompressor struct {
    level int
}

func (gc *GZIPCompressor) Compress(data []byte) ([]byte, error) {
    var buffer bytes.Buffer
    writer, err := gzip.NewWriterLevel(&buffer, gc.level)
    if err != nil {
        return nil, err
    }
    
    if _, err := writer.Write(data); err != nil {
        return nil, err
    }
    
    if err := writer.Close(); err != nil {
        return nil, err
    }
    
    return buffer.Bytes(), nil
}

func (gc *GZIPCompressor) Decompress(data []byte) ([]byte, error) {
    reader, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    var buffer bytes.Buffer
    if _, err := buffer.ReadFrom(reader); err != nil {
        return nil, err
    }
    
    return buffer.Bytes(), nil
}

// LZ4压缩器
type LZ4Compressor struct{}

func (lc *LZ4Compressor) Compress(data []byte) ([]byte, error) {
    compressed := make([]byte, lz4.CompressBound(len(data)))
    n, err := lz4.CompressBlock(data, compressed, nil)
    if err != nil {
        return nil, err
    }
    
    return compressed[:n], nil
}

func (lc *LZ4Compressor) Decompress(data []byte) ([]byte, error) {
    // 需要知道原始大小，这里简化处理
    decompressed := make([]byte, len(data)*4) // 假设压缩比至少为4:1
    n, err := lz4.DecompressBlock(data, decompressed)
    if err != nil {
        return nil, err
    }
    
    return decompressed[:n], nil
}

// 压缩存储
type CompressedStorage struct {
    storage   Storage
    compressor Compressor
}

func (cs *CompressedStorage) Put(key string, value []byte) error {
    compressed, err := cs.compressor.Compress(value)
    if err != nil {
        return err
    }
    
    return cs.storage.Put(key, compressed)
}

func (cs *CompressedStorage) Get(key string) ([]byte, error) {
    compressed, err := cs.storage.Get(key)
    if err != nil {
        return nil, err
    }
    
    return cs.compressor.Decompress(compressed)
}
```

## 7. 数据备份

### 7.1 备份策略

```go
// 备份管理器
type BackupManager struct {
    storage    Storage
    strategy   BackupStrategy
    scheduler  *BackupScheduler
}

type BackupStrategy struct {
    Type        BackupType
    Frequency   time.Duration
    Retention   time.Duration
    Compression bool
    Encryption  bool
}

type BackupType int

const (
    FullBackup BackupType = iota
    IncrementalBackup
    DifferentialBackup
)

// 备份调度器
type BackupScheduler struct {
    jobs    map[string]*BackupJob
    ticker  *time.Ticker
    done    chan bool
}

type BackupJob struct {
    ID       string
    Strategy BackupStrategy
    LastRun  time.Time
    NextRun  time.Time
    Status   JobStatus
}

type JobStatus int

const (
    JobPending JobStatus = iota
    JobRunning
    JobCompleted
    JobFailed
)

func (bm *BackupManager) ScheduleBackup(job *BackupJob) error {
    bm.scheduler.jobs[job.ID] = job
    return nil
}

func (bm *BackupManager) ExecuteBackup(jobID string) error {
    job, exists := bm.scheduler.jobs[jobID]
    if !exists {
        return errors.New("job not found")
    }
    
    job.Status = JobRunning
    job.LastRun = time.Now()
    
    // 执行备份
    if err := bm.performBackup(job); err != nil {
        job.Status = JobFailed
        return err
    }
    
    job.Status = JobCompleted
    job.NextRun = time.Now().Add(job.Strategy.Frequency)
    
    return nil
}

func (bm *BackupManager) performBackup(job *BackupJob) error {
    switch job.Strategy.Type {
    case FullBackup:
        return bm.fullBackup(job)
    case IncrementalBackup:
        return bm.incrementalBackup(job)
    case DifferentialBackup:
        return bm.differentialBackup(job)
    default:
        return errors.New("unknown backup type")
    }
}

func (bm *BackupManager) fullBackup(job *BackupJob) error {
    // 完整备份实现
    backupPath := fmt.Sprintf("backup/full/%s_%s", job.ID, time.Now().Format("20060102_150405"))
    
    // 获取所有数据
    keys, err := bm.storage.List("")
    if err != nil {
        return err
    }
    
    backup := make(map[string][]byte)
    for _, key := range keys {
        data, err := bm.storage.Get(key)
        if err != nil {
            continue
        }
        backup[key] = data
    }
    
    // 序列化备份
    backupData, err := json.Marshal(backup)
    if err != nil {
        return err
    }
    
    // 压缩和加密
    if job.Strategy.Compression {
        backupData, err = bm.compress(backupData)
        if err != nil {
            return err
        }
    }
    
    if job.Strategy.Encryption {
        backupData, err = bm.encrypt(backupData)
        if err != nil {
            return err
        }
    }
    
    // 存储备份
    return bm.storage.Put(backupPath, backupData)
}

func (bm *BackupManager) compress(data []byte) ([]byte, error) {
    var buffer bytes.Buffer
    writer := gzip.NewWriter(&buffer)
    
    if _, err := writer.Write(data); err != nil {
        return nil, err
    }
    
    if err := writer.Close(); err != nil {
        return nil, err
    }
    
    return buffer.Bytes(), nil
}

func (bm *BackupManager) encrypt(data []byte) ([]byte, error) {
    // 简化的加密实现
    return data, nil
}
```

## 8. 存储优化

### 8.1 存储优化策略

```go
// 存储优化器
type StorageOptimizer struct {
    storage Storage
    metrics *StorageMetrics
    policies []OptimizationPolicy
}

type StorageMetrics struct {
    ReadLatency   time.Duration
    WriteLatency  time.Duration
    Throughput    float64
    Utilization   float64
    ErrorRate     float64
}

type OptimizationPolicy struct {
    Name        string
    Type        OptimizationType
    Threshold   float64
    Action      OptimizationAction
}

type OptimizationType int

const (
    LatencyOptimization OptimizationType = iota
    ThroughputOptimization
    SpaceOptimization
    CostOptimization
)

type OptimizationAction struct {
    Type   string
    Params map[string]interface{}
}

func (so *StorageOptimizer) Optimize() error {
    // 收集指标
    metrics := so.collectMetrics()
    
    // 应用优化策略
    for _, policy := range so.policies {
        if so.shouldApplyPolicy(policy, metrics) {
            if err := so.applyPolicy(policy); err != nil {
                return err
            }
        }
    }
    
    return nil
}

func (so *StorageOptimizer) collectMetrics() *StorageMetrics {
    // 收集存储指标
    return &StorageMetrics{
        ReadLatency:  time.Millisecond * 10,
        WriteLatency: time.Millisecond * 20,
        Throughput:   1000.0,
        Utilization:  0.75,
        ErrorRate:    0.001,
    }
}

func (so *StorageOptimizer) shouldApplyPolicy(policy OptimizationPolicy, metrics *StorageMetrics) bool {
    switch policy.Type {
    case LatencyOptimization:
        return metrics.ReadLatency > time.Duration(policy.Threshold)*time.Millisecond
    case ThroughputOptimization:
        return metrics.Throughput < policy.Threshold
    case SpaceOptimization:
        return metrics.Utilization > policy.Threshold
    case CostOptimization:
        return true // 总是应用成本优化
    default:
        return false
    }
}

func (so *StorageOptimizer) applyPolicy(policy OptimizationPolicy) error {
    switch policy.Action.Type {
    case "compression":
        return so.enableCompression(policy.Action.Params)
    case "partitioning":
        return so.enablePartitioning(policy.Action.Params)
    case "caching":
        return so.enableCaching(policy.Action.Params)
    case "deduplication":
        return so.enableDeduplication(policy.Action.Params)
    default:
        return errors.New("unknown optimization action")
    }
}

func (so *StorageOptimizer) enableCompression(params map[string]interface{}) error {
    // 启用压缩
    return nil
}

func (so *StorageOptimizer) enablePartitioning(params map[string]interface{}) error {
    // 启用分区
    return nil
}

func (so *StorageOptimizer) enableCaching(params map[string]interface{}) error {
    // 启用缓存
    return nil
}

func (so *StorageOptimizer) enableDeduplication(params map[string]interface{}) error {
    // 启用去重
    return nil
}
```

## 总结

本章详细介绍了数据存储的核心技术，包括：

1. **理论基础**：存储模型、存储接口
2. **分布式存储**：分布式存储架构、数据分片
3. **数据仓库**：数据仓库架构、ETL处理
4. **数据湖**：数据湖架构、数据湖操作
5. **缓存系统**：缓存架构、分布式缓存
6. **数据压缩**：压缩算法、压缩存储
7. **数据备份**：备份策略、备份管理
8. **存储优化**：存储优化策略、性能优化

这些技术为构建高效、可靠的数据存储系统提供了完整的理论基础和实现方案。

---

**相关链接**：

- [01-数据摄入](../01-Data-Ingestion.md)
- [02-数据处理](../02-Data-Processing.md)
- [03-数据分析](../03-Data-Analytics.md)
- [05-数据可视化](../05-Data-Visualization.md)
- [06-机器学习](../06-Machine-Learning.md)
- [07-实时计算](../07-Real-Time-Computing.md)
- [08-数据治理](../08-Data-Governance.md)
