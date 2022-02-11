package xlexer

import (
	"fmt"

	"github.com/zhangzheheng12345/flowscript/lex_tools"
)

func Lex(str []string) []Token {
	result := make([]Token, 0)
	for index := 0; index < len(str); index++ {
		if str[index] == " " || str[index] == "\t" {
			// Do nothing
		} else if str[index] == "_" || lextools.IsSingleAlpha(str[index]) {
			var symbol string
			symbol, index = lextools.PickSymbol(str, index)
			result = append(result, Token{SYMBOL, symbol})
		} else if lextools.IsSingleNum(str[index]) {
			var num string
			num, index = lextools.PickNum(str, index)
			result = append(result, Token{NUM, num})
		} else if str[index] == "(" {
			result = append(result, Token{FPAREN, ""})
		} else if str[index] == ")" {
			result = append(result, Token{BPAREN, ""})
		} else if str[index] == "+" {
			result = append(result, Token{ADD, ""})
		} else if str[index] == "-" {
			if len(result) > 0 && (result[len(result)-1].Type() == NUM ||
				result[len(result)-1].Type() == SYMBOL ||
				result[len(result)-1].Type() == TMP ||
				result[len(result)-1].Type() == BPAREN) {
				result = append(result, Token{SUB, ""})
			} else if index+1 < len(str) && lextools.IsSingleNum(str[index+1]) {
				var num string
				num, index = lextools.PickNum(str, index)
				result = append(result, Token{NUM, num})
			} else {
				result = append(result, Token{TMP, ""})
			}
		} else if str[index] == "*" {
			result = append(result, Token{MULTI, ""})
		} else if str[index] == "/" {
			result = append(result, Token{DIV, ""})
		} else {
			fmt.Println("Warn: unknown token in X expression. Token: ", index)
		}
	}
	return result
}
