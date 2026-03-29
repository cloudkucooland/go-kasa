// Look ma! using chatgpt to write unit tests!
// of course it was only 90% of the way there, but what the heck

package kasa

import (
	"encoding/binary"
	"testing"
)

func TestScrambleUnscramble_RoundTrip(t *testing.T) {
	input := `{"system":{"get_sysinfo":{}}}`

	scrambled := Scramble(input)
	unscrambled := Unscramble(scrambled)

	if string(unscrambled) != string(input) {
		t.Fatalf("expected %s, got %s", input, unscrambled)
	}
}

func TestScrambleUnscramble_Empty(t *testing.T) {
	input := ""

	scrambled := Scramble(input)
	unscrambled := Unscramble(scrambled)

	if len(unscrambled) != 0 {
		t.Fatalf("expected empty result, got %v", unscrambled)
	}
}

func TestScramble_KnownValue(t *testing.T) {
	input := `{"system":{"get_sysinfo":{}}}`

	expected := []byte{
		208, 242, 129, 248, 139, 255, 154, 247, 213, 239, 148, 182, 209, 180, 192, 159, 236, 149, 230, 143, 225, 135, 232, 202, 240, 139, 246, 139, 246,
	}

	actual := Scramble(input)

	if len(actual) != len(expected) {
		t.Fatalf("length mismatch: got %d want %d", len(actual), len(expected))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("byte mismatch at %d: got %d want %d", i, actual[i], expected[i])
		}
	}
}

func TestUnscramble_InvalidData(t *testing.T) {
	// Not a valid scrambled payload
	input := []byte{0x00, 0x01, 0x02}

	output := Unscramble(input)

	// Should not panic and should return something deterministic
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func TestScramble_NotIdentity(t *testing.T) {
	input := "hello"

	scrambled := Scramble(input)

	if string(scrambled) == string(input) {
		t.Fatal("scramble should modify input")
	}
}

func TestScramble_HasLengthPrefix(t *testing.T) {
	input := `{"test":1}`

	out := ScrambleTCP(input)

	if len(out) < 4 {
		t.Fatal("expected length prefix")
	}

	payloadLen := int(out[0])<<24 | int(out[1])<<16 | int(out[2])<<8 | int(out[3])

	if payloadLen != len(out[4:]) {
		t.Fatalf("length prefix mismatch: got %d want %d", payloadLen, len(out[4:]))
	}
}

func TestScrambleUnscramble_Table(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "on"},
		{"json", `{"a":1}`},
		{"long", string(make([]byte, 256))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.input

			out := Unscramble(Scramble(in))

			if string(out) != tt.input {
				t.Fatalf("roundtrip failed: got %q want %q", out, tt.input)
			}
		})
	}
}

func TestScrambleAndUnscramble(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"empty string", ""},
		{"simple text", "hello world"},
		{"numbers", "1234567890"},
		{"special chars", "!@#$%^&*()_+-="},
		{"unicode", "こんにちは世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udp := Scramble(tt.text)
			if len(udp) != len(tt.text) {
				t.Fatalf("UDP scramble length mismatch, got %d, want %d", len(udp), len(tt.text))
			}

			tcp := ScrambleTCP(tt.text)
			if len(tcp) != len(tt.text)+4 {
				t.Fatalf("TCP scramble length mismatch, got %d, want %d", len(tcp), len(tt.text)+4)
			}

			// Check TCP header length
			headerLen := binary.BigEndian.Uint32(tcp[:4])
			if int(headerLen) != len(tt.text) {
				t.Fatalf("TCP header length mismatch, got %d, want %d", headerLen, len(tt.text))
			}

			// Unscramble UDP
			copyBuf := make([]byte, len(udp))
			copy(copyBuf, udp)
			res := Unscramble(copyBuf)
			if string(res) != tt.text {
				t.Fatalf("UDP unscramble mismatch, got %q, want %q", string(res), tt.text)
			}

			// Unscramble TCP (skip header)
			copyBufTCP := make([]byte, len(tcp)-4)
			copy(copyBufTCP, tcp[4:])
			resTCP := Unscramble(copyBufTCP)
			if string(resTCP) != tt.text {
				t.Fatalf("TCP unscramble mismatch, got %q, want %q", string(resTCP), tt.text)
			}
		})
	}
}

func TestScrambleReversibility(t *testing.T) {
	samples := []string{
		"simple",
		"another test string",
		"1234567890",
		"!@#$%^&*()",
		"",
	}

	for _, s := range samples {
		t.Run(s, func(t *testing.T) {
			udp := Scramble(s)
			unscrambled := Unscramble(append([]byte(nil), udp...))
			if string(unscrambled) != s {
				t.Errorf("UDP reversible failed: got %q, want %q", string(unscrambled), s)
			}

			tcp := ScrambleTCP(s)
			// skip 4-byte header
			unscrambledTCP := Unscramble(append([]byte(nil), tcp[4:]...))
			if string(unscrambledTCP) != s {
				t.Errorf("TCP reversible failed: got %q, want %q", string(unscrambledTCP), s)
			}
		})
	}
}

func TestScrambleTCPHeaderIntegrity(t *testing.T) {
	text := "testheader"
	tcp := ScrambleTCP(text)

	if len(tcp) < 4 {
		t.Fatal("TCP buffer too short")
	}

	headerLen := binary.BigEndian.Uint32(tcp[:4])
	if int(headerLen) != len(text) {
		t.Errorf("TCP header length mismatch, got %d, want %d", headerLen, len(text))
	}
}
