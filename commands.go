package kasa

import (
	"encoding/json"
	"fmt"
)

func (d *Device) SetRelayState(newstate bool) error {
	// fmt.Printf("setting kasa hardware state for [%s] to [%t]", a.Name, newstate)
	state := 0
	if newstate {
		state = 1
	}
	cmd := fmt.Sprintf(`{"system":{"set_relay_state":{"state":%d}}}`, state)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetRelayStateChild(childID string, newstate bool) error {
	// fmt.Printf("setting kasa hardware state for [%s] to [%t]", a.Name, newstate)
	state := 0
	if newstate {
		state = 1
	}
	cmd := fmt.Sprintf(`{"context":{"child_ids":["%s"]},"system":{"set_relay_state":{"state":%d}}}`, childID, state)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetBrightness(newval int) error {
	cmd := fmt.Sprintf(`{"smartlife.iot.dimmer":{"set_brightness":{"brightness":%d}}}`, newval)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

const sysinfo = `{"system":{"get_sysinfo":{}}}`

func (d *Device) GetSettings() (*Sysinfo, error) {
	res, err := d.sendTCP(sysinfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Println(res)
	}

	var kd kasaDevice
	if err = json.Unmarshal([]byte(res), &kd); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Printf("%+v\n", kd)
	}
	return &kd.GetSysinfo.Sysinfo, nil
}

const emeter = `{"emeter":{"get_realtime":{}}}`

func (d *Device) GetEmeter() (*emeterRealtime, error) {
	res, err := d.sendTCP(emeter)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Println(res)
	}

	var em emeterTop
	if err = json.Unmarshal([]byte(res), &em); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Printf("%+v\n", em)
	}
	return &em.E.Realtime, nil
}

const emeterGetDaystat = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`

func (d *Device) GetEmeterMonth(month, year int) (*emeterDaystat, error) {
	q := fmt.Sprintf(emeterGetDaystat, month, year)

	res, err := d.sendTCP(q)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Println(res)
	}

	var em emeterTop
	if err = json.Unmarshal([]byte(res), &em); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if d.Debug {
		fmt.Printf("%+v\n", em)
	}
	return &em.E.DayStat, nil
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

// forget any cloud settings
func (d *Device) DisableCloud() error {
	err := d.sendUDP(`{"cnCloud":{"unbind":null}}`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) Reboot() error {
	err := d.sendUDP(`{"system":{"reboot":{"delay":2}}}`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetLEDOff(t bool) error {
	off := 0
	if t {
		off = 1
	}
	cmd := fmt.Sprintf(`{"system":{"set_led_off":{"off":%d}}}`, off)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetAlias(s string) error {
	cmd := fmt.Sprintf(`{"system":{"set_dev_alias":{"alias":"%s"}}}`, s)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetChildAlias(childID, s string) error {
	cmd := fmt.Sprintf(`{"context":{"child_ids":["%s"]},"system":{"set_dev_alias":{"alias":"%s"}}}`, childID, s)
	err := d.sendUDP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Device) SetMode(m string) error {
	cmd := fmt.Sprintf(`{"system":{"set_mode":{"mode":"%m"}}}`, m)
	res, err := d.sendTCP(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(res)
	return nil
}

func (d *Device) GetWIFIStatus() (string, error) {
	res, err := d.sendTCP(`{"netif":{"get_stainfo":{}}}`)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return res, nil
}

func (d *Device) GetDimmerParameters() (string, error) {
	res, err := d.sendTCP(`{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return res, nil
}

func (d *Device) GetRules() (string, error) {
	res, err := d.sendTCP(`{"smartlife.iot.common.schedule":{"get_rules":{}}}`)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return res, nil
}

// https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents

// when I get bored, set myself up as the cloud server... -- make it as responsive as the shellies
// {"cnCloud":{"set_server_url":{"server":"devs.tplinkcloud.com"}}}
// {"cnCloud":{"bind":{"username":alice@home.com, "password":"secret"}}}
