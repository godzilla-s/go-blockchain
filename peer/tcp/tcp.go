package tcp

import (
	"errors"
	"go-blockchain/event"
	"go-blockchain/peer"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	errNotWritable = errors.New("not writable")
	errNotReadable = errors.New("not readable")
)

type TCPConn struct {
	mux        sync.Mutex
	ID         string
	self       *net.TCPAddr
	connpool   map[string]*connection
	exit       chan struct{}
	broadcast  event.Event
	closeEvent event.Event
}

// NewTCPConn 创建一个TCP连接
func NewTCPConn(id, self string, peers []string) *TCPConn {
	var t TCPConn
	laddr, err := net.ResolveTCPAddr("tcp4", self)
	if err != nil {
		return nil
	}

	t.self = laddr
	t.exit = make(chan struct{})
	t.ID = id

	go t.loopAccept()
	go t.loopDail(peers)
	return &t
}

func (t *TCPConn) loopAccept() {
	lsn, err := net.ListenTCP("tcp", t.self)
	if err != nil {
		panic(err)
	}
	log.Println("listen up .....", t.ID)
	for {
		conn, err := lsn.AcceptTCP()
		if err != nil {
			continue
		}
		go t.procConnect(conn)
	}
}

// 处理链接和接收数据
func (t *TCPConn) procConnect(c *net.TCPConn) {
	defer c.Close()
	// 握手确认
	_, err := t.ackHandshake(c)
	if err != nil {
		log.Println("ackHandshake err:", err)
		return
	}
	//conn := connection{conn: c, readable: true}
	buf := make([]byte, 1028)
	//c.CloseRead()
	for {
		n, err := c.Read(buf)
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			// 断开连接
			log.Println("connect closed:", c.RemoteAddr())
			t.closeEvent.Send(t.ID)
			return
		}
		log.Println("read:", buf[:n])
	}
}

func (t *TCPConn) ackHandshake(c *net.TCPConn) (string, error) {
	buf := make([]byte, 256)
	// c.SetReadDeadline()
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}

	var msg peer.Message
	err = msg.Decode(buf[:n])
	if err != nil {
		return "", err
	}

	if msg.MsgType != peer.PackHandshake {
		return "", errors.New("not handshake type")
	}

	id := msg.ID
	data := make([]byte, 128)
	var msg2 peer.Message
	msg2.MsgType = peer.PackAckHandshake
	msg2.ID = t.ID
	msg2.Data = nil
	err = msg2.Decode(data)
	if err != nil {
		return "", err
	}

	_, err = c.Write(data)
	return id, err
}

// 连接其他节点
func (t *TCPConn) connectPeer(laddr *net.TCPAddr) (*connection, error) {
	conn := tryConnect(laddr, 30)
	if conn == nil {
		return nil, errors.New("not connect peer")
	}

	id, err := t.encHandshake(conn)
	if err != nil {
		log.Println("encHandshake err:", err)
		return nil, err
	}

	c := &connection{conn: conn, writable: true}
	c.id = id
	c.message = make(chan peer.Message, 10)
	c.closed = make(chan string)
	c.exit = make(chan struct{})
	t.closeEvent.Subcribe(c.closed)
	t.broadcast.Subcribe(c.message)
	go c.loop()
	return c, nil
}

func tryConnect(laddr *net.TCPAddr, timeout int64) *net.TCPConn {
	for {
		conn, err := net.DialTCP("tcp4", nil, laddr)
		if err == nil {
			return conn
		}
		time.Sleep(2 * time.Second)
		timeout -= 2
		if timeout <= 0 {
			log.Println("connect timeout")
			return nil
		}
	}
}

// 握手确认
func (t *TCPConn) encHandshake(c *net.TCPConn) (string, error) {
	msg := peer.Message{MsgType: peer.PackHandshake, ID: t.ID}
	buf := msg.Encode()
	_, err := c.Write(buf)
	if err != nil {
		return "", err
	}

	data := make([]byte, 256)
	n, err := c.Read(data)
	if err != nil {
		return "", err
	}

	var m peer.Message
	err = m.Decode(data[:n])
	if err != nil {
		return "", err
	}
	if m.MsgType != peer.PackAckHandshake {
		return "", errors.New("not ackhandshake")
	}

	return m.ID, nil
}

// 连接处理
func (t *TCPConn) loopDail(addrs []string) {
	for _, addr := range addrs {
		taddr, err := net.ResolveTCPAddr("tcp4", addr)
		if err != nil {
			return
		}
		_, err = t.connectPeer(taddr)
		if err != nil {
			log.Println("connect -", addr, " fail")
			return
		}
		log.Println("connect -", addr, " ok")
	}
Fin:
	for {
		select {
		case <-t.exit:
			// TODO 通知所有连接线程断开
			t.closeEvent.Send(struct{}{})
			break Fin
		}
	}
}

// SendMsg 发送数据
func (t *TCPConn) SendMsg(msg peer.Message) {
	t.broadcast.Send(msg)
}

func (t *TCPConn) addConnectPool(id string) {
	connect := connection{}
	t.connpool[id] = &connect
}
