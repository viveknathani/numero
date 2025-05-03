package main

import (
	"github.com/viveknathani/numero/nlog"
	"github.com/viveknathani/numero/nparser"
)

func main() {
	nlog.Info("hello from numero!")
	nlog.Debug("hello from numero again!")
	nlog.Error("hello from numero again!")

	parser := nparser.New("23/46")
	nlog.Debug(parser.Run())
}
