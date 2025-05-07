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

	p1 := nparser.New("sin(max(2, 333))")
	fmt.Println(p1.Run())
}
