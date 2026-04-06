package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/cloudkucooland/go-kasa"
	"github.com/urfave/cli/v3"
)

var addcountdown = &cli.Command{
	Name:      "addcountdown",
	Usage:     "add device countdown",
	Before:    RequireDevice,
	ArgsUsage: "host duration True|False",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.IntArg{Name: "duration"},
		&cli.StringArg{Name: "target"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		dur := cmd.Int("duration")
		if dur < 1 || dur > 3600 {
			return fmt.Errorf("invalid duration (1-3600)")
		}
		b, err := strconv.ParseBool(cmd.String("target"))
		if err != nil {
			return err
		}
		return k.AddCountdownRuleCtx(ctx, dur, b, "auto")
	},
}

var cleancountdown = &cli.Command{
	Name:      "cleancountdown",
	Usage:     "remove countdown rules",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.ClearCountdownRulesCtx(ctx)
	},
}

var countdown = &cli.Command{
	Name:      "getcountdown",
	Usage:     "view countdown rules",
	Before:    RequireDevice,
	ArgsUsage: "host",
	Arguments: []cli.Argument{&cli.StringArg{Name: "host"}},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		res, err := k.GetCountdownRulesCtx(ctx)
		if err != nil {
			return err
		}

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tabwrite, "ID\tName\tEnable\tDelay\tActive\tRemaining\n")
		for _, r := range res {
			fmt.Fprintf(tabwrite, "%s\t%s\t%d\t%d\t%d\t%d\n", r.ID, r.Name, r.Enable, r.Delay, r.Active, r.Remaining)
		}
		_ = tabwrite.Flush()
		return nil
	},
}

var setmode = &cli.Command{
	Name:      "setmode",
	Usage:     "set operating mode",
	Before:    RequireDevice,
	ArgsUsage: "host mode",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host"},
		&cli.StringArg{Name: "mode"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k := ctx.Value("kasaDev").(*kasa.Device)
		return k.SetModeCtx(ctx, cmd.String("mode"))
	},
}
