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
