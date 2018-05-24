package discover

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Table struct {
	db        *nodeDB
	self      *Node
	udp       *udp
	mux       sync.Mutex
	nursery   []*Node
	stable    []*Node
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

	tab.self = NewNode(t.Id, t.self.IP, t.self.UDP, t.self.TCP)

	if len(cfg.Bootnodes) > 0 {
		tab.nursery = cfg.Bootnodes
	}

	go tab.loop()

	return tab
}

// 探寻
func (tab *Table) loop() {
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	tab.explore()
	for {
		tab.bondAll(tab.nursery)
		select {
		case <-tab.exit:
			return
		case <-timer.C:
			tab.explore()
			timer.Reset(10 * time.Second)
		}
	}
}

// 探寻节点
func (tab *Table) explore() {
	var asked = make(map[NodeID]bool)
	//var replys = make(chan []*Node, 5)
	nodes := tab.getSeedNodes()
	for _, n := range nodes {
		asked[n.ID] = false
	}

	fmt.Println("len nodes:", len(nodes))
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		//log.Println(n.ID, tab.self.ID)
		if n.ID == tab.self.ID {
			continue
		}

		asked[n.ID] = true
		// 查看节点状态
		// fail := tab.db.findFails(n.ID)

		go func() {
			rn := tab.udp.findnode(n.Addr())
			log.Println("rn:==>", rn)
			tab.bondAll(rn)
		}()
	}
}

func (tab *Table) getSeedNodes() []*Node {
	var nodes []*Node

	ns := tab.db.query()
	nodes = append(nodes, ns...)
	if len(nodes) > 5 {
		return nodes
	}

	rest := 5 - len(nodes)
	if len(tab.nursery) < rest {
		nodes = append(nodes, tab.nursery...)
	} else {
		nodes = append(nodes, tab.nursery[:rest]...)
	}

	return nodes
}

func (tab *Table) bondNode(n *Node) {
	if tab.hasBond(n) {
		return
	}

	//tab.bondslots <- struct{}{}
	// tab.db.saveNode(n)
	tab.addNersury(n)
}

// 添加到托管中
func (tab *Table) addNersury(n *Node) {
	tab.mux.Lock()
	defer tab.mux.Unlock()
	tab.nursery = append(tab.nursery, n)
}

func (tab *Table) hasBond(n *Node) bool {
	if time.Now().Unix()-tab.db.bondTime(n.ID) > 60*10 {
		return false
	}

	return true
}

// 获取附近的节点
func (tab *Table) closest() []*Node {
	nodes := tab.db.query()
	return nodes
}

// 对发现的几点进行ping/pong校验
func (tab *Table) bondAll(nodes []*Node) {
	for _, n := range nodes {
		if n.ID == tab.self.ID {
			continue
		}
		tab.bond(n)
	}
}

func (tab *Table) bond(n *Node) {
	// <-tab.bondslots
	// defer func() {
	// 	tab.bondslots <- struct{}{}
	// }()
	errc := make(chan error, 1)
	// check
	lastping := tab.db.lastPing(n.ID)
	if lastping > time.Now().Unix()-3*60 {
		// TODO
		return
	}

	go tab.pingpong(n.Addr(), errc)
	if err := <-errc; err == nil {
		tab.add(n)
	} else {
		log.Println("fail to pingpong", err)
	}
}

// 简单ping节点
func (tab *Table) pingpong(to *net.UDPAddr, errc chan<- error) {
	err := tab.udp.ping(to)
	if err != nil {
		log.Println("ping error:", err)
	}
	errc <- err
}

func (tab *Table) add(n *Node) {
	tab.db.saveNode(n)
}
