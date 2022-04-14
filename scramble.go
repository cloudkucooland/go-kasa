package kasa

import (
	"bytes"
	"encoding/binary"
)

// ScrambleTCP is for TCP, It writes the length in the first byte. It uses a binary buffer and writer.
func ScrambleTCP(plaintext string) []byte {
	var buf bytes.Buffer

	n := len(plaintext)
	buf.Grow(n + 4)

	// write the length as a 32-bit big-endian uint
	binary.Write(&buf, binary.BigEndian, uint32(n))

	key := byte(0xAB)
	for i := 0; i < n; i++ {
		key = plaintext[i] ^ key
		if err := buf.WriteByte(key); err != nil {
			klogger.Printf(err.Error())
			break
		}
	}

	return buf.Bytes()
}

// Scramble is simpler. UDP doesn't require the length header, just allocates and write to a slice.
func Scramble(plaintext string) []byte {
	n := len(plaintext)
	payload := make([]byte, n)

	key := byte(0xAB)
	for i := 0; i < n; i++ {
		payload[i] = plaintext[i] ^ key
		key = payload[i]
	}

	return payload
}

// Unscramble turns the response from the Kasa into parsable JSON
// it works in place -- be careful with your buffers
func Unscramble(ciphertext []byte) []byte {
	key := byte(0xAB)
	var nextKey byte

	for i := 0; i < len(ciphertext); i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}
	return ciphertext
}
