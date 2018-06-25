package tcpv2

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	handshakeRespTime = 5 * time.Second
)

type message struct {
}

type conn struct {
	id      string
	fd      *net.TCPConn
	readCh  chan message
	writeCh chan message
	closed  chan struct{}
}

func newConn(fd *net.TCPConn) *conn {
	c := new(conn)
	c.fd = fd
	c.readCh = make(chan message)
	c.writeCh = make(chan message)
	c.closed = make(chan struct{})

	return c
}

func (c *conn) readLoop() {
	buf := make([]byte, 1028)
	for {
		_, err := c.fd.Read(buf)
		if err != nil && err != io.EOF {
			continue
		}
		if err == io.EOF {
			c.closed <- struct{}{}
			return
		}
		// handle buf
	}
}

func (c *conn) loop() {
	go c.readLoop()
	for {
		select {
		case <-c.closed:
			//
		case <-c.readCh:
			//
		case <-c.writeCh:
		}
	}
}

// 握手协议
func encHandshake(c *conn) {
	c.fd.Write([]byte(""))
	c.fd.SetReadDeadline(time.Now().Add(handshakeRespTime))
	buf := make([]byte, 512)
	c.fd.Read(buf)
}

// 确认握手
func ackHandshake(c *conn) error {
	buf := make([]byte, 512)
	c.fd.SetReadDeadline(time.Now().Add(handshakeRespTime))
	n, err := c.fd.Read(buf)
	if err != nil {
		return err
	}

	// TODO 数据确认
	fmt.Println(buf[:n])
	c.fd.Write([]byte(""))

	return nil
}
