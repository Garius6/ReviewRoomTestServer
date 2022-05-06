package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
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

func getUserIdFromToken(r *http.Request) (float64, error) {
	accessTokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if accessTokenString == "" {
		return 0, errors.New("access token string cannot be empty")
	}

	token, err := jwt.ParseWithClaims(accessTokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})
	if err != nil {
		return 0, err
	}

	if uc, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return uc.Id, nil
	} else {
		return 0, errors.New("token is invalid")
	}
}
