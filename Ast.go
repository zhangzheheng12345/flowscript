package parser

import (
	"fmt"
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
)

type Ast interface {
	run() interface{}
}

type Add_ struct {
	op1  Value
	op2  Value
	line int
}

func (add_ Add_) run() interface{} {
	errlog.Line = add_.line
	v1 := add_.op1.get()
	v2 := add_.op2.get()
	switch v := v1.(type) {
	case int:
		return v + WantInt(v2)
	case byte:
		return int(v) + WantInt(v2)
	case []interface{}:
		/* Connect two lists */
		return append(v, WantList(v2)...)
	case string:
		return v + WantString(v2)
	default:
		errlog.Err("runtime", add_.line, "Try to add a unkown type value with another")
		return nil
	}
}

type Sub_ struct {
	op1  Value
	op2  Value
	line int
}

func (sub_ Sub_) run() interface{} {
	errlog.Line = sub_.line
	return WantInt(sub_.op1.get()) - WantInt(sub_.op2.get())
}

type Multi_ struct {
	op1  Value
	op2  Value
	line int
}

func (multi_ Multi_) run() interface{} {
	errlog.Line = multi_.line
	return WantInt(multi_.op1.get()) * WantInt(multi_.op2.get())
}

type Div_ struct {
	op1  Value
	op2  Value
	line int
}

func (div_ Div_) run() interface{} {
	return WantInt(div_.op1.get()) / WantInt(div_.op2.get())
}

type Mod_ struct {
	op1  Value
	op2  Value
	line int
}

func (mod_ Mod_) run() interface{} {
	errlog.Line = mod_.line
	return WantInt(mod_.op1.get()) % WantInt(mod_.op2.get())
}

type Bigr_ struct {
	op1  Value
	op2  Value
	line int
}

func (bigr_ Bigr_) run() interface{} {
	errlog.Line = bigr_.line
	return BoolToInt(WantInt(bigr_.op1.get()) > WantInt(bigr_.op2.get()))
}

type Smlr_ struct {
	op1  Value
	op2  Value
	line int
}

func (smlr_ Smlr_) run() interface{} {
	errlog.Line = smlr_.line
	return BoolToInt(WantInt(smlr_.op1.get()) < WantInt(smlr_.op2.get()))
}

type Equal_ struct {
	op1  Value
	op2  Value
	line int
}

func (equal_ Equal_) run() interface{} {
	errlog.Line = equal_.line
	v1 := equal_.op1.get()
	v2 := equal_.op2.get()
	switch v := v1.(type) {
	case int:
		return BoolToInt(v == WantInt(v2))
	case byte:
		return BoolToInt(int(v) == WantInt(v2))
	case []interface{}:
		/* Compare two lists */
		li2 := WantList(v2)
		for i, value := range v {
			if value != li2[i] {
				return 0
			}
		}
		return 1
	case string:
		return BoolToInt(v == WantString(v2))
	default:
		errlog.Err("runtime", equal_.line, "Try to compare a unkown type value to another. Operation: equal")
		return nil
	}
}

type And_ struct {
	op1  Value
	op2  Value
	line int
}

func (and_ And_) run() interface{} {
	errlog.Line = and_.line
	return WantInt(and_.op1.get()) & WantInt(and_.op2.get())
}

type Xor_ struct {
	op1  Value
	op2  Value
	line int
}

func (xor_ Xor_) run() interface{} {
	errlog.Line = xor_.line
	return WantInt(xor_.op1.get()) ^ WantInt(xor_.op1.get())
}

type Or_ struct {
	op1  Value
	op2  Value
	line int
}

func (or_ Or_) run() interface{} {
	errlog.Line = or_.line
	return WantInt(or_.op1.get()) | WantInt(or_.op1.get())
}

type Not_ struct {
	op   Value
	line int
}

func (not_ Not_) run() interface{} {
	errlog.Line = not_.line
	return BoolToInt(WantInt(not_.op.get()) == 0)
}

type Expr_ struct {
	op   Value
	line int
}

func (expr_ Expr_) run() interface{} {
	errlog.Line = expr_.line
	return expr_.op.get()
}

type Var_ struct {
	name string
	op   Value
	line int
}

func (var_ Var_) run() interface{} {
	errlog.Line = var_.line
	if var_.op != nil {
		/* give a start value */
		Scope.Add(var_.name, var_.op.get())
	} else {
		/* default 0 */
		Scope.Add(var_.name, 0)
	}
	return 0
}

/* type cmd returns a string which tells the operand's type */
type Type_ struct {
	op   Value
	line int
}

func (type_ Type_) run() interface{} {
	errlog.Line = type_.line
	switch v := type_.op.get().(type) {
	case int:
		_ = v // FIXME: why must i write like this
		return "int"
	case byte:
		return "char"
	case []interface{}:
		return "list"
	case string:
		return "string"
	case Func_:
		return "function"
	case Struct:
		return "struct"
	default:
		return "?unknown_type?"
	}
}

type Len_ struct {
	op   Value
	line int
}

func (len_ Len_) run() interface{} {
	errlog.Line = len_.line
	switch v := len_.op.get().(type) {
	case int, byte:
		errlog.Err("runtime", len_.line, "Try to get the length of a int or char value.")
		return 0
	case []interface{}:
		return len(v)
	case string:
		return len(strings.Split(v, ""))
	default:
		errlog.Err("runtime", len_.line, "Try to get the length of a unknown type value.")
		return 0
	}
}

type Index_ struct {
	op    Value
	index Value
	line  int
}

func (index_ Index_) run() interface{} {
	errlog.Line = index_.line
	index := WantInt(index_.index.get())
	switch v := index_.op.get().(type) {
	case int, byte:
		errlog.Err("runtime", index_.line, "Try to index a int or char value.")
		return 0
	case []interface{}:
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
		errlog.Err("runtime", index_.line, "Try to index a unknown type value.")
		return 0
	}
}

type App_ struct {
	op1  Value
	op2  Value
	line int
}

func (app_ App_) run() interface{} {
	errlog.Line = app_.line
	switch v := app_.op1.get().(type) {
	case int, byte:
		errlog.Err("runtime", app_.line, "Try to append a value after a int or char value.")
		return 0
	case []interface{}:
		return append(v, app_.op2.get())
	case string:
		return v + string([]byte{byte(WantInt(app_.op2.get()))})
	default:
		errlog.Err("runtime", app_.line, "Try to append sth after a unknown type value.")
		return 0
	}
}

type Slice_ struct {
	op1  Value
	op2  Value
	op3  Value
	line int
}

func (slice_ Slice_) run() interface{} {
	errlog.Line = slice_.line
	begin := WantInt(slice_.op2.get())
	end := WantInt(slice_.op3.get())
	var err1, err2 bool
	switch v := slice_.op1.get().(type) {
	case int, byte:
		errlog.Err("runtime", slice_.line, "Try to slice a int or char value.")
		return 0
	case []interface{}:
		begin, err1 = AbsIndex(len(v)+1, begin)
		end, err2 = AbsIndex(len(v)+1, end)
		if err1 || err2 {
			return []interface{}{}
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
		errlog.Err("runtime", slice_.line, "Try to slice a unknown type value.")
		return 0
	}
}

/* Split the string by words */
type Words_ struct {
	op   Value
	line int
}

func (words_ Words_) run() interface{} {
	errlog.Line = words_.line
	str := WantString(words_.op.get())
	return strings.Split(str, " ")
}

/* Split the string by lines*/
type Lines_ struct {
	op   Value
	line int
}

func (lines_ Lines_) run() interface{} {
	errlog.Line = lines_.line
	str := WantString(lines_.op.get())
	return strings.Split(str, "\n")
}

type List_ struct {
	ops  []Value
	line int
}

func (list_ List_) run() interface{} {
	errlog.Line = list_.line
	res := make([]interface{}, 0)
	for _, v := range list_.ops {
		res = append(res, v.get())
	}
	return res
}

type Def_ struct {
	name  string
	args  []string
	codes []Ast
	line  int
}

func (def_ Def_) run() interface{} {
	errlog.Line = def_.line
	res := Func_{Scope, def_.args, def_.codes}
	Scope.Add(def_.name, res)
	return res
}

type Lambda_ struct {
	args  []string
	codes []Ast
	line  int
}

func (lambda_ Lambda_) run() interface{} {
	errlog.Line = lambda_.line
	return Func_{Scope, lambda_.args, lambda_.codes}
}

/* Struct_ ( with one underline ) is the command */
type Struct_ struct {
	codes []Ast
	line  int
}

func (struct_ Struct_) run() interface{} {
	errlog.Line = struct_.line
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
	line  int
}

func (send_ Send_) run() interface{} {
	errlog.Line = send_.line
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
	line  int
}

func (block_ Block_) run() interface{} {
	errlog.Line = block_.line
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
	line      int
}

func (if_ If_) run() interface{} {
	errlog.Line = if_.line
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
	line int
}

func (call_ Call_) run() interface{} {
	errlog.Line = call_.line
	argsValue := make([]interface{}, 0)
	for _, arg := range call_.args {
		argsValue = append(argsValue, arg.get())
	}
	return WantFunc(call_.name.get()).run(argsValue)
}

/* output to stdout */
type Echo_ struct {
	op   Value
	line int
}

func (echo_ Echo_) run() interface{} {
	errlog.Line = echo_.line
	op := echo_.op.get()
	v, ok := op.(byte)
	if ok {
		fmt.Print(string([]byte{v}))
	} else {
		fmt.Print(op)
	}
	return 0
}

/* output to stdout */
type Echoln_ struct {
	op   Value
	line int
}

func (echoln_ Echoln_) run() interface{} {
	errlog.Line = echoln_.line
	op := echoln_.op.get()
	v, ok := op.(byte)
	if ok {
		fmt.Println(string([]byte{v}))
	} else {
		fmt.Println(op)
	}
	return 0
}

/* Read from stdin */
type Input_ struct {
	// op is a reminding string
	op   Value
	line int
}

func (input_ Input_) run() interface{} {
	// errlog.Line = input_.line is not needed because Echo_{input_.op, input_.line}.run() has already done it
	Echo_(input_).run()
	var res string
	_, err := fmt.Scan(&res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return res
}
