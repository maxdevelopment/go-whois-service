package service

import (
	"fmt"
	"net"
	"time"
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

func (ci *clientInfo) fetch() {
	fmt.Println("FETCHING")
	fetch.getData()
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
				fmt.Println(cacheData.validThru)

				if cacheData.isValidCache() {
					fmt.Println("VALID CACHE")
				} else {
					fmt.Println("NOT VALID CACHE")
					//delete(cache.IPs, ip)
					cacheData.fetch()
				}
			} else {
				fmt.Println("NEW IP")
				ci := &clientInfo{
					validThru: time.Now().Add(time.Second * 10),
				}
				ci.fetch()
				cache.IPs[ip] = ci
			}
		}
	}()
}
