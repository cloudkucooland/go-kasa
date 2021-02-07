package kasa

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func BroadcastDiscovery(timeout, probes int) (map[string]*Sysinfo, error) {
	m := make(map[string]*Sysinfo)

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
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					fmt.Printf("discovery failed: %s\n", err.Error())
					return
				}
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
		m[addr.IP.String()] = &kd.GetSysinfo.Sysinfo
	}
	return m, nil
}

// func BroadcastDimmerParam(timeout, probes int) (map[string]*kasaSysinfo, error) {
func BroadcastDimmerParameters(timeout, probes int) (*map[string]*dimmerParameters, error) {
	m := make(map[string]*dimmerParameters)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		fmt.Printf("unable to start discovery listener: %s", err.Error())
		return &m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(`{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`)
		for i := 0; i < probes; i++ {
			// fmt.Println("sending broadcast")
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					fmt.Printf("discovery failed: %s\n", err.Error())
					return
				}
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

		// fmt.Printf("%s\n", res)
		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			fmt.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.Dimmer.ErrCode != 0 {
			// fmt.Printf("%s\n", kd.Dimmer.ErrMsg)
			continue
		}
		// fmt.Printf("%+v\n", kd.Dimmer.Parameters)
		m[addr.IP.String()] = &(kd.Dimmer.Parameters)
	}
	return &m, nil
}

func BroadcastWifiParameters(timeout, probes int) (*map[string]*stainfo, error) {
	m := make(map[string]*stainfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		fmt.Printf("unable to start discovery listener: %s", err.Error())
		return &m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(`{"netif":{"get_stainfo":{}}}`)
		for i := 0; i < probes; i++ {
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					fmt.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])
		// fmt.Println(string(res))

		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			fmt.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.NetIf.ErrCode != 0 {
			fmt.Printf("%s\n", kd.NetIf.ErrMsg)
			continue
		}
		// fmt.Printf("%+v\n", kd.NetIf.StaInfo)
		m[addr.IP.String()] = &(kd.NetIf.StaInfo)
	}
	return &m, nil
}

func BroadcastEmeter(timeout, probes int) (*map[string]string, error) {
	m := make(map[string]string)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		fmt.Printf("unable to start discovery listener: %s", err.Error())
		return &m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(`{"emeter":{"get_realtime":{}}}`)
		for i := 0; i < probes; i++ {
			// fmt.Println("sending broadcast")
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					fmt.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])

		fmt.Printf("%s\n", res)

		// I don't have anything to test with yet
		/* var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			fmt.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.Dimmer.ErrCode != 0 {
			// fmt.Printf("%s\n", kd.Dimmer.ErrMsg)
			continue
		}
		fmt.Printf("%+v\n", kd.Dimmer.Parameters) */
		m[addr.IP.String()] = res
	}
	return &m, nil
}
