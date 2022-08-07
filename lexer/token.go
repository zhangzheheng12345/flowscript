package lexer

/* token kinds */
const (
	/* Basic tokens */
	VAR = iota
	FILL
	ENUM
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

/* Buildin cmd to token type. This map should not be changed */
var BuildinCmdMap = map[string]byte{
	"var":    VAR,
	"fill":   FILL,
	"if":     IF,
	"else":   ELSE,
	"send":   SEND,
	"def":    DEF,
	"lambda": LAMBDA,
	"struct": STRUCT,
	"begin":  BEGIN,
	"end":    END,
	"enum":   ENUM,
}
