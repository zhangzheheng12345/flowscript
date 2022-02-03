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
	VAR
	IF
	SEND
	DEF
	BEGIN
	END
	XEXP
	STOP
	NUM
	SYMBOL
	STR
	CHAR
	TMP
	/* Buildin commands */
	ECHO
	INPUT
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
