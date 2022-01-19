package tools

import (
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

func PickSymbol(str []string, index int) (string, int) {
    begin := index
	for ; index < len(str) &&
		(IsSingleAlpha(str[index]) ||
		 IsSingleNum(str[index]) ||
		 str[index] == "_"); index++ {
	}
    return strings.Join(str[begin:index], ""), index - 1
}