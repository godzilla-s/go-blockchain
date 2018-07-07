package tcpv2

import (
	"go-blockchain/event"
	"net"
	"sync"
)

type NodeID string

type TCPConn struct {
	//self      *conn
	addr      *net.TCPAddr
	Id        string
	broadcast event.Event
	wg        sync.WaitGroup
	connpool  map[string]*conn
	exit      chan struct{}
}

func New(id, addr string) *TCPConn {
	laddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}

	tcpc := new(TCPConn)
	tcpc.addr = laddr
	tcpc.exit = make(chan struct{})
	tcpc.connpool = make(map[string]*conn)

	return tcpc
}

func (t *TCPConn) loopAccept() {
	lsn, err := net.ListenTCP("tcp", t.addr)
	if err != nil {
		return
	}
	for {
		con, err := lsn.AcceptTCP()
		if err != nil {
			continue
		}

		t.wg.Add(1)
		go t.handleConn(con)
	}
}

func (t *TCPConn) handleConn(fd *net.TCPConn) {

}

func (t *TCPConn) Dial(addr string) {
	raddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return
	}

	net.DialTCP("tcp", nil, raddr)
}
