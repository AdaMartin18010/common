package token

import (
	"errors"
	"runtime"
	"sync"
	"time"

	cm "common/model"
)

const (
// waitForToken = 150 * time.Millisecond
)

var (
	ErrWaitedTimeOut = errors.New("timeout waiting for token")
)

// Token defines the interface for the tokens used to indicate when
// actions have completed.
type Token interface {
	// Done returns a channel that is closed when the flow associated
	// with the Token completes. Clients should call Error after the
	// channel is closed to check if the flow completed successfully.
	//
	// Done is provided for use in select statements. Simple use cases may
	// use Wait or WaitTimeout.
	Done() <-chan struct{}
	Reset()
}

type TokenSetter interface {
	// Set error of  the error. And it's not be completed.
	SetErr(error)
	// Set completed state of the Token.
	Completed()
}

type TokenGetter interface {
	// if it's not completed, get the error
	Err() error
	// Wait will wait indefinitely for the Token to complete, ie the Publish
	// to be sent and confirmed receipt from the broker.
	Wait() bool
	// WaitTimeout takes a time.Duration to wait for the flow associated with the
	// Token to complete, returns true if it returned before the timeout or
	// returns false if the timeout occurred. In the case of a timeout the Token
	// does not have an error set in case the caller wishes to wait again.
	WaitTimeout(time.Duration) bool
}

/***********************************************************************
TokenCompletor 的行为模式:
创建:
 1. 不可在同一个goroutine中使用 否则堵塞 无有效方法检测
 2. 创建者可以重复使用Reset创建新的token

设置:
 1. 成功只有Completed()函数调用
 2. 失败 直接调用 SetErr(error)

等待:
 1. 可使用Done()获取chan 主动select
 2. 可以使用Wait 函数组-阻塞 和 非阻塞
 	非阻塞-超时逻辑:
		1. 超时会自动设置错误ErrWaitedTimeOut同时设置完成-失败
		2. 如果后续完成,再次使用Wait函数组 依然可以获取到成功的状态
 3. 完成-成功: 只有时间上 第一次Wait获取到true  后续的Wait都会获取到false
 4. 完成-失败: 设置者使用SetError(error) 所有的Wait都会获取到false
**************************************************************************/

type TokenCompletor interface {
	Token
	TokenSetter
	TokenGetter
}

type BaseToken struct {
	rwl      *sync.RWMutex
	complete chan struct{}
	err      error
}

func NewBaseToken() *BaseToken {
	return &BaseToken{
		rwl:      &sync.RWMutex{},
		complete: make(chan struct{}),
	}
}

func (b *BaseToken) Reset() {
	b.rwl.Lock()
	defer b.rwl.Unlock()
	//-- for runtime GC only
	b.flowComplete()
	runtime.Gosched()
	b.complete = nil
	//--
	b.complete = make(chan struct{})
	b.err = nil
}

// Wait implements the Token Wait method.
// return completed value meaning that Work is completed Or not.
func (b *BaseToken) Wait() bool {
	if b.complete == nil {
		return false
	}
	_, ok := <-b.complete
	return ok
}

// WaitTimeout implements the Token WaitTimeout method.
// return completed value meaning that Work is completed Or not.
func (b *BaseToken) WaitTimeout(d time.Duration) bool {
	if b.complete == nil {
		return false
	}

	timer := cm.TimerPool.Get(d)
	defer cm.TimerPool.Put(timer)
	select {
	case _, ok := <-b.complete:
		if !timer.Stop() {
			<-timer.C
		}
		return ok
	case <-timer.C:
		if b.err == nil {
			b.err = ErrWaitedTimeOut
		}
		return false
	}
}

func (b *BaseToken) Done() <-chan struct{} {
	b.rwl.Lock()
	defer b.rwl.Unlock()
	return b.complete
}

func (b *BaseToken) flowComplete() {
	if b.complete != nil {
		select {
		case <-b.complete:
		default:
			close(b.complete)
		}
	}
}

func (b *BaseToken) Err() error {
	b.rwl.RLock()
	defer b.rwl.RUnlock()
	err := b.err
	b.err = nil
	return err
}

func (b *BaseToken) SetErr(e error) {
	b.rwl.Lock()
	defer b.rwl.Unlock()
	b.err = e
	b.flowComplete()
}

func (b *BaseToken) Completed() {
	b.rwl.Lock()
	defer b.rwl.Unlock()
	if b.complete != nil {
		b.complete <- struct{}{}
		close(b.complete)
	}
}
