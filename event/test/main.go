package main

import (
	"fmt"
	"go-blockchain/event"
	"time"
)

func main() {
	ev := new(event.Event) // 创建一个订阅事件
	c1 := make(chan int, 2)
	ev.Subcribe(c1)
	c2 := make(chan int, 2)
	ev.Subcribe(c2)

	go func() {
		for {
			ev.Send(int(time.Now().Unix()))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for v := range c1 {
			fmt.Println(v)
		}
	}()
	for v := range c2 {
		fmt.Println(v)
	}
}
