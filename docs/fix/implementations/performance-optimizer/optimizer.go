package performance

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"go.uber.org/zap"
)

// PerformanceOptimizer 性能优化器
type PerformanceOptimizer struct {
	collector  *MetricsCollector
	profiler   *Profiler
	pool       *ObjectPool
	workerPool *WorkerPool
	logger     *zap.Logger
	config     *OptimizerConfig
}

// OptimizerConfig 优化器配置
type OptimizerConfig struct {
	EnableProfiling   bool          `json:"enable_profiling"`
	EnableMetrics     bool          `json:"enable_metrics"`
	EnableObjectPool  bool          `json:"enable_object_pool"`
	EnableWorkerPool  bool          `json:"enable_worker_pool"`
	WorkerPoolSize    int           `json:"worker_pool_size"`
	ObjectPoolSize    int           `json:"object_pool_size"`
	ProfilingInterval time.Duration `json:"profiling_interval"`
	MetricsInterval   time.Duration `json:"metrics_interval"`
}

// NewPerformanceOptimizer 创建性能优化器
func NewPerformanceOptimizer(config *OptimizerConfig) *PerformanceOptimizer {
	optimizer := &PerformanceOptimizer{
		collector:  NewMetricsCollector(),
		profiler:   NewProfiler(),
		pool:       NewObjectPool(config.ObjectPoolSize),
		workerPool: NewWorkerPool(config.WorkerPoolSize),
		logger:     zap.L().Named("performance-optimizer"),
		config:     config,
	}

	if config.EnableMetrics {
		optimizer.startMetricsCollection()
	}

	if config.EnableProfiling {
		optimizer.startProfiling()
	}

	return optimizer
}

// startMetricsCollection 启动指标收集
func (po *PerformanceOptimizer) startMetricsCollection() {
	go func() {
		ticker := time.NewTicker(po.config.MetricsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				po.collectMetrics()
			}
		}
	}()
}

// startProfiling 启动性能分析
func (po *PerformanceOptimizer) startProfiling() {
	go func() {
		ticker := time.NewTicker(po.config.ProfilingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				po.collectProfiles()
			}
		}
	}()
}

// collectMetrics 收集指标
func (po *PerformanceOptimizer) collectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	po.collector.RecordMetric("memory.alloc", float64(m.Alloc))
	po.collector.RecordMetric("memory.total_alloc", float64(m.TotalAlloc))
	po.collector.RecordMetric("memory.sys", float64(m.Sys))
	po.collector.RecordMetric("memory.num_gc", float64(m.NumGC))
	po.collector.RecordMetric("goroutine.count", float64(runtime.NumGoroutine()))

	po.logger.Debug("metrics collected",
		zap.Uint64("alloc", m.Alloc),
		zap.Uint64("total_alloc", m.TotalAlloc),
		zap.Int("goroutines", runtime.NumGoroutine()))
}

// collectProfiles 收集性能分析数据
func (po *PerformanceOptimizer) collectProfiles() {
	po.profiler.WriteMemoryProfile("memory.prof")
	po.profiler.WriteGoroutineProfile("goroutine.prof")
	po.logger.Debug("profiles collected")
}

// GetMetrics 获取指标
func (po *PerformanceOptimizer) GetMetrics() map[string]float64 {
	return po.collector.GetMetrics()
}

// GetObjectPool 获取对象池
func (po *PerformanceOptimizer) GetObjectPool() *ObjectPool {
	return po.pool
}

// GetWorkerPool 获取工作池
func (po *PerformanceOptimizer) GetWorkerPool() *WorkerPool {
	return po.workerPool
}

// MonitorFunction 监控函数执行
func (po *PerformanceOptimizer) MonitorFunction(name string, fn func() error) error {
	start := time.Now()

	// 记录开始时的指标
	startGoroutines := runtime.NumGoroutine()
	var startMem runtime.MemStats
	runtime.ReadMemStats(&startMem)

	err := fn()

	// 记录结束时的指标
	duration := time.Since(start)
	endGoroutines := runtime.NumGoroutine()
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	// 记录指标
	po.collector.RecordMetric(fmt.Sprintf("function.%s.duration", name), duration.Seconds())
	po.collector.RecordMetric(fmt.Sprintf("function.%s.goroutines", name), float64(endGoroutines-startGoroutines))
	po.collector.RecordMetric(fmt.Sprintf("function.%s.memory", name), float64(endMem.Alloc-startMem.Alloc))

	if err != nil {
		po.collector.RecordMetric(fmt.Sprintf("function.%s.errors", name), 1)
		po.logger.Error("function execution failed",
			zap.String("name", name),
			zap.Duration("duration", duration),
			zap.Error(err))
	} else {
		po.collector.RecordMetric(fmt.Sprintf("function.%s.success", name), 1)
		po.logger.Debug("function executed",
			zap.String("name", name),
			zap.Duration("duration", duration))
	}

	return err
}

// ObjectPool 对象池
type ObjectPool struct {
	pool   chan interface{}
	new    func() interface{}
	reset  func(interface{})
	size   int
	logger *zap.Logger
}

// NewObjectPool 创建对象池
func NewObjectPool(size int) *ObjectPool {
	return &ObjectPool{
		pool:   make(chan interface{}, size),
		size:   size,
		logger: zap.L().Named("object-pool"),
	}
}

// SetFactory 设置对象工厂
func (op *ObjectPool) SetFactory(new func() interface{}, reset func(interface{})) {
	op.new = new
	op.reset = reset
}

// Get 获取对象
func (op *ObjectPool) Get() interface{} {
	select {
	case obj := <-op.pool:
		if op.reset != nil {
			op.reset(obj)
		}
		return obj
	default:
		if op.new != nil {
			return op.new()
		}
		return nil
	}
}

// Put 归还对象
func (op *ObjectPool) Put(obj interface{}) {
	if obj == nil {
		return
	}

	select {
	case op.pool <- obj:
		// 对象已归还到池中
	default:
		// 池已满，丢弃对象
		op.logger.Debug("object pool full, discarding object")
	}
}

// Size 获取池大小
func (op *ObjectPool) Size() int {
	return len(op.pool)
}

// Capacity 获取池容量
func (op *ObjectPool) Capacity() int {
	return op.size
}

// WorkerPool 工作池
type WorkerPool struct {
	workers  chan chan Job
	jobQueue chan Job
	quit     chan bool
	wg       sync.WaitGroup
	size     int
	logger   *zap.Logger
}

// Job 任务
type Job struct {
	ID       string
	Execute  func() error
	Callback func(error)
}

// NewWorkerPool 创建工作池
func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		workers:  make(chan chan Job, size),
		jobQueue: make(chan Job, size*2),
		quit:     make(chan bool),
		size:     size,
		logger:   zap.L().Named("worker-pool"),
	}

	for i := 0; i < size; i++ {
		worker := NewWorker(i, pool.workers)
		worker.Start()
	}

	go pool.dispatch()
	return pool
}

// Submit 提交任务
func (wp *WorkerPool) Submit(job Job) error {
	select {
	case wp.jobQueue <- job:
		wp.logger.Debug("job submitted", zap.String("id", job.ID))
		return nil
	default:
		wp.logger.Warn("job queue full, job rejected", zap.String("id", job.ID))
		return fmt.Errorf("job queue full")
	}
}

// SubmitAsync 异步提交任务
func (wp *WorkerPool) SubmitAsync(job Job) {
	go func() {
		if err := wp.Submit(job); err != nil {
			wp.logger.Error("failed to submit job", zap.String("id", job.ID), zap.Error(err))
		}
	}()
}

// dispatch 分发任务
func (wp *WorkerPool) dispatch() {
	for {
		select {
		case job := <-wp.jobQueue:
			go func() {
				worker := <-wp.workers
				worker <- job
			}()
		case <-wp.quit:
			return
		}
	}
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
	wp.logger.Info("worker pool stopped")
}

// Worker 工作者
type Worker struct {
	id       int
	jobQueue chan Job
	workers  chan chan Job
	quit     chan bool
	logger   *zap.Logger
}

// NewWorker 创建工作者
func NewWorker(id int, workers chan chan Job) *Worker {
	return &Worker{
		id:       id,
		jobQueue: make(chan Job),
		workers:  workers,
		quit:     make(chan bool),
		logger:   zap.L().Named(fmt.Sprintf("worker-%d", id)),
	}
}

// Start 启动工作者
func (w *Worker) Start() {
	go func() {
		for {
			w.workers <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				w.logger.Debug("processing job", zap.String("id", job.ID))

				start := time.Now()
				err := job.Execute()
				duration := time.Since(start)

				if err != nil {
					w.logger.Error("job execution failed",
						zap.String("id", job.ID),
						zap.Duration("duration", duration),
						zap.Error(err))
				} else {
					w.logger.Debug("job completed",
						zap.String("id", job.ID),
						zap.Duration("duration", duration))
				}

				if job.Callback != nil {
					job.Callback(err)
				}

			case <-w.quit:
				w.logger.Info("worker stopped")
				return
			}
		}
	}()
}

// Stop 停止工作者
func (w *Worker) Stop() {
	close(w.quit)
}

// LockFreeQueue 无锁队列
type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// node 节点
type node struct {
	value interface{}
	next  unsafe.Pointer
}

// NewLockFreeQueue 创建无锁队列
func NewLockFreeQueue() *LockFreeQueue {
	n := &node{}
	return &LockFreeQueue{
		head: unsafe.Pointer(n),
		tail: unsafe.Pointer(n),
	}
}

// Enqueue 入队
func (q *LockFreeQueue) Enqueue(value interface{}) {
	n := &node{value: value}

	for {
		tail := (*node)(atomic.LoadPointer(&q.tail))
		next := (*node)(atomic.LoadPointer(&tail.next))

		if tail == (*node)(atomic.LoadPointer(&q.tail)) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&tail.next, nil, unsafe.Pointer(n)) {
					atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(n))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			}
		}
	}
}

// Dequeue 出队
func (q *LockFreeQueue) Dequeue() interface{} {
	for {
		head := (*node)(atomic.LoadPointer(&q.head))
		tail := (*node)(atomic.LoadPointer(&q.tail))
		next := (*node)(atomic.LoadPointer(&head.next))

		if head == (*node)(atomic.LoadPointer(&q.head)) {
			if head == tail {
				if next == nil {
					return nil
				}
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			} else {
				value := next.value
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
					return value
				}
			}
		}
	}
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
	logger  *zap.Logger
}

// Metric 指标
type Metric struct {
	Name  string
	Value float64
	Type  MetricType
	Count int64
	Sum   float64
	Min   float64
	Max   float64
	mu    sync.RWMutex
}

// MetricType 指标类型
type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
)

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
		logger:  zap.L().Named("metrics-collector"),
	}
}

// RecordMetric 记录指标
func (mc *MetricsCollector) RecordMetric(name string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		metric = &Metric{
			Name: name,
			Type: MetricTypeGauge,
		}
		mc.metrics[name] = metric
	}

	metric.mu.Lock()
	defer metric.mu.Unlock()

	metric.Value = value
	metric.Count++
	metric.Sum += value

	if metric.Count == 1 {
		metric.Min = value
		metric.Max = value
	} else {
		if value < metric.Min {
			metric.Min = value
		}
		if value > metric.Max {
			metric.Max = value
		}
	}
}

// GetMetrics 获取指标
func (mc *MetricsCollector) GetMetrics() map[string]float64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := make(map[string]float64)
	for name, metric := range mc.metrics {
		metric.mu.RLock()
		metrics[name] = metric.Value
		metric.mu.RUnlock()
	}

	return metrics
}

// GetMetricDetails 获取指标详情
func (mc *MetricsCollector) GetMetricDetails(name string) *Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return mc.metrics[name]
}

// Profiler 性能分析器
type Profiler struct {
	logger *zap.Logger
}

// NewProfiler 创建性能分析器
func NewProfiler() *Profiler {
	return &Profiler{
		logger: zap.L().Named("profiler"),
	}
}

// WriteMemoryProfile 写入内存分析
func (p *Profiler) WriteMemoryProfile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create memory profile: %w", err)
	}
	defer file.Close()

	if err := pprof.WriteHeapProfile(file); err != nil {
		return fmt.Errorf("failed to write memory profile: %w", err)
	}

	p.logger.Debug("memory profile written", zap.String("file", filename))
	return nil
}

// WriteGoroutineProfile 写入goroutine分析
func (p *Profiler) WriteGoroutineProfile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create goroutine profile: %w", err)
	}
	defer file.Close()

	if err := pprof.Lookup("goroutine").WriteTo(file, 0); err != nil {
		return fmt.Errorf("failed to write goroutine profile: %w", err)
	}

	p.logger.Debug("goroutine profile written", zap.String("file", filename))
	return nil
}

// StartCPUProfile 开始CPU分析
func (p *Profiler) StartCPUProfile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CPU profile: %w", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		file.Close()
		return fmt.Errorf("failed to start CPU profile: %w", err)
	}

	p.logger.Info("CPU profiling started", zap.String("file", filename))
	return nil
}

// StopCPUProfile 停止CPU分析
func (p *Profiler) StopCPUProfile() {
	pprof.StopCPUProfile()
	p.logger.Info("CPU profiling stopped")
}
