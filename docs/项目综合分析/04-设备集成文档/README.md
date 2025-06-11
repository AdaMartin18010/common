# 4. 设备集成文档

## 4.1 设备集成概述

### 4.1.1 设备类型

葛洲坝船闸导航系统集成了以下主要设备类型：

1. **雷达设备**
   - 船舶位置检测
   - 速度测量
   - 轨迹跟踪

2. **云台设备 (PTZ)**
   - 摄像头控制
   - 角度调整
   - 图像采集

3. **LED显示屏**
   - 信息显示
   - 状态提示
   - 警告信息

4. **开关量设备**
   - 船闸状态检测
   - 设备开关控制
   - 信号采集

### 4.1.2 设备集成架构

```text
┌─────────────────────────────────────────────────────────┐
│                    应用层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 雷达服务     │ │ 云台服务     │ │ LED服务      │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    驱动层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 雷达驱动     │ │ 云台驱动     │ │ LED驱动      │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    通信层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ TCP/IP      │ │ 串口通信     │ │ Modbus      │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    设备层                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ 雷达设备     │ │ 云台设备     │ │ LED显示屏    │        │
│  └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────┘
```

## 4.2 雷达设备集成

### 4.2.1 雷达设备概述

雷达设备是船闸导航系统的核心传感器，主要用于：

- 船舶位置实时检测
- 船舶速度计算
- 船舶轨迹跟踪
- 禁停区域监测

### 4.2.2 雷达设备规格

```text
设备型号：XXX雷达
工作频率：24GHz
探测距离：0.5-3000米
角度范围：120度
精度：±0.1米
数据输出：TCP/IP
协议：自定义协议
```

### 4.2.3 雷达通信协议

```go
// 雷达数据包结构
type RadarPacket struct {
    Header     [4]byte  // 包头 0xAA 0x55 0xAA 0x55
    Length     uint16   // 数据长度
    Command    uint8    // 命令字
    Data       []byte   // 数据内容
    Checksum   uint8    // 校验和
}

// 雷达目标数据结构
type RadarTarget struct {
    TargetID   uint16   // 目标ID
    Distance   float32  // 距离 (米)
    Angle      float32  // 角度 (度)
    Speed      float32  // 速度 (m/s)
    Quality    uint8    // 信号质量
    Timestamp  uint32   // 时间戳
}

// 雷达命令定义
const (
    CMD_GET_STATUS    = 0x01  // 获取状态
    CMD_GET_TARGETS   = 0x02  // 获取目标
    CMD_SET_CONFIG    = 0x03  // 设置配置
    CMD_START_SCAN    = 0x04  // 开始扫描
    CMD_STOP_SCAN     = 0x05  // 停止扫描
)
```

### 4.2.4 雷达驱动实现

```go
// 雷达驱动接口
type RadarDriver interface {
    Connect() error
    Disconnect() error
    GetStatus() (*RadarStatus, error)
    GetTargets() ([]RadarTarget, error)
    SetConfig(config RadarConfig) error
    StartScan() error
    StopScan() error
}

// 雷达驱动实现
type RadarDriverImpl struct {
    connection net.Conn
    config     RadarConfig
    targets    []RadarTarget
    mu         sync.RWMutex
}

func (r *RadarDriverImpl) Connect() error {
    conn, err := net.Dial("tcp", r.config.Address)
    if err != nil {
        return fmt.Errorf("failed to connect to radar: %w", err)
    }
    
    r.connection = conn
    
    // 启动数据接收协程
    go r.receiveData()
    
    return nil
}

func (r *RadarDriverImpl) GetTargets() ([]RadarTarget, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    return r.targets, nil
}

func (r *RadarDriverImpl) receiveData() {
    buffer := make([]byte, 1024)
    
    for {
        n, err := r.connection.Read(buffer)
        if err != nil {
            log.Printf("Failed to read radar data: %v", err)
            break
        }
        
        // 解析数据包
        packet := r.parsePacket(buffer[:n])
        if packet != nil {
            r.processPacket(packet)
        }
    }
}

func (r *RadarDriverImpl) parsePacket(data []byte) *RadarPacket {
    // 查找包头
    headerIndex := bytes.Index(data, []byte{0xAA, 0x55, 0xAA, 0x55})
    if headerIndex == -1 {
        return nil
    }
    
    if len(data) < headerIndex+6 {
        return nil
    }
    
    // 解析数据包
    packet := &RadarPacket{}
    copy(packet.Header[:], data[headerIndex:headerIndex+4])
    packet.Length = binary.BigEndian.Uint16(data[headerIndex+4:headerIndex+6])
    packet.Command = data[headerIndex+6]
    
    if len(data) < headerIndex+7+int(packet.Length)+1 {
        return nil
    }
    
    packet.Data = data[headerIndex+7 : headerIndex+7+int(packet.Length)]
    packet.Checksum = data[headerIndex+7+int(packet.Length)]
    
    // 验证校验和
    if !r.validateChecksum(packet) {
        return nil
    }
    
    return packet
}
```

### 4.2.5 雷达数据处理

```go
// 雷达数据处理服务
type RadarDataProcessor struct {
    targets map[uint16]*RadarTarget
    mu      sync.RWMutex
}

func (p *RadarDataProcessor) ProcessTarget(target *RadarTarget) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // 更新目标信息
    p.targets[target.TargetID] = target
    
    // 判断是否为船舶
    if p.isShip(target) {
        p.processShipTarget(target)
    }
}

func (p *RadarDataProcessor) isShip(target *RadarTarget) bool {
    // 根据目标特征判断是否为船舶
    // 例如：速度范围、轨迹特征等
    return target.Speed > 0.1 && target.Speed < 20.0
}

func (p *RadarDataProcessor) processShipTarget(target *RadarTarget) {
    // 转换为船舶位置
    shipPosition := &ShipPosition{
        X:         p.calculateX(target),
        Y:         p.calculateY(target),
        Speed:     target.Speed,
        Direction: p.calculateDirection(target),
        Timestamp: time.Now(),
    }
    
    // 发布船舶位置事件
    p.publishShipPosition(shipPosition)
}

func (p *RadarDataProcessor) calculateX(target *RadarTarget) float64 {
    // 根据距离和角度计算X坐标
    return target.Distance * math.Cos(target.Angle*math.Pi/180)
}

func (p *RadarDataProcessor) calculateY(target *RadarTarget) float64 {
    // 根据距离和角度计算Y坐标
    return target.Distance * math.Sin(target.Angle*math.Pi/180)
}
```

## 4.3 云台设备集成

### 4.3.1 云台设备概述

云台设备用于控制摄像头的角度和位置，主要功能包括：

- 摄像头角度调整
- 自动跟踪船舶
- 图像采集
- 远程控制

### 4.3.2 云台设备规格

```text
设备型号：XXX云台
控制协议：Pelco-D/P
通信方式：RS485/RS232
旋转范围：水平360度，垂直-90到+90度
旋转速度：0.1-300度/秒
预置位：128个
```

### 4.3.3 云台控制协议

```go
// Pelco-D协议命令结构
type PelcoDCommand struct {
    Sync    byte   // 同步字节 0xFF
    Address byte   // 设备地址
    Command1 byte  // 命令1
    Command2 byte  // 命令2
    Data1   byte   // 数据1
    Data2   byte   // 数据2
    Checksum byte  // 校验和
}

// 云台命令定义
const (
    CMD_PAN_LEFT     = 0x04  // 左转
    CMD_PAN_RIGHT    = 0x02  // 右转
    CMD_TILT_UP      = 0x08  // 上转
    CMD_TILT_DOWN    = 0x10  // 下转
    CMD_ZOOM_IN      = 0x20  // 放大
    CMD_ZOOM_OUT     = 0x40  // 缩小
    CMD_STOP         = 0x00  // 停止
    CMD_PRESET_SET   = 0x03  // 设置预置位
    CMD_PRESET_CALL  = 0x07  // 调用预置位
)
```

### 4.3.4 云台驱动实现

```go
// 云台驱动接口
type PTZDriver interface {
    Connect() error
    Disconnect() error
    PanLeft(speed byte) error
    PanRight(speed byte) error
    TiltUp(speed byte) error
    TiltDown(speed byte) error
    ZoomIn() error
    ZoomOut() error
    Stop() error
    SetPreset(preset byte) error
    CallPreset(preset byte) error
    GetPosition() (*PTZPosition, error)
}

// 云台驱动实现
type PTZDriverImpl struct {
    connection net.Conn
    address    byte
    position   *PTZPosition
    mu         sync.RWMutex
}

func (p *PTZDriverImpl) PanLeft(speed byte) error {
    cmd := &PelcoDCommand{
        Sync:     0xFF,
        Address:  p.address,
        Command1: CMD_PAN_LEFT,
        Command2: 0x00,
        Data1:    speed,
        Data2:    0x00,
    }
    cmd.Checksum = p.calculateChecksum(cmd)
    
    return p.sendCommand(cmd)
}

func (p *PTZDriverImpl) sendCommand(cmd *PelcoDCommand) error {
    data := []byte{
        cmd.Sync,
        cmd.Address,
        cmd.Command1,
        cmd.Command2,
        cmd.Data1,
        cmd.Data2,
        cmd.Checksum,
    }
    
    _, err := p.connection.Write(data)
    return err
}

func (p *PTZDriverImpl) calculateChecksum(cmd *PelcoDCommand) byte {
    return byte((cmd.Address + cmd.Command1 + cmd.Command2 + cmd.Data1 + cmd.Data2) % 256)
}
```

### 4.3.5 云台自动跟踪

```go
// 云台自动跟踪服务
type PTZAutoTrack struct {
    ptzDriver PTZDriver
    targets   map[uint16]*RadarTarget
    mu        sync.RWMutex
}

func (t *PTZAutoTrack) TrackTarget(targetID uint16) error {
    t.mu.RLock()
    target, exists := t.targets[targetID]
    t.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("target %d not found", targetID)
    }
    
    // 计算云台角度
    panAngle := t.calculatePanAngle(target)
    tiltAngle := t.calculateTiltAngle(target)
    
    // 控制云台转动
    return t.moveToPosition(panAngle, tiltAngle)
}

func (t *PTZAutoTrack) calculatePanAngle(target *RadarTarget) float64 {
    // 根据目标位置计算水平角度
    return math.Atan2(target.Distance*math.Sin(target.Angle*math.Pi/180), 
                     target.Distance*math.Cos(target.Angle*math.Pi/180)) * 180 / math.Pi
}

func (t *PTZAutoTrack) moveToPosition(panAngle, tiltAngle float64) error {
    // 获取当前云台位置
    currentPos, err := t.ptzDriver.GetPosition()
    if err != nil {
        return err
    }
    
    // 计算角度差
    panDiff := panAngle - currentPos.Pan
    tiltDiff := tiltAngle - currentPos.Tilt
    
    // 控制云台转动
    if math.Abs(panDiff) > 1.0 {
        if panDiff > 0 {
            t.ptzDriver.PanRight(byte(math.Min(255, math.Abs(panDiff))))
        } else {
            t.ptzDriver.PanLeft(byte(math.Min(255, math.Abs(panDiff))))
        }
    }
    
    if math.Abs(tiltDiff) > 1.0 {
        if tiltDiff > 0 {
            t.ptzDriver.TiltUp(byte(math.Min(255, math.Abs(tiltDiff))))
        } else {
            t.ptzDriver.TiltDown(byte(math.Min(255, math.Abs(tiltDiff))))
        }
    }
    
    return nil
}
```

## 4.4 LED显示屏集成

### 4.4.1 LED显示屏概述

LED显示屏用于显示船闸状态、船舶信息和警告信息，主要功能包括：

- 船闸状态显示
- 船舶信息显示
- 警告信息显示
- 时间日期显示

### 4.4.2 LED显示屏规格

```text
设备型号：XXX LED显示屏
分辨率：64x32像素
颜色：单色/双色/全彩
通信方式：RS485/RS232
协议：自定义协议
显示内容：文字、图形、动画
```

### 4.4.3 LED显示协议

```go
// LED显示命令结构
type LEDCommand struct {
    Header    [2]byte  // 包头 0xAA 0x55
    Length    uint16   // 数据长度
    Command   uint8    // 命令字
    Address   uint8    // 设备地址
    Data      []byte   // 数据内容
    Checksum  uint8    // 校验和
}

// LED命令定义
const (
    CMD_CLEAR_SCREEN  = 0x01  // 清屏
    CMD_SET_TEXT      = 0x02  // 设置文字
    CMD_SET_GRAPHIC   = 0x03  // 设置图形
    CMD_SET_BRIGHT    = 0x04  // 设置亮度
    CMD_SET_COLOR     = 0x05  // 设置颜色
    CMD_PLAY_ANIM     = 0x06  // 播放动画
)
```

### 4.4.4 LED驱动实现

```go
// LED驱动接口
type LEDDriver interface {
    Connect() error
    Disconnect() error
    ClearScreen() error
    SetText(text string, x, y int, color byte) error
    SetGraphic(data []byte, x, y int) error
    SetBrightness(brightness byte) error
    PlayAnimation(animation []byte) error
}

// LED驱动实现
type LEDDriverImpl struct {
    connection net.Conn
    address    byte
    width      int
    height     int
}

func (l *LEDDriverImpl) SetText(text string, x, y int, color byte) error {
    // 将文字转换为点阵数据
    fontData := l.textToFontData(text)
    
    cmd := &LEDCommand{
        Header:   [2]byte{0xAA, 0x55},
        Command:  CMD_SET_TEXT,
        Address:  l.address,
        Data:     append([]byte{byte(x), byte(y), color}, fontData...),
    }
    cmd.Length = uint16(len(cmd.Data))
    cmd.Checksum = l.calculateChecksum(cmd)
    
    return l.sendCommand(cmd)
}

func (l *LEDDriverImpl) textToFontData(text string) []byte {
    // 将文字转换为点阵数据
    // 这里需要实现字体转换逻辑
    var fontData []byte
    
    for _, char := range text {
        charData := l.getCharFontData(char)
        fontData = append(fontData, charData...)
    }
    
    return fontData
}

func (l *LEDDriverImpl) sendCommand(cmd *LEDCommand) error {
    data := []byte{
        cmd.Header[0], cmd.Header[1],
        byte(cmd.Length >> 8), byte(cmd.Length),
        cmd.Command,
        cmd.Address,
    }
    data = append(data, cmd.Data...)
    data = append(data, cmd.Checksum)
    
    _, err := l.connection.Write(data)
    return err
}
```

### 4.4.5 LED显示服务

```go
// LED显示服务
type LEDDisplayService struct {
    ledDriver LEDDriver
    content   *DisplayContent
    mu        sync.RWMutex
}

type DisplayContent struct {
    Text    string
    Graphic []byte
    Color   byte
    X       int
    Y       int
}

func (s *LEDDisplayService) DisplayShipInfo(ship *Ship) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 格式化船舶信息
    text := fmt.Sprintf("Ship: %s Speed: %.1fm/s", ship.ShipID, ship.Speed)
    
    // 设置显示内容
    s.content = &DisplayContent{
        Text:  text,
        Color: 0xFF, // 白色
        X:     0,
        Y:     0,
    }
    
    // 发送到LED显示屏
    return s.ledDriver.SetText(s.content.Text, s.content.X, s.content.Y, s.content.Color)
}

func (s *LEDDisplayService) DisplayWarning(warning string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 设置警告内容
    s.content = &DisplayContent{
        Text:  warning,
        Color: 0xF0, // 红色
        X:     0,
        Y:     16,
    }
    
    // 发送到LED显示屏
    return s.ledDriver.SetText(s.content.Text, s.content.X, s.content.Y, s.content.Color)
}

func (s *LEDDisplayService) DisplayLockStatus(status *LockStatus) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 格式化船闸状态
    text := fmt.Sprintf("Lock: %s Status: %s", status.LockID, status.Status)
    
    // 设置显示内容
    s.content = &DisplayContent{
        Text:  text,
        Color: 0x0F, // 绿色
        X:     0,
        Y:     8,
    }
    
    // 发送到LED显示屏
    return s.ledDriver.SetText(s.content.Text, s.content.X, s.content.Y, s.content.Color)
}
```

## 4.5 开关量设备集成

### 4.5.1 开关量设备概述

开关量设备用于检测和控制船闸的各种状态，主要功能包括：

- 船闸开关状态检测
- 设备运行状态检测
- 安全信号检测
- 控制信号输出

### 4.5.2 开关量设备规格

```text
设备型号：XXX开关量模块
输入通道：16路
输出通道：8路
通信方式：Modbus RTU
协议：Modbus RTU
电压：24V DC
```

### 4.5.3 Modbus协议实现

```go
// Modbus RTU帧结构
type ModbusFrame struct {
    Address  byte   // 设备地址
    Function byte   // 功能码
    Data     []byte // 数据
    CRC      uint16 // CRC校验
}

// Modbus功能码
const (
    FUNC_READ_COILS     = 0x01  // 读线圈
    FUNC_READ_DISCRETE  = 0x02  // 读离散输入
    FUNC_READ_HOLDING   = 0x03  // 读保持寄存器
    FUNC_READ_INPUT     = 0x04  // 读输入寄存器
    FUNC_WRITE_COIL     = 0x05  // 写单个线圈
    FUNC_WRITE_REGISTER = 0x06  // 写单个寄存器
    FUNC_WRITE_COILS    = 0x0F  // 写多个线圈
    FUNC_WRITE_REGISTERS = 0x10 // 写多个寄存器
)
```

### 4.5.4 开关量驱动实现

```go
// 开关量驱动接口
type DigitalIODriver interface {
    Connect() error
    Disconnect() error
    ReadInput(channel int) (bool, error)
    ReadAllInputs() ([]bool, error)
    WriteOutput(channel int, value bool) error
    WriteAllOutputs(values []bool) error
    GetStatus() (*DigitalIOStatus, error)
}

// 开关量驱动实现
type DigitalIODriverImpl struct {
    connection net.Conn
    address    byte
    inputs     []bool
    outputs    []bool
    mu         sync.RWMutex
}

func (d *DigitalIODriverImpl) ReadInput(channel int) (bool, error) {
    frame := &ModbusFrame{
        Address:  d.address,
        Function: FUNC_READ_DISCRETE,
        Data:     []byte{byte(channel >> 8), byte(channel), 0x00, 0x01},
    }
    frame.CRC = d.calculateCRC(frame)
    
    response, err := d.sendFrame(frame)
    if err != nil {
        return false, err
    }
    
    // 解析响应
    if len(response.Data) < 1 {
        return false, fmt.Errorf("invalid response")
    }
    
    return response.Data[0] != 0, nil
}

func (d *DigitalIODriverImpl) WriteOutput(channel int, value bool) error {
    frame := &ModbusFrame{
        Address:  d.address,
        Function: FUNC_WRITE_COIL,
        Data:     []byte{byte(channel >> 8), byte(channel), 0x00, boolToByte(value)},
    }
    frame.CRC = d.calculateCRC(frame)
    
    _, err := d.sendFrame(frame)
    return err
}

func (d *DigitalIODriverImpl) calculateCRC(frame *ModbusFrame) uint16 {
    data := append([]byte{frame.Address, frame.Function}, frame.Data...)
    
    crc := uint16(0xFFFF)
    for _, b := range data {
        crc ^= uint16(b)
        for i := 0; i < 8; i++ {
            if crc&0x0001 != 0 {
                crc >>= 1
                crc ^= 0xA001
            } else {
                crc >>= 1
            }
        }
    }
    
    return crc
}

func boolToByte(b bool) byte {
    if b {
        return 0xFF
    }
    return 0x00
}
```

### 4.5.5 开关量监控服务

```go
// 开关量监控服务
type DigitalIOMonitor struct {
    driver DigitalIODriver
    inputs map[int]*InputChannel
    mu     sync.RWMutex
}

type InputChannel struct {
    Channel   int
    Name      string
    Value     bool
    LastValue bool
    Timestamp time.Time
}

func (m *DigitalIOMonitor) StartMonitoring() error {
    ticker := time.NewTicker(100 * time.Millisecond)
    
    go func() {
        for range ticker.C {
            m.checkInputs()
        }
    }()
    
    return nil
}

func (m *DigitalIOMonitor) checkInputs() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    for channel, input := range m.inputs {
        value, err := m.driver.ReadInput(channel)
        if err != nil {
            log.Printf("Failed to read input %d: %v", channel, err)
            continue
        }
        
        // 检查状态变化
        if value != input.Value {
            input.LastValue = input.Value
            input.Value = value
            input.Timestamp = time.Now()
            
            // 发布状态变化事件
            m.publishInputChange(input)
        }
    }
}

func (m *DigitalIOMonitor) publishInputChange(input *InputChannel) {
    event := &InputChangeEvent{
        Channel:   input.Channel,
        Name:      input.Name,
        OldValue:  input.LastValue,
        NewValue:  input.Value,
        Timestamp: input.Timestamp,
    }
    
    // 发布到消息总线
    // publishEvent("digital_io.input_change", event)
}
```

## 4.6 设备管理

### 4.6.1 设备管理器

```go
// 设备管理器
type DeviceManager struct {
    devices map[string]Device
    mu      sync.RWMutex
}

type Device interface {
    GetID() string
    GetType() string
    Connect() error
    Disconnect() error
    GetStatus() (*DeviceStatus, error)
}

func (dm *DeviceManager) AddDevice(device Device) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    dm.devices[device.GetID()] = device
    return nil
}

func (dm *DeviceManager) RemoveDevice(deviceID string) error {
    dm.mu.Lock()
    defer dm.mu.Unlock()
    
    delete(dm.devices, deviceID)
    return nil
}

func (dm *DeviceManager) GetDevice(deviceID string) (Device, error) {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    device, exists := dm.devices[deviceID]
    if !exists {
        return nil, fmt.Errorf("device %s not found", deviceID)
    }
    
    return device, nil
}

func (dm *DeviceManager) GetAllDevices() []Device {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    
    devices := make([]Device, 0, len(dm.devices))
    for _, device := range dm.devices {
        devices = append(devices, device)
    }
    
    return devices
}
```

### 4.6.2 设备状态监控

```go
// 设备状态监控
type DeviceMonitor struct {
    deviceManager *DeviceManager
    statusMap     map[string]*DeviceStatus
    mu            sync.RWMutex
}

func (dm *DeviceMonitor) StartMonitoring() error {
    ticker := time.NewTicker(5 * time.Second)
    
    go func() {
        for range ticker.C {
            dm.checkAllDevices()
        }
    }()
    
    return nil
}

func (dm *DeviceMonitor) checkAllDevices() {
    devices := dm.deviceManager.GetAllDevices()
    
    for _, device := range devices {
        status, err := device.GetStatus()
        if err != nil {
            log.Printf("Failed to get status for device %s: %v", device.GetID(), err)
            continue
        }
        
        dm.mu.Lock()
        dm.statusMap[device.GetID()] = status
        dm.mu.Unlock()
        
        // 检查设备健康状态
        if !status.IsHealthy {
            dm.handleDeviceFailure(device, status)
        }
    }
}

func (dm *DeviceMonitor) handleDeviceFailure(device Device, status *DeviceStatus) {
    // 记录设备故障
    log.Printf("Device %s is unhealthy: %s", device.GetID(), status.ErrorMessage)
    
    // 尝试重新连接
    go func() {
        device.Disconnect()
        time.Sleep(5 * time.Second)
        device.Connect()
    }()
    
    // 发布设备故障事件
    event := &DeviceFailureEvent{
        DeviceID:    device.GetID(),
        DeviceType:  device.GetType(),
        ErrorMessage: status.ErrorMessage,
        Timestamp:   time.Now(),
    }
    
    // publishEvent("device.failure", event)
}
```

## 4.7 设备配置管理

### 4.7.1 设备配置结构

```go
// 设备配置
type DeviceConfig struct {
    ID       string                 `json:"id"`
    Type     string                 `json:"type"`
    Name     string                 `json:"name"`
    Address  string                 `json:"address"`
    Port     int                    `json:"port"`
    Protocol string                 `json:"protocol"`
    Params   map[string]interface{} `json:"params"`
}

// 雷达配置
type RadarConfig struct {
    Address     string  `json:"address"`
    Port        int     `json:"port"`
    ScanRange   float64 `json:"scan_range"`
    ScanAngle   float64 `json:"scan_angle"`
    UpdateRate  int     `json:"update_rate"`
}

// 云台配置
type PTZConfig struct {
    Address     string  `json:"address"`
    Port        int     `json:"port"`
    Protocol    string  `json:"protocol"`
    PanSpeed    byte    `json:"pan_speed"`
    TiltSpeed   byte    `json:"tilt_speed"`
    ZoomSpeed   byte    `json:"zoom_speed"`
}

// LED配置
type LEDConfig struct {
    Address     string `json:"address"`
    Port        int    `json:"port"`
    Width       int    `json:"width"`
    Height      int    `json:"height"`
    Brightness  byte   `json:"brightness"`
}
```

### 4.7.2 配置加载器

```go
// 设备配置加载器
type DeviceConfigLoader struct {
    configPath string
}

func (l *DeviceConfigLoader) LoadConfig() (map[string]*DeviceConfig, error) {
    data, err := os.ReadFile(l.configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var configs map[string]*DeviceConfig
    err = json.Unmarshal(data, &configs)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return configs, nil
}

func (l *DeviceConfigLoader) ValidateConfig(config *DeviceConfig) error {
    if config.ID == "" {
        return fmt.Errorf("device ID is required")
    }
    
    if config.Type == "" {
        return fmt.Errorf("device type is required")
    }
    
    if config.Address == "" {
        return fmt.Errorf("device address is required")
    }
    
    return nil
}
```

## 4.8 设备测试

### 4.8.1 设备测试框架

```go
// 设备测试接口
type DeviceTester interface {
    TestConnection() error
    TestFunctionality() error
    TestPerformance() error
    GenerateReport() *TestReport
}

// 设备测试实现
type DeviceTesterImpl struct {
    device Device
    report *TestReport
}

func (t *DeviceTesterImpl) TestConnection() error {
    start := time.Now()
    
    err := t.device.Connect()
    if err != nil {
        t.report.ConnectionTest = &TestResult{
            Success:   false,
            Error:     err.Error(),
            Duration:  time.Since(start),
        }
        return err
    }
    
    defer t.device.Disconnect()
    
    t.report.ConnectionTest = &TestResult{
        Success:  true,
        Duration: time.Since(start),
    }
    
    return nil
}

func (t *DeviceTesterImpl) TestFunctionality() error {
    // 测试设备基本功能
    status, err := t.device.GetStatus()
    if err != nil {
        t.report.FunctionalityTest = &TestResult{
            Success: false,
            Error:   err.Error(),
        }
        return err
    }
    
    t.report.FunctionalityTest = &TestResult{
        Success: true,
        Data:    status,
    }
    
    return nil
}

func (t *DeviceTesterImpl) GenerateReport() *TestReport {
    t.report.DeviceID = t.device.GetID()
    t.report.DeviceType = t.device.GetType()
    t.report.Timestamp = time.Now()
    
    return t.report
}
```

### 4.8.2 测试报告

```go
// 测试报告
type TestReport struct {
    DeviceID           string      `json:"device_id"`
    DeviceType         string      `json:"device_type"`
    Timestamp          time.Time   `json:"timestamp"`
    ConnectionTest     *TestResult `json:"connection_test"`
    FunctionalityTest  *TestResult `json:"functionality_test"`
    PerformanceTest    *TestResult `json:"performance_test"`
    OverallResult      string      `json:"overall_result"`
}

type TestResult struct {
    Success   bool        `json:"success"`
    Error     string      `json:"error,omitempty"`
    Duration  time.Duration `json:"duration,omitempty"`
    Data      interface{} `json:"data,omitempty"`
}
```
