package common

import (
	"fmt"
	"testing"
	"time"
)

// 测试优化后的IOT组件架构

func TestOptimizedCtrlSt(t *testing.T) {
	ctrl := NewOptimizedCtrlSt(nil)

	// 测试初始状态
	if ctrl.DebugInfo() != "(optimized_ctrl)[initialized]" {
		t.Errorf("Expected initialized state, got %s", ctrl.DebugInfo())
	}

	// 测试取消操作
	ctrl.Cancel()
	if ctrl.DebugInfo() != "(optimized_ctrl)[canceled]" {
		t.Errorf("Expected canceled state, got %s", ctrl.DebugInfo())
	}
}

func TestLockFreeComponent(t *testing.T) {
	comp := NewLockFreeComponent("test_id", "test_kind")

	// 测试初始状态
	if comp.Id() != "test_id" {
		t.Errorf("Expected test_id, got %s", comp.Id())
	}

	if comp.Kind() != "test_kin" {
		t.Errorf("Expected test_kin, got %s", comp.Kind())
	}

	if comp.IsRunning() {
		t.Error("Component should not be running initially")
	}

	// 测试启动
	if err := comp.Start(); err != nil {
		t.Errorf("Failed to start component: %v", err)
	}

	if !comp.IsRunning() {
		t.Error("Component should be running after start")
	}

	// 测试停止
	if err := comp.Stop(); err != nil {
		t.Errorf("Failed to stop component: %v", err)
	}

	if comp.IsRunning() {
		t.Error("Component should not be running after stop")
	}
}

func TestComponentPool(t *testing.T) {
	pool := NewComponentPool()

	// 测试获取组件
	comp1 := pool.Get()
	comp2 := pool.Get()

	if comp1 == comp2 {
		t.Error("Pool should return different components")
	}

	// 测试放回组件
	pool.Put(comp1)
	pool.Put(comp2)

	// 再次获取应该能复用
	comp3 := pool.Get()
	if comp3 == nil {
		t.Error("Pool should return a component")
	}
}

func TestBatchComponentManager(t *testing.T) {
	bcm := NewBatchComponentManager()

	// 添加组件
	comp1 := bcm.AddComponent("device1", "sensor")
	comp2 := bcm.AddComponent("device2", "actuator")
	comp3 := bcm.AddComponent("device3", "controller")

	// 测试批量启动
	if err := bcm.StartBatch(); err != nil {
		t.Errorf("Failed to start batch: %v", err)
	}

	// 检查运行状态
	if !comp1.IsRunning() || !comp2.IsRunning() || !comp3.IsRunning() {
		t.Error("All components should be running after batch start")
	}

	// 测试批量停止
	if err := bcm.StopBatch(); err != nil {
		t.Errorf("Failed to stop batch: %v", err)
	}

	// 检查停止状态
	if comp1.IsRunning() || comp2.IsRunning() || comp3.IsRunning() {
		t.Error("All components should be stopped after batch stop")
	}

	// 测试运行计数
	if count := bcm.GetRunningCount(); count != 0 {
		t.Errorf("Expected 0 running components, got %d", count)
	}
}

func TestIOTDevice(t *testing.T) {
	device := NewIOTDevice("iot_device", "sensor", "temperature")

	// 测试启动
	if err := device.Start(); err != nil {
		t.Errorf("Failed to start IOT device: %v", err)
	}

	if !device.IsRunning() {
		t.Error("IOT device should be running after start")
	}

	// 测试数据处理
	data := "test_data"
	device.dataChan <- data

	// 等待处理
	time.Sleep(10 * time.Millisecond)

	// 测试停止
	if err := device.Stop(); err != nil {
		t.Errorf("Failed to stop IOT device: %v", err)
	}

	if device.IsRunning() {
		t.Error("IOT device should not be running after stop")
	}
}

func TestSensorDevice(t *testing.T) {
	sensor := NewSensorDevice("temp_sensor", "temperature")

	// 测试初始值
	if sensor.ReadValue() != 0.0 {
		t.Errorf("Expected initial value 0.0, got %f", sensor.ReadValue())
	}

	// 测试设置值
	expectedValue := 25.5
	sensor.SetValue(expectedValue)

	if sensor.ReadValue() != expectedValue {
		t.Errorf("Expected value %f, got %f", expectedValue, sensor.ReadValue())
	}
}

func TestActuatorDevice(t *testing.T) {
	actuator := NewActuatorDevice("relay_actuator", "relay")

	// 测试初始状态
	if actuator.GetState() {
		t.Error("Actuator should be off initially")
	}

	// 测试设置状态
	actuator.SetState(true)

	if !actuator.GetState() {
		t.Error("Actuator should be on after setting state")
	}

	// 测试关闭状态
	actuator.SetState(false)

	if actuator.GetState() {
		t.Error("Actuator should be off after setting state to false")
	}
}

func TestControllerDevice(t *testing.T) {
	controller := NewControllerDevice("main_controller", "automation")

	// 创建传感器和执行器
	sensor := NewSensorDevice("temp_sensor", "temperature")
	actuator := NewActuatorDevice("heater_actuator", "heater")

	// 添加设备到控制器
	controller.AddSensor(sensor)
	controller.AddActuator(actuator)

	// 测试控制所有设备
	controller.ControlAllDevices(true)

	// 检查执行器状态
	if !actuator.GetState() {
		t.Error("Actuator should be on after controller command")
	}

	// 测试关闭命令
	controller.ControlAllDevices(false)

	if actuator.GetState() {
		t.Error("Actuator should be off after controller command")
	}
}

func TestPerformanceMonitor(t *testing.T) {
	monitor := NewSimplePerformanceMonitor()

	// 测试记录性能数据
	monitor.RecordComponentCreation(1 * time.Millisecond)
	monitor.RecordComponentStart(500 * time.Microsecond)
	monitor.RecordMemoryUsage(1024)
	monitor.RecordCPUUsage(25.5)

	// 获取统计数据
	stats := monitor.GetStats()

	// 验证记录的数据
	if stats["component_creation_time"] != 1*time.Millisecond {
		t.Error("Component creation time not recorded correctly")
	}

	if stats["component_start_time"] != 500*time.Microsecond {
		t.Error("Component start time not recorded correctly")
	}

	if stats["memory_usage"] != int64(1024) {
		t.Error("Memory usage not recorded correctly")
	}

	if stats["cpu_usage"] != 25.5 {
		t.Error("CPU usage not recorded correctly")
	}
}

// 基准测试：比较原始组件和优化组件的性能

func BenchmarkOriginalComponentCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 这里应该测试原始组件的创建
		// 由于原始组件在model包中，这里只是占位
		_ = i
	}
}

func BenchmarkOptimizedComponentCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comp := NewLockFreeComponent("test_id", "test_kind")
		_ = comp
	}
}

func BenchmarkComponentPool(b *testing.B) {
	pool := NewComponentPool()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp := pool.Get()
		pool.Put(comp)
	}
}

func BenchmarkBatchOperations(b *testing.B) {
	bcm := NewBatchComponentManager()

	// 预添加组件
	for i := 0; i < 100; i++ {
		bcm.AddComponent("device"+string(rune(i)), "sensor")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bcm.StartBatch()
		bcm.StopBatch()
	}
}

func BenchmarkConcurrentComponentAccess(b *testing.B) {
	comp := NewLockFreeComponent("test_id", "test_kind")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			comp.IsRunning()
			comp.Id()
			comp.Kind()
		}
	})
}

func BenchmarkLargeBatchOperations(b *testing.B) {
	bcm := NewBatchComponentManager()
	for i := 0; i < 100000; i++ {
		bcm.AddComponent(fmt.Sprintf("device%d", i), "sensor")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bcm.StartBatch()
		bcm.StopBatch()
	}
}

// 高并发事件通道基准测试
func BenchmarkHighConcurrencyDataFlow(b *testing.B) {
	// 创建100个IOT设备，模拟大规模IOT网络
	devices := make([]*IOTDevice, 100)
	for i := 0; i < 100; i++ {
		devices[i] = NewIOTDevice(fmt.Sprintf("device%d", i), "sensor", "sensor")
		devices[i].Start()
		defer devices[i].Stop()
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 模拟高并发数据发送
			for _, device := range devices {
				select {
				case device.dataChan <- "test_data":
				default:
					// 通道满时跳过
				}
			}
		}
	})
}

// 事件通道压力测试
func BenchmarkEventChannelStress(b *testing.B) {
	device := NewIOTDevice("stress_device", "sensor", "sensor")
	device.Start()
	defer device.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 模拟高频率事件发送
			select {
			case device.controlChan <- true:
			case device.dataChan <- "stress_data":
			default:
				// 通道满时跳过
			}
		}
	})
}

// 事件分发器基准测试
func BenchmarkEventDispatcher(b *testing.B) {
	dispatcher := NewEventDispatcher()

	// 添加100个设备
	for i := 0; i < 100; i++ {
		device := NewIOTDevice(fmt.Sprintf("device%d", i), "sensor", "sensor")
		device.Start()
		defer device.Stop()
		dispatcher.AddDevice(device)
	}

	dispatcher.Start()
	defer dispatcher.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dispatcher.SendEvent("test_device", "data", "test_event_data")
		}
	})
}

// 批量事件处理器基准测试
func BenchmarkBatchEventProcessor(b *testing.B) {
	// 创建批量事件处理器
	processor := NewBatchEventProcessor(100, func(events []*Event) {
		// 模拟批量处理逻辑
		for _, event := range events {
			_ = event // 避免未使用变量警告
		}
	})

	processor.Start()
	defer processor.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 从对象池获取事件
			event := eventPool.Get().(*Event)
			event.DeviceID = "test_device"
			event.Type = "data"
			event.Data = "test_event_data"

			processor.SendEvent(event)
		}
	})
}

// 智能事件路由器基准测试
func BenchmarkEventRouter(b *testing.B) {
	router := NewEventRouter()

	// 添加不同类型的设备
	for i := 0; i < 50; i++ {
		sensor := NewIOTDevice(fmt.Sprintf("sensor%d", i), "sensor", "sensor")
		sensor.Start()
		defer sensor.Stop()
		router.AddDevice(sensor, fmt.Sprintf("sensor%d:data", i))
	}

	for i := 0; i < 50; i++ {
		actuator := NewIOTDevice(fmt.Sprintf("actuator%d", i), "actuator", "actuator")
		actuator.Start()
		defer actuator.Stop()
		router.AddDevice(actuator, fmt.Sprintf("actuator%d:control", i))
	}

	router.Start()
	defer router.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 发送到传感器
			router.SendEvent("sensor0", "data", "sensor_data")
			// 发送到执行器
			router.SendEvent("actuator0", "control", "control_command")
		}
	})
}

// 自适应事件分发器基准测试
func BenchmarkAdaptiveEventDispatcher(b *testing.B) {
	dispatcher := NewAdaptiveEventDispatcher(10, 1000) // 10秒窗口，1000事件/秒限制

	// 添加100个设备
	for i := 0; i < 100; i++ {
		device := NewIOTDevice(fmt.Sprintf("device%d", i), "sensor", "sensor")
		device.Start()
		defer device.Stop()
		dispatcher.AddDevice(device)
	}

	dispatcher.Start()
	defer dispatcher.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dispatcher.SendEvent("test_device", "data", "test_event_data")
		}
	})
}

// 内存使用测试
func TestMemoryUsage(t *testing.T) {
	// 创建大量组件测试内存使用
	components := make([]*LockFreeComponent, 1000)

	for i := 0; i < 1000; i++ {
		components[i] = NewLockFreeComponent(fmt.Sprintf("device%d", i), "sensor")
	}

	// 验证所有组件都创建成功
	for i, comp := range components {
		if comp == nil {
			t.Errorf("Component %d is nil", i)
		}
		expectedId := fmt.Sprintf("device%d", i)
		if comp.Id() != expectedId {
			t.Errorf("Component %d has wrong ID, expected %s, got %s", i, expectedId, comp.Id())
		}
	}
}

// 并发安全性测试
func TestConcurrentSafety(t *testing.T) {
	comp := NewLockFreeComponent("test_id", "test_kind")

	// 启动多个goroutine同时访问组件
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				comp.IsRunning()
				comp.Id()
				comp.Kind()
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证组件状态仍然正确
	if comp.Id() != "test_id" {
		t.Error("Component ID changed during concurrent access")
	}

	if comp.Kind() != "test_kin" {
		t.Error("Component Kind changed during concurrent access")
	}
}

// 背压控制测试
func TestBackpressureController(t *testing.T) {
	bc := NewBackpressureController(5, 10) // 5秒窗口，10事件/秒限制

	// 测试正常速率
	for i := 0; i < 5; i++ {
		if !bc.ShouldAccept() {
			t.Error("Should accept events at normal rate")
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 测试高速率（应该触发背压控制）
	for i := 0; i < 50; i++ {
		bc.ShouldAccept()
	}

	// 验证自适应速率已降低
	adaptiveRate := bc.GetAdaptiveRate()
	currentRate := bc.GetCurrentRate()

	// 如果当前速率超过限制，自适应速率应该降低
	if currentRate > 10.0 && adaptiveRate >= 10.0 {
		t.Errorf("Expected adaptive rate to be reduced when current rate is %f, got %f", currentRate, adaptiveRate)
	}
}
