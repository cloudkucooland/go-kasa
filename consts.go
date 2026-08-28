package kasa

// https://lib.dr.iastate.edu/cgi/viewcontent.cgi?article=1424&context=creativecomponents
// https://github.com/whitslack/kasa/blob/master/API.md

/* -- found in firmware
smartlife.common.debug
smartlife.iot.LAS
smartlife.iot.PIR
smartlife.iot.sensor_trigger
smartlife.iot.smartpowerstrip.manage
*/

// Request strings
const (
	CmdSetRelayState    = `{"system":{"set_relay_state":{"state":%d}}}` // 0 or 1
	CmdGetSysinfo       = `{"system":{"get_sysinfo":{}}}`
	CmdGetBtnCheckRes   = `{"system":{"get_btn_check_res":{}}}`
	CmdGetTestMode      = `{"system":{"get_test_mode":{}}}`
	CmdReboot           = `{"system":{"reboot":{"delay":2}}}`
	CmdLEDOff           = `{"system":{"set_led_off":{"off":%d}}}` // off = 1, on = 0
	CmdDeviceAlias      = `{"system":{"set_dev_alias":{"alias":"%s"}}}`
	CmdSetMode          = `{"system":{"set_mode":{"mode":"%s"}}}` // "none", "count_down", ???
	CmdGetOnboarding    = `{"system":{"get_onboarding_status":{}}}`
	CmdSetOnboarding    = `{"system":{"set_onboarding_status":{"value":"%s"}}}`
	CmdCheckConfig      = `{"system":{"check_new_config":null}}`
	CmdCheckUboot       = `{"system":{"test_check_uboot":null}}`
	CmdReset            = `{"system":{"reset":{"delay":1}}}`
	CmdSetMAC           = `{"system":{"set_mac_addr":{"mac":"%s"}}}`                                                                            // 50-C7-BF-01-02-03
	CmdSetHWID          = `{"system":{"set_hw_id":{"hwId":%s}}}`                                                                                // "0123456789ABCDEF0123456789ABCDEF"
	CmdSetDevID         = `{"system":{"set_device_id":{"deviceId":%s}}}`                                                                        //  "0123456789ABCDEF0123456789ABCDEF01234567"
	CmdDownloadFirmware = `{"system":{"download_firmware":{"url":%s}}}`                                                                         // http://
	CmdSetLocation      = `{"system":{"device_location_change":{"latitude":%d,"longitude":%d,"latitude_i":%d,"longitude_i":%d,"timezone":%d}}}` //ints

	CmdGetDebug         = `{"smartlife.common.debug":{"get_diagnose_status":{}}}`
	CmdGetMCUDiagnose   = `{"smartlife.common.debug":{"get_mcu_diagnose":{}}}`
	CmdSetLocationFloat = `{"smartlife.iot.common.system": {"device_location_change": {"latitude": %f, "longitude": %f, "timezone": %d }}}` // floats -- untested, no other source found

	CmdGetEmeter           = `{"emeter":{"get_realtime":{}}}`
	CmdGetEmeterGetDaystat = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`
	CmdGetEmeterVgain      = `{"emeter":{"get_vgain_igain":{}}}`
	CmdSetEmeterVgain      = `{"emeter":{"set_vgain_igain":{"vgain":%d,"igain":%d}}}`       // int, int
	CmdEmeterCalibration   = `{"emeter":{"start_calibration":{"vtarget":%d,"itarget":%d}}}` // int, int
	CmdEmeterGetMonth      = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`            // 1-12, 4-digit-year
	CmdEmeterGetYear       = `{"emeter":{"get_monthstat":{"year":%d}}}`                     // 4-digit-year
	CmdEmeterErase         = `{"emeter":{"erase_emeter_stat":null}}`

	CmdGetEmeterChild      = `{"context":{"child_ids":["%s"]},"emeter":{"get_realtime":{}}}`
	CmdGetEmeterMonthChild = `{"context":{"child_ids":["%s"]},"emeter":{"get_daystat":{"month":%d,"year":%d}}}`

	CmdWifiStainfo    = `{"netif":{"get_stainfo":{}}}`
	CmdWifiScanInfo   = `{"netif":{"get_scaninfo":{"refresh":%d}}}`
	CmdWifiSetStainfo = `{"netif":{"set_stainfo":{"ssid":"%s","password":"%s","key_type":%d}}}` // string, string, int

	CmdSetRelayStateChild      = `{"context":{"child_ids":["%s"]},"system":{"set_relay_state":{"state":%d}}}` // index (e.g. ".....00"), 0/1
	CmdSetRelayStateChildMulti = `{"context":{"child_ids":[%s]},"system":{"set_relay_state":{"state":%d}}}`   // indexes (e.g. `"....00","....03"`), 0/1
	CmdChildAlias              = `{"context":{"child_ids":["%s"]},"system":{"set_dev_alias":{"alias":"%s"}}}` // index (e.g. "....01"), name

	CmdGetDimmer        = `{"smartlife.iot.dimmer":{"get_dimmer_parameters":{}}}`
	CmdSetBrightness    = `{"smartlife.iot.dimmer":{"set_brightness":{"brightness":%d}}}`    // 0-100
	CmdSetFadeOffTime   = `{"smartlife.iot.dimmer":{"set_fade_off_time":{"fadeTime":%d}}}`   // ms
	CmdSetFadeOnTime    = `{"smartlife.iot.dimmer":{"set_fade_on_time":{"fadeTime":%d}}}`    // ms
	CmdSetGentleOffTime = `{"smartlife.iot.dimmer":{"set_gentle_off_time":{"fadeTime":%d}}}` // ms
	CmdSetGentleOnTime  = `{"smartlife.iot.dimmer":{"set_gentle_on_time":{"fadeTime":%d}}}`  // ms

	CmdGetCountdownRules = `{"count_down":{"get_rules":{}}}`
	CmdDeleteAllRules    = `{"count_down":{"delete_all_rules":{}}}`
	CmdAddCountdownRule  = `{"count_down":{"add_rule":{"enable":1,"delay":%d,"act":%d,"name":"%s"}}}` // 0-3600, 0/1, string

	CmdCloudUnbind    = `{"cnCloud":{"unbind":null}}`
	CmdGetCloudInfo   = `{"cnCloud":{"get_info":{}}}`
	CmdGetIntlFwList  = `{"cnCloud":{"get_intl_fw_list":null}}`
	CmdSetServerURL   = `{"cnCloud":{"set_server_url":{"server":"%s"}}}`          // bare hostname, no protocol spec
	CmdSetServerCreds = `{"cnCloud":{"bind":{"username":"%s", "password":"%s"}}}` // alice@home.com / mikeisagoat
	CmdGetSefInfo     = `{"cnCloud":{"get_sefinfo": {}}}`
	CmdSetSefURL      = `{"cnCloud":{"set_sefserver_url":{"server":"%s"}}}` // bare hostname, no protocol spec

	CmdGetLightSensorConfig = `{"smartlife.iot.LAS":{"get_config":{}}}`
	CmdGetCurrentBrightness = `{"smartlife.iot.LAS":{"get_current_brt":{}}}`
	CmdSetBrightnessLevel   = `{"smartlife.iot.LAS":{"set_brt_level":{"index":%d,"value":%d}}}` // int, int
	CmdSetDarkIndex         = `{"smartlife.iot.LAS":{"set_dark_index":{"dark_index":%d}}}`      // int
	CmdSetLightSensorEnable = `{"smartlife.iot.LAS":{"set_enable":{"enable":%d}}}`              // 0/1
	CmdGetDark              = `{"smartlife.iot.LAS":{"get_dark_status":{}}}`

	CmdGetPIRConfig      = `{"smartlife.iot.PIR":{"get_config":{}}}`
	CmdSetPIRColdTime    = `{"smartlife.iot.PIR":{"set_cold_time":{"cold_time":%d}}}`           // int
	CmdSetPIREnable      = `{"smartlife.iot.PIR":{"set_enable":{"enable":%d}}}`                 // 0/1
	CmdSetPIRSensitivity = `{"smartlife.iot.PIR":{"set_trigger_sens":{"index":%d,"value":%d}}}` // int, int (~decimeters)
	CmdGetPIRADC         = `{"smartlife.iot.PIR": {"get_adc_value":{}}}`

	CmdGetTime     = `{"time":{"get_time":{}}}`
	CmdGetTimezone = `{"time":{"get_timezone":{}}}`
	CmdSetTimezone = `{"time":{"set_timezone":{"year":%d,"month":%d,"mday":%d,"hour":%d,"min":%d,"sec":%d}}}`

	CmdGetSensorRoutine       = `{"smartlife.iot.sensor_trigger":{"get_weekday_routine":{}}}`
	CmdGetSensorMode          = `{"smartlife.iot.sensor_trigger":{"get_mode":{}}}`
	CmdSetSensorMode          = `{"smartlife.iot.sensor_trigger":{"set_mode":"%s"}}`
	CmdAddSensorRoutine       = `{"smartlife.iot.sensor_trigger":{"edit_weekday_routine":%s}}` // JSON struct string
	CmdDeleteSensorRoutine    = `{"smartlife.iot.sensor_trigger":{"delete_weekday_routine":{"id":"%s"}}}`
	CmdGetDefaultManualAction = `{"smartlife.iot.sensor_trigger":{"get_default_manual_action":{}}}`
	CmdSetDefaultManualAction = `{"smartlife.iot.sensor_trigger":{"set_default_manual_action":{"offToS":%d}}}`

	CmdGetScheduleRules       = `{"schedule":{"get_rules":null}}`
	CmdAddScheduleRule        = `{"schedule":{"edit_rule":%s}}`
	CmdDeleteScheduleRule     = `{"schedule":{"delete_rule":{"id":"%s"}}}`
	CmdDeleteAllScheduleRules = `{"schedule":{"delete_all_rules":null}}`
	CmdSetScheduleEnabled     = `{"schedule":{"set_overall_enable":{"enable":%d}}}`

	CmdGetLightState        = `{"smartlife.iot.smartbulb.lightingservice":{"get_light_state":{}}}`
	CmdTransitionLightState = `{"smartlife.iot.smartbulb.lightingservice":{"transition_light_state":%s}}`
	CmdGetPreferredState    = `{"smartlife.iot.smartbulb.lightingservice":{"get_preferred_state":{}}}`
	CmdSetPreferredState    = `{"smartlife.iot.smartbulb.lightingservice":{"set_preferred_state":%s}}`
	CmdGetDefaultBehavior   = `{"smartlife.iot.smartbulb.lightingservice":{"get_default_behavior":{}}}`
	CmdSetDefaultBehavior   = `{"smartlife.iot.smartbulb.lightingservice":{"set_default_behavior":%s}}`
	CmdGetLightDetails      = `{"smartlife.iot.smartbulb.lightingservice":{"get_light_details":{}}}`
)
