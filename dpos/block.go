package dpos

import (
	"bytes"
	"errors"
	"log"
	"sync"
	"time"
)

type Block struct {
	Data      []byte
	Timestamp int64
	Sign      []byte
	PrevHash  []byte
	Number    int64
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

func (b *Block) Hash() []byte {
	return []byte{}
}

// 区块链
type BlockChain struct {
	mux    sync.Mutex
	length int
	Blocks []Block
}

func (bc *BlockChain) check(b Block) error {
	num := b.Number
	prevBlock := bc.Blocks[num-1]
	if prevBlock.Number != num-1 {
		return errors.New("invalid blockchain number")
	}

	if bytes.Equal(prevBlock.Hash(), b.PrevHash) {
		return errors.New("invalid blockchain hash")
	}

	return nil
}

func (bc *BlockChain) add(b Block) {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	// 验证
	err := bc.check(b)
	if err != nil {
		log.Println("valid block error:", err)
		return
	}
	bc.Blocks = append(bc.Blocks, b)
}
