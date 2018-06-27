package main

import (
	"flag"
	"fmt"
	"go-blockchain/peer/putils"
	"go-blockchain/peer/tcpv2"
	"time"
)

var index int

var addresses = []string{
	"192.168.1.195:5001",
	"192.168.1.195:5002",
	"192.168.1.195:5003",
}

func main() {
	flag.IntVar(&index, "i", 0, "listen port")
	flag.Parse()
	if index < 0 || index >= len(addresses) {
		panic("invalid index")
	}
	id := putils.Rand().Hex()
	conn := tcpv2.New(id, addresses[index])
	for i := 0; i < len(addresses); i++ {
		if i != index {
			err := conn.AddPeer(addresses[i])
			fmt.Println(err)
		}
	}

	go func() {
		for {
			conn.SendMsg(fmt.Sprintf("hello from %s", id))
			time.Sleep(2 * time.Second)
		}
	}()

	r := make(chan string, 10)
	conn.ReadMsg(r)
	fmt.Println("read begin")
	for v := range r {
		fmt.Println("read:", v)
	}
}
