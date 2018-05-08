package event

import (
	"fmt"
	"reflect"
	"sync"
)

type Event struct {
	mux      sync.Mutex
	inbox    []reflect.SelectCase // 被订阅的channel
	sendCase []reflect.SelectCase // 待发送数据的channel
	etype    reflect.Type         // 事件的数据类型
}

// 订阅事件
func (e *Event) Subcribe(ch interface{}) {
	chVal := reflect.ValueOf(ch)
	chTyp := chVal.Type()

	// channel应有Recieve状态
	if chTyp.Kind() != reflect.Chan || chTyp.ChanDir() == reflect.SendDir {
		panic("invalid channel")
	}

	e.mux.Lock()
	defer e.mux.Unlock()

	// check channel的数据类型
	if e.etype == nil {
		e.etype = chTyp.Elem()
	} else {
		if e.etype != chTyp.Elem() {
			panic("invalid channel type")
		}
	}

	cas := reflect.SelectCase{Dir: reflect.SelectSend, Chan: chVal}
	e.inbox = append(e.inbox, cas)
}

// 发送数据
func (e *Event) Send(val interface{}) {
	rval := reflect.ValueOf(val)

	e.mux.Lock()
	e.sendCase = append(e.sendCase, e.inbox...)

	if rval.Type() != e.etype {
		fmt.Println("chan type not match")
	}
	e.mux.Unlock()

	// 设置send value
	for i := 1; i < len(e.sendCase); i++ {
		e.sendCase[i].Send = rval
	}

	sendcase := e.sendCase
	for i := 0; i < len(sendcase); i++ {
		if sendcase[i].Chan.TrySend(rval) {
			//
		}
	}
}
