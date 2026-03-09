package value

import "fmt"

type Kind byte

const (
	IntKind Kind = iota
	BoolKind
	NilKind
	TupleKind
)

type Value struct {
	Kind  Kind
	I     int64
	B     bool
	Items []Value
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

func NewNil() Value {
	return Value{Kind: NilKind}
}

func NewTuple(items []Value) Value {
	out := make([]Value, len(items))
	copy(out, items)
	return Value{Kind: TupleKind, Items: out}
}

func (v Value) String() string {
	switch v.Kind {
	case IntKind:
		return fmt.Sprintf("%d", v.I)
	case BoolKind:
		return fmt.Sprintf("%t", v.B)
	case NilKind:
		return "nil"
	case TupleKind:
		return fmt.Sprintf("%v", v.Items)
	default:
		return "unknown"
	}
}
