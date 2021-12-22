package xlexer

const (
    ADD = iota
    SUB
    MULTI
    DIV
    FPAREN
    BPAREN
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