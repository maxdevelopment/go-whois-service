package service
//
//import (
//	"time"
//	"encoding/json"
//	"fmt"
//	"github.com/maxdevelopment/go-whois-service/ws"
//	"github.com/maxdevelopment/go-whois-service/config"
//	"net"
//)
//
//type clientInfo struct {
//	ValidThru time.Time `json:"valid_thru"`
//	Ip        string    `json:"ip"`
//	City      string    `json:"city"`
//	Country   string    `json:"country"`
//	Cached    bool      `json:"cached"`
//	Link      string    `json:"link"`
//}
//
//func (ci *clientInfo) isValidCache() bool {
//	diff := time.Now().Sub(ci.ValidThru)
//	if diff >= 0 {
//		ci.Cached = false
//		return false
//	}
//	ci.Cached = true
//	return true
//}
//
//func (ci *clientInfo) fetch() {
//	Fetch.getData(ci)
//}
//
//type cached struct {
//	Clients map[string]*clientInfo `json:"clients"`
//}
//
//func (c *cached) broadcastData() {
//	cache, _ := json.Marshal(c)
//	Broadcast <- cache
//}
//
//type ConnectedClient struct {
//	Ip string
//	Id string
//}
//
////var ConnectedClients = make(chan *ConnectedClient)
//var ConnectedClients = make(chan *ws.Client)
//var DisconnectedClients = make(chan *ws.Client)
//
//var Cache = cached{
//	Clients: make(map[string]*clientInfo),
//}
//
//var Broadcast = make(chan []byte)
//
//func Listen() {
//	go func() {
//		for connectedClient := range ConnectedClients {
//			fmt.Println(connectedClient)
//		}
//	}()
//
//	go func() {
//		for disconnectedClient := range DisconnectedClients {
//			fmt.Println(disconnectedClient)
//		}
//	}()
//
//	go func() {
//		for connClient := range ConnectedClients {
//			ip, _, err := net.SplitHostPort(connClient.Ip)
//			if err != nil {
//				continue
//			}
//
//			if ci, ok := Cache.Clients[connClient.Id]; ok {
//
//				if !ci.isValidCache() {
//					ci.ValidThru = time.Now().Add(time.Second * config.Get.ValidThru)
//					ci.fetch()
//				}
//			} else {
//				ci := &clientInfo{
//					ValidThru:   time.Now().Add(time.Second * config.Get.ValidThru),
//					Ip:          ip,
//					Cached:      false,
//				}
//				ci.fetch()
//				Cache.Clients[connClient.Id] = ci
//			}
//
//			Cache.broadcastData()
//		}
//	}()
//
//	go func() {
//		for disConnClientId := range DisconnectedClients {
//
//			if _, ok := Cache.Clients[disConnClientId]; ok {
//				fmt.Println("present", disConnClientId)
//				delete(Cache.Clients, disConnClientId)
//				Cache.broadcastData()
//			}
//		}
//	}()
//}
