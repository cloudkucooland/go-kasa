package kasa

import (
	"fmt"
	"net"
)

// http://rat.admin.lv/wp-content/uploads/2018/08/TR17_fgont_-iot_tp_link_hacking.pdf

// Device is the primary type, commands are called from the device
type Device struct {
	IP     string
	parsed net.IP
}

func NewDevice(ip string) (*Device, error) {
	d := Device{IP: ip}
	d.parsed = net.ParseIP(ip)
	if d.parsed == nil {
		addrs, err := net.LookupHost(ip)
		if err != nil {
			return nil, err
		}
		ax := len(addrs)
		if ax == 0 {
			return nil, fmt.Errorf("unknown host: %s", ip)
		}
		d.IP = addrs[0] // XXX make this smarter
		d.parsed = net.ParseIP(d.IP)
	}
	return &d, nil
}

// defined by kasa devices
type kasaDevice struct {
	System kasaSystem `json:"system"`
}

// defined by kasa devices
type kasaSystem struct {
	Sysinfo kasaSysinfo `json:"get_sysinfo"`
}

// defined by kasa devices
type kasaSysinfo struct {
	SWVersion  string `json:"sw_ver"`
	HWVersion  string `json:"hw_ver"`
	Model      string `json:"model"`
	DeviceID   string `json:"deviceId"`
	OEMID      string `json:"oemId"`
	HWID       string `json:"hwId"`
	RSSI       int    `json:"rssi"`
	Longitude  int    `json:"longitude_i"`
	Latitude   int    `json:"latitude_i"`
	Alias      string `json:"alias"`
	Status     string `json:"status"`
	MIC        string `json:"mic_type"`
	Feature    string `json:"feature"`
	MAC        string `json:"mac"`
	Updating   int    `json""updating"`
	LEDOff     int    `json:"led_off"`
	RelayState int    `json:"relay_state"`
	Brightness int    `json:"brightness"`
	OnTime     int    `json:"on_time"`
	ActiveMode string `json:"active_mode"`
	DevName    string `json:"dev_name"`
}
