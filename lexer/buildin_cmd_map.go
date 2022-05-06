package lexer

/* Buildin cmd to token type. This map should not be changed */
var BuildinCmdMap = map[string]byte{
	"equal":  EQUAL,
	"bigr":   EQUAL,
	"smlr":   SMLR,
	"not":    NOT,
	"and":    AND,
	"or":     OR,
	"xor":    XOR,
	"expr":   EXPR,
	"var":    VAR,
	"if":     IF,
	"else":   ELSE,
	"send":   SEND,
	"def":    DEF,
	"lambda": LAMBDA,
	"struct": STRUCT,
	"begin":  BEGIN,
	"end":    END,
	"type":   TYPE,
	"echo":   ECHO,
	"echoln": ECHOLN,
	"input":  INPUT,
	"list":   LIST,
	"len":    LEN,
	"index":  INDEX,
	"app":    APP,
	"slice":  SLICE,
	"words":  WORDS,
	"lines":  LINES,
}
