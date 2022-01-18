package xlexer

import (
	"strings"
	"testing"
)

func TestXLex(t *testing.T) {
	const t1 = "7--1"
	var res1 = []Token{
		Token{NUM, "7"},
		Token{SUB, ""},
		Token{NUM, "-1"},
	}
	result := Lex(strings.Split(t1, ""))
	for key, value := range result {
		if res1[key].Type() != value.Type() ||
			res1[key].Value() != value.Value() {
			t.Fatal("expecting: ", res1[key], "but got: ", value, "token: ", key)
		}
	}
}
