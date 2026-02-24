package bytecode

type OpCode byte

const (
	OP_CONST OpCode = iota
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_TRUE
	OP_FALSE
	OP_POP

	OP_EQUAL
	OP_NOT_EQUAL
	OP_GREATER
	OP_LESS
	OP_GREATER_EQUAL
	OP_LESS_EQUAL

	OP_NOT
	OP_AND
	OP_OR

	OP_JUMP
	OP_JUMP_IF_FALSE
)

func (op OpCode) String() string {
	switch op {
	case OP_CONST:
		return "OP_CONST"
	case OP_ADD:
		return "OP_ADD"
	case OP_SUB:
		return "OP_SUB"
	case OP_MUL:
		return "OP_MUL"
	case OP_DIV:
		return "OP_DIV"
	case OP_TRUE:
		return "OP_TRUE"
	case OP_FALSE:
		return "OP_FALSE"
	case OP_POP:
		return "OP_POP"
	case OP_EQUAL:
		return "OP_EQUAL"
	case OP_NOT_EQUAL:
		return "OP_NOT_EQUAL"
	case OP_GREATER:
		return "OP_GREATER"
	case OP_LESS:
		return "OP_LESS"
	case OP_GREATER_EQUAL:
		return "OP_GREATER_EQUAL"
	case OP_LESS_EQUAL:
		return "OP_LESS_EQUAL"
	case OP_NOT:
		return "OP_NOT"
	case OP_AND:
		return "OP_AND"
	case OP_OR:
		return "OP_OR"
	case OP_JUMP:
		return "OP_JUMP"
	case OP_JUMP_IF_FALSE:
		return "OP_JUMP_IF_FALSE"
	default:
		return "OP_UNKNOWN"
	}
}
