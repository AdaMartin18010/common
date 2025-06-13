// SubPub另外一种实现
//	只使用golang的基础语言特性

package eventchans

import (
	"context"
	"time"
)

// EventChans is interface for global (subscribe, publish, control) EventChans behavior
type EventChans interface {
	ChanController
	ChanSubscriber
	ChanPublisher
}

// ChanSubscriber defines subscription-related EventChans behavior
// <-chan any is buffered chan
type ChanSubscriber interface {
	Subscribe(topic string) <-chan any
	UnSubscribe(topic string, ch <-chan any) error
}

// ChanPublisher defines publishing-related EventChans behavior
type ChanPublisher interface {
	Publish(topic string, msgs ...any) bool
	PublishAsync(ctx context.Context, tm time.Duration, topic string, msgs ...any) error
}

// ChanController defines EventChans control behavior (checking topic presence, synchronization)
type ChanController interface {
	Topics() string
	HasChansLen(topic string) int
	Close()
	WaitAsync()
}

// todo : implement channel topic  to publish async 可以指定独立的goroutine来分发消息
