package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker 工作器结构
type Worker struct {
	ID       int
	JobChan  chan Job
	QuitChan chan bool
	Wg       *sync.WaitGroup
}

// Job 任务结构
type Job struct {
	ID       int
	Data     string
	Duration time.Duration
}

// NewWorker 创建新的工作器
func NewWorker(id int, jobChan chan Job, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:       id,
		JobChan:  jobChan,
		QuitChan: make(chan bool),
		Wg:       wg,
	}
}

// Start 启动工作器
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobChan:
				w.processJob(job)
			case <-w.QuitChan:
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop 停止工作器
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// processJob 处理任务
func (w *Worker) processJob(job Job) {
	defer w.Wg.Done()
	fmt.Printf("Worker %d processing job %d: %s\n", w.ID, job.ID, job.Data)
	time.Sleep(job.Duration)
	fmt.Printf("Worker %d completed job %d\n", w.ID, job.ID)
}

// WorkerPool 工作池结构
type WorkerPool struct {
	Workers    []*Worker
	JobChan    chan Job
	Wg         *sync.WaitGroup
	NumWorkers int
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(numWorkers int, jobChanSize int) *WorkerPool {
	pool := &WorkerPool{
		Workers:    make([]*Worker, numWorkers),
		JobChan:    make(chan Job, jobChanSize),
		Wg:         &sync.WaitGroup{},
		NumWorkers: numWorkers,
	}

	// 创建并启动工作器
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, pool.JobChan, pool.Wg)
		pool.Workers[i] = worker
		worker.Start()
	}

	return pool
}

// SubmitJob 提交任务
func (wp *WorkerPool) SubmitJob(job Job) {
	wp.Wg.Add(1)
	wp.JobChan <- job
}

// Wait 等待所有任务完成
func (wp *WorkerPool) Wait() {
	wp.Wg.Wait()
}

// Stop 停止所有工作器
func (wp *WorkerPool) Stop() {
	for _, worker := range wp.Workers {
		worker.Stop()
	}
}

func main() {
	// 创建工作池
	pool := NewWorkerPool(3, 10)
	defer pool.Stop()

	// 提交任务
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:       i,
			Data:     fmt.Sprintf("Task data %d", i),
			Duration: time.Duration(i*100) * time.Millisecond,
		}
		pool.SubmitJob(job)
	}

	// 等待所有任务完成
	pool.Wait()
	fmt.Println("All jobs completed!")
}
