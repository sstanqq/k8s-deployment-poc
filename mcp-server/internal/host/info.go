package host

import (
	"errors"
	"fmt"
	"net"
	"os"
)

var ErrIPNotFound = errors.New("IP address not found")

type SystemHost struct{}

func (SystemHost) GetNodeInfo() (string, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", err
	}

	ip, err := getLocalIP()
	if err != nil {
		return "", "", err
	}

	return hostname, ip, nil
}

func getLocalIP() (string, error) {
	var ip string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip, fmt.Errorf("getting system's unicast interface addresses")
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			ip = ipnet.IP.String()
			break
		}
	}

	if len(ip) == 0 {
		return ip, ErrIPNotFound
	}

	return ip, nil
}
