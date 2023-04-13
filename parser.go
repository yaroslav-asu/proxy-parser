package main

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"proxy-parser/internal/utils/db"
	"proxy-parser/models"
)

type Parser struct {
	DB *gorm.DB
}

func NewParser() Parser {
	return Parser{
		DB: db.Connect(),
	}
}

func (p *Parser) ParseProxies() {
	resp, err := soup.Get("https://free-proxy-list.net/")
	if err != nil {
		zap.L().Error("Failed to get proxy list")
	}
	doc := soup.HTMLParse(resp)
	for _, proxyRow := range doc.Find("tbody").FindAll("tr") {
		proxyRowElements := proxyRow.FindAll("td")
		proxyValuesText := make([]string, 8)
		for i, proxyValue := range proxyRowElements {
			proxyValuesText[i] = proxyValue.Text()
		}
		p.DB.FirstOrCreate(&models.Proxy{}, models.NewProxyFromArray(proxyValuesText))
	}
}
