# 适配器模式 Go 语言实现

## 概述

本文档提供适配器模式的完整 Go 语言实现，包括对象适配器和类适配器两种实现方式，以及详细的性能分析和测试验证。

## 核心接口定义

### 目标接口和被适配接口

```go
package adapter

import (
    "fmt"
    "reflect"
    "time"
)

// Target 目标接口 - 客户端期望的接口
type Target interface {
    Request() string
    GetData() map[string]interface{}
    Process(input []byte) ([]byte, error)
}

// Adaptee 被适配接口 - 需要适配的现有接口
type Adaptee interface {
    SpecificRequest() string
    GetSpecificData() map[string]string
    ProcessSpecific(input string) (string, error)
}

// 具体被适配类
type ConcreteAdaptee struct {
    name string
    data map[string]string
}

// NewConcreteAdaptee 创建具体被适配对象
func NewConcreteAdaptee(name string) *ConcreteAdaptee {
    return &ConcreteAdaptee{
        name: name,
        data: make(map[string]string),
    }
}

// SpecificRequest 具体请求方法
func (a *ConcreteAdaptee) SpecificRequest() string {
    return fmt.Sprintf("Specific request from %s", a.name)
}

// GetSpecificData 获取具体数据
func (a *ConcreteAdaptee) GetSpecificData() map[string]string {
    return a.data
}

// ProcessSpecific 处理具体数据
func (a *ConcreteAdaptee) ProcessSpecific(input string) (string, error) {
    if input == "" {
        return "", fmt.Errorf("empty input")
    }
    return fmt.Sprintf("Processed: %s", input), nil
}

// SetData 设置数据
func (a *ConcreteAdaptee) SetData(key, value string) {
    a.data[key] = value
}
```

## 对象适配器实现

### 基本对象适配器

```go
// ObjectAdapter 对象适配器 - 通过组合实现适配
type ObjectAdapter struct {
    adaptee Adaptee
}

// NewObjectAdapter 创建对象适配器
func NewObjectAdapter(adaptee Adaptee) *ObjectAdapter {
    return &ObjectAdapter{
        adaptee: adaptee,
    }
}

// Request 实现目标接口的请求方法
func (a *ObjectAdapter) Request() string {
    // 调用被适配对象的方法并转换结果
    specificResult := a.adaptee.SpecificRequest()
    return fmt.Sprintf("Adapted: %s", specificResult)
}

// GetData 实现目标接口的数据获取方法
func (a *ObjectAdapter) GetData() map[string]interface{} {
    // 将 map[string]string 转换为 map[string]interface{}
    specificData := a.adaptee.GetSpecificData()
    result := make(map[string]interface{})
    
    for key, value := range specificData {
        result[key] = value
    }
    
    return result
}

// Process 实现目标接口的处理方法
func (a *ObjectAdapter) Process(input []byte) ([]byte, error) {
    // 将 []byte 转换为 string，调用被适配对象的方法
    inputStr := string(input)
    result, err := a.adaptee.ProcessSpecific(inputStr)
    if err != nil {
        return nil, err
    }
    
    // 将结果转换回 []byte
    return []byte(result), nil
}
```

### 高级对象适配器

```go
// AdvancedObjectAdapter 高级对象适配器 - 支持更多功能
type AdvancedObjectAdapter struct {
    adaptee    Adaptee
    cache      map[string]interface{}
    statistics map[string]int
}

// NewAdvancedObjectAdapter 创建高级对象适配器
func NewAdvancedObjectAdapter(adaptee Adaptee) *AdvancedObjectAdapter {
    return &AdvancedObjectAdapter{
        adaptee:    adaptee,
        cache:      make(map[string]interface{}),
        statistics: make(map[string]int),
    }
}

// Request 带缓存的请求方法
func (a *AdvancedObjectAdapter) Request() string {
    cacheKey := "request"
    
    // 检查缓存
    if cached, exists := a.cache[cacheKey]; exists {
        a.statistics["cache_hits"]++
        return cached.(string)
    }
    
    // 调用被适配对象
    result := a.adaptee.SpecificRequest()
    adaptedResult := fmt.Sprintf("Advanced Adapted: %s", result)
    
    // 缓存结果
    a.cache[cacheKey] = adaptedResult
    a.statistics["cache_misses"]++
    
    return adaptedResult
}

// GetData 带统计的数据获取方法
func (a *AdvancedObjectAdapter) GetData() map[string]interface{} {
    a.statistics["data_requests"]++
    
    specificData := a.adaptee.GetSpecificData()
    result := make(map[string]interface{})
    
    for key, value := range specificData {
        result[key] = value
    }
    
    // 添加统计信息
    result["statistics"] = a.statistics
    
    return result
}

// Process 带重试的处理方法
func (a *AdvancedObjectAdapter) Process(input []byte) ([]byte, error) {
    maxRetries := 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := a.adaptee.ProcessSpecific(string(input))
        if err == nil {
            a.statistics["process_success"]++
            return []byte(result), nil
        }
        
        lastErr = err
        a.statistics["process_retries"]++
        
        // 简单的重试延迟
        time.Sleep(time.Duration(i+1) * 10 * time.Millisecond)
    }
    
    a.statistics["process_failures"]++
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// GetStatistics 获取统计信息
func (a *AdvancedObjectAdapter) GetStatistics() map[string]int {
    return a.statistics
}

// ClearCache 清除缓存
func (a *AdvancedObjectAdapter) ClearCache() {
    a.cache = make(map[string]interface{})
}
```

## 类适配器实现

### 通过嵌入实现类适配器

```go
// ClassAdapter 类适配器 - 通过嵌入实现适配
type ClassAdapter struct {
    *ConcreteAdaptee // 嵌入被适配类
}

// NewClassAdapter 创建类适配器
func NewClassAdapter(name string) *ClassAdapter {
    return &ClassAdapter{
        ConcreteAdaptee: NewConcreteAdaptee(name),
    }
}

// Request 实现目标接口的请求方法
func (a *ClassAdapter) Request() string {
    // 直接调用嵌入对象的方法
    specificResult := a.SpecificRequest()
    return fmt.Sprintf("Class Adapted: %s", specificResult)
}

// GetData 实现目标接口的数据获取方法
func (a *ClassAdapter) GetData() map[string]interface{} {
    specificData := a.GetSpecificData()
    result := make(map[string]interface{})
    
    for key, value := range specificData {
        result[key] = value
    }
    
    return result
}

// Process 实现目标接口的处理方法
func (a *ClassAdapter) Process(input []byte) ([]byte, error) {
    result, err := a.ProcessSpecific(string(input))
    if err != nil {
        return nil, err
    }
    
    return []byte(result), nil
}
```

## 工厂模式创建适配器

```go
// AdapterType 适配器类型
type AdapterType int

const (
    ObjectAdapterType AdapterType = iota
    ClassAdapterType
    AdvancedAdapterType
)

// AdapterFactory 适配器工厂
type AdapterFactory struct{}

// NewAdapterFactory 创建适配器工厂
func NewAdapterFactory() *AdapterFactory {
    return &AdapterFactory{}
}

// CreateAdapter 创建适配器
func (f *AdapterFactory) CreateAdapter(adapterType AdapterType, adaptee Adaptee) Target {
    switch adapterType {
    case ObjectAdapterType:
        return NewObjectAdapter(adaptee)
    case ClassAdapterType:
        if concreteAdaptee, ok := adaptee.(*ConcreteAdaptee); ok {
            return &ClassAdapter{ConcreteAdaptee: concreteAdaptee}
        }
        // 如果不是 ConcreteAdaptee，回退到对象适配器
        return NewObjectAdapter(adaptee)
    case AdvancedAdapterType:
        return NewAdvancedObjectAdapter(adaptee)
    default:
        return NewObjectAdapter(adaptee)
    }
}

// CreateAdapterWithConfig 根据配置创建适配器
func (f *AdapterFactory) CreateAdapterWithConfig(config AdapterConfig) Target {
    adaptee := NewConcreteAdaptee(config.AdapteeName)
    
    // 设置初始数据
    for key, value := range config.InitialData {
        adaptee.SetData(key, value)
    }
    
    return f.CreateAdapter(config.AdapterType, adaptee)
}

// AdapterConfig 适配器配置
type AdapterConfig struct {
    AdapterType  AdapterType
    AdapteeName  string
    InitialData  map[string]string
    EnableCache  bool
    MaxRetries   int
}
```

## 使用示例

### 基本使用示例

```go
// ExampleBasicUsage 基本使用示例
func ExampleBasicUsage() {
    // 创建被适配对象
    adaptee := NewConcreteAdaptee("LegacySystem")
    adaptee.SetData("version", "1.0")
    adaptee.SetData("status", "active")
    
    // 创建对象适配器
    adapter := NewObjectAdapter(adaptee)
    
    // 使用目标接口
    fmt.Println(adapter.Request())
    fmt.Println(adapter.GetData())
    
    // 处理数据
    result, err := adapter.Process([]byte("test data"))
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %s\n", string(result))
    }
}
```

### 高级使用示例

```go
// ExampleAdvancedUsage 高级使用示例
func ExampleAdvancedUsage() {
    // 创建工厂
    factory := NewAdapterFactory()
    
    // 配置适配器
    config := AdapterConfig{
        AdapterType: AdvancedAdapterType,
        AdapteeName: "ProductionSystem",
        InitialData: map[string]string{
            "environment": "production",
            "version":     "2.0",
        },
        EnableCache: true,
        MaxRetries:  3,
    }
    
    // 创建高级适配器
    adapter := factory.CreateAdapterWithConfig(config)
    
    // 多次调用以测试缓存
    for i := 0; i < 3; i++ {
        fmt.Printf("Request %d: %s\n", i+1, adapter.Request())
    }
    
    // 获取统计信息
    if advancedAdapter, ok := adapter.(*AdvancedObjectAdapter); ok {
        fmt.Printf("Statistics: %+v\n", advancedAdapter.GetStatistics())
    }
}
```

## 性能分析

### 基准测试

```go
// BenchmarkObjectAdapter 对象适配器基准测试
func BenchmarkObjectAdapter(b *testing.B) {
    adaptee := NewConcreteAdaptee("BenchmarkSystem")
    adapter := NewObjectAdapter(adaptee)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        adapter.Request()
    }
}

// BenchmarkClassAdapter 类适配器基准测试
func BenchmarkClassAdapter(b *testing.B) {
    adapter := NewClassAdapter("BenchmarkSystem")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        adapter.Request()
    }
}

// BenchmarkAdvancedAdapter 高级适配器基准测试
func BenchmarkAdvancedAdapter(b *testing.B) {
    adaptee := NewConcreteAdaptee("BenchmarkSystem")
    adapter := NewAdvancedObjectAdapter(adaptee)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        adapter.Request()
    }
}

// BenchmarkDirectCall 直接调用基准测试
func BenchmarkDirectCall(b *testing.B) {
    adaptee := NewConcreteAdaptee("BenchmarkSystem")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        adaptee.SpecificRequest()
    }
}
```

### 性能对比

| 操作 | 对象适配器 | 类适配器 | 高级适配器 | 直接调用 |
|------|------------|----------|------------|----------|
| 创建时间 | O(1) | O(1) | O(1) | O(1) |
| 方法调用 | O(1) | O(1) | O(1) | O(1) |
| 内存使用 | 中等 | 低 | 高 | 最低 |
| 缓存命中 | 无 | 无 | O(1) | 无 |

## 测试验证

### 单元测试

```go
// TestObjectAdapter 测试对象适配器
func TestObjectAdapter(t *testing.T) {
    adaptee := NewConcreteAdaptee("TestSystem")
    adaptee.SetData("test_key", "test_value")
    
    adapter := NewObjectAdapter(adaptee)
    
    // 测试 Request 方法
    result := adapter.Request()
    expected := "Adapted: Specific request from TestSystem"
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
    
    // 测试 GetData 方法
    data := adapter.GetData()
    if data["test_key"] != "test_value" {
        t.Errorf("Expected test_value, got %v", data["test_key"])
    }
    
    // 测试 Process 方法
    resultBytes, err := adapter.Process([]byte("test input"))
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    
    expectedResult := "Processed: test input"
    if string(resultBytes) != expectedResult {
        t.Errorf("Expected %s, got %s", expectedResult, string(resultBytes))
    }
}

// TestClassAdapter 测试类适配器
func TestClassAdapter(t *testing.T) {
    adapter := NewClassAdapter("TestSystem")
    adapter.SetData("test_key", "test_value")
    
    // 测试 Request 方法
    result := adapter.Request()
    expected := "Class Adapted: Specific request from TestSystem"
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
    
    // 测试数据转换
    data := adapter.GetData()
    if data["test_key"] != "test_value" {
        t.Errorf("Expected test_value, got %v", data["test_key"])
    }
}

// TestAdvancedAdapter 测试高级适配器
func TestAdvancedAdapter(t *testing.T) {
    adaptee := NewConcreteAdaptee("TestSystem")
    adapter := NewAdvancedObjectAdapter(adaptee)
    
    // 测试缓存功能
    result1 := adapter.Request()
    result2 := adapter.Request()
    
    if result1 != result2 {
        t.Error("Cached results should be identical")
    }
    
    // 测试统计信息
    stats := adapter.GetStatistics()
    if stats["cache_hits"] != 1 || stats["cache_misses"] != 1 {
        t.Errorf("Unexpected statistics: %+v", stats)
    }
}

// TestAdapterFactory 测试适配器工厂
func TestAdapterFactory(t *testing.T) {
    factory := NewAdapterFactory()
    
    adaptee := NewConcreteAdaptee("TestSystem")
    
    // 测试对象适配器创建
    objectAdapter := factory.CreateAdapter(ObjectAdapterType, adaptee)
    if objectAdapter == nil {
        t.Error("Object adapter should not be nil")
    }
    
    // 测试高级适配器创建
    advancedAdapter := factory.CreateAdapter(AdvancedAdapterType, adaptee)
    if advancedAdapter == nil {
        t.Error("Advanced adapter should not be nil")
    }
}
```

### 集成测试

```go
// TestAdapterIntegration 集成测试
func TestAdapterIntegration(t *testing.T) {
    // 创建多个适配器
    adaptees := []Adaptee{
        NewConcreteAdaptee("System1"),
        NewConcreteAdaptee("System2"),
        NewConcreteAdaptee("System3"),
    }
    
    adapters := []Target{
        NewObjectAdapter(adaptees[0]),
        NewClassAdapter("System2"),
        NewAdvancedObjectAdapter(adaptees[2]),
    }
    
    // 测试所有适配器
    for i, adapter := range adapters {
        result := adapter.Request()
        if result == "" {
            t.Errorf("Adapter %d returned empty result", i)
        }
        
        data := adapter.GetData()
        if data == nil {
            t.Errorf("Adapter %d returned nil data", i)
        }
    }
}
```

## 错误处理

```go
// AdapterError 适配器错误类型
type AdapterError struct {
    Op   string
    Err  error
    Code string
}

func (e *AdapterError) Error() string {
    return fmt.Sprintf("adapter %s: %v (code: %s)", e.Op, e.Err, e.Code)
}

func (e *AdapterError) Unwrap() error {
    return e.Err
}

// SafeAdapter 安全适配器 - 包含错误处理
type SafeAdapter struct {
    adapter Target
    logger  Logger
}

// Logger 日志接口
type Logger interface {
    Log(level, message string)
}

// NewSafeAdapter 创建安全适配器
func NewSafeAdapter(adapter Target, logger Logger) *SafeAdapter {
    return &SafeAdapter{
        adapter: adapter,
        logger:  logger,
    }
}

// Request 带错误处理的请求方法
func (s *SafeAdapter) Request() string {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Log("ERROR", fmt.Sprintf("Panic in Request: %v", r))
        }
    }()
    
    result := s.adapter.Request()
    s.logger.Log("INFO", "Request completed successfully")
    return result
}

// GetData 带错误处理的数据获取方法
func (s *SafeAdapter) GetData() map[string]interface{} {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Log("ERROR", fmt.Sprintf("Panic in GetData: %v", r))
        }
    }()
    
    data := s.adapter.GetData()
    s.logger.Log("INFO", "GetData completed successfully")
    return data
}

// Process 带错误处理的处理方法
func (s *SafeAdapter) Process(input []byte) ([]byte, error) {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Log("ERROR", fmt.Sprintf("Panic in Process: %v", r))
        }
    }()
    
    result, err := s.adapter.Process(input)
    if err != nil {
        s.logger.Log("ERROR", fmt.Sprintf("Process failed: %v", err))
        return nil, &AdapterError{
            Op:   "Process",
            Err:  err,
            Code: "PROCESS_ERROR",
        }
    }
    
    s.logger.Log("INFO", "Process completed successfully")
    return result, nil
}
```

---

**构建状态**: ✅ 完成  
**最后更新**: 2024-01-06  
**版本**: v1.0.0  

<(￣︶￣)↗[GO!] Go实现，实践之本！
