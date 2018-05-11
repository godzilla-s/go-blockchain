package p2p

import (
	"net"
	"sync"
	"time"
)

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

	db.cache[n.ID] = nodeCache{node: n}
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

type Table struct {
	db        *nodeDB
	self      *Node
	udp       *udp
	bondslots chan struct{}
	exit      chan struct{}
}

func newTable(t *udp, cfg Config) *Table {
	tab := &Table{
		db:        newNodeDB(),
		udp:       t,
		exit:      make(chan struct{}),
		bondslots: make(chan struct{}, 5), // 最多处理5个并发事件
	}

	go tab.explore()

	return tab
}

// 探寻
func (tab *Table) explore() {
	for {
		select {
		case <-tab.exit:
			return
		case <-time.After(2 * time.Second):
			//
		}
	}
}

func (tab *Table) bondNode(n *Node) {
	if tab.hasBond(n) {
		return
	}

	tab.bond(n)
}

func (tab *Table) hasBond(n *Node) bool {
	if time.Now().Unix()-tab.db.bondTime(n.ID) > 60*10 {
		return false
	}

	return true
}

// 获取附近的节点
func (tab *Table) closest() []*Node {
	return []*Node{}
}

func (tab *Table) bond(n *Node) {
	<-tab.bondslots
	defer func() {
		tab.bondslots <- struct{}{}
	}()
	errc := make(chan error, 1)
	tab.pingpong(n.ToNetAddr(), errc)
	if err := <-errc; err == nil {
		tab.add(n)
	}
}

// 简单ping节点
func (tab *Table) pingpong(to *net.UDPAddr, errc chan<- error) {
	err := tab.udp.ping(to)
	if err != nil {
		return
	}

	errc <- tab.udp.waitping()
}

func (tab *Table) add(n *Node) {

}
