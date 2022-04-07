package lexer

import (
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	lextools "github.com/zhangzheheng12345/flowscript/lex_tools"
)

/* Core lexer code */
func Lex(str []string) []Token {
	result := make([]Token, 0)
	line := 0
	for index := 0; index < len(str); index++ {
		if str[index] == " " || str[index] == "\t" {
			// Do nothing
		} else if str[index] == "#" {
			for ; index < len(str) && str[index] != "\n"; index++ {
				// Do nothing
			}
			index--
		} else if lextools.IsSingleAlpha(str[index]) || str[index] == "_" {
			var word string
			word, index = lextools.PickSymbol(str, index)
			v, ok := buildinCmdMap[word]
			if ok {
				result = append(result, Token{v, "", line})
			} else {
				result = append(result, Token{SYMBOL, word, line})
			}
		} else if lextools.IsSingleNum(str[index]) {
			var num string
			num, index = lextools.PickNum(str, index)
			result = append(result, Token{NUM, num, line})
		} else if str[index] == ">" {
			result = append(result, Token{SEND, "", line})
		} else if str[index] == "-" {
			if index+1 < len(str) && lextools.IsSingleNum(str[index+1]) {
				var num string
				num, index = lextools.PickNum(str, index)
				result = append(result, Token{NUM, num, line})
			} else if index+1 < len(str) && str[index+1] == "." {
				index++
				var chain string
				chain, index = lextools.PickSymbol(str, index)
				result = append(result, Token{SYMBOL, "-" + chain, line})
			} else {
				result = append(result, Token{SYMBOL, "-", line})
			}
		} else if str[index] == ";" {
			// ; => end, obviously
			result = append(result, Token{STOP, "", line})
		} else if str[index] == "\n" {
			if len(result) > 0 && result[len(result)-1].Type() != SEND && result[len(result)-1].Type() != BEGIN {
				/*
				   begin \n ... => not end
				   > \n ... => not end
				   or will autoly end
				*/
				result = append(result, Token{STOP, "", line})
			}
			line++
		} else if str[index] == "(" {
			index++
			begin := index
			count := 1
			for ; index < len(str) && count != 0; index++ {
				if str[index] == "(" {
					count++
				} else if str[index] == ")" {
					count--
				}
			}
			index--
			result = append(result, Token{XEXP, strings.Join(str[begin:index], ""), line})
		} else if str[index] == ")" {
			errlog.Err("lexer", line, "too much back parenthesis.")
		} else if str[index] == "\"" {
			index++
			// TODO: Avoid connecting string very often ( res is connected very often
			res := ""
			for index < len(str) && str[index] != "\"" {
				if str[index] == "\\" {
					/* Escape character */
					index++
					escChar := lextools.PickEscapeChar(str, index, line)
					res += escChar
				} else {
					res += str[index]
				}
				if str[index] == "\n" {
					line++
				}
				index++
			}
			if index >= len(str) && str[index-1] != "\"" {
				errlog.Err("lexer", line, "lose \" while giving a string value.")
			}
			result = append(result, Token{STR, res, line})
		} else if str[index] == "'" {
			/* Char */
			index++
			if index+1 < len(str) {
				if len(str[index]) > 1 {
					errlog.Err("lexer", line, "a char value must be ascii but not other wide codeset.")
				} else {
					res := ""
					if str[index] == "\\" {
						/* Escape character */
						index++
						res = lextools.PickEscapeChar(str, index, line)
					} else {
						res = str[index]
					}
					if str[index] == "\n" {
						line++
					}
					index++
					if str[index] == "'" {
						result = append(result, Token{CHAR, res, line})
					} else {
						errlog.Err("lexer", line, "lose ' while giving a char value")
					}
				}
			} else {
				errlog.Err("lexer", line, "lose ' while giving a char value")
			}
		} else {
			errlog.Err("lexer", line, "unexpected charcter")
		}
	}
	return result
}
