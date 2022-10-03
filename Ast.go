package parser

type Ast interface {
	run(*Context) interface{}
}

type Var_ struct {
	name string
	op   Value
	line int
}

func (var_ Var_) run(ctx *Context) interface{} {
	ctx.Line = var_.line
	if var_.op != nil {
		/* give a start value */
		ctx.scope.Add(var_.name, var_.op.get(ctx), ctx)
	} else {
		/* default 0 */
		ctx.scope.Add(var_.name, 0, ctx)
	}
	return 0
}

type Fill_ struct {
	fn   Ast // A expression which returns func
	op   Value
	line int
}

func (fill_ Fill_) run(ctx *Context) interface{} {
	ctx.Line = fill_.line
	var fn Func_
	if fill_.fn == nil {
		fn = WantFunc(fill_.op.get(ctx), ctx)
	} else {
		fn = WantFunc(fill_.fn.run(ctx), ctx)
	}
	var argsLen int
	if tmpQueue.Size() < fn.argsNum() || fn.argsNum() < 0 {
		argsLen = tmpQueue.Size()
	} else {
		argsLen = fn.argsNum()
	}
	args := make([]interface{}, argsLen)
	for i := 0; i < argsLen; i++ {
		args[i] = tmpQueue.Get(ctx)
		tmpQueue.Pop()
	}
	return fn.run(args, ctx)
}

type Enum_ struct {
	names []string
	line  int
}

func (enum_ Enum_) run(ctx *Context) interface{} {
	ctx.Line = enum_.line
	for _, v := range enum_.names {
		ctx.scope.Add(v, ctx.scope.enumCounter, ctx)
		ctx.scope.enumCounter++
	}
	return 0
}

type Def_ struct {
	name  string
	args  []string
	codes []Ast
	line  int
}

func (def_ Def_) run(ctx *Context) interface{} {
	ctx.Line = def_.line
	res := FlowFunc{ctx.scope, def_.args, def_.codes}
	ctx.scope.Add(def_.name, res, ctx)
	return res
}

type Lambda_ struct {
	args  []string
	codes []Ast
	line  int
}

func (lambda_ Lambda_) run(ctx *Context) interface{} {
	ctx.Line = lambda_.line
	return FlowFunc{ctx.scope, lambda_.args, lambda_.codes}
}

/* Struct_ ( with one underline ) is the command */
type Struct_ struct {
	codes []Ast
	line  int
}

func (struct_ Struct_) run(ctx *Context) interface{} {
	ctx.Line = struct_.line
	ctx.scope = MakeScope(ctx.scope, ctx.scope)
	for _, code := range struct_.codes {
		code.run(ctx)
	}
	result := ctx.scope.vars
	ctx.scope = ctx.scope.Back()
	return Struct{result}
}

/* a send sequence */
type Send_ struct {
	codes []Ast
	line  int
}

func (send_ Send_) run(ctx *Context) interface{} {
	ctx.Line = send_.line
	tmpQueue = MakeTmpQueue(tmpQueue, len(send_.codes))
	for _, code := range send_.codes {
		tmpQueue.Add(code.run(ctx))
	}
	result := tmpQueue.Get(ctx)
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

func (block_ Block_) run(ctx *Context) interface{} {
	ctx.Line = block_.line
	ctx.scope = MakeScope(ctx.scope, ctx.scope)
	var result interface{}
	for _, code := range block_.codes {
		result = code.run(ctx)
	}
	ctx.scope = ctx.scope.Back()
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

func (if_ If_) run(ctx *Context) interface{} {
	ctx.Line = if_.line
	if WantInt(if_.condition.get(ctx), ctx) != 0 {
		ctx.scope = MakeScope(ctx.scope, ctx.scope)
		var result interface{}
		for _, code := range if_.ifcodes {
			result = code.run(ctx)
		}
		ctx.scope = ctx.scope.Back()
		return result
	} else {
		ctx.scope = MakeScope(ctx.scope, ctx.scope)
		var result interface{} = 0
		for _, code := range if_.elsecodes {
			result = code.run(ctx)
		}
		ctx.scope = ctx.scope.Back()
		return result
	}
}

/* call a user defined function*/
type Call_ struct {
	name Value
	args []Value
	line int
}

func (call_ Call_) run(ctx *Context) interface{} {
	ctx.Line = call_.line
	argsValue := make([]interface{}, 0)
	for _, arg := range call_.args {
		argsValue = append(argsValue, arg.get(ctx))
	}
	return WantFunc(call_.name.get(ctx), ctx).run(argsValue, ctx)
}
