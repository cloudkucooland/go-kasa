package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var brightness = &cli.Command{
	Name:      "brightness",
	Usage:     "set brightness",
	UsageText: "kasa brightness host (value: 0-100)",
	ArgsUsage: "host brightness",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "brightness"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)

		b := cmd.IntArg("brightness")
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
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.StringArg("host") == "" {
			return getalldimmer.Action(ctx, cmd)
		}

		RequireDevice(ctx, cmd)
		k := ctx.Value("kasaDev").(*kasa.Device)
		res, err := k.GetDimmerParametersCtx(ctx)
		if err != nil {
			return err
		}

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabwrite, "Device\tIP\tMin\tFade On\tFade Off\tGentle On\tGentle Off\tRamp Rate\n")
		fmt.Fprintf(tabwrite, "%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n", "", k.IP, res.MinThreshold, res.FadeOnTime, res.FadeOffTime, res.GentleOnTime, res.GentleOffTime, res.RampRate)
		tabwrite.Flush()
		return nil
	},
}

var getalldimmer = &cli.Command{
	Name:  "getalldimmer",
	Usage: "get dimmer status for all devices",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		bctx, cancel := context.WithTimeout(ctx, time.Duration(cmd.Int("timeout"))*time.Second)
		m, err := kasa.BroadcastDimmerParameters(bctx, int(cmd.Int("repeats")))
		if err != nil {
			return err
		}
		defer cancel()

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabwrite, "Device\tIP\tMin\tFade On\tFade Off\tGentle On\tGentle Off\tRamp Rate\n")

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

				fmt.Fprintf(tabwrite, "%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n", s.Alias, k, v.MinThreshold, v.FadeOnTime, v.FadeOffTime, v.GentleOnTime, v.GentleOffTime, v.RampRate)
			}
		}
		tabwrite.Flush()
		return nil
	},
}

var setfadeontime = &cli.Command{
	Name:      "setfadeontime",
	Usage:     "set fade on time",
	ArgsUsage: "time in ms",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "time"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetFadeOnTimeCtx(ctx, cmd.IntArg("time"))
	},
}

var setfadeofftime = &cli.Command{
	Name:      "setfadeofftime",
	Usage:     "set fade off time",
	ArgsUsage: "time in ms",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "time"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetFadeOffTimeCtx(ctx, cmd.IntArg("time"))
	},
}

var setgentleontime = &cli.Command{
	Name:      "setgentleontime",
	Usage:     "set gentle on time",
	ArgsUsage: "time in ms",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "time"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetGentleOnTimeCtx(ctx, cmd.IntArg("time"))
	},
}

var setgentleofftime = &cli.Command{
	Name:      "setgentleofftime",
	Usage:     "set gentle off time",
	ArgsUsage: "time in ms",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "time"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetGentleOffTimeCtx(ctx, cmd.IntArg("time"))
	},
}
