package tcpv2

import (
	"bytes"
)

type Packet interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	String() string
}

type message struct {
	Type byte
	Data Packet
}

const (
	pingPack = iota + 1
	pongPack
	msgPack
)

func (p *message) msgEncode() ([]byte, error) {
	data, err := p.Data.Encode()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(p.Type)
	buf.Write(data)
	return buf.Bytes(), nil
}

func (p *message) msgDecode(buf []byte) (err error) {
	p.Type = buf[0]
	switch p.Type {
	case pingPack:
	case pongPack:
	case msgPack:
		var m msgSender
		err = m.Decode(buf)
		p.Data = &m
	}
	return
}

type msgSender struct {
	val string
}

func newMsg(data string) *msgSender {
	return &msgSender{
		val: data,
	}
}
func (m *msgSender) Encode() ([]byte, error) {
	return []byte(m.val), nil
}

func (m *msgSender) Decode(data []byte) error {
	return nil
}

func (m *msgSender) String() string {
	return m.val
}
