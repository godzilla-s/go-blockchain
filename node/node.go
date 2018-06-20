package node

import "net"

type endpoint struct {
	IP   net.IP
	Port int
}

type Node struct {
	config *Config
	ID     string
	self   *net.TCPAddr
	exit   chan struct{}
}

func NewNode(cfg *Config) *Node {
	var node Node

	return &node
}

func (n *Node) Start() {

}
