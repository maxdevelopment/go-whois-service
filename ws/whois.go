package ws

import (
	"net/http"
	"time"
	"github.com/maxdevelopment/go-whois-service/config"
	"fmt"
	"strings"
	"encoding/json"
)

type whoiser struct {
	client  *http.Client
	servers map[string]*fetchServers
}

type cached struct {
	validThru time.Time
	city      string
	country   string
	link      string
}

type fetchServers struct {
	link   string
	usedAt time.Time
}

type respData struct {
	City    string
	Country string
}

var dataCache = make(map[string]*cached)

var WhoIs = &whoiser{
	client:  &http.Client{Timeout: 10 * time.Second},
	servers: make(map[string]*fetchServers),
}

func (wh *whoiser) flushCache() {
	for key, cache := range dataCache {
		diff := time.Now().Sub(cache.validThru)
		if diff > 0 {
			delete(dataCache, key)
		}
	}
}

func (wh *whoiser) SetServers() {
	ticker := time.NewTicker(time.Second * config.Get.ValidThru)
	go func() {
		for {
			select {
			case <-ticker.C:
				wh.flushCache()
			}
		}
	}()


	for _, link := range config.Get.Servers {
		wh.servers[link] = &fetchServers{
			link:   link,
			usedAt: time.Now(),
		}
	}
}

func (wh *whoiser) isValidCache(client *Client) (*cached, bool) {

	if _, ok := dataCache[client.RemoteAddr]; !ok {
		return nil, false
	}

	diff := time.Now().Sub(dataCache[client.RemoteAddr].validThru)
	if diff >= 0 {
		delete(dataCache, client.RemoteAddr)
		return nil, false
	}

	return dataCache[client.RemoteAddr], true
}

func (wh *whoiser) getLink(client *Client) string {

	var du time.Duration
	var link string

	for key, srv := range wh.servers {
		diff := time.Now().Sub(srv.usedAt)
		if diff > du {
			du = diff
			link = key
		}
	}

	srv := wh.servers[link]
	srv.usedAt = time.Now()

	return strings.Replace(link, "$IP$", client.RemoteAddr, -1)
}

func (wh *whoiser) fetchData(client *Client) {
	link := wh.getLink(client)

	resp, err := wh.client.Get(link)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	rd := respData{}
	json.NewDecoder(resp.Body).Decode(&rd)
	dataCache[client.RemoteAddr] = &cached{
		validThru: time.Now().Add(time.Second * config.Get.ValidThru),
		city:      rd.City,
		country:   rd.Country,
		link:      link,
	}
}

func (wh *whoiser) getData(client *Client) {
	if _, ok := wh.isValidCache(client); !ok {
		client.Cached = false
		wh.fetchData(client)
	}
	cache := dataCache[client.RemoteAddr]
	client.City = cache.city
	client.Country = cache.country
	client.Link = cache.link
}
