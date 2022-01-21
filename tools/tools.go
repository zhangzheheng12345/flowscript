package tools

import (
	"fmt"
	"strings"
)

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
	switch v := value.(type) {
	case int:
		return v
	case byte:
		return int(v)
	case []int:
		fmt.Println("Error: int or char type value wanted, but got a int list")
		return 0
	case string:
		fmt.Println("Error: int or char type value wanted, but got a string")
		return 0
	default:
		fmt.Println("Error: int or char type value wanted, but got other unknown type value")
		return 0
	}
}

func WantIntList(value interface{}) []int {
	switch v := value.(type) {
	case int:
		fmt.Println("Error: int list wanted, but got a int value")
		return make([]int, 0)
	case byte:
		fmt.Println("Error: int list wanted, but got a char value")
		return make([]int, 0)
	case []int:
		return v
	case string:
		fmt.Println("Error: int list wanted, but got a string")
		return make([]int, 0)
	default:
		fmt.Println("Error: istringe wanted, but got other unknown type value")
		return make([]int, 0)
	}
}

func WantString(value interface{}) string {
	switch v := value.(type) {
	case int:
		fmt.Println("Error: string wanted, but got a int value")
		return ""
	case byte:
		fmt.Println("Error: string wanted, but got a char value")
		return ""
	case []int:
		fmt.Println("Error: string wanted, but got a char value")
		return ""
	case string:
		return v
	default:
		fmt.Println("Error: string wanted, but got other unknown type value")
		return ""
	}
}
