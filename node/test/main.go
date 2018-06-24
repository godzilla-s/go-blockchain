package main

import "go-blockchain/node"

func main() {
	cfg := node.Config{
		Key:  "ace241235af867a876a87c9e0d149",
		Addr: "192.168.1.100:8002",
	}
	node := node.NewNode(&cfg)
	node.Start()
	node.Stop()
}
