package nstack

// Nstack is a simple stack
type Nstack[T any] struct {
	elements []T
}

// New creates a new stack
func New[T any]() *Nstack[T] {
	return &Nstack[T]{
		elements: make([]T, 0),
	}
}

// Push pushes the element to the stack
func (nstack *Nstack[T]) Push(element T) {
	nstack.elements = append(nstack.elements, element)
}

// Pop removes the top most element from the stack and returns it
func (nstack *Nstack[T]) Pop() (T, bool) {
	if len(nstack.elements) == 0 {
		var zero T
		return zero, false
	}
	element := nstack.elements[len(nstack.elements)-1]
	nstack.elements = nstack.elements[:len(nstack.elements)-1]
	return element, true
}
