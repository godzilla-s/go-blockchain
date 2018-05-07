package pow

import (
	"fmt"
	"go-blockchain/run"
	"strconv"
)

// for test
func init() {
	run.Register("pow", Run)
}

func Run() {
	fmt.Println("run ...")

	bc := NewBlockChain()

	bc.AddBlock("send 1 utc to Jamie")
	bc.AddBlock("send 2 more utc to Allen")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewPoW(block)
		fmt.Printf("PoW: %s, nonce:%d\n", strconv.FormatBool(pow.Validate()), block.Nonce)
		fmt.Println()
	}
}
