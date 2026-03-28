// Look ma! using chatgpt to write unit tests!
// of course it was only 90% of the way there, but what the heck

package kasa_test

import (
	// "fmt"
	"testing"

	"github.com/cloudkucooland/go-kasa"
)

func TestScrambleUnscramble_RoundTrip(t *testing.T) {
	input := `{"system":{"get_sysinfo":{}}}`

	scrambled := kasa.Scramble(input)
	unscrambled := kasa.Unscramble(scrambled)

	if string(unscrambled) != string(input) {
		t.Fatalf("expected %s, got %s", input, unscrambled)
	}
}

func TestScrambleUnscramble_Empty(t *testing.T) {
	input := ""

	scrambled := kasa.Scramble(input)
	unscrambled := kasa.Unscramble(scrambled)

	if len(unscrambled) != 0 {
		t.Fatalf("expected empty result, got %v", unscrambled)
	}
}

func TestScramble_KnownValue(t *testing.T) {
	input := `{"system":{"get_sysinfo":{}}}`

	expected := []byte{
		208, 242, 129, 248, 139, 255, 154, 247, 213, 239, 148, 182, 209, 180, 192, 159, 236, 149, 230, 143, 225, 135, 232, 202, 240, 139, 246, 139, 246,
	}

	actual := kasa.Scramble(input)

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

	output := kasa.Unscramble(input)

	// Should not panic and should return something deterministic
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func TestScramble_NotIdentity(t *testing.T) {
	input := "hello"

	scrambled := kasa.Scramble(input)

	if string(scrambled) == string(input) {
		t.Fatal("scramble should modify input")
	}
}

func TestScramble_HasLengthPrefix(t *testing.T) {
	input := `{"test":1}`

	out := kasa.ScrambleTCP(input)

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

			out := kasa.Unscramble(kasa.Scramble(in))

			if string(out) != tt.input {
				t.Fatalf("roundtrip failed: got %q want %q", out, tt.input)
			}
		})
	}
}
