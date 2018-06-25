package node

import (
	"go-blockchain/console"
	"go-blockchain/peer"
	"go-blockchain/peer/tcp"
	"regexp"
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

var reg = regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}:\\d{1,5}")

// 启动节点服务
func (n *Node) Start() {
	n.console.Start()
	for {
		select {
		case <-n.exit:
			return
		case <-n.console.Exit():
			return
		case s := <-n.console.Read():
			//fmt.Println("read:", s)
			if reg.MatchString(s) {
				v := reg.FindString(s)
				n.peer.Add(v)
				break
			}
			msg := tcp.Message{MsgType: tcp.PackHeartbeat, ID: n.peer.ID, Data: []byte(s)}
			n.peer.Send(msg)
		}
	}
}

func (n *Node) Wait() {
	for {
		select {
		case <-n.console.Exit():
			return
		}
	}
}

func (n *Node) handleString(s string) {
}
