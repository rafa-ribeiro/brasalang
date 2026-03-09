package vm

import (
	"fmt"

	"github.com/rafa-ribeiro/brasalang/internal/bytecode"
	"github.com/rafa-ribeiro/brasalang/internal/value"
)

type callFrame struct {
	returnIP int // Where to return after function call
	base     int // Base index in the stack for this function's local variables
	fnIndex  int // Index of the function in execution
}

type VM struct {
	stack   Stack           // Store the values in execution
	ip      int             // Points to the current bytecode instruction being executed
	chunk   *bytecode.Chunk // Current chunk of bytecode (or block of code) being executed
	globals []value.Value   // Global variables storage
	frames  []callFrame     // Call stack frames for function calls
}

func New() *VM {
	return &VM{}
}

func (vm *VM) Run(chunk *bytecode.Chunk) {
	vm.chunk = chunk
	vm.ip = 0
	vm.frames = vm.frames[:0]

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

		case bytecode.OP_DEFINE_GLOBAL:
			vm.opDefineGlobal()

		case bytecode.OP_GET_GLOBAL:
			vm.opGetGlobal()

		case bytecode.OP_DEFINE_LOCAL:
			vm.opDefineLocal()

		case bytecode.OP_GET_LOCAL:
			vm.opGetLocal()

		case bytecode.OP_CALL:
			vm.opCall()

		case bytecode.OP_BUILD_TUPLE:
			vm.opBuildTuple()

		case bytecode.OP_RUNTIME_ERROR:
			panic("function reached end without explicit return")

		case bytecode.OP_RETURN:
			vm.opReturn()

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
	b, a := vm.stack.Pop(), vm.stack.Pop()

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

func (vm *VM) opDefineGlobal() {
	slot := int(vm.chunk.Code[vm.ip])
	vm.ip++

	for len(vm.globals) <= slot {
		vm.globals = append(vm.globals, value.Value{})
	}

	vm.globals[slot] = vm.stack.Pop()
}

func (vm *VM) opGetGlobal() {
	slot := int(vm.chunk.Code[vm.ip])
	vm.ip++

	if slot >= len(vm.globals) {
		panic("global slot not initialized")
	}

	vm.stack.Push(vm.globals[slot])
}

func (vm *VM) opDefineLocal() {
	slot := int(vm.chunk.Code[vm.ip])
	vm.ip++
	if len(vm.frames) == 0 {
		panic("local declaration outside function")
	}
	frame := vm.frames[len(vm.frames)-1]
	idx := frame.base + slot
	v := vm.stack.Pop()

	if idx >= vm.stack.Size() {
		missing := idx + 1 - vm.stack.Size()
		vm.stack.PushAll(make([]value.Value, missing))
	}

	vm.stack.Set(idx, v)
}

func (vm *VM) opGetLocal() {
	slot := int(vm.chunk.Code[vm.ip])
	vm.ip++
	frame := vm.frames[len(vm.frames)-1]
	vm.stack.Push(vm.stack.Get(frame.base + slot))
}

func (vm *VM) opCall() {
	fnIndex := int(vm.chunk.Code[vm.ip])
	argc := int(vm.chunk.Code[vm.ip+1])
	vm.ip += 2

	fn := vm.chunk.Functions[fnIndex]
	if argc != int(fn.Arity) {
		panic(fmt.Sprintf("function %s expects %d args, got %d", fn.Name, fn.Arity, argc))
	}

	base := vm.stack.Size() - argc
	vm.frames = append(vm.frames, callFrame{returnIP: vm.ip, base: base, fnIndex: fnIndex})
	for i := argc; i < int(fn.LocalCount); i++ {
		vm.stack.Push(value.Value{})
	}
	vm.ip = int(fn.Entry)
}

func (vm *VM) opBuildTuple() {
	count := int(vm.chunk.Code[vm.ip])
	vm.ip++

	items := make([]value.Value, count)
	for i := count - 1; i >= 0; i-- {
		items[i] = vm.stack.Pop()
	}

	vm.stack.Push(value.NewTuple(items))
}

func (vm *VM) opReturn() {
	ret := vm.stack.Pop()
	if len(vm.frames) == 0 {
		// Returning from main function, just push the value and exit
		vm.stack.Push(ret)
		vm.ip = len(vm.chunk.Code)
		return
	}

	frame := vm.frames[len(vm.frames)-1]
	vm.frames = vm.frames[:len(vm.frames)-1] // pop the call frame
	vm.stack.Truncate(frame.base)            // Remove every thing that belongs to the executed function from the stack
	vm.stack.Push(ret)                       // Push the return value of the function to the stack for the caller to use
	vm.ip = frame.returnIP                   // Return to the instruction after the call
}

func (vm *VM) readUint16() uint16 {
	high := uint16(vm.chunk.Code[vm.ip])
	low := uint16(vm.chunk.Code[vm.ip+1])

	vm.ip += 2

	return (high << 8) | low

}
