package kasa

import (
	"net"
	"strings"
)

func broadcastAddresses() ([]net.IP, error) {
	var broadcasts []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return broadcasts, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return broadcasts, err
		}

		for _, addr := range addrs {
			as := addr.String()

			// ignore IPv6 and loopback since Kasa devices are v4 only (for now)
			if !strings.Contains(as, ":") && !strings.HasPrefix(as, "127.") {
				_, ipnet, err := net.ParseCIDR(as)
				if err != nil {
					return broadcasts, err
				}
				broadcast := net.IP(make([]byte, 4))
				for i := range ipnet.IP {
					broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
				}
				broadcasts = append(broadcasts, broadcast)
			}
		}
	}
	return broadcasts, nil
}
