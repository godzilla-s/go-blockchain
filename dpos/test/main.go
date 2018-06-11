package main

import (
	"flag"
	"go-blockchain/dpos"
)

var idx int

func main() {
	flag.IntVar(&idx, "i", 0, "index of node")
	flag.Parse()

	config := dpos.GetConfig("config.yaml")
	node := dpos.NewNode(idx, config)
	node.Start()
}
