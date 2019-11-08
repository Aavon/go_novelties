package confunc

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

type Confunc struct {
	// 事件推送队列
	eventCh chan interface{}
	// 退出
	close chan struct{}
	// 最大并发数
	max chan struct{}
	wg  sync.WaitGroup
	// 处理函数
	hdrs     []reflect.Value
	isclosed int32
}

func NewConfunc(maxEvent, maxFunc int, hdrs ...interface{}) (*Confunc, error) {
	fns := []reflect.Value{}
	for _, h := range hdrs {
		fn := reflect.ValueOf(h)
		if fn.Kind() == reflect.Func {
			fns = append(fns, fn)
		} else {
			return nil, errors.New("invalid hdrs.")
		}
	}
	cf := &Confunc{
		eventCh:  make(chan interface{}, maxEvent),
		close:    make(chan struct{}, 1),
		max:      make(chan struct{}, maxFunc),
		hdrs:     fns,
		isclosed: 1,
	}
	return cf, nil
}

func (cf *Confunc) Handle(event interface{}) {
	if atomic.LoadInt32(&cf.isclosed) == 0 {
		cf.eventCh <- event
	}
}

func (cf *Confunc) Start() {
	atomic.StoreInt32(&cf.isclosed, 0)
	go func() {
		for e := range cf.eventCh {
			args := []reflect.Value{reflect.ValueOf(e)}
			for _, h := range cf.hdrs {
				cf.max <- struct{}{}
				cf.wg.Add(1)
				go func() {
					defer func() {
						<-cf.max
						cf.wg.Done()
					}()
					h.Call(args)
				}()
			}
		}
		cf.close <- struct{}{}
	}()
}

func (cf *Confunc) Stop() {
	if atomic.CompareAndSwapInt32(&cf.isclosed, 0, 1) {
		close(cf.eventCh)
		<-cf.close
		cf.wg.Wait()
	}
}

func (cf *Confunc) IsClosed() bool {
	return atomic.LoadInt32(&cf.isclosed) == 1
}
