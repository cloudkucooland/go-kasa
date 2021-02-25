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
	Countdown  countdown  `json:"count_down"`
	Emeter     emeterSub  `json:"emeter"`
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

// {"emeter":{"get_realtime":{"current_ma":1799,"voltage_mv":121882,"power_mw":174545,"total_wh":547,"err_code":0}}}
// {"emeter":{"get_daystat":{"day_list":[{"year":2021,"month":2,"day":6,"energy_wh":842},{"year":2021,"month":2,"day":7,"energy_wh":1142}],"err_code":0}}}

type emeterSub struct {
	Realtime emeterRealtime `json:"get_realtime"`
	DayStat  emeterDaystat  `json:"get_daystat"`
}

type emeterRealtime struct {
	CurrentMA uint   `json:"current_ma"`
	VoltageMV uint   `json:"voltage_mv"`
	PowerMW   uint   `json:"power_mw"`
	TotalWH   uint   `json:"total_wh"`
	ErrCode   int8   `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
}

type emeterDaystat struct {
	List    []emeterDay `json:"day_list"`
	ErrCode int8        `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
}

type emeterDay struct {
	Year  uint16 `json:"year"`
	Month uint8  `json:"month"`
	Day   uint8  `json:"day"`
	WH    uint16 `json:"energy_wh"`
}

type countdown struct {
	GetRules getRules `json:"get_rules"`
	DelRules delRules `json:"delete_all_rules"`
	AddRule  addRule  `json:"add_rule"`
}

type getRules struct {
	RuleList     []rule `json:"rule_list"`
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}

type rule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Enable    uint8  `json:"enable"`
	Delay     uint16 `json:"delay"`
	Active    uint8  `json:"act"`
	Remaining uint16 `json:"remain"`
}

type delRules struct {
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}

type addRule struct {
	ID           string `json:"id"`
	ErrorCode    int8   `json:"err_code"`
	ErrorMessage string `json:"err_msg"`
}
