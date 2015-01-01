package cpic

import "container/list"

type stack struct {
	stack *list.List
}
type structure int

//newLinear initates stack struct.
func newStack() *stack {
	return &stack{
		stack: list.New(),
	}
}

func (s *stack) top() interface{} {
	l := s.stack.Back()
	return l.Value
}
func (s *stack) push(v interface{}) {
	s.stack.PushBack(v)
}

func (s *stack) pop() interface{} {
	l := s.stack.Back()
	s.stack.Remove(l)
	return l.Value
}
