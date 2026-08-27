package kasa

import (
	"encoding/json"
	"net"
	"testing"
)

func TestNewDevice(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"valid ip", "192.168.1.1", false},
		{"invalid ip/host", "invalid-host-name-that-should-not-exist", true},
		{"localhost", "127.0.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDevice(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDeviceIP(t *testing.T) {
	ip := net.ParseIP("192.168.1.1")
	dev, err := NewDeviceIP(ip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dev.IP.Equal(ip) {
		t.Errorf("got IP %v, want %v", dev.IP, ip)
	}
}

func TestDevice_Addr(t *testing.T) {
	dev, _ := NewDevice("192.168.1.1")
	dev.Port = 9999
	want := "192.168.1.1:9999"
	if got := dev.Addr(); got != want {
		t.Errorf("Addr() = %v, want %v", got, want)
	}
}

func TestKasaErr_OK(t *testing.T) {
	tests := []struct {
		name    string
		err     KasaErr
		wantErr bool
	}{
		{"no error", KasaErr{ErrCode: 0}, false},
		{"with error", KasaErr{ErrCode: -1, ErrMsg: "fail"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.OK(); (err != nil) != tt.wantErr {
				t.Errorf("OK() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiagnoseUnmarshaling(t *testing.T) {
	data := `{"smartlife.common.debug":{"get_diagnose_status":{"result":{"sysinfo":{"uptime":1654005,"free_mem":157472,"rbt_flag":0,"low":92512},"wireless":{"rssi":-47,"channel":1,"ssid":"IoT8417","bssid":"5A91E30B5CDE","auth_type":3,"cipherType":0,"reconn_count":17,"max_reconn_time":168,"total_disconnect_time":234},"cloud":{"send_fail":0,"hb_timeout":4,"recv_eof":7,"cloud_fsm":4,"acc_fsm":0}},"err_code":0}}}`
	var kd KasaDevice
	if err := json.Unmarshal([]byte(data), &kd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kd.Debug.Status.Result.Sysinfo.Uptime != 1654005 {
		t.Errorf("got Uptime %d, want 1654005", kd.Debug.Status.Result.Sysinfo.Uptime)
	}
	if kd.Debug.Status.Result.Wireless.RSSI != -47 {
		t.Errorf("got RSSI %d, want -47", kd.Debug.Status.Result.Wireless.RSSI)
	}
}

func TestOnboardingUnmarshaling(t *testing.T) {
	data := `{"system":{"get_onboarding_status":{"value":"pending","err_code":0}}}`
	var kd KasaDevice
	if err := json.Unmarshal([]byte(data), &kd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kd.GetSysinfo.Onboarding.Value != "pending" {
		t.Errorf("got Value %s, want pending", kd.GetSysinfo.Onboarding.Value)
	}
}
