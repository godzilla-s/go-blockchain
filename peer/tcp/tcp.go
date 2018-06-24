package tcp

import (
	"errors"
	"fmt"
	"go-blockchain/event"
	"log"
	"net"
	"sync"
	"time"
)

var (
	errNotWritable = errors.New("not writable")
	errNotReadable = errors.New("not readable")
	errTimeout     = errors.New("time out")
	errNullRead    = errors.New("read empty byte")
)

type cState uint8

const (
	readableConn cState = iota + 1
	writableConn
)

type TCPConn struct {
	mux      sync.Mutex
	ID       string
	self     *net.TCPAddr
	connpool map[string]*connection

	broadcast  event.Event // 广播事件
	closeEvent event.Event // 关闭事件
	wg         sync.WaitGroup
	exit       chan struct{}
	addPeer    chan string // 添加节点
	delPeer    chan string
}

// NewTCPConn 创建一个TCP连接
func NewConn(id, self string) *TCPConn {
	var t TCPConn
	laddr, err := net.ResolveTCPAddr("tcp4", self)
	if err != nil {
		return nil
	}

	fmt.Println("addr:", laddr.String())
	t.self = laddr
	t.exit = make(chan struct{})
	t.ID = id
	t.connpool = make(map[string]*connection)
	t.addPeer = make(chan string, 5)
	t.delPeer = make(chan string, 5)

	go t.loopAccept()
	go t.loopDail()
	return &t
}

func (t *TCPConn) loopAccept() {
	lsn, err := net.ListenTCP("tcp", t.self)
	if err != nil {
		panic(err)
	}
	log.Println("listen up .....", t.ID, lsn.Addr())
	for {
		conn, err := lsn.AcceptTCP()
		if err != nil {
			continue
		}
		t.wg.Add(1)
		go t.procConnect(conn)
	}
}

// 处理链接和接收数据
func (t *TCPConn) procConnect(c *net.TCPConn) {
	defer t.wg.Done()
	// 握手确认
	id, err := t.ackHandshake(c)
	if err != nil {
		log.Println("ackHandshake err:", err)
		c.Close()
		return
	}

	conn := t.addConnPool(id, c, readableConn)
	buf := make([]byte, 1028)
	for {
		n, err := conn.read(buf, 10*time.Second)
		if isBreakErr(err) {
			conn.close()
			return
		}
		if isContinueErr(err) {
			continue
		}
		//var msg peer.Message
		//msg.Decode(buf[:n])
		log.Println("read:", buf[:n])
	}
}

func (t *TCPConn) ackHandshake(c *net.TCPConn) (string, error) {
	buf := make([]byte, 256)
	c.SetReadDeadline(time.Now().Add(5 * time.Second)) // 有效时间5秒
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}

	//log.Println("data=>", buf[:n])
	var msg Message
	err = msg.Decode(buf[:n])
	if err != nil {
		return "", err
	}

	if msg.MsgType != PackHandshake {
		return "", errors.New("not handshake type")
	}

	id := msg.ID
	msg.MsgType = PackAckHandshake
	msg.ID = t.ID
	data := msg.Encode()
	if err != nil {
		return "", err
	}

	_, err = c.Write(data)
	return id, err
}

// 连接其他节点
func (t *TCPConn) connectPeers(addr string) (*connection, error) {
	laddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}
	conn := tryConnect(laddr, 30)
	if conn == nil {
		return nil, errors.New("not connect peer")
	}

	id, err := t.encHandshake(conn)
	if err != nil {
		log.Println("encHandshake err:", err)
		return nil, err
	}

	c := t.addConnPool(id, conn, writableConn)
	t.wg.Add(1)
	go c.loop()
	return c, nil
}

func tryConnect(laddr *net.TCPAddr, timeout int64) *net.TCPConn {
	for {
		conn, err := net.DialTCP("tcp4", nil, laddr)
		if err == nil {
			return conn
		}
		time.Sleep(3 * time.Second)
		timeout -= 3
		if timeout <= 0 {
			log.Println("connect timeout")
			return nil
		}
	}
}

// 握手确认
func (t *TCPConn) encHandshake(c *net.TCPConn) (string, error) {
	// 发送验证包
	msg := Message{MsgType: PackHandshake, ID: t.ID}
	buf := msg.Encode()
	_, err := c.Write(buf)
	if err != nil {
		return "", err
	}

	// 接收并验证确认包
	data := make([]byte, 256)
	c.SetReadDeadline(time.Now().Add(5 * time.Second)) // 5秒确认
	n, err := c.Read(data)
	if err != nil {
		return "", err
	}

	var m Message
	err = m.Decode(data[:n])
	if err != nil {
		return "", err
	}
	if m.MsgType != PackAckHandshake {
		return "", errors.New("not ackhandshake")
	}

	return m.ID, nil
}

// 连接处理
func (t *TCPConn) loopDail() {
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()
	flashConnpool := func() {
		// TODO 连接清理
	}
FinLoop:
	for {
		select {
		case <-t.exit:
			// TODO 通知所有连接线程断开
			t.closeEvent.Send(struct{}{})
			break FinLoop
		case addr := <-t.addPeer:
			_, err := t.connectPeers(addr)
			if err != nil {
				log.Println("connect peer fail:", err)
			}
			log.Println("add connection ", addr, " ok")
		case id := <-t.delPeer:
			conn := t.connpool[id]
			if conn != nil {
				conn.close()
				delete(t.connpool, id)
			}
		case <-timer.C:
			flashConnpool()
		}
	}
}

// 添加连接池
func (t *TCPConn) addConnPool(id string, c *net.TCPConn, typ cState) *connection {
	t.mux.Lock()
	conn := t.connpool[id]
	if conn == nil {
		conn = newConnection(t, id)
		t.connpool[id] = conn
		conn.self = t.ID
	}
	if typ == readableConn {
		if conn.rconn == nil {
			c.CloseWrite()
			conn.rconn = c
		}
	}
	if typ == writableConn {
		if conn.wconn == nil {
			c.CloseRead()
			conn.wconn = c
		}
	}
	t.mux.Unlock()
	return conn
}

// SendMsg 发送数据
func (t *TCPConn) SendMsg(msg Message) {
	t.broadcast.Send(msg)
}

// AddPeer 添加节点
func (t *TCPConn) AddPeer(addr string) {
	select {
	case <-t.exit:
	case t.addPeer <- addr:
	}
}

func (t *TCPConn) DelPeer(id string) {
	select {
	case <-t.exit:
	case t.delPeer <- id:
	}
}

// GetPeersInfo 获取节点信息
func (t *TCPConn) GetPeers() []peerInfo {
	var peers []peerInfo
	for _, conn := range t.connpool {
		var p peerInfo
		p.ID = conn.id
		p.LocalAddr = conn.rconn.LocalAddr().String()
		p.RemoteAddr = conn.rconn.RemoteAddr().String()
		peers = append(peers, p)
	}
	return peers
}

// Close 关闭节点
func (t *TCPConn) Close() {
	// 先断开连接池
	// 再断开主链接
	t.exit <- struct{}{}
	t.wg.Wait()
}

type peerInfo struct {
	ID         string `json:"id"`
	LocalAddr  string `json:"localAddr"`
	RemoteAddr string `json:"remoteAddr"`
}
