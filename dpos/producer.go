package dpos

import (
	"log"
	"time"
)

// 区块产生

type producer struct {
	self     string
	order    []string
	orderNum int
	blockCh  chan *Block
	bchain   *BlockChain
	exit     chan struct{}
}

// 区块生产者
func (n *Node) newProducer() *producer {
	p := &producer{
		self:    n.ID,
		blockCh: make(chan *Block, 30),
		exit:    make(chan struct{}),
	}

	for _, ns := range n.config.Nodes {
		p.order = append(p.order, ns.ID)
	}

	p.orderNum = len(p.order)
	p.bchain = n.blockChain
	// go p.produce()
	return p
}

// 产生区块
func (p *producer) produce() {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	getIdx := func() int {
		return int(time.Now().Unix() % int64(p.orderNum))
	}
	for {
		select {
		case <-p.exit:
			return
		case <-timer.C:
			// 判断顺序轮换生产区块
			idx := getIdx()
			// log.Println("generate index:", idx, " self:", p.self)
			if p.self == p.order[idx] {
				log.Println("genrate block")
				block := p.bchain.createBlock(nil)
				if block != nil {
					p.blockCh <- block
				}
			}
			timer.Reset(1 * time.Second)
		}
	}
}
