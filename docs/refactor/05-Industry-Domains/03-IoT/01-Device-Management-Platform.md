# 01-设备管理平台 (Device Management Platform)

## 目录

1. [理论基础](#1-理论基础)
2. [形式化定义](#2-形式化定义)
3. [架构设计](#3-架构设计)
4. [Go语言实现](#4-go语言实现)
5. [安全机制](#5-安全机制)

## 1. 理论基础

### 1.1 设备管理平台定义

设备管理平台是物联网系统的核心组件，负责设备的注册、监控、配置、固件更新和生命周期管理。

**形式化定义**：

```math
设备管理平台定义为七元组：
DMP = (D, R, M, C, F, S, L)

其中：
- D: 设备集合，D = \{d_1, d_2, ..., d_n\}
- R: 注册函数，R: D \rightarrow \mathbb{B}
- M: 监控函数，M: D \rightarrow Status
- C: 配置函数，C: D \times Config \rightarrow \mathbb{B}
- F: 固件管理，F: D \times Firmware \rightarrow \mathbb{B}
- S: 安全机制，S: D \rightarrow SecurityLevel
- L: 生命周期，L: D \rightarrow LifecycleState
```

### 1.2 核心功能

1. **设备注册**: 设备身份验证和注册
2. **设备监控**: 实时状态监控和告警
3. **设备配置**: 远程配置和参数设置
4. **固件管理**: OTA固件更新
5. **安全管理**: 设备认证和授权

## 2. 形式化定义

### 2.1 设备模型

```math
设备定义为六元组：
Device = (ID, Type, Status, Config, Firmware, Metadata)

其中：
- ID: 设备唯一标识
- Type: 设备类型
- Status: 设备状态 (在线、离线、故障等)
- Config: 设备配置
- Firmware: 固件信息
- Metadata: 设备元数据

设备状态转换：
Status: \{Online, Offline, Fault, Maintenance, Updating\}
```

### 2.2 设备注册协议

```math
注册协议定义为：
Register(d, credentials) = 
  if VerifyCredentials(d, credentials) then
    if ValidateDevice(d) then
      RegisterDevice(d) \land AssignID(d)
    else
      RejectRegistration(d)
  else
    AuthenticationFailed(d)
```

## 3. 架构设计

### 3.1 分层架构

```go
// DeviceManagementPlatform 设备管理平台
type DeviceManagementPlatform struct {
    // 接入层 - 设备连接和协议适配
    AccessLayer *AccessLayer
    // 管理层 - 设备管理和业务逻辑
    ManagementLayer *ManagementLayer
    // 服务层 - 核心服务
    ServiceLayer *ServiceLayer
    // 数据层 - 数据存储和访问
    DataLayer *DataLayer
}

// AccessLayer 接入层
type AccessLayer struct {
    MQTTHandler    *MQTTHandler
    HTTPHandler    *HTTPHandler
    CoAPHandler    *CoAPHandler
    ProtocolAdapter *ProtocolAdapter
}

// ManagementLayer 管理层
type ManagementLayer struct {
    DeviceManager    *DeviceManager
    ConfigManager    *ConfigManager
    FirmwareManager  *FirmwareManager
    SecurityManager  *SecurityManager
}

// ServiceLayer 服务层
type ServiceLayer struct {
    RegistrationService *RegistrationService
    MonitoringService   *MonitoringService
    ConfigurationService *ConfigurationService
    UpdateService       *UpdateService
}

// DataLayer 数据层
type DataLayer struct {
    DeviceDB      *DeviceDatabase
    ConfigDB      *ConfigDatabase
    FirmwareDB    *FirmwareDatabase
    LogDB         *LogDatabase
}
```

## 4. Go语言实现

### 4.1 设备定义

```go
// Device 设备定义
type Device struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        DeviceType        `json:"type"`
    Status      DeviceStatus      `json:"status"`
    Config      *DeviceConfig     `json:"config"`
    Firmware    *FirmwareInfo     `json:"firmware"`
    Metadata    map[string]string `json:"metadata"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    LastSeen    time.Time         `json:"last_seen"`
}

// DeviceType 设备类型
type DeviceType int

const (
    DeviceTypeSensor DeviceType = iota
    DeviceTypeActuator
    DeviceTypeGateway
    DeviceTypeController
)

// DeviceStatus 设备状态
type DeviceStatus int

const (
    DeviceStatusOffline DeviceStatus = iota
    DeviceStatusOnline
    DeviceStatusFault
    DeviceStatusMaintenance
    DeviceStatusUpdating
)

// DeviceConfig 设备配置
type DeviceConfig struct {
    Parameters map[string]interface{} `json:"parameters"`
    Settings   map[string]interface{} `json:"settings"`
    Version    string                 `json:"version"`
}

// FirmwareInfo 固件信息
type FirmwareInfo struct {
    Version     string    `json:"version"`
    BuildDate   time.Time `json:"build_date"`
    Checksum    string    `json:"checksum"`
    Size        int64     `json:"size"`
    URL         string    `json:"url"`
}
```

### 4.2 设备管理器

```go
// DeviceManager 设备管理器
type DeviceManager struct {
    devices     map[string]*Device
    db          *sql.DB
    cache       *Cache
    validator   *DeviceValidator
    auditor     *Auditor
    mu          sync.RWMutex
}

// NewDeviceManager 创建设备管理器
func NewDeviceManager(db *sql.DB, cache *Cache) *DeviceManager {
    return &DeviceManager{
        devices:   make(map[string]*Device),
        db:        db,
        cache:     cache,
        validator: NewDeviceValidator(),
        auditor:   NewAuditor(),
    }
}

// RegisterDevice 注册设备
func (dm *DeviceManager) RegisterDevice(device *Device, credentials *Credentials) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    // 验证设备信息
    if err := dm.validator.ValidateDevice(device); err != nil {
        return err
    }
    
    // 验证凭据
    if err := dm.validateCredentials(device, credentials); err != nil {
        return err
    }
    
    // 生成设备ID
    device.ID = dm.generateDeviceID()
    device.CreatedAt = time.Now()
    device.UpdatedAt = time.Now()
    device.Status = DeviceStatusOffline
    
    // 保存到数据库
    if err := dm.saveDevice(device); err != nil {
        return err
    }
    
    // 更新缓存
    dm.cache.Set(device.ID, device, 24*time.Hour)
    
    // 记录审计日志
    dm.auditor.LogDeviceRegistration(device)
    
    return nil
}

// GetDevice 获取设备
func (dm *DeviceManager) GetDevice(deviceID string) (*Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    // 先从缓存获取
    if cached, found := dm.cache.Get(deviceID); found {
        return cached.(*Device), nil
    }
    
    // 从数据库获取
    device, err := dm.loadDevice(deviceID)
    if err != nil {
        return nil, err
    }
    
    // 更新缓存
    dm.cache.Set(deviceID, device, 24*time.Hour)
    
    return device, nil
}

// UpdateDeviceStatus 更新设备状态
func (dm *DeviceManager) UpdateDeviceStatus(deviceID string, status DeviceStatus) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    device, err := dm.GetDevice(deviceID)
    if err != nil {
        return err
    }
    
    device.Status = status
    device.UpdatedAt = time.Now()
    device.LastSeen = time.Now()
    
    // 保存到数据库
    if err := dm.saveDevice(device); err != nil {
        return err
    }
    
    // 更新缓存
    dm.cache.Set(deviceID, device, 24*time.Hour)
    
    // 记录审计日志
    dm.auditor.LogDeviceStatusChange(device, status)
    
    return nil
}

// saveDevice 保存设备到数据库
func (dm *DeviceManager) saveDevice(device *Device) error {
    query := `
        INSERT INTO devices (id, name, type, status, config, firmware, metadata, created_at, updated_at, last_seen)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        name = VALUES(name),
        type = VALUES(type),
        status = VALUES(status),
        config = VALUES(config),
        firmware = VALUES(firmware),
        metadata = VALUES(metadata),
        updated_at = VALUES(updated_at),
        last_seen = VALUES(last_seen)
    `
    
    configJSON, err := json.Marshal(device.Config)
    if err != nil {
        return err
    }
    
    firmwareJSON, err := json.Marshal(device.Firmware)
    if err != nil {
        return err
    }
    
    metadataJSON, err := json.Marshal(device.Metadata)
    if err != nil {
        return err
    }
    
    _, err = dm.db.Exec(query,
        device.ID,
        device.Name,
        device.Type,
        device.Status,
        configJSON,
        firmwareJSON,
        metadataJSON,
        device.CreatedAt,
        device.UpdatedAt,
        device.LastSeen,
    )
    
    return err
}

// loadDevice 从数据库加载设备
func (dm *DeviceManager) loadDevice(deviceID string) (*Device, error) {
    query := `
        SELECT id, name, type, status, config, firmware, metadata, created_at, updated_at, last_seen
        FROM devices WHERE id = ?
    `
    
    var device Device
    var configJSON, firmwareJSON, metadataJSON []byte
    
    err := dm.db.QueryRow(query, deviceID).Scan(
        &device.ID,
        &device.Name,
        &device.Type,
        &device.Status,
        &configJSON,
        &firmwareJSON,
        &metadataJSON,
        &device.CreatedAt,
        &device.UpdatedAt,
        &device.LastSeen,
    )
    
    if err != nil {
        return nil, err
    }
    
    // 解析配置
    if err := json.Unmarshal(configJSON, &device.Config); err != nil {
        return nil, err
    }
    
    // 解析固件信息
    if err := json.Unmarshal(firmwareJSON, &device.Firmware); err != nil {
        return nil, err
    }
    
    // 解析元数据
    if err := json.Unmarshal(metadataJSON, &device.Metadata); err != nil {
        return nil, err
    }
    
    return &device, nil
}
```

### 4.3 配置管理器

```go
// ConfigManager 配置管理器
type ConfigManager struct {
    db          *sql.DB
    cache       *Cache
    validator   *ConfigValidator
    auditor     *Auditor
    mu          sync.RWMutex
}

// NewConfigManager 创建配置管理器
func NewConfigManager(db *sql.DB, cache *Cache) *ConfigManager {
    return &ConfigManager{
        db:        db,
        cache:     cache,
        validator: NewConfigValidator(),
        auditor:   NewAuditor(),
    }
}

// UpdateDeviceConfig 更新设备配置
func (cm *ConfigManager) UpdateDeviceConfig(deviceID string, config *DeviceConfig) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    // 验证配置
    if err := cm.validator.ValidateConfig(config); err != nil {
        return err
    }
    
    // 保存配置
    if err := cm.saveConfig(deviceID, config); err != nil {
        return err
    }
    
    // 更新缓存
    cm.cache.Set(fmt.Sprintf("config_%s", deviceID), config, 1*time.Hour)
    
    // 记录审计日志
    cm.auditor.LogConfigUpdate(deviceID, config)
    
    return nil
}

// GetDeviceConfig 获取设备配置
func (cm *ConfigManager) GetDeviceConfig(deviceID string) (*DeviceConfig, error) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    // 先从缓存获取
    cacheKey := fmt.Sprintf("config_%s", deviceID)
    if cached, found := cm.cache.Get(cacheKey); found {
        return cached.(*DeviceConfig), nil
    }
    
    // 从数据库获取
    config, err := cm.loadConfig(deviceID)
    if err != nil {
        return nil, err
    }
    
    // 更新缓存
    cm.cache.Set(cacheKey, config, 1*time.Hour)
    
    return config, nil
}

// saveConfig 保存配置到数据库
func (cm *ConfigManager) saveConfig(deviceID string, config *DeviceConfig) error {
    query := `
        INSERT INTO device_configs (device_id, parameters, settings, version, updated_at)
        VALUES (?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        parameters = VALUES(parameters),
        settings = VALUES(settings),
        version = VALUES(version),
        updated_at = VALUES(updated_at)
    `
    
    parametersJSON, err := json.Marshal(config.Parameters)
    if err != nil {
        return err
    }
    
    settingsJSON, err := json.Marshal(config.Settings)
    if err != nil {
        return err
    }
    
    _, err = cm.db.Exec(query,
        deviceID,
        parametersJSON,
        settingsJSON,
        config.Version,
        time.Now(),
    )
    
    return err
}

// loadConfig 从数据库加载配置
func (cm *ConfigManager) loadConfig(deviceID string) (*DeviceConfig, error) {
    query := `
        SELECT parameters, settings, version
        FROM device_configs WHERE device_id = ?
    `
    
    var config DeviceConfig
    var parametersJSON, settingsJSON []byte
    
    err := cm.db.QueryRow(query, deviceID).Scan(
        &parametersJSON,
        &settingsJSON,
        &config.Version,
    )
    
    if err != nil {
        return nil, err
    }
    
    // 解析参数
    if err := json.Unmarshal(parametersJSON, &config.Parameters); err != nil {
        return nil, err
    }
    
    // 解析设置
    if err := json.Unmarshal(settingsJSON, &config.Settings); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 4.4 固件管理器

```go
// FirmwareManager 固件管理器
type FirmwareManager struct {
    db          *sql.DB
    storage     *FirmwareStorage
    validator   *FirmwareValidator
    auditor     *Auditor
    mu          sync.RWMutex
}

// NewFirmwareManager 创建固件管理器
func NewFirmwareManager(db *sql.DB, storage *FirmwareStorage) *FirmwareManager {
    return &FirmwareManager{
        db:        db,
        storage:   storage,
        validator: NewFirmwareValidator(),
        auditor:   NewAuditor(),
    }
}

// UploadFirmware 上传固件
func (fm *FirmwareManager) UploadFirmware(firmware *Firmware) error {
    fm.mu.Lock()
    defer fm.mu.Unlock()
    
    // 验证固件
    if err := fm.validator.ValidateFirmware(firmware); err != nil {
        return err
    }
    
    // 计算校验和
    checksum, err := fm.calculateChecksum(firmware.Data)
    if err != nil {
        return err
    }
    firmware.Checksum = checksum
    
    // 保存到存储
    if err := fm.storage.Save(firmware); err != nil {
        return err
    }
    
    // 保存到数据库
    if err := fm.saveFirmware(firmware); err != nil {
        return err
    }
    
    // 记录审计日志
    fm.auditor.LogFirmwareUpload(firmware)
    
    return nil
}

// UpdateDeviceFirmware 更新设备固件
func (fm *FirmwareManager) UpdateDeviceFirmware(deviceID string, firmwareVersion string) error {
    fm.mu.Lock()
    defer fm.mu.Unlock()
    
    // 获取固件
    firmware, err := fm.GetFirmware(firmwareVersion)
    if err != nil {
        return err
    }
    
    // 创建更新任务
    updateTask := &FirmwareUpdateTask{
        DeviceID:        deviceID,
        FirmwareVersion: firmwareVersion,
        Status:          UpdateStatusPending,
        CreatedAt:       time.Now(),
    }
    
    // 保存更新任务
    if err := fm.saveUpdateTask(updateTask); err != nil {
        return err
    }
    
    // 记录审计日志
    fm.auditor.LogFirmwareUpdate(deviceID, firmwareVersion)
    
    return nil
}

// GetFirmware 获取固件
func (fm *FirmwareManager) GetFirmware(version string) (*Firmware, error) {
    fm.mu.RLock()
    defer fm.mu.RUnlock()
    
    // 从数据库获取固件信息
    firmware, err := fm.loadFirmware(version)
    if err != nil {
        return nil, err
    }
    
    // 从存储获取固件数据
    data, err := fm.storage.Load(firmware.ID)
    if err != nil {
        return nil, err
    }
    
    firmware.Data = data
    return firmware, nil
}

// calculateChecksum 计算校验和
func (fm *FirmwareManager) calculateChecksum(data []byte) (string, error) {
    hash := sha256.New()
    hash.Write(data)
    return hex.EncodeToString(hash.Sum(nil)), nil
}

// saveFirmware 保存固件信息到数据库
func (fm *FirmwareManager) saveFirmware(firmware *Firmware) error {
    query := `
        INSERT INTO firmwares (id, version, build_date, checksum, size, url, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        build_date = VALUES(build_date),
        checksum = VALUES(checksum),
        size = VALUES(size),
        url = VALUES(url)
    `
    
    _, err := fm.db.Exec(query,
        firmware.ID,
        firmware.Version,
        firmware.BuildDate,
        firmware.Checksum,
        firmware.Size,
        firmware.URL,
        firmware.CreatedAt,
    )
    
    return err
}

// loadFirmware 从数据库加载固件信息
func (fm *FirmwareManager) loadFirmware(version string) (*Firmware, error) {
    query := `
        SELECT id, version, build_date, checksum, size, url, created_at
        FROM firmwares WHERE version = ?
    `
    
    var firmware Firmware
    
    err := fm.db.QueryRow(query, version).Scan(
        &firmware.ID,
        &firmware.Version,
        &firmware.BuildDate,
        &firmware.Checksum,
        &firmware.Size,
        &firmware.URL,
        &firmware.CreatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &firmware, nil
}
```

### 4.5 监控服务

```go
// MonitoringService 监控服务
type MonitoringService struct {
    deviceManager *DeviceManager
    alertManager  *AlertManager
    metrics       *MetricsCollector
    interval      time.Duration
    stop          chan struct{}
}

// NewMonitoringService 创建监控服务
func NewMonitoringService(deviceManager *DeviceManager, interval time.Duration) *MonitoringService {
    return &MonitoringService{
        deviceManager: deviceManager,
        alertManager:  NewAlertManager(),
        metrics:       NewMetricsCollector(),
        interval:      interval,
        stop:          make(chan struct{}),
    }
}

// Start 启动监控
func (ms *MonitoringService) Start() {
    ticker := time.NewTicker(ms.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            ms.checkDevices()
        case <-ms.stop:
            return
        }
    }
}

// Stop 停止监控
func (ms *MonitoringService) Stop() {
    close(ms.stop)
}

// checkDevices 检查设备状态
func (ms *MonitoringService) checkDevices() {
    devices, err := ms.deviceManager.GetAllDevices()
    if err != nil {
        log.Printf("Failed to get devices: %v", err)
        return
    }
    
    for _, device := range devices {
        // 检查设备是否超时
        if ms.isDeviceTimeout(device) {
            ms.handleDeviceTimeout(device)
        }
        
        // 检查设备状态
        if device.Status == DeviceStatusFault {
            ms.handleDeviceFault(device)
        }
        
        // 收集指标
        ms.metrics.CollectDeviceMetrics(device)
    }
}

// isDeviceTimeout 检查设备是否超时
func (ms *MonitoringService) isDeviceTimeout(device *Device) bool {
    timeout := 5 * time.Minute
    return time.Since(device.LastSeen) > timeout
}

// handleDeviceTimeout 处理设备超时
func (ms *MonitoringService) handleDeviceTimeout(device *Device) {
    // 更新设备状态为离线
    if device.Status == DeviceStatusOnline {
        ms.deviceManager.UpdateDeviceStatus(device.ID, DeviceStatusOffline)
    }
    
    // 发送告警
    ms.alertManager.SendAlert(&Alert{
        Type:      AlertTypeDeviceTimeout,
        DeviceID:  device.ID,
        Message:   fmt.Sprintf("Device %s is timeout", device.Name),
        Timestamp: time.Now(),
    })
}

// handleDeviceFault 处理设备故障
func (ms *MonitoringService) handleDeviceFault(device *Device) {
    // 发送告警
    ms.alertManager.SendAlert(&Alert{
        Type:      AlertTypeDeviceFault,
        DeviceID:  device.ID,
        Message:   fmt.Sprintf("Device %s is in fault state", device.Name),
        Timestamp: time.Now(),
    })
}
```

## 5. 安全机制

### 5.1 设备认证

```go
// SecurityManager 安全管理器
type SecurityManager struct {
    certManager  *CertificateManager
    keyManager   *KeyManager
    authManager  *AuthManager
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager() *SecurityManager {
    return &SecurityManager{
        certManager: NewCertificateManager(),
        keyManager:  NewKeyManager(),
        authManager: NewAuthManager(),
    }
}

// AuthenticateDevice 设备认证
func (sm *SecurityManager) AuthenticateDevice(deviceID string, credentials *Credentials) error {
    // 验证设备证书
    if err := sm.certManager.ValidateCertificate(deviceID, credentials.Certificate); err != nil {
        return err
    }
    
    // 验证设备密钥
    if err := sm.keyManager.ValidateKey(deviceID, credentials.Key); err != nil {
        return err
    }
    
    // 验证设备令牌
    if err := sm.authManager.ValidateToken(deviceID, credentials.Token); err != nil {
        return err
    }
    
    return nil
}

// GenerateDeviceCredentials 生成设备凭据
func (sm *SecurityManager) GenerateDeviceCredentials(deviceID string) (*Credentials, error) {
    // 生成证书
    cert, err := sm.certManager.GenerateCertificate(deviceID)
    if err != nil {
        return nil, err
    }
    
    // 生成密钥
    key, err := sm.keyManager.GenerateKey(deviceID)
    if err != nil {
        return nil, err
    }
    
    // 生成令牌
    token, err := sm.authManager.GenerateToken(deviceID)
    if err != nil {
        return nil, err
    }
    
    return &Credentials{
        DeviceID:   deviceID,
        Certificate: cert,
        Key:        key,
        Token:      token,
    }, nil
}
```

## 总结

设备管理平台是物联网系统的核心组件，提供设备注册、监控、配置、固件更新等关键功能。本文档提供了完整的理论基础、形式化定义和Go语言实现。

### 关键要点

1. **设备生命周期管理**: 完整的设备注册到注销流程
2. **实时监控**: 设备状态监控和告警机制
3. **安全认证**: 多层次的安全保护
4. **配置管理**: 远程配置和参数设置
5. **固件更新**: OTA固件更新机制

### 扩展阅读

- [数据采集系统](./02-Data-Collection-System.md)
- [边缘计算](./03-Edge-Computing.md)
- [传感器网络](./04-Sensor-Network.md)
