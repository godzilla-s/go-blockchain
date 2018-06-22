package peer

import "go-blockchain/peer/tcp"

// https://github.com/y13i/j2y 转化

type Peer struct {
	ID  string
	net NetWorking
}

type NetWorking interface {
	SendMsg(msg tcp.Message)
	AddPeer(addr string)
	DelPeer(id string)
	Close()
}

// NewPeer
func NewPeer(id, addr string) *Peer {
	conn := tcp.NewTCPConn(id, addr)
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
	p.net.AddPeer(addr)
}

// Close 关闭节点
func (p *Peer) Close() {
	p.net.Close()
}
