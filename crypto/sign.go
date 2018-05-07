// 签名
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

// 用私钥签名
func Sign(hash []byte, prk *ecdsa.PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, prk, hash)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

// 验证签名
func VerifySign(puk *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	if len(sig) != pukLen {
		return false
	}
	half := len(sig) / 2
	r := new(big.Int).SetBytes(sig[:half])
	s := new(big.Int).SetBytes(sig[half:])
	return ecdsa.Verify(puk, hash, r, s)
}
