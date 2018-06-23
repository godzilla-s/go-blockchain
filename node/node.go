package node

import (
	"go-blockchain/console"
	"go-blockchain/peer"
)

type Node struct {
	peer    *peer.Peer
	console *console.Console
	exit    chan struct{}
}

func NewNode(cfg *Config) *Node {
	node := new(Node)
	node.peer = peer.NewPeer(cfg.Key, cfg.addr)
	node.exit = make(chan struct{})
	node.console = console.New()
	return node
}

// 启动节点服务
func (n *Node) Start() {
}
