package lexer

/* token kinds */
const (
	/* Basic tokens */
	EQUAL = iota
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
	TYPE
	ECHO
	ECHOLN
	INPUT
	LIST
	LEN
	INDEX
	APP
	SLICE
	WORDS
	LINES
)

/* token structure, including kind and value */
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
