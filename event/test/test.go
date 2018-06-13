package main

import (
	"fmt"
	"go-blockchain/event"
)

func main() {
	ev := new(event.Event) // 创建一个订阅事件
	c1 := make(chan int, 1)
	ev.Subcribe(c1)
	c2 := make(chan int, 1)
	ev.Subcribe(c2)

	go ev.Send(20)
	fmt.Println(<-c1)
	fmt.Println(<-c2)
}
