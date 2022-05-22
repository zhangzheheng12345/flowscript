package xlexer

const (
	ADD = iota
	SUB
	MULTI
	DIV
	MOD
	AND
	OR
	XOR
	FPAREN
	BPAREN
	NUM
	SYMBOL
	TMP
)

type Token struct {
	kind  byte
	value string
	line  int
}

func (token Token) Type() byte {
	return token.kind
}

func (token Token) Value() string {
	return token.value
}

func (token Token) Line() int {
	return token.line
}
