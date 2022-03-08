package kasa

import (
	"io"
	"net"
	"time"
)

// better would be to read the first 4 bytes, convert to uint32, allocate that much, then read the rest of the stream
func (d *Device) sendTCP(cmd string) ([]byte, error) {
	payload := encryptTCP(cmd)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: d.parsed, Port: 9999})
	if err != nil {
		klogger.Printf("Cannot connnect to device: %s", err.Error())
		return nil, err
	}
	defer conn.Close()

	if _, err = conn.Write(payload); err != nil {
		klogger.Printf("Cannot send command to device: %s", err.Error())
		return nil, err
	}

	blocksize := 1024
	bufsize := 10 * blocksize
	bytesread := 0
	data := make([]byte, 0, bufsize)
	tmp := make([]byte, blocksize)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))

		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			return nil, err
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
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: d.parsed, Port: 9999})
	if err != nil {
		return err
	}
	defer conn.Close()

	payload := encryptUDP(cmd)
	if _, err = conn.Write(payload); err != nil {
		return err
	}
	return nil
}
