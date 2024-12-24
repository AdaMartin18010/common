package eventchans_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	ec "common/model/eventchans"

	"go.uber.org/goleak"
)

func TestChanCapLen(t *testing.T) {
	ch := make(chan any, 1000)
	t.Logf("ch cap:%d,len:%d", cap(ch), len(ch))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	ecs := ec.NewEvtChans(1)
	if ecs == nil {
		t.Log("New EventChans not created!")
		t.Fail()
	}
}

func TestPubSync(t *testing.T) {
	ecs := ec.NewEvtChans(1000)
	topic1 := "topic1"
	topic2 := "topic2"
	//wg := &sync.WaitGroup{}
	beginChan := make(chan chan struct{}, 1)
	endChan := make(chan chan struct{}, 1)

	subSyncFunc := func(topic string, t *testing.T) {
		//defer wg.Done()
		<-beginChan
		c := ecs.Subscribe(topic)
		if c == nil {
			t.Logf("%#v:%#v", topic, c)
		} else {
			defer ecs.UnSubscribe(topic, c)
			t.Logf("evinfo:%#v", ecs.Topics())
			// 同步发送 如果缓存小于发布者的总体消息数 则阻塞
			ecs.Publish(topic, nil)
			for {
				select {
				case <-c:
				case <-endChan:
					{
						return
					}
				}
			}
		}
	}
	for i := 0; i < 1000; i++ {
		//wg.Add(1)
		go subSyncFunc(topic1, t)
	}
	for i := 0; i < 1000; i++ {
		//wg.Add(1)
		go subSyncFunc(topic2, t)
	}

	close(beginChan)
	t.Log("Begin")
	time.Sleep(3 * time.Second)
	close(endChan)
	t.Log("End")

	//wg.Wait()
	ecs.WaitAsync()
	t.Logf("evinfo:%#v", ecs.Topics())
}

func TestSubscribe(t *testing.T) {
	ecs := ec.NewEvtChans(1)
	if ecs.Subscribe("topic") == nil {
		t.Fail()
	}

	if ecs.Subscribe("topic") == nil {
		t.Fail()
	}

	if ecs.HasChansLen("topic") == 0 {
		t.Fail()
	}

	if ecs.HasChansLen("topic") != 2 {
		t.Fail()
	}
}

func TestSubUnSub(t *testing.T) {
	buflen := 10
	ecs := ec.NewEvtChans(uint(buflen))
	event := "topic"
	chan0 := ecs.Subscribe(event)
	chan1 := ecs.Subscribe(event)
	ecs.UnSubscribe(event, chan0)
	ecs.UnSubscribe(event, chan1)
	ecs.WaitAsync()
	t.Log("OK")
}

func TestNullClose(t *testing.T) {
	buflen := 10
	ecs := ec.NewEvtChans(uint(buflen))
	ecs.Close()
	ecs.WaitAsync()
	t.Log("OK")
}

func TestClose(t *testing.T) {
	buflen := 10
	ecs := ec.NewEvtChans(uint(buflen))
	event := "topic"
	chan0 := ecs.Subscribe(event)
	chan1 := ecs.Subscribe(event)
	ecs.Publish(event, nil)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		chan0closed := false
		chan1closed := false
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case v, ok := <-chan0:
				{
					if !ok {
						if !chan0closed {
							ecs.UnSubscribe(event, chan0)
							chan0closed = true
						} else {
							continue
						}
					} else {
						t.Logf("chan0 value: %+v", v)
					}

				}
			case v, ok := <-chan1:
				{
					if !ok {
						if !chan1closed {
							ecs.UnSubscribe(event, chan1)
							chan1closed = true
						} else {
							continue
						}
					} else {
						t.Logf("chan1 value: %+v", v)
					}
				}
			case <-ticker.C:
				{
					if chan0closed && chan1closed {
						t.Log("chan0 chan1 all close.")
						return
					}
				}
			}
		}
	}()
	ecs.Close()
	ecs.WaitAsync()
	t.Log("OK")
	wg.Wait()
}

func TestSubUnSubPub(t *testing.T) {
	buflen := 10
	ecs := ec.NewEvtChans(uint(buflen))
	event := "topic"
	chan0 := ecs.Subscribe(event)
	chan1 := ecs.Subscribe(event)
	ecs.Publish(event, nil)
	go func() {
		chan0closed := false
		chan1closed := false
		ticker := time.NewTicker(1 * time.Millisecond)
		for {
			select {
			case v, ok := <-chan0:
				{
					if !ok {
						if !chan0closed {
							ecs.UnSubscribe(event, chan0)
							chan0 = nil
							chan0closed = true
						} else {
							continue
						}
					} else {
						t.Logf("chan0 value: %+v", v)
					}
				}
			case v, ok := <-chan1:
				{
					if !ok {
						if !chan1closed {
							ecs.UnSubscribe(event, chan1)
							chan1 = nil
							chan1closed = true
						} else {
							continue
						}
					} else {
						t.Logf("chan1 value: %+v", v)

					}
				}
			case <-ticker.C:
				{
					if chan0closed && chan1closed {
						t.Log("cha0 chan1 all close.")
						return
					}
				}
			}
		}
	}()
	t.Logf("EvtChans topics: %s\n", ecs.Topics())
	t.Logf("ecs.Publish:%t.", ecs.Publish(event, "Hello world!"))
	t.Logf("ecs.PublishAsync:%+v.", ecs.PublishAsync(context.Background(), 1*time.Microsecond, event, 1, 2, 3, 4, 5))
	ecs.Close()
	t.Logf("ecs.Publish:%t.", ecs.Publish(event, "Hello world!"))
	t.Logf("ecs.PublishAsync:%+v.", ecs.PublishAsync(context.Background(), 1*time.Microsecond, event, 6, 7, 8, 9, 10))
	t.Logf("EvtChans topics: %s\n", ecs.Topics())
	ecs.WaitAsync()
	time.Sleep(10 * time.Millisecond)
	t.Logf("EvtChans topics: %s\n", ecs.Topics())
	t.Log("OK")
}

func TestManySubscribe(t *testing.T) {
	buflen := 10
	ecs := ec.NewEvtChans(uint(buflen))
	fmt.Printf("EvtChans topics: %s\n", ecs.Topics())
	event := "topic"
	chan0 := ecs.Subscribe(event)
	chan1 := ecs.Subscribe(event)

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	mu.Lock()
	flag := 0
	mu.Unlock()

	recFunc := func(c <-chan any) {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case _, ok := <-c:
				if ok {
					ecs.Topics()
					mu.Lock()
					flag++
					mu.Unlock()
				} else {
					ecs.UnSubscribe(event, c)
					return
				}
			case <-ticker.C:
			}
		}
	}

	wg.Add(1)
	go recFunc(chan0)
	runtime.Gosched()
	wg.Add(1)
	go recFunc(chan1)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)
	chan2 := ecs.Subscribe(event)
	chan3 := ecs.Subscribe(event)
	wg.Add(1)
	go recFunc(chan2)
	wg.Add(1)
	go recFunc(chan3)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)

	fmt.Printf("EvtChans topic chan len: %d\n", ecs.HasChansLen(event))
	fmt.Printf("EvtChans topics: %s\n", ecs.Topics())

	//for test args nil
	sendOk := ecs.Publish(event)
	fmt.Printf("EvtChans topic sendOk: %t\n", sendOk)
	//for test many args
	sendOk = ecs.Publish(event, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	fmt.Printf("EvtChans topic sendOk: %t\n", sendOk)

	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)
	ecs.Close()
	wg.Wait()

	//for test many args
	sendOk = ecs.Publish(event, 1, 2, 3)
	fmt.Printf("EvtChans topic sendOk: %t\n", sendOk)

	//for code cover
	ecs.Close()
	//for test Publish false
	sendOk = ecs.Publish(event, 1, 2, 3)
	fmt.Printf("EvtChans topic sendOk: %t\n", sendOk)

	ecs.UnSubscribe(event, chan0)
	ecs.UnSubscribe(event, chan1)
	ecs.WaitAsync()
	fmt.Printf("EvtChans topics: %s\n", ecs.Topics())
	mu.Lock()
	fmt.Printf("flag: %d\n", flag)
	if flag != int(4*buflen) {
		t.Fail()
	}
	mu.Unlock()
}
