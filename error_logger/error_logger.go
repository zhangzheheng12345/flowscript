package lexlog

import (
	"fmt"
)

func Err(kind string, line int, msg ...interface{}) {
	if len(msg) == 0 {
		fmt.Println("errlog: len(msg) == 0 !!!")
	}
	fmt.Print(kind, ":", line, ": ")
	fmt.Println(msg...)
}

/* The line recorder for runtime */
var Line int = 0
