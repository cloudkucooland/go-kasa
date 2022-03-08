package kasa

import (
	"encoding/json"
	"fmt"
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
	cmd := fmt.Sprintf(`{"system":{"set_relay_state":{"state":%d}}}`, state)
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
	cmd := fmt.Sprintf(`{"context":{"child_ids":["%s"]},"system":{"set_relay_state":{"state":%d}}}`, childID, state)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetBrightness adjust the brightness setting on a dimmer-capable device (1-100)
func (d *Device) SetBrightness(newval int) error {
	cmd := fmt.Sprintf(`{"smartlife.iot.dimmer":{"set_brightness":{"brightness":%d}}}`, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

const sysinfo = `{"system":{"get_sysinfo":{}}}`

// GetSettings gets the device sys info
func (d *Device) GetSettings() (*Sysinfo, error) {
	res, err := d.sendTCP(sysinfo)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var kd kasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Printf("%+v\n", kd)
	}
	return &kd.GetSysinfo.Sysinfo, nil
}

const emeter = `{"emeter":{"get_realtime":{}}}`

// GetEmeter returns emeter data from the device
func (d *Device) GetEmeter() (*emeterRealtime, error) {
	res, err := d.sendTCP(emeter)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var k kasaDevice
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

const emeterGetDaystat = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`

// GetEmeterMonth returns a single month's emeter data from the device
func (d *Device) GetEmeterMonth(month, year int) (*emeterDaystat, error) {
	q := fmt.Sprintf(emeterGetDaystat, month, year)

	res, err := d.sendTCP(q)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var k kasaDevice
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

/*
Get EMeter VGain and IGain Settings
{"emeter":{"get_vgain_igain":{}}}

Set EMeter VGain and Igain
{"emeter":{"set_vgain_igain":{"vgain":13462,"igain":16835}}}

Start EMeter Calibration
{"emeter":{"start_calibration":{"vtarget":13462,"itarget":16835}}}
*/

/*
Get Daily Statistic for given Month
{"emeter":{"get_daystat":{"month":1,"year":2016}}}

Get Montly Statistic for given Year
{"emeter":{"get_monthstat":{"year":2016}}}

Erase All EMeter Statistics
{"emeter":{"erase_emeter_stat":null}}
*/

// DisableCloud sets the device to "local only" mode.
// TODO: forget any cloud settings
func (d *Device) DisableCloud() error {
	err := d.sendUDP(`{"cnCloud":{"unbind":null}}`)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// Reboot instructs the device to reboot
func (d *Device) Reboot() error {
	err := d.sendUDP(`{"system":{"reboot":{"delay":2}}}`)
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
	cmd := fmt.Sprintf(`{"system":{"set_led_off":{"off":%d}}}`, off)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetAlias sets a device name
func (d *Device) SetAlias(s string) error {
	cmd := fmt.Sprintf(`{"system":{"set_dev_alias":{"alias":"%s"}}}`, s)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetChildAlias sets the name of an individual relay on a multi-relay device, I don't think this works
func (d *Device) SetChildAlias(childID, s string) error {
	cmd := fmt.Sprintf(`{"context":{"child_ids":["%s"]},"system":{"set_dev_alias":{"alias":"%s"}}}`, childID, s)
	err := d.sendUDP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// SetMode sets the target mode of the system
func (d *Device) SetMode(m string) error {
	cmd := fmt.Sprintf(`{"system":{"set_mode":{"mode":"%s"}}}`, m)
	res, err := d.sendTCP(cmd)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	klogger.Println(res)
	return nil
}

// GetWIFIStatus returns the WiFi station info
func (d *Device) GetWIFIStatus() (string, error) {
	res, err := d.sendTCP(`{"netif":{"get_stainfo":{}}}`)
	if err != nil {
		klogger.Println(err.Error())
		return "", err
	}
	return string(res), nil
}

// GetDimmerParameters returns the dimmer parameters from dimmer-capable devices
func (d *Device) GetDimmerParameters() (string, error) {
	res, err := d.sendTCP(`{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`)
	if err != nil {
		klogger.Println(err.Error())
		return "", err
	}
	return string(res), nil
}

// GetRules returns the rule information from a device
func (d *Device) GetRules() (string, error) {
	res, err := d.sendTCP(`{"smartlife.iot.common.schedule":{"get_rules":{}}}`)
	if err != nil {
		klogger.Println(err.Error())
		return "", err
	}
	return string(res), nil
}

// GetCountdownRules returns a list of the countdown timers on a device
func (d *Device) GetCountdownRules() (*[]rule, error) {
	res, err := d.sendTCP(`{"count_down":{"get_rules":{}}}`)
	if err != nil {
		klogger.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		klogger.Println(res)
	}

	var c kasaDevice
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
	err := d.sendUDP(`{"count_down":{"delete_all_rules":{}}}`)
	if err != nil {
		klogger.Println(err.Error())
		return err
	}
	return nil
}

// https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents

// when I get bored, set myself up as the cloud server... -- make it as responsive as the shellies
// {"cnCloud":{"set_server_url":{"server":"devs.tplinkcloud.com"}}}
// {"cnCloud":{"bind":{"username":alice@home.com, "password":"secret"}}}
