package kasa

import (
	"context"
)

type mockDevice struct {
	Device
	sendTCPFunc func(ctx context.Context, cmd string) ([]byte, error)
	sendUDPFunc func(ctx context.Context, cmd string) error
}

func (m *mockDevice) sendTCP(ctx context.Context, cmd string) ([]byte, error) {
	return m.sendTCPFunc(ctx, cmd)
}

func (m *mockDevice) sendUDP(ctx context.Context, cmd string) error {
	return m.sendUDPFunc(ctx, cmd)
}
