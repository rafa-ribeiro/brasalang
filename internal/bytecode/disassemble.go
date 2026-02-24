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
		case OP_CONST:
			if i >= len(c.Code) {
				out.WriteString("<missing const index>\n")
				continue
			}
			idx := int(c.Code[i])
			i++

			if idx < 0 || idx >= len(c.Constants) {
				fmt.Fprintf(&out, "%d <invalid const index>\n", idx)
				continue
			}
			fmt.Fprintf(&out, "%d (%s)\n", idx, c.Constants[idx])

		case OP_JUMP, OP_JUMP_IF_FALSE:
			if i+1 >= len(c.Code) {
				out.WriteString("<missing jump offset>\n")
				continue
			}

			offset := (uint16(c.Code[i]) << 8) | uint16(c.Code[i+1])
			target := i + 2 + int(offset)
			i += 2
			fmt.Fprintf(&out, "%d -> %04d\n", offset, target)

		default:
			out.WriteByte('\n')

		}
	}

	return out.String()

}
