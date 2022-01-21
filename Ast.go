package parser

import (
	"fmt"

	"github.com/zhangzheheng12345/FlowScript/tools"
)

/* Tool function, for cannot convert bool to int directly */
type Ast interface {
	run() interface{}
}

// TODO:Support more type functions. add, multi should support string and list, and len, index... functions
type Add_ struct {
	op1 Value
	op2 Value
}

func (add_ Add_) run() interface{} {
	v1 := add_.op1.get()
	v2 := add_.op2.get()
	switch v := v1.(type) {
	case int:
		return v + tools.WantInt(v2)
    case byte:
		return int(v) + tools.WantInt(v2)
	case []int:
		/* Connect two lists */
		return append(v, tools.WantIntList(v2)...)
	case string:
		return v + tools.WantString(v2)
    default:
        fmt.Println("Error: try to add a unkown type value with another")
        return nil
	}
}

type Sub_ struct {
	op1 Value
	op2 Value
}

func (sub_ Sub_) run() interface{} {
	return tools.WantInt(sub_.op1.get()) - tools.WantInt(sub_.op2.get())
}

type Multi_ struct {
	op1 Value
	op2 Value
}

func (multi_ Multi_) run() interface{} {
	return tools.WantInt(multi_.op1.get()) * tools.WantInt(multi_.op2.get())
}

type Div_ struct {
	op1 Value
	op2 Value
}

func (div_ Div_) run() interface{} {
	return tools.WantInt(div_.op1.get()) / tools.WantInt(div_.op2.get())
}

type Mod_ struct {
	op1 Value
	op2 Value
}

func (mod_ Mod_) run() interface{} {
	return tools.WantInt(mod_.op1.get()) % tools.WantInt(mod_.op2.get())
}

type Bigr_ struct {
	op1 Value
	op2 Value
}

func (bigr_ Bigr_) run() interface{} {
	return tools.BoolToInt(tools.WantInt(bigr_.op1.get()) > tools.WantInt(bigr_.op2.get()))
}

type Smlr_ struct {
	op1 Value
	op2 Value
}

func (smlr_ Smlr_) run() interface{} {
	return tools.BoolToInt(tools.WantInt(smlr_.op1.get()) < tools.WantInt(smlr_.op2.get()))
}

type Equal_ struct {
	op1 Value
	op2 Value
}

func (equal_ Equal_) run() interface{} {
	v1 := equal_.op1.get()
	v2 := equal_.op2.get()
	switch v := v1.(type) {
	case int:
		return tools.BoolToInt(v == tools.WantInt(v2))
    case byte:
		return tools.BoolToInt(int(v) == tools.WantInt(v2))
	case []int:
		/* Compare two lists */
		li2 := tools.WantIntList(v2)
		for i, value := range v {
			if value != li2[i] {
				return 0
			}
		}
		return 1
	case string:
		return tools.BoolToInt(v == tools.WantString(v2))
    default:
        fmt.Println("Error: try to compare a unkown type value to another. Operation: equal")
        return nil
	}
}

type And_ struct {
	op1 Value
	op2 Value
}

func (and_ And_) run() interface{} {
	return tools.BoolToInt((tools.WantInt(and_.op1.get()) != 0) && (tools.WantInt(and_.op2.get()) != 0))
}

type Or_ struct {
	op1 Value
	op2 Value
}

func (or_ Or_) run() interface{} {
	return tools.BoolToInt((tools.WantInt(or_.op1.get()) != 0) || (tools.WantInt(or_.op1.get()) != 0))
}

type Not_ struct {
	op Value
}

func (not_ Not_) run() interface{} {
	return tools.BoolToInt(tools.WantInt(not_.op.get()) == 0)
}

type Var_ struct {
	name string
	op   Value
}

func (var_ Var_) run() interface{} {
	if var_.op != nil {
		/* give a start value */
		Scope.Add(var_.name, var_.op.get())
	} else {
		if tmpQueue.Size() == 0 {
			/* default 0 */
			Scope.Add(var_.name, 0)
		} else {
			/* autoly pick value from send queue */
			Scope.Add(var_.name, tmpQueue.Get())
			tmpQueue.Pop()
		}
	}
	return 0
}

type Def_ struct {
	name  string
	args  []string
	codes []Ast
}

func (def_ Def_) run() interface{} {
	FuncScope.Add(def_.name, Func_{Scope, FuncScope, def_.args, def_.codes})
	return 0
}

/* a send sequence */
type Send_ struct {
	codes []Ast
}

func (send_ Send_) run() interface{} {
	tmpQueue = MakeTmpQueue(tmpQueue)
	for _, code := range send_.codes {
		tmpQueue.Add(code.run())
	}
	result := tmpQueue.Get()
	tmpQueue = tmpQueue.Clear()
	return result
}

/*
begin
    ...
end
*/
type Block_ struct {
	codes []Ast
}

func (block_ Block_) run() interface{} {
	Scope = MakeScope(Scope, Scope)
	FuncScope = MakeFuncScope(FuncScope, FuncScope)
	var result interface{}
	for _, code := range block_.codes {
		result = code.run()
	}
	Scope = Scope.Back()
	FuncScope = FuncScope.Back()
	return result
}

/*
if condition begin
    ...
end
*/
type If_ struct {
	condition Value
	codes     []Ast
}

func (if_ If_) run() interface{} {
	if if_.condition.get() != 0 {
		Scope = MakeScope(Scope, Scope)
		FuncScope = MakeFuncScope(FuncScope, FuncScope)
		var result interface{}
		for _, code := range if_.codes {
			result = code.run()
		}
		Scope = Scope.Back()
		FuncScope = FuncScope.Back()
		return result
	}
	return 0
}

/* call a user defined function*/
type Call_ struct {
	name string
	args []Value
}

func (call_ Call_) run() interface{} {
	argsValue := make([]interface{}, 0)
	for _, arg := range call_.args {
		argsValue = append(argsValue, arg.get())
	}
	return FuncScope.Find(call_.name).run(argsValue)
}

/* output to stdout */
type Echo_ struct {
	op Value
}

func (echo_ Echo_) run() interface{} {
    op := echo_.op.get()
    v, ok := op.(byte)
    if ok {
        fmt.Println(string([]byte{v}))
    } else {
	    fmt.Println(op)
    }
	return 0
}
