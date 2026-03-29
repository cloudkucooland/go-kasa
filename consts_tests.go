package kasa

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCmdFormatting(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		args []any
	}{
		{"SetRelayState On", CmdSetRelayState, []any{1}},
		{"SetRelayState Off", CmdSetRelayState, []any{0}},
		{"Device Alias", CmdDeviceAlias, []any{"Living Room"}},
		{"SetMode CountDown", CmdSetMode, []any{"count_down"}},
		{"LED Off", CmdLEDOff, []any{1}},
		{"LED On", CmdLEDOff, []any{0}},
		{"Child Alias", CmdChildAlias, []any{"01", "Kitchen"}},
		{"SetRelayStateChild", CmdSetRelayStateChild, []any{"01", 1}},
		{"SetRelayStateChildMulti", CmdSetRelayStateChildMulti, []any{`"01","02"`, 0}},
		{"AddCountdownRule", CmdAddCountdownRule, []any{300, 1, "Timer1"}},
		{"SetServerCreds", CmdSetServerCreds, []any{"alice@home.com", "password123"}},
		{"SetServerURL", CmdSetServerURL, []any{"cloud.kasaplugin.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdStr := fmt.Sprintf(tt.cmd, tt.args...)
			if len(cmdStr) == 0 {
				t.Fatalf("Formatted command is empty")
			}

			// Quick JSON validation
			var js map[string]any
			if err := json.Unmarshal([]byte(cmdStr), &js); err != nil {
				t.Fatalf("Invalid JSON produced: %v", err)
			}
		})
	}
}

func TestCmdEmeterMonthYearFormatting(t *testing.T) {
	cmd := fmt.Sprintf(CmdEmeterGetMonth, 7, 2025)
	expected := `{"emeter":{"get_daystat":{"month":7,"year":2025}}}`
	if cmd != expected {
		t.Errorf("Got %q, want %q", cmd, expected)
	}

	cmdYear := fmt.Sprintf(CmdEmeterGetYear, 2025)
	expectedYear := `{"emeter":{"get_monthstat":{"year":2025}}}`
	if cmdYear != expectedYear {
		t.Errorf("Got %q, want %q", cmdYear, expectedYear)
	}
}

func TestCmdWifiFormatting(t *testing.T) {
	ssid := "MyWiFi"
	pass := "s3cr3t"
	keyType := 3
	cmd := fmt.Sprintf(CmdWifiSetStainfo, ssid, pass, keyType)
	var js map[string]any
	if err := json.Unmarshal([]byte(cmd), &js); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
}

func TestCmdCountdownRulesFormatting(t *testing.T) {
	cmd := fmt.Sprintf(CmdAddCountdownRule, 120, 1, "NapTimer")
	var js map[string]any
	if err := json.Unmarshal([]byte(cmd), &js); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
}

func TestCmdSetRelayStateChildMultiEscaping(t *testing.T) {
	children := `"01","02","03"`
	state := 1
	cmd := fmt.Sprintf(CmdSetRelayStateChildMulti, children, state)
	var js map[string]any
	if err := json.Unmarshal([]byte(cmd), &js); err != nil {
		t.Fatalf("Invalid JSON for multi-child relay: %v", err)
	}
}
