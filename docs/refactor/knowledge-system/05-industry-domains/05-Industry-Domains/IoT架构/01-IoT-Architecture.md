# IoT架构 (Internet of Things Architecture)

## 1. 基本概念

### 1.1 IoT定义

**物联网 (Internet of Things, IoT)** 是指通过互联网连接各种物理设备、传感器、执行器等，使它们能够收集、传输和处理数据，实现智能化管理和控制的技术体系。

### 1.2 核心特征

- **感知层**: 通过各种传感器收集物理世界的数据
- **网络层**: 通过多种通信协议传输数据
- **平台层**: 提供数据处理、存储和分析能力
- **应用层**: 提供各种业务应用和服务

### 1.3 架构层次

```text
┌─────────────────────────────────────────┐
│                应用层                    │
│  (业务应用、数据分析、用户界面)           │
├─────────────────────────────────────────┤
│                平台层                    │
│  (数据处理、存储、分析、管理)             │
├─────────────────────────────────────────┤
│                网络层                    │
│  (通信协议、网关、路由)                  │
├─────────────────────────────────────────┤
│                感知层                    │
│  (传感器、执行器、设备)                  │
└─────────────────────────────────────────┘
```

## 2. IoT架构模式

### 2.1 三层架构

```go
// 感知层 - 传感器设备
type SensorDevice struct {
    id       string
    type     string
    location string
    data     map[string]interface{}
    mu       sync.RWMutex
}

func (sd *SensorDevice) CollectData() map[string]interface{} {
    sd.mu.Lock()
    defer sd.mu.Unlock()
    
    // 模拟数据收集
    sd.data["temperature"] = rand.Float64() * 50
    sd.data["humidity"] = rand.Float64() * 100
    sd.data["timestamp"] = time.Now()
    
    return sd.data
}

func (sd *SensorDevice) SendData(data map[string]interface{}) error {
    // 发送数据到网关
    return sd.sendToGateway(data)
}

// 网络层 - 网关
type Gateway struct {
    id       string
    devices  map[string]*SensorDevice
    platform *Platform
    mu       sync.RWMutex
}

func (g *Gateway) ReceiveData(deviceID string, data map[string]interface{}) error {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    // 数据预处理
    processedData := g.preprocessData(data)
    
    // 发送到平台层
    return g.platform.ProcessData(deviceID, processedData)
}

func (g *Gateway) preprocessData(data map[string]interface{}) map[string]interface{} {
    // 数据清洗、格式转换等
    processed := make(map[string]interface{})
    for key, value := range data {
        switch v := value.(type) {
        case float64:
            processed[key] = math.Round(v*100) / 100
        default:
            processed[key] = value
        }
    }
    return processed
}

// 平台层 - 数据处理平台
type Platform struct {
    id       string
    storage  DataStorage
    analyzer DataAnalyzer
    mu       sync.RWMutex
}

func (p *Platform) ProcessData(deviceID string, data map[string]interface{}) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // 存储数据
    if err := p.storage.Store(deviceID, data); err != nil {
        return err
    }
    
    // 分析数据
    go p.analyzer.Analyze(deviceID, data)
    
    return nil
}
```

### 2.2 边缘计算架构

```go
// 边缘节点
type EdgeNode struct {
    id       string
    devices  map[string]*SensorDevice
    compute  EdgeCompute
    storage  EdgeStorage
    mu       sync.RWMutex
}

type EdgeCompute struct {
    processor *Processor
    memory    *Memory
}

type EdgeStorage struct {
    localDB   *LocalDatabase
    cache     *Cache
}

func (en *EdgeNode) ProcessLocally(deviceID string, data map[string]interface{}) error {
    en.mu.Lock()
    defer en.mu.Unlock()
    
    // 本地数据处理
    processedData := en.compute.processor.Process(data)
    
    // 本地存储
    if err := en.storage.localDB.Store(deviceID, processedData); err != nil {
        return err
    }
    
    // 缓存热点数据
    en.storage.cache.Set(deviceID, processedData)
    
    // 判断是否需要上传到云端
    if en.needUploadToCloud(processedData) {
        go en.uploadToCloud(deviceID, processedData)
    }
    
    return nil
}

func (en *EdgeNode) needUploadToCloud(data map[string]interface{}) bool {
    // 根据数据重要性、大小等判断是否需要上传
    if importance, ok := data["importance"].(int); ok && importance > 7 {
        return true
    }
    return false
}
```

## 3. 通信协议

### 3.1 MQTT协议

```go
// MQTT客户端
type MQTTClient struct {
    client   mqtt.Client
    broker   string
    clientID string
    topics   map[string]mqtt.MessageHandler
}

func (mc *MQTTClient) Connect() error {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(mc.broker)
    opts.SetClientID(mc.clientID)
    opts.SetDefaultPublishHandler(mc.defaultMessageHandler)
    
    mc.client = mqtt.NewClient(opts)
    if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    
    return nil
}

func (mc *MQTTClient) Publish(topic string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    token := mc.client.Publish(topic, 0, false, data)
    token.Wait()
    return token.Error()
}

func (mc *MQTTClient) Subscribe(topic string, handler mqtt.MessageHandler) error {
    token := mc.client.Subscribe(topic, 0, handler)
    token.Wait()
    return token.Error()
}

func (mc *MQTTClient) defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
    log.Printf("Received message: %s from topic: %s", string(msg.Payload()), msg.Topic())
}

// MQTT消息处理器
type MQTTHandler struct {
    client *MQTTClient
}

func (mh *MQTTHandler) HandleSensorData(client mqtt.Client, msg mqtt.Message) {
    var sensorData map[string]interface{}
    if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
        log.Printf("Failed to unmarshal sensor data: %v", err)
        return
    }
    
    // 处理传感器数据
    mh.processSensorData(sensorData)
}

func (mh *MQTTHandler) processSensorData(data map[string]interface{}) {
    // 数据验证
    if !mh.validateData(data) {
        log.Printf("Invalid sensor data: %v", data)
        return
    }
    
    // 数据转换
    processedData := mh.transformData(data)
    
    // 存储数据
    if err := mh.storeData(processedData); err != nil {
        log.Printf("Failed to store data: %v", err)
    }
}
```

### 3.2 CoAP协议

```go
// CoAP客户端
type CoAPClient struct {
    client *coap.Client
    server string
}

func (cc *CoAPClient) Get(resource string) ([]byte, error) {
    req, err := coap.NewRequest(coap.GET, resource)
    if err != nil {
        return nil, err
    }
    
    resp, err := cc.client.Send(req)
    if err != nil {
        return nil, err
    }
    
    return resp.Payload, nil
}

func (cc *CoAPClient) Post(resource string, payload []byte) error {
    req, err := coap.NewRequest(coap.POST, resource)
    if err != nil {
        return err
    }
    
    req.SetPayload(payload)
    
    _, err = cc.client.Send(req)
    return err
}

// CoAP服务器
type CoAPServer struct {
    server *coap.Server
    handlers map[string]coap.HandlerFunc
}

func (cs *CoAPServer) Start(addr string) error {
    cs.server = &coap.Server{
        Addr: addr,
        Handler: cs.handleRequest,
    }
    
    return cs.server.ListenAndServe()
}

func (cs *CoAPServer) handleRequest(w coap.ResponseWriter, r *coap.Request) {
    handler, exists := cs.handlers[r.URL.Path]
    if exists {
        handler(w, r)
    } else {
        w.WriteHeader(coap.NotFound)
    }
}

func (cs *CoAPServer) RegisterHandler(path string, handler coap.HandlerFunc) {
    cs.handlers[path] = handler
}
```

### 3.3 HTTP/HTTPS协议

```go
// HTTP客户端
type HTTPClient struct {
    client *http.Client
    baseURL string
}

func (hc *HTTPClient) SendData(endpoint string, data map[string]interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    url := hc.baseURL + endpoint
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := hc.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
    }
    
    return nil
}

// HTTP服务器
type HTTPServer struct {
    server *http.Server
    router *mux.Router
}

func (hs *HTTPServer) Start(addr string) error {
    hs.router = mux.NewRouter()
    hs.setupRoutes()
    
    hs.server = &http.Server{
        Addr:    addr,
        Handler: hs.router,
    }
    
    return hs.server.ListenAndServe()
}

func (hs *HTTPServer) setupRoutes() {
    hs.router.HandleFunc("/api/sensor/data", hs.handleSensorData).Methods("POST")
    hs.router.HandleFunc("/api/device/status", hs.handleDeviceStatus).Methods("GET")
    hs.router.HandleFunc("/api/device/control", hs.handleDeviceControl).Methods("POST")
}

func (hs *HTTPServer) handleSensorData(w http.ResponseWriter, r *http.Request) {
    var sensorData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&sensorData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 处理传感器数据
    if err := hs.processSensorData(sensorData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

## 4. 设备管理

### 4.1 设备注册

```go
// 设备管理器
type DeviceManager struct {
    devices map[string]*Device
    mu      sync.RWMutex
}

type Device struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Location    string            `json:"location"`
    Status      string            `json:"status"`
    Properties  map[string]string `json:"properties"`
    LastSeen    time.Time         `json:"last_seen"`
    CreatedAt   time.Time         `json:"created_at"`
}

func (dm *DeviceManager) RegisterDevice(device *Device) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    // 检查设备是否已存在
    if _, exists := dm.devices[device.ID]; exists {
        return fmt.Errorf("device %s already exists", device.ID)
    }
    
    // 设置创建时间
    device.CreatedAt = time.Now()
    device.LastSeen = time.Now()
    
    // 注册设备
    dm.devices[device.ID] = device
    
    // 保存到数据库
    return dm.saveToDatabase(device)
}

func (dm *DeviceManager) GetDevice(deviceID string) (*Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    device, exists := dm.devices[deviceID]
    if !exists {
        return nil, fmt.Errorf("device %s not found", deviceID)
    }
    
    return device, nil
}

func (dm *DeviceManager) UpdateDeviceStatus(deviceID string, status string) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    device, exists := dm.devices[deviceID]
    if !exists {
        return fmt.Errorf("device %s not found", deviceID)
    }
    
    device.Status = status
    device.LastSeen = time.Now()
    
    return dm.updateInDatabase(device)
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
```

### 4.2 设备监控

```go
// 设备监控器
type DeviceMonitor struct {
    devices map[string]*Device
    alerts  chan Alert
    mu      sync.RWMutex
}

type Alert struct {
    DeviceID  string    `json:"device_id"`
    Type      string    `json:"type"`
    Message   string    `json:"message"`
    Severity  string    `json:"severity"`
    Timestamp time.Time `json:"timestamp"`
}

func (dm *DeviceMonitor) StartMonitoring() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        dm.checkDevices()
    }
}

func (dm *DeviceMonitor) checkDevices() {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    now := time.Now()
    for deviceID, device := range dm.devices {
        // 检查设备是否离线
        if now.Sub(device.LastSeen) > 5*time.Minute {
            alert := Alert{
                DeviceID:  deviceID,
                Type:      "offline",
                Message:   fmt.Sprintf("Device %s is offline", deviceID),
                Severity:  "high",
                Timestamp: now,
            }
            dm.alerts <- alert
        }
        
        // 检查设备状态
        if device.Status == "error" {
            alert := Alert{
                DeviceID:  deviceID,
                Type:      "error",
                Message:   fmt.Sprintf("Device %s is in error state", deviceID),
                Severity:  "medium",
                Timestamp: now,
            }
            dm.alerts <- alert
        }
    }
}

func (dm *DeviceMonitor) ProcessAlerts() {
    for alert := range dm.alerts {
        dm.handleAlert(alert)
    }
}

func (dm *DeviceMonitor) handleAlert(alert Alert) {
    log.Printf("Alert: %s - %s", alert.Type, alert.Message)
    
    // 根据告警类型处理
    switch alert.Type {
    case "offline":
        dm.handleOfflineAlert(alert)
    case "error":
        dm.handleErrorAlert(alert)
    default:
        log.Printf("Unknown alert type: %s", alert.Type)
    }
}

func (dm *DeviceMonitor) handleOfflineAlert(alert Alert) {
    // 尝试重新连接设备
    log.Printf("Attempting to reconnect device: %s", alert.DeviceID)
    
    // 发送通知
    dm.sendNotification(alert)
}

func (dm *DeviceMonitor) handleErrorAlert(alert Alert) {
    // 记录错误日志
    log.Printf("Device error: %s", alert.Message)
    
    // 发送通知
    dm.sendNotification(alert)
}

func (dm *DeviceMonitor) sendNotification(alert Alert) {
    // 发送邮件、短信或其他通知
    log.Printf("Sending notification for alert: %s", alert.Message)
}
```

## 5. 数据处理

### 5.1 数据流处理

```go
// 数据流处理器
type DataStreamProcessor struct {
    input   chan DataPoint
    output  chan ProcessedData
    filters []DataFilter
    mu      sync.RWMutex
}

type DataPoint struct {
    DeviceID  string                 `json:"device_id"`
    Timestamp time.Time              `json:"timestamp"`
    Values    map[string]interface{} `json:"values"`
}

type ProcessedData struct {
    DeviceID  string                 `json:"device_id"`
    Timestamp time.Time              `json:"timestamp"`
    Values    map[string]interface{} `json:"values"`
    Metadata  map[string]interface{} `json:"metadata"`
}

type DataFilter interface {
    Filter(data DataPoint) (DataPoint, error)
}

func (dsp *DataStreamProcessor) Start() {
    go dsp.process()
}

func (dsp *DataStreamProcessor) process() {
    for dataPoint := range dsp.input {
        processedData, err := dsp.applyFilters(dataPoint)
        if err != nil {
            log.Printf("Failed to process data point: %v", err)
            continue
        }
        
        dsp.output <- processedData
    }
}

func (dsp *DataStreamProcessor) applyFilters(dataPoint DataPoint) (ProcessedData, error) {
    dsp.mu.RLock()
    defer dsp.mu.RUnlock()
    
    processed := dataPoint
    
    // 应用过滤器
    for _, filter := range dsp.filters {
        filtered, err := filter.Filter(processed)
        if err != nil {
            return ProcessedData{}, err
        }
        processed = filtered
    }
    
    // 转换为处理后的数据
    return ProcessedData{
        DeviceID:  processed.DeviceID,
        Timestamp: processed.Timestamp,
        Values:    processed.Values,
        Metadata:  dsp.generateMetadata(processed),
    }, nil
}

func (dsp *DataStreamProcessor) generateMetadata(dataPoint DataPoint) map[string]interface{} {
    return map[string]interface{}{
        "processed_at": time.Now(),
        "filter_count": len(dsp.filters),
        "data_size":    len(dataPoint.Values),
    }
}

// 数据过滤器示例
type OutlierFilter struct {
    threshold float64
}

func (of *OutlierFilter) Filter(data DataPoint) (DataPoint, error) {
    filtered := data
    filtered.Values = make(map[string]interface{})
    
    for key, value := range data.Values {
        if num, ok := value.(float64); ok {
            if math.Abs(num) <= of.threshold {
                filtered.Values[key] = value
            }
        } else {
            filtered.Values[key] = value
        }
    }
    
    return filtered, nil
}
```

### 5.2 数据分析

```go
// 数据分析器
type DataAnalyzer struct {
    storage DataStorage
    mu      sync.RWMutex
}

func (da *DataAnalyzer) AnalyzeTrends(deviceID string, timeRange time.Duration) ([]Trend, error) {
    da.mu.Lock()
    defer da.mu.Unlock()
    
    // 获取历史数据
    data, err := da.storage.GetHistoricalData(deviceID, timeRange)
    if err != nil {
        return nil, err
    }
    
    // 分析趋势
    trends := da.calculateTrends(data)
    
    return trends, nil
}

type Trend struct {
    Metric    string  `json:"metric"`
    Direction string  `json:"direction"` // "increasing", "decreasing", "stable"
    Slope     float64 `json:"slope"`
    R2        float64 `json:"r2"`
}

func (da *DataAnalyzer) calculateTrends(data []DataPoint) []Trend {
    trends := make([]Trend, 0)
    
    // 按指标分组
    metrics := da.groupByMetric(data)
    
    for metric, values := range metrics {
        trend := da.calculateTrend(metric, values)
        trends = append(trends, trend)
    }
    
    return trends
}

func (da *DataAnalyzer) groupByMetric(data []DataPoint) map[string][]float64 {
    metrics := make(map[string][]float64)
    
    for _, point := range data {
        for metric, value := range point.Values {
            if num, ok := value.(float64); ok {
                metrics[metric] = append(metrics[metric], num)
            }
        }
    }
    
    return metrics
}

func (da *DataAnalyzer) calculateTrend(metric string, values []float64) Trend {
    if len(values) < 2 {
        return Trend{
            Metric:    metric,
            Direction: "stable",
            Slope:     0,
            R2:        0,
        }
    }
    
    // 简单线性回归
    n := len(values)
    sumX := 0.0
    sumY := 0.0
    sumXY := 0.0
    sumX2 := 0.0
    
    for i, y := range values {
        x := float64(i)
        sumX += x
        sumY += y
        sumXY += x * y
        sumX2 += x * x
    }
    
    slope := (float64(n)*sumXY - sumX*sumY) / (float64(n)*sumX2 - sumX*sumX)
    
    // 确定趋势方向
    direction := "stable"
    if slope > 0.01 {
        direction = "increasing"
    } else if slope < -0.01 {
        direction = "decreasing"
    }
    
    return Trend{
        Metric:    metric,
        Direction: direction,
        Slope:     slope,
        R2:        0.8, // 简化的R²值
    }
}
```

## 6. 安全与隐私

### 6.1 设备认证

```go
// 设备认证器
type DeviceAuthenticator struct {
    devices map[string]*DeviceCredentials
    mu      sync.RWMutex
}

type DeviceCredentials struct {
    DeviceID   string `json:"device_id"`
    PublicKey  string `json:"public_key"`
    PrivateKey string `json:"private_key"`
    Token      string `json:"token"`
    ExpiresAt  time.Time `json:"expires_at"`
}

func (da *DeviceAuthenticator) AuthenticateDevice(deviceID, token string) (bool, error) {
    da.mu.RLock()
    defer da.mu.RUnlock()
    
    credentials, exists := da.devices[deviceID]
    if !exists {
        return false, fmt.Errorf("device %s not found", deviceID)
    }
    
    // 检查token是否有效
    if credentials.Token != token {
        return false, fmt.Errorf("invalid token for device %s", deviceID)
    }
    
    // 检查token是否过期
    if time.Now().After(credentials.ExpiresAt) {
        return false, fmt.Errorf("token expired for device %s", deviceID)
    }
    
    return true, nil
}

func (da *DeviceAuthenticator) GenerateToken(deviceID string) (string, error) {
    da.mu.Lock()
    defer da.mu.Unlock()
    
    credentials, exists := da.devices[deviceID]
    if !exists {
        return "", fmt.Errorf("device %s not found", deviceID)
    }
    
    // 生成新的token
    token := da.generateRandomToken()
    credentials.Token = token
    credentials.ExpiresAt = time.Now().Add(24 * time.Hour)
    
    return token, nil
}

func (da *DeviceAuthenticator) generateRandomToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

### 6.2 数据加密

```go
// 数据加密器
type DataEncryptor struct {
    key []byte
}

func (de *DataEncryptor) Encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(de.key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    return gcm.Seal(nonce, nonce, data, nil), nil
}

func (de *DataEncryptor) Decrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(de.key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
```

## 总结

IoT架构是一个复杂的系统，涉及多个层次和组件。成功的IoT系统需要：

**关键设计原则**：

1. 可扩展性：支持大量设备接入
2. 可靠性：确保数据传输和处理的可靠性
3. 安全性：保护设备和数据的安全
4. 实时性：支持实时数据处理和响应
5. 互操作性：支持不同设备和协议的互操作

**常见挑战**：

1. 设备管理和监控
2. 数据安全和隐私保护
3. 网络连接和通信协议
4. 数据处理和分析
5. 系统集成和互操作性

IoT架构的设计需要根据具体的应用场景和需求来选择合适的组件和技术，确保系统的可靠性、安全性和可扩展性。
