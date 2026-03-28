package kasa

import (
	"net"
)

// BroadcastAddresses - probably belongs in its own library, get all broadcast addresses
func BroadcastAddresses() ([]net.IP, error) {
	broadcasts := make([]net.IP, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return broadcasts, err
	}

	seen := make(map[string]struct{})

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
			// return broadcasts, err
		}

		for _, addr := range addrs {
			switch a := addr.(type) {
			case *net.IPNet:
				// skip non-IPv4 and loopback
				v4 := a.IP.To4()
				if v4 == nil || v4.IsLoopback() {
					continue
				}

				for j := 0; j < 4; j++ {
					v4[j] = v4[j] | ^a.Mask[j]
				}

				key := v4.String()
				if _, exists := seen[key]; !exists {
					seen[key] = struct{}{}
					broadcasts = append(broadcasts, v4)
				}
			default:
				// skip all non-IP
				continue
			}

		}
	}
	return broadcasts, nil
}
