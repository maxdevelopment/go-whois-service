package service

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/maxdevelopment/go-whois-service/config"
	"strings"
)

type fetcher struct {
	client  *http.Client
	servers map[string]*fetchServer
}

type RespData struct {
	City    string
	Country string
}

func (f *fetcher) getLink(ip string) string {

	var du time.Duration
	var link string

	for key, srv := range f.servers {
		diff := time.Now().Sub(srv.usedAt)
		if diff > du {
			du = diff
			link = key
		}
	}

	srv := f.servers[link]
	srv.usedAt = time.Now()

	return strings.Replace(link, "$IP$", ip, -1)
}

func (f *fetcher) getData(ci *clientInfo) error {
	link := f.getLink(ci.Ip)
	resp, err := f.client.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rd := RespData{}
	json.NewDecoder(resp.Body).Decode(&rd)
	ci.City = rd.City
	ci.Country = rd.Country
	ci.Link = link

	return nil
}

var Fetch = fetcher{
	client:  &http.Client{Timeout: 10 * time.Second},
	servers: make(map[string]*fetchServer),
}

func (f *fetcher) SetServers() {
	for _, link := range config.Get.Servers {
		f.servers[link] = &fetchServer{
			link:   link,
			usedAt: time.Now(),
		}
	}
}

type fetchServer struct {
	link   string
	usedAt time.Time
}
