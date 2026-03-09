package parser

import (
	"testing"

	"github.com/rafa-ribeiro/brasalang/internal/ast"
	"github.com/rafa-ribeiro/brasalang/internal/token"
)

func TestParseProgramExpressionStatements(t *testing.T) {
	p := NewFromSource("true\n42\nname\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	if got, want := len(program.Statements), 3; got != want {
		t.Fatalf("statements count mismatch: got=%d want=%d", got, want)
	}
}

func TestParseTypedVarDeclaration(t *testing.T) {
	p := NewFromSource("idade int = 42\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	decl, ok := program.Statements[0].(*ast.VarDeclStmt)
	if !ok {
		t.Fatalf("expected VarDeclStmt, got %T", program.Statements[0])
	}
	if decl.Name.Lexeme != "idade" || decl.TypeName.Lexeme != "int" {
		t.Fatalf("unexpected declaration: name=%s type=%s", decl.Name.Lexeme, decl.TypeName.Lexeme)
	}
}

func TestParseBlockStatement(t *testing.T) {
	p := NewFromSource("{\n  1 + 2\n  true\n}\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	block, ok := program.Statements[0].(*ast.BlockStmt)
	if !ok {
		t.Fatalf("expected BlockStmt, got %T", program.Statements[0])
	}
	if len(block.Statements) != 2 {
		t.Fatalf("expected 2 statements in block, got %d", len(block.Statements))
	}
}

func TestParseExpressionPrecedence(t *testing.T) {
	p := NewFromSource("1 + 2 * 3 == 7\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	stmt := program.Statements[0].(*ast.ExprStmt)
	equality, ok := stmt.Expression.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected top-level binary expression, got %T", stmt.Expression)
	}
	if equality.Operator.Type != token.EQUAL_EQUAL {
		t.Fatalf("expected top-level operator '==', got %s", equality.Operator.Type)
	}

	leftAdd, ok := equality.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected left side to be binary expression, got %T", equality.Left)
	}
	if leftAdd.Operator.Type != token.PLUS {
		t.Fatalf("expected left operator '+', got %s", leftAdd.Operator.Type)
	}

	rightMul, ok := leftAdd.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected right side of '+' to be binary expression, got %T", leftAdd.Right)
	}
	if rightMul.Operator.Type != token.STAR {
		t.Fatalf("expected right operator '*', got %s", rightMul.Operator.Type)
	}
}

func TestParseFunctionDeclarationAndCall(t *testing.T) {
	p := NewFromSource("def sum(a int, b int) -> int {\n  return a + b\n}\nsum(1, 2)\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	fn, ok := program.Statements[0].(*ast.FuncDeclStmt)
	if !ok {
		t.Fatalf("expected FuncDeclStmt, got %T", program.Statements[0])
	}
	if fn.Name.Lexeme != "sum" || len(fn.Params) != 2 || len(fn.ReturnTypes) != 1 {
		t.Fatalf("unexpected function declaration: %#v", fn)
	}

	stmt := program.Statements[1].(*ast.ExprStmt)
	call, ok := stmt.Expression.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", stmt.Expression)
	}
	if call.Callee.Lexeme != "sum" || len(call.Arguments) != 2 {
		t.Fatalf("unexpected function call: %#v", call)
	}
}

func TestParseFunctionWithoutReturnTypeAndBareReturn(t *testing.T) {
	p := NewFromSource("def noop() {\n  return\n}\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	fn := program.Statements[0].(*ast.FuncDeclStmt)
	if len(fn.ReturnTypes) != 0 {
		t.Fatalf("expected no return types, got %d", len(fn.ReturnTypes))
	}

	ret := fn.Body.Statements[0].(*ast.ReturnStmt)
	if len(ret.Values) != 0 {
		t.Fatalf("expected bare return, got %d values", len(ret.Values))
	}
}

func TestParseMultipleReturnTypes(t *testing.T) {
	p := NewFromSource("def pair(a int) -> (int, bool) {\n  return a, true\n}\n")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	fn := program.Statements[0].(*ast.FuncDeclStmt)
	if got := len(fn.ReturnTypes); got != 2 {
		t.Fatalf("expected 2 return types, got %d", got)
	}
}
