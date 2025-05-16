package nparser

// ErrUnexpectedChar represents an error when an unexpected character is encountered
type ErrUnexpectedChar struct {
	Char byte
}

func (e ErrUnexpectedChar) Error() string {
	return "unexpected character: " + string(e.Char)
}

// ErrMisplacedComma represents an error when a comma is misplaced or parentheses are mismatched
type ErrMisplacedComma struct{}

func (e ErrMisplacedComma) Error() string {
	return "misplaced comma or mismatched parentheses"
}

// ErrMismatchedParentheses represents an error when parentheses are mismatched
type ErrMismatchedParentheses struct{}

func (e ErrMismatchedParentheses) Error() string {
	return "mismatched parentheses"
}

// ErrUnaryMinusMissingOperand represents an error when unary minus is missing an operand
type ErrUnaryMinusMissingOperand struct{}

func (e ErrUnaryMinusMissingOperand) Error() string {
	return "invalid expression: unary minus missing operand"
}

// ErrNotEnoughOperandsForFunction represents an error when a function has insufficient operands
type ErrNotEnoughOperandsForFunction struct {
	Function string
}

func (e ErrNotEnoughOperandsForFunction) Error() string {
	return "not enough operands for function: " + e.Function
}

// ErrNotEnoughOperands represents an error when an expression has insufficient operands
type ErrNotEnoughOperands struct{}

func (e ErrNotEnoughOperands) Error() string {
	return "invalid expression: not enough operands"
}

// ErrUnsupportedOperator represents an error when an unsupported operator is encountered
type ErrUnsupportedOperator struct {
	Operator string
}

func (e ErrUnsupportedOperator) Error() string {
	return "unsupported operator: " + e.Operator
}

// ErrUndefinedVariable represents an error when an undefined variable is referenced
type ErrUndefinedVariable struct {
	Variable string
}

func (e ErrUndefinedVariable) Error() string {
	return "undefined variable: " + e.Variable
}

// ErrEmptyStack represents an error when the stack is empty at the end of evaluation
type ErrEmptyStack struct{}

func (e ErrEmptyStack) Error() string {
	return "invalid expression: empty stack at the end"
}
