package semantic

import "github.com/rafa-ribeiro/brasalang/internal/ast"

type Type string

const (
	TypeInt  Type = "int"
	TypeBool Type = "bool"
)

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(program *ast.Program) []error {
	_ = program
	// Step 1 infra: semantic rules will be implemented incrementally.
	return nil
}
