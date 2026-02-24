package vm

import (
	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
	"github.com/rafa-ribeiro/brasalang/internal/value"
)

// stack -> Store the values in execution
// ip -> Points to the current bytecode
// chunck -> Holds the program (or block of code) in execution

type VM struct {
	stack Stack
	ip    int // Instruction Pointer
	chunk *bytecode.Chunk
}

func New() *VM {
	return &VM{}
}

func (vm *VM) Run(chunk *bytecode.Chunk) {
	vm.chunk = chunk
	vm.ip = 0

	for vm.ip < len(vm.chunk.Code) {

		op := bytecode.OpCode(vm.chunk.Code[vm.ip])
		vm.ip++

		switch op {
		case bytecode.OP_CONST:
			vm.opConst()

		case bytecode.OP_ADD:
			vm.binaryIntOp(func(a, b int64) int64 { return a + b })

		case bytecode.OP_SUB:
			vm.binaryIntOp(func(a, b int64) int64 { return a - b })

		case bytecode.OP_MUL:
			vm.binaryIntOp(func(a, b int64) int64 { return a * b })

		case bytecode.OP_DIV:
			vm.binaryIntOp(func(a, b int64) int64 { return a / b })

		case bytecode.OP_TRUE:
			vm.stack.Push(value.NewBool(true))

		case bytecode.OP_FALSE:
			vm.stack.Push(value.NewBool(false))

		case bytecode.OP_EQUAL:
			vm.binaryCompareOp(func(a, b int64) bool { return a == b })

		case bytecode.OP_NOT_EQUAL:
			vm.binaryCompareOp(func(a, b int64) bool { return a != b })

		case bytecode.OP_GREATER:
			vm.binaryCompareOp(func(a, b int64) bool { return a > b })

		case bytecode.OP_LESS:
			vm.binaryCompareOp(func(a, b int64) bool { return a < b })

		case bytecode.OP_GREATER_EQUAL:
			vm.binaryCompareOp(func(a, b int64) bool { return a >= b })

		case bytecode.OP_LESS_EQUAL:
			vm.binaryCompareOp(func(a, b int64) bool { return a <= b })

		case bytecode.OP_NOT:
			vm.opNot()

		case bytecode.OP_AND:
			vm.binaryBoolOp(func(a, b bool) bool { return a && b })

		case bytecode.OP_OR:
			vm.binaryBoolOp(func(a, b bool) bool { return a || b })

		case bytecode.OP_JUMP:
			vm.opJump()

		case bytecode.OP_JUMP_IF_FALSE:
			vm.opJumpIfFalse()

		case bytecode.OP_POP:
			vm.stack.Pop()

		default:
			panic("unknown opcode")
		}

	}
}

func (vm *VM) StackTop() value.Value {
	return vm.stack.Peek()
}

func (vm *VM) opConst() {
	index := vm.chunk.Code[vm.ip]
	vm.ip++
	constant := vm.chunk.Constants[index]
	vm.stack.Push(constant)
}

func (vm *VM) binaryIntOp(op func(int64, int64) int64) {
	b := vm.stack.Pop()
	a := vm.stack.Pop()

	result := op(a.I, b.I)
	vm.stack.Push(value.NewInt(result))
}

func (vm *VM) binaryCompareOp(op func(int64, int64) bool) {
	b := vm.stack.Pop()
	a := vm.stack.Pop()

	result := op(a.I, b.I)
	vm.stack.Push(value.NewBool(result))
}

func (vm *VM) opNot() {
	v := vm.stack.Pop()

	if v.Kind != value.BoolKind {
		panic("OP_NOT requires boolean")
	}

	vm.stack.Push(value.NewBool(!v.B))
}

// TODO rafael: Change this to implement short-circuit
func (vm *VM) binaryBoolOp(op func(bool, bool) bool) {
	b := vm.stack.Pop()
	a := vm.stack.Pop()

	if a.Kind != value.BoolKind || b.Kind != value.BoolKind {
		panic("Boolean operation requires bool operands")
	}

	result := op(a.B, b.B)
	vm.stack.Push(value.NewBool(result))
}

func (vm *VM) opJump() {
	offset := vm.readUint16()
	vm.ip += int(offset)
}

func (vm *VM) opJumpIfFalse() {
	offset := vm.readUint16()

	condition := vm.stack.Peek()

	if condition.Kind != value.BoolKind {
		panic("OP_JUMP_IF_FALSE requires boolean")
	}

	if !condition.B {
		vm.ip += int(offset)
	}
}

func (vm *VM) readUint16() uint16 {
	high := uint16(vm.chunk.Code[vm.ip])
	low := uint16(vm.chunk.Code[vm.ip+1])

	vm.ip += 2

	return (high << 8) | low

}
