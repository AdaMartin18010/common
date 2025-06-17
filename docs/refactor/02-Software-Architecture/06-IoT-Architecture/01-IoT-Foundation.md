# 1. IoT架构基础

## 1.1 IoT理论基础

### 1.1.1 IoT系统形式化定义

**定义 1.1** (IoT系统): IoT系统是一个五元组 $\mathcal{I} = (D, G, C, P, A)$，其中：

- $D$ 是设备集合 (Devices)
- $G$ 是网关集合 (Gateways)
- $C$ 是云平台集合 (Cloud)
- $P$ 是协议集合 (Protocols)
- $A$ 是应用集合 (Applications)

**IoT架构层次**:

```latex
\text{IoTArchitecture} = \text{DeviceLayer} \times \text{GatewayLayer} \times \text{CloudLayer} \times \text{ApplicationLayer}
```

### 1.1.2 IoT设备模型

**定义 1.2** (IoT设备): IoT设备是一个四元组 $\text{Device} = (S, A, P, C)$，其中：

- $S$ 是传感器集合
- $A$ 是执行器集合
- $P$ 是处理单元
- $C$ 是通信模块

**设备能力模型**:

```latex
\text{DeviceCapability} = \text{Sensing} \cup \text{Actuation} \cup \text{Processing} \cup \text{Communication}
```

## 1.2 Go语言IoT设备实现

### 1.2.1 设备抽象层

```go
// Device IoT设备接口
type Device interface {
    // 设备标识
    GetID() string
    GetName() string
    GetType() DeviceType
    
    // 设备状态
    GetStatus() DeviceStatus
    Start() error
    Stop() error
    
    // 数据操作
    ReadSensor(sensorID string) (SensorData, error)
    WriteActuator(actuatorID string, value interface{}) error
    
    // 配置管理
    GetConfig() DeviceConfig
    UpdateConfig(config DeviceConfig) error
}

// DeviceType 设备类型
type DeviceType string

const (
    DeviceTypeSensor    DeviceType = "SENSOR"
    DeviceTypeActuator  DeviceType = "ACTUATOR"
    DeviceTypeGateway   DeviceType = "GATEWAY"
    DeviceTypeController DeviceType = "CONTROLLER"
)

// DeviceStatus 设备状态
type DeviceStatus struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        DeviceType `json:"type"`
    Status      string    `json:"status"`
    LastSeen    time.Time `json:"last_seen"`
    Battery     float64   `json:"battery"`
    Signal      float64   `json:"signal"`
    Temperature float64   `json:"temperature"`
    Errors      []string  `json:"errors"`
}

// DeviceConfig 设备配置
type DeviceConfig struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Type         DeviceType             `json:"type"`
    Location     Location               `json:"location"`
    Parameters   map[string]interface{} `json:"parameters"`
    Sensors      []SensorConfig         `json:"sensors"`
    Actuators    []ActuatorConfig       `json:"actuators"`
    Communication CommunicationConfig   `json:"communication"`
}

// Location 位置信息
type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Altitude  float64 `json:"altitude"`
    Address   string  `json:"address"`
}

// BaseDevice 基础设备实现
type BaseDevice struct {
    id       string
    name     string
    deviceType DeviceType
    config   DeviceConfig
    status   DeviceStatus
    sensors  map[string]Sensor
    actuators map[string]Actuator
    comm     CommunicationModule
    mutex    sync.RWMutex
    stopChan chan struct{}
}

// NewBaseDevice 创建基础设备
func NewBaseDevice(id, name string, deviceType DeviceType) *BaseDevice {
    return &BaseDevice{
        id:         id,
        name:       name,
        deviceType: deviceType,
        sensors:    make(map[string]Sensor),
        actuators:  make(map[string]Actuator),
        stopChan:   make(chan struct{}),
    }
}

// GetID 获取设备ID
func (bd *BaseDevice) GetID() string {
    return bd.id
}

// GetName 获取设备名称
func (bd *BaseDevice) GetName() string {
    return bd.name
}

// GetType 获取设备类型
func (bd *BaseDevice) GetType() DeviceType {
    return bd.deviceType
}

// GetStatus 获取设备状态
func (bd *BaseDevice) GetStatus() DeviceStatus {
    bd.mutex.RLock()
    defer bd.mutex.RUnlock()
    
    status := bd.status
    status.LastSeen = time.Now()
    return status
}

// Start 启动设备
func (bd *BaseDevice) Start() error {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.status.Status = "RUNNING"
    
    // 启动通信模块
    if bd.comm != nil {
        if err := bd.comm.Start(); err != nil {
            return fmt.Errorf("failed to start communication: %w", err)
        }
    }
    
    // 启动传感器
    for _, sensor := range bd.sensors {
        if err := sensor.Start(); err != nil {
            return fmt.Errorf("failed to start sensor: %w", err)
        }
    }
    
    // 启动执行器
    for _, actuator := range bd.actuators {
        if err := actuator.Start(); err != nil {
            return fmt.Errorf("failed to start actuator: %w", err)
        }
    }
    
    return nil
}

// Stop 停止设备
func (bd *BaseDevice) Stop() error {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.status.Status = "STOPPED"
    close(bd.stopChan)
    
    // 停止通信模块
    if bd.comm != nil {
        if err := bd.comm.Stop(); err != nil {
            return fmt.Errorf("failed to stop communication: %w", err)
        }
    }
    
    // 停止传感器
    for _, sensor := range bd.sensors {
        if err := sensor.Stop(); err != nil {
            return fmt.Errorf("failed to stop sensor: %w", err)
        }
    }
    
    // 停止执行器
    for _, actuator := range bd.actuators {
        if err := actuator.Stop(); err != nil {
            return fmt.Errorf("failed to stop actuator: %w", err)
        }
    }
    
    return nil
}

// ReadSensor 读取传感器数据
func (bd *BaseDevice) ReadSensor(sensorID string) (SensorData, error) {
    bd.mutex.RLock()
    defer bd.mutex.RUnlock()
    
    sensor, exists := bd.sensors[sensorID]
    if !exists {
        return SensorData{}, fmt.Errorf("sensor not found: %s", sensorID)
    }
    
    return sensor.Read()
}

// WriteActuator 写入执行器
func (bd *BaseDevice) WriteActuator(actuatorID string, value interface{}) error {
    bd.mutex.RLock()
    defer bd.mutex.RUnlock()
    
    actuator, exists := bd.actuators[actuatorID]
    if !exists {
        return fmt.Errorf("actuator not found: %s", actuatorID)
    }
    
    return actuator.Write(value)
}

// GetConfig 获取设备配置
func (bd *BaseDevice) GetConfig() DeviceConfig {
    bd.mutex.RLock()
    defer bd.mutex.RUnlock()
    
    return bd.config
}

// UpdateConfig 更新设备配置
func (bd *BaseDevice) UpdateConfig(config DeviceConfig) error {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.config = config
    return nil
}

// AddSensor 添加传感器
func (bd *BaseDevice) AddSensor(sensor Sensor) {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.sensors[sensor.GetID()] = sensor
}

// AddActuator 添加执行器
func (bd *BaseDevice) AddActuator(actuator Actuator) {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.actuators[actuator.GetID()] = actuator
}

// SetCommunication 设置通信模块
func (bd *BaseDevice) SetCommunication(comm CommunicationModule) {
    bd.mutex.Lock()
    defer bd.mutex.Unlock()
    
    bd.comm = comm
}
```

### 1.2.2 传感器实现

```go
// Sensor 传感器接口
type Sensor interface {
    // 传感器标识
    GetID() string
    GetName() string
    GetType() SensorType
    
    // 传感器操作
    Start() error
    Stop() error
    Read() (SensorData, error)
    
    // 配置管理
    GetConfig() SensorConfig
    UpdateConfig(config SensorConfig) error
}

// SensorType 传感器类型
type SensorType string

const (
    SensorTypeTemperature SensorType = "TEMPERATURE"
    SensorTypeHumidity    SensorType = "HUMIDITY"
    SensorTypePressure    SensorType = "PRESSURE"
    SensorTypeLight       SensorType = "LIGHT"
    SensorTypeMotion      SensorType = "MOTION"
    SensorTypeGas         SensorType = "GAS"
    SensorTypeSound       SensorType = "SOUND"
)

// SensorData 传感器数据
type SensorData struct {
    ID        string                 `json:"id"`
    Type      SensorType             `json:"type"`
    Value     interface{}            `json:"value"`
    Unit      string                 `json:"unit"`
    Timestamp time.Time              `json:"timestamp"`
    Quality   DataQuality            `json:"quality"`
    Metadata  map[string]interface{} `json:"metadata"`
}

// DataQuality 数据质量
type DataQuality string

const (
    DataQualityGood    DataQuality = "GOOD"
    DataQualityBad     DataQuality = "BAD"
    DataQualityUnknown DataQuality = "UNKNOWN"
)

// SensorConfig 传感器配置
type SensorConfig struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        SensorType             `json:"type"`
    Unit        string                 `json:"unit"`
    Range       SensorRange            `json:"range"`
    Calibration map[string]interface{} `json:"calibration"`
    Parameters  map[string]interface{} `json:"parameters"`
}

// SensorRange 传感器量程
type SensorRange struct {
    Min float64 `json:"min"`
    Max float64 `json:"max"`
}

// BaseSensor 基础传感器实现
type BaseSensor struct {
    id       string
    name     string
    sensorType SensorType
    config   SensorConfig
    running  bool
    mutex    sync.RWMutex
    stopChan chan struct{}
}

// NewBaseSensor 创建基础传感器
func NewBaseSensor(id, name string, sensorType SensorType) *BaseSensor {
    return &BaseSensor{
        id:         id,
        name:       name,
        sensorType: sensorType,
        stopChan:   make(chan struct{}),
    }
}

// GetID 获取传感器ID
func (bs *BaseSensor) GetID() string {
    return bs.id
}

// GetName 获取传感器名称
func (bs *BaseSensor) GetName() string {
    return bs.name
}

// GetType 获取传感器类型
func (bs *BaseSensor) GetType() SensorType {
    return bs.sensorType
}

// Start 启动传感器
func (bs *BaseSensor) Start() error {
    bs.mutex.Lock()
    defer bs.mutex.Unlock()
    
    bs.running = true
    return nil
}

// Stop 停止传感器
func (bs *BaseSensor) Stop() error {
    bs.mutex.Lock()
    defer bs.mutex.Unlock()
    
    bs.running = false
    close(bs.stopChan)
    return nil
}

// GetConfig 获取传感器配置
func (bs *BaseSensor) GetConfig() SensorConfig {
    bs.mutex.RLock()
    defer bs.mutex.RUnlock()
    
    return bs.config
}

// UpdateConfig 更新传感器配置
func (bs *BaseSensor) UpdateConfig(config SensorConfig) error {
    bs.mutex.Lock()
    defer bs.mutex.Unlock()
    
    bs.config = config
    return nil
}

// TemperatureSensor 温度传感器
type TemperatureSensor struct {
    *BaseSensor
    temperature float64
}

// NewTemperatureSensor 创建温度传感器
func NewTemperatureSensor(id, name string) *TemperatureSensor {
    return &TemperatureSensor{
        BaseSensor: NewBaseSensor(id, name, SensorTypeTemperature),
    }
}

// Read 读取温度数据
func (ts *TemperatureSensor) Read() (SensorData, error) {
    ts.mutex.RLock()
    defer ts.mutex.RUnlock()
    
    if !ts.running {
        return SensorData{}, fmt.Errorf("sensor not running")
    }
    
    // 模拟温度读取
    ts.temperature = 20.0 + rand.Float64()*10.0
    
    return SensorData{
        ID:        ts.id,
        Type:      ts.sensorType,
        Value:     ts.temperature,
        Unit:      "°C",
        Timestamp: time.Now(),
        Quality:   DataQualityGood,
        Metadata: map[string]interface{}{
            "sensor_model": "DHT22",
            "accuracy":     "±0.5°C",
        },
    }, nil
}

// HumiditySensor 湿度传感器
type HumiditySensor struct {
    *BaseSensor
    humidity float64
}

// NewHumiditySensor 创建湿度传感器
func NewHumiditySensor(id, name string) *HumiditySensor {
    return &HumiditySensor{
        BaseSensor: NewBaseSensor(id, name, SensorTypeHumidity),
    }
}

// Read 读取湿度数据
func (hs *HumiditySensor) Read() (SensorData, error) {
    hs.mutex.RLock()
    defer hs.mutex.RUnlock()
    
    if !hs.running {
        return SensorData{}, fmt.Errorf("sensor not running")
    }
    
    // 模拟湿度读取
    hs.humidity = 40.0 + rand.Float64()*30.0
    
    return SensorData{
        ID:        hs.id,
        Type:      hs.sensorType,
        Value:     hs.humidity,
        Unit:      "%",
        Timestamp: time.Now(),
        Quality:   DataQualityGood,
        Metadata: map[string]interface{}{
            "sensor_model": "DHT22",
            "accuracy":     "±2%",
        },
    }, nil
}
```

### 1.2.3 执行器实现

```go
// Actuator 执行器接口
type Actuator interface {
    // 执行器标识
    GetID() string
    GetName() string
    GetType() ActuatorType
    
    // 执行器操作
    Start() error
    Stop() error
    Write(value interface{}) error
    Read() (ActuatorState, error)
    
    // 配置管理
    GetConfig() ActuatorConfig
    UpdateConfig(config ActuatorConfig) error
}

// ActuatorType 执行器类型
type ActuatorType string

const (
    ActuatorTypeRelay    ActuatorType = "RELAY"
    ActuatorTypeMotor    ActuatorType = "MOTOR"
    ActuatorTypeValve    ActuatorType = "VALVE"
    ActuatorTypeLED      ActuatorType = "LED"
    ActuatorTypeServo    ActuatorType = "SERVO"
    ActuatorTypePump     ActuatorType = "PUMP"
)

// ActuatorState 执行器状态
type ActuatorState struct {
    ID        string                 `json:"id"`
    Type      ActuatorType           `json:"type"`
    Value     interface{}            `json:"value"`
    Status    string                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}

// ActuatorConfig 执行器配置
type ActuatorConfig struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        ActuatorType           `json:"type"`
    Range       ActuatorRange          `json:"range"`
    Parameters  map[string]interface{} `json:"parameters"`
    Safety      SafetyConfig           `json:"safety"`
}

// ActuatorRange 执行器范围
type ActuatorRange struct {
    Min interface{} `json:"min"`
    Max interface{} `json:"max"`
}

// SafetyConfig 安全配置
type SafetyConfig struct {
    MaxOnTime    time.Duration `json:"max_on_time"`
    MinOffTime   time.Duration `json:"min_off_time"`
    MaxCycles    int           `json:"max_cycles"`
    Temperature  float64       `json:"temperature"`
}

// BaseActuator 基础执行器实现
type BaseActuator struct {
    id         string
    name       string
    actuatorType ActuatorType
    config     ActuatorConfig
    state      ActuatorState
    running    bool
    mutex      sync.RWMutex
    stopChan   chan struct{}
}

// NewBaseActuator 创建基础执行器
func NewBaseActuator(id, name string, actuatorType ActuatorType) *BaseActuator {
    return &BaseActuator{
        id:           id,
        name:         name,
        actuatorType: actuatorType,
        stopChan:     make(chan struct{}),
    }
}

// GetID 获取执行器ID
func (ba *BaseActuator) GetID() string {
    return ba.id
}

// GetName 获取执行器名称
func (ba *BaseActuator) GetName() string {
    return ba.name
}

// GetType 获取执行器类型
func (ba *BaseActuator) GetType() ActuatorType {
    return ba.actuatorType
}

// Start 启动执行器
func (ba *BaseActuator) Start() error {
    ba.mutex.Lock()
    defer ba.mutex.Unlock()
    
    ba.running = true
    ba.state.Status = "READY"
    return nil
}

// Stop 停止执行器
func (ba *BaseActuator) Stop() error {
    ba.mutex.Lock()
    defer ba.mutex.Unlock()
    
    ba.running = false
    ba.state.Status = "STOPPED"
    close(ba.stopChan)
    return nil
}

// Read 读取执行器状态
func (ba *BaseActuator) Read() (ActuatorState, error) {
    ba.mutex.RLock()
    defer ba.mutex.RUnlock()
    
    if !ba.running {
        return ActuatorState{}, fmt.Errorf("actuator not running")
    }
    
    return ba.state, nil
}

// GetConfig 获取执行器配置
func (ba *BaseActuator) GetConfig() ActuatorConfig {
    ba.mutex.RLock()
    defer ba.mutex.RUnlock()
    
    return ba.config
}

// UpdateConfig 更新执行器配置
func (ba *BaseActuator) UpdateConfig(config ActuatorConfig) error {
    ba.mutex.Lock()
    defer ba.mutex.Unlock()
    
    ba.config = config
    return nil
}

// RelayActuator 继电器执行器
type RelayActuator struct {
    *BaseActuator
    isOn bool
}

// NewRelayActuator 创建继电器执行器
func NewRelayActuator(id, name string) *RelayActuator {
    return &RelayActuator{
        BaseActuator: NewBaseActuator(id, name, ActuatorTypeRelay),
    }
}

// Write 写入继电器状态
func (ra *RelayActuator) Write(value interface{}) error {
    ra.mutex.Lock()
    defer ra.mutex.Unlock()
    
    if !ra.running {
        return fmt.Errorf("actuator not running")
    }
    
    // 解析布尔值
    var isOn bool
    switch v := value.(type) {
    case bool:
        isOn = v
    case string:
        isOn = strings.ToLower(v) == "true" || v == "1" || v == "on"
    case int:
        isOn = v != 0
    default:
        return fmt.Errorf("invalid value type for relay: %T", value)
    }
    
    // 更新状态
    ra.isOn = isOn
    ra.state.Value = isOn
    ra.state.Status = "ACTIVE"
    ra.state.Timestamp = time.Now()
    
    // 模拟硬件控制
    log.Printf("Relay %s set to: %v", ra.id, isOn)
    
    return nil
}

// Read 读取继电器状态
func (ra *RelayActuator) Read() (ActuatorState, error) {
    ra.mutex.RLock()
    defer ra.mutex.RUnlock()
    
    if !ra.running {
        return ActuatorState{}, fmt.Errorf("actuator not running")
    }
    
    state := ra.state
    state.Value = ra.isOn
    state.Timestamp = time.Now()
    
    return state, nil
}

// LEDActuator LED执行器
type LEDActuator struct {
    *BaseActuator
    brightness float64
    color      string
}

// NewLEDActuator 创建LED执行器
func NewLEDActuator(id, name string) *LEDActuator {
    return &LEDActuator{
        BaseActuator: NewBaseActuator(id, name, ActuatorTypeLED),
        brightness:   0.0,
        color:        "white",
    }
}

// Write 写入LED状态
func (la *LEDActuator) Write(value interface{}) error {
    la.mutex.Lock()
    defer la.mutex.Unlock()
    
    if !la.running {
        return fmt.Errorf("actuator not running")
    }
    
    // 解析LED控制值
    switch v := value.(type) {
    case map[string]interface{}:
        if brightness, ok := v["brightness"].(float64); ok {
            la.brightness = math.Max(0.0, math.Min(1.0, brightness))
        }
        if color, ok := v["color"].(string); ok {
            la.color = color
        }
    case float64:
        la.brightness = math.Max(0.0, math.Min(1.0, v))
    case int:
        la.brightness = math.Max(0.0, math.Min(1.0, float64(v)/100.0))
    default:
        return fmt.Errorf("invalid value type for LED: %T", value)
    }
    
    // 更新状态
    la.state.Value = map[string]interface{}{
        "brightness": la.brightness,
        "color":      la.color,
    }
    la.state.Status = "ACTIVE"
    la.state.Timestamp = time.Now()
    
    // 模拟硬件控制
    log.Printf("LED %s set to brightness: %.2f, color: %s", la.id, la.brightness, la.color)
    
    return nil
}
```

## 1.3 通信协议实现

### 1.3.1 通信模块抽象

```go
// CommunicationModule 通信模块接口
type CommunicationModule interface {
    // 通信操作
    Start() error
    Stop() error
    Send(data []byte) error
    Receive() (<-chan Message, error)
    
    // 配置管理
    GetConfig() CommunicationConfig
    UpdateConfig(config CommunicationConfig) error
    
    // 状态查询
    GetStatus() CommunicationStatus
}

// CommunicationConfig 通信配置
type CommunicationConfig struct {
    Protocol    ProtocolType           `json:"protocol"`
    Address     string                 `json:"address"`
    Port        int                    `json:"port"`
    Parameters  map[string]interface{} `json:"parameters"`
    Security    SecurityConfig         `json:"security"`
}

// ProtocolType 协议类型
type ProtocolType string

const (
    ProtocolTypeMQTT    ProtocolType = "MQTT"
    ProtocolTypeHTTP    ProtocolType = "HTTP"
    ProtocolTypeCoAP    ProtocolType = "COAP"
    ProtocolTypeModbus  ProtocolType = "MODBUS"
    ProtocolTypeZigbee  ProtocolType = "ZIGBEE"
    ProtocolTypeLoRa    ProtocolType = "LORA"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
    Encryption  string `json:"encryption"`
    Key         string `json:"key"`
    Certificate string `json:"certificate"`
    Username    string `json:"username"`
    Password    string `json:"password"`
}

// CommunicationStatus 通信状态
type CommunicationStatus struct {
    Connected   bool      `json:"connected"`
    LastSeen    time.Time `json:"last_seen"`
    BytesSent   int64     `json:"bytes_sent"`
    BytesReceived int64   `json:"bytes_received"`
    Errors      []string  `json:"errors"`
}

// Message 消息结构
type Message struct {
    ID        string                 `json:"id"`
    Topic     string                 `json:"topic"`
    Payload   []byte                 `json:"payload"`
    QoS       int                    `json:"qos"`
    Retain    bool                   `json:"retain"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}

// MQTTCommunication MQTT通信模块
type MQTTCommunication struct {
    config     CommunicationConfig
    client     mqtt.Client
    status     CommunicationStatus
    messageChan chan Message
    mutex      sync.RWMutex
    stopChan   chan struct{}
}

// NewMQTTCommunication 创建MQTT通信模块
func NewMQTTCommunication(config CommunicationConfig) *MQTTCommunication {
    return &MQTTCommunication{
        config:      config,
        messageChan: make(chan Message, 100),
        stopChan:    make(chan struct{}),
    }
}

// Start 启动MQTT通信
func (mc *MQTTCommunication) Start() error {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    // 创建MQTT客户端选项
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mc.config.Address, mc.config.Port))
    opts.SetClientID(fmt.Sprintf("iot_device_%s", uuid.New().String()))
    
    // 设置安全配置
    if mc.config.Security.Username != "" {
        opts.SetUsername(mc.config.Security.Username)
        opts.SetPassword(mc.config.Security.Password)
    }
    
    // 设置连接回调
    opts.SetOnConnectHandler(func(client mqtt.Client) {
        mc.status.Connected = true
        mc.status.LastSeen = time.Now()
        log.Printf("MQTT connected to %s:%d", mc.config.Address, mc.config.Port)
    })
    
    // 设置连接丢失回调
    opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
        mc.status.Connected = false
        mc.status.Errors = append(mc.status.Errors, err.Error())
        log.Printf("MQTT connection lost: %v", err)
    })
    
    // 创建客户端
    mc.client = mqtt.NewClient(opts)
    
    // 连接到服务器
    if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
        return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
    }
    
    return nil
}

// Stop 停止MQTT通信
func (mc *MQTTCommunication) Stop() error {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    if mc.client != nil && mc.client.IsConnected() {
        mc.client.Disconnect(250)
    }
    
    mc.status.Connected = false
    close(mc.stopChan)
    
    return nil
}

// Send 发送消息
func (mc *MQTTCommunication) Send(data []byte) error {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    if !mc.status.Connected {
        return fmt.Errorf("not connected to MQTT broker")
    }
    
    // 创建消息
    message := Message{
        ID:        uuid.New().String(),
        Topic:     "iot/data",
        Payload:   data,
        QoS:       1,
        Retain:    false,
        Timestamp: time.Now(),
    }
    
    // 发布消息
    token := mc.client.Publish(message.Topic, byte(message.QoS), message.Retain, message.Payload)
    if token.Wait() && token.Error() != nil {
        return fmt.Errorf("failed to publish message: %w", token.Error())
    }
    
    mc.status.BytesSent += int64(len(data))
    mc.status.LastSeen = time.Now()
    
    return nil
}

// Receive 接收消息
func (mc *MQTTCommunication) Receive() (<-chan Message, error) {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    if !mc.status.Connected {
        return nil, fmt.Errorf("not connected to MQTT broker")
    }
    
    // 订阅主题
    token := mc.client.Subscribe("iot/command", 1, func(client mqtt.Client, msg mqtt.Message) {
        message := Message{
            ID:        uuid.New().String(),
            Topic:     msg.Topic(),
            Payload:   msg.Payload(),
            QoS:       int(msg.Qos()),
            Retain:    msg.Retained(),
            Timestamp: time.Now(),
        }
        
        select {
        case mc.messageChan <- message:
        default:
            log.Printf("Message channel full, dropping message")
        }
        
        mc.status.BytesReceived += int64(len(msg.Payload()))
        mc.status.LastSeen = time.Now()
    })
    
    if token.Wait() && token.Error() != nil {
        return nil, fmt.Errorf("failed to subscribe: %w", token.Error())
    }
    
    return mc.messageChan, nil
}

// GetConfig 获取通信配置
func (mc *MQTTCommunication) GetConfig() CommunicationConfig {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    return mc.config
}

// UpdateConfig 更新通信配置
func (mc *MQTTCommunication) UpdateConfig(config CommunicationConfig) error {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    mc.config = config
    return nil
}

// GetStatus 获取通信状态
func (mc *MQTTCommunication) GetStatus() CommunicationStatus {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    return mc.status
}
```

## 1.4 IoT网关实现

### 1.4.1 网关核心功能

```go
// Gateway IoT网关接口
type Gateway interface {
    // 网关操作
    Start() error
    Stop() error
    
    // 设备管理
    RegisterDevice(device Device) error
    UnregisterDevice(deviceID string) error
    GetDevices() map[string]Device
    
    // 数据处理
    ProcessData(data DeviceData) error
    ForwardData(data DeviceData) error
    
    // 配置管理
    GetConfig() GatewayConfig
    UpdateConfig(config GatewayConfig) error
}

// GatewayConfig 网关配置
type GatewayConfig struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Location    Location               `json:"location"`
    Protocols   []ProtocolConfig       `json:"protocols"`
    Processing  ProcessingConfig       `json:"processing"`
    Storage     StorageConfig          `json:"storage"`
    Security    SecurityConfig         `json:"security"`
}

// ProtocolConfig 协议配置
type ProtocolConfig struct {
    Type       ProtocolType           `json:"type"`
    Enabled    bool                   `json:"enabled"`
    Parameters map[string]interface{} `json:"parameters"`
}

// ProcessingConfig 处理配置
type ProcessingConfig struct {
    EnableFiltering    bool `json:"enable_filtering"`
    EnableAggregation  bool `json:"enable_aggregation"`
    EnableTransformation bool `json:"enable_transformation"`
    BatchSize          int  `json:"batch_size"`
    BatchTimeout       time.Duration `json:"batch_timeout"`
}

// StorageConfig 存储配置
type StorageConfig struct {
    Type       string                 `json:"type"`
    Parameters map[string]interface{} `json:"parameters"`
}

// DeviceData 设备数据
type DeviceData struct {
    DeviceID   string                 `json:"device_id"`
    SensorID   string                 `json:"sensor_id"`
    Value      interface{}            `json:"value"`
    Unit       string                 `json:"unit"`
    Timestamp  time.Time              `json:"timestamp"`
    Quality    DataQuality            `json:"quality"`
    Metadata   map[string]interface{} `json:"metadata"`
}

// IoTGateway IoT网关实现
type IoTGateway struct {
    config      GatewayConfig
    devices     map[string]Device
    protocols   map[ProtocolType]CommunicationModule
    processor   DataProcessor
    storage     DataStorage
    mutex       sync.RWMutex
    stopChan    chan struct{}
}

// NewIoTGateway 创建IoT网关
func NewIoTGateway(config GatewayConfig) *IoTGateway {
    return &IoTGateway{
        config:    config,
        devices:   make(map[string]Device),
        protocols: make(map[ProtocolType]CommunicationModule),
        stopChan:  make(chan struct{}),
    }
}

// Start 启动网关
func (ig *IoTGateway) Start() error {
    ig.mutex.Lock()
    defer ig.mutex.Unlock()
    
    // 启动协议模块
    for _, protocolConfig := range ig.config.Protocols {
        if protocolConfig.Enabled {
            if err := ig.startProtocol(protocolConfig); err != nil {
                return fmt.Errorf("failed to start protocol %s: %w", protocolConfig.Type, err)
            }
        }
    }
    
    // 启动数据处理
    if ig.processor != nil {
        if err := ig.processor.Start(); err != nil {
            return fmt.Errorf("failed to start data processor: %w", err)
        }
    }
    
    // 启动数据存储
    if ig.storage != nil {
        if err := ig.storage.Start(); err != nil {
            return fmt.Errorf("failed to start data storage: %w", err)
        }
    }
    
    // 启动设备监控
    go ig.monitorDevices()
    
    return nil
}

// Stop 停止网关
func (ig *IoTGateway) Stop() error {
    ig.mutex.Lock()
    defer ig.mutex.Unlock()
    
    close(ig.stopChan)
    
    // 停止协议模块
    for _, protocol := range ig.protocols {
        if err := protocol.Stop(); err != nil {
            log.Printf("Failed to stop protocol: %v", err)
        }
    }
    
    // 停止数据处理
    if ig.processor != nil {
        if err := ig.processor.Stop(); err != nil {
            log.Printf("Failed to stop data processor: %v", err)
        }
    }
    
    // 停止数据存储
    if ig.storage != nil {
        if err := ig.storage.Stop(); err != nil {
            log.Printf("Failed to stop data storage: %v", err)
        }
    }
    
    return nil
}

// RegisterDevice 注册设备
func (ig *IoTGateway) RegisterDevice(device Device) error {
    ig.mutex.Lock()
    defer ig.mutex.Unlock()
    
    ig.devices[device.GetID()] = device
    
    // 启动设备
    if err := device.Start(); err != nil {
        return fmt.Errorf("failed to start device: %w", err)
    }
    
    log.Printf("Device %s registered and started", device.GetID())
    return nil
}

// UnregisterDevice 注销设备
func (ig *IoTGateway) UnregisterDevice(deviceID string) error {
    ig.mutex.Lock()
    defer ig.mutex.Unlock()
    
    device, exists := ig.devices[deviceID]
    if !exists {
        return fmt.Errorf("device not found: %s", deviceID)
    }
    
    // 停止设备
    if err := device.Stop(); err != nil {
        return fmt.Errorf("failed to stop device: %w", err)
    }
    
    delete(ig.devices, deviceID)
    log.Printf("Device %s unregistered and stopped", deviceID)
    
    return nil
}

// GetDevices 获取所有设备
func (ig *IoTGateway) GetDevices() map[string]Device {
    ig.mutex.RLock()
    defer ig.mutex.RUnlock()
    
    result := make(map[string]Device)
    for k, v := range ig.devices {
        result[k] = v
    }
    return result
}

// ProcessData 处理数据
func (ig *IoTGateway) ProcessData(data DeviceData) error {
    // 数据过滤
    if ig.config.Processing.EnableFiltering {
        if !ig.filterData(data) {
            return nil // 数据被过滤掉
        }
    }
    
    // 数据转换
    if ig.config.Processing.EnableTransformation {
        data = ig.transformData(data)
    }
    
    // 数据聚合
    if ig.config.Processing.EnableAggregation {
        if err := ig.aggregateData(data); err != nil {
            return fmt.Errorf("failed to aggregate data: %w", err)
        }
    }
    
    // 存储数据
    if ig.storage != nil {
        if err := ig.storage.Store(data); err != nil {
            return fmt.Errorf("failed to store data: %w", err)
        }
    }
    
    // 转发数据
    if err := ig.ForwardData(data); err != nil {
        return fmt.Errorf("failed to forward data: %w", err)
    }
    
    return nil
}

// ForwardData 转发数据
func (ig *IoTGateway) ForwardData(data DeviceData) error {
    // 序列化数据
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }
    
    // 通过所有启用的协议转发
    for protocolType, protocol := range ig.protocols {
        if err := protocol.Send(jsonData); err != nil {
            log.Printf("Failed to forward data via %s: %v", protocolType, err)
        }
    }
    
    return nil
}

// GetConfig 获取网关配置
func (ig *IoTGateway) GetConfig() GatewayConfig {
    ig.mutex.RLock()
    defer ig.mutex.RUnlock()
    
    return ig.config
}

// UpdateConfig 更新网关配置
func (ig *IoTGateway) UpdateConfig(config GatewayConfig) error {
    ig.mutex.Lock()
    defer ig.mutex.Unlock()
    
    ig.config = config
    return nil
}

// startProtocol 启动协议模块
func (ig *IoTGateway) startProtocol(config ProtocolConfig) error {
    var protocol CommunicationModule
    
    switch config.Type {
    case ProtocolTypeMQTT:
        commConfig := CommunicationConfig{
            Protocol:   ProtocolTypeMQTT,
            Address:    config.Parameters["address"].(string),
            Port:       int(config.Parameters["port"].(float64)),
            Parameters: config.Parameters,
        }
        protocol = NewMQTTCommunication(commConfig)
    default:
        return fmt.Errorf("unsupported protocol: %s", config.Type)
    }
    
    if err := protocol.Start(); err != nil {
        return err
    }
    
    ig.protocols[config.Type] = protocol
    return nil
}

// monitorDevices 监控设备
func (ig *IoTGateway) monitorDevices() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ig.stopChan:
            return
        case <-ticker.C:
            ig.checkDeviceHealth()
        }
    }
}

// checkDeviceHealth 检查设备健康状态
func (ig *IoTGateway) checkDeviceHealth() {
    ig.mutex.RLock()
    devices := make(map[string]Device)
    for k, v := range ig.devices {
        devices[k] = v
    }
    ig.mutex.RUnlock()
    
    for deviceID, device := range devices {
        status := device.GetStatus()
        
        // 检查设备是否响应
        if time.Since(status.LastSeen) > 5*time.Minute {
            log.Printf("Device %s appears to be unhealthy", deviceID)
            // 可以在这里实现设备重启逻辑
        }
    }
}

// filterData 过滤数据
func (ig *IoTGateway) filterData(data DeviceData) bool {
    // 简单的数据质量过滤
    if data.Quality == DataQualityBad {
        return false
    }
    
    // 可以添加更多过滤规则
    return true
}

// transformData 转换数据
func (ig *IoTGateway) transformData(data DeviceData) DeviceData {
    // 添加网关元数据
    data.Metadata["gateway_id"] = ig.config.ID
    data.Metadata["gateway_name"] = ig.config.Name
    data.Metadata["processed_at"] = time.Now()
    
    return data
}

// aggregateData 聚合数据
func (ig *IoTGateway) aggregateData(data DeviceData) error {
    // 实现数据聚合逻辑
    // 这里可以按时间窗口、设备组等进行聚合
    return nil
}
```

## 1.5 总结

IoT架构基础模块涵盖了以下核心内容：

1. **理论基础**: 形式化定义IoT系统的数学模型
2. **设备抽象**: Go语言实现的设备、传感器、执行器抽象层
3. **通信协议**: 支持多种IoT通信协议（MQTT、HTTP、CoAP等）
4. **网关功能**: 设备管理、数据处理、协议转换等核心功能
5. **配置管理**: 灵活的配置管理和动态更新机制

这个设计提供了一个完整的IoT架构框架，支持各种IoT设备的接入、数据处理和云端通信。 