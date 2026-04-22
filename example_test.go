package kasa_test

import (
	"context"
	"fmt"
	"github.com/cloudkucooland/go-kasa"
)

func ExampleNewDevice() {
	// Initialize a new device with its IP address
	dev, err := kasa.NewDevice("192.168.1.100")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Device IP: %s\n", dev.IP.String())
	// Output: Device IP: 192.168.1.100
}

func ExampleDevice_SetRelayState() {
	dev, _ := kasa.NewDevice("192.168.1.100")
	dev.OverrideUDP = func(ctx context.Context, cmd string) error {
		fmt.Printf("Command: %s\n", cmd)
		return nil
	}

	// Wrapper for SetRelayStateCtx(context.Background(), true)
	_ = dev.SetRelayState(true)
	// Output: Command: {"system":{"set_relay_state":{"state":1}}}
}

func ExampleDevice_SetRelayStateCtx() {
	// This example uses a mock to avoid actual network calls
	dev, _ := kasa.NewDevice("192.168.1.100")
	dev.OverrideUDP = func(ctx context.Context, cmd string) error {
		// In a real scenario, this would send a UDP packet to the device
		fmt.Printf("Sending command: %s\n", cmd)
		return nil
	}

	err := dev.SetRelayStateCtx(context.Background(), true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output: Sending command: {"system":{"set_relay_state":{"state":1}}}
}

func ExampleDevice_GetSettings() {
	dev, _ := kasa.NewDevice("192.168.1.100")
	dev.OverrideTCP = func(ctx context.Context, cmd string) ([]byte, error) {
		return []byte(`{"system":{"get_sysinfo":{"alias":"Smart Plug","err_code":0}}}`), nil
	}

	settings, _ := dev.GetSettings()
	fmt.Println(settings.Alias)
	// Output: Smart Plug
}

func ExampleDevice_GetSettingsCtx() {
	dev, _ := kasa.NewDevice("192.168.1.100")
	dev.OverrideTCP = func(ctx context.Context, cmd string) ([]byte, error) {
		return []byte(`{"system":{"get_sysinfo":{"alias":"My Smart Plug","model":"HS100","err_code":0}}}`), nil
	}

	settings, err := dev.GetSettingsCtx(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Device Alias: %s\n", settings.Alias)
	fmt.Printf("Device Model: %s\n", settings.Model)
	// Output:
	// Device Alias: My Smart Plug
	// Device Model: HS100
}

func ExampleDevice_GetEmeter() {
	dev, _ := kasa.NewDevice("192.168.1.100")
	dev.OverrideTCP = func(ctx context.Context, cmd string) ([]byte, error) {
		return []byte(`{"emeter":{"get_realtime":{"power_mw":1500,"err_code":0}}}`), nil
	}

	emeter, _ := dev.GetEmeter()
	fmt.Printf("Power: %d mW\n", emeter.PowerMW)
	// Output: Power: 1500 mW
}
