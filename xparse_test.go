package parser

import (
    "testing"
    "xlexer"
    "strings"
    "fmt"
)

func TestExp(t *testing.T) {
    table := []struct{ input string, output int }{
        {"3*2+4-1",9},
        {"4+6*(3-1)-6/3",14},
        {"(6+4)-6*(7/1)",0-32},
    }
    for _, v := range table {
        if Exp_{xlexer.Lex(strings.Split(v.input,""))}.get() != v.output {
            fmt.Println(v.input,"expecting",v.ouput,"but not")
        }
    }
}