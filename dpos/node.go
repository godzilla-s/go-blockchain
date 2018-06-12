package dpos

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"path"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	packReqGetID = iota + 1
	packRspGetID
	packHeartBeat
	packBlockData
)

var (
	errConnClosed = errors.New("connect is closed")
)

type nodeInfo struct {
	Index int
	ID    string
	Addr  string
}
type Config struct {
	ProduceBlockSlot    uint64
	ProduceBlocksByTurn uint64
	Nodes               []nodeInfo
}

func GetConfig(filename string) Config {
	_, filestr, _, _ := runtime.Caller(1)
	file := path.Join(path.Dir(filestr), filename)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}
	return config
}

type Node struct {
	ID       string
	self     *net.TCPAddr
	config   Config
	isLeader bool
	order    []string
	pool     *connPool
	broad    event
	exit     chan struct{}
}

// 消息
type message struct {
	MsgTyp byte
	ID     string
	Data   interface{}
}

func (msg *message) encodeMsg() []byte {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(msg)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func decodeMsg(buf []byte) (msg message, err error) {
	reader := bytes.NewReader(buf)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&msg)
	return msg, err
}

// NewNode 创建一个Node
func NewNode(idx int, cfg Config) *Node {
	if idx > len(cfg.Nodes) {
		panic("invalid index: out of range")
	}

	ninfo := cfg.Nodes[idx]
	hostaddr, err := net.ResolveTCPAddr("tcp4", ninfo.Addr)
	if err != nil {
		panic(err)
	}

	node := &Node{
		ID:     ninfo.ID,
		self:   hostaddr,
		config: cfg,
		exit:   make(chan struct{}),
	}
	node.pool = newConnPool()
	return node
}

// Start 启动节点
func (n *Node) Start() {
	log.Println("node start:", n.self.String())
	go n.initConnPool()
	n.startListen()
}

// 启动监听
func (n *Node) startListen() {
	lsn, err := net.ListenTCP("tcp", n.self)
	if err != nil {
		log.Fatal("listen error:", err)
		return
	}

	defer lsn.Close()

	for {
		c, err := lsn.Accept()
		if err != nil {
			continue
		}
		go n.handleAccept(c)
	}
}

// 处理连接
func (n *Node) handleAccept(c net.Conn) {
	buf := make([]byte, 512)
	defer c.Close()

	rid, err := n.handshake(c)
	if err != nil {
		log.Println("handshake fail:", err)
		return
	}

	conn := n.pool.add(rid, 1, c)
	for {
		err := conn.recv(buf)
		if err == errConnClosed {
			// handle close
		} else if err != nil {
			continue
		}
		n.handleMessage(conn, buf)
	}
}

func (n *Node) handleMessage(conn *connection, buf []byte) error {
	msg, err := decodeMsg(buf)
	if err != nil {
		log.Println("decode msg error:", err)
		return err
	}
	switch msg.MsgTyp {
	case packReqGetID:
		msg := message{MsgTyp: packRspGetID, ID: n.ID}
		//conn.send(msg.encodeMsg())
		conn.read.Write(msg.encodeMsg())
	case packBlockData:
		// 区块处理

	}
	return nil
}

// 简单的握手
func (n *Node) handshake(c net.Conn) (string, error) {
	req := message{MsgTyp: packReqGetID, ID: n.ID}
	c.Write(req.encodeMsg())

	buf := make([]byte, 256)
	nbyte, err := c.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	rsp, err := decodeMsg(buf[:nbyte])
	if err != nil {
		return "", nil
	}
	if rsp.MsgTyp != packRspGetID {
		return "", errors.New("invalid response id")
	}
	return rsp.ID, nil
}

// 初始化连接池
func (n *Node) initConnPool() {
	for _, ns := range n.config.Nodes {
		if ns.ID == n.ID {
			continue
		}
		go n.connect(ns)
	}
}

func (n *Node) connect(ninfo nodeInfo) {
	var c net.Conn
	var err error
	for {
		c, err = net.DialTimeout("tcp", ninfo.Addr, 30*time.Second)
		if err != nil {
			log.Println("dial error:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	log.Println("connect to", ninfo.Addr, "ok")
	defer c.Close()

	rid, err := n.handshake(c)
	if err != nil {
		log.Println("handshake fail:", err)
		return
	}
	conn := n.pool.add(rid, 2, c)
	// 订阅事件
	n.broad.Subscribe(conn.broadcast)
	for {
		select {
		case <-n.exit:
			return
		case msg := <-conn.broadcast:
			data := msg.encodeMsg()
			conn.send(data)
		}
	}
}

//------------------------------
// 消息组成： 类型 + ID + 数据部分
//------------------------------
func packingData(typ byte, id string, data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(typ)
	buf.WriteString(id)
	if data != nil {
		buf.Write(data)
	}
	return buf.Bytes()
}

type connPool struct {
	mux sync.Mutex
	set map[string]*connection
}

func newConnPool() *connPool {
	return &connPool{
		set: make(map[string]*connection),
	}
}

func (cp *connPool) add(id string, ctyp int, c net.Conn) *connection {
	cp.mux.Lock()
	defer cp.mux.Unlock()
	var conn *connection
	var ok bool
	conn, ok = cp.set[id]
	if !ok {
		conn = new(connection)
		conn.broadcast = make(chan message)
	}
	if ctyp == 1 {
		conn.read = c
		conn.readable = true
	}
	if ctyp == 2 {
		conn.write = c
		conn.writable = true
	}
	cp.set[id] = conn
	return conn
}

type connection struct {
	read      net.Conn // 读取
	readable  bool
	write     net.Conn // 写入
	writable  bool
	addr      string
	broadcast chan message
}

func (c connection) send(data []byte) error {
	if c.writable {
		_, err := c.write.Write(data)
		return err
	}
	return errors.New("conn is unwritable")
}

func (c connection) recv(data []byte) error {
	if c.readable {
		for {
			n, err := c.read.Read(data)
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				return errConnClosed
			}
			data = data[:n]
			return nil
		}
	}
	return errors.New("conn is unreadale")
}
