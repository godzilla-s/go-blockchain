// 实现命令行
package console

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/peterh/liner"
)

const PerfixInput = "> "

type Console struct {
	prompter *liner.State
	prompt   string // 输出前缀
	read     chan string
	exit     chan struct{}
}

// New 创建
func New() *Console {
	c := new(Console)
	c.prompter = liner.NewLiner()
	c.prompt = PerfixInput
	c.exit = make(chan struct{})
	c.read = make(chan string)
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
				if err == liner.ErrPromptAborted {
					strCh <- ""
					continue
				}
				close(strCh)
				return
			}
			strCh <- line
		}
	}()

	// 处理命令行得到的数据
	go func() {
		abort := make(chan os.Signal, 1)
		signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM) // 捕捉终止信号
		for {
			strCh <- c.prompt
			select {
			case <-abort:
				return
			case s, ok := <-strCh:
				if !ok {
					return
				}
				if isExit(s) {
					//fmt.Println("exit...")
					c.exit <- struct{}{}
					return
				}
				if !validInput(s) {
					break
				}
				//fmt.Println("s:", s)
				c.read <- s
			}
		}
	}()
}

func (c *Console) Read() <-chan string {
	return c.read
}

func (c *Console) Exit() <-chan struct{} {
	return c.exit
}

func isExit(s string) bool {
	return s == "exit" || s == "Exit" || s == "EXIT"
}

func validInput(s string) bool {
	return len(s) > 0
}
