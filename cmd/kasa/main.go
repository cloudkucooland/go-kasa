package main

import (
	"flag"
	"fmt"
	"github.com/cloudkucooland/go-kasa"
	"sort"
	"strconv"
)

func main() {
	var command, host, value string
	repeats := flag.Int("r", 1, "UDP repeats")
	timeout := flag.Int("t", 2, "timeout")
	child := flag.String("c", "", "child")

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
		if err != nil {
			panic(err)
		}
		s, err := k.GetSettings()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Alias:\t\t%s\n", s.Alias)
		fmt.Printf("DevName:\t%s\n", s.DevName)
		fmt.Printf("Model:\t\t%s [%s]\n", s.Model, s.HWVersion)
		fmt.Printf("Device ID:\t%s\n", s.DeviceID)
		fmt.Printf("OEM ID:\t\t%s\n", s.OEMID)
		fmt.Printf("Hardware ID:\t%s\n", s.HWID)
		fmt.Printf("Software:\t%s\n", s.SWVersion)
		fmt.Printf("MIC:\t\t%s\n", s.MIC)
		fmt.Printf("MAC:\t\t%s\n", s.MAC)
		fmt.Printf("LED Off:\t%d\n", s.LEDOff)
		fmt.Printf("Active Mode:\t%s\n", s.ActiveMode)
        if s.NumChildren > 0 {
            for _, v := range s.Children {
		        fmt.Printf("Outlet [%s]:\t%d\t\t(%s)\n", v.Alias, v.RelayState, v.ID)
            }
        } else {
		    fmt.Printf("Relay:\t%d\tBrightness:\t%d%%\n", s.RelayState, s.Brightness)
        }
	case "status":
		if host == "" {
			fmt.Println("usage: kasa status [host]")
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		s, err := k.GetSettings()
		if err != nil {
			panic(err)
		}
        if s.NumChildren > 0 {
            for _, v := range s.Children {
		        fmt.Printf("[%s].[%s]:\t%d\n", s.Alias, v.Alias, v.RelayState)
            }
        } else {
		    fmt.Printf("[%s]\tRelay:\t%d\tBrightness:\t%d%%\n", s.Alias, s.RelayState, s.Brightness)
        }
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
        if *child != "" {
            k.SetRelayStateChild(*child, b)
        } else {
		    err = k.SetRelayState(b)
        }
		if err != nil {
			panic(err)
		}
	case "ledoff":
		if host == "" || value == "" {
			fmt.Println("usage: kasa ledoff [host] [true|false]")
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
		err = k.SetLEDOff(b)
		if err != nil {
			panic(err)
		}
	case "discover":
		if argc > 1 {
			fmt.Println("ignoring host, discover always broadcasts")
		}
		m, err := kasa.BroadcastDiscovery(*timeout, *repeats)
		if err != nil {
			panic(err)
		}

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Printf("found %d devices\n", len(m))
		for _, k := range keys {
			v := m[k]
			fmt.Printf("%s: %s [state: %d] [brightness: %3d] %s\n", k, v.Model, v.RelayState, v.Brightness, v.Alias)
		}
	case "reboot":
		if host == "" {
			fmt.Println("usage: kasa reboot [host]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		err = k.Reboot()
		if err != nil {
			panic(err)
		}
	case "alias":
		if host == "" || value == "" {
			fmt.Println("usage: kasa alias [host] [NewName]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		err = k.SetAlias(value)
		if err != nil {
			panic(err)
		}
	case "wifistatus":
		if host == "" {
			fmt.Println("usage: kasa wifistatus [host]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		res, err := k.GetWIFIStatus()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	case "getdimmer":
		if host == "" {
			fmt.Println("usage: kasa getdimmer [host]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		res, err := k.GetDimmerParameters()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	case "getalldimmer":
		m, err := kasa.BroadcastDimmerParameters(*timeout, *repeats)
		if err != nil {
			panic(err)
		}
		for k, v := range *m {
			if v.ErrCode == 0 {
				fmt.Printf("%s: %+v\n", k, v)
			}
		}
	case "getsched":
		if host == "" {
			fmt.Println("usage: kasa getsched [host]")
			return
		}
		k, err := kasa.NewDevice(host)
		if err != nil {
			panic(err)
		}
		res, err := k.GetRules()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	default:
		fmt.Println("Valid commands: info, status, brightness, nocloud, switch, ledoff, discover, reboot, alias, wifistatus, getdimmer, getsched")
	}
}
