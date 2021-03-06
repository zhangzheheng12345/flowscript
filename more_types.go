package parser

import errlog "github.com/zhangzheheng12345/flowscript/error_logger"

/*
This interface units GoFunc and FlowFunc, enabling you to add both go functions & FlowScript function.
*/
type Func_ interface {
	run([]interface{}) interface{}
}

/*
Structure FlowFunc contains all the information for a function to run.
FlowFunc.args contains the arguments' name in order to add them in a new variable scope after calling the function.
FlowFunc.codes contains the codes in the function to run later.
*/
type FlowFunc struct {
	fathers *Scope_
	args    []string
	codes   []Ast
}

/*
FlowFunc.run([]int) run the codes in the function.
*/
func (flowFunc FlowFunc) run(args []interface{}) interface{} {
	/* delete the local variables and change the scope before leave the function */
	defer func() { Scope = Scope.Back() }()
	Scope = MakeScope(flowFunc.fathers, Scope)
	if len(args) > len(flowFunc.args) {
		errlog.Err("runtime", errlog.Line, "Too many arguments while calling function.")
		return 0
	}
	/* add the arguments to the local scope*/
	for key, arg := range args {
		Scope.Add(flowFunc.args[key], arg)
	}
	if len(args) < len(flowFunc.args) {
		return FlowFunc{Scope, flowFunc.args[len(args):], flowFunc.codes}
	}
	/* run the codes.
	   the default return value is 0.
	*/
	var result interface{}
	for _, code := range flowFunc.codes {
		result = code.run()
	}
	return result
}

/* Struct ( no underlines ) is a type */
type Struct struct {
	members map[string]interface{}
}

func (struct_ Struct) Member(name string) interface{} {
	v, ok := struct_.members[name]
	if ok {
		return v
	} else {
		errlog.Err("runtime", errlog.Line, "Try to find a no-existing member in the structure.")
		return nil // nil is safe, as the type checkings (Want...) will process them
	}
}
