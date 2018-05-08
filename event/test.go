package event

import (
	"fmt"
	"go-blockchain/run"
)

func init() {
	run.Register("event", Run)
}

func Run() {
	event := new(Event) // 创建一个订阅事件
	c1 := make(chan int, 1)
	event.Subcribe(c1)
	c2 := make(chan int, 1)
	event.Subcribe(c2)

	go event.Send(20)
	fmt.Println(<-c1)
	fmt.Println(<-c2)
}
