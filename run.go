package parser

import (
	"strings"

	errlog "github.com/zhangzheheng12345/flowscript/error_logger"
	"github.com/zhangzheheng12345/flowscript/lexer"
)

func Build(str string) []Ast {
	tokens := lexer.Lex(strings.Split(str, ""))
	codes, index := Parse(tokens)
	if index < len(tokens) {
		errlog.Err("parser", tokens[len(tokens)-1].Line(), "unexpected 'end'.")
	}
	return codes
}

/*
RunBlock(string) receives a text script ( string ) and run it directly.
The runtime status would not be saved as the script runs in a Block_.
*/
func RunBlock(str string) {
	Block_{Build(str), 0}.run()
}

/*
RunModule(string) receives a text script (string) and run it directly.
The interpreter will add a structure named 'globalName' in global scope, which contains all the variables the module defined.
*/
func RunModule(str string, moduleName string) {
	Scope.Add(moduleName, Struct_{Build(str), 0}.run())
}

/*
GoFunc wraps native Go functions to enable you to call Go functions inside FlowScript.
*/
type GoFunc struct {
	fn       func([]interface{}) interface{}
	argnum   int
	foreArgs []interface{} // for curried
}

func (goFunc GoFunc) run(args []interface{}) interface{} {
	if goFunc.argnum < 0 {
		if -goFunc.argnum > len(args) {
			return GoFunc{goFunc.fn, goFunc.argnum + len(args), args}
		}
	} else if goFunc.argnum < len(args) {
		errlog.Err("runtime", errlog.Line, "Too many arguments while calling function.")
		return 0
	} else if goFunc.argnum > len(args) {
		return GoFunc{goFunc.fn, goFunc.argnum - len(args), args}
	}
	return goFunc.fn(append(goFunc.foreArgs, args...))
}

/*
Add a native Go function to global scope. The function will be named as `name`
If argnum is smaller than 0, it means the function needs at least -argnum arguments
*/
func AddGoFunc(name string, fn func([]interface{}) interface{}, argnum int) {
	Scope.Add(name, GoFunc{fn, argnum, make([]interface{}, 0)})
}

func init() {
	AddBuildinFuncs()
}
