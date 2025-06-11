# 5. 系统配置文档

## 5.1 配置概述

### 5.1.1 配置管理架构

葛洲坝船闸导航系统采用分层配置管理架构：

```text
┌─────────────────────────────────────────────────────────┐
│                    应用配置                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 系统配置     │ │ 业务配置     │ │ 设备配置     │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    配置加载器                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 文件加载     │ │ 环境变量     │ │ 远程配置     │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    配置存储                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ YAML文件     │ │ JSON文件     │ │ 数据库       │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
```

### 5.1.2 配置分类

1. **系统级配置**
   - 日志配置
   - 网络配置
   - 数据库配置
   - 安全配置

2. **业务级配置**
   - 船闸参数
   - 船舶参数
   - 业务规则
   - 阈值设置

3. **设备级配置**
   - 雷达配置
   - 云台配置
   - LED配置
   - 开关量配置

## 5.2 配置文件结构

### 5.2.1 主配置文件 (config.yaml)

```yaml
# 系统配置
system:
  name: "葛洲坝船闸导航系统"
  version: "1.0.0"
  environment: "production"
  debug: false

# 日志配置
logging:
  level: "info"
  format: "json"
  output_path: "./logs"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true

# 网络配置
network:
  nats:
    url: "nats://localhost:4222"
    cluster_id: "navlock-cluster"
    client_id: "gezhouba-navlock"
  http:
    port: 8080
    host: "0.0.0.0"
    timeout: 30s

# 数据库配置
database:
  local_db:
    driver: "sqlite3"
    dsn: "./data/local.db"
    max_open_conns: 10
    max_idle_conns: 5
    conn_max_lifetime: 300s
  
  produce_db:
    driver: "mysql"
    host: "localhost"
    port: 3306
    username: "navlock"
    password: "password"
    database: "navlock_prod"
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s
  
  web_db:
    driver: "mysql"
    host: "localhost"
    port: 3306
    username: "navlock_web"
    password: "password"
    database: "navlock_web"
    max_open_conns: 10
    max_idle_conns: 3
    conn_max_lifetime: 300s

# 安全配置
security:
  jwt_secret: "your-jwt-secret"
  token_expiry: 24h
  cors_origins: ["http://localhost:3000"]
  rate_limit: 100

# 业务配置
business:
  lock:
    id: "gezhouba-1"
    name: "葛洲坝1号船闸"
    type: "navigation"
    capacity: 10000
    max_ships: 5
  
  ship:
    max_speed: 15.0
    min_speed: 0.1
    max_length: 300.0
    max_width: 50.0
  
  zones:
    no_stop_zones:
      - name: "上游禁停区"
        points:
          - {x: 0, y: 0}
          - {x: 100, y: 0}
          - {x: 100, y: 50}
          - {x: 0, y: 50}
      - name: "下游禁停区"
        points:
          - {x: 200, y: 0}
          - {x: 300, y: 0}
          - {x: 300, y: 50}
          - {x: 200, y: 50}
    
    speed_limit_zones:
      - name: "上游限速区"
        max_speed: 8.0
        points:
          - {x: 0, y: 0}
          - {x: 150, y: 0}
          - {x: 150, y: 50}
          - {x: 0, y: 50}

# 设备配置
devices:
  radar:
    - id: "radar-upstream"
      name: "上游雷达"
      type: "radar"
      address: "192.168.1.100"
      port: 8001
      protocol: "tcp"
      scan_range: 3000.0
      scan_angle: 120.0
      update_rate: 10
    
    - id: "radar-downstream"
      name: "下游雷达"
      type: "radar"
      address: "192.168.1.101"
      port: 8002
      protocol: "tcp"
      scan_range: 3000.0
      scan_angle: 120.0
      update_rate: 10
  
  ptz:
    - id: "ptz-upstream"
      name: "上游云台"
      type: "ptz"
      address: "192.168.1.200"
      port: 8003
      protocol: "pelco-d"
      pan_speed: 100
      tilt_speed: 100
      zoom_speed: 50
    
    - id: "ptz-downstream"
      name: "下游云台"
      type: "ptz"
      address: "192.168.1.201"
      port: 8004
      protocol: "pelco-d"
      pan_speed: 100
      tilt_speed: 100
      zoom_speed: 50
  
  led:
    - id: "led-upstream"
      name: "上游LED"
      type: "led"
      address: "192.168.1.300"
      port: 8005
      protocol: "custom"
      width: 64
      height: 32
      brightness: 255
    
    - id: "led-downstream"
      name: "下游LED"
      type: "led"
      address: "192.168.1.301"
      port: 8006
      protocol: "custom"
      width: 64
      height: 32
      brightness: 255
  
  digital_io:
    - id: "io-controller"
      name: "开关量控制器"
      type: "digital_io"
      address: "192.168.1.400"
      port: 8007
      protocol: "modbus"
      input_channels: 16
      output_channels: 8
```

### 5.2.2 环境特定配置

#### 开发环境 (config.dev.yaml)

```yaml
system:
  environment: "development"
  debug: true

logging:
  level: "debug"
  output_path: "stdout"

database:
  local_db:
    dsn: ":memory:"
  
  produce_db:
    host: "localhost"
    username: "dev_user"
    password: "dev_password"

devices:
  radar:
    - address: "127.0.0.1"
      port: 9001
```

#### 测试环境 (config.test.yaml)

```yaml
system:
  environment: "testing"
  debug: true

logging:
  level: "debug"
  output_path: "./test_logs"

database:
  local_db:
    dsn: "./test_data/test.db"
  
  produce_db:
    host: "test-db"
    username: "test_user"
    password: "test_password"

devices:
  radar:
    - address: "test-radar"
      port: 9001
```

#### 生产环境 (config.prod.yaml)

```yaml
system:
  environment: "production"
  debug: false

logging:
  level: "info"
  output_path: "/var/log/navlock"

database:
  local_db:
    dsn: "/var/lib/navlock/local.db"
  
  produce_db:
    host: "prod-db-cluster"
    username: "prod_user"
    password: "${DB_PASSWORD}"

devices:
  radar:
    - address: "10.0.1.100"
      port: 8001
```

## 5.3 配置加载机制

### 5.3.1 配置加载器实现

```go
// 配置加载器
type ConfigLoader struct {
    viper *viper.Viper
}

func NewConfigLoader() *ConfigLoader {
    v := viper.New()
    
    // 设置默认值
    v.SetDefault("system.name", "船闸导航系统")
    v.SetDefault("system.version", "1.0.0")
    v.SetDefault("system.environment", "development")
    v.SetDefault("system.debug", false)
    
    v.SetDefault("logging.level", "info")
    v.SetDefault("logging.format", "json")
    v.SetDefault("logging.output_path", "./logs")
    
    v.SetDefault("network.http.port", 8080)
    v.SetDefault("network.http.host", "0.0.0.0")
    
    return &ConfigLoader{viper: v}
}

func (l *ConfigLoader) Load() (*SystemConfig, error) {
    // 1. 读取主配置文件
    l.viper.SetConfigName("config")
    l.viper.SetConfigType("yaml")
    l.viper.AddConfigPath(".")
    l.viper.AddConfigPath("./config")
    l.viper.AddConfigPath("/etc/navlock")
    
    if err := l.viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    // 2. 读取环境特定配置
    env := l.viper.GetString("system.environment")
    if env != "" {
        l.viper.SetConfigName(fmt.Sprintf("config.%s", env))
        if err := l.viper.MergeInConfig(); err != nil {
            log.Printf("Warning: failed to merge environment config: %v", err)
        }
    }
    
    // 3. 读取环境变量
    l.viper.AutomaticEnv()
    l.viper.SetEnvPrefix("NAVLOCK")
    l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // 4. 解析配置
    var config SystemConfig
    if err := l.viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &config, nil
}
```

### 5.3.2 配置验证

```go
// 配置验证器
type ConfigValidator struct{}

func (v *ConfigValidator) Validate(config *SystemConfig) error {
    var errors []string
    
    // 验证系统配置
    if err := v.validateSystemConfig(config.System); err != nil {
        errors = append(errors, fmt.Sprintf("system: %v", err))
    }
    
    // 验证网络配置
    if err := v.validateNetworkConfig(config.Network); err != nil {
        errors = append(errors, fmt.Sprintf("network: %v", err))
    }
    
    // 验证数据库配置
    if err := v.validateDatabaseConfig(config.Database); err != nil {
        errors = append(errors, fmt.Sprintf("database: %v", err))
    }
    
    // 验证设备配置
    if err := v.validateDevicesConfig(config.Devices); err != nil {
        errors = append(errors, fmt.Sprintf("devices: %v", err))
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
    }
    
    return nil
}

func (v *ConfigValidator) validateSystemConfig(config SystemConfig) error {
    if config.Name == "" {
        return fmt.Errorf("system name is required")
    }
    
    if config.Version == "" {
        return fmt.Errorf("system version is required")
    }
    
    return nil
}

func (v *ConfigValidator) validateNetworkConfig(config NetworkConfig) error {
    if config.NATS.URL == "" {
        return fmt.Errorf("NATS URL is required")
    }
    
    if config.HTTP.Port <= 0 || config.HTTP.Port > 65535 {
        return fmt.Errorf("invalid HTTP port")
    }
    
    return nil
}

func (v *ConfigValidator) validateDatabaseConfig(config DatabaseConfig) error {
    if err := v.validateDBConfig(config.LocalDB); err != nil {
        return fmt.Errorf("local database: %v", err)
    }
    
    if err := v.validateDBConfig(config.ProduceDB); err != nil {
        return fmt.Errorf("produce database: %v", err)
    }
    
    return nil
}

func (v *ConfigValidator) validateDBConfig(config DBConfig) error {
    if config.Driver == "" {
        return fmt.Errorf("driver is required")
    }
    
    if config.DSN == "" {
        return fmt.Errorf("DSN is required")
    }
    
    return nil
}
```

## 5.4 配置热更新

### 5.4.1 配置监听器

```go
// 配置监听器
type ConfigWatcher struct {
    viper    *viper.Viper
    handlers []ConfigChangeHandler
    mu       sync.RWMutex
}

type ConfigChangeHandler func(key string, oldValue, newValue interface{})

func (w *ConfigWatcher) AddHandler(handler ConfigChangeHandler) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    w.handlers = append(w.handlers, handler)
}

func (w *ConfigWatcher) StartWatching() error {
    w.viper.WatchConfig()
    w.viper.OnConfigChange(func(e fsnotify.Event) {
        w.handleConfigChange(e)
    })
    
    return nil
}

func (w *ConfigWatcher) handleConfigChange(e fsnotify.Event) {
    log.Printf("Configuration file changed: %s", e.Name)
    
    // 重新加载配置
    if err := w.viper.ReadInConfig(); err != nil {
        log.Printf("Failed to reload config: %v", err)
        return
    }
    
    // 通知所有处理器
    w.mu.RLock()
    defer w.mu.RUnlock()
    
    for _, handler := range w.handlers {
        go handler("", nil, nil)
    }
}
```

### 5.4.2 动态配置更新

```go
// 动态配置管理器
type DynamicConfigManager struct {
    config  *SystemConfig
    watcher *ConfigWatcher
    mu      sync.RWMutex
}

func (m *DynamicConfigManager) UpdateConfig(key string, value interface{}) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // 更新内存中的配置
    if err := m.updateConfigValue(m.config, key, value); err != nil {
        return err
    }
    
    // 保存到文件
    if err := m.saveConfigToFile(); err != nil {
        return err
    }
    
    // 通知配置变更
    m.notifyConfigChange(key, value)
    
    return nil
}

func (m *DynamicConfigManager) updateConfigValue(config interface{}, key string, value interface{}) error {
    // 使用反射更新配置值
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    keys := strings.Split(key, ".")
    return m.setNestedValue(v, keys, value)
}

func (m *DynamicConfigManager) setNestedValue(v reflect.Value, keys []string, value interface{}) error {
    if len(keys) == 0 {
        return nil
    }
    
    key := keys[0]
    remainingKeys := keys[1:]
    
    switch v.Kind() {
    case reflect.Struct:
        field := v.FieldByName(key)
        if !field.IsValid() {
            return fmt.Errorf("field %s not found", key)
        }
        
        if len(remainingKeys) == 0 {
            return m.setValue(field, value)
        }
        
        return m.setNestedValue(field, remainingKeys, value)
        
    case reflect.Map:
        mapKey := reflect.ValueOf(key)
        mapValue := v.MapIndex(mapKey)
        
        if len(remainingKeys) == 0 {
            v.SetMapIndex(mapKey, reflect.ValueOf(value))
            return nil
        }
        
        if !mapValue.IsValid() {
            mapValue = reflect.New(v.Type().Elem()).Elem()
        }
        
        if err := m.setNestedValue(mapValue, remainingKeys, value); err != nil {
            return err
        }
        
        v.SetMapIndex(mapKey, mapValue)
        return nil
    }
    
    return fmt.Errorf("unsupported type: %v", v.Kind())
}
```

## 5.5 配置加密

### 5.5.1 敏感配置加密

```go
// 配置加密器
type ConfigEncryptor struct {
    key []byte
}

func NewConfigEncryptor(key string) *ConfigEncryptor {
    return &ConfigEncryptor{
        key: []byte(key),
    }
}

func (e *ConfigEncryptor) Encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
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

func (e *ConfigEncryptor) Decrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
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

### 5.5.2 加密配置示例

```yaml
# 敏感配置使用加密值
database:
  produce_db:
    password: "ENC(AES256_GCM,encrypted_password_here)"
  
  web_db:
    password: "ENC(AES256_GCM,encrypted_password_here)"

security:
  jwt_secret: "ENC(AES256_GCM,encrypted_jwt_secret_here)"
```

## 5.6 配置模板

### 5.6.1 配置模板生成器

```go
// 配置模板生成器
type ConfigTemplateGenerator struct {
    templatePath string
    outputPath   string
}

func (g *ConfigTemplateGenerator) GenerateTemplate() error {
    // 读取模板文件
    templateData, err := os.ReadFile(g.templatePath)
    if err != nil {
        return fmt.Errorf("failed to read template: %w", err)
    }
    
    // 解析模板
    tmpl, err := template.New("config").Parse(string(templateData))
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    // 准备模板数据
    data := g.prepareTemplateData()
    
    // 生成配置文件
    file, err := os.Create(g.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()
    
    if err := tmpl.Execute(file, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }
    
    return nil
}

func (g *ConfigTemplateGenerator) prepareTemplateData() map[string]interface{} {
    return map[string]interface{}{
        "SystemName":    "葛洲坝船闸导航系统",
        "SystemVersion": "1.0.0",
        "Environment":   "production",
        "DatabaseHost":  "localhost",
        "DatabasePort":  3306,
        "DatabaseName":  "navlock",
        "NATSUrl":       "nats://localhost:4222",
        "HTTPPort":      8080,
    }
}
```

### 5.6.2 配置模板示例

```yaml
# config.template.yaml
system:
  name: "{{.SystemName}}"
  version: "{{.SystemVersion}}"
  environment: "{{.Environment}}"
  debug: false

logging:
  level: "info"
  format: "json"
  output_path: "./logs"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true

network:
  nats:
    url: "{{.NATSUrl}}"
    cluster_id: "navlock-cluster"
    client_id: "gezhouba-navlock"
  http:
    port: {{.HTTPPort}}
    host: "0.0.0.0"
    timeout: 30s

database:
  local_db:
    driver: "sqlite3"
    dsn: "./data/local.db"
    max_open_conns: 10
    max_idle_conns: 5
    conn_max_lifetime: 300s
  
  produce_db:
    driver: "mysql"
    host: "{{.DatabaseHost}}"
    port: {{.DatabasePort}}
    username: "navlock"
    password: "{{.DatabasePassword}}"
    database: "{{.DatabaseName}}"
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s
```

## 5.7 配置管理最佳实践

### 5.7.1 配置分层原则

1. **默认配置**
   - 提供合理的默认值
   - 确保系统可以正常启动

2. **环境配置**
   - 根据环境调整配置
   - 避免硬编码环境特定值

3. **用户配置**
   - 允许用户自定义配置
   - 提供配置验证和提示

### 5.7.2 配置安全原则

1. **敏感信息加密**
   - 密码、密钥等敏感信息加密存储
   - 使用环境变量传递敏感信息

2. **访问控制**
   - 限制配置文件访问权限
   - 记录配置变更日志

3. **配置备份**
   - 定期备份配置文件
   - 版本控制配置文件

### 5.7.3 配置维护原则

1. **配置文档化**
   - 详细说明每个配置项的作用
   - 提供配置示例和说明

2. **配置验证**
   - 启动时验证配置有效性
   - 提供配置检查工具

3. **配置监控**
   - 监控配置变更
   - 告警异常配置

## 5.8 配置工具

### 5.8.1 配置检查工具

```go
// 配置检查工具
type ConfigChecker struct {
    config *SystemConfig
}

func (c *ConfigChecker) CheckAll() *CheckResult {
    result := &CheckResult{
        Checks: make(map[string]CheckItem),
    }
    
    // 检查系统配置
    result.Checks["system"] = c.checkSystemConfig()
    
    // 检查网络配置
    result.Checks["network"] = c.checkNetworkConfig()
    
    // 检查数据库配置
    result.Checks["database"] = c.checkDatabaseConfig()
    
    // 检查设备配置
    result.Checks["devices"] = c.checkDevicesConfig()
    
    // 计算总体结果
    result.CalculateOverallResult()
    
    return result
}

func (c *ConfigChecker) checkSystemConfig() CheckItem {
    item := CheckItem{Name: "System Configuration"}
    
    if c.config.System.Name == "" {
        item.Status = "error"
        item.Message = "System name is required"
    } else if c.config.System.Version == "" {
        item.Status = "error"
        item.Message = "System version is required"
    } else {
        item.Status = "ok"
        item.Message = "System configuration is valid"
    }
    
    return item
}

func (c *ConfigChecker) checkNetworkConfig() CheckItem {
    item := CheckItem{Name: "Network Configuration"}
    
    if c.config.Network.NATS.URL == "" {
        item.Status = "error"
        item.Message = "NATS URL is required"
    } else if c.config.Network.HTTP.Port <= 0 {
        item.Status = "error"
        item.Message = "Invalid HTTP port"
    } else {
        item.Status = "ok"
        item.Message = "Network configuration is valid"
    }
    
    return item
}
```

### 5.8.2 配置生成工具

```go
// 配置生成工具
type ConfigGenerator struct {
    templatePath string
    outputPath   string
}

func (g *ConfigGenerator) GenerateFromTemplate(data map[string]interface{}) error {
    // 读取模板
    templateContent, err := os.ReadFile(g.templatePath)
    if err != nil {
        return fmt.Errorf("failed to read template: %w", err)
    }
    
    // 解析模板
    tmpl, err := template.New("config").Parse(string(templateContent))
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    // 生成配置
    file, err := os.Create(g.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()
    
    if err := tmpl.Execute(file, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }
    
    return nil
}

func (g *ConfigGenerator) GenerateFromDefaults() error {
    // 生成默认配置
    defaultConfig := &SystemConfig{
        System: SystemConfig{
            Name:        "船闸导航系统",
            Version:     "1.0.0",
            Environment: "development",
            Debug:       true,
        },
        Logging: LoggingConfig{
            Level:      "info",
            Format:     "json",
            OutputPath: "./logs",
        },
        Network: NetworkConfig{
            NATS: NATSConfig{
                URL:       "nats://localhost:4222",
                ClusterID: "navlock-cluster",
                ClientID:  "navlock-client",
            },
            HTTP: HTTPConfig{
                Port:    8080,
                Host:    "0.0.0.0",
                Timeout: 30 * time.Second,
            },
        },
    }
    
    // 序列化为YAML
    data, err := yaml.Marshal(defaultConfig)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }
    
    // 写入文件
    if err := os.WriteFile(g.outputPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }
    
    return nil
}
```
