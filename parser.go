package main

import (
	"github.com/anaskhan96/soup"
	"github.com/tcnksm/go-httpstat"
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

const requestTimeout = 10 * time.Second

type Parser struct {
	DB *gorm.DB
}

func NewParser() Parser {
	return Parser{
		DB: db.Connect(),
	}
}

func (p *Parser) WorkingProxiesCount() int64 {
	var workingProxiesCount int64
	p.DB.Where("is_working = true").Count(&workingProxiesCount)
	return workingProxiesCount
}

func (p *Parser) ParseProxies() {
	zap.L().Info("Starting to parse proxies")
	resp, err := soup.Get("https://free-proxy-list.net/")
	if err != nil {
		zap.L().Error("Failed to get proxy list")
		return
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

func (p *Parser) CheckProxy(proxy *models.Proxy, checkingSites []models.Site) bool {
	zap.L().Info("Checking is proxy works")
	parsedProxy, _ := url.Parse(proxy.Url())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxy)}, Timeout: requestTimeout}
	cBool := make(chan bool, len(checkingSites))
	for _, site := range checkingSites {
		go func(site models.Site, cBool chan bool) {
			req, err := http.NewRequest("GET", site.Url, nil)
			if err != nil {
				zap.L().Error("Failed to create request")
			}
			var result httpstat.Result
			ctx := httpstat.WithHTTPStat(req.Context(), &result)
			req = req.WithContext(ctx)
			res, err := client.Do(req)
			if err != nil {
				zap.L().Info("failed to connect destination site with error: " + err.Error())
			} else {
				defer res.Body.Close()
				result.End(time.Now())
				proxy.RequestTime = result.StartTransfer
				zap.L().Info("Success")
			}
			cBool <- err == nil && res.StatusCode == 200
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
	zap.L().Info("Updating proxies")
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
			IsWorkingNow := p.CheckProxy(&proxy, checkingSites)
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
	zap.L().Info("Starting to remove essences proxies")
	var proxies []models.Proxy
	p.DB.Where("last_working <= ? and is_working = false", time.Now().Add(-1*time.Hour)).Find(&proxies)
	for _, proxy := range proxies {
		proxy.Delete(p.DB)
	}
}

func (p *Parser) Deconstruct() {
	zap.L().Info("Deconstructing parser")
	db.Close(p.DB)
}
