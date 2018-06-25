package main

import (
	"flag"
	"fmt"
	"go-blockchain/node"
	"go-blockchain/peer/putils"
)

var port int

func main() {
	flag.IntVar(&port, "p", 0, "listen port")
	flag.Parse()

	if port == 0 {
		panic("invalid port")
	}
	localAddr := putils.GetLocalIP().String()
	cfg := node.Config{
		Key:  putils.Rand().Hex(),
		Addr: fmt.Sprintf("%s:%d", localAddr, port),
	}
	node := node.NewNode(&cfg)
	node.Start()
	//node.Wait()
}
