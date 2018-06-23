// 实现命令行
package console

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterh/liner"
)

const PerfixInput = "> "

type Console struct {
	prompter *liner.State
	prompt   string // 输出前缀
	exit     chan struct{}
}

// New 创建
func New() *Console {
	c := new(Console)
	c.prompter = liner.NewLiner()
	c.prompt = PerfixInput
	c.exit = make(chan struct{})
	c.init()
	return c
}

func (c *Console) init() {
	// TODO
}

// Start 启动命令交互
func (c *Console) Start() {
	var strCh = make(chan string)
	// 从命令行读取数据
	go func() {
		for {
			line, err := c.prompter.Prompt(<-strCh)
			if err != nil {
				fmt.Println("err", err)
				return
			}
			strCh <- line
		}
	}()

	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM) // 捕捉终止信号

	// 处理命令行得到的数据
	for {
		strCh <- c.prompt
		select {
		case <-abort:
			return
		case s := <-strCh:
			if s == "exit" {
				fmt.Println("exit...")
				return
			}
			if len(s) == 0 {
				break
			}
			fmt.Println("read:", s)
		}
	}
}
