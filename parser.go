package main

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"proxy-parser/internal/utils/db"
	"proxy-parser/models"
	"strconv"
	"sync"
	"time"
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

func (p *Parser) CheckProxy(proxy models.Proxy, checkingSites []models.Site) bool {
	parsedProxy, _ := url.Parse(proxy.Url())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxy)}, Timeout: 20 * time.Second}
	cBool := make(chan bool, len(checkingSites))
	for _, site := range checkingSites {
		go func(site models.Site, cBool chan bool) {
			r, err := client.Get(site.Url)
			cBool <- err == nil && r.StatusCode == 200
		}(site, cBool)
	}
	for range checkingSites {
		isWorking := <-cBool
		if !isWorking {
			return false
		}
	}
	return true
}

func (p *Parser) UpdateProxiesWorking() {
	var proxies []models.Proxy
	p.DB.Find(&proxies)
	var wg sync.WaitGroup
	var checkingSites []models.Site
	p.DB.Find(&checkingSites)
	for _, proxy := range proxies {
		go func(proxy models.Proxy) {
			wg.Add(1)
			IsWorkingNow := p.CheckProxy(proxy, checkingSites)
			if proxy.IsWorking != IsWorkingNow {
				proxy.IsWorking = IsWorkingNow
				p.DB.Save(&proxy)
			}
			zap.L().Info("Updated work of: " + proxy.Url() + " to: " + strconv.FormatBool(proxy.IsWorking))
			defer wg.Done()
		}(proxy)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}
