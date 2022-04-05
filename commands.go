package kasa

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SetRelayState changes the relay state of the device -- for multi-relay devices use SetRelayStateChild
func (d *Device) SetRelayState(newstate bool) error {
	if d.Debug {
		klogger.Printf("setting kasa hardware state for [%s] to [%t]", d.IP, newstate)
	}

	state := 0
	if newstate {
		state = 1
	}
	cmd := fmt.Sprintf(CmdSetRelayState, state)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetRelayStateChild adjusts a single relay on a multi-relay device
func (d *Device) SetRelayStateChild(childID string, newstate bool) error {
	if d.Debug {
		klogger.Printf("setting kasa hardware state for [%s] to [%t]", d.IP, newstate)
	}

	state := 0
	if newstate {
		state = 1
	}
	cmd := fmt.Sprintf(CmdSetRelayStateChild, childID, state)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetRelayStateChildMulti adjusts a multiple relays on a multi-relay device
func (d *Device) SetRelayStateChildMulti(newstate bool, children ...string) error {
	if d.Debug {
		klogger.Printf("setting kasa hardware state for [%s] to [%t]", d.IP, newstate)
	}

	state := 0
	if newstate {
		state = 1
	}

	var cc strings.Builder
	first := true
	for _, child := range children {
		if first {
			cc.WriteRune(44) // ,
			first = false
		}
		cc.WriteRune(34) // "
		cc.WriteString(child)
		cc.WriteRune(34) // "
	}

	if d.Debug {
		klogger.Printf("relays [%s]", cc.String())
	}

	cmd := fmt.Sprintf(CmdSetRelayStateChildMulti, cc.String(), state)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetBrightness adjust the brightness setting on a dimmer-capable device (1-100)
func (d *Device) SetBrightness(newval int) error {
	cmd := fmt.Sprintf(CmdSetBrightness, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetFadeOffTime(newval int) error {
	cmd := fmt.Sprintf(CmdSetFadeOffTime, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetFadeOnTime(newval int) error {
	cmd := fmt.Sprintf(CmdSetFadeOnTime, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetGentleOffTime(newval int) error {
	cmd := fmt.Sprintf(CmdSetGentleOffTime, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetGentleOnTime(newval int) error {
	cmd := fmt.Sprintf(CmdSetGentleOnTime, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// GetSettings gets the device sys info
func (d *Device) GetSettings() (*Sysinfo, error) {
	res, err := d.sendTCP(CmdGetSysinfo)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Printf("%+v\n", kd)
	}
	return &kd.GetSysinfo.Sysinfo, nil
}

// GetEmeter returns emeter data from the device
func (d *Device) GetEmeter() (*EmeterRealtime, error) {
	res, err := d.sendTCP(CmdGetEmeter)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var k KasaDevice
	if err = json.Unmarshal(res, &k); err != nil {
		klogger.Println(err.Error())
		klogger.Println(res)
		return nil, err
	}

	if d.Debug {
		klogger.Printf("%+v\n", k)
	}
	return &k.Emeter.Realtime, nil
}

// GetEmeterMonth returns a single month's emeter data from the device
func (d *Device) GetEmeterMonth(month, year int) (*EmeterDaystat, error) {
	q := fmt.Sprintf(CmdEmeterGetMonth, month, year)

	res, err := d.sendTCP(q)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var k KasaDevice
	if err = json.Unmarshal(res, &k); err != nil {
		klogger.Println(err.Error())
		klogger.Println(res)
		return nil, err
	}

	if d.Debug {
		klogger.Printf("%+v\n", k)
	}
	return &k.Emeter.DayStat, nil
}

// DisableCloud sets the device to "local only" mode.
// TODO: forget any cloud settings
func (d *Device) DisableCloud() error {
	err := d.sendUDP(CmdCloudUnbind)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// Reboot instructs the device to reboot
func (d *Device) Reboot() error {
	err := d.sendUDP(CmdReboot)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetLEDOff is insanely named... it should be SetLED, but I'm just going with what TP-Link called these things internally...
func (d *Device) SetLEDOff(t bool) error {
	off := 0
	if t {
		off = 1
	}
	cmd := fmt.Sprintf(CmdLEDOff, off)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetAlias sets a device name
func (d *Device) SetAlias(s string) error {
	cmd := fmt.Sprintf(CmdDeviceAlias, s)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetChildAlias sets the name of an individual relay on a multi-relay device, I don't think this works
func (d *Device) SetChildAlias(childID, s string) error {
	cmd := fmt.Sprintf(CmdChildAlias, childID, s)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetMode sets the target mode of the system
func (d *Device) SetMode(m string) error {
	cmd := fmt.Sprintf(CmdSetMode, m)
	res, err := d.sendTCP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	klogger.Println("SetMode: ", string(res))
	return nil
}

// GetWIFIStatus returns the WiFi station info
func (d *Device) GetWIFIStatus() (string, error) {
	res, err := d.sendTCP(CmdWifiStainfo)
	if err != nil {
		klogger.Println(err.Error())
		return "", err
	}
	return string(res), nil
}

// GetDimmerParameters returns the dimmer parameters from dimmer-capable devices
func (d *Device) GetDimmerParameters() (*DimmerParameters, error) {
	res, err := d.sendTCP(CmdGetDimmer)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}
	var kd KasaDevice
	if err := json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	return &kd.Dimmer.Parameters, nil
}

// GetRules returns the rule information from a device
func (d *Device) GetRules() (string, error) {
	res, err := d.sendTCP(CmdGetRules)
	if err != nil {
		klogger.Println(err.Error())
		return "", err
	}
	return string(res), nil
}

// GetCountdownRules returns a list of the countdown timers on a device
func (d *Device) GetCountdownRules() (*[]Rule, error) {
	res, err := d.sendTCP(CmdGetCountdownRules)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var c KasaDevice
	if err = json.Unmarshal(res, &c); err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Printf("%+v\n", c)
	}
	return &c.Countdown.GetRules.RuleList, nil
}

// ClearCountdownRules resets all countdown rules on the device
func (d *Device) ClearCountdownRules() error {
	err := d.sendUDP(CmdDeleteAllRules)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// AddCountdownRule adds a new countdown
func (d *Device) AddCountdownRule(dur int, target bool, name string) error {
	state := 0
	if target {
		state = 1
	}

	cmd := fmt.Sprintf(CmdAddCountdownRule, dur, state, name)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}
