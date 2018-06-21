package tcp

import (
	"go-blockchain/event"
	"go-blockchain/peer"
	"log"
	"net"
)

type connection struct {
	id         string
	conn       *net.TCPConn
	readable   bool
	writable   bool
	message    chan peer.Message
	messageSub event.Subcription
	exit       chan struct{}
	closed     chan string
}

func newConnection(t *TCPConn, id string, conn *net.TCPConn) *connection {
	c := new(connection)
	c.id = id
	c.conn = conn
	c.message = make(chan peer.Message, 10)
	c.messageSub = t.broadcast.Subcribe(c.message)
	c.exit = make(chan struct{})
	c.closed = make(chan string)

	return c
}

func (c *connection) loop() {
Loop:
	for {
		select {
		case <-c.exit:
			break Loop
		case id := <-c.closed: // 断开连接
			log.Println("going to close:", id)
			if id == c.id {
				break Loop
			}
		case msg := <-c.message: // 广播数据
			c.write(msg)
		}
	}
	// 退出
	close(c.message)
	close(c.exit)
	c.conn.Close()
	c.messageSub.Unsubcribe()
	log.Println("close send connection")
}

// 发送
func (c *connection) write(msg peer.Message) error {
	//data := make([]byte, 128)
	//r, w := io.Pipe()
	if c.writable {
		data := msg.Encode()
		c.conn.Write(data)
	}
	return errNotWritable
}
