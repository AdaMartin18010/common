package component

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	mdl "common/model"

	uuid "github.com/google/uuid"
)

var (
	//Verify Satisfies interfaces
	_ CptRoot           = (*CptMetaSt)(nil)
	_ Cpt               = (*CptMetaSt)(nil)
	_ mdl.WorkerRecover = (*CptMetaSt)(nil)
)

type CptMetaSt struct {
	mu    *sync.Mutex
	ctlSt *mdl.CtrlSt

	IdStr   IdName
	KindStr KindName
	State   *atomic.Value
	mdl.WorkerRecover
}

// 该函数会直接copy创建cm.ControlStruct
// Accepted type: IdName KindName context.Context context.CancelFunc or *common.CommonStruct
func NewCpt(v ...any) *CptMetaSt {
	cpbd := &CptMetaSt{
		mu:    &sync.Mutex{},
		ctlSt: nil,
		State: &atomic.Value{},
	}

	for i := range v {
		switch v[i].(type) {
		case IdName:
			cpbd.IdStr = v[i].(IdName)
		case KindName:
			cpbd.KindStr = v[i].(KindName)
		// case context.Context:
		// 	cpbd.CtlSt.ctx = v[i].(context.Context)
		case mdl.WorkerRecover:
			cpbd.WorkerRecover = v[i].(mdl.WorkerRecover)
		case *mdl.CtrlSt:
			cpbd.mu.Lock()
			cpbd.ctlSt = v[i].(*mdl.CtrlSt)
			cpbd.mu.Unlock()
		}
	}

	if cpbd.WorkerRecover == nil {
		cpbd.WorkerRecover = cpbd
	}

	cpbd.State.Store(false)
	if len(cpbd.KindStr) == 0 {
		cpbd.reflectKind()
	}

	if len(cpbd.IdStr) == 0 {
		cpbd.Id()
	}

	mdl.L.Sugar().Debugf("component-%s:initialized.", cpbd.CmptInfo())
	return cpbd
}

// 该函数会自己检查和创建 cm.ControlStruct有默认的行为
// Accepted type: IdName KindName context.Context context.CancelFunc or *common.CommonStruct
func NewCptMetaSt(v ...any) *CptMetaSt {
	cpbd := &CptMetaSt{
		mu:    &sync.Mutex{},
		ctlSt: nil,
		State: &atomic.Value{},
	}

	for i := range v {
		switch v[i].(type) {
		case IdName:
			cpbd.IdStr = v[i].(IdName)
		case KindName:
			cpbd.KindStr = v[i].(KindName)
		// case context.Context:
		// 	cpbd.ctlSt.Ctx = v[i].(context.Context)
		// case context.CancelFunc:
		// 	cpbd.ctlSt.Cancel = v[i].(context.CancelFunc)
		// case *sync.WaitGroup:
		// 	cpbd.ctlSt.WorkerWaitGroup = v[i].(*sync.WaitGroup)
		case mdl.WorkerRecover:
			cpbd.WorkerRecover = v[i].(mdl.WorkerRecover)
		case *mdl.CtrlSt:
			cpbd.mu.Lock()
			cpbd.ctlSt = v[i].(*mdl.CtrlSt)
			cpbd.mu.Unlock()
		}
	}

	//这部分的逻辑 可以直接在组件外部 创建后赋值
	//内部额外的逻辑 只是为了节省必要的工作代码
	if cpbd.ctlSt == nil {
		cpbd.ctlSt = mdl.NewCtrlSt(context.Background())
	}
	// else
	// if cpbd.ctlSt.Context() == nil {
	// 	cpbd.ctlSt.WithCtx(context.Background())
	// } else {
	// 	//初始化如果ctx有值则直接使用其生成ctx cancel 用于组件自身的控制
	// 	cpbd.ctlSt.WithCtx(cpbd.CtlSt.Context())
	// }
	// // 生成默认的waitgroup 否则使用传入的
	// if cpbd.ctlSt.WorkerWaitGroup == nil {
	// 	cpbd.ctlSt.WorkerWaitGroup = cm.NewWorkerWaitGroup()
	// }

	if cpbd.WorkerRecover == nil {
		cpbd.WorkerRecover = cpbd
	}

	cpbd.State.Store(false)

	if len(cpbd.KindStr) == 0 {
		cpbd.reflectKind()
	}

	if len(cpbd.IdStr) == 0 {
		cpbd.Id()
	}

	mdl.L.Sugar().Debugf("%s-%s:initialized.", cpbd.CmptInfo(), cpbd.Ctrl().DebugInfo())
	return cpbd
}

func (cpbd *CptMetaSt) CmptInfo() string {
	return fmt.Sprintf("(cmpt)[Kd:%s,Id:%s]", cpbd.KindStr, cpbd.IdStr)
}

func (cpbd *CptMetaSt) Ctrl() *mdl.CtrlSt {
	cpbd.mu.Lock()
	defer cpbd.mu.Unlock()
	return cpbd.ctlSt

}

func (cpbd *CptMetaSt) Id() IdName {
	if len(cpbd.IdStr) == 0 {
		if len(cpbd.KindStr) == 0 {
			cpbd.reflectKind()
		}

		if uuid, err := uuid.NewUUID(); err == nil {
			cpbd.IdStr = IdName(uuid.String())
		} else {
			mdl.L.Sugar().Debugf("Error generating id: %+v", err)
			cpbd.IdStr = (IdName)(fmt.Sprintf("%s_%X", cpbd.Kind(), rand.Intn(int(^uint(0)>>1))))
		}
	}

	return cpbd.IdStr
}

func (cpbd *CptMetaSt) reflectKind() {
	cpbd.KindStr = KindName(reflect.TypeOf(cpbd).Name())
	if len(cpbd.KindStr) == 0 {
		cpbd.KindStr = KindName(reflect.TypeOf(cpbd).String())
	}
}

func (cpbd *CptMetaSt) Kind() KindName {
	if cpbd.KindStr == "" {
		cpbd.reflectKind()
	}

	return cpbd.KindStr
}

func (cpbd *CptMetaSt) IsRunning() bool {
	return cpbd.State.Load().(bool)
}

// todo: implement this method of each Component
func (cpbd *CptMetaSt) Work() (err error) {
	<-cpbd.Ctrl().Context().Done()
	if err = cpbd.Ctrl().Context().Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			mdl.L.Sugar().Debugf("Components-%s  Work timeout error : %+v", cpbd.CmptInfo(), err)
			return nil
		}
		return err
	}
	return nil
}

func (cpbd *CptMetaSt) Recover() {
	if rc := recover(); rc != nil {
		var buf [8196]byte
		//只打印出本golang内部的调用栈
		n := runtime.Stack(buf[:], false)
		mdl.L.Sugar().Warnf(`%s Worker Recover :%+v,Stack Trace: %s`, cpbd.CmptInfo(), rc, buf[:n])
	}
}

func (cpbd *CptMetaSt) Start() error {
	if cpbd.State.Load().(bool) {
		return fmt.Errorf("component:%s is already Running", cpbd.CmptInfo())
	}

	if cpbd.WorkerRecover != nil {
		cpbd.Ctrl().WaitGroup().StartingWait(cpbd.WorkerRecover)
	} else {
		cpbd.Ctrl().WaitGroup().StartingWait(cpbd)
	}
	cpbd.Ctrl().WaitGroup().StartAsync()
	cpbd.State.Store(true)
	return nil
}

func (cpbd *CptMetaSt) Stop() error {
	if !cpbd.State.Load().(bool) {
		return fmt.Errorf("component:%s is already Stopping", cpbd.CmptInfo())
	}
	cpbd.Ctrl().Cancel()
	<-cpbd.Ctrl().Context().Done()
	cpbd.State.Store(false)
	return nil
}

func (cpbd *CptMetaSt) Finalize() error {
	<-cpbd.Ctrl().Context().Done()
	if err := cpbd.Ctrl().Context().Err(); err != nil {
		if !(errors.Is(err, context.Canceled) && errors.Is(err, context.DeadlineExceeded)) {
			return err
		}
	}
	cpbd.Ctrl().WaitGroup().WaitAsync()
	return nil
}
