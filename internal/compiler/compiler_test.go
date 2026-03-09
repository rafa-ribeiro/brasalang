package compiler

import (
	"strings"
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

func TestCompileAndRunFunctionCallWithLocalsAndGlobals(t *testing.T) {
	srcCode := `
	base int = 10
	def sum_with_base(a int, b int) -> int {
		tmp int = a + b
		return tmp + base
	}
	sum_with_base(1, 2)
	`

	result := compileAndRun(t, srcCode)
	if result.Kind != value.IntKind || result.I != 13 {
		t.Fatalf("unexpected result: got=%v", result)
	}
}

func TestCompileVoidFunctionReturnNil(t *testing.T) {
	src := `
	def noop() {
		return
	}
	noop()
	`
	result := compileAndRun(t, src)
	if result.Kind != value.NilKind {
		t.Fatalf("expected nil result, got=%v", result)
	}
}

func TestCompileRejectsValueReturnInVoidFunction(t *testing.T) {
	src := `
	def noop() {
		return 1
	}
	`
	p := parser.NewFromSource(src)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	_, err := New().Compile(program)
	if err == nil || !strings.Contains(err.Error(), "cannot return a value") {
		t.Fatalf("expected compile error for returning value in void function, got=%v", err)
	}
}

func TestCompileAndRunTupleReturn(t *testing.T) {
	src := `
	def pair(a int, b int) -> (int, int) {
		return a, b
	}
	pair(4, 5)
	`
	result := compileAndRun(t, src)
	if result.Kind != value.TupleKind || len(result.Items) != 2 {
		t.Fatalf("expected tuple result, got=%v", result)
	}
	if result.Items[0].I != 4 || result.Items[1].I != 5 {
		t.Fatalf("unexpected tuple items: %v", result.Items)
	}
}

func TestRuntimeErrorWhenFunctionWithReturnTypeOmitsReturn(t *testing.T) {
	src := `
	def bad(a int) -> int {
		a + 1
	}
	bad(1)
	`
	p := parser.NewFromSource(src)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	chunk, err := New().Compile(program)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}

	machine := vm.New()
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("expected runtime panic")
		}
		if !strings.Contains(recovered.(string), "without explicit return") {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	machine.Run(chunk)
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
