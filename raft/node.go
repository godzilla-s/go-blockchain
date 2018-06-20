package raft

import "net"

type State uint8

const (
	Follower State = iota + 1 // 初始状态
	Candidate
	Leader
)

type endpoint struct {
	IP   net.IP
	Port int
}

type Node struct {
	ep    *endpoint
	state State
}

func newEndpoint(addr string) *endpoint {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}

	return &endpoint{
		IP:   tcpAddr.IP,
		Port: tcpAddr.Port,
	}
}

func NewNode(id, addr string) *Node {
	return &Node{
		ep:    newEndpoint(addr),
		state: Follower,
	}
}

func (n *Node) Start() {

}
