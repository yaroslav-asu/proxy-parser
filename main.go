package main

import (
	"github.com/yaroslav-asu/proxy-parser/internal/parser"
	"github.com/yaroslav-asu/proxy-parser/internal/utils/functions"
)

func main() {
	functions.Init()
	parser.StartProxiesParsing()
}
