package p2p

import (
	"sync"
)

type nodeDB struct {
	mux   sync.Mutex
	cache map[string]nodeCache
}

type nodeCache struct {
	node     *Node
	lastPing int64
	fails    int
}

func newNodeDB() *nodeDB {
	return &nodeDB{
		cache: make(map[string]nodeCache),
	}
}

func (db *nodeDB) saveNode(n *Node) {
	db.mux.Lock()
	defer db.mux.Unlock()

	db.cache[n.ID] = nodeCache{node: n}
}

func (db *nodeDB) findFails(id string) int {
	db.mux.Lock()
	defer db.mux.Unlock()

	nc, ok := db.cache[id]
	if !ok {
		return 0
	}
	return nc.fails
}

func (db *nodeDB) updateFails(id string, fails int) {
	db.mux.Lock()
	defer db.mux.Unlock()

	nc, ok := db.cache[id]
	if !ok {
		return
	}
	nc.fails = fails
	db.cache[id] = nc
}

func (db *nodeDB) deleteNode(id string) {
	db.mux.Lock()
	defer db.mux.Unlock()

	delete(db.cache, id)
}

func (db *nodeDB) lastPing(id string) int64 {
	db.mux.Lock()
	defer db.mux.Unlock()

	nc, ok := db.cache[id]
	if !ok {
		return 0
	}
	return nc.lastPing
}

func (db *nodeDB) updateLastPing(id string, now int64) {
	db.mux.Lock()
	defer db.mux.Unlock()

	nc, ok := db.cache[id]
	if !ok {
		return
	}
	nc.lastPing = now
	db.cache[id] = nc
}

type Table struct {
	db   nodeDB
	exit chan struct{}
}

func newTable() {

}

// 探寻
func (tab *Table) explore() {
	for {

	}
}

func (tab *Table) bondNode(n *Node) {

}
