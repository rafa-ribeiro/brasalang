package compiler

import (
	"fmt"

	"github.com/rafa-ribeiro/brasalang/internal/ast"
	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
	"github.com/rafa-ribeiro/brasalang/internal/token"
	"github.com/rafa-ribeiro/brasalang/internal/value"
)

type functionMeta struct {
	Name       string
	Arity      byte
	Entry      uint16
	LocalCount byte
	Private    bool
}

type Compiler struct {
	globals   map[string]byte
	functions map[string]byte
}

func New() *Compiler {
	return &Compiler{globals: map[string]byte{}, functions: map[string]byte{}}
}

func (c *Compiler) Compile(program *ast.Program) (*bytecode.Chunk, error) {
	chunk := &bytecode.Chunk{}

	metas := make([]bytecode.FunctionMeta, 0)
	funcDecls := make([]*ast.FuncDeclStmt, 0)
	mainStmts := make([]ast.Stmt, 0)

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncDeclStmt); ok {
			if _, exists := c.functions[fn.Name.Lexeme]; exists {
				return nil, fmt.Errorf("function %q already declared", fn.Name.Lexeme)
			}
			idx := byte(len(c.functions))
			c.functions[fn.Name.Lexeme] = idx
			funcDecls = append(funcDecls, fn)
			continue
		}
		mainStmts = append(mainStmts, stmt)
	}

	seenGlobals := map[string]bool{}
	for _, stmt := range mainStmts {
		decl, ok := stmt.(*ast.VarDeclStmt)
		if !ok {
			continue
		}
		if seenGlobals[decl.Name.Lexeme] {
			return nil, fmt.Errorf("variable %q already declared", decl.Name.Lexeme)
		}
		seenGlobals[decl.Name.Lexeme] = true
		if _, exists := c.globals[decl.Name.Lexeme]; !exists {
			c.globals[decl.Name.Lexeme] = byte(len(c.globals))
		}
	}

	jumpPos := chunk.EmitJump(bytecode.OP_JUMP)

	for _, fn := range funcDecls {
		meta, err := c.emitFunction(chunk, fn)
		if err != nil {
			return nil, err
		}
		metas = append(metas, meta)
	}

	chunk.PatchJump(jumpPos)

	for i, stmt := range mainStmts {
		if err := c.emitStmt(chunk, stmt, nil); err != nil {
			return nil, err
		}

		if i < len(mainStmts)-1 {
			if _, ok := stmt.(*ast.ExprStmt); ok {
				chunk.Write(bytecode.OP_POP)
			}
		}
	}

	chunk.Functions = metas
	return chunk, nil
}

func (c *Compiler) emitFunction(chunk *bytecode.Chunk, fn *ast.FuncDeclStmt) (bytecode.FunctionMeta, error) {
	locals := map[string]byte{}
	for i, p := range fn.Params {
		locals[p.Name.Lexeme] = byte(i)
	}

	entry := uint16(len(chunk.Code))
	localCount := byte(len(locals))
	for i, stmt := range fn.Body.Statements {
		if err := c.emitStmt(chunk, stmt, locals); err != nil {
			return bytecode.FunctionMeta{}, fmt.Errorf("in function %s: %w", fn.Name.Lexeme, err)
		}
		if i < len(fn.Body.Statements)-1 {
			if _, ok := stmt.(*ast.ExprStmt); ok {
				chunk.Write(bytecode.OP_POP)
			}
		}
	}

	// default return value when no explicit return exists.
	chunk.WriteConst(value.NewInt(0))
	chunk.Write(bytecode.OP_RETURN)

	return bytecode.FunctionMeta{Name: fn.Name.Lexeme, Arity: byte(len(fn.Params)), Entry: entry, LocalCount: localCount, Private: fn.Private}, nil
}

func (c *Compiler) emitStmt(chunk *bytecode.Chunk, stmt ast.Stmt, locals map[string]byte) error {
	switch node := stmt.(type) {
	case *ast.ExprStmt:
		return c.emitExpr(chunk, node.Expression, locals)

	case *ast.ReturnStmt:
		if locals == nil {
			return fmt.Errorf("return statement is only allowed inside functions")
		}
		if err := c.emitExpr(chunk, node.Value, locals); err != nil {
			return err
		}
		chunk.Write(bytecode.OP_RETURN)
		return nil

	case *ast.VarDeclStmt:

		if locals == nil {
			if err := c.emitExpr(chunk, node.Initializer, locals); err != nil {
				return err
			}
			slot, exists := c.globals[node.Name.Lexeme]
			if !exists {
				slot = byte(len(c.globals))
				c.globals[node.Name.Lexeme] = slot
			}
			chunk.Write(bytecode.OP_DEFINE_GLOBAL)
			chunk.WriteByte(slot)
			return nil
		}

		if _, exists := locals[node.Name.Lexeme]; exists {
			return fmt.Errorf("local variable %q already declared", node.Name.Lexeme)
		}

		if err := c.emitExpr(chunk, node.Initializer, locals); err != nil {
			return err
		}

		slot := byte(len(locals))
		locals[node.Name.Lexeme] = slot
		chunk.Write(bytecode.OP_DEFINE_LOCAL)
		chunk.WriteByte(slot)
		return nil
	default:
		return fmt.Errorf("unsupported statement type %T", stmt)
	}

}

func (c *Compiler) emitExpr(chunk *bytecode.Chunk, expr ast.Expr, locals map[string]byte) error {
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
		if locals != nil {
			if slot, ok := locals[node.Name]; ok {
				chunk.Write(bytecode.OP_GET_LOCAL)
				chunk.WriteByte(slot)
				return nil
			}
		}
		slot, ok := c.globals[node.Name]
		if !ok {
			return fmt.Errorf("identifier %q is not declared", node.Name)
		}
		chunk.Write(bytecode.OP_GET_GLOBAL)
		chunk.WriteByte(slot)
		return nil

	case *ast.CallExpr:
		for _, arg := range node.Arguments {
			if err := c.emitExpr(chunk, arg, locals); err != nil {
				return err
			}
		}
		fnIdx, ok := c.functions[node.Callee.Lexeme]
		if !ok {
			return fmt.Errorf("function %q is not declared", node.Callee.Lexeme)
		}
		chunk.Write(bytecode.OP_CALL)
		chunk.WriteByte(fnIdx)
		chunk.WriteByte(byte(len(node.Arguments)))
		return nil

	case *ast.UnaryExpr:
		switch node.Operator.Type {
		case token.NOT:
			if err := c.emitExpr(chunk, node.Right, locals); err != nil {
				return err
			}
			chunk.Write(bytecode.OP_NOT)
		case token.MINUS:
			chunk.WriteConst(value.NewInt(0))
			if err := c.emitExpr(chunk, node.Right, locals); err != nil {
				return err
			}
			chunk.Write(bytecode.OP_SUB)
		default:
			return fmt.Errorf("unsupported unary operator %s", node.Operator.Type)
		}
		return nil

	case *ast.BinaryExpr:
		if err := c.emitExpr(chunk, node.Left, locals); err != nil {
			return err
		}
		if err := c.emitExpr(chunk, node.Right, locals); err != nil {
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
