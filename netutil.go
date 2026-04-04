package kasa

import (
	"net"
)

// BroadcastAddresses - probably belongs in its own library, get all broadcast addresses
func BroadcastAddresses() ([]net.IP, error) {
	var broadcasts []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	seen := make(map[[4]byte]struct{})

	for _, i := range ifaces {
		// Only check interfaces that are UP and support BROADCAST
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagBroadcast == 0 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			v4 := ipNet.IP.To4()
			if v4 == nil || v4.IsLoopback() {
				continue
			}

			// Calculate broadcast: IP | ^Mask
			// We use a fresh 4-byte slice to be safe
			bcast := make(net.IP, 4)
			for j := 0; j < 4; j++ {
				bcast[j] = v4[j] | ^ipNet.Mask[j]
			}

			var key [4]byte
			copy(key[:], v4)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				broadcasts = append(broadcasts, bcast)
			}
		}
	}
	return broadcasts, nil
}
