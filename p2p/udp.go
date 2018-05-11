package p2p

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
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

var (
	expirationTime = 30 * time.Second
)

var (
	errPacketTimeout = errors.New("timeout")
	errPacketHandle  = errors.New("handle packet error")
)

type udp struct {
	conn     conn
	pending  chan *pending
	gotreply chan gotreply
	self     endpoint
	tab      *Table
	Id       NodeID
	exit     chan struct{}
}

//待处理
type pending struct {
	typ      byte
	deadline int64
	callback func(v interface{}) bool
	errch    chan error
}

// 返回处理
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
		Id:       strToNodeID(cfg.Id),
		exit:     make(chan struct{}),
	}

	tab := newTable(&t, cfg)
	t.tab = tab
	return &t
}

func (t *udp) taskLoop() {
	var plist = list.New()
	for {
		//
		select {
		case p := <-t.pending:
			p.deadline = time.Now().Add(1 * time.Minute).Unix() //有效时间
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
	pack, fromID, err := decodePacket(buf)
	if err != nil {
		return err
	}
	err = pack.handle(t, fromID, to)
	return err
}

func (t *udp) sendMessage(typ byte, to *net.UDPAddr, pack packet) {
	data := encodePacket(t.Id, typ, pack)
	t.conn.WriteToUDP(data, to)
}

func (t *udp) addPending(typ byte, call func(v interface{}) bool) <-chan error {
	ch := make(chan error, 1)
	select {
	case t.pending <- &pending{typ: typ, callback: call, errch: ch}:
		// todo
	case <-t.exit:
		ch <- errors.New("udp exit")
	}
	return ch
}

// 处理返回的结果
func (t *udp) handleReply(typ byte, pack packet) bool {
	ch := make(chan bool, 1)
	select {
	case t.gotreply <- gotreply{typ: typ, data: pack, matched: ch}:
		return <-ch
	case <-t.exit:
		return true
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
		FromID string
		Expire int64
	}

	replynode struct {
		Nodes  []*Node
		Expire int64
	}
)

type packet interface {
	handle(t *udp, fromID NodeID, to *net.UDPAddr) error
}

func (p *ping) handle(t *udp, fromID NodeID, to *net.UDPAddr) error {
	if expire(p.Expire) {
		return errPacketTimeout
	}

	t.addPending(pongPacket, func(v interface{}) bool { return true })

	reply := pong{Expire: time.Now().Add(expirationTime).Unix()}

	t.sendMessage(pingPacket, to, &reply)
	return nil
}

func (p *pong) handle(t *udp, fromID NodeID, to *net.UDPAddr) error {
	if expire(p.Expire) {
		return errPacketTimeout
	}

	if !t.handleReply(pongPacket, p) {
		return errPacketHandle
	}
	return nil
}

func (p *findnode) handle(t *udp, fromID NodeID, to *net.UDPAddr) error {
	// todo
	if expire(p.Expire) {
		return errPacketTimeout
	}

	n := NewNode(fromID, to.IP, to.Port, to.Port)
	t.tab.bondNode(n)
	// 返回reply
	reply := replynode{
		Expire: time.Now().Add(expirationTime).Unix(),
	}

	// 取附近的node
	reply.Nodes = t.tab.closest()

	t.sendMessage(replynodePacket, to, &reply)
	return nil
}

func (p *replynode) handle(t *udp, fromID NodeID, to *net.UDPAddr) error {
	if expire(p.Expire) {
		return errPacketTimeout
	}

	if !t.handleReply(replynodePacket, p) {
		return errPacketHandle
	}

	return nil
}

// 编码
func encodePacket(id NodeID, typ byte, pack packet) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(typ)
	// 添加ID
	buf.Write(id[:])
	encoder := json.NewEncoder(buf)
	encoder.Encode(pack)

	return buf.Bytes()
}

// 解码
func decodePacket(buf []byte) (packet, NodeID, error) {
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

	// 获取发送方ID
	var fromID NodeID
	copy(fromID[:], buf[1:17])

	buffer := bytes.NewBuffer(buf[17:])
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(pack)
	if err != nil {
		return nil, fromID, err
	}
	return pack, fromID, nil
}

func expire(ts int64) bool {
	return true
}

func (t *udp) findnode(to *net.UDPAddr) []Node {
	var nodes []Node
	t.addPending(replynodePacket, func(v interface{}) bool {
		// todo nodes
		return true
	})

	p := findnode{
		Expire: time.Now().Add(expirationTime).Unix(),
	}

	t.sendMessage(findnodePacket, to, &p)
	return nodes
}

func (t *udp) ping(to *net.UDPAddr) error {
	t.addPending(pongPacket, func(v interface{}) bool {
		return true
	})

	p := ping{
		Expire: time.Now().Add(expirationTime).Unix(),
	}
	t.sendMessage(pingPacket, to, &p)
	return nil
}

func (t *udp) waitping() error {
	return <-t.addPending(pongPacket, func(v interface{}) bool { return true })
}
