package main

import (
	"context"
	"strconv"

	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

var nocloud = &cli.Command{
	Name:      "nocloud",
	Usage:     "disable the TP-Link cloud connection",
	UsageText: "kasa nocloud host",
	ArgsUsage: "host",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.DisableCloudCtx(ctx)
	},
}

var cloud = &cli.Command{
	Name:      "cloud",
	Usage:     "configure the TP-Link cloud connection",
	UsageText: "kasa cloud host username password",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.EnableCloudCtx(ctx, cmd.Args().Get(1), cmd.Args().Get(2))
	},
}

var ledoff = &cli.Command{
	Name:      "ledoff",
	Usage:     "disable status LED",
	ArgsUsage: "host true|false",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "state"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		b, err := strconv.ParseBool(cmd.String("state"))
		if err != nil {
			return err
		}
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetLEDOffCtx(ctx, b)
	},
}

var reboot = &cli.Command{
	Name:      "reboot",
	Usage:     "reboot device",
	ArgsUsage: "host",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.RebootCtx(ctx)
	},
}

var alias = &cli.Command{
	Name:      "alias",
	Usage:     "update device name (alias)",
	ArgsUsage: "host new-name",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "newname"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		child := cmd.String("child")
		if child != "" {
			return k.SetChildAliasCtx(ctx, child, cmd.String("newname"))
		}
		return k.SetAlias(cmd.String("newname"))
	},
}

var raw = &cli.Command{
	Name:      "raw",
	Usage:     "send raw command",
	ArgsUsage: "host command",
	Before:    RequireDevice,
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "command"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SendRawCommandCtx(ctx, cmd.String("command"))
	},
}
