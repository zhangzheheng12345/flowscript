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
    args []string
    codes []Ast
}

/*
Func_.run([]int) run the codes in the function.
*/
func (func_ Func_) run(args []int) int {
    /* provide a independence scope for the function*/
    Scope = MakeScope(Scope)
    FuncScope = MakeFuncScope(FuncScope)
    /* add the arguments to the local scope*/
    for key, arg := range args {
        Scope.Add(func_.args[key],arg)
    }
    /* run the codes. 
    the last expression will give the return value. 
    the default return value is 0. 
    */
    result := 0
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
var Scope *Scope_ = &Scope_{nil, make(map[string]int)}
var FuncScope *FuncScope_ = &FuncScope_{nil, make(map[string]int)}

/*
exp1 > exp2 > exp3
The values of the send sequence expressions will be pushed to this queue.
*/
var tmpQueue Queue_

/* scope system ( only for variables ) */
type Scope_ struct {
    father *Scope_
    vars map[string]int
}

func MakeScope(father *Scope_) *Scope_ {
    return &Scope_{father,make(map[string]int)}
}

func (scope_ Scope_) Add(key string, value int) {
    _, ok := scope_.vars[key]
    if ok {
        fmt.Println("Error: Try to add a variable that has been added. var: " + key)
    } else {
        scope_.vars[key] = value
    }
}

func (scope_ Scope_) Find(key string) int {
    v, ok := scope_.vars[key]
    if ok {
        return v
    } else if scope_.father != nil{
        return scope_.father.Find(key)
    } else {
        fmt.Println("Error: Try to read a variable that hasn't been added. var: " + key)
    }
    return 0
}

func (scope_ Scope_) Back() *Scope_ {
    return scope_.father
}

/* function scope system ( only for functions ) */
type FuncScope_ struct {
    father *FuncScope_
    vars map[string]int
}

func MakeFuncScope(father *FuncScope_) *FuncScope_ {
    return &FuncScope_{father,make(map[string]int)}
}

func (funcScope_ FuncScope_) Add(key string, value int) {
    _, ok := funcScope_.vars[key]
    if ok {
        fmt.Println("Error: Try to define a function that has been defined. func: " + key)
    } else {
        scope_.vars[key] = value
    }
}

func (funcScope_ FuncScope_) Find(key string) int {
    v, ok := funcScope_.vars[key]
    if ok {
        return v
    } else if funcScope_.father != nil{
        return funcScope_.father.Find(key)
    } else {
        fmt.Println("Error: Try to use a function that hasn't been defined. func: " + key)
    }
    return 0
}

func (funcScope_ FuncScope_) Back() *FuncScope_ {
    return funcScope_.father
}

/* Queue_ ( for tmpQueue ) */
type Queue_ struct {
    data []int
}

func (queue_ Queue_) Add(value int) {
    queue_.data = append(queue_.data,value)
}

func (queue_ Queue) Get() int {
    return queue_.data[0]
}

func (queue_ Queue) Pop() {
    queue_.data = queue_.data[1:]
}

func (queue_ Queue_) Clear() {
    queue_.data = make([]int)
}