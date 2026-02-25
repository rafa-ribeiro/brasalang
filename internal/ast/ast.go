package ast

import (
	"github.com/rafa-ribeiro/brasalang/internal/token"
)

type Node interface {
	Pos() token.Position
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Program struct {
	Statements []Stmt
}

type IntLiteral struct {
	Token token.Token
	Value int64
}

func (node *IntLiteral) Pos() token.Position {
	return node.Token.Position
}

func (node *IntLiteral) exprNode()

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (node *BoolLiteral) Pos() token.Position {
	return node.Token.Position
}

func (node *BoolLiteral) exprNode()

type Identifier struct {
	Token token.Token
	Name  string
}

func (node *Identifier) Pos() token.Position {
	return node.Token.Position
}

func (node *Identifier) exprNode()

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (node *UnaryExpr) Pos() token.Position {
	return node.Operator.Position
}

func (node *UnaryExpr) exprNode()

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node *BinaryExpr) Pos() token.Position {
	return node.Operator.Position
}

func (node *BinaryExpr) exprNode()

type ExprStmt struct {
	Expression Expr
	Semicolon  token.Token
}

func (node *ExprStmt) Pos() token.Position {
	return node.Expression.Pos()
}
func (node *ExprStmt) stmtNode()
