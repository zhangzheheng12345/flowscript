package parser

import (
	"github.com/zhangzheheng12345/FlowScript/xlexer"
)

type Value interface {
	get() interface{}
}

/* int surface value */
type Int_ struct {
	value int
}

func (int_ Int_) get() interface{} {
	return int_.value
}

/* variables */
type Symbol_ struct {
	// A.B.C makes a struct visiting chain
	base  string
	after []string
}

func (symbol_ Symbol_) get() interface{} {
	// TODO: Will len(symbol_.name) == 0 ?
	base := Scope.Find(symbol_.base)
	// Iterating the chain
	for _, memberName := range symbol_.after {
		base = WantStruct(base).Member(memberName)
	}
	return base
}

type Str_ struct {
	value string
}

func (str_ Str_) get() interface{} {
	return str_.value
}

type Char_ struct {
	value byte
}

func (char_ Char_) get() interface{} {
	return char_.value
}

/* get value from the send value queue*/
type Tmp_ struct {
	after []string
}

func (tmp_ Tmp_) get() interface{} {
	base := tmpQueue.Get()
	tmpQueue.Pop()
	for _, v := range tmp_.after {
		base = WantStruct(base).Member(v)
	}
	return base
}

// TODO: Find out a way to avoid call useless Want... while using Exp_
/* X-expression only support int value */
type Exp_ struct {
	tokens []xlexer.Token
}

func (exp_ Exp_) get() interface{} {
	/* in xparse */
	return E_(T(exp_.tokens))
}
