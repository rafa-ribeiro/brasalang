package main

import (
	"fmt"

	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
	"github.com/rafa-ribeiro/brasalang/internal/value"
	"github.com/rafa-ribeiro/brasalang/internal/vm"
)

func main() {
	fmt.Println("Brasa VM starting...")

	chunk := &bytecode.Chunk{}

	// Testing Binary operators

	// Add constants
	// index2 := chunk.AddConstant(value.NewInt(2))
	// index3 := chunk.AddConstant(value.NewInt(1))

	// // Write to bytecode
	// chunk.Write(bytecode.OP_CONST)
	// chunk.Write(bytecode.OpCode(index2))

	// chunk.Write(bytecode.OP_CONST)
	// chunk.Write(bytecode.OpCode(index3))

	// chunk.Write(bytecode.OP_GREATER_EQUAL)

	// Testing NOT operator
	// indexTrue := chunk.AddConstant(value.NewBool(false))

	// chunk.Write(bytecode.OP_CONST)
	// chunk.Write(bytecode.OpCode(indexTrue))
	// chunk.Write(bytecode.OP_NOT)

	// Testing boolean operators
	// trueValue := chunk.AddConstant(value.NewBool(true))
	// falseValue := chunk.AddConstant(value.NewBool(false))

	// chunk.Write(bytecode.OP_CONST)
	// chunk.Write(bytecode.OpCode(falseValue))

	// chunk.Write(bytecode.OP_CONST)
	// chunk.Write(bytecode.OpCode(trueValue))

	// chunk.Write(bytecode.OP_OR)

	// Testing If statement

	// condição
	chunk.WriteConst(value.NewBool(false))

	// JUMP_IF_FALSE
	jumpIfFalsePos := chunk.EmitJump(bytecode.OP_JUMP_IF_FALSE)

	// THEN
	chunk.WriteConst(value.NewInt(1))
	chunk.Write(bytecode.OP_POP)

	// JUMP para pular o ELSE
	jumpPos := chunk.EmitJump(bytecode.OP_JUMP)

	// ELSE começa aqui
	chunk.PatchJump(jumpIfFalsePos)

	chunk.WriteConst(value.NewInt(2))
	chunk.Write(bytecode.OP_POP)

	// fim
	chunk.PatchJump(jumpPos)

	// Create the VM
	machine := vm.New()

	// Print the bytecode
	fmt.Println(chunk.Code)

	fmt.Println(chunk.Constants)

	// RUN code
	machine.Run(chunk)

	// Print the Top of the stack
	result := machine.StackTop()
	fmt.Println("Result: ", result)
}
