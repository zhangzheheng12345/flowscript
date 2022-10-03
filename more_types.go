package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/*
This interface units GoFunc and FlowFunc, enabling you to add both go functions & FlowScript function.
*/
type Func_ interface {
	run([]interface{}, *Context) interface{}
	argsNum() int
}

/*
Structure FlowFunc contains all the information for a function to run.
FlowFunc.args contains the arguments' name in order to add them in a new variable scope after calling the function.
FlowFunc.codes contains the codes in the function to run.
*/
type FlowFunc struct {
	fathers *Scope_
	args    []string
	codes   []Ast
}

/*
Runs the flowscript codes in the function.
*/
func (flowFunc FlowFunc) run(args []interface{}, ctx *Context) interface{} {
	/* delete the local variables and change the scope before leaving the function */
	defer func() { ctx.scope = ctx.scope.Back() }()
	ctx.scope = MakeScope(flowFunc.fathers, ctx.scope)
	if len(args) > len(flowFunc.args) {
		errlog.Err("runtime", ctx.Line, "Too many arguments while calling function.")
		return 0
	}
	/* add the arguments to the local scope*/
	for key, arg := range args {
		ctx.scope.Add(flowFunc.args[key], arg, ctx)
	}
	if len(args) < len(flowFunc.args) {
		return FlowFunc{ctx.scope, flowFunc.args[len(args):], flowFunc.codes}
	}
	/* run the codes.
	   the default return value is 0.
	*/
	var result interface{}
	for _, code := range flowFunc.codes {
		result = code.run(ctx)
	}
	return result
}

func (flowFunc FlowFunc) argsNum() int {
	return len(flowFunc.args)
}

/* Struct ( with no underline ) is a type */
type Struct struct {
	members map[string]interface{}
}

func (struct_ Struct) Member(name string, ctx *Context) interface{} {
	v, ok := struct_.members[name]
	if ok {
		return v
	} else {
		errlog.Err("runtime", ctx.Line, "Try to find a no-existing member in the structure.")
		return nil // nil is safe, as the type checkings (Want...) will process them
	}
}
