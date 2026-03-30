package kasa

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	cmd := fmt.Sprintf(CmdSetRelayStateChild, childID, boolToInt(newstate))
	return d.sendUDP(ctx, cmd)
}

// SetRelayStateChildMulti adjusts multiple relays on a multi-relay device
func (d *Device) SetRelayStateChildMulti(newstate bool, children ...string) error {
	return d.SetRelayStateChildMultiCtx(context.Background(), newstate, children...)
}

func (d *Device) SetRelayStateChildMultiCtx(ctx context.Context, newstate bool, children ...string) error {
	quoted := make([]string, len(children))
	for i, c := range children {
		quoted[i] = `"` + c + `"`
	}

	cmd := fmt.Sprintf(CmdSetRelayStateChildMulti, strings.Join(quoted, ","), boolToInt(newstate))
	return d.sendUDP(ctx, cmd)
}

func (d *Device) SendRawCommand(cmd string) error {
	return d.SendRawCommandCtx(context.Background(), cmd)
}

func (d *Device) SendRawCommandCtx(ctx context.Context, cmd string) error {
	_, err := d.sendTCP(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
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
	q := fmt.Sprintf(CmdGetEmeterChild, child)

	res, err := d.sendTCP(ctx, q)
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
	q := fmt.Sprintf(CmdGetEmeterMonthChild, child, month, year)

	res, err := d.sendTCP(ctx, q)
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
	cmd := fmt.Sprintf(CmdSetServerCreds, username, password)
	return d.sendUDP(ctx, cmd)
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
	cmd := fmt.Sprintf(CmdDeviceAlias, s)
	return d.sendUDP(ctx, cmd)
}

// SetChildAlias sets the name of an individual relay on a multi-relay device, I don't think this works
func (d *Device) SetChildAlias(childID, s string) error {
	return d.SetChildAliasCtx(context.Background(), childID, s)
}

func (d *Device) SetChildAliasCtx(ctx context.Context, childID, s string) error {
	cmd := fmt.Sprintf(CmdChildAlias, childID, s)
	return d.sendUDP(ctx, cmd)
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
	var ksta StaInfo
	if err := json.Unmarshal(res, &ksta); err != nil {
		return nil, err
	}

	if err := ksta.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &ksta, nil
}

// SetWIFI configures the WiFi station info
func (d *Device) SetWIFI(ssid string, key string) (*StaInfo, error) {
	return d.SetWIFICtx(context.Background(), ssid, key)
}

func (d *Device) SetWIFICtx(ctx context.Context, ssid string, key string) (*StaInfo, error) {
	cmd := fmt.Sprintf(CmdWifiSetStainfo, ssid, key, 3)
	res, err := d.sendTCP(ctx, cmd)
	if err != nil {
		return nil, err
	}
	var ksta StaInfo
	if err := json.Unmarshal(res, &ksta); err != nil {
		return nil, err
	}

	if err := ksta.KasaErr.OK(); err != nil {
		return nil, err
	}

	return &ksta, nil
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
func (d *Device) GetRules() (string, error) {
	return d.GetRulesCtx(context.Background())
}

func (d *Device) GetRulesCtx(ctx context.Context) (string, error) {
	res, err := d.sendTCP(ctx, CmdGetRules)
	return string(res), err
}

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
	cmd := fmt.Sprintf(CmdAddCountdownRule, dur, boolToInt(target), name)
	return d.sendUDP(ctx, cmd)
}
