package nqueue

// ErrEmptyQueue is returned when attempting to access elements from an empty queue
type ErrEmptyQueue struct{}

func (e *ErrEmptyQueue) Error() string {
	return "queue is empty"
}
