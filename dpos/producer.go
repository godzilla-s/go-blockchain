package dpos

import "time"

// 区块产生
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
