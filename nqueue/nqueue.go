package nqueue

// NQueue is a simple queue
type NQueue[T any] struct {
	elements []T
}

// New creates a new queue
func New[T any]() *NQueue[T] {
	return &NQueue[T]{
		elements: make([]T, 0),
	}
}

// Enqueue adds an element to the end of the queue
func (nqueue *NQueue[T]) Enqueue(element T) {
	nqueue.elements = append(nqueue.elements, element)
}

// Dequeue removes an element from the front of the queue
func (nqueue *NQueue[T]) Dequeue() (T, error) {
	if len(nqueue.elements) == 0 {
		var zero T
		return zero, &ErrEmptyQueue{}
	}
	element := nqueue.elements[0]
	nqueue.elements = nqueue.elements[1:]
	return element, nil
}

// Peek returns an element from the front of the queue
func (nqueue *NQueue[T]) Peek() (T, error) {
	if len(nqueue.elements) == 0 {
		var zero T
		return zero, &ErrEmptyQueue{}
	}
	element := nqueue.elements[0]
	return element, nil
}

// IsEmpty checks if queue is empty
func (nqueue *NQueue[T]) IsEmpty() bool {
	return len(nqueue.elements) == 0
}

// Len returns length of the queue
func (nqueue *NQueue[T]) Len() int {
	return len(nqueue.elements)
}
