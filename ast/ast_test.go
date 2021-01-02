package ast

import (
	"emo/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	expected := `let myVar = anotherVar;`
	result := program.String()
	if result != expected {
		t.Errorf("program.String returns incorrect value.\nExpected %s\ngot %s\n", expected, result) 
	}
}