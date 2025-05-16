package nparser

import (

	"math"
	"strconv"

	"github.com/viveknathani/numero/nqueue"
	"github.com/viveknathani/numero/nstack"
)

// Operator is an operator type
type Operator string

// Expression is an expression type
type Expression string

// Token is a token type
type Token string

// Variables is a map of variable names to their values
type Variables map[string]float64

// Function is a function type
type Function func(...float64) float64

// FunctionDesc is a function description
type FunctionDesc struct {
	arity int
	fn    Function
}

// FunctionList is a map of function names to their descriptions
type FunctionList map[string]FunctionDesc

const (
	// LPAREN is left parenthesis
	LPAREN = "(" //'

	// RPAREN is right parenthesis
	RPAREN = ")"

	// COMMA is well, a comma
	COMMA = ","

	// UMINUS is unary minus
	UMINUS = "u-"

	// PLUS is addition operator
	PLUS = "+"

	// MINUS is subtraction operator
	MINUS = "-"

	// MUL is multiplication operator
	MUL = "*"

	// DIV is division operator
	DIV = "/"

	// POW is power operator
	POW = "^"
)

var operatorList = []Operator{
	PLUS, MINUS, MUL, DIV, POW, UMINUS,
}

var precedence = map[Operator]int{
	PLUS:   1,
	MINUS:  1,
	MUL:    2,
	DIV:    2,
	POW:    4,
	UMINUS: 3,
}

var isLeftAssociative = map[Operator]bool{
	PLUS:   true,
	MINUS:  true,
	MUL:    true,
	DIV:    true,
	POW:    false,
	UMINUS: false,
}

var functionList = map[string]FunctionDesc{
	"sin":   {arity: 1, fn: func(args ...float64) float64 { return math.Sin(args[0]) }},
	"cos":   {arity: 1, fn: func(args ...float64) float64 { return math.Cos(args[0]) }},
	"tan":   {arity: 1, fn: func(args ...float64) float64 { return math.Tan(args[0]) }},
	"cosec": {arity: 1, fn: func(args ...float64) float64 { return 1.0 / math.Sin(args[0]) }},
	"sec":   {arity: 1, fn: func(args ...float64) float64 { return 1.0 / math.Cos(args[0]) }},
	"cot":   {arity: 1, fn: func(args ...float64) float64 { return 1.0 / math.Tan(args[0]) }},
	"log":   {arity: 1, fn: func(args ...float64) float64 { return math.Log(args[0]) }},
	"log10": {arity: 1, fn: func(args ...float64) float64 { return math.Log10(args[0]) }},
	"log2":  {arity: 1, fn: func(args ...float64) float64 { return math.Log2(args[0]) }},
	"sqrt":  {arity: 1, fn: func(args ...float64) float64 { return math.Sqrt(args[0]) }},
	"max": {arity: 2, fn: func(args ...float64) float64 {
		return math.Max(args[0], args[1])
	}},
	"min": {arity: 2, fn: func(args ...float64) float64 {
		return math.Min(args[0], args[1])
	}},
}

// Nparser is a better parser
type Nparser struct {
	pointer    int
	expression Expression
	variables  Variables
}

// New creates a new Nparser
func New(expression string) *Nparser {
	return &Nparser{
		expression: Expression(expression),
		variables:  make(Variables),
	}
}

// SetVariable assigns a value to a variable
func (np *Nparser) SetVariable(name string, value float64) {
	np.variables[name] = value
}

// isAnOperator checks if a token is an operator
func (np *Nparser) isAnOperator(token Token) bool {
	for _, op := range operatorList {
		if token == Token(op) {
			return true
		}
	}
	return false
}

// next returns the next token, whether it was a valid token, and an error if any
func (np *Nparser) next() (Token, bool, error) {

	np.skipSpaces()

	if np.isEndOfExpression() {
		return "", false, nil
	}

	ch := np.expression[np.pointer]

	if np.isAnOperator(Token(ch)) ||
		string(ch) == LPAREN ||
		string(ch) == RPAREN ||
		string(ch) == COMMA {
		np.pointer++
		return Token(ch), true, nil
	}

	if np.isPartOfNumber(ch) {
		startIndex := np.pointer
		for np.pointer < len(np.expression) &&
			((np.expression[np.pointer] >= '0' && np.expression[np.pointer] <= '9') || np.expression[np.pointer] == '.') {
			np.pointer++
		}
		return Token(np.expression[startIndex:np.pointer]), true, nil
	}

	if np.isStartOfVariable(ch) {
		startIndex := np.pointer
		for np.pointer < len(np.expression) &&
			((np.expression[np.pointer] >= 'a' &&
				np.expression[np.pointer] <= 'z') ||
				(np.expression[np.pointer] >= 'A' &&
					np.expression[np.pointer] <= 'Z') ||
				(np.expression[np.pointer] >= '0' &&
					np.expression[np.pointer] <= '9') ||
				(np.expression[np.pointer] == '.')) {
			np.pointer++
		}
		return Token(np.expression[startIndex:np.pointer]), true, nil
	}

	return "", false, ErrUnexpectedChar{Char: ch}
}

// skipSpaces skips all spaces
func (np *Nparser) skipSpaces() {
	for np.pointer < len(np.expression) && np.expression[np.pointer] == ' ' {
		np.pointer++
	}
}

// isEndOfExpression checks if the pointer is at the end of the expression
func (np *Nparser) isEndOfExpression() bool {
	return np.pointer >= len(np.expression)
}

// isPartOfNumber checks if the character is part of a number
func (np *Nparser) isPartOfNumber(ch byte) bool {
	return (ch >= '0' && ch <= '9') || ch == '.'
}

// isStartOfVariable checks if the character is the start of a variable
func (np *Nparser) isStartOfVariable(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// shouldPop checks if the second operator should be popped from the stack
func (np *Nparser) shouldPop(o1, o2 Operator) bool {
	return (precedence[o2] > precedence[o1]) || (precedence[o2] == precedence[o1] && !isLeftAssociative[o1])
}

// Run runs the parser
func (np *Nparser) Run() (float64, error) {

	var prevToken Token
	outputQueue := nqueue.New[Token]()
	operatorStack := nstack.New[Token]()

	for {
		token, ok, err := np.next()
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}

		if token == MINUS {
			if prevToken == "" || prevToken == LPAREN || np.isAnOperator(prevToken) {
				token = UMINUS
			}
		}

		if token == COMMA {
			for {
				topMostOperator, err := operatorStack.Top()
				if err != nil {
					return 0, ErrMisplacedComma{}
				}
				if topMostOperator == LPAREN {
					break
				}
				operatorStack.Pop()
				outputQueue.Enqueue(topMostOperator)
			}
			continue
		} else if np.isAnOperator(token) {
			for {
				topMostOperator, err := operatorStack.Top()
				if err != nil {
					break
				}
				if topMostOperator == LPAREN {
					break
				}
				if np.shouldPop(Operator(token), Operator(topMostOperator)) {
					operatorStack.Pop()
					outputQueue.Enqueue(topMostOperator)
				} else {
					break
				}
			}
			operatorStack.Push(token)
		} else if token == LPAREN {
			operatorStack.Push(token)
		} else if token == RPAREN {
			for {
				topMostOperator, err := operatorStack.Top()
				if err != nil {
					return 0, ErrMismatchedParentheses{}
				}
				if topMostOperator == LPAREN {
					operatorStack.Pop()
					break
				}
				operatorStack.Pop()
				outputQueue.Enqueue(topMostOperator)
			}
		} else if _, isFunction := functionList[string(token)]; isFunction {
			operatorStack.Push(token)
		} else {
			outputQueue.Enqueue(token)
		}

		prevToken = token
	}

	for {
		topMostOperator, err := operatorStack.Pop()
		if err != nil {
			break
		}
		outputQueue.Enqueue(topMostOperator)
	}

	return np.eval(outputQueue)
}

func (np *Nparser) eval(rpn *nqueue.NQueue[Token]) (float64, error) {
	stack := nstack.New[float64]()

	for {
		token, err := rpn.Dequeue()
		if err != nil {
			break
		}

		if token == UMINUS {
			a, err := stack.Pop()
			if err != nil {
				return 0, ErrUnaryMinusMissingOperand{}
			}
			stack.Push(-a)
			continue
		}

		if fn, isFunc := functionList[string(token)]; isFunc {
			arity := fn.arity
			args := make([]float64, arity)
			for i := arity - 1; i >= 0; i-- {
				arg, err := stack.Pop()
				if err != nil {
					return 0, ErrNotEnoughOperandsForFunction{Function: string(token)}
				}
				args[i] = arg
			}
			result := fn.fn(args...)
			stack.Push(result)
			continue
		}

		if np.isAnOperator(token) {
			// pop two numbers (b first, then a)
			b, err1 := stack.Pop()
			a, err2 := stack.Pop()
			if err1 != nil || err2 != nil {
				return 0, ErrNotEnoughOperands{}
			}

			var res float64
			switch token {
			case PLUS:
				res = a + b
			case MINUS:
				res = a - b
			case MUL:
				res = a * b
			case DIV:
				res = a / b
			case POW:
				res = math.Pow(a, b)
			default:
				return 0, ErrUnsupportedOperator{Operator: string(token)}
			}

			stack.Push(res)
		} else {
			num, err := strconv.ParseFloat(string(token), 64)
			if err == nil {
				stack.Push(num)
			} else {
				val, ok := np.variables[string(token)]
				if !ok {
					return 0, ErrUndefinedVariable{Variable: string(token)}
				}
				stack.Push(val)
			}
		}
	}

	// final result
	result, err := stack.Pop()
	if err != nil {
		return 0, ErrEmptyStack{}
	}

	return result, nil
}
