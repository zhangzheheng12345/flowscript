package lextools

import (
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
)

/* Check if a character is a number. */
func IsSingleNum(str string) bool {
	return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 57
}

/* Check if a character is a alpha */
func IsSingleAlpha(str string) bool {
	return len(str) != 1 || ([]byte(str)[0] >= 65 && []byte(str)[0] <= 122)
}

/* Return a string of the symbol's name, if there is a symbol. */
func PickSymbol(str []string, index int) (string, int) {
	begin := index
	for ; index < len(str) &&
		(IsSingleAlpha(str[index]) ||
			IsSingleNum(str[index]) ||
			str[index] == "_" ||
			str[index] == "."); index++ {
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

func PickEscapeChar(str []string, index int, line int) string {
	if index < len(str) {
		if str[index] == "\"" {
			return "\""
		} else if str[index] == "\\" {
			return "\\"
		} else if str[index] == "n" {
			return "\n"
		} else if str[index] == "r" {
			return "\r"
		} else if str[index] == "t" {
			return "\t"
		} else if str[index] == "a" {
			return "\a"
		} else if str[index] == "b" {
			return "\b"
		}
	}
	errlog.Err("lexer", line, "unexpected escape character.")
	return ""
}
