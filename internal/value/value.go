package value

import "fmt"

type Kind byte

const (
	IntKind Kind = iota
	BoolKind
)

type Value struct {
	Kind Kind
	I    int64
	B    bool
}

func NewInt(v int64) Value {
	return Value{
		Kind: IntKind,
		I:    v,
	}
}

func NewBool(v bool) Value {
	return Value{
		Kind: BoolKind,
		B:    v,
	}
}

func (v Value) String() string {
	switch v.Kind {
	case IntKind:
		return fmt.Sprintf("%d", v.I)
	case BoolKind:
		return fmt.Sprintf("%t", v.B)
	default:
		return "unknown"
	}
}
