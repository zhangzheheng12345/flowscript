package lexer

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

/* Core lexer code */
func Lex(str []string) []Token {
    result := make([]Token)
    for index := 0;index < len(str);index++ {
        if str[index] == " " || stt[index] == "\t" {
            // Do nothing
        } else if str[index] == "#" {
            for ;index < len(str) && str[index] != "\n"; index++ {
                // Do nothing
            }
            index--
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
            } else if word == "sub" {
                result = append(result,Token{SUB,""})
            } else if word == "multi" {
                result = append(result,Token{MULTI,""})
            } else if word == "div" {
                result = append(result,Token{DIV,""})
            } else if word == "mod" {
                result = append(result,Token{MOD,""})
            } else if word == "equal" {
                result = append(result,Token{EQUAL,""})
            } else if word == "bigr" {
                result = append(result,Token{BIGR,""})
            } else if word == "smlr" {
                result = append(result,Token{SMLR,""})
            } else if word == "not" {
                result = append(result,Token{NOT,""})
            } else if word == "and" {
                result = append(result,Token{AND,""})
            } else if word == "or" {
                result = append(result,Token{OR,""})
            } else if word == "var" {
                result = append(result,Token{VAR,""})
            } else if word == "if" {
                result = append(result,Token{IF,""})
            } else if word == "def" {
                result = append(result,Token{DEF,""})
            } else if word == "begin" {
                result = append(result,Token{BEGIN,""})
            } else if word == "end" {
                result = append(result,Token{END,""})
            } else if word == "echo" {
                result = append(result,Token{ECHO,""})
            } else {
                result = append(result,Token{SYMBOL,word})
            }
        } else if IsSingleNum(str[index]) {
            begin := index
            for ; index < len(str) && IsSingleNum(str[index]); index++ {
                // Do nothing
            }
            result = append(result,Token{NUM,strings.Join(str[begin:index],"")})
            index--
        } else if str[index] == ">" {
            result = append(result,Token{SEND,""})
        } else if str[index] == "-" {
            result = append(result,Token{TMP,""})
        } else if str[index] == ";" ||
                 (str[index] == "\n" &&
                  result[len[result - 1].type() != SEND &&
                  result[len[result - 1].type() != BEGIN) {
            /*
            begin \n ... => not end
            > \n ... => not end
            ; => end, obviously
            or will autoly end
            */
            result = append(result,Token{STOP,""})
        } else if str[index] == "(" {
            begin := index
            count := 1
            for ;index < len(str) && count != 0;index++ {
                if str[index] == "(" {
                    count++
                } else {
                    count--
                }
            }
            result = append(result, Token{XEXP, strings.Join(str[begin:index],"")})
            index--
        } else if str[index] == ")" {
            fmt.Println("Error: too much back parenthesis. Letter: ", index)
        }else {
            fmt.Println("Warn: unexpected token of: ",index + 1)
        }
    }
    return result
}