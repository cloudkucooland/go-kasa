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
	GetSysinfo getSysinfo `json:"system"`
	Dimmer     dimmer     `json:"smartlife.iot.dimmer"`
}

// defined by kasa devices
type getSysinfo struct {
	Sysinfo Sysinfo `json:"get_sysinfo"`
}

// defined by kasa devices
type Sysinfo struct {
	SWVersion   string  `json:"sw_ver"`
	HWVersion   string  `json:"hw_ver"`
	Model       string  `json:"model"`
	DeviceID    string  `json:"deviceId"`
	OEMID       string  `json:"oemId"`
	HWID        string  `json:"hwId"`
	RSSI        int     `json:"rssi"`
	Longitude   int     `json:"longitude_i"`
	Latitude    int     `json:"latitude_i"`
	Alias       string  `json:"alias"`
	Status      string  `json:"status"`
	MIC         string  `json:"mic_type"`
	Feature     string  `json:"feature"`
	MAC         string  `json:"mac"`
	Updating    int     `json""updating"`
	LEDOff      int     `json:"led_off"`
	RelayState  int     `json:"relay_state"`
	Brightness  int     `json:"brightness"`
	OnTime      int     `json:"on_time"`
	ActiveMode  string  `json:"active_mode"`
	DevName     string  `json:"dev_name"`
	Children    []child `json:"children"`
	NumChildren int     `json:"child_num"`
	NTCState    int     `json:"ntc_state"`
}

type dimmer struct {
	Parameters dimmerParameters `json:"get_dimmer_parameters"`
	ErrCode    int              `json:"err_code"`
	ErrMsg     string           `json:"err_msg"`
}

type dimmerParameters struct {
	MinThreshold  int    `json:"minThreshold"`
	FadeOnTime    int    `json:"fadeOnTime"`
	FadeOffTime   int    `json:"fadeOffTime"`
	GentleOnTime  int    `json:"gentleOnTime"`
	GentleOffTime int    `json:"gentleOffTime"`
	RampRate      int    `json:"rampRate"`
	BulbType      int    `json:"bulb_type"`
	ErrCode       int    `json:"err_code"`
	ErrMsg        string `json:"err_msg"`
}

type child struct {
	ID         string `json:"id"`
	RelayState int    `json:"state"`
	Alias      string `json:"alias"`
	OnTime     int    `json:"on_time"`
	// NextAction
}
