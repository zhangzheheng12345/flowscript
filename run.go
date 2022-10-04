package parser

import (
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/lexer"
)

/*
Provide a independent context to run codes.
*/
type Context struct {
	scope *Scope_
	Line  int
}

/*
Build(str string) receives a string of FlowScript code and build it to runnable AST.
*/
func Build(str string) []Ast {
	tokens := lexer.Lex(strings.Split(str, ""))
	codes, index := Parse(tokens)
	if index < len(tokens) {
		errlog.Err("parser", tokens[len(tokens)-1].Line(), "unexpected 'end'.")
	}
	return codes
}

/*
RunBlock(string) receives a text script ( string ) and run it in a seperate env.
The runtime status would not be saved as the script runs in a Block_.
*/
func (ctx *Context) RunBlock(str string) interface{} {
	return Block_{Build(str), 0}.run(&Context{MakeScope(ctx.scope, nil), 1})
}

/*
GoFunc wraps native Go functions to enable you to call Go functions inside FlowScript.
*/
type GoFuncForFS func([]interface{}, *Context) interface{}

type GoFunc struct {
	fn       GoFuncForFS
	argnum   int
	foreArgs []interface{} // for curried
}

func (goFunc GoFunc) run(args []interface{}, ctx *Context) interface{} {
	if goFunc.argnum < 0 {
		// argnum < 0, argnum = -n => at least n-1 arguments
		if -goFunc.argnum-1 > len(args) {
			return GoFunc{goFunc.fn, goFunc.argnum + len(args), args}
		}
	} else if goFunc.argnum < len(args) {
		errlog.Err("runtime", ctx.Line, "Too many arguments while calling function.")
		return 0
	} else if goFunc.argnum > len(args) {
		return GoFunc{goFunc.fn, goFunc.argnum - len(args), args}
	}
	return goFunc.fn(append(goFunc.foreArgs, args...), ctx) // call & return
}

func (goFunc GoFunc) argsNum() int {
	return goFunc.argnum
}

/*
Add a native Go function to global scope. The function will be named as `name`
If argnum is smaller than 0, it means the function needs at least -argnum arguments
*/
func AddGoFunc(name string, fn GoFuncForFS, argnum int) {
	Global.scope.Add(name, GoFunc{fn, argnum, make([]interface{}, 0)}, Global)
}

func init() {
	AddBuildinFuncs()
}
