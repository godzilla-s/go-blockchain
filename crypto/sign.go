// 签名
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
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

// 根据签名导出公钥
func SignToPubkey(sig, hash []byte) {
	// TODO
}

type Signature struct {
	R *big.Int
	S *big.Int
}

// 解析签名
func recoverSign(sig []byte) *Signature {
	if len(sig) != 64 {
		panic("invalid signature")
	}

	var s Signature
	s.R = new(big.Int).SetBytes(sig[:32])
	s.S = new(big.Int).SetBytes(sig[32:])
	return &s
}

// r: 签名中的r
func recoverKeyFromSign(curve elliptic.Curve, r, s *big.Int) {
	// iter / 2
	rx := new(big.Int).Mul(curve.Params().N, new(big.Int).SetInt64(int64(0)))
	rx.Add(rx, r)
	if rx.Cmp(curve.Params().P) != -1 {
		return
	}
}
