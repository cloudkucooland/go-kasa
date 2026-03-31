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

	err := md.SendRawCommandCtx(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
