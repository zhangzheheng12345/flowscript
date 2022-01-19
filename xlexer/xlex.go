package xlexer

import (
	"fmt"
	"strings"
    "https://github.com/zhangzheheng12345/FlowScript/tools"
)

func Lex(str []string) []Token {
	result := make([]Token, 0)
	for index := 0; index < len(str); index++ {
		if str[index] == " " || str[index] == "\t" {
			// Do nothing
		} else if str[index] == "_" || tools.IsSingleAlpha(str[index]) {
			begin := index
			for ; index < len(str) &&
				(tools.IsSingleAlpha(str[index]) ||
				 tools.IsSingleNum(str[index]) ||
				 tools.str[index] == "_"); index++ {
			}
			result = append(result, Token{SYMBOL, strings.Join(str[begin:index], "")})
			index--
		} else if tools.IsSingleNum(str[index]) {
			begin := index
			for ; index < len(str) && tools.IsSingleNum(str[index]); index++ {
				// Do nothing
			}
			if len(result) != 0 && result[len(result)-1].Type() == TMP {
				result[len(result)-1] = Token{NUM, "-" + strings.Join(str[begin:index], "")}
			} else {
				result = append(result, Token{NUM, strings.Join(str[begin:index], "")})
			}
			index--
		} else if str[index] == "(" {
			result = append(result, Token{FPAREN, ""})
		} else if str[index] == ")" {
			result = append(result, Token{BPAREN, ""})
		} else if str[index] == "+" {
			result = append(result, Token{ADD, ""})
		} else if str[index] == "-" {
			if result[len(result)-1].Type() == NUM ||
				result[len(result)-1].Type() == SYMBOL ||
				result[len(result)-1].Type() == TMP ||
				result[len(result)-1].Type() == BPAREN {
				result = append(result, Token{SUB, ""})
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
