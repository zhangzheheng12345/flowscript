package parser

import (
	"fmt"
	"strings"
)

type Ast interface {
	run() interface{}
}

type Add_ struct {
	op1 Value
	op2 Value
}

func (add_ Add_) run() interface{} {
	v1 := add_.op1.get()
	v2 := add_.op2.get()
	switch v := v1.(type) {
	case int:
		return v + WantInt(v2)
	case byte:
		return int(v) + WantInt(v2)
	case []int:
		/* Connect two lists */
		return append(v, WantIntList(v2)...)
	case string:
		return v + WantString(v2)
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
	return WantInt(sub_.op1.get()) - WantInt(sub_.op2.get())
}

type Multi_ struct {
	op1 Value
	op2 Value
}

func (multi_ Multi_) run() interface{} {
	return WantInt(multi_.op1.get()) * WantInt(multi_.op2.get())
}

type Div_ struct {
	op1 Value
	op2 Value
}

func (div_ Div_) run() interface{} {
	return WantInt(div_.op1.get()) / WantInt(div_.op2.get())
}

type Mod_ struct {
	op1 Value
	op2 Value
}

func (mod_ Mod_) run() interface{} {
	return WantInt(mod_.op1.get()) % WantInt(mod_.op2.get())
}

type Bigr_ struct {
	op1 Value
	op2 Value
}

func (bigr_ Bigr_) run() interface{} {
	return BoolToInt(WantInt(bigr_.op1.get()) > WantInt(bigr_.op2.get()))
}

type Smlr_ struct {
	op1 Value
	op2 Value
}

func (smlr_ Smlr_) run() interface{} {
	return BoolToInt(WantInt(smlr_.op1.get()) < WantInt(smlr_.op2.get()))
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
		return BoolToInt(v == WantInt(v2))
	case byte:
		return BoolToInt(int(v) == WantInt(v2))
	case []int:
		/* Compare two lists */
		li2 := WantIntList(v2)
		for i, value := range v {
			if value != li2[i] {
				return 0
			}
		}
		return 1
	case string:
		return BoolToInt(v == WantString(v2))
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
	return WantInt(and_.op1.get()) & WantInt(and_.op2.get())
}

type Xor_ struct {
	op1 Value
	op2 Value
}

func (xor_ Xor_) run() interface{} {
	return WantInt(xor_.op1.get()) ^ WantInt(xor_.op1.get())
}

type Or_ struct {
	op1 Value
	op2 Value
}

func (or_ Or_) run() interface{} {
	return WantInt(or_.op1.get()) | WantInt(or_.op1.get())
}

type Not_ struct {
	op Value
}

func (not_ Not_) run() interface{} {
	return BoolToInt(WantInt(not_.op.get()) == 0)
}

type Expr_ struct {
	op Value
}

func (expr_ Expr_) run() interface{} {
	return expr_.op.get()
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

type Len_ struct {
	op Value
}

func (len_ Len_) run() interface{} {
	switch v := len_.op.get().(type) {
	case int, byte:
		fmt.Println("Error: Try to get the length of a int or char value.")
		return 0
	case []int:
		return len(v)
	case string:
		return len(strings.Split(v, ""))
	default:
		fmt.Println("Error: Try to get the length of a unknown type value.")
		return 0
	}
}

type Index_ struct {
	op    Value
	index Value
}

func (index_ Index_) run() interface{} {
	index := WantInt(index_.index.get())
	switch v := index_.op.get().(type) {
	case int, byte:
		fmt.Println("Error: Try to index a int or char value.")
		return 0
	case []int:
		abs, err := AbsIndex(len(v), index)
		if err {
			return 0
		} else {
			return v[abs]
		}
	case string:
		abs, err := AbsIndex(len(strings.Split(v, "")), index)
		if err {
			return ""
		} else {
			return strings.Split(v, "")[abs]
		}
	default:
		fmt.Println("Error: Try to index a unknown type value.")
		return 0
	}
}

type App_ struct {
	op1 Value
	op2 Value
}

func (app_ App_) run() interface{} {
	switch v := app_.op1.get().(type) {
	case int, byte:
		fmt.Println("Error: Try to append a value after a int or char value.")
		return 0
	case []int:
		return append(v, WantInt(app_.op2.get()))
	case string:
		return v + string([]byte{byte(WantInt(app_.op2.get()))})
	default:
		fmt.Println("Error: Try to append sth after a unknown type value.")
		return 0
	}
}

type Slice_ struct {
	op1 Value
	op2 Value
	op3 Value
}

func (slice_ Slice_) run() interface{} {
	// TODO: Support minus index in slice
	begin := WantInt(slice_.op2.get())
	end := WantInt(slice_.op3.get())
	var err1, err2 bool
	switch v := slice_.op1.get().(type) {
	case int, byte:
		fmt.Println("Error: Try to slice a int or char value.")
		return 0
	case []int:
		begin, err1 = AbsIndex(len(v), begin)
		end, err2 = AbsIndex(len(v), end)
		if err1 || err2 {
			return 0
		} else {
			return v[begin:end]
		}
	case string:
		splited := strings.Split(v, "")
		begin, err1 = AbsIndex(len(splited), begin)
		end, err2 = AbsIndex(len(splited), end)
		if err1 || err2 {
			return 0
		} else {
			return strings.Join(splited[begin:end], "")
		}
	default:
		fmt.Println("Error: Try to slice a unknown type value.")
		return 0
	}
}

type List_ struct {
	ops []Value
}

func (list_ List_) run() interface{} {
	res := make([]int, 0)
	for _, v := range list_.ops {
		res = append(res, WantInt(v.get()))
	}
	return res
}

type Def_ struct {
	name  string
	args  []string
	codes []Ast
}

func (def_ Def_) run() interface{} {
	res := Func_{Scope, def_.args, def_.codes}
	Scope.Add(def_.name, res)
	return res
}

type Lambda_ struct {
	args  []string
	codes []Ast
}

func (lambda_ Lambda_) run() interface{} {
	return Func_{Scope, lambda_.args, lambda_.codes}
}

/* Struct_ ( with one underline ) is the command */
type Struct_ struct {
	codes []Ast
}

func (struct_ Struct_) run() interface{} {
	Scope = MakeScope(Scope, Scope)
	for _, code := range struct_.codes {
		code.run()
	}
	result := Scope.vars
	Scope = Scope.Back()
	return Struct{result}
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
	var result interface{}
	for _, code := range block_.codes {
		result = code.run()
	}
	Scope = Scope.Back()
	return result
}

/*
if condition begin
    ...
end
*/
type If_ struct {
	condition Value
	ifcodes   []Ast
	elsecodes []Ast
}

func (if_ If_) run() interface{} {
	if WantInt(if_.condition.get()) != 0 {
		Scope = MakeScope(Scope, Scope)
		var result interface{}
		for _, code := range if_.ifcodes {
			result = code.run()
		}
		Scope = Scope.Back()
		return result
	} else {
		Scope = MakeScope(Scope, Scope)
		var result interface{} = 0
		for _, code := range if_.elsecodes {
			result = code.run()
		}
		Scope = Scope.Back()
		return result
	}
}

/* call a user defined function*/
type Call_ struct {
	name Value
	args []Value
}

func (call_ Call_) run() interface{} {
	argsValue := make([]interface{}, 0)
	for _, arg := range call_.args {
		argsValue = append(argsValue, arg.get())
	}
	return WantFunc(call_.name.get()).run(argsValue)
}

/* output to stdout */
type Echo_ struct {
	op Value
}

func (echo_ Echo_) run() interface{} {
	op := echo_.op.get()
	v, ok := op.(byte)
	if ok {
		fmt.Print(string([]byte{v}))
	} else {
		fmt.Print(op)
	}
	return 0
}

/* Read from stdin */
type Input_ struct {
	// op is a reminding string
	op Value
}

func (input_ Input_) run() interface{} {
	Echo_{input_.op}.run()
	var res string
	_, err := fmt.Scan(&res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return res
}
