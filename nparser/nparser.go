package nparser

import (
	"strconv"

	"github.com/viveknathani/numero/nstack"
)

// Nparser will parse mathematical expressions
type Nparser struct {
	operatorList []string
	expression   string
	pointer      int
	variables    map[string]float64
}

// New returns a new Nparser
func New(expression string) *Nparser {
	return &Nparser{
		operatorList: []string{"+", "-", "*", "/"},
		expression:   expression,
		pointer:      0,
		variables:    make(map[string]float64),
	}
}

// SetVariable assings a value to a variable
func (np *Nparser) SetVariable(name string, value float64) {
	np.variables[name] = value
}

// Run runs the parser
func (np *Nparser) Run() float64 {

	outputQueue := make([]string, 0)
	operatorStack := nstack.New[string]()

	for {
		token, ok := np.next()
		if !ok {
			break
		}
		if np.isAnOperator(token) {

			for {
				topMostOperator, allOk := operatorStack.Top()
				if !allOk {
					break
				}
				if topMostOperator == "(" {
					break
				}
				if np.shouldPop(token, topMostOperator) {
					operatorStack.Pop()
					outputQueue = append(outputQueue, topMostOperator)
				} else {
					break
				}
			}
			operatorStack.Push(token)

		} else if token == "(" {
			operatorStack.Push(token)
		} else if token == ")" {
			for {
				topMostOperator, allOk := operatorStack.Top()
				if !allOk {
					break
				}
				if topMostOperator == "(" {
					operatorStack.Pop()
					break
				}
				operatorStack.Pop()
				outputQueue = append(outputQueue, topMostOperator)
			}
		} else {
			outputQueue = append(outputQueue, token)
		}
	}

	for {
		topMostOperator, ok := operatorStack.Pop()
		if !ok {
			break
		}
		outputQueue = append(outputQueue, topMostOperator)
	}

	return np.eval(outputQueue)
}

func (np *Nparser) next() (string, bool) {
	for np.pointer < len(np.expression) &&
		np.expression[np.pointer] == ' ' {
		np.pointer++
	}

	if np.pointer >= len(np.expression) {
		return "", false
	}

	ch := np.expression[np.pointer]

	// Operator or parenthesis
	if np.isAnOperator(string(ch)) || ch == '(' || ch == ')' {
		np.pointer++
		return string(ch), true
	}

	// Number (digits, optional '.')
	if (ch >= '0' && ch <= '9') || ch == '.' {
		startIndex := np.pointer
		for np.pointer < len(np.expression) &&
			((np.expression[np.pointer] >= '0' && np.expression[np.pointer] <= '9') || np.expression[np.pointer] == '.') {
			np.pointer++
		}
		return np.expression[startIndex:np.pointer], true
	}

	// Variable (letters: a-z, A-Z)
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
		startIndex := np.pointer
		for np.pointer < len(np.expression) &&
			((np.expression[np.pointer] >= 'a' && np.expression[np.pointer] <= 'z') || (np.expression[np.pointer] >= 'A' && np.expression[np.pointer] <= 'Z')) {
			np.pointer++
		}
		return np.expression[startIndex:np.pointer], true
	}

	panic("unexpected character: " + string(ch))
}

func (np *Nparser) isAnOperator(token string) bool {
	for _, op := range np.operatorList {
		if token == op {
			return true
		}
	}
	return false
}

func (np *Nparser) shouldPop(o1, o2 string) bool {
	prec := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	leftAssoc := map[string]bool{"+": true, "-": true, "*": true, "/": true}

	p1 := prec[o1]
	p2 := prec[o2]

	return (p2 > p1) || (p2 == p1 && leftAssoc[o1])
}

func (np *Nparser) eval(rpn []string) float64 {
	stack := nstack.New[float64]()

	for _, token := range rpn {
		if np.isAnOperator(token) {
			// Pop two numbers (b first, then a)
			b, ok1 := stack.Pop()
			a, ok2 := stack.Pop()
			if !ok1 || !ok2 {
				panic("invalid expression: not enough operands")
			}

			var res float64
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				res = a / b
			default:
				panic("unsupported operator: " + token)
			}
			stack.Push(res)
		} else {
			// Convert to float and push
			num, err := strconv.ParseFloat(token, 64)
			if err == nil {
				stack.Push(num)
			} else {
				val, ok := np.variables[token]
				if !ok {
					panic("undefined variable: " + token)
				}
				stack.Push(val)
			}
		}
	}

	// Final result
	result, ok := stack.Pop()
	if !ok {
		panic("invalid expression: empty stack at the end")
	}
	return result
}
