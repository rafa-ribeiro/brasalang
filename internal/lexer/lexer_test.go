package lexer

import (
	"testing"

	"github.com/rafa-ribeiro/brasalang/internal/token"
)

func TestTokensBasicProgram(t *testing.T) {
	l := New("if (true) { 10; } else { 20; }")
	got := l.Tokens()

	wantTypes := []token.Type{
		token.IF,
		token.LPAREN,
		token.TRUE,
		token.RPAREN,
		token.LBRACE,
		token.INT,
		token.SEMICOLON,
		token.RBRACE,
		token.ELSE,
		token.LBRACE,
		token.INT,
		token.SEMICOLON,
		token.RBRACE,
		token.EOF,
	}

	if len(got) != len(wantTypes) {
		t.Fatalf("token count mismatch: got=%d want=%d", len(got), len(wantTypes))
	}

	for i, want := range wantTypes {
		if got[i].Type != want {
			t.Fatalf("token[%d] = %s, want %s", i, got[i].Type, want)
		}
	}
}
