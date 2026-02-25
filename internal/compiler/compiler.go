package compiler

import (
	"errors"

	"github.com/rafa-ribeiro/brasalang/internal/ast"
	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
)

type Compiler struct{}

func New() *Compiler {
	return &Compiler{}
}

func (C *Compiler) Compile(program *ast.Program) (*bytecode.Chunk, error) {
	_ = program
	return nil, errors.New("Compiler codegen not implemented yet")
}
