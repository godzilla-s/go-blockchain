// 使用btcec

package crypto

import (
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/btcec"
)

func BtcGenerateKey() (*ecdsa.PrivateKey, error) {
	prk, err := btcec.NewPrivateKey(btcec.S256())
	return (*ecdsa.PrivateKey)(prk), err
}

func BtcSign(hash []byte, prk *ecdsa.PrivateKey) ([]byte, error) {
	if prk.Curve != btcec.S256() {
		return nil, errors.New("private curve is not secp256k1")
	}

	sig, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(prk), hash, false)
	if err != nil {
		return nil, err
	}

	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

func BtcSignToPubkey(hash, sig []byte) (*ecdsa.PublicKey, error) {
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig[:])

	pubkey, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, hash)
	return (*ecdsa.PublicKey)(pubkey), err
}
