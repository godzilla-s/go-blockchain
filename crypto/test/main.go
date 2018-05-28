package main

import (
	"fmt"
	"go-blockchain/crypto"
)

func testSign() {
	prk, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("generate key error:", err)
		return
	}

	fmt.Printf("prk:\n\tD:%d\n\tX:%d\n\tY:%d\n", prk.D, prk.X, prk.Y)
	params := prk.Curve.Params()
	fmt.Printf("params:\n\tP:%d\n\tN:%d\n\tB:%d\n\tGx:%d\n\tGy:%d\n", params.P, params.N, params.B, params.Gx, params.Gy)

	pukid := crypto.PubkeyID(prk)
	puk, err := pukid.PublibKey()
	if err != nil {
		fmt.Println("public key id error:", err)
		return
	}

	fmt.Printf("puk:\n\tX:%d\n\tY:%d\n", puk.X, puk.Y)

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

func main() {
	testSign()
	//testHash()
}
