package p2p

import (
	"sync"
)

type nodeDB struct {
	mux   sync.Mutex
	cache map[string]*Node
}

func newNodeDB() *nodeDB {
	return &nodeDB{
		cache: make(map[string]*Node),
	}
}

func (db *nodeDB) save(n *Node) {
	db.mux.Lock()
	defer db.mux.Unlock()

	db.cache[n.ID] = n
}

func (db *nodeDB) updateLastTime(id string, now int64) {

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
