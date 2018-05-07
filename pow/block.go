package pow

import "time"

type Block struct {
	Timestamp     int64  // 时间戳
	Data          []byte // 数据
	PrevBlockHash []byte // 上一个区块hash
	Hash          []byte // 当前区块hash
	Nonce         int    // 难度值
}

// 创建一个区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		prevBlockHash,
		[]byte{},
		0,
	}

	pow := NewPoW(block)

	// 计算PoW
	nonce, hash := pow.Calc()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("genesis block", []byte("00000000000000000000000000000000"))
}

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
