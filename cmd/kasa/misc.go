package main

import (
	"context"
	"fmt"
	"os"
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
		b, err := strconv.ParseBool(cmd.StringArg("state"))
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

		nn := cmd.StringArg("newname")
		if nn == "" {
			return fmt.Errorf("need a valid name")
		}

		child := cmd.String("child")
		if child != "" {
			fmt.Fprintf(os.Stderr, "using child %s", child)
			return k.SetChildAliasCtx(ctx, child, nn)
		}

		return k.SetAlias(nn)
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
		b, err := k.SendRawCommandCtx(ctx, cmd.StringArg("command"))
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}
