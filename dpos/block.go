package dpos

import (
	"bytes"
	"encoding/json"
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

func NewBlock(prevHash []byte, data []byte) *Block {
	return &Block{
		Data:      data,
		Timestamp: time.Now().Unix(),
		PrevHash:  prevHash,
	}
}

// SignBlock 签名
func (b *Block) SignBlock() {

}

// hash
func (b *Block) Hash() []byte {
	return []byte{}
}

// encode
func (b *Block) Encode() []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(b)
	return buf.Bytes()
}

func (b *Block) Decode(buf []byte) error {
	reader := bytes.NewReader(buf)
	dec := json.NewDecoder(reader)
	return dec.Decode(b)
}

// 区块链
type BlockChain struct {
	mux    sync.Mutex
	length int64
	Blocks []Block
}

func (bc *BlockChain) check(b Block) error {
	num := b.Number
	prevBlock, err := bc.getBlockByNumber(num - 1)
	if err != nil {
		return err
	}

	if prevBlock == nil {
		return nil
	}
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
	bc.length++
}

// 获取区块链上最后一个number
func (bc *BlockChain) getLastNumber() int64 {
	return bc.length
}

func (bc *BlockChain) getBlockByNumber(num int64) (*Block, error) {
	if num < 0 {
		return nil, errors.New("invalid number")
	}

	if num == 0 {
		return nil, nil
	}

	if bc.length < num {
		return nil, errors.New("blockchain has less blocks")
	}
	return &bc.Blocks[num-1], nil
}

func (bc *BlockChain) pending(b Block) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.Blocks = append(bc.Blocks, b)
	bc.length++
}

func (bc *BlockChain) createBlock(data []byte) *Block {
	currNum := bc.getLastNumber()
	lastBlock, err := bc.getBlockByNumber(currNum)
	if err != nil {
		return nil
	}
	var prevHash []byte
	if lastBlock == nil {
		prevHash = make([]byte, 32)
	} else {
		prevHash = lastBlock.Hash()
	}
	block := &Block{
		Data:      data,
		Timestamp: time.Now().Unix(),
		Number:    currNum + 1,
		PrevHash:  prevHash,
	}
	block.SignBlock()
	return block
}
