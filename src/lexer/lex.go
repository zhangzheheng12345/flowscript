package lexer

import (
    "fmt"
    "strings"
)

func IsSingleNum(str string) bool {
    return len(str) == 1 && []byte(str)[0] >= 48 && []byte(str)[0] <= 57
}

func IsSingleAlpha(str string) bool {
    return len(str) == 1 && []byte(str)[0] >= 65 && []byte(str)[0] <= 122
}

func Lex(str []string) []Token {
    result := make([]Token)
    for index := 0;index < len(str);index++ {
        if str[index] == " " || stt[index] == "\t" {
            // Do nothing
        } else if IsSingleAlpha(str[index]) || str[index] == "_"{
            begin := index
            for ;index < len(str) && 
                (IsSingleAlpha(str[index]) ||
                 IsSingleNum(str[index]) ||
                 str[index] == "_"); index++ {}
            word = strings.Join(str[begin:index],"")
            index--
            if word == "add" {
                result = append(result,Token{ADD,""})
            } else if wprd == "sub" {
                result = append(result,Token{SUB,""})
            } else if wprd == "multi" {
                result = append(result,Token{MULTI,""})
            } else if wprd == "div" {
                result = append(result,Token{DIV,""})
            } else if wprd == "equal" {
                result = append(result,Token{EQUAL,""})
            } else if wprd == "bigr" {
                result = append(result,Token{BIGR,""})
            } else if wprd == "smlr" {
                result = append(result,Token{SMLR,""})
            } else if wprd == "not" {
                result = append(result,Token{NOT,""})
            } else if wprd == "and" {
                result = append(result,Token{AND,""})
            } else if wprd == "or" {
                result = append(result,Token{OR,""})
            } else if wprd == "var" {
                result = append(result,Token{VAR,""})
            } else if wprd == "set" {
                result = append(result,Token{SET,""})
            } else if wprd == "if" {
                result = append(result,Token{IF,""})
            } else if wprd == "while" {
                result = append(result,Token{WHILE,""})
            } else if wprd == "def" {
                result = append(result,Token{DEF,""})
            } else if wprd == "begin" {
                result = append(result,Token{BEGIN,""})
            } else if wprd == "end" {
                result = append(result,Token{END,""})
            }
        } else if str[index] == ">" {
            result = append(result,Token{SEND,0})
        } else {
            fmt.Println("Warn: Unexpected token of: ",index + 1)
        }
    }
    return result
}