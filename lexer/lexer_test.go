package lexer

import (
	"testing"
	"emo/token"
)

func TestNextToken(t *testing.T) {
	input := `
	let cheers = 9 / 3;
	let candles = 4 * 4;
	
	let add = f(x, y) {
		return x + y;
	}

	let eq = f(x, y) {
		let a = x > y;
		let b = x < y;
		if (a) {
			return true;
		} else {
			return false;
		}
	}

	let test = !eq(cheers, candles);
	let pilots = add(add(cheers, candles));

	test == pilots
	test != pilots
	`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "cheers"},
		{token.ASSIGN, "="},
		{token.INT, "9"},
		{token.FWDSLASH, "/"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "candles"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.ASTRX, "*"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "eq"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.RBRACE, "}"},

		{token.LET, "let"},
		{token.IDENT, "test"},
		{token.ASSIGN, "="},
		{token.BANG, "!"},
		{token.IDENT, "eq"},
		{token.LPAREN, "("},
		{token.IDENT, "cheers"},
		{token.COMMA, ","},
		{token.IDENT, "candles"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "pilots"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "cheers"},
		{token.COMMA, ","},
		{token.IDENT, "candles"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "test"},
		{token.EQ, "=="},
		{token.IDENT, "pilots"},
		
		{token.IDENT, "test"},
		{token.NOTEQ, "!="},
		{token.IDENT, "pilots"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
			i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
			i, tt.expectedLiteral, tok.Literal)
		}
	}
}