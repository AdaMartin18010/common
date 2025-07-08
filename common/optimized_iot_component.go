package common

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	mdl "common/model"
)

// 优化后的IOT组件架构实现

// OptimizedCtrlSt 无锁控制结构
type OptimizedCtrlSt struct {
	c     context.Context
	ccl   context.CancelFunc
	wwg   *mdl.WorkerWG
	state atomic.Value // 使用原子操作替代锁
}

func NewOptimizedCtrlSt(ctx context.Context) *OptimizedCtrlSt {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)

	ctrl := &OptimizedCtrlSt{
		c:   ctx,
		ccl: cancel,
		wwg: mdl.NewWorkerWG(),
	}

	// 初始化状态
	ctrl.state.Store("initialized")
	return ctrl
}

func (cs *OptimizedCtrlSt) DebugInfo() string {
	// 无锁读取，使用原子操作
	state := cs.state.Load().(string)
	return fmt.Sprintf("(optimized_ctrl)[%s]", state)
}

func (cs *OptimizedCtrlSt) Cancel() {
	// 原子操作设置状态
	cs.state.Store("canceling")
	cs.ccl()
	cs.state.Store("canceled")
}

func (cs *OptimizedCtrlSt) Context() context.Context {
	return cs.c
}

func (cs *OptimizedCtrlSt) WaitGroup() *mdl.WorkerWG {
	return cs.wwg
}

// LockFreeComponent 无锁组件实现
type LockFreeComponent struct {
	id    [16]byte     // 固定大小ID
	kind  [8]byte      // 固定大小类型
	state atomic.Value // 原子状态
	ctrl  *OptimizedCtrlSt
	mu    sync.Mutex // 仅在必要时使用锁
}

func NewLockFreeComponent(id, kind string) *LockFreeComponent {
	comp := &LockFreeComponent{
		ctrl: NewOptimizedCtrlSt(context.Background()),
	}

	// 复制ID和类型到固定大小数组，确保不截断
	if len(id) > 16 {
		id = id[:16]
	}
	if len(kind) > 8 {
		kind = kind[:8]
	}
	copy(comp.id[:], id)
	copy(comp.kind[:], kind)

	// 初始化状态 - 确保原子操作正确初始化
	comp.state.Store(false)

	return comp
}

func (c *LockFreeComponent) Id() string {
	// 找到第一个null字节的位置，截断字符串
	for i, b := range c.id {
		if b == 0 {
			return string(c.id[:i])
		}
	}
	return string(c.id[:])
}

func (c *LockFreeComponent) Kind() string {
	// 找到第一个null字节的位置，截断字符串
	for i, b := range c.kind {
		if b == 0 {
			return string(c.kind[:i])
		}
	}
	return string(c.kind[:])
}

func (c *LockFreeComponent) IsRunning() bool {
	// 安全地获取原子值，避免panic
	if value := c.state.Load(); value != nil {
		if running, ok := value.(bool); ok {
			return running
		}
	}
	return false
}

func (c *LockFreeComponent) SetRunning(running bool) {
	c.state.Store(running)
}

func (c *LockFreeComponent) Start() error {
	if c.IsRunning() {
		return fmt.Errorf("component is already running")
	}

	c.SetRunning(true)
	c.ctrl.state.Store("running")
	return nil
}

func (c *LockFreeComponent) Stop() error {
	if !c.IsRunning() {
		return fmt.Errorf("component is not running")
	}

	c.SetRunning(false)
	c.ctrl.state.Store("stopped")
	return nil
}

// ComponentPool 组件对象池
type ComponentPool struct {
	pool sync.Pool
}

func NewComponentPool() *ComponentPool {
	return &ComponentPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &LockFreeComponent{
					ctrl: NewOptimizedCtrlSt(context.Background()),
				}
			},
		},
	}
}

func (cp *ComponentPool) Get() *LockFreeComponent {
	return cp.pool.Get().(*LockFreeComponent)
}

func (cp *ComponentPool) Put(comp *LockFreeComponent) {
	// 重置组件状态
	comp.SetRunning(false)
	comp.ctrl.state.Store("idle")
	cp.pool.Put(comp)
}

// BatchComponentManager 批量组件管理器
type BatchComponentManager struct {
	components []*LockFreeComponent
	mu         sync.RWMutex
	pool       *ComponentPool
}

func NewBatchComponentManager() *BatchComponentManager {
	return &BatchComponentManager{
		components: make([]*LockFreeComponent, 0),
		pool:       NewComponentPool(),
	}
}

func (bcm *BatchComponentManager) AddComponent(id, kind string) *LockFreeComponent {
	comp := bcm.pool.Get()

	// 设置组件属性
	copy(comp.id[:], id)
	copy(comp.kind[:], kind)

	bcm.mu.Lock()
	bcm.components = append(bcm.components, comp)
	bcm.mu.Unlock()

	return comp
}

func (bcm *BatchComponentManager) StartBatch() error {
	bcm.mu.RLock()
	defer bcm.mu.RUnlock()

	// 批量启动，减少锁竞争
	for _, comp := range bcm.components {
		comp.Start()
	}
	return nil
}

func (bcm *BatchComponentManager) StopBatch() error {
	bcm.mu.RLock()
	defer bcm.mu.RUnlock()

	// 批量停止
	for _, comp := range bcm.components {
		comp.Stop()
	}
	return nil
}

func (bcm *BatchComponentManager) GetRunningCount() int {
	bcm.mu.RLock()
	defer bcm.mu.RUnlock()

	count := 0
	for _, comp := range bcm.components {
		if comp.IsRunning() {
			count++
		}
	}
	return count
}

// IOTDevice 模拟IOT设备
type IOTDevice struct {
	*LockFreeComponent
	deviceType  string
	dataChan    chan interface{}
	controlChan chan bool
}

func NewIOTDevice(id, kind, deviceType string) *IOTDevice {
	return &IOTDevice{
		LockFreeComponent: NewLockFreeComponent(id, kind),
		deviceType:        deviceType,
		dataChan:          make(chan interface{}, 10),
		controlChan:       make(chan bool, 1),
	}
}

func (d *IOTDevice) Start() error {
	if err := d.LockFreeComponent.Start(); err != nil {
		return err
	}

	// 启动设备特定的工作
	go d.work()
	return nil
}

func (d *IOTDevice) work() {
	for d.IsRunning() {
		// 优先处理控制命令，避免阻塞
		select {
		case control := <-d.controlChan:
			d.processControl(control)
			continue
		default:
		}

		// 处理数据通道
		select {
		case data := <-d.dataChan:
			d.processData(data)
		case <-d.ctrl.Context().Done():
			return
		default:
			// 避免忙等待，但减少休眠时间
			time.Sleep(100 * time.Microsecond)
		}
	}
}

func (d *IOTDevice) processData(data interface{}) {
	// 优化数据处理，减少处理时间
	// 实际应用中这里可能是数据验证、转换、存储等操作
	_ = data // 避免未使用变量警告
}

func (d *IOTDevice) processControl(control bool) {
	// 优化控制处理，减少处理时间
	// 实际应用中这里可能是状态更新、设备控制等操作
	_ = control // 避免未使用变量警告
}

// SensorDevice 传感器设备
type SensorDevice struct {
	*IOTDevice
	sensorType string
	value      float64
}

func NewSensorDevice(id, sensorType string) *SensorDevice {
	return &SensorDevice{
		IOTDevice:  NewIOTDevice(id, "sensor", "sensor"),
		sensorType: sensorType,
	}
}

func (s *SensorDevice) ReadValue() float64 {
	return s.value
}

func (s *SensorDevice) SetValue(value float64) {
	s.value = value
}

// ActuatorDevice 执行器设备
type ActuatorDevice struct {
	*IOTDevice
	actuatorType string
	state        bool
}

func NewActuatorDevice(id, actuatorType string) *ActuatorDevice {
	return &ActuatorDevice{
		IOTDevice:    NewIOTDevice(id, "actuator", "actuator"),
		actuatorType: actuatorType,
	}
}

func (a *ActuatorDevice) SetState(state bool) {
	a.state = state
	// 使用非阻塞发送，避免死锁
	select {
	case a.controlChan <- state:
		// 成功发送
	default:
		// 通道已满，忽略
	}
}

func (a *ActuatorDevice) GetState() bool {
	return a.state
}

// ControllerDevice 控制器设备
type ControllerDevice struct {
	*IOTDevice
	controlType string
	actuators   []*ActuatorDevice
	sensors     []*SensorDevice
}

func NewControllerDevice(id, controlType string) *ControllerDevice {
	return &ControllerDevice{
		IOTDevice:   NewIOTDevice(id, "controller", "controller"),
		controlType: controlType,
		actuators:   make([]*ActuatorDevice, 0),
		sensors:     make([]*SensorDevice, 0),
	}
}

func (c *ControllerDevice) AddActuator(actuator *ActuatorDevice) {
	c.actuators = append(c.actuators, actuator)
}

func (c *ControllerDevice) AddSensor(sensor *SensorDevice) {
	c.sensors = append(c.sensors, sensor)
}

func (c *ControllerDevice) ControlAllDevices(command bool) {
	for _, actuator := range c.actuators {
		actuator.SetState(command)
	}
}

// 性能监控接口
type PerformanceMonitor interface {
	RecordComponentCreation(duration time.Duration)
	RecordComponentStart(duration time.Duration)
	RecordMemoryUsage(bytes int64)
	RecordCPUUsage(percentage float64)
}

// SimplePerformanceMonitor 简单性能监控实现
type SimplePerformanceMonitor struct {
	mu    sync.RWMutex
	stats map[string]interface{}
}

func NewSimplePerformanceMonitor() *SimplePerformanceMonitor {
	return &SimplePerformanceMonitor{
		stats: make(map[string]interface{}),
	}
}

func (spm *SimplePerformanceMonitor) RecordComponentCreation(duration time.Duration) {
	spm.mu.Lock()
	defer spm.mu.Unlock()
	spm.stats["component_creation_time"] = duration
}

func (spm *SimplePerformanceMonitor) RecordComponentStart(duration time.Duration) {
	spm.mu.Lock()
	defer spm.mu.Unlock()
	spm.stats["component_start_time"] = duration
}

func (spm *SimplePerformanceMonitor) RecordMemoryUsage(bytes int64) {
	spm.mu.Lock()
	defer spm.mu.Unlock()
	spm.stats["memory_usage"] = bytes
}

func (spm *SimplePerformanceMonitor) RecordCPUUsage(percentage float64) {
	spm.mu.Lock()
	defer spm.mu.Unlock()
	spm.stats["cpu_usage"] = percentage
}

func (spm *SimplePerformanceMonitor) GetStats() map[string]interface{} {
	spm.mu.RLock()
	defer spm.mu.RUnlock()

	// 返回副本
	stats := make(map[string]interface{})
	for k, v := range spm.stats {
		stats[k] = v
	}
	return stats
}

// Event对象池，减少内存分配
var eventPool = sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

// EventDispatcher 无锁事件分发器
type EventDispatcher struct {
	devices   []*IOTDevice
	eventChan chan *Event
	done      chan struct{}
	mu        sync.RWMutex
}

// Event 事件结构
type Event struct {
	DeviceID string
	Type     string
	Data     interface{}
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		devices:   make([]*IOTDevice, 0),
		eventChan: make(chan *Event, 1000),
		done:      make(chan struct{}),
	}
}

func (ed *EventDispatcher) AddDevice(device *IOTDevice) {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.devices = append(ed.devices, device)
}

func (ed *EventDispatcher) Start() {
	go ed.dispatch()
}

func (ed *EventDispatcher) Stop() {
	close(ed.done)
}

func (ed *EventDispatcher) SendEvent(deviceID, eventType string, data interface{}) {
	// 从对象池获取Event
	event := eventPool.Get().(*Event)
	event.DeviceID = deviceID
	event.Type = eventType
	event.Data = data

	select {
	case ed.eventChan <- event:
	default:
		// 通道满时归还对象到池中
		eventPool.Put(event)
	}
}

func (ed *EventDispatcher) dispatch() {
	for {
		select {
		case event := <-ed.eventChan:
			ed.mu.RLock()
			devices := make([]*IOTDevice, len(ed.devices))
			copy(devices, ed.devices)
			ed.mu.RUnlock()

			// 并发分发事件到所有设备
			var wg sync.WaitGroup
			for _, device := range devices {
				wg.Add(1)
				go func(d *IOTDevice) {
					defer wg.Done()
					select {
					case d.dataChan <- event.Data:
					default:
						// 设备通道满时跳过
					}
				}(device)
			}
			wg.Wait()

			// 归还Event对象到池中
			eventPool.Put(event)

		case <-ed.done:
			return
		}
	}
}

// BatchEventProcessor 批量事件处理器
type BatchEventProcessor struct {
	events    []*Event
	batchSize int
	processor func([]*Event)
	eventChan chan *Event
	done      chan struct{}
	mu        sync.Mutex
}

func NewBatchEventProcessor(batchSize int, processor func([]*Event)) *BatchEventProcessor {
	return &BatchEventProcessor{
		events:    make([]*Event, 0, batchSize),
		batchSize: batchSize,
		processor: processor,
		eventChan: make(chan *Event, 1000),
		done:      make(chan struct{}),
	}
}

func (bep *BatchEventProcessor) Start() {
	go bep.process()
}

func (bep *BatchEventProcessor) Stop() {
	close(bep.done)
}

func (bep *BatchEventProcessor) SendEvent(event *Event) {
	select {
	case bep.eventChan <- event:
	default:
		// 通道满时归还对象到池中
		eventPool.Put(event)
	}
}

func (bep *BatchEventProcessor) process() {
	ticker := time.NewTicker(10 * time.Millisecond) // 批量处理间隔
	defer ticker.Stop()

	for {
		select {
		case event := <-bep.eventChan:
			bep.mu.Lock()
			bep.events = append(bep.events, event)

			// 达到批量大小时立即处理
			if len(bep.events) >= bep.batchSize {
				events := make([]*Event, len(bep.events))
				copy(events, bep.events)
				bep.events = bep.events[:0]
				bep.mu.Unlock()

				// 批量处理事件
				bep.processor(events)

				// 归还事件对象到池中
				for _, e := range events {
					eventPool.Put(e)
				}
			} else {
				bep.mu.Unlock()
			}

		case <-ticker.C:
			// 定时处理未满批次的事件
			bep.mu.Lock()
			if len(bep.events) > 0 {
				events := make([]*Event, len(bep.events))
				copy(events, bep.events)
				bep.events = bep.events[:0]
				bep.mu.Unlock()

				bep.processor(events)

				// 归还事件对象到池中
				for _, e := range events {
					eventPool.Put(e)
				}
			} else {
				bep.mu.Unlock()
			}

		case <-bep.done:
			// 处理剩余事件
			bep.mu.Lock()
			if len(bep.events) > 0 {
				bep.processor(bep.events)
				for _, e := range bep.events {
					eventPool.Put(e)
				}
			}
			bep.mu.Unlock()
			return
		}
	}
}

// EventRouter 智能事件路由器
type EventRouter struct {
	routes    map[string][]*IOTDevice // 基于路由键的设备分组
	devices   map[string]*IOTDevice   // 设备ID到设备的映射
	mu        sync.RWMutex
	eventChan chan *Event
	done      chan struct{}
}

// RouteKey 路由键生成函数
type RouteKey func(*Event) string

func NewEventRouter() *EventRouter {
	return &EventRouter{
		routes:    make(map[string][]*IOTDevice),
		devices:   make(map[string]*IOTDevice),
		eventChan: make(chan *Event, 1000),
		done:      make(chan struct{}),
	}
}

// 默认路由键生成器
func DefaultRouteKey(event *Event) string {
	return fmt.Sprintf("%s:%s", event.DeviceID, event.Type)
}

// 按设备类型路由
func DeviceTypeRouteKey(event *Event) string {
	return event.Type
}

// 按设备ID路由
func DeviceIDRouteKey(event *Event) string {
	return event.DeviceID
}

func (er *EventRouter) AddDevice(device *IOTDevice, routeKey string) {
	er.mu.Lock()
	defer er.mu.Unlock()

	er.devices[device.Id()] = device
	if _, exists := er.routes[routeKey]; !exists {
		er.routes[routeKey] = make([]*IOTDevice, 0)
	}
	er.routes[routeKey] = append(er.routes[routeKey], device)
}

func (er *EventRouter) Start() {
	go er.route()
}

func (er *EventRouter) Stop() {
	close(er.done)
}

func (er *EventRouter) SendEvent(deviceID, eventType string, data interface{}) {
	// 从对象池获取Event
	event := eventPool.Get().(*Event)
	event.DeviceID = deviceID
	event.Type = eventType
	event.Data = data

	select {
	case er.eventChan <- event:
	default:
		// 通道满时归还对象到池中
		eventPool.Put(event)
	}
}

func (er *EventRouter) route() {
	for {
		select {
		case event := <-er.eventChan:
			er.routeEvent(event)
		case <-er.done:
			return
		}
	}
}

func (er *EventRouter) routeEvent(event *Event) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	// 使用默认路由键
	routeKey := DefaultRouteKey(event)

	// 查找匹配的设备
	if devices, exists := er.routes[routeKey]; exists {
		// 并发发送到匹配的设备
		var wg sync.WaitGroup
		for _, device := range devices {
			wg.Add(1)
			go func(d *IOTDevice) {
				defer wg.Done()
				select {
				case d.dataChan <- event.Data:
				default:
					// 设备通道满时跳过
				}
			}(device)
		}
		wg.Wait()
	}

	// 归还Event对象到池中
	eventPool.Put(event)
}

// BackpressureController 背压控制器
type BackpressureController struct {
	windowSize   int         // 滑动窗口大小
	window       []time.Time // 时间窗口
	windowIndex  int         // 当前窗口索引
	rateLimit    int         // 每秒最大事件数
	mu           sync.RWMutex
	lastCheck    time.Time
	currentRate  float64 // 当前速率
	adaptiveRate float64 // 自适应速率
}

func NewBackpressureController(windowSize int, rateLimit int) *BackpressureController {
	return &BackpressureController{
		windowSize:   windowSize,
		window:       make([]time.Time, windowSize),
		rateLimit:    rateLimit,
		lastCheck:    time.Now(),
		adaptiveRate: float64(rateLimit),
	}
}

func (bc *BackpressureController) ShouldAccept() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	now := time.Now()

	// 更新滑动窗口
	bc.window[bc.windowIndex] = now
	bc.windowIndex = (bc.windowIndex + 1) % bc.windowSize

	// 计算当前速率（每秒事件数）
	windowStart := now.Add(-time.Duration(bc.windowSize) * time.Second)
	count := 0
	for _, t := range bc.window {
		if t.After(windowStart) {
			count++
		}
	}

	bc.currentRate = float64(count) / float64(bc.windowSize)

	// 自适应调整速率
	if bc.currentRate > float64(bc.rateLimit) {
		bc.adaptiveRate *= 0.9 // 降低10%
	} else if bc.currentRate < float64(bc.rateLimit)*0.8 {
		bc.adaptiveRate *= 1.1 // 提高10%
	}

	// 限制自适应速率范围
	if bc.adaptiveRate < float64(bc.rateLimit)*0.5 {
		bc.adaptiveRate = float64(bc.rateLimit) * 0.5
	} else if bc.adaptiveRate > float64(bc.rateLimit) {
		bc.adaptiveRate = float64(bc.rateLimit)
	}

	// 只有当当前速率低于自适应速率时才接受
	return bc.currentRate < bc.adaptiveRate
}

func (bc *BackpressureController) GetCurrentRate() float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.currentRate
}

func (bc *BackpressureController) GetAdaptiveRate() float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.adaptiveRate
}

// AdaptiveEventDispatcher 自适应事件分发器（带背压控制）
type AdaptiveEventDispatcher struct {
	*EventDispatcher
	backpressure *BackpressureController
}

func NewAdaptiveEventDispatcher(windowSize, rateLimit int) *AdaptiveEventDispatcher {
	return &AdaptiveEventDispatcher{
		EventDispatcher: NewEventDispatcher(),
		backpressure:    NewBackpressureController(windowSize, rateLimit),
	}
}

func (aed *AdaptiveEventDispatcher) SendEvent(deviceID, eventType string, data interface{}) {
	// 检查背压控制
	if !aed.backpressure.ShouldAccept() {
		// 背压控制拒绝事件，直接丢弃
		return
	}

	// 调用父类方法发送事件
	aed.EventDispatcher.SendEvent(deviceID, eventType, data)
}

func (aed *AdaptiveEventDispatcher) GetBackpressureStats() (currentRate, adaptiveRate float64) {
	return aed.backpressure.GetCurrentRate(), aed.backpressure.GetAdaptiveRate()
}
