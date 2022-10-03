package parser

import (
	"strconv"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/xlexer"
)

func B_(tokens []xlexer.Token, value int, ctx *Context) int {
	if len(tokens) == 0 {
		return value
	}
	rest, res := E(tokens[1:], ctx)
	switch tokens[0].Type() {
	case xlexer.AND:
		return value & B_(rest, res, ctx)
	case xlexer.OR:
		return value | B_(rest, res, ctx)
	case xlexer.XOR:
		return value ^ B_(rest, res, ctx)
	default:
		errlog.Err("xparser", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return 0
	}
}

func E(tokens []xlexer.Token, ctx *Context) ([]xlexer.Token, int) {
	rest, res := T(tokens, ctx)
	return E_(rest, res, ctx)
}

func E_(tokens []xlexer.Token, value int, ctx *Context) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, value
	}
	switch tokens[0].Type() {
	case xlexer.AND, xlexer.OR, xlexer.XOR:
		return tokens, value
	case xlexer.ADD:
		tail, v := T(tokens[1:], ctx)
		return E_(tail, value+v, ctx)
	case xlexer.SUB:
		tail, v := T(tokens[1:], ctx)
		return E_(tail, value-v, ctx)
	default:
		errlog.Err("xparse", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}

func T(tokens []xlexer.Token, ctx *Context) ([]xlexer.Token, int) {
	rest, res := F(tokens, ctx)
	return T_(rest, res, ctx)
}

func T_(tokens []xlexer.Token, value int, ctx *Context) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, value
	}
	switch tokens[0].Type() {
	case xlexer.ADD, xlexer.SUB, xlexer.AND, xlexer.OR, xlexer.XOR:
		return tokens, value
	}
	tail, v := F(tokens[1:], ctx)
	switch tokens[0].Type() {
	case xlexer.MULTI:
		return T_(tail, value*v, ctx)
	case xlexer.DIV:
		if v == 0 {
			errlog.Err("runtime", ctx.Line, "Cannot div 0")
			return T_(tail, 0, ctx)
		}
		return T_(tail, value/v, ctx)
	case xlexer.MOD:
		if v == 0 {
			errlog.Err("runtime", ctx.Line, "Cannot mod 0")
			return T_(tail, 0, ctx)
		}
		return T_(tail, value%v, ctx)
	default:
		errlog.Err("xparser", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}

func F(tokens []xlexer.Token, ctx *Context) ([]xlexer.Token, int) {
	if len(tokens) == 0 {
		return tokens, 1
	}
	switch tokens[0].Type() {
	case xlexer.SYMBOL:
		return tokens[1:], WantInt(ctx.scope.Find(tokens[0].Value(), ctx), ctx)
	case xlexer.TMP:
		result := tmpQueue.Get(ctx)
		tmpQueue.Pop()
		return tokens[1:], WantInt(result, ctx)
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
		return tokens[index:], WantInt(Exp_{tokens[1 : index-1]}.get(ctx), ctx)
	default:
		errlog.Err("xparser", tokens[0].Line(), "unexpected token in xparser:", tokens[0].Type())
		return tokens[1:], 0
	}
}
