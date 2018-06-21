package main

import (
	"flag"
	"go-blockchain/peer"
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

	id := peer.Rand()

	self := addresses[index-1]
	peers := append(addresses[:index-1], addresses[index:]...)
	c := tcp.NewTCPConn(id.Hex(), self, peers)
	if c == nil {
		panic("fail new tcpconn")
	}

	time.Sleep(10 * time.Second)
	msg := peer.Message{MsgType: peer.PackHeartbeat, ID: "gargfgarvsdfad", Data: []byte("1243251")}
	if index == 1 {
		c.SendMsg(msg)
	}
	time.Sleep(3 * time.Second)
	if index == 2 {
		c.SendMsg(msg)
	}
	time.Sleep(3 * time.Second)
	if index == 3 {
		c.SendMsg(msg)
	}
	for {
	}
}
