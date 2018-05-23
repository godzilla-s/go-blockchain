package main

import (
	"go-blockchain/p2p/discover"
	"os"
)

func main() {
	cfg := discover.NewConfig(os.Args[1:])
	server := discover.NewP2PServer(cfg)
	if server == nil {
		return
	}

	server.Start()
	server.Stop()
}
