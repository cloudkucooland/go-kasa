package kasa

import (
	"fmt"
	"net"
)

// things to read to learn the protocol:
// http://rat.admin.lv/wp-content/uploads/2018/08/TR17_fgont_-iot_tp_link_hacking.pdf
// https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/#TP-Link%20Smart%20Home%20Protocol
// https://medium.com/@hu3vjeen/reverse-engineering-tp-link-kc100-bac4641bf1cd
// https://machinekoder.com/controlling-tp-link-hs100110-smart-plugs-with-machinekit/
// https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents
// https://github.com/p-doyle/Python-KasaSmartPowerStrip
// https://community.hubitat.com/t/release-tp-link-kasa-plug-switch-and-bulb-integration/1675/482
// https://github.com/whitslack/kasa/blob/master/API.md

// Device is the primary type, commands are called from the device
type Device struct {
	IP   net.IP
	Port int
}

// NewDevice sets up a new Kasa device for polling
func NewDevice(ip string) (*Device, error) {
	d := Device{Port: 9999}

	d.IP = net.ParseIP(ip)

	// if not an IP address, it might be a hostname, try looking it up
	if d.IP == nil {
		ips, err := net.LookupIP(ip)
		if err != nil {
			return nil, err
		}

		for _, ip := range ips {
			v4 := ip.To4()
			if v4 == nil || v4.IsLoopback() {
				continue
			}
			d.IP = v4
			// stop after first found
			break
		}

		if d.IP == nil {
			return nil, fmt.Errorf("unknown host: %s", ip)
		}
	}

	return &d, nil
}

func NewDeviceIP(ip net.IP) (*Device, error) {
	d := Device{
		IP:   ip,
		Port: 9999,
	}
	return &d, nil
}

func (d *Device) Addr() string {
	return net.JoinHostPort(d.IP.String(), fmt.Sprintf("%d", d.Port))
}

type KasaErr struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (e KasaErr) OK() error {
	if e.ErrCode != 0 {
		return fmt.Errorf("kasa error %d: %s", e.ErrCode, e.ErrMsg)
	}
	return nil
}

// KasaDevice is the primary type, defined by kasa devices
type KasaDevice struct {
	GetSysinfo GetSysinfo `json:"system"`
	Dimmer     Dimmer     `json:"smartlife.iot.dimmer"`
	NetIf      NetIf      `json:"netif"`
	Countdown  Countdown  `json:"count_down"`
	Emeter     EmeterSub  `json:"emeter"`
}

// GetSysinfo is defined by kasa devices
type GetSysinfo struct {
	Sysinfo Sysinfo `json:"get_sysinfo"`
}

// Sysinfo is defined by kasa devices
type Sysinfo struct {
	SWVersion      string   `json:"sw_ver"`
	HWVersion      string   `json:"hw_ver"`
	Model          string   `json:"model"`
	DeviceID       string   `json:"deviceId"`
	OEMID          string   `json:"oemId"`
	HWID           string   `json:"hwId"`
	RSSI           int      `json:"rssi"`
	Longitude      int      `json:"longitude_i"`
	Latitude       int      `json:"latitude_i"`
	Alias          string   `json:"alias"`
	Status         string   `json:"status"`
	MIC            string   `json:"mic_type"`
	Feature        string   `json:"feature"`
	MAC            string   `json:"mac"`
	Updating       uint     `json:"updating"`
	LEDOff         uint     `json:"led_off"`
	RelayState     uint     `json:"relay_state"`
	Brightness     uint     `json:"brightness"`
	OnTime         int      `json:"on_time"`
	ActiveMode     string   `json:"active_mode"`
	DevName        string   `json:"dev_name"`
	Children       []Child  `json:"children"`
	NumChildren    uint     `json:"child_num"`
	NTCState       int      `json:"ntc_state"`
	PreferredState []Preset `json:"preferred_state"`
	KasaErr
}

// Dimmer is defined by kasa devices
type Dimmer struct {
	Parameters DimmerParameters `json:"get_dimmer_parameters"`
	KasaErr
}

// DimmerParameters is defined by kasa devices
type DimmerParameters struct {
	MinThreshold  uint `json:"minThreshold"`
	FadeOnTime    uint `json:"fadeOnTime"`
	FadeOffTime   uint `json:"fadeOffTime"`
	GentleOnTime  uint `json:"gentleOnTime"`
	GentleOffTime uint `json:"gentleOffTime"`
	RampRate      uint `json:"rampRate"`
	BulbType      uint `json:"bulb_type"`
	KasaErr
}

// Child is defined by kasa devices
type Child struct {
	ID         string `json:"id"`
	RelayState uint   `json:"state"`
	Alias      string `json:"alias"`
	OnTime     int    `json:"on_time"`
	// NextAction
}

// Preset is defined by kasa devices
type Preset struct {
	OnOff      int    `json:"on_off"`
	Index      uint   `json:"index"`
	Brightness uint   `json:"brightness"`
	Mode       string `json:"mode"`
	Hue        int    `json:"hue"`
	Saturation int    `json:"saturation"`
	ColorTemp  int    `json:"color_temp"`
}

// NetIf is defined by kasa devices
// {"netif":{"get_stainfo":{"ssid":"IoT8417","key_type":3,"rssi":-61,"err_code":0}}}
type NetIf struct {
	StaInfo StaInfo `json:"get_stainfo"`
	KasaErr
}

// StaInfo is defined by kasa devices
type StaInfo struct {
	SSID    string `json:"ssid"`
	KeyType int    `json:"key_type"`
	RSSI    int    `json:"rssi"`
	KasaErr
}

// {"emeter":{"get_realtime":{"current_ma":1799,"voltage_mv":121882,"power_mw":174545,"total_wh":547,"err_code":0}}}
// {"emeter":{"get_daystat":{"day_list":[{"year":2021,"month":2,"day":6,"energy_wh":842},{"year":2021,"month":2,"day":7,"energy_wh":1142}],"err_code":0}}}

// EmeterSub is defined by kasa devices
type EmeterSub struct {
	Realtime EmeterRealtime `json:"get_realtime"`
	DayStat  EmeterDaystat  `json:"get_daystat"`
	KasaErr
}

// EmeterRealtime is defined by kasa devices
type EmeterRealtime struct {
	Slot      uint `json:"slot_id"`
	CurrentMA uint `json:"current_ma"`
	VoltageMV uint `json:"voltage_mv"`
	PowerMW   uint `json:"power_mw"`
	TotalWH   uint `json:"total_wh"`
	KasaErr
}

// EmeterDaystat is defined by kasa devices
type EmeterDaystat struct {
	List []EmeterDay `json:"day_list"`
	KasaErr
}

// EmeterDay is defined by kasa devices
type EmeterDay struct {
	Year  uint `json:"year"`
	Month uint `json:"month"`
	Day   uint `json:"day"`
	WH    uint `json:"energy_wh"`
}

// Countdown is defined by kasa devices
type Countdown struct {
	GetRules GetRules `json:"get_rules"`
	DelRules DelRules `json:"delete_all_rules"`
	AddRule  AddRule  `json:"add_rule"`
}

// GetRules is defined by kasa devices
type GetRules struct {
	RuleList []Rule `json:"rule_list"`
	KasaErr
}

// Rule is defined by kasa devices
type Rule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Enable    uint   `json:"enable"`
	Delay     uint   `json:"delay"`
	Active    uint   `json:"act"`
	Remaining uint   `json:"remain"`
}

// DelRules is defined by kasa devices
type DelRules struct {
	KasaErr
}

// AddRule is defined by kasa devices
type AddRule struct {
	ID string `json:"id"`
	KasaErr
}
