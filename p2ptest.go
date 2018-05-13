package main

import (
	"go-blockchain/p2p"
	"os"
)

func main() {
	cfg := p2p.NewConfig(os.Args[1:])
	server := p2p.NewP2PServer(cfg)
	if server == nil {
		return
	}

	server.Start()
	server.Stop()
}
