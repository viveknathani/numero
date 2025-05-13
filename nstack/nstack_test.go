package nstack

import "testing"

func TestNewStack(t *testing.T) {
	stack := New[int]()
	if stack == nil {
		t.Error("New() returned nil")
	}
	if len(stack.elements) != 0 {
		t.Errorf("New() stack not empty, got len %d", len(stack.elements))
	}
}

func TestPushAndTop(t *testing.T) {
	stack := New[string]()

	// test Top on empty stack
	_, err := stack.Top()
	if err == nil {
		t.Error("Top() on empty stack should return error")
	}

	// test Push and Top
	stack.Push("hello")
	val, err := stack.Top()
	if err != nil {
		t.Errorf("Top() after Push() returned error: %v", err)
	}
	if val != "hello" {
		t.Errorf("Top() = %v, want %v", val, "hello")
	}
}

func TestPop(t *testing.T) {
	stack := New[int]()

	// test Pop on empty stack
	_, err := stack.Pop()
	if err == nil {
		t.Error("Pop() on empty stack should return error")
	}

	// test Push and Pop
	stack.Push(42)
	stack.Push(43)

	val, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val != 43 {
		t.Errorf("Pop() = %v, want %v", val, 43)
	}

	val, err = stack.Pop()
	if err != nil {
		t.Errorf("Pop() returned error: %v", err)
	}
	if val != 42 {
		t.Errorf("Pop() = %v, want %v", val, 42)
	}

	// stack should be empty now
	_, err = stack.Pop()
	if err == nil {
		t.Error("Pop() on empty stack should return error")
	}
}

func TestStackOrder(t *testing.T) {
	stack := New[int]()

	// push multiple elements
	for i := 0; i < 5; i++ {
		stack.Push(i)
	}

	// verify LIFO order
	for i := 4; i >= 0; i-- {
		val, err := stack.Pop()
		if err != nil {
			t.Errorf("Pop() returned error: %v", err)
		}
		if val != i {
			t.Errorf("Pop() = %v, want %v", val, i)
		}
	}
}
