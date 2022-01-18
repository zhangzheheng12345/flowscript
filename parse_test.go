package parser

import (
	"testing"
)

func TestRunCode1(t *testing.T) {
	const ts = "echo (7--1)"
	RunCode(ts)
}
