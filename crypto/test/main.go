package main

import (
	"fmt"
	"go-blockchain/crypto"
)

func main() {
	prk, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("generate key error:", err)
		return
	}

	pukid := crypto.PubkeyID(prk)

	puk, err := pukid.PublibKey()
	if err != nil {
		fmt.Println("public key id error:", err)
		return
	}

	data := []byte("1243254215456")
	sig, err := crypto.Sign(data, prk)
	if err != nil {
		fmt.Println("fail to sign:", err)
		return
	}

	if crypto.VerifySign(puk, data, sig) {
		fmt.Println("verify ok")
	} else {
		fmt.Println("verify fail")
	}
}
