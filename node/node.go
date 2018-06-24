package node

import (
	"go-blockchain/console"
	"go-blockchain/peer"
	"go-blockchain/peer/tcp"
)

type Node struct {
	peer    *peer.Peer
	console *console.Console
	exit    chan struct{}
}

func NewNode(cfg *Config) *Node {
	node := new(Node)
	node.peer = peer.NewPeer(cfg.Key, cfg.Addr)
	node.exit = make(chan struct{})
	node.console = console.New()
	return node
}

// 启动节点服务
func (n *Node) Start() {
	go n.console.Start()
	for {
		select {
		case <-n.exit:
			return
		case s := <-n.console.Read():
			msg := tcp.Message{MsgType: tcp.PackHeartbeat, ID: n.peer.ID, Data: []byte(s)}
			n.peer.Send(msg)
		}
	}
}

func (n *Node) Stop() {
	for {
		select {
		case <-n.console.Exit():
			return
		}
	}
}
