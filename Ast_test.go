package parser

import (
    "testing"
    "fmt"
)

func TestFunc(t *testing.T) {
    fmt.Println("Ast:t1-----")
    FuncScope.Add("func",Func_{
        Scope,FuncScope,
        []string{
            "a",
            "b",
        },
        []Ast{
            Echo_{Symbol_{"a"}},
            Echo_{Symbol_{"b"}},
            Send_{[]Ast{
                    Add_{Symbol_{"a"}, Symbol_{"b"}},
                    Echo_{Tmp_{}},
            }},
        },
    })
    FuncScope.Find("func").run([]int{1,2})
    FuncScope.Find("func").run([]int{100,892})
    fmt.Println("Ast:t1-----")
}