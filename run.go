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
		errlog.Err("parser", "unexpected 'end'.")
	}
	return codes
}

/*
RunBlock(string) receives a text script ( string ) and run it directly.
The runtime status would not be saved as the script runs in a Block_.
*/
func RunBlock(str string) {
	Block_{Build(str)}.run()
}

/*
RunModule(string) receives a text script ( string ) and run it directly.
The interpreter will add a structure named 'globalName' in global scope, which contains all the variables the module defined.
*/
func RunModule(str string, moduleName string) {
	Scope.Add(moduleName, Struct_{Build(str)}.run())
}
