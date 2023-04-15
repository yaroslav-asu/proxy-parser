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
		proxy := models.NewProxyFromArray(proxyValuesText)
		p.DB.Where("ip = ? and port = ?", proxy.Ip, proxy.Port).FirstOrCreate(&proxy)
	}
}

func (p *Parser) CheckProxy(proxy models.Proxy, checkingSites []models.Site) bool {
	parsedProxy, _ := url.Parse(proxy.Url())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxy)}, Timeout: 20 * time.Second}
	cBool := make(chan bool, len(checkingSites))
	for _, site := range checkingSites {
		go func(site models.Site, cBool chan bool) {
			r, err := client.Get(site.Url)
			if err != nil {
				zap.L().Error(err.Error())
			} else {
				zap.L().Info("Success")
			}
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
	p.DB.Order("last_working DESC").Find(&proxies)
	var wg sync.WaitGroup
	var checkingSites []models.Site
	p.DB.Find(&checkingSites)
	counter := 0
	for _, proxy := range proxies {
		counter++
		go func(proxy models.Proxy) {
			wg.Add(1)
			IsWorkingNow := p.CheckProxy(proxy, checkingSites)
			if proxy.IsWorking != IsWorkingNow {
				proxy.ToggleWorking()
				p.DB.Save(&proxy)
				zap.L().Info("Updated work of: " + proxy.Url() + " to: " + strconv.FormatBool(proxy.IsWorking))
			} else {
				zap.L().Info(proxy.Url() + " work value still: " + strconv.FormatBool(proxy.IsWorking))
			}
			defer wg.Done()
		}(proxy)
		if counter == 8 {
			wg.Wait()
			counter = 0
		}
	}
}

func (p *Parser) RemoveEssencesProxies() {
	var proxies []models.Proxy
	p.DB.Where("last_working <= ? and is_working = false", time.Now().Add(-3*time.Hour)).Find(&proxies)
	for _, proxy := range proxies {
		proxy.Delete(p.DB)
	}
}
