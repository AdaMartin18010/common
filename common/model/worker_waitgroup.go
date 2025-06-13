package common

import (
	"fmt"
	"runtime"
	"sync"
)

// 无交互的goroutine接口
type Worker interface {
	Work() error
}

type Recover interface {
	Recover()
}

type WorkerRecover interface {
	Worker
	Recover
}

// type WorkerWrapper interface {
// 	StartingWait(worker WorkerRecover)
// 	Started() bool
// }

type WorkerWG struct {
	wg              *sync.WaitGroup
	startWaiting    chan struct{}
	startChanClosed bool
	wrwm            *sync.RWMutex // guards worker
	wm              *sync.Mutex   // guards waitgroup wait and add(+n)
}

func NewWorkerWG() *WorkerWG {
	return &WorkerWG{
		wg:   &sync.WaitGroup{},
		wrwm: &sync.RWMutex{},
		wm:   &sync.Mutex{},
	}
}

func (w *WorkerWG) DebugInfo() string {
	w.wrwm.RLock()
	defer w.wrwm.RUnlock()
	return fmt.Sprintf("wg:%+v,startWaiting:%+v,startChanClosed:%t", w.wg, w.startWaiting, w.startChanClosed)
}

// FanOut Implements
// Only for the same goroutine invoking these methods
func (w *WorkerWG) StartingWait(worker WorkerRecover) {
	w.wrwm.Lock()
	if w.startWaiting == nil {
		w.startWaiting = make(chan struct{}, 1)
		w.startChanClosed = false
	}
	w.wrwm.Unlock()

	//guarding wait() and add(+n)
	w.wm.Lock()
	defer w.wm.Unlock()
	w.wg.Add(1)

	go func() {
		//protect workers panic wg.Done() is not executed
		defer w.wg.Done()
		runtime.Gosched()
		startchan := (<-chan struct{})(nil)
		w.wrwm.RLock()
		if w.startWaiting != nil && !w.startChanClosed {
			startchan = w.startWaiting
		}
		w.wrwm.RUnlock()
		if startchan != nil {
			<-startchan
		}
		defer worker.Recover()
		runtime.Gosched()
		worker.Work()
	}()
}

func (w *WorkerWG) StartAsync() {
	w.wrwm.Lock()
	defer w.wrwm.Unlock()
	if w.startWaiting != nil && !w.startChanClosed {
		close(w.startWaiting)
		w.startChanClosed = true
	}
}

func (w *WorkerWG) WaitAsync() {
	//guarding wait() and add(+n)
	w.wm.Lock()
	w.wg.Wait()
	w.wm.Unlock()

	// don't lock or wait on mutex for long time.
	w.wrwm.Lock()
	defer w.wrwm.Unlock()
	if w.startWaiting != nil && w.startChanClosed {
		w.startWaiting = nil
		w.startChanClosed = false
	}
}
