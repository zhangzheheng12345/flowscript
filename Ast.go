package parser

import (
    "fmt"
)

/* Tool function, for cannot convert bool to int directly */
func booltoint(value bool) int {
    if value {
        return 1
    }
    return 0
}

type Ast interface {
    run() int
}

type Add_ struct {
    op1 Value
    op2 Value
}

func (add_ Add_) run() int {
    return add_.op1.get() + add_.op2.get()
}

type Sub_ struct {
    op1 Value
    op2 Value
}

func (sub_ Sub_) run() int {
    return sub_.op1.get() + sub_.op2.get()
}

type Multi_ struct {
    op1 Value
    op2 Value
}

func (multi_ Multi_) run() int {
    return multi_.op1.get() * multi_.op2.get()
}

type Div_ struct {
    op1 Value
    op2 Value
}

func (div_ Div_) run() int {
    return div_.op1.get() + div_.op2.get()
}

type Mod_ struct {
    op1 Value
    op2 Value
}

func (mod_ Mod_) run() int {
    return mod_.op1.get() % mod_.op2.get()
}

type Bigr_ struct {
    op1 Value
    op2 Value
}

func (bigr_ Bigr_) run() int {
    return booltoint(bigr_.op1.get() > bigr_.op2.get())
}

type Smlr_ struct {
    op1 Value
    op2 Value
}

func (smlr_ Smlr_) run() int {
    return booltoint(smlr_.op1.get() < smlr_.op2.get())
}

type Equal_ struct {
    op1 Value
    op2 Value
}

func (equal_ Equal_) run() int {
    return booltoint(equal_.op1.get() == equal_.op2.get())
}

type And_ struct {
    op1 Value
    op2 Value
}

func (and_ And_) run() int {
    return booltoint((and_.op1.get() != 0) && (and_.op2.get() != 0))
}

type Or_ struct {
    op1 Value
    op2 Value
}

func (or_ Or_) run() int {
    return booltoint((or_.op1.get() != 0) || (or_.op2.get() != 0))
}

type Not_ struct {
    op Value
}

func (not_ Not_) run() int {
    return booltoint(not_.op.get() == 0)
}

type Var_ struct {
    name string
    op Value
}

func (var_ Var_) run() int {
    if var_.op != nil {
        /* give a start value */
        Scope.Add(var_.name,var_.op.get())
    } else {
        if tmpQueue.Size() == 0 {
            /* default 0 */
            Scope.Add(var_.name,0)
        } else {
            /* autoly pick value from send queue */
            Scope.Add(var_.name,tmpQueue.Get())
            tmpQueue.Pop()
        }
    }
    return 0
}

type Def_ struct {
    name string
    args []string
    codes []Ast
}

func (def_ Def_) run() int {
    FuncScope.Add(def_.name,Func_{Scope,FuncScope,def_.args,def_.codes})
    return 0
}

/* a send sequence */
type Send_ struct {
    codes []Ast
}

func (send_ Send_) run() int {
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

func (block_ Block_) run() int {
    Scope = MakeScope(Scope,Scope)
    FuncScope = MakeFuncScope(FuncScope,FuncScope)
    result := 0
    for _, code := range block_.codes {
        result = code.run()
    }
    Scope = Scope.Back()
    FuncScope = FuncScope.Back()
    return result
}

/*
if condition begin
    ...
end
*/
type If_ struct {
    condition Value
    codes []Ast
}

func (if_ If_) run() int {
    if if_.condition.get() != 0 {
        Scope = MakeScope(Scope,Scope)
        FuncScope = MakeFuncScope(FuncScope,FuncScope)
        result := 1
        for _, code := range if_.codes {
            result = code.run()
        }
        Scope = Scope.Back()
        FuncScope = FuncScope.Back()
        return result
    }
    return 0
}

/* call a user defined function*/
type Call_ struct {
    name string
    args []Value
}

func (call_ Call_) run() int {
    argsValue := make([]int,0)
    for _, arg := range call_.args {
        argsValue = append(argsValue,arg.get())
    }
    return FuncScope.Find(call_.name).run(argsValue)
}

/* output to stdout */
type Echo_ struct {
    op Value
}

func (echo_ Echo_) run() int {
    fmt.Println(echo_.op.get())
    return 0
}