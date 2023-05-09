package main

import (
	"github.com/yaroslav-asu/proxy-parser/internal/utils/functions"
	"github.com/yaroslav-asu/proxy-parser/parser"
)

func main() {
	functions.Init()
	parser.StartProxiesParsing()
}
