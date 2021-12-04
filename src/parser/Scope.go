package parser

import (
    "fmt"
)

/* function object */
type Func_ struct {
    args []string
    codes []Ast
}

/* function running */
func (func_ Func_) run(args []int) int {
    Scope = MakeScope(Scope)
    FuncScope = MakeFuncScope(FuncScope)
    for key, arg := range args {
        Scope.Add(func_.args[key],arg)
    }
    result := 0
    for _, code := range func_.codes {
        result = code.run()
    }
    Scope = Scope.Back()
    FuncScope = FuncScope.Back()
    return result
}

/* scopes */
var Scope *Scope_ = &Scope_{nil, make(map[string]int)}
var FuncScope *FuncScope_ = &FuncScope_{nil, make(map[string]int)}

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