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
	/* provide a independence scope for the function*/
	Scope = MakeScope(flowFunc.fathers, Scope)
	/* add the arguments to the local scope*/
	for key, arg := range args {
		Scope.Add(flowFunc.args[key], arg)
	}
	/* run the codes.
	   the last expression will give the return value.
	   the default return value is 0.
	*/
	var result interface{}
	for _, code := range flowFunc.codes {
		result = code.run()
	}
	/* delete the local variables and change the scope before leave the function*/
	Scope = Scope.Back()
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
