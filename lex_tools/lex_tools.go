package lextools

import (
	"strconv"
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
)

/* Check if a character is a number. */
func IsSingleNum(str string) bool {
	return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 57
}

func IsSingleHexNum(str string) bool {
	return len(str) == 1 && (IsSingleNum(str) || ([]byte(str)[0] >= 97 && []byte(str)[0] <= 102))
}

func IsSingleOctNum(str string) bool {
	return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 55
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
	if str[index] == "0" && index+1 < len(str) {
		index++
		if str[index] == "x" {
			index++
			for ; index < len(str) && IsSingleHexNum(str[index]); index++ {
				// Do nothing
			}
		} else if IsSingleOctNum(str[index]) {
			index++
			for ; index < len(str) && IsSingleOctNum(str[index]); index++ {
				// Do nothing
			}
		}
	}
	for ; index < len(str) && IsSingleNum(str[index]); index++ {
		// Do nothing
	}
	return strings.Join(str[begin:index], ""), index - 1
}

func PickEscapeChar(str []string, index int, line int) (string, int) {
	if index < len(str) {
		if str[index] == "\"" {
			return "\"", 0
		} else if str[index] == "\\" {
			return "\\", 0
		} else if str[index] == "n" {
			return "\n", 0
		} else if str[index] == "r" {
			return "\r", 0
		} else if str[index] == "t" {
			return "\t", 0
		} else if str[index] == "a" {
			return "\a", 0
		} else if str[index] == "b" {
			return "\b", 0
		} else if IsSingleNum(str[index]) {
			str, newIndex := PickNum(str, index)
			num, _ := strconv.ParseInt(str, 0, 0)
			return string(num), newIndex - index
		}
	}
	errlog.Err("lexer", line, "unexpected escape character.")
	return "", 0
}
