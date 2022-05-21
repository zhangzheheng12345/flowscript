package parser

import (
	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
)

type Ast interface {
	run() interface{}
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

type Enum_ struct {
	names []string
	line  int
}

func (enum_ Enum_) run() interface{} {
	errlog.Line = enum_.line
	for _, v := range enum_.names {
		Scope.Add(v, Scope.enumCounter)
		Scope.enumCounter++
	}
	return 0
}

type Def_ struct {
	name  string
	args  []string
	codes []Ast
	line  int
}

func (def_ Def_) run() interface{} {
	errlog.Line = def_.line
	res := FlowFunc{Scope, def_.args, def_.codes}
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
	return FlowFunc{Scope, lambda_.args, lambda_.codes}
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
	tmpQueue = MakeTmpQueue(tmpQueue, len(send_.codes))
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
