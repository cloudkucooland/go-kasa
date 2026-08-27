package kasa

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSetRelayStateCtx(t *testing.T) {
	tests := []struct {
		name      string
		input     bool
		wantCmd   string
		shouldErr bool
	}{
		{"turn on", true, fmt.Sprintf(CmdSetRelayState, 1), false},
		{"turn off", false, fmt.Sprintf(CmdSetRelayState, 0), false},
		{"udp error", true, fmt.Sprintf(CmdSetRelayState, 1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideUDP: func(ctx context.Context, cmd string) error {
					if tt.shouldErr {
						return errors.New("udp error")
					}
					if cmd != tt.wantCmd {
						t.Errorf("got cmd %q, want %q", cmd, tt.wantCmd)
					}
					return nil
				},
			}

			err := md.SetRelayStateCtx(context.Background(), tt.input)

			if tt.shouldErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetSettingsCtx(t *testing.T) {
	tests := []struct {
		name      string
		response  string
		shouldErr bool
	}{
		{
			name: "valid response",
			response: `{
				"system": {
					"get_sysinfo": {
						"alias": "plug1",
						"err_code": 0
					}
				}
			}`,
			shouldErr: false,
		},
		{
			name: "kasa error",
			response: `{
				"system": {
					"get_sysinfo": {
						"err_code": -1,
						"err_msg": "failure"
					}
				}
			}`,
			shouldErr: true,
		},
		{
			name:      "invalid json",
			response:  `{invalid}`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideTCP: func(ctx context.Context, cmd string) ([]byte, error) {
					return []byte(tt.response), nil
				},
			}

			res, err := md.GetSettingsCtx(context.Background())

			if tt.shouldErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res == nil {
				t.Fatal("expected result, got nil")
			}
		})
	}
}

func TestGetEmeterCtx(t *testing.T) {
	tests := []struct {
		name      string
		response  string
		shouldErr bool
	}{
		{
			name: "valid",
			response: `{
				"emeter": {
					"get_realtime": {
						"power_mw": 12345,
						"err_code": 0
					}
				}
			}`,
			shouldErr: false,
		},
		{
			name: "device error",
			response: `{
				"emeter": {
					"get_realtime": {
						"err_code": -1
					}
				}
			}`,
			shouldErr: true,
		},
		{
			name:      "invalid json",
			response:  `{invalid}`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideTCP: func(ctx context.Context, cmd string) ([]byte, error) {
					return []byte(tt.response), nil
				},
			}

			res, err := md.GetEmeterCtx(context.Background())

			if tt.shouldErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.shouldErr && res == nil {
				t.Fatal("expected result")
			}
		})
	}
}

func TestSetRelayStateChildMultiCtx(t *testing.T) {
	tests := []struct {
		name     string
		children []string
		state    bool
	}{
		{
			name:     "single child",
			children: []string{"a"},
			state:    true,
		},
		{
			name:     "multiple children",
			children: []string{"a", "b"},
			state:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideUDP: func(ctx context.Context, cmd string) error {
					for _, c := range tt.children {
						if !strings.Contains(cmd, c) {
							t.Fatalf("expected child %q in cmd %q", c, cmd)
						}
					}
					return nil
				},
			}

			err := md.SetRelayStateChildMultiCtx(context.Background(), tt.state, tt.children...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSendRawCommandCtx(t *testing.T) {
	md := &Device{
		OverrideTCP: func(ctx context.Context, cmd string) ([]byte, error) {
			return []byte(`ok`), nil
		},
		OverrideUDP: func(ctx context.Context, cmd string) error {
			return nil
		},
	}

	_, err := md.SendRawCommandCtx(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAliasEscaping(t *testing.T) {
	tests := []struct {
		name     string
		call     func(d *Device) error
		contains []string
		useTCP   bool
	}{
		{
			name:     "SetAlias with quotes",
			call:     func(d *Device) error { return d.SetAliasCtx(context.Background(), `My "Smart" Plug`) },
			contains: []string{`"alias":"My \"Smart\" Plug"`},
		},
		{
			name:     "SetAlias with backslash",
			call:     func(d *Device) error { return d.SetAliasCtx(context.Background(), `C:\Path`) },
			contains: []string{`"alias":"C:\\Path"`},
		},
		{
			name:     "SetChildAlias with quotes",
			call:     func(d *Device) error { return d.SetChildAliasCtx(context.Background(), `id"1`, `Lamp "2"`) },
			contains: []string{`"child_ids":["id\"1"]`, `"alias":"Lamp \"2\""`},
		},
		{
			name: "SetWIFI with special chars",
			call: func(d *Device) error {
				_, err := d.SetWIFICtx(context.Background(), `My "SSID"`, `P@$$w0rd\`)
				return err
			},
			contains: []string{`"ssid":"My \"SSID\""`, `"password":"P@$$w0rd\\"`},
			useTCP:   true,
		},
		{
			name:     "EnableCloud with quotes",
			call:     func(d *Device) error { return d.EnableCloudCtx(context.Background(), `user"name`, `p"ss`) },
			contains: []string{`"username":"user\"name"`, `"password":"p\"ss"`},
		},
		{
			name:     "AddCountdownRule with quotes",
			call:     func(d *Device) error { return d.AddCountdownRuleCtx(context.Background(), 60, true, `Timer "1"`) },
			contains: []string{`"name":"Timer \"1\""`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{}
			if tt.useTCP {
				md.OverrideTCP = func(ctx context.Context, cmd string) ([]byte, error) {
					for _, want := range tt.contains {
						if !strings.Contains(cmd, want) {
							t.Errorf("cmd %q does not contain %q", cmd, want)
						}
					}
					return []byte(`{"netif":{"set_stainfo":{"err_code":0}}}`), nil
				}
			} else {
				md.OverrideUDP = func(ctx context.Context, cmd string) error {
					for _, want := range tt.contains {
						if !strings.Contains(cmd, want) {
							t.Errorf("cmd %q does not contain %q", cmd, want)
						}
					}
					return nil
				}
			}
			err := tt.call(md)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUDPCommands(t *testing.T) {
	tests := []struct {
		name string
		call func(d *Device) error
		want string
	}{
		{
			name: "SetBrightness",
			call: func(d *Device) error { return d.SetBrightnessCtx(context.Background(), 50) },
			want: fmt.Sprintf(CmdSetBrightness, 50),
		},
		{
			name: "SetAlias",
			call: func(d *Device) error { return d.SetAliasCtx(context.Background(), "new-alias") },
			want: `{"system":{"set_dev_alias":{"alias":"new-alias"}}}`,
		},
		{
			name: "SetLEDOff true",
			call: func(d *Device) error { return d.SetLEDOffCtx(context.Background(), true) },
			want: fmt.Sprintf(CmdLEDOff, 1),
		},
		{
			name: "SetLEDOff false",
			call: func(d *Device) error { return d.SetLEDOffCtx(context.Background(), false) },
			want: fmt.Sprintf(CmdLEDOff, 0),
		},
		{
			name: "Reboot",
			call: func(d *Device) error { return d.RebootCtx(context.Background()) },
			want: CmdReboot,
		},
		{
			name: "DisableCloud",
			call: func(d *Device) error { return d.DisableCloudCtx(context.Background()) },
			want: CmdCloudUnbind,
		},
		{
			name: "SetFadeOnTime",
			call: func(d *Device) error { return d.SetFadeOnTimeCtx(context.Background(), 1000) },
			want: fmt.Sprintf(CmdSetFadeOnTime, 1000),
		},
		{
			name: "EnableCloud",
			call: func(d *Device) error { return d.EnableCloudCtx(context.Background(), "user", "pass") },
			want: `{"cnCloud":{"bind":{"username":"user","password":"pass"}}}`,
		},
		{
			name: "ClearCountdownRules",
			call: func(d *Device) error { return d.ClearCountdownRulesCtx(context.Background()) },
			want: CmdDeleteAllRules,
		},
		{
			name: "AddCountdownRule",
			call: func(d *Device) error { return d.AddCountdownRuleCtx(context.Background(), 10, true, "test") },
			want: `{"count_down":{"add_rule":{"enable":1,"delay":10,"act":1,"name":"test"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideUDP: func(ctx context.Context, cmd string) error {
					if cmd != tt.want {
						t.Errorf("got %q, want %q", cmd, tt.want)
					}
					return nil
				},
			}
			err := tt.call(md)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTCPCommands(t *testing.T) {
	tests := []struct {
		name     string
		call     func(d *Device) error
		wantCmd  string
		response string
	}{
		{
			name:     "SetMode",
			call:     func(d *Device) error { return d.SetModeCtx(context.Background(), "none") },
			wantCmd:  fmt.Sprintf(CmdSetMode, "none"),
			response: `{"system":{"get_sysinfo":{"err_code":0}}}`,
		},
		{
			name:     "GetWIFIStatus",
			call:     func(d *Device) error { _, err := d.GetWIFIStatusCtx(context.Background()); return err },
			wantCmd:  CmdWifiStainfo,
			response: `{"netif":{"get_stainfo":{"err_code":0}}}`,
		},
		{
			name:     "GetDimmerParameters",
			call:     func(d *Device) error { _, err := d.GetDimmerParametersCtx(context.Background()); return err },
			wantCmd:  CmdGetDimmer,
			response: `{"smartlife.iot.dimmer":{"get_dimmer_parameters":{"err_code":0}}}`,
		},
		{
			name:     "GetCountdownRules",
			call:     func(d *Device) error { _, err := d.GetCountdownRulesCtx(context.Background()); return err },
			wantCmd:  CmdGetCountdownRules,
			response: `{"count_down":{"get_rules":{"err_code":0}}}`,
		},
		{
			name:     "SetWIFI",
			call:     func(d *Device) error { _, err := d.SetWIFICtx(context.Background(), "ssid", "pass"); return err },
			wantCmd:  `{"netif":{"set_stainfo":{"ssid":"ssid","password":"pass","key_type":4}}}`,
			response: `{"netif":{"set_stainfo":{"err_code":0}}}`,
		},
		{
			name:     "GetLightSensorConfig",
			call:     func(d *Device) error { _, err := d.GetLightSensorConfigCtx(context.Background()); return err },
			wantCmd:  CmdGetLightSensorConfig,
			response: `{"get_config":{"ver":"1.0","err_code":0}}`,
		},
		{
			name:     "GetCurrentBrightness",
			call:     func(d *Device) error { _, err := d.GetCurrentBrightnessCtx(context.Background()); return err },
			wantCmd:  CmdGetCurrentBrightness,
			response: `{"get_current_brt":{"value":10,"err_code":0}}`,
		},
		{
			name:     "GetEmeterMonth",
			call:     func(d *Device) error { _, err := d.GetEmeterMonthCtx(context.Background(), 1, 2025); return err },
			wantCmd:  fmt.Sprintf(CmdEmeterGetMonth, 1, 2025),
			response: `{"emeter":{"get_daystat":{"err_code":0}}}`,
		},
		{
			name:     "GetEmeterChild",
			call:     func(d *Device) error { _, err := d.GetEmeterChildCtx(context.Background(), "child1"); return err },
			wantCmd:  `{"context":{"child_ids":["child1"]},"emeter":{"get_realtime":{}}}`,
			response: `{"emeter":{"get_realtime":{"err_code":0}}}`,
		},
		{
			name: "GetEmeterChildMonth",
			call: func(d *Device) error {
				_, err := d.GetEmeterChildMonthCtx(context.Background(), 1, 2025, "child1")
				return err
			},
			wantCmd:  `{"context":{"child_ids":["child1"]},"emeter":{"get_daystat":{"month":1,"year":2025}}}`,
			response: `{"emeter":{"get_daystat":{"err_code":0}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &Device{
				OverrideTCP: func(ctx context.Context, cmd string) ([]byte, error) {
					if cmd != tt.wantCmd {
						t.Errorf("got %q, want %q", cmd, tt.wantCmd)
					}
					return []byte(tt.response), nil
				},
			}
			err := tt.call(md)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
