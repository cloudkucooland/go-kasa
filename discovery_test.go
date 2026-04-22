package kasa

import (
	"testing"
)

func TestBroadcastAddresses(t *testing.T) {
	addrs, err := BroadcastAddresses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should at least find loopback if nothing else, but netutil.go skips loopback.
	// Most systems running tests will have at least one active interface.
	if len(addrs) == 0 {
		t.Log("no broadcast addresses found (normal in some CI environments)")
	}

	for _, addr := range addrs {
		if addr.To4() == nil {
			t.Errorf("found non-IPv4 broadcast address: %v", addr)
		}
	}
}
