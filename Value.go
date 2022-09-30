package parser

import (
	"github.com/zhangzheheng12345/flowscript/xlexer"
)

type Value interface {
	get(*Context) interface{}
}

/* int surface value */
type Int_ struct {
	value int
}

func (int_ Int_) get(ctx *Context) interface{} {
	return int_.value
}

/* variables */
type Symbol_ struct {
	// A.B.C makes a struct visiting chain
	base  string
	after []string
}

func (symbol_ Symbol_) get(ctx *Context) interface{} {
	var base interface{}
	if symbol_.base == "-" {
		base = tmpQueue.Get()
		tmpQueue.Pop()
	} else {
		base = ctx.scope.Find(symbol_.base)
	}
	// Iterating the chain
	for _, memberName := range symbol_.after {
		base = WantStruct(base).Member(memberName)
	}
	return base
}

type Str_ struct {
	value string
}

func (str_ Str_) get(ctx *Context) interface{} {
	return str_.value
}

type Char_ struct {
	value byte
}

func (char_ Char_) get(ctx *Context) interface{} {
	return char_.value
}

// TODO: Find out a way to avoid calling useless Want... while using Exp_
/* X-expression only support int value */
type Exp_ struct {
	tokens []xlexer.Token
}

func (exp_ Exp_) get(ctx *Context) interface{} {
	/* in xparse */
	return B_(E(exp_.tokens))
}
