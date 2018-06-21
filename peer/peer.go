package peer

// https://github.com/y13i/j2y 转化

type Peer struct {
	ID   string
	conn Connect
}

type Connect interface {
	SendMsg(msg Message)
	Close()
}

func NewPeer() *Peer {
	return nil
}

func (p *Peer) Send(msg Message) {
	p.conn.SendMsg(msg)
}