package main

import (
	"fmt"

	"github.com/viveknathani/numero/nlog"
	"github.com/viveknathani/numero/nparser"
)

func main() {
	nlog.Info("hello from numero!")
	nlog.Debug("hello from numero again!")
	nlog.Error("hello from numero again!")

	p1 := nparser.New("2 ^ 3")
	fmt.Println(p1.Run()) // expect 8

	p2 := nparser.New("2 ^ 3 ^ 2")
	fmt.Println(p2.Run()) // expect 512 (2 ^ (3 ^ 2))

	p3 := nparser.New("(2 ^ 3) ^ 2")
	fmt.Println(p3.Run()) // expect 64 ((2 ^ 3) ^ 2)

	p4 := nparser.New("2 * x ^ 3")
	p4.SetVariable("x", 2)
	fmt.Println(p4.Run()) // expect 2 * 8 = 16
}
