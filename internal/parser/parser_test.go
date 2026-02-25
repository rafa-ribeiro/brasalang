package parser

import "testing"

func TestParseProgramExpressionStatements(t *testing.T) {
	p := NewFromSource("true; 42; name;")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	if got, want := len(program.Statements), 3; got != want {
		t.Fatalf("statements count mismatch: got=%d want=%d", got, want)
	}
}
