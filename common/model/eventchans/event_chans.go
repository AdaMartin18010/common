package eventchans

import (
	"context"
	"errors"
	"sync"
	"time"

	mdl "common/model"
)

const (
	defaultChanBufferSize = 10
)

var (
	//Verify EvtChans satisfied interface - EventChans
	_ EventChans = (*EvtChans)(nil)
)

var (
	ErrTopicEmpty       = errors.New("evtchans: topic empty")
	ErrTopicChanNotFind = errors.New("evtchans: topic and chan is not found")
	ErrChanNil          = errors.New("evtchans: chan is nil")
	ErrChansClose       = errors.New("evtchans: events chan closed")
	ErrAsyncTimeOut     = errors.New("evtchans: time out")
)

// 只封装最简单最基础的应用 每个订阅者都需要判断chan的是否close
// 订阅和发布者都有较大的自由空间来控制发布订阅策略
// 消息分发可以选择 每个主题 一个goroutine分发模式
type EvtChans struct {
	wg *sync.WaitGroup // todo: implement work events

	rwmu       *sync.RWMutex // guards map[string][]chan any
	subs       map[string][]chan any
	chanBufLen uint

	// mu     *sync.Mutex // guards closed
	closed bool
}

func NewEvtChans(chanlen uint) *EvtChans {
	evtcs := &EvtChans{
		rwmu: &sync.RWMutex{},
		//mu:         &sync.Mutex{},
		wg:         &sync.WaitGroup{},
		chanBufLen: chanlen,
		subs:       make(map[string][]chan any),
	}

	if evtcs.chanBufLen < defaultChanBufferSize {
		evtcs.chanBufLen = defaultChanBufferSize
	}

	return evtcs
}

// create a new channel for the given topic
func (ecs *EvtChans) Subscribe(topic string) <-chan any {
	if topic == "" {
		return nil
	}

	ecs.rwmu.Lock()
	clsd := ecs.closed
	ecs.rwmu.Unlock()

	if clsd {
		return nil
	}

	ecs.rwmu.Lock()
	defer ecs.rwmu.Unlock()
	ch := make(chan any, ecs.chanBufLen)
	// 	// todo: check this is race condition
	// 	==================
	// WARNING: DATA RACE
	// Write at 0x00c000212090 by goroutine 162:
	//   runtime.mapassign_faststr()
	//       C:/_program_swap_/Go/src/runtime/map_faststr.go:203 +0x0
	//   navigate/common/model/eventchans.(*EvtChans).Subscribe()
	//       X:/gopath/src/navigate/navigate/common/model/eventchans/event_chans.go:85 +0x2d2
	// ????????????? golang 1.20.2
	ecs.subs[topic] = append(ecs.subs[topic], ch)
	ecs.wg.Add(1)
	return ch
}

// Unsubscribe from topic and the event channel
func (ecs *EvtChans) UnSubscribe(topic string, ch <-chan any) error {
	if topic == "" {
		return ErrTopicEmpty
	}
	if ch == nil {
		return ErrChanNil
	}
	ecs.rwmu.Lock()
	defer ecs.rwmu.Unlock()
	//remove all equal channel that are equal to the topic and channel
	for i := 0; i < len(ecs.subs[topic]); i++ {
		if ecs.subs[topic][i] == ch {
			ecs.subs[topic] = append(ecs.subs[topic][:i], ecs.subs[topic][i+1:]...)
			i--
			ecs.wg.Done()
			if len(ecs.subs[topic]) == 0 {
				delete(ecs.subs, topic)
			}
			return nil
		}
	}
	return ErrTopicChanNotFind
}

// 如果整个EvtChans关闭 不再发送消息 返回false
func (ecs *EvtChans) Publish(topic string, msgs ...any) bool {
	if topic == "" {
		return false
	}

	if len(msgs) == 0 {
		return false
	}

	ecs.rwmu.Lock()
	clsd := ecs.closed
	ecs.rwmu.Unlock()

	if clsd {
		return false
	}

	ecs.rwmu.RLock()
	defer ecs.rwmu.RUnlock()
	for _, ch := range ecs.subs[topic] {
		for _, msg := range msgs {
			ch <- msg
		}
	}
	return true
}

func (ecs *EvtChans) PublishAsync(ctx context.Context, tm time.Duration, topic string, msgs ...any) error {
	if topic == "" {
		return ErrTopicEmpty
	}
	if len(msgs) == 0 {
		return nil
	}

	ecs.rwmu.Lock()
	clsd := ecs.closed
	ecs.rwmu.Unlock()
	if clsd {
		return ErrChansClose
	}

	timer := mdl.TimerPool.Get(tm)
	defer mdl.TimerPool.Put(timer)
	ecs.rwmu.RLock()
	defer ecs.rwmu.RUnlock()
	for _, ch := range ecs.subs[topic] {
		for _, msg := range msgs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- msg:
			case <-timer.C:
				return ErrAsyncTimeOut
			}
		}
	}
	return nil
}

func (ecs *EvtChans) Topics() string {
	type topics struct {
		Topic string `json:"topic"`
		Cap   int    `json:"cap"`
		Len   int    `json:"len"`
	}
	type topicJson struct {
		Topics []topics `json:"Topics"`
	}

	sl := topicJson{Topics: []topics{}}

	ecs.rwmu.RLock()
	for k := range ecs.subs {
		sl.Topics = append(sl.Topics, topics{Topic: k, Cap: ecs.topicCap(k), Len: ecs.topicLen(k)})
	}
	ecs.rwmu.RUnlock()
	rs, _ := mdl.Json.MarshalToString(sl)
	return rs
}

func (ecs *EvtChans) HasChansLen(topic string) int {
	ecs.rwmu.RLock()
	defer ecs.rwmu.RUnlock()
	return ecs.topicLen(topic)
}

// if not find return -1
func (ecs *EvtChans) topicLen(topic string) int {
	tpcs, ok := ecs.subs[topic]
	if ok {
		return len(tpcs)
	}
	return -1
}

// if not find return -1
func (ecs *EvtChans) topicCap(topic string) int {
	tpcs, ok := ecs.subs[topic]
	if ok {
		return cap(tpcs)
	}
	return -1
}

func (ecs *EvtChans) Close() {
	ecs.rwmu.Lock()
	clsd := ecs.closed
	ecs.rwmu.Unlock()
	if clsd {
		return
	}

	ecs.rwmu.Lock()
	for _, subs := range ecs.subs {
		for _, ch := range subs {
			close(ch)
		}
	}
	ecs.closed = true
	ecs.rwmu.Unlock()
}

func (ecs *EvtChans) WaitAsync() {
	// Wait until all channels are UnSubscribed
	ecs.wg.Wait()
	// ecs.rwmu.RLock()
	// defer ecs.rwmu.RUnlock()
}
