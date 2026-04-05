package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var host string
var value string
var secondary string

func main() {
	cmd := &cli.Command{
		Name:      "kasa",
		Version:   "v0.3.5",
		Copyright: "(c) 2026 Scot Bontrager",
		Usage:     "control TP-Link kasa devices",
		UsageText: "kasa command",

		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "child",
				Usage:   "child device ID",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
			},
			&cli.IntFlag{
				Name:    "repeats",
				Aliases: []string{"r"},
				Value:   1,
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   2,
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   9999,
				Usage:   "alternate port if using port-forwarding",
				Aliases: []string{"p"},
			},
		},

		Commands: []*cli.Command{
			discover,
			{
				Name:      "info",
				Usage:     "show basic info",
				UsageText: "kasa info host",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}

					s, err := k.GetSettingsCtx(ctx)
					if err != nil {
						return err
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
					return nil
				},
			},
			{
				Name:      "status",
				Usage:     "current device status",
				UsageText: "kasa status host",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					s, err := k.GetSettingsCtx(ctx)
					if err != nil {
						return err
					}
					if s.NumChildren > 0 {
						for _, v := range s.Children {
							fmt.Printf("[%s].[%s]:\t%d\n", s.Alias, v.Alias, v.RelayState)
						}
					} else {
						fmt.Printf("[%s]\tRelay:\t%d\tBrightness:\t%d%%\n", s.Alias, s.RelayState, s.Brightness)
					}
					return nil
				},
			},
			{
				Name:      "switch",
				Usage:     "toggle a relay's state",
				ArgsUsage: "[-c child ID] host true|false",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "state", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					b, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					child := cmd.String("child")
					if child != "" {
						return k.SetRelayStateChildCtx(ctx, child, b)
					}
					return k.SetRelayStateCtx(ctx, b)
				},
			},
			alias,
			getallemeter,
			emeter,
			getalldimmer,
			getdimmer,
			brightness,
			getallwifi,
			wifistatus,
			setwifi,
			setfadeontime,
			setfadeofftime,
			setgentleontime,
			setgentleofftime,
			reboot,
			nocloud,
			cloud,
			ledoff,
			addcountdown,
			cleancountdown,
			getcountdown,
			setmode,
			getlightsensorbrightness,
			getlightsensorconfig,
			raw,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func getKasaDevice(cmd *cli.Command) (*kasa.Device, error) {
	if host == "" {
		return nil, fmt.Errorf("missing host")
	}

	k, err := kasa.NewDevice(host)
	if err != nil {
		return nil, err
	}
	k.Port = int(cmd.Int("port"))
	return k, nil
}
