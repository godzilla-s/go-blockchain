package tcp

import (
	"bytes"
	"encoding/json"
)

const (
	PackHandshake    = iota + 1 // 握手
	PackAckHandshake            // 握手确认
	PackPing                    // ping
	PackPong                    // pong
	PackBlock                   // 区块数据
)

type Message struct {
	MsgType byte
	ID      string
	Data    []byte
}

func NewMessage() *Message {
	return nil
}

func (m *Message) Encode() []byte {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (m *Message) Decode(data []byte) error {
	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)
	return decoder.Decode(m)
}
