package dpos

import (
	"time"
)

// 区块产生

type producer struct {
	self     string
	order    []string
	orderNum int
	blockCh  chan *Block
}

func newProducer(cfg Config) *producer {
	p := &producer{
		blockCh: make(chan *Block, 10),
	}

	for _, ns := range cfg.Nodes {
		p.order = append(p.order, ns.ID)
	}

	return p
}

// 产生区块
func (p *producer) produce() {
	now := time.Now().Unix()
	idx := now % int64(p.orderNum)
	// 判断顺序
	if p.self != p.order[idx] {
		return
	}

	block := NewBlock(nil)

	p.blockCh <- block
}
