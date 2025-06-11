# 6. 测试验证文档

## 6.1 测试概述

### 6.1.1 测试策略

葛洲坝船闸导航系统采用分层测试策略，确保系统的可靠性、稳定性和安全性：

```text
┌─────────────────────────────────────────────────────────┐
│                    端到端测试                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 业务流程测试 │ │ 系统集成测试 │ │ 性能压力测试 │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    集成测试                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 组件集成测试 │ │ 接口集成测试 │ │ 设备集成测试 │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    单元测试                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 函数单元测试 │ │ 组件单元测试 │ │ 算法单元测试 │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
```

### 6.1.2 测试目标

1. **功能验证**
   - 验证系统功能是否符合需求
   - 确保业务流程正确执行
   - 验证异常处理机制

2. **性能验证**
   - 验证系统性能指标
   - 确保系统在高负载下稳定运行
   - 验证资源使用效率

3. **安全验证**
   - 验证系统安全性
   - 确保数据完整性
   - 验证访问控制机制

4. **可靠性验证**
   - 验证系统稳定性
   - 确保故障恢复能力
   - 验证数据一致性

## 6.2 单元测试

### 6.2.1 测试框架

系统使用Go语言标准测试框架和第三方测试库：

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)
```

### 6.2.2 组件单元测试

```go
// 雷达数据处理单元测试
func TestRadarDataProcessor_ProcessTarget(t *testing.T) {
    // 准备测试数据
    processor := &RadarDataProcessor{
        targets: make(map[uint16]*RadarTarget),
    }
    
    target := &RadarTarget{
        TargetID:  1,
        Distance:  100.0,
        Angle:     45.0,
        Speed:     5.0,
        Quality:   80,
        Timestamp: 1234567890,
    }
    
    // 执行测试
    processor.ProcessTarget(target)
    
    // 验证结果
    assert.Equal(t, 1, len(processor.targets))
    assert.Equal(t, target, processor.targets[1])
}

// 船舶位置计算单元测试
func TestShipPositionCalculator_CalculatePosition(t *testing.T) {
    calculator := &ShipPositionCalculator{}
    
    tests := []struct {
        name     string
        target   *RadarTarget
        expected *ShipPosition
    }{
        {
            name: "正常位置计算",
            target: &RadarTarget{
                Distance: 100.0,
                Angle:    45.0,
                Speed:    5.0,
            },
            expected: &ShipPosition{
                X: 70.71, // 100 * cos(45°)
                Y: 70.71, // 100 * sin(45°)
                Speed: 5.0,
            },
        },
        {
            name: "零距离测试",
            target: &RadarTarget{
                Distance: 0.0,
                Angle:    0.0,
                Speed:    0.0,
            },
            expected: &ShipPosition{
                X: 0.0,
                Y: 0.0,
                Speed: 0.0,
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calculator.CalculatePosition(tt.target)
            
            assert.InDelta(t, tt.expected.X, result.X, 0.01)
            assert.InDelta(t, tt.expected.Y, result.Y, 0.01)
            assert.InDelta(t, tt.expected.Speed, result.Speed, 0.01)
        })
    }
}

// 禁停区域判断单元测试
func TestNoStopZoneChecker_IsInNoStopZone(t *testing.T) {
    checker := &NoStopZoneChecker{}
    
    // 定义禁停区域
    noStopZones := []NoStopZone{
        {
            Name: "测试禁停区",
            Boundary: []Point{
                {X: 0, Y: 0},
                {X: 100, Y: 0},
                {X: 100, Y: 50},
                {X: 0, Y: 50},
            },
        },
    }
    
    tests := []struct {
        name     string
        position ShipPosition
        expected bool
    }{
        {
            name: "在禁停区内",
            position: ShipPosition{X: 50, Y: 25},
            expected: true,
        },
        {
            name: "在禁停区外",
            position: ShipPosition{X: 150, Y: 25},
            expected: false,
        },
        {
            name: "在边界上",
            position: ShipPosition{X: 100, Y: 25},
            expected: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.IsInNoStopZone(tt.position, noStopZones)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 6.2.3 算法单元测试

```go
// 速度计算算法测试
func TestSpeedCalculator_CalculateSpeed(t *testing.T) {
    calculator := &SpeedCalculator{}
    
    tests := []struct {
        name      string
        positions []ShipPosition
        expected  float64
    }{
        {
            name: "正常速度计算",
            positions: []ShipPosition{
                {X: 0, Y: 0, Timestamp: time.Unix(0, 0)},
                {X: 10, Y: 0, Timestamp: time.Unix(1, 0)},
            },
            expected: 10.0, // 10米/秒
        },
        {
            name: "零速度",
            positions: []ShipPosition{
                {X: 0, Y: 0, Timestamp: time.Unix(0, 0)},
                {X: 0, Y: 0, Timestamp: time.Unix(1, 0)},
            },
            expected: 0.0,
        },
        {
            name: "复杂轨迹",
            positions: []ShipPosition{
                {X: 0, Y: 0, Timestamp: time.Unix(0, 0)},
                {X: 3, Y: 4, Timestamp: time.Unix(1, 0)},
            },
            expected: 5.0, // sqrt(3²+4²) = 5米/秒
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calculator.CalculateSpeed(tt.positions)
            assert.InDelta(t, tt.expected, result, 0.01)
        })
    }
}

// 轨迹分析算法测试
func TestTrajectoryAnalyzer_AnalyzeTrajectory(t *testing.T) {
    analyzer := &TrajectoryAnalyzer{}
    
    trajectory := []ShipPosition{
        {X: 0, Y: 0, Timestamp: time.Unix(0, 0)},
        {X: 10, Y: 0, Timestamp: time.Unix(1, 0)},
        {X: 20, Y: 0, Timestamp: time.Unix(2, 0)},
        {X: 30, Y: 0, Timestamp: time.Unix(3, 0)},
    }
    
    result := analyzer.AnalyzeTrajectory(trajectory)
    
    assert.Equal(t, "straight", result.Type)
    assert.InDelta(t, 10.0, result.AverageSpeed, 0.01)
    assert.Equal(t, 0.0, result.Direction)
}
```

### 6.2.4 配置单元测试

```go
// 配置加载器测试
func TestConfigLoader_Load(t *testing.T) {
    // 创建临时配置文件
    tempFile, err := os.CreateTemp("", "test_config_*.yaml")
    assert.NoError(t, err)
    defer os.Remove(tempFile.Name())
    
    configContent := `
system:
  name: "测试系统"
  version: "1.0.0"
  environment: "test"
logging:
  level: "debug"
  output_path: "stdout"
`
    _, err = tempFile.WriteString(configContent)
    assert.NoError(t, err)
    tempFile.Close()
    
    // 测试配置加载
    loader := NewConfigLoader()
    loader.viper.SetConfigFile(tempFile.Name())
    
    config, err := loader.Load()
    assert.NoError(t, err)
    assert.Equal(t, "测试系统", config.System.Name)
    assert.Equal(t, "1.0.0", config.System.Version)
    assert.Equal(t, "test", config.System.Environment)
    assert.Equal(t, "debug", config.Logging.Level)
}

// 配置验证器测试
func TestConfigValidator_Validate(t *testing.T) {
    validator := &ConfigValidator{}
    
    tests := []struct {
        name   string
        config *SystemConfig
        hasErr bool
    }{
        {
            name: "有效配置",
            config: &SystemConfig{
                System: SystemConfig{
                    Name:    "测试系统",
                    Version: "1.0.0",
                },
                Network: NetworkConfig{
                    NATS: NATSConfig{URL: "nats://localhost:4222"},
                    HTTP: HTTPConfig{Port: 8080},
                },
            },
            hasErr: false,
        },
        {
            name: "缺少系统名称",
            config: &SystemConfig{
                System: SystemConfig{
                    Version: "1.0.0",
                },
            },
            hasErr: true,
        },
        {
            name: "无效端口",
            config: &SystemConfig{
                System: SystemConfig{
                    Name:    "测试系统",
                    Version: "1.0.0",
                },
                Network: NetworkConfig{
                    NATS: NATSConfig{URL: "nats://localhost:4222"},
                    HTTP: HTTPConfig{Port: 70000}, // 无效端口
                },
            },
            hasErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.Validate(tt.config)
            if tt.hasErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 6.3 集成测试

### 6.3.1 组件集成测试

```go
// 组件集成测试套件
type ComponentIntegrationTestSuite struct {
    suite.Suite
    serverMaster *ServiceMasterAdapter
    mockNATS     *MockNATSConnection
    mockDB       *MockDatabase
}

func (suite *ComponentIntegrationTestSuite) SetupSuite() {
    // 设置测试环境
    suite.mockNATS = NewMockNATSConnection()
    suite.mockDB = NewMockDatabase()
    
    // 创建服务主控
    ctrl := mdl.NewCtrlSt(context.Background())
    suite.serverMaster = NewServiceMasterAdapter(ctrl)
}

func (suite *ComponentIntegrationTestSuite) TestComponentStartup() {
    // 测试组件启动
    err := suite.serverMaster.Start()
    suite.NoError(err)
    
    // 验证组件状态
    suite.True(suite.serverMaster.IsRunning())
    
    // 验证子组件启动
    components := suite.serverMaster.GetComponents()
    suite.NotEmpty(components)
    
    for _, component := range components {
        suite.True(component.IsRunning())
    }
}

func (suite *ComponentIntegrationTestSuite) TestComponentCommunication() {
    // 启动组件
    err := suite.serverMaster.Start()
    suite.NoError(err)
    defer suite.serverMaster.Stop()
    
    // 发送测试消息
    testMessage := &TestMessage{
        Type: "test",
        Data: "test data",
    }
    
    err = suite.mockNATS.Publish("test.topic", testMessage)
    suite.NoError(err)
    
    // 等待消息处理
    time.Sleep(100 * time.Millisecond)
    
    // 验证消息处理
    receivedMessages := suite.mockNATS.GetReceivedMessages()
    suite.NotEmpty(receivedMessages)
}

func (suite *ComponentIntegrationTestSuite) TestComponentShutdown() {
    // 启动组件
    err := suite.serverMaster.Start()
    suite.NoError(err)
    
    // 停止组件
    err = suite.serverMaster.Stop()
    suite.NoError(err)
    
    // 验证组件状态
    suite.False(suite.serverMaster.IsRunning())
    
    // 验证子组件停止
    components := suite.serverMaster.GetComponents()
    for _, component := range components {
        suite.False(component.IsRunning())
    }
}

func TestComponentIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(ComponentIntegrationTestSuite))
}
```

### 6.3.2 设备集成测试

```go
// 设备集成测试
func TestDeviceIntegration(t *testing.T) {
    // 创建模拟设备
    mockRadar := NewMockRadarDevice()
    mockPTZ := NewMockPTZDevice()
    mockLED := NewMockLEDDevice()
    
    // 创建设备管理器
    deviceManager := NewDeviceManager()
    deviceManager.AddDevice(mockRadar)
    deviceManager.AddDevice(mockPTZ)
    deviceManager.AddDevice(mockLED)
    
    // 测试设备连接
    err := deviceManager.ConnectAll()
    assert.NoError(t, err)
    
    // 测试雷达数据采集
    radarData := &RadarData{
        Targets: []RadarTarget{
            {TargetID: 1, Distance: 100.0, Angle: 45.0, Speed: 5.0},
        },
    }
    
    mockRadar.SetTestData(radarData)
    
    // 等待数据处理
    time.Sleep(100 * time.Millisecond)
    
    // 验证数据处理
    processedData := deviceManager.GetProcessedData()
    assert.NotEmpty(t, processedData)
    
    // 测试云台控制
    err = mockPTZ.PanLeft(50)
    assert.NoError(t, err)
    
    // 验证云台状态
    status, err := mockPTZ.GetStatus()
    assert.NoError(t, err)
    assert.Equal(t, "panning_left", status.State)
    
    // 测试LED显示
    err = mockLED.SetText("Test Message", 0, 0, 255)
    assert.NoError(t, err)
    
    // 验证LED显示内容
    displayContent := mockLED.GetDisplayContent()
    assert.Equal(t, "Test Message", displayContent.Text)
}
```

### 6.3.3 数据库集成测试

```go
// 数据库集成测试
func TestDatabaseIntegration(t *testing.T) {
    // 创建测试数据库
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)
    defer db.Close()
    
    // 创建测试表
    _, err = db.Exec(`
        CREATE TABLE ships (
            id INTEGER PRIMARY KEY,
            ship_id TEXT NOT NULL,
            ship_name TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    assert.NoError(t, err)
    
    // 创建数据访问层
    shipRepo := NewShipRepository(db)
    
    // 测试数据插入
    ship := &Ship{
        ShipID:   "TEST001",
        ShipName: "测试船舶",
    }
    
    err = shipRepo.Create(ship)
    assert.NoError(t, err)
    
    // 测试数据查询
    retrievedShip, err := shipRepo.GetByID("TEST001")
    assert.NoError(t, err)
    assert.Equal(t, "TEST001", retrievedShip.ShipID)
    assert.Equal(t, "测试船舶", retrievedShip.ShipName)
    
    // 测试数据更新
    ship.ShipName = "更新后的船舶名称"
    err = shipRepo.Update(ship)
    assert.NoError(t, err)
    
    updatedShip, err := shipRepo.GetByID("TEST001")
    assert.NoError(t, err)
    assert.Equal(t, "更新后的船舶名称", updatedShip.ShipName)
    
    // 测试数据删除
    err = shipRepo.Delete("TEST001")
    assert.NoError(t, err)
    
    _, err = shipRepo.GetByID("TEST001")
    assert.Error(t, err) // 应该返回错误，因为数据已删除
}
```

## 6.4 系统测试

### 6.4.1 业务流程测试

```go
// 业务流程测试
func TestShipPassageWorkflow(t *testing.T) {
    // 创建测试系统
    system := NewTestSystem()
    defer system.Cleanup()
    
    // 启动系统
    err := system.Start()
    assert.NoError(t, err)
    
    // 模拟船舶进入上行航道
    ship := &Ship{
        ShipID:   "TEST001",
        ShipName: "测试船舶",
        Length:   100.0,
        Width:    20.0,
    }
    
    // 模拟雷达检测到船舶
    radarData := &RadarData{
        Targets: []RadarTarget{
            {
                TargetID:  1,
                Distance:  1000.0,
                Angle:     45.0,
                Speed:     5.0,
                Quality:   90,
                Timestamp: time.Now().Unix(),
            },
        },
    }
    
    system.SimulateRadarData(radarData)
    
    // 等待系统处理
    time.Sleep(1 * time.Second)
    
    // 验证船舶被识别
    detectedShips := system.GetDetectedShips()
    assert.NotEmpty(t, detectedShips)
    assert.Equal(t, "TEST001", detectedShips[0].ShipID)
    
    // 模拟船舶移动到禁停区域
    radarData.Targets[0].Distance = 500.0
    system.SimulateRadarData(radarData)
    
    time.Sleep(1 * time.Second)
    
    // 验证禁停警告
    warnings := system.GetWarnings()
    assert.NotEmpty(t, warnings)
    assert.Contains(t, warnings[0].Message, "禁停区域")
    
    // 模拟船舶超速
    radarData.Targets[0].Speed = 20.0
    system.SimulateRadarData(radarData)
    
    time.Sleep(1 * time.Second)
    
    // 验证超速警告
    warnings = system.GetWarnings()
    assert.NotEmpty(t, warnings)
    assert.Contains(t, warnings[0].Message, "超速")
    
    // 模拟船舶正常通过
    radarData.Targets[0].Speed = 5.0
    radarData.Targets[0].Distance = 100.0
    system.SimulateRadarData(radarData)
    
    time.Sleep(1 * time.Second)
    
    // 验证船闸操作
    lockOperations := system.GetLockOperations()
    assert.NotEmpty(t, lockOperations)
    assert.Equal(t, "open", lockOperations[0].Action)
}
```

### 6.4.2 异常处理测试

```go
// 异常处理测试
func TestExceptionHandling(t *testing.T) {
    system := NewTestSystem()
    defer system.Cleanup()
    
    err := system.Start()
    assert.NoError(t, err)
    
    // 测试设备故障
    system.SimulateDeviceFailure("radar-upstream")
    
    time.Sleep(1 * time.Second)
    
    // 验证故障检测
    failures := system.GetDeviceFailures()
    assert.NotEmpty(t, failures)
    assert.Equal(t, "radar-upstream", failures[0].DeviceID)
    
    // 验证系统降级运行
    assert.True(t, system.IsRunning())
    
    // 测试设备恢复
    system.SimulateDeviceRecovery("radar-upstream")
    
    time.Sleep(1 * time.Second)
    
    // 验证设备恢复
    failures = system.GetDeviceFailures()
    assert.Empty(t, failures)
    
    // 测试网络中断
    system.SimulateNetworkFailure()
    
    time.Sleep(1 * time.Second)
    
    // 验证网络故障处理
    networkStatus := system.GetNetworkStatus()
    assert.False(t, networkStatus.IsConnected)
    
    // 验证本地模式运行
    assert.True(t, system.IsRunning())
    
    // 测试网络恢复
    system.SimulateNetworkRecovery()
    
    time.Sleep(1 * time.Second)
    
    // 验证网络恢复
    networkStatus = system.GetNetworkStatus()
    assert.True(t, networkStatus.IsConnected)
}
```

## 6.5 性能测试

### 6.5.1 负载测试

```go
// 负载测试
func TestSystemLoad(t *testing.T) {
    system := NewTestSystem()
    defer system.Cleanup()
    
    err := system.Start()
    assert.NoError(t, err)
    
    // 模拟高负载场景
    const numShips = 100
    const testDuration = 30 * time.Second
    
    startTime := time.Now()
    
    // 并发模拟多艘船舶
    var wg sync.WaitGroup
    for i := 0; i < numShips; i++ {
        wg.Add(1)
        go func(shipID int) {
            defer wg.Done()
            
            ship := &Ship{
                ShipID:   fmt.Sprintf("LOAD%03d", shipID),
                ShipName: fmt.Sprintf("负载测试船舶%d", shipID),
            }
            
            // 模拟船舶通过流程
            system.SimulateShipPassage(ship)
        }(i)
    }
    
    wg.Wait()
    
    endTime := time.Now()
    duration := endTime.Sub(startTime)
    
    // 验证性能指标
    stats := system.GetPerformanceStats()
    
    t.Logf("测试持续时间: %v", duration)
    t.Logf("处理船舶数量: %d", stats.ProcessedShips)
    t.Logf("平均处理时间: %v", stats.AverageProcessingTime)
    t.Logf("最大内存使用: %d MB", stats.MaxMemoryUsage)
    t.Logf("CPU使用率: %.2f%%", stats.CPUUsage)
    
    // 性能断言
    assert.Less(t, duration, testDuration)
    assert.Equal(t, numShips, stats.ProcessedShips)
    assert.Less(t, stats.AverageProcessingTime, 100*time.Millisecond)
    assert.Less(t, stats.MaxMemoryUsage, 1024) // 小于1GB
    assert.Less(t, stats.CPUUsage, 80.0) // 小于80%
}
```

### 6.5.2 压力测试

```go
// 压力测试
func TestSystemStress(t *testing.T) {
    system := NewTestSystem()
    defer system.Cleanup()
    
    err := system.Start()
    assert.NoError(t, err)
    
    // 压力测试参数
    const numConcurrentShips = 1000
    const testDuration = 60 * time.Second
    
    startTime := time.Now()
    
    // 创建压力测试场景
    stressTest := NewStressTest(system)
    stressTest.SetConcurrentShips(numConcurrentShips)
    stressTest.SetDuration(testDuration)
    
    // 执行压力测试
    results := stressTest.Run()
    
    endTime := time.Now()
    duration := endTime.Sub(startTime)
    
    // 验证压力测试结果
    t.Logf("压力测试持续时间: %v", duration)
    t.Logf("总处理船舶数: %d", results.TotalShips)
    t.Logf("成功处理数: %d", results.SuccessfulShips)
    t.Logf("失败处理数: %d", results.FailedShips)
    t.Logf("平均响应时间: %v", results.AverageResponseTime)
    t.Logf("最大响应时间: %v", results.MaxResponseTime)
    t.Logf("吞吐量: %.2f ships/sec", results.Throughput)
    
    // 压力测试断言
    assert.Greater(t, results.SuccessfulShips, results.TotalShips*9/10) // 成功率>90%
    assert.Less(t, results.AverageResponseTime, 200*time.Millisecond)
    assert.Greater(t, results.Throughput, 10.0) // 吞吐量>10 ships/sec
}
```

### 6.5.3 内存泄漏测试

```go
// 内存泄漏测试
func TestMemoryLeak(t *testing.T) {
    system := NewTestSystem()
    defer system.Cleanup()
    
    err := system.Start()
    assert.NoError(t, err)
    
    // 获取初始内存使用
    initialMemory := system.GetMemoryUsage()
    
    // 执行多次操作
    const numIterations = 1000
    for i := 0; i < numIterations; i++ {
        ship := &Ship{
            ShipID:   fmt.Sprintf("MEM%04d", i),
            ShipName: fmt.Sprintf("内存测试船舶%d", i),
        }
        
        system.SimulateShipPassage(ship)
        
        // 强制垃圾回收
        if i%100 == 0 {
            runtime.GC()
        }
    }
    
    // 强制垃圾回收
    runtime.GC()
    
    // 获取最终内存使用
    finalMemory := system.GetMemoryUsage()
    memoryIncrease := finalMemory - initialMemory
    
    t.Logf("初始内存使用: %d MB", initialMemory)
    t.Logf("最终内存使用: %d MB", finalMemory)
    t.Logf("内存增长: %d MB", memoryIncrease)
    
    // 验证内存泄漏
    assert.Less(t, memoryIncrease, 100) // 内存增长应小于100MB
}
```

## 6.6 安全测试

### 6.6.1 输入验证测试

```go
// 输入验证测试
func TestInputValidation(t *testing.T) {
    validator := NewInputValidator()
    
    tests := []struct {
        name  string
        input interface{}
        valid bool
    }{
        {
            name:  "有效船舶ID",
            input: "SHIP001",
            valid: true,
        },
        {
            name:  "无效船舶ID（包含特殊字符）",
            input: "SHIP@001",
            valid: false,
        },
        {
            name:  "有效速度",
            input: 10.5,
            valid: true,
        },
        {
            name:  "无效速度（负数）",
            input: -5.0,
            valid: false,
        },
        {
            name:  "有效坐标",
            input: Point{X: 100.0, Y: 200.0},
            valid: true,
        },
        {
            name:  "无效坐标（超出范围）",
            input: Point{X: 10000.0, Y: 20000.0},
            valid: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.Validate(tt.input)
            if tt.valid {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
        })
    }
}
```

### 6.6.2 权限控制测试

```go
// 权限控制测试
func TestAccessControl(t *testing.T) {
    authManager := NewAuthManager()
    
    // 创建测试用户
    adminUser := &User{
        ID:       "admin",
        Username: "admin",
        Role:     "admin",
    }
    
    operatorUser := &User{
        ID:       "operator",
        Username: "operator",
        Role:     "operator",
    }
    
    viewerUser := &User{
        ID:       "viewer",
        Username: "viewer",
        Role:     "viewer",
    }
    
    tests := []struct {
        name     string
        user     *User
        action   string
        resource string
        allowed  bool
    }{
        {
            name:     "管理员可以执行所有操作",
            user:     adminUser,
            action:   "write",
            resource: "system_config",
            allowed:  true,
        },
        {
            name:     "操作员可以查看配置",
            user:     operatorUser,
            action:   "read",
            resource: "system_config",
            allowed:  true,
        },
        {
            name:     "操作员不能修改配置",
            user:     operatorUser,
            action:   "write",
            resource: "system_config",
            allowed:  false,
        },
        {
            name:     "查看者只能查看数据",
            user:     viewerUser,
            action:   "read",
            resource: "ship_data",
            allowed:  true,
        },
        {
            name:     "查看者不能修改数据",
            user:     viewerUser,
            action:   "write",
            resource: "ship_data",
            allowed:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            allowed := authManager.CheckPermission(tt.user, tt.action, tt.resource)
            assert.Equal(t, tt.allowed, allowed)
        })
    }
}
```

## 6.7 测试报告

### 6.7.1 测试报告生成

```go
// 测试报告生成器
type TestReportGenerator struct {
    results []TestResult
}

type TestResult struct {
    TestName    string
    Status      string
    Duration    time.Duration
    Error       string
    Metrics     map[string]interface{}
}

func (g *TestReportGenerator) AddResult(result TestResult) {
    g.results = append(g.results, result)
}

func (g *TestReportGenerator) GenerateReport() *TestReport {
    report := &TestReport{
        Timestamp: time.Now(),
        Results:   g.results,
    }
    
    // 计算统计信息
    totalTests := len(g.results)
    passedTests := 0
    failedTests := 0
    totalDuration := time.Duration(0)
    
    for _, result := range g.results {
        totalDuration += result.Duration
        if result.Status == "PASS" {
            passedTests++
        } else {
            failedTests++
        }
    }
    
    report.Summary = TestSummary{
        TotalTests:    totalTests,
        PassedTests:   passedTests,
        FailedTests:   failedTests,
        SuccessRate:   float64(passedTests) / float64(totalTests) * 100,
        TotalDuration: totalDuration,
        AverageDuration: totalDuration / time.Duration(totalTests),
    }
    
    return report
}

func (g *TestReportGenerator) ExportToHTML(report *TestReport, filepath string) error {
    // 生成HTML报告
    template := `
<!DOCTYPE html>
<html>
<head>
    <title>测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .summary { background-color: #f0f0f0; padding: 15px; border-radius: 5px; }
        .test-result { margin: 10px 0; padding: 10px; border-left: 4px solid #ccc; }
        .pass { border-left-color: #4CAF50; }
        .fail { border-left-color: #f44336; }
    </style>
</head>
<body>
    <h1>测试报告</h1>
    <div class="summary">
        <h2>测试摘要</h2>
        <p>总测试数: {{.Summary.TotalTests}}</p>
        <p>通过测试: {{.Summary.PassedTests}}</p>
        <p>失败测试: {{.Summary.FailedTests}}</p>
        <p>成功率: {{.Summary.SuccessRate}}%</p>
        <p>总耗时: {{.Summary.TotalDuration}}</p>
    </div>
    
    <h2>详细结果</h2>
    {{range .Results}}
    <div class="test-result {{.Status}}">
        <h3>{{.TestName}}</h3>
        <p>状态: {{.Status}}</p>
        <p>耗时: {{.Duration}}</p>
        {{if .Error}}
        <p>错误: {{.Error}}</p>
        {{end}}
    </div>
    {{end}}
</body>
</html>
    `
    
    tmpl, err := template.New("report").Parse(template)
    if err != nil {
        return err
    }
    
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return tmpl.Execute(file, report)
}
```

### 6.7.2 测试覆盖率报告

```go
// 测试覆盖率报告
func TestCoverage(t *testing.T) {
    // 运行测试覆盖率
    cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
    err := cmd.Run()
    assert.NoError(t, err)
    
    // 生成HTML覆盖率报告
    cmd = exec.Command("go", "tool", "cover", "-html=coverage.out", "-o=coverage.html")
    err = cmd.Run()
    assert.NoError(t, err)
    
    // 读取覆盖率数据
    coverageData, err := os.ReadFile("coverage.out")
    assert.NoError(t, err)
    
    // 解析覆盖率
    coverage := parseCoverage(coverageData)
    
    t.Logf("测试覆盖率: %.2f%%", coverage.Percentage)
    t.Logf("覆盖行数: %d", coverage.CoveredLines)
    t.Logf("总行数: %d", coverage.TotalLines)
    
    // 覆盖率断言
    assert.Greater(t, coverage.Percentage, 80.0) // 覆盖率应大于80%
}
```

## 6.8 自动化测试

### 6.8.1 持续集成测试

```yaml
# .github/workflows/test.yml
name: 测试

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: 设置Go环境
      uses: actions/setup-go@v2
      with:
        go-version: 1.24
    
    - name: 下载依赖
      run: go mod download
    
    - name: 运行单元测试
      run: go test -v ./...
    
    - name: 运行集成测试
      run: go test -v -tags=integration ./...
    
    - name: 生成测试覆盖率
      run: go test -coverprofile=coverage.out ./...
    
    - name: 上传覆盖率报告
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

### 6.8.2 测试脚本

```bash
#!/bin/bash
# test.sh

echo "开始运行测试..."

# 运行单元测试
echo "运行单元测试..."
go test -v ./...

# 运行集成测试
echo "运行集成测试..."
go test -v -tags=integration ./...

# 运行性能测试
echo "运行性能测试..."
go test -v -tags=performance ./...

# 生成测试报告
echo "生成测试报告..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 检查测试结果
if [ $? -eq 0 ]; then
    echo "所有测试通过！"
    exit 0
else
    echo "测试失败！"
    exit 1
fi
```
