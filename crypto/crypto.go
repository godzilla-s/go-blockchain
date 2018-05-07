package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

const pukLen = 64

type PublicID [pukLen]byte

// 生成私钥
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curveFunc(), rand.Reader)
}

// 使用P-256的曲线
func curveFunc() elliptic.Curve {
	return elliptic.P256()
}

// 导出公钥ID
func PubkeyID(prk *ecdsa.PrivateKey) PublicID {
	puk := prk.PublicKey
	buf := elliptic.Marshal(curveFunc(), puk.X, puk.Y)
	var id PublicID
	copy(id[:], buf[1:])
	return id
}

// 公钥ID转公钥
func (id PublicID) PublibKey() (*ecdsa.PublicKey, error) {
	if len(id) != pukLen {
		return nil, fmt.Errorf("invalid PublicID length")
	}

	half := pukLen / 2
	puk := &ecdsa.PublicKey{Curve: curveFunc(), X: new(big.Int), Y: new(big.Int)}
	puk.X.SetBytes(id[:half])
	puk.Y.SetBytes(id[half:])
	if !puk.Curve.IsOnCurve(puk.X, puk.Y) {
		return nil, fmt.Errorf("invalid pubkeyID")
	}
	return puk, nil
}

func (id PublicID) String() string {
	return fmt.Sprintf("%x", id[:])
}
