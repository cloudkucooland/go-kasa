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
	Debug  bool
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
	NetIf      netif      `json:"netif"`
}

// defined by kasa devices
type getSysinfo struct {
	Sysinfo Sysinfo `json:"get_sysinfo"`
}

// defined by kasa devices
type Sysinfo struct {
	SWVersion      string   `json:"sw_ver"`
	HWVersion      string   `json:"hw_ver"`
	Model          string   `json:"model"`
	DeviceID       string   `json:"deviceId"`
	OEMID          string   `json:"oemId"`
	HWID           string   `json:"hwId"`
	RSSI           int8     `json:"rssi"`
	Longitude      int      `json:"longitude_i"`
	Latitude       int      `json:"latitude_i"`
	Alias          string   `json:"alias"`
	Status         string   `json:"status"`
	MIC            string   `json:"mic_type"`
	Feature        string   `json:"feature"`
	MAC            string   `json:"mac"`
	Updating       uint8    `json""updating"`
	LEDOff         uint8    `json:"led_off"`
	RelayState     uint8    `json:"relay_state"`
	Brightness     uint8    `json:"brightness"`
	OnTime         int      `json:"on_time"`
	ActiveMode     string   `json:"active_mode"`
	DevName        string   `json:"dev_name"`
	Children       []child  `json:"children"`
	NumChildren    uint8    `json:"child_num"`
	NTCState       int      `json:"ntc_state"`
	PreferredState []preset `json:"preferred_state"`
	ErrCode        int8     `json:"error_code"`
}

type dimmer struct {
	Parameters dimmerParameters `json:"get_dimmer_parameters"`
	ErrCode    int8             `json:"err_code"`
	ErrMsg     string           `json:"err_msg"`
}

type dimmerParameters struct {
	MinThreshold  uint16 `json:"minThreshold"`
	FadeOnTime    uint16 `json:"fadeOnTime"`
	FadeOffTime   uint16 `json:"fadeOffTime"`
	GentleOnTime  uint16 `json:"gentleOnTime"`
	GentleOffTime uint16 `json:"gentleOffTime"`
	RampRate      uint16 `json:"rampRate"`
	BulbType      uint8  `json:"bulb_type"`
	ErrCode       int8   `json:"err_code"`
	ErrMsg        string `json:"err_msg"`
}

type child struct {
	ID         string `json:"id"`
	RelayState uint8  `json:"state"`
	Alias      string `json:"alias"`
	OnTime     int    `json:"on_time"`
	// NextAction
}

type preset struct {
	Index      uint8 `json:"index"`
	Brightness uint8 `json:"brightness"`
}

//{"netif":{"get_stainfo":{"ssid":"IoT8417","key_type":3,"rssi":-61,"err_code":0}}}
type netif struct {
	StaInfo stainfo `json:"get_stainfo"`
	ErrCode int8    `json:"err_code"`
	ErrMsg  string  `json:"err_msg"`
}

type stainfo struct {
	SSID    string `json:"ssid"`
	KeyType int8   `json:"key_type"`
	RSSI    int8   `json:"rssi"`
	ErrCode int8   `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
