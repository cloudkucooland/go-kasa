package kasa

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
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
	if err := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
		klogger.Println(err.Error())
	}

	go func() {
		payload := Scramble(CmdGetSysinfo)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := BroadcastAddresses()
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
		if errors.Is(err, os.ErrDeadlineExceeded) {
			break
		}
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := Unscramble(buffer[:n])
		// klogger.Printf("%s:\n%s\n", addr.IP.String(), res)

		var kd KasaDevice
		if err = json.Unmarshal(res, &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			return nil, err
		}
		m[addr.IP.String()] = &kd.GetSysinfo.Sysinfo
	}
	return m, nil
}

// BroadcastDimmerParameters  queries all devices on all attached subnets for dimmer state
func BroadcastDimmerParameters(timeout, probes int) (map[string]*DimmerParameters, error) {
	m := make(map[string]*DimmerParameters)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start dimmer listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
		klogger.Println(err.Error())
	}

	go func() {
		payload := Scramble(CmdGetDimmer)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := BroadcastAddresses()
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
		if errors.Is(err, os.ErrDeadlineExceeded) {
			break
		}
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := Unscramble(buffer[:n])

		// klogger.Printf("%s\n", res)
		var kd KasaDevice
		if err = json.Unmarshal(res, &kd); err != nil {
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
func BroadcastWifiParameters(timeout, probes int) (map[string]*StaInfo, error) {
	m := make(map[string]*StaInfo)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start wifi parameter listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
		klogger.Println(err.Error())
	}

	go func() {
		payload := Scramble(CmdWifiStainfo)
		for i := 0; i < probes; i++ {
			bcast, _ := BroadcastAddresses()
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
		if errors.Is(err, os.ErrDeadlineExceeded) {
			break
		}
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := Unscramble(buffer[:n])
		// res := Unscramble(buffer)
		// klogger.Println(string(res))

		var kd KasaDevice
		if err = json.Unmarshal(res, &kd); err != nil {
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
func BroadcastEmeter(timeout, probes int) (map[string]KasaDevice, error) {
	m := make(map[string]KasaDevice)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start emeter listener: %s", err.Error())
		return m, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
		klogger.Println(err.Error())
	}

	go func() {
		payload := Scramble(CmdGetEmeter)
		for i := 0; i < probes; i++ {
			// klogger.Println("sending broadcast")
			bcast, _ := BroadcastAddresses()
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
		if errors.Is(err, os.ErrDeadlineExceeded) {
			break
		}
		if err != nil {
			klogger.Println(err.Error())
			break
		}
		res := Unscramble(buffer[:n])

		if strings.Contains(string(res), "module not support") {
			continue
		}

		// 192.168.12.180: {"emeter":{"get_realtime":{"current_ma":48,"voltage_mv":126368,"power_mw":4211,"total_wh":3190,"err_code":0}}}
		// 192.168.12.168: {"emeter":{"get_realtime":{"slot_id":0,"current_ma":119,"voltage_mv":125533,"power_mw":9389,"total_wh":4150,"err_code":0}}}

		// I don't have anything to test with yet -- I do now, I need to write this
		var kd KasaDevice
		if err = json.Unmarshal(res, &kd); err != nil {
			klogger.Printf("unmarshal: %s\n", err.Error())
			continue
		}
		if kd.Emeter.ErrCode != 0 || kd.Emeter.Realtime.ErrCode != 0 {
			klogger.Printf("+v\n", string(res), kd.Emeter)
			continue
		}
		m[addr.IP.String()] = kd
	}
	return m, nil
}
