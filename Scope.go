package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/*
Global scope
*/
var Global *Context = &Context{
	&Scope_{nil, nil, make(map[string]interface{}), 0},
	1,
}

/*
exp1 > exp2 > exp3
The values of the send sequence expressions will be pushed to this queue.
*/
var tmpQueue *Queue_ = MakeTmpQueue(nil, 0)

/* scope system */
type Scope_ struct {
	last        *Scope_
	father      *Scope_
	vars        map[string]interface{}
	enumCounter int
}

func MakeScope(father *Scope_, last *Scope_) *Scope_ {
	return &Scope_{last, father, make(map[string]interface{}), father.enumCounter}
}

func (scope_ Scope_) Add(key string, value interface{}, ctx *Context) {
	_, ok := scope_.vars[key]
	if ok {
		errlog.Err("runtime", ctx.Line, "Try to add a variable that has been added. name:", key)
	} else {
		scope_.vars[key] = value
	}
}

func (scope_ Scope_) Find(key string, ctx *Context) interface{} {
	v, ok := scope_.vars[key]
	if ok {
		return v
	} else if scope_.father != nil {
		return scope_.father.Find(key, ctx)
	} else {
		errlog.Err("runtime", ctx.Line, "Try to read a variable that hasn't been added. name:", key)
	}
	return 0
}

func (scope_ Scope_) Back() *Scope_ {
	return scope_.last
}

type Queue_ struct {
	father  *Queue_
	data    []interface{}
	dataLen int
}

func MakeTmpQueue(father *Queue_, initSize int) *Queue_ {
	return &Queue_{father, make([]interface{}, initSize), 0}
}

func (queue_ *Queue_) Add(value interface{}) {
	queue_.dataLen++
	queue_.data[queue_.dataLen-1] = value
}

func (queue_ Queue_) Get(ctx *Context) interface{} {
	if queue_.dataLen > 0 {
		return queue_.data[0]
	}
	errlog.Err("runtime", ctx.Line, "Try to get a value from temp queue while it's empty")
	return nil
}

func (queue_ *Queue_) Pop() {
	queue_.dataLen--
	queue_.data = queue_.data[1:]
}

func (queue_ *Queue_) Clear() *Queue_ {
	queue_.data = nil
	return queue_.father
}

func (queue_ Queue_) Size() int {
	return queue_.dataLen
}
