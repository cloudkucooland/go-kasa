package kasa

// https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents

// Request strings
const (
	CmdSetRelayState = `{"system":{"set_relay_state":{"state":%d}}}` // 0 or 1
	CmdGetSysinfo    = `{"system":{"get_sysinfo":{}}}`
	CmdReboot        = `{"system":{"reboot":{"delay":2}}}`
	CmdLEDOff        = `{"system":{"set_led_off":{"off":%d}}}` // off = 1, on = 0
	CmdDeviceAlias   = `{"system":{"set_dev_alias":{"alias":"%s"}}}`
	CmdSetMode       = `{"system":{"set_mode":{"mode":"%s"}}}` // "none", "count_down", ???

	CmdGetEmeter           = `{"emeter":{"get_realtime":{}}}`
	CmdGetEmeterGetDaystat = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`
	CmdGetEmeterVgain      = `{"emeter":{"get_vgain_igain":{}}}`
	CmdSetEmeterVgain      = `{"emeter":{"set_vgain_igain":{"vgain":%d,"igain":%d}}}`       // int, int
	CmdEmeterCalibration   = `{"emeter":{"start_calibration":{"vtarget":%d,"itarget":%d}}}` // int, int
	CmdEmeterGetMonth      = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`            // 1-12, 4-digit-year
	CmdEmeterGetYear       = `{"emeter":{"get_monthstat":{"year":%d}}}`                     // 4-digit-year
	CmdEmeterErase         = `{"emeter":{"erase_emeter_stat":null}}`

	CmdGetEmeterChild = `{"context":{"child_ids":["%s"]},"emeter":{"get_realtime":{}}}`

	CmdWifiStainfo = `{"netif":{"get_stainfo":{}}}`

	CmdSetRelayStateChild      = `{"context":{"child_ids":["%s"]},"system":{"set_relay_state":{"state":%d}}}` // index (e.g. ".....00"), 0/1
	CmdSetRelayStateChildMulti = `{"context":{"child_ids":[%s]},"system":{"set_relay_state":{"state":%d}}}`   // indexes (e.g. `"....00","....03"`), 0/1
	CmdChildAlias              = `{"context":{"child_ids":["%s"]},"system":{"set_dev_alias":{"alias":"%s"}}}` // index (e.g. "....01"), name

	CmdGetDimmer        = `{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`
	CmdSetBrightness    = `{"smartlife.iot.dimmer":{"set_brightness":{"brightness":%d}}}`    // 0-100
	CmdSetFadeOffTime   = `{"smartlife.iot.dimmer":{"set_fade_off_time":{"fadeTime":%d}}}`   // ms
	CmdSetFadeOnTime    = `{"smartlife.iot.dimmer":{"set_fade_on_time":{"fadeTime":%d}}}`    // ms
	CmdSetGentleOffTime = `{"smartlife.iot.dimmer":{"set_gentle_off_time":{"fadeTime":%d}}}` // ms
	CmdSetGentleOnTime  = `{"smartlife.iot.dimmer":{"set_gentle_on_time":{"fadeTime":%d}}}`  // ms

	CmdGetRules = `{"smartlife.iot.common.schedule":{"get_rules":{}}}`

	CmdGetCountdownRules = `{"count_down":{"get_rules":{}}}`
	CmdDeleteAllRules    = `{"count_down":{"delete_all_rules":{}}}`
	CmdAddCountdownRule  = `{"count_down":{"add_rule":{"enable":1,"delay":%d,"act":%d,"name":"%s"}}}` // 0-3600, 0/1, string

	// CmdGetCountdownRules = `{"smartlife.iot.common.count_down":{"get_rules":{}}}`
	// CmdDeleteAllRules    = `{"smartlife.iot.common.count_down":{"delete_all_rules":{}}}`
	// CmdAddCountdownRule  = `{"smartlife.iot.common.count_down":{"add_rule":{"enable":1,"delay":%d,"act":%d,"name":"%s"}}}`

	CmdCloudUnbind    = `{"cnCloud":{"unbind":null}}`
	CmdSetServerURL   = `{"cnCloud":{"set_server_url":{"server":"%s"}}}`          // bare hostname, no protocol spec
	CmdSetServerCreds = `{"cnCloud":{"bind":{"username":"%s", "password":"%s"}}}` // alice@home.com / mikeisagoat
)
