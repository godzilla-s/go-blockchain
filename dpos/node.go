package dpos

import (
	"bytes"
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

type nodeInfo struct {
	Index string
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
	ID     string
	self   *net.TCPAddr
	config Config
	pool   *connPool
	exit   chan struct{}
}

// 创建一个Node
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

func (n *Node) Start() {
	log.Println("node start:", n.self.String())
	go n.initConnPool()
	n.startListen()
}

// 初始化连接池
func (n *Node) initConnPool() {
	for _, ns := range n.config.Nodes {
		if ns.ID == n.ID {
			continue
		}

		c, err := net.DialTimeout("tcp", ns.Addr, 30*time.Second)
		if err != nil {
			log.Println("fail to dial err:", err)
			continue
		}

		data := packData(packReqGetID, n.ID, nil)
		c.Write(data)
	}
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

		go n.pool.handleConn(c)
	}
}

// 广播数据
func (n *Node) broadcast(typ byte, data []byte) {
	sendData := packData(typ, n.ID, data)
	for _, c := range n.pool.set {
		c.Write(sendData)
	}
}

//------------------------------
// 消息组成： 类型 + ID + 数据部分
//------------------------------
func packData(typ byte, id string, data []byte) []byte {
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
	set map[string]net.Conn
}

func newConnPool() *connPool {
	return &connPool{
		set: make(map[string]net.Conn),
	}
}

func (cp *connPool) handleConn(c net.Conn) {
	buf := make([]byte, 512)
	for {
		nbyte, err := c.Read(buf)
		if err != nil {
			continue
		}

		cp.handleBuffer(buf[:nbyte])
	}
}

func (cp *connPool) add(Id string, conn net.Conn) {
	cp.mux.Lock()
	defer cp.mux.Unlock()
	cp.set[Id] = conn
}

// 处理接受数据
func (cp *connPool) handleBuffer(buf []byte) {
	switch buf[0] {
	case packReqGetID:
		//
	case packHeartBeat:
		//
	case packBlockData:
	}
}
