package lexer

import (
    "testing"
    "fmt"
)

func TestLex(t *testing.T) {
    const t1 = 
        "add 1 var1 > var var2" +
    const res1 = []Token{
        Token{ADD, ""},
        Token{NUM, "1"},
        Token{SYMBOL, "var1"},
        Token{TMP, ""}
        Token{VAR, ""}
        Token{SYMBOL, "var2"}
    }
    result := Lex(t1)
    for key, value := range result {
        if res1[key].Type() != value.Type() ||
           res1[key].Value() != value.Value() {
           fmt.Println("expecting: ",res1[key])
           fmt.Println("but got: ",value)
        }
    }
}