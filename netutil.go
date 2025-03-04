package kasa

import (
	"net"
)

// BroadcastAddresses - probably belongs in its own library, get all broadcast addresses
func BroadcastAddresses() ([]net.IP, error) {
	broadcasts := make([]net.IP, 0, 4)
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
			switch a := addr.(type) {
			case *net.IPNet:
				// skip non-IPv4 and non-loopback
				v4 := a.IP.To4()
				if v4 == nil || v4[0] == 127 {
					continue
				}

				for i := 0; i < 4; i++ {
					v4[i] = v4[i] | ^a.Mask[i]
				}
				broadcasts = append(broadcasts, v4)
			default:
				// skip all non-IP
				continue
			}

		}
	}
	return broadcasts, nil
}
