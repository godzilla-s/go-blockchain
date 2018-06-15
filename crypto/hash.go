package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Hash [32]byte

var EmptyHash = Hash{}

func ToHash(v interface{}) (Hash, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	if err != nil {
		return EmptyHash, err
	}
	return CalcHash(buffer.Bytes()), nil
}

// CalcHash 计算hash
func CalcHash(data []byte) (h Hash) {
	d := sha256.New()
	_, err := d.Write(data)
	if err != nil {
		panic(err)
	}
	hash := d.Sum(nil)
	copy(h[:], hash[:])
	return
}

func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

// Empty 为空
func (h Hash) Empty() bool {
	if bytes.Equal(h[:], EmptyHash[:]) {
		return true
	}
	return false
}

// Equal 比较相等
func (h Hash) Equal(c Hash) bool {
	return bytes.Equal(h[:], c[:])
}
