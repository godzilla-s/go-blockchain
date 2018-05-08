package p2p

import (
	"fmt"
	"net"
)

type Node struct {
	ID  string
	IP  net.IP
	UDP int
	TCP int
}

func NewNode(id string, ip net.IP, udp, tcp int) Node {
	return Node{
		ID:  id,
		IP:  ip,
		UDP: udp,
		TCP: tcp,
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("node://%s@%v:%d", n.ID, n.IP, n.UDP)
}

func ParseNode(s string) *Node {
	return nil
}
