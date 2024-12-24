package common

import (
	"context"
	"fmt"
	"sync"
	"time"

	tmrp "common/model/timerpool"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	// publish  use zap.NewProduction() that error print structure is different.
	L, _ = zap.NewDevelopment()
	//L, _ = zap.NewProduction()
	TimerPool = &tmrp.TimerPool{}
	Json      = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {

}

// 组件component的控制结构 component内部控制和外界控制
// context cancel pair with control flow
type CtrlSt struct {
	c   context.Context
	ccl context.CancelFunc
	wwg *WorkerWG

	rwm *sync.RWMutex
}

func NewCtrlSt(ctx context.Context) *CtrlSt {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)
	return &CtrlSt{
		c:   ctx,
		ccl: cancel,
		wwg: NewWorkerWG(),
		rwm: &sync.RWMutex{},
	}
}

func (cs *CtrlSt) DebugInfo() string {
	cs.rwm.RLock()
	// don't use # format string placeholders because that raises runtime race condition
	ss := fmt.Sprintf("(ctrl)[Ctx:%+v;Ccl:%+v;Wg:%s]", cs.c, cs.ccl, cs.wwg.DebugInfo())
	cs.rwm.RUnlock()
	return ss
}

func (cs *CtrlSt) Cancel() {
	cs.rwm.Lock()
	defer cs.rwm.Unlock()
	if cs.ccl != nil {
		cs.ccl()
	}
}

func (cs *CtrlSt) Context() context.Context {
	cs.rwm.RLock()
	defer cs.rwm.RUnlock()
	return cs.c
}

func (cs *CtrlSt) WaitGroup() *WorkerWG {
	cs.rwm.RLock()
	defer cs.rwm.RUnlock()
	return cs.wwg
}

func (cs *CtrlSt) ForkCtxWg() *CtrlSt {
	cs.rwm.RLock()
	defer cs.rwm.RUnlock()
	ctx, cancel := context.WithCancel(cs.c)
	return &CtrlSt{
		c:   ctx,
		ccl: cancel,
		wwg: cs.wwg,
		rwm: &sync.RWMutex{},
	}
}

func (cs *CtrlSt) ForkCtxWgTimeout(tm time.Duration) *CtrlSt {
	cs.rwm.RLock()
	defer cs.rwm.RUnlock()

	ctx, cancel := context.WithTimeout(cs.c, tm)
	return &CtrlSt{
		c:   ctx,
		ccl: cancel,
		wwg: cs.wwg,
		rwm: &sync.RWMutex{},
	}
}

// implemented just for testing
func (cs *CtrlSt) WithCtrl(ctrl *CtrlSt) *CtrlSt {
	cs.rwm.Lock()
	defer cs.rwm.Unlock()
	cs = ctrl
	return cs
}

func (cs *CtrlSt) WithCtx(ctx context.Context) *CtrlSt {
	cs.rwm.Lock()
	defer cs.rwm.Unlock()
	ctx0, cancel0 := context.WithCancel(ctx)
	cs.c = ctx0
	cs.ccl = cancel0
	return cs
}

func (cs *CtrlSt) WithTimeout(ctx context.Context, tm time.Duration) *CtrlSt {
	cs.rwm.Lock()
	defer cs.rwm.Unlock()
	ctx, cancel := context.WithTimeout(ctx, tm)
	cs.c = ctx
	cs.ccl = cancel
	return cs
}
