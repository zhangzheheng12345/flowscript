package parser

import (
    "fmt"
)

/*
Structure Func_ contains all the information for a function to run.
Func_.args contains the arguments' name in order to add them in a new variable scope after calling the function.
Func_.codes contains the codes in the function to run later.
*/
type Func_ struct {
	fathers *Scope_
	fatherf *FuncScope_
	args    []string
	codes   []Ast
}

/*
Func_.run([]int) run the codes in the function.
*/
func (func_ Func_) run(args []interface{}) interface{} {
	/* provide a independence scope for the function*/
	Scope = MakeScope(func_.fathers, Scope)
	FuncScope = MakeFuncScope(func_.fatherf, FuncScope)
	/* add the arguments to the local scope*/
	for key, arg := range args {
		Scope.Add(func_.args[key], arg)
	}
	/* run the codes.
	   the last expression will give the return value.
	   the default return value is 0.
	*/
	var result interface{}
	for _, code := range func_.codes {
		result = code.run()
	}
	/* delete the local variables and change the scope before leave the function*/
	Scope = Scope.Back()
	FuncScope = FuncScope.Back()
	return result
}

type Struct_ struct {
	members map[string]interface{}
}

func (struct_ Struct_) Member(name string) interface{} {
	v, ok := struct_.members[string]
	if ok {
		return v
	} else {
		fmt.Println("Error: Try to find a no-existing member in the structure.")
		return nil // TODO:Is nil safe?
	}
}