package ws

import (
	"net/http"
	"time"
	"github.com/maxdevelopment/go-whois-service/config"
	"fmt"
)

type whoiser struct {
	client  *http.Client
	servers map[string]*fetchServers
}

type cached struct {
	validThru time.Time
}

type fetchServers struct {
	link   string
	usedAt time.Time
}

var dataCache = make(map[string]*cached)

var WhoIs = &whoiser{
	client:  &http.Client{Timeout: 10 * time.Second},
	servers: make(map[string]*fetchServers),
}

func (wh *whoiser) SetServers() {
	for _, link := range config.Get.Servers {
		wh.servers[link] = &fetchServers{
			link:   link,
			usedAt: time.Now(),
		}
	}
}

func (wh *whoiser) isValidCache(client *Client) (*cached, bool) {

	if _, ok := dataCache[client.remoteAddr]; !ok {
		return nil, false
	}

	diff := time.Now().Sub(dataCache[client.remoteAddr].validThru)
	if diff >= 0 {
		delete(dataCache, client.remoteAddr)
		return nil, false
	}

	return dataCache[client.remoteAddr], true
}

func (wh *whoiser) getData(client *Client) {
	fmt.Println(client)
	if _, ok := wh.isValidCache(client); !ok {
		fmt.Println("NOT OK")
	}
}