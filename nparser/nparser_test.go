package nparser

import (
	"math"
	"os"
	"testing"
)

func TestWithSimpleExpression(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("2 + 2")
	result, err := nparser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != 4 {
		t.Errorf("expected 4, got %f", result)
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithVariables(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("x + y")
	nparser.SetVariable("x", 2)
	nparser.SetVariable("y", 2)
	result, err := nparser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != 4 {
		t.Errorf("expected 4, got %f", result)
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithFunctions(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("sin(x)")
	nparser.SetVariable("x", math.Pi/2)
	result, err := nparser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != 1 {
		t.Errorf("expected 1, got %f", result)
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithComplicatedExpression(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("2 + 2 * (3 + 4) / 5")
	result, err := nparser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != 4.8 {
		t.Errorf("expected 4.8, got %f", result)
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithMismatchedParentheses(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("(2 + 2 * 3 + 4) / 5)")
	_, err := nparser.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithUndefinedVariable(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("x + y")
	_, err := nparser.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	os.Unsetenv("LOG_LEVEL")
}

func TestWithMultipleFunctions(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	nparser := New("sin(max(2, 333))")
	result, err := nparser.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != math.Sin(333) {
		t.Errorf("expected %f, got %f", math.Sin(333), result)
	}
	os.Unsetenv("LOG_LEVEL")
}
