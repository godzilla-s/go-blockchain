package crypto

import (
	"fmt"
	"go-blockchain/run"
)

// for test
func init() {
	run.Register("crypto", Run)
}

func Run() {
	prk, err := GenerateKey()
	if err != nil {
		fmt.Println("generate key error:", err)
		return
	}

	pukid := PubkeyID(prk)

	puk, err := pukid.PublibKey()
	if err != nil {
		fmt.Println("public key id error:", err)
		return
	}

	data := []byte("1243254215456")
	sig, err := Sign(data, prk)
	if err != nil {
		fmt.Println("fail to sign:", err)
		return
	}

	if VerifySign(puk, data, sig) {
		fmt.Println("verify ok")
	} else {
		fmt.Println("verify fail")
	}
}
