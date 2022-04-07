package xlexer

import (
	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	lextools "github.com/zhangzheheng12345/flowscript/lex_tools"
)

func Lex(str []string, startLine int) []Token {
	result := make([]Token, 0)
	for index := 0; index < len(str); index++ {
		if str[index] == " " || str[index] == "\t" {
			// Do nothing
		} else if str[index] == "_" || lextools.IsSingleAlpha(str[index]) {
			var symbol string
			symbol, index = lextools.PickSymbol(str, index)
			result = append(result, Token{SYMBOL, symbol, startLine})
		} else if lextools.IsSingleNum(str[index]) {
			var num string
			num, index = lextools.PickNum(str, index)
			result = append(result, Token{NUM, num, startLine})
		} else if str[index] == "(" {
			result = append(result, Token{FPAREN, "", startLine})
		} else if str[index] == ")" {
			result = append(result, Token{BPAREN, "", startLine})
		} else if str[index] == "+" {
			result = append(result, Token{ADD, "", startLine})
		} else if str[index] == "-" {
			if len(result) > 0 && (result[len(result)-1].Type() == NUM ||
				result[len(result)-1].Type() == SYMBOL ||
				result[len(result)-1].Type() == TMP ||
				result[len(result)-1].Type() == BPAREN) {
				result = append(result, Token{SUB, "", startLine})
			} else if index+1 < len(str) && lextools.IsSingleNum(str[index+1]) {
				var num string
				num, index = lextools.PickNum(str, index)
				result = append(result, Token{NUM, num, startLine})
			} else {
				result = append(result, Token{TMP, "", startLine})
			}
		} else if str[index] == "*" {
			result = append(result, Token{MULTI, "", startLine})
		} else if str[index] == "/" {
			result = append(result, Token{DIV, "", startLine})
		} else if str[index] == "\n" {
			startLine++
		} else {
			errlog.Err("xlexer", startLine, "unexpected character in X expression")
		}
	}
	return result
}
