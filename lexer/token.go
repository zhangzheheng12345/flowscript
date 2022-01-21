package lexer

/* token kinds */
const (
    ADD = iota
    SUB
    MULTI
    DIV
    MOD
    EQUAL
    BIGR
    SMLR
    NOT
    AND
    OR
    VAR
    IF
    SEND
    DEF
    BEGIN
    END
    XEXP
    STOP
    ECHO
    NUM
    SYMBOL
    STR
    CHAR
    TMP
)

/* token structure, including kind and value */
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