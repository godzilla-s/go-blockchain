package tcpv2

import (
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
