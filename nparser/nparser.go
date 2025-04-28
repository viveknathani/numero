package nparser

import "github.com/viveknathani/numero/nlog"

// Nparser will parse mathematical expressions
type Nparser struct {
	operatorList []string
	expression   string
	pointer      int
}

// New returns a new Nparser
func New(expression string) *Nparser {
	return &Nparser{
		operatorList: []string{"+", "-", "*", "/", "(", ")", "^"},
		expression:   expression,
		pointer:      0,
	}
}

// Run runs the parser
func (np *Nparser) Run() {

	for {
		token, ok := np.next()
		if !ok {
			break
		}
		nlog.Debug(token)
	}
}

func (np *Nparser) next() (string, bool) {
	for np.pointer < len(np.expression) &&
		np.expression[np.pointer] == ' ' {
		np.pointer++
	}

	if np.pointer >= len(np.expression) {
		return "", false
	}

	token := string(np.expression[np.pointer])

	if np.isAnOperator(token) {
		np.pointer++
		return token, true
	}

	startIndex := np.pointer
	for np.pointer < len(np.expression) &&
		(np.expression[np.pointer] >= '0' &&
			np.expression[np.pointer] <= '9' ||
			np.expression[np.pointer] == '.') {
		np.pointer++
	}

	return np.expression[startIndex:np.pointer], true
}

func (np *Nparser) isAnOperator(token string) bool {
	for _, op := range np.operatorList {
		if token == op {
			return true
		}
	}
	return false
}
