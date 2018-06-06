package service

import (
	"net"
	"time"
	"github.com/maxdevelopment/go-whois-service/config"
	"encoding/json"
)

type clientInfo struct {
	ValidThru time.Time `json:"valid_thru"`
	Ip        string    `json:"ip"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Cached    bool      `json:"cached"`
	Link      string    `json:"link"`
}

func (ci *clientInfo) isValidCache() bool {
	diff := time.Now().Sub(ci.ValidThru)
	if diff >= 0 {
		ci.Cached = false
		return false
	}
	ci.Cached = true
	return true
}

func (ci *clientInfo) fetch() {
	Fetch.getData(ci)
}

type cached struct {
	IPs map[string]*clientInfo `json:"ips"`
}

func (c *cached) broadcastData() {
	cache, _ := json.Marshal(c)
	Broadcast <- cache
}

type whois struct {
	ClientIPs chan string
}

var WH = whois{
	ClientIPs: make(chan string),
}

var Cache = cached{
	IPs: make(map[string]*clientInfo),
}

var Broadcast = make(chan []byte)

func (wh *whois) Listen() {
	go func() {
		for ip := range wh.ClientIPs {
			ip, _, err := net.SplitHostPort(ip)
			if err != nil {
				continue
			}

			if ci, ok := Cache.IPs[ip]; ok {

				if !ci.isValidCache() {
					ci.ValidThru = time.Now().Add(time.Second * config.Get.ValidThru)
					ci.fetch()
				}
			} else {
				ci := &clientInfo{
					ValidThru: time.Now().Add(time.Second * config.Get.ValidThru),
					Ip:        ip,
					Cached:    false,
				}
				ci.fetch()
				Cache.IPs[ip] = ci
			}

			Cache.broadcastData()
		}
	}()
}
