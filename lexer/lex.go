package lexer

import (
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	lextools "github.com/zhangzheheng12345/flowscript/lex_tools"
)

/* Core lexer code */
func Lex(str []string) []Token {
	result := make([]Token, 0)
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
			if word == "add" {
				result = append(result, Token{ADD, ""})
			} else if word == "sub" {
				result = append(result, Token{SUB, ""})
			} else if word == "multi" {
				result = append(result, Token{MULTI, ""})
			} else if word == "div" {
				result = append(result, Token{DIV, ""})
			} else if word == "mod" {
				result = append(result, Token{MOD, ""})
			} else if word == "equal" {
				result = append(result, Token{EQUAL, ""})
			} else if word == "bigr" {
				result = append(result, Token{BIGR, ""})
			} else if word == "smlr" {
				result = append(result, Token{SMLR, ""})
			} else if word == "not" {
				result = append(result, Token{NOT, ""})
			} else if word == "and" {
				result = append(result, Token{AND, ""})
			} else if word == "or" {
				result = append(result, Token{OR, ""})
			} else if word == "xor" {
				result = append(result, Token{XOR, ""})
			} else if word == "expr" {
				result = append(result, Token{EXPR, ""})
			} else if word == "var" {
				result = append(result, Token{VAR, ""})
			} else if word == "if" {
				result = append(result, Token{IF, ""})
			} else if word == "else" {
				result = append(result, Token{ELSE, ""})
			} else if word == "def" {
				result = append(result, Token{DEF, ""})
			} else if word == "lambda" {
				result = append(result, Token{LAMBDA, ""})
			} else if word == "struct" {
				result = append(result, Token{STRUCT, ""})
			} else if word == "begin" {
				result = append(result, Token{BEGIN, ""})
			} else if word == "end" {
				result = append(result, Token{END, ""})
			} else if word == "type" {
				result = append(result, Token{TYPE, ""})
			} else if word == "echo" {
				result = append(result, Token{ECHO, ""})
			} else if word == "echoln" {
				result = append(result, Token{ECHOLN, ""})
			} else if word == "input" {
				result = append(result, Token{INPUT, ""})
			} else if word == "list" {
				result = append(result, Token{LIST, ""})
			} else if word == "len" {
				result = append(result, Token{LEN, ""})
			} else if word == "index" {
				result = append(result, Token{INDEX, ""})
			} else if word == "app" {
				result = append(result, Token{APP, ""})
			} else if word == "slice" {
				result = append(result, Token{SLICE, ""})
			} else if word == "words" {
				result = append(result, Token{WORDS, ""})
			} else if word == "lines" {
				result = append(result, Token{LINES, ""})
			} else {
				result = append(result, Token{SYMBOL, word})
			}
		} else if lextools.IsSingleNum(str[index]) {
			var num string
			num, index = lextools.PickNum(str, index)
			result = append(result, Token{NUM, num})
		} else if str[index] == ">" {
			result = append(result, Token{SEND, ""})
		} else if str[index] == "-" {
			if index+1 < len(str) && lextools.IsSingleNum(str[index+1]) {
				var num string
				num, index = lextools.PickNum(str, index)
				result = append(result, Token{NUM, num})
			} else if index+1 < len(str) && str[index+1] == "." {
				index++
				var chain string
				chain, index = lextools.PickSymbol(str, index)
				result = append(result, Token{SYMBOL, "-" + chain})
			} else {
				result = append(result, Token{SYMBOL, "-"})
			}
		} else if str[index] == ";" {
			result = append(result, Token{STOP, ""})
		} else if str[index] == "\n" {
			if len(result) > 0 && result[len(result)-1].Type() != SEND && result[len(result)-1].Type() != BEGIN {
				/*
				   begin \n ... => not end
				   > \n ... => not end
				   or will autoly end
				*/
				result = append(result, Token{STOP, ""})
			}
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
			result = append(result, Token{XEXP, strings.Join(str[begin:index], "")})
		} else if str[index] == ")" {
			errlog.Err("lexer", "too much back parenthesis.")
		} else if str[index] == "\n" {
			// ; => end, obviously
			result = append(result, Token{STOP, ""})
		} else if str[index] == "\"" {
			index++
			// TODO: Avoid connecting string very often ( res is connected very often
			res := ""
			for index < len(str) && str[index] != "\"" {
				if str[index] == "\\" {
					/* Escape character */
					index++
					escChar := lextools.PickEscapeChar(str, index)
					res += escChar
				} else {
					res += str[index]
				}
				index++
			}
			if index >= len(str) && str[index-1] == "\"" {
				errlog.Err("lexer", "lose \" while giving a string value.")
			}
			result = append(result, Token{STR, res})
		} else if str[index] == "'" {
			/* Char */
			index++
			if index+1 < len(str) {
				if len(str[index]) > 1 {
					errlog.Err("lexer", "a char value must be ascii but not other wide codeset.")
				} else {
					res := ""
					if str[index] == "\\" {
						/* Escape character */
						index++
						res = lextools.PickEscapeChar(str, index)
					} else {
						res = str[index]
					}
					index++
					if str[index] == "'" {
						result = append(result, Token{CHAR, res})
					} else {
						errlog.Err("lexer", "lose ' while giving a char value")
					}
				}
			} else {
				errlog.Err("lexer", "lose ' while giving a char value")
			}
		} else {
			errlog.Err("lexer", "unexpected token")
		}
	}
	return result
}
