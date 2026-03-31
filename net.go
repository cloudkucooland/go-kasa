package kasa

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func (d *Device) sendTCP(ctx context.Context, cmd string) ([]byte, error) {
	if d.OverrideTCP != nil {
		return d.OverrideTCP(ctx, cmd)
	}

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

	// send the command with the uint32 "header"
	payload := ScrambleTCP(cmd)
	if _, err = conn.Write(payload); err != nil {
		klogger.Printf("cannot send command to device: %s", err.Error())
		return nil, err
	}

	// read the uint32 "header" to get the size of the rest of the block
	header := make([]byte, 4)
	n, err := conn.Read(header)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		err := fmt.Errorf("header not 32 bits (4 bytes): %d", n)
		klogger.Printf(err.Error())
		return nil, err
	}
	size := binary.BigEndian.Uint32(header)

	// read the entire rest of the block, then close the connection
	// we could leave the connection open and send subsequent requests
	// but for one-shot, this is enough
	data := make([]byte, size)
	totalread := 0
	for {
		n, err = conn.Read(data[totalread:])
		if err != nil {
			return nil, err
		}
		totalread = totalread + n

		if totalread >= int(size) {
			break
		}
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
