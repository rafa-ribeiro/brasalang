package vm

import "github.com/rafa-ribeiro/brasalang/internal/value"

type Stack struct {
	values []value.Value
}

func (s *Stack) Push(v value.Value) {
	s.values = append(s.values, v)
}

func (s *Stack) Pop() value.Value {
	if len(s.values) == 0 {
		panic("Stack Underflow")
	}

	lastIndex := len(s.values) - 1
	v := s.values[lastIndex]
	s.values = s.values[:lastIndex]

	return v
}

func (s *Stack) Peek() value.Value {
	if len(s.values) == 0 {
		panic("stack is empty")
	}
	return s.values[len(s.values)-1]
}

func (s *Stack) Size() int {
	return len(s.values)
}
