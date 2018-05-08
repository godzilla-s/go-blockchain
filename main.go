package main

import (
	"flag"
	"fmt"
	_ "go-blockchain/crypto"
	_ "go-blockchain/event"
	_ "go-blockchain/pow"
	"go-blockchain/run"
)

var fname string

func main() {
	flag.StringVar(&fname, "f", "", "function name")
	flag.Parse()

	functions := run.GetFunctions()

	if f, ok := functions[fname]; ok {
		f()
	} else {
		fmt.Printf("function %s not register\n", fname)
	}
}
