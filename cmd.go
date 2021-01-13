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

func (d *Device) GetSettings() (*kasaSysinfo, error) {
	res, err := d.sendTCP(sysinfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var kd kasaDevice
	if err = json.Unmarshal([]byte(res), &kd); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &kd.System.Sysinfo, nil
}

// forget any cloud settings
func (d *Device) DisableCloud() error {
	err := d.sendUDP(`{"cnCloud":{"unbind":null}}`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// when I get bored, set myself up as the cloud server... -- make it as responsive as the shellies
// {"cnCloud":{"set_server_url":{"server":"devs.tplinkcloud.com"}}}
// {"cnCloud":{"bind":{"username":alice@home.com, "password":"secret"}}}
