package parser

import (
	"strconv"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/xlexer"
)

func B_(tokens []xlexer.Token, value int) int {
	if len(tokens) == 0 {
		return value
	}
	switch tokens[0].Type() {
	case xlexer.AND:
		return value & B_(E(tokens[1:]))
	case xlexer.OR:
		return value | B_(E(tokens[1:]))
	case xlexer.XOR:
		return value ^ B_(E(tokens[1:]))
	default:
		errlog.Err("xparser", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return 0
	}
}

func E(tokens []xlexer.Token) ([]xlexer.Token, int) {
	return E_(T(tokens))
}

func E_(tokens []xlexer.Token, value int) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, value
	}
	switch tokens[0].Type() {
	case xlexer.AND, xlexer.OR, xlexer.XOR:
		return tokens, value
	case xlexer.ADD:
		tail, v := E_(T(tokens[1:]))
		return tail, value + v
	case xlexer.SUB:
		tail, v := T(tokens[1:])
		return E_(tail, value-v)
	default:
		errlog.Err("xparse", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}

func T(tokens []xlexer.Token) ([]xlexer.Token, int) {
	return T_(F(tokens))
}

func T_(tokens []xlexer.Token, value int) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, value
	}
	switch tokens[0].Type() {
	case xlexer.ADD, xlexer.SUB, xlexer.AND, xlexer.OR, xlexer.XOR:
		return tokens, value
	case xlexer.MULTI:
		tail, v := T_(F(tokens[1:]))
		return tail, value * v
	case xlexer.DIV:
		tail, v := F(tokens[1:])
		if v == 0 {
			errlog.Err("runtime", errlog.Line, "Cannot div 0")
			return T_(tail, 0)
		}
		return T_(tail, value/v)
	case xlexer.MOD:
		tail, v := F(tokens[1:])
		if v == 0 {
			errlog.Err("runtime", errlog.Line, "Cannot mod 0")
			return T_(tail, 0)
		}
		return T_(tail, value%v)
	default:
		errlog.Err("xparse", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}

func F(tokens []xlexer.Token) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, 1
	}
	switch tokens[0].Type() {
	case xlexer.SYMBOL:
		return tokens[1:], WantInt(Scope.Find(tokens[0].Value()))
	case xlexer.TMP:
		result := tmpQueue.Get()
		tmpQueue.Pop()
		return tokens[1:], WantInt(result)
	case xlexer.NUM:
		num, _ := strconv.Atoi(tokens[0].Value())
		return tokens[1:], num
	case xlexer.FPAREN:
		count := 1
		index := 1
		for count != 0 && index < len(tokens) {
			if tokens[index].Type() == xlexer.FPAREN {
				count++
			} else if tokens[index].Type() == xlexer.BPAREN {
				count--
			}
			index++
		}
		if index >= len(tokens) && tokens[index-1].Type() != xlexer.BPAREN {
			errlog.Err("xparser", tokens[index-1].Line(), "lose back parenthesis in xparser")
		}
		return tokens[index:], WantInt(Exp_{tokens[1 : index-1]}.get())
	default:
		errlog.Err("xparse", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}
