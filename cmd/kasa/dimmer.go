package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var brightness = &cli.Command{
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
}
var getdimmer = &cli.Command{
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
}

var getalldimmer = &cli.Command{
	Name:  "getalldimmer",
	Usage: "get dimmer status for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		m, err := kasa.BroadcastDimmerParameters(int(cmd.Int("timeout")), int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		// ctx already canceled

		for k, v := range m {
			if v.ErrCode == 0 {
				kd, err := kasa.NewDevice(k)
				if err != nil {
					return err
				}
				s, err := kd.GetSettingsCtx(context.Background())
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
}
var setfadeontime = &cli.Command{
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
}
var setfadeofftime = &cli.Command{
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
}
var setgentleontime = &cli.Command{
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
}
var setgentleofftime = &cli.Command{
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
}
