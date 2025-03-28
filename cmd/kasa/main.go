package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var host string
var value string
var secondary string

func main() {
	cmd := &cli.Command{
		Name:      "kasa",
		Version:   "v0.2.0",
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}

					s, err := k.GetSettings()
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					s, err := k.GetSettings()
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "brightness", Min: 1, Max: 1, Destination: &value},
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
					return k.SetBrightness(b)
				},
			},
			{
				Name:      "nocloud",
				Usage:     "disable the TP-Link cloud connection",
				UsageText: "kasa nocloud host",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.DisableCloud()
				},
			},
			{
				Name:      "cloud",
				Usage:     "configure the TP-Link cloud connection",
				UsageText: "kasa cloud host username password",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.EnableCloud(cmd.Args().Get(1), cmd.Args().Get(2))
				},
			},
			{
				Name:      "switch",
				Usage:     "toggle a relay's state",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "state", Min: 1, Max: 1, Destination: &value},
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
						return k.SetRelayStateChild(child, b)
					}
					return k.SetRelayState(b)
				},
			},
			{
				Name:      "ledoff",
				Usage:     "disable status LED",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "state", Min: 1, Max: 1, Destination: &value},
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
					return k.SetLEDOff(b)
				},
			},
			{
				Name:  "discover",
				Usage: "discover local devices",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := kasa.BroadcastDiscovery(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
					if err != nil {
						return err
					}

					keys := make([]string, 0, len(m))
					for key := range m {
						keys = append(keys, key)
					}
					sort.Strings(keys)

					fmt.Printf("found %d devices\n", len(m))
					for _, k := range keys {
						v := m[k]
						if len(v.Children) == 0 {
							fmt.Printf("%15s: %s %32s [state: %d] [brightness: %3d]\n", k, v.Model, v.Alias, v.RelayState, v.Brightness)
						} else {
							fmt.Printf("%15s: %s %s\n", k, v.Model, v.Alias)
							for _, c := range v.Children {
								fmt.Printf("    ID: %40s%s %26s [state: %d]\n", v.DeviceID, c.ID, c.Alias, c.RelayState)
							}
						}
					}
					return nil
				},
			},
			{
				Name:      "countdown",
				Usage:     "adjust device countdowns",
				ArgsUsage: "host [delete]",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "delete", Min: 0, Max: 1, Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					if value == "delete" {
						return k.ClearCountdownRules()
					}
					rules, err := k.GetCountdownRules()
					if err != nil {
						return err
					}
					for _, r := range *rules {
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.Reboot()
					// return nil
				},
			},
			{
				Name:      "alias",
				Usage:     "update device name (alias)",
				ArgsUsage: "host new-name",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "newname", Min: 1, Max: 1, Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					child := cmd.String("child")
					if child != "" {
						return k.SetChildAlias(child, value)
					}
					return k.SetAlias(value)
				},
			},
			{
				Name:      "raw",
				Usage:     "send raw command",
				ArgsUsage: "host command",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "command", Min: 1, Max: 1, Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.SendRawCommand(value)
				},
			},
			{
				Name:      "setwifi",
				Usage:     "configure wifi",
				ArgsUsage: "host ssid key",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "ssid", Min: 1, Max: 1, Destination: &value},
					&cli.StringArg{Name: "key", Min: 1, Max: 1, Destination: &secondary},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					_, err = k.SetWIFI(value, secondary)
					return err
				},
			},
			{
				Name:      "wifistatus",
				Usage:     "check device wifi status",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetWIFIStatus()
					if err != nil {
						return err
					}
					fmt.Println(res) // make this prettier
					return nil
				},
			},
			{
				Name:      "addcountdown",
				Usage:     "add device countdown",
				ArgsUsage: "host duration True|False",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "duration", Min: 1, Max: 1, Destination: &value},
					&cli.StringArg{Name: "target", Min: 1, Max: 1, Destination: &secondary},
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
					return k.AddCountdownRule(dur, b, "auto")
				},
			},
			{
				Name:      "cleancountdown",
				Usage:     "remove coundown rules",
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.ClearCountdownRules()
				},
			},
			{
				Name:      "getcountdown",
				Usage:     "view countdown rules",
				ArgsUsage: "host",
				Arguments: []cli.Argument{&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host}},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetCountdownRules()
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetDimmerParameters()
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					res, err := k.GetRules()
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "mode", Min: 1, Max: 1, Destination: &value},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					return k.SetMode(value)
				},
			},
			{
				Name:  "getallwifi",
				Usage: "get wifi stats for all devices",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := kasa.BroadcastWifiParameters(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
					if err != nil {
						return err
					}
					for k, v := range m {
						fmt.Printf("%s: %+v\n", k, v)
					}
					return nil
				},
			},
			{
				Name:  "getallemeter",
				Usage: "get emeter stats for all devices",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := kasa.BroadcastEmeter(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
					if err != nil {
						return err
					}
					for k, v := range m {
						fmt.Printf("%s slot %d\n", k, v.Emeter.Realtime.Slot)
						fmt.Printf("Current:\t%dmA\n", v.Emeter.Realtime.CurrentMA)
						fmt.Printf("Voltage:\t%2.2fV\n", float64(v.Emeter.Realtime.VoltageMV)/1000)
						fmt.Printf("Power:\t\t%2.2fW\n", float64(v.Emeter.Realtime.PowerMW)/1000)
						fmt.Printf("Total:\t\t%2.2fkWh\n", float64(v.Emeter.Realtime.TotalWH)/1000)
					}
					return nil
				},
			},
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
							fmt.Printf("%s: %+v\n", k, v)
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
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "time", Min: 1, Max: 1, Destination: &value},
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
					return k.SetFadeOnTime(fade)
				},
			},
			{
				Name:      "setfadeoofftime",
				Usage:     "set fade off time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "time", Min: 1, Max: 1, Destination: &value},
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
					return k.SetFadeOffTime(fade)
				},
			},
			{
				Name:      "setgentleontime",
				Usage:     "set gentle on time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "time", Min: 1, Max: 1, Destination: &value},
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
					return k.SetGentleOnTime(fade)
				},
			},
			{
				Name:      "setgentleoofftime",
				Usage:     "set gentle off time",
				ArgsUsage: "time in ms",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "time", Min: 1, Max: 1, Destination: &value},
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
					return k.SetGentleOffTime(fade)
				},
			},
			{
				Name:      "emeter",
				Usage:     "check energy usage",
				ArgsUsage: "host month year",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host", Min: 1, Max: 1, Destination: &host},
					&cli.StringArg{Name: "month", Min: 0, Max: 1, Destination: &value},
					&cli.StringArg{Name: "year", Min: 0, Max: 1, Destination: &secondary},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k, err := getKasaDevice(cmd)
					if err != nil {
						return err
					}
					month := 0
					year := 0

					if value != "" {
						month, err = strconv.Atoi(value)
						if err != nil {
							return err
						}
						if month < 1 || month > 12 {
							return fmt.Errorf("invalid month")
						}
					}

					yy := cmd.String("year")
					if yy != "" {
						year, err = strconv.Atoi(yy)
						if err != nil {
							return err
						}
					}

					if month == 0 {
						var em *kasa.EmeterRealtime

						if cmd.String("child") != "" {
							em, err = k.GetEmeterChild(cmd.String("child"))
						} else {
							em, err = k.GetEmeter()
						}

						if err != nil {
							return err
						}

						fmt.Printf("Current:\t%dmA\n", em.CurrentMA)
						fmt.Printf("Voltage:\t%2.2fV\n", float64(em.VoltageMV)/1000)
						fmt.Printf("Power:\t\t%2.2fW\n", float64(em.PowerMW)/1000)
						fmt.Printf("Total:\t\t%2.2fkWh\n", float64(em.TotalWH)/1000)
						return nil
					}

					if year == 0 {
						year = 2025 // make this auto-determine the current year
					}
					// get month/year date range
					em, err := k.GetEmeterMonth(month, year)
					if err != nil {
						return err
					}
					for _, v := range em.List {
						fmt.Printf("%d-%02d-%02d:\t%2.2fkWh\n", v.Year, v.Month, v.Day, float64(v.WH)/1000)
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
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
	k.Debug = cmd.Bool("debug")
	k.Port = int(cmd.Int("port"))
	return k, nil
}
