package main

import (
	"flag"
	"fmt"
	"go-blockchain/peer"
	"go-blockchain/peer/putils"
	"go-blockchain/peer/tcp"
	"time"
)

var index int
var addresses = []string{
	"192.168.1.195:5001",
	"192.168.1.195:5002",
	"192.168.1.195:5003",
}

func main() {
	flag.IntVar(&index, "i", 0, "index of address")
	flag.Parse()

	if index <= 0 || index > len(addresses) {
		panic("invalid index")
	}

	id := putils.Rand()

	self := addresses[index-1]

	peer := peer.NewPeer(id.Hex(), self)
	for _, addr := range addresses {
		if addr != self {
			peer.Add(addr)
		}
	}

	time.Sleep(5 * time.Second)
	fmt.Println("send test pack 1")
	msg := tcp.Message{MsgType: tcp.PackHeartbeat, ID: "gargfgarvsdfad", Data: []byte("1243251")}
	if index == 1 {
		peer.Send(msg)
	}
	fmt.Println("send test pack 2")
	time.Sleep(3 * time.Second)
	if index == 2 {
		peer.Send(msg)
	}
	fmt.Println("send test pack 3")
	time.Sleep(4 * time.Second)
	if index == 1 {
		fmt.Println("send ....")
		peer.Send(msg)
	}
	for {
	}
}
