package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"

	"github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "kasa",
		Version:   "v0.3.6",
		Copyright: "(c) 2026 Scot Bontrager",
		Usage:     "control TP-Link kasa devices",
		UsageText: "kasa command",

		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "child",
				Usage:   "child device ID",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
			},
			&cli.BoolFlag{
				Name:    "no-header",
				Aliases: []string{"n"},
			},
			&cli.IntFlag{
				Name:    "repeats",
				Aliases: []string{"r"},
				Value:   1,
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   2,
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   9999,
				Usage:   "alternate port if using port-forwarding",
				Aliases: []string{"p"},
			},
		},

		Commands: []*cli.Command{
			discover,
			{
				Name:      "info",
				Usage:     "show basic info",
				UsageText: "kasa info host",
				Before:    RequireDevice,
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k := ctx.Value("kasaDev").(*kasa.Device)

					s, err := k.GetSettingsCtx(ctx)
					if err != nil {
						return err
					}

					tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

					fmt.Fprintf(tabwrite, "Alias:\t%s\n", s.Alias)
					fmt.Fprintf(tabwrite, "DevName:\t%s\n", s.DevName)
					fmt.Fprintf(tabwrite, "Model:\t%s (%s)\n", s.Model, s.HWVersion)
					fmt.Fprintf(tabwrite, "Device ID:\t%s\n", s.DeviceID)
					fmt.Fprintf(tabwrite, "OEM ID:\t%s\n", s.OEMID)
					fmt.Fprintf(tabwrite, "Hardware ID:\t%s\n", s.HWID)
					fmt.Fprintf(tabwrite, "Software:\t%s\n", s.SWVersion)
					fmt.Fprintf(tabwrite, "MIC:\t%s\n", s.MIC)
					fmt.Fprintf(tabwrite, "MAC:\t%s\n", s.MAC)
					fmt.Fprintf(tabwrite, "LED Off:\t%d\n", s.LEDOff)
					fmt.Fprintf(tabwrite, "Active Mode:\t%s\n", s.ActiveMode)

					fmt.Fprintf(tabwrite, "Outlet\tRelay State\tBrightness\n")
					if s.NumChildren > 0 {
						for _, v := range s.Children {
							fmt.Fprintf(tabwrite, "%s\t%d\t\n", v.Alias, v.RelayState)
						}
					} else {
						fmt.Fprintf(tabwrite, "\t%d\t%d\n", s.RelayState, s.Brightness)
					}
					_ = tabwrite.Flush()
					return nil
				},
			},
			{
				Name:      "status",
				Usage:     "current device status",
				UsageText: "kasa status host",
				Before:    RequireDevice,
				ArgsUsage: "host",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "host"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					k := ctx.Value("kasaDev").(*kasa.Device)
					s, err := k.GetSettingsCtx(ctx)
					if err != nil {
						return err
					}
					tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
					fmt.Fprintf(tabwrite, "Device\tOutlet\tRelay State\tBrightness\n")
					if s.NumChildren > 0 {
						for _, v := range s.Children {
							fmt.Fprintf(tabwrite, "%s\t%s\t%d\t\n", s.Alias, v.Alias, v.RelayState)
						}
					} else {
						fmt.Fprintf(tabwrite, "%s\t\t%d\t%d\n", s.Alias, s.RelayState, s.Brightness)
					}
					_ = tabwrite.Flush()
					return nil
				},
			},
			{
				Name:      "switch",
				Usage:     "toggle a relay's state",
				Before:    RequireDevice,
				ArgsUsage: "[-c child ID] host true|false",
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
					child := cmd.String("child")
					if child != "" {
						return k.SetRelayStateChildCtx(ctx, child, b)
					}
					return k.SetRelayStateCtx(ctx, b)
				},
			},
			alias,
			emeter,
			allemeter,
			dimmer,
			alldimmer,
			brightness,
			wifi,
			allwifi,
			setwifi,
			setfadeontime,
			setfadeofftime,
			setgentleontime,
			setgentleofftime,
			reboot,
			nocloud,
			cloud,
			ledoff,
			addcountdown,
			cleancountdown,
			countdown,
			setmode,
			lightsensorbrightness,
			lightsensorconfig,
			raw,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := cmd.Run(ctx, os.Args); err != nil {
		if cmd.Bool("json") {
			status := map[string]any{"success": false, "error": err.Error()}
			_ = json.NewEncoder(os.Stdout).Encode(status)
		}
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	/* if cmd.Bool("json") {
		status := map[string]any{"success": true}
		json.NewEncoder(os.Stdout).Encode(status)
	} */
}

func RequireDevice(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	host := cmd.Args().Get(0)
	if host == "" {
		return ctx, fmt.Errorf("host argument is required for this command")
	}

	k, err := kasa.NewDevice(host)
	if err != nil {
		return ctx, fmt.Errorf("failed to initialize device: %w", err)
	}
	k.Port = int(cmd.Int("port"))

	return context.WithValue(ctx, "kasaDev", k), nil
}

func formatOutput(cmd *cli.Command, data any, pretty func()) error {
	if cmd.Bool("json") {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}
	pretty()
	return nil
}
