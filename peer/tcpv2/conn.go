package tcpv2

import (
	"fmt"
	"go-blockchain/event"
	"net"
	"time"
)

const (
	handshakeRespTime = 5 * time.Second
)

type conn struct {
	id         string
	fd         *net.TCPConn
	readCh     chan message // 读取消息
	writeCh    chan message // 写入消息
	closing    chan struct{}
	writeSub   event.Subcription
	lastactive int64 // 最后活跃时间
}

func (c *conn) writeMsg(msg message) {
	data, err := msg.msgEncode()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("write:", data, "id", c.id)
	_, err = c.fd.Write(data)
	fmt.Println("write err:", err)
}

func (c *conn) readMsg() (msg message, err error) {
	buf := make([]byte, 1028)
	n := 0
	n, err = c.fd.Read(buf)
	if err != nil {
		fmt.Println("read err:", err)
		return
	}
	fmt.Println("read:", buf[:n])
	err = msg.msgDecode(buf[:n])
	return
}

func (c *conn) close() {
	c.writeSub.Unsubcribe()
	close(c.readCh)
	close(c.writeCh)
	c.fd.Close()
}
