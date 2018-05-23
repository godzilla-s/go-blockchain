package discover

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
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

func (n *Node) Addr() *net.UDPAddr {
	return &net.UDPAddr{IP: n.IP, Port: n.UDP}
}

func (n *Node) Validate() bool {
	// TODO
	return true
}

func ParseNode(rawUrl string) (*Node, error) {
	if !strings.HasPrefix(rawUrl, "node://") {
		return nil, errors.New("invalid node: prefix")
	}

	data := strings.Split(strings.TrimPrefix(rawUrl, "node://"), "@")
	var node Node

	node.ID = StringID(data[0])

	addr := strings.Split(data[1], ":")

	ip := net.ParseIP(addr[0])
	if ip == nil {
		return nil, errors.New("invalid node: ip")
	}

	node.IP = net.ParseIP(addr[0])

	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, errors.New("invalid node: port")
	}

	if port > 65535 || port < 1025 {
		return nil, errors.New("invalid node: port")
	}
	node.UDP = port
	return &node, nil
}

func MustParseNode(raw string) *Node {
	node, err := ParseNode(raw)
	if err != nil {
		panic(err)
	}
	return node
}

type NodeID [16]byte

func StringID(s string) NodeID {
	b := fmt.Sprintf("%016s", s)
	var id NodeID
	copy(id[:], []byte(b)[:])
	return id
}

// 随机ID
func RandomID() NodeID {
	var id NodeID
	rand.Seed(time.Now().UnixNano())
	rand.Read(id[:])
	return id
}

// bytes
func BytesID(v []byte) NodeID {
	var id NodeID
	if len(v) >= len(id) {
		copy(id[:], v[:len(id)])
	} else {
		copy(id[len(id)-len(v):], v[:])
	}
	return id
}

func (id NodeID) Equal(dst NodeID) bool {
	return bytes.Equal(id[:], dst[:])
}

func (id NodeID) String() string {
	return string(id[:])
}

func (id NodeID) IntV() int64 {
	return 0
}

//++++++++++++++++++++++++++
// DB 存储
type nodeDB struct {
	mux   sync.Mutex
	cache map[NodeID]nodeCache
}

type nodeCache struct {
	node     *Node
	lastPing int64
	fails    int
	bondTime int64
}

func newNodeDB() *nodeDB {
	return &nodeDB{
		cache: make(map[NodeID]nodeCache),
	}
}

func (db *nodeDB) saveNode(n *Node) {
	db.mux.Lock()
	defer db.mux.Unlock()

	now := time.Now().Unix()
	db.cache[n.ID] = nodeCache{node: n, lastPing: now, bondTime: now}
}

func (db *nodeDB) getNode(id NodeID) *nodeCache {
	obj, ok := db.cache[id]
	if ok {
		return &obj
	}
	return nil
}

func (db *nodeDB) findFails(id NodeID) int {
	obj := db.getNode(id)
	if obj == nil {
		return 0
	}
	return obj.fails
}

func (db *nodeDB) updateFails(id NodeID, fails int) {
	db.mux.Lock()
	defer db.mux.Unlock()
	obj := db.getNode(id)
	if obj != nil {
		obj.fails = fails
		db.cache[id] = *obj
	}
}

func (db *nodeDB) deleteNode(id NodeID) {
	db.mux.Lock()
	defer db.mux.Unlock()

	delete(db.cache, id)
}

func (db *nodeDB) lastPing(id NodeID) int64 {
	obj := db.getNode(id)
	if obj == nil {
		return 0
	}
	return obj.lastPing
}

func (db *nodeDB) updateLastping(id NodeID, now int64) {
	db.mux.Lock()
	defer db.mux.Unlock()
	obj := db.getNode(id)
	if obj != nil {
		obj.lastPing = now
		db.cache[id] = *obj
	}
}

func (db *nodeDB) bondTime(id NodeID) int64 {
	v, ok := db.cache[id]
	if ok {
		return v.bondTime
	}
	return 0
}

func (db *nodeDB) updateBondtime(id NodeID, now int64) {
	db.mux.Lock()
	defer db.mux.Unlock()
	obj := db.getNode(id)
	if obj == nil {
		obj.bondTime = now
		db.cache[id] = *obj
	}
}

func (db *nodeDB) query() []*Node {
	var nodes []*Node
	for _, v := range db.cache {
		nodes = append(nodes, v.node)
	}
	return nodes
}
