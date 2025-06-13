package token_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	token "common/model/token"
)

func TestNilChan(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	ticker := time.NewTicker(time.Duration(500) * time.Millisecond)
	ci := (chan int)(nil)
	for {
		select {
		//从 nil chan 接收 永久阻塞
		case _, ok := <-ci:
			if !timer.Stop() {
				<-timer.C
			}
			t.Logf(" receive form nil chan: %t", ok)
		case <-ticker.C:
			{
				t.Log("Ticker")
			}

		case <-timer.C:
			t.Logf(" time out: %+v", timer)
			return
		}
	}

}

func TestClosedChan01(t *testing.T) {
	timer := time.NewTimer(time.Second)
	ci := make(chan int, 1)
	ci <- 1
	close(ci)
	for {
		select {
		// 从closed chan 接收 先接收完数据后续 接收返回false 和 零值
		case v, ok := <-ci:
			t.Logf(" receive form nil chan value:%#v ,ok: %t", v, ok)
			time.Sleep(300 * time.Millisecond)
		case <-timer.C:
			t.Logf(" time out: %+v", timer)
			goto exit
		}
	}
exit:
}

func TestClosedChan02(t *testing.T) {
	timer := time.NewTimer(time.Second)
	ci := make(chan int, 1)
	ci <- 1
	close(ci)
	for {
		select {
		// 从closed chan 接收 先接收完数据后续 接收一直返回零值
		case v := <-ci:
			t.Logf(" receive form closed chan value:%+v", v)
			time.Sleep(300 * time.Millisecond)
		case <-timer.C:
			t.Logf(" time out: %+v", timer)
			goto exit
		}
	}
exit:
}

func TestWaitTimeout(t *testing.T) {
	b := token.BaseToken{}
	if b.Wait() {
		t.Fatal("Should have failed")
	}

	if b.WaitTimeout(time.Second) {
		t.Fatal("Should have failed")
	}

	//*********************
	tk := token.NewBaseToken()
	go func(bt *token.BaseToken) {
		bt.SetErr(errors.New("test error"))
	}(tk)

	if tk.WaitTimeout(5 * time.Second) {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.Wait() {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	tk.Reset()
	go func(bt *token.BaseToken) {
		time.Sleep(time.Second)
		bt.SetErr(nil)
	}(tk)

	runtime.Gosched()
	if tk.Wait() {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.WaitTimeout(2 * time.Second) {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	tk.Reset()
	go func(bt *token.BaseToken) {
		time.Sleep(3 * time.Second)
		bt.SetErr(nil)
	}(tk)

	runtime.Gosched()
	if tk.WaitTimeout(2 * time.Second) {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.Wait() {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	tk.Reset()
	go func(bt *token.BaseToken) {
		time.Sleep(1 * time.Second)
		bt.Completed()
	}(tk)
	runtime.Gosched()
	if tk.WaitTimeout(2 * time.Second) {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.Wait() {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	tk.Reset()
	go func(bt *token.BaseToken) {
		time.Sleep(3 * time.Second)
		bt.Completed()
	}(tk)
	runtime.Gosched()
	if tk.WaitTimeout(2 * time.Second) {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.Wait() {
		t.Logf("succeeded , err: %#v", tk.Err())
	} else {
		t.Logf("not succeeded , err: %#v", tk.Err())
	}

}

func BenchmarkEncode(b *testing.B) {

	tk := token.NewBaseToken()
	go func(bt *token.BaseToken) {
		bt.SetErr(errors.New("test error"))
	}(tk)

	if tk.WaitTimeout(5 * time.Second) {
		b.Logf("succeeded , err: %#v", tk.Err())
	} else {
		b.Logf("not succeeded , err: %#v", tk.Err())
	}

	if tk.Wait() {
		b.Logf("succeeded , err: %#v", tk.Err())
	} else {
		b.Logf("not succeeded , err: %#v", tk.Err())
	}

	for i := 0; i < b.N; i++ {
		tk.Reset()
		go func(bt *token.BaseToken) {
			bt.SetErr(fmt.Errorf("%d err", i))
		}(tk)

		if tk.WaitTimeout(5 * time.Second) {
			b.Logf("succeeded , err: %#v", tk.Err())
		} else {
			b.Logf("not succeeded , err: %#v", tk.Err())
		}

		if tk.Wait() {
			b.Logf("succeeded , err: %#v", tk.Err())
		} else {
			b.Logf("not succeeded , err: %#v", tk.Err())
		}

		tk.Reset()
		go func(bt *token.BaseToken) {
			bt.Completed()
		}(tk)

		if tk.WaitTimeout(5 * time.Second) {
			b.Logf("succeeded , err: %#v", tk.Err())
		} else {
			b.Logf("not succeeded , err: %#v", tk.Err())
		}

		if tk.Wait() {
			b.Logf("succeeded , err: %#v", tk.Err())
		} else {
			b.Logf("not succeeded , err: %#v", tk.Err())
		}

	}
}
