package models

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Proxy struct {
	Id          uint `gorm:"primaryKey"`
	Ip          string
	Port        string
	Code        string
	Country     string
	Https       bool
	IsWorking   bool
	RequestTime time.Duration
	LastWorking time.Time
}

func (p *Proxy) Url() string {
	return fmt.Sprintf("http://%s:%s", p.Ip, p.Port)
}

func convertHttpsToBool(httpsValue string) bool {
	switch strings.ToLower(httpsValue) {
	case "yes":
		return true
	case "no":
		return false
	default:
		zap.L().Error("Failed to convert: " + strings.ToLower(httpsValue) + " to bool")
		return false
	}
}

func NewProxyFromArray(proxyValues []string) *Proxy {
	return &Proxy{
		Ip:          proxyValues[0],
		Port:        proxyValues[1],
		Code:        proxyValues[2],
		Country:     proxyValues[3],
		Https:       convertHttpsToBool(proxyValues[6]),
		IsWorking:   false,
		RequestTime: -1,
		LastWorking: time.Now(),
	}
}

func (p *Proxy) Delete(db *gorm.DB) {
	db.Delete(&p)
}

func (p *Proxy) ToggleWorking() {
	p.IsWorking = !p.IsWorking
	if p.IsWorking {
		p.LastWorking = time.Now()
	}
}
