package parser

import (
    "testing"
    "github.com/zhangzheheng12345/FlowScript/xlexer"
    "strings"
    "fmt"
)

func TestExp(t *testing.T) {
    table := []struct{ input string; output int }{
        {"1+2",3},
        {"2*3",6},
        {"3+5*6",33},
        {"3*2+4-1",9},
        {"4+6*(3-1)-6/3",14},
        {"(6+4)-6*(7/1)",0-32},
        {"3+4*5-6-7",10},
    }
    for k, v := range table {
        fmt.Println(k)
        value := Exp_{xlexer.Lex(strings.Split(v.input,""))}.get()
        if (value != v.output) {
            fmt.Println(v.input,"expecting",v.output,"but got",value)
        }
    }
}