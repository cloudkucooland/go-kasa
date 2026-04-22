package kasa

import (
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
