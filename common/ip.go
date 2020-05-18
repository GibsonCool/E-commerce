package common

import (
	"errors"
	"fmt"
	"net"
)

func GetIntranceIp() (string, error) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range adders {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String())
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("")
}
