package service

import (
	"fmt"
	"net"
	"time"
	"github.com/maxdevelopment/go-whois-service/config"
)

type clientInfo struct {
	validThru time.Time
}

func (ci *clientInfo) isValidCache() bool {
	diff := time.Now().Sub(ci.validThru)
	if diff >= 0 {
		return false
	}
	return true
}

func (ci *clientInfo) fetch(ip string) {
	fmt.Println("FETCHING")
	Fetch.getData(ip)
}

type cached struct {
	IPs map[string]*clientInfo
}

type whois struct {
	ClientIPs chan net.Addr
}

var WH = whois{
	ClientIPs: make(chan net.Addr),
}

var cache = cached{
	IPs: make(map[string]*clientInfo),
}

func (wh *whois) Listen() {
	go func() {
		for ip := range wh.ClientIPs {
			ip, _, err := net.SplitHostPort(ip.String())
			if err != nil {
				continue
			}
			fmt.Println(ip)
			if cacheData, ok := cache.IPs[ip]; ok {
				fmt.Println("ip present in the cache")

				if cacheData.isValidCache() {
					fmt.Println("VALID CACHE")
				} else {
					fmt.Println("NOT VALID CACHE")
					cacheData.validThru = time.Now().Add(time.Second * config.Get.ValidThru)
					cacheData.fetch(ip)
				}
			} else {
				fmt.Println("new ip")
				ci := &clientInfo{
					validThru: time.Now().Add(time.Second * config.Get.ValidThru),
				}
				ci.fetch(ip)
				cache.IPs[ip] = ci
			}
		}
	}()
}
