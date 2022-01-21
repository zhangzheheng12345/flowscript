package tools

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

/* Check if a character is a number. */
func IsSingleNum(str string) bool {
	return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 57
}

/* Check if a character is a alpha */
func IsSingleAlpha(str string) bool {
	return len(str) == 1 && []byte(str)[0] >= 65 && []byte(str)[0] <= 122
}

/* Return a string of the symbol's name, if there is a symbol. */
func PickSymbol(str []string, index int) (string, int) {
	begin := index
	for ; index < len(str) &&
		(IsSingleAlpha(str[index]) ||
			IsSingleNum(str[index]) ||
			str[index] == "_"); index++ {
		// Do nothing
	}
	return strings.Join(str[begin:index], ""), index - 1
}

/* Return a string of the number, if there is a number. */
func PickNum(str []string, index int) (string, int) {
	begin := index
	if str[index] == "-" {
		index++
	}
	for ; index < len(str) && IsSingleNum(str[index]); index++ {
		// Do nothing
	}
	return strings.Join(str[begin:index], ""), index - 1
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
