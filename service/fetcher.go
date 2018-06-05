package service

import (
	"net/http"
	"time"
	"fmt"
	"encoding/json"
)

type fetcher struct {
	client *http.Client
}

type RespData struct {
	City string
	Country string
}

func (f *fetcher) getData() error {
	//test
	//resp, err := f.client.Get("http://ip-api.com/json/5.61.45.181")
	//resp, err := f.client.Get("http://ipapi.co/5.61.45.181/json/")
	resp, err := f.client.Get("http://ipfind.co/?ip=5.61.45.181")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rd := RespData{}
	json.NewDecoder(resp.Body).Decode(&rd)

	fmt.Println(rd.City)
	fmt.Println(rd.Country)

	return nil
}

var fetch = fetcher{
	client: &http.Client{Timeout: 10 * time.Second},
}
