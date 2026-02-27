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

func TestCompileKeepsOnlyLastStatementResult(t *testing.T) {
	result := compileAndRun(t, "1 + 1\n10 - 3\n")
	if result.Kind != value.IntKind || result.I != 7 {
		t.Fatalf("unexpected result: got=%v", result)
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
