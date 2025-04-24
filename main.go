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
	fmt.Println(nparser.Parse("2*6*4/5"))
}
