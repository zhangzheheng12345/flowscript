package parser

import (
    "fmt"
    "xlexer"
)

func E_(tokens []xlexer.Token, value int) int {
    switch token[0].Type() {
        case xlexer.ADD:
            return value + E_(T(tokens[1:]))
        case xlexer.SUB:
            return value - E_(T(tokens[1:]))
        default:
            fmt.Println("Error: unexpected token in xparser")
    }
}

func T(tokens []xlexer.Token) ([]xlexer.Token, int) {
    return T_(F(tokens))
}

func T_(tokens []xlexer.Token, value int) ([]xlexer.Token, int) {
    if tokens == nil {
        return nil, value
    }
    switch tokens[0].Type() {
        case ADD, SUB:
            return tokens, value
        case MULTI:
            tail, v := T_(F(tokens))
            return tail, value * v
        case DIV:
            tail, v := T_(F(tokens))
            return tail, value / v
        default:
            fmt.Println("Error: unexpected token in xparser")
    }
}

func F(tokens []xlexer.Token) ([]xlexer.Token, int) {
    if tokens == nil {
        return nil, 1
    }
    switch tokens[0].Type() {
        case xlexer.SYMBOL:
            return tokens[1:], Scope.Find(tokens[0].Value())
        case xlexer.TMP:
            result := tmpQueue.Get()
            tmpQueue.Pop()
            return tokens[1:], result
        case xlexer.NUM:
            return tokens[1:], strconv.Atoi(tokens[0].Value())
        case xlexer.FPAREN:
            count := 1
            index := 1
            for count != 0 && index < len(tokens) {
                if tokens[index].Type() == xlexer.FPAREN {
                    count++
                } else if tokens[index].Type() == xlexer.BPAREN {
                    count--
                }
                index++
            }
            if index >= len(tokens) {
                fmt.Println("Error: lose back parenthesis in xparser")
            }
            return tokens[index:], Exp_{ tokens[1:index] }.get()
        default:
            fmt.Println("Error: X expression mistake")
    }
}