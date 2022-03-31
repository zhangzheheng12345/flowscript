package lexlog

import (
	"fmt"
)

func Err(kind string, msg ...interface{}) {
	fmt.Println(kind, "\b:", msg)
}
