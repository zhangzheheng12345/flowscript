package parser

import (
	"fmt"
	"strings"
)

/* Tool function, for cannot convert bool to int directly */
func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func WantInt(value interface{}) int {
	v, ok := value.(int)
	if !ok {
		v, ok := value.(byte)
		if !ok {
			fmt.Println("Error: int or char value wanted, but got other type value")
		}
		return int(v)
	}
	return v
}

func WantIntList(value interface{}) []int {
	v, ok := value.([]int)
	if !ok {
		fmt.Println("Error: int list wanted, but got other type value")
	}
	return v
}

func WantString(value interface{}) string {
	v, ok := value.(string)
	if !ok {
		fmt.Println("Error: string wanted, but got other type value")
	}
	return v
}

func AbsIndex(length int, index int) (int, bool) {
    if index >= 0 && index < length {
        return index, false
    } else if index < 0 && -index <= length {
        return length + index, false
    } else {
        fmt.Println("Error: index out of range. Index: ", index, " length of the list: ", length)
        return 0, true
    }
}
