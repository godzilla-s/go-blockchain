package pow

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 难度调整值： 越大算的时间越长
const targetBit = 16

type PoW struct {
	block  *Block
	target *big.Int
}

// NewPoW pow
func NewPoW(b *Block) *PoW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBit))

	pow := &PoW{
		b,
		target,
	}

	return pow
}

// 预处理数据
func (pow *PoW) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			big.NewInt(pow.block.Timestamp).Bytes(),
			big.NewInt(targetBit).Bytes(),
			big.NewInt(int64(nonce)).Bytes(),
		},
		[]byte{},
	)

	return data
}

// Calc 计算PoW
func (pow *PoW) Calc() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

// Validate 验证工作量证明
func (pow *PoW) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
