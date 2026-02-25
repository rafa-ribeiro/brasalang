package bytecode

import "github.com/rafa-ribeiro/brasalang/internal/value"

const JumpInstructionOperandWidth int = 2

type Chunk struct {
	Code      []byte
	Constants []value.Value
}

// Write the provided OpCode to the Chunk
func (c *Chunk) Write(op OpCode) {
	c.Code = append(c.Code, byte(op))
}

func (c *Chunk) WriteByte(b byte) {
	c.Code = append(c.Code, b)
}

func (c *Chunk) WriteConst(v value.Value) {
	index := c.AddConstant(v)
	c.Write(OP_CONST)
	c.WriteByte(byte(index))
}

func (c *Chunk) AddConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)
	return len(c.Constants) - 1
}

func (c *Chunk) EmitJump(op OpCode) int {
	c.Write(op)

	// reserve two bytes for offset
	c.WriteByte(0)
	c.WriteByte(0)

	return len(c.Code) - JumpInstructionOperandWidth // posição do primeiro byte do offset
}

func (c *Chunk) PatchJump(offsetPos int) {
	jump := len(c.Code) - (offsetPos + JumpInstructionOperandWidth)

	c.Code[offsetPos] = byte(jump >> 8)
	c.Code[offsetPos+1] = byte(jump & 0xff)
}
