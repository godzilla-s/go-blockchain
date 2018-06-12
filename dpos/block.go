package dpos

import "time"

type Block struct {
	Data      []byte
	Timestamp int64
	Sign      []byte
	PrevHash  []byte
}

func NewBlock(data []byte) *Block {
	return &Block{
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// SignBlock 签名
func (b *Block) SignBlock() {

}

func (b *Block) CheckBlock() {
	// 拿到上一个区块， 并对比区块的时间
}

// 区块链
type BlockChain struct {
}

func (bc *BlockChain) check() {

}

func (bc *BlockChain) add() {

}
