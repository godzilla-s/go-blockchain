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
	defer func() {
		close(p.blockCh)
		close(p.exit)
		timer.Stop()
	}()

	for {
		select {
		case <-p.exit:
			return
		case <-timer.C:
			// 判断顺序轮换生产区块
			block := p.generateBlock()
			if block != nil {
				p.blockCh <- block
			}
			timer.Reset(1 * time.Second)
		}
	}
}

func (p *producer) generateBlock() *Block {
	idx := p.getOrderSlot()
	if p.self != p.order[idx] {
		return nil
	}
	log.Println("produce new block")
	return p.bchain.createBlock(nil)
}

// 获取顺序
func (p *producer) getOrderSlot() int {
	now := time.Now().Unix()
	return int(now) % p.orderNum
}
