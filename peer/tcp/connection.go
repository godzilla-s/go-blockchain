package tcp

import (
	"go-blockchain/peer"
	"log"
	"net"
)

type connection struct {
	id       string
	conn     *net.TCPConn
	readable bool
	writable bool
	message  chan peer.Message
	exit     chan struct{}
	closed   chan string
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
	log.Println("close send connection")
}

// 发送
func (c *connection) write(msg peer.Message) error {
	//data := make([]byte, 128)
	//r, w := io.Pipe()
	if c.writable {
		data := msg.Encode()
		//fmt.Println("==>", data)
		c.conn.Write(data)
	}
	return errNotWritable
}
