# 02-数据处理 (Data Processing)

## 概述

数据处理是大数据系统的核心环节，包括数据清洗、转换、聚合、分析等操作。本章将介绍大数据处理的理论基础、技术架构和Go语言实现。

## 目录

1. [理论基础](#1-理论基础)
2. [数据清洗](#2-数据清洗)
3. [数据转换](#3-数据转换)
4. [数据聚合](#4-数据聚合)
5. [流式处理](#5-流式处理)
6. [批处理](#6-批处理)
7. [实时处理](#7-实时处理)
8. [性能优化](#8-性能优化)

## 1. 理论基础

### 1.1 数据处理模型

数据处理遵循以下数学模型：

**定义 1.1** (数据处理函数)
设 ```latex
D
``` 为数据集，```latex
F
``` 为处理函数集合，数据处理函数 ```latex
f: D \rightarrow D'
``` 满足：

$
f(d) = \begin{cases}
\text{transform}(d) & \text{if } \text{valid}(d) \\
\text{clean}(d) & \text{if } \text{dirty}(d) \\
\text{drop}(d) & \text{if } \text{invalid}(d)
\end{cases}
$

**定理 1.1** (数据处理幂等性)
对于数据处理函数 ```latex
f
```，如果 ```latex
f
``` 是幂等的，则：

$```latex
f(f(d)) = f(d)
```$

### 1.2 处理模式

```go
// 数据处理接口定义
type DataProcessor interface {
    Process(data []byte) ([]byte, error)
    Validate(data []byte) bool
    Transform(data []byte) ([]byte, error)
    Clean(data []byte) ([]byte, error)
}

// 基础处理器
type BaseProcessor struct {
    validators []Validator
    transformers []Transformer
    cleaners []Cleaner
}

func (bp *BaseProcessor) Process(data []byte) ([]byte, error) {
    // 验证数据
    if !bp.Validate(data) {
        return nil, errors.New("invalid data")
    }

    // 清洗数据
    cleaned, err := bp.Clean(data)
    if err != nil {
        return nil, err
    }

    // 转换数据
    transformed, err := bp.Transform(cleaned)
    if err != nil {
        return nil, err
    }

    return transformed, nil
}
```

## 2. 数据清洗

### 2.1 清洗策略

数据清洗包括以下操作：

1. **去重**：移除重复数据
2. **补全**：填充缺失值
3. **修正**：修正错误数据
4. **标准化**：统一数据格式

```go
// 数据清洗器
type DataCleaner struct {
    deduplicator Deduplicator
    imputer      Imputer
    validator    Validator
    normalizer   Normalizer
}

// 去重器
type Deduplicator interface {
    Deduplicate(data [][]byte) [][]byte
}

// 哈希去重实现
type HashDeduplicator struct {
    seen map[string]bool
}

func (hd *HashDeduplicator) Deduplicate(data [][]byte) [][]byte {
    result := make([][]byte, 0)
    for _, item := range data {
        hash := sha256.Sum256(item)
        hashStr := hex.EncodeToString(hash[:])

        if !hd.seen[hashStr] {
            hd.seen[hashStr] = true
            result = append(result, item)
        }
    }
    return result
}

// 缺失值填充器
type Imputer interface {
    Impute(data []byte, strategy ImputeStrategy) ([]byte, error)
}

type ImputeStrategy int

const (
    MeanImpute ImputeStrategy = iota
    MedianImpute
    ModeImpute
    ForwardFill
    BackwardFill
)
```

### 2.2 清洗算法

**算法 2.1** (数据清洗算法)

```go
func CleanData(data []Record) []Record {
    cleaned := make([]Record, 0)

    for _, record := range data {
        // 1. 验证记录
        if !ValidateRecord(record) {
            continue
        }

        // 2. 清洗记录
        cleanedRecord := CleanRecord(record)

        // 3. 标准化记录
        normalizedRecord := NormalizeRecord(cleanedRecord)

        cleaned = append(cleaned, normalizedRecord)
    }

    return cleaned
}

func ValidateRecord(record Record) bool {
    // 检查必填字段
    if record.ID == "" || record.Timestamp.IsZero() {
        return false
    }

    // 检查数据类型
    if !IsValidDataType(record.Data) {
        return false
    }

    return true
}

func CleanRecord(record Record) Record {
    // 去除空白字符
    record.Data = strings.TrimSpace(record.Data)

    // 修正格式
    record.Data = CorrectFormat(record.Data)

    // 填充缺失值
    if record.Data == "" {
        record.Data = DefaultValue
    }

    return record
}
```

## 3. 数据转换

### 3.1 转换类型

数据转换包括：

1. **格式转换**：JSON ↔ XML ↔ CSV
2. **编码转换**：UTF-8 ↔ ASCII ↔ Base64
3. **类型转换**：字符串 ↔ 数字 ↔ 布尔值
4. **结构转换**：扁平化 ↔ 嵌套化

```go
// 数据转换器
type DataTransformer struct {
    formatConverters map[string]FormatConverter
    encodingConverters map[string]EncodingConverter
    typeConverters map[string]TypeConverter
}

// 格式转换器
type FormatConverter interface {
    Convert(data []byte, fromFormat, toFormat string) ([]byte, error)
}

// JSON转换器
type JSONConverter struct{}

func (jc *JSONConverter) Convert(data []byte, fromFormat, toFormat string) ([]byte, error) {
    switch {
    case fromFormat == "json" && toFormat == "xml":
        return jc.JSONToXML(data)
    case fromFormat == "json" && toFormat == "csv":
        return jc.JSONToCSV(data)
    case fromFormat == "xml" && toFormat == "json":
        return jc.XMLToJSON(data)
    case fromFormat == "csv" && toFormat == "json":
        return jc.CSVToJSON(data)
    default:
        return nil, errors.New("unsupported format conversion")
    }
}

func (jc *JSONConverter) JSONToXML(data []byte) ([]byte, error) {
    var jsonData interface{}
    if err := json.Unmarshal(data, &jsonData); err != nil {
        return nil, err
    }

    // 转换为XML
    xmlData, err := xml.MarshalIndent(jsonData, "", "  ")
    if err != nil {
        return nil, err
    }

    return xmlData, nil
}
```

### 3.2 转换管道

```go
// 转换管道
type TransformPipeline struct {
    stages []TransformStage
}

type TransformStage interface {
    Transform(data []byte) ([]byte, error)
    Name() string
}

// 管道执行
func (tp *TransformPipeline) Execute(data []byte) ([]byte, error) {
    result := data

    for _, stage := range tp.stages {
        transformed, err := stage.Transform(result)
        if err != nil {
            return nil, fmt.Errorf("stage %s failed: %w", stage.Name(), err)
        }
        result = transformed
    }

    return result, nil
}

// 具体转换阶段
type JSONFlattenStage struct{}

func (jfs *JSONFlattenStage) Transform(data []byte) ([]byte, error) {
    var jsonData map[string]interface{}
    if err := json.Unmarshal(data, &jsonData); err != nil {
        return nil, err
    }

    flattened := FlattenJSON(jsonData)
    return json.Marshal(flattened)
}

func (jfs *JSONFlattenStage) Name() string {
    return "JSONFlatten"
}

func FlattenJSON(data map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})

    for key, value := range data {
        switch v := value.(type) {
        case map[string]interface{}:
            flattened := FlattenJSON(v)
            for fk, fv := range flattened {
                result[key+"."+fk] = fv
            }
        default:
            result[key] = v
        }
    }

    return result
}
```

## 4. 数据聚合

### 4.1 聚合函数

**定义 4.1** (聚合函数)
聚合函数 ```latex
f: D^n \rightarrow D
``` 满足：

$```latex
f(d_1, d_2, ..., d_n) = \text{aggregate}(d_1, d_2, ..., d_n)
```$

常见的聚合函数包括：

1. **统计聚合**：COUNT, SUM, AVG, MIN, MAX
2. **分组聚合**：GROUP BY
3. **窗口聚合**：滑动窗口聚合
4. **自定义聚合**：用户定义的聚合逻辑

```go
// 聚合器接口
type Aggregator interface {
    Aggregate(data []interface{}) (interface{}, error)
    Type() string
}

// 统计聚合器
type StatisticalAggregator struct {
    aggType string
}

func (sa *StatisticalAggregator) Aggregate(data []interface{}) (interface{}, error) {
    switch sa.aggType {
    case "count":
        return len(data), nil
    case "sum":
        return sa.sum(data)
    case "avg":
        return sa.average(data)
    case "min":
        return sa.minimum(data)
    case "max":
        return sa.maximum(data)
    default:
        return nil, errors.New("unknown aggregation type")
    }
}

func (sa *StatisticalAggregator) sum(data []interface{}) (float64, error) {
    sum := 0.0
    for _, item := range data {
        if num, ok := item.(float64); ok {
            sum += num
        } else {
            return 0, errors.New("non-numeric data")
        }
    }
    return sum, nil
}

func (sa *StatisticalAggregator) average(data []interface{}) (float64, error) {
    sum, err := sa.sum(data)
    if err != nil {
        return 0, err
    }
    return sum / float64(len(data)), nil
}

func (sa *StatisticalAggregator) Type() string {
    return sa.aggType
}
```

### 4.2 分组聚合

```go
// 分组聚合器
type GroupAggregator struct {
    keyFunc   func(interface{}) string
    aggregator Aggregator
}

func (ga *GroupAggregator) GroupAggregate(data []interface{}) (map[string]interface{}, error) {
    groups := make(map[string][]interface{})

    // 分组
    for _, item := range data {
        key := ga.keyFunc(item)
        groups[key] = append(groups[key], item)
    }

    // 聚合
    result := make(map[string]interface{})
    for key, groupData := range groups {
        aggregated, err := ga.aggregator.Aggregate(groupData)
        if err != nil {
            return nil, err
        }
        result[key] = aggregated
    }

    return result, nil
}

// 使用示例
func ExampleGroupAggregation() {
    data := []User{
        {Name: "Alice", Age: 25, Salary: 50000},
        {Name: "Bob", Age: 30, Salary: 60000},
        {Name: "Alice", Age: 28, Salary: 55000},
        {Name: "Bob", Age: 32, Salary: 65000},
    }

    // 按姓名分组，计算平均年龄
    keyFunc := func(item interface{}) string {
        return item.(User).Name
    }

    ageAggregator := &StatisticalAggregator{aggType: "avg"}
    groupAggregator := &GroupAggregator{
        keyFunc: keyFunc,
        aggregator: ageAggregator,
    }

    result, err := groupAggregator.GroupAggregate(interface{}(data))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Average age by name: %+v\n", result)
}
```

## 5. 流式处理

### 5.1 流处理模型

**定义 5.1** (流处理)
流处理是对连续数据流的实时处理，满足：

$```latex
\text{Stream}(t) = \text{process}(\text{data}_t, \text{state}_{t-1})
```$

```go
// 流处理器
type StreamProcessor struct {
    windowSize time.Duration
    processors []StreamStage
    state      map[string]interface{}
}

type StreamStage interface {
    Process(data []byte, state map[string]interface{}) ([]byte, error)
    Name() string
}

// 滑动窗口处理器
type SlidingWindowProcessor struct {
    windowSize time.Duration
    buffer     []DataPoint
}

func (swp *SlidingWindowProcessor) Process(data []byte, state map[string]interface{}) ([]byte, error) {
    var point DataPoint
    if err := json.Unmarshal(data, &point); err != nil {
        return nil, err
    }

    // 添加新数据点
    swp.buffer = append(swp.buffer, point)

    // 移除过期数据
    cutoff := time.Now().Add(-swp.windowSize)
    validBuffer := make([]DataPoint, 0)
    for _, p := range swp.buffer {
        if p.Timestamp.After(cutoff) {
            validBuffer = append(validBuffer, p)
        }
    }
    swp.buffer = validBuffer

    // 计算窗口统计
    stats := swp.calculateStats(swp.buffer)

    return json.Marshal(stats)
}

func (swp *SlidingWindowProcessor) calculateStats(data []DataPoint) WindowStats {
    if len(data) == 0 {
        return WindowStats{}
    }

    sum := 0.0
    min := data[0].Value
    max := data[0].Value

    for _, point := range data {
        sum += point.Value
        if point.Value < min {
            min = point.Value
        }
        if point.Value > max {
            max = point.Value
        }
    }

    return WindowStats{
        Count: len(data),
        Sum:   sum,
        Avg:   sum / float64(len(data)),
        Min:   min,
        Max:   max,
    }
}

func (swp *SlidingWindowProcessor) Name() string {
    return "SlidingWindow"
}
```

### 5.2 流处理管道

```go
// 流处理管道
type StreamPipeline struct {
    stages []StreamStage
    input  chan []byte
    output chan []byte
}

func (sp *StreamPipeline) Start() {
    go func() {
        for data := range sp.input {
            result := data

            for _, stage := range sp.stages {
                processed, err := stage.Process(result, make(map[string]interface{}))
                if err != nil {
                    log.Printf("Stage %s failed: %v", stage.Name(), err)
                    continue
                }
                result = processed
            }

            sp.output <- result
        }
    }()
}

// 使用示例
func ExampleStreamProcessing() {
    pipeline := &StreamPipeline{
        stages: []StreamStage{
            &JSONParser{},
            &DataValidator{},
            &SlidingWindowProcessor{windowSize: 5 * time.Minute},
            &Aggregator{},
        },
        input:  make(chan []byte, 100),
        output: make(chan []byte, 100),
    }

    pipeline.Start()

    // 发送数据
    go func() {
        for i := 0; i < 100; i++ {
            data := fmt.Sprintf(`{"timestamp":"%s","value":%d}`,
                time.Now().Format(time.RFC3339), i)
            pipeline.input <- []byte(data)
            time.Sleep(100 * time.Millisecond)
        }
        close(pipeline.input)
    }()

    // 接收结果
    for result := range pipeline.output {
        fmt.Printf("Processed: %s\n", string(result))
    }
}
```

## 6. 批处理

### 6.1 批处理模型

**定义 6.1** (批处理)
批处理是对大量数据的批量处理，满足：

$```latex
\text{Batch}(D) = \bigcup_{i=1}^{n} \text{process}(D_i)
```$

其中 ```latex
D = \{D_1, D_2, ..., D_n\}
``` 是数据分片。

```go
// 批处理器
type BatchProcessor struct {
    batchSize int
    workers   int
    processor DataProcessor
}

func (bp *BatchProcessor) ProcessBatch(data [][]byte) ([][]byte, error) {
    // 分片数据
    chunks := bp.chunkData(data, bp.batchSize)

    // 并行处理
    results := make([][][]byte, len(chunks))
    var wg sync.WaitGroup
    errChan := make(chan error, len(chunks))

    for i, chunk := range chunks {
        wg.Add(1)
        go func(index int, chunkData [][]byte) {
            defer wg.Done()

            processed, err := bp.processChunk(chunkData)
            if err != nil {
                errChan <- err
                return
            }
            results[index] = processed
        }(i, chunk)
    }

    wg.Wait()
    close(errChan)

    // 检查错误
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }

    // 合并结果
    return bp.mergeResults(results), nil
}

func (bp *BatchProcessor) chunkData(data [][]byte, size int) [][][]byte {
    chunks := make([][][]byte, 0)
    for i := 0; i < len(data); i += size {
        end := i + size
        if end > len(data) {
            end = len(data)
        }
        chunks = append(chunks, data[i:end])
    }
    return chunks
}

func (bp *BatchProcessor) processChunk(chunk [][]byte) ([][]byte, error) {
    results := make([][]byte, 0)

    for _, item := range chunk {
        processed, err := bp.processor.Process(item)
        if err != nil {
            return nil, err
        }
        results = append(results, processed)
    }

    return results, nil
}

func (bp *BatchProcessor) mergeResults(results [][][]byte) [][]byte {
    merged := make([][]byte, 0)
    for _, result := range results {
        merged = append(merged, result...)
    }
    return merged
}
```

### 6.2 MapReduce模式

```go
// MapReduce处理器
type MapReduceProcessor struct {
    mapper   Mapper
    reducer  Reducer
    workers  int
}

type Mapper interface {
    Map(key string, value []byte) ([]KeyValue, error)
}

type Reducer interface {
    Reduce(key string, values [][]byte) ([]byte, error)
}

type KeyValue struct {
    Key   string
    Value []byte
}

func (mrp *MapReduceProcessor) Execute(data map[string][]byte) (map[string][]byte, error) {
    // Map阶段
    mapped := make(map[string][][]byte)

    for key, value := range data {
        keyValues, err := mrp.mapper.Map(key, value)
        if err != nil {
            return nil, err
        }

        for _, kv := range keyValues {
            mapped[kv.Key] = append(mapped[kv.Key], kv.Value)
        }
    }

    // Reduce阶段
    result := make(map[string][]byte)

    for key, values := range mapped {
        reduced, err := mrp.reducer.Reduce(key, values)
        if err != nil {
            return nil, err
        }
        result[key] = reduced
    }

    return result, nil
}

// 单词计数示例
type WordCountMapper struct{}

func (wcm *WordCountMapper) Map(key string, value []byte) ([]KeyValue, error) {
    words := strings.Fields(string(value))
    result := make([]KeyValue, 0)

    for _, word := range words {
        word = strings.ToLower(strings.Trim(word, ".,!?"))
        if word != "" {
            result = append(result, KeyValue{
                Key:   word,
                Value: []byte("1"),
            })
        }
    }

    return result, nil
}

type WordCountReducer struct{}

func (wcr *WordCountReducer) Reduce(key string, values [][]byte) ([]byte, error) {
    count := 0
    for _, value := range values {
        if val, err := strconv.Atoi(string(value)); err == nil {
            count += val
        }
    }

    return []byte(strconv.Itoa(count)), nil
}
```

## 7. 实时处理

### 7.1 实时处理架构

**定义 7.1** (实时处理)
实时处理要求在时间约束 ```latex
T
``` 内完成处理：

$```latex
\text{RealTime}(data_t) = \text{process}(data_t) \text{ s.t. } \text{latency} \leq T
```$

```go
// 实时处理器
type RealTimeProcessor struct {
    maxLatency time.Duration
    processor  DataProcessor
    metrics    *Metrics
}

type Metrics struct {
    processedCount int64
    errorCount     int64
    avgLatency     time.Duration
    maxLatency     time.Duration
    mu             sync.RWMutex
}

func (rtp *RealTimeProcessor) Process(data []byte) ([]byte, error) {
    start := time.Now()

    // 处理数据
    result, err := rtp.processor.Process(data)

    // 记录指标
    latency := time.Since(start)
    rtp.recordMetrics(err, latency)

    // 检查延迟约束
    if latency > rtp.maxLatency {
        return nil, fmt.Errorf("latency %v exceeds limit %v", latency, rtp.maxLatency)
    }

    return result, err
}

func (rtp *RealTimeProcessor) recordMetrics(err error, latency time.Duration) {
    rtp.metrics.mu.Lock()
    defer rtp.metrics.mu.Unlock()

    rtp.metrics.processedCount++
    if err != nil {
        rtp.metrics.errorCount++
    }

    // 更新平均延迟
    totalLatency := rtp.metrics.avgLatency * time.Duration(rtp.metrics.processedCount-1)
    rtp.metrics.avgLatency = (totalLatency + latency) / time.Duration(rtp.metrics.processedCount)

    // 更新最大延迟
    if latency > rtp.metrics.maxLatency {
        rtp.metrics.maxLatency = latency
    }
}
```

### 7.2 实时流处理

```go
// 实时流处理器
type RealTimeStreamProcessor struct {
    input     chan DataEvent
    output    chan ProcessedEvent
    processor RealTimeProcessor
    buffer    *RingBuffer
}

type DataEvent struct {
    ID        string
    Data      []byte
    Timestamp time.Time
}

type ProcessedEvent struct {
    EventID   string
    Result    []byte
    Latency   time.Duration
    Timestamp time.Time
}

type RingBuffer struct {
    buffer []DataEvent
    head   int
    tail   int
    size   int
    mu     sync.RWMutex
}

func (rb *RingBuffer) Push(event DataEvent) bool {
    rb.mu.Lock()
    defer rb.mu.Unlock()

    next := (rb.tail + 1) % rb.size
    if next == rb.head {
        return false // 缓冲区满
    }

    rb.buffer[rb.tail] = event
    rb.tail = next
    return true
}

func (rb *RingBuffer) Pop() (DataEvent, bool) {
    rb.mu.Lock()
    defer rb.mu.Unlock()

    if rb.head == rb.tail {
        return DataEvent{}, false // 缓冲区空
    }

    event := rb.buffer[rb.head]
    rb.head = (rb.head + 1) % rb.size
    return event, true
}

func (rtsp *RealTimeStreamProcessor) Start() {
    go func() {
        for event := range rtsp.input {
            // 尝试处理
            result, err := rtsp.processor.Process(event.Data)

            processed := ProcessedEvent{
                EventID:   event.ID,
                Result:    result,
                Latency:   time.Since(event.Timestamp),
                Timestamp: time.Now(),
            }

            if err != nil {
                // 处理失败，放入缓冲区
                if !rtsp.buffer.Push(event) {
                    log.Printf("Buffer full, dropping event %s", event.ID)
                }
            } else {
                rtsp.output <- processed
            }
        }
    }()

    // 后台处理缓冲区
    go rtsp.processBuffer()
}

func (rtsp *RealTimeStreamProcessor) processBuffer() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        for {
            event, ok := rtsp.buffer.Pop()
            if !ok {
                break
            }

            result, err := rtsp.processor.Process(event.Data)
            if err != nil {
                log.Printf("Failed to process buffered event %s: %v", event.ID, err)
                continue
            }

            processed := ProcessedEvent{
                EventID:   event.ID,
                Result:    result,
                Latency:   time.Since(event.Timestamp),
                Timestamp: time.Now(),
            }

            rtsp.output <- processed
        }
    }
}
```

## 8. 性能优化

### 8.1 并行处理

```go
// 并行处理器
type ParallelProcessor struct {
    workers   int
    processor DataProcessor
    pool      *WorkerPool
}

type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultChan chan Result
    wg         sync.WaitGroup
}

type Job struct {
    ID   string
    Data []byte
}

type Result struct {
    JobID  string
    Data   []byte
    Error  error
}

func (pp *ParallelProcessor) ProcessParallel(data [][]byte) ([][]byte, error) {
    // 创建作业
    jobs := make([]Job, len(data))
    for i, item := range data {
        jobs[i] = Job{
            ID:   fmt.Sprintf("job_%d", i),
            Data: item,
        }
    }

    // 提交作业
    for _, job := range jobs {
        pp.pool.jobQueue <- job
    }

    // 收集结果
    results := make([][]byte, len(jobs))
    for i := 0; i < len(jobs); i++ {
        result := <-pp.pool.resultChan
        if result.Error != nil {
            return nil, result.Error
        }
        results[i] = result.Data
    }

    return results, nil
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()

    for job := range wp.jobQueue {
        // 处理作业
        result := Result{
            JobID: job.ID,
        }

        // 这里应该调用实际的处理器
        // result.Data, result.Error = processor.Process(job.Data)

        wp.resultChan <- result
    }
}
```

### 8.2 内存优化

```go
// 内存优化的处理器
type MemoryOptimizedProcessor struct {
    processor DataProcessor
    pool      *sync.Pool
}

func (mop *MemoryOptimizedProcessor) Process(data []byte) ([]byte, error) {
    // 从对象池获取缓冲区
    buffer := mop.pool.Get().([]byte)
    defer mop.pool.Put(buffer)

    // 确保缓冲区足够大
    if cap(buffer) < len(data) {
        buffer = make([]byte, len(data))
    }
    buffer = buffer[:len(data)]

    // 复制数据
    copy(buffer, data)

    // 处理数据
    return mop.processor.Process(buffer)
}

// 使用示例
func NewMemoryOptimizedProcessor(processor DataProcessor) *MemoryOptimizedProcessor {
    return &MemoryOptimizedProcessor{
        processor: processor,
        pool: &sync.Pool{
            New: func() interface{} {
                return make([]byte, 1024)
            },
        },
    }
}
```

### 8.3 缓存优化

```go
// 缓存处理器
type CachedProcessor struct {
    processor DataProcessor
    cache     Cache
}

type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, value []byte) error
}

type MemoryCache struct {
    data map[string][]byte
    mu   sync.RWMutex
}

func (mc *MemoryCache) Get(key string) ([]byte, bool) {
    mc.mu.RLock()
    defer mc.mu.RUnlock()

    value, exists := mc.data[key]
    return value, exists
}

func (mc *MemoryCache) Set(key string, value []byte) error {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    mc.data[key] = value
    return nil
}

func (cp *CachedProcessor) Process(data []byte) ([]byte, error) {
    // 生成缓存键
    key := generateCacheKey(data)

    // 检查缓存
    if cached, exists := cp.cache.Get(key); exists {
        return cached, nil
    }

    // 处理数据
    result, err := cp.processor.Process(data)
    if err != nil {
        return nil, err
    }

    // 缓存结果
    cp.cache.Set(key, result)

    return result, nil
}

func generateCacheKey(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
```

## 总结

本章详细介绍了大数据处理的核心技术，包括：

1. **理论基础**：数据处理函数、处理模式
2. **数据清洗**：去重、补全、修正、标准化
3. **数据转换**：格式转换、编码转换、类型转换
4. **数据聚合**：统计聚合、分组聚合、窗口聚合
5. **流式处理**：实时流处理、滑动窗口
6. **批处理**：批量处理、MapReduce模式
7. **实时处理**：延迟约束、实时流处理
8. **性能优化**：并行处理、内存优化、缓存优化

这些技术为构建高效、可靠的大数据处理系统提供了完整的理论基础和实现方案。

---

**相关链接**：

- [01-数据摄入](../01-Data-Ingestion.md)
- [03-数据分析](../03-Data-Analytics.md)
- [04-数据存储](../04-Data-Storage.md)
- [05-数据可视化](../05-Data-Visualization.md)
- [06-机器学习](../06-Machine-Learning.md)
- [07-实时计算](../07-Real-Time-Computing.md)
- [08-数据治理](../08-Data-Governance.md)
