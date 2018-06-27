package main

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/robertkrimen/otto"
)

func main() {
	chantest()
	//writeBuf()
	//ottoTest()
}

func writeBuf() {
	buf := new(bytes.Buffer)
	buf.WriteByte(255)
	buf.WriteString("abcdefg")
	fmt.Println(buf.Bytes())
}

func regexpTest() {
	// reg := regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}:\\d{1,5}")
	// fmt.Println(reg.MatchString("hello"))
	// fmt.Println(reg.MatchString("192.168.1.195:9001"))
	// fmt.Println(reg.FindString("192.168.1.195:9001"))

	reg := regexp.MustCompile("([a-z])+[.]([a-z])+")
	a := "admin.add(\"\")"
	fmt.Println(reg.MatchString(a))
	fmt.Println(reg.FindString(a))

}
func ottoTest() {
	vm := otto.New()
	vm.Run(`abc = 2 + 2; console.log("value is " + abc)`)
}

func chantest() {
	ch := make(chan int)
	nch := make(chan int)

	go func() {
		for v := range ch {
			select {
			case nch <- v:
			default:
				fmt.Println("default running")
			}
		}
	}()

	for {
		select {
		case <-time.After(2 * time.Second):
			ch <- int(time.Now().Unix())
		case v := <-nch:
			fmt.Println("nch:", v)
		}
	}
}
