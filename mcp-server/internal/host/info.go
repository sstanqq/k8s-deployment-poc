package host

import (
	"errors"
	"fmt"
	"net"
	"os"
)

var ErrIPNotFound = errors.New("IP address not found")

type SystemHost struct {
	NodeName string
	NodeIP   string
}

func NewSystemHost(nodeName, nodeIP string) *SystemHost {
	return &SystemHost{
		NodeName: nodeName,
		NodeIP:   nodeIP,
	}
}

func (h *SystemHost) GetNodeInfo() (string, string, error) {
	if h.NodeName == "" {
		hname, err := os.Hostname()
		if err != nil {
			return "", "", fmt.Errorf("cannot get hostname: %w", err)
		}
		h.NodeName = hname
	}

	if h.NodeIP == "" {
		localIP, err := getLocalIP()
		if err != nil {
			return "", "", fmt.Errorf("cannot get local IP: %w", err)
		}
		h.NodeIP = localIP
	}

	return h.NodeName, h.NodeIP, nil
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("getting system's unicast interface addresses: %w", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", ErrIPNotFound
}
