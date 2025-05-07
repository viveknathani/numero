package main

import (
	"fmt"
	"math"

	"github.com/viveknathani/numero/nlog"
	"github.com/viveknathani/numero/nparser"
)

func main() {
	nlog.Info("hello from numero!")
	nlog.Debug("hello from numero again!")
	nlog.Error("hello from numero again!")

	p := nparser.New("-5 + 3")
	fmt.Println(p.Run()) // expect -2

	p2 := nparser.New("-(2 + 3)")
	fmt.Println(p2.Run()) // expect -5

	p3 := nparser.New("2 * -x")
	p3.SetVariable("x", 4)
	fmt.Println(p3.Run()) // expect -8

	p4 := nparser.New("-sin(-x)")
	p4.SetVariable("x", math.Pi/2)
	fmt.Println(p4.Run()) // expect -1
}
