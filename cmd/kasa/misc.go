package main

import (
	"context"
	"strconv"

	// "github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

var nocloud = &cli.Command{
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
}

var cloud = &cli.Command{
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
}

var ledoff = &cli.Command{
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
}

var reboot = &cli.Command{
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
}

var alias = &cli.Command{
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
}

var raw = &cli.Command{
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
}
