package nparser

import (
	"strconv"

	"github.com/viveknathani/numero/nlog"
)

type parser struct {
	input    string
	position int
}

// Parse will go through an expression and parse it mathematically
// We start with a vanilla version that supports sum and product operations over just numbers
// Sample input: 2*6+4*5
func Parse(expression string) (float64, error) {
	internalParser := parser{input: expression, position: 0}
	return internalParser.run()
}

func (internalParser *parser) run() (float64, error) {
	sum := 0.0

	err := internalParser.split('+', func() error {
		product := 1.0
		err := internalParser.split('*', func() error {
			if internalParser.peak() == '(' {
				internalParser.next()
				internalSum, internalErr := internalParser.run()
				if internalErr != nil {
					return internalErr
				}
				product *= internalSum
				internalParser.next()
			} else {
				number, err := internalParser.parseNumber()
				if err != nil {
					return err
				}
				product *= number
			}
			return nil
		})
		if err != nil {
			nlog.Error(err)
		}
		sum += product
		return nil
	})
	if err != nil {
		nlog.Error(err)
		return 0, err
	}

	return sum, nil
}

func (internalParser *parser) next() {
	if internalParser.position < len(internalParser.input) {
		internalParser.position++
	}
}

func (internalParser *parser) peak() byte {
	if internalParser.position >= len(internalParser.input) {
		return ' '
	}
	return internalParser.input[internalParser.position]
}

func (internalParser *parser) split(operator byte, callback func() error) error {
	for {
		err := callback()
		if err != nil {
			return err
		}
		if internalParser.peak() == operator {
			internalParser.next()
		} else {
			break
		}
	}

	return nil
}

func (internalParser *parser) parseNumber() (float64, error) {
	start := internalParser.position
	for internalParser.position < len(internalParser.input) && (internalParser.input[internalParser.position] >= '0' && internalParser.input[internalParser.position] <= '9' || internalParser.input[internalParser.position] == '.') {
		internalParser.position++
	}
	numStr := internalParser.input[start:internalParser.position]
	return strconv.ParseFloat(numStr, 64)
}
