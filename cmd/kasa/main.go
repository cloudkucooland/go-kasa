package main

import (
	"flag"
	"fmt"
	"github.com/cloudkucooland/go-kasa"
	"strconv"
)

func main() {
	var command, host, value string
	// host := flag.String("host", "255.255.255.255", "host")
	// command := flag.String("command", "discover", "Param name")
	// value := flag.String("value", "", "Param value. Empty means only get")

	flag.Parse()
	args := flag.Args()
	argc := len(args)
	if argc == 0 {
		command = "unset"
	}
	if argc >= 1 {
		command = args[0]
	}
	if argc > 1 {
		host = args[1]
	}
	if argc > 2 {
		value = args[2]
	}

	switch command {
	case "args":
		fmt.Printf("command: %s ; host: %s ; value: %s\n", command, host, value)
	case "info":
		k, err := kasa.NewDevice(host)
		s, err := k.GetSettings()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", s)
	case "brightness":
		if host == "" || value == "" {
			fmt.Println("usage: kasa brightness [host] [1-100]")
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		b, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		if b < 1 {
			b = 1
		}
		if b > 100 {
			b = 100
		}
		err = k.SetBrightness(b)
		if err != nil {
			panic(err)
		}
	case "nocloud":
		if host == "" {
			fmt.Println("usage: kasa nocloud [host]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		err = k.DisableCloud()
		if err != nil {
			panic(err)
		}
	case "switch":
		if host == "" || value == "" {
			fmt.Println("usage: kasa switch [host] [true|false]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		err = k.SetRelayState(b)
		if err != nil {
			panic(err)
		}
	case "discover":
		if argc > 1 {
			fmt.Println("ignoring host, discover always broadcasts")
		}
		m, err := kasa.BroadcastDiscovery()
		if err != nil {
			panic(err)
		}
		fmt.Printf("found %d devices\n", len(m))
		for k, v := range m {
			fmt.Printf("%s: %s [state: %d] [brightness: %3d] %s\n", k, v.Model, v.RelayState, v.Brightness, v.Alias)
		}
	default:
		fmt.Println("usage: kasa [command] [host] [value]")
	}
}
