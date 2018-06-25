package tcp

import (
	"fmt"
	"go-blockchain/event"
	"go-blockchain/peer/putils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type connection struct {
	id         string
	self       string
	rconn      *net.TCPConn      // 网络读通道
	wconn      *net.TCPConn      // 网络写通道
	message    chan Message      // 待发送的消息通道
	messageSub event.Subcription // 订阅消息
	closed     chan struct{}
	lastactive int64 // 最后活跃时间
	trytimes   int   // 尝试连接次数
	wg         *sync.WaitGroup
}

func newConnection(t *TCPConn, id string) *connection {
	c := new(connection)
	c.id = id
	c.message = make(chan Message, 10)
	c.messageSub = t.broadcast.Subcribe(c.message)
	c.closed = make(chan struct{})
	c.lastactive = time.Now().Unix()
	c.wg = &t.wg
	return c
}

func (c *connection) loop() {
Loop:
	for {
		select {
		case <-c.closed: // 断开连接
			break Loop
		case msg := <-c.message: // 广播数据
			log.Println("get message tp broadcast to:", c.id)
			if c.id != c.self {
				err := c.write(msg)
				if err != nil {
					fmt.Println("write err:", err)
				}
			}
		}
	}
	c.wg.Done()
	// 退出
	log.Println("close send connection")
}

// write 发送
func (c *connection) write(msg Message) error {
	//data := make([]byte, 128)
	//r, w := io.Pipe()
	if c.wconn != nil {
		data := msg.Encode()
		_, err := c.wconn.Write(data)
		return err
	}
	return errNotWritable
}

// read 读取数据
func (c *connection) read(data []byte, timeout time.Duration) (int, error) {
	if c.rconn != nil {
		c.rconn.SetReadDeadline(time.Now().Add(timeout)) // 读取超时
		n, err := c.rconn.Read(data)
		if err != nil {
			if putils.ErrContain(err, "i/o timeout") {
				return 0, errTimeout
			}
			return 0, err
		}
		if n == 0 {
			return 0, errNullRead
		}
		return n, nil
	}
	return 0, errNotWritable
}

// close 关闭连接
func (c *connection) close() {
	c.closed <- struct{}{}
	if c.rconn != nil {
		c.rconn.Close()
		c.rconn = nil
	}

	if c.wconn != nil {
		c.wconn.Close()
		c.wconn = nil
	}

	c.messageSub.Unsubcribe()
	close(c.message)
	close(c.closed)
	log.Println("connection close")
}

func isBreakErr(err error) bool {
	return err == errNotReadable || err == errNotWritable || err == io.EOF
}

func isContinueErr(err error) bool {
	return err == errTimeout || err == errNullRead
}
