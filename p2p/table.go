package p2p

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Table struct {
	db        *nodeDB
	self      *Node
	udp       *udp
	nursery   []*Node
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
	tab.explore()
	for {
		select {
		case <-tab.exit:
			return
		case <-time.After(30 * time.Second):
			tab.explore()
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
			fmt.Println("==>", n)
			rn := tab.udp.findnode(n.Addr())
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

	nodes = append(nodes, tab.nursery...)

	return nodes
}

func (tab *Table) bondNode(n *Node) {
	if tab.hasBond(n) {
		return
	}

	//tab.bondslots <- struct{}{}
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
	nodes := tab.db.query()
	return nodes
}

func (tab *Table) bondAll(nodes []*Node) {
	for _, n := range nodes {
		if n.ID == tab.self.ID {
			continue
		}
		//tab.bond(n)
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
		return
	}

	log.Println("wait for ping begin")
	//errc <- err
	errc <- tab.udp.waitping()
	log.Println("wait for ping end")
}

func (tab *Table) add(n *Node) {
	tab.db.saveNode(n)
}
