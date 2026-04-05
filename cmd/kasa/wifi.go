package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cloudkucooland/go-kasa"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var setwifi = &cli.Command{
	Name:      "setwifi",
	Usage:     "configure wifi",
	ArgsUsage: "host ssid key",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
		&cli.StringArg{Name: "ssid", Destination: &value},
		&cli.StringArg{Name: "key", Destination: &secondary},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		_, err = k.SetWIFICtx(ctx, value, secondary)
		return err
	},
}

var wifistatus = &cli.Command{
	Name:      "wifistatus",
	Usage:     "check device wifi status",
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		res, err := k.GetWIFIStatusCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Println(res) // make this prettier
		return nil
	},
}

var getallwifi = &cli.Command{
	Name:  "getallwifi",
	Usage: "get wifi stats for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		bctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmd.Int("timeout"))*time.Second)
		m, err := kasa.BroadcastWifiParameters(bctx, int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		defer cancel()

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintf(tabwrite, "Device\tIP\tSSID\t%s\t%s\n", color.GreenString("Key Type"), color.GreenString("RSSI"))
		for k, v := range m {
			kd, err := kasa.NewDevice(k)
			if err != nil {
				return err
			}
			s, err := kd.GetSettingsCtx(ctx)
			if err != nil {
				continue
			}

			fmt.Fprintf(tabwrite, "%s\t%s\t%s\t%s\t%s\n", s.Alias, k, v.SSID, keyType(v.KeyType), colorRSSI(v.RSSI))
		}
		tabwrite.Flush()
		return nil
	},
}

func colorRSSI(rssi int) string {
	s := fmt.Sprintf("%ddB", rssi)
	colored := ""
	switch {
	case rssi < -96.0:
		colored = color.RedString(s)
	case (rssi > -96.0 && rssi <= -85):
		colored = color.YellowString(s)
	default:
		colored = color.GreenString(s)
	}
	return colored
}

func keyType(t int) string {
	switch t {
	case 3:
		return color.GreenString("WPA3")
	case 2:
		return color.YellowString("WEP")
	default:
		return color.RedString("Unknown")
	}
}
