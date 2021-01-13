package kasa

import (
	"fmt"
	"net"
)

func (d *Device) sendTCP(cmd string) (string, error) {
	payload := encrypt(cmd)
	r := net.TCPAddr{
		IP:   d.parsed,
		Port: 9999,
	}

	conn, err := net.DialTCP("tcp", nil, &r)
	if err != nil {
		fmt.Printf("Cannot connnect to device: %s", err.Error())
		return "", err
	}
	n, err := conn.Write(payload)
	if err != nil {
		fmt.Printf("Cannot send command to device: %s", err.Error())
		return "", err
	}

	data := make([]byte, 1024)
	n, err = conn.Read(data)
	if err != nil {
		fmt.Println("Cannot read data from device:", err)
		return "", err
	}
	result := decrypt(data[4:n]) // start reading at 4, go to total bytes read
	return result, nil
}

func (d *Device) sendUDP(cmd string) error {
	// we don't need ListenUDP, DialUDP would be better
	payload := encryptUDP(cmd)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	_, err = conn.WriteToUDP(payload, &net.UDPAddr{IP: d.parsed, Port: 9999})
	if err != nil {
		return err
	}
	return nil
}

func broadcastCmd(cmd string) error {
	// get interface subnets and real local broadcast addresses
	d := Device{IP: "255.255.255.255"}
	return d.sendUDP(cmd)
}
