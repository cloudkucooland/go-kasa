package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var host string
var value string
var secondary string

func main() {
	cmd := &cli.Command{
		Name:      "kasa",
		Version:   "v0.3.0",
		Copyright: "(c) 2025 Scot Bontrager",
		Usage:     "control TP-Link kasa devices",
		UsageText: "kasa command",

		UseShortOptionHandling: true,
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
				Name:      "brightness",
				Usage:     "set brightness",
				UsageText: "kasa brightness host (value: 0-100)",
				ArgsUsage: "host brightness",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "brightness", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}

					b, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if b < 1 {
						b = 1
					}
					if b > 100 {
						b = 100
					}
					return k.SetBrightnessCtx(ctx, b)
				},
			},
			{
				Name:      "nocloud",
				Usage:     "disable the TP-Link cloud connection",
				UsageText: "kasa nocloud host",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.DisableCloudCtx(ctx)
				},
			},
			{
				Name:      "cloud",
				Usage:     "configure the TP-Link cloud connection",
				UsageText: "kasa cloud host username password",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.EnableCloudCtx(ctx, cmd.Args().Get(1), cmd.Args().Get(2))
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
			{
				Name:      "ledoff",
				Usage:     "disable status LED",
				ArgsUsage: "host true|false",
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
					return k.SetLEDOffCtx(ctx, b)
				},
			},
            discover,
			{
				Name:      "countdown",
				Usage:     "adjust device countdowns",
				ArgsUsage: "host [delete]",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "delete", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					if value == "delete" {
						return k.ClearCountdownRulesCtx(ctx)
					}
					rules, err := k.GetCountdownRulesCtx(ctx)
					if err != nil {
						return err
					}
					for _, r := range rules {
						fmt.Printf("%+v\n", r)
					}
					return nil
				},
			},
			{
				Name:      "reboot",
				Usage:     "reboot device",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.RebootCtx(ctx)
					// return nil
				},
			},
			{
				Name:      "alias",
				Usage:     "update device name (alias)",
				ArgsUsage: "host new-name",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "newname", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					child := cmd.String("child")
					if child != "" {
						return k.SetChildAliasCtx(ctx, child, value)
					}
					return k.SetAlias(value)
				},
			},
			{
				Name:      "raw",
				Usage:     "send raw command",
				ArgsUsage: "host command",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "command", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.SendRawCommandCtx(ctx, value)
				},
			},
            setwifi,
            wifistatus,
            getallwifi,
			{
				Name:      "addcountdown",
				Usage:     "add device countdown",
				ArgsUsage: "host duration True|False",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "duration", Destination: &value},
					&cli.StringArg{Name: "target", Destination: &secondary},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					dur, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if dur < 1 || dur > 3600 {
						return fmt.Errorf("invalid duration (1-3600)")
					}
					b, err := strconv.ParseBool(secondary)
					if err != nil {
						return err
					}
					return k.AddCountdownRuleCtx(ctx, dur, b, "auto")
				},
			},
			{
				Name:      "cleancountdown",
				Usage:     "remove countdown rules",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.ClearCountdownRulesCtx(ctx)
				},
			},
			{
				Name:      "getcountdown",
				Usage:     "view countdown rules",
				ArgsUsage: "host",
				Arguments: []cli.Argument{&cli.StringArg{Name: "host", Destination: &host}},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetCountdownRulesCtx(ctx)
					if err != nil {
						return err
					}
					fmt.Println(res)
					return nil
				},
			},
			{
				Name:      "getdimmer",
				Usage:     "check dimmer parameters",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetDimmerParametersCtx(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("%+v\n", res)
					return nil
				},
			},
			{
				Name:      "getrules",
				Usage:     "check running rules",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetRulesCtx(ctx)
					if err != nil {
						return err
					}
					fmt.Println(res)
					return nil
				},
			},
			{
				Name:      "setmode",
				Usage:     "set operating mode",
				ArgsUsage: "host mode",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "mode", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.SetModeCtx(ctx, value)
				},
			},
            getallemeter,
			{
				Name:  "getalldimmer",
				Usage: "get dimmer status for all devices",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := kasa.BroadcastDimmerParameters(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
					if err != nil {
						return err
					}

					for k, v := range m {
						if v.ErrCode == 0 {
							kd, err := kasa.NewDevice(k)
							if err != nil {
								return err
							}
							s, err := kd.GetSettingsCtx(ctx)
							if err != nil {
								return err
							}

							fmt.Printf("[%s] %s\n", k, s.Alias)
							fmt.Printf("Min Threshold: %d\t", v.MinThreshold)
							fmt.Printf("Fade On: %dms\t\t", v.FadeOnTime)
							fmt.Printf("Fade Off: %dms\n", v.FadeOffTime)
							fmt.Printf("Gentle On: %dms\t", v.GentleOnTime)
							fmt.Printf("Gentle Off: %dms\t", v.GentleOffTime)
							fmt.Printf("Ramp Rate: %dms\n", v.RampRate)
						}
					}
					return nil
				},
			},
			{
				Name:      "setfadeontime",
				Usage:     "set fade on time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "time", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					fade, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					return k.SetFadeOnTimeCtx(ctx, fade)
				},
			},
			{
				Name:      "setfadeofftime",
				Usage:     "set fade off time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "time", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					fade, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					return k.SetFadeOffTimeCtx(ctx, fade)
				},
			},
			{
				Name:      "setgentleontime",
				Usage:     "set gentle on time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "time", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					fade, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					return k.SetGentleOnTimeCtx(ctx, fade)
				},
			},
			{
				Name:      "setgentleofftime",
				Usage:     "set gentle off time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Destination: &host},
					&cli.StringArg{Name: "time", Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					fade, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					return k.SetGentleOffTimeCtx(ctx, fade)
				},
			},
            emeter,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmd.Int("timeout"))*time.Second)
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
