package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"
)

var emlog *slog.Logger

func main() {
	cmd := &cli.Command{
		Name:      "emeterlog",
		Version:   "v0.0.1",
		Copyright: "(c) 2026 Scot Bontrager",
		Usage:     "log kasa emeter data to influxdb",
		UsageText: "emeterlog",

		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
			},
		},

		Commands: []*cli.Command{
			startup,
		},
	}

	emlog = slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := cmd.Run(ctx, os.Args); err != nil {
		emlog.Error("error", "error", err.Error())
	}
}
