package event

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Event struct {
	once      sync.Once
	mux       sync.Mutex
	removeSub chan interface{} // 移除订阅
	sendLock  chan struct{}
	inbox     caseList     // 被订阅的channel
	sendCase  caseList     // 待发送数据的channel
	etype     reflect.Type // 发送数据类型
}

func (e *Event) init() {
	e.removeSub = make(chan interface{})
	e.sendLock = make(chan struct{}, 1)
	e.sendLock <- struct{}{}
	e.sendCase = caseList{{Chan: reflect.ValueOf(e.removeSub), Dir: reflect.SelectRecv}}
}

// 订阅事件
func (e *Event) Subcribe(ch interface{}) Subcription {
	e.once.Do(e.init) // 初始化
	chVal := reflect.ValueOf(ch)
	chTyp := chVal.Type()

	// channel应有Recieve状态
	if chTyp.Kind() != reflect.Chan || chTyp.ChanDir() == reflect.SendDir {
		panic("invalid channel")
	}

	e.mux.Lock()
	defer e.mux.Unlock()
	// check channel的数据类型
	if !e.typeCheck(chTyp.Elem()) {
		panic(errors.New("type check fail"))
	}

	sub := &feedSub{event: e, ch: chVal, err: make(chan error, 1)}

	cas := reflect.SelectCase{Dir: reflect.SelectSend, Chan: chVal}
	e.inbox = append(e.inbox, cas)
	return sub
}

func (e *Event) typeCheck(rTyp reflect.Type) bool {
	if e.etype == nil {
		e.etype = rTyp
		return true
	}
	return e.etype == rTyp
}

// 发送数据
func (e *Event) Send(val interface{}) (nsent int) {
	e.once.Do(e.init)
	<-e.sendLock

	e.mux.Lock()
	e.sendCase = append(e.sendCase, e.inbox...)
	e.inbox = nil
	e.mux.Unlock()

	rval := reflect.ValueOf(val)
	if !e.typeCheck(rval.Type()) {
		e.sendLock <- struct{}{}
		fmt.Println("send fail: type check")
		return
	}

	// 设置send value
	for i := 1; i < len(e.sendCase); i++ {
		e.sendCase[i].Send = rval
	}

	cases := e.sendCase
	for {
		for i := 1; i < len(cases); i++ {
			// 发送
			if cases[i].Chan.TrySend(rval) {
				// TODO
				nsent++
				cases = cases.deactivate(i)
				i--
			}
		}

		if len(cases) == 1 {
			break
		}

		chosen, recv, _ := reflect.Select(cases)
		if chosen == 0 {
			idx := e.sendCase.find(recv.Interface())
			e.sendCase = e.sendCase.delete(idx)
			if idx >= 0 && idx < len(cases) {
				cases = e.sendCase[:len(cases)-1]
			}
		} else {
			cases = cases.deactivate(chosen)
			nsent++
		}
	}

	// 清除数据
	for i := 1; i < len(e.sendCase); i++ {
		e.sendCase[i].Send = reflect.Value{}
	}
	// 解锁
	e.sendLock <- struct{}{}
	return nsent
}

func (e *Event) remove(sub *feedSub) {
	ch := sub.ch.Interface()
	e.mux.Lock()
	idx := e.inbox.find(ch)
	if idx != -1 {
		e.inbox = e.inbox.delete(idx)
	}
	e.mux.Unlock()

	select {
	case e.removeSub <- ch:
		//
	case <-e.sendLock:
		//
		e.sendCase = e.sendCase.delete(e.sendCase.find(ch))
	}
}

type Subcription interface {
	Unsubcribe()
}

type feedSub struct {
	event   *Event
	ch      reflect.Value
	err     chan error
	errOnce sync.Once
}

func (s *feedSub) Unsubcribe() {
	s.errOnce.Do(func() {
		s.event.remove(s)
		close(s.err)
	})
}

type caseList []reflect.SelectCase

func (cl caseList) find(ch interface{}) int {
	for i, cas := range cl {
		if cas.Chan.Interface() == ch {
			return i
		}
	}
	return -1
}

func (cl caseList) delete(idx int) caseList {
	if idx == -1 {
		return cl
	}
	return append(cl[:idx], cl[idx+1:]...)
}

func (cl caseList) deactivate(idx int) caseList {
	last := len(cl) - 1
	cl[idx], cl[last] = cl[last], cl[idx]
	return cl[:last]
}
