package peer

import (
	"fmt"
	"go-blockchain/peer/tcp"
	"time"
)

// https://github.com/y13i/j2y 转化

type Peer struct {
	ID  string
	net netWorking
}

type netWorking interface {
	SendMsg(msg tcp.Message)
	AddPeer(addr string)
	DelPeer(id string)
	Close()
}

// NewPeer
func NewPeer(id, addr string) *Peer {
	if id == "" || addr == "" {
		panic("empty id or address")
	}
	conn := tcp.NewConn(id, addr)
	p := new(Peer)
	p.net = conn
	return p
}

// Send 发送消息
func (p *Peer) Send(msg tcp.Message) {
	p.net.SendMsg(msg)
}

// Add 添加节点
func (p *Peer) Add(addr string) {
	fmt.Println(addr)
	p.net.AddPeer(addr)
}

// Close 关闭节点
func (p *Peer) Close() {
	p.net.Close()
}

func (p *Peer) GetPeers() {
}

// ping 心跳包
func (p *Peer) ping() {
	ping := time.NewTicker(5 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-ping.C:
			msg := tcp.Message{MsgType: tcp.PackPing}
			p.Send(msg)
		}
	}
}
