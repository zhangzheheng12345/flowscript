package xlexer

import (
    "testing"
    "fmt"
    "strings"
)

func TestLex(t *testing.T) {
    const t1 = 
        "1+var1/ (7 - -)*8-7-7"
    var res1 = []Token{
        Token{NUM,"1"},
        Token{ADD,""},
        Token{SYMBOL,"var1"},
        Token{DIV,""},
        Token{FPAREN,""},
        Token{NUM,"7"},
        Token{SUB,""},
        Token{TMP,""},
        Token{BPAREN,""},
        Token{MULTI,""},
        Token{NUM,"8"},
        Token{SUB,""},
        Token{NUM,"7"},
        Token{SUB,""},
        Token{NUM,"7"},
    }
    result := Lex(strings.Split(t1,""))
    for key, value := range result {
        if res1[key].Type() != value.Type() ||
           res1[key].Value() != value.Value() {
            fmt.Println("expecting: ",res1[key])
            fmt.Println("but got: ",value)
             fmt.Println("token: ", key)
        }
    }
}