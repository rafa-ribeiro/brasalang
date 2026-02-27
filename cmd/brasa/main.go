package main

import (
	"fmt"
	"log"

	"github.com/rafa-ribeiro/brasalang/internal/compiler"
	"github.com/rafa-ribeiro/brasalang/internal/parser"
	"github.com/rafa-ribeiro/brasalang/internal/vm"
)

func main() {
	fmt.Println("Brasa VM starting...")

	// sourceCode := "(1 + 2) * 3 == 9"

	sourceCode := `
	a int = 2 * 5 
	b int = 13 - 3
	a == b
	`

	p := parser.NewFromSource(sourceCode)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		log.Fatalf("parse errors: %v", p.Errors())
	}

	c := compiler.New()
	chunk, err := c.Compile(program)
	if err != nil {
		log.Fatalf("compile error: %v", err)
	}

	fmt.Println("Bytecode:")
	fmt.Println(chunk.Disassemble())

	machine := vm.New()
	machine.Run(chunk)

	fmt.Printf("Source: %s\n", sourceCode)
	fmt.Printf("Result: %s\n", machine.StackTop())

}
