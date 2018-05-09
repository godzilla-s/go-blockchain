package p2p

import (
	"bytes"
	"container/list"
	"encoding/json"
	"io"
	"net"
	"time"
)

const (
	pingPacket = iota + 1
	pongPacket
	findnodePacket
	replynodePacket
)

type udp struct {
	conn     conn
	pending  chan *pending
	gotreply chan gotreply
}

type pending struct {
	typ      byte
	deadline int64
	callback func(v interface{}) bool
	errch    chan error
}

type gotreply struct {
	typ     byte
	data    interface{}
	matched chan bool
}

type endpoint struct {
	IP       net.IP
	UDP, TCP int
}

type conn interface {
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	WriteToUDP(b []byte, to *net.UDPAddr) (int, error)
	Close() error
	LocalAddr() net.Addr
}

func ListenUDP(cfg Config) {
	laddr, err := net.ResolveUDPAddr("udp", cfg.Laddr)
	if err != nil {
		return
	}

	c, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return
	}

	newUDP(c, cfg)
}

func newUDP(c conn, cfg Config) *udp {
	t := udp{
		conn:     c,
		pending:  make(chan *pending, 10),
		gotreply: make(chan gotreply, 10),
	}

	return &t
}

func (t *udp) taskLoop() {
	var plist = list.New()

	for {
		//
		select {
		case p := <-t.pending:
			//
			plist.PushBack(p)
		case r := <-t.gotreply:
			for pl := plist.Front(); pl != nil; pl = pl.Next() {
				v := pl.Value.(*pending)
				if r.typ == v.typ {
					//
				}
			}
		}
	}
}

func (t *udp) readLoop() {
	buf := make([]byte, 1028)
	for {
		nbytes, from, err := t.conn.ReadFromUDP(buf)
		if err != nil && err != io.EOF {
			return
		}

		if err == io.EOF {
			continue
		}

		if err = t.handleRequest(buf[:nbytes], from); err != nil {
			// 处理失败
		}
	}
}

func (t *udp) handleRequest(buf []byte, to *net.UDPAddr) error {
	pack, err := decodePacket(buf)
	if err != nil {
		return err
	}
	err = pack.handle(t, to)
	return err
}

func (t *udp) send(typ byte, to *net.UDPAddr, pack packet) {
	data := encodePacket(typ, pack)
	t.conn.WriteToUDP(data, to)
}

func (t *udp) addPending(typ byte, call func(v interface{}) bool) {
	select {
	case t.pending <- &pending{typ: typ, callback: call}:
		// todo
	}
}

// 处理返回的结果
func (t *udp) handleReply(pack packet) bool {
	ch := make(chan bool, 1)
	select {
	case t.gotreply <- gotreply{data: pack, matched: ch}:
		return <-ch
	}
}

type (
	ping struct {
		From   endpoint
		To     endpoint
		Expire int64
	}

	pong struct {
		To     endpoint
		Expire int64
	}

	findnode struct {
		Expire int64
	}

	replynode struct {
		Nodes  []*Node
		Expire int64
	}
)

type packet interface {
	handle(t *udp, to *net.UDPAddr) error
}

func (p *ping) handle(t *udp, to *net.UDPAddr) error {
	// todo
	return nil
}

func (p *pong) handle(t *udp, to *net.UDPAddr) error {
	// todo
	return nil
}

func (p *findnode) handle(t *udp, to *net.UDPAddr) error {
	// todo
	return nil
}

func (p *replynode) handle(t *udp, to *net.UDPAddr) error {
	// todo
	return nil
}

// 编码
func encodePacket(typ byte, pack packet) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(typ)

	encoder := json.NewEncoder(buf)
	encoder.Encode(pack)

	return buf.Bytes()
}

// 解码
func decodePacket(buf []byte) (packet, error) {
	var pack packet
	typ := buf[0]
	switch typ {
	case pingPacket:
		pack = new(ping)
	case pongPacket:
		pack = new(pong)
	case findnodePacket:
		pack = new(findnode)
	case replynodePacket:
		pack = new(replynode)
	}

	buffer := bytes.NewBuffer(buf[1:])
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(pack)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (t *udp) findnode(to *net.UDPAddr) []Node {
	var nodes []Node

	t.addPending(replynodePacket, func(v interface{}) bool {
		return true
	})

	p := findnode{
		Expire: time.Now().Unix(),
	}
	t.send(findnodePacket, to, &p)
	return nodes
}
