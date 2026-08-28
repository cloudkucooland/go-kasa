package kasa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// SetRelayState changes the relay state of the device -- for multi-relay devices use SetRelayStateChild
func (d *Device) SetRelayState(newstate bool) error {
	return d.SetRelayStateCtx(context.Background(), newstate)
}

func (d *Device) SetRelayStateCtx(ctx context.Context, newstate bool) error {
	cmd := fmt.Sprintf(CmdSetRelayState, boolToInt(newstate))
	return d.sendUDP(ctx, cmd)
}

// SetRelayStateChild adjusts a single relay on a multi-relay device
func (d *Device) SetRelayStateChild(childID string, newstate bool) error {
	return d.SetRelayStateChildCtx(context.Background(), childID, newstate)
}

func (d *Device) SetRelayStateChildCtx(ctx context.Context, childID string, newstate bool) error {
	type cmd struct {
		Context struct {
			ChildIDs []string `json:"child_ids"`
		} `json:"context"`
		System struct {
			SetRelayState struct {
				State int `json:"state"`
			} `json:"set_relay_state"`
		} `json:"system"`
	}
	var c cmd
	c.Context.ChildIDs = []string{childID}
	c.System.SetRelayState.State = boolToInt(newstate)
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

// SetRelayStateChildMulti adjusts multiple relays on a multi-relay device
func (d *Device) SetRelayStateChildMulti(newstate bool, children ...string) error {
	return d.SetRelayStateChildMultiCtx(context.Background(), newstate, children...)
}

func (d *Device) SetRelayStateChildMultiCtx(ctx context.Context, newstate bool, children ...string) error {
	type cmd struct {
		Context struct {
			ChildIDs []string `json:"child_ids"`
		} `json:"context"`
		System struct {
			SetRelayState struct {
				State int `json:"state"`
			} `json:"set_relay_state"`
		} `json:"system"`
	}
	var c cmd
	c.Context.ChildIDs = children
	c.System.SetRelayState.State = boolToInt(newstate)
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

func (d *Device) SendRawCommand(cmd string) ([]byte, error) {
	return d.SendRawCommandCtx(context.Background(), cmd)
}

func (d *Device) SendRawCommandCtx(ctx context.Context, cmd string) ([]byte, error) {
	result, err := d.sendTCP(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SetBrightness adjust the brightness setting on a dimmer-capable device (1-100)
func (d *Device) SetBrightness(newval int) error {
	return d.SetBrightnessCtx(context.Background(), newval)
}

func (d *Device) SetBrightnessCtx(ctx context.Context, newval int) error {
	cmd := fmt.Sprintf(CmdSetBrightness, newval)
	return d.sendUDP(ctx, cmd)
}

func (d *Device) SetFadeOffTime(newval int) error {
	return d.SetFadeOffTimeCtx(context.Background(), newval)
}

func (d *Device) SetFadeOffTimeCtx(ctx context.Context, newval int) error {
	cmd := fmt.Sprintf(CmdSetFadeOffTime, newval)
	return d.sendUDP(ctx, cmd)
}

func (d *Device) SetFadeOnTime(newval int) error {
	return d.SetFadeOnTimeCtx(context.Background(), newval)
}

func (d *Device) SetFadeOnTimeCtx(ctx context.Context, newval int) error {
	cmd := fmt.Sprintf(CmdSetFadeOnTime, newval)
	return d.sendUDP(ctx, cmd)
}

func (d *Device) SetGentleOffTime(newval int) error {
	return d.SetGentleOffTimeCtx(context.Background(), newval)
}

func (d *Device) SetGentleOffTimeCtx(ctx context.Context, newval int) error {
	cmd := fmt.Sprintf(CmdSetGentleOffTime, newval)
	return d.sendUDP(ctx, cmd)
}

func (d *Device) SetGentleOnTime(newval int) error {
	return d.SetGentleOnTimeCtx(context.Background(), newval)
}

func (d *Device) SetGentleOnTimeCtx(ctx context.Context, newval int) error {
	cmd := fmt.Sprintf(CmdSetGentleOnTime, newval)
	return d.sendUDP(ctx, cmd)
}

// GetSettings gets the device sys info
func (d *Device) GetSettings() (*Sysinfo, error) {
	return d.GetSettingsCtx(context.Background())
}

func (d *Device) GetSettingsCtx(ctx context.Context) (*Sysinfo, error) {
	res, err := d.sendTCP(ctx, CmdGetSysinfo)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.GetSysinfo.Sysinfo.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.GetSysinfo.Sysinfo, nil
}

// GetBtnCheckRes gets the device button status
func (d *Device) GetBtnCheckRes() (*BtnCheckRes, error) {
	return d.GetBtnCheckResCtx(context.Background())
}

func (d *Device) GetBtnCheckResCtx(ctx context.Context) (*BtnCheckRes, error) {
	res, err := d.sendTCP(ctx, CmdGetBtnCheckRes)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.GetSysinfo.BtnCheck.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.GetSysinfo.BtnCheck, nil
}

// GetTestMode gets the device test mode status
func (d *Device) GetTestMode() (*TestModeRes, error) {
	return d.GetTestModeCtx(context.Background())
}

func (d *Device) GetTestModeCtx(ctx context.Context) (*TestModeRes, error) {
	res, err := d.sendTCP(ctx, CmdGetTestMode)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.GetSysinfo.TestMode.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.GetSysinfo.TestMode, nil
}

// GetOnboardingStatus returns onboarding status
func (d *Device) GetOnboardingStatus() (*OnboardingStatus, error) {
	return d.GetOnboardingStatusCtx(context.Background())
}

func (d *Device) GetOnboardingStatusCtx(ctx context.Context) (*OnboardingStatus, error) {
	res, err := d.sendTCP(ctx, CmdGetOnboarding)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.GetSysinfo.Onboarding.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.GetSysinfo.Onboarding, nil
}

// GetEmeter returns emeter data from the device
func (d *Device) GetEmeter() (*EmeterRealtime, error) {
	return d.GetEmeterCtx(context.Background())
}

func (d *Device) GetEmeterCtx(ctx context.Context) (*EmeterRealtime, error) {
	res, err := d.sendTCP(ctx, CmdGetEmeter)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Emeter.Realtime.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Emeter.Realtime, nil
}

// GetEmeterMonth returns a single month's emeter data from the device
func (d *Device) GetEmeterMonth(month, year int) (*EmeterDaystat, error) {
	return d.GetEmeterMonthCtx(context.Background(), month, year)
}

func (d *Device) GetEmeterMonthCtx(ctx context.Context, month, year int) (*EmeterDaystat, error) {
	q := fmt.Sprintf(CmdEmeterGetMonth, month, year)

	res, err := d.sendTCP(ctx, q)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Emeter.DayStat.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Emeter.DayStat, nil
}

// GetEmeter returns emeter data from the device
func (d *Device) GetEmeterChild(child string) (*EmeterRealtime, error) {
	return d.GetEmeterChildCtx(context.Background(), child)
}

func (d *Device) GetEmeterChildCtx(ctx context.Context, child string) (*EmeterRealtime, error) {
	type cmd struct {
		Context struct {
			ChildIDs []string `json:"child_ids"`
		} `json:"context"`
		Emeter struct {
			GetRealtime struct{} `json:"get_realtime"`
		} `json:"emeter"`
	}
	var c cmd
	c.Context.ChildIDs = []string{child}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	res, err := d.sendTCP(ctx, string(b))
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Emeter.Realtime.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Emeter.Realtime, nil
}

func (d *Device) GetEmeterChildMonth(month int, year int, child string) (*EmeterDaystat, error) {
	return d.GetEmeterChildMonthCtx(context.Background(), month, year, child)
}

func (d *Device) GetEmeterChildMonthCtx(ctx context.Context, month int, year int, child string) (*EmeterDaystat, error) {
	type cmd struct {
		Context struct {
			ChildIDs []string `json:"child_ids"`
		} `json:"context"`
		Emeter struct {
			GetDaystat struct {
				Month int `json:"month"`
				Year  int `json:"year"`
			} `json:"get_daystat"`
		} `json:"emeter"`
	}
	var c cmd
	c.Context.ChildIDs = []string{child}
	c.Emeter.GetDaystat.Month = month
	c.Emeter.GetDaystat.Year = year
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	res, err := d.sendTCP(ctx, string(b))
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Emeter.DayStat.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Emeter.DayStat, nil
}

// DisableCloud sets the device to "local only" mode.
// TODO: forget any cloud settings
func (d *Device) DisableCloud() error {
	return d.DisableCloudCtx(context.Background())
}

func (d *Device) DisableCloudCtx(ctx context.Context) error {
	return d.sendUDP(ctx, CmdCloudUnbind)
}

// Enable/Configure Cloud
func (d *Device) EnableCloud(username, password string) error {
	return d.EnableCloudCtx(context.Background(), username, password)
}

func (d *Device) EnableCloudCtx(ctx context.Context, username, password string) error {
	type cmd struct {
		CNCloud struct {
			Bind struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"bind"`
		} `json:"cnCloud"`
	}
	var c cmd
	c.CNCloud.Bind.Username = username
	c.CNCloud.Bind.Password = password
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

// Reboot instructs the device to reboot
func (d *Device) Reboot() error {
	return d.RebootCtx(context.Background())
}

func (d *Device) RebootCtx(ctx context.Context) error {
	return d.sendUDP(ctx, CmdReboot)
}

// SetLEDOff is insanely named... it should be SetLED, but I'm just going with what TP-Link called these things internally...
func (d *Device) SetLEDOff(t bool) error {
	return d.SetLEDOffCtx(context.Background(), t)
}

func (d *Device) SetLEDOffCtx(ctx context.Context, t bool) error {
	cmd := fmt.Sprintf(CmdLEDOff, boolToInt(t))
	return d.sendUDP(ctx, cmd)
}

// SetAlias sets a device name
func (d *Device) SetAlias(s string) error {
	return d.SetAliasCtx(context.Background(), s)
}

func (d *Device) SetAliasCtx(ctx context.Context, s string) error {
	type cmd struct {
		System struct {
			SetDevAlias struct {
				Alias string `json:"alias"`
			} `json:"set_dev_alias"`
		} `json:"system"`
	}
	var c cmd
	c.System.SetDevAlias.Alias = s
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

// SetChildAlias sets the name of an individual relay on a multi-relay device, I don't think this works
func (d *Device) SetChildAlias(childID, s string) error {
	return d.SetChildAliasCtx(context.Background(), childID, s)
}

func (d *Device) SetChildAliasCtx(ctx context.Context, childID, s string) error {
	type cmd struct {
		Context struct {
			ChildIDs []string `json:"child_ids"`
		} `json:"context"`
		System struct {
			SetDevAlias struct {
				Alias string `json:"alias"`
			} `json:"set_dev_alias"`
		} `json:"system"`
	}
	var c cmd
	c.Context.ChildIDs = []string{childID}
	c.System.SetDevAlias.Alias = s
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

// SetMode sets the target mode of the system
func (d *Device) SetMode(m string) error {
	return d.SetModeCtx(context.Background(), m)
}

func (d *Device) SetModeCtx(ctx context.Context, m string) error {
	cmd := fmt.Sprintf(CmdSetMode, m)
	res, err := d.sendTCP(ctx, cmd)
	if err != nil {
		return err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return err
	}

	if err := kd.GetSysinfo.Sysinfo.KasaErr.OK(); err != nil {
		return err
	}

	return nil
}

// GetWIFIStatus returns the WiFi station info
func (d *Device) GetWIFIStatus() (*StaInfo, error) {
	return d.GetWIFIStatusCtx(context.Background())
}

func (d *Device) GetWIFIStatusCtx(ctx context.Context) (*StaInfo, error) {
	res, err := d.sendTCP(ctx, CmdWifiStainfo)
	if err != nil {
		return nil, err
	}

	var ksta KasaDevice
	if err := json.Unmarshal(res, &ksta); err != nil {
		return nil, err
	}

	if err := ksta.NetIf.StaInfo.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &ksta.NetIf.StaInfo, nil
}

// GetScanInfo returns the WiFi scan info
func (d *Device) GetScanInfo(refresh int) (*ScanInfo, error) {
	return d.GetScanInfoCtx(context.Background(), refresh)
}

func (d *Device) GetScanInfoCtx(ctx context.Context, refresh int) (*ScanInfo, error) {
	res, err := d.sendTCP(ctx, fmt.Sprintf(CmdWifiScanInfo, refresh))
	if err != nil {
		return nil, err
	}

	var ksta KasaDevice
	if err := json.Unmarshal(res, &ksta); err != nil {
		return nil, err
	}

	if err := ksta.NetIf.ScanInfo.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &ksta.NetIf.ScanInfo, nil
}

// SetWIFI configures the WiFi station info
func (d *Device) SetWIFI(ssid string, key string) (*SetStaInfo, error) {
	return d.SetWIFICtx(context.Background(), ssid, key)
}

func (d *Device) SetWIFICtx(ctx context.Context, ssid string, key string) (*SetStaInfo, error) {
	if ssid == "" {
		return nil, fmt.Errorf("no ssid specified")
	}
	if key == "" {
		return nil, fmt.Errorf("no key specified")
	}

	type cmd struct {
		NetIf struct {
			SetStainfo struct {
				SSID     string `json:"ssid"`
				Password string `json:"password"`
				KeyType  int    `json:"key_type"`
			} `json:"set_stainfo"`
		} `json:"netif"`
	}
	var c cmd
	c.NetIf.SetStainfo.SSID = ssid
	c.NetIf.SetStainfo.Password = key
	c.NetIf.SetStainfo.KeyType = 4
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	res, err := d.sendTCP(ctx, string(b))
	if err != nil {
		return nil, err
	}

	var ksta KasaDevice
	if err := json.Unmarshal(res, &ksta); err != nil {
		return nil, err
	}

	if err := ksta.NetIf.SetStaInfo.KasaErr.OK(); err != nil {
		return nil, err
	}
	if err := ksta.NetIf.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &ksta.NetIf.SetStaInfo, nil
}

// GetDimmerParameters returns the dimmer parameters from dimmer-capable devices
func (d *Device) GetDimmerParameters() (*DimmerParameters, error) {
	return d.GetDimmerParametersCtx(context.Background())
}

func (d *Device) GetDimmerParametersCtx(ctx context.Context) (*DimmerParameters, error) {
	res, err := d.sendTCP(ctx, CmdGetDimmer)
	if err != nil {
		return nil, err
	}
	var kd KasaDevice
	if err := json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Dimmer.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Dimmer.Parameters, nil
}

// GetRules returns the rule information from a device
/*
func (d *Device) GetRules() (string, error) {
	return d.GetRulesCtx(context.Background())
}

func (d *Device) GetRulesCtx(ctx context.Context) (string, error) {
	res, err := d.sendTCP(ctx, CmdGetRules)
	if err != nil {
		return nil, err
	}
	var kd KasaDevice
	if err := json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Rules.KasaErr.OK(); err != nil {
		return nil, err
	}
	return kd.Rules, err
} */

// GetCountdownRules returns a list of the countdown timers on a device
func (d *Device) GetCountdownRules() ([]Rule, error) {
	return d.GetCountdownRulesCtx(context.Background())
}

func (d *Device) GetCountdownRulesCtx(ctx context.Context) ([]Rule, error) {
	res, err := d.sendTCP(ctx, CmdGetCountdownRules)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Countdown.GetRules.OK(); err != nil {
		return nil, err
	}

	return kd.Countdown.GetRules.RuleList, nil
}

// ClearCountdownRules resets all countdown rules on the device
func (d *Device) ClearCountdownRules() error {
	return d.ClearCountdownRulesCtx(context.Background())
}

func (d *Device) ClearCountdownRulesCtx(ctx context.Context) error {
	return d.sendUDP(ctx, CmdDeleteAllRules)
}

// AddCountdownRule adds a new countdown
func (d *Device) AddCountdownRule(dur int, target bool, name string) error {
	return d.AddCountdownRuleCtx(context.Background(), dur, target, name)
}

func (d *Device) AddCountdownRuleCtx(ctx context.Context, dur int, target bool, name string) error {
	type cmd struct {
		Countdown struct {
			AddRule struct {
				Enable int    `json:"enable"`
				Delay  int    `json:"delay"`
				Act    int    `json:"act"`
				Name   string `json:"name"`
			} `json:"add_rule"`
		} `json:"count_down"`
	}
	var c cmd
	c.Countdown.AddRule.Enable = 1
	c.Countdown.AddRule.Delay = dur
	c.Countdown.AddRule.Act = boolToInt(target)
	c.Countdown.AddRule.Name = name
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return d.sendUDP(ctx, string(b))
}

func (d *Device) GetLightSensorConfig() (*LightSensorConfig, error) {
	return d.GetLightSensorConfigCtx(context.Background())
}

func (d *Device) GetLightSensorConfigCtx(ctx context.Context) (*LightSensorConfig, error) {
	res, err := d.sendTCP(ctx, CmdGetLightSensorConfig)
	if err != nil {
		return nil, err
	}

	var ls LightSensor
	if err = json.Unmarshal(res, &ls); err != nil {
		return nil, err
	}
	if err := ls.GetConfig.OK(); err != nil {
		return nil, err
	}
	if ls.GetConfig.Version == "" {
		return nil, errors.New("light sensor module not present")
	}
	return &ls.GetConfig, nil
}

func (d *Device) GetCurrentBrightness() (uint, error) {
	return d.GetCurrentBrightnessCtx(context.Background())
}

func (d *Device) GetCurrentBrightnessCtx(ctx context.Context) (uint, error) {
	res, err := d.sendTCP(ctx, CmdGetCurrentBrightness)
	if err != nil {
		return 0, err
	}

	var ls LightSensor
	if err = json.Unmarshal(res, &ls); err != nil {
		return 0, err
	}

	if err := ls.GetBrightness.OK(); err != nil {
		return 0, err
	}
	return ls.GetBrightness.Value, nil
}

// GetCloudInfo returns cloud configuration data
func (d *Device) GetCloudInfo() (*CloudInfo, error) {
	return d.GetCloudInfoCtx(context.Background())
}

func (d *Device) GetCloudInfoCtx(ctx context.Context) (*CloudInfo, error) {
	res, err := d.sendTCP(ctx, CmdGetCloudInfo)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.CNCloud.Info.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.CNCloud.Info, nil
}

// GetSefInfo returns cloud SEF configuration data
func (d *Device) GetSefInfo() (*SefInfo, error) {
	return d.GetSefInfoCtx(context.Background())
}

func (d *Device) GetSefInfoCtx(ctx context.Context) (*SefInfo, error) {
	res, err := d.sendTCP(ctx, CmdGetSefInfo)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.CNCloud.SefInfo.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.CNCloud.SefInfo, nil
}

// SetOnboardingStatus sets onboarding status
func (d *Device) SetOnboardingStatus(status string) error {
	return d.SetOnboardingStatusCtx(context.Background(), status)
}

func (d *Device) SetOnboardingStatusCtx(ctx context.Context, status string) error {
	res, err := d.sendTCP(ctx, fmt.Sprintf(CmdSetOnboarding, status))
	if err != nil {
		return err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return err
	}

	return kd.GetSysinfo.SetOnb.KasaErr.OK()
}

// GetDiagnoseStatus returns diagnose status
func (d *Device) GetDiagnoseStatus() (*DiagnoseResult, error) {
	return d.GetDiagnoseStatusCtx(context.Background())
}

func (d *Device) GetDiagnoseStatusCtx(ctx context.Context) (*DiagnoseResult, error) {
	res, err := d.sendTCP(ctx, CmdGetDebug)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Debug.Status.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Debug.Status.Result, nil
}

// GetMCUDiagnose returns mcu diagnose status
func (d *Device) GetMCUDiagnose() (*MCUDiagnose, error) {
	return d.GetMCUDiagnoseCtx(context.Background())
}

func (d *Device) GetMCUDiagnoseCtx(ctx context.Context) (*MCUDiagnose, error) {
	res, err := d.sendTCP(ctx, CmdGetMCUDiagnose)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Debug.MCUDiagnose.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Debug.MCUDiagnose, nil
}

// GetPIRConfig returns PIR configuration
func (d *Device) GetPIRConfig() (*PIRSensorConfig, error) {
	return d.GetPIRConfigCtx(context.Background())
}

func (d *Device) GetPIRConfigCtx(ctx context.Context) (*PIRSensorConfig, error) {
	res, err := d.sendTCP(ctx, CmdGetPIRConfig)
	if err != nil {
		return nil, err
	}

	var p PIRSensor
	if err = json.Unmarshal(res, &p); err != nil {
		return nil, err
	}

	if err := p.GetConfig.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &p.GetConfig, nil
}

// SetPIREnable enables/disables PIR
func (d *Device) SetPIREnable(enable bool) error {
	return d.SetPIREnableCtx(context.Background(), enable)
}

func (d *Device) SetPIREnableCtx(ctx context.Context, enable bool) error {
	cmd := fmt.Sprintf(CmdSetPIREnable, boolToInt(enable))
	return d.sendUDP(ctx, cmd)
}

// SetPIRColdTime sets PIR cold time
func (d *Device) SetPIRColdTime(t int) error {
	return d.SetPIRColdTimeCtx(context.Background(), t)
}

func (d *Device) SetPIRColdTimeCtx(ctx context.Context, t int) error {
	cmd := fmt.Sprintf(CmdSetPIRColdTime, t)
	return d.sendUDP(ctx, cmd)
}

// GetFirmwareList returns available firmware updates
func (d *Device) GetFirmwareList() ([]Firmware, error) {
	return d.GetFirmwareListCtx(context.Background())
}

func (d *Device) GetFirmwareListCtx(ctx context.Context) ([]Firmware, error) {
	res, err := d.sendTCP(ctx, CmdGetIntlFwList)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.CNCloud.FwList.KasaErr.OK(); err != nil {
		return nil, err
	}

	return kd.CNCloud.FwList.List, nil
}

// GetTime returns the current device time
func (d *Device) GetTime() (*TimeData, error) {
	return d.GetTimeCtx(context.Background())
}

func (d *Device) GetTimeCtx(ctx context.Context) (*TimeData, error) {
	res, err := d.sendTCP(ctx, CmdGetTime)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Time.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Time.Time, nil
}

// GetTimezone returns the current device timezone configuration
func (d *Device) GetTimezone() (*TimezoneData, error) {
	return d.GetTimezoneCtx(context.Background())
}

func (d *Device) GetTimezoneCtx(ctx context.Context) (*TimezoneData, error) {
	res, err := d.sendTCP(ctx, CmdGetTimezone)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Time.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Time.Timezone, nil
}

// SetTimezone sets the device time
func (d *Device) SetTimezone(year, month, mday, hour, min, sec int) error {
	return d.SetTimezoneCtx(context.Background(), year, month, mday, hour, min, sec)
}

func (d *Device) SetTimezoneCtx(ctx context.Context, year, month, mday, hour, min, sec int) error {
	cmd := fmt.Sprintf(CmdSetTimezone, year, month, mday, hour, min, sec)
	return d.sendUDP(ctx, cmd)
}

// GetSensorRoutines returns all routines
func (d *Device) GetSensorRoutines() ([]Routine, error) {
	return d.GetSensorRoutinesCtx(context.Background())
}

func (d *Device) GetSensorRoutinesCtx(ctx context.Context) ([]Routine, error) {
	res, err := d.sendTCP(ctx, CmdGetSensorRoutine)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.SensorTrig.KasaErr.OK(); err != nil {
		return nil, err
	}

	return kd.SensorTrig.Routines.RoutineList, nil
}

// DeleteSensorRoutine removes a routine by ID
func (d *Device) DeleteSensorRoutine(id string) error {
	return d.DeleteSensorRoutineCtx(context.Background(), id)
}

func (d *Device) DeleteSensorRoutineCtx(ctx context.Context, id string) error {
	cmd := fmt.Sprintf(CmdDeleteSensorRoutine, id)
	return d.sendUDP(ctx, cmd)
}

// GetDefaultManualAction returns the default action configuration
func (d *Device) GetDefaultManualAction() (int, error) {
	return d.GetDefaultManualActionCtx(context.Background())
}

func (d *Device) GetDefaultManualActionCtx(ctx context.Context) (int, error) {
	res, err := d.sendTCP(ctx, CmdGetDefaultManualAction)
	if err != nil {
		return 0, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return 0, err
	}

	if err := kd.SensorTrig.KasaErr.OK(); err != nil {
		return 0, err
	}

	return kd.SensorTrig.Manual.OffToS, nil
}

// GetSensorMode returns the sensor mode
func (d *Device) GetSensorMode() (*SensorMode, error) {
	return d.GetSensorModeCtx(context.Background())
}

func (d *Device) GetSensorModeCtx(ctx context.Context) (*SensorMode, error) {
	res, err := d.sendTCP(ctx, CmdGetSensorMode)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.SensorTrig.Mode.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.SensorTrig.Mode, nil
}

// SetSensorMode sets the sensor mode
func (d *Device) SetSensorMode(mode string) error {
	return d.SetSensorModeCtx(context.Background(), mode)
}

func (d *Device) SetSensorModeCtx(ctx context.Context, mode string) error {
	cmd := fmt.Sprintf(CmdSetSensorMode, mode)
	return d.sendUDP(ctx, cmd)
}

// SetDefaultManualAction sets the default action configuration
func (d *Device) SetDefaultManualAction(offToS int) error {
	return d.SetDefaultManualActionCtx(context.Background(), offToS)
}

func (d *Device) SetDefaultManualActionCtx(ctx context.Context, offToS int) error {
	cmd := fmt.Sprintf(CmdSetDefaultManualAction, offToS)
	return d.sendUDP(ctx, cmd)
}

// GetScheduleRules returns all schedule rules
func (d *Device) GetScheduleRules() ([]SchedRule, error) {
	return d.GetScheduleRulesCtx(context.Background())
}

func (d *Device) GetScheduleRulesCtx(ctx context.Context) ([]SchedRule, error) {
	res, err := d.sendTCP(ctx, CmdGetScheduleRules)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Schedule.KasaErr.OK(); err != nil {
		return nil, err
	}

	return kd.Schedule.Rules.RuleList, nil
}

// DeleteScheduleRule removes a rule by ID
func (d *Device) DeleteScheduleRule(id string) error {
	return d.DeleteScheduleRuleCtx(context.Background(), id)
}

func (d *Device) DeleteScheduleRuleCtx(ctx context.Context, id string) error {
	cmd := fmt.Sprintf(CmdDeleteScheduleRule, id)
	return d.sendUDP(ctx, cmd)
}

// DeleteAllScheduleRules removes all rules
func (d *Device) DeleteAllScheduleRules() error {
	return d.DeleteAllScheduleRulesCtx(context.Background())
}

func (d *Device) DeleteAllScheduleRulesCtx(ctx context.Context) error {
	return d.sendUDP(ctx, CmdDeleteAllScheduleRules)
}

// SetScheduleEnabled enables or disables the schedule
func (d *Device) SetScheduleEnabled(enable bool) error {
	return d.SetScheduleEnabledCtx(context.Background(), enable)
}

func (d *Device) SetScheduleEnabledCtx(ctx context.Context, enable bool) error {
	cmd := fmt.Sprintf(CmdSetScheduleEnabled, boolToInt(enable))
	return d.sendUDP(ctx, cmd)
}

// GetLightState returns current light state
func (d *Device) GetLightState() (*LightState, error) {
	return d.GetLightStateCtx(context.Background())
}

func (d *Device) GetLightStateCtx(ctx context.Context) (*LightState, error) {
	res, err := d.sendTCP(ctx, CmdGetLightState)
	if err != nil {
		return nil, err
	}

	var kd KasaDevice
	if err = json.Unmarshal(res, &kd); err != nil {
		return nil, err
	}

	if err := kd.Bulb.State.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &kd.Bulb.State, nil
}

// TransitionLightState updates light settings
func (d *Device) TransitionLightState(onOff, brightness, hue, saturation, colorTemp, transitionPeriod, ignoreDefault int) error {
	return d.TransitionLightStateCtx(context.Background(), onOff, brightness, hue, saturation, colorTemp, transitionPeriod, ignoreDefault)
}

func (d *Device) TransitionLightStateCtx(ctx context.Context, onOff, brightness, hue, saturation, colorTemp, transitionPeriod, ignoreDefault int) error {
	type transitionState struct {
		OnOff            int    `json:"on_off"`
		Mode             string `json:"mode"`
		Hue              int    `json:"hue"`
		Saturation       int    `json:"saturation"`
		ColorTemp        int    `json:"color_temp"`
		Brightness       int    `json:"brightness"`
		TransitionPeriod int    `json:"transition_period"`
		IgnoreDefault    int    `json:"ignore_default"`
	}
	ts := transitionState{
		OnOff:            onOff,
		Mode:             "normal",
		Hue:              hue,
		Saturation:       saturation,
		ColorTemp:        colorTemp,
		Brightness:       brightness,
		TransitionPeriod: transitionPeriod,
		IgnoreDefault:    ignoreDefault,
	}
	b, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf(CmdTransitionLightState, string(b))
	return d.sendUDP(ctx, cmd)
}
