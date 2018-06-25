package console

import (
	"regexp"
)

var funcReg = regexp.MustCompile("([a-z])+[.]([a-z])+")
