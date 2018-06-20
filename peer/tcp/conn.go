package tcp

import (
	"errors"
	"go-blockchain/event"
	"go-blockchain/peer"
	"io"
	"net"
)

var (
	errNotWritable = errors.New("not writable")
	errNotReadable = errors.New("not readable")
)

type TCPConn struct {
	ID        string
	self      *net.TCPAddr
	broadcast *event.Subcription
	exit      chan struct{}
}

func NewConn(conf *peer.Config) *TCPConn {
	var tcp TCPConn

	return &tcp
}

func (t *TCPConn) Start() {
	lsn, err := net.ListenTCP("tcp", t.self)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lsn.AcceptTCP()
		if err != nil {
			continue
		}
		go t.procConnect(conn)
	}

}

// 处理链接
func (t *TCPConn) procConnect(c *net.TCPConn) {
	// 握手确认

	for {

	}
}

func (t *TCPConn) ackhandshake(c *net.TCPConn) {
	buf := make([]byte, 256)
	// c.SetReadDeadline()
	n, err := c.Read(buf)
	if err != nil && err != io.EOF {
		return
	}

	var msg peer.Message
	err = msg.Decode(buf[:n])
	if err != nil {
		return
	}
}

type connection struct {
	read     net.Conn
	readable bool // 是否可读
	write    net.Conn
	writable bool
	addr     *net.TCPAddr
	exit     chan struct{}
}

func (c connection) send(data []byte) error {
	if c.writable {
		_, err := c.write.Write(data)
		if err != nil {
			return err
		}
	}
	return errNotWritable
}

func (c connection) recv(data []byte) error {
	if c.readable {
		n, err := c.read.Read(data)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			// TODO
		}
		data = data[:n]
	}
	return errNotReadable
}
