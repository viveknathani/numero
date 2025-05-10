package nparser

import (
	"math"
	"strconv"

	"github.com/viveknathani/numero/nstack"
)

// Nparser will parse mathematical expressions
type Nparser struct {
	operatorList  []string
	expression    string
	pointer       int
	variables     Variables
	functions     map[string]func(...float64) float64
	functionArity map[string]int
}

// New returns a new Nparser
func New(expression string) *Nparser {
	return &Nparser{
		operatorList: []string{"+", "-", "*", "/", "^", "u-"},
		expression:   expression,
		pointer:      0,
		variables:    make(Variables),
		functions: map[string]func(...float64) float64{
			"sin":   func(args ...float64) float64 { return math.Sin(args[0]) },
			"cos":   func(args ...float64) float64 { return math.Cos(args[0]) },
			"tan":   func(args ...float64) float64 { return math.Tan(args[0]) },
			"cosec": func(args ...float64) float64 { return 1.0 / math.Sin(args[0]) },
			"sec":   func(args ...float64) float64 { return 1.0 / math.Cos(args[0]) },
			"cot":   func(args ...float64) float64 { return 1.0 / math.Tan(args[0]) },
			"log":   func(args ...float64) float64 { return math.Log(args[0]) },
			"log10": func(args ...float64) float64 { return math.Log10(args[0]) },
			"log2":  func(args ...float64) float64 { return math.Log2(args[0]) },
			"sqrt":  func(args ...float64) float64 { return math.Sqrt(args[0]) },
			"max": func(args ...float64) float64 {
				if len(args) != 2 {
					panic("max expects 2 args")
				}
				if args[0] > args[1] {
					return args[0]
				}
				return args[1]
			},
			"min": func(args ...float64) float64 {
				if len(args) != 2 {
					panic("max expects 2 args")
				}
				if args[0] < args[1] {
					return args[0]
				}
				return args[1]
			},
		},
		functionArity: map[string]int{
			"sin":   1,
			"cos":   1,
			"tan":   1,
			"cosec": 1,
			"sec":   1,
			"cot":   1,
			"sqrt":  1,
			"log":   1,
			"log10": 1,
			"log2":  1,
			"max":   2,
			"min":   2,
		},
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
	prevToken := ""

	for {
		token, ok := np.next()
		if !ok {
			break
		}

		if token == "-" {
			isUnary := false
			if prevToken == "" || prevToken == "(" || np.isAnOperator(prevToken) {
				isUnary = true
			}
			if isUnary {
				token = "u-"
			}
		}

		if token == "," {
			// Pop operators until "("
			for {
				topMostOperator, allOk := operatorStack.Top()
				if !allOk {
					panic("misplaced comma or mismatched parentheses")
				}
				if topMostOperator == "(" {
					break
				}
				operatorStack.Pop()
				outputQueue = append(outputQueue, topMostOperator)
			}
			continue
		}

		if np.isAnOperator(token) || token == "u-" {

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
		} else if _, isFunction := np.functions[token]; isFunction {
			operatorStack.Push(token)
		} else {
			outputQueue = append(outputQueue, token)
		}

		prevToken = token
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
	if np.isAnOperator(string(ch)) || ch == '(' || ch == ')' || ch == ',' {
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
			((np.expression[np.pointer] >= 'a' && np.expression[np.pointer] <= 'z') ||
				(np.expression[np.pointer] >= 'A' && np.expression[np.pointer] <= 'Z') ||
				(np.expression[np.pointer] >= '0' && np.expression[np.pointer] <= '9') ||
				(np.expression[np.pointer] == '.')) {
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
	prec := map[string]int{
		"+":  1,
		"-":  1,
		"*":  2,
		"/":  2,
		"^":  4,
		"u-": 3,
	}
	leftAssoc := map[string]bool{
		"+":  true,
		"-":  true,
		"*":  true,
		"/":  true,
		"^":  false,
		"u-": false,
	}

	p1 := prec[o1]
	p2 := prec[o2]

	return (p2 > p1) || (p2 == p1 && leftAssoc[o1])
}

func (np *Nparser) eval(rpn []string) float64 {
	stack := nstack.New[float64]()

	for _, token := range rpn {
		if token == "u-" {
			a, ok := stack.Pop()
			if !ok {
				panic("invalid expression: unary minus missing operand")
			}
			stack.Push(-a)
			continue
		}

		if fn, isFunc := np.functions[token]; isFunc {
			arity := np.functionArity[token]
			args := make([]float64, arity)
			for i := arity - 1; i >= 0; i-- { // reverse order
				arg, ok := stack.Pop()
				if !ok {
					panic("not enough operands for function: " + token)
				}
				args[i] = arg
			}
			result := fn(args...)
			stack.Push(result)
			continue
		}

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
			case "^":
				res = math.Pow(a, b)
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
