package kasa

import (
	"io"
	"net"
	"time"
)

func (d *Device) sendTCP(cmd string) (string, error) {
	payload := encrypt(cmd)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP:   d.parsed, Port: 9999})
	if err != nil {
		klogger.Printf("Cannot connnect to device: %s", err.Error())
		return "", err
	}
	_, err = conn.Write(payload)
	if err != nil {
		klogger.Printf("Cannot send command to device: %s", err.Error())
		return "", err
	}

	blocksize := 1024
	bufsize := 10 * blocksize
	bytesread := 0
	data := make([]byte, 0, bufsize)
	tmp := make([]byte, blocksize)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		defer conn.SetDeadline(time.Time{})

		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			return "", err
		}
		data = append(data, tmp[:n]...)
		bytesread += n
		if err == io.EOF || n != blocksize {
			break
		}
		// we read faster than the kasa fills its own buffers
		// 100 works some of the time, 150 seems better
		time.Sleep(time.Millisecond * 150)
	}

	result := decrypt(data[4:bytesread]) // start reading at 4, go to total bytes read
	return result, nil
}

func (d *Device) sendUDP(cmd string) error {
	payload := encryptUDP(cmd)
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: d.parsed, Port: 9999})
	if err != nil {
		return err
	}
	_, err = conn.Write(payload)
	if err != nil {
		return err
	}
	return nil
}
