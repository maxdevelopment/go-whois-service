package service

import (
	"fmt"
	"net"
)

type whois struct {
	ClientIPs chan net.Addr
}

var WH = whois{
	ClientIPs: make(chan net.Addr),
}

func (wh *whois) Listen() {
	go func() {
		for ip := range wh.ClientIPs {
			ip, _, err := net.SplitHostPort(ip.String())
			if err != nil {
				continue
			}
			fmt.Println(ip)
		}
	}()
}