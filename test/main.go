package main

import (
	"fmt"
	"regexp"

	"github.com/robertkrimen/otto"
)

func main() {
	// reg := regexp.MustCompile("(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}:\\d{1,5}")
	// fmt.Println(reg.MatchString("hello"))
	// fmt.Println(reg.MatchString("192.168.1.195:9001"))
	// fmt.Println(reg.FindString("192.168.1.195:9001"))

	reg := regexp.MustCompile("([a-z])+[.]([a-z])+")
	a := "admin.add(\"\")"
	fmt.Println(reg.MatchString(a))
	fmt.Println(reg.FindString(a))

	ottoTest()
}

func ottoTest() {
	vm := otto.New()
	vm.Run(`abc = 2 + 2; console.log("value is " + abc)`)
}
