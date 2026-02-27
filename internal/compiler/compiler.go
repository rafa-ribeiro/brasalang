package compiler

import (
	"fmt"

	"github.com/rafa-ribeiro/brasalang/internal/ast"
	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
	"github.com/rafa-ribeiro/brasalang/internal/token"
	"github.com/rafa-ribeiro/brasalang/internal/value"
)

type Compiler struct{}

func New() *Compiler {
	return &Compiler{}
}

func (C *Compiler) Compile(program *ast.Program) (*bytecode.Chunk, error) {
	chunk := &bytecode.Chunk{}

	for i, stmt := range program.Statements {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			return nil, fmt.Errorf("unsupported statement type %T", stmt)
		}

		if err := C.emitExpr(chunk, exprStmt.Expression); err != nil {
			return nil, err
		}

		if i < len(program.Statements)-1 {
			chunk.Write(bytecode.OP_POP)
		}
	}

	return chunk, nil
}

func (C *Compiler) emitExpr(chunk *bytecode.Chunk, expr ast.Expr) error {
	switch node := expr.(type) {
	case *ast.IntLiteral:
		chunk.WriteConst(value.NewInt(node.Value))
		return nil
	case *ast.BoolLiteral:
		if node.Value {
			chunk.Write(bytecode.OP_TRUE)
		} else {
			chunk.Write(bytecode.OP_FALSE)
		}
		return nil
	case *ast.Identifier:
		return fmt.Errorf("identifier %q is not supported yet", node.Name)
	case *ast.UnaryExpr:
		switch node.Operator.Type {
		case token.NOT:
			if err := C.emitExpr(chunk, node.Right); err != nil {
				return err
			}
			chunk.Write(bytecode.OP_NOT)
		case token.MINUS:
			chunk.WriteConst(value.NewInt(0))
			if err := C.emitExpr(chunk, node.Right); err != nil {
				return err
			}
			chunk.Write(bytecode.OP_SUB)
		default:
			return fmt.Errorf("unsupported unary operator %s", node.Operator.Type)
		}
		return nil
	case *ast.BinaryExpr:
		if err := C.emitExpr(chunk, node.Left); err != nil {
			return err
		}
		if err := C.emitExpr(chunk, node.Right); err != nil {
			return err
		}

		op, err := mapBinaryOperator(node.Operator.Type)
		if err != nil {
			return err
		}
		chunk.Write(op)
		return nil
	default:
		return fmt.Errorf("unsupported expression type %T", expr)
	}
}

func mapBinaryOperator(op token.Type) (bytecode.OpCode, error) {
	switch op {
	case token.PLUS:
		return bytecode.OP_ADD, nil
	case token.MINUS:
		return bytecode.OP_SUB, nil
	case token.STAR:
		return bytecode.OP_MUL, nil
	case token.SLASH:
		return bytecode.OP_DIV, nil
	case token.EQUAL_EQUAL:
		return bytecode.OP_EQUAL, nil
	case token.NOT_EQUAL:
		return bytecode.OP_NOT_EQUAL, nil
	case token.GREATER:
		return bytecode.OP_GREATER, nil
	case token.GREATER_EQ:
		return bytecode.OP_GREATER_EQUAL, nil
	case token.LESS:
		return bytecode.OP_LESS, nil
	case token.LESS_EQ:
		return bytecode.OP_LESS_EQUAL, nil
	case token.AND_AND:
		return bytecode.OP_AND, nil
	case token.OR_OR:
		return bytecode.OP_OR, nil
	default:
		return 0, fmt.Errorf("unsupported binary operator %s", op)
	}
}
