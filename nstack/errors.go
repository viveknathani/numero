package nstack

// ErrEmptyStack is returned when attempting to access elements from an empty stack
type ErrEmptyStack struct{}

func (e *ErrEmptyStack) Error() string {
	return "stack is empty"
}
