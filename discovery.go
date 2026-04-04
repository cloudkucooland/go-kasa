package kasa

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"time"
)

const bufsize = 2048 // 6-outlet strips cross the 1k mark, double to 2k

// BroadcastDiscovery pulls every attached subnet for kasa devices and returns whatever is discovered
func BroadcastDiscovery(ctx context.Context, probes int) (map[string]*Sysinfo, error) {
	result := make(map[string]*Sysinfo)

	err := discover(ctx, probes, CmdGetSysinfo, func(addr *net.UDPAddr, kd *KasaDevice) error {
		if err := kd.GetSysinfo.Sysinfo.KasaErr.OK(); err != nil {
			klogger.Println(err)
			return nil
		}

		info := kd.GetSysinfo.Sysinfo
		result[addr.IP.String()] = &info
		return nil
	})

	return result, err
}

// BroadcastDimmerParameters  queries all devices on all attached subnets for dimmer state
func BroadcastDimmerParameters(ctx context.Context, probes int) (map[string]*DimmerParameters, error) {
	result := make(map[string]*DimmerParameters)

	err := discover(ctx, probes, CmdGetDimmer, func(addr *net.UDPAddr, kd *KasaDevice) error {
		if err := kd.Dimmer.KasaErr.OK(); err != nil {
			klogger.Println(err)
			return nil
		}
		dimmer := kd.Dimmer.Parameters
		result[addr.IP.String()] = &dimmer
		return nil
	})

	return result, err
}

// BroadcastWifiParameters polls all devices on all attached subnets for wifi status. This is handy when you have one device that never wants to respond, seeing how its wifi status changes over time
func BroadcastWifiParameters(ctx context.Context, probes int) (map[string]*StaInfo, error) {
	result := make(map[string]*StaInfo)

	err := discover(ctx, probes, CmdWifiStainfo, func(addr *net.UDPAddr, kd *KasaDevice) error {
		if err := kd.NetIf.KasaErr.OK(); err != nil {
			klogger.Println(err)
			return nil
		}
		stainfo := kd.NetIf.StaInfo
		result[addr.IP.String()] = &stainfo
		return nil
	})
	return result, err
}

// BroadcastEmeter pulls all devices on all attached subnets for emeter data
func BroadcastEmeter(ctx context.Context, probes int) (map[string]*KasaDevice, error) {
	result := make(map[string]*KasaDevice)

	err := discover(ctx, probes, CmdGetEmeter, func(addr *net.UDPAddr, kd *KasaDevice) error {
		if err := kd.Emeter.KasaErr.OK(); err != nil {
			klogger.Println(err)
			return nil
		}
		if err := kd.Emeter.Realtime.KasaErr.OK(); err != nil {
			klogger.Println(err)
			return nil
		}

		device := *kd
		result[addr.IP.String()] = &device
		return nil
	})
	return result, err
}

func sendBroadcasts(ctx context.Context, cmd string, conn *net.UDPConn, interval time.Duration) {
	payload := Scramble(cmd)
	bcast, err := BroadcastAddresses()
	if err != nil {
		klogger.Println(err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		for _, b := range bcast {
			_, err := conn.WriteToUDP(payload, &net.UDPAddr{IP: b, Port: 9999})
			if err != nil {
				klogger.Println(err)
				return
			}
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func discover(ctx context.Context, probes int, cmd string, handler func(addr *net.UDPAddr, kd *KasaDevice) error) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		klogger.Printf("unable to start listener: %s", err.Error())
		return err
	}
	defer conn.Close()

	buffer := make([]byte, bufsize)

	if probes <= 0 {
		probes = 1
	}

	remaining := 2 * time.Second // default if ctx doesn't have it set
	if deadline, ok := ctx.Deadline(); ok {
		remaining = time.Until(deadline)
	}

	interval := remaining / time.Duration(probes)

	go sendBroadcasts(ctx, cmd, conn, interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue
			}
			return err
		}

		res := Unscramble(buffer[:n])
		if bytes.Contains(res, []byte("module not support")) {
			continue
		}

		var kd KasaDevice
		if err := json.Unmarshal(res, &kd); err != nil {
			klogger.Println(err)
			continue
		}

		if err := handler(addr, &kd); err != nil {
			return err
		}
	}
}
