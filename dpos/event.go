package dpos

import (
	"reflect"
	"sync"
)

type event struct {
	once      sync.Once
	mux       sync.Mutex
	sendLock  chan struct{}
	sendCases []reflect.SelectCase
	inbox     []reflect.SelectCase
	etype     reflect.Type
	closed    bool
}

type caseList []reflect.SelectCase

func (cs caseList) find(ch interface{}) int {
	for i, cas := range cs {
		if cas.Chan.Interface() == ch {
			return i
		}
	}
	return -1
}

func (cs caseList) delete(index int) caseList {
	return append(cs[:index], cs[index+1:]...)
}

func (e *event) Subscribe(channel interface{}) {
	e.once.Do(func() {})

	chanval := reflect.ValueOf(channel)
	chantype := chanval.Type()

	if chantype.Kind() != reflect.Chan || chantype.ChanDir() == reflect.SendDir {
		panic("invalid channel dir type")
	}

	e.mux.Lock()
	defer e.mux.Unlock()

	if !e.typecheck(chantype.Elem()) {
		panic("chan type invalid")
	}

	cas := reflect.SelectCase{Dir: reflect.SelectSend, Chan: chanval}
	e.inbox = append(e.inbox, cas)
}

func (e *event) typecheck(typ reflect.Type) bool {
	if e.etype == nil {
		e.etype = typ
		return true
	}
	return e.etype == typ
}

func (e *event) Send(value interface{}) {

}

func (e *event) Unsubcribe() {

}
