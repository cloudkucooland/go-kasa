package main

import (
	"context"
	"fmt"

	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

var getlightsensorbrightness = &cli.Command{
	Name:      "ambient",
	Usage:     "get ambient brightness",
	UsageText: "kasa ambient host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		b, err := k.GetCurrentBrightnessCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Ambient Brightness: %d", b)
		return nil
	},
}

var getlightsensorconfig = &cli.Command{
	Name:      "lightsensor",
	Usage:     "get light sensor config",
	UsageText: "kasa lightsensor host",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)

		c, err := k.GetLightSensorConfigCtx(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", c)
		return nil
	},
}
