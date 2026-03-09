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

func (node *IntLiteral) exprNode() {}

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (node *BoolLiteral) Pos() token.Position {
	return node.Token.Position
}

func (node *BoolLiteral) exprNode() {}

type NilLiteral struct {
	Token token.Token
}

func (node *NilLiteral) Pos() token.Position {
	return node.Token.Position
}

func (node *NilLiteral) exprNode() {}

type Identifier struct {
	Token token.Token
	Name  string
}

func (node *Identifier) Pos() token.Position {
	return node.Token.Position
}

func (node *Identifier) exprNode() {}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (node *UnaryExpr) Pos() token.Position {
	return node.Operator.Position
}

func (node *UnaryExpr) exprNode() {}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node *BinaryExpr) Pos() token.Position {
	return node.Operator.Position
}

func (node *BinaryExpr) exprNode() {}

// CallExpr represents a function call
type CallExpr struct {
	Callee    token.Token // the function to be called
	Arguments []Expr
}

func (node *CallExpr) Pos() token.Position {
	return node.Callee.Position
}

func (node *CallExpr) exprNode() {}

type ExprStmt struct {
	Expression Expr
	Semicolon  token.Token
}

func (node *ExprStmt) Pos() token.Position {
	return node.Expression.Pos()
}

func (node *ExprStmt) stmtNode() {}

type VarDeclStmt struct {
	Name        token.Token
	TypeName    token.Token
	Initializer Expr
}

func (node *VarDeclStmt) Pos() token.Position {
	return node.Name.Position
}

func (node *VarDeclStmt) stmtNode() {}

type BlockStmt struct {
	LBrace     token.Token
	Statements []Stmt
}

func (node *BlockStmt) Pos() token.Position {
	return node.LBrace.Position
}

func (node *BlockStmt) stmtNode() {}

// ReturnStmt represents the keyword return used in functions
type ReturnStmt struct {
	Return token.Token
	Values []Expr
}

func (node *ReturnStmt) Pos() token.Position {
	return node.Return.Position
}

func (node *ReturnStmt) stmtNode() {}

// Param defines a function parameter
type Param struct {
	Name token.Token
	Type token.Token
}

// FuncDeclStmt represents a function definition
type FuncDeclStmt struct {
	DefToken    token.Token
	Name        token.Token
	Params      []Param
	ReturnTypes []token.Token
	Body        *BlockStmt
	Private     bool
}

func (node *FuncDeclStmt) Pos() token.Position {
	return node.DefToken.Position
}

func (node *FuncDeclStmt) stmtNode() {}
