package kasa

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func (d *Device) sendTCP(ctx context.Context, cmd string) ([]byte, error) {
	if d.OverrideTCP != nil {
		return d.OverrideTCP(ctx, cmd)
	}

	// is this needed if we do it on the connection based on ctx?
	dialer := &net.Dialer{
		Timeout:  1 * time.Second,
		Deadline: time.Now().Add(2 * time.Second),
	}

	conn, err := dialer.DialContext(ctx, "tcp4", d.Addr())
	if err != nil {
		klogger.Printf("cannot connnect to device: %s", err.Error())
		return nil, err
	}
	defer conn.Close()

	if d, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(d)
	}

	// send the command with the uint32 "header"
	payload := ScrambleTCP(cmd)
	if _, err = conn.Write(payload); err != nil {
		klogger.Printf("cannot send command to device: %s", err.Error())
		return nil, err
	}

	// read the uint32 "header" to get the size of the rest of the block
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	size := binary.BigEndian.Uint32(header)

	data := make([]byte, size)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return Unscramble(data), nil
}

func (d *Device) sendUDP(ctx context.Context, cmd string) error {
	if d.OverrideUDP != nil {
		return d.OverrideUDP(ctx, cmd)
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: d.IP, Port: d.Port})
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := Scramble(cmd)
	if _, err = conn.Write(payload); err != nil {
		return err
	}
	return nil
}
