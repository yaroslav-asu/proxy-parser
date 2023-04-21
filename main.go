package main

import (
	"go.uber.org/zap"
	"proxy-parser/internal/utils/functions"
	"time"
)

func main() {
	functions.Init()
	parser := NewParser()
	defer parser.Deconstruct()
	for {
		zap.L().Info("Checking is proxies amount less than 5")
		if parser.WorkingProxiesCount() < 5 {
			zap.L().Info("Starting to update proxies")
			parser.RemoveEssencesProxies()
			parser.ParseProxies()
			parser.UpdateProxiesWorking()
			zap.L().Info("Finished to update proxies")
		} else {
			zap.L().Info("Don't need to update proxies")
		}
		time.Sleep(5 * time.Minute)
	}
}
