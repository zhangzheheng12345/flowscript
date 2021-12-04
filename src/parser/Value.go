package parser

import (
    "fmt"
)

type Value interface {
    get() int
}

type Int_ struct {
    value int
}

func (int_ Int_) get() int {
    return int_.value
}

type Symbol_ struct {
    name string
}

func (symbol_ Symbol_) get() int {
    return Scope.Find(symbol_.name)
}

type Tmp_ struct {
    // no content
}

func (tmp_ Tmp_) get() int {
    result := tmpQueue.Get()
    tmpQueue.Pop()
    return result
}