package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zhangzheheng12345/FlowScript/lexer"
	"github.com/zhangzheheng12345/FlowScript/xlexer"
)

/*
ParseValue(lexer.Token) receive a single token (lexer.Token) ,
and translate it into a value which is able to get directly (Value_).
*/
func ParseValue(token lexer.Token) Value {
	if token.Type() == lexer.NUM {
		num, _ := strconv.Atoi(token.Value())
		return Int_{num}
	} else if token.Type() == lexer.SYMBOL {
		return Symbol_{strings.Split(token.Value(), ".")}
	} else if token.Type() == lexer.TMP {
        // TODO: make tmp mark support member visiting chain -.A.B.C...
		return Tmp_{}
	} else if token.Type() == lexer.XEXP {
		return Exp_{xlexer.Lex(strings.Split(token.Value(), ""))}
	} else if token.Type() == lexer.STR {
		return Str_{token.Value()}
	} else if token.Type() == lexer.CHAR {
		return Char_{token.Value()[0]}
	} else {
		fmt.Println("Error: expecting a Value but got kind: ", token.Type())
		return nil
	}
}

/*
Parse([]lexer.Token) receives a token sequence (lexer.Token) ,
and translates it into AST sequence (Ast_).
It returns the processed AST sequence and a int number.
The int number describes where Parse() function stop parsing. (Slice index)
*/
func Parse(tokens []lexer.Token) ([]Ast, int) {
	codes := make([]Ast, 0)
	sendList := make([]Ast, 0)
	index := 0
	checkSend := func() {
		if len(sendList) > 0 {
			sendList = append(sendList, codes[len(codes)-1])
			codes = codes[:len(codes)-1]
			codes = append(codes, Send_{sendList})
			sendList = make([]Ast, 0)
		}
	}
	for ; index < len(tokens); index++ {
		if tokens[index].Type() == lexer.ADD {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Add_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to add. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.SUB {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Sub_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to substrict. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.MULTI {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Multi_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to multiply. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.DIV {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Div_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to divide. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.MOD {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Mod_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to mod. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.BIGR {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Bigr_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to compare. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.SMLR {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Smlr_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to compare. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.EQUAL {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Equal_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to compare. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.AND {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, And_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to and. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.OR {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Or_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to or. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.NOT {
			if index+1 < len(tokens) {
				index++
				codes = append(codes, Not_{ParseValue(tokens[index])})
			} else {
				fmt.Println("Error: lose argumanet to not. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.VAR {
			/*
			   var name => initialize with 0
			   var name v => initialize with v
			*/
			if index+1 < len(tokens) {
				index++
				var name string
				if tokens[index].Type() == lexer.SYMBOL {
					name = tokens[index].Value()
				} else {
					fmt.Println("Error: unallowed variable name. Token: ", index)
					name = ""
				}
				var op Value
				if index+1 < len(tokens) {
					index++
					switch tokens[index].Type() {
					case lexer.NUM, lexer.SYMBOL, lexer.STR, lexer.TMP:
						/* initialize with given value */
						op = ParseValue(tokens[index])
					default:
						/* initialize with 0 */
						op = nil
					}
				} else {
					/* initialize with 0 */
					op = nil
				}
				codes = append(codes, Var_{name, op})
			} else {
				fmt.Println("Error: lost the variable's name while trying to define one")
			}
		} else if tokens[index].Type() == lexer.LEN {
			if index+1 < len(tokens) {
				index++
				op := ParseValue(tokens[index])
				codes = append(codes, Len_{op})
			} else {
				fmt.Println("Error: lose argumanet to not. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.INDEX {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, Index_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to index. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.APP {
			if index+2 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				codes = append(codes, App_{op1, op2})
			} else {
				fmt.Println("Error: lose argumanet to index. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.SLICE {
			if index+3 < len(tokens) {
				index++
				op1 := ParseValue(tokens[index])
				index++
				op2 := ParseValue(tokens[index])
				index++
				op3 := ParseValue(tokens[index])
				codes = append(codes, Slice_{op1, op2, op3})
			} else {
				fmt.Println("Error: lose argumanet to index. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.LIST {
			index++
			ops := make([]Value, 0)
			for index < len(tokens) &&
				(tokens[index].Type() == lexer.NUM ||
					tokens[index].Type() == lexer.STR ||
					tokens[index].Type() == lexer.CHAR ||
					tokens[index].Type() == lexer.TMP ||
					tokens[index].Type() == lexer.SYMBOL ||
					tokens[index].Type() == lexer.XEXP) {
				ops = append(ops, ParseValue(tokens[index]))
				index++
			}
			codes = append(codes, List_{ops})
			index--
		} else if tokens[index].Type() == lexer.IF {
			/*
			   if codition begin ... end
			*/
			if index+3 < len(tokens) {
				index++
				condition := ParseValue(tokens[index])
				index++
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInBlock, i := Parse(tokens[index:])
					index += i
					if index < len(tokens) && tokens[index].Type() == lexer.END {
						codes = append(codes, If_{condition, codesInBlock})
					} else {
						fmt.Println("Error: lost ' end ' at the end of the block. Token ", index)
					}
				} else {
					fmt.Println("Error: lost ' begin ' at the start of the block. Token: ", index)
				}
			} else {
				fmt.Println("Error: not complete if block. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.SEND {
			/*
			   > ... > ... # wrong
			   ... > ... > ... # right
			*/
			if index == 0 || tokens[index-1].Type() == lexer.STOP {
				fmt.Println("Error: use ' > ' at the start of a sentence. There's nothing to send. Token: ", index)
			} else if len(codes) > 0 {
				sendList = append(sendList, codes[len(codes)-1])
				codes = codes[:len(codes)-1]
			}
		} else if tokens[index].Type() == lexer.DEF {
			if index+3 < len(tokens) {
				index++
				var name string
				if tokens[index].Type() == lexer.SYMBOL {
					name = tokens[index].Value()
				} else {
					fmt.Println("Error: unallowed function name. Token: ", index)
					name = ""
				}
				index++
				argList := make([]string, 0)
				for tokens[index].Type() == lexer.SYMBOL && index < len(tokens) {
					argList = append(argList, tokens[index].Value())
					index++
				}
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInFunc, i := Parse(tokens[index:])
					index += i
					if index < len(tokens) && tokens[index].Type() == lexer.END {
						codes = append(codes, Def_{name, argList, codesInFunc})
					} else {
						fmt.Println("Error: lost ' end ' at the end of the block. Token: ", index)
					}
				} else {
					fmt.Println("Error: lost ' begin ' at the start of the block. Token: ", index)
				}
			} else {
				fmt.Println("Error: not complete def block. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.BEGIN {
			if index+1 < len(tokens) {
				index++
				codesInBlock, i := Parse(tokens[index:])
				index += i
				if index < len(tokens) && tokens[index].Type() == lexer.END {
					codes = append(codes, Block_{codesInBlock})
				} else {
					fmt.Println("Error: lost ' end ' at the end of the block. Token: ", index)
				}
			} else {
				fmt.Println("Error: lost ' end ' at the end of the block. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.STRUCT {
			if index+2 < len(tokens) {
				index++
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInStruct, i := Parse(tokens[index:])
					index += i
					if index < len(tokens) && tokens[index].Type() == lexer.END {
						codes = append(codes, Struct_{codesInStruct})
					} else {
						fmt.Println("Error: lost ' end ' at the end of the block. Token: ", index)
					}
				} else {
					fmt.Println("Error: lost ' begin ' at the start of the block. Token: ", index)
				}
			} else {
				fmt.Println("Error: not complete def block. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.ECHO {
			if index+1 < len(tokens) {
				index++
				codes = append(codes, Echo_{ParseValue(tokens[index])})
			} else {
				fmt.Println("Error: lose argument to echo. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.INPUT {
			if index+1 < len(tokens) {
				index++
				codes = append(codes, Input_{ParseValue(tokens[index])})
			} else {
				fmt.Println("Error: lose argument to echo. Token: ", index)
			}
		} else if tokens[index].Type() == lexer.NUM {
			fmt.Println("Error: unexpected constant number. Token: ", index)
		} else if tokens[index].Type() == lexer.TMP {
			fmt.Println("Error: unexpected ' - ' (tmp mark) . Token: ", index)
		} else if tokens[index].Type() == lexer.SYMBOL {
			name := tokens[index].Value()
			index++
			args := make([]Value, 0)
			for ; index < len(tokens) && (tokens[index].Type() == lexer.NUM ||
				tokens[index].Type() == lexer.SYMBOL ||
				tokens[index].Type() == lexer.TMP); index++ {
				args = append(args, ParseValue(tokens[index]))
			}
			index--
			codes = append(codes, Call_{name, args})
		} else if tokens[index].Type() == lexer.STOP {
			checkSend()
		} else if tokens[index].Type() == lexer.END {
			checkSend()
			return codes, index
		}
	}
	checkSend()
	return codes, index
}

/*
RunCode(string) receives a text script ( string ) and run it directly.
The runtime status would not be saved as the script runs in a Block_.
*/
func RunCode(str string) {
	tokens := lexer.Lex(strings.Split(str, ""))
	codes, index := Parse(tokens)
	if index < len(tokens) {
		fmt.Println("Error: unexpected ' end '. Token: ", index)
	} else {
		Block_{codes}.run()
	}
}
