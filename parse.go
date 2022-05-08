package parser

import (
	"strconv"
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/lexer"
	"github.com/zhangzheheng12345/flowscript/xlexer"
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
		splited := strings.Split(token.Value(), ".")
		return Symbol_{splited[0], splited[1:]}
	} else if token.Type() == lexer.XEXP {
		return Exp_{xlexer.Lex(strings.Split(token.Value(), ""), token.Line())}
	} else if token.Type() == lexer.STR {
		return Str_{token.Value()}
	} else if token.Type() == lexer.CHAR {
		return Char_{token.Value()[0]}
	} else {
		t := token.Type()
		s := strconv.Itoa(int(t))
		for k, v := range lexer.BuildinCmdMap {
			if v == t {
				s = k
				break
			}
		}
		errlog.Err("parser", token.Line(), "expecting a Value but got kind:", s)
		return Int_{0} // to avoid err
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
	var parseLine int
	checkSend := func() {
		if len(sendList) > 0 {
			sendList = append(sendList, codes[len(codes)-1])
			codes = codes[:len(codes)-1]
			codes = append(codes, Send_{sendList, parseLine})
			sendList = make([]Ast, 0)
		}
	}
	for ; index < len(tokens); index++ {
		parseLine = tokens[index].Line()
		if tokens[index].Type() == lexer.VAR {
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
					errlog.Err("parser", tokens[index].Line(), "unallowed variable name.")
					name = ""
				}
				var op Value
				if index+1 < len(tokens) {
					index++
					switch tokens[index].Type() {
					case lexer.NUM, lexer.SYMBOL, lexer.STR:
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
				codes = append(codes, Var_{name, op, parseLine})
			} else {
				errlog.Err("parser", tokens[index].Line(), "lost the variable's name while trying to define one")
			}
		} else if tokens[index].Type() == lexer.LIST {
			index++
			ops := make([]Value, 0)
			for index < len(tokens) &&
				(tokens[index].Type() == lexer.NUM ||
					tokens[index].Type() == lexer.STR ||
					tokens[index].Type() == lexer.CHAR ||
					tokens[index].Type() == lexer.SYMBOL ||
					tokens[index].Type() == lexer.XEXP) {
				ops = append(ops, ParseValue(tokens[index]))
				index++
			}
			codes = append(codes, List_{ops, parseLine})
			index--
		} else if tokens[index].Type() == lexer.IF {
			/*
			   if codition begin ... end
			*/
			var ifnode If_
			if index+3 < len(tokens) {
				index++
				condition := ParseValue(tokens[index])
				index++
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInBlock, i := Parse(tokens[index:])
					index += i

					if index < len(tokens) {
						if tokens[index].Type() == lexer.END {
							ifnode.condition = condition
							ifnode.ifcodes = codesInBlock
						} else {
							errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
						}
					} else {
						errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
				}
			} else {
				errlog.Err("parser", tokens[index].Line(), "not complete if block.")
			}
			if index+3 < len(tokens) && tokens[index+1].Type() == lexer.ELSE {
				index += 2
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInBlock, i := Parse(tokens[index:])
					index += i
					if index < len(tokens) {
						if tokens[index].Type() == lexer.END {
							ifnode.elsecodes = codesInBlock
						} else {
							errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
						}
					} else {
						errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
				}
			} else {
				ifnode.elsecodes = make([]Ast, 0)
			}
			codes = append(codes, ifnode)
		} else if tokens[index].Type() == lexer.SEND {
			/*
			   > ... > ... # wrong
			   ... > ... > ... # right
			*/
			if index == 0 || tokens[index-1].Type() == lexer.STOP {
				errlog.Err("parser", tokens[index].Line(), "use ' > ' at the start of a sentence. There's nothing to send.")
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
					errlog.Err("parser", tokens[index].Line(), "unallowed function name.")
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
					if index < len(tokens) {
						if tokens[index].Type() == lexer.END {
							codes = append(codes, Def_{name, argList, codesInFunc, parseLine})
						} else {
							errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
						}
					} else {
						errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
				}
			} else {
				errlog.Err("parser", tokens[index].Line(), "not complete def block.")
			}
		} else if tokens[index].Type() == lexer.LAMBDA {
			if index+2 < len(tokens) {
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
					if index < len(tokens) {
						if tokens[index].Type() == lexer.END {
							codes = append(codes, Lambda_{argList, codesInFunc, parseLine})
						} else {
							errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
						}
					} else {
						errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
				}
			} else {
				errlog.Err("parser", tokens[index].Line(), "not complete lambda block.")
			}
		} else if tokens[index].Type() == lexer.BEGIN {
			if index+1 < len(tokens) {
				index++
				codesInBlock, i := Parse(tokens[index:])
				index += i
				if index < len(tokens) {
					if tokens[index].Type() == lexer.END {
						codes = append(codes, Block_{codesInBlock, parseLine})
					} else {
						errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
				}
			} else {
				errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
			}
		} else if tokens[index].Type() == lexer.STRUCT {
			if index+2 < len(tokens) {
				index++
				if tokens[index].Type() == lexer.BEGIN {
					index++
					codesInStruct, i := Parse(tokens[index:])
					index += i
					if index < len(tokens) {
						if tokens[index].Type() == lexer.END {
							codes = append(codes, Struct_{codesInStruct, parseLine})
						} else {
							errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
						}
					} else {
						errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
					}
				} else {
					errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
				}
			} else {
				errlog.Err("parser", tokens[index].Line(), "not complete def block.")
			}
		} else if tokens[index].Type() == lexer.NUM {
			errlog.Err("parser", tokens[index].Line(), "unexpected constant number.")
		} else if tokens[index].Type() == lexer.SYMBOL {
			name := ParseValue(tokens[index])
			index++
			args := make([]Value, 0)
			for ; index < len(tokens) && (tokens[index].Type() == lexer.NUM ||
				tokens[index].Type() == lexer.STR ||
				tokens[index].Type() == lexer.CHAR ||
				tokens[index].Type() == lexer.XEXP ||
				tokens[index].Type() == lexer.SYMBOL); index++ {
				args = append(args, ParseValue(tokens[index]))
			}
			index--
			codes = append(codes, Call_{name, args, parseLine})
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
