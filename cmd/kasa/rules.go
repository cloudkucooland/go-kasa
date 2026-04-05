package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
)

var addcountdown = &cli.Command{
	Name:      "addcountdown",
	Usage:     "add device countdown",
	ArgsUsage: "host duration True|False",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
		&cli.StringArg{Name: "duration", Destination: &value},
		&cli.StringArg{Name: "target", Destination: &secondary},
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
		return k.AddCountdownRuleCtx(ctx, dur, b, "auto")
	},
}

var cleancountdown = &cli.Command{
	Name:      "cleancountdown",
	Usage:     "remove countdown rules",
	ArgsUsage: "host",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		return k.ClearCountdownRulesCtx(ctx)
	},
}

var getcountdown = &cli.Command{
	Name:      "getcountdown",
	Usage:     "view countdown rules",
	ArgsUsage: "host",
	Arguments: []cli.Argument{&cli.StringArg{Name: "host", Destination: &host}},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		res, err := k.GetCountdownRulesCtx(ctx)
		if err != nil {
			return err
		}

		tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tabwrite, "ID\tName\tEnable\tDelay\tActive\tRemaining\n")
		for _, r := range res {
			fmt.Fprintf(tabwrite, "%s\t%s\t%s\t%d\t%d\t%d\n", "", r.ID, r.Name, r.Enable, r.Delay, r.Active, r.Remaining)
		}
		tabwrite.Flush()
		return nil
	},
}

var setmode = &cli.Command{
	Name:      "setmode",
	Usage:     "set operating mode",
	ArgsUsage: "host mode",
	Arguments: []cli.Argument{
		&cli.StringArg{Name: "host", Destination: &host},
		&cli.StringArg{Name: "mode", Destination: &value},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		k, err := getKasaDevice(cmd)
		if err != nil {
			return err
		}
		return k.SetModeCtx(ctx, value)
	},
}
