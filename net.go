package kasa

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func (d *Device) sendTCP(cmd string) ([]byte, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: d.parsed, Port: 9999})
	if err != nil {
		klogger.Printf("Cannot connnect to device: %s", err.Error())
		return nil, err
	}
	defer conn.Close()
	// assume we are on the same LAN, one second is enough
	conn.SetReadDeadline(time.Now().Add(time.Second))

	// send the command with the uint32 "header"
	payload := ScrambleTCP(cmd)
	if _, err = conn.Write(payload); err != nil {
		klogger.Printf("Cannot send command to device: %s", err.Error())
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
	n, err = conn.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		err := fmt.Errorf("not all bytes read from host %s: %d/%d, %s", d.parsed, n, size, Unscramble(data))
		klogger.Printf(err.Error())
		return nil, err
	}

	return Unscramble(data), nil
}

func (d *Device) sendUDP(cmd string) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: d.parsed, Port: 9999})
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
