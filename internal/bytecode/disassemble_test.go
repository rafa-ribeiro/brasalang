package bytecode

import (
	"strings"
	"testing"

	"github.com/rafa-ribeiro/brasalang/internal/value"
)

func TestDisassembleIncludesConstantsAndJumpTargets(t *testing.T) {
	chunk := &Chunk{}

	chunk.WriteConst(value.NewBool(false))
	jumpIfFalsePos := chunk.EmitJump(OP_JUMP_IF_FALSE)
	chunk.WriteConst(value.NewInt(1))
	chunk.Write(OP_POP)
	jumpPos := chunk.EmitJump(OP_JUMP)
	chunk.PatchJump(jumpIfFalsePos)
	chunk.WriteConst(value.NewInt(2))
	chunk.Write(OP_POP)
	chunk.PatchJump(jumpPos)

	got := chunk.Disassemble()

	checks := []string{
		"OP_CONST",
		"0 (false)",
		"OP_JUMP_IF_FALSE",
		"OP_JUMP",
		"1 (1)",
		"2 (2)",
	}

	for _, check := range checks {
		if !strings.Contains(got, check) {
			t.Fatalf("disassembly missing %q\n%s", check, got)
		}
	}
}
