package kasa

import (
	"encoding/json"
	"net"
	"time"
)

const bufsize = 2048 // 6-outlet strips cross the 1k mark, double to 2k

// BroadcastDiscovery pulls every attached subnet for kasa devices and returns whatever is discovered
func BroadcastDiscovery(timeout, probes int) (map[string]*Sysinfo, error) {
	m := make(map[string]*Sysinfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start discovery listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(sysinfo)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					klogger.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, bufsize)
	// klogger.Printf("probing %d times in %d seconds (rate: %d)\n", probes, timeout, timeout / (probes + 1) )
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])
		// klogger.Printf("%s:\n%s\n", addr.IP.String(), res)

		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			return nil, err
		}
		m[addr.IP.String()] = &kd.GetSysinfo.Sysinfo
	}
	return m, nil
}

// func BroadcastDimmerParam(timeout, probes int) (map[string]*kasaSysinfo, error) {

// BroadcastDimmerParameters queries all devices on all attached subnets for dimmer state
func BroadcastDimmerParameters(timeout, probes int) (map[string]*dimmerParameters, error) {
	m := make(map[string]*dimmerParameters)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start discovery listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(`{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					klogger.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, bufsize)
	// klogger.Printf("probing %d times in %d seconds (rate: %d)\n", probes, timeout, timeout / (probes + 1) )
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])

		// klogger.Printf("%s\n", res)
		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.Dimmer.ErrCode != 0 {
			// klogger.Printf("%s\n", kd.Dimmer.ErrMsg)
			continue
		}
		// klogger.Printf("%+v\n", kd.Dimmer.Parameters)
		m[addr.IP.String()] = &(kd.Dimmer.Parameters)
	}
	return m, nil
}

// BroadcastWifiParameters polls all devices on all attached subnets for wifi status. This is handy when you have one device that never wants to respond, seeing how its wifi status changes over time
func BroadcastWifiParameters(timeout, probes int) (map[string]*stainfo, error) {
	m := make(map[string]*stainfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start discovery listener: %s", err.Error())
		return m, err
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
					klogger.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, bufsize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])
		// klogger.Println(string(res))

		var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.NetIf.ErrCode != 0 {
			klogger.Printf("%s\n", kd.NetIf.ErrMsg)
			continue
		}
		// klogger.Printf("%+v\n", kd.NetIf.StaInfo)
		m[addr.IP.String()] = &(kd.NetIf.StaInfo)
	}
	return m, nil
}

// BroadcastEmeter pulls all devices on all attached subnets for emeter data
func BroadcastEmeter(timeout, probes int) (map[string]string, error) {
	m := make(map[string]string)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start discovery listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	go func() {
		payload := encryptUDP(`{"emeter":{"get_realtime":{}}}`)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := broadcastAddresses()
			for _, b := range bcast {
				_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
				if err != nil {
					klogger.Printf("discovery failed: %s\n", err.Error())
					return
				}
			}
			time.Sleep(time.Second * time.Duration(timeout/(probes+1)))
		}
	}()

	buffer := make([]byte, bufsize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := decrypt(buffer[:n])

		klogger.Printf("%s\n", res)

		// I don't have anything to test with yet -- I do now, I need to write this
		/* var kd kasaDevice
		if err = json.Unmarshal([]byte(res), &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.Dimmer.ErrCode != 0 {
			// klogger.Printf("%s\n", kd.Dimmer.ErrMsg)
			continue
		}
		klogger.Printf("%+v\n", kd.Dimmer.Parameters) */
		m[addr.IP.String()] = res
	}
	return m, nil
}
