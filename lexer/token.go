package lexer

/* token kinds */
const (
	/* Basic tokens */
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
	XOR
	EXPR
	VAR
	IF
	ELSE
	SEND
	DEF
	LAMBDA
	STRUCT
	BEGIN
	END
	XEXP
	STOP
	NUM
	SYMBOL
	STR
	CHAR
	/* Buildin commands */
	ECHO
	INPUT
	LIST
	LEN
	INDEX
	APP
	SLICE
)

/* token structure, including kind and value */
type Token struct {
	kind  byte
	value string
}

func (token Token) Type() byte {
	return token.kind
}

func (token Token) Value() string {
	return token.value
}
