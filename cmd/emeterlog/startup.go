package main

import (
	"context"
	"time"

	"github.com/urfave/cli/v3"
)

var startup = &cli.Command{
	Name:  "startup",
	Usage: "start the daemon",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if err := setupdb(ctx, cmd); err != nil {
			return err
		}
		results := make(chan emeterdata, 100)
		go startDBWriter(ctx, results)

		ticker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-ctx.Done():
				close(results)
				return nil
			case <-ticker.C:
				// drop any lingering attempts before the next tick
				runCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
				if err := queryall(runCtx, results); err != nil {
					emlog.Error("query error", "err", err)
				}
			}
		}
		return nil
	},
}
