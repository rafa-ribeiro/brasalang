package compiler

import (
	"testing"

	"github.com/rafa-ribeiro/brasalang/internal/parser"
	"github.com/rafa-ribeiro/brasalang/internal/value"
	"github.com/rafa-ribeiro/brasalang/internal/vm"
)

func TestCompileAndRunArithmeticExpression(t *testing.T) {
	result := compileAndRun(t, "1 + 2 * 3\n")
	if result.Kind != value.IntKind || result.I != 7 {
		t.Fatalf("unexpected result: got=%v", result)
	}
}

func TestCompileAndRunBooleanExpression(t *testing.T) {
	result := compileAndRun(t, "!(1 + 2 <= 3) && true\n")
	if result.Kind != value.BoolKind || result.B {
		t.Fatalf("unexpected result: got=%v", result)
	}
}

func TestCompileAndRunGlobalVariables(t *testing.T) {
	srcCode := `
	a int = 2 * 5
	b int = 13 - 3
	a == b
	`
	result := compileAndRun(t, srcCode)
	if result.Kind != value.BoolKind || !result.B {
		t.Fatalf("unexpected result: got=%v", result)
	}
}

func TestCompileKeepsOnlyLastStatementResult(t *testing.T) {
	result := compileAndRun(t, "1 + 1\n10 - 3\n")
	if result.Kind != value.IntKind || result.I != 7 {
		t.Fatalf("unexpected result: got=%v", result)
	}
}

func TestCompileFailsForUndeclaredIdentifier(t *testing.T) {
	p := parser.NewFromSource("a == 1\n")
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	_, err := New().Compile(program)
	if err == nil {
		t.Fatalf("expected compile error for undeclared identifier")
	}
}

func compileAndRun(t *testing.T, src string) value.Value {
	t.Helper()

	p := parser.NewFromSource(src)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	c := New()
	chunk, err := c.Compile(program)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}

	machine := vm.New()
	machine.Run(chunk)

	return machine.StackTop()
}
