# Go语言在物联网中的应用 (Go Language in IoT)

## 概述

Go语言在物联网(IoT)领域凭借其高性能、低资源消耗、并发处理能力和跨平台特性，成为构建IoT设备、网关、云平台和边缘计算系统的理想选择。从传感器数据采集到设备管理，从边缘计算到云端分析，Go语言为IoT生态系统提供了稳定、高效的技术基础。

## 基本概念

### 核心特征

- **高性能**：编译型语言，执行效率高，适合实时数据处理
- **低资源消耗**：内存占用小，适合资源受限的IoT设备
- **并发处理**：原生goroutine和channel，支持高并发数据流
- **跨平台**：支持多种硬件平台，便于IoT设备部署
- **网络编程**：强大的网络库，支持多种IoT通信协议
- **安全性**：内置安全特性，保护IoT设备和数据

### 应用场景

- **IoT设备**：传感器节点、执行器、网关设备
- **设备管理**：设备注册、配置、监控、固件更新
- **数据采集**：传感器数据收集、预处理、传输
- **边缘计算**：本地数据处理、实时分析、决策
- **云平台**：IoT云服务、数据分析、可视化
- **协议网关**：MQTT、CoAP、HTTP等协议转换

## 核心组件

### IoT设备管理平台 (IoT Device Management Platform)

```go
// 设备信息
type Device struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Model       string            `json:"model"`
    Version     string            `json:"version"`
    Status      DeviceStatus      `json:"status"`
    Location    Location          `json:"location"`
    Properties  map[string]interface{} `json:"properties"`
    Tags        []string          `json:"tags"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    LastSeen    time.Time         `json:"last_seen"`
}

// 设备状态
type DeviceStatus string

const (
    StatusOnline  DeviceStatus = "online"
    StatusOffline DeviceStatus = "offline"
    StatusError   DeviceStatus = "error"
    StatusMaintenance DeviceStatus = "maintenance"
)

// 位置信息
type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Altitude  float64 `json:"altitude"`
    Address   string  `json:"address"`
}

// 设备管理器
type DeviceManager struct {
    devices map[string]*Device
    mu      sync.RWMutex
}

func NewDeviceManager() *DeviceManager {
    return &DeviceManager{
        devices: make(map[string]*Device),
    }
}

func (dm *DeviceManager) RegisterDevice(device *Device) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    if device.ID == "" {
        device.ID = generateDeviceID()
    }
    
    device.CreatedAt = time.Now()
    device.UpdatedAt = time.Now()
    device.LastSeen = time.Now()
    device.Status = StatusOnline
    
    dm.devices[device.ID] = device
    return nil
}

func (dm *DeviceManager) GetDevice(deviceID string) (*Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    device, exists := dm.devices[deviceID]
    if !exists {
        return nil, fmt.Errorf("device not found: %s", deviceID)
    }
    
    return device, nil
}

func (dm *DeviceManager) UpdateDevice(device *Device) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    if _, exists := dm.devices[device.ID]; !exists {
        return fmt.Errorf("device not found: %s", device.ID)
    }
    
    device.UpdatedAt = time.Now()
    dm.devices[device.ID] = device
    return nil
}

func (dm *DeviceManager) DeleteDevice(deviceID string) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    delete(dm.devices, deviceID)
    return nil
}

func (dm *DeviceManager) ListDevices() []*Device {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    devices := make([]*Device, 0, len(dm.devices))
    for _, device := range dm.devices {
        devices = append(devices, device)
    }
    
    return devices
}

func (dm *DeviceManager) GetDevicesByType(deviceType string) []*Device {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    var devices []*Device
    for _, device := range dm.devices {
        if device.Type == deviceType {
            devices = append(devices, device)
        }
    }
    
    return devices
}

func (dm *DeviceManager) UpdateDeviceStatus(deviceID string, status DeviceStatus) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    device, exists := dm.devices[deviceID]
    if !exists {
        return fmt.Errorf("device not found: %s", deviceID)
    }
    
    device.Status = status
    device.UpdatedAt = time.Now()
    if status == StatusOnline {
        device.LastSeen = time.Now()
    }
    
    return nil
}

func (dm *DeviceManager) Heartbeat(deviceID string) error {
    return dm.UpdateDeviceStatus(deviceID, StatusOnline)
}

func generateDeviceID() string {
    return fmt.Sprintf("device_%d", time.Now().UnixNano())
}
```

### 传感器数据采集系统 (Sensor Data Collection System)

```go
// 传感器数据
type SensorData struct {
    ID        string                 `json:"id"`
    DeviceID  string                 `json:"device_id"`
    SensorID  string                 `json:"sensor_id"`
    Type      string                 `json:"type"`
    Value     float64                `json:"value"`
    Unit      string                 `json:"unit"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}

// 传感器配置
type SensorConfig struct {
    ID           string  `json:"id"`
    DeviceID     string  `json:"device_id"`
    Type         string  `json:"type"`
    Name         string  `json:"name"`
    Unit         string  `json:"unit"`
    MinValue     float64 `json:"min_value"`
    MaxValue     float64 `json:"max_value"`
    SamplingRate int     `json:"sampling_rate"` // 采样率(Hz)
    Enabled      bool    `json:"enabled"`
}

// 数据采集器
type DataCollector struct {
    deviceManager *DeviceManager
    sensors       map[string]*SensorConfig
    dataChan      chan *SensorData
    running       bool
    mu            sync.RWMutex
}

func NewDataCollector(deviceManager *DeviceManager) *DataCollector {
    return &DataCollector{
        deviceManager: deviceManager,
        sensors:       make(map[string]*SensorConfig),
        dataChan:      make(chan *SensorData, 1000),
        running:       false,
    }
}

func (dc *DataCollector) AddSensor(config *SensorConfig) error {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    dc.sensors[config.ID] = config
    return nil
}

func (dc *DataCollector) RemoveSensor(sensorID string) error {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    delete(dc.sensors, sensorID)
    return nil
}

func (dc *DataCollector) Start() {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    if dc.running {
        return
    }
    
    dc.running = true
    go dc.collectionLoop()
}

func (dc *DataCollector) Stop() {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    dc.running = false
    close(dc.dataChan)
}

func (dc *DataCollector) collectionLoop() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for dc.running {
        select {
        case <-ticker.C:
            dc.collectData()
        }
    }
}

func (dc *DataCollector) collectData() {
    dc.mu.RLock()
    sensors := make([]*SensorConfig, 0, len(dc.sensors))
    for _, sensor := range dc.sensors {
        if sensor.Enabled {
            sensors = append(sensors, sensor)
        }
    }
    dc.mu.RUnlock()
    
    for _, sensor := range sensors {
        go dc.collectSensorData(sensor)
    }
}

func (dc *DataCollector) collectSensorData(sensor *SensorConfig) {
    // 模拟传感器数据采集
    value := dc.simulateSensorReading(sensor)
    
    data := &SensorData{
        ID:        generateDataID(),
        DeviceID:  sensor.DeviceID,
        SensorID:  sensor.ID,
        Type:      sensor.Type,
        Value:     value,
        Unit:      sensor.Unit,
        Timestamp: time.Now(),
        Metadata: map[string]interface{}{
            "sensor_name": sensor.Name,
            "sampling_rate": sensor.SamplingRate,
        },
    }
    
    // 发送数据
    select {
    case dc.dataChan <- data:
    default:
        log.Printf("Data channel full, dropping sensor data")
    }
}

func (dc *DataCollector) simulateSensorReading(sensor *SensorConfig) float64 {
    // 模拟不同类型的传感器数据
    switch sensor.Type {
    case "temperature":
        return 20 + rand.Float64()*10 // 20-30°C
    case "humidity":
        return 40 + rand.Float64()*30 // 40-70%
    case "pressure":
        return 1000 + rand.Float64()*50 // 1000-1050 hPa
    case "light":
        return rand.Float64() * 1000 // 0-1000 lux
    default:
        return sensor.MinValue + rand.Float64()*(sensor.MaxValue-sensor.MinValue)
    }
}

func (dc *DataCollector) GetDataChannel() <-chan *SensorData {
    return dc.dataChan
}

func generateDataID() string {
    return fmt.Sprintf("data_%d", time.Now().UnixNano())
}
```

### MQTT通信网关 (MQTT Communication Gateway)

```go
// MQTT客户端配置
type MQTTConfig struct {
    Broker   string `json:"broker"`
    Port     int    `json:"port"`
    ClientID string `json:"client_id"`
    Username string `json:"username"`
    Password string `json:"password"`
    QoS      int    `json:"qos"`
}

// MQTT消息
type MQTTMessage struct {
    Topic   string                 `json:"topic"`
    Payload map[string]interface{} `json:"payload"`
    QoS     int                    `json:"qos"`
    Retain  bool                   `json:"retain"`
}

// MQTT网关
type MQTTGateway struct {
    config     *MQTTConfig
    client     mqtt.Client
    dataChan   chan *SensorData
    messageChan chan *MQTTMessage
    running    bool
    mu         sync.RWMutex
}

func NewMQTTGateway(config *MQTTConfig) *MQTTGateway {
    return &MQTTGateway{
        config:      config,
        dataChan:    make(chan *SensorData, 1000),
        messageChan: make(chan *MQTTMessage, 1000),
        running:     false,
    }
}

func (mg *MQTTGateway) Connect() error {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mg.config.Broker, mg.config.Port))
    opts.SetClientID(mg.config.ClientID)
    opts.SetUsername(mg.config.Username)
    opts.SetPassword(mg.config.Password)
    opts.SetDefaultPublishHandler(mg.messageHandler)
    opts.SetOnConnectHandler(mg.connectHandler)
    opts.SetConnectionLostHandler(mg.connectionLostHandler)
    
    mg.client = mqtt.NewClient(opts)
    token := mg.client.Connect()
    if token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    return nil
}

func (mg *MQTTGateway) Disconnect() {
    if mg.client != nil {
        mg.client.Disconnect(250)
    }
}

func (mg *MQTTGateway) Start() {
    mg.mu.Lock()
    defer mg.mu.Unlock()
    
    if mg.running {
        return
    }
    
    mg.running = true
    go mg.publishLoop()
}

func (mg *MQTTGateway) Stop() {
    mg.mu.Lock()
    defer mg.mu.Unlock()
    
    mg.running = false
}

func (mg *MQTTGateway) messageHandler(client mqtt.Client, msg mqtt.Message) {
    // 处理接收到的MQTT消息
    var payload map[string]interface{}
    if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
        log.Printf("Error unmarshaling MQTT message: %v", err)
        return
    }
    
    mqttMessage := &MQTTMessage{
        Topic:   msg.Topic(),
        Payload: payload,
        QoS:     int(msg.Qos()),
        Retain:  msg.Retained(),
    }
    
    select {
    case mg.messageChan <- mqttMessage:
    default:
        log.Printf("Message channel full, dropping MQTT message")
    }
}

func (mg *MQTTGateway) connectHandler(client mqtt.Client) {
    log.Printf("Connected to MQTT broker: %s", mg.config.Broker)
    
    // 订阅相关主题
    topics := map[string]byte{
        "device/+/command": byte(mg.config.QoS),
        "device/+/config":  byte(mg.config.QoS),
        "system/+/status":  byte(mg.config.QoS),
    }
    
    token := client.SubscribeMultiple(topics, nil)
    if token.Wait() && token.Error() != nil {
        log.Printf("Error subscribing to topics: %v", token.Error())
    }
}

func (mg *MQTTGateway) connectionLostHandler(client mqtt.Client, err error) {
    log.Printf("Connection lost to MQTT broker: %v", err)
    
    // 尝试重新连接
    go func() {
        for {
            time.Sleep(5 * time.Second)
            if token := client.Connect(); token.Wait() && token.Error() != nil {
                log.Printf("Failed to reconnect: %v", token.Error())
                continue
            }
            log.Printf("Reconnected to MQTT broker")
            break
        }
    }()
}

func (mg *MQTTGateway) publishLoop() {
    for mg.running {
        select {
        case data := <-mg.dataChan:
            mg.publishSensorData(data)
        case msg := <-mg.messageChan:
            mg.handleMessage(msg)
        }
    }
}

func (mg *MQTTGateway) publishSensorData(data *SensorData) {
    topic := fmt.Sprintf("device/%s/sensor/%s/data", data.DeviceID, data.SensorID)
    
    payload, err := json.Marshal(data)
    if err != nil {
        log.Printf("Error marshaling sensor data: %v", err)
        return
    }
    
    token := mg.client.Publish(topic, byte(mg.config.QoS), false, payload)
    if token.Wait() && token.Error() != nil {
        log.Printf("Error publishing sensor data: %v", token.Error())
    }
}

func (mg *MQTTGateway) handleMessage(msg *MQTTMessage) {
    // 处理不同类型的消息
    switch {
    case strings.HasPrefix(msg.Topic, "device/") && strings.HasSuffix(msg.Topic, "/command"):
        mg.handleDeviceCommand(msg)
    case strings.HasPrefix(msg.Topic, "device/") && strings.HasSuffix(msg.Topic, "/config"):
        mg.handleDeviceConfig(msg)
    case strings.HasPrefix(msg.Topic, "system/") && strings.HasSuffix(msg.Topic, "/status"):
        mg.handleSystemStatus(msg)
    default:
        log.Printf("Unknown message topic: %s", msg.Topic)
    }
}

func (mg *MQTTGateway) handleDeviceCommand(msg *MQTTMessage) {
    // 处理设备命令
    deviceID := strings.Split(msg.Topic, "/")[1]
    command := msg.Payload["command"].(string)
    
    log.Printf("Received command for device %s: %s", deviceID, command)
    
    // 这里应该执行具体的设备命令
    switch command {
    case "restart":
        log.Printf("Restarting device: %s", deviceID)
    case "update_config":
        log.Printf("Updating config for device: %s", deviceID)
    case "get_status":
        mg.publishDeviceStatus(deviceID)
    }
}

func (mg *MQTTGateway) handleDeviceConfig(msg *MQTTMessage) {
    // 处理设备配置更新
    deviceID := strings.Split(msg.Topic, "/")[1]
    log.Printf("Received config update for device: %s", deviceID)
}

func (mg *MQTTGateway) handleSystemStatus(msg *MQTTMessage) {
    // 处理系统状态消息
    systemID := strings.Split(msg.Topic, "/")[1]
    log.Printf("Received system status from: %s", systemID)
}

func (mg *MQTTGateway) publishDeviceStatus(deviceID string) {
    topic := fmt.Sprintf("device/%s/status", deviceID)
    
    status := map[string]interface{}{
        "device_id": deviceID,
        "status":    "online",
        "timestamp": time.Now().Unix(),
        "uptime":    3600, // 示例数据
    }
    
    payload, err := json.Marshal(status)
    if err != nil {
        log.Printf("Error marshaling device status: %v", err)
        return
    }
    
    token := mg.client.Publish(topic, byte(mg.config.QoS), false, payload)
    if token.Wait() && token.Error() != nil {
        log.Printf("Error publishing device status: %v", token.Error())
    }
}

func (mg *MQTTGateway) SetDataChannel(dataChan chan *SensorData) {
    mg.dataChan = dataChan
}

func (mg *MQTTGateway) GetMessageChannel() <-chan *MQTTMessage {
    return mg.messageChan
}
```

### 边缘计算处理器 (Edge Computing Processor)

```go
// 数据处理任务
type DataTask struct {
    ID       string                 `json:"id"`
    Type     string                 `json:"type"`
    Config   map[string]interface{} `json:"config"`
    Priority int                    `json:"priority"`
    Created  time.Time              `json:"created"`
}

// 处理结果
type ProcessingResult struct {
    TaskID   string                 `json:"task_id"`
    Success  bool                   `json:"success"`
    Data     map[string]interface{} `json:"data"`
    Error    string                 `json:"error,omitempty"`
    Duration time.Duration          `json:"duration"`
}

// 边缘处理器
type EdgeProcessor struct {
    tasks       map[string]*DataTask
    processors  map[string]DataProcessor
    resultChan  chan *ProcessingResult
    running     bool
    mu          sync.RWMutex
}

// 数据处理器接口
type DataProcessor interface {
    Process(data *SensorData, config map[string]interface{}) (map[string]interface{}, error)
    Type() string
}

// 基础数据处理器
type BaseDataProcessor struct {
    processorType string
}

func (bdp *BaseDataProcessor) Type() string {
    return bdp.processorType
}

// 平均值处理器
type AverageProcessor struct {
    BaseDataProcessor
    windowSize int
    dataBuffer []float64
    mu         sync.RWMutex
}

func NewAverageProcessor(windowSize int) *AverageProcessor {
    return &AverageProcessor{
        BaseDataProcessor: BaseDataProcessor{processorType: "average"},
        windowSize:        windowSize,
        dataBuffer:        make([]float64, 0, windowSize),
    }
}

func (ap *AverageProcessor) Process(data *SensorData, config map[string]interface{}) (map[string]interface{}, error) {
    ap.mu.Lock()
    defer ap.mu.Unlock()
    
    // 添加新数据到缓冲区
    ap.dataBuffer = append(ap.dataBuffer, data.Value)
    
    // 保持窗口大小
    if len(ap.dataBuffer) > ap.windowSize {
        ap.dataBuffer = ap.dataBuffer[1:]
    }
    
    // 计算平均值
    if len(ap.dataBuffer) == 0 {
        return nil, fmt.Errorf("no data available")
    }
    
    sum := 0.0
    for _, value := range ap.dataBuffer {
        sum += value
    }
    average := sum / float64(len(ap.dataBuffer))
    
    return map[string]interface{}{
        "average":     average,
        "window_size": len(ap.dataBuffer),
        "min":         ap.getMin(),
        "max":         ap.getMax(),
    }, nil
}

func (ap *AverageProcessor) getMin() float64 {
    if len(ap.dataBuffer) == 0 {
        return 0
    }
    min := ap.dataBuffer[0]
    for _, value := range ap.dataBuffer {
        if value < min {
            min = value
        }
    }
    return min
}

func (ap *AverageProcessor) getMax() float64 {
    if len(ap.dataBuffer) == 0 {
        return 0
    }
    max := ap.dataBuffer[0]
    for _, value := range ap.dataBuffer {
        if value > max {
            max = value
        }
    }
    return max
}

// 阈值检测处理器
type ThresholdProcessor struct {
    BaseDataProcessor
    threshold float64
    operator  string // "gt", "lt", "eq", "gte", "lte"
}

func NewThresholdProcessor(threshold float64, operator string) *ThresholdProcessor {
    return &ThresholdProcessor{
        BaseDataProcessor: BaseDataProcessor{processorType: "threshold"},
        threshold:         threshold,
        operator:          operator,
    }
}

func (tp *ThresholdProcessor) Process(data *SensorData, config map[string]interface{}) (map[string]interface{}, error) {
    var triggered bool
    
    switch tp.operator {
    case "gt":
        triggered = data.Value > tp.threshold
    case "lt":
        triggered = data.Value < tp.threshold
    case "eq":
        triggered = data.Value == tp.threshold
    case "gte":
        triggered = data.Value >= tp.threshold
    case "lte":
        triggered = data.Value <= tp.threshold
    default:
        return nil, fmt.Errorf("unknown operator: %s", tp.operator)
    }
    
    return map[string]interface{}{
        "triggered":  triggered,
        "value":      data.Value,
        "threshold":  tp.threshold,
        "operator":   tp.operator,
        "timestamp":  data.Timestamp,
    }, nil
}

// 异常检测处理器
type AnomalyProcessor struct {
    BaseDataProcessor
    mean     float64
    stdDev   float64
    dataCount int
    mu       sync.RWMutex
}

func NewAnomalyProcessor() *AnomalyProcessor {
    return &AnomalyProcessor{
        BaseDataProcessor: BaseDataProcessor{processorType: "anomaly"},
    }
}

func (ap *AnomalyProcessor) Process(data *SensorData, config map[string]interface{}) (map[string]interface{}, error) {
    ap.mu.Lock()
    defer ap.mu.Unlock()
    
    // 更新统计信息
    ap.dataCount++
    delta := data.Value - ap.mean
    ap.mean += delta / float64(ap.dataCount)
    delta2 := data.Value - ap.mean
    ap.stdDev += delta * delta2
    
    if ap.dataCount > 1 {
        ap.stdDev = math.Sqrt(ap.stdDev / float64(ap.dataCount-1))
    }
    
    // 检测异常（超过2个标准差）
    isAnomaly := false
    if ap.dataCount > 10 && ap.stdDev > 0 {
        zScore := math.Abs((data.Value - ap.mean) / ap.stdDev)
        isAnomaly = zScore > 2.0
    }
    
    return map[string]interface{}{
        "is_anomaly": isAnomaly,
        "value":      data.Value,
        "mean":       ap.mean,
        "std_dev":    ap.stdDev,
        "z_score":    (data.Value - ap.mean) / ap.stdDev,
        "data_count": ap.dataCount,
    }, nil
}

func NewEdgeProcessor() *EdgeProcessor {
    return &EdgeProcessor{
        tasks:      make(map[string]*DataTask),
        processors: make(map[string]DataProcessor),
        resultChan: make(chan *ProcessingResult, 1000),
        running:    false,
    }
}

func (ep *EdgeProcessor) RegisterProcessor(processor DataProcessor) {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    
    ep.processors[processor.Type()] = processor
}

func (ep *EdgeProcessor) AddTask(task *DataTask) error {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    
    if task.ID == "" {
        task.ID = generateTaskID()
    }
    
    task.Created = time.Now()
    ep.tasks[task.ID] = task
    return nil
}

func (ep *EdgeProcessor) Start() {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    
    if ep.running {
        return
    }
    
    ep.running = true
    go ep.processingLoop()
}

func (ep *EdgeProcessor) Stop() {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    
    ep.running = false
}

func (ep *EdgeProcessor) processingLoop() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for ep.running {
        select {
        case <-ticker.C:
            ep.processTasks()
        }
    }
}

func (ep *EdgeProcessor) processTasks() {
    ep.mu.RLock()
    tasks := make([]*DataTask, 0, len(ep.tasks))
    for _, task := range ep.tasks {
        tasks = append(tasks, task)
    }
    processors := make(map[string]DataProcessor)
    for k, v := range ep.processors {
        processors[k] = v
    }
    ep.mu.RUnlock()
    
    for _, task := range tasks {
        go ep.processTask(task, processors)
    }
}

func (ep *EdgeProcessor) processTask(task *DataTask, processors map[string]DataProcessor) {
    processor, exists := processors[task.Type]
    if !exists {
        result := &ProcessingResult{
            TaskID:  task.ID,
            Success: false,
            Error:   fmt.Sprintf("processor not found: %s", task.Type),
        }
        ep.resultChan <- result
        return
    }
    
    // 模拟传感器数据
    sensorData := &SensorData{
        ID:        generateDataID(),
        DeviceID:  "device_1",
        SensorID:  "sensor_1",
        Type:      "temperature",
        Value:     25 + rand.Float64()*10,
        Unit:      "°C",
        Timestamp: time.Now(),
    }
    
    startTime := time.Now()
    
    result, err := processor.Process(sensorData, task.Config)
    
    duration := time.Since(startTime)
    
    processingResult := &ProcessingResult{
        TaskID:   task.ID,
        Success:  err == nil,
        Data:     result,
        Error:    "",
        Duration: duration,
    }
    
    if err != nil {
        processingResult.Error = err.Error()
    }
    
    ep.resultChan <- processingResult
}

func (ep *EdgeProcessor) GetResultChannel() <-chan *ProcessingResult {
    return ep.resultChan
}

func generateTaskID() string {
    return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
```

## 设计原则

### 1. 设备管理设计

- **设备注册**：自动发现和注册IoT设备
- **状态监控**：实时监控设备状态和健康度
- **配置管理**：远程配置和参数更新
- **固件更新**：安全的固件升级机制

### 2. 数据采集设计

- **多协议支持**：支持MQTT、CoAP、HTTP等协议
- **数据缓存**：本地缓存防止数据丢失
- **批量传输**：批量数据传输提高效率
- **数据压缩**：减少网络传输量

### 3. 边缘计算设计

- **本地处理**：在设备端进行数据处理
- **实时分析**：实时数据分析和决策
- **规则引擎**：可配置的数据处理规则
- **机器学习**：边缘AI和机器学习

### 4. 安全性设计

- **设备认证**：设备身份验证和授权
- **数据加密**：端到端数据加密
- **访问控制**：基于角色的访问控制
- **安全更新**：安全的软件更新机制

## 实现示例

```go
func main() {
    // 创建设备管理器
    deviceManager := NewDeviceManager()
    
    // 创建数据采集器
    dataCollector := NewDataCollector(deviceManager)
    
    // 创建MQTT网关
    mqttConfig := &MQTTConfig{
        Broker:   "localhost",
        Port:     1883,
        ClientID: "iot_gateway",
        Username: "user",
        Password: "password",
        QoS:      1,
    }
    mqttGateway := NewMQTTGateway(mqttConfig)
    
    // 创建边缘处理器
    edgeProcessor := NewEdgeProcessor()
    
    // 注册数据处理器
    edgeProcessor.RegisterProcessor(NewAverageProcessor(10))
    edgeProcessor.RegisterProcessor(NewThresholdProcessor(30.0, "gt"))
    edgeProcessor.RegisterProcessor(NewAnomalyProcessor())
    
    // 连接MQTT网关
    if err := mqttGateway.Connect(); err != nil {
        log.Printf("Failed to connect to MQTT broker: %v", err)
        return
    }
    defer mqttGateway.Disconnect()
    
    // 设置数据通道
    mqttGateway.SetDataChannel(dataCollector.GetDataChannel())
    
    // 启动系统
    dataCollector.Start()
    mqttGateway.Start()
    edgeProcessor.Start()
    
    // 添加示例设备
    device := &Device{
        ID:     "device_1",
        Name:   "Temperature Sensor",
        Type:   "sensor",
        Model:  "TEMP-001",
        Version: "1.0.0",
        Location: Location{
            Latitude:  37.7749,
            Longitude: -122.4194,
            Address:   "San Francisco, CA",
        },
        Properties: map[string]interface{}{
            "manufacturer": "IoT Corp",
            "battery":      85,
        },
        Tags: []string{"temperature", "outdoor"},
    }
    deviceManager.RegisterDevice(device)
    
    // 添加传感器配置
    sensorConfig := &SensorConfig{
        ID:           "sensor_1",
        DeviceID:     "device_1",
        Type:         "temperature",
        Name:         "Temperature Sensor",
        Unit:         "°C",
        MinValue:     0,
        MaxValue:     50,
        SamplingRate: 1,
        Enabled:      true,
    }
    dataCollector.AddSensor(sensorConfig)
    
    // 添加处理任务
    task := &DataTask{
        Type: "average",
        Config: map[string]interface{}{
            "window_size": 10,
        },
        Priority: 1,
    }
    edgeProcessor.AddTask(task)
    
    // 处理结果
    go func() {
        for result := range edgeProcessor.GetResultChannel() {
            log.Printf("Processing result: %+v", result)
        }
    }()
    
    // 等待一段时间
    time.Sleep(30 * time.Second)
    
    // 停止系统
    dataCollector.Stop()
    mqttGateway.Stop()
    edgeProcessor.Stop()
    
    fmt.Println("IoT system stopped")
}
```

## 总结

Go语言在物联网领域具有显著优势，特别适合构建高性能、低资源消耗的IoT设备和系统。

### 关键要点

1. **高性能**：编译型语言提供优秀的执行效率
2. **低资源消耗**：内存占用小，适合资源受限设备
3. **并发处理**：原生支持高并发数据处理
4. **跨平台**：支持多种硬件平台部署
5. **网络编程**：强大的网络库支持多种协议

### 发展趋势

- **5G IoT**：5G网络支持的大规模IoT部署
- **边缘AI**：边缘计算和人工智能结合
- **数字孪生**：物理世界的数字映射
- **工业IoT**：工业自动化和智能制造
- **智能城市**：智慧城市和基础设施管理
