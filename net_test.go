package kasa

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestDevice_SendTCP(t *testing.T) {
	// Start a local TCP server to mock a Kasa device
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Read header
		head := make([]byte, 4)
		_, _ = conn.Read(head)
		size := binary.BigEndian.Uint32(head)

		// Read body
		body := make([]byte, size)
		_, _ = conn.Read(body)

		// Respond with "ok" scrambled
		resp := ScrambleTCP("ok")
		_, _ = conn.Write(resp)
	}()

	addr := ln.Addr().(*net.TCPAddr)
	dev, _ := NewDeviceIP(addr.IP)
	dev.Port = addr.Port

	res, err := dev.sendTCP(context.Background(), "test")
	if err != nil {
		t.Fatalf("sendTCP failed: %v", err)
	}

	if string(res) != "ok" {
		t.Errorf("got %q, want %q", string(res), "ok")
	}
}

func TestDevice_SendUDP(t *testing.T) {
	// Start a local UDP server
	pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer pc.Close()

	addr := pc.LocalAddr().(*net.UDPAddr)
	dev, _ := NewDeviceIP(addr.IP)
	dev.Port = addr.Port

	err = dev.sendUDP(context.Background(), "test")
	if err != nil {
		t.Fatalf("sendUDP failed: %v", err)
	}
}

func TestDevice_SendTCP_Fail(t *testing.T) {
	// Use a non-existent port
	dev, _ := NewDeviceIP(net.ParseIP("127.0.0.1"))
	dev.Port = 1 // unlikely to be open

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := dev.sendTCP(ctx, "test")
	if err == nil {
		t.Error("expected error for closed port, got nil")
	}
}
