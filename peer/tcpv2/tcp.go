package tcpv2

import (
	"bytes"
	"errors"
	"fmt"
	"go-blockchain/event"
	"go-blockchain/peer/putils"
	"log"
	"net"
	"sync"
	"time"
)

type NodeID string

type TCPConn struct {
	addr      *net.TCPAddr
	ID        string
	broadcast event.Event
	mux       sync.Mutex
	wg        sync.WaitGroup
	connpool  map[string]*conn
	exit      chan struct{}
}

func New(id, addr string) *TCPConn {
	laddr := putils.ParseTCPAddr(addr)
	t := new(TCPConn)
	t.addr = laddr
	t.exit = make(chan struct{})
	t.connpool = make(map[string]*conn)
	t.ID = id

	go t.loopAccept()
	return t
}

func (t *TCPConn) loopAccept() {
	lsn, err := net.ListenTCP("tcp", t.addr)
	if err != nil {
		return
	}
	log.Println("listen up ...")
	for {
		fd, err := lsn.AcceptTCP()
		if err != nil {
			continue
		}

		//t.wg.Add(1)
		go t.handleConn(fd)
	}
}

// 处理新节点的链接
func (t *TCPConn) handleConn(fd *net.TCPConn) {
	//defer fd.Close()
	conn := t.newConn(fd, false)
	if conn == nil {
		fd.Close()
		return
	}
	t.add(conn)
}

// Connect 连接端点
func (t *TCPConn) connect(addr string, errCh chan<- error) {
	raddr := putils.ParseTCPAddr(addr)
	fd, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		errCh <- err
		return
	}
	conn := t.newConn(fd, true)
	if conn == nil {
		fd.Close()
		errCh <- errors.New("fail to connect")
		return
	}
	t.add(conn)
	errCh <- nil
}

// 创建一个双工连接
func (t *TCPConn) newConn(fd *net.TCPConn, passive bool) *conn {
	c := new(conn)
	c.fd = fd
	var err error
	// 握手协议确认
	if passive {
		err = t.procHandshake(c)
	} else {
		err = t.ackHandshake(c)
	}
	if err != nil {
		return nil
	}

	log.Println("new connect handshake ok:", c.id)

	c.readCh = make(chan message, 1)
	c.writeCh = make(chan message, 1)
	c.closing = make(chan struct{})
	c.lastactive = time.Now().Unix()
	c.writeSub = t.broadcast.Subcribe(c.writeCh)

	//go c.readLoop()
	//go c.loop()
	return c
}

// 添加到连接池
func (t *TCPConn) add(c *conn) {
	t.mux.Lock()
	fmt.Println("id", c.id)
	if conn := t.connpool[c.id]; conn == nil {
		t.connpool[c.id] = c
	} else {
		//TODO
	}
	t.mux.Unlock()
}

// 发送消息
func (t *TCPConn) SendMsg(data string) {
	for _, c := range t.connpool {
		msg := message{msgPack, newMsg(data)}
		c.writeMsg(msg)
	}
}

func (t *TCPConn) ReadMsg(r chan<- string) {
	readFunc := func(c *conn) {
		for {
			msg, err := c.readMsg()
			if err != nil {
				return
			}

			if msg.Type == msgPack {
				v := msg.Data.(*msgSender)
				r <- v.val
			}
		}
	}
	for len(t.connpool) == 0 {
		time.Sleep(500 * time.Millisecond)
	}
	for _, c := range t.connpool {
		go readFunc(c)
	}
}

// 添加节点
func (t *TCPConn) AddPeer(addr string) error {
	errCh := make(chan error, 1)
	go t.connect(addr, errCh)
	return <-errCh
}

// 删除节点
func (t *TCPConn) DelPeer(id string) {

}

func (t *TCPConn) procHandshake(c *conn) error {
	data := encHandshakePack(t.ID)
	_, err := c.fd.Write(data)
	if err != nil {
		return err
	}
	c.fd.SetReadDeadline(time.Now().Add(handshakeRespTime))
	buf := make([]byte, 256)
	n, err := c.fd.Read(buf)
	if err != nil {
		return err
	}
	ver, id := decHandshakePack(buf[:n])
	if ver != protoVer {
		return errors.New("protoc version not match")
	}
	if c.id == "" {
		c.id = id
	} else {
		return errors.New("not match")
	}

	return nil
}

func (t *TCPConn) ackHandshake(c *conn) error {
	buf := make([]byte, 512)
	c.fd.SetReadDeadline(time.Now().Add(handshakeRespTime))
	n, err := c.fd.Read(buf)
	if err != nil {
		return err
	}

	ver, id := decHandshakePack(buf[:n])
	if ver != protoVer {
		return errors.New("protoc version not match")
	}
	c.id = id
	// TODO 数据确认
	data := encHandshakePack(t.ID)
	c.fd.Write(data)

	return nil
}

var protoVer byte = 5 // 小于255的数

// 编码握手包
// 版本号：1位
// 节点ID: 20位
func encHandshakePack(id string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(protoVer)
	buf.WriteString(id)
	return buf.Bytes()
}

// 解码握手包
func decHandshakePack(buf []byte) (ver byte, id string) {
	ver = buf[0]
	id = string(buf[1:])
	return
}
