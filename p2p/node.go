package p2p

import (
	"fmt"
	"net"
)

type Node struct {
	ID  NodeID
	IP  net.IP
	UDP int
	TCP int
}

func NewNode(id NodeID, ip net.IP, udp, tcp int) *Node {
	return &Node{
		ID:  id,
		IP:  ip,
		UDP: udp,
		TCP: tcp,
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("node://%s@%v:%d", n.ID, n.IP, n.UDP)
}

func (n *Node) ToNetAddr() *net.UDPAddr {
	return &net.UDPAddr{IP: n.IP, Port: n.UDP}
}

func ParseNode(s string) *Node {
	return nil
}

type NodeID [16]byte

func strToNodeID(s string) NodeID {
	b := fmt.Sprintf("%016s", s)
	var id NodeID
	copy(id[:], []byte(b)[:])
	return id
}

func (id NodeID) String() string {
	return string(id[:])
}
