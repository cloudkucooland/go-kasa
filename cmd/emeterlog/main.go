package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// "github.com/cloudkucooland/go-kasa"

	"github.com/urfave/cli/v3"
)

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
