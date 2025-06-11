# 3. 技术实现文档

## 3.1 代码架构分析

### 3.1.1 项目结构

```text
navigate/
├── cmd/server/navlock-gezhouba/     # 服务入口
├── common/                          # 公共组件
│   ├── model/                       # 模型定义
│   ├── runtime/                     # 运行时
│   └── utility.go                   # 工具函数
├── initialize/                      # 初始化模块
│   └── iot/navlock/gezhouba/        # 葛洲坝船闸初始化
├── iot/                            # IoT模块
│   └── navlock/                     # 船闸相关
├── config/                         # 配置管理
├── global/                         # 全局定义
└── docs/gezhouba/                  # 项目文档
```

### 3.1.2 核心模块分析

1. **服务入口模块**
   - 位置：`cmd/server/navlock-gezhouba/main.go`
   - 职责：系统启动入口，生命周期管理
   - 特点：简洁的入口设计，委托给服务主控

2. **初始化模块**
   - 位置：`initialize/iot/navlock/gezhouba/`
   - 职责：系统初始化，组件创建，配置加载
   - 特点：集中管理初始化逻辑

3. **IoT模块**
   - 位置：`iot/navlock/`
   - 职责：船闸业务逻辑，设备控制
   - 特点：业务逻辑封装，设备抽象

4. **公共模块**
   - 位置：`common/`
   - 职责：通用功能，模型定义，工具函数
   - 特点：可复用组件，标准化接口

## 3.2 核心算法实现

### 3.2.1 船舶位置计算算法

```go
// 雷达数据处理算法
func ProcessRadarData(rawData []byte) (*ShipPosition, error) {
    // 1. 数据解析
    radarData := ParseRadarData(rawData)
    
    // 2. 坐标转换
    worldCoord := ConvertToWorldCoordinates(radarData)
    
    // 3. 位置计算
    position := CalculateShipPosition(worldCoord)
    
    // 4. 轨迹分析
    trajectory := AnalyzeTrajectory(position)
    
    return &ShipPosition{
        X: position.X,
        Y: position.Y,
        Speed: position.Speed,
        Direction: position.Direction,
        Timestamp: time.Now(),
    }, nil
}
```

### 3.2.2 船舶速度计算算法

```go
// 速度计算算法
func CalculateShipSpeed(positions []ShipPosition) float64 {
    if len(positions) < 2 {
        return 0.0
    }
    
    // 计算位移
    dx := positions[len(positions)-1].X - positions[0].X
    dy := positions[len(positions)-1].Y - positions[0].Y
    distance := math.Sqrt(dx*dx + dy*dy)
    
    // 计算时间间隔
    timeDiff := positions[len(positions)-1].Timestamp.Sub(positions[0].Timestamp)
    
    // 计算速度 (m/s)
    speed := distance / timeDiff.Seconds()
    
    return speed
}
```

### 3.2.3 禁停区域判断算法

```go
// 禁停区域判断
func IsInNoStopZone(position ShipPosition, noStopZones []NoStopZone) bool {
    for _, zone := range noStopZones {
        if IsPointInPolygon(position.X, position.Y, zone.Boundary) {
            return true
        }
    }
    return false
}

// 点在多边形内判断
func IsPointInPolygon(x, y float64, polygon []Point) bool {
    n := len(polygon)
    inside := false
    
    j := n - 1
    for i := 0; i < n; i++ {
        if ((polygon[i].Y > y) != (polygon[j].Y > y)) &&
            (x < (polygon[j].X-polygon[i].X)*(y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X) {
            inside = !inside
        }
        j = i
    }
    
    return inside
}
```

## 3.3 数据库设计

### 3.3.1 数据库架构

系统采用多数据库架构：

1. **本地数据库 (LocalDB)**
   - 类型：SQLite
   - 用途：缓存、临时数据、配置
   - 特点：轻量级，本地存储

2. **生产数据库 (ProduceDB)**
   - 类型：MySQL
   - 用途：业务数据、历史数据
   - 特点：高性能，可靠性

3. **Web数据库 (ProduceWebDB)**
   - 类型：MySQL
   - 用途：展示数据、报表数据
   - 特点：专门用于Web展示

### 3.3.2 核心数据表设计

```sql
-- 船舶信息表
CREATE TABLE ships (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ship_id VARCHAR(50) NOT NULL,
    ship_name VARCHAR(100),
    ship_type VARCHAR(50),
    length DECIMAL(10,2),
    width DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 船舶位置记录表
CREATE TABLE ship_positions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ship_id VARCHAR(50) NOT NULL,
    x_coordinate DECIMAL(10,6),
    y_coordinate DECIMAL(10,6),
    speed DECIMAL(10,2),
    direction DECIMAL(5,2),
    timestamp TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 船闸状态表
CREATE TABLE lock_status (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    lock_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    direction VARCHAR(10),
    water_level DECIMAL(10,2),
    gate_position VARCHAR(20),
    timestamp TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 异常事件表
CREATE TABLE events (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    event_type VARCHAR(50) NOT NULL,
    event_level VARCHAR(20),
    ship_id VARCHAR(50),
    description TEXT,
    position_x DECIMAL(10,6),
    position_y DECIMAL(10,6),
    timestamp TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3.3.3 数据访问层

```go
// GORM模型定义
type Ship struct {
    ID        uint      `gorm:"primaryKey"`
    ShipID    string    `gorm:"uniqueIndex;not null"`
    ShipName  string
    ShipType  string
    Length    decimal.Decimal
    Width     decimal.Decimal
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ShipPosition struct {
    ID         uint      `gorm:"primaryKey"`
    ShipID     string    `gorm:"index;not null"`
    XCoordinate decimal.Decimal
    YCoordinate decimal.Decimal
    Speed      decimal.Decimal
    Direction  decimal.Decimal
    Timestamp  time.Time
    CreatedAt  time.Time
}

// 数据访问接口
type ShipRepository interface {
    Create(ship *Ship) error
    GetByID(shipID string) (*Ship, error)
    Update(ship *Ship) error
    Delete(shipID string) error
    List(offset, limit int) ([]*Ship, error)
}
```

## 3.4 接口实现

### 3.4.1 NATS消息接口

```go
// 消息发布接口
type MessagePublisher interface {
    Publish(topic string, data []byte) error
    PublishAsync(topic string, data []byte) error
}

// 消息订阅接口
type MessageSubscriber interface {
    Subscribe(topic string, handler func(msg *nats.Msg)) error
    Unsubscribe(topic string) error
}

// 消息实现
type NATSMessageService struct {
    conn *nats.Conn
}

func (n *NATSMessageService) Publish(topic string, data []byte) error {
    return n.conn.Publish(topic, data)
}

func (n *NATSMessageService) Subscribe(topic string, handler func(msg *nats.Msg)) error {
    _, err := n.conn.Subscribe(topic, handler)
    return err
}
```

### 3.4.2 HTTP API接口

```go
// HTTP服务器
type HTTPServer struct {
    router *gin.Engine
    port   string
}

// 船舶状态API
func (s *HTTPServer) GetShipStatus(c *gin.Context) {
    shipID := c.Param("id")
    
    ship, err := s.shipService.GetShipStatus(shipID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ship not found"})
        return
    }
    
    c.JSON(http.StatusOK, ship)
}

// 船闸状态API
func (s *HTTPServer) GetLockStatus(c *gin.Context) {
    lockID := c.Param("id")
    
    status, err := s.lockService.GetLockStatus(lockID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, status)
}
```

### 3.4.3 设备控制接口

```go
// 设备控制接口
type DeviceController interface {
    Connect() error
    Disconnect() error
    SendCommand(cmd Command) error
    GetStatus() (*DeviceStatus, error)
}

// 雷达设备控制
type RadarController struct {
    connection net.Conn
    config     RadarConfig
}

func (r *RadarController) SendCommand(cmd Command) error {
    data := cmd.Serialize()
    _, err := r.connection.Write(data)
    return err
}

func (r *RadarController) GetStatus() (*DeviceStatus, error) {
    // 发送状态查询命令
    cmd := Command{Type: "STATUS_QUERY"}
    err := r.SendCommand(cmd)
    if err != nil {
        return nil, err
    }
    
    // 读取响应
    response := make([]byte, 1024)
    n, err := r.connection.Read(response)
    if err != nil {
        return nil, err
    }
    
    return ParseDeviceStatus(response[:n])
}
```

## 3.5 配置管理实现

### 3.5.1 配置结构定义

```go
// 系统配置结构
type SystemConfig struct {
    Database DatabaseConfig `mapstructure:"database"`
    Network  NetworkConfig  `mapstructure:"network"`
    Devices  DevicesConfig  `mapstructure:"devices"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type DatabaseConfig struct {
    LocalDB    DBConfig `mapstructure:"local_db"`
    ProduceDB  DBConfig `mapstructure:"produce_db"`
    WebDB      DBConfig `mapstructure:"web_db"`
}

type DBConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
}

type NetworkConfig struct {
    NATSURL string `mapstructure:"nats_url"`
    HTTPPort int   `mapstructure:"http_port"`
}

type DevicesConfig struct {
    Radar RadarConfig `mapstructure:"radar"`
    LED   LEDConfig   `mapstructure:"led"`
    PTZ   PTZConfig   `mapstructure:"ptz"`
}
```

### 3.5.2 配置加载实现

```go
// 配置加载器
type ConfigLoader struct {
    viper *viper.Viper
}

func NewConfigLoader() *ConfigLoader {
    v := viper.New()
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath(".")
    v.AddConfigPath("./config")
    
    return &ConfigLoader{viper: v}
}

func (c *ConfigLoader) Load() (*SystemConfig, error) {
    err := c.viper.ReadInConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config SystemConfig
    err = c.viper.Unmarshal(&config)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &config, nil
}

func (c *ConfigLoader) Validate(config *SystemConfig) error {
    // 验证数据库配置
    if err := c.validateDatabaseConfig(config.Database); err != nil {
        return fmt.Errorf("database config error: %w", err)
    }
    
    // 验证网络配置
    if err := c.validateNetworkConfig(config.Network); err != nil {
        return fmt.Errorf("network config error: %w", err)
    }
    
    // 验证设备配置
    if err := c.validateDevicesConfig(config.Devices); err != nil {
        return fmt.Errorf("devices config error: %w", err)
    }
    
    return nil
}
```

## 3.6 日志系统实现

### 3.6.1 日志配置

```go
// 日志配置
type LoggingConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    OutputPath string `mapstructure:"output_path"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

// 日志初始化
func InitLogger(config LoggingConfig) (*zap.Logger, error) {
    // 配置日志级别
    level, err := zap.ParseAtomicLevel(config.Level)
    if err != nil {
        return nil, fmt.Errorf("invalid log level: %w", err)
    }
    
    // 配置输出
    var outputPaths []string
    if config.OutputPath != "" {
        outputPaths = append(outputPaths, config.OutputPath)
    } else {
        outputPaths = append(outputPaths, "stdout")
    }
    
    // 创建配置
    zapConfig := zap.NewProductionConfig()
    zapConfig.Level = level
    zapConfig.OutputPaths = outputPaths
    zapConfig.Encoding = config.Format
    
    // 创建日志记录器
    logger, err := zapConfig.Build()
    if err != nil {
        return nil, fmt.Errorf("failed to build logger: %w", err)
    }
    
    return logger, nil
}
```

### 3.6.2 结构化日志

```go
// 业务日志记录
func LogShipEvent(logger *zap.Logger, event ShipEvent) {
    logger.Info("Ship event occurred",
        zap.String("ship_id", event.ShipID),
        zap.String("event_type", event.Type),
        zap.Float64("position_x", event.Position.X),
        zap.Float64("position_y", event.Position.Y),
        zap.Float64("speed", event.Speed),
        zap.Time("timestamp", event.Timestamp),
    )
}

func LogLockOperation(logger *zap.Logger, operation LockOperation) {
    logger.Info("Lock operation executed",
        zap.String("lock_id", operation.LockID),
        zap.String("operation", operation.Type),
        zap.String("direction", operation.Direction),
        zap.String("operator", operation.Operator),
        zap.Time("timestamp", operation.Timestamp),
    )
}
```

## 3.7 错误处理实现

### 3.7.1 错误类型定义

```go
// 应用错误类型
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Cause   error  `json:"cause,omitempty"`
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 错误代码定义
const (
    ErrCodeInvalidInput     = "INVALID_INPUT"
    ErrCodeDeviceNotFound   = "DEVICE_NOT_FOUND"
    ErrCodeDeviceOffline    = "DEVICE_OFFLINE"
    ErrCodeDatabaseError    = "DATABASE_ERROR"
    ErrCodeNetworkError     = "NETWORK_ERROR"
    ErrCodePermissionDenied = "PERMISSION_DENIED"
)

// 错误工厂函数
func NewInvalidInputError(message string) *AppError {
    return &AppError{
        Code:    ErrCodeInvalidInput,
        Message: message,
    }
}

func NewDeviceNotFoundError(deviceID string) *AppError {
    return &AppError{
        Code:    ErrCodeDeviceNotFound,
        Message: fmt.Sprintf("Device %s not found", deviceID),
    }
}
```

### 3.7.2 错误处理中间件

```go
// HTTP错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            // 转换为应用错误
            if appErr, ok := err.(*AppError); ok {
                c.JSON(getStatusCode(appErr.Code), appErr)
                return
            }
            
            // 默认错误处理
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Internal server error",
            })
        }
    }
}

func getStatusCode(errorCode string) int {
    switch errorCode {
    case ErrCodeInvalidInput:
        return http.StatusBadRequest
    case ErrCodeDeviceNotFound:
        return http.StatusNotFound
    case ErrCodePermissionDenied:
        return http.StatusForbidden
    default:
        return http.StatusInternalServerError
    }
}
```

## 3.8 性能优化实现

### 3.8.1 连接池管理

```go
// 数据库连接池
type ConnectionPool struct {
    pool *sql.DB
    config DBConfig
}

func NewConnectionPool(config DBConfig) (*ConnectionPool, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        config.Username, config.Password, config.Host, config.Port, config.Database)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // 配置连接池
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return &ConnectionPool{
        pool:   db,
        config: config,
    }, nil
}
```

### 3.8.2 缓存实现

```go
// 内存缓存
type MemoryCache struct {
    cache map[string]interface{}
    mu    sync.RWMutex
    ttl   map[string]time.Time
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{
        cache: make(map[string]interface{}),
        ttl:   make(map[string]time.Time),
    }
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = value
    if ttl > 0 {
        c.ttl[key] = time.Now().Add(ttl)
    }
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    value, exists := c.cache[key]
    if !exists {
        return nil, false
    }
    
    // 检查TTL
    if expiry, hasTTL := c.ttl[key]; hasTTL && time.Now().After(expiry) {
        delete(c.cache, key)
        delete(c.ttl, key)
        return nil, false
    }
    
    return value, true
}
```

## 3.9 监控和指标

### 3.9.1 性能指标收集

```go
// 性能指标
type Metrics struct {
    RequestCount    int64
    ErrorCount      int64
    ResponseTime    time.Duration
    ActiveConnections int64
}

type MetricsCollector struct {
    metrics *Metrics
    mu      sync.RWMutex
}

func (m *MetricsCollector) IncrementRequestCount() {
    atomic.AddInt64(&m.metrics.RequestCount, 1)
}

func (m *MetricsCollector) IncrementErrorCount() {
    atomic.AddInt64(&m.metrics.ErrorCount, 1)
}

func (m *MetricsCollector) RecordResponseTime(duration time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.metrics.ResponseTime = duration
}

func (m *MetricsCollector) GetMetrics() *Metrics {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    return &Metrics{
        RequestCount:      atomic.LoadInt64(&m.metrics.RequestCount),
        ErrorCount:        atomic.LoadInt64(&m.metrics.ErrorCount),
        ResponseTime:      m.metrics.ResponseTime,
        ActiveConnections: atomic.LoadInt64(&m.metrics.ActiveConnections),
    }
}
```

### 3.9.2 健康检查

```go
// 健康检查
type HealthChecker struct {
    db     *sql.DB
    nats   *nats.Conn
    devices map[string]DeviceController
}

func (h *HealthChecker) CheckHealth() *HealthStatus {
    status := &HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Checks:    make(map[string]CheckResult),
    }
    
    // 检查数据库
    if err := h.db.Ping(); err != nil {
        status.Status = "unhealthy"
        status.Checks["database"] = CheckResult{
            Status:  "failed",
            Message: err.Error(),
        }
    } else {
        status.Checks["database"] = CheckResult{
            Status:  "passed",
            Message: "Database connection is healthy",
        }
    }
    
    // 检查NATS连接
    if h.nats.IsConnected() {
        status.Checks["nats"] = CheckResult{
            Status:  "passed",
            Message: "NATS connection is healthy",
        }
    } else {
        status.Status = "unhealthy"
        status.Checks["nats"] = CheckResult{
            Status:  "failed",
            Message: "NATS connection is not healthy",
        }
    }
    
    // 检查设备状态
    for deviceID, device := range h.devices {
        if deviceStatus, err := device.GetStatus(); err != nil {
            status.Status = "unhealthy"
            status.Checks[deviceID] = CheckResult{
                Status:  "failed",
                Message: err.Error(),
            }
        } else {
            status.Checks[deviceID] = CheckResult{
                Status:  "passed",
                Message: "Device is healthy",
            }
        }
    }
    
    return status
}
```
