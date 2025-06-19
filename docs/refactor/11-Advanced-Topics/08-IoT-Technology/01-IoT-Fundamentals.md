# 11.8.1 物联网基础理论

## 11.8.1.1 概述

物联网（Internet of Things, IoT）是通过网络连接各种物理设备，实现数据采集、传输、处理和应用的技术体系。

### 11.8.1.1.1 基本概念

**定义 11.8.1.1** (物联网)
物联网是一个由相互关联的计算设备、机械和数字机器、物体、动物或人组成的系统，这些系统具有唯一标识符，能够通过网络传输数据。

**定义 11.8.1.2** (IoT设备)
IoT设备是具有感知能力、通信能力、处理能力和标识能力的物理实体。

### 11.8.1.1.2 IoT架构层次

```go
// IoT架构层次枚举
type IoTLayer int

const (
    PerceptionLayer IoTLayer = iota    // 感知层
    NetworkLayer                       // 网络层
    PlatformLayer                      // 平台层
    ApplicationLayer                   // 应用层
)

// IoT设备类型
type DeviceType int

const (
    SensorDevice DeviceType = iota     // 传感器设备
    ActuatorDevice                     // 执行器设备
    GatewayDevice                      // 网关设备
    EdgeDevice                         // 边缘设备
    CloudDevice                        // 云设备
)
```

## 11.8.1.2 IoT架构理论

### 11.8.1.2.1 三层架构模型

**定义 11.8.1.3** (IoT三层架构)
IoT系统采用三层架构：
1. **感知层**: 负责数据采集和物理世界感知
2. **网络层**: 负责数据传输和通信
3. **应用层**: 负责数据处理和业务应用

**定理 11.8.1.1** (IoT架构可扩展性)
在IoT三层架构中，每层的独立性和标准化程度决定了系统的可扩展性。

**证明**:
设系统复杂度为 $C$，层数为 $n$，每层标准化程度为 $s_i$。
则系统可扩展性为：
$$S = \prod_{i=1}^n s_i \cdot \frac{1}{C}$$

### 11.8.1.2.2 Go实现架构框架

```go
// IoT架构接口
type IoTLayer interface {
    Process(data interface{}) error
    GetLayerType() IoTLayer
    GetStatus() LayerStatus
}

// 层状态
type LayerStatus struct {
    IsActive    bool
    Performance float64
    ErrorRate   float64
    LastUpdate  time.Time
}

// 感知层实现
type PerceptionLayer struct {
    sensors    map[string]Sensor
    actuators  map[string]Actuator
    status     LayerStatus
}

// 传感器接口
type Sensor interface {
    GetID() string
    GetType() string
    Read() (interface{}, error)
    GetStatus() DeviceStatus
}

// 执行器接口
type Actuator interface {
    GetID() string
    GetType() string
    Execute(command interface{}) error
    GetStatus() DeviceStatus
}

// 设备状态
type DeviceStatus struct {
    IsOnline    bool
    Battery     float64
    Signal      float64
    LastSeen    time.Time
    ErrorCount  int
}

// 创建感知层
func NewPerceptionLayer() *PerceptionLayer {
    return &PerceptionLayer{
        sensors:   make(map[string]Sensor),
        actuators: make(map[string]Actuator),
        status: LayerStatus{
            IsActive:    true,
            Performance: 1.0,
            ErrorRate:   0.0,
            LastUpdate:  time.Now(),
        },
    }
}

// 添加传感器
func (pl *PerceptionLayer) AddSensor(sensor Sensor) {
    pl.sensors[sensor.GetID()] = sensor
}

// 添加执行器
func (pl *PerceptionLayer) AddActuator(actuator Actuator) {
    pl.actuators[actuator.GetID()] = actuator
}

// 处理数据
func (pl *PerceptionLayer) Process(data interface{}) error {
    // 读取传感器数据
    sensorData := make(map[string]interface{})
    for id, sensor := range pl.sensors {
        if reading, err := sensor.Read(); err == nil {
            sensorData[id] = reading
        } else {
            pl.status.ErrorRate += 0.1
        }
    }
    
    // 处理执行器命令
    if commands, ok := data.(map[string]interface{}); ok {
        for id, command := range commands {
            if actuator, exists := pl.actuators[id]; exists {
                if err := actuator.Execute(command); err != nil {
                    pl.status.ErrorRate += 0.1
                }
            }
        }
    }
    
    pl.status.LastUpdate = time.Now()
    return nil
}

// 获取层类型
func (pl *PerceptionLayer) GetLayerType() IoTLayer {
    return PerceptionLayer
}

// 获取状态
func (pl *PerceptionLayer) GetStatus() LayerStatus {
    return pl.status
}
```

## 11.8.1.3 IoT协议栈

### 11.8.1.3.1 协议层次结构

**定义 11.8.1.4** (IoT协议栈)
IoT协议栈分为四层：
1. **应用层**: MQTT、CoAP、HTTP
2. **传输层**: TCP、UDP、DTLS
3. **网络层**: IPv6、6LoWPAN
4. **物理层**: WiFi、蓝牙、LoRa、NB-IoT

### 11.8.1.3.2 MQTT协议

**定义 11.8.1.5** (MQTT协议)
MQTT（Message Queuing Telemetry Transport）是一种轻量级的发布/订阅消息传输协议。

**定理 11.8.1.2** (MQTT效率)
MQTT协议的消息开销比HTTP协议低约80%。

### 11.8.1.3.3 Go实现MQTT

```go
// MQTT客户端
type MQTTClient struct {
    clientID   string
    broker     string
    port       int
    username   string
    password   string
    client     mqtt.Client
    topics     map[string]MessageHandler
    isConnected bool
}

// 消息处理器
type MessageHandler func(topic string, payload []byte)

// 创建MQTT客户端
func NewMQTTClient(clientID, broker string, port int) *MQTTClient {
    return &MQTTClient{
        clientID: clientID,
        broker:   broker,
        port:     port,
        topics:   make(map[string]MessageHandler),
    }
}

// 连接到MQTT代理
func (mc *MQTTClient) Connect() error {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mc.broker, mc.port))
    opts.SetClientID(mc.clientID)
    
    if mc.username != "" {
        opts.SetUsername(mc.username)
        opts.SetPassword(mc.password)
    }
    
    mc.client = mqtt.NewClient(opts)
    if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    mc.isConnected = true
    return nil
}

// 订阅主题
func (mc *MQTTClient) Subscribe(topic string, handler MessageHandler) error {
    if !mc.isConnected {
        return fmt.Errorf("client not connected")
    }
    
    mc.topics[topic] = handler
    
    token := mc.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
        if handler, exists := mc.topics[msg.Topic()]; exists {
            handler(msg.Topic(), msg.Payload())
        }
    })
    
    if token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    return nil
}

// 发布消息
func (mc *MQTTClient) Publish(topic string, payload []byte) error {
    if !mc.isConnected {
        return fmt.Errorf("client not connected")
    }
    
    token := mc.client.Publish(topic, 0, false, payload)
    if token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    return nil
}

// 断开连接
func (mc *MQTTClient) Disconnect() {
    if mc.isConnected {
        mc.client.Disconnect(250)
        mc.isConnected = false
    }
}
```

## 11.8.1.4 设备管理

### 11.8.1.4.1 设备生命周期

**定义 11.8.1.6** (设备生命周期)
IoT设备的生命周期包括注册、认证、配置、监控、维护和退役。

### 11.8.1.4.2 Go实现设备管理

```go
// 设备管理器
type DeviceManager struct {
    devices    map[string]*Device
    registry   *DeviceRegistry
    monitor    *DeviceMonitor
    config     *DeviceConfig
}

// 设备结构
type Device struct {
    ID           string
    Name         string
    Type         DeviceType
    Status       DeviceStatus
    Capabilities []string
    Config       map[string]interface{}
    LastSeen     time.Time
}

// 设备注册表
type DeviceRegistry struct {
    devices map[string]*Device
    mutex   sync.RWMutex
}

// 创建设备注册表
func NewDeviceRegistry() *DeviceRegistry {
    return &DeviceRegistry{
        devices: make(map[string]*Device),
    }
}

// 注册设备
func (dr *DeviceRegistry) RegisterDevice(device *Device) error {
    dr.mutex.Lock()
    defer dr.mutex.Unlock()
    
    if _, exists := dr.devices[device.ID]; exists {
        return fmt.Errorf("device already registered")
    }
    
    dr.devices[device.ID] = device
    return nil
}

// 获取设备
func (dr *DeviceRegistry) GetDevice(id string) (*Device, error) {
    dr.mutex.RLock()
    defer dr.mutex.RUnlock()
    
    device, exists := dr.devices[id]
    if !exists {
        return nil, fmt.Errorf("device not found")
    }
    
    return device, nil
}

// 更新设备状态
func (dr *DeviceRegistry) UpdateDeviceStatus(id string, status DeviceStatus) error {
    dr.mutex.Lock()
    defer dr.mutex.Unlock()
    
    device, exists := dr.devices[id]
    if !exists {
        return fmt.Errorf("device not found")
    }
    
    device.Status = status
    device.LastSeen = time.Now()
    return nil
}

// 创建设备管理器
func NewDeviceManager() *DeviceManager {
    return &DeviceManager{
        devices:  make(map[string]*Device),
        registry: NewDeviceRegistry(),
        monitor:  NewDeviceMonitor(),
        config:   NewDeviceConfig(),
    }
}

// 注册设备
func (dm *DeviceManager) RegisterDevice(device *Device) error {
    return dm.registry.RegisterDevice(device)
}

// 更新设备状态
func (dm *DeviceManager) UpdateDeviceStatus(deviceID string, status DeviceStatus) error {
    return dm.registry.UpdateDeviceStatus(deviceID, status)
}
```

## 11.8.1.5 数据采集和处理

### 11.8.1.5.1 数据采集模型

**定义 11.8.1.7** (数据采集)
数据采集是从IoT设备获取原始数据的过程。

**定理 11.8.1.3** (数据采集效率)
数据采集的效率取决于采样频率和传输带宽的平衡。

**证明**:
设采样频率为 $f$，传输带宽为 $B$，数据大小为 $S$。
则效率为：
$$E = \min\left(\frac{B}{S \cdot f}, 1\right)$$

### 11.8.1.5.2 Go实现数据处理

```go
// 数据采集器
type DataCollector struct {
    sensors    map[string]Sensor
    interval   time.Duration
    dataChan   chan SensorData
    isRunning  bool
    wg         sync.WaitGroup
}

// 传感器数据
type SensorData struct {
    SensorID   string
    Timestamp  time.Time
    Value      interface{}
    Unit       string
    Quality    float64
}

// 创建数据采集器
func NewDataCollector(interval time.Duration) *DataCollector {
    return &DataCollector{
        sensors:  make(map[string]Sensor),
        interval: interval,
        dataChan: make(chan SensorData, 1000),
    }
}

// 添加传感器
func (dc *DataCollector) AddSensor(sensor Sensor) {
    dc.sensors[sensor.GetID()] = sensor
}

// 开始采集
func (dc *DataCollector) Start() {
    if dc.isRunning {
        return
    }
    
    dc.isRunning = true
    dc.wg.Add(1)
    
    go func() {
        defer dc.wg.Done()
        ticker := time.NewTicker(dc.interval)
        defer ticker.Stop()
        
        for dc.isRunning {
            select {
            case <-ticker.C:
                dc.collectData()
            }
        }
    }()
}

// 停止采集
func (dc *DataCollector) Stop() {
    dc.isRunning = false
    dc.wg.Wait()
    close(dc.dataChan)
}

// 采集数据
func (dc *DataCollector) collectData() {
    for id, sensor := range dc.sensors {
        if value, err := sensor.Read(); err == nil {
            data := SensorData{
                SensorID:  id,
                Timestamp: time.Now(),
                Value:     value,
                Quality:   1.0,
            }
            
            select {
            case dc.dataChan <- data:
            default:
                // 通道满，丢弃数据
            }
        }
    }
}

// 获取数据通道
func (dc *DataCollector) GetDataChannel() <-chan SensorData {
    return dc.dataChan
}
```

## 11.8.1.6 总结

本章详细介绍了物联网基础理论，包括：

1. **IoT架构**: 三层架构模型和Go实现
2. **协议栈**: MQTT协议的理论和实现
3. **设备管理**: 设备注册和状态管理
4. **数据处理**: 数据采集和处理管道

通过Go语言实现，展示了IoT系统的核心组件和功能，为构建完整的IoT解决方案提供了理论基础和实践指导。

---

**相关链接**:
- [11.8.2 IoT安全](../02-IoT-Security/README.md)
- [11.8.3 IoT边缘计算](../03-IoT-Edge-Computing/README.md)
- [11.8.4 IoT应用](../04-IoT-Applications/README.md)
- [11.9 其他高级主题](../README.md) 