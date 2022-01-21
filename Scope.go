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

/*
Global variable, Scope, contains all the variable at present.
If someone want to add global variable by Golang in a embeding way,
he should Add(string,int) variable to Scope directly,
because all the codes run with RunCode(string) will be inside a independence Block_,
so you can't add global variable / function by native FlowScript.
*/
var Scope *Scope_ = &Scope_{nil, nil, make(map[string]interface{})}
var FuncScope *FuncScope_ = &FuncScope_{nil, nil, make(map[string]Func_)}

/*
exp1 > exp2 > exp3
The values of the send sequence expressions will be pushed to this queue.
*/
var tmpQueue *Queue_ = nil

/* scope system ( only for variables ) */
type Scope_ struct {
	last   *Scope_
	father *Scope_
	vars   map[string]interface{}
}

func MakeScope(father *Scope_, last *Scope_) *Scope_ {
	return &Scope_{last, father, make(map[string]interface{})}
}

func (scope_ Scope_) Add(key string, value interface{}) {
	_, ok := scope_.vars[key]
	if ok {
		fmt.Println("Error: Try to add a variable that has been added. var: " + key)
	} else {
		scope_.vars[key] = value
	}
}

func (scope_ Scope_) Find(key string) interface{} {
	v, ok := scope_.vars[key]
	if ok {
		return v
	} else if scope_.father != nil {
		return scope_.father.Find(key)
	} else {
		fmt.Println("Error: Try to read a variable that hasn't been added. var: " + key)
	}
	return 0
}

func (scope_ Scope_) Back() *Scope_ {
	return scope_.last
}

/* function scope system ( only for functions ) */
type FuncScope_ struct {
	last   *FuncScope_
	father *FuncScope_
	vars   map[string]Func_
}

func MakeFuncScope(father *FuncScope_, last *FuncScope_) *FuncScope_ {
	return &FuncScope_{last, father, make(map[string]Func_)}
}

func (funcScope_ FuncScope_) Add(key string, value Func_) {
	_, ok := funcScope_.vars[key]
	if ok {
		fmt.Println("Error: Try to define a function that has been defined. func: " + key)
	} else {
		funcScope_.vars[key] = value
	}
}

func (funcScope_ FuncScope_) Find(key string) Func_ {
	v, ok := funcScope_.vars[key]
	if ok {
		return v
	} else if funcScope_.father != nil {
		return funcScope_.father.Find(key)
	} else {
		fmt.Println("Error: Try to use a function that hasn't been defined. func: " + key)
		return Func_{nil, nil, []string{}, []Ast{}}
	}
}

func (funcScope_ FuncScope_) Back() *FuncScope_ {
	return funcScope_.last
}

/* Queue_ ( for tmpQueue ) */

// Give a proper beginning size of a queue
const maxBufferSize = 5

type Queue_ struct {
	father  *Queue_
	data    []interface{}
	dataLen int
}

func MakeTmpQueue(father *Queue_) *Queue_ {
	return &Queue_{father, make([]interface{}, maxBufferSize), 0}
}

func (queue_ *Queue_) Add(value interface{}) {
	queue_.dataLen++
	if queue_.dataLen > maxBufferSize {
		queue_.data = append(queue_.data, value)
	} else {
		queue_.data[queue_.dataLen-1] = value
	}
}

func (queue_ Queue_) Get() interface{} {
	return queue_.data[0]
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
