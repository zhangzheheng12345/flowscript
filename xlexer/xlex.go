package xlexer

import (
    "fmt"
    "strings"
)

/* Check if a character is a number. */
func IsSingleNum(str string) bool {
    return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 57
}

/* Check if a character is a alpha */
func IsSingleAlpha(str string) bool {
    return len(str) == 1 && []byte(str)[0] >= 65 && []byte(str)[0] <= 122
}

func Lex(str []string) []Token {
    result := make([]Token,0)
    for index := 0; index < len(str); index++ {
        if str[index] == " " || str[index] == "\t" {
            // Do nothing
        } else if str[index] == "_" || IsSingleAlpha(str[index]) {
            begin := index
            for ;index < len(str) &&
                ( IsSingleAlpha(str[index]) ||
                IsSingleNum(str[index]) ||
                str[index] == "_" );index++ {}
            result = append(result,Token{SYMBOL, strings.Join(str[begin:index], "")})
            index--
        } else if IsSingleNum(str[index]) {
            begin := index
            for ;index < len(str) && IsSingleNum(str[index]);index++ {}
            result = append(result,Token{NUM, strings.Join(str[begin:index], "")})
            index--
        } else if str[index] == "(" {
            result = append(result,Token{FPAREN,""})
        }  else if str[index] == ")" {
            result = append(result,Token{BPAREN,""})
        } else if str[index] == "+" {
            result = append(result,Token{ADD,""})
        } else if str[index] == "-" {
            if result[len(result) - 1].Type() == NUM ||
                result[len(result) - 1].Type() == SYMBOL ||
                result[len(result) - 1].Type() == TMP ||
                result[len(result) - 1].Type() == BPAREN {
                result = append(result,Token{SUB,""})
            } else {
                result = append(result,Token{TMP,""})
            }
        } else if str[index] == "*" {
            result = append(result,Token{MULTI,""})
        } else if str[index] == "/" {
            result = append(result,Token{DIV,""})
        } else {
            fmt.Println("Warn: unknown token in X expression. Token: ", index)
        }
    }
    return result
}