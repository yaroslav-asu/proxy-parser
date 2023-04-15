package main

import "proxy-parser/internal/utils/functions"

func main() {
	functions.Init()
	parser := NewParser()
	parser.RemoveEssencesProxies()
	parser.ParseProxies()
	parser.UpdateProxiesWorking()
}
