package kasa

import (
	"context"
	"fmt"
	"net"
	"strconv"
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

	OverrideTCP func(ctx context.Context, cmd string) ([]byte, error)
	OverrideUDP func(ctx context.Context, cmd string) error
}

// NewDevice sets up a new Kasa device for polling
func NewDevice(ip string) (*Device, error) {
	d := Device{Port: 9999}

	d.IP = net.ParseIP(ip)

	// if not an IP address, it might be a hostname, try looking it up
	if d.IP == nil {
		// ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", ip)
		ips, err := net.LookupIP(ip)
		if err != nil {
			return nil, fmt.Errorf("lookup failed for host %q: %w", ip, err)
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
	return net.JoinHostPort(d.IP.String(), strconv.Itoa(d.Port))
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
	GetSysinfo  GetSysinfo  `json:"system"`
	Dimmer      Dimmer      `json:"smartlife.iot.dimmer"`
	NetIf       NetIf       `json:"netif"`
	Countdown   Countdown   `json:"count_down"`
	Emeter      EmeterSub   `json:"emeter"`
	LightSensor LightSensor `json:"smartlife.iot.LAS"`
	CNCloud     CNCloud     `json:"cnCloud"`
	Time        Time        `json:"time"`
	SensorTrig  SensorTrig  `json:"smartlife.iot.sensor_trigger"`
	Schedule    Schedule    `json:"schedule"`
	Bulb        Bulb        `json:"smartlife.iot.smartbulb.lightingservice"`
	Debug       Diagnose    `json:"smartlife.common.debug"`
}

type Bulb struct {
	State           LightState      `json:"get_light_state"`
	PreferredState  PreferredState  `json:"get_preferred_state"`
	DefaultBehavior DefaultBehavior `json:"get_default_behavior"`
	Details         LightDetails    `json:"get_light_details"`
	KasaErr
}

type LightState struct {
	OnOff      int `json:"on_off"`
	Brightness int `json:"brightness"`
	Hue        int `json:"hue"`
	Saturation int `json:"saturation"`
	ColorTemp  int `json:"color_temp"`
	KasaErr
}

type PreferredState struct {
	Index      int `json:"index"`
	Brightness int `json:"brightness"`
	Hue        int `json:"hue"`
	Saturation int `json:"saturation"`
	ColorTemp  int `json:"color_temp"`
	KasaErr
}

type DefaultBehavior struct {
	SoftOn Behavior `json:"soft_on"`
	HardOn Behavior `json:"hard_on"`
	KasaErr
}

type Behavior struct {
	Mode  string `json:"mode"`
	Index int    `json:"index"`
}

type LightDetails struct {
	Wattage   int `json:"wattage"`
	Lumens    int `json:"lumens"`
	BeamAngle int `json:"beam_angle"`
	ColorTemp int `json:"color_temp"`
	KasaErr
}

type Schedule struct {
	Rules     GetSchedRules  `json:"get_rules"`
	DayStat   SchedDayStat   `json:"get_daystat"`
	MonthStat SchedMonthStat `json:"get_monthstat"`
	KasaErr
}

type GetSchedRules struct {
	RuleList []SchedRule `json:"rule_list"`
	KasaErr
}

type SchedRule struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Enable   int    `json:"enable"`
	Wday     []int  `json:"wday"`
	STimeOpt int    `json:"stime_opt"`
	SOffset  int    `json:"soffset"`
	SMin     int    `json:"smin"`
	ETimeOpt int    `json:"etime_opt"`
	EOffset  int    `json:"eoffset"`
	EMin     int    `json:"emin"`
	Freq     int    `json:"frequency"`
	Repeat   int    `json:"repeat"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	Force    int    `json:"force"`
	Duration int    `json:"duration"`
	LastFor  int    `json:"lastFor"`
}

type SchedDayStat struct {
	List []SchedDay `json:"day_list"`
	KasaErr
}

type SchedDay struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
	WH    int `json:"energy_wh"`
}

type SchedMonthStat struct {
	List []SchedMonth `json:"month_list"`
	KasaErr
}

type SchedMonth struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	WH    int `json:"energy_wh"`
}

type SensorTrig struct {
	Routines GetRoutines `json:"get_weekday_routine"`
	Manual   ManualAct   `json:"get_default_manual_action"`
	KasaErr
}

type GetRoutines struct {
	RoutineList []Routine `json:"routine_list"`
	KasaErr
}

type Routine struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	En    int           `json:"en"`
	Wday  []int         `json:"wday"`
	Array []RoutineItem `json:"array"`
}

type RoutineItem struct {
	ST    int `json:"sT"`
	ET    int `json:"eT"`
	Clr   int `json:"clr"`
	OnTT  int `json:"onTT"`
	OffTT int `json:"offTT"`
	OffWT int `json:"offWT"`
}

type ManualAct struct {
	OffToS int `json:"offToS"`
	KasaErr
}

type Time struct {
	Time     TimeData     `json:"get_time"`
	Timezone TimezoneData `json:"get_timezone"`
	KasaErr
}

type TimeData struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Mday  int `json:"mday"`
	Hour  int `json:"hour"`
	Min   int `json:"min"`
	Sec   int `json:"sec"`
}

type TimezoneData struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Mday  int `json:"mday"`
	Hour  int `json:"hour"`
	Min   int `json:"min"`
	Sec   int `json:"sec"`
}

type CNCloud struct {
	Info    CloudInfo    `json:"get_info"`
	SefInfo SefInfo      `json:"get_sefinfo"`
	FwList  FirmwareList `json:"get_intl_fw_list"`
	KasaErr
}

type SefInfo struct {
	SefServer     string `json:"sefServer"`
	DefaultServer string `json:"defaultServer"`
	CachedServer  string `json:"cachedServer"`
	KasaErr
}

type CloudInfo struct {
	Username string `json:"username"`
	Server   string `json:"server"`
	Bind     int    `json:"bind"`
	KasaErr
}

type FirmwareList struct {
	List []Firmware `json:"fw_list"`
	KasaErr
}

type Firmware struct {
	Ver string `json:"ver"`
	Rel string `json:"rel"`
}

// GetSysinfo is defined by kasa devices
type GetSysinfo struct {
	Sysinfo    Sysinfo          `json:"get_sysinfo"`
	BtnCheck   BtnCheckRes      `json:"get_btn_check_res,omitempty"`
	TestMode   TestModeRes      `json:"get_test_mode,omitempty"`
	Onboarding OnboardingStatus `json:"get_onboarding_status,omitempty"`
	SetOnb     SetOnboarding    `json:"set_onboarding_status,omitempty"`
}

type BtnCheckRes struct {
	ResetBtn  bool `json:"reset_btn"`
	SwitchBtn bool `json:"switch_btn"`
	KasaErr
}

type TestModeRes struct {
	FactoryMode int `json:"factory_mode"`
	KasaErr
}

type SetOnboarding struct {
	KasaErr
}

type OnboardingStatus struct {
	Value string `json:"value"`
	KasaErr
}

// Sysinfo is defined by kasa devices
type Sysinfo struct {
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
	OBDSrc     string `json:"obd_src"`
	MIC        string `json:"mic_type"`
	Feature    string `json:"feature"` // "TIM" "TIM:ENE"
	MAC        string `json:"mac"`
	Updating   uint   `json:"updating"`
	LEDOff     uint   `json:"led_off"`
	RelayState uint   `json:"relay_state"`
	Brightness uint   `json:"brightness"`
	OnTime     int    `json:"on_time"`
	IconHash   string `json:"icon_hash"`
	ActiveMode string `json:"active_mode"`
	DevName    string `json:"dev_name"`
	// NextAction     ...      `json:"next_action"`
	Children       []Child  `json:"children"`
	NumChildren    uint     `json:"child_num"`
	NTCState       int      `json:"ntc_state"`
	PreferredState []Preset `json:"preferred_state"`
	KasaErr
}

// "next_action":{"type":-1}

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
// {"netif":{"set_stainfo":{"err_code":0}}}
type NetIf struct {
	StaInfo    StaInfo    `json:"get_stainfo"`
	SetStaInfo SetStaInfo `json:"set_stainfo"`

	KasaErr
}

// StaInfo is defined by kasa devices
type StaInfo struct {
	SSID    string `json:"ssid"`
	KeyType int    `json:"key_type"`
	RSSI    int    `json:"rssi"`
	KasaErr
}

type SetStaInfo struct {
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

// { "smartlife.iot.LAS": { "get_config": { "devs": [ { "hw_id": 0, "enable": 1, "dark_index": 0, "min_adc": 0, "max_adc": 2450, "level_array": [ { "name": "cloudy", "adc": 390, "value": 15 } ] } ], "ver": "1.0", "err_code": 0 } } }
// { "smartlife.iot.LAS": { "get_current_brt": { "value": 0, "err_code": 0 } } }

type LightSensor struct {
	GetConfig     LightSensorConfig     `json:"get_config"`
	GetBrightness LightSensorBrightness `json:"get_current_brt"`
}

type LightSensorConfig struct {
	Devs    []LightSensorDev `json:"devs"`
	Version string           `json:"ver"`
	KasaErr
}

type LightSensorDev struct {
	ID        string             `json:"hw_id"`
	Enable    uint               `json:"enable"`
	DarkIndex uint               `json:"dark_index"`
	MinADC    uint               `json:"min_adc"`
	MaxADC    uint               `json:"max_adc"`
	Levels    []LightSensorLevel `json:"level_array"`
}

type LightSensorLevel struct {
	Name  string `json:"name"`
	ADC   uint   `json:"adc"`
	Value uint   `json:"value"`
}

type LightSensorBrightness struct {
	Value uint `json:"value"`
	KasaErr
}

type PIRSensor struct {
	GetConfig PIRSensorConfig `json:"get_config"`
}

type PIRSensorConfig struct {
	Enable       uint   `json:"enable"`
	Version      string `json:"version"`
	TriggerIndex uint   `json:"trigger_index"`
	ColdTime     uint   `json:"cold_time"`
	MinADC       uint   `json:"min_adc"`
	MaxADC       uint   `json:"max_adc"`
	Data         []uint `json:"array"`
	KasaErr
}

// { "smartlife.iot.PIR": { "get_config": { "enable": 1, "version": "1.0", "trigger_index": 1, "cold_time": 60000, "min_adc": 0, "max_adc": 4095, "array": [80, 50, 20, 0], "err_code": 0 } } }

type Diagnose struct {
	Status      DiagnoseStatus `json:"get_diagnose_status"`
	MCUDiagnose MCUDiagnose    `json:"get_mcu_diagnose"`
}

type MCUDiagnose struct {
	I2CTotalNum       int `json:"i2c_total_num"`
	MCUI2CResetNum    int `json:"mcu_i2c_reset_num"`
	MasterI2CAbnormal int `json:"master_i2c_abnormal_num"`
	KasaErr
}

type DiagnoseStatus struct {
	Result DiagnoseResult `json:"result"`
	KasaErr
}

type DiagnoseResult struct {
	Sysinfo  DiagnoseSysinfo  `json:"sysinfo"`
	Wireless DiagnoseWireless `json:"wireless"`
	Cloud    DiagnoseCloud    `json:"cloud"`
}

type DiagnoseSysinfo struct {
	Uptime  int `json:"uptime"`
	FreeMem int `json:"free_mem"`
	RbtFlag int `json:"rbt_flag"`
	Low     int `json:"low"`
}

type DiagnoseWireless struct {
	RSSI                int    `json:"rssi"`
	Channel             int    `json:"channel"`
	SSID                string `json:"ssid"`
	BSSID               string `json:"bssid"`
	AuthType            int    `json:"auth_type"`
	CipherType          int    `json:"cipherType"`
	ReconnCount         int    `json:"reconn_count"`
	MaxReconnTime       int    `json:"max_reconn_time"`
	TotalDisconnectTime int    `json:"total_disconnect_time"`
}

type DiagnoseCloud struct {
	SendFail  int `json:"send_fail"`
	HbTimeout int `json:"hb_timeout"`
	RecvEOF   int `json:"recv_eof"`
	CloudFsm  int `json:"cloud_fsm"`
	AccFsm    int `json:"acc_fsm"`
}
