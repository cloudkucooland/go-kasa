package main

import (
	"context"
	"fmt"

	"github.com/cloudkucooland/go-kasa"

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
		m, err := kasa.BroadcastWifiParameters(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		for k, v := range m {
			nctx := context.Background()

			kd, err := kasa.NewDevice(k)
			if err != nil {
				return err
			}
			s, err := kd.GetSettingsCtx(nctx)
			if err != nil {
				continue
				// return err
			}

			fmt.Printf("[%s]\t", k)
			fmt.Printf("SSID: %s\t", v.SSID)
			fmt.Printf("Key Type: %d\t", v.KeyType)
			fmt.Printf("RSSI: %d\t", v.RSSI)
			fmt.Printf("%s\n", s.Alias)
		}
		return nil
	},
}
