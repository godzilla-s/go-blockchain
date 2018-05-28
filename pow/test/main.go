package main

import (
	"fmt"
	"go-blockchain/pow"
	"strconv"
)

func main() {
	fmt.Println("run ...")

	bc := pow.NewBlockChain()

	bc.AddBlock("send 1 utc to Jamie")
	bc.AddBlock("send 2 more utc to Allen")

	for _, block := range bc.Blocks() {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := pow.NewPoW(block)
		fmt.Printf("PoW: %s, nonce:%d\n", strconv.FormatBool(pow.Validate()), block.Nonce)
		fmt.Println()
	}
}
