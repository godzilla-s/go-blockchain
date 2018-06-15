package dpos

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-blockchain/crypto"
	"log"
	"sync"
	"time"
)

type Block struct {
	Data      []byte
	Timestamp int64
	Sign      []byte
	PrevHash  crypto.Hash
	Number    uint64
}

// NewBlock
func NewBlock(prevHash crypto.Hash, data []byte) *Block {
	return &Block{
		Data:      data,
		Timestamp: time.Now().Unix(),
		PrevHash:  prevHash,
	}
}

// SignBlock 签名
func (b *Block) SignBlock() {
	// TODO
}

// Hash 计算Block hash
func (b *Block) Hash() crypto.Hash {
	buf := b.Encode()
	return crypto.CalcHash(buf)
}

// encode
func (b *Block) Encode() []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(b)
	return buf.Bytes()
}

// decode
func (b *Block) Decode(buf []byte) error {
	reader := bytes.NewReader(buf)
	dec := json.NewDecoder(reader)
	return dec.Decode(b)
}

// 区块链
type BlockChain struct {
	mux    sync.Mutex
	length uint64
	Blocks []Block
}

func (bc *BlockChain) check(b Block) error {
	num := b.Number
	prevBlock, err := bc.getBlockByNumber(uint64(num - 1))
	if err != nil {
		return err
	}

	if prevBlock == nil {
		return nil
	}
	if prevBlock.Number != num-1 {
		return errors.New("invalid blockchain number")
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
func (bc *BlockChain) getLastNumber() uint64 {
	return bc.length
}

func (bc *BlockChain) getBlockByNumber(num uint64) (*Block, error) {
	if num == 0 {
		return nil, nil
	}

	if bc.length < num {
		log.Println("Error: blockChain length:", bc.length, " getNum:", num)
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

	block := &Block{
		Data:      data,
		Timestamp: time.Now().Unix(),
		Number:    currNum + 1,
	}

	prevBlock, err := bc.getBlockByNumber(currNum)
	if err != nil {
		log.Println("get block by number :", err)
		return nil
	}
	if prevBlock == nil {
		block.PrevHash = crypto.EmptyHash
	} else {
		block.PrevHash = prevBlock.Hash()
	}
	return block
}
