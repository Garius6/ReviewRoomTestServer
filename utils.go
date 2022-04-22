package main

import (
	"fmt"
	"net"
	"strings"
)

func getLocalIp() {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if strings.Contains(ip.String(), "192") {
				fmt.Print("Host = ")
				fmt.Println(ip)
			}
		}
	}
}
