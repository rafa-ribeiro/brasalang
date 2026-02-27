package lexer

import (
	"testing"

	"github.com/rafa-ribeiro/brasalang/internal/token"
)

func TestTokensBasicProgram(t *testing.T) {
	l := New("x int = 10\n{\n  true\n}\n")
	got := l.Tokens()

	wantTypes := []token.Type{
		token.IDENT,
		token.IDENT,
		token.EQUAL,
		token.INT,
		token.NEWLINE,
		token.LBRACE,
		token.NEWLINE,
		token.TRUE,
		token.NEWLINE,
		token.RBRACE,
		token.NEWLINE,
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
