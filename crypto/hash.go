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
	return toHash(buffer.Bytes())
}

// 计算hash
func toHash(data []byte) (Hash, error) {
	d := sha256.New()
	_, err := d.Write(data)
	if err != nil {
		return EmptyHash, err
	}
	h := d.Sum(nil)
	var hash Hash
	copy(hash[:], h[:])
	return hash, nil
}

func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

func (h Hash) Empty() bool {
	if bytes.Equal(h[:], EmptyHash[:]) {
		return true
	}
	return false
}

func testHash() {
	hash, _ := toHash([]byte("1243252465"))
	fmt.Println(hash)

	type Data struct {
		Id   int
		Buf  string
		self []byte
	}
	var d = Data{
		Id:   100,
		Buf:  "hello world",
		self: []byte("123456"),
	}

	h, _ := ToHash(&d)
	fmt.Println(h)
}
