package nqueue

import "testing"

func TestNewQueue(t *testing.T) {
	queue := New[int]()
	if queue == nil {
		t.Error("New() returned nil")
	}
	if len(queue.elements) != 0 {
		t.Errorf("New() queue not empty, got len %d", len(queue.elements))
	}
}

func TestEnqueueAndPeek(t *testing.T) {
	queue := New[string]()
	
	// Test Peek on empty queue
	_, err := queue.Peek()
	if err == nil {
		t.Error("Peek() on empty queue should return error")
	}
	
	// Test Enqueue and Peek
	queue.Enqueue("hello")
	val, err := queue.Peek()
	if err != nil {
		t.Errorf("Peek() after Enqueue() returned error: %v", err)
	}
	if val != "hello" {
		t.Errorf("Peek() = %v, want %v", val, "hello")
	}
}

func TestDequeue(t *testing.T) {
	queue := New[int]()
	
	// Test Dequeue on empty queue
	_, err := queue.Dequeue()
	if err == nil {
		t.Error("Dequeue() on empty queue should return error")
	}
	
	// Test Enqueue and Dequeue
	queue.Enqueue(42)
	queue.Enqueue(43)
	
	val, err := queue.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() returned error: %v", err)
	}
	if val != 42 {
		t.Errorf("Dequeue() = %v, want %v", val, 42)
	}
	
	val, err = queue.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() returned error: %v", err)
	}
	if val != 43 {
		t.Errorf("Dequeue() = %v, want %v", val, 43)
	}
	
	// Queue should be empty now
	_, err = queue.Dequeue()
	if err == nil {
		t.Error("Dequeue() on empty queue should return error")
	}
}

func TestQueueOrder(t *testing.T) {
	queue := New[int]()
	
	// Enqueue multiple elements
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	
	// Verify FIFO order
	for i := 0; i < 5; i++ {
		val, err := queue.Dequeue()
		if err != nil {
			t.Errorf("Dequeue() returned error: %v", err)
		}
		if val != i {
			t.Errorf("Dequeue() = %v, want %v", val, i)
		}
	}
}
