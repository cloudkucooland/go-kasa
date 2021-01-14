package kasa

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func BroadcastDiscovery(timeout, probes int) (map[string]*kasaSysinfo, error) {
	m := make(map[string]*kasaSysinfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		fmt.Printf("unable to start discovery listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(sysinfo)
		for i := 0; i < probes; i++ {
			// fmt.Println("sending broadcast")
			_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 9999})
			if err != nil {
				fmt.Printf("discovery failed: %s\n", err.Error())
				return
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, 1024)
	// fmt.Printf("probing %d times in %d seconds (rate: %d)\n", probes, timeout, timeout / (probes + 1) )
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])
		// fmt.Printf("%s:\n%s\n", addr.IP.String(), res)

		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			fmt.Printf("unmarshal: %s\n", err.Error())
			return nil, err
		}
		m[addr.IP.String()] = &kd.System.Sysinfo
	}
	return m, nil
}
