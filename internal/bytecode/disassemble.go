package bytecode

import (
	"fmt"
	"strings"
)

// Disassemble returns a human-readable view of the chunk opcodes.
// It is useful while evolving the compiler/VM because it makes control-flow
// and constants easy to inspect.
func (c *Chunk) Disassemble() string {
	var out strings.Builder
	i := 0

	for i < len(c.Code) {
		op := OpCode(c.Code[i])
		fmt.Fprintf(&out, "%04d %-16s", i, op)
		i++

		switch op {
		case OP_CONST, OP_DEFINE_GLOBAL, OP_GET_GLOBAL, OP_DEFINE_LOCAL, OP_GET_LOCAL, OP_BUILD_TUPLE:
			if i >= len(c.Code) {
				out.WriteString("<missing operand>\n")
				continue
			}
			idx := int(c.Code[i])
			i++

			if op == OP_CONST {
				if idx < 0 || idx >= len(c.Constants) {
					fmt.Fprintf(&out, "%d <invalid const index>\n", idx)
					continue
				}
				fmt.Fprintf(&out, "%d (%s)\n", idx, c.Constants[idx])
				continue
			}

			if op == OP_BUILD_TUPLE {
				fmt.Fprintf(&out, "count=%d\n", idx)
				continue
			}

			fmt.Fprintf(&out, "%d\n", idx)

		case OP_CALL:
			if i+1 >= len(c.Code) {
				out.WriteString("<missing call operands>\n")
				continue
			}

			fnIdx := c.Code[i]
			argc := c.Code[i+1]
			i += 2
			fmt.Fprintf(&out, "fn=%d argc=%d\n", fnIdx, argc)

		case OP_JUMP, OP_JUMP_IF_FALSE:
			if i+1 >= len(c.Code) {
				out.WriteString("<missing jump offset>\n")
				continue
			}

			offset := (uint16(c.Code[i]) << 8) | uint16(c.Code[i+1])
			target := i + JumpInstructionOperandWidth + int(offset)
			i += JumpInstructionOperandWidth
			fmt.Fprintf(&out, "%d -> %04d\n", offset, target)

		default:
			out.WriteByte('\n')

		}
	}

	return out.String()

}
