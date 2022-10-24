package parser

import (
	"strconv"
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/lexer"
	"github.com/zhangzheheng12345/flowscript/xlexer"
)

func isValue(token lexer.Token) bool {
	return token.Type() == lexer.NUM ||
		token.Type() == lexer.STR ||
		token.Type() == lexer.CHAR ||
		token.Type() == lexer.XEXP ||
		token.Type() == lexer.SYMBOL
}

/*
ParseValue(lexer.Token) receive a single token,
and translate it into a value which is able to get directly.
*/
func ParseValue(token lexer.Token) Value {
	if token.Type() == lexer.NUM {
		num, _ := strconv.ParseInt(token.Value(), 0, 0)
		return Int_{int(num)}
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

func parseBlock(tokens []lexer.Token, index int) ([]Ast, int) {
	/* begin ... end */
	if tokens[index].Type() != lexer.BEGIN {
		errlog.Err("parser", tokens[index].Line(), "lost ' begin ' at the start of the block.")
		return nil, index
	}
	index++
	codesInCode, i := Parse(tokens[index:])
	index += i
	if index >= len(tokens) {
		errlog.Err("parser", tokens[len(tokens)-1].Line(), "lost 'end' at the end of the block.")
		return nil, index
	}
	if tokens[index].Type() != lexer.END {
		errlog.Err("parser", tokens[index].Line(), "lost 'end' at the end of the block.")
		return nil, index
	}
	return codesInCode, index
}

func parseVar(tokens []lexer.Token, index int) (Ast, int) {
	/*
	   var name => initialize with 0
	   var name v => initialize with v
	*/
	if index+1 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "lost the variable's name while trying to define one")
		return nil, index // nil: failed
	}
	index++
	var name string
	if tokens[index].Type() == lexer.SYMBOL {
		name = tokens[index].Value()
	} else {
		errlog.Err("parser", tokens[index].Line(), "unallowed variable name.")
		return nil, index
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
	return Var_{name, op, tokens[index].Line()}, index
}

func parseFill(tokens []lexer.Token, index int) (Ast, int) {
	if index+1 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "lost the function while using fill")
		return nil, index
	}
	index++
	if index+1 < len(tokens) && isValue(tokens[index+1]) {
		/* A not complete function call => curried */
		var fn Ast
		fn, index = parseCall(tokens, index)
		return Fill_{fn, nil, tokens[index].Line()}, index
	} else {
		// avoid that cannot fill echo, echoln, etc. as they have no arg num limit
		// which makes them never return a curried func even if you give no args
		return Fill_{
			nil, ParseValue(tokens[index]), tokens[index].Line(),
		}, index
	}
}

func parseEnum(tokens []lexer.Token, index int) (Ast, int) {
	if index+1 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "cannot use enum without delcaring any enum value")
		return nil, index
	}
	index++
	var names []string
	for ; index < len(tokens) && tokens[index].Type() == lexer.SYMBOL; index++ {
		names = append(names, tokens[index].Value())
	}
	return Enum_{names, tokens[index].Line()}, index
}

func parseIf(tokens []lexer.Token, index int) (Ast, int) {
	/* if codition begin ... end */
	if index+3 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "not complete if block.")
		return nil, index
	}
	var ifnode If_
	index++
	condition := ParseValue(tokens[index])
	index++
	block, index := parseBlock(tokens, index)
	ifnode = If_{condition, block, nil, tokens[index].Line()}
	if index+3 < len(tokens) && tokens[index+1].Type() == lexer.ELSE {
		index += 2
		ifnode.elsecodes, index = parseBlock(tokens, index)
	} else {
		ifnode.elsecodes = make([]Ast, 0)
	}
	return ifnode, index
}

func parseDef(tokens []lexer.Token, index int) (Ast, int) {
	/* def name arg1 arg2 ... begin ... end */
	if index+3 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "not complete def block.")
		return nil, index
	}
	index++
	var name string
	if tokens[index].Type() == lexer.SYMBOL {
		name = tokens[index].Value()
	} else {
		errlog.Err("parser", tokens[index].Line(), "unallowed function name.")
		return nil, index
	}
	index++
	argList := make([]string, 0)
	for tokens[index].Type() == lexer.SYMBOL && index < len(tokens) {
		argList = append(argList, tokens[index].Value())
		index++
	}
	block, index := parseBlock(tokens, index)
	return Def_{name, argList, block, tokens[index].Line()}, index
}

func parseLambda(tokens []lexer.Token, index int) (Ast, int) {
	/* lambda arg1 arg2 ... begin ... end */
	if index+2 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "not complete lambda block.")
		return nil, index
	}
	index++
	argList := make([]string, 0)
	for tokens[index].Type() == lexer.SYMBOL && index < len(tokens) {
		argList = append(argList, tokens[index].Value())
		index++
	}
	block, index := parseBlock(tokens, index)
	return Lambda_{argList, block, tokens[index].Line()}, index
}

func parseStruct(tokens []lexer.Token, index int) (Ast, int) {
	/* struct begin ... end */
	if index+2 >= len(tokens) {
		errlog.Err("parser", tokens[index].Line(), "not complete struct block.")
		return nil, index
	}
	index++
	block, index := parseBlock(tokens, index)
	return Struct_{block, tokens[index].Line()}, index
}

func parseCall(tokens []lexer.Token, index int) (Ast, int) {
	name := ParseValue(tokens[index])
	index++
	args := make([]Value, 0)
	for ; index < len(tokens) && isValue(tokens[index]); index++ {
		args = append(args, ParseValue(tokens[index]))
	}
	index--
	return Call_{name, args, tokens[index].Line()}, index
}

/*
Parse([]lexer.Token) receives a token sequence,
and translates it into AST sequence.
It returns the processed []AST and a int number,
which describes where Parse() function stop parsing. (Slice index)
*/
func Parse(tokens []lexer.Token) ([]Ast, int) {
	var TokenParseFunc = map[byte]func([]lexer.Token, int) (Ast, int){
		lexer.VAR:    parseVar,
		lexer.FILL:   parseFill,
		lexer.ENUM:   parseEnum,
		lexer.IF:     parseIf,
		lexer.DEF:    parseDef,
		lexer.STRUCT: parseStruct,
		lexer.LAMBDA: parseLambda,
		lexer.SYMBOL: parseCall,
	}
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
		if tokens[index].Type() == lexer.SEND {
			/*
			   > ... > ... # wrong
			   ... > ... > ... # right
			*/
			if index == 0 || tokens[index-1].Type() == lexer.STOP {
				errlog.Err("parser", tokens[index].Line(), "use ' > ' at the start of a sentence. There's nothing to send.")
				return nil, index
			}
			sendList = append(sendList, codes[len(codes)-1])
			codes = codes[:len(codes)-1]
		} else if tokens[index].Type() == lexer.BEGIN {
			var block []Ast
			block, index = parseBlock(tokens, index)
			codes = append(codes, Block_{block, tokens[index].Line()})
		} else if tokens[index].Type() == lexer.NUM {
			errlog.Err("parser", tokens[index].Line(), "unexpected constant number.")
			return nil, index
		} else if tokens[index].Type() == lexer.STOP {
			checkSend()
		} else if tokens[index].Type() == lexer.END {
			checkSend()
			return codes, index
		} else {
			f, ok := TokenParseFunc[tokens[index].Type()]
			if !ok {
				errlog.Err("parser", tokens[index].Line(), "Unknown token type", tokens[index].Type())
				return nil, index
			}
			var astNode Ast
			astNode, index = f(tokens, index)
			if astNode == nil { // parse failed, err has been reported
				return nil, index
			}
			codes = append(codes, astNode)
		}
	}
	checkSend()
	return codes, index
}
