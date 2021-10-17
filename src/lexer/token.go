package lexer

const (
    ADD = iota
    SUB
    MULTI
    DIV
    EQUAL
    BIGR
    SMLR
    NOT
    AND
    OR
    VAR
    SET
    IF
    WHILE
    SEND
    DEF
    BEGIN
    END
    STOP
    NUM
    SYMBOL
    TMP
)

type Token struct {
    kind byte
    value string
}

func (token Token) Type() byte {
    return token.kind
}

func (token Token) Value() string {
    return token.value
}