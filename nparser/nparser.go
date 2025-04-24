package nparser

import (
	"strconv"
	"strings"

	"github.com/viveknathani/numero/nlog"
)

// Parse will go through an expression and parse it mathematically
// We start with a vanilla version that supports sum and product operations over just numbers
// Sample input: 2*6+4*5
func Parse(expression string) (float64, error) {
	sum := 0.0

	for _, term := range strings.Split(expression, "+") {
		product := 1.0
		for _, factor := range strings.Split(term, "*") {
			number, err := strconv.ParseFloat(factor, 64)
			if err != nil {
				nlog.Error(err)
				return 0, err
			}
			product *= number
		}

		sum += product
	}

	return sum, nil
}
